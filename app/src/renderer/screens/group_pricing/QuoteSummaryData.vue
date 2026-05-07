<template>
  <v-container>
    <v-row class="mb-4" dense>
      <v-col v-for="item in keyInfo" :key="item.label" cols="4" md="2">
        <v-card variant="outlined" height="100%" class="summary-tile">
          <v-card-text class="tile-body">
            <div class="tile-label-row">
              <span class="text-overline text-grey-darken-1 tile-label">{{
                item.label
              }}</span>
              <v-tooltip v-if="item.info" location="bottom" max-width="420">
                <template #activator="{ props: tipProps }">
                  <v-icon
                    v-bind="tipProps"
                    size="small"
                    color="info"
                    class="tile-info-icon"
                  >
                    mdi-information-outline
                  </v-icon>
                </template>
                <div>
                  <template v-if="fclSourceLabel">
                    <div class="font-weight-medium mb-1">FCL source</div>
                    <div class="text-caption mb-1">{{ fclSourceLabel }}</div>
                    <div
                      v-if="overrideFcl > 0 && overrideFcl > appliedFcl"
                      class="text-caption mb-2"
                    >
                      Original override: {{ quoteData.currency }}
                      {{ overrideFcl.toLocaleString() }}.
                    </div>
                    <v-divider class="my-2"></v-divider>
                  </template>
                  <div class="font-weight-medium mb-1">
                    Maximum cover allowed per benefit
                  </div>
                  <div class="text-caption mb-2">
                    The Free Cover Limit is also subject to these per-benefit
                    caps, computed as the lower of the restriction limit and
                    the reinsurer limit (which depends on scheme size).
                  </div>
                  <table v-if="coverCaps.length" class="caps-table">
                    <tr v-for="r in coverCaps" :key="r.benefit">
                      <td class="pe-3">{{ r.benefit }}</td>
                      <td class="text-right">{{ r.cap }}</td>
                    </tr>
                  </table>
                  <div v-else class="text-caption">
                    No per-benefit caps recorded on this quote yet.
                  </div>
                </div>
              </v-tooltip>
            </div>
            <div class="quote-summary-value">{{ item.value }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <!-- Win Probability card -->
      <v-col cols="4" md="2">
        <v-card variant="outlined" height="100%" class="summary-tile">
          <v-card-text class="tile-body">
            <div class="tile-label-row">
              <span class="text-overline text-grey-darken-1 tile-label"
                >Win Probability</span
              >
            </div>
            <div class="quote-summary-value">
              <ProbabilityBadge
                :score="winProb"
                :loading="winProbLoading"
                size="small"
              />
            </div>
            <div v-if="topFeatures.length" class="mt-2 d-flex flex-wrap gap-1">
              <v-chip
                v-for="feat in topFeatures"
                :key="feat.name"
                size="x-small"
                :color="feat.contribution >= 0 ? 'success' : 'error'"
                variant="tonal"
              >
                {{ feat.name.replace(/_/g, ' ') }}
              </v-chip>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12" md="8">
        <v-card>
          <v-card-title class="text-h6">Scheme & Financials</v-card-title>
          <v-divider></v-divider>
          <v-card-text>
            <div class="text-subtitle-1 font-weight-bold text-primary mb-3"
              >Scheme Details</div
            >
            <v-row dense>
              <v-col
                v-for="detail in schemeDetails"
                :key="detail.label"
                cols="12"
                sm="6"
              >
                <div class="d-flex justify-space-between py-1">
                  <span class="text-grey-darken-1 mx-3">{{
                    detail.label
                  }}</span>
                  <span class="font-weight-medium">{{ detail.value }}</span>
                </div>
              </v-col>
            </v-row>

            <v-divider class="my-6"></v-divider>

            <div class="text-subtitle-1 font-weight-bold text-primary mb-3"
              >Configuration</div
            >
            <v-row dense>
              <v-col
                v-for="config in configuration"
                :key="config.label"
                cols="12"
                sm="6"
              >
                <div class="d-flex justify-space-between py-1">
                  <span class="text-grey-darken-1 mx-3">{{
                    config.label
                  }}</span>
                  <span class="font-weight-medium">{{ config.value }}</span>
                </div>
              </v-col>
            </v-row>

            <v-divider class="my-6"></v-divider>

            <div class="text-subtitle-1 font-weight-bold text-primary mb-3"
              >Financials & Loadings</div
            >
            <v-row dense>
              <v-col
                v-for="fin in financials"
                :key="fin.label"
                cols="12"
                sm="6"
              >
                <div class="d-flex justify-space-between py-1">
                  <span class="text-grey-darken-1 mx-3">{{ fin.label }}</span>
                  <span class="font-weight-medium">{{ fin.value }}</span>
                </div>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="4">
        <v-card>
          <v-card-title class="text-h6">Benefit Structure</v-card-title>
          <v-divider></v-divider>
          <v-card-text>
            <v-select
              v-model="selectedCategory"
              :items="categoryItems"
              label="Scheme Category"
              variant="outlined"
              density="compact"
              hide-details
              class="mb-4"
              @update:model-value="updateBenefitsActivated"
            ></v-select>
            <div
              v-if="selectedCategory && selectedCategoryRegion"
              class="d-flex align-center mb-3 px-1"
            >
              <v-icon size="small" class="mr-2 text-grey-darken-1"
                >mdi-map-marker</v-icon
              >
              <span class="text-grey-darken-1 text-body-2">Region:</span>
              <span class="font-weight-medium text-body-2 ml-2">{{
                selectedCategoryRegion
              }}</span>
            </div>
            <v-list lines="one">
              <div v-for="(benefit, i) in benefits" :key="benefit.name">
                <v-list-item>
                  <v-list-item-title class="font-weight-medium">{{
                    benefit.name
                  }}</v-list-item-title>
                  <template #append>
                    <v-chip
                      :color="benefit.active ? 'green' : 'grey'"
                      :text="benefit.active ? 'Active' : 'Inactive'"
                      size="small"
                      variant="tonal"
                    ></v-chip>
                  </template>
                </v-list-item>
                <v-divider v-if="Number(i) < benefits.length - 1"></v-divider>
              </div>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
// import BaseCard from '@/renderer/components/BaseCard.vue'
import formatDateString from '@/renderer/utils/helpers'
import { computed, onMounted, ref } from 'vue'
import ProbabilityBadge from '@/renderer/components/ProbabilityBadge.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

const props = defineProps({
  quote: {
    type: Object,
    required: true
  },
  resultSummaries: {
    type: Array,
    default: () => []
  }
})

const winProb = ref<number | null>(null)
const winProbLoading = ref(true)
const topFeatures = ref<{ name: string; contribution: number }[]>([])

onMounted(() => {
  if (props.quote?.id) {
    GroupPricingService.getQuoteWinProbability(props.quote.id)
      .then((res) => {
        const data = res.data?.data
        if (data) {
          winProb.value = data.score_pct
          try {
            topFeatures.value = JSON.parse(data.top_features || '[]')
          } catch {
            topFeatures.value = []
          }
        }
      })
      .catch(() => {})
      .finally(() => {
        winProbLoading.value = false
      })
  } else {
    winProbLoading.value = false
  }
})

// Mocking stuff
const quoteData = ref({
  schemeName: props.quote.scheme_name || 'Not Provided',
  industry: props.quote.industry || 'Not Provided',
  distributionChannel: props.quote.distribution_channel || 'broker',
  broker: props.quote.quote_broker?.name || 'N/A',
  startDate: props.quote.commencement_date
    ? formatDateString(props.quote.commencement_date, true, true, true)
    : 'Not Provided',
  contact: {
    name: props.quote.scheme_contact || 'Not Provided',
    email: props.quote.scheme_email || 'Not Provided'
  },
  quoteType: props.quote.quote_type || 'Not Provided',
  obligationType: props.quote.obligation_type || 'Not Provided',
  retirementAge: props.quote.normal_retirement_age || 0,
  experienceRating: props.quote.experience_rating || 'No',
  basis: props.quote.basis || 'Not Provided',
  globalSalaryMultiple: props.quote.use_global_salary_multiple ? 'Yes' : 'No',
  currency: props.quote.currency || 'ZAR'
})

// Loading rates displayed in the "Financials & Loadings" panel are taken
// from the MemberRatingResultSummary rows so they reflect what the rating
// engine actually applied. Loading fractions and final_scheme_total_commission
// are scheme-level totals mirrored onto every category, so the first row is
// sufficient for those. final_total_annual_premium is PER-CATEGORY, so the
// scheme premium is the sum across all summaries — this is the denominator
// for the implied commission rate. The summary stores loadings as fractions
// (0.05 == 5%) and discount as a negative fraction; quote.loadings.* (the
// input form) stores percentages (5 == 5%), which is why the fallback path
// skips the *100.
const round2 = (n: number) => Math.round(n * 100) / 100
const firstSummary = computed<any>(
  () => (props.resultSummaries as any[])?.[0] ?? null
)
const schemeFinalTotalAnnualPremium = computed(() =>
  (props.resultSummaries as any[]).reduce(
    (acc, s) => acc + (s?.final_total_annual_premium || 0),
    0
  )
)

const impliedCommission = computed(() => {
  const s = firstSummary.value
  const denom = schemeFinalTotalAnnualPremium.value
  if (s && denom > 0) {
    return round2((s.final_scheme_total_commission / denom) * 100)
  }
  return props.quote.loadings?.commission_loading ?? 0
})

const expenseLoading = computed(() => {
  const s = firstSummary.value
  return s
    ? round2(s.expense_loading * 100)
    : (props.quote.loadings?.expense_loading ?? 0)
})

const adminLoading = computed(() => {
  const s = firstSummary.value
  return s
    ? round2(s.admin_loading * 100)
    : (props.quote.loadings?.admin_loading ?? 0)
})

const otherLoading = computed(() => {
  const s = firstSummary.value
  return s
    ? round2(s.other_loading * 100)
    : (props.quote.loadings?.other_loading ?? 0)
})

const profitLoading = computed(() => {
  const s = firstSummary.value
  return s
    ? round2(s.profit_loading * 100)
    : (props.quote.loadings?.profit_loading ?? 0)
})

const overallDiscount = computed(() => {
  const s = firstSummary.value
  // Summary stores discount as a negative fraction; flip sign for display.
  return s ? round2(-s.discount * 100) : (props.quote.loadings?.discount ?? 0)
})

const benefits: any = ref([])

const selectedCategory = ref(null)
const selectedCategoryRegion = ref<string>('')

const updateBenefitsActivated = () => {
  // This function would update the benefits based on the selected category
  // For now, it just logs the selected category

  const schemeCategory = props.quote.scheme_categories.find(
    (cat: any) => cat.scheme_category === selectedCategory.value
  )
  benefits.value = []
  selectedCategoryRegion.value = ''
  if (schemeCategory) {
    selectedCategoryRegion.value = schemeCategory.region || ''
    benefits.value.push({
      name: schemeCategory.gla_alias || 'GLA',
      active: schemeCategory.gla_benefit
    })
    benefits.value.push({
      name: schemeCategory.sgla_alias || 'SGLA',
      active: schemeCategory.sgla_benefit
    })
    benefits.value.push({
      name: schemeCategory.phi_alias || 'PHI',
      active: schemeCategory.phi_benefit
    })
    benefits.value.push({
      name: schemeCategory.ci_alias || 'CI',
      active: schemeCategory.ci_benefit
    })
    benefits.value.push({
      name: schemeCategory.ptd_alias || 'PTD',
      active: schemeCategory.ptd_benefit
    })
    benefits.value.push({
      name: schemeCategory.ttd_alias || 'TTD',
      active: schemeCategory.ttd_benefit
    })
    benefits.value.push({
      name: schemeCategory.family_funeral_alias || 'GFF',
      active: schemeCategory.family_funeral_benefit
    })
  }
}

// --- COMPUTED PROPERTIES ---
// Using computed properties makes the template cleaner and more readable.
const distributionChannelLabel = (channel: string): string => {
  switch (channel) {
    case 'broker':
      return 'Broker'
    case 'direct':
      return 'Direct'
    case 'binder':
      return 'Binder'
    case 'tied_agent':
      return 'Tied Agent'
    default:
      return channel || '—'
  }
}

// Pick the smaller positive cap, treating 0 as "no cap on this side".
// Per-benefit max cover allowed is min(restriction, reinsurer); if either
// half is unconfigured (0) the other one wins; if both are 0 the benefit
// has no cap on this scheme.
const minNonZero = (a: number, b: number): number => {
  if (a > 0 && b > 0) return Math.min(a, b)
  if (a > 0) return a
  if (b > 0) return b
  return 0
}

// Per-benefit caps shown in the FCL info tooltip. Sourced from the first
// MemberRatingResultSummary row (scheme-level fields are mirrored across
// every category row, same convention used by `firstSummary` below).
const coverCaps = computed(() => {
  const s: any = (props.resultSummaries as any[])[0] || {}
  const ccy = quoteData.value.currency
  const fmt = (n: number, suffix = ''): string =>
    `${ccy} ${n.toLocaleString()}${suffix}`

  const rows: { benefit: string; cap: string }[] = []
  const push = (
    benefit: string,
    restriction: number,
    reins: number,
    suffix = ''
  ) => {
    const cap = minNonZero(Number(restriction) || 0, Number(reins) || 0)
    if (cap > 0) rows.push({ benefit, cap: fmt(cap, suffix) })
  }

  push('GLA', s.maximum_gla_cover, s.reins_max_gla_cover)
  push('PTD', s.maximum_ptd_cover, s.reins_max_ptd_cover)
  push('CI', s.severe_illness_maximum_benefit, s.reins_max_ci_cover)
  push('Spouse GLA', s.spouse_gla_maximum_benefit, s.reins_max_sgla_cover)
  push('TTD', s.ttd_maximum_monthly_benefit, s.reins_max_ttd_cover, '/mo')
  push('PHI', s.phi_maximum_monthly_benefit, s.reins_max_phi_cover, '/mo')
  push('Funeral', 0, s.reins_max_fun_cover)

  return rows
})

// FCL is computed/applied per scheme_category by the rating engine. The
// quote-level `free_cover_limit` is the optional user override (only set
// when "Enforce FCL" was ticked on the General Input form); the engine
// always populates `scheme_categories[i].free_cover_limit` regardless,
// clamping the override to MaximumAllowedFCL when it exceeds the cap.
// We trust the per-category applied value as the source of truth.
const overrideFcl = computed<number>(
  () => Number(props.quote.free_cover_limit) || 0
)

const selectedSchemeCategoryData = computed<any | null>(() => {
  if (!selectedCategory.value) return null
  return (
    (props.quote.scheme_categories as any[])?.find(
      (cat: any) => cat.scheme_category === selectedCategory.value
    ) ?? null
  )
})

const appliedFcl = computed<number>(() => {
  const cat = selectedSchemeCategoryData.value
  return cat ? Number(cat.free_cover_limit) || 0 : 0
})

const fclSourceLabel = computed<string>(() => {
  if (!selectedSchemeCategoryData.value || appliedFcl.value <= 0) return ''
  if (overrideFcl.value > 0) {
    if (overrideFcl.value > appliedFcl.value) {
      return 'User-enforced override (clamped to maximum allowed FCL)'
    }
    return 'User-enforced override'
  }
  return 'Calculated by FCL method (percentile / outlier)'
})

const fclDisplayValue = computed<string>(() => {
  if (!selectedCategory.value) return 'Select a category'
  if (appliedFcl.value > 0) {
    return `${quoteData.value.currency} ${appliedFcl.value.toLocaleString()}`
  }
  return 'No FCL enforced'
})

interface KeyInfoItem {
  label: string
  value: string
  info?: boolean
}

const keyInfo = computed<KeyInfoItem[]>(() => [
  { label: 'Industry', value: quoteData.value.industry },
  {
    label: 'Distribution Channel',
    value: distributionChannelLabel(quoteData.value.distributionChannel)
  },
  { label: 'Broker', value: quoteData.value.broker },
  { label: 'Start Date', value: quoteData.value.startDate },
  {
    label: 'Free Cover Limit',
    value: fclDisplayValue.value,
    info: true
  }
])

const schemeDetails = computed(() => [
  { label: 'Scheme Contact', value: quoteData.value.contact.name },
  { label: 'Contact Email', value: quoteData.value.contact.email },
  { label: 'Quote Type', value: quoteData.value.quoteType },
  { label: 'Obligation Type', value: quoteData.value.obligationType }
])

const configuration = computed(() => [
  { label: 'Normal Retirement Age', value: quoteData.value.retirementAge },
  { label: 'Experience Rating', value: quoteData.value.experienceRating },
  { label: 'Basis', value: quoteData.value.basis },
  {
    label: 'Global Salary Multiple',
    value: quoteData.value.globalSalaryMultiple
  }
])

const financials = computed(() => [
  { label: 'Currency', value: quoteData.value.currency },
  { label: 'Implied Commission', value: `${impliedCommission.value}%` },
  { label: 'Expense Loading', value: `${expenseLoading.value}%` },
  { label: 'Admin Loading', value: `${adminLoading.value}%` },
  { label: 'Other Loading', value: `${otherLoading.value}%` },
  { label: 'Profit Loading', value: `${profitLoading.value}%` },
  { label: 'Overall Premium Discount', value: `${overallDiscount.value}%` }
])

const categoryItems = computed(() => {
  // This could be dynamic based on the quote data, but for now, it's static

  const categories = props.quote.scheme_categories.map(
    (cat: any) => cat.scheme_category
  )
  return categories.length > 0
    ? categories
    : ['Management', 'General Staff', 'Executive']
})

//

// const dashIfEmpty = (value: any) => value || '–'
</script>

<style scoped>
/* Top "Quote Summary" tiles. Layout target: 3 cards per row on small
   screens (two rows total), 6 cards per row on md+ (one row). Card height
   is unified so neighbours match regardless of label/value length, and
   the label/value/extras stack predictably via flex inside v-card-text. */
.summary-tile {
  border-radius: 10px;
  transition:
    box-shadow 0.18s ease,
    transform 0.18s ease,
    border-color 0.18s ease;
}
.summary-tile:hover {
  box-shadow: 0 6px 14px rgba(0, 0, 0, 0.06);
  border-color: rgba(0, 0, 0, 0.18);
  transform: translateY(-1px);
}

.tile-body {
  display: flex;
  flex-direction: column;
  height: 100%;
  gap: 6px;
  padding: 14px 16px;
}

.tile-label-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 6px;
  min-height: 32px;
}

.tile-label {
  flex: 1;
  line-height: 1.25;
  letter-spacing: 0.6px;
  word-break: break-word;
}

.tile-info-icon {
  flex-shrink: 0;
  margin-top: 2px;
}

/* The min-height matches the v-chip used by ProbabilityBadge (small chip ≈
   24px) so the chip and plain text values sit at the same vertical centre. */
.quote-summary-value {
  display: flex;
  align-items: center;
  min-height: 28px;
  font-weight: 500;
  font-size: 0.95rem;
  line-height: 1.25;
  word-break: break-word;
}

.caps-table td {
  padding: 2px 0;
  white-space: nowrap;
}
</style>
