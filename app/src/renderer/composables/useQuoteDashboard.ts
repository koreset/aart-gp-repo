import { ref } from 'vue'

import QuoteDashboardService, {
  type DashboardFilter,
  type OpenUserFlagBody,
  type QuoteExtractFilter,
  type QuoteSlaTarget,
  type QuoteUserFlag,
  type UserFlagsFilter
} from '@/renderer/api/QuoteDashboardService'
import { useFilterPersistence } from '@/renderer/composables/useFilterPersistence'

export interface QuotePerformanceKpis {
  user_name: string
  total_quotes: number
  draft_count: number
  submitted_count: number
  approved_count: number
  rejected_count: number
  accepted_count: number
  in_force_count: number
  approval_rate: number
  acceptance_rate: number
  conversion_rate: number
  rejection_rate: number
  avg_time_to_submit_hours: number
  avg_time_to_approve_hours: number
  avg_time_to_accept_hours: number
  avg_total_cycle_hours: number
  sla_breach_count: number
  sla_transition_count: number
  sla_compliance_pct: number
  total_annual_premium: number
  pipeline_annual_premium: number
  avg_quote_value: number
  open_flags?: QuoteUserFlag[]
}

export interface FunnelStage {
  stage: string
  count: number
  avg_dwell_hours: number
}

export interface TrendBucket {
  bucket: string
  submitted: number
  approved: number
  accepted: number
  rejected: number
}

export interface SlaBreachByTransition {
  from_status: string
  to_status: string
  breach_count: number
  transition_count: number
  target_hours: number
}

export interface SlaBreachByUser {
  user_name: string
  breach_count: number
  transition_count: number
}

export interface SlaBreachSummary {
  breaches_by_transition: SlaBreachByTransition[]
  breaches_by_user: SlaBreachByUser[]
}

const defaultFilter: DashboardFilter = {
  from: null,
  to: null,
  users: [],
  region: [],
  quote_type: [],
  distribution_channel: []
}

// useQuoteDashboard is the single composable that backs the
// QuotePerformanceDashboard screen. Holds filter state (persisted across
// navigation via useFilterPersistence), exposes refresh helpers for each
// data slice, and surfaces loading / error refs for the UI.
export function useQuoteDashboard() {
  const { filters, resetFilters } = useFilterPersistence<DashboardFilter>(
    'quote-performance-dashboard',
    defaultFilter
  )

  const kpis = ref<QuotePerformanceKpis[]>([])
  const funnel = ref<FunnelStage[]>([])
  const trend = ref<TrendBucket[]>([])
  const sla = ref<SlaBreachSummary>({
    breaches_by_transition: [],
    breaches_by_user: []
  })
  const loading = ref(false)
  const error = ref<string | null>(null)
  const trendBucket = ref<'day' | 'week' | 'month'>('day')

  async function refreshAll() {
    loading.value = true
    error.value = null
    try {
      const [k, f, t, s] = await Promise.all([
        QuoteDashboardService.getKpis(filters.value),
        QuoteDashboardService.getFunnel(filters.value),
        QuoteDashboardService.getTrend(filters.value, trendBucket.value),
        QuoteDashboardService.getSlaBreaches(filters.value)
      ])
      kpis.value = k.data || []
      funnel.value = f.data || []
      trend.value = t.data || []
      sla.value = s.data || {
        breaches_by_transition: [],
        breaches_by_user: []
      }
    } catch (e: any) {
      error.value = e?.response?.data?.error || e?.message || String(e)
    } finally {
      loading.value = false
    }
  }

  return {
    filters,
    resetFilters,
    kpis,
    funnel,
    trend,
    sla,
    loading,
    error,
    trendBucket,
    refreshAll
  }
}

// useQuoteExtract drives the management extract grid: filters, paged
// rows, total count, and the xlsx-download trigger.
export function useQuoteExtract() {
  const filter = ref<QuoteExtractFilter>({
    page: 1,
    page_size: 50,
    order_by: 'creation_date desc',
    status: [],
    region: [],
    created_by: [],
    quote_type: [],
    industry: [],
    distribution_channel: []
  })
  const rows = ref<any[]>([])
  const total = ref(0)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const exporting = ref(false)

  async function refresh() {
    loading.value = true
    error.value = null
    try {
      const resp = await QuoteDashboardService.extract(filter.value)
      rows.value = resp.data.rows || []
      total.value = resp.data.total || 0
    } catch (e: any) {
      error.value = e?.response?.data?.error || e?.message || String(e)
    } finally {
      loading.value = false
    }
  }

  async function downloadXlsx() {
    exporting.value = true
    error.value = null
    try {
      const resp = await QuoteDashboardService.downloadExtractXlsx(filter.value)
      const blob = new Blob([resp.data], {
        type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
      })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `quote-performance-extract-${new Date()
        .toISOString()
        .replace(/[:.]/g, '-')}.xlsx`
      document.body.appendChild(a)
      a.click()
      a.remove()
      URL.revokeObjectURL(url)
    } catch (e: any) {
      const status = e?.response?.status
      if (status === 413) {
        error.value =
          'Extract exceeds 50,000 rows. Please narrow the filter and try again.'
      } else {
        error.value = e?.response?.data?.error || e?.message || String(e)
      }
    } finally {
      exporting.value = false
    }
  }

  return {
    filter,
    rows,
    total,
    loading,
    error,
    exporting,
    refresh,
    downloadXlsx
  }
}

// useSlaTargets drives the SLA target settings screen.
export function useSlaTargets() {
  const targets = ref<QuoteSlaTarget[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function refresh() {
    loading.value = true
    error.value = null
    try {
      const resp = await QuoteDashboardService.getSlaTargets()
      targets.value = resp.data || []
    } catch (e: any) {
      error.value = e?.response?.data?.error || e?.message || String(e)
    } finally {
      loading.value = false
    }
  }

  async function save(target: QuoteSlaTarget) {
    error.value = null
    try {
      await QuoteDashboardService.upsertSlaTarget(target)
      await refresh()
    } catch (e: any) {
      error.value = e?.response?.data?.error || e?.message || String(e)
      throw e
    }
  }

  async function remove(id: number) {
    error.value = null
    try {
      await QuoteDashboardService.deleteSlaTarget(id)
      await refresh()
    } catch (e: any) {
      error.value = e?.response?.data?.error || e?.message || String(e)
      throw e
    }
  }

  return { targets, loading, error, refresh, save, remove }
}

// useUserFlags drives the User Flags admin screen and the inline
// open/resolve actions on the dashboard's leaderboard card. Mirrors
// the SLA-target composable shape so the screen-level glue is
// familiar to anyone who has touched that feature.
export function useUserFlags() {
  const flags = ref<QuoteUserFlag[]>([])
  const loading = ref(false)
  const saving = ref(false)
  const error = ref<string | null>(null)
  const filter = ref<UserFlagsFilter>({ status: 'open' })

  async function refresh() {
    loading.value = true
    error.value = null
    try {
      const resp = await QuoteDashboardService.listUserFlags(filter.value)
      flags.value = resp.data || []
    } catch (e: any) {
      error.value = e?.response?.data?.error || e?.message || String(e)
    } finally {
      loading.value = false
    }
  }

  async function openFlag(body: OpenUserFlagBody) {
    saving.value = true
    error.value = null
    try {
      const resp = await QuoteDashboardService.openUserFlag(body)
      return resp.data as QuoteUserFlag
    } catch (e: any) {
      error.value = e?.response?.data?.error || e?.message || String(e)
      throw e
    } finally {
      saving.value = false
    }
  }

  async function resolveFlag(id: number, resolutionNote: string) {
    saving.value = true
    error.value = null
    try {
      const resp = await QuoteDashboardService.resolveUserFlag(
        id,
        resolutionNote
      )
      return resp.data as QuoteUserFlag
    } catch (e: any) {
      error.value = e?.response?.data?.error || e?.message || String(e)
      throw e
    } finally {
      saving.value = false
    }
  }

  return {
    flags,
    loading,
    saving,
    error,
    filter,
    refresh,
    openFlag,
    resolveFlag
  }
}
