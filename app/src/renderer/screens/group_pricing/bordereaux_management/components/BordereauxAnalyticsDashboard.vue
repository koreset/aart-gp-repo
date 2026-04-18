<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center justify-space-between">
              <div>
                <span class="headline">Bordereaux Analytics</span>
                <p class="text-subtitle-1 text-medium-emphasis mt-2">
                  Live metrics aggregated from generated bordereaux,
                  confirmations, and group-scheme claims.
                </p>
              </div>
              <div class="d-flex align-center gap-2">
                <v-select
                  v-model="selectedPeriod"
                  :items="periodOptions"
                  label="Period"
                  variant="outlined"
                  density="compact"
                  hide-details
                  @update:model-value="loadAnalytics"
                />
                <v-btn
                  color="grey"
                  variant="outlined"
                  prepend-icon="mdi-arrow-left"
                  @click="$router.push('/group-pricing/bordereaux-management')"
                >
                  Back to Dashboard
                </v-btn>
              </div>
            </div>
          </template>
          <template #default>
            <v-alert
              v-if="errorMessage"
              type="error"
              variant="tonal"
              class="mb-4"
              closable
              @click:close="errorMessage = ''"
            >
              {{ errorMessage }}
            </v-alert>

            <!-- Top KPI cards -->
            <v-row class="mb-6">
              <v-col cols="12" sm="6" lg="3">
                <v-card variant="outlined" class="h-100">
                  <v-card-text>
                    <div class="d-flex align-center justify-space-between">
                      <div>
                        <p class="text-caption text-medium-emphasis mb-1">
                          Total Premium Volume
                        </p>
                        <p class="text-h4 font-weight-bold text-primary">
                          {{ formatCurrency(topMetrics.total_premium_volume) }}
                        </p>
                        <p class="text-caption text-medium-emphasis">
                          {{ periodLabel }}
                        </p>
                      </div>
                      <v-icon size="50" color="primary"
                        >mdi-currency-usd</v-icon
                      >
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card variant="outlined" class="h-100">
                  <v-card-text>
                    <div class="d-flex align-center justify-space-between">
                      <div>
                        <p class="text-caption text-medium-emphasis mb-1">
                          Loss Ratio
                        </p>
                        <p
                          class="text-h4 font-weight-bold"
                          :class="`text-${getLossRatioColor(topMetrics.loss_ratio)}`"
                        >
                          {{ topMetrics.loss_ratio.toFixed(1) }}%
                        </p>
                        <p class="text-caption text-medium-emphasis">
                          claims / premium, {{ periodLabel }}
                        </p>
                      </div>
                      <v-icon size="50" color="success"
                        >mdi-scale-balance</v-icon
                      >
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card variant="outlined" class="h-100">
                  <v-card-text>
                    <div class="d-flex align-center justify-space-between">
                      <div>
                        <p class="text-caption text-medium-emphasis mb-1">
                          Total Claims
                        </p>
                        <p class="text-h4 font-weight-bold text-info">
                          {{ topMetrics.claims_count.toLocaleString() }}
                        </p>
                        <p class="text-caption text-medium-emphasis">
                          {{ formatCurrency(topMetrics.total_claim_volume) }}
                        </p>
                      </div>
                      <v-icon size="50" color="info">
                        mdi-chart-timeline-variant
                      </v-icon>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card variant="outlined" class="h-100">
                  <v-card-text>
                    <div class="d-flex align-center justify-space-between">
                      <div>
                        <p class="text-caption text-medium-emphasis mb-1">
                          Avg Claim Size
                        </p>
                        <p class="text-h4 font-weight-bold text-purple">
                          {{ formatCurrency(topMetrics.avg_claim_size) }}
                        </p>
                        <p class="text-caption text-medium-emphasis">
                          {{ topMetrics.bordereaux_generated }} bordereaux
                          generated
                        </p>
                      </div>
                      <v-icon size="50" color="purple">mdi-calculator</v-icon>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Monthly trend chart -->
            <v-row class="mb-6">
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title class="text-h6 font-weight-bold">
                    Monthly Premium vs Claim Volume
                  </v-card-title>
                  <v-card-text>
                    <div
                      v-if="loading"
                      class="d-flex align-center justify-center"
                      style="height: 320px"
                    >
                      <v-progress-circular indeterminate color="primary" />
                    </div>
                    <ag-charts
                      v-else-if="monthlyChartOptions"
                      :options="monthlyChartOptions"
                      style="height: 320px; width: 100%"
                    />
                    <div
                      v-else
                      class="d-flex align-center justify-center bg-grey-lighten-5 rounded"
                      style="height: 320px"
                    >
                      <div class="text-center">
                        <v-icon size="60" color="grey">mdi-chart-line</v-icon>
                        <p class="text-grey mt-2"> No data in this period </p>
                      </div>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Benefit Performance -->
            <v-row class="mb-6">
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title class="text-h6 font-weight-bold">
                    Benefit Performance
                  </v-card-title>
                  <v-card-text>
                    <v-data-table
                      :headers="benefitHeaders"
                      :items="benefitPerformance"
                      :items-per-page="10"
                      :loading="loading"
                      density="compact"
                      no-data-text="No benefit data in this period"
                    >
                      <template #[`item.benefit`]="{ item }">
                        <div class="d-flex align-center">
                          <v-icon
                            :color="getBenefitIconColor(item.benefit)"
                            class="me-2"
                          >
                            {{ getBenefitIcon(item.benefit) }}
                          </v-icon>
                          {{ item.benefit }}
                        </div>
                      </template>

                      <template #[`item.claim_count`]="{ item }">
                        {{ item.claim_count.toLocaleString() }}
                      </template>

                      <template #[`item.claim_volume`]="{ item }">
                        {{ formatCurrency(item.claim_volume) }}
                      </template>

                      <template #[`item.avg_claim_amount`]="{ item }">
                        {{ formatCurrency(item.avg_claim_amount) }}
                      </template>

                      <template #[`item.approval_rate`]="{ item }">
                        <div class="d-flex align-center">
                          <v-progress-linear
                            :model-value="item.approval_rate"
                            color="success"
                            height="6"
                            rounded
                            class="me-2"
                            style="width: 60px"
                          />
                          <span>{{ item.approval_rate.toFixed(1) }}%</span>
                        </div>
                      </template>
                    </v-data-table>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Insurer Performance -->
            <v-row class="mb-6">
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title class="text-h6 font-weight-bold">
                    Insurer Performance
                  </v-card-title>
                  <v-card-text>
                    <v-data-table
                      :headers="insurerHeaders"
                      :items="insurerMetrics"
                      :items-per-page="10"
                      :loading="loading"
                      density="compact"
                      no-data-text="No insurer activity in this period"
                    >
                      <template #[`item.insurer_name`]="{ item }">
                        <div class="d-flex align-center">
                          <v-avatar size="32" class="me-3" color="primary">
                            {{ item.insurer_name.charAt(0) }}
                          </v-avatar>
                          {{ item.insurer_name }}
                        </div>
                      </template>

                      <template #[`item.market_share_pct`]="{ item }">
                        <div class="d-flex align-center">
                          <v-progress-linear
                            :model-value="item.market_share_pct"
                            color="primary"
                            height="6"
                            rounded
                            class="me-2"
                            style="width: 60px"
                          />
                          <span>{{ item.market_share_pct.toFixed(1) }}%</span>
                        </div>
                      </template>

                      <template #[`item.bordereaux_count`]="{ item }">
                        {{ item.bordereaux_count.toLocaleString() }}
                      </template>

                      <template #[`item.record_count`]="{ item }">
                        {{ item.record_count.toLocaleString() }}
                      </template>

                      <template #[`item.avg_match_score`]="{ item }">
                        <v-chip
                          :color="getMatchScoreColor(item.avg_match_score)"
                          size="small"
                        >
                          {{ item.avg_match_score.toFixed(1) }}%
                        </v-chip>
                      </template>

                      <template #[`item.discrepancy_count`]="{ item }">
                        {{ item.discrepancy_count.toLocaleString() }}
                      </template>
                    </v-data-table>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Monthly KPI Performance -->
            <v-row class="mb-6">
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title class="text-h6 font-weight-bold">
                    Monthly KPI Performance
                  </v-card-title>
                  <v-card-text>
                    <v-data-table
                      :headers="kpiHeaders"
                      :items="monthlyKPIs"
                      :items-per-page="12"
                      :loading="loading"
                      density="compact"
                      no-data-text="No monthly data in this period"
                    >
                      <template #[`item.month`]="{ item }">
                        <span class="font-weight-bold">{{ item.month }}</span>
                      </template>

                      <template #[`item.bordereaux_count`]="{ item }">
                        {{ item.bordereaux_count.toLocaleString() }}
                      </template>

                      <template #[`item.premium_volume`]="{ item }">
                        {{ formatCurrency(item.premium_volume) }}
                      </template>

                      <template #[`item.claim_volume`]="{ item }">
                        {{ formatCurrency(item.claim_volume) }}
                      </template>

                      <template #[`item.avg_match_score`]="{ item }">
                        <v-chip
                          :color="getMatchScoreColor(item.avg_match_score)"
                          size="small"
                        >
                          {{ item.avg_match_score.toFixed(1) }}%
                        </v-chip>
                      </template>

                      <template #[`item.discrepancy_count`]="{ item }">
                        {{ item.discrepancy_count.toLocaleString() }}
                      </template>
                    </v-data-table>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { AgCharts } from 'ag-charts-vue3'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useFlashStore } from '@/renderer/store/flash'

interface TopMetrics {
  total_premium_volume: number
  total_claim_volume: number
  loss_ratio: number
  claims_count: number
  avg_claim_size: number
  bordereaux_generated: number
  avg_match_score: number
}

interface BenefitRow {
  benefit: string
  claim_count: number
  claim_volume: number
  avg_claim_amount: number
  approved_count: number
  approval_rate: number
}

interface InsurerRow {
  insurer_name: string
  bordereaux_count: number
  record_count: number
  market_share_pct: number
  avg_match_score: number
  discrepancy_count: number
}

interface MonthlyRow {
  month: string
  bordereaux_count: number
  premium_volume: number
  claim_volume: number
  avg_match_score: number
  discrepancy_count: number
}

const flash = useFlashStore()
const selectedPeriod = ref('last_30_days')
const loading = ref(false)
const errorMessage = ref('')

const topMetrics = ref<TopMetrics>({
  total_premium_volume: 0,
  total_claim_volume: 0,
  loss_ratio: 0,
  claims_count: 0,
  avg_claim_size: 0,
  bordereaux_generated: 0,
  avg_match_score: 0
})
const benefitPerformance = ref<BenefitRow[]>([])
const insurerMetrics = ref<InsurerRow[]>([])
const monthlyKPIs = ref<MonthlyRow[]>([])

const periodOptions = [
  { title: 'Last 7 Days', value: 'last_7_days' },
  { title: 'Last 30 Days', value: 'last_30_days' },
  { title: 'Last Quarter', value: 'last_quarter' },
  { title: 'Last Year', value: 'last_year' },
  { title: 'Year to Date', value: 'ytd' }
]

const periodLabel = computed(() => {
  return (
    periodOptions.find((o) => o.value === selectedPeriod.value)?.title ||
    selectedPeriod.value
  )
})

const benefitHeaders = [
  { title: 'Benefit', key: 'benefit', sortable: true },
  { title: 'Claim Count', key: 'claim_count', sortable: true },
  { title: 'Claim Volume', key: 'claim_volume', sortable: true },
  { title: 'Avg Claim', key: 'avg_claim_amount', sortable: true },
  { title: 'Approval Rate', key: 'approval_rate', sortable: true }
]

const insurerHeaders = [
  { title: 'Insurer', key: 'insurer_name', sortable: true },
  { title: 'Market Share', key: 'market_share_pct', sortable: true },
  { title: 'Bordereaux', key: 'bordereaux_count', sortable: true },
  { title: 'Records', key: 'record_count', sortable: true },
  { title: 'Avg Match Score', key: 'avg_match_score', sortable: true },
  { title: 'Discrepancies', key: 'discrepancy_count', sortable: true }
]

const kpiHeaders = [
  { title: 'Month', key: 'month', sortable: true },
  { title: 'Bordereaux', key: 'bordereaux_count', sortable: true },
  { title: 'Premium Volume', key: 'premium_volume', sortable: true },
  { title: 'Claim Volume', key: 'claim_volume', sortable: true },
  { title: 'Avg Match Score', key: 'avg_match_score', sortable: true },
  { title: 'Discrepancies', key: 'discrepancy_count', sortable: true }
]

const monthlyChartOptions = computed<any>(() => {
  if (!monthlyKPIs.value.length) return null
  return {
    title: { enabled: false, text: '' },
    data: monthlyKPIs.value,
    axes: [
      { type: 'category', position: 'bottom', keys: ['month'] },
      {
        type: 'number',
        position: 'left',
        title: { text: 'ZAR' },
        label: {
          formatter: (params: any) => formatCompactCurrency(params.value)
        }
      }
    ],
    series: [
      {
        type: 'line',
        xKey: 'month',
        yKey: 'premium_volume',
        yName: 'Premium Volume',
        stroke: '#1976d2',
        marker: { enabled: true, fill: '#1976d2' }
      },
      {
        type: 'line',
        xKey: 'month',
        yKey: 'claim_volume',
        yName: 'Claim Volume',
        stroke: '#d32f2f',
        marker: { enabled: true, fill: '#d32f2f' }
      }
    ],
    legend: { enabled: true, position: 'bottom' }
  }
})

const formatCurrency = (amount: number): string => {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(amount || 0)
}

const formatCompactCurrency = (amount: number): string => {
  const abs = Math.abs(amount)
  if (abs >= 1_000_000) return `R${(amount / 1_000_000).toFixed(1)}M`
  if (abs >= 1_000) return `R${(amount / 1_000).toFixed(0)}k`
  return `R${amount}`
}

const getBenefitIcon = (benefit: string): string => {
  const icons: Record<string, string> = {
    GLA: 'mdi-shield-account',
    CI: 'mdi-medical-bag',
    PTD: 'mdi-wheelchair-accessibility',
    TTD: 'mdi-account-injury',
    SGLA: 'mdi-shield-account-outline',
    PHI: 'mdi-hospital-building',
    FuneralMM: 'mdi-flower',
    FuneralSP: 'mdi-flower-outline',
    FuneralCH: 'mdi-flower-outline',
    FuneralPAR: 'mdi-flower-outline',
    FuneralDEP: 'mdi-flower-outline'
  }
  return icons[benefit] || 'mdi-shield'
}

const getBenefitIconColor = (benefit: string): string => {
  const colors: Record<string, string> = {
    GLA: 'primary',
    CI: 'error',
    PTD: 'warning',
    TTD: 'info',
    SGLA: 'deep-purple',
    PHI: 'teal'
  }
  return colors[benefit] || 'grey'
}

const getLossRatioColor = (ratio: number): string => {
  if (!ratio) return 'medium-emphasis'
  if (ratio <= 70) return 'success'
  if (ratio <= 85) return 'warning'
  return 'error'
}

const getMatchScoreColor = (score: number): string => {
  if (score >= 95) return 'success'
  if (score >= 80) return 'warning'
  return 'error'
}

const loadAnalytics = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await GroupPricingService.getBordereauxAnalytics({
      period: selectedPeriod.value
    })
    const data = response.data || {}
    topMetrics.value = data.top_metrics || topMetrics.value
    benefitPerformance.value = data.benefit_performance || []
    insurerMetrics.value = data.insurer_metrics || []
    monthlyKPIs.value = data.monthly_kpis || []
  } catch (error: any) {
    const msg =
      error.response?.data?.error ||
      error.message ||
      'Failed to load bordereaux analytics'
    errorMessage.value = msg
    flash.show(msg, 'error')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadAnalytics()
})
</script>

<style scoped>
.h-100 {
  height: 100%;
}
</style>
