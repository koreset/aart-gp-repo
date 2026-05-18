<template>
  <v-card variant="outlined" rounded="lg">
    <v-card-title class="d-flex align-center font-weight-bold">
      <span>Medical disclosure</span>
      <v-spacer />
      <v-btn-toggle
        v-model="form.form_variant"
        density="compact"
        variant="outlined"
        divided
        mandatory
      >
        <v-btn value="short">Short-form</v-btn>
        <v-btn value="long">Tele-UW long-form</v-btn>
      </v-btn-toggle>
    </v-card-title>
    <v-card-text>
      <v-alert
        v-if="latest"
        type="info"
        variant="tonal"
        density="compact"
        class="mb-3"
        icon="mdi-history"
      >
        Last submitted by {{ latest.submitted_by || '—' }} on
        {{ formatDate(latest.submitted_at) }} ({{
          latest.form_variant || 'short'
        }}). BMI: <strong>{{ latest.bmi?.toFixed(1) || '—' }}</strong>
      </v-alert>

      <v-row dense>
        <v-col cols="12" md="6">
          <p class="text-subtitle-2 font-weight-bold mb-2">Build</p>
          <v-row dense>
            <v-col cols="6">
              <v-text-field
                v-model.number="form.height"
                label="Height (cm)"
                type="number"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model.number="form.weight"
                label="Weight (kg)"
                type="number"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12">
              <p class="text-caption text-grey"
                >Computed BMI:
                <strong>{{ computedBMI || '—' }}</strong></p
              >
            </v-col>
          </v-row>
        </v-col>

        <v-col cols="12" md="6">
          <p class="text-subtitle-2 font-weight-bold mb-2">Lifestyle</p>
          <v-checkbox
            v-model="form.smoker"
            label="Smoker"
            density="compact"
            hide-details
          />
          <v-text-field
            v-if="form.smoker"
            v-model.number="form.cigarettes_per_day"
            label="Cigarettes per day"
            type="number"
            density="compact"
            variant="outlined"
            class="mt-2"
          />
          <v-text-field
            v-model.number="form.alcohol_units_per_week"
            label="Alcohol units per week"
            type="number"
            density="compact"
            variant="outlined"
            class="mt-2"
          />
          <v-checkbox
            v-model="form.has_hazardous_hobbies"
            label="Hazardous hobbies"
            density="compact"
            hide-details
          />
          <v-text-field
            v-if="form.has_hazardous_hobbies"
            v-model="form.hazardous_hobbies"
            label="Describe (e.g. scuba, climbing)"
            density="compact"
            variant="outlined"
            class="mt-2"
          />
        </v-col>

        <v-col cols="12">
          <p class="text-subtitle-2 font-weight-bold mb-2 mt-3"
            >Disclosed conditions</p
          >
          <v-combobox
            v-model="form.disclosed_conditions"
            label="Disclosed conditions (type or pick)"
            chips
            multiple
            clearable
            :items="conditionSuggestions"
            density="compact"
            variant="outlined"
            persistent-hint
            hint="Free text accepted; preferred codes appear in the dropdown. Phase 6 will validate against the UWConditionCode catalogue."
          />
        </v-col>

        <v-col v-if="form.form_variant === 'long'" cols="12">
          <p class="text-subtitle-2 font-weight-bold mb-2 mt-3"
            >Occupation risk questions</p
          >
          <v-row dense>
            <v-col
              v-for="(q, key) in occupationQuestions"
              :key="key"
              cols="12"
              md="6"
            >
              <v-checkbox
                v-model="form.occupation_risk_answers[key]"
                :label="q"
                density="compact"
                hide-details
              />
            </v-col>
          </v-row>
        </v-col>

        <v-col cols="12">
          <v-textarea
            v-model="form.additional_notes"
            label="Additional notes"
            rows="2"
            density="compact"
            variant="outlined"
            class="mt-3"
          />
        </v-col>
      </v-row>

      <v-divider class="my-3" />
      <p class="text-subtitle-2 font-weight-bold mb-2"
        >POPIA consent for medical information</p
      >
      <v-checkbox
        v-model="consentGranted"
        label="The member consents to the processing of medical information for underwriting."
        density="compact"
        hide-details
      />
      <v-text-field
        v-model="consentGrantedByName"
        label="Signed by (typed name)"
        density="compact"
        variant="outlined"
        class="mt-2"
        :disabled="!consentGranted"
      />

      <v-alert
        v-if="error"
        type="error"
        variant="tonal"
        density="compact"
        class="mt-3"
        >{{ error }}</v-alert
      >
      <v-btn
        class="mt-3"
        color="primary"
        :loading="busy"
        :disabled="!canSubmit"
        @click="submit"
        >Save disclosure</v-btn
      >
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

const props = defineProps<{
  caseId: number
}>()
const emit = defineEmits<{
  (e: 'submitted'): void
}>()

interface DisclosureRow {
  id: number
  bmi: number
  submitted_by: string
  submitted_at: string
  form_variant: string
}

const latest = ref<DisclosureRow | null>(null)
const consentGranted = ref(false)
const consentGrantedByName = ref('')
const busy = ref(false)
const error = ref('')

const form = ref({
  height: 0,
  weight: 0,
  smoker: false,
  cigarettes_per_day: 0,
  alcohol_units_per_week: 0,
  has_hazardous_hobbies: false,
  hazardous_hobbies: '',
  occupation_risk_answers: {} as Record<string, any>,
  disclosed_conditions: [] as string[],
  additional_notes: '',
  form_variant: 'short' as 'short' | 'long'
})

const formatDate = (s: string) => (s ? new Date(s).toLocaleString() : '—')

const conditionSuggestions = [
  'diabetes_type1',
  'diabetes_type2',
  'hypertension',
  'asthma',
  'cancer_history',
  'cardiovascular_disease',
  'kidney_disease',
  'mental_health_treatment'
]

const occupationQuestions: Record<string, string> = {
  works_at_height: 'Works at height above 3m regularly?',
  works_underground: 'Works underground (mining)?',
  handles_hazardous_materials: 'Handles hazardous materials?',
  operates_heavy_machinery: 'Operates heavy machinery?',
  drives_for_a_living: 'Drives more than 30,000 km/year for work?'
}

const computedBMI = computed(() => {
  if (form.value.height <= 0 || form.value.weight <= 0) return ''
  const m = form.value.height / 100
  return (form.value.weight / (m * m)).toFixed(1)
})

const canSubmit = computed(
  () =>
    consentGranted.value &&
    consentGrantedByName.value.trim().length > 0 &&
    form.value.height > 0 &&
    form.value.weight > 0
)

const load = async () => {
  try {
    const res = await GroupPricingService.getMemberDisclosure(props.caseId)
    latest.value = res.data
  } catch (err: any) {
    if (err?.response?.status !== 404) console.warn(err)
  }
}

const submit = async () => {
  error.value = ''
  busy.value = true
  try {
    // Record consent first so the POPIA gate on the backend passes.
    await GroupPricingService.submitConsent({
      case_id: props.caseId,
      consent_type: 'medical_info',
      granted_by_name: consentGrantedByName.value.trim()
    })
    await GroupPricingService.submitMemberDisclosure(props.caseId, form.value)
    emit('submitted')
    await load()
  } catch (err: any) {
    error.value =
      err?.response?.data || err?.message || 'Failed to save disclosure'
  } finally {
    busy.value = false
  }
}

watch(
  () => props.caseId,
  () => load(),
  { immediate: false }
)

onMounted(load)
</script>
