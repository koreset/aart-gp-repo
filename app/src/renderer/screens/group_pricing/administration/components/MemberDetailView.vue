<template>
  <v-container>
    <v-row>
      <!-- Member Details Cards -->
      <v-col cols="12" md="4">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-primary text-white">
            Personal Information
          </v-card-title>
          <v-card-text>
            <v-list density="compact">
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Full Name</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  member?.member_name || 'N/A'
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >ID Number</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  member?.member_id_number || 'N/A'
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Gender</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  member?.gender || 'N/A'
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Date of Birth</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  formatDate(member?.date_of_birth)
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Age</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2"
                  >{{
                    calculateAge(member?.date_of_birth)
                  }}
                  years</v-list-item-subtitle
                >
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Employment Information -->
      <v-col cols="12" md="4">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-info text-white">
            Employment Details
          </v-card-title>
          <v-card-text>
            <v-list density="compact">
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Scheme</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  member?.scheme_name || 'N/A'
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Scheme Category</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  member?.scheme_category || 'N/A'
                }}</v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Annual Salary</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  formatCurrency(member?.annual_salary)
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Entry Date</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  formatDate(member?.entry_date)
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Employee Number</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  member?.employee_number || 'N/A'
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item v-if="member?.effective_exit_date">
                <v-list-item-title class="text-caption text-grey"
                  >Effective Exit Date</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  formatDate(member?.effective_exit_date)
                }}</v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Status</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">
                  <v-chip :color="getStatusColor(member?.status)" size="small">
                    {{ (member?.status || 'active').toUpperCase() }}
                  </v-chip>
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Contact Information -->
      <v-col cols="12" md="4">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-success text-white">
            Contact Information
          </v-card-title>
          <v-card-text>
            <v-list density="compact">
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Email</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  member?.email || 'Not provided'
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Phone</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  member?.phone_number || 'Not provided'
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Address</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  formatAddress(member)
                }}</v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Beneficiaries Summary -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title
            class="text-subtitle-1 bg-warning d-flex justify-space-between align-center"
          >
            <span>Beneficiaries ({{ beneficiaries.length }})</span>
            <v-btn
              color="primary"
              rounded
              size="small"
              variant="elevated"
              @click="$emit('manage-beneficiaries')"
            >
              Manage Beneficiaries
            </v-btn>
          </v-card-title>
          <v-card-text>
            <div
              v-if="beneficiaries.length === 0"
              class="text-center text-grey py-4"
            >
              No beneficiaries added yet
            </div>
            <v-row v-else>
              <v-col
                v-for="beneficiary in beneficiaries"
                :key="beneficiary.id"
                cols="12"
                md="6"
                lg="4"
              >
                <v-card variant="outlined" size="small">
                  <v-card-text class="pa-3">
                    <div class="d-flex justify-space-between align-start mb-2">
                      <div class="text-subtitle-2 font-weight-bold">{{
                        beneficiary.full_name
                      }}</div>
                      <v-chip size="x-small" color="primary"
                        >{{ beneficiary.allocation_percentage }}%</v-chip
                      >
                    </div>
                    <div class="text-caption text-grey mb-1">{{
                      formatRelationship(beneficiary.relationship)
                    }}</div>
                    <div class="text-caption">{{ beneficiary.id_number }}</div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Total Allocation Alert -->
            <v-alert
              v-if="beneficiaries.length > 0"
              :type="totalAllocation === 100 ? 'success' : 'warning'"
              variant="tonal"
              class="mt-3"
            >
              Total beneficiary allocation: {{ totalAllocation }}%
              <span v-if="totalAllocation !== 100"> (Should be 100%)</span>
            </v-alert>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Quick Actions -->
      <v-col cols="12">
        <v-card variant="outlined">
          <v-card-title
            class="text-subtitle-1 bg-grey-lighten-3 d-flex justify-space-between align-center"
          >
            <span>Quick Actions</span>
            <v-btn
              rounded
              color="primary"
              size="small"
              variant="elevated"
              @click="handleEditMember"
            >
              <v-icon left>mdi-pencil</v-icon>
              Edit Member
            </v-btn>
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="6" md="3">
                <v-btn
                  size="small"
                  rounded
                  color="info"
                  variant="outlined"
                  block
                  @click="$emit('view-claims')"
                >
                  <v-icon left>mdi-file-document-outline</v-icon>
                  View Claims
                </v-btn>
              </v-col>
              <v-col cols="6" md="3">
                <v-btn
                  size="small"
                  rounded
                  color="success"
                  variant="outlined"
                  block
                  @click="downloadMemberCertificate"
                >
                  <v-icon left>mdi-certificate</v-icon>
                  Certificate
                </v-btn>
              </v-col>
              <v-col cols="6" md="3">
                <v-btn
                  size="small"
                  rounded
                  color="warning"
                  variant="outlined"
                  block
                  @click="viewBenefitSummary"
                >
                  <v-icon left>mdi-chart-line</v-icon>
                  Benefit Summary
                </v-btn>
              </v-col>
              <v-col cols="6" md="3">
                <v-btn
                  size="small"
                  rounded
                  color="primary"
                  variant="outlined"
                  block
                  @click="viewMemberHistory"
                >
                  <v-icon left>mdi-history</v-icon>
                  History
                </v-btn>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Member Benefit Summary Dialog -->
    <member-benefit-summary
      v-model="showBenefitSummary"
      :member-id="member?.id || null"
      :member-name="member?.member_name || ''"
    />

    <!-- Member Edit Dialog -->
    <member-edit-dialog
      v-model="showEditDialog"
      :member="member"
      @member-updated="handleMemberUpdated"
    />

    <!-- Member History Dialog -->
    <member-history-dialog
      v-model="showHistoryDialog"
      :member-id="member?.id || null"
      :member-name="member?.member_name || ''"
    />
  </v-container>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import MemberBenefitSummary from './MemberBenefitSummary.vue'
import MemberEditDialog from './MemberEditDialog.vue'
import MemberHistoryDialog from './MemberHistoryDialog.vue'

interface Member {
  id?: number
  member_name: string
  member_id_number: string
  gender: string
  date_of_birth: string | Date | null
  email?: string
  phone_number?: string
  scheme_name: string
  scheme_category: string
  annual_salary: number
  entry_date: string | Date | null
  effective_exit_date?: string | Date | null
  employee_number?: string
  status?: string
  address_line_1?: string
  address_line_2?: string
  city?: string
  province?: string
  postal_code?: string
}

interface Beneficiary {
  id: number
  full_name: string
  relationship: string
  id_number: string
  allocation_percentage: number
}

interface Props {
  member: Member | null
  beneficiaries: Beneficiary[]
}

interface Emits {
  (e: 'manage-beneficiaries'): void
  (e: 'view-claims'): void
  (e: 'edit-member'): void
  (e: 'member-updated', member: Member): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// State for benefit summary dialog
const showBenefitSummary = ref(false)

// State for edit dialog
const showEditDialog = ref(false)

// State for history dialog
const showHistoryDialog = ref(false)

// Computed properties
const totalAllocation = computed(() => {
  return props.beneficiaries.reduce(
    (total: number, beneficiary: Beneficiary) =>
      total + (beneficiary.allocation_percentage || 0),
    0
  )
})

// Methods
const formatDate = (date: string | Date | null | undefined) => {
  if (!date) return 'Not provided'
  return new Date(date).toLocaleDateString('en-ZA', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const calculateAge = (dateOfBirth: string | Date | null | undefined) => {
  if (!dateOfBirth) return 'Unknown'

  const today = new Date()
  const birthDate = new Date(dateOfBirth)
  let age = today.getFullYear() - birthDate.getFullYear()
  const monthDiff = today.getMonth() - birthDate.getMonth()

  if (
    monthDiff < 0 ||
    (monthDiff === 0 && today.getDate() < birthDate.getDate())
  ) {
    age--
  }

  return age
}

const formatCurrency = (amount: number | null | undefined) => {
  if (!amount) return 'Not specified'
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR'
  }).format(amount)
}

const formatAddress = (member: Member | null) => {
  if (!member) return 'Not provided'

  const addressParts = [
    member.address_line_1,
    member.address_line_2,
    member.city,
    member.province,
    member.postal_code
  ].filter(Boolean)

  return addressParts.length > 0 ? addressParts.join(', ') : 'Not provided'
}

const formatRelationship = (relationship: string) => {
  return relationship
    .split('_')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ')
}

const getStatusColor = (status: string | undefined) => {
  switch ((status || 'active').toLowerCase()) {
    case 'active':
      return 'success'
    case 'inactive':
      return 'error'
    case 'pending':
      return 'warning'
    case 'suspended':
      return 'orange'
    default:
      return 'grey'
  }
}

const downloadMemberCertificate = () => {
  // Implementation for downloading member certificate
  console.log('Download certificate for member:', props.member?.member_name)
}

const viewBenefitSummary = () => {
  if (props.member?.id) {
    showBenefitSummary.value = true
  }
}

const viewMemberHistory = () => {
  if (props.member?.id) {
    showHistoryDialog.value = true
  }
  console.log('Opening history dialog for member:', props.member?.member_name)
}

const handleEditMember = () => {
  showEditDialog.value = true
  // Also emit the event for parent components that might need it
  emit('edit-member')
}

const handleMemberUpdated = (updatedMember: Member) => {
  emit('member-updated', updatedMember)
}
</script>

<style scoped>
.v-card-title {
  padding: 12px 16px;
}

.v-card-text {
  padding: 16px;
}

.v-list-item-title {
  font-size: 0.75rem;
  margin-bottom: 4px;
}

.v-list-item-subtitle {
  font-size: 0.875rem;
  font-weight: 500;
}
</style>
