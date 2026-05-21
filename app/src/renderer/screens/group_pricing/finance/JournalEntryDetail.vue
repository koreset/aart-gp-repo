<template>
  <v-container>
    <base-card :show-actions="false">
      <template #header>
        <div class="d-flex align-center flex-wrap ga-2">
          <v-btn
            icon="mdi-arrow-left"
            variant="text"
            class="mr-2"
            @click="$router.back()"
          />
          <span class="headline">
            {{ entry?.entry_number || (entry ? `Draft #${entry.id}` : 'Journal Entry') }}
          </span>
          <v-chip
            v-if="entry"
            :color="statusColor(entry.status)"
            size="small"
            variant="tonal"
            class="ml-2 text-uppercase"
          >
            {{ statusLabel(entry.status) }}
          </v-chip>
          <v-spacer />

          <!-- Draft owner controls -->
          <v-btn
            v-if="canEditDraft"
            color="primary"
            prepend-icon="mdi-pencil"
            variant="text"
            @click="$emit('edit-draft')"
          >
            Edit (use New Manual Journal dialog)
          </v-btn>
          <v-btn
            v-if="canSubmit"
            color="info"
            prepend-icon="mdi-send"
            @click="runAction(() => GLService.submitManualJournal(entry!.id!))"
          >
            Submit for Approval
          </v-btn>
          <v-btn
            v-if="canDiscard"
            color="error"
            prepend-icon="mdi-delete"
            variant="text"
            @click="discardDialog = true"
          >
            Discard
          </v-btn>

          <!-- Approval -->
          <v-btn
            v-if="canApprove"
            color="success"
            prepend-icon="mdi-check"
            @click="runAction(() => GLService.approveManualJournal(entry!.id!))"
          >
            Approve
          </v-btn>

          <!-- Posting -->
          <v-btn
            v-if="canPost"
            color="primary"
            prepend-icon="mdi-book-arrow-right"
            @click="runAction(() => GLService.postApprovedJournal(entry!.id!))"
          >
            Post to Ledger
          </v-btn>

          <!-- Reversal -->
          <v-btn
            v-if="canRequestReverse"
            color="warning"
            prepend-icon="mdi-undo"
            @click="reverseDialog = true"
          >
            Request Reversal
          </v-btn>
          <v-btn
            v-if="canApproveReverse"
            color="warning"
            prepend-icon="mdi-undo-variant"
            @click="runAction(approveReverseAndShowResult)"
          >
            Approve Reversal
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

        <v-row v-if="entry" dense>
          <v-col cols="12" md="6">
            <v-list density="compact">
              <v-list-item v-if="entry.posted_at">
                <v-list-item-title>Posted</v-list-item-title>
                <v-list-item-subtitle>
                  {{ new Date(entry.posted_at).toLocaleString() }}
                </v-list-item-subtitle>
              </v-list-item>
              <v-list-item v-if="entry.posted_by">
                <v-list-item-title>Posted by</v-list-item-title>
                <v-list-item-subtitle>{{ entry.posted_by }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title>Period</v-list-item-title>
                <v-list-item-subtitle>#{{ entry.period_id }}</v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-col>
          <v-col cols="12" md="6">
            <v-list density="compact">
              <v-list-item>
                <v-list-item-title>Source</v-list-item-title>
                <v-list-item-subtitle>
                  {{ entry.source_type }}{{ entry.source_id ? ` #${entry.source_id}` : '' }}
                </v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title>Description</v-list-item-title>
                <v-list-item-subtitle>{{ entry.description }}</v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-col>

          <v-col cols="12">
            <v-card variant="tonal" color="surface" class="my-3" density="compact">
              <v-card-title class="text-subtitle-2">Audit Trail</v-card-title>
              <v-card-text class="text-body-2">
                <div v-if="entry.created_by || entry.created_at">
                  <strong>Drafted by</strong> {{ entry.created_by || '—' }}
                  <span v-if="entry.created_at">
                    on {{ new Date(entry.created_at).toLocaleString() }}
                  </span>
                </div>
                <div v-if="entry.updated_by && entry.updated_by !== entry.created_by">
                  <strong>Last updated by</strong> {{ entry.updated_by }}
                  <span v-if="entry.updated_at">
                    on {{ new Date(entry.updated_at).toLocaleString() }}
                  </span>
                </div>
                <div v-if="entry.submitted_by">
                  <strong>Submitted by</strong> {{ entry.submitted_by }}
                  <span v-if="entry.submitted_at">
                    on {{ new Date(entry.submitted_at).toLocaleString() }}
                  </span>
                </div>
                <div v-if="entry.approved_by">
                  <strong>Approved by</strong> {{ entry.approved_by }}
                  <span v-if="entry.approved_at">
                    on {{ new Date(entry.approved_at).toLocaleString() }}
                  </span>
                </div>
                <div v-if="entry.posted_by && entry.status === 'posted'">
                  <strong>Posted by</strong> {{ entry.posted_by }}
                  <span v-if="entry.posted_at">
                    on {{ new Date(entry.posted_at).toLocaleString() }}
                  </span>
                </div>
                <div v-if="entry.reversal_requested_by">
                  <strong>Reversal requested by</strong>
                  {{ entry.reversal_requested_by }}
                  <span v-if="entry.reversal_requested_at">
                    on {{ new Date(entry.reversal_requested_at).toLocaleString() }}
                  </span>
                  <em v-if="entry.reversal_reason"> — “{{ entry.reversal_reason }}”</em>
                </div>
                <div v-if="entry.reversal_approved_by">
                  <strong>Reversal approved by</strong>
                  {{ entry.reversal_approved_by }}
                  <span v-if="entry.reversal_approved_at">
                    on {{ new Date(entry.reversal_approved_at).toLocaleString() }}
                  </span>
                </div>
              </v-card-text>
            </v-card>

            <v-alert
              v-if="entry.is_reversed"
              color="warning"
              variant="tonal"
              icon="mdi-undo"
              class="my-2"
            >
              This entry has been reversed by entry #{{ entry.reversed_by_entry_id }}.
            </v-alert>

            <v-data-table
              :headers="headers"
              :items="entry.lines"
              density="compact"
              hide-default-footer
              items-per-page="-1"
            >
              <template #[`item.account`]="{ item }">
                {{ accountLabel(item.account_id) }}
              </template>
              <template #[`item.debit`]="{ value }">
                <span class="text-right d-block">{{ format(value) }}</span>
              </template>
              <template #[`item.credit`]="{ value }">
                <span class="text-right d-block">{{ format(value) }}</span>
              </template>
            </v-data-table>
            <div class="d-flex justify-end mt-2 text-subtitle-2">
              <span class="mr-6">
                Total Debit: <strong>{{ format(entry.total_debit || 0) }}</strong>
              </span>
              <span>
                Total Credit: <strong>{{ format(entry.total_credit || 0) }}</strong>
              </span>
            </div>
          </v-col>
        </v-row>
      </template>
    </base-card>

    <v-dialog v-model="reverseDialog" max-width="500">
      <v-card>
        <v-card-title>Request Reversal</v-card-title>
        <v-card-text>
          <v-textarea
            v-model="reverseReason"
            label="Reason"
            rows="2"
            variant="outlined"
            density="compact"
            required
          />
          <v-alert color="info" variant="tonal" density="compact" icon="mdi-information-outline" class="mt-2">
            The reversal will not post until a different user approves it.
          </v-alert>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="reverseDialog = false">Cancel</v-btn>
          <v-btn
            color="warning"
            :loading="reversing"
            :disabled="!reverseReason.trim()"
            @click="requestReverse"
          >
            Submit Reversal Request
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="discardDialog" max-width="420">
      <v-card>
        <v-card-title>Discard draft?</v-card-title>
        <v-card-text>
          This will delete the draft permanently. The action is recorded in the audit log.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="discardDialog = false">Cancel</v-btn>
          <v-btn color="error" @click="confirmDiscard">Discard</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import GLService, {
  GLAccount,
  JournalEntry,
  JournalStatus
} from '@/renderer/api/GeneralLedgerService'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

const { hasPermission } = usePermissionCheck()
const route = useRoute()
const router = useRouter()

const currentUserName = computed<string>(() => {
  const u: any = window.mainApi?.sendSync('msgGetAuthenticatedUser')
  return u?.full_name || u?.username || ''
})

const headers = [
  { title: 'Account', key: 'account' },
  { title: 'Description', key: 'description' },
  { title: 'Debit', key: 'debit', align: 'end' as const },
  { title: 'Credit', key: 'credit', align: 'end' as const }
]

const entry = ref<JournalEntry | null>(null)
const accounts = ref<GLAccount[]>([])
const reverseDialog = ref(false)
const reverseReason = ref('')
const reversing = ref(false)
const discardDialog = ref(false)
const actionError = ref('')

const format = (n: number) =>
  new Intl.NumberFormat(undefined, { minimumFractionDigits: 2 }).format(n || 0)

const accountLabel = (id: number) => {
  const a = accounts.value.find((x) => x.id === id)
  return a ? `${a.code} — ${a.name}` : `#${id}`
}

const sameUser = (a?: string, b?: string) =>
  Boolean(a && b && a.trim().toLowerCase() === b.trim().toLowerCase())

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

const canEditDraft = computed(
  () =>
    entry.value?.status === 'draft' &&
    hasPermission('gl:draft_journal') &&
    sameUser(entry.value.created_by, currentUserName.value)
)
const canSubmit = computed(
  () => entry.value?.status === 'draft' && hasPermission('gl:draft_journal')
)
const canDiscard = computed(
  () =>
    (entry.value?.status === 'draft' || entry.value?.status === 'submitted') &&
    hasPermission('gl:draft_journal')
)
const canApprove = computed(
  () =>
    entry.value?.status === 'submitted' &&
    hasPermission('gl:approve_journal') &&
    !sameUser(entry.value.submitted_by, currentUserName.value)
)
const canPost = computed(
  () => entry.value?.status === 'approved' && hasPermission('gl:post_journal')
)
const canRequestReverse = computed(
  () =>
    entry.value?.status === 'posted' &&
    !entry.value?.is_reversed &&
    hasPermission('gl:reverse')
)
const canApproveReverse = computed(
  () =>
    entry.value?.status === 'reversal_pending' &&
    hasPermission('gl:approve_reverse') &&
    !sameUser(entry.value.reversal_requested_by, currentUserName.value)
)

const load = async () => {
  const id = Number(route.params.id)
  const [e, a] = await Promise.all([
    GLService.getJournalEntry(id),
    GLService.listAccounts()
  ])
  entry.value = e
  accounts.value = a
}

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

const requestReverse = async () => {
  if (!entry.value?.id) return
  reversing.value = true
  try {
    await GLService.requestReverseJournal(entry.value.id, reverseReason.value.trim())
    reverseDialog.value = false
    reverseReason.value = ''
    await load()
  } catch (e: any) {
    actionError.value =
      e?.response?.data?.error || e?.message || 'Failed to request reversal'
  } finally {
    reversing.value = false
  }
}

const approveReverseAndShowResult = async () => {
  if (!entry.value?.id) return
  const created = await GLService.approveReverseJournal(entry.value.id)
  if (created?.id && created.id !== entry.value.id) {
    router.replace({
      name: 'group-pricing-gl-journal-detail',
      params: { id: created.id }
    })
  }
}

const confirmDiscard = async () => {
  if (!entry.value?.id) return
  await GLService.discardDraftJournal(entry.value.id)
  discardDialog.value = false
  router.back()
}

onMounted(load)
</script>
