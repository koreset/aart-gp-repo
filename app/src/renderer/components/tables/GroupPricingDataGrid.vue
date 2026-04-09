<!-- eslint-disable vue/attribute-hyphenation -->
<template>
  <v-container>
    <v-row class="align-center mb-1">
      <v-col class="d-flex align-center justify-space-between py-1">
        <h4 v-if="tableTitle" class="h4-title">{{ tableTitle }}</h4>
        <div v-else></div>
        <div class="d-flex align-center ga-1">
          <!-- Three parallel export icons -->
          <v-btn
            v-if="showExport"
            icon
            size="x-small"
            variant="tonal"
            color="success"
            title="Export to CSV"
            @click="exportDataCsv"
            ><v-icon size="small">mdi-file-delimited</v-icon></v-btn
          >
          <v-btn
            v-if="showExport"
            icon
            size="x-small"
            variant="tonal"
            color="error"
            title="Export to PDF"
            @click="exportDataPdf"
            ><v-icon size="small">mdi-file-pdf-box</v-icon></v-btn
          >
          <v-btn
            v-if="showExport"
            icon
            size="x-small"
            variant="tonal"
            color="success"
            title="Export to Excel"
            @click="exportDataExcelGrid"
            ><v-icon size="small">mdi-microsoft-excel</v-icon></v-btn
          >
          <!-- Legacy full-export button (valuation results) -->
          <v-btn
            v-if="showFullExport"
            icon
            size="x-small"
            variant="tonal"
            color="primary"
            :loading="exportLoader"
            title="Export All"
            @click="exportDataExcel"
            ><v-icon size="small">mdi-download-circle</v-icon></v-btn
          >
          <v-btn
            v-if="showExportAllCsv"
            icon
            size="x-small"
            variant="tonal"
            color="primary"
            :loading="exportAllCsvLoader"
            title="Export All to CSV"
            @click="exportAllDataCsv"
            ><v-icon size="small">mdi-download</v-icon></v-btn
          >
          <!-- Delete selected -->
          <v-btn
            v-if="showDeleteButton"
            size="small"
            variant="outlined"
            rounded
            color="red"
            class="ml-2"
            @click="deleteRow"
            >Delete Selected</v-btn
          >
          <!-- Close -->
          <v-btn
            v-if="showCloseButton"
            variant="plain"
            size="small"
            @click="clearRowData"
            >Close</v-btn
          >
        </div>
      </v-col>
    </v-row>
    <p v-if="exportLoader" class="text-caption text-medium-emphasis mb-1"
      >Processing data, this may take a moment…</p
    >
    <v-row>
      <v-col>
        <ag-grid-vue
          class="ag-theme-balham"
          :readOnlyEdit="true"
          :enableRangeSelection="enableRangeSelection"
          :enableCharts="enableCharts"
          :statusBar="statusBar"
          :style="{ height: gridHeight }"
          :class="gridThemeClass"
          :rowData="props.useInfiniteModel ? undefined : localRowData"
          :rowHeight="gridOptions.rowHeight"
          :headerHeight="gridOptions.headerHeight"
          :columnDefs="localColumnDefs"
          :autoGroupColumnDef="autoGroupColumnDef"
          :rowModelType="props.useInfiniteModel ? 'serverSide' : 'clientSide'"
          :cacheBlockSize="props.useInfiniteModel ? 1000 : undefined"
          :cacheOverflowSize="props.useInfiniteModel ? 2 : undefined"
          :maxConcurrentDatasourceRequests="
            props.useInfiniteModel ? 2 : undefined
          "
          :infiniteInitialRowCount="props.useInfiniteModel ? 1 : undefined"
          :maxBlocksInCache="props.useInfiniteModel ? 2 : undefined"
          :components="gridComponents"
          @row-selected="onRowSelected"
          @row-clicked="onRowClicked"
          @grid-ready="onGridReady"
          @cell-edit-request="OnCellEditRequest"
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
import CountStatusBarComponent from './countStatusBarComponent.vue'
import jsPDF from 'jspdf'
import 'jspdf-autotable'

const props = defineProps<{
  rowData?: any[]
  columnDefs: any[]
  rowModel?: string
  pagination?: boolean
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
  showExportAllCsv?: boolean
  runId?: string
  productCode?: string
  runName?: string
  useInfiniteModel?: boolean
  dataSource?: any
  suppressAutoSize?: boolean
  density?: 'default' | 'compact' | 'comfortable'
}>()

// const emit = defineEmits(['delete-row'])

const emit = defineEmits<{
  (e: 'update:row-deleted', value: any): void
  (e: 'update:clear-data', value: any): void
  (e: 'update-column', value: { colId: string; newValue: any }): void
  (e: 'export-all-csv'): void
  (e: 'row-clicked', value: any): void
}>()

// Grid variables
const showDeleteButton = ref(false)
const selectedRow = ref(null)
const exportLoader = ref(false)
const exportAllCsvLoader = ref(false)

// const localRowData = ref(props.rowData)
// const localColumnDefs = ref(props.columnDefs)
const localRowData = computed(() => props.rowData)
const localColumnDefs = computed(() => props.columnDefs)

const localShowExport = ref(true)

const gridApi: any = ref(null)
const columnApi: any = ref(null)

// Define the autoSizeStrategy
// const autoSizeStrategy = ref({
//   type: 'fitCellContents'
//   // You can also provide 'skipHeader: false' here, but it's the default.
//   // skipHeader: false
// })

const autoSizeAll = (skipHeader: boolean) => {
  const allColumnIds: string[] = []
  gridApi.value!.getColumns()!.forEach((column) => {
    allColumnIds.push(column.getId())
  })
  gridApi.value!.autoSizeColumns(allColumnIds, skipHeader)
}

const onGridReady = (params) => {
  console.log('Grid ready, useInfiniteModel:', props.useInfiniteModel)
  gridApi.value = params.api
  columnApi.value = params.columnApi

  if (!props.useInfiniteModel && !props.suppressAutoSize) {
    gridApi.value.autoSizeAllColumns()
    autoSizeAll(false)
  }

  // Set datasource if using infinite model
  if (props.useInfiniteModel && props.dataSource) {
    console.log('Setting datasource for infinite model in onGridReady')
    // gridApi.value.setDatasource(props.dataSource)
    gridApi.value.setGridOption('serverSideDatasource', props.dataSource)
  }
}

const onRowClicked = (event) => {
  // Emit row-clicked event to parent component
  emit('row-clicked', event)
}

const OnCellEditRequest = (event) => {
  const colId = event.column.getColId()
  const newValue = event.newValue

  console.log(`Request to update entire '${colId}' column to '${newValue}'`)

  // 3. Emit the change request to the parent component
  emit('update-column', { colId, newValue })
}

// Grid options based on density
const gridOptions = computed(() => {
  const baseOptions = {
    rowHeight: 22,
    headerHeight: 26
  }
  return baseOptions
})

// Grid height based on density and content
const gridHeight = computed(() => {
  const maxHeight =
    props.density === 'compact'
      ? 650
      : props.density === 'comfortable'
        ? 750
        : 600
  const minHeight = 300 // Minimum height to ensure grid is usable

  if (!localRowData.value || localRowData.value.length === 0) {
    return `${minHeight}px`
  }

  const rowHeight = gridOptions.value.rowHeight
  const headerHeight = gridOptions.value.headerHeight
  const statusBarHeight = 30 // Approximate height of status bar
  const padding = 20 // Some padding for borders and scrollbars

  // Calculate content height: header + (number of rows * row height) + status bar + padding
  const contentHeight =
    headerHeight +
    localRowData.value.length * rowHeight +
    statusBarHeight +
    padding

  // Use the smaller of content height or max height, but not less than min height
  const calculatedHeight = Math.max(
    minHeight,
    Math.min(contentHeight, maxHeight)
  )

  return `${calculatedHeight}px`
})

const gridThemeClass = computed(() => {
  return 'ag-theme-balham'
})

const autoGroupColumnDef = { cellRendererParams: { suppressCount: true } }

watch(
  () => props.showExport,
  (newVal) => {
    localShowExport.value = newVal
    // emit('update:showExport', newVal);
  }
)

// Watch for dataSource changes to set it on the grid
watch(
  () => props.dataSource,
  (newDataSource) => {
    console.log(
      'DataSource watcher triggered:',
      !!newDataSource,
      'Grid ready:',
      !!gridApi.value,
      'Use infinite:',
      props.useInfiniteModel
    )
    if (gridApi.value && props.useInfiniteModel && newDataSource) {
      console.log('Setting datasource via watcher')
      gridApi.value.setDatasource(newDataSource)
    }
  }
)
const gridComponents = {
  countComponent: CountStatusBarComponent
}

const enableRangeSelection = true
const enableCharts = true

const statusBar = {
  statusPanels: [
    // { statusPanel: 'agTotalAndFilteredRowCountComponent', align: 'left' },
    {
      statusPanel: 'countComponent',
      align: 'left',
      statusPanelParams: {
        isInfinite: props.useInfiniteModel
      }
    },
    { statusPanel: 'agTotalRowCountComponent' },
    // { statusPanel: 'agFilteredRowCountComponent' },
    // { statusPanel: 'agSelectedRowCountComponent' },
    {
      statusPanel: 'agAggregationComponent',
      statusPanelParams: {
        // possible values are: 'count', 'sum', 'min', 'max', 'avg'
        aggFuncs: ['avg', 'sum'],
        suppressAggFilteredOnly: true
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

const exportDataPdf = () => {
  const visibleCols = (localColumnDefs.value || []).filter(
    (c: any) => !c.hide && !c.rowGroup
  )
  const columns = visibleCols.map((c: any) => ({
    header: c.headerName || c.field || '',
    dataKey: c.field
  }))
  const colCount = columns.length

  const formatVal = (v: any) => {
    if (v === null || v === undefined || v === '') return ''
    const n = Number(v)
    if (!isNaN(n) && v !== '')
      return n.toLocaleString('en-ZA', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      })
    return String(v)
  }

  const body: any[] = []
  gridApi.value?.forEachNodeAfterFilterAndSort((node: any) => {
    if (node.group) {
      body.push([
        {
          content: node.key ?? '',
          colSpan: colCount,
          styles: {
            fontStyle: 'bold',
            fillColor: [45, 85, 110],
            textColor: [255, 255, 255]
          }
        }
      ])
    } else if (node.data) {
      body.push(visibleCols.map((c: any) => formatVal(node.data[c.field])))
    }
  })

  // eslint-disable-next-line new-cap
  const doc: any = new jsPDF({ orientation: 'landscape' })
  if (props.tableTitle) {
    doc.setFontSize(11)
    doc.setFont('helvetica', 'bold')
    doc.text(props.tableTitle, 14, 16)
  }
  doc.autoTable({
    startY: props.tableTitle ? 22 : 14,
    columns,
    body,
    theme: 'striped',
    styles: { fontSize: 7.5, cellPadding: 1.5 },
    headStyles: { fillColor: [70, 100, 120], textColor: 255, fontStyle: 'bold' }
  })
  const pageCount = doc.internal.getNumberOfPages()
  for (let i = 1; i <= pageCount; i++) {
    doc.setPage(i)
    doc.setFontSize(7)
    doc.setFont('helvetica', 'normal')
    doc.text(
      `Page ${i} of ${pageCount}  |  Generated ${new Date().toLocaleDateString()}`,
      14,
      doc.internal.pageSize.getHeight() - 6
    )
  }
  const fileName = (props.tableTitle || 'export').replace(/\s+/g, '_')
  doc.save(`${fileName}.pdf`)
}

const exportDataExcelGrid = () => {
  const fileName = (props.tableTitle || 'export').replace(/\s+/g, '_')
  gridApi.value?.exportDataAsExcel({ fileName: `${fileName}.xlsx` })
}

const exportAllDataCsv = () => {
  exportAllCsvLoader.value = true
  emit('export-all-csv')
  // The parent component will handle the actual export and reset the loader
}

// Expose a method to reset the loader from parent
const resetExportAllCsvLoader = () => {
  exportAllCsvLoader.value = false
}

// Expose the method so parent can call it
defineExpose({
  resetExportAllCsvLoader
})

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
  if (event.api.getSelectedRows().length > 0) {
    showDeleteButton.value = true
    selectedRow.value = event.api.getSelectedRows()[0]
  } else {
    showDeleteButton.value = false
    selectedRow.value = null
  }
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
