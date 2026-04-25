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

              <v-row v-if="isRatingTablesTab && pendingChangeCount > 0">
                <v-col>
                  <v-alert
                    type="info"
                    variant="tonal"
                    density="compact"
                    class="mb-2"
                  >
                    <div class="d-flex align-center">
                      <span class="flex-grow-1">
                        {{ pendingChangeCount }} unsaved
                        {{ pendingChangeCount === 1 ? 'change' : 'changes' }}
                        to the Required configuration. Nothing is persisted or
                        recorded in the audit log until you click Save.
                      </span>
                      <v-btn
                        size="small"
                        variant="text"
                        color="primary"
                        :disabled="savingChanges"
                        class="mr-2"
                        @click="discardRequiredChanges"
                      >
                        Discard
                      </v-btn>
                      <v-btn
                        size="small"
                        variant="flat"
                        color="primary"
                        :loading="savingChanges"
                        @click="saveRequiredChanges"
                      >
                        Save changes
                      </v-btn>
                    </div>
                  </v-alert>
                </v-col>
              </v-row>

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
                    <template #[`item.is_required`]="{ item: rawItemR }">
                      <div class="d-flex align-center justify-center">
                        <v-checkbox
                          :model-value="effectiveRequired(rawItemR)"
                          density="compact"
                          hide-details
                          color="primary"
                          :disabled="savingChanges"
                          @update:model-value="
                            (val) => onRequiredChange(rawItemR, val)
                          "
                        />
                        <v-tooltip
                          v-if="hasPendingChange(rawItemR)"
                          location="top"
                        >
                          <template #activator="{ props }">
                            <v-icon
                              v-bind="props"
                              size="small"
                              color="warning"
                              class="ml-1"
                            >
                              mdi-circle-medium
                            </v-icon>
                          </template>
                          <span>Unsaved change — click Save to apply.</span>
                        </v-tooltip>
                        <v-tooltip location="top" max-width="320">
                          <template #activator="{ props }">
                            <v-icon
                              v-bind="props"
                              size="small"
                              color="grey"
                              class="ml-1"
                            >
                              mdi-information-outline
                            </v-icon>
                          </template>
                          <span v-if="effectiveRequired(rawItemR)">
                            Required. This table will be loaded and used by
                            downstream calculations. If no data has been
                            uploaded, the status will show "Empty".
                          </span>
                          <span v-else>
                            Not required. This table will not be read from the
                            database; any downstream variables that depend on it
                            will resolve to zero. The "Empty" warning is
                            suppressed.
                          </span>
                        </v-tooltip>
                      </div>
                    </template>
                    <template #[`item.populated`]="{ item: rawItem2 }">
                      <div class="text-center">
                        <v-chip
                          v-if="(rawItem2 as any).is_required === false"
                          color="default"
                          size="small"
                          variant="tonal"
                          label
                        >
                          Not required
                        </v-chip>
                        <v-chip
                          v-else
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
                        <v-btn
                          variant="text"
                          size="small"
                          color="primary"
                          :title="`View configuration history for ${(rawItem3 as any).table_type}`"
                          @click.stop="viewAuditHistory(rawItem3)"
                        >
                          <v-icon start size="small">mdi-history</v-icon>
                          History
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
                              'Group Pricing Parameters' &&
                            (rawItem3 as any).table_type !== 'Age Bands'
                          "
                          @uploadFile="handleUpload"
                        ></file-updater>
                        <v-btn
                          variant="text"
                          size="small"
                          color="error"
                          @click.stop="deleteClicked(rawItem3)"
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
                v-if="isRatingTablesTab && tableData.length > 0 && !loadingData"
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
    <confirmation-dialog ref="confirmSaveDialog" />
    <table-config-audit-dialog ref="auditDialog" />
    <v-dialog v-model="riskCodeDialog" persistent max-width="550px">
      <base-card>
        <template #header>
          <span class="headline">{{ pickerLabel }}</span>
        </template>
        <template #default>
          <v-row>
            <v-col>
              <v-select
                v-model="selectedRiskRateCode"
                variant="outlined"
                density="compact"
                :label="pickerFieldLabel"
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
import TableConfigAuditDialog from '@/renderer/components/TableConfigAuditDialog.vue'
import BinderFeeManagement from '@/renderer/screens/group_pricing/administration/BinderFeeManagement.vue'
import CommissionStructureManagement from '@/renderer/screens/group_pricing/administration/CommissionStructureManagement.vue'

import GroupPricingService from '@/renderer/api/GroupPricingService'
import formatValues from '@/renderer/utils/format_values'
import { ref, computed, onMounted, watch } from 'vue'
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
const auditDialog: any = ref(null)
const confirmSaveDialog: any = ref(null)
// Tentative checkbox edits keyed by canonical table_key (falls back to
// table_type display name). Empty = no unsaved changes; nothing is sent to
// the API or recorded in the audit log until the user clicks Save.
const pendingRequired = ref<Record<string, boolean>>({})
const savingChanges = ref(false)

const tableHeaders = [
  { title: 'Table Name', key: 'table_type', width: '40%' },
  {
    title: 'Required',
    key: 'is_required',
    width: '14%',
    align: 'center' as const,
    sortable: false
  },
  { title: 'Status', key: 'populated', width: '14%', align: 'center' as const },
  {
    title: 'Actions',
    key: 'actions',
    width: '32%',
    align: 'center' as const,
    sortable: false
  }
]

const viewHeader = computed(() => {
  if (activeTab.value === 'reinsurance') return 'Reinsurance Rating Tables'
  if (activeTab.value === 'binder_fees') return 'Binder Fees'
  if (activeTab.value === 'commission_structures')
    return 'Commission Structures'
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
  // Only count required-but-unpopulated tables — opting a table out of the
  // configuration intentionally suppresses the "Empty" warning for it.
  hasEmpty.value = tables.value.some(
    (t: any) =>
      !t.populated &&
      t.is_required !== false &&
      (t.category || 'group_pricing') === activeTab.value
  )
}

// Stable lookup key for an item. Falls back to display name when table_key
// is missing (older API payloads).
const itemKey = (item: any): string =>
  String(item?.table_key || item?.table_type || '')

// True if the item currently has an unsaved override.
const hasPendingChange = (item: any): boolean => {
  const k = itemKey(item)
  return Object.prototype.hasOwnProperty.call(pendingRequired.value, k)
}

// Effective (visible) required state — the unsaved override if present,
// otherwise the saved server-side value.
const effectiveRequired = (item: any): boolean => {
  const k = itemKey(item)
  if (Object.prototype.hasOwnProperty.call(pendingRequired.value, k)) {
    return pendingRequired.value[k]
  }
  return item?.is_required !== false
}

// Total unsaved changes across all tabs. Drives the Save banner visibility.
const pendingChangeCount = computed(
  () => Object.keys(pendingRequired.value).length
)

// Toggle handler: stores the new value as a pending edit. If the user toggles
// back to the original value, the entry is dropped so the banner clears.
const onRequiredChange = (item: any, newValue: boolean | null) => {
  const k = itemKey(item)
  if (!k) return
  const desired = newValue === true
  const original = item.is_required !== false
  const next = { ...pendingRequired.value }
  if (desired === original) {
    delete next[k]
  } else {
    next[k] = desired
  }
  pendingRequired.value = next
}

// Discard all unsaved overrides — checkboxes snap back to the saved values
// because the bindings re-read item.is_required.
const discardRequiredChanges = () => {
  pendingRequired.value = {}
}

// Persist all unsaved overrides. One PATCH per change (server writes one
// audit row per call). Surface a single confirmation dialog summarising the
// changes before any network call fires.
const saveRequiredChanges = async () => {
  const entries = Object.entries(pendingRequired.value)
  if (entries.length === 0) return

  const summary = entries
    .map(([key, val]) => {
      const display =
        tables.value.find((t: any) => itemKey(t) === key)?.table_type || key
      return `• ${display} → ${val ? 'Required' : 'Not required'}`
    })
    .join('\n')

  let confirmed = false
  try {
    confirmed = (await confirmSaveDialog.value.open(
      `Save ${entries.length} configuration ${
        entries.length === 1 ? 'change' : 'changes'
      }?`,
      `The following changes will be persisted and recorded in the audit log:\n\n${summary}`
    )) as boolean
  } catch {
    confirmed = false
  }
  if (!confirmed) return

  savingChanges.value = true
  let failed = 0
  try {
    for (const [key, val] of entries) {
      try {
        await GroupPricingService.updateTableConfiguration(key, val)
      } catch (e) {
        failed += 1
        console.error('Failed to save table configuration', key, e)
      }
    }

    // Refresh metadata to pick up the authoritative server state (including
    // the new updated_by / updated_at and audit history) and clear pending.
    const res = await GroupPricingService.getTableMetaData()
    tables.value = res.data.associated_tables
    pendingRequired.value = {}
    checkPopulated()

    if (failed === 0) {
      snackbarText.value = `Saved ${entries.length} ${
        entries.length === 1 ? 'change' : 'changes'
      } successfully.`
      timeout.value = 2500
    } else {
      snackbarText.value = `Saved ${
        entries.length - failed
      } of ${entries.length} changes — ${failed} failed.`
      timeout.value = 4000
    }
    snackbar.value = true
  } finally {
    savingChanges.value = false
  }
}

const viewAuditHistory = (item: any) => {
  const tableKey = item.table_key || item.table_type
  if (!tableKey || !auditDialog.value) return
  auditDialog.value.open(tableKey, item.table_type)
}

// Recompute the "some tables not populated" banner whenever the user
// switches tabs. Without this, hasEmpty reflects whichever tab was active
// at mount time.
watch(activeTab, () => {
  checkPopulated()
})

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
  // Prefer the authoritative delete_key from the backend; fall back to
  // deriving it from the display name for older API payloads.
  selectedTable.value =
    item.delete_key || item.table_type.replace(/\s+/g, '').toLowerCase()

  availableRiskRateCodes.value = []
  getRiskRateCodes()
  riskCodeDialog.value = true
}

// Delete button — routes through the standard picker. For Age Bands the
// picker shows distinct types (the backend returns them under the same
// risk-code slot) so the user can delete one type at a time.
const deleteClicked = (item: any) => {
  chooseRiskRateCode(item)
}

const pickerLabel = computed(() =>
  selectedTableText.value === 'Age Bands'
    ? 'Choose the Age Band Type to delete'
    : 'Choose the relevant Risk Code'
)
const pickerFieldLabel = computed(() =>
  selectedTableText.value === 'Age Bands'
    ? 'Select an age band type'
    : 'Select a risk rate code'
)

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
    }
  } catch (error: any) {
    let errorMessage = 'An error occurred while uploading the file'

    console.log('Upload error details:', error.response?.data)

    if (error.response?.status === 400 && error.response?.data?.error) {
      errorMessage = error.response.data.error
    }

    snackbarText.value = errorMessage
    snackbar.value = true
    timeout.value = 3000
    uploadComplete.value = true
  } finally {
    // update generic risk rate codes
    getGenericRiskRateCodes()
  }

  // Refresh metadata first, then recompute the empty-banner state so it
  // reflects the row just uploaded (not the stale pre-upload snapshot).
  const res = await GroupPricingService.getTableMetaData()
  tables.value = res.data.associated_tables
  checkPopulated()
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
      }
    }
  } catch (error) {
    console.log(error)
  }
  // Refresh metadata first so checkPopulated() sees the latest populated
  // flags; otherwise the banner would reflect the pre-delete snapshot.
  const res = await GroupPricingService.getTableMetaData()
  tables.value = res.data.associated_tables
  checkPopulated()
}

const viewTable = (item: any) => {
  loadingData.value = true
  tableData.value = []

  // Prefer the authoritative delete_key from the backend (matches the
  // switch-case slug on the server); fall back to deriving it from the
  // display name for older API payloads.
  const tableType =
    item.delete_key || item.table_type.replace(/ /g, '').toLowerCase()
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
        createColumnDefs(tableData.value, item.table_type)
      }
      selectedTable.value = item.table_type
    }
    loadingData.value = false
  })
}

// Per-table explicit column ordering. Keys not listed fall back to the order
// returned by the backend (which for tables serialised via map[string]interface{}
// is alphabetical). Unknown fields in a preferred list are skipped; remaining
// fields are appended after the preferred ones in their original order.
const preferredColumnOrder: Record<string, string[]> = {
  'Age Bands': ['name', 'type', 'min_age', 'max_age'],
  'Tax Retirement': [
    'risk_rate_code',
    'lower_bound',
    'upper_bound',
    'tax_rate',
    'cumulative_tax_relief'
  ]
}

// System fields that should never render as columns.
const hiddenColumns = new Set(['id'])

const createColumnDefs = (data: any, tableType?: string) => {
  columnDefs.value = []
  const keys = Object.keys(data[0]).filter((k) => !hiddenColumns.has(k))
  const preferred = tableType ? preferredColumnOrder[tableType] : undefined
  const ordered = preferred
    ? [
        ...preferred.filter((k) => keys.includes(k)),
        ...keys.filter((k) => !preferred.includes(k))
      ]
    : keys
  ordered.forEach((element) => {
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
