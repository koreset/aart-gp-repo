import { contextBridge, ipcRenderer, IpcRendererEvent } from 'electron'

// Whitelist of valid channels used for IPC communication (Send message from Renderer to Main)
const mainAvailChannels: string[] = [
  'msgGetAppVersion',
  'msgOpenExternalLink',
  'msgGetAppStatus',
  'msgGetAuthenticatedUser',
  'msgResizeWindow',
  'msgStoreAuthenticatedUser',
  'msgRestartApplication',
  'msgGetBaseUrl',
  'msgGetAccessToken',
  'msgGetMachineFingerprint',
  'msgSetUserLicense',
  'msgGetUserLicense',
  'msgSaveBaseUrl',
  'msgGetLicenseServerUrl',
  'msgSetLicenseServerUrl',
  'msgCheckLicenseValidity',
  'msgActivateLicense',
  'msgGetEnvironment',
  'msgGetEntitlements',
  'msgCacheEntitlements',
  'msgGetCachedEntitlements',
  'msgLogout',
  'log-message'
]
const rendererAvailChannels: string[] = [
  'update_available',
  'update_not_available',
  'update_downloaded',
  'download_progress',
  'update_error',
  'logout'
]

contextBridge.exposeInMainWorld('mainApi', {
  send: (channel: string, ...data: any[]): void => {
    if (mainAvailChannels.includes(channel)) {
      ipcRenderer.send.apply(null, [channel, ...data])
      if (process.env.NODE_ENV === 'development') {
        console.log({ type: 'send', channel, request: data })
      }
    } else {
      throw new Error(`Unknown ipc channel name: ${channel}`)
    }
  },
  sendSync: (channel: string, ...data: any[]): any => {
    if (mainAvailChannels.includes(channel)) {
      return ipcRenderer.sendSync.apply(null, [channel, ...data])
    }

    throw new Error(`Unknown ipc channel name: ${channel}`)
  },
  on: (
    channel: string,
    listener: (event: IpcRendererEvent, ...args: any[]) => void
  ): void => {
    if (rendererAvailChannels.includes(channel)) {
      ipcRenderer.on(channel, listener)
    } else {
      throw new Error(`Unknown ipc channel name: ${channel}`)
    }
  },
  once: (
    channel: string,
    listener: (event: IpcRendererEvent, ...args: any[]) => void
  ): void => {
    if (rendererAvailChannels.includes(channel)) {
      ipcRenderer.once(channel, listener)
    } else {
      throw new Error(`Unknown ipc channel name: ${channel}`)
    }
  },
  off: (
    channel: string,
    listener: (event: IpcRendererEvent, ...args: any[]) => void
  ): void => {
    if (rendererAvailChannels.includes(channel)) {
      ipcRenderer.off(channel, listener)
    } else {
      throw new Error(`Unknown ipc channel name: ${channel}`)
    }
  },
  invoke: async (channel: string, ...data: any[]): Promise<any> => {
    if (mainAvailChannels.includes(channel)) {
      const result = await ipcRenderer.invoke.apply(null, [channel, ...data])
      if (process.env.NODE_ENV === 'development') {
        console.log({ type: 'invoke', channel, request: data, result })
      }
      return result
    }

    throw new Error(`Unknown ipc channel name: ${channel}`)
  }
})

// Add any Node.js APIs you need to expose
contextBridge.exposeInMainWorld('nodeApi', {
  setImmediate: (callback: () => void) => setImmediate(callback),
  clearImmediate: (immediate: NodeJS.Immediate) => clearImmediate(immediate)
})

contextBridge.exposeInMainWorld('electronAPI', {
  onNavigate: (callback: (routeName: string) => void) =>
    ipcRenderer.on('navigate', (_event, routeName) => callback(routeName))
})
