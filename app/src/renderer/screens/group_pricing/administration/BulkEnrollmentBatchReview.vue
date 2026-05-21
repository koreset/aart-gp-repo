<template>
  <v-container fluid>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex justify-space-between align-center w-100">
              <div class="d-flex align-center">
                <v-btn
                  icon="mdi-arrow-left"
                  size="small"
                  variant="text"
                  class="mr-2"
                  @click="goBack"
                />
                <div>
                  <span class="headline">Bulk Enrollment Batch</span>
                  <div class="text-caption breadcrumb-on-dark">
                    Member Management → Pending Approvals → Batch
                    {{ batchId }}
                  </div>
                </div>
              </div>
              <div class="d-flex gap-2">
                <v-btn
                  v-if="canRunExternalCheck"
                  size="small"
                  variant="outlined"
                  color="white"
                  rounded
                  prepend-icon="mdi-check-decagram"
                  :loading="externalChecking"
                  @click="runExternalCheck"
                >
                  Run External CheckID
                </v-btn>
                <v-btn
                  size="small"
                  color="error"
                  variant="outlined"
                  rounded
                  :disabled="!isPending || working"
                  @click="rejectDialog = true"
                >
                  Reject
                </v-btn>
                <v-btn
                  size="small"
                  color="success"
                  rounded
                  :disabled="!canApprove || working"
                  :loading="approving"
                  @click="approve"
                >
                  Approve
                </v-btn>
              </div>
            </div>
          </template>

          <template #default>
            <v-row v-if="batch">
              <v-col cols="12" md="3">
                <v-card variant="outlined" class="pa-3 text-center">
                  <div class="text-caption">Status</div>
                  <div class="text-h6 mt-1">{{ batch.status }}</div>
                </v-card>
              </v-col>
              <v-col cols="6" md="3">
                <v-card variant="outlined" class="pa-3 text-center">
                  <div class="text-caption">Uploaded By</div>
                  <div class="text-body-1 mt-1">{{ batch.uploaded_by }}</div>
                  <div class="text-caption mt-1">{{
                    formatDate(batch.uploaded_at)
                  }}</div>
                </v-card>
              </v-col>
              <v-col cols="6" md="3">
                <v-card variant="outlined" class="pa-3 text-center">
                  <div class="text-caption">Approved By</div>
                  <div class="text-body-1 mt-1">{{
                    batch.approved_by || '—'
                  }}</div>
                  <div class="text-caption mt-1">{{
                    formatDate(batch.approved_at)
                  }}</div>
                </v-card>
              </v-col>
              <v-col cols="6" md="3">
                <v-card variant="outlined" class="pa-3 text-center">
                  <div class="text-caption">Rejected By</div>
                  <div class="text-body-1 mt-1">{{
                    batch.rejected_by || '—'
                  }}</div>
                  <div class="text-caption mt-1">{{
                    formatDate(batch.rejected_at)
                  }}</div>
                </v-card>
              </v-col>
            </v-row>

            <v-row v-if="batch" class="mt-2">
              <v-col cols="6" md="3">
                <v-card variant="outlined" class="pa-3 text-center">
                  <div class="text-h4 text-success">{{
                    batch.valid_count
                  }}</div>
                  <div class="text-caption">Valid</div>
                </v-card>
              </v-col>
              <v-col cols="6" md="3">
                <v-card variant="outlined" class="pa-3 text-center">
                  <div class="text-h4 text-error">{{
                    batch.blocking_count
                  }}</div>
                  <div class="text-caption">Blocking</div>
                </v-card>
              </v-col>
              <v-col cols="6" md="3">
                <v-card variant="outlined" class="pa-3 text-center">
                  <div class="text-h4 text-warning">{{
                    batch.soft_error_count
                  }}</div>
                  <div class="text-caption">Soft</div>
                </v-card>
              </v-col>
              <v-col cols="6" md="3">
                <v-card variant="outlined" class="pa-3 text-center">
                  <div class="text-h4 text-info">{{ batch.member_count }}</div>
                  <div class="text-caption">Total</div>
                </v-card>
              </v-col>
            </v-row>

            <v-alert
              v-if="batch && batch.blocking_count > 0"
              class="my-4"
              type="error"
              variant="tonal"
              density="compact"
            >
              This batch cannot be approved while it has blocking validation
              errors. Reject and ask the uploader to fix the source file.
            </v-alert>

            <v-row class="mb-2" align="center">
              <v-col cols="auto">
                <v-btn-toggle
                  v-model="filterMode"
                  color="primary"
                  density="compact"
                  variant="outlined"
                  mandatory
                >
                  <v-btn value="all">All ({{ members.length }})</v-btn>
                  <v-btn value="errors"
                    >With errors ({{ erroredMembers.length }})</v-btn
                  >
                  <v-btn value="valid"
                    >Valid only ({{ validMembers.length }})</v-btn
                  >
                </v-btn-toggle>
              </v-col>
            </v-row>

            <div :style="{ height: gridHeight, width: '100%' }">
              <data-grid
                :column-defs="memberColumnDefs"
                :row-data="filteredMembers"
                :pagination="false"
                :loading="loading"
                style="height: 100%; width: 100%"
              />
            </div>

            <MemberUploadErrorReport
              v-if="report.length > 0"
              class="mt-4"
              :blocking-errors="blockingRows"
              :soft-errors="softRows"
              :context-label="`Batch ${batchId}`"
              filename-prefix="batch-validation"
            />
          </template>
        </base-card>
      </v-col>
    </v-row>

    <v-dialog v-model="rejectDialog" max-width="600px" persistent>
      <base-card>
        <template #header>
          <span class="headline">Reject Batch</span>
        </template>
        <template #default>
          <v-textarea
            v-model="rejectReason"
            label="Reason"
            variant="outlined"
            density="compact"
            rows="4"
            :rules="[(v) => !!v?.trim() || 'Reason is required']"
            autofocus
          />
          <div class="d-flex justify-end gap-2 mt-2">
            <v-btn
              size="small"
              variant="outlined"
              rounded
              :disabled="rejecting"
              @click="rejectDialog = false"
              >Cancel</v-btn
            >
            <v-btn
              size="small"
              color="error"
              rounded
              :loading="rejecting"
              :disabled="!rejectReason?.trim()"
              @click="reject"
              >Confirm Reject</v-btn
            >
          </div>
        </template>
      </base-card>
    </v-dialog>

    <v-snackbar v-model="snackbar" :timeout="4000" :color="snackbarColor">{{
      snackbarMessage
    }}</v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import MemberUploadErrorReport from '@/renderer/screens/group_pricing/shared/MemberUploadErrorReport.vue'
import { useGridHeight } from '@/renderer/composables/useGridHeight'

interface Props {
  batchId: string
}
const props = defineProps<Props>()

const router = useRouter()
const gridHeight = useGridHeight(380)

const batch = ref<any>(null)
const members = ref<any[]>([])
const report = ref<any[]>([])
const loading = ref(false)
const approving = ref(false)
const rejecting = ref(false)
const externalChecking = ref(false)
const filterMode = ref<'all' | 'errors' | 'valid'>('all')
const rejectDialog = ref(false)
const rejectReason = ref('')
const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref<'success' | 'error' | 'info'>('success')

const isPending = computed(() => batch.value?.status === 'pending_approval')
const canApprove = computed(
  () => isPending.value && (batch.value?.blocking_count ?? 0) === 0
)
const canRunExternalCheck = computed(
  () => isPending.value && !batch.value?.external_id_check_run
)
const working = computed(
  () => approving.value || rejecting.value || externalChecking.value
)

const errorRowsByRow = computed(() => {
  const map = new Map<number, any[]>()
  for (const r of report.value) {
    if (!map.has(r.row)) map.set(r.row, [])
    map.get(r.row)!.push(r)
  }
  return map
})

const erroredMembers = computed(() =>
  members.value.filter((m) => errorRowsByRow.value.has(m.row_index))
)
const validMembers = computed(() =>
  members.value.filter((m) => !errorRowsByRow.value.has(m.row_index))
)
const filteredMembers = computed(() => {
  if (filterMode.value === 'errors') return erroredMembers.value
  if (filterMode.value === 'valid') return validMembers.value
  return members.value
})

const blockingRows = computed(() =>
  report.value
    .filter((r) => r.severity === 'blocking')
    .map((r) => ({
      row: r.row,
      field: r.field,
      message: r.message,
      member_id_number: r.member_id_number,
      member_name: r.member_name
    }))
)
const softRows = computed(() =>
  report.value
    .filter((r) => r.severity === 'soft')
    .map((r) => ({ row: r.row, message: r.message }))
)

const memberColumnDefs = [
  { headerName: 'Row', field: 'row_index', minWidth: 70, width: 80 },
  {
    headerName: 'Status',
    field: 'validation_status',
    minWidth: 130,
    cellRenderer: (p: any) => {
      const v = p.value || 'valid'
      const colors: Record<string, string> = {
        valid: '#2e7d32',
        soft_error: '#ef6c00',
        blocking_error: '#c62828'
      }
      const c = colors[v] || '#616161'
      return `<span style="background:${c};color:#fff;padding:2px 8px;border-radius:10px;font-size:11px;">${v.replace('_', ' ')}</span>`
    }
  },
  { headerName: 'Member Name', field: 'member_name', minWidth: 200 },
  { headerName: 'ID Number', field: 'member_id_number', minWidth: 150 },
  { headerName: 'ID Type', field: 'member_id_type', minWidth: 110 },
  { headerName: 'Gender', field: 'gender', minWidth: 90 },
  { headerName: 'Annual Salary', field: 'annual_salary', minWidth: 130 },
  { headerName: 'Scheme Category', field: 'scheme_category', minWidth: 160 },
  { headerName: 'Occupation', field: 'occupation', minWidth: 160 },
  {
    headerName: 'Issues',
    field: '_issues',
    minWidth: 280,
    valueGetter: (p: any) => {
      const issues = errorRowsByRow.value.get(p.data?.row_index) ?? []
      return issues
        .map((i: any) => `${i.field || ''}: ${i.message}`)
        .join(' | ')
    }
  }
]

const formatDate = (v: any) => {
  if (!v) return ''
  const d = new Date(v)
  return isNaN(d.getTime()) ? '' : d.toLocaleString()
}

const showSnack = (
  msg: string,
  color: 'success' | 'error' | 'info' = 'success'
) => {
  snackbarMessage.value = msg
  snackbarColor.value = color
  snackbar.value = true
}

const load = async () => {
  loading.value = true
  try {
    const res = await GroupPricingService.getBulkEnrollmentBatch(props.batchId)
    batch.value = res?.data?.batch ?? null
    members.value = res?.data?.members ?? []
    report.value = res?.data?.validation_report ?? []
  } catch (err: any) {
    console.error('Failed to load batch', err)
    showSnack(err?.response?.data?.error || 'Failed to load batch', 'error')
  } finally {
    loading.value = false
  }
}

const approve = async () => {
  approving.value = true
  try {
    const res = await GroupPricingService.approveBulkEnrollmentBatch(
      props.batchId
    )
    batch.value = res?.data?.batch ?? batch.value
    showSnack('Batch approved — members are now active', 'success')
    setTimeout(goBack, 800)
  } catch (err: any) {
    showSnack(err?.response?.data?.error || 'Approval failed', 'error')
  } finally {
    approving.value = false
  }
}

const reject = async () => {
  if (!rejectReason.value?.trim()) return
  rejecting.value = true
  try {
    const res = await GroupPricingService.rejectBulkEnrollmentBatch(
      props.batchId,
      rejectReason.value.trim()
    )
    batch.value = res?.data?.batch ?? batch.value
    rejectDialog.value = false
    showSnack('Batch rejected and drafts removed', 'success')
    setTimeout(goBack, 800)
  } catch (err: any) {
    showSnack(err?.response?.data?.error || 'Rejection failed', 'error')
  } finally {
    rejecting.value = false
  }
}

const runExternalCheck = async () => {
  externalChecking.value = true
  try {
    const res = await GroupPricingService.runExternalIdCheckOnBatch(
      props.batchId
    )
    batch.value = res?.data?.batch ?? batch.value
    report.value = res?.data?.validation_report ?? report.value
    showSnack('External CheckID validation complete', 'success')
    await load()
  } catch (err: any) {
    showSnack(err?.response?.data?.error || 'External check failed', 'error')
  } finally {
    externalChecking.value = false
  }
}

const goBack = () => {
  router.push({ name: 'group-pricing-member-management' })
}

onMounted(load)
</script>

<style scoped>
.headline {
  font-size: 1.25rem;
  font-weight: 500;
}

.breadcrumb-on-dark {
  color: rgba(255, 255, 255, 0.75);
}
</style>
