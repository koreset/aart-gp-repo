import { app, WebContents, RenderProcessGoneDetails } from 'electron'
import Constants from './utils/Constants'
import { createErrorWindow, createMainWindow } from './MainRunner'

let mainWindow
let errorWindow

process.on('unhandledRejection', (reason, promise) => {
  console.error('Unhandled Promise Rejection at:', promise, 'reason:', reason)
})

process.on('uncaughtException', (error) => {
  console.error('Uncaught Exception:', error)
})

const gotTheLock = app.requestSingleInstanceLock()

if (!gotTheLock) {
  app.quit()
} else {
  app.on('second-instance', () => {
    if (mainWindow) {
      if (mainWindow.isMinimized()) mainWindow.restore()
      mainWindow.focus()
    }
  })

  app.on('ready', async () => {
    if (Constants.IS_DEV_ENV) {
      import('./index.dev')
    }

    mainWindow = await createMainWindow(mainWindow)
  })

  app.on('activate', async () => {
    if (!mainWindow) {
      mainWindow = await createMainWindow(mainWindow)
    }
  })

  app.on('window-all-closed', () => {
    mainWindow = null
    errorWindow = null

    if (!Constants.IS_MAC) {
      app.quit()
    }
  })

  app.on(
    'render-process-gone',
    (
      event: Event,
      webContents: WebContents,
      details: RenderProcessGoneDetails
    ) => {
      errorWindow = createErrorWindow(errorWindow, mainWindow, details)
    }
  )
}
