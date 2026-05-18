<template>
  <base-card :show-actions="false">
    <template #header>
      <span class="headline">Calculation Results & Analysis</span>
    </template>
    <template #default>
      <v-expansion-panels
        v-if="zeroRateBenefitWarnings.length > 0"
        class="mb-3"
      >
        <v-expansion-panel>
          <v-expansion-panel-title
            class="warnings-banner-title font-weight-medium"
          >
            <v-icon class="mr-2" color="warning">mdi-alert</v-icon>
            <v-chip size="small" color="warning" variant="flat" class="mr-3">
              {{ zeroRateBenefitWarnings.length }}
            </v-chip>
            <span>
              Benefit warning{{
                zeroRateBenefitWarnings.length === 1 ? '' : 's'
              }}
              — chosen benefit has a zero risk rate
            </span>
            <v-spacer />
            <span class="text-caption text-medium-emphasis mr-2">
              Click to expand
            </span>
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <ul class="mb-0 ps-4">
              <li v-for="(msg, i) in zeroRateWarningsCapped" :key="'zr-' + i">
                {{ msg }}
              </li>
            </ul>
            <div v-if="zeroRateWarningsOverflow > 0" class="text-caption mt-1">
              +{{ zeroRateWarningsOverflow }} more — fix the ones above and
              recalculate.
            </div>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
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
            <span
              v-if="item.populated && calculationCompletedLabel"
              class="ml-2 text-medium-emphasis"
            >
              · {{ calculationCompletedLabel }}
            </span>
          </v-list-item-subtitle>

          <template #append>
            <v-btn
              rounded
              color="primary"
              class="mr-2"
              variant="outlined"
              size="small"
              :disabled="!item.populated || loadingItem === item.value"
              :loading="loadingItem === item.value"
              @click="viewTable(item)"
            >
              <v-icon left color="primary">mdi-information</v-icon>
              <span>View</span></v-btn
            >
          </template>
        </v-list-item>
      </v-list>

      <loading-indicator :loadingData="loadingData" :label="loadingLabel" />

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
          <v-switch
            v-if="supportsColumnToggle"
            v-model="showAllColumns"
            density="compact"
            color="primary"
            class="d-inline-flex ml-3 align-middle"
            hide-details
            :label="
              showAllColumns ? 'Showing all columns' : 'Showing summary columns'
            "
            :disabled="loadingData"
            @update:model-value="reloadCurrentTable"
          />
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
  quote: { type: Object, required: true },
  resultSummaries: { type: Array, required: false, default: () => [] }
})

const zeroRateBenefitWarnings = computed<string[]>(() => {
  const cats = (props.quote as any)?.scheme_categories ?? []
  const findCat = (name: string) =>
    cats.find((c: any) => c.scheme_category === name)

  const checks: Array<{ flag: string; rate: string; label: string }> = [
    { flag: 'gla_benefit', rate: 'total_gla_risk_rate', label: 'GLA' },
    { flag: 'ptd_benefit', rate: 'total_ptd_risk_rate', label: 'PTD' },
    { flag: 'ci_benefit', rate: 'total_ci_risk_rate', label: 'CI' },
    { flag: 'sgla_benefit', rate: 'total_sgla_risk_rate', label: 'SGLA' },
    { flag: 'phi_benefit', rate: 'total_phi_risk_rate', label: 'PHI' }
  ]

  const out: string[] = []
  for (const rs of (props.resultSummaries ?? []) as any[]) {
    const cat = findCat(rs.category)
    if (!cat) continue
    for (const { flag, rate, label } of checks) {
      const enabled = cat[flag] === true || cat[flag] === 1 || cat[flag] === '1'
      const value = Number(rs[rate] ?? 0)
      if (enabled && (!isFinite(value) || value === 0)) {
        out.push(
          `${rs.category}: ${label} benefit is enabled but the risk rate is zero. Check that base rates exist for the scheme's risk rate code.`
        )
      }
    }
  }
  return out
})

const zeroRateWarningsCapped = computed(() =>
  zeroRateBenefitWarnings.value.slice(0, 5)
)
const zeroRateWarningsOverflow = computed(() =>
  Math.max(0, zeroRateBenefitWarnings.value.length - 5)
)

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

// "Show all columns" toggle state. Default is summary mode (only the
// renderer-visible columns + a few extras); flipping to full pulls every
// column the Go struct declares. Only the three tables that have a
// summary projection on the backend benefit from this toggle.
const showAllColumns = ref(false)
const currentItem: any = ref(null)
const tablesWithColumnToggle = new Set([
  'member_rating_results',
  'member_data',
  'member_premium_schedules'
])
const supportsColumnToggle = computed(
  () =>
    !!currentItem.value &&
    tablesWithColumnToggle.has(currentItem.value.value) &&
    !!selectedTable.value
)
const currentFieldSet = computed<'summary' | 'full'>(() =>
  showAllColumns.value ? 'full' : 'summary'
)
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

// Tracks which list item the user is currently loading so the View button
// for that row shows a spinner. Bordereaux is computed on-the-fly from
// MemberRatingResult, which can take longer than a plain DB read — the
// per-row spinner plus the contextual loadingLabel (below) make it clear
// the wait is due to live computation, not a stalled fetch.
const loadingItem = ref<string | null>(null)
const loadingLabel = computed(() => {
  if (loadingItem.value === 'bordereaux') {
    return 'Computing bordereaux from latest pricing inputs…'
  }
  return 'Loading data…'
})

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
  // Reset the column-toggle so the next "View" starts in summary mode.
  showAllColumns.value = false
  currentItem.value = null
}

// reloadCurrentTable re-fetches the currently-open table with the latest
// fieldSet (driven by the "Show all columns" switch). For the infinite
// model AG Grid re-requests blocks via the dataSource closure (which
// reads currentFieldSet at call time), so we just bump the grid key to
// force the column definitions to refresh from a fresh sample.
const reloadCurrentTable = async () => {
  if (!currentItem.value) return
  await viewTable(currentItem.value)
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
// Formatted timestamp shown beside each populated result row so the user
// can see at a glance when the underlying calculation last produced these
// records. Returns empty string when the quote has never been calculated.
const calculationCompletedLabel = computed(() => {
  const ts = props.quote?.calculation_completed_at
  if (!ts) return ''
  const d = new Date(ts)
  if (isNaN(d.getTime())) return ''
  const datePart = d.toLocaleDateString(undefined, {
    day: '2-digit',
    month: 'short',
    year: 'numeric'
  })
  const timePart = d.toLocaleTimeString(undefined, {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
  return `Calculated ${datePart} at ${timePart}`
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
  loadingItem.value = item.value
  currentItem.value = item

  try {
    currentTableType.value = item.value
    console.log('Item to view:', item)

    // Use AG Grid's infinite-row model for anything above the backend's
    // default-page size. The backend caps page size at 500 and returns
    // `has_more` so the grid can keep paging.
    const shouldUseInfiniteModel = item.count > 500
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
          limit: 1,
          fields: currentFieldSet.value
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
              { offset: startRow, limit, fields: currentFieldSet.value }
            )
            rowCount.value = endRow
            const rows: any[] =
              res.data && Array.isArray(res.data.data) ? res.data.data : []
            const orderedData = rows.map((item) =>
              _.fromPairs(
                (res.data?.json_tags || []).map((key) => [key, item[key]])
              )
            )
            // AG Grid SSRM (LazyCache) expects { rowData, rowCount? }.
            // rowCount is the total when known — present it whenever the
            // backend signals `has_more: false` OR the page came back
            // short. Otherwise omit it so AG Grid keeps requesting pages.
            const isLastPage =
              res.data?.has_more === false || rows.length < limit
            const response: { rowData: any[]; rowCount?: number } = {
              rowData: orderedData
            }
            if (isLastPage) {
              response.rowCount = startRow + orderedData.length
            }
            params.success(response)
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
      // Use original client-side approach for smaller datasets — still
      // pass `fields` so the payload stays trim unless the user has
      // toggled "Show all columns".
      const res = await GroupPricingService.getQuoteTable(
        props.quote.id,
        item.value,
        { fields: currentFieldSet.value }
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
  } finally {
    loadingItem.value = null
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

<style scoped>
.warnings-banner-title {
  background-color: rgb(var(--v-theme-warning) / 0.12);
}
</style>
