<template>
  <v-container fluid>
    <base-card :show-actions="false">
      <template #header>
        <div class="d-flex justify-space-between align-center flex-wrap">
          <span class="headline">GL & Cash Audit Log</span>
          <v-btn
            variant="text"
            prepend-icon="mdi-refresh"
            @click="load"
          >
            Refresh
          </v-btn>
        </div>
      </template>
      <template #default>
        <v-row dense class="mb-2">
          <v-col cols="12" md="3">
            <v-select
              v-model="filters.event_type"
              :items="eventTypeOptions"
              label="Event"
              density="compact"
              variant="outlined"
              clearable
              @update:model-value="load"
            />
          </v-col>
          <v-col cols="12" md="3">
            <v-select
              v-model="filters.object_type"
              :items="objectTypeOptions"
              label="Object"
              density="compact"
              variant="outlined"
              clearable
              @update:model-value="load"
            />
          </v-col>
          <v-col cols="12" md="2">
            <v-select
              v-model="filters.changed_by"
              :items="userOptions"
              label="Changed by"
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
          :items="rows"
          :loading="loading"
          density="compact"
          items-per-page="100"
        >
          <template #[`item.changed_at`]="{ value }">
            {{ new Date(value).toLocaleString() }}
          </template>
          <template #[`item.event_type`]="{ value }">
            <v-chip size="x-small" :color="eventColour(value)" variant="tonal">
              {{ value }}
            </v-chip>
          </template>
          <template #[`item.details`]="{ value }">
            <code class="text-caption d-inline-block" style="max-width: 480px; overflow: hidden; text-overflow: ellipsis; vertical-align: middle">
              {{ value || '' }}
            </code>
          </template>
        </v-data-table>
      </template>
    </base-card>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import GLService, { GLAuditLogEntry } from '@/renderer/api/GeneralLedgerService'

const headers = [
  { title: 'When', key: 'changed_at' },
  { title: 'Event', key: 'event_type' },
  { title: 'Object', key: 'object_type' },
  { title: 'Reference', key: 'object_name' },
  { title: 'Changed by', key: 'changed_by' },
  { title: 'Details', key: 'details', sortable: false }
]

const eventTypeOptions = [
  'journal_drafted',
  'journal_draft_updated',
  'journal_draft_discarded',
  'journal_submitted',
  'journal_approved',
  'journal_posted',
  'journal_reversal_requested',
  'journal_reversal_approved',
  'period_opened',
  'period_close_requested',
  'period_closed',
  'gl_account_change_requested',
  'gl_account_change_approved',
  'posting_rule_change_requested',
  'posting_rule_change_approved',
  'bank_account_change_requested',
  'bank_account_change_approved',
  'bank_statement_imported',
  'statement_line_matched',
  'statement_line_ignored',
  'statement_line_match_reviewed',
  'statement_line_review_rejected'
]

const objectTypeOptions = [
  'journal_entry',
  'accounting_period',
  'gl_account',
  'posting_rule',
  'bank_account',
  'bank_statement_line'
]

const rows = ref<GLAuditLogEntry[]>([])
const knownUsers = ref<string[]>([])
const loading = ref(false)
const filters = reactive({
  event_type: undefined as string | undefined,
  object_type: undefined as string | undefined,
  changed_by: undefined as string | undefined,
  from: '',
  to: ''
})

// Pool of users that have appeared in any audit row this session — used to
// populate the "Changed by" dropdown. We keep the union (rather than just the
// current `rows`) so a user filtering by one person doesn't lose the others.
const userOptions = computed(() => knownUsers.value)

const refreshKnownUsers = (latest: GLAuditLogEntry[]) => {
  const pool = new Set(knownUsers.value)
  for (const r of latest) {
    if (r.changed_by) pool.add(r.changed_by)
  }
  knownUsers.value = Array.from(pool).sort((a, b) =>
    a.localeCompare(b, undefined, { sensitivity: 'base' })
  )
}

const eventColour = (e: string): string => {
  if (e.endsWith('_approved') || e.endsWith('_reviewed') || e === 'journal_posted') return 'success'
  if (e.endsWith('_requested') || e.endsWith('_submitted') || e === 'journal_drafted') return 'info'
  if (e.endsWith('_rejected') || e.endsWith('_discarded') || e === 'journal_reversal_requested') return 'warning'
  return 'default'
}

const load = async () => {
  loading.value = true
  try {
    const params: any = {}
    if (filters.event_type) params.event_type = filters.event_type
    if (filters.object_type) params.object_type = filters.object_type
    if (filters.changed_by) params.changed_by = filters.changed_by
    if (filters.from) params.from = filters.from
    if (filters.to) params.to = filters.to
    rows.value = await GLService.listAuditLog(params)
    refreshKnownUsers(rows.value)
  } finally {
    loading.value = false
  }
}

// Seed the user dropdown from the dedicated endpoint (union of users who
// have acted + all system users), so the filter is useful even before any
// audit rows exist. Falls back gracefully if the endpoint fails.
onMounted(async () => {
  try {
    const users = await GLService.listAuditLogUsers()
    if (users.length) {
      knownUsers.value = users
    }
  } catch {
    // non-fatal — load() will still populate from the visible rows
  }
  await load()
})
</script>
