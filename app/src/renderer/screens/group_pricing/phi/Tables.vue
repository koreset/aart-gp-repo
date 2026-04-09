<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :showActions="false">
          <template #header> {{ viewHeader }} </template>
          <template #default>
            <v-container fluid>
              <v-row>
                <v-col>
                  <v-table hover class="rating-tables">
                    <thead>
                      <tr>
                        <th style="width: 50%">Table Name</th>
                        <th style="width: 10%; text-align: center">Status</th>
                        <th style="width: 40%; text-align: center">Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="item in tables" :key="item.table_type">
                        <td>
                          <div class="d-flex align-center py-1">
                            <v-icon
                              :color="
                                item.populated ? 'success' : 'grey-lighten-1'
                              "
                              size="small"
                              class="mr-3"
                            >
                              {{
                                item.populated
                                  ? 'mdi-table-check'
                                  : 'mdi-table-alert'
                              }}
                            </v-icon>
                            <span>{{ item.table_type }}</span>
                          </div>
                        </td>
                        <td style="text-align: center">
                          <v-chip
                            :color="item.populated ? 'success' : 'warning'"
                            size="small"
                            variant="tonal"
                            label
                          >
                            {{ item.populated ? 'Populated' : 'Empty' }}
                          </v-chip>
                        </td>
                        <td style="text-align: center">
                          <div class="d-flex align-center justify-center ga-2">
                            <v-btn
                              variant="text"
                              size="small"
                              color="primary"
                              @click.stop="viewTable(item)"
                            >
                              <v-icon start size="small"
                                >mdi-information-outline</v-icon
                              >
                              Info
                            </v-btn>
                            <file-updater
                              :show-year="true"
                              :show-version="true"
                              :uploadComplete="uploadComplete"
                              :tableType="item.table_type"
                              @uploadFile="handleUpload"
                            ></file-updater>
                            <v-btn
                              variant="text"
                              size="small"
                              color="error"
                              @click.stop="deleteTableDataFlow(item)"
                            >
                              <v-icon start size="small"
                                >mdi-delete-outline</v-icon
                              >
                              Delete
                            </v-btn>
                          </div>
                        </td>
                      </tr>
                    </tbody>
                  </v-table>
                </v-col>
              </v-row>
              <v-row v-if="tableData.length > 0 && !loadingData">
                <v-col>
                  <data-grid
                    :show-close-button="true"
                    :columnDefs="columnDefs"
                    :rowData="tableData"
                    :table-title="selectedTable"
                    :pagination="true"
                    :rowCount="rowCount"
                    @update:clear-data="clearData"
                  />
                </v-col>
              </v-row>
              <v-row v-if="loadingData">
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
              <v-row>
                <v-col>
                  <v-expansion-panels variant="inset" class="my-4">
                    <v-expansion-panel title="PHI Model Points">
                      <v-expansion-panel-text>
                        <v-row class="mb-2">
                          <v-col class="d-flex justify-end">
                            <file-updater
                              :show-year="true"
                              :show-version="true"
                              :upload-complete="uploadPhiMpComplete"
                              table-type="PHI Model Points"
                              action-name="Upload Model Points"
                              @upload-file="handlePhiMpUpload"
                            ></file-updater>
                          </v-col>
                        </v-row>
                        <v-row>
                          <v-col cols="3">
                            <v-select
                              v-model="selectedPhiYear"
                              :loading="loadingPhiMps"
                              variant="outlined"
                              density="compact"
                              label="Select an MP year"
                              :items="sortedUniquePhiYears"
                              @update:model-value="getPhiVersions"
                            ></v-select>
                          </v-col>
                        </v-row>
                        <v-row v-if="phiMpVersions.length > 0">
                          <v-col>
                            <v-data-table
                              :headers="phiMpVersionHeaders"
                              :items="phiMpVersions"
                              :items-per-page="5"
                              :items-per-page-options="[5, 10, 25, 50]"
                              class="elevation-1"
                              item-value="version"
                            >
                              <template #[`item.actions`]="{ item }">
                                <v-btn
                                  variant="text"
                                  class="mr-1"
                                  :loading="
                                    loadingPhiExcelData.has(
                                      `${(item as any).year}_${(item as any).version}`
                                    )
                                  "
                                  size="small"
                                  color="primary"
                                  @click="getPhiModelPointsExcel(item as any)"
                                >
                                  <v-icon start size="small"
                                    >mdi-download</v-icon
                                  >
                                  Download
                                </v-btn>
                                <v-btn
                                  variant="text"
                                  class="mr-1"
                                  size="small"
                                  color="primary"
                                  @click="getPhiModelPoints(item as any)"
                                >
                                  <v-icon start size="small"
                                    >mdi-eye-outline</v-icon
                                  >
                                  View
                                </v-btn>
                                <v-btn
                                  variant="text"
                                  size="small"
                                  color="error"
                                  @click="deletePhiModelPoints(item as any)"
                                >
                                  <v-icon start size="small"
                                    >mdi-delete-outline</v-icon
                                  >
                                  Delete
                                </v-btn>
                              </template>
                            </v-data-table>
                          </v-col>
                        </v-row>
                        <v-row v-if="loadingPhiData">
                          <v-col>
                            <p class="mt-3">Loading model points...</p>
                            <v-progress-linear
                              buffer-value="20"
                              color="primary"
                              stream
                              value="10"
                            ></v-progress-linear>
                          </v-col>
                        </v-row>
                        <v-row v-if="phiMpData.length > 0 && !loadingPhiData">
                          <v-col>
                            <data-grid
                              :show-close-button="true"
                              :column-defs="phiColumnDefs"
                              :row-data="phiMpData"
                              :table-title="`PHI Model Points`"
                              @update:clear-data="clearPhiData"
                            />
                          </v-col>
                        </v-row>
                      </v-expansion-panel-text>
                    </v-expansion-panel>
                  </v-expansion-panels>
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
    <v-dialog v-model="deleteTableDialog" persistent max-width="600">
      <base-card>
        <template #header
          ><span class="headline"
            >Available years for {{ selectedTable.table_type }}</span
          ></template
        >
        <template #default>
          <v-row v-if="availableTableYears.length > 0" class="mt-5">
            <v-col>
              <v-select
                v-model="selectedTableYear"
                variant="outlined"
                density="compact"
                label="Data Year"
                placeholder="Select an existing Year"
                :items="availableTableYears"
                item-title="Year"
                item-value="Year"
                @update:model-value="getYearVersions"
              ></v-select>
            </v-col>
          </v-row>
          <v-row
            v-if="availableYearVersions && availableYearVersions.length > 0"
            class="mt-5"
          >
            <v-col>
              <v-select
                v-model="selectedYearVersion"
                variant="outlined"
                density="compact"
                label="Version"
                placeholder="Select a Version"
                :items="availableYearVersions"
              ></v-select>
            </v-col>
          </v-row>
        </template>
        <template #actions>
          <v-btn
            color="primary darken-1"
            variant="text"
            @click="deleteTableData()"
            >Proceed</v-btn
          >
          <v-btn
            color="primary darken-1"
            variant="text"
            @click="cancelDeleteTableDialog()"
            >Cancel</v-btn
          >
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

import PhiValuationService from '@/renderer/api/PhiValuationService'
import formatValues from '@/renderer/utils/format_values'
import { ref, onMounted } from 'vue'
import { DataPayload } from '@/renderer/components/types'

// data
// const otherTableDialog: any = ref(false)
const viewHeader: string = 'PHI Tables'
const confirmDeleteDialog: any = ref()
const deleteTableDialog: any = ref(false)
const loadingData: any = ref(false)
const rowCount: any = ref(0)
const tableData: any = ref([])
const selectedTable: any = ref('')
const selectedTableYear: any = ref(null)
const selectedYearVersion: any = ref(null)
const availableYearVersions: any = ref([])
const tables: any = ref([
  { name: 'Parameters' },
  { name: 'Yield Curve' },
  { name: 'Margins' },
  { name: 'Shocks' }
])

const availableTableYears: any = ref([])

const modelPoints: any = ref([])

// PHI Model Points data
const phiModelPointCount: any = ref([])
const uniquePhiYears: any = ref([])
const sortedUniquePhiYears: any = ref([])
const selectedPhiYear: any = ref(null)
const phiMpVersions: any = ref([])
const phiMpData: any = ref([])
const phiColumnDefs: any = ref([])
const loadingPhiMps = ref(false)
const loadingPhiData = ref(false)
const loadingPhiExcelData: any = ref(new Set())

const text: any = ref('')
// const yieldCurveDataDialog: any = ref(false)
const columnDefs: any = ref([])
// const dialog: any = ref(false)

const uploadComplete = ref(false)
const uploadPhiMpComplete = ref(false)
const snackbarText: any = ref(null)
const timeout: any = ref(3000)
const snackbar: any = ref(false)

// Headers for the PHI model points versions data table
const phiMpVersionHeaders = ref([
  { title: 'Version', key: 'version', align: 'start' as const },
  { title: 'Data Count', key: 'count', align: 'start' as const },
  {
    title: 'Actions',
    key: 'actions',
    sortable: false,
    align: 'center' as const
  }
])

const loadPhiModelPointCount = async () => {
  loadingPhiMps.value = true
  try {
    const phiMpResponse = await PhiValuationService.getPhiModelPointCount()
    if (phiMpResponse.data && phiMpResponse.data.length > 0) {
      phiModelPointCount.value = phiMpResponse.data
      uniquePhiYears.value = Array.from(
        new Set(phiModelPointCount.value.map((item: any) => item.year))
      )
      sortedUniquePhiYears.value = [...uniquePhiYears.value].sort(
        (a: any, b: any) => b - a
      )
    } else {
      phiModelPointCount.value = []
      uniquePhiYears.value = []
      sortedUniquePhiYears.value = []
    }
  } catch (error) {
    console.error('Error loading PHI model point count:', error)
    phiModelPointCount.value = []
    uniquePhiYears.value = []
    sortedUniquePhiYears.value = []
  } finally {
    loadingPhiMps.value = false
  }
}

onMounted(async () => {
  try {
    const res = await PhiValuationService.getTableMetaData()
    // Exclude PHI Model Points from the generic tables list — handled in its own expansion panel
    tables.value = res.data.table_meta_data.filter(
      (t: any) => t.table_type !== 'PHI Model Points'
    )
    modelPoints.value = res.data.model_points
  } catch (error) {
    console.error('Error loading PHI table metadata:', error)
  }

  await loadPhiModelPointCount()
})

// methods
const handlePhiMpUpload = (payload: DataPayload) => {
  uploadPhiMpComplete.value = false
  const formdata: any = new FormData()
  formdata.append('file', payload.file)
  formdata.append('table_type', 'PHI Model Points')
  formdata.append('year', payload.selectedYear)
  formdata.append('version', payload.version)
  PhiValuationService.uploadTables(formdata)
    .then(async (res: any) => {
      if (res.status === 200) {
        snackbarText.value = 'PHI model points uploaded successfully'
        snackbar.value = true
        uploadPhiMpComplete.value = true
        // Refresh the count list so the new version appears
        await loadPhiModelPointCount()
        // If the same year is still selected, refresh its version list
        if (selectedPhiYear.value) {
          phiMpVersions.value = phiModelPointCount.value.filter(
            (item: any) => item.year === selectedPhiYear.value
          )
        }
      }
    })
    .catch((error: any) => {
      let errorMessage = 'An error occurred while uploading PHI model points'
      if (
        error.response &&
        error.response.status === 400 &&
        error.response.data &&
        error.response.data.error
      ) {
        errorMessage = error.response.data.error
      }
      snackbarText.value = errorMessage
      snackbar.value = true
      uploadPhiMpComplete.value = true
    })
}

const handleUpload = (payload: DataPayload) => {
  uploadComplete.value = false
  const formdata: any = new FormData()
  formdata.append('file', payload.file)
  formdata.append('table_type', payload.selectedType)
  formdata.append('year', payload.selectedYear)
  formdata.append('version', payload.version)
  PhiValuationService.uploadTables(formdata)
    .then((res: any) => {
      if (res.status === 200) {
        snackbarText.value = 'File uploaded successfully'
        snackbar.value = true
        timeout.value = 3000
        uploadComplete.value = true
      }
    })
    .catch((error: any) => {
      let errorMessage = 'An error occurred while uploading the file ooooo'

      console.log('Upload error details:', error)

      if (
        error.response &&
        error.response.status === 400 &&
        error.response.data &&
        error.response.data.error
      ) {
        errorMessage = error.response.data.error
      }

      snackbarText.value = errorMessage
      snackbar.value = true
      timeout.value = 3000
      uploadComplete.value = true
    })
}

const cancelDeleteTableDialog = () => {
  selectedTableYear.value = null
  selectedYearVersion.value = null
  deleteTableDialog.value = false
}

const getYearVersions = (year: any) => {
  if (selectedTable.value === null) {
    text.value = 'Please select a table to get year versions'
    snackbar.value = true
    return
  }
  PhiValuationService.getAvailableVersionsForTableYear(
    selectedTable.value.table_name,
    year
  ).then((response) => {
    availableYearVersions.value = response.data
  })
}

const deleteTableDataFlow = async (table: any) => {
  // get the available years for the table as we delete by year and version
  const res = await PhiValuationService.getAvailableYearsForTable(
    table.table_name
  )

  availableTableYears.value = res.data
  availableYearVersions.value = []
  selectedTable.value = table
  deleteTableDialog.value = true

  // if (table.table_type === 'Yield Curve') {
  //   console.log('deleting yield curve data')
  //   yieldCurveDataDialog.value = true
  // }
  // if (table.table_type === 'Shocks') {
  //   console.log('deleting shocks data')
  //   // shocksDialog.value = true
  // }

  // if (table.table_type !== 'Yield Curve' && table.table_type !== 'Shocks') {
  //   try {
  //     const result = await confirmDeleteDialog.value.open(
  //       'Deleting Data for ' + table.table_type + ' table',
  //       'Are you sure you want to delete this data?'
  //     )
  //     if (result) {
  //       PhiValuationService.deleteTable(table.table_type).then((response) => {
  //         text.value = response.data
  //         snackbar.value = true
  //         tableData.value = []
  //         selectedTable.value = ''
  //       })
  //     }
  //   } catch (error) {
  //     console.log(error)
  //   }
  // }
}

const deleteTableData = async () => {
  try {
    deleteTableDialog.value = false
    console.log(
      'Deleting data for',
      selectedTable.value.table_name,
      'in year',
      selectedTableYear.value,
      'and version',
      selectedYearVersion.value
    )
    if (
      selectedTableYear.value === null ||
      selectedYearVersion.value === null ||
      selectedTable.value === null
    ) {
      text.value = 'Please select a year to delete data'
      snackbar.value = true
      return
    }

    const res = await confirmDeleteDialog.value.open(
      'Deleting Data for ' +
        selectedTable.value.table_name +
        ' table for year ' +
        selectedTableYear.value +
        ' version ' +
        selectedYearVersion.value,
      'Are you sure you want to delete this data?'
    )

    if (!res) {
      return
    }

    PhiValuationService.deleteTable(
      selectedTable.value.table_name,
      selectedTableYear.value,
      selectedYearVersion.value
    ).then((response) => {
      text.value = response.data
      snackbar.value = true
      tableData.value = []
      selectedTable.value = ''
      selectedTableYear.value = null
      selectedYearVersion.value = null
    })
  } catch (err) {}
}

// methods
const clearData = () => {
  tableData.value = []
  selectedTable.value = ''
}

const viewTable = (item: any) => {
  loadingData.value = true
  tableData.value = []

  const tableType = item.table_type.replace(/ /g, '').toLowerCase()
  PhiValuationService.getDataForTable(tableType).then((res) => {
    if (res.data === null) {
      res.data = []
    }

    if (res.data.length === 0) {
      text.value = 'No data available for this table'
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

// PHI Model Points Methods
const getPhiVersions = async () => {
  phiMpVersions.value = []
  phiMpData.value = []
  const year = selectedPhiYear.value
  phiMpVersions.value = phiModelPointCount.value.filter(
    (item: { year: any }) => item.year === year
  )
}

const clearPhiData = () => {
  phiMpData.value = []
}

const getPhiModelPoints = async (item: any) => {
  phiMpData.value = []
  loadingPhiData.value = true

  try {
    const response = await PhiValuationService.getPhiModelPointsForYear(
      item.year,
      item.version
    )
    if (response.data !== null && response.data.length > 0) {
      phiColumnDefs.value = []
      phiMpData.value = []
      createPhiColumnDefs(response.data)

      response.data.forEach((item: any) => {
        const transformed: any = {}
        const keys = Object.keys(item)
        keys.forEach((key) => {
          if (isNaN(item[key])) {
            transformed[key] = item[key]
          } else {
            const value = Number(item[key])
            transformed[key] = value
          }
        })
        phiMpData.value.push(transformed)
      })
    } else {
      snackbarText.value =
        'No PHI model points available for this year and version'
      snackbar.value = true
    }
  } catch (error) {
    console.error('Error loading PHI model points:', error)
    snackbarText.value = 'Error loading PHI model points'
    snackbar.value = true
  }

  loadingPhiData.value = false
}

const getPhiModelPointsExcel = async (item: any) => {
  const itemKey = `${item.year}_${item.version}`
  loadingPhiExcelData.value.add(itemKey)

  try {
    const response = await PhiValuationService.getPhiModelPointsExcel(
      item.year,
      item.version
    )
    loadingPhiExcelData.value.delete(itemKey)

    if (response.data !== null) {
      const fileURL = window.URL.createObjectURL(new Blob([response.data]))
      const fileName = `phi_model_points_${item.year}_${item.version}.xlsx`
      const link = document.createElement('a')
      link.href = fileURL
      link.setAttribute('download', fileName)
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    }
  } catch (error) {
    loadingPhiExcelData.value.delete(itemKey)
    console.error('Error downloading PHI model points Excel file:', error)
    snackbarText.value = 'Error downloading PHI model points file'
    snackbar.value = true
  }
}

const deletePhiModelPoints = async (item: any) => {
  try {
    const result = await confirmDeleteDialog.value.open(
      'Delete PHI Model Points',
      `Are you sure you want to delete the PHI model points for year ${item.year} version ${item.version}?`
    )

    if (result) {
      const response = await PhiValuationService.deletePhiModelPoints(
        item.year,
        item.version
      )
      if (response.status === 200) {
        phiMpVersions.value = phiMpVersions.value.filter(
          (i: any) => i.version !== item.version
        )
        snackbarText.value = 'PHI model points deleted successfully'
        snackbar.value = true
        // Clear the displayed data if it was showing the deleted version
        if (phiMpData.value.length > 0) {
          phiMpData.value = []
        }
      }
    }
  } catch (error) {
    console.error('Error deleting PHI model points:', error)
    snackbarText.value = 'Error deleting PHI model points'
    snackbar.value = true
  }
}

const createPhiColumnDefs = (data: any) => {
  phiColumnDefs.value = []
  Object.keys(data[0]).forEach((element) => {
    const header: any = {}
    header.headerName = element
    header.field = element
    header.valueFormatter = formatValues
    header.minWidth = 220
    header.filter = true
    header.resizable = true
    header.sortable = true
    phiColumnDefs.value.push(header)
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
