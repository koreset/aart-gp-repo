<template>
  <v-container fluid class="broker-detail pa-6">
    <!-- Hero header -->
    <v-row class="align-center mb-4" no-gutters>
      <v-col cols="auto" class="me-3">
        <v-btn
          variant="tonal"
          size="small"
          color="primary"
          prepend-icon="mdi-arrow-left"
          @click="goBack"
        >
          Back
        </v-btn>
      </v-col>
      <v-col cols="auto" class="me-3">
        <v-avatar size="48" color="primary" class="text-uppercase">
          <span class="text-h6 font-weight-bold">{{ brokerInitials }}</span>
        </v-avatar>
      </v-col>
      <v-col>
        <h2 class="text-h5 font-weight-bold mb-0">
          {{ broker?.name || 'Broker' }}
        </h2>
        <span class="text-caption text-medium-emphasis">
          Performance for {{ year }}
        </span>
      </v-col>
    </v-row>

    <v-row v-if="loadingHeader">
      <v-col><v-progress-linear indeterminate color="primary" /></v-col>
    </v-row>

    <!-- Broker info card -->
    <v-row v-if="broker">
      <v-col cols="12">
        <v-card elevation="1" class="rounded-lg pa-4 broker-info-card">
          <div class="d-flex align-center mb-3 flex-wrap">
            <v-icon color="primary" class="me-2">mdi-account-tie</v-icon>
            <div class="text-subtitle-1 font-weight-bold">Broker details</div>
          </div>
          <v-row dense>
            <v-col
              v-for="field in brokerFields"
              :key="field.label"
              cols="12"
              sm="6"
              md="auto"
              class="flex-grow-1"
            >
              <div class="d-flex align-start info-cell">
                <v-icon
                  size="small"
                  color="grey-darken-1"
                  class="me-2 mt-1"
                >
                  {{ field.icon }}
                </v-icon>
                <div>
                  <div class="text-caption text-medium-emphasis">
                    {{ field.label }}
                  </div>
                  <div class="text-body-2 font-weight-medium">
                    {{ field.value || '—' }}
                  </div>
                </div>
              </div>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>

    <!-- Summary metrics -->
    <v-row v-if="summary" class="mt-2">
      <v-col cols="12">
        <v-card elevation="1" class="rounded-lg pa-4">
          <div class="d-flex align-center mb-3 flex-wrap">
            <v-icon color="primary" class="me-2">mdi-chart-box-outline</v-icon>
            <div class="text-subtitle-1 font-weight-bold">Summary metrics</div>
            <v-spacer />
            <ChartMenu :title="summaryTitle" :data="summaryCsvData" />
          </div>
          <v-row dense>
            <v-col
              v-for="card in summaryCards"
              :key="card.label"
              cols="6"
              md="2"
            >
              <v-card
                variant="flat"
                class="metric-tile pa-3"
                :class="card.tone ? `metric-tile--${card.tone}` : ''"
              >
                <div class="d-flex align-center mb-1">
                  <v-icon
                    size="small"
                    :color="card.iconColor || 'grey-darken-1'"
                    class="me-2"
                  >
                    {{ card.icon }}
                  </v-icon>
                  <span class="text-caption text-medium-emphasis">
                    {{ card.label }}
                  </span>
                  <v-spacer />
                  <v-tooltip location="top" max-width="260">
                    <template #activator="{ props: tipProps }">
                      <v-icon
                        v-bind="tipProps"
                        size="x-small"
                        color="grey"
                        icon="mdi-information-outline"
                        style="cursor: help"
                      />
                    </template>
                    <span>{{ card.tooltip }}</span>
                  </v-tooltip>
                </div>
                <div class="text-h6 font-weight-bold metric-value">
                  {{ card.value }}
                </div>
              </v-card>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>

    <!-- Exposure tabs -->
    <v-row class="mt-2">
      <v-col>
        <v-card elevation="1" class="rounded-lg exposure-card">
          <div class="d-flex align-center px-4 pt-4 flex-wrap">
            <v-icon color="primary" class="me-2">mdi-chart-bar</v-icon>
            <div class="text-subtitle-1 font-weight-bold">
              Exposure breakdown
            </div>
            <v-spacer />
            <ChartMenu
              :chart-ref="exposureChart"
              :title="exposureChartTitle"
              :data="activeTabState.buckets"
            />
          </div>
          <v-tabs v-model="activeTab" color="primary" density="compact">
            <v-tab value="age">
              <v-icon start size="small">mdi-cake-variant-outline</v-icon>
              Age Band
            </v-tab>
            <v-tab value="sum_assured">
              <v-icon start size="small">mdi-cash-multiple</v-icon>
              Sum Assured
            </v-tab>
            <v-tab value="occupation_class">
              <v-icon start size="small">mdi-briefcase-outline</v-icon>
              Occupation Class
            </v-tab>
          </v-tabs>
          <v-divider />
          <v-card-text>
            <div v-if="activeTabState.loading">
              <v-progress-linear indeterminate color="primary" />
            </div>
            <div
              v-else-if="!activeTabState.buckets.length"
              class="text-center py-8"
            >
              <v-icon size="40" color="grey-lighten-1" class="mb-2"
                >mdi-database-off-outline</v-icon
              >
              <div class="text-medium-emphasis">
                No exposure data for this dimension.
              </div>
            </div>
            <div v-else>
              <ag-charts ref="exposureChart" :options="chartOptions" />
              <v-divider class="my-4" />
              <v-table density="compact" hover class="exposure-table">
                <thead>
                  <tr>
                    <th class="text-left">{{ activeTabHeader }}</th>
                    <th class="text-right">{{ countColumnHeader }}</th>
                    <th class="text-right">Total Sum Assured</th>
                    <th class="text-right">Male</th>
                    <th class="text-right">Female</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="b in activeTabState.buckets" :key="b.label">
                    <td>{{ b.label }}</td>
                    <td class="text-right">{{ b.record_count }}</td>
                    <td class="text-right"
                      >R{{ formatNumber(b.total_sum_assured) }}</td
                    >
                    <td class="text-right"
                      >R{{ formatNumber(b.male_sum_assured) }}</td
                    >
                    <td class="text-right"
                      >R{{ formatNumber(b.female_sum_assured) }}</td
                    >
                  </tr>
                </tbody>
              </v-table>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { AgCharts } from 'ag-charts-vue3'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useErrorHandler } from '@/renderer/composables/useErrorHandler'
import { barChartOptions } from '@/renderer/composables/useChartDefaults'
import ChartMenu from '@/renderer/components/charts/ChartMenu.vue'

type Dimension = 'age' | 'sum_assured' | 'occupation_class'

interface ExposureBucket {
  label: string
  sort_order: number
  record_count: number
  total_sum_assured: number
  male_sum_assured: number
  female_sum_assured: number
}

interface BrokerSummary {
  broker_id: number
  year: number
  total_quotes: number
  accepted_quotes: number
  conversion_rate: number
  total_premium: number
  in_force_premium: number
  approved_claims: number
  expected_claims: number
  expected_loss_ratio: number
  actual_loss_ratio: number
}

const props = defineProps<{
  brokerId: string | number
  year: string | number
}>()

const router = useRouter()
const { handleApiError } = useErrorHandler()

const brokerIdNum = computed(() => Number(props.brokerId))
const yearNum = computed(() => Number(props.year))

const broker = ref<any>(null)
const summary = ref<BrokerSummary | null>(null)
const loadingHeader = ref(true)
const exposureChart: any = ref(null)

const tabState = reactive<
  Record<
    Dimension,
    { loaded: boolean; loading: boolean; buckets: ExposureBucket[] }
  >
>({
  age: { loaded: false, loading: false, buckets: [] },
  sum_assured: { loaded: false, loading: false, buckets: [] },
  occupation_class: { loaded: false, loading: false, buckets: [] }
})

const activeTab = ref<Dimension>('age')
const activeTabState = computed(() => tabState[activeTab.value])
const activeTabHeader = computed(() => {
  switch (activeTab.value) {
    case 'age':
      return 'Age band'
    case 'sum_assured':
      return 'Sum assured band'
    case 'occupation_class':
      return 'Occupation class'
    default:
      return ''
  }
})

// Header for the count column. For age/occupation_class the count is sourced
// from member_rating_results so each row equals one member-quote record;
// for sum_assured it's the count of exposure rows in the band.
const countColumnHeader = computed(() =>
  activeTab.value === 'sum_assured' ? 'Records' : 'Members'
)

const formatNumber = (num: number | null | undefined) => {
  if (num == null) return '0'
  if (num >= 1_000_000)
    return (num / 1_000_000).toFixed(1).replace(/\.0$/, '') + 'M'
  if (num >= 1_000) return (num / 1_000).toFixed(1).replace(/\.0$/, '') + 'K'
  return num.toString()
}

const brokerInitials = computed(() => {
  const name = broker.value?.name || ''
  return name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((p: string) => p[0])
    .join('')
    .toUpperCase()
})

const brokerFields = computed(() => {
  if (!broker.value) return []
  return [
    {
      icon: 'mdi-email-outline',
      label: 'Contact email',
      value: broker.value.contact_email
    },
    {
      icon: 'mdi-phone-outline',
      label: 'Contact number',
      value: broker.value.contact_number
    },
    {
      icon: 'mdi-identifier',
      label: 'FSP number',
      value: broker.value.fsp_number
    },
    {
      icon: 'mdi-shield-check-outline',
      label: 'FSP category',
      value: broker.value.fsp_category
    },
    {
      icon: 'mdi-file-document-outline',
      label: 'Binder ref',
      value: broker.value.binder_agreement_ref
    }
  ]
})

// Compare actual vs expected loss ratio to drive a colour cue:
// > 5pp worse than expected → warning, > 5pp better → success.
const lossRatioDelta = computed(() => {
  if (!summary.value) return 0
  return summary.value.actual_loss_ratio - summary.value.expected_loss_ratio
})

const summaryCards = computed(() => {
  if (!summary.value) return []
  const lrTone =
    lossRatioDelta.value > 5
      ? 'danger'
      : lossRatioDelta.value < -5
        ? 'success'
        : ''
  const conversionTone =
    summary.value.conversion_rate >= 30
      ? 'success'
      : summary.value.conversion_rate <= 5
        ? 'danger'
        : ''
  return [
    {
      label: 'Total quotes',
      value: summary.value.total_quotes,
      icon: 'mdi-file-document-multiple-outline',
      iconColor: 'primary',
      tone: '',
      tooltip:
        'Number of quotes raised by this broker that fall within the selected financial year.'
    },
    {
      label: 'Accepted',
      value: summary.value.accepted_quotes,
      icon: 'mdi-check-circle-outline',
      iconColor: 'success',
      tone: '',
      tooltip:
        'Quotes from this broker that were accepted and went on risk during the year.'
    },
    {
      label: 'Conversion',
      value: `${summary.value.conversion_rate}%`,
      icon: 'mdi-percent-outline',
      iconColor: conversionTone || 'primary',
      tone: conversionTone,
      tooltip:
        'Accepted ÷ Total quotes. Higher is better — measures how often this broker’s quotes win business.'
    },
    {
      label: 'Total premium',
      value: `R${formatNumber(summary.value.total_premium)}`,
      icon: 'mdi-currency-usd',
      iconColor: 'primary',
      tone: '',
      tooltip:
        'Sum of annual premium across all of this broker’s quotes for the year (in scheme currency, ZAR-equivalent shown).'
    },
    {
      label: 'Expected LR',
      value: `${summary.value.expected_loss_ratio}%`,
      icon: 'mdi-trending-up',
      iconColor: 'primary',
      tone: '',
      tooltip:
        'Expected loss ratio: predicted claims ÷ annual premium, based on the rating-engine output for this broker’s book.'
    },
    {
      label: 'Actual LR',
      value: `${summary.value.actual_loss_ratio}%`,
      icon: 'mdi-chart-line',
      iconColor: lrTone || 'primary',
      tone: lrTone,
      tooltip:
        'Actual loss ratio: approved claims ÷ in-force premium. Highlights red if it exceeds the expected LR by more than 5pp.'
    }
  ]
})

const summaryCsvData = computed(() => {
  if (!summary.value) return []
  return [
    {
      broker: broker.value?.name ?? '',
      year: summary.value.year,
      total_quotes: summary.value.total_quotes,
      accepted_quotes: summary.value.accepted_quotes,
      conversion_rate: summary.value.conversion_rate,
      total_premium: summary.value.total_premium,
      in_force_premium: summary.value.in_force_premium,
      approved_claims: summary.value.approved_claims,
      expected_claims: summary.value.expected_claims,
      expected_loss_ratio: summary.value.expected_loss_ratio,
      actual_loss_ratio: summary.value.actual_loss_ratio
    }
  ]
})

const summaryTitle = computed(
  () => `${broker.value?.name ?? 'Broker'} - Summary (${yearNum.value})`
)

const exposureChartTitle = computed(() => {
  const brokerName = broker.value?.name ?? 'Broker'
  const dim =
    activeTab.value === 'age'
      ? 'Age Band'
      : activeTab.value === 'sum_assured'
        ? 'Sum Assured'
        : 'Occupation Class'
  return `${brokerName} - ${dim} Exposure (${yearNum.value})`
})

const chartOptions = computed<any>(() =>
  barChartOptions({
    data: activeTabState.value.buckets,
    series: [
      {
        type: 'bar',
        xKey: 'label',
        yKey: 'total_sum_assured',
        yName: 'Total sum assured',
        cornerRadius: 4,
        tooltip: {
          renderer: (params: any) =>
            `<div style="padding: 8px 10px;">
              <div style="font-weight:600;margin-bottom:4px;">${params.datum.label}</div>
              <div>Total: <b>R${formatNumber(params.datum.total_sum_assured)}</b></div>
              <div style="color:#555;">Male: R${formatNumber(params.datum.male_sum_assured)}</div>
              <div style="color:#555;">Female: R${formatNumber(params.datum.female_sum_assured)}</div>
              <div style="color:#888;font-size:11px;margin-top:4px;">${params.datum.record_count} ${activeTab.value === 'sum_assured' ? 'records' : 'members'}</div>
            </div>`
        }
      }
    ],
    axes: [
      { type: 'category', position: 'bottom' },
      {
        type: 'number',
        position: 'left',
        label: { formatter: (p: any) => `R${formatNumber(p.value)}` }
      }
    ],
    legend: { enabled: false },
    height: 320
  })
)

const loadHeader = async () => {
  loadingHeader.value = true
  try {
    const [b, s] = await Promise.all([
      GroupPricingService.getBroker(brokerIdNum.value),
      GroupPricingService.getBrokerPerformanceSummary(
        brokerIdNum.value,
        yearNum.value
      )
    ])
    broker.value = b.data
    summary.value = s.data
  } catch (err) {
    handleApiError(err)
  } finally {
    loadingHeader.value = false
  }
}

const loadDimension = async (dim: Dimension) => {
  if (tabState[dim].loaded || tabState[dim].loading) return
  tabState[dim].loading = true
  try {
    const res = await GroupPricingService.getBrokerExposureBreakdown(
      brokerIdNum.value,
      yearNum.value,
      dim
    )
    tabState[dim].buckets = res.data || []
    tabState[dim].loaded = true
  } catch (err) {
    handleApiError(err)
  } finally {
    tabState[dim].loading = false
  }
}

const goBack = () => {
  router.push({ name: 'group-pricing-dashboard' })
}

watch(activeTab, (dim) => loadDimension(dim), { immediate: false })

onMounted(async () => {
  await loadHeader()
  await loadDimension('age')
})
</script>

<style scoped>
.broker-detail {
  background: rgb(var(--v-theme-background));
}

/* Subtle primary tint so the broker details card sits visually behind the
   metric/chart panels — same brand teal that the chart palette uses. */
.broker-info-card {
  background: rgba(0, 63, 88, 0.06);
  border: 1px solid rgba(0, 63, 88, 0.12);
}

/* Exposure card: same teal family but a touch lighter, so the visual
   hierarchy steps down from the broker header without going pure white. */
.exposure-card {
  background: rgba(0, 63, 88, 0.04);
  border: 1px solid rgba(0, 63, 88, 0.10);
}
.exposure-card :deep(.v-card-text),
.exposure-card :deep(.v-table) {
  background: transparent;
}

.info-cell {
  padding: 4px 12px 4px 0;
  min-width: 180px;
}

.metric-tile {
  border-radius: 10px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: rgb(var(--v-theme-surface));
  transition:
    box-shadow 0.18s ease,
    transform 0.18s ease,
    border-color 0.18s ease;
}
.metric-tile:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
  border-color: rgba(0, 0, 0, 0.18);
  transform: translateY(-1px);
}
.metric-tile--success {
  border-color: rgba(76, 175, 80, 0.4);
  background: rgba(76, 175, 80, 0.06);
}
.metric-tile--success .metric-value {
  color: rgb(var(--v-theme-success));
}
.metric-tile--danger {
  border-color: rgba(244, 67, 54, 0.4);
  background: rgba(244, 67, 54, 0.06);
}
.metric-tile--danger .metric-value {
  color: rgb(var(--v-theme-error));
}

.exposure-table :deep(th) {
  font-weight: 600;
  white-space: nowrap;
  background: rgba(0, 0, 0, 0.02);
}
.exposure-table :deep(tbody tr:hover) {
  background: rgba(0, 63, 88, 0.04);
}
</style>
