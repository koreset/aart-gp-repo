<template>
  <v-container fluid>
    <base-card :show-actions="false">
      <template #header>
        <h3 class="mb-0">Premium Dashboard</h3>
      </template>
      <template #default>
        <!-- Year selector -->
        <v-row class="mb-2" align="center">
          <v-col cols="12" md="3">
            <v-select
              v-model="selectedYear"
              :items="availableYears"
              label="Year"
              variant="outlined"
              density="compact"
              @update:model-value="loadDashboard"
            />
          </v-col>
          <v-col>
            <v-btn
              variant="outlined"
              color="primary"
              prepend-icon="mdi-refresh"
              @click="loadDashboard"
            >
              Refresh
            </v-btn>
          </v-col>
        </v-row>

        <!-- KPI Cards -->
        <v-row class="mb-4">
          <v-col
            v-for="card in kpiCards"
            :key="card.title"
            cols="12"
            sm="6"
            md="3"
          >
            <stat-card
              :title="card.title"
              :value="card.value"
              :icon="card.icon"
              color="primary"
              :loading="loading"
            />
          </v-col>
        </v-row>

        <!-- Charts Row -->
        <v-row class="mb-4">
          <v-col cols="12" md="8">
            <v-card variant="outlined">
              <v-card-title class="text-subtitle-1 pa-3"
                >Premiums Due vs Collected (12 months)</v-card-title
              >
              <v-card-text>
                <ag-charts
                  v-if="trendOptions"
                  :options="trendOptions"
                  style="height: 260px"
                />
                <v-skeleton-loader v-else type="image" height="260" />
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="4">
            <v-card variant="outlined">
              <v-card-title class="text-subtitle-1 pa-3"
                >Collection Status</v-card-title
              >
              <v-card-text>
                <ag-charts
                  v-if="statusOptions"
                  :options="statusOptions"
                  style="height: 260px"
                />
                <v-skeleton-loader v-else type="image" height="260" />
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <!-- Top Outstanding Table -->
        <v-row>
          <v-col cols="12">
            <v-card variant="outlined">
              <v-card-title class="text-subtitle-1 pa-3"
                >Top Outstanding Schemes</v-card-title
              >
              <v-card-text class="pa-0">
                <v-data-table
                  :headers="outstandingHeaders"
                  :items="topOutstanding"
                  density="compact"
                  hide-default-footer
                  :items-per-page="10"
                >
                  <template #[`item.balance`]="{ item }">
                    <span class="text-error font-weight-medium">{{
                      formatCurrency(item.balance)
                    }}</span>
                  </template>
                  <template #[`item.status`]="{ item }">
                    <v-chip
                      :color="statusColor(item.status)"
                      size="x-small"
                      variant="tonal"
                    >
                      {{ item.status }}
                    </v-chip>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </template>
    </base-card>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { AgCharts } from 'ag-charts-vue3'
import type { AgChartOptions } from 'ag-charts-community'
import PremiumManagementService from '@/renderer/api/PremiumManagementService'
import BaseCard from '@/renderer/components/BaseCard.vue'
import StatCard from '@/renderer/components/StatCard.vue'
import { CHART_COLORS } from '@/renderer/constants/designTokens'
import { fmtCurrency } from '@/renderer/utils/formatters'
import { useStatusBarStore } from '@/renderer/store/statusBar'

const statusBarStore = useStatusBarStore()

const currentYear = new Date().getFullYear()
const selectedYear = ref(currentYear)
const availableYears = Array.from({ length: 5 }, (_, i) => currentYear - i)

const loading = ref(false)
const kpis = ref<any>(null)
const monthlyTrend = ref<any[]>([])
const statusBreakdown = ref<any[]>([])
const topOutstanding = ref<any[]>([])

const kpiCards = computed(() => [
  {
    title: 'Due This Month',
    value: kpis.value ? formatCurrency(kpis.value.due_this_month) : '—',
    icon: 'mdi-cash-multiple'
  },
  {
    title: 'Collected',
    value: kpis.value ? formatCurrency(kpis.value.collected) : '—',
    icon: 'mdi-check-circle-outline'
  },
  {
    title: 'Collection Rate',
    value: kpis.value ? kpis.value.collection_rate.toFixed(1) + '%' : '—',
    icon: 'mdi-percent-outline'
  },
  {
    title: 'Overdue',
    value: kpis.value ? formatCurrency(kpis.value.overdue) : '—',
    icon: 'mdi-alert-circle-outline'
  }
])

const trendOptions = computed(() => {
  if (!monthlyTrend.value.length) return null
  return {
    title: { text: '' },
    data: monthlyTrend.value,
    series: [
      {
        type: 'bar',
        xKey: 'month',
        yKey: 'due',
        yName: 'Due',
        fill: CHART_COLORS[0]
      },
      {
        type: 'bar',
        xKey: 'month',
        yKey: 'collected',
        yName: 'Collected',
        fill: CHART_COLORS[2]
      }
    ],
    axes: [
      { type: 'category', position: 'bottom' },
      {
        type: 'number',
        position: 'left',
        label: { formatter: (p: any) => formatShort(p.value) }
      }
    ],
    legend: { position: 'bottom' }
  } as unknown as AgChartOptions
})

const statusOptions = computed(() => {
  if (!statusBreakdown.value.length) return null
  return {
    title: { text: '' },
    data: statusBreakdown.value.filter((s) => s.count > 0),
    series: [
      {
        type: 'pie',
        angleKey: 'amount',
        calloutLabelKey: 'status',
        innerRadiusRatio: 0.6
      }
    ]
  } as unknown as AgChartOptions
})

const outstandingHeaders = [
  { title: 'Scheme', key: 'scheme_name' },
  {
    title: 'Due',
    key: 'amount_due',
    value: (item: any) => formatCurrency(item.amount_due)
  },
  {
    title: 'Paid',
    key: 'amount_paid',
    value: (item: any) => formatCurrency(item.amount_paid)
  },
  { title: 'Balance', key: 'balance' },
  { title: 'Status', key: 'status' }
]

function formatCurrency(val: number) {
  return fmtCurrency(val)
}

function formatShort(val: number) {
  if (val >= 1_000_000) return `R${(val / 1_000_000).toFixed(1)}M`
  if (val >= 1_000) return `R${(val / 1_000).toFixed(0)}K`
  return `R${val}`
}

function statusColor(status: string) {
  const map: Record<string, string> = {
    current: 'success',
    in_arrears: 'warning',
    suspended: 'error'
  }
  return map[status] ?? 'grey'
}

async function loadDashboard() {
  loading.value = true
  try {
    const res = await PremiumManagementService.getPremiumDashboard(
      selectedYear.value
    )
    const data = res.data.data
    kpis.value = data.kpis
    monthlyTrend.value = data.monthly_trend ?? []
    statusBreakdown.value = data.status_breakdown ?? []
    topOutstanding.value = data.top_outstanding ?? []
  } catch (e) {
    console.error('Failed to load premium dashboard', e)
  } finally {
    loading.value = false
    if (kpis.value) {
      const rate = kpis.value.collection_rate ?? 0
      const overdue = kpis.value.overdue ?? 0
      statusBarStore.set([
        { icon: 'mdi-calendar', text: `Year: ${selectedYear.value}` },
        {
          icon: 'mdi-percent',
          text: `Collection: ${rate.toFixed(1)}%`,
          severity: rate < 80 ? 'warn' : 'info'
        },
        {
          icon: 'mdi-alert-circle-outline',
          text: `Overdue: ${fmtShort(overdue)}`,
          severity: overdue > 0 ? 'error' : 'info'
        }
      ])
    }
  }
}

function fmtShort(val: number) {
  if (val >= 1_000_000) return `R${(val / 1_000_000).toFixed(1)}M`
  if (val >= 1_000) return `R${(val / 1_000).toFixed(0)}K`
  return `R${val}`
}

onMounted(loadDashboard)
onUnmounted(() => statusBarStore.clear())
</script>
