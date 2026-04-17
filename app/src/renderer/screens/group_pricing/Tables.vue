<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :showActions="false">
          <template #header> {{ viewHeader }} </template>
          <template #default>
            <v-container fluid>
              <v-tabs v-model="activeTab" color="primary" class="mb-4">
                <v-tab value="group_pricing">Group Pricing</v-tab>
                <v-tab value="reinsurance">Reinsurance</v-tab>
                <v-tab value="binder_fees">Binder Fees</v-tab>
                <v-tab value="commission_structures"
                  >Commission Structures</v-tab
                >
              </v-tabs>

              <binder-fee-management v-if="activeTab === 'binder_fees'" />
              <commission-structure-management
                v-if="activeTab === 'commission_structures'"
              />

              <template v-if="isRatingTablesTab && hasEmpty">
                <v-row>
                  <v-col>
                    <v-alert
                      type="warning"
                      variant="tonal"
                      density="compact"
                      class="mb-2"
                    >
                      Some tables have not been populated yet. Upload data to
                      complete your setup.
                    </v-alert>
                  </v-col>
                </v-row>
              </template>

              <v-row v-if="isRatingTablesTab">
                <v-col>
                  <v-data-table
                    :headers="tableHeaders"
                    :items="filteredTables"
                    :items-per-page="5"
                    item-value="table_type"
                    hover
                    class="rating-tables"
                  >
                    <template #[`item.table_type`]="{ item: rawItem, index }">
                      <div class="d-flex align-center py-1">
                        <v-icon
                          :color="
                            (rawItem as any).populated
                              ? 'success'
                              : 'grey-lighten-1'
                          "
                          size="small"
                          class="mr-3"
                        >
                          {{
                            (rawItem as any).populated
                              ? 'mdi-table-check'
                              : 'mdi-table-alert'
                          }}
                        </v-icon>
                        <span :class="index === 0 ? 'font-weight-bold' : ''">
                          {{ (rawItem as any).table_type }}
                        </span>
                      </div>
                    </template>
                    <template #[`item.populated`]="{ item: rawItem2 }">
                      <div class="text-center">
                        <v-chip
                          :color="
                            (rawItem2 as any).populated ? 'success' : 'warning'
                          "
                          size="small"
                          variant="tonal"
                          label
                        >
                          {{
                            (rawItem2 as any).populated ? 'Populated' : 'Empty'
                          }}
                        </v-chip>
                      </div>
                    </template>
                    <template #[`item.actions`]="{ item: rawItem3 }">
                      <div class="d-flex align-center justify-center ga-2">
                        <v-btn
                          variant="text"
                          size="small"
                          color="primary"
                          @click.stop="viewTable(rawItem3)"
                        >
                          <v-icon start size="small"
                            >mdi-information-outline</v-icon
                          >
                          Info
                        </v-btn>
                        <file-updater
                          :show-risk-rate-code="
                            (rawItem3 as any).table_type ===
                            'Group Pricing Parameters'
                          "
                          :show-year="false"
                          :uploadComplete="uploadComplete"
                          :tableType="(rawItem3 as any).table_type"
                          :gp-risk-rate-codes="genericRiskRateCodes"
                          :show-risk-rate-code-select="
                            (rawItem3 as any).table_type !==
                            'Group Pricing Parameters'
                          "
                          @uploadFile="handleUpload"
                        ></file-updater>
                        <v-btn
                          variant="text"
                          size="small"
                          color="error"
                          @click.stop="chooseRiskRateCode(rawItem3)"
                        >
                          <v-icon start size="small">mdi-delete-outline</v-icon>
                          Delete
                        </v-btn>
                      </div>
                    </template>
                  </v-data-table>
                </v-col>
              </v-row>
              <v-row
                v-if="
                  isRatingTablesTab &&
                  tableData.length > 0 &&
                  !loadingData
                "
              >
                <v-col>
                  <data-grid
                    :columnDefs="columnDefs"
                    :show-close-button="true"
                    :rowData="tableData"
                    :table-title="selectedTable"
                    :pagination="true"
                    :rowCount="rowCount"
                    @update:clear-data="clearData"
                  />
                </v-col>
              </v-row>
              <v-row v-if="isRatingTablesTab && loadingData">
                <v-col>
                  <p class="mt-3">Loading...</p>
                  <v-progress-linear
                    buffer-value="20"
                    color="primary"
                    stream
                    value="10"
                  ></v-progress-linear>
                </v-col>
              </v-row>
            </v-container>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <v-snackbar
      v-model="snackbar"
      centered
      :timeout="timeout"
      :multi-line="true"
    >
      {{ snackbarText }}
      <v-btn rounded color="red" variant="text" @click="snackbar = false"
        >Close</v-btn
      >
    </v-snackbar>
    <confirmation-dialog ref="confirmDeleteDialog" />
    <v-dialog v-model="riskCodeDialog" persistent max-width="550px">
      <base-card>
        <template #header>
          <span class="headline">Choose the relevant Risk Code</span>
        </template>
        <template #default>
          <v-row>
            <v-col>
              <v-select
                v-model="selectedRiskRateCode"
                variant="outlined"
                density="compact"
                label="Select a risk rate code"
                :items="availableRiskRateCodes"
                item-title="year"
                item-value="year"
              ></v-select>
            </v-col>
          </v-row>
        </template>
        <template #actions>
          <v-spacer></v-spacer>
          <v-btn rounded variant="text" @click="closeDialog">Ok</v-btn>
        </template>
      </base-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import ConfirmationDialog from '@/renderer/components/ConfirmDialog.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import FileUpdater from '@/renderer/components/FileUpdater.vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import BinderFeeManagement from '@/renderer/screens/group_pricing/administration/BinderFeeManagement.vue'
import CommissionStructureManagement from '@/renderer/screens/group_pricing/administration/CommissionStructureManagement.vue'

import GroupPricingService from '@/renderer/api/GroupPricingService'
import formatValues from '@/renderer/utils/format_values'
import { ref, computed, onMounted } from 'vue'
import { DataPayload } from '@/renderer/components/types'
import * as pako from 'pako'

// data
const activeTab = ref('group_pricing')
const riskCodeDialog: any = ref(false)
const availableDataYears: any = ref([])
const selectedYear: any = ref('')
const confirmDeleteDialog: any = ref()
const loadingData: any = ref(false)
const rowCount: any = ref(0)
const tableData: any = ref([])
const selectedTable: any = ref('')
const selectedTableText: any = ref('')
const tables: any = ref([])

const columnDefs: any = ref([])

const uploadComplete = ref(false)
const snackbarText: any = ref(null)
const timeout: any = ref(3000)
const snackbar: any = ref(false)
const availableRiskRateCodes: any = ref([])
const genericRiskRateCodes: any = ref([])
const selectedRiskRateCode: any = ref('')
const hasEmpty: any = ref(false)

const tableHeaders = [
  { title: 'Table Name', key: 'table_type', width: '50%' },
  { title: 'Status', key: 'populated', width: '10%', align: 'center' as const },
  {
    title: 'Actions',
    key: 'actions',
    width: '40%',
    align: 'center' as const,
    sortable: false
  }
]

const viewHeader = computed(() => {
  if (activeTab.value === 'reinsurance') return 'Reinsurance Rating Tables'
  if (activeTab.value === 'binder_fees') return 'Binder Fees'
  if (activeTab.value === 'commission_structures') return 'Commission Structures'
  return 'Group Pricing Rating Tables'
})

const isRatingTablesTab = computed(
  () =>
    activeTab.value !== 'binder_fees' &&
    activeTab.value !== 'commission_structures'
)

const filteredTables = computed(() => {
  return tables.value.filter(
    (t: any) => (t.category || 'group_pricing') === activeTab.value
  )
})

onMounted(async () => {
  loadingData.value = true

  const [response] = await Promise.all([
    GroupPricingService.getTableMetaData(),
    getGenericRiskRateCodes()
  ])

  tables.value = response.data.associated_tables
  checkPopulated()
  loadingData.value = false
})

const checkPopulated = () => {
  hasEmpty.value = tables.value.some(
    (t: any) =>
      !t.populated && (t.category || 'group_pricing') === activeTab.value
  )
}

const closeDialog = () => {
  riskCodeDialog.value = false
  deleteTableData()
}

// methods
const clearData = () => {
  tableData.value = []
  selectedTable.value = ''
}

const getRiskRateCodes = async () => {
  const tableType = selectedTable.value.replace(/\s+/g, '').toLowerCase()
  const response = await GroupPricingService.getRiskRateCodes(tableType)
  availableRiskRateCodes.value = response.data
}

const getGenericRiskRateCodes = async () => {
  const response = await GroupPricingService.getRiskRateCodes(
    'grouppricingparameters'
  )
  genericRiskRateCodes.value = response.data
  console.log('genericRiskRateCodes', genericRiskRateCodes.value)
}

const chooseRiskRateCode = async (item: any) => {
  console.log('item', item)
  availableDataYears.value = []
  selectedYear.value = null
  selectedRiskRateCode.value = null
  selectedTableText.value = item.table_type
  selectedTable.value = item.table_type.replace(/\s+/g, '').toLowerCase()

  availableRiskRateCodes.value = []
  getRiskRateCodes()
  riskCodeDialog.value = true
}

const handleUpload = async (payload: DataPayload) => {
  uploadComplete.value = false

  try {
    // Read the file as array buffer for compression
    const fileBuffer = await payload.file.arrayBuffer()
    const uint8Array = new Uint8Array(fileBuffer)

    // Compress the file using gzip
    const compressed = pako.gzip(uint8Array)

    // Create a new File object with the compressed data
    const compressedFile = new File([compressed], payload.file.name + '.gz', {
      type: 'application/gzip'
    })

    const formdata: any = new FormData()
    formdata.append('file', compressedFile)
    formdata.append('table_type', payload.selectedType)
    formdata.append('risk_rate_code', payload.riskRateCode)
    if (payload.schemeName) {
      formdata.append('scheme_name', payload.schemeName)
    }
    formdata.append('original_filename', payload.file.name)
    formdata.append('compression', 'gzip')

    const response = await GroupPricingService.uploadTables(formdata)

    if (response.status === 200) {
      snackbarText.value = 'File uploaded successfully (compressed)'
      snackbar.value = true
      timeout.value = 3000
      uploadComplete.value = true
      checkPopulated()
    }
  } catch (error: any) {
    let errorMessage = 'An error occurred while uploading the file'

    console.log('Upload error details:', error.data)

    if (error.status === 400 && error.data && error.data.error) {
      errorMessage = error.data.error
    }

    snackbarText.value = errorMessage
    snackbar.value = true
    timeout.value = 3000
    uploadComplete.value = true
  } finally {
    // update generic risk rate codes
    getGenericRiskRateCodes()
  }

  GroupPricingService.getTableMetaData().then((res) => {
    tables.value = res.data.associated_tables
    console.log(tables.value)
  })
}

const deleteTableData = async () => {
  try {
    const result = await confirmDeleteDialog.value.open(
      'Deleting Data for ' + selectedTableText.value + ' table',
      'Are you sure you want to delete this data?'
    )
    if (result) {
      const response = await GroupPricingService.deleteTable(
        selectedTable.value,
        selectedRiskRateCode.value
      )
      if (response.status === 200) {
        snackbarText.value = 'Data deleted successfully'
        snackbar.value = true
        timeout.value = 3000
        tableData.value = []
        selectedTable.value = ''
        checkPopulated()
      }
    }
  } catch (error) {
    console.log(error)
  }
  GroupPricingService.getTableMetaData().then((res) => {
    tables.value = res.data.associated_tables
  })
}

const viewTable = (item: any) => {
  loadingData.value = true
  tableData.value = []

  const tableType = item.table_type.replace(/ /g, '').toLowerCase()
  GroupPricingService.getDataForTable(tableType).then((res) => {
    if (res.data === null) {
      res.data = []
    }

    if (res.data.length === 0) {
      snackbarText.value = 'No data available for this table'
      snackbar.value = true
    } else {
      tableData.value = res.data
      if (tableData.value.length > 0) {
        createColumnDefs(tableData.value)
      }
      selectedTable.value = item.table_type
    }
    loadingData.value = false
  })
}

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
.rating-tables :deep(th) {
  font-size: 0.8rem !important;
  font-weight: 600 !important;
  text-transform: uppercase;
  letter-spacing: 0.025em;
  color: rgba(0, 0, 0, 0.6);
}

.rating-tables :deep(td) {
  border-bottom: 1px solid rgba(0, 0, 0, 0.06) !important;
}

.rating-tables :deep(tr:last-child td) {
  border-bottom: none !important;
}
</style>
