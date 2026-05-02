<template>
  <v-card variant="outlined">
    <v-card-text class="pa-0">
      <ag-grid-vue
        class="ag-theme-balham"
        :style="{ height: gridHeight, width: '100%' }"
        :column-defs="columnDefs"
        :row-data="rows"
        :default-col-def="defaultColDef"
        row-selection="single"
        :animate-rows="true"
      />
      <empty-state
        v-if="rows.length === 0"
        title="No in-force schemes"
        message="Once schemes go in-force they will appear here with their performance metrics."
        icon="mdi-shield-outline"
      />
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-balham.css'
import { AgGridVue } from 'ag-grid-vue3'
import EmptyState from '@/renderer/components/EmptyState.vue'
import { useGridHeight } from '@/renderer/composables/useGridHeight'

interface SchemePerformanceRow {
  scheme_id: number
  scheme_name: string
  status: string
  cover_start_date?: string
  months_in_force: number
  member_count: number
  annual_premium: number
  earned_premium: number
  itd_earned_premium: number
  commission_pct: number
  itd_claims_paid: number
  itd_claims_count: number
  avg_claim_severity: number | null
  claims_frequency: number | null
  expected_loss_ratio: number
  actual_loss_ratio: number | null
  loss_ratio_delta: number | null
  r12m_claims_paid: number
  r12m_claims_count: number
  r12m_alr: number | null
}

defineProps<{ rows: SchemePerformanceRow[] }>()
const router = useRouter()
const gridHeight = useGridHeight(420)

const fmtCurrency = (v: number | null | undefined) =>
  v == null
    ? '—'
    : new Intl.NumberFormat('en-ZA', {
        style: 'currency',
        currency: 'ZAR',
        maximumFractionDigits: 0
      }).format(v)
const fmtPct = (v: number | null | undefined) =>
  v == null ? '—' : `${v.toFixed(1)}%`
const fmtNum = (v: number | null | undefined, dp = 0) =>
  v == null ? '—' : v.toFixed(dp)

const ragColor = (alr: number | null, elr: number) => {
  if (alr == null) return undefined
  if (alr > 100) return '#ffd6d6'
  if (alr > elr) return '#fff3cd'
  return '#d4edda'
}

const defaultColDef = {
  sortable: true,
  filter: true,
  resizable: true,
  flex: 1,
  minWidth: 100
}

const columnDefs = computed(() => [
  {
    headerName: 'Scheme',
    field: 'scheme_name',
    minWidth: 200,
    pinned: 'left',
    cellStyle: { fontWeight: '600', cursor: 'pointer', color: '#1976D2' },
    onCellClicked: (p: any) =>
      router.push({
        name: 'group-pricing-schemes-detail',
        params: { id: p.data.scheme_id }
      })
  },
  { headerName: 'Status', field: 'status', maxWidth: 110 },
  {
    headerName: 'Cover Start',
    field: 'cover_start_date',
    valueFormatter: (p: any) =>
      p.value ? new Date(p.value).toISOString().slice(0, 10) : '—'
  },
  { headerName: 'Months', field: 'months_in_force', maxWidth: 100 },
  {
    headerName: 'Members',
    field: 'member_count',
    valueFormatter: (p: any) => fmtNum(p.value)
  },
  {
    headerName: 'Annual Premium',
    field: 'annual_premium',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: 'ITD Earned Prem.',
    field: 'itd_earned_premium',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: 'Commission %',
    field: 'commission_pct',
    valueFormatter: (p: any) => fmtPct(p.value)
  },
  {
    headerName: 'ITD Claims Paid',
    field: 'itd_claims_paid',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: 'Claims Count',
    field: 'itd_claims_count',
    maxWidth: 130
  },
  {
    headerName: 'Avg Severity',
    field: 'avg_claim_severity',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: 'Freq /1000',
    field: 'claims_frequency',
    valueFormatter: (p: any) => fmtNum(p.value, 2)
  },
  {
    headerName: 'ELR %',
    field: 'expected_loss_ratio',
    valueFormatter: (p: any) => fmtPct(p.value)
  },
  {
    headerName: 'ITD ALR %',
    field: 'actual_loss_ratio',
    valueFormatter: (p: any) => fmtPct(p.value),
    cellStyle: (p: any) => {
      const c = ragColor(p.value, p.data.expected_loss_ratio)
      return c ? { backgroundColor: c, fontWeight: '600' } : null
    },
    sort: 'desc'
  },
  {
    headerName: 'R12M ALR %',
    field: 'r12m_alr',
    valueFormatter: (p: any) => fmtPct(p.value),
    cellStyle: (p: any) => {
      const c = ragColor(p.value, p.data.expected_loss_ratio)
      return c ? { backgroundColor: c, fontWeight: '600' } : null
    },
    headerTooltip:
      'Rolling-12-month Actual Loss Ratio: claims paid in the last 12 months ÷ time-weighted annual premium. Compare against ITD ALR to spot recent deterioration.'
  },
  {
    headerName: 'Δ (ALR-ELR)',
    field: 'loss_ratio_delta',
    valueFormatter: (p: any) =>
      p.value == null
        ? '—'
        : `${p.value > 0 ? '+' : ''}${p.value.toFixed(1)}pp`,
    cellStyle: (p: any) => {
      if (p.value == null) return null
      if (p.value > 20) return { backgroundColor: '#ffd6d6', fontWeight: '600' }
      if (p.value > 0) return { backgroundColor: '#fff3cd' }
      return { backgroundColor: '#d4edda' }
    }
  }
])
</script>
