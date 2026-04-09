<template>
  <div>
    <!-- Chart rendered automatically from chart_config -->
    <ag-charts
      v-if="chartOptions"
      :options="chartOptions as any"
      :style="{ height: chartHeight, width: '100%' }"
    />

    <!-- Data grid below the chart -->
    <ag-grid-vue
      v-if="rowData.length > 0"
      class="ag-theme-balham"
      :style="{ height: gridHeight, width: '100%' }"
      :rowData="rowData"
      :columnDefs="columnDefs"
      :rowHeight="22"
      :headerHeight="26"
      :defaultColDef="defaultColDef"
      @grid-ready="onGridReady"
    />

    <div v-if="rowData.length === 0" class="no-data-message">
      No data available for this summary.
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { AgCharts } from 'ag-charts-vue3'
import { AgGridVue } from 'ag-grid-vue3'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-balham.css'

interface ChartSeries {
  key: string
  label: string
  type?: string
  axis?: string
  // box-plot fields
  min?: number
  mean?: number
  median?: number
  max?: number
  pct?: number
}

interface ChartConfig {
  chart_type:
    | 'bar'
    | 'line'
    | 'stacked_bar'
    | 'dual_axis'
    | 'tornado'
    | 'box_plot'
  x_key?: string
  series?: ChartSeries[]
  reference_lines?: { value: number; label: string }[]
}

interface TableWithChart {
  table_name: string
  chart_config?: ChartConfig
  data: Record<string, unknown>[]
}

const props = defineProps<{
  table: TableWithChart
  chartHeight?: string
  gridHeight?: string
}>()

const gridApi = ref<any>(null)

const chartHeight = computed(() => props.chartHeight ?? '320px')
const gridHeight = computed(() => props.gridHeight ?? '260px')

const rowData = computed<Record<string, unknown>[]>(() => {
  const d = props.table?.data
  if (!d) return []
  return Array.isArray(d) ? d : [d as Record<string, unknown>]
})

// Auto-generate column definitions from data keys
const columnDefs = computed(() => {
  if (rowData.value.length === 0) return []
  return Object.keys(rowData.value[0]).map((key) => ({
    headerName: key.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase()),
    field: key,
    sortable: true,
    filter: true,
    resizable: true,
    minWidth: 100,
    valueFormatter: (params: any) => {
      if (typeof params.value === 'number') {
        if (Math.abs(params.value) >= 1000) {
          return new Intl.NumberFormat('en-ZA', {
            minimumFractionDigits: 0,
            maximumFractionDigits: 0
          }).format(params.value)
        }
        return Number(params.value).toFixed(4)
      }
      return params.value ?? ''
    }
  }))
})

const defaultColDef = { flex: 1, minWidth: 100 }

function onGridReady(params: any) {
  gridApi.value = params.api
  params.api.autoSizeAllColumns()
}

// ── Chart options factory ─────────────────────────────────────────────────────
// Using `any` return type to avoid overly restrictive AgChartOptions union
// mismatches — the library runtime handles validation.
const chartOptions = computed<any | null>(() => {
  const cfg = props.table?.chart_config
  if (!cfg || rowData.value.length === 0) return null

  const title = {
    text: props.table.table_name,
    fontSize: 14,
    fontWeight: 'bold'
  }
  const data = rowData.value

  switch (cfg.chart_type) {
    case 'bar':
      return buildBarChart(cfg, data, title, false)
    case 'stacked_bar':
      return buildBarChart(cfg, data, title, true)
    case 'line':
      return buildLineChart(cfg, data, title)
    case 'dual_axis':
      return buildDualAxisChart(cfg, data, title)
    case 'tornado':
      return buildTornadoChart(cfg, title)
    case 'box_plot':
      return buildBoxPlotChart(cfg, title)
    default:
      return null
  }
})

// ── Bar / Stacked bar ─────────────────────────────────────────────────────────
function buildBarChart(
  cfg: ChartConfig,
  data: any[],
  title: any,
  stacked: boolean
) {
  const series = (cfg.series ?? []).map((s) => ({
    type: 'bar',
    xKey: cfg.x_key ?? 'accident_year',
    yKey: s.key,
    yName: s.label,
    stacked
  }))
  return {
    title,
    data,
    series,
    axes: [
      { type: 'category', position: 'bottom' },
      { type: 'number', position: 'left' }
    ],
    legend: { position: 'top', enabled: true }
  }
}

// ── Line (with optional reference lines as extra constant series) ─────────────
function buildLineChart(cfg: ChartConfig, data: any[], title: any) {
  const mainSeries = (cfg.series ?? []).map((s) => ({
    type: 'line',
    xKey: cfg.x_key ?? 'accident_year',
    yKey: s.key,
    yName: s.label,
    marker: { enabled: true, size: 6 }
  }))

  // Inject reference lines as constant-value line series
  const refData = data.map((row) => {
    const copy: any = { ...row }
    for (const rl of cfg.reference_lines ?? []) {
      copy[`_ref_${rl.label.replace(/\s+/g, '_')}`] = rl.value
    }
    return copy
  })
  const refSeries = (cfg.reference_lines ?? []).map((rl) => ({
    type: 'line',
    xKey: cfg.x_key ?? 'accident_year',
    yKey: `_ref_${rl.label.replace(/\s+/g, '_')}`,
    yName: rl.label,
    stroke: '#e53935',
    strokeDasharray: [6, 4],
    marker: { enabled: false }
  }))

  return {
    title,
    data: refData,
    series: [...mainSeries, ...refSeries],
    axes: [
      { type: 'category', position: 'bottom' },
      { type: 'number', position: 'left' }
    ],
    legend: { position: 'top', enabled: true }
  }
}

// ── Dual axis (primary bars + secondary lines) ────────────────────────────────
function buildDualAxisChart(cfg: ChartConfig, data: any[], title: any) {
  const primarySeries: any[] = []
  const secondarySeries: any[] = []

  for (const s of cfg.series ?? []) {
    if (s.axis === 'secondary') {
      secondarySeries.push({
        type: 'line',
        xKey: cfg.x_key ?? 'accident_year',
        yKey: s.key,
        yName: s.label,
        yAxisKey: 'right-axis',
        marker: { enabled: true, size: 6 }
      })
    } else {
      primarySeries.push({
        type: s.type === 'line' ? 'line' : 'bar',
        xKey: cfg.x_key ?? 'accident_year',
        yKey: s.key,
        yName: s.label
      })
    }
  }

  return {
    title,
    data,
    series: [...primarySeries, ...secondarySeries],
    axes: [
      { type: 'category', position: 'bottom' },
      {
        type: 'number',
        position: 'left',
        keys: primarySeries.map((s) => s.yKey)
      },
      {
        type: 'number',
        position: 'right',
        id: 'right-axis',
        keys: secondarySeries.map((s) => s.yKey),
        title: { text: secondarySeries.map((s) => s.yName).join(' / ') }
      }
    ],
    legend: { position: 'top', enabled: true }
  }
}

// ── Tornado (horizontal bar — % deviation from BEL) ──────────────────────────
function buildTornadoChart(cfg: ChartConfig, title: any) {
  const row = rowData.value[0] ?? {}
  const methodData = (cfg.series ?? []).map((s) => ({
    method: s.label,
    deviation: (row[s.key] as number) ?? 0
  }))

  return {
    title,
    data: methodData,
    series: [
      {
        type: 'bar',
        direction: 'horizontal',
        xKey: 'method',
        yKey: 'deviation',
        yName: '% Deviation from BEL',
        formatter: ({ datum }: any) => ({
          fill: datum.deviation >= 0 ? '#1976d2' : '#e53935'
        })
      }
    ],
    axes: [
      { type: 'category', position: 'left' },
      {
        type: 'number',
        position: 'bottom',
        title: { text: '% Deviation from BEL' }
      }
    ],
    legend: { enabled: false }
  }
}

// ── Box-plot rendered as grouped bar (min / mean / median / max / nth-pct) ────
// Uses a simple grouped bar instead of enterprise box-plot for broad compatibility.
function buildBoxPlotChart(cfg: ChartConfig, title: any) {
  const methodData = (cfg.series ?? []).map((s) => ({
    method: s.label ?? '',
    Minimum: s.min ?? 0,
    Mean: s.mean ?? 0,
    Median: s.median ?? 0,
    NthPercentile: s.pct ?? 0,
    Maximum: s.max ?? 0
  }))

  const belValue = (cfg.reference_lines ?? [])[0]?.value ?? 0
  // Add BEL reference line as extra column
  const dataWithRef = methodData.map((r) => ({ ...r, BEL: belValue }))

  return {
    title,
    data: dataWithRef,
    series: [
      {
        type: 'bar',
        xKey: 'method',
        yKey: 'Minimum',
        yName: 'Minimum',
        fill: '#90caf9'
      },
      {
        type: 'bar',
        xKey: 'method',
        yKey: 'Mean',
        yName: 'Mean',
        fill: '#1976d2'
      },
      {
        type: 'bar',
        xKey: 'method',
        yKey: 'Median',
        yName: 'Median',
        fill: '#42a5f5'
      },
      {
        type: 'bar',
        xKey: 'method',
        yKey: 'NthPercentile',
        yName: 'Nth Percentile',
        fill: '#0d47a1'
      },
      {
        type: 'bar',
        xKey: 'method',
        yKey: 'Maximum',
        yName: 'Maximum',
        fill: '#bbdefb'
      },
      {
        type: 'line',
        xKey: 'method',
        yKey: 'BEL',
        yName: 'BEL (Deterministic)',
        stroke: '#e53935',
        strokeDasharray: [6, 4],
        marker: { enabled: false }
      }
    ],
    axes: [
      { type: 'category', position: 'bottom' },
      { type: 'number', position: 'left', title: { text: 'Reserve (ZAR)' } }
    ],
    legend: { position: 'top', enabled: true }
  }
}
</script>

<style scoped>
.no-data-message {
  padding: 16px;
  color: #888;
  font-style: italic;
  text-align: center;
}
</style>
