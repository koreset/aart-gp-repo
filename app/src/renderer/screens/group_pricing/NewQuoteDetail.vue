<template>
  <v-container v-if="quote">
    <!-- Calculation progress overlay -->
    <v-overlay
      :model-value="awaitingCalculation"
      persistent
      class="align-start justify-center pt-16"
      scrim="rgba(0,0,0,0.4)"
    >
      <v-card width="400" class="pa-6 text-center" rounded="lg" elevation="8">
        <v-card-title class="text-h6 mb-2">
          {{
            calcProgress?.phase === 'queued'
              ? 'Calculation Queued'
              : calcProgress?.phase === 'failed'
                ? 'Calculation Failed'
                : 'Calculating Quote'
          }}
        </v-card-title>
        <v-card-text>
          <v-progress-linear
            v-if="calcProgress?.phase !== 'queued'"
            :model-value="progressPercent"
            color="primary"
            height="8"
            rounded
            class="mb-3"
          />
          <v-progress-linear
            v-else
            indeterminate
            color="primary"
            height="8"
            rounded
            class="mb-3"
          />
          <div
            v-if="calcProgress?.phase !== 'queued'"
            class="text-body-1 font-weight-medium mb-1"
            >{{ progressPercent }}%</div
          >
          <div class="text-body-2 text-medium-emphasis">
            {{ phaseLabel }}
            <span v-if="calcProgress?.currentCategory">
              — {{ calcProgress.currentCategory }}
            </span>
          </div>
          <div
            v-if="calcProgress && calcProgress.totalCategories > 0"
            class="text-caption text-medium-emphasis mt-1"
          >
            {{ calcProgress.completedCategories }} /
            {{ calcProgress.totalCategories }} categories
          </div>
        </v-card-text>
        <v-card-actions class="justify-center">
          <v-btn variant="text" size="small" @click="dismissCalcOverlay">
            Dismiss
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-overlay>

    <base-card :show-actions="false">
      <template #header>
        <h4 class="text-h6 d-flex align-center">
          {{ quote.scheme_name }}
          <v-chip :color="statusColor(quote.status)" class="ml-4" label>
            {{ quote.status }}
          </v-chip>
        </h4>
      </template>
      <template #default>
        <v-row class="mb-4">
          <v-col cols="12" md="2">
            <v-btn
              class="mr-2"
              size="small"
              variant="text"
              @click="router.go(-1)"
              >Back</v-btn
            >
          </v-col>
          <v-col cols="12" md="10" class="d-flex text-md-right justify-end">
            <v-tooltip
              location="top"
              text="Edit the quote's initial parameters."
            >
              <template #activator="{ props: tooltipProps }">
                <div v-bind="tooltipProps" class="d-inline-block">
                  <v-btn
                    :disabled="
                      quote.status === 'accepted' || quote.status === 'in_force'
                    "
                    class="mr-2"
                    size="small"
                    rounded
                    variant="outlined"
                    color="primary"
                    @click="editQuote(quote.id)"
                    >Edit</v-btn
                  >
                </div>
              </template>
            </v-tooltip>

            <v-tooltip
              location="top"
              :text="
                customTirMissing
                  ? 'Custom tiered income replacement table is missing. The administrator must upload the table before calculations can be run.'
                  : quote.experience_rating === 'Override' &&
                      (quote.experience_rate_overrides_count || 0) === 0
                    ? 'Run with no overrides to see the baseline loaded rates per benefit, then add overrides and re-run.'
                    : 'Requires member data and (if Yes) a claims-experience upload.'
              "
            >
              <template #activator="{ props: tooltipProps }">
                <div v-bind="tooltipProps" class="d-inline-block">
                  <v-btn
                    :disabled="
                      customTirMissing ||
                      quote.member_data_count === 0 ||
                      (quote.experience_rating === 'Yes' &&
                        quote.claims_experience_count === 0) ||
                      quote.status === 'in_force' ||
                      quote.status === 'accepted'
                    "
                    rounded
                    class="mr-2"
                    size="small"
                    color="primary"
                    :loading="loading"
                    @click="basisDialog = true"
                    >Run Calculations</v-btn
                  >
                  <div
                    v-if="dismissedCalcStatusLabel"
                    class="text-caption text-medium-emphasis mt-1 text-center mr-2"
                  >
                    {{ dismissedCalcStatusLabel }}
                  </div>
                </div>
              </template>
            </v-tooltip>

            <v-tooltip
              location="top"
              text="Apply a discount to the calculated premiums. Not available once the quote has been approved."
            >
              <template #activator="{ props: tooltipProps }">
                <div v-bind="tooltipProps" class="d-inline-block">
                  <v-btn
                    :disabled="
                      quote.member_rating_result_count === 0 ||
                      quote.status === 'approved' ||
                      quote.status === 'in_force' ||
                      quote.status === 'accepted'
                    "
                    rounded
                    class="mr-2"
                    size="small"
                    color="secondary"
                    variant="outlined"
                    @click="openDiscountDialog"
                    >Apply Discount</v-btn
                  >
                </div>
              </template>
            </v-tooltip>

            <v-tooltip
              location="top"
              text="All data must be populated and results generated."
            >
              <template #activator="{ props: tooltipProps }">
                <div v-bind="tooltipProps" class="d-inline-block">
                  <v-btn
                    v-if="hasPermission('quote:approve')"
                    class="mr-2"
                    size="small"
                    color="primary"
                    rounded
                    :disabled="
                      quote.status === 'in_force' ||
                      hasEmptyQuoteTables ||
                      quote.status === 'accepted'
                    "
                    @click="approveQuote"
                    >Approve</v-btn
                  >
                </div>
              </template>
            </v-tooltip>

            <v-tooltip
              location="top"
              text="Quote must be in 'Approved' status."
            >
              <template #activator="{ props: tooltipProps }">
                <div v-bind="tooltipProps" class="d-inline-block">
                  <v-btn
                    v-if="hasPermission('quote:accept')"
                    class="mr-2"
                    :disabled="
                      quote.status !== 'approved' || quote.status === 'InForce'
                    "
                    size="small"
                    rounded
                    color="primary"
                    :loading="acceptQuoteLoading"
                    @click="acceptQuote"
                    >Accept Quote</v-btn
                  >
                </div>
              </template>
            </v-tooltip>
            <v-menu location="bottom end">
              <template #activator="{ props: menuProps }">
                <v-btn
                  v-bind="menuProps"
                  size="small"
                  rounded
                  color="primary"
                  variant="outlined"
                  append-icon="mdi-menu-down"
                  :loading="
                    isGeneratingOnRiskLetter || isGeneratingTemplatedOnRisk
                  "
                >
                  Export
                </v-btn>
              </template>
              <v-list density="compact">
                <v-list-item
                  :disabled="
                    (quote.status !== 'accepted' &&
                      quote.status !== 'in_force' &&
                      quote.status !== 'Accepted' &&
                      quote.status !== 'InForce') ||
                    isGeneratingOnRiskLetter
                  "
                  prepend-icon="mdi-file-document-outline"
                  title="On Risk Letter"
                  subtitle="Available after quote acceptance."
                  @click="generateOnRiskLetterManual"
                />
                <v-list-item
                  :disabled="
                    (quote.status !== 'accepted' &&
                      quote.status !== 'in_force' &&
                      quote.status !== 'Accepted' &&
                      quote.status !== 'InForce') ||
                    isGeneratingTemplatedOnRisk
                  "
                  prepend-icon="mdi-file-document-edit-outline"
                  title="On Risk Letter (from Template)"
                  subtitle="Uses the insurer's uploaded Word template."
                  @click="generateOnRiskLetterFromBackend"
                />
                <v-list-item
                  :disabled="hasEmptyQuoteTables"
                  prepend-icon="mdi-file-pdf-box"
                  title="Quote PDF"
                  subtitle="Quote must be in 'Accepted' status."
                  @click="generatePdf(quote.id)"
                />
              </v-list>
            </v-menu>
          </v-col>
        </v-row>
        <v-tabs v-model="tab" color="primary" class="mb-5">
          <v-tab value="summary">Quote Summary</v-tab>
          <v-tab value="benefits">Benefits & Config</v-tab>
          <v-tab value="data">Data Management</v-tab>
          <v-tab v-if="hasPermission('quote:view_results')" value="results"
            >Results & Analysis</v-tab
          >
          <v-tab v-if="resultSummaries !== null" value="outputsummary"
            >Output Summary</v-tab
          >
          <v-tab value="benefitssummary">Premiums Summary</v-tab>
          <v-tab value="reinsurancepremiumsummary"
            >Reinsurance Premium Summary</v-tab
          >
          <v-tab value="additionalglacover">Additional GLA Cover</v-tab>
        </v-tabs>
        <v-window v-model="tab">
          <v-window-item value="summary" eager>
            <QuoteSummaryData
              :quote="quote"
              :result-summaries="resultSummaries || []"
            />
          </v-window-item>
          <v-window-item value="benefits" eager>
            <QuoteBenefitsConfiguration
              :quote="quote"
              :result-summaries="resultSummaries || []"
            />
          </v-window-item>
          <v-window-item value="data" eager>
            <QuoteDataTableManager
              :quote="quote"
              :result-summaries="resultSummaries || []"
              @quote-updated="loadQuote"
              @indicative-data-updated="handleIndicativeDataUpdate"
              @navigate-to-general-input="editQuote(quote.id)"
            />
          </v-window-item>
          <v-window-item
            v-if="hasPermission('quote:view_results')"
            value="results"
            eager
          >
            <QuoteResults :quote="quote" @quote-updated="loadQuote" />
          </v-window-item>
          <v-window-item value="outputsummary" eager>
            <OutputSummary
              v-if="resultSummaries !== null && resultSummaries.length > 0"
              :key="`output-${resultsRefreshKey}`"
              :quote="quote"
              :resultSummaries="resultSummaries"
            />
          </v-window-item>
          <v-window-item value="benefitssummary" eager>
            <QuoteBenefitSummary
              v-if="resultSummaries !== null"
              :key="`premium-${resultsRefreshKey}`"
              :resultSummaries="resultSummaries"
              :quote="quote"
            />
          </v-window-item>
          <v-window-item value="reinsurancepremiumsummary" eager>
            <QuoteReinsurancePremiumSummary
              v-if="resultSummaries !== null"
              :key="`reins-${resultsRefreshKey}`"
              :resultSummaries="resultSummaries"
              :quote="quote"
            />
          </v-window-item>
          <v-window-item value="additionalglacover" eager>
            <AdditionalGlaCoverSummary
              v-if="resultSummaries !== null"
              :key="`agla-${resultsRefreshKey}`"
              :resultSummaries="resultSummaries"
              :quote="quote"
              @quote-updated="loadQuote"
            />
          </v-window-item>
        </v-window>
      </template>
    </base-card>

    <v-dialog v-model="basisDialog" persistent max-width="550px">
      <base-card>
        <template #header>
          <span class="headline">Choose a parameter basis</span>
        </template>
        <template #default>
          <v-row>
            <v-col>
              <v-select
                v-model:model-value="quote.basis"
                clearable
                variant="outlined"
                density="compact"
                placeholder="Select a Basis"
                label="Basis"
                :items="parameterBases"
                item-title="basis"
                item-children="basis"
              ></v-select>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <p class="text-body-2 text-medium-emphasis mb-2">
                Binder & outsource fees
                <span v-if="!isBinderChannel">
                  &mdash; only editable when the quote's distribution channel is
                  <strong>Binder</strong>.
                </span>
                <span v-else-if="binderFeeCapMissing" class="text-warning">
                  &mdash; no binder fee cap configured for this broker and risk
                  rate code. Ask the administrator to add one before running
                  calculations.
                </span>
              </p>
            </v-col>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model.number="quote.loadings.binder_fee"
                type="number"
                variant="outlined"
                density="compact"
                label="Binder fee (%)"
                suffix="%"
                :min="0"
                :max="maxBinderFee ?? undefined"
                :disabled="!isBinderChannel"
                :hint="
                  isBinderChannel && maxBinderFee !== null
                    ? `Max: ${maxBinderFee}%`
                    : ''
                "
                persistent-hint
                :rules="[
                  (v) => v >= 0 || 'Cannot be negative',
                  (v) =>
                    maxBinderFee === null ||
                    v <= maxBinderFee ||
                    `Max allowed is ${maxBinderFee}%`
                ]"
              ></v-text-field>
            </v-col>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model.number="quote.loadings.outsource_fee"
                type="number"
                variant="outlined"
                density="compact"
                label="Outsource fee (%)"
                suffix="%"
                :min="0"
                :max="maxOutsourceFee ?? undefined"
                :disabled="!isBinderChannel"
                :hint="
                  isBinderChannel && maxOutsourceFee !== null
                    ? `Max: ${maxOutsourceFee}%`
                    : ''
                "
                persistent-hint
                :rules="[
                  (v) => v >= 0 || 'Cannot be negative',
                  (v) =>
                    maxOutsourceFee === null ||
                    v <= maxOutsourceFee ||
                    `Max allowed is ${maxOutsourceFee}%`
                ]"
              ></v-text-field>
            </v-col>
          </v-row>
        </template>
        <template #actions>
          <v-spacer></v-spacer>
          <v-btn
            rounded
            variant="text"
            :disabled="binderFeeInvalid"
            @click="closeBasisDialog(true)"
            >Ok</v-btn
          >
          <v-btn rounded variant="text" @click="closeBasisDialog(false)"
            >Cancel</v-btn
          >
        </template>
      </base-card>
    </v-dialog>

    <!-- Apply Discount Dialog -->
    <v-dialog v-model="discountDialog" persistent max-width="450px">
      <base-card>
        <template #header>
          <span class="headline">Apply Discount</span>
        </template>
        <template #default>
          <v-row>
            <v-col cols="12">
              <p class="text-body-2 text-medium-emphasis mb-3">
                Enter a discount percentage to apply to the calculated premiums.
                Maximum allowed: <strong>{{ maxDiscount }}%</strong>
              </p>
              <v-text-field
                v-model.number="discountRate"
                type="number"
                variant="outlined"
                density="compact"
                label="Discount (%)"
                :min="0"
                :max="maxDiscount"
                :rules="[
                  (v) => v >= 0 || 'Discount cannot be negative',
                  (v) =>
                    v <= maxDiscount ||
                    `Maximum allowed discount is ${maxDiscount}%`
                ]"
                suffix="%"
              ></v-text-field>
            </v-col>
          </v-row>
        </template>
        <template #actions>
          <v-spacer></v-spacer>
          <v-btn
            rounded
            variant="text"
            color="primary"
            :loading="discountLoading"
            :disabled="
              discountLoading || discountRate < 0 || discountRate > maxDiscount
            "
            @click="confirmApplyDiscount"
            >Apply</v-btn
          >
          <v-btn
            rounded
            variant="text"
            :disabled="discountLoading"
            @click="discountDialog = false"
            >Cancel</v-btn
          >
        </template>
      </base-card>
    </v-dialog>

    <confirm-dialog ref="confirmAction" />
    <v-snackbar
      v-model="snackbar"
      centered
      :timeout="snackbarTimeout"
      :multi-line="true"
    >
      {{ snackbarText }}
      <v-btn rounded color="red" variant="text" @click="snackbar = false"
        >Close</v-btn
      >
    </v-snackbar>
  </v-container>
  <v-container
    v-else
    class="d-flex justify-center align-center"
    style="height: 80vh"
  >
    <v-progress-circular
      indeterminate
      color="primary"
      size="64"
    ></v-progress-circular>
  </v-container>
  <!-- Accept Quote Dialog -->
  <v-dialog v-model="acceptQuoteDialog" persistent max-width="550px">
    <base-card>
      <template #header>
        <span class="headline">Accept Quote</span>
      </template>
      <template #default>
        <v-row>
          <v-col cols="12">
            <v-date-input
              v-model="acceptQuoteParams.commencementDate"
              hide-actions
              locale="en-ZA"
              view-mode="month"
              prepend-icon=""
              prepend-inner-icon="$calendar"
              variant="outlined"
              density="compact"
              label="Commencement Date"
              placeholder="Select a date"
            ></v-date-input>
          </v-col>
          <v-col cols="12">
            <v-text-field
              v-model="acceptQuoteParams.term"
              type="number"
              variant="outlined"
              density="compact"
              label="Term (months)"
              placeholder="Enter term in months"
            ></v-text-field>
          </v-col>
        </v-row>
      </template>
      <template #actions>
        <v-spacer></v-spacer>
        <v-btn
          rounded
          variant="text"
          :loading="acceptQuoteLoading"
          :disabled="acceptQuoteLoading"
          @click="confirmAcceptQuote"
          >Accept</v-btn
        >
        <v-btn
          rounded
          variant="text"
          :disabled="acceptQuoteLoading"
          @click="acceptQuoteDialog = false"
          >Cancel</v-btn
        >
      </template>
    </base-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import BaseCard from '@/renderer/components/BaseCard.vue'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

// Import Child Components
import QuoteSummaryData from './QuoteSummaryData.vue'
import QuoteBenefitsConfiguration from './QuoteBenefitsConfiguration.vue'
import QuoteDataTableManager from './QuoteDataTableManager.vue'
import QuoteResults from './QuoteResults.vue'
import QuoteBenefitSummary from './QuoteBenefitSummary.vue'
import QuoteReinsurancePremiumSummary from './QuoteReinsurancePremiumSummary.vue'
import AdditionalGlaCoverSummary from './AdditionalGlaCoverSummary.vue'
import formatDateString from '@/renderer/utils/helpers'

// Import Other Components
import ConfirmDialog from '@/renderer/components/ConfirmDialog.vue'
import OutputSummary from './OutputSummary.vue'
import { VDateInput } from 'vuetify/labs/VDateInput'
import { useOnRiskLetterGeneration } from '@/renderer/composables/useOnRiskLetterGeneration'
import { useCalculationProgress } from '@/renderer/composables/useCalculationProgress'

const { hasPermission } = usePermissionCheck()

const {
  isGenerating: isGeneratingOnRiskLetter,
  generateOnRiskLetterDocx,
  generateOnRiskLetterPdf
} = useOnRiskLetterGeneration()

const {
  progress: calcProgress,
  phaseLabel,
  progressPercent,
  startTracking,
  stopTracking
} = useCalculationProgress()

// Only show overlay and react to events when this component initiated the calculation.
const awaitingCalculation = ref(false)
// True from the moment the user clicks Run Calculations until the calculation
// terminates (completed/failed). Stays true after the user dismisses the
// overlay so the success snackbar still fires and the inline status label
// keeps reporting progress under the Run Calculations button.
const isOurCalc = ref(false)

const dismissedCalcStatusLabel = computed(() => {
  if (!isOurCalc.value || awaitingCalculation.value) return ''
  const phase = calcProgress.value?.phase
  if (!phase || phase === 'completed' || phase === 'failed') return ''
  if (phase === 'queued') {
    const pos = calcProgress.value?.queuePosition
    return pos ? `Queued (${pos})` : 'Queued'
  }
  return 'In progress'
})

// Watchdog: if the WebSocket connection drops mid-calculation (e.g. a long
// debugger pause or a real-world network blip exceeds the server's pongWait)
// the final `completed` event can be lost, leaving the overlay stuck. The
// timer below polls the result summary after 30s of silence and dismisses
// the overlay if results have landed. The Dismiss button is the universal
// manual escape hatch.
const SILENCE_THRESHOLD_MS = 30_000
const WATCHDOG_TICK_MS = 5_000
let lastProgressEventAt = 0
let watchdogTimer: ReturnType<typeof setInterval> | null = null

function startCalcWatchdog() {
  lastProgressEventAt = Date.now()
  if (watchdogTimer) clearInterval(watchdogTimer)
  watchdogTimer = setInterval(async () => {
    if (!isOurCalc.value) {
      stopCalcWatchdog()
      return
    }
    if (Date.now() - lastProgressEventAt < SILENCE_THRESHOLD_MS) return
    // Only auto-recover when progress events indicated all categories were
    // already finished — i.e. only the final `completed` event was missed.
    // Earlier-phase silences fall through to the manual Dismiss button to
    // avoid declaring success on a calculation that may still be running.
    const p = calcProgress.value
    const lookedDone =
      !!p && p.totalCategories > 0 && p.completedCategories >= p.totalCategories
    if (!lookedDone) {
      lastProgressEventAt = Date.now()
      return
    }
    try {
      const resp = await GroupPricingService.getResultSummary(props.id)
      const hasResults = resp?.status === 200 && resp.data != null
      if (hasResults) {
        snackbarText.value = 'Calculations Successful'
        snackbar.value = true
        awaitingCalculation.value = false
        isOurCalc.value = false
        stopTracking()
        await loadQuote()
        stopCalcWatchdog()
      } else {
        lastProgressEventAt = Date.now()
      }
    } catch {
      lastProgressEventAt = Date.now()
    }
  }, WATCHDOG_TICK_MS)
}

function stopCalcWatchdog() {
  if (watchdogTimer) {
    clearInterval(watchdogTimer)
    watchdogTimer = null
  }
}

// Dismiss only hides the overlay. Tracking and the watchdog keep running so
// the success snackbar still fires on completion and the inline status label
// under Run Calculations stays current.
function dismissCalcOverlay() {
  awaitingCalculation.value = false
}

watch(calcProgress, (val) => {
  if (!isOurCalc.value) return
  lastProgressEventAt = Date.now()
  if (val?.phase === 'completed') {
    snackbarText.value = 'Calculations Successful'
    snackbar.value = true
    awaitingCalculation.value = false
    isOurCalc.value = false
    stopCalcWatchdog()
    loadQuote()
  }
  if (val?.phase === 'failed') {
    snackbarText.value =
      'Calculations failed. Please contact your administrator.'
    snackbar.value = true
    awaitingCalculation.value = false
    isOurCalc.value = false
    stopCalcWatchdog()
  }
})

onUnmounted(() => {
  stopCalcWatchdog()
})

const props = defineProps({
  id: {
    type: String,
    required: true
  }
})

const router = useRouter()
const quote: any = ref(null)
const loading = ref(false)
const acceptQuoteLoading = ref(false)
const tab = ref('summary')
const confirmAction = ref()
const broker = ref(null)
const parameterBases = ref([]) // Array to hold basis options
const resultSummaries: any = ref(null)
// Bumped on every successful refetch of the result summary so the bound
// child components (Output Summary, Premiums Summary, Reinsurance Premium
// Summary) are forcibly remounted. Belt-and-braces against any deep-watch
// or prop-diffing edge case that might leave a child rendering against
// stale data after a recalculation.
const resultsRefreshKey = ref(0)
const acceptQuoteParams: any = ref({
  commencementDate: '',
  term: 0
})

// State for dialogs and snackbar
const basisDialog = ref(false)
const acceptQuoteDialog = ref(false)
const discountDialog = ref(false)
const customTirMissing = ref(false)
const discountRate = ref(0)
const maxDiscount = ref(0)
const discountLoading = ref(false)
const snackbar = ref(false)
const snackbarText = ref('')
const snackbarTimeout = ref(4000) // 4 seconds for confirmations

// Binder / outsource fee input state — surfaced on the basis dialog only when
// the quote's distribution channel is "binder". Caps come from the binder_fees
// table keyed by (binderholder_name = broker name, risk_rate_code) and bound
// the values the binderholder can iterate on to hit a competitive rate.
const maxBinderFee = ref<number | null>(null)
const maxOutsourceFee = ref<number | null>(null)
const binderFeeCapMissing = ref(false)

const isBinderChannel = computed(
  () => quote.value?.distribution_channel === 'binder'
)

const loadBinderFeeCaps = async () => {
  maxBinderFee.value = null
  maxOutsourceFee.value = null
  binderFeeCapMissing.value = false
  if (!isBinderChannel.value) return
  const brokerName = quote.value?.quote_broker?.name?.trim()
  const riskRateCode = quote.value?.risk_rate_code?.trim()
  if (!brokerName || !riskRateCode) {
    binderFeeCapMissing.value = true
    return
  }
  try {
    const res = await GroupPricingService.getBinderFees()
    const rows: any[] = res.data || []
    const match = rows.find(
      (r) =>
        r.binderholder_name?.trim() === brokerName &&
        r.risk_rate_code?.trim() === riskRateCode
    )
    if (match) {
      // binder_fees caps are persisted as decimals (0.075 means 7.5%).
      // The basis dialog works in whole-percent units, so convert on read.
      maxBinderFee.value =
        Math.round(Number(match.maximum_binder_fee) * 100 * 10000) / 10000 || 0
      maxOutsourceFee.value =
        Math.round(Number(match.maximum_outsource_fee) * 100 * 10000) / 10000 ||
        0
    } else {
      binderFeeCapMissing.value = true
    }
  } catch {
    binderFeeCapMissing.value = true
  }
}

watch(basisDialog, (open) => {
  if (open) loadBinderFeeCaps()
})

// Computed properties (can be simplified if child components manage their own data)
const hasEmptyQuoteTables = computed(() => {
  if (!quote.value) return true
  return (
    quote.value.member_data_count === 0 ||
    (quote.value.experience_rating === 'Yes' &&
      quote.value.claims_experience_count === 0) ||
    quote.value.member_rating_result_count === 0
  )
})

const statusColor = (status: string) => {
  switch (status) {
    case 'InForce':
      return 'success'
    case 'Approved':
      return 'info'
    case 'Draft':
      return 'warning'
    default:
      return 'grey'
  }
}

const binderFeeInvalid = computed(() => {
  if (!isBinderChannel.value) return false
  const bf = Number(quote.value?.loadings?.binder_fee ?? 0)
  const of = Number(quote.value?.loadings?.outsource_fee ?? 0)
  if (bf < 0 || of < 0) return true
  if (maxBinderFee.value !== null && bf > maxBinderFee.value) return true
  if (maxOutsourceFee.value !== null && of > maxOutsourceFee.value) return true
  return false
})

const closeBasisDialog = async (value) => {
  if (value) {
    if (binderFeeInvalid.value) {
      snackbarText.value =
        'Binder fee or outsource fee exceeds the maximum allowed for this binderholder.'
      snackbar.value = true
      return
    }
    // Persist the quote (including binder/outsource fees on loadings) so the
    // enqueued calculation reads the latest values from the DB.
    if (isBinderChannel.value) {
      try {
        await GroupPricingService.changeQuoteStatus(quote.value)
      } catch (error) {
        console.error(
          'Failed to save binder/outsource fees before calc:',
          error
        )
        snackbarText.value = 'Could not save the binder fees. Please try again.'
        snackbar.value = true
        return
      }
    }
    runQuoteCalculations()
  }
  basisDialog.value = false
}

const loadQuote = async () => {
  try {
    const res = await GroupPricingService.getQuote(props.id)
    quote.value = res.data
    broker.value = quote.value.quoteBroker
    // Ensure loadings object exists with the new binder/outsource fields so
    // v-model inputs don't trip on undefined when rendering the dialog.
    if (!quote.value.loadings) quote.value.loadings = {}
    if (quote.value.loadings.binder_fee == null) {
      quote.value.loadings.binder_fee = 0
    }
    if (quote.value.loadings.outsource_fee == null) {
      quote.value.loadings.outsource_fee = 0
    }

    const res1 = await GroupPricingService.getQuoteTable(
      quote.value.id,
      'group_pricing_parameters'
    )

    if (res1.data !== null && res1.data.data.length > 0) {
      parameterBases.value = res1.data.data.map((item: any) => {
        if (
          item.risk_rate_code !== '' &&
          item.risk_rate_code === quote.value.risk_rate_code
        ) {
          return {
            basis: item.basis
          }
        }
      })
    }

    // strip out empty entries
    parameterBases.value = parameterBases.value.filter(
      (item: any) => item !== undefined && item !== null
    )

    const resp = await GroupPricingService.getResultSummary(props.id)
    if (resp.status === 200) {
      // Force a clean remount of the summary children even if the new
      // payload happens to deep-equal the old one (or if Vue reactivity
      // misses the swap for any reason). The :key on each summary panel
      // is bound to resultsRefreshKey, so bumping it guarantees a fresh
      // render against the new data.
      resultSummaries.value = resp.data
      resultsRefreshKey.value++
    } else {
      resultSummaries.value = null
    }

    // Pre-flight: check whether the custom TIR table has been uploaded
    try {
      const tirCheck = await GroupPricingService.getCustomTirStatus(props.id)
      const tirData = tirCheck.data?.data
      customTirMissing.value = !!(
        tirData?.needs_custom_tir && !tirData?.has_table
      )
    } catch {
      customTirMissing.value = false
    }
  } catch (error) {
    console.log('Error:', error)
  }
}

const runQuoteCalculations = async () => {
  if (quote.value.basis !== null && quote.value.basis !== '') {
    loading.value = true
    awaitingCalculation.value = true
    isOurCalc.value = true
    startTracking(String(quote.value.id))
    startCalcWatchdog()
    try {
      const res = await GroupPricingService.runQuoteCalculations(
        quote.value.id,
        quote.value.basis
      )
      if (res.status === 202 || res.status === 201) {
        // Job has been queued — progress updates arrive via WebSocket.
        // The button stays in loading state; the overlay shows queue/progress.
        // We stop the button spinner since the overlay takes over.
        loading.value = false
      }
    } catch (error: any) {
      console.error('Error:', error)
      stopTracking()
      stopCalcWatchdog()
      awaitingCalculation.value = false
      isOurCalc.value = false
      const apiMessage = error?.response?.data?.message || error?.response?.data
      snackbarText.value =
        typeof apiMessage === 'string' && apiMessage
          ? apiMessage
          : 'Calculations failed. Please contact your administrator.'
      snackbar.value = true
      loading.value = false
    }
  } else {
    snackbarText.value = 'Please select a basis'
    snackbar.value = true
  }
}

const approveQuote = async () => {
  try {
    // Show confirmation dialog
    const result = await confirmAction.value.open(
      'Approve Quote',
      'Are you sure you want to approve this quote?'
    )

    // If user cancels the confirmation, return without approving
    if (!result) {
      return
    }

    await GroupPricingService.approveQuote(quote.value.id)
    snackbarText.value = 'Quote has been approved successfully.'
    snackbarTimeout.value = 4000
    snackbar.value = true

    // Refresh the quote data to update the view
    await loadQuote()
  } catch (error) {
    console.error('Error:', error)
    snackbarText.value = 'Quote Approval Failed'
    snackbarTimeout.value = 2000
    snackbar.value = true
  }
}

const acceptQuote = async () => {
  // Check if the quote is already InForce
  if (quote.value.status === 'in_force') {
    snackbarText.value = 'Quote is already InForce'
    snackbar.value = true
    return
  }

  // Check if quote can be accepted (status should be Approved)
  if (quote.value.status !== 'approved') {
    snackbarText.value = 'Quote must be approved before it can be accepted'
    snackbar.value = true
    return
  }

  // Show the accept quote dialog
  acceptQuoteDialog.value = true
}

const generatePdf = async (quoteId: string) => {
  try {
    router.push({
      name: 'group-pricing-quotes-generation',
      params: { quoteId }
    })
  } catch (error) {
    console.log('Error generating PDF:', error)
    snackbarText.value = 'Failed to generate PDF'
    snackbar.value = true
  }
}

const confirmAcceptQuote = async () => {
  acceptQuoteLoading.value = true
  try {
    // Format the date for the API if it's a Date object
    const formattedDate = formatDateString(
      acceptQuoteParams.value.commencementDate,
      true,
      true,
      true
    )
    // const formattedDate =
    // acceptQuoteParams.value.commencementDate instanceof Date
    //   ? acceptQuoteParams.value.commencementDate.toISOString().split('T')[0]
    //   : acceptQuoteParams.value.commencementDate

    // Make the API call with the form parameters

    await GroupPricingService.acceptQuote(
      quote.value.id,
      formattedDate,
      acceptQuoteParams.value.term.toString()
    )

    snackbarText.value = 'Quote has been accepted successfully.'
    snackbarTimeout.value = 4000
    snackbar.value = true

    // Close the dialog
    acceptQuoteDialog.value = false

    // Refresh the quote data to update the view
    await loadQuote()
  } catch (error: any) {
    console.error('Error:', error.data || error)
    snackbarText.value = 'Quote Acceptance Failed'
    snackbarTimeout.value = 2000
    snackbar.value = true
  } finally {
    acceptQuoteLoading.value = false
  }
}

/** Backend-templated On Risk letter generation. Falls back gracefully
 * when no template is configured for the insurer. */
const isGeneratingTemplatedOnRisk = ref(false)
const generateOnRiskLetterFromBackend = async () => {
  if (!quote.value?.id) return
  isGeneratingTemplatedOnRisk.value = true
  try {
    // Make sure a letter record exists (creates one if needed)
    await GroupPricingService.createOnRiskLetter(quote.value.id)
    const r = await GroupPricingService.getOnRiskLetterDocx(quote.value.id)
    let filename = `${quote.value.scheme_name || 'On_Risk_Letter'}_On_Risk_Letter_Backend.docx`
    const cd = r.headers?.['content-disposition']
    if (cd) {
      const m = /filename="?([^";]+)"?/i.exec(cd)
      if (m?.[1]) filename = m[1]
    }
    const { saveAs } = await import('file-saver')
    saveAs(r.data, filename)
    snackbarText.value = 'On Risk letter generated from template'
    snackbarTimeout.value = 4000
    snackbar.value = true
  } catch (err: any) {
    console.error('Backend On Risk letter generation failed:', err)
    if (err?.response?.status === 404) {
      snackbarText.value =
        'No On Risk letter template uploaded for this insurer. Upload one in Group Pricing Configuration → On Risk Letter Template.'
      snackbarTimeout.value = 6000
    } else {
      snackbarText.value = 'On Risk letter generation failed'
      snackbarTimeout.value = 3000
    }
    snackbar.value = true
  } finally {
    isGeneratingTemplatedOnRisk.value = false
  }
}

/** Manual button handler for re-generating the On Risk letter. */
const generateOnRiskLetterManual = async () => {
  try {
    // Ensure a letter record exists (creates one if not already present)
    await GroupPricingService.createOnRiskLetter(quote.value.id)
    const orlRes = await GroupPricingService.getOnRiskLetterData(quote.value.id)
    await generateOnRiskLetterDocx(orlRes.data)
    await generateOnRiskLetterPdf(orlRes.data)
    snackbarText.value = 'On Risk letter generated successfully'
    snackbarTimeout.value = 4000
    snackbar.value = true
  } catch (err: any) {
    console.error('On Risk letter generation failed:', err)
    snackbarText.value = 'On Risk letter generation failed'
    snackbarTimeout.value = 2000
    snackbar.value = true
  }
}

const editQuote = (quoteId: string) => {
  router.push({
    name: 'group-pricing-quote-generation-edit',
    params: { id: quoteId }
  })
}

const openDiscountDialog = async () => {
  discountRate.value = 0
  maxDiscount.value = 0
  try {
    const res = await GroupPricingService.getDiscountAuthority(
      quote.value.risk_rate_code
    )
    maxDiscount.value = (res.data?.max_discount ?? 0) * 100
  } catch {
    maxDiscount.value = 0
  }
  discountDialog.value = true
}

const confirmApplyDiscount = async () => {
  if (discountRate.value < 0 || discountRate.value > maxDiscount.value) return
  discountLoading.value = true
  try {
    await GroupPricingService.applyDiscount(quote.value.id, discountRate.value)
    snackbarText.value = 'Discount applied successfully'
    snackbar.value = true
    discountDialog.value = false
    loadQuote()
  } catch (error) {
    console.error('Error applying discount:', error)
    snackbarText.value = 'Failed to apply discount'
    snackbar.value = true
  } finally {
    discountLoading.value = false
  }
}

onMounted(() => {
  loadQuote()
})

// Handler for indicative-data-updated event (currently does nothing)
const handleIndicativeDataUpdate = (e) => {
  // You can add logic here if needed, or leave empty if not required
  // update the quote data to reflect the changes
  GroupPricingService.updateIndicativeDataFlag(
    quote.value.id,
    e.indicativeData
  ).then(() => {
    loadQuote()
  })
}
</script>
