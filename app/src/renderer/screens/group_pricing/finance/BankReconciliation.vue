<template>
  <v-container fluid>
    <base-card :show-actions="false">
      <template #header>
        <div class="d-flex justify-space-between align-center flex-wrap gap-3">
          <span class="headline">Bank Reconciliation</span>
          <div class="d-flex align-center gap-2">
            <v-select
              v-model="selectedBankId"
              :items="bankAccounts"
              item-title="name"
              item-value="id"
              label="Bank account"
              density="compact"
              variant="outlined"
              style="min-width: 260px"
              @update:model-value="load"
            />
            <v-btn
              v-if="hasPermission('gl:bank_rec') && selectedBankId"
              prepend-icon="mdi-upload"
              variant="outlined"
              @click="importDialog = true"
            >
              Import Statement
            </v-btn>
          </div>
        </div>
      </template>
      <template #default>
        <v-alert color="info" variant="tonal" density="compact" class="mb-3" icon="mdi-shield-check-outline">
          Matches and ignores are recorded immediately, but the statement
          line is not considered reconciled until a different user signs off
          on the action under <strong>Pending Review</strong>.
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

        <v-tabs v-if="selectedBankId" v-model="activeTab" density="compact" color="primary">
          <v-tab value="match">Match</v-tab>
          <v-tab value="review">
            Pending Review
            <v-chip size="x-small" class="ml-2">{{ pendingReview.length }}</v-chip>
          </v-tab>
        </v-tabs>

        <v-window v-if="selectedBankId" v-model="activeTab" class="mt-3">
          <v-window-item value="match">
            <v-row>
              <v-col cols="12" md="6">
                <div class="text-subtitle-2 mb-2">
                  Statement lines
                  <v-chip size="x-small" class="ml-2">{{ unmatched.length }}</v-chip>
                </div>
                <v-data-table
                  v-model="selectedStmt"
                  :headers="stmtHeaders"
                  :items="unmatched"
                  :loading="loading"
                  density="compact"
                  items-per-page="50"
                  show-select
                  return-object
                  item-value="id"
                  select-strategy="single"
                >
                  <template #[`item.statement_date`]="{ value }">
                    {{ new Date(value).toLocaleDateString() }}
                  </template>
                  <template #[`item.amount`]="{ value }">
                    {{ format(value) }}
                  </template>
                </v-data-table>
              </v-col>
              <v-col cols="12" md="6">
                <div class="text-subtitle-2 mb-2">
                  Posted bank lines
                  <v-chip size="x-small" class="ml-2">{{ bankLedger.length }}</v-chip>
                </div>
                <v-data-table
                  v-model="selectedLedger"
                  :headers="ledgerHeaders"
                  :items="bankLedger"
                  density="compact"
                  items-per-page="50"
                  show-select
                  return-object
                  item-value="entry_id"
                  select-strategy="single"
                >
                  <template #[`item.posted_at`]="{ value }">
                    {{ new Date(value).toLocaleDateString() }}
                  </template>
                  <template #[`item.debit`]="{ value }">{{ format(value) }}</template>
                  <template #[`item.credit`]="{ value }">{{ format(value) }}</template>
                </v-data-table>
              </v-col>
              <v-col cols="12" class="d-flex justify-end gap-2">
                <v-btn
                  v-if="hasPermission('gl:bank_rec')"
                  variant="outlined"
                  :disabled="selectedStmt.length === 0"
                  @click="ignore"
                >
                  Ignore Statement Line
                </v-btn>
                <v-btn
                  color="primary"
                  :disabled="selectedStmt.length === 0 || selectedLedger.length === 0"
                  @click="match"
                >
                  Match Selected
                </v-btn>
              </v-col>
            </v-row>
          </v-window-item>

          <v-window-item value="review">
            <v-data-table
              :headers="reviewHeaders"
              :items="pendingReview"
              :loading="loading"
              density="compact"
              items-per-page="50"
            >
              <template #[`item.statement_date`]="{ value }">
                {{ new Date(value).toLocaleDateString() }}
              </template>
              <template #[`item.amount`]="{ value }">{{ format(value) }}</template>
              <template #[`item.match_status`]="{ value }">
                <v-chip
                  :color="value === 'matched' ? 'primary' : 'grey'"
                  size="x-small"
                  variant="tonal"
                  class="text-uppercase"
                >
                  {{ value }}
                </v-chip>
              </template>
              <template #[`item.matched_by`]="{ item }">
                <div class="text-body-2">
                  {{ item.matched_by }}
                </div>
                <div v-if="item.matched_at" class="text-caption text-medium-emphasis">
                  {{ new Date(item.matched_at).toLocaleString() }}
                </div>
              </template>
              <template #[`item.actions`]="{ item }">
                <div class="d-flex ga-1 justify-end">
                  <template
                    v-if="
                      hasPermission('gl:review_reconciliation') &&
                      !sameUser(item.matched_by, currentUserName)
                    "
                  >
                    <v-btn
                      size="x-small"
                      color="success"
                      variant="tonal"
                      @click="review(item, 'reviewed')"
                    >
                      Sign off
                    </v-btn>
                    <v-btn
                      size="x-small"
                      color="warning"
                      variant="tonal"
                      @click="review(item, 'rejected')"
                    >
                      Reject
                    </v-btn>
                  </template>
                  <v-chip
                    v-else-if="sameUser(item.matched_by, currentUserName)"
                    size="x-small"
                    variant="tonal"
                    color="warning"
                  >
                    Awaiting review
                  </v-chip>
                </div>
              </template>
            </v-data-table>
          </v-window-item>
        </v-window>

        <v-alert v-if="!selectedBankId" color="info" variant="tonal" class="my-4">
          Select a bank account to start reconciling.
        </v-alert>
      </template>
    </base-card>

    <v-dialog v-model="importDialog" max-width="600">
      <v-card>
        <v-card-title>Paste statement rows (CSV)</v-card-title>
        <v-card-text>
          <v-alert color="info" variant="tonal" density="compact" class="mb-3">
            Format: <code>YYYY-MM-DD, value-date, description, amount, reference</code>.
            value-date and reference are optional.
          </v-alert>
          <v-textarea
            v-model="csvText"
            rows="10"
            variant="outlined"
            density="compact"
            placeholder="2026-05-15,, Salary payment, -125000, REF123"
          />
          <v-alert
            v-if="importError"
            color="error"
            variant="tonal"
            class="mt-2"
          >
            {{ importError }}
          </v-alert>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="importDialog = false">Cancel</v-btn>
          <v-btn color="primary" :loading="importing" @click="doImport">
            Import
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import GLService, {
  BankAccount,
  BankStatementLine,
  LedgerRow,
  StatementImportRow
} from '@/renderer/api/GeneralLedgerService'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

const { hasPermission } = usePermissionCheck()

const currentUserName = computed<string>(() => {
  const u: any = window.mainApi?.sendSync('msgGetAuthenticatedUser')
  return u?.full_name || u?.username || ''
})

const stmtHeaders = [
  { title: 'Date', key: 'statement_date' },
  { title: 'Description', key: 'description' },
  { title: 'Amount', key: 'amount', align: 'end' as const },
  { title: 'Reference', key: 'reference' }
]

const ledgerHeaders = [
  { title: 'Date', key: 'posted_at' },
  { title: 'Entry', key: 'entry_number' },
  { title: 'Description', key: 'description' },
  { title: 'Debit', key: 'debit', align: 'end' as const },
  { title: 'Credit', key: 'credit', align: 'end' as const }
]

const reviewHeaders = [
  { title: 'Date', key: 'statement_date' },
  { title: 'Description', key: 'description' },
  { title: 'Amount', key: 'amount', align: 'end' as const },
  { title: 'Reference', key: 'reference' },
  { title: 'Action', key: 'match_status' },
  { title: 'By', key: 'matched_by' },
  { title: '', key: 'actions', sortable: false, align: 'end' as const }
]

const bankAccounts = ref<BankAccount[]>([])
const selectedBankId = ref<number | null>(null)
const allLines = ref<BankStatementLine[]>([])
const bankLedger = ref<(LedgerRow & { journal_line_id?: number })[]>([])
const selectedStmt = ref<BankStatementLine[]>([])
const selectedLedger = ref<LedgerRow[]>([])
const loading = ref(false)
const importDialog = ref(false)
const importing = ref(false)
const csvText = ref('')
const importError = ref('')
const actionError = ref('')
const activeTab = ref<'match' | 'review'>('match')

const unmatched = computed(() =>
  allLines.value.filter((l) => l.match_status === 'unmatched')
)
const pendingReview = computed(() =>
  allLines.value.filter((l) => l.review_status === 'pending_review')
)

const sameUser = (a?: string, b?: string) =>
  Boolean(a && b && a.trim().toLowerCase() === b.trim().toLowerCase())

const format = (n: number) =>
  new Intl.NumberFormat(undefined, { minimumFractionDigits: 2 }).format(n || 0)

const load = async () => {
  if (!selectedBankId.value) return
  loading.value = true
  try {
    const bank = bankAccounts.value.find((x) => x.id === selectedBankId.value)
    const [s, b] = await Promise.all([
      GLService.listStatementLines(selectedBankId.value),
      bank
        ? GLService.getAccountLedger(bank.gl_account_id)
        : Promise.resolve([] as LedgerRow[])
    ])
    allLines.value = s
    bankLedger.value = b
    selectedStmt.value = []
    selectedLedger.value = []
  } finally {
    loading.value = false
  }
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

const match = async () => {
  if (!selectedStmt.value.length || !selectedLedger.value.length) return
  await runAction(async () => {
    const entry = selectedLedger.value[0]
    const full = await GLService.getJournalEntry(entry.entry_id)
    const bank = bankAccounts.value.find((x) => x.id === selectedBankId.value)!
    const line = (full.lines || []).find((l) => l.account_id === bank.gl_account_id)
    if (!line?.id) throw new Error('Selected entry has no line for this bank account')
    await GLService.matchStatementLine(selectedStmt.value[0].id, line.id)
  })
}

const ignore = () => {
  if (!selectedStmt.value.length) return
  runAction(() => GLService.ignoreStatementLine(selectedStmt.value[0].id))
}

const review = (line: BankStatementLine, outcome: 'reviewed' | 'rejected') => {
  runAction(() => GLService.reviewStatementLine(line.id, outcome))
}

const doImport = async () => {
  importError.value = ''
  if (!selectedBankId.value) return
  const rows: StatementImportRow[] = []
  for (const raw of csvText.value.split(/\r?\n/)) {
    const line = raw.trim()
    if (!line) continue
    const parts = line.split(',').map((p) => p.trim())
    if (parts.length < 3) {
      importError.value = `Need at least date, description, amount: "${line}"`
      return
    }
    rows.push({
      statement_date: parts[0],
      value_date: parts[1] || undefined,
      description: parts[2] || '',
      amount: Number(parts[3] ?? parts[2]),
      reference: parts[4] || ''
    })
  }
  if (!rows.length) {
    importError.value = 'No rows to import'
    return
  }
  importing.value = true
  try {
    await GLService.importStatement(selectedBankId.value, rows)
    importDialog.value = false
    csvText.value = ''
    await load()
  } catch (e: any) {
    importError.value = e?.response?.data?.error || 'Import failed'
  } finally {
    importing.value = false
  }
}

onMounted(async () => {
  const data = await GLService.listBankAccounts()
  bankAccounts.value = data
  if (data.length === 1) {
    selectedBankId.value = data[0].id || null
    await load()
  }
})
</script>
