import { ipcMain, shell, IpcMainEvent, BrowserWindow, app } from 'electron'
import Constants from './utils/Constants'
import Store from 'electron-store'
import _ from 'lodash'
import { encrypt, decrypt } from './utils/encryption'
import { machine } from 'node-unique-machine-id'
import axios from 'axios'
import crypto from 'crypto'
const { autoUpdater } = require('electron-updater')
const log = require('electron-log')

// import { generateMachineFingerprint } from './utils/fingerprint'
const store = new Store()
const { snakeCase } = _

// Cache for internet connectivity status to avoid excessive requests
let internetConnectionCache: { status: boolean; timestamp: number } | null =
  null
const CONNECTIVITY_CACHE_DURATION = 30000 // 30 seconds

const toSnakeCaseKeys = (obj) => {
  if (typeof obj !== 'object' || obj === null) return obj
  if (Array.isArray(obj)) return obj.map(toSnakeCaseKeys)
  return Object.keys(obj).reduce((result, key) => {
    result[snakeCase(key)] = toSnakeCaseKeys(obj[key])
    return result
  }, {})
}

const decode = (token) => {
  const decodedString = decodeURIComponent(
    atob(token.split('.')[1].replace('-', '+').replace('_', '/'))
      .split('')
      .map((c) => `%${('00' + c.charCodeAt(0).toString(16)).slice(-2)}`)
      .join('')
  )

  const decodedObject = JSON.parse(decodedString)

  return toSnakeCaseKeys(decodedObject)
}

// --- Offline license certificate helpers ---

/**
 * Stores an offline checkout certificate returned by keygen-v2
 * in activation or validation responses.
 */
function cacheOfflineCertificate(certData: {
  certificate: string
  ttl: number
  expiry: string
  issued: string
}): void {
  if (!certData?.certificate) return
  store.set('offline_certificate', encrypt(JSON.stringify(certData)))
  log.info('Offline certificate cached, expires:', certData.expiry)
}

/**
 * Fetches the Keygen Ed25519 public key from keygen-v2 and caches it.
 * Only needs to succeed once — the key never changes.
 */
async function fetchAndCachePublicKey(
  licenseServer: string
): Promise<string | null> {
  const cached = store.get('keygen_public_key', null) as string | null
  if (cached) return cached

  try {
    const resp = await axios.get(licenseServer + '/public-key', {
      timeout: 5000
    })
    if (resp.data?.public_key) {
      store.set('keygen_public_key', resp.data.public_key)
      return resp.data.public_key
    }
  } catch (err: any) {
    log.warn('Failed to fetch public key:', err.message)
  }
  return null
}

/**
 * Verifies a cached machine file certificate offline using Ed25519.
 * Returns the validation status based on signature validity and expiry.
 *
 * Keygen machine file format (PEM):
 *   -----BEGIN MACHINE FILE-----
 *   <base64 payload>
 *   -----END MACHINE FILE-----
 *
 * The base64 payload decodes to JSON: { enc, sig, alg }
 * - sig is base64-encoded Ed25519 signature over enc
 * - enc is base64-encoded (optionally encrypted) JSON with license/machine data
 */
function verifyOfflineCertificate(
  publicKeyHex: string,
  fingerprint: string
): string {
  try {
    const certData = store.get('offline_certificate', null)
    if (!certData) return 'NO_CERTIFICATE'

    const parsed = JSON.parse(decrypt(certData))
    const certificate: string = parsed.certificate
    const expiry: string = parsed.expiry

    if (!certificate) return 'NO_CERTIFICATE'

    // Check expiry from the checkout response metadata first (quick check)
    if (expiry) {
      const expiryDate = new Date(expiry)
      if (expiryDate <= new Date()) {
        log.info('Offline certificate expired at', expiry)
        return 'CERTIFICATE_EXPIRED'
      }
    }

    // Strip PEM headers and decode the base64 payload
    const pemBody = certificate
      .replace(/-----BEGIN MACHINE FILE-----/, '')
      .replace(/-----END MACHINE FILE-----/, '')
      .trim()
    const decoded = JSON.parse(Buffer.from(pemBody, 'base64').toString('utf8'))

    const { enc, sig, alg } = decoded

    log.info('Certificate alg:', alg)
    log.info('Public key (first 20 chars):', publicKeyHex.substring(0, 20))
    log.info('Public key length:', publicKeyHex.length)
    log.info('enc length:', enc?.length, 'sig length:', sig?.length)

    if (!enc || !sig) return 'INVALID_CERTIFICATE'

    // Verify the Ed25519 signature
    // Detect key encoding: 64 chars = hex-encoded 32-byte key, otherwise try base64
    let pubKeyBytes: Buffer
    if (/^[0-9a-fA-F]+$/.test(publicKeyHex) && publicKeyHex.length === 64) {
      pubKeyBytes = Buffer.from(publicKeyHex, 'hex')
    } else {
      // Try base64 decoding
      pubKeyBytes = Buffer.from(publicKeyHex, 'base64')
    }

    log.info('Public key decoded bytes:', pubKeyBytes.length)

    const pubKey = crypto.createPublicKey({
      key: Buffer.concat([
        // Ed25519 DER prefix for a 32-byte public key
        Buffer.from('302a300506032b6570032100', 'hex'),
        pubKeyBytes
      ]),
      format: 'der',
      type: 'spki'
    })

    const signatureBuffer = Buffer.from(sig, 'base64')
    // Keygen signs "machine/" + enc for machine file certificates
    const dataBuffer = Buffer.from('machine/' + enc)

    const isValid = crypto.verify(null, dataBuffer, pubKey, signatureBuffer)

    if (!isValid) {
      log.warn('Offline certificate signature verification failed')
      return 'INVALID_CERTIFICATE'
    }

    log.info('Offline certificate verified successfully, valid until', expiry)
    return 'VALID'
  } catch (err: any) {
    log.error('Offline certificate verification error:', err.message)
    return 'INVALID_CERTIFICATE'
  }
}

async function hasInternetConnection(): Promise<boolean> {
  // Check cache first to avoid excessive requests
  const now = Date.now()
  if (
    internetConnectionCache &&
    now - internetConnectionCache.timestamp < CONNECTIVITY_CACHE_DURATION
  ) {
    return internetConnectionCache.status
  }

  // List of reliable connectivity check endpoints
  const connectivityEndpoints = [
    {
      url: 'https://www.cloudflare.com/cdn-cgi/trace',
      expectedStatus: 200,
      timeout: 3000
    },
    {
      url: 'https://httpbin.org/status/200',
      expectedStatus: 200,
      timeout: 3000
    },
    {
      url: 'https://detectportal.firefox.com/canonical.html',
      expectedStatus: 200,
      timeout: 3000
    },
    {
      url: 'https://clients3.google.com/generate_204',
      expectedStatus: 204,
      timeout: 3000
    }
  ]

  // Try each endpoint with a shorter timeout for faster fallback
  for (const endpoint of connectivityEndpoints) {
    try {
      const response = await axios.get(endpoint.url, {
        timeout: endpoint.timeout,
        validateStatus: (status) => status === endpoint.expectedStatus,
        headers: {
          'User-Agent': 'Mozilla/5.0 (compatible; ConnectivityCheck/1.0)',
          'Cache-Control': 'no-cache',
          Pragma: 'no-cache'
        },
        maxRedirects: 0,
        // Disable response data parsing to minimize overhead
        responseType: 'text',
        // Limit response size
        maxContentLength: 1024
      })

      // If we get here, the request was successful
      const result = true
      internetConnectionCache = { status: result, timestamp: now }
      return result
    } catch (error) {
      // Log the specific error for debugging but continue to next endpoint
      console.debug(
        `Connectivity check failed for ${endpoint.url}:`,
        error.code || error.message
      )
      continue
    }
  }

  // If all endpoints fail, perform a final DNS resolution test
  try {
    const dns = require('dns').promises
    await dns.resolve('google.com', 'A')
    const result = true
    internetConnectionCache = { status: result, timestamp: now }
    return result
  } catch (dnsError) {
    console.debug('DNS resolution test failed:', dnsError.message)
  }

  // All connectivity checks failed
  const result = false
  internetConnectionCache = { status: result, timestamp: now }
  return result
}

/*
 * IPC Communications
 * */
export default class IPCs {
  static initialize(): void {
    // Get application version
    ipcMain.on('msgGetAppVersion', (event: IpcMainEvent) => {
      event.returnValue = Constants.APP_VERSION
    })

    // Logout user
    ipcMain.on('msgLogout', (event: IpcMainEvent) => {
      console.log('Logout user')
      store.delete('authenticatedUser')
      store.delete('access_token')
      store.delete('user_profile')
      event.returnValue = 'success'
      // Restart the app to go back to login screen
      // app.relaunch()
      app.exit()
    })

    // Open url via web browser
    ipcMain.on(
      'msgOpenExternalLink',
      async (event: IpcMainEvent, url: string) => {
        await shell.openExternal(url)
      }
    )

    // check if the application has been activated
    ipcMain.on('msgGetAppStatus', (event: IpcMainEvent) => {
      event.returnValue = store.get('activated', false)
    })

    // check if an authenticated user exists
    ipcMain.on('msgGetAuthenticatedUser', (event: IpcMainEvent) => {
      event.returnValue = store.get('authenticatedUser', null)
    })

    // store authenticated user by decoding the passed token. If the token is invalid or null, then remove the authenticated user
    ipcMain.on(
      'msgStoreAuthenticatedUser',
      (event: IpcMainEvent, token: string) => {
        try {
          if (!token) {
            store.delete('authenticatedUser')
            event.returnValue = null
            return
          }
          const decodedJWT = decode(token)
          const authenticatedUser = decodedJWT.user

          store.set('authenticatedUser', authenticatedUser)
          const hashedProfile = encrypt(JSON.stringify(decodedJWT.user))
          store.set('user_profile', hashedProfile)
          store.set('access_token', token)
          event.returnValue = authenticatedUser
        } catch (e) {
          store.delete('authenticatedUser')
          store.delete('user_profile')
          store.delete('access_token')
          event.returnValue = null
        }
      }
    )

    // restart the application
    ipcMain.on(
      'msgRestartApplication',
      (event: IpcMainEvent, updateIndicator: boolean) => {
        event.returnValue = 'success'

        if (updateIndicator) {
          autoUpdater.quitAndInstall()
        } else {
          app.relaunch()
          app.exit()
        }
      }
    )

    // start downloading an available update.
    // Triggered when the user confirms the "Update Available" dialog in
    // the renderer. Once the download completes, electron-updater fires
    // `update-downloaded` and the renderer prompts the user to restart,
    // which sends `msgRestartApplication(true)` -> `quitAndInstall()`.
    ipcMain.on('msgStartDownload', (event: IpcMainEvent) => {
      event.returnValue = 'success'
      autoUpdater.downloadUpdate().catch((err) => {
        log.error('Failed to start update download:', err)
      })
    })

    // resize window
    ipcMain.on(
      'msgResizeWindow',
      (
        event: IpcMainEvent,
        width: number,
        height: number,
        fullscreen: boolean
      ) => {
        event.returnValue = 'success'
        const mainWindow = BrowserWindow.getFocusedWindow()
        if (fullscreen) {
          mainWindow.maximize()
        } else {
          mainWindow.setSize(width, height, true)
          mainWindow.center()
        }
      }
    )

    // save the base url to the store
    ipcMain.on('msgSaveBaseUrl', (event: IpcMainEvent, baseUrl: string) => {
      store.set('baseUrl', baseUrl + '/')
      event.returnValue = 'success'
    })

    // get the baseUrl
    ipcMain.on('msgGetBaseUrl', (event: IpcMainEvent) => {
      // placeholder for now
      // store.set('baseUrl', 'http://localhost:9091/')
      event.returnValue = store.get('baseUrl', null)
    })

    // get the license server url from environment
    ipcMain.on('msgGetLicenseServerUrl', (event: IpcMainEvent) => {
      const licenseServer = JSON.parse(
        decrypt(store.get('license_server', null))
      )
      event.returnValue = licenseServer
    })

    ipcMain.on('msgSetLicenseServerUrl', (event: IpcMainEvent, url: string) => {
      if (url) {
        store.set('license_server', encrypt(JSON.stringify(url)))
      }
      event.returnValue = 'success'
    })

    ipcMain.on('msgCheckLicenseValidity', async (event: IpcMainEvent) => {
      let licenseStatus = 'NO_LICENSE'

      if (store.get('license', null) === null) {
        event.returnValue = licenseStatus
        return
      }

      // Extract stored license data
      let licenseKey: string | null = null
      let fingerprint: string | null = null
      let licenseServer: string | null = null
      try {
        const license = JSON.parse(decrypt(store.get('license')))
        licenseKey =
          license?.data?.attributes?.key ||
          license?.attributes?.key ||
          license?.key
        fingerprint = await machine()
        licenseServer = JSON.parse(decrypt(store.get('license_server', null)))
      } catch {
        event.returnValue = 'INVALID'
        return
      }

      if (!licenseKey) {
        console.error('Stored license data has no key — re-activation required')
        event.returnValue = 'INVALID'
        return
      }

      // Check connectivity
      const isOnline = await hasInternetConnection()

      if (!isOnline) {
        // --- OFFLINE PATH: verify cached machine certificate ---
        log.info('Offline — attempting local certificate verification')
        const publicKey = store.get('keygen_public_key', null) as string | null
        if (!publicKey) {
          log.warn('No cached public key for offline verification')
          event.returnValue = 'NO_INTERNET'
          return
        }
        const certStatus = verifyOfflineCertificate(publicKey, fingerprint!)
        switch (certStatus) {
          case 'VALID':
            event.returnValue = 'VALID'
            return
          case 'CERTIFICATE_EXPIRED':
            event.returnValue = 'EXPIRED'
            return
          default:
            // NO_CERTIFICATE, INVALID_CERTIFICATE
            event.returnValue = 'NO_INTERNET'
            return
        }
      }

      // --- ONLINE PATH: validate against keygen-v2 ---
      // Send machine_id so the server can generate an offline certificate in the same response
      const storedMachineId = store.get('machine_id', null) as string | null
      try {
        const validationResult = await axios.post(
          licenseServer + '/validate-license',
          {
            key: licenseKey,
            fingerprint,
            machine_id: storedMachineId || undefined
          },
          {
            timeout: 10000,
            headers: {
              Authorization: 'Bearer ' + licenseKey
            }
          }
        )

        switch (validationResult.data.status) {
          case 'VALID':
            licenseStatus = 'VALID'
            break
          case 'EXPIRED':
            licenseStatus = 'EXPIRED'
            break
          case 'SUSPENDED':
            licenseStatus = 'SUSPENDED'
            break
          case 'OVERDUE':
            licenseStatus = 'OVERDUE'
            break
          default:
            licenseStatus = 'INVALID'
        }

        // Store the offline certificate if one was returned
        if (licenseStatus === 'VALID' && validationResult.data.certificate) {
          cacheOfflineCertificate(validationResult.data.certificate)
          // Also ensure we have the public key cached for future offline use
          fetchAndCachePublicKey(licenseServer!).catch(() => {})
        }
      } catch (error: any) {
        console.error('License validation failed:', error)
        if (error.response) {
          licenseStatus = 'INVALID'
        } else {
          // Server unreachable — try offline certificate as fallback
          log.info(
            'License server unreachable, attempting offline certificate fallback'
          )
          const publicKey = store.get('keygen_public_key', null) as
            | string
            | null
          log.info('Public key cached:', !!publicKey)
          const hasCert = store.get('offline_certificate', null)
          log.info('Offline certificate cached:', !!hasCert)
          if (publicKey) {
            const certStatus = verifyOfflineCertificate(publicKey, fingerprint!)
            log.info('Offline certificate verification result:', certStatus)
            if (certStatus === 'VALID') {
              log.info(
                'License server unreachable but offline certificate is valid'
              )
              event.returnValue = 'VALID'
              return
            }
          }
          licenseStatus = 'SERVER_UNREACHABLE'
        }
      }

      event.returnValue = licenseStatus
    })

    // get the user access token
    ipcMain.on('msgGetAccessToken', (event: IpcMainEvent) => {
      event.returnValue = store.get('access_token', null)
    })

    // get machine fingerprint
    ipcMain.on('msgGetMachineFingerprint', (event: IpcMainEvent) => {
      machine().then((id) => {
        event.returnValue = id
      })
    })

    // set the user license (and optional machine data for offline checkout)
    ipcMain.on(
      'msgSetUserLicense',
      (event: IpcMainEvent, license: any, machineData?: { id: string }) => {
        store.set('license', encrypt(JSON.stringify(license)))
        store.set('activated', true)
        if (machineData?.id) {
          store.set('machine_id', machineData.id)
        }
        machine().then((id) => {
          store.set('finger_print', id)
        })
        event.returnValue = 'success'
      }
    )

    // get the current enviroment

    ipcMain.on('msgGetEnvironment', (event: IpcMainEvent) => {
      event.returnValue = process.env.NODE_ENV
    })

    ipcMain.on('msgGetUserLicense', (event: IpcMainEvent) => {
      event.returnValue = JSON.parse(decrypt(store.get('license', null)))
    })

    // cache entitlements (encrypted) for offline use
    ipcMain.on(
      'msgCacheEntitlements',
      (event: IpcMainEvent, entitlements: string[]) => {
        store.set('cached_entitlements', encrypt(JSON.stringify(entitlements)))
        event.returnValue = 'success'
      }
    )

    // retrieve cached entitlements
    ipcMain.on('msgGetCachedEntitlements', (event: IpcMainEvent) => {
      try {
        const cached = store.get('cached_entitlements', null)
        if (cached) {
          event.returnValue = JSON.parse(decrypt(cached))
        } else {
          event.returnValue = null
        }
      } catch (e) {
        event.returnValue = null
      }
    })

    // logging messages
    ipcMain.on('log-message', (event, { level, message }) => {
      if (log[level]) {
        log[level](...message)
      } else {
        log.info('[Renderer]', ...message)
      }
    })

    // activate license
    ipcMain.on(
      'msgActivateLicense',
      async (event: IpcMainEvent, licenseKey: string) => {
        try {
          const fingerprint = await machine()
          const licenseServer = JSON.parse(
            decrypt(store.get('license_server', null))
          )
          const os = require('os')

          const validation = await fetch(licenseServer + '/activate-key', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              Accept: 'application/json'
            },
            body: JSON.stringify({
              key: licenseKey,
              fingerprint,
              name: os.hostname(),
              platform: `${os.type()} ${os.release()} (${os.arch()})`
            })
          })

          const rs = await validation.json()

          // Store machine ID and offline certificate from activation response
          if (rs.valid && rs.machine?.id) {
            store.set('machine_id', rs.machine.id)

            if (rs.certificate) {
              cacheOfflineCertificate(rs.certificate)
            }

            // Ensure public key is cached for future offline verification
            fetchAndCachePublicKey(licenseServer).catch(() => {})
          }

          event.returnValue = rs
        } catch (error) {
          console.error('License activation failed:', error)
          event.returnValue = null
        }
      }
    )
  }
}
