<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center justify-space-between">
              <span class="headline"
                >RI Claims Notifications &amp; Cat Register</span
              >
              <v-btn
                size="small"
                rounded
                color="white"
                variant="outlined"
                prepend-icon="mdi-arrow-left"
                @click="$router.push('/group-pricing/bordereaux-management')"
              >
                Back
              </v-btn>
            </div>
          </template>

          <template #default>
            <v-tabs v-model="activeTab" class="mb-4">
              <v-tab value="notices">Large Claim Notifications</v-tab>
              <v-tab value="cat">Cat Event Register</v-tab>
            </v-tabs>

            <v-tabs-window v-model="activeTab">
              <!-- ======= TAB 1: Large Claim Notifications ======= -->
              <v-tabs-window-item value="notices">
                <!-- Stats chips -->
                <v-row class="mb-4">
                  <v-col cols="auto">
                    <v-chip color="primary" variant="tonal"
                      >Total: {{ stats.total }}</v-chip
                    >
                  </v-col>
                  <v-col cols="auto">
                    <v-chip color="warning" variant="tonal"
                      >Pending: {{ stats.pending }}</v-chip
                    >
                  </v-col>
                  <v-col cols="auto">
                    <v-chip color="blue" variant="tonal"
                      >Sent: {{ stats.sent }}</v-chip
                    >
                  </v-col>
                  <v-col cols="auto">
                    <v-chip color="success" variant="tonal"
                      >Acknowledged: {{ stats.acknowledged }}</v-chip
                    >
                  </v-col>
                  <v-col cols="auto">
                    <v-chip color="error" variant="tonal"
                      >Late: {{ stats.late }}</v-chip
                    >
                  </v-col>
                  <v-col cols="auto">
                    <v-chip color="orange" variant="tonal"
                      >Queried: {{ stats.queried }}</v-chip
                    >
                  </v-col>
                </v-row>

                <!-- Filters -->
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
                          @update:model-value="loadNotices"
                        />
                      </v-col>
                      <v-col cols="12" sm="3">
                        <v-text-field
                          v-model="filters.scheme_id"
                          label="Scheme ID"
                          variant="outlined"
                          density="compact"
                          clearable
                          type="number"
                          @update:model-value="loadNotices"
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
                          @update:model-value="loadNotices"
                        />
                      </v-col>
                      <v-col cols="auto" class="d-flex align-center gap-2">
                        <v-btn
                          size="small"
                          rounded
                          color="success"
                          prepend-icon="mdi-refresh"
                          :loading="loading"
                          @click="loadNotices"
                        >
                          Refresh
                        </v-btn>
                        <v-btn
                          size="small"
                          rounded
                          color="deep-orange"
                          prepend-icon="mdi-radar"
                          @click="showMonitorDialog = true"
                        >
                          Run Monitor
                        </v-btn>
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>

                <!-- Grid -->
                <div style="height: 500px">
                  <ag-grid-vue
                    class="ag-theme-material"
                    style="height: 100%; width: 100%"
                    :row-data="notices"
                    :column-defs="columnDefs"
                    :default-col-def="{
                      resizable: true,
                      sortable: true,
                      filter: true
                    }"
                    :pagination="true"
                    :pagination-page-size="25"
                    :overlay-loading-template="loadingOverlay"
                    :overlay-no-rows-template="`<span class='ag-overlay-no-rows-center'>No large-claim notices in this period.</span>`"
                    @grid-ready="onNoticesGridReady"
                  />
                </div>
              </v-tabs-window-item>

              <!-- ======= TAB 2: Cat Event Register ======= -->
              <v-tabs-window-item value="cat">
                <!-- Cat KPI chips -->
                <v-row class="mb-4">
                  <v-col cols="auto">
                    <v-chip color="primary" variant="tonal"
                      >Total Cat Rows: {{ catRows.length }}</v-chip
                    >
                  </v-col>
                  <v-col cols="auto">
                    <v-chip color="error" variant="tonal"
                      >Large Loss:
                      {{
                        catRows.filter((r) => r.large_loss_flag).length
                      }}</v-chip
                    >
                  </v-col>
                  <v-col cols="auto">
                    <v-chip color="warning" variant="tonal"
                      >Events: {{ uniqueCatCodes.length }}</v-chip
                    >
                  </v-col>
                </v-row>

                <!-- Cat Filters -->
                <v-card variant="outlined" class="mb-4">
                  <v-card-text>
                    <v-row>
                      <v-col cols="12" sm="3">
                        <v-text-field
                          v-model="catFilters.cat_event_code"
                          label="Cat Event Code"
                          variant="outlined"
                          density="compact"
                          clearable
                        />
                      </v-col>
                      <v-col cols="12" sm="3">
                        <v-select
                          v-model="catFilters.treaty_id"
                          :items="activeTreaties"
                          item-title="treaty_name"
                          item-value="id"
                          label="Treaty"
                          variant="outlined"
                          density="compact"
                          clearable
                        />
                      </v-col>
                      <v-col cols="12" sm="2">
                        <v-text-field
                          v-model="catFilters.period_from"
                          label="Period From"
                          type="date"
                          variant="outlined"
                          density="compact"
                        />
                      </v-col>
                      <v-col cols="12" sm="2">
                        <v-text-field
                          v-model="catFilters.period_to"
                          label="Period To"
                          type="date"
                          variant="outlined"
                          density="compact"
                        />
                      </v-col>
                      <v-col cols="auto" class="d-flex align-center gap-2">
                        <v-btn
                          size="small"
                          rounded
                          color="primary"
                          prepend-icon="mdi-magnify"
                          :loading="catLoading"
                          @click="loadCatRows"
                        >
                          Search
                        </v-btn>
                        <v-btn
                          size="small"
                          rounded
                          variant="outlined"
                          @click="resetCatFilters"
                        >
                          Reset
                        </v-btn>
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>

                <!-- Cat Grid -->
                <div style="height: 500px">
                  <ag-grid-vue
                    class="ag-theme-material"
                    style="height: 100%; width: 100%"
                    :row-data="catRows"
                    :column-defs="catColumnDefs"
                    :default-col-def="{
                      resizable: true,
                      sortable: true,
                      filter: true
                    }"
                    :pagination="true"
                    :pagination-page-size="25"
                    :overlay-loading-template="loadingOverlay"
                    :overlay-no-rows-template="`<span class='ag-overlay-no-rows-center'>No cat events recorded for these filters.</span>`"
                    @grid-ready="onCatGridReady"
                  />
                </div>
              </v-tabs-window-item>
            </v-tabs-window>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Monitor Dialog -->
    <v-dialog v-model="showMonitorDialog" max-width="440">
      <v-card>
        <v-card-title>Monitor Large Claims</v-card-title>
        <v-card-text>
          <p class="text-body-2 mb-3 text-medium-emphasis">
            Scans all claims for the selected treaty's linked schemes and
            creates notification records for any claims above the large claims
            threshold.
          </p>
          <v-select
            v-model="monitorTreatyId"
            :items="activeTreaties"
            item-title="treaty_name"
            item-value="id"
            label="Treaty *"
            variant="outlined"
            density="compact"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showMonitorDialog = false"
            >Cancel</v-btn
          >
          <v-btn color="deep-orange" :loading="monitoring" @click="runMonitor"
            >Run Monitor</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Update Notice Dialog -->
    <v-dialog v-model="showUpdateDialog" max-width="500">
      <v-card>
        <v-card-title
          >Update Notification —
          {{ updatingNotice?.claim_number }}</v-card-title
        >
        <v-card-text>
          <v-row>
            <v-col cols="12">
              <v-select
                v-model="updateForm.status"
                :items="statusOptions"
                label="Status"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="updateForm.query_details"
                label="Query Details"
                variant="outlined"
                density="compact"
                rows="2"
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="updateForm.response_notes"
                label="Response Notes"
                variant="outlined"
                density="compact"
                rows="2"
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showUpdateDialog = false">Cancel</v-btn>
          <v-btn color="primary" :loading="updating" @click="saveUpdate"
            >Save</v-btn
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
import { ref, computed, onMounted, watch } from 'vue'
import { AgGridVue } from 'ag-grid-vue3'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import BaseCard from '@/renderer/components/BaseCard.vue'

const activeTab = ref('notices')

// AG-Grid overlay plumbing: the templates render when the grid calls
// showLoadingOverlay / showNoRowsOverlay. The per-grid watchers below flip
// between overlays based on the matching loading ref.
const loadingOverlay = '<span class="ag-overlay-loading-center">Loading…</span>'
const noticesGridApi = ref(null)
const catGridApi = ref(null)
const onNoticesGridReady = (p) => {
  noticesGridApi.value = p.api
  applyNoticesOverlay()
}
const onCatGridReady = (p) => {
  catGridApi.value = p.api
  applyCatOverlay()
}
const applyNoticesOverlay = () => {
  if (!noticesGridApi.value) return
  if (loading.value) noticesGridApi.value.showLoadingOverlay()
  else if (!notices.value?.length) noticesGridApi.value.showNoRowsOverlay()
  else noticesGridApi.value.hideOverlay()
}
const applyCatOverlay = () => {
  if (!catGridApi.value) return
  if (catLoading.value) catGridApi.value.showLoadingOverlay()
  else if (!catRows.value?.length) catGridApi.value.showNoRowsOverlay()
  else catGridApi.value.hideOverlay()
}

const notices = ref([])
const activeTreaties = ref([])
const loading = ref(false)
const monitoring = ref(false)
const updating = ref(false)
const stats = ref({
  total: 0,
  pending: 0,
  sent: 0,
  acknowledged: 0,
  late: 0,
  queried: 0
})
const filters = ref({ treaty_id: null, scheme_id: '', status: '' })
const showMonitorDialog = ref(false)
const showUpdateDialog = ref(false)
const monitorTreatyId = ref(null)
const updatingNotice = ref(null)
const updateForm = ref({ status: '', query_details: '', response_notes: '' })
const snackbar = ref({ show: false, message: '', color: 'success' })

watch(loading, applyNoticesOverlay)
watch(notices, applyNoticesOverlay, { deep: true })

// Cat Event Register state
const catRows = ref([])
const catLoading = ref(false)
watch(catLoading, applyCatOverlay)
watch(catRows, applyCatOverlay, { deep: true })
const catFilters = ref({
  cat_event_code: '',
  treaty_id: null,
  period_from: '',
  period_to: ''
})
const uniqueCatCodes = computed(() => [
  ...new Set(catRows.value.map((r) => r.catastrophe_event_code).filter(Boolean))
])

const statusOptions = [
  { title: 'Pending', value: 'pending' },
  { title: 'Sent', value: 'sent' },
  { title: 'Acknowledged', value: 'acknowledged' },
  { title: 'Late', value: 'late' },
  { title: 'Queried', value: 'queried' }
]

const statusColor = {
  pending: '#f57c00',
  sent: '#1976d2',
  acknowledged: '#388e3c',
  late: '#d32f2f',
  queried: '#ef6c00'
}

const formatCurrency = (v) =>
  v != null
    ? 'R ' + Number(v).toLocaleString('en-ZA', { minimumFractionDigits: 0 })
    : '—'

const slaDaysColor = (dueDate) => {
  if (!dueDate) return '#9e9e9e'
  const diff = Math.ceil((new Date(dueDate) - new Date()) / 86400000)
  if (diff < 0) return '#d32f2f'
  if (diff <= 7) return '#f57c00'
  return '#388e3c'
}

const columnDefs = [
  { field: 'treaty_number', headerName: 'Treaty', width: 130 },
  { field: 'claim_number', headerName: 'Claim No.', width: 130 },
  { field: 'scheme_name', headerName: 'Scheme', flex: 1, minWidth: 140 },
  { field: 'reinsurer_name', headerName: 'Reinsurer', width: 140 },
  { field: 'event_date', headerName: 'Event Date', width: 110 },
  {
    field: 'due_date',
    headerName: 'Due Date / SLA',
    width: 140,
    cellRenderer: (p) => {
      const color = slaDaysColor(p.value)
      const diff = p.value
        ? Math.ceil((new Date(p.value) - new Date()) / 86400000)
        : null
      const label =
        diff !== null
          ? `${p.value} (${diff > 0 ? diff + 'd left' : Math.abs(diff) + 'd late'})`
          : p.value
      return `<span style="color:${color};font-size:12px;font-weight:500">${label || '—'}</span>`
    }
  },
  {
    field: 'gross_claim_amount',
    headerName: 'Gross Claim',
    width: 130,
    valueFormatter: (p) => formatCurrency(p.value)
  },
  {
    field: 'estimated_ceded_amount',
    headerName: 'Est. Ceded',
    width: 120,
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
  {
    field: 'late_flag',
    headerName: 'Late',
    width: 70,
    cellRenderer: (p) => (p.value ? '⚠️' : '')
  },
  {
    headerName: 'Actions',
    width: 200,
    sortable: false,
    filter: false,
    cellRenderer: (p) => {
      const id = p.data.id
      window[`updateNotice_${id}`] = () => openUpdateDialog(p.data)
      window[`markSentNotice_${id}`] = () => quickMark(p.data, 'sent')
      window[`markAckNotice_${id}`] = () => quickMark(p.data, 'acknowledged')
      const isPending = p.data.status === 'pending' || p.data.status === 'late'
      const isSent = p.data.status === 'sent'
      return `<div style="display:flex;gap:4px;align-items:center;height:100%">
        ${isPending ? `<button onclick="markSentNotice_${id}()" style="padding:2px 8px;border-radius:4px;border:1px solid #1976d2;color:#1976d2;background:none;cursor:pointer;font-size:11px">Mark Sent</button>` : ''}
        ${isSent ? `<button onclick="markAckNotice_${id}()" style="padding:2px 8px;border-radius:4px;border:1px solid #388e3c;color:#388e3c;background:none;cursor:pointer;font-size:11px">Acknowledge</button>` : ''}
        <button onclick="updateNotice_${id}()" style="padding:2px 8px;border-radius:4px;border:1px solid #616161;color:#616161;background:none;cursor:pointer;font-size:11px">Update</button>
      </div>`
    }
  }
]

const catColumnDefs = [
  {
    field: 'catastrophe_event_code',
    headerName: 'Cat Event Code',
    width: 150,
    cellRenderer: (p) =>
      `<span style="padding:2px 10px;border-radius:12px;font-size:12px;background:#ef5350aa;color:#fff;font-weight:600">${p.value || '—'}</span>`
  },
  { field: 'run_id', headerName: 'Run ID', width: 120 },
  { field: 'claim_number', headerName: 'Claim No.', width: 130 },
  { field: 'member_id_number', headerName: 'Member ID', width: 130 },
  { field: 'member_name', headerName: 'Member', flex: 1, minWidth: 140 },
  { field: 'date_of_event', headerName: 'Event Date', width: 110 },
  { field: 'benefit_code', headerName: 'Benefit', width: 100 },
  {
    field: 'gross_claim_amount',
    headerName: 'Gross Claim',
    width: 130,
    valueFormatter: (p) => formatCurrency(p.value)
  },
  {
    field: 'ceded_claim_amount',
    headerName: 'Ceded Claim',
    width: 130,
    valueFormatter: (p) => formatCurrency(p.value)
  },
  {
    field: 'large_loss_flag',
    headerName: 'Large Loss',
    width: 100,
    cellRenderer: (p) =>
      p.value
        ? '<span style="color:#d32f2f;font-weight:600">⚠ Yes</span>'
        : '<span style="color:#9e9e9e">—</span>'
  },
  {
    field: 'gross_paid_losses',
    headerName: 'Paid Losses',
    width: 120,
    valueFormatter: (p) => (p.value ? formatCurrency(p.value) : '—')
  },
  {
    field: 'gross_outstanding_reserve',
    headerName: 'Outstanding Res.',
    width: 140,
    valueFormatter: (p) => (p.value ? formatCurrency(p.value) : '—')
  }
]

const notify = (message, color = 'success') => {
  snackbar.value = { show: true, message, color }
}

async function loadNotices() {
  loading.value = true
  try {
    const params = {}
    if (filters.value.treaty_id) params.treaty_id = filters.value.treaty_id
    if (filters.value.scheme_id) params.scheme_id = filters.value.scheme_id
    if (filters.value.status) params.status = filters.value.status
    const res = await GroupPricingService.getLargeClaimNotices(params)
    notices.value = res.data?.data || []
  } catch {
    notify('Failed to load notices', 'error')
  } finally {
    loading.value = false
  }
}

async function loadStats() {
  try {
    const params = {}
    if (filters.value.treaty_id) params.treaty_id = filters.value.treaty_id
    const res = await GroupPricingService.getLargeClaimStats(params)
    stats.value = res.data?.data || stats.value
  } catch {}
}

async function loadActiveTreaties() {
  try {
    const res = await GroupPricingService.getTreaties({ status: 'active' })
    activeTreaties.value = res.data?.data || []
  } catch {}
}

async function runMonitor() {
  if (!monitorTreatyId.value) {
    notify('Please select a treaty', 'warning')
    return
  }
  monitoring.value = true
  try {
    const res = await GroupPricingService.monitorLargeClaims({
      treaty_id: monitorTreatyId.value
    })
    const count = res.data?.data?.created || 0
    notify(`Monitor complete — ${count} new notice(s) created`)
    showMonitorDialog.value = false
    loadNotices()
    loadStats()
  } catch (e) {
    notify(e.response?.data?.error || 'Monitor failed', 'error')
  } finally {
    monitoring.value = false
  }
}

function openUpdateDialog(notice) {
  updatingNotice.value = notice
  updateForm.value = {
    status: notice.status,
    query_details: notice.query_details,
    response_notes: notice.response_notes
  }
  showUpdateDialog.value = true
}

async function saveUpdate() {
  updating.value = true
  try {
    await GroupPricingService.updateLargeClaimNotice(
      updatingNotice.value.id,
      updateForm.value
    )
    notify('Notice updated')
    showUpdateDialog.value = false
    loadNotices()
    loadStats()
  } catch {
    notify('Failed to update', 'error')
  } finally {
    updating.value = false
  }
}

async function quickMark(notice, status) {
  try {
    await GroupPricingService.updateLargeClaimNotice(notice.id, { status })
    notify(`Marked as ${status}`)
    loadNotices()
    loadStats()
  } catch {
    notify('Failed to update', 'error')
  }
}

async function loadCatRows() {
  catLoading.value = true
  try {
    const params = {}
    if (catFilters.value.cat_event_code)
      params.cat_event_code = catFilters.value.cat_event_code
    if (catFilters.value.treaty_id)
      params.treaty_id = catFilters.value.treaty_id
    if (catFilters.value.period_from)
      params.period_from = catFilters.value.period_from
    if (catFilters.value.period_to)
      params.period_to = catFilters.value.period_to
    const res = await GroupPricingService.getCatastropheClaimsRows(params)
    catRows.value = res.data?.data || []
  } catch {
    notify('Failed to load cat event rows', 'error')
  } finally {
    catLoading.value = false
  }
}

function resetCatFilters() {
  catFilters.value = {
    cat_event_code: '',
    treaty_id: null,
    period_from: '',
    period_to: ''
  }
  loadCatRows()
}

onMounted(() => {
  loadNotices()
  loadStats()
  loadActiveTreaties()
  loadCatRows()
})
</script>
