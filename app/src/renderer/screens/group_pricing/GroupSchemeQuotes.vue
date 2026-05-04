<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <span class="headline">Quote List</span>
          </template>
          <template #default>
            <v-row>
              <v-col cols="4">
                <v-text-field
                  v-model="search"
                  class="search-box mb-2"
                  label="Search"
                  prepend-inner-icon="mdi-magnify"
                  variant="outlined"
                  density="compact"
                  hide-details
                  single-line
                ></v-text-field>
              </v-col>
              <v-col cols="8">
                <v-btn
                  v-if="hasPermission('quote:access_new_business')"
                  class="mt-1"
                  color="primary"
                  size="small"
                  rounded
                  @click="goToQuoteCreation"
                  >New Quote</v-btn
                >
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-data-table
                  class="table-row"
                  density="compact"
                  :headers="headers"
                  :items="quotes"
                  :search="search"
                >
                  <!-- Slot for Status-->
                  <template #[`item.status`]="{ item }">
                    <v-chip
                      :color="getStatusColor(item.status)"
                      text-color="white"
                      small
                      class="ma-1"
                      variant="flat"
                    >
                      {{ snakeToTitleCase(item.status) }}
                    </v-chip>
                  </template>
                  <!-- Slot for Win Probability -->
                  <template #[`item.win_probability`]="{ item }">
                    <div
                      style="cursor: pointer"
                      @dblclick="openExplainer(item)"
                    >
                      <ProbabilityBadge
                        :score="winProbabilities[item.id]?.score_pct ?? null"
                        size="small"
                      />
                    </div>
                  </template>
                  <!-- Slot for Actions Column -->
                  <template #[`item.actions`]="{ item }">
                    <v-btn
                      icon
                      size="small"
                      variant="plain"
                      color="primary"
                      :disabled="item.status === 'accepted'"
                      @click="editItem(item)"
                    >
                      <v-icon>mdi-pencil</v-icon>
                    </v-btn>
                    <v-btn
                      icon
                      variant="plain"
                      size="small"
                      color="error"
                      :disabled="item.status === 'accepted'"
                      @click="deleteItem(item)"
                    >
                      <v-icon>mdi-delete</v-icon>
                    </v-btn>
                    <v-tooltip>
                      <template #activator="{ props }">
                        <v-btn
                          icon
                          color="primary"
                          variant="plain"
                          size="small"
                          v-bind="props"
                          @click="viewItem(item)"
                        >
                          <v-icon>mdi-eye</v-icon>
                        </v-btn>
                      </template>
                      <span>View Item</span>
                    </v-tooltip>
                    <v-tooltip>
                      <template #activator="{ props }">
                        <v-btn
                          icon
                          color="primary"
                          variant="plain"
                          size="small"
                          v-bind="props"
                          :disabled="item.status === 'accepted'"
                          @click="submitReview(item)"
                        >
                          <v-icon>mdi-file-eye-outline</v-icon>
                        </v-btn>
                      </template>
                      <span>Submit for Review</span>
                    </v-tooltip>
                    <v-tooltip
                      v-if="hasPermission('quote:access_new_business')"
                    >
                      <template #activator="{ props }">
                        <v-btn
                          icon
                          color="primary"
                          variant="plain"
                          size="small"
                          v-bind="props"
                          :disabled="item.status === 'accepted'"
                          @click="submitQuoteGeneration(item)"
                        >
                          <v-icon>mdi-file-send-outline</v-icon>
                        </v-btn>
                      </template>
                      <span>Generate Quote</span>
                    </v-tooltip>
                  </template>
                </v-data-table>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-col>
    </v-row>
    <v-dialog v-model="dialog" max-width="400">
      <base-card>
        <template #header>
          <span class="headline">Submit for Review</span>
        </template>
        <template #default>
          <!-- Select Field -->
          <v-select
            v-model="selectedReviewer"
            variant="outlined"
            density="compact"
            :items="reviewers"
            item-title="name"
            item-value="name"
            label="Select Reviewer"
            return-object
          ></v-select>
        </template>
        <template #actions>
          <!-- Cancel Button -->
          <v-btn variant="plain" @click="dialog = false">Cancel</v-btn>
          <!-- Ok Button -->
          <v-btn color="primary" @click="submitForReview">Ok</v-btn>
        </template>
      </base-card>
    </v-dialog>
    <confirm-dialog ref="confirmDeleteDialog"></confirm-dialog>

    <!-- Win Probability Explainer Dialog -->
    <v-dialog v-model="explainDialog" max-width="680" scrollable>
      <v-card>
        <v-card-title class="d-flex align-center ga-3 pt-4">
          <span>Win Probability — Algorithm Explainer</span>
          <v-spacer />
          <ProbabilityBadge
            :score="selectedExplainData?.score_pct ?? null"
            size="default"
          />
        </v-card-title>
        <v-card-subtitle class="pb-0">
          <span class="font-weight-medium">{{
            selectedExplainItem?.scheme_name
          }}</span>
          <span
            v-if="selectedExplainData?.band"
            class="ml-2 text-capitalize text-grey-darken-1"
            >· {{ selectedExplainData.band.replace(/_/g, ' ') }}</span
          >
          <span
            v-if="selectedExplainData?.scored_at"
            class="ml-2 text-grey-darken-1"
            >· Scored
            {{
              new Date(selectedExplainData.scored_at).toLocaleDateString()
            }}</span
          >
        </v-card-subtitle>
        <v-divider class="mt-3" />
        <v-card-text>
          <!-- How it works note -->
          <v-alert
            type="info"
            variant="tonal"
            density="compact"
            class="mb-4 text-body-2"
          >
            Each bar shows how much a feature
            <strong>pushes the score up</strong> (green) or
            <strong>down</strong> (red). Longer bars have more influence. The
            final score is the sum of all contributions passed through a sigmoid
            function.
          </v-alert>

          <!-- Feature waterfall bars -->
          <div v-if="explainFeatures.length">
            <div
              v-for="feat in explainFeatures"
              :key="feat.name"
              class="feat-row mb-3"
            >
              <div class="feat-label text-body-2">
                <span class="font-weight-medium">{{
                  featureLabel(feat.name)
                }}</span>
                <span class="text-grey-darken-1 ml-1 text-caption"
                  >({{ feat.contribution >= 0 ? '+' : ''
                  }}{{ feat.contribution.toFixed(3) }})</span
                >
              </div>
              <div class="feat-desc text-caption text-grey-darken-1 mb-1">
                {{ featureDescription(feat.name) }}
              </div>
              <div class="waterfall-track">
                <div
                  class="waterfall-bar"
                  :class="
                    feat.contribution >= 0 ? 'bar-positive' : 'bar-negative'
                  "
                  :style="waterfallStyle(feat.contribution)"
                />
              </div>
            </div>
          </div>
          <div v-else class="text-grey text-body-2 text-center py-4">
            No feature data available for this quote.
          </div>
        </v-card-text>
        <v-divider />
        <v-card-actions>
          <span class="text-caption text-grey-darken-1 ml-2"
            >Model weights learned from historical accepted/rejected quotes via
            logistic regression.</span
          >
          <v-spacer />
          <v-btn variant="text" @click="explainDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="3000">
      {{ snackbarMessage }}
      <template #actions>
        <v-btn color="white" variant="text" @click="hideNotification"
          >Close</v-btn
        >
      </template>
    </v-snackbar>
  </v-container>
</template>
<script setup lang="ts">
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useAppStore } from '@/renderer/store/app'
import { useRouter } from 'vue-router'
import { computed, onMounted, ref, watchEffect } from 'vue'
import { useGroupPricingStore } from '@/renderer/store/group_pricing'
import ConfirmDialog from '@/renderer/components/ConfirmDialog.vue'
import ProbabilityBadge from '@/renderer/components/ProbabilityBadge.vue'
import { useNotifications } from '@/renderer/composables/useNotifications'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

interface Quote {
  id: number
  scheme_name: string
  quote_type: string
  commencement_date: string
  basis: string
  status: string
  quote_broker: { name: string }
  obligation_type: string
  sgla_benefit: number
  phi_benefit: number
  ttd_benefit: number
  ptd_benefit: number
  ci_benefit: number
  family_funeral_benefit: number
  creation_date: string
  created_by: string
  reviewer: string
}

const confirmDeleteDialog = ref()
const {
  snackbar,
  snackbarMessage,
  snackbarColor,
  showSuccess,
  showError,
  hideNotification
} = useNotifications()
const { hasPermission } = usePermissionCheck()
const groupStore = useGroupPricingStore()
const router = useRouter()
const appStore = useAppStore()
const quotes = ref<Quote[]>([])
const search = ref('')
const selectedReviewer: any = ref(null)
const reviewers: any = ref([])
const dialog = ref(false)
const selectedQuote: any = ref({})
const organization: any = ref(null)
const benefitMaps: any = ref([])
const winProbabilities = ref<Record<number, any>>({})
const explainDialog = ref(false)
const selectedExplainItem = ref<any>(null)
const selectedExplainData = ref<any>(null)

const headers = ref([
  {
    title: 'Actions',
    value: 'actions',
    align: 'center' as 'center',
    sortable: false
  },
  {
    title: 'Scheme Name',
    value: 'scheme_name',
    key: 'scheme_name',
    width: '120px'
  },
  {
    title: 'Quote Name',
    value: 'quote_name',
    key: 'quote_name',
    width: '150px'
  },
  { title: 'Type', value: 'quote_type', key: 'quote_type', width: '20%' },
  {
    title: 'Commencement Date',
    key: 'commencement_date',
    width: '20%',
    value: (item: any) => parseDateString(item.commencement_date)
  },
  { title: 'Basis', value: 'basis' },
  { title: 'Status', key: 'status', value: 'status' },
  {
    title: 'Win %',
    key: 'win_probability',
    value: 'win_probability',
    sortable: false
  },
  {
    title: 'Channel',
    value: 'distribution_channel',
    key: 'distribution_channel'
  },
  {
    title: 'Broker',
    value: 'quote_broker.name',
    key: 'quote_broker',
    width: '20%'
  },
  { title: 'Type', value: 'obligation_type' },
  {
    title: 'Creation Date',
    key: 'creation_date',
    width: '20%',
    value: (item: any) => parseDateString(item.creation_date)
  },
  { title: 'Submitted By', value: 'created_by' },
  { title: 'Reviewer', value: 'reviewer' }
])

const snakeToTitleCase = (str: string) => {
  // Return null or an empty string if the input is falsy
  if (!str) {
    return ''
  }

  // Split the string by underscores, capitalize each word, and join with a space
  const result = str
    .split('_')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ')

  console.log('Converted string:', result) // Debugging log
  return result
}

const getStatusColor = (status) => {
  switch (status) {
    case 'pending_review':
      return '#FFE699'
    case 'in_force':
      return '#70AD47'
    case 'approved':
      return '#E2EFDA'
    case 'in_progress':
      return '#F8CBAD'
    case 'accepted':
      return '#70AD47'
    default:
      return 'black'
  }
}

const parseDateString = (dateString) => {
  // const date = new Date(dateString)
  // const formattedDate = date.toISOString().split('T')[0]
  return dateString.substring(0, 10)
  // return formattedDate
}

const goToQuoteCreation = () => {
  router.push({ name: 'group-pricing-quote-generation' })
}

const editItem = (item) => {
  router.push({
    name: 'group-pricing-quote-generation-edit',
    params: { id: item.id }
  })
}

const deleteItem = async (item) => {
  if (item.status === 'InForce') {
    showError('Cannot delete a quote that is InForce')
    return
  }
  try {
    const confirmed = await confirmDeleteDialog.value.open(
      'Delete Quote',
      'Are you sure you want to delete this quote?'
    )
    if (!confirmed) return

    await GroupPricingService.deleteQuote(item.id)
    quotes.value = quotes.value.filter((quote: any) => quote.id !== item.id)
    showSuccess('Quote deleted successfully')
  } catch (error: any) {
    const msg =
      error?.response?.data?.message ||
      error?.message ||
      'Failed to delete quote'
    showError(msg)
  }
}

const viewItem = (item) => {
  groupStore.selectedQuote = item
  router.push({ name: 'group-pricing-scheme-details', params: { id: item.id } })
}
const submitReview = (item) => {
  selectedQuote.value = item
  dialog.value = true
}

const submitForReview = () => {
  selectedQuote.value.reviewer = selectedReviewer.value.name
  selectedQuote.value.status = 'pending_review'
  GroupPricingService.changeQuoteStatus(selectedQuote.value)
  dialog.value = false
}

const submitQuoteGeneration = (item) => {
  router.push({
    name: 'group-pricing-quotes-generation',
    params: { quoteId: item.id }
  })
}

onMounted(() => {
  // Start fetching data that doesn't depend on organization.value immediately
  GroupPricingService.getBenefitMaps()
    .then((res) => {
      benefitMaps.value = res.data
      headers.value = headers.value.map((header) => {
        const bff = benefitMaps.value.find(
          (map) => map.benefit_code === header.title
        )
        if (bff && bff.benefit_alias !== '') {
          return {
            ...header,
            title: bff.benefit_alias
          }
        }
        return header
      })
    })
    .catch((error) => {
      console.error('Error fetching benefit maps:', error)
      // Handle error appropriately, e.g., set default headers or show a message
    })

  GroupPricingService.getAllQuotes()
    .then((res) => {
      if (res.data && res.data.length > 0) {
        quotes.value = res.data
        console.log('Quotes:', quotes.value)
        // Fetch win probabilities in parallel (fire-and-forget per quote)
        const ids: number[] = quotes.value.map((q: any) => q.id)
        Promise.allSettled(
          ids.map((id) =>
            GroupPricingService.getQuoteWinProbability(id)
              .then((r) => {
                if (r.data?.data) {
                  winProbabilities.value = {
                    ...winProbabilities.value,
                    [id]: r.data.data
                  }
                }
              })
              .catch(() => {})
          )
        )
      } else {
        quotes.value = []
      }
    })
    .catch((error) => {
      console.error('Error fetching all quotes:', error)
      quotes.value = [] // Ensure quotes is an array on error
    })

  // Re-fetch reviewers whenever the resolved organisation name becomes
  // available or changes (e.g. after license activation completes).
  watchEffect(() => {
    const orgName = appStore.getOrganisationName
    if (!orgName || organization.value === orgName) return
    organization.value = orgName

    GroupPricingService.getOrgUsers({ name: orgName })
      .then((res) => {
        if (res && Array.isArray(res.data)) {
          const uniqueData = Array.from(
            new Map(res.data.map((entry) => [entry.name, entry])).values()
          )
          reviewers.value = uniqueData
        } else {
          console.warn(
            'Org Users response is not as expected or data is missing:',
            res
          )
          reviewers.value = []
        }
      })
      .catch((error) => {
        console.error('Error fetching org users:', error)
        reviewers.value = []
      })
  })
})

// --- Win Probability Explainer ---

const FEATURE_META: Record<string, { label: string; description: string }> = {
  is_renewal: {
    label: 'Renewal Quote',
    description:
      'Renewal quotes historically convert at higher rates than new business'
  },
  distribution_channel: {
    label: 'Distribution Channel',
    description:
      'How the quote reaches the client — broker, direct, binder, or tied agent'
  },
  discount_pct: {
    label: 'Premium Discount',
    description:
      'Overall discount applied to the gross premium — higher discounts can improve win rate'
  },
  commission_pct: {
    label: 'Commission Rate',
    description:
      'Broker commission rate; very high rates may signal aggressive deals'
  },
  total_loading_pct: {
    label: 'Total Loading',
    description:
      'Sum of all expense, profit, and contingency loadings minus discount — proxy for competitiveness'
  },
  member_count: {
    label: 'Member Count',
    description:
      'Number of members in the scheme (log-scaled); larger schemes tend to be more competitive'
  },
  avg_age: {
    label: 'Average Age',
    description:
      'Average age of scheme members — older cohorts carry higher mortality risk'
  },
  avg_income: {
    label: 'Average Income',
    description:
      'Average member income (log-scaled) — higher-income schemes attract more competition'
  },
  gender_ratio: {
    label: 'Gender Ratio (Male)',
    description:
      'Proportion of male members (0–1); influences mortality and disability risk pricing'
  },
  has_experience_data: {
    label: 'Has Claims Experience',
    description:
      'Whether historical claims data exists — experience-rated quotes are more defensible'
  },
  expected_loss_ratio: {
    label: 'Expected Loss Ratio',
    description:
      'Expected claims ÷ premium — values closer to 1.0 indicate tighter, more competitive pricing'
  },
  broker_historic_rate: {
    label: 'Broker Win Rate',
    description:
      "This broker's historical conversion rate across all past quotes"
  },
  days_to_commencement: {
    label: 'Days to Commencement',
    description:
      'Lead time from quote creation to policy start — longer lead times give more time to negotiate'
  },
  benefit_count: {
    label: 'Benefit Count',
    description:
      'Number of benefit types (GLA, PHI, CI, etc.) included — comprehensive schemes tend to convert better'
  }
}

const featureLabel = (name: string) =>
  FEATURE_META[name]?.label ?? name.replace(/_/g, ' ')
const featureDescription = (name: string) =>
  FEATURE_META[name]?.description ?? ''

const explainFeatures = computed<
  { name: string; contribution: number; weight: number }[]
>(() => {
  if (!selectedExplainData.value?.top_features) return []
  try {
    const parsed = JSON.parse(selectedExplainData.value.top_features)
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
})

// The maximum absolute contribution — used to scale bar widths to 100%
const maxAbsContrib = computed(() => {
  const vals = explainFeatures.value.map((f) => Math.abs(f.contribution))
  return vals.length ? Math.max(...vals) : 1
})

const waterfallStyle = (contribution: number) => {
  const pct =
    maxAbsContrib.value > 0
      ? (Math.abs(contribution) / maxAbsContrib.value) * 100
      : 0
  return { width: `${pct}%` }
}

const openExplainer = (item: any) => {
  selectedExplainItem.value = item
  selectedExplainData.value = winProbabilities.value[item.id] ?? null
  explainDialog.value = true
}
</script>
<style lang="css" scoped>
.table-row {
  white-space: nowrap;
}

::v-deep(.v-data-table thead th) {
  background-color: #223f54 !important;
  color: white;
  text-align: center;
  font-weight: bold;
  white-space: nowrap;
  min-width: 150px;
}

.search-box {
  width: 100%;
}
.v-table__wrapper > table > thead {
  background-color: #223f54 !important;
  color: white;
  white-space: nowrap;
}

/* Win probability waterfall bars */
.feat-row {
  padding-bottom: 2px;
}
.feat-label {
  display: flex;
  align-items: baseline;
}
.waterfall-track {
  background: #f0f0f0;
  border-radius: 4px;
  height: 10px;
  overflow: hidden;
  width: 100%;
}
.waterfall-bar {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease;
}
.bar-positive {
  background: #4caf50;
}
.bar-negative {
  background: #f44336;
}
</style>
