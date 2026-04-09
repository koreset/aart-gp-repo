<template>
  <v-container fluid>
    <base-card :show-actions="false">
      <template #header>
        <h3 class="mb-0">Invoices</h3>
      </template>
      <template #default>
        <!-- Quick Stats -->
        <v-row class="mb-4">
          <v-col
            v-for="card in statCards"
            :key="card.title"
            cols="12"
            sm="6"
            md="3"
          >
            <stat-card
              :title="card.title"
              :value="card.value"
              :icon="card.icon"
              :color="card.iconColor"
              :loading="loading"
            />
          </v-col>
        </v-row>

        <!-- Filter Bar -->
        <v-row class="mb-3" align="center">
          <v-col cols="12" md="2">
            <v-select
              v-model="filters.schemeId"
              label="Scheme"
              :items="inForceSchemes"
              item-title="name"
              item-value="id"
              variant="outlined"
              density="compact"
              clearable
              prepend-inner-icon="mdi-office-building-outline"
              :loading="schemesLoading"
            />
          </v-col>
          <v-col cols="12" md="2">
            <v-select
              v-model="filters.status"
              label="Status"
              :items="statusOptions"
              variant="outlined"
              density="compact"
              clearable
            />
          </v-col>
          <v-col cols="12" md="2">
            <v-text-field
              v-model="filters.from"
              label="From"
              type="date"
              variant="outlined"
              density="compact"
            />
          </v-col>
          <v-col cols="12" md="2">
            <v-text-field
              v-model="filters.to"
              label="To"
              type="date"
              variant="outlined"
              density="compact"
            />
          </v-col>
          <v-col cols="12" md="4" class="d-flex ga-2 align-center">
            <v-btn
              color="primary"
              variant="outlined"
              prepend-icon="mdi-magnify"
              @click="loadInvoices"
            >
              Search
            </v-btn>
            <v-btn
              variant="outlined"
              prepend-icon="mdi-refresh"
              @click="loadInvoices"
              >Refresh</v-btn
            >
            <v-btn
              variant="outlined"
              prepend-icon="mdi-filter-remove"
              @click="
                () => {
                  resetFilters()
                  loadInvoices()
                }
              "
            >
              Clear Filters
            </v-btn>
          </v-col>
        </v-row>

        <!-- Invoice Grid -->
        <v-row>
          <v-col cols="12">
            <v-card variant="outlined">
              <v-card-text class="pa-0">
                <ag-grid-vue
                  class="ag-theme-balham"
                  :style="{ height: gridHeight, width: '100%' }"
                  :column-defs="columnDefs"
                  :row-data="invoices"
                  :default-col-def="defaultColDef"
                  :loading="loading"
                  row-selection="single"
                  :get-row-class="getRowClass"
                  @row-clicked="onRowClicked"
                />
                <empty-state
                  v-if="!loading && invoices.length === 0"
                  title="No invoices found"
                  message="Finalize a premium schedule to generate invoices."
                  icon="mdi-receipt-text-outline"
                />
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </template>
    </base-card>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-balham.css'
import { AgGridVue } from 'ag-grid-vue3'
import PremiumManagementService from '@/renderer/api/PremiumManagementService'
import BaseCard from '@/renderer/components/BaseCard.vue'
import StatCard from '@/renderer/components/StatCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import { useGridHeight } from '@/renderer/composables/useGridHeight'
import { useFilterPersistence } from '@/renderer/composables/useFilterPersistence'
import { statusCellRenderer } from '@/renderer/utils/statusCellRenderer'
import { dateFormatter } from '@/renderer/utils/formatters'
import { useStatusBarStore } from '@/renderer/store/statusBar'

const statusBarStore = useStatusBarStore()
const router = useRouter()
const gridHeight = useGridHeight(340)
const { filters, resetFilters } = useFilterPersistence('invoices', {
  schemeId: null as number | null,
  status: null as string | null,
  from: '',
  to: ''
})

const loading = ref(false)
const invoices = ref<any[]>([])
const schemesLoading = ref(false)
const inForceSchemes = ref<{ id: number; name: string }[]>([])
const statusOptions = ['draft', 'sent', 'partial', 'paid', 'overdue']

const today = new Date().toISOString().slice(0, 10)
const currentMonth = new Date().getMonth() + 1
const currentYear = new Date().getFullYear()

const statCards = computed(() => {
  const thisMonth = invoices.value.filter(
    (i) => i.month === currentMonth && i.year === currentYear
  )
  const invoicedThisMonth = thisMonth.reduce(
    (s: number, i: any) => s + i.net_payable,
    0
  )
  const collected = thisMonth.reduce(
    (s: number, i: any) => s + i.paid_amount,
    0
  )
  const outstanding = invoices.value.reduce(
    (s: number, i: any) => s + i.balance,
    0
  )
  const overdue = invoices.value
    .filter((i: any) => i.balance > 0 && i.due_date < today)
    .reduce((s: number, i: any) => s + i.balance, 0)

  return [
    {
      title: 'Invoiced This Month',
      value: fmtCurrency(invoicedThisMonth),
      icon: 'mdi-file-document-outline',
      color: 'text-primary',
      iconColor: 'primary'
    },
    {
      title: 'Collected',
      value: fmtCurrency(collected),
      icon: 'mdi-check-circle-outline',
      color: 'text-success',
      iconColor: 'success'
    },
    {
      title: 'Outstanding',
      value: fmtCurrency(outstanding),
      icon: 'mdi-clock-outline',
      color: 'text-warning',
      iconColor: 'warning'
    },
    {
      title: 'Overdue',
      value: fmtCurrency(overdue),
      icon: 'mdi-alert-circle-outline',
      color: 'text-error',
      iconColor: 'error'
    }
  ]
})

const defaultColDef = { sortable: true, filter: true, resizable: true, flex: 1 }
const columnDefs = [
  { headerName: 'Invoice #', field: 'invoice_number', minWidth: 160 },
  { headerName: 'Scheme', field: 'scheme_name', minWidth: 180 },
  {
    headerName: 'Period',
    valueGetter: (p: any) => `${p.data.month}/${p.data.year}`
  },
  {
    headerName: 'Gross Amount',
    field: 'gross_amount',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: 'Paid',
    field: 'paid_amount',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  {
    headerName: 'Balance',
    field: 'balance',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  { headerName: 'Due Date', field: 'due_date', valueFormatter: dateFormatter },
  {
    headerName: 'Status',
    field: 'status',
    cellRenderer: (p: any) => statusCellRenderer(p.value)
  },
  {
    headerName: '',
    cellRenderer: () =>
      `<span style="cursor:pointer;color:#1976D2;font-size:12px;font-weight:600">View →</span>`,
    onCellClicked: (p: any) =>
      router.push({
        name: 'group-pricing-invoice-detail',
        params: { invoiceId: p.data.id }
      }),
    maxWidth: 80,
    sortable: false,
    filter: false
  }
]

function getRowClass(params: any) {
  if (params.data?.balance > 0 && params.data?.due_date < today) {
    return 'ag-row-overdue'
  }
  return ''
}

function fmtCurrency(val: number) {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR',
    maximumFractionDigits: 0
  }).format(val ?? 0)
}

function onRowClicked(e: any) {
  router.push({
    name: 'group-pricing-invoice-detail',
    params: { invoiceId: e.data.id }
  })
}

async function loadInvoices() {
  loading.value = true
  try {
    const params: any = {}
    if (filters.value.schemeId) params.scheme_id = filters.value.schemeId
    if (filters.value.status) params.status = filters.value.status
    if (filters.value.from) params.from = filters.value.from
    if (filters.value.to) params.to = filters.value.to
    const res = await PremiumManagementService.getInvoices(params)
    invoices.value = res.data.data ?? []
    const outstanding = invoices.value.reduce(
      (s: number, i: any) => s + i.balance,
      0
    )
    const overdue = invoices.value
      .filter((i: any) => i.balance > 0 && i.due_date < today)
      .reduce((s: number, i: any) => s + i.balance, 0)
    const overdueCount = invoices.value.filter(
      (i: any) => i.balance > 0 && i.due_date < today
    ).length
    statusBarStore.set([
      {
        icon: 'mdi-clock-outline',
        text: `Outstanding: ${fmtCurrency(outstanding)}`
      },
      {
        icon: 'mdi-alert-circle-outline',
        text: `Overdue: ${overdueCount} · ${fmtCurrency(overdue)}`,
        severity: overdue > 0 ? 'error' : 'info'
      }
    ])
  } catch (e) {
    console.error('Failed to load invoices', e)
  } finally {
    loading.value = false
  }
}

async function loadInForceSchemes() {
  schemesLoading.value = true
  try {
    const res = await PremiumManagementService.getInForceSchemes()
    inForceSchemes.value = (res.data ?? []).map((s: any) => ({
      id: s.id,
      name: s.name
    }))
  } catch (e) {
    console.error('Failed to load in-force schemes', e)
  } finally {
    schemesLoading.value = false
  }
}

onMounted(() => {
  loadInForceSchemes()
  loadInvoices()
})
onUnmounted(() => statusBarStore.clear())
</script>
