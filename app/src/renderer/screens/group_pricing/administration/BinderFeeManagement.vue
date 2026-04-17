<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex justify-space-between align-center">
              <div class="d-flex align-center">
                <span class="headline">Binder Fee Management</span>
                <v-tooltip location="right">
                  <template #activator="{ props }">
                    <v-icon
                      v-bind="props"
                      size="small"
                      color="info"
                      class="ml-2"
                    >
                      mdi-information-outline
                    </v-icon>
                  </template>
                  <span
                    >Maximum Binder Fee Percentage and Maximum Outsourced Fee
                    Percentage are both exclusive of VAT.</span
                  >
                </v-tooltip>
              </div>
              <v-btn color="primary" @click="openAddDialog">
                Add Binder Fee
              </v-btn>
            </div>
          </template>
          <template #default>
            <data-grid
              :column-defs="columnDefs"
              :row-data="binderFees"
              :pagination="true"
            />
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Add / Edit dialog -->
    <v-dialog v-model="formDialog" persistent max-width="600px">
      <base-card>
        <template #header>
          <span class="headline">
            {{ isEditMode ? 'Edit Binder Fee' : 'Add Binder Fee' }}
          </span>
        </template>
        <template #default>
          <v-form @submit.prevent="saveBinderFee">
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="binderholderName"
                  v-bind="binderholderNameAttrs"
                  variant="outlined"
                  density="compact"
                  label="Binderholder Name"
                  placeholder="Enter binderholder name"
                  :error-messages="errors.binderholder_name"
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-select
                  v-model="riskRateCode"
                  v-bind="riskRateCodeAttrs"
                  variant="outlined"
                  density="compact"
                  label="Risk Rate Code"
                  placeholder="Select a risk rate code"
                  :items="riskRateCodes"
                  :error-messages="errors.risk_rate_code"
                ></v-select>
              </v-col>
              <v-col cols="6">
                <v-text-field
                  v-model="maximumBinderFee"
                  v-bind="maximumBinderFeeAttrs"
                  type="number"
                  step="0.01"
                  min="0"
                  max="100"
                  suffix="%"
                  variant="outlined"
                  density="compact"
                  label="Maximum Binder Fee Percentage"
                  placeholder="e.g. 15"
                  hint="Exclusive of VAT"
                  persistent-hint
                  :error-messages="errors.maximum_binder_fee"
                ></v-text-field>
              </v-col>
              <v-col cols="6">
                <v-text-field
                  v-model="maximumOutsourceFee"
                  v-bind="maximumOutsourceFeeAttrs"
                  type="number"
                  step="0.01"
                  min="0"
                  max="100"
                  suffix="%"
                  variant="outlined"
                  density="compact"
                  label="Maximum Outsourced Fee Percentage"
                  placeholder="e.g. 5"
                  hint="Exclusive of VAT"
                  persistent-hint
                  :error-messages="errors.maximum_outsource_fee"
                ></v-text-field>
              </v-col>
            </v-row>
          </v-form>
        </template>
        <template #actions>
          <v-spacer></v-spacer>
          <v-btn color="grey" @click="closeFormDialog">Cancel</v-btn>
          <v-btn color="primary" :loading="saving" @click="saveBinderFee">
            {{ isEditMode ? 'Update' : 'Create' }}
          </v-btn>
        </template>
      </base-card>
    </v-dialog>

    <!-- Delete confirmation -->
    <v-dialog v-model="deleteDialog" max-width="420px">
      <v-card>
        <v-card-title>Delete Binder Fee</v-card-title>
        <v-card-text>
          Delete the binder fee for
          <strong>{{ pendingDelete?.binderholder_name }}</strong>
          ({{ pendingDelete?.risk_rate_code }})? This action cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey" @click="deleteDialog = false">Cancel</v-btn>
          <v-btn color="error" :loading="deleting" @click="confirmDelete">
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar.show" :timeout="4000" :color="snackbar.color">
      {{ snackbar.message }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted, h, markRaw } from 'vue'
import { useForm } from 'vee-validate'
import * as yup from 'yup'
import BaseCard from '@/renderer/components/BaseCard.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

interface BinderFee {
  id?: number
  binderholder_name: string
  risk_rate_code: string
  maximum_binder_fee: number
  maximum_outsource_fee: number
  creation_date?: string
  created_by?: string
}

const binderFees = ref<BinderFee[]>([])
const riskRateCodes = ref<string[]>([])

const formDialog = ref(false)
const isEditMode = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)

const deleteDialog = ref(false)
const pendingDelete = ref<BinderFee | null>(null)
const deleting = ref(false)

const snackbar = ref<{ show: boolean; message: string; color: string }>({
  show: false,
  message: '',
  color: 'success'
})

// ---------------------------------------------------------------------------
// Form
// ---------------------------------------------------------------------------
const validationSchema = yup.object({
  binderholder_name: yup
    .string()
    .required('Binderholder name is required')
    .max(255, 'Binderholder name must be 255 characters or fewer'),
  risk_rate_code: yup.string().required('Risk rate code is required'),
  maximum_binder_fee: yup
    .number()
    .typeError('Maximum binder fee percentage must be a number')
    .required('Maximum binder fee percentage is required')
    .min(0, 'Maximum binder fee percentage must be 0 or greater')
    .max(100, 'Maximum binder fee percentage must be 100 or less'),
  maximum_outsource_fee: yup
    .number()
    .typeError('Maximum outsourced fee percentage must be a number')
    .required('Maximum outsourced fee percentage is required')
    .min(0, 'Maximum outsourced fee percentage must be 0 or greater')
    .max(100, 'Maximum outsourced fee percentage must be 100 or less')
})

const { defineField, errors, resetForm, validate } = useForm({
  validationSchema,
  initialValues: {
    binderholder_name: '',
    risk_rate_code: '',
    maximum_binder_fee: 0,
    maximum_outsource_fee: 0
  }
})

const [binderholderName, binderholderNameAttrs] =
  defineField('binderholder_name')
const [riskRateCode, riskRateCodeAttrs] = defineField('risk_rate_code')
const [maximumBinderFee, maximumBinderFeeAttrs] =
  defineField('maximum_binder_fee')
const [maximumOutsourceFee, maximumOutsourceFeeAttrs] = defineField(
  'maximum_outsource_fee'
)

// ---------------------------------------------------------------------------
// Percent ↔ decimal conversion
// ---------------------------------------------------------------------------
// Decimals are the canonical representation on the server (e.g. 0.15).
// The form and grid display percentages (e.g. 15).
function decimalToPercent(value: unknown): number {
  const n = Number(value)
  if (!Number.isFinite(n)) return 0
  // Round to 4 decimals to avoid floating-point noise (e.g. 0.10000000000001).
  return Math.round(n * 100 * 10000) / 10000
}

function percentToDecimal(value: unknown): number {
  const n = Number(value)
  if (!Number.isFinite(n)) return 0
  return Math.round((n / 100) * 1000000) / 1000000
}

// ---------------------------------------------------------------------------
// Grid
// ---------------------------------------------------------------------------
const percentFormatter = (params: any) => {
  const value = params.value
  if (value === null || value === undefined || value === '') return ''
  const pct = decimalToPercent(value)
  return (
    pct.toLocaleString(undefined, {
      minimumFractionDigits: 2,
      maximumFractionDigits: 4
    }) + '%'
  )
}

const dateFormatter = (params: any) => {
  const value = params.value
  if (!value) return ''
  const d = new Date(value)
  if (isNaN(d.getTime())) return String(value)
  return d.toLocaleString(undefined, {
    year: 'numeric',
    month: 'short',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// A tiny inline cell renderer that emits edit/delete events.
const ActionsCell = markRaw({
  props: ['params'],
  emits: [],
  setup(props: any) {
    const buttonStyle =
      'background: transparent; border: 0; cursor: pointer; padding: 0 8px; font: inherit; line-height: 1;'
    return () =>
      h(
        'div',
        {
          style:
            'display: flex; align-items: center; justify-content: center; gap: 12px; height: 100%; width: 100%;'
        },
        [
          h(
            'button',
            {
              type: 'button',
              class: 'text-primary',
              style: buttonStyle,
              onClick: () => openEditDialog(props.params.data)
            },
            'Edit'
          ),
          h(
            'button',
            {
              type: 'button',
              class: 'text-error',
              style: buttonStyle,
              onClick: () => requestDelete(props.params.data)
            },
            'Delete'
          )
        ]
      )
  }
})

const columnDefs = [
  {
    field: 'binderholder_name',
    headerName: 'Binderholder Name',
    sortable: true,
    filter: true,
    flex: 2
  },
  {
    field: 'risk_rate_code',
    headerName: 'Risk Rate Code',
    sortable: true,
    filter: true,
    flex: 1
  },
  {
    field: 'maximum_binder_fee',
    headerName: 'Maximum Binder Fee Percentage',
    headerTooltip: 'Exclusive of VAT',
    sortable: true,
    filter: 'agNumberColumnFilter',
    type: 'rightAligned',
    valueFormatter: percentFormatter,
    flex: 1
  },
  {
    field: 'maximum_outsource_fee',
    headerName: 'Maximum Outsourced Fee Percentage',
    headerTooltip: 'Exclusive of VAT',
    sortable: true,
    filter: 'agNumberColumnFilter',
    type: 'rightAligned',
    valueFormatter: percentFormatter,
    flex: 1
  },
  {
    field: 'created_by',
    headerName: 'Created By',
    sortable: true,
    filter: true,
    flex: 1
  },
  {
    field: 'creation_date',
    headerName: 'Creation Date',
    sortable: true,
    filter: 'agDateColumnFilter',
    valueFormatter: dateFormatter,
    flex: 1
  },
  {
    headerName: 'Actions',
    cellRenderer: ActionsCell,
    cellRendererParams: (params: any) => ({ params }),
    sortable: false,
    filter: false,
    width: 160,
    headerClass: 'ag-center-header',
    cellStyle: {
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: 0
    }
  }
]

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------
function openAddDialog() {
  isEditMode.value = false
  editingId.value = null
  resetForm({
    values: {
      binderholder_name: '',
      risk_rate_code: '',
      maximum_binder_fee: 0,
      maximum_outsource_fee: 0
    }
  })
  formDialog.value = true
}

function openEditDialog(row: BinderFee) {
  isEditMode.value = true
  editingId.value = row.id ?? null
  resetForm({
    values: {
      binderholder_name: row.binderholder_name,
      risk_rate_code: row.risk_rate_code,
      // Stored as decimals (0.15), displayed as percentages (15).
      maximum_binder_fee: decimalToPercent(row.maximum_binder_fee),
      maximum_outsource_fee: decimalToPercent(row.maximum_outsource_fee)
    }
  })
  formDialog.value = true
}

function closeFormDialog() {
  formDialog.value = false
}

async function saveBinderFee() {
  const result = await validate()
  if (!result.valid) return

  const payload: BinderFee = {
    binderholder_name: String(binderholderName.value ?? '').trim(),
    risk_rate_code: String(riskRateCode.value ?? '').trim(),
    // Form values are percentages (e.g. 15); backend stores decimals (0.15).
    maximum_binder_fee: percentToDecimal(maximumBinderFee.value),
    maximum_outsource_fee: percentToDecimal(maximumOutsourceFee.value)
  }

  saving.value = true
  try {
    if (isEditMode.value && editingId.value !== null) {
      await GroupPricingService.updateBinderFee(editingId.value, payload)
      showMessage('Binder fee updated.', 'success')
    } else {
      await GroupPricingService.createBinderFee(payload)
      showMessage('Binder fee created.', 'success')
    }
    formDialog.value = false
    await loadBinderFees()
  } catch (err: any) {
    const status = err?.response?.status
    if (status === 409) {
      showMessage(
        'A binder fee for this binderholder and risk rate code already exists.',
        'error'
      )
    } else {
      showMessage(
        err?.response?.data?.error ||
          err?.response?.data ||
          'Failed to save binder fee.',
        'error'
      )
    }
  } finally {
    saving.value = false
  }
}

function requestDelete(row: BinderFee) {
  pendingDelete.value = row
  deleteDialog.value = true
}

async function confirmDelete() {
  if (!pendingDelete.value?.id) {
    deleteDialog.value = false
    return
  }
  deleting.value = true
  try {
    await GroupPricingService.deleteBinderFee(pendingDelete.value.id)
    showMessage('Binder fee deleted.', 'success')
    deleteDialog.value = false
    pendingDelete.value = null
    await loadBinderFees()
  } catch (err: any) {
    showMessage(
      err?.response?.data?.error ||
        err?.response?.data ||
        'Failed to delete binder fee.',
      'error'
    )
  } finally {
    deleting.value = false
  }
}

function showMessage(message: string, color: 'success' | 'error') {
  snackbar.value = { show: true, message, color }
}

// ---------------------------------------------------------------------------
// Data loading
// ---------------------------------------------------------------------------
async function loadBinderFees() {
  try {
    const response = await GroupPricingService.getBinderFees()
    binderFees.value = response.data ?? []
  } catch (err: any) {
    showMessage(
      err?.response?.data?.error || 'Failed to load binder fees.',
      'error'
    )
  }
}

async function loadRiskRateCodes() {
  try {
    // Source of truth for binder fees: the risk rate codes registered on
    // group_pricing_parameters (see GetGPTableRiskCodes in the backend).
    const response = await GroupPricingService.getRiskRateCodes(
      'grouppricingparameters'
    )
    riskRateCodes.value = response.data ?? []
  } catch {
    // Dropdown falls back to empty — user can retry; risk_rate_code is a
    // free-form string on the backend.
    riskRateCodes.value = []
  }
}

onMounted(async () => {
  await Promise.all([loadBinderFees(), loadRiskRateCodes()])
})
</script>

<style>
/* Centre the "Actions" column header label. Unscoped so it reaches the
   AG Grid shadow-free internals. */
.ag-center-header .ag-header-cell-label {
  justify-content: center;
}
</style>
