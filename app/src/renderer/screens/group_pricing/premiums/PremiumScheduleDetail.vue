<template>
  <v-container fluid>
    <base-card :show-actions="false">
      <template #header>
        <h3 class="mb-0">Premium Schedule Detail</h3>
      </template>
      <template #default>
        <page-header
          title="Schedule Detail"
          :subtitle="
            detail
              ? `${detail.scheme_name} — ${fmtDate(detail.period_start ?? `${detail.year}-${String(detail.month).padStart(2, '0')}-01`)}`
              : ''
          "
          icon="mdi-calendar-check"
          :breadcrumbs="[
            {
              text: 'Premium Schedules',
              to: { name: 'group-pricing-premium-schedules' }
            },
            { text: `Schedule ${scheduleId}` }
          ]"
        />

        <!-- Linked Submission chip -->
        <v-row v-if="detail?.linked_submission_id" class="mb-2" align="center">
          <v-col cols="auto">
            <v-chip
              color="deep-purple"
              variant="tonal"
              size="small"
              prepend-icon="mdi-inbox-arrow-down"
            >
              Source: Employer Submission #{{ detail.linked_submission_id }}
            </v-chip>
            <v-btn
              size="x-small"
              variant="text"
              color="deep-purple"
              class="ml-1"
              append-icon="mdi-arrow-right"
              @click="
                $router.push({
                  name: 'group-pricing-bordereaux-inbound-detail',
                  params: { submissionId: detail.linked_submission_id }
                })
              "
            >
              View
            </v-btn>
          </v-col>
        </v-row>

        <!-- KPI chips -->
        <v-row class="mb-3">
          <v-col v-for="chip in kpiChips" :key="chip.label" cols="auto">
            <v-chip color="primary" variant="tonal" size="small">
              <strong class="mr-1">{{ chip.label }}:</strong> {{ chip.value }}
            </v-chip>
          </v-col>
        </v-row>
        <!-- Audit trail card -->
        <v-row v-if="hasAuditData" class="mb-3">
          <v-col cols="12">
            <v-card variant="outlined" density="compact">
              <v-table density="compact">
                <thead>
                  <tr>
                    <th>Stage</th>
                    <th>User</th>
                    <th>Date</th>
                    <th>Note</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-if="detail.reviewed_by">
                    <td>Reviewed</td>
                    <td>{{ detail.reviewed_by }}</td>
                    <td>{{ fmtDate(detail.reviewed_at) }}</td>
                    <td>—</td>
                  </tr>
                  <tr v-if="detail.approved_by">
                    <td>Approved</td>
                    <td>{{ detail.approved_by }}</td>
                    <td>{{ fmtDate(detail.approved_at) }}</td>
                    <td>—</td>
                  </tr>
                  <tr v-if="detail.finalized_by">
                    <td>Finalized</td>
                    <td>{{ detail.finalized_by }}</td>
                    <td>{{ fmtDate(detail.finalized_at) }}</td>
                    <td>—</td>
                  </tr>
                  <tr v-if="detail.voided_by">
                    <td>{{
                      detail.status === 'cancelled' ? 'Cancelled' : 'Voided'
                    }}</td>
                    <td>{{ detail.voided_by }}</td>
                    <td>{{ fmtDate(detail.voided_at) }}</td>
                    <td>{{ detail.void_reason ?? '—' }}</td>
                  </tr>
                </tbody>
              </v-table>
            </v-card>
          </v-col>
        </v-row>
        <!-- Terminal state banner -->
        <v-row
          v-if="detail?.status === 'void' || detail?.status === 'cancelled'"
          class="mb-3"
        >
          <v-col cols="12">
            <v-alert
              :type="detail.status === 'void' ? 'error' : 'warning'"
              variant="tonal"
            >
              This schedule has been
              <strong>{{
                detail.status === 'void' ? 'voided' : 'cancelled'
              }}</strong
              >.
              <span v-if="detail.void_reason">
                Reason: {{ detail.void_reason }}</span
              >
            </v-alert>
          </v-col>
        </v-row>
        <!-- Member rows grid -->
        <v-row class="mb-4">
          <v-col cols="12">
            <v-card variant="outlined">
              <v-card-title
                class="text-subtitle-1 pa-3 d-flex justify-space-between align-center"
              >
                <span>Member Premium Breakdown</span>
                <div class="d-flex ga-2 align-center">
                  <v-btn
                    variant="outlined"
                    size="small"
                    prepend-icon="mdi-download"
                    :loading="exporting"
                    @click="exportCSV"
                  >
                    Export CSV
                  </v-btn>

                  <!-- Draft actions -->
                  <template v-if="detail?.status === 'draft'">
                    <v-btn
                      variant="outlined"
                      size="small"
                      color="secondary"
                      prepend-icon="mdi-refresh"
                      :loading="regenerating"
                      @click="regenerateDialog = true"
                    >
                      Regenerate
                    </v-btn>
                    <v-btn
                      color="info"
                      size="small"
                      prepend-icon="mdi-send"
                      :loading="transitioning"
                      @click="handleReview"
                    >
                      Submit for Review
                    </v-btn>
                    <v-btn
                      color="error"
                      size="small"
                      variant="outlined"
                      prepend-icon="mdi-cancel"
                      @click="confirmCancelDialog = true"
                    >
                      Cancel Schedule
                    </v-btn>
                  </template>

                  <!-- Reviewed actions -->
                  <template v-else-if="detail?.status === 'reviewed'">
                    <v-btn
                      color="success"
                      size="small"
                      prepend-icon="mdi-check"
                      :loading="transitioning"
                      @click="handleApprove"
                    >
                      Approve
                    </v-btn>
                    <v-btn
                      variant="outlined"
                      size="small"
                      prepend-icon="mdi-undo"
                      :loading="transitioning"
                      @click="handleReturnToDraft"
                    >
                      Return to Draft
                    </v-btn>
                    <v-btn
                      color="error"
                      size="small"
                      variant="outlined"
                      prepend-icon="mdi-cancel"
                      @click="confirmCancelDialog = true"
                    >
                      Cancel Schedule
                    </v-btn>
                  </template>

                  <!-- Approved actions -->
                  <template v-else-if="detail?.status === 'approved'">
                    <v-btn
                      color="primary"
                      size="small"
                      prepend-icon="mdi-check-circle"
                      :loading="transitioning"
                      @click="finalizeDialog = true"
                    >
                      Finalize
                    </v-btn>
                    <v-btn
                      variant="outlined"
                      size="small"
                      prepend-icon="mdi-undo"
                      :loading="transitioning"
                      @click="handleReturnToDraft"
                    >
                      Return to Draft
                    </v-btn>
                    <v-btn
                      color="error"
                      size="small"
                      variant="outlined"
                      prepend-icon="mdi-cancel"
                      @click="confirmCancelDialog = true"
                    >
                      Cancel Schedule
                    </v-btn>
                  </template>

                  <!-- Finalized actions -->
                  <template v-else-if="detail?.status === 'finalized'">
                    <v-btn
                      color="primary"
                      size="small"
                      prepend-icon="mdi-file-document"
                      :loading="transitioning"
                      @click="openInvoiceDialog"
                    >
                      Generate Invoice
                    </v-btn>
                    <v-btn
                      color="error"
                      size="small"
                      variant="outlined"
                      prepend-icon="mdi-cancel"
                      @click="confirmVoidDialog = true"
                    >
                      Void Schedule
                    </v-btn>
                  </template>

                  <!-- Invoiced actions -->
                  <template v-else-if="detail?.status === 'invoiced'">
                    <v-btn
                      color="error"
                      size="small"
                      variant="outlined"
                      prepend-icon="mdi-cancel"
                      @click="confirmVoidDialog = true"
                    >
                      Void Schedule
                    </v-btn>
                  </template>
                </div>
              </v-card-title>
              <v-card-text class="pa-0">
                <ag-grid-vue
                  class="ag-theme-balham"
                  style="height: 480px; width: 100%"
                  :column-defs="columnDefs"
                  :row-data="detail?.members ?? []"
                  :default-col-def="defaultColDef"
                  :loading="loading"
                  :get-row-class="getRowClass"
                  :stop-editing-when-cells-loses-focus="true"
                  @cell-value-changed="onCellValueChanged"
                />
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </template>
    </base-card>

    <!-- Finalize confirm dialog -->
    <v-dialog v-model="finalizeDialog" max-width="440" persistent>
      <v-card>
        <v-card-title>Finalize Schedule</v-card-title>
        <v-card-text>
          <v-alert type="warning" variant="tonal">
            Finalizing will lock the schedule. You can then generate an invoice
            separately. This action cannot be undone without voiding.
          </v-alert>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="finalizeDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="transitioning"
            @click="handleFinalize"
            >Finalize</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Generate Invoice dialog with due date picker -->
    <v-dialog v-model="invoiceDialog" max-width="480" persistent>
      <v-card>
        <v-card-title>Generate Invoice</v-card-title>
        <v-card-text>
          <v-alert type="info" variant="tonal" class="mb-4">
            An invoice will be generated for {{ detail?.scheme_name }} —
            {{ monthName(detail?.month) }} {{ detail?.year }}.
          </v-alert>
          <v-text-field
            v-model="invoiceDueDate"
            label="Due Date"
            type="date"
            variant="outlined"
            hint="Defaults to the last day of the invoice period month"
            persistent-hint
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="invoiceDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="transitioning"
            @click="handleGenerateInvoice"
          >
            Generate
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Void confirm dialog -->
    <v-dialog v-model="confirmVoidDialog" max-width="480" persistent>
      <v-card>
        <v-card-title>Void Schedule</v-card-title>
        <v-card-text>
          <v-alert type="error" variant="tonal" class="mb-4">
            Voiding is irreversible. Please provide a mandatory reason.
          </v-alert>
          <v-textarea
            v-model="voidReason"
            label="Reason *"
            variant="outlined"
            rows="3"
            auto-grow
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="confirmVoidDialog = false"
            >Cancel</v-btn
          >
          <v-btn
            color="error"
            :disabled="!voidReason.trim()"
            :loading="transitioning"
            @click="handleVoid"
          >
            Void
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Cancel confirm dialog -->
    <v-dialog v-model="confirmCancelDialog" max-width="480" persistent>
      <v-card>
        <v-card-title>Cancel Schedule</v-card-title>
        <v-card-text>
          <v-alert type="warning" variant="tonal" class="mb-4">
            Cancelling will abandon this schedule. Please provide a reason.
          </v-alert>
          <v-textarea
            v-model="cancelReason"
            label="Reason *"
            variant="outlined"
            rows="3"
            auto-grow
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="confirmCancelDialog = false"
            >Cancel</v-btn
          >
          <v-btn
            color="warning"
            :disabled="!cancelReason.trim()"
            :loading="transitioning"
            @click="handleCancel"
          >
            Cancel Schedule
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Regenerate confirm dialog -->
    <v-dialog v-model="regenerateDialog" max-width="480" persistent>
      <v-card>
        <v-card-title>Regenerate Schedule</v-card-title>
        <v-card-text>
          <v-alert type="info" variant="tonal">
            This will replace all current member rows with a fresh calculation
            based on the latest member data and premium rates. Any manual rate
            edits will be lost.
          </v-alert>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="regenerateDialog = false"
            >Cancel</v-btn
          >
          <v-btn
            color="secondary"
            :loading="regenerating"
            @click="handleRegenerate"
          >
            Regenerate
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Remove row confirm dialog -->
    <v-dialog v-model="removeRowDialog" max-width="400" persistent>
      <v-card>
        <v-card-title>Remove Member</v-card-title>
        <v-card-text
          >Are you sure you want to remove this member from the
          schedule?</v-card-text
        >
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="removeRowDialog = false">Cancel</v-btn>
          <v-btn color="error" :loading="transitioning" @click="handleRemoveRow"
            >Remove</v-btn
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
import { ref, computed, onMounted } from 'vue'
import { onBeforeRouteLeave } from 'vue-router'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-balham.css'
import { AgGridVue } from 'ag-grid-vue3'
import PremiumManagementService from '@/renderer/api/PremiumManagementService'
import BaseCard from '@/renderer/components/BaseCard.vue'
import PageHeader from '@/renderer/components/PageHeader.vue'
import { fmtDate, fmtCurrency } from '@/renderer/utils/formatters'

const props = defineProps<{ scheduleId: string }>()

const loading = ref(false)
const transitioning = ref(false)
const regenerating = ref(false)
const exporting = ref(false)
const finalizeDialog = ref(false)
const confirmVoidDialog = ref(false)
const confirmCancelDialog = ref(false)
const regenerateDialog = ref(false)
const removeRowDialog = ref(false)
const invoiceDialog = ref(false)
const invoiceDueDate = ref('')
const snackbar = ref(false)
const snackbarText = ref('')
const snackbarColor = ref('success')

const detail = ref<any>(null)
const selectedRemoveRowId = ref<number | null>(null)
const voidReason = ref('')
const cancelReason = ref('')
const isDirty = ref(false)

onBeforeRouteLeave((_to, _from, next) => {
  if (!isDirty.value) return next()
  const ok = confirm('You have unsaved edits. Leave anyway?')
  ok ? next() : next(false)
})

const MONTH_NAMES = [
  'Jan',
  'Feb',
  'Mar',
  'Apr',
  'May',
  'Jun',
  'Jul',
  'Aug',
  'Sep',
  'Oct',
  'Nov',
  'Dec'
]

function monthName(m: number | undefined) {
  return m ? MONTH_NAMES[m - 1] : ''
}

function lastDayOfMonth(year: number, month: number): string {
  // month is 1-based; Date with day=0 of next month gives last day of target month
  const d = new Date(year, month, 0)
  const mm = String(d.getMonth() + 1).padStart(2, '0')
  const dd = String(d.getDate()).padStart(2, '0')
  return `${d.getFullYear()}-${mm}-${dd}`
}

const hasAuditData = computed(() => {
  if (!detail.value) return false
  return (
    detail.value.reviewed_by ||
    detail.value.approved_by ||
    detail.value.finalized_by ||
    detail.value.voided_by
  )
})

const kpiChips = computed(() => {
  if (!detail.value) return []
  return [
    { label: 'Members', value: detail.value.member_count },
    { label: 'New Joiners', value: detail.value.new_joiners },
    { label: 'Exits', value: detail.value.exits },
    { label: 'Gross Premium', value: fmtCurrency(detail.value.gross_premium) },
    { label: 'Net Payable', value: fmtCurrency(detail.value.net_payable) }
  ]
})

const defaultColDef = { sortable: true, filter: true, resizable: true }

const columnDefs = computed(() => {
  const isDraft = detail.value?.status === 'draft'
  const cols: any[] = [
    { headerName: 'Member ID', field: 'member_id', maxWidth: 110 },
    { headerName: 'Name', field: 'member_name', minWidth: 180 },
    { headerName: 'ID Number', field: 'member_id_number', minWidth: 140 },
    { headerName: 'Benefit', field: 'benefit' },
    {
      headerName: 'Annual Salary',
      field: 'annual_salary',
      valueFormatter: (p: any) => fmtCurrency(p.value)
    },
    {
      headerName: 'Full Month Prem',
      field: 'full_month_premium',
      valueFormatter: (p: any) => fmtCurrency(p.value)
    },
    { headerName: 'Pro-Rata Days', field: 'pro_rata_days', maxWidth: 130 },
    {
      headerName: 'Rate',
      field: 'rate',
      maxWidth: 110,
      editable: isDraft,
      valueFormatter: (p: any) => (p.value ? `${p.value}%` : '—'),
      cellClass: isDraft ? 'editable-cell' : ''
    },
    {
      headerName: 'Actual Premium',
      field: 'actual_premium',
      editable: isDraft,
      valueFormatter: (p: any) => fmtCurrency(p.value),
      cellClass: isDraft ? 'editable-cell' : ''
    },
    {
      headerName: 'Employer',
      field: 'employer_contribution',
      valueFormatter: (p: any) => fmtCurrency(p.value)
    },
    {
      headerName: 'Employee',
      field: 'employee_contribution',
      valueFormatter: (p: any) => fmtCurrency(p.value)
    }
  ]
  // if (isDraft) {
  //   cols.push({
  //     headerName: '',
  //     field: 'id',
  //     maxWidth: 90,
  //     sortable: false,
  //     filter: false,
  //     cellRenderer: () =>
  //       `<button style="color:#F44336;background:none;border:none;cursor:pointer;font-size:12px">Remove</button>`,
  //     onCellClicked: (p: any) => {
  //       selectedRemoveRowId.value = p.data.id
  //       removeRowDialog.value = true
  //     }
  //   })
  // }
  return cols
})

function getRowClass(params: any) {
  return params.data?.is_pro_rata ? 'ag-row-pro-rata' : ''
}

function showSuccess(msg: string) {
  snackbarText.value = msg
  snackbarColor.value = 'success'
  snackbar.value = true
}

function showError(e: any, fallback: string) {
  snackbarText.value = e?.response?.data?.message ?? fallback
  snackbarColor.value = 'error'
  snackbar.value = true
}

async function loadDetail() {
  loading.value = true
  try {
    const res = await PremiumManagementService.getScheduleDetail(
      Number(props.scheduleId)
    )
    detail.value = res.data.data
  } catch (e) {
    console.error('Failed to load schedule detail', e)
  } finally {
    loading.value = false
  }
}

async function handleReview() {
  transitioning.value = true
  try {
    await PremiumManagementService.reviewSchedule(Number(props.scheduleId))
    showSuccess('Schedule submitted for review')
    await loadDetail()
  } catch (e) {
    showError(e, 'Failed to submit for review')
  } finally {
    transitioning.value = false
  }
}

async function handleApprove() {
  transitioning.value = true
  try {
    await PremiumManagementService.approveSchedule(Number(props.scheduleId))
    showSuccess('Schedule approved')
    await loadDetail()
  } catch (e) {
    showError(e, 'Failed to approve schedule')
  } finally {
    transitioning.value = false
  }
}

async function handleReturnToDraft() {
  transitioning.value = true
  try {
    await PremiumManagementService.returnScheduleToDraft(
      Number(props.scheduleId)
    )
    showSuccess('Schedule returned to draft')
    await loadDetail()
  } catch (e) {
    showError(e, 'Failed to return to draft')
  } finally {
    transitioning.value = false
  }
}

async function handleFinalize() {
  transitioning.value = true
  finalizeDialog.value = false
  try {
    await PremiumManagementService.finalizeSchedule(Number(props.scheduleId))
    showSuccess('Schedule finalized')
    await loadDetail()
  } catch (e) {
    showError(e, 'Failed to finalize schedule')
  } finally {
    transitioning.value = false
  }
}

function openInvoiceDialog() {
  if (detail.value) {
    invoiceDueDate.value = lastDayOfMonth(detail.value.year, detail.value.month)
  }
  invoiceDialog.value = true
}

async function handleGenerateInvoice() {
  transitioning.value = true
  invoiceDialog.value = false
  try {
    await PremiumManagementService.generateInvoice(
      Number(props.scheduleId),
      invoiceDueDate.value || undefined
    )
    showSuccess('Invoice generated')
    await loadDetail()
  } catch (e) {
    showError(e, 'Failed to generate invoice')
  } finally {
    transitioning.value = false
  }
}

async function handleVoid() {
  transitioning.value = true
  try {
    await PremiumManagementService.voidSchedule(Number(props.scheduleId), {
      reason: voidReason.value
    })
    confirmVoidDialog.value = false
    voidReason.value = ''
    showSuccess('Schedule voided')
    await loadDetail()
  } catch (e) {
    showError(e, 'Failed to void schedule')
  } finally {
    transitioning.value = false
  }
}

async function handleCancel() {
  transitioning.value = true
  try {
    await PremiumManagementService.cancelSchedule(Number(props.scheduleId), {
      reason: cancelReason.value
    })
    confirmCancelDialog.value = false
    cancelReason.value = ''
    showSuccess('Schedule cancelled')
    await loadDetail()
  } catch (e) {
    showError(e, 'Failed to cancel schedule')
  } finally {
    transitioning.value = false
  }
}

async function handleRegenerate() {
  regenerating.value = true
  regenerateDialog.value = false
  try {
    await PremiumManagementService.regenerateSchedule(Number(props.scheduleId))
    showSuccess('Schedule regenerated')
    await loadDetail()
  } catch (e) {
    showError(e, 'Failed to regenerate schedule')
  } finally {
    regenerating.value = false
  }
}

async function handleRemoveRow() {
  if (!selectedRemoveRowId.value) return
  transitioning.value = true
  try {
    await PremiumManagementService.removeScheduleMember(
      Number(props.scheduleId),
      selectedRemoveRowId.value
    )
    removeRowDialog.value = false
    selectedRemoveRowId.value = null
    showSuccess('Member removed')
    await loadDetail()
  } catch (e) {
    showError(e, 'Failed to remove member')
  } finally {
    transitioning.value = false
  }
}

async function onCellValueChanged(params: any) {
  if (detail.value?.status !== 'draft') return
  isDirty.value = true
  const rowId = params.data.id
  await PremiumManagementService.updateScheduleMemberRow(
    Number(props.scheduleId),
    rowId,
    {
      rate: params.data.rate,
      actual_premium: params.data.actual_premium,
      employer_contribution: params.data.employer_contribution,
      employee_contribution: params.data.employee_contribution
    }
  )
  isDirty.value = false
  await loadDetail()
}

async function exportCSV() {
  exporting.value = true
  try {
    const res = await PremiumManagementService.exportSchedule(
      Number(props.scheduleId)
    )
    const url = URL.createObjectURL(new Blob([res.data], { type: 'text/csv' }))
    const a = document.createElement('a')
    a.href = url
    a.download = `schedule_${props.scheduleId}.csv`
    a.click()
    URL.revokeObjectURL(url)
  } catch (e) {
    console.error('Export failed', e)
  } finally {
    exporting.value = false
  }
}

onMounted(loadDetail)
</script>

<style>
.editable-cell {
  background: rgba(33, 150, 243, 0.08) !important;
  border-left: 2px solid #2196f3 !important;
}
</style>
