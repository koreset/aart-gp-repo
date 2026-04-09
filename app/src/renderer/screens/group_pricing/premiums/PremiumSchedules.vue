<template>
  <v-container fluid>
    <base-card :show-actions="false">
      <template #header>
        <h3 class="mb-0">Premium Schedules</h3>
      </template>
      <template #default>
        <!-- Filter Bar -->
        <v-row class="mb-3" align="center">
          <v-col cols="12" md="3">
            <v-select
              v-model="filters.schemeId"
              label="Scheme"
              :items="inForceSchemes"
              item-title="name"
              item-value="id"
              variant="outlined"
              density="compact"
              clearable
              prepend-inner-icon="mdi-office-building-outline"
              :loading="schemesLoading"
            />
          </v-col>
          <v-col cols="12" md="3">
            <v-select
              v-model.number="filters.month"
              label="Month"
              :items="months"
              item-title="label"
              item-value="value"
              variant="outlined"
              density="compact"
              clearable
            />
          </v-col>
          <v-col cols="12" md="3">
            <v-text-field
              v-model.number="filters.year"
              label="Year"
              type="number"
              variant="outlined"
              density="compact"
              clearable
            />
          </v-col>
          <v-col cols="12" md="3">
            <v-select
              v-model="filters.status"
              label="Status"
              :items="statusOptions"
              variant="outlined"
              density="compact"
              clearable
            />
          </v-col>
        </v-row>
        <v-row class="mt-n9 mb-8">
          <v-col cols="12" md="12" class="d-flex justify-end">
            <v-btn
              variant="outlined"
              rounded
              size="small"
              class="mr-2"
              prepend-icon="mdi-filter-remove"
              @click="
                () => {
                  resetFilters()
                  loadSchedules()
                }
              "
            >
              Clear Filters
            </v-btn>
            <v-btn
              variant="outlined"
              rounded
              size="small"
              class="mr-2"
              color="primary"
              prepend-icon="mdi-refresh"
              @click="loadSchedules"
            >
              Refresh
            </v-btn>
            <v-btn
              v-if="hasPermission('premiums:generate_schedule')"
              rounded
              size="small"
              class="mr-2"
              color="primary"
              prepend-icon="mdi-plus"
              @click="openGenerateDialog"
            >
              Generate
            </v-btn>
            <v-btn
              variant="outlined"
              rounded
              size="small"
              color="secondary"
              prepend-icon="mdi-grid"
              @click="openCoverageOverview"
            >
              Coverage Overview
            </v-btn>
          </v-col>
        </v-row>
        <!-- Schedules Grid -->
        <v-row>
          <v-col cols="12">
            <v-card variant="outlined">
              <v-card-text class="pa-0">
                <ag-grid-vue
                  class="ag-theme-balham"
                  :style="{ height: gridHeight, width: '100%' }"
                  :column-defs="columnDefs"
                  :row-data="schedules"
                  :default-col-def="defaultColDef"
                  :loading="loading"
                  row-selection="single"
                  @row-clicked="onRowClicked"
                />
                <empty-state
                  v-if="!loading && schedules.length === 0"
                  title="No schedules found"
                  message="Generate a premium schedule or adjust your filters."
                  icon="mdi-calendar-month-outline"
                />
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </template>
    </base-card>

    <!-- Generate Schedule Dialog -->
    <v-dialog v-model="generateDialog" max-width="500" persistent>
      <v-card>
        <v-card-title>Generate Premium Schedule</v-card-title>
        <v-card-text>
          <!-- All-schemes toggle -->
          <v-switch
            v-model="genForm.allSchemes"
            label="Generate for all in-force schemes"
            color="primary"
            density="compact"
            class="mb-2"
            hide-details
          />

          <v-alert
            v-if="genForm.allSchemes"
            type="info"
            variant="tonal"
            class="mb-4"
          >
            A schedule will be generated for every in-force scheme. Schemes that
            already have a schedule for the selected period will be skipped
            automatically.
          </v-alert>
          <v-alert v-else type="info" variant="tonal" class="mb-4">
            Select a scheme and period to generate the monthly premium schedule.
          </v-alert>

          <v-row>
            <v-col v-if="!genForm.allSchemes" cols="12">
              <v-select
                v-model="genForm.schemeId"
                label="Scheme *"
                :items="inForceSchemes"
                item-title="name"
                item-value="id"
                variant="outlined"
                density="compact"
                :loading="schemesLoading"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-select
                v-model.number="genForm.month"
                label="Month *"
                :items="months"
                item-title="label"
                item-value="value"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model.number="genForm.year"
                label="Year *"
                type="number"
                variant="outlined"
                density="compact"
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="generateDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="generating"
            :disabled="
              (!genForm.allSchemes && !genForm.schemeId) ||
              !genForm.month ||
              !genForm.year
            "
            @click="handleGenerate"
          >
            Generate
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Bulk Results Dialog -->
    <v-dialog v-model="resultsDialog" max-width="680" scrollable>
      <v-card>
        <v-card-title>Bulk Generation Results</v-card-title>
        <v-card-text>
          <!-- Summary chips -->
          <v-row class="mb-4" dense>
            <v-col cols="auto">
              <v-chip color="primary" variant="tonal" size="small">
                Total: {{ bulkResult?.total ?? 0 }}
              </v-chip>
            </v-col>
            <v-col cols="auto">
              <v-chip color="success" variant="tonal" size="small">
                Succeeded: {{ bulkResult?.succeeded ?? 0 }}
              </v-chip>
            </v-col>
            <v-col cols="auto">
              <v-chip color="warning" variant="tonal" size="small">
                Skipped: {{ bulkResult?.skipped ?? 0 }}
              </v-chip>
            </v-col>
            <v-col cols="auto">
              <v-chip color="error" variant="tonal" size="small">
                Failed: {{ bulkResult?.failed ?? 0 }}
              </v-chip>
            </v-col>
          </v-row>

          <!-- Per-scheme results list -->
          <v-list density="compact" lines="one">
            <v-list-item
              v-for="row in bulkResult?.results ?? []"
              :key="row.scheme_id"
              :prepend-icon="statusIcon(row.status)"
              :base-color="statusResultColor(row.status)"
            >
              <v-list-item-title>{{ row.scheme_name }}</v-list-item-title>
              <v-list-item-subtitle v-if="row.message">{{
                row.message
              }}</v-list-item-subtitle>
              <template #append>
                <v-chip
                  :color="statusResultColor(row.status)"
                  size="x-small"
                  variant="tonal"
                >
                  {{ row.status }}
                </v-chip>
              </template>
            </v-list-item>
          </v-list>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="primary" @click="resultsDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Coverage Overview Dialog -->
    <v-dialog v-model="coverageDialog" max-width="1200" scrollable>
      <v-card>
        <v-card-title class="d-flex align-center">
          <v-icon class="mr-2">mdi-grid</v-icon>
          Schedule Coverage Overview (Last 12 Months)
        </v-card-title>
        <v-card-text>
          <v-alert
            v-if="!coverageLoading && coverageMatrix && missingCount > 0"
            type="warning"
            variant="tonal"
            class="mb-4"
            density="compact"
          >
            {{ missingCount }} scheme/month combination(s) without a premium
            schedule detected.
          </v-alert>
          <v-alert
            v-else-if="!coverageLoading && coverageMatrix && missingCount === 0"
            type="success"
            variant="tonal"
            class="mb-4"
            density="compact"
          >
            All in-force schemes have premium schedules for every month in the
            window.
          </v-alert>

          <div v-if="coverageLoading" class="d-flex justify-center py-8">
            <v-progress-circular indeterminate color="primary" />
          </div>

          <div v-else-if="coverageMatrix" style="overflow-x: auto">
            <table class="coverage-table">
              <thead>
                <tr>
                  <th class="sticky-col scheme-col">Scheme</th>
                  <th class="sticky-col commencement-col">Commencement</th>
                  <th
                    v-for="m in coverageMatrix.months"
                    :key="m"
                    class="month-col"
                    >{{ m }}</th
                  >
                </tr>
              </thead>
              <tbody>
                <tr v-for="row in coverageMatrix.schemes" :key="row.scheme_id">
                  <td class="sticky-col scheme-col">{{ row.scheme_name }}</td>
                  <td class="sticky-col commencement-col">{{
                    row.commencement_date || '—'
                  }}</td>
                  <td
                    v-for="(cell, idx) in row.cells"
                    :key="idx"
                    class="month-col text-center"
                    :class="{ 'before-commencement': cell.before_commencement }"
                  >
                    <v-tooltip
                      :text="
                        cell.before_commencement
                          ? 'Before commencement date'
                          : cell.exists
                            ? `${cell.status} (#${cell.schedule_id})`
                            : 'No schedule'
                      "
                      location="top"
                    >
                      <template #activator="{ props }">
                        <v-icon
                          v-bind="props"
                          :color="
                            cell.before_commencement
                              ? 'grey-lighten-1'
                              : cell.exists
                                ? coverageCellColor(cell.status)
                                : 'error'
                          "
                          :icon="
                            cell.before_commencement
                              ? 'mdi-minus-circle-outline'
                              : cell.exists
                                ? 'mdi-check-circle'
                                : 'mdi-close-circle'
                          "
                          size="20"
                        />
                      </template>
                    </v-tooltip>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="primary" @click="coverageDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar
      v-model="snackbar"
      :color="snackbarColor"
      :timeout="3500"
      centered
    >
      {{ snackbarText }}
      <template #actions>
        <v-btn variant="text" @click="snackbar = false">Close</v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-balham.css'
import { AgGridVue } from 'ag-grid-vue3'
import PremiumManagementService from '@/renderer/api/PremiumManagementService'
import BaseCard from '@/renderer/components/BaseCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import { useGridHeight } from '@/renderer/composables/useGridHeight'
import { useFilterPersistence } from '@/renderer/composables/useFilterPersistence'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'
import { statusCellRenderer } from '@/renderer/utils/statusCellRenderer'
import { dateFormatter } from '@/renderer/utils/formatters'
import { useStatusBarStore } from '@/renderer/store/statusBar'

const statusBarStore = useStatusBarStore()
const { hasPermission } = usePermissionCheck()

const router = useRouter()
const gridHeight = useGridHeight(380)
const { filters, resetFilters } = useFilterPersistence('premium-schedules', {
  schemeId: null as number | null,
  month: null as number | null,
  year: null as number | null,
  status: null as string | null
})

const loading = ref(false)
const generating = ref(false)
const generateDialog = ref(false)
const resultsDialog = ref(false)
const schedules = ref<any[]>([])
const inForceSchemes = ref<{ id: number; name: string }[]>([])
const schemesLoading = ref(false)
const snackbar = ref(false)
const snackbarText = ref('')
const snackbarColor = ref('success')
const bulkResult = ref<any>(null)
const coverageDialog = ref(false)
const coverageLoading = ref(false)
const coverageMatrix = ref<any>(null)

const missingCount = computed(() => {
  if (!coverageMatrix.value?.schemes) return 0
  let count = 0
  for (const row of coverageMatrix.value.schemes) {
    for (const cell of row.cells) {
      if (!cell.exists && !cell.before_commencement) count++
    }
  }
  return count
})

// filters and resetFilters are provided by useFilterPersistence above
const genForm = ref({
  allSchemes: false,
  schemeId: null as number | null,
  month: new Date().getMonth() + 1,
  year: new Date().getFullYear()
})

const months = [
  { label: 'January', value: 1 },
  { label: 'February', value: 2 },
  { label: 'March', value: 3 },
  { label: 'April', value: 4 },
  { label: 'May', value: 5 },
  { label: 'June', value: 6 },
  { label: 'July', value: 7 },
  { label: 'August', value: 8 },
  { label: 'September', value: 9 },
  { label: 'October', value: 10 },
  { label: 'November', value: 11 },
  { label: 'December', value: 12 }
]
const statusOptions = [
  'draft',
  'reviewed',
  'approved',
  'finalized',
  'invoiced',
  'void',
  'cancelled'
]

const defaultColDef = { sortable: true, filter: true, resizable: true, flex: 1 }
const columnDefs = [
  { headerName: 'ID', field: 'id', maxWidth: 80 },
  { headerName: 'Scheme', field: 'scheme_name', minWidth: 180 },
  {
    headerName: 'Period',
    valueGetter: (p: any) => `${p.data.month}/${p.data.year}`
  },
  { headerName: 'Members', field: 'member_count', maxWidth: 110 },
  {
    headerName: 'Gross Premium',
    field: 'gross_premium',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: 'Net Payable',
    field: 'net_payable',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: 'Status',
    field: 'status',
    cellRenderer: (p: any) => statusCellRenderer(p.value)
  },
  {
    headerName: 'Generated',
    field: 'generated_date',
    valueFormatter: dateFormatter
  },
  { headerName: 'By', field: 'generated_by' }
]

function fmtCurrency(val: number) {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR',
    maximumFractionDigits: 0
  }).format(val)
}

function statusIcon(status: string) {
  if (status === 'success') return 'mdi-check-circle'
  if (status === 'skipped') return 'mdi-skip-next-circle'
  return 'mdi-alert-circle'
}

function statusResultColor(status: string) {
  if (status === 'success') return 'success'
  if (status === 'skipped') return 'warning'
  return 'error'
}

function onRowClicked(e: any) {
  router.push({
    name: 'group-pricing-premium-schedule-detail',
    params: { scheduleId: e.data.id }
  })
}

function coverageCellColor(status: string) {
  const colors: Record<string, string> = {
    draft: 'grey',
    reviewed: 'blue',
    approved: 'light-green',
    finalized: 'green',
    invoiced: 'teal'
  }
  return colors[status] ?? 'green'
}

async function openCoverageOverview() {
  coverageDialog.value = true
  coverageLoading.value = true
  try {
    const res = await PremiumManagementService.getScheduleCoverageMatrix()
    coverageMatrix.value = res.data.data ?? null
  } catch (e) {
    console.error('Failed to load coverage matrix', e)
    snackbarText.value = 'Failed to load coverage overview'
    snackbarColor.value = 'error'
    snackbar.value = true
  } finally {
    coverageLoading.value = false
  }
}

function openGenerateDialog() {
  genForm.value.allSchemes = false
  genForm.value.schemeId = null
  genForm.value.month = new Date().getMonth() + 1
  genForm.value.year = new Date().getFullYear()
  generateDialog.value = true
}

async function loadSchedules() {
  loading.value = true
  try {
    const params: any = {}
    if (filters.value.schemeId) params.scheme_id = filters.value.schemeId
    if (filters.value.month) params.month = filters.value.month
    if (filters.value.year) params.year = filters.value.year
    if (filters.value.status) params.status = filters.value.status
    const res = await PremiumManagementService.getPremiumSchedules(params)
    schedules.value = res.data.data ?? []
    const finalized = schedules.value.filter(
      (s) => s.status === 'finalized'
    ).length
    const draft = schedules.value.filter((s) => s.status === 'draft').length
    const invoiced = schedules.value.filter(
      (s) => s.status === 'invoiced'
    ).length
    statusBarStore.set([
      { icon: 'mdi-calendar-check', text: `Finalized: ${finalized}` },
      {
        icon: 'mdi-file-outline',
        text: `Draft: ${draft}`,
        severity: draft > 0 ? 'warn' : 'info'
      },
      { icon: 'mdi-receipt-text-outline', text: `Invoiced: ${invoiced}` }
    ])
  } catch (e) {
    console.error('Failed to load schedules', e)
  } finally {
    loading.value = false
  }
}

async function handleGenerate() {
  generating.value = true
  try {
    if (genForm.value.allSchemes) {
      const res = await PremiumManagementService.generateAllSchedules({
        month: genForm.value.month,
        year: genForm.value.year
      })
      generateDialog.value = false
      bulkResult.value = res.data.data
      resultsDialog.value = true
      await loadSchedules()
    } else {
      if (!genForm.value.schemeId) return
      await PremiumManagementService.generateSchedule({
        scheme_id: genForm.value.schemeId!,
        month: genForm.value.month,
        year: genForm.value.year
      })
      generateDialog.value = false
      snackbarText.value = 'Schedule generated successfully'
      snackbarColor.value = 'success'
      snackbar.value = true
      await loadSchedules()
    }
  } catch (e: any) {
    snackbarText.value =
      e?.response?.data?.message ?? 'Failed to generate schedule'
    snackbarColor.value = 'error'
    snackbar.value = true
  } finally {
    generating.value = false
  }
}

async function loadInForceSchemes() {
  schemesLoading.value = true
  try {
    const res = await PremiumManagementService.getInForceSchemes()
    inForceSchemes.value = (res.data ?? []).map((s: any) => ({
      id: s.id,
      name: s.name
    }))
  } catch (e) {
    console.error('Failed to load in-force schemes', e)
  } finally {
    schemesLoading.value = false
  }
}

onMounted(() => {
  loadInForceSchemes()
  loadSchedules()
})
onUnmounted(() => statusBarStore.clear())
</script>

<style scoped>
.coverage-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}
.coverage-table th,
.coverage-table td {
  border: 1px solid #e0e0e0;
  padding: 6px 10px;
  white-space: nowrap;
}
.coverage-table thead th {
  background: #f5f5f5;
  font-weight: 600;
  position: sticky;
  top: 0;
  z-index: 2;
}
.coverage-table .sticky-col {
  position: sticky;
  background: #fff;
  z-index: 1;
}
.coverage-table .scheme-col {
  left: 0;
  min-width: 200px;
  max-width: 280px;
  overflow: hidden;
  text-overflow: ellipsis;
}
.coverage-table .commencement-col {
  left: 200px;
  min-width: 120px;
}
.coverage-table thead .sticky-col {
  z-index: 3;
}
.coverage-table .month-col {
  min-width: 80px;
  text-align: center;
}
.coverage-table tbody tr:hover {
  background: #f9f9f9;
}
.coverage-table td.before-commencement {
  background: #f0f0f0;
  opacity: 0.5;
}
</style>
