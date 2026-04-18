<template>
  <v-container fluid>
    <base-card :show-actions="false">
      <template #header>
        <h3 class="mb-0">Arrears Management</h3>
      </template>
      <template #default>
        <!-- Summary KPIs -->
        <v-row class="mb-4">
          <v-col cols="12" sm="4">
            <stat-card
              title="Total Outstanding"
              :value="fmtCurrency(totalOutstanding)"
              icon="mdi-cash-multiple"
              color="error"
              :loading="loading"
            />
          </v-col>
          <v-col cols="12" sm="4">
            <stat-card
              title="In Arrears"
              :value="String(inArrearsCount)"
              icon="mdi-clock-alert-outline"
              color="warning"
              :loading="loading"
            />
          </v-col>
          <v-col cols="12" sm="4">
            <stat-card
              title="Suspended"
              :value="String(suspendedCount)"
              icon="mdi-cancel"
              color="error"
              :loading="loading"
            />
          </v-col>
        </v-row>

        <!-- Filter Bar -->
        <v-row class="mb-3" align="center">
          <v-col cols="12" md="3">
            <v-select
              v-model="filters.status"
              label="Status"
              :items="statusOptions"
              variant="outlined"
              density="compact"
              clearable
            />
          </v-col>
          <v-col cols="12" md="3">
            <v-select
              v-model="agingBucket"
              label="Aging Bucket"
              :items="agingBuckets"
              variant="outlined"
              density="compact"
            />
          </v-col>
          <v-col cols="12" md="3" class="d-flex align-center">
            <v-btn
              variant="outlined"
              color="primary"
              prepend-icon="mdi-refresh"
              @click="loadArrears"
            >
              Refresh
            </v-btn>
          </v-col>
        </v-row>

        <!-- Aging Table -->
        <v-row class="mb-4">
          <v-col cols="12">
            <v-card variant="outlined">
              <v-card-title class="text-subtitle-1 pa-3"
                >Premium Arrears Aging</v-card-title
              >
              <v-card-text class="pa-0">
                <ag-grid-vue
                  class="ag-theme-balham"
                  :style="{ height: gridHeight, width: '100%' }"
                  :column-defs="columnDefs"
                  :row-data="filteredArrearsRecords"
                  :default-col-def="{
                    sortable: true,
                    resizable: true,
                    flex: 1
                  }"
                  :loading="loading"
                  :get-row-class="getRowClass"
                  @cell-clicked="onCellClicked"
                />
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </template>
    </base-card>

    <!-- Reminder Dialog -->
    <v-dialog v-model="reminderDialog" max-width="480" persistent>
      <v-card>
        <v-card-title>Send Payment Reminder</v-card-title>
        <v-card-text>
          <div class="text-body-2 mb-3">
            Scheme: <strong>{{ selectedScheme?.scheme_name }}</strong>
          </div>
          <v-textarea
            v-model="reminderMessage"
            label="Message (optional)"
            variant="outlined"
            density="compact"
            rows="3"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="reminderDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="actionLoading"
            @click="handleReminder"
            >Send Reminder</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Suspend Dialog -->
    <v-dialog v-model="suspendDialog" max-width="500" persistent>
      <v-card>
        <v-card-title class="text-error">Suspend Cover</v-card-title>
        <v-card-text>
          <v-alert type="error" variant="tonal" class="mb-4">
            <strong
              >This will cease cover for all members of
              {{ selectedScheme?.scheme_name }}.</strong
            >
          </v-alert>
          <v-row>
            <v-col cols="12">
              <v-text-field
                v-model="suspendForm.effective_date"
                label="Effective Date *"
                type="date"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12">
              <v-text-field
                v-model="suspendForm.reason"
                label="Reason *"
                variant="outlined"
                density="compact"
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="suspendDialog = false">Cancel</v-btn>
          <v-btn color="error" :loading="actionLoading" @click="handleSuspend"
            >Suspend Cover</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Reinstate Dialog -->
    <v-dialog v-model="reinstateDialog" max-width="500" persistent>
      <v-card>
        <v-card-title>Reinstate Cover</v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12">
              <v-text-field
                v-model="reinstateForm.reinstatement_date"
                label="Reinstatement Date *"
                type="date"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12">
              <v-text-field
                v-model.number="reinstateForm.back_premium"
                label="Back Premium Amount (0 = none)"
                type="number"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="reinstateForm.notes"
                label="Notes"
                variant="outlined"
                density="compact"
                rows="2"
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="reinstateDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="actionLoading"
            @click="handleReinstate"
            >Reinstate</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Payment Plan Dialog -->
    <v-dialog v-model="planDialog" max-width="600" persistent>
      <v-card>
        <v-card-title class="d-flex align-center justify-space-between">
          <span>Record Payment Plan — {{ selectedScheme?.scheme_name }}</span>
          <v-chip
            :color="
              planRemaining < 0
                ? 'error'
                : planRemaining === 0
                  ? 'success'
                  : 'warning'
            "
            variant="tonal"
            size="small"
          >
            Remaining: {{ fmtCurrency(planRemaining) }}
          </v-chip>
        </v-card-title>
        <v-card-text>
          <v-row v-for="(inst, i) in planInstalments" :key="i" align="center">
            <v-col cols="5">
              <v-text-field
                v-model="inst.date"
                :label="`Instalment ${i + 1} Date`"
                type="date"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="5">
              <v-text-field
                v-model.number="inst.amount"
                :label="`Amount ${i + 1}`"
                type="number"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="2">
              <v-btn
                icon="mdi-minus-circle"
                variant="plain"
                color="error"
                @click="planInstalments.splice(i, 1)"
              />
            </v-col>
          </v-row>
          <v-btn
            variant="text"
            prepend-icon="mdi-plus"
            @click="planInstalments.push({ date: '', amount: 0 })"
          >
            Add Instalment
          </v-btn>
          <v-textarea
            v-model="planNotes"
            label="Notes"
            variant="outlined"
            density="compact"
            rows="2"
            class="mt-3"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="planDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="actionLoading"
            :disabled="!planInstalments.length"
            @click="handlePlan"
          >
            Save Plan
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- History Drawer -->
    <v-navigation-drawer
      v-model="historyDrawer"
      location="right"
      width="420"
      temporary
    >
      <v-card flat>
        <v-card-title class="d-flex align-center justify-space-between pa-3">
          <span>Arrears History — {{ selectedScheme?.scheme_name }}</span>
          <v-btn
            icon="mdi-close"
            variant="plain"
            size="small"
            @click="historyDrawer = false"
          />
        </v-card-title>
        <v-divider />
        <v-card-text>
          <!-- Payment Plans section -->
          <div v-if="paymentPlans.length" class="mb-4">
            <div class="text-subtitle-2 mb-2">Payment Plans</div>
            <v-card
              v-for="plan in paymentPlans"
              :key="plan.id"
              variant="outlined"
              class="mb-2"
            >
              <v-card-title
                class="text-body-2 d-flex align-center justify-space-between pa-2"
              >
                <div>
                  <v-chip
                    :color="plan.status === 'active' ? 'primary' : 'grey'"
                    size="x-small"
                    class="mr-2"
                    >{{ plan.status }}</v-chip
                  >
                  <span class="text-caption"
                    >Plan #{{ plan.id }} ·
                    {{ plan.instalments?.length ?? 0 }} instalment{{
                      (plan.instalments?.length ?? 0) === 1 ? '' : 's'
                    }}
                    · Total {{ fmtCurrency(planSum(plan)) }}</span
                  >
                </div>
                <v-btn
                  size="x-small"
                  variant="text"
                  :icon="
                    expandedPlans[plan.id]
                      ? 'mdi-chevron-up'
                      : 'mdi-chevron-down'
                  "
                  @click="expandedPlans[plan.id] = !expandedPlans[plan.id]"
                />
              </v-card-title>
              <v-divider v-if="expandedPlans[plan.id]" />
              <v-card-text v-if="expandedPlans[plan.id]" class="pa-2">
                <div v-if="plan.notes" class="text-caption mb-2 text-grey">
                  {{ plan.notes }}
                </div>
                <v-table density="compact" class="text-caption">
                  <thead>
                    <tr>
                      <th class="text-left">Date</th>
                      <th class="text-right">Amount</th>
                      <th class="text-left">Status</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="inst in plan.instalments" :key="inst.id">
                      <td>{{ inst.date }}</td>
                      <td class="text-right">{{ fmtCurrency(inst.amount) }}</td>
                      <td>
                        <v-chip
                          :color="instalmentColor(inst.status)"
                          size="x-small"
                          >{{ inst.status }}</v-chip
                        >
                      </td>
                    </tr>
                  </tbody>
                </v-table>
              </v-card-text>
            </v-card>
          </div>

          <!-- Event timeline -->
          <div
            v-if="paymentPlans.length && historyItems.length"
            class="text-subtitle-2 mb-2"
            >Events</div
          >
          <v-timeline v-if="historyItems.length" density="compact" side="end">
            <v-timeline-item
              v-for="h in historyItems"
              :key="h.id"
              :dot-color="eventColor(h.event_type)"
              size="small"
            >
              <div class="text-caption text-grey"
                >{{ h.event_date }} · {{ h.performed_by }}</div
              >
              <div class="text-body-2 font-weight-medium">{{
                h.event_type
              }}</div>
              <div class="text-body-2">{{ h.description }}</div>
            </v-timeline-item>
          </v-timeline>
          <v-skeleton-loader
            v-else-if="historyLoading"
            type="list-item-two-line@5"
          />
          <div
            v-else-if="!paymentPlans.length"
            class="text-body-2 text-grey text-center mt-4"
            >No history found.</div
          >
        </v-card-text>
      </v-card>
    </v-navigation-drawer>

    <v-snackbar
      v-model="snackbar"
      :color="snackbarColor"
      :timeout="3500"
      centered
    >
      {{ snackbarText }}
      <template #actions
        ><v-btn variant="text" @click="snackbar = false">Close</v-btn></template
      >
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-balham.css'
import { AgGridVue } from 'ag-grid-vue3'
import PremiumManagementService from '@/renderer/api/PremiumManagementService'
import BaseCard from '@/renderer/components/BaseCard.vue'
import StatCard from '@/renderer/components/StatCard.vue'
import { useGridHeight } from '@/renderer/composables/useGridHeight'
import { statusCellRenderer } from '@/renderer/utils/statusCellRenderer'
import { fmtCurrency } from '@/renderer/utils/formatters'
import { useStatusBarStore } from '@/renderer/store/statusBar'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

const statusBarStore = useStatusBarStore()
const { hasPermission } = usePermissionCheck()

const gridHeight = useGridHeight(340)
const agingBucket = ref('All')
const agingBuckets = ['All', '30+', '60+', '90+', '120+']

const loading = ref(false)
const actionLoading = ref(false)
const historyLoading = ref(false)
const reminderDialog = ref(false)
const suspendDialog = ref(false)
const reinstateDialog = ref(false)
const planDialog = ref(false)
const historyDrawer = ref(false)
const snackbar = ref(false)
const snackbarText = ref('')
const snackbarColor = ref('success')

const arrearsRecords = ref<any[]>([])
const selectedScheme = ref<any>(null)
const historyItems = ref<any[]>([])
const paymentPlans = ref<any[]>([])
const expandedPlans = ref<Record<number, boolean>>({})
const reminderMessage = ref('')
const planInstalments = ref<Array<{ date: string; amount: number }>>([])
const planNotes = ref('')

const today = new Date().toISOString().slice(0, 10)
const suspendForm = ref({ effective_date: today, reason: '' })
const reinstateForm = ref({
  reinstatement_date: today,
  back_premium: 0,
  notes: ''
})

const filters = ref({ status: null as string | null })
const statusOptions = ['current', 'in_arrears', 'suspended']

const filteredArrearsRecords = computed(() => {
  if (agingBucket.value === 'All') return arrearsRecords.value
  const thresholds: Record<string, number> = {
    '30+': 30,
    '60+': 60,
    '90+': 90,
    '120+': 120
  }
  const days = thresholds[agingBucket.value] ?? 0
  return arrearsRecords.value.filter((r: any) => (r.days_overdue ?? 0) >= days)
})

const totalOutstanding = computed(() =>
  arrearsRecords.value.reduce((s: number, r: any) => s + r.total_outstanding, 0)
)
const planTotal = computed(() =>
  planInstalments.value.reduce((s, inst) => s + (inst.amount || 0), 0)
)
const planRemaining = computed(
  () => (selectedScheme.value?.total_outstanding ?? 0) - planTotal.value
)
const inArrearsCount = computed(
  () =>
    arrearsRecords.value.filter((r: any) => r.status === 'in_arrears').length
)
const suspendedCount = computed(
  () => arrearsRecords.value.filter((r: any) => r.status === 'suspended').length
)

const columnDefs = [
  { headerName: 'Scheme', field: 'scheme_name', minWidth: 180 },
  {
    headerName: '0-30 days',
    field: 'days_0_to_30',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: '31-60 days',
    field: 'days_31_to_60',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: '61-90 days',
    field: 'days_61_to_90',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: '90+ days',
    field: 'days_over_90',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: 'Total',
    field: 'total_outstanding',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: 'Status',
    field: 'status',
    cellRenderer: (p: any) => statusCellRenderer(p.value)
  },
  ...(hasPermission('premiums:manage_arrears')
    ? [
        {
          headerName: 'Actions',
          colId: 'actions',
          cellRenderer: (params: any) => {
            const s = params.data?.status
            const btn = (action: string, label: string, color = '#333') =>
              `<button data-action="${action}" style="margin:2px;padding:2px 6px;font-size:11px;cursor:pointer;color:${color}">${label}</button>`
            const actions = [btn('remind', 'Remind'), btn('plan', 'Plan')]
            if (s !== 'suspended') {
              actions.push(btn('suspend', 'Suspend', '#F44336'))
            } else {
              actions.push(btn('reinstate', 'Reinstate', '#4CAF50'))
            }
            actions.push(btn('history', 'History'))
            return actions.join('')
          },
          minWidth: 260
        }
      ]
    : [])
]

function onCellClicked(e: any) {
  if (e.column?.getColId() !== 'actions') return
  const action = (e.event?.target as HTMLElement)?.dataset?.action
  if (!action) return
  selectedScheme.value = e.data
  if (action === 'remind') {
    reminderDialog.value = true
  } else if (action === 'suspend') {
    suspendDialog.value = true
  } else if (action === 'reinstate') {
    reinstateDialog.value = true
  } else if (action === 'plan') {
    planInstalments.value = [{ date: '', amount: 0 }]
    planNotes.value = ''
    planDialog.value = true
  } else if (action === 'history') {
    historyDrawer.value = true
    historyLoading.value = true
    historyItems.value = []
    paymentPlans.value = []
    expandedPlans.value = {}
    Promise.all([
      PremiumManagementService.getArrearsHistory(e.data.scheme_id),
      PremiumManagementService.getPaymentPlans(e.data.scheme_id)
    ])
      .then(([histRes, plansRes]: any[]) => {
        historyItems.value = histRes.data.data ?? []
        paymentPlans.value = plansRes.data.data ?? []
      })
      .finally(() => {
        historyLoading.value = false
      })
  }
}

function getRowClass(params: any) {
  if (params.data?.status === 'suspended') return 'ag-row-error'
  if (params.data?.status === 'in_arrears') return 'ag-row-warning'
  return ''
}

function eventColor(type: string) {
  const map: Record<string, string> = {
    REMINDER: 'info',
    PAYMENT_PLAN: 'primary',
    SUSPENDED: 'error',
    REINSTATED: 'success'
  }
  return map[type] ?? 'grey'
}

function planSum(plan: any) {
  return (plan.instalments ?? []).reduce(
    (sum: number, i: any) => sum + (Number(i.amount) || 0),
    0
  )
}

function instalmentColor(status: string) {
  const map: Record<string, string> = {
    pending: 'warning',
    paid: 'success',
    overdue: 'error'
  }
  return map[status] ?? 'grey'
}

async function loadArrears() {
  loading.value = true
  try {
    const params: any = {}
    if (filters.value.status) params.status = filters.value.status
    const res = await PremiumManagementService.getArrearsAging(params)
    arrearsRecords.value = res.data.data ?? []
    const outstanding = arrearsRecords.value.reduce(
      (s: number, r: any) => s + r.total_outstanding,
      0
    )
    const inArrears = arrearsRecords.value.filter(
      (r: any) => r.status === 'in_arrears'
    ).length
    const suspended = arrearsRecords.value.filter(
      (r: any) => r.status === 'suspended'
    ).length
    statusBarStore.set([
      {
        icon: 'mdi-cash-multiple',
        text: `Outstanding: ${fmtCurrency(outstanding)}`,
        severity: outstanding > 0 ? 'warn' : 'info'
      },
      {
        icon: 'mdi-clock-alert-outline',
        text: `In arrears: ${inArrears}`,
        severity: inArrears > 0 ? 'warn' : 'info'
      },
      {
        icon: 'mdi-cancel',
        text: `Suspended: ${suspended}`,
        severity: suspended > 0 ? 'error' : 'info'
      }
    ])
  } catch (e) {
    console.error('Failed to load arrears', e)
  } finally {
    loading.value = false
  }
}

async function handleReminder() {
  if (!selectedScheme.value) return
  actionLoading.value = true
  try {
    await PremiumManagementService.sendReminder(
      selectedScheme.value.scheme_id,
      {
        message: reminderMessage.value
      }
    )
    reminderDialog.value = false
    showSnack('Reminder logged')
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Failed', 'error')
  } finally {
    actionLoading.value = false
  }
}

async function handleSuspend() {
  if (!selectedScheme.value) return
  actionLoading.value = true
  try {
    await PremiumManagementService.suspendCover(
      selectedScheme.value.scheme_id,
      suspendForm.value
    )
    suspendDialog.value = false
    showSnack('Cover suspended')
    await loadArrears()
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Failed', 'error')
  } finally {
    actionLoading.value = false
  }
}

async function handleReinstate() {
  if (!selectedScheme.value) return
  actionLoading.value = true
  try {
    await PremiumManagementService.reinstateCover(
      selectedScheme.value.scheme_id,
      reinstateForm.value
    )
    reinstateDialog.value = false
    showSnack('Cover reinstated')
    await loadArrears()
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Failed', 'error')
  } finally {
    actionLoading.value = false
  }
}

async function handlePlan() {
  if (!selectedScheme.value) return
  actionLoading.value = true
  try {
    await PremiumManagementService.recordPaymentPlan(
      selectedScheme.value.scheme_id,
      {
        instalments: planInstalments.value,
        notes: planNotes.value
      }
    )
    planDialog.value = false
    showSnack('Payment plan recorded')
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Failed', 'error')
  } finally {
    actionLoading.value = false
  }
}

function showSnack(msg: string, color = 'success') {
  snackbarText.value = msg
  snackbarColor.value = color
  snackbar.value = true
}

onMounted(() => {
  loadArrears()
})
onUnmounted(() => statusBarStore.clear())
</script>
