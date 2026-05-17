<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center ga-3 flex-wrap">
              <v-btn
                icon="mdi-arrow-left"
                variant="text"
                class="mr-2"
                @click="goBack"
              />
              <span class="headline">Payment Exceptions</span>
              <v-spacer />
              <v-btn
                rounded
                size="small"
                variant="outlined"
                prepend-icon="mdi-download"
                :loading="exporting"
                :disabled="rows.length === 0"
                @click="exportCSV"
              >
                Export CSV
              </v-btn>
              <v-btn
                rounded
                size="small"
                variant="text"
                prepend-icon="mdi-refresh"
                :loading="loading"
                @click="load"
              >
                Refresh
              </v-btn>
            </div>
          </template>

          <template #default>
            <v-row dense class="mb-4">
              <v-col cols="12" sm="4">
                <v-card
                  variant="outlined"
                  rounded="lg"
                  class="summary-card pa-3 h-100 d-flex flex-column"
                >
                  <div class="summary-card__label">Failed</div>
                  <div class="summary-card__value">{{
                    summary.outstanding_failed
                  }}</div>
                  <div class="summary-card__hint"
                    >Bank rejected, awaiting retry</div
                  >
                </v-card>
              </v-col>
              <v-col cols="12" sm="4">
                <v-card
                  variant="outlined"
                  rounded="lg"
                  class="summary-card pa-3 h-100 d-flex flex-column"
                >
                  <div class="summary-card__label">Unmatched</div>
                  <div class="summary-card__value">{{
                    summary.outstanding_unmatched
                  }}</div>
                  <div class="summary-card__hint"
                    >Bank ack'd but couldn't be matched</div
                  >
                </v-card>
              </v-col>
              <v-col cols="12" sm="4">
                <v-card
                  variant="outlined"
                  rounded="lg"
                  class="summary-card summary-card--accent pa-3 h-100 d-flex flex-column"
                >
                  <div class="summary-card__label">Outstanding value</div>
                  <div
                    class="summary-card__value summary-card__value--warning"
                    >{{ formatCurrency(summary.outstanding_value) }}</div
                  >
                  <div class="summary-card__hint">Sum still to recover</div>
                </v-card>
              </v-col>
            </v-row>

            <div class="d-flex align-center flex-wrap ga-6 mb-4">
              <v-select
                v-model="statusFilter"
                :items="[
                  { title: 'Failed (default)', value: 'failed' },
                  { title: 'Unmatched', value: 'unmatched' },
                  { title: 'All', value: '' }
                ]"
                label="Status"
                variant="outlined"
                density="compact"
                style="max-width: 240px; min-width: 200px"
                hide-details
                @update:model-value="load"
              />
              <v-switch
                v-model="includeResolved"
                inset
                color="success"
                base-color="grey-lighten-1"
                label="Show resolved"
                density="compact"
                hide-details
                @update:model-value="load"
              />
            </div>

            <v-data-table
              :headers="headers"
              :items="rows"
              :loading="loading"
              density="compact"
              hover
            >
              <template #[`item.schedule_number`]="{ item }">
                <router-link
                  class="text-primary"
                  :to="{
                    name: 'group-pricing-claim-payment-schedule-claims',
                    params: { scheduleId: item.schedule_id }
                  }"
                >
                  {{ item.schedule_number }}
                </router-link>
              </template>
              <template #[`item.amount`]="{ item }">
                {{ formatCurrency(item.amount) }}
              </template>
              <template #[`item.status`]="{ item }">
                <v-chip
                  size="x-small"
                  variant="flat"
                  :color="
                    item.resolved
                      ? 'success'
                      : item.status === 'failed'
                        ? 'error'
                        : 'warning'
                  "
                >
                  {{ item.resolved ? 'resolved' : item.status }}
                </v-chip>
              </template>
              <template #[`item.created_at`]="{ item }">
                {{ formatDate(item.created_at) }}
              </template>
            </v-data-table>
          </template>
        </base-card>
      </v-col>
    </v-row>
    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="3500">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

const router = useRouter()

const headers = [
  { title: 'Claim #', key: 'claim_number', sortable: true },
  { title: 'Member', key: 'member_name', sortable: true },
  { title: 'Benefit', key: 'benefit_name', sortable: true },
  { title: 'Schedule', key: 'schedule_number', sortable: true },
  { title: 'Bank acct', key: 'account_number', sortable: false },
  { title: 'Amount', key: 'amount', sortable: true, align: 'end' as const },
  { title: 'Status', key: 'status', sortable: true },
  { title: 'Reason', key: 'failure_reason', sortable: false },
  { title: 'Failed at', key: 'created_at', sortable: true }
]

const rows = ref<any[]>([])
const loading = ref(false)
const exporting = ref(false)
const statusFilter = ref('failed')
const includeResolved = ref(false)
const summary = reactive({
  outstanding_failed: 0,
  outstanding_unmatched: 0,
  outstanding_value: 0
})

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

function unwrap(res: any) {
  const body = res?.data
  if (body && typeof body === 'object' && 'success' in body && 'data' in body) {
    return body.data
  }
  return body
}

function formatCurrency(val: number) {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR'
  }).format(val ?? 0)
}

function formatDate(val: string) {
  if (!val) return '—'
  return new Date(val).toLocaleString('en-ZA', {
    month: 'short',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function goBack() {
  router.push({ name: 'group-pricing-claim-payment-schedules' })
}

async function load() {
  loading.value = true
  try {
    const [listRes, summaryRes] = await Promise.all([
      GroupPricingService.listPaymentExceptions({
        status: statusFilter.value,
        includeResolved: includeResolved.value
      }),
      GroupPricingService.getPaymentExceptionsSummary()
    ])
    rows.value = unwrap(listRes) ?? []
    const s = unwrap(summaryRes) ?? {}
    summary.outstanding_failed = s.outstanding_failed ?? 0
    summary.outstanding_unmatched = s.outstanding_unmatched ?? 0
    summary.outstanding_value = s.outstanding_value ?? 0
  } catch (e: any) {
    snackbar.value = true
    snackbarColor.value = 'error'
    snackbarMessage.value =
      e?.response?.data?.message ?? 'Failed to load payment exceptions'
  } finally {
    loading.value = false
  }
}

async function exportCSV() {
  exporting.value = true
  try {
    const res = await GroupPricingService.exportPaymentExceptionsCSV({
      status: statusFilter.value,
      includeResolved: includeResolved.value,
      limit: 500
    })
    const blob = new Blob([res.data], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    const stamp = new Date().toISOString().slice(0, 19).replace(/[:T]/g, '')
    a.download = `payment_exceptions_${stamp}.csv`
    a.click()
    URL.revokeObjectURL(url)
  } catch (e: any) {
    snackbar.value = true
    snackbarColor.value = 'error'
    snackbarMessage.value =
      e?.response?.data?.message ?? 'Failed to export payment exceptions'
  } finally {
    exporting.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.summary-card {
  min-height: 110px;
  transition: border-color 0.15s ease;
}

.summary-card:hover {
  border-color: rgba(var(--v-theme-primary), 0.4);
}

.summary-card--accent {
  border-color: rgba(var(--v-theme-warning), 0.55);
}

.summary-card__label {
  font-size: 0.7rem;
  font-weight: 500;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  color: rgba(var(--v-theme-on-surface), 0.6);
  margin-bottom: 4px;
}

.summary-card__value {
  font-size: 1.5rem;
  font-weight: 700;
  line-height: 1.2;
  color: rgba(var(--v-theme-on-surface), 0.95);
}

.summary-card__value--warning {
  color: rgb(var(--v-theme-warning));
}

.summary-card__hint {
  margin-top: auto;
  padding-top: 4px;
  font-size: 0.72rem;
  color: rgba(var(--v-theme-on-surface), 0.55);
}
</style>
