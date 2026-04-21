<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center">
              <v-btn
                class="mr-3"
                size="small"
                variant="text"
                prepend-icon="mdi-arrow-left"
                @click="$router.back()"
              >
                Back
              </v-btn>
              <span class="headline">RI Submission Register</span>
            </div>
          </template>

          <template #default>
            <!-- KPI chips -->
            <v-row class="mb-4">
              <v-col cols="auto">
                <v-chip
                  color="primary"
                  variant="tonal"
                  prepend-icon="mdi-file-document-multiple"
                >
                  Total Runs: {{ runs.length }}
                </v-chip>
              </v-col>
              <v-col cols="auto">
                <v-chip
                  color="success"
                  variant="tonal"
                  prepend-icon="mdi-check-circle"
                >
                  Acknowledged: {{ acknowledgedCount }}
                </v-chip>
              </v-col>
              <v-col cols="auto">
                <v-chip
                  color="warning"
                  variant="tonal"
                  prepend-icon="mdi-clock-outline"
                >
                  Pending: {{ pendingCount }}
                </v-chip>
              </v-col>
              <v-col cols="auto">
                <v-chip
                  color="deep-purple"
                  variant="tonal"
                  prepend-icon="mdi-pencil-plus"
                >
                  Amendments: {{ amendmentCount }}
                </v-chip>
              </v-col>
            </v-row>

            <!-- Filter bar -->
            <v-card variant="outlined" class="mb-4">
              <v-card-text>
                <v-row>
                  <v-col cols="12" sm="3">
                    <v-text-field
                      v-model="filters.bpr"
                      label="BPR"
                      variant="outlined"
                      density="compact"
                      clearable
                      prepend-inner-icon="mdi-magnify"
                      @update:model-value="applyFilters"
                    />
                  </v-col>
                  <v-col cols="12" sm="3">
                    <v-select
                      v-model="filters.status"
                      :items="statusOptions"
                      label="Status"
                      variant="outlined"
                      density="compact"
                      clearable
                      @update:model-value="applyFilters"
                    />
                  </v-col>
                  <v-col cols="12" sm="3">
                    <v-select
                      v-model="filters.type"
                      :items="typeOptions"
                      label="Run Type"
                      variant="outlined"
                      density="compact"
                      clearable
                      @update:model-value="applyFilters"
                    />
                  </v-col>
                  <v-col cols="auto" class="d-flex align-center">
                    <v-btn
                      size="small"
                      rounded
                      color="success"
                      prepend-icon="mdi-refresh"
                      :loading="loading"
                      @click="loadRuns"
                    >
                      Refresh
                    </v-btn>
                  </v-col>
                </v-row>
              </v-card-text>
            </v-card>

            <!-- Grid -->
            <div style="height: 520px">
              <ag-grid-vue
                class="ag-theme-material"
                style="height: 100%; width: 100%"
                :row-data="filteredRuns"
                :column-defs="columnDefs"
                :default-col-def="{
                  resizable: true,
                  sortable: true,
                  filter: true
                }"
                :pagination="true"
                :pagination-page-size="20"
              />
            </div>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Acknowledge Receipt Dialog -->
    <v-dialog v-model="showAckDialog" max-width="460">
      <v-card>
        <v-card-title>Log Reinsurer Receipt</v-card-title>
        <v-card-text>
          <p class="mb-3 text-body-2">
            Confirm that run <strong>{{ ackingRun?.run_id }}</strong> (BPR:
            <strong>{{ ackingRun?.bpr }}</strong
            >) has been received and acknowledged by
            {{ ackingRun?.reinsurer_name }}.
          </p>
          <v-text-field
            v-model="ackForm.received_date"
            label="Date Received by Reinsurer"
            type="date"
            variant="outlined"
            density="compact"
          />
          <v-textarea
            v-model="ackForm.notes"
            label="Notes (optional)"
            variant="outlined"
            density="compact"
            rows="2"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showAckDialog = false">Cancel</v-btn>
          <v-btn color="success" :loading="acking" @click="performAck"
            >Confirm Receipt</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Amend Dialog -->
    <v-dialog v-model="showAmendDialog" max-width="500">
      <v-card>
        <v-card-title>Create Amendment</v-card-title>
        <v-card-text>
          <p class="mb-3 text-body-2">
            This will create a new
            <strong>Version {{ (amendingRun?.run_version ?? 0) + 1 }}</strong>
            amendment run for <strong>{{ amendingRun?.run_id }}</strong> ({{
              amendingRun?.period_label
            }}), copying all
            {{ amendingRun?.type === 'member_census' ? 'member' : 'claims' }}
            rows from the parent for editing.
          </p>
          <v-textarea
            v-model="amendForm.amendment_notes"
            label="Amendment Reason *"
            variant="outlined"
            density="compact"
            rows="3"
            hint="Describe what changed in this amendment"
            persistent-hint
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showAmendDialog = false">Cancel</v-btn>
          <v-btn
            color="deep-purple"
            :loading="amending"
            :disabled="!amendForm.amendment_notes"
            @click="performAmend"
          >
            Create Amendment
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3500">
      {{ snackbar.message }}
    </v-snackbar>
  </v-container>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { AgGridVue } from 'ag-grid-vue3'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import BaseCard from '@/renderer/components/BaseCard.vue'

const runs = ref([])
const loading = ref(false)
const acking = ref(false)
const amending = ref(false)
const filters = ref({ bpr: '', status: '', type: '' })
const showAckDialog = ref(false)
const showAmendDialog = ref(false)
const ackingRun = ref(null)
const amendingRun = ref(null)
const ackForm = ref({ received_date: '', notes: '' })
const amendForm = ref({ amendment_notes: '' })
const snackbar = ref({ show: false, message: '', color: 'success' })

const statusOptions = [
  { title: 'Generated', value: 'generated' },
  { title: 'Validating', value: 'validating' },
  { title: 'Validated', value: 'validated' },
  { title: 'Validation Failed', value: 'validation_failed' },
  { title: 'Submitted', value: 'submitted' },
  { title: 'Acknowledged', value: 'acknowledged' },
  { title: 'Queried', value: 'queried' },
  { title: 'Settled', value: 'settled' }
]

const typeOptions = [
  { title: 'Member Census', value: 'member_census' },
  { title: 'Claims Run', value: 'claims_run' }
]

const statusColor = {
  draft: '#9e9e9e',
  generated: '#1976d2',
  validating: '#7b1fa2',
  validated: '#2e7d32',
  validation_failed: '#c62828',
  submitted: '#f57c00',
  acknowledged: '#388e3c',
  queried: '#ef6c00',
  settled: '#00695c'
}

const formatCurrency = (v) =>
  v != null
    ? 'R ' +
      Number(v).toLocaleString('en-ZA', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      })
    : '—'

const acknowledgedCount = computed(
  () =>
    runs.value.filter(
      (r) => r.status === 'acknowledged' || r.status === 'settled'
    ).length
)
const pendingCount = computed(
  () =>
    runs.value.filter((r) =>
      ['generated', 'validating', 'validated', 'submitted'].includes(r.status)
    ).length
)
const amendmentCount = computed(
  () => runs.value.filter((r) => r.parent_run_id).length
)

const filteredRuns = computed(() => {
  return runs.value.filter((r) => {
    if (
      filters.value.bpr &&
      !r.bpr?.toLowerCase().includes(filters.value.bpr.toLowerCase())
    )
      return false
    if (filters.value.status && r.status !== filters.value.status) return false
    if (filters.value.type && r.type !== filters.value.type) return false
    return true
  })
})

const columnDefs = [
  {
    field: 'bpr',
    headerName: 'BPR',
    width: 160,
    cellRenderer: (p) =>
      p.value
        ? `<span style="font-family:monospace;font-weight:600;color:#1565c0">${p.value}</span>`
        : '<span style="color:#9e9e9e">—</span>'
  },
  { field: 'run_id', headerName: 'Run ID', width: 170 },
  { field: 'treaty_number', headerName: 'Treaty', width: 120 },
  { field: 'reinsurer_name', headerName: 'Reinsurer', width: 150 },
  { field: 'period_label', headerName: 'Period', width: 120 },
  {
    field: 'type',
    headerName: 'Type',
    width: 120,
    cellRenderer: (p) => {
      const label = p.value === 'member_census' ? 'Member' : 'Claims'
      const color = p.value === 'member_census' ? '#1565c0' : '#b71c1c'
      return `<span style="padding:2px 8px;border-radius:12px;font-size:11px;background:${color}22;color:${color};font-weight:500">${label}</span>`
    }
  },
  {
    field: 'run_version',
    headerName: 'Ver.',
    width: 65,
    cellRenderer: (p) => {
      const v = p.value || 1
      const color = v > 1 ? '#7b1fa2' : '#9e9e9e'
      return `<span style="font-weight:600;color:${color}">v${v}</span>`
    }
  },
  {
    field: 'status',
    headerName: 'Status',
    width: 140,
    cellRenderer: (p) => {
      const c = statusColor[p.value] || '#9e9e9e'
      return `<span style="padding:2px 8px;border-radius:12px;font-size:11px;background:${c}22;color:${c};font-weight:500">${p.value?.replace(/_/g, ' ')}</span>`
    }
  },
  {
    field: 'received_date',
    headerName: 'Received',
    width: 110,
    cellRenderer: (p) => p.value || '<span style="color:#9e9e9e">—</span>'
  },
  {
    field: 'acknowledged_by',
    headerName: "Ack'd By",
    width: 120,
    cellRenderer: (p) => p.value || '<span style="color:#9e9e9e">—</span>'
  },
  {
    field: 'ceded_premium',
    headerName: 'Ceded Premium',
    width: 140,
    valueFormatter: (p) => formatCurrency(p.value)
  },
  { field: 'generated_by', headerName: 'Generated By', width: 130 },
  {
    headerName: 'Actions',
    width: 280,
    sortable: false,
    filter: false,
    cellRenderer: (p) => {
      const id = p.data.run_id
      const key = id.replace(/-/g, '_')
      window[`regAck_${key}`] = () => openAckDialog(p.data)
      window[`regAmend_${key}`] = () => openAmendDialog(p.data)
      window[`regView_${key}`] = () => goToGeneration(p.data)
      const canAck = p.data.status === 'submitted'
      const canAmend =
        p.data.status === 'submitted' || p.data.status === 'acknowledged'
      return `<div style="display:flex;gap:4px;align-items:center;height:100%">
        <button onclick="regView_${key}()" style="padding:2px 8px;border-radius:4px;border:1px solid #1976d2;color:#1976d2;background:none;cursor:pointer;font-size:12px">View</button>
        ${canAck ? `<button onclick="regAck_${key}()" style="padding:2px 8px;border-radius:4px;border:1px solid #388e3c;color:#388e3c;background:none;cursor:pointer;font-size:12px">Log Receipt</button>` : ''}
        ${canAmend ? `<button onclick="regAmend_${key}()" style="padding:2px 8px;border-radius:4px;border:1px solid #7b1fa2;color:#7b1fa2;background:none;cursor:pointer;font-size:12px">Amend</button>` : ''}
      </div>`
    }
  }
]

const notify = (message, color = 'success') => {
  snackbar.value = { show: true, message, color }
}

async function loadRuns() {
  loading.value = true
  try {
    const res = await GroupPricingService.getRIBordereauxRuns({})
    runs.value = res.data?.data || []
  } catch {
    notify('Failed to load runs', 'error')
  } finally {
    loading.value = false
  }
}

function applyFilters() {
  // filteredRuns computed handles client-side filtering
}

function openAckDialog(run) {
  ackingRun.value = run
  ackForm.value = {
    received_date: new Date().toISOString().split('T')[0],
    notes: ''
  }
  showAckDialog.value = true
}

async function performAck() {
  acking.value = true
  try {
    await GroupPricingService.acknowledgeRIBordereauxReceipt(
      ackingRun.value.run_id,
      ackForm.value
    )
    notify('Receipt logged successfully')
    showAckDialog.value = false
    loadRuns()
  } catch (e) {
    notify(e.response?.data?.error || 'Failed to log receipt', 'error')
  } finally {
    acking.value = false
  }
}

function openAmendDialog(run) {
  amendingRun.value = run
  amendForm.value = { amendment_notes: '' }
  showAmendDialog.value = true
}

async function performAmend() {
  amending.value = true
  try {
    const res = await GroupPricingService.amendRIBordereaux(
      amendingRun.value.run_id,
      amendForm.value
    )
    const newRun = res.data?.data
    notify(`Amendment created: ${newRun?.run_id} (v${newRun?.run_version})`)
    showAmendDialog.value = false
    loadRuns()
  } catch (e) {
    notify(e.response?.data?.error || 'Failed to create amendment', 'error')
  } finally {
    amending.value = false
  }
}

function goToGeneration() {
  // Navigate to RI bordereaux generation screen where runs can be edited/validated
  window.__vue_router__?.push(
    '/group-pricing/bordereaux-management/ri-bordereaux'
  )
}

onMounted(() => {
  loadRuns()
})
</script>
