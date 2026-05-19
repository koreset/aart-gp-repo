<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center w-100">
              <span class="headline">Underwriting Cases</span>
              <v-spacer />
              <v-btn
                v-if="!quoteIdLocked"
                size="small"
                variant="text"
                prepend-icon="mdi-view-grid-outline"
                @click="goToQuoteQueue"
                >Quote queue</v-btn
              >
            </div>
          </template>
          <template #default>
            <v-row class="mb-4">
              <v-col v-if="!quoteIdLocked" cols="12" md="3">
                <v-text-field
                  v-model="quoteIdFilter"
                  label="Quote ID"
                  prepend-inner-icon="mdi-pound"
                  variant="outlined"
                  density="compact"
                  type="number"
                  clearable
                  @update:model-value="loadCases"
                />
              </v-col>
              <v-col cols="12" :md="quoteIdLocked ? 4 : 3">
                <v-select
                  v-model="statusFilter"
                  :items="statusOptions"
                  label="Status"
                  item-title="text"
                  item-value="value"
                  variant="outlined"
                  density="compact"
                  clearable
                  @update:model-value="loadCases"
                />
              </v-col>
              <v-col cols="12" :md="quoteIdLocked ? 4 : 3">
                <v-select
                  v-model="tierFilter"
                  :items="tierOptions"
                  label="Tier"
                  item-title="text"
                  item-value="value"
                  variant="outlined"
                  density="compact"
                  clearable
                  @update:model-value="loadCases"
                />
              </v-col>
              <v-col cols="12" :md="quoteIdLocked ? 4 : 3">
                <v-text-field
                  v-model="assigneeFilter"
                  label="Assignee email"
                  prepend-inner-icon="mdi-account"
                  variant="outlined"
                  density="compact"
                  clearable
                  @update:model-value="loadCases"
                />
              </v-col>
            </v-row>

            <v-row class="mb-3 align-center">
              <v-col cols="auto">
                <v-btn-toggle
                  v-model="takeoverFilter"
                  density="compact"
                  variant="outlined"
                  divided
                  mandatory
                >
                  <v-btn value="all">All cases</v-btn>
                  <v-btn value="takeover">Takeover only</v-btn>
                  <v-btn value="non_takeover">New underwriting only</v-btn>
                </v-btn-toggle>
              </v-col>
              <v-col v-if="takeoverFilter === 'takeover'" cols="auto">
                <span class="text-caption text-grey">
                  Showing cases where the engine outcome starts with
                  <code>continuation_</code>.
                </span>
              </v-col>
            </v-row>

            <v-row class="mb-4">
              <v-col cols="6" md="3">
                <stat-card
                  title="Pending evidence"
                  :value="statsByStatus.pending_evidence"
                  color="warning"
                  icon="mdi-clipboard-clock-outline"
                />
              </v-col>
              <v-col cols="6" md="3">
                <stat-card
                  title="In review"
                  :value="statsByStatus.in_review"
                  color="info"
                  icon="mdi-eye-outline"
                />
              </v-col>
              <v-col cols="6" md="3">
                <stat-card
                  title="Decided"
                  :value="statsByStatus.decided"
                  color="success"
                  icon="mdi-check-circle-outline"
                />
              </v-col>
              <v-col cols="6" md="3">
                <stat-card
                  title="Declined"
                  :value="statsByStatus.declined"
                  color="error"
                  icon="mdi-close-circle-outline"
                />
              </v-col>
              <v-col cols="6" md="3">
                <stat-card
                  title="Auto-accepted"
                  :value="statsByStatus.auto_accepted"
                  color="grey"
                  icon="mdi-shield-check-outline"
                />
              </v-col>
            </v-row>

            <div :style="{ height: gridHeight, width: '100%' }">
              <data-grid
                :column-defs="columnDefs"
                :row-data="cases"
                :loading="loading"
                style="height: 100%; width: 100%"
                @row-double-clicked="openDetail"
              />
            </div>
            <empty-state
              v-if="!loading && cases.length === 0"
              icon="mdi-clipboard-text-off-outline"
              title="No underwriting cases"
              message="Cases are created automatically when a quote contains members above the free cover limit."
            />
          </template>
        </base-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import StatCard from '@/renderer/components/StatCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useGridHeight } from '@/renderer/composables/useGridHeight'

const gridHeight = useGridHeight(420)
const router = useRouter()

interface UnderwritingCase {
  id: number
  quote_id: number
  member_name: string
  member_id_number: string
  category: string
  tier: number
  fcl_excess_ratio: number
  gla_sum_assured: number
  free_cover_limit: number
  status: string
  assigned_underwriter_email: string
  engine_outcome: string
  creation_date: string
}

// When `quoteId` is supplied (programmatic embedding), the list is
// scoped to that quote and the quote-id filter is hidden. When the
// screen is reached as a top-level route, it also honours a
// `?quote_id=` query param — that's how the quote-grouped queue drills
// into a specific quote's cases without re-opening the quote itself.
const props = defineProps<{
  quoteId?: number | string | null
}>()
const route = useRoute()
const queryQuoteId = computed(() => {
  const raw = route?.query?.quote_id
  if (Array.isArray(raw)) return raw[0] ?? null
  return typeof raw === 'string' && raw !== '' ? raw : null
})
const quoteIdLocked = computed(
  () =>
    (props.quoteId !== undefined &&
      props.quoteId !== null &&
      props.quoteId !== '') ||
    queryQuoteId.value !== null
)

const loading = ref(false)
const allCases = ref<UnderwritingCase[]>([])
const quoteIdFilter = ref<string | null>(
  props.quoteId !== undefined && props.quoteId !== null && props.quoteId !== ''
    ? String(props.quoteId)
    : queryQuoteId.value
)
watch(
  () => props.quoteId,
  (next) => {
    if (next !== undefined && next !== null && next !== '') {
      quoteIdFilter.value = String(next)
      loadCases()
    }
  }
)
watch(
  () => queryQuoteId.value,
  (next) => {
    if (
      props.quoteId === undefined ||
      props.quoteId === null ||
      props.quoteId === ''
    ) {
      quoteIdFilter.value = next
      loadCases()
    }
  }
)
const statusFilter = ref<string | null>(null)
const tierFilter = ref<string | null>(null)
const assigneeFilter = ref<string | null>(null)
const takeoverFilter = ref<'all' | 'takeover' | 'non_takeover'>('all')

const cases = computed<UnderwritingCase[]>(() => {
  if (takeoverFilter.value === 'all') return allCases.value
  if (takeoverFilter.value === 'takeover') {
    return allCases.value.filter((c) =>
      String(c.engine_outcome || '').startsWith('continuation_')
    )
  }
  return allCases.value.filter(
    (c) => !String(c.engine_outcome || '').startsWith('continuation_')
  )
})

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

const statsByStatus = computed(() => {
  const out: Record<string, number> = {
    pending_evidence: 0,
    in_review: 0,
    decided: 0,
    postponed: 0,
    declined: 0,
    auto_accepted: 0
  }
  for (const c of cases.value) {
    if (out[c.status] !== undefined) out[c.status]++
  }
  return out
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
const statusColor = (status: string) => {
  if (status === 'decided') return 'success'
  if (status === 'declined') return 'error'
  if (status === 'in_review') return 'info'
  if (status === 'postponed') return 'grey'
  if (status === 'auto_accepted') return 'grey'
  return 'warning'
}

const columnDefs = [
  {
    headerName: 'Case #',
    field: 'id',
    sortable: true,
    minWidth: 90,
    maxWidth: 110
  },
  {
    headerName: 'Quote',
    field: 'quote_id',
    sortable: true,
    minWidth: 90,
    maxWidth: 110
  },
  {
    headerName: 'Member',
    field: 'member_name',
    filter: true,
    sortable: true,
    minWidth: 180
  },
  {
    headerName: 'ID Number',
    field: 'member_id_number',
    filter: true,
    sortable: true,
    minWidth: 140
  },
  {
    headerName: 'Category',
    field: 'category',
    filter: true,
    sortable: true,
    minWidth: 130
  },
  {
    headerName: 'Tier',
    field: 'tier',
    sortable: true,
    minWidth: 130,
    cellRenderer: (params: any) =>
      `<span class="text-${tierColor(params.value)}">${tierLabel(params.value)}</span>`
  },
  {
    headerName: 'SA / FCL',
    field: 'fcl_excess_ratio',
    sortable: true,
    minWidth: 110,
    valueFormatter: (p: any) =>
      p.value ? `${(p.value as number).toFixed(2)}×` : '—'
  },
  {
    headerName: 'GLA SA',
    field: 'gla_sum_assured',
    sortable: true,
    minWidth: 140,
    valueFormatter: (p: any) =>
      p.value ? Number(p.value).toLocaleString() : '—'
  },
  {
    headerName: 'Status',
    field: 'status',
    sortable: true,
    minWidth: 140,
    cellRenderer: (params: any) =>
      `<span class="text-${statusColor(params.value)}">${params.value.replace(/_/g, ' ')}</span>`
  },
  {
    headerName: 'Assignee',
    field: 'assigned_underwriter_email',
    filter: true,
    sortable: true,
    minWidth: 200
  }
]

const loadCases = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {}
    if (quoteIdFilter.value) params.quote_id = quoteIdFilter.value
    if (statusFilter.value) params.status = statusFilter.value
    if (tierFilter.value) params.tier = tierFilter.value
    if (assigneeFilter.value) params.assignee = assigneeFilter.value
    const res = await GroupPricingService.listUnderwritingCases(params)
    allCases.value = res.data || []
  } catch (err) {
    console.error('Failed to load underwriting cases', err)
    allCases.value = []
  } finally {
    loading.value = false
  }
}

const openDetail = (event: any) => {
  const c: UnderwritingCase = event.data
  router.push({
    name: 'group-pricing-underwriting-case-detail',
    params: { caseId: String(c.id) }
  })
}

const goToQuoteQueue = () => {
  router.push({ name: 'group-pricing-underwriting-cases' })
}

onMounted(() => {
  loadCases()
})
</script>
