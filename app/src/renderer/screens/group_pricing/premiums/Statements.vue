<template>
  <v-container fluid>
    <base-card :show-actions="false">
      <template #header>
        <h3 class="mb-0">Statements</h3>
      </template>
      <template #default>
        <v-tabs v-model="activeTab" color="primary" class="mb-4">
          <v-tab value="employer">Employer Statements</v-tab>
          <v-tab value="broker">Broker Commission Statements</v-tab>
        </v-tabs>

        <v-window v-model="activeTab">
          <!-- ── Employer Statements Tab ─────────────────────────────── -->
          <v-window-item value="employer">
            <v-row class="mb-4" align="center">
              <v-col cols="12" md="3">
                <v-select
                  v-model="employerFilters.schemeId"
                  label="Scheme *"
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
                <v-text-field
                  v-model="employerFilters.from"
                  label="From"
                  type="date"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" md="2">
                <v-text-field
                  v-model="employerFilters.to"
                  label="To"
                  type="date"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" md="5" class="d-flex align-center ga-2">
                <v-btn
                  color="primary"
                  prepend-icon="mdi-file-find"
                  :loading="loadingEmployer"
                  :disabled="
                    !employerFilters.schemeId ||
                    !employerFilters.from ||
                    !employerFilters.to
                  "
                  @click="loadEmployerStatement"
                >
                  Generate
                </v-btn>
                <v-btn
                  v-if="employerStatement"
                  variant="outlined"
                  prepend-icon="mdi-download"
                  @click="downloadEmployerCSV"
                >
                  Download CSV
                </v-btn>
                <v-btn
                  v-if="employerStatement"
                  variant="outlined"
                  prepend-icon="mdi-file-pdf-box"
                  @click="downloadEmployerPDF"
                >
                  Download PDF
                </v-btn>
              </v-col>
            </v-row>

            <!-- Statement Header -->
            <template v-if="employerStatement">
              <v-row class="mb-2">
                <v-col cols="12">
                  <v-card variant="tonal" color="primary" class="pa-3">
                    <div class="text-h6 mb-1">{{
                      employerStatement.scheme_name
                    }}</div>
                    <div class="text-body-2"
                      >Period: {{ employerStatement.period }}</div
                    >
                    <v-row class="mt-2">
                      <v-col
                        v-for="stat in statSummary"
                        :key="stat.label"
                        cols="auto"
                      >
                        <v-chip size="small" variant="flat">
                          <strong class="mr-1">{{ stat.label }}:</strong>
                          {{ fmtCurrency(stat.value) }}
                        </v-chip>
                      </v-col>
                    </v-row>
                  </v-card>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="12">
                  <v-card variant="outlined">
                    <v-card-title class="text-subtitle-1 pa-3"
                      >Ledger</v-card-title
                    >
                    <v-card-text class="pa-0">
                      <v-data-table
                        :headers="ledgerHeaders"
                        :items="employerStatement.line_items ?? []"
                        density="compact"
                        hide-default-footer
                        :items-per-page="-1"
                      >
                        <template #[`item.date`]="{ item }: { item: any }">{{
                          fmtDate(item.date)
                        }}</template>
                        <template #[`item.debit`]="{ item }: { item: any }">{{
                          item.debit ? fmtCurrency(item.debit) : ''
                        }}</template>
                        <template #[`item.credit`]="{ item }: { item: any }">{{
                          item.credit ? fmtCurrency(item.credit) : ''
                        }}</template>
                        <template #[`item.balance`]="{ item }: { item: any }">
                          <span
                            :class="
                              item.balance > 0 ? 'text-error' : 'text-success'
                            "
                          >
                            {{ fmtCurrency(item.balance) }}
                          </span>
                        </template>
                      </v-data-table>
                    </v-card-text>
                  </v-card>
                </v-col>
              </v-row>
            </template>
            <v-alert v-else type="info" variant="tonal" class="mt-4">
              Select a scheme and date range, then click Generate to view the
              statement.
            </v-alert>
          </v-window-item>

          <!-- ── Broker Commission Tab ───────────────────────────────── -->
          <v-window-item value="broker">
            <v-row class="mb-4" align="center">
              <v-col cols="12" md="3">
                <v-text-field
                  v-model.number="brokerFilters.brokerId"
                  label="Broker ID *"
                  type="number"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" md="2">
                <v-text-field
                  v-model="brokerFilters.from"
                  label="From"
                  type="date"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" md="2">
                <v-text-field
                  v-model="brokerFilters.to"
                  label="To"
                  type="date"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" md="5" class="d-flex align-center">
                <v-btn
                  color="primary"
                  prepend-icon="mdi-file-find"
                  :loading="loadingBroker"
                  :disabled="
                    !brokerFilters.brokerId ||
                    !brokerFilters.from ||
                    !brokerFilters.to
                  "
                  @click="loadBrokerStatement"
                >
                  Generate
                </v-btn>
              </v-col>
            </v-row>

            <template v-if="brokerStatement">
              <v-row class="mb-2">
                <v-col cols="12">
                  <v-card variant="tonal" color="primary" class="pa-3">
                    <div class="text-h6 mb-1">{{
                      brokerStatement.broker_name
                    }}</div>
                    <div class="text-body-2"
                      >Period: {{ brokerStatement.period }}</div
                    >
                    <v-chip class="mt-2" size="small">
                      <strong class="mr-1">Total Earned:</strong>
                      {{ fmtCurrency(brokerStatement.total_earned) }}
                    </v-chip>
                  </v-card>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="12">
                  <v-card variant="outlined">
                    <v-card-title class="text-subtitle-1 pa-3"
                      >Commission by Scheme</v-card-title
                    >
                    <v-card-text class="pa-0">
                      <v-data-table
                        :headers="commissionHeaders"
                        :items="brokerStatement.schemes ?? []"
                        density="compact"
                        hide-default-footer
                        :items-per-page="-1"
                      >
                        <template
                          #[`item.premium_collected`]="{ item }: { item: any }"
                          >{{ fmtCurrency(item.premium_collected) }}</template
                        >
                        <template
                          #[`item.commission_rate`]="{ item }: { item: any }"
                          >{{
                            (item.commission_rate * 100).toFixed(2)
                          }}%</template
                        >
                        <template
                          #[`item.commission_earned`]="{ item }: { item: any }"
                          >{{ fmtCurrency(item.commission_earned) }}</template
                        >
                      </v-data-table>
                    </v-card-text>
                  </v-card>
                </v-col>
              </v-row>
            </template>
            <v-alert v-else type="info" variant="tonal" class="mt-4">
              Enter a Broker ID and date range, then click Generate to view the
              commission statement.
            </v-alert>
          </v-window-item>
        </v-window>
      </template>
    </base-card>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import PremiumManagementService from '@/renderer/api/PremiumManagementService'
import BaseCard from '@/renderer/components/BaseCard.vue'
import { fmtDate } from '@/renderer/utils/formatters'
import jsPDF from 'jspdf'
import { applyPlugin } from 'jspdf-autotable'

applyPlugin(jsPDF)

const activeTab = ref('employer')
const loadingEmployer = ref(false)
const loadingBroker = ref(false)
const schemesLoading = ref(false)
const inForceSchemes = ref<{ id: number; name: string }[]>([])

const employerStatement = ref<any>(null)
const brokerStatement = ref<any>(null)

const now = new Date()
const firstOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
  .toISOString()
  .slice(0, 10)
const today = now.toISOString().slice(0, 10)

const employerFilters = ref({
  schemeId: null as number | null,
  from: firstOfMonth,
  to: today
})
const brokerFilters = ref({
  brokerId: null as number | null,
  from: firstOfMonth,
  to: today
})

const statSummary = computed(() => {
  if (!employerStatement.value) return []
  const s = employerStatement.value
  return [
    { label: 'Opening Balance', value: s.opening_balance },
    { label: 'Invoiced', value: s.invoiced_amount },
    { label: 'Received', value: s.received },
    { label: 'Adjustments', value: s.adjustments },
    { label: 'Closing Balance', value: s.closing_balance }
  ]
})

const ledgerHeaders = [
  { title: 'Date', key: 'date' },
  { title: 'Description', key: 'description' },
  { title: 'Debit', key: 'debit', align: 'end' as const },
  { title: 'Credit', key: 'credit', align: 'end' as const },
  { title: 'Balance', key: 'balance', align: 'end' as const }
]

const commissionHeaders = [
  { title: 'Scheme', key: 'scheme_name' },
  {
    title: 'Premium Collected',
    key: 'premium_collected',
    align: 'end' as const
  },
  { title: 'Commission Rate', key: 'commission_rate', align: 'end' as const },
  {
    title: 'Commission Earned',
    key: 'commission_earned',
    align: 'end' as const
  },
  { title: 'Status', key: 'status' }
]

function fmtCurrency(val: number | undefined) {
  if (val == null) return '—'
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR',
    maximumFractionDigits: 0
  }).format(val)
}

async function loadEmployerStatement() {
  if (!employerFilters.value.schemeId) return
  loadingEmployer.value = true
  try {
    const res = await PremiumManagementService.getEmployerStatement(
      employerFilters.value.schemeId,
      { from: employerFilters.value.from, to: employerFilters.value.to }
    )
    employerStatement.value = res.data.data
  } catch (e) {
    console.error('Failed to load employer statement', e)
  } finally {
    loadingEmployer.value = false
  }
}

async function loadBrokerStatement() {
  if (!brokerFilters.value.brokerId) return
  loadingBroker.value = true
  try {
    const res = await PremiumManagementService.getBrokerStatement(
      brokerFilters.value.brokerId,
      {
        from: brokerFilters.value.from,
        to: brokerFilters.value.to
      }
    )
    brokerStatement.value = res.data.data
  } catch (e) {
    console.error('Failed to load broker statement', e)
  } finally {
    loadingBroker.value = false
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
})

function downloadEmployerCSV() {
  if (!employerStatement.value) return
  const rows = [
    ['Date', 'Description', 'Debit', 'Credit', 'Balance'],
    ...(employerStatement.value.line_items ?? []).map((li: any) => [
      li.date,
      li.description,
      li.debit ?? '',
      li.credit ?? '',
      li.balance
    ])
  ]
  const csv = rows.map((r) => r.join(',')).join('\n')
  const url = URL.createObjectURL(new Blob([csv], { type: 'text/csv' }))
  const a = document.createElement('a')
  a.href = url
  a.download = `statement_scheme_${employerFilters.value.schemeId}.csv`
  a.click()
  URL.revokeObjectURL(url)
}

function downloadEmployerPDF() {
  if (!employerStatement.value) return
  const s = employerStatement.value

  // eslint-disable-next-line new-cap
  const doc: any = new jsPDF()
  const pageWidth = doc.internal.pageSize.width
  const margin = 15

  // Header bar
  doc.setFillColor(30, 58, 95)
  doc.rect(0, 0, pageWidth, 22, 'F')
  doc.setFontSize(16)
  doc.setFont('helvetica', 'bold')
  doc.setTextColor(255, 255, 255)
  doc.text('Employer Statement', margin, 15)

  // Scheme name & period
  let y = 32
  doc.setTextColor(30, 58, 95)
  doc.setFontSize(14)
  doc.text(s.scheme_name ?? '', margin, y)
  y += 7
  doc.setFontSize(10)
  doc.setFont('helvetica', 'normal')
  doc.setTextColor(80, 80, 80)
  doc.text(
    `Period: ${s.period ?? `${employerFilters.value.from} to ${employerFilters.value.to}`}`,
    margin,
    y
  )
  y += 10

  // Summary chips
  const summaryItems = [
    { label: 'Opening Balance', value: s.opening_balance },
    { label: 'Invoiced', value: s.invoiced_amount },
    { label: 'Received', value: s.received },
    { label: 'Adjustments', value: s.adjustments },
    { label: 'Closing Balance', value: s.closing_balance }
  ]

  doc.autoTable({
    startY: y,
    head: [summaryItems.map((i) => i.label)],
    body: [summaryItems.map((i) => fmtCurrency(i.value))],
    theme: 'grid',
    headStyles: {
      fillColor: [30, 58, 95],
      textColor: 255,
      fontStyle: 'bold',
      fontSize: 8,
      halign: 'center'
    },
    bodyStyles: {
      fontSize: 9,
      halign: 'center',
      fontStyle: 'bold'
    },
    margin: { left: margin, right: margin }
  })

  y = doc.lastAutoTable.finalY + 12

  // Ledger section title
  doc.setFontSize(12)
  doc.setFont('helvetica', 'bold')
  doc.setTextColor(30, 58, 95)
  doc.text('Ledger', margin, y)
  y += 4

  // Ledger table
  const ledgerData = (s.line_items ?? []).map((li: any) => [
    li.date,
    li.description,
    li.debit ? fmtCurrency(li.debit) : '',
    li.credit ? fmtCurrency(li.credit) : '',
    fmtCurrency(li.balance)
  ])

  doc.autoTable({
    startY: y,
    head: [['Date', 'Description', 'Debit', 'Credit', 'Balance']],
    body: ledgerData,
    theme: 'striped',
    headStyles: {
      fillColor: [30, 58, 95],
      textColor: 255,
      fontStyle: 'bold',
      fontSize: 9
    },
    bodyStyles: {
      fontSize: 8
    },
    columnStyles: {
      0: { cellWidth: 28 },
      1: { cellWidth: 60 },
      2: { halign: 'right', cellWidth: 30 },
      3: { halign: 'right', cellWidth: 30 },
      4: { halign: 'right', cellWidth: 30 }
    },
    alternateRowStyles: {
      fillColor: [245, 248, 252]
    },
    margin: { left: margin, right: margin }
  })

  // Footer with generation date
  const finalY = doc.lastAutoTable.finalY + 10
  doc.setFontSize(7)
  doc.setFont('helvetica', 'normal')
  doc.setTextColor(150, 150, 150)
  doc.text(
    `Generated on ${new Date().toLocaleDateString('en-ZA')}`,
    margin,
    finalY
  )

  doc.save(
    `statement_${s.scheme_name?.replace(/\s+/g, '_') ?? employerFilters.value.schemeId}_${employerFilters.value.from}_${employerFilters.value.to}.pdf`
  )
}
</script>
