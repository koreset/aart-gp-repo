<template>
  <v-container>
    <base-card :show-actions="false">
      <template #header>
        <div class="d-flex justify-space-between align-center flex-wrap">
          <span class="headline">Journal Entries</span>
          <v-btn
            v-if="hasPermission('gl:draft_journal')"
            color="primary"
            prepend-icon="mdi-plus"
            @click="showManual = true"
          >
            New Manual Journal
          </v-btn>
        </div>
      </template>
      <template #default>
        <v-row dense class="mb-2">
          <v-col cols="12" md="3">
            <v-select
              v-model="filters.status"
              :items="statusOptions"
              item-title="label"
              item-value="value"
              label="Status"
              density="compact"
              variant="outlined"
              clearable
              @update:model-value="load"
            />
          </v-col>
          <v-col cols="12" md="3">
            <v-select
              v-model="filters.period_id"
              :items="periodOptions"
              item-title="name"
              item-value="id"
              label="Period"
              density="compact"
              variant="outlined"
              clearable
              @update:model-value="load"
            />
          </v-col>
          <v-col cols="12" md="2">
            <v-select
              v-model="filters.source_type"
              :items="sourceTypes"
              label="Source"
              density="compact"
              variant="outlined"
              clearable
              @update:model-value="load"
            />
          </v-col>
          <v-col cols="12" md="2">
            <v-text-field
              v-model="filters.from"
              type="date"
              label="From"
              density="compact"
              variant="outlined"
              @change="load"
            />
          </v-col>
          <v-col cols="12" md="2">
            <v-text-field
              v-model="filters.to"
              type="date"
              label="To"
              density="compact"
              variant="outlined"
              @change="load"
            />
          </v-col>
        </v-row>

        <v-data-table
          :headers="headers"
          :items="entries"
          :loading="loading"
          density="compact"
          items-per-page="50"
          @click:row="openDetail"
        >
          <template #[`item.status`]="{ value }">
            <v-chip
              :color="statusColor(value)"
              size="x-small"
              variant="tonal"
              class="text-uppercase"
            >
              {{ statusLabel(value) }}
            </v-chip>
          </template>
          <template #[`item.entry_number`]="{ item }">
            {{ item.entry_number || `(draft #${item.id})` }}
          </template>
          <template #[`item.posted_at`]="{ value }">
            {{ value ? new Date(value).toLocaleDateString() : '' }}
          </template>
          <template #[`item.total_debit`]="{ value }">
            {{ formatAmount(value) }}
          </template>
          <template #[`item.actions`]="{ item }">
            <div class="d-flex ga-1 justify-end">
              <v-btn
                v-if="canApprove(item)"
                size="x-small"
                color="success"
                variant="tonal"
                @click.stop="approveEntry(item)"
              >
                Approve
              </v-btn>
              <v-btn
                v-if="canPost(item)"
                size="x-small"
                color="primary"
                variant="tonal"
                @click.stop="postEntry(item)"
              >
                Post
              </v-btn>
              <v-btn
                v-if="canApproveReversal(item)"
                size="x-small"
                color="warning"
                variant="tonal"
                @click.stop="approveReversal(item)"
              >
                Approve Reversal
              </v-btn>
            </div>
          </template>
        </v-data-table>

        <v-alert
          v-if="actionError"
          class="mt-3"
          color="error"
          variant="tonal"
          density="compact"
          icon="mdi-alert"
          closable
          @click:close="actionError = ''"
        >
          {{ actionError }}
        </v-alert>
      </template>
    </base-card>

    <ManualJournalDialog
      v-model="showManual"
      :accounts="accountOptions"
      :periods="periodOptions"
      @posted="load"
    />
  </v-container>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import GLService, {
  AccountingPeriod,
  GLAccount,
  JournalEntry,
  JournalStatus
} from '@/renderer/api/GeneralLedgerService'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'
import ManualJournalDialog from './ManualJournalDialog.vue'

const { hasPermission } = usePermissionCheck()
const router = useRouter()

const currentUserName = computed<string>(() => {
  const u: any = window.mainApi?.sendSync('msgGetAuthenticatedUser')
  return u?.full_name || u?.username || ''
})

const headers = [
  { title: 'Status', key: 'status' },
  { title: 'Entry #', key: 'entry_number', sortable: true },
  { title: 'Posted', key: 'posted_at' },
  { title: 'Source', key: 'source_type' },
  { title: 'Source ID', key: 'source_id' },
  { title: 'Description', key: 'description' },
  { title: 'Debit', key: 'total_debit', align: 'end' as const },
  { title: '', key: 'actions', align: 'end' as const, sortable: false }
]

const statusOptions = [
  { label: 'Draft', value: 'draft' },
  { label: 'Submitted', value: 'submitted' },
  { label: 'Approved', value: 'approved' },
  { label: 'Posted', value: 'posted' },
  { label: 'Reversal Pending', value: 'reversal_pending' },
  { label: 'Reversed', value: 'reversed' }
]

const sourceTypes = [
  'manual',
  'claim_payment',
  'premium_allocation',
  'premium_allocation_reversal',
  'write_off',
  'refund',
  'reversal'
]

const entries = ref<JournalEntry[]>([])
const periods = ref<AccountingPeriod[]>([])
const accounts = ref<GLAccount[]>([])
const loading = ref(false)
const showManual = ref(false)
const actionError = ref('')

const filters = reactive({
  status: undefined as string | undefined,
  period_id: undefined as number | undefined,
  source_type: undefined as string | undefined,
  from: '',
  to: ''
})

const periodOptions = computed(() => periods.value)
const accountOptions = computed(() =>
  accounts.value.filter((a) => a.is_active)
)

const formatAmount = (n: number) =>
  new Intl.NumberFormat(undefined, { minimumFractionDigits: 2 }).format(n || 0)

const statusLabel = (s?: JournalStatus | string) => {
  switch (s) {
    case 'reversal_pending':
      return 'Reversal Pending'
    case 'reversal_approved':
      return 'Reversal Approved'
    default:
      return s || 'posted'
  }
}
const statusColor = (s?: JournalStatus | string) => {
  switch (s) {
    case 'draft':
      return 'grey'
    case 'submitted':
      return 'info'
    case 'approved':
      return 'primary'
    case 'posted':
      return 'success'
    case 'reversal_pending':
    case 'reversal_approved':
      return 'warning'
    case 'reversed':
      return 'error'
    default:
      return 'default'
  }
}

const sameUser = (a?: string, b?: string) =>
  Boolean(a && b && a.trim().toLowerCase() === b.trim().toLowerCase())

const canApprove = (e: JournalEntry) =>
  e.status === 'submitted' &&
  hasPermission('gl:approve_journal') &&
  !sameUser(e.submitted_by, currentUserName.value)

const canPost = (e: JournalEntry) =>
  e.status === 'approved' && hasPermission('gl:post_journal')

const canApproveReversal = (e: JournalEntry) =>
  e.status === 'reversal_pending' &&
  hasPermission('gl:approve_reverse') &&
  !sameUser(e.reversal_requested_by, currentUserName.value)

const runAction = async (fn: () => Promise<unknown>) => {
  actionError.value = ''
  try {
    await fn()
    await load()
  } catch (e: any) {
    actionError.value =
      e?.response?.data?.error || e?.message || 'Action failed'
  }
}

const approveEntry = (e: JournalEntry) =>
  e.id && runAction(() => GLService.approveManualJournal(e.id!))
const postEntry = (e: JournalEntry) =>
  e.id && runAction(() => GLService.postApprovedJournal(e.id!))
const approveReversal = (e: JournalEntry) =>
  e.id && runAction(() => GLService.approveReverseJournal(e.id!))

const load = async () => {
  loading.value = true
  try {
    const params: any = {}
    if (filters.status) params.status = filters.status
    if (filters.period_id) params.period_id = filters.period_id
    if (filters.source_type) params.source_type = filters.source_type
    if (filters.from) params.from = filters.from
    if (filters.to) params.to = filters.to
    entries.value = await GLService.listJournals(params)
  } finally {
    loading.value = false
  }
}

const openDetail = (_evt: unknown, row: { item: JournalEntry }) => {
  if (row.item.id) {
    router.push({
      name: 'group-pricing-gl-journal-detail',
      params: { id: row.item.id }
    })
  }
}

onMounted(async () => {
  const [p, a] = await Promise.all([
    GLService.listPeriods(),
    GLService.listAccounts()
  ])
  periods.value = p
  accounts.value = a
  await load()
})
</script>
