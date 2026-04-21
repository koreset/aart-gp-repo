<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center">
              <v-btn
                class="mr-3"
                size="small"
                variant="text"
                prepend-icon="mdi-arrow-left"
                @click="$router.back()"
              >
                Back
              </v-btn>
              <span class="headline">RI Bordereaux KPI Dashboard</span>
            </div>
          </template>

          <template #default>
            <!-- Filter bar -->
            <v-card variant="outlined" class="mb-5">
              <v-card-text>
                <v-row align="center">
                  <v-col cols="12" sm="3">
                    <v-select
                      v-model="filters.treaty_id"
                      :items="treaties"
                      item-title="treaty_name"
                      item-value="id"
                      label="Treaty (all if blank)"
                      variant="outlined"
                      density="compact"
                      clearable
                    />
                  </v-col>
                  <v-col cols="12" sm="3">
                    <v-text-field
                      v-model="filters.period_from"
                      label="Period From"
                      type="date"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                  <v-col cols="12" sm="3">
                    <v-text-field
                      v-model="filters.period_to"
                      label="Period To"
                      type="date"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                  <v-col cols="auto" class="d-flex gap-2 align-center">
                    <v-btn
                      color="primary"
                      :loading="loading"
                      prepend-icon="mdi-calculator"
                      @click="computeKPIs"
                    >
                      Compute
                    </v-btn>
                    <v-btn variant="outlined" @click="resetFilters"
                      >Reset</v-btn
                    >
                  </v-col>
                </v-row>
              </v-card-text>
            </v-card>

            <div
              v-if="!kpis && !loading"
              class="text-center py-10 text-medium-emphasis"
            >
              <v-icon size="64" color="grey-lighten-1">mdi-chart-bar</v-icon>
              <p class="mt-3"
                >Select filters and click <strong>Compute</strong> to generate
                the KPI report.</p
              >
            </div>

            <div v-if="loading" class="text-center py-10">
              <v-progress-circular indeterminate color="primary" size="48" />
              <p class="mt-3 text-medium-emphasis">Computing KPIs…</p>
            </div>

            <template v-if="kpis && !loading">
              <!-- Meta row -->
              <v-row class="mb-2">
                <v-col>
                  <p class="text-caption text-medium-emphasis">
                    Computed
                    {{
                      kpis.computed_at?.replace('T', ' ').replace('Z', ' UTC')
                    }}
                    · Period {{ kpis.period_from }} → {{ kpis.period_to }}
                    <span v-if="kpis.treaty_id">
                      · Treaty ID {{ kpis.treaty_id }}</span
                    >
                  </p>
                </v-col>
              </v-row>

              <!-- KPI Cards: row 1 (percentage KPIs) -->
              <v-row class="mb-4">
                <v-col
                  v-for="kpi in pctKPIs"
                  :key="kpi.key"
                  cols="12"
                  sm="6"
                  md="4"
                  lg=""
                >
                  <v-card variant="outlined" :class="kpiCardClass(kpi)">
                    <v-card-text class="pa-4">
                      <div
                        class="d-flex align-start justify-space-between mb-2"
                      >
                        <div>
                          <p class="text-caption text-medium-emphasis mb-1">{{
                            kpi.label
                          }}</p>
                          <p
                            :class="[
                              'text-h4 font-weight-bold',
                              kpiValueColor(kpi)
                            ]"
                          >
                            {{ kpi.value.toFixed(1) }}%
                          </p>
                        </div>
                        <v-icon :color="kpiIconColor(kpi)" size="36">{{
                          kpi.icon
                        }}</v-icon>
                      </div>
                      <v-progress-linear
                        :model-value="Math.min(kpi.value, 100)"
                        :color="kpiBarColor(kpi)"
                        bg-color="grey-lighten-3"
                        rounded
                        height="6"
                        class="mb-2"
                      />
                      <div class="d-flex justify-space-between">
                        <span class="text-caption text-medium-emphasis"
                          >Target: {{ kpi.target }}%</span
                        >
                        <v-chip
                          :color="kpiMeetsTarget(kpi) ? 'success' : 'error'"
                          size="x-small"
                          variant="tonal"
                        >
                          {{
                            kpiMeetsTarget(kpi)
                              ? '✓ On Track'
                              : '✗ Below Target'
                          }}
                        </v-chip>
                      </div>
                      <p class="text-caption text-medium-emphasis mt-1">{{
                        kpi.detail
                      }}</p>
                    </v-card-text>
                  </v-card>
                </v-col>
              </v-row>

              <!-- KPI Cards: row 2 (non-percentage KPIs) -->
              <v-row class="mb-5">
                <!-- KPI 4: Avg Error Resolution Days -->
                <v-col cols="12" sm="6" md="4">
                  <v-card
                    variant="outlined"
                    :class="
                      kpis.avg_error_resolution_days <=
                      kpis.avg_error_resolution_target
                        ? 'border-success'
                        : 'border-error'
                    "
                  >
                    <v-card-text class="pa-4">
                      <div
                        class="d-flex align-start justify-space-between mb-2"
                      >
                        <div>
                          <p class="text-caption text-medium-emphasis mb-1"
                            >Avg. Error Resolution</p
                          >
                          <p
                            :class="[
                              'text-h4 font-weight-bold',
                              kpis.avg_error_resolution_days <=
                              kpis.avg_error_resolution_target
                                ? 'text-success'
                                : 'text-error'
                            ]"
                          >
                            {{ kpis.avg_error_resolution_days.toFixed(1) }}
                            <span class="text-h6 font-weight-regular"
                              >days</span
                            >
                          </p>
                        </div>
                        <v-icon
                          :color="
                            kpis.avg_error_resolution_days <=
                            kpis.avg_error_resolution_target
                              ? 'success'
                              : 'error'
                          "
                          size="36"
                        >
                          mdi-timer-outline
                        </v-icon>
                      </div>
                      <div class="d-flex justify-space-between">
                        <span class="text-caption text-medium-emphasis"
                          >Target: ≤{{
                            kpis.avg_error_resolution_target
                          }}
                          days</span
                        >
                        <v-chip
                          :color="
                            kpis.avg_error_resolution_days <=
                            kpis.avg_error_resolution_target
                              ? 'success'
                              : 'error'
                          "
                          size="x-small"
                          variant="tonal"
                        >
                          {{
                            kpis.avg_error_resolution_days <=
                            kpis.avg_error_resolution_target
                              ? '✓ On Track'
                              : '✗ Slow'
                          }}
                        </v-chip>
                      </div>
                      <p class="text-caption text-medium-emphasis mt-1"
                        >{{ kpis.error_resolution_samples }} validated run(s)
                        sampled</p
                      >
                    </v-card-text>
                  </v-card>
                </v-col>

                <!-- KPI 7: Open Query Backlog -->
                <v-col cols="12" sm="6" md="4">
                  <v-card
                    variant="outlined"
                    :class="
                      kpis.open_query_backlog < kpis.open_query_backlog_target
                        ? 'border-success'
                        : 'border-error'
                    "
                  >
                    <v-card-text class="pa-4">
                      <div
                        class="d-flex align-start justify-space-between mb-2"
                      >
                        <div>
                          <p class="text-caption text-medium-emphasis mb-1"
                            >Open Query Backlog</p
                          >
                          <p
                            :class="[
                              'text-h4 font-weight-bold',
                              kpis.open_query_backlog <
                              kpis.open_query_backlog_target
                                ? 'text-success'
                                : 'text-error'
                            ]"
                          >
                            {{ kpis.open_query_backlog }}
                            <span class="text-h6 font-weight-regular"
                              >runs</span
                            >
                          </p>
                        </div>
                        <v-icon
                          :color="
                            kpis.open_query_backlog <
                            kpis.open_query_backlog_target
                              ? 'success'
                              : 'error'
                          "
                          size="36"
                        >
                          mdi-alert-circle-outline
                        </v-icon>
                      </div>
                      <div class="d-flex justify-space-between">
                        <span class="text-caption text-medium-emphasis"
                          >Target: &lt;{{ kpis.open_query_backlog_target }} per
                          quarter</span
                        >
                        <v-chip
                          :color="
                            kpis.open_query_backlog <
                            kpis.open_query_backlog_target
                              ? 'success'
                              : 'error'
                          "
                          size="x-small"
                          variant="tonal"
                        >
                          {{
                            kpis.open_query_backlog <
                            kpis.open_query_backlog_target
                              ? '✓ On Track'
                              : '✗ Backlog'
                          }}
                        </v-chip>
                      </div>
                      <p class="text-caption text-medium-emphasis mt-1"
                        >Runs in validation_failed &gt;30 days</p
                      >
                    </v-card-text>
                  </v-card>
                </v-col>

                <!-- Overall compliance score -->
                <v-col cols="12" sm="6" md="4">
                  <v-card
                    variant="outlined"
                    class="h-100"
                    style="
                      background: linear-gradient(
                        135deg,
                        #e8f5e9 0%,
                        #f3e5f5 100%
                      );
                    "
                  >
                    <v-card-text
                      class="pa-4 d-flex flex-column align-center justify-center"
                      style="min-height: 130px"
                    >
                      <p class="text-caption text-medium-emphasis mb-1"
                        >Overall Compliance Score</p
                      >
                      <p
                        :class="[
                          'text-h3 font-weight-bold',
                          overallScore >= 80
                            ? 'text-success'
                            : overallScore >= 60
                              ? 'text-warning'
                              : 'text-error'
                        ]"
                      >
                        {{ overallScore }}/7
                      </p>
                      <p class="text-caption text-medium-emphasis"
                        >KPIs meeting target</p
                      >
                      <v-chip
                        :color="
                          overallScore >= 6
                            ? 'success'
                            : overallScore >= 4
                              ? 'warning'
                              : 'error'
                        "
                        variant="tonal"
                        class="mt-2"
                      >
                        {{
                          overallScore >= 6
                            ? 'Compliant'
                            : overallScore >= 4
                              ? 'Partial'
                              : 'Non-Compliant'
                        }}
                      </v-chip>
                    </v-card-text>
                  </v-card>
                </v-col>
              </v-row>

              <!-- AG Charts: Actual vs Target bar chart for % KPIs -->
              <v-card variant="outlined" class="mb-4">
                <v-card-title class="text-subtitle-1 font-weight-bold pa-4">
                  KPI Performance vs Target — Percentage Metrics
                </v-card-title>
                <v-card-text>
                  <ag-charts
                    :options="chartOptions"
                    style="height: 320px; width: 100%"
                  />
                </v-card-text>
              </v-card>

              <!-- Management report table -->
              <v-card variant="outlined">
                <v-card-title
                  class="text-subtitle-1 font-weight-bold pa-4 d-flex justify-space-between"
                >
                  Management Report Summary
                  <v-btn
                    size="small"
                    variant="outlined"
                    prepend-icon="mdi-download"
                    @click="exportCSV"
                  >
                    Export CSV
                  </v-btn>
                </v-card-title>
                <v-card-text class="pa-0">
                  <v-table density="compact">
                    <thead>
                      <tr>
                        <th>KPI</th>
                        <th>Actual</th>
                        <th>Target</th>
                        <th>Sample</th>
                        <th>Status</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="row in reportRows" :key="row.name">
                        <td>{{ row.name }}</td>
                        <td :class="row.pass ? 'text-success' : 'text-error'">
                          <strong>{{ row.actual }}</strong>
                        </td>
                        <td>{{ row.target }}</td>
                        <td class="text-medium-emphasis">{{ row.sample }}</td>
                        <td>
                          <v-chip
                            :color="row.pass ? 'success' : 'error'"
                            size="x-small"
                            variant="tonal"
                          >
                            {{ row.pass ? 'Pass' : 'Fail' }}
                          </v-chip>
                        </td>
                      </tr>
                    </tbody>
                  </v-table>
                </v-card-text>
              </v-card>
            </template>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3500">
      {{ snackbar.message }}
    </v-snackbar>
  </v-container>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { AgCharts } from 'ag-charts-vue3'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import BaseCard from '@/renderer/components/BaseCard.vue'

const loading = ref(false)
const kpis = ref(null)
const treaties = ref([])
const snackbar = ref({ show: false, message: '', color: 'success' })

const now = new Date()
const firstOfYear = `${now.getFullYear()}-01-01`
const today = now.toISOString().split('T')[0]

const filters = ref({
  treaty_id: null,
  period_from: firstOfYear,
  period_to: today
})

const notify = (message, color = 'success') => {
  snackbar.value = { show: true, message, color }
}

// ── KPI card config ────────────────────────────────────────────────────────
const pctKPIs = computed(() => {
  if (!kpis.value) return []
  const k = kpis.value
  return [
    {
      key: 'submission',
      label: 'Submission Timeliness',
      value: k.submission_timeliness_pct,
      target: k.submission_timeliness_target,
      detail: `${k.submission_on_time} of ${k.submission_total} submitted within 10 days of period end`,
      icon: 'mdi-send-clock'
    },
    {
      key: 'processing',
      label: 'Processing Timeliness',
      value: k.processing_timeliness_pct,
      target: k.processing_timeliness_target,
      detail: `${k.processing_on_time} of ${k.processing_total} acknowledged within 10 days of submission`,
      icon: 'mdi-clock-fast'
    },
    {
      key: 'firsttime',
      label: 'First-Time Acceptance',
      value: k.first_time_acceptance_pct,
      target: k.first_time_acceptance_target,
      detail: `${k.first_time_accepted} of ${k.first_time_total} runs required no amendment`,
      icon: 'mdi-check-decagram'
    },
    {
      key: 'settlement',
      label: 'Settlement Timeliness',
      value: k.settlement_timeliness_pct,
      target: k.settlement_timeliness_target,
      detail: `${k.settlement_on_time} of ${k.settlement_total} accounts settled within 30 days of agreement`,
      icon: 'mdi-bank-check'
    },
    {
      key: 'claims',
      label: 'Claims Register Completeness',
      value: k.claims_completeness_pct,
      target: k.claims_completeness_target,
      detail: `${k.claims_with_notice} of ${k.claims_large_loss_total} large-loss claims have a notification`,
      icon: 'mdi-clipboard-check-outline'
    }
  ]
})

const kpiMeetsTarget = (kpi) => kpi.value >= kpi.target
const kpiValueColor = (kpi) =>
  kpiMeetsTarget(kpi) ? 'text-success' : 'text-error'
const kpiIconColor = (kpi) => (kpiMeetsTarget(kpi) ? 'success' : 'error')
const kpiBarColor = (kpi) => (kpiMeetsTarget(kpi) ? 'success' : 'error')
const kpiCardClass = (kpi) =>
  kpiMeetsTarget(kpi) ? 'border-success h-100' : 'border-error h-100'

const overallScore = computed(() => {
  if (!kpis.value) return 0
  const k = kpis.value
  let score = 0
  if (k.submission_timeliness_pct >= k.submission_timeliness_target) score++
  if (k.processing_timeliness_pct >= k.processing_timeliness_target) score++
  if (k.first_time_acceptance_pct >= k.first_time_acceptance_target) score++
  if (k.avg_error_resolution_days <= k.avg_error_resolution_target) score++
  if (k.settlement_timeliness_pct >= k.settlement_timeliness_target) score++
  if (k.claims_completeness_pct >= k.claims_completeness_target) score++
  if (k.open_query_backlog < k.open_query_backlog_target) score++
  return score
})

// ── AG Charts config ───────────────────────────────────────────────────────
const chartOptions = computed(() => {
  if (!kpis.value) return {}
  const labels = [
    'Submission\nTimeliness',
    'Processing\nTimeliness',
    'First-Time\nAcceptance',
    'Settlement\nTimeliness',
    'Claims\nCompleteness'
  ]
  const actuals = [
    kpis.value.submission_timeliness_pct,
    kpis.value.processing_timeliness_pct,
    kpis.value.first_time_acceptance_pct,
    kpis.value.settlement_timeliness_pct,
    kpis.value.claims_completeness_pct
  ]
  const targets = [
    kpis.value.submission_timeliness_target,
    kpis.value.processing_timeliness_target,
    kpis.value.first_time_acceptance_target,
    kpis.value.settlement_timeliness_target,
    kpis.value.claims_completeness_target
  ]
  const data = labels.map((label, i) => ({
    kpi: label,
    actual: actuals[i],
    target: targets[i]
  }))
  return {
    data,
    series: [
      {
        type: 'bar',
        xKey: 'kpi',
        yKey: 'actual',
        yName: 'Actual %',
        fill: '#1976d2',
        tooltip: {
          renderer: (p) => ({
            content: `${p.datum.kpi}: ${p.datum.actual.toFixed(1)}%`
          })
        }
      },
      {
        type: 'bar',
        xKey: 'kpi',
        yKey: 'target',
        yName: 'Target %',
        fill: '#e0e0e0',
        stroke: '#9e9e9e',
        tooltip: {
          renderer: (p) => ({ content: `Target: ${p.datum.target}%` })
        }
      }
    ],
    axes: [
      { type: 'category', position: 'bottom' },
      {
        type: 'number',
        position: 'left',
        min: 0,
        max: 100,
        label: { formatter: (p) => `${p.value}%` }
      }
    ],
    legend: { enabled: true, position: 'bottom' }
  }
})

// ── Management report rows ─────────────────────────────────────────────────
const reportRows = computed(() => {
  if (!kpis.value) return []
  const k = kpis.value
  return [
    {
      name: 'KPI 1 — Submission Timeliness',
      actual: `${k.submission_timeliness_pct.toFixed(1)}%`,
      target: `≥${k.submission_timeliness_target}%`,
      sample: `${k.submission_total} runs`,
      pass: k.submission_timeliness_pct >= k.submission_timeliness_target
    },
    {
      name: 'KPI 2 — Processing Timeliness',
      actual: `${k.processing_timeliness_pct.toFixed(1)}%`,
      target: `≥${k.processing_timeliness_target}%`,
      sample: `${k.processing_total} runs`,
      pass: k.processing_timeliness_pct >= k.processing_timeliness_target
    },
    {
      name: 'KPI 3 — First-Time Acceptance Rate',
      actual: `${k.first_time_acceptance_pct.toFixed(1)}%`,
      target: `≥${k.first_time_acceptance_target}%`,
      sample: `${k.first_time_total} concluded runs`,
      pass: k.first_time_acceptance_pct >= k.first_time_acceptance_target
    },
    {
      name: 'KPI 4 — Avg Error Resolution',
      actual: `${k.avg_error_resolution_days.toFixed(1)} days`,
      target: `≤${k.avg_error_resolution_target} days`,
      sample: `${k.error_resolution_samples} validated runs`,
      pass: k.avg_error_resolution_days <= k.avg_error_resolution_target
    },
    {
      name: 'KPI 5 — Settlement Timeliness',
      actual: `${k.settlement_timeliness_pct.toFixed(1)}%`,
      target: `≥${k.settlement_timeliness_target}%`,
      sample: `${k.settlement_total} settled accounts`,
      pass: k.settlement_timeliness_pct >= k.settlement_timeliness_target
    },
    {
      name: 'KPI 6 — Claims Register Completeness',
      actual: `${k.claims_completeness_pct.toFixed(1)}%`,
      target: `≥${k.claims_completeness_target}%`,
      sample: `${k.claims_large_loss_total} large-loss claims`,
      pass: k.claims_completeness_pct >= k.claims_completeness_target
    },
    {
      name: 'KPI 7 — Open Query Backlog',
      actual: `${k.open_query_backlog} runs`,
      target: `<${k.open_query_backlog_target} per quarter`,
      sample: '>30 days in validation_failed',
      pass: k.open_query_backlog < k.open_query_backlog_target
    }
  ]
})

// ── Actions ────────────────────────────────────────────────────────────────
async function computeKPIs() {
  loading.value = true
  kpis.value = null
  try {
    const params = {}
    if (filters.value.treaty_id) params.treaty_id = filters.value.treaty_id
    if (filters.value.period_from)
      params.period_from = filters.value.period_from
    if (filters.value.period_to) params.period_to = filters.value.period_to
    const res = await GroupPricingService.getRIBordereauxKPIs(params)
    kpis.value = res.data?.data
  } catch (e) {
    notify(e.response?.data?.error || 'Failed to compute KPIs', 'error')
  } finally {
    loading.value = false
  }
}

function resetFilters() {
  filters.value = {
    treaty_id: null,
    period_from: firstOfYear,
    period_to: today
  }
  kpis.value = null
}

function exportCSV() {
  if (!reportRows.value.length) return
  const header = 'KPI,Actual,Target,Sample,Status'
  const rows = reportRows.value.map(
    (r) =>
      `"${r.name}","${r.actual}","${r.target}","${r.sample}","${r.pass ? 'Pass' : 'Fail'}"`
  )
  const csv = [header, ...rows].join('\n')
  const blob = new Blob([csv], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `ri_kpi_report_${kpis.value?.period_from}_${kpis.value?.period_to}.csv`
  a.click()
  URL.revokeObjectURL(url)
}

async function loadTreaties() {
  try {
    const res = await GroupPricingService.getTreaties({ status: 'active' })
    treaties.value = res.data?.data || []
  } catch {}
}

onMounted(() => {
  loadTreaties()
})
</script>

<style scoped>
.border-success {
  border-color: rgb(var(--v-theme-success)) !important;
}
.border-error {
  border-color: rgb(var(--v-theme-error)) !important;
}
</style>
