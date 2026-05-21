<template>
  <div ref="dashboardRoot" class="claims-analytics">
    <!-- Sticky control bar -->
    <div class="claims-analytics__controls no-print">
      <div class="claims-analytics__filters">
        <v-select
          v-model="selectedPeriod"
          :items="periods"
          label="Period"
          variant="outlined"
          density="compact"
          hide-details
          class="filter-input"
        />
        <v-select
          v-model="selectedScheme"
          :items="schemes"
          item-title="name"
          item-value="id"
          label="Scheme"
          variant="outlined"
          density="compact"
          hide-details
          clearable
          class="filter-input filter-input--wide"
        />
        <v-select
          v-model="selectedBenefitType"
          :items="benefitTypeOptions"
          label="Benefit type"
          variant="outlined"
          density="compact"
          hide-details
          clearable
          class="filter-input filter-input--wide"
        />
        <template v-if="selectedPeriod === 'custom'">
          <v-text-field
            v-model="startDate"
            label="From"
            type="date"
            variant="outlined"
            density="compact"
            hide-details
            class="filter-input"
          />
          <v-text-field
            v-model="endDate"
            label="To"
            type="date"
            variant="outlined"
            density="compact"
            hide-details
            class="filter-input"
          />
        </template>
      </div>
      <div class="claims-analytics__actions">
        <span v-if="lastUpdatedAt" class="claims-analytics__updated">
          Updated {{ lastUpdatedLabel }}
        </span>
        <v-tooltip text="Refresh" location="bottom">
          <template #activator="{ props: tipProps }">
            <v-btn
              v-bind="tipProps"
              icon
              size="small"
              variant="text"
              :loading="loading"
              @click="loadAnalyticsData()"
            >
              <v-icon>mdi-refresh</v-icon>
            </v-btn>
          </template>
        </v-tooltip>
        <v-menu location="bottom end">
          <template #activator="{ props: menuProps }">
            <v-btn
              v-bind="menuProps"
              size="small"
              variant="outlined"
              prepend-icon="mdi-download"
              :disabled="loading || exporting"
              :loading="exporting"
            >
              Export
            </v-btn>
          </template>
          <v-list density="compact">
            <v-list-item
              prepend-icon="mdi-file-pdf-box"
              title="Snapshot (PDF)"
              @click="exportPdf()"
            />
            <v-list-item
              prepend-icon="mdi-file-excel"
              title="Workbook (Excel)"
              @click="exportExcel()"
            />
          </v-list>
        </v-menu>
      </div>
    </div>

    <v-alert
      v-if="errorMessage"
      type="error"
      variant="tonal"
      density="compact"
      class="mb-3"
      closable
      @click:close="errorMessage = null"
    >
      {{ errorMessage }}
    </v-alert>

    <!-- Section 1: Financial KPIs -->
    <div class="kpi-strip">
      <stat-card
        title="Total Claims"
        :value="formatNumber(analyticsData.totalClaims)"
        icon="mdi-file-document-multiple-outline"
        color="primary"
        :subtitle="periodLabel"
        :loading="loading"
        class="kpi-clickable"
        @click="drillDown({})"
      />
      <stat-card
        title="Total Paid"
        :value="formatCurrency(analyticsData.totalPaid)"
        icon="mdi-cash-multiple"
        color="accent"
        subtitle="Approved + paid claims"
        :loading="loading"
        class="kpi-clickable"
        @click="drillDown({ status: 'paid' })"
      />
      <stat-card
        title="Avg Claim Amount"
        :value="formatCurrency(analyticsData.avgClaimAmount)"
        icon="mdi-trending-up"
        color="accent"
        subtitle="Mean of approved / paid"
        :loading="loading"
      />
      <stat-card
        title="Outstanding Exposure"
        :value="formatCurrency(analyticsData.totalExposure)"
        icon="mdi-alert-circle-outline"
        color="warning"
        subtitle="Open claim liability"
        :loading="loading"
        class="kpi-clickable"
        @click="drillDown({ status: 'pending_assessment' })"
      />
    </div>

    <!-- Section 2: Operational KPIs -->
    <div class="kpi-strip">
      <stat-card
        title="Approval Rate"
        :value="formatPercent(analyticsData.approvalRate)"
        icon="mdi-check-decagram-outline"
        color="success"
        subtitle="Approved or paid"
        :loading="loading"
        class="kpi-clickable"
        @click="drillDown({ status: 'approved' })"
      />
      <stat-card
        title="Decline Rate"
        :value="formatPercent(analyticsData.declineRate)"
        icon="mdi-close-octagon-outline"
        color="error"
        subtitle="Declined claims"
        :loading="loading"
        class="kpi-clickable"
        @click="drillDown({ status: 'declined' })"
      />
      <stat-card
        title="Avg Processing"
        :value="`${formatNumber(analyticsData.avgProcessingDays, 1)} days`"
        icon="mdi-clock-outline"
        color="info"
        :subtitle="`Median ${formatNumber(analyticsData.processingTimeMedian, 1)} days`"
        :loading="loading"
      />
      <stat-card
        title="SLA Compliance"
        :value="formatPercent(analyticsData.slaCompliance)"
        icon="mdi-timer-check-outline"
        :color="slaTone"
        :subtitle="`Within ${analyticsData.slaDays} days`"
        :loading="loading"
      />
    </div>

    <!-- Section 3: Operational health -->
    <v-row class="analytics-row">
      <v-col cols="12" md="6">
        <claims-analytics-chart-card
          title="Claims by status"
          subtitle="Click a slice to view those claims"
          :loading="loading"
          :empty="!statusDonutOptions"
          empty-text="No claims in the selected period"
          body-height="320px"
          :actions="chartActions.status"
        >
          <ag-charts
            v-if="statusDonutOptions"
            ref="statusChartRef"
            :options="statusDonutOptions"
            class="ag-fill"
          />
        </claims-analytics-chart-card>
      </v-col>
      <v-col cols="12" md="6">
        <claims-analytics-chart-card
          title="Aging of open claims"
          subtitle="Days since claim creation"
          :loading="loading"
          :empty="!agingChartOptions"
          empty-text="No open claims to age"
          body-height="320px"
          :actions="chartActions.aging"
        >
          <ag-charts
            v-if="agingChartOptions"
            ref="agingChartRef"
            :options="agingChartOptions"
            class="ag-fill"
          />
        </claims-analytics-chart-card>
      </v-col>
    </v-row>

    <!-- Section 4: Trend -->
    <v-row class="analytics-row">
      <v-col cols="12">
        <claims-analytics-chart-card
          title="Claims trend over time"
          subtitle="Volume and value by month"
          :loading="loading"
          :empty="!trendChartOptions"
          empty-text="No trend data for the selected period"
          body-height="380px"
          :actions="chartActions.trend"
        >
          <ag-charts
            v-if="trendChartOptions"
            ref="trendChartRef"
            :options="trendChartOptions"
            class="ag-fill"
          />
        </claims-analytics-chart-card>
      </v-col>
    </v-row>

    <!-- Section 5: Mix analysis -->
    <v-row class="analytics-row">
      <v-col cols="12" md="4">
        <claims-analytics-chart-card
          title="By benefit type"
          :loading="loading"
          :empty="!benefitChartOptions"
          empty-text="No benefit data"
          body-height="300px"
          :actions="chartActions.benefit"
        >
          <ag-charts
            v-if="benefitChartOptions"
            ref="benefitChartRef"
            :options="benefitChartOptions"
            class="ag-fill"
          />
        </claims-analytics-chart-card>
      </v-col>
      <v-col cols="12" md="4">
        <claims-analytics-chart-card
          title="By cause type"
          :loading="loading"
          :empty="!causeChartOptions"
          empty-text="No cause data"
          body-height="300px"
          :actions="chartActions.cause"
        >
          <ag-charts
            v-if="causeChartOptions"
            ref="causeChartRef"
            :options="causeChartOptions"
            class="ag-fill"
          />
        </claims-analytics-chart-card>
      </v-col>
      <v-col cols="12" md="4">
        <claims-analytics-chart-card
          title="By member type"
          subtitle="Principal / spouse / child"
          :loading="loading"
          :empty="!memberChartOptions"
          empty-text="No member data"
          body-height="300px"
          :actions="chartActions.member"
        >
          <ag-charts
            v-if="memberChartOptions"
            ref="memberChartRef"
            :options="memberChartOptions"
            class="ag-fill"
          />
        </claims-analytics-chart-card>
      </v-col>
    </v-row>

    <!-- Section 6: Processing & SLA -->
    <v-row class="analytics-row">
      <v-col cols="12" md="6">
        <claims-analytics-chart-card
          title="Processing time distribution"
          subtitle="Days to terminal status (approved / declined / paid)"
          :loading="loading"
          :empty="!processingChartOptions"
          empty-text="No processing data"
          body-height="320px"
          :actions="chartActions.processing"
        >
          <ag-charts
            v-if="processingChartOptions"
            ref="processingChartRef"
            :options="processingChartOptions"
            class="ag-fill"
          />
        </claims-analytics-chart-card>
      </v-col>
      <v-col cols="12" md="6">
        <claims-analytics-chart-card
          title="Claimant age distribution"
          subtitle="Principal member age at claim"
          :loading="loading"
          :empty="!ageChartOptions"
          empty-text="No age data"
          body-height="320px"
          :actions="chartActions.age"
        >
          <ag-charts
            v-if="ageChartOptions"
            ref="ageChartRef"
            :options="ageChartOptions"
            class="ag-fill"
          />
        </claims-analytics-chart-card>
      </v-col>
    </v-row>

    <!-- Section 7: Risk & quality -->
    <v-row class="analytics-row">
      <v-col cols="12" md="4">
        <stat-card
          title="Large claims watch"
          :value="formatNumber(analyticsData.largeClaimsCount)"
          icon="mdi-shield-alert-outline"
          color="warning"
          :subtitle="`≥ ${formatCurrency(analyticsData.largeClaimThreshold)}`"
          :loading="loading"
        />
      </v-col>
      <v-col cols="12" md="4">
        <stat-card
          title="Finance rejection rate"
          :value="formatPercent(analyticsData.financeRejectionRate)"
          icon="mdi-bank-remove"
          color="error"
          subtitle="Approved claims sent back"
          :loading="loading"
          class="kpi-clickable"
          @click="drillDown({ status: 'finance_rejected' })"
        />
      </v-col>
      <v-col cols="12" md="4">
        <stat-card
          title="Reopened claims"
          :value="formatNumber(analyticsData.reopenCount)"
          icon="mdi-history"
          color="info"
          subtitle="Terminal → reopened in period"
          :loading="loading"
        />
      </v-col>
    </v-row>

    <!-- Section 8: Tables -->
    <v-row class="analytics-row">
      <v-col cols="12" md="8">
        <claims-analytics-chart-card
          title="Largest claims in period"
          subtitle="Sorted by claim amount"
          :loading="loading"
          :empty="!loading && topClaims.length === 0"
          empty-text="No claims in the selected period"
          body-height="auto"
        >
          <group-pricing-data-grid
            :rowData="topClaims"
            :columnDefs="topClaimsColumns"
            :showExport="true"
            :tableTitle="''"
            tableName="largest_claims"
            density="compact"
          />
        </claims-analytics-chart-card>
      </v-col>
      <v-col cols="12" md="4">
        <claims-analytics-chart-card
          title="Top schemes"
          subtitle="By claim volume"
          :loading="loading"
          :empty="!loading && topSchemes.length === 0"
          empty-text="No scheme data"
          body-height="auto"
        >
          <div class="top-list">
            <div
              v-for="(item, idx) in topSchemes"
              :key="item.scheme_name + idx"
              class="top-list__row"
              role="button"
              tabindex="0"
              @click="drillDown({ scheme_name: item.scheme_name })"
              @keydown.enter="drillDown({ scheme_name: item.scheme_name })"
            >
              <div class="top-list__rank">{{ idx + 1 }}</div>
              <div class="top-list__main">
                <div class="top-list__title">{{ item.scheme_name }}</div>
                <div class="top-list__subtitle">
                  {{ formatNumber(item.count) }} claims ·
                  {{ formatCurrency(item.total_amount) }}
                </div>
                <div class="top-list__bar">
                  <div
                    class="top-list__bar-fill"
                    :style="{
                      width:
                        (topSchemes[0]?.count
                          ? (item.count / topSchemes[0].count) * 100
                          : 0) + '%'
                    }"
                  />
                </div>
              </div>
            </div>
          </div>
        </claims-analytics-chart-card>
      </v-col>
    </v-row>

    <!-- Section 9: Assessor & decline insights -->
    <v-row class="analytics-row">
      <v-col cols="12" md="6">
        <claims-analytics-chart-card
          title="Top assessor performance"
          subtitle="Assessments completed in period"
          :loading="loading"
          :empty="!loading && analyticsData.topAssessors.length === 0"
          empty-text="No assessments recorded"
          body-height="auto"
        >
          <div class="top-list">
            <div
              v-for="(a, idx) in analyticsData.topAssessors"
              :key="a.name + idx"
              class="top-list__row top-list__row--static"
            >
              <div class="top-list__rank top-list__rank--accent">
                {{ idx + 1 }}
              </div>
              <div class="top-list__main">
                <div class="top-list__title">{{ a.name }}</div>
                <div class="top-list__subtitle"
                  >{{ formatNumber(a.claimsProcessed) }} assessments</div
                >
                <div class="top-list__bar">
                  <div
                    class="top-list__bar-fill top-list__bar-fill--accent"
                    :style="{
                      width:
                        (analyticsData.topAssessors[0]?.claimsProcessed
                          ? (a.claimsProcessed /
                              analyticsData.topAssessors[0].claimsProcessed) *
                            100
                          : 0) + '%'
                    }"
                  />
                </div>
              </div>
            </div>
          </div>
        </claims-analytics-chart-card>
      </v-col>
      <v-col cols="12" md="6">
        <claims-analytics-chart-card
          title="Top decline reasons"
          subtitle="Sized by frequency"
          :loading="loading"
          :empty="!loading && analyticsData.topDeclineReasons.length === 0"
          empty-text="No declined claims in the selected period"
          body-height="auto"
        >
          <div class="reason-cloud">
            <div
              v-for="r in analyticsData.topDeclineReasons"
              :key="r.reason"
              class="reason-chip"
              :style="{
                fontSize: reasonChipSize(r.count) + 'rem',
                opacity: reasonChipOpacity(r.count)
              }"
            >
              {{ r.reason }}
              <span class="reason-chip__count">{{ r.count }}</span>
            </div>
          </div>
        </claims-analytics-chart-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRouter } from 'vue-router'
import { AgCharts } from 'ag-charts-vue3'
import type { AgChartOptions } from 'ag-charts-community'
import JsPdf from 'jspdf'
import html2canvas from 'html2canvas'
import * as XLSX from 'xlsx'

import GroupPricingService from '@/renderer/api/GroupPricingService'
import StatCard from '@/renderer/components/StatCard.vue'
import GroupPricingDataGrid from '@/renderer/components/tables/GroupPricingDataGrid.vue'
import ClaimsAnalyticsChartCard from './ClaimsAnalyticsChartCard.vue'

interface Props {
  schemes: Array<{ id: number; name: string }>
}
defineProps<Props>()

const router = useRouter()

// --- Filters / state ---
const selectedScheme = ref<number | null>(null)
const selectedPeriod = ref<string>('last_12_months')
const selectedBenefitType = ref<string | null>(null)
const startDate = ref<string>('')
const endDate = ref<string>('')

const loading = ref(false)
const exporting = ref(false)
const errorMessage = ref<string | null>(null)
const lastUpdatedAt = ref<Date | null>(null)
const dashboardRoot = ref<HTMLElement>()

const periods = [
  { title: 'Last 30 days', value: 'last_30_days' },
  { title: 'Last 3 months', value: 'last_3_months' },
  { title: 'Last 12 months', value: 'last_12_months' },
  { title: 'Year to date', value: 'ytd' },
  { title: 'Custom', value: 'custom' }
]

const benefitTypeOptions = [
  { title: 'Group Life Assurance (GLA)', value: 'GLA' },
  { title: 'Spouse Group Life Assurance (SGLA)', value: 'SGLA' },
  { title: 'Permanent Total Disability (PTD)', value: 'PTD' },
  { title: 'Critical Illness (CI)', value: 'CI' },
  { title: 'Temporary Total Disability (TTD)', value: 'TTD' },
  { title: 'Personal Health Insurance (PHI)', value: 'PHI' },
  { title: 'Group Family Funeral (GFF)', value: 'GFF' }
]

interface AnalyticsState {
  totalClaims: number
  totalPaid: number
  avgClaimAmount: number
  totalExposure: number
  avgProcessingDays: number
  processingTimeMedian: number
  approvalRate: number
  declineRate: number
  financeRejectionRate: number
  largeClaimsCount: number
  largeClaimThreshold: number
  reopenCount: number
  slaCompliance: number
  slaDays: number
  efficiencyScore: number
  closedInPeriod: number
  wipOpenClaims: number
  throughputPerWeek: number
  byStatus: Array<{ status: string; count: number }>
  byBenefit: Array<{ label: string; count: number }>
  byCauseType: Array<{ label: string; count: number }>
  byMemberType: Array<{ label: string; count: number }>
  agingBuckets: Array<{ label: string; count: number }>
  claimsTrend: Array<{ month: string; count: number; amount: number }>
  claimantAgeDistribution: Array<{ label: string; count: number }>
  processingTimeDistribution: Array<{ label: string; count: number }>
  topSchemes: Array<{
    scheme_name: string
    count: number
    total_amount: number
  }>
  topAssessors: Array<{ name: string; claimsProcessed: number }>
  topDeclineReasons: Array<{ reason: string; count: number }>
}

const blankAnalytics = (): AnalyticsState => ({
  totalClaims: 0,
  totalPaid: 0,
  avgClaimAmount: 0,
  totalExposure: 0,
  avgProcessingDays: 0,
  processingTimeMedian: 0,
  approvalRate: 0,
  declineRate: 0,
  financeRejectionRate: 0,
  largeClaimsCount: 0,
  largeClaimThreshold: 1000000,
  reopenCount: 0,
  slaCompliance: 0,
  slaDays: 30,
  efficiencyScore: 0,
  closedInPeriod: 0,
  wipOpenClaims: 0,
  throughputPerWeek: 0,
  byStatus: [],
  byBenefit: [],
  byCauseType: [],
  byMemberType: [],
  agingBuckets: [],
  claimsTrend: [],
  claimantAgeDistribution: [],
  processingTimeDistribution: [],
  topSchemes: [],
  topAssessors: [],
  topDeclineReasons: []
})

const analyticsData = ref<AnalyticsState>(blankAnalytics())
const topClaims = ref<any[]>([])
const topSchemes = computed(() => analyticsData.value.topSchemes)

// --- Theme palette helper ---
// Vuetify stores colours as space-separated RGB triples on root vars.
const themeColor = (token: string, fallback: string): string => {
  if (typeof window === 'undefined') return fallback
  const raw = getComputedStyle(document.documentElement)
    .getPropertyValue(`--v-theme-${token}`)
    .trim()
  if (!raw) return fallback
  const parts = raw.split(/\s+/).map((p) => Number(p))
  if (parts.length === 3 && parts.every((n) => !Number.isNaN(n))) {
    return `rgb(${parts[0]}, ${parts[1]}, ${parts[2]})`
  }
  return raw
}

const palette = computed(() => ({
  primary: themeColor('primary', '#003F58'),
  accent: themeColor('accent', '#006C8C'),
  info: themeColor('info', '#4338CA'),
  success: themeColor('success', '#059669'),
  warning: themeColor('warning', '#D97706'),
  error: themeColor('error', '#DC2626'),
  muted: 'rgba(100,116,139,0.8)'
}))

const statusPalette = computed<Record<string, string>>(() => ({
  draft: palette.value.muted,
  pending: palette.value.info,
  pending_assessment: palette.value.info,
  under_assessment: palette.value.warning,
  additional_info_required: palette.value.warning,
  approved: palette.value.success,
  paid: palette.value.accent,
  submitted_for_payment: palette.value.accent,
  finance_rejected: palette.value.error,
  declined: palette.value.error,
  cancelled: palette.value.muted
}))

const formatStatusLabel = (s: string) =>
  s
    .split('_')
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(' ')

// --- Formatters ---
const formatCurrency = (n: number) =>
  new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR',
    maximumFractionDigits: 0
  }).format(n || 0)

const formatNumber = (n: number, digits = 0) =>
  new Intl.NumberFormat('en-ZA', {
    minimumFractionDigits: digits,
    maximumFractionDigits: digits
  }).format(n || 0)

const formatPercent = (n: number) => `${(n || 0).toFixed(1)}%`

const periodLabel = computed(() => {
  const found = periods.find((p) => p.value === selectedPeriod.value)
  return found ? found.title : ''
})

const slaTone = computed(() => {
  const v = analyticsData.value.slaCompliance
  if (v >= 90) return 'success'
  if (v >= 75) return 'warning'
  return 'error'
})

// --- Last-updated label ---
// Tick refreshes the relative label every 30s without thrashing layout.
const tick = ref(Date.now())
let tickTimer: any = null
onMounted(() => {
  tickTimer = setInterval(() => {
    tick.value = Date.now()
  }, 30000)
})
onBeforeUnmount(() => {
  if (tickTimer) clearInterval(tickTimer)
})

const lastUpdatedLabel = computed(() => {
  if (!lastUpdatedAt.value) return ''
  const seconds = Math.floor(
    (tick.value - lastUpdatedAt.value.getTime()) / 1000
  )
  if (seconds < 5) return 'just now'
  if (seconds < 60) return `${seconds}s ago`
  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return `${minutes}m ago`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}h ago`
  return lastUpdatedAt.value.toLocaleString()
})

// --- Drill-down ---
const drillDown = (filters: Record<string, string | number | undefined>) => {
  const query: Record<string, string> = {}
  if (filters.status) query.status = String(filters.status)
  if (selectedScheme.value) query.scheme_id = String(selectedScheme.value)
  if (selectedBenefitType.value)
    query.benefit_alias = String(selectedBenefitType.value)
  if (filters.scheme_name && !query.scheme_id)
    query.search = String(filters.scheme_name)
  if (filters.benefit_alias) query.benefit_alias = String(filters.benefit_alias)
  router.push({
    name: 'group-pricing-claims-management',
    query
  })
}

// --- Reason chip sizing ---
const reasonChipSize = (count: number) => {
  const max =
    Math.max(
      ...analyticsData.value.topDeclineReasons.map((r) => r.count || 0),
      1
    ) || 1
  return 0.75 + (count / max) * 0.65
}
const reasonChipOpacity = (count: number) => {
  const max =
    Math.max(
      ...analyticsData.value.topDeclineReasons.map((r) => r.count || 0),
      1
    ) || 1
  return 0.6 + (count / max) * 0.4
}

// --- Data fetch ---
let debounceTimer: any = null
const queueLoad = () => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => loadAnalyticsData(), 300)
}

const loadAnalyticsData = async () => {
  loading.value = true
  errorMessage.value = null
  try {
    const filters: any = {
      scheme_id: selectedScheme.value || undefined,
      period: selectedPeriod.value,
      benefit_type: selectedBenefitType.value || undefined,
      limit: 10
    }
    if (selectedPeriod.value === 'custom') {
      if (startDate.value) filters.from = startDate.value
      if (endDate.value) filters.to = endDate.value
    }

    const response = await GroupPricingService.getClaimsAnalytics(filters)
    const data = response.data || {}

    const next = blankAnalytics()
    next.totalClaims = Number(data.total_claims || 0)
    next.totalPaid = Number(data.total_paid_amount || 0)
    next.avgClaimAmount = Number(data.avg_claim_amount || 0)
    next.totalExposure = Number(data.total_exposure || 0)
    next.avgProcessingDays = Number(data.avg_processing_days || 0)
    next.approvalRate = Math.round((data.approval_rate || 0) * 1000) / 10
    next.declineRate = Math.round((data.decline_rate || 0) * 1000) / 10
    next.financeRejectionRate =
      Math.round((data.finance_rejection_rate || 0) * 1000) / 10
    next.largeClaimsCount = Number(data.large_claims_count || 0)
    next.largeClaimThreshold = Number(data.large_claim_threshold || 1000000)
    next.reopenCount = Number(data.reopen_count || 0)

    if (data.processing_efficiency) {
      const e = data.processing_efficiency
      next.slaDays = Number(e.sla_days || 30)
      next.slaCompliance = Math.round((e.on_time_rate || 0) * 1000) / 10
      next.efficiencyScore = Math.round((e.efficiency_score || 0) * 1000) / 10
      next.closedInPeriod = Number(e.closed_in_period || 0)
      next.wipOpenClaims = Number(e.wip_open_claims || 0)
      next.throughputPerWeek = Number(e.throughput_per_week || 0)
      next.processingTimeMedian = Number(e.processing_time_median || 0)
      next.processingTimeDistribution = e.processing_time_distribution || []
    }

    next.byStatus = Object.entries(data.by_status || {}).map(([k, v]) => ({
      status: k,
      count: Number(v) || 0
    }))
    next.byBenefit = Object.entries(data.by_benefit || {})
      .map(([k, v]) => ({ label: k || 'Unspecified', count: Number(v) || 0 }))
      .sort((a, b) => b.count - a.count)
    next.byCauseType = (data.by_cause_type || []).map((x: any) => ({
      label: x.label,
      count: Number(x.count) || 0
    }))
    next.byMemberType = (data.by_member_type || []).map((x: any) => ({
      label: x.label,
      count: Number(x.count) || 0
    }))
    next.agingBuckets = (data.aging_buckets || []).map((x: any) => ({
      label: x.label,
      count: Number(x.count) || 0
    }))
    next.claimsTrend = (data.claims_trend || []).map((x: any) => ({
      month: x.month,
      count: Number(x.count) || 0,
      amount: Number(x.amount) || 0
    }))
    next.claimantAgeDistribution = (data.claimant_age_distribution || []).map(
      (x: any) => ({ label: x.label, count: Number(x.count) || 0 })
    )
    next.topSchemes = (data.top_schemes || []).map((x: any) => ({
      scheme_name: x.scheme_name,
      count: Number(x.count) || 0,
      total_amount: Number(x.total_amount) || 0
    }))
    next.topAssessors = (data.top_assessors || []).map((x: any) => ({
      name: x.name,
      claimsProcessed: Number(x.count) || 0
    }))
    next.topDeclineReasons = (data.topDeclineReasons || []).map((x: any) => ({
      reason: x.reason,
      count: Number(x.count) || 0
    }))

    analyticsData.value = next
    topClaims.value = (data.top_claims || []).map((c: any) => ({
      id: c.id,
      claim_number: c.claim_number,
      member_name: c.member,
      scheme_name: c.scheme,
      benefit_type: c.benefit_type,
      claim_amount: c.amount,
      status: c.status,
      date_notified: c.date_notified
    }))

    lastUpdatedAt.value = new Date()
  } catch (err: any) {
    console.error('Failed to load claims analytics', err)
    const body = err?.response?.data
    const backend =
      typeof body === 'string'
        ? body
        : body?.message || body?.error || body?.detail
    errorMessage.value =
      backend ||
      err?.message ||
      'Failed to load claims analytics.'
    analyticsData.value = blankAnalytics()
    topClaims.value = []
  } finally {
    loading.value = false
  }
}

watch(
  [selectedScheme, selectedPeriod, selectedBenefitType, startDate, endDate],
  () => {
    queueLoad()
  }
)

onMounted(() => {
  loadAnalyticsData()
})

// --- Chart instance refs ---
// AG Charts Vue exposes the AgChartInstance via `chartRef.value.chart`.
// We keep a ref per chart so the header menu can call .download().
const statusChartRef = ref<any>(null)
const agingChartRef = ref<any>(null)
const trendChartRef = ref<any>(null)
const benefitChartRef = ref<any>(null)
const causeChartRef = ref<any>(null)
const memberChartRef = ref<any>(null)
const processingChartRef = ref<any>(null)
const ageChartRef = ref<any>(null)

const downloadChartImage = (
  chartRef: any,
  fileName: string,
  fileFormat: 'png' | 'jpg'
) => {
  const inst = chartRef?.value?.chart
  if (inst && typeof inst.download === 'function') {
    inst.download({ fileName, fileFormat })
  }
}

// --- Chart download / data export helpers ---
const downloadCsv = (rows: any[], name: string) => {
  if (!rows || rows.length === 0) return
  const headers = Object.keys(rows[0])
  const escape = (v: any) => {
    if (v == null) return ''
    const s = String(v)
    return /[",\n\r]/.test(s) ? `"${s.replace(/"/g, '""')}"` : s
  }
  const csv = [
    headers.join(','),
    ...rows.map((r) => headers.map((h) => escape(r[h])).join(','))
  ].join('\n')
  const blob = new Blob(['﻿', csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${name}.csv`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

const copyTsv = async (rows: any[]) => {
  if (!rows || rows.length === 0) return
  const headers = Object.keys(rows[0])
  const tsv = [
    headers.join('\t'),
    ...rows.map((r) => headers.map((h) => (r[h] == null ? '' : String(r[h]))).join('\t'))
  ].join('\n')
  try {
    await navigator.clipboard.writeText(tsv)
  } catch {
    /* clipboard denied — silently no-op */
  }
}

// Build header-menu actions for a chart card. Used by the hamburger
// (mdi-menu) menu in the top-right of each chart's header.
const buildChartActions = (
  name: string,
  rows: any[],
  chartRef: any
) => [
  {
    label: 'Download as PNG',
    icon: 'mdi-file-image-outline',
    onClick: () => downloadChartImage(chartRef, name, 'png')
  },
  {
    label: 'Download as JPG',
    icon: 'mdi-file-image-outline',
    onClick: () => downloadChartImage(chartRef, name, 'jpg')
  },
  {
    label: 'Download data (CSV)',
    icon: 'mdi-file-delimited-outline',
    onClick: () => downloadCsv(rows, name)
  },
  {
    label: 'Copy data (TSV)',
    icon: 'mdi-content-copy',
    onClick: () => copyTsv(rows)
  }
]

// Build a context-menu config that extends the AG Charts default menu with
// chart-image and data download actions. `'defaults'` includes the built-in
// "Download" item; `'separator'` draws a divider before our custom items.
const buildContextMenu = (name: string, rows: any[]): any => ({
  enabled: true,
  items: [
    'defaults',
    'separator',
    {
      label: 'Download chart as PNG',
      action: (e: any) =>
        e?.chart?.download?.({ fileName: name, fileFormat: 'png' })
    },
    {
      label: 'Download chart as JPG',
      action: (e: any) =>
        e?.chart?.download?.({ fileName: name, fileFormat: 'jpg' })
    },
    {
      label: 'Download data as CSV',
      action: () => downloadCsv(rows, name)
    },
    {
      label: 'Copy data to clipboard (TSV)',
      action: () => copyTsv(rows)
    }
  ]
})

// Header-menu actions, one entry per chart. Keys match the section names.
const chartActions = computed(() => ({
  status: buildChartActions(
    'claims-by-status',
    analyticsData.value.byStatus.map((r) => ({
      status: formatStatusLabel(r.status),
      count: r.count
    })),
    statusChartRef
  ),
  aging: buildChartActions(
    'open-claims-aging',
    analyticsData.value.agingBuckets.map((r) => ({
      bucket_days: r.label,
      open_claims: r.count
    })),
    agingChartRef
  ),
  trend: buildChartActions(
    'claims-trend',
    analyticsData.value.claimsTrend.map((r) => ({
      month: r.month,
      claims: r.count,
      amount_zar: r.amount
    })),
    trendChartRef
  ),
  benefit: buildChartActions(
    'claims-by-benefit',
    analyticsData.value.byBenefit.map((r) => ({
      benefit: r.label,
      count: r.count
    })),
    benefitChartRef
  ),
  cause: buildChartActions(
    'claims-by-cause',
    analyticsData.value.byCauseType.map((r) => ({
      cause: r.label,
      count: r.count
    })),
    causeChartRef
  ),
  member: buildChartActions(
    'claims-by-member-type',
    analyticsData.value.byMemberType.map((r) => ({
      member_type: r.label,
      count: r.count
    })),
    memberChartRef
  ),
  processing: buildChartActions(
    'processing-time-distribution',
    analyticsData.value.processingTimeDistribution.map((r) => ({
      days_bucket: r.label,
      count: r.count
    })),
    processingChartRef
  ),
  age: buildChartActions(
    'claimant-age-distribution',
    analyticsData.value.claimantAgeDistribution.map((r) => ({
      age_band: r.label,
      count: r.count
    })),
    ageChartRef
  )
}))

// --- Chart options ---
const statusDonutOptions = computed<AgChartOptions | null>(() => {
  const rows = analyticsData.value.byStatus
  if (!rows.length) return null
  const data = rows.map((r) => ({
    status: r.status,
    label: formatStatusLabel(r.status),
    count: r.count,
    fill: statusPalette.value[r.status] || palette.value.muted
  }))
  return {
    data,
    series: [
      {
        type: 'donut',
        angleKey: 'count',
        calloutLabelKey: 'label',
        sectorLabelKey: 'count',
        innerRadiusRatio: 0.55,
        fills: data.map((d) => d.fill),
        strokes: ['#ffffff'],
        strokeWidth: 2,
        sectorLabel: {
          color: '#ffffff',
          fontWeight: 'bold',
          fontSize: 11
        },
        listeners: {
          nodeClick: (event: any) => {
            const status = event.datum?.status
            if (status) drillDown({ status })
          }
        }
      } as any
    ],
    legend: { position: 'right', enabled: true },
    background: { fill: 'transparent' },
    contextMenu: buildContextMenu(
      'claims-by-status',
      data.map((d) => ({ status: d.label, count: d.count }))
    )
  }
})

const agingChartOptions = computed<AgChartOptions | null>(() => {
  const rows = analyticsData.value.agingBuckets
  if (!rows.length || rows.every((r) => r.count === 0)) return null
  // Tone severity by bucket index
  const tones = [
    palette.value.info,
    palette.value.info,
    palette.value.warning,
    palette.value.warning,
    palette.value.error
  ]
  const data = rows.map((r, i) => ({
    label: r.label,
    count: r.count,
    fill: tones[i] || palette.value.muted
  }))
  return {
    data,
    series: [
      {
        type: 'bar',
        direction: 'horizontal',
        xKey: 'label',
        yKey: 'count',
        yName: 'Open claims',
        fills: data.map((d) => d.fill),
        cornerRadius: 4,
        label: { enabled: true, color: '#ffffff' }
      } as any
    ],
    axes: [
      {
        type: 'category',
        position: 'left',
        title: { text: 'Age bucket (days)' }
      },
      { type: 'number', position: 'bottom', title: { text: 'Open claims' } }
    ],
    legend: { enabled: false },
    background: { fill: 'transparent' },
    contextMenu: buildContextMenu(
      'open-claims-aging',
      rows.map((r) => ({ bucket_days: r.label, open_claims: r.count }))
    )
  }
})

const trendChartOptions = computed<AgChartOptions | null>(() => {
  const rows = analyticsData.value.claimsTrend
  if (!rows.length) return null
  const data = rows.map((r) => {
    const d = new Date(r.month + '-01')
    return {
      month: d.toLocaleDateString('en-ZA', { month: 'short', year: '2-digit' }),
      count: r.count,
      amount: r.amount / 1_000_000
    }
  })
  return {
    data,
    series: [
      {
        type: 'bar',
        xKey: 'month',
        yKey: 'count',
        yName: 'Claims',
        fill: palette.value.primary,
        cornerRadius: 4
      } as any,
      {
        type: 'line',
        xKey: 'month',
        yKey: 'amount',
        yName: 'Amount (R millions)',
        stroke: palette.value.accent,
        strokeWidth: 3,
        marker: { fill: palette.value.accent, size: 8, strokeWidth: 0 }
      } as any
    ],
    axes: [
      { type: 'category', position: 'bottom' },
      {
        type: 'number',
        position: 'left',
        keys: ['count'],
        title: { text: 'Claims' }
      },
      {
        type: 'number',
        position: 'right',
        keys: ['amount'],
        title: { text: 'Amount (R millions)' },
        label: {
          formatter: (p: any) => `R${Number(p.value).toFixed(1)}M`
        }
      }
    ],
    legend: { position: 'top', enabled: true },
    background: { fill: 'transparent' },
    contextMenu: buildContextMenu(
      'claims-trend',
      analyticsData.value.claimsTrend.map((r) => ({
        month: r.month,
        claims: r.count,
        amount_zar: r.amount
      }))
    )
  }
})

const benefitChartOptions = computed<AgChartOptions | null>(() => {
  const rows = analyticsData.value.byBenefit
  if (!rows.length) return null
  return {
    data: rows,
    series: [
      {
        type: 'bar',
        direction: 'horizontal',
        xKey: 'label',
        yKey: 'count',
        yName: 'Claims',
        fill: palette.value.accent,
        cornerRadius: 4,
        listeners: {
          nodeClick: (event: any) => {
            const label = event.datum?.label
            if (label) drillDown({ benefit_alias: label })
          }
        }
      } as any
    ],
    axes: [
      { type: 'category', position: 'left' },
      { type: 'number', position: 'bottom' }
    ],
    legend: { enabled: false },
    background: { fill: 'transparent' },
    contextMenu: buildContextMenu(
      'claims-by-benefit',
      rows.map((r) => ({ benefit: r.label, count: r.count }))
    )
  }
})

const causeChartOptions = computed<AgChartOptions | null>(() => {
  const rows = analyticsData.value.byCauseType
  if (!rows.length) return null
  return {
    data: rows,
    series: [
      {
        type: 'bar',
        direction: 'horizontal',
        xKey: 'label',
        yKey: 'count',
        yName: 'Claims',
        fill: palette.value.info,
        cornerRadius: 4
      } as any
    ],
    axes: [
      { type: 'category', position: 'left' },
      { type: 'number', position: 'bottom' }
    ],
    legend: { enabled: false },
    background: { fill: 'transparent' },
    contextMenu: buildContextMenu(
      'claims-by-cause',
      rows.map((r) => ({ cause: r.label, count: r.count }))
    )
  }
})

const memberChartOptions = computed<AgChartOptions | null>(() => {
  const rows = analyticsData.value.byMemberType
  if (!rows.length) return null
  const fills = [
    palette.value.primary,
    palette.value.accent,
    palette.value.info,
    palette.value.warning,
    palette.value.success,
    palette.value.error
  ]
  return {
    data: rows.map((r, i) => ({
      ...r,
      fill: fills[i % fills.length]
    })),
    series: [
      {
        type: 'donut',
        angleKey: 'count',
        calloutLabelKey: 'label',
        sectorLabelKey: 'count',
        innerRadiusRatio: 0.55,
        fills,
        strokes: ['#ffffff'],
        strokeWidth: 2,
        sectorLabel: { color: '#ffffff', fontWeight: 'bold', fontSize: 11 }
      } as any
    ],
    legend: { position: 'right', enabled: true },
    background: { fill: 'transparent' },
    contextMenu: buildContextMenu(
      'claims-by-member-type',
      rows.map((r) => ({ member_type: r.label, count: r.count }))
    )
  }
})

const processingChartOptions = computed<AgChartOptions | null>(() => {
  const rows = analyticsData.value.processingTimeDistribution
  if (!rows.length || rows.every((r) => r.count === 0)) return null
  return {
    data: rows,
    series: [
      {
        type: 'bar',
        xKey: 'label',
        yKey: 'count',
        yName: 'Claims',
        fill: palette.value.warning,
        cornerRadius: 4
      } as any
    ],
    axes: [
      { type: 'category', position: 'bottom', title: { text: 'Days bucket' } },
      { type: 'number', position: 'left', title: { text: 'Claims' } }
    ],
    legend: { enabled: false },
    background: { fill: 'transparent' },
    contextMenu: buildContextMenu(
      'processing-time-distribution',
      rows.map((r) => ({ days_bucket: r.label, count: r.count }))
    )
  }
})

const ageChartOptions = computed<AgChartOptions | null>(() => {
  const rows = analyticsData.value.claimantAgeDistribution
  if (!rows.length || rows.every((r) => r.count === 0)) return null
  return {
    data: rows,
    series: [
      {
        type: 'bar',
        xKey: 'label',
        yKey: 'count',
        yName: 'Claims',
        fill: palette.value.primary,
        cornerRadius: 4
      } as any
    ],
    axes: [
      { type: 'category', position: 'bottom', title: { text: 'Age band' } },
      { type: 'number', position: 'left', title: { text: 'Claims' } }
    ],
    legend: { enabled: false },
    background: { fill: 'transparent' },
    contextMenu: buildContextMenu(
      'claimant-age-distribution',
      rows.map((r) => ({ age_band: r.label, count: r.count }))
    )
  }
})

// --- Top claims grid columns ---
const topClaimsColumns = [
  {
    field: 'claim_number',
    headerName: 'Claim #',
    sortable: true,
    filter: true,
    width: 140
  },
  {
    field: 'member_name',
    headerName: 'Member',
    sortable: true,
    filter: true,
    flex: 1
  },
  {
    field: 'scheme_name',
    headerName: 'Scheme',
    sortable: true,
    filter: true,
    flex: 1
  },
  {
    field: 'benefit_type',
    headerName: 'Benefit',
    sortable: true,
    filter: true,
    width: 110
  },
  {
    field: 'claim_amount',
    headerName: 'Amount',
    sortable: true,
    filter: true,
    width: 140,
    type: 'numericColumn',
    cellRenderer: (p: any) => (p.value != null ? formatCurrency(p.value) : '')
  },
  {
    field: 'status',
    headerName: 'Status',
    sortable: true,
    filter: true,
    width: 150,
    cellRenderer: (p: any) => {
      if (!p.value) return ''
      const color = statusPalette.value[p.value] || palette.value.muted
      const label = formatStatusLabel(p.value)
      return `<span style="background:${color};color:#fff;padding:2px 10px;border-radius:12px;font-size:11px;font-weight:500;">${label}</span>`
    }
  },
  {
    field: 'date_notified',
    headerName: 'Notified',
    sortable: true,
    filter: true,
    width: 130
  }
]

// --- Exports ---
const exportPdf = async () => {
  if (!dashboardRoot.value) return
  exporting.value = true
  try {
    const canvas = await html2canvas(dashboardRoot.value, {
      backgroundColor: '#ffffff',
      scale: 2,
      useCORS: true,
      windowWidth: dashboardRoot.value.scrollWidth,
      windowHeight: dashboardRoot.value.scrollHeight
    })
    const imgData = canvas.toDataURL('image/jpeg', 0.92)
    const pdf = new JsPdf({
      orientation: 'landscape',
      unit: 'pt',
      format: 'a3'
    })
    const pageWidth = pdf.internal.pageSize.getWidth()
    const pageHeight = pdf.internal.pageSize.getHeight()
    const imgWidth = pageWidth - 40
    const imgHeight = (canvas.height * imgWidth) / canvas.width
    const topMargin = 20
    let remaining = imgHeight
    if (imgHeight <= pageHeight - 40) {
      pdf.addImage(imgData, 'JPEG', 20, topMargin, imgWidth, imgHeight)
    } else {
      let y = 0
      while (remaining > 0) {
        pdf.addImage(imgData, 'JPEG', 20, topMargin - y, imgWidth, imgHeight)
        remaining -= pageHeight - 40
        y += pageHeight - 40
        if (remaining > 0) pdf.addPage('a3', 'landscape')
      }
    }
    const ts = new Date().toISOString().replace(/[:T]/g, '-').split('.')[0]
    pdf.save(`claims-analytics_${ts}.pdf`)
  } catch (err) {
    console.error('PDF export failed', err)
    errorMessage.value = 'Could not generate PDF export.'
  } finally {
    exporting.value = false
  }
}

const exportExcel = () => {
  exporting.value = true
  try {
    const wb = XLSX.utils.book_new()

    const kpis = [
      ['Period', periodLabel.value],
      [
        'Last updated',
        lastUpdatedAt.value ? lastUpdatedAt.value.toISOString() : ''
      ],
      [],
      ['Total claims', analyticsData.value.totalClaims],
      ['Total paid', analyticsData.value.totalPaid],
      ['Avg claim amount', analyticsData.value.avgClaimAmount],
      ['Outstanding exposure', analyticsData.value.totalExposure],
      ['Approval rate (%)', analyticsData.value.approvalRate],
      ['Decline rate (%)', analyticsData.value.declineRate],
      ['Finance rejection rate (%)', analyticsData.value.financeRejectionRate],
      ['Avg processing days', analyticsData.value.avgProcessingDays],
      ['SLA compliance (%)', analyticsData.value.slaCompliance],
      ['SLA days', analyticsData.value.slaDays],
      ['WIP open claims', analyticsData.value.wipOpenClaims],
      ['Throughput per week', analyticsData.value.throughputPerWeek],
      ['Closed in period', analyticsData.value.closedInPeriod],
      ['Large claims count', analyticsData.value.largeClaimsCount],
      ['Large claim threshold', analyticsData.value.largeClaimThreshold],
      ['Reopened claims', analyticsData.value.reopenCount]
    ]
    XLSX.utils.book_append_sheet(wb, XLSX.utils.aoa_to_sheet(kpis), 'KPIs')

    const sheet = (name: string, rows: any[]) => {
      if (!rows.length) return
      XLSX.utils.book_append_sheet(
        wb,
        XLSX.utils.json_to_sheet(rows),
        name.slice(0, 31)
      )
    }
    sheet('By Status', analyticsData.value.byStatus)
    sheet('By Benefit', analyticsData.value.byBenefit)
    sheet('By Cause', analyticsData.value.byCauseType)
    sheet('By Member Type', analyticsData.value.byMemberType)
    sheet('Aging Buckets', analyticsData.value.agingBuckets)
    sheet('Claims Trend', analyticsData.value.claimsTrend)
    sheet('Age Distribution', analyticsData.value.claimantAgeDistribution)
    sheet('Processing Time', analyticsData.value.processingTimeDistribution)
    sheet('Top Schemes', analyticsData.value.topSchemes)
    sheet('Top Assessors', analyticsData.value.topAssessors)
    sheet('Decline Reasons', analyticsData.value.topDeclineReasons)
    sheet('Largest Claims', topClaims.value)

    const ts = new Date().toISOString().replace(/[:T]/g, '-').split('.')[0]
    XLSX.writeFile(wb, `claims-analytics_${ts}.xlsx`)
  } catch (err) {
    console.error('Excel export failed', err)
    errorMessage.value = 'Could not generate Excel export.'
  } finally {
    exporting.value = false
  }
}
</script>

<style scoped>
.claims-analytics {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.claims-analytics__controls {
  position: sticky;
  top: 0;
  z-index: 4;
  background: rgb(var(--v-theme-surface));
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  border-radius: 10px;
  padding: 10px 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.claims-analytics__filters {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  flex: 1 1 auto;
  min-width: 0;
}

.filter-input {
  width: 170px;
  min-width: 140px;
}

.filter-input--wide {
  width: 220px;
}

.claims-analytics__actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.claims-analytics__updated {
  font-size: 0.75rem;
  color: rgba(var(--v-theme-on-surface), 0.55);
  margin-right: 4px;
}

.kpi-strip {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

@media (max-width: 1100px) {
  .kpi-strip {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 600px) {
  .kpi-strip {
    grid-template-columns: 1fr;
  }
}

.kpi-clickable {
  cursor: pointer;
}

.analytics-row {
  margin: 0;
}

.ag-fill {
  height: 100%;
  width: 100%;
  display: block;
}

/* Top schemes / assessors list */
.top-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 4px 0;
}

.top-list__row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 4px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.15s ease;
}

.top-list__row:hover {
  background: rgba(var(--v-theme-on-surface), 0.04);
}

.top-list__row--static {
  cursor: default;
}

.top-list__row--static:hover {
  background: transparent;
}

.top-list__rank {
  font-size: 0.75rem;
  font-weight: 700;
  color: rgba(var(--v-theme-on-surface), 0.55);
  background: rgba(var(--v-theme-on-surface), 0.08);
  width: 22px;
  height: 22px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.top-list__rank--accent {
  color: #fff;
  background: rgb(var(--v-theme-accent));
}

.top-list__main {
  flex: 1;
  min-width: 0;
}

.top-list__title {
  font-size: 0.85rem;
  font-weight: 600;
  color: rgba(var(--v-theme-on-surface), 0.9);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.top-list__subtitle {
  font-size: 0.72rem;
  color: rgba(var(--v-theme-on-surface), 0.55);
  margin-top: 2px;
}

.top-list__bar {
  margin-top: 6px;
  height: 4px;
  width: 100%;
  background: rgba(var(--v-theme-on-surface), 0.08);
  border-radius: 2px;
  overflow: hidden;
}

.top-list__bar-fill {
  height: 100%;
  background: rgb(var(--v-theme-primary));
  transition: width 0.3s ease;
}

.top-list__bar-fill--accent {
  background: rgb(var(--v-theme-accent));
}

/* Decline reason cloud */
.reason-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 8px 4px;
  align-items: center;
}

.reason-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 14px;
  background: rgba(var(--v-theme-error), 0.08);
  color: rgb(var(--v-theme-error));
  border: 1px solid rgba(var(--v-theme-error), 0.18);
  font-weight: 500;
  line-height: 1.1;
}

.reason-chip__count {
  font-size: 0.7rem;
  padding: 0 6px;
  background: rgb(var(--v-theme-error));
  color: #fff;
  border-radius: 8px;
  font-weight: 700;
}

@media print {
  .no-print {
    display: none !important;
  }
  .claims-analytics {
    gap: 8px;
  }
}
</style>
