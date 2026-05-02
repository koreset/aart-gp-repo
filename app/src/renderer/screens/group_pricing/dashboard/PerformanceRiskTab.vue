<template>
  <div>
    <!-- Banner -->
    <v-alert
      type="info"
      variant="tonal"
      color="primary"
      density="compact"
      class="mb-4"
      icon="mdi-chart-line"
    >
      <strong>Performance &amp; Risk</strong> — Per-scheme view of the in-force
      book. Loss ratios are inception-to-date (claims since each scheme's cover
      start ÷ ITD earned premium).
    </v-alert>

    <div v-if="loading" class="d-flex justify-center my-8">
      <v-progress-circular indeterminate color="primary" />
    </div>

    <template v-else>
      <!-- Headline KPI cards -->
      <v-row class="d-flex justify-center">
        <v-col v-for="card in kpiCards" :key="card.title" cols="6" md="2">
          <v-card
            variant="tonal"
            :color="card.color || 'primary'"
            class="dash-card"
          >
            <v-card-subtitle class="d-flex align-center">
              <h5 class="flex-grow-1">{{ card.title }}</h5>
              <v-tooltip location="top" max-width="320">
                <template #activator="{ props: tipProps }">
                  <v-icon
                    v-bind="tipProps"
                    size="x-small"
                    icon="mdi-information-outline"
                    class="ml-1"
                    style="cursor: help"
                  />
                </template>
                <span>{{ card.tooltip }}</span>
              </v-tooltip>
            </v-card-subtitle>
            <v-card-text>
              <v-row>
                <v-col class="d-flex justify-center">
                  <h2>{{ card.value }}</h2>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-divider class="my-6" />

      <!-- Refresh + section title for table -->
      <v-row class="align-center mb-2">
        <v-col class="d-flex align-center">
          <span class="text-subtitle-1 font-weight-bold"
            >Scheme Performance</span
          >
          <v-tooltip location="top" max-width="360">
            <template #activator="{ props: tipProps }">
              <v-icon
                v-bind="tipProps"
                size="small"
                icon="mdi-information-outline"
                class="ml-2"
                style="cursor: help"
              />
            </template>
            <span>
              One row per in-force scheme. ELR = Expected Loss Ratio (priced at
              quote time). ITD ALR = Actual Loss Ratio inception-to-date (claims
              paid ÷ earned premium since cover start). Δ = ALR − ELR; positive
              means the scheme is running worse than priced.
            </span>
          </v-tooltip>
        </v-col>
        <v-col cols="auto">
          <v-btn
            variant="tonal"
            color="primary"
            size="small"
            prepend-icon="mdi-refresh"
            :loading="loading"
            @click="loadAll"
          >
            Refresh
          </v-btn>
        </v-col>
      </v-row>

      <SchemePerformanceTable :rows="performance.rows" />

      <v-divider class="my-6" />

      <LossRatioTrend :trend="trend" />

      <v-divider class="my-6" />

      <LossRatioDistribution
        :buckets="profile.loss_ratio_buckets"
        :worst="profile.top10_worst"
      />

      <v-divider class="my-6" />

      <div class="d-flex align-center mb-3">
        <div class="text-subtitle-1 font-weight-bold">Risk Profile</div>
        <v-tooltip location="top" max-width="360">
          <template #activator="{ props: tipProps }">
            <v-icon
              v-bind="tipProps"
              size="small"
              icon="mdi-information-outline"
              class="ml-2"
              style="cursor: help"
            />
          </template>
          <span>
            Portfolio-level risk lens: where premium is concentrated, which
            industry × region segments are running hot, the
            frequency-vs-severity profile per scheme, and the watchlist of
            schemes whose ALR has breached the trigger thresholds.
          </span>
        </v-tooltip>
      </div>
      <RiskProfilePanel
        :profile="profile"
        :rows="performance.rows"
        :deteriorating="derivedDeteriorating"
        v-model:alr-ceiling="alrCeiling"
        v-model:alr-delta="alrDelta"
        :is-custom-view="isCustomView"
        @reset="resetThresholds"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import SchemePerformanceTable from './SchemePerformanceTable.vue'
import LossRatioDistribution from './LossRatioDistribution.vue'
import LossRatioTrend from './LossRatioTrend.vue'
import RiskProfilePanel from './RiskProfilePanel.vue'

const loading = ref(false)
const performance = ref<any>({
  rows: [],
  total_schemes: 0,
  total_premium: 0,
  total_itd_claims: 0,
  portfolio_alr: null,
  total_r12m_claims: 0,
  portfolio_r12m_alr: null
})
const trend = ref<any>({ months: [], portfolio: [], schemes: [] })
const profile = ref<any>({
  concentration: {
    top5_premium_share: 0,
    top10_premium_share: 0,
    hhi: 0,
    total_schemes: 0
  },
  pareto: [],
  loss_ratio_buckets: [],
  top10_worst: [],
  industry_region: [],
  frequency_severity: [],
  deteriorating: [],
  thresholds: { alr_ceiling_pct: 100, alr_delta_pp: 20 }
})

// Per-session threshold overrides — seeded from profile.thresholds (the
// company defaults from settings) and re-derived client-side as the analyst
// dials them. Resetting just snaps these back to the company defaults.
const alrCeiling = ref<number>(100)
const alrDelta = ref<number>(20)

const isCustomView = computed(
  () =>
    alrCeiling.value !== profile.value.thresholds.alr_ceiling_pct ||
    alrDelta.value !== profile.value.thresholds.alr_delta_pp
)

const resetThresholds = () => {
  alrCeiling.value = profile.value.thresholds.alr_ceiling_pct
  alrDelta.value = profile.value.thresholds.alr_delta_pp
}

// Re-derive the deteriorating watchlist locally from the per-scheme rows so
// dialing the thresholds gives instant feedback with no backend round-trip.
// The backend's profile.deteriorating remains the canonical company-default
// view but we don't render it directly — re-deriving keeps the UI list
// consistent with whatever the analyst has in the threshold inputs.
const derivedDeteriorating = computed(() => {
  const ceiling = alrCeiling.value
  const delta = alrDelta.value
  const out: any[] = []
  for (const r of performance.value.rows ?? []) {
    const reasons: string[] = []
    if (r.actual_loss_ratio != null && r.actual_loss_ratio > ceiling) {
      reasons.push(`ALR > ${ceiling}%`)
    }
    if (r.loss_ratio_delta != null && r.loss_ratio_delta > delta) {
      reasons.push(`ALR exceeds ELR by >${delta}pp`)
    }
    if (reasons.length === 0) continue
    out.push({
      scheme_id: r.scheme_id,
      scheme_name: r.scheme_name,
      trigger_reasons: reasons,
      expected_loss_ratio: r.expected_loss_ratio,
      actual_loss_ratio: r.actual_loss_ratio,
      loss_ratio_delta: r.loss_ratio_delta,
      itd_claims_paid: r.itd_claims_paid,
      last_claim_date: r.last_claim_date
    })
  }
  out.sort((a, b) => (b.actual_loss_ratio ?? 0) - (a.actual_loss_ratio ?? 0))
  return out
})

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

const kpiCards = computed(() => [
  {
    title: 'Schemes In Force',
    value: performance.value.total_schemes,
    tooltip:
      'Number of schemes currently active (in_force = true) on the Group Risk book.'
  },
  {
    title: 'Annual Premium',
    value: fmtCurrency(performance.value.total_premium),
    tooltip:
      'Sum of contracted annual premiums (GWP) across all in-force schemes.'
  },
  {
    title: 'ITD Claims Paid',
    value: fmtCurrency(performance.value.total_itd_claims),
    tooltip:
      'Inception-to-date claims paid: total of approved + paid claims summed since each scheme’s cover start date.'
  },
  {
    title: 'Portfolio ITD ALR',
    value: fmtPct(performance.value.portfolio_alr),
    color: portfolioAlrColor.value,
    tooltip:
      'Inception-to-date Actual Loss Ratio for the whole book = Σ ITD claims paid ÷ Σ ITD earned premium × 100. Green ≤ 80%, amber 80–100%, red > 100%.'
  },
  {
    title: 'Portfolio R12M ALR',
    value: fmtPct(performance.value.portfolio_r12m_alr),
    color: portfolioR12mAlrColor.value,
    tooltip:
      'Rolling-12-month Actual Loss Ratio for the whole book = Σ claims in last 12 months ÷ Σ time-weighted annual premium × 100. Compare to Portfolio ITD ALR — a higher R12M means the book is running hotter recently than its long-run history.'
  },
  {
    title: 'Deteriorating',
    value: derivedDeteriorating.value.length,
    color: derivedDeteriorating.value.length > 0 ? 'error' : 'success',
    tooltip: `Count of schemes flagged because their ITD ALR exceeds ${alrCeiling.value}% or exceeds Expected LR by more than ${alrDelta.value}pp. Adjust the thresholds below the deteriorating panel for what-if analysis; the company defaults are configured in MetaData.`
  },
  {
    title: 'Top-5 Share',
    value: fmtPct(profile.value.concentration.top5_premium_share),
    color: concentrationColor.value,
    tooltip:
      'Share of total annual premium contributed by the 5 largest schemes — a quick read on premium concentration risk.'
  }
])

const portfolioAlrColor = computed(() => {
  const a = performance.value.portfolio_alr
  if (a == null) return 'primary'
  if (a > 100) return 'error'
  if (a > 80) return 'warning'
  return 'success'
})

const portfolioR12mAlrColor = computed(() => {
  const a = performance.value.portfolio_r12m_alr
  if (a == null) return 'primary'
  if (a > 100) return 'error'
  if (a > 80) return 'warning'
  return 'success'
})

const concentrationColor = computed(() => {
  const t5 = profile.value.concentration.top5_premium_share
  if (t5 > 60) return 'error'
  if (t5 > 40) return 'warning'
  return 'primary'
})

const loadAll = async () => {
  loading.value = true
  try {
    const [perfRes, profileRes, trendRes] = await Promise.all([
      GroupPricingService.getSchemePerformance(),
      GroupPricingService.getRiskProfile(),
      GroupPricingService.getLossRatioTrend()
    ])
    performance.value = perfRes.data
    profile.value = profileRes.data
    trend.value = trendRes.data
    // Seed the per-session threshold inputs from the company defaults
    // returned by the backend on every refresh — this also acts as
    // "Reset to defaults" if the user re-fetches.
    if (profile.value?.thresholds) {
      alrCeiling.value = Number(profile.value.thresholds.alr_ceiling_pct)
      alrDelta.value = Number(profile.value.thresholds.alr_delta_pp)
    }
  } catch (err) {
    console.error('Failed to load performance & risk data', err)
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)
</script>

<style scoped>
.dash-card {
  min-height: 110px;
}
.dash-card h2 {
  font-size: 1.4rem;
  word-break: break-word;
  text-align: center;
}
</style>
