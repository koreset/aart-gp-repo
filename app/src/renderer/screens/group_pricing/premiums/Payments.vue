<template>
  <v-container fluid>
    <base-card :show-actions="false">
      <template #header>
        <h3 class="mb-0">Payments</h3>
      </template>
      <template #default>
        <!-- Filter Bar -->
        <v-row class="mb-3" align="center">
          <v-col cols="12" md="2">
            <v-select
              v-model="filters.schemeId"
              label="Scheme"
              :items="inForceSchemes"
              item-title="name"
              item-value="id"
              variant="outlined"
              density="compact"
              clearable
              prepend-inner-icon="mdi-office-building-outline"
              :loading="schemesLoading"
            />
          </v-col>
          <v-col cols="12" md="2">
            <v-select
              v-model="filters.method"
              label="Method"
              :items="methodOptions"
              variant="outlined"
              density="compact"
              clearable
            />
          </v-col>
          <v-col cols="12" md="2">
            <v-select
              v-model="filters.status"
              label="Status"
              :items="statusOptions"
              variant="outlined"
              density="compact"
              clearable
            />
          </v-col>
          <v-col cols="12" md="6" class="d-flex ga-2 align-center justify-end">
            <v-btn
              variant="outlined"
              prepend-icon="mdi-refresh"
              @click="loadPayments"
              >Refresh</v-btn
            >
            <v-btn
              v-if="hasPermission('premiums:bulk_import')"
              variant="outlined"
              color="primary"
              prepend-icon="mdi-upload"
              @click="bulkDialog = true"
            >
              Bulk Import
            </v-btn>
            <v-btn
              v-if="hasPermission('premiums:record_payment')"
              color="primary"
              prepend-icon="mdi-plus"
              @click="recordDialog = true"
            >
              Record Payment
            </v-btn>
          </v-col>
        </v-row>

        <!-- Payments Grid -->
        <v-row>
          <v-col cols="12">
            <v-card variant="outlined">
              <v-card-text class="pa-0">
                <ag-grid-vue
                  class="ag-theme-balham"
                  :style="{ height: gridHeight, width: '100%' }"
                  :column-defs="columnDefs"
                  :row-data="payments"
                  :default-col-def="defaultColDef"
                  :loading="loading"
                  :get-row-class="getRowClass"
                />
                <empty-state
                  v-if="!loading && payments.length === 0"
                  title="No payments found"
                  message="Record a payment or import via CSV."
                  icon="mdi-bank-transfer"
                />
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </template>
    </base-card>

    <!-- Record Payment Dialog -->
    <v-dialog v-model="recordDialog" max-width="520" persistent>
      <v-card>
        <v-card-title>Record Payment</v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12" md="6">
              <v-select
                v-model="payForm.scheme_id"
                label="Scheme *"
                :items="inForceSchemes"
                item-title="name"
                item-value="id"
                variant="outlined"
                density="compact"
                prepend-inner-icon="mdi-office-building-outline"
                :loading="schemesLoading"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model.number="payForm.invoice_id"
                label="Invoice ID (optional)"
                type="number"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="payForm.payment_date"
                label="Payment Date *"
                type="date"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-select
                v-model="payForm.method"
                label="Method *"
                :items="methodOptions"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model.number="payForm.amount"
                label="Amount *"
                type="number"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="payForm.bank_reference"
                label="Bank Reference"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="payForm.notes"
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
          <v-btn variant="plain" @click="recordDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="savingPayment"
            @click="handleRecordPayment"
            >Save</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Bulk Import Dialog -->
    <v-dialog v-model="bulkDialog" max-width="600" persistent>
      <v-card>
        <v-card-title>Bulk Import Payments</v-card-title>
        <v-card-text>
          <v-alert type="info" variant="tonal" class="mb-4">
            Upload a CSV with columns: Reference, Amount, Bank Reference,
            Invoice Number.
          </v-alert>
          <v-file-input
            v-model="importFile"
            label="Select CSV file"
            accept=".csv"
            variant="outlined"
            density="compact"
            prepend-icon="mdi-paperclip"
          />
          <!-- Import results -->
          <div v-if="importResult" class="mt-4">
            <v-chip color="success" class="mr-2"
              >Matched: {{ importResult.matched }}</v-chip
            >
            <v-chip color="warning" class="mr-2"
              >Unmatched: {{ importResult.unmatched }}</v-chip
            >
            <v-chip color="error">Errors: {{ importResult.errors }}</v-chip>
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="resetData()">Close</v-btn>
          <v-btn
            color="primary"
            :loading="importing"
            :disabled="!importFile"
            @click="handleBulkImport"
          >
            Import
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Void Confirm Dialog -->
    <v-dialog v-model="voidDialog" max-width="400" persistent>
      <v-card>
        <v-card-title>Void Payment</v-card-title>
        <v-card-text>
          <v-alert type="warning" variant="tonal">
            This will void the payment and reverse any invoice balance update.
            Are you sure?
          </v-alert>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="voidDialog = false">Cancel</v-btn>
          <v-btn color="error" :loading="voiding" @click="handleVoid"
            >Void Payment</v-btn
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
      <template #actions>
        <v-btn variant="text" @click="snackbar = false">Close</v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-balham.css'
import { AgGridVue } from 'ag-grid-vue3'
import PremiumManagementService from '@/renderer/api/PremiumManagementService'
import { useStatusBarStore } from '@/renderer/store/statusBar'
import BaseCard from '@/renderer/components/BaseCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import { useGridHeight } from '@/renderer/composables/useGridHeight'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'
import { useFilterPersistence } from '@/renderer/composables/useFilterPersistence'
import { statusCellRenderer } from '@/renderer/utils/statusCellRenderer'
import { dateFormatter } from '@/renderer/utils/formatters'

const statusBarStore = useStatusBarStore()
const { hasPermission } = usePermissionCheck()

const gridHeight = useGridHeight(280)
const { filters } = useFilterPersistence('payments', {
  schemeId: null as number | null,
  method: null as string | null,
  status: null as string | null
})

const loading = ref(false)
const savingPayment = ref(false)
const schemesLoading = ref(false)
const inForceSchemes = ref<{ id: number; name: string }[]>([])
const importing = ref(false)
const voiding = ref(false)
const recordDialog = ref(false)
const bulkDialog = ref(false)
const voidDialog = ref(false)
const snackbar = ref(false)
const snackbarText = ref('')
const snackbarColor = ref('success')

const payments = ref<any[]>([])
const importFile = ref<File | null>(null)
const importResult = ref<any>(null)
const selectedPaymentId = ref<number | null>(null)

const methodOptions = ['eft', 'cheque', 'cash', 'debit_order']
const statusOptions = ['pending', 'matched', 'unmatched', 'voided']

const today = new Date().toISOString().slice(0, 10)
const payForm = ref({
  scheme_id: null as number | null,
  invoice_id: null as number | null,
  payment_date: today,
  method: 'eft',
  amount: 0,
  bank_reference: '',
  notes: ''
})

const defaultColDef = { sortable: true, filter: true, resizable: true, flex: 1 }
const baseColumns = [
  { headerName: 'Scheme', field: 'scheme_name', minWidth: 160 },
  { headerName: 'Invoice #', field: 'invoice_number' },
  { headerName: 'Date', field: 'payment_date', valueFormatter: dateFormatter },
  { headerName: 'Method', field: 'method' },
  {
    headerName: 'Amount',
    field: 'amount',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  { headerName: 'Bank Ref', field: 'bank_reference' },
  {
    headerName: 'Status',
    field: 'status',
    cellRenderer: (p: any) => statusCellRenderer(p.value)
  }
]
const voidColumn = {
  headerName: 'Actions',
  cellRenderer: () =>
    `<v-btn size="x-small" color="error" variant="text">Void</v-btn>`,
  onCellClicked: (p: any) => {
    selectedPaymentId.value = p.data.id
    voidDialog.value = true
  },
  maxWidth: 100
}
const columnDefs = hasPermission('premiums:void_payment')
  ? [...baseColumns, voidColumn]
  : baseColumns

function getRowClass(params: any) {
  return params.data?.status === 'unmatched' ? 'ag-row-warning' : ''
}

function fmtCurrency(val: number) {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR',
    maximumFractionDigits: 0
  }).format(val ?? 0)
}

async function loadPayments() {
  loading.value = true
  try {
    const params: any = {}
    if (filters.value.schemeId) params.scheme_id = filters.value.schemeId
    if (filters.value.method) params.method = filters.value.method
    if (filters.value.status) params.status = filters.value.status
    const res = await PremiumManagementService.getPayments(params)
    payments.value = res.data.data ?? []
    const unmatched = payments.value.filter(
      (p: any) => p.status === 'unmatched'
    ).length
    const unmatchedTotal = payments.value
      .filter((p: any) => p.status === 'unmatched')
      .reduce((s: number, p: any) => s + p.amount, 0)
    statusBarStore.set([
      {
        icon: 'mdi-bank-outline',
        text:
          unmatched > 0
            ? `Unmatched: ${unmatched} · ${fmtCurrency(unmatchedTotal)}`
            : 'All payments matched',
        severity: unmatched > 0 ? 'warn' : 'info'
      }
    ])
  } catch (e) {
    console.error('Failed to load payments', e)
  } finally {
    loading.value = false
  }
}

async function handleRecordPayment() {
  if (!payForm.value.scheme_id) return
  savingPayment.value = true
  try {
    await PremiumManagementService.recordPayment({
      scheme_id: payForm.value.scheme_id!,
      invoice_id: payForm.value.invoice_id || undefined,
      payment_date: payForm.value.payment_date,
      method: payForm.value.method,
      amount: payForm.value.amount,
      bank_reference: payForm.value.bank_reference,
      notes: payForm.value.notes
    })
    recordDialog.value = false
    await loadPayments()
    showSnack('Payment recorded')
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Failed to record payment', 'error')
  } finally {
    savingPayment.value = false
  }
}

async function handleBulkImport() {
  if (!importFile.value) return
  importing.value = true
  try {
    const formData = new FormData()
    formData.append('file', importFile.value)
    const res = await PremiumManagementService.bulkImportPayments(formData)
    importResult.value = res.data.data
    await loadPayments()
    showSnack(`Import complete: ${importResult.value.matched} matched`)
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Import failed', 'error')
  } finally {
    importing.value = false
  }
}

async function handleVoid() {
  if (!selectedPaymentId.value) return
  voiding.value = true
  try {
    await PremiumManagementService.voidPayment(selectedPaymentId.value)
    voidDialog.value = false
    await loadPayments()
    showSnack('Payment voided')
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Failed to void payment', 'error')
  } finally {
    voiding.value = false
  }
}

function showSnack(msg: string, color = 'success') {
  snackbarText.value = msg
  snackbarColor.value = color
  snackbar.value = true
}

const resetData = () => {
  bulkDialog.value = false
  importResult.value = null
}

async function loadInForceSchemes() {
  schemesLoading.value = true
  try {
    const res = await PremiumManagementService.getInForceSchemes()
    inForceSchemes.value = (res.data ?? []).map((s: any) => ({
      id: s.id,
      name: s.name
    }))
  } catch (e) {
    console.error('Failed to load in-force schemes', e)
  } finally {
    schemesLoading.value = false
  }
}

onMounted(() => {
  loadInForceSchemes()
  loadPayments()
})
onUnmounted(() => statusBarStore.clear())
</script>
