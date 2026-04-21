<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex justify-space-between align-center">
              <div class="d-flex align-center">
                <v-btn
                  v-if="memberId"
                  icon="mdi-arrow-left"
                  variant="text"
                  class="mr-2"
                  @click="goBackToMemberDetails"
                />
                <span class="headline">
                  {{
                    memberId ? 'Beneficiary Management' : 'Beneficiary Overview'
                  }}
                </span>
              </div>
              <div v-if="memberId">
                <v-btn color="primary" @click="addBeneficiaryDialog = true">
                  Add Beneficiary
                </v-btn>
              </div>
            </div>
          </template>
          <template #default>
            <!-- Member Information Summary (only show when viewing specific member) -->
            <v-row v-if="memberId" class="mb-4">
              <v-col>
                <v-card variant="outlined">
                  <v-card-title class="text-subtitle-1 bg-primary text-white">
                    Member Information
                  </v-card-title>
                  <v-card-text>
                    <v-row>
                      <v-col cols="12" md="3">
                        <strong>Name:</strong>
                        {{ memberInfo?.member_name || 'Loading...' }}
                      </v-col>
                      <v-col cols="12" md="3">
                        <strong>ID Number:</strong>
                        {{ memberInfo?.member_id_number || 'Loading...' }}
                      </v-col>
                      <v-col cols="12" md="3">
                        <strong>Scheme:</strong>
                        {{ memberInfo?.scheme_name || 'Loading...' }}
                      </v-col>
                      <v-col cols="12" md="3">
                        <strong>Total Allocation:</strong>
                        <v-chip
                          :color="
                            totalAllocation === 100 ? 'success' : 'warning'
                          "
                          size="small"
                        >
                          {{ totalAllocation }}%
                        </v-chip>
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Beneficiaries List -->
            <v-row>
              <v-col>
                <data-grid
                  :columnDefs="beneficiaryColumnDefs"
                  :rowData="beneficiaries"
                  :pagination="true"
                  :loading="loading"
                  @row-clicked="handleBeneficiaryClick"
                />
              </v-col>
            </v-row>

            <!-- Validation Messages -->
            <v-row v-if="validationMessages.length > 0" class="mt-4">
              <v-col>
                <v-alert type="warning" variant="tonal">
                  <div class="text-subtitle-2">Validation Issues:</div>
                  <ul class="mt-2">
                    <li v-for="message in validationMessages" :key="message">{{
                      message
                    }}</li>
                  </ul>
                </v-alert>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Add/Edit Beneficiary Dialog -->
    <v-dialog v-model="addBeneficiaryDialog" persistent max-width="800px">
      <base-card>
        <template #header>
          <span class="headline">
            {{ isEditMode ? 'Edit Beneficiary' : 'Add New Beneficiary' }}
          </span>
        </template>
        <template #default>
          <beneficiary-form
            :beneficiary="selectedBeneficiary"
            :member-id="memberId"
            :is-edit-mode="isEditMode"
            :existing-beneficiaries="beneficiaries"
            @save="handleBeneficiarySave"
            @cancel="closeBeneficiaryDialog"
          />
        </template>
      </base-card>
    </v-dialog>

    <!-- Beneficiary Details Dialog -->
    <v-dialog v-model="beneficiaryDetailsDialog" persistent max-width="600px">
      <base-card>
        <template #header>
          <div class="d-flex justify-space-between align-center">
            <span class="headline">Beneficiary Details</span>
            <div>
              <v-btn
                color="white"
                variant="outlined"
                class="mr-2"
                @click="editBeneficiary"
              >
                Edit
              </v-btn>
              <v-btn
                color="error"
                variant="outlined"
                @click="deleteBeneficiary"
              >
                Delete
              </v-btn>
            </div>
          </div>
        </template>
        <template #default>
          <beneficiary-detail-view :beneficiary="selectedBeneficiary" />
        </template>
        <template #actions>
          <v-btn color="grey" @click="beneficiaryDetailsDialog = false"
            >Close</v-btn
          >
        </template>
      </base-card>
    </v-dialog>

    <!-- Confirmation Dialog -->
    <v-dialog v-model="confirmDialog" max-width="400px">
      <v-card>
        <v-card-title>Confirm Action</v-card-title>
        <v-card-text>{{ confirmMessage }}</v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey" @click="confirmDialog = false">Cancel</v-btn>
          <v-btn color="error" @click="confirmAction">Confirm</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Snackbar -->
    <v-snackbar v-model="snackbar" :timeout="4000" :color="snackbarColor">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import BeneficiaryForm from './components/BeneficiaryForm.vue'
import BeneficiaryDetailView from './components/BeneficiaryDetailView.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

const route = useRoute()
const router = useRouter()
const memberId = computed(() => {
  const id = route.params.memberId
  return id ? parseInt(id as string) : null
})

// Interfaces
interface MemberInfo {
  member_name: string
  member_id_number: string
  scheme_name: string
}

interface Beneficiary {
  id?: number
  full_name: string
  relationship: string
  id_type: string
  id_number: string
  gender: string
  date_of_birth: Date | null
  contact_number?: string
  email?: string
  address?: string
  allocation_percentage: number
  benefit_types: string[]
  guardian_name?: string
  guardian_relationship?: string
  guardian_id_number?: string
  guardian_contact?: string
  bank_name?: string
  branch_code?: string
  account_number?: string
  account_type?: string
  status?: string
  memberId?: number
}

// interface Beneficiary {
//   id: number
//   full_name: string
//   relationship: string
//   id_number: string
//   date_of_birth: string
//   allocation_percentage: number
//   contact_number: string
//   status: string
// }

// State
const loading = ref(false)
const memberInfo = ref<MemberInfo | null>(null)
const beneficiaries = ref<Beneficiary[]>([])
const selectedBeneficiary = ref<Beneficiary | null>(null)

// Dialog states
const addBeneficiaryDialog = ref(false)
const beneficiaryDetailsDialog = ref(false)
const confirmDialog = ref(false)
const isEditMode = ref(false)

// Confirmation
const confirmMessage = ref('')
const confirmCallback = ref<(() => void) | null>(null)

// Snackbar
const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

// Column definitions for beneficiaries grid
const beneficiaryColumnDefs = [
  {
    headerName: 'Name',
    field: 'full_name',
    filter: true,
    sortable: true,
    minWidth: 200
  },
  {
    headerName: 'Relationship',
    field: 'relationship',
    filter: true,
    sortable: true,
    minWidth: 150
  },
  {
    headerName: 'ID Number',
    field: 'id_number',
    filter: true,
    sortable: true,
    minWidth: 150
  },
  {
    headerName: 'Date of Birth',
    field: 'date_of_birth',
    filter: true,
    sortable: true,
    minWidth: 120,
    valueFormatter: (params: any) => {
      return params.value ? new Date(params.value).toLocaleDateString() : ''
    }
  },
  {
    headerName: 'Allocation %',
    field: 'allocation_percentage',
    filter: true,
    sortable: true,
    minWidth: 120,
    valueFormatter: (params: any) => `${params.value}%`
  },
  {
    headerName: 'Contact Number',
    field: 'contact_number',
    filter: true,
    sortable: true,
    minWidth: 150
  },
  {
    headerName: 'Status',
    field: 'status',
    filter: true,
    sortable: true,
    minWidth: 100,
    cellRenderer: (params: any) => {
      const status = params.value || 'active'
      const color = status === 'active' ? 'success' : 'error'
      return `<v-chip size="small" color="${color}">${status.toUpperCase()}</v-chip>`
    }
  },
  {
    headerName: 'Actions',
    field: 'actions',
    sortable: false,
    filter: false,
    minWidth: 120,
    cellRenderer: () => {
      return `<v-btn size="small" color="primary" variant="text">View Details</v-btn>`
    }
  }
]

// Computed properties
const totalAllocation = computed(() => {
  return beneficiaries.value.reduce(
    (total: number, beneficiary: any) =>
      total + (beneficiary.allocation_percentage || 0),
    0
  )
})

const validationMessages = computed((): string[] => {
  const messages: string[] = []

  if (totalAllocation.value !== 100) {
    messages.push(
      `Total allocation is ${totalAllocation.value}%, should be 100%`
    )
  }

  const relationships = beneficiaries.value.map((b: any) => b.relationship)
  const duplicateRelationships = relationships.filter(
    (rel: string, index: number) => relationships.indexOf(rel) !== index
  )

  if (duplicateRelationships.length > 0) {
    messages.push(
      'Duplicate relationships found: ' + duplicateRelationships.join(', ')
    )
  }

  const minorBeneficiaries = beneficiaries.value.filter((b: any) => {
    if (!b.date_of_birth) return false
    const age =
      new Date().getFullYear() - new Date(b.date_of_birth).getFullYear()
    return age < 18
  })

  if (minorBeneficiaries.length > 0 && !hasGuardianAppointed()) {
    messages.push('Minor beneficiaries require guardian appointment')
  }

  return messages
})

// Methods
const loadMemberInfo = async () => {
  console.log('Loading member info for memberId:', memberId.value)
  if (!memberId.value) {
    memberInfo.value = null
    return
  }

  try {
    // This would be implemented when the API is available
    const response = await GroupPricingService.getMemberInfo(memberId.value)
    memberInfo.value = response.data
    console.log('Member info response:', memberInfo.value)

    // Placeholder data
    // memberInfo.value = {
    //   member_name: 'John Doe',
    //   member_id_number: '8001015009087',
    //   scheme_name: 'ABC Company Group Scheme'
    // }
  } catch (error) {
    console.error('Error loading member info:', error)
    showSnackbar('Error loading member information', 'error')
  }
}

const loadBeneficiaries = async () => {
  loading.value = true
  try {
    if (memberId.value) {
      // Load beneficiaries for specific member
      // This would be implemented when the API is available
      const response = await GroupPricingService.getMemberBeneficiaries(
        memberId.value
      )
      beneficiaries.value = response.data
    }
  } catch (error) {
    console.error('Error loading beneficiaries:', error)
    showSnackbar('Error loading beneficiaries', 'error')
  } finally {
    loading.value = false
  }
}

const handleBeneficiaryClick = (event: any) => {
  selectedBeneficiary.value = event.data
  beneficiaryDetailsDialog.value = true
}

const handleBeneficiarySave = async (beneficiaryData: any) => {
  console.log('Saving beneficiary data:', beneficiaryData)
  try {
    if (isEditMode.value) {
      // Update beneficiary
      await GroupPricingService.updateBeneficiary(
        memberId.value,
        beneficiaryData
      )
      showSnackbar('Beneficiary updated successfully', 'success')
    } else {
      // Add new beneficiary
      await GroupPricingService.addBeneficiaryToMember(
        memberId.value,
        beneficiaryData
      )
      showSnackbar('Beneficiary added successfully', 'success')
    }
    await loadBeneficiaries()
    closeBeneficiaryDialog()
  } catch (error) {
    console.error('Error saving beneficiary:', error)
    showSnackbar('Error saving beneficiary', 'error')
  }
}

const closeBeneficiaryDialog = () => {
  addBeneficiaryDialog.value = false
  isEditMode.value = false
  selectedBeneficiary.value = null
}

const editBeneficiary = () => {
  isEditMode.value = true
  beneficiaryDetailsDialog.value = false
  addBeneficiaryDialog.value = true
}

const deleteBeneficiary = () => {
  if (!selectedBeneficiary.value) return

  confirmMessage.value = `Are you sure you want to delete beneficiary "${selectedBeneficiary.value.full_name}"?`
  confirmCallback.value = confirmDeleteBeneficiary
  confirmDialog.value = true
}

const confirmDeleteBeneficiary = async () => {
  if (!selectedBeneficiary.value?.id) return

  try {
    await GroupPricingService.deleteBeneficiary(
      memberId.value,
      selectedBeneficiary.value.id
    )
    showSnackbar('Beneficiary deleted successfully', 'success')
    await loadBeneficiaries()
    beneficiaryDetailsDialog.value = false
    confirmDialog.value = false
  } catch (error) {
    console.error('Error deleting beneficiary:', error)
    showSnackbar('Error deleting beneficiary', 'error')
  }
}

const confirmAction = () => {
  if (confirmCallback.value) {
    confirmCallback.value()
  }
}

const hasGuardianAppointed = () => {
  // Check if any beneficiary has guardian information
  return beneficiaries.value.some(
    (b: any) => (b as any).guardian_name && (b as any).guardian_contact
  )
}

const goBackToMemberDetails = () => {
  // Navigate back to member management with the member ID as a query parameter
  // This will allow the MemberManagement view to auto-open the member details dialog
  router.push({
    name: 'group-pricing-member-management',
    query: {
      openMemberDetails: memberId.value?.toString()
    }
  })
}

const showSnackbar = (message: string, color: string = 'success') => {
  snackbarMessage.value = message
  snackbarColor.value = color
  snackbar.value = true
}

// Watchers
watch(
  () => route.params.memberId,
  () => {
    loadMemberInfo()
    loadBeneficiaries()
  }
)

// Lifecycle
onMounted(() => {
  loadMemberInfo()
  loadBeneficiaries()
})
</script>

<style scoped></style>
