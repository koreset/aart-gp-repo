<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center">
              <span class="headline">Payment Authority Matrix</span>
              <v-spacer />
              <v-btn
                color="primary"
                prepend-icon="mdi-plus"
                size="small"
                @click="openCreate"
              >
                Add row
              </v-btn>
            </div>
          </template>

          <template #default>
            <v-alert
              type="info"
              variant="tonal"
              density="compact"
              icon="mdi-information-outline"
              class="mb-4"
            >
              <div class="text-body-2">
                Each row authorises a role to perform a payment-schedule action
                when the schedule's net payable amount falls between the min and
                max bounds. Leave the matrix empty to allow all holders of the
                relevant permission slug to act (bootstrap mode).
              </div>
            </v-alert>

            <v-data-table
              :headers="headers"
              :items="rows"
              :loading="loading"
              density="compact"
              hover
            >
              <template #[`item.action`]="{ item }">
                <span class="font-weight-medium">{{
                  ACTION_LABELS[item.action] ?? item.action
                }}</span>
              </template>

              <template #[`item.min_amount`]="{ item }">
                {{ formatCurrency(item.min_amount) }}
              </template>

              <template #[`item.max_amount`]="{ item }">
                <span v-if="item.max_amount < 0" class="text-medium-emphasis"
                  >No limit</span
                >
                <span v-else>{{ formatCurrency(item.max_amount) }}</span>
              </template>

              <template #[`item.is_active`]="{ item }">
                <v-chip
                  :color="item.is_active ? 'success' : 'grey'"
                  size="x-small"
                  variant="flat"
                >
                  {{ item.is_active ? 'Active' : 'Disabled' }}
                </v-chip>
              </template>

              <template #[`item.actions`]="{ item }">
                <v-btn
                  size="x-small"
                  variant="text"
                  icon="mdi-pencil"
                  @click="openEdit(item)"
                />
                <v-btn
                  size="x-small"
                  variant="text"
                  color="error"
                  icon="mdi-delete"
                  @click="remove(item)"
                />
              </template>
            </v-data-table>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Edit / Create dialog -->
    <v-dialog v-model="editDialog" max-width="540px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2">
          {{ editing ? 'Edit authority row' : 'Add authority row' }}
        </v-card-title>
        <v-card-text>
          <v-text-field
            v-model="form.role"
            label="Role name *"
            variant="outlined"
            density="compact"
            class="mb-2"
            hint="Must match a GP role name exactly (e.g. 'Head of Claims')"
            persistent-hint
          />
          <v-select
            v-model="form.action"
            :items="ACTION_OPTIONS"
            label="Action *"
            variant="outlined"
            density="compact"
            class="mb-2"
          />
          <v-row dense>
            <v-col cols="6">
              <v-text-field
                v-model.number="form.min_amount"
                type="number"
                label="Min amount"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model.number="form.max_amount"
                type="number"
                label="Max amount (-1 = no limit)"
                variant="outlined"
                density="compact"
              />
            </v-col>
          </v-row>
          <v-switch
            v-model="form.is_active"
            color="success"
            label="Active"
            density="compact"
            hide-details
          />
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="editDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="saving"
            :disabled="!form.role || !form.action"
            @click="save"
          >
            Save
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="3500">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

interface AuthorityRow {
  id: number
  role: string
  action: string
  min_amount: number
  max_amount: number
  is_active: boolean
}

const ACTION_LABELS: Record<string, string> = {
  signoff_schedule: 'Head of Claims sign-off',
  finance_review: 'Finance review',
  authorise_first: '1st finance authorisation',
  authorise_second: '2nd finance authorisation',
  archive_schedule: 'Archive schedule'
}

const ACTION_OPTIONS = Object.entries(ACTION_LABELS).map(([value, title]) => ({
  value,
  title
}))

const headers = [
  { title: 'Role', key: 'role', sortable: true },
  { title: 'Action', key: 'action', sortable: true },
  {
    title: 'Min amount',
    key: 'min_amount',
    sortable: true,
    align: 'end' as const
  },
  {
    title: 'Max amount',
    key: 'max_amount',
    sortable: true,
    align: 'end' as const
  },
  { title: 'Active', key: 'is_active', sortable: false },
  { title: '', key: 'actions', sortable: false }
]

const rows = ref<AuthorityRow[]>([])
const loading = ref(false)
const saving = ref(false)
const editDialog = ref(false)
const editing = ref<AuthorityRow | null>(null)

const form = reactive({
  role: '',
  action: 'signoff_schedule',
  min_amount: 0,
  max_amount: -1,
  is_active: true
})

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

function notify(message: string, color = 'success') {
  snackbarMessage.value = message
  snackbarColor.value = color
  snackbar.value = true
}

function formatCurrency(val: number) {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR'
  }).format(val ?? 0)
}

function unwrap(res: any) {
  const body = res?.data
  if (body && typeof body === 'object' && 'success' in body && 'data' in body) {
    return body.data
  }
  return body
}

async function load() {
  loading.value = true
  try {
    const res = await GroupPricingService.listAuthorityMatrix()
    rows.value = unwrap(res) ?? []
  } catch (e: any) {
    notify('Failed to load authority matrix', 'error')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = null
  Object.assign(form, {
    role: '',
    action: 'signoff_schedule',
    min_amount: 0,
    max_amount: -1,
    is_active: true
  })
  editDialog.value = true
}

function openEdit(row: AuthorityRow) {
  editing.value = row
  Object.assign(form, {
    role: row.role,
    action: row.action,
    min_amount: row.min_amount,
    max_amount: row.max_amount,
    is_active: row.is_active
  })
  editDialog.value = true
}

async function save() {
  saving.value = true
  try {
    if (editing.value) {
      await GroupPricingService.updateAuthorityMatrixRow(editing.value.id, form)
      notify('Authority row updated.')
    } else {
      await GroupPricingService.createAuthorityMatrixRow(form)
      notify('Authority row created.')
    }
    editDialog.value = false
    await load()
  } catch (e: any) {
    notify('Failed to save authority row', 'error')
  } finally {
    saving.value = false
  }
}

async function remove(row: AuthorityRow) {
  if (!confirm(`Delete authority row for ${row.role} → ${row.action}?`)) return
  try {
    await GroupPricingService.deleteAuthorityMatrixRow(row.id)
    notify('Authority row deleted.')
    await load()
  } catch (e: any) {
    notify('Failed to delete authority row', 'error')
  }
}

onMounted(load)
</script>
