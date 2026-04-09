<template>
  <v-container>
    <v-row>
      <!-- Basic Information -->
      <v-col cols="12" md="6">
        <v-card variant="outlined">
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
                  beneficiary?.full_name
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Relationship</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">
                  <v-chip size="small" color="info">{{
                    formatRelationship(beneficiary?.relationship)
                  }}</v-chip>
                </v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >ID Number</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  beneficiary?.id_number
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Gender</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  beneficiary?.gender
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Date of Birth</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">
                  {{ formatDate(beneficiary?.date_of_birth) }}
                  <v-chip
                    v-if="isMinor"
                    size="small"
                    color="warning"
                    class="ml-2"
                    >Minor</v-chip
                  >
                </v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Age</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2"
                  >{{
                    calculateAge(beneficiary?.date_of_birth)
                  }}
                  years</v-list-item-subtitle
                >
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Contact Information -->
      <v-col cols="12" md="6">
        <v-card variant="outlined">
          <v-card-title class="text-subtitle-1 bg-info text-white">
            Contact Information
          </v-card-title>
          <v-card-text>
            <v-list density="compact">
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Contact Number</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">
                  {{ beneficiary?.contact_number || 'Not provided' }}
                </v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Email</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">
                  {{ beneficiary?.email || 'Not provided' }}
                </v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Address</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">
                  {{ beneficiary?.address || 'Not provided' }}
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Allocation and Benefits -->
      <v-col cols="12">
        <v-card variant="outlined">
          <v-card-title class="text-subtitle-1 bg-success text-white">
            Benefit Allocation
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <div class="text-center">
                  <div class="text-h4 text-success"
                    >{{ beneficiary?.allocation_percentage }}%</div
                  >
                  <div class="text-subtitle-2 text-grey"
                    >Allocation Percentage</div
                  >
                </div>
              </v-col>
              <v-col cols="12" md="6">
                <div class="text-subtitle-2 text-grey mb-2"
                  >Applicable Benefits:</div
                >
                <div class="d-flex flex-wrap ga-1">
                  <v-chip
                    v-for="benefit in beneficiary?.benefit_types"
                    :key="benefit"
                    size="small"
                    color="success"
                    variant="outlined"
                  >
                    {{ benefit }}
                  </v-chip>
                </div>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Guardian Information (if minor) -->
      <v-col v-if="isMinor" cols="12">
        <v-card variant="outlined" color="warning">
          <v-card-title class="text-subtitle-1 bg-warning">
            Guardian Information
          </v-card-title>
          <v-card-text>
            <v-list density="compact">
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Guardian Name</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">
                  {{ beneficiary?.guardian_name || 'Not provided' }}
                </v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Guardian Relationship</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">
                  {{
                    formatRelationship(beneficiary?.guardian_relationship) ||
                    'Not provided'
                  }}
                </v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Guardian ID Number</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">
                  {{ beneficiary?.guardian_id_number || 'Not provided' }}
                </v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Guardian Contact</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">
                  {{ beneficiary?.guardian_contact || 'Not provided' }}
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Banking Information -->
      <v-col v-if="hasBankingInfo" cols="12">
        <v-card variant="outlined">
          <v-card-title class="text-subtitle-1 bg-grey-lighten-2">
            Banking Information
          </v-card-title>
          <v-card-text>
            <v-list density="compact">
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Bank Name</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  beneficiary?.bank_name
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Branch Code</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  beneficiary?.branch_code
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Account Number</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">{{
                  beneficiary?.account_number
                }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title class="text-caption text-grey"
                  >Account Type</v-list-item-title
                >
                <v-list-item-subtitle class="text-body-2">
                  {{ formatAccountType(beneficiary?.account_type) }}
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Status and Metadata -->
      <v-col cols="12">
        <v-card variant="outlined">
          <v-card-title class="text-subtitle-1 bg-grey-lighten-3">
            Status Information
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <div class="text-subtitle-2 text-grey mb-2">Status:</div>
                <v-chip
                  :color="getStatusColor(beneficiary?.status)"
                  size="large"
                >
                  {{ (beneficiary?.status || 'active').toUpperCase() }}
                </v-chip>
              </v-col>
              <v-col cols="12" md="6">
                <div
                  v-if="beneficiary?.created_date"
                  class="text-subtitle-2 text-grey mb-2"
                >
                  Added: {{ formatDate(beneficiary?.created_date) }}
                </div>
                <div
                  v-if="beneficiary?.updated_date"
                  class="text-subtitle-2 text-grey"
                >
                  Last Updated: {{ formatDate(beneficiary?.updated_date) }}
                </div>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Beneficiary {
  id?: number
  full_name: string
  relationship: string
  id_number: string
  gender: string
  date_of_birth: string | Date | null
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
  created_date?: string | Date
  updated_date?: string | Date
}

interface Props {
  beneficiary: Beneficiary | null
}

const props = defineProps<Props>()

// Computed properties
const isMinor = computed(() => {
  if (!props.beneficiary?.date_of_birth) return false
  const today = new Date()
  const birthDate = new Date(props.beneficiary.date_of_birth)
  const age = today.getFullYear() - birthDate.getFullYear()
  const monthDiff = today.getMonth() - birthDate.getMonth()

  if (
    monthDiff < 0 ||
    (monthDiff === 0 && today.getDate() < birthDate.getDate())
  ) {
    return age - 1 < 18
  }
  return age < 18
})

const hasBankingInfo = computed(() => {
  return (
    props.beneficiary?.bank_name ||
    props.beneficiary?.account_number ||
    props.beneficiary?.branch_code
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

const formatRelationship = (relationship: string | undefined) => {
  if (!relationship) return 'Unknown'

  return relationship
    .split('_')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ')
}

const formatAccountType = (accountType: string | undefined) => {
  if (!accountType) return 'Not specified'

  const types: Record<string, string> = {
    savings: 'Savings Account',
    current: 'Current/Cheque Account',
    transmission: 'Transmission Account'
  }

  return types[accountType] || accountType
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
