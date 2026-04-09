// stores/networkStatus.ts
import { defineStore } from 'pinia'
import axios from 'axios'

export const useNetworkStatusStore = defineStore('networkStatus', {
  state: () => ({
    isOnline: navigator.onLine,
    isServiceAvailable: true, // Assume true initially
    checkInterval: null as number | null
  }),
  actions: {
    updateOnlineStatus() {
      this.isOnline = navigator.onLine
      if (this.isOnline) {
        this.startCheckingService()
      } else {
        this.stopCheckingService()
      }
    },
    async checkServiceAvailability() {
      try {
        const apiServerUrl = await window.mainApi?.sendSync('msgGetBaseUrl')

        const response = await axios.get(apiServerUrl + 'health')
        this.isServiceAvailable = response.status === 200
      } catch (error) {
        this.isServiceAvailable = false
      }
    },
    startCheckingService() {
      this.checkServiceAvailability() // Initial check
      if (!this.checkInterval) {
        this.checkInterval = window.setInterval(
          this.checkServiceAvailability,
          10000
        ) // Check every 10 seconds
      }
    },
    stopCheckingService() {
      if (this.checkInterval) {
        clearInterval(this.checkInterval)
        this.checkInterval = null
      }
    }
  }
})
