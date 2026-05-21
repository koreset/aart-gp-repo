<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div
              class="d-flex justify-space-between align-center flex-wrap gap-2"
            >
              <div class="d-flex align-center">
                <v-btn
                  icon="mdi-arrow-left"
                  variant="text"
                  class="mr-2"
                  @click="goBack"
                />
                <div>
                  <span class="headline">
                    Claims on Payment Schedule
                    <span v-if="schedule">{{ schedule.schedule_number }}</span>
                  </span>
                  <div
                    v-if="schedule"
                    class="text-caption d-flex align-center gap-1 mt-1"
                  >
                    <span class="header-meta-label">Status:</span>
                    <v-chip
                      :color="statusColor(schedule.status)"
                      size="x-small"
                      label
                      variant="flat"
                    >
                      {{ statusLabel(schedule.status) }}
                    </v-chip>
                  </div>
                </div>
              </div>
              <div class="d-flex align-center gap-2 flex-wrap">
                <v-chip
                  v-if="schedule"
                  size="small"
                  variant="tonal"
                  color="primary"
                  prepend-icon="mdi-cash-multiple"
                >
                  {{ formatCurrency(schedule.total_amount) }} ·
                  {{ schedule.claims_count }} claims
                </v-chip>
                <v-btn
                  v-if="schedule"
                  rounded
                  size="small"
                  variant="outlined"
                  prepend-icon="mdi-download"
                  :disabled="!filteredItems.length"
                  @click="downloadCsv"
                >
                  Download
                </v-btn>
              </div>
            </div>
          </template>

          <template #default>
            <section class="page-section">
              <v-text-field
                v-model="search"
                variant="outlined"
                density="compact"
                prepend-inner-icon="mdi-magnify"
                placeholder="Search claim number, ID / policy number, or member name"
                clearable
                hide-details
                class="mb-3"
              />

              <v-progress-linear
                v-if="loading"
                indeterminate
                color="primary"
                class="mb-2"
              />

              <empty-state
                v-if="
                  !loading && (!schedule?.items || schedule.items.length === 0)
                "
                icon="mdi-file-document-outline"
                title="No claims on this schedule"
                message="Line items will appear here once the schedule has been generated."
              />

              <empty-state
                v-else-if="!loading && filteredItems.length === 0"
                icon="mdi-magnify-close"
                title="No matches"
                message="No claims match your search."
              />

              <v-data-table
                v-else-if="!loading"
                :headers="headers"
                :items="filteredItems"
                density="compact"
                hover
              >
                <template #[`item.claim_amount`]="{ item }">
                  {{ formatCurrency(item.claim_amount) }}
                </template>
                <template #[`item.net_payable`]="{ item }">
                  {{ formatCurrency(item.net_payable || item.claim_amount) }}
                </template>
                <template #[`item.line_status`]="{ item }">
                  <v-chip
                    v-if="item.line_status"
                    :color="lineStatusColor(item.line_status)"
                    size="x-small"
                    label
                    variant="flat"
                  >
                    {{ item.line_status }}
                  </v-chip>
                </template>
              </v-data-table>
            </section>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <v-snackbar v-model="snackbar" :color="snackbarColor" timeout="3500">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useScheduleLifecycle } from '@/renderer/composables/useScheduleLifecycle'

interface ScheduleItem {
  id: number
  claim_id: number
  claim_number: string
  member_name: string
  member_id_number: string
  benefit_name: string
  scheme_name: string
  claim_amount: number
  net_payable: number
  line_status?: string
}

interface PaymentSchedule {
  id: number
  schedule_number: string
  status: string
  total_amount: number
  claims_count: number
  created_by: string
  created_at: string
  items?: ScheduleItem[]
}

const props = defineProps<{ scheduleId: string | number }>()

const router = useRouter()
const schedule = ref<PaymentSchedule | null>(null)
const loading = ref(false)
const search = ref('')

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

const scheduleStatus = computed(() => schedule.value?.status ?? null)
const { statusLabel } = useScheduleLifecycle(scheduleStatus)

const headers = [
  { title: 'Claim #', key: 'claim_number', sortable: true },
  { title: 'Member', key: 'member_name', sortable: true },
  { title: 'ID / Policy #', key: 'member_id_number', sortable: true },
  { title: 'Benefit', key: 'benefit_name', sortable: true },
  { title: 'Scheme', key: 'scheme_name', sortable: true },
  { title: 'Status', key: 'line_status', sortable: true },
  {
    title: 'Amount',
    key: 'net_payable',
    sortable: true,
    align: 'end' as const
  }
]

const filteredItems = computed<ScheduleItem[]>(() => {
  const items = schedule.value?.items ?? []
  const q = search.value?.trim().toLowerCase()
  if (!q) return items
  return items.filter(
    (it) =>
      it.claim_number?.toLowerCase().includes(q) ||
      it.member_id_number?.toLowerCase().includes(q) ||
      it.member_name?.toLowerCase().includes(q)
  )
})

function formatCurrency(val: number) {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR'
  }).format(val ?? 0)
}

function statusColor(status: string) {
  const map: Record<string, string> = {
    draft: 'grey',
    claims_signed_off: 'info',
    finance_in_review: 'info',
    finance_first_authorised: 'warning',
    finance_second_authorised: 'warning',
    submitted_to_bank: 'warning',
    confirmed: 'success',
    archived: 'default'
  }
  return map[status] ?? 'default'
}

function lineStatusColor(status: string) {
  const map: Record<string, string> = {
    pending: 'grey',
    verified: 'success',
    queried: 'orange',
    rejected: 'error'
  }
  return map[status] ?? 'default'
}

function unwrap(res: any) {
  return res?.data?.data ?? res?.data ?? res
}

function notify(message: string, color = 'error') {
  snackbarMessage.value = message
  snackbarColor.value = color
  snackbar.value = true
}

async function load() {
  const id = Number(props.scheduleId)
  if (!id || Number.isNaN(id)) {
    notify('Invalid schedule id')
    return
  }
  loading.value = true
  try {
    const res = await GroupPricingService.getPaymentSchedule(id)
    schedule.value = unwrap(res)
  } catch {
    notify('Failed to load schedule')
  } finally {
    loading.value = false
  }
}

function goBack() {
  router.push({ name: 'group-pricing-claim-my-submissions' })
}

function csvCell(val: unknown): string {
  if (val === null || val === undefined) return ''
  const s = String(val)
  if (/[",\n\r]/.test(s)) return '"' + s.replace(/"/g, '""') + '"'
  return s
}

function downloadCsv() {
  if (!schedule.value || filteredItems.value.length === 0) return
  const header = [
    'Claim #',
    'Member',
    'ID / Policy #',
    'Benefit',
    'Scheme',
    'Status',
    'Amount'
  ]
  const rows = filteredItems.value.map((it) => [
    it.claim_number,
    it.member_name,
    it.member_id_number,
    it.benefit_name,
    it.scheme_name,
    it.line_status ?? '',
    (it.net_payable || it.claim_amount || 0).toFixed(2)
  ])
  const csv =
    '﻿' + [header, ...rows].map((r) => r.map(csvCell).join(',')).join('\r\n')
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  const today = new Date().toISOString().slice(0, 10)
  a.download = `claims-on-${schedule.value.schedule_number}-${today}.csv`
  document.body.appendChild(a)
  a.click()
  a.remove()
  URL.revokeObjectURL(url)
}

onMounted(load)
</script>

<style scoped>
.page-section {
  margin-bottom: 28px;
}
.header-meta-label {
  color: rgba(255, 255, 255, 0.75);
}
</style>
