<template>
  <v-container fluid>
    <v-row>
      <v-col>
        <base-card>
          <template #header> Shock Settings </template>
          <template #default>
            <v-container>
              <form @submit.prevent="addToScenarios">
                <h4 class="mb-2">Create PHI Shock Setting</h4>
                <v-row>
                  <v-col cols="3">
                    <v-text-field
                      v-model="settingName.value.value"
                      variant="outlined"
                      density="compact"
                      placeholder="Enter a name for this scenario"
                      :error-messages="settingName.errorMessage.value"
                    ></v-text-field>
                  </v-col>
                </v-row>
                <v-row>
                  <v-col>
                    <v-textarea
                      v-model="settingDescription.value.value"
                      variant="outlined"
                      rows="3"
                      label="Provide a description for this scenario"
                      :error-messages="settingDescription.errorMessage.value"
                    ></v-textarea>
                  </v-col>
                </v-row>
                <v-row>
                  <v-col>
                    <p
                      >Select the desired combination for this setting from the
                      options below</p
                    >
                  </v-col>
                </v-row>
                <v-row>
                  <v-col cols="3">
                    <v-checkbox
                      v-model="mortality"
                      :label="`Mortality`"
                    ></v-checkbox>
                  </v-col>
                  <v-col cols="3">
                    <v-checkbox
                      v-model="recovery"
                      :label="`Recovery`"
                    ></v-checkbox>
                  </v-col>
                  <v-col cols="3">
                    <v-checkbox
                      v-model="expenses"
                      :label="`Expenses`"
                    ></v-checkbox>
                  </v-col>
                  <v-col cols="3">
                    <v-checkbox
                      v-model="mortalityCatastrophe"
                      :label="`Mortality Catastrophe`"
                    ></v-checkbox>
                  </v-col>
                  <v-col cols="3">
                    <v-checkbox
                      v-model="morbidityCatastrophe"
                      :label="`Morbidity Catastrophe`"
                    ></v-checkbox>
                  </v-col>
                  <v-col cols="3">
                    <v-checkbox
                      v-model="realYieldCurve"
                      :label="`Real Yield Curve`"
                    ></v-checkbox>
                  </v-col>
                  <v-col cols="3">
                    <v-checkbox
                      v-model="nominalYieldCurve"
                      :label="`Nominal Yield Curve`"
                    ></v-checkbox>
                  </v-col>
                  <v-col cols="3">
                    <v-checkbox
                      v-model="inflation"
                      :label="`Inflation`"
                    ></v-checkbox>
                  </v-col>
                </v-row>
                <v-row>
                  <v-col cols="3">
                    <v-select
                      v-model="selectedShockBasis.value.value"
                      variant="outlined"
                      density="compact"
                      placeholder="Select a Shock Basis"
                      label="Shock Basis"
                      :items="shockBases"
                      item-title="Basis"
                      item-value="Basis"
                      :error-messages="selectedShockBasis.errorMessage.value"
                    >
                    </v-select>
                  </v-col>
                </v-row>
                <v-row>
                  <v-col>
                    <v-btn
                      type="submit"
                      rounded
                      width="250"
                      size="small"
                      variant="outlined"
                      color="primary"
                      >Add to Shock Settings</v-btn
                    >
                  </v-col>
                </v-row>
              </form>
              <v-row>
                <v-col>
                  <v-divider></v-divider>
                </v-col>
              </v-row>
              <h4 class="mt-4">Available Scenarios</h4>
              <v-row v-if="shockScenarios.length > 0">
                <v-col>
                  <data-grid
                    ref="table"
                    :pagination="true"
                    :rowSelection="rowSelection"
                    :columnDefs="columnDefs"
                    :showExport="showExport"
                    :rowData="shockScenarios"
                    @update:row-deleted="deleteShockSetting"
                  />
                </v-col>
              </v-row>
            </v-container>
            <confirmation-dialog ref="confirmAction" />
          </template>
        </base-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import PhiValuationService from '../../../api/PhiValuationService'
import { onMounted, ref } from 'vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import ConfirmationDialog from '@/renderer/components/ConfirmDialog.vue'
import { useField, useForm } from 'vee-validate'
import * as yup from 'yup'

const schema = yup.object({
  settingName: yup.string().required('Please enter a name for this setting'),
  settingDescription: yup
    .string()
    .required('Please provide a description for this setting'),
  selectedShockBasis: yup.string().nullable()
})

const { handleSubmit, handleReset } = useForm({
  validationSchema: schema
})

const settingName = useField('settingName')
const settingDescription = useField('settingDescription')
const selectedShockBasis = useField('selectedShockBasis')

const confirmAction = ref()
const columnDefs = ref([])
const shockBases: any = ref([])
const mortality = ref(false)
const disability = ref(false)
const recovery = ref(false)
const criticalIllness = ref(false)
const nominalYieldCurve = ref(false)
const realYieldCurve = ref(false)
const expenses = ref(false)
const inflation = ref(false)
const mortalityCatastrophe = ref(false)
const morbidityCatastrophe = ref(false)
const retrenchment = ref(false)
const rowSelection = ref(null)
const shockScenarios: any = ref([])
const showExport = ref(true)

onMounted(() => {
  PhiValuationService.getShockSettings().then((response) => {
    if (response.data !== null && response.data.length > 0) {
      console.log(response.data)
      shockScenarios.value = []
      shockScenarios.value = response.data
      columnDefs.value = createColumnDefs(
        shockScenarios.value,
        'shockScenarios'
      )
    }
  })
  getAvailableShockBases()
})

const getAvailableShockBases = () => {
  PhiValuationService.getAvailableShockBases().then((response) => {
    if (response.data !== null && response.data.length > 0) {
      response.data.forEach((item: any) => {
        shockBases.value.push(item)
      })
    }
  })
}

const deleteShockSetting = async (scenario: any) => {
  console.log('delete', scenario)
  try {
    const res = await confirmAction.value.open(
      `Deleting ${scenario.name} `,
      'Are you sure you want to delete this shock scenario? This action is not undoable.'
    )

    if (res) {
      console.log('confirming delete', scenario)
      PhiValuationService.deleteShockSetting(scenario.id).then((res) => {
        console.log(res)
        if (res.status === 204) {
          console.log('deleted')
          shockScenarios.value = shockScenarios.value.filter((item) => {
            return item.id !== scenario.id
          })
        }
      })
    }
  } catch (error) {
    console.log(error)
  }
}

const resetForm = () => {
  settingName.value.value = ''
  settingDescription.value.value = ''
  mortality.value = false
  disability.value = false
  recovery.value = false
  criticalIllness.value = false
  nominalYieldCurve.value = false
  realYieldCurve.value = false
  expenses.value = false
  inflation.value = false
  mortalityCatastrophe.value = false
  morbidityCatastrophe.value = false
  retrenchment.value = false
  selectedShockBasis.value.value = null
  handleReset()
}
const addToScenarios = handleSubmit(() => {
  console.log('valid')

  const setting: any = {}
  setting.name = settingName.value.value
  setting.description = settingDescription.value.value
  setting.mortality = mortality.value
  setting.recovery = recovery.value
  setting.critical_illness = criticalIllness.value
  setting.nominal_yield_curve = nominalYieldCurve.value
  setting.real_yield_curve = realYieldCurve.value
  setting.expense = expenses.value
  setting.retrenchment = retrenchment.value
  setting.shock_basis = selectedShockBasis.value.value
  setting.inflation = inflation.value
  setting.morbidity_catastrophe = morbidityCatastrophe.value
  setting.mortality_catastrophe = mortalityCatastrophe.value

  //   this.shockScenarios.push(setting);
  PhiValuationService.saveShockSetting(setting).then((res) => {
    if (res.status === 201) {
      shockScenarios.value.push(res.data)
      resetForm()
    }
  })
})

const createColumnDefs = (data: any, tableName: string) => {
  console.log(tableName)
  const cDefs: any = []
  if (typeof data !== 'undefined') {
    Object.keys(data[0]).forEach((element) => {
      if (element !== 'id' && element !== 'year') {
        const header: any = {}
        header.headerName = element
        header.field = element
        header.minWidth = 200
        header.sortable = true
        header.filter = true
        header.resizable = true

        cDefs.push(header)
        if (element === 'name') {
          header.checkboxSelection = true
          header.pinned = 'left'
          header.minWidth = 250
          header.sortable = true
          header.filter = true
          header.resizable = true
        }
      }
    })
  }
  return cDefs
}
</script>

<style scoped></style>
