<template>
  <div>
    <!-- Concentration KPIs + Pareto -->
    <v-row>
      <v-col cols="12" md="5">
        <v-card variant="outlined" class="pa-3" style="height: 100%">
          <div class="d-flex align-center mb-3">
            <div class="text-subtitle-1 font-weight-bold flex-grow-1">
              Premium Concentration
            </div>
            <v-tooltip location="top" max-width="320">
              <template #activator="{ props: tipProps }">
                <v-icon
                  v-bind="tipProps"
                  size="small"
                  icon="mdi-information-outline"
                  style="cursor: help"
                />
              </template>
              <span>
                Top-N share = % of total annual premium from the N largest
                schemes. HHI (Herfindahl–Hirschman Index) = Σ(scheme share)² on
                a 0–10 000 scale: &lt; 1 500 = low, 1 500–2 500 = moderate, &gt;
                2 500 = highly concentrated book.
              </span>
            </v-tooltip>
          </div>
          <v-row dense>
            <v-col cols="6">
              <div class="text-caption text-medium-emphasis">
                Top-5 premium share
              </div>
              <div class="text-h5 font-weight-bold">
                {{ fmtPct(profile.concentration.top5_premium_share) }}
              </div>
            </v-col>
            <v-col cols="6">
              <div class="text-caption text-medium-emphasis">
                Top-10 premium share
              </div>
              <div class="text-h5 font-weight-bold">
                {{ fmtPct(profile.concentration.top10_premium_share) }}
              </div>
            </v-col>
            <v-col cols="6" class="mt-3">
              <div class="text-caption text-medium-emphasis">
                HHI (0–10 000)
              </div>
              <div
                class="text-h5 font-weight-bold"
                :style="{ color: hhiColor(profile.concentration.hhi) }"
              >
                {{ profile.concentration.hhi.toFixed(0) }}
              </div>
              <div class="text-caption">{{
                hhiLabel(profile.concentration.hhi)
              }}</div>
            </v-col>
            <v-col cols="6" class="mt-3">
              <div class="text-caption text-medium-emphasis">
                Schemes in book
              </div>
              <div class="text-h5 font-weight-bold">
                {{ profile.concentration.total_schemes }}
              </div>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
      <v-col cols="12" md="7">
        <v-card variant="outlined" class="pa-3" style="height: 100%">
          <div class="d-flex align-center mb-2">
            <div class="text-subtitle-1 font-weight-bold flex-grow-1">
              Premium Pareto (cumulative share by scheme rank)
            </div>
            <v-tooltip location="top" max-width="320">
              <template #activator="{ props: tipProps }">
                <v-icon
                  v-bind="tipProps"
                  size="small"
                  icon="mdi-information-outline"
                  class="mr-1"
                  style="cursor: help"
                />
              </template>
              <span>
                Schemes ranked by annual premium (largest first). The line shows
                cumulative share of total premium — a steep early climb means
                the book leans heavily on a few schemes.
              </span>
            </v-tooltip>
            <ChartMenu
              :chart-ref="paretoChart"
              title="Premium Pareto"
              :data="profile.pareto"
            />
          </div>
          <ag-charts
            v-if="paretoOptions"
            ref="paretoChart"
            :options="paretoOptions"
          />
        </v-card>
      </v-col>
    </v-row>

    <!-- Industry / Region heatmap -->
    <v-row class="mt-2">
      <v-col cols="12">
        <v-card variant="outlined" class="pa-3">
          <div class="d-flex align-center mb-2">
            <div class="text-subtitle-1 font-weight-bold flex-grow-1">
              Industry × Region — Loss Ratio Heatmap
            </div>
            <v-tooltip location="top" max-width="360">
              <template #activator="{ props: tipProps }">
                <v-icon
                  v-bind="tipProps"
                  size="small"
                  icon="mdi-information-outline"
                  class="mr-1"
                  style="cursor: help"
                />
              </template>
              <span>
                Each cell groups in-force schemes by Industry (from the quote)
                and Region (from the scheme’s first category). Cell value is ITD
                claims paid ÷ annual premium for that segment; deeper red =
                paying more in claims than the segment is earning.
              </span>
            </v-tooltip>
            <ChartMenu
              title="Industry × Region heatmap"
              :data="profile.industry_region"
            />
          </div>
          <div class="text-caption text-medium-emphasis mb-2">
            Cell colour reflects ITD claims paid ÷ annual premium for schemes in
            that industry &amp; region. Deeper red = paying more in claims than
            the segment is earning.
          </div>
          <div
            v-if="heatmapMatrix.regions.length === 0"
            class="text-medium-emphasis"
          >
            No industry/region data.
          </div>
          <table v-else class="heatmap-table">
            <thead>
              <tr>
                <th></th>
                <th v-for="r in heatmapMatrix.regions" :key="r">{{ r }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="ind in heatmapMatrix.industries" :key="ind">
                <td class="industry-cell">{{ ind }}</td>
                <td
                  v-for="r in heatmapMatrix.regions"
                  :key="r + ind"
                  :style="cellStyle(heatmapMatrix.cells[ind + '||' + r])"
                  :title="cellTooltip(heatmapMatrix.cells[ind + '||' + r])"
                >
                  <span v-if="heatmapMatrix.cells[ind + '||' + r]">
                    {{
                      fmtPct(heatmapMatrix.cells[ind + '||' + r]!.loss_ratio)
                    }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </v-card>
      </v-col>
    </v-row>

    <!-- Frequency / severity scatter -->
    <v-row class="mt-2">
      <v-col cols="12">
        <v-card variant="outlined" class="pa-3">
          <div class="d-flex align-center mb-2">
            <div class="text-subtitle-1 font-weight-bold flex-grow-1">
              Claims Frequency vs Severity (per scheme)
            </div>
            <v-tooltip location="top" max-width="360">
              <template #activator="{ props: tipProps }">
                <v-icon
                  v-bind="tipProps"
                  size="small"
                  icon="mdi-information-outline"
                  class="mr-1"
                  style="cursor: help"
                />
              </template>
              <span>
                Frequency = ITD claims count ÷ members × 1 000 (claims per 1 000
                lives). Severity = avg ITD claim amount. Bubble size = annual
                premium. Top-right bubbles need the closest attention.
              </span>
            </v-tooltip>
            <ChartMenu
              :chart-ref="scatterChart"
              title="Frequency vs Severity"
              :data="profile.frequency_severity"
            />
          </div>
          <div class="text-caption text-medium-emphasis mb-2">
            Each bubble is one in-force scheme. Bubble size scales with annual
            premium. Top-right = high frequency &amp; high severity (most urgent
            risks).
          </div>
          <ag-charts
            v-if="scatterOptions"
            ref="scatterChart"
            :options="scatterOptions"
          />
        </v-card>
      </v-col>
    </v-row>

    <!-- Deteriorating schemes -->
    <v-row class="mt-2">
      <v-col cols="12">
        <v-card variant="outlined" class="pa-3">
          <div class="d-flex align-center mb-2">
            <div class="text-subtitle-1 font-weight-bold">
              Deteriorating Schemes
            </div>
            <v-chip size="small" color="error" variant="tonal" class="ml-2">
              {{ deteriorating.length }}
            </v-chip>
            <v-tooltip location="top" max-width="360">
              <template #activator="{ props: tipProps }">
                <v-icon
                  v-bind="tipProps"
                  size="small"
                  icon="mdi-information-outline"
                  class="ml-2"
                  style="cursor: help"
                />
              </template>
              <span>
                Schemes where ITD ALR &gt; ALR ceiling % OR ITD ALR exceeds
                Expected LR by more than the configured pp delta. Defaults come
                from MetaData → Risk Watchlist Thresholds; tweak the inputs
                below for what-if exploration. Click a row to drill through to
                the scheme.
              </span>
            </v-tooltip>
            <v-spacer />
            <ChartMenu
              :title="`Deteriorating Schemes (ALR>${alrCeiling}%, Δ>${alrDelta}pp)`"
              :data="deteriorating"
            />
          </div>

          <!-- Threshold controls -->
          <v-row dense class="align-center mb-2">
            <v-col cols="12" md="3">
              <v-text-field
                :model-value="alrCeiling"
                label="ALR ceiling (%)"
                type="number"
                min="0"
                max="1000"
                step="1"
                variant="outlined"
                density="compact"
                hide-details
                @update:model-value="onCeilingInput"
              />
            </v-col>
            <v-col cols="12" md="3">
              <v-text-field
                :model-value="alrDelta"
                label="ALR − ELR delta (pp)"
                type="number"
                min="0"
                max="1000"
                step="1"
                variant="outlined"
                density="compact"
                hide-details
                @update:model-value="onDeltaInput"
              />
            </v-col>
            <v-col cols="auto">
              <v-tooltip
                location="top"
                :text="
                  isCustomView
                    ? 'Click to reset to company defaults'
                    : `Company defaults from MetaData: ${profile.thresholds.alr_ceiling_pct}% / ${profile.thresholds.alr_delta_pp}pp`
                "
              >
                <template #activator="{ props: tipProps }">
                  <v-chip
                    v-bind="tipProps"
                    size="small"
                    :color="isCustomView ? 'warning' : 'success'"
                    variant="tonal"
                    :prepend-icon="
                      isCustomView ? 'mdi-restore' : 'mdi-check-circle-outline'
                    "
                    :style="
                      isCustomView ? { cursor: 'pointer' } : { cursor: 'help' }
                    "
                    @click="isCustomView ? emit('reset') : null"
                  >
                    {{
                      isCustomView ? 'Custom view — reset' : 'Company defaults'
                    }}
                  </v-chip>
                </template>
              </v-tooltip>
            </v-col>
            <v-col cols="auto" class="text-caption text-medium-emphasis">
              {{ deteriorating.length }} of {{ rowsWithAlrCount }} scored
              schemes flagged.
              <span v-if="rowsWithAlrCount > 0 && deteriorating.length === 0">
                Highest ALR: {{ fmtPct(maxAlr) }} · max Δ:
                {{ maxDelta == null ? '—' : maxDelta.toFixed(1) + 'pp' }}.
              </span>
            </v-col>
          </v-row>

          <v-table density="compact">
            <thead>
              <tr>
                <th>Scheme</th>
                <th>Trigger</th>
                <th class="text-right">ELR</th>
                <th class="text-right">ALR</th>
                <th class="text-right">Δ</th>
                <th class="text-right">ITD Claims</th>
                <th>Last Claim</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="d in deteriorating"
                :key="d.scheme_id"
                style="cursor: pointer"
                @click="goToScheme(d.scheme_id)"
              >
                <td>{{ d.scheme_name }}</td>
                <td>
                  <v-chip
                    v-for="r in d.trigger_reasons"
                    :key="r"
                    size="x-small"
                    color="error"
                    variant="tonal"
                    class="mr-1"
                  >
                    {{ r }}
                  </v-chip>
                </td>
                <td class="text-right">{{ fmtPct(d.expected_loss_ratio) }}</td>
                <td class="text-right">{{ fmtPct(d.actual_loss_ratio) }}</td>
                <td class="text-right">
                  {{
                    d.loss_ratio_delta == null
                      ? '—'
                      : (d.loss_ratio_delta > 0 ? '+' : '') +
                        d.loss_ratio_delta.toFixed(1) +
                        'pp'
                  }}
                </td>
                <td class="text-right">{{ fmtCurrency(d.itd_claims_paid) }}</td>
                <td>{{ d.last_claim_date || '—' }}</td>
              </tr>
              <tr v-if="deteriorating.length === 0">
                <td colspan="7" class="text-center text-medium-emphasis">
                  No schemes flagged at the current thresholds.
                </td>
              </tr>
            </tbody>
          </v-table>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { AgCharts } from 'ag-charts-vue3'
import ChartMenu from '@/renderer/components/charts/ChartMenu.vue'

const paretoChart: any = ref(null)
const scatterChart: any = ref(null)

interface ConcentrationKPIs {
  top5_premium_share: number
  top10_premium_share: number
  hhi: number
  total_schemes: number
}
interface ParetoPoint {
  rank: number
  scheme_name: string
  premium: number
  cumulative_share: number
}
interface IndustryRegionCell {
  industry: string
  region: string
  premium: number
  claims_paid: number
  loss_ratio: number
  member_count: number
}
interface FreqSeverityPoint {
  scheme_id: number
  scheme_name: string
  frequency: number
  avg_severity: number
  annual_premium: number
}
interface DeterioratingScheme {
  scheme_id: number
  scheme_name: string
  trigger_reasons: string[]
  expected_loss_ratio: number
  actual_loss_ratio: number | null
  loss_ratio_delta: number | null
  itd_claims_paid: number
  last_claim_date?: string
}
interface RiskProfile {
  concentration: ConcentrationKPIs
  pareto: ParetoPoint[]
  industry_region: IndustryRegionCell[]
  frequency_severity: FreqSeverityPoint[]
  deteriorating: DeterioratingScheme[]
}

const props = defineProps<{
  profile: RiskProfile & {
    thresholds: { alr_ceiling_pct: number; alr_delta_pp: number }
  }
  // Re-derived deteriorating list from the parent — kept separate from
  // profile.deteriorating (which is the company-default backend list) so
  // the panel always reflects whatever thresholds are in the inputs below.
  deteriorating: DeterioratingScheme[]
  alrCeiling: number
  alrDelta: number
  isCustomView: boolean
  // All performance rows so we can compute helpful stats (max ALR / max Δ)
  // for the empty-state hint when the watchlist is empty.
  rows: any[]
}>()
const emit = defineEmits<{
  (e: 'update:alrCeiling', value: number): void
  (e: 'update:alrDelta', value: number): void
  (e: 'reset'): void
}>()
const router = useRouter()

// Coerce textfield strings → finite numbers; ignore empty/invalid input so
// the field doesn't get stuck on 0 mid-edit.
const onCeilingInput = (v: string | number | null) => {
  const n = Number(v)
  if (Number.isFinite(n)) emit('update:alrCeiling', n)
}
const onDeltaInput = (v: string | number | null) => {
  const n = Number(v)
  if (Number.isFinite(n)) emit('update:alrDelta', n)
}

const rowsWithAlrCount = computed(
  () => (props.rows ?? []).filter((r) => r.actual_loss_ratio != null).length
)
const maxAlr = computed(() => {
  let m: number | null = null
  for (const r of props.rows ?? []) {
    if (r.actual_loss_ratio == null) continue
    if (m == null || r.actual_loss_ratio > m) m = r.actual_loss_ratio
  }
  return m
})
const maxDelta = computed<number | null>(() => {
  let m: number | null = null
  for (const r of props.rows ?? []) {
    if (r.loss_ratio_delta == null) continue
    if (m == null || r.loss_ratio_delta > m) m = r.loss_ratio_delta
  }
  return m
})

const fmtCurrency = (v: number | null | undefined) =>
  v == null
    ? '—'
    : new Intl.NumberFormat('en-ZA', {
        style: 'currency',
        currency: 'ZAR',
        maximumFractionDigits: 0
      }).format(v)
const fmtPct = (v: number | null | undefined) =>
  v == null ? '—' : `${v.toFixed(1)}%`

const goToScheme = (id: number) =>
  router.push({ name: 'group-pricing-schemes-detail', params: { id } })

const hhiColor = (h: number) => {
  if (h > 2500) return '#c62828'
  if (h > 1500) return '#ef6c00'
  return '#2e7d32'
}
const hhiLabel = (h: number) => {
  if (h > 2500) return 'Highly concentrated'
  if (h > 1500) return 'Moderately concentrated'
  return 'Low concentration'
}

const paretoOptions = computed<any>(() => ({
  data: props.profile.pareto.map((p) => ({
    rank: p.rank,
    scheme: p.scheme_name,
    cumulative: p.cumulative_share * 100
  })),
  background: { fill: 'transparent' },
  height: 280,
  series: [
    {
      type: 'line',
      xKey: 'rank',
      yKey: 'cumulative',
      yName: 'Cumulative premium %',
      tooltip: {
        renderer: (p: any) =>
          `<div style="padding:6px"><b>${p.datum.scheme}</b><br/>Rank ${p.datum.rank}: ${p.datum.cumulative.toFixed(1)}% of GWP</div>`
      }
    }
  ],
  axes: [
    { type: 'number', position: 'bottom', title: { text: 'Scheme rank' } },
    {
      type: 'number',
      position: 'left',
      min: 0,
      max: 100,
      title: { text: 'Cumulative GWP %' }
    }
  ],
  legend: { enabled: false }
}))

interface HeatmapMatrix {
  industries: string[]
  regions: string[]
  cells: Record<string, IndustryRegionCell>
}

const heatmapMatrix = computed<HeatmapMatrix>(() => {
  const industriesSet = new Set<string>()
  const regionsSet = new Set<string>()
  const cells: Record<string, IndustryRegionCell> = {}
  for (const c of props.profile.industry_region) {
    industriesSet.add(c.industry)
    regionsSet.add(c.region)
    cells[c.industry + '||' + c.region] = c
  }
  return {
    industries: Array.from(industriesSet).sort(),
    regions: Array.from(regionsSet).sort(),
    cells
  }
})

const cellStyle = (cell?: IndustryRegionCell) => {
  if (!cell) return { backgroundColor: '#f5f5f5', color: '#bbb' }
  const lr = cell.loss_ratio
  let bg = '#d4edda'
  if (lr > 100) bg = '#ef5350'
  else if (lr > 80) bg = '#ffb74d'
  else if (lr > 50) bg = '#fff59d'
  return {
    backgroundColor: bg,
    fontWeight: '600',
    textAlign: 'center' as const,
    padding: '6px 8px'
  }
}
const cellTooltip = (cell?: IndustryRegionCell) => {
  if (!cell) return ''
  return `${cell.industry} / ${cell.region}\nPremium ${fmtCurrency(cell.premium)}\nClaims ${fmtCurrency(cell.claims_paid)}\nLR ${fmtPct(cell.loss_ratio)}`
}

const scatterOptions = computed<any>(() => ({
  data: props.profile.frequency_severity,
  background: { fill: 'transparent' },
  height: 320,
  series: [
    {
      type: 'bubble',
      xKey: 'frequency',
      xName: 'Frequency (claims per 1000 lives)',
      yKey: 'avg_severity',
      yName: 'Avg severity',
      sizeKey: 'annual_premium',
      sizeName: 'Annual premium',
      labelKey: 'scheme_name',
      labelName: 'Scheme',
      tooltip: {
        renderer: (p: any) =>
          `<div style="padding:6px"><b>${p.datum.scheme_name}</b><br/>Freq: ${p.datum.frequency.toFixed(2)}<br/>Severity: ${fmtCurrency(p.datum.avg_severity)}<br/>Premium: ${fmtCurrency(p.datum.annual_premium)}</div>`
      }
    }
  ],
  axes: [
    {
      type: 'number',
      position: 'bottom',
      title: { text: 'Claims frequency (per 1 000 lives)' }
    },
    {
      type: 'number',
      position: 'left',
      title: { text: 'Average claim severity' }
    }
  ],
  legend: { enabled: false }
}))
</script>

<style scoped>
.heatmap-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
}
.heatmap-table th,
.heatmap-table td {
  border: 1px solid #e0e0e0;
  padding: 6px 8px;
  text-align: center;
}
.heatmap-table th {
  background: #fafafa;
  font-weight: 600;
}
.heatmap-table td.industry-cell {
  text-align: left;
  font-weight: 600;
  background: #fafafa;
}
</style>
