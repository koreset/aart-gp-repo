<template>
  <v-container fluid>
    <base-card :show-actions="false">
      <template #header>
        <div class="d-flex align-center justify-space-between">
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
            <span class="headline">Inbound Employer Submissions</span>
          </div>
          <div class="d-flex align-center gap-2">
            <v-btn
              class="mr-2"
              variant="outlined"
              size="small"
              rounded
              prepend-icon="mdi-refresh"
              :loading="loading"
              @click="loadSubmissions"
            >
              Refresh
            </v-btn>
            <v-btn
              size="small"
              variant="outlined"
              rounded
              prepend-icon="mdi-plus"
              @click="createDialog = true"
            >
              New Submission
            </v-btn>
          </div>
        </div>
      </template>
      <template #default>
        <!-- Filter Bar -->
        <v-row class="mb-3" align="center">
          <v-col cols="12" md="4">
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
          <v-col cols="12" md="2">
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

        <!-- Grid -->
        <v-row>
          <v-col cols="12">
            <v-card variant="outlined">
              <v-card-text class="pa-0">
                <ag-grid-vue
                  class="ag-theme-balham"
                  :style="{ height: gridHeight, width: '100%' }"
                  :column-defs="columnDefs"
                  :row-data="submissions"
                  :default-col-def="defaultColDef"
                  :loading="loading"
                  row-selection="single"
                  @row-clicked="onRowClicked"
                />
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </template>
    </base-card>

    <!-- New Submission Dialog -->
    <v-dialog v-model="createDialog" max-width="540" persistent>
      <v-card>
        <v-card-title>New Employer Submission</v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12">
              <v-select
                v-model="form.schemeId"
                label="Scheme *"
                :items="inForceSchemes"
                item-title="name"
                item-value="id"
                variant="outlined"
                density="compact"
                :loading="schemesLoading"
              />
            </v-col>
            <v-col cols="6">
              <v-select
                v-model.number="form.month"
                label="Month *"
                :items="months"
                item-title="label"
                item-value="value"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model.number="form.year"
                label="Year *"
                type="number"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12">
              <v-text-field
                v-model="form.dueDate"
                label="Due Date"
                type="date"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12">
              <v-text-field
                v-model="form.submittedBy"
                label="Submitted By"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="form.notes"
                label="Notes"
                variant="outlined"
                density="compact"
                rows="2"
                auto-grow
              />
            </v-col>
            <v-col cols="12">
              <v-checkbox
                v-model="form.isRetro"
                label="Retrospective Submission"
                density="compact"
                hide-details
              />
            </v-col>
            <v-col v-if="form.isRetro" cols="12">
              <v-text-field
                v-model="form.retroEffectiveDate"
                label="Retro Effective Date *"
                type="date"
                variant="outlined"
                density="compact"
                hint="The earliest date from which catch-up premiums are due"
                persistent-hint
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="createDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="creating"
            :disabled="
              !form.schemeId ||
              !form.month ||
              !form.year ||
              (form.isRetro && !form.retroEffectiveDate)
            "
            @click="handleCreate"
          >
            Create
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="4000">
      {{ snackbar.message }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-balham.css'
import { AgGridVue } from 'ag-grid-vue3'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import PremiumManagementService from '@/renderer/api/PremiumManagementService'
import { statusCellRenderer } from '@/renderer/utils/statusCellRenderer'
import { dateFormatter } from '@/renderer/utils/formatters'
import { useFilterPersistence } from '@/renderer/composables/useFilterPersistence'
import { useGridHeight } from '@/renderer/composables/useGridHeight'
import { useStatusBarStore } from '@/renderer/store/statusBar'

const statusBarStore = useStatusBarStore()

const gridHeight = useGridHeight(340)

const router = useRouter()

const loading = ref(false)
const creating = ref(false)
const schemesLoading = ref(false)
const createDialog = ref(false)
const submissions = ref<any[]>([])
const inForceSchemes = ref<any[]>([])

const { filters } = useFilterPersistence('bordereaux-inbound-submissions', {
  schemeId: null as number | null,
  month: null as number | null,
  year: null as number | null,
  status: ''
})
const form = ref({
  schemeId: null as number | null,
  month: null as number | null,
  year: new Date().getFullYear(),
  dueDate: '',
  submittedBy: '',
  notes: '',
  isRetro: false,
  retroEffectiveDate: ''
})
const snackbar = ref({ show: false, message: '', color: 'success' })

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
  'pending_receipt',
  'received',
  'under_review',
  'queries_raised',
  'accepted',
  'rejected'
]

const defaultColDef = {
  sortable: true,
  filter: true,
  resizable: true,
  suppressMovable: false
}

const columnDefs = [
  { headerName: 'ID', field: 'id', width: 70 },
  { headerName: 'Scheme', field: 'scheme_name', flex: 1, minWidth: 140 },
  {
    headerName: 'Period',
    flex: 1,
    minWidth: 100,
    valueGetter: (p: any) => (p.data ? `${p.data.month}/${p.data.year}` : '')
  },
  {
    headerName: 'Due Date',
    field: 'due_date',
    width: 120,
    valueFormatter: dateFormatter
  },
  {
    headerName: 'Received',
    field: 'received_date',
    width: 130,
    valueFormatter: dateFormatter
  },
  {
    headerName: 'Status',
    field: 'status',
    width: 150,
    cellRenderer: (p: any) => statusCellRenderer(p.value)
  },
  {
    headerName: 'Type',
    field: 'is_retro',
    width: 80,
    sortable: false,
    filter: false,
    cellRenderer: (p: any) =>
      p.value
        ? `<span style="background:#7B1FA222;color:#7B1FA2;padding:2px 6px;border-radius:4px;font-size:11px;font-weight:600">RETRO</span>`
        : ''
  },
  {
    headerName: 'Records',
    field: 'record_count',
    width: 95,
    type: 'numericColumn'
  },
  {
    headerName: 'Valid',
    field: 'valid_count',
    width: 80,
    type: 'numericColumn'
  },
  {
    headerName: 'Invalid',
    field: 'invalid_count',
    width: 85,
    type: 'numericColumn'
  },
  {
    headerName: 'Actions',
    width: 90,
    sortable: false,
    filter: false,
    cellRenderer: () => `<v-btn size="small" variant="text">View</v-btn>`
  }
]

const loadSubmissions = async () => {
  loading.value = true
  try {
    const params: any = {}
    if (filters.value.schemeId) params.scheme_id = filters.value.schemeId
    if (filters.value.month) params.month = filters.value.month
    if (filters.value.year) params.year = filters.value.year
    if (filters.value.status) params.status = filters.value.status
    const res = await GroupPricingService.getEmployerSubmissions(params)
    submissions.value = res.data?.data ?? []
    const pending = submissions.value.filter((s: any) =>
      ['pending_receipt', 'received', 'under_review'].includes(s.status)
    ).length
    statusBarStore.set([
      {
        icon: 'mdi-inbox-arrow-down',
        text: `Submissions: ${submissions.value.length} total`
      },
      {
        icon: 'mdi-clock-outline',
        text: `Pending: ${pending}`,
        severity: pending > 0 ? 'warn' : 'info'
      }
    ])
  } catch {
    snackbar.value = {
      show: true,
      message: 'Failed to load submissions',
      color: 'error'
    }
  } finally {
    loading.value = false
  }
}

const loadInForceSchemes = async () => {
  schemesLoading.value = true
  try {
    const res = await PremiumManagementService.getInForceSchemes()
    inForceSchemes.value = (res.data ?? []).map((s: any) => ({
      id: s.id,
      name: s.name
    }))
  } catch {
    // non-fatal
  } finally {
    schemesLoading.value = false
  }
}

const handleCreate = async () => {
  creating.value = true
  try {
    const res = await GroupPricingService.createEmployerSubmission({
      scheme_id: form.value.schemeId,
      month: form.value.month,
      year: form.value.year,
      due_date: form.value.dueDate,
      submitted_by: form.value.submittedBy,
      notes: form.value.notes,
      is_retro: form.value.isRetro,
      retro_effective_date: form.value.isRetro
        ? form.value.retroEffectiveDate
        : ''
    })
    createDialog.value = false
    const id = res.data?.data?.id
    if (id) {
      router.push({
        name: 'group-pricing-bordereaux-inbound-detail',
        params: { submissionId: id }
      })
    } else {
      await loadSubmissions()
    }
  } catch {
    snackbar.value = {
      show: true,
      message: 'Failed to create submission',
      color: 'error'
    }
  } finally {
    creating.value = false
  }
}

const onRowClicked = (e: any) => {
  if (e.data?.id) {
    router.push({
      name: 'group-pricing-bordereaux-inbound-detail',
      params: { submissionId: e.data.id }
    })
  }
}

onMounted(() => {
  loadInForceSchemes()
  loadSubmissions()
})
onUnmounted(() => statusBarStore.clear())
</script>
