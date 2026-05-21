<template>
  <v-container>
    <base-card :show-actions="false">
      <template #header>
        <div class="d-flex align-center flex-wrap gap-3">
          <v-btn
            icon="mdi-arrow-left"
            variant="text"
            class="mr-2"
            @click="$router.back()"
          />
          <span class="headline">
            {{ account ? `${account.code} — ${account.name}` : 'General Ledger' }}
          </span>
          <v-spacer />
          <v-text-field
            v-model="from"
            type="date"
            label="From"
            density="compact"
            variant="outlined"
            style="max-width: 180px"
            @change="load"
          />
          <v-text-field
            v-model="to"
            type="date"
            label="To"
            density="compact"
            variant="outlined"
            style="max-width: 180px"
            @change="load"
          />
        </div>
      </template>
      <template #default>
        <v-data-table
          :headers="headers"
          :items="rows"
          :loading="loading"
          density="compact"
          items-per-page="100"
          @click:row="openJournal"
        >
          <template #[`item.posted_at`]="{ value }">{{
            new Date(value).toLocaleDateString()
          }}</template>
          <template #[`item.debit`]="{ value }">{{ format(value) }}</template>
          <template #[`item.credit`]="{ value }">{{ format(value) }}</template>
          <template #[`item.running_balance`]="{ value }">
            <strong>{{ format(value) }}</strong>
          </template>
        </v-data-table>
      </template>
    </base-card>
  </v-container>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import GLService, {
  GLAccount,
  LedgerRow
} from '@/renderer/api/GeneralLedgerService'

const route = useRoute()
const router = useRouter()

const headers = [
  { title: 'Date', key: 'posted_at' },
  { title: 'Entry', key: 'entry_number' },
  { title: 'Source', key: 'source_type' },
  { title: 'Description', key: 'description' },
  { title: 'Debit', key: 'debit', align: 'end' as const },
  { title: 'Credit', key: 'credit', align: 'end' as const },
  { title: 'Running', key: 'running_balance', align: 'end' as const }
]

const account = ref<GLAccount | null>(null)
const rows = ref<LedgerRow[]>([])
const loading = ref(false)
const from = ref('')
const to = ref('')

const format = (n: number) =>
  new Intl.NumberFormat(undefined, { minimumFractionDigits: 2 }).format(n || 0)

const load = async () => {
  const id = Number(route.params.id)
  loading.value = true
  try {
    const params: { from?: string; to?: string } = {}
    if (from.value) params.from = from.value
    if (to.value) params.to = to.value
    const [a, l] = await Promise.all([
      GLService.getAccount(id),
      GLService.getAccountLedger(id, params)
    ])
    account.value = a
    rows.value = l
  } finally {
    loading.value = false
  }
}

const openJournal = (_evt: unknown, row: { item: LedgerRow }) => {
  router.push({
    name: 'group-pricing-gl-journal-detail',
    params: { id: row.item.entry_id }
  })
}

onMounted(load)
</script>
