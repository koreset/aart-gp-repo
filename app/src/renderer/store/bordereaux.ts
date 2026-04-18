import { defineStore } from 'pinia'
import GroupPricingService from '@/renderer/api/GroupPricingService'

// Shape of the scheme rows returned by getSchemesInforcev2. Kept loose so
// consumers that only read a couple of fields don't pull in the full model.
interface SchemeSummary {
  id: number
  name: string
  [key: string]: any
}

interface BordereauxTemplate {
  id: number
  name: string
  type: string
  [key: string]: any
}

interface BordereauxConfiguration {
  id: number
  name: string
  [key: string]: any
}

interface BordereauxFieldDef {
  display_name: string
  field_name: string
}

interface BordereauxDashboardStats {
  generated_this_month: number
  pending_submissions: number
  reconciled_this_week: number
  active_templates: number
}

interface LoadOptions {
  force?: boolean
}

// Pinia store caching the read-heavy lookups shared across every screen in the
// bordereaux module. Prior to this store each component re-fetched schemes /
// templates / configurations on every mount; a user flipping between screens
// triggered 3-5 redundant round trips per minute. The store exposes loader
// actions that return the cached value unless `force: true` is passed, plus
// `invalidate*` helpers that components call after a successful mutation.
export const useBordereauxStore = defineStore('bordereaux', {
  state: () => ({
    schemes: [] as SchemeSummary[],
    schemesLoaded: false,
    schemesLoading: null as Promise<SchemeSummary[]> | null,

    templates: [] as BordereauxTemplate[],
    templatesLoaded: false,
    templatesLoading: null as Promise<BordereauxTemplate[]> | null,

    configurations: [] as BordereauxConfiguration[],
    configurationsLoaded: false,
    configurationsLoading: null as Promise<BordereauxConfiguration[]> | null,

    dashboardStats: null as BordereauxDashboardStats | null,
    dashboardStatsLoading: null as Promise<BordereauxDashboardStats | null> | null,

    // Field catalogue is type-scoped (member | premium | claim) so we cache by key.
    fieldsByType: {} as Record<string, BordereauxFieldDef[]>,
    fieldsByTypeLoading: {} as Record<string, Promise<BordereauxFieldDef[]> | null>
  }),

  actions: {
    /**
     * Return the cached scheme list, fetching on first call. Concurrent callers
     * during the in-flight fetch share the same promise so we never fan out
     * duplicate requests.
     */
    async loadSchemes(opts: LoadOptions = {}): Promise<SchemeSummary[]> {
      if (!opts.force && this.schemesLoaded) return this.schemes
      if (this.schemesLoading) return this.schemesLoading
      this.schemesLoading = (async () => {
        try {
          const res = await GroupPricingService.getSchemesInforcev2()
          this.schemes = (res.data ?? []) as SchemeSummary[]
          this.schemesLoaded = true
          return this.schemes
        } finally {
          this.schemesLoading = null
        }
      })()
      return this.schemesLoading
    },

    async loadTemplates(opts: LoadOptions = {}): Promise<BordereauxTemplate[]> {
      if (!opts.force && this.templatesLoaded) return this.templates
      if (this.templatesLoading) return this.templatesLoading
      this.templatesLoading = (async () => {
        try {
          const res = await GroupPricingService.getBordereauxTemplates()
          this.templates = (res.data ?? []) as BordereauxTemplate[]
          this.templatesLoaded = true
          return this.templates
        } finally {
          this.templatesLoading = null
        }
      })()
      return this.templatesLoading
    },

    async loadConfigurations(
      opts: LoadOptions = {}
    ): Promise<BordereauxConfiguration[]> {
      if (!opts.force && this.configurationsLoaded) return this.configurations
      if (this.configurationsLoading) return this.configurationsLoading
      this.configurationsLoading = (async () => {
        try {
          const res = await GroupPricingService.getBordereauxConfigurations()
          this.configurations = (res.data ?? []) as BordereauxConfiguration[]
          this.configurationsLoaded = true
          return this.configurations
        } finally {
          this.configurationsLoading = null
        }
      })()
      return this.configurationsLoading
    },

    async loadDashboardStats(
      opts: LoadOptions = {}
    ): Promise<BordereauxDashboardStats | null> {
      if (!opts.force && this.dashboardStats) return this.dashboardStats
      if (this.dashboardStatsLoading) return this.dashboardStatsLoading
      this.dashboardStatsLoading = (async () => {
        try {
          const res = await GroupPricingService.getBordereauxDashboardStats()
          this.dashboardStats = res.data as BordereauxDashboardStats
          return this.dashboardStats
        } finally {
          this.dashboardStatsLoading = null
        }
      })()
      return this.dashboardStatsLoading
    },

    async loadFieldsByType(
      bordereauType: string,
      opts: LoadOptions = {}
    ): Promise<BordereauxFieldDef[]> {
      if (!bordereauType) return []
      if (!opts.force && this.fieldsByType[bordereauType]) {
        return this.fieldsByType[bordereauType]
      }
      if (this.fieldsByTypeLoading[bordereauType]) {
        return this.fieldsByTypeLoading[bordereauType] as Promise<
          BordereauxFieldDef[]
        >
      }
      const pending = (async () => {
        try {
          const res = await GroupPricingService.getBordereauxFields(
            bordereauType
          )
          const fields = (res.data ?? []) as BordereauxFieldDef[]
          this.fieldsByType[bordereauType] = fields
          return fields
        } finally {
          this.fieldsByTypeLoading[bordereauType] = null
        }
      })()
      this.fieldsByTypeLoading[bordereauType] = pending
      return pending
    },

    // Invalidators — call after a successful mutation so the next read reloads
    // from the server instead of serving a stale cache.
    invalidateTemplates() {
      this.templatesLoaded = false
      this.templates = []
    },
    invalidateConfigurations() {
      this.configurationsLoaded = false
      this.configurations = []
    },
    invalidateDashboardStats() {
      this.dashboardStats = null
    },
    invalidateAll() {
      this.$reset()
    }
  }
})
