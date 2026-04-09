<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  useNotificationStore,
  type AppNotification
} from '@/renderer/store/notifications'

const router = useRouter()
const store = useNotificationStore()

const unreadCount = computed(() => store.unreadCount)
const recentNotifications = computed(() => store.notifications.slice(0, 8))

const typeIcons: Record<string, string> = {
  quote_submitted: 'mdi-file-send',
  quote_approved: 'mdi-check-circle',
  quote_rejected: 'mdi-close-circle',
  quote_accepted: 'mdi-thumb-up',
  message_received: 'mdi-message-text',
  mention: 'mdi-at',
  scheme_status_change: 'mdi-swap-horizontal',
  schedule_reviewed: 'mdi-file-eye',
  schedule_approved: 'mdi-file-check',
  schedule_finalized: 'mdi-file-lock',
  schedule_voided: 'mdi-file-remove',
  schedule_cancelled: 'mdi-file-cancel',
  scheme_suspended: 'mdi-pause-circle',
  scheme_reinstated: 'mdi-play-circle',
  submission_reviewed: 'mdi-file-eye',
  submission_query_raised: 'mdi-help-circle',
  submission_accepted: 'mdi-file-check',
  submission_rejected: 'mdi-file-remove',
  bordereaux_reviewed: 'mdi-file-eye',
  bordereaux_approved: 'mdi-file-check',
  bordereaux_submitted: 'mdi-file-send',
  ri_bordereaux_submitted: 'mdi-file-send',
  ri_bordereaux_acknowledged: 'mdi-check-decagram',
  csm_run_reviewed: 'mdi-calculator-variant',
  csm_run_approved: 'mdi-calculator-variant',
  amendment_approved: 'mdi-check-circle',
  settlement_updated: 'mdi-bank-transfer',
  settlement_dispute_resolved: 'mdi-handshake',
  claim_payment_summary: 'mdi-cash-check'
}

const typeColors: Record<string, string> = {
  quote_submitted: 'info',
  quote_approved: 'success',
  quote_rejected: 'error',
  quote_accepted: 'success',
  message_received: 'primary',
  mention: 'warning',
  scheme_status_change: 'info',
  schedule_reviewed: 'info',
  schedule_approved: 'success',
  schedule_finalized: 'success',
  schedule_voided: 'error',
  schedule_cancelled: 'error',
  scheme_suspended: 'error',
  scheme_reinstated: 'success',
  submission_reviewed: 'info',
  submission_query_raised: 'warning',
  submission_accepted: 'success',
  submission_rejected: 'error',
  bordereaux_reviewed: 'info',
  bordereaux_approved: 'success',
  bordereaux_submitted: 'info',
  ri_bordereaux_submitted: 'info',
  ri_bordereaux_acknowledged: 'success',
  csm_run_reviewed: 'info',
  csm_run_approved: 'success',
  amendment_approved: 'success',
  settlement_updated: 'info',
  settlement_dispute_resolved: 'success',
  claim_payment_summary: 'info'
}

function getIcon(type: string) {
  return typeIcons[type] || 'mdi-bell'
}

function getColor(type: string) {
  return typeColors[type] || 'grey'
}

function timeAgo(dateStr: string): string {
  const now = new Date()
  const date = new Date(dateStr)
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000)
  if (seconds < 60) return 'just now'
  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return `${minutes}m ago`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  return `${days}d ago`
}

function onNotificationClick(n: AppNotification) {
  if (!n.is_read) {
    store.markAsRead(n.id)
  }
  // Navigate based on object type
  if (n.object_type === 'quote' && n.object_id) {
    router.push({
      name: 'group-pricing-scheme-details',
      params: { id: n.object_id }
    })
  } else if (n.object_type === 'conversation' && n.object_id) {
    router.push({
      name: 'messages-inbox',
      query: { conversation: n.object_id }
    })
  } else if (n.object_type === 'scheme' && n.object_id) {
    router.push({
      name: 'group-pricing-schemes-detail',
      params: { id: n.object_id }
    })
  } else if (n.object_type === 'premium_schedule' && n.object_id) {
    router.push({
      name: 'group-pricing-premium-schedule-detail',
      params: { id: n.object_id }
    })
  } else if (n.object_type === 'employer_submission' && n.object_id) {
    router.push({
      name: 'group-pricing-bordereaux-inbound-detail',
      params: { id: n.object_id }
    })
  } else if (n.object_type === 'bordereaux') {
    router.push({ name: 'group-pricing-bordereaux-tracking' })
  } else if (n.object_type === 'ri_bordereaux') {
    router.push({ name: 'group-pricing-ri-bordereaux' })
  } else if (n.object_type === 'technical_account') {
    router.push({ name: 'group-pricing-ri-settlement' })
  } else if (n.object_type === 'claim_payment') {
    router.push({ name: 'group-pricing-claims-management' })
  }
}

function viewAll() {
  router.push({ name: 'notification-center' })
}

onMounted(() => {
  store.fetchUnreadCount()
  store.fetchNotifications({ page: 1, page_size: 8 })
})
</script>

<template>
  <v-menu offset-y :close-on-content-click="false" max-width="380">
    <template #activator="{ props }">
      <v-btn icon variant="text" v-bind="props" size="small">
        <v-badge
          :content="unreadCount"
          :model-value="unreadCount > 0"
          color="error"
          overlap
        >
          <v-icon>mdi-bell-outline</v-icon>
        </v-badge>
      </v-btn>
    </template>

    <v-card min-width="340" max-width="380">
      <v-card-title
        class="d-flex align-center py-2 px-4"
        style="font-size: 14px"
      >
        <span>Notifications</span>
        <v-spacer />
        <v-btn
          v-if="unreadCount > 0"
          variant="text"
          size="x-small"
          color="primary"
          @click="store.markAllAsRead()"
        >
          Mark all read
        </v-btn>
      </v-card-title>
      <v-divider />

      <v-list
        v-if="recentNotifications.length > 0"
        density="compact"
        class="py-0"
      >
        <v-list-item
          v-for="n in recentNotifications"
          :key="n.id"
          :class="{ 'bg-grey-lighten-4': !n.is_read }"
          style="cursor: pointer"
          @click="onNotificationClick(n)"
        >
          <template #prepend>
            <v-icon :color="getColor(n.type)" size="20">{{
              getIcon(n.type)
            }}</v-icon>
          </template>
          <v-list-item-title style="font-size: 13px; font-weight: 500">
            {{ n.title }}
          </v-list-item-title>
          <v-list-item-subtitle style="font-size: 12px">
            {{ n.body }}
          </v-list-item-subtitle>
          <template #append>
            <span style="font-size: 11px; color: #999; white-space: nowrap">
              {{ timeAgo(n.created_at) }}
            </span>
          </template>
        </v-list-item>
      </v-list>

      <v-card-text
        v-else
        class="text-center text-grey py-6"
        style="font-size: 13px"
      >
        No notifications
      </v-card-text>

      <v-divider />
      <v-card-actions class="justify-center py-1">
        <v-btn variant="text" size="small" color="primary" @click="viewAll">
          View All
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-menu>
</template>
