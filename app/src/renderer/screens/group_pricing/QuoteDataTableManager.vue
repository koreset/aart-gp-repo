<template>
  <experience-rate-overrides
    v-if="quote && quote.experience_rating === 'Override'"
    :quote="quote"
    :result-summaries="resultSummaries"
    class="mb-4"
    @overrides-updated="handleOverridesUpdated"
  />
  <base-card :show-actions="false">
    <template #header>
      <span class="headline">Input Data Tables</span>
    </template>
    <p
      v-if="indicativeDataEnabled"
      class="text-body-2 text-medium-emphasis mt-2 mb-2"
    >
      Indicative data lets you price a quote without uploading a member file by
      entering category-level summary statistics. Premiums are then derived from
      category salary multiples rather than per-member salaries.
    </p>
    <v-alert
      v-if="indicativeDataEnabled && !props.quote.use_global_salary_multiple"
      type="info"
      variant="tonal"
      density="compact"
      prepend-icon="mdi-information-outline"
      class="mb-3"
    >
      To use indicative data, enable <strong>Global Salary Multiple</strong> on
      this quote. Without it, benefit amounts cannot be computed from
      category-level averages.
      <template #append>
        <v-btn
          size="small"
          variant="text"
          color="primary"
          append-icon="mdi-arrow-right"
          @click="goToGeneralInput"
        >
          Open General Input
        </v-btn>
      </template>
    </v-alert>
    <v-row v-if="indicativeDataEnabled" class="mt-1" align="start" no-gutters>
      <v-col cols="12" md="4" class="pe-md-2 pb-2">
        <v-select
          v-model="selectedCategory"
          variant="outlined"
          density="compact"
          label="Category"
          clearable
          :items="availableCategories"
        ></v-select>
      </v-col>
      <v-col cols="12" md="4" class="pe-md-2 pb-2">
        <v-text-field
          v-model="memberAverageAge"
          label="Average Age"
          type="number"
          density="compact"
          variant="outlined"
          hide-details
        ></v-text-field>
      </v-col>
      <v-col cols="12" md="4" class="pe-md-2 pb-2">
        <v-text-field
          v-model="memberAverageIncome"
          label="Average Annual Salary"
          type="number"
          density="compact"
          variant="outlined"
          hide-details="auto"
          hint="Annual gross salary, averaged across members in the category"
          persistent-hint
        ></v-text-field>
      </v-col>
      <v-col cols="12" md="4" class="pe-md-2 pb-2">
        <v-text-field
          v-model="memberDataCount"
          label="Member Count"
          type="number"
          density="compact"
          variant="outlined"
          hide-details
        ></v-text-field>
      </v-col>
      <v-col cols="12" md="4" class="pe-md-2 pb-2">
        <v-text-field
          v-model="memberMaleFemaleDistribution"
          label="Male Proportion (%)"
          density="compact"
          type="number"
          variant="outlined"
          hide-details="auto"
          hint="Percentage of male members (0–100). Female proportion is the remainder."
          persistent-hint
          :rules="[malePropRule]"
          min="0"
          max="100"
          suffix="%"
        ></v-text-field>
      </v-col>
      <v-col cols="12" md="4" class="pe-md-2 pb-2">
        <v-btn rounded color="primary" size="small" @click="addToDataSet"
          >Add / Update</v-btn
        >
      </v-col>
    </v-row>
    <v-row v-if="indicativeDataEnabled">
      <v-col cols="12">
        <v-table v-if="indicativeDataSet.length > 0" density="compact">
          <thead>
            <tr>
              <th class="text-left">Category</th>
              <th class="text-left">Average Age</th>
              <th class="text-left">Average Annual Salary</th>
              <th class="text-left">Member Count</th>
              <th class="text-left">Male Proportion (%)</th>
              <th class="text-left">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in indicativeDataSet" :key="index">
              <td>{{ item.scheme_category }}</td>
              <td>{{ item.member_average_age }}</td>
              <td>{{ item.member_average_income }}</td>
              <td>{{ item.member_data_count }}</td>
              <td>{{ item.member_male_female_distribution }}</td>
              <td>
                <v-btn
                  size="small"
                  color="error"
                  variant="text"
                  @click="removeIndicativeData(index)"
                >
                  <v-icon>mdi-delete</v-icon>
                </v-btn>
              </td>
            </tr>
          </tbody>
        </v-table>
      </v-col>
    </v-row>
    <v-row v-if="indicativeDataEnabled && indicativeDataSet.length > 0">
      <v-col cols="12" md="4">
        <v-btn
          rounded
          color="primary"
          size="small"
          @click="uploadIndicativeData"
          >Save Data</v-btn
        >
        <v-btn
          rounded
          class="ml-4"
          color="red"
          size="small"
          @click="deleteIndicativeData"
          >Delete Data</v-btn
        >
      </v-col>
    </v-row>
    <v-divider v-if="indicativeDataEnabled" class="my-4"></v-divider>
    <v-list lines="two">
      <v-list-item v-for="item in relatedTables" :key="item.value">
        <template #prepend>
          <v-icon :color="item.populated ? 'success' : 'error'" class="mr-4">
            {{
              item.populated ? 'mdi-check-circle' : 'mdi-alert-circle-outline'
            }}
          </v-icon>
        </template>
        <v-list-item-title>{{ item.table_type }}</v-list-item-title>
        <v-list-item-subtitle
          v-if="item.table_type !== 'Group Pricing Parameters'"
          :class="item.populated ? '' : 'text-error'"
        >
          <span v-if="item.table_type === 'Member Data' && item.isIndicative">
            Using indicative data
          </span>
          <span v-else>
            {{
              item.populated ? `${item.count} records loaded` : 'Data required'
            }}
          </span>
          <span
            v-if="
              item.table_type === 'Member Data' &&
              item.populated &&
              !item.isIndicative &&
              memberGenderSplit &&
              memberGenderSplit.total_count > 0
            "
            class="ml-2 text-medium-emphasis"
          >
            · Male {{ memberGenderSplit.male_count }} ({{
              (memberGenderSplit.male_prop * 100).toFixed(1)
            }}%) · Female {{ memberGenderSplit.female_count }} ({{
              (memberGenderSplit.female_prop * 100).toFixed(1)
            }}%)<span v-if="memberGenderSplit.other_count > 0">
              · Other {{ memberGenderSplit.other_count }}</span
            >
          </span>
        </v-list-item-subtitle>

        <template #append>
          <div
            v-if="
              item.table_type === 'Member Data' &&
              quote.quote_type === 'New Business'
            "
            class="d-flex align-center mr-4"
          >
            <v-checkbox
              v-model="indicativeDataEnabled"
              label="Indicative Data"
              hide-details
              density="compact"
              @update:modelValue="toggleIndicativeProp"
            ></v-checkbox>
          </div>
          <v-btn
            :disabled="
              !item.populated ||
              (item.table_type === 'Member Data' && indicativeDataEnabled)
            "
            rounded
            color="primary"
            class="mr-2"
            variant="outlined"
            size="small"
            @click="viewTable(item)"
          >
            <v-icon left color="primary">mdi-information</v-icon>
            <span>View</span></v-btn
          >
          <v-btn
            v-if="item.table_type !== 'Group Pricing Parameters'"
            :disabled="
              quote.status === 'accepted' ||
              quote.status === 'in_force' ||
              (item.table_type === 'Member Data' && indicativeDataEnabled)
            "
            class="mr-2"
            variant="outlined"
            size="small"
            rounded
            color="primary"
            @click="openUploadDialog(item)"
          >
            <v-icon color="accent">mdi-upload</v-icon>
            <span>Upload</span></v-btn
          >
          <v-btn
            v-if="item.table_type !== 'Group Pricing Parameters'"
            :disabled="isDeleteButtonDisabled(item)"
            variant="outlined"
            color="error"
            rounded
            size="small"
            @click="deleteTable(item)"
          >
            <v-icon color="accent">mdi-delete</v-icon>
            <span>Delete</span></v-btn
          >
        </template>
        <!-- Indicative Data Fields for Member Data -->
      </v-list-item>
    </v-list>

    <loading-indicator :loadingData="loadingData || uploadingData" />

    <v-row
      v-if="
        (resultTableData.length > 0 || useInfiniteModel) &&
        !loadingData &&
        selectedTable
      "
    >
      <v-col>
        <group-pricing-data-grid
          ref="dataGridRef"
          :columnDefs="columnDefs"
          :show-close-button="true"
          :rowData="useInfiniteModel ? undefined : resultTableData"
          :table-title="tableTitle"
          :pagination="!useInfiniteModel"
          :rowCount="rowCount"
          :useInfiniteModel="useInfiniteModel"
          :dataSource="dataSource"
          :show-export-all-csv="shouldExportAll"
          @update:clear-data="clearData"
          @update-column="handleColumnUpdate"
          @export-all-csv="handleExportAllCsv"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <file-upload-dialog
          :yearLabel="yearLabel"
          :isDialogOpen="isDialogOpen"
          :showModelPoint="showModelPoint"
          :mpLabel="mpLabel"
          :table="'undefined'"
          :uploadTitle="uploadTitle"
          :years="years"
          @upload="handleUpload"
          @update:isDialogOpen="updateDialog"
        />
      </v-col>
    </v-row>

    <!-- Snackbar for notifications -->
    <v-snackbar
      v-model="snackbar"
      :timeout="4000"
      color="primary"
      location="bottom"
    >
      {{ snackbarText }}
      <template #actions>
        <v-btn color="white" variant="text" @click="snackbar = false">
          Close
        </v-btn>
      </template>
    </v-snackbar>
  </base-card>
  <confirm-dialog ref="confirmationDialog" />
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import ConfirmDialog from '@/renderer/components/ConfirmDialog.vue'
// import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import formatValues from '@/renderer/utils/format_values'
import _ from 'lodash'
import FileUploadDialog from '@/renderer/components/FileUploadDialog.vue'
import LoadingIndicator from '@/renderer/components/LoadingIndicator.vue'
import GroupPricingDataGrid from '@/renderer/components/tables/GroupPricingDataGrid.vue'
import ExperienceRateOverrides from '@/renderer/components/grouppricing/ExperienceRateOverrides.vue'

// Import necessary components like BaseCard, DataGrid, FileUploadDialog

// const props = defineProps({
//   quoteId: { type: String, required: true },
//   quoteStatus: { type: String, required: true }
// })

const props = defineProps({
  quote: {
    type: Object,
    required: true
  },
  resultSummaries: {
    type: Array as () => any[],
    default: () => []
  }
})

const emit = defineEmits([
  'quote-updated',
  'indicative-data-updated',
  'navigate-to-general-input'
])
const confirmationDialog: any = ref(null)
const tableData = ref([])
const columnDefs: any = ref([])
const selectedTable: any = ref(null)
const loadingData = ref(false)
const isDialogOpen = ref(false)
const resultTableData = ref([])
const displaySummary = ref(false)
const snackbar = ref(false)
const snackbarText = ref('')
const glaBenefitTitle = ref('Group Life Assurance (GLA)')
const sglaBenefitTitle = ref('Spouse Group Life Assurance (GLA)')
const ptdBenefitTitle = ref('Permanent Total Disability')
const ciBenefitTitle = ref('Critical Illness')
const phiBenefitTitle = ref('Personal Health Insurance')
const ttdBenefitTitle = ref('Temporary Total Disability')
const familyFuneralBenefitTitle = ref('Family Funeral')
const benefitMaps: any = ref([])
const yearLabel = ref('') // 'Select a year'
const uploadTitle = ref('')
const mpLabel = ref('')
const showModelPoint = ref(false)
const indicativeDataEnabled = ref(false)
const indicativeDataSet = ref<any[]>([])
const memberAverageAge = ref<number | null>(null)
const memberAverageIncome = ref<number | null>(null)
const memberDataCount = ref<number | null>(null)
const memberMaleFemaleDistribution = ref<number | null>(null)
const malePropRule = (v: number | string | null) => {
  if (v === null || v === '') return true
  const n = Number(v)
  return (n >= 0 && n <= 100) || 'Enter a value between 0 and 100'
}
const goToGeneralInput = () => {
  emit('navigate-to-general-input')
}
const memberGenderSplit = ref<{
  male_count: number
  female_count: number
  other_count: number
  total_count: number
  male_prop: number
  female_prop: number
} | null>(null)

async function refreshMemberGenderSplit() {
  if (!props.quote || !props.quote.quote_id) {
    memberGenderSplit.value = null
    return
  }
  if (!(props.quote.member_data_count && props.quote.member_data_count > 0)) {
    memberGenderSplit.value = null
    return
  }
  try {
    const res = await GroupPricingService.getQuoteMemberGenderSplit(
      props.quote.quote_id
    )
    if (res?.data?.success) {
      memberGenderSplit.value = res.data.data
    }
  } catch {
    memberGenderSplit.value = null
  }
}
const selectedCategory = ref<string | null>(null)
const useInfiniteModel = ref(false)
const currentTableType = ref('')
const dataSource: any = ref(null)
const dataGridRef: any = ref(null)
const rowCount: any = ref(0)
const shouldExportAll = ref(true)
const uploadingData = ref(false)

// const availableCategories = ref<string[]>([])

const years = ref<number[]>(
  Array.from({ length: 10 }, (v, k) => new Date().getFullYear() - k)
)
const updateDialog = (value: boolean) => {
  isDialogOpen.value = value
}

const toggleIndicativeProp = () => {
  emit('indicative-data-updated', {
    tableType: 'Member Data',
    indicativeData: indicativeDataEnabled.value
  })
}

const handleOverridesUpdated = (count: number) => {
  // Re-fetch the quote so experience_rate_overrides_count flows into the
  // Run Calculations gate without requiring a full page reload.
  emit('quote-updated', {
    tableType: 'Experience Rate Overrides',
    experienceRateOverridesCount: count
  })
}

const clearData = () => {
  tableData.value = []
  resultTableData.value = []
  selectedTable.value = ''
  useInfiniteModel.value = false
  dataSource.value = null
  currentTableType.value = ''
}

const isDeleteButtonDisabled = (item: any) => {
  return (
    !item.populated ||
    props.quote.status === 'accepted' ||
    props.quote.status === 'in_force' ||
    props.quote.quote_type === 'Renewal' ||
    (item.table_type === 'Member Data' && indicativeDataEnabled.value)
  )
}

const handleExportAllCsv = async () => {
  try {
    const response = await GroupPricingService.exportQuoteTableCsv(
      props.quote.id,
      currentTableType.value
    )

    // Create download link
    const fileURL = window.URL.createObjectURL(new Blob([response.data]))
    const fileLink = document.createElement('a')

    fileLink.href = fileURL
    fileLink.setAttribute(
      'download',
      `${currentTableType.value}_${props.quote.id}.xlsx`
    )
    document.body.appendChild(fileLink)

    fileLink.click()

    // Clean up
    document.body.removeChild(fileLink)
    window.URL.revokeObjectURL(fileURL)
  } catch (error) {
    snackbarText.value = 'Error exporting CSV file'
    snackbar.value = true
  } finally {
    // Reset the loader in the data grid component
    if (dataGridRef.value && dataGridRef.value.resetExportAllCsvLoader) {
      dataGridRef.value.resetExportAllCsvLoader()
    }
  }
}
const handleColumnUpdate = async ({ colId, newValue }) => {
  // your existing logic...
  // on success, emit event to refresh parent
  emit('quote-updated')
}

const relatedTables = computed(() => {
  const tables: any = []
  if (!props.quote) {
    return tables
  }
  if (props.quote.member_data_count > 0 || indicativeDataEnabled.value) {
    tables.push({
      table_type: 'Member Data',
      value: 'member_data',
      populated: true,
      count: props.quote.member_data_count,
      isIndicative: indicativeDataEnabled.value
    })
  } else {
    tables.push({
      table_type: 'Member Data',
      value: 'member_data',
      populated: false,
      count: 0,
      isIndicative: false
    })
  }
  // Only show Claims Experience table if experience rating is enabled
  if (
    props.quote.experience_rating === 'Yes' &&
    indicativeDataEnabled.value === false
  ) {
    if (props.quote.claims_experience_count > 0) {
      tables.push({
        table_type: 'Claims Experience',
        value: 'claims_experience',
        populated: true,
        count: props.quote.claims_experience_count
      })
    } else {
      tables.push({
        table_type: 'Claims Experience',
        value: 'claims_experience',
        populated: false,
        count: 0
      })
    }
  }

  tables.push({
    table_type: 'Group Pricing Parameters',
    value: 'group_pricing_parameters',
    populated: true
  })

  return tables
})

const availableCategories = computed(() => {
  const categories: string[] = []

  // get scheme categories from the quote's schemes
  if (
    props.quote &&
    props.quote.selected_scheme_categories &&
    props.quote.selected_scheme_categories.length > 0
  ) {
    props.quote.selected_scheme_categories.forEach((category) => {
      if (category && !categories.includes(category)) {
        categories.push(category)
      }
    })
  }
  return categories
})

const deleteTable = async (item) => {
  const res = await confirmationDialog.value.open(
    'Confirm Deletion',
    `Are you sure you want to delete the ${item.table_type} table? This action cannot be undone.`
  )

  if (!res) {
    return
  }
  // your delete logic here...
  // on success:
  GroupPricingService.deleteQuoteTable(props.quote.id, item.table_type)
    .then((res) => {
      snackbarText.value = 'Table deleted successfully'
      snackbar.value = true
      emit('quote-updated', {
        tableType: item.table_type,
        count: 0,
        updateType: 'delete'
      })
    })
    .catch(() => {
      snackbarText.value = 'Failed to delete table'
      snackbar.value = true
    })
}
const addToDataSet = () => {
  if (!selectedCategory.value) {
    snackbarText.value = 'Please select a category before adding to data set'
    snackbar.value = true
    return
  }

  const indicativeData = {
    quote_id: props.quote.id,
    scheme_category: selectedCategory.value,
    member_average_age: memberAverageAge.value
      ? Number(memberAverageAge.value)
      : null,
    member_average_income: memberAverageIncome.value
      ? Number(memberAverageIncome.value)
      : null,
    member_data_count: memberDataCount.value
      ? Number(memberDataCount.value)
      : null,
    member_male_female_distribution: memberMaleFemaleDistribution.value
      ? Number(memberMaleFemaleDistribution.value)
      : null
  }

  // if indicativeData with the selectedCategory already exists in indicativeDataSet, update it
  const existingIndex = indicativeDataSet.value.findIndex(
    (data) => data.scheme_category === selectedCategory.value
  )
  if (existingIndex !== -1) {
    indicativeDataSet.value[existingIndex] = indicativeData
  } else {
    indicativeDataSet.value.push(indicativeData)
  }

  // console.log('Adding to Data Set:', indicativeDataSet.value)
  // reset the input fields
  selectedCategory.value = null
  memberAverageAge.value = null
  memberAverageIncome.value = null
  memberDataCount.value = null
  memberMaleFemaleDistribution.value = null
}

const uploadIndicativeData = async () => {
  if (indicativeDataSet.value.length === 0) {
    snackbarText.value = 'No indicative data to upload'
    snackbar.value = true
    return
  }
  if (
    indicativeDataSet.value.length <
    props.quote.selected_scheme_categories.length
  ) {
    snackbarText.value =
      'Please add indicative data for all selected scheme categories before saving'
    snackbar.value = true
    return
  }

  GroupPricingService.sendIndicativeMemberData(indicativeDataSet.value)
    .then((res) => {
      snackbarText.value = 'Indicative data uploaded successfully'
      snackbar.value = true
      emit('quote-updated', {
        tableType: 'Member Data',
        updateType: 'indicative',
        indicativeData: true,
        indicativeDataSet: [...indicativeDataSet.value]
      })
    })
    .catch(() => {
      snackbarText.value = 'Failed to upload indicative data'
      snackbar.value = true
    })
}

const removeIndicativeData = async (index: number) => {
  // await uploadIndicativeData()
  indicativeDataSet.value.splice(index, 1)
  snackbarText.value = 'Indicative data entry removed'
  snackbar.value = true
}

const deleteIndicativeData = async () => {
  const res = await confirmationDialog.value.open(
    'Confirm Deletion',
    `Are you sure you want to delete all indicative data? This action cannot be undone.`
  )

  if (!res) {
    return
  }

  GroupPricingService.deleteIndicativeMemberData(props.quote.id)
    .then((res) => {
      indicativeDataSet.value = []
      snackbarText.value = 'Indicative data deleted successfully'
      snackbar.value = true
      emit('quote-updated', {
        tableType: 'Member Data',
        updateType: 'delete_indicative',
        indicativeData: false,
        indicativeDataSet: []
      })
    })
    .catch(() => {
      snackbarText.value = 'Failed to delete indicative data'
      snackbar.value = true
    })
}

const openUploadDialog = (item: any) => {
  if (props.quote.status === 'accepted' || props.quote.status === 'in_force') {
    snackbarText.value = 'Cannot upload data to accepted or in-force quotes'
    snackbar.value = true
    return
  }
  selectedTable.value = item
  yearLabel.value = 'Select a year'
  uploadTitle.value = 'Upload Data for ' + item.table_type + ' Table (csv)'
  isDialogOpen.value = true
}

const handleUpload = async (payload: any) => {
  uploadingData.value = true
  const formdata = new FormData()
  formdata.append('file', payload.file)
  formdata.append('quote_id', props.quote.id)
  formdata.append('table_type', selectedTable.value.table_type)

  if (!payload?.file) {
    snackbarText.value = 'Choose file to upload'
    snackbar.value = true
    return
  }

  if (payload?.file.size <= 0) {
    snackbarText.value = 'Cannot upload an empty file'
    snackbar.value = true
    return
  }

  GroupPricingService.uploadQuoteTable(formdata)
    .then((res) => {
      const count = res.data

      snackbarText.value = 'Upload Successful'
      snackbar.value = true
      emit('quote-updated', {
        tableType: selectedTable.value.table_type,
        count,
        updateType: 'upload'
      })
      if (selectedTable.value?.table_type === 'Member Data') {
        // Refresh the gender split readout shown under the Member Data row
        // as soon as the upload completes.
        refreshMemberGenderSplit()
      }
      isDialogOpen.value = false
      uploadingData.value = false
    })
    .catch((error) => {
      // console.log('Error:', error)
      const errMsg = error.response?.data?.error
      snackbarText.value =
        errMsg && errMsg.trim() !== '' ? errMsg : 'Failed to upload table'
      snackbar.value = true
      uploadingData.value = false
    })
}
const tableTitle = ref('')
const viewTable = async (item: any) => {
  // special case for output summary
  loadingData.value = true
  tableTitle.value = item.table_type
  try {
    currentTableType.value = item.value

    // Check if we should use infinite model (for large datasets)
    const shouldUseInfiniteModel = item.count > 10000
    useInfiniteModel.value = shouldUseInfiniteModel

    if (shouldUseInfiniteModel) {
      // Create datasource for infinite row model

      // First, fetch a small sample to get column definitions
      const sampleRes = await GroupPricingService.getQuoteTable(
        props.quote.id,
        item.value,
        {
          offset: 0,
          limit: 1
        }
      )

      if (sampleRes.data.data && sampleRes.data.data.length > 0) {
        const sampleData = sampleRes.data.data.map((item) =>
          _.fromPairs(sampleRes.data.json_tags.map((key) => [key, item[key]]))
        )
        createColumnDefs(sampleData)
      }

      dataSource.value = {
        getRows: async (params: any) => {
          try {
            const startRow = params.request.startRow
            const endRow = params.request.endRow
            const limit = endRow - startRow

            const res = await GroupPricingService.getQuoteTable(
              props.quote.id,
              currentTableType.value,
              { offset: startRow, limit }
            )
            rowCount.value = endRow
            if (res.data.data !== null && res.data.data.length > 0) {
              const orderedData = res.data.data.map((item) =>
                _.fromPairs(res.data.json_tags.map((key) => [key, item[key]]))
              )

              // Calculate last row - if we got less data than requested, we've reached the end
              const lastRow =
                res.data.data.length < limit
                  ? startRow + res.data.data.length
                  : -1
              params.success({ rowData: orderedData, lastRow })
            } else {
              params.success([], 0)
            }
          } catch (error) {
            params.fail()
          }
        }
      }
      selectedTable.value = item.table_type
      displaySummary.value = false
      resultTableData.value = [] // Clear client-side data when using infinite model
      loadingData.value = false
    } else {
      // Use original client-side approach for smaller datasets
      const res = await GroupPricingService.getQuoteTable(
        props.quote.id,
        item.value
      )
      if (res.data.data !== null && res.data.data.length > 0) {
        const orderedData = res.data.data.map((item) =>
          _.fromPairs(res.data.json_tags.map((key) => [key, item[key]]))
        )

        displaySummary.value = false

        if (
          item.value === 'member_data' ||
          item.value === 'group_pricing_parameters'
        ) {
          resultTableData.value = orderedData
        } else {
          tableData.value = orderedData
        }

        selectedTable.value = item.table_type
        createColumnDefs(orderedData)
        loadingData.value = false
      } else {
        tableData.value = []
        resultTableData.value = []
        selectedTable.value = ''
        snackbarText.value = 'No data found for this table'
        snackbar.value = true
        loadingData.value = false
      }
    }
  } catch (error) {
    loadingData.value = false
  }
}
//     if (res.data.data !== null && res.data.data.length > 0) {
//       const orderedData = res.data.data.map((item) =>
//         _.fromPairs(res.data.json_tags.map((key) => [key, item[key]]))
//       )

//       displaySummary.value = false

//       if (
//         item.value === 'member_rating_results' ||
//         item.value === 'member_premium_schedules' ||
//         item.value === 'bordereaux'
//       ) {
//         resultTableData.value = orderedData
//       } else {
//         tableData.value = orderedData
//       }

//       selectedTable.value = item.table_type
//       createColumnDefs(orderedData)
//       loadingData.value = false
//     } else {
//       tableData.value = []
//       resultTableData.value = []
//       selectedTable.value = ''
//       snackbarText.value = 'No data found for this table'
//       snackbar.value = true
//     }
//   } catch (error) {
//     console.log('Error:', error)
//   }
// }

onMounted(() => {
  refreshMemberGenderSplit()
  GroupPricingService.getBenefitMaps().then((res) => {
    benefitMaps.value = res.data
    const glaBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'GLA'
    )
    if (glaBenefit.benefit_alias !== '') {
      glaBenefitTitle.value = glaBenefit.benefit_alias
    }
    const sglaBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'SGLA'
    )
    if (sglaBenefit.benefit_alias !== '') {
      sglaBenefitTitle.value = sglaBenefit.benefit_alias
    }
    const ptdBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'PTD'
    )
    if (ptdBenefit.benefit_alias !== '') {
      ptdBenefitTitle.value = ptdBenefit.benefit_alias
    }
    const ciBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'CI'
    )
    if (ciBenefit.benefit_alias !== '') {
      ciBenefitTitle.value = ciBenefit.benefit_alias
    }
    const phiBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'PHI'
    )
    if (phiBenefit.benefit_alias !== '') {
      phiBenefitTitle.value = phiBenefit.benefit_alias
    }
    const ttdBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'TTD'
    )
    if (ttdBenefit.benefit_alias !== '') {
      ttdBenefitTitle.value = ttdBenefit.benefit_alias
    }
    const familyFuneralBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'GFF'
    )
    if (familyFuneralBenefit.benefit_alias !== '') {
      familyFuneralBenefitTitle.value = familyFuneralBenefit.benefit_alias
    }
  })

  // load indicative data if available
  indicativeDataEnabled.value = props.quote.member_indicative_data
  if (indicativeDataEnabled.value) {
    indicativeDataSet.value = props.quote.member_indicative_data_set || []
  }
})

const replaceSubstring = (input, search, replacement) => {
  if (input.includes(search)) {
    benefitMaps.value.forEach((item) => {
      if (
        item.benefit_code.toLowerCase() === search &&
        item.benefit_alias_code !== ''
      ) {
        input = input.replace(search, item.benefit_alias_code.toLowerCase())
      }
    })
    // return input
  }
  return input
}

const createColumnDefs = (data: any) => {
  columnDefs.value = []
  Object.keys(data[0]).forEach((element) => {
    const column: any = {}
    if (element.startsWith('ci')) {
      column.headerName = replaceSubstring(element, 'ci', 'severe')
    } else if (element.startsWith('phi')) {
      column.headerName = replaceSubstring(element, 'phi', 'phi2')
    } else if (element.startsWith('gla')) {
      column.headerName = replaceSubstring(element, 'gla', 'gla2')
    } else if (element.startsWith('sgla')) {
      column.headerName = replaceSubstring(element, 'sgla', 'sgla2')
    } else if (element.startsWith('ptd')) {
      column.headerName = replaceSubstring(element, 'ptd', 'ptd2')
    } else if (element.startsWith('ttd')) {
      column.headerName = replaceSubstring(element, 'ttd', 'ttd2')
    } else if (element.startsWith('gff')) {
      column.headerName = replaceSubstring(element, 'gff', 'gff2')
    } else {
      column.headerName = element
    }

    // header.headerName = element
    column.field = element
    column.valueFormatter = formatValues
    column.minWidth = 200
    column.sortable = true
    column.filter = true
    column.resizable = true
    if (column.field === 'id' || column.field === 'quote_id') {
      column.hide = true
    }
    if (column.field === 'benefits') {
      column.valueFormatter = (params: any) => {
        const b = params.value
        if (!b || typeof b !== 'object') return ''
        const parts: string[] = []
        const add = (code: string, enabled: any, mult: any) => {
          const m = Number(mult ?? 0)
          if (enabled || m > 0) parts.push(`${code} ${m}x`)
        }
        add('GLA', b.gla_enabled, b.gla_multiple)
        add('SGLA', b.sgla_enabled, b.sgla_multiple)
        add('PTD', b.ptd_enabled, b.ptd_multiple)
        add('CI', b.ci_enabled, b.ci_multiple)
        add('TTD', b.ttd_enabled, b.ttd_multiple)
        add('PHI', b.phi_enabled, b.phi_multiple)
        if (b.gff_enabled) parts.push('GFF')
        return parts.length ? parts.join(' · ') : '—'
      }
      column.minWidth = 280
      column.filter = false
    }
    // if column.field is exp_credibility, then set it to editable
    if (column.field === 'manually_added_credibility') {
      column.editable = true
    } else {
      column.editable = false
    }
    columnDefs.value.push(column)
  })
}
</script>
