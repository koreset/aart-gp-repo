<template>
  <v-card variant="outlined" class="pa-3">
    <div class="d-flex align-center mb-2 flex-wrap">
      <div class="text-subtitle-1 font-weight-bold">
        Loss-Ratio Trend (Rolling 12 months)
      </div>
      <v-tooltip location="top" max-width="360">
        <template #activator="{ props: tipProps }">
          <v-icon
            v-bind="tipProps"
            size="small"
            icon="mdi-information-outline"
            class="ml-2 mr-1"
            style="cursor: help"
          />
        </template>
        <span>
          Each line is one scheme's rolling-12-month ALR (claims in trailing 12
          months ÷ time-weighted annual premium) at the end of each month. The
          dashed line is the portfolio aggregate. Default selection is the top 5
          schemes by current R12M ALR; pick others from the dropdown.
        </span>
      </v-tooltip>
      <v-spacer />
      <ChartMenu
        :chart-ref="trendChart"
        title="Loss Ratio Trend"
        :data="csvData"
      />
    </div>

    <v-row dense class="align-center mb-2">
      <v-col cols="12" md="9">
        <v-autocomplete
          v-model="selectedSchemeIds"
          :items="schemeOptions"
          item-title="label"
          item-value="scheme_id"
          label="Schemes plotted"
          variant="outlined"
          density="compact"
          chips
          closable-chips
          multiple
          hide-details
        />
      </v-col>
      <v-col cols="auto">
        <v-checkbox
          v-model="showPortfolio"
          label="Portfolio line"
          density="compact"
          hide-details
        />
      </v-col>
      <v-col cols="auto">
        <v-btn
          variant="text"
          size="small"
          color="primary"
          prepend-icon="mdi-restore"
          @click="resetSelection"
        >
          Reset to top 5
        </v-btn>
      </v-col>
    </v-row>

    <ag-charts v-if="chartOptions" ref="trendChart" :options="chartOptions" />
    <div
      v-if="trend.months.length === 0"
      class="text-center text-medium-emphasis my-4"
    >
      No trend data available.
    </div>
  </v-card>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { AgCharts } from 'ag-charts-vue3'
import ChartMenu from '@/renderer/components/charts/ChartMenu.vue'

interface TrendSeries {
  scheme_id: number
  scheme_name: string
  annual_premium: number
  current_r12m_alr: number | null
  values: (number | null)[]
}
interface LossRatioTrend {
  months: string[]
  portfolio: (number | null)[]
  schemes: TrendSeries[]
}

const props = defineProps<{ trend: LossRatioTrend }>()
const trendChart: any = ref(null)

const schemeOptions = computed(() =>
  props.trend.schemes.map((s) => ({
    scheme_id: s.scheme_id,
    label:
      s.current_r12m_alr != null
        ? `${s.scheme_name} (R12M ${s.current_r12m_alr.toFixed(1)}%)`
        : s.scheme_name
  }))
)

const defaultSelection = computed(() =>
  [...props.trend.schemes]
    .filter((s) => s.current_r12m_alr != null)
    .sort((a, b) => (b.current_r12m_alr ?? 0) - (a.current_r12m_alr ?? 0))
    .slice(0, 5)
    .map((s) => s.scheme_id)
)

const selectedSchemeIds = ref<number[]>([])
const showPortfolio = ref(true)

const resetSelection = () => {
  selectedSchemeIds.value = [...defaultSelection.value]
}

watch(
  () => props.trend,
  () => {
    // Re-seed defaults whenever a fresh payload arrives.
    selectedSchemeIds.value = [...defaultSelection.value]
  },
  { immediate: true, deep: false }
)

// Build chart data: each row = { month, portfolio, [scheme_<id>]: value }
const chartData = computed(() => {
  const selected = new Set(selectedSchemeIds.value)
  const seriesById = new Map(props.trend.schemes.map((s) => [s.scheme_id, s]))
  return props.trend.months.map((m, i) => {
    const row: Record<string, any> = { month: m }
    if (showPortfolio.value) {
      row.portfolio = props.trend.portfolio[i]
    }
    for (const id of selected) {
      const s = seriesById.get(id)
      if (!s) continue
      row[`s_${id}`] = s.values[i]
    }
    return row
  })
})

const seriesDefs = computed(() => {
  const out: any[] = []
  if (showPortfolio.value) {
    out.push({
      type: 'line',
      xKey: 'month',
      yKey: 'portfolio',
      yName: 'Portfolio R12M ALR',
      stroke: '#1976D2',
      strokeWidth: 3,
      marker: { enabled: true, size: 6 },
      lineDash: [6, 4]
    })
  }
  for (const id of selectedSchemeIds.value) {
    const s = props.trend.schemes.find((x) => x.scheme_id === id)
    if (!s) continue
    out.push({
      type: 'line',
      xKey: 'month',
      yKey: `s_${id}`,
      yName: s.scheme_name,
      strokeWidth: 2,
      marker: { enabled: true, size: 4 }
    })
  }
  return out
})

const chartOptions = computed<any>(() => ({
  data: chartData.value,
  background: { fill: 'transparent' },
  height: 380,
  series: seriesDefs.value,
  axes: [
    {
      type: 'category',
      position: 'bottom',
      title: { text: 'Month' },
      label: { rotation: -45 }
    },
    {
      type: 'number',
      position: 'left',
      title: { text: 'Rolling-12-Month ALR (%)' },
      crossLines: [
        {
          type: 'line',
          value: 100,
          stroke: '#c62828',
          lineDash: [4, 4],
          strokeWidth: 1,
          label: { text: '100%', position: 'right', color: '#c62828' }
        }
      ]
    }
  ],
  legend: { enabled: true, position: 'bottom' }
}))

// Flatten for CSV export: one row per month per active series.
const csvData = computed(() => {
  const out: Record<string, any>[] = []
  const seriesById = new Map(props.trend.schemes.map((s) => [s.scheme_id, s]))
  for (let i = 0; i < props.trend.months.length; i++) {
    const month = props.trend.months[i]
    if (showPortfolio.value) {
      out.push({
        month,
        scheme: 'PORTFOLIO',
        rolling_12m_alr_pct: props.trend.portfolio[i]
      })
    }
    for (const id of selectedSchemeIds.value) {
      const s = seriesById.get(id)
      if (!s) continue
      out.push({
        month,
        scheme: s.scheme_name,
        rolling_12m_alr_pct: s.values[i]
      })
    }
  }
  return out
})
</script>
