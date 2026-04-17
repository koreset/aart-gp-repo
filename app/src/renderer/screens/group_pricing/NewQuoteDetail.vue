<template>
  <v-container v-if="quote">
    <!-- Calculation progress overlay -->
    <v-overlay
      :model-value="awaitingCalculation"
      contained
      persistent
      class="align-center justify-center"
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
                    :disabled="quote.status === 'InForce'"
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
                  : 'Requires member data and claims experience (if experience rating is enabled) to be uploaded.'
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
            <v-tooltip
              location="top"
              text="Generate On Risk letter confirming cover is active. Available after quote acceptance."
            >
              <template #activator="{ props: tooltipProps }">
                <div v-bind="tooltipProps" class="d-inline-block">
                  <v-btn
                    v-if="hasPermission('quote:generate_on_risk_letter')"
                    :disabled="
                      (quote.status !== 'accepted' &&
                        quote.status !== 'in_force' &&
                        quote.status !== 'Accepted' &&
                        quote.status !== 'InForce') ||
                      isGeneratingOnRiskLetter
                    "
                    :loading="isGeneratingOnRiskLetter"
                    size="small"
                    rounded
                    color="success"
                    class="mr-2"
                    @click="generateOnRiskLetterManual"
                    >Generate On Risk Letter</v-btn
                  >
                </div>
              </template>
            </v-tooltip>
            <v-tooltip
              location="top"
              text="Generate the On Risk letter from the insurer's uploaded Word template (if any)."
            >
              <template #activator="{ props: tooltipProps }">
                <div v-bind="tooltipProps" class="d-inline-block">
                  <v-btn
                    v-if="hasPermission('quote:generate_on_risk_letter')"
                    :disabled="
                      (quote.status !== 'accepted' &&
                        quote.status !== 'in_force' &&
                        quote.status !== 'Accepted' &&
                        quote.status !== 'InForce') ||
                      isGeneratingTemplatedOnRisk
                    "
                    :loading="isGeneratingTemplatedOnRisk"
                    size="small"
                    rounded
                    color="success"
                    variant="outlined"
                    class="mr-2"
                    @click="generateOnRiskLetterFromBackend"
                    >Generate On Risk Letter (Backend)</v-btn
                  >
                </div>
              </template>
            </v-tooltip>

            <v-tooltip
              location="top"
              text="Quote must be in 'Accepted' status."
            >
              <template #activator="{ props: tooltipProps }">
                <div v-bind="tooltipProps" class="d-inline-block">
                  <v-btn
                    v-if="hasPermission('quote:generate_pdf')"
                    :disabled="hasEmptyQuoteTables"
                    size="small"
                    rounded
                    color="primary"
                    @click="generatePdf(quote.id)"
                    >Generate PDF</v-btn
                  >
                </div>
              </template>
            </v-tooltip>
          </v-col>
        </v-row>
        <v-tabs v-model="tab" color="primary" class="mb-5">
          <v-tab value="summary">Quote Summary</v-tab>
          <v-tab value="benefits">Benefits & Config</v-tab>
          <v-tab value="data">Data Management</v-tab>
          <v-tab v-if="hasPermission('quote:view_results')" value="results"
            >Results & Analysis</v-tab
          >
          <v-tab
            v-if="
              resultSummaries !== null &&
              hasPermission('quote:view_output_summary')
            "
            value="outputsummary"
            >Output Summary</v-tab
          >
          <v-tab
            v-if="hasPermission('quote:view_premium_summary')"
            value="benefitssummary"
            >Premiums Summary</v-tab
          >
        </v-tabs>
        <v-window v-model="tab">
          <v-window-item value="summary" eager>
            <QuoteSummaryData :quote="quote" />
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
              @quote-updated="loadQuote"
              @indicative-data-updated="handleIndicativeDataUpdate"
            />
          </v-window-item>
          <v-window-item
            v-if="hasPermission('quote:view_results')"
            value="results"
            eager
          >
            <QuoteResults :quote="quote" @quote-updated="loadQuote" />
          </v-window-item>
          <v-window-item
            v-if="hasPermission('quote:view_output_summary')"
            value="outputsummary"
            eager
          >
            <OutputSummary
              v-if="resultSummaries !== null && resultSummaries.length > 0"
              :quote="quote"
              :resultSummaries="resultSummaries"
            />
          </v-window-item>
          <v-window-item
            v-if="hasPermission('quote:view_premium_summary')"
            value="benefitssummary"
            eager
          >
            <QuoteBenefitSummary
              v-if="resultSummaries !== null"
              :resultSummaries="resultSummaries"
              :quote="quote"
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
        </template>
        <template #actions>
          <v-spacer></v-spacer>
          <v-btn rounded variant="text" @click="closeBasisDialog(true)"
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
import { ref, watch, onMounted, computed } from 'vue'
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

watch(calcProgress, (val) => {
  if (!awaitingCalculation.value) return
  if (val?.phase === 'completed') {
    snackbarText.value = 'Calculations Successful'
    snackbar.value = true
    awaitingCalculation.value = false
    loadQuote()
  }
  if (val?.phase === 'failed') {
    snackbarText.value =
      'Calculations failed. Please contact your administrator.'
    snackbar.value = true
    awaitingCalculation.value = false
  }
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

const closeBasisDialog = (value) => {
  if (value) {
    runQuoteCalculations()
  }
  basisDialog.value = false
}

const loadQuote = async () => {
  try {
    const res = await GroupPricingService.getQuote(props.id)
    quote.value = res.data
    broker.value = quote.value.quoteBroker

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
      resultSummaries.value = resp.data
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
    startTracking(String(quote.value.id))
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
      awaitingCalculation.value = false
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

    // Auto-generate On Risk letter after successful acceptance
    try {
      const orlRes = await GroupPricingService.getOnRiskLetterData(
        quote.value.id
      )
      await generateOnRiskLetterDocx(orlRes.data)
      await generateOnRiskLetterPdf(orlRes.data)
      snackbarText.value = 'Quote accepted and On Risk letter generated'
      snackbarTimeout.value = 4000
      snackbar.value = true
    } catch (orlErr: any) {
      console.error('On Risk letter generation failed:', orlErr)
    }
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
