<template>
  <!-- Top tab strip (Premium-Receipts-style) -->
  <v-tabs
    :model-value="activeTab"
    color="primary"
    density="compact"
    show-arrows
    class="schedule-tabs"
  >
    <v-tab
      value="group-pricing-claim-payment-schedule-claims"
      :to="tabRoute('group-pricing-claim-payment-schedule-claims')"
      prepend-icon="mdi-file-document-multiple"
    >
      Claims
    </v-tab>
    <v-tab
      v-if="hasPermission('claims_pay:generate_acb')"
      value="group-pricing-claim-payment-schedule-acb"
      :to="tabRoute('group-pricing-claim-payment-schedule-acb')"
      prepend-icon="mdi-bank-transfer"
    >
      ACB Files
    </v-tab>
    <v-tab
      value="group-pricing-claim-payment-schedule-queries"
      :to="tabRoute('group-pricing-claim-payment-schedule-queries')"
      prepend-icon="mdi-comment-alert-outline"
    >
      Queries
    </v-tab>
    <v-tab
      v-if="hasPermission('claims_pay:upload_response')"
      value="group-pricing-claim-payment-schedule-reconciliation"
      :to="tabRoute('group-pricing-claim-payment-schedule-reconciliation')"
      prepend-icon="mdi-scale-balance"
    >
      Reconciliation
    </v-tab>
    <v-tab
      v-if="hasPermission('claims_pay:upload_response')"
      value="group-pricing-claim-payment-schedule-proofs"
      :to="tabRoute('group-pricing-claim-payment-schedule-proofs')"
      prepend-icon="mdi-receipt-text-check"
    >
      Proof of Payment
    </v-tab>
  </v-tabs>

  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center flex-wrap gap-2">
              <v-btn
                icon="mdi-arrow-left"
                variant="text"
                class="mr-2"
                @click="goBack"
              />
              <span class="headline">
                Payment Schedule
                <template v-if="schedule">
                  — {{ schedule.schedule_number }}
                </template>
              </span>
            </div>
          </template>

          <template #default>
            <!-- Loading -->
            <v-row v-if="loading && !schedule">
              <v-col cols="12" class="text-center py-8">
                <v-progress-circular indeterminate color="primary" />
                <div class="mt-2 text-body-2 text-medium-emphasis"
                  >Loading payment schedule...</div
                >
              </v-col>
            </v-row>

            <!-- Not found -->
            <empty-state
              v-else-if="!schedule"
              icon="mdi-file-document-remove-outline"
              title="Payment schedule not found"
              message="This schedule may have been deleted or you do not have access."
              action-label="Back to Payment Schedules"
              :action-fn="goBack"
            />

            <template v-else>
              <!-- ── Section: Workflow ───────────────────── -->
              <section class="page-section">
                <div class="section-header">
                  <span class="section-label">Workflow</span>
                  <span class="section-divider" />
                </div>
                <div class="workflow-stepper">
                  <div
                    v-for="(step, idx) in pipelineSteps"
                    :key="step.status"
                    class="workflow-stepper__step"
                    :class="`workflow-stepper__step--${stepTone(idx)}`"
                  >
                    <div class="workflow-stepper__num">
                      <v-icon
                        v-if="idx < currentStepIndex"
                        size="16"
                        icon="mdi-check"
                      />
                      <template v-else>{{ idx + 1 }}</template>
                    </div>
                    <div class="workflow-stepper__text">
                      <div class="workflow-stepper__label">{{
                        step.label
                      }}</div>
                      <div class="workflow-stepper__sub">{{ step.sub }}</div>
                    </div>
                    <v-icon
                      v-if="idx < pipelineSteps.length - 1"
                      class="workflow-stepper__arrow"
                      size="18"
                      >mdi-chevron-right</v-icon
                    >
                  </div>
                </div>
              </section>

              <!-- Conditional notices (locked, duplicates, sanctions) -->
              <v-alert
                v-if="schedule.locked_at"
                type="info"
                variant="tonal"
                density="compact"
                class="mb-4"
                icon="mdi-lock-outline"
              >
                <div class="text-body-2">
                  <strong>Schedule locked.</strong> Line items can only be
                  removed (which sends the claim back to the approval queue for
                  the next cut-off). Amounts and banking details cannot be
                  edited on this schedule.
                </div>
              </v-alert>

              <v-alert
                v-if="outstandingDuplicates > 0"
                type="warning"
                variant="tonal"
                density="compact"
                class="mb-4"
                icon="mdi-account-multiple-outline"
              >
                <div class="text-body-2">
                  <strong>Duplicate beneficiary detected.</strong>
                  {{ outstandingDuplicates }} line(s) share a beneficiary with
                  another line in this schedule. Review each flagged line on the
                  Claims tab and explicitly clear if the duplicate is
                  intentional. First finance authorisation is blocked until
                  every flag is cleared.
                </div>
              </v-alert>

              <v-alert
                v-if="sanctionsBlockers > 0 || reinsuranceOutstanding > 0"
                type="warning"
                variant="tonal"
                density="compact"
                class="mb-4"
                icon="mdi-shield-alert-outline"
              >
                <div class="text-body-2">
                  <span v-if="sanctionsBlockers > 0">
                    <strong>{{ sanctionsBlockers }}</strong> line(s) outstanding
                    on sanctions / PEP screening.
                  </span>
                  <span v-if="reinsuranceOutstanding > 0">
                    <strong>{{ reinsuranceOutstanding }}</strong> reinsurance
                    recovery(ies) not yet raised.
                  </span>
                  These block first finance authorisation.
                </div>
              </v-alert>

              <!-- ── Section: Actions ────────────────────── -->
              <section class="page-section">
                <div class="section-header">
                  <span class="section-label">Actions</span>
                  <span class="section-divider" />
                </div>
                <div class="d-flex align-center flex-wrap ga-3">
                  <!-- Lifecycle gate buttons (Phase 1) -->
                  <v-btn
                    v-if="canSignOff"
                    variant="flat"
                    size="small"
                    rounded
                    color="primary"
                    prepend-icon="mdi-clipboard-check-outline"
                    :loading="signingOff"
                    @click="signOff"
                  >
                    Sign Off (Head of Claims)
                  </v-btn>
                  <v-btn
                    v-if="canStartReview"
                    variant="flat"
                    size="small"
                    rounded
                    color="primary"
                    prepend-icon="mdi-magnify"
                    :loading="startingReview"
                    @click="startFinanceReview"
                  >
                    Start Finance Review
                  </v-btn>
                  <v-btn
                    v-if="canAuthoriseFirst"
                    variant="flat"
                    size="small"
                    rounded
                    color="primary"
                    prepend-icon="mdi-check-decagram"
                    :loading="authorising === 'first'"
                    @click="authoriseFirst"
                  >
                    Authorise (1st)
                  </v-btn>
                  <v-btn
                    v-if="canAuthoriseSecond"
                    variant="flat"
                    size="small"
                    rounded
                    color="primary"
                    prepend-icon="mdi-shield-check"
                    :loading="authorising === 'second'"
                    @click="authoriseSecond"
                  >
                    Authorise (2nd)
                  </v-btn>
                  <v-btn
                    variant="outlined"
                    size="small"
                    rounded
                    prepend-icon="mdi-download"
                    :loading="exporting"
                    @click="exportSchedule"
                  >
                    Export CSV
                  </v-btn>
                  <v-tooltip
                    v-if="hasPermission('claims_pay:generate_acb')"
                    :disabled="!acbBlockedReason"
                    location="bottom"
                  >
                    <template #activator="{ props: tipProps }">
                      <div v-bind="tipProps">
                        <v-btn
                          variant="outlined"
                          size="small"
                          rounded
                          prepend-icon="mdi-file-document-outline"
                          :disabled="!!acbBlockedReason"
                          @click="openACBDialog"
                        >
                          Generate ACB
                        </v-btn>
                      </div>
                    </template>
                    {{ acbBlockedReason }}
                  </v-tooltip>
                  <v-btn
                    v-if="
                      schedule.acb_file_generated &&
                      hasPermission('claims_pay:upload_response')
                    "
                    variant="outlined"
                    size="small"
                    rounded
                    prepend-icon="mdi-file-upload"
                    @click="openResponseDialog"
                  >
                    Upload Bank Response
                  </v-btn>
                  <v-btn
                    v-if="canArchive"
                    variant="outlined"
                    size="small"
                    rounded
                    prepend-icon="mdi-archive-outline"
                    :loading="archiving"
                    @click="archive"
                  >
                    Archive
                  </v-btn>
                  <v-spacer />
                  <v-btn
                    v-if="
                      schedule.status !== 'confirmed' &&
                      hasPermission('claims_pay:upload_response')
                    "
                    variant="flat"
                    size="small"
                    rounded
                    color="success"
                    prepend-icon="mdi-upload"
                    @click="openProofDialog"
                  >
                    Upload Proof of Payment
                  </v-btn>
                </div>
              </section>

              <!-- ── Section: Schedule summary ───────────── -->
              <section class="page-section">
                <div class="section-header">
                  <span class="section-label">Schedule summary</span>
                  <span class="section-divider" />
                </div>
                <v-row dense>
                  <v-col cols="12" sm="6" md="3">
                    <v-card
                      variant="outlined"
                      rounded="lg"
                      class="meta-card h-100 pa-3 d-flex flex-column"
                    >
                      <div class="meta-card__label">Claims</div>
                      <div class="meta-card__value">{{
                        schedule.claims_count
                      }}</div>
                      <div class="meta-card__hint">
                        claim{{ schedule.claims_count === 1 ? '' : 's' }} in
                        this schedule
                      </div>
                    </v-card>
                  </v-col>
                  <v-col cols="12" sm="6" md="3">
                    <v-card
                      variant="outlined"
                      rounded="lg"
                      class="meta-card h-100 pa-3 d-flex flex-column"
                    >
                      <div class="meta-card__label">Total Amount</div>
                      <div class="meta-card__value">{{
                        formatCurrency(schedule.total_amount)
                      }}</div>
                      <div class="meta-card__hint">payment run total</div>
                    </v-card>
                  </v-col>
                  <v-col cols="12" sm="6" md="3">
                    <v-card
                      variant="outlined"
                      rounded="lg"
                      class="meta-card h-100 pa-3 d-flex flex-column"
                    >
                      <div class="meta-card__label">Created By</div>
                      <div class="meta-card__value meta-card__value--text">{{
                        schedule.created_by || '—'
                      }}</div>
                      <div class="meta-card__hint">{{
                        formatDate(schedule.created_at)
                      }}</div>
                    </v-card>
                  </v-col>
                  <v-col cols="12" sm="6" md="3">
                    <v-card
                      variant="outlined"
                      rounded="lg"
                      class="meta-card h-100 pa-3 d-flex flex-column"
                    >
                      <div class="meta-card__label">ACB Generated By</div>
                      <div class="meta-card__value meta-card__value--text">{{
                        schedule.acb_generated_by || '—'
                      }}</div>
                      <div class="meta-card__hint">{{
                        schedule.acb_generated_at
                          ? formatDate(schedule.acb_generated_at)
                          : 'Not generated'
                      }}</div>
                    </v-card>
                  </v-col>
                </v-row>

                <v-card
                  v-if="schedule.description"
                  variant="tonal"
                  color="grey-lighten-4"
                  class="mt-3 pa-3"
                  flat
                >
                  <div class="text-caption text-medium-emphasis"
                    >Description</div
                  >
                  <div class="text-body-2">{{ schedule.description }}</div>
                </v-card>
              </section>

              <!-- Active tab content -->
              <router-view />
            </template>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- ── ACB Generate Dialog ── -->
    <v-dialog v-model="acbDialog" persistent max-width="540px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2">Generate ACB File</v-card-title>
        <v-card-text>
          <v-alert
            type="warning"
            variant="tonal"
            density="compact"
            class="mb-3"
          >
            <div class="text-body-2">
              Generating an ACB file authorises a BankServ payment run for the
              full schedule total. Verify the bank profile, action date, and
              totals below before proceeding.
            </div>
          </v-alert>
          <v-table density="compact" class="border rounded mb-3">
            <tbody>
              <tr>
                <th class="text-left">Schedule</th>
                <td
                  ><strong>{{ schedule?.schedule_number }}</strong></td
                >
              </tr>
              <tr>
                <th class="text-left">Claims</th>
                <td>{{ schedule?.claims_count }}</td>
              </tr>
              <tr>
                <th class="text-left">Total Amount</th>
                <td
                  ><strong>{{
                    formatCurrency(schedule?.total_amount ?? 0)
                  }}</strong></td
                >
              </tr>
            </tbody>
          </v-table>
          <v-select
            v-model="acbProfileId"
            :items="bankProfiles"
            item-title="profile_name"
            item-value="id"
            label="Bank Profile *"
            variant="outlined"
            density="compact"
            class="mb-3"
            :loading="loadingProfiles"
          />
          <v-text-field
            v-model="acbActionDate"
            label="Action Date *"
            type="date"
            variant="outlined"
            density="compact"
            hint="Date the bank should process the payments"
            persistent-hint
            class="mb-3"
          />
          <v-checkbox
            v-model="acbConfirmed"
            density="compact"
            color="teal"
            hide-details
            label="I have verified the bank profile, action date, and schedule totals. Authorise ACB generation."
          />
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="acbDialog = false">Cancel</v-btn>
          <v-btn
            color="teal"
            :loading="generatingACB"
            :disabled="!acbProfileId || !acbActionDate || !acbConfirmed"
            @click="generateACB"
          >
            Generate ACB
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- ── Bank Response Upload Dialog ── -->
    <v-dialog v-model="responseDialog" persistent max-width="500px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2"
          >Upload Bank Response</v-card-title
        >
        <v-card-text>
          <p class="text-body-2 mb-3">
            Upload the bank response file for schedule
            <strong>{{ schedule?.schedule_number }}</strong> to reconcile
            payments.
          </p>
          <v-select
            v-model="responseACBFileId"
            :items="responseACBFiles"
            item-title="file_name"
            item-value="id"
            label="ACB File *"
            variant="outlined"
            density="compact"
            class="mb-3"
            :loading="loadingResponseACBFiles"
          />
          <v-file-input
            v-model="responseFile"
            label="Bank Response File"
            prepend-icon="mdi-file-upload"
            variant="outlined"
            density="compact"
            accept=".txt,.csv"
            hint="ACB response (.txt) or CSV format (.csv)"
            persistent-hint
          />
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="responseDialog = false">Cancel</v-btn>
          <v-btn
            color="deep-purple"
            :loading="processingResponse"
            :disabled="!responseACBFileId || !responseFile"
            @click="processResponse"
          >
            Reconcile
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- ── Upload Proof Dialog ── -->
    <v-dialog v-model="proofDialog" persistent max-width="540px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2"
          >Upload Proof of Payment</v-card-title
        >
        <v-card-text>
          <v-alert
            type="error"
            variant="tonal"
            density="compact"
            class="mb-3"
            icon="mdi-alert-octagon-outline"
          >
            <div class="text-body-2 font-weight-medium mb-1"
              >This action is irreversible.</div
            >
            <div class="text-body-2">
              Uploading a proof of payment will mark
              <strong
                >all {{ schedule?.claims_count }} claim(s) ({{
                  formatCurrency(schedule?.total_amount ?? 0)
                }})</strong
              >
              in schedule
              <strong>{{ schedule?.schedule_number }}</strong> as
              <strong>Paid</strong> and confirm the schedule. Member payment
              statuses cannot be reverted from this screen.
            </div>
          </v-alert>
          <v-file-input
            v-model="proofFile"
            label="Proof of Payment Document"
            prepend-icon="mdi-file-upload"
            variant="outlined"
            density="compact"
            accept=".pdf,.csv,.xlsx,.xls,.png,.jpg,.jpeg"
            class="mb-3"
          />
          <v-textarea
            v-model="proofNotes"
            label="Notes (optional)"
            variant="outlined"
            density="compact"
            rows="3"
            class="mb-3"
          />
          <v-checkbox
            v-model="proofConfirmed"
            density="compact"
            color="success"
            hide-details
            label="I confirm payment has cleared and authorise marking these claims as Paid."
          />
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="proofDialog = false">Cancel</v-btn>
          <v-btn
            color="success"
            :loading="uploadingProof"
            :disabled="!proofFile || !proofConfirmed"
            @click="uploadProof"
          >
            Confirm Payment
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Snackbar -->
    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="4000">
      {{ snackbarMessage }}
      <template #actions>
        <v-btn variant="text" color="white" @click="snackbar = false"
          >Close</v-btn
        >
      </template>
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, provide } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'
import {
  PAYMENT_SCHEDULE_CONTEXT,
  type PaymentSchedule,
  type PaymentProof,
  type BankProfile,
  type ACBFile,
  type ReconResult,
  type ReconSummary,
  type ScheduleItem
} from './payment_schedule_context'

const props = defineProps<{
  scheduleId: string | number
}>()

const route = useRoute()
const router = useRouter()
const { hasPermission } = usePermissionCheck()

// ── State ──────────────────────────────────────────────
const loading = ref(false)
const schedule = ref<PaymentSchedule | null>(null)

const exporting = ref(false)
const downloadingProof = ref<number | null>(null)

// ACB generate
const acbDialog = ref(false)
const acbProfileId = ref<number | null>(null)
const acbActionDate = ref('')
const acbConfirmed = ref(false)
const generatingACB = ref(false)

// ACB files (cached for this schedule)
const acbFiles = ref<ACBFile[]>([])
const loadingACBFiles = ref(false)
const downloadingACB = ref<number | null>(null)

// Bank response upload
const responseDialog = ref(false)
const responseACBFileId = ref<number | null>(null)
const responseACBFiles = ref<ACBFile[]>([])
const loadingResponseACBFiles = ref(false)
const responseFile = ref<File | null>(null)
const processingResponse = ref(false)

// Reconciliation
const reconResults = ref<ReconResult[]>([])
const reconSummary = ref<ReconSummary | null>(null)
const loadingRecon = ref(false)
const retrying = ref(false)

// Bank profiles (for ACB dialog dropdown)
const bankProfiles = ref<BankProfile[]>([])
const loadingProfiles = ref(false)

// Proof upload
const proofDialog = ref(false)
const proofFile = ref<File | null>(null)
const proofNotes = ref('')
const proofConfirmed = ref(false)
const uploadingProof = ref(false)

// Snackbar
const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

// ── Computed ────────────────────────────────────────────
const scheduleId = computed((): number => {
  return typeof props.scheduleId === 'string'
    ? parseInt(props.scheduleId, 10)
    : props.scheduleId
})

const itemsMissingBank = computed((): ScheduleItem[] => {
  if (!schedule.value?.items) return []
  return schedule.value.items.filter(
    (i) => !i.bank_account_number || !i.bank_branch_code
  )
})

const acbBlockedReason = computed((): string | null => {
  if (!schedule.value) return null
  if (schedule.value.status === 'confirmed') {
    return 'Schedule already paid — ACB generation is closed'
  }
  if (schedule.value.acb_file_generated) {
    return 'ACB already generated for this schedule — use Retry Failed Payments for unsuccessful items'
  }
  if (itemsMissingBank.value.length > 0) {
    return `${itemsMissingBank.value.length} item(s) missing bank details — resolve before generating ACB`
  }
  return null
})

const activeTab = computed(() => (route.name as string) ?? '')

function tabRoute(name: string) {
  return { name, params: { scheduleId: props.scheduleId } }
}

// ── Helpers ─────────────────────────────────────────────
function formatCurrency(val: number) {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR'
  }).format(val ?? 0)
}

function formatDate(val?: string) {
  if (!val) return '—'
  return new Date(val).toLocaleDateString('en-ZA', {
    year: 'numeric',
    month: 'short',
    day: '2-digit'
  })
}

// Ordered pipeline steps for the workflow stepper. Each entry maps to a real
// schedule status and supplies a label + caption shown under the step number.
const pipelineSteps = [
  { status: 'draft', label: 'Draft', sub: 'Created' },
  {
    status: 'claims_signed_off',
    label: 'Payment Schedule Signed Off',
    sub: 'Head of Claims'
  },
  { status: 'finance_in_review', label: 'Finance Review', sub: 'In review' },
  {
    status: 'finance_first_authorised',
    label: '1st Authorisation',
    sub: 'Finance'
  },
  {
    status: 'finance_second_authorised',
    label: '2nd Authorisation',
    sub: 'Finance'
  },
  { status: 'submitted_to_bank', label: 'Submitted to Bank', sub: 'ACB run' },
  { status: 'confirmed', label: 'Paid / Confirmed', sub: 'Proof uploaded' }
]

const PIPELINE_STATUSES = pipelineSteps.map((s) => s.status)

const currentStepIndex = computed(() => {
  if (!schedule.value) return -1
  const idx = PIPELINE_STATUSES.indexOf(schedule.value.status)
  if (idx >= 0) return idx
  // Legacy mapping for old "submitted" rows.
  if (schedule.value.status === 'submitted') return 5
  if (schedule.value.status === 'archived') return PIPELINE_STATUSES.length
  return -1
})

function stepTone(idx: number) {
  if (idx < currentStepIndex.value) return 'success'
  if (idx === currentStepIndex.value) return 'current'
  return 'muted'
}

function reconStatusColor(status: string) {
  const map: Record<string, string> = {
    paid: 'success',
    failed: 'error',
    unmatched: 'orange'
  }
  return map[status] ?? 'default'
}

function notify(message: string, color: string = 'success') {
  snackbarMessage.value = message
  snackbarColor.value = color
  snackbar.value = true
}

function unwrap(res: any) {
  const body = res?.data
  if (body && typeof body === 'object' && 'success' in body && 'data' in body) {
    return body.data
  }
  return body
}

function downloadBlob(data: any, filename: string, type: string) {
  const blob = new Blob([data], { type })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}

// ── Navigation ──────────────────────────────────────────
function goBack() {
  const target = hasPermission('claims_pay:finance_review')
    ? 'group-pricing-claim-payment-schedules'
    : 'group-pricing-claim-my-submissions'
  router.push({ name: target })
}

// ── Data loading ────────────────────────────────────────
async function loadSchedule() {
  const id = scheduleId.value
  if (!id || Number.isNaN(id)) {
    notify('Invalid schedule id', 'error')
    return
  }
  loading.value = true
  try {
    const res = await GroupPricingService.getPaymentSchedule(id)
    schedule.value = unwrap(res)
  } catch (e: any) {
    notify('Failed to load payment schedule', 'error')
    schedule.value = null
  } finally {
    loading.value = false
  }
}

async function loadBankProfiles() {
  loadingProfiles.value = true
  try {
    const res = await GroupPricingService.getBankProfiles()
    bankProfiles.value = (unwrap(res) ?? []).filter(
      (p: BankProfile) => p.is_active
    )
  } catch (e: any) {
    notify('Failed to load bank profiles', 'error')
  } finally {
    loadingProfiles.value = false
  }
}

async function loadACBFiles() {
  if (!schedule.value) return
  loadingACBFiles.value = true
  try {
    const res = await GroupPricingService.getACBFileRecords(schedule.value.id)
    acbFiles.value = unwrap(res) ?? []
  } catch {
    acbFiles.value = []
  } finally {
    loadingACBFiles.value = false
  }
}

async function loadReconData() {
  if (!schedule.value) return
  loadingRecon.value = true
  try {
    const summaryRes = await GroupPricingService.getReconciliationSummary(
      schedule.value.id
    )
    reconSummary.value = unwrap(summaryRes)

    const filesRes = await GroupPricingService.getACBFileRecords(
      schedule.value.id
    )
    const files: ACBFile[] = unwrap(filesRes) ?? []
    const allResults: ReconResult[] = []
    for (const f of files) {
      if (f.status === 'reconciled') {
        const res = await GroupPricingService.getReconciliationResults(f.id)
        allResults.push(...(unwrap(res) ?? []))
      }
    }
    reconResults.value = allResults
  } catch {
    reconResults.value = []
    reconSummary.value = null
  } finally {
    loadingRecon.value = false
  }
}

// ── Export ──────────────────────────────────────────────
async function exportSchedule() {
  if (!schedule.value) return
  exporting.value = true
  try {
    const res = await GroupPricingService.exportPaymentScheduleCSV(
      schedule.value.id
    )
    downloadBlob(
      res.data,
      `payment_schedule_${schedule.value.schedule_number}.csv`,
      'text/csv'
    )
    notify('Payment schedule exported.')
    await loadSchedule()
  } catch (e: any) {
    notify('Failed to export schedule', 'error')
  } finally {
    exporting.value = false
  }
}

// ── ACB Generation ──────────────────────────────────────
async function openACBDialog() {
  acbProfileId.value = null
  acbActionDate.value = new Date(Date.now() + 2 * 86400000)
    .toISOString()
    .split('T')[0]
  acbConfirmed.value = false
  await loadBankProfiles()
  acbDialog.value = true
}

async function generateACB() {
  if (!schedule.value || !acbProfileId.value || !acbActionDate.value) return
  generatingACB.value = true
  try {
    const res = await GroupPricingService.generateACBFile(schedule.value.id, {
      bank_profile_id: acbProfileId.value,
      action_date: acbActionDate.value
    })
    acbDialog.value = false
    notify('ACB file generated successfully.')

    const acbRecord = unwrap(res)
    const dlRes = await GroupPricingService.downloadACBFile(acbRecord.id)
    downloadBlob(dlRes.data, acbRecord.file_name, 'text/plain')

    await loadSchedule()
    await loadACBFiles()
  } catch (e: any) {
    notify(
      e?.response?.data?.message ??
        e?.response?.data ??
        'Failed to generate ACB file',
      'error'
    )
  } finally {
    generatingACB.value = false
  }
}

async function downloadACBFile(acb: ACBFile) {
  downloadingACB.value = acb.id
  try {
    const res = await GroupPricingService.downloadACBFile(acb.id)
    downloadBlob(res.data, acb.file_name, 'text/plain')
  } catch {
    notify('Failed to download ACB file', 'error')
  } finally {
    downloadingACB.value = null
  }
}

// ── Bank Response Upload ────────────────────────────────
async function openResponseDialog() {
  if (!schedule.value) return
  responseACBFileId.value = null
  responseFile.value = null
  loadingResponseACBFiles.value = true
  responseDialog.value = true
  try {
    const res = await GroupPricingService.getACBFileRecords(schedule.value.id)
    responseACBFiles.value = (unwrap(res) ?? []).filter(
      (f: ACBFile) => f.status === 'generated'
    )
  } catch {
    responseACBFiles.value = []
  } finally {
    loadingResponseACBFiles.value = false
  }
}

async function processResponse() {
  if (!responseACBFileId.value || !responseFile.value) return
  processingResponse.value = true
  try {
    const formData = new FormData()
    formData.append('file', responseFile.value as File)
    const res = await GroupPricingService.processBankResponse(
      responseACBFileId.value,
      formData
    )
    responseDialog.value = false

    const s = unwrap(res)
    notify(
      `Reconciliation complete: ${s.paid} paid, ${s.failed} failed, ${s.unmatched} unmatched`,
      s.failed > 0 ? 'warning' : 'success'
    )
    await loadSchedule()
    await loadACBFiles()
    await loadReconData()
  } catch (e: any) {
    notify(
      e?.response?.data?.message ??
        e?.response?.data ??
        'Failed to process bank response',
      'error'
    )
  } finally {
    processingResponse.value = false
  }
}

// ── Retry Failed ────────────────────────────────────────
async function retryFailed() {
  if (!schedule.value) return
  retrying.value = true
  try {
    const filesRes = await GroupPricingService.getACBFileRecords(
      schedule.value.id
    )
    const reconciledFile = (unwrap(filesRes) ?? []).find(
      (f: ACBFile) => f.status === 'reconciled'
    )
    if (!reconciledFile) {
      notify('No reconciled ACB file found to retry from', 'error')
      return
    }
    const res = await GroupPricingService.retryFailedPayments(
      reconciledFile.id,
      { item_ids: [] }
    )
    notify('Retry ACB file generated.')

    const retryRecord = unwrap(res)
    const dlRes = await GroupPricingService.downloadACBFile(retryRecord.id)
    downloadBlob(dlRes.data, retryRecord.file_name, 'text/plain')

    await loadReconData()
    await loadACBFiles()
    await loadSchedule()
  } catch (e: any) {
    notify(e?.response?.data ?? 'Failed to retry failed payments', 'error')
  } finally {
    retrying.value = false
  }
}

// ── Proof of Payment ────────────────────────────────────
function openProofDialog() {
  proofFile.value = null
  proofNotes.value = ''
  proofConfirmed.value = false
  proofDialog.value = true
}

async function uploadProof() {
  if (!proofFile.value || !schedule.value) return
  uploadingProof.value = true
  try {
    const formData = new FormData()
    formData.append('file', proofFile.value as File)
    formData.append('notes', proofNotes.value)
    await GroupPricingService.uploadPaymentProof(schedule.value.id, formData)
    proofDialog.value = false
    notify('Proof of payment uploaded. All claims marked as Paid.')
    await loadSchedule()
  } catch (e: any) {
    notify(e?.data ?? 'Failed to upload proof of payment', 'error')
  } finally {
    uploadingProof.value = false
  }
}

async function downloadProof(proof: PaymentProof) {
  downloadingProof.value = proof.id
  try {
    const res = await GroupPricingService.downloadPaymentProof(proof.id)
    downloadBlob(
      res.data,
      proof.file_name,
      proof.content_type || 'application/octet-stream'
    )
  } catch (e: any) {
    notify('Failed to download proof', 'error')
  } finally {
    downloadingProof.value = null
  }
}

// ── Lifecycle actions (Phase 1) ───────────────────────────
const signingOff = ref(false)
const startingReview = ref(false)
const authorising = ref<'' | 'first' | 'second'>('')
const archiving = ref(false)

function errMessage(e: any, fallback: string) {
  return (
    e?.response?.data?.message ??
    (typeof e?.response?.data === 'string' ? e.response.data : null) ??
    e?.message ??
    fallback
  )
}

async function signOff() {
  if (!schedule.value) return
  signingOff.value = true
  try {
    await GroupPricingService.signOffPaymentSchedule(schedule.value.id)
    notify('Head of Claims sign-off recorded.')
    await loadSchedule()
  } catch (e: any) {
    notify(errMessage(e, 'Failed to sign off schedule'), 'error')
  } finally {
    signingOff.value = false
  }
}

async function startFinanceReview() {
  if (!schedule.value) return
  startingReview.value = true
  try {
    await GroupPricingService.startFinanceReview(schedule.value.id)
    notify('Finance review started.')
    await loadSchedule()
  } catch (e: any) {
    notify(errMessage(e, 'Failed to start finance review'), 'error')
  } finally {
    startingReview.value = false
  }
}

async function verifyLineItem(itemId: number) {
  if (!schedule.value) return
  try {
    await GroupPricingService.verifyScheduleLineItem(schedule.value.id, itemId)
    notify('Line item verified.')
    await loadSchedule()
  } catch (e: any) {
    notify(errMessage(e, 'Failed to verify line item'), 'error')
  }
}

async function queryLineItem(
  itemId: number,
  reasonCode: string,
  notes: string
) {
  if (!schedule.value) return
  try {
    await GroupPricingService.queryScheduleLineItem(schedule.value.id, itemId, {
      reason_code: reasonCode,
      notes
    })
    notify('Line queried and returned to claims.', 'warning')
    await loadSchedule()
    await loadQueries()
  } catch (e: any) {
    notify(errMessage(e, 'Failed to query line item'), 'error')
  }
}

async function rejectLineItem(
  itemId: number,
  reasonCode: string,
  notes: string
) {
  if (!schedule.value) return
  try {
    await GroupPricingService.rejectScheduleLineItem(
      schedule.value.id,
      itemId,
      {
        reason_code: reasonCode,
        notes
      }
    )
    notify('Line rejected and returned to claims.', 'warning')
    await loadSchedule()
    await loadQueries()
  } catch (e: any) {
    notify(errMessage(e, 'Failed to reject line item'), 'error')
  }
}

async function authoriseFirst() {
  if (!schedule.value) return
  authorising.value = 'first'
  try {
    await GroupPricingService.firstAuthorisePaymentSchedule(schedule.value.id)
    notify('First authorisation recorded.')
    await loadSchedule()
  } catch (e: any) {
    notify(errMessage(e, 'Failed to authorise (1st)'), 'error')
  } finally {
    authorising.value = ''
  }
}

async function authoriseSecond() {
  if (!schedule.value) return
  authorising.value = 'second'
  try {
    await GroupPricingService.secondAuthorisePaymentSchedule(schedule.value.id)
    notify('Second authorisation recorded — schedule ready for ACB.')
    await loadSchedule()
  } catch (e: any) {
    notify(errMessage(e, 'Failed to authorise (2nd)'), 'error')
  } finally {
    authorising.value = ''
  }
}

async function archive() {
  if (!schedule.value) return
  archiving.value = true
  try {
    await GroupPricingService.archivePaymentSchedule(schedule.value.id)
    notify('Schedule archived.')
    await loadSchedule()
  } catch (e: any) {
    notify(errMessage(e, 'Failed to archive schedule'), 'error')
  } finally {
    archiving.value = false
  }
}

// Queries tab state
const queries = ref<any[]>([])
const loadingQueries = ref(false)

async function loadQueries() {
  if (!schedule.value) return
  loadingQueries.value = true
  try {
    const res = await GroupPricingService.getScheduleQueries(schedule.value.id)
    queries.value = unwrap(res) ?? []
  } catch {
    queries.value = []
  } finally {
    loadingQueries.value = false
  }
}

// ── Sanctions / reinsurance / duplicates (Phase 3) ────────
const sanctions = ref<any[]>([])

async function loadSanctions() {
  if (!schedule.value) return
  try {
    const res = await GroupPricingService.listScheduleSanctions(
      schedule.value.id
    )
    sanctions.value = unwrap(res) ?? []
  } catch {
    sanctions.value = []
  }
}

async function screenLineItem(itemId: number) {
  if (!schedule.value) return
  try {
    await GroupPricingService.screenScheduleLineItem(schedule.value.id, itemId)
    notify('Sanctions screening run — record outcome to clear.', 'info')
    await Promise.all([loadSanctions(), loadSchedule()])
  } catch (e: any) {
    notify(errMessage(e, 'Failed to run sanctions screening'), 'error')
  }
}

async function recordSanctionsOutcome(
  itemId: number,
  status: string,
  notes: string
) {
  if (!schedule.value) return
  try {
    await GroupPricingService.recordSanctionsOutcome(
      schedule.value.id,
      itemId,
      { status, notes }
    )
    notify(`Sanctions outcome recorded: ${status}.`)
    await Promise.all([loadSanctions(), loadSchedule()])
  } catch (e: any) {
    notify(errMessage(e, 'Failed to record sanctions outcome'), 'error')
  }
}

async function setReinsuranceRecovery(
  itemId: number,
  required: boolean,
  amount: number
) {
  if (!schedule.value) return
  try {
    await GroupPricingService.setReinsuranceRecovery(
      schedule.value.id,
      itemId,
      { required, amount }
    )
    notify(
      required
        ? 'Reinsurance recovery flagged.'
        : 'Reinsurance recovery cleared.'
    )
    await loadSchedule()
  } catch (e: any) {
    notify(errMessage(e, 'Failed to update reinsurance recovery'), 'error')
  }
}

async function confirmReinsuranceRaised(itemId: number) {
  if (!schedule.value) return
  try {
    await GroupPricingService.confirmReinsuranceRecoveryRaised(
      schedule.value.id,
      itemId
    )
    notify('Reinsurance recovery raised confirmed.')
    await loadSchedule()
  } catch (e: any) {
    notify(errMessage(e, 'Failed to confirm reinsurance raised'), 'error')
  }
}

async function clearDuplicateBeneficiary(itemId: number) {
  if (!schedule.value) return
  try {
    await GroupPricingService.clearDuplicateBeneficiary(
      schedule.value.id,
      itemId
    )
    notify('Duplicate beneficiary flag cleared.')
    await loadSchedule()
  } catch (e: any) {
    notify(errMessage(e, 'Failed to clear duplicate flag'), 'error')
  }
}

// ── Tax certificates (Phase 4) ──────────────────────────
const taxCertificates = ref<any[]>([])

async function loadTaxCertificates() {
  if (!schedule.value) return
  try {
    const res = await GroupPricingService.listScheduleTaxCertificates(
      schedule.value.id
    )
    taxCertificates.value = unwrap(res) ?? []
  } catch {
    taxCertificates.value = []
  }
}

async function downloadTaxCertificate(cert: any) {
  try {
    const res = await GroupPricingService.downloadTaxCertificate(cert.id)
    downloadBlob(res.data, cert.file_name, cert.content_type || 'text/html')
  } catch (e: any) {
    notify(errMessage(e, 'Failed to download tax certificate'), 'error')
  }
}

// ── Action gating helpers (used in template) ──────────────
const canSignOff = computed(() => {
  if (!schedule.value) return false
  return (
    schedule.value.status === 'draft' &&
    hasPermission('claims_pay:signoff_schedule')
  )
})

const canStartReview = computed(() => {
  if (!schedule.value) return false
  return (
    schedule.value.status === 'claims_signed_off' &&
    hasPermission('claims_pay:finance_review')
  )
})

const canAuthoriseFirst = computed(() => {
  if (!schedule.value) return false
  if (schedule.value.status !== 'finance_in_review') return false
  if (!hasPermission('claims_pay:authorise_first')) return false
  // Block while any line is still pending.
  return !schedule.value.items?.some(
    (i) => !i.line_status || i.line_status === 'pending'
  )
})

const canAuthoriseSecond = computed(() => {
  if (!schedule.value) return false
  return (
    schedule.value.status === 'finance_first_authorised' &&
    hasPermission('claims_pay:authorise_second')
  )
})

const canArchive = computed(() => {
  if (!schedule.value) return false
  return (
    schedule.value.status === 'confirmed' && hasPermission('claims_pay:archive')
  )
})

// ── Phase 3 banner counters ───────────────────────────────
const outstandingDuplicates = computed(() => {
  if (!schedule.value?.items) return 0
  return schedule.value.items.filter(
    (i: any) =>
      i.duplicate_beneficiary_flag === true &&
      i.duplicate_beneficiary_cleared !== true &&
      (i.line_status === 'pending' ||
        i.line_status === 'verified' ||
        !i.line_status)
  ).length
})

const sanctionsBlockers = computed(() => {
  if (!schedule.value?.items) return 0
  const itemsById = new Map<number, string[]>()
  for (const row of sanctions.value) {
    const list = itemsById.get(row.schedule_item_id) ?? []
    list.push(row.status)
    itemsById.set(row.schedule_item_id, list)
  }
  let count = 0
  for (const item of schedule.value.items) {
    if (
      item.line_status !== 'pending' &&
      item.line_status !== 'verified' &&
      item.line_status
    ) {
      continue
    }
    const statuses = itemsById.get(item.id) ?? []
    // Block when item has no clear / manual_clear status across providers.
    const cleared = statuses.some((s) => s === 'clear' || s === 'manual_clear')
    const blocking = statuses.some(
      (s) => s === 'hit' || s === 'pending' || s === 'skipped'
    )
    if (blocking || (!cleared && statuses.length === 0)) {
      count++
    }
  }
  return count
})

const reinsuranceOutstanding = computed(() => {
  if (!schedule.value?.items) return 0
  return schedule.value.items.filter(
    (i: any) =>
      i.reinsurance_recovery_required &&
      !i.reinsurance_recovery_raised_at &&
      (i.line_status === 'pending' ||
        i.line_status === 'verified' ||
        !i.line_status)
  ).length
})

// Provide context to child tab components
provide(PAYMENT_SCHEDULE_CONTEXT, {
  schedule,
  formatCurrency,
  formatDate,
  hasPermission,
  notify,
  acbFiles,
  loadingACBFiles,
  downloadingACB,
  loadACBFiles,
  downloadACBFile,
  reconResults,
  reconSummary,
  loadingRecon,
  retrying,
  loadReconData,
  retryFailed,
  reconStatusColor,
  downloadingProof,
  downloadProof,
  signOff,
  startFinanceReview,
  verifyLineItem,
  queryLineItem,
  rejectLineItem,
  authoriseFirst,
  authoriseSecond,
  archive,
  refreshSchedule: loadSchedule,
  queries,
  loadingQueries,
  loadQueries,
  sanctions,
  loadSanctions,
  screenLineItem,
  recordSanctionsOutcome,
  setReinsuranceRecovery,
  confirmReinsuranceRaised,
  clearDuplicateBeneficiary,
  taxCertificates,
  loadTaxCertificates,
  downloadTaxCertificate
})

onMounted(async () => {
  await loadSchedule()
  await Promise.all([loadSanctions(), loadTaxCertificates()])
})
</script>

<style scoped>
.schedule-tabs {
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}

/* ── Page sections ────────────────────────────────────── */
.page-section {
  margin-bottom: 28px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.section-label {
  font-size: 0.7rem;
  font-weight: 700;
  letter-spacing: 1.2px;
  text-transform: uppercase;
  color: rgb(var(--v-theme-primary));
}

.section-divider {
  flex: 1;
  height: 1px;
  background: linear-gradient(
    to right,
    rgba(var(--v-theme-primary), 0.18),
    rgba(var(--v-theme-primary), 0.02)
  );
}

.meta-card {
  min-height: 96px;
  transition: border-color 0.15s ease;
}

.meta-card:hover {
  border-color: rgba(var(--v-theme-primary), 0.4);
}

.meta-card__label {
  font-size: 0.7rem;
  font-weight: 500;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  color: rgba(var(--v-theme-on-surface), 0.6);
  margin-bottom: 4px;
}

.meta-card__value {
  font-size: 1.5rem;
  font-weight: 700;
  line-height: 1.2;
  color: rgba(var(--v-theme-on-surface), 0.95);
  word-break: break-word;
  overflow-wrap: break-word;
}

.meta-card__value--text {
  font-size: 1.05rem;
  font-weight: 600;
  line-height: 1.3;
}

.meta-card__hint {
  margin-top: auto;
  padding-top: 4px;
  font-size: 0.72rem;
  color: rgba(var(--v-theme-on-surface), 0.55);
}

/* ── Workflow stepper ─────────────────────────────────── */
.workflow-stepper {
  display: flex;
  align-items: stretch;
  flex-wrap: wrap;
  gap: 8px;
  padding: 14px 16px;
  background: rgba(var(--v-theme-primary), 0.04);
  border: 1px solid rgba(var(--v-theme-primary), 0.1);
  border-radius: 10px;
}

.workflow-stepper__step {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1 1 0;
  min-width: 150px;
}

.workflow-stepper__num {
  flex: 0 0 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.78rem;
  font-weight: 700;
  color: #fff;
  background: rgba(var(--v-theme-on-surface), 0.35);
}

.workflow-stepper__step--muted .workflow-stepper__num {
  background: rgba(var(--v-theme-on-surface), 0.28);
}
.workflow-stepper__step--current .workflow-stepper__num {
  background: rgb(var(--v-theme-primary));
  box-shadow: 0 0 0 4px rgba(var(--v-theme-primary), 0.18);
}
.workflow-stepper__step--success .workflow-stepper__num {
  background: rgb(var(--v-theme-success));
}

.workflow-stepper__text {
  flex: 1;
  min-width: 0;
}

.workflow-stepper__label {
  font-size: 0.82rem;
  font-weight: 600;
  color: rgba(var(--v-theme-on-surface), 0.88);
  line-height: 1.2;
}

.workflow-stepper__step--muted .workflow-stepper__label {
  color: rgba(var(--v-theme-on-surface), 0.55);
}

.workflow-stepper__step--current .workflow-stepper__label {
  color: rgb(var(--v-theme-primary));
}

.workflow-stepper__sub {
  font-size: 0.7rem;
  color: rgba(var(--v-theme-on-surface), 0.55);
  line-height: 1.2;
}

.workflow-stepper__arrow {
  color: rgba(var(--v-theme-primary), 0.35);
  flex: 0 0 auto;
}
</style>
