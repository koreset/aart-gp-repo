<template>
  <v-container>
    <base-card :show-actions="false">
      <template #header>
        <div class="d-flex justify-space-between align-center flex-wrap">
          <span class="headline">Posting Rules</span>
          <v-btn
            v-if="hasPermission('gl:manage_rules')"
            color="primary"
            prepend-icon="mdi-plus"
            @click="openCreate"
          >
            New Rule
          </v-btn>
        </div>
      </template>
      <template #default>
        <v-alert color="info" variant="tonal" density="compact" class="mb-3" icon="mdi-shield-check-outline">
          Posting rules drive automatic GL postings. Every change is staged
          for a second user to approve — the live rule keeps its current
          values until released.
        </v-alert>
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
          :items="rules"
          :loading="loading"
          density="compact"
          items-per-page="50"
        >
          <template #[`item.debit_account_id`]="{ value }">
            {{ accountLabel(value) }}
          </template>
          <template #[`item.credit_account_id`]="{ value }">
            {{ accountLabel(value) }}
          </template>
          <template #[`item.is_active`]="{ value }">
            <v-icon :color="value ? 'success' : 'grey'">
              {{ value ? 'mdi-check-circle' : 'mdi-pause-circle' }}
            </v-icon>
          </template>
          <template #[`item.approval_status`]="{ item }">
            <v-chip
              v-if="item.approval_status && item.approval_status !== 'active'"
              size="x-small"
              color="warning"
              variant="tonal"
              class="text-uppercase"
            >
              {{ item.approval_status.replace('pending_', '') }} pending
            </v-chip>
            <span
              v-if="item.pending_requested_by"
              class="text-caption text-medium-emphasis ml-2"
            >
              by {{ item.pending_requested_by }}
            </span>
          </template>
          <template #[`item.actions`]="{ item }">
            <v-btn
              v-if="canApprove(item)"
              size="x-small"
              color="success"
              variant="tonal"
              @click="approve(item)"
            >
              Approve change
            </v-btn>
            <v-btn
              v-else-if="
                item.approval_status &&
                item.approval_status !== 'active' &&
                sameUser(item.pending_requested_by, currentUserName)
              "
              size="x-small"
              variant="tonal"
              color="warning"
              disabled
            >
              Awaiting approval
            </v-btn>
            <template v-else>
              <v-btn
                v-if="hasPermission('gl:manage_rules')"
                icon="mdi-pencil"
                variant="text"
                size="small"
                @click="openEdit(item)"
              />
              <v-btn
                v-if="hasPermission('gl:manage_rules')"
                icon="mdi-delete"
                variant="text"
                size="small"
                @click="remove(item)"
              />
            </template>
          </template>
        </v-data-table>
      </template>
    </base-card>

    <v-dialog v-model="dialog" max-width="640">
      <v-card>
        <v-card-title>
          {{ form.id ? 'Request Update' : 'Request New' }} Posting Rule
        </v-card-title>
        <v-card-text>
          <v-row dense>
            <v-col cols="12">
              <v-text-field
                v-model="form.event_key"
                label="Event key (e.g. claim_payment.confirmed)"
                :disabled="!!form.id"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-autocomplete
                v-model="form.debit_account_id"
                :items="accounts"
                :item-title="accountLabel"
                item-value="id"
                label="Debit account"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-autocomplete
                v-model="form.credit_account_id"
                :items="accounts"
                :item-title="accountLabel"
                item-value="id"
                label="Credit account"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="form.notes"
                label="Notes"
                rows="2"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12">
              <v-switch
                v-model="form.is_active"
                color="success"
                label="Active"
                inset
                hide-details
              />
            </v-col>
          </v-row>
          <v-alert color="info" variant="tonal" density="compact" icon="mdi-shield-check-outline">
            This change is staged and will only take effect once a different user approves it.
          </v-alert>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="dialog = false">Cancel</v-btn>
          <v-btn color="primary" :loading="saving" @click="save">
            Submit for Approval
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import GLService, {
  GLAccount,
  PostingRule
} from '@/renderer/api/GeneralLedgerService'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

const { hasPermission } = usePermissionCheck()

const currentUserName = computed<string>(() => {
  const u: any = window.mainApi?.sendSync('msgGetAuthenticatedUser')
  return u?.full_name || u?.username || ''
})

const headers = [
  { title: 'Event Key', key: 'event_key' },
  { title: 'Debit', key: 'debit_account_id' },
  { title: 'Credit', key: 'credit_account_id' },
  { title: 'Active', key: 'is_active' },
  { title: 'Notes', key: 'notes' },
  { title: 'Pending', key: 'approval_status', sortable: false },
  { title: '', key: 'actions', sortable: false, align: 'end' as const }
]

const rules = ref<PostingRule[]>([])
const accounts = ref<GLAccount[]>([])
const loading = ref(false)
const dialog = ref(false)
const saving = ref(false)
const actionError = ref('')
const form = reactive<PostingRule>({
  event_key: '',
  debit_account_id: 0,
  credit_account_id: 0,
  is_active: true,
  notes: ''
})

const sameUser = (a?: string, b?: string) =>
  Boolean(a && b && a.trim().toLowerCase() === b.trim().toLowerCase())

const accountLabel = (a: GLAccount | number) => {
  const acc =
    typeof a === 'number' ? accounts.value.find((x) => x.id === a) : a
  return acc ? `${acc.code} — ${acc.name}` : ''
}

const canApprove = (r: PostingRule) =>
  r.approval_status &&
  r.approval_status !== 'active' &&
  hasPermission('gl:approve_rule') &&
  !sameUser(r.pending_requested_by, currentUserName.value)

const load = async () => {
  loading.value = true
  try {
    const [r, a] = await Promise.all([
      GLService.listPostingRules(),
      GLService.listAccounts()
    ])
    rules.value = r
    accounts.value = a
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  Object.assign(form, {
    id: undefined,
    event_key: '',
    debit_account_id: 0,
    credit_account_id: 0,
    is_active: true,
    notes: ''
  })
  dialog.value = true
}

const openEdit = (r: PostingRule) => {
  Object.assign(form, r)
  dialog.value = true
}

const save = async () => {
  saving.value = true
  actionError.value = ''
  try {
    if (form.id) {
      await GLService.requestUpdatePostingRule(form.id, form)
    } else {
      await GLService.requestCreatePostingRule(form)
    }
    dialog.value = false
    await load()
  } catch (e: any) {
    actionError.value =
      e?.response?.data?.error || e?.message || 'Failed to stage change'
  } finally {
    saving.value = false
  }
}

const remove = async (r: PostingRule) => {
  if (!r.id) return
  if (
    !confirm(
      `Request deletion of posting rule for ${r.event_key}? A different user must approve.`
    )
  )
    return
  try {
    await GLService.requestDeletePostingRule(r.id)
    await load()
  } catch (e: any) {
    actionError.value =
      e?.response?.data?.error || e?.message || 'Failed to request deletion'
  }
}

const approve = async (r: PostingRule) => {
  if (!r.id) return
  try {
    await GLService.approvePostingRuleChange(r.id)
    await load()
  } catch (e: any) {
    actionError.value =
      e?.response?.data?.error || e?.message || 'Failed to approve change'
  }
}

onMounted(load)
</script>
