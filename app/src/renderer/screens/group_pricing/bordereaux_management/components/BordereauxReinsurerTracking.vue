<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center justify-space-between">
              <div>
                <span class="headline"
                  >Reinsurer Acceptance &amp; Recovery</span
                >
              </div>
              <v-btn
                size="small"
                rounded
                color="white"
                variant="outlined"
                prepend-icon="mdi-arrow-left"
                @click="$router.push('/group-pricing/bordereaux-management')"
              >
                Back to Dashboard
              </v-btn>
            </div>
          </template>

          <template #default>
            <!-- Filter bar -->
            <v-card variant="outlined" class="mb-4">
              <v-card-text>
                <v-row>
                  <v-col cols="12" sm="6" md="4">
                    <v-text-field
                      v-model="filters.generatedId"
                      label="Bordereaux ID"
                      variant="outlined"
                      density="compact"
                      clearable
                      prepend-inner-icon="mdi-magnify"
                      @update:model-value="loadAll"
                    />
                  </v-col>
                  <v-col cols="12" sm="6" md="3">
                    <v-select
                      v-model="filters.status"
                      :items="currentStatusOptions"
                      label="Status"
                      variant="outlined"
                      density="compact"
                      clearable
                      @update:model-value="loadAll"
                    />
                  </v-col>
                  <v-col v-if="activeTab === 1" cols="12" sm="6" md="3">
                    <v-text-field
                      v-model="filters.claimRef"
                      label="Claim Reference"
                      variant="outlined"
                      density="compact"
                      clearable
                      @update:model-value="loadRecoveries"
                    />
                  </v-col>
                  <v-col cols="auto" class="d-flex align-center gap-2">
                    <v-btn
                      size="small"
                      rounded
                      color="success"
                      prepend-icon="mdi-refresh"
                      :loading="loading"
                      @click="loadAll"
                    >
                      Refresh
                    </v-btn>
                    <v-btn
                      v-if="activeTab === 0"
                      size="small"
                      rounded
                      color="primary"
                      prepend-icon="mdi-plus"
                      @click="openAddAcceptanceDialog"
                    >
                      Add Acceptance
                    </v-btn>
                    <v-btn
                      v-if="activeTab === 1"
                      size="small"
                      rounded
                      color="primary"
                      prepend-icon="mdi-plus"
                      @click="openAddRecoveryDialog"
                    >
                      Record Recovery
                    </v-btn>
                  </v-col>
                </v-row>
              </v-card-text>
            </v-card>

            <!-- Tabs -->
            <v-tabs v-model="activeTab" color="primary" class="mb-4">
              <v-tab :value="0">
                <v-icon class="me-2">mdi-handshake</v-icon>
                Reinsurer Acceptances
                <v-chip
                  v-if="acceptanceSummary.queried > 0"
                  size="x-small"
                  color="warning"
                  class="ms-2"
                >
                  {{ acceptanceSummary.queried }} queried
                </v-chip>
              </v-tab>
              <v-tab :value="1">
                <v-icon class="me-2">mdi-cash-refund</v-icon>
                Recovery Tracking
                <v-chip
                  v-if="recoverySummary.pending > 0"
                  size="x-small"
                  color="warning"
                  class="ms-2"
                >
                  {{ recoverySummary.pending }} pending
                </v-chip>
              </v-tab>
            </v-tabs>

            <!-- ── ACCEPTANCES TAB ── -->
            <v-window v-model="activeTab">
              <v-window-item :value="0">
                <!-- Acceptance KPI chips -->
                <v-row class="mb-3">
                  <v-col cols="auto">
                    <v-chip color="grey" size="small">
                      Total: {{ acceptanceSummary.total }}
                    </v-chip>
                  </v-col>
                  <v-col cols="auto">
                    <v-chip color="blue" size="small">
                      Pending: {{ acceptanceSummary.pending }}
                    </v-chip>
                  </v-col>
                  <v-col cols="auto">
                    <v-chip color="success" size="small">
                      Accepted: {{ acceptanceSummary.accepted }}
                    </v-chip>
                  </v-col>
                  <v-col cols="auto">
                    <v-chip color="warning" size="small">
                      Queried: {{ acceptanceSummary.queried }}
                    </v-chip>
                  </v-col>
                  <v-col cols="auto">
                    <v-chip color="error" size="small">
                      Rejected: {{ acceptanceSummary.rejected }}
                    </v-chip>
                  </v-col>
                </v-row>

                <DataGrid
                  :row-data="acceptances"
                  :column-defs="acceptanceColumnDefs"
                  table-title="Reinsurer Acceptances"
                  :show-export="true"
                  density="compact"
                />
              </v-window-item>

              <!-- ── RECOVERY TAB ── -->
              <v-window-item :value="1">
                <!-- Recovery KPI chips -->
                <v-row class="mb-3">
                  <v-col cols="auto">
                    <v-chip color="grey" size="small">
                      Records: {{ recoverySummary.total }}
                    </v-chip>
                  </v-col>
                  <v-col cols="auto">
                    <v-chip color="primary" size="small">
                      Claimed:
                      {{ fmtCurrency(recoverySummary.totalClaimAmount) }}
                    </v-chip>
                  </v-col>
                  <v-col cols="auto">
                    <v-chip color="success" size="small">
                      Recovered:
                      {{ fmtCurrency(recoverySummary.totalRecovered) }}
                    </v-chip>
                  </v-col>
                  <v-col cols="auto">
                    <v-chip color="error" size="small">
                      Outstanding:
                      {{ fmtCurrency(recoverySummary.totalOutstanding) }}
                    </v-chip>
                  </v-col>
                </v-row>

                <DataGrid
                  :row-data="recoveries"
                  :column-defs="recoveryColumnDefs"
                  table-title="Recovery Tracking"
                  :show-export="true"
                  density="compact"
                />
              </v-window-item>
            </v-window>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- ── Add Acceptance Dialog ── -->
    <v-dialog v-model="showAddAcceptanceDialog" max-width="560">
      <v-card>
        <v-card-title class="text-h6 font-weight-bold bg-primary text-white">
          <v-icon class="me-2">mdi-handshake</v-icon>
          Add Reinsurer Acceptance
        </v-card-title>
        <v-card-text class="pt-6">
          <v-row>
            <v-col cols="12">
              <v-text-field
                v-model="acceptanceForm.generated_bordereaux_id"
                label="Bordereaux ID *"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="acceptanceForm.reinsurer_name"
                label="Reinsurer Name *"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="acceptanceForm.reinsurer_code"
                label="Reinsurer Code"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model.number="acceptanceForm.submitted_amount"
                label="Submitted Amount"
                variant="outlined"
                density="compact"
                type="number"
                prefix="R"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="acceptanceForm.due_date"
                label="Due Date"
                variant="outlined"
                density="compact"
                placeholder="YYYY-MM-DD"
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="acceptanceForm.notes"
                label="Notes"
                variant="outlined"
                density="compact"
                rows="2"
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            size="small"
            rounded
            color="grey"
            variant="outlined"
            @click="showAddAcceptanceDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn
            size="small"
            rounded
            color="primary"
            variant="flat"
            :loading="saving"
            :disabled="
              !acceptanceForm.generated_bordereaux_id ||
              !acceptanceForm.reinsurer_name
            "
            @click="saveAcceptance"
          >
            Save
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- ── Update Acceptance Dialog ── -->
    <v-dialog v-model="showUpdateAcceptanceDialog" max-width="560">
      <v-card v-if="selectedAcceptance">
        <v-card-title class="text-h6 font-weight-bold bg-teal text-white">
          <v-icon class="me-2">mdi-pencil</v-icon>
          Update Acceptance — {{ selectedAcceptance.reinsurer_name }}
        </v-card-title>
        <v-card-text class="pt-6">
          <v-row>
            <v-col cols="12" md="6">
              <v-select
                v-model="updateAcceptanceForm.status"
                :items="acceptanceStatusOptions"
                label="Status"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model.number="updateAcceptanceForm.accepted_amount"
                label="Accepted Amount"
                variant="outlined"
                density="compact"
                type="number"
                prefix="R"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="updateAcceptanceForm.received_date"
                label="Received Date"
                variant="outlined"
                density="compact"
                placeholder="YYYY-MM-DD"
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="updateAcceptanceForm.query_details"
                label="Query Details"
                variant="outlined"
                density="compact"
                rows="2"
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="updateAcceptanceForm.notes"
                label="Notes"
                variant="outlined"
                density="compact"
                rows="2"
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            size="small"
            rounded
            color="grey"
            variant="outlined"
            @click="showUpdateAcceptanceDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn
            size="small"
            rounded
            color="teal"
            variant="flat"
            :loading="saving"
            @click="performUpdateAcceptance"
          >
            Update
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- ── Add Recovery Dialog ── -->
    <v-dialog v-model="showAddRecoveryDialog" max-width="560">
      <v-card>
        <v-card-title class="text-h6 font-weight-bold bg-primary text-white">
          <v-icon class="me-2">mdi-cash-refund</v-icon>
          Record Claim Recovery
        </v-card-title>
        <v-card-text class="pt-6">
          <v-row>
            <v-col cols="12">
              <v-text-field
                v-model="recoveryForm.generated_bordereaux_id"
                label="Bordereaux ID *"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="recoveryForm.claim_reference"
                label="Claim Reference *"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="recoveryForm.reinsurer_name"
                label="Reinsurer Name *"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="recoveryForm.reinsurer_code"
                label="Reinsurer Code"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model.number="recoveryForm.claim_amount"
                label="Claim Amount"
                variant="outlined"
                density="compact"
                type="number"
                prefix="R"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model.number="recoveryForm.recovered_amount"
                label="Recovered Amount"
                variant="outlined"
                density="compact"
                type="number"
                prefix="R"
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="recoveryForm.notes"
                label="Notes"
                variant="outlined"
                density="compact"
                rows="2"
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            size="small"
            rounded
            color="grey"
            variant="outlined"
            @click="showAddRecoveryDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn
            size="small"
            rounded
            color="primary"
            variant="flat"
            :loading="saving"
            :disabled="
              !recoveryForm.generated_bordereaux_id ||
              !recoveryForm.claim_reference ||
              !recoveryForm.reinsurer_name
            "
            @click="saveRecovery"
          >
            Save
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- ── Update Recovery Dialog ── -->
    <v-dialog v-model="showUpdateRecoveryDialog" max-width="560">
      <v-card v-if="selectedRecovery">
        <v-card-title class="text-h6 font-weight-bold bg-teal text-white">
          <v-icon class="me-2">mdi-pencil</v-icon>
          Update Recovery — {{ selectedRecovery.claim_reference }}
        </v-card-title>
        <v-card-text class="pt-6">
          <v-row>
            <v-col cols="12" md="6">
              <v-text-field
                v-model.number="updateRecoveryForm.recovered_amount"
                label="Recovered Amount"
                variant="outlined"
                density="compact"
                type="number"
                prefix="R"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-select
                v-model="updateRecoveryForm.status"
                :items="recoveryStatusOptions"
                label="Status"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="updateRecoveryForm.received_date"
                label="Received Date"
                variant="outlined"
                density="compact"
                placeholder="YYYY-MM-DD"
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="updateRecoveryForm.notes"
                label="Notes"
                variant="outlined"
                density="compact"
                rows="2"
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            size="small"
            rounded
            color="grey"
            variant="outlined"
            @click="showUpdateRecoveryDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn
            size="small"
            rounded
            color="teal"
            variant="flat"
            :loading="saving"
            @click="performUpdateRecovery"
          >
            Update
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Snackbar -->
    <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
      {{ snackbar.message }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

// ── State ──────────────────────────────────────────────────────────────────
const loading = ref(false)
const saving = ref(false)
const activeTab = ref(0)

const acceptances: any = ref([])
const recoveries: any = ref([])

const acceptanceSummary = ref({
  total: 0,
  pending: 0,
  accepted: 0,
  queried: 0,
  rejected: 0,
  totalSubmitted: 0,
  totalAccepted: 0,
  totalVariance: 0
})

const recoverySummary = ref({
  total: 0,
  pending: 0,
  partial: 0,
  full: 0,
  disputed: 0,
  totalClaimAmount: 0,
  totalRecovered: 0,
  totalOutstanding: 0
})

const filters = ref({ generatedId: '', status: '', claimRef: '' })

// Dialogs
const showAddAcceptanceDialog = ref(false)
const showUpdateAcceptanceDialog = ref(false)
const selectedAcceptance: any = ref(null)
const showAddRecoveryDialog = ref(false)
const showUpdateRecoveryDialog = ref(false)
const selectedRecovery: any = ref(null)

const acceptanceForm = ref({
  generated_bordereaux_id: '',
  reinsurer_name: '',
  reinsurer_code: '',
  submitted_amount: 0,
  due_date: '',
  notes: ''
})

const updateAcceptanceForm: any = ref({
  status: '',
  accepted_amount: 0,
  query_details: '',
  notes: '',
  received_date: ''
})

const recoveryForm = ref({
  generated_bordereaux_id: '',
  claim_reference: '',
  reinsurer_name: '',
  reinsurer_code: '',
  claim_amount: 0,
  recovered_amount: 0,
  notes: ''
})

const updateRecoveryForm: any = ref({
  recovered_amount: 0,
  status: '',
  notes: '',
  received_date: ''
})

const snackbar = ref({ show: false, message: '', color: 'success' })

// ── Static data ────────────────────────────────────────────────────────────
const acceptanceStatusOptions = [
  { title: 'Pending', value: 'pending' },
  { title: 'Accepted', value: 'accepted' },
  { title: 'Queried', value: 'queried' },
  { title: 'Rejected', value: 'rejected' }
]

const recoveryStatusOptions = [
  { title: 'Pending', value: 'pending' },
  { title: 'Partial', value: 'partial' },
  { title: 'Full', value: 'full' },
  { title: 'Disputed', value: 'disputed' }
]

const currentStatusOptions = computed(() =>
  activeTab.value === 0 ? acceptanceStatusOptions : recoveryStatusOptions
)

// ── Column definitions ─────────────────────────────────────────────────────
const acceptanceColumnDefs = ref([
  {
    headerName: 'Bordereaux ID',
    field: 'generated_bordereaux_id',
    sortable: true,
    filter: true,
    width: 180
  },
  {
    headerName: 'Reinsurer',
    field: 'reinsurer_name',
    sortable: true,
    filter: true
  },
  { headerName: 'Code', field: 'reinsurer_code', width: 90 },
  {
    headerName: 'Status',
    field: 'status',
    sortable: true,
    filter: true,
    width: 120,
    cellRenderer: (p: any) => {
      const c = acceptanceStatusColor(p.value)
      return `<span style="background:${c}22;color:${c};padding:2px 8px;border-radius:4px;font-size:12px">${p.value}</span>`
    }
  },
  {
    headerName: 'Submitted (R)',
    field: 'submitted_amount',
    sortable: true,
    width: 130,
    cellRenderer: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: 'Accepted (R)',
    field: 'accepted_amount',
    sortable: true,
    width: 130,
    cellRenderer: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: 'Variance (R)',
    field: 'variance',
    sortable: true,
    width: 130,
    cellRenderer: (p: any) => {
      const v = p.value || 0
      const colour = v > 0 ? '#F44336' : v < 0 ? '#4CAF50' : '#9E9E9E'
      return `<span style="color:${colour};font-weight:500">${fmtCurrency(Math.abs(v))}</span>`
    }
  },
  { headerName: 'Due Date', field: 'due_date', width: 110 },
  {
    headerName: 'Received',
    field: 'received_date',
    width: 120,
    cellRenderer: (p: any) => (p.value ? fmtDate(p.value) : '—')
  },
  {
    headerName: 'Actions',
    field: 'actions',
    sortable: false,
    filter: false,
    width: 90,
    cellRenderer: (p: any) => {
      const item = p.data
      return `<div class="ag-actions-cell">
        <button onclick="window.editAcceptance_${item.id}()" class="ag-btn ag-btn-primary" title="Edit">✏️</button>
      </div>`
    }
  }
])

const recoveryColumnDefs = ref([
  {
    headerName: 'Bordereaux ID',
    field: 'generated_bordereaux_id',
    sortable: true,
    filter: true,
    width: 180
  },
  {
    headerName: 'Claim Ref',
    field: 'claim_reference',
    sortable: true,
    filter: true
  },
  {
    headerName: 'Reinsurer',
    field: 'reinsurer_name',
    sortable: true,
    filter: true
  },
  {
    headerName: 'Status',
    field: 'status',
    sortable: true,
    filter: true,
    width: 110,
    cellRenderer: (p: any) => {
      const c = recoveryStatusColor(p.value)
      return `<span style="background:${c}22;color:${c};padding:2px 8px;border-radius:4px;font-size:12px">${p.value}</span>`
    }
  },
  {
    headerName: 'Claim (R)',
    field: 'claim_amount',
    sortable: true,
    width: 120,
    cellRenderer: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: 'Recovered (R)',
    field: 'recovered_amount',
    sortable: true,
    width: 130,
    cellRenderer: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: 'Recovery %',
    field: 'recovery_percentage',
    sortable: true,
    width: 120,
    cellRenderer: (p: any) => {
      const pct = Math.min(p.value || 0, 100)
      const colour = pct >= 100 ? '#4CAF50' : pct >= 50 ? '#FF9800' : '#F44336'
      return `<div style="display:flex;align-items:center;gap:6px">
        <div style="width:60px;height:6px;background:#e0e0e0;border-radius:3px;overflow:hidden">
          <div style="width:${pct}%;height:100%;background:${colour}"></div>
        </div>
        <span style="font-size:12px">${pct.toFixed(1)}%</span>
      </div>`
    }
  },
  {
    headerName: 'Received',
    field: 'received_date',
    width: 120,
    cellRenderer: (p: any) => (p.value ? fmtDate(p.value) : '—')
  },
  {
    headerName: 'Actions',
    field: 'actions',
    sortable: false,
    filter: false,
    width: 90,
    cellRenderer: (p: any) => {
      const item = p.data
      return `<div class="ag-actions-cell">
        <button onclick="window.editRecovery_${item.id}()" class="ag-btn ag-btn-primary" title="Edit">✏️</button>
      </div>`
    }
  }
])

// ── Helpers ────────────────────────────────────────────────────────────────
const fmtCurrency = (v: number) =>
  (v || 0).toLocaleString('en-ZA', {
    style: 'currency',
    currency: 'ZAR',
    maximumFractionDigits: 2
  })

const fmtDate = (d: string) =>
  new Date(d).toLocaleDateString('en-ZA', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })

const acceptanceStatusColor = (s: string) =>
  ({
    pending: '#9E9E9E',
    accepted: '#4CAF50',
    queried: '#FF9800',
    rejected: '#F44336'
  })[s] ?? '#9E9E9E'

const recoveryStatusColor = (s: string) =>
  ({
    pending: '#9E9E9E',
    partial: '#FF9800',
    full: '#4CAF50',
    disputed: '#F44336'
  })[s] ?? '#9E9E9E'

const showSnack = (message: string, color = 'success') => {
  snackbar.value = { show: true, message, color }
}

// ── Data loaders ───────────────────────────────────────────────────────────
const loadAcceptances = async () => {
  try {
    const res = await GroupPricingService.getReinsurerAcceptances({
      generated_id: filters.value.generatedId || undefined,
      status: filters.value.status || undefined
    })
    acceptances.value = res.data?.data || []
    computeAcceptanceSummary()
    setupAcceptanceActions()
  } catch {
    // silently degrade
  }
}

const loadRecoveries = async () => {
  try {
    const res = await GroupPricingService.getReinsurerRecoveries({
      generated_id: filters.value.generatedId || undefined,
      claim_ref: filters.value.claimRef || undefined,
      status: filters.value.status || undefined
    })
    recoveries.value = res.data?.data || []
    computeRecoverySummary()
    setupRecoveryActions()
  } catch {
    // silently degrade
  }
}

const loadAll = async () => {
  loading.value = true
  await Promise.all([loadAcceptances(), loadRecoveries()])
  loading.value = false
}

const computeAcceptanceSummary = () => {
  const rows: any[] = acceptances.value
  acceptanceSummary.value = {
    total: rows.length,
    pending: rows.filter((r) => r.status === 'pending').length,
    accepted: rows.filter((r) => r.status === 'accepted').length,
    queried: rows.filter((r) => r.status === 'queried').length,
    rejected: rows.filter((r) => r.status === 'rejected').length,
    totalSubmitted: rows.reduce((s, r) => s + (r.submitted_amount || 0), 0),
    totalAccepted: rows.reduce((s, r) => s + (r.accepted_amount || 0), 0),
    totalVariance: rows.reduce((s, r) => s + (r.variance || 0), 0)
  }
}

const computeRecoverySummary = () => {
  const rows: any[] = recoveries.value
  recoverySummary.value = {
    total: rows.length,
    pending: rows.filter((r) => r.status === 'pending').length,
    partial: rows.filter((r) => r.status === 'partial').length,
    full: rows.filter((r) => r.status === 'full').length,
    disputed: rows.filter((r) => r.status === 'disputed').length,
    totalClaimAmount: rows.reduce((s, r) => s + (r.claim_amount || 0), 0),
    totalRecovered: rows.reduce((s, r) => s + (r.recovered_amount || 0), 0),
    totalOutstanding: rows.reduce(
      (s, r) => s + (r.claim_amount - r.recovered_amount || 0),
      0
    )
  }
}

// ── Global action handlers ─────────────────────────────────────────────────
const setupAcceptanceActions = () => {
  acceptances.value.forEach((item: any) => {
    ;(window as any)[`editAcceptance_${item.id}`] = () =>
      openUpdateAcceptanceDialog(item)
  })
}

const setupRecoveryActions = () => {
  recoveries.value.forEach((item: any) => {
    ;(window as any)[`editRecovery_${item.id}`] = () =>
      openUpdateRecoveryDialog(item)
  })
}

watch(() => acceptances.value, setupAcceptanceActions, { deep: true })
watch(() => recoveries.value, setupRecoveryActions, { deep: true })

// ── Dialog handlers ────────────────────────────────────────────────────────
const openAddAcceptanceDialog = () => {
  acceptanceForm.value = {
    generated_bordereaux_id: filters.value.generatedId || '',
    reinsurer_name: '',
    reinsurer_code: '',
    submitted_amount: 0,
    due_date: '',
    notes: ''
  }
  showAddAcceptanceDialog.value = true
}

const saveAcceptance = async () => {
  saving.value = true
  try {
    await GroupPricingService.createReinsurerAcceptance(acceptanceForm.value)
    showSnack('Acceptance record created')
    showAddAcceptanceDialog.value = false
    await loadAcceptances()
  } catch (e: any) {
    showSnack(e.response?.data?.error || 'Failed to create acceptance', 'error')
  } finally {
    saving.value = false
  }
}

const openUpdateAcceptanceDialog = (item: any) => {
  selectedAcceptance.value = item
  updateAcceptanceForm.value = {
    status: item.status,
    accepted_amount: item.accepted_amount,
    query_details: item.query_details,
    notes: item.notes,
    received_date: item.received_date ? item.received_date.substring(0, 10) : ''
  }
  showUpdateAcceptanceDialog.value = true
}

const performUpdateAcceptance = async () => {
  saving.value = true
  try {
    const payload: any = { ...updateAcceptanceForm.value }
    if (payload.received_date) {
      payload.received_date = new Date(payload.received_date).toISOString()
    }
    await GroupPricingService.updateReinsurerAcceptance(
      selectedAcceptance.value.id,
      payload
    )
    showSnack('Acceptance updated')
    showUpdateAcceptanceDialog.value = false
    await loadAcceptances()
  } catch (e: any) {
    showSnack(e.response?.data?.error || 'Failed to update acceptance', 'error')
  } finally {
    saving.value = false
  }
}

const openAddRecoveryDialog = () => {
  recoveryForm.value = {
    generated_bordereaux_id: filters.value.generatedId || '',
    claim_reference: '',
    reinsurer_name: '',
    reinsurer_code: '',
    claim_amount: 0,
    recovered_amount: 0,
    notes: ''
  }
  showAddRecoveryDialog.value = true
}

const saveRecovery = async () => {
  saving.value = true
  try {
    await GroupPricingService.createReinsurerRecovery(recoveryForm.value)
    showSnack('Recovery record created')
    showAddRecoveryDialog.value = false
    await loadRecoveries()
  } catch (e: any) {
    showSnack(e.response?.data?.error || 'Failed to create recovery', 'error')
  } finally {
    saving.value = false
  }
}

const openUpdateRecoveryDialog = (item: any) => {
  selectedRecovery.value = item
  updateRecoveryForm.value = {
    recovered_amount: item.recovered_amount,
    status: item.status,
    notes: item.notes,
    received_date: item.received_date ? item.received_date.substring(0, 10) : ''
  }
  showUpdateRecoveryDialog.value = true
}

const performUpdateRecovery = async () => {
  saving.value = true
  try {
    const payload: any = { ...updateRecoveryForm.value }
    if (payload.received_date) {
      payload.received_date = new Date(payload.received_date).toISOString()
    }
    await GroupPricingService.updateReinsurerRecovery(
      selectedRecovery.value.id,
      payload
    )
    showSnack('Recovery updated')
    showUpdateRecoveryDialog.value = false
    await loadRecoveries()
  } catch (e: any) {
    showSnack(e.response?.data?.error || 'Failed to update recovery', 'error')
  } finally {
    saving.value = false
  }
}

onMounted(loadAll)
</script>

<style scoped>
:deep(.ag-actions-cell) {
  display: flex;
  align-items: center;
  gap: 4px;
  height: 100%;
}

:deep(.ag-btn) {
  border: none;
  background: transparent;
  cursor: pointer;
  padding: 4px 6px;
  border-radius: 4px;
  font-size: 14px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

:deep(.ag-btn-primary) {
  color: rgb(var(--v-theme-primary));
}

:deep(.ag-btn-primary:hover) {
  background-color: rgba(var(--v-theme-primary), 0.1);
}
</style>
