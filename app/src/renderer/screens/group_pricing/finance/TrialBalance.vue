<template>
  <v-container>
    <base-card :show-actions="false">
      <template #header>
        <div class="d-flex justify-space-between align-center flex-wrap">
          <span class="headline">Trial Balance</span>
          <v-select
            v-model="periodId"
            :items="periods"
            item-title="name"
            item-value="id"
            label="Period"
            density="compact"
            variant="outlined"
            style="max-width: 220px"
            clearable
            @update:model-value="load"
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
          @click:row="onRowClick"
        >
          <template #[`item.total_debit`]="{ value }">{{ format(value) }}</template>
          <template #[`item.total_credit`]="{ value }">{{ format(value) }}</template>
          <template #[`item.net_balance`]="{ value, item }">
            <span :class="netColour(value, item.normal_balance)">
              {{ format(value) }}
            </span>
          </template>
        </v-data-table>

        <div class="d-flex justify-end mt-3 text-subtitle-2">
          <span class="mr-6">Total Debit: <strong>{{ format(totalDr) }}</strong></span>
          <span class="mr-6">Total Credit: <strong>{{ format(totalCr) }}</strong></span>
          <v-chip
            :color="Math.abs(totalDr - totalCr) < 0.005 ? 'success' : 'warning'"
            size="small"
            variant="tonal"
            >{{
              Math.abs(totalDr - totalCr) < 0.005 ? 'In balance' : 'Out of balance'
            }}</v-chip
          >
        </div>
      </template>
    </base-card>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import GLService, {
  AccountingPeriod,
  TrialBalanceRow
} from '@/renderer/api/GeneralLedgerService'

const router = useRouter()

const headers = [
  { title: 'Code', key: 'account_code' },
  { title: 'Account', key: 'account_name' },
  { title: 'Type', key: 'account_type' },
  { title: 'Debit', key: 'total_debit', align: 'end' as const },
  { title: 'Credit', key: 'total_credit', align: 'end' as const },
  { title: 'Net', key: 'net_balance', align: 'end' as const }
]

const periods = ref<AccountingPeriod[]>([])
const periodId = ref<number | undefined>(undefined)
const rows = ref<TrialBalanceRow[]>([])
const loading = ref(false)

const totalDr = computed(() =>
  rows.value.reduce((sum, r) => sum + r.total_debit, 0)
)
const totalCr = computed(() =>
  rows.value.reduce((sum, r) => sum + r.total_credit, 0)
)

const format = (n: number) =>
  new Intl.NumberFormat(undefined, { minimumFractionDigits: 2 }).format(n || 0)

const netColour = (n: number, normal: string) => {
  if (Math.abs(n) < 0.005) return 'text-grey'
  return n > 0 ? `text-${normal === 'debit' ? 'blue' : 'green'}` : 'text-red'
}

const load = async () => {
  loading.value = true
  try {
    rows.value = await GLService.getTrialBalance(periodId.value)
  } finally {
    loading.value = false
  }
}

const onRowClick = (_evt: unknown, row: { item: TrialBalanceRow }) => {
  router.push({
    name: 'group-pricing-gl-account-ledger',
    params: { id: row.item.account_id }
  })
}

onMounted(async () => {
  periods.value = await GLService.listPeriods()
  await load()
})
</script>
