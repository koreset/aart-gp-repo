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
              <span class="headline">RI Member &amp; Claims Census</span>
            </div>
          </template>

          <template #default>
            <!-- KPI Cards -->
            <v-row class="mb-4">
              <v-col cols="6" sm="3">
                <v-card variant="outlined">
                  <v-card-text class="text-center">
                    <div class="text-h5 font-weight-bold text-primary">{{
                      stats.total_runs
                    }}</div>
                    <div class="text-caption text-medium-emphasis"
                      >Total Runs</div
                    >
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="6" sm="3">
                <v-card variant="outlined">
                  <v-card-text class="text-center">
                    <div class="text-h5 font-weight-bold text-indigo">{{
                      stats.member_census
                    }}</div>
                    <div class="text-caption text-medium-emphasis"
                      >Member Census</div
                    >
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="6" sm="3">
                <v-card variant="outlined">
                  <v-card-text class="text-center">
                    <div class="text-h5 font-weight-bold text-deep-orange">{{
                      stats.claims_runs
                    }}</div>
                    <div class="text-caption text-medium-emphasis"
                      >Claims Runs</div
                    >
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="6" sm="3">
                <v-card variant="outlined">
                  <v-card-text class="text-center">
                    <div class="text-h5 font-weight-bold text-success">{{
                      stats.submitted
                    }}</div>
                    <div class="text-caption text-medium-emphasis"
                      >Submitted</div
                    >
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Tabs -->
            <v-tabs v-model="activeTab" class="mb-4">
              <v-tab value="member">Member Census</v-tab>
              <v-tab value="claims">Claims Runs</v-tab>
            </v-tabs>

            <!-- Filter bar -->
            <v-card variant="outlined" class="mb-4">
              <v-card-text>
                <v-row>
                  <v-col cols="12" sm="4">
                    <v-select
                      v-model="filters.treaty_id"
                      :items="activeTreaties"
                      item-title="treaty_name"
                      item-value="id"
                      label="Treaty"
                      variant="outlined"
                      density="compact"
                      clearable
                      @update:model-value="loadRuns"
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
                      @update:model-value="loadRuns"
                    />
                  </v-col>
                  <v-col cols="auto" class="d-flex align-center gap-2">
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
                    <v-btn
                      size="small"
                      rounded
                      color="primary"
                      prepend-icon="mdi-plus"
                      @click="openGenerateDialog"
                    >
                      Generate
                    </v-btn>
                  </v-col>
                </v-row>
              </v-card-text>
            </v-card>

            <!-- Grid -->
            <div style="height: 460px">
              <data-grid
                :row-data="filteredRuns"
                :column-defs="columnDefs"
                :default-col-def="{
                  resizable: true,
                  sortable: true,
                  filter: true
                }"
                :pagination="true"
                :pagination-page-size="20"
                :loading="loading"
                no-rows-message="No RI bordereaux runs yet for the selected filters."
              />
            </div>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Generate Dialog -->
    <v-dialog v-model="showGenerateDialog" max-width="560">
      <v-card>
        <v-card-title>Generate RI Bordereaux</v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12">
              <v-select
                v-model="genForm.treaty_id"
                :items="activeTreaties"
                item-title="treaty_name"
                item-value="id"
                label="Treaty *"
                variant="outlined"
                density="compact"
                @update:model-value="onGenTreatySelected"
              />
            </v-col>
            <v-col cols="6">
              <v-select
                v-model="genForm.month"
                :items="months"
                label="Month *"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="6">
              <v-select
                v-model="genForm.year"
                :items="years"
                label="Year *"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12">
              <v-autocomplete
                v-model="genForm.scheme_ids"
                :items="linkedSchemes"
                :item-title="(s) => `${s.scheme_name} (ID: ${s.scheme_id})`"
                item-value="scheme_id"
                label="Schemes"
                variant="outlined"
                density="compact"
                multiple
                chips
                closable-chips
                clearable
                :loading="loadingSchemes"
                :disabled="!genForm.treaty_id"
                :no-data-text="
                  genForm.treaty_id
                    ? 'No schemes linked to this treaty'
                    : 'Select a treaty first'
                "
                hint="Leave empty to include all treaty-linked schemes"
                persistent-hint
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showGenerateDialog = false"
            >Cancel</v-btn
          >
          <v-btn color="primary" :loading="generating" @click="generate"
            >Generate</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Submit Dialog -->
    <v-dialog v-model="showSubmitDialog" max-width="440">
      <v-card>
        <v-card-title>Submit RI Bordereaux</v-card-title>
        <v-card-text>
          <p class="mb-3"
            >Submit run <strong>{{ submittingRun?.run_id }}</strong> to the
            reinsurer?</p
          >
          <v-textarea
            v-model="submitMessage"
            label="Cover Message (optional)"
            variant="outlined"
            density="compact"
            rows="2"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showSubmitDialog = false">Cancel</v-btn>
          <v-btn color="warning" :loading="submitting" @click="performSubmit"
            >Submit</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Detail Dialog -->
    <v-dialog v-model="showDetailDialog" max-width="1000" scrollable>
      <v-card>
        <v-card-title>Run Detail — {{ detailRun?.run_id }}</v-card-title>
        <v-card-text>
          <div style="height: 420px">
            <data-grid
              :row-data="detailRows"
              :column-defs="detailColDefs"
              :default-col-def="{
                resizable: true,
                sortable: true,
                filter: true
              }"
              :pagination="true"
              :pagination-page-size="25"
            />
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showDetailDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Validation Results Dialog -->
    <v-dialog v-model="showValidationDialog" max-width="1100" scrollable>
      <v-card>
        <v-card-title>
          Validation Report — {{ validatingRun?.run_id }}
          <v-chip
            v-if="validationSummary"
            :color="
              validationSummary.status === 'validated' ? 'success' : 'error'
            "
            size="small"
            class="ml-2"
          >
            {{ validationSummary.status }}
          </v-chip>
        </v-card-title>
        <v-card-text>
          <!-- Summary chips -->
          <v-row v-if="validationSummary" class="mb-3">
            <v-col cols="auto">
              <v-chip
                color="error"
                variant="tonal"
                prepend-icon="mdi-alert-circle"
              >
                Critical: {{ validationSummary.critical }}
              </v-chip>
            </v-col>
            <v-col cols="auto">
              <v-chip color="warning" variant="tonal" prepend-icon="mdi-alert">
                Major: {{ validationSummary.major }}
              </v-chip>
            </v-col>
            <v-col cols="auto">
              <v-chip
                color="info"
                variant="tonal"
                prepend-icon="mdi-information"
              >
                Minor: {{ validationSummary.minor }}
              </v-chip>
            </v-col>
            <v-col cols="auto">
              <v-chip color="grey" variant="tonal">
                Total: {{ validationSummary.total }}
              </v-chip>
            </v-col>
          </v-row>
          <div
            v-if="validationSummary && validationSummary.total === 0"
            class="text-center py-6 text-success"
          >
            <v-icon size="48" color="success">mdi-check-circle</v-icon>
            <p class="mt-2"
              >All checks passed — bordereaux is ready for submission.</p
            >
          </div>
          <div v-else style="height: 380px">
            <ag-grid-vue
              class="ag-theme-material"
              style="height: 100%; width: 100%"
              :row-data="validationSummary ? validationSummary.results : []"
              :column-defs="validationColDefs"
              :default-col-def="{
                resizable: true,
                sortable: true,
                filter: true
              }"
              :pagination="true"
              :pagination-page-size="25"
            />
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showValidationDialog = false"
            >Close</v-btn
          >
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
import DataGrid from '@/renderer/components/tables/DataGrid.vue'

const activeTab = ref('member')
const runs = ref([])
const activeTreaties = ref([])
const linkedSchemes = ref([])
const loading = ref(false)
const generating = ref(false)
const loadingSchemes = ref(false)
const submitting = ref(false)
const validating = ref(false)
const stats = ref({
  total_runs: 0,
  member_census: 0,
  claims_runs: 0,
  submitted: 0
})
const filters = ref({ treaty_id: null, status: '' })
const showGenerateDialog = ref(false)
const showSubmitDialog = ref(false)
const showDetailDialog = ref(false)
const showValidationDialog = ref(false)
const submittingRun = ref(null)
const submitMessage = ref('')
const detailRun = ref(null)
const detailRows = ref([])
const detailColDefs = ref([])
const validatingRun = ref(null)
const validationSummary = ref(null)
const snackbar = ref({ show: false, message: '', color: 'success' })

const now = new Date()
const months = Array.from({ length: 12 }, (_, i) => ({
  title: new Date(2000, i).toLocaleString('en', { month: 'long' }),
  value: i + 1
}))
const years = Array.from({ length: 5 }, (_, i) => now.getFullYear() - 2 + i)

const genForm = ref({
  treaty_id: null,
  month: now.getMonth() + 1,
  year: now.getFullYear(),
  scheme_ids: []
})

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

const validationColDefs = [
  {
    field: 'severity',
    headerName: 'Severity',
    width: 110,
    cellRenderer: (p) => {
      const colors = { critical: '#c62828', major: '#e65100', minor: '#1565c0' }
      const c = colors[p.value] || '#9e9e9e'
      return `<span style="padding:2px 10px;border-radius:12px;font-size:12px;background:${c}22;color:${c};font-weight:600">${p.value}</span>`
    }
  },
  {
    field: 'level',
    headerName: 'Level',
    width: 75,
    valueFormatter: (p) => `L${p.value}`
  },
  { field: 'row_type', headerName: 'Row Type', width: 100 },
  { field: 'row_index', headerName: 'Row #', width: 80 },
  { field: 'error_code', headerName: 'Code', width: 90 },
  { field: 'field_name', headerName: 'Field', width: 160 },
  { field: 'message', headerName: 'Message', flex: 1, minWidth: 300 }
]

const formatCurrency = (v) =>
  v != null
    ? 'R ' +
      Number(v).toLocaleString('en-ZA', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      })
    : '—'

const filteredRuns = computed(() =>
  runs.value.filter((r) =>
    activeTab.value === 'member'
      ? r.type === 'member_census'
      : r.type === 'claims_run'
  )
)

const columnDefs = [
  { field: 'run_id', headerName: 'Run ID', width: 160 },
  { field: 'treaty_number', headerName: 'Treaty', width: 130 },
  { field: 'reinsurer_name', headerName: 'Reinsurer', width: 150 },
  { field: 'period_label', headerName: 'Period', width: 120 },
  { field: 'total_lives', headerName: 'Lives', width: 90 },
  { field: 'total_ceded_lives', headerName: 'Ceded Lives', width: 110 },
  {
    field: 'gross_premium',
    headerName: 'Gross Premium',
    width: 140,
    valueFormatter: (p) => formatCurrency(p.value)
  },
  {
    field: 'ceded_premium',
    headerName: 'Ceded Premium',
    width: 140,
    valueFormatter: (p) => formatCurrency(p.value)
  },
  {
    field: 'status',
    headerName: 'Status',
    width: 120,
    cellRenderer: (p) => {
      const color = statusColor[p.value] || '#9e9e9e'
      return `<span style="padding:2px 10px;border-radius:12px;font-size:12px;background:${color}22;color:${color};font-weight:500">${p.value}</span>`
    }
  },
  { field: 'run_version', headerName: 'Version', width: 80 },
  { field: 'generated_by', headerName: 'Generated By', width: 140 },
  {
    headerName: 'Actions',
    width: 80,
    pinned: 'right',
    sortable: false,
    filter: false,
    cellRenderer: (p) => {
      const key = p.data.run_id.replace(/-/g, '_')
      window[`showRunMenu_${key}`] = (event) => showContextMenu(event, p.data)
      return `<div style="display:flex;align-items:center;justify-content:center;height:100%">
        <button onclick="showRunMenu_${key}(event)" title="Actions" style="background:none;border:none;cursor:pointer;padding:4px 8px;border-radius:4px;color:#616161;display:flex;align-items:center;justify-content:center">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="5" r="2"/><circle cx="12" cy="12" r="2"/><circle cx="12" cy="19" r="2"/></svg>
        </button>
      </div>`
    }
  }
]

let activeMenuCleanup = null

function showContextMenu(event, data) {
  if (activeMenuCleanup) activeMenuCleanup()

  const canValidate =
    data.status === 'generated' || data.status === 'validation_failed'
  const hasResults = ['validating', 'validated', 'validation_failed'].includes(
    data.status
  )
  const canSubmit = data.status === 'validated'
  const canAck = data.status === 'submitted'

  const menuItems = [
    { label: 'View', color: '#1976d2', fn: () => viewDetail(data) },
    canValidate
      ? { label: 'Validate', color: '#7b1fa2', fn: () => runValidation(data) }
      : null,
    hasResults
      ? {
          label: 'Results',
          color: '#616161',
          fn: () => openValidationResults(data)
        }
      : null,
    canSubmit
      ? { label: 'Submit', color: '#f57c00', fn: () => openSubmitDialog(data) }
      : null,
    canAck
      ? { label: 'Acknowledge', color: '#388e3c', fn: () => acknowledge(data) }
      : null
  ].filter(Boolean)

  const menu = document.createElement('div')
  menu.style.cssText =
    'position:fixed;background:#fff;border:1px solid #e0e0e0;border-radius:8px;' +
    'box-shadow:0 4px 16px rgba(0,0,0,0.14);z-index:9999;min-width:160px;padding:4px 0;'

  menuItems.forEach(({ label, color, fn }) => {
    const item = document.createElement('div')
    item.textContent = label
    item.style.cssText = `padding:8px 16px;cursor:pointer;font-size:13px;color:${color};`
    item.addEventListener(
      'mouseenter',
      () => (item.style.background = '#f5f5f5')
    )
    item.addEventListener('mouseleave', () => (item.style.background = ''))
    item.addEventListener('click', () => {
      cleanup()
      fn()
    })
    menu.appendChild(item)
  })

  document.body.appendChild(menu)

  const btn = event.currentTarget || event.target
  const rect = btn.getBoundingClientRect()
  menu.style.top = `${rect.bottom + 4}px`
  menu.style.left = `${rect.left}px`

  const mr = menu.getBoundingClientRect()
  if (mr.right > window.innerWidth - 8)
    menu.style.left = `${rect.right - mr.width}px`
  if (mr.bottom > window.innerHeight - 8)
    menu.style.top = `${rect.top - mr.height - 4}px`

  function cleanup() {
    menu.remove()
    document.removeEventListener('click', outsideClick, true)
    activeMenuCleanup = null
  }
  activeMenuCleanup = cleanup

  function outsideClick(e) {
    if (!menu.contains(e.target) && e.target !== btn) cleanup()
  }
  setTimeout(() => document.addEventListener('click', outsideClick, true), 0)
}

const notify = (message, color = 'success') => {
  snackbar.value = { show: true, message, color }
}

async function loadRuns() {
  loading.value = true
  try {
    const params = {}
    if (filters.value.treaty_id) params.treaty_id = filters.value.treaty_id
    if (filters.value.status) params.status = filters.value.status
    const res = await GroupPricingService.getRIBordereauxRuns(params)
    runs.value = res.data?.data || []
  } catch {
    notify('Failed to load runs', 'error')
  } finally {
    loading.value = false
  }
}

async function loadStats() {
  try {
    const res = await GroupPricingService.getRIBordereauxStats({})
    stats.value = res.data?.data || stats.value
  } catch {}
}

async function loadActiveTreaties() {
  try {
    const res = await GroupPricingService.getTreaties({ status: 'active' })
    activeTreaties.value = res.data?.data || []
  } catch {}
}

function openGenerateDialog() {
  genForm.value = {
    treaty_id: null,
    month: now.getMonth() + 1,
    year: now.getFullYear(),
    scheme_ids: []
  }
  linkedSchemes.value = []
  showGenerateDialog.value = true
}

async function onGenTreatySelected(treatyId) {
  genForm.value.scheme_ids = []
  linkedSchemes.value = []
  if (!treatyId) return
  loadingSchemes.value = true
  try {
    const res = await GroupPricingService.getTreatySchemeLinks(treatyId)
    linkedSchemes.value = res.data?.data || []
  } catch {
    notify('Failed to load linked schemes', 'error')
  } finally {
    loadingSchemes.value = false
  }
}

async function generate() {
  if (!genForm.value.treaty_id) {
    notify('Please select a treaty', 'warning')
    return
  }
  generating.value = true
  try {
    const payload = {
      treaty_id: genForm.value.treaty_id,
      month: genForm.value.month,
      year: genForm.value.year,
      scheme_ids: genForm.value.scheme_ids
    }
    if (activeTab.value === 'member') {
      await GroupPricingService.generateRIMemberBordereaux(payload)
    } else {
      await GroupPricingService.generateRIClaimsBordereaux(payload)
    }
    notify('RI bordereaux generated successfully')
    showGenerateDialog.value = false
    loadRuns()
    loadStats()
  } catch (e) {
    notify(e.response?.data?.error || 'Generation failed', 'error')
  } finally {
    generating.value = false
  }
}

function openSubmitDialog(run) {
  submittingRun.value = run
  submitMessage.value = ''
  showSubmitDialog.value = true
}

async function performSubmit() {
  submitting.value = true
  try {
    await GroupPricingService.submitRIBordereaux({
      run_id: submittingRun.value.run_id,
      message: submitMessage.value
    })
    notify('Submitted successfully')
    showSubmitDialog.value = false
    loadRuns()
    loadStats()
  } catch (e) {
    notify(e.response?.data?.error || 'Submission failed', 'error')
  } finally {
    submitting.value = false
  }
}

async function acknowledge(run) {
  try {
    await GroupPricingService.acknowledgeRIBordereaux(run.run_id)
    notify('Marked as acknowledged')
    loadRuns()
    loadStats()
  } catch {
    notify('Failed to acknowledge', 'error')
  }
}

async function runValidation(run) {
  validating.value = true
  notify(`Validating ${run.run_id}…`, 'info')
  try {
    const res = await GroupPricingService.validateRIBordereaux(run.run_id)
    const summary = res.data?.data
    validatingRun.value = run
    validationSummary.value = summary
    showValidationDialog.value = true
    loadRuns()
    loadStats()
    if (summary.status === 'validated') {
      notify(
        `Validation passed — ${run.run_id} is ready for submission`,
        'success'
      )
    } else {
      notify(
        `Validation failed — ${summary.critical} critical issue(s) found`,
        'error'
      )
    }
  } catch (e) {
    notify(e.response?.data?.error || 'Validation failed', 'error')
  } finally {
    validating.value = false
  }
}

async function openValidationResults(run) {
  validatingRun.value = run
  validationSummary.value = null
  showValidationDialog.value = true
  try {
    const res = await GroupPricingService.getRIValidationResults(run.run_id)
    validationSummary.value = res.data?.data
  } catch {
    notify('Failed to load validation results', 'error')
  }
}

async function viewDetail(run) {
  detailRun.value = run
  detailRows.value = []
  showDetailDialog.value = true
  try {
    if (run.type === 'member_census') {
      const res = await GroupPricingService.getRIBordereauxMemberRows(
        run.run_id
      )
      const rows = res.data?.data || []
      detailRows.value = rows
      detailColDefs.value = [
        { field: 'member_id_number', headerName: 'ID Number', width: 140 },
        { field: 'member_name', headerName: 'Name', flex: 1 },
        { field: 'scheme_name', headerName: 'Scheme', width: 140 },
        { field: 'benefit_code', headerName: 'Benefit', width: 100 },
        { field: 'policy_type', headerName: 'Policy Type', width: 110 },
        { field: 'treaty_section', headerName: 'Treaty Section', width: 130 },
        {
          field: 'endorsement_number',
          headerName: 'Endorsement No.',
          width: 140
        },
        {
          field: 'sum_assured',
          headerName: 'Sum Assured',
          width: 130,
          valueFormatter: (p) => formatCurrency(p.value)
        },
        {
          field: 'ceded_amount',
          headerName: 'Ceded SA',
          width: 120,
          valueFormatter: (p) => formatCurrency(p.value)
        },
        {
          field: 'ceded_premium',
          headerName: 'Ceded Prem',
          width: 120,
          valueFormatter: (p) => formatCurrency(p.value)
        },
        {
          field: 'cumulative_premium_ytd',
          headerName: 'Cumul. Prem YTD',
          width: 150,
          valueFormatter: (p) => formatCurrency(p.value)
        },
        {
          field: 'period_adjustment',
          headerName: 'Period Adj.',
          width: 120,
          valueFormatter: (p) => formatCurrency(p.value)
        },
        {
          field: 'exchange_rate',
          headerName: 'FX Rate',
          width: 90,
          valueFormatter: (p) => p.value?.toFixed(4) ?? '—'
        },
        {
          field: 'sanctions_screening_status',
          headerName: 'Sanctions',
          width: 110,
          cellRenderer: (p) => {
            const colors = {
              cleared: '#388e3c',
              pending: '#f57c00',
              flagged: '#d32f2f'
            }
            const c = colors[p.value] || '#9e9e9e'
            return `<span style="padding:2px 8px;border-radius:12px;font-size:11px;background:${c}22;color:${c};font-weight:500">${p.value}</span>`
          }
        },
        { field: 'cession_basis', headerName: 'Basis', width: 110 },
        { field: 'change_type', headerName: 'Change', width: 90 }
      ]
    } else {
      const res = await GroupPricingService.getRIBordereauxClaimsRows(
        run.run_id
      )
      detailRows.value = res.data?.data || []
      detailColDefs.value = [
        { field: 'claim_number', headerName: 'Claim No.', width: 130 },
        {
          field: 'reinsurer_claim_reference',
          headerName: 'RI Claim Ref',
          width: 130
        },
        { field: 'member_name', headerName: 'Member', flex: 1 },
        { field: 'scheme_name', headerName: 'Scheme', width: 140 },
        { field: 'benefit_code', headerName: 'Benefit', width: 100 },
        { field: 'cause_of_loss', headerName: 'Cause of Loss', width: 130 },
        {
          field: 'gross_claim_amount',
          headerName: 'Gross Claim',
          width: 130,
          valueFormatter: (p) => formatCurrency(p.value)
        },
        {
          field: 'gross_paid_losses',
          headerName: 'Gross Paid',
          width: 120,
          valueFormatter: (p) => formatCurrency(p.value)
        },
        {
          field: 'gross_outstanding_reserve',
          headerName: 'Outstanding Res.',
          width: 140,
          valueFormatter: (p) => formatCurrency(p.value)
        },
        {
          field: 'ceded_claim_amount',
          headerName: 'Ceded',
          width: 120,
          valueFormatter: (p) => formatCurrency(p.value)
        },
        {
          field: 'excess_retention',
          headerName: 'Retention',
          width: 120,
          valueFormatter: (p) => formatCurrency(p.value)
        },
        {
          field: 'recoveries',
          headerName: 'Recoveries',
          width: 120,
          valueFormatter: (p) => formatCurrency(p.value)
        },
        {
          field: 'prior_period_movement',
          headerName: 'Prior Period Mvt',
          width: 140,
          valueFormatter: (p) => formatCurrency(p.value)
        },
        {
          field: 'cumulative_loss_ytd',
          headerName: 'Cumul. Loss YTD',
          width: 140,
          valueFormatter: (p) => formatCurrency(p.value)
        },
        {
          field: 'ibnr_flag',
          headerName: 'IBNR',
          width: 80,
          cellRenderer: (p) =>
            p.value
              ? '<span style="color:#f57c00;font-weight:500">IBNR</span>'
              : '—'
        },
        {
          field: 'large_loss_flag',
          headerName: 'Large Loss',
          width: 100,
          cellRenderer: (p) =>
            p.value
              ? '<span style="color:#d32f2f;font-weight:500">Large</span>'
              : '—'
        },
        { field: 'catastrophe_event_code', headerName: 'Cat Code', width: 100 },
        { field: 'claim_status', headerName: 'Status', width: 110 }
      ]
    }
  } catch {
    notify('Failed to load detail rows', 'error')
  }
}

onMounted(() => {
  loadRuns()
  loadStats()
  loadActiveTreaties()
})
</script>
