<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex justify-space-between align-center">
              <span class="headline">Bordereaux Management Dashboard</span>
            </div>
          </template>
          <template #default>
            <!-- Quick Stats Cards -->
            <v-row class="mb-6">
              <v-col v-if="error.stats" cols="12">
                <v-alert
                  type="error"
                  density="compact"
                  class="mb-4"
                  :text="error.stats"
                  dismissible
                  @click:close="error.stats = null"
                ></v-alert>
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <stat-card
                  title="Generated This Month"
                  :value="String(stats.thisMonth.generated)"
                  icon="mdi-file-document-multiple"
                  color="primary"
                  :loading="loading.stats"
                />
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <stat-card
                  title="Pending Submissions"
                  :value="String(stats.pending.submissions)"
                  icon="mdi-clock-outline"
                  color="warning"
                  :loading="loading.stats"
                />
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <stat-card
                  title="Reconciled This Week"
                  :value="String(stats.reconciled.count)"
                  icon="mdi-check-circle"
                  color="success"
                  :loading="loading.stats"
                />
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <stat-card
                  title="Active Templates"
                  :value="String(stats.templates.active)"
                  icon="mdi-file-cog"
                  color="info"
                  :loading="loading.stats"
                />
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <stat-card
                  title="Inbound Pending"
                  :value="String(stats.inbound.pending)"
                  icon="mdi-inbox-arrow-down"
                  color="secondary"
                  :loading="loading.stats"
                />
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <stat-card
                  title="Deadlines Overdue"
                  :value="String(stats.deadlines.overdue)"
                  icon="mdi-calendar-clock"
                  color="error"
                  :loading="loading.stats"
                />
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <stat-card
                  title="Pending Notifications"
                  :value="String(stats.notifications.pending)"
                  icon="mdi-bell-ring"
                  color="warning"
                  :loading="loading.stats"
                />
              </v-col>
            </v-row>

            <!-- Main Actions Grid -->
            <!-- Outbound Section -->
            <template v-if="hasPermission('bordereaux:generate_outbound')">
            <div
              class="text-subtitle-2 font-weight-bold mb-2 text-medium-emphasis text-uppercase"
              >Outbound</div
            >
            <v-row class="mb-2">
              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="navigateToGeneration"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="primary" class="mb-4"
                      >mdi-file-document-plus</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >Generate Bordereaux</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      Create member, premium, claims, or benefit bordereaux
                      using configurable templates
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="navigateToSubmissionTracking"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="info" class="mb-4"
                      >mdi-send-check</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >Submission Tracking</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      Track bordereaux submissions, monitor delivery status, and
                      manage scheme responses
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="navigateToReconciliation"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="success" class="mb-4"
                      >mdi-checkbox-multiple-marked</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >Reconciliation</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      Match submissions with scheme confirmations and resolve
                      discrepancies
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="navigateToTemplateManager"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="orange" class="mb-4"
                      >mdi-file-cog</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >Template Manager</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      Manage bordereaux templates for different insurers and
                      customize field mappings
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="navigateToAnalytics"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="purple" class="mb-4"
                      >mdi-chart-line</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >Analytics Dashboard</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      View bordereaux metrics, processing times, and compliance
                      reports
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="openComplianceDialog"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="teal" class="mb-4"
                      >mdi-shield-check</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >Compliance Center</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      FSP Act compliance, SARS reporting, and regulatory
                      documentation
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            </template>

            <v-divider class="my-4"></v-divider>
            <template v-if="hasPermission('bordereaux:submit_inbound')">
            <div
              class="text-subtitle-2 font-weight-bold mb-2 text-medium-emphasis text-uppercase"
              >Inbound</div
            >
            <v-row class="mb-2">
              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="navigateToInboundSubmissions"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="deep-purple" class="mb-4"
                      >mdi-inbox-arrow-down</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >Inbound Submissions</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      Review and process inbound employer member data
                      submissions before premium schedule generation
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="navigateToDeadlineCalendar"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="red" class="mb-4"
                      >mdi-calendar-clock</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >Deadline Calendar</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      Manage and monitor submission deadlines for all in-force
                      schemes; track overdue items and generate monthly
                      calendars
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            </template>

            <v-divider class="my-4"></v-divider>
            <template v-if="hasPermission('reinsurance:view')">
            <div
              class="text-subtitle-2 font-weight-bold mb-2 text-medium-emphasis text-uppercase"
              >Reinsurance</div
            >
            <v-row class="mb-6">
              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="navigateToReinsurerTracking"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="indigo" class="mb-4"
                      >mdi-handshake</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >Reinsurer Tracking</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      Track reinsurer acceptance responses and claim recovery
                      amounts for submitted bordereaux
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="navigateToClaimNotifications"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="amber-darken-2" class="mb-4"
                      >mdi-bell-ring</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >Claim Notifications</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      Track and manage the cadence of claim notifications to
                      reinsurers; generate month-end status updates for open
                      claims
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="navigateToRITreaties"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="deep-purple" class="mb-4"
                      >mdi-file-sign</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >RI Treaty Management</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      Define and manage reinsurance treaties (quota share,
                      surplus, XL), link schemes, and track treaty terms and
                      thresholds
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="navigateToRIBordereaux"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="teal-darken-1" class="mb-4"
                      >mdi-table-large</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >RI Bordereaux Generation</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      Generate member census and claims bordereaux runs per
                      treaty; submit and acknowledge RI bordereaux with full
                      audit trail
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="navigateToRIClaims"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="red-darken-2" class="mb-4"
                      >mdi-alert-circle</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >RI Large Claims SLA</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      Monitor large claim notification obligations to
                      reinsurers; track SLA deadlines, send status and run the
                      claims monitor
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="navigateToRISettlement"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="blue-grey-darken-1" class="mb-4"
                      >mdi-bank-transfer</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >Technical Accounts &amp; Settlement</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      Generate quarterly technical accounts, record RI
                      settlement payments, and track net balances owed between
                      cedant and reinsurer
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="navigateToRISubmissionRegister"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="cyan-darken-2" class="mb-4"
                      >mdi-format-list-checks</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >RI Submission Register</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      Full audit register of all RI bordereaux runs with BPR
                      references, receipt confirmations, version history and
                      amendment controls
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" lg="3">
                <v-card
                  variant="outlined"
                  class="h-100 action-card"
                  @click="navigateToRIKPIDashboard"
                >
                  <v-card-text class="text-center pa-6">
                    <v-icon size="60" color="green-darken-2" class="mb-4"
                      >mdi-chart-bar</v-icon
                    >
                    <h3 class="text-h6 font-weight-bold mb-2"
                      >RI KPI Dashboard</h3
                    >
                    <p class="text-body-2 text-medium-emphasis">
                      Live management KPI report: submission timeliness,
                      first-time acceptance, settlement performance and open
                      query backlog against §8.2 targets
                    </p>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            </template>

            <!-- Recent Activity -->
            <v-row>
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title
                    class="text-h6 font-weight-bold d-flex align-center justify-space-between"
                  >
                    <span>Recent Bordereaux Activity</span>
                    <v-btn
                      size="small"
                      variant="outlined"
                      prepend-icon="mdi-refresh"
                      :loading="loading.activity"
                      @click="fetchRecentActivity"
                    >
                      Refresh
                    </v-btn>
                  </v-card-title>
                  <v-card-text>
                    <v-alert
                      v-if="error.activity"
                      type="error"
                      density="compact"
                      class="mb-4"
                      :text="error.activity"
                    ></v-alert>
                    <data-grid
                      :rowData="recentActivity"
                      :columnDefs="activityColumnDefs"
                      :pagination="true"
                      :showExport="true"
                      :density="'compact'"
                      tableTitle=""
                      @row-clicked="handleRowClick"
                    />
                    <div v-if="loading.activity" class="text-center pa-4">
                      <v-progress-circular
                        indeterminate
                        color="primary"
                      ></v-progress-circular>
                      <p class="text-body-2 mt-2">Loading recent activity...</p>
                    </div>
                    <div
                      v-if="!loading.activity && recentActivity.length === 0"
                      class="text-center pa-4"
                    >
                      <p class="text-body-2 text-medium-emphasis"
                        >No recent bordereaux activity</p
                      >
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Compliance Dialog -->
    <v-dialog v-model="complianceDialog" max-width="600">
      <v-card>
        <v-card-title class="text-h6 font-weight-bold bg-teal text-white">
          South African Compliance Center
        </v-card-title>
        <v-card-text class="pt-4">
          <v-list>
            <v-list-item>
              <template #prepend>
                <v-icon color="success">mdi-check-circle</v-icon>
              </template>
              <v-list-item-title>FSP Act Compliance</v-list-item-title>
              <v-list-item-subtitle
                >All bordereaux meet FSP regulatory
                requirements</v-list-item-subtitle
              >
            </v-list-item>

            <v-list-item>
              <template #prepend>
                <v-icon color="success">mdi-check-circle</v-icon>
              </template>
              <v-list-item-title>SARS Reporting</v-list-item-title>
              <v-list-item-subtitle
                >Tax-compliant premium and benefit
                documentation</v-list-item-subtitle
              >
            </v-list-item>

            <v-list-item>
              <template #prepend>
                <v-icon color="warning">mdi-clock</v-icon>
              </template>
              <v-list-item-title>FSCA Submission</v-list-item-title>
              <v-list-item-subtitle
                >Next quarterly report due in 15 days</v-list-item-subtitle
              >
            </v-list-item>

            <v-list-item>
              <template #prepend>
                <v-icon color="info">mdi-information</v-icon>
              </template>
              <v-list-item-title>POPIA Compliance</v-list-item-title>
              <v-list-item-subtitle
                >Data protection measures implemented</v-list-item-subtitle
              >
            </v-list-item>
          </v-list>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="grey"
            variant="outlined"
            @click="complianceDialog = false"
          >
            Close
          </v-btn>
          <v-btn color="teal" variant="flat" @click="generateComplianceReport">
            Generate Report
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import StatCard from '@/renderer/components/StatCard.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useStatusBarStore } from '@/renderer/store/statusBar'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

const { hasPermission } = usePermissionCheck()
const statusBarStore = useStatusBarStore()

const router = useRouter()

// Reactive data
const complianceDialog = ref(false)
const loading = ref({
  stats: false,
  activity: false
})
const error: any = ref({
  stats: null,
  activity: null
})
const stats = ref({
  thisMonth: {
    generated: 0
  },
  pending: {
    submissions: 0
  },
  reconciled: {
    count: 0
  },
  templates: {
    active: 0
  },
  inbound: {
    pending: 0
  },
  deadlines: {
    overdue: 0
  },
  notifications: {
    pending: 0
  }
})

const activityColumnDefs = [
  {
    headerName: 'Bordereaux ID',
    field: 'generated_id',
    sortable: true,
    width: 150,
    filter: true
  },
  {
    headerName: 'Type',
    field: 'type',
    sortable: true,
    width: 120,
    filter: true,
    cellRenderer: (params: any) => {
      if (params.value) {
        const color = getBordereauTypeColor(params.value)
        const formattedType = formatBordereauType(params.value)
        return `<span style="padding:2px 10px;border-radius:12px;font-size:12px;background:${color};color:#fff;font-weight:500">${formattedType}</span>`
      }
      return ''
    }
  },
  {
    headerName: 'Scheme',
    field: 'scheme_name',
    sortable: true,
    width: 200,
    filter: true
  },
  {
    headerName: 'Status',
    field: 'status',
    sortable: true,
    width: 120,
    filter: true,
    cellRenderer: (params: any) => {
      if (params.value) {
        const color = getStatusColor(params.value)
        const formattedStatus = formatStatus(params.value)
        return `<span style="padding:2px 10px;border-radius:12px;font-size:12px;background:${color};color:#fff;font-weight:500">${formattedStatus}</span>`
      }
      return ''
    }
  },
  {
    headerName: 'Created',
    field: 'created_at',
    sortable: true,
    width: 160,
    filter: true,
    valueFormatter: (params: any) => {
      if (params.value) {
        return formatDate(params.value)
      }
      return ''
    }
  },
  {
    headerName: 'Actions',
    field: 'actions',
    sortable: false,
    width: 80,
    pinned: 'right',
    cellRenderer: (params: any) => {
      const key = String(params.data.id).replace(/-/g, '_')
      ;(window as any)[`showActivityMenu_${key}`] = (event: MouseEvent) =>
        showActivityContextMenu(event, params.data)
      return `<div style="display:flex;align-items:center;justify-content:center;height:100%">
        <button onclick="showActivityMenu_${key}(event)" title="Actions" style="background:none;border:none;cursor:pointer;padding:4px 8px;border-radius:4px;color:#616161;display:flex;align-items:center;justify-content:center">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="5" r="2"/><circle cx="12" cy="12" r="2"/><circle cx="12" cy="19" r="2"/></svg>
        </button>
      </div>`
    }
  }
]

const recentActivity: any = ref([])

// Data fetching functions
const fetchDashboardStats = async () => {
  loading.value.stats = true
  error.value.stats = null

  try {
    const response = await GroupPricingService.getBordereauxDashboardStats()
    console.log('Dashboard Stats Response:', response.data)
    if (response.data) {
      stats.value = {
        thisMonth: {
          generated: response.data.generated_this_month || 0
        },
        pending: {
          submissions: response.data.pending_submissions || 0
        },
        reconciled: {
          count: response.data.reconciled_this_week || 0
        },
        templates: {
          active: response.data.active_templates || 0
        },
        inbound: {
          pending: response.data.pending_inbound_submissions || 0
        },
        deadlines: {
          overdue: response.data.overdue_deadlines || 0
        },
        notifications: {
          pending: response.data.pending_notifications || 0
        }
      }
    }
  } catch (err) {
    error.value.stats = 'Failed to load dashboard statistics'
    console.error('Error fetching dashboard stats:', err)
    // Keep default values (0) if API fails
  } finally {
    loading.value.stats = false
  }
}

const fetchRecentActivity = async () => {
  loading.value.activity = true
  error.value.activity = null

  try {
    const response = await GroupPricingService.getBordereauxActivity({
      limit: 10
    })
    console.log('Recent Activity Response:', response)
    if (response.data) {
      recentActivity.value = response.data || []
    }
  } catch (err) {
    error.value.activity = 'Failed to load recent activity'
    console.error('Error fetching recent activity:', err)
  } finally {
    loading.value.activity = false
  }
}

const fetchNotificationStats = async () => {
  try {
    const res = await GroupPricingService.getNotificationStats({})
    if (res.data?.data) {
      stats.value.notifications.pending =
        (res.data.data.pending ?? 0) + (res.data.data.overdue ?? 0)
    }
  } catch {
    // non-fatal
  }
}

const fetchDeadlineStats = async () => {
  try {
    const res = await GroupPricingService.getDeadlineStats()
    if (res.data?.data) {
      stats.value.deadlines.overdue = res.data.data.overdue_count ?? 0
    }
  } catch {
    // non-fatal
  }
}

const loadDashboardData = async () => {
  await Promise.all([
    fetchDashboardStats(),
    fetchRecentActivity(),
    fetchDeadlineStats(),
    fetchNotificationStats()
  ])
  const items = [
    {
      icon: 'mdi-inbox-arrow-down',
      text: `Inbound pending: ${stats.value.inbound.pending}`,
      severity: stats.value.inbound.pending > 0 ? 'warn' : ('info' as any)
    },
    {
      icon: 'mdi-send-clock-outline',
      text: `Outbound pending: ${stats.value.pending.submissions}`,
      severity: stats.value.pending.submissions > 0 ? 'warn' : ('info' as any)
    }
  ]
  if (stats.value.deadlines.overdue > 0) {
    items.push({
      icon: 'mdi-calendar-alert',
      text: `Overdue deadlines: ${stats.value.deadlines.overdue}`,
      severity: 'error' as any
    })
  }
  statusBarStore.set(items)
}

// Navigation methods
const navigateToGeneration = () => {
  router.push('/group-pricing/bordereaux-management/generation')
}

const navigateToSubmissionTracking = () => {
  router.push('/group-pricing/bordereaux-management/tracking')
}

const navigateToReconciliation = () => {
  router.push('/group-pricing/bordereaux-management/reconciliation')
}

const navigateToTemplateManager = () => {
  router.push('/group-pricing/bordereaux-management/templates')
}

const navigateToAnalytics = () => {
  router.push('/group-pricing/bordereaux-management/analytics')
}

const navigateToInboundSubmissions = () => {
  router.push('/group-pricing/bordereaux-management/inbound-submissions')
}

const navigateToDeadlineCalendar = () => {
  router.push('/group-pricing/bordereaux-management/deadline-calendar')
}

const navigateToReinsurerTracking = () => {
  router.push('/group-pricing/bordereaux-management/reinsurer-tracking')
}

const navigateToClaimNotifications = () => {
  router.push('/group-pricing/bordereaux-management/claim-notifications')
}

const navigateToRITreaties = () => {
  router.push('/group-pricing/bordereaux-management/ri-treaties')
}

const navigateToRIBordereaux = () => {
  router.push('/group-pricing/bordereaux-management/ri-bordereaux')
}

const navigateToRIClaims = () => {
  router.push('/group-pricing/bordereaux-management/ri-claims')
}

const navigateToRISettlement = () => {
  router.push('/group-pricing/bordereaux-management/ri-settlement')
}

const navigateToRISubmissionRegister = () => {
  router.push('/group-pricing/bordereaux-management/ri-submission-register')
}

const navigateToRIKPIDashboard = () => {
  router.push('/group-pricing/bordereaux-management/ri-kpi-dashboard')
}

const openComplianceDialog = () => {
  complianceDialog.value = true
}

const generateComplianceReport = () => {
  // TODO: Implement compliance report generation
  console.log('Generating compliance report...')
  complianceDialog.value = false
}

// Utility methods
const getBordereauTypeColor = (type: string): string => {
  const colors: Record<string, string> = {
    member: '#1976d2',
    premium: '#388e3c',
    claims: '#f57c00',
    benefits: '#7b1fa2'
  }
  return colors[type] || '#757575'
}

const formatBordereauType = (type: string): string => {
  const types: Record<string, string> = {
    member: 'Member',
    premium: 'Premium',
    claims: 'Claims',
    benefits: 'Benefits'
  }
  return types[type] || type
}

const getStatusColor = (status: string): string => {
  const colors: Record<string, string> = {
    draft: '#757575',
    generated: '#1976d2',
    submitted: '#f57c00',
    confirmed: '#00796b',
    reconciled: '#388e3c',
    pending: '#fbc02d',
    rejected: '#d32f2f'
  }
  return colors[status] || '#757575'
}

const formatStatus = (status: string): string => {
  return status.charAt(0).toUpperCase() + status.slice(1)
}

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('en-ZA', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

let activeActivityMenuCleanup: (() => void) | null = null

function showActivityContextMenu(event: MouseEvent, data: any) {
  if (activeActivityMenuCleanup) activeActivityMenuCleanup()

  const menuItems = [
    { label: 'View', color: '#1976d2', fn: () => viewBordereaux(data) }
  ]

  const menu = document.createElement('div')
  menu.style.cssText =
    'position:fixed;background:#fff;border:1px solid #e0e0e0;border-radius:8px;' +
    'box-shadow:0 4px 16px rgba(0,0,0,0.14);z-index:9999;min-width:160px;padding:4px 0;'

  menuItems.forEach(({ label, color, fn }) => {
    const item = document.createElement('div')
    item.textContent = label
    item.style.cssText = `padding:8px 16px;cursor:pointer;font-size:13px;color:${color};`
    item.addEventListener(
      'mouseenter',
      () => (item.style.background = '#f5f5f5')
    )
    item.addEventListener('mouseleave', () => (item.style.background = ''))
    item.addEventListener('click', () => {
      cleanup()
      fn()
    })
    menu.appendChild(item)
  })

  document.body.appendChild(menu)

  const btn = (event.currentTarget || event.target) as HTMLElement
  const rect = btn.getBoundingClientRect()
  menu.style.top = `${rect.bottom + 4}px`
  menu.style.left = `${rect.left}px`

  const mr = menu.getBoundingClientRect()
  if (mr.right > window.innerWidth - 8)
    menu.style.left = `${rect.right - mr.width}px`
  if (mr.bottom > window.innerHeight - 8)
    menu.style.top = `${rect.top - mr.height - 4}px`

  function cleanup() {
    menu.remove()
    document.removeEventListener('click', outsideClick, true)
    activeActivityMenuCleanup = null
  }
  activeActivityMenuCleanup = cleanup

  function outsideClick(e: MouseEvent) {
    if (!menu.contains(e.target as Node) && e.target !== btn) cleanup()
  }
  setTimeout(() => document.addEventListener('click', outsideClick, true), 0)
}

const handleRowClick = (event: any) => {
  if (event.data && event.data.id) {
    viewBordereaux(event.data)
  }
}

const viewBordereaux = (item: any) => {
  let bordereauId: string

  // Handle both object with id property and direct ID string
  if (typeof item === 'string') {
    bordereauId = item
  } else if (item && item.id) {
    bordereauId = item.id
  } else {
    console.error('Invalid bordereaux item:', item)
    return
  }

  // Navigate to submission tracking page with the specific bordereaux ID
  router.push({
    path: '/group-pricing/bordereaux-management/tracking',
    query: { bordereaux_id: bordereauId }
  })
}

// Lifecycle
onMounted(() => {
  loadDashboardData()
})
onUnmounted(() => statusBarStore.clear())
</script>

<style scoped>
.action-card {
  cursor: pointer;
  transition: all 0.2s ease-in-out;
}

.action-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
</style>
