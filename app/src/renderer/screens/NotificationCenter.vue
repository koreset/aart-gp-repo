<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  useNotificationStore,
  type AppNotification
} from '@/renderer/store/notifications'
import BaseCard from '@/renderer/components/BaseCard.vue'

const router = useRouter()
const store = useNotificationStore()

const filterType = ref('')
const filterRead = ref<string>('all')
const page = ref(1)
const pageSize = 20

const typeOptions = [
  { title: 'All Types', value: '' },
  { title: 'Quote Submitted', value: 'quote_submitted' },
  { title: 'Quote Approved', value: 'quote_approved' },
  { title: 'Quote Rejected', value: 'quote_rejected' },
  { title: 'Quote Accepted', value: 'quote_accepted' },
  { title: 'Message Received', value: 'message_received' },
  { title: 'Mention', value: 'mention' },
  { title: 'Scheme Status', value: 'scheme_status_change' },
  { title: 'Schedule Reviewed', value: 'schedule_reviewed' },
  { title: 'Schedule Approved', value: 'schedule_approved' },
  { title: 'Schedule Finalized', value: 'schedule_finalized' },
  { title: 'Schedule Voided', value: 'schedule_voided' },
  { title: 'Schedule Cancelled', value: 'schedule_cancelled' },
  { title: 'Scheme Suspended', value: 'scheme_suspended' },
  { title: 'Scheme Reinstated', value: 'scheme_reinstated' },
  { title: 'Submission Reviewed', value: 'submission_reviewed' },
  { title: 'Submission Query', value: 'submission_query_raised' },
  { title: 'Submission Accepted', value: 'submission_accepted' },
  { title: 'Submission Rejected', value: 'submission_rejected' },
  { title: 'Bordereaux Reviewed', value: 'bordereaux_reviewed' },
  { title: 'Bordereaux Approved', value: 'bordereaux_approved' },
  { title: 'Bordereaux Submitted', value: 'bordereaux_submitted' },
  { title: 'RI Bordereaux Submitted', value: 'ri_bordereaux_submitted' },
  { title: 'RI Bordereaux Acknowledged', value: 'ri_bordereaux_acknowledged' },
  { title: 'CSM Run Reviewed', value: 'csm_run_reviewed' },
  { title: 'CSM Run Approved', value: 'csm_run_approved' },
  { title: 'Amendment Approved', value: 'amendment_approved' },
  { title: 'Settlement Updated', value: 'settlement_updated' },
  { title: 'Dispute Resolved', value: 'settlement_dispute_resolved' },
  { title: 'Claim Payment Summary', value: 'claim_payment_summary' }
]

const readOptions = [
  { title: 'All', value: 'all' },
  { title: 'Unread', value: 'unread' },
  { title: 'Read', value: 'read' }
]

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

function fetchData() {
  const params: any = { page: page.value, page_size: pageSize }
  if (filterType.value) params.type = filterType.value
  if (filterRead.value === 'unread') params.is_read = false
  else if (filterRead.value === 'read') params.is_read = true
  store.fetchNotifications(params)
}

function onNotificationClick(n: AppNotification) {
  if (!n.is_read) store.markAsRead(n.id)
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

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleString()
}

const totalPages = ref(0)
watch(
  () => store.total,
  (t) => {
    totalPages.value = Math.ceil(t / pageSize)
  }
)

watch([filterType, filterRead], () => {
  page.value = 1
  fetchData()
})

watch(page, () => fetchData())

onMounted(() => {
  fetchData()
  store.fetchUnreadCount()
})
</script>

<template>
  <v-container fluid>
    <base-card :show-actions="false">
      <template #header>
        <div class="d-flex align-center" style="width: 100%">
          <span class="headline">Notifications</span>
          <v-spacer />
          <v-btn
            variant="outlined"
            size="small"
            color="white"
            :disabled="store.unreadCount === 0"
            @click="store.markAllAsRead()"
          >
            Mark All Read
          </v-btn>
        </div>
      </template>
      <template #default>
        <v-row class="mb-4">
          <v-col cols="3">
            <v-select
              v-model="filterRead"
              :items="readOptions"
              label="Status"
              density="compact"
              variant="outlined"
              hide-details
            />
          </v-col>
          <v-col cols="3">
            <v-select
              v-model="filterType"
              :items="typeOptions"
              label="Type"
              density="compact"
              variant="outlined"
              hide-details
            />
          </v-col>
        </v-row>

        <v-list v-if="store.notifications.length > 0" density="compact">
          <v-list-item
            v-for="n in store.notifications"
            :key="n.id"
            :class="{ 'bg-grey-lighten-4': !n.is_read }"
            style="cursor: pointer"
            @click="onNotificationClick(n)"
          >
            <template #prepend>
              <v-icon
                :color="typeColors[n.type] || 'grey'"
                size="22"
                class="mr-3"
              >
                {{ typeIcons[n.type] || 'mdi-bell' }}
              </v-icon>
            </template>
            <v-list-item-title style="font-size: 13px; font-weight: 500">
              {{ n.title }}
            </v-list-item-title>
            <v-list-item-subtitle style="font-size: 12px">
              {{ n.body }}
            </v-list-item-subtitle>
            <template #append>
              <div class="d-flex flex-column align-end">
                <span style="font-size: 11px; color: #999">{{
                  formatDate(n.created_at)
                }}</span>
                <span v-if="n.sender_name" style="font-size: 11px; color: #666">
                  from {{ n.sender_name }}
                </span>
              </div>
            </template>
          </v-list-item>
        </v-list>

        <div v-else class="text-center text-grey py-8">
          No notifications found
        </div>

        <v-row v-if="totalPages > 1" class="mt-4 justify-center">
          <v-pagination v-model="page" :length="totalPages" size="small" />
        </v-row>
      </template>
    </base-card>
  </v-container>
</template>
