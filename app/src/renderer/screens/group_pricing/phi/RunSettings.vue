<template>
  <v-container fluid>
    <v-row>
      <v-col>
        <base-card>
          <template #header>
            <span class="headline">PHI Run Settings</span>
          </template>
          <template #default>
            <v-container>
              <form @submit.prevent="addToRunJobs">
                <v-row align="center">
                  <v-col cols="6">
                    <v-select
                      v-model="selectedSavedConfig"
                      :items="savedConfigs"
                      item-title="name"
                      return-object
                      clearable
                      label="Load a Saved Configuration"
                      density="compact"
                      variant="outlined"
                    />
                  </v-col>
                  <v-col cols="auto">
                    <v-btn
                      size="small"
                      variant="outlined"
                      :disabled="!selectedSavedConfig"
                      @click="loadConfig"
                      >Load</v-btn
                    >
                  </v-col>
                  <v-col cols="auto">
                    <v-btn
                      v-if="selectedSavedConfig"
                      icon
                      size="small"
                      color="red"
                      variant="plain"
                      @click="deleteConfig(selectedSavedConfig)"
                      ><v-icon>mdi-delete</v-icon></v-btn
                    >
                  </v-col>
                </v-row>
                <v-row>
                  <v-col cols="6">
                    <v-text-field
                      v-model="settingRunName"
                      :error-messages="errors.settingRunName"
                      v-bind="settingRunNameAttrs"
                      variant="outlined"
                      density="compact"
                      label="Enter a name for this run"
                    ></v-text-field>
                  </v-col>
                  <v-col cols="6">
                    <v-date-input
                      v-model="runDate"
                      hide-actions
                      locale="en-in"
                      view-mode="month"
                      prepend-icon=""
                      prepend-inner-icon="$calendar"
                      variant="outlined"
                      density="compact"
                      label="Run Date"
                    ></v-date-input>
                  </v-col>
                </v-row>
                <v-row>
                  <v-col>
                    <v-textarea
                      v-model="settingDescription"
                      variant="outlined"
                      rows="3"
                      placeholder="Enter a description for this run"
                    ></v-textarea>
                  </v-col>
                </v-row>
                <v-row>
                  <v-col cols="4">
                    <v-select
                      v-model="selectedModelPointYear"
                      density="compact"
                      variant="outlined"
                      label="Select Model Point Year"
                      :items="availableModelPointYears"
                      item-title="Year"
                      item-value="Year"
                      @update:model-value="getModelPointVersions"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model="selectedModelPointVersion"
                      density="compact"
                      variant="outlined"
                      label="Select Model Point Version"
                      :items="availableModelPointVersions"
                      item-title="Year"
                      item-value="Year"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model="selectedParameterYear"
                      density="compact"
                      variant="outlined"
                      label="Select Parameter Year"
                      :items="availableParameterYears"
                      item-title="Year"
                      item-value="Year"
                      @update:model-value="getParameterVersions"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model="selectedParameterVersion"
                      density="compact"
                      variant="outlined"
                      label="Parameter Version"
                      placeholder="Select a Parameter Version"
                      :items="availableParameterVersions"
                      item-title="name"
                      item-value="name"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model="selectedMortalityYear"
                      density="compact"
                      variant="outlined"
                      label="Select Mortality Year"
                      :items="availableMortalityYears"
                      item-title="Year"
                      item-value="Year"
                      @update:model-value="getMortalityVersions"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model="selectedMortalityVersion"
                      density="compact"
                      variant="outlined"
                      label="Mortality Version"
                      placeholder="Select a Mortalty Version"
                      :items="availableMortalityVersions"
                      item-title="name"
                      item-value="name"
                    ></v-select>
                  </v-col>

                  <v-col cols="4">
                    <v-select
                      v-model="selectedRecoveryYear"
                      density="compact"
                      variant="outlined"
                      label="Recovery Year"
                      placeholder="Select a Recovery Year"
                      :items="availableRecoveryYears"
                      item-title="Year"
                      item-value="Year"
                      @update:model-value="getRecoveryVersions"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model="selectedRecoveryVersion"
                      density="compact"
                      variant="outlined"
                      label="Recovery Version"
                      placeholder="Select a Recovery Version"
                      :items="availableRecoveryVersions"
                      item-title="name"
                      item-value="name"
                    ></v-select>
                  </v-col>

                  <v-col cols="4">
                    <v-select
                      v-model="selectedYieldCurveYear"
                      density="compact"
                      variant="outlined"
                      label="Select Yield Curve Year"
                      :items="availableYieldCurveYears"
                      item-title="Year"
                      item-value="Year"
                      @update:model-value="getYieldCurveVersions"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model="selectedYieldCurveVersion"
                      density="compact"
                      variant="outlined"
                      label="Select Yield Curve Version"
                      :items="availableYieldCurveVersions"
                      item-title="name"
                      item-value="name"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model="selectedShock"
                      density="compact"
                      variant="outlined"
                      label="Select Applicable Shock"
                      :items="shockData"
                      item-title="name"
                      item-value="name"
                      return-object
                    ></v-select>
                  </v-col>
                  <v-col cols="4"
                    ><v-text-field
                      v-model="aggPeriod"
                      density="compact"
                      variant="outlined"
                      type="number"
                      label="Aggregation Period"
                      placeholder="Enter Aggregation period in months"
                    ></v-text-field>
                  </v-col>
                  <v-col cols="4"
                    ><v-text-field
                      v-model="yearEndMonth"
                      density="compact"
                      variant="outlined"
                      type="number"
                      max="12"
                      min="1"
                      label="Year End Month"
                      placeholder="Enter Year End Month (1 - 12)"
                    ></v-text-field>
                  </v-col>
                </v-row>
                <v-row>
                  <v-col cols="3">
                    <v-checkbox
                      v-model="runSingle"
                      :label="`Use a Single Model Point`"
                    ></v-checkbox>
                  </v-col>
                </v-row>
                <v-row>
                  <v-col>
                    <v-btn
                      rounded
                      type="submit"
                      variant="outlined"
                      size="small"
                      class="primary"
                      >Add to Run Jobs</v-btn
                    >
                  </v-col>
                </v-row>
              </form>
              <v-row v-if="runJobs.length > 0">
                <v-divider class="my-5"></v-divider>
              </v-row>

              <v-row v-if="runJobs.length > 0">
                <v-col>
                  <v-table class="trans-tables">
                    <thead>
                      <tr class="table-row">
                        <th class="text-left table-col">Run Name</th>
                        <th class="text-left table-col">Run Date</th>
                        <th class="text-left table-col">Description</th>
                        <th class="text-left table-col">Model Point Year</th>
                        <th class="text-left table-col">Model Point Version</th>
                        <th class="text-left table-col">Parameter Year</th>
                        <th class="text-left table-col">Parameter Version</th>
                        <th class="text-left table-col">Yield Curve Year</th>
                        <th class="text-left table-col">Yield Curve Version</th>
                        <th class="text-left table-col">Recovery Year</th>
                        <th class="text-left table-col">Recovery Version</th>
                        <th class="text-left table-col">Mortality Year </th>
                        <th class="text-left table-col">Mortality Version </th>
                        <th class="text-left table-col">Shock</th>
                        <th class="text-left">Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="item in runJobs" :key="item.run_name">
                        <td>{{ item.run_name }}</td>
                        <td>{{ item.run_date }}</td>
                        <td>{{ item.run_description }}</td>
                        <td>{{ item.modelpoint_year }}</td>
                        <td>{{ item.modelpoint_version }}</td>
                        <td>{{ item.parameter_year }}</td>
                        <td>{{ item.parameter_version }}</td>
                        <td>{{ item.yield_curve_year }}</td>
                        <td>{{ item.yield_curve_version }}</td>
                        <td>{{ item.recovery_year }}</td>
                        <td>{{ item.recovery_version }}</td>
                        <td>{{ item.mortality_year }}</td>
                        <td>{{ item.mortality_version }}</td>
                        <td>
                          {{
                            item.shock_settings !== null
                              ? item.shock_settings.name
                              : 'N/A'
                          }}
                        </td>
                        <td>
                          <v-btn
                            color="red"
                            size="small"
                            variant="plain"
                            icon
                            @click="removeFromRunJobs(item.run_name)"
                          >
                            <v-icon>mdi-delete</v-icon>
                          </v-btn>
                        </td>
                      </tr>
                    </tbody>
                  </v-table>
                </v-col>
              </v-row>
            </v-container>
          </template>
          <template #actions>
            <v-btn
              color="secondary"
              variant="outlined"
              @click="saveConfigDialog = true"
              >Save Configuration</v-btn
            >
            <v-btn
              v-if="runJobs.length > 0"
              color="primary"
              @click="runProjections"
              >Run Valuations</v-btn
            >
          </template>
        </base-card>
      </v-col>
    </v-row>
    <v-snackbar v-model="snackbar" :timeout="timeout">
      {{ snackbarText }}
      <template #actions>
        <v-btn color="white" variant="text" @click="snackbar = false">
          Close
        </v-btn>
      </template>
    </v-snackbar>
    <v-dialog v-model="saveConfigDialog" persistent max-width="600">
      <base-card class="rounded-lg">
        <template #header>
          <span class="headline">Save Run Configuration</span>
        </template>
        <template #default>
          <v-container>
            <v-text-field
              v-model="saveConfigName"
              variant="outlined"
              density="compact"
              label="Configuration Name"
            ></v-text-field>
            <v-list density="compact">
              <v-list-item
                >Model Point: {{ selectedModelPointYear }} /
                {{ selectedModelPointVersion }}</v-list-item
              >
              <v-list-item
                >Parameters: {{ selectedParameterYear }} /
                {{ selectedParameterVersion }}</v-list-item
              >
              <v-list-item
                >Mortality: {{ selectedMortalityYear }} /
                {{ selectedMortalityVersion }}</v-list-item
              >
              <v-list-item
                >Recovery: {{ selectedRecoveryYear }} /
                {{ selectedRecoveryVersion }}</v-list-item
              >
              <v-list-item
                >Yield Curve: {{ selectedYieldCurveYear }} /
                {{ selectedYieldCurveVersion }}</v-list-item
              >
              <v-list-item
                >Shock: {{ selectedShock?.name ?? 'None' }}</v-list-item
              >
              <v-list-item>Aggregation Period: {{ aggPeriod }}</v-list-item>
              <v-list-item>Year End Month: {{ yearEndMonth }}</v-list-item>
            </v-list>
          </v-container>
        </template>
        <template #actions>
          <v-btn
            @click="
              () => {
                saveConfigName = ''
                saveConfigDialog = false
              }
            "
            >Cancel</v-btn
          >
          <v-btn
            color="primary"
            :loading="saveConfigLoading"
            @click="saveConfig"
            >Save</v-btn
          >
        </template>
      </base-card>
    </v-dialog>
    <confirm-dialog ref="confirmDeleteAction" />
  </v-container>
</template>

<script setup lang="ts">
import BaseCard from '@/renderer/components/BaseCard.vue'
import { ref, onMounted } from 'vue'
import PhiValuationService from '@/renderer/api/PhiValuationService'
import { VDateInput } from 'vuetify/labs/VDateInput'
import { useForm } from 'vee-validate'
import * as yup from 'yup'
import formatDateString from '@/renderer/utils/helpers'
import { useRouter } from 'vue-router'
import ConfirmDialog from '@/renderer/components/ConfirmDialog.vue'

const $router = useRouter()

const validationSchema = yup.object({
  settingRunName: yup.string().required('Run name is required'),
  selectedProducts: yup
    .array()
    .min(1, 'You need to select at least one product'),
  selectedModelPointYear: yup.string().required('Model point year is required')
})

const { defineField, errors } = useForm({
  validationSchema
})

const [settingRunName, settingRunNameAttrs] = defineField('settingRunName')

const confirmDeleteAction = ref()
const snackbarText = ref('')
const timeout = ref(0)
const snackbar = ref(false)

const savedConfigs: any = ref([])
const selectedSavedConfig: any = ref(null)
const saveConfigDialog = ref(false)
const saveConfigName = ref('')
const saveConfigLoading = ref(false)
const runSingle = ref(false)

const selectedYieldCurveMonth: any = ref(null)
const selectedShock: any = ref(null)
const aggPeriod: any = ref(null)
const yearEndMonth: any = ref(12)

const selectedModelPointYear = ref(null)
const availableModelPointYears: any = ref([])

const selectedModelPointVersion: any = ref(null)
const availableModelPointVersions: any = ref([])

const selectedParameterYear: any = ref(null)
const availableParameterYears: any = ref([])

const selectedParameterVersion: any = ref(null)
const availableParameterVersions: any = ref([])

const selectedMortalityYear: any = ref(null)
const availableMortalityYears: any = ref([])

const selectedMortalityVersion: any = ref(null)
const availableMortalityVersions: any = ref([])

const selectedRecoveryYear: any = ref(null)
const availableRecoveryYears: any = ref([])

const selectedRecoveryVersion: any = ref(null)
const availableRecoveryVersions: any = ref([])

const selectedYieldCurveYear: any = ref(null)
const availableYieldCurveYears: any = ref([])

const selectedYieldCurveVersion: any = ref(null)
const availableYieldCurveVersions: any = ref([])

const selectedBasis: any = ref(null)
const runDate = ref(null)
const shockData: any = ref([])
const runJobs: any = ref([])
const settingDescription = ref('')

const removeFromRunJobs = (runName) => {
  runJobs.value = runJobs.value.filter((item) => item.run_name !== runName)
}

// const executeJobs = () => {
//   runProjections()
// }

const returnValidationError = (text) => {
  snackbarText.value = text
  timeout.value = 3000
  snackbar.value = true
}

const getModelPointVersions = async () => {
  if (selectedModelPointYear.value !== null) {
    availableModelPointVersions.value = []
    selectedModelPointVersion.value = null
    const resp = await PhiValuationService.getModelPointVersionsForYear(
      selectedModelPointYear.value
    )
    if (resp.data !== null) {
      availableModelPointVersions.value = resp.data
    }
  }
}

const getMortalityVersions = async () => {
  if (selectedMortalityYear.value !== null) {
    availableMortalityVersions.value = []
    selectedMortalityVersion.value = null
    const resp = await PhiValuationService.getAvailableMortalityVersions(
      selectedMortalityYear.value
    )
    if (resp.data !== null) {
      availableMortalityVersions.value = resp.data
    }
  }
}

const getRecoveryVersions = async () => {
  if (selectedRecoveryYear.value !== null) {
    availableRecoveryVersions.value = []
    selectedRecoveryVersion.value = null
    const resp = await PhiValuationService.getAvailableRecoveryVersions(
      selectedRecoveryYear.value
    )
    if (resp.data !== null) {
      availableRecoveryVersions.value = resp.data
    }
  }
}

const getYieldCurveVersions = async () => {
  if (selectedYieldCurveYear.value !== null) {
    selectedYieldCurveMonth.value = null
    const resp = await PhiValuationService.getAvailableYieldCurveVersions(
      selectedYieldCurveYear.value
    )
    if (resp.data !== null) {
      availableYieldCurveVersions.value = resp.data
    }
  }
}

const getParameterVersions = async () => {
  availableParameterVersions.value = []
  const res = await PhiValuationService.getAvailableParameterVersions(
    selectedParameterYear.value
  )

  if (res.data !== null && res.data.length > 0) {
    availableParameterVersions.value = res.data
  }
}

const loadConfig = async () => {
  const cfg = selectedSavedConfig.value
  if (!cfg) return
  selectedModelPointYear.value = cfg.model_point_year
  selectedParameterYear.value = cfg.parameter_year
  selectedMortalityYear.value = cfg.mortality_year
  selectedRecoveryYear.value = cfg.recovery_year
  selectedYieldCurveYear.value = cfg.yield_curve_year
  settingRunName.value = cfg.run_name

  await Promise.all([
    getModelPointVersions(),
    getParameterVersions(),
    getMortalityVersions(),
    getRecoveryVersions(),
    getYieldCurveVersions()
  ])

  selectedModelPointVersion.value = cfg.model_point_version
  selectedParameterVersion.value = cfg.parameter_version
  selectedMortalityVersion.value = cfg.mortality_version
  selectedRecoveryVersion.value = cfg.recovery_version
  selectedYieldCurveVersion.value = cfg.yield_curve_version
  aggPeriod.value = cfg.aggregation_period
  yearEndMonth.value = cfg.year_end_month
  runSingle.value = cfg.run_single

  if (cfg.shock_settings_id === 0) {
    selectedShock.value = shockData.value.find(
      (s: any) => s.name === 'N/A'
    ) ?? { name: 'N/A' }
  } else {
    selectedShock.value =
      shockData.value.find((s: any) => s.id === cfg.shock_settings_id) ?? null
  }
}

const saveConfig = async () => {
  if (!saveConfigName.value) return
  saveConfigLoading.value = true
  try {
    const payload = {
      name: saveConfigName.value,
      run_name: settingRunName.value,
      model_point_year: selectedModelPointYear.value,
      model_point_version: selectedModelPointVersion.value,
      parameter_year: selectedParameterYear.value,
      parameter_version: selectedParameterVersion.value,
      mortality_year: selectedMortalityYear.value,
      mortality_version: selectedMortalityVersion.value,
      recovery_year: selectedRecoveryYear.value,
      recovery_version: selectedRecoveryVersion.value,
      yield_curve_year: selectedYieldCurveYear.value,
      yield_curve_version: selectedYieldCurveVersion.value,
      shock_settings_id: selectedShock.value?.id ?? 0,
      shock_settings_name: selectedShock.value?.name ?? 'N/A',
      aggregation_period: aggPeriod.value ? parseInt(aggPeriod.value) : 0,
      year_end_month: yearEndMonth.value ? parseInt(yearEndMonth.value) : 12,
      run_single: runSingle.value
    }

    console.log('Saving config with payload:', payload)
    const resp = await PhiValuationService.savePhiRunConfig(payload)
    savedConfigs.value.unshift(resp.data)
    saveConfigName.value = ''
    saveConfigDialog.value = false
  } catch (err: any) {
    snackbarText.value =
      err?.response?.data?.error ?? 'Failed to save configuration'
    timeout.value = 4000
    snackbar.value = true
  } finally {
    saveConfigLoading.value = false
  }
}

const deleteConfig = async (cfg: any) => {
  confirmDeleteAction.value.open(
    'Delete Configuration',
    `Delete "${cfg.name}"?`,
    { color: 'red' },
    async () => {
      await PhiValuationService.deletePhiRunConfig(cfg.id)
      savedConfigs.value = savedConfigs.value.filter(
        (c: any) => c.id !== cfg.id
      )
      if (selectedSavedConfig.value?.id === cfg.id) {
        selectedSavedConfig.value = null
      }
    }
  )
}

const addToRunJobs = async () => {
  if (!settingRunName.value)
    return returnValidationError('Run name is required')
  if (!runDate.value) return returnValidationError('Run date is required')
  if (!selectedModelPointYear.value)
    return returnValidationError('Model point year is required')
  if (!selectedShock.value)
    return returnValidationError('Please select a shock (or N/A)')

  const job: any = {}
  const ids: any = []
  const names: any = []
  job.prod_ids = ids
  job.prod_names = names
  job.jobs_template_id = 0
  job.run_name = settingRunName.value
  job.run_description = settingDescription.value
  job.run_single = runSingle.value
  job.run_date = formatDateString(runDate.value, true, true, false)
  job.modelpoint_year = selectedModelPointYear.value
  job.modelpoint_version = selectedModelPointVersion.value
  job.yield_curve_year = selectedYieldCurveYear.value
  job.yield_curve_version = selectedYieldCurveVersion.value

  job.parameter_year = selectedParameterYear.value
  job.parameter_version = selectedParameterVersion.value
  job.mortality_year = selectedMortalityYear.value
  job.mortality_version = selectedMortalityVersion.value
  job.recovery_year = selectedRecoveryYear.value
  job.recovery_version = selectedRecoveryVersion.value
  job.shock_settings_id = selectedShock.value?.id ?? 0
  job.shock_settings = selectedShock.value ?? null
  if (!isNaN(parseInt(aggPeriod.value)) && parseInt(aggPeriod.value) > 0) {
    if (parseInt(aggPeriod.value) > 1440) {
      job.aggregation_period = 1440
    } else {
      job.aggregation_period = parseInt(aggPeriod.value)
    }
  }
  job.year_end_month = yearEndMonth.value

  job.run_basis = selectedBasis.value

  runJobs.value.push(job)
  resetFields()
}

const resetFields = () => {
  settingRunName.value = ''
  runDate.value = null
  selectedModelPointYear.value = null
  selectedModelPointVersion.value = null
  selectedYieldCurveYear.value = null
  selectedYieldCurveMonth.value = null
  selectedParameterYear.value = null
  selectedMortalityYear.value = null
  selectedShock.value = null
  aggPeriod.value = null
  yearEndMonth.value = null
  selectedBasis.value = null
  settingDescription.value = ''
  runSingle.value = false
}

const runProjections = () => {
  const payload: any = {}

  payload.jobs = runJobs.value
  payload.jobs_template_id = 0

  PhiValuationService.runProjections(payload)
    .then((resp) => {
      console.log(resp)
      timeout.value = 3000
      snackbar.value = true
      snackbarText.value = 'the run has been started.'
      setTimeout(() => {
        $router.push({ name: 'group-pricing-phi-run-results' })
      }, 3500)
    })
    .catch((err) => {
      console.log(err)
      timeout.value = 5000
      snackbar.value = true
      snackbarText.value =
        err?.response?.data?.message ?? 'Failed to run projections'
    })
}

onMounted(async () => {
  const [rMps, rParams, rMortality, rRecovery, rYieldCurveYears] =
    await Promise.all([
      PhiValuationService.getAvailableModelPointYears(),
      PhiValuationService.getAvailableParameterYears(),
      PhiValuationService.getAvailableMortalityYears(),
      PhiValuationService.getAvailableRecoveryYears(),
      PhiValuationService.getAvailableYieldCurveYears()
    ])

  if (rMps.data !== null && rMps.data.length > 0) {
    availableModelPointYears.value = rMps.data
  }
  if (rParams.data !== null && rParams.data.length > 0) {
    availableParameterYears.value = rParams.data
  }
  if (rMortality.data !== null && rMortality.data.length > 0) {
    availableMortalityYears.value = rMortality.data
  }
  if (rRecovery.data !== null && rRecovery.data.length > 0) {
    availableRecoveryYears.value = rRecovery.data
  }
  if (rYieldCurveYears.data !== null && rYieldCurveYears.data.length > 0) {
    availableYieldCurveYears.value = rYieldCurveYears.data
  }

  PhiValuationService.getPhiRunConfigs()
    .then((r) => {
      if (r.data !== null && r.data.length > 0) {
        savedConfigs.value = r.data
      }
    })
    .catch(() => {})

  PhiValuationService.getShockSettings().then((response) => {
    shockData.value = []
    const naObj: any = { name: 'N/A' }
    if (response.data !== null && response.data.length > 0) {
      shockData.value = response.data
    }
    shockData.value.unshift(naObj)
  })
})
</script>

<style scoped>
.responsive-table {
  display: inline-block;
}

.v-table {
  border: 1px solid grey;
}
</style>
