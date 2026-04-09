<template>
  <v-dialog v-model="dialog" max-width="1150px">
    <base-card :show-actions="false">
      <template #header>
        <span class="headline">Credibility Analysis Dashboard</span>
        <v-btn
          size="small"
          icon
          class="float-right"
          style="
            position: absolute;
            top: -5px;
            right: 8px;
            z-index: 2;
            background-color: transparent;
            box-shadow: none;
          "
          @click="close"
        >
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </template>
      <template #default>
        <div v-if="loading" class="text-center my-4">
          <v-progress-circular indeterminate color="primary" />
        </div>
        <div v-else>
          <!-- Filters row -->
          <v-row class="mb-2" dense>
            <v-col cols="12" md="6">
              <v-range-slider
                v-model="experiencePeriodRange"
                :min="experiencePeriodMin"
                :max="experiencePeriodMax"
                :step="1"
                label="Experience Period"
                thumb-label
                density="compact"
                class="mt-2"
              />
            </v-col>
            <v-col cols="12" md="3">
              <v-select
                v-model="selectedYear"
                variant="outlined"
                density="compact"
                :items="yearOptions"
                label="Year"
                clearable
                class="mt-2"
              />
            </v-col>
            <v-col cols="12" md="3">
              <v-select
                v-model="selectedBenefitType"
                variant="outlined"
                density="compact"
                :items="benefitTypeOptions"
                label="Benefit Type (CI panel)"
                class="mt-2"
              />
            </v-col>
          </v-row>

          <!-- Panel tabs -->
          <v-tabs
            v-model="selectedTab"
            class="mb-3"
            density="compact"
            color="primary"
          >
            <v-tab value="curve">
              <v-icon start size="small"
                >mdi-chart-bell-curve-cumulative</v-icon
              >
              Credibility Curve
            </v-tab>
            <v-tab value="gap">
              <v-icon start size="small">mdi-chart-bar</v-icon>
              Gap Analysis
            </v-tab>
            <v-tab value="ci">
              <v-icon start size="small">mdi-chart-error-bar</v-icon>
              Experience CI
            </v-tab>
            <v-tab value="exposure">
              <v-icon start size="small">mdi-gauge</v-icon>
              Exposure Adequacy
            </v-tab>
          </v-tabs>

          <!-- No data state -->
          <div
            v-if="filteredData.length === 0"
            class="text-center pa-8 text-medium-emphasis"
          >
            <v-icon size="48" class="mb-2"
              >mdi-chart-timeline-variant-shimmer</v-icon
            >
            <div>No credibility data available for the selected filters.</div>
          </div>

          <template v-else>
            <!-- Panel 1: Credibility Curve -->
            <div v-if="selectedTab === 'curve'">
              <div class="panel-description">
                Schemes plotted on the theoretical credibility curve
                <code>sqrt(WLY / FCT)</code>. Manual overrides shown as separate
                points — deviation above the curve is conservative, below is
                aggressive.
              </div>
              <ag-charts
                v-if="curveChartOptions"
                :options="curveChartOptions"
                style="height: 420px; width: 100%"
              />
              <div v-else class="text-center text-medium-emphasis pa-4">
                Weighted life years data not yet available — recalculate quotes
                to populate.
              </div>
            </div>

            <!-- Panel 2: Gap Analysis -->
            <div v-if="selectedTab === 'gap'">
              <div class="panel-description">
                Calculated vs manually applied credibility per scheme. Orange
                bars indicate the actuary set a higher value (conservative); red
                indicates a lower override (aggressive).
              </div>
              <ag-charts
                v-if="gapChartOptions"
                :options="gapChartOptions"
                style="height: 420px; width: 100%"
              />
            </div>

            <!-- Panel 3: Experience CI -->
            <div v-if="selectedTab === 'ci'">
              <div class="panel-description">
                Annual experience rate with Poisson 95% confidence intervals
                <code>rate ± 1.96 × √claims / (members × period)</code>. Wide
                bands = noisy data; narrow bands = credible experience. The
                dashed line is the pooled average across all filtered schemes.
              </div>
              <ag-charts
                v-if="ciChartOptions"
                :options="ciChartOptions"
                style="height: 420px; width: 100%"
              />
            </div>

            <!-- Panel 4: Exposure Adequacy -->
            <div v-if="selectedTab === 'exposure'">
              <div class="panel-description">
                Progress of each scheme toward the full credibility threshold
                (weighted life years / FCT). Schemes below 25% need
                substantially more exposure before experience is reliable.
              </div>
              <ag-charts
                v-if="exposureChartOptions"
                :options="exposureChartOptions"
                style="height: 420px; width: 100%"
              />
              <div v-else class="text-center text-medium-emphasis pa-4">
                Full credibility threshold data not yet available — recalculate
                quotes to populate.
              </div>
            </div>
          </template>
        </div>
      </template>
      <template #actions>
        <v-spacer></v-spacer>
        <v-btn rounded variant="text" @click="close">Close</v-btn>
      </template>
    </base-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { AgCharts } from 'ag-charts-vue3'
import BaseCard from '../BaseCard.vue'
import GroupPricingService from '../../api/GroupPricingService'

const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits(['update:modelValue'])

const dialog = ref(props.modelValue)
const loading = ref(false)
const selectedTab = ref('curve')
const selectedYear = ref<number | null>(null)
const selectedBenefitType = ref<'gla' | 'ptd' | 'ci'>('gla')
const experiencePeriodRange = ref<[number, number]>([0, 10])
const experiencePeriodMin = ref(0)
const experiencePeriodMax = ref(10)
const yearOptions = ref<number[]>([])
const benefitTypeOptions = [
  { title: 'GLA', value: 'gla' },
  { title: 'PTD', value: 'ptd' },
  { title: 'CI', value: 'ci' }
]

let allData: any[] = []

// ─── Watchers ────────────────────────────────────────────────────────────────

watch(
  () => props.modelValue,
  (val) => {
    dialog.value = val
    if (val) fetchData()
  }
)
watch(dialog, (val) => emit('update:modelValue', val))

// ─── Filtered data ───────────────────────────────────────────────────────────

const filteredData = computed(() => {
  if (!allData.length) return []
  return allData.filter((d) => {
    const inPeriod =
      d.experience_period >= experiencePeriodRange.value[0] &&
      d.experience_period <= experiencePeriodRange.value[1]
    const inYear = selectedYear.value === null || d.year === selectedYear.value
    return inPeriod && inYear
  })
})

// ─── Fetch ───────────────────────────────────────────────────────────────────

async function fetchData() {
  loading.value = true
  try {
    const res = await GroupPricingService.getHistoricalCredibilityData()
    allData = res.data || []
    yearOptions.value = Array.from(
      new Set(allData.map((d: any) => d.year))
    ).sort((a: any, b: any) => b - a)
    const periods = allData.map((d: any) => d.experience_period)
    experiencePeriodMin.value = periods.length
      ? Math.floor(Math.min(...periods))
      : 0
    experiencePeriodMax.value = periods.length
      ? Math.ceil(Math.max(...periods))
      : 10
    experiencePeriodRange.value = [
      experiencePeriodMin.value,
      experiencePeriodMax.value
    ]
  } catch {
    allData = []
  } finally {
    loading.value = false
  }
}

// ─── Panel 1: Credibility Curve ───────────────────────────────────────────────

const curveChartOptions = computed<any>(() => {
  const data = filteredData.value
  if (!data.length) return null

  // Require weighted_life_years and full_credibility_threshold to draw the curve
  const hasCurveData = data.some(
    (d) => d.full_credibility_threshold > 0 && d.weighted_life_years >= 0
  )
  if (!hasCurveData) return null

  // Use the max FCT found in the dataset
  const fct = Math.max(
    ...data
      .map((d: any) => d.full_credibility_threshold)
      .filter((v: number) => v > 0)
  )
  const maxWly =
    Math.max(...data.map((d: any) => d.weighted_life_years), fct) * 1.1

  // Build smooth theoretical curve (60 points)
  const curvePoints = Array.from({ length: 61 }, (_, i) => {
    const wly = (i / 60) * maxWly
    return { wly, curve_credibility: Math.min(Math.sqrt(wly / fct), 1) }
  })

  // Deduplicate by scheme — credibility and WLY are scheme-level values identical
  // across all category rows for the same quote; take the latest year per scheme.
  const bySchemeLatest = new Map<string, any>()
  for (const d of data) {
    const existing = bySchemeLatest.get(d.scheme_name)
    if (!existing || d.year > existing.year)
      bySchemeLatest.set(d.scheme_name, d)
  }
  const deduped = Array.from(bySchemeLatest.values())

  // Scheme points: calculated (one point per scheme)
  const calcPoints = deduped.map((d: any) => ({
    wly: d.weighted_life_years,
    calc_credibility: d.calculated_credibility,
    scheme_name: d.scheme_name,
    claim_count: d.claim_count,
    experience_period: d.experience_period,
    member_count: d.member_count,
    manual: d.manually_added_credibility
  }))

  // Scheme points: manual (one point per scheme, only where manually set)
  const manualPoints = deduped
    .filter((d: any) => d.manually_added_credibility > 0)
    .map((d: any) => ({
      wly: d.weighted_life_years,
      manual_credibility: d.manually_added_credibility,
      scheme_name: d.scheme_name,
      calc: d.calculated_credibility
    }))

  return {
    title: {
      text: 'Credibility Curve — Weighted Life Years vs Credibility Factor'
    },
    axes: [
      {
        type: 'number',
        position: 'bottom',
        title: { text: 'Weighted Life Years (Exposure)' },
        nice: true
      },
      {
        type: 'number',
        position: 'left',
        title: { text: 'Credibility Factor' },
        min: 0,
        max: 1.05,
        nice: false
      }
    ],
    series: [
      // Theoretical curve
      {
        type: 'line',
        data: curvePoints,
        xKey: 'wly',
        yKey: 'curve_credibility',
        yName: 'Theoretical Curve √(WLY/FCT)',
        stroke: '#90a4ae',
        strokeWidth: 2,
        marker: { enabled: false },
        tooltip: { enabled: false }
      },
      // Calculated scheme points
      {
        type: 'scatter',
        data: calcPoints,
        xKey: 'wly',
        yKey: 'calc_credibility',
        yName: 'Calculated Credibility',
        marker: {
          size: 9,
          fill: '#1976d2',
          stroke: '#0d47a1',
          strokeWidth: 1.5
        },
        tooltip: {
          renderer: (params: any) => {
            const d = params.datum
            const delta =
              d.manual > 0
                ? ` | Manual Δ: ${d.manual - d.calc_credibility >= 0 ? '+' : ''}${(d.manual - d.calc_credibility).toFixed(3)}`
                : ''
            return {
              content: `<b>${d.scheme_name}</b><br/>
                WLY: ${d.wly.toFixed(1)}<br/>
                Calculated: ${d.calc_credibility.toFixed(4)}<br/>
                Claims: ${d.claim_count} | Members: ${d.member_count}<br/>
                Exp. Period: ${d.experience_period.toFixed(1)} yrs${delta}`
            }
          }
        }
      },
      // Manual override points
      {
        type: 'scatter',
        data: manualPoints,
        xKey: 'wly',
        yKey: 'manual_credibility',
        yName: 'Manual Override',
        marker: {
          size: 10,
          fill: '#e65100',
          shape: 'diamond',
          stroke: '#bf360c',
          strokeWidth: 1.5
        },
        tooltip: {
          renderer: (params: any) => {
            const d = params.datum
            const diff = d.manual_credibility - d.calc
            const direction = diff > 0 ? 'Conservative (+)' : 'Aggressive (−)'
            return {
              content: `<b>${d.scheme_name}</b> — Manual Override<br/>
                Manual: ${d.manual_credibility.toFixed(4)}<br/>
                Calculated: ${d.calc.toFixed(4)}<br/>
                Δ = ${diff >= 0 ? '+' : ''}${diff.toFixed(4)} (${direction})`
            }
          }
        }
      }
    ],
    annotations: [
      // Full credibility horizontal line
      {
        type: 'line',
        direction: 'horizontal',
        value: 1.0,
        stroke: '#43a047',
        strokeWidth: 1.5,
        lineDash: [6, 4],
        label: {
          text: 'Full Credibility (1.0)',
          position: 'right',
          fontSize: 11,
          color: '#43a047'
        }
      },
      // Half credibility line
      {
        type: 'line',
        direction: 'horizontal',
        value: 0.5,
        stroke: '#fb8c00',
        strokeWidth: 1.5,
        lineDash: [4, 4],
        label: {
          text: '50% Credibility',
          position: 'right',
          fontSize: 11,
          color: '#fb8c00'
        }
      },
      // Full credibility threshold vertical line
      {
        type: 'line',
        direction: 'vertical',
        value: fct,
        stroke: '#43a047',
        strokeWidth: 1.5,
        lineDash: [6, 4],
        label: {
          text: `FCT = ${fct.toFixed(0)}`,
          position: 'top',
          fontSize: 11,
          color: '#43a047'
        }
      }
    ],
    legend: { enabled: true, position: 'bottom' }
  }
})

// ─── Panel 2: Gap Analysis ────────────────────────────────────────────────────

const gapChartOptions = computed<any>(() => {
  const data = filteredData.value
  if (!data.length) return null

  // Deduplicate by scheme_name — take most recent record per scheme
  const byScheme = new Map<string, any>()
  for (const d of data) {
    const existing = byScheme.get(d.scheme_name)
    if (!existing || d.year > existing.year) byScheme.set(d.scheme_name, d)
  }
  const chartData = Array.from(byScheme.values()).map((d) => ({
    scheme_name: d.scheme_name,
    calculated_credibility: d.calculated_credibility,
    manual_credibility:
      d.manually_added_credibility > 0 ? d.manually_added_credibility : null,
    delta:
      d.manually_added_credibility > 0
        ? parseFloat(
            (d.manually_added_credibility - d.calculated_credibility).toFixed(4)
          )
        : 0
  }))

  return {
    title: { text: 'Calculated vs Manual Credibility per Scheme' },
    data: chartData,
    axes: [
      {
        type: 'category',
        position: 'bottom',
        title: { text: 'Scheme' },
        label: { rotation: -30 }
      },
      {
        type: 'number',
        position: 'left',
        title: { text: 'Credibility Factor' },
        min: 0,
        max: 1.05,
        nice: false
      }
    ],
    series: [
      {
        type: 'bar',
        xKey: 'scheme_name',
        yKey: 'calculated_credibility',
        yName: 'Calculated',
        fill: '#1976d2',
        tooltip: {
          renderer: (params: any) => ({
            content: `<b>${params.datum.scheme_name}</b><br/>Calculated: ${params.datum.calculated_credibility.toFixed(4)}`
          })
        }
      },
      {
        type: 'bar',
        xKey: 'scheme_name',
        yKey: 'manual_credibility',
        yName: 'Manual Override',
        fill: '#e65100',
        tooltip: {
          renderer: (params: any) => {
            const d = params.datum
            if (d.manual_credibility === null)
              return { content: 'No manual override set' }
            const direction = d.delta > 0 ? 'Conservative' : 'Aggressive'
            return {
              content: `<b>${d.scheme_name}</b><br/>
                Manual: ${d.manual_credibility.toFixed(4)}<br/>
                Δ = ${d.delta >= 0 ? '+' : ''}${d.delta} (${direction})`
            }
          }
        }
      }
    ],
    legend: { enabled: true, position: 'bottom' }
  }
})

// ─── Panel 3: Experience CI ───────────────────────────────────────────────────

const ciChartOptions = computed<any>(() => {
  const data = filteredData.value
  if (!data.length) return null

  const rateKey =
    selectedBenefitType.value === 'gla'
      ? 'annual_gla_experience_rate'
      : selectedBenefitType.value === 'ptd'
        ? 'annual_ptd_experience_rate'
        : 'annual_ci_experience_rate'

  const benefitLabel =
    selectedBenefitType.value === 'gla'
      ? 'GLA'
      : selectedBenefitType.value === 'ptd'
        ? 'PTD'
        : 'CI'

  // Deduplicate by scheme_name (latest year)
  const byScheme = new Map<string, any>()
  for (const d of data) {
    const existing = byScheme.get(d.scheme_name)
    if (!existing || d.year > existing.year) byScheme.set(d.scheme_name, d)
  }

  const z = 1.96
  const chartData = Array.from(byScheme.values())
    .filter((d) => d[rateKey] > 0)
    .map((d) => {
      const rate = d[rateKey]
      const exposure = Math.max(d.member_count * d.experience_period, 1)
      const ciHalf =
        d.claim_count > 0
          ? (z * Math.sqrt(d.claim_count)) / exposure
          : rate * 0.5
      return {
        scheme_name: d.scheme_name,
        rate,
        ci_low: Math.max(rate - ciHalf, 0),
        ci_high: rate + ciHalf,
        claim_count: d.claim_count,
        member_count: d.member_count,
        experience_period: d.experience_period
      }
    })

  if (!chartData.length) return null

  // Pooled average rate (simple mean across schemes)
  const pooledRate =
    chartData.reduce((sum, d) => sum + d.rate, 0) / chartData.length

  return {
    title: {
      text: `${benefitLabel} Annual Experience Rate with 95% Poisson Confidence Intervals`
    },
    data: chartData,
    axes: [
      {
        type: 'category',
        position: 'bottom',
        title: { text: 'Scheme' },
        label: { rotation: -30 }
      },
      {
        type: 'number',
        position: 'left',
        title: { text: `Annual ${benefitLabel} Rate` },
        nice: true
      }
    ],
    series: [
      // CI range bars (background)
      {
        type: 'range-bar',
        xKey: 'scheme_name',
        yLowKey: 'ci_low',
        yHighKey: 'ci_high',
        yLowName: '95% CI Low',
        yHighName: '95% CI High',
        fill: '#bbdefb',
        stroke: '#90caf9',
        strokeWidth: 1,
        cornerRadius: 2,
        tooltip: {
          renderer: (params: any) => {
            const d = params.datum
            return {
              content: `<b>${d.scheme_name}</b><br/>
                95% CI: [${d.ci_low.toFixed(5)}, ${d.ci_high.toFixed(5)}]<br/>
                Claims: ${d.claim_count} | Members: ${d.member_count}<br/>
                Period: ${d.experience_period.toFixed(1)} yrs`
            }
          }
        }
      },
      // Actual rate points
      {
        type: 'scatter',
        xKey: 'scheme_name',
        yKey: 'rate',
        yName: `${benefitLabel} Experience Rate`,
        marker: {
          size: 9,
          fill: '#1565c0',
          stroke: '#0d47a1',
          strokeWidth: 1.5
        },
        tooltip: {
          renderer: (params: any) => {
            const d = params.datum
            const width = d.ci_high - d.ci_low
            const reliability =
              width < d.rate * 0.5
                ? 'High'
                : width < d.rate
                  ? 'Moderate'
                  : 'Low'
            return {
              content: `<b>${d.scheme_name}</b><br/>
                Rate: ${d.rate.toFixed(5)}<br/>
                CI Width: ${width.toFixed(5)} (${reliability} reliability)`
            }
          }
        }
      }
    ],
    annotations: [
      {
        type: 'line',
        direction: 'horizontal',
        value: pooledRate,
        stroke: '#6a1b9a',
        strokeWidth: 1.5,
        lineDash: [6, 4],
        label: {
          text: `Pooled avg: ${pooledRate.toFixed(5)}`,
          position: 'right',
          fontSize: 11,
          color: '#6a1b9a'
        }
      }
    ],
    legend: { enabled: true, position: 'bottom' }
  }
})

// ─── Panel 4: Exposure Adequacy ───────────────────────────────────────────────

const exposureChartOptions = computed<any>(() => {
  const data = filteredData.value
  if (!data.length) return null

  const hasFct = data.some((d) => d.full_credibility_threshold > 0)
  if (!hasFct) return null

  // Deduplicate by scheme (latest year)
  const byScheme = new Map<string, any>()
  for (const d of data) {
    const existing = byScheme.get(d.scheme_name)
    if (!existing || d.year > existing.year) byScheme.set(d.scheme_name, d)
  }

  const chartData = Array.from(byScheme.values())
    .filter((d) => d.full_credibility_threshold > 0)
    .map((d) => {
      const pct = Math.min(
        (d.weighted_life_years / d.full_credibility_threshold) * 100,
        110
      )
      return {
        scheme_name: d.scheme_name,
        pct,
        wly: d.weighted_life_years,
        fct: d.full_credibility_threshold,
        remaining: Math.max(
          d.full_credibility_threshold - d.weighted_life_years,
          0
        ),
        credibility: d.calculated_credibility
      }
    })
    .sort((a, b) => b.pct - a.pct)

  if (!chartData.length) return null

  return {
    title: {
      text: 'Exposure Adequacy — Progress Toward Full Credibility Threshold (FCT)'
    },
    data: chartData,
    axes: [
      {
        type: 'number',
        position: 'bottom',
        title: { text: '% of Full Credibility Threshold' },
        min: 0,
        max: 115,
        nice: false
      },
      {
        type: 'category',
        position: 'left',
        title: { text: 'Scheme' }
      }
    ],
    series: [
      {
        type: 'bar',
        direction: 'horizontal',
        xKey: 'scheme_name',
        yKey: 'pct',
        yName: '% of FCT',
        formatter: (params: any) => {
          const pct = params.datum.pct
          return {
            fill:
              pct >= 100
                ? '#2e7d32'
                : pct >= 75
                  ? '#558b2f'
                  : pct >= 25
                    ? '#f57f17'
                    : '#c62828'
          }
        },
        tooltip: {
          renderer: (params: any) => {
            const d = params.datum
            const statusMsg =
              d.pct >= 100
                ? 'Full credibility achieved'
                : `${d.remaining.toFixed(0)} more WLY needed for full credibility`
            return {
              content: `<b>${d.scheme_name}</b><br/>
                WLY: ${d.wly.toFixed(1)} / FCT: ${d.fct.toFixed(0)}<br/>
                Progress: ${d.pct.toFixed(1)}%<br/>
                Credibility: ${d.credibility.toFixed(4)}<br/>
                ${statusMsg}`
            }
          }
        }
      }
    ],
    annotations: [
      {
        type: 'line',
        direction: 'vertical',
        value: 100,
        stroke: '#2e7d32',
        strokeWidth: 2,
        lineDash: [6, 4],
        label: {
          text: 'Full Credibility',
          position: 'top',
          fontSize: 11,
          color: '#2e7d32'
        }
      },
      {
        type: 'line',
        direction: 'vertical',
        value: 25,
        stroke: '#c62828',
        strokeWidth: 1.5,
        lineDash: [4, 4],
        label: {
          text: '25%',
          position: 'top',
          fontSize: 10,
          color: '#c62828'
        }
      },
      {
        type: 'line',
        direction: 'vertical',
        value: 75,
        stroke: '#f57f17',
        strokeWidth: 1.5,
        lineDash: [4, 4],
        label: {
          text: '75%',
          position: 'top',
          fontSize: 10,
          color: '#f57f17'
        }
      }
    ],
    legend: { enabled: false }
  }
})

// ─── Close ────────────────────────────────────────────────────────────────────

function close() {
  dialog.value = false
}
</script>

<style scoped>
.panel-description {
  font-size: 0.8rem;
  color: rgba(0, 0, 0, 0.6);
  margin-bottom: 8px;
  padding: 6px 8px;
  background: rgba(0, 0, 0, 0.03);
  border-left: 3px solid #1976d2;
  border-radius: 0 4px 4px 0;
}
</style>
