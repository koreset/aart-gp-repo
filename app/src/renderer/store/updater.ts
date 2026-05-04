import { defineStore } from 'pinia'

// Phases of the auto-update lifecycle that the UI cares about.
//   idle        — no update activity, status bar shows nothing
//   downloading — zip is being pulled; show progress in the status bar
//   downloaded  — zip is on disk, awaiting user-confirmed restart
//   error       — install/download failed. The renderer logs to electron-log
//                 and records the message here for debugging but does NOT
//                 pop a dialog or show a status-bar indicator — interrupting
//                 the user with raw electron-updater error text was noisy
//                 and unhelpful (most "errors" are just a missing manifest
//                 on the server, which is filtered out before reaching
//                 this state anyway).
export type UpdaterPhase = 'idle' | 'downloading' | 'downloaded' | 'error'

interface UpdaterState {
  phase: UpdaterPhase
  version: string
  percent: number
  bytesPerSecond: number
  errorMessage: string
}

export const useUpdaterStore = defineStore('updater', {
  state: (): UpdaterState => ({
    phase: 'idle',
    version: '',
    percent: 0,
    bytesPerSecond: 0,
    errorMessage: ''
  }),
  actions: {
    setProgress(percent: number, bytesPerSecond: number, version?: string) {
      this.phase = 'downloading'
      // Round to integer so we don't re-render the status bar on every
      // sub-percent tick (electron-updater fires download-progress roughly
      // once a second; integer percent gives ~100 renders per download
      // instead of thousands).
      this.percent = Math.round(percent)
      this.bytesPerSecond = bytesPerSecond
      if (version) this.version = version
    },
    setDownloaded(version: string) {
      this.phase = 'downloaded'
      this.version = version
      this.percent = 100
    },
    setError(message: string) {
      this.phase = 'error'
      this.errorMessage = message
    },
    reset() {
      this.phase = 'idle'
      this.version = ''
      this.percent = 0
      this.bytesPerSecond = 0
      this.errorMessage = ''
    }
  }
})
