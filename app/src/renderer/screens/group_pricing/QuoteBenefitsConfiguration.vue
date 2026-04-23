<template>
  <base-card :show-actions="false">
    <template #header>
      <span class="headline"> Benefit Details </span>
    </template>
    <template #default>
      <div class="d-flex align-end mb-1">
        <div style="width: 340px; flex-shrink: 0">
          <v-select
            v-model="selectedCategory"
            :items="quote.selected_scheme_categories"
            placeholder="Select a Scheme Category"
            label="Scheme Category"
            variant="outlined"
            density="compact"
            clearable
            hide-details
            @update:model-value="updatedBenefitDisplay"
          ></v-select>
        </div>
        <v-spacer />
        <div class="d-flex align-end" style="gap: 8px">
          <v-btn
            size="small"
            variant="outlined"
            color="primary"
            :loading="isGenerating"
            :disabled="
              !props.resultSummaries || props.resultSummaries.length === 0
            "
            prepend-icon="mdi-file-word-outline"
            @click="downloadWord"
          >
            Download Word
          </v-btn>
          <v-btn
            size="small"
            variant="outlined"
            color="error"
            :loading="isGenerating"
            :disabled="
              !props.resultSummaries || props.resultSummaries.length === 0
            "
            prepend-icon="mdi-file-pdf-box"
            @click="downloadPdf"
          >
            Download PDF
          </v-btn>
        </div>
      </div>
      <div
        v-if="selectedCategoryDetails?.region"
        class="d-flex align-center mb-2 mt-2"
      >
        <v-icon size="small" class="mr-2 text-grey-darken-1"
          >mdi-map-marker</v-icon
        >
        <span class="text-grey-darken-1 text-body-2">Region:</span>
        <span class="font-weight-medium text-body-2 ml-2">{{
          selectedCategoryDetails.region
        }}</span>
      </div>
      <v-expansion-panels
        v-if="selectedCategoryDetails"
        variant="inset"
        class="mt-4"
      >
        <v-expansion-panel>
          <v-expansion-panel-title>
            {{ selectedCategoryDetails.gla_alias || 'GLA' }}
            <v-chip
              :color="selectedCategoryDetails.gla_benefit ? 'info' : 'grey'"
              size="small"
              class="ml-4"
              label
            >
              {{ selectedCategoryDetails.gla_benefit ? 'Active' : 'Inactive' }}
            </v-chip>
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <v-row>
              <v-col>
                <v-list density="compact">
                  <v-row>
                    <v-col cols="6" md="4">
                      <v-list-item
                        title="Benefit Type"
                        :subtitle="selectedCategoryDetails.gla_benefit_type"
                      ></v-list-item>
                    </v-col>
                    <v-col cols="6" md="4">
                      <v-list-item
                        :title="glaEducatorTitle"
                        :subtitle="selectedCategoryDetails.gla_educator_benefit"
                      ></v-list-item>
                    </v-col>
                    <v-col cols="6" md="4">
                      <v-list-item
                        title="Salary Multiple"
                        :subtitle="glaSalaryMultiple"
                      ></v-list-item>
                    </v-col>
                    <v-col cols="6" md="4">
                      <v-list-item
                        title="Terminal Illness Benefit"
                        :subtitle="
                          selectedCategoryDetails.gla_terminal_illness_benefit
                        "
                      ></v-list-item>
                    </v-col>
                    <v-col cols="6" md="4">
                      <v-list-item
                        title="Waiting Period"
                        :subtitle="selectedCategoryDetails.gla_waiting_period"
                      ></v-list-item>
                    </v-col>
                    <v-col
                      v-if="selectedCategoryDetails.tax_saver_benefit"
                      cols="6"
                      md="4"
                    >
                      <v-list-item
                        title="Tax Saver"
                        subtitle="Enabled"
                      ></v-list-item>
                    </v-col>
                  </v-row>
                </v-list>
              </v-col> </v-row
          ></v-expansion-panel-text>
        </v-expansion-panel>

        <v-expansion-panel>
          <v-expansion-panel-title>
            {{ selectedCategoryDetails.sgla_alias || 'SGLA' }}
            <v-chip
              :color="selectedCategoryDetails.sgla_benefit ? 'info' : 'grey'"
              size="small"
              class="ml-4"
              label
            >
              {{ selectedCategoryDetails.sgla_benefit ? 'Active' : 'Inactive' }}
            </v-chip>
          </v-expansion-panel-title>

          <v-expansion-panel-text>
            <v-row>
              <v-col>
                <v-list density="compact">
                  <v-row>
                    <v-col cols="6" md="4">
                      <v-list-item
                        title="Salary Multiple"
                        :subtitle="sglaSalaryMultiple"
                      ></v-list-item>
                    </v-col>
                  </v-row>
                </v-list>
              </v-col>
            </v-row>
          </v-expansion-panel-text>
        </v-expansion-panel>
        <v-expansion-panel>
          <v-expansion-panel-title>
            {{ selectedCategoryDetails.ptd_alias || 'PTD' }}
            <v-chip
              :color="selectedCategoryDetails.ptd_benefit ? 'info' : 'grey'"
              size="small"
              class="ml-4"
              label
            >
              {{ selectedCategoryDetails.ptd_benefit ? 'Active' : 'Inactive' }}
            </v-chip>
          </v-expansion-panel-title>

          <v-expansion-panel-text>
            <v-row>
              <v-col>
                <v-list density="compact">
                  <v-row>
                    <v-col cols="6" md="4">
                      <v-list-item
                        title="Benefit Type"
                        :subtitle="selectedCategoryDetails.ptd_benefit_type"
                      ></v-list-item>
                    </v-col>
                    <v-col cols="6" md="4">
                      <v-list-item
                        title="Deferred Period"
                        :subtitle="selectedCategoryDetails.ptd_deferred_period"
                      ></v-list-item>
                    </v-col>
                    <v-col cols="6" md="4">
                      <v-list-item
                        title="Disability Definition"
                        :subtitle="
                          selectedCategoryDetails.ptd_disability_definition
                        "
                      ></v-list-item>
                    </v-col>
                    <v-col cols="6" md="4">
                      <v-list-item
                        :title="ptdEducatorTitle"
                        :subtitle="selectedCategoryDetails.ptd_educator_benefit"
                      ></v-list-item>
                    </v-col>
                    <v-col cols="6" md="4">
                      <v-list-item
                        title="Risk Type"
                        :subtitle="selectedCategoryDetails.ptd_risk_type"
                      ></v-list-item>
                    </v-col>
                    <v-col cols="6" md="4">
                      <v-list-item
                        title="Salary Multiple"
                        :subtitle="ptdSalaryMultiple"
                      ></v-list-item>
                    </v-col>
                  </v-row>
                </v-list>
              </v-col>
            </v-row>
          </v-expansion-panel-text>
        </v-expansion-panel>
        <v-expansion-panel>
          <v-expansion-panel-title>
            {{ selectedCategoryDetails.ci_alias || 'CI' }}
            <v-chip
              :color="selectedCategoryDetails.ci_benefit ? 'info' : 'grey'"
              size="small"
              class="ml-4"
              label
            >
              {{ selectedCategoryDetails.ci_benefit ? 'Active' : 'Inactive' }}
            </v-chip>
          </v-expansion-panel-title>

          <v-expansion-panel-text>
            <v-row>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Benefit Definition"
                  :subtitle="selectedCategoryDetails.ci_benefit_definition"
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Benefit Structure"
                  :subtitle="selectedCategoryDetails.ci_benefit_structure"
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Salary Multiple"
                  :subtitle="ciCriticalIllnessSalaryMultiple"
                ></v-list-item>
              </v-col>
            </v-row>
          </v-expansion-panel-text>
        </v-expansion-panel>
        <v-expansion-panel>
          <v-expansion-panel-title>
            {{ selectedCategoryDetails.ttd_alias || 'TTD' }}
            <v-chip
              :color="selectedCategoryDetails.ttd_benefit ? 'info' : 'grey'"
              size="small"
              class="ml-4"
              label
            >
              {{ selectedCategoryDetails.ttd_benefit ? 'Active' : 'Inactive' }}
            </v-chip>
          </v-expansion-panel-title>

          <v-expansion-panel-text>
            <v-row>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Deferred Period"
                  :subtitle="selectedCategoryDetails.ttd_deferred_period"
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Disability Definition"
                  :subtitle="selectedCategoryDetails.ttd_disability_definition"
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Income Replacement Percentage"
                  :subtitle="
                    selectedCategoryDetails.ttd_income_replacement_percentage
                  "
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Tiered Income Replacement Ratio"
                  :subtitle="
                    selectedCategoryDetails.ttd_use_tiered_income_replacement_ratio
                      ? `Yes (${selectedCategoryDetails.ttd_tiered_income_replacement_type === 'custom' ? 'Custom' : 'Standard'})`
                      : 'No'
                  "
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Premium Waiver Percentage"
                  :subtitle="
                    selectedCategoryDetails.ttd_premium_waiver_percentage
                  "
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Risk Type"
                  :subtitle="selectedCategoryDetails.ttd_risk_type"
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Waiting Period"
                  :subtitle="selectedCategoryDetails.ttd_waiting_period"
                ></v-list-item>
              </v-col> </v-row
          ></v-expansion-panel-text>
        </v-expansion-panel>
        <v-expansion-panel>
          <v-expansion-panel-title>
            {{ selectedCategoryDetails.phi_alias || 'PHI' }}
            <v-chip
              :color="selectedCategoryDetails.phi_benefit ? 'info' : 'grey'"
              size="small"
              class="ml-4"
              label
            >
              {{ selectedCategoryDetails.phi_benefit ? 'Active' : 'Inactive' }}
            </v-chip>
          </v-expansion-panel-title>

          <v-expansion-panel-text
            ><v-row>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Benefit Escalation Option"
                  :subtitle="selectedCategoryDetails.phi_benefit_escalation"
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Deferred Period"
                  :subtitle="selectedCategoryDetails.phi_deferred_period"
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Disability Definition"
                  :subtitle="selectedCategoryDetails.phi_disability_definition"
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Income Replacement Percentage"
                  :subtitle="
                    selectedCategoryDetails.phi_income_replacement_percentage
                  "
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Tiered Income Replacement Ratio"
                  :subtitle="
                    selectedCategoryDetails.phi_use_tiered_income_replacement_ratio
                      ? `Yes (${selectedCategoryDetails.phi_tiered_income_replacement_type === 'custom' ? 'Custom' : 'Standard'})`
                      : 'No'
                  "
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Medical Aid Premium Waiver"
                  :subtitle="
                    selectedCategoryDetails.phi_medical_aid_premium_waiver
                  "
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Normal Retirement Age"
                  :subtitle="selectedCategoryDetails.phi_normal_retirement_age"
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Premium Waiver"
                  :subtitle="selectedCategoryDetails.phi_premium_waiver"
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Risk Type"
                  :subtitle="selectedCategoryDetails.phi_risk_type"
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Waiting Period"
                  :subtitle="selectedCategoryDetails.phi_waiting_period"
                ></v-list-item>
              </v-col>
            </v-row>
          </v-expansion-panel-text>
        </v-expansion-panel>
        <v-expansion-panel>
          <v-expansion-panel-title>
            {{ selectedCategoryDetails.family_funeral_alias || 'GFF' }}
            <v-chip
              :color="
                selectedCategoryDetails.family_funeral_benefit ? 'info' : 'grey'
              "
              size="small"
              class="ml-4"
              label
            >
              {{
                selectedCategoryDetails.family_funeral_benefit
                  ? 'Active'
                  : 'Inactive'
              }}
            </v-chip>
          </v-expansion-panel-title>

          <v-expansion-panel-text
            ><v-row>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Main Member Sum Assured"
                  :subtitle="
                    formatNumber(
                      selectedCategoryDetails.family_funeral_main_member_funeral_sum_assured
                    )
                  "
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Spouse Sum Assured"
                  :subtitle="
                    formatNumber(
                      selectedCategoryDetails.family_funeral_spouse_funeral_sum_assured
                    )
                  "
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Child Sum Assured"
                  :subtitle="
                    formatNumber(
                      selectedCategoryDetails.family_funeral_children_funeral_sum_assured
                    )
                  "
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Parent Sum Assured"
                  :subtitle="
                    formatNumber(
                      selectedCategoryDetails.family_funeral_parent_funeral_sum_assured
                    )
                  "
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Adult Dependent Sum Assured"
                  :subtitle="
                    formatNumber(
                      selectedCategoryDetails.family_funeral_adult_dependant_sum_assured
                    )
                  "
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Maximum Number of Children"
                  :subtitle="
                    selectedCategoryDetails.family_funeral_max_number_children
                  "
                ></v-list-item>
              </v-col>
              <v-col cols="6" md="4">
                <v-list-item
                  title="Maximum Number of Adult Dependents"
                  :subtitle="
                    selectedCategoryDetails.family_funeral_max_number_adult_dependants
                  "
                ></v-list-item>
              </v-col>
            </v-row>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
    </template>
  </base-card>
</template>

<script setup lang="ts">
import BaseCard from '@/renderer/components/BaseCard.vue'
import { ref, computed } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useBenefitScheduleGeneration } from '@/renderer/composables/useBenefitScheduleGeneration'

const {
  isGenerating,
  generateBenefitScheduleDocx,
  generateBenefitSchedulePdf
} = useBenefitScheduleGeneration()

const selectedCategory = ref(null)
const selectedCategoryDetails: any = ref(null)

const glaSalaryMultiple = ref('')
const sglaSalaryMultiple = ref('')
const ptdSalaryMultiple = ref('')
const ciCriticalIllnessSalaryMultiple = ref('')

const benefitMaps = ref([])
const glaBenefitTitle = ref('')
const sglaBenefitTitle = ref('')
const ptdBenefitTitle = ref('')
const ciBenefitTitle = ref('')
const phiBenefitTitle = ref('')
const ttdBenefitTitle = ref('')
const familyFuneralBenefitTitle = ref('')
const glaEducatorTitle = ref('GLA Educator')
const ptdEducatorTitle = ref('PTD Educator')

const props = defineProps({
  quote: {
    type: Object,
    required: true
  },
  resultSummaries: {
    type: Array,
    default: () => []
  }
})

const benefitTitles = computed(() => ({
  glaBenefitTitle: glaBenefitTitle.value,
  sglaBenefitTitle: sglaBenefitTitle.value,
  ptdBenefitTitle: ptdBenefitTitle.value,
  ciBenefitTitle: ciBenefitTitle.value,
  phiBenefitTitle: phiBenefitTitle.value,
  ttdBenefitTitle: ttdBenefitTitle.value,
  familyFuneralBenefitTitle: familyFuneralBenefitTitle.value
}))

const downloadWord = () => {
  generateBenefitScheduleDocx(
    props.quote,
    props.resultSummaries as any[],
    benefitTitles.value
  )
}

const downloadPdf = () => {
  generateBenefitSchedulePdf(
    props.quote,
    props.resultSummaries as any[],
    benefitTitles.value
  )
}

GroupPricingService.getBenefitMaps().then((res: any) => {
  benefitMaps.value = res.data
  const resolve = (code: string) => {
    const b: any = benefitMaps.value.find((m: any) => m.benefit_code === code)
    return b ? b.benefit_alias || b.benefit_name : code
  }
  glaBenefitTitle.value = resolve('GLA')
  sglaBenefitTitle.value = resolve('SGLA')
  ptdBenefitTitle.value = resolve('PTD')
  ciBenefitTitle.value = resolve('CI')
  phiBenefitTitle.value = resolve('PHI')
  ttdBenefitTitle.value = resolve('TTD')
  familyFuneralBenefitTitle.value = resolve('GFF')
  glaEducatorTitle.value = resolve('GLA_EDU')
  ptdEducatorTitle.value = resolve('PTD_EDU')
})

const formatNumber = (value: any) => {
  if (value === null || value === undefined || value === '') {
    return ''
  }
  const num = Number(value)
  if (isNaN(num)) {
    return value
  }
  return num.toLocaleString()
}

const updatedBenefitDisplay = () => {
  // Logic to update the display based on selected category
  selectedCategoryDetails.value = props.quote.scheme_categories.find(
    (category: any) => category.scheme_category === selectedCategory.value
  )
  console.log('Selected Category Details:', selectedCategoryDetails.value)

  if (props.quote.use_global_salary_multiple === false) {
    glaSalaryMultiple.value =
      sglaSalaryMultiple.value =
      ptdSalaryMultiple.value =
      ciCriticalIllnessSalaryMultiple.value =
        'This is read from Member Data.'
  } else {
    // :subtitle="selectedCategoryDetails.gla_salary_multiple"
    glaSalaryMultiple.value = selectedCategoryDetails.value.gla_salary_multiple
    sglaSalaryMultiple.value =
      selectedCategoryDetails.value.sgla_salary_multiple
    ptdSalaryMultiple.value = selectedCategoryDetails.value.ptd_salary_multiple
    ciCriticalIllnessSalaryMultiple.value =
      selectedCategoryDetails.value.ci_critical_illness_salary_multiple
  }
}

// You can still have the logic to get dynamic benefit titles here if needed
</script>
