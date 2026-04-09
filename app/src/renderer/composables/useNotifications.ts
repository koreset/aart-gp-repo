import { ref } from 'vue'

export function useNotifications() {
  const snackbar = ref(false)
  const snackbarMessage = ref('')
  const snackbarColor = ref('success')

  const showSuccess = (message: string) => {
    snackbarMessage.value = message
    snackbarColor.value = 'success'
    snackbar.value = true
  }

  const showError = (message: string) => {
    snackbarMessage.value = message
    snackbarColor.value = 'error'
    snackbar.value = true
  }

  const showWarning = (message: string) => {
    snackbarMessage.value = message
    snackbarColor.value = 'warning'
    snackbar.value = true
  }

  const hideNotification = () => {
    snackbar.value = false
  }

  return {
    snackbar,
    snackbarMessage,
    snackbarColor,
    showSuccess,
    showError,
    showWarning,
    hideNotification
  }
}
