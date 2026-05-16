<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center">
              <v-btn
                icon="mdi-arrow-left"
                variant="text"
                class="mr-2"
                @click="goBack"
              />
              <span class="headline">Payment Exceptions</span>
              <v-spacer />
              <v-btn
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
            <v-row dense class="mb-3">
              <v-col cols="12" sm="4">
                <v-card variant="outlined" rounded="lg" class="pa-3">
                  <div class="text-caption text-medium-emphasis">Failed</div>
                  <div class="text-h5 font-weight-bold">{{
                    summary.outstanding_failed
                  }}</div>
                  <div class="text-caption text-medium-emphasis"
                    >Bank rejected, awaiting retry</div
                  >
                </v-card>
              </v-col>
              <v-col cols="12" sm="4">
                <v-card variant="outlined" rounded="lg" class="pa-3">
                  <div class="text-caption text-medium-emphasis">Unmatched</div>
                  <div class="text-h5 font-weight-bold">{{
                    summary.outstanding_unmatched
                  }}</div>
                  <div class="text-caption text-medium-emphasis"
                    >Bank ack'd but couldn't be matched</div
                  >
                </v-card>
              </v-col>
              <v-col cols="12" sm="4">
                <v-card
                  variant="outlined"
                  rounded="lg"
                  class="pa-3"
                  color="warning"
                  style="border-color: rgb(var(--v-theme-warning))"
                >
                  <div class="text-caption font-weight-medium"
                    >Outstanding value</div
                  >
                  <div class="text-h5 font-weight-bold">{{
                    formatCurrency(summary.outstanding_value)
                  }}</div>
                </v-card>
              </v-col>
            </v-row>

            <div class="d-flex align-center gap-2 mb-3 flex-wrap">
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
                style="max-width: 240px"
                hide-details
                @update:model-value="load"
              />
              <v-switch
                v-model="includeResolved"
                color="success"
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

onMounted(load)
</script>
