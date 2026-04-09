import { app, BrowserWindow, RenderProcessGoneDetails, Menu } from 'electron'
import Constants from './utils/Constants'
import IPCs from './IPCs'
const { autoUpdater } = require('electron-updater')
const log = require('electron-log')

autoUpdater.setFeedURL({
  provider: 'generic',
  url: 'https://updates.aart-enterprise.com/update/group_risk/'
})

autoUpdater.logger = log
autoUpdater.logger.transports.file.level = 'debug'
autoUpdater.forceDevUpdateConfig = true

const exitApp = (mainWindow: BrowserWindow): void => {
  if (mainWindow && !mainWindow.isDestroyed()) {
    mainWindow.hide()
  }
  mainWindow.destroy()
  autoUpdater.quitAndInstall()
  app.exit()
}

const isMac = Constants.IS_MAC
const isProduction = !Constants.IS_DEV_ENV

export const createMainWindow = async (
  mainWindow: BrowserWindow
): Promise<BrowserWindow> => {
  mainWindow = new BrowserWindow({
    title: Constants.APP_NAME,
    show: false,
    width: 1024,
    height: 600,
    useContentSize: true,
    webPreferences: Constants.DEFAULT_WEB_PREFERENCES
  })

  const template: any = [
    ...(isMac
      ? [
          {
            label: app.name,
            submenu: [
              { role: 'about' },
              { type: 'separator' },
              { role: 'services' },
              { type: 'separator' },
              { role: 'hide' },
              { role: 'hideOthers' },
              { role: 'unhide' },
              { type: 'separator' },
              {
                label: 'Log Out',
                click: () => {
                  mainWindow.webContents.send('logout')
                }
              },
              { role: 'quit' }
            ]
          }
        ]
      : []),
    {
      label: 'File',
      submenu: [isMac ? { role: 'close' } : { role: 'quit' }]
    },
    {
      label: 'Edit',
      submenu: [
        { role: 'undo' },
        { role: 'redo' },
        { type: 'separator' },
        { role: 'cut' },
        { role: 'copy' },
        { role: 'paste' },
        ...(isMac
          ? [
              { role: 'pasteAndMatchStyle' },
              { role: 'delete' },
              { role: 'selectAll' },
              { type: 'separator' },
              {
                label: 'Speech',
                submenu: [{ role: 'startSpeaking' }, { role: 'stopSpeaking' }]
              }
            ]
          : [{ role: 'delete' }, { type: 'separator' }, { role: 'selectAll' }]),
        { type: 'separator' },
        {
          label: isMac ? 'Preferences…' : 'Preferences',
          accelerator: 'CmdOrCtrl+,',
          click: () => mainWindow.webContents.send('navigate', 'app-settings')
        }
      ]
    },
    {
      label: 'Go',
      submenu: [
        {
          label: 'Dashboard',
          accelerator: 'CmdOrCtrl+Shift+D',
          click: () =>
            mainWindow.webContents.send('navigate', 'group-pricing-dashboard')
        },
        { type: 'separator' },
        {
          label: 'Quotes',
          click: () =>
            mainWindow.webContents.send('navigate', 'group-pricing-quotes')
        },
        {
          label: 'Scheme Management',
          click: () =>
            mainWindow.webContents.send('navigate', 'group-pricing-schemes')
        },
        {
          label: 'Member Management',
          click: () =>
            mainWindow.webContents.send(
              'navigate',
              'group-pricing-member-management'
            )
        },
        {
          label: 'Claims Management',
          click: () =>
            mainWindow.webContents.send(
              'navigate',
              'group-pricing-claims-management'
            )
        },
        {
          label: 'Claims Analytics',
          click: () =>
            mainWindow.webContents.send(
              'navigate',
              'group-pricing-claims-analytics'
            )
        },
        {
          label: 'Bordereaux Management',
          click: () =>
            mainWindow.webContents.send(
              'navigate',
              'group-pricing-bordereaux-management'
            )
        },
        { type: 'separator' },
        {
          label: 'PHI',
          submenu: [
            {
              label: 'Tables',
              click: () =>
                mainWindow.webContents.send(
                  'navigate',
                  'group-pricing-phi-tables'
                )
            },
            {
              label: 'Run Settings',
              click: () =>
                mainWindow.webContents.send(
                  'navigate',
                  'group-pricing-phi-run-settings'
                )
            },
            {
              label: 'Shock Settings',
              click: () =>
                mainWindow.webContents.send(
                  'navigate',
                  'group-pricing-phi-shock-settings'
                )
            },
            {
              label: 'Run Results',
              click: () =>
                mainWindow.webContents.send(
                  'navigate',
                  'group-pricing-phi-run-results'
                )
            }
          ]
        },
        {
          label: 'Premiums',
          submenu: [
            {
              label: 'Dashboard',
              click: () =>
                mainWindow.webContents.send(
                  'navigate',
                  'group-pricing-premium-dashboard'
                )
            },
            {
              label: 'Premium Schedules',
              click: () =>
                mainWindow.webContents.send(
                  'navigate',
                  'group-pricing-premium-schedules'
                )
            },
            {
              label: 'Invoices',
              click: () =>
                mainWindow.webContents.send(
                  'navigate',
                  'group-pricing-invoices'
                )
            },
            {
              label: 'Payments',
              click: () =>
                mainWindow.webContents.send(
                  'navigate',
                  'group-pricing-payments'
                )
            },
            {
              label: 'Reconciliation',
              click: () =>
                mainWindow.webContents.send(
                  'navigate',
                  'group-pricing-premium-reconciliation'
                )
            },
            {
              label: 'Arrears',
              click: () =>
                mainWindow.webContents.send('navigate', 'group-pricing-arrears')
            },
            {
              label: 'Statements',
              click: () =>
                mainWindow.webContents.send(
                  'navigate',
                  'group-pricing-statements'
                )
            }
          ]
        },
        { type: 'separator' },
        {
          label: 'Tables',
          click: () =>
            mainWindow.webContents.send('navigate', 'group-pricing-tables')
        },
        {
          label: 'Metadata',
          click: () =>
            mainWindow.webContents.send('navigate', 'group-pricing-metadata')
        },
        { type: 'separator' },
        {
          label: 'User Management',
          submenu: [
            {
              label: 'Users',
              click: () =>
                mainWindow.webContents.send('navigate', 'user-management-list')
            },
            {
              label: 'Roles',
              click: () =>
                mainWindow.webContents.send('navigate', 'user-management-roles')
            }
          ]
        }
      ]
    },
    {
      label: 'View',
      submenu: [
        { role: 'reload' },
        { role: 'forceReload' },
        ...(isProduction ? [] : [{ role: 'toggleDevTools' }]),
        { type: 'separator' },
        { role: 'resetZoom' },
        { role: 'zoomIn' },
        { role: 'zoomOut' },
        { type: 'separator' },
        { role: 'togglefullscreen' }
      ]
    },
    {
      label: 'Window',
      submenu: [
        { role: 'minimize' },
        { role: 'zoom' },
        ...(isMac
          ? [
              { type: 'separator' },
              { role: 'front' },
              { type: 'separator' },
              { role: 'window' }
            ]
          : [{ role: 'close' }])
      ]
    },
    {
      role: 'help',
      submenu: [
        {
          label: 'Documentation',
          click: () => mainWindow.webContents.send('navigate', 'documentation')
        },
        ...(!isProduction
          ? [
              { type: 'separator' },
              {
                label: 'Learn More',
                click: async () => {
                  const { shell } = require('electron')
                  await shell.openExternal('https://electronjs.org')
                }
              }
            ]
          : [])
      ]
    }
  ]

  const menu = Menu.buildFromTemplate(template)
  Menu.setApplicationMenu(menu)

  mainWindow.maximize()

  mainWindow.on('close', (event: Event): void => {
    event.preventDefault()
    exitApp(mainWindow)
  })

  mainWindow.webContents.on('did-frame-finish-load', (): void => {
    if (Constants.IS_DEV_ENV) {
      mainWindow.webContents.openDevTools()
    }
  })

  mainWindow.once('ready-to-show', (): void => {
    mainWindow.setAlwaysOnTop(true)
    mainWindow.show()
    mainWindow.focus()
    mainWindow.setAlwaysOnTop(false)
  })

  IPCs.initialize()

  if (Constants.IS_DEV_ENV) {
    await mainWindow.loadURL(Constants.APP_INDEX_URL_DEV)
  } else {
    await mainWindow.loadFile(Constants.APP_INDEX_URL_PROD)
  }

  autoUpdater.on('update-available', () => {
    mainWindow.webContents.send('update_available')
  })

  autoUpdater.on('update-not-available', () => {
    mainWindow.webContents.send('update_not_available')
  })

  autoUpdater.on('download-progress', (progressObj) => {
    mainWindow.webContents.send('download_progress', progressObj)
  })

  autoUpdater.on('update-downloaded', (info) => {
    mainWindow.webContents.send('update_downloaded', info)
  })

  autoUpdater.on('error', (error) => {
    mainWindow.webContents.send('update_error', error)
  })

  autoUpdater.checkForUpdatesAndNotify()

  return mainWindow
}

export const createErrorWindow = async (
  errorWindow: BrowserWindow,
  mainWindow: BrowserWindow,
  details?: RenderProcessGoneDetails
): Promise<BrowserWindow> => {
  if (!Constants.IS_DEV_ENV) {
    mainWindow?.hide()
  }

  errorWindow = new BrowserWindow({
    title: Constants.APP_NAME,
    show: false,
    resizable: Constants.IS_DEV_ENV,
    webPreferences: Constants.DEFAULT_WEB_PREFERENCES
  })

  errorWindow.setMenu(null)

  if (Constants.IS_DEV_ENV) {
    await errorWindow.loadURL(`${Constants.APP_INDEX_URL_DEV}#/error`)
  } else {
    await errorWindow.loadFile(Constants.APP_INDEX_URL_PROD, { hash: 'error' })
  }

  errorWindow.on('ready-to-show', (): void => {
    if (!Constants.IS_DEV_ENV && mainWindow && !mainWindow.isDestroyed()) {
      mainWindow.destroy()
    }
    errorWindow.show()
    errorWindow.focus()
  })

  errorWindow.webContents.on('did-frame-finish-load', (): void => {
    if (Constants.IS_DEV_ENV) {
      errorWindow.webContents.openDevTools()
    }
  })

  return errorWindow
}
