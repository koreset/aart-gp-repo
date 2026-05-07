export function useErrorHandler() {
  const handleApiError = (error: any, defaultMessage = 'An error occurred') => {
    console.error('API Error:', error)

    // Api.js interceptor rejects with error.response (the response object),
    // so status lives directly on `error`, and body lives on `error.data`.
    const status = error?.status ?? error?.response?.status

    const body = error?.data ?? error?.response?.data
    const serverMessage =
      body && typeof body === 'object' ? (body.error ?? body.message) : null
    if (serverMessage) {
      return serverMessage
    }

    if (status === 409) {
      return 'This item already exists'
    }

    if (status === 404) {
      return 'Item not found'
    }

    if (status === 422) {
      return 'Invalid data provided'
    }

    if (status >= 500) {
      return 'Server error. Please try again later.'
    }

    if (error?.message) {
      return error.message
    }

    return defaultMessage
  }

  const handleValidationError = (error: any) => {
    if (error.name === 'ValidationError') {
      return error.errors?.[0] || 'Validation failed'
    }
    return handleApiError(error, 'Validation failed')
  }

  return {
    handleApiError,
    handleValidationError
  }
}
