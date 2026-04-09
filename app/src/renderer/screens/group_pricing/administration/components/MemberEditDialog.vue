<template>
  <v-dialog v-model="dialog" max-width="900" persistent scrollable>
    <v-card>
      <v-card-title
        class="d-flex justify-space-between align-center bg-primary text-white"
      >
        <div class="d-flex align-center">
          <v-icon class="mr-2">mdi-pencil</v-icon>
          <span>Edit Member - {{ member?.member_name }}</span>
        </div>
        <v-btn icon="mdi-close" variant="text" @click="closeDialog" />
      </v-card-title>

      <v-card-text class="pa-0">
        <v-container v-if="loading" class="text-center py-8">
          <v-progress-circular indeterminate color="primary" size="64" />
          <div class="mt-4 text-h6">Updating member...</div>
        </v-container>

        <v-form v-else ref="form" @submit.prevent="handleSubmit">
          <v-container class="py-4">
            <!-- Personal Information -->
            <v-row>
              <v-col cols="12">
                <v-card variant="outlined" class="mb-4">
                  <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
                    Personal Information
                  </v-card-title>
                  <v-card-text>
                    <v-row>
                      <v-col cols="12" md="6">
                        <v-text-field
                          v-model="formData.member_name"
                          label="Full Name *"
                          variant="outlined"
                          density="compact"
                          :rules="[rules.required]"
                          required
                        />
                      </v-col>
                      <v-col cols="12" md="6">
                        <v-text-field
                          v-model="formData.member_id_number"
                          label="ID Number"
                          variant="outlined"
                          density="compact"
                          readonly
                          hint="ID Number cannot be changed"
                          persistent-hint
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
            </v-row>

            <!-- Contact Information -->
            <v-row>
              <v-col cols="12">
                <v-card variant="outlined" class="mb-4">
                  <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
                    Contact Information
                  </v-card-title>
                  <v-card-text>
                    <v-row>
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
                      <v-col cols="12" md="6">
                        <v-text-field
                          v-model="formData.phone_number"
                          label="Phone Number"
                          variant="outlined"
                          density="compact"
                        />
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Employment Details -->
            <v-row>
              <v-col cols="12">
                <v-card variant="outlined" class="mb-4">
                  <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
                    Employment Details
                  </v-card-title>
                  <v-card-text>
                    <v-row>
                      <v-col cols="12" md="6">
                        <v-text-field
                          v-model="formData.employee_number"
                          label="Employee Number"
                          variant="outlined"
                          density="compact"
                        />
                      </v-col>
                      <v-col cols="12" md="6">
                        <v-text-field
                          v-model.number="formData.annual_salary"
                          label="Annual Salary *"
                          type="number"
                          prefix="R"
                          variant="outlined"
                          density="compact"
                          :rules="[rules.required, rules.salary]"
                          required
                        />
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col cols="12" md="6">
                        <v-date-input
                          v-model="formData.entry_date"
                          label="Entry Date *"
                          variant="outlined"
                          density="compact"
                          :rules="[rules.required]"
                          hide-actions
                          required
                        />
                      </v-col>
                      <v-col cols="12" md="6">
                        <v-select
                          v-model="formData.status"
                          :items="statusOptions"
                          label="Status"
                          variant="outlined"
                          density="compact"
                        />
                      </v-col>
                    </v-row>
                    <v-row v-if="formData.status === 'INACTIVE'">
                      <v-col cols="12" md="6">
                        <v-date-input
                          v-model="formData.effective_exit_date"
                          label="Effective Exit Date"
                          variant="outlined"
                          density="compact"
                          hide-actions
                        />
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Address Information -->
            <v-row>
              <v-col cols="12">
                <v-card variant="outlined" class="mb-4">
                  <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
                    Address Information
                  </v-card-title>
                  <v-card-text>
                    <v-row>
                      <v-col cols="12">
                        <v-text-field
                          v-model="formData.address_line_1"
                          label="Address Line 1"
                          variant="outlined"
                          density="compact"
                        />
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col cols="12" md="6">
                        <v-text-field
                          v-model="formData.address_line_2"
                          label="Address Line 2"
                          variant="outlined"
                          density="compact"
                        />
                      </v-col>
                      <v-col cols="12" md="6">
                        <v-text-field
                          v-model="formData.city"
                          label="City"
                          variant="outlined"
                          density="compact"
                        />
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col cols="12" md="6">
                        <v-select
                          v-model="formData.province"
                          :items="provinces"
                          label="Province"
                          variant="outlined"
                          density="compact"
                        />
                      </v-col>
                      <v-col cols="12" md="6">
                        <v-text-field
                          v-model="formData.postal_code"
                          label="Postal Code"
                          variant="outlined"
                          density="compact"
                        />
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-container>
        </v-form>
      </v-card-text>

      <v-card-actions class="px-4 py-3">
        <v-btn
          color="grey"
          variant="outlined"
          :disabled="loading"
          @click="closeDialog"
        >
          Cancel
        </v-btn>
        <v-spacer />
        <v-btn
          color="primary"
          variant="elevated"
          :loading="loading"
          :disabled="!isFormValid"
          @click="handleSubmit"
        >
          Update Member
        </v-btn>
      </v-card-actions>
    </v-card>

    <!-- Success/Error Snackbar -->
    <v-snackbar v-model="snackbar" :timeout="4000" :color="snackbarColor">
      {{ snackbarMessage }}
    </v-snackbar>

    <!-- Salary Change Confirmation Dialog -->
    <v-dialog v-model="showSalaryConfirmation" max-width="500" persistent>
      <v-card>
        <v-card-title class="d-flex align-center bg-warning text-white">
          <v-icon class="mr-2">mdi-alert</v-icon>
          <span>Salary Change Confirmation</span>
        </v-card-title>

        <v-card-text class="py-6">
          <div class="text-body-1 mb-4">
            <strong>Warning:</strong> Changing this member's annual salary will
            impact several connected systems and calculations:
          </div>

          <v-list dense class="mb-4">
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title
                  >• Premium calculations and quotes</v-list-item-title
                >
              </v-list-item-content>
            </v-list-item>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title
                  >• Benefit calculations and coverage
                  amounts</v-list-item-title
                >
              </v-list-item-content>
            </v-list-item>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title
                  >• Risk assessments and underwriting</v-list-item-title
                >
              </v-list-item-content>
            </v-list-item>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title
                  >• Actuarial reports and analytics</v-list-item-title
                >
              </v-list-item-content>
            </v-list-item>
          </v-list>

          <div class="text-body-2 text-grey-darken-1">
            <strong>Current Salary:</strong> R{{
              originalSalary.toLocaleString()
            }}<br />
            <strong>New Salary:</strong> R{{
              Number(formData.annual_salary).toLocaleString()
            }}
          </div>

          <div class="mt-4 text-body-1">
            Are you sure you want to proceed with this salary change?
          </div>
        </v-card-text>

        <v-card-actions class="px-4 py-3">
          <v-btn color="grey" variant="outlined" @click="cancelSalaryChange">
            Cancel
          </v-btn>
          <v-spacer />
          <v-btn
            color="warning"
            variant="elevated"
            @click="confirmSalaryChange"
          >
            Yes, Update Salary
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { VDateInput } from 'vuetify/labs/VDateInput'
import GroupPricingService from '@/renderer/api/GroupPricingService'

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

interface Props {
  modelValue: boolean
  member: Member | null
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'member-updated', member: Member): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// State
const loading = ref(false)
const form = ref(null)
const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')
const showSalaryConfirmation = ref(false)
const originalSalary = ref<number>(0)

// Form data
const formData = ref<Member>({
  member_name: '',
  member_id_number: '',
  gender: '',
  date_of_birth: null,
  email: '',
  phone_number: '',
  scheme_name: '',
  scheme_category: '',
  annual_salary: 0,
  entry_date: null,
  effective_exit_date: null,
  employee_number: '',
  status: 'ACTIVE',
  address_line_1: '',
  address_line_2: '',
  city: '',
  province: '',
  postal_code: ''
})

// Computed
const dialog = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const isFormValid = computed(() => {
  return (
    formData.value.member_name &&
    formData.value.gender &&
    formData.value.date_of_birth &&
    formData.value.entry_date &&
    formData.value.annual_salary &&
    formData.value.annual_salary > 0
  )
})

// Form options
const genderOptions = [
  { title: 'Male', value: 'Male' },
  { title: 'Female', value: 'Female' }
]

const statusOptions = [
  { title: 'Active', value: 'ACTIVE' },
  { title: 'Inactive', value: 'INACTIVE' },
  { title: 'Suspended', value: 'SUSPENDED' },
  { title: 'Pending', value: 'PENDING' }
]

const provinces = [
  { title: 'Western Cape', value: 'WC' },
  { title: 'Eastern Cape', value: 'EC' },
  { title: 'Northern Cape', value: 'NC' },
  { title: 'Free State', value: 'FS' },
  { title: 'KwaZulu-Natal', value: 'KZN' },
  { title: 'North West', value: 'NW' },
  { title: 'Gauteng', value: 'GP' },
  { title: 'Mpumalanga', value: 'MP' },
  { title: 'Limpopo', value: 'LP' }
]

// Validation rules
const rules = {
  required: (value: any) => !!value || 'Field is required',
  email: (value: string) => {
    if (!value) return true
    const pattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return pattern.test(value) || 'Invalid email format'
  },
  dateOfBirth: (value: Date | null) => {
    if (!value) return 'Date of birth is required'
    const today = new Date()
    const birthDate = new Date(value)
    const age = today.getFullYear() - birthDate.getFullYear()
    return (age >= 16 && age <= 75) || 'Age must be between 16 and 75 years'
  },
  salary: (value: number | null) => {
    if (!value) return 'Annual salary is required'
    return value > 0 || 'Salary must be greater than 0'
  }
}

// Watch for dialog opening and populate form
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue && props.member) {
      populateForm()
    }
  }
)

watch(
  () => props.member,
  (newMember) => {
    if (newMember && props.modelValue) {
      populateForm()
    }
  }
)

// Methods
const populateForm = () => {
  if (!props.member) return

  formData.value = {
    ...props.member,
    date_of_birth: props.member.date_of_birth
      ? new Date(props.member.date_of_birth)
      : null,
    entry_date: props.member.entry_date
      ? new Date(props.member.entry_date)
      : null,
    effective_exit_date: props.member.effective_exit_date
      ? new Date(props.member.effective_exit_date)
      : null
  }

  // Store original salary for comparison
  originalSalary.value = props.member.annual_salary || 0
}

const handleSubmit = async () => {
  if (!form.value || !props.member?.id) return

  const { valid } = await (form.value as any).validate()
  if (!valid) return

  // Check if salary has changed
  const currentSalary = Number(formData.value.annual_salary)
  if (currentSalary !== originalSalary.value) {
    showSalaryConfirmation.value = true
    return
  }

  // If salary hasn't changed, proceed with update
  await submitMemberUpdate()
}

const submitMemberUpdate = async () => {
  loading.value = true

  try {
    // Prepare data for API
    const memberData = {
      ...formData.value,
      annual_salary: Number(formData.value.annual_salary)
    }

    const response = await GroupPricingService.updateMember(
      props.member!.id!,
      memberData
    )

    showSnackbar('Member updated successfully', 'success')
    emit('member-updated', response.data)

    // Close dialog after short delay
    setTimeout(() => {
      closeDialog()
    }, 1000)
  } catch (error: any) {
    console.error('Error updating member:', error)
    const errorMessage =
      error.response?.data?.message || 'Failed to update member'
    showSnackbar(errorMessage, 'error')
  } finally {
    loading.value = false
  }
}

const confirmSalaryChange = async () => {
  showSalaryConfirmation.value = false
  await submitMemberUpdate()
}

const cancelSalaryChange = () => {
  showSalaryConfirmation.value = false
}

const closeDialog = () => {
  dialog.value = false
  // Reset form after dialog closes
  setTimeout(() => {
    resetForm()
  }, 300)
}

const resetForm = () => {
  formData.value = {
    member_name: '',
    member_id_number: '',
    gender: '',
    date_of_birth: null,
    email: '',
    phone_number: '',
    annual_salary: 0,
    entry_date: null,
    scheme_name: '',
    scheme_category: '',
    effective_exit_date: null,
    employee_number: '',
    status: 'ACTIVE',
    address_line_1: '',
    address_line_2: '',
    city: '',
    province: '',
    postal_code: ''
  }

  // Reset original salary tracking
  originalSalary.value = 0
}

const showSnackbar = (message: string, color: string) => {
  snackbarMessage.value = message
  snackbarColor.value = color
  snackbar.value = true
}
</script>

<style scoped>
.v-card-title {
  padding: 12px 16px;
}

.v-card-text {
  padding: 16px;
}

.v-card-actions {
  padding: 12px 16px;
}

@media (max-width: 600px) {
  .v-dialog {
    margin: 8px;
  }
}
</style>
