<template>
  <v-dialog v-model="dialog" persistent max-width="550px">
    <template #activator="{ props: slotProps }">
      <v-btn
        :disabled="disabled"
        variant="outlined"
        rounded
        size="small"
        v-bind="slotProps"
      >
        <v-icon color="accent">mdi-upload</v-icon>
        <span>{{ actionName }}</span>
      </v-btn>
    </template>
    <base-card v-if="!uploadSuccess">
      <template #header>
        <span class="headline">Upload data for {{ tableType }} (csv)</span>
      </template>
      <template #default>
        <v-row>
          <v-col>
            <v-file-input
              v-model="file"
              variant="outlined"
              density="compact"
              show-size
              accept=".csv"
              placeholder="Click to choose a file"
              :label="tableType + ' File'"
              prepend-icon="mdi-paperclip"
            >
            </v-file-input>
          </v-col>
        </v-row>
        <v-row v-if="showYieldCurveCode && tableType === 'Yield Curve'">
          <v-col>
            <v-text-field
              v-model="yieldCurveCode"
              variant="outlined"
              density="compact"
              placeholder="Enter a yield curve code"
              label="Yield Curve Code"
            ></v-text-field>
          </v-col>
        </v-row>
        <v-row v-if="showYear">
          <v-col>
            <v-select
              v-model="selectedYear"
              variant="outlined"
              density="compact"
              label="Select a year for this data set"
              :items="availableYears"
              item-title="year"
              item-value="year"
            ></v-select>
          </v-col>
        </v-row>
        <v-row v-if="showVersion">
          <v-col>
            <v-text-field
              v-model="dataVersion"
              variant="outlined"
              density="compact"
              placeholder="Enter a version name"
              label="Version Name"
            ></v-text-field>
          </v-col>
        </v-row>
        <v-row v-if="showRiskRateCode">
          <v-col>
            <v-text-field
              v-model="dataRiskRateCode"
              variant="outlined"
              density="compact"
              placeholder="Enter a risk rate code"
              label="Risk Rate Code"
            ></v-text-field>
          </v-col>
        </v-row>
        <v-row v-if="showRiskRateCodeSelect">
          <v-col>
            <v-select
              v-model="dataRiskRateCode"
              variant="outlined"
              density="compact"
              label="Select a risk rate code"
              :items="gpRiskRateCodes"
            >
            </v-select>
          </v-col>
        </v-row>
        <v-row v-if="tableType === 'Custom Tiered Income Replacement'">
          <v-col>
            <v-text-field
              v-model="dataSchemeName"
              variant="outlined"
              density="compact"
              placeholder="Enter the scheme name"
              label="Scheme Name"
            ></v-text-field>
          </v-col>
        </v-row>

        <v-row>
          <v-col v-if="showMonth && tableType === 'Yield Curve'">
            <v-select
              v-model="selectedMonth"
              variant="outlined"
              density="compact"
              label="Select a month for this yield curve table"
              :items="availableMonths"
              item-title="month"
              item-value="month"
            ></v-select>
          </v-col>
        </v-row>
        <v-row v-if="showExpTableTypes">
          <v-col>
            <v-select
              v-model="selectedTableType"
              variant="outlined"
              density="compact"
              label="Select a table type"
              :items="expTableTypes"
              item-title="name"
              item-value="value"
            ></v-select>
          </v-col>
        </v-row>
        <v-row v-if="showIbnrTableTypes">
          <v-col>
            <v-select
              v-model="selectedTableType"
              variant="outlined"
              density="compact"
              label="Select a table type"
              :items="ibnrTableTypes"
              item-title="name"
              item-value="value"
            ></v-select>
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-btn
              :loading="uploadInProgress"
              :disabled="uploadDisabled"
              rounded
              size="small"
              color="primary"
              class="primary mr-7"
              @click="uploadFile"
              >Upload</v-btn
            >
            <v-progress-circular
              class="mr-7"
              indeterminate
              :size="loaderSize"
              :width="2"
              color="primary"
            ></v-progress-circular>
          </v-col>
        </v-row>
      </template>
      <template #actions>
        <v-spacer></v-spacer>
        <v-btn rounded variant="text" @click="closeDialog">Close</v-btn>
      </template>
    </base-card>
    <base-card v-else>
      <template #header>
        <span class="headline">Upload Status for {{ fileName }}</span>
      </template>
      <template #default>
        <h4>Table was successfully uploaded</h4>
      </template>
      <template #actions>
        <v-spacer></v-spacer>
        <v-btn rounded variant="text" @click="closeDialog">Close</v-btn>
      </template>
    </base-card>
  </v-dialog>
</template>
<script setup lang="ts">
import { computed, defineProps, ref, watch } from 'vue'
import { DataPayload } from './types'
import BaseCard from './BaseCard.vue'

const props = defineProps({
  disabled: { type: Boolean, default: false },
  tableType: { type: String, default: '' },
  gpRiskRateCodes: { type: Array, default: () => [] },
  actionName: { type: String, default: 'Upload' },
  showYear: { type: Boolean, default: false },
  showMonth: { type: Boolean, default: false },
  showRiskRateCode: { type: Boolean, default: false },
  showRiskRateCodeSelect: { type: Boolean, default: false },
  showVersion: { type: Boolean, default: false },
  showExpTableTypes: { type: Boolean, default: false },
  showIbnrTableTypes: { type: Boolean, default: false },
  showYieldCurveCode: { type: Boolean, default: false },
  showPaaModelPoints: { type: Boolean, default: false },
  assumptionType: { type: String, default: '' },
  uploadComplete: { type: Boolean, default: false }
})

const emit = defineEmits<{
  (e: 'uploadFile', payload: DataPayload): void
}>()

// data
const dataVersion: any = ref(null)
const dataRiskRateCode: any = ref(null)
const dataSchemeName: any = ref(null)
const uploadInProgress: any = ref(false)
const uploadDisabled: any = ref(false)
const loaderSize: any = ref(0)
const uploadSuccess: any = ref(false)
const file: any = ref(null)
const dialog: any = ref(false)
const selectedYear: any = ref(null)
const selectedMonth: any = ref(null)
const yieldCurveCode: any = ref(null)
const selectedTableType: any = ref(null)
const fileName: any = ref(null)
const expTableTypes = [
  { name: 'Actual Data', value: 'actual_data' },
  { name: 'Exposure Data', value: 'exposure_data' }
]
const ibnrTableTypes = [
  { name: 'Claims Data', value: 'claims_data' },
  { name: 'Earned Premium', value: 'earned_premium' }
]

// computed
const availableYears: any = computed(() => {
  const years: any = []
  const currentYear: any = new Date().getFullYear()

  for (let i = 0; i < 20; i++) {
    years.push(currentYear - i)
  }
  return years
})

const availableMonths: any = computed(() => {
  const months: any = []
  for (let i = 1; i < 13; i++) {
    months.push({ month: i })
  }
  return months
})

// methods
const closeDialog = () => {
  dialog.value = false
  uploadSuccess.value = false
  uploadInProgress.value = false
  selectedYear.value = null
  file.value = null
  yieldCurveCode.value = null
  selectedMonth.value = null
  uploadDisabled.value = false
  selectedTableType.value = null
  dataVersion.value = null
}

// processing of the captured data will be delegated to the parent caller
const uploadFile = () => {
  if (file.value !== null) {
    const payload: DataPayload = {
      file: file.value,
      selectedYear: selectedYear.value,
      selectedMonth: selectedMonth.value,
      yieldCurveCode: yieldCurveCode.value,
      selectedType: props.tableType,
      fileName: file.value.name,
      version: dataVersion.value,
      riskRateCode: dataRiskRateCode.value,
      dataTableType: selectedTableType.value,
      schemeName: dataSchemeName.value
    }
    uploadInProgress.value = true
    uploadDisabled.value = true

    emit('uploadFile', payload)
  }
}

watch(
  () => props.uploadComplete,
  (value) => {
    if (value) {
      closeDialog()
      uploadInProgress.value = false
      uploadDisabled.value = false
    }
  }
)
</script>
