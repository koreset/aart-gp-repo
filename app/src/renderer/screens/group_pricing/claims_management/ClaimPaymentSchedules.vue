<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex justify-space-between align-center flex-wrap gap-2">
              <div class="d-flex align-center">
                <v-btn
                  icon="mdi-arrow-left"
                  variant="text"
                  class="mr-2"
                  @click="goBack"
                />
                <span class="headline">Payment Schedules</span>
              </div>
              <div class="d-flex gap-2">
                <v-btn
                  v-if="hasPermission('claims_pay:manage_bank_profiles')"
                  rounded
                  size="small"
                  variant="outlined"
                  class="mr-2"
                  prepend-icon="mdi-bank"
                  @click="openBankProfilesDialog"
                >
                  Bank Profiles
                </v-btn>
                <v-btn
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
            <!-- ── Section: Overview ───────────────────────── -->
            <section class="page-section">
              <div class="section-header">
                <span class="section-label">Overview</span>
                <span class="section-divider" />
              </div>

              <v-row class="mb-2">
                <v-col cols="12" sm="6" md="3">
                  <div class="kpi-card kpi-card--primary">
                    <div class="kpi-card__stripe" />
                    <div class="kpi-card__body">
                      <div class="kpi-card__head">
                        <span class="kpi-card__label">Active Schedules</span>
                        <v-icon class="kpi-card__icon" size="20"
                          >mdi-clipboard-clock-outline</v-icon
                        >
                      </div>
                      <div class="kpi-card__value">
                        {{ paymentStats.activeCount }}
                      </div>
                      <div class="kpi-card__hint"
                        >Draft or submitted, awaiting confirmation</div
                      >
                    </div>
                  </div>
                </v-col>
                <v-col cols="12" sm="6" md="3">
                  <div class="kpi-card kpi-card--accent">
                    <div class="kpi-card__stripe" />
                    <div class="kpi-card__body">
                      <div class="kpi-card__head">
                        <span class="kpi-card__label"
                          >Pending Payment Value</span
                        >
                        <v-icon class="kpi-card__icon" size="20"
                          >mdi-cash-clock</v-icon
                        >
                      </div>
                      <div class="kpi-card__value">
                        {{ formatCurrency(paymentStats.pendingValue) }}
                      </div>
                      <div class="kpi-card__hint"
                        >Total awaiting confirmation</div
                      >
                    </div>
                  </div>
                </v-col>
                <v-col cols="12" sm="6" md="3">
                  <div class="kpi-card kpi-card--success">
                    <div class="kpi-card__stripe" />
                    <div class="kpi-card__body">
                      <div class="kpi-card__head">
                        <span class="kpi-card__label">Confirmed This Month</span>
                        <v-icon class="kpi-card__icon" size="20"
                          >mdi-check-decagram-outline</v-icon
                        >
                      </div>
                      <div class="kpi-card__value">
                        {{ paymentStats.confirmedThisMonth }}
                      </div>
                      <div class="kpi-card__hint"
                        >Schedules paid in current month</div
                      >
                    </div>
                  </div>
                </v-col>
                <v-col cols="12" sm="6" md="3">
                  <div
                    class="kpi-card"
                    :class="
                      paymentStats.missingBank > 0
                        ? 'kpi-card--warning'
                        : 'kpi-card--muted'
                    "
                  >
                    <div class="kpi-card__stripe" />
                    <div class="kpi-card__body">
                      <div class="kpi-card__head">
                        <span class="kpi-card__label">Items Missing Bank</span>
                        <v-icon class="kpi-card__icon" size="20">{{
                          paymentStats.missingBank > 0
                            ? 'mdi-bank-off-outline'
                            : 'mdi-bank-check-outline'
                        }}</v-icon>
                      </div>
                      <div class="kpi-card__value">
                        {{ paymentStats.missingBank }}
                      </div>
                      <div class="kpi-card__hint">{{
                        paymentStats.missingBank > 0
                          ? 'Block ACB generation until resolved'
                          : 'All schedule items have bank details'
                      }}</div>
                    </div>
                  </div>
                </v-col>
              </v-row>
            </section>

            <!-- ── Section: Workflow ───────────────────────── -->
            <section class="page-section">
              <div class="section-header">
                <span class="section-label">Payment Lifecycle</span>
                <span class="section-divider" />
              </div>

              <div class="workflow-stepper">
                <div
                  v-for="(step, idx) in workflowSteps"
                  :key="step.label"
                  class="workflow-stepper__step"
                  :class="`workflow-stepper__step--${step.tone}`"
                >
                  <div class="workflow-stepper__num">{{ idx + 1 }}</div>
                  <div class="workflow-stepper__text">
                    <div class="workflow-stepper__label">{{ step.label }}</div>
                    <div class="workflow-stepper__sub">{{ step.sub }}</div>
                  </div>
                  <v-icon
                    v-if="idx < workflowSteps.length - 1"
                    class="workflow-stepper__arrow"
                    size="18"
                    >mdi-chevron-right</v-icon
                  >
                </div>
              </div>
            </section>

            <!-- ── Section: Schedules ──────────────────────── -->
            <section class="page-section page-section--last">
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
                message="Create a payment schedule from approved claims to begin a payment run."
                action-label="Create your first payment schedule"
                :action-fn="openCreateDialog"
              />

              <template v-if="schedules.length > 0">
                <!-- Schedule card strip -->
                <div class="schedule-card-strip mb-4">
                  <div class="schedule-card-strip__inner">
                    <div
                      v-for="schedule in schedules"
                      :key="'card-' + schedule.id"
                      class="schedule-card-strip__item"
                    >
                      <div
                        class="ps-card"
                        :class="`ps-card--${statusColor(schedule.status)}`"
                        @click="openSchedule(schedule)"
                      >
                        <div class="ps-card__stripe" />
                        <div class="ps-card__body">
                          <div class="ps-card__head">
                            <span class="ps-card__num">{{
                              schedule.schedule_number
                            }}</span>
                            <div class="d-flex gap-1">
                              <v-chip
                                v-if="schedule.acb_file_generated"
                                color="accent"
                                size="x-small"
                                label
                                variant="flat"
                              >
                                ACB
                              </v-chip>
                              <v-chip
                                :color="statusColor(schedule.status)"
                                size="x-small"
                                label
                              >
                                {{ statusLabel(schedule.status) }}
                              </v-chip>
                            </div>
                          </div>
                          <div class="ps-card__metrics">
                            <div>
                              <div class="ps-card__metric-label">Claims</div>
                              <div class="ps-card__metric-value">{{
                                schedule.claims_count
                              }}</div>
                            </div>
                            <div class="text-right">
                              <div class="ps-card__metric-label">Total</div>
                              <div class="ps-card__metric-value">{{
                                formatCurrency(schedule.total_amount)
                              }}</div>
                            </div>
                          </div>
                          <div class="ps-card__footer">
                            <span>{{ formatDate(schedule.created_at) }}</span>
                            <span
                              v-if="schedule.created_by"
                              class="ps-card__footer-divider"
                              >·</span
                            >
                            <span v-if="schedule.created_by">{{
                              schedule.created_by
                            }}</span>
                          </div>
                          <v-chip
                            v-if="itemsMissingBanking(schedule).length > 0"
                            color="warning"
                            size="x-small"
                            variant="tonal"
                            class="mt-2"
                            prepend-icon="mdi-alert-outline"
                          >
                            {{ itemsMissingBanking(schedule).length }} missing
                            bank
                          </v-chip>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Schedules data grid in subtle panel -->
                <div class="grid-panel">
                  <DataGrid
                    :row-data="schedules"
                    :column-defs="scheduleColumnDefs"
                    density="compact"
                    :pagination="true"
                    :pagination-page-size="20"
                    @row-double-clicked="onScheduleRowDoubleClicked"
                  />
                </div>
                <div class="text-caption text-medium-emphasis mt-2">
                  <v-icon size="14" class="mr-1"
                    >mdi-information-outline</v-icon
                  >
                  Double-click a row, click the View action, or click any card
                  above to open the schedule.
                </div>
              </template>
            </section>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- ── Create Schedule Dialog ── -->
    <v-dialog v-model="createDialog" persistent max-width="800px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2"
          >New Payment Schedule</v-card-title
        >
        <v-card-text>
          <v-text-field
            v-model="newScheduleDescription"
            label="Description (optional)"
            variant="outlined"
            density="compact"
            class="mb-4"
          />

          <div class="d-flex align-center justify-space-between mb-2">
            <span class="text-subtitle-2">Select Approved Claims</span>
            <v-chip size="small" color="primary" variant="tonal">
              {{ selectedClaimIDs.length }} selected
            </v-chip>
          </div>
          <v-row dense class="mb-2">
            <v-col cols="6">
              <v-text-field
                v-model="claimFilter"
                label="Filter by Claim Number"
                prepend-inner-icon="mdi-magnify"
                variant="outlined"
                density="compact"
                clearable
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model="benefitFilter"
                label="Filter by Benefit"
                prepend-inner-icon="mdi-filter-outline"
                variant="outlined"
                density="compact"
                clearable
              />
            </v-col>
          </v-row>
          <DataGrid
            :row-data="filteredApprovedClaims"
            :column-defs="claimsColumnDefs"
            row-selection="multiple"
            density="compact"
            :pagination="false"
            @row-selection-changed="onClaimSelectionChanged"
          />
          <div class="text-subtitle-2 text-right mt-2">
            Subtotal: <strong>{{ formatCurrency(selectedTotal) }}</strong>
          </div>
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="createDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="creating"
            :disabled="selectedClaimIDs.length === 0"
            @click="createSchedule"
          >
            Create Schedule
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- ── Bank Profiles Management Dialog ── -->
    <v-dialog v-model="profilesDialog" max-width="700px">
      <v-card rounded="lg">
        <v-card-title
          class="text-h6 pa-4 pb-2 d-flex justify-space-between align-center"
        >
          Bank Profiles
          <v-btn
            size="small"
            color="primary"
            prepend-icon="mdi-plus"
            @click="openCreateProfileDialog"
          >
            New Profile
          </v-btn>
        </v-card-title>
        <v-card-text>
          <v-progress-linear
            v-if="loadingProfiles"
            indeterminate
            color="teal"
            class="mb-2"
          />
          <v-alert
            v-if="!loadingProfiles && bankProfiles.length === 0"
            type="info"
            variant="tonal"
            density="compact"
          >
            No bank profiles configured yet.
          </v-alert>
          <v-list v-else density="compact" class="border rounded">
            <v-list-item
              v-for="profile in bankProfiles"
              :key="profile.id"
              :subtitle="`${profile.bank_name} | Account: ${profile.user_account_number} | Gen #${profile.generation_number}`"
              :title="profile.profile_name"
            >
              <template #prepend>
                <v-icon color="teal">mdi-bank</v-icon>
              </template>
              <template #append>
                <v-chip
                  :color="profile.is_active ? 'success' : 'grey'"
                  size="x-small"
                  label
                  class="mr-2"
                >
                  {{ profile.is_active ? 'Active' : 'Inactive' }}
                </v-chip>
                <v-btn
                  size="x-small"
                  variant="text"
                  color="primary"
                  icon="mdi-pencil"
                  @click="editProfile(profile)"
                />
                <v-btn
                  size="x-small"
                  variant="text"
                  color="error"
                  icon="mdi-delete"
                  @click="deleteProfile(profile)"
                />
              </template>
            </v-list-item>
          </v-list>
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="profilesDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- ── Create/Edit Profile Dialog ── -->
    <v-dialog v-model="profileFormDialog" persistent max-width="550px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2">
          {{ editingProfile ? 'Edit Bank Profile' : 'New Bank Profile' }}
        </v-card-title>
        <v-card-text>
          <v-text-field
            v-model="profileForm.profile_name"
            label="Profile Name *"
            variant="outlined"
            density="compact"
            class="mb-2"
          />
          <v-select
            v-model="profileForm.bank_name"
            :items="bankNameOptions"
            label="Bank *"
            variant="outlined"
            density="compact"
            class="mb-2"
            @update:model-value="onProfileBankSelected"
          />
          <v-text-field
            v-model="profileForm.user_code"
            label="BankServ User Code *"
            variant="outlined"
            density="compact"
            class="mb-2"
            hint="4-character code assigned by the bank"
            persistent-hint
          />
          <v-row dense>
            <v-col cols="6">
              <v-text-field
                v-model="profileForm.user_branch_code"
                label="Source Branch Code"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model="profileForm.user_account_number"
                label="Source Account Number *"
                variant="outlined"
                density="compact"
              />
            </v-col>
          </v-row>
          <v-row dense>
            <v-col cols="6">
              <v-select
                v-model="profileForm.user_account_type"
                :items="accountTypeOptions"
                label="Account Type"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="6">
              <v-select
                v-model="profileForm.service_type"
                :items="serviceTypeOptions"
                label="Service Type"
                variant="outlined"
                density="compact"
              />
            </v-col>
          </v-row>
          <v-text-field
            v-model="profileForm.bank_type_code"
            label="Bank Type Code"
            variant="outlined"
            density="compact"
            hint="Default: 04 (standard)"
            persistent-hint
          />
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="profileFormDialog = false"
            >Cancel</v-btn
          >
          <v-btn
            color="primary"
            :loading="savingProfile"
            :disabled="
              !profileForm.profile_name ||
              !profileForm.bank_name ||
              !profileForm.user_code
            "
            @click="saveProfile"
          >
            {{ editingProfile ? 'Update' : 'Create' }}
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
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import type { ColDef } from 'ag-grid-community'
import BaseCard from '@/renderer/components/BaseCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
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
  notes: string
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

interface Claim {
  id: number
  claim_number: string
  member_name: string
  member_id_number: string
  scheme_name: string
  benefit_alias: string
  claim_amount: number
  status: string
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

const router = useRouter()
const { hasPermission } = usePermissionCheck()

// ── State ──────────────────────────────────────────────
const loading = ref(false)
const schedules = ref<PaymentSchedule[]>([])

const approvedClaims = ref<Claim[]>([])
const loadingApproved = ref(false)
const claimFilter = ref('')
const benefitFilter = ref('')

const createDialog = ref(false)
const newScheduleDescription = ref('')
const selectedClaimIDs = ref<number[]>([])
const creating = ref(false)

const profilesDialog = ref(false)
const bankProfiles = ref<BankProfile[]>([])
const loadingProfiles = ref(false)
const profileFormDialog = ref(false)
const editingProfile = ref<BankProfile | null>(null)
const savingProfile = ref(false)
const profileForm = ref({
  profile_name: '',
  bank_name: '',
  user_code: '',
  user_branch_code: '',
  user_account_number: '',
  user_account_type: '1',
  bank_type_code: '04',
  service_type: 'two_day'
})

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

// ── Constants ─────────────────────────────────────────
const bankNameOptions = [
  'FNB',
  'Standard Bank',
  'ABSA',
  'Nedbank',
  'Capitec',
  'Investec',
  'African Bank',
  'TymeBank',
  'Discovery Bank',
  'Bank Zero'
]

const universalBranchCodes: Record<string, string> = {
  FNB: '250655',
  'Standard Bank': '051001',
  ABSA: '632005',
  Nedbank: '198765',
  Capitec: '470010',
  Investec: '580105',
  'African Bank': '430000',
  TymeBank: '678910',
  'Discovery Bank': '679000',
  'Bank Zero': '888000'
}

const accountTypeOptions = [
  { title: 'Current/Cheque', value: '1' },
  { title: 'Savings', value: '2' },
  { title: 'Transmission', value: '3' }
]

const serviceTypeOptions = [
  { title: 'Same Day', value: 'same_day' },
  { title: 'One Day', value: 'one_day' },
  { title: 'Two Day', value: 'two_day' }
]

const workflowSteps = [
  { label: 'Draft', sub: 'Created', tone: 'muted' },
  { label: 'Submitted', sub: 'For payment', tone: 'warning' },
  { label: 'ACB Generated', sub: 'Sent to bank', tone: 'accent' },
  { label: 'Reconciled', sub: 'Bank response', tone: 'info' },
  { label: 'Paid', sub: 'Confirmed', tone: 'success' }
]

// ── Computed ────────────────────────────────────────────
const filteredApprovedClaims = computed(() => {
  let result = approvedClaims.value
  if (claimFilter.value) {
    const q = claimFilter.value.toLowerCase()
    result = result.filter((c) => c.claim_number.toLowerCase().includes(q))
  }
  if (benefitFilter.value) {
    const b = benefitFilter.value.toLowerCase()
    result = result.filter((c) => c.benefit_alias?.toLowerCase().includes(b))
  }
  return result
})

const selectedTotal = computed(() => {
  return approvedClaims.value
    .filter((c) => selectedClaimIDs.value.includes(c.id))
    .reduce((sum, c) => sum + c.claim_amount, 0)
})

const paymentStats = computed(() => {
  const now = new Date()
  const currentMonth = now.getMonth()
  const currentYear = now.getFullYear()

  let activeCount = 0
  let pendingValue = 0
  let confirmedThisMonth = 0
  let missingBank = 0

  for (const s of schedules.value) {
    if (s.status !== 'confirmed') {
      activeCount++
      pendingValue += s.total_amount || 0
    } else {
      const d = new Date(s.created_at)
      if (d.getMonth() === currentMonth && d.getFullYear() === currentYear) {
        confirmedThisMonth++
      }
    }
    missingBank += itemsMissingBanking(s).length
  }

  return { activeCount, pendingValue, confirmedThisMonth, missingBank }
})

const claimsColumnDefs: ColDef<Claim>[] = [
  {
    checkboxSelection: true,
    headerCheckboxSelection: true,
    width: 48,
    minWidth: 48,
    maxWidth: 48,
    pinned: 'left',
    suppressMovable: true,
    resizable: false
  },
  {
    headerName: 'Claim #',
    field: 'claim_number',
    sortable: true,
    minWidth: 120
  },
  { headerName: 'Member', field: 'member_name', sortable: true, minWidth: 140 },
  { headerName: 'ID Number', field: 'member_id_number', minWidth: 120 },
  { headerName: 'Scheme', field: 'scheme_name', sortable: true, minWidth: 140 },
  { headerName: 'Benefit', field: 'benefit_alias', minWidth: 130 },
  {
    headerName: 'Amount',
    field: 'claim_amount',
    minWidth: 120,
    type: 'rightAligned',
    valueFormatter: (p) => formatCurrency(p.value)
  },
  {
    headerName: 'Status',
    field: 'status',
    minWidth: 100,
    cellRenderer: (p: any) =>
      `<span class="v-chip v-chip--label v-chip--density-comfortable bg-${statusColor(p.value)}">${p.value}</span>`
  }
]

const scheduleColumnDefs: ColDef<PaymentSchedule>[] = [
  {
    headerName: 'Schedule #',
    field: 'schedule_number',
    sortable: true,
    minWidth: 160
  },
  {
    headerName: 'Status',
    field: 'status',
    sortable: true,
    minWidth: 130,
    cellRenderer: (p: any) => {
      const color = statusColor(p.value)
      const label = statusLabel(p.value)
      return `<span class="v-chip v-chip--label v-chip--size-small bg-${color}" style="font-size:11px;padding:0 8px;height:22px;display:inline-flex;align-items:center">${label}</span>`
    }
  },
  {
    headerName: 'Description',
    field: 'description',
    sortable: true,
    minWidth: 180,
    flex: 1
  },
  {
    headerName: 'Claims',
    field: 'claims_count',
    sortable: true,
    minWidth: 80,
    type: 'rightAligned'
  },
  {
    headerName: 'Total Amount',
    field: 'total_amount',
    sortable: true,
    minWidth: 140,
    type: 'rightAligned',
    valueFormatter: (p) => formatCurrency(p.value)
  },
  {
    headerName: 'Created By',
    field: 'created_by',
    sortable: true,
    minWidth: 120
  },
  {
    headerName: 'Created',
    field: 'created_at',
    sortable: true,
    minWidth: 120,
    valueFormatter: (p) => formatDate(p.value)
  },
  {
    headerName: 'ACB',
    field: 'acb_file_generated',
    minWidth: 70,
    maxWidth: 70,
    cellRenderer: (p: any) =>
      p.value
        ? '<span class="v-chip v-chip--label v-chip--size-x-small bg-teal" style="font-size:10px;padding:0 6px;height:18px;display:inline-flex;align-items:center;color:#fff">ACB</span>'
        : ''
  },
  {
    headerName: 'Actions',
    width: 110,
    minWidth: 110,
    sortable: false,
    filter: false,
    resizable: false,
    cellRenderer: () => {
      const viewBtn = `
        background:#1976D2;color:#fff;border:none;border-radius:4px;
        padding:3px 10px;font-size:11px;font-weight:600;cursor:pointer;
        line-height:1.6;
      `
      return `<button data-action="view" style="${viewBtn}">View</button>`
    },
    onCellClicked: (params: any) => {
      if (params.data) openSchedule(params.data as PaymentSchedule)
    }
  },
  {
    headerName: 'Notes',
    field: 'notes',
    sortable: false,
    filter: true,
    minWidth: 200,
    flex: 1,
    editable: true,
    singleClickEdit: true,
    cellEditor: 'agTextCellEditor',
    cellEditorParams: { maxLength: 30 },
    headerTooltip: 'Double-click to edit · max 30 characters',
    cellStyle: { fontStyle: 'italic', color: '#475569' },
    valueFormatter: (p: any) => p.value || 'Click to add notes…',
    onCellValueChanged: (params: any) => {
      onNotesChanged(params.data?.id, params.newValue || '', params.oldValue || '')
    }
  }
]

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

function itemsMissingBanking(schedule: PaymentSchedule): ScheduleItem[] {
  if (!schedule.items) return []
  return schedule.items.filter(
    (i) => !i.bank_account_number || !i.bank_branch_code
  )
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

// ── Navigation ──────────────────────────────────────────
function goBack() {
  router.push({ name: 'group-pricing-claims-management' })
}

function openSchedule(schedule: PaymentSchedule) {
  router.push({
    name: 'group-pricing-claim-payment-schedule-detail',
    params: { scheduleId: schedule.id }
  })
}

function onScheduleRowDoubleClicked(event: any) {
  const schedule = event?.data as PaymentSchedule | undefined
  if (schedule?.id) openSchedule(schedule)
}

async function onNotesChanged(
  scheduleId: number | undefined,
  newValue: string,
  oldValue: string
) {
  if (!scheduleId || newValue === oldValue) return
  const trimmed = (newValue || '').slice(0, 30)
  try {
    await GroupPricingService.updatePaymentScheduleNotes(scheduleId, trimmed)
    const row = schedules.value.find((s) => s.id === scheduleId)
    if (row) row.notes = trimmed
    notify('Notes saved.')
  } catch (e: any) {
    // revert in-memory value on failure
    const row = schedules.value.find((s) => s.id === scheduleId)
    if (row) row.notes = oldValue
    notify('Failed to save notes', 'error')
  }
}

// ── Data loading ────────────────────────────────────────
async function loadSchedules() {
  loading.value = true
  try {
    const res = await GroupPricingService.getPaymentSchedules()
    schedules.value = unwrap(res) ?? []
  } catch (e: any) {
    notify('Failed to load payment schedules', 'error')
  } finally {
    loading.value = false
  }
}

async function loadApprovedClaims() {
  loadingApproved.value = true
  try {
    const res = await GroupPricingService.getClaims()
    const all: Claim[] = unwrap(res) ?? []
    approvedClaims.value = all.filter((c) => c.status === 'approved')
  } catch (e: any) {
    notify('Failed to load claims', 'error')
  } finally {
    loadingApproved.value = false
  }
}

async function loadBankProfiles() {
  loadingProfiles.value = true
  try {
    const res = await GroupPricingService.getBankProfiles()
    bankProfiles.value = unwrap(res) ?? []
  } catch (e: any) {
    notify('Failed to load bank profiles', 'error')
  } finally {
    loadingProfiles.value = false
  }
}

// ── AG Grid handlers ─────────────────────────────────────
function onClaimSelectionChanged(rows: Claim[]) {
  selectedClaimIDs.value = rows.map((c) => c.id)
}

// ── Schedule creation ───────────────────────────────────
async function openCreateDialog() {
  selectedClaimIDs.value = []
  newScheduleDescription.value = ''
  claimFilter.value = ''
  benefitFilter.value = ''
  await loadApprovedClaims()
  createDialog.value = true
}

async function createSchedule() {
  creating.value = true
  try {
    await GroupPricingService.createPaymentSchedule({
      claim_ids: selectedClaimIDs.value,
      description: newScheduleDescription.value
    })
    createDialog.value = false
    notify(
      'Payment schedule created. Claims moved to "Submitted for Payment".'
    )
    await loadSchedules()
  } catch (e: any) {
    notify(e?.response?.data ?? 'Failed to create schedule', 'error')
  } finally {
    creating.value = false
  }
}

// ── Bank Profile Management ──────────────────────────────
async function openBankProfilesDialog() {
  await loadBankProfiles()
  profilesDialog.value = true
}

function onProfileBankSelected(bankName: string) {
  profileForm.value.user_branch_code = universalBranchCodes[bankName] || ''
}

function openCreateProfileDialog() {
  editingProfile.value = null
  profileForm.value = {
    profile_name: '',
    bank_name: '',
    user_code: '',
    user_branch_code: '',
    user_account_number: '',
    user_account_type: '1',
    bank_type_code: '04',
    service_type: 'two_day'
  }
  profileFormDialog.value = true
}

function editProfile(profile: BankProfile) {
  editingProfile.value = profile
  profileForm.value = {
    profile_name: profile.profile_name,
    bank_name: profile.bank_name,
    user_code: profile.user_code,
    user_branch_code: profile.user_branch_code,
    user_account_number: profile.user_account_number,
    user_account_type: profile.user_account_type || '1',
    bank_type_code: profile.bank_type_code || '04',
    service_type: profile.service_type || 'two_day'
  }
  profileFormDialog.value = true
}

async function saveProfile() {
  savingProfile.value = true
  try {
    if (editingProfile.value) {
      await GroupPricingService.updateBankProfile(
        editingProfile.value.id,
        profileForm.value
      )
      notify('Bank profile updated.')
    } else {
      await GroupPricingService.createBankProfile(profileForm.value)
      notify('Bank profile created.')
    }
    profileFormDialog.value = false
    await loadBankProfiles()
  } catch (e: any) {
    notify(e?.response?.data ?? 'Failed to save bank profile', 'error')
  } finally {
    savingProfile.value = false
  }
}

async function deleteProfile(profile: BankProfile) {
  try {
    await GroupPricingService.deleteBankProfile(profile.id)
    notify('Bank profile deleted.')
    await loadBankProfiles()
  } catch (e: any) {
    notify('Failed to delete bank profile', 'error')
  }
}

onMounted(loadSchedules)
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

/* ── KPI cards ────────────────────────────────────────── */
.kpi-card {
  position: relative;
  display: flex;
  background: rgb(var(--v-theme-surface));
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  border-radius: 10px;
  overflow: hidden;
  min-height: 108px;
  transition:
    transform 0.15s ease,
    box-shadow 0.15s ease,
    border-color 0.15s ease;
}

.kpi-card:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 14px rgba(0, 63, 88, 0.08);
  border-color: rgba(var(--v-theme-on-surface), 0.14);
}

.kpi-card__stripe {
  flex: 0 0 4px;
  background: rgba(var(--v-theme-on-surface), 0.2);
}

.kpi-card__body {
  flex: 1;
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.kpi-card__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.kpi-card__label {
  font-size: 0.72rem;
  font-weight: 600;
  letter-spacing: 0.4px;
  text-transform: uppercase;
  color: rgba(var(--v-theme-on-surface), 0.6);
}

.kpi-card__icon {
  color: rgba(var(--v-theme-on-surface), 0.45);
}

.kpi-card__value {
  font-size: 1.6rem;
  font-weight: 700;
  line-height: 1.15;
  color: rgba(var(--v-theme-on-surface), 0.95);
  word-break: break-word;
}

.kpi-card__hint {
  margin-top: auto;
  padding-top: 4px;
  font-size: 0.72rem;
  color: rgba(var(--v-theme-on-surface), 0.5);
  line-height: 1.3;
}

/* KPI tone variants — accent stripe + icon tint within theme family */
.kpi-card--primary .kpi-card__stripe {
  background: rgb(var(--v-theme-primary));
}
.kpi-card--primary .kpi-card__icon {
  color: rgb(var(--v-theme-primary));
}

.kpi-card--accent .kpi-card__stripe {
  background: rgb(var(--v-theme-accent));
}
.kpi-card--accent .kpi-card__icon {
  color: rgb(var(--v-theme-accent));
}

.kpi-card--success .kpi-card__stripe {
  background: rgb(var(--v-theme-success));
}
.kpi-card--success .kpi-card__icon {
  color: rgb(var(--v-theme-success));
}

.kpi-card--warning .kpi-card__stripe {
  background: rgb(var(--v-theme-warning));
}
.kpi-card--warning .kpi-card__icon {
  color: rgb(var(--v-theme-warning));
}
.kpi-card--warning .kpi-card__value {
  color: rgb(var(--v-theme-warning));
}

.kpi-card--muted .kpi-card__stripe {
  background: rgba(var(--v-theme-on-surface), 0.18);
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
  background: rgba(var(--v-theme-on-surface), 0.35);
}
.workflow-stepper__step--warning .workflow-stepper__num {
  background: rgb(var(--v-theme-warning));
}
.workflow-stepper__step--accent .workflow-stepper__num {
  background: rgb(var(--v-theme-accent));
}
.workflow-stepper__step--info .workflow-stepper__num {
  background: rgb(var(--v-theme-info));
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

.workflow-stepper__sub {
  font-size: 0.7rem;
  color: rgba(var(--v-theme-on-surface), 0.55);
  line-height: 1.2;
}

.workflow-stepper__arrow {
  color: rgba(var(--v-theme-primary), 0.35);
  flex: 0 0 auto;
}

/* ── Schedule card strip ──────────────────────────────── */
.schedule-card-strip {
  overflow-x: auto;
  overflow-y: hidden;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: thin;
  padding-bottom: 4px;
}

.schedule-card-strip__inner {
  display: flex;
  gap: 12px;
}

.schedule-card-strip__item {
  flex: 0 0 320px;
  min-width: 320px;
}

.ps-card {
  position: relative;
  display: flex;
  height: 100%;
  background: rgb(var(--v-theme-surface));
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  border-radius: 10px;
  overflow: hidden;
  cursor: pointer;
  transition:
    transform 0.15s ease,
    box-shadow 0.15s ease,
    border-color 0.15s ease;
}

.ps-card:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 14px rgba(0, 63, 88, 0.1);
  border-color: rgba(var(--v-theme-primary), 0.3);
}

.ps-card__stripe {
  flex: 0 0 4px;
  background: rgba(var(--v-theme-on-surface), 0.2);
}

.ps-card--warning .ps-card__stripe {
  background: rgb(var(--v-theme-warning));
}
.ps-card--success .ps-card__stripe {
  background: rgb(var(--v-theme-success));
}
.ps-card--info .ps-card__stripe {
  background: rgb(var(--v-theme-info));
}
.ps-card--error .ps-card__stripe {
  background: rgb(var(--v-theme-error));
}

.ps-card__body {
  flex: 1;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
}

.ps-card__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 10px;
}

.ps-card__num {
  font-family:
    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.85rem;
  font-weight: 600;
  color: rgb(var(--v-theme-primary));
  letter-spacing: 0.3px;
}

.ps-card__metrics {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.ps-card__metric-label {
  font-size: 0.68rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.4px;
  color: rgba(var(--v-theme-on-surface), 0.55);
  margin-bottom: 2px;
}

.ps-card__metric-value {
  font-size: 0.95rem;
  font-weight: 700;
  color: rgba(var(--v-theme-on-surface), 0.9);
}

.ps-card__footer {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.72rem;
  color: rgba(var(--v-theme-on-surface), 0.55);
}

.ps-card__footer-divider {
  opacity: 0.5;
}

/* ── Grid panel ───────────────────────────────────────── */
.grid-panel {
  background: rgb(var(--v-theme-surface));
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  border-radius: 10px;
  padding: 4px;
}
</style>
