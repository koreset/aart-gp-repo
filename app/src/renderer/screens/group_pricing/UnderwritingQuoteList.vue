<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center w-100">
              <span class="headline">Underwriting — Quote Queue</span>
              <v-spacer />
              <v-btn
                size="small"
                variant="text"
                prepend-icon="mdi-format-list-bulleted"
                @click="openFlatList"
                >Flat case list</v-btn
              >
            </div>
          </template>

          <template #default>
            <v-row class="mb-4">
              <v-col cols="12" md="4">
                <v-select
                  v-model="statusFilter"
                  :items="statusOptions"
                  label="Case status"
                  item-title="text"
                  item-value="value"
                  variant="outlined"
                  density="compact"
                  clearable
                  @update:model-value="loadSummaries"
                />
              </v-col>
              <v-col cols="12" md="4">
                <v-select
                  v-model="tierFilter"
                  :items="tierOptions"
                  label="Tier"
                  item-title="text"
                  item-value="value"
                  variant="outlined"
                  density="compact"
                  clearable
                  @update:model-value="loadSummaries"
                />
              </v-col>
              <v-col cols="12" md="4">
                <v-text-field
                  v-model="assigneeFilter"
                  label="Assignee email"
                  prepend-inner-icon="mdi-account"
                  variant="outlined"
                  density="compact"
                  clearable
                  @update:model-value="loadSummaries"
                />
              </v-col>
            </v-row>

            <v-row class="mb-4">
              <v-col cols="6" md="3">
                <stat-card
                  title="Quotes with cases"
                  :value="summaries.length"
                  color="primary"
                  icon="mdi-file-document-multiple-outline"
                />
              </v-col>
              <v-col cols="6" md="3">
                <stat-card
                  title="Pending evidence"
                  :value="totals.pendingQuotes"
                  color="warning"
                  icon="mdi-clipboard-clock-outline"
                />
              </v-col>
              <v-col cols="6" md="3">
                <stat-card
                  title="In review"
                  :value="totals.inReviewQuotes"
                  color="info"
                  icon="mdi-eye-outline"
                />
              </v-col>
              <v-col cols="6" md="3">
                <stat-card
                  title="Fully decided"
                  :value="totals.decidedQuotes"
                  color="success"
                  icon="mdi-check-circle-outline"
                />
              </v-col>
            </v-row>

            <div :style="{ height: gridHeight, width: '100%' }">
              <data-grid
                :column-defs="columnDefs"
                :row-data="summaries"
                :loading="loading"
                style="height: 100%; width: 100%"
                @row-double-clicked="openQuote"
              />
            </div>

            <empty-state
              v-if="!loading && summaries.length === 0"
              icon="mdi-clipboard-text-off-outline"
              title="No underwriting workload"
              message="Quotes with members above the free cover limit will appear here once they're calculated."
            />
          </template>
        </base-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import StatCard from '@/renderer/components/StatCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useGridHeight } from '@/renderer/composables/useGridHeight'

const gridHeight = useGridHeight(420)
const router = useRouter()

interface QuoteSummary {
  quote_id: number
  quote_name: string
  scheme_name: string
  broker_name: string
  quote_status: string
  total_cases: number
  pending_evidence_count: number
  in_review_count: number
  decided_count: number
  postponed_count: number
  declined_count: number
  auto_accepted_count: number
  top_tier: number
  latest_activity_at: string | null
}

const loading = ref(false)
const summaries = ref<QuoteSummary[]>([])
const statusFilter = ref<string | null>(null)
const tierFilter = ref<string | null>(null)
const assigneeFilter = ref<string | null>(null)

const statusOptions = [
  { text: 'Pending evidence', value: 'pending_evidence' },
  { text: 'In review', value: 'in_review' },
  { text: 'Decided', value: 'decided' },
  { text: 'Postponed', value: 'postponed' },
  { text: 'Declined', value: 'declined' },
  { text: 'Auto-accepted', value: 'auto_accepted' }
]
const tierOptions = [
  { text: 'Tier 1 — short-form', value: '1' },
  { text: 'Tier 2 — full UW', value: '2' }
]

const totals = computed(() => {
  let pendingQuotes = 0
  let inReviewQuotes = 0
  let decidedQuotes = 0
  for (const q of summaries.value) {
    if (q.pending_evidence_count > 0) pendingQuotes++
    if (q.in_review_count > 0) inReviewQuotes++
    if (
      q.total_cases > 0 &&
      q.pending_evidence_count === 0 &&
      q.in_review_count === 0
    )
      decidedQuotes++
  }
  return { pendingQuotes, inReviewQuotes, decidedQuotes }
})

const tierLabel = (tier: number) => {
  if (tier === 2) return 'Full UW'
  if (tier === 1) return 'Short-form'
  return 'Within FCL'
}
const tierColor = (tier: number) => {
  if (tier === 2) return 'error'
  if (tier === 1) return 'warning'
  return 'success'
}

const formatDate = (s: string | null) =>
  s ? new Date(s).toLocaleString() : '—'

const columnDefs = [
  {
    headerName: 'Quote',
    field: 'quote_name',
    sortable: true,
    filter: true,
    minWidth: 180,
    cellRenderer: (p: any) =>
      `<span class="font-weight-medium text-primary">${
        p.value || `#${p.data?.quote_id}`
      }</span>`
  },
  {
    headerName: 'Scheme',
    field: 'scheme_name',
    sortable: true,
    filter: true,
    minWidth: 180
  },
  {
    headerName: 'Broker',
    field: 'broker_name',
    sortable: true,
    filter: true,
    minWidth: 160
  },
  {
    headerName: 'Top tier',
    field: 'top_tier',
    sortable: true,
    minWidth: 130,
    cellRenderer: (p: any) =>
      `<span class="text-${tierColor(p.value)}">${tierLabel(p.value)}</span>`
  },
  {
    headerName: 'Pending ev.',
    field: 'pending_evidence_count',
    sortable: true,
    minWidth: 120,
    maxWidth: 140,
    cellStyle: (p: any) =>
      p.value > 0 ? { color: '#b45309', fontWeight: '600' } : null
  },
  {
    headerName: 'In review',
    field: 'in_review_count',
    sortable: true,
    minWidth: 110,
    maxWidth: 130,
    cellStyle: (p: any) =>
      p.value > 0 ? { color: '#1d4ed8', fontWeight: '600' } : null
  },
  {
    headerName: 'Decided',
    field: 'decided_count',
    sortable: true,
    minWidth: 110,
    maxWidth: 130,
    cellStyle: (p: any) =>
      p.value > 0 ? { color: '#15803d', fontWeight: '600' } : null
  },
  {
    headerName: 'Postponed',
    field: 'postponed_count',
    sortable: true,
    minWidth: 110,
    maxWidth: 130
  },
  {
    headerName: 'Declined',
    field: 'declined_count',
    sortable: true,
    minWidth: 110,
    maxWidth: 130,
    cellStyle: (p: any) =>
      p.value > 0 ? { color: '#b91c1c', fontWeight: '600' } : null
  },
  {
    headerName: 'Auto-accepted',
    field: 'auto_accepted_count',
    sortable: true,
    minWidth: 130,
    maxWidth: 150
  },
  {
    headerName: 'Total',
    field: 'total_cases',
    sortable: true,
    minWidth: 90,
    maxWidth: 110
  },
  {
    headerName: 'Last activity',
    field: 'latest_activity_at',
    sortable: true,
    minWidth: 180,
    valueFormatter: (p: any) => formatDate(p.value)
  }
]

const loadSummaries = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {}
    if (statusFilter.value) params.status = statusFilter.value
    if (tierFilter.value) params.tier = tierFilter.value
    if (assigneeFilter.value) params.assignee = assigneeFilter.value
    const res =
      await GroupPricingService.listUnderwritingCaseQuoteSummaries(params)
    summaries.value = res.data || []
  } catch (err) {
    console.error('Failed to load underwriting quote summaries', err)
    summaries.value = []
  } finally {
    loading.value = false
  }
}

const openQuote = (event: any) => {
  // Drill into the flat case list scoped to this quote. The underwriting
  // workflow is intentionally independent of the quote pages so different
  // teams can manage access without touching quote permissions.
  const row: QuoteSummary = event.data
  router.push({
    name: 'group-pricing-underwriting-cases-flat',
    query: { quote_id: String(row.quote_id) }
  })
}

const openFlatList = () => {
  router.push({ name: 'group-pricing-underwriting-cases-flat' })
}

onMounted(loadSummaries)
</script>
