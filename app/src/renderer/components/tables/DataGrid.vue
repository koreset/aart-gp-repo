<!-- eslint-disable vue/attribute-hyphenation -->
<template>
  <v-container>
    <v-row v-if="tableTitle">
      <v-col cols="12">
        <h4 class="h4-title title-display"
          ><span>{{ tableTitle }} Data</span
          ><span v-if="showCloseButton" class="text-left"
            ><v-btn variant="plain" size="small" @click="clearRowData"
              >Close</v-btn
            ></span
          ></h4
        >
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="6">
        <v-btn
          v-if="showExport"
          size="small"
          color="primary"
          rounded
          class="custom-btn primary white--text mt-4"
          @click="exportDataCsv"
          >Export to Csv</v-btn
        >
        <v-btn
          v-if="showFullExport"
          size="small"
          color="primary"
          rounded
          :loading="exportLoader"
          class="custom-btn ml-4 primary white--text mt-4"
          @click="exportDataExcel"
          >Export All</v-btn
        >
        <p v-if="exportLoader" class="mt-4"
          >processing data. this will take a while...</p
        >

        <v-btn
          v-if="showDeleteButton"
          size="small"
          variant="outlined"
          rounded
          color="red"
          class="custom-btn primary white--text ml-4 mt-4"
          @click="deleteRow"
          >Delete Selected</v-btn
        >
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <ag-grid-vue
          :enableRangeSelection="enableRangeSelection"
          :enableCharts="enableCharts"
          :statusBar="statusBar"
          :style="{ height: gridHeight }"
          :class="gridThemeClass"
          :rowData="localRowData"
          :columnDefs="localColumnDefs"
          :defaultColDef="props.defaultColDef"
          :autoSizeStrategy="autoSizeStrategy"
          :rowHeight="gridOptions.rowHeight"
          :headerHeight="gridOptions.headerHeight"
          :rowSelection="props.rowSelection || 'multiple'"
          :pagination="props.pagination"
          :paginationPageSize="props.paginationPageSize"
          @row-selected="onRowSelected"
          @row-double-clicked="onRowDoubleClicked"
          @grid-ready="onGridReady"
        >
        </ag-grid-vue>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
// import { ref } from 'vue';
import 'ag-grid-community/styles/ag-grid.css' // Core CSS
import 'ag-grid-community/styles/ag-theme-balham.css' // Theme

import { AgGridVue } from 'ag-grid-vue3' // Vue Grid Logic
import { ref, watch, computed } from 'vue'
import PhiValuationService from '@/renderer/api/PhiValuationService'

const props = defineProps<{
  rowData?: any[]
  columnDefs: any[]
  defaultColDef?: Record<string, any>
  rowModel?: string
  pagination?: boolean
  paginationPageSize?: number
  rowSelection?: string | null
  tableTitle?: string | null
  rowCount?: number
  tableName?: string | null
  chartTitle?: string | null
  chartXAxisTitle?: string | null
  chartYAxisTitle?: string | null
  showExport?: boolean
  showCloseButton?: boolean
  showFullExport?: boolean
  showDelete?: boolean
  runId?: string
  productCode?: string
  runName?: string
  density?: 'default' | 'compact' | 'comfortable'
}>()

// const emit = defineEmits(['delete-row'])

const emit = defineEmits<{
  (e: 'update:row-deleted', value: any): void
  (e: 'update:clear-data', value: any): void
  (e: 'row-double-clicked', value: any): void
  (e: 'row-selection-changed', value: any[]): void
}>()

// Grid variables
const showDeleteButton = ref(false)
const selectedRow = ref(null)
const exportLoader = ref(false)

// const localRowData = ref(props.rowData)
// const localColumnDefs = ref(props.columnDefs)
const localRowData = computed(() => props.rowData)
const localColumnDefs = computed(() => props.columnDefs)

const localShowExport = ref(true)

const gridApi: any = ref(null)
const columnApi: any = ref(null)

// Define the autoSizeStrategy
const autoSizeStrategy = ref({
  type: 'fitCellContents'
  // You can also provide 'skipHeader: false' here, but it's the default.
  // skipHeader: false
})

// Grid options based on density
const gridOptions = computed(() => {
  const baseOptions = {
    rowHeight: 28,
    headerHeight: 32
  }
  return baseOptions
})

// Grid height based on density and content
const gridHeight = computed(() => {
  const maxHeight =
    props.density === 'compact'
      ? 400
      : props.density === 'comfortable'
        ? 600
        : 500
  const minHeight = 200 // Minimum height to ensure grid is usable

  if (!localRowData.value || localRowData.value.length === 0) {
    return `${minHeight}px`
  }

  const rowHeight = gridOptions.value.rowHeight
  const headerHeight = gridOptions.value.headerHeight
  const statusBarHeight = 30
  const paginationBarHeight = props.pagination ? 40 : 0
  const padding = 20

  // When pagination is on, size to fit the page rather than all rows
  const visibleRows = props.pagination
    ? Math.min(localRowData.value.length, props.paginationPageSize ?? 20)
    : localRowData.value.length

  const contentHeight =
    headerHeight +
    visibleRows * rowHeight +
    statusBarHeight +
    paginationBarHeight +
    padding

  const calculatedHeight = Math.max(
    minHeight,
    Math.min(contentHeight, maxHeight)
  )

  return `${calculatedHeight}px`
})

// Grid theme class based on density
const gridThemeClass = computed(() => {
  return props.density === 'compact'
    ? 'ag-theme-balham ag-theme-compact'
    : 'ag-theme-balham'
})

const autoSizeAll = (skipHeader: boolean) => {
  const allColumnIds: string[] = []
  gridApi.value!.getColumns()!.forEach((column) => {
    allColumnIds.push(column.getId())
  })
  gridApi.value!.autoSizeColumns(allColumnIds, skipHeader)
}

const onGridReady = (params) => {
  gridApi.value = params.api
  columnApi.value = params.columnApi
  gridApi.value.autoSizeAllColumns()

  autoSizeAll(false)
}

// watch(
//   () => props.rowData,
//   (newVal) => {
//     localRowData.value = newVal
//     // emit('update:rowData', newVal);
//   }
// )

// watch(
//   () => props.columnDefs,
//   (newVal) => {
//     localColumnDefs.value = newVal
//     // emit('update:columnDefs', newVal);
//   }
// )

watch(
  () => props.showExport,
  (newVal) => {
    localShowExport.value = newVal
    // emit('update:showExport', newVal);
  }
)

const enableRangeSelection = true
const enableCharts = true
const statusBar = {
  statusPanels: [
    { statusPanel: 'agTotalAndFilteredRowCountComponent', align: 'left' },
    {
      statusPanel: 'agAggregationComponent',
      statusPanelParams: {
        // possible values are: 'count', 'sum', 'min', 'max', 'avg'
        aggFuncs: ['avg', 'sum']
      }
    }
  ]
}

const exportDataCsv = () => {
  gridApi.value.exportDataAsCsv({
    suppressQuotes: true,
    allColumns: true
  })
}

const exportDataExcel = async () => {
  exportLoader.value = true
  const response = await PhiValuationService.getExcelResults(props.runId, null)

  exportLoader.value = false

  const fileURL = window.URL.createObjectURL(new Blob([response.data]))
  const fileLink = document.createElement('a')

  fileLink.href = fileURL
  fileLink.setAttribute(
    'download',
    'agg-results_' + props.runName + '_' + props.productCode + '.xlsx'
  )
  document.body.appendChild(fileLink)

  fileLink.click()
}

const onRowSelected = (event) => {
  const selectedRows = event.api.getSelectedRows()

  if (selectedRows.length > 0) {
    showDeleteButton.value = true
    selectedRow.value = selectedRows[0]
  } else {
    showDeleteButton.value = false
    selectedRow.value = null
  }

  // Emit selection change to parent
  emit('row-selection-changed', selectedRows)
}

const onRowDoubleClicked = (event) => {
  // Emit row-double-clicked event to parent component
  emit('row-double-clicked', event)
}

const deleteRow = () => {
  if (selectedRow.value == null) return
  // this.gridApi.applyTransaction({ remove: [this.selectedRow] });
  // emit event to parent component
  emit('update:row-deleted', selectedRow.value)
}

const clearRowData = () => {
  emit('update:clear-data', null)
}
</script>

<style scoped>
.title-display {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.h4-title {
  width: 100%;
  border-bottom: 1px solid #b4b3b3;
}
</style>
