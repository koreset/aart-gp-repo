<template>
  <v-container>
    <!-- Add scheme type selector -->
    <v-row>
      <v-col cols="4">
        <v-select
          v-model="selectedSchemeType"
          variant="outlined"
          density="compact"
          :items="groupStore.group_pricing_quote.selected_scheme_categories"
          label="Scheme Category"
          placeholder="Select a scheme category"
          @update:model-value="onSchemeTypeChange"
        />
      </v-col>
      <v-col cols="4">
        <v-select
          v-model="selectedRegion"
          variant="outlined"
          density="compact"
          :items="availableRegions"
          label="Region"
          placeholder="Select a region"
          clearable
          prepend-inner-icon="mdi-map-marker"
          :disabled="!selectedSchemeType || availableRegions.length === 0"
          :hint="
            availableRegions.length === 0
              ? 'No regions for this risk rate code'
              : ''
          "
          persistent-hint
        />
      </v-col>
      <v-col cols="4">
        <v-btn
          size="small"
          rounded
          color="primary"
          :disabled="!selectedSchemeType"
          @click="saveCurrentSchemeCategory"
        >
          Save Scheme Category
        </v-btn>
      </v-col>
    </v-row>

    <!-- Snackbar for save confirmation -->
    <v-snackbar v-model="snackbar.show" :timeout="4000" color="success">
      {{ snackbar.message }}
    </v-snackbar>

    <v-tabs v-model="currentTab" color="primary" class="mb-5">
      <v-tab value="gla" :class="{ 'error-tab': tabStatus.gla }">
        {{ glaLabel }}
      </v-tab>
      <v-tab value="ptd" :class="{ 'error-tab': tabStatus.ptd }">{{
        ptdLabel
      }}</v-tab>
      <v-tab value="ci" :class="{ 'error-tab': tabStatus.ci }">{{
        ciLabel
      }}</v-tab>
      <v-tab value="sgla" :class="{ 'error-tab': tabStatus.sgla }">{{
        sglaLabel
      }}</v-tab>
      <v-tab value="phi" :class="{ 'error-tab': tabStatus.phi }">{{
        phiLabel
      }}</v-tab>
      <v-tab value="ttd" :class="{ 'error-tab': tabStatus.ttd }">{{
        ttdLabel
      }}</v-tab>
      <v-tab
        value="familyFuneral"
        :class="{ 'error-tab': tabStatus.family_funeral }"
        >{{ familyFuneralLabel }}</v-tab
      >
    </v-tabs>

    <v-window v-model="currentTab">
      <v-window-item class="pa-2" value="gla">
        <base-card :show-actions="false">
          <template #header>
            <span class="headline">{{ glaLabel }} Benefit</span>
          </template>
          <template #default>
            <v-switch
              v-model="glaBenefit"
              color="primary"
              :label="`Enable ${glaLabel} Benefit`"
              :disabled="!selectedSchemeType"
            ></v-switch>
            <v-divider class="mb-5"></v-divider>
            <v-row>
              <v-col cols="4">
                <v-select
                  v-model="glaBenefitType"
                  v-bind="glaBenefitTypeAttrs"
                  variant="outlined"
                  density="compact"
                  label="Benefit Type"
                  placeholder="Select Benefit Type"
                  :error-messages="errors.gla_benefit_type"
                  :items="glaBenefitTypes"
                  :disabled="!glaBenefit || !selectedSchemeType"
                ></v-select>
              </v-col>
              <v-col
                v-if="groupStore.group_pricing_quote.use_global_salary_multiple"
                cols="4"
              >
                <v-text-field
                  v-model:model-value="glaSalaryMultiple"
                  v-bind="glaSalaryMultipleAttrs"
                  type="number"
                  variant="outlined"
                  density="compact"
                  placeholder="Enter a value"
                  :error-messages="errors.gla_salary_multiple"
                  label="GLA Salary Multiple"
                  :disabled="!glaBenefit || !selectedSchemeType"
                ></v-text-field>
              </v-col>
              <v-col cols="4">
                <v-select
                  v-model="glaTerminalIllnessBenefit"
                  v-bind="glaTerminalIllnessBenefitAttrs"
                  variant="outlined"
                  density="compact"
                  label="Terminal Illness Benefit"
                  placeholder="Add Terminal Illness Benefit"
                  :error-messages="errors.gla_terminal_illness_benefit"
                  :items="groupStore.terminalIllnessBenefits"
                  :disabled="!glaBenefit || !selectedSchemeType"
                ></v-select>
              </v-col>
              <v-col cols="4">
                <v-select
                  v-model="glaWaitingPeriod"
                  v-bind="glaWaitingPeriodAttrs"
                  variant="outlined"
                  density="compact"
                  label="Waiting Period"
                  placeholder="Select Waiting Period"
                  :error-messages="errors.gla_waiting_period"
                  :items="waitingPeriods"
                  :disabled="!glaBenefit || !selectedSchemeType"
                ></v-select>
              </v-col>
              <v-col cols="4">
                <v-select
                  v-model="glaEducatorBenefit"
                  v-bind="glaEducatorBenefitAttrs"
                  variant="outlined"
                  density="compact"
                  label="Educator Benefit"
                  placeholder="Enable Educator Benefit"
                  :error-messages="errors.gla_educator_benefit"
                  :items="groupStore.yesNoItems"
                  :disabled="!glaBenefit || !selectedSchemeType"
                ></v-select>
              </v-col>
              <v-col v-if="glaEducatorBenefit === 'Yes'" cols="4">
                <v-select
                  v-model="glaEducatorBenefitType"
                  v-bind="glaEducatorBenefitTypeAttrs"
                  variant="outlined"
                  density="compact"
                  label="Educator Benefit Type"
                  placeholder="Select Educator Benefit Type"
                  :error-messages="errors.gla_educator_benefit_type"
                  :items="educatorBenefitTypes"
                  :disabled="!glaBenefit || !selectedSchemeType"
                ></v-select>
              </v-col>
              <v-col cols="4">
                <v-checkbox
                  v-model="glaConversionOnWithdrawal"
                  v-bind="glaConversionOnWithdrawalAttrs"
                  variant="outlined"
                  density="compact"
                  label="Conversion on Withdrawal"
                  :disabled="!glaBenefit || !selectedSchemeType"
                ></v-checkbox>
              </v-col>
              <v-col cols="4">
                <v-checkbox
                  v-model="glaConversionOnRetirement"
                  v-bind="glaConversionOnRetirementAttrs"
                  variant="outlined"
                  density="compact"
                  label="Conversion on Retirement"
                  :disabled="!glaBenefit || !selectedSchemeType"
                ></v-checkbox>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-window-item>
      <v-window-item class="pa-2" value="ptd">
        <base-card :show-actions="false">
          <template #header>
            <span class="headline">{{ ptdLabel }} Benefit</span>
          </template>
          <template #default>
            <v-switch
              v-model="ptdBenefit"
              color="primary"
              :label="`Enable ${ptdLabel} Benefit`"
              :disabled="!selectedSchemeType"
            ></v-switch>
            <v-divider class="mb-5"></v-divider>

            <v-row>
              <v-col cols="4">
                <v-select
                  v-model="ptdRiskType"
                  v-bind="ptdRiskTypeAttrs"
                  placeholder="Choose a Risk Type"
                  label="Risk Type"
                  variant="outlined"
                  density="compact"
                  :error-messages="errors.ptd_risk_type"
                  :items="ptdRiskTypes"
                  :disabled="!ptdBenefit || !selectedSchemeType"
                ></v-select>
              </v-col>
              <v-col cols="4">
                <v-select
                  v-model="ptdBenefitType"
                  v-bind="ptdBenefitTypeAttrs"
                  :error-messages="errors.ptd_benefit_type"
                  placeholder="Choose a Benefit Type"
                  label="Benefit Type"
                  variant="outlined"
                  density="compact"
                  :items="groupStore.benefitTypes"
                  :disabled="!ptdBenefit || !selectedSchemeType"
                ></v-select>
              </v-col>
              <v-col
                v-if="groupStore.group_pricing_quote.use_global_salary_multiple"
                cols="4"
              >
                <v-text-field
                  v-model:model-value="ptdSalaryMultiple"
                  v-bind="ptdSalaryMultipleAttrs"
                  :error-messages="errors.ptd_salary_multiple"
                  type="number"
                  variant="outlined"
                  density="compact"
                  placeholder="Enter a value"
                  label="PTD Salary Multiple"
                  :disabled="!ptdBenefit || !selectedSchemeType"
                ></v-text-field>
              </v-col>
              <v-col cols="4">
                <v-select
                  v-model:model-value="ptdDeferredPeriod"
                  v-bind="ptdDeferredPeriodAttrs"
                  :error-messages="errors.ptd_deferred_period"
                  :items="ptdDeferredPeriods"
                  variant="outlined"
                  density="compact"
                  placeholder="Select a Deferred Period"
                  label="Deferred Period (Months)"
                  :disabled="!ptdBenefit || !selectedSchemeType"
                ></v-select>
              </v-col>
              <v-col cols="4">
                <v-select
                  v-model="ptdDisabilityDefinition"
                  v-bind="ptdDisabilityDefinitionAttrs"
                  :error-messages="errors.ptd_disability_definition"
                  placeholder="Choose a Definition"
                  label="Disability Definition"
                  variant="outlined"
                  density="compact"
                  :items="ptdDisabilityDefinitions"
                  :disabled="!ptdBenefit || !selectedSchemeType"
                ></v-select>
              </v-col>
              <v-col cols="4">
                <v-select
                  v-model="ptdEducatorBenefit"
                  v-bind="ptdEducatorBenefitAttrs"
                  :error-messages="errors.ptd_educator_benefit"
                  variant="outlined"
                  density="compact"
                  label="Educator Benefit"
                  placeholder="Enable Educator Benefit"
                  :items="groupStore.yesNoItems"
                  :disabled="!ptdBenefit || !selectedSchemeType"
                ></v-select>
              </v-col>
              <v-col v-if="ptdEducatorBenefit === 'Yes'" cols="4">
                <v-select
                  v-model="ptdEducatorBenefitType"
                  v-bind="ptdEducatorBenefitTypeAttrs"
                  :error-messages="errors.ptd_educator_benefit_type"
                  variant="outlined"
                  density="compact"
                  label="Educator Benefit Type"
                  placeholder="Select Educator Benefit Type"
                  :items="educatorBenefitTypes"
                  :disabled="!ptdBenefit || !selectedSchemeType"
                ></v-select>
              </v-col>
              <v-col cols="4">
                <v-checkbox
                  v-model="ptdConversionOnWithdrawal"
                  v-bind="ptdConversionOnWithdrawalAttrs"
                  variant="outlined"
                  density="compact"
                  label="Conversion on Withdrawal"
                  :disabled="!ptdBenefit || !selectedSchemeType"
                ></v-checkbox>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-window-item>
      <v-window-item class="pa-2" value="ci">
        <v-row>
          <v-col>
            <base-card :show-actions="false">
              <template #header>
                <span class="headline">{{ ciLabel }} Benefit</span>
              </template>
              <template #default>
                <v-switch
                  v-model="ciBenefit"
                  color="primary"
                  :label="`Enable ${ciLabel} Benefit`"
                  :disabled="!selectedSchemeType"
                ></v-switch>
                <v-divider class="mb-5"></v-divider>

                <v-row>
                  <v-col cols="4">
                    <v-select
                      v-model="ciBenefitStructure"
                      v-bind="ciBenefitStructureAttrs"
                      :error-messages="errors.ci_benefit_structure"
                      placeholder="Choose a Benefit Structure"
                      label="Benefit Structure"
                      variant="outlined"
                      density="compact"
                      :items="groupStore.benefitStructures"
                      :disabled="!ciBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model="ciBenefitDefinition"
                      v-bind="ciBenefitDefinitionAttrs"
                      :error-messages="errors.ci_benefit_definition"
                      placeholder="Choose a Benefit Definition"
                      label="Benefit Definition"
                      variant="outlined"
                      density="compact"
                      :items="benefitDefinitions"
                      :disabled="!ciBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>

                  <v-col
                    v-if="
                      groupStore.group_pricing_quote.use_global_salary_multiple
                    "
                    cols="4"
                  >
                    <v-text-field
                      v-model:model-value="ciCriticalIllnessSalaryMultiple"
                      v-bind="ciCriticalIllnessSalaryMultipleAttrs"
                      :error-messages="
                        errors.ci_critical_illness_salary_multiple
                      "
                      placeholder="Enter a value"
                      label="Critical Illness Salary Multiple"
                      variant="outlined"
                      density="compact"
                      type="number"
                      :disabled="!ciBenefit || !selectedSchemeType"
                    ></v-text-field>
                  </v-col>
                  <v-col cols="4">
                    <v-checkbox
                      v-model="ciConversionOnWithdrawal"
                      v-bind="ciConversionOnWithdrawalAttrs"
                      variant="outlined"
                      density="compact"
                      label="Conversion on Withdrawal"
                      :disabled="!ciBenefit || !selectedSchemeType"
                    ></v-checkbox>
                  </v-col> </v-row
              ></template>
            </base-card>
          </v-col>
        </v-row>
      </v-window-item>
      <v-window-item class="pa-2" value="sgla">
        <v-row>
          <v-col>
            <base-card :show-actions="false">
              <template #header>
                <span class="headline">{{ sglaLabel }} Benefit</span>
              </template>
              <template #default>
                <v-switch
                  v-model="sglaBenefit"
                  color="primary"
                  :label="`Enable ${sglaLabel} Benefit`"
                  :disabled="!selectedSchemeType"
                ></v-switch>
                <v-divider class="mb-5"></v-divider>
                <v-row>
                  <v-col
                    v-if="
                      groupStore.group_pricing_quote.use_global_salary_multiple
                    "
                    cols="4"
                  >
                    <v-text-field
                      v-model:model-value="sglaSalaryMultiple"
                      v-bind="sglaSalaryMultipleAttrs"
                      :error-messages="errors.sgla_salary_multiple"
                      placeholder="Enter a value"
                      label="SGLA Salary Multiple"
                      variant="outlined"
                      density="compact"
                      type="number"
                      :disabled="!sglaBenefit || !selectedSchemeType"
                    ></v-text-field>
                  </v-col>
                </v-row>
              </template>
            </base-card>
          </v-col>
        </v-row>
      </v-window-item>
      <v-window-item class="pa-2" value="phi">
        <v-row>
          <v-col>
            <base-card :show-actions="false">
              <template #header
                >family_funeral_benefit
                <span class="headline">{{ phiLabel }} Benefit</span>
              </template>
              <template #default>
                <v-switch
                  v-model="phiBenefit"
                  color="primary"
                  :label="`Enable ${phiLabel} Benefit`"
                  :disabled="!selectedSchemeType"
                ></v-switch>
                <v-divider class="mb-5"></v-divider>
                <v-row>
                  <v-col cols="4">
                    <v-select
                      v-model="phiRiskType"
                      v-bind="phiRiskTypeAttrs"
                      :error-messages="errors.phi_risk_type"
                      placeholder="Choose a Risk Type"
                      label="Risk Type"
                      variant="outlined"
                      density="compact"
                      :items="phiRiskTypes"
                      :disabled="!phiBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-checkbox
                      v-model="phiUseTieredIncomeReplacementRatio"
                      v-bind="phiUseTieredIncomeReplacementRatioAttrs"
                      variant="outlined"
                      density="compact"
                      label="Use PHI Tiered Income Replacement Ratio"
                      :disabled="!phiBenefit || !selectedSchemeType"
                    ></v-checkbox>
                  </v-col>
                  <v-col v-if="phiUseTieredIncomeReplacementRatio" cols="4">
                    <v-select
                      v-model="phiTieredIncomeReplacementType"
                      v-bind="phiTieredIncomeReplacementTypeAttrs"
                      label="Tiered Table Type"
                      :items="tieredIncomeReplacementTypes"
                      variant="outlined"
                      density="compact"
                      :disabled="!phiBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>
                  <v-col
                    v-if="
                      phiUseTieredIncomeReplacementRatio &&
                      phiTieredIncomeReplacementType === 'custom'
                    "
                    cols="12"
                  >
                    <v-alert
                      v-if="!phiCustomTableExists"
                      type="warning"
                      density="compact"
                      variant="tonal"
                    >
                      Missing custom tiered income replacement table for this
                      scheme — this needs super admin attention. Calculations
                      cannot be run until the custom table is uploaded.
                    </v-alert>
                    <v-alert
                      v-else
                      type="success"
                      density="compact"
                      variant="tonal"
                    >
                      Custom tiered income replacement table is available for
                      this scheme.
                    </v-alert>
                  </v-col>
                  <v-col v-if="!phiUseTieredIncomeReplacementRatio" cols="4">
                    <v-text-field
                      v-model:model-value="phiIncomeReplacementPercentage"
                      v-bind="phiIncomeReplacementPercentageAttrs"
                      :error-messages="errors.phi_income_replacement_percentage"
                      placeholder="Enter a value"
                      label="Income Replacement %"
                      variant="outlined"
                      density="compact"
                      type="number"
                      :disabled="!phiBenefit || !selectedSchemeType"
                    ></v-text-field>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model:model-value="phiPremiumWaiver"
                      v-bind="phiPremiumWaiverAttrs"
                      :error-messages="errors.phi_premium_waiver"
                      placeholder="Enable Premium Waiver Benefit?"
                      label="Premium Waiver Benefit"
                      variant="outlined"
                      density="compact"
                      :items="groupStore.yesNoItems"
                      :disabled="!phiBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model:model-value="phiMedicalAidPremiumWaiver"
                      v-bind="phiMedicalAidPremiumWaiverAttrs"
                      :error-messages="errors.phi_medical_aid_premium_waiver"
                      placeholder="Enable Medical Aid Premium Waiver Benefit?"
                      label="Medical Aid Premium Waiver Benefit"
                      variant="outlined"
                      density="compact"
                      :items="groupStore.yesNoItems"
                      :disabled="!phiBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>

                  <v-col cols="4">
                    <v-select
                      v-model:model-value="phiBenefitEscalation"
                      v-bind="phiBenefitEscalationAttrs"
                      :error-messages="errors.phi_benefit_escalation"
                      placeholder="Choose an Escalation Option"
                      label="Benefit Escalation Option"
                      variant="outlined"
                      density="compact"
                      :items="incomeEscalations"
                      :disabled="!phiBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model:model-value="phiWaitingPeriod"
                      v-bind="phiWaitingPeriodAttrs"
                      :error-messages="errors.phi_waiting_period"
                      placeholder="Select Waiting Period"
                      label="Waiting Period (Months)"
                      variant="outlined"
                      density="compact"
                      :items="phiWaitingPeriods"
                      :disabled="!phiBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model:model-value="phiNormalRetirementAge"
                      v-bind="phiNormalRetirementAgeAttrs"
                      :error-messages="errors.phi_normal_retirement_age"
                      placeholder="Select Retirement Age"
                      label="Normal Retirement Age"
                      variant="outlined"
                      density="compact"
                      :items="normalRetirementAges"
                      :disabled="!phiBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model:model-value="phiDeferredPeriod"
                      v-bind="phiDeferredPeriodAttrs"
                      :error-messages="errors.phi_deferred_period"
                      placeholder="Select a Deferred Period"
                      label="Deferred Period (Months)"
                      variant="outlined"
                      density="compact"
                      :items="phiDeferredPeriods"
                      :disabled="!phiBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model="phiDisabilityDefinition"
                      v-bind="phiDisabilityDefinitionAttrs"
                      :error-messages="errors.phi_disability_definition"
                      placeholder="Choose a Definition"
                      label="Disability Definition"
                      variant="outlined"
                      density="compact"
                      :items="phiDisabilityDefinitions"
                      :disabled="!phiBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>
                </v-row>
              </template>
            </base-card>
          </v-col>
        </v-row>
      </v-window-item>
      <v-window-item class="pa-2" value="ttd">
        <v-row>
          <v-col>
            <base-card :show-actions="false">
              <template #header>
                <span class="headline">{{ ttdLabel }} Benefit</span>
              </template>
              <template #default>
                <v-switch
                  v-model="ttdBenefit"
                  color="primary"
                  :label="`Enable ${ttdLabel} Benefit`"
                  :disabled="!selectedSchemeType"
                ></v-switch>
                <v-divider class="mb-5"></v-divider>
                <v-row>
                  <v-col cols="4">
                    <v-select
                      v-model="ttdRiskType"
                      v-bind="ttdRiskTypeAttrs"
                      :error-messages="errors.ttd_risk_type"
                      placeholder="Choose a Risk Type"
                      label="Risk Type"
                      variant="outlined"
                      density="compact"
                      :items="ttdRiskTypes"
                      :disabled="!ttdBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-checkbox
                      v-model="ttdUseTieredIncomeReplacementRatio"
                      v-bind="ttdUseTieredIncomeReplacementRatioAttrs"
                      variant="outlined"
                      density="compact"
                      label="Use TTD Tiered Income Replacement Ratio"
                      :disabled="!ttdBenefit || !selectedSchemeType"
                    ></v-checkbox>
                  </v-col>
                  <v-col v-if="ttdUseTieredIncomeReplacementRatio" cols="4">
                    <v-select
                      v-model="ttdTieredIncomeReplacementType"
                      v-bind="ttdTieredIncomeReplacementTypeAttrs"
                      label="Tiered Table Type"
                      :items="tieredIncomeReplacementTypes"
                      variant="outlined"
                      density="compact"
                      :disabled="!ttdBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>
                  <v-col
                    v-if="
                      ttdUseTieredIncomeReplacementRatio &&
                      ttdTieredIncomeReplacementType === 'custom'
                    "
                    cols="12"
                  >
                    <v-alert
                      v-if="!ttdCustomTableExists"
                      type="warning"
                      density="compact"
                      variant="tonal"
                    >
                      Missing custom tiered income replacement table for this
                      scheme — this needs super admin attention. Calculations
                      cannot be run until the custom table is uploaded.
                    </v-alert>
                    <v-alert
                      v-else
                      type="success"
                      density="compact"
                      variant="tonal"
                    >
                      Custom tiered income replacement table is available for
                      this scheme.
                    </v-alert>
                  </v-col>
                  <v-col v-if="!ttdUseTieredIncomeReplacementRatio" cols="4">
                    <v-text-field
                      v-model:model-value="ttdIncomeReplacementPercentage"
                      v-bind="ttdIncomeReplacementPercentageAttrs"
                      :error-messages="errors.ttd_income_replacement_percentage"
                      placeholder="Enter a value"
                      label="Income Replacement %"
                      variant="outlined"
                      density="compact"
                      type="number"
                      :disabled="!ttdBenefit || !selectedSchemeType"
                    ></v-text-field>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model:model-value="ttdWaitingPeriod"
                      v-bind="ttdWaitingPeriodAttrs"
                      :error-messages="errors.ttd_waiting_period"
                      placeholder="Select Waiting Period"
                      label="Waiting Period (Months)"
                      variant="outlined"
                      density="compact"
                      :items="ttdWaitingPeriods"
                      :disabled="!ttdBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model:model-value="ttdDeferredPeriod"
                      v-bind="ttdDeferredPeriodAttrs"
                      :error-messages="errors.ttd_deferred_period"
                      placeholder="Select a Deferred Period"
                      label="Deferred Period (Months)"
                      variant="outlined"
                      density="compact"
                      :items="ttdDeferredPeriods"
                      :disabled="!ttdBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>
                  <v-col cols="4">
                    <v-select
                      v-model="ttdDisabilityDefinition"
                      v-bind="ttdDisabilityDefinitionAttrs"
                      :error-messages="errors.ttd_disability_definition"
                      placeholder="Choose a Definition"
                      label="Disability Definition"
                      variant="outlined"
                      density="compact"
                      :items="ttdDisabilityDefinitions"
                      :disabled="!ttdBenefit || !selectedSchemeType"
                    ></v-select>
                  </v-col>
                </v-row>
              </template>
            </base-card>
          </v-col>
        </v-row>
      </v-window-item>
      <v-window-item value="familyFuneral" class="pa-2">
        <v-row>
          <v-col>
            <base-card :show-actions="false">
              <template #header>
                <span class="headline">{{ familyFuneralLabel }} Benefit</span>
              </template>
              <template #default>
                <v-switch
                  v-model="familyFuneralBenefit"
                  color="primary"
                  :label="`Enable ${familyFuneralLabel} Benefit`"
                  :disabled="!selectedSchemeType"
                ></v-switch>
                <v-divider class="mb-5"></v-divider>
                <v-row>
                  <v-col cols="4">
                    <v-text-field
                      v-model:model-value="
                        formattedFamilyFuneralMainMemberFuneralSumAssured
                      "
                      v-bind="familyFuneralMainMemberFuneralSumAssuredAttrs"
                      :error-messages="
                        errors.family_funeral_main_member_funeral_sum_assured
                      "
                      type="text"
                      variant="outlined"
                      density="compact"
                      placeholder="Enter a value"
                      label="Main Member Sum Assured"
                      :disabled="!familyFuneralBenefit || !selectedSchemeType"
                    ></v-text-field>
                  </v-col>
                  <v-col cols="4">
                    <v-text-field
                      v-model:model-value="
                        formattedFamilyFuneralSpouseFuneralSumAssured
                      "
                      v-bind="familyFuneralSpouseFuneralSumAssuredAttrs"
                      :error-messages="
                        errors.family_funeral_spouse_funeral_sum_assured
                      "
                      type="text"
                      variant="outlined"
                      density="compact"
                      placeholder="Enter a value"
                      label="Spouse Funeral Sum Assured"
                      :disabled="!familyFuneralBenefit || !selectedSchemeType"
                    ></v-text-field>
                  </v-col>
                  <v-col cols="4">
                    <v-text-field
                      v-model:model-value="
                        formattedFamilyFuneralChildrenFuneralSumAssured
                      "
                      v-bind="familyFuneralChildrenFuneralSumAssuredAttrs"
                      :error-messages="
                        errors.family_funeral_children_funeral_sum_assured
                      "
                      type="text"
                      variant="outlined"
                      density="compact"
                      placeholder="Enter a value"
                      label="Children Funeral Sum Assured"
                      :disabled="!familyFuneralBenefit || !selectedSchemeType"
                    ></v-text-field>
                  </v-col>
                  <v-col cols="4">
                    <v-text-field
                      v-model:model-value="
                        formattedFamilyFuneralAdultDependantSumAssured
                      "
                      v-bind="familyFuneralAdultDependantSumAssuredAttrs"
                      :error-messages="
                        errors.family_funeral_adult_dependant_sum_assured
                      "
                      type="text"
                      variant="outlined"
                      density="compact"
                      placeholder="Enter a value"
                      label="Dependant Sum Assured"
                      :disabled="!familyFuneralBenefit || !selectedSchemeType"
                    ></v-text-field>
                  </v-col>
                  <v-col cols="4">
                    <v-text-field
                      v-model:model-value="
                        formattedFamilyFuneralParentFuneralSumAssured
                      "
                      v-bind="familyFuneralParentFuneralSumAssuredAttrs"
                      :error-messages="
                        errors.family_funeral_parent_funeral_sum_assured
                      "
                      type="text"
                      variant="outlined"
                      density="compact"
                      placeholder="Enter a value"
                      label="Parent Funeral Sum Assured"
                      :disabled="!familyFuneralBenefit || !selectedSchemeType"
                    ></v-text-field>
                  </v-col>

                  <v-col cols="4">
                    <v-text-field
                      v-model:model-value="familyFuneralMaxNumberChildren"
                      v-bind="familyFuneralMaxNumberChildrenAttrs"
                      :error-messages="
                        errors.family_funeral_max_number_children
                      "
                      type="number"
                      variant="outlined"
                      density="compact"
                      placeholder="Enter a value"
                      label="Maximum Number of Children"
                      :disabled="!familyFuneralBenefit || !selectedSchemeType"
                    ></v-text-field>
                  </v-col>
                  <v-col cols="4">
                    <v-text-field
                      v-model:model-value="
                        familyFuneralMaxNumberAdultDependants
                      "
                      v-bind="familyFuneralMaxNumberAdultDependantsAttrs"
                      :error-messages="
                        errors.family_funeral_max_number_adult_dependants
                      "
                      type="number"
                      variant="outlined"
                      density="compact"
                      placeholder="Enter a value"
                      label="Maximum Number of Dependants"
                      :disabled="!familyFuneralBenefit || !selectedSchemeType"
                    ></v-text-field>
                  </v-col>
                </v-row>
              </template>
            </base-card>
          </v-col>
        </v-row>
      </v-window-item>
    </v-window>
  </v-container>
</template>
<script setup lang="ts">
import { useGroupPricingStore } from '@/renderer/store/group_pricing'
import { computed, onBeforeMount, onMounted, ref, watch } from 'vue'
import BaseCard from '../BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useForm } from 'vee-validate'
import * as yup from 'yup'

const groupStore = useGroupPricingStore()
const benefitMaps: any = ref([])
const selectedSchemeType = ref(null)
const currentTab = ref('gla')

// formatted values
const formattedFamilyFuneralMainMemberFuneralSumAssured = computed({
  get: () => {
    if (typeof familyFuneralMainMemberFuneralSumAssured.value === 'number') {
      return familyFuneralMainMemberFuneralSumAssured.value.toLocaleString()
    }
    return familyFuneralMainMemberFuneralSumAssured.value
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      familyFuneralMainMemberFuneralSumAssured.value = isNaN(parsed)
        ? 0
        : parsed
    } else {
      familyFuneralMainMemberFuneralSumAssured.value = 0
    }
  }
})

const formattedFamilyFuneralSpouseFuneralSumAssured = computed({
  get: () => {
    if (typeof familyFuneralSpouseFuneralSumAssured.value === 'number') {
      return familyFuneralSpouseFuneralSumAssured.value.toLocaleString()
    }
    return familyFuneralSpouseFuneralSumAssured.value
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      familyFuneralSpouseFuneralSumAssured.value = isNaN(parsed) ? 0 : parsed
    } else {
      familyFuneralSpouseFuneralSumAssured.value = 0
    }
  }
})

const formattedFamilyFuneralChildrenFuneralSumAssured = computed({
  get: () => {
    if (typeof familyFuneralChildrenFuneralSumAssured.value === 'number') {
      return familyFuneralChildrenFuneralSumAssured.value.toLocaleString()
    }
    return familyFuneralChildrenFuneralSumAssured.value
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      familyFuneralChildrenFuneralSumAssured.value = isNaN(parsed) ? 0 : parsed
    } else {
      familyFuneralChildrenFuneralSumAssured.value = 0
    }
  }
})

const formattedFamilyFuneralAdultDependantSumAssured = computed({
  get: () => {
    if (typeof familyFuneralAdultDependantSumAssured.value === 'number') {
      return familyFuneralAdultDependantSumAssured.value.toLocaleString()
    }
    return familyFuneralAdultDependantSumAssured.value
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      familyFuneralAdultDependantSumAssured.value = isNaN(parsed) ? 0 : parsed
    } else {
      familyFuneralAdultDependantSumAssured.value = 0
    }
  }
})

const formattedFamilyFuneralParentFuneralSumAssured = computed({
  get: () => {
    if (typeof familyFuneralParentFuneralSumAssured.value === 'number') {
      return familyFuneralParentFuneralSumAssured.value.toLocaleString()
    }
    return familyFuneralParentFuneralSumAssured.value
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      familyFuneralParentFuneralSumAssured.value = isNaN(parsed) ? 0 : parsed
    } else {
      familyFuneralParentFuneralSumAssured.value = 0
    }
  }
})

// Define a type for scheme category objects
interface SchemeCategory {
  scheme_category: string
  [key: string]: any
}

// Helper: Find benefit data for a scheme type in scheme_categories
function getSchemeCategoryData(
  schemeCategory: string
): SchemeCategory | undefined {
  return (
    groupStore.group_pricing_quote.scheme_categories as SchemeCategory[]
  ).find((cat) => cat.scheme_category === schemeCategory)
}

// Helper: Save or update benefit data for a scheme type
function saveSchemeCategoryData(
  schemeCategory: string,
  benefitData: Record<string, any>
): void {
  const categories = groupStore.group_pricing_quote
    .scheme_categories as SchemeCategory[]
  const idx = categories.findIndex(
    (cat) => cat.scheme_category === schemeCategory
  )
  if (idx !== -1) {
    // Merge new benefitData with existing object
    categories[idx] = {
      ...categories[idx],
      ...benefitData
    }
  } else {
    categories.push({
      scheme_category: schemeCategory,
      ...benefitData
    })
  }
}

// When scheme type changes, load its data into the form
function onSchemeTypeChange(schemeType) {
  selectedSchemeType.value = schemeType
  const data: any = getSchemeCategoryData(schemeType)
  if (data) {
    // Populate form fields with data (example for a few fields)
    glaBenefit.value = data.gla_benefit || false
    ptdBenefit.value = data.ptd_benefit || false
    ciBenefit.value = data.ci_benefit || false
    sglaBenefit.value = data.sgla_benefit || false
    phiBenefit.value = data.phi_benefit || false
    ttdBenefit.value = data.ttd_benefit || false
    familyFuneralBenefit.value = data.family_funeral_benefit || false

    // Populate GLA fields
    glaBenefitType.value = data.gla_benefit_type || ''
    glaSalaryMultiple.value = data.gla_salary_multiple || 0
    glaTerminalIllnessBenefit.value = data.gla_terminal_illness_benefit || ''
    glaWaitingPeriod.value = data.gla_waiting_period ?? null
    glaEducatorBenefit.value = data.gla_educator_benefit || ''
    glaEducatorBenefitType.value = data.gla_educator_benefit_type || null
    glaConversionOnWithdrawal.value = data.gla_conversion_on_withdrawal || false
    glaConversionOnRetirement.value = data.gla_conversion_on_retirement || false

    // Populate PTD fields
    ptdRiskType.value = data.ptd_risk_type || ''
    ptdBenefitType.value = data.ptd_benefit_type || ''
    ptdSalaryMultiple.value = data.ptd_salary_multiple || 0
    ptdDeferredPeriod.value = data.ptd_deferred_period ?? null
    ptdDisabilityDefinition.value = data.ptd_disability_definition || ''
    ptdEducatorBenefit.value = data.ptd_educator_benefit || ''
    ptdEducatorBenefitType.value = data.ptd_educator_benefit_type || null
    ptdConversionOnWithdrawal.value = data.ptd_conversion_on_withdrawal || false

    // Populate CI fields
    ciBenefitStructure.value = data.ci_benefit_structure || ''
    ciBenefitDefinition.value = data.ci_benefit_definition || ''
    ciCriticalIllnessSalaryMultiple.value =
      data.ci_critical_illness_salary_multiple || 0
    ciConversionOnWithdrawal.value = data.ci_conversion_on_withdrawal || false
    // Populate SGLA fields
    sglaSalaryMultiple.value = data.sgla_salary_multiple || 0
    // Populate PHI fields
    phiRiskType.value = data.phi_risk_type || ''
    phiIncomeReplacementPercentage.value =
      data.phi_income_replacement_percentage || 0
    phiUseTieredIncomeReplacementRatio.value =
      data.phi_use_tiered_income_replacement_ratio || false
    phiTieredIncomeReplacementType.value =
      data.phi_tiered_income_replacement_type || 'standard'
    if (phiTieredIncomeReplacementType.value === 'custom') {
      checkCustomTableExists('phi')
    }
    phiPremiumWaiver.value = data.phi_premium_waiver || ''
    phiMedicalAidPremiumWaiver.value = data.phi_medical_aid_premium_waiver || ''
    phiBenefitEscalation.value = data.phi_benefit_escalation || ''
    phiWaitingPeriod.value = data.phi_waiting_period ?? null
    phiNormalRetirementAge.value = data.phi_normal_retirement_age ?? null
    phiDeferredPeriod.value = data.phi_deferred_period ?? null
    phiDisabilityDefinition.value = data.phi_disability_definition || ''
    // Populate TTD fields
    ttdRiskType.value = data.ttd_risk_type || ''
    ttdIncomeReplacementPercentage.value =
      data.ttd_income_replacement_percentage || 0
    ttdUseTieredIncomeReplacementRatio.value =
      data.ttd_use_tiered_income_replacement_ratio || false
    ttdTieredIncomeReplacementType.value =
      data.ttd_tiered_income_replacement_type || 'standard'
    if (ttdTieredIncomeReplacementType.value === 'custom') {
      checkCustomTableExists('ttd')
    }
    // ttdPremiumWaiverPercentage.value = data.ttd_premium_waiver_percentage || 0
    ttdWaitingPeriod.value = data.ttd_waiting_period ?? null
    ttdDeferredPeriod.value = data.ttd_deferred_period ?? null
    ttdDisabilityDefinition.value = data.ttd_disability_definition || ''
    // Populate Family Funeral fields
    familyFuneralMainMemberFuneralSumAssured.value =
      data.family_funeral_main_member_funeral_sum_assured || 0
    familyFuneralSpouseFuneralSumAssured.value =
      data.family_funeral_spouse_funeral_sum_assured || 0
    familyFuneralChildrenFuneralSumAssured.value =
      data.family_funeral_children_funeral_sum_assured || 0
    familyFuneralAdultDependantSumAssured.value =
      data.family_funeral_adult_dependant_sum_assured || 0
    familyFuneralParentFuneralSumAssured.value =
      data.family_funeral_parent_funeral_sum_assured || 0
    familyFuneralMaxNumberChildren.value =
      data.family_funeral_max_number_children || 0
    familyFuneralMaxNumberAdultDependants.value =
      data.family_funeral_max_number_adult_dependants || 0
    selectedRegion.value = data.region || ''
  } else {
    // Reset form fields
    selectedRegion.value = ''
    glaSalaryMultiple.value = 0
    glaTerminalIllnessBenefit.value = null
    glaWaitingPeriod.value = null
    glaEducatorBenefit.value = null
    glaEducatorBenefitType.value = null
    glaConversionOnWithdrawal.value = false
    glaConversionOnRetirement.value = false
    glaBenefit.value = false
    ptdBenefit.value = false
    ciBenefit.value = false
    sglaBenefit.value = false
    phiBenefit.value = false
    ttdBenefit.value = false
    familyFuneralBenefit.value = false
    // ...repeat for all benefit fields as needed
    // reset PTD fields
    ptdRiskType.value = null
    ptdBenefitType.value = null
    ptdSalaryMultiple.value = 0
    ptdDeferredPeriod.value = null
    ptdDisabilityDefinition.value = null
    ptdEducatorBenefit.value = null
    ptdEducatorBenefitType.value = null
    ptdConversionOnWithdrawal.value = false
    // reset CI fields
    ciBenefitStructure.value = null
    ciBenefitDefinition.value = null
    ciCriticalIllnessSalaryMultiple.value = 0
    ciConversionOnWithdrawal.value = false
    // reset SGLA fields
    sglaSalaryMultiple.value = 0
    // reset PHI fields
    phiRiskType.value = null
    phiIncomeReplacementPercentage.value = 0
    phiUseTieredIncomeReplacementRatio.value = false
    phiTieredIncomeReplacementType.value = 'standard'
    phiPremiumWaiver.value = null
    phiMedicalAidPremiumWaiver.value = null
    phiBenefitEscalation.value = null
    phiWaitingPeriod.value = null
    phiNormalRetirementAge.value = 0
    phiDeferredPeriod.value = null
    phiDisabilityDefinition.value = null
    // reset TTD fields
    ttdRiskType.value = null
    ttdIncomeReplacementPercentage.value = 0
    ttdUseTieredIncomeReplacementRatio.value = false
    ttdTieredIncomeReplacementType.value = 'standard'
    // ttdPremiumWaiverPercentage.value = 0
    ttdWaitingPeriod.value = null
    ttdDeferredPeriod.value = null
    ttdDisabilityDefinition.value = null
    // reset Family Funeral fields
    familyFuneralMainMemberFuneralSumAssured.value = 0
    familyFuneralSpouseFuneralSumAssured.value = 0
    familyFuneralChildrenFuneralSumAssured.value = 0
    familyFuneralAdultDependantSumAssured.value = 0
    familyFuneralParentFuneralSumAssured.value = 0
    familyFuneralMaxNumberChildren.value = 0
    familyFuneralMaxNumberAdultDependants.value = 0
  }
}

const savedSchemeCategories = ref<string[]>([])
const selectedRegion = ref<string>('')
const availableRegions = ref<string[]>([])

// Save current form data for selected scheme type
function saveCurrentSchemeCategory() {
  if (!selectedSchemeType.value || typeof selectedSchemeType.value !== 'string')
    return
  // Use vee-validate validate() to check form
  validate().then((result) => {
    if (result.valid) {
      const benefitData = {
        scheme_category: selectedSchemeType.value,
        region: selectedRegion.value || '',
        ptd_benefit: ptdBenefit.value,
        gla_benefit: glaBenefit.value,
        ci_benefit: ciBenefit.value,
        sgla_benefit: sglaBenefit.value,
        phi_benefit: phiBenefit.value,
        ttd_benefit: ttdBenefit.value,
        family_funeral_benefit: familyFuneralBenefit.value,
        ...(glaBenefit.value && {
          gla_benefit_type: glaBenefitType.value,
          gla_salary_multiple: Number(glaSalaryMultiple.value),
          gla_terminal_illness_benefit: glaTerminalIllnessBenefit.value,
          gla_waiting_period: Number(glaWaitingPeriod.value),
          gla_educator_benefit: glaEducatorBenefit.value,
          ...(glaEducatorBenefit.value === 'Yes' && {
            gla_educator_benefit_type: glaEducatorBenefitType.value
          }),
          gla_conversion_on_withdrawal: !!glaConversionOnWithdrawal.value,
          gla_conversion_on_retirement: !!glaConversionOnRetirement.value
        }),
        ...(ptdBenefit.value && {
          ptd_risk_type: ptdRiskType.value,
          ptd_benefit_type: ptdBenefitType.value,
          ptd_salary_multiple: Number(ptdSalaryMultiple.value),
          ptd_deferred_period: Number(ptdDeferredPeriod.value),
          ptd_disability_definition: ptdDisabilityDefinition.value,
          ptd_educator_benefit: ptdEducatorBenefit.value,
          ...(ptdEducatorBenefit.value === 'Yes' && {
            ptd_educator_benefit_type: ptdEducatorBenefitType.value
          }),
          ptd_conversion_on_withdrawal: !!ptdConversionOnWithdrawal.value
        }),
        ...(ciBenefit.value && {
          ci_benefit_structure: ciBenefitStructure.value,
          ci_benefit_definition: ciBenefitDefinition.value,
          ci_critical_illness_salary_multiple: Number(
            ciCriticalIllnessSalaryMultiple.value
          ),
          ci_conversion_on_withdrawal: !!ciConversionOnWithdrawal.value
        }),
        ...(sglaBenefit.value && {
          sgla_salary_multiple: Number(sglaSalaryMultiple.value)
        }),
        ...(phiBenefit.value && {
          phi_risk_type: phiRiskType.value,
          phi_use_tiered_income_replacement_ratio:
            phiUseTieredIncomeReplacementRatio.value,
          phi_tiered_income_replacement_type:
            phiTieredIncomeReplacementType.value || 'standard',
          phi_income_replacement_percentage: Number(
            phiIncomeReplacementPercentage.value
          ),
          phi_premium_waiver: phiPremiumWaiver.value,
          phi_medical_aid_premium_waiver: phiMedicalAidPremiumWaiver.value,
          phi_benefit_escalation: phiBenefitEscalation.value,
          phi_waiting_period: Number(phiWaitingPeriod.value),
          phi_normal_retirement_age: Number(phiNormalRetirementAge.value),
          phi_deferred_period: Number(phiDeferredPeriod.value),
          phi_disability_definition: phiDisabilityDefinition.value
        }),
        ...(ttdBenefit.value && {
          ttd_risk_type: ttdRiskType.value,
          ttd_use_tiered_income_replacement_ratio:
            ttdUseTieredIncomeReplacementRatio.value,
          ttd_tiered_income_replacement_type:
            ttdTieredIncomeReplacementType.value || 'standard',
          ttd_income_replacement_percentage: Number(
            ttdIncomeReplacementPercentage.value
          ),
          ttd_waiting_period: Number(ttdWaitingPeriod.value),
          ttd_deferred_period: Number(ttdDeferredPeriod.value),
          ttd_disability_definition: ttdDisabilityDefinition.value
        }),
        ...(familyFuneralBenefit.value && {
          family_funeral_main_member_funeral_sum_assured: Number(
            familyFuneralMainMemberFuneralSumAssured.value
          ),
          family_funeral_spouse_funeral_sum_assured: Number(
            familyFuneralSpouseFuneralSumAssured.value
          ),
          family_funeral_children_funeral_sum_assured: Number(
            familyFuneralChildrenFuneralSumAssured.value
          ),
          family_funeral_adult_dependant_sum_assured: Number(
            familyFuneralAdultDependantSumAssured.value
          ),
          family_funeral_parent_funeral_sum_assured: Number(
            familyFuneralParentFuneralSumAssured.value
          ),
          family_funeral_max_number_children: Number(
            familyFuneralMaxNumberChildren.value
          ),
          family_funeral_max_number_adult_dependants: Number(
            familyFuneralMaxNumberAdultDependants.value
          )
        })
      }
      if (typeof selectedSchemeType.value === 'string') {
        saveSchemeCategoryData(selectedSchemeType.value, benefitData)
        snackbar.value.show = true

        // Show snackbar confirmation
        snackbar.value.message = `Scheme category saved: ${selectedSchemeType.value}. Benefits enabled: ${getEnabledBenefitsMessage(benefitData)}`

        if (!savedSchemeCategories.value.includes(selectedSchemeType.value)) {
          savedSchemeCategories.value.push(selectedSchemeType.value)
        }
        if (
          savedSchemeCategories.value.length ===
          groupStore.group_pricing_quote.selected_scheme_categories.length
        ) {
          emit('all_schemes_saved', selectedSchemeType.value, true)
          console.log('all saved')
        }
      }
    }
    // else: do not save, errors will be shown by vee-validate
  })
}

const ptdLabel = ref('PTD')
const ciLabel = ref('CI')
const sglaLabel = ref('SGLA')
const phiLabel = ref('PHI')
const ttdLabel = ref('TTD')
const glaLabel = ref('GLA')
const familyFuneralLabel = ref('Family Funeral')
const benefitDefinitions: any = ref(['Lump Sum', 'Monthly'])
const incomeEscalations: any = ref([])
const ptdDisabilityDefinitions: any = ref([])
const phiDisabilityDefinitions: any = ref([])
const ttdDisabilityDefinitions: any = ref([])
const ciWaitingPeriods: any = ref([])
const ptdWaitingPeriods: any = ref([])
const ttdWaitingPeriods: any = ref([])
const phiWaitingPeriods: any = ref([])
const phiDeferredPeriods: any = ref([])
const ptdDeferredPeriods: any = ref([])
const ttdDeferredPeriods: any = ref([])
const normalRetirementAges: any = ref([])
const ttdRiskTypes: any = ref([])
const phiRiskTypes: any = ref([])
const ptdRiskTypes: any = ref([])
const waitingPeriods: any = ref([])
const glaBenefitTypes: any = ref([])
const educatorBenefitTypes = ref<string[]>([])

const tieredIncomeReplacementTypes = [
  { title: 'Standard', value: 'standard' },
  { title: 'Custom', value: 'custom' }
]
const phiCustomTableExists = ref(false)
const ttdCustomTableExists = ref(false)

const validationSchema = yup.object({
  gla_benefit: yup.boolean().nullable(),

  gla_benefit_type: yup.string().when('gla_benefit', {
    is: true,
    then: (schema) => schema.required('Benefit type is required'),
    otherwise: (schema) => schema.nullable()
  }),

  gla_salary_multiple: yup.number().when(['gla_benefit'], {
    is: (glaBenefit) => {
      // Access groupStore directly for the second condition
      return (
        glaBenefit === true &&
        groupStore.group_pricing_quote.use_global_salary_multiple
      )
    },
    then: (schema) =>
      schema
        .required('Salary multiple is required')
        .positive('Salary multiple must be a positive number'),
    otherwise: (schema) => schema.nullable()
  }),

  gla_terminal_illness_benefit: yup.string().when('gla_benefit', {
    is: true,
    then: (schema) => schema.required('Terminal illness benefit is required'),
    otherwise: (schema) => schema.nullable()
  }),
  gla_waiting_period: yup.number().when('gla_benefit', {
    is: true,
    then: (schema) =>
      schema
        .required('Waiting period is required')
        .min(0, 'Waiting period must be at least 0'),
    otherwise: (schema) => schema.nullable()
  }),
  gla_educator_benefit: yup.string().when('gla_benefit', {
    is: true,
    then: (schema) => schema.required('Educator benefit is required'),
    otherwise: (schema) => schema.nullable()
  }),
  gla_educator_benefit_type: yup
    .string()
    .when(['gla_benefit', 'gla_educator_benefit'], {
      is: (glaBenefit: boolean, glaEdu: string) =>
        glaBenefit === true && glaEdu === 'Yes',
      then: (schema) => schema.required('Educator benefit type is required'),
      otherwise: (schema) => schema.nullable()
    }),
  gla_conversion_on_withdrawal: yup.boolean().nullable(),
  gla_conversion_on_retirement: yup.boolean().nullable(),
  ptd_conversion_on_withdrawal: yup.boolean().nullable(),
  ci_conversion_on_withdrawal: yup.boolean().nullable(),
  ptd_benefit: yup.boolean().nullable(),
  ci_benefit: yup.boolean().nullable(),
  sgla_benefit: yup.boolean().nullable(),
  phi_benefit: yup.boolean().nullable(),
  ttd_benefit: yup.boolean().nullable(),
  phi_use_tiered_income_replacement_ratio: yup.boolean().nullable(),
  phi_tiered_income_replacement_type: yup.string().nullable(),
  ttd_use_tiered_income_replacement_ratio: yup.boolean().nullable(),
  ttd_tiered_income_replacement_type: yup.string().nullable(),
  family_funeral_benefit: yup.boolean().nullable(),
  ptd_salary_multiple: yup.number().when(['ptd_benefit'], {
    is: (ptdBenefit) => {
      // Access groupStore directly for the second condition
      return (
        ptdBenefit === true &&
        groupStore.group_pricing_quote.use_global_salary_multiple
      )
    },
    then: (schema) =>
      schema
        .required('Salary multiple is required')
        .positive('Salary multiple must be a positive number'),
    otherwise: (schema) => schema.nullable()
  }),
  ptd_risk_type: yup.string().when('ptd_benefit', {
    is: true,
    then: (schema) => schema.required('Risk type is required'),
    otherwise: (schema) => schema.nullable()
  }),
  ptd_benefit_type: yup.string().when('ptd_benefit', {
    is: true,
    then: (schema) => schema.required('Benefit type is required'),
    otherwise: (schema) => schema.nullable()
  }),
  ptd_deferred_period: yup.number().when('ptd_benefit', {
    is: true,
    then: (schema) =>
      schema
        .required('Deferred period is required')
        .min(0, 'Deferred period must be at least 0'),
    otherwise: (schema) => schema.nullable()
  }),
  ptd_disability_definition: yup.string().when('ptd_benefit', {
    is: true,
    then: (schema) => schema.required('Disability definition is required'),
    otherwise: (schema) => schema.nullable()
  }),
  ptd_educator_benefit: yup.string().when('ptd_benefit', {
    is: true,
    then: (schema) => schema.required('Educator benefit is required'),
    otherwise: (schema) => schema.nullable()
  }),
  ptd_educator_benefit_type: yup
    .string()
    .when(['ptd_benefit', 'ptd_educator_benefit'], {
      is: (ptdBenefit: boolean, ptdEdu: string) =>
        ptdBenefit === true && ptdEdu === 'Yes',
      then: (schema) => schema.required('Educator benefit type is required'),
      otherwise: (schema) => schema.nullable()
    }),
  ci_benefit_structure: yup.string().when('ci_benefit', {
    is: true,
    then: (schema) => schema.required('Benefit structure is required'),
    otherwise: (schema) => schema.nullable()
  }),
  ci_benefit_definition: yup.string().when('ci_benefit', {
    is: true,
    then: (schema) => schema.required('Benefit definition is required'),
    otherwise: (schema) => schema.nullable()
  }),
  ci_critical_illness_salary_multiple: yup.number().when(['ci_benefit'], {
    is: (ciBenefit) => {
      return (
        ciBenefit === true &&
        groupStore.group_pricing_quote.use_global_salary_multiple
      )
    },
    then: (schema) =>
      schema
        .required('Critical illness salary multiple is required')
        .positive('Critical illness salary multiple must be a positive number'),
    otherwise: (schema) => schema.nullable()
  }),
  sgla_salary_multiple: yup.number().when(['sgla_benefit'], {
    is: (sglaBenefit) => {
      return (
        sglaBenefit === true &&
        groupStore.group_pricing_quote.use_global_salary_multiple
      )
    },
    then: (schema) =>
      schema
        .required('SGLA salary multiple is required')
        .positive('SGLA salary multiple must be a positive number'),
    otherwise: (schema) => schema.nullable()
  }),
  phi_income_replacement_percentage: yup
    .number()
    .when(['phi_benefit', 'phi_use_tiered_income_replacement_ratio'], {
      is: (phiBenefit: boolean, useTiered: boolean) =>
        phiBenefit === true && useTiered !== true,
      then: (schema) =>
        schema
          .required('Income replacement percentage is required')
          .positive('Income replacement percentage must be a positive number'),
      otherwise: (schema) => schema.nullable()
    }),
  phi_premium_waiver: yup.string().when('phi_benefit', {
    is: true,
    then: (schema) => schema.required('Premium waiver is required'),
    otherwise: (schema) => schema.nullable()
  }),
  phi_medical_aid_premium_waiver: yup.string().when('phi_benefit', {
    is: true,
    then: (schema) => schema.required('Medical aid premium waiver is required'),
    otherwise: (schema) => schema.nullable()
  }),
  phi_benefit_escalation: yup.string().when('phi_benefit', {
    is: true,
    then: (schema) => schema.required('Benefit escalation is required'),
    otherwise: (schema) => schema.nullable()
  }),
  phi_waiting_period: yup.number().when('phi_benefit', {
    is: true,
    then: (schema) =>
      schema
        .required('Waiting period is required')
        .min(0, 'Waiting period must be at least 0'),
    otherwise: (schema) => schema.nullable()
  }),
  phi_normal_retirement_age: yup.number().when('phi_benefit', {
    is: true,
    then: (schema) =>
      schema
        .required('A normal retirement age is required')
        .positive('Normal retirement age must be a positive number'),
    otherwise: (schema) => schema.nullable()
  }),
  phi_deferred_period: yup.number().when('phi_benefit', {
    is: true,
    then: (schema) =>
      schema
        .required('Deferred period is required')
        .min(0, 'Deferred period must be at least 0'),
    otherwise: (schema) => schema.nullable()
  }),
  phi_disability_definition: yup.string().when('phi_benefit', {
    is: true,
    then: (schema) => schema.required('Disability definition is required'),
    otherwise: (schema) => schema.nullable()
  }),
  phi_risk_type: yup.string().when('phi_benefit', {
    is: true,
    then: (schema) => schema.required('Risk type is required'),
    otherwise: (schema) => schema.nullable()
  }),
  ttd_income_replacement_percentage: yup
    .number()
    .when(['ttd_benefit', 'ttd_use_tiered_income_replacement_ratio'], {
      is: (ttdBenefit: boolean, useTiered: boolean) =>
        ttdBenefit === true && useTiered !== true,
      then: (schema) =>
        schema
          .required('Income replacement percentage is required')
          .positive('Income replacement percentage must be a positive number'),
      otherwise: (schema) => schema.nullable()
    }),
  ttd_waiting_period: yup.number().when('ttd_benefit', {
    is: true,
    then: (schema) =>
      schema
        .required('Waiting period is required')
        .min(0, 'Waiting period must be at least 0'),
    otherwise: (schema) => schema.nullable()
  }),
  ttd_deferred_period: yup.number().when('ttd_benefit', {
    is: true,
    then: (schema) =>
      schema
        .required('Deferred period is required')
        .min(0, 'Deferred period must be at least 0'),
    otherwise: (schema) => schema.nullable()
  }),
  ttd_disability_definition: yup.string().when('ttd_benefit', {
    is: true,
    then: (schema) => schema.required('Disability definition is required'),
    otherwise: (schema) => schema.nullable()
  }),
  ttd_risk_type: yup.string().when('ttd_benefit', {
    is: true,
    then: (schema) => schema.required('Risk type is required'),
    otherwise: (schema) => schema.nullable()
  }),
  family_funeral_main_member_funeral_sum_assured: yup
    .number()
    .when('family_funeral_benefit', {
      is: true,
      then: (schema) =>
        schema
          .required('Main member funeral sum assured is required')
          .positive(
            'Main member funeral sum assured must be a positive number'
          ),
      otherwise: (schema) => schema.nullable()
    }),
  family_funeral_spouse_funeral_sum_assured: yup
    .number()
    .when('family_funeral_benefit', {
      is: true,
      then: (schema) =>
        schema
          .required('Spouse funeral sum assured is required')
          .min(0, 'Spouse funeral sum assured must be at least 0'),

      otherwise: (schema) => schema.nullable()
    }),
  family_funeral_children_funeral_sum_assured: yup
    .number()
    .when('family_funeral_benefit', {
      is: true,
      then: (schema) =>
        schema
          .required('Children funeral sum assured is required')
          .min(0, 'Children funeral sum assured must be at least 0'),

      otherwise: (schema) => schema.nullable()
    }),
  family_funeral_adult_dependant_sum_assured: yup
    .number()
    .when('family_funeral_benefit', {
      is: true,
      then: (schema) =>
        schema
          .required('Adult dependant sum assured is required')
          .min(0, 'Adult dependant sum assured must be at least 0'),
      otherwise: (schema) => schema.nullable()
    }),
  family_funeral_parent_funeral_sum_assured: yup
    .number()
    .when('family_funeral_benefit', {
      is: true,
      then: (schema) =>
        schema
          .required('Parent funeral sum assured is required')
          .min(0, 'Parent funeral sum assured must be at least 0'),
      otherwise: (schema) => schema.nullable()
    }),
  family_funeral_max_number_children: yup
    .number()
    .when('family_funeral_benefit', {
      is: true,
      then: (schema) =>
        schema
          .required('Maximum number of children is required')
          .min(0, 'Maximum number of children must be at least 0'),
      otherwise: (schema) => schema.nullable()
    })
    .when('family_funeral_children_funeral_sum_assured', {
      is: (value) => value > 0,
      then: (schema) =>
        schema.min(
          1,
          'Maximum number of children must be at least 1 when children funeral sum assured is greater than 0'
        )
    }),
  family_funeral_max_number_adult_dependants: yup
    .number()
    .when('family_funeral_benefit', {
      is: true,
      then: (schema) =>
        schema
          .required('Maximum number of adult dependants is required')
          .min(0, 'Maximum number of adult dependants must be at least 0'),
      otherwise: (schema) => schema.nullable()
    })
    .when('family_funeral_adult_dependant_sum_assured', {
      is: (value) => value > 0,
      then: (schema) =>
        schema.min(
          1,
          'Maximum number of adult dependants must be at least 1 when adult dependant sum assured is greater than 0'
        )
    })
})

const { handleSubmit, defineField, errors, validate } = useForm({
  validationSchema,
  initialValues: {
    use_global_salary_multiple:
      groupStore.group_pricing_quote.use_global_salary_multiple,
    gla_salary_multiple:
      groupStore.scheme_category_template.gla_salary_multiple,
    gla_terminal_illness_benefit:
      groupStore.scheme_category_template.gla_terminal_illness_benefit,
    gla_waiting_period: groupStore.scheme_category_template.gla_waiting_period,
    gla_educator_benefit:
      groupStore.scheme_category_template.gla_educator_benefit,
    gla_educator_benefit_type:
      groupStore.scheme_category_template.gla_educator_benefit_type,
    gla_benefit_type: groupStore.scheme_category_template.gla_benefit_type,
    gla_conversion_on_withdrawal:
      groupStore.scheme_category_template.gla_conversion_on_withdrawal,
    gla_conversion_on_retirement:
      groupStore.scheme_category_template.gla_conversion_on_retirement,
    ptd_conversion_on_withdrawal:
      groupStore.scheme_category_template.ptd_conversion_on_withdrawal,
    ci_conversion_on_withdrawal:
      groupStore.scheme_category_template.ci_conversion_on_withdrawal,

    gla_benefit: groupStore.scheme_category_template.gla_benefit,
    ptd_benefit: groupStore.scheme_category_template.ptd_benefit,
    ci_benefit: groupStore.scheme_category_template.ci_benefit,
    sgla_benefit: groupStore.scheme_category_template.sgla_benefit,
    phi_benefit: groupStore.scheme_category_template.phi_benefit,
    ttd_benefit: groupStore.scheme_category_template.ttd_benefit,
    family_funeral_benefit:
      groupStore.scheme_category_template.family_funeral_benefit,
    ptd_risk_type: groupStore.scheme_category_template.ptd_risk_type,
    ptd_benefit_type: groupStore.scheme_category_template.ptd_benefit_type,
    ptd_salary_multiple:
      groupStore.scheme_category_template.ptd_salary_multiple,
    ptd_deferred_period:
      groupStore.scheme_category_template.ptd_deferred_period,
    ptd_disability_definition:
      groupStore.scheme_category_template.ptd_disability_definition,
    ptd_educator_benefit:
      groupStore.scheme_category_template.ptd_educator_benefit,
    ptd_educator_benefit_type:
      groupStore.scheme_category_template.ptd_educator_benefit_type,
    ci_benefit_structure:
      groupStore.scheme_category_template.ci_benefit_structure,
    ci_benefit_definition:
      groupStore.scheme_category_template.ci_benefit_definition,
    ci_critical_illness_salary_multiple:
      groupStore.scheme_category_template.ci_critical_illness_salary_multiple,
    sgla_salary_multiple:
      groupStore.scheme_category_template.sgla_salary_multiple,
    phi_income_replacement_percentage:
      groupStore.scheme_category_template.phi_income_replacement_percentage,
    phi_use_tiered_income_replacement_ratio:
      groupStore.scheme_category_template
        .phi_use_tiered_income_replacement_ratio,
    phi_premium_waiver: groupStore.scheme_category_template.phi_premium_waiver,
    phi_medical_aid_premium_waiver:
      groupStore.scheme_category_template.phi_medical_aid_premium_waiver,
    phi_benefit_escalation:
      groupStore.scheme_category_template.phi_benefit_escalation,
    phi_waiting_period: groupStore.scheme_category_template.phi_waiting_period,
    phi_normal_retirement_age:
      groupStore.scheme_category_template.phi_normal_retirement_age,
    phi_deferred_period:
      groupStore.scheme_category_template.phi_deferred_period,
    phi_disability_definition:
      groupStore.scheme_category_template.phi_disability_definition,
    phi_risk_type: groupStore.scheme_category_template.phi_risk_type,
    ttd_income_replacement_percentage:
      groupStore.scheme_category_template.ttd_income_replacement_percentage,
    ttd_use_tiered_income_replacement_ratio:
      groupStore.scheme_category_template
        .ttd_use_tiered_income_replacement_ratio,
    ttd_premium_waiver_percentage:
      groupStore.scheme_category_template.ttd_premium_waiver_percentage,
    ttd_waiting_period: groupStore.scheme_category_template.ttd_waiting_period,
    ttd_deferred_period:
      groupStore.scheme_category_template.ttd_deferred_period,
    ttd_disability_definition:
      groupStore.scheme_category_template.ttd_disability_definition,
    ttd_risk_type: groupStore.scheme_category_template.ttd_risk_type,
    family_funeral_main_member_funeral_sum_assured:
      groupStore.scheme_category_template
        .family_funeral_main_member_funeral_sum_assured,
    family_funeral_spouse_funeral_sum_assured:
      groupStore.scheme_category_template
        .family_funeral_spouse_funeral_sum_assured,
    family_funeral_children_funeral_sum_assured:
      groupStore.scheme_category_template
        .family_funeral_children_funeral_sum_assured,
    family_funeral_adult_dependant_sum_assured:
      groupStore.scheme_category_template
        .family_funeral_adult_dependant_sum_assured,
    family_funeral_parent_funeral_sum_assured:
      groupStore.scheme_category_template
        .family_funeral_parent_funeral_sum_assured,
    family_funeral_max_number_children:
      groupStore.scheme_category_template.family_funeral_max_number_children,
    family_funeral_max_number_adult_dependants:
      groupStore.scheme_category_template
        .family_funeral_max_number_adult_dependants
  }
})

const [glaSalaryMultiple, glaSalaryMultipleAttrs] = defineField(
  'gla_salary_multiple'
)
const [glaTerminalIllnessBenefit, glaTerminalIllnessBenefitAttrs] = defineField(
  'gla_terminal_illness_benefit'
)
const [glaWaitingPeriod, glaWaitingPeriodAttrs] =
  defineField('gla_waiting_period')
const [glaEducatorBenefit, glaEducatorBenefitAttrs] = defineField(
  'gla_educator_benefit'
)
const [glaEducatorBenefitType, glaEducatorBenefitTypeAttrs] = defineField(
  'gla_educator_benefit_type'
)
const [glaBenefitType, glaBenefitTypeAttrs] = defineField('gla_benefit_type')
const [glaConversionOnWithdrawal, glaConversionOnWithdrawalAttrs] = defineField(
  'gla_conversion_on_withdrawal'
)
const [glaConversionOnRetirement, glaConversionOnRetirementAttrs] = defineField(
  'gla_conversion_on_retirement'
)
const [ptdConversionOnWithdrawal, ptdConversionOnWithdrawalAttrs] = defineField(
  'ptd_conversion_on_withdrawal'
)
const [ciConversionOnWithdrawal, ciConversionOnWithdrawalAttrs] = defineField(
  'ci_conversion_on_withdrawal'
)
const [glaBenefit] = defineField('gla_benefit')
const [ptdBenefit] = defineField('ptd_benefit')
const [ciBenefit] = defineField('ci_benefit')
const [sglaBenefit] = defineField('sgla_benefit')
const [phiBenefit] = defineField('phi_benefit')
const [ttdBenefit] = defineField('ttd_benefit')
const [familyFuneralBenefit] = defineField('family_funeral_benefit')
const [ptdRiskType, ptdRiskTypeAttrs] = defineField('ptd_risk_type')
const [ptdBenefitType, ptdBenefitTypeAttrs] = defineField('ptd_benefit_type')
const [ptdSalaryMultiple, ptdSalaryMultipleAttrs] = defineField(
  'ptd_salary_multiple'
)
const [ptdDeferredPeriod, ptdDeferredPeriodAttrs] = defineField(
  'ptd_deferred_period'
)
const [ptdDisabilityDefinition, ptdDisabilityDefinitionAttrs] = defineField(
  'ptd_disability_definition'
)
const [ptdEducatorBenefit, ptdEducatorBenefitAttrs] = defineField(
  'ptd_educator_benefit'
)
const [ptdEducatorBenefitType, ptdEducatorBenefitTypeAttrs] = defineField(
  'ptd_educator_benefit_type'
)
const [ciBenefitStructure, ciBenefitStructureAttrs] = defineField(
  'ci_benefit_structure'
)
const [ciBenefitDefinition, ciBenefitDefinitionAttrs] = defineField(
  'ci_benefit_definition'
)
const [ciCriticalIllnessSalaryMultiple, ciCriticalIllnessSalaryMultipleAttrs] =
  defineField('ci_critical_illness_salary_multiple')
const [sglaSalaryMultiple, sglaSalaryMultipleAttrs] = defineField(
  'sgla_salary_multiple'
)
const [phiIncomeReplacementPercentage, phiIncomeReplacementPercentageAttrs] =
  defineField('phi_income_replacement_percentage')
const [
  phiUseTieredIncomeReplacementRatio,
  phiUseTieredIncomeReplacementRatioAttrs
] = defineField('phi_use_tiered_income_replacement_ratio')
const [phiTieredIncomeReplacementType, phiTieredIncomeReplacementTypeAttrs] =
  defineField('phi_tiered_income_replacement_type' as any)
const [phiPremiumWaiver, phiPremiumWaiverAttrs] =
  defineField('phi_premium_waiver')
const [phiMedicalAidPremiumWaiver, phiMedicalAidPremiumWaiverAttrs] =
  defineField('phi_medical_aid_premium_waiver')
const [phiBenefitEscalation, phiBenefitEscalationAttrs] = defineField(
  'phi_benefit_escalation'
)
const [phiWaitingPeriod, phiWaitingPeriodAttrs] =
  defineField('phi_waiting_period')
const [phiNormalRetirementAge, phiNormalRetirementAgeAttrs] = defineField(
  'phi_normal_retirement_age'
)
const [phiDeferredPeriod, phiDeferredPeriodAttrs] = defineField(
  'phi_deferred_period'
)
const [phiDisabilityDefinition, phiDisabilityDefinitionAttrs] = defineField(
  'phi_disability_definition'
)
const [phiRiskType, phiRiskTypeAttrs] = defineField('phi_risk_type')
const [ttdIncomeReplacementPercentage, ttdIncomeReplacementPercentageAttrs] =
  defineField('ttd_income_replacement_percentage')
const [
  ttdUseTieredIncomeReplacementRatio,
  ttdUseTieredIncomeReplacementRatioAttrs
] = defineField('ttd_use_tiered_income_replacement_ratio')
const [ttdTieredIncomeReplacementType, ttdTieredIncomeReplacementTypeAttrs] =
  defineField('ttd_tiered_income_replacement_type' as any)
// const [ttdPremiumWaiverPercentage, ttdPremiumWaiverPercentageAttrs] = defineField(
//   'ttd_premium_waiver_percentage'
// )

const [ttdWaitingPeriod, ttdWaitingPeriodAttrs] =
  defineField('ttd_waiting_period')
const [ttdDeferredPeriod, ttdDeferredPeriodAttrs] = defineField(
  'ttd_deferred_period'
)
const [ttdDisabilityDefinition, ttdDisabilityDefinitionAttrs] = defineField(
  'ttd_disability_definition'
)
const [ttdRiskType, ttdRiskTypeAttrs] = defineField('ttd_risk_type')
const [
  familyFuneralMainMemberFuneralSumAssured,
  familyFuneralMainMemberFuneralSumAssuredAttrs
] = defineField('family_funeral_main_member_funeral_sum_assured')
const [
  familyFuneralSpouseFuneralSumAssured,
  familyFuneralSpouseFuneralSumAssuredAttrs
] = defineField('family_funeral_spouse_funeral_sum_assured')
const [
  familyFuneralChildrenFuneralSumAssured,
  familyFuneralChildrenFuneralSumAssuredAttrs
] = defineField('family_funeral_children_funeral_sum_assured')
const [
  familyFuneralAdultDependantSumAssured,
  familyFuneralAdultDependantSumAssuredAttrs
] = defineField('family_funeral_adult_dependant_sum_assured')
const [
  familyFuneralParentFuneralSumAssured,
  familyFuneralParentFuneralSumAssuredAttrs
] = defineField('family_funeral_parent_funeral_sum_assured')
const [familyFuneralMaxNumberChildren, familyFuneralMaxNumberChildrenAttrs] =
  defineField('family_funeral_max_number_children')
const [
  familyFuneralMaxNumberAdultDependants,
  familyFuneralMaxNumberAdultDependantsAttrs
] = defineField('family_funeral_max_number_adult_dependants')

onBeforeMount(async () => {
  // get benefit definitions from the API
  const benefitResponse = await GroupPricingService.getBenefitDefinitions()
  benefitDefinitions.value = benefitResponse.data
  const res = await GroupPricingService.getBenefitMaps()
  benefitMaps.value = res.data
  glaLabel.value = getBenefitAlias('GLA')
  ptdLabel.value = getBenefitAlias('PTD')
  ciLabel.value = getBenefitAlias('CI')
  sglaLabel.value = getBenefitAlias('SGLA')
  phiLabel.value = getBenefitAlias('PHI')
  ttdLabel.value = getBenefitAlias('TTD')
  familyFuneralLabel.value = getBenefitAlias('GFF')
})

const getBenefitAlias = (benefit: any) => {
  const benefitMap = benefitMaps.value.find(
    (map: any) => map.benefit_code === benefit
  )
  return benefitMap.benefit_alias !== ''
    ? benefitMap.benefit_alias
    : benefitMap.benefit_name
}

// Snackbar state
const snackbar = ref({
  show: false,
  message: ''
})

// Utility: Get enabled benefits message
function getEnabledBenefitsMessage(benefitData: Record<string, any>): string {
  const benefitMap = [
    { key: 'gla_benefit', label: glaLabel.value },
    { key: 'ptd_benefit', label: ptdLabel.value },
    { key: 'ci_benefit', label: ciLabel.value },
    { key: 'sgla_benefit', label: sglaLabel.value },
    { key: 'phi_benefit', label: phiLabel.value },
    { key: 'ttd_benefit', label: ttdLabel.value },
    { key: 'family_funeral_benefit', label: familyFuneralLabel.value }
  ]
  const enabled = benefitMap
    .filter((b) => benefitData[b.key])
    .map((b) => b.label)
  return enabled.length ? enabled.join(', ') : 'None'
}

const validateForm = handleSubmit((values) => {
  return true
})

// Check custom tiered income replacement table existence when type changes to "custom"
const checkCustomTableExists = async (benefit: 'phi' | 'ttd') => {
  const schemeName = groupStore.group_pricing_quote.scheme_name
  const riskRateCode = groupStore.group_pricing_quote.risk_rate_code
  if (!schemeName || !riskRateCode) return
  try {
    const response = await GroupPricingService.checkCustomTieredTableExists(
      schemeName,
      riskRateCode
    )
    if (benefit === 'phi') {
      phiCustomTableExists.value = response.data.exists
    } else {
      ttdCustomTableExists.value = response.data.exists
    }
  } catch {
    if (benefit === 'phi') {
      phiCustomTableExists.value = false
    } else {
      ttdCustomTableExists.value = false
    }
  }
}

const requestCustomTable = async (benefit: 'phi' | 'ttd') => {
  const quote = groupStore.group_pricing_quote
  if (!quote.scheme_name || !quote.risk_rate_code) return
  try {
    await GroupPricingService.requestCustomTieredTable({
      scheme_name: quote.scheme_name,
      scheme_id: quote.scheme_id || 0,
      risk_rate_code: quote.risk_rate_code
    })
  } catch {
    // Notification request is best-effort
  }
}

watch(glaEducatorBenefit, (newVal) => {
  if (newVal !== 'Yes') {
    glaEducatorBenefitType.value = null
  }
})

watch(ptdEducatorBenefit, (newVal) => {
  if (newVal !== 'Yes') {
    ptdEducatorBenefitType.value = null
  }
})

watch(phiTieredIncomeReplacementType, (newVal, oldVal) => {
  if (newVal === 'custom') {
    checkCustomTableExists('phi')
    if (oldVal !== 'custom') {
      requestCustomTable('phi')
    }
  }
})

watch(ttdTieredIncomeReplacementType, (newVal, oldVal) => {
  if (newVal === 'custom') {
    checkCustomTableExists('ttd')
    if (oldVal !== 'custom') {
      requestCustomTable('ttd')
    }
  }
})

onMounted(() => {
  const rrc = groupStore.group_pricing_quote.risk_rate_code
  GroupPricingService.getWaitingPeriods('gla_rates', rrc).then((response) => {
    waitingPeriods.value = response.data
  })
  GroupPricingService.getGlaBenefitTypes(rrc).then((response) => {
    glaBenefitTypes.value = response.data
  })
  GroupPricingService.getEducatorBenefitTypes(rrc).then((response) => {
    educatorBenefitTypes.value = response.data
  })

  // Load income escalations from the store
  GroupPricingService.getBenefitEscalations(rrc).then((response) => {
    incomeEscalations.value = response.data
  })
  GroupPricingService.getPtdDisabilityDefinitions(rrc).then((response) => {
    ptdDisabilityDefinitions.value = response.data
  })
  GroupPricingService.getTtdDisabilityDefinitions(rrc).then((response) => {
    ttdDisabilityDefinitions.value = response.data
  })
  GroupPricingService.getPhiDisabilityDefinitions(rrc).then((response) => {
    phiDisabilityDefinitions.value = response.data
  })

  GroupPricingService.getWaitingPeriods('ci_rates', rrc).then((response) => {
    ciWaitingPeriods.value = response.data
  })
  GroupPricingService.getWaitingPeriods('ptd_rates', rrc).then((response) => {
    ptdWaitingPeriods.value = response.data
  })
  GroupPricingService.getWaitingPeriods('ttd_rates', rrc).then((response) => {
    ttdWaitingPeriods.value = response.data
  })
  GroupPricingService.getWaitingPeriods('phi_rates', rrc).then((response) => {
    phiWaitingPeriods.value = response.data
  })
  GroupPricingService.getDeferredPeriods('phi_rates', rrc).then((response) => {
    phiDeferredPeriods.value = response.data
  })
  GroupPricingService.getDeferredPeriods('ttd_rates', rrc).then((response) => {
    ttdDeferredPeriods.value = response.data
  })
  GroupPricingService.getDeferredPeriods('ptd_rates', rrc).then((response) => {
    ptdDeferredPeriods.value = response.data
  })
  GroupPricingService.getNormalRetirementAges().then((response) => {
    normalRetirementAges.value = response.data
  })
  GroupPricingService.getRiskTypes().then((response) => {
    ttdRiskTypes.value = response.data.ttd_rates
    phiRiskTypes.value = response.data.phi_rates
    ptdRiskTypes.value = response.data.ptd_rates
  })
  if (rrc) {
    GroupPricingService.getRegionsForRiskCode(rrc).then((response) => {
      availableRegions.value = response.data.data ?? []
    })
  }
})

const tabStatus = computed(() => {
  const currentErrors = Object.keys(errors.value)

  return {
    gla: currentErrors.some((k) => k.startsWith('gla')),
    ptd: currentErrors.some((k) => k.startsWith('ptd')),
    ci: currentErrors.some((k) => k.startsWith('ci')),
    sgla: currentErrors.some((k) => k.startsWith('sgla')),
    phi: currentErrors.some((k) => k.startsWith('phi')),
    ttd: currentErrors.some((k) => k.startsWith('ttd')),
    family_funeral: currentErrors.some((k) => k.startsWith('family_funeral'))
  }
})

defineExpose({
  validateForm
})

const emit = defineEmits<{
  (e: 'all_schemes_saved', schemeId: string, isSaved: boolean): void
}>()
</script>
<style lang="css" scoped>
.error-tab {
  color: #b00020 !important;
  border-bottom: 2px solid #b00020;
  font-weight: bold;
}
</style>
