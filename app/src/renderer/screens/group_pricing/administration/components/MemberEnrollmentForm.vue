<template>
  <v-form ref="form" @submit.prevent="handleSubmit">
    <v-row>
      <!-- Basic Information Section -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
            Basic Information
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
                <v-select
                  v-model="formData.member_id_type"
                  :items="idTypes"
                  label="ID Type *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  required
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.member_id_number"
                  :label="idNumberLabel"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required, rules.idNumber]"
                  required
                />
              </v-col>
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
            </v-row>
            <v-row>
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
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.phone_number"
                  label="Phone Number"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.employee_number"
                  label="Employee Number"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Scheme and Employment Details -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
            Scheme & Employment Details
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <v-autocomplete
                  v-model="formData.scheme_id"
                  :items="schemes"
                  item-title="name"
                  item-value="id"
                  label="Scheme *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  required
                  @update:model-value="onSchemeChange"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.scheme_category"
                  :items="schemeCategories"
                  label="Scheme Category *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  required
                  @update:model-value="onSchemeCategoryChange"
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
                  :min="minDate"
                  :rules="[rules.required, rules.dateOfEntry]"
                  hide-actions
                  required
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
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Address Information -->
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

      <!-- Benefit Selection -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
            Benefit Selection
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12">
                <v-alert type="info" variant="tonal" class="mb-4">
                  Select the benefits this member should be enrolled in. Benefit
                  amounts will be calculated based on salary and scheme
                  configuration.
                </v-alert>
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-checkbox
                  v-model="formData.benefits.gla_enabled"
                  :label="getBenefitLabel('GLA')"
                  :disabled="isCompulsory"
                  color="primary"
                />
                <v-text-field
                  v-if="formData.benefits.gla_enabled"
                  v-model.number="formData.benefits.gla_multiple"
                  label="Salary Multiple"
                  type="number"
                  variant="outlined"
                  density="compact"
                  class="ml-8"
                  step="0.5"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-checkbox
                  v-model="formData.benefits.sgla_enabled"
                  :label="getBenefitLabel('SGLA')"
                  :disabled="isCompulsory"
                  color="primary"
                />
                <v-text-field
                  v-if="formData.benefits.sgla_enabled"
                  v-model.number="formData.benefits.sgla_multiple"
                  label="Salary Multiple"
                  type="number"
                  variant="outlined"
                  density="compact"
                  class="ml-8"
                  step="0.5"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-checkbox
                  v-model="formData.benefits.ptd_enabled"
                  :label="getBenefitLabel('PTD')"
                  :disabled="isCompulsory"
                  color="primary"
                />
                <v-text-field
                  v-if="formData.benefits.ptd_enabled"
                  v-model.number="formData.benefits.ptd_multiple"
                  label="Salary Multiple"
                  type="number"
                  variant="outlined"
                  density="compact"
                  class="ml-8"
                  step="0.5"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-checkbox
                  v-model="formData.benefits.ci_enabled"
                  :label="getBenefitLabel('CI')"
                  :disabled="isCompulsory"
                  color="primary"
                />
                <v-text-field
                  v-if="formData.benefits.ci_enabled"
                  v-model.number="formData.benefits.ci_multiple"
                  label="Salary Multiple"
                  type="number"
                  variant="outlined"
                  density="compact"
                  class="ml-8"
                  step="0.5"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-checkbox
                  v-model="formData.benefits.ttd_enabled"
                  :label="getBenefitLabel('TTD')"
                  :disabled="isCompulsory"
                  color="primary"
                />
                <v-text-field
                  v-if="formData.benefits.ttd_enabled"
                  v-model.number="formData.benefits.ttd_multiple"
                  :label="getBenefitLabel('TTD')"
                  type="number"
                  variant="outlined"
                  density="compact"
                  class="ml-8"
                  step="0.5"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-checkbox
                  v-model="formData.benefits.phi_enabled"
                  :label="getBenefitLabel('PHI')"
                  :disabled="isCompulsory"
                  color="primary"
                />
                <v-text-field
                  v-if="formData.benefits.phi_enabled"
                  v-model.number="formData.benefits.phi_multiple"
                  label="Salary Multiple"
                  type="number"
                  variant="outlined"
                  density="compact"
                  class="ml-8"
                  step="0.5"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-checkbox
                  v-model="formData.benefits.gff_enabled"
                  :label="getBenefitLabel('GFF')"
                  :disabled="isCompulsory"
                  color="primary"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Action Buttons -->
      <v-col cols="12">
        <v-card-actions class="justify-end">
          <v-btn
            rounded
            size="small"
            color="grey"
            variant="outlined"
            @click="$emit('cancel')"
            >Cancel</v-btn
          >
          <v-btn
            rounded
            size="small"
            color="primary"
            type="submit"
            :loading="loading"
            :disabled="!isFormValid"
          >
            {{ isEditMode ? 'Update Member' : 'Add Member' }}
          </v-btn>
        </v-card-actions>
      </v-col>
    </v-row>
  </v-form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { VDateInput } from 'vuetify/labs/VDateInput'
import GroupPricingService from '@/renderer/api/GroupPricingService'

interface Member {
  id?: number
  member_name: string
  member_id_number: string
  member_id_type: string
  gender: string
  date_of_birth: Date | null
  email?: string
  phone_number?: string
  employee_number?: string
  scheme_id: number | null
  scheme_category: string | null
  entry_date: Date | null
  annual_salary: number | null
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

interface Props {
  member?: Member | null
  schemes: Array<any>
  isEditMode: boolean
  preselectedSchemeId?: number | null
}

interface Emits {
  (e: 'save', member: Member): void
  (e: 'cancel'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const form = ref(null)
const loading = ref(false)
const benefitMaps: any = ref([])

// Form data with default values
const formData = ref<Member>({
  member_name: '',
  member_id_number: '',
  member_id_type: 'RSA_ID',
  gender: '',
  date_of_birth: null,
  email: '',
  phone_number: '',
  employee_number: '',
  scheme_id: null,
  scheme_category: '',
  entry_date: null,
  annual_salary: null,
  address_line_1: '',
  address_line_2: '',
  city: '',
  province: '',
  postal_code: '',
  benefits: {
    gla_enabled: false,
    gla_multiple: 2,
    sgla_enabled: false,
    sgla_multiple: 1,
    ptd_enabled: false,
    ptd_multiple: 3,
    ci_enabled: false,
    ci_multiple: 2,
    ttd_enabled: false,
    ttd_multiple: 0.75,
    phi_enabled: false,
    phi_multiple: 0.75,
    gff_enabled: false
  }
})

// Form options
const idTypes = [
  { title: 'RSA ID', value: 'RSA_ID' },
  { title: 'Passport', value: 'PASSPORT' },
  { title: 'Other', value: 'OTHER' }
]

const genderOptions = [
  { title: 'Male', value: 'Male' },
  { title: 'Female', value: 'Female' }
]

const schemeCategories = ref()

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
  idNumber: (value: string) => {
    if (!value) return 'ID number is required'
    if (formData.value.member_id_type === 'RSA_ID') {
      return value.length === 13 || 'RSA ID number must be 13 digits'
    }
    return true
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
  },
  dateOfEntry: (value: Date | null) => {
    if (!value) return 'Entry is required'
    if (!selectedSchemeCommencementDate.value) return true
    const entryDate = new Date(value).setHours(0, 0, 0, 0)
    const commencementDate = new Date(
      selectedSchemeCommencementDate.value
    ).setHours(0, 0, 0, 0)
    return (
      entryDate >= commencementDate ||
      'Entry Date must be on or after scheme commencement date'
    )
  }
}

// Computed properties
const isFormValid = computed(() => {
  return (
    formData.value.member_name &&
    formData.value.member_id_number &&
    formData.value.gender &&
    formData.value.date_of_birth &&
    formData.value.scheme_id &&
    formData.value.entry_date &&
    formData.value.annual_salary &&
    formData.value.annual_salary > 0
  )
})

const idNumberLabel = computed(() => {
  switch (formData.value.member_id_type) {
    case 'RSA_ID':
      return 'RSA ID Number *'
    case 'PASSPORT':
      return 'Passport Number *'
    case 'OTHER':
      return 'ID Number *'
    default:
      return 'ID Number *'
  }
})

// Methods
const handleSubmit = async () => {
  if (!form.value) return

  const { valid } = await (form.value as any).validate()
  if (!valid) return

  loading.value = true
  try {
    // Prepare the data for submission
    const memberData = {
      ...formData.value,
      annual_salary: Number(formData.value.annual_salary),
      benefit_salary_multiple:
        Number(formData.value.benefits.gla_multiple) || 0,
      benefits: {
        ...formData.value.benefits,
        gla_multiple: Number(formData.value.benefits.gla_multiple) || 0,
        sgla_multiple: Number(formData.value.benefits.sgla_multiple) || 0,
        ptd_multiple: Number(formData.value.benefits.ptd_multiple) || 0,
        ci_multiple: Number(formData.value.benefits.ci_multiple) || 0,
        ttd_multiple: Number(formData.value.benefits.ttd_multiple) || 0,
        phi_multiple: Number(formData.value.benefits.phi_multiple) || 0
      }
    }

    emit('save', memberData)
  } finally {
    loading.value = false
  }
}
const selectedSchemeCommencementDate = ref('')
const isCompulsory = ref(false)
const minDate = ref('')

const getBenefitLabel = (benefitCode: string) => {
  const benefitMap = benefitMaps.value.find(
    (benefit: any) => benefit.benefit_code === benefitCode
  )

  return (
    benefitMap?.benefit_alias ||
    benefitMap?.benefit_code ||
    benefitMap?.benefit_alias_code ||
    benefitMap?.benefit_name
  )
}

const onSchemeCategoryChange = async (category: string) => {
  const selectedSchemeCategory = rawCategories.value.find(
    (cat) => cat.scheme_category === category
  )

  if (selectedSchemeCategory) {
    if (isCompulsory.value) {
      formData.value.benefits.gla_enabled = selectedSchemeCategory?.gla_benefit
      formData.value.benefits.sgla_enabled =
        selectedSchemeCategory?.sgla_benefit
      formData.value.benefits.ptd_enabled = selectedSchemeCategory?.ptd_benefit
      formData.value.benefits.ci_enabled = selectedSchemeCategory?.ci_benefit
      formData.value.benefits.ttd_enabled = selectedSchemeCategory?.ttd_benefit
      formData.value.benefits.phi_enabled = selectedSchemeCategory?.phi_benefit
      formData.value.benefits.gff_enabled =
        selectedSchemeCategory?.family_funeral_benefit
    }
  }
}

interface SchemeCategory {
  scheme_category: string
  gla_benefit: boolean
  sgla_benefit: boolean
  ptd_benefit: boolean
  ci_benefit: boolean
  ttd_benefit: boolean
  phi_benefit: boolean
  family_funeral_benefit: boolean
}

const selectedScheme = ref()
const rawCategories = ref<SchemeCategory[]>([])
const onSchemeChange = async (schemeId: number) => {
  formData.value.scheme_category = null

  console.log('Schemes:', props.schemes)
  selectedScheme.value = props.schemes.find((s) => s.id === schemeId)

  console.log('Selected Scheme:', selectedScheme)

  const response = await GroupPricingService.getSchemeCategories(
    selectedScheme.value.quote_id
  )
  schemeCategories.value = response.data.map((c) => c.scheme_category)
  rawCategories.value = response.data

  isCompulsory.value =
    selectedScheme.value.quote.obligation_type === 'Compulsory'

  if (selectedScheme.value) {
    selectedSchemeCommencementDate.value =
      selectedScheme.value.commencement_date
    const commencementDate = new Date(selectedSchemeCommencementDate.value)
    // commencementDate.setDate(commencementDate.getDate() - 1)
    minDate.value = commencementDate.toISOString()
    console.log('Commencement date:', selectedScheme.value.commencement_date)
    console.log('Minimum Date:', minDate.value)
    const categories: any = []
    selectedScheme.value.active_scheme_categories.forEach(
      (category: string) => {
        categories.push({ title: category, value: category })
      }
    )

    // formData.value.scheme_category = schemeCategories.value[0]
    schemeCategories.value = categories
  }
}

// Watchers
watch(
  () => props.member,
  (newMember) => {
    if (newMember && props.isEditMode) {
      // Populate form with existing member data
      formData.value = {
        ...formData.value,
        ...newMember,
        date_of_birth: newMember.date_of_birth
          ? new Date(newMember.date_of_birth)
          : null,
        entry_date: newMember.entry_date
          ? new Date(newMember.entry_date)
          : null,
        benefits: {
          ...formData.value.benefits,
          ...(newMember.benefits || {})
        }
      }
    } else if (!props.isEditMode) {
      // Reset form for new member
      formData.value = {
        member_name: '',
        member_id_number: '',
        member_id_type: 'RSA_ID',
        gender: '',
        date_of_birth: null,
        email: '',
        phone_number: '',
        employee_number: '',
        scheme_id: null,
        scheme_category: '',
        entry_date: null,
        annual_salary: null,
        address_line_1: '',
        address_line_2: '',
        city: '',
        province: '',
        postal_code: '',
        benefits: {
          gla_enabled: false,
          gla_multiple: 2,
          sgla_enabled: false,
          sgla_multiple: 1,
          ptd_enabled: false,
          ptd_multiple: 3,
          ci_enabled: false,
          ci_multiple: 2,
          ttd_enabled: false,
          ttd_multiple: 0.75,
          phi_enabled: false,
          phi_multiple: 0.75,
          gff_enabled: false
        }
      }
    }
  },
  { immediate: true }
)

onMounted(async () => {
  // Set default entry date to today
  if (!props.isEditMode && !formData.value.entry_date) {
    formData.value.entry_date = new Date()
  }

  // Set preselected scheme if provided
  if (props.preselectedSchemeId && !props.isEditMode) {
    formData.value.scheme_id = props.preselectedSchemeId
    // Trigger scheme change to update categories if needed
    onSchemeChange(props.preselectedSchemeId)
  }

  const res = await GroupPricingService.getBenefitMaps()
  console.log('Benefit Maps:', res.data)
  benefitMaps.value = res.data
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
