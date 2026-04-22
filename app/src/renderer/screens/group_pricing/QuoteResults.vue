<template>
  <base-card :show-actions="false">
    <template #header>
      <span class="headline">Calculation Results & Analysis</span>
    </template>
    <template #default>
      <v-list lines="two">
        <v-list-item v-for="item in relatedResultTables" :key="item.table_type">
          <template #prepend>
            <v-icon :color="item.populated ? 'success' : 'error'" class="mr-4">
              {{
                item.populated ? 'mdi-check-circle' : 'mdi-alert-circle-outline'
              }}
            </v-icon>
          </template>
          <v-list-item-title>{{ item.table_type }}</v-list-item-title>
          <v-list-item-subtitle
            v-if="item.table_type !== 'Group Pricing Parameters'"
            :class="item.populated ? '' : 'text-error'"
          >
            {{
              item.populated ? `${item.count} records loaded` : 'Data required'
            }}
          </v-list-item-subtitle>

          <template #append>
            <v-btn
              rounded
              color="primary"
              class="mr-2"
              variant="outlined"
              size="small"
              :disabled="!item.populated"
              @click="viewTable(item)"
            >
              <v-icon left color="primary">mdi-information</v-icon>
              <span>View</span></v-btn
            >
          </template>
        </v-list-item>
      </v-list>

      <loading-indicator :loadingData="loadingData" />

      <v-row
        v-if="
          (resultTableData.length > 0 || useInfiniteModel) &&
          !loadingData &&
          selectedTable
        "
      >
        <v-col>
          <v-btn
            v-if="$props.quote.experience_rating === 'Yes'"
            class="mb-2"
            variant="outlined"
            rounded
            color="primary"
            @click="credibilityDialog = true"
            >Credibility Chart</v-btn
          >
          <v-btn
            v-if="$props.quote.experience_rating === 'Yes'"
            class="mb-2 ml-2"
            variant="outlined"
            rounded
            color="secondary"
            @click="manualCredibilityDialog = true"
            >Manual Credibility</v-btn
          >
          <group-pricing-data-grid
            ref="dataGridRef"
            :key="`results-grid-${gridRemountKey}`"
            :columnDefs="columnDefs"
            :show-close-button="true"
            :rowData="useInfiniteModel ? undefined : resultTableData"
            :table-title="selectedTable"
            :pagination="!useInfiniteModel"
            :rowCount="rowCount"
            :useInfiniteModel="useInfiniteModel"
            :dataSource="dataSource"
            :show-export-all-csv="shouldExportAll"
            @update:clear-data="clearData"
            @update-column="handleColumnUpdate"
            @export-all-csv="handleExportAllCsv"
          />
        </v-col>
      </v-row>

      <historical-credibility-chart-dialog
        v-model="credibilityDialog"
        :quote-id="props.quote?.id"
      />

      <!-- Manual Credibility Dialog -->
      <v-dialog v-model="manualCredibilityDialog" max-width="500px" persistent>
        <v-card>
          <v-card-title>
            <span class="headline">Manual Credibility Calculation</span>
          </v-card-title>
          <v-card-text>
            <v-form
              ref="manualCredibilityForm"
              v-model="manualCredibilityFormValid"
            >
              <v-text-field
                v-model="manualCredibilityValue"
                class="mb-7"
                variant="outlined"
                density="compact"
                label="Credibility Value"
                type="number"
                step="0.01"
                min="0"
                max="1"
                :rules="credibilityRules"
                required
                hint="Enter a value between 0 and 1"
                persistent-hint
              ></v-text-field>

              <v-select
                v-model="selectedBasis"
                variant="outlined"
                density="compact"
                :items="parameterBases"
                item-title="name"
                item-value="value"
                label="Calculation Basis"
                :rules="basisRules"
                required
              ></v-select>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              color="grey-darken-1"
              variant="text"
              @click="closeManualCredibilityDialog"
            >
              Cancel
            </v-btn>
            <v-btn
              color="primary"
              variant="text"
              :loading="calculatingManualCredibility"
              :disabled="!manualCredibilityFormValid"
              @click="calculateManualCredibility"
            >
              Calculate
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </template>
  </base-card>
  <!-- Calculation progress overlay -->
  <v-overlay
    :model-value="awaitingManualCredibility"
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
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useCalculationProgress } from '@/renderer/composables/useCalculationProgress'
import _ from 'lodash'
import formatValues from '@/renderer/utils/format_values'
import HistoricalCredibilityChartDialog from '@/renderer/components/charts/HistoricalCredibilityChartDialog.vue'
import LoadingIndicator from '@/renderer/components/LoadingIndicator.vue'
import GroupPricingDataGrid from '@/renderer/components/tables/GroupPricingDataGrid.vue'

// import all necessary components and services

const props = defineProps({
  quote: { type: Object, required: true }
})

const emit = defineEmits(['quote-updated'])
const rowCount: any = ref(0)
const shouldExportAll = ref(true)

const columnDefs: any = ref([])
const selectedTable = ref('')
const loadingData = ref(false)
const credibilityDialog = ref(false)
const resultTableData: any = ref([])
const displaySummary = ref(false)
const snackbar = ref(false)
const snackbarTimeout = 3000
const snackbarText = ref('No data found for this table')
const benefitMaps: any = ref([])
const glaBenefitTitle = ref('')
const sglaBenefitTitle = ref('')
const ptdBenefitTitle = ref('')
const ciBenefitTitle = ref('')
const phiBenefitTitle = ref('')
const ttdBenefitTitle = ref('')
const familyFuneralBenefitTitle = ref('')
const useInfiniteModel = ref(false)
const dataSource: any = ref(null)
const currentTableType = ref('')
const dataGridRef: any = ref(null)
// Bumped when the underlying data for the currently-viewed table changes
// (e.g. the user re-runs calculations) so the grid remounts and drops its
// cached blocks. AG Grid's server-side / infinite row model caches rows by
// block index and will keep serving stale data after a re-run otherwise —
// the user would see the pre-run numbers even though the DB has the fresh
// ones. The currently selected table is preserved because only the inner
// grid child is remounted, not QuoteResults itself.
const gridRemountKey = ref(0)

// Manual Credibility Dialog state
const manualCredibilityDialog = ref(false)
const manualCredibilityValue = ref('')
const selectedBasis: any = ref('')
const parameterBases = ref([])
const manualCredibilityFormValid = ref(false)
const calculatingManualCredibility = ref(false)

const {
  progress: calcProgress,
  phaseLabel,
  progressPercent,
  startTracking,
  stopTracking
} = useCalculationProgress()

// Only react to progress events when this component initiated the calculation.
const awaitingManualCredibility = ref(false)

watch(calcProgress, (val) => {
  if (!awaitingManualCredibility.value) return
  if (val?.phase === 'completed') {
    snackbarText.value = 'Manual credibility calculation completed successfully'
    snackbar.value = true
    awaitingManualCredibility.value = false
    emit('quote-updated')
  }
  if (val?.phase === 'failed') {
    snackbarText.value =
      'Calculations failed. Please contact your administrator.'
    snackbar.value = true
    awaitingManualCredibility.value = false
  }
})

// When the parent reloads the quote after a re-run the *_count fields on
// the quote change (new member rating rows, fresh bordereaux, etc.). Bump
// gridRemountKey so the inner AG Grid remounts and drops its cached blocks;
// otherwise a user who was already looking at a table keeps seeing the
// pre-run numbers served from the infinite-model cache.
watch(
  [
    () => props.quote?.member_rating_result_count,
    () => props.quote?.member_data_count,
    () => props.quote?.claims_experience_count,
    () => props.quote?.member_premium_schedule_count,
    () => props.quote?.bordereaux_count
  ],
  () => {
    gridRemountKey.value++
  }
)

const manualCredibilityForm: any = ref(null)

// Validation rules
const credibilityRules = [
  (v: any) => !!v || 'Credibility value is required',
  (v: any) =>
    (!isNaN(parseFloat(v)) && parseFloat(v) >= 0 && parseFloat(v) <= 1) ||
    'Credibility must be between 0 and 1'
]

const basisRules = [(v: any) => !!v || 'Calculation basis is required']
// Define related result tables and their properties
const tableData = ref([]) // This will hold the data for the selected table

// ... other state and methods related to results

const handleColumnUpdate = async ({ colId, newValue }) => {
  // your existing logic...
  // on success, emit event to refresh parent
  console.log('Column updated:', colId, newValue)
  emit('quote-updated')
}

const clearData = () => {
  tableData.value = []
  resultTableData.value = []
  selectedTable.value = ''
  useInfiniteModel.value = false
  dataSource.value = null
  currentTableType.value = ''
}

const handleExportAllCsv = async () => {
  try {
    console.log('Exporting all CSV data for table:', currentTableType.value)

    const response = await GroupPricingService.exportQuoteTableCsv(
      props.quote.id,
      currentTableType.value
    )

    // Create download link
    const fileURL = window.URL.createObjectURL(new Blob([response.data]))
    const fileLink = document.createElement('a')

    fileLink.href = fileURL
    fileLink.setAttribute(
      'download',
      `${currentTableType.value}_${props.quote.id}.xlsx`
    )
    document.body.appendChild(fileLink)

    fileLink.click()

    // Clean up
    document.body.removeChild(fileLink)
    window.URL.revokeObjectURL(fileURL)
  } catch (error) {
    console.error('Error exporting CSV:', error)
    snackbarText.value = 'Error exporting CSV file'
    snackbar.value = true
  } finally {
    // Reset the loader in the data grid component
    if (dataGridRef.value && dataGridRef.value.resetExportAllCsvLoader) {
      dataGridRef.value.resetExportAllCsvLoader()
    }
  }
}

const closeManualCredibilityDialog = () => {
  manualCredibilityDialog.value = false
  manualCredibilityValue.value = ''
  selectedBasis.value = ''
  if (manualCredibilityForm.value) {
    manualCredibilityForm.value.reset()
  }
}

const calculateManualCredibility = async () => {
  if (!manualCredibilityFormValid.value) return

  calculatingManualCredibility.value = true
  awaitingManualCredibility.value = true
  startTracking(String(props.quote.id))

  try {
    await GroupPricingService.runQuoteCalculationsWithCredibility(
      props.quote.id,
      selectedBasis.value,
      parseFloat(manualCredibilityValue.value)
    )

    // Job has been queued — progress and completion arrive via WebSocket.
    closeManualCredibilityDialog()
  } catch (error: any) {
    console.error('Error calculating manual credibility:', error)
    stopTracking()
    awaitingManualCredibility.value = false
    snackbarText.value =
      error.response?.data?.message || 'Error calculating manual credibility'
    snackbar.value = true
  } finally {
    calculatingManualCredibility.value = false
  }
}

onMounted(async () => {
  // Fetch result summary or other initial result data
  console.log('Quote ID:', props.quote)

  // Load parameter bases for manual credibility calculation
  try {
    const basesResponse = await GroupPricingService.getParameterBases()
    parameterBases.value = basesResponse.data || []
  } catch (error) {
    console.error('Error loading parameter bases:', error)
  }

  GroupPricingService.getBenefitMaps().then((res) => {
    benefitMaps.value = res.data
    const glaBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'GLA'
    )
    if (glaBenefit.benefit_alias !== '') {
      glaBenefitTitle.value = glaBenefit.benefit_alias
    }
    const sglaBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'SGLA'
    )
    if (sglaBenefit.benefit_alias !== '') {
      sglaBenefitTitle.value = sglaBenefit.benefit_alias
    }
    const ptdBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'PTD'
    )
    if (ptdBenefit.benefit_alias !== '') {
      ptdBenefitTitle.value = ptdBenefit.benefit_alias
    }
    const ciBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'CI'
    )
    if (ciBenefit.benefit_alias !== '') {
      ciBenefitTitle.value = ciBenefit.benefit_alias
    }
    const phiBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'PHI'
    )
    if (phiBenefit.benefit_alias !== '') {
      phiBenefitTitle.value = phiBenefit.benefit_alias
    }
    const ttdBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'TTD'
    )
    if (ttdBenefit.benefit_alias !== '') {
      ttdBenefitTitle.value = ttdBenefit.benefit_alias
    }
    const familyFuneralBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'GFF'
    )
    if (familyFuneralBenefit.benefit_alias !== '') {
      familyFuneralBenefitTitle.value = familyFuneralBenefit.benefit_alias
    }
  })
})
const relatedResultTables = computed(() => {
  const tables: any = []

  if (props.quote.member_rating_result_count > 0) {
    tables.push({
      table_type: 'Member Rating Results',
      value: 'member_rating_results',
      populated: true,
      count: props.quote.member_rating_result_count
    })
  } else {
    tables.push({
      table_type: 'Member Rating Results',
      value: 'member_rating_results',
      populated: false,
      count: 0
    })
  }

  if (props.quote.member_premium_schedule_count > 0) {
    tables.push({
      table_type: 'Member Premium Schedules',
      value: 'member_premium_schedules',
      populated: true,
      count: props.quote.member_premium_schedule_count
    })
  } else {
    tables.push({
      table_type: 'Member Premium Schedules',
      value: 'member_premium_schedules',
      populated: false,
      count: 0
    })
  }
  if (props.quote.bordereaux_count > 0) {
    tables.push({
      table_type: 'Bordereaux',
      value: 'bordereaux',
      populated: true,
      count: props.quote.bordereaux_count
    })
  } else {
    tables.push({
      table_type: 'Bordereaux',
      value: 'bordereaux',
      populated: false,
      count: 0
    })
  }

  return tables
})

const viewTable = async (item: any) => {
  // special case for output summary
  loadingData.value = true

  try {
    currentTableType.value = item.value
    console.log('Item to view:', item)

    // Check if we should use infinite model (for large datasets)
    const shouldUseInfiniteModel = item.count > 20000
    useInfiniteModel.value = shouldUseInfiniteModel

    if (shouldUseInfiniteModel) {
      // Create datasource for infinite row model
      console.log('Using infinite model for', item.table_type)

      // First, fetch a small sample to get column definitions
      const sampleRes = await GroupPricingService.getQuoteTable(
        props.quote.id,
        item.value,
        {
          offset: 0,
          limit: 1
        }
      )

      if (sampleRes.data.data && sampleRes.data.data.length > 0) {
        const sampleData = sampleRes.data.data.map((item) =>
          _.fromPairs(sampleRes.data.json_tags.map((key) => [key, item[key]]))
        )
        createColumnDefs(sampleData)
      }

      dataSource.value = {
        getRows: async (params: any) => {
          console.log('Fetching rows for infinite model:', params)
          try {
            const startRow = params.request.startRow
            const endRow = params.request.endRow
            const limit = endRow - startRow

            console.log('Fetching rows from', startRow, 'to', endRow)

            const res = await GroupPricingService.getQuoteTable(
              props.quote.id,
              currentTableType.value,
              { offset: startRow, limit }
            )
            rowCount.value = endRow
            if (res.data.data !== null && res.data.data.length > 0) {
              const orderedData = res.data.data.map((item) =>
                _.fromPairs(res.data.json_tags.map((key) => [key, item[key]]))
              )

              // Calculate last row - if we got less data than requested, we've reached the end
              const lastRow =
                res.data.data.length < limit
                  ? startRow + res.data.data.length
                  : -1
              params.success({ rowData: orderedData, lastRow })
            } else {
              params.success([], 0)
            }
          } catch (error) {
            console.error('Error fetching data:', error)
            params.fail()
          }
        }
      }
      console.log(dataSource.value)
      selectedTable.value = item.table_type
      displaySummary.value = false
      resultTableData.value = [] // Clear client-side data when using infinite model
      loadingData.value = false
    } else {
      // Use original client-side approach for smaller datasets
      const res = await GroupPricingService.getQuoteTable(
        props.quote.id,
        item.value
      )
      if (res.data.data !== null && res.data.data.length > 0) {
        const orderedData = res.data.data.map((item) =>
          _.fromPairs(res.data.json_tags.map((key) => [key, item[key]]))
        )

        displaySummary.value = false

        if (
          item.value === 'member_rating_results' ||
          item.value === 'member_premium_schedules' ||
          item.value === 'bordereaux'
        ) {
          resultTableData.value = orderedData
        } else {
          tableData.value = orderedData
        }

        selectedTable.value = item.table_type
        createColumnDefs(orderedData)
        loadingData.value = false
      } else {
        tableData.value = []
        resultTableData.value = []
        selectedTable.value = ''
        snackbarText.value = 'No data found for this table'
        snackbar.value = true
        loadingData.value = false
      }
    }
  } catch (error) {
    console.log('Error:', error)
    loadingData.value = false
  }
}

const replaceSubstring = (input, search) => {
  if (input.includes(search)) {
    benefitMaps.value.forEach((item) => {
      if (
        item.benefit_code.toLowerCase() === search &&
        item.benefit_alias_code !== ''
      ) {
        input = input.replace(search, item.benefit_alias_code.toLowerCase())
      }
    })
    // return input
  }
  return input
}

const createColumnDefs = (data: any) => {
  columnDefs.value = []
  Object.keys(data[0]).forEach((element) => {
    const column: any = {}
    if (element.includes('ci')) {
      column.headerName = replaceSubstring(element, 'ci')
    } else if (element.includes('sgla')) {
      column.headerName = replaceSubstring(element, 'sgla')
    } else if (element.includes('phi')) {
      column.headerName = replaceSubstring(element, 'phi')
    } else if (element.includes('ptd')) {
      column.headerName = replaceSubstring(element, 'ptd')
    } else if (element.includes('ttd')) {
      column.headerName = replaceSubstring(element, 'ttd')
    } else if (element.includes('gff')) {
      column.headerName = replaceSubstring(element, 'gff')
    } else if (element.includes('gla')) {
      column.headerName = replaceSubstring(element, 'gla')
    } else {
      column.headerName = element
    }

    // header.headerName = element
    column.field = element
    column.valueFormatter = formatValues
    column.minWidth = 200
    column.sortable = true
    column.filter = true
    column.resizable = true
    if (column.field === 'id' || column.field === 'quote_id') {
      column.hide = true
    }
    // if column.field is exp_credibility, then set it to editable
    if (column.field === 'manually_added_credibility') {
      column.editable = true
    } else {
      column.editable = false
    }
    console.log('Column:', column)
    columnDefs.value.push(column)
  })
}
</script>
