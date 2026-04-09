<template>
  <v-container fluid>
    <v-row>
      <v-col>
        <base-card>
          <template #header>
            <span class="headline">PHI Valuation Run Results</span>
          </template>
          <template #default>
            <v-row
              ><v-col>
                <v-btn
                  rounded
                  size="small"
                  class="mb-2"
                  variant="outlined"
                  @click="getSelectStatus"
                  >{{ selectBtnText }}</v-btn
                >
                <v-btn
                  v-if="selectedItems.length > 0 && showSelect"
                  rounded
                  size="small"
                  class="ml-2 mb-2"
                  variant="outlined"
                  @click="selectedItems = []"
                  >Clear Selection</v-btn
                >
                <v-btn
                  v-if="selectedItems.length > 0 && showSelect"
                  rounded
                  size="small"
                  color="red"
                  class="ml-2 mb-2"
                  variant="outlined"
                  @click="deleteRuns(selectedItems)"
                  >Delete Selected</v-btn
                >
              </v-col></v-row
            >
            <loading-indicator :loadingData="loadingData" />
            <div
              v-if="!loadingData && runJobs.length === 0"
              class="text-center pa-8"
            >
              <v-icon size="48" color="grey-lighten-1"
                >mdi-folder-open-outline</v-icon
              >
              <p class="text-grey mt-2">No run results found.</p>
              <router-link to="/group-pricing/phi/run-settings"
                >Go to Run Settings</router-link
              >
              to configure and start a valuation run.
            </div>
            <v-expansion-panels v-if="!loadingData && runJobs.length > 0">
              <v-expansion-panel v-for="job in paginatedJobs" :key="job.id">
                <v-expansion-panel-title class="custom-panel-title px-3">
                  <template #default>
                    <v-row no-gutters>
                      <v-col class="d-flex align-center justify-start" cols="3">
                        <v-checkbox
                          v-if="showSelect"
                          density="compact"
                          class="mr-2 no-padding-checkbox"
                          hide-details
                          :modelValue="isSelected(job.id)"
                          @click.stop="toggleSelect(job.id)"
                        />
                        {{ job.run_name }}
                        <v-chip
                          :color="statusColor(job.status)"
                          size="x-small"
                          class="ml-2"
                          >{{ job.status }}</v-chip
                        >
                      </v-col>
                      <v-col class="d-flex align-center" cols="9">
                        <v-list-item-subtitle
                          v-if="job.status == 'In Progress'"
                        >
                          <span>{{ job.status }}</span>
                          <v-progress-linear
                            :model-value="
                              getProgressPercentage(
                                job.points_done,
                                job.model_point_count
                              )
                            "
                            color="orange"
                            height="8"
                            rounded
                            class="mt-1 mb-1"
                          ></v-progress-linear>
                          <span class="text-caption">
                            {{
                              getProgressPercentage(
                                job.points_done,
                                job.model_point_count
                              )
                            }}% ({{ job.points_done }} /
                            {{ job.model_point_count }} model points)
                          </span>
                        </v-list-item-subtitle>
                        <v-list-item-subtitle v-else>
                          Start: {{ formatDateString(job.creation_date) }} |
                          Duration: {{ toMinutes(job.run_time) }} | User:
                          {{ job.user_name }}
                        </v-list-item-subtitle>
                      </v-col>
                    </v-row>
                  </template>
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <v-container>
                    <v-row>
                      <v-col cols="8"></v-col>
                      <v-col cols="4">
                        <v-btn
                          v-if="job.status !== 'Failed'"
                          variant="outlined"
                          rounded
                          size="small"
                          color="primary"
                          @click.stop="viewControl(job.id)"
                          >View Control</v-btn
                        >
                        <v-btn
                          v-if="job.status !== 'Failed'"
                          variant="outlined"
                          rounded
                          class="ml-2 mr-2"
                          size="small"
                          color="primary"
                          :to="'/group-pricing/phi/run-detail/' + job.id"
                          >View Results</v-btn
                        >
                      </v-col>
                    </v-row>
                  </v-container>

                  <v-divider></v-divider>
                  <v-row>
                    <v-col class="d-flex justify-space-between">
                      <v-btn
                        variant="plain"
                        :loading="actionLoading"
                        rounded
                        size="small"
                        color="red"
                        @click="deleteRun(job.id)"
                        >Delete {{ job.run_name }}</v-btn
                      >
                      <v-btn
                        variant="plain"
                        :loading="actionLoading"
                        rounded
                        size="small"
                        color="primary"
                        @click="exportAllProducts(job.id)"
                        >Export All Results</v-btn
                      >
                    </v-col>
                  </v-row>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
          </template>
          <template #actions>
            <v-row v-if="!loadingData">
              <v-col>
                <v-pagination
                  v-if="totalPages > 1"
                  v-model="currentPage"
                  :length="totalPages"
                ></v-pagination>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <file-info
          :tableTitle="tableTitle"
          :rowData="rowData"
          :columnDefs="columnDefs"
          :onUpdate:isInfoDialogOpen="closeInfoBox"
          :isDialogOpen="infoDialog"
          :show-export="true"
        />
      </v-col>
    </v-row>
  </v-container>
  <confirm-dialog ref="confirmDelete" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/renderer/components/ConfirmDialog.vue'
import BaseCard from '../../../components/BaseCard.vue'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import PhiValuationService from '../../../api/PhiValuationService'
import { DateTime } from 'luxon'
import FileInfo from '@/renderer/components/FileInfo.vue'
import LoadingIndicator from '@/renderer/components/LoadingIndicator.vue'
import formatValues from '@/renderer/utils/format_values'
import { useStatusBarStore } from '@/renderer/store/statusBar'

const statusBarStore = useStatusBarStore()

let pollTimer: any = null

// selection variables
const selectedItems: any = ref([])
const showSelect = ref(false)

const selectBtnText = computed(() =>
  showSelect.value ? 'Hide Selection' : 'Show Selection'
)

const infoDialog = ref(false)
const loadingData = ref(false)
const actionLoading = ref(false)

const tableTitle = ref('')
const rowData: any = ref([])
const columnDefs: any = ref([])
const confirmDelete = ref()
const runJobs = ref([])
const pageSize = 10
const currentPage = ref(1)

const totalPages = computed(() => Math.ceil(runJobs.value.length / pageSize))

const statusColor = (s: string) =>
  ({
    Completed: '#4CAF50',
    Failed: '#F44336',
    'In Progress': '#FF9800',
    Queued: '#9E9E9E'
  })[s] ?? '#9E9E9E'

const paginatedJobs: any = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  const end = start + pageSize
  return runJobs.value.slice(start, end)
})

const formatDateString = (dateString: any) => {
  return DateTime.fromISO(dateString).toLocaleString(DateTime.DATETIME_MED)
}

const toMinutes = (number: any) => {
  number = number * 60
  const minutes = Math.floor(number / 60)
  let seconds = ((number % 60) / 100) * 60
  seconds = Math.round(seconds)
  return minutes + ' m, ' + seconds + ' s'
}

const exportAllProducts = async (jobId: any) => {
  actionLoading.value = true
  const response = await PhiValuationService.getJobExcelResults(jobId)
  actionLoading.value = false

  const fileURL = window.URL.createObjectURL(new Blob([response.data]))
  const fileName = 'valuation_run_' + jobId + '.xlsx'
  const link = document.createElement('a')
  link.href = fileURL
  link.setAttribute('download', fileName)
  document.body.appendChild(link)
  link.click()
}

const isSelected = (id) => {
  return selectedItems.value.includes(id)
}
const closeInfoBox = (value: any) => {
  infoDialog.value = value
}

const toggleSelect = (id: any) => {
  if (selectedItems.value.includes(id)) {
    selectedItems.value = selectedItems.value.filter((item) => item !== id)
  } else {
    selectedItems.value.push(id)
  }
}

const getProgressPercentage = (
  pointsDone: number,
  totalPoints: number
): string => {
  if (totalPoints === 0) return '0.00'
  return (Math.round((pointsDone / totalPoints) * 100 * 100) / 100).toFixed(2)
}

const getSelectStatus = () => {
  showSelect.value = !showSelect.value
}

const loadRuns = async () => {
  const res = await PhiValuationService.getAllRunResults()
  runJobs.value = res.data ?? []
  console.log(runJobs.value)
  const inProgress = runJobs.value.filter(
    (j: any) => j.status === 'In Progress' || j.status === 'Queued'
  ).length
  const completed = runJobs.value.filter(
    (j: any) => j.status === 'Completed'
  ).length
  statusBarStore.set([
    { icon: 'mdi-format-list-bulleted', text: `Runs: ${runJobs.value.length}` },
    { icon: 'mdi-check-circle-outline', text: `Completed: ${completed}` },
    ...(inProgress > 0
      ? [
          {
            icon: 'mdi-progress-clock',
            text: `In progress: ${inProgress}`,
            severity: 'warn' as any
          }
        ]
      : [])
  ])
}

const deleteRun = async (runId: any) => {
  const result = await confirmDelete.value.open(
    'Deleting PHI Job Run',
    'Are you sure you want to delete the selected valuation job run?'
  )
  if (!result) {
    return
  }

  actionLoading.value = true
  await PhiValuationService.deleteProjectionJob(runId)
  await loadRuns()
  actionLoading.value = false
}

const deleteRuns = async (runIds: any) => {
  try {
    const result = await confirmDelete.value.open(
      'Delete Valuation Jobs',
      'Are you sure you want to delete the selected valuation jobs?'
    )
    if (!result) {
      return
    }
    actionLoading.value = true
    await PhiValuationService.deleteValuationJobs(runIds)
    await loadRuns()
    actionLoading.value = false
  } catch (error) {
    actionLoading.value = false
  }
}

const viewControl = async (jobId: any) => {
  const resp = await PhiValuationService.getControlRunResult(jobId)
  rowData.value = resp.data
  createColumnDefs(rowData.value)
  tableTitle.value = 'Run Name Control Results'
  infoDialog.value = true
}

const createColumnDefs = (data) => {
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

onMounted(async () => {
  loadingData.value = true
  await loadRuns()
  loadingData.value = false

  if (
    runJobs.value.length > 0 &&
    runJobs.value.some(
      (job: any) => job.status === 'In Progress' || job.status === 'Queued'
    )
  ) {
    pollTimer = setInterval(() => {
      if (
        runJobs.value.some(
          (job: any) => job.status === 'In Progress' || job.status === 'Queued'
        )
      ) {
        PhiValuationService.getAllRunResults().then((response) => {
          runJobs.value = response.data
        })
      } else {
        clearInterval(pollTimer)
      }
    }, 3000)
  }
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
  statusBarStore.clear()
})
</script>

<style scoped>
/* Custom styling to minimize the spacing */
.custom-panel-title {
  padding: 0; /* Remove default padding */
  height: auto; /* Let height adjust automatically */
  display: flex;
  align-items: center; /* Align items vertically in the center */
}

/* Remove padding and margin from checkbox */
.no-padding-checkbox {
  padding: 0;
  margin: 0;
}

.v-checkbox {
  margin-top: 0;
  margin-bottom: 0;
}
</style>
