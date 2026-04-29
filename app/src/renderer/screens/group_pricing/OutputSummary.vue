<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <span class="headline"
              >Group Risk Rating Summary for {{ quote.scheme_name }}</span
            >
          </template>
          <template #default>
            <v-row>
              <v-col cols="3"><p>Scheme Name</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  quote.scheme_name
                }}</p></v-col
              >
              <v-col cols="3"><p>Start Date</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  parseDateString(quote.commencement_date)
                }}</p></v-col
              >
            </v-row>
            <v-row>
              <v-col cols="3"><p>Industry</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  quote.industry
                }}</p></v-col
              >
              <v-col cols="3"><p>Renewal/New</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  quote.quote_type
                }}</p></v-col
              >
            </v-row>
            <v-row>
              <v-col cols="3"><p>Distribution Channel</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  quote.distribution_channel || 'broker'
                }}</p></v-col
              >
              <v-col cols="3"><p>Broker</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  quote.quote_broker?.name || 'N/A'
                }}</p></v-col
              >
            </v-row>
            <v-row>
              <v-col cols="3"><p>Voluntary/Compulsory</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  quote.obligation_type
                }}</p></v-col
              >
            </v-row>
            <v-row>
              <v-col cols="3"><p>Currency</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  quote.currency
                }}</p></v-col
              >
              <v-col cols="3"><p>Exchange Rate</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  quote.exchangeRate
                }}</p></v-col
              >
            </v-row>
            <v-row>
              <v-col cols="3"><p>Expense Loading</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  quote.loadings.expense_loading + '%'
                }}</p></v-col
              >
              <v-col cols="3"><p>Commission Rate</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  quote.loadings.commission_loading + '%'
                }}</p></v-col
              >
            </v-row>
            <v-row>
              <v-col cols="3"><p>Contingency Loading</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  quote.loadings.contingency_loading + '%'
                }}</p></v-col
              >
              <v-col cols="3"><p>Admin Loading</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  quote.loadings.admin_loading + '%'
                }}</p></v-col
              >
            </v-row>
            <v-row>
              <v-col cols="3"><p>Other Loading</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  quote.loadings.other_loading + '%'
                }}</p></v-col
              >

              <v-col cols="3"><p>Profit Loading</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  quote.loadings.profit_loading + '%'
                }}</p></v-col
              >
            </v-row>
            <v-row>
              <v-col cols="3"><p>Member Count</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  resultSummary.member_count
                }}</p></v-col
              >

              <v-col cols="3"><p>Overall Discount</p></v-col>
              <v-col cols="3"
                ><p class="text-right content-bg">{{
                  quote.loadings.discount + '%'
                }}</p></v-col
              >
            </v-row>
            <v-divider class="my-5"></v-divider>
            <v-row>
              <v-col cols="3">
                <v-select
                  v-model="resultSummary"
                  variant="outlined"
                  density="compact"
                  :items="validResultSummaries"
                  item-title="category"
                  item-value="id"
                  return-object
                  label="Scheme Category"
                  placeholder="Select a Category"
                  @update:model-value="handleCategoryChange"
                ></v-select>
              </v-col>
            </v-row>
            <v-row v-if="selectedCategory" class="mb-2">
              <v-col class="d-flex justify-end gap-2">
                <v-btn
                  size="small"
                  color="success"
                  variant="tonal"
                  prepend-icon="mdi-microsoft-excel"
                  @click="exportAllToExcel"
                  >Export All to Excel</v-btn
                >
                <v-btn
                  size="small"
                  color="error"
                  variant="tonal"
                  prepend-icon="mdi-file-pdf-box"
                  @click="exportAllToPDF"
                  >Export All to PDF</v-btn
                >
              </v-col>
            </v-row>
            <v-expansion-panels v-if="selectedCategory">
              <v-expansion-panel elevation="1" tile>
                <v-expansion-panel-title
                  ><v-row align="center"
                    ><v-col cols="6">
                      <p
                        >{{ glaBenefitTitle }}
                        <v-chip
                          :color="
                            selectedCategory[0].gla_benefit ? 'green' : 'grey'
                          "
                          :text="
                            selectedCategory[0].gla_benefit
                              ? 'Active'
                              : 'Inactive'
                          "
                          size="small"
                          variant="tonal"
                        ></v-chip>
                      </p>
                    </v-col>
                    <v-col
                      cols="6"
                      class="d-flex justify-end align-center gap-1"
                    >
                      <v-btn
                        icon
                        size="x-small"
                        color="success"
                        variant="tonal"
                        title="Export to Excel"
                        :disabled="!selectedCategory[0].gla_benefit"
                        @click.stop="
                          exportBenefitToExcel('GLA', glaBenefitTitle)
                        "
                        ><v-icon size="small"
                          >mdi-microsoft-excel</v-icon
                        ></v-btn
                      >
                      <v-btn
                        icon
                        size="x-small"
                        color="error"
                        variant="tonal"
                        title="Export to PDF"
                        :disabled="!selectedCategory[0].gla_benefit"
                        @click.stop="exportBenefitToPDF('GLA', glaBenefitTitle)"
                        ><v-icon size="small">mdi-file-pdf-box</v-icon></v-btn
                      >
                    </v-col>
                  </v-row></v-expansion-panel-title
                >
                <v-expansion-panel-text>
                  <v-container>
                    <v-row>
                      <v-col cols="3"></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg"
                          >Technical Rate</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg"
                          >Experience Rated</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">Discounted</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n9">
                      <v-col cols="3"
                        ><p><b>GLA SUM ASSURED</b></p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Minimum Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.min_gla_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.min_gla_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.min_gla_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Maximum Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_gla_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_gla_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_gla_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Maximum FCL Capped Sum Assured</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_gla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_gla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_gla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Total Capped Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_gla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_gla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_gla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Average Covered Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_gla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_gla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_gla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>

                    <v-row class="mb-n9 mt-8">
                      <v-col cols="3"
                        ><p><b>RISK PREMIUM</b></p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Expected Number of Claims</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_gla_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_gla_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_gla_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Annual Risk Premium</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_gla_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_gla_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_gla_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Unit Rate per 1000 Sum Assured</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.gla_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_gla_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_gla_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Risk Premium as % of Annual Salary</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.proportion_gla_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_proportion_gla_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_proportion_gla_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n9 mt-8">
                      <v-col cols="3">
                        <p>
                          <b>OFFICE PREMIUM</b>
                          <v-tooltip location="end" max-width="500">
                            <template #activator="{ props: tProps }">
                              <v-icon
                                v-bind="tProps"
                                size="small"
                                class="ml-1"
                                color="primary"
                                style="cursor: help"
                                >mdi-information-outline</v-icon
                              >
                            </template>
                            <div style="white-space: pre-line">{{
                              commissionExplanation
                            }}</div>
                          </v-tooltip>
                        </p>
                      </v-col>
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Annual Office Premium</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              officePremiumWithCommission(
                                resultSummary.total_gla_annual_risk_premium,
                                resultSummary.total_gla_annual_commission_amount,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              officePremiumWithCommission(
                                resultSummary.exp_total_gla_annual_risk_premium,
                                resultSummary.exp_total_gla_annual_commission_amount,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              finalFieldValue(
                                resultSummary,
                                'final_gla_annual_office_premium'
                              )
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p
                          >Unit Office Premium Rate per 1000 Covered Sum
                          Assured</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeRateFromRiskRate(
                                resultSummary.gla_risk_rate_per_1000_sa,
                                resultSummary.total_gla_annual_commission_amount,
                                resultSummary.total_gla_capped_sum_assured,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeRateFromRiskRate(
                                resultSummary.exp_gla_risk_rate_per_1000_sa,
                                resultSummary.exp_total_gla_annual_commission_amount,
                                resultSummary.total_gla_capped_sum_assured,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              ratePer1000(
                                resultSummary.final_gla_annual_office_premium,
                                resultSummary.total_gla_capped_sum_assured
                              )
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p
                          >Office Premium Premium as % of Annual Salary</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeProportionFromRiskProportion(
                                resultSummary.proportion_gla_annual_risk_premium_salary,
                                resultSummary.total_gla_annual_commission_amount,
                                resultSummary.total_annual_salary,
                                resultSummary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeProportionFromRiskProportion(
                                resultSummary.exp_proportion_gla_annual_risk_premium_salary,
                                resultSummary.exp_total_gla_annual_commission_amount,
                                resultSummary.total_annual_salary,
                                resultSummary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              proportionOfSalary(
                                resultSummary.final_gla_annual_office_premium,
                                resultSummary.total_annual_salary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                    </v-row>
                  </v-container>
                </v-expansion-panel-text>
              </v-expansion-panel>
              <v-expansion-panel elevation="1" tile>
                <v-expansion-panel-title
                  ><v-row align="center"
                    ><v-col cols="6">
                      <p
                        >{{ ptdBenefitTitle }}
                        <v-chip
                          :color="
                            selectedCategory[0].ptd_benefit ? 'green' : 'grey'
                          "
                          :text="
                            selectedCategory[0].ptd_benefit
                              ? 'Active'
                              : 'Inactive'
                          "
                          size="small"
                          variant="tonal"
                        ></v-chip>
                      </p>
                    </v-col>
                    <v-col
                      cols="6"
                      class="d-flex justify-end align-center gap-1"
                    >
                      <v-btn
                        icon
                        size="x-small"
                        color="success"
                        variant="tonal"
                        title="Export to Excel"
                        :disabled="!selectedCategory[0].ptd_benefit"
                        @click.stop="
                          exportBenefitToExcel('PTD', ptdBenefitTitle)
                        "
                        ><v-icon size="small"
                          >mdi-microsoft-excel</v-icon
                        ></v-btn
                      >
                      <v-btn
                        icon
                        size="x-small"
                        color="error"
                        variant="tonal"
                        title="Export to PDF"
                        :disabled="!selectedCategory[0].ptd_benefit"
                        @click.stop="exportBenefitToPDF('PTD', ptdBenefitTitle)"
                        ><v-icon size="small">mdi-file-pdf-box</v-icon></v-btn
                      >
                    </v-col>
                  </v-row></v-expansion-panel-title
                >
                <v-expansion-panel-text>
                  <v-container>
                    <v-row>
                      <v-col cols="3"></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg"
                          >Technical Rate</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg"
                          >Experience Rated</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">Discounted</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n9">
                      <v-col cols="3"
                        ><p><b>SUM ASSURED</b></p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Minimum Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.min_ptd_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.min_ptd_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.min_ptd_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Maximum Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_ptd_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_ptd_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_ptd_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Maximum FCL Capped Sum Assured</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_ptd_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_ptd_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_ptd_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Total Capped Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_ptd_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_ptd_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_ptd_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Average Covered Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_ptd_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_ptd_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_ptd_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>

                    <v-row class="mb-n9 mt-8">
                      <v-col cols="3"
                        ><p><b>RISK PREMIUM</b></p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Expected Number of Claims</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_ptd_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_ptd_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_ptd_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Annual Risk Premium</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_ptd_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_ptd_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_ptd_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Unit Rate per 1000 Sum Assured</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.ptd_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_ptd_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_ptd_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Risk Premium as % of Annual Salary</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.proportion_ptd_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_proportion_ptd_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_proportion_ptd_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n9 mt-8">
                      <v-col cols="3">
                        <p>
                          <b>OFFICE PREMIUM</b>
                          <v-tooltip location="end" max-width="500">
                            <template #activator="{ props: tProps }">
                              <v-icon
                                v-bind="tProps"
                                size="small"
                                class="ml-1"
                                color="primary"
                                style="cursor: help"
                                >mdi-information-outline</v-icon
                              >
                            </template>
                            <div style="white-space: pre-line">{{
                              commissionExplanation
                            }}</div>
                          </v-tooltip>
                        </p>
                      </v-col>
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Annual Office Premium</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              officePremiumWithCommission(
                                resultSummary.total_ptd_annual_risk_premium,
                                resultSummary.total_ptd_annual_commission_amount,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              officePremiumWithCommission(
                                resultSummary.exp_total_ptd_annual_risk_premium,
                                resultSummary.exp_total_ptd_annual_commission_amount,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              finalFieldValue(
                                resultSummary,
                                'final_ptd_annual_office_premium'
                              )
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p
                          >Unit Office Premium Rate per 1000 Covered Sum
                          Assured</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeRateFromRiskRate(
                                resultSummary.ptd_risk_rate_per_1000_sa,
                                resultSummary.total_ptd_annual_commission_amount,
                                resultSummary.total_ptd_capped_sum_assured,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeRateFromRiskRate(
                                resultSummary.exp_ptd_risk_rate_per_1000_sa,
                                resultSummary.exp_total_ptd_annual_commission_amount,
                                resultSummary.total_ptd_capped_sum_assured,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              ratePer1000(
                                resultSummary.final_ptd_annual_office_premium,
                                resultSummary.total_ptd_capped_sum_assured
                              )
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p
                          >Office Premium Premium as % of Annual Salary</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeProportionFromRiskProportion(
                                resultSummary.proportion_ptd_annual_risk_premium_salary,
                                resultSummary.total_ptd_annual_commission_amount,
                                resultSummary.total_annual_salary,
                                resultSummary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeProportionFromRiskProportion(
                                resultSummary.exp_proportion_ptd_annual_risk_premium_salary,
                                resultSummary.exp_total_ptd_annual_commission_amount,
                                resultSummary.total_annual_salary,
                                resultSummary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              proportionOfSalary(
                                resultSummary.final_ptd_annual_office_premium,
                                resultSummary.total_annual_salary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                    </v-row>
                  </v-container>
                </v-expansion-panel-text>
              </v-expansion-panel>
              <v-expansion-panel elevation="1" tile>
                <v-expansion-panel-title
                  ><v-row align="center"
                    ><v-col cols="6">
                      <p
                        >{{ ciBenefitTitle }}
                        <v-chip
                          :color="
                            selectedCategory[0].ci_benefit ? 'green' : 'grey'
                          "
                          :text="
                            selectedCategory[0].ci_benefit
                              ? 'Active'
                              : 'Inactive'
                          "
                          size="small"
                          variant="tonal"
                        ></v-chip>
                      </p>
                    </v-col>
                    <v-col
                      cols="6"
                      class="d-flex justify-end align-center gap-1"
                    >
                      <v-btn
                        icon
                        size="x-small"
                        color="success"
                        variant="tonal"
                        title="Export to Excel"
                        :disabled="!selectedCategory[0].ci_benefit"
                        @click.stop="exportBenefitToExcel('CI', ciBenefitTitle)"
                        ><v-icon size="small"
                          >mdi-microsoft-excel</v-icon
                        ></v-btn
                      >
                      <v-btn
                        icon
                        size="x-small"
                        color="error"
                        variant="tonal"
                        title="Export to PDF"
                        :disabled="!selectedCategory[0].ci_benefit"
                        @click.stop="exportBenefitToPDF('CI', ciBenefitTitle)"
                        ><v-icon size="small">mdi-file-pdf-box</v-icon></v-btn
                      >
                    </v-col>
                  </v-row></v-expansion-panel-title
                >
                <v-expansion-panel-text>
                  <v-container>
                    <v-row>
                      <v-col cols="3"></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg"
                          >Technical Rate</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg"
                          >Experience Rated</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">Discounted</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n9">
                      <v-col cols="3"
                        ><p><b>SUM ASSURED</b></p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Minimum Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.min_ci_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.min_ci_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.min_ci_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Maximum Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_ci_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_ci_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_ci_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Maximum FCL Capped Sum Assured</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_ci_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_ci_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_ci_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Total Capped Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_ci_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_ci_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_ci_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Average Covered Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_ci_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_ci_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_ci_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>

                    <v-row class="mb-n9 mt-8">
                      <v-col cols="3"
                        ><p><b>RISK PREMIUM</b></p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Expected Number of Claims</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_ci_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_ci_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_ci_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Annual Risk Premium</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_ci_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_ci_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_ci_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Unit Rate per 1000 Sum Assured</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.ci_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_ci_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_ci_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Risk Premium as % of Annual Salary</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.proportion_ci_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_proportion_ci_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_proportion_ci_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n9 mt-8">
                      <v-col cols="3">
                        <p>
                          <b>OFFICE PREMIUM</b>
                          <v-tooltip location="end" max-width="500">
                            <template #activator="{ props: tProps }">
                              <v-icon
                                v-bind="tProps"
                                size="small"
                                class="ml-1"
                                color="primary"
                                style="cursor: help"
                                >mdi-information-outline</v-icon
                              >
                            </template>
                            <div style="white-space: pre-line">{{
                              commissionExplanation
                            }}</div>
                          </v-tooltip>
                        </p>
                      </v-col>
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Annual Office Premium</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              officePremiumWithCommission(
                                resultSummary.total_ci_annual_risk_premium,
                                resultSummary.total_ci_annual_commission_amount,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              officePremiumWithCommission(
                                resultSummary.exp_total_ci_annual_risk_premium,
                                resultSummary.exp_total_ci_annual_commission_amount,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              finalFieldValue(
                                resultSummary,
                                'final_ci_annual_office_premium'
                              )
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p
                          >Unit Office Premium Rate per 1000 Covered Sum
                          Assured</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeRateFromRiskRate(
                                resultSummary.ci_risk_rate_per_1000_sa,
                                resultSummary.total_ci_annual_commission_amount,
                                resultSummary.total_ci_capped_sum_assured,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeRateFromRiskRate(
                                resultSummary.exp_ci_risk_rate_per_1000_sa,
                                resultSummary.exp_total_ci_annual_commission_amount,
                                resultSummary.total_ci_capped_sum_assured,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              ratePer1000(
                                resultSummary.final_ci_annual_office_premium,
                                resultSummary.total_ci_capped_sum_assured
                              )
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p
                          >Office Premium Premium as % of Annual Salary</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeProportionFromRiskProportion(
                                resultSummary.proportion_ci_annual_risk_premium_salary,
                                resultSummary.total_ci_annual_commission_amount,
                                resultSummary.total_annual_salary,
                                resultSummary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeProportionFromRiskProportion(
                                resultSummary.exp_proportion_ci_annual_risk_premium_salary,
                                resultSummary.exp_total_ci_annual_commission_amount,
                                resultSummary.total_annual_salary,
                                resultSummary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              proportionOfSalary(
                                resultSummary.final_ci_annual_office_premium,
                                resultSummary.total_annual_salary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                    </v-row>
                  </v-container>
                </v-expansion-panel-text>
              </v-expansion-panel>
              <v-expansion-panel elevation="1" tile>
                <v-expansion-panel-title
                  ><v-row align="center"
                    ><v-col cols="6">
                      <p
                        >{{ phiBenefitTitle }}
                        <v-chip
                          :color="
                            selectedCategory[0].phi_benefit ? 'green' : 'grey'
                          "
                          :text="
                            selectedCategory[0].phi_benefit
                              ? 'Active'
                              : 'Inactive'
                          "
                          size="small"
                          variant="tonal"
                        ></v-chip
                      ></p>
                    </v-col>
                    <v-col
                      cols="6"
                      class="d-flex justify-end align-center gap-1"
                    >
                      <v-btn
                        icon
                        size="x-small"
                        color="success"
                        variant="tonal"
                        title="Export to Excel"
                        :disabled="!selectedCategory[0].phi_benefit"
                        @click.stop="
                          exportBenefitToExcel('PHI', phiBenefitTitle)
                        "
                        ><v-icon size="small"
                          >mdi-microsoft-excel</v-icon
                        ></v-btn
                      >
                      <v-btn
                        icon
                        size="x-small"
                        color="error"
                        variant="tonal"
                        title="Export to PDF"
                        :disabled="!selectedCategory[0].phi_benefit"
                        @click.stop="exportBenefitToPDF('PHI', phiBenefitTitle)"
                        ><v-icon size="small">mdi-file-pdf-box</v-icon></v-btn
                      >
                    </v-col>
                  </v-row></v-expansion-panel-title
                >
                <v-expansion-panel-text>
                  <v-container>
                    <v-row>
                      <v-col cols="3"></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg"
                          >Technical Rate</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg"
                          >Experience Rated</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">Discounted</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n9">
                      <v-col cols="3"
                        ><p><b>Regular Income</b></p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Minimum Income</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(resultSummary.min_phi_income)
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(resultSummary.min_phi_income)
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(resultSummary.min_phi_income)
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Maximum Income</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(resultSummary.max_phi_income)
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(resultSummary.max_phi_income)
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(resultSummary.max_phi_income)
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Maximum Capped Income</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_phi_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_phi_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_phi_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Total Capped Income</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_phi_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_phi_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_phi_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Average Covered Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_phi_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_phi_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_phi_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>

                    <v-row class="mb-n9 mt-8">
                      <v-col cols="3"
                        ><p><b>RISK PREMIUM</b></p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Expected Number of Claims</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_phi_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_phi_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_phi_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Annual Risk Premium</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_phi_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_phi_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_phi_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Unit Rate per 1000 Sum Assured</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.phi_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_phi_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_phi_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Risk Premium as % of Annual Salary</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.proportion_phi_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_proportion_phi_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_proportion_phi_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n9 mt-8">
                      <v-col cols="3">
                        <p>
                          <b>OFFICE PREMIUM</b>
                          <v-tooltip location="end" max-width="500">
                            <template #activator="{ props: tProps }">
                              <v-icon
                                v-bind="tProps"
                                size="small"
                                class="ml-1"
                                color="primary"
                                style="cursor: help"
                                >mdi-information-outline</v-icon
                              >
                            </template>
                            <div style="white-space: pre-line">{{
                              commissionExplanation
                            }}</div>
                          </v-tooltip>
                        </p>
                      </v-col>
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Annual Office Premium</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              officePremiumWithCommission(
                                resultSummary.total_phi_annual_risk_premium,
                                resultSummary.total_phi_annual_commission_amount,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              officePremiumWithCommission(
                                resultSummary.exp_total_phi_annual_risk_premium,
                                resultSummary.exp_total_phi_annual_commission_amount,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              finalFieldValue(
                                resultSummary,
                                'final_phi_annual_office_premium'
                              )
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p
                          >Unit Office Premium Rate per 1000 Covered Sum
                          Assured</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">0.00</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">0.00</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">0.00</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p
                          >Office Premium Premium as % of Annual Salary</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeProportionFromRiskProportion(
                                resultSummary.proportion_phi_annual_risk_premium_salary,
                                resultSummary.total_phi_annual_commission_amount,
                                resultSummary.total_annual_salary,
                                resultSummary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeProportionFromRiskProportion(
                                resultSummary.exp_proportion_phi_annual_risk_premium_salary,
                                resultSummary.exp_total_phi_annual_commission_amount,
                                resultSummary.total_annual_salary,
                                resultSummary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              proportionOfSalary(
                                resultSummary.final_phi_annual_office_premium,
                                resultSummary.total_annual_salary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                    </v-row>
                  </v-container>
                </v-expansion-panel-text>
              </v-expansion-panel>
              <v-expansion-panel elevation="1" tile>
                <v-expansion-panel-title
                  ><v-row align="center"
                    ><v-col cols="6">
                      <p
                        >{{ ttdBenefitTitle }}
                        <v-chip
                          :color="
                            selectedCategory[0].ttd_benefit ? 'green' : 'grey'
                          "
                          :text="
                            selectedCategory[0].ttd_benefit
                              ? 'Active'
                              : 'Inactive'
                          "
                          size="small"
                          variant="tonal"
                        ></v-chip
                      ></p>
                    </v-col>
                    <v-col
                      cols="6"
                      class="d-flex justify-end align-center gap-1"
                    >
                      <v-btn
                        icon
                        size="x-small"
                        color="success"
                        variant="tonal"
                        title="Export to Excel"
                        :disabled="!selectedCategory[0].ttd_benefit"
                        @click.stop="
                          exportBenefitToExcel('TTD', ttdBenefitTitle)
                        "
                        ><v-icon size="small"
                          >mdi-microsoft-excel</v-icon
                        ></v-btn
                      >
                      <v-btn
                        icon
                        size="x-small"
                        color="error"
                        variant="tonal"
                        title="Export to PDF"
                        :disabled="!selectedCategory[0].ttd_benefit"
                        @click.stop="exportBenefitToPDF('TTD', ttdBenefitTitle)"
                        ><v-icon size="small">mdi-file-pdf-box</v-icon></v-btn
                      >
                    </v-col>
                  </v-row></v-expansion-panel-title
                >
                <v-expansion-panel-text>
                  <v-container>
                    <v-row>
                      <v-col cols="3"></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg"
                          >Technical Rate</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg"
                          >Experience Rated</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">Discounted</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n9">
                      <v-col cols="3"
                        ><p><b>Regular Income</b></p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p> Minimum Income</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(resultSummary.min_ttd_income)
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(resultSummary.min_ttd_income)
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(resultSummary.min_ttd_income)
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Maximum Income</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(resultSummary.max_ttd_income)
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(resultSummary.max_ttd_income)
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(resultSummary.max_ttd_income)
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Maximum Capped Income</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_ttd_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_ttd_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_ttd_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Total Capped Income</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_ttd_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_ttd_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_ttd_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Average Capped Income</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_ttd_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_ttd_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_ttd_capped_income
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>

                    <v-row class="mb-n9 mt-8">
                      <v-col cols="3"
                        ><p><b>RISK PREMIUM</b></p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Expected Number of Claims</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_ttd_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_ttd_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_ttd_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Annual Risk Premium</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_ttd_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_ttd_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_ttd_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Unit Rate per 1000 Sum Assured</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.ttd_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_ttd_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_ttd_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Risk Premium as % of Annual Salary</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.proportion_ttd_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_proportion_ttd_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_proportion_ttd_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n9 mt-8">
                      <v-col cols="3">
                        <p>
                          <b>OFFICE PREMIUM</b>
                          <v-tooltip location="end" max-width="500">
                            <template #activator="{ props: tProps }">
                              <v-icon
                                v-bind="tProps"
                                size="small"
                                class="ml-1"
                                color="primary"
                                style="cursor: help"
                                >mdi-information-outline</v-icon
                              >
                            </template>
                            <div style="white-space: pre-line">{{
                              commissionExplanation
                            }}</div>
                          </v-tooltip>
                        </p>
                      </v-col>
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Annual Office Premium</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              officePremiumWithCommission(
                                resultSummary.total_ttd_annual_risk_premium,
                                resultSummary.total_ttd_annual_commission_amount,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              officePremiumWithCommission(
                                resultSummary.exp_total_ttd_annual_risk_premium,
                                resultSummary.exp_total_ttd_annual_commission_amount,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              finalFieldValue(
                                resultSummary,
                                'final_ttd_annual_office_premium'
                              )
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p
                          >Unit Office Premium Rate per 1000 Covered Sum
                          Assured</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">0.00</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">0.00</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">0.00</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p
                          >Office Premium Premium as % of Annual Salary</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeProportionFromRiskProportion(
                                resultSummary.proportion_ttd_annual_risk_premium_salary,
                                resultSummary.total_ttd_annual_commission_amount,
                                resultSummary.total_annual_salary,
                                resultSummary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeProportionFromRiskProportion(
                                resultSummary.exp_proportion_ttd_annual_risk_premium_salary,
                                resultSummary.exp_total_ttd_annual_commission_amount,
                                resultSummary.total_annual_salary,
                                resultSummary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              proportionOfSalary(
                                resultSummary.final_ttd_annual_office_premium,
                                resultSummary.total_annual_salary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                    </v-row>
                  </v-container>
                </v-expansion-panel-text>
              </v-expansion-panel>
              <v-expansion-panel elevation="1" tile>
                <v-expansion-panel-title
                  ><v-row align="center"
                    ><v-col cols="6">
                      <p
                        >{{ sglaBenefitTitle }}
                        <v-chip
                          :color="
                            selectedCategory[0].sgla_benefit ? 'green' : 'grey'
                          "
                          :text="
                            selectedCategory[0].sgla_benefit
                              ? 'Active'
                              : 'Inactive'
                          "
                          size="small"
                          variant="tonal"
                        ></v-chip
                      ></p>
                    </v-col>
                    <v-col
                      cols="6"
                      class="d-flex justify-end align-center gap-1"
                    >
                      <v-btn
                        icon
                        size="x-small"
                        color="success"
                        variant="tonal"
                        title="Export to Excel"
                        :disabled="!selectedCategory[0].sgla_benefit"
                        @click.stop="
                          exportBenefitToExcel('SGLA', sglaBenefitTitle)
                        "
                        ><v-icon size="small"
                          >mdi-microsoft-excel</v-icon
                        ></v-btn
                      >
                      <v-btn
                        icon
                        size="x-small"
                        color="error"
                        variant="tonal"
                        title="Export to PDF"
                        :disabled="!selectedCategory[0].sgla_benefit"
                        @click.stop="
                          exportBenefitToPDF('SGLA', sglaBenefitTitle)
                        "
                        ><v-icon size="small">mdi-file-pdf-box</v-icon></v-btn
                      >
                    </v-col>
                  </v-row></v-expansion-panel-title
                >
                <v-expansion-panel-text>
                  <v-container>
                    <v-row>
                      <v-col cols="3"></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg"
                          >Technical Rate</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg"
                          >Experience Rated</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">Discounted</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n9">
                      <v-col cols="3"
                        ><p><b>SUM ASSURED</b></p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Minimum Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.min_sgla_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.min_sgla_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.min_sgla_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Maximum Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_sgla_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_sgla_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_sgla_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Maximum FCL Capped Sum Assured</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_sgla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_sgla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.max_sgla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Total Capped Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_sgla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_sgla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_sgla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Average Covered Sum Assured</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_sgla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_sgla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.average_sgla_capped_sum_assured
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>

                    <v-row class="mb-n9 mt-8">
                      <v-col cols="3"
                        ><p><b>RISK PREMIUM</b></p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Expected Number of Claims</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_sgla_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_sgla_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_sgla_risk_rate
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Annual Risk Premium</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_sgla_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_sgla_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_sgla_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Unit Rate per 1000 Sum Assured</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.sgla_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_sgla_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_sgla_risk_rate_per_1000_sa
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Risk Premium as % of Annual Salary</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.proportion_sgla_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_proportion_sgla_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_proportion_sgla_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n9 mt-8">
                      <v-col cols="3">
                        <p>
                          <b>OFFICE PREMIUM</b>
                          <v-tooltip location="end" max-width="500">
                            <template #activator="{ props: tProps }">
                              <v-icon
                                v-bind="tProps"
                                size="small"
                                class="ml-1"
                                color="primary"
                                style="cursor: help"
                                >mdi-information-outline</v-icon
                              >
                            </template>
                            <div style="white-space: pre-line">{{
                              commissionExplanation
                            }}</div>
                          </v-tooltip>
                        </p>
                      </v-col>
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Annual Office Premium</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              officePremiumWithCommission(
                                resultSummary.total_sgla_annual_risk_premium,
                                resultSummary.total_sgla_annual_commission_amount,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              officePremiumWithCommission(
                                resultSummary.exp_total_sgla_annual_risk_premium,
                                resultSummary.exp_total_sgla_annual_commission_amount,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              finalFieldValue(
                                resultSummary,
                                'final_sgla_annual_office_premium'
                              )
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p
                          >Unit Office Premium Rate per 1000 Covered Sum
                          Assured</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeRateFromRiskRate(
                                resultSummary.sgla_risk_rate_per_1000_sa,
                                resultSummary.total_sgla_annual_commission_amount,
                                resultSummary.total_sgla_capped_sum_assured,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeRateFromRiskRate(
                                resultSummary.exp_sgla_risk_rate_per_1000_sa,
                                resultSummary.exp_total_sgla_annual_commission_amount,
                                resultSummary.total_sgla_capped_sum_assured,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              ratePer1000(
                                resultSummary.final_sgla_annual_office_premium,
                                resultSummary.total_sgla_capped_sum_assured
                              )
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p
                          >Office Premium Premium as % of Annual Salary</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeProportionFromRiskProportion(
                                resultSummary.proportion_sgla_annual_risk_premium_salary,
                                resultSummary.total_sgla_annual_commission_amount,
                                resultSummary.total_annual_salary,
                                resultSummary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeProportionFromRiskProportion(
                                resultSummary.exp_proportion_sgla_annual_risk_premium_salary,
                                resultSummary.exp_total_sgla_annual_commission_amount,
                                resultSummary.total_annual_salary,
                                resultSummary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              proportionOfSalary(
                                resultSummary.final_sgla_annual_office_premium,
                                resultSummary.total_annual_salary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                    </v-row>
                  </v-container>
                </v-expansion-panel-text>
              </v-expansion-panel>
              <v-expansion-panel elevation="1" tile>
                <v-expansion-panel-title
                  ><v-row align="center"
                    ><v-col cols="6">
                      <p
                        >{{ familyFuneralBenefitTitle }}
                        <v-chip
                          :color="
                            selectedCategory[0].family_funeral_benefit
                              ? 'green'
                              : 'grey'
                          "
                          :text="
                            selectedCategory[0].family_funeral_benefit
                              ? 'Active'
                              : 'Inactive'
                          "
                          size="small"
                          variant="tonal"
                        ></v-chip
                      ></p>
                    </v-col>
                    <v-col
                      cols="6"
                      class="d-flex justify-end align-center gap-1"
                    >
                      <v-btn
                        icon
                        size="x-small"
                        color="success"
                        variant="tonal"
                        title="Export to Excel"
                        :disabled="!selectedCategory[0].family_funeral_benefit"
                        @click.stop="
                          exportBenefitToExcel('FUN', familyFuneralBenefitTitle)
                        "
                        ><v-icon size="small"
                          >mdi-microsoft-excel</v-icon
                        ></v-btn
                      >
                      <v-btn
                        icon
                        size="x-small"
                        color="error"
                        variant="tonal"
                        title="Export to PDF"
                        :disabled="!selectedCategory[0].family_funeral_benefit"
                        @click.stop="
                          exportBenefitToPDF('FUN', familyFuneralBenefitTitle)
                        "
                        ><v-icon size="small">mdi-file-pdf-box</v-icon></v-btn
                      >
                    </v-col>
                  </v-row></v-expansion-panel-title
                >
                <v-expansion-panel-text>
                  <v-container>
                    <v-row>
                      <v-col cols="3"></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg"
                          >Technical Rate</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg"
                          >Experience Rated</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">Discounted</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n9 mt-8">
                      <v-col cols="3"
                        ><p><b>RISK PREMIUM</b></p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Annual Risk Premium</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_fun_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_fun_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_fun_annual_risk_premium
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Risk Premium per Member per Month</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.total_fun_annual_risk_premium /
                                (resultSummary.member_count * 12)
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_fun_annual_risk_premium /
                                (resultSummary.member_count * 12)
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_fun_annual_risk_premium /
                                (resultSummary.member_count * 12)
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Risk Premium as % of Annual Salary</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.proportion_fun_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_proportion_fun_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_proportion_fun_annual_risk_premium_salary *
                                100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n9 mt-8">
                      <v-col cols="3">
                        <p>
                          <b>OFFICE PREMIUM</b>
                          <v-tooltip location="end" max-width="500">
                            <template #activator="{ props: tProps }">
                              <v-icon
                                v-bind="tProps"
                                size="small"
                                class="ml-1"
                                color="primary"
                                style="cursor: help"
                                >mdi-information-outline</v-icon
                              >
                            </template>
                            <div style="white-space: pre-line">{{
                              commissionExplanation
                            }}</div>
                          </v-tooltip>
                        </p>
                      </v-col>
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"><p>Annual Office Premium</p></v-col>
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              officePremiumWithCommission(
                                resultSummary.total_fun_annual_risk_premium,
                                resultSummary.total_fun_annual_commission_amount,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              officePremiumWithCommission(
                                resultSummary.exp_total_fun_annual_risk_premium,
                                resultSummary.exp_total_fun_annual_commission_amount,
                                resultSummary
                              )
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              finalFieldValue(
                                resultSummary,
                                'final_fun_annual_office_premium'
                              )
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p>Office Premium per Member per Month</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_fun_monthly_premium_per_member
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_fun_monthly_premium_per_member
                            )
                          )
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              resultSummary.exp_total_fun_monthly_premium_per_member
                            )
                          )
                        }}</p></v-col
                      >
                    </v-row>
                    <v-row class="mb-n8">
                      <v-col cols="3"
                        ><p
                          >Office Premium Premium as % of Annual Salary</p
                        ></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeProportionFromRiskProportion(
                                resultSummary.proportion_fun_annual_risk_premium_salary,
                                resultSummary.total_fun_annual_commission_amount,
                                resultSummary.total_annual_salary,
                                resultSummary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              expOfficeProportionFromRiskProportion(
                                resultSummary.exp_proportion_fun_annual_risk_premium_salary,
                                resultSummary.exp_total_fun_annual_commission_amount,
                                resultSummary.total_annual_salary,
                                resultSummary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                      <v-col cols="3"
                        ><p class="text-center content-bg">{{
                          dashIfEmpty(
                            roundUpToTwoDecimals(
                              proportionOfSalary(
                                resultSummary.final_fun_annual_office_premium,
                                resultSummary.total_annual_salary
                              ) * 100
                            )
                          ) + '%'
                        }}</p></v-col
                      >
                    </v-row>
                  </v-container>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
          </template>
        </base-card>
      </v-col>
    </v-row>
  </v-container>
</template>
<script setup lang="ts">
import BaseCard from '@/renderer/components/BaseCard.vue'
import { onMounted, ref, watch, computed } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import * as XLSX from 'xlsx'
import jsPDF from 'jspdf'
import { applyPlugin } from 'jspdf-autotable'
import {
  computeOfficePremium,
  officePremiumWithCommission,
  expOfficeRateFromRiskRate,
  expOfficeProportionFromRiskProportion,
  officeRateFromRiskRate,
  officeProportionFromRiskProportion,
  ratePer1000,
  proportionOfSalary,
  finalFieldValue
} from '@/renderer/utils/quoteDataHelpers'
applyPlugin(jsPDF)

const glaBenefitTitle = ref('Group Life Assurance (GLA)')
const sglaBenefitTitle = ref('Spouse Group Life Assurance (GLA)')
const ptdBenefitTitle = ref('Permanent Total Disability')
const ciBenefitTitle = ref('Critical Illness')
const phiBenefitTitle = ref('Personal Health Insurance')
const ttdBenefitTitle = ref('Temporary Total Disability')
const familyFuneralBenefitTitle = ref('Family Funeral')
const benefitMaps: any = ref([])
const resultSummary: any = ref(null)
const selectedCategory: any = ref(null)

// Tooltip text for the OFFICE PREMIUM section info icon. Explains progressive
// commission and how experience rating propagates through the scheme total —
// notably why a benefit with experience adjustment == 1 can still see a
// Technical vs Experience commission gap (other benefits' experience shifts
// the scheme total, which moves the progressive band, which changes EVERY
// benefit's commission slice).
const commissionExplanation = `Progressive Commission & Experience Rating

Commission is allocated using a progressive rate table on the SCHEME total premium (not per benefit):

• Technical: book risk / (1 − loadings) + commission slice from the BOOK scheme total band.

• Experience-Rated: experience-blended risk / (1 − loadings) + commission slice from the EXPERIENCE scheme total band.

• Discounted: same as Experience-Rated but with the discount applied on top of the loadings.

The commission rate depends on the SCHEME total. Even when one benefit's experience adjustment is 1 (book risk == experience risk), other benefits' experience can shift the scheme total enough to land in a different progressive band — which changes every benefit's commission slice.

Distribution: scheme commission → categories (split by each category's share of premium) → benefits within a category (split by each benefit's share of premium). The last category and last benefit absorb any rounding residual so totals match exactly.`

const props = defineProps({
  quote: {
    type: Object,
    default: () => ({})
  },
  resultSummaries: {
    type: Array,
    default: () => []
  }
})

// Watch for resultSummaries changes and update resultSummary
watch(
  () => props.resultSummaries,
  (newVal) => {
    if (Array.isArray(newVal) && newVal.length > 0) {
      resultSummary.value = newVal[0]
    } else {
      resultSummary.value = null
    }
    // Reset selectedCategory when resultSummaries change
    selectedCategory.value = null
  },
  { immediate: true, deep: true }
)

// Watch for quote changes
watch(
  () => props.quote,
  (newVal) => {
    // console.log('Quote updated:', newVal)
  },
  { immediate: true, deep: true }
)

// Computed property to ensure we have valid resultSummaries
const validResultSummaries = computed(() => {
  return Array.isArray(props.resultSummaries) ? props.resultSummaries : []
})

const parseDateString = (dateString) => {
  const date = new Date(dateString)
  const formattedDate = date.toISOString().split('T')[0]
  return formattedDate
}

const roundUpToTwoDecimals = (num) => {
  const roundedNum = Math.ceil(num * 100) / 100 // Round up to two decimal places
  return roundedNum
    .toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    })
    .replace(/,/g, ' ') // Replace commas with spaces for accounting format }
}

const dashIfEmpty = (value: any) => {
  return value || '-'
}

const handleCategoryChange = (event: any) => {
  if (resultSummary.value && props.quote.scheme_categories) {
    selectedCategory.value = props.quote.scheme_categories.filter(
      (category) => category.scheme_category === resultSummary.value.category
    )
  }
}
// === EXPORT UTILITIES ===

const fmtNum = (val: any, isPercent = false) => {
  const n = isPercent ? val * 100 : val
  if (n == null || isNaN(n)) return '-'
  return (
    (Math.ceil(n * 100) / 100)
      .toLocaleString('en-US', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      })
      .replace(/,/g, ' ') + (isPercent ? '%' : '')
  )
}

const SECTION_LABELS = [
  'GLA SUM ASSURED',
  'SUM ASSURED',
  'REGULAR INCOME',
  'RISK PREMIUM',
  'OFFICE PREMIUM'
]

const getBenefitRows = (rs: any, code: string): string[][] => {
  const rows: string[][] = []
  // For office-premium / office-rate / office-proportion fields the canonical
  // value is no longer persisted on the summary — derive on the fly from the
  // matching risk-side field. The helpers detect the suffix and resolve it
  // automatically so the field-name strings in the row tables stay readable.
  const resolve = (field: string): number => {
    if (field.includes('_office_premium')) {
      const risk = field
        .replace('_office_premium', '_risk_premium')
        .replace('annual_risk_premium_salary', 'annual_risk_premium_salary')
      // Total fields keep "annual"; proportions need different mapping.
      return computeOfficePremium(
        rs[risk] ?? rs[risk.replace('_risk_premium', '_annual_risk_premium')],
        rs
      )
    }
    if (field.includes('office_rate_per_1000_sa')) {
      const risk = field.replace(
        'office_rate_per_1000_sa',
        'risk_rate_per_1000_sa'
      )
      return officeRateFromRiskRate(rs[risk] ?? 0, rs)
    }
    if (field.includes('office_premium_salary')) {
      // proportion_<b>_office_premium_salary -> proportion_<b>_annual_risk_premium_salary
      // (or for educator/conv slices: proportion_<b>_risk_premium_salary)
      let risk = field.replace(
        'office_premium_salary',
        'annual_risk_premium_salary'
      )
      if (rs[risk] === undefined) {
        risk = field.replace('office_premium_salary', 'risk_premium_salary')
      }
      return officeProportionFromRiskProportion(rs[risk] ?? 0, rs)
    }
    return rs[field] ?? 0
  }
  const th = (field: string) => fmtNum(resolve(field))
  const ex = (field: string) => fmtNum(resolve(field))
  const thP = (field: string) => fmtNum(resolve(field), true)
  const exP = (field: string) => fmtNum(resolve(field), true)
  const f = (field: string) => fmtNum(finalFieldValue(rs, field))
  // Office premium WITH commission slice. Used for both the Technical
  // column (book risk + book commission) and the Experience-Rated column
  // (exp risk + exp commission). The helper itself is tier-neutral; the
  // caller picks which risk and commission fields to feed it.
  const withC = (riskField: string, commissionField: string) =>
    fmtNum(
      officePremiumWithCommission(
        rs[riskField] ?? 0,
        rs[commissionField] ?? 0,
        rs
      )
    )
  const exCRate = (
    riskRateField: string,
    commissionField: string,
    cappedField: string
  ) =>
    fmtNum(
      expOfficeRateFromRiskRate(
        rs[riskRateField] ?? 0,
        rs[commissionField] ?? 0,
        rs[cappedField] ?? 0,
        rs
      )
    )
  const exCProp = (
    riskPropField: string,
    commissionField: string,
    salaryField: string
  ) =>
    fmtNum(
      expOfficeProportionFromRiskProportion(
        rs[riskPropField] ?? 0,
        rs[commissionField] ?? 0,
        rs[salaryField] ?? 0,
        rs
      ),
      true
    )
  // Discounted-column rate / proportion helpers. Both derive directly from
  // the persisted final_*_annual_office_premium so the export reconciles to
  // the same Final amount the on-screen "Annual Office Premium" cell shows
  // (post-discount AND post-commission).
  const fRate = (finalField: string, cappedField: string) =>
    fmtNum(ratePer1000(rs[finalField] ?? 0, rs[cappedField] ?? 0))
  const fProp = (finalField: string, salaryField: string) =>
    fmtNum(proportionOfSalary(rs[finalField] ?? 0, rs[salaryField] ?? 0), true)
  const row = (
    label: string,
    thVal: string,
    expVal: string,
    finalVal?: string
  ) => rows.push([label, thVal, expVal, finalVal ?? expVal])
  const sec = (label: string) => rows.push([label, '', '', ''])

  switch (code) {
    case 'GLA':
      sec('GLA SUM ASSURED')
      row(
        'Minimum Sum Assured',
        th('min_gla_sum_assured'),
        ex('min_gla_sum_assured')
      )
      row(
        'Maximum Sum Assured',
        th('max_gla_sum_assured'),
        ex('max_gla_sum_assured')
      )
      row(
        'Maximum FCL Capped Sum Assured',
        th('max_gla_capped_sum_assured'),
        ex('max_gla_capped_sum_assured')
      )
      row(
        'Total Capped Sum Assured',
        th('total_gla_capped_sum_assured'),
        ex('total_gla_capped_sum_assured')
      )
      row(
        'Average Covered Sum Assured',
        th('average_gla_capped_sum_assured'),
        ex('average_gla_capped_sum_assured')
      )
      sec('RISK PREMIUM')
      row(
        'Expected Number of Claims',
        th('total_gla_risk_rate'),
        ex('exp_total_gla_risk_rate')
      )
      row(
        'Annual Risk Premium',
        th('total_gla_annual_risk_premium'),
        ex('exp_total_gla_annual_risk_premium')
      )
      row(
        'Unit Rate per 1000 Sum Assured',
        th('gla_risk_rate_per_1000_sa'),
        ex('exp_gla_risk_rate_per_1000_sa')
      )
      row(
        'Risk Premium as % of Annual Salary',
        thP('proportion_gla_annual_risk_premium_salary'),
        exP('exp_proportion_gla_annual_risk_premium_salary')
      )
      sec('OFFICE PREMIUM')
      row(
        'Annual Office Premium',
        withC(
          'total_gla_annual_risk_premium',
          'total_gla_annual_commission_amount'
        ),
        withC(
          'exp_total_gla_annual_risk_premium',
          'exp_total_gla_annual_commission_amount'
        ),
        f('final_gla_annual_office_premium')
      )
      row(
        'Unit Office Premium Rate per 1000 Covered Sum Assured',
        exCRate(
          'gla_risk_rate_per_1000_sa',
          'total_gla_annual_commission_amount',
          'total_gla_capped_sum_assured'
        ),
        exCRate(
          'exp_gla_risk_rate_per_1000_sa',
          'exp_total_gla_annual_commission_amount',
          'total_gla_capped_sum_assured'
        ),
        fRate('final_gla_annual_office_premium', 'total_gla_capped_sum_assured')
      )
      row(
        'Office Premium as % of Annual Salary',
        exCProp(
          'proportion_gla_annual_risk_premium_salary',
          'total_gla_annual_commission_amount',
          'total_annual_salary'
        ),
        exCProp(
          'exp_proportion_gla_annual_risk_premium_salary',
          'exp_total_gla_annual_commission_amount',
          'total_annual_salary'
        ),
        fProp('final_gla_annual_office_premium', 'total_annual_salary')
      )
      break
    case 'PTD':
      sec('SUM ASSURED')
      row(
        'Minimum Sum Assured',
        th('min_ptd_sum_assured'),
        ex('min_ptd_sum_assured')
      )
      row(
        'Maximum Sum Assured',
        th('max_ptd_sum_assured'),
        ex('max_ptd_sum_assured')
      )
      row(
        'Maximum FCL Capped Sum Assured',
        th('max_ptd_capped_sum_assured'),
        ex('max_ptd_capped_sum_assured')
      )
      row(
        'Total Capped Sum Assured',
        th('total_ptd_capped_sum_assured'),
        ex('total_ptd_capped_sum_assured')
      )
      row(
        'Average Covered Sum Assured',
        th('average_ptd_capped_sum_assured'),
        ex('average_ptd_capped_sum_assured')
      )
      sec('RISK PREMIUM')
      row(
        'Expected Number of Claims',
        th('total_ptd_risk_rate'),
        ex('exp_total_ptd_risk_rate')
      )
      row(
        'Annual Risk Premium',
        th('total_ptd_annual_risk_premium'),
        ex('exp_total_ptd_annual_risk_premium')
      )
      row(
        'Unit Rate per 1000 Sum Assured',
        th('ptd_risk_rate_per_1000_sa'),
        ex('exp_ptd_risk_rate_per_1000_sa')
      )
      row(
        'Risk Premium as % of Annual Salary',
        thP('proportion_ptd_annual_risk_premium_salary'),
        exP('exp_proportion_ptd_annual_risk_premium_salary')
      )
      sec('OFFICE PREMIUM')
      row(
        'Annual Office Premium',
        withC(
          'total_ptd_annual_risk_premium',
          'total_ptd_annual_commission_amount'
        ),
        withC(
          'exp_total_ptd_annual_risk_premium',
          'exp_total_ptd_annual_commission_amount'
        ),
        f('final_ptd_annual_office_premium')
      )
      row(
        'Unit Office Premium Rate per 1000 Covered Sum Assured',
        exCRate(
          'ptd_risk_rate_per_1000_sa',
          'total_ptd_annual_commission_amount',
          'total_ptd_capped_sum_assured'
        ),
        exCRate(
          'exp_ptd_risk_rate_per_1000_sa',
          'exp_total_ptd_annual_commission_amount',
          'total_ptd_capped_sum_assured'
        ),
        fRate('final_ptd_annual_office_premium', 'total_ptd_capped_sum_assured')
      )
      row(
        'Office Premium as % of Annual Salary',
        exCProp(
          'proportion_ptd_annual_risk_premium_salary',
          'total_ptd_annual_commission_amount',
          'total_annual_salary'
        ),
        exCProp(
          'exp_proportion_ptd_annual_risk_premium_salary',
          'exp_total_ptd_annual_commission_amount',
          'total_annual_salary'
        ),
        fProp('final_ptd_annual_office_premium', 'total_annual_salary')
      )
      break
    case 'CI':
      sec('SUM ASSURED')
      row(
        'Minimum Sum Assured',
        th('min_ci_sum_assured'),
        ex('min_ci_sum_assured')
      )
      row(
        'Maximum Sum Assured',
        th('max_ci_sum_assured'),
        ex('max_ci_sum_assured')
      )
      row(
        'Maximum FCL Capped Sum Assured',
        th('max_ci_capped_sum_assured'),
        ex('max_ci_capped_sum_assured')
      )
      row(
        'Total Capped Sum Assured',
        th('total_ci_capped_sum_assured'),
        ex('total_ci_capped_sum_assured')
      )
      row(
        'Average Covered Sum Assured',
        th('average_ci_capped_sum_assured'),
        ex('average_ci_capped_sum_assured')
      )
      sec('RISK PREMIUM')
      row(
        'Expected Number of Claims',
        th('total_ci_risk_rate'),
        ex('exp_total_ci_risk_rate')
      )
      row(
        'Annual Risk Premium',
        th('total_ci_annual_risk_premium'),
        ex('exp_total_ci_annual_risk_premium')
      )
      row(
        'Unit Rate per 1000 Sum Assured',
        th('ci_risk_rate_per_1000_sa'),
        ex('exp_ci_risk_rate_per_1000_sa')
      )
      row(
        'Risk Premium as % of Annual Salary',
        thP('proportion_ci_annual_risk_premium_salary'),
        exP('exp_proportion_ci_annual_risk_premium_salary')
      )
      sec('OFFICE PREMIUM')
      row(
        'Annual Office Premium',
        withC(
          'total_ci_annual_risk_premium',
          'total_ci_annual_commission_amount'
        ),
        withC(
          'exp_total_ci_annual_risk_premium',
          'exp_total_ci_annual_commission_amount'
        ),
        f('final_ci_annual_office_premium')
      )
      row(
        'Unit Office Premium Rate per 1000 Covered Sum Assured',
        exCRate(
          'ci_risk_rate_per_1000_sa',
          'total_ci_annual_commission_amount',
          'total_ci_capped_sum_assured'
        ),
        exCRate(
          'exp_ci_risk_rate_per_1000_sa',
          'exp_total_ci_annual_commission_amount',
          'total_ci_capped_sum_assured'
        ),
        fRate('final_ci_annual_office_premium', 'total_ci_capped_sum_assured')
      )
      row(
        'Office Premium as % of Annual Salary',
        exCProp(
          'proportion_ci_annual_risk_premium_salary',
          'total_ci_annual_commission_amount',
          'total_annual_salary'
        ),
        exCProp(
          'exp_proportion_ci_annual_risk_premium_salary',
          'exp_total_ci_annual_commission_amount',
          'total_annual_salary'
        ),
        fProp('final_ci_annual_office_premium', 'total_annual_salary')
      )
      break
    case 'SGLA':
      sec('SUM ASSURED')
      row(
        'Minimum Sum Assured',
        th('min_sgla_sum_assured'),
        ex('min_sgla_sum_assured')
      )
      row(
        'Maximum Sum Assured',
        th('max_sgla_sum_assured'),
        ex('max_sgla_sum_assured')
      )
      row(
        'Maximum FCL Capped Sum Assured',
        th('max_sgla_capped_sum_assured'),
        ex('max_sgla_capped_sum_assured')
      )
      row(
        'Total Capped Sum Assured',
        th('total_sgla_capped_sum_assured'),
        ex('total_sgla_capped_sum_assured')
      )
      row(
        'Average Covered Sum Assured',
        th('average_sgla_capped_sum_assured'),
        ex('average_sgla_capped_sum_assured')
      )
      sec('RISK PREMIUM')
      row(
        'Expected Number of Claims',
        th('total_sgla_risk_rate'),
        ex('exp_total_sgla_risk_rate')
      )
      row(
        'Annual Risk Premium',
        th('total_sgla_annual_risk_premium'),
        ex('exp_total_sgla_annual_risk_premium')
      )
      row(
        'Unit Rate per 1000 Sum Assured',
        th('sgla_risk_rate_per_1000_sa'),
        ex('exp_sgla_risk_rate_per_1000_sa')
      )
      row(
        'Risk Premium as % of Annual Salary',
        thP('proportion_sgla_annual_risk_premium_salary'),
        exP('exp_proportion_sgla_annual_risk_premium_salary')
      )
      sec('OFFICE PREMIUM')
      row(
        'Annual Office Premium',
        withC(
          'total_sgla_annual_risk_premium',
          'total_sgla_annual_commission_amount'
        ),
        withC(
          'exp_total_sgla_annual_risk_premium',
          'exp_total_sgla_annual_commission_amount'
        ),
        f('final_sgla_annual_office_premium')
      )
      row(
        'Unit Office Premium Rate per 1000 Covered Sum Assured',
        exCRate(
          'sgla_risk_rate_per_1000_sa',
          'total_sgla_annual_commission_amount',
          'total_sgla_capped_sum_assured'
        ),
        exCRate(
          'exp_sgla_risk_rate_per_1000_sa',
          'exp_total_sgla_annual_commission_amount',
          'total_sgla_capped_sum_assured'
        ),
        fRate(
          'final_sgla_annual_office_premium',
          'total_sgla_capped_sum_assured'
        )
      )
      row(
        'Office Premium as % of Annual Salary',
        exCProp(
          'proportion_sgla_annual_risk_premium_salary',
          'total_sgla_annual_commission_amount',
          'total_annual_salary'
        ),
        exCProp(
          'exp_proportion_sgla_annual_risk_premium_salary',
          'exp_total_sgla_annual_commission_amount',
          'total_annual_salary'
        ),
        fProp('final_sgla_annual_office_premium', 'total_annual_salary')
      )
      break
    case 'PHI':
      sec('REGULAR INCOME')
      row('Minimum Income', th('min_phi_income'), ex('min_phi_income'))
      row('Maximum Income', th('max_phi_income'), ex('max_phi_income'))
      row(
        'Maximum Capped Income',
        th('max_phi_capped_income'),
        ex('max_phi_capped_income')
      )
      row(
        'Total Capped Income',
        th('total_phi_capped_income'),
        ex('total_phi_capped_income')
      )
      row(
        'Average Covered Sum Assured',
        th('average_phi_capped_income'),
        ex('average_phi_capped_income')
      )
      sec('RISK PREMIUM')
      row(
        'Expected Number of Claims',
        th('total_phi_risk_rate'),
        ex('exp_total_phi_risk_rate')
      )
      row(
        'Annual Risk Premium',
        th('total_phi_annual_risk_premium'),
        ex('exp_total_phi_annual_risk_premium')
      )
      row(
        'Unit Rate per 1000 Sum Assured',
        th('phi_risk_rate_per_1000_sa'),
        ex('exp_phi_risk_rate_per_1000_sa')
      )
      row(
        'Risk Premium as % of Annual Salary',
        thP('proportion_phi_annual_risk_premium_salary'),
        exP('exp_proportion_phi_annual_risk_premium_salary')
      )
      sec('OFFICE PREMIUM')
      row(
        'Annual Office Premium',
        withC(
          'total_phi_annual_risk_premium',
          'total_phi_annual_commission_amount'
        ),
        withC(
          'exp_total_phi_annual_risk_premium',
          'exp_total_phi_annual_commission_amount'
        ),
        f('final_phi_annual_office_premium')
      )
      row(
        'Unit Office Premium Rate per 1000 Covered Sum Assured',
        '0.00',
        '0.00',
        '0.00'
      )
      row(
        'Office Premium as % of Annual Salary',
        exCProp(
          'proportion_phi_annual_risk_premium_salary',
          'total_phi_annual_commission_amount',
          'total_annual_salary'
        ),
        exCProp(
          'exp_proportion_phi_annual_risk_premium_salary',
          'exp_total_phi_annual_commission_amount',
          'total_annual_salary'
        ),
        fProp('final_phi_annual_office_premium', 'total_annual_salary')
      )
      break
    case 'TTD':
      sec('REGULAR INCOME')
      row('Minimum Income', th('min_ttd_income'), ex('min_ttd_income'))
      row('Maximum Income', th('max_ttd_income'), ex('max_ttd_income'))
      row(
        'Maximum Capped Income',
        th('max_ttd_capped_income'),
        ex('max_ttd_capped_income')
      )
      row(
        'Total Capped Income',
        th('total_ttd_capped_income'),
        ex('total_ttd_capped_income')
      )
      row(
        'Average Covered Income',
        th('average_ttd_capped_income'),
        ex('average_ttd_capped_income')
      )
      sec('RISK PREMIUM')
      row(
        'Expected Number of Claims',
        th('total_ttd_risk_rate'),
        ex('exp_total_ttd_risk_rate')
      )
      row(
        'Annual Risk Premium',
        th('total_ttd_annual_risk_premium'),
        ex('exp_total_ttd_annual_risk_premium')
      )
      row(
        'Unit Rate per 1000 Sum Assured',
        th('ttd_risk_rate_per_1000_sa'),
        ex('exp_ttd_risk_rate_per_1000_sa')
      )
      row(
        'Risk Premium as % of Annual Salary',
        thP('proportion_ttd_annual_risk_premium_salary'),
        exP('exp_proportion_ttd_annual_risk_premium_salary')
      )
      sec('OFFICE PREMIUM')
      row(
        'Annual Office Premium',
        withC(
          'total_ttd_annual_risk_premium',
          'total_ttd_annual_commission_amount'
        ),
        withC(
          'exp_total_ttd_annual_risk_premium',
          'exp_total_ttd_annual_commission_amount'
        ),
        f('final_ttd_annual_office_premium')
      )
      row(
        'Unit Office Premium Rate per 1000 Covered Sum Assured',
        '0.00',
        '0.00',
        '0.00'
      )
      row(
        'Office Premium as % of Annual Salary',
        exCProp(
          'proportion_ttd_annual_risk_premium_salary',
          'total_ttd_annual_commission_amount',
          'total_annual_salary'
        ),
        exCProp(
          'exp_proportion_ttd_annual_risk_premium_salary',
          'exp_total_ttd_annual_commission_amount',
          'total_annual_salary'
        ),
        fProp('final_ttd_annual_office_premium', 'total_annual_salary')
      )
      break
    case 'FUN': {
      const mc = rs.member_count || 1
      sec('RISK PREMIUM')
      row(
        'Annual Risk Premium',
        th('total_fun_annual_risk_premium'),
        ex('exp_total_fun_annual_risk_premium')
      )
      row(
        'Risk Premium per Member per Month',
        fmtNum(rs.total_fun_annual_risk_premium / (mc * 12)),
        fmtNum(rs.exp_total_fun_annual_risk_premium / (mc * 12))
      )
      row(
        'Risk Premium as % of Annual Salary',
        thP('proportion_fun_annual_risk_premium_salary'),
        exP('exp_proportion_fun_annual_risk_premium_salary')
      )
      sec('OFFICE PREMIUM')
      row(
        'Annual Office Premium',
        withC(
          'total_fun_annual_risk_premium',
          'total_fun_annual_commission_amount'
        ),
        withC(
          'exp_total_fun_annual_risk_premium',
          'exp_total_fun_annual_commission_amount'
        ),
        f('final_fun_annual_office_premium')
      )
      row(
        'Office Premium per Member per Month',
        ex('exp_total_fun_monthly_premium_per_member'),
        ex('exp_total_fun_monthly_premium_per_member')
      )
      row(
        'Office Premium as % of Annual Salary',
        exCProp(
          'proportion_fun_annual_risk_premium_salary',
          'total_fun_annual_commission_amount',
          'total_annual_salary'
        ),
        exCProp(
          'exp_proportion_fun_annual_risk_premium_salary',
          'exp_total_fun_annual_commission_amount',
          'total_annual_salary'
        ),
        fProp('final_fun_annual_office_premium', 'total_annual_salary')
      )
      break
    }
  }
  return rows
}

const BENEFIT_LIST = [
  { code: 'GLA', titleRef: glaBenefitTitle, flagKey: 'gla_benefit' },
  { code: 'SGLA', titleRef: sglaBenefitTitle, flagKey: 'sgla_benefit' },
  { code: 'PTD', titleRef: ptdBenefitTitle, flagKey: 'ptd_benefit' },
  { code: 'CI', titleRef: ciBenefitTitle, flagKey: 'ci_benefit' },
  { code: 'PHI', titleRef: phiBenefitTitle, flagKey: 'phi_benefit' },
  { code: 'TTD', titleRef: ttdBenefitTitle, flagKey: 'ttd_benefit' },
  {
    code: 'FUN',
    titleRef: familyFuneralBenefitTitle,
    flagKey: 'family_funeral_benefit'
  }
]

const activeBenefitList = computed(() => {
  const categories: any[] = props.quote.scheme_categories || []
  if (!categories.length) return BENEFIT_LIST
  return BENEFIT_LIST.filter(({ flagKey }) =>
    categories.some((cat: any) => cat[flagKey])
  )
})

const getExportSummaries = (): any[] => {
  if (resultSummary.value) return [resultSummary.value]
  return validResultSummaries.value as any[]
}

const buildExcelSheet = (
  wb: any,
  code: string,
  title: string,
  summaries: any[]
) => {
  const wsData: string[][] = []
  wsData.push([`${props.quote.scheme_name} — ${title}`, '', '', ''])
  wsData.push(['', '', '', ''])
  wsData.push(['Metric', 'Technical Rate', 'Experience Rated', 'Discounted'])
  summaries.forEach((rs: any) => {
    if (summaries.length > 1)
      wsData.push([`Category: ${rs.category}`, '', '', ''])
    getBenefitRows(rs, code).forEach((r) => wsData.push(r))
    wsData.push(['', '', '', ''])
  })
  const ws = XLSX.utils.aoa_to_sheet(wsData)
  ws['!cols'] = [{ wch: 55 }, { wch: 20 }, { wch: 20 }, { wch: 20 }]
  XLSX.utils.book_append_sheet(wb, ws, code.substring(0, 31))
}

const buildPDFSection = (
  doc: any,
  code: string,
  title: string,
  summaries: any[],
  addPageBreak: boolean
) => {
  if (addPageBreak) doc.addPage()
  doc.setFontSize(13)
  doc.setFont('helvetica', 'bold')
  doc.text(props.quote.scheme_name, 14, 20)
  doc.setFontSize(10)
  doc.text(title, 14, 27)
  let y = 33
  summaries.forEach((rs: any) => {
    if (summaries.length > 1) {
      doc.setFontSize(9)
      doc.setFont('helvetica', 'bold')
      doc.text(`Category: ${rs.category}`, 14, y)
      y += 5
    }
    const body = getBenefitRows(rs, code)
    doc.autoTable({
      startY: y,
      head: [['Metric', 'Technical Rate', 'Experience Rated', 'Discounted']],
      body,
      theme: 'striped',
      styles: { fontSize: 7.5, cellPadding: 1.5 },
      headStyles: {
        fillColor: [70, 100, 120],
        textColor: 255,
        fontStyle: 'bold'
      },
      columnStyles: {
        0: { cellWidth: 90 },
        1: { cellWidth: 32 },
        2: { cellWidth: 32 },
        3: { cellWidth: 32 }
      },
      didParseCell: (data: any) => {
        if (
          data.column.index === 0 &&
          SECTION_LABELS.includes(data.cell.raw as string)
        ) {
          data.cell.styles.fontStyle = 'bold'
          data.cell.styles.fillColor = [210, 215, 220]
        }
      }
    })
    y = (doc as any).lastAutoTable.finalY + 6
  })
}

const exportBenefitToExcel = (code: string, titleRef: any) => {
  const summaries = getExportSummaries()
  if (!summaries.length) return
  const title = titleRef?.value ?? titleRef ?? code
  const wb = XLSX.utils.book_new()
  buildExcelSheet(wb, code, title, summaries)
  XLSX.writeFile(wb, `${code}_summary_${props.quote.scheme_name}.xlsx`)
}

const exportBenefitToPDF = (code: string, titleRef: any) => {
  const summaries = getExportSummaries()
  if (!summaries.length) return
  const title = titleRef?.value ?? titleRef ?? code
  // eslint-disable-next-line new-cap
  const doc: any = new jsPDF({ orientation: 'landscape' })
  buildPDFSection(doc, code, title, summaries, false)
  const pageCount = doc.internal.getNumberOfPages()
  for (let i = 1; i <= pageCount; i++) {
    doc.setPage(i)
    doc.setFontSize(7)
    doc.setFont('helvetica', 'normal')
    doc.text(
      `Page ${i} of ${pageCount}  |  Generated ${new Date().toLocaleDateString()}`,
      14,
      doc.internal.pageSize.getHeight() - 6
    )
  }
  doc.save(`${code}_summary_${props.quote.scheme_name}.pdf`)
}

const exportAllToExcel = () => {
  const summaries = getExportSummaries()
  if (!summaries.length) return
  const wb = XLSX.utils.book_new()
  activeBenefitList.value.forEach(({ code, titleRef }) =>
    buildExcelSheet(wb, code, titleRef.value || code, summaries)
  )
  XLSX.writeFile(wb, `output_summary_${props.quote.scheme_name}.xlsx`)
}

const exportAllToPDF = () => {
  const summaries = getExportSummaries()
  if (!summaries.length) return
  // eslint-disable-next-line new-cap
  const doc: any = new jsPDF({ orientation: 'landscape' })
  activeBenefitList.value.forEach(({ code, titleRef }, i) =>
    buildPDFSection(doc, code, titleRef.value || code, summaries, i > 0)
  )
  const pageCount = doc.internal.getNumberOfPages()
  for (let i = 1; i <= pageCount; i++) {
    doc.setPage(i)
    doc.setFontSize(7)
    doc.setFont('helvetica', 'normal')
    doc.text(
      `Page ${i} of ${pageCount}  |  Generated ${new Date().toLocaleDateString()}`,
      14,
      doc.internal.pageSize.getHeight() - 6
    )
  }
  doc.save(`output_summary_${props.quote.scheme_name}.pdf`)
}

// === END EXPORT UTILITIES ===

onMounted(async () => {
  GroupPricingService.getBenefitMaps().then((res) => {
    benefitMaps.value = res.data

    const glaBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'GLA'
    )
    if (glaBenefit.benefit_alias !== '') {
      glaBenefitTitle.value = glaBenefit.benefit_alias
    }
    const sglaBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'SGLA'
    )
    if (sglaBenefit.benefit_alias !== '') {
      sglaBenefitTitle.value = sglaBenefit.benefit_alias
    }
    const ptdBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'PTD'
    )
    if (ptdBenefit.benefit_alias !== '') {
      ptdBenefitTitle.value = ptdBenefit.benefit_alias
    }
    const ciBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'CI'
    )
    if (ciBenefit.benefit_alias !== '') {
      ciBenefitTitle.value = ciBenefit.benefit_alias
    }
    const phiBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'PHI'
    )
    if (phiBenefit.benefit_alias !== '') {
      phiBenefitTitle.value = phiBenefit.benefit_alias
    }
    const ttdBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'TTD'
    )
    if (ttdBenefit.benefit_alias !== '') {
      ttdBenefitTitle.value = ttdBenefit.benefit_alias
    }
    const familyFuneralBenefit = benefitMaps.value.find(
      (item) => item.benefit_code === 'GFF'
    )
    if (familyFuneralBenefit.benefit_alias !== '') {
      familyFuneralBenefitTitle.value = familyFuneralBenefit.benefit_alias
    }
  })
})
</script>
<style lang="css" scoped></style>
