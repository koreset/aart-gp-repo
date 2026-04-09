<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center justify-space-between">
              <span class="headline"
                >RI Technical Accounts &amp; Settlement</span
              >
              <v-btn
                size="small"
                rounded
                color="white"
                variant="outlined"
                prepend-icon="mdi-arrow-left"
                @click="$router.push('/group-pricing/bordereaux-management')"
              >
                Back
              </v-btn>
            </div>
          </template>

          <template #default>
            <!-- KPI Cards -->
            <v-row class="mb-4">
              <v-col cols="6" sm="2">
                <v-card variant="outlined">
                  <v-card-text class="text-center pa-3">
                    <div class="text-h6 font-weight-bold text-primary">{{
                      stats.total
                    }}</div>
                    <div class="text-caption">Total</div>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="6" sm="2">
                <v-card variant="outlined">
                  <v-card-text class="text-center pa-3">
                    <div class="text-h6 font-weight-bold text-grey">{{
                      stats.draft
                    }}</div>
                    <div class="text-caption">Draft</div>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="6" sm="2">
                <v-card variant="outlined">
                  <v-card-text class="text-center pa-3">
                    <div class="text-h6 font-weight-bold text-warning">{{
                      stats.agreed
                    }}</div>
                    <div class="text-caption">Agreed</div>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="6" sm="2">
                <v-card variant="outlined">
                  <v-card-text class="text-center pa-3">
                    <div class="text-h6 font-weight-bold text-success">{{
                      stats.settled
                    }}</div>
                    <div class="text-caption">Settled</div>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="6" sm="2">
                <v-card variant="outlined">
                  <v-card-text class="text-center pa-3">
                    <div class="text-caption text-error font-weight-bold">{{
                      formatCurrency(stats.net_owed_to_ri)
                    }}</div>
                    <div class="text-caption">Owed to RI</div>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="6" sm="2">
                <v-card variant="outlined">
                  <v-card-text class="text-center pa-3">
                    <div class="text-caption text-success font-weight-bold">{{
                      formatCurrency(stats.net_owed_by_cedant)
                    }}</div>
                    <div class="text-caption">Due from RI</div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Filters -->
            <v-card variant="outlined" class="mb-4">
              <v-card-text>
                <v-row>
                  <v-col cols="12" sm="4">
                    <v-select
                      v-model="filters.treaty_id"
                      :items="activeTreaties"
                      item-title="treaty_name"
                      item-value="id"
                      label="Treaty"
                      variant="outlined"
                      density="compact"
                      clearable
                      @update:model-value="loadAccounts"
                    />
                  </v-col>
                  <v-col cols="12" sm="3">
                    <v-select
                      v-model="filters.status"
                      :items="statusOptions"
                      label="Status"
                      variant="outlined"
                      density="compact"
                      clearable
                      @update:model-value="loadAccounts"
                    />
                  </v-col>
                  <v-col cols="auto" class="d-flex align-center gap-2">
                    <v-btn
                      size="small"
                      rounded
                      color="success"
                      prepend-icon="mdi-refresh"
                      :loading="loading"
                      @click="loadAccounts"
                    >
                      Refresh
                    </v-btn>
                    <v-btn
                      size="small"
                      rounded
                      color="teal"
                      prepend-icon="mdi-plus"
                      @click="showGenerateDialog = true"
                    >
                      Generate Account
                    </v-btn>
                  </v-col>
                </v-row>
              </v-card-text>
            </v-card>

            <!-- Grid -->
            <div style="height: 460px">
              <data-grid
                :row-data="accounts"
                :column-defs="columnDefs"
                :default-col-def="{
                  resizable: true,
                  sortable: true,
                  filter: true
                }"
                :pagination="true"
                :pagination-page-size="20"
              />
            </div>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Generate Account Dialog -->
    <v-dialog v-model="showGenerateDialog" max-width="520">
      <v-card>
        <v-card-title>Generate Technical Account</v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12">
              <v-select
                v-model="genForm.treaty_id"
                :items="activeTreaties"
                item-title="treaty_name"
                item-value="id"
                label="Treaty *"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model="genForm.period_start"
                label="Period Start *"
                type="date"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model="genForm.period_end"
                label="Period End *"
                type="date"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="genForm.notes"
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
          <v-btn variant="text" @click="showGenerateDialog = false"
            >Cancel</v-btn
          >
          <v-btn color="teal" :loading="generating" @click="generateAccount"
            >Generate</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Payments Dialog -->
    <v-dialog v-model="showPaymentsDialog" max-width="700" scrollable>
      <v-card>
        <v-card-title
          >Payments — {{ selectedAccount?.account_number }}</v-card-title
        >
        <v-card-subtitle v-if="selectedAccount">
          Net Balance:
          <span
            :style="{
              color: selectedAccount.net_balance > 0 ? '#d32f2f' : '#388e3c',
              fontWeight: 'bold'
            }"
          >
            {{ formatNetBalance(selectedAccount) }}
          </span>
        </v-card-subtitle>
        <v-card-text>
          <v-table density="compact" class="mb-4">
            <thead>
              <tr>
                <th>Date</th>
                <th>Amount</th>
                <th>Direction</th>
                <th>Method</th>
                <th>Reference</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in payments" :key="p.id">
                <td>{{ p.payment_date }}</td>
                <td>{{ formatCurrency(p.amount) }}</td>
                <td>{{
                  p.direction === 'cedant_to_ri' ? 'Cedant → RI' : 'RI → Cedant'
                }}</td>
                <td>{{ p.payment_method || '—' }}</td>
                <td>{{ p.reference || '—' }}</td>
              </tr>
              <tr v-if="payments.length === 0">
                <td colspan="5" class="text-center text-medium-emphasis"
                  >No payments recorded</td
                >
              </tr>
            </tbody>
          </v-table>

          <v-divider class="mb-3" />
          <p class="text-subtitle-2 mb-2">Record Payment</p>
          <v-row>
            <v-col cols="6">
              <v-text-field
                v-model="payForm.payment_date"
                label="Payment Date *"
                type="date"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model.number="payForm.amount"
                label="Amount (ZAR) *"
                type="number"
                variant="outlined"
                density="compact"
                prefix="R"
              />
            </v-col>
            <v-col cols="6">
              <v-select
                v-model="payForm.direction"
                :items="directionOptions"
                label="Direction *"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model="payForm.payment_method"
                label="Payment Method"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="8">
              <v-text-field
                v-model="payForm.reference"
                label="Reference"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="4" class="d-flex align-center">
              <v-btn
                color="teal"
                :loading="recordingPayment"
                @click="recordPayment"
              >
                Record
              </v-btn>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showPaymentsDialog = false"
            >Close</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Escalate Dispute Dialog -->
    <v-dialog v-model="showEscalateDialog" max-width="480">
      <v-card>
        <v-card-title
          >Escalate Dispute —
          {{ escalatingAccount?.account_number }}</v-card-title
        >
        <v-card-subtitle class="text-medium-emphasis">
          Current dispute stage:
          {{ escalatingAccount?.dispute_status || 'none' }}
        </v-card-subtitle>
        <v-card-text>
          <v-row>
            <v-col cols="12">
              <v-select
                v-model="escalateForm.stage"
                :items="disputeStageOptions"
                label="Escalate to Stage *"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="escalateForm.notes"
                label="Escalation Notes"
                variant="outlined"
                density="compact"
                rows="3"
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showEscalateDialog = false"
            >Cancel</v-btn
          >
          <v-btn color="error" :loading="escalating" @click="performEscalate"
            >Escalate</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Resolve Dispute Dialog -->
    <v-dialog v-model="showResolveDialog" max-width="440">
      <v-card>
        <v-card-title
          >Resolve Dispute —
          {{ resolvingAccount?.account_number }}</v-card-title
        >
        <v-card-text>
          <p class="text-body-2 text-medium-emphasis mb-3">
            Resolving the dispute will return the account to
            <strong>Agreed</strong> status and mark the dispute stage as
            resolved.
          </p>
          <v-textarea
            v-model="resolveForm.notes"
            label="Resolution Notes"
            variant="outlined"
            density="compact"
            rows="3"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showResolveDialog = false"
            >Cancel</v-btn
          >
          <v-btn color="success" :loading="resolving" @click="performResolve"
            >Resolve</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3500">
      {{ snackbar.message }}
    </v-snackbar>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import BaseCard from '@/renderer/components/BaseCard.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'

const accounts = ref([])
const activeTreaties = ref([])
const payments = ref([])
const loading = ref(false)
const generating = ref(false)
const recordingPayment = ref(false)
const escalating = ref(false)
const resolving = ref(false)
const stats = ref({
  total: 0,
  draft: 0,
  issued: 0,
  agreed: 0,
  settled: 0,
  disputed: 0,
  net_owed_to_ri: 0,
  net_owed_by_cedant: 0
})
const filters = ref({ treaty_id: null, status: '' })
const showGenerateDialog = ref(false)
const showPaymentsDialog = ref(false)
const showEscalateDialog = ref(false)
const showResolveDialog = ref(false)
const selectedAccount = ref(null)
const escalatingAccount = ref(null)
const resolvingAccount = ref(null)
const escalateForm = ref({ stage: 'stage1', notes: '' })
const resolveForm = ref({ notes: '' })
const snackbar = ref({ show: false, message: '', color: 'success' })

const genForm = ref({
  treaty_id: null,
  period_start: '',
  period_end: '',
  notes: ''
})
const payForm = ref({
  payment_date: '',
  amount: 0,
  direction: 'cedant_to_ri',
  payment_method: '',
  reference: ''
})

const statusOptions = [
  { title: 'Draft', value: 'draft' },
  { title: 'Issued', value: 'issued' },
  { title: 'Agreed', value: 'agreed' },
  { title: 'Settled', value: 'settled' },
  { title: 'Disputed', value: 'disputed' }
]
const directionOptions = [
  { title: 'Cedant → Reinsurer', value: 'cedant_to_ri' },
  { title: 'Reinsurer → Cedant', value: 'ri_to_cedant' }
]
const disputeStageOptions = [
  { title: 'Stage 1 — Internal Review', value: 'stage1' },
  { title: 'Stage 2 — Formal Dispute', value: 'stage2' },
  { title: 'Stage 3 — Arbitration', value: 'stage3' }
]

const statusColor = {
  draft: '#9e9e9e',
  issued: '#1976d2',
  agreed: '#f57c00',
  settled: '#388e3c',
  disputed: '#d32f2f'
}

const formatCurrency = (v) =>
  v != null
    ? 'R ' +
      Number(v).toLocaleString('en-ZA', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      })
    : '—'

const formatNetBalance = (account) => {
  if (!account) return '—'
  const b = account.net_balance
  if (b > 0)
    return `R ${Number(b).toLocaleString('en-ZA', { minimumFractionDigits: 2 })} owed to RI`
  return `R ${Number(Math.abs(b)).toLocaleString('en-ZA', { minimumFractionDigits: 2 })} due from RI`
}

const columnDefs = [
  { field: 'account_number', headerName: 'Account No.', width: 190 },
  { field: 'treaty_number', headerName: 'Treaty', width: 130 },
  { field: 'reinsurer_name', headerName: 'Reinsurer', width: 140 },
  { field: 'period_start', headerName: 'Period Start', width: 115 },
  { field: 'period_end', headerName: 'Period End', width: 115 },
  {
    field: 'ceded_premium_earned',
    headerName: 'Ceded Prem',
    width: 130,
    valueFormatter: (p) => formatCurrency(p.value)
  },
  {
    field: 'reinsurance_commission',
    headerName: 'RI Commission',
    width: 130,
    valueFormatter: (p) => formatCurrency(p.value)
  },
  {
    field: 'total_ceded_claims',
    headerName: 'Ceded Claims',
    width: 130,
    valueFormatter: (p) => formatCurrency(p.value)
  },
  {
    field: 'net_balance',
    headerName: 'Net Balance',
    width: 150,
    cellRenderer: (p) => {
      const color = p.value > 0 ? '#d32f2f' : '#388e3c'
      const label =
        p.value > 0
          ? `R ${Number(p.value).toLocaleString('en-ZA', { minimumFractionDigits: 2 })} (owed)`
          : `R ${Number(Math.abs(p.value)).toLocaleString('en-ZA', { minimumFractionDigits: 2 })} (due)`
      return `<span style="color:${color};font-weight:600;font-size:12px">${label}</span>`
    }
  },
  {
    field: 'status',
    headerName: 'Status',
    width: 110,
    cellRenderer: (p) => {
      const color = statusColor[p.value] || '#9e9e9e'
      return `<span style="padding:2px 10px;border-radius:12px;font-size:12px;background:${color}22;color:${color};font-weight:500">${p.value}</span>`
    }
  },
  {
    field: 'dispute_status',
    headerName: 'Dispute Stage',
    width: 130,
    cellRenderer: (p) => {
      if (!p.value || p.value === 'none')
        return '<span style="color:#9e9e9e;font-size:12px">—</span>'
      const color = p.value === 'resolved' ? '#388e3c' : '#d32f2f'
      return `<span style="padding:2px 8px;border-radius:12px;font-size:11px;background:${color}22;color:${color};font-weight:500">${p.value}</span>`
    }
  },
  {
    headerName: 'Actions',
    width: 80,
    pinned: 'right',
    sortable: false,
    filter: false,
    cellRenderer: (p) => {
      const key = String(p.data.id).replace(/-/g, '_')
      window[`showAccountMenu_${key}`] = (event) =>
        showContextMenu(event, p.data)
      return `<div style="display:flex;align-items:center;justify-content:center;height:100%">
        <button onclick="showAccountMenu_${key}(event)" title="Actions" style="background:none;border:none;cursor:pointer;padding:4px 8px;border-radius:4px;color:#616161;display:flex;align-items:center;justify-content:center">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="5" r="2"/><circle cx="12" cy="12" r="2"/><circle cx="12" cy="19" r="2"/></svg>
        </button>
      </div>`
    }
  }
]

const notify = (message, color = 'success') => {
  snackbar.value = { show: true, message, color }
}

let activeMenuCleanup = null

function showContextMenu(event, data) {
  if (activeMenuCleanup) activeMenuCleanup()

  const canIssue = data.status === 'draft'
  const canAgree = data.status === 'issued'
  const isDisputed = data.status === 'disputed'
  const canResolve = isDisputed && data.dispute_status !== 'resolved'

  const menuItems = [
    { label: 'Payments', color: '#616161', fn: () => openPaymentsDialog(data) },
    canIssue
      ? { label: 'Issue', color: '#1976d2', fn: () => advance(data, 'issued') }
      : null,
    canAgree
      ? { label: 'Agree', color: '#f57c00', fn: () => advance(data, 'agreed') }
      : null,
    data.status === 'agreed'
      ? {
          label: 'Settle',
          color: '#388e3c',
          fn: () => openPaymentsDialog(data)
        }
      : null,
    isDisputed
      ? {
          label: 'Escalate',
          color: '#d32f2f',
          fn: () => openEscalateDialog(data)
        }
      : null,
    canResolve
      ? {
          label: 'Resolve Dispute',
          color: '#388e3c',
          fn: () => openResolveDialog(data)
        }
      : null
  ].filter(Boolean)

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

  const btn = event.currentTarget || event.target
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
    activeMenuCleanup = null
  }
  activeMenuCleanup = cleanup

  function outsideClick(e) {
    if (!menu.contains(e.target) && e.target !== btn) cleanup()
  }
  setTimeout(() => document.addEventListener('click', outsideClick, true), 0)
}

async function loadAccounts() {
  loading.value = true
  try {
    const params = {}
    if (filters.value.treaty_id) params.treaty_id = filters.value.treaty_id
    if (filters.value.status) params.status = filters.value.status
    const res = await GroupPricingService.getTechnicalAccounts(params)
    accounts.value = res.data?.data || []
  } catch {
    notify('Failed to load accounts', 'error')
  } finally {
    loading.value = false
  }
}

async function loadStats() {
  try {
    const params = {}
    if (filters.value.treaty_id) params.treaty_id = filters.value.treaty_id
    const res = await GroupPricingService.getSettlementStats(params)
    stats.value = res.data?.data || stats.value
  } catch {}
}

async function loadActiveTreaties() {
  try {
    const res = await GroupPricingService.getTreaties({ status: 'active' })
    activeTreaties.value = res.data?.data || []
  } catch {}
}

async function generateAccount() {
  if (
    !genForm.value.treaty_id ||
    !genForm.value.period_start ||
    !genForm.value.period_end
  ) {
    notify('Please fill required fields', 'warning')
    return
  }
  generating.value = true
  try {
    await GroupPricingService.generateTechnicalAccount(genForm.value)
    notify('Technical account generated')
    showGenerateDialog.value = false
    loadAccounts()
    loadStats()
  } catch (e) {
    notify(e.response?.data?.error || 'Failed to generate account', 'error')
  } finally {
    generating.value = false
  }
}

async function advance(account, status) {
  try {
    await GroupPricingService.updateTechnicalAccount(account.id, { status })
    notify(`Account marked as ${status}`)
    loadAccounts()
    loadStats()
  } catch {
    notify('Failed to update account', 'error')
  }
}

async function openPaymentsDialog(account) {
  selectedAccount.value = account
  payForm.value = {
    payment_date: '',
    amount: 0,
    direction: 'cedant_to_ri',
    payment_method: '',
    reference: ''
  }
  showPaymentsDialog.value = true
  try {
    const res = await GroupPricingService.getSettlementPayments({
      account_id: account.id
    })
    payments.value = res.data?.data || []
  } catch {
    notify('Failed to load payments', 'error')
  }
}

async function recordPayment() {
  if (
    !payForm.value.payment_date ||
    !payForm.value.amount ||
    !payForm.value.direction
  ) {
    notify('Please fill required fields', 'warning')
    return
  }
  recordingPayment.value = true
  try {
    await GroupPricingService.recordSettlementPayment({
      ...payForm.value,
      technical_account_id: selectedAccount.value.id
    })
    notify('Payment recorded')
    payForm.value = {
      payment_date: '',
      amount: 0,
      direction: 'cedant_to_ri',
      payment_method: '',
      reference: ''
    }
    const res = await GroupPricingService.getSettlementPayments({
      account_id: selectedAccount.value.id
    })
    payments.value = res.data?.data || []
    loadAccounts()
    loadStats()
  } catch {
    notify('Failed to record payment', 'error')
  } finally {
    recordingPayment.value = false
  }
}

function openEscalateDialog(account) {
  escalatingAccount.value = account
  escalateForm.value = { stage: 'stage1', notes: '' }
  showEscalateDialog.value = true
}

async function performEscalate() {
  if (!escalateForm.value.stage) {
    notify('Please select a stage', 'warning')
    return
  }
  escalating.value = true
  try {
    await GroupPricingService.escalateSettlementDispute(
      escalatingAccount.value.id,
      escalateForm.value
    )
    notify(`Dispute escalated to ${escalateForm.value.stage}`)
    showEscalateDialog.value = false
    loadAccounts()
    loadStats()
  } catch (e) {
    notify(e.response?.data?.error || 'Failed to escalate dispute', 'error')
  } finally {
    escalating.value = false
  }
}

function openResolveDialog(account) {
  resolvingAccount.value = account
  resolveForm.value = { notes: '' }
  showResolveDialog.value = true
}

async function performResolve() {
  resolving.value = true
  try {
    await GroupPricingService.resolveSettlementDispute(
      resolvingAccount.value.id,
      resolveForm.value
    )
    notify('Dispute resolved — account returned to Agreed')
    showResolveDialog.value = false
    loadAccounts()
    loadStats()
  } catch (e) {
    notify(e.response?.data?.error || 'Failed to resolve dispute', 'error')
  } finally {
    resolving.value = false
  }
}

onMounted(() => {
  loadAccounts()
  loadStats()
  loadActiveTreaties()
})
</script>
