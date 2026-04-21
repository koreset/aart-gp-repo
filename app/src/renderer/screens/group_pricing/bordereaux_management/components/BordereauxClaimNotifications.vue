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
          <span class="headline">Claim Notification Cadence</span>
        </div>
      </template>
      <template #default>
        <!-- Stat chips -->
        <v-row class="mb-3">
          <v-col cols="12">
            <div class="d-flex flex-wrap gap-3">
              <v-chip
                class="mr-2"
                color="primary"
                variant="tonal"
                size="small"
                prepend-icon="mdi-bell-ring"
              >
                Total: {{ stats.total }}
              </v-chip>
              <v-chip
                class="mr-2"
                color="warning"
                variant="tonal"
                size="small"
                prepend-icon="mdi-clock-outline"
              >
                Pending: {{ stats.pending }}
              </v-chip>
              <v-chip
                class="mr-2"
                color="blue"
                variant="tonal"
                size="small"
                prepend-icon="mdi-send"
              >
                Sent: {{ stats.sent }}
              </v-chip>
              <v-chip
                class="mr-2"
                color="success"
                variant="tonal"
                size="small"
                prepend-icon="mdi-check-circle"
              >
                Acknowledged: {{ stats.acknowledged }}
              </v-chip>
              <v-chip
                class="mr-2"
                color="error"
                variant="tonal"
                size="small"
                prepend-icon="mdi-alert-circle"
              >
                Overdue: {{ stats.overdue }}
              </v-chip>
            </div>
          </v-col>
        </v-row>

        <!-- Filter bar -->
        <v-row class="mb-3">
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
              :loading="schemesLoading"
            />
          </v-col>
          <v-col cols="12" md="3">
            <v-select
              v-model="filters.notificationType"
              label="Type"
              :items="notificationTypeOptions"
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
          <v-col cols="12" md="3">
            <v-text-field
              v-model="filters.claimNumber"
              label="Claim Number"
              variant="outlined"
              density="compact"
              clearable
            />
          </v-col>
          <v-col cols="12" md="12" class="d-flex gap-2 justify-end">
            <v-btn
              variant="outlined"
              rounded
              class="mr-3"
              size="small"
              color="primary"
              prepend-icon="mdi-refresh"
              :loading="loading"
              @click="loadNotifications"
            >
              Refresh
            </v-btn>
            <v-btn
              rounded
              size="small"
              class="mr-3"
              color="amber-darken-2"
              prepend-icon="mdi-calendar-refresh"
              @click="generateDialog = true"
            >
              Generate Month-End
            </v-btn>
            <v-btn
              rounded
              size="small"
              color="primary"
              class="mr-3"
              prepend-icon="mdi-plus"
              @click="addDialog = true"
            >
              Add Notification
            </v-btn>
            <v-btn
              variant="outlined"
              rounded
              size="small"
              color="teal"
              prepend-icon="mdi-file-download-outline"
              :loading="exporting"
              @click="handleExport"
            >
              Export CSV
            </v-btn>
          </v-col>
        </v-row>

        <!-- Grid -->
        <v-row>
          <v-col cols="12">
            <v-card variant="outlined">
              <v-card-text class="pa-0">
                <ag-grid-vue
                  class="ag-theme-balham"
                  style="height: 520px; width: 100%"
                  :column-defs="columnDefs"
                  :row-data="notifications"
                  :default-col-def="defaultColDef"
                  :loading="loading"
                  :overlay-loading-template="loadingOverlay"
                  :overlay-no-rows-template="`<span class='ag-overlay-no-rows-center'>No claim notifications match these filters.</span>`"
                  @grid-ready="onGridReady"
                />
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </template>
    </base-card>

    <!-- Add Notification Dialog -->
    <v-dialog v-model="addDialog" max-width="560" persistent>
      <v-card>
        <v-card-title>Add Claim Notification</v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12">
              <v-select
                v-model="form.schemeId"
                label="Scheme"
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
                v-model="form.claimNumber"
                label="Claim Number *"
                :items="schemeClaims"
                :item-title="(c: any) => `${c.claim_number} — ${c.member_name}`"
                item-value="claim_number"
                variant="outlined"
                density="compact"
                :loading="schemeClaimsLoading"
                :disabled="!form.schemeId"
                :hint="!form.schemeId ? 'Select a scheme first' : ''"
                persistent-hint
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-select
                v-model="form.notificationType"
                label="Notification Type *"
                :items="notificationTypeOptions"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="form.reinsurerName"
                label="Reinsurer Name"
                variant="outlined"
                density="compact"
                readonly
                :hint="form.reinsurerName ? '' : 'Auto-filled from treaty'"
                persistent-hint
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="form.reinsurerCode"
                label="Reinsurer Code"
                variant="outlined"
                density="compact"
                readonly
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
              <v-textarea
                v-model="form.notes"
                label="Notes"
                variant="outlined"
                density="compact"
                rows="2"
                auto-grow
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="addDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="saving"
            :disabled="!form.claimNumber || !form.notificationType"
            @click="handleAdd"
          >
            Add
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Generate Month-End Dialog -->
    <v-dialog v-model="generateDialog" max-width="480" persistent>
      <v-card>
        <v-card-title>Generate Month-End Notifications</v-card-title>
        <v-card-text>
          <p class="text-body-2 mb-4 text-medium-emphasis">
            Auto-generates initial, status update, and final notifications for
            all open claims on the selected scheme for the specified period.
          </p>
          <v-row>
            <v-col cols="12">
              <v-select
                v-model="generateForm.schemeId"
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
                v-model.number="generateForm.month"
                label="Month"
                :items="months"
                item-title="label"
                item-value="value"
                variant="outlined"
                density="compact"
                clearable
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model.number="generateForm.year"
                label="Year"
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
            color="amber-darken-2"
            :loading="generating"
            :disabled="!generateForm.schemeId"
            @click="handleGenerate"
          >
            Generate
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Mark Sent Dialog -->
    <v-dialog v-model="sentDialog.show" max-width="440" persistent>
      <v-card>
        <v-card-title>Mark as Sent</v-card-title>
        <v-card-text>
          <p class="text-body-2 mb-3 text-medium-emphasis">
            Notification #{{ sentDialog.id }} — {{ sentDialog.claimNumber }}
          </p>
          <v-textarea
            v-model="sentDialog.notes"
            label="Notes (optional)"
            variant="outlined"
            density="compact"
            rows="2"
            auto-grow
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="sentDialog.show = false">Cancel</v-btn>
          <v-btn color="blue" :loading="actionLoading" @click="handleMarkSent"
            >Mark Sent</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Mark Acknowledged Dialog -->
    <v-dialog v-model="ackDialog.show" max-width="440" persistent>
      <v-card>
        <v-card-title>Mark as Acknowledged</v-card-title>
        <v-card-text>
          <p class="text-body-2 mb-3 text-medium-emphasis">
            Notification #{{ ackDialog.id }} — {{ ackDialog.claimNumber }}
          </p>
          <v-textarea
            v-model="ackDialog.notes"
            label="Acknowledgement Notes (optional)"
            variant="outlined"
            density="compact"
            rows="2"
            auto-grow
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="ackDialog.show = false">Cancel</v-btn>
          <v-btn
            color="success"
            :loading="actionLoading"
            @click="handleMarkAcknowledged"
          >
            Acknowledge
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog.show" max-width="400" persistent>
      <v-card>
        <v-card-title>Delete Notification</v-card-title>
        <v-card-text>
          <p class="text-body-2 text-medium-emphasis">
            Are you sure you want to delete notification #{{
              deleteDialog.id
            }}
            — {{ deleteDialog.claimNumber }}? This cannot be undone.
          </p>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="deleteDialog.show = false"
            >Cancel</v-btn
          >
          <v-btn color="error" :loading="actionLoading" @click="handleDelete"
            >Delete</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Context Menu -->
    <teleport to="body">
      <v-card
        v-if="contextMenu.show"
        v-click-outside="() => (contextMenu.show = false)"
        elevation="8"
        rounded="lg"
        :style="{
          position: 'fixed',
          left: contextMenu.x + 'px',
          top: contextMenu.y + 'px',
          zIndex: 9999,
          minWidth: '180px'
        }"
      >
        <v-list density="compact" nav>
          <v-list-item
            v-for="(item, idx) in contextMenu.items"
            :key="idx"
            @click="handleContextAction(item.fn)"
          >
            <v-list-item-title
              :class="'text-' + item.color"
              style="font-size: 13px"
            >
              {{ item.label }}
            </v-list-item-title>
          </v-list-item>
        </v-list>
      </v-card>
    </teleport>

    <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="4000">
      {{ snackbar.message }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onBeforeUnmount } from 'vue'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-balham.css'
import { AgGridVue } from 'ag-grid-vue3'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import PremiumManagementService from '@/renderer/api/PremiumManagementService'

const loading = ref(false)
const saving = ref(false)
const loadingOverlay = '<span class="ag-overlay-loading-center">Loading…</span>'
const gridApi = ref<any>(null)
const onGridReady = (p: any) => {
  gridApi.value = p.api
  applyOverlay()
}
const applyOverlay = () => {
  if (!gridApi.value) return
  if (loading.value) gridApi.value.showLoadingOverlay()
  else if (!notifications.value?.length) gridApi.value.showNoRowsOverlay()
  else gridApi.value.hideOverlay()
}
watch(loading, applyOverlay)
const generating = ref(false)
const actionLoading = ref(false)
const schemesLoading = ref(false)
const exporting = ref(false)
const addDialog = ref(false)
const generateDialog = ref(false)
const notifications = ref<any[]>([])
watch(notifications, applyOverlay, { deep: true })
const inForceSchemes = ref<any[]>([])
const schemeClaims = ref<any[]>([])
const schemeClaimsLoading = ref(false)
const stats = ref({
  total: 0,
  pending: 0,
  sent: 0,
  acknowledged: 0,
  overdue: 0
})

const filters = ref({
  schemeId: null as number | null,
  notificationType: '',
  status: '',
  claimNumber: ''
})

const form = ref({
  schemeId: null as number | null,
  claimNumber: '',
  reinsurerName: '',
  reinsurerCode: '',
  notificationType: '',
  dueDate: '',
  notes: ''
})

const generateForm = ref({
  schemeId: null as number | null,
  month: new Date().getMonth() + 1,
  year: new Date().getFullYear()
})

const sentDialog = ref({ show: false, id: 0, claimNumber: '', notes: '' })
const ackDialog = ref({ show: false, id: 0, claimNumber: '', notes: '' })
const deleteDialog = ref({ show: false, id: 0, claimNumber: '' })
const snackbar = ref({ show: false, message: '', color: 'success' })

const notificationTypeOptions = ['initial', 'status_update', 'final']
const statusOptions = ['pending', 'sent', 'acknowledged', 'overdue']

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

const typeColors: Record<string, string> = {
  initial: '#9C27B0',
  status_update: '#FF9800',
  final: '#2196F3'
}

const statusColors: Record<string, string> = {
  pending: '#FF9800',
  sent: '#2196F3',
  acknowledged: '#4CAF50',
  overdue: '#F44336'
}

const badge = (value: string, colorMap: Record<string, string>) => {
  const c = colorMap[value] ?? '#9E9E9E'
  return `<span style="background:${c}22;color:${c};padding:2px 8px;border-radius:4px;font-size:12px;font-weight:500">${(value ?? '').replace(/_/g, ' ')}</span>`
}

const fmtDate = (v: string | null) => {
  if (!v) return ''
  return new Date(v).toLocaleDateString('en-ZA', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const defaultColDef = { sortable: true, filter: true, resizable: true }

const contextMenu = ref({
  show: false,
  x: 0,
  y: 0,
  items: [] as { label: string; color: string; fn: () => void }[]
})

const MENU_WIDTH = 200
const MENU_ITEM_HEIGHT = 40
const MENU_PAD = 8

function showContextMenu(event: MouseEvent, data: any) {
  const items: { label: string; color: string; fn: () => void }[] = []
  if (data.status === 'pending' || data.status === 'overdue') {
    items.push({
      label: 'Mark as Sent',
      color: 'blue',
      fn: () => openSentDialog(data.id, data.claim_number ?? '')
    })
  }
  if (data.status === 'sent') {
    items.push({
      label: 'Mark as Acknowledged',
      color: 'success',
      fn: () => openAckDialog(data.id, data.claim_number ?? '')
    })
  }
  if (data.status === 'pending') {
    items.push({
      label: 'Delete',
      color: 'error',
      fn: () => openDeleteDialog(data.id, data.claim_number ?? '')
    })
  }
  if (!items.length) return

  const btn = (event.currentTarget || event.target) as HTMLElement
  const rect = btn.getBoundingClientRect()
  const menuHeight = items.length * MENU_ITEM_HEIGHT + MENU_PAD * 2

  // Default: below and left-aligned to button
  let x = rect.left
  let y = rect.bottom + 4

  // If overflows right, anchor to left of button
  if (x + MENU_WIDTH > window.innerWidth - MENU_PAD) {
    x = rect.right - MENU_WIDTH
  }
  // If still overflows right, clamp
  if (x + MENU_WIDTH > window.innerWidth - MENU_PAD) {
    x = window.innerWidth - MENU_WIDTH - MENU_PAD
  }
  // If overflows bottom, show above
  if (y + menuHeight > window.innerHeight - MENU_PAD) {
    y = rect.top - menuHeight - 4
  }
  if (x < MENU_PAD) x = MENU_PAD
  if (y < MENU_PAD) y = MENU_PAD

  contextMenu.value = { show: true, x, y, items }
}

function handleContextAction(fn: () => void) {
  contextMenu.value.show = false
  fn()
}

const columnDefs = [
  { headerName: 'ID', field: 'id', width: 70 },
  { headerName: 'Claim No', field: 'claim_number', width: 140 },
  { headerName: 'Scheme', field: 'scheme_name', flex: 1, minWidth: 120 },
  { headerName: 'Reinsurer', field: 'reinsurer_name', flex: 1, minWidth: 120 },
  {
    headerName: 'Type',
    field: 'notification_type',
    width: 140,
    cellRenderer: (p: any) => badge(p.value, typeColors)
  },
  {
    headerName: 'Status',
    field: 'status',
    width: 130,
    cellRenderer: (p: any) => badge(p.value, statusColors)
  },
  {
    headerName: 'Due Date',
    field: 'due_date',
    width: 115,
    valueFormatter: (p: any) => fmtDate(p.value)
  },
  {
    headerName: 'Sent',
    field: 'sent_at',
    width: 115,
    valueFormatter: (p: any) => fmtDate(p.value)
  },
  {
    headerName: 'Acknowledged',
    field: 'acknowledged_at',
    width: 130,
    valueFormatter: (p: any) => fmtDate(p.value)
  },
  {
    headerName: 'Actions',
    width: 80,
    pinned: 'right',
    sortable: false,
    filter: false,
    cellRenderer: (p: any) => {
      const btn = document.createElement('button')
      btn.title = 'Actions'
      btn.style.cssText =
        'background:none;border:none;cursor:pointer;padding:4px 8px;border-radius:4px;color:#616161;display:flex;align-items:center;justify-content:center'
      btn.innerHTML =
        '<svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="5" r="2"/><circle cx="12" cy="12" r="2"/><circle cx="12" cy="19" r="2"/></svg>'
      btn.addEventListener('click', (e: MouseEvent) =>
        showContextMenu(e, p.data)
      )
      const wrapper = document.createElement('div')
      wrapper.style.cssText =
        'display:flex;align-items:center;justify-content:center;height:100%'
      wrapper.appendChild(btn)
      return wrapper
    }
  }
]

const loadNotifications = async () => {
  loading.value = true
  try {
    const params: any = { page: 1, page_size: 500 }
    if (filters.value.schemeId) params.scheme_id = filters.value.schemeId
    if (filters.value.notificationType)
      params.notification_type = filters.value.notificationType
    if (filters.value.status) params.status = filters.value.status
    if (filters.value.claimNumber)
      params.claim_number = filters.value.claimNumber
    const res = await GroupPricingService.getClaimNotifications(params)
    notifications.value = res.data?.data ?? []
  } catch {
    snackbar.value = {
      show: true,
      message: 'Failed to load notifications',
      color: 'error'
    }
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const params: any = {}
    if (filters.value.schemeId) params.scheme_id = filters.value.schemeId
    const res = await GroupPricingService.getNotificationStats(params)
    const d = res.data?.data
    if (d) {
      stats.value = {
        total: d.total ?? 0,
        pending: d.pending ?? 0,
        sent: d.sent ?? 0,
        acknowledged: d.acknowledged ?? 0,
        overdue: d.overdue ?? 0
      }
    }
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

// When scheme changes in add form, fetch claims and treaty info
watch(
  () => form.value.schemeId,
  async (schemeId) => {
    form.value.claimNumber = ''
    form.value.reinsurerName = ''
    form.value.reinsurerCode = ''
    schemeClaims.value = []
    if (!schemeId) return

    schemeClaimsLoading.value = true
    try {
      const [claimsRes, treatyRes] = await Promise.all([
        GroupPricingService.getClaimsByScheme(schemeId),
        GroupPricingService.getActiveTreatiesForScheme(schemeId)
      ])
      schemeClaims.value = claimsRes.data?.data ?? []
      const treaties = treatyRes.data?.data ?? treatyRes.data ?? []
      if (treaties.length > 0) {
        form.value.reinsurerName = treaties[0].reinsurer_name ?? ''
        form.value.reinsurerCode = treaties[0].reinsurer_code ?? ''
      }
    } catch {
      // non-fatal
    } finally {
      schemeClaimsLoading.value = false
    }
  }
)

const handleAdd = async () => {
  saving.value = true
  try {
    const scheme = inForceSchemes.value.find(
      (s) => s.id === form.value.schemeId
    )
    await GroupPricingService.createClaimNotification({
      scheme_id: form.value.schemeId,
      scheme_name: scheme?.name ?? '',
      claim_number: form.value.claimNumber,
      reinsurer_name: form.value.reinsurerName,
      reinsurer_code: form.value.reinsurerCode,
      notification_type: form.value.notificationType,
      due_date: form.value.dueDate,
      notes: form.value.notes
    })
    addDialog.value = false
    form.value = {
      schemeId: null,
      claimNumber: '',
      reinsurerName: '',
      reinsurerCode: '',
      notificationType: '',
      dueDate: '',
      notes: ''
    }
    snackbar.value = {
      show: true,
      message: 'Notification added',
      color: 'success'
    }
    await Promise.all([loadNotifications(), loadStats()])
  } catch {
    snackbar.value = {
      show: true,
      message: 'Failed to add notification',
      color: 'error'
    }
  } finally {
    saving.value = false
  }
}

const handleGenerate = async () => {
  generating.value = true
  try {
    const res = await GroupPricingService.generateMonthEndNotifications({
      scheme_id: generateForm.value.schemeId,
      month: generateForm.value.month,
      year: generateForm.value.year
    })
    generateDialog.value = false
    const count = res.data?.count ?? 0
    snackbar.value = {
      show: true,
      message: `Generated ${count} notification(s)`,
      color: 'success'
    }
    await Promise.all([loadNotifications(), loadStats()])
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

const openSentDialog = (id: number, claimNumber: string) => {
  sentDialog.value = { show: true, id, claimNumber, notes: '' }
}

const openAckDialog = (id: number, claimNumber: string) => {
  ackDialog.value = { show: true, id, claimNumber, notes: '' }
}

const openDeleteDialog = (id: number, claimNumber: string) => {
  deleteDialog.value = { show: true, id, claimNumber }
}

const handleDelete = async () => {
  actionLoading.value = true
  try {
    await GroupPricingService.deleteClaimNotification(deleteDialog.value.id)
    deleteDialog.value.show = false
    snackbar.value = {
      show: true,
      message: 'Notification deleted',
      color: 'success'
    }
    await Promise.all([loadNotifications(), loadStats()])
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Delete failed',
      color: 'error'
    }
  } finally {
    actionLoading.value = false
  }
}

const handleExport = async () => {
  exporting.value = true
  try {
    const params: any = {}
    if (filters.value.schemeId) params.scheme_id = filters.value.schemeId
    if (filters.value.notificationType)
      params.notification_type = filters.value.notificationType
    if (filters.value.status) params.status = filters.value.status
    const res = await GroupPricingService.exportNotificationsCSV(params)
    const url = window.URL.createObjectURL(new Blob([res.data]))
    const a = document.createElement('a')
    a.href = url
    a.download = 'claim_notifications.csv'
    a.click()
    window.URL.revokeObjectURL(url)
  } catch {
    snackbar.value = { show: true, message: 'Export failed', color: 'error' }
  } finally {
    exporting.value = false
  }
}

const handleMarkSent = async () => {
  actionLoading.value = true
  try {
    await GroupPricingService.markNotificationSent(sentDialog.value.id, {
      notes: sentDialog.value.notes
    })
    sentDialog.value.show = false
    snackbar.value = {
      show: true,
      message: 'Notification marked as sent',
      color: 'success'
    }
    await Promise.all([loadNotifications(), loadStats()])
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Action failed',
      color: 'error'
    }
  } finally {
    actionLoading.value = false
  }
}

const handleMarkAcknowledged = async () => {
  actionLoading.value = true
  try {
    await GroupPricingService.markNotificationAcknowledged(ackDialog.value.id, {
      notes: ackDialog.value.notes
    })
    ackDialog.value.show = false
    snackbar.value = {
      show: true,
      message: 'Notification acknowledged',
      color: 'success'
    }
    await Promise.all([loadNotifications(), loadStats()])
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Action failed',
      color: 'error'
    }
  } finally {
    actionLoading.value = false
  }
}

let filterDebounce: ReturnType<typeof setTimeout> | null = null
watch(
  filters,
  () => {
    if (filterDebounce) clearTimeout(filterDebounce)
    filterDebounce = setTimeout(() => {
      loadNotifications()
      loadStats()
    }, 400)
  },
  { deep: true }
)

onMounted(() => {
  loadInForceSchemes()
  loadNotifications()
  loadStats()
})

onBeforeUnmount(() => {
  if (filterDebounce) clearTimeout(filterDebounce)
})
</script>
