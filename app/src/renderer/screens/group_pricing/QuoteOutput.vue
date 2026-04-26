<template>
  <v-container>
    <v-row v-if="quote">
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <span class="headline"
              >Generate Quotation PDF for {{ quote.scheme_name }}</span
            >
          </template>
          <template #default>
            <v-row>
              <v-col cols="12" md="6">
                <v-btn
                  class="mr-2"
                  size="small"
                  variant="text"
                  @click="router.go(-1)"
                  >Back</v-btn
                >
              </v-col>
            </v-row>
            <v-row v-if="quote" class="mb-4 mx-7">
              <v-col cols="12">
                <v-card outlined elevation="1">
                  <v-card-title class="font-weight-bold"
                    >Quote Summary</v-card-title
                  >
                  <v-card-text>
                    <v-row>
                      <v-col cols="12" md="6">
                        <strong>Scheme Name:</strong> {{ quote.scheme_name }}
                      </v-col>
                      <v-col cols="12" md="6">
                        <strong>Quote Date:</strong>
                        {{
                          formatDateString(
                            quote.creation_date,
                            true,
                            true,
                            true
                          )
                        }}
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col cols="12" md="6">
                        <strong>Inception Date:</strong>
                        {{
                          formatDateString(
                            quote.commencement_date,
                            true,
                            true,
                            true
                          )
                        }}
                      </v-col>
                      <v-col cols="12" md="6">
                        <strong>Coverage Period:</strong> 1 year
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col cols="12" md="6">
                        <strong>Experience Rating:</strong>
                        {{ quote.experience_rating }}
                      </v-col>
                      <v-col cols="12" md="6">
                        <strong>Status:</strong> {{ quote.status || 'N/A' }}
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <v-row v-if="insurer" class="mx-7">
              <v-col>
                <v-textarea
                  v-model="insurer.introductory_text"
                  label="Introductory Text"
                  rows="4"
                  variant="outlined"
                  density="compact"
                ></v-textarea>
              </v-col>
            </v-row>
            <v-row v-if="insurer" class="mx-7">
              <v-col>
                <v-textarea
                  v-model="insurer.general_provisions_text"
                  label="Underwriting and General Provisions"
                  rows="4"
                  variant="outlined"
                  density="compact"
                ></v-textarea>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-btn size="small" color="primary" rounded @click="generatePDF"
                  >Generate and Display Quotation PDF</v-btn
                >
                <v-btn
                  class="ml-2"
                  size="small"
                  color="info"
                  rounded
                  :loading="isGeneratingWord"
                  :disabled="!quote || !insurer"
                  @click="generateWord"
                  >Generate and Download Quotation Word</v-btn
                >
                <v-btn
                  class="ml-2"
                  size="small"
                  color="info"
                  variant="outlined"
                  rounded
                  :loading="isGeneratingBackendWord"
                  :disabled="!quote"
                  @click="generateWordFromBackend"
                  >Generate Word (Backend)</v-btn
                >
                <v-btn
                  class="ml-2"
                  size="small"
                  color="primary"
                  rounded
                  @click="exportBenefitDataToExcel"
                  >Export Benefits to Excel</v-btn
                >
                <v-btn
                  class="ml-2"
                  size="small"
                  color="secondary"
                  rounded
                  @click="showGridView = !showGridView"
                  >{{ showGridView ? 'Hide' : 'Show' }} Grid View</v-btn
                >
              </v-col>
            </v-row>
            <v-row v-if="wordErrorMessage" class="mx-7">
              <v-col>
                <v-alert
                  type="error"
                  variant="tonal"
                  closable
                  @click:close="wordErrorMessage = ''"
                >
                  {{ wordErrorMessage }}
                </v-alert>
              </v-col>
            </v-row>

            <!-- AG Grid Section -->
            <v-row v-if="showGridView" class="mx-7 mb-4">
              <v-col cols="12">
                <v-card outlined elevation="1">
                  <v-card-title class="font-weight-bold"
                    >Benefits Summary Grid</v-card-title
                  >
                  <v-card-text>
                    <div
                      class="ag-theme-balham"
                      style="height: 600px; width: 100%"
                    >
                      <AgGridVue
                        :key="`quote-output-grid-${gridRowDataKey}`"
                        :rowData="gridRowData"
                        :columnDefs="gridColumnDefs"
                        :defaultColDef="defaultColDef"
                        :groupDefaultExpanded="1"
                        :autoGroupColumnDef="autoGroupColumnDef"
                        :suppressAggFuncInHeader="true"
                        :rowGroupPanelShow="'never'"
                        :animateRows="true"
                        :suppressRowClickSelection="true"
                        :rowHeight="35"
                        :headerHeight="40"
                        @grid-ready="onGridReady"
                      />
                    </div>
                    <v-row v-if="gridRowData.length === 0" class="mt-2">
                      <v-col>
                        <v-alert type="info" variant="tonal">
                          No data available for grid display. Please ensure
                          quote data is loaded.
                        </v-alert>
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <div v-if="pdfUrl" class="pdf-viewer">
              <iframe
                :src="pdfUrl"
                width="100%"
                height="800px"
                frameborder="0"
              ></iframe>
            </div>
          </template>
        </base-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import BaseCard from '@/renderer/components/BaseCard.vue'
import { onMounted, ref, watch } from 'vue'
import jsPDF from 'jspdf'
import { applyPlugin } from 'jspdf-autotable'
// import logo from '@/renderer/assets/aart-logo-01.png'
import GroupPricingService from '@/renderer/api/GroupPricingService'
// import { useRouter } from 'vue-router'
import formatDateString from '@/renderer/utils/helpers.js'
import { useRouter } from 'vue-router'
import * as XLSX from 'xlsx'
import { saveAs } from 'file-saver'
import { AgGridVue } from 'ag-grid-vue3'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-balham.css'
import {
  safeGetValue as _safeGetValue,
  roundUpToTwoDecimalsAccounting as _roundUpToTwoDecimalsAccounting,
  dashIfEmpty as _dashIfEmpty,
  officeProportionFromRiskProportion,
  finalFieldValue
} from '@/renderer/utils/quoteDataHelpers'
import { useDocxQuoteGeneration } from '@/renderer/composables/useDocxQuoteGeneration'

const props = defineProps({
  quoteId: {
    type: String,
    default: ''
  }
})

const router = useRouter()
const benefitMaps: any = ref([])
const glaBenefitTitle = ref('')
const sglaBenefitTitle = ref('')
const ptdBenefitTitle = ref('')
const ciBenefitTitle = ref('')
const phiBenefitTitle = ref('')
const ttdBenefitTitle = ref('')
const familyFuneralBenefitTitle = ref('')
const additionalAccidentalGlaBenefitTitle = ref('Additional Accidental GLA')
const additionalGlaCoverBenefitTitle = ref('Additional GLA Cover')
const quote: any = ref(null)
const resultSummaries: any = ref([])
const insurer: any = ref(null)
const categoryEducatorBenefits: any = ref([])

// const router = useRouter()

applyPlugin(jsPDF)

const pdfUrl = ref('')

// DOCX generation composable
const {
  isGenerating: isGeneratingWord,
  errorMessage: wordErrorMessage,
  generateDocxQuote
} = useDocxQuoteGeneration()

const generateWord = async () => {
  try {
    await generateDocxQuote(
      quote.value,
      resultSummaries.value,
      insurer.value,
      categoryEducatorBenefits.value,
      benefitMaps.value,
      {
        glaBenefitTitle: glaBenefitTitle.value,
        sglaBenefitTitle: sglaBenefitTitle.value,
        ptdBenefitTitle: ptdBenefitTitle.value,
        ciBenefitTitle: ciBenefitTitle.value,
        phiBenefitTitle: phiBenefitTitle.value,
        ttdBenefitTitle: ttdBenefitTitle.value,
        familyFuneralBenefitTitle: familyFuneralBenefitTitle.value,
        additionalAccidentalGlaBenefitTitle:
          additionalAccidentalGlaBenefitTitle.value,
        additionalGlaCoverBenefitTitle: additionalGlaCoverBenefitTitle.value
      }
    )
  } catch (error) {
    console.error('Error generating Word document:', error)
  }
}

// Backend-generated Word document (PoC — parallel path to the client-side generator above)
const isGeneratingBackendWord = ref(false)
const generateWordFromBackend = async () => {
  if (!quote.value?.id) return
  isGeneratingBackendWord.value = true
  try {
    const response = await GroupPricingService.getQuoteDocx(quote.value.id)
    // Try to pull the filename from Content-Disposition; fall back to a sensible default.
    let filename = `${quote.value.scheme_name || 'Quotation'}_Quotation_Backend.docx`
    const cd = response.headers?.['content-disposition']
    if (cd) {
      const match = /filename="?([^";]+)"?/i.exec(cd)
      if (match?.[1]) filename = match[1]
    }
    saveAs(response.data, filename)
  } catch (error) {
    console.error('Error generating backend Word document:', error)
  } finally {
    isGeneratingBackendWord.value = false
  }
}

// AG Grid variables
const showGridView = ref(false)
const gridApi = ref(null)
const gridRowData: any = ref([])
// Bumped every time gridRowData is rebuilt so AgGridVue remounts. Without
// this, AG Grid cannot reconcile rows across data swaps (benefit titles
// resolve asynchronously from getBenefitMaps after the first watcher fire)
// and some cells render blank even though gridRowData holds the values.
const gridRowDataKey = ref(0)
const gridColumnDefs = ref([
  {
    field: 'category',
    headerName: 'Category',
    rowGroup: true,
    hide: true
  },
  {
    field: 'benefit',
    headerName: 'Benefit',
    minWidth: 200,
    cellStyle: (params) => {
      if (params.data?.isSubtotal) {
        return { fontWeight: 'bold', backgroundColor: '#f0f0f0' }
      }
      if (params.data?.isSectionHeader) {
        return {
          fontWeight: 'bold',
          backgroundColor: '#e3f2fd',
          fontStyle: 'italic'
        }
      }
      return { fontWeight: 'bold' }
    }
  },
  {
    field: 'totalSumAssured',
    headerName: 'Total Sum Assured',
    minWidth: 150,
    valueFormatter: (params) => {
      if (
        params.value === null ||
        params.value === undefined ||
        params.value === ''
      )
        return '-'
      return typeof params.value === 'string'
        ? params.value
        : roundUpToTwoDecimalsAccounting(params.value)
    },
    type: 'rightAligned',
    cellStyle: (params) => {
      if (params.data?.isSubtotal) {
        return { backgroundColor: '#f0f0f0', textAlign: 'right' }
      }
      if (params.data?.isSectionHeader) {
        return { backgroundColor: '#e3f2fd', textAlign: 'right' }
      }
      return { textAlign: 'right' }
    }
  },
  {
    field: 'annualPremium',
    headerName: 'Annual Premium',
    minWidth: 150,
    valueFormatter: (params) => {
      if (
        params.value === null ||
        params.value === undefined ||
        params.value === ''
      )
        return '-'
      return typeof params.value === 'string'
        ? params.value
        : roundUpToTwoDecimalsAccounting(params.value)
    },
    type: 'rightAligned',
    cellStyle: (params) => {
      if (params.data?.isSubtotal) {
        return { backgroundColor: '#f0f0f0', textAlign: 'right' }
      }
      if (params.data?.isSectionHeader) {
        return { backgroundColor: '#e3f2fd', textAlign: 'right' }
      }
      return { textAlign: 'right' }
    }
  },
  {
    field: 'percentSalary',
    headerName: '% of Salary',
    minWidth: 120,
    type: 'rightAligned',
    cellStyle: (params) => {
      if (params.data?.isSubtotal) {
        return { backgroundColor: '#f0f0f0', textAlign: 'right' }
      }
      if (params.data?.isSectionHeader) {
        return { backgroundColor: '#e3f2fd', textAlign: 'right' }
      }
      return { textAlign: 'right' }
    }
  }
])

const defaultColDef = ref({
  sortable: true,
  filter: true,
  resizable: true,
  flex: 1
})

const autoGroupColumnDef = ref({
  headerName: 'Category',
  minWidth: 200,
  cellRendererParams: {
    suppressCount: true
  }
})

// Delegate to shared helpers (keeps all call sites in this file working unchanged)
const safeGetValue = _safeGetValue
const roundUpToTwoDecimalsAccounting = _roundUpToTwoDecimalsAccounting
const dashIfEmpty = _dashIfEmpty

// AG Grid event handler
const onGridReady = (params) => {
  gridApi.value = params.api
}

// Function to convert Excel data to AG Grid format
const convertExcelDataToGridData = () => {
  if (!resultSummaries.value || resultSummaries.value.length === 0) return []

  const gridData: any = []

  resultSummaries.value.forEach((resultSummary) => {
    const category = resultSummary.category

    // Add benefit rows
    gridData.push({
      category,
      benefit: glaBenefitTitle.value,
      totalSumAssured: resultSummary.total_gla_capped_sum_assured,
      annualPremium: finalFieldValue(
        resultSummary,
        'final_gla_annual_office_premium'
      ),
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(resultSummary.exp_proportion_gla_annual_risk_premium_salary, resultSummary) * 100)}%`
    })

    if (resultSummary.total_additional_accidental_gla_capped_sum_assured > 0) {
      gridData.push({
        category,
        benefit: additionalAccidentalGlaBenefitTitle.value,
        totalSumAssured:
          resultSummary.total_additional_accidental_gla_capped_sum_assured,
        annualPremium: finalFieldValue(
          resultSummary,
          'final_additional_accidental_gla_annual_office_premium'
        ),
        percentSalary: `${roundUpToTwoDecimalsAccounting((officeProportionFromRiskProportion(resultSummary.exp_proportion_additional_accidental_gla_annual_risk_premium_salary, resultSummary) || 0) * 100)}%`
      })
    }

    gridData.push({
      category,
      benefit: ptdBenefitTitle.value,
      totalSumAssured: resultSummary.total_ptd_capped_sum_assured,
      annualPremium: finalFieldValue(
        resultSummary,
        'final_ptd_annual_office_premium'
      ),
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(resultSummary.exp_proportion_ptd_annual_risk_premium_salary, resultSummary) * 100)}%`
    })

    gridData.push({
      category,
      benefit: ciBenefitTitle.value,
      totalSumAssured: resultSummary.total_ci_capped_sum_assured,
      annualPremium: finalFieldValue(
        resultSummary,
        'final_ci_annual_office_premium'
      ),
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(resultSummary.exp_proportion_ci_annual_risk_premium_salary, resultSummary) * 100)}%`
    })

    gridData.push({
      category,
      benefit: sglaBenefitTitle.value,
      totalSumAssured: resultSummary.total_sgla_capped_sum_assured,
      annualPremium: finalFieldValue(
        resultSummary,
        'final_sgla_annual_office_premium'
      ),
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(resultSummary.exp_proportion_sgla_annual_risk_premium_salary, resultSummary) * 100)}%`
    })

    gridData.push({
      category,
      benefit: phiBenefitTitle.value,
      totalSumAssured: resultSummary.total_phi_capped_income,
      annualPremium: finalFieldValue(
        resultSummary,
        'final_phi_annual_office_premium'
      ),
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(resultSummary.exp_proportion_phi_annual_risk_premium_salary, resultSummary) * 100)}%`
    })

    gridData.push({
      category,
      benefit: ttdBenefitTitle.value,
      totalSumAssured: resultSummary.total_ttd_capped_income,
      annualPremium: finalFieldValue(
        resultSummary,
        'final_ttd_annual_office_premium'
      ),
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(resultSummary.exp_proportion_ttd_annual_risk_premium_salary, resultSummary) * 100)}%`
    })

    // Add subtotal row
    gridData.push({
      category,
      benefit: 'Sub Total/Total Premiums',
      totalSumAssured: resultSummary.total_gla_capped_sum_assured,
      annualPremium: finalFieldValue(resultSummary, 'final_total_annual_premium_excl_funeral'),
      percentSalary: `${roundUpToTwoDecimalsAccounting(resultSummary.proportion_exp_total_premium_excl_funeral_salary * 100)}%`,
      isSubtotal: true
    })

    // Add Group Funeral section header
    gridData.push({
      category,
      benefit: 'Group Funeral',
      totalSumAssured: '',
      annualPremium: '',
      percentSalary: '',
      isSectionHeader: true
    })

    // Add Group Funeral rows
    gridData.push({
      category,
      benefit: 'Monthly Premium per Member',
      totalSumAssured: '',
      annualPremium: roundUpToTwoDecimalsAccounting(
        resultSummary.exp_total_fun_monthly_premium_per_member
      ),
      percentSalary: ''
    })

    gridData.push({
      category,
      benefit: 'Annual Premium per Member',
      totalSumAssured: '',
      annualPremium: roundUpToTwoDecimalsAccounting(
        resultSummary.exp_total_fun_annual_premium_per_member
      ),
      percentSalary: ''
    })

    gridData.push({
      category,
      benefit: 'Total Annual Premium',
      totalSumAssured: '',
      annualPremium: roundUpToTwoDecimalsAccounting(
        finalFieldValue(
          resultSummary,
          'final_fun_annual_office_premium'
        )
      ),
      percentSalary: ''
    })
  })

  return gridData
}

// Watch for data changes and update grid
watch(
  [
    resultSummaries,
    glaBenefitTitle,
    sglaBenefitTitle,
    ptdBenefitTitle,
    ciBenefitTitle,
    phiBenefitTitle,
    ttdBenefitTitle
  ],
  () => {
    if (resultSummaries.value.length > 0) {
      gridRowData.value = convertExcelDataToGridData()
      gridRowDataKey.value++
    }
  },
  { deep: true }
)

onMounted(async () => {
  // Generate PDF on component mount
  GroupPricingService.getBenefitMaps().then((res) => {
    benefitMaps.value = res.data
    console.log('Benefit Maps:', benefitMaps.value)
    const glaBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'GLA'
    )
    if (glaBenefit.benefit_alias !== '') {
      glaBenefitTitle.value = glaBenefit.benefit_alias
    } else {
      glaBenefitTitle.value = glaBenefit.benefit_name
    }
    const sglaBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'SGLA'
    )
    if (sglaBenefit.benefit_alias !== '') {
      sglaBenefitTitle.value = sglaBenefit.benefit_alias
    } else {
      sglaBenefitTitle.value = sglaBenefit.benefit_name
    }
    const ptdBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'PTD'
    )
    if (ptdBenefit.benefit_alias !== '') {
      ptdBenefitTitle.value = ptdBenefit.benefit_alias
    } else {
      ptdBenefitTitle.value = ptdBenefit.benefit_name
    }
    const ciBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'CI'
    )
    if (ciBenefit.benefit_alias !== '') {
      ciBenefitTitle.value = ciBenefit.benefit_alias
    } else {
      ciBenefitTitle.value = ciBenefit.benefit_name
    }
    const phiBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'PHI'
    )
    if (phiBenefit.benefit_alias !== '') {
      phiBenefitTitle.value = phiBenefit.benefit_alias
    } else {
      phiBenefitTitle.value = phiBenefit.benefit_name
    }
    const ttdBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'TTD'
    )
    if (ttdBenefit.benefit_alias !== '') {
      ttdBenefitTitle.value = ttdBenefit.benefit_alias
    } else {
      ttdBenefitTitle.value = ttdBenefit.benefit_name
    }
    const familyFuneralBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'GFF'
    )
    if (familyFuneralBenefit.benefit_alias !== '') {
      familyFuneralBenefitTitle.value = familyFuneralBenefit.benefit_alias
    } else {
      familyFuneralBenefitTitle.value = familyFuneralBenefit.benefit_name
    }
    const additionalAccidentalGlaBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'AAGLA'
    )
    if (additionalAccidentalGlaBenefit) {
      additionalAccidentalGlaBenefitTitle.value =
        additionalAccidentalGlaBenefit.benefit_alias ||
        additionalAccidentalGlaBenefit.benefit_name
    }
    const additionalGlaCoverBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'AGLA'
    )
    if (additionalGlaCoverBenefit) {
      additionalGlaCoverBenefitTitle.value =
        additionalGlaCoverBenefit.benefit_alias ||
        additionalGlaCoverBenefit.benefit_name
    }
  })

  GroupPricingService.getQuote(props.quoteId).then((res) => {
    quote.value = res.data
    console.log('Quote data:', quote.value)
  })
  GroupPricingService.getResultSummary(props.quoteId).then((res) => {
    resultSummaries.value = res.data
    console.log('Result Summary:', resultSummaries.value)
  })
  GroupPricingService.getCategoryEducatorBenefits(props.quoteId).then((res) => {
    categoryEducatorBenefits.value = res.data
    console.log('Category Educator Benefits:', categoryEducatorBenefits.value)
  })

  GroupPricingService.getInsurer().then((res) => {
    insurer.value = res.data
    console.log('Insurer data:', insurer.value)
  })
})

const exportBenefitDataToExcel = () => {
  if (
    !quote.value ||
    !resultSummaries.value ||
    resultSummaries.value.length === 0
  )
    return

  const wsData: any = []
  wsData.push(['Benefit Summary for ' + quote.value.scheme_name])
  wsData.push([]) // Empty row for spacing
  console.log('Exporting benefit data to Excel...', resultSummaries.value)

  resultSummaries.value.forEach((resultSummary) => {
    console.log('Result Summary:', resultSummary.category)
    // Example: adjust this structure to match your screenshot layout
    console.log('Result Summary:', resultSummary.category)
    const wsDataCategory = [
      [`Benefit Summary for ${resultSummary.category}`],
      // Headers
      ['Benefit', 'Total Sum Assured', 'Annual Premium', '% of Salary'],
      // Data rows
      [
        glaBenefitTitle.value,
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            resultSummary.total_gla_capped_sum_assured
          )
        ),
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            finalFieldValue(
              resultSummary,
              'final_gla_annual_office_premium'
            )
          )
        ),
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            officeProportionFromRiskProportion(
              resultSummary.exp_proportion_gla_annual_risk_premium_salary,
              resultSummary
            ) * 100
          ) + '%'
        )
      ],
      [
        ptdBenefitTitle.value,
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            resultSummary.total_ptd_capped_sum_assured
          )
        ),
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            finalFieldValue(
              resultSummary,
              'final_ptd_annual_office_premium'
            )
          )
        ),
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            officeProportionFromRiskProportion(
              resultSummary.exp_proportion_ptd_annual_risk_premium_salary,
              resultSummary
            ) * 100
          ) + '%'
        )
      ],
      [
        ciBenefitTitle.value,
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            resultSummary.total_ci_capped_sum_assured
          )
        ),
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            finalFieldValue(
              resultSummary,
              'final_ci_annual_office_premium'
            )
          )
        ),
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            officeProportionFromRiskProportion(
              resultSummary.exp_proportion_ci_annual_risk_premium_salary,
              resultSummary
            ) * 100
          ) + '%'
        )
      ],
      [
        sglaBenefitTitle.value,
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            resultSummary.total_sgla_capped_sum_assured
          )
        ),
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            finalFieldValue(
              resultSummary,
              'final_sgla_annual_office_premium'
            )
          )
        ),
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            officeProportionFromRiskProportion(
              resultSummary.exp_proportion_sgla_annual_risk_premium_salary,
              resultSummary
            ) * 100
          ) + '%'
        )
      ],
      [
        phiBenefitTitle.value,
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(resultSummary.total_phi_capped_income)
        ),
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            finalFieldValue(
              resultSummary,
              'final_phi_annual_office_premium'
            )
          )
        ),
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            officeProportionFromRiskProportion(
              resultSummary.exp_proportion_phi_annual_risk_premium_salary,
              resultSummary
            ) * 100
          ) + '%'
        )
      ],
      [
        ttdBenefitTitle.value,
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(resultSummary.total_ttd_capped_income)
        ),
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            finalFieldValue(
              resultSummary,
              'final_ttd_annual_office_premium'
            )
          )
        ),
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            officeProportionFromRiskProportion(
              resultSummary.exp_proportion_ttd_annual_risk_premium_salary,
              resultSummary
            ) * 100
          ) + '%'
        )
      ],
      [
        'Sub Total/Total Premiums',
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            resultSummary.total_gla_capped_sum_assured
          )
        ),
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            finalFieldValue(resultSummary, 'final_total_annual_premium_excl_funeral')
          )
        ),
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            resultSummary.proportion_exp_total_premium_excl_funeral_salary * 100
          ) + '%'
        )
      ],
      [],
      // Group Funeral section
      ['Group Funeral'],
      [
        'Monthly Premium',
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            resultSummary.exp_total_fun_monthly_premium_per_member
          )
        )
      ],
      [
        'Annual Premium',
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            resultSummary.exp_total_fun_annual_premium_per_member
          )
        )
      ],
      [
        'Total Annual Premium',
        dashIfEmpty(
          roundUpToTwoDecimalsAccounting(
            finalFieldValue(
              resultSummary,
              'final_fun_annual_office_premium'
            )
          )
        )
      ]
    ]

    console.log('wsdata: ', wsData)
    // Add the category data to the main worksheet data
    wsData.push(...wsDataCategory)
    wsData.push([], []) // Empty row for spacing between categories
  })

  // Create worksheet and workbook
  const ws = XLSX.utils.aoa_to_sheet(wsData)
  // Merge the title row across all columns
  ws['!merges'] = [
    { s: { r: 0, c: 0 }, e: { r: 0, c: 3 } } // Merge A1:D1
  ]
  // Optionally set column widths
  ws['!cols'] = [{ wch: 30 }, { wch: 20 }, { wch: 20 }, { wch: 15 }]
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, 'Benefit Summary')
  // Export to file
  XLSX.writeFile(wb, `benefits_summary_${quote.value.scheme_name}.xlsx`)
}

const generatePDF = async () => {
  // eslint-disable-next-line new-cap
  const doc: any = new jsPDF()

  // Design constants
  const colors = {
    primary: [52, 73, 94], // Dark blue-gray
    secondary: [52, 73, 94], // Dark blue-gray
    accent: [231, 76, 60], // Red accent
    light: [236, 240, 241], // Light gray
    dark: [44, 62, 80], // Dark gray
    white: [255, 255, 255]
  }

  const fonts = {
    primary: 'helvetica',
    sizes: {
      title: 16,
      heading: 14,
      subheading: 12,
      body: 10,
      caption: 9
    }
  }

  // --- Header and Logo Section ---
  const pageWidth = doc.internal.pageSize.getWidth()
  const logoWidth = 80
  const logoHeight = 30
  const logoX = pageWidth - logoWidth - 14 // Right margin

  // Improved margins and spacing
  const topMargin = 20
  let currentY = topMargin
  const leftMargin = 18
  const rightMargin = 18
  const contentWidth = pageWidth - leftMargin - rightMargin

  // Enhanced text wrapping function
  const wrapText = (doc, text, x, y, maxWidth) => {
    if (!text || typeof text !== 'string') {
      return y
    }
    const lines = doc.splitTextToSize(text, maxWidth)
    doc.text(lines, x, y)
    const lineHeight = doc.getLineHeight() / doc.internal.scaleFactor
    return y + lines.length * lineHeight
  }

  const addPageHeaderFooter = (doc, pageNumber) => {
    const pageWidth = doc.internal.pageSize.getWidth()
    const pageHeight = doc.internal.pageSize.getHeight()
    // Header
    doc.setFillColor(...colors.light)
    doc.rect(0, 0, pageWidth, 15, 'F')
    doc.setTextColor(...colors.secondary)
    doc.setFontSize(fonts.sizes.caption)
    doc.setFont(fonts.primary, 'normal')
    doc.text(`${quote.value.scheme_name} - Quotation`, leftMargin, 10)

    // Footer
    doc.setFillColor(...colors.light)
    doc.rect(0, pageHeight - 15, pageWidth, 15, 'F')
    doc.setTextColor(...colors.secondary)
    doc.text(`Page ${pageNumber}`, pageWidth - 25, pageHeight - 7)
    doc.text(
      `Generated on ${formatDateString(new Date(), true, true, true)}`,
      leftMargin,
      pageHeight - 7
    )
  }

  // Conservative page break detection function - only breaks when absolutely necessary
  const checkPageBreak = (doc, currentY, requiredHeight) => {
    const pageHeight = doc.internal.pageSize.getHeight()
    const bottomMargin = 20 // Standard bottom margin

    if (currentY + requiredHeight > pageHeight - bottomMargin) {
      return true // Indicates a break is needed
    }
    return false // No page break needed
  }

  let pageNumber = 1
  addPageHeaderFooter(doc, pageNumber)

  // Add subtle header background
  doc.setFillColor(...colors.light)
  doc.rect(0, 0, pageWidth, 55, 'F')

  // 3. Construct the Base64 Data URL
  const logoUrl = `data:${insurer.value.logo_mime_type};base64,${insurer.value.logo}`

  const logoData = {
    imageFormat: insurer.value.logo_mime_type.split('/')[1], // Extract format from MIME type
    imageWidth: logoWidth,
    imageHeight: logoHeight
  }

  doc.addImage(
    logoUrl,
    logoData.imageFormat.toUpperCase(),
    logoX,
    currentY,
    logoWidth,
    logoHeight
  )

  // Enhanced header text styling
  doc.setTextColor(...colors.dark)
  doc.setFont(fonts.primary, 'bold')
  doc.setFontSize(fonts.sizes.subheading)
  doc.text(`${insurer.value.name}`, leftMargin, currentY + 2)

  doc.setFont(fonts.primary, 'normal')
  doc.setFontSize(fonts.sizes.caption)
  doc.setTextColor(...colors.secondary)
  currentY += 10
  doc.text(
    `${insurer.value.address_line_1}, ${insurer.value.address_line_2}`,
    leftMargin,
    currentY
  )

  if (
    insurer.value.address_line_3 &&
    insurer.value.address_line_3.trim() !== ''
  ) {
    currentY += 4
    doc.text(`${insurer.value.address_line_3}`, leftMargin, currentY)
  }
  currentY += 4
  doc.text(
    `${insurer.value.city}, ${insurer.value.province}, ${insurer.value.post_code}`,
    leftMargin,
    currentY
  )
  currentY += 4
  doc.text(`Tel: ${insurer.value.telephone}`, leftMargin, currentY)
  currentY += 4
  doc.text(`Email: ${insurer.value.email}`, leftMargin, currentY)
  // --------------------------------

  // **Enhanced Title Section**
  currentY += 20

  // doc.setFillColor(...colors.primary)
  // doc.rect(leftMargin - 5, currentY - 8, contentWidth + 10, 15, 'F')

  doc.setTextColor(...colors.dark)
  doc.setFontSize(fonts.sizes.title)
  doc.setFont(fonts.primary, 'bold')
  const text = 'Group Risk Quotation'
  const textWidth = doc.getTextWidth(text)
  const centerX = (pageWidth - textWidth) / 2
  doc.text('Group Risk Quotation', centerX, currentY + 5)

  // Reset text color for body content
  doc.setTextColor(...colors.dark)
  doc.setFontSize(fonts.sizes.body)
  doc.setFont(fonts.primary, 'normal')
  currentY += 12

  if (insurer.value.introductory_text.trim() !== '') {
    currentY = wrapText(
      doc,
      insurer.value.introductory_text,
      leftMargin,
      currentY,
      contentWidth + 3
    )
    currentY += 5 // Adjust margin after wrapping text
  } else {
    doc.text(
      'We are pleased to submit for your consideration the quotation you requested for the above scheme.',
      leftMargin,
      currentY,
      { align: 'justify', maxWidth: contentWidth }
    )
    currentY += 5
  }

  // **Initial Information Table**
  const totalLivesCovered = resultSummaries.value.reduce(
    (sum, item) => sum + item.member_count,
    0
  )
  const totalSumAssured = resultSummaries.value.reduce(
    (sum, item) => sum + item.total_sum_assured,
    0
  )
  const totalAnnualSalary = resultSummaries.value.reduce(
    (sum, item) => sum + item.total_annual_salary,
    0
  )
  const totalAnnualPremium = resultSummaries.value.reduce(
    (sum, item) => sum + item.total_annual_premium,
    0
  )

  // Cover page executive summary table data with enhanced styling
  const initialInfo = [
    ['Type of Policy:', 'Group Risk Assurance'],
    ['Quote Number:', `${quote.value.quote_name}`],
    [
      'Quote Date:',
      `${formatDateString(quote.value.creation_date, true, true, true)}`
    ],
    ['Scheme Name:', `${quote.value.scheme_name}`],
    [
      'Inception Date:',
      `${formatDateString(quote.value.commencement_date, true, true, true)}`
    ],
    // ['Coverage Period:', '1 year'],
    ['Number of Lives Covered:', `${totalLivesCovered}`],
    [
      'Total Sum Assured:',
      `${roundUpToTwoDecimalsAccounting(totalSumAssured)}`
    ],
    [
      'Total Annual Salary:',
      `${roundUpToTwoDecimalsAccounting(totalAnnualSalary)}`
    ],
    [
      'Total Annual Premium:',
      `${roundUpToTwoDecimalsAccounting(totalAnnualPremium)}`
    ]
  ]

  // Add the initial information table with enhanced styling
  doc.autoTable({
    startY: currentY,
    body: initialInfo,
    theme: 'grid',
    styles: {
      fontSize: fonts.sizes.body,
      cellPadding: { top: 4, right: 8, bottom: 4, left: 8 },
      lineColor: colors.light,
      lineWidth: 0.5
    },
    columnStyles: {
      0: {
        fontStyle: 'bold',
        fillColor: colors.light,
        textColor: colors.dark,
        halign: 'left'
      },
      1: {
        textColor: colors.secondary
      }
    },
    alternateRowStyles: {
      fillColor: [250, 250, 250]
    },
    didDrawPage: (data) => {
      // Add header and footer to each page
      addPageHeaderFooter(doc, pageNumber)
    }
  })

  // Add page break before Premium Summary section
  doc.addPage('a4', 'landscape')
  pageNumber++
  addPageHeaderFooter(doc, pageNumber)
  currentY = topMargin + 10

  doc.setFontSize(fonts.sizes.body)
  doc.setFont(fonts.primary, 'normal')
  doc.setTextColor(...colors.secondary)

  // Check if any category has non-funeral benefits enabled
  const hasAnyNonFuneralBenefits = resultSummaries.value.some((item) => {
    return (
      safeGetValue(item, 'total_gla_capped_sum_assured', 0) > 0 ||
      safeGetValue(item, 'total_ptd_capped_sum_assured', 0) > 0 ||
      safeGetValue(item, 'total_ci_capped_sum_assured', 0) > 0 ||
      safeGetValue(item, 'total_sgla_capped_sum_assured', 0) > 0 ||
      safeGetValue(item, 'total_phi_capped_income', 0) > 0 ||
      safeGetValue(item, 'total_ttd_capped_income', 0) > 0
    )
  })

  // Only display Premium Summary if any category has non-funeral benefits
  if (hasAnyNonFuneralBenefits) {
    doc.text(
      'The following table provides a summary of the benefits and premiums for the scheme:',
      leftMargin,
      currentY
    )

    currentY += 10
    // Enhanced Premium Summary Header
    doc.setFontSize(fonts.sizes.heading)
    doc.setFont(fonts.primary, 'bold')
    doc.setTextColor(...colors.primary)
    doc.text('Premium Summary', leftMargin, currentY)

    // Add underline accent
    doc.setDrawColor(...colors.primary)
    doc.setLineWidth(1)
    doc.line(leftMargin, currentY + 2, leftMargin + 60, currentY + 2)

    currentY += 8

    // Premium Summary Table
    const premiumSummaryData = [
      [
        'Category',
        'No of Lives',
        'Total Salary',
        'Total Sum Assured',
        'Annual Premium',
        '%Salary'
      ]
    ]

    resultSummaries.value.forEach((item) => {
      const hasBenefits =
        safeGetValue(item, 'total_gla_capped_sum_assured', 0) > 0 ||
        safeGetValue(item, 'total_ptd_capped_sum_assured', 0) > 0 ||
        safeGetValue(item, 'total_ci_capped_sum_assured', 0) > 0 ||
        safeGetValue(item, 'total_sgla_capped_sum_assured', 0) > 0 ||
        safeGetValue(item, 'total_phi_capped_income', 0) > 0 ||
        safeGetValue(item, 'total_ttd_capped_income', 0) > 0

      if (hasBenefits) {
        premiumSummaryData.push([
          item.category,
          item.member_count.toString(),
          roundUpToTwoDecimalsAccounting(item.total_annual_salary),
          roundUpToTwoDecimalsAccounting(item.total_sum_assured),
          roundUpToTwoDecimalsAccounting(
            finalFieldValue(item, 'final_total_annual_premium_excl_funeral')
          ),
          `${roundUpToTwoDecimalsAccounting(item.proportion_exp_total_premium_excl_funeral_salary * 100)}%`
        ])
      }
    })

    premiumSummaryData.push([
      'Total',
      resultSummaries.value.reduce((sum, item) => sum + item.member_count, 0),
      roundUpToTwoDecimalsAccounting(
        resultSummaries.value.reduce(
          (sum, item) => sum + item.total_annual_salary,
          0
        )
      ),
      roundUpToTwoDecimalsAccounting(
        resultSummaries.value.reduce(
          (sum, item) => sum + item.total_sum_assured,
          0
        )
      ),
      roundUpToTwoDecimalsAccounting(
        resultSummaries.value.reduce(
          (sum, item) => sum + finalFieldValue(item, 'final_total_annual_premium_excl_funeral'),
          0
        )
      ),
      `${roundUpToTwoDecimalsAccounting(
        (resultSummaries.value.reduce(
          (sum, item) => sum + finalFieldValue(item, 'final_total_annual_premium_excl_funeral'),
          0
        ) /
          resultSummaries.value.reduce(
            (sum, item) => sum + item.total_annual_salary,
            0
          )) *
          100
      )}%`
    ])

    doc.autoTable({
      startY: currentY,
      head: premiumSummaryData.slice(0, 1),
      body: premiumSummaryData.slice(1),
      theme: 'striped',
      headStyles: {
        fillColor: colors.primary,
        textColor: colors.white,
        fontSize: fonts.sizes.body,
        fontStyle: 'bold',
        halign: 'center'
      },
      styles: {
        fontSize: fonts.sizes.caption,
        cellPadding: { top: 4, right: 6, bottom: 4, left: 6 },
        lineColor: colors.light,
        lineWidth: 0.3
      },
      columnStyles: {
        0: { fontStyle: 'bold', fillColor: [248, 249, 250] },
        1: { halign: 'center' },
        2: { halign: 'right' },
        3: { halign: 'right' },
        4: { halign: 'right', fontStyle: 'bold' },
        5: { halign: 'center' }
      },
      alternateRowStyles: {
        fillColor: [252, 253, 254]
      },
      didDrawPage: (data) => {
        addPageHeaderFooter(doc, pageNumber)
        const pageHeight = doc.internal.pageSize.getHeight()
        const bottomMargin = 60 // A reasonable margin

        if (data.cursor.y > pageHeight - bottomMargin) {
          doc.addPage()
          pageNumber++
          addPageHeaderFooter(doc, pageNumber)
          currentY = topMargin // Reset Y position for the new page
        }
      }
    })

    currentY = doc.lastAutoTable.finalY + 10
  }

  // **Enhanced Group Funeral Section**
  doc.setFontSize(fonts.sizes.heading)
  doc.setFont(fonts.primary, 'bold')
  doc.setTextColor(...colors.primary)
  // Adjust currentY based on whether Premium Summary was displayed
  currentY = hasAnyNonFuneralBenefits
    ? doc.lastAutoTable.finalY + 15
    : currentY + 15
  doc.text('Group Funeral', leftMargin, currentY)

  // Add underline accent
  doc.setDrawColor(...colors.primary)
  doc.setLineWidth(1)
  doc.line(leftMargin, currentY + 2, leftMargin + 50, currentY + 2)

  currentY += 8
  const groupFuneralData = [
    [
      'Category',
      'No of Lives',
      'Monthly Premium',
      'Annual Premium',
      'Total Annual Premium'
    ]
  ]
  resultSummaries.value.forEach((item) => {
    groupFuneralData.push([
      item.category,
      item.member_count.toString(),
      roundUpToTwoDecimalsAccounting(
        item.exp_total_fun_monthly_premium_per_member
      ),
      roundUpToTwoDecimalsAccounting(
        item.exp_total_fun_annual_premium_per_member
      ),
      roundUpToTwoDecimalsAccounting(
        finalFieldValue(item, 'final_fun_annual_office_premium')
      )
    ])
  })
  groupFuneralData.push([
    'Total',
    resultSummaries.value.reduce((sum, item) => sum + item.member_count, 0),
    roundUpToTwoDecimalsAccounting(
      resultSummaries.value.reduce(
        (sum, item) => sum + item.exp_total_fun_monthly_premium_per_member,
        0
      )
    ),
    roundUpToTwoDecimalsAccounting(
      resultSummaries.value.reduce(
        (sum, item) => sum + item.exp_total_fun_annual_premium_per_member,
        0
      )
    ),
    roundUpToTwoDecimalsAccounting(
      resultSummaries.value.reduce(
        (sum, item) =>
          sum +
          finalFieldValue(item, 'final_fun_annual_office_premium'),
        0
      )
    )
  ])

  doc.autoTable({
    startY: currentY,
    head: groupFuneralData.slice(0, 1),
    body: groupFuneralData.slice(1),
    theme: 'striped',
    headStyles: {
      fillColor: colors.secondary,
      textColor: colors.white,
      fontSize: fonts.sizes.body,
      fontStyle: 'bold',
      halign: 'center'
    },
    styles: {
      fontSize: fonts.sizes.caption,
      cellPadding: { top: 4, right: 6, bottom: 4, left: 6 },
      lineColor: colors.light,
      lineWidth: 0.3
    },
    columnStyles: {
      0: { fontStyle: 'bold', fillColor: [248, 249, 250] },
      1: { halign: 'center' },
      2: { halign: 'right' },
      3: { halign: 'right' },
      4: { halign: 'right', fontStyle: 'bold' }
    },
    alternateRowStyles: {
      fillColor: [252, 253, 254]
    },
    didDrawPage: (data) => {
      addPageHeaderFooter(doc, pageNumber)
    }
  })

  // New Section for Premium Breakdown

  doc.addPage('a4', 'landscape')
  pageNumber++
  addPageHeaderFooter(doc, pageNumber)

  // reset currentY for new page with header space
  currentY = topMargin + 10

  // Enhanced page 2 title
  doc.setTextColor(...colors.dark)
  doc.setFontSize(fonts.sizes.heading)
  doc.setFont(fonts.primary, 'bold')
  doc.text('Premium Breakdown', leftMargin, currentY)

  // Add underline accent
  doc.setDrawColor(...colors.primary)
  doc.setLineWidth(1)
  doc.line(leftMargin, currentY + 2, leftMargin + 50, currentY + 2)

  // **Premium Breakdown Table per category (resultSummaries)**
  doc.setTextColor(...colors.dark)
  resultSummaries.value.forEach((item, index) => {
    currentY += 15

    // Category header with background
    doc.setTextColor(...colors.dark)
    doc.setFontSize(fonts.sizes.subheading)
    doc.setFont(fonts.primary, 'bold')
    if (hasAnyNonFuneralBenefits) {
      doc.text(`${item.category} Category`, leftMargin, currentY)
    }

    const premiumBreakdownData = [
      ['Benefit', 'Total Sum Assured', 'Annual Premium', '% Salary']
    ]

    premiumBreakdownData.push([
      glaBenefitTitle.value,
      roundUpToTwoDecimalsAccounting(item.total_gla_capped_sum_assured),
      roundUpToTwoDecimalsAccounting(
        finalFieldValue(item, 'final_gla_annual_office_premium')
      ),
      `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(item.exp_proportion_gla_annual_risk_premium_salary, item) * 100)}%`
    ])
    premiumBreakdownData.push([
      sglaBenefitTitle.value,
      roundUpToTwoDecimalsAccounting(item.total_sgla_capped_sum_assured),
      roundUpToTwoDecimalsAccounting(
        finalFieldValue(item, 'final_sgla_annual_office_premium')
      ),
      `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(item.exp_proportion_sgla_annual_risk_premium_salary, item) * 100)}%`
    ])
    premiumBreakdownData.push([
      ptdBenefitTitle.value,
      roundUpToTwoDecimalsAccounting(item.total_ptd_capped_sum_assured),
      roundUpToTwoDecimalsAccounting(
        finalFieldValue(item, 'final_ptd_annual_office_premium')
      ),
      `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(item.exp_proportion_ptd_annual_risk_premium_salary, item) * 100)}%`
    ])
    premiumBreakdownData.push([
      ciBenefitTitle.value,
      roundUpToTwoDecimalsAccounting(item.total_ci_capped_sum_assured),
      roundUpToTwoDecimalsAccounting(
        finalFieldValue(item, 'final_ci_annual_office_premium')
      ),
      `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(item.exp_proportion_ci_annual_risk_premium_salary, item) * 100)}%`
    ])
    premiumBreakdownData.push([
      phiBenefitTitle.value,
      roundUpToTwoDecimalsAccounting(item.total_phi_capped_income),
      roundUpToTwoDecimalsAccounting(
        finalFieldValue(item, 'final_phi_annual_office_premium')
      ),
      `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(item.exp_proportion_phi_annual_risk_premium_salary, item) * 100)}%`
    ])
    premiumBreakdownData.push([
      ttdBenefitTitle.value,
      roundUpToTwoDecimalsAccounting(item.total_ttd_capped_income),
      roundUpToTwoDecimalsAccounting(
        finalFieldValue(item, 'final_ttd_annual_office_premium')
      ),
      `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(item.exp_proportion_ttd_annual_risk_premium_salary, item) * 100)}%`
    ])

    if (hasAnyNonFuneralBenefits) {
      currentY += 4
      doc.autoTable({
        startY: currentY,
        head: premiumBreakdownData.slice(0, 1),
        body: premiumBreakdownData.slice(1),
        theme: 'grid',
        headStyles: {
          fillColor: colors.primary,
          textColor: colors.white,
          fontSize: fonts.sizes.body,
          fontStyle: 'bold',
          halign: 'center'
        },
        styles: {
          fontSize: fonts.sizes.caption,
          cellPadding: { top: 3, right: 5, bottom: 3, left: 5 },
          lineColor: colors.light,
          lineWidth: 0.3
        },
        columnStyles: {
          0: { fontStyle: 'bold', fillColor: [248, 249, 250] },
          1: { halign: 'right' },
          2: { halign: 'right', fontStyle: 'bold' },
          3: { halign: 'center' }
        },
        alternateRowStyles: {
          fillColor: [253, 253, 254]
        },
        didDrawPage: (data) => {
          addPageHeaderFooter(doc, pageNumber)
          const pageHeight = doc.internal.pageSize.getHeight()
          const bottomMargin = 60 // A reasonable margin

          if (data.cursor.y > pageHeight - bottomMargin) {
            doc.addPage()
            pageNumber++
            addPageHeaderFooter(doc, pageNumber)
            currentY = topMargin // Reset Y position for the new page
          }
        }
      })
      currentY = doc.lastAutoTable.finalY + 10
    }

    // Group Funeral subsection with styled header
    doc.setTextColor(...colors.dark)
    doc.setFontSize(fonts.sizes.body)
    doc.setFont(fonts.primary, 'bold')
    doc.text(`${item.category} - Group Funeral`, leftMargin, currentY)
    currentY += 4

    const groupFuneralBreakdownData = [
      [
        'Monthly Premium per Member',
        roundUpToTwoDecimalsAccounting(
          item.exp_total_fun_monthly_premium_per_member
        )
      ],
      [
        'Annual Premium per Member',
        roundUpToTwoDecimalsAccounting(
          item.exp_total_fun_annual_premium_per_member
        )
      ],
      [
        'Total Annual Premium',
        roundUpToTwoDecimalsAccounting(
          finalFieldValue(item, 'final_fun_annual_office_premium')
        )
      ]
    ]

    doc.autoTable({
      startY: currentY,
      body: groupFuneralBreakdownData,
      theme: 'grid',
      styles: {
        fontSize: fonts.sizes.caption,
        cellPadding: { top: 3, right: 5, bottom: 3, left: 5 },
        lineColor: colors.light,
        lineWidth: 0.3
      },
      columnStyles: {
        0: { fontStyle: 'bold', fillColor: colors.light },
        1: { halign: 'right', fontStyle: 'bold' }
      },
      didDrawPage: (data) => {
        console.log('Drawing page for category:', data)
        addPageHeaderFooter(doc, pageNumber)
        const pageHeight = doc.internal.pageSize.getHeight()
        const bottomMargin = 60 // A reasonable margin

        if (data.cursor.y > pageHeight - bottomMargin) {
          doc.addPage()
          pageNumber++
          addPageHeaderFooter(doc, pageNumber)
          currentY = topMargin // Reset Y position for the new page
        } else {
          currentY = data.cursor.y + 6 // Update currentY based on where the table ended
        }
      }
    })

    // currentY = doc.lastAutoTable.finalY + 20

    // No manual page break needed - smart detection handles this automatically
  })

  // doc.addPage('a4', 'landscape')
  // pageNumber++
  // const landscapePageWidth = doc.internal.pageSize.getWidth()
  // const landscapeContentWidth = landscapePageWidth - leftMargin * 2
  // addPageHeaderFooter(doc, pageNumber)

  if (!hasAnyNonFuneralBenefits) {
    console.log('Adding new page for Benefits and Definitions section')
    doc.addPage('a4', 'landscape')
    pageNumber++
    addPageHeaderFooter(doc, pageNumber)
    currentY = topMargin + 10
  } else {
    currentY = topMargin + 10
  }

  // Enhanced landscape page title

  doc.setTextColor(...colors.dark)
  doc.setFontSize(fonts.sizes.heading)
  doc.setFont(fonts.primary, 'bold')
  doc.text('Benefits and Definitions of the Cover', leftMargin, currentY)
  // Add underline accent
  doc.setDrawColor(...colors.primary)
  doc.setLineWidth(1)
  doc.line(leftMargin, currentY + 2, leftMargin + 120, currentY + 2)

  currentY += 15
  doc.setTextColor(...colors.dark)

  quote.value.scheme_categories.forEach((item, index) => {
    // const pageHeight = doc.internal.pageSize.getHeight()
    // Only check for page break if we're very close to the bottom (conservative approach)
    if (checkPageBreak(doc, currentY, 60)) {
      doc.addPage('a4', 'landscape')
      pageNumber++
      addPageHeaderFooter(doc, pageNumber)
      currentY = topMargin
    }

    // Category header with styling
    doc.setTextColor(...colors.dark)
    doc.setFontSize(fonts.sizes.subheading)
    doc.setFont(fonts.primary, 'bold')
    doc.text(`${item.scheme_category}`, leftMargin, currentY)
    currentY += 5

    const glaEducatorLabel =
      benefitMaps.value
        .find((b: any) => b.benefit_code === 'GLA_EDU')
        ?.benefit_alias?.trim() || 'GLA Educator'
    const ptdEducatorLabel =
      benefitMaps.value
        .find((b: any) => b.benefit_code === 'PTD_EDU')
        ?.benefit_alias?.trim() || 'PTD Educator'
    const categoryBenefitCommonData = [
      ['Terminal Illness', `${item.gla_terminal_illness_benefit}`],
      ['Free Cover Limit', quote.value.free_cover_limit],
      [glaEducatorLabel, item.gla_educator_benefit],
      [ptdEducatorLabel, item.ptd_educator_benefit],
      ['Retirement Premium Waiver', item.phi_premium_waiver || 'No'],
      [
        'Medical Aid Premium Waiver',
        item.phi_medical_aid_premium_waiver || 'No'
      ]
    ]

    if (hasAnyNonFuneralBenefits) {
      doc.autoTable({
        startY: currentY,
        body: categoryBenefitCommonData,
        theme: 'grid',
        styles: {
          fontSize: fonts.sizes.caption,
          cellPadding: { top: 3, right: 5, bottom: 3, left: 5 },
          lineColor: colors.light,
          lineWidth: 0.3
        },
        columnStyles: {
          0: { fontStyle: 'bold', fillColor: colors.light },
          1: { fontStyle: 'normal' }
        },
        didDrawPage: (data) => {
          addPageHeaderFooter(doc, pageNumber)
          const pageHeight = doc.internal.pageSize.getHeight()
          const bottomMargin = 60 // A reasonable margin

          if (data.cursor.y > pageHeight - bottomMargin) {
            doc.addPage('a4', 'landscape')
            pageNumber++
            addPageHeaderFooter(doc, pageNumber)
            currentY = topMargin // Reset Y position for the new page
          } else {
            currentY = data.cursor.y + 6 // Update currentY based on where the table ended
          }
        }
      })
      currentY = doc.lastAutoTable.finalY + 8
    }

    const categoryBenefitData = [
      [
        'Benefit',
        'Salary Multiple',
        'Benefit Structure',
        'Waiting Period',
        'Deferred Period',
        'Cover Definition',
        'Risk Type'
      ]
    ]
    categoryBenefitData.push([
      glaBenefitTitle.value,
      quote.value.use_global_salary_multiple
        ? item.gla_salary_multiple
        : 'varies',
      'standalone',
      item.gla_waiting_period,
      'n.a',
      'n.a',
      'n.a'
    ])
    categoryBenefitData.push([
      sglaBenefitTitle.value,
      quote.value.use_global_salary_multiple
        ? item.sgla_salary_multiple
        : 'varies',
      'rider',
      item.sgla_waiting_period,
      'n.a',
      'n.a',
      'n.a'
    ])
    categoryBenefitData.push([
      ptdBenefitTitle.value,
      quote.value.use_global_salary_multiple
        ? item.ptd_salary_multiple
        : 'varies',
      item.ptd_benefit_type,
      '0',
      item.ptd_deferred_period,
      item.ptd_disability_definition,
      item.ptd_risk_type
    ])
    categoryBenefitData.push([
      ciBenefitTitle.value,
      quote.value.use_global_salary_multiple
        ? item.ci_critical_illness_salary_multiple
        : 'varies',
      item.ci_benefit_structure,
      item.ci_waiting_period,
      item.ci_deferred_period,
      item.ci_benefit_definition,
      'n.a'
    ])
    categoryBenefitData.push([
      phiBenefitTitle.value,
      item.phi_income_replacement_percentage / 100,
      'n.a',
      item.phi_waiting_period,
      item.phi_deferred_period,
      item.phi_disability_definition,
      item.phi_risk_type
    ])
    categoryBenefitData.push([
      ttdBenefitTitle.value,
      item.ttd_income_replacement_percentage / 100,
      'n.a',
      item.ttd_waiting_period,
      item.ttd_deferred_period,
      item.ttd_disability_definition,
      item.ttd_risk_type
    ])

    if (hasAnyNonFuneralBenefits) {
      doc.autoTable({
        startY: currentY + 5, // Add some space before the table
        head: categoryBenefitData.slice(0, 1),
        body: categoryBenefitData.slice(1),
        theme: 'striped',
        headStyles: {
          fillColor: colors.primary,
          textColor: colors.white,
          fontSize: fonts.sizes.caption,
          fontStyle: 'bold',
          halign: 'center'
        },
        styles: {
          fontSize: 8, // Smaller font for landscape table
          cellPadding: { top: 2, right: 3, bottom: 2, left: 3 },
          lineColor: colors.light,
          lineWidth: 0.2
        },
        columnStyles: {
          0: {
            fontStyle: 'bold',
            fillColor: [248, 249, 250],
            minCellWidth: 25
          },
          1: { halign: 'center', minCellWidth: 20 },
          2: { halign: 'center', minCellWidth: 20 },
          3: { halign: 'center', minCellWidth: 20 },
          4: { halign: 'center', minCellWidth: 20 },
          5: { halign: 'center', minCellWidth: 20 },
          6: { halign: 'center', minCellWidth: 15 }
        },
        alternateRowStyles: {
          fillColor: [252, 253, 254]
        },
        didDrawPage: (data) => {
          addPageHeaderFooter(doc, pageNumber)
          const pageHeight = doc.internal.pageSize.getHeight()
          const bottomMargin = 60 // A reasonable margin

          if (data.cursor.y > pageHeight - bottomMargin) {
            doc.addPage('a4', 'landscape')
            pageNumber++
            addPageHeaderFooter(doc, pageNumber)
            currentY = topMargin // Reset Y position for the new page
          }
        }
      })
    }
    // currentY = doc.lastAutoTable.finalY + 8

    // Group Funeral section with styled header
    console.log('Adding Group Funeral section for category:', currentY)
    doc.setTextColor(...colors.dark)
    doc.setFontSize(fonts.sizes.body)
    doc.setFont(fonts.primary, 'bold')
    doc.text(`${item.scheme_category} - Group Funeral`, leftMargin, currentY)

    const groupFuneralBenefitData = [
      ['Member', 'Sum Assured', 'Maximum Number Covered']
    ]
    groupFuneralBenefitData.push(
      ['Main Member', item.family_funeral_main_member_funeral_sum_assured, 1],
      ['Spouse', item.family_funeral_spouse_funeral_sum_assured, 1],
      [
        'Child',
        item.family_funeral_children_funeral_sum_assured,
        item.family_funeral_max_number_children
      ],
      [
        'Parent',
        item.family_funeral_parent_funeral_sum_assured,
        item.family_funeral_parent_maximum_number_covered
      ],
      [
        'Dependant',
        item.family_funeral_adult_dependant_sum_assured,
        item.family_funeral_max_number_adult_dependants
      ]
    )

    doc.autoTable({
      startY: currentY + 5, // Add some space before the table
      head: groupFuneralBenefitData.slice(0, 1),
      body: groupFuneralBenefitData.slice(1),
      theme: 'grid',
      headStyles: {
        fillColor: colors.secondary,
        textColor: colors.white,
        fontSize: fonts.sizes.caption,
        fontStyle: 'bold',
        halign: 'center'
      },
      styles: {
        fontSize: fonts.sizes.caption,
        cellPadding: { top: 3, right: 5, bottom: 3, left: 5 },
        lineColor: colors.light,
        lineWidth: 0.3
      },
      columnStyles: {
        0: { fontStyle: 'bold', fillColor: colors.light },
        1: { halign: 'right' },
        2: { halign: 'center' }
      },
      didDrawPage: (data) => {
        addPageHeaderFooter(doc, pageNumber)
        const pageHeight = doc.internal.pageSize.getHeight()
        const bottomMargin = 60 // A reasonable margin

        if (data.cursor.y > pageHeight - bottomMargin) {
          doc.addPage('a4', 'landscape')
          pageNumber++
          addPageHeaderFooter(doc, pageNumber)
          currentY = topMargin // Reset Y position for the new page
        } else {
          currentY = data.cursor.y + 6 // Update currentY based on where the table ended
        }
      }
    })
    currentY = doc.lastAutoTable.finalY + 8

    // educator benefits
    if (categoryEducatorBenefits.value.length > 0 && hasAnyNonFuneralBenefits) {
      // Educator Benefits section header
      doc.setTextColor(...colors.dark)
      doc.setFontSize(fonts.sizes.body)
      doc.setFont(fonts.primary, 'bold')
      doc.text('Educator Benefits', leftMargin, currentY)
      currentY += 4

      const educatorBenefitsData = [
        [
          'Education Level',
          'Maximum Tuition per Year',
          'Maximum Coverage Period'
        ]
      ]

      const catItem = categoryEducatorBenefits.value.find(
        (educatorBenefit) =>
          educatorBenefit.scheme_category === item.scheme_category
      )

      console.log('Educator Benefits for Category:', catItem)

      educatorBenefitsData.push(
        [
          'Grade 0',
          catItem?.educator_benefit_structure.grade0_max_tuition_per_year ||
            'n.a',
          catItem?.educator_benefit_structure.grade0_max_coverage_years || 'n.a'
        ],
        [
          'Grade  1 - 7',
          catItem?.educator_benefit_structure.grade17_max_tuition_per_year ||
            'n.a',
          catItem?.educator_benefit_structure.grade17_max_coverage_years ||
            'n.a'
        ],
        [
          'Grade  8 - 12',
          catItem?.educator_benefit_structure.grade812_max_tuition_per_year ||
            'n.a',
          catItem?.educator_benefit_structure.grade812_max_coverage_years ||
            'n.a'
        ],
        [
          'Tertiary Education',
          catItem?.educator_benefit_structure.tertiary_max_tuition_per_year ||
            'n.a',
          catItem?.educator_benefit_structure.tertiary_max_coverage_years ||
            'n.a'
        ]
      )

      doc.autoTable({
        startY: currentY,
        head: educatorBenefitsData.slice(0, 1),
        body: educatorBenefitsData.slice(1),
        theme: 'grid',
        headStyles: {
          fillColor: colors.primary,
          textColor: colors.white,
          fontSize: fonts.sizes.caption,
          fontStyle: 'bold',
          halign: 'center'
        },
        styles: {
          fontSize: fonts.sizes.caption,
          cellPadding: { top: 3, right: 5, bottom: 3, left: 5 },
          lineColor: colors.light,
          lineWidth: 0.3
        },
        columnStyles: {
          0: { fontStyle: 'bold', fillColor: colors.light },
          1: { halign: 'right' },
          2: { halign: 'center' }
        },
        didDrawPage: (data) => {
          addPageHeaderFooter(doc, pageNumber)
          const pageHeight = doc.internal.pageSize.getHeight()
          const bottomMargin = 60 // A reasonable margin

          if (data.cursor.y > pageHeight - bottomMargin) {
            doc.addPage()
            pageNumber++
            addPageHeaderFooter(doc, pageNumber)
            currentY = topMargin // Reset Y position for the new page
          } else {
            currentY = data.cursor.y + 6 // Update currentY based on where the table ended
          }
        }
      })
      if (hasAnyNonFuneralBenefits) {
        currentY = doc.lastAutoTable.finalY + 30
      }
    } else {
      currentY += 30
    }

    // Smart page break detection handles this automatically
  })

  // **Enhanced Final Page for Underwriting and General Provisions**

  doc.addPage('a4', 'portrait')
  currentY = topMargin
  pageNumber++
  addPageHeaderFooter(doc, pageNumber)

  doc.setTextColor(...colors.dark)
  doc.setFontSize(fonts.sizes.heading)
  doc.setFont(fonts.primary, 'bold')
  doc.text('Underwriting and General Provisions', leftMargin, currentY)

  currentY += 12
  doc.setTextColor(...colors.dark)
  doc.setFontSize(fonts.sizes.body)
  doc.setFont(fonts.primary, 'normal')

  // Enhanced text wrapping with better spacing
  wrapText(
    doc,
    insurer.value.general_provisions_text,
    leftMargin,
    currentY,
    contentWidth
  )

  // Add professional closing section
  const finalY = doc.internal.pageSize.getHeight() - 60
  doc.setFillColor(...colors.light)
  doc.rect(leftMargin - 5, finalY - 5, contentWidth + 10, 40, 'F')

  doc.setTextColor(...colors.secondary)
  doc.setFontSize(fonts.sizes.body)
  doc.setFont(fonts.primary, 'italic')
  doc.text('Thank you for considering our quotation.', leftMargin, finalY + 5)
  doc.text(
    'We look forward to the opportunity to serve your insurance needs.',
    leftMargin,
    finalY + 15
  )

  doc.setFont(fonts.primary, 'bold')
  doc.text(`${insurer.value.name}`, leftMargin, finalY + 25)

  // **Further pages for  the acceptance form**
  doc.addPage('a4', 'portrait')
  currentY = topMargin
  pageNumber++
  addPageHeaderFooter(doc, pageNumber)

  const margin = 15

  // Colors
  const navyBlue = [30, 58, 95] // #1e3a5f
  const darkGray = [26, 26, 26]
  const mediumGray = [88, 96, 105]
  const lightGray = [225, 228, 232]
  const orange = [212, 118, 0] // #d47600
  const lightOrange = [255, 248, 240]
  const lightBlue = [241, 248, 255]

  // Helper functions
  function drawRect(x, y, w, h, fillColor, borderColor, borderWidth = 0.2) {
    if (fillColor) {
      doc.setFillColor(fillColor[0], fillColor[1], fillColor[2])
      doc.rect(x, y, w, h, 'F')
    }
    if (borderColor) {
      doc.setDrawColor(borderColor[0], borderColor[1], borderColor[2])
      doc.setLineWidth(borderWidth)
      doc.rect(x, y, w, h, 'S')
    }
  }

  function drawLine(x1, y1, x2, y2, color = mediumGray, width = 0.3) {
    doc.setDrawColor(color[0], color[1], color[2])
    doc.setLineWidth(width)
    doc.line(x1, y1, x2, y2)
  }

  function setFont(size, style = 'normal', color = darkGray) {
    doc.setFont('helvetica', style)
    doc.setFontSize(size)
    doc.setTextColor(color[0], color[1], color[2])
  }

  // ===== HEADER SECTION =====
  const headerHeight = 25
  drawRect(0, 0, pageWidth, headerHeight, navyBlue, null)

  // Title
  setFont(16, 'bold', [255, 255, 255])
  doc.text('ACCEPTANCE OF QUOTATION', margin, 12)

  // Subtitle
  setFont(8, 'normal', [230, 230, 230])
  doc.text('POPIA Compliant', margin, 18)

  // Logo placeholder
  doc.setFillColor(255, 255, 255)
  doc.setDrawColor(255, 255, 255)
  doc.setLineWidth(0.3)
  doc.setLineDash([1, 1], 0)
  doc.roundedRect(pageWidth - margin - 25, 5, 20, 12, 2, 2, 'D')
  setFont(7, 'normal', [255, 255, 255])
  doc.text('[LOGO]', pageWidth - margin - 22, 12)
  doc.setLineDash([], 0)

  currentY = headerHeight + 8

  // ===== POLICY DETAILS SECTION =====
  const sectionHeight = 29
  drawRect(
    margin,
    currentY,
    contentWidth,
    sectionHeight,
    [250, 251, 252],
    lightGray
  )

  // Section title
  setFont(7, 'bold', navyBlue)
  doc.text('POLICY DETAILS', margin + 3, currentY + 5)
  drawLine(margin + 3, currentY + 7, margin + 35, currentY + 7, navyBlue, 1)

  // Labels and input lines
  setFont(6, 'normal', mediumGray)

  // Row 1
  doc.text('EMPLOYER / SCHEME NAME', margin + 3, currentY + 13)
  doc.text('QUOTE NUMBER', margin + 105, currentY + 13)

  drawLine(margin + 3, currentY + 19, margin + 70, currentY + 19)
  drawLine(margin + 105, currentY + 19, margin + 160, currentY + 19)

  // Row 2
  doc.text('DATE OF QUOTE', margin + 3, currentY + 24)
  doc.text('COMMENCEMENT DATE', margin + 105, currentY + 24)

  drawLine(margin + 3, currentY + 28, margin + 70, currentY + 28)
  drawLine(margin + 105, currentY + 28, margin + 160, currentY + 28)

  // Placeholder text
  setFont(7, 'normal', [180, 180, 180])
  doc.text('DD/MM/YYYY', margin + 3, currentY + 26)
  doc.text('DD/MM/YYYY', margin + 105, currentY + 26)

  currentY += sectionHeight + 38

  setFont(10, 'normal', darkGray)
  const profileText =
    'If the member data profile at the quotation implementation date differ by 7% or more from that on which the quotation was based, we reserve the right to revise the rates and Automatic Acceptance Limit. The Employer/Scheme will be notified accordingly and must provide acceptance before implementation proceeds.'
  const splitprofile = doc.splitTextToSize(profileText, contentWidth - 2)
  const textOptions = { align: 'justify', maxWidth: contentWidth - 2 }
  doc.text(splitprofile, margin, currentY - 26, textOptions)

  setFont(10, 'normal', darkGray)
  const acknowledgementText =
    'By signing this quotation, the Employer/Scheme acknowledges that they have read, understood, and agree to be bound by all the terms and conditions of this quotation.'
  const splitacknowledgement = doc.splitTextToSize(
    acknowledgementText,
    contentWidth - 2
  )
  doc.text(splitacknowledgement, margin, currentY - 7, textOptions)

  // ===== EMPLOYER AUTHORISATION SECTION =====
  const empSectionHeight = 41
  drawRect(
    margin,
    currentY,
    contentWidth,
    empSectionHeight,
    [250, 251, 252],
    lightGray
  )

  setFont(7, 'bold', navyBlue)
  doc.text('EMPLOYER AUTHORISATION', margin + 3, currentY + 5)
  drawLine(margin + 3, currentY + 7, margin + 50, currentY + 7, navyBlue, 1)

  // Signature box
  const sigBoxHeight = 22
  drawRect(
    margin + 3,
    currentY + 10,
    contentWidth - 6,
    sigBoxHeight,
    [255, 255, 255],
    mediumGray
  )

  setFont(7, 'bold', navyBlue)
  doc.text('Duly Authorised Signatory', margin + 5, currentY + 15)

  // Signature line area
  setFont(6, 'normal', [150, 150, 150])
  doc.text('Sign here', margin + 5, currentY + 25)
  drawLine(
    margin + 5,
    currentY + 28,
    margin + 80,
    currentY + 28,
    [200, 200, 200],
    0.5
  )

  // Clear signature text
  setFont(6, 'normal', [200, 100, 100])
  doc.text('Clear Signature', margin + 5, currentY + 30)

  // Input fields below signature
  const inputY = currentY + 38
  setFont(6, 'normal', mediumGray)

  doc.text('FULL NAME', margin + 3, inputY)
  doc.text('CAPACITY', margin + 85, inputY)
  doc.text('DATE', margin + 130, inputY)

  drawLine(margin + 3, inputY + 2, margin + 80, inputY + 2)
  drawLine(margin + 85, inputY + 2, margin + 125, inputY + 2)
  drawLine(margin + 130, inputY + 2, margin + 170, inputY + 2)

  // Placeholders
  setFont(6, 'normal', [180, 180, 180])
  doc.text('e.g., Director', margin + 85, inputY + 2)
  doc.text('DD/MM/YYYY', margin + 130, inputY + 2)

  currentY += empSectionHeight + 6

  // ===== INTERMEDIARY DETAILS SECTION =====
  const intSectionHeight = 41
  drawRect(
    margin,
    currentY,
    contentWidth,
    intSectionHeight,
    [250, 251, 252],
    lightGray
  )

  setFont(7, 'bold', navyBlue)
  doc.text('INTERMEDIARY DETAILS', margin + 3, currentY + 5)
  drawLine(margin + 3, currentY + 7, margin + 45, currentY + 7, navyBlue, 1)

  // Signature box
  drawRect(
    margin + 3,
    currentY + 10,
    contentWidth - 6,
    sigBoxHeight,
    [255, 255, 255],
    mediumGray
  )

  setFont(7, 'bold', navyBlue)
  doc.text('Intermediary / FAIS Representative', margin + 5, currentY + 15)

  // Signature line area
  setFont(6, 'normal', [150, 150, 150])
  doc.text('Sign here', margin + 5, currentY + 25)
  drawLine(
    margin + 5,
    currentY + 28,
    margin + 80,
    currentY + 28,
    [200, 200, 200],
    0.5
  )

  // Clear signature text
  setFont(6, 'normal', [200, 100, 100])
  doc.text('Clear Signature', margin + 5, currentY + 30)

  // Input fields
  const intInputY = currentY + 38
  setFont(6, 'normal', mediumGray)

  doc.text('FULL NAME', margin + 3, intInputY)
  doc.text('FAIS REG NO.', margin + 85, intInputY)
  doc.text('DATE', margin + 130, intInputY)

  drawLine(margin + 3, intInputY + 2, margin + 80, intInputY + 2)
  drawLine(margin + 85, intInputY + 2, margin + 125, intInputY + 2)
  drawLine(margin + 130, intInputY + 2, margin + 170, intInputY + 2)

  // Placeholder
  setFont(6, 'normal', [180, 180, 180])
  doc.text('DD/MM/YYYY', margin + 130, intInputY + 2)

  currentY += intSectionHeight + 6

  // ===== POPIA SECTION =====
  const popiaHeight = 32
  drawRect(
    margin,
    currentY,
    contentWidth,
    popiaHeight,
    lightBlue,
    [200, 225, 255]
  )

  setFont(7, 'bold', [3, 102, 214]) // Blue color
  doc.text('POPIA Consent & Data Protection', margin + 3, currentY + 5)

  // POPIA text
  setFont(7, 'normal', darkGray)
  const popiaText = `In terms of the Protection of Personal Information Act 4 of 2013 (POPIA), the Employer consents to the processing of personal information of employees and scheme members for the purpose of underwriting, administering, and processing claims under this Group Risk policy. Information will be processed lawfully, minimally, and only for the specific purpose stated.`

  const splitPopia = doc.splitTextToSize(popiaText, contentWidth - 6)
  doc.text(splitPopia, margin + 3, currentY + 11)

  // Checkbox
  const checkboxY = currentY + 22
  drawRect(margin + 3, checkboxY - 3, 4, 4, [255, 255, 255], mediumGray, 0.3)

  setFont(6.5, 'normal', darkGray)
  const consentText = `I confirm that the Employer has obtained necessary consent from data subjects (employees/members) for the processing of their personal information as required by POPIA, and warrants that all information provided is true and complete.`

  const splitConsent = doc.splitTextToSize(consentText, contentWidth - 12)
  doc.text(splitConsent, margin + 9, checkboxY)

  currentY += popiaHeight + 6

  // ===== FOR OFFICE USE ONLY SECTION =====
  const officeHeight = 25
  drawRect(
    margin,
    currentY,
    contentWidth,
    officeHeight,
    lightOrange,
    [255, 209, 168]
  )

  // Lock icon and title
  // setFont(9, 'normal', orange)
  // doc.text('🔒', margin + 3, currentY + 6)

  setFont(7, 'bold', orange)
  doc.text('FOR OFFICE USE ONLY', margin + 3, currentY + 6)

  // Table grid
  const tableY = currentY + 10
  const colWidth = contentWidth / 3

  // Horizontal lines
  drawLine(margin, tableY, margin + contentWidth, tableY, [255, 180, 120], 0.3)
  drawLine(
    margin,
    tableY + 7,
    margin + contentWidth,
    tableY + 7,
    [255, 180, 120],
    0.3
  )

  // Vertical lines
  for (let i = 0; i <= 3; i++) {
    drawLine(
      margin + colWidth * i,
      tableY,
      margin + colWidth * i,
      tableY + 14,
      [255, 180, 120],
      0.3
    )
  }

  // Headers
  setFont(6, 'normal', mediumGray)
  doc.text('RECEIVED BY', margin + 3, tableY + 4)
  doc.text('DATE RECEIVED', margin + colWidth + 3, tableY + 4)
  doc.text('POLICY NUMBER', margin + colWidth * 2 + 3, tableY + 4)

  doc.text('UNDERWRITER', margin + 3, tableY + 11)
  doc.text('APPROVED BY', margin + colWidth + 3, tableY + 11)
  doc.text('DATE APPROVED', margin + colWidth * 2 + 3, tableY + 11)

  // // ===== FOOTER =====
  // const footerY = pageHeight - 12;
  // drawRect(0, footerY - 5, pageWidth, 17, [246, 248, 250], [225, 228, 232]);

  // setFont(6.5, 'normal', mediumGray);
  // doc.text('Doc Ref: GRA-SA-2024', margin, footerY);
  // doc.text('Page 1 of 1', margin + 40, footerY);

  // setFont(6.5, 'normal', mediumGray);
  // const footerText = 'This acceptance is subject to the terms and conditions of the Group Risk Policy';
  // const textWidth = doc.getTextWidth(footerText);
  // doc.text(footerText, pageWidth - margin - textWidth, footerY);

  // End of Acceptance Form - can add more content here as needed

  const pdfBlob = doc.output('blob')
  pdfUrl.value = URL.createObjectURL(pdfBlob)
}
</script>

<style scoped>
/*
button {
  padding: 10px 20px;
  font-size: 16px;
  cursor: pointer;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 5px;
  margin: 20px;
}

button:hover {
  background-color: #38a873;
}
  */

.pdf-viewer {
  margin-top: 20px;
  border: 1px solid #ccc;
}
.pdf-render-content {
  font-family: Helvetica, sans-serif;
  font-size: 10pt !important; /* Use !important to override other styles */
  line-height: 1.5;
}

/* You can also target specific tags inside it */
.pdf-render-content :deep(*) {
  font-size: 10pt !important;
}

.pdf-render-content h1 {
  font-size: 14pt !important;
}

/* AG Grid custom styles */
.ag-theme-balham {
  --ag-header-background-color: #1976d2;
  --ag-header-foreground-color: white;
  --ag-odd-row-background-color: #f8f9fa;
  --ag-row-hover-color: #e3f2fd;
}

.ag-theme-balham .ag-header-cell-label {
  font-weight: bold;
}

.ag-theme-balham .ag-row-group {
  font-weight: bold;
  background-color: #e8f5e8 !important;
}

.ag-theme-balham .ag-group-expanded .ag-icon {
  color: #1976d2;
}

.ag-theme-balham .ag-group-contracted .ag-icon {
  color: #1976d2;
}
</style>
