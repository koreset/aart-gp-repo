<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex justify-space-between align-center">
              <span class="headline">Member Management</span>
              <div class="d-flex gap-2">
                <v-btn
                  color="white"
                  class="mr-2"
                  rounded
                  size="small"
                  variant="outlined"
                  prepend-icon="mdi-file-upload"
                  @click="bulkEnrollmentDialog = true"
                >
                  Bulk Enrollment
                </v-btn>
                <v-btn
                  size="small"
                  color="white"
                  variant="outlined"
                  rounded
                  prepend-icon="mdi-account-plus"
                  @click="addMemberDialog = true"
                >
                  Add Member
                </v-btn>
              </div>
            </div>
          </template>
          <template #default>
            <!-- Empty state when no members loaded -->
            <v-row
              v-if="!loading && members.length === 0 && !hasSearched"
              class="mb-4"
            >
              <v-col cols="12" class="text-center py-8">
                <v-icon size="64" color="grey-lighten-1" class="mb-4"
                  >mdi-account-search</v-icon
                >
                <h3 class="text-h6 text-grey-darken-1 mb-2"
                  >No Members Loaded</h3
                >
                <p class="text-body-1 text-grey mb-4">
                  Use the search and filters above, then click "Search Members"
                  to load member data.
                  <br />
                  <small
                    >This prevents loading potentially large datasets
                    automatically.</small
                  >
                </p>
              </v-col>
            </v-row>

            <!-- Loading Progress Bar -->
            <v-row v-if="loading" class="mb-4">
              <v-col cols="12">
                <v-card variant="outlined" class="pa-4">
                  <div class="d-flex justify-space-between align-center mb-2">
                    <span>{{ loadingMessage }}</span>
                    <v-btn
                      size="small"
                      color="error"
                      variant="outlined"
                      @click="cancelLoading"
                    >
                      Cancel
                    </v-btn>
                  </div>
                  <v-progress-linear
                    :model-value="loadingProgress"
                    color="primary"
                    height="8"
                    rounded
                  />
                </v-card>
              </v-col>
            </v-row>

            <!-- Search and Filter Bar -->
            <v-row class="mb-4">
              <v-col cols="12" md="3">
                <v-text-field
                  v-model="searchQuery"
                  label="Search Members"
                  prepend-inner-icon="mdi-magnify"
                  variant="outlined"
                  density="compact"
                  clearable
                  @input="searchMembers"
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
                  @update:model-value="filterByScheme"
                />
              </v-col>
              <v-col cols="12" md="3">
                <v-select
                  v-model="selectedStatus"
                  :items="memberStatuses"
                  label="Filter by Status"
                  variant="outlined"
                  density="compact"
                  clearable
                  @update:model-value="filterByStatus"
                />
              </v-col>
              <v-col cols="12" md="3" class="d-flex gap-2">
                <v-btn
                  class="mr-2 mt-1"
                  rounded
                  size="small"
                  color="primary"
                  variant="outlined"
                  :loading="loading"
                  @click="reloadMembers"
                >
                  Search Members
                </v-btn>
                <v-btn
                  class="mt-1"
                  rounded
                  size="small"
                  color="info"
                  variant="outlined"
                  @click="exportMembers"
                >
                  Export
                </v-btn>
              </v-col>
            </v-row>

            <!-- Members Data Grid -->
            <v-row>
              <v-col>
                <div :style="{ height: gridHeight, width: '100%' }">
                  <data-grid
                    :columnDefs="memberColumnDefs"
                    :rowData="filteredMembers"
                    :pagination="false"
                    :loading="loading"
                    style="height: 100%; width: 100%"
                    @row-double-clicked="handleRowClick"
                  />
                </div>
              </v-col>
            </v-row>

            <!-- Pagination Controls -->
            <v-row class="mt-4">
              <v-col cols="12" md="6">
                <v-card variant="outlined" class="pa-3">
                  <div class="text-body-2 text-medium-emphasis">
                    Showing {{ paginationInfo.displayedMembers }} of
                    {{ paginationInfo.totalMembers }} members
                    <span v-if="paginationInfo.hasMore"
                      >({{
                        paginationInfo.totalMembers -
                        paginationInfo.displayedMembers
                      }}
                      more available)</span
                    >
                  </div>
                </v-card>
              </v-col>
              <v-col cols="12" md="6" class="d-flex justify-end">
                <v-btn
                  v-if="paginationInfo.hasMore"
                  rounded
                  size="small"
                  :loading="loading"
                  :disabled="loading"
                  color="primary"
                  variant="outlined"
                  @click="loadMoreMembers"
                >
                  Load More Members
                </v-btn>
                <v-btn
                  v-if="members.length > 0"
                  rounded
                  size="small"
                  :loading="loading"
                  :disabled="loading"
                  color="info"
                  variant="outlined"
                  class="ml-2"
                  @click="reloadMembers"
                >
                  Refresh
                </v-btn>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Add/Edit Member Dialog -->
    <v-dialog v-model="addMemberDialog" persistent max-width="800px">
      <base-card>
        <template #header>
          <span class="headline">{{
            isEditMode ? 'Edit Member' : 'Add New Member'
          }}</span>
        </template>
        <template #default>
          <member-enrollment-form
            :member="selectedMember"
            :schemes="schemes"
            :is-edit-mode="isEditMode"
            :preselected-scheme-id="selectedScheme"
            @save="handleMemberSave"
            @cancel="closeAddMemberDialog"
          />
        </template>
      </base-card>
    </v-dialog>

    <!-- Bulk Enrollment Dialog -->
    <v-dialog v-model="bulkEnrollmentDialog" persistent max-width="800px">
      <base-card>
        <template #header>
          <span class="headline">Bulk Member Enrollment</span>
        </template>
        <template #default>
          <bulk-member-enrollment
            :schemes="schemes"
            @upload-complete="handleBulkUploadComplete"
            @cancel="bulkEnrollmentDialog = false"
          />
        </template>
      </base-card>
    </v-dialog>

    <!-- Member Details Dialog -->
    <v-dialog v-model="memberDetailsDialog" persistent max-width="1200px">
      <base-card>
        <template #header>
          <div class="d-flex justify-space-between align-center">
            <span class="headline"
              >Member Details - {{ selectedMember?.member_name }}</span
            >
            <div>
              <v-btn
                size="small"
                rounded
                color="white"
                variant="outlined"
                class="mr-2"
                @click="editMember"
              >
                Edit Member
              </v-btn>
              <v-btn
                rounded
                size="small"
                color="error"
                variant="outlined"
                @click="deactivateMember"
              >
                Deactivate
              </v-btn>
            </div>
          </div>
        </template>
        <template #default>
          <member-detail-view
            :member="selectedMember"
            :beneficiaries="memberBeneficiaries"
            @manage-beneficiaries="openBeneficiaryManagement"
            @view-claims="viewMemberClaims"
            @member-updated="handleMemberUpdated"
          />
        </template>
        <template #actions>
          <v-btn color="grey" @click="memberDetailsDialog = false">Close</v-btn>
        </template>
      </base-card>
    </v-dialog>

    <!-- Exit Date Confirmation Dialog -->
    <v-dialog v-model="exitDateDialog" persistent max-width="500px">
      <base-card>
        <template #header>
          <span class="headline">Confirm Member Deactivation</span>
        </template>
        <template #default>
          <div class="pa-4">
            <div class="mb-4">
              <p class="text-body-1 mb-2">
                You are about to deactivate
                <strong>{{ selectedMember?.member_name }}</strong> from the
                <strong>{{ selectedMember?.scheme_name }}</strong> scheme.
              </p>
              <p class="text-body-2 text-medium-emphasis">
                Please specify the effective exit date for this member.
              </p>
            </div>

            <v-form
              v-model="exitDateValid"
              @submit.prevent="confirmDeactivation"
            >
              <v-text-field
                v-model="exitDate"
                label="Effective Exit Date"
                type="date"
                variant="outlined"
                :rules="exitDateRules"
                :min="
                  selectedMember?.entry_date
                    ? new Date(selectedMember.entry_date)
                        .toISOString()
                        .substr(0, 10)
                    : undefined
                "
                required
                class="mb-4"
              />

              <div class="text-caption text-medium-emphasis mb-4">
                <v-icon size="small" class="mr-1">mdi-information</v-icon>
                Entry Date:
                {{
                  selectedMember?.entry_date
                    ? new Date(selectedMember.entry_date).toLocaleDateString()
                    : 'Not specified'
                }}
              </div>
            </v-form>
          </div>
        </template>
        <template #actions>
          <v-btn color="grey" @click="cancelDeactivation">Cancel</v-btn>
          <v-btn
            color="error"
            :disabled="!exitDateValid"
            @click="confirmDeactivation"
          >
            Deactivate Member
          </v-btn>
        </template>
      </base-card>
    </v-dialog>

    <!-- Snackbar for notifications -->
    <v-snackbar v-model="snackbar" :timeout="4000" :color="snackbarColor">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import MemberEnrollmentForm from './components/MemberEnrollmentForm.vue'
import BulkMemberEnrollment from './components/BulkMemberEnrollment.vue'
import MemberDetailView from './components/MemberDetailView.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import formatValues from '@/renderer/utils/format_values'
import { statusCellRenderer } from '@/renderer/utils/statusCellRenderer'
import { useGridHeight } from '@/renderer/composables/useGridHeight'
import { useStatusBarStore } from '@/renderer/store/statusBar'

const gridHeight = useGridHeight(340)
const statusBarStore = useStatusBarStore()

const router = useRouter()
const route = useRoute()

// Interfaces
// interface Member {
//   id?: number
//   member_name: string
//   member_id_number: string
//   scheme_name: string
//   scheme_id: number
//   status?: string
//   annual_salary: number
//   entry_date: string
//   gender: string
//   effective_exit_date?: string | null
// }

interface Member {
  id?: number
  member_name: string
  member_id_number: string
  member_id_type: string
  scheme_name: string
  gender: string
  date_of_birth: Date | null
  email?: string
  phone_number?: string
  employee_number?: string
  scheme_id: number | null
  scheme_category: string
  entry_date: Date | null
  annual_salary: number
  status?: string
  effective_exit_date?: string | Date | null
  occupation?: string
  occupational_class?: string
  address_line_1?: string
  address_line_2?: string
  city?: string
  province?: string
  postal_code?: string
  benefits: {
    gla_enabled: boolean
    gla_multiple?: number
    sgla_enabled: boolean
    sgla_multiple?: number
    ptd_enabled: boolean
    ptd_multiple?: number
    ci_enabled: boolean
    ci_multiple?: number
    ttd_enabled: boolean
    ttd_multiple?: number
    phi_enabled: boolean
    phi_multiple?: number
    gff_enabled: boolean
  }
}

interface Scheme {
  id: number
  name: string
}

// State
const loading = ref(false)
const loadingProgress = ref(0)
const loadingMessage = ref('')
const searchQuery = ref('')
const selectedScheme = ref<number | null>(null)
const selectedStatus = ref<string | null>(null)
const members = ref<Member[]>([])
const schemes = ref<Scheme[]>([])
const selectedMember = ref<Member | null>(null)
const memberBeneficiaries = ref<any[]>([])

// Pagination and performance state
const totalMembers = ref(0)
const hasMoreMembers = ref(true)
const hasSearched = ref(false)
const loadAbortController = ref<AbortController | null>(null)

// Server-side filters
const serverFilters = ref({
  search: '',
  schemeId: null as number | null,
  status: null as string | null,
  page: 1,
  pageSize: 100
})

// Dialog states
const addMemberDialog = ref(false)
const bulkEnrollmentDialog = ref(false)
const memberDetailsDialog = ref(false)
const exitDateDialog = ref(false)
const isEditMode = ref(false)

// Exit date form state
const exitDate = ref('')
const exitDateValid = ref(false)
const exitDateRules = [
  (v: string) => !!v || 'Exit date is required',
  (v: string) => {
    if (!v) return true
    const selectedDate = new Date(v)
    const today = new Date()
    const thirtyDaysAgo = new Date(today.getTime() - 30 * 24 * 60 * 60 * 1000)
    const entryDate = selectedMember.value?.entry_date
      ? new Date(selectedMember.value.entry_date)
      : null

    if (selectedDate < thirtyDaysAgo) {
      return 'Exit date cannot be more than 30 days in the past'
    }
    if (entryDate && selectedDate <= entryDate) {
      return 'Exit date must be after entry date'
    }
    return true
  }
]

// Snackbar
const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

// Filter options
const memberStatuses = [
  { title: 'ACTIVE', value: 'ACTIVE' },
  { title: 'INACTIVE', value: 'INACTIVE' },
  { title: 'PENDING', value: 'PENDING' },
  { title: 'SUSPENDED', value: 'SUSPENDED' }
]

// Column definitions for the data grid
const memberColumnDefs = [
  {
    headerName: 'Member Name',
    field: 'member_name',
    filter: true,
    sortable: true,
    minWidth: 200
  },
  {
    headerName: 'Employee Number',
    field: 'employee_number',
    filter: true,
    sortable: true,
    minWidth: 150
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
    minWidth: 180
  },
  {
    headerName: 'Scheme Category',
    field: 'scheme_category',
    filter: true,
    sortable: true,
    minWidth: 180
  },

  {
    headerName: 'Status',
    field: 'status',
    filter: true,
    sortable: true,
    minWidth: 220,
    cellRenderer: (params: any) => {
      const pill = statusCellRenderer(params.value || 'active')
      const exit = params.data?.effective_exit_date
      if (!exit) return pill
      const d = new Date(exit)
      if (isNaN(d.getTime())) return pill
      const today = new Date()
      today.setHours(0, 0, 0, 0)
      const dDay = new Date(d)
      dDay.setHours(0, 0, 0, 0)
      const verb = dDay > today ? 'Exits' : 'Exited'
      return `
        <div style="display:flex; align-items:center; gap:6px;">
          ${pill}
          <span style="font-size:10px; color:#757575; white-space:nowrap;">${verb} ${d.toLocaleDateString()}</span>
        </div>
      `
    }
  },
  {
    headerName: 'Annual Salary',
    field: 'annual_salary',
    filter: true,
    sortable: true,
    valueFormatter: formatValues,
    minWidth: 150
  },
  {
    headerName: 'Entry Date',
    field: 'entry_date',
    filter: true,
    sortable: true,
    minWidth: 120,
    valueFormatter: (params: any) => {
      return params.value ? new Date(params.value).toLocaleDateString() : ''
    }
  },
  {
    headerName: 'Gender',
    field: 'gender',
    filter: true,
    sortable: true,
    minWidth: 100
  },
  {
    headerName: 'Actions',
    field: 'actions',
    sortable: false,
    filter: false,
    minWidth: 120,
    cellRenderer: () => {
      return `
        <v-btn size="small" color="primary" variant="text">
          View Details
        </v-btn>
      `
    }
  }
]

// Computed properties
const filteredMembers = computed((): Member[] => {
  // With server-side filtering, we just return the loaded members
  // All filtering is now done on the server
  return members.value
})

// Computed property for pagination info
const paginationInfo = computed(() => ({
  currentPage: Math.ceil(members.value.length / serverFilters.value.pageSize),
  totalPages: Math.ceil(totalMembers.value / serverFilters.value.pageSize),
  totalMembers: totalMembers.value,
  displayedMembers: members.value.length,
  hasMore: hasMoreMembers.value
}))

// Methods
const loadMembers = async (append: boolean = false) => {
  // Cancel any existing load operation
  if (loadAbortController.value) {
    loadAbortController.value.abort()
  }

  loadAbortController.value = new AbortController()
  loading.value = true
  loadingProgress.value = 0
  loadingMessage.value = 'Loading schemes...'

  try {
    // Load schemes first (lightweight operation)
    if (schemes.value.length === 0) {
      const schemesResponse = await GroupPricingService.getSchemesInforce()
      schemes.value = schemesResponse.data
      loadingProgress.value = 10
    }

    // Use server-side pagination and filtering (with fallback)
    loadingMessage.value = `Loading members (page ${serverFilters.value.page})...`

    let membersResponse
    let newMembers, total, hasMore

    try {
      // Try the new paginated API first
      membersResponse = await GroupPricingService.getMembersPaginated({
        page: serverFilters.value.page,
        pageSize: serverFilters.value.pageSize,
        search: serverFilters.value.search,
        schemeId: serverFilters.value.schemeId,
        status: serverFilters.value.status,
        signal: loadAbortController.value.signal
      })

      const response = membersResponse.data
      newMembers = response.data
      total = response.total
      hasMore = response.hasMore
    } catch (paginationError) {
      // Fallback to old method with client-side filtering (for backward compatibility)
      console.warn(
        'Paginated API not available, falling back to client-side filtering'
      )
      loadingMessage.value = 'Loading all members (fallback mode)...'

      const allMembers: Member[] = []
      let processedSchemes = 0

      const schemesToLoad = serverFilters.value.schemeId
        ? schemes.value.filter((s) => s.id === serverFilters.value.schemeId)
        : schemes.value

      for (const scheme of schemesToLoad) {
        try {
          const membersResponse = await GroupPricingService.getMembersInForce(
            scheme.id
          )
          const schemeMembers: Member[] = membersResponse.data.map(
            (member: any) => ({
              ...member,
              scheme_name: scheme.name,
              scheme_id: scheme.id,
              status: member.status || 'active'
            })
          )

          // Apply client-side filtering
          let filteredSchemeMembers = schemeMembers

          if (serverFilters.value.search) {
            const query = serverFilters.value.search.toLowerCase()
            filteredSchemeMembers = filteredSchemeMembers.filter(
              (member) =>
                member.member_name?.toLowerCase().includes(query) ||
                member.member_id_number?.toLowerCase().includes(query)
            )
          }

          if (serverFilters.value.status) {
            filteredSchemeMembers = filteredSchemeMembers.filter(
              (member) => member.status === serverFilters.value.status
            )
          }

          allMembers.push(...filteredSchemeMembers)
          processedSchemes++

          // Update progress
          loadingProgress.value = Math.min(
            90,
            (processedSchemes / schemesToLoad.length) * 80 + 10
          )
          loadingMessage.value = `Processed ${processedSchemes}/${schemesToLoad.length} schemes...`
        } catch (schemeError) {
          console.error(
            `Error loading members for scheme ${scheme.id}:`,
            schemeError
          )
          // Continue with other schemes
        }
      }

      // Apply pagination to client-filtered results
      const startIndex =
        (serverFilters.value.page - 1) * serverFilters.value.pageSize
      const endIndex = startIndex + serverFilters.value.pageSize

      newMembers = allMembers.slice(startIndex, endIndex)
      total = allMembers.length
      hasMore = endIndex < allMembers.length
    }

    // Map scheme names to members
    const membersWithSchemeNames = newMembers.map((member: any) => {
      const scheme = schemes.value.find((s) => s.id === member.scheme_id)
      return {
        ...member,
        scheme_name: scheme?.name || 'Unknown Scheme'
      }
    })

    if (append) {
      members.value = [...members.value, ...membersWithSchemeNames]
    } else {
      members.value = membersWithSchemeNames
    }

    totalMembers.value = total
    hasMoreMembers.value = hasMore
    loadingProgress.value = 100
    loadingMessage.value = 'Loading complete'
  } catch (error: any) {
    if (error.name === 'AbortError') {
      loadingMessage.value = 'Loading cancelled'
      return
    }
    console.error('Error loading members:', error)
    showSnackbar('Error loading members', 'error')
  } finally {
    loading.value = false
    loadAbortController.value = null
  }
}

// Load more members (for infinite scroll or load more button)
const loadMoreMembers = async () => {
  if (!hasMoreMembers.value || loading.value) return

  serverFilters.value.page += 1
  await loadMembers(true)
}

// Reset and reload with new filters
const reloadMembers = async () => {
  serverFilters.value.page = 1
  hasSearched.value = true
  await loadMembers(false)
}

// Cancel current loading operation
const cancelLoading = () => {
  if (loadAbortController.value) {
    loadAbortController.value.abort()
    loading.value = false
    loadingMessage.value = 'Loading cancelled'
  }
}

const handleRowClick = (event: any) => {
  console.log('Row clicked:', event.data)
  selectedMember.value = event.data
  loadMemberBeneficiaries(event.data.id)
  memberDetailsDialog.value = true
}

const loadMemberBeneficiaries = async (memberId: number) => {
  try {
    // This would be implemented when beneficiary API is available
    if (memberId == null) return

    const response = await GroupPricingService.getMemberBeneficiaries(memberId)
    memberBeneficiaries.value = response.data
  } catch (error) {
    console.error('Error loading beneficiaries:', error)
  }

  // loading.value = true
  //   try {
  //     if (memberId.value) {
  //       // Load beneficiaries for specific member
  //       // This would be implemented when the API is available
  //       const response = await GroupPricingService.getMemberBeneficiaries(memberId.value)
  //       beneficiaries.value = response.data
  //     }
  //   } catch (error) {
  //     console.error('Error loading beneficiaries:', error)
  //     showSnackbar('Error loading beneficiaries', 'error')
  //   } finally {
  //     loading.value = false
  //   }
}

const handleMemberSave = async (memberData: any) => {
  console.log('Saving member data:', memberData)
  try {
    if (isEditMode.value) {
      await GroupPricingService.updateMember(
        selectedMember.value?.id,
        memberData
      )
      showSnackbar('Member updated successfully', 'success')
    } else {
      await GroupPricingService.addMember(memberData)
      showSnackbar('Member added successfully', 'success')
    }
    await loadMembers()
    closeAddMemberDialog()
  } catch (error: any) {
    console.error('Error saving member:', error)
    showSnackbar(
      error?.response?.data || error?.message || 'Failed to save member',
      'error'
    )
  }
}

const closeAddMemberDialog = () => {
  addMemberDialog.value = false
  isEditMode.value = false
  selectedMember.value = null
}

const editMember = () => {
  isEditMode.value = true
  memberDetailsDialog.value = false
  addMemberDialog.value = true
}

const deactivateMember = () => {
  if (!selectedMember.value) return

  // Initialize exit date to today as default
  exitDate.value = new Date().toISOString().substr(0, 10)
  exitDateValid.value = false
  exitDateDialog.value = true
}

const confirmDeactivation = async () => {
  if (!selectedMember.value || !exitDateValid.value) return

  try {
    const today = new Date().toISOString().slice(0, 10)
    const exitInFuture = exitDate.value > today

    const memberToUpdate: any = {
      ...selectedMember.value,
      effective_exit_date: exitDate.value,
      status: exitInFuture ? 'ACTIVE' : 'INACTIVE'
    }
    await GroupPricingService.removeMemberFromScheme(
      selectedMember.value.scheme_id,
      memberToUpdate
    )
    showSnackbar(
      exitInFuture
        ? `Deactivation scheduled for ${exitDate.value}`
        : 'Member deactivated successfully',
      'success'
    )
    await loadMembers()
    memberDetailsDialog.value = false
    exitDateDialog.value = false
  } catch (error) {
    console.error('Error deactivating member:', error)
    showSnackbar('Error deactivating member', 'error')
  }
}

const cancelDeactivation = () => {
  exitDateDialog.value = false
  exitDate.value = ''
  exitDateValid.value = false
}

const handleBulkUploadComplete = (result: any) => {
  showSnackbar(
    `Bulk upload complete. Success: ${result.success}, Failed: ${result.failed}`,
    result.failed > 0 ? 'warning' : 'success'
  )

  if (result.failed > 0) {
    setTimeout(() => {
      result.errors.forEach((error) => {
        showSnackbar(error.message, 'error')
      })
    }, 5000)
  }

  bulkEnrollmentDialog.value = false
  loadMembers()
}

// Debounced search to avoid too many API calls
let searchTimeout: ReturnType<typeof setTimeout> | null = null
const searchMembers = () => {
  if (searchTimeout) clearTimeout(searchTimeout)

  searchTimeout = setTimeout(async () => {
    serverFilters.value.search = searchQuery.value
    await reloadMembers()
  }, 500) // 500ms delay
}

const filterByScheme = async () => {
  serverFilters.value.schemeId = selectedScheme.value
  await reloadMembers()
}

const filterByStatus = async () => {
  serverFilters.value.status = selectedStatus.value
  await reloadMembers()
}

const exportMembers = () => {
  // Implementation for exporting members to CSV/Excel
  const csvData = filteredMembers.value.map((member) => ({
    'Member Name': member.member_name,
    'ID Number': member.member_id_number,
    Scheme: member.scheme_name,
    Status: member.status,
    'Annual Salary': member.annual_salary,
    'Entry Date': member.entry_date,
    Gender: member.gender
  }))

  const csvRows = [
    Object.keys(csvData[0]).join(','),
    ...csvData.map((row) => Object.values(row).join(','))
  ].join('\n')

  const blob = new Blob(['\uFEFF' + csvRows], {
    type: 'text/csv;charset=utf-8;'
  })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.setAttribute('href', url)
  link.setAttribute(
    'download',
    `members_export_${new Date().toISOString().split('T')[0]}.csv`
  )
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

const openBeneficiaryManagement = () => {
  if (!selectedMember.value?.id) return
  router.push({
    name: 'group-pricing-beneficiary-management',
    params: { memberId: selectedMember.value.id }
  })
}

const viewMemberClaims = () => {
  if (!selectedMember.value?.id) return
  router.push({
    name: 'group-pricing-claims-management',
    query: { memberId: selectedMember.value.id }
  })
}

const handleMemberUpdated = async (updatedMember: any) => {
  // Update the selected member with new data
  selectedMember.value = updatedMember

  // Refresh the member list to show updated information
  await loadMembers()

  // Show success message
  showSnackbar('Member updated successfully', 'success')
}

const showSnackbar = (message: string, color: string = 'success') => {
  snackbarMessage.value = message
  snackbarColor.value = color
  snackbar.value = true
}

// Handle auto-opening member details when returning from beneficiary management
const handleAutoOpenMemberDetails = async () => {
  const memberIdToOpen = route.query.openMemberDetails
  if (memberIdToOpen) {
    const memberId = parseInt(memberIdToOpen as string, 10)
    // Find the member in the loaded members or load it specifically
    let member = members.value.find((m) => m.id === memberId)

    if (!member && hasSearched.value) {
      // If member not found in current results, try to load it specifically
      try {
        const memberResponse = await GroupPricingService.getMemberInfo(memberId)
        member = memberResponse.data
        // Add scheme name if missing
        if (member && !member.scheme_name) {
          const scheme = schemes.value.find((s) => s.id === member?.scheme_id)
          member.scheme_name = scheme?.name || 'Unknown Scheme'
        }
      } catch (error) {
        console.error('Error loading specific member:', error)
        showSnackbar('Member not found', 'error')
        return
      }
    }

    if (member && member.id) {
      selectedMember.value = member
      await loadMemberBeneficiaries(member.id)
      memberDetailsDialog.value = true

      // Clear the query parameter to prevent reopening on refresh
      router.replace({
        name: route.name,
        query: { ...route.query, openMemberDetails: undefined }
      })
    }
  }
}

// Watch member count and update status bar whenever data changes
watch([totalMembers, schemes], () => {
  statusBarStore.set([
    {
      icon: 'mdi-office-building-outline',
      text: `Schemes: ${schemes.value.length}`
    },
    ...(totalMembers.value > 0
      ? [
          {
            icon: 'mdi-account-group',
            text: `Members: ${totalMembers.value.toLocaleString()}`
          }
        ]
      : [])
  ])
})

// Lifecycle
onMounted(async () => {
  // Load schemes first for the filter dropdown
  try {
    const schemesResponse = await GroupPricingService.getSchemesInforce()
    schemes.value = schemesResponse.data
    console.log('loaded schemes:', schemes.value)

    // Check if we came from GroupSchemeDetail with a schemeId
    const schemeId = route.query.schemeId
    if (schemeId) {
      // Pre-select the scheme
      selectedScheme.value = parseInt(schemeId as string, 10)
      // Filter by the pre-selected scheme
      await filterByScheme()
      // Auto-open the Add Member dialog
      setTimeout(() => {
        addMemberDialog.value = true
      }, 500) // Small delay to ensure data is loaded
    }

    // Check if we need to auto-open member details (returning from beneficiary management)
    const openMemberDetails = route.query.openMemberDetails
    if (openMemberDetails) {
      // Load members first if not already loaded
      if (!hasSearched.value) {
        await reloadMembers()
      }
      setTimeout(() => {
        handleAutoOpenMemberDetails()
      }, 500)
    }
  } catch (error) {
    console.error('Error loading schemes:', error)
    showSnackbar('Error loading schemes', 'error')
  }

  // Don't automatically load members - let user search/filter first
  // This prevents loading potentially millions of records on component mount
})

onUnmounted(() => {
  if (loadAbortController.value) {
    loadAbortController.value.abort()
  }
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  statusBarStore.clear()
})
</script>

<style scoped></style>
