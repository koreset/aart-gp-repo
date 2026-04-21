<template>
  <v-container fluid>
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
          <span class="headline">Submission Deadline Calendar</span>
        </div>
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
          <v-col cols="12" md="12" class="d-flex gap-2 justify-end">
            <v-btn
              variant="outlined"
              rounded
              size="small"
              class="mr-3"
              color="primary"
              prepend-icon="mdi-refresh"
              :loading="loading"
              @click="loadDeadlines"
            >
              Refresh
            </v-btn>
            <v-btn
              rounded
              size="small"
              class="mr-3"
              color="primary"
              prepend-icon="mdi-calendar-plus"
              @click="generateDialog = true"
            >
              Generate Deadlines
            </v-btn>
            <v-btn
              rounded
              size="small"
              class="mr-3"
              variant="outlined"
              color="primary"
              prepend-icon="mdi-plus"
              @click="createDialog = true"
            >
              Add Single
            </v-btn>
          </v-col>
        </v-row>

        <!-- Summary stat chips -->
        <v-row class="mb-3">
          <v-col cols="auto">
            <v-chip
              color="error"
              variant="tonal"
              size="small"
              prepend-icon="mdi-alert-circle"
            >
              Overdue: {{ stats.overdue_count ?? 0 }}
            </v-chip>
          </v-col>
          <v-col cols="auto">
            <v-chip
              color="warning"
              variant="tonal"
              size="small"
              prepend-icon="mdi-clock-outline"
            >
              Pending: {{ stats.pending_count ?? 0 }}
            </v-chip>
          </v-col>
          <v-col cols="auto">
            <v-chip
              color="success"
              variant="tonal"
              size="small"
              prepend-icon="mdi-check-circle"
            >
              Received: {{ stats.received_count ?? 0 }}
            </v-chip>
          </v-col>
          <v-col cols="auto">
            <v-chip
              color="grey"
              variant="tonal"
              size="small"
              prepend-icon="mdi-minus-circle"
            >
              Waived: {{ stats.waived_count ?? 0 }}
            </v-chip>
          </v-col>
        </v-row>

        <!-- Deadlines Grid -->
        <v-row>
          <v-col cols="12">
            <v-card variant="outlined">
              <v-card-text class="pa-0">
                <ag-grid-vue
                  class="ag-theme-balham"
                  style="height: 540px; width: 100%"
                  :column-defs="columnDefs"
                  :row-data="deadlines"
                  :default-col-def="defaultColDef"
                  :loading="loading"
                  :row-class-rules="rowClassRules"
                  @cell-clicked="onCellClicked"
                />
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </template>
    </base-card>

    <!-- Generate Deadlines Dialog -->
    <v-dialog v-model="generateDialog" max-width="480" persistent>
      <v-card>
        <v-card-title>Generate Submission Deadlines</v-card-title>
        <v-card-text>
          <v-alert type="info" variant="tonal" class="mb-4" density="compact">
            Creates a deadline for every in-force scheme for the selected
            period. Existing deadlines are skipped.
          </v-alert>
          <v-row>
            <v-col cols="6">
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
            <v-col cols="6">
              <v-text-field
                v-model.number="genForm.year"
                label="Year *"
                type="number"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model.number="genForm.dueDayOfMonth"
                label="Due Day of Month"
                type="number"
                variant="outlined"
                density="compact"
                hint="1–28, default 15"
                persistent-hint
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model.number="genForm.gracePeriodDays"
                label="Grace Period (days)"
                type="number"
                variant="outlined"
                density="compact"
                hint="0 = no grace period"
                persistent-hint
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
            :disabled="!genForm.month || !genForm.year"
            @click="handleGenerate"
          >
            Generate
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Create Single Deadline Dialog -->
    <v-dialog v-model="createDialog" max-width="480" persistent>
      <v-card>
        <v-card-title>Add Deadline</v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12">
              <v-select
                v-model="createForm.schemeId"
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
                v-model.number="createForm.month"
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
                v-model.number="createForm.year"
                label="Year *"
                type="number"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="8">
              <v-text-field
                v-model="createForm.dueDate"
                label="Due Date *"
                type="date"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="4">
              <v-text-field
                v-model.number="createForm.gracePeriodDays"
                label="Grace Days"
                type="number"
                variant="outlined"
                density="compact"
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
              !createForm.schemeId ||
              !createForm.month ||
              !createForm.year ||
              !createForm.dueDate
            "
            @click="handleCreate"
          >
            Create
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Waive Deadline Dialog -->
    <v-dialog v-model="waiveDialog.show" max-width="420" persistent>
      <v-card>
        <v-card-title>Waive Deadline</v-card-title>
        <v-card-text>
          <v-textarea
            v-model="waiveDialog.reason"
            label="Waiver Reason"
            variant="outlined"
            density="compact"
            rows="3"
            auto-grow
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="waiveDialog.show = false"
            >Cancel</v-btn
          >
          <v-btn
            color="warning"
            :loading="waiveDialog.loading"
            @click="handleWaive"
            >Waive</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="4000">
      {{ snackbar.message }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-balham.css'
import { AgGridVue } from 'ag-grid-vue3'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import PremiumManagementService from '@/renderer/api/PremiumManagementService'

const router = useRouter()

const loading = ref(false)
const generating = ref(false)
const creating = ref(false)
const schemesLoading = ref(false)
const deadlines = ref<any[]>([])
const inForceSchemes = ref<any[]>([])
const stats = ref<any>({})
const generateDialog = ref(false)
const createDialog = ref(false)

const filters = ref({
  schemeId: null as number | null,
  month: null as number | null,
  year: null as number | null,
  status: ''
})
const genForm = ref({
  month: new Date().getMonth() + 1,
  year: new Date().getFullYear(),
  dueDayOfMonth: 15,
  gracePeriodDays: 5
})
const createForm = ref({
  schemeId: null as number | null,
  month: null as number | null,
  year: new Date().getFullYear(),
  dueDate: '',
  gracePeriodDays: 5
})
const waiveDialog = ref({ show: false, loading: false, id: 0, reason: '' })
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
const statusOptions = ['pending', 'received', 'overdue', 'waived']

const statusColors: Record<string, string> = {
  pending: '#FF9800',
  received: '#4CAF50',
  overdue: '#F44336',
  waived: '#9E9E9E'
}
const statusBadge = (value: string) => {
  const c = statusColors[value] ?? '#9E9E9E'
  return `<span style="background:${c}22;color:${c};padding:2px 8px;border-radius:4px;font-size:12px;font-weight:500">${value}</span>`
}

const graceExpiry = (dueDate: string, graceDays: number) => {
  if (!dueDate || !graceDays) return dueDate
  const d = new Date(dueDate)
  d.setDate(d.getDate() + graceDays)
  return d.toISOString().slice(0, 10)
}

const defaultColDef = { sortable: true, filter: true, resizable: true }

const rowClassRules = {
  'row-overdue': (p: any) => p.data?.status === 'overdue'
}

const columnDefs = [
  { headerName: 'ID', field: 'id', width: 70 },
  { headerName: 'Scheme', field: 'scheme_name', flex: 1, minWidth: 140 },
  {
    headerName: 'Period',
    flex: 1,
    minWidth: 90,
    valueGetter: (p: any) => (p.data ? `${p.data.month}/${p.data.year}` : '')
  },
  {
    headerName: 'Type',
    field: 'type',
    width: 150,
    valueFormatter: (p: any) => (p.value ?? '').replace(/_/g, ' ')
  },
  { headerName: 'Due Date', field: 'due_date', width: 120 },
  {
    headerName: 'Grace Expiry',
    width: 120,
    valueGetter: (p: any) =>
      p.data ? graceExpiry(p.data.due_date, p.data.grace_period_days) : ''
  },
  {
    headerName: 'Status',
    field: 'status',
    width: 120,
    cellRenderer: (p: any) => statusBadge(p.value)
  },
  {
    headerName: 'Submission',
    field: 'linked_submission_id',
    width: 120,
    cellRenderer: (p: any) => {
      if (!p.value)
        return '<span class="text-medium-emphasis text-caption">—</span>'
      return `<span style="color:#1976d2;cursor:pointer" data-sub-id="${p.value}">#${p.value}</span>`
    }
  },
  {
    headerName: 'Actions',
    width: 160,
    sortable: false,
    filter: false,
    cellRenderer: (p: any) => {
      if (!p.data) return ''
      if (p.data.status === 'waived') {
        return `<button class="action-btn" data-action="reopen" data-id="${p.data.id}" style="font-size:12px;padding:2px 8px;border:1px solid #9E9E9E;border-radius:4px;cursor:pointer">Reopen</button>`
      }
      if (p.data.status === 'pending' || p.data.status === 'overdue') {
        return `<button class="action-btn" data-action="waive" data-id="${p.data.id}" style="font-size:12px;padding:2px 8px;border:1px solid #FF9800;border-radius:4px;cursor:pointer;color:#FF9800">Waive</button>`
      }
      return ''
    }
  }
]

const loadDeadlines = async () => {
  loading.value = true
  try {
    const params: any = {}
    if (filters.value.schemeId) params.scheme_id = filters.value.schemeId
    if (filters.value.month) params.month = filters.value.month
    if (filters.value.year) params.year = filters.value.year
    if (filters.value.status) params.status = filters.value.status
    const res = await GroupPricingService.getBordereauxDeadlines(params)
    deadlines.value = res.data?.data ?? []
  } catch {
    snackbar.value = {
      show: true,
      message: 'Failed to load deadlines',
      color: 'error'
    }
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const res = await GroupPricingService.getDeadlineStats()
    stats.value = res.data?.data ?? {}
  } catch {
    // non-fatal
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

const handleGenerate = async () => {
  generating.value = true
  try {
    const res = await GroupPricingService.generateBordereauxDeadlines({
      month: genForm.value.month,
      year: genForm.value.year,
      due_day_of_month: genForm.value.dueDayOfMonth,
      grace_period_days: genForm.value.gracePeriodDays
    })
    const r = res.data?.data
    generateDialog.value = false
    snackbar.value = {
      show: true,
      message: `Generated ${r?.created ?? 0} deadlines, skipped ${r?.skipped ?? 0}`,
      color: 'success'
    }
    await Promise.all([loadDeadlines(), loadStats()])
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Generation failed',
      color: 'error'
    }
  } finally {
    generating.value = false
  }
}

const handleCreate = async () => {
  creating.value = true
  try {
    await GroupPricingService.createBordereauxDeadline({
      scheme_id: createForm.value.schemeId,
      month: createForm.value.month,
      year: createForm.value.year,
      due_date: createForm.value.dueDate,
      grace_period_days: createForm.value.gracePeriodDays
    })
    createDialog.value = false
    snackbar.value = {
      show: true,
      message: 'Deadline created',
      color: 'success'
    }
    await Promise.all([loadDeadlines(), loadStats()])
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Create failed',
      color: 'error'
    }
  } finally {
    creating.value = false
  }
}

const handleWaive = async () => {
  waiveDialog.value.loading = true
  try {
    await GroupPricingService.updateDeadlineStatus(waiveDialog.value.id, {
      status: 'waived',
      waiver_reason: waiveDialog.value.reason
    })
    waiveDialog.value.show = false
    waiveDialog.value.reason = ''
    snackbar.value = {
      show: true,
      message: 'Deadline waived',
      color: 'warning'
    }
    await Promise.all([loadDeadlines(), loadStats()])
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Action failed',
      color: 'error'
    }
  } finally {
    waiveDialog.value.loading = false
  }
}

const onCellClicked = (e: any) => {
  const target = e.event?.target as HTMLElement | null
  if (!target) return
  const btn = target.closest('.action-btn') as HTMLElement | null
  if (btn) {
    const action = btn.dataset.action
    const id = parseInt(btn.dataset.id ?? '0')
    if (action === 'waive') {
      waiveDialog.value = { show: true, loading: false, id, reason: '' }
    } else if (action === 'reopen') {
      GroupPricingService.updateDeadlineStatus(id, {
        status: 'pending',
        waiver_reason: ''
      })
        .then(() => {
          snackbar.value = {
            show: true,
            message: 'Deadline reopened',
            color: 'info'
          }
          return Promise.all([loadDeadlines(), loadStats()])
        })
        .catch(() => {
          snackbar.value = {
            show: true,
            message: 'Failed to reopen deadline',
            color: 'error'
          }
        })
    }
    return
  }
  // Submission link click
  const subLink = target.closest('[data-sub-id]') as HTMLElement | null
  if (subLink) {
    const subId = subLink.dataset.subId
    if (subId) {
      router.push({
        name: 'group-pricing-bordereaux-inbound-detail',
        params: { submissionId: subId }
      })
    }
  }
}

onMounted(() => {
  loadInForceSchemes()
  loadDeadlines()
  loadStats()
})
</script>

<style scoped>
:deep(.row-overdue) {
  background-color: #fff5f5 !important;
}
</style>
