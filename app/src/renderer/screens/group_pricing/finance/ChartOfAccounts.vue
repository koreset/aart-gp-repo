<template>
  <v-container>
    <base-card :show-actions="false">
      <template #header>
        <div class="d-flex justify-space-between align-center flex-wrap">
          <span class="headline">Chart of Accounts</span>
          <v-btn
            v-if="hasPermission('gl:manage_accounts')"
            color="primary"
            prepend-icon="mdi-plus"
            @click="openCreate"
          >
            New Account
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
        <v-alert color="info" variant="tonal" density="compact" class="mb-3" icon="mdi-shield-check-outline">
          Every change to the chart of accounts is staged for a second user to approve.
        </v-alert>

        <v-data-table
          :headers="headers"
          :items="accounts"
          :loading="loading"
          density="compact"
          items-per-page="50"
        >
          <template #[`item.account_type`]="{ value }">
            <v-chip size="small" :color="typeColour(value)" variant="tonal">
              {{ value }}
            </v-chip>
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
          {{ form.id ? 'Request Update' : 'Request New' }} Account
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
              <v-select
                v-model="form.account_type"
                :items="['asset', 'liability', 'equity', 'income', 'expense']"
                label="Type"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-select
                v-model="form.normal_balance"
                :items="['debit', 'credit']"
                label="Normal Balance"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="form.description"
                label="Description"
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
import GLService, { GLAccount } from '@/renderer/api/GeneralLedgerService'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

const { hasPermission } = usePermissionCheck()

const currentUserName = computed<string>(() => {
  const u: any = window.mainApi?.sendSync('msgGetAuthenticatedUser')
  return u?.full_name || u?.username || ''
})

const headers = [
  { title: 'Code', key: 'code', sortable: true },
  { title: 'Name', key: 'name', sortable: true },
  { title: 'Type', key: 'account_type' },
  { title: 'Normal', key: 'normal_balance' },
  { title: 'Active', key: 'is_active' },
  { title: 'Pending', key: 'approval_status', sortable: false },
  { title: '', key: 'actions', sortable: false, align: 'end' as const }
]

const accounts = ref<GLAccount[]>([])
const loading = ref(false)
const dialog = ref(false)
const saving = ref(false)
const actionError = ref('')
const form = reactive<GLAccount>({
  code: '',
  name: '',
  account_type: 'asset',
  normal_balance: 'debit',
  is_active: true,
  description: ''
})

const sameUser = (a?: string, b?: string) =>
  Boolean(a && b && a.trim().toLowerCase() === b.trim().toLowerCase())

const typeColour = (t: string) => {
  return (
    {
      asset: 'blue',
      liability: 'red',
      equity: 'purple',
      income: 'green',
      expense: 'orange'
    } as Record<string, string>
  )[t] || 'grey'
}

const canApprove = (a: GLAccount) =>
  a.approval_status &&
  a.approval_status !== 'active' &&
  hasPermission('gl:approve_account') &&
  !sameUser(a.pending_requested_by, currentUserName.value)

const load = async () => {
  loading.value = true
  try {
    accounts.value = await GLService.listAccounts()
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  Object.assign(form, {
    id: undefined,
    code: '',
    name: '',
    account_type: 'asset',
    normal_balance: 'debit',
    is_active: true,
    description: ''
  })
  dialog.value = true
}

const openEdit = (acc: GLAccount) => {
  Object.assign(form, acc)
  dialog.value = true
}

const save = async () => {
  saving.value = true
  actionError.value = ''
  try {
    if (form.id) {
      await GLService.requestUpdateAccount(form.id, form)
    } else {
      await GLService.requestCreateAccount(form)
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

const deactivate = async (acc: GLAccount) => {
  if (!acc.id) return
  if (
    !confirm(
      `Request deactivation of ${acc.code} — ${acc.name}? A different user must approve.`
    )
  )
    return
  try {
    await GLService.requestDeactivateAccount(acc.id)
    await load()
  } catch (e: any) {
    actionError.value =
      e?.response?.data?.error || e?.message || 'Failed to request deactivation'
  }
}

const approve = async (acc: GLAccount) => {
  if (!acc.id) return
  try {
    await GLService.approveAccountChange(acc.id)
    await load()
  } catch (e: any) {
    actionError.value =
      e?.response?.data?.error || e?.message || 'Failed to approve change'
  }
}

onMounted(load)
</script>
