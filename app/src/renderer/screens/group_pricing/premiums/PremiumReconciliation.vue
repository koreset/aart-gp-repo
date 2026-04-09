<template>
  <v-container fluid>
    <base-card :show-actions="false">
      <template #header>
        <h3 class="mb-0">Premium Reconciliation</h3>
      </template>
      <template #default>
        <!-- Summary Dashboard -->
        <v-row class="mb-4">
          <v-col cols="12" sm="6" md="3">
            <v-card variant="tonal" color="warning">
              <v-card-text class="d-flex align-center justify-space-between">
                <div>
                  <div class="text-caption">Unallocated Payments</div>
                  <div class="text-h6 font-weight-bold">{{
                    summary.unallocated_payment_count
                  }}</div>
                  <div class="text-caption">{{
                    fmtCurrency(summary.total_unallocated_payments)
                  }}</div>
                </div>
                <v-icon size="36">mdi-bank-outline</v-icon>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" sm="6" md="3">
            <v-card variant="tonal" color="error">
              <v-card-text class="d-flex align-center justify-space-between">
                <div>
                  <div class="text-caption">Unpaid Invoices</div>
                  <div class="text-h6 font-weight-bold">{{
                    summary.unpaid_invoice_count
                  }}</div>
                  <div class="text-caption">{{
                    fmtCurrency(summary.total_unpaid_invoices)
                  }}</div>
                </div>
                <v-icon size="36">mdi-file-document-alert-outline</v-icon>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" sm="6" md="3">
            <v-card variant="tonal" color="info">
              <v-card-text class="d-flex align-center justify-space-between">
                <div>
                  <div class="text-caption">Suspense Balance</div>
                  <div class="text-h6 font-weight-bold">{{
                    summary.suspense_count
                  }}</div>
                  <div class="text-caption">{{
                    fmtCurrency(summary.suspense_balance)
                  }}</div>
                </div>
                <v-icon size="36">mdi-clock-alert-outline</v-icon>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" sm="6" md="3">
            <v-card variant="tonal" color="deep-orange">
              <v-card-text class="d-flex align-center justify-space-between">
                <div>
                  <div class="text-caption">Aged &gt; 90 Days</div>
                  <div class="text-h6 font-weight-bold">{{
                    fmtCurrency(summary.aged_over_90_days)
                  }}</div>
                  <div class="text-caption"
                    >&gt;60d: {{ fmtCurrency(summary.aged_over_60_days) }}</div
                  >
                </div>
                <v-icon size="36">mdi-alert-octagon-outline</v-icon>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <!-- Tabs -->
        <v-tabs v-model="activeTab" class="mb-4">
          <v-tab value="workspace">Workspace</v-tab>
          <v-tab value="runs">Reconciliation Runs</v-tab>
          <v-tab value="rules">Matching Rules</v-tab>
        </v-tabs>

        <!-- Tab: Workspace -->
        <div v-if="activeTab === 'workspace'">
          <!-- Action Buttons -->
          <v-row class="mb-3">
            <v-col class="d-flex ga-2 flex-wrap">
              <v-btn
                v-if="hasPermission('premiums:reconcile')"
                color="primary"
                prepend-icon="mdi-eye-outline"
                :loading="previewing"
                @click="handlePreview"
              >
                Preview Auto-Match
              </v-btn>
              <v-btn
                v-if="hasPermission('premiums:reconcile')"
                color="primary"
                prepend-icon="mdi-auto-fix"
                :loading="autoMatching"
                @click="handleAutoMatch"
              >
                Run Auto-Match
              </v-btn>
              <v-btn
                variant="outlined"
                color="primary"
                prepend-icon="mdi-link-variant"
                :disabled="
                  selectedPaymentItems.length === 0 ||
                  selectedInvoiceItems.length === 0
                "
                @click="allocateDialog = true"
              >
                Allocate Selected
              </v-btn>
              <v-btn
                variant="outlined"
                prepend-icon="mdi-pencil-off-outline"
                :disabled="!selectedItem"
                @click="openWriteOff"
              >
                Write-Off
              </v-btn>
              <v-btn
                variant="outlined"
                prepend-icon="mdi-cash-refund"
                :disabled="!selectedPaymentItem"
                @click="openRefund"
              >
                Refund
              </v-btn>
            </v-col>
          </v-row>

          <!-- Preview Results -->
          <v-alert
            v-if="previewResult"
            type="info"
            variant="tonal"
            closable
            class="mb-4"
            @click:close="previewResult = null"
          >
            <strong>Preview:</strong>
            {{ previewResult.total_matched }} allocations totalling
            {{ fmtCurrency(previewResult.total_allocated) }},
            {{ previewResult.total_remaining }} payments remaining.
          </v-alert>

          <!-- Two-Panel Layout -->
          <v-row>
            <v-col cols="12" md="6">
              <v-card variant="outlined">
                <v-card-title class="text-subtitle-1 pa-3 d-flex align-center">
                  <span>Unallocated Payments</span>
                  <v-spacer />
                  <v-chip size="small" color="warning">{{
                    paymentItems.length
                  }}</v-chip>
                </v-card-title>
                <v-card-text class="pa-0">
                  <ag-grid-vue
                    class="ag-theme-balham"
                    style="height: 400px; width: 100%"
                    :column-defs="paymentItemCols"
                    :row-data="paymentItems"
                    :default-col-def="{
                      sortable: true,
                      resizable: true,
                      flex: 1
                    }"
                    :loading="loading"
                    row-selection="multiple"
                    @selection-changed="onPaymentItemsSelected"
                  />
                </v-card-text>
              </v-card>
            </v-col>
            <v-col cols="12" md="6">
              <v-card variant="outlined">
                <v-card-title class="text-subtitle-1 pa-3 d-flex align-center">
                  <span>Unpaid / Partial Invoices</span>
                  <v-spacer />
                  <v-chip size="small" color="error">{{
                    invoiceItems.length
                  }}</v-chip>
                </v-card-title>
                <v-card-text class="pa-0">
                  <ag-grid-vue
                    class="ag-theme-balham"
                    style="height: 400px; width: 100%"
                    :column-defs="invoiceItemCols"
                    :row-data="invoiceItems"
                    :default-col-def="{
                      sortable: true,
                      resizable: true,
                      flex: 1
                    }"
                    :loading="loading"
                    row-selection="multiple"
                    @selection-changed="onInvoiceItemsSelected"
                  />
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
        </div>

        <!-- Tab: Reconciliation Runs -->
        <div v-if="activeTab === 'runs'">
          <ag-grid-vue
            class="ag-theme-balham"
            style="height: 500px; width: 100%"
            :column-defs="runCols"
            :row-data="runs"
            :default-col-def="{ sortable: true, resizable: true, flex: 1 }"
            :loading="runsLoading"
            row-selection="single"
            @selection-changed="onRunSelected"
          />
          <v-row v-if="selectedRun" class="mt-3">
            <v-col>
              <v-btn
                variant="outlined"
                color="error"
                prepend-icon="mdi-undo"
                :loading="rollingBack"
                :disabled="selectedRun.status === 'rolled_back'"
                @click="handleRollback"
              >
                Rollback Run #{{ selectedRun.id }}
              </v-btn>
            </v-col>
          </v-row>
        </div>

        <!-- Tab: Matching Rules -->
        <div v-if="activeTab === 'rules'">
          <v-row class="mb-3">
            <v-col class="d-flex ga-2">
              <v-btn
                color="primary"
                prepend-icon="mdi-plus"
                @click="openRuleDialog()"
                >Add Rule</v-btn
              >
            </v-col>
          </v-row>
          <ag-grid-vue
            class="ag-theme-balham"
            style="height: 400px; width: 100%"
            :column-defs="ruleCols"
            :row-data="matchingRules"
            :default-col-def="{ sortable: true, resizable: true, flex: 1 }"
            :loading="rulesLoading"
          />
        </div>
      </template>
    </base-card>

    <!-- Allocate Dialog -->
    <v-dialog v-model="allocateDialog" max-width="560" persistent>
      <v-card>
        <v-card-title>Allocate Payment to Invoices</v-card-title>
        <v-card-text>
          <p v-if="selectedPaymentItems.length === 1" class="mb-3">
            Payment:
            <strong>{{
              fmtCurrency(selectedPaymentItems[0].unallocated_amount)
            }}</strong>
            ({{ selectedPaymentItems[0].scheme_name }})
          </p>
          <v-table density="compact">
            <thead>
              <tr>
                <th>Invoice</th>
                <th>Scheme</th>
                <th>Outstanding</th>
                <th>Allocate</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(inv, idx) in selectedInvoiceItems" :key="inv.id">
                <td>{{ inv.invoice_id }}</td>
                <td>{{ inv.scheme_name }}</td>
                <td>{{ fmtCurrency(inv.unallocated_amount) }}</td>
                <td>
                  <v-text-field
                    v-model.number="allocAmounts[idx]"
                    type="number"
                    variant="outlined"
                    density="compact"
                    hide-details
                    style="max-width: 120px"
                  />
                </td>
              </tr>
            </tbody>
          </v-table>
          <v-text-field
            v-model="allocNotes"
            label="Notes"
            variant="outlined"
            density="compact"
            class="mt-3"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="allocateDialog = false">Cancel</v-btn>
          <v-btn color="primary" :loading="allocating" @click="handleAllocate"
            >Allocate</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Write-Off Dialog -->
    <v-dialog v-model="writeOffDialog" max-width="420" persistent>
      <v-card>
        <v-card-title>Write Off Balance</v-card-title>
        <v-card-text>
          <p class="mb-3" v-if="selectedItem">
            Unallocated:
            <strong>{{ fmtCurrency(selectedItem.unallocated_amount) }}</strong>
          </p>
          <v-text-field
            v-model.number="writeOffForm.amount"
            label="Amount *"
            type="number"
            variant="outlined"
            density="compact"
          />
          <v-text-field
            v-model="writeOffForm.reason"
            label="Reason *"
            variant="outlined"
            density="compact"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="writeOffDialog = false">Cancel</v-btn>
          <v-btn color="primary" :loading="writingOff" @click="handleWriteOff"
            >Write Off</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Refund Dialog -->
    <v-dialog v-model="refundDialog" max-width="420" persistent>
      <v-card>
        <v-card-title>Refund Overpayment</v-card-title>
        <v-card-text>
          <p class="mb-3" v-if="selectedPaymentItem">
            Unallocated:
            <strong>{{
              fmtCurrency(selectedPaymentItem.unallocated_amount)
            }}</strong>
          </p>
          <v-text-field
            v-model.number="refundForm.amount"
            label="Amount *"
            type="number"
            variant="outlined"
            density="compact"
          />
          <v-text-field
            v-model="refundForm.reason"
            label="Reason *"
            variant="outlined"
            density="compact"
          />
          <v-select
            v-model="refundForm.refund_method"
            label="Method *"
            :items="['eft', 'cheque']"
            variant="outlined"
            density="compact"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="refundDialog = false">Cancel</v-btn>
          <v-btn color="primary" :loading="refunding" @click="handleRefund"
            >Refund</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Matching Rule Dialog -->
    <v-dialog v-model="ruleDialog" max-width="560" persistent>
      <v-card>
        <v-card-title
          >{{ ruleForm.id ? 'Edit' : 'New' }} Matching Rule</v-card-title
        >
        <v-card-text>
          <v-row>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="ruleForm.name"
                label="Name *"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model.number="ruleForm.priority"
                label="Priority *"
                type="number"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12">
              <v-text-field
                v-model="ruleForm.description"
                label="Description"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-select
                v-model="ruleForm.strategy"
                label="Strategy *"
                :items="strategyOptions"
                item-title="label"
                item-value="value"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="6" md="3">
              <v-select
                v-model="ruleForm.tolerance_type"
                label="Tolerance Type"
                :items="['absolute', 'percentage']"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="6" md="3">
              <v-text-field
                v-model.number="ruleForm.tolerance_value"
                label="Tolerance Val"
                type="number"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="6">
              <v-checkbox
                v-model="ruleForm.allow_partial"
                label="Allow Partial"
                density="compact"
              />
            </v-col>
            <v-col cols="6">
              <v-checkbox
                v-model="ruleForm.allow_multi_invoice"
                label="Allow Multi-Invoice"
                density="compact"
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="ruleDialog = false">Cancel</v-btn>
          <v-btn color="primary" :loading="savingRule" @click="handleSaveRule"
            >Save</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

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
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-balham.css'
import { AgGridVue } from 'ag-grid-vue3'
import PremiumManagementService from '@/renderer/api/PremiumManagementService'
import BaseCard from '@/renderer/components/BaseCard.vue'
import { useStatusBarStore } from '@/renderer/store/statusBar'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

const statusBarStore = useStatusBarStore()
const { hasPermission } = usePermissionCheck()

// ─── State ──────────────────────────────────────────────────────────────────

const loading = ref(false)
const autoMatching = ref(false)
const previewing = ref(false)
const allocating = ref(false)
const writingOff = ref(false)
const refunding = ref(false)
const runsLoading = ref(false)
const rulesLoading = ref(false)
const savingRule = ref(false)
const rollingBack = ref(false)
const snackbar = ref(false)
const snackbarText = ref('')
const snackbarColor = ref('success')

const activeTab = ref('workspace')

const summary = ref({
  total_unallocated_payments: 0,
  unallocated_payment_count: 0,
  total_unpaid_invoices: 0,
  unpaid_invoice_count: 0,
  suspense_balance: 0,
  suspense_count: 0,
  aged_over_30_days: 0,
  aged_over_60_days: 0,
  aged_over_90_days: 0,
  recent_runs: [] as any[]
})

const paymentItems = ref<any[]>([])
const invoiceItems = ref<any[]>([])
const selectedPaymentItems = ref<any[]>([])
const selectedInvoiceItems = ref<any[]>([])
const previewResult = ref<any>(null)

// Allocate
const allocateDialog = ref(false)
const allocAmounts = ref<number[]>([])
const allocNotes = ref('')

// Write-off
const writeOffDialog = ref(false)
const writeOffForm = ref({ amount: 0, reason: '' })

// Refund
const refundDialog = ref(false)
const refundForm = ref({ amount: 0, reason: '', refund_method: 'eft' })

// Runs
const runs = ref<any[]>([])
const selectedRun = ref<any>(null)

// Matching Rules
const matchingRules = ref<any[]>([])
const ruleDialog = ref(false)
const ruleForm = ref({
  id: 0,
  name: '',
  description: '',
  priority: 10,
  strategy: 'exact_reference',
  tolerance_type: 'absolute',
  tolerance_value: 0.01,
  allow_partial: true,
  allow_multi_invoice: true
})

const strategyOptions = [
  { label: 'Exact Reference Match', value: 'exact_reference' },
  { label: 'Scheme + Exact Amount', value: 'scheme_amount' },
  { label: 'Scheme + Amount (Tolerance)', value: 'scheme_amount_tolerance' },
  { label: 'Scheme + Date Range', value: 'scheme_date_range' },
  { label: 'Amount Only', value: 'amount_only' }
]

// ─── Computed ───────────────────────────────────────────────────────────────

const selectedItem = computed(
  () => selectedPaymentItems.value[0] ?? selectedInvoiceItems.value[0] ?? null
)
const selectedPaymentItem = computed(
  () => selectedPaymentItems.value[0] ?? null
)

// ─── Column Definitions ─────────────────────────────────────────────────────

const paymentItemCols = [
  { headerName: 'Scheme', field: 'scheme_name', minWidth: 120 },
  {
    headerName: 'Original',
    field: 'original_amount',
    valueFormatter: (p: any) => fmtCurrency(p.value),
    minWidth: 110
  },
  {
    headerName: 'Allocated',
    field: 'allocated_amount',
    valueFormatter: (p: any) => fmtCurrency(p.value),
    minWidth: 110
  },
  {
    headerName: 'Unallocated',
    field: 'unallocated_amount',
    valueFormatter: (p: any) => fmtCurrency(p.value),
    minWidth: 120
  },
  { headerName: 'Status', field: 'status', minWidth: 80 },
  { headerName: 'Age (days)', field: 'age_in_days', minWidth: 90 },
  { headerName: 'Priority', field: 'priority', minWidth: 80 },
  { headerName: 'Assigned', field: 'assigned_to', minWidth: 100 }
]

const invoiceItemCols = [
  { headerName: 'Scheme', field: 'scheme_name', minWidth: 120 },
  {
    headerName: 'Original',
    field: 'original_amount',
    valueFormatter: (p: any) => fmtCurrency(p.value),
    minWidth: 110
  },
  {
    headerName: 'Allocated',
    field: 'allocated_amount',
    valueFormatter: (p: any) => fmtCurrency(p.value),
    minWidth: 110
  },
  {
    headerName: 'Unallocated',
    field: 'unallocated_amount',
    valueFormatter: (p: any) => fmtCurrency(p.value),
    minWidth: 120
  },
  { headerName: 'Status', field: 'status', minWidth: 80 },
  { headerName: 'Age (days)', field: 'age_in_days', minWidth: 90 }
]

const runCols = [
  { headerName: 'ID', field: 'id', minWidth: 60 },
  { headerName: 'Date', field: 'run_date', minWidth: 110 },
  { headerName: 'Type', field: 'run_type', minWidth: 90 },
  { headerName: 'Status', field: 'status', minWidth: 100 },
  { headerName: 'Processed', field: 'total_processed', minWidth: 120 },
  { headerName: 'Matched', field: 'total_matched', minWidth: 110 },
  { headerName: 'Unmatched', field: 'total_unmatched', minWidth: 120 },
  {
    headerName: 'Allocated',
    field: 'total_allocated',
    valueFormatter: (p: any) => fmtCurrency(p.value),
    minWidth: 110
  },
  { headerName: 'By', field: 'initiated_by', minWidth: 100 },
  { headerName: 'Rule Set', field: 'matching_rule_set', minWidth: 100 }
]

const ruleCols = [
  { headerName: 'Priority', field: 'priority', maxWidth: 80 },
  { headerName: 'Name', field: 'name', minWidth: 160 },
  { headerName: 'Strategy', field: 'strategy', minWidth: 140 },
  { headerName: 'Tolerance', field: 'tolerance_value', maxWidth: 90 },
  { headerName: 'Type', field: 'tolerance_type', maxWidth: 90 },
  {
    headerName: 'Partial',
    field: 'allow_partial',
    maxWidth: 70,
    valueFormatter: (p: any) => (p.value ? 'Yes' : 'No')
  },
  {
    headerName: 'Multi',
    field: 'allow_multi_invoice',
    maxWidth: 70,
    valueFormatter: (p: any) => (p.value ? 'Yes' : 'No')
  },
  {
    headerName: 'Active',
    field: 'is_active',
    maxWidth: 70,
    valueFormatter: (p: any) => (p.value ? 'Yes' : 'No')
  }
]

// ─── Helpers ────────────────────────────────────────────────────────────────

function fmtCurrency(val: number) {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR',
    maximumFractionDigits: 0
  }).format(val ?? 0)
}

function showSnack(msg: string, color = 'success') {
  snackbarText.value = msg
  snackbarColor.value = color
  snackbar.value = true
}

function onPaymentItemsSelected(e: any) {
  selectedPaymentItems.value = e.api.getSelectedRows()
}

function onInvoiceItemsSelected(e: any) {
  selectedInvoiceItems.value = e.api.getSelectedRows()
}

function onRunSelected(e: any) {
  const rows = e.api.getSelectedRows()
  selectedRun.value = rows[0] ?? null
}

// ─── Data Loading ───────────────────────────────────────────────────────────

async function loadSummary() {
  try {
    const res = await PremiumManagementService.getReconciliationSummary()
    summary.value = res.data.data ?? summary.value
    statusBarStore.set([
      {
        icon: 'mdi-bank-outline',
        text: `Unallocated: ${summary.value.unallocated_payment_count} · ${fmtCurrency(summary.value.total_unallocated_payments)}`,
        severity: summary.value.unallocated_payment_count > 0 ? 'warn' : 'info'
      },
      {
        icon: 'mdi-file-document-alert-outline',
        text: `Unpaid invoices: ${summary.value.unpaid_invoice_count} · ${fmtCurrency(summary.value.total_unpaid_invoices)}`,
        severity: summary.value.unpaid_invoice_count > 0 ? 'error' : 'info'
      }
    ])
  } catch (e) {
    console.error('Failed to load summary', e)
  }
}

async function loadItems() {
  loading.value = true
  try {
    const [payRes, invRes] = await Promise.all([
      PremiumManagementService.getReconciliationItems({
        type: 'payment',
        status: ''
      }),
      PremiumManagementService.getReconciliationItems({
        type: 'invoice',
        status: ''
      })
    ])
    paymentItems.value = (payRes.data.data ?? []).filter((i: any) =>
      ['open', 'partial'].includes(i.status)
    )
    invoiceItems.value = (invRes.data.data ?? []).filter((i: any) =>
      ['open', 'partial'].includes(i.status)
    )
  } catch (e) {
    console.error('Failed to load items', e)
  } finally {
    loading.value = false
  }
}

async function loadRuns() {
  runsLoading.value = true
  try {
    const res = await PremiumManagementService.getReconciliationRuns()
    runs.value = res.data.data ?? []
  } catch (e) {
    console.error('Failed to load runs', e)
  } finally {
    runsLoading.value = false
  }
}

async function loadRules() {
  rulesLoading.value = true
  try {
    const res = await PremiumManagementService.getMatchingRules()
    matchingRules.value = res.data.data ?? []
  } catch (e) {
    console.error('Failed to load rules', e)
  } finally {
    rulesLoading.value = false
  }
}

async function refreshAll() {
  await Promise.all([loadSummary(), loadItems()])
}

// ─── Actions ────────────────────────────────────────────────────────────────

async function handlePreview() {
  previewing.value = true
  try {
    const res = await PremiumManagementService.runAutoMatchV2({ dry_run: true })
    previewResult.value = res.data.data
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Preview failed', 'error')
  } finally {
    previewing.value = false
  }
}

async function handleAutoMatch() {
  autoMatching.value = true
  try {
    const res = await PremiumManagementService.runAutoMatchV2()
    const r = res.data.data
    showSnack(
      `Auto-matched ${r.total_matched} allocations. ${r.total_unmatched} remaining.`
    )
    previewResult.value = null
    await refreshAll()
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Auto-match failed', 'error')
  } finally {
    autoMatching.value = false
  }
}

async function handleAllocate() {
  if (selectedPaymentItems.value.length === 0) return
  allocating.value = true
  try {
    const payment = selectedPaymentItems.value[0]
    const allocations = selectedInvoiceItems.value
      .map((inv: any, idx: number) => ({
        invoice_id: inv.invoice_id,
        amount: allocAmounts.value[idx] || 0
      }))
      .filter((a: any) => a.amount > 0)

    await PremiumManagementService.allocatePayment({
      payment_id: payment.payment_id,
      allocations,
      notes: allocNotes.value
    })
    allocateDialog.value = false
    showSnack('Payment allocated successfully')
    await refreshAll()
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Allocation failed', 'error')
  } finally {
    allocating.value = false
  }
}

function openWriteOff() {
  if (!selectedItem.value) return
  writeOffForm.value = {
    amount: selectedItem.value.unallocated_amount,
    reason: ''
  }
  writeOffDialog.value = true
}

async function handleWriteOff() {
  if (!selectedItem.value) return
  writingOff.value = true
  try {
    await PremiumManagementService.writeOffBalance({
      reconciliation_item_id: selectedItem.value.id,
      amount: writeOffForm.value.amount,
      reason: writeOffForm.value.reason,
      invoice_id: selectedItem.value.invoice_id
    })
    writeOffDialog.value = false
    showSnack('Balance written off')
    await refreshAll()
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Write-off failed', 'error')
  } finally {
    writingOff.value = false
  }
}

function openRefund() {
  if (!selectedPaymentItem.value) return
  refundForm.value = {
    amount: selectedPaymentItem.value.unallocated_amount,
    reason: '',
    refund_method: 'eft'
  }
  refundDialog.value = true
}

async function handleRefund() {
  if (!selectedPaymentItem.value) return
  refunding.value = true
  try {
    await PremiumManagementService.refundOverpayment({
      reconciliation_item_id: selectedPaymentItem.value.id,
      amount: refundForm.value.amount,
      reason: refundForm.value.reason,
      refund_method: refundForm.value.refund_method
    })
    refundDialog.value = false
    showSnack('Refund recorded')
    await refreshAll()
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Refund failed', 'error')
  } finally {
    refunding.value = false
  }
}

async function handleRollback() {
  if (!selectedRun.value) return
  rollingBack.value = true
  try {
    await PremiumManagementService.rollbackRun(selectedRun.value.id)
    showSnack('Run rolled back successfully')
    selectedRun.value = null
    await Promise.all([loadRuns(), refreshAll()])
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Rollback failed', 'error')
  } finally {
    rollingBack.value = false
  }
}

function openRuleDialog(rule?: any) {
  if (rule) {
    ruleForm.value = { ...rule }
  } else {
    ruleForm.value = {
      id: 0,
      name: '',
      description: '',
      priority: matchingRules.value.length + 1,
      strategy: 'exact_reference',
      tolerance_type: 'absolute',
      tolerance_value: 0.01,
      allow_partial: true,
      allow_multi_invoice: true
    }
  }
  ruleDialog.value = true
}

async function handleSaveRule() {
  savingRule.value = true
  try {
    await PremiumManagementService.saveMatchingRule(ruleForm.value)
    ruleDialog.value = false
    showSnack('Rule saved')
    await loadRules()
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Failed to save rule', 'error')
  } finally {
    savingRule.value = false
  }
}

// ─── Watchers ───────────────────────────────────────────────────────────────

watch(allocateDialog, (val) => {
  if (val) {
    allocAmounts.value = selectedInvoiceItems.value.map(
      (inv: any) => inv.unallocated_amount
    )
    allocNotes.value = ''
  }
})

watch(activeTab, (val) => {
  if (val === 'runs') loadRuns()
  if (val === 'rules') loadRules()
})

// ─── Init ───────────────────────────────────────────────────────────────────

onMounted(() => {
  refreshAll()
})
onUnmounted(() => statusBarStore.clear())
</script>
