<template>
  <v-form ref="form" @submit.prevent="handleSubmit">
    <v-row>
      <!-- Basic Information -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
            Beneficiary Information
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.full_name"
                  label="Full Name *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  required
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.relationship"
                  :items="relationshipOptions"
                  label="Relationship *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required, rules.uniqueRelationship]"
                  required
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.id_type"
                  :items="idTypes"
                  label="ID Type *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  required
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.id_number"
                  label="ID Number *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required, rules.idNumber]"
                  required
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.gender"
                  :items="genderOptions"
                  label="Gender *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  required
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-date-input
                  v-model="formData.date_of_birth"
                  label="Date of Birth *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required, rules.dateOfBirth]"
                  hide-actions
                  required
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Contact Information -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
            Contact Information
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.contact_number"
                  label="Contact Number"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.phone]"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.email"
                  label="Email Address"
                  type="email"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.email]"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="formData.address"
                  label="Address"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Allocation and Benefits -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
            Benefit Allocation
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model.number="formData.allocation_percentage"
                  label="Allocation Percentage *"
                  type="number"
                  suffix="%"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required, rules.percentage]"
                  :hint="`Remaining allocation: ${remainingAllocation}%`"
                  persistent-hint
                  required
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.benefit_types"
                  :items="benefitTypeOptions"
                  label="Applicable Benefits"
                  variant="outlined"
                  density="compact"
                  multiple
                  chips
                  :hint="'Select which benefits this beneficiary applies to'"
                  persistent-hint
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Guardian Information (for minors) -->
      <v-col v-if="isMinor" cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-warning-lighten-4">
            Guardian Information (Required for Minors)
          </v-card-title>
          <v-card-text>
            <v-alert type="info" variant="tonal" class="mb-4">
              This beneficiary is under 18 years old. Guardian information is
              required.
            </v-alert>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.guardian_name"
                  label="Guardian Full Name *"
                  variant="outlined"
                  density="compact"
                  :rules="isMinor ? [rules.required] : []"
                  :required="isMinor"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.guardian_relationship"
                  :items="guardianRelationshipOptions"
                  label="Guardian Relationship *"
                  variant="outlined"
                  density="compact"
                  :rules="isMinor ? [rules.required] : []"
                  :required="isMinor"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.guardian_id_number"
                  label="Guardian ID Number *"
                  variant="outlined"
                  density="compact"
                  :rules="isMinor ? [rules.required, rules.idNumber] : []"
                  :required="isMinor"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.guardian_contact"
                  label="Guardian Contact Number *"
                  variant="outlined"
                  density="compact"
                  :rules="isMinor ? [rules.required, rules.phone] : []"
                  :required="isMinor"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Banking Information -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
            Banking Information (Optional)
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.bank_name"
                  label="Bank Name"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.branch_code"
                  label="Branch Code"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.account_number"
                  label="Account Number"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.account_type"
                  :items="accountTypes"
                  label="Account Type"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Action Buttons -->
      <v-col cols="12">
        <v-card-actions class="justify-end">
          <v-btn color="grey" variant="outlined" @click="$emit('cancel')"
            >Cancel</v-btn
          >
          <v-btn
            color="primary"
            type="submit"
            :loading="loading"
            :disabled="!isFormValid"
          >
            {{ isEditMode ? 'Update Beneficiary' : 'Add Beneficiary' }}
          </v-btn>
        </v-card-actions>
      </v-col>
    </v-row>
  </v-form>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { VDateInput } from 'vuetify/labs/VDateInput'

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
}

interface Props {
  beneficiary?: Beneficiary | null
  memberId: number | null
  isEditMode: boolean
  existingBeneficiaries: Beneficiary[]
}

interface Emits {
  (e: 'save', beneficiary: Beneficiary): void
  (e: 'cancel'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const form = ref(null)
const loading = ref(false)

// Form data with default values
const formData = ref<Beneficiary>({
  full_name: '',
  relationship: '',
  id_type: 'RSA_ID',
  id_number: '',
  gender: '',
  date_of_birth: null,
  contact_number: '',
  email: '',
  address: '',
  allocation_percentage: 0,
  benefit_types: ['GLA', 'GFF'],
  guardian_name: '',
  guardian_relationship: '',
  guardian_id_number: '',
  guardian_contact: '',
  bank_name: '',
  branch_code: '',
  account_number: '',
  account_type: '',
  status: 'active'
})

// Form options
const relationshipOptions = [
  { title: 'Spouse', value: 'spouse' },
  { title: 'Child', value: 'child' },
  { title: 'Parent', value: 'parent' },
  { title: 'Sibling', value: 'sibling' },
  { title: 'Grandparent', value: 'grandparent' },
  { title: 'Grandchild', value: 'grandchild' },
  { title: 'Other Family', value: 'other_family' },
  { title: 'Friend', value: 'friend' },
  { title: 'Estate', value: 'estate' },
  { title: 'Trust', value: 'trust' },
  { title: 'Charity', value: 'charity' }
]

const idTypes = [
  { title: 'RSA ID', value: 'RSA_ID' },
  { title: 'Passport', value: 'PASSPORT' },
  { title: 'Birth Certificate', value: 'BIRTH_CERT' }
]

const genderOptions = [
  { title: 'Male', value: 'Male' },
  { title: 'Female', value: 'Female' }
]

const benefitTypeOptions = [
  { title: 'Group Life Assurance (GLA)', value: 'GLA' },
  { title: 'Spouse Group Life Assurance (SGLA)', value: 'SGLA' },
  { title: 'Permanent Total Disability (PTD)', value: 'PTD' },
  { title: 'Critical Illness (CI)', value: 'CI' },
  { title: 'Temporary Total Disability (TTD)', value: 'TTD' },
  { title: 'Personal Health Insurance (PHI)', value: 'PHI' },
  { title: 'Group Family Funeral (GFF)', value: 'GFF' }
]

const guardianRelationshipOptions = [
  { title: 'Parent', value: 'parent' },
  { title: 'Legal Guardian', value: 'legal_guardian' },
  { title: 'Grandparent', value: 'grandparent' },
  { title: 'Uncle/Aunt', value: 'uncle_aunt' },
  { title: 'Other Relative', value: 'other_relative' }
]

const accountTypes = [
  { title: 'Savings', value: 'savings' },
  { title: 'Current/Cheque', value: 'current' },
  { title: 'Transmission', value: 'transmission' }
]

// Computed properties
const isMinor = computed(() => {
  if (!formData.value.date_of_birth) return false
  const today = new Date()
  const birthDate = new Date(formData.value.date_of_birth)
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

const remainingAllocation = computed(() => {
  const totalExisting = props.existingBeneficiaries
    .filter((b) => (props.isEditMode ? b.id !== props.beneficiary?.id : true))
    .reduce((sum, b) => sum + (b.allocation_percentage || 0), 0)

  return 100 - totalExisting
})

const isFormValid = computed(() => {
  return (
    formData.value.full_name &&
    formData.value.relationship &&
    formData.value.id_number &&
    formData.value.gender &&
    formData.value.date_of_birth &&
    formData.value.allocation_percentage > 0 &&
    formData.value.allocation_percentage <= remainingAllocation.value &&
    (!isMinor.value ||
      (formData.value.guardian_name && formData.value.guardian_contact))
  )
})

// Validation rules
const rules = {
  required: (value: any) => !!value || 'Field is required',
  email: (value: string) => {
    if (!value) return true
    const pattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return pattern.test(value) || 'Invalid email format'
  },
  phone: (value: string) => {
    if (!value) return true
    const pattern = /^[\d\s\-+()]{10,15}$/
    return pattern.test(value) || 'Invalid phone number format'
  },
  idNumber: (value: string) => {
    if (!value) return 'ID number is required'
    if (formData.value.id_type === 'RSA_ID') {
      return value.length === 13 || 'RSA ID number must be 13 digits'
    }
    return true
  },
  dateOfBirth: (value: Date | null) => {
    if (!value) return 'Date of birth is required'
    const today = new Date()
    const birthDate = new Date(value)
    return birthDate <= today || 'Date of birth cannot be in the future'
  },
  percentage: (value: number) => {
    if (!value || value <= 0) return 'Allocation must be greater than 0'
    if (remainingAllocation.value <= 0) {
      return 'No remaining allocation available'
    }
    if (value > remainingAllocation.value) {
      return `Allocation cannot exceed ${remainingAllocation.value}%`
    }
    return true
  },
  uniqueRelationship: (value: string) => {
    if (!value) return 'Relationship is required'

    // Check for unique relationships that should only have one beneficiary
    const uniqueRelationships = ['spouse', 'estate', 'trust']
    if (uniqueRelationships.includes(value)) {
      const existing = props.existingBeneficiaries.find(
        (b) =>
          b.relationship === value &&
          (props.isEditMode ? b.id !== props.beneficiary?.id : true)
      )
      if (existing) {
        return `Only one ${value} beneficiary is allowed`
      }
    }
    return true
  }
}

// Methods
const handleSubmit = async () => {
  if (!form.value) return

  const { valid } = await (form.value as any).validate()
  if (!valid) return

  loading.value = true
  try {
    const beneficiaryData = {
      ...formData.value,
      member_id: props.memberId,
      allocation_percentage: Number(formData.value.allocation_percentage)
    }

    emit('save', beneficiaryData)
  } finally {
    loading.value = false
  }
}

// Watchers
watch(
  () => props.beneficiary,
  (newBeneficiary) => {
    if (newBeneficiary && props.isEditMode) {
      // Populate form with existing beneficiary data
      formData.value = {
        ...formData.value,
        ...newBeneficiary,
        date_of_birth: newBeneficiary.date_of_birth
          ? new Date(newBeneficiary.date_of_birth)
          : null
      }
    } else if (!props.isEditMode) {
      // Reset form for new beneficiary
      formData.value = {
        full_name: '',
        relationship: '',
        id_type: 'RSA_ID',
        id_number: '',
        gender: '',
        date_of_birth: null,
        contact_number: '',
        email: '',
        address: '',
        allocation_percentage: remainingAllocation.value,
        benefit_types: ['GLA', 'GFF'],
        guardian_name: '',
        guardian_relationship: '',
        guardian_id_number: '',
        guardian_contact: '',
        bank_name: '',
        branch_code: '',
        account_number: '',
        account_type: '',
        status: 'active'
      }
    }
  },
  { immediate: true }
)

// Set default allocation to remaining percentage for new beneficiaries
watch(remainingAllocation, (newValue) => {
  if (!props.isEditMode && formData.value.allocation_percentage === 0) {
    formData.value.allocation_percentage = newValue
  }
})
</script>

<style scoped>
.v-card-title {
  padding: 12px 16px;
}

.v-card-text {
  padding: 16px;
}
</style>
