<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex justify-space-between align-center">
              <div class="d-flex align-center">
                <span class="headline">Commission Structures</span>
                <v-tooltip location="right" max-width="320">
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
                  <span>
                    Bands are applied progressively — each band's rate covers
                    only the portion of annualised premium that falls within it.
                    The applicable rate may not exceed the regulated maximum.
                    Direct channel is always 0%.
                  </span>
                </v-tooltip>
              </div>
              <div class="d-flex ga-2">
                <v-btn
                  color="secondary"
                  variant="outlined"
                  @click="uploadDialog = true"
                >
                  Upload CSV
                </v-btn>
                <v-btn color="primary" @click="openAddDialog"> Add Band </v-btn>
              </div>
            </div>
          </template>
          <template #default>
            <v-tabs
              v-model="activeChannel"
              color="primary"
              class="mb-4"
              @update:model-value="onChannelChange"
            >
              <v-tab v-for="ch in channels" :key="ch.value" :value="ch.value">{{
                ch.title
              }}</v-tab>
              <v-tab value="direct" disabled>Direct (0% — no scale)</v-tab>
            </v-tabs>

            <v-row class="mb-2">
              <v-col cols="12" md="6">
                <v-select
                  v-model="selectedHolder"
                  :items="holderOptions"
                  item-title="title"
                  item-value="value"
                  variant="outlined"
                  density="compact"
                  label="Holder"
                  hide-details
                  @update:model-value="loadBands"
                ></v-select>
              </v-col>
            </v-row>

            <v-alert
              v-if="contiguityIssues.length > 0"
              type="warning"
              variant="tonal"
              density="compact"
              class="mb-3"
            >
              <div class="text-subtitle-2">
                Structure issues for
                {{ channelLabel(activeChannel) }} —
                {{ selectedHolder === '' ? 'Default' : selectedHolder }}:
              </div>
              <ul class="mt-1">
                <li v-for="msg in contiguityIssues" :key="msg">{{ msg }}</li>
              </ul>
            </v-alert>

            <data-grid
              :column-defs="columnDefs"
              :row-data="bands"
              :pagination="false"
            />
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Add / Edit dialog -->
    <v-dialog v-model="formDialog" persistent max-width="640px">
      <base-card>
        <template #header>
          <span class="headline">
            {{ isEditMode ? 'Edit Band' : 'Add Band' }}
          </span>
        </template>
        <template #default>
          <v-form @submit.prevent="saveBand">
            <v-row>
              <v-col cols="6">
                <v-text-field
                  :model-value="channelLabel(formChannel)"
                  readonly
                  variant="outlined"
                  density="compact"
                  label="Channel"
                ></v-text-field>
              </v-col>
              <v-col cols="6">
                <v-autocomplete
                  v-model="formHolder"
                  :items="holderOptions"
                  item-title="title"
                  item-value="value"
                  variant="outlined"
                  density="compact"
                  label="Holder"
                  hint="Leave as Default to apply to all holders on this channel"
                  persistent-hint
                ></v-autocomplete>
              </v-col>
              <v-col cols="6">
                <v-text-field
                  v-model="formattedLowerBound"
                  v-bind="lowerBoundAttrs"
                  type="text"
                  inputmode="numeric"
                  variant="outlined"
                  density="compact"
                  label="Lower Bound"
                  prefix="R"
                  :error-messages="errors.lower_bound"
                ></v-text-field>
              </v-col>
              <v-col cols="6">
                <v-text-field
                  v-if="!unboundedUpper"
                  v-model="formattedUpperBound"
                  v-bind="upperBoundAttrs"
                  type="text"
                  inputmode="numeric"
                  variant="outlined"
                  density="compact"
                  label="Upper Bound"
                  prefix="R"
                  :error-messages="errors.upper_bound"
                ></v-text-field>
                <v-text-field
                  v-else
                  model-value="Unbounded"
                  readonly
                  variant="outlined"
                  density="compact"
                  label="Upper Bound"
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-checkbox
                  v-model="unboundedUpper"
                  label="Final band (unbounded — covers all premium above the lower bound)"
                  density="compact"
                  hide-details
                ></v-checkbox>
              </v-col>
              <v-col cols="6">
                <v-text-field
                  v-model="maximumCommission"
                  v-bind="maximumCommissionAttrs"
                  type="number"
                  step="0.01"
                  min="0"
                  max="100"
                  suffix="%"
                  variant="outlined"
                  density="compact"
                  label="Maximum Commission Percentage"
                  placeholder="e.g. 7.5"
                  hint="Regulated cap — applicable rate cannot exceed this"
                  persistent-hint
                  :error-messages="errors.maximum_commission"
                ></v-text-field>
              </v-col>
              <v-col cols="6">
                <v-text-field
                  v-model="applicableRate"
                  v-bind="applicableRateAttrs"
                  type="number"
                  step="0.01"
                  min="0"
                  max="100"
                  suffix="%"
                  variant="outlined"
                  density="compact"
                  label="Applicable Rate Percentage"
                  placeholder="e.g. 7.5"
                  hint="What pricing actually uses"
                  persistent-hint
                  :error-messages="errors.applicable_rate"
                ></v-text-field>
              </v-col>
            </v-row>
          </v-form>
        </template>
        <template #actions>
          <v-spacer></v-spacer>
          <v-btn color="grey" @click="closeFormDialog">Cancel</v-btn>
          <v-btn color="primary" :loading="saving" @click="saveBand">
            {{ isEditMode ? 'Update' : 'Create' }}
          </v-btn>
        </template>
      </base-card>
    </v-dialog>

    <!-- Delete confirmation -->
    <v-dialog v-model="deleteDialog" max-width="440px">
      <v-card>
        <v-card-title>Delete Commission Band</v-card-title>
        <v-card-text>
          Delete the
          <strong>{{ channelLabel(pendingDelete?.channel) }}</strong>
          band starting at
          <strong>{{ formatCurrency(pendingDelete?.lower_bound) }}</strong
          >? This may leave a gap in the structure until you fix it.
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

    <!-- CSV upload dialog -->
    <v-dialog v-model="uploadDialog" max-width="560px">
      <v-card>
        <v-card-title>Upload Commission Structure CSV</v-card-title>
        <v-card-text>
          <p class="text-body-2 mb-2">
            Columns:
            <code
              >channel, holder_name, lower_bound, upper_bound,
              maximum_commission, applicable_rate</code
            >.
          </p>
          <p class="text-body-2 mb-3">
            Values are decimals (e.g. <code>0.075</code> for 7.5%). Leave
            <code>upper_bound</code> blank on the final band for "unbounded".
            Leave <code>holder_name</code> blank for a
            <strong>default</strong> scale that applies to any holder on that
            channel; supply a holder name (as it appears in the Broker list
            under Metadata) to override for that holder only. Uploading replaces
            existing bands for every <em>(channel, holder)</em> group present in
            the file.
          </p>
          <div class="d-flex align-center mb-3">
            <v-btn
              variant="text"
              color="primary"
              size="small"
              prepend-icon="mdi-download"
              @click="downloadTemplate"
            >
              Download template
            </v-btn>
            <span class="text-caption text-grey ml-2">
              CSV with the expected columns and sample rows.
            </span>
          </div>
          <v-file-input
            v-model="uploadFile"
            accept=".csv"
            label="Select CSV file"
            density="compact"
            variant="outlined"
            show-size
            prepend-icon="mdi-file-delimited-outline"
          ></v-file-input>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey" @click="uploadDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="uploading"
            :disabled="!uploadFile"
            @click="confirmUpload"
          >
            Upload
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar.show" :timeout="5000" :color="snackbar.color">
      {{ snackbar.message }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted, h, markRaw, computed } from 'vue'
import { useForm } from 'vee-validate'
import * as yup from 'yup'
import BaseCard from '@/renderer/components/BaseCard.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

interface CommissionBand {
  id?: number
  channel: string
  holder_name: string
  lower_bound: number
  upper_bound: number | null
  maximum_commission: number
  applicable_rate: number
  creation_date?: string
  created_by?: string
}

interface BrokerRecord {
  id: number
  name: string
}

// ---------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------
const channels = [
  { title: 'Broker', value: 'broker' },
  { title: 'Binder', value: 'binder' },
  { title: 'Tied Agent', value: 'tied_agent' }
]

const activeChannel = ref<string>('broker')
const bands = ref<CommissionBand[]>([])

// Holder dropdown state — populated from the Broker metadata list.
// Empty string means "Default (applies to all holders on this channel)".
const brokers = ref<BrokerRecord[]>([])
const selectedHolder = ref<string>('')
const holderOptions = computed<Array<{ title: string; value: string }>>(() => [
  { title: 'Default (applies to all holders)', value: '' },
  ...brokers.value.map((b) => ({ title: b.name, value: b.name }))
])

const formDialog = ref(false)
const isEditMode = ref(false)
const editingId = ref<number | null>(null)
const formChannel = ref<string>('broker')
const formHolder = ref<string>('')
const unboundedUpper = ref(false)
const saving = ref(false)

const deleteDialog = ref(false)
const pendingDelete = ref<CommissionBand | null>(null)
const deleting = ref(false)

const uploadDialog = ref(false)
const uploadFile = ref<File | File[] | null>(null)
const uploading = ref(false)

const snackbar = ref<{ show: boolean; message: string; color: string }>({
  show: false,
  message: '',
  color: 'success'
})

// ---------------------------------------------------------------------------
// Form
// ---------------------------------------------------------------------------
const validationSchema = yup.object({
  lower_bound: yup
    .number()
    .typeError('Lower bound must be a number')
    .required('Lower bound is required')
    .min(0, 'Lower bound must be 0 or greater'),
  upper_bound: yup
    .number()
    .nullable()
    .transform((v, o) => (o === '' || o === null ? null : v))
    .typeError('Upper bound must be a number')
    .when('lower_bound', (lowerBound: any, schema) => {
      const lb = Array.isArray(lowerBound) ? lowerBound[0] : lowerBound
      return schema.test(
        'gt-lower',
        'Upper bound must be greater than lower bound',
        function (value) {
          if (unboundedUpper.value) return true
          if (value === null || value === undefined) return false
          return Number(value) > Number(lb ?? 0)
        }
      )
    }),
  maximum_commission: yup
    .number()
    .typeError('Maximum commission must be a number')
    .required('Maximum commission is required')
    .min(0, 'Maximum commission must be 0 or greater')
    .max(100, 'Maximum commission must be 100 or less'),
  applicable_rate: yup
    .number()
    .typeError('Applicable rate must be a number')
    .required('Applicable rate is required')
    .min(0, 'Applicable rate must be 0 or greater')
    .max(100, 'Applicable rate must be 100 or less')
    .test(
      'le-max',
      'Applicable rate cannot exceed the maximum commission',
      function (value) {
        const max = this.parent.maximum_commission
        if (value === null || value === undefined || max === null) return true
        return Number(value) <= Number(max)
      }
    )
})

const { defineField, errors, resetForm, validate } = useForm({
  validationSchema,
  initialValues: {
    lower_bound: 0,
    upper_bound: 0,
    maximum_commission: 0,
    applicable_rate: 0
  }
})

const [lowerBound, lowerBoundAttrs] = defineField('lower_bound')
const [upperBound, upperBoundAttrs] = defineField('upper_bound')
const [maximumCommission, maximumCommissionAttrs] =
  defineField('maximum_commission')
const [applicableRate, applicableRateAttrs] = defineField('applicable_rate')

// ---------------------------------------------------------------------------
// Accounting-format helpers for the bound inputs. The underlying refs still
// hold plain numbers (so yup + the save handler are unaffected) — these
// wrappers only change how values are rendered into the <v-text-field>.
// Pattern mirrors `formattedFreeCoverLimit` in Generalnput.vue.
// ---------------------------------------------------------------------------
function formatAccountingInput(v: unknown): string {
  if (v === null || v === undefined || v === '') return ''
  const n = typeof v === 'number' ? v : Number(String(v).replace(/,/g, ''))
  if (!Number.isFinite(n)) return String(v)
  return n.toLocaleString(undefined, { maximumFractionDigits: 2 })
}

function parseAccountingInput(v: unknown): number {
  if (v === null || v === undefined || v === '') return 0
  const cleaned = String(v).replace(/,/g, '').trim()
  const n = parseFloat(cleaned)
  return Number.isNaN(n) ? 0 : n
}

const formattedLowerBound = computed({
  get: () => formatAccountingInput(lowerBound.value),
  set: (val: string | number | null) => {
    lowerBound.value = parseAccountingInput(val)
  }
})

const formattedUpperBound = computed({
  get: () => formatAccountingInput(upperBound.value),
  set: (val: string | number | null) => {
    upperBound.value = parseAccountingInput(val)
  }
})

// ---------------------------------------------------------------------------
// Percent ↔ decimal conversion (rates only)
// ---------------------------------------------------------------------------
function decimalToPercent(value: unknown): number {
  const n = Number(value)
  if (!Number.isFinite(n)) return 0
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
function formatCurrency(value: unknown): string {
  if (value === null || value === undefined || value === '') return ''
  return (
    'R ' +
    Number(value).toLocaleString(undefined, {
      minimumFractionDigits: 0,
      maximumFractionDigits: 2
    })
  )
}

function formatBound(params: any) {
  return formatCurrency(params.value)
}

function formatUpperBound(params: any) {
  const v = params.value
  if (v === null || v === undefined || v === '') return 'Unbounded'
  return formatCurrency(v)
}

function percentFormatter(params: any) {
  const v = params.value
  if (v === null || v === undefined || v === '') return ''
  return (
    decimalToPercent(v).toLocaleString(undefined, {
      minimumFractionDigits: 2,
      maximumFractionDigits: 4
    }) + '%'
  )
}

function dateFormatter(params: any) {
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

const ActionsCell = markRaw({
  props: ['params'],
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

function holderCellFormatter(params: any) {
  const v = params.value
  if (v === null || v === undefined || v === '') return 'Default'
  return String(v)
}

const columnDefs = [
  {
    field: 'holder_name',
    headerName: 'Holder',
    sortable: true,
    filter: true,
    valueFormatter: holderCellFormatter,
    flex: 1
  },
  {
    field: 'lower_bound',
    headerName: 'Lower Bound',
    sortable: true,
    valueFormatter: formatBound,
    flex: 1
  },
  {
    field: 'upper_bound',
    headerName: 'Upper Bound',
    sortable: true,
    valueFormatter: formatUpperBound,
    flex: 1
  },
  {
    field: 'maximum_commission',
    headerName: 'Maximum Commission %',
    sortable: true,
    type: 'rightAligned',
    valueFormatter: percentFormatter,
    flex: 1
  },
  {
    field: 'applicable_rate',
    headerName: 'Applicable Rate %',
    sortable: true,
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
// Client-side contiguity check (server still validates authoritatively)
// ---------------------------------------------------------------------------
const contiguityIssues = computed<string[]>(() => {
  const issues: string[] = []
  const rows = [...bands.value].sort((a, b) => a.lower_bound - b.lower_bound)
  if (rows.length === 0) return issues

  for (let i = 0; i < rows.length; i++) {
    const b = rows[i]
    if (b.applicable_rate > b.maximum_commission) {
      issues.push(
        `Band starting at ${formatCurrency(b.lower_bound)}: applicable rate exceeds the maximum commission.`
      )
    }
    if (i === 0 && b.lower_bound !== 0) {
      issues.push(
        `First band does not start at R 0 (starts at ${formatCurrency(b.lower_bound)}).`
      )
    }
    if (i > 0) {
      const prev = rows[i - 1]
      if (prev.upper_bound === null) {
        issues.push(
          `There is a band after an unbounded band — the unbounded band must be last.`
        )
      } else if (prev.upper_bound !== b.lower_bound) {
        issues.push(
          `Gap or overlap between ${formatCurrency(prev.upper_bound)} and ${formatCurrency(b.lower_bound)}.`
        )
      }
    }
  }
  const last = rows[rows.length - 1]
  if (last.upper_bound !== null) {
    issues.push(
      `Last band has a finite upper bound (${formatCurrency(last.upper_bound)}) — consider adding a final unbounded band.`
    )
  }
  return issues
})

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------
function channelLabel(v?: string | null): string {
  if (!v) return ''
  return channels.find((c) => c.value === v)?.title ?? v
}

async function onChannelChange() {
  // Reset holder filter to Default when switching channel so the user
  // sees that channel's default scale first.
  selectedHolder.value = ''
  await loadBands()
}

function openAddDialog() {
  isEditMode.value = false
  editingId.value = null
  formChannel.value = activeChannel.value
  formHolder.value = selectedHolder.value
  unboundedUpper.value = false
  // Default new lower_bound to the current last upper_bound, if any.
  const rows = [...bands.value].sort((a, b) => a.lower_bound - b.lower_bound)
  const last = rows[rows.length - 1]
  const nextLower = last && last.upper_bound !== null ? last.upper_bound : 0
  resetForm({
    values: {
      lower_bound: nextLower,
      upper_bound: 0,
      maximum_commission: 0,
      applicable_rate: 0
    }
  })
  formDialog.value = true
}

function openEditDialog(row: CommissionBand) {
  isEditMode.value = true
  editingId.value = row.id ?? null
  formChannel.value = row.channel
  formHolder.value = row.holder_name ?? ''
  unboundedUpper.value = row.upper_bound === null
  resetForm({
    values: {
      lower_bound: row.lower_bound,
      upper_bound: row.upper_bound ?? 0,
      maximum_commission: decimalToPercent(row.maximum_commission),
      applicable_rate: decimalToPercent(row.applicable_rate)
    }
  })
  formDialog.value = true
}

function closeFormDialog() {
  formDialog.value = false
}

async function saveBand() {
  const result = await validate()
  if (!result.valid) return

  const payload: CommissionBand = {
    channel: formChannel.value,
    holder_name: (formHolder.value ?? '').trim(),
    lower_bound: Number(lowerBound.value) || 0,
    upper_bound: unboundedUpper.value ? null : Number(upperBound.value) || 0,
    maximum_commission: percentToDecimal(maximumCommission.value),
    applicable_rate: percentToDecimal(applicableRate.value)
  }

  saving.value = true
  try {
    if (isEditMode.value && editingId.value !== null) {
      await GroupPricingService.updateCommissionBand(editingId.value, payload)
      showMessage('Band updated.', 'success')
    } else {
      await GroupPricingService.createCommissionBand(payload)
      showMessage('Band created.', 'success')
    }
    formDialog.value = false
    await loadBands()
  } catch (err: any) {
    const status = err?.response?.status
    const msg =
      err?.response?.data?.error ||
      err?.response?.data ||
      'Failed to save band.'
    if (status === 409) {
      showMessage('A band with this lower bound already exists.', 'error')
    } else {
      showMessage(msg, 'error')
    }
  } finally {
    saving.value = false
  }
}

function requestDelete(row: CommissionBand) {
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
    await GroupPricingService.deleteCommissionBand(pendingDelete.value.id)
    showMessage('Band deleted.', 'success')
    deleteDialog.value = false
    pendingDelete.value = null
    await loadBands()
  } catch (err: any) {
    showMessage(
      err?.response?.data?.error ||
        err?.response?.data ||
        'Failed to delete band.',
      'error'
    )
  } finally {
    deleting.value = false
  }
}

function downloadTemplate() {
  const headers = [
    'channel',
    'holder_name',
    'lower_bound',
    'upper_bound',
    'maximum_commission',
    'applicable_rate'
  ]
  // Illustrative sample: the first 5 broker rows leave holder_name blank
  // to define the default FSCA scale for all brokers. The "ABC Brokers"
  // block shows how to override for a specific holder. Binder/tied_agent
  // defaults follow. Empty upper_bound on the last row in a group means
  // "unbounded".
  const rows: Array<Array<string | number>> = [
    ['broker', '', 0, 200000, 0.075, 0.075],
    ['broker', '', 200000, 300000, 0.05, 0.05],
    ['broker', '', 300000, 600000, 0.03, 0.03],
    ['broker', '', 600000, 2000000, 0.02, 0.02],
    ['broker', '', 2000000, '', 0.01, 0.01],
    ['broker', 'ABC Brokers', 0, 500000, 0.06, 0.06],
    ['broker', 'ABC Brokers', 500000, '', 0.03, 0.03],
    ['binder', '', 0, '', 0.05, 0.05],
    ['tied_agent', '', 0, '', 0.05, 0.05]
  ]

  const escape = (v: string | number) => {
    const s = String(v ?? '')
    return /[",\n\r]/.test(s) ? `"${s.replace(/"/g, '""')}"` : s
  }

  const csv =
    [headers, ...rows]
      .map((row) => row.map((cell) => escape(cell as any)).join(','))
      .join('\n') + '\n'

  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = 'commission-structure-template.csv'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

async function confirmUpload() {
  const file =
    uploadFile.value instanceof File
      ? uploadFile.value
      : Array.isArray(uploadFile.value)
        ? uploadFile.value[0]
        : null
  if (!file) return
  uploading.value = true
  try {
    const form = new FormData()
    form.append('file', file)
    form.append('table_type', 'Commission Structure')
    form.append('risk_rate_code', '')
    await GroupPricingService.uploadTables(form)
    showMessage('Commission structure uploaded.', 'success')
    uploadDialog.value = false
    uploadFile.value = null
    await loadBands()
  } catch (err: any) {
    showMessage(
      err?.response?.data?.error || err?.response?.data || 'Upload failed.',
      'error'
    )
  } finally {
    uploading.value = false
  }
}

function showMessage(message: string, color: 'success' | 'error') {
  snackbar.value = { show: true, message, color }
}

// ---------------------------------------------------------------------------
// Data loading
// ---------------------------------------------------------------------------
async function loadBands() {
  try {
    const response = await GroupPricingService.getCommissionBands(
      activeChannel.value,
      selectedHolder.value
    )
    bands.value = (response.data ?? []) as CommissionBand[]
  } catch (err: any) {
    showMessage(
      err?.response?.data?.error || 'Failed to load commission bands.',
      'error'
    )
  }
}

async function loadBrokers() {
  try {
    const response = await GroupPricingService.getBrokers()
    brokers.value = (response.data ?? []) as BrokerRecord[]
  } catch {
    brokers.value = []
  }
}

onMounted(async () => {
  await Promise.all([loadBands(), loadBrokers()])
})
</script>

<style>
.ag-center-header .ag-header-cell-label {
  justify-content: center;
}
</style>
