<template>
  <v-container>
    <base-card :show-actions="false">
      <template #header>
        <div class="d-flex justify-space-between align-center flex-wrap">
          <span class="headline">Accounting Periods</span>
          <v-btn
            v-if="hasPermission('gl:manage_accounts')"
            color="primary"
            prepend-icon="mdi-plus"
            @click="dialog = true"
          >
            Open New Period
          </v-btn>
        </div>
      </template>
      <template #default>
        <v-alert
          v-if="actionError"
          color="error"
          variant="tonal"
          density="compact"
          icon="mdi-alert"
          class="mb-3"
          closable
          @click:close="actionError = ''"
        >
          {{ actionError }}
        </v-alert>

        <v-data-table
          :headers="headers"
          :items="periods"
          :loading="loading"
          density="compact"
          items-per-page="50"
        >
          <template #[`item.status`]="{ value }">
            <v-chip
              :color="statusColor(value)"
              size="small"
              variant="tonal"
              class="text-uppercase"
            >
              {{ statusLabel(value) }}
            </v-chip>
          </template>
          <template #[`item.start_date`]="{ value }">
            {{ new Date(value).toLocaleDateString() }}
          </template>
          <template #[`item.end_date`]="{ value }">
            {{ new Date(value).toLocaleDateString() }}
          </template>
          <template #[`item.close_audit`]="{ item }">
            <div class="text-caption">
              <div v-if="item.close_requested_by">
                Close requested by <strong>{{ item.close_requested_by }}</strong>
                <span v-if="item.close_requested_at">
                  on {{ new Date(item.close_requested_at).toLocaleString() }}
                </span>
              </div>
              <div v-if="item.closed_by">
                Closed by <strong>{{ item.closed_by }}</strong>
                <span v-if="item.closed_at">
                  on {{ new Date(item.closed_at).toLocaleString() }}
                </span>
              </div>
            </div>
          </template>
          <template #[`item.actions`]="{ item }">
            <v-btn
              v-if="item.status === 'open' && hasPermission('gl:request_close_period')"
              color="warning"
              variant="text"
              size="small"
              prepend-icon="mdi-flag"
              @click="requestClose(item)"
            >
              Request Close
            </v-btn>
            <v-btn
              v-if="
                item.status === 'close_requested' &&
                hasPermission('gl:close_period') &&
                !sameUser(item.close_requested_by, currentUserName)
              "
              color="error"
              variant="tonal"
              size="small"
              prepend-icon="mdi-lock"
              @click="confirmClose(item)"
            >
              Approve & Close
            </v-btn>
            <v-tooltip
              v-else-if="
                item.status === 'close_requested' &&
                sameUser(item.close_requested_by, currentUserName)
              "
              text="Another user must approve and close (you requested it)"
              location="top"
            >
              <template #activator="{ props: tip }">
                <v-chip v-bind="tip" size="small" variant="tonal" color="warning">
                  Awaiting approval
                </v-chip>
              </template>
            </v-tooltip>
          </template>
        </v-data-table>
      </template>
    </base-card>

    <v-dialog v-model="dialog" max-width="400">
      <v-card>
        <v-card-title>Open New Period</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="newStart"
            type="date"
            label="Any date inside the target month"
            density="compact"
            variant="outlined"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="dialog = false">Cancel</v-btn>
          <v-btn color="primary" :loading="saving" @click="create">Open</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import GLService, {
  AccountingPeriod
} from '@/renderer/api/GeneralLedgerService'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

const { hasPermission } = usePermissionCheck()

const currentUserName = computed<string>(() => {
  const u: any = window.mainApi?.sendSync('msgGetAuthenticatedUser')
  return u?.full_name || u?.username || ''
})

const headers = [
  { title: 'Name', key: 'name' },
  { title: 'Start', key: 'start_date' },
  { title: 'End', key: 'end_date' },
  { title: 'Status', key: 'status' },
  { title: 'Close audit', key: 'close_audit', sortable: false },
  { title: '', key: 'actions', sortable: false, align: 'end' as const }
]

const periods = ref<AccountingPeriod[]>([])
const loading = ref(false)
const dialog = ref(false)
const saving = ref(false)
const newStart = ref('')
const actionError = ref('')

const sameUser = (a?: string, b?: string) =>
  Boolean(a && b && a.trim().toLowerCase() === b.trim().toLowerCase())

const statusLabel = (s?: string) => {
  if (s === 'close_requested') return 'Close Requested'
  return s || 'open'
}
const statusColor = (s?: string) => {
  if (s === 'open') return 'success'
  if (s === 'close_requested') return 'warning'
  return 'grey'
}

const load = async () => {
  loading.value = true
  try {
    periods.value = await GLService.listPeriods()
  } finally {
    loading.value = false
  }
}

const create = async () => {
  if (!newStart.value) return
  saving.value = true
  try {
    await GLService.createPeriod({ start_date: newStart.value })
    dialog.value = false
    newStart.value = ''
    await load()
  } catch (e: any) {
    actionError.value =
      e?.response?.data?.error || e?.message || 'Failed to open period'
  } finally {
    saving.value = false
  }
}

const requestClose = async (p: AccountingPeriod) => {
  if (!p.id) return
  if (
    !confirm(
      `Request close of period ${p.name}? A different user will need to approve and commit the close.`
    )
  )
    return
  try {
    await GLService.requestClosePeriod(p.id)
    await load()
  } catch (e: any) {
    actionError.value =
      e?.response?.data?.error || e?.message || 'Failed to request close'
  }
}

const confirmClose = async (p: AccountingPeriod) => {
  if (!p.id) return
  if (
    !confirm(
      `Approve close of period ${p.name}? Postings to this period will be blocked.`
    )
  )
    return
  try {
    await GLService.closePeriod(p.id)
    await load()
  } catch (e: any) {
    actionError.value =
      e?.response?.data?.error || e?.message || 'Failed to close period'
  }
}

onMounted(load)
</script>
