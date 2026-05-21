<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex justify-space-between align-center flex-wrap">
              <span class="headline">CPI Index</span>
              <div class="d-flex align-center gap-2">
                <v-btn
                  size="small"
                  variant="outlined"
                  class="mr-2"
                  rounded
                  prepend-icon="mdi-download"
                  @click="downloadTemplate"
                >
                  Download Template
                </v-btn>
                <v-btn
                  v-if="hasPermission('claims:manage_cpi_indices')"
                  size="small"
                  variant="outlined"
                  class="mr-2"
                  rounded
                  prepend-icon="mdi-upload"
                  @click="uploadDialog = true"
                >
                  Upload
                </v-btn>
                <v-btn
                  v-if="hasPermission('claims:manage_cpi_indices')"
                  size="small"
                  variant="outlined"
                  rounded
                  prepend-icon="mdi-plus"
                  @click="openAddDialog"
                >
                  Add Index
                </v-btn>
              </div>
            </div>
          </template>
          <template #default>
            <v-row class="mb-4">
              <v-col cols="12" md="3">
                <v-text-field
                  v-model="searchQuery"
                  label="Search"
                  prepend-inner-icon="mdi-magnify"
                  variant="outlined"
                  density="compact"
                  clearable
                />
              </v-col>
              <v-col cols="12" md="3">
                <v-select
                  v-model="selectedYear"
                  :items="yearOptions"
                  label="Year"
                  variant="outlined"
                  density="compact"
                  clearable
                />
              </v-col>
            </v-row>

            <div :style="{ height: gridHeight, width: '100%' }">
              <data-grid
                :column-defs="columnDefs"
                :row-data="filteredRows"
                :loading="loading"
                style="height: 100%; width: 100%"
              />
            </div>
            <empty-state
              v-if="!loading && filteredRows.length === 0"
              icon="mdi-table-off"
              title="No CPI rows"
              message="Upload a CPI file or add a row to get started."
            />
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Add / edit single row -->
    <v-dialog v-model="formDialog" persistent max-width="420px">
      <v-card>
        <v-card-title class="text-h6">Add CPI Index</v-card-title>
        <v-card-text>
          <v-row dense>
            <v-col cols="6">
              <v-text-field
                v-model.number="form.year_index"
                label="Year"
                type="number"
                density="compact"
                variant="outlined"
                :error-messages="formErrors.year_index"
              />
            </v-col>
            <v-col cols="6">
              <v-select
                v-model.number="form.month_index"
                :items="monthOptions"
                label="Month"
                item-title="text"
                item-value="value"
                density="compact"
                variant="outlined"
                :error-messages="formErrors.month_index"
              />
            </v-col>
            <v-col cols="12">
              <v-text-field
                v-model.number="form.cpi_index"
                label="CPI Index"
                type="number"
                step="0.001"
                density="compact"
                variant="outlined"
                :error-messages="formErrors.cpi_index"
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="formDialog = false">Cancel</v-btn>
          <v-btn color="primary" :loading="saving" @click="saveSingle">
            Save
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Upload file -->
    <v-dialog v-model="uploadDialog" persistent max-width="520px">
      <v-card>
        <v-card-title class="text-h6">Upload CPI File</v-card-title>
        <v-card-text>
          <p class="text-body-2 mb-3">
            Accepts .xlsx or .csv. Required columns:
            <strong>year_index</strong>, <strong>month_index</strong>,
            <strong>cpi_index</strong>. Re-uploading overwrites existing rows
            for the same (year, month).
            <a
              href="#"
              class="text-primary"
              @click.prevent="downloadTemplate"
            >Download template</a>.
          </p>
          <v-file-input
            v-model="uploadFile"
            accept=".csv,.xlsx"
            label="Choose file"
            variant="outlined"
            density="compact"
            show-size
            :error-messages="uploadError"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="closeUploadDialog">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="uploading"
            :disabled="!uploadFile"
            @click="submitUpload"
          >
            Upload
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="3000">
      {{ snackbarMessage }}
      <template #actions>
        <v-btn color="white" variant="text" @click="snackbar = false"
          >Close</v-btn
        >
      </template>
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { dateFormatter } from '@/renderer/utils/formatters'
import { useGridHeight } from '@/renderer/composables/useGridHeight'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

interface CpiRow {
  id: number
  year_index: number
  month_index: number
  cpi_index: number
  created: string
  created_by: string
}

const gridHeight = useGridHeight(360)
const { hasPermission } = usePermissionCheck()

const loading = ref(false)
const saving = ref(false)
const uploading = ref(false)
const rows = ref<CpiRow[]>([])

const searchQuery = ref('')
const selectedYear = ref<number | null>(null)

const formDialog = ref(false)
const form = reactive({
  year_index: new Date().getFullYear(),
  month_index: new Date().getMonth() + 1,
  cpi_index: 0
})
const formErrors = reactive<Record<string, string>>({})

const uploadDialog = ref(false)
const uploadFile = ref<File | File[] | null>(null)
const uploadError = ref<string>('')

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

const monthOptions = [
  { text: 'January', value: 1 },
  { text: 'February', value: 2 },
  { text: 'March', value: 3 },
  { text: 'April', value: 4 },
  { text: 'May', value: 5 },
  { text: 'June', value: 6 },
  { text: 'July', value: 7 },
  { text: 'August', value: 8 },
  { text: 'September', value: 9 },
  { text: 'October', value: 10 },
  { text: 'November', value: 11 },
  { text: 'December', value: 12 }
]

const columnDefs = [
  { headerName: 'ID', field: 'id', sortable: true, minWidth: 90, maxWidth: 110 },
  {
    headerName: 'Year',
    field: 'year_index',
    sortable: true,
    minWidth: 100,
    sort: 'desc' as const
  },
  {
    headerName: 'Month',
    field: 'month_index',
    sortable: true,
    minWidth: 100,
    valueFormatter: (p: any) => monthLabel(p.value)
  },
  {
    headerName: 'CPI Index',
    field: 'cpi_index',
    sortable: true,
    minWidth: 130,
    valueFormatter: (p: any) =>
      typeof p.value === 'number' ? p.value.toFixed(3) : '—'
  },
  {
    headerName: 'Created',
    field: 'created',
    sortable: true,
    minWidth: 160,
    valueFormatter: dateFormatter
  },
  {
    headerName: 'Created By',
    field: 'created_by',
    sortable: true,
    filter: true,
    minWidth: 160
  }
]

const yearOptions = computed(() =>
  Array.from(new Set(rows.value.map((r) => r.year_index))).sort((a, b) => b - a)
)

const filteredRows = computed(() => {
  let out = [...rows.value]
  if (selectedYear.value !== null && selectedYear.value !== undefined) {
    out = out.filter((r) => r.year_index === selectedYear.value)
  }
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    out = out.filter(
      (r) =>
        String(r.year_index).includes(q) ||
        String(r.month_index).includes(q) ||
        String(r.cpi_index).includes(q) ||
        monthLabel(r.month_index).toLowerCase().includes(q)
    )
  }
  return out
})

function monthLabel(m: number): string {
  const opt = monthOptions.find((o) => o.value === m)
  return opt ? opt.text : String(m)
}

function showSnackbar(msg: string, color: string = 'success') {
  snackbarMessage.value = msg
  snackbarColor.value = color
  snackbar.value = true
}

function downloadTemplate() {
  const headers = ['year_index', 'month_index', 'cpi_index']
  const sampleYear = new Date().getFullYear()
  const sampleRows = [
    [sampleYear, 1, 100.0],
    [sampleYear, 2, 100.5],
    [sampleYear, 3, 101.1]
  ]
  const csv =
    headers.join(',') +
    '\n' +
    sampleRows.map((r) => r.join(',')).join('\n') +
    '\n'
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = 'cpi_index_template.csv'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

async function load() {
  loading.value = true
  try {
    const res = await GroupPricingService.getCpiIndices()
    rows.value = res.data || []
  } catch (e) {
    console.error('Error loading CPI indices:', e)
    showSnackbar('Error loading CPI indices', 'error')
    rows.value = []
  } finally {
    loading.value = false
  }
}

function openAddDialog() {
  form.year_index = new Date().getFullYear()
  form.month_index = new Date().getMonth() + 1
  form.cpi_index = 0
  for (const k of Object.keys(formErrors)) delete formErrors[k]
  formDialog.value = true
}

function validateForm(): boolean {
  for (const k of Object.keys(formErrors)) delete formErrors[k]
  let ok = true
  if (!form.year_index || form.year_index < 1900 || form.year_index > 2100) {
    formErrors.year_index = 'Year must be between 1900 and 2100'
    ok = false
  }
  if (!form.month_index || form.month_index < 1 || form.month_index > 12) {
    formErrors.month_index = 'Month is required'
    ok = false
  }
  if (!form.cpi_index || form.cpi_index <= 0) {
    formErrors.cpi_index = 'CPI Index must be positive'
    ok = false
  }
  return ok
}

async function saveSingle() {
  if (!validateForm()) return
  saving.value = true
  try {
    await GroupPricingService.upsertCpiIndex({
      year_index: form.year_index,
      month_index: form.month_index,
      cpi_index: form.cpi_index
    })
    showSnackbar('CPI Index saved')
    formDialog.value = false
    await load()
  } catch (e: any) {
    console.error('Error saving CPI index:', e)
    showSnackbar(e?.response?.data?.error || 'Error saving CPI Index', 'error')
  } finally {
    saving.value = false
  }
}

function closeUploadDialog() {
  uploadDialog.value = false
  uploadFile.value = null
  uploadError.value = ''
}

async function submitUpload() {
  uploadError.value = ''
  const f = Array.isArray(uploadFile.value)
    ? uploadFile.value[0]
    : uploadFile.value
  if (!f) {
    uploadError.value = 'Please choose a file'
    return
  }
  uploading.value = true
  try {
    const res = await GroupPricingService.uploadCpiIndices(f)
    const n = res?.data?.upserted ?? 0
    showSnackbar(`Uploaded ${n} CPI rows`)
    closeUploadDialog()
    await load()
  } catch (e: any) {
    console.error('Error uploading CPI file:', e)
    uploadError.value =
      e?.response?.data?.error || 'Upload failed. Check the file and retry.'
  } finally {
    uploading.value = false
  }
}

onMounted(load)
</script>
