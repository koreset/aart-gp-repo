<template>
  <v-container>
    <v-row class="mb-4">
      <v-col v-for="item in keyInfo" :key="item.label" cols="12" sm="6" md="2">
        <v-card variant="outlined">
          <v-card-text>
            <div class="text-overline text-grey-darken-1">{{ item.label }}</div>
            <div class="quote-summary-value">{{ item.value }}</div>
          </v-card-text>
        </v-card>
      </v-col>
      <!-- Win Probability card -->
      <v-col cols="12" sm="6" md="2">
        <v-card variant="outlined">
          <v-card-text>
            <div class="text-overline text-grey-darken-1">Win Probability</div>
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
  freeCoverLimit: 300000,
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
  currency: props.quote.currency || 'ZAR',
  commissionRate: props.quote.loadings?.commission_loading || 0,
  expenseLoading: props.quote.loadings?.expense_loading || 0,
  profitLoading: props.quote.loadings?.profit_loading || 0,
  adminLoading: props.quote.loadings?.admin_loading || 0,
  contingencyLoading: props.quote.loadings?.contingency_loading || 0,
  otherLoading: props.quote.loadings?.other_loading || 0,
  overallDiscount: props.quote.loadings?.discount || 0
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
const keyInfo = computed(() => [
  { label: 'Industry', value: quoteData.value.industry },
  { label: 'Broker', value: quoteData.value.broker },
  { label: 'Start Date', value: quoteData.value.startDate },
  {
    label: 'Free Cover Limit',
    value: `${quoteData.value.currency} ${quoteData.value.freeCoverLimit.toLocaleString()}`
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
  { label: 'Commission Rate', value: `${quoteData.value.commissionRate}%` },
  { label: 'Expense Loading', value: `${quoteData.value.expenseLoading}%` },
  { label: 'Admin Loading', value: `${quoteData.value.adminLoading}%` },
  {
    label: 'Contingency Loading',
    value: `${quoteData.value.contingencyLoading}%`
  },
  { label: 'Other Loading', value: `${quoteData.value.otherLoading}%` },
  { label: 'Profit Loading', value: `${quoteData.value.profitLoading}%` },
  {
    label: 'Overall Premium Discount',
    value: `${quoteData.value.overallDiscount}%`
  }
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
/* Shared value row for the top "Quote Summary" cards. The min-height matches
   the v-chip used by ProbabilityBadge (small chip ≈ 24px) so the chip and
   plain text values sit at the same vertical centre. Without this the chip
   extends below the text baseline and looks misaligned next to its
   neighbours. */
.quote-summary-value {
  display: flex;
  align-items: center;
  min-height: 28px;
  font-weight: 500;
  line-height: 1.2;
}
</style>
