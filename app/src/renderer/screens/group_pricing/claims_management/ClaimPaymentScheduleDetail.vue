<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div
              class="d-flex justify-space-between align-center flex-wrap gap-2"
            >
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
                    - {{ schedule.schedule_number }}
                  </template>
                </span>
                <v-chip
                  v-if="schedule"
                  :color="statusColor(schedule.status)"
                  size="small"
                  label
                >
                  {{ statusLabel(schedule.status) }}
                </v-chip>
                <v-chip
                  v-if="schedule?.acb_file_generated"
                  color="teal"
                  size="small"
                  label
                >
                  ACB Generated
                </v-chip>
              </div>
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
              <!-- Action buttons strip -->
              <v-card variant="tonal" color="grey-lighten-4" class="pa-3 mb-4">
                <div class="d-flex align-center flex-wrap gap-2">
                  <span
                    class="text-body-2 font-weight-medium text-medium-emphasis mr-2"
                    >Actions:</span
                  >
                  <v-btn
                    variant="outlined"
                    size="small"
                    rounded
                    color="info"
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
                          color="teal"
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
                    color="deep-purple"
                    prepend-icon="mdi-file-upload"
                    @click="openResponseDialog"
                  >
                    Upload Bank Response
                  </v-btn>
                  <v-btn
                    v-if="schedule.status !== 'confirmed'"
                    variant="flat"
                    size="small"
                    rounded
                    color="success"
                    prepend-icon="mdi-upload"
                    @click="openProofDialog"
                  >
                    Upload Proof of Payment
                  </v-btn>
                  <v-spacer />
                  <v-chip
                    v-if="acbBlockedReason && hasPermission('claims_pay:generate_acb')"
                    :color="
                      schedule.acb_file_generated || schedule.status === 'confirmed'
                        ? 'info'
                        : 'warning'
                    "
                    size="small"
                    variant="tonal"
                    prepend-icon="mdi-information-outline"
                  >
                    {{ acbBlockedReason }}
                  </v-chip>
                </div>
              </v-card>

              <!-- Metadata grid -->
              <v-row dense class="mb-4">
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
                      claim{{ schedule.claims_count === 1 ? '' : 's' }} in this
                      schedule
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

              <!-- Description -->
              <v-card
                v-if="schedule.description"
                variant="tonal"
                color="grey-lighten-4"
                class="mb-4 pa-3"
                flat
              >
                <div class="text-caption text-medium-emphasis"
                  >Description</div
                >
                <div class="text-body-2">{{ schedule.description }}</div>
              </v-card>

              <!-- Tabs -->
              <v-tabs v-model="viewTab" density="compact" class="mb-3">
                <v-tab value="claims">
                  <v-icon size="18" class="mr-1">mdi-file-document-multiple</v-icon>
                  Claims
                </v-tab>
                <v-tab value="acb">
                  <v-icon size="18" class="mr-1">mdi-bank-transfer</v-icon>
                  ACB Files
                </v-tab>
                <v-tab value="reconciliation">
                  <v-icon size="18" class="mr-1">mdi-scale-balance</v-icon>
                  Reconciliation
                </v-tab>
                <v-tab value="proofs">
                  <v-icon size="18" class="mr-1">mdi-receipt-text-check</v-icon>
                  Proof of Payment
                </v-tab>
              </v-tabs>

              <v-tabs-window v-model="viewTab">
                <!-- Claims Tab -->
                <v-tabs-window-item value="claims">
                  <v-table density="compact" class="border rounded">
                    <thead>
                      <tr>
                        <th>Claim #</th>
                        <th>Member</th>
                        <th>ID Number</th>
                        <th>Scheme</th>
                        <th>Benefit</th>
                        <th>Bank</th>
                        <th class="text-right">Amount</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="item in schedule.items" :key="item.id">
                        <td>{{ item.claim_number }}</td>
                        <td>{{ item.member_name }}</td>
                        <td>{{ item.member_id_number }}</td>
                        <td>{{ item.scheme_name }}</td>
                        <td>{{ item.benefit_name }}</td>
                        <td>
                          <v-chip
                            v-if="item.bank_account_number"
                            size="x-small"
                            color="teal"
                            variant="tonal"
                          >
                            {{ item.bank_name || 'Set' }}
                          </v-chip>
                          <v-chip
                            v-else
                            size="x-small"
                            color="orange"
                            variant="tonal"
                          >
                            Missing
                          </v-chip>
                        </td>
                        <td class="text-right">{{
                          formatCurrency(item.claim_amount)
                        }}</td>
                      </tr>
                    </tbody>
                  </v-table>
                </v-tabs-window-item>

                <!-- ACB Files Tab -->
                <v-tabs-window-item value="acb">
                  <v-progress-linear
                    v-if="loadingACBFiles"
                    indeterminate
                    color="teal"
                    class="mb-2"
                  />
                  <empty-state
                    v-if="!loadingACBFiles && acbFiles.length === 0"
                    icon="mdi-file-document-outline"
                    title="No ACB files generated yet"
                    message="Generate an ACB file to start a BankServ payment run."
                  />
                  <v-list v-else density="compact" class="border rounded">
                    <v-list-item
                      v-for="acb in acbFiles"
                      :key="acb.id"
                      :subtitle="`Generated by ${acb.generated_by} on ${formatDate(acb.generated_at)} | ${acb.transaction_count} transactions | ${formatCurrency(acb.total_amount)}`"
                      :title="acb.file_name"
                    >
                      <template #prepend>
                        <v-icon
                          :color="
                            acb.status === 'reconciled' ? 'success' : 'teal'
                          "
                        >
                          {{
                            acb.status === 'reconciled'
                              ? 'mdi-check-circle'
                              : 'mdi-file-document'
                          }}
                        </v-icon>
                      </template>
                      <template #append>
                        <div class="d-flex gap-1 align-center">
                          <v-chip
                            v-if="acb.is_retry"
                            size="x-small"
                            color="orange"
                            variant="tonal"
                            class="mr-1"
                          >
                            Retry
                          </v-chip>
                          <v-chip
                            :color="
                              acb.status === 'reconciled' ? 'success' : 'grey'
                            "
                            size="x-small"
                            label
                          >
                            {{ acb.status }}
                          </v-chip>
                          <v-btn
                            size="x-small"
                            variant="text"
                            color="primary"
                            icon="mdi-download"
                            :loading="downloadingACB === acb.id"
                            @click="downloadACBFile(acb)"
                          />
                        </div>
                      </template>
                    </v-list-item>
                  </v-list>
                </v-tabs-window-item>

                <!-- Reconciliation Tab -->
                <v-tabs-window-item value="reconciliation">
                  <v-progress-linear
                    v-if="loadingRecon"
                    indeterminate
                    color="deep-purple"
                    class="mb-2"
                  />

                  <div
                    v-if="reconSummary"
                    class="d-flex gap-2 mb-3 flex-wrap"
                  >
                    <v-chip color="default" variant="tonal">
                      Total: {{ reconSummary.total_transactions }}
                    </v-chip>
                    <v-chip color="success" variant="tonal">
                      Paid: {{ reconSummary.paid }} ({{
                        formatCurrency(reconSummary.total_paid)
                      }})
                    </v-chip>
                    <v-chip color="error" variant="tonal">
                      Failed: {{ reconSummary.failed }} ({{
                        formatCurrency(reconSummary.total_failed)
                      }})
                    </v-chip>
                    <v-chip color="orange" variant="tonal">
                      Unmatched: {{ reconSummary.unmatched }}
                    </v-chip>
                  </div>

                  <empty-state
                    v-if="!loadingRecon && reconResults.length === 0"
                    icon="mdi-scale-balance"
                    title="No reconciliation data yet"
                    message="Upload a bank response file to reconcile this payment run."
                  />
                  <v-table v-else density="compact" class="border rounded mb-3">
                    <thead>
                      <tr>
                        <th>Claim #</th>
                        <th>Account</th>
                        <th class="text-right">Amount</th>
                        <th>Status</th>
                        <th>Failure Reason</th>
                        <th>Bank Ref</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="r in reconResults" :key="r.id">
                        <td>{{ r.claim_number || '—' }}</td>
                        <td>{{ r.account_number }}</td>
                        <td class="text-right">{{
                          formatCurrency(r.amount)
                        }}</td>
                        <td>
                          <v-chip
                            :color="reconStatusColor(r.status)"
                            size="x-small"
                            label
                            variant="flat"
                          >
                            {{ r.status }}
                          </v-chip>
                        </td>
                        <td class="text-caption">{{
                          r.failure_reason || '—'
                        }}</td>
                        <td class="text-caption">{{
                          r.bank_reference || '—'
                        }}</td>
                      </tr>
                    </tbody>
                  </v-table>

                  <v-btn
                    v-if="
                      hasPermission('claims_pay:retry_failed') &&
                      reconResults.some((r: any) => r.status === 'failed')
                    "
                    color="orange"
                    size="small"
                    variant="outlined"
                    prepend-icon="mdi-refresh"
                    :loading="retrying"
                    @click="retryFailed"
                  >
                    Retry Failed Payments
                  </v-btn>
                </v-tabs-window-item>

                <!-- Proofs Tab -->
                <v-tabs-window-item value="proofs">
                  <empty-state
                    v-if="
                      !schedule.proof_of_payments ||
                      schedule.proof_of_payments.length === 0
                    "
                    icon="mdi-receipt-text-outline"
                    title="No proof of payment uploaded yet"
                    message="Upload proof of payment to confirm this schedule and mark all claims as Paid."
                  />
                  <v-list v-else density="compact" class="border rounded">
                    <v-list-item
                      v-for="proof in schedule.proof_of_payments"
                      :key="proof.id"
                      :subtitle="`Uploaded by ${proof.uploaded_by} on ${formatDate(proof.uploaded_at)}${proof.notes ? ' — ' + proof.notes : ''}`"
                      :title="proof.file_name"
                    >
                      <template #prepend>
                        <v-icon color="success">mdi-receipt-text-check</v-icon>
                      </template>
                      <template #append>
                        <v-btn
                          size="x-small"
                          variant="text"
                          color="primary"
                          icon="mdi-download"
                          :loading="downloadingProof === proof.id"
                          @click="downloadProof(proof)"
                        />
                      </template>
                    </v-list-item>
                  </v-list>
                </v-tabs-window-item>
              </v-tabs-window>
            </template>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- ── ACB Generate Dialog ── -->
    <v-dialog v-model="acbDialog" persistent max-width="540px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2"
          >Generate ACB File</v-card-title
        >
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
                >all
                {{ schedule?.claims_count }} claim(s) ({{
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

interface ScheduleItem {
  id: number
  claim_id: number
  claim_number: string
  member_name: string
  member_id_number: string
  benefit_name: string
  scheme_name: string
  scheme_id: number
  claim_amount: number
  bank_name?: string
  bank_branch_code?: string
  bank_account_number?: string
  bank_account_type?: string
  account_holder_name?: string
}

interface PaymentProof {
  id: number
  schedule_id: number
  file_name: string
  content_type: string
  size_bytes: number
  notes: string
  uploaded_by: string
  uploaded_at: string
}

interface PaymentSchedule {
  id: number
  schedule_number: string
  description: string
  status: string
  total_amount: number
  claims_count: number
  exported_at?: string
  exported_by?: string
  acb_file_generated?: boolean
  acb_generated_at?: string
  acb_generated_by?: string
  created_by: string
  created_at: string
  items: ScheduleItem[]
  proof_of_payments: PaymentProof[]
}

interface BankProfile {
  id: number
  profile_name: string
  bank_name: string
  user_code: string
  user_branch_code: string
  user_account_number: string
  user_account_type: string
  bank_type_code: string
  service_type: string
  generation_number: number
  is_active: boolean
}

interface ACBFile {
  id: number
  schedule_id: number
  file_name: string
  action_date: string
  transaction_count: number
  total_amount: number
  status: string
  is_retry: boolean
  generated_by: string
  generated_at: string
}

interface ReconResult {
  id: number
  claim_number: string
  account_number: string
  amount: number
  status: string
  failure_reason: string
  bank_reference: string
}

const props = defineProps<{
  scheduleId: string | number
}>()

const router = useRouter()
const { hasPermission } = usePermissionCheck()

// ── State ──────────────────────────────────────────────
const loading = ref(false)
const schedule = ref<PaymentSchedule | null>(null)

const viewTab = ref('claims')

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
const reconSummary = ref<any>(null)
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

function statusLabel(status: string) {
  const map: Record<string, string> = {
    submitted: 'Submitted for Payment',
    confirmed: 'Paid / Confirmed',
    draft: 'Draft'
  }
  return map[status] ?? status
}

function statusColor(status: string) {
  const map: Record<string, string> = {
    submitted: 'warning',
    confirmed: 'success',
    draft: 'grey',
    approved: 'info',
    submitted_for_payment: 'warning',
    paid: 'success',
    payment_failed: 'error',
    pending: 'default',
    declined: 'error'
  }
  return map[status] ?? 'default'
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
  router.push({ name: 'group-pricing-claim-payment-schedules' })
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

// Lazy-load tab data on first switch
watch(viewTab, async (tab) => {
  if (!schedule.value) return
  if (tab === 'acb' && acbFiles.value.length === 0) {
    await loadACBFiles()
  } else if (tab === 'reconciliation' && reconResults.value.length === 0) {
    await loadReconData()
  }
})

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

onMounted(loadSchedule)
</script>

<style scoped>
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
</style>
