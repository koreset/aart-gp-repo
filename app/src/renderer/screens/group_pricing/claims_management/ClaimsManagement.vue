<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex justify-space-between align-center flex-wrap">
              <span class="headline">Claims Management</span>
              <div class="d-flex gap-2">
                <v-btn
                  size="small"
                  variant="outlined"
                  class="mr-2"
                  rounded
                  prepend-icon="mdi-upload"
                  @click="bulkUploadDialog = true"
                >
                  Bulk Upload
                </v-btn>
                <v-btn
                  v-if="hasPermission('claims:view_analytics')"
                  size="small"
                  variant="outlined"
                  class="mr-2"
                  rounded
                  prepend-icon="mdi-chart-line"
                  @click="
                    router.push({ name: 'group-pricing-claims-analytics' })
                  "
                >
                  Analytics
                </v-btn>
                <v-btn
                  v-if="hasPermission('claims_pay:create_schedule')"
                  size="small"
                  variant="outlined"
                  class="mr-2"
                  rounded
                  prepend-icon="mdi-cash-check"
                  @click="
                    router.push({
                      name: 'group-pricing-claim-payment-schedules'
                    })
                  "
                >
                  Payment Schedules
                </v-btn>
                <v-btn
                  v-if="hasPermission('claims:lodge')"
                  size="small"
                  variant="outlined"
                  rounded
                  prepend-icon="mdi-file-plus"
                  @click="newClaimDialog = true"
                >
                  New Claim
                </v-btn>
              </div>
            </div>
          </template>
          <template #default>
            <!-- Search and Filter Bar -->
            <v-row class="mb-4">
              <v-col cols="12" md="3">
                <v-text-field
                  v-model="searchQuery"
                  label="Search Claims"
                  prepend-inner-icon="mdi-magnify"
                  variant="outlined"
                  density="compact"
                  clearable
                />
              </v-col>
              <v-col cols="12" md="3">
                <v-select
                  v-model="selectedStatus"
                  :items="claimStatuses"
                  label="Filter by Status"
                  item-title="text"
                  item-value="value"
                  variant="outlined"
                  density="compact"
                  clearable
                />
              </v-col>
              <v-col cols="12" md="3">
                <v-select
                  v-model="selectedBenefitType"
                  :items="benefitTypes"
                  label="Filter by Benefit Type"
                  item-title="title"
                  item-value="value"
                  variant="outlined"
                  density="compact"
                  clearable
                />
              </v-col>
              <v-col cols="12" md="3">
                <v-select
                  v-model="selectedScheme"
                  :items="schemes"
                  label="Filter by Scheme"
                  item-title="name"
                  item-value="id"
                  variant="outlined"
                  density="compact"
                  clearable
                />
              </v-col>
            </v-row>

            <!-- Quick Stats Cards -->
            <v-row class="mb-4">
              <v-col cols="12" sm="6" md="4" lg="2">
                <stat-card
                  title="Drafts"
                  :value="String(claimsStats.drafts)"
                  icon="mdi-file-edit-outline"
                  color="grey"
                  :loading="loading"
                />
              </v-col>
              <v-col cols="12" sm="6" md="4" lg="2">
                <stat-card
                  title="Pending Assessment"
                  :value="String(claimsStats.pendingAssessment)"
                  icon="mdi-clock-outline"
                  color="info"
                  :loading="loading"
                />
              </v-col>
              <v-col cols="12" sm="6" md="4" lg="2">
                <stat-card
                  title="Under Assessment"
                  :value="String(claimsStats.underAssessment)"
                  icon="mdi-magnify"
                  color="warning"
                  :loading="loading"
                />
              </v-col>
              <v-col cols="12" sm="6" md="4" lg="3">
                <stat-card
                  title="Approved This Month"
                  :value="String(claimsStats.approved)"
                  icon="mdi-check-circle-outline"
                  color="success"
                  :loading="loading"
                />
              </v-col>
              <v-col cols="12" sm="6" md="4" lg="3">
                <stat-card
                  title="Declined This Month"
                  :value="String(claimsStats.declined)"
                  icon="mdi-close-circle-outline"
                  color="error"
                  :loading="loading"
                />
              </v-col>
            </v-row>

            <!-- Claims Data Grid -->
            <div :style="{ height: gridHeight, width: '100%' }">
              <data-grid
                :column-defs="claimsColumnDefs"
                :row-data="filteredClaims"
                :loading="loading"
                style="height: 100%; width: 100%"
                @row-double-clicked="viewClaimDetails"
              />
            </div>
            <empty-state
              v-if="!loading && filteredClaims.length === 0"
              icon="mdi-clipboard-text-off-outline"
              title="No claims found"
              message="No claims match the current filters."
            />
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- New Claim Dialog -->
    <register-claim-dialog
      v-model="newClaimDialog"
      :schemes="schemes"
      @save="handleNewClaimSave"
      @cancel="newClaimDialog = false"
    />

    <!-- Bulk Claims Upload Dialog -->
    <v-dialog v-model="bulkUploadDialog" persistent max-width="900px">
      <base-card>
        <template #header>
          <span class="headline">Bulk Claims Upload</span>
        </template>
        <template #default>
          <bulk-claims-upload
            :schemes="schemes"
            @upload-complete="handleBulkUploadComplete"
            @cancel="bulkUploadDialog = false"
          />
        </template>
      </base-card>
    </v-dialog>

    <!-- Confirmation Dialog -->
    <v-dialog v-model="confirmDialog" persistent max-width="400px">
      <v-card>
        <v-card-title class="text-h6">{{ confirmTitle }}</v-card-title>
        <v-card-text>{{ confirmMessage }}</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="grey" variant="text" @click="confirmDialog = false"
            >Cancel</v-btn
          >
          <v-btn color="primary" @click="confirmAction">Confirm</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Snackbar for notifications -->
    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="3000">
      {{ snackbarMessage }}
      <template #actions>
        <v-btn color="white" variant="text" @click="snackbar = false"
          >Close</v-btn
        >
      </template>
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'
import StatCard from '@/renderer/components/StatCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import RegisterClaimDialog from './components/RegisterClaimDialog.vue'
import BulkClaimsUpload from './components/BulkClaimsUpload.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { statusCellRenderer } from '@/renderer/utils/statusCellRenderer'
import { currencyFormatter, dateFormatter } from '@/renderer/utils/formatters'
import { useGridHeight } from '@/renderer/composables/useGridHeight'
import { useStatusBarStore } from '@/renderer/store/statusBar'
import {
  isEditableClaimStatus,
  isSubmittableClaimStatus
} from '@/renderer/utils/claimStatus'

const gridHeight = useGridHeight(380)
const statusBarStore = useStatusBarStore()
const { hasPermission } = usePermissionCheck()

// Interfaces
interface Claim {
  id?: number
  claim_number: string
  member_name: string
  member_id_number: string
  scheme_name: string
  benefit_alias: string
  benefit_code: string
  benefit_name: string
  claim_amount: number
  date_of_event: string
  date_notified: string
  status: string
  assessor_name?: string
  priority: string
}

interface ClaimsStats {
  drafts: number
  pendingAssessment: number
  underAssessment: number
  approved: number
  declined: number
}

interface Scheme {
  id: number
  name: string
}

// State
const loading = ref(false)
const claims = ref<Claim[]>([])
const schemes = ref<Scheme[]>([])

const router = useRouter()

// Dialog states
const newClaimDialog = ref(false)
const bulkUploadDialog = ref(false)
const confirmDialog = ref(false)

// Filter states
const searchQuery = ref('')
const selectedStatus = ref<string | null>(null)
const selectedBenefitType: any = ref(null)
const selectedScheme = ref<number | null>(null)

const benefitMaps = ref<any>([])

// Confirmation
const confirmTitle = ref('')
const confirmMessage = ref('')
const confirmCallback = ref<(() => void) | null>(null)

// Snackbar
const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

// Options
const claimStatuses = [
  { text: 'Draft', value: 'draft' },
  { text: 'Pending Assessment', value: 'pending_assessment' },
  { text: 'Under Assessment', value: 'under_assessment' },
  { text: 'Additional Info Required', value: 'additional_info_required' },
  { text: 'Approved', value: 'approved' },
  { text: 'Submitted for Payment', value: 'submitted_for_payment' },
  { text: 'Paid', value: 'paid' },
  { text: 'Declined', value: 'declined' },
  { text: 'Cancelled', value: 'cancelled' }
]

const benefitTypes = computed(() => {
  if (!benefitMaps.value || benefitMaps.value.length === 0) {
    return []
  }

  return benefitMaps.value
    .filter((benefit: any) => benefit.is_mapped)
    .map((benefit: any) => {
      return {
        title: benefit.benefit_alias,
        value: benefit
      }
    })
})

// const benefitTypes = [
//   { text: 'Group Life Assurance (GLA)', value: 'Group Life Assurance (GLA)' },
//   { text: 'Spouse Group Life Assurance (SGLA)', value: 'Spouse Group Life Assurance (SGLA)' },
//   { text: 'Permanent Total Disability (PTD)', value: 'Permanent Total Disability (PTD)' },
//   { text: 'Critical Illness (CI)', value: 'Critical Illness (CI)' },
//   { text: 'Temporary Total Disability (TTD)', value: 'Temporary Total Disability (TTD)' },
//   { text: 'Personal Health Insurance (PHI)', value: 'Personal Health Insurance (PHI)' },
//   { text: 'Group Family Funeral (GFF)', value: 'Group Family Funeral (GFF)' }
// ]

// Column definitions for claims grid
const claimsColumnDefs = [
  {
    headerName: 'Claim Number',
    field: 'claim_number',
    filter: true,
    sortable: true,
    minWidth: 150,
    cellRenderer: (params: any) => {
      return `<span class="font-weight-medium text-primary cursor-pointer">${params.value}</span>`
    }
  },
  {
    headerName: 'Member',
    field: 'member_name',
    filter: true,
    sortable: true,
    minWidth: 180
  },
  {
    headerName: 'ID Number',
    field: 'member_id_number',
    filter: true,
    sortable: true,
    minWidth: 150
  },
  {
    headerName: 'Scheme',
    field: 'scheme_name',
    filter: true,
    sortable: true,
    minWidth: 150
  },
  {
    headerName: 'Benefit Type',
    field: 'benefit_alias',
    filter: true,
    sortable: true,
    minWidth: 180
  },
  {
    headerName: 'Member Type',
    field: 'member_type',
    filter: true,
    sortable: true,
    minWidth: 180
  },
  {
    headerName: 'Claim Amount',
    field: 'claim_amount',
    filter: true,
    sortable: true,
    minWidth: 130,
    valueFormatter: currencyFormatter
  },
  {
    headerName: 'Date of Event',
    field: 'date_of_event',
    filter: true,
    sortable: true,
    minWidth: 130,
    valueFormatter: dateFormatter
  },
  {
    headerName: 'Status',
    field: 'status',
    filter: true,
    sortable: true,
    minWidth: 150,
    cellRenderer: (params: any) => statusCellRenderer(params.value)
  },
  {
    headerName: 'Priority',
    field: 'priority',
    filter: true,
    sortable: true,
    minWidth: 100,
    cellRenderer: (params: any) => statusCellRenderer(params.value)
  },
  {
    headerName: 'Date of Creation',
    field: 'creation_date',
    filter: true,
    sortable: true,
    minWidth: 130
  },
  {
    headerName: 'Actions',
    pinned: 'right' as const,
    width: 220,
    sortable: false,
    filter: false,
    resizable: false,
    cellRenderer: (params: any) => {
      const viewBtn = `
        background:#1976D2;color:#fff;border:none;border-radius:4px;
        padding:3px 10px;font-size:11px;font-weight:600;cursor:pointer;
        line-height:1.6;
      `
      const editBtn = `
        background:#fff;color:#1976D2;border:1px solid #1976D2;border-radius:4px;
        padding:3px 10px;font-size:11px;font-weight:600;cursor:pointer;
        line-height:1.6;margin-left:4px;
      `
      const submitBtn = `
        background:#2E7D32;color:#fff;border:none;border-radius:4px;
        padding:3px 10px;font-size:11px;font-weight:600;cursor:pointer;
        line-height:1.6;margin-left:4px;
      `
      const showEdit = isEditableClaimStatus(params.data?.status)
      const showSubmit = isSubmittableClaimStatus(params.data?.status)
      return (
        `<button data-action="view" style="${viewBtn}">View</button>` +
        (showEdit
          ? `<button data-action="edit" style="${editBtn}">Edit</button>`
          : '') +
        (showSubmit
          ? `<button data-action="submit" style="${submitBtn}">Submit</button>`
          : '')
      )
    },
    onCellClicked: (params: any) => {
      const action = (
        params.event?.target as HTMLElement | undefined
      )?.getAttribute?.('data-action')
      if (action === 'edit') {
        editClaim(params)
      } else if (action === 'submit') {
        promptSubmitForAssessment(params)
      } else {
        viewClaimDetails(params)
      }
    }
  }
]

// Computed properties
const filteredClaims = computed(() => {
  let filtered = [...claims.value]

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(
      (claim) =>
        claim.claim_number.toLowerCase().includes(query) ||
        claim.member_name.toLowerCase().includes(query) ||
        claim.member_id_number.includes(query) ||
        claim.scheme_name.toLowerCase().includes(query)
    )
  }

  if (selectedStatus.value) {
    filtered = filtered.filter((claim) => claim.status === selectedStatus.value)
  }

  if (selectedBenefitType.value) {
    filtered = filtered.filter(
      (claim) => claim.benefit_alias === selectedBenefitType.value.benefit_alias
    )
  }

  if (selectedScheme.value) {
    filtered = filtered.filter(
      (claim) =>
        claim.scheme_name ===
        schemes.value.find((s) => s.id === selectedScheme.value)?.name
    )
  }

  return filtered
})

const claimsStats = computed((): ClaimsStats => {
  const currentMonth = new Date().getMonth()
  const currentYear = new Date().getFullYear()
  return {
    drafts: claims.value.filter((c) => c.status === 'draft').length,
    pendingAssessment: claims.value.filter(
      (c) => c.status === 'pending_assessment'
    ).length,
    underAssessment: claims.value.filter((c) => c.status === 'under_assessment')
      .length,
    approved: claims.value.filter((c) => {
      const claimDate = new Date(c.date_notified)
      return (
        c.status === 'approved' &&
        claimDate.getMonth() === currentMonth &&
        claimDate.getFullYear() === currentYear
      )
    }).length,
    declined: claims.value.filter((c) => {
      const claimDate = new Date(c.date_notified)
      return (
        c.status === 'declined' &&
        claimDate.getMonth() === currentMonth &&
        claimDate.getFullYear() === currentYear
      )
    }).length
  }
})

// Methods
const loadClaims = async () => {
  loading.value = true
  try {
    const response = await GroupPricingService.getClaims()
    claims.value = response.data || []
    const drafts = claims.value.filter((c: any) => c.status === 'draft').length
    const pendingAssessment = claims.value.filter(
      (c: any) => c.status === 'pending_assessment'
    ).length
    const underAssessment = claims.value.filter(
      (c: any) => c.status === 'under_assessment'
    ).length
    statusBarStore.set([
      {
        icon: 'mdi-file-edit-outline',
        text: `Drafts: ${drafts}`,
        severity: 'info'
      },
      {
        icon: 'mdi-clock-outline',
        text: `Pending assessment: ${pendingAssessment}`,
        severity: pendingAssessment > 0 ? 'warn' : 'info'
      },
      { icon: 'mdi-magnify', text: `Under assessment: ${underAssessment}` }
    ])
  } catch (error) {
    console.error('Error loading claims:', error)
    showSnackbar('Error loading claims', 'error')
    claims.value = []
  } finally {
    loading.value = false
  }
}

const loadSchemes = async () => {
  try {
    const response = await GroupPricingService.getSchemesInforce()
    schemes.value = response.data || []
  } catch (error) {
    console.error('Error loading schemes:', error)
    showSnackbar('Error loading schemes', 'error')
    // Fallback to empty array if API fails
    schemes.value = []
  }
}

const viewClaimDetails = (claim: any) => {
  const claimRow = claim.data
  if (!claimRow?.id) return
  router.push({
    name: 'group-pricing-claim-details',
    params: { id: claimRow.id }
  })
}

const editClaim = (claim: any) => {
  const claimRow = claim.data
  if (!claimRow?.id) return
  router.push({
    name: 'group-pricing-claim-edit',
    params: { id: claimRow.id }
  })
}

const promptSubmitForAssessment = (claim: any) => {
  const claimRow = claim.data
  if (!claimRow?.id) return
  confirmTitle.value = 'Submit for Assessment'
  confirmMessage.value = `Send claim ${claimRow.claim_number || claimRow.id} to the assessment queue? It will no longer appear in Drafts.`
  confirmCallback.value = () => submitClaimForAssessment(claimRow.id)
  confirmDialog.value = true
}

const submitClaimForAssessment = async (claimId: number) => {
  try {
    await GroupPricingService.submitClaimForAssessment(claimId)
    showSnackbar('Claim submitted for assessment', 'success')
    await loadClaims()
  } catch (error: any) {
    const status = error?.response?.status
    if (status === 422) {
      const missing: string[] = error?.response?.data?.missing || []
      showSnackbar(
        missing.length
          ? `Cannot submit: missing ${missing.join(', ')}.`
          : 'Claim is incomplete and cannot be submitted.',
        'error'
      )
    } else if (status === 409) {
      showSnackbar(
        'This claim cannot be submitted in its current status.',
        'error'
      )
    } else {
      console.error('Error submitting claim for assessment:', error)
      showSnackbar('Error submitting claim. Please try again.', 'error')
    }
  }
}

const handleNewClaimSave = async (claimData: any) => {
  loading.value = true
  try {
    // Show different message if files are being uploaded
    if (
      claimData.supporting_documents &&
      claimData.supporting_documents.length > 0
    ) {
      showSnackbar(
        `Uploading claim with ${claimData.supporting_documents.length} document(s)...`,
        'info'
      )
    }
    await GroupPricingService.submitClaim(claimData)
    showSnackbar('Claim registered successfully', 'success')
    newClaimDialog.value = false
    await loadClaims()
  } catch (error) {
    console.error('Error saving claim:', error)
    showSnackbar('Error registering claim. Please try again.', 'error')
  } finally {
    loading.value = false
  }
}

const handleBulkUploadComplete = async (result: any) => {
  const { successful, failed, errors } = result

  if (successful > 0) {
    showSnackbar(
      `Bulk upload completed! ${successful} claims uploaded successfully${
        failed > 0 ? `, ${failed} failed` : ''
      }.`,
      failed > 0 ? 'warning' : 'success'
    )
  } else {
    showSnackbar(
      'Bulk upload failed. Please check your file and try again.',
      'error'
    )
  }

  // Log errors for debugging if any
  if (errors && errors.length > 0) {
    console.error('Bulk upload errors:', errors)
  }

  bulkUploadDialog.value = false
  await loadClaims() // Refresh claims list
}

const showSnackbar = (message: string, color: string = 'success') => {
  snackbarMessage.value = message
  snackbarColor.value = color
  snackbar.value = true
}

const confirmAction = () => {
  if (confirmCallback.value) {
    confirmCallback.value()
    confirmCallback.value = null
  }
  confirmDialog.value = false
}

// Lifecycle
onMounted(async () => {
  const res = await GroupPricingService.getBenefitMaps()
  benefitMaps.value = res.data

  loadClaims()
  loadSchemes()
})
onUnmounted(() => statusBarStore.clear())
</script>

<style scoped>
.cursor-pointer {
  cursor: pointer;
}

.v-card {
  transition: transform 0.2s;
}

.v-card:hover {
  transform: translateY(-2px);
}
</style>
