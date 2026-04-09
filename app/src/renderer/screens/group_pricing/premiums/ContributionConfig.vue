<template>
  <base-card :show-actions="false">
    <template #header>
      <h3 class="mb-0">Contribution Configuration</h3>
    </template>
    <template #default>
      <!-- Loading State -->
      <v-row v-if="loading">
        <v-col class="text-center py-8">
          <v-progress-circular indeterminate color="primary" />
          <div class="mt-2 text-body-2 text-grey">Loading configuration...</div>
        </v-col>
      </v-row>

      <v-form v-else ref="form">
        <!-- Contribution Mode Toggle (AC: 2) -->
        <v-row>
          <v-col cols="12">
            <div class="text-subtitle-1 font-weight-medium mb-3"
              >Contribution Mode</div
            >
            <v-btn-toggle
              v-model="contributionMode"
              mandatory
              color="primary"
              variant="outlined"
              density="compact"
            >
              <v-btn value="employer_only">Employer Only</v-btn>
              <v-btn value="split">Split</v-btn>
            </v-btn-toggle>
          </v-col>
        </v-row>

        <!-- Employer / Employee % fields — shown only in Split mode (AC: 2, 3) -->
        <v-row v-if="contributionMode === 'split'" class="mt-4">
          <v-col cols="12" md="4">
            <v-text-field
              v-model.number="employerPct"
              type="number"
              label="Employer %"
              variant="outlined"
              density="compact"
              :rules="[rules.required, rules.percentage]"
              :error-messages="apiErrors.employerPct"
              min="0"
              max="100"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              v-model.number="employeePct"
              type="number"
              label="Employee %"
              variant="outlined"
              density="compact"
              :rules="[rules.required, rules.percentage]"
              :error-messages="splitSumErrors"
              min="0"
              max="100"
            />
          </v-col>
          <v-col cols="12" md="4" class="d-flex align-center">
            <v-chip
              :color="pctSumValid ? 'success' : 'error'"
              variant="tonal"
              size="small"
            >
              Total: {{ employerPct + employeePct }}%
            </v-chip>
          </v-col>
        </v-row>

        <!-- Effective Date (AC: 5) -->
        <v-row class="mt-4">
          <v-col cols="12" md="4">
            <v-date-input
              v-model="effectiveDate"
              label="Effective Date *"
              variant="outlined"
              density="compact"
              hide-actions
              prepend-icon=""
              prepend-inner-icon="$calendar"
              locale="en-ZA"
              view-mode="month"
              :rules="[rules.requiredDate]"
              :error-messages="apiErrors.effectiveDate"
            />
          </v-col>
        </v-row>

        <!-- Per-Benefit Override Accordion (AC: 4) -->
        <v-row class="mt-6">
          <v-col>
            <div class="text-subtitle-1 font-weight-medium mb-3"
              >Per-Benefit Overrides</div
            >
            <v-expansion-panels
              v-if="benefitOverrides.length > 0"
              multiple
              variant="accordion"
            >
              <v-expansion-panel
                v-for="(override, idx) in benefitOverrides"
                :key="override.benefitId"
              >
                <v-expansion-panel-title>
                  <div class="d-flex align-center gap-3">
                    <span>{{ override.benefitName }}</span>
                    <v-chip
                      v-if="override.overrideEnabled"
                      size="small"
                      color="primary"
                      variant="tonal"
                    >
                      Override Active
                    </v-chip>
                  </div>
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <v-row>
                    <v-col cols="12">
                      <v-switch
                        v-model="override.overrideEnabled"
                        label="Enable override for this benefit"
                        color="primary"
                        density="compact"
                        hide-details
                      />
                    </v-col>
                  </v-row>
                  <v-row v-if="override.overrideEnabled" class="mt-3">
                    <v-col cols="12" md="4">
                      <v-text-field
                        v-model.number="override.employerPct"
                        type="number"
                        label="Employer %"
                        variant="outlined"
                        density="compact"
                        :rules="[rules.required, rules.percentage]"
                        min="0"
                        max="100"
                      />
                    </v-col>
                    <v-col cols="12" md="4">
                      <v-text-field
                        v-model.number="override.employeePct"
                        type="number"
                        label="Employee %"
                        variant="outlined"
                        density="compact"
                        :rules="[rules.required, rules.percentage]"
                        :error-messages="benefitSumError(idx)"
                        min="0"
                        max="100"
                      />
                    </v-col>
                    <v-col cols="12" md="4" class="d-flex align-center">
                      <v-chip
                        :color="
                          override.employerPct + override.employeePct === 100
                            ? 'success'
                            : 'error'
                        "
                        variant="tonal"
                        size="small"
                      >
                        Total:
                        {{ override.employerPct + override.employeePct }}%
                      </v-chip>
                    </v-col>
                  </v-row>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
            <div v-else class="text-body-2 text-grey">
              No benefits available for override configuration.
            </div>
          </v-col>
        </v-row>

        <!-- Save Button (AC: 3, 6) -->
        <v-row class="mt-6">
          <v-col>
            <v-btn
              color="primary"
              :disabled="!canSave"
              :loading="saving"
              @click="onSave"
            >
              Save Configuration
            </v-btn>
          </v-col>
        </v-row>
      </v-form>
    </template>
  </base-card>

  <!-- Snackbar (AC: 6) -->
  <v-snackbar
    v-model="snackbar"
    :timeout="3000"
    :color="snackbarColor"
    centered
  >
    {{ snackbarText }}
    <template #actions>
      <v-btn variant="text" @click="snackbar = false">Close</v-btn>
    </template>
  </v-snackbar>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { VDateInput } from 'vuetify/labs/VDateInput'
import PremiumManagementService from '@/renderer/api/PremiumManagementService'
import BaseCard from '@/renderer/components/BaseCard.vue'

const props = defineProps<{
  schemeId: number
}>()

// ─── Types ──────────────────────────────────────────────────────────────────

interface BenefitOverride {
  benefitId: string
  benefitName: string
  overrideEnabled: boolean
  employerPct: number
  employeePct: number
}

interface ApiFieldErrors {
  contributionMode?: string[]
  employerPct?: string[]
  employeePct?: string[]
  effectiveDate?: string[]
}

// ─── Form State ──────────────────────────────────────────────────────────────

const loading = ref(false)
const saving = ref(false)
const form = ref()

const contributionMode = ref<'employer_only' | 'split'>('employer_only')
const employerPct = ref<number>(100)
const employeePct = ref<number>(0)
const effectiveDate = ref<Date | null>(null)
const benefitOverrides = ref<BenefitOverride[]>([])
const apiErrors = ref<ApiFieldErrors>({})

// ─── Snackbar ────────────────────────────────────────────────────────────────

const snackbar = ref(false)
const snackbarText = ref('')
const snackbarColor = ref<string>('success')

// ─── Validation Rules ────────────────────────────────────────────────────────

const rules = {
  required: (v: any) =>
    (v !== null && v !== undefined && v !== '') || 'This field is required',
  requiredDate: (v: any) => !!v || 'Effective date is required',
  percentage: (v: number) => (v >= 0 && v <= 100) || 'Must be between 0 and 100'
}

// ─── Computed Validation ─────────────────────────────────────────────────────

const pctSumValid = computed(() => {
  if (contributionMode.value === 'employer_only') return true
  return employerPct.value + employeePct.value === 100
})

const splitSumErrors = computed<string[]>(() => {
  if (contributionMode.value === 'split' && !pctSumValid.value) {
    return ['Employer % + Employee % must equal 100']
  }
  if (apiErrors.value.employeePct?.length) {
    return apiErrors.value.employeePct
  }
  return []
})

const benefitSumError = (idx: number): string[] => {
  const o = benefitOverrides.value[idx]
  if (o.overrideEnabled && o.employerPct + o.employeePct !== 100) {
    return ['Must equal 100']
  }
  return []
}

const benefitOverridesValid = computed(() => {
  return benefitOverrides.value.every((o) => {
    if (!o.overrideEnabled) return true
    return o.employerPct + o.employeePct === 100
  })
})

const canSave = computed(() => {
  return (
    pctSumValid.value && benefitOverridesValid.value && !!effectiveDate.value
  )
})

// ─── Date Helpers ─────────────────────────────────────────────────────────────

function parseDateStringToDate(dateString: string | null): Date | null {
  if (!dateString) return null
  const [year, month, day] = dateString.split('-').map(Number)
  return new Date(year, month - 1, day)
}

function formatDateToYMD(dateObj: Date | null): string {
  if (!dateObj) return ''
  const year = dateObj.getFullYear()
  const month = String(dateObj.getMonth() + 1).padStart(2, '0')
  const day = String(dateObj.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// ─── API Methods ──────────────────────────────────────────────────────────────

const loadConfig = async () => {
  loading.value = true
  try {
    const res = await PremiumManagementService.getContributionConfig(
      props.schemeId
    )
    const data = res.data?.data ?? res.data
    contributionMode.value = data.contribution_type || 'employer_only'
    employerPct.value = data.employer_percent ?? 100
    employeePct.value = data.employee_percent ?? 0
    effectiveDate.value = parseDateStringToDate(data.effective_date ?? null)
    benefitOverrides.value = (data.benefit_overrides || []).map((b: any) => ({
      benefitId: b.benefit_code,
      benefitName: b.benefit_code,
      overrideEnabled: !!(b.employer_percent || b.employee_percent),
      employerPct: b.employer_percent ?? 0,
      employeePct: b.employee_percent ?? 0
    }))
  } catch (err) {
    console.error('Failed to load contribution config:', err)
  } finally {
    loading.value = false
  }
}

// ─── Save Handler (AC: 6, 7) ─────────────────────────────────────────────────

const onSave = async () => {
  const result = await form.value?.validate()
  if (!result?.valid) return
  if (!canSave.value) return

  saving.value = true
  apiErrors.value = {}

  const payload = {
    scheme_id: props.schemeId,
    contribution_type: contributionMode.value,
    employer_percent:
      contributionMode.value === 'split' ? employerPct.value : 100,
    employee_percent:
      contributionMode.value === 'split' ? employeePct.value : 0,
    effective_date: formatDateToYMD(effectiveDate.value),
    benefit_overrides: benefitOverrides.value.map((o) => ({
      benefit_code: o.benefitId,
      employer_percent: o.employerPct,
      employee_percent: o.employeePct
    }))
  }

  try {
    await PremiumManagementService.saveContributionConfig(
      props.schemeId,
      payload
    )
    snackbarText.value = 'Contribution configuration saved successfully'
    snackbarColor.value = 'success'
    snackbar.value = true
    await loadConfig()
  } catch (err: any) {
    if (err?.response?.data?.errors) {
      apiErrors.value = err.response.data.errors
      snackbarText.value = 'Please correct the highlighted errors and try again'
    } else {
      snackbarText.value = 'Failed to save configuration'
    }
    snackbarColor.value = 'error'
    snackbar.value = true
  } finally {
    saving.value = false
  }
}

// ─── Lifecycle ────────────────────────────────────────────────────────────────

onMounted(() => {
  loadConfig()
})
</script>
