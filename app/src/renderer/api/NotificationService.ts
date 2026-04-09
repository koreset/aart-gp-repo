import Api from '@/renderer/api/Api'

export default {
  getNotifications(params: {
    page?: number
    page_size?: number
    is_read?: boolean
    type?: string
  }) {
    return Api.get('/notifications', { params })
  },

  getUnreadCount() {
    return Api.get('/notifications/unread-count')
  },

  markAsRead(id: number) {
    return Api.patch(`/notifications/${id}/read`)
  },

  markAllAsRead() {
    return Api.post('/notifications/read-all')
  },

  deleteNotification(id: number) {
    return Api.delete(`/notifications/${id}`)
  }
}
