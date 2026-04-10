<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <span class="headline">Valuation Result Detail</span>
          </template>
          <template #default>
            <v-row class="mb-3">
              <v-col>
                <v-btn variant="plain" :to="'/group-pricing/phi/run-results'">
                  {{ backButtonTitle }}
                </v-btn>
              </v-col>
            </v-row>
            <loading-indicator
              class="mb-5"
              :loading-data="loadingData"
            ></loading-indicator>
            <base-card v-if="!loadingData">
              <template #header>
                <span class="headline">Aggregated Results</span>
              </template>
              <template #default>
                <v-row v-if="showSpCodeSelect">
                  <v-col cols="3">
                    <v-select
                      v-model="selectedSpCode"
                      variant="outlined"
                      density="compact"
                      placeholder="Select an SP Code"
                      label="SP Code"
                      :items="spCodes"
                      @update:modelValue="displayAggDataForSpCode"
                    ></v-select>
                  </v-col>
                </v-row>
                <loading-indicator
                  :loadingData="loadingSpCodeData"
                ></loading-indicator>
                <data-grid
                  v-if="rowData.length > 0"
                  :run-name="runName"
                  :run-id="runId"
                  :show-full-export="true"
                  :showExport="true"
                  :rowData="rowData"
                  :columnDefs="cDefs"
                ></data-grid>
              </template>
            </base-card>
            <base-card v-if="sapRowData.length > 0 && !loadingData">
              <template #header>
                <span class="headline">Scoped Aggregated Results</span>
              </template>
              <template #default>
                <v-row v-if="showSapSpCodeSelect">
                  <v-col cols="3">
                    <v-select
                      v-model="selectedIfrs17Group"
                      variant="outlined"
                      density="compact"
                      placeholder="Select a group"
                      label="IFRS17 Group"
                      :items="ifrs17groups"
                      @update:modelValue="displaySapAggDataForIFRS17Group"
                    ></v-select>
                  </v-col>
                </v-row>

                <data-grid
                  :showExport="showExport"
                  :rowData="sapRowData"
                  :columnDefs="spCDefs"
                ></data-grid>
              </template>
            </base-card>
            <base-card
              v-if="runSettings !== null && !loadingData"
              :showActions="false"
            >
              <template #header>
                <span class="headline">Run Settings</span>
              </template>
              <template #default>
                <v-table>
                  <thead>
                    <tr class="table-row">
                      <th class="text-left table-col">Run Name</th>
                      <th class="text-left table-col">Run Date</th>
                      <th class="text-left table-col">Model Points</th>
                      <th class="text-left table-col">Model Point Version</th>
                      <th class="text-left table-col">Yield Curve</th>
                      <th class="text-left table-col">Yield Curve Month</th>
                      <th class="text-left table-col">Parameters</th>
                      <th class="text-left table-col">Transitions</th>
                      <th class="text-left table-col">Morbidity</th>
                      <th class="text-left table-col">Mortality</th>
                      <th class="text-left table-col">Lapse</th>
                      <th class="text-left table-col">Lapse Margin</th>
                      <th class="text-left table-col">Retrenchment</th>
                      <th class="text-left table-col">IFRS17</th>
                      <th class="text-left table-col">Shock Setting</th>
                      <th class="text-left table-col">Yield Curve Basis</th>
                      <th class="text-left table-col">Run Basis</th>
                      <th class="text-left table-col">Single Run</th>
                      <th class="text-left table-col">Aggregation Period</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr>
                      <td>{{ runSettings.run_name }}</td>
                      <td>{{ runSettings.run_date }}</td>
                      <td>{{ runSettings.modelpoint_year }}</td>
                      <td>{{ runSettings.mp_version }}</td>
                      <td>{{ runSettings.yieldcurve_year }}</td>
                      <td>{{ runSettings.yieldcurve_month }}</td>
                      <td>{{ runSettings.parameter_year }}</td>
                      <td>{{ runSettings.transition_year }}</td>
                      <td>{{ runSettings.morbidity_year }}</td>
                      <td>{{ runSettings.mortality_year }}</td>
                      <td>{{ runSettings.lapse_year }}</td>
                      <td>{{ runSettings.lapse_margin_year }}</td>
                      <td>{{ runSettings.retrenchment_year }}</td>
                      <td>{{ runSettings.ifrs17_indicator }}</td>
                      <td>{{ runSettings.shock_setting_name }}</td>
                      <td>{{ runSettings.yield_curve_basis }}</td>
                      <td>{{ runSettings.run_basis }}</td>
                      <td>{{ runSettings.run_single }}</td>
                      <td>{{ runSettings.aggregation_period }}</td>
                    </tr>
                  </tbody>
                </v-table>
              </template>
            </base-card>
          </template>
        </base-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import BaseCard from '@/renderer/components/BaseCard.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import { onMounted, ref, onBeforeMount } from 'vue'
import { useRoute } from 'vue-router'
import PhiValuationService from '@/renderer/api/PhiValuationService'
import LoadingIndicator from '@/renderer/components/LoadingIndicator.vue'

const route = useRoute()

const runName: any = ref(null)
const showExport = ref(true)
const backButtonTitle = '<  Back to Run List'
const runId: any = ref(null)
const loadingData = ref(false)
const gridOptions: any = ref(null)
const rowData: any = ref([])
const sapRowData: any = ref([])
const runSettings: any = ref(null)
const spCodes: any = ref([])
// const runErrors: any = ref(null)
// const prodId: any = ref(null)
const prodName: any = ref(null)
const cDefs = ref([])
const spCDefs = ref([])
const selectedVariable: any = ref('reserves')
const selectedSpCode = ref(null)
const aggData: any = ref([])
const sapAggData: any = ref([])
const showSpCodeSelect = ref(false)
const showSapSpCodeSelect = ref(false)
const selectedIfrs17Group = ref(null)
const ifrs17groups: any = ref([])
const loadingSpCodeData = ref(false)

const chartSeries: any = ref({
  data: [],
  name: null,
  color: null
})

const chartOptions: any = ref({
  credits: {
    enabled: false
  },
  chart: {
    type: 'spline'
  },
  title: {
    text: ''
  },
  xAxis: {
    categories: [],
    title: { text: 'Projection Month' }
  },
  yAxis: {
    title: {
      text: 'Reserves'
    }
  },
  series: []
})

const displayAggDataForSpCode = () => {
  rowData.value = []
  cDefs.value = []
  loadingSpCodeData.value = true
  PhiValuationService.getValuationJobWithSpCode(
    runId.value,
    selectedSpCode.value
  ).then((resp) => {
    rowData.value = resp.data.projections
    cDefs.value = createColumnDefs(rowData.value)
    loadingSpCodeData.value = false
  })
  // rowData.value = aggData.value.filter((item: any) => item.sp_code === spCodes.value[0])
  // cDefs.value = createColumnDefs(rowData.value)
}

const displaySapAggDataForIFRS17Group = () => {
  sapRowData.value = []
  sapRowData.value = sapAggData.value.filter(
    (item: any) => item.IFRS17Group === selectedIfrs17Group.value
  )
  spCDefs.value = createSapColumnDefs(sapRowData.value)
}

// const camelToSnakeCase = (str: String) => str.replace(/[A-Z]/g, (letter) => `_${letter.toLowerCase()}`);

const createColumnDefs = (rowData: any) => {
  if (rowData === null || rowData.length === 0) {
    return []
  }
  if (rowData !== null && rowData.length > 0) {
    const columnDefs: any = []
    const keys = Object.keys(rowData[0])
    keys.forEach((key) => {
      columnDefs.push({
        headerName: key,
        field: key,
        sortable: true,
        filter: true,
        resizable: true,
        width: 150
      })
    })
    gridOptions.value.columnDefs = columnDefs
    return columnDefs
  }
}

const createSapColumnDefs = (rowData: any) => {
  const columnDefs: any = []
  const keys = Object.keys(rowData[0])
  keys.forEach((key) => {
    columnDefs.push({
      headerName: key,
      field: key,
      sortable: true,
      filter: true,
      resizable: true,
      width: 150
    })
  })
  gridOptions.value.sapColumnDefs = columnDefs
  return columnDefs
}

const getAggregatedVariableV2 = (variable, spcode) => {
  chartSeries.value.data = []
  if (aggData.value !== null) {
    aggData.value.forEach((elem: any) => {
      if (
        chartOptions.value.xAxis.categories.indexOf(elem.ProjectionMonth) ===
          -1 &&
        elem.sp_code === spcode
      ) {
        chartOptions.value.xAxis.categories.push(elem.ProjectionMonth)
      }
      if (elem.sp_code === spcode) {
        chartSeries.value.data.push(elem[variable])
      }
    })
    chartSeries.value.name = prodName.value
    chartSeries.value.color = generateHexColorExcludingWhite()
    chartOptions.value.series = []
    chartOptions.value.series.push(chartSeries.value)
    chartOptions.value.yAxis.title.text = selectedVariable.value
  }
}

const generateHexColorExcludingWhite = () => {
  const hexChars = '0123456789ABCDEF'
  let color
  do {
    color = '#'
    for (let i = 0; i < 6; i++) {
      const randomIndex = Math.floor(Math.random() * hexChars.length)
      color += hexChars[randomIndex]
    }
  } while (color === '#FFFFFF')
  return color
}

onBeforeMount(() => {
  getAggregatedVariableV2('reserves', null)
})

onMounted(async () => {
  runId.value = route.params.jobId
  console.log('route params', route.params)

  console.log('runId', runId.value)
  loadingData.value = true
  gridOptions.value = {}

  const resp = await PhiValuationService.getRunResult(runId.value)

  console.log('Response for Phi Valuation Run Result:', resp.data)
  aggData.value = resp.data.aggregated_projections
  rowData.value = aggData.value

  console.log('Row Data:', rowData.value)

  if (rowData.value.length > 0) {
    cDefs.value = createColumnDefs(rowData.value)
  }

  loadingData.value = false
})
</script>

<style scoped>
.table-col {
  min-width: 120px;
  font-size: 12px;
  white-space: nowrap;
}
</style>
