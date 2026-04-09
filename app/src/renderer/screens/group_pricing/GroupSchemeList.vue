<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="true">
          <template #header>
            <span class="headline">Group Pricing Schemes In Force</span>
          </template>
          <template #default>
            <!-- KPI Summary Cards -->
            <v-row class="mb-4">
              <v-col cols="12" sm="6" md="3">
                <v-card variant="outlined" rounded="lg" class="pa-4 h-100">
                  <div class="text-overline text-medium-emphasis mb-1"
                    >ACTIVE SCHEMES</div
                  >
                  <div class="text-h4 font-weight-bold">{{
                    activeSchemeCount
                  }}</div>
                  <div class="text-caption text-success mt-1">
                    <v-icon size="x-small">mdi-arrow-up</v-icon>
                    {{ newThisMonth }} new this month
                  </div>
                </v-card>
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <v-card variant="outlined" rounded="lg" class="pa-4 h-100">
                  <div class="text-overline text-medium-emphasis mb-1"
                    >TOTAL MEMBERS</div
                  >
                  <div class="text-h4 font-weight-bold">{{
                    formatNumber(totalMembers)
                  }}</div>
                  <div class="text-caption text-info mt-1">
                    <v-icon size="x-small">mdi-account-group</v-icon>
                    Across all active schemes
                  </div>
                </v-card>
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <v-card variant="outlined" rounded="lg" class="pa-4 h-100">
                  <div class="text-overline text-medium-emphasis mb-1"
                    >PENDING CLAIMS</div
                  >
                  <div class="text-h4 font-weight-bold">{{
                    pendingClaimsCount
                  }}</div>
                  <div
                    v-if="flaggedClaimsCount > 0"
                    class="text-caption text-warning mt-1"
                  >
                    <v-icon size="x-small">mdi-arrow-up</v-icon>
                    {{ flaggedClaimsCount }} flagged
                  </div>
                  <div v-else class="text-caption text-success mt-1">
                    <v-icon size="x-small">mdi-check</v-icon>
                    All on track
                  </div>
                </v-card>
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <v-card variant="outlined" rounded="lg" class="pa-4 h-100">
                  <div class="text-overline text-medium-emphasis mb-1"
                    >RENEWAL DUE</div
                  >
                  <div class="text-h4 font-weight-bold">{{
                    renewalDueCount
                  }}</div>
                  <div
                    v-if="renewalDueCount > 0"
                    class="text-caption text-warning mt-1"
                  >
                    <v-icon size="x-small">mdi-alert</v-icon>
                    Within 30 days
                  </div>
                  <div v-else class="text-caption text-success mt-1">
                    <v-icon size="x-small">mdi-check</v-icon>
                    None due soon
                  </div>
                </v-card>
              </v-col>
            </v-row>

            <v-row v-if="schemes">
              <v-col>
                <!-- Search Input -->
                <v-row class="mb-4">
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="search"
                      label="Search schemes..."
                      prepend-inner-icon="mdi-magnify"
                      variant="outlined"
                      density="compact"
                      clearable
                      hide-details
                      class="search-box"
                    ></v-text-field>
                  </v-col>
                </v-row>
                <v-data-table
                  v-model:model-value="selectedScheme"
                  show-select
                  select-strategy="single"
                  class="table-row"
                  density="compact"
                  :headers="headers"
                  :items="schemes"
                  :search="search"
                  return-object
                >
                  <!-- Slot for Treaty Link Column -->
                  <template #[`item.has_treaty_link`]="{ item }: { item: any }">
                    <v-chip
                      :color="item.has_treaty_link ? 'success' : 'grey'"
                      size="small"
                      variant="flat"
                    >
                      {{ item.has_treaty_link ? 'Linked' : 'Not Linked' }}
                    </v-chip>
                  </template>
                  <!-- Slot for Status Column -->
                  <template #[`item.status`]="{ item }: { item: any }">
                    <status-chip :status="item.status" />
                  </template>
                  <!-- Slot for Actions Column -->
                  <template #[`item.actions`]="{ item }: { item: any }">
                    <v-tooltip>
                      <template #activator="{ props }">
                        <v-btn
                          v-if="selectedScheme?.[0]?.id === item.id"
                          icon
                          size="small"
                          variant="plain"
                          color="primary"
                          v-bind="props"
                          @click="editScheme(item)"
                        >
                          <v-icon>mdi-wrench</v-icon>
                        </v-btn>
                      </template>
                      <span>Maintain Scheme</span>
                    </v-tooltip>
                  </template>
                </v-data-table>
              </v-col>
            </v-row>
          </template>
          <template #actions>
            <v-btn
              v-if="hasPermission('reports:export')"
              color="primary"
              variant="outlined"
              prepend-icon="mdi-microsoft-excel"
              :disabled="!schemes || schemes.length === 0"
              @click="exportToExcel"
            >
              Export to Excel
            </v-btn>
          </template>
        </base-card>
      </v-col>
    </v-row>
    <confirmation-dialog ref="confirmDialog"></confirmation-dialog>
  </v-container>
</template>

<script setup lang="ts">
import BaseCard from '../../components/BaseCard.vue'
import StatusChip from '../../components/StatusChip.vue'
import { computed, onMounted, ref } from 'vue'
import ConfirmationDialog from '../../components/ConfirmDialog.vue'
import GroupPricingService from '../../api/GroupPricingService'
import {
  formatValues,
  roundUpToTwoDecimalsAccounting
} from '@/renderer/utils/format_values'
import { useRouter } from 'vue-router'
import formatDateString from '@/renderer/utils/helpers'
import * as XLSX from 'xlsx'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

const router = useRouter()
const { hasPermission } = usePermissionCheck()
const confirmDialog = ref()
const schemes: any = ref([])
const selectedScheme: any = ref([])
const claims: any = ref([])

const columnDefs: any = ref([])
const rowData = ref([])
const search = ref('')

// KPI computeds
const activeSchemeCount = computed(() => {
  return (
    schemes.value.filter(
      (s: any) =>
        s.status?.toLowerCase() === 'in_force' ||
        s.status?.toLowerCase() === 'active' ||
        s.in_force
    ).length || schemes.value.length
  )
})

const totalMembers = computed(() => {
  return schemes.value.reduce(
    (sum: number, s: any) => sum + (s.member_count || 0),
    0
  )
})

const newThisMonth = computed(() => {
  const now = new Date()
  const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
  return schemes.value.filter((s: any) => {
    if (!s.commencement_date) return false
    return new Date(s.commencement_date) >= startOfMonth
  }).length
})

const pendingClaimsCount = computed(() => {
  return claims.value.filter(
    (c: any) =>
      c.status?.toLowerCase() === 'pending' ||
      c.status?.toLowerCase() === 'submitted'
  ).length
})

const flaggedClaimsCount = computed(() => {
  return claims.value.filter(
    (c: any) =>
      c.priority?.toLowerCase() === 'high' ||
      c.priority?.toLowerCase() === 'critical'
  ).length
})

const renewalDueCount = computed(() => {
  const now = new Date()
  const thirtyDays = new Date(now.getTime() + 30 * 24 * 60 * 60 * 1000)
  return schemes.value.filter((s: any) => {
    if (!s.renewal_date) return false
    const renewal = new Date(s.renewal_date)
    return renewal >= now && renewal <= thirtyDays
  }).length
})

function formatNumber(val: number) {
  return new Intl.NumberFormat('en-ZA').format(val)
}

const headers = computed(() => {
  const baseHeaders: any = [
    { title: 'Scheme Name', value: 'name', sortable: true },
    { title: 'Treaty', value: 'has_treaty_link', sortable: true },
    { title: 'Quote In Force', value: 'quote_in_force', sortable: true },
    { title: 'Status', value: 'status', sortable: true },
    {
      title: 'Commencement Date',
      key: 'commencement_date',
      value: (item: any) =>
        formatDateString(item.commencement_date, true, true, true)
    },
    {
      title: 'Annual Premium',
      key: 'annual_premium',
      value: (item: any) => roundUpToTwoDecimalsAccounting(item.annual_premium)
    },
    {
      title: 'Duration in Force',
      value: 'duration_in_force_days',
      sortable: true
    },
    {
      title: 'Renewal Date',
      key: 'renewal_date',
      width: '20%',
      value: (item: any) =>
        formatDateString(item.renewal_date, true, true, true)
    },
    {
      title: 'Member Count',
      key: 'member_count',
      value: (item: any) =>
        new Intl.NumberFormat('fr-FR').format(item.member_count)
    },
    {
      title: 'Earned Premium',
      key: 'earned_premium',
      value: (item: any) => roundUpToTwoDecimalsAccounting(item.earned_premium)
    },
    {
      title: 'Cover Start Date',
      key: 'cover_start_date',
      value: (item: any) =>
        formatDateString(item.cover_start_date, true, true, true)
    },
    {
      title: 'Cover End Date',
      key: 'cover_end_date',
      value: (item: any) =>
        formatDateString(item.cover_end_date, true, true, true)
    }
  ]

  if (selectedScheme.value.length > 0) {
    baseHeaders.unshift({
      title: 'Actions',
      value: 'actions',
      align: 'center',
      sortable: false
    })
  }

  return baseHeaders
})

// const parseDateString = (dateString) => {
//   if (!dateString) return ''
//   console.log('dateString', dateString)
//   const date = new Date(dateString)
//   if (isNaN(date.getTime())) return ''
//   const formattedDate = date.toISOString().split('T')[0]
//   console.log('formattedDate', formattedDate)
//   return formattedDate
// }

// const createScheme = () => {
//   const schemPayload = {
//     name: schemeName.value
//   }
//   GroupPricingService.createScheme(schemPayload).then((res) => {
//     schemes.value.push(res.data)
//     schemeName.value = null
//   })
// }

// const deleteBroker = async (id: number) => {
// }

const editScheme = (item) => {
  router.push({ name: 'group-pricing-schemes-detail', params: { id: item.id } })
}

const exportToExcel = () => {
  if (!schemes.value || schemes.value.length === 0) return

  // Get filtered data based on search
  let dataToExport = schemes.value
  if (search.value) {
    const searchLower = search.value.toLowerCase()
    dataToExport = schemes.value.filter((item: any) => {
      return Object.values(item).some(
        (value: any) =>
          value && value.toString().toLowerCase().includes(searchLower)
      )
    })
  }

  // Create header row based on the headers computed property
  const headerRow = headers.value
    .filter((header) => header.value !== 'actions') // Exclude actions column
    .map((header) => header.title)

  // Create data rows
  const dataRows = dataToExport.map((item: any) => {
    return headers.value
      .filter((header) => header.value !== 'actions') // Exclude actions column
      .map((header) => {
        if (header.value) {
          // Handle computed values (functions)
          if (typeof header.value === 'function') {
            return header.value(item)
          }
          // Handle direct property access
          return item[header.value] || ''
        }
        return ''
      })
  })

  // Combine header and data
  const wsData = [headerRow, ...dataRows]

  // Create worksheet and workbook
  const ws = XLSX.utils.aoa_to_sheet(wsData)

  // Set column widths for better readability
  ws['!cols'] = headerRow.map(() => ({ wch: 20 }))

  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, 'Group Pricing Schemes')

  // Generate filename with current date
  const currentDate = new Date().toISOString().split('T')[0]
  const filename = `group_pricing_schemes_${currentDate}.xlsx`

  // Export to file
  XLSX.writeFile(wb, filename)
}

onMounted(() => {
  GroupPricingService.getSchemesInforce().then((res) => {
    const data = res.data?.data ?? res.data
    if (data && data.length > 0) {
      rowData.value = data
      schemes.value = data
      createColumnDefs(data)
    } else {
      schemes.value = []
    }
  })
  GroupPricingService.getClaims()
    .then((res) => {
      claims.value = res.data?.data ?? res.data ?? []
    })
    .catch(() => {
      claims.value = []
    })
})

const createColumnDefs = (data: any) => {
  columnDefs.value = []
  Object.keys(data[0]).forEach((element) => {
    const header: any = {}
    header.headerName = element
    header.field = element
    header.valueFormatter = formatValues
    header.minWidth = 200
    header.sortable = true
    header.filter = true
    header.resizable = true
    columnDefs.value.push(header)
  })
}
</script>

<style scoped>
.table-row {
  white-space: nowrap;
}

::v-deep(.v-data-table thead th) {
  background-color: #223f54 !important;
  color: white;
  text-align: center;
  font-weight: bold;
  white-space: nowrap;
  min-width: 150px;
}

.search-box {
  width: 100%;
}
.v-table__wrapper > table > thead {
  background-color: #223f54 !important;
  color: white;
  white-space: nowrap;
}
</style>
