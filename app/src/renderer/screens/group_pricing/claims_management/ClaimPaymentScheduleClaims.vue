<template>
  <div v-if="schedule">
    <!-- ── Section: Totals ─────────────────────── -->
    <section class="page-section">
      <div class="section-header">
        <span class="section-label">Totals</span>
        <span class="section-divider" />
      </div>
      <v-row dense>
        <v-col cols="12" sm="4">
          <v-card
            variant="outlined"
            rounded="lg"
            class="totals-card h-100 pa-3 d-flex flex-column"
          >
            <div class="totals-card__label">Gross</div>
            <div class="totals-card__value">{{
              formatCurrency(schedule.gross_total ?? schedule.total_amount)
            }}</div>
          </v-card>
        </v-col>
        <v-col cols="12" sm="4">
          <v-card
            variant="outlined"
            rounded="lg"
            class="totals-card h-100 pa-3 d-flex flex-column"
          >
            <div class="totals-card__label">Deductions</div>
            <div class="totals-card__value">{{
              formatCurrency(schedule.deductions_total ?? 0)
            }}</div>
          </v-card>
        </v-col>
        <v-col cols="12" sm="4">
          <v-card
            variant="outlined"
            rounded="lg"
            class="totals-card totals-card--accent h-100 pa-3 d-flex flex-column"
          >
            <div class="totals-card__label">Net payable</div>
            <div class="totals-card__value totals-card__value--primary">{{
              formatCurrency(schedule.net_total ?? schedule.total_amount)
            }}</div>
          </v-card>
        </v-col>
      </v-row>
    </section>

    <!-- ── Section: Claim lines ─────────────────── -->
    <section class="page-section">
      <div class="section-header">
        <span class="section-label">Claim lines</span>
        <span class="section-divider" />
      </div>
      <v-table density="compact" class="border rounded">
        <thead>
          <tr>
            <th>Claim #</th>
            <th>Beneficiary</th>
            <th>Member</th>
            <th>Scheme / Benefit</th>
            <th>Bank</th>
            <th class="text-right">Gross</th>
            <th class="text-right">Deductions</th>
            <th class="text-right">Net</th>
            <th>Flags</th>
            <th>Sanctions</th>
            <th>Reinsurance</th>
            <th>Status</th>
            <th v-if="isFinanceReview">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in schedule.items" :key="item.id">
            <td>
              <div class="font-weight-medium">{{ item.claim_number }}</div>
              <div
                v-if="item.approval_reference"
                class="text-caption text-medium-emphasis"
                >Approved by {{ item.approval_reference }}</div
              >
            </td>
            <td>
              <div>{{ item.beneficiary_name || item.member_name }}</div>
              <div class="text-caption text-medium-emphasis">
                {{ item.beneficiary_id_number || item.member_id_number }}
                <v-tooltip
                  v-if="beneficiaryMismatch(item)"
                  text="Beneficiary differs from member — verify before paying"
                >
                  <template #activator="{ props: tipProps }">
                    <v-icon
                      v-bind="tipProps"
                      size="14"
                      color="warning"
                      icon="mdi-alert-circle"
                    />
                  </template>
                </v-tooltip>
              </div>
            </td>
            <td>
              <div>{{ item.member_name }}</div>
              <div class="text-caption text-medium-emphasis">{{
                item.member_id_number
              }}</div>
            </td>
            <td>
              <div>{{ item.scheme_name }}</div>
              <div class="text-caption text-medium-emphasis">{{
                item.benefit_name
              }}</div>
            </td>
            <td>
              <v-chip
                v-if="item.bank_account_number"
                size="x-small"
                color="teal"
                variant="tonal"
              >
                {{ item.bank_name || 'Set' }}
              </v-chip>
              <v-chip v-else size="x-small" color="orange" variant="tonal">
                Missing
              </v-chip>
            </td>
            <td class="text-right">{{
              formatCurrency(item.gross_amount ?? item.claim_amount)
            }}</td>
            <td class="text-right">
              <v-tooltip v-if="totalDeductions(item) > 0" location="top">
                <template #activator="{ props: tipProps }">
                  <span v-bind="tipProps">{{
                    formatCurrency(totalDeductions(item))
                  }}</span>
                </template>
                <div class="text-caption">
                  <div
                    >Premium arrears:
                    {{
                      formatCurrency(item.premium_arrears_deduction ?? 0)
                    }}</div
                  >
                  <div
                    >Policy loan:
                    {{ formatCurrency(item.policy_loan_deduction ?? 0) }}</div
                  >
                  <div
                    >Tax withheld:
                    {{ formatCurrency(item.tax_withheld ?? 0) }}</div
                  >
                </div>
              </v-tooltip>
              <span v-else>—</span>
            </td>
            <td class="text-right font-weight-medium">{{
              formatCurrency(item.net_payable ?? item.claim_amount)
            }}</td>
            <td>
              <div class="d-flex flex-wrap gap-1">
                <v-chip
                  v-if="riskFlag(item, 'banking_change_30d')"
                  size="x-small"
                  color="warning"
                  variant="flat"
                  title="Banking details changed in last 30 days"
                  >Bank 30d</v-chip
                >
                <v-chip
                  v-if="riskFlag(item, 'contestable')"
                  size="x-small"
                  color="orange"
                  variant="flat"
                  title="Within the contestability window"
                  >Contestable</v-chip
                >
                <v-chip
                  v-if="riskFlag(item, 'recent_reinstatement')"
                  size="x-small"
                  color="amber"
                  variant="flat"
                  title="Policy reinstated recently"
                  >Reinstated</v-chip
                >
                <v-chip
                  v-if="fraudLevel(item)"
                  size="x-small"
                  :color="fraudColor(fraudLevel(item))"
                  variant="flat"
                >
                  {{ fraudLevel(item) }}
                </v-chip>
                <v-chip
                  v-if="
                    item.duplicate_beneficiary_flag &&
                    !item.duplicate_beneficiary_cleared
                  "
                  size="x-small"
                  color="warning"
                  variant="flat"
                  title="Same beneficiary appears on another line in this schedule"
                >
                  Duplicate
                </v-chip>
              </div>
            </td>
            <td>
              <div class="d-flex align-center gap-1">
                <v-chip
                  :color="sanctionsChipColor(item.id)"
                  size="x-small"
                  variant="flat"
                  :title="sanctionsHitSummary(item.id)"
                >
                  {{ sanctionsLabel(item.id) }}
                </v-chip>
                <v-btn
                  v-if="canScreenSanctions"
                  size="x-small"
                  variant="text"
                  icon="mdi-magnify-scan"
                  title="Run / re-run sanctions screening"
                  @click="onScreen(item.id)"
                />
                <v-btn
                  v-if="canScreenSanctions"
                  size="x-small"
                  variant="text"
                  icon="mdi-pencil-outline"
                  title="Record sanctions outcome"
                  @click="openSanctionsDialog(item)"
                />
              </div>
            </td>
            <td>
              <div v-if="item.reinsurance_recovery_required">
                <v-chip
                  :color="
                    item.reinsurance_recovery_raised_at ? 'success' : 'warning'
                  "
                  size="x-small"
                  variant="flat"
                >
                  {{
                    item.reinsurance_recovery_raised_at
                      ? `Raised ${formatCurrency(item.reinsurance_recovery_amount ?? 0)}`
                      : `Required ${formatCurrency(item.reinsurance_recovery_amount ?? 0)}`
                  }}
                </v-chip>
                <div class="d-flex gap-1 mt-1">
                  <v-btn
                    v-if="
                      !item.reinsurance_recovery_raised_at && canConfirmRecovery
                    "
                    size="x-small"
                    variant="outlined"
                    color="success"
                    @click="onConfirmRaised(item.id)"
                  >
                    Mark raised
                  </v-btn>
                  <v-btn
                    v-if="canEditReinsurance"
                    size="x-small"
                    variant="text"
                    icon="mdi-pencil-outline"
                    @click="openReinsuranceDialog(item)"
                  />
                </div>
              </div>
              <v-btn
                v-else-if="canEditReinsurance"
                size="x-small"
                variant="text"
                prepend-icon="mdi-plus"
                @click="openReinsuranceDialog(item)"
              >
                Flag
              </v-btn>
              <span v-else class="text-medium-emphasis">—</span>
            </td>
            <td>
              <div class="d-flex align-center gap-1">
                <v-chip
                  :color="lineStatusColor(item.line_status)"
                  size="x-small"
                  variant="flat"
                >
                  {{ item.line_status || 'pending' }}
                </v-chip>
                <template v-if="taxCertForItem(item.id)">
                  <v-btn
                    size="x-small"
                    variant="text"
                    icon="mdi-certificate-outline"
                    :title="`IT3(a) ${taxCertForItem(item.id)!.certificate_ref}`"
                    @click="downloadTaxCertificate(taxCertForItem(item.id)!)"
                  />
                </template>
                <v-btn
                  v-if="canGenerateLetter"
                  size="x-small"
                  variant="text"
                  icon="mdi-file-document-check-outline"
                  title="Payment confirmation letter"
                  @click="openLetterDialog(item)"
                />
              </div>
            </td>
            <td v-if="isFinanceReview">
              <div class="d-flex gap-1">
                <v-btn
                  v-if="canVerify(item)"
                  size="x-small"
                  variant="flat"
                  color="success"
                  icon="mdi-check"
                  title="Verify"
                  @click="verifyLineItem(item.id)"
                />
                <v-btn
                  v-if="canQuery(item)"
                  size="x-small"
                  variant="outlined"
                  color="warning"
                  icon="mdi-comment-alert-outline"
                  title="Query"
                  @click="openQueryDialog(item, 'query')"
                />
                <v-btn
                  v-if="canQuery(item)"
                  size="x-small"
                  variant="outlined"
                  color="error"
                  icon="mdi-close"
                  title="Reject"
                  @click="openQueryDialog(item, 'reject')"
                />
                <v-btn
                  v-if="
                    item.duplicate_beneficiary_flag &&
                    !item.duplicate_beneficiary_cleared
                  "
                  size="x-small"
                  variant="outlined"
                  color="indigo"
                  icon="mdi-account-check-outline"
                  title="Clear duplicate beneficiary flag"
                  @click="onClearDuplicate(item.id)"
                />
              </div>
            </td>
          </tr>
        </tbody>
      </v-table>
    </section>

    <!-- Sanctions outcome dialog (Phase 3) -->
    <v-dialog v-model="sanctionsDialog" max-width="520px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2"
          >Record sanctions / PEP outcome</v-card-title
        >
        <v-card-text>
          <div class="text-body-2 text-medium-emphasis mb-3">
            Beneficiary:
            <strong>{{ dialogItem?.beneficiary_name }}</strong>
          </div>
          <v-select
            v-model="sanctionsStatus"
            :items="[
              { title: 'Clear', value: 'clear' },
              { title: 'Hit (blocks authorisation)', value: 'hit' },
              {
                title: 'Manually cleared (false-positive reviewed)',
                value: 'manual_clear'
              }
            ]"
            label="Outcome *"
            variant="outlined"
            density="compact"
            class="mb-3"
          />
          <v-textarea
            v-model="sanctionsNotes"
            label="Notes"
            variant="outlined"
            density="compact"
            rows="3"
            placeholder="Provider ref, matched list, false-positive reasoning..."
          />
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="sanctionsDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :disabled="!sanctionsStatus"
            :loading="savingSanctions"
            @click="submitSanctions"
          >
            Save
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Reinsurance recovery dialog (Phase 3) -->
    <v-dialog v-model="reinsuranceDialog" max-width="480px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2"
          >Reinsurance recovery</v-card-title
        >
        <v-card-text>
          <v-switch
            v-model="reinsuranceRequired"
            color="primary"
            label="Required for this claim"
            density="compact"
            hide-details
            class="mb-3"
          />
          <v-text-field
            v-model.number="reinsuranceAmount"
            type="number"
            label="Recovery amount (ZAR)"
            variant="outlined"
            density="compact"
            :disabled="!reinsuranceRequired"
          />
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="reinsuranceDialog = false"
            >Cancel</v-btn
          >
          <v-btn
            color="primary"
            :loading="savingReinsurance"
            @click="submitReinsurance"
          >
            Save
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Query / Reject dialog -->
    <v-dialog v-model="queryDialog" max-width="540px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2">
          {{ dialogMode === 'reject' ? 'Reject line item' : 'Query line item' }}
        </v-card-title>
        <v-card-text>
          <v-alert
            type="warning"
            variant="tonal"
            density="compact"
            class="mb-3"
            icon="mdi-alert-outline"
          >
            <div class="text-body-2">
              {{
                dialogMode === 'reject'
                  ? 'Rejecting this line removes it from the schedule. The claim returns to the Claims team for resolution and will not appear on the next cut-off until they re-approve it.'
                  : 'Querying this line removes it from the schedule. The claim returns to the Claims team and is eligible for the next cut-off once resolved.'
              }}
            </div>
          </v-alert>
          <v-select
            v-model="reasonCode"
            :items="REASON_CODES"
            label="Reason code *"
            variant="outlined"
            density="compact"
            class="mb-3"
          />
          <v-textarea
            v-model="reasonNotes"
            label="Notes"
            variant="outlined"
            density="compact"
            rows="3"
            placeholder="What did finance find? Be specific — these notes are visible to claims."
          />
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="queryDialog = false">Cancel</v-btn>
          <v-btn
            :color="dialogMode === 'reject' ? 'error' : 'warning'"
            :disabled="!reasonCode"
            :loading="submittingQuery"
            @click="submitQuery"
          >
            {{ dialogMode === 'reject' ? 'Reject line' : 'Query line' }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <PaymentLetterDialog
      v-model="letterDialog"
      :claim-id="letterClaimId"
      :claimant-email="letterClaimContact.email"
      :claimant-phone="letterClaimContact.phone"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, inject, ref } from 'vue'
import {
  PAYMENT_SCHEDULE_CONTEXT,
  type ScheduleItem,
  type RiskFlags
} from './payment_schedule_context'
import PaymentLetterDialog from './components/PaymentLetterDialog.vue'

const ctx = inject(PAYMENT_SCHEDULE_CONTEXT)
if (!ctx) {
  throw new Error(
    'ClaimPaymentScheduleClaims must be rendered inside ClaimPaymentScheduleLayout'
  )
}

const {
  schedule,
  formatCurrency,
  hasPermission,
  verifyLineItem,
  queryLineItem,
  rejectLineItem,
  sanctions,
  screenLineItem,
  recordSanctionsOutcome,
  setReinsuranceRecovery,
  confirmReinsuranceRaised,
  clearDuplicateBeneficiary,
  taxCertificates,
  downloadTaxCertificate
} = ctx

function taxCertForItem(itemId: number) {
  return (taxCertificates.value ?? []).find(
    (c: any) => c.schedule_item_id === itemId
  )
}

const canScreenSanctions = computed(() =>
  hasPermission('claims_pay:screen_sanctions')
)
const canEditReinsurance = computed(() =>
  hasPermission('claims_pay:finance_review')
)
const canConfirmRecovery = computed(() =>
  hasPermission('claims_pay:recovery_confirm')
)

const REASON_CODES = [
  { title: 'Banking details mismatch', value: 'BANKING_DETAILS_MISMATCH' },
  {
    title: 'Missing call-back evidence on changed banking',
    value: 'NO_CALLBACK_EVIDENCE'
  },
  { title: 'ID number does not match claimant', value: 'ID_MISMATCH' },
  {
    title: 'Deduction amount disputed (premium/loan/tax)',
    value: 'DEDUCTION_MISMATCH'
  },
  { title: 'Sanctions / PEP hit', value: 'SANCTIONS_HIT' },
  { title: 'Duplicate payment', value: 'DUPLICATE' },
  { title: 'Other (see notes)', value: 'OTHER' }
]

const isFinanceReview = computed(
  () =>
    schedule.value?.status === 'finance_in_review' &&
    hasPermission('claims_pay:finance_review')
)

const canGenerateLetter = computed(
  () =>
    (schedule.value?.status === 'confirmed' ||
      schedule.value?.status === 'archived') &&
    hasPermission('claims_pay:generate_letter')
)

const letterDialog = ref(false)
const letterClaimId = ref<number | null>(null)
const letterClaimContact = ref<{ email: string; phone: string }>({
  email: '',
  phone: ''
})

function openLetterDialog(item: ScheduleItem) {
  letterClaimId.value = item.claim_id
  // Schedule items hold the bank/beneficiary snapshot, not the claim contact
  // details. The dialog falls back to "no recipient" — the user can type one
  // or jump out to edit the claim record.
  letterClaimContact.value = { email: '', phone: '' }
  letterDialog.value = true
}

function parseRiskFlags(item: ScheduleItem): RiskFlags {
  const raw = item.risk_flags
  if (!raw) return {}
  if (typeof raw === 'string') {
    try {
      return JSON.parse(raw)
    } catch {
      return {}
    }
  }
  return raw
}

function riskFlag(item: ScheduleItem, key: keyof RiskFlags): boolean {
  const f = parseRiskFlags(item)
  return Boolean(f[key])
}

function fraudLevel(item: ScheduleItem): string {
  return parseRiskFlags(item).fraud_risk_level ?? ''
}

function fraudColor(level: string) {
  const map: Record<string, string> = {
    low: 'success',
    medium: 'warning',
    high: 'error',
    critical: 'red-darken-4'
  }
  return map[level?.toLowerCase()] ?? 'grey'
}

function totalDeductions(item: ScheduleItem) {
  return (
    (item.premium_arrears_deduction ?? 0) +
    (item.policy_loan_deduction ?? 0) +
    (item.tax_withheld ?? 0)
  )
}

function beneficiaryMismatch(item: ScheduleItem) {
  if (!item.beneficiary_id_number || !item.member_id_number) return false
  return item.beneficiary_id_number !== item.member_id_number
}

function lineStatusColor(status?: string) {
  const map: Record<string, string> = {
    pending: 'grey',
    verified: 'success',
    queried: 'warning',
    rejected: 'error'
  }
  return map[status ?? 'pending'] ?? 'default'
}

function canVerify(item: ScheduleItem) {
  if (!isFinanceReview.value) return false
  return !item.line_status || item.line_status === 'pending'
}

function canQuery(item: ScheduleItem) {
  if (!isFinanceReview.value) return false
  return (
    !item.line_status ||
    item.line_status === 'pending' ||
    item.line_status === 'verified'
  )
}

// Query / Reject dialog
const queryDialog = ref(false)
const dialogMode = ref<'query' | 'reject'>('query')
const dialogItemId = ref<number | null>(null)
const reasonCode = ref('')
const reasonNotes = ref('')
const submittingQuery = ref(false)

function openQueryDialog(item: ScheduleItem, mode: 'query' | 'reject') {
  dialogItemId.value = item.id
  dialogMode.value = mode
  reasonCode.value = ''
  reasonNotes.value = ''
  queryDialog.value = true
}

async function submitQuery() {
  if (!dialogItemId.value || !reasonCode.value) return
  submittingQuery.value = true
  try {
    if (dialogMode.value === 'reject') {
      await rejectLineItem(
        dialogItemId.value,
        reasonCode.value,
        reasonNotes.value
      )
    } else {
      await queryLineItem(
        dialogItemId.value,
        reasonCode.value,
        reasonNotes.value
      )
    }
    queryDialog.value = false
  } finally {
    submittingQuery.value = false
  }
}

// ── Phase 3: sanctions chip / dialog ─────────────────────────
function latestSanctionsForItem(itemId: number) {
  return (sanctions.value ?? [])
    .filter((s: any) => s.schedule_item_id === itemId)
    .sort(
      (a: any, b: any) =>
        new Date(b.updated_at || b.created_at).getTime() -
        new Date(a.updated_at || a.created_at).getTime()
    )[0]
}

function sanctionsLabel(itemId: number) {
  const row = latestSanctionsForItem(itemId)
  if (!row) return 'Not screened'
  if (row.status === 'manual_clear') return 'Manual clear'
  return row.status
}

function sanctionsChipColor(itemId: number) {
  const row = latestSanctionsForItem(itemId)
  const map: Record<string, string> = {
    clear: 'success',
    manual_clear: 'success',
    pending: 'orange',
    hit: 'error',
    skipped: 'warning'
  }
  if (!row) return 'grey'
  return map[row.status] ?? 'default'
}

function sanctionsHitSummary(itemId: number) {
  const row = latestSanctionsForItem(itemId)
  return row?.hit_summary || row?.notes || ''
}

const dialogItem = ref<ScheduleItem | null>(null)

const sanctionsDialog = ref(false)
const sanctionsStatus = ref('')
const sanctionsNotes = ref('')
const savingSanctions = ref(false)

function openSanctionsDialog(item: ScheduleItem) {
  dialogItem.value = item
  const existing = latestSanctionsForItem(item.id)
  sanctionsStatus.value =
    existing?.status === 'pending' ? '' : (existing?.status ?? '')
  sanctionsNotes.value = existing?.notes ?? ''
  sanctionsDialog.value = true
}

async function submitSanctions() {
  if (!dialogItem.value || !sanctionsStatus.value) return
  savingSanctions.value = true
  try {
    await recordSanctionsOutcome(
      dialogItem.value.id,
      sanctionsStatus.value,
      sanctionsNotes.value
    )
    sanctionsDialog.value = false
  } finally {
    savingSanctions.value = false
  }
}

async function onScreen(itemId: number) {
  await screenLineItem(itemId)
}

// ── Phase 3: reinsurance dialog ──────────────────────────────
const reinsuranceDialog = ref(false)
const reinsuranceRequired = ref(false)
const reinsuranceAmount = ref(0)
const savingReinsurance = ref(false)

function openReinsuranceDialog(item: ScheduleItem) {
  dialogItem.value = item
  reinsuranceRequired.value = Boolean(item.reinsurance_recovery_required)
  reinsuranceAmount.value = item.reinsurance_recovery_amount ?? 0
  reinsuranceDialog.value = true
}

async function submitReinsurance() {
  if (!dialogItem.value) return
  savingReinsurance.value = true
  try {
    await setReinsuranceRecovery(
      dialogItem.value.id,
      reinsuranceRequired.value,
      reinsuranceAmount.value
    )
    reinsuranceDialog.value = false
  } finally {
    savingReinsurance.value = false
  }
}

async function onConfirmRaised(itemId: number) {
  await confirmReinsuranceRaised(itemId)
}

async function onClearDuplicate(itemId: number) {
  await clearDuplicateBeneficiary(itemId)
}
</script>

<style scoped>
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

.totals-card {
  min-height: 90px;
  transition: border-color 0.15s ease;
}

.totals-card:hover {
  border-color: rgba(var(--v-theme-primary), 0.4);
}

.totals-card--accent {
  border-color: rgba(var(--v-theme-primary), 0.55);
}

.totals-card__label {
  font-size: 0.7rem;
  font-weight: 500;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  color: rgba(var(--v-theme-on-surface), 0.6);
  margin-bottom: 4px;
}

.totals-card__value {
  font-size: 1.35rem;
  font-weight: 700;
  line-height: 1.2;
  color: rgba(var(--v-theme-on-surface), 0.95);
}

.totals-card__value--primary {
  color: rgb(var(--v-theme-primary));
}
</style>
