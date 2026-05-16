import Api from '@/renderer/api/Api'

// Filter envelope shared by the dashboard's KPI / funnel / trend / SLA
// endpoints. Empty arrays and nulls mean "no filter on this dimension".
export interface DashboardFilter {
  from?: string | null
  to?: string | null
  users?: string[]
  region?: string[]
  quote_type?: string[]
  distribution_channel?: string[]
}

// QuoteExtractFilter mirrors the Go struct of the same name. Extends the
// dashboard filter with extract-specific fields (premium range, status
// multi-select, paging).
export interface QuoteExtractFilter extends DashboardFilter {
  created_by?: string[]
  reviewer?: string[]
  status?: string[]
  industry?: string[]
  min_annual_premium?: number | null
  max_annual_premium?: number | null
  page?: number
  page_size?: number
  order_by?: string
}

export interface QuoteSlaTarget {
  id?: number
  from_status: string
  to_status: string
  target_hours: number
  warning_pct_of_sla?: number
  quote_type?: string
  active?: boolean
  updated_by?: string
  updated_at?: string
}

// QuoteUserFlag mirrors the Go struct. Internal `note` and
// `resolution_note` come back empty when the caller lacks
// quote:manage_user_flags — the dashboard's chip rendering only needs
// reason + open/resolved state, so non-managers still get the chip.
export interface QuoteUserFlag {
  id: number
  user_name: string
  user_email: string
  flag_reason: 'coaching' | 'capacity'
  note: string
  opened_by: string
  opened_by_name: string
  opened_at: string
  resolved_by?: string | null
  resolved_by_name?: string | null
  resolved_at?: string | null
  resolution_note?: string | null
}

export interface UserFlagsFilter {
  status?: 'open' | 'resolved' | 'all'
  user_name?: string
  reason?: 'coaching' | 'capacity'
}

export interface OpenUserFlagBody {
  user_name: string
  user_email?: string
  flag_reason: 'coaching' | 'capacity'
  note: string
}

// toParams converts a DashboardFilter into the URLSearchParams-style
// shape axios expects for query strings, expanding array fields into
// repeated keys (users=alice&users=bob) — matching gin's c.QueryArray.
function toParams(filter: DashboardFilter) {
  const params: Record<string, unknown> = {}
  for (const [key, value] of Object.entries(
    filter as Record<string, unknown>
  )) {
    if (value === null || value === undefined) continue
    if (Array.isArray(value)) {
      if (value.length === 0) continue
      params[key] = value
    } else {
      params[key] = value
    }
  }
  return params
}

export default {
  getKpis(filter: DashboardFilter) {
    return Api.get('/group-pricing/dashboard/kpis', {
      params: toParams(filter),
      paramsSerializer: { indexes: null }
    })
  },

  getFunnel(filter: DashboardFilter) {
    return Api.get('/group-pricing/dashboard/funnel', {
      params: toParams(filter),
      paramsSerializer: { indexes: null }
    })
  },

  getTrend(filter: DashboardFilter, bucket: 'day' | 'week' | 'month' = 'day') {
    return Api.get('/group-pricing/dashboard/trend', {
      params: { ...toParams(filter), bucket },
      paramsSerializer: { indexes: null }
    })
  },

  getSlaBreaches(filter: DashboardFilter) {
    return Api.get('/group-pricing/dashboard/sla-breaches', {
      params: toParams(filter),
      paramsSerializer: { indexes: null }
    })
  },

  // Extract grid: POST because the filter payload (multi-select arrays)
  // can exceed URL length on large selections.
  extract(filter: QuoteExtractFilter) {
    return Api.post('/group-pricing/dashboard/extract', filter)
  },

  // Build a download URL the browser can hit via window.open or an
  // <a href download> tag. The xlsx endpoint reads the filter from
  // query params, so we serialise the same way as the GET endpoints.
  buildExtractXlsxUrl(filter: QuoteExtractFilter): string {
    const sp = new URLSearchParams()
    const addArray = (key: string, vals?: string[] | null) => {
      if (!vals) return
      for (const v of vals) sp.append(key, v)
    }
    addArray('created_by', filter.created_by)
    addArray('reviewer', filter.reviewer)
    addArray('region', filter.region)
    addArray('quote_type', filter.quote_type)
    addArray('industry', filter.industry)
    addArray('distribution_channel', filter.distribution_channel)
    addArray('status', filter.status)
    if (filter.min_annual_premium != null) {
      sp.set('min_annual_premium', String(filter.min_annual_premium))
    }
    if (filter.max_annual_premium != null) {
      sp.set('max_annual_premium', String(filter.max_annual_premium))
    }
    if (filter.from) sp.set('from', filter.from)
    if (filter.to) sp.set('to', filter.to)
    if (filter.order_by) sp.set('order_by', filter.order_by)
    return `/group-pricing/dashboard/extract.xlsx?${sp.toString()}`
  },

  // Download xlsx as a Blob so callers can wire it into a file-save
  // dialog without leaving the SPA. Returns the raw axios response.
  downloadExtractXlsx(filter: QuoteExtractFilter) {
    return Api.get(this.buildExtractXlsxUrl(filter), { responseType: 'blob' })
  },

  getQuoteStatusHistory(quoteId: number | string) {
    return Api.get(`/group-pricing/quotes/${quoteId}/status-history`)
  },

  getSlaTargets() {
    return Api.get('/group-pricing/dashboard/sla-targets')
  },

  upsertSlaTarget(target: QuoteSlaTarget) {
    return Api.post('/group-pricing/dashboard/sla-targets', target)
  },

  deleteSlaTarget(id: number) {
    return Api.delete(`/group-pricing/dashboard/sla-targets/${id}`)
  },

  listUserFlags(filter: UserFlagsFilter = {}) {
    return Api.get('/group-pricing/dashboard/user-flags', { params: filter })
  },

  openUserFlag(body: OpenUserFlagBody) {
    return Api.post('/group-pricing/dashboard/user-flags', body)
  },

  resolveUserFlag(id: number, resolutionNote: string) {
    return Api.post(`/group-pricing/dashboard/user-flags/${id}/resolve`, {
      resolution_note: resolutionNote
    })
  }
}
