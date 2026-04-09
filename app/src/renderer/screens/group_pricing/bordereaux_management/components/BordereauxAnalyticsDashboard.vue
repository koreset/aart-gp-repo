<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center justify-space-between">
              <div>
                <span class="headline">South African Group Risk Analytics</span>
                <p class="text-subtitle-1 text-medium-emphasis mt-2">
                  Comprehensive group risk metrics, premium analytics, and
                  regulatory compliance insights
                </p>
              </div>
              <div class="d-flex align-center gap-2">
                <v-select
                  v-model="selectedPeriod"
                  :items="periodOptions"
                  label="Period"
                  variant="outlined"
                  density="compact"
                  @update:model-value="loadAnalytics"
                />
                <v-btn
                  color="success"
                  variant="outlined"
                  prepend-icon="mdi-export"
                  @click="exportReport"
                >
                  Export Report
                </v-btn>
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
            <!-- Key Financial Metrics -->
            <v-row class="mb-6">
              <v-col cols="12" sm="6" lg="3">
                <v-card variant="outlined" class="h-100">
                  <v-card-text>
                    <div class="d-flex align-center justify-space-between">
                      <div>
                        <p class="text-caption text-medium-emphasis mb-1"
                          >Total Premium Volume</p
                        >
                        <p class="text-h4 font-weight-bold text-primary">{{
                          formatCurrency(metrics.totalPremiumVolume)
                        }}</p>
                        <p
                          class="text-caption"
                          :class="getChangeClass(metrics.premiumVolumeChange)"
                        >
                          {{ formatChange(metrics.premiumVolumeChange) }}% vs
                          last period
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
                        <p class="text-caption text-medium-emphasis mb-1"
                          >Loss Ratio</p
                        >
                        <p class="text-h4 font-weight-bold text-success"
                          >{{ metrics.lossRatio }}%</p
                        >
                        <p
                          class="text-caption"
                          :class="getChangeClass(-metrics.lossRatioChange)"
                        >
                          {{ formatChange(metrics.lossRatioChange) }}% vs last
                          period
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
                        <p class="text-caption text-medium-emphasis mb-1"
                          >Claims Frequency</p
                        >
                        <p class="text-h4 font-weight-bold text-info"
                          >{{ metrics.claimsFrequency }} per 1000</p
                        >
                        <p
                          class="text-caption"
                          :class="getChangeClass(-metrics.frequencyChange)"
                        >
                          {{ formatChange(metrics.frequencyChange) }} vs last
                          period
                        </p>
                      </div>
                      <v-icon size="50" color="info"
                        >mdi-chart-timeline-variant</v-icon
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
                        <p class="text-caption text-medium-emphasis mb-1"
                          >Avg Claim Size</p
                        >
                        <p class="text-h4 font-weight-bold text-purple">{{
                          formatCurrency(metrics.avgClaimSize)
                        }}</p>
                        <p
                          class="text-caption"
                          :class="getChangeClass(metrics.claimSizeChange)"
                        >
                          {{ formatChange(metrics.claimSizeChange) }}% vs last
                          period
                        </p>
                      </div>
                      <v-icon size="50" color="purple">mdi-calculator</v-icon>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Premium & Member Analytics -->
            <v-row class="mb-6">
              <v-col cols="12" lg="8">
                <v-card variant="outlined">
                  <v-card-title class="text-h6 font-weight-bold">
                    Premium per Member per Month (ZAR)
                  </v-card-title>
                  <v-card-text>
                    <div ref="premiumTrendChart" style="height: 300px">
                      <div
                        class="d-flex align-center justify-center h-100 bg-grey-lighten-5 rounded"
                      >
                        <div class="text-center">
                          <v-icon size="60" color="grey">mdi-chart-line</v-icon>
                          <p class="text-grey mt-2"
                            >Premium trend analysis by benefit type</p
                          >
                        </div>
                      </div>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" lg="4">
                <v-card variant="outlined">
                  <v-card-title class="text-h6 font-weight-bold">
                    Benefit Mix Distribution
                  </v-card-title>
                  <v-card-text>
                    <div ref="benefitMixChart" style="height: 300px">
                      <div
                        class="d-flex align-center justify-center h-100 bg-grey-lighten-5 rounded"
                      >
                        <div class="text-center">
                          <v-icon size="60" color="grey"
                            >mdi-chart-donut</v-icon
                          >
                          <p class="text-grey mt-2"
                            >Benefit distribution by premium volume</p
                          >
                        </div>
                      </div>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Benefit Performance Analysis -->
            <v-row class="mb-6">
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title class="text-h6 font-weight-bold">
                    Benefit Performance Analysis
                  </v-card-title>
                  <v-card-text>
                    <v-data-table
                      :headers="benefitHeaders"
                      :items="benefitPerformance"
                      :items-per-page="10"
                      density="compact"
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

                      <template #[`item.premium_volume`]="{ item }">
                        {{ formatCurrency(item.premium_volume) }}
                      </template>

                      <template #[`item.member_count`]="{ item }">
                        {{ item.member_count.toLocaleString() }}
                      </template>

                      <template #[`item.avg_sum_assured`]="{ item }">
                        {{ formatCurrency(item.avg_sum_assured) }}
                      </template>

                      <template #[`item.claim_ratio`]="{ item }">
                        <div class="d-flex align-center">
                          <v-progress-linear
                            :model-value="item.claim_ratio"
                            :color="getLossRatioColor(item.claim_ratio)"
                            height="6"
                            rounded
                            class="me-2"
                            style="width: 60px"
                          />
                          <span>{{ item.claim_ratio }}%</span>
                        </div>
                      </template>

                      <template #[`item.unit_rate`]="{ item }">
                        {{ item.unit_rate }} per 1000
                      </template>
                    </v-data-table>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Insurer Performance & Market Share -->
            <v-row class="mb-6">
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title class="text-h6 font-weight-bold">
                    Insurer Performance & Market Share
                  </v-card-title>
                  <v-card-text>
                    <v-data-table
                      :headers="insurerHeaders"
                      :items="insurerMetrics"
                      :items-per-page="10"
                      density="compact"
                    >
                      <template #[`item.insurer_name`]="{ item }">
                        <div class="d-flex align-center">
                          <v-avatar size="32" class="me-3" color="primary">
                            {{ item.insurer_name.charAt(0) }}
                          </v-avatar>
                          {{ item.insurer_name }}
                        </div>
                      </template>

                      <template #[`item.market_share`]="{ item }">
                        <div class="d-flex align-center">
                          <v-progress-linear
                            :model-value="item.market_share"
                            color="primary"
                            height="6"
                            rounded
                            class="me-2"
                            style="width: 60px"
                          />
                          <span>{{ item.market_share }}%</span>
                        </div>
                      </template>

                      <template #[`item.premium_volume`]="{ item }">
                        {{ formatCurrency(item.premium_volume) }}
                      </template>

                      <template #[`item.member_count`]="{ item }">
                        {{ item.member_count.toLocaleString() }}
                      </template>

                      <template #[`item.loss_ratio`]="{ item }">
                        <v-chip
                          :color="getLossRatioColor(item.loss_ratio)"
                          size="small"
                        >
                          {{ item.loss_ratio }}%
                        </v-chip>
                      </template>

                      <template #[`item.avg_processing_hours`]="{ item }">
                        <span
                          :class="
                            getProcessingTimeClass(item.avg_processing_hours)
                          "
                        >
                          {{ item.avg_processing_hours }}h
                        </span>
                      </template>

                      <template #[`item.sla_compliance`]="{ item }">
                        <v-chip
                          :color="getSLAColor(item.sla_compliance)"
                          size="small"
                        >
                          {{ item.sla_compliance }}%
                        </v-chip>
                      </template>
                    </v-data-table>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- SA-Specific Error Analysis & Compliance -->
            <v-row class="mb-6">
              <v-col cols="12" md="6">
                <v-card variant="outlined">
                  <v-card-title class="text-h6 font-weight-bold text-error">
                    <v-icon class="me-2">mdi-alert-circle</v-icon>
                    SA-Specific Data Issues
                  </v-card-title>
                  <v-card-text>
                    <v-list>
                      <v-list-item
                        v-for="error in saSpecificErrors"
                        :key="error.type"
                        class="px-0"
                      >
                        <template #prepend>
                          <v-avatar
                            size="32"
                            :color="getErrorSeverityColor(error.severity)"
                          >
                            {{ error.count }}
                          </v-avatar>
                        </template>
                        <v-list-item-title>{{
                          error.description
                        }}</v-list-item-title>
                        <v-list-item-subtitle>
                          {{ error.type }} • {{ error.frequency }}% of failures
                          • {{ error.impact }}
                        </v-list-item-subtitle>
                        <template #append>
                          <v-chip
                            :color="getErrorSeverityColor(error.severity)"
                            size="small"
                            variant="tonal"
                          >
                            {{ error.severity }}
                          </v-chip>
                        </template>
                      </v-list-item>
                    </v-list>
                  </v-card-text>
                </v-card>
              </v-col>

              <!-- Enhanced SA Compliance Overview -->
              <v-col cols="12" md="6">
                <v-card variant="outlined">
                  <v-card-title class="text-h6 font-weight-bold text-teal">
                    <v-icon class="me-2">mdi-shield-check</v-icon>
                    SA Regulatory Compliance
                  </v-card-title>
                  <v-card-text>
                    <v-list>
                      <v-list-item class="px-0">
                        <template #prepend>
                          <v-icon color="success">mdi-check-circle</v-icon>
                        </template>
                        <v-list-item-title
                          >FSP Act Compliance</v-list-item-title
                        >
                        <v-list-item-subtitle
                          >Financial Services Provider Act
                          requirements</v-list-item-subtitle
                        >
                        <template #append>
                          <v-chip color="success" size="small"
                            >{{ compliance.fspCompliance }}%</v-chip
                          >
                        </template>
                      </v-list-item>

                      <v-list-item class="px-0">
                        <template #prepend>
                          <v-icon color="success">mdi-check-circle</v-icon>
                        </template>
                        <v-list-item-title
                          >SARS Tax Compliance</v-list-item-title
                        >
                        <v-list-item-subtitle
                          >Premium documentation for tax
                          purposes</v-list-item-subtitle
                        >
                        <template #append>
                          <v-chip color="success" size="small"
                            >{{ compliance.sarsTaxCompliance }}%</v-chip
                          >
                        </template>
                      </v-list-item>

                      <v-list-item class="px-0">
                        <template #prepend>
                          <v-icon
                            :color="
                              compliance.popiaCover >= 95
                                ? 'success'
                                : 'warning'
                            "
                          >
                            {{
                              compliance.popiaCover >= 95
                                ? 'mdi-check-circle'
                                : 'mdi-clock'
                            }}
                          </v-icon>
                        </template>
                        <v-list-item-title
                          >POPIA Data Protection</v-list-item-title
                        >
                        <v-list-item-subtitle
                          >Personal data protection
                          measures</v-list-item-subtitle
                        >
                        <template #append>
                          <v-chip
                            :color="
                              compliance.popiaCover >= 95
                                ? 'success'
                                : 'warning'
                            "
                            size="small"
                          >
                            {{ compliance.popiaCover }}%
                          </v-chip>
                        </template>
                      </v-list-item>

                      <v-list-item class="px-0">
                        <template #prepend>
                          <v-icon
                            :color="
                              compliance.idValidationRate >= 90
                                ? 'success'
                                : 'warning'
                            "
                          >
                            {{
                              compliance.idValidationRate >= 90
                                ? 'mdi-check-circle'
                                : 'mdi-alert'
                            }}
                          </v-icon>
                        </template>
                        <v-list-item-title
                          >SA ID Validation Rate</v-list-item-title
                        >
                        <v-list-item-subtitle
                          >South African ID number validation
                          success</v-list-item-subtitle
                        >
                        <template #append>
                          <v-chip
                            :color="
                              compliance.idValidationRate >= 90
                                ? 'success'
                                : 'warning'
                            "
                            size="small"
                          >
                            {{ compliance.idValidationRate }}%
                          </v-chip>
                        </template>
                      </v-list-item>

                      <v-list-item class="px-0">
                        <template #prepend>
                          <v-icon color="info">mdi-bank</v-icon>
                        </template>
                        <v-list-item-title
                          >Banking Details Accuracy</v-list-item-title
                        >
                        <v-list-item-subtitle
                          >SA banking validation for premium
                          collection</v-list-item-subtitle
                        >
                        <template #append>
                          <v-chip color="info" size="small"
                            >{{ compliance.bankingAccuracy }}%</v-chip
                          >
                        </template>
                      </v-list-item>
                    </v-list>
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
                      density="compact"
                    >
                      <template #[`item.month`]="{ item }">
                        <span class="font-weight-bold">{{ item.month }}</span>
                      </template>

                      <template #[`item.premium_volume`]="{ item }">
                        {{ formatCurrency(item.premium_volume) }}
                      </template>

                      <template #[`item.premium_per_member`]="{ item }">
                        {{ formatCurrency(item.premium_per_member) }}
                      </template>

                      <template #[`item.loss_ratio`]="{ item }">
                        <div class="d-flex align-center">
                          <v-progress-linear
                            :model-value="item.loss_ratio"
                            :color="getLossRatioColor(item.loss_ratio)"
                            height="6"
                            rounded
                            class="me-2"
                            style="width: 60px"
                          />
                          <span>{{ item.loss_ratio }}%</span>
                        </div>
                      </template>

                      <template #[`item.member_count`]="{ item }">
                        {{ item.member_count.toLocaleString() }}
                      </template>

                      <template #[`item.processing_time`]="{ item }">
                        <span
                          :class="getProcessingTimeClass(item.processing_time)"
                        >
                          {{ item.processing_time }}h
                        </span>
                      </template>

                      <template #[`item.compliance_score`]="{ item }">
                        <v-chip
                          :color="getComplianceColor(item.compliance_score)"
                          size="small"
                        >
                          {{ item.compliance_score }}%
                        </v-chip>
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
import { ref, onMounted } from 'vue'
import BaseCard from '@/renderer/components/BaseCard.vue'

// Reactive data
const selectedPeriod = ref('last_30_days')
const premiumTrendChart = ref<HTMLElement>()
const benefitMixChart = ref<HTMLElement>()

// Enhanced metrics for SA Group Risk
const metrics = ref({
  totalPremiumVolume: 28400000,
  premiumVolumeChange: 8.7,
  lossRatio: 67.8,
  lossRatioChange: -2.1,
  claimsFrequency: 2.4,
  frequencyChange: 0.3,
  avgClaimSize: 142500,
  claimSizeChange: 5.2
})

// SA Compliance metrics
const compliance = ref({
  fspCompliance: 100.0,
  sarsTaxCompliance: 96.8,
  popiaCover: 98.5,
  idValidationRate: 92.3,
  bankingAccuracy: 94.7
})

const periodOptions = [
  { title: 'Last 7 Days', value: 'last_7_days' },
  { title: 'Last 30 Days', value: 'last_30_days' },
  { title: 'Last Quarter', value: 'last_quarter' },
  { title: 'Last Year', value: 'last_year' },
  { title: 'Year to Date', value: 'ytd' }
]

// Enhanced table headers
const benefitHeaders = [
  { title: 'Benefit Type', key: 'benefit', sortable: true },
  { title: 'Premium Volume', key: 'premium_volume', sortable: true },
  { title: 'Member Count', key: 'member_count', sortable: true },
  { title: 'Avg Sum Assured', key: 'avg_sum_assured', sortable: true },
  { title: 'Loss Ratio', key: 'claim_ratio', sortable: true },
  { title: 'Unit Rate', key: 'unit_rate', sortable: true }
]

const insurerHeaders = [
  { title: 'Insurer', key: 'insurer_name', sortable: true },
  { title: 'Market Share', key: 'market_share', sortable: true },
  { title: 'Premium Volume', key: 'premium_volume', sortable: true },
  { title: 'Members', key: 'member_count', sortable: true },
  { title: 'Loss Ratio', key: 'loss_ratio', sortable: true },
  { title: 'Processing Time', key: 'avg_processing_hours', sortable: true },
  { title: 'SLA Compliance', key: 'sla_compliance', sortable: true }
]

const kpiHeaders = [
  { title: 'Month', key: 'month', sortable: true },
  { title: 'Premium Volume', key: 'premium_volume', sortable: true },
  { title: 'Premium/Member', key: 'premium_per_member', sortable: true },
  { title: 'Loss Ratio', key: 'loss_ratio', sortable: true },
  { title: 'Members', key: 'member_count', sortable: true },
  { title: 'Processing Time', key: 'processing_time', sortable: true },
  { title: 'Compliance Score', key: 'compliance_score', sortable: true }
]

// SA Group Risk benefit performance data
const benefitPerformance = ref([
  {
    benefit: 'Group Life Assurance',
    premium_volume: 18500000,
    member_count: 12847,
    avg_sum_assured: 485000,
    claim_ratio: 65.2,
    unit_rate: 3.4
  },
  {
    benefit: 'Critical Illness',
    premium_volume: 4200000,
    member_count: 8534,
    avg_sum_assured: 325000,
    claim_ratio: 78.9,
    unit_rate: 5.8
  },
  {
    benefit: 'Permanent Total Disability',
    premium_volume: 3100000,
    member_count: 7821,
    avg_sum_assured: 280000,
    claim_ratio: 72.1,
    unit_rate: 4.2
  },
  {
    benefit: 'Temporary Total Disability',
    premium_volume: 1800000,
    member_count: 5642,
    avg_sum_assured: 180000,
    claim_ratio: 89.3,
    unit_rate: 8.1
  },
  {
    benefit: 'Funeral Benefit',
    premium_volume: 800000,
    member_count: 15200,
    avg_sum_assured: 25000,
    claim_ratio: 45.3,
    unit_rate: 2.1
  }
])

// SA Insurer performance data
const insurerMetrics = ref([
  {
    insurer_name: 'Old Mutual',
    market_share: 28.4,
    premium_volume: 8400000,
    member_count: 4250,
    loss_ratio: 68.5,
    avg_processing_hours: 16.2,
    sla_compliance: 94.2
  },
  {
    insurer_name: 'Liberty Life',
    market_share: 24.1,
    premium_volume: 7100000,
    member_count: 3890,
    loss_ratio: 72.1,
    avg_processing_hours: 14.8,
    sla_compliance: 91.7
  },
  {
    insurer_name: 'Momentum',
    market_share: 18.7,
    premium_volume: 5500000,
    member_count: 3200,
    loss_ratio: 75.8,
    avg_processing_hours: 19.3,
    sla_compliance: 88.4
  },
  {
    insurer_name: 'Discovery Life',
    market_share: 15.2,
    premium_volume: 4500000,
    member_count: 2800,
    loss_ratio: 69.2,
    avg_processing_hours: 13.5,
    sla_compliance: 96.1
  },
  {
    insurer_name: 'Sanlam',
    market_share: 13.6,
    premium_volume: 4000000,
    member_count: 2500,
    loss_ratio: 71.6,
    avg_processing_hours: 17.8,
    sla_compliance: 89.7
  }
])

// SA-specific error analysis
const saSpecificErrors = ref([
  {
    type: 'SA_ID_VALIDATION',
    description: 'Invalid South African ID numbers',
    count: 23,
    frequency: 28.4,
    severity: 'high',
    impact: 'Prevents policy issuance'
  },
  {
    type: 'BANKING_DETAILS',
    description: 'Invalid SA banking details',
    count: 18,
    frequency: 22.2,
    severity: 'high',
    impact: 'Prevents premium collection'
  },
  {
    type: 'BENEFIT_LIMITS',
    description: 'Benefits exceed regulatory limits',
    count: 15,
    frequency: 18.5,
    severity: 'medium',
    impact: 'Compliance risk'
  },
  {
    type: 'TAX_DIRECTIVE',
    description: 'Missing SARS tax directive',
    count: 12,
    frequency: 14.8,
    severity: 'medium',
    impact: 'Tax compliance risk'
  },
  {
    type: 'FORMAT_ERROR',
    description: 'Incorrect bordereaux format',
    count: 8,
    frequency: 9.9,
    severity: 'low',
    impact: 'Processing delays'
  },
  {
    type: 'FSP_DISCLOSURE',
    description: 'Missing FSP disclosures',
    count: 5,
    frequency: 6.2,
    severity: 'medium',
    impact: 'Regulatory compliance'
  }
])

// Monthly KPI data
const monthlyKPIs = ref([
  {
    month: 'Dec 2024',
    premium_volume: 28400000,
    premium_per_member: 485.5,
    loss_ratio: 67.8,
    member_count: 48750,
    processing_time: 16.2,
    compliance_score: 94.7
  },
  {
    month: 'Nov 2024',
    premium_volume: 26800000,
    premium_per_member: 478.2,
    loss_ratio: 69.2,
    member_count: 46820,
    processing_time: 17.8,
    compliance_score: 92.3
  },
  {
    month: 'Oct 2024',
    premium_volume: 25200000,
    premium_per_member: 472.8,
    loss_ratio: 71.5,
    member_count: 44390,
    processing_time: 18.9,
    compliance_score: 91.1
  },
  {
    month: 'Sep 2024',
    premium_volume: 24600000,
    premium_per_member: 468.3,
    loss_ratio: 73.1,
    member_count: 43750,
    processing_time: 20.1,
    compliance_score: 89.8
  },
  {
    month: 'Aug 2024',
    premium_volume: 23800000,
    premium_per_member: 461.9,
    loss_ratio: 74.8,
    member_count: 42980,
    processing_time: 21.4,
    compliance_score: 88.5
  },
  {
    month: 'Jul 2024',
    premium_volume: 22900000,
    premium_per_member: 455.6,
    loss_ratio: 76.3,
    member_count: 41820,
    processing_time: 22.7,
    compliance_score: 87.2
  }
])

// Enhanced methods
const formatCurrency = (amount: number): string => {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(amount)
}

const formatChange = (change: number): string => {
  const sign = change >= 0 ? '+' : ''
  return `${sign}${change.toFixed(1)}`
}

const getChangeClass = (change: number): string => {
  if (change > 0) return 'text-success'
  if (change < 0) return 'text-error'
  return 'text-medium-emphasis'
}

const getBenefitIcon = (benefit: string): string => {
  const icons: Record<string, string> = {
    'Group Life Assurance': 'mdi-shield-account',
    'Critical Illness': 'mdi-medical-bag',
    'Permanent Total Disability': 'mdi-wheelchair-accessibility',
    'Temporary Total Disability': 'mdi-account-injury',
    'Funeral Benefit': 'mdi-flower'
  }
  return icons[benefit] || 'mdi-shield'
}

const getBenefitIconColor = (benefit: string): string => {
  const colors: Record<string, string> = {
    'Group Life Assurance': 'primary',
    'Critical Illness': 'error',
    'Permanent Total Disability': 'warning',
    'Temporary Total Disability': 'info',
    'Funeral Benefit': 'purple'
  }
  return colors[benefit] || 'grey'
}

const getLossRatioColor = (ratio: number): string => {
  if (ratio <= 70) return 'success'
  if (ratio <= 85) return 'warning'
  return 'error'
}

const getProcessingTimeClass = (time: number): string => {
  if (time <= 12) return 'text-success font-weight-bold'
  if (time <= 24) return 'text-warning font-weight-bold'
  return 'text-error font-weight-bold'
}

const getSLAColor = (compliance: number): string => {
  if (compliance >= 90) return 'success'
  if (compliance >= 75) return 'warning'
  return 'error'
}

const getComplianceColor = (score: number): string => {
  if (score >= 95) return 'success'
  if (score >= 85) return 'warning'
  return 'error'
}

const getErrorSeverityColor = (severity: string): string => {
  const colors: Record<string, string> = {
    high: 'error',
    medium: 'warning',
    low: 'info'
  }
  return colors[severity] || 'grey'
}

const loadAnalytics = () => {
  // TODO: Load analytics data based on selected period
  console.log(
    'Loading SA Group Risk analytics for period:',
    selectedPeriod.value
  )
}

const exportReport = () => {
  // TODO: Export comprehensive SA Group Risk analytics report
  console.log(
    'Exporting SA Group Risk analytics report for period:',
    selectedPeriod.value
  )
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
