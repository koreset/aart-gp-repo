import { defineStore } from 'pinia'
import { ref } from 'vue'
import NotificationService from '@/renderer/api/NotificationService'

export interface AppNotification {
  id: number
  recipient_email: string
  sender_email: string
  sender_name: string
  type: string
  title: string
  body: string
  object_type: string
  object_id: number
  is_read: boolean
  read_at: string | null
  created_at: string
}

export const useNotificationStore = defineStore('notifications', () => {
  const notifications = ref<AppNotification[]>([])
  const unreadCount = ref(0)
  const loading = ref(false)
  const total = ref(0)

  async function fetchNotifications(params?: {
    page?: number
    page_size?: number
    is_read?: boolean
    type?: string
  }) {
    loading.value = true
    try {
      const res = await NotificationService.getNotifications(params || {})
      const data = res.data?.data || res.data
      notifications.value = data.notifications || []
      total.value = data.total || 0
    } catch (e) {
      console.error('Failed to fetch notifications:', e)
    } finally {
      loading.value = false
    }
  }

  async function fetchUnreadCount() {
    try {
      const res = await NotificationService.getUnreadCount()
      const data = res.data?.data || res.data
      unreadCount.value = data.unread_count || 0
    } catch (e) {
      console.error('Failed to fetch unread count:', e)
    }
  }

  function addNotification(n: AppNotification) {
    notifications.value.unshift(n)
    if (!n.is_read) {
      unreadCount.value++
    }
  }

  async function markAsRead(id: number) {
    try {
      await NotificationService.markAsRead(id)
      const n = notifications.value.find((n) => n.id === id)
      if (n && !n.is_read) {
        n.is_read = true
        n.read_at = new Date().toISOString()
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
    } catch (e) {
      console.error('Failed to mark as read:', e)
    }
  }

  async function markAllAsRead() {
    try {
      await NotificationService.markAllAsRead()
      notifications.value.forEach((n) => {
        n.is_read = true
        n.read_at = new Date().toISOString()
      })
      unreadCount.value = 0
    } catch (e) {
      console.error('Failed to mark all as read:', e)
    }
  }

  async function deleteNotification(id: number) {
    try {
      await NotificationService.deleteNotification(id)
      const idx = notifications.value.findIndex((n) => n.id === id)
      if (idx !== -1) {
        if (!notifications.value[idx].is_read) {
          unreadCount.value = Math.max(0, unreadCount.value - 1)
        }
        notifications.value.splice(idx, 1)
      }
    } catch (e) {
      console.error('Failed to delete notification:', e)
    }
  }

  return {
    notifications,
    unreadCount,
    loading,
    total,
    fetchNotifications,
    fetchUnreadCount,
    addNotification,
    markAsRead,
    markAllAsRead,
    deleteNotification
  }
})
