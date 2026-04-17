<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center justify-space-between">
              <div>
                <span class="headline">Submission Tracking</span>
              </div>
              <v-btn
                size="small"
                rounded
                color="white"
                variant="outlined"
                prepend-icon="mdi-arrow-left"
                @click="$router.push('/group-pricing/bordereaux-management')"
              >
                Back to Dashboard
              </v-btn>
            </div>
          </template>
          <template #default>
            <!-- Filter Controls -->
            <v-row class="mb-4">
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title class="text-h6 font-weight-bold">
                    <v-icon class="me-2">mdi-filter</v-icon>
                    Filter Submissions
                  </v-card-title>
                  <v-card-text>
                    <v-row>
                      <v-col cols="12" sm="6" md="3">
                        <v-select
                          v-model="filters.status"
                          :items="statusOptions"
                          label="Status"
                          variant="outlined"
                          density="compact"
                          clearable
                          multiple
                          chips
                        />
                      </v-col>
                      <v-col cols="12" sm="6" md="3">
                        <v-select
                          v-model="filters.type"
                          :items="bordereauTypes"
                          label="Bordereaux Type"
                          variant="outlined"
                          density="compact"
                          clearable
                        />
                      </v-col>
                      <v-col cols="12" sm="6" md="3">
                        <v-text-field
                          v-model="filters.search"
                          label="Search ID or Reference"
                          variant="outlined"
                          density="compact"
                          prepend-inner-icon="mdi-magnify"
                          clearable
                        />
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Submissions Table -->
            <v-row>
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title
                    class="d-flex align-center justify-space-between"
                  >
                    <span class="text-h6 font-weight-bold"
                      >Bordereaux Submissions</span
                    >
                    <div class="d-flex align-center gap-2">
                      <v-btn
                        v-if="hasPermission('bordereaux:submit_outbound')"
                        color="primary"
                        size="small"
                        rounded
                        class="mr-2"
                        prepend-icon="mdi-plus"
                        :disabled="selectedItems.length === 0"
                        @click="showBatchSubmissionDialog = true"
                      >
                        Batch Submit
                      </v-btn>
                      <v-btn
                        color="success"
                        size="small"
                        rounded
                        prepend-icon="mdi-refresh"
                        :loading="loading"
                        @click="refreshSubmissions"
                      >
                        Refresh
                      </v-btn>
                    </div>
                  </v-card-title>
                  <v-card-text>
                    <DataGrid
                      :row-data="filteredSubmissions"
                      :column-defs="columnDefs"
                      table-title="Bordereaux Submissions"
                      :show-export="false"
                      density="compact"
                      row-selection="multiple"
                      @row-clicked="onRowClicked"
                      @row-selection-changed="onRowSelectionChanged"
                    />
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Batch Submission Dialog -->
    <v-dialog v-model="showBatchSubmissionDialog" max-width="800">
      <v-card>
        <v-card-title class="text-h6 font-weight-bold bg-primary text-white">
          <v-icon class="me-2">mdi-send-check</v-icon>
          Batch Submission
        </v-card-title>
        <v-card-text class="pt-6">
          <v-alert color="info" variant="tonal" class="mb-4">
            <p class="font-weight-bold"
              >Submit multiple bordereaux to schemes</p
            >
            <p class="mb-0"
              >Select bordereaux and configure submission settings</p
            >
          </v-alert>

          <v-row>
            <v-col cols="12" md="6">
              <v-select
                v-model="batchSubmission.delivery_method"
                :items="deliveryMethods"
                label="Delivery Method *"
                variant="outlined"
                density="compact"
                required
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="batchSubmission.message"
                label="Submission Message"
                variant="outlined"
                density="compact"
                rows="3"
                hint="Optional message to include with submission"
                persistent-hint
              />
            </v-col>
          </v-row>

          <v-divider class="my-4"></v-divider>

          <p class="text-subtitle-2 font-weight-bold mb-3"
            >Selected Submissions ({{ selectedItems.length }})</p
          >
          <v-chip-group>
            <v-chip
              v-for="item in selectedItems"
              :key="item.id"
              size="small"
              color="primary"
              variant="tonal"
            >
              {{ item.generated_id }}
            </v-chip>
          </v-chip-group>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            size="small"
            rounded
            color="grey"
            variant="outlined"
            @click="showBatchSubmissionDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn
            size="small"
            rounded
            color="primary"
            variant="flat"
            :disabled="selectedItems.length === 0"
            :loading="loading"
            @click="performBatchSubmission"
          >
            Submit {{ selectedItems.length }} Items
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Submission Detail Dialog -->
    <v-dialog v-model="showDetailDialog" max-width="900">
      <v-card v-if="selectedSubmission">
        <v-card-title class="text-h6 font-weight-bold bg-info text-white">
          <v-icon class="me-2">mdi-file-document-outline</v-icon>
          Submission Details - {{ selectedSubmission.bordereaux_id }}
        </v-card-title>
        <v-card-text class="pt-6">
          <v-row>
            <v-col cols="12" md="6">
              <h3 class="text-h6 font-weight-bold mb-3">Basic Information</h3>
              <v-list density="compact">
                <v-list-item>
                  <v-list-item-title>Type</v-list-item-title>
                  <v-list-item-subtitle>{{
                    formatBordereauType(selectedSubmission.type)
                  }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Scheme</v-list-item-title>
                  <v-list-item-subtitle>{{
                    selectedSubmission.scheme_name
                  }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Scheme</v-list-item-title>
                  <v-list-item-subtitle>{{
                    selectedSubmission.scheme_name
                  }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Record Count</v-list-item-title>
                  <v-list-item-subtitle>{{
                    selectedSubmission.record_count?.toLocaleString()
                  }}</v-list-item-subtitle>
                </v-list-item>
              </v-list>
            </v-col>
            <v-col cols="12" md="6">
              <h3 class="text-h6 font-weight-bold mb-3">Status & Timeline</h3>
              <v-timeline density="compact" direction="vertical">
                <v-timeline-item
                  v-for="event in selectedSubmission.timeline || []"
                  :key="event.date"
                  size="small"
                  :dot-color="getEventColor(event.type)"
                >
                  <template #opposite>
                    <span class="text-caption">{{
                      formatDateTime(event.date)
                    }}</span>
                  </template>
                  <div>
                    <p class="text-body-2 font-weight-bold">{{
                      event.title
                    }}</p>
                    <p class="text-caption text-medium-emphasis">{{
                      event.description
                    }}</p>
                  </div>
                </v-timeline-item>
              </v-timeline>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="info"
            variant="outlined"
            prepend-icon="mdi-table-eye"
            @click="viewBordereauData(selectedSubmission)"
          >
            View Data
          </v-btn>
          <v-btn
            color="success"
            variant="outlined"
            prepend-icon="mdi-download"
            @click="downloadFile(selectedSubmission)"
          >
            Download
          </v-btn>
          <v-btn
            v-if="selectedSubmission.status === 'submitted'"
            color="warning"
            variant="outlined"
            prepend-icon="mdi-email-sync"
            @click="checkStatus(selectedSubmission)"
          >
            Check Status
          </v-btn>
          <v-spacer></v-spacer>
          <v-btn
            color="grey"
            variant="outlined"
            @click="showDetailDialog = false"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Scheme List Dialog -->
    <v-dialog v-model="showSchemeDialog" max-width="600">
      <v-card v-if="selectedSchemeSubmission">
        <v-card-title class="text-h6 font-weight-bold bg-blue text-white">
          <v-icon class="me-2">mdi-list-box</v-icon>
          Schemes for
          {{
            selectedSchemeSubmission.bordereaux_id ||
            selectedSchemeSubmission.generated_id
          }}
        </v-card-title>
        <v-card-text class="pt-4">
          <v-alert color="info" variant="tonal" class="mb-4">
            <p class="mb-0"
              >All schemes associated with this bordereaux submission:</p
            >
          </v-alert>

          <v-list>
            <v-list-item
              v-for="(scheme, index) in getSchemeList(selectedSchemeSubmission)"
              :key="index"
              class="px-0"
            >
              <template #prepend>
                <v-icon color="primary">mdi-format-list-bulleted</v-icon>
              </template>
              <v-list-item-title>{{ scheme }}</v-list-item-title>
            </v-list-item>
          </v-list>

          <v-divider class="my-4"></v-divider>

          <div
            class="d-flex align-center justify-space-between text-caption text-medium-emphasis"
          >
            <span
              >Total Schemes:
              {{ getSchemeList(selectedSchemeSubmission).length }}</span
            >
            <span
              >Bordereaux ID:
              {{
                selectedSchemeSubmission.bordereaux_id ||
                selectedSchemeSubmission.generated_id
              }}</span
            >
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="grey"
            variant="outlined"
            @click="showSchemeDialog = false"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Review Dialog -->
    <v-dialog v-model="showReviewDialog" max-width="500">
      <v-card v-if="selectedApprovalItem">
        <v-card-title class="text-h6 font-weight-bold bg-teal text-white">
          <v-icon class="me-2">mdi-check-circle-outline</v-icon>
          Mark as Reviewed
        </v-card-title>
        <v-card-text class="pt-6">
          <v-alert color="teal" variant="tonal" class="mb-4">
            <p class="mb-0">
              Mark <strong>{{ selectedApprovalItem.generated_id }}</strong> as
              reviewed. This will allow it to be approved before submission.
            </p>
          </v-alert>
          <v-textarea
            v-model="approvalNotes"
            label="Review Notes (optional)"
            variant="outlined"
            density="compact"
            rows="3"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            size="small"
            rounded
            color="grey"
            variant="outlined"
            @click="showReviewDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn
            size="small"
            rounded
            color="teal"
            variant="flat"
            :loading="loading"
            @click="performReview"
          >
            Confirm Review
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Approve Dialog -->
    <v-dialog v-model="showApproveDialog" max-width="500">
      <v-card v-if="selectedApprovalItem">
        <v-card-title class="text-h6 font-weight-bold bg-success text-white">
          <v-icon class="me-2">mdi-thumb-up-outline</v-icon>
          Approve Bordereaux
        </v-card-title>
        <v-card-text class="pt-6">
          <v-alert color="success" variant="tonal" class="mb-4">
            <p class="mb-0">
              Approve
              <strong>{{ selectedApprovalItem.generated_id }}</strong> for
              submission. Once approved it can be batch-submitted to the
              reinsurer.
            </p>
          </v-alert>
          <v-textarea
            v-model="approvalNotes"
            label="Approval Notes (optional)"
            variant="outlined"
            density="compact"
            rows="3"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            size="small"
            rounded
            color="grey"
            variant="outlined"
            @click="showApproveDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn
            size="small"
            rounded
            color="success"
            variant="flat"
            :loading="loading"
            @click="performApprove"
          >
            Approve
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Return to Draft Dialog -->
    <v-dialog v-model="showReturnDialog" max-width="500">
      <v-card v-if="selectedApprovalItem">
        <v-card-title class="text-h6 font-weight-bold bg-warning text-white">
          <v-icon class="me-2">mdi-undo</v-icon>
          Return to Draft
        </v-card-title>
        <v-card-text class="pt-6">
          <v-alert color="warning" variant="tonal" class="mb-4">
            <p class="mb-0">
              Return <strong>{{ selectedApprovalItem.generated_id }}</strong> to
              draft status for rework.
            </p>
          </v-alert>
          <v-textarea
            v-model="returnReason"
            label="Reason for returning to draft *"
            variant="outlined"
            density="compact"
            rows="3"
            :rules="[(v) => !!v || 'Reason is required']"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            size="small"
            rounded
            color="grey"
            variant="outlined"
            @click="showReturnDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn
            size="small"
            rounded
            color="warning"
            variant="flat"
            :loading="loading"
            :disabled="!returnReason.trim()"
            @click="performReturnToDraft"
          >
            Return to Draft
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Data View Dialog -->
    <v-dialog v-model="showDataViewDialog" max-width="1200" scrollable>
      <v-card v-if="selectedSubmission">
        <v-card-title class="text-h6 font-weight-bold bg-info text-white">
          <v-icon class="me-2">mdi-table-eye</v-icon>
          Bordereau Data -
          {{
            selectedSubmission.bordereaux_id || selectedSubmission.generated_id
          }}
        </v-card-title>
        <v-card-text class="pt-4" style="height: 70vh">
          <v-alert color="info" variant="tonal" class="mb-4">
            <div class="d-flex align-center justify-space-between">
              <div>
                <p class="font-weight-bold mb-1"
                  >{{ formatBordereauType(selectedSubmission.type) }} Data</p
                >
                <p class="mb-0 text-caption"
                  >Total Records:
                  {{ bordereauData.length?.toLocaleString() }}</p
                >
              </div>
              <v-chip
                size="small"
                :color="getStatusColor(selectedSubmission.status)"
              >
                {{ formatStatus(selectedSubmission.status) }}
              </v-chip>
            </div>
          </v-alert>

          <div v-if="loading" class="text-center py-8">
            <v-progress-circular
              indeterminate
              color="primary"
              size="64"
            ></v-progress-circular>
            <p class="mt-4 text-body-2">Loading bordereau data...</p>
          </div>

          <div v-else-if="bordereauData.length === 0" class="text-center py-8">
            <v-icon size="64" color="grey-lighten-1">mdi-table-off</v-icon>
            <p class="mt-4 text-body-1">No data available</p>
          </div>

          <DataGrid
            v-else
            :row-data="bordereauData"
            :column-defs="dataColumnDefs"
            :table-title="`${formatBordereauType(selectedSubmission.type)} Data`"
            :show-export="true"
            density="compact"
            :pagination="true"
            :page-size="100"
            style="height: calc(70vh - 120px)"
          />
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="success"
            variant="outlined"
            prepend-icon="mdi-download"
            :disabled="bordereauData.length === 0"
            @click="downloadFile(selectedSubmission)"
          >
            Download
          </v-btn>
          <v-spacer></v-spacer>
          <v-btn color="grey" variant="outlined" @click="closeDataView">
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useFlashStore } from '@/renderer/store/flash'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

// Initialize stores
const flash = useFlashStore()
const { hasPermission } = usePermissionCheck()

// Reactive data
const loading = ref(false)
const selectedItems: any = ref([])
const showBatchSubmissionDialog = ref(false)
const showDetailDialog = ref(false)
const showSchemeDialog = ref(false)
const selectedSubmission: any = ref(null)
const selectedSchemeSubmission: any = ref(null)
const showDataViewDialog = ref(false)
const bordereauData: any = ref([])
const dataColumnDefs: any = ref([])

// Approval workflow
const selectedApprovalItem: any = ref(null)
const showReviewDialog = ref(false)
const showApproveDialog = ref(false)
const showReturnDialog = ref(false)
const approvalNotes = ref('')
const returnReason = ref('')

const filters: any = ref({
  status: [],
  scheme: null,
  type: '',
  search: ''
})

const batchSubmission = ref({
  delivery_method: 'email',
  message: ''
})

// Static data
const columnDefs = ref([
  {
    headerName: '',
    checkboxSelection: true,
    headerCheckboxSelection: true,
    width: 50,
    pinned: 'left',
    lockPosition: true,
    suppressMenu: true,
    suppressSorting: true,
    suppressFilter: true,
    suppressResize: true
  },
  {
    headerName: 'Bordereaux ID',
    field: 'generated_id',
    sortable: true,
    filter: true
  },
  {
    headerName: 'Type',
    field: 'type',
    sortable: true,
    filter: true,
    cellRenderer: (params: any) => {
      const color = getBordereauTypeColor(params.value)
      const label = formatBordereauType(params.value)
      return `<span style="padding:2px 10px;border-radius:12px;font-size:12px;background:${color};color:#fff;font-weight:500">${label}</span>`
    }
  },
  {
    headerName: 'Scheme',
    field: 'scheme_name',
    sortable: true,
    filter: true,
    cellRenderer: (params: any) => {
      // Get scheme data - could be scheme_name field or schemes array
      let schemesData =
        params.data.schemes || params.data.scheme_name || params.value

      // Handle string with comma-separated values (like "DevX, KoreBiz, Blah1, Blah2, Blah3, BBlah4")
      if (typeof schemesData === 'string') {
        schemesData = schemesData.split(',').map((s) => s.trim())
      }

      // Ensure we have an array
      const schemes = Array.isArray(schemesData) ? schemesData : [schemesData]

      if (schemes.length <= 3) {
        return schemes.join(', ')
      } else {
        const truncated = schemes.slice(0, 3).join(', ')
        return `<span class="scheme-cell" style="cursor: pointer; color: #1976d2; text-decoration: underline;" onclick="window.showSchemes_${params.data.id}()">${truncated} ... (+${schemes.length - 3} more)</span>`
      }
    }
  },
  {
    headerName: 'Status',
    field: 'status',
    sortable: true,
    filter: true,
    cellRenderer: (params: any) => {
      const color = getStatusColor(params.value)
      const label = formatStatus(params.value)
      return `<span style="padding:2px 10px;border-radius:12px;font-size:12px;background:${color};color:#fff;font-weight:500">${label}</span>`
    }
  },
  {
    headerName: 'Progress',
    field: 'progress',
    sortable: true,
    filter: true,
    cellRenderer: (params: any) => {
      const color = getProgressColor(params.value)
      return `<div style="display: flex; align-items: center; gap: 8px;"><div style="width: 60px; height: 6px; background-color: #e0e0e0; border-radius: 3px; overflow: hidden;"><div style="width: ${params.value}%; height: 100%; background-color: ${color}; transition: width 0.3s ease;"></div></div><span style="font-size: 12px;">${params.value}%</span></div>`
    }
  },
  {
    headerName: 'Submitted',
    field: 'submission_date',
    sortable: true,
    filter: true,
    cellRenderer: (params: any) => formatDateTime(params.value)
  },
  {
    headerName: 'Updated',
    field: 'last_updated',
    sortable: true,
    filter: true,
    cellRenderer: (params: any) => formatDateTime(params.value)
  },
  {
    headerName: 'SLA',
    field: 'sla_status',
    sortable: true,
    filter: true,
    cellRenderer: (params: any) => {
      const color = getSLAColor(params.value)
      return `<span style="padding:2px 10px;border-radius:12px;font-size:12px;background:${color};color:#fff;font-weight:500">${params.value}</span>`
    }
  },
  {
    headerName: 'Actions',
    field: 'actions',
    sortable: false,
    filter: false,
    width: 80,
    pinned: 'right',
    cellRenderer: (params: any) => {
      const key = String(params.data.id).replace(/-/g, '_')
      ;(window as any)[`showSubmissionMenu_${key}`] = (event: MouseEvent) =>
        showContextMenu(event, params.data)
      return `<div style="display:flex;align-items:center;justify-content:center;height:100%">
        <button onclick="showSubmissionMenu_${key}(event)" title="Actions" style="background:none;border:none;cursor:pointer;padding:4px 8px;border-radius:4px;color:#616161;display:flex;align-items:center;justify-content:center">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="5" r="2"/><circle cx="12" cy="12" r="2"/><circle cx="12" cy="19" r="2"/></svg>
        </button>
      </div>`
    }
  }
])

const statusOptions = [
  { title: 'Draft', value: 'draft' },
  { title: 'Generated', value: 'generated' },
  { title: 'Reviewed', value: 'reviewed' },
  { title: 'Approved', value: 'approved' },
  { title: 'Submitted', value: 'submitted' },
  { title: 'Delivered', value: 'delivered' },
  { title: 'Confirmed', value: 'confirmed' },
  { title: 'Processing', value: 'processing' },
  { title: 'Reconciled', value: 'reconciled' },
  { title: 'Failed', value: 'failed' },
  { title: 'Cancelled', value: 'cancelled' }
]

const bordereauTypes = [
  { title: 'Member', value: 'member' },
  { title: 'Premium', value: 'premium' },
  { title: 'Claims', value: 'claims' },
  { title: 'Benefits', value: 'benefits' }
]

const deliveryMethods = [
  { title: 'Email', value: 'email' },
  { title: 'SFTP Upload', value: 'sftp' },
  { title: 'API Submission', value: 'api' },
  { title: 'Portal Upload', value: 'portal' }
]

const schemes = ref([
  { id: 1, name: 'Liberty Life' },
  { id: 2, name: 'Old Mutual' },
  { id: 3, name: 'Momentum' },
  { id: 4, name: 'Discovery Life' }
])

const submissions: any = ref([])

// Computed properties
const filteredSubmissions = computed(() => {
  let filtered = submissions.value

  if (filters.value.status.length > 0) {
    filtered = filtered.filter((item) =>
      filters.value.status.includes(item.status)
    )
  }

  if (filters.value.scheme) {
    filtered = filtered.filter(
      (item) =>
        item.scheme_name ===
        schemes.value.find((i) => i.id === filters.value.scheme)?.name
    )
  }

  if (filters.value.type) {
    filtered = filtered.filter((item) => item.type === filters.value.type)
  }

  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    filtered = filtered.filter(
      (item) =>
        item.bordereaux_id.toLowerCase().includes(search) ||
        item.scheme_name.toLowerCase().includes(search)
    )
  }

  return filtered
})

// Methods
const getBordereauTypeColor = (type: string): string => {
  const colors: Record<string, string> = {
    member: '#1976d2',
    premium: '#388e3c',
    claims: '#f57c00',
    benefits: '#7b1fa2'
  }
  return colors[type] || '#757575'
}

const formatBordereauType = (type: string): string => {
  const types: Record<string, string> = {
    member: 'Member',
    premium: 'Premium',
    claims: 'Claims',
    benefits: 'Benefits'
  }
  return types[type] || type
}

const getStatusColor = (status: string): string => {
  status = status.toLowerCase()
  const colors: Record<string, string> = {
    draft: '#757575',
    generated: '#1976d2',
    reviewed: '#00796b',
    approved: '#388e3c',
    submitted: '#f57c00',
    delivered: '#00796b',
    confirmed: '#388e3c',
    processing: '#0288d1',
    reconciled: '#388e3c',
    failed: '#d32f2f',
    cancelled: '#f57c00'
  }
  return colors[status] || '#757575'
}

const formatStatus = (status: string): string => {
  return status.charAt(0).toUpperCase() + status.slice(1)
}

const getProgressColor = (progress: number): string => {
  if (progress === 100) return '#388e3c'
  if (progress >= 75) return '#0288d1'
  if (progress >= 50) return '#f57c00'
  return '#d32f2f'
}

const getSLAColor = (sla: string): string => {
  const colors: Record<string, string> = {
    'On Time': '#388e3c',
    Warning: '#f57c00',
    Overdue: '#d32f2f',
    Failed: '#d32f2f'
  }
  return colors[sla] || '#757575'
}

const getEventColor = (type: string): string => {
  const colors: Record<string, string> = {
    generated: 'blue',
    submitted: 'orange',
    delivered: 'teal',
    confirmed: 'success',
    processing: 'info',
    failed: 'error'
  }
  return colors[type] || 'grey'
}

const formatDateTime = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('en-ZA', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const refreshSubmissions = async () => {
  loading.value = true
  try {
    const response = await GroupPricingService.getBordereauxActivity()
    submissions.value = response.data || []
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to fetch bordereaux submissions',
      'error'
    )
  } finally {
    loading.value = false
  }
}

const viewSubmission = (item: any) => {
  selectedSubmission.value = item
  showDetailDialog.value = true
}

const submitToScheme = async (item: any) => {
  if (!item?.id) return
  if (item.status !== 'approved') {
    flash.show(
      'Bordereaux must be approved before it can be submitted to the scheme',
      'warning'
    )
    return
  }
  try {
    loading.value = true
    const response = await GroupPricingService.submitBordereauxBatch({
      bordereaux_ids: [item.id],
      delivery_method: 'email',
      message: '',
      submission_date: new Date().toISOString()
    })
    if (response.status === 200) {
      const idx = submissions.value.findIndex((s) => s.id === item.id)
      if (idx !== -1) {
        submissions.value[idx].status = 'submitted'
        submissions.value[idx].submission_date = new Date().toISOString()
        submissions.value[idx].last_updated = new Date().toISOString()
        submissions.value[idx].progress = 66
      }
      flash.show(
        `Bordereaux ${item.generated_id || item.id} submitted to scheme`,
        'success'
      )
      await refreshSubmissions()
    }
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to submit bordereaux to scheme',
      'error'
    )
  } finally {
    loading.value = false
  }
}

const checkStatus = async (item: any) => {
  const generatedId = item?.generated_id
  if (!generatedId) {
    flash.show('Cannot check status: missing generated_id', 'error')
    return
  }
  try {
    loading.value = true
    const response = await GroupPricingService.getBordereauxById(generatedId)
    const fresh = response.data
    if (fresh) {
      const idx = submissions.value.findIndex((s) => s.id === item.id)
      if (idx !== -1) {
        submissions.value[idx] = { ...submissions.value[idx], ...fresh }
      }
      if (selectedSubmission.value?.id === item.id) {
        selectedSubmission.value = {
          ...selectedSubmission.value,
          ...fresh
        }
      }
      flash.show(
        `Status: ${fresh.status}${
          typeof fresh.progress === 'number' ? ` (${fresh.progress}%)` : ''
        }`,
        'info'
      )
    }
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to fetch bordereaux status',
      'error'
    )
  } finally {
    loading.value = false
  }
}

const downloadFile = async (item: any) => {
  try {
    loading.value = true

    // Show downloading notification
    flash.show(
      `Downloading ${item.generated_id || item.bordereaux_id}...`,
      'info'
    )

    // Make API call to get file blob using file_name
    const blob = await GroupPricingService.downloadBordereaux(item.file_name)

    // Check if we received a blob
    if (blob.data instanceof Blob) {
      // Create download link
      const url = window.URL.createObjectURL(blob.data)
      const link = document.createElement('a')
      link.href = url
      link.download =
        item.file_name || `${item.generated_id || item.bordereaux_id}.xlsx`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)

      // Show success message
      flash.show(
        `Successfully downloaded ${item.generated_id || item.bordereaux_id}`,
        'success'
      )
    } else {
      throw new Error('Invalid file data received from server')
    }
  } catch (error: any) {
    console.error('Download failed:', error)

    // Show user-friendly error message
    let errorMessage = `Failed to download ${item.generated_id || item.bordereaux_id}`
    if (error.response?.data?.message) {
      errorMessage = error.response.data.message
    } else if (
      error.message &&
      error.message !== 'Invalid file data received from server'
    ) {
      errorMessage += `: ${error.message}`
    }

    flash.show(errorMessage, 'error')
  } finally {
    loading.value = false
  }
}

const viewBordereauData = async (item: any) => {
  console.log('Viewing bordereau data for:', item)
  try {
    loading.value = true
    bordereauData.value = []
    dataColumnDefs.value = []

    // Show loading message
    flash.show(
      `Loading data for ${item.generated_id || item.bordereaux_id}...`,
      'info'
    )

    // Make API call to get bordereau data
    const response = await GroupPricingService.getBordereauxData(
      item.generated_id
    )

    console.log('Bordereau data response:', response.data)

    if (response.data && Array.isArray(response.data)) {
      bordereauData.value = response.data

      // Generate column definitions dynamically from the first record
      if (bordereauData.value.length > 0) {
        const firstRecord = bordereauData.value[0]
        dataColumnDefs.value = Object.keys(firstRecord).map((key) => ({
          headerName: key
            .replace(/_/g, ' ')
            .replace(/\b\w/g, (l) => l.toUpperCase()),
          field: key,
          sortable: true,
          filter: true,
          resizable: true,
          width: getColumnWidth(key, firstRecord[key])
        }))
      }

      // Show the data dialog
      showDataViewDialog.value = true
      flash.show(`Loaded ${bordereauData.value.length} records`, 'success')
    } else {
      throw new Error('No data found for this bordereau')
    }
  } catch (error: any) {
    console.error('Failed to load bordereau data:', error)

    // Show user-friendly error message
    let errorMessage = `Failed to load data for ${item.generated_id || item.bordereaux_id}`
    if (error.response?.data?.message) {
      errorMessage = error.response.data.message
    } else if (
      error.message &&
      error.message !== 'No data found for this bordereau'
    ) {
      errorMessage += `: ${error.message}`
    }

    flash.show(errorMessage, 'error')
  } finally {
    loading.value = false
  }
}

const closeDataView = () => {
  showDataViewDialog.value = false
  bordereauData.value = []
  dataColumnDefs.value = []
}

const getColumnWidth = (columnName: string, sampleValue: any): number => {
  // Set appropriate column widths based on content type and name
  const name = columnName.toLowerCase()

  if (name.includes('id') || name.includes('reference')) return 150
  if (name.includes('name') || name.includes('description')) return 200
  if (name.includes('email') || name.includes('address')) return 250
  if (name.includes('date') || name.includes('time')) return 140
  if (
    name.includes('amount') ||
    name.includes('premium') ||
    name.includes('value')
  )
    return 130
  if (name.includes('status') || name.includes('type')) return 120

  // Default width based on sample value type
  if (typeof sampleValue === 'number') return 100
  if (typeof sampleValue === 'boolean') return 80

  return 150
}

const performBatchSubmission = async () => {
  if (selectedItems.value.length === 0) {
    return
  }

  // Validate required fields
  if (!batchSubmission.value.delivery_method) {
    flash.show('Please select a delivery method', 'error')
    return
  }

  try {
    // Show loading state
    loading.value = true

    // Prepare submission data
    const submissionData = {
      bordereaux_ids: selectedItems.value.map((item) => item.id),
      delivery_method: batchSubmission.value.delivery_method,
      message: batchSubmission.value.message || '',
      submission_date: new Date().toISOString()
    }

    console.log('Submitting batch with data:', submissionData)

    // Make API call to submit batch
    const response =
      await GroupPricingService.submitBordereauxBatch(submissionData)

    if (response.data && response.data.success) {
      // Update local data to reflect submission status
      selectedItems.value.forEach((item) => {
        const submissionIndex = submissions.value.findIndex(
          (s) => s.id === item.id
        )
        if (submissionIndex !== -1) {
          submissions.value[submissionIndex].status = 'submitted'
          submissions.value[submissionIndex].submission_date =
            new Date().toISOString()
          submissions.value[submissionIndex].last_updated =
            new Date().toISOString()
          submissions.value[submissionIndex].progress = 50

          // Add timeline event
          if (!submissions.value[submissionIndex].timeline) {
            submissions.value[submissionIndex].timeline = []
          }
          submissions.value[submissionIndex].timeline.push({
            date: new Date().toISOString(),
            type: 'submitted',
            title: 'Batch Submitted',
            description: `Submitted via ${batchSubmission.value.delivery_method} with ${selectedItems.value.length - 1} other bordereaux`
          })
        }
      })

      // Show success message
      flash.show(
        `Successfully submitted ${selectedItems.value.length} bordereaux via ${batchSubmission.value.delivery_method}`,
        'success'
      )

      // Reset form and close dialog
      batchSubmission.value = {
        delivery_method: 'email',
        message: ''
      }
      showBatchSubmissionDialog.value = false
      selectedItems.value = []
    } else {
      // Handle API error response
      const errorMessage = response.data?.message || 'Failed to submit batch'
      flash.show(errorMessage, 'error')
    }
  } catch (error: any) {
    console.error('Batch submission failed:', error)

    // Show user-friendly error message
    let errorMessage = 'Failed to submit batch. Please try again.'
    if (error.response?.data?.message) {
      errorMessage = error.response.data.message
    } else if (error.message) {
      errorMessage = error.message
    }

    flash.show(errorMessage, 'error')
  } finally {
    loading.value = false
  }
}

// Approval workflow actions
const reviewItem = (item: any) => {
  selectedApprovalItem.value = item
  approvalNotes.value = ''
  showReviewDialog.value = true
}

const approveItem = (item: any) => {
  selectedApprovalItem.value = item
  approvalNotes.value = ''
  showApproveDialog.value = true
}

const returnToDraftItem = (item: any) => {
  selectedApprovalItem.value = item
  returnReason.value = ''
  showReturnDialog.value = true
}

const performReview = async () => {
  if (!selectedApprovalItem.value) return
  try {
    loading.value = true
    const response = await GroupPricingService.reviewGeneratedBordereaux(
      selectedApprovalItem.value.generated_id,
      { notes: approvalNotes.value }
    )
    if (response.data?.success) {
      flash.show('Bordereaux marked as reviewed', 'success')
      showReviewDialog.value = false
      approvalNotes.value = ''
      await refreshSubmissions()
    }
  } catch (error: any) {
    flash.show(
      error.response?.data?.error || 'Failed to review bordereaux',
      'error'
    )
  } finally {
    loading.value = false
  }
}

const performApprove = async () => {
  if (!selectedApprovalItem.value) return
  try {
    loading.value = true
    const response = await GroupPricingService.approveGeneratedBordereaux(
      selectedApprovalItem.value.generated_id,
      { notes: approvalNotes.value }
    )
    if (response.data?.success) {
      flash.show('Bordereaux approved', 'success')
      showApproveDialog.value = false
      approvalNotes.value = ''
      await refreshSubmissions()
    }
  } catch (error: any) {
    flash.show(
      error.response?.data?.error || 'Failed to approve bordereaux',
      'error'
    )
  } finally {
    loading.value = false
  }
}

const performReturnToDraft = async () => {
  if (!selectedApprovalItem.value || !returnReason.value.trim()) {
    flash.show('A reason is required', 'error')
    return
  }
  try {
    loading.value = true
    const response = await GroupPricingService.returnBordereauxToDraft(
      selectedApprovalItem.value.generated_id,
      { reason: returnReason.value }
    )
    if (response.data?.success) {
      flash.show('Bordereaux returned to draft', 'success')
      showReturnDialog.value = false
      returnReason.value = ''
      await refreshSubmissions()
    }
  } catch (error: any) {
    flash.show(
      error.response?.data?.error || 'Failed to return bordereaux to draft',
      'error'
    )
  } finally {
    loading.value = false
  }
}

const onRowClicked = (event: any) => {
  // Handle row click from DataGrid
  const rowData = event.data
  console.log('Row clicked:', rowData)
}

const onRowSelectionChanged = (selectedRows: any[]) => {
  // Update selectedItems when row selection changes
  selectedItems.value = selectedRows
  console.log('Selected rows:', selectedRows.length)
}

let activeSubmissionMenuCleanup: (() => void) | null = null

function showContextMenu(event: MouseEvent, data: any) {
  if (activeSubmissionMenuCleanup) activeSubmissionMenuCleanup()

  const isDraftOrGenerated =
    data.status === 'draft' || data.status === 'generated'
  const isReviewed = data.status === 'reviewed'
  const isApproved = data.status === 'approved'
  const isSubmitted = data.status === 'submitted'

  const menuItems = [
    { label: 'View', color: '#0288d1', fn: () => viewSubmission(data) },
    isDraftOrGenerated && hasPermission('bordereaux:approve_outbound')
      ? {
          label: 'Mark as Reviewed',
          color: '#1976d2',
          fn: () => reviewItem(data)
        }
      : null,
    isReviewed && hasPermission('bordereaux:approve_outbound')
      ? { label: 'Approve', color: '#388e3c', fn: () => approveItem(data) }
      : null,
    isApproved && hasPermission('bordereaux:submit_outbound')
      ? { label: 'Submit', color: '#1976d2', fn: () => submitToScheme(data) }
      : null,
    isReviewed || isApproved
      ? {
          label: 'Return to Draft',
          color: '#f57c00',
          fn: () => returnToDraftItem(data)
        }
      : null,
    isSubmitted
      ? { label: 'Check Status', color: '#f57c00', fn: () => checkStatus(data) }
      : null,
    { label: 'Download', color: '#388e3c', fn: () => downloadFile(data) },
    { label: 'View Data', color: '#616161', fn: () => viewBordereauData(data) }
  ].filter(Boolean) as { label: string; color: string; fn: () => void }[]

  const menu = document.createElement('div')
  menu.style.cssText =
    'position:fixed;background:#fff;border:1px solid #e0e0e0;border-radius:8px;' +
    'box-shadow:0 4px 16px rgba(0,0,0,0.14);z-index:9999;min-width:180px;padding:4px 0;'

  menuItems.forEach(({ label, color, fn }) => {
    const item = document.createElement('div')
    item.textContent = label
    item.style.cssText = `padding:8px 16px;cursor:pointer;font-size:13px;color:${color};`
    item.addEventListener(
      'mouseenter',
      () => (item.style.background = '#f5f5f5')
    )
    item.addEventListener('mouseleave', () => (item.style.background = ''))
    item.addEventListener('click', () => {
      cleanup()
      fn()
    })
    menu.appendChild(item)
  })

  document.body.appendChild(menu)

  const btn = (event.currentTarget || event.target) as HTMLElement
  const rect = btn.getBoundingClientRect()
  menu.style.top = `${rect.bottom + 4}px`
  menu.style.left = `${rect.left}px`

  const mr = menu.getBoundingClientRect()
  if (mr.right > window.innerWidth - 8)
    menu.style.left = `${rect.right - mr.width}px`
  if (mr.bottom > window.innerHeight - 8)
    menu.style.top = `${rect.top - mr.height - 4}px`

  function cleanup() {
    menu.remove()
    document.removeEventListener('click', outsideClick, true)
    activeSubmissionMenuCleanup = null
  }
  activeSubmissionMenuCleanup = cleanup

  function outsideClick(e: MouseEvent) {
    if (!menu.contains(e.target as Node) && e.target !== btn) cleanup()
  }
  setTimeout(() => document.addEventListener('click', outsideClick, true), 0)
}

// Setup global functions for action buttons in AG Grid
const setupGlobalActions = () => {
  filteredSubmissions.value.forEach((item) => {
    ;(window as any)[`showSchemes_${item.id}`] = () => showSchemes(item)
  })
}

const showSchemes = (item: any) => {
  selectedSchemeSubmission.value = item
  showSchemeDialog.value = true
}

const getSchemeList = (submission: any) => {
  if (!submission) return []

  const schemesData = submission.schemes || submission.scheme_name

  // Handle string with comma-separated values
  if (typeof schemesData === 'string') {
    return schemesData.split(',').map((s) => s.trim())
  }

  // Handle array
  if (Array.isArray(schemesData)) {
    return schemesData
  }

  // Fallback to single item array
  return [schemesData]
}

// Watch for changes in filtered submissions to setup global actions
watch(
  () => filteredSubmissions.value,
  () => {
    setupGlobalActions()
  },
  { deep: true }
)

onMounted(async () => {
  await refreshSubmissions()
  setupGlobalActions()
})
</script>

<style scoped>
.progress-container {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 100px;
}

/* AG Grid action button styles */
:deep(.ag-actions-cell) {
  display: flex;
  align-items: center;
  gap: 4px;
  height: 100%;
}

:deep(.ag-btn) {
  border: none;
  background: transparent;
  cursor: pointer;
  padding: 4px 6px;
  border-radius: 4px;
  font-size: 14px;
  min-width: 24px;
  min-height: 24px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
}

:deep(.ag-btn-info) {
  color: rgb(var(--v-theme-info));
}

:deep(.ag-btn-info:hover) {
  background-color: rgba(var(--v-theme-info), 0.1);
}

:deep(.ag-btn-primary) {
  color: rgb(var(--v-theme-primary));
}

:deep(.ag-btn-primary:hover) {
  background-color: rgba(var(--v-theme-primary), 0.1);
}

:deep(.ag-btn-warning) {
  color: rgb(var(--v-theme-warning));
}

:deep(.ag-btn-warning:hover) {
  background-color: rgba(var(--v-theme-warning), 0.1);
}

:deep(.ag-btn-success) {
  color: rgb(var(--v-theme-success));
}

:deep(.ag-btn-success:hover) {
  background-color: rgba(var(--v-theme-success), 0.1);
}

:deep(.ag-btn-grey) {
  color: rgb(var(--v-theme-on-surface));
}

:deep(.ag-btn-grey:hover) {
  background-color: rgba(var(--v-theme-on-surface), 0.1);
}
</style>
