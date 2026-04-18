<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center justify-space-between">
              <span class="headline">Reinsurance Treaty Management</span>
              <div class="d-flex gap-2">
                <v-btn
                  size="small"
                  rounded
                  color="primary"
                  prepend-icon="mdi-plus"
                  @click="openCreateDialog"
                >
                  New Treaty
                </v-btn>
                <v-btn
                  size="small"
                  rounded
                  color="white"
                  variant="outlined"
                  prepend-icon="mdi-arrow-left"
                  @click="$router.push('/group-pricing/bordereaux-management')"
                >
                  Back
                </v-btn>
              </div>
            </div>
          </template>

          <template #default>
            <!-- Stats chips -->
            <v-row class="mb-4">
              <v-col cols="auto">
                <v-chip
                  color="primary"
                  variant="tonal"
                  prepend-icon="mdi-file-sign"
                >
                  Total: {{ stats.total }}
                </v-chip>
              </v-col>
              <v-col cols="auto">
                <v-chip
                  color="success"
                  variant="tonal"
                  prepend-icon="mdi-check-circle"
                >
                  Active: {{ stats.active }}
                </v-chip>
              </v-col>
              <v-col cols="auto">
                <v-chip color="grey" variant="tonal" prepend-icon="mdi-pencil">
                  Draft: {{ stats.draft }}
                </v-chip>
              </v-col>
              <v-col cols="auto">
                <v-chip
                  color="error"
                  variant="tonal"
                  prepend-icon="mdi-clock-alert"
                >
                  Expired: {{ stats.expired }}
                </v-chip>
              </v-col>
              <v-col cols="auto">
                <v-chip
                  color="warning"
                  variant="tonal"
                  prepend-icon="mdi-alert"
                >
                  Expiring &lt;60 days: {{ stats.expiring_in_60_days }}
                </v-chip>
              </v-col>
              <v-col cols="auto">
                <v-chip
                  color="orange"
                  variant="tonal"
                  prepend-icon="mdi-archive-clock"
                  style="cursor: pointer"
                  @click="
                    () => {
                      filters.run_off_only = !filters.run_off_only
                      loadTreaties()
                    }
                  "
                >
                  Run-Off: {{ runOffCount }}
                </v-chip>
              </v-col>
            </v-row>

            <!-- Filters -->
            <v-card variant="outlined" class="mb-4">
              <v-card-text>
                <v-row>
                  <v-col cols="12" sm="4" md="3">
                    <v-select
                      v-model="filters.status"
                      :items="statusOptions"
                      label="Status"
                      variant="outlined"
                      density="compact"
                      clearable
                      @update:model-value="loadTreaties"
                    />
                  </v-col>
                  <v-col cols="12" sm="4" md="3">
                    <v-select
                      v-model="filters.type"
                      :items="typeOptions"
                      label="Treaty Type"
                      variant="outlined"
                      density="compact"
                      clearable
                      @update:model-value="loadTreaties"
                    />
                  </v-col>
                  <v-col cols="12" sm="4" md="3">
                    <v-text-field
                      v-model="filters.reinsurer"
                      label="Reinsurer"
                      variant="outlined"
                      density="compact"
                      clearable
                      prepend-inner-icon="mdi-magnify"
                      @update:model-value="loadTreaties"
                    />
                  </v-col>
                  <v-col cols="12" sm="2" class="d-flex align-center">
                    <v-switch
                      v-model="filters.run_off_only"
                      label="Run-Off Only"
                      density="compact"
                      color="warning"
                      hide-details
                      @update:model-value="loadTreaties"
                    />
                  </v-col>
                  <v-col cols="auto" class="d-flex align-center">
                    <v-btn
                      size="small"
                      rounded
                      color="success"
                      prepend-icon="mdi-refresh"
                      :loading="loading"
                      @click="loadTreaties"
                    >
                      Refresh
                    </v-btn>
                  </v-col>
                </v-row>
              </v-card-text>
            </v-card>

            <!-- AG Grid -->
            <data-grid
              class="ag-theme-material"
              style="height: 100%; width: 100%"
              :row-data="treaties"
              :column-defs="columnDefs"
              :default-col-def="{
                resizable: true,
                sortable: true,
                filter: true
              }"
              :pagination="true"
              :pagination-page-size="20"
              @grid-ready="onGridReady"
            />
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Create/Edit Dialog -->
    <v-dialog v-model="showFormDialog" max-width="800" scrollable>
      <v-card>
        <v-card-title>{{
          editingTreaty ? 'Edit Treaty' : 'New Reinsurance Treaty'
        }}</v-card-title>
        <v-card-text>
          <v-form ref="treatyForm">
            <v-row>
              <v-col cols="12" sm="6">
                <v-text-field
                  v-model="form.treaty_number"
                  label="Treaty Number *"
                  variant="outlined"
                  density="compact"
                  :disabled="!!editingTreaty"
                  :rules="rules.required"
                />
              </v-col>
              <v-col cols="12" sm="6">
                <v-text-field
                  v-model="form.treaty_name"
                  label="Treaty Name *"
                  variant="outlined"
                  density="compact"
                  :rules="rules.required"
                />
              </v-col>
              <v-col cols="12" sm="6">
                <v-autocomplete
                  v-model="form.reinsurer_name"
                  :items="reinsurers"
                  item-title="name"
                  item-value="name"
                  label="Reinsurer Name *"
                  variant="outlined"
                  density="compact"
                  :rules="rules.required"
                  clearable
                  no-data-text="No reinsurers found"
                  @update:model-value="onReinsurerSelected"
                />
              </v-col>
              <v-col cols="12" sm="6">
                <v-text-field
                  v-model="form.reinsurer_code"
                  label="Reinsurer Code"
                  variant="outlined"
                  density="compact"
                  :disabled="!!form.reinsurer_name"
                  hint="Auto-filled when reinsurer is selected"
                  persistent-hint
                />
              </v-col>
              <v-col cols="12" sm="6">
                <v-autocomplete
                  v-model="form.broker_name"
                  :items="brokers"
                  item-title="name"
                  item-value="name"
                  label="Broker"
                  variant="outlined"
                  density="compact"
                  clearable
                  no-data-text="No brokers found"
                />
              </v-col>
              <v-col cols="12" sm="6">
                <v-select
                  v-model="form.treaty_type"
                  :items="typeOptions"
                  label="Treaty Type *"
                  variant="outlined"
                  density="compact"
                  :rules="rules.required"
                />
              </v-col>
              <v-col cols="12" sm="6">
                <v-select
                  v-model="form.line_of_business"
                  :items="lobOptions"
                  label="Line of Business"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" sm="6">
                <v-select
                  v-model="form.status"
                  :items="statusOptions"
                  label="Status"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model="form.effective_date"
                  label="Effective Date *"
                  type="date"
                  variant="outlined"
                  density="compact"
                  :rules="rules.required"
                />
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model="form.expiry_date"
                  label="Expiry Date *"
                  type="date"
                  variant="outlined"
                  density="compact"
                  :rules="rules.expiryDate"
                />
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model="form.renewal_date"
                  label="Renewal Date"
                  type="date"
                  variant="outlined"
                  density="compact"
                  :rules="rules.renewalDate"
                />
              </v-col>

              <!-- XL fields -->
              <template
                v-if="
                  form.treaty_type === 'xl_risk' ||
                  form.treaty_type === 'xl_event'
                "
              >
                <v-col cols="12" sm="6">
                  <v-text-field
                    v-model:model-value="formattedXlRetention"
                    label="XL Retention (ZAR) *"
                    type="text"
                    variant="outlined"
                    density="compact"
                    prefix="R"
                    :rules="rules.positiveRequired"
                  />
                </v-col>
                <v-col cols="12" sm="6">
                  <v-text-field
                    v-model:model-value="formattedXlLimit"
                    label="XL Limit (ZAR) *"
                    type="text"
                    variant="outlined"
                    density="compact"
                    prefix="R"
                    :rules="rules.xlLimit"
                  />
                </v-col>
              </template>

              <!-- Stop Loss -->
              <template
                v-if="
                  form.treaty_type === 'stop_loss' ||
                  form.treaty_type === 'catastrophe_xl'
                "
              >
                <v-col cols="12" sm="6">
                  <v-text-field
                    v-model:model-value="formattedAggregateAnnualLimit"
                    label="Aggregate Annual Limit (ZAR) *"
                    type="text"
                    variant="outlined"
                    density="compact"
                    prefix="R"
                    :rules="rules.positiveRequired"
                  />
                </v-col>
              </template>

              <v-col cols="12" sm="4">
                <v-text-field
                  v-model.number="form.reinsurance_commission_rate"
                  label="RI Commission %"
                  type="number"
                  variant="outlined"
                  density="compact"
                  suffix="%"
                  :rules="rules.percentage"
                />
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model.number="form.profit_commission_rate"
                  label="Profit Commission %"
                  type="number"
                  variant="outlined"
                  density="compact"
                  suffix="%"
                  :rules="rules.percentage"
                />
              </v-col>
              <v-col cols="12" sm="4">
                <v-select
                  v-model="form.premium_payment_frequency"
                  :items="['monthly', 'quarterly', 'annual']"
                  label="Payment Frequency"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model.number="form.claims_notification_days"
                  label="Claims Notification SLA (days)"
                  type="number"
                  variant="outlined"
                  density="compact"
                  :rules="rules.positiveDays"
                />
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model:model-value="formattedLargeClaimsThreshold"
                  label="Large Claims Threshold (ZAR)"
                  type="text"
                  variant="outlined"
                  density="compact"
                  prefix="R"
                  :rules="rules.nonNegative"
                />
              </v-col>
              <v-col cols="12" sm="4" class="d-flex align-center">
                <v-checkbox
                  v-model="form.delta_reporting"
                  label="Delta Reporting (changes only)"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" sm="4" class="d-flex align-center">
                <v-checkbox
                  v-model="form.is_run_off"
                  label="Run-Off Treaty"
                  density="compact"
                />
              </v-col>
              <v-col v-if="form.is_run_off" cols="12" sm="4">
                <v-text-field
                  v-model="form.run_off_start_date"
                  label="Run-Off Start Date *"
                  type="date"
                  variant="outlined"
                  density="compact"
                  :rules="rules.required"
                />
              </v-col>

              <!-- Three-Tier Reinsurance Structure -->
              <v-col cols="12">
                <v-divider class="my-2" />
                <p class="text-subtitle-2 font-weight-bold mb-2"
                  >Three-Tier Reinsurance Structure</p
                >
              </v-col>

              <v-col cols="12" sm="4">
                <v-text-field
                  v-model="form.treaty_code"
                  label="Treaty Code"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" sm="4">
                <v-select
                  v-model="form.risk_premium_basis_indicator"
                  :items="['default', 'flat_rate', 'custom']"
                  label="Risk Premium Basis"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model.number="form.flat_annual_reins_prem_rate"
                  label="Flat Annual Reins Prem Rate"
                  type="number"
                  variant="outlined"
                  density="compact"
                  step="0.01"
                  :disabled="form.risk_premium_basis_indicator !== 'flat_rate'"
                />
              </v-col>

              <!-- Sum Assured Level 1 -->
              <v-col cols="12">
                <p class="text-subtitle-2 font-weight-bold mb-2"
                  >Sum Assured Tiers</p
                >
              </v-col>
              <v-col cols="12">
                <p class="text-caption text-medium-emphasis"
                  >Level 1 - Sum Assured Tier (Required)</p
                >
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model:model-value="formattedLevel1Lowerbound"
                  label="Level 1 Lower Bound *"
                  type="text"
                  variant="outlined"
                  density="compact"
                  prefix="R"
                  :rules="rules.required"
                />
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model:model-value="formattedLevel1Upperbound"
                  label="Level 1 Upper Bound *"
                  type="text"
                  variant="outlined"
                  density="compact"
                  prefix="R"
                  :rules="rules.level1UpperBound"
                />
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model.number="form.level1_ceded_proportion"
                  label="Level 1 Ceded % *"
                  type="number"
                  variant="outlined"
                  density="compact"
                  suffix="%"
                  :rules="rules.required"
                />
              </v-col>

              <!-- Add Level 2 Button -->
              <v-col v-if="!showSumAssuredLevel2" cols="12">
                <v-btn
                  size="small"
                  variant="outlined"
                  color="primary"
                  prepend-icon="mdi-plus"
                  @click="showSumAssuredIncomeLevel2"
                >
                  Add Sum Assured Level 2
                </v-btn>
              </v-col>

              <!-- Sum Assured Level 2 -->
              <template v-if="showSumAssuredLevel2">
                <v-col
                  cols="12"
                  class="d-flex align-center justify-space-between"
                >
                  <p class="text-caption text-medium-emphasis mb-0"
                    >Level 2 - Sum Assured Tier (Optional)</p
                  >
                  <v-btn
                    size="x-small"
                    variant="text"
                    color="error"
                    icon="mdi-close"
                    @click="showLevel2And3"
                  />
                </v-col>
                <v-col cols="12" sm="4">
                  <v-text-field
                    v-model:model-value="formattedLevel2Lowerbound"
                    label="Level 2 Lower Bound"
                    type="text"
                    variant="outlined"
                    density="compact"
                    prefix="R"
                    readonly
                    hint="Auto-calculated: Level 1 Upper + 1"
                    persistent-hint
                  />
                </v-col>
                <v-col cols="12" sm="4">
                  <v-text-field
                    v-model:model-value="formattedLevel2Upperbound"
                    label="Level 2 Upper Bound"
                    type="text"
                    variant="outlined"
                    density="compact"
                    prefix="R"
                    :rules="rules.level2UpperBound"
                  />
                </v-col>
                <v-col cols="12" sm="4">
                  <v-text-field
                    v-model.number="form.level2_ceded_proportion"
                    label="Level 2 Ceded %"
                    type="number"
                    variant="outlined"
                    density="compact"
                    suffix="%"
                  />
                </v-col>

                <!-- Add Level 3 Button -->
                <v-col v-if="!showSumAssuredLevel3" cols="12">
                  <v-btn
                    size="small"
                    variant="outlined"
                    color="primary"
                    prepend-icon="mdi-plus"
                    @click="showLevel3"
                  >
                    Add Sum Assured Level 3
                  </v-btn>
                </v-col>
              </template>

              <!-- Sum Assured Level 3 -->
              <template v-if="showSumAssuredLevel3">
                <v-col
                  cols="12"
                  class="d-flex align-center justify-space-between"
                >
                  <p class="text-caption text-medium-emphasis mb-0"
                    >Level 3 - Sum Assured Tier (Optional)</p
                  >
                  <v-btn
                    size="x-small"
                    variant="text"
                    color="error"
                    icon="mdi-close"
                    @click="showLevel3"
                  />
                </v-col>
                <v-col cols="12" sm="4">
                  <v-text-field
                    v-model:model-value="formattedLevel3Lowerbound"
                    label="Level 3 Lower Bound"
                    type="text"
                    variant="outlined"
                    density="compact"
                    prefix="R"
                    readonly
                    hint="Auto-calculated: Level 2 Upper + 1"
                    persistent-hint
                  />
                </v-col>
                <v-col cols="12" sm="4">
                  <v-text-field
                    v-model:model-value="formattedLevel3Upperbound"
                    label="Level 3 Upper Bound"
                    type="text"
                    variant="outlined"
                    density="compact"
                    prefix="R"
                    :rules="rules.level3UpperBound"
                  />
                </v-col>
                <v-col cols="12" sm="4">
                  <v-text-field
                    v-model.number="form.level3_ceded_proportion"
                    label="Level 3 Ceded %"
                    type="number"
                    variant="outlined"
                    density="compact"
                    suffix="%"
                  />
                </v-col>
              </template>

              <!-- Income Level 1 -->
              <v-col cols="12">
                <v-divider class="my-2" />
                <p class="text-subtitle-2 font-weight-bold mb-2"
                  >Income Tiers (monthly income)</p
                >
              </v-col>
              <v-col cols="12">
                <p class="text-caption text-medium-emphasis"
                  >Income Level 1 - Income Tier (monthly income) (Required)</p
                >
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model:model-value="formattedIncomeLevel1Lowerbound"
                  label="Income Level 1 Lower Bound *"
                  type="text"
                  variant="outlined"
                  density="compact"
                  prefix="R"
                  :rules="rules.required"
                />
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model:model-value="formattedIncomeLevel1Upperbound"
                  label="Income Level 1 Upper Bound *"
                  type="text"
                  variant="outlined"
                  density="compact"
                  prefix="R"
                  :rules="rules.incomeLevel1UpperBound"
                />
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model.number="form.income_level1_ceded_proportion"
                  label="Income Level 1 Ceded % *"
                  type="number"
                  variant="outlined"
                  density="compact"
                  suffix="%"
                  :rules="rules.required"
                />
              </v-col>

              <!-- Income Level 2 automatically shown with Sum Assured Level 2 -->

              <!-- Income Level 2 -->
              <template v-if="showIncomeLevel2">
                <v-col
                  cols="12"
                  class="d-flex align-center justify-space-between"
                >
                  <p class="text-caption text-medium-emphasis mb-0"
                    >Income Level 2 - Income Tier (monthly income) (Optional)</p
                  >
                  <v-btn
                    size="x-small"
                    variant="text"
                    color="error"
                    icon="mdi-close"
                    @click="hideLevel2And3"
                  />
                </v-col>
                <v-col cols="12" sm="4">
                  <v-text-field
                    v-model:model-value="formattedIncomeLevel2Lowerbound"
                    label="Income Level 2 Lower Bound"
                    type="text"
                    variant="outlined"
                    density="compact"
                    prefix="R"
                    readonly
                    hint="Auto-calculated: Income Level 1 Upper + 1"
                    persistent-hint
                  />
                </v-col>
                <v-col cols="12" sm="4">
                  <v-text-field
                    v-model:model-value="formattedIncomeLevel2Upperbound"
                    label="Income Level 2 Upper Bound"
                    type="text"
                    variant="outlined"
                    density="compact"
                    prefix="R"
                    :rules="rules.incomeLevel2UpperBound"
                  />
                </v-col>
                <v-col cols="12" sm="4">
                  <v-text-field
                    v-model.number="form.income_level2_ceded_proportion"
                    label="Income Level 2 Ceded %"
                    type="number"
                    variant="outlined"
                    density="compact"
                    suffix="%"
                  />
                </v-col>

                <!-- Income Level 3 automatically shown with Sum Assured Level 3 -->
              </template>

              <!-- Income Level 3 -->
              <template v-if="showIncomeLevel3">
                <v-col
                  cols="12"
                  class="d-flex align-center justify-space-between"
                >
                  <p class="text-caption text-medium-emphasis mb-0"
                    >Income Level 3 - Income Tier (monthly income) (Optional)</p
                  >
                  <v-btn
                    size="x-small"
                    variant="text"
                    color="error"
                    icon="mdi-close"
                    @click="hideLevel3"
                  />
                </v-col>
                <v-col cols="12" sm="4">
                  <v-text-field
                    v-model:model-value="formattedIncomeLevel3Lowerbound"
                    label="Income Level 3 Lower Bound"
                    type="text"
                    variant="outlined"
                    density="compact"
                    prefix="R"
                    readonly
                    hint="Auto-calculated: Income Level 2 Upper + 1"
                    persistent-hint
                  />
                </v-col>
                <v-col cols="12" sm="4">
                  <v-text-field
                    v-model:model-value="formattedIncomeLevel3Upperbound"
                    label="Income Level 3 Upper Bound"
                    type="text"
                    variant="outlined"
                    density="compact"
                    prefix="R"
                    :rules="rules.incomeLevel3UpperBound"
                  />
                </v-col>
                <v-col cols="12" sm="4">
                  <v-text-field
                    v-model.number="form.income_level3_ceded_proportion"
                    label="Income Level 3 Ceded %"
                    type="number"
                    variant="outlined"
                    density="compact"
                    suffix="%"
                  />
                </v-col>
              </template>

              <!-- Multi-Reinsurer Structure -->
              <v-col cols="12">
                <v-divider class="my-2" />
                <p class="text-subtitle-2 font-weight-bold mb-2"
                  >Multi-Reinsurer Structure</p
                >
              </v-col>

              <v-col cols="12" sm="4">
                <v-autocomplete
                  v-model="form.lead_reinsurer_code"
                  :items="reinsurers"
                  item-title="name"
                  item-value="code"
                  label="Lead Reinsurer"
                  variant="outlined"
                  density="compact"
                  clearable
                  no-data-text="No reinsurers found"
                />
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model.number="form.lead_reinsurer_share"
                  label="Lead Reinsurer Share %"
                  type="number"
                  variant="outlined"
                  density="compact"
                  suffix="%"
                />
              </v-col>
              <v-col cols="12" sm="4">
                <v-text-field
                  v-model.number="form.ceding_commission"
                  label="Ceding Commission %"
                  type="number"
                  variant="outlined"
                  density="compact"
                  suffix="%"
                />
              </v-col>

              <!-- Non-Lead Reinsurers (Expandable) -->
              <v-col cols="12">
                <v-expansion-panels variant="accordion">
                  <v-expansion-panel>
                    <v-expansion-panel-title>
                      <span class="text-body-2">
                        <v-icon size="small" class="mr-1"
                          >mdi-account-multiple</v-icon
                        >
                        Non-Lead Reinsurers (Optional)
                      </span>
                    </v-expansion-panel-title>
                    <v-expansion-panel-text>
                      <v-row>
                        <v-col cols="12" sm="6">
                          <v-autocomplete
                            v-model="form.non_lead_reinsurer1_code"
                            :items="reinsurers"
                            item-title="name"
                            item-value="code"
                            label="Non-Lead Reinsurer 1"
                            variant="outlined"
                            density="compact"
                            clearable
                            no-data-text="No reinsurers found"
                          />
                        </v-col>
                        <v-col cols="12" sm="6">
                          <v-text-field
                            v-model.number="form.non_lead_reinsurer1_share"
                            label="Non-Lead Reinsurer 1 Share %"
                            type="number"
                            variant="outlined"
                            density="compact"
                            suffix="%"
                          />
                        </v-col>

                        <v-col cols="12" sm="6">
                          <v-autocomplete
                            v-model="form.non_lead_reinsurer2_code"
                            :items="reinsurers"
                            item-title="name"
                            item-value="code"
                            label="Non-Lead Reinsurer 2"
                            variant="outlined"
                            density="compact"
                            clearable
                            no-data-text="No reinsurers found"
                          />
                        </v-col>
                        <v-col cols="12" sm="6">
                          <v-text-field
                            v-model.number="form.non_lead_reinsurer2_share"
                            label="Non-Lead Reinsurer 2 Share %"
                            type="number"
                            variant="outlined"
                            density="compact"
                            suffix="%"
                          />
                        </v-col>

                        <v-col cols="12" sm="6">
                          <v-autocomplete
                            v-model="form.non_lead_reinsurer3_code"
                            :items="reinsurers"
                            item-title="name"
                            item-value="code"
                            label="Non-Lead Reinsurer 3"
                            variant="outlined"
                            density="compact"
                            clearable
                            no-data-text="No reinsurers found"
                          />
                        </v-col>
                        <v-col cols="12" sm="6">
                          <v-text-field
                            v-model.number="form.non_lead_reinsurer3_share"
                            label="Non-Lead Reinsurer 3 Share %"
                            type="number"
                            variant="outlined"
                            density="compact"
                            suffix="%"
                          />
                        </v-col>
                      </v-row>
                    </v-expansion-panel-text>
                  </v-expansion-panel>
                </v-expansion-panels>
              </v-col>

              <v-col cols="12">
                <v-textarea
                  v-model="form.notes"
                  label="Notes"
                  variant="outlined"
                  density="compact"
                  rows="2"
                />
              </v-col>
            </v-row>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showFormDialog = false">Cancel</v-btn>
          <v-btn color="primary" :loading="saving" @click="saveTreaty">
            {{ editingTreaty ? 'Update' : 'Create' }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Scheme Links Dialog -->
    <v-dialog v-model="showSchemesDialog" max-width="1100" scrollable>
      <v-card>
        <v-card-title class="d-flex align-center">
          Linked Schemes — {{ selectedTreaty?.treaty_name }}
          <v-chip class="ml-3" size="small" color="primary" variant="tonal">
            {{ schemeLinks.length }} linked
          </v-chip>
          <v-spacer />
          <v-btn
            v-if="selectedLinkRows.length"
            color="error"
            size="small"
            variant="tonal"
            class="mr-2"
            @click="confirmBulkUnlink"
          >
            <v-icon start>mdi-link-off</v-icon>
            Unlink {{ selectedLinkRows.length }} Selected
          </v-btn>
        </v-card-title>
        <v-card-text>
          <DataGrid
            :column-defs="schemeLinkCols"
            :row-data="schemeLinks"
            :pagination="true"
            :pagination-page-size="50"
            row-selection="multiple"
            density="compact"
            @row-selection-changed="onLinkRowsSelected"
          />

          <v-divider class="my-4" />
          <p class="text-subtitle-2 mb-2">Add Schemes</p>
          <v-row>
            <v-col cols="12">
              <v-autocomplete
                v-model="selectedSchemes"
                :items="unlinkedSchemes"
                :item-title="(s) => `${s.name} (ID: ${s.id})`"
                item-value="id"
                label="Search and select schemes *"
                variant="outlined"
                density="compact"
                clearable
                multiple
                chips
                closable-chips
                no-data-text="No unlinked schemes found"
                hint="Only schemes not already linked are shown"
                persistent-hint
              />
            </v-col>
            <v-col cols="4">
              <v-text-field
                v-model.number="linkForm.cession_override"
                label="Cession Override %"
                type="number"
                variant="outlined"
                density="compact"
                suffix="%"
              />
            </v-col>
            <v-col cols="4">
              <v-text-field
                v-model="linkForm.effective_date"
                label="Effective Date"
                type="date"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="4" class="d-flex align-center">
              <v-btn
                color="primary"
                :loading="savingLink"
                :disabled="!selectedSchemes.length"
                @click="addSchemeLinks"
              >
                Link
                {{
                  selectedSchemes.length > 1
                    ? selectedSchemes.length + ' Schemes'
                    : 'Scheme'
                }}
              </v-btn>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showSchemesDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Bulk Unlink Confirm -->
    <v-dialog v-model="showBulkUnlinkDialog" max-width="460">
      <v-card>
        <v-card-title>Remove Scheme Links</v-card-title>
        <v-card-text>
          Remove <strong>{{ selectedLinkRows.length }}</strong> selected
          scheme(s) from this treaty? This cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showBulkUnlinkDialog = false"
            >Cancel</v-btn
          >
          <v-btn
            color="error"
            :loading="deletingLinkLoading"
            @click="executeBulkUnlink"
            >Remove</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete Confirm -->
    <v-dialog v-model="showDeleteDialog" max-width="420">
      <v-card>
        <v-card-title>Delete Treaty</v-card-title>
        <v-card-text>
          Are you sure you want to delete
          <strong>{{ deletingTreaty?.treaty_name }}</strong
          >? This cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showDeleteDialog = false">Cancel</v-btn>
          <v-btn color="error" :loading="deleting" @click="confirmDelete"
            >Delete</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3500">
      {{ snackbar.message }}
    </v-snackbar>
  </v-container>
</template>

<script setup>
import { ref, computed, nextTick, onMounted, watch } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import BaseCard from '@/renderer/components/BaseCard.vue'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import { useBordereauxStore } from '@/renderer/store/bordereaux'

const bordereauxStore = useBordereauxStore()

const treaties = ref([])
const brokers = ref([])
const reinsurers = ref([])
const runOffCount = computed(
  () => treaties.value.filter((t) => t.is_run_off).length
)
const treatyForm = ref(null)
const loading = ref(false)
const saving = ref(false)
const savingLink = ref(false)
const deleting = ref(false)
const deletingLinkLoading = ref(false)
const stats = ref({
  total: 0,
  active: 0,
  draft: 0,
  expired: 0,
  expiring_in_60_days: 0
})
const schemeLinks = ref([])
const filters = ref({
  status: '',
  type: '',
  reinsurer: '',
  run_off_only: false
})
const showFormDialog = ref(false)
const showSchemesDialog = ref(false)
const showDeleteDialog = ref(false)
const showBulkUnlinkDialog = ref(false)
const editingTreaty = ref(null)
const selectedTreaty = ref(null)
const deletingTreaty = ref(null)
const selectedLinkRows = ref([])
const snackbar = ref({ show: false, message: '', color: 'success' })

// Visibility state for optional tier levels
const showSumAssuredLevel2 = ref(false)
const showSumAssuredLevel3 = ref(false)
const showIncomeLevel2 = ref(false)
const showIncomeLevel3 = ref(false)

const showLevel2And3 = () => {
  showSumAssuredLevel2.value = true
  showSumAssuredLevel3.value = true
  showIncomeLevel2.value = true
  showIncomeLevel3.value = true
}

const showLevel3 = () => {
  showSumAssuredLevel3.value = true
  showIncomeLevel3.value = true
}

const hideLevel2And3 = () => {
  showSumAssuredLevel2.value = false
  showSumAssuredLevel3.value = false
  showIncomeLevel2.value = false
  showIncomeLevel3.value = false
}

const hideLevel3 = () => {
  showSumAssuredLevel3.value = false
  showIncomeLevel3.value = false
}

const defaultForm = () => ({
  treaty_number: '',
  treaty_name: '',
  reinsurer_name: '',
  reinsurer_code: '',
  broker_name: '',
  treaty_type: 'proportional',
  line_of_business: 'group_life',
  effective_date: '',
  expiry_date: '',
  renewal_date: '',
  status: 'draft',
  currency: 'ZAR',
  retention_type: 'percentage',
  retention_amount: 0,
  retention_percentage: 0,
  surplus_lines: 0,
  xl_retention: 0,
  xl_limit: 0,
  xl_layer_from: 0,
  xl_layer_to: 0,
  aggregate_annual_limit: 0,
  profit_commission_rate: 0,
  reinsurance_commission_rate: 0,
  premium_payment_frequency: 'monthly',
  claims_notification_days: 30,
  large_claims_threshold: 0,
  delta_reporting: false,
  is_run_off: false,
  run_off_start_date: '',
  notes: '',
  // Three-Tier Reinsurance Structure
  treaty_code: '',
  risk_premium_basis_indicator: 'default',
  flat_annual_reins_prem_rate: 0,
  level1_ceded_proportion: 0,
  level1_lowerbound: 0,
  level1_upperbound: 0,
  level2_ceded_proportion: 0,
  level2_lowerbound: 0,
  level2_upperbound: 0,
  level3_ceded_proportion: 0,
  level3_lowerbound: 0,
  level3_upperbound: 0,
  income_level1_ceded_proportion: 0,
  income_level1_lowerbound: 0,
  income_level1_upperbound: 0,
  income_level2_ceded_proportion: 0,
  income_level2_lowerbound: 0,
  income_level2_upperbound: 0,
  income_level3_ceded_proportion: 0,
  income_level3_lowerbound: 0,
  income_level3_upperbound: 0,
  lead_reinsurer_share: 0,
  lead_reinsurer_code: '',
  non_lead_reinsurer1_share: 0,
  non_lead_reinsurer1_code: '',
  non_lead_reinsurer2_share: 0,
  non_lead_reinsurer2_code: '',
  non_lead_reinsurer3_share: 0,
  non_lead_reinsurer3_code: '',
  ceding_commission: 0
})

const form = ref(defaultForm())
const linkForm = ref({ cession_override: 0, effective_date: '' })
const availableSchemes = ref([])
const selectedSchemes = ref([])

const schemeLinkCols = [
  {
    headerName: 'Scheme ID',
    field: 'scheme_id',
    sortable: true,
    filter: true,
    minWidth: 110
  },
  {
    headerName: 'Scheme Name',
    field: 'scheme_name',
    sortable: true,
    filter: true,
    flex: 1,
    minWidth: 200
  },
  {
    headerName: 'Cession Override %',
    field: 'cession_override',
    sortable: true,
    filter: true,
    minWidth: 160,
    valueFormatter: (p) => p.value || '—'
  },
  {
    headerName: 'Effective Date',
    field: 'effective_date',
    sortable: true,
    filter: true,
    minWidth: 140,
    valueFormatter: (p) => p.value || '—'
  }
]

const unlinkedSchemes = computed(() => {
  const linkedIds = new Set(schemeLinks.value.map((l) => l.scheme_id))
  return availableSchemes.value.filter((s) => !linkedIds.has(s.id))
})

// Formatted currency fields with accounting format (thousand separators)
const formattedXlRetention = computed({
  get: () => {
    if (typeof form.value.xl_retention === 'number') {
      return form.value.xl_retention.toLocaleString()
    }
    return form.value.xl_retention
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.xl_retention = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.xl_retention = 0
    }
  }
})

const showSumAssuredIncomeLevel2 = () => {
  showSumAssuredLevel2.value = true
  showIncomeLevel2.value = true
}

const formattedXlLimit = computed({
  get: () => {
    if (typeof form.value.xl_limit === 'number') {
      return form.value.xl_limit.toLocaleString()
    }
    return form.value.xl_limit
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.xl_limit = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.xl_limit = 0
    }
  }
})

const formattedAggregateAnnualLimit = computed({
  get: () => {
    if (typeof form.value.aggregate_annual_limit === 'number') {
      return form.value.aggregate_annual_limit.toLocaleString()
    }
    return form.value.aggregate_annual_limit
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.aggregate_annual_limit = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.aggregate_annual_limit = 0
    }
  }
})

const formattedLargeClaimsThreshold = computed({
  get: () => {
    if (typeof form.value.large_claims_threshold === 'number') {
      return form.value.large_claims_threshold.toLocaleString()
    }
    return form.value.large_claims_threshold
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.large_claims_threshold = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.large_claims_threshold = 0
    }
  }
})

// Sum Assured Level bounds
const formattedLevel1Lowerbound = computed({
  get: () => {
    if (typeof form.value.level1_lowerbound === 'number') {
      return form.value.level1_lowerbound.toLocaleString()
    }
    return form.value.level1_lowerbound
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.level1_lowerbound = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.level1_lowerbound = 0
    }
  }
})

const formattedLevel1Upperbound = computed({
  get: () => {
    if (typeof form.value.level1_upperbound === 'number') {
      return form.value.level1_upperbound.toLocaleString()
    }
    return form.value.level1_upperbound
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.level1_upperbound = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.level1_upperbound = 0
    }
  }
})

const formattedLevel2Lowerbound = computed({
  get: () => {
    if (typeof form.value.level2_lowerbound === 'number') {
      return form.value.level2_lowerbound.toLocaleString()
    }
    return form.value.level2_lowerbound
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.level2_lowerbound = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.level2_lowerbound = 0
    }
  }
})

const formattedLevel2Upperbound = computed({
  get: () => {
    if (typeof form.value.level2_upperbound === 'number') {
      return form.value.level2_upperbound.toLocaleString()
    }
    return form.value.level2_upperbound
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.level2_upperbound = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.level2_upperbound = 0
    }
  }
})

const formattedLevel3Lowerbound = computed({
  get: () => {
    if (typeof form.value.level3_lowerbound === 'number') {
      return form.value.level3_lowerbound.toLocaleString()
    }
    return form.value.level3_lowerbound
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.level3_lowerbound = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.level3_lowerbound = 0
    }
  }
})

const formattedLevel3Upperbound = computed({
  get: () => {
    if (typeof form.value.level3_upperbound === 'number') {
      return form.value.level3_upperbound.toLocaleString()
    }
    return form.value.level3_upperbound
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.level3_upperbound = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.level3_upperbound = 0
    }
  }
})

// Income Level bounds
const formattedIncomeLevel1Lowerbound = computed({
  get: () => {
    if (typeof form.value.income_level1_lowerbound === 'number') {
      return form.value.income_level1_lowerbound.toLocaleString()
    }
    return form.value.income_level1_lowerbound
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.income_level1_lowerbound = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.income_level1_lowerbound = 0
    }
  }
})

const formattedIncomeLevel1Upperbound = computed({
  get: () => {
    if (typeof form.value.income_level1_upperbound === 'number') {
      return form.value.income_level1_upperbound.toLocaleString()
    }
    return form.value.income_level1_upperbound
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.income_level1_upperbound = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.income_level1_upperbound = 0
    }
  }
})

const formattedIncomeLevel2Lowerbound = computed({
  get: () => {
    if (typeof form.value.income_level2_lowerbound === 'number') {
      return form.value.income_level2_lowerbound.toLocaleString()
    }
    return form.value.income_level2_lowerbound
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.income_level2_lowerbound = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.income_level2_lowerbound = 0
    }
  }
})

const formattedIncomeLevel2Upperbound = computed({
  get: () => {
    if (typeof form.value.income_level2_upperbound === 'number') {
      return form.value.income_level2_upperbound.toLocaleString()
    }
    return form.value.income_level2_upperbound
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.income_level2_upperbound = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.income_level2_upperbound = 0
    }
  }
})

const formattedIncomeLevel3Lowerbound = computed({
  get: () => {
    if (typeof form.value.income_level3_lowerbound === 'number') {
      return form.value.income_level3_lowerbound.toLocaleString()
    }
    return form.value.income_level3_lowerbound
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.income_level3_lowerbound = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.income_level3_lowerbound = 0
    }
  }
})

const formattedIncomeLevel3Upperbound = computed({
  get: () => {
    if (typeof form.value.income_level3_upperbound === 'number') {
      return form.value.income_level3_upperbound.toLocaleString()
    }
    return form.value.income_level3_upperbound
  },
  set: (val) => {
    if (val) {
      const parsed = parseFloat(val.replace(/,/g, ''))
      form.value.income_level3_upperbound = isNaN(parsed) ? 0 : parsed
    } else {
      form.value.income_level3_upperbound = 0
    }
  }
})

const statusOptions = [
  { title: 'Draft', value: 'draft' },
  { title: 'Active', value: 'active' },
  { title: 'Expired', value: 'expired' },
  { title: 'Cancelled', value: 'cancelled' },
  { title: 'Under Negotiation', value: 'under_negotiation' }
]
const typeOptions = [
  { title: 'Proportional', value: 'proportional' },
  { title: 'XL Risk', value: 'xl_risk' },
  { title: 'XL Event', value: 'xl_event' },
  { title: 'Stop Loss', value: 'stop_loss' },
  { title: 'Catastrophe XL', value: 'catastrophe_xl' }
]
const lobOptions = [
  { title: 'Group Life', value: 'group_life' },
  { title: 'Group Disability', value: 'group_disability' },
  { title: 'Funeral', value: 'funeral' },
  { title: 'Group Health', value: 'group_health' },
  { title: 'Credit Life', value: 'credit_life' }
]

// Helper function to parse formatted currency values
const parseFormattedNumber = (val) => {
  if (val === null || val === undefined || val === '') return null
  const parsed = parseFloat(String(val).replace(/,/g, ''))
  return isNaN(parsed) ? null : parsed
}

const rules = {
  required: [
    (v) =>
      (v !== null && v !== undefined && v !== '') || 'This field is required'
  ],
  positiveRequired: [
    (v) =>
      (v !== null && v !== undefined && v !== '') || 'This field is required',
    (v) => {
      const num = parseFormattedNumber(v)
      return num > 0 || 'Must be greater than 0'
    }
  ],
  nonNegative: [
    (v) => {
      if (v === '' || v === null) return true
      const num = parseFormattedNumber(v)
      return num >= 0 || 'Cannot be negative'
    }
  ],
  percentage: [
    (v) => {
      if (v === '' || v === null) return true
      const num = parseFormattedNumber(v)
      return (num >= 0 && num <= 100) || 'Must be between 0 and 100'
    }
  ],
  cessionPct: [
    (v) =>
      (v !== null && v !== undefined && v !== '') || 'Cession % is required',
    (v) => {
      const num = parseFormattedNumber(v)
      return num > 0 || 'Must be greater than 0'
    },
    (v) => {
      const num = parseFormattedNumber(v)
      return num <= 100 || 'Cannot exceed 100%'
    }
  ],
  surplusLines: [
    (v) =>
      (v !== null && v !== undefined && v !== '') ||
      'Surplus lines is required',
    (v) => {
      const num = parseFormattedNumber(v)
      return Number.isInteger(num) || 'Must be a whole number'
    },
    (v) => {
      const num = parseFormattedNumber(v)
      return num >= 1 || 'Must be at least 1'
    }
  ],
  xlLimit: [
    (v) =>
      (v !== null && v !== undefined && v !== '') || 'XL Limit is required',
    (v) => {
      const num = parseFormattedNumber(v)
      return num > 0 || 'Must be greater than 0'
    },
    (v) => {
      const num = parseFormattedNumber(v)
      return (
        num > form.value.xl_retention || 'XL Limit must exceed XL Retention'
      )
    }
  ],
  positiveDays: [
    (v) => {
      if (v === '' || v === null) return true
      const num = parseFormattedNumber(v)
      return num > 0 || 'Must be greater than 0'
    }
  ],
  expiryDate: [
    (v) => !!v || 'Expiry date is required',
    (v) =>
      !form.value.effective_date ||
      v > form.value.effective_date ||
      'Must be after effective date'
  ],
  renewalDate: [
    (v) =>
      !v ||
      !form.value.effective_date ||
      v >= form.value.effective_date ||
      'Must be on or after effective date'
  ],
  // Sum Assured Level Upper Bounds
  level1UpperBound: [
    (v) =>
      (v !== null && v !== undefined && v !== '') || 'This field is required',
    (v) => {
      const num = parseFormattedNumber(v)
      return (
        num >= form.value.level1_lowerbound ||
        'Upper bound must be >= lower bound'
      )
    }
  ],
  level2UpperBound: [
    (v) => {
      if (!v) return true
      const num = parseFormattedNumber(v)
      return (
        num >= form.value.level2_lowerbound ||
        'Upper bound must be >= lower bound'
      )
    }
  ],
  level3UpperBound: [
    (v) => {
      if (!v) return true
      const num = parseFormattedNumber(v)
      return (
        num >= form.value.level3_lowerbound ||
        'Upper bound must be >= lower bound'
      )
    }
  ],
  // Income Level Upper Bounds
  incomeLevel1UpperBound: [
    (v) =>
      (v !== null && v !== undefined && v !== '') || 'This field is required',
    (v) => {
      const num = parseFormattedNumber(v)
      return (
        num >= form.value.income_level1_lowerbound ||
        'Upper bound must be >= lower bound'
      )
    }
  ],
  incomeLevel2UpperBound: [
    (v) => {
      if (!v) return true
      const num = parseFormattedNumber(v)
      return (
        num >= form.value.income_level2_lowerbound ||
        'Upper bound must be >= lower bound'
      )
    }
  ],
  incomeLevel3UpperBound: [
    (v) => {
      if (!v) return true
      const num = parseFormattedNumber(v)
      return (
        num >= form.value.income_level3_lowerbound ||
        'Upper bound must be >= lower bound'
      )
    }
  ]
}

const statusColor = {
  draft: 'grey',
  active: 'success',
  expired: 'error',
  cancelled: 'error',
  under_negotiation: 'warning'
}
const typeLabel = {
  proportional: 'Proportional',
  quota_share: 'Proportional', // legacy value
  surplus: 'Proportional', // legacy value
  xl_risk: 'XL Risk',
  xl_event: 'XL Event',
  stop_loss: 'Stop Loss',
  catastrophe_xl: 'Cat XL'
}
const formatCurrency = (v) =>
  v != null
    ? 'R ' + Number(v).toLocaleString('en-ZA', { minimumFractionDigits: 0 })
    : '—'

const columnDefs = [
  { field: 'treaty_number', headerName: 'Treaty No.', width: 140 },
  { field: 'treaty_name', headerName: 'Name', flex: 1, minWidth: 160 },
  { field: 'reinsurer_name', headerName: 'Reinsurer', width: 160 },
  {
    field: 'treaty_type',
    headerName: 'Type',
    width: 120,
    cellRenderer: (p) =>
      `<span style="padding:2px 8px;border-radius:12px;background:#e3f2fd;font-size:12px">${typeLabel[p.value] || p.value}</span>`
  },
  { field: 'line_of_business', headerName: 'Line', width: 130 },
  { field: 'effective_date', headerName: 'Effective', width: 110 },
  { field: 'expiry_date', headerName: 'Expiry', width: 110 },
  {
    field: 'status',
    headerName: 'Status',
    width: 120,
    cellRenderer: (p) => {
      const c = statusColor[p.value] || 'grey'
      return `<v-chip style="padding:2px 8px;border-radius:12px;font-size:12px;background:var(--v-theme-${c}20);color:var(--v-theme-${c})">${p.value}</v-chip>`
    }
  },
  {
    field: 'large_claims_threshold',
    headerName: 'Large Claims Threshold',
    width: 180,
    valueFormatter: (p) => formatCurrency(p.value)
  },
  {
    field: 'is_run_off',
    headerName: 'Run-Off',
    width: 95,
    cellRenderer: (p) =>
      p.value
        ? '<span style="padding:2px 8px;border-radius:12px;font-size:11px;background:#ef6c0022;color:#ef6c00;font-weight:500">Run-Off</span>'
        : '—'
  },
  {
    headerName: 'Actions',
    width: 80,
    pinned: 'right',
    sortable: false,
    filter: false,
    cellRenderer: (p) => {
      const key = String(p.data.id)
      window[`showTreatyMenu_${key}`] = (event) =>
        showContextMenu(event, p.data)
      return `<div style="display:flex;align-items:center;justify-content:center;height:100%">
        <button onclick="showTreatyMenu_${key}(event)" title="Actions" style="background:none;border:none;cursor:pointer;padding:4px 8px;border-radius:4px;color:#616161;display:flex;align-items:center;justify-content:center">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><circle cx="12" cy="5" r="2"/><circle cx="12" cy="12" r="2"/><circle cx="12" cy="19" r="2"/></svg>
        </button>
      </div>`
    }
  }
]

const onGridReady = () => {}

let activeMenuCleanup = null

function showContextMenu(event, data) {
  if (activeMenuCleanup) activeMenuCleanup()

  const menuItems = [
    { label: 'Edit', color: '#1976d2', fn: () => openEditDialog(data) },
    { label: 'Schemes', color: '#388e3c', fn: () => openSchemesDialog(data) },
    { label: 'Delete', color: '#d32f2f', fn: () => openDeleteDialog(data) }
  ]

  const menu = document.createElement('div')
  menu.style.cssText =
    'position:fixed;background:#fff;border:1px solid #e0e0e0;border-radius:8px;' +
    'box-shadow:0 4px 16px rgba(0,0,0,0.14);z-index:9999;min-width:160px;padding:4px 0;'

  menuItems.forEach(({ label, color, fn }) => {
    const item = document.createElement('div')
    item.textContent = label
    item.style.cssText = `padding:8px 16px;cursor:pointer;font-size:13px;color:${color};`
    item.addEventListener(
      'mouseenter',
      () => (item.style.background = '#f5f5f5')
    )
    item.addEventListener('mouseleave', () => (item.style.background = ''))
    item.addEventListener('click', () => {
      cleanup()
      fn()
    })
    menu.appendChild(item)
  })

  document.body.appendChild(menu)

  const btn = event.currentTarget || event.target
  const rect = btn.getBoundingClientRect()
  menu.style.top = `${rect.bottom + 4}px`
  menu.style.left = `${rect.left}px`

  const mr = menu.getBoundingClientRect()
  if (mr.right > window.innerWidth - 8)
    menu.style.left = `${rect.right - mr.width}px`
  if (mr.bottom > window.innerHeight - 8)
    menu.style.top = `${rect.top - mr.height - 4}px`

  function cleanup() {
    menu.remove()
    document.removeEventListener('click', outsideClick, true)
    activeMenuCleanup = null
  }
  activeMenuCleanup = cleanup

  function outsideClick(e) {
    if (!menu.contains(e.target) && e.target !== btn) cleanup()
  }
  setTimeout(() => document.addEventListener('click', outsideClick, true), 0)
}

const notify = (message, color = 'success') => {
  snackbar.value = { show: true, message, color }
}

async function loadTreaties() {
  loading.value = true
  try {
    let res
    if (filters.value.run_off_only) {
      const params = {}
      if (filters.value.status) params.status = filters.value.status
      res = await GroupPricingService.getRunOffTreaties(params)
    } else {
      const params = {}
      if (filters.value.status) params.status = filters.value.status
      if (filters.value.type) params.type = filters.value.type
      if (filters.value.reinsurer) params.reinsurer = filters.value.reinsurer
      res = await GroupPricingService.getTreaties(params)
    }
    treaties.value = res.data?.data || []
  } catch {
    notify('Failed to load treaties', 'error')
  } finally {
    loading.value = false
  }
}

async function loadStats() {
  try {
    const res = await GroupPricingService.getTreatyStats()
    stats.value = res.data?.data || stats.value
  } catch {}
}

function openCreateDialog() {
  editingTreaty.value = null
  form.value = defaultForm()
  // Reset tier visibility flags
  showSumAssuredLevel2.value = false
  showSumAssuredLevel3.value = false
  showIncomeLevel2.value = false
  showIncomeLevel3.value = false
  showFormDialog.value = true
  nextTick(() => treatyForm.value?.resetValidation())
}

function openEditDialog(treaty) {
  editingTreaty.value = treaty
  form.value = { ...treaty }
  // Show tier levels if they have data
  showSumAssuredLevel2.value = !!(
    treaty.level2_upperbound || treaty.level2_ceded_proportion
  )
  showSumAssuredLevel3.value = !!(
    treaty.level3_upperbound || treaty.level3_ceded_proportion
  )
  showIncomeLevel2.value = !!(
    treaty.income_level2_upperbound || treaty.income_level2_ceded_proportion
  )
  showIncomeLevel3.value = !!(
    treaty.income_level3_upperbound || treaty.income_level3_ceded_proportion
  )
  showFormDialog.value = true
  nextTick(() => treatyForm.value?.resetValidation())
}

async function saveTreaty() {
  const { valid } = await treatyForm.value.validate()
  if (!valid) return
  saving.value = true
  try {
    if (editingTreaty.value) {
      await GroupPricingService.updateTreaty(editingTreaty.value.id, form.value)
      notify('Treaty updated successfully')
    } else {
      await GroupPricingService.createTreaty(form.value)
      notify('Treaty created successfully')
    }
    showFormDialog.value = false
    loadTreaties()
    loadStats()
  } catch (e) {
    notify(e.response?.data?.error || 'Failed to save treaty', 'error')
  } finally {
    saving.value = false
  }
}

function openDeleteDialog(treaty) {
  deletingTreaty.value = treaty
  showDeleteDialog.value = true
}

async function confirmDelete() {
  deleting.value = true
  try {
    await GroupPricingService.deleteTreaty(deletingTreaty.value.id)
    notify('Treaty deleted')
    showDeleteDialog.value = false
    loadTreaties()
    loadStats()
  } catch {
    notify('Failed to delete treaty', 'error')
  } finally {
    deleting.value = false
  }
}

async function loadAvailableSchemes() {
  if (availableSchemes.value.length > 0) return
  try {
    availableSchemes.value = await bordereauxStore.loadSchemes()
  } catch {}
}

async function openSchemesDialog(treaty) {
  selectedTreaty.value = treaty
  linkForm.value = { cession_override: 0, effective_date: '' }
  selectedSchemes.value = []
  selectedLinkRows.value = []
  showSchemesDialog.value = true
  await Promise.all([loadSchemeLinks(treaty.id), loadAvailableSchemes()])
}

async function loadSchemeLinks(treatyId) {
  try {
    const res = await GroupPricingService.getTreatySchemeLinks(treatyId)
    schemeLinks.value = res.data?.data || []
  } catch {
    notify('Failed to load scheme links', 'error')
  }
}

async function addSchemeLinks() {
  if (!selectedSchemes.value.length) {
    notify('Please select at least one scheme', 'warning')
    return
  }
  savingLink.value = true
  try {
    const res = await GroupPricingService.bulkLinkSchemesToTreaty(
      selectedTreaty.value.id,
      {
        scheme_ids: selectedSchemes.value,
        cession_override: linkForm.value.cession_override,
        effective_date: linkForm.value.effective_date
      }
    )
    const created = res.data?.data ?? selectedSchemes.value.length
    notify(`${created} scheme(s) linked successfully`)
    linkForm.value = { cession_override: 0, effective_date: '' }
    selectedSchemes.value = []
    await loadSchemeLinks(selectedTreaty.value.id)
  } catch (e) {
    notify(e.response?.data?.error || 'Failed to link schemes', 'error')
  } finally {
    savingLink.value = false
  }
}

function onLinkRowsSelected(rows) {
  selectedLinkRows.value = rows
}

function confirmBulkUnlink() {
  if (!selectedLinkRows.value.length) return
  showBulkUnlinkDialog.value = true
}

async function executeBulkUnlink() {
  deletingLinkLoading.value = true
  try {
    const linkIds = selectedLinkRows.value.map((r) => r.id)
    await GroupPricingService.bulkRemoveSchemeLinks(selectedTreaty.value.id, {
      link_ids: linkIds
    })
    notify(`${linkIds.length} scheme link(s) removed`)
    showBulkUnlinkDialog.value = false
    selectedLinkRows.value = []
    await loadSchemeLinks(selectedTreaty.value.id)
  } catch {
    notify('Failed to remove links', 'error')
  } finally {
    deletingLinkLoading.value = false
  }
}

async function loadBrokers() {
  try {
    const res = await GroupPricingService.getBrokers()
    brokers.value = res.data || []
  } catch {}
}

async function loadReinsurers() {
  try {
    const res = await GroupPricingService.getReinsurers()
    reinsurers.value = res.data || []
  } catch {}
}

function onReinsurerSelected(reinsurerName) {
  if (reinsurerName && reinsurers.value.length > 0) {
    const selectedReinsurer = reinsurers.value.find(
      (r) => r.name === reinsurerName
    )
    if (selectedReinsurer) {
      form.value.reinsurer_code = selectedReinsurer.code
    }
  } else {
    form.value.reinsurer_code = ''
  }
}

// Watchers for three-tier reinsurance structure - auto-cascade lower bounds
// Sum Assured Levels
watch(
  () => form.value.level1_upperbound,
  (newVal) => {
    if (
      newVal !== null &&
      newVal !== undefined &&
      newVal !== '' &&
      newVal > 0
    ) {
      form.value.level2_lowerbound = Number(newVal) + 1
    }
  }
)

watch(
  () => form.value.level2_upperbound,
  (newVal) => {
    if (
      newVal !== null &&
      newVal !== undefined &&
      newVal !== '' &&
      newVal > 0
    ) {
      form.value.level3_lowerbound = Number(newVal) + 1
    }
  }
)

// Income Levels
watch(
  () => form.value.income_level1_upperbound,
  (newVal) => {
    if (
      newVal !== null &&
      newVal !== undefined &&
      newVal !== '' &&
      newVal > 0
    ) {
      form.value.income_level2_lowerbound = Number(newVal) + 1
    }
  }
)

watch(
  () => form.value.income_level2_upperbound,
  (newVal) => {
    if (
      newVal !== null &&
      newVal !== undefined &&
      newVal !== '' &&
      newVal > 0
    ) {
      form.value.income_level3_lowerbound = Number(newVal) + 1
    }
  }
)

// Clear data when levels are hidden
watch(showSumAssuredLevel2, (newVal) => {
  if (!newVal) {
    form.value.level2_lowerbound = 0
    form.value.level2_upperbound = 0
    form.value.level2_ceded_proportion = 0
    showSumAssuredLevel3.value = false
  }
})

watch(showSumAssuredLevel3, (newVal) => {
  if (!newVal) {
    form.value.level3_lowerbound = 0
    form.value.level3_upperbound = 0
    form.value.level3_ceded_proportion = 0
  }
})

watch(showIncomeLevel2, (newVal) => {
  if (!newVal) {
    form.value.income_level2_lowerbound = 0
    form.value.income_level2_upperbound = 0
    form.value.income_level2_ceded_proportion = 0
    showIncomeLevel3.value = false
  }
})

watch(showIncomeLevel3, (newVal) => {
  if (!newVal) {
    form.value.income_level3_lowerbound = 0
    form.value.income_level3_upperbound = 0
    form.value.income_level3_ceded_proportion = 0
  }
})

onMounted(() => {
  loadTreaties()
  loadStats()
  loadBrokers()
  loadReinsurers()
})
</script>
