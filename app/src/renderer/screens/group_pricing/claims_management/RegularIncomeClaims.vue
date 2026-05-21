<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex justify-space-between align-center flex-wrap">
              <span class="headline">Regular Income Claims</span>
              <div class="d-flex align-center gap-2">
                <span class="text-caption metrics-timestamp mr-2">
                  Metrics refreshed: {{ lastUpdatedLabel }}
                </span>
                <v-btn
                  size="small"
                  variant="outlined"
                  rounded
                  prepend-icon="mdi-refresh"
                  :loading="loading || metricsLoading"
                  @click="refreshAll"
                >
                  Refresh
                </v-btn>
                <v-btn
                  size="small"
                  variant="outlined"
                  rounded
                  prepend-icon="mdi-arrow-left"
                  @click="
                    router.push({ name: 'group-pricing-claims-management' })
                  "
                >
                  Back to Claims
                </v-btn>
              </div>
            </div>
          </template>
          <template #default>
            <v-row class="mb-4">
              <v-col cols="12" sm="6" md="3">
                <stat-card
                  title="Total Claims"
                  :value="String(metrics.total ?? 0)"
                  icon="mdi-clipboard-list-outline"
                  color="primary"
                  :loading="metricsLoading"
                  :subtitle="`${benefitTypeFilter || 'PHI + TTD'}`"
                />
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <stat-card
                  title="By Display Status"
                  :value="topStatus.value"
                  icon="mdi-progress-clock"
                  color="info"
                  :loading="metricsLoading"
                  :subtitle="topStatus.label"
                />
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <stat-card
                  title="Top Region"
                  :value="topRegion.value"
                  icon="mdi-map-marker-outline"
                  color="accent"
                  :loading="metricsLoading"
                  :subtitle="topRegion.label"
                />
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <stat-card
                  title="Top Occupation Class"
                  :value="topOccupation.value"
                  icon="mdi-briefcase-outline"
                  color="warning"
                  :loading="metricsLoading"
                  :subtitle="topOccupation.label"
                />
              </v-col>
            </v-row>

            <v-row class="mb-4">
              <v-col cols="12" md="3">
                <v-text-field
                  v-model="searchQuery"
                  label="Search"
                  prepend-inner-icon="mdi-magnify"
                  variant="outlined"
                  density="compact"
                  clearable
                />
              </v-col>
              <v-col cols="12" md="3">
                <v-select
                  v-model="selectedScheme"
                  :items="schemes"
                  label="Scheme"
                  item-title="name"
                  item-value="id"
                  variant="outlined"
                  density="compact"
                  clearable
                  @update:model-value="loadAll"
                />
              </v-col>
              <v-col cols="12" md="2">
                <v-select
                  v-model="selectedDisplayStatus"
                  :items="displayStatusOptions"
                  label="Status"
                  item-title="text"
                  item-value="value"
                  variant="outlined"
                  density="compact"
                  clearable
                  @update:model-value="loadAll"
                />
              </v-col>
              <v-col cols="12" md="2">
                <v-select
                  v-model="selectedRegion"
                  :items="regionOptions"
                  label="Region"
                  variant="outlined"
                  density="compact"
                  clearable
                  @update:model-value="loadAll"
                />
              </v-col>
              <v-col cols="12" md="2">
                <v-select
                  v-model="selectedOccupation"
                  :items="occupationOptions"
                  label="Occupation Class"
                  variant="outlined"
                  density="compact"
                  clearable
                  @update:model-value="loadAll"
                />
              </v-col>
            </v-row>

            <div :style="{ height: gridHeight, width: '100%' }">
              <data-grid
                :column-defs="columnDefs"
                :row-data="filteredClaims"
                :loading="loading"
                style="height: 100%; width: 100%"
                @row-double-clicked="viewClaimDetails"
              />
            </div>
            <empty-state
              v-if="!loading && filteredClaims.length === 0"
              icon="mdi-clipboard-text-off-outline"
              title="No regular income claims"
              message="No PHI or TTD claims match the current filters."
            />
          </template>
        </base-card>
      </v-col>
    </v-row>

    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="3000">
      {{ snackbarMessage }}
      <template #actions>
        <v-btn color="white" variant="text" @click="snackbar = false"
          >Close</v-btn
        >
      </template>
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { DateTime } from 'luxon'
import BaseCard from '@/renderer/components/BaseCard.vue'
import StatCard from '@/renderer/components/StatCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { statusCellRenderer } from '@/renderer/utils/statusCellRenderer'
import { currencyFormatter, dateFormatter } from '@/renderer/utils/formatters'
import { useGridHeight } from '@/renderer/composables/useGridHeight'

interface RegularIncomeClaim {
  id: number
  claim_number: string
  member_name: string
  member_id_number: string
  gender: string
  date_of_birth: string
  date_of_event: string
  date_notified: string
  termination_date: string
  deferred_period: number
  duration_in_force_months: number
  display_status: string
  persisted_status: string
  escalation_option: string
  normal_retirement_age: number
  tiered_income_replacement: string
  benefit_escalation_month: string
  benefit_type: string
  region: string
  occupational_class: string
  scheme_id: number
  scheme_name: string
  claim_amount: number
}

interface MetricsResponse {
  total: number
  by_region: Record<string, number>
  by_occupation_class: Record<string, number>
  by_display_status: Record<string, number>
  generated_at: string
}

interface Scheme {
  id: number
  name: string
}

const router = useRouter()
const gridHeight = useGridHeight(420)

const loading = ref(false)
const metricsLoading = ref(false)
const claims = ref<RegularIncomeClaim[]>([])
const schemes = ref<Scheme[]>([])
const metrics = ref<MetricsResponse>({
  total: 0,
  by_region: {},
  by_occupation_class: {},
  by_display_status: {},
  generated_at: ''
})

const searchQuery = ref('')
const selectedScheme = ref<number | null>(null)
const selectedRegion = ref<string | null>(null)
const selectedOccupation = ref<string | null>(null)
const selectedDisplayStatus = ref<string | null>(null)

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

const displayStatusOptions = [
  { text: 'Notified', value: 'notified' },
  { text: 'Pending', value: 'pending' },
  { text: 'In Payment', value: 'in_payment' },
  { text: 'Ombud Claim', value: 'ombud_claim' }
]

const benefitTypeFilter = computed(() => '')

const columnDefs = [
  {
    headerName: 'Claim Number',
    field: 'claim_number',
    filter: true,
    sortable: true,
    minWidth: 150,
    pinned: 'left' as const,
    cellRenderer: (params: any) =>
      `<span class="font-weight-medium text-primary">${params.value ?? ''}</span>`
  },
  { headerName: 'Member', field: 'member_name', filter: true, sortable: true, minWidth: 170 },
  { headerName: 'Gender', field: 'gender', sortable: true, minWidth: 90 },
  {
    headerName: 'Date of Birth',
    field: 'date_of_birth',
    sortable: true,
    minWidth: 120,
    valueFormatter: dateFormatter
  },
  {
    headerName: 'Date of Event',
    field: 'date_of_event',
    sortable: true,
    minWidth: 120,
    valueFormatter: dateFormatter
  },
  {
    headerName: 'Date Notified',
    field: 'date_notified',
    sortable: true,
    minWidth: 120,
    valueFormatter: dateFormatter
  },
  {
    headerName: 'Termination Date',
    field: 'termination_date',
    sortable: true,
    minWidth: 130,
    valueFormatter: dateFormatter
  },
  {
    headerName: 'Deferred Period (mo)',
    field: 'deferred_period',
    sortable: true,
    minWidth: 130
  },
  {
    headerName: 'Duration In Force (mo)',
    field: 'duration_in_force_months',
    sortable: true,
    minWidth: 150
  },
  {
    headerName: 'Status',
    field: 'display_status',
    sortable: true,
    minWidth: 140,
    cellRenderer: (params: any) => statusCellRenderer(params.value)
  },
  {
    headerName: 'Escalation Option',
    field: 'escalation_option',
    sortable: true,
    minWidth: 160
  },
  {
    headerName: 'Normal Retirement Age',
    field: 'normal_retirement_age',
    sortable: true,
    minWidth: 140
  },
  {
    headerName: 'Tiered Income Replacement',
    field: 'tiered_income_replacement',
    sortable: true,
    minWidth: 180
  },
  {
    headerName: 'Benefit Escalation Month',
    field: 'benefit_escalation_month',
    sortable: true,
    minWidth: 170
  },
  { headerName: 'Benefit', field: 'benefit_type', sortable: true, minWidth: 100 },
  { headerName: 'Region', field: 'region', sortable: true, minWidth: 130 },
  {
    headerName: 'Occupation Class',
    field: 'occupational_class',
    sortable: true,
    minWidth: 140
  },
  { headerName: 'Scheme', field: 'scheme_name', sortable: true, minWidth: 160 },
  {
    headerName: 'Claim Amount',
    field: 'claim_amount',
    sortable: true,
    minWidth: 130,
    valueFormatter: currencyFormatter
  }
]

const regionOptions = computed(() =>
  Array.from(
    new Set(claims.value.map((c) => c.region).filter((r) => r))
  ).sort()
)

const occupationOptions = computed(() =>
  Array.from(
    new Set(
      claims.value.map((c) => c.occupational_class).filter((o) => o)
    )
  ).sort()
)

const filteredClaims = computed(() => {
  if (!searchQuery.value) return claims.value
  const q = searchQuery.value.toLowerCase()
  return claims.value.filter(
    (c) =>
      (c.claim_number ?? '').toLowerCase().includes(q) ||
      (c.member_name ?? '').toLowerCase().includes(q) ||
      (c.member_id_number ?? '').includes(q) ||
      (c.scheme_name ?? '').toLowerCase().includes(q)
  )
})

const lastUpdatedLabel = computed(() => {
  if (!metrics.value.generated_at) return '—'
  const dt = DateTime.fromISO(metrics.value.generated_at)
  if (!dt.isValid) return metrics.value.generated_at
  return dt.toLocal().toFormat('dd MMM yyyy HH:mm')
})

function topEntry(map: Record<string, number>) {
  let bestKey = '—'
  let bestValue = 0
  let total = 0
  for (const [k, v] of Object.entries(map ?? {})) {
    total += v
    if (v > bestValue) {
      bestKey = k
      bestValue = v
    }
  }
  return { key: bestKey, value: bestValue, total }
}

const topStatus = computed(() => {
  const t = topEntry(metrics.value.by_display_status)
  return {
    value: t.key === '—' ? '—' : labelFor(t.key),
    label: `${t.value} of ${t.total}`
  }
})

const topRegion = computed(() => {
  const t = topEntry(metrics.value.by_region)
  return { value: t.key, label: `${t.value} of ${t.total}` }
})

const topOccupation = computed(() => {
  const t = topEntry(metrics.value.by_occupation_class)
  return { value: t.key, label: `${t.value} of ${t.total}` }
})

function labelFor(key: string): string {
  const map: Record<string, string> = {
    notified: 'Notified',
    pending: 'Pending',
    in_payment: 'In Payment',
    ombud_claim: 'Ombud Claim'
  }
  return map[key] ?? key
}

function currentFilters() {
  return {
    scheme_id: selectedScheme.value ?? undefined,
    region: selectedRegion.value ?? undefined,
    occupational_class: selectedOccupation.value ?? undefined,
    display_status: selectedDisplayStatus.value ?? undefined
  }
}

async function loadClaims() {
  loading.value = true
  try {
    const res = await GroupPricingService.getRegularIncomeClaims(
      currentFilters()
    )
    claims.value = res.data || []
  } catch (e) {
    console.error('Error loading regular income claims:', e)
    showSnackbar('Error loading regular income claims', 'error')
    claims.value = []
  } finally {
    loading.value = false
  }
}

async function loadMetrics() {
  metricsLoading.value = true
  try {
    const res = await GroupPricingService.getRegularIncomeClaimsMetrics(
      currentFilters()
    )
    metrics.value = res.data || {
      total: 0,
      by_region: {},
      by_occupation_class: {},
      by_display_status: {},
      generated_at: new Date().toISOString()
    }
  } catch (e) {
    console.error('Error loading metrics:', e)
    showSnackbar('Error loading metrics', 'error')
  } finally {
    metricsLoading.value = false
  }
}

async function loadSchemes() {
  try {
    const res = await GroupPricingService.getSchemesWithCoverageHistory()
    schemes.value = res.data || []
  } catch (e) {
    console.error('Error loading schemes:', e)
    schemes.value = []
  }
}

async function loadAll() {
  await Promise.all([loadClaims(), loadMetrics()])
}

async function refreshAll() {
  await loadAll()
  showSnackbar('Refreshed', 'success')
}

function viewClaimDetails(row: any) {
  const claim = row.data
  if (!claim?.id) return
  router.push({ name: 'group-pricing-claim-details', params: { id: claim.id } })
}

function showSnackbar(message: string, color: string = 'success') {
  snackbarMessage.value = message
  snackbarColor.value = color
  snackbar.value = true
}

onMounted(async () => {
  await loadSchemes()
  await loadAll()
})
</script>

<style scoped>
.cursor-pointer {
  cursor: pointer;
}

.metrics-timestamp {
  color: rgba(255, 255, 255, 0.85);
}
</style>
