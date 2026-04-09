import { defineStore } from 'pinia'

export const useFlashStore = defineStore('flash', {
  state: () => ({
    message: '',
    color: 'info', // Vuetify color class (e.g., 'error', 'success')
    visible: false
  }),
  actions: {
    setMessage(message: string) {
      this.message = message
    },
    show(
      message: string,
      color: 'success' | 'error' | 'info' | 'warning' = 'info'
    ) {
      this.message = message
      this.color = color
      this.visible = true

      // Auto-hide after 4 seconds
      setTimeout(() => {
        this.visible = false
      }, 4000)
    },
    hide() {
      this.visible = false
    }
  }
})
