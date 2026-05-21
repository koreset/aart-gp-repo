<template>
  <v-container>
    <base-card :show-actions="false">
      <template #header>
        <div class="d-flex justify-space-between align-center flex-wrap">
          <span class="headline">Bank Accounts</span>
          <v-btn
            v-if="hasPermission('gl:manage_accounts')"
            color="primary"
            prepend-icon="mdi-plus"
            @click="openCreate"
          >
            New Bank Account
          </v-btn>
        </div>
      </template>
      <template #default>
        <v-alert color="info" variant="tonal" density="compact" class="mb-3" icon="mdi-shield-check-outline">
          Bank account changes (number, GL link, name) are payment-fraud risks
          and are staged for a second user to approve. The live record keeps
          its current values until the approver releases the change.
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
          :items="rows"
          :loading="loading"
          density="compact"
          items-per-page="50"
        >
          <template #[`item.gl_account_id`]="{ value }">
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
                v-if="hasPermission('gl:manage_accounts')"
                icon="mdi-pencil"
                variant="text"
                size="small"
                @click="openEdit(item)"
              />
              <v-btn
                v-if="hasPermission('gl:manage_accounts') && item.is_active"
                icon="mdi-pause"
                variant="text"
                size="small"
                @click="deactivate(item)"
              />
            </template>
          </template>
        </v-data-table>
      </template>
    </base-card>

    <v-dialog v-model="dialog" max-width="640">
      <v-card>
        <v-card-title>
          {{ form.id ? 'Request Update' : 'Request New' }} Bank Account
        </v-card-title>
        <v-card-text>
          <v-row dense>
            <v-col cols="12" md="4">
              <v-text-field
                v-model="form.code"
                label="Code"
                :disabled="!!form.id"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12" md="8">
              <v-text-field
                v-model="form.name"
                label="Name"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="form.bank_name"
                label="Bank"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="form.account_number"
                label="Account number"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12" md="8">
              <v-autocomplete
                v-model="form.gl_account_id"
                :items="bankGLAccounts"
                :item-title="accountLabel"
                item-value="id"
                label="GL account (asset)"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12" md="4">
              <v-text-field
                v-model="form.currency"
                label="Currency"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12">
              <v-switch
                v-model="form.is_active"
                label="Active"
                color="success"
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
  BankAccount,
  GLAccount
} from '@/renderer/api/GeneralLedgerService'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

const { hasPermission } = usePermissionCheck()

const currentUserName = computed<string>(() => {
  const u: any = window.mainApi?.sendSync('msgGetAuthenticatedUser')
  return u?.full_name || u?.username || ''
})

const headers = [
  { title: 'Code', key: 'code' },
  { title: 'Name', key: 'name' },
  { title: 'Bank', key: 'bank_name' },
  { title: 'Account #', key: 'account_number' },
  { title: 'GL Account', key: 'gl_account_id' },
  { title: 'Currency', key: 'currency' },
  { title: 'Active', key: 'is_active' },
  { title: 'Pending', key: 'approval_status', sortable: false },
  { title: '', key: 'actions', sortable: false, align: 'end' as const }
]

const rows = ref<BankAccount[]>([])
const accounts = ref<GLAccount[]>([])
const loading = ref(false)
const dialog = ref(false)
const saving = ref(false)
const actionError = ref('')
const form = reactive<BankAccount>({
  code: '',
  name: '',
  bank_name: '',
  account_number: '',
  gl_account_id: 0,
  currency: 'ZAR',
  is_active: true
})

const sameUser = (a?: string, b?: string) =>
  Boolean(a && b && a.trim().toLowerCase() === b.trim().toLowerCase())

const bankGLAccounts = computed(() =>
  accounts.value.filter((a) => a.account_type === 'asset' && a.is_active)
)

const accountLabel = (a: GLAccount | number) => {
  const acc =
    typeof a === 'number' ? accounts.value.find((x) => x.id === a) : a
  return acc ? `${acc.code} — ${acc.name}` : ''
}

const canApprove = (b: BankAccount) =>
  b.approval_status &&
  b.approval_status !== 'active' &&
  hasPermission('gl:approve_bank_account') &&
  !sameUser(b.pending_requested_by, currentUserName.value)

const load = async () => {
  loading.value = true
  try {
    const [b, a] = await Promise.all([
      GLService.listBankAccounts(),
      GLService.listAccounts()
    ])
    rows.value = b
    accounts.value = a
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  Object.assign(form, {
    id: undefined,
    code: '',
    name: '',
    bank_name: '',
    account_number: '',
    gl_account_id: 0,
    currency: 'ZAR',
    is_active: true
  })
  dialog.value = true
}

const openEdit = (b: BankAccount) => {
  Object.assign(form, b)
  dialog.value = true
}

const save = async () => {
  saving.value = true
  actionError.value = ''
  try {
    if (form.id) {
      await GLService.requestUpdateBankAccount(form.id, form)
    } else {
      await GLService.requestCreateBankAccount(form)
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

const deactivate = async (b: BankAccount) => {
  if (!b.id) return
  if (
    !confirm(
      `Request deactivation of ${b.code} — ${b.name}? A different user must approve.`
    )
  )
    return
  try {
    await GLService.requestDeactivateBankAccount(b.id)
    await load()
  } catch (e: any) {
    actionError.value =
      e?.response?.data?.error || e?.message || 'Failed to request deactivation'
  }
}

const approve = async (b: BankAccount) => {
  if (!b.id) return
  try {
    await GLService.approveBankAccountChange(b.id)
    await load()
  } catch (e: any) {
    actionError.value =
      e?.response?.data?.error || e?.message || 'Failed to approve change'
  }
}

onMounted(load)
</script>
