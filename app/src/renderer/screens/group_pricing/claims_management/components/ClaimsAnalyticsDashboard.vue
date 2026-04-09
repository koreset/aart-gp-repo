<template>
  <v-container>
    <v-row>
      <!-- Summary Cards -->
      <v-col cols="12" md="3">
        <v-card class="text-center pa-4" color="primary" variant="tonal">
          <div class="text-h6 font-weight-bold">{{
            analyticsData.totalClaims
          }}</div>
          <div class="text-caption">Total Claims (YTD)</div>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card class="text-center pa-4" color="success" variant="tonal">
          <div class="text-h6 font-weight-bold">{{
            formatCurrency(analyticsData.totalPaid)
          }}</div>
          <div class="text-caption">Total Reported (YTD)</div>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card class="text-center pa-4" color="warning" variant="tonal">
          <div class="text-h6 font-weight-bold">{{
            analyticsData.avgProcessingDays
          }}</div>
          <div class="text-caption">Avg Processing Days</div>
        </v-card>
      </v-col>
      <v-col cols="12" md="3">
        <v-card class="text-center pa-4" color="info" variant="tonal">
          <div class="text-h6 font-weight-bold"
            >{{ analyticsData.approvalRate }}%</div
          >
          <div class="text-caption">Approval Rate</div>
        </v-card>
      </v-col>
    </v-row>

    <!-- Filter Controls -->
    <v-row class="mb-4">
      <v-col cols="12" md="3">
        <v-select
          v-model="selectedPeriod"
          :items="periods"
          label="Time Period"
          variant="outlined"
          density="compact"
        />
      </v-col>

      <v-col cols="12" md="3">
        <v-select
          v-model="selectedScheme"
          :items="schemes"
          item-title="name"
          item-value="id"
          label="Filter by Scheme"
          variant="outlined"
          density="compact"
          clearable
        />
      </v-col>
      <v-col cols="12" md="3">
        <v-select
          v-model="selectedBenefitType"
          :items="benefitTypes"
          label="Filter by Benefit Type"
          variant="outlined"
          density="compact"
          clearable
        />
      </v-col>
      <v-col cols="12" md="3">
        <v-btn
          color="primary"
          variant="outlined"
          prepend-icon="mdi-filter"
          @click="loadAnalyticsData"
        >
          Apply Filters
        </v-btn>
      </v-col>
    </v-row>

    <!-- Custom Date Range Row (shows when Custom Range is selected) -->
    <v-row v-if="selectedPeriod === 'custom'" class="mb-4">
      <v-col cols="12" md="4">
        <v-text-field
          v-model="startDate"
          label="Start Date"
          type="date"
          variant="outlined"
          density="compact"
        />
      </v-col>
      <v-col cols="12" md="4">
        <v-text-field
          v-model="endDate"
          label="End Date"
          type="date"
          variant="outlined"
          density="compact"
        />
      </v-col>
      <v-col cols="12" md="4" class="d-flex align-center">
        <v-btn
          color="secondary"
          variant="outlined"
          prepend-icon="mdi-calendar-clear"
          @click="clearCustomDates"
        >
          Clear Dates
        </v-btn>
      </v-col>
    </v-row>

    <!-- Charts Row 1 -->
    <v-row>
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>Claims by Status</v-card-title>
          <v-card-text>
            <div ref="claimsStatusChart" style="height: 300px"></div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>Claims by Benefit Type</v-card-title>
          <v-card-text>
            <div ref="benefitTypeChart" style="height: 300px"></div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Charts Row 2 -->
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-text>
            <ag-charts
              v-if="claimsTrendChartOptions"
              :options="claimsTrendChartOptions"
              style="height: 400px; width: 100%"
            />
            <div
              v-else
              style="
                height: 200px;
                display: flex;
                align-items: center;
                justify-content: center;
                color: #666;
              "
            >
              No trend data available for the selected period
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Charts Row 3 -->
    <v-row>
      <v-col cols="12" md="6">
        <v-card>
          <v-card-text>
            <ag-charts
              v-if="processingTimeChartOptions"
              :options="processingTimeChartOptions"
              style="height: 300px; width: 100%"
            />
            <div
              v-else
              style="
                height: 240px;
                display: flex;
                align-items: center;
                justify-content: center;
                color: #666;
              "
            >
              No processing time data available for the selected period
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="6">
        <v-card>
          <v-card-text>
            <ag-charts
              v-if="ageDistributionChartOptions"
              :options="ageDistributionChartOptions"
              style="height: 300px; width: 100%"
            />
            <div
              v-else
              style="
                height: 240px;
                display: flex;
                align-items: center;
                justify-content: center;
                color: #666;
              "
            >
              No age distribution data available for the selected period
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Top Claims Table -->
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>Largest Claims This Period</v-card-title>
          <v-card-text>
            <GroupPricingDataGrid
              :rowData="topClaims"
              :columnDefs="claimsHeaders"
              :showExport="true"
              :tableTitle="''"
              :density="'compact'"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- KPI Cards -->
    <v-row>
      <v-col cols="12" md="4">
        <v-card>
          <v-card-title class="bg-error text-white"
            >Declined Claims Analysis</v-card-title
          >
          <v-card-text>
            <div class="text-h6">{{ analyticsData.declineRate }}%</div>
            <div class="text-caption">Decline Rate</div>
            <v-divider class="my-2" />
            <div class="text-body-2 mb-2">Top Decline Reasons:</div>
            <v-chip
              v-for="reason in analyticsData.topDeclineReasons"
              :key="reason.reason"
              size="small"
              class="ma-1"
            >
              {{ reason.reason }} ({{ reason.count }})
            </v-chip>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="4">
        <v-card>
          <v-card-title class="bg-warning text-white"
            >Processing Efficiency</v-card-title
          >
          <v-card-text>
            <div class="text-h6">{{ analyticsData.efficiencyScore || 0 }}%</div>
            <div class="text-caption">Overall Efficiency Score</div>
            <v-divider class="my-2" />

            <!-- SLA Performance -->
            <div class="text-body-2 mb-1"
              >SLA Performance ({{ analyticsData.slaDays || 30 }} days):</div
            >
            <v-progress-linear
              :model-value="analyticsData.slaCompliance"
              color="success"
              height="8"
              class="mb-2"
            />
            <div class="text-caption mb-3"
              >{{ analyticsData.slaCompliance }}% on-time</div
            >

            <!-- Additional Metrics -->
            <div class="d-flex justify-space-between text-body-2 mb-1">
              <span>Closed this period:</span>
              <span class="font-weight-bold">{{
                analyticsData.closedInPeriod || 0
              }}</span>
            </div>
            <div class="d-flex justify-space-between text-body-2 mb-1">
              <span>Open claims (WIP):</span>
              <span class="font-weight-bold">{{
                analyticsData.wipOpenClaims || 0
              }}</span>
            </div>
            <div class="d-flex justify-space-between text-body-2">
              <span>Avg processing time:</span>
              <span class="font-weight-bold"
                >{{ analyticsData.avgProcessingDays }} days</span
              >
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="4">
        <v-card>
          <v-card-title class="bg-info text-white"
            >Assessor Performance</v-card-title
          >
          <v-card-text>
            <div class="text-body-2 mb-2">Top Performers:</div>
            <v-list density="compact">
              <v-list-item
                v-for="assessor in analyticsData.topAssessors"
                :key="assessor.name"
              >
                <v-list-item-title>{{ assessor.name }}</v-list-item-title>
                <v-list-item-subtitle
                  >{{ assessor.claimsProcessed }} claims</v-list-item-subtitle
                >
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { AgCharts } from 'ag-charts-vue3'
import type { AgChartOptions } from 'ag-charts-community'
import GroupPricingDataGrid from '@/renderer/components/tables/GroupPricingDataGrid.vue'

interface Props {
  schemes: Array<any>
}

interface Emits {
  (e: 'close'): void
}

defineProps<Props>()
defineEmits<Emits>()

// Chart refs
const claimsStatusChart = ref<HTMLElement>()
const benefitTypeChart = ref<HTMLElement>()

// Claims trend chart options for AG Charts
const claimsTrendChartOptions = ref<AgChartOptions | null>(null)

// Claims age distribution chart options for AG Charts
const ageDistributionChartOptions = ref<AgChartOptions | null>(null)

// Processing time distribution chart options for AG Charts
const processingTimeChartOptions = ref<AgChartOptions | null>(null)

// Filter states
const selectedScheme = ref<number | null>(null)
const selectedPeriod = ref('last_12_months')
const selectedBenefitType = ref<string | null>(null)
const startDate = ref<string>('')
const endDate = ref<string>('')
const loading = ref(false)

// Data
const analyticsData: any = ref({
  totalClaims: 0,
  totalPaid: 0,
  avgProcessingDays: 0,
  approvalRate: 0,
  declineRate: 8,
  slaCompliance: 94,
  efficiencyScore: 0,
  closedInPeriod: 0,
  slaDays: 30,
  wipOpenClaims: 0,
  throughputPerWeek: 0,
  byStatus: {} as Record<string, number>,
  byBenefit: {} as Record<string, number>,
  claimsTrend: [] as Array<{ month: string; amount: number; count: number }>,
  claimantAgeDistribution: [] as Array<{ label: string; count: number }>,
  processingTimeDistribution: [] as Array<{ label: string; count: number }>,
  topDeclineReasons: [],
  topAssessors: []
})

// Options
const periods = [
  { title: 'Last 30 Days', value: 'last_30_days' },
  { title: 'Last 3 Months', value: 'last_3_months' },
  { title: 'Last 6 Months', value: 'last_6_months' },
  { title: 'Last 12 Months', value: 'last_12_months' },
  { title: 'Year to Date', value: 'ytd' },
  { title: 'Custom Range', value: 'custom' }
]

const benefitTypes = [
  'Group Life Assurance (GLA)',
  'Spouse Group Life Assurance (SGLA)',
  'Permanent Total Disability (PTD)',
  'Critical Illness (CI)',
  'Temporary Total Disability (TTD)',
  'Personal Health Insurance (PHI)',
  'Group Family Funeral (GFF)'
]

// Table headers for AG Grid
const claimsHeaders = [
  {
    field: 'claim_number',
    headerName: 'Claim Number',
    sortable: true,
    filter: true,
    width: 150
  },
  {
    field: 'member_name',
    headerName: 'Member',
    sortable: true,
    filter: true,
    width: 200
  },
  {
    field: 'benefit_type',
    headerName: 'Benefit Type',
    sortable: true,
    filter: true,
    width: 150
  },
  {
    field: 'claim_amount',
    headerName: 'Amount',
    sortable: true,
    filter: true,
    width: 150,
    type: 'numericColumn',
    cellRenderer: (params: any) => {
      if (params.value) {
        return new Intl.NumberFormat('en-ZA', {
          style: 'currency',
          currency: 'ZAR',
          minimumFractionDigits: 0
        }).format(params.value)
      }
      return ''
    }
  },
  {
    field: 'status',
    headerName: 'Status',
    sortable: true,
    filter: true,
    width: 130,
    cellRenderer: (params: any) => {
      if (params.value) {
        const status = params.value
        const formattedStatus =
          status
            ?.split('_')
            .map((word: string) => word.charAt(0).toUpperCase() + word.slice(1))
            .join(' ') || ''

        const colorMap: Record<string, string> = {
          pending: '#2196f3',
          under_assessment: '#ff9800',
          approved: '#4caf50',
          declined: '#f44336',
          paid: '#009688'
        }

        const color = colorMap[status] || '#666'
        return `<span style="background: ${color}; color: white; padding: 2px 8px; border-radius: 12px; font-size: 11px; font-weight: 500;">${formattedStatus}</span>`
      }
      return ''
    }
  },
  {
    field: 'date_notified',
    headerName: 'Date',
    sortable: true,
    filter: true,
    width: 120
  }
]

// Mock data
const topClaims = ref([
  {
    claim_number: 'CLM-2024-001',
    member_name: 'John Doe',
    benefit_type: 'GLA',
    claim_amount: 2500000,
    status: 'approved',
    date_notified: '2024-01-15'
  },
  {
    claim_number: 'CLM-2024-002',
    member_name: 'Jane Smith',
    benefit_type: 'PTD',
    claim_amount: 2000000,
    status: 'under_assessment',
    date_notified: '2024-01-20'
  },
  {
    claim_number: 'CLM-2024-003',
    member_name: 'Mike Johnson',
    benefit_type: 'CI',
    claim_amount: 1800000,
    status: 'approved',
    date_notified: '2024-01-25'
  }
])

// Methods
const loadAnalyticsData = async () => {
  loading.value = true
  try {
    const filters: any = {
      scheme_id: selectedScheme.value,
      period: selectedPeriod.value,
      benefit_type: selectedBenefitType.value
    }

    // Add custom date range if period is 'custom'
    if (selectedPeriod.value === 'custom') {
      if (startDate.value) filters.from = startDate.value
      if (endDate.value) filters.to = endDate.value
    }

    const response = await GroupPricingService.getClaimsAnalytics(filters)
    const data = response.data

    // Map API response to local data structure
    analyticsData.value.totalClaims = data.total_claims || 0
    analyticsData.value.totalPaid = data.total_paid_amount || 0
    analyticsData.value.avgProcessingDays = data.avg_processing_days || 0
    analyticsData.value.approvalRate = Math.round(
      (data.approval_rate || 0) * 100
    )
    analyticsData.value.declineRate = Math.round((data.decline_rate || 0) * 100)
    analyticsData.value.byStatus = data.by_status || {}
    analyticsData.value.byBenefit = data.by_benefit || {}
    analyticsData.value.claimsTrend = data.claims_trend || []
    analyticsData.value.claimantAgeDistribution =
      data.claimant_age_distribution || []

    // Create AG Charts options for claims trend
    if (
      analyticsData.value.claimsTrend &&
      analyticsData.value.claimsTrend.length > 0
    ) {
      const chartData = analyticsData.value.claimsTrend.map((item: any) => {
        const date = new Date(item.month + '-01')
        const monthName = date.toLocaleDateString('en-US', {
          year: 'numeric',
          month: 'short'
        })
        return {
          month: monthName,
          count: item.count,
          amount: item.amount / 1000000 // Convert to millions for readability
        }
      })

      claimsTrendChartOptions.value = {
        title: {
          text: 'Claims Trend Over Time',
          fontSize: 16,
          fontWeight: 'bold'
        },
        data: chartData,
        series: [
          {
            type: 'bar',
            xKey: 'month',
            yKey: 'count',
            yName: 'Number of Claims',
            fill: '#2196F3',
            stroke: '#1976D2',
            strokeWidth: 1
          },
          {
            type: 'line',
            xKey: 'month',
            yKey: 'amount',
            yName: 'Amount (Millions ZAR)',
            stroke: '#FF5722',
            strokeWidth: 3,
            marker: {
              fill: '#FF5722',
              strokeWidth: 2,
              size: 8
            }
          }
        ],
        axes: [
          {
            type: 'category',
            position: 'bottom',
            title: {
              text: 'Month'
            }
          },
          {
            type: 'number',
            position: 'left',
            title: {
              text: 'Number of Claims'
            }
          },
          {
            type: 'number',
            position: 'right',
            title: {
              text: 'Amount (Millions ZAR)'
            },
            label: {
              formatter: (params: any) => `R${params.value.toFixed(1)}M`
            }
          }
        ],
        legend: {
          position: 'top',
          enabled: true
        }
      }
    } else {
      claimsTrendChartOptions.value = null
    }

    // Create AG Charts options for claimant age distribution
    if (
      analyticsData.value.claimantAgeDistribution &&
      analyticsData.value.claimantAgeDistribution.length > 0
    ) {
      const ageData = analyticsData.value.claimantAgeDistribution

      ageDistributionChartOptions.value = {
        title: {
          text: 'Distribution of Claimant Ages',
          fontSize: 16,
          fontWeight: 'bold'
        },
        data: ageData,
        series: [
          {
            type: 'bar',
            xKey: 'label',
            yKey: 'count',
            yName: 'Frequency',
            fill: '#8E24AA',
            stroke: '#7B1FA2',
            strokeWidth: 1
          }
        ],
        axes: [
          {
            type: 'category',
            position: 'bottom',
            title: {
              text: 'Age of Claimant'
            }
          },
          {
            type: 'number',
            position: 'left',
            title: {
              text: 'Frequency'
            }
          }
        ],
        legend: {
          enabled: false
        }
      }
    } else {
      ageDistributionChartOptions.value = null
    }

    // Map processing efficiency data from API response
    if (data.processing_efficiency) {
      const efficiency = data.processing_efficiency
      analyticsData.value.slaCompliance = Math.round(
        (efficiency.on_time_rate || 0) * 100
      )
      analyticsData.value.efficiencyScore = Math.round(
        (efficiency.efficiency_score || 0) * 100
      )
      analyticsData.value.closedInPeriod = efficiency.closed_in_period || 0
      analyticsData.value.slaDays = efficiency.sla_days || 30
      analyticsData.value.wipOpenClaims = efficiency.wip_open_claims || 0
      analyticsData.value.throughputPerWeek =
        efficiency.throughput_per_week || 0

      // Set processing time distribution data
      analyticsData.value.processingTimeDistribution =
        efficiency.processing_time_distribution || []

      // Create AG Charts options for processing time distribution
      if (
        analyticsData.value.processingTimeDistribution &&
        analyticsData.value.processingTimeDistribution.length > 0
      ) {
        processingTimeChartOptions.value = {
          title: {
            text: 'Processing Time Distribution',
            fontSize: 16,
            fontWeight: 'bold'
          },
          data: analyticsData.value.processingTimeDistribution,
          series: [
            {
              type: 'bar',
              xKey: 'label',
              yKey: 'count',
              yName: 'Frequency',
              fill: '#FF9800',
              stroke: '#F57C00',
              strokeWidth: 1
            }
          ],
          axes: [
            {
              type: 'category',
              position: 'bottom',
              title: {
                text: 'Processing Time (Days)'
              }
            },
            {
              type: 'number',
              position: 'left',
              title: {
                text: 'Frequency'
              }
            }
          ],
          legend: {
            enabled: false
          }
        }
      } else {
        processingTimeChartOptions.value = null
      }
    }
    if (data.top_claims && Array.isArray(data.top_claims)) {
      topClaims.value = data.top_claims.map((claim) => ({
        claim_number: claim.claim_number,
        member_name: claim.member,
        benefit_type: claim.benefit_type,
        claim_amount: claim.amount,
        status: claim.status,
        date_notified: claim.date_notified
      }))
    }

    // Map top_assessors from API to expected structure
    if (data.top_assessors && Array.isArray(data.top_assessors)) {
      analyticsData.value.topAssessors = data.top_assessors.map((assessor) => ({
        name: assessor.name,
        claimsProcessed: assessor.count
      }))
    }

    // Map topDeclineReasons from API to expected structure
    if (data.topDeclineReasons && Array.isArray(data.topDeclineReasons)) {
      analyticsData.value.topDeclineReasons = data.topDeclineReasons.map(
        (reason) => ({
          reason: reason.reason,
          count: reason.count
        })
      )
    }

    // Update charts with new data
    await nextTick()
    updateCharts()
  } catch (error) {
    console.error('Error loading analytics data:', error)
    // Keep existing mock data on error
  } finally {
    loading.value = false
  }
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR',
    minimumFractionDigits: 0
  }).format(amount)
}

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    pending: 'info',
    under_assessment: 'warning',
    approved: 'success',
    declined: 'error',
    paid: 'teal'
  }
  return colors[status] || 'default'
}

const formatStatus = (status: string) => {
  return (
    status
      ?.split('_')
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ') || ''
  )
}

const createClaimsStatusChart = () => {
  if (!claimsStatusChart.value) return

  // Use real API data
  const data = analyticsData.value.byStatus
  const maxValue = Math.max(...Object.values(data).map((v) => Number(v) || 0))

  claimsStatusChart.value.innerHTML = `
    <div style="display: flex; flex-direction: column; gap: 10px; padding: 20px;">
      ${Object.entries(data)
        .map(([status, count]) => {
          const countNum = Number(count) || 0
          return `
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span>${formatStatus(status)}</span>
          <div style="display: flex; align-items: center; gap: 10px;">
            <div style="width: 100px; height: 8px; background: #e0e0e0; border-radius: 4px;">
              <div style="width: ${maxValue > 0 ? (countNum / maxValue) * 100 : 0}%; height: 100%; background: ${getStatusColor(status) === 'success' ? '#4caf50' : getStatusColor(status) === 'warning' ? '#ff9800' : getStatusColor(status) === 'error' ? '#f44336' : '#2196f3'}; border-radius: 4px;"></div>
            </div>
            <span style="font-weight: bold;">${countNum}</span>
          </div>
        </div>
      `
        })
        .join('')}
    </div>
  `
}

const createBenefitTypeChart = () => {
  if (!benefitTypeChart.value) return

  // Use real API data
  const data = analyticsData.value.byBenefit
  const maxValue = Math.max(...Object.values(data).map((v) => Number(v) || 0))

  benefitTypeChart.value.innerHTML = `
    <div style="display: flex; flex-direction: column; gap: 8px; padding: 20px;">
      ${Object.entries(data)
        .map(([type, count]) => {
          const countNum = Number(count) || 0
          return `
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span title="${type}">${type.length > 30 ? type.substring(0, 30) + '...' : type}</span>
          <div style="display: flex; align-items: center; gap: 10px;">
            <div style="width: 80px; height: 6px; background: #e0e0e0; border-radius: 3px;">
              <div style="width: ${maxValue > 0 ? (countNum / maxValue) * 100 : 0}%; height: 100%; background: #2196f3; border-radius: 3px;"></div>
            </div>
            <span style="font-weight: bold; min-width: 30px;">${countNum}</span>
          </div>
        </div>
      `
        })
        .join('')}
    </div>
  `
}

const updateCharts = () => {
  nextTick(() => {
    createClaimsStatusChart()
    createBenefitTypeChart()
    // Processing time and age distribution charts are now handled by AG Charts components
  })
}

// const handleFilterChange = () => {
//   loadAnalyticsData()
// }

const clearCustomDates = () => {
  startDate.value = ''
  endDate.value = ''
  loadAnalyticsData()
}

// const exportReport = () => {
//   // This would generate and download an analytics report
//   console.log('Exporting analytics report...')
// }

// Lifecycle
onMounted(() => {
  loadAnalyticsData()
})
</script>

<style scoped>
.v-card-title {
  font-size: 1.1rem;
  font-weight: 500;
}
</style>
