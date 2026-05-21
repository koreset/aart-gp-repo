<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div
              class="d-flex justify-space-between align-center flex-wrap gap-2"
            >
              <div class="d-flex align-center">
                <v-btn
                  icon="mdi-arrow-left"
                  variant="text"
                  class="mr-2"
                  @click="goBack"
                />
                <span class="headline">Payment Schedules</span>
              </div>
              <div class="d-flex ga-3 align-center flex-wrap">
                <v-chip
                  v-if="nextCutoff"
                  color="indigo"
                  size="small"
                  variant="tonal"
                  prepend-icon="mdi-clock-outline"
                  :title="`Next auto cut-off: ${nextCutoff}`"
                >
                  Next cut-off {{ nextCutoff }}
                </v-chip>
                <v-btn
                  v-if="hasPermission('claims_pay:run_cutoff')"
                  rounded
                  size="small"
                  variant="outlined"
                  prepend-icon="mdi-play-circle-outline"
                  :loading="runningCutoff"
                  @click="runCutoffNow"
                >
                  Run cut-off now
                </v-btn>
                <v-btn
                  v-if="hasPermission('claims_pay:admin_cutoff')"
                  rounded
                  size="small"
                  variant="outlined"
                  prepend-icon="mdi-cog-outline"
                  :to="{ name: 'group-pricing-payment-cutoff-settings' }"
                >
                  Cut-off settings
                </v-btn>
                <v-switch
                  v-model="showArchived"
                  inset
                  color="success"
                  base-color="grey-lighten-1"
                  density="compact"
                  hide-details
                  label="Show archived"
                  @update:model-value="loadSchedules"
                />
                <v-btn
                  v-if="hasPermission('claims_pay:create_schedule')"
                  rounded
                  size="small"
                  color="primary"
                  variant="flat"
                  prepend-icon="mdi-plus"
                  @click="openCreateDialog"
                >
                  New Payment Schedule
                </v-btn>
              </div>
            </div>
          </template>

          <template #default>
            <!-- ── Section: Payment lifecycle (static — educational) ── -->
            <section class="page-section">
              <div class="section-header">
                <span class="section-label">Payment Lifecycle</span>
                <span class="section-divider" />
              </div>
              <div class="workflow-stepper">
                <div
                  v-for="(step, idx) in pipelineSteps"
                  :key="step.status"
                  class="workflow-stepper__step"
                  :class="
                    isClaimsStep(step.status)
                      ? 'workflow-stepper__step--current'
                      : 'workflow-stepper__step--muted'
                  "
                >
                  <div class="workflow-stepper__num">{{ idx + 1 }}</div>
                  <div class="workflow-stepper__text">
                    <div class="workflow-stepper__label">
                      {{ step.label }}
                    </div>
                    <div class="workflow-stepper__sub d-flex align-center ga-1">
                      <v-icon
                        v-if="!isClaimsStep(step.status)"
                        size="12"
                        icon="mdi-bank"
                      />
                      <span>
                        {{ isClaimsStep(step.status) ? step.sub : 'Finance' }}
                      </span>
                    </div>
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

            <section class="page-section">
              <div class="section-header">
                <span class="section-label">Schedules</span>
                <v-chip
                  v-if="schedules.length > 0"
                  size="x-small"
                  variant="tonal"
                  color="primary"
                  class="ml-2"
                >
                  {{ schedules.length }}
                </v-chip>
                <span class="section-divider" />
              </div>

              <v-progress-linear
                v-if="loading"
                indeterminate
                color="primary"
                class="mb-2"
              />

              <empty-state
                v-if="!loading && schedules.length === 0"
                icon="mdi-cash-multiple"
                title="No payment schedules yet"
                message="Generate a payment schedule from approved claims to submit it to finance."
                action-label="New payment schedule"
                :action-fn="openCreateDialog"
              />

              <v-data-table
                v-if="schedules.length > 0"
                :headers="headers"
                :items="schedules"
                :loading="loading"
                density="compact"
                hover
                @click:row="(_e: any, row: any) => openDrawer(row.item)"
              >
                <template #[`item.status`]="{ item }">
                  <v-chip
                    :color="statusColor(item.status)"
                    size="x-small"
                    label
                    variant="flat"
                  >
                    {{ statusLabel(item.status) }}
                  </v-chip>
                </template>

                <template #[`item.total_amount`]="{ item }">
                  {{ formatCurrency(item.total_amount) }}
                </template>

                <template #[`item.created_at`]="{ item }">
                  {{ formatDate(item.created_at) }}
                </template>

                <template #[`item.actions`]="{ item }">
                  <div class="d-flex ga-1 justify-end">
                    <v-btn
                      size="x-small"
                      variant="text"
                      prepend-icon="mdi-format-list-bulleted"
                      @click.stop="openClaimsList(item)"
                    >
                      Claims
                    </v-btn>
                    <v-btn
                      size="x-small"
                      variant="text"
                      prepend-icon="mdi-eye-outline"
                      @click.stop="openDrawer(item)"
                    >
                      View
                    </v-btn>
                  </div>
                </template>
              </v-data-table>
            </section>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Detail drawer (read-only) -->
    <v-navigation-drawer
      v-model="drawer"
      location="right"
      temporary
      width="640"
    >
      <template v-if="active">
        <div class="pa-4">
          <div class="d-flex justify-space-between align-center mb-3">
            <div>
              <div class="text-overline text-medium-emphasis">Schedule</div>
              <div class="text-h6">{{ active.schedule_number }}</div>
            </div>
            <v-btn icon="mdi-close" variant="text" @click="drawer = false" />
          </div>

          <div class="d-flex align-center flex-wrap gap-2 mb-3">
            <v-chip
              :color="statusColor(active.status)"
              size="small"
              label
              variant="flat"
            >
              {{ statusLabel(active.status) }}
            </v-chip>
            <v-btn
              v-if="
                active.status === 'draft' &&
                hasPermission('claims_pay:signoff_schedule')
              "
              size="small"
              color="primary"
              variant="flat"
              rounded
              prepend-icon="mdi-clipboard-check-outline"
              :loading="signingOff"
              @click="signOff"
            >
              Sign Off
            </v-btn>
            <v-btn
              v-if="
                active.status === 'draft' &&
                hasPermission('claims_pay:create_schedule')
              "
              size="small"
              color="error"
              variant="outlined"
              rounded
              prepend-icon="mdi-trash-can-outline"
              :loading="discarding"
              @click="confirmDiscard = true"
            >
              Discard draft
            </v-btn>
          </div>

          <div class="workflow-stepper workflow-stepper--vertical mb-4">
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
                <div class="workflow-stepper__label">{{ step.label }}</div>
                <div class="workflow-stepper__sub">{{ step.sub }}</div>
              </div>
            </div>
          </div>

          <v-table density="compact" class="border rounded mb-4">
            <tbody>
              <tr>
                <th class="text-left">Claims</th>
                <td>{{ active.claims_count }}</td>
              </tr>
              <tr>
                <th class="text-left">Total amount</th>
                <td
                  ><strong>{{
                    formatCurrency(active.total_amount)
                  }}</strong></td
                >
              </tr>
              <tr>
                <th class="text-left">Created by</th>
                <td>{{ active.created_by || '—' }}</td>
              </tr>
              <tr>
                <th class="text-left">Created at</th>
                <td>{{ formatDate(active.created_at) }}</td>
              </tr>
              <tr v-if="active.acb_generated_by">
                <th class="text-left">ACB generated</th>
                <td>
                  {{ active.acb_generated_by }}
                  <span class="text-medium-emphasis text-caption">
                    · {{ formatDate(active.acb_generated_at) }}
                  </span>
                </td>
              </tr>
            </tbody>
          </v-table>

          <div class="d-flex align-center mb-2">
            <div class="text-subtitle-2">Queries & follow-ups</div>
            <v-spacer />
            <v-btn
              size="x-small"
              variant="text"
              prepend-icon="mdi-refresh"
              :loading="loadingQueries"
              @click="loadQueries"
            >
              Refresh
            </v-btn>
          </div>

          <v-card variant="outlined" rounded="lg" class="pa-3 mb-3">
            <empty-state
              v-if="!loadingQueries && queries.length === 0"
              icon="mdi-comment-check-outline"
              title="No queries yet"
              message="Finance has not raised any queries on this submission."
            />
            <v-list v-else density="compact">
              <v-list-item
                v-for="q in queries"
                :key="q.id"
                class="px-0"
                lines="two"
              >
                <template #prepend>
                  <v-chip
                    :color="outcomeColor(q.outcome)"
                    size="x-small"
                    variant="flat"
                    label
                  >
                    {{ q.outcome }}
                  </v-chip>
                </template>
                <v-list-item-title class="text-body-2">
                  <span class="font-weight-medium">{{ q.reason_code }}</span>
                  <span v-if="q.claim_number" class="text-medium-emphasis">
                    · {{ q.claim_number }}
                  </span>
                </v-list-item-title>
                <v-list-item-subtitle class="text-caption">
                  {{ q.notes }}
                  <span class="text-medium-emphasis">
                    · {{ q.raised_by }} on {{ formatDate(q.raised_at) }}
                  </span>
                  <div
                    v-if="q.resolution_notes"
                    class="text-caption text-success mt-1"
                  >
                    <v-icon size="12" icon="mdi-reply" class="mr-1" />
                    <strong>{{ q.resolved_by }}:</strong>
                    {{ q.resolution_notes }}
                    <span v-if="q.resolved_at" class="text-medium-emphasis">
                      · {{ formatDate(q.resolved_at) }}
                    </span>
                  </div>
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card>

          <div class="text-subtitle-2 mb-2">Post a follow-up to finance</div>
          <v-textarea
            v-model="followupNotes"
            variant="outlined"
            density="compact"
            rows="3"
            placeholder="e.g. Please confirm expected payment date for this run."
            :disabled="!canFollowup"
            :hint="
              canFollowup
                ? ''
                : 'Follow-ups can be posted once the schedule has been signed off.'
            "
            persistent-hint
          />
          <div class="d-flex justify-end mt-2">
            <v-btn
              color="primary"
              size="small"
              :loading="postingFollowup"
              :disabled="!followupNotes.trim() || !canFollowup"
              @click="postFollowup"
            >
              Send follow-up
            </v-btn>
          </div>
        </div>
      </template>
    </v-navigation-drawer>

    <!-- Shared create-schedule dialog -->
    <CreatePaymentScheduleDialog
      v-model="createDialog"
      @created="onScheduleCreated"
      @error="(msg: string) => notify(msg, 'error')"
    />

    <!-- Discard draft confirmation -->
    <v-dialog v-model="confirmDiscard" max-width="460px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2">Discard draft?</v-card-title>
        <v-card-text>
          <p class="text-body-2 mb-2">
            This will permanently delete draft schedule
            <strong>{{ active?.schedule_number }}</strong> and return its claims
            to the approved queue for the next cut-off.
          </p>
          <p class="text-body-2 text-medium-emphasis">
            This cannot be undone.
          </p>
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="confirmDiscard = false">Cancel</v-btn>
          <v-btn color="error" :loading="discarding" @click="discardDraft">
            Discard draft
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar" :color="snackbarColor" timeout="3500">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'
import {
  useScheduleLifecycle,
  isClaimsStep
} from '@/renderer/composables/useScheduleLifecycle'
import CreatePaymentScheduleDialog from './CreatePaymentScheduleDialog.vue'

interface PaymentSchedule {
  id: number
  schedule_number: string
  status: string
  total_amount: number
  claims_count: number
  created_by: string
  created_at: string
  acb_generated_at?: string
  acb_generated_by?: string
}

interface ScheduleQuery {
  id: number
  schedule_id: number
  schedule_item_id: number
  claim_number: string
  reason_code: string
  notes: string
  outcome: string
  raised_by: string
  raised_at: string
  resolution_notes?: string
  resolved_by?: string
  resolved_at?: string | null
}

const router = useRouter()
const { hasPermission } = usePermissionCheck()

const loading = ref(false)
const schedules = ref<PaymentSchedule[]>([])
const showArchived = ref(false)

const createDialog = ref(false)
const runningCutoff = ref(false)
const nextCutoff = ref<string | null>(null)

const drawer = ref(false)
const active = ref<PaymentSchedule | null>(null)
const queries = ref<ScheduleQuery[]>([])
const loadingQueries = ref(false)
const followupNotes = ref('')
const postingFollowup = ref(false)
const signingOff = ref(false)
const discarding = ref(false)
const confirmDiscard = ref(false)

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

const activeStatus = computed(() => active.value?.status ?? null)
const { pipelineSteps, currentStepIndex, stepTone, statusLabel } =
  useScheduleLifecycle(activeStatus)

const canFollowup = computed(() => {
  if (!active.value) return false
  return active.value.status !== 'draft'
})

const headers = [
  { title: 'Schedule', key: 'schedule_number', sortable: true },
  { title: 'Status', key: 'status', sortable: true },
  {
    title: 'Claims',
    key: 'claims_count',
    sortable: true,
    align: 'end' as const
  },
  {
    title: 'Total',
    key: 'total_amount',
    sortable: true,
    align: 'end' as const
  },
  { title: 'Submitted by', key: 'created_by', sortable: true },
  { title: 'Submitted at', key: 'created_at', sortable: true },
  { title: '', key: 'actions', sortable: false, align: 'end' as const }
]

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

function statusColor(status: string) {
  const map: Record<string, string> = {
    draft: 'grey',
    claims_signed_off: 'info',
    finance_in_review: 'info',
    finance_first_authorised: 'warning',
    finance_second_authorised: 'warning',
    submitted_to_bank: 'warning',
    confirmed: 'success',
    archived: 'default'
  }
  return map[status] ?? 'default'
}

function outcomeColor(outcome: string) {
  const map: Record<string, string> = {
    open: 'orange',
    queried: 'orange',
    rejected: 'error',
    resolved: 'success',
    cancelled: 'grey'
  }
  return map[outcome] ?? 'default'
}

function notify(message: string, color: string = 'success') {
  snackbarMessage.value = message
  snackbarColor.value = color
  snackbar.value = true
}

function unwrap(res: any) {
  return res?.data?.data ?? res?.data ?? res
}

async function loadSchedules() {
  loading.value = true
  try {
    const res = await GroupPricingService.getPaymentSchedules(
      showArchived.value,
      'claims'
    )
    schedules.value = unwrap(res) ?? []
  } catch (e) {
    notify('Failed to load submissions', 'error')
  } finally {
    loading.value = false
  }
}

async function loadQueries() {
  if (!active.value) return
  loadingQueries.value = true
  try {
    const res = await GroupPricingService.getScheduleQueries(active.value.id)
    queries.value = unwrap(res) ?? []
  } catch (e) {
    notify('Failed to load queries', 'error')
  } finally {
    loadingQueries.value = false
  }
}

function openDrawer(schedule: PaymentSchedule) {
  active.value = schedule
  followupNotes.value = ''
  queries.value = []
  drawer.value = true
  loadQueries()
}

function openClaimsList(schedule: PaymentSchedule) {
  router.push({
    name: 'group-pricing-claim-my-submissions-claims',
    params: { scheduleId: String(schedule.id) }
  })
}

function openCreateDialog() {
  createDialog.value = true
}

async function onScheduleCreated() {
  notify('Payment schedule created. Submitted to your draft queue.')
  await loadSchedules()
}

async function loadNextCutoff() {
  try {
    const res = await GroupPricingService.getNextPaymentCutoff()
    const data = unwrap(res)
    nextCutoff.value = data?.next_cutoff ?? data?.formatted ?? null
  } catch {
    nextCutoff.value = null
  }
}

async function runCutoffNow() {
  runningCutoff.value = true
  try {
    await GroupPricingService.runPaymentCutoffNow()
    notify('Cut-off run started. Schedules will refresh shortly.')
    await loadSchedules()
  } catch (e: any) {
    const msg = e?.response?.data?.error || 'Failed to run cut-off'
    notify(typeof msg === 'string' ? msg : 'Failed to run cut-off', 'error')
  } finally {
    runningCutoff.value = false
  }
}

function goBack() {
  router.push({ name: 'group-pricing-claims-management' })
}

async function discardDraft() {
  if (!active.value) return
  discarding.value = true
  try {
    await GroupPricingService.discardPaymentSchedule(active.value.id)
    notify('Draft discarded. Claims returned to the approved queue.')
    confirmDiscard.value = false
    drawer.value = false
    active.value = null
    await loadSchedules()
  } catch (e: any) {
    const msg = e?.response?.data?.error || 'Failed to discard draft'
    notify(typeof msg === 'string' ? msg : 'Failed to discard draft', 'error')
  } finally {
    discarding.value = false
  }
}

async function signOff() {
  if (!active.value) return
  signingOff.value = true
  try {
    const res = await GroupPricingService.signOffPaymentSchedule(
      active.value.id
    )
    const updated = unwrap(res)
    if (updated) {
      active.value = { ...active.value, ...updated }
    }
    notify('Payment schedule signed off and submitted to finance.')
    await loadSchedules()
  } catch (e: any) {
    const msg = e?.response?.data?.error || 'Failed to sign off schedule'
    notify(
      typeof msg === 'string' ? msg : 'Failed to sign off schedule',
      'error'
    )
  } finally {
    signingOff.value = false
  }
}

async function postFollowup() {
  if (!active.value || !followupNotes.value.trim()) return
  postingFollowup.value = true
  try {
    await GroupPricingService.postScheduleFollowup(
      active.value.id,
      followupNotes.value.trim()
    )
    followupNotes.value = ''
    notify('Follow-up sent to finance.')
    await loadQueries()
  } catch (e: any) {
    const msg = e?.response?.data?.error || 'Failed to send follow-up'
    notify(msg, 'error')
  } finally {
    postingFollowup.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadSchedules(), loadNextCutoff()])
})
</script>

<style scoped>
/* ── Page sections ────────────────────────────────────── */
.page-section {
  margin-bottom: 28px;
}
.page-section--last {
  margin-bottom: 0;
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

/* ── Workflow stepper ─────────────────────────────────── */
.workflow-stepper {
  display: flex;
  align-items: stretch;
  flex-wrap: nowrap;
  gap: 8px;
  padding: 14px 16px;
  background: rgba(var(--v-theme-primary), 0.04);
  border: 1px solid rgba(var(--v-theme-primary), 0.1);
  border-radius: 10px;
  overflow-x: auto;
  scrollbar-width: thin;
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

.workflow-stepper--vertical {
  flex-direction: column;
  align-items: stretch;
  background: transparent;
  border: none;
  padding: 0;
  gap: 6px;
}
.workflow-stepper--vertical .workflow-stepper__step {
  min-width: 0;
}
</style>
