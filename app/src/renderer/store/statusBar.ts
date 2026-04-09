import { defineStore } from 'pinia'

export interface StatusBarItem {
  icon?: string
  text: string
  severity?: 'info' | 'warn' | 'error'
}

export const useStatusBarStore = defineStore('statusBar', {
  state: () => ({
    items: [] as StatusBarItem[]
  }),
  actions: {
    set(items: StatusBarItem[]) {
      this.items = items
    },
    clear() {
      this.items = []
    }
  }
})
