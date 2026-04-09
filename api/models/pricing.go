package models

import (
	"time"
)

type PricingPoint struct {
	ID                                                   int     `json:"id" gorm:"primary_key"`
	SensitivityID                                        int     `json:"sensitivity_id"`
	JobProductID                                         int     `json:"job_product_id" gorm:"index"`
	ScenarioId                                           int     `json:"scenario_id" gorm:"index"`
	RunId                                                int     `json:"run_id" gorm:"index"`
	RunType                                              int     `json:"run_type" gorm:"index"` // 0 for valuation, 1 for pricing
	ProductCode                                          string  `json:"product_code" gorm:"index"`
	PolicyNumber                                         string  `json:"policy_number" gorm:"index"`
	ProjectionMonth                                      int     `json:"projection_month" gorm:"index"`
	ProjectionYear                                       int     `json:"projection_year"`
	ValuationTimeMonth                                   int     `json:"valuation_time_month"`
	ValuationTimeYear                                    float64 `json:"valuation_time_year"`
	MainMemberAgeNextBirthday                            int     `json:"main_member_age_next_birthday"`
	SpouseAgeNextBirthday                                int     `json:"spouse_age_next_birthday"`
	AgeNextBirthday                                      int     `json:"age_next_birthday"`
	BenefitInForce                                       float64 `json:"benefit_in_force"`
	AccidentProportion                                   float64 `json:"accident_proportion"`
	InflationFactor                                      float64 `json:"inflation_factor"`
	InflationFactorAdjusted                              float64 `json:"inflation_factor_adjusted"`
	PremiumEscalation                                    float64 `json:"premium_escalation"`
	SumAssuredEscalation                                 float64 `json:"sum_assured_escalation"`
	LapseMargin                                          float64 `json:"lapse_margin"`
	PremiumWaiverOnFactor                                float64 `json:"premium_waiver_on_fac"`
	PaidUpOnFactor                                       float64 `json:"paid_up_on_fac"`
	MainMemberMortalityRate                              float64 `json:"main_member_mortality_rate"`
	MainMemberMortalityRateAdjusted                      float64 `json:"main_member_mortality_rate_adjusted"`
	SpouseMortalityRate                                  float64 `json:"spouse_mortality_rate"`
	SpouseMortalityRateAdjusted                          float64 `json:"spouse_mortality_rate_adjusted"`
	ChildMortalityRate                                   float64 `json:"child_mortality_rate"`
	ChildMortalityRateAdjusted                           float64 `json:"child_mortality_rate_adjusted"`
	ChildAccidentalProportion                            float64 `json:"child_accidental_proportion"`
	ChildIndependentMortalityMonthly                     float64 `json:"child_independent_mortality_monthly"`
	ChildIndependentMortalityAdjustedMonthly             float64 `json:"child_independent_mortality_adjusted_monthly"`
	BaseLapse                                            float64 `json:"base_lapse"`
	BaseLapseAdjusted                                    float64 `json:"base_lapse_adjusted"`
	ContractingPartyAlivePortion                         float64 `json:"contracting_party_alive_portion"`
	ContractingPartyAlivePortionAdjusted                 float64 `json:"contracting_party_alive_portion_adjusted"`
	ContractingPartyPolicyLapse                          float64 `json:"contracting_party_policy_lapse"`
	ContractingPartyPolicyLapseAdjusted                  float64 `json:"contracting_party_policy_lapse_adjusted"`
	ChildIndependentLapse                                float64 `json:"child_independent_lapse"`
	ChildIndependentLapseAdjusted                        float64 `json:"child_independent_lapse_adjusted"`
	BaseMortalityRate                                    float64 `json:"base_mortality_rate"`
	BaseMortalityRateAdjusted                            float64 `json:"base_mortality_rate_adjusted"`
	BaseSpouseIndependentLapse                           float64 `json:"base_spouse_independent_lapse"`
	BaseSpouseIndependentLapseAdjusted                   float64 `json:"base_spouse_independent_lapse_adjusted"`
	BaseIndependentRetrenchmentRate                      float64 `json:"base_independent_retrenchment_rate"`
	BaseIndependentRetrenchmentRateAdjusted              float64 `json:"base_independent_retrenchment_rate_adjusted"`
	BaseIndependentDisabilityIncrement                   float64 `json:"base_independent_disability_increment"`
	BaseIndependentDisabilityIncrementAdjusted           float64 `json:"base_independent_disability_increment_adjusted"`
	MainMemberMortalityRateByMonth                       float64 `json:"main_member_mortality_rate_by_month"`
	MainMemberMortalityRateAdjustedByMonth               float64 `json:"main_member_mortality_rate_adjusted_by_month"`
	IndependentSpouseMortalityRateByMonth                float64 `json:"independent_spouse_mortality_rate_by_month"`
	IndependentSpouseMortalityRateAdjustedByMonth        float64 `json:"independent_spouse_mortality_rate_adjusted_by_month"`
	IndependentMortalityRateMonthly                      float64 `json:"independent_mortality_rate_monthly"`
	IndependentMortalityRateAdjustedByMonth              float64 `json:"independent_mortality_rate_adjusted_by_month"`
	IndependentLapseMonthly                              float64 `json:"independent_lapse_monthly"`
	IndependentLapseMonthlyAdjusted                      float64 `json:"independent_lapse_monthly_adjusted"`
	IndependentSpouseLapseMonthly                        float64 `json:"independent_spouse_lapse_monthly"`
	IndependentSpouseLapseMonthlyAdjusted                float64 `json:"independent_spouse_lapse_monthly_adjusted"`
	ChildIndependentLapseMonthly                         float64 `json:"child_independent_lapse_monthly"`
	ChildIndependentLapseAdjustedMonthly                 float64 `json:"child_independent_lapse_adjusted_monthly"`
	IndependentRetrenchmentMonthly                       float64 `json:"independent_retrenchment_monthly"`
	IndependentRetrenchmentMonthlyAdjusted               float64 `json:"independent_retrenchment_monthly_adjusted"`
	IndependentDisabilityMonthly                         float64 `json:"independent_disability_monthly"`
	IndependentDisabilityMonthlyAdjusted                 float64 `json:"independent_disability_monthly_adjusted"`
	MonthlyDependentMortality                            float64 `json:"monthly_dependent_mortality"`
	MonthlyDependentMortalityAdjusted                    float64 `json:"monthly_dependent_mortality_adjusted"`
	MonthlySpouseDependentMortality                      float64 `json:"monthly_spouse_dependent_mortality"`
	MonthlySpouseDependentMortalityAdjusted              float64 `json:"monthly_spouse_dependent_mortality_adjusted"`
	MonthlyChildDependentMortality                       float64 `json:"monthly_child_dependent_mortality"`
	MonthlyChildDependentMortalityAdjusted               float64 `json:"monthly_child_dependent_mortality_adjusted"`
	MonthlyDependentLapse                                float64 `json:"monthly_dependent_lapse"`
	MonthlyDependentLapseAdjusted                        float64 `json:"monthly_dependent_lapse_adjusted"`
	MonthlySpouseDependentLapse                          float64 `json:"monthly_spouse_dependent_lapse"`
	MonthlySpouseDependentLapseAdjusted                  float64 `json:"monthly_spouse_dependent_lapse_adjusted"`
	MonthlyChildDependentLapse                           float64 `json:"monthly_child_dependent_lapse"`
	MonthlyChildDependentLapseAdjusted                   float64 `json:"monthly_child_dependent_lapse_adjusted"`
	MonthlyDependentRetrenchment                         float64 `json:"monthly_dependent_retrenchment"`
	MonthlyDependentRetrenchmentAdjusted                 float64 `json:"monthly_dependent_retrenchment_adjusted"`
	MonthlyDependentDisability                           float64 `json:"monthly_dependent_disability"`
	MonthlyDependentDisabilityAdjusted                   float64 `json:"monthly_dependent_disability_adjusted"`
	SpouseSurvivalRate                                   float64 `json:"spouse_survival_rate"`
	SpouseSurvivalRateAdjusted                           float64 `json:"spouse_survival_rate_adjusted"`
	ChildSurvivalRate                                    float64 `json:"child_survival_rate"`
	ChildSurvivalRateAdjusted                            float64 `json:"child_survival_rate_adjusted"`
	InitialPolicy                                        float64 `gorm:"index" json:"initial_policy"`
	InitialPolicyAdjusted                                float64 `gorm:"index" json:"initial_policy_adjusted"`
	InitialPolicyInclRetrenchmentDecrement               float64 `json:"initial_policy_incl_retrenchment_decrement"`
	InitialPolicyInclRetrenchmentDecrementAdjusted       float64 `json:"initial_policy_incl_retrenchment_decrement_adjusted"`
	SpouseNumberPolicies                                 float64 `json:"spouse_number_policies"`
	SpouseNumberPoliciesAdjusted                         float64 `json:"spouse_number_policies_adjusted"`
	ChildNumberPolicies                                  float64 `json:"child_number_policies"`
	ChildNumberPoliciesAdjusted                          float64 `json:"child_number_policies_adjusted"`
	NumberPaidUp                                         float64 `json:"number_paid_up"`
	NumberPaidUpAdjusted                                 float64 `json:"number_paid_up_adjusted"`
	IncrementalPaidUp                                    float64 `json:"incremental_paid_up"`
	IncrementalPaidUpAdjusted                            float64 `json:"incremental_paid_up_adjusted"`
	PostPaidupPolicyCount                                float64 `json:"post_paidup_policy_count"`
	PostPaidupPolicyCountAdjusted                        float64 `json:"post_paidup_policy_count_adjusted"`
	LastSurvivorPostPaidupPolicyCount                    float64 `json:"last_survivor_post_paidup_policy_count"`
	LastSurvivorPostPaidupPolicyCountAdjusted            float64 `json:"last_survivor_post_paidup_policy_count_adjusted"`
	SpouseNumberOfPaidUps                                float64 `json:"spouse_number_of_paid_ups"`
	SpouseNumberOfPaidUpsAdjusted                        float64 `json:"spouse_number_of_paid_ups_adjusted"`
	ChildNumberOfPaidUps                                 float64 `json:"child_number_of_paid_ups"`
	ChildNumberOfPaidUpsAdjusted                         float64 `json:"child_number_of_paid_ups_adjusted"`
	SpouseNumberOfPremiumWaivers                         float64 `json:"spouse_number_of_premium_waivers"`
	SpouseNumberOfPremiumWaiversAdjusted                 float64 `json:"spouse_number_of_premium_waivers_adjusted"`
	ChildNumberOfPremiumWaivers                          float64 `json:"child_number_of_premium_waivers"`
	ChildNumberOfPremiumWaiversAdjusted                  float64 `json:"child_number_of_premium_waivers_adjusted"`
	SpouseNumberOfPaidUpNaturalDeaths                    float64 `json:"spouse_number_of_paid_up_natural_deaths"`
	SpouseNumberOfPaidUpNaturalDeathsAdjusted            float64 `json:"spouse_number_of_paid_up_natural_deaths_adjusted"`
	ChildNumberOfPaidUpNaturalDeaths                     float64 `json:"child_number_of_paid_up_natural_deaths"`
	ChildNumberOfPaidUpNaturalDeathsAdjusted             float64 `json:"child_number_of_paid_up_natural_deaths_adjusted"`
	ChildNumberOfPremiumWaiverNaturalDeaths              float64 `json:"child_number_of_premium_waiver_natural_deaths"`
	ChildNumberOfPremiumWaiverNaturalDeathsAdjusted      float64 `json:"child_number_of_premium_waiver_natural_deaths_adjusted"`
	SpouseNumberOfPaidUpAccidentalDeaths                 float64 `json:"spouse_number_of_paid_up_accidental_deaths"`
	SpouseNumberOfPaidUpAccidentalDeathsAdjusted         float64 `json:"spouse_number_of_paid_up_accidental_deaths_adjusted"`
	SpouseNumberOfPremiumWaiverNaturalDeaths             float64 `json:"spouse_number_of_premium_waiver_natural_deaths"`
	SpouseNumberOfPremiumWaiverNaturalDeathsAdjusted     float64 `json:"spouse_number_of_premium_waiver_natural_deaths_adjusted"`
	SpouseNumberOfPremiumWaiverAccidentalDeaths          float64 `json:"spouse_number_of_premium_waiver_accidental_deaths"`
	SpouseNumberOfPremiumWaiveAccidentalDeathsAdjusted   float64 `json:"spouse_number_of_premium_waive_accidental_deaths_adjusted"`
	ChildNumberOfPaidUpAccidentalDeaths                  float64 `json:"child_number_of_paid_up_accidental_deaths"`
	ChildNumberOfPaidUpAccidentalDeathsAdjusted          float64 `json:"child_number_of_paid_up_accidental_deaths_adjusted"`
	ChildNumberOfPremiumWaiverAccidentalDeaths           float64 `json:"child_number_of_premium_waiver_accidental_deaths"`
	ChildNumberOfPremiumWaiverAccidentalDeathsAdjusted   float64 `json:"child_number_of_premium_waiver_accidental_deaths_adjusted"`
	IncrementalPremiumWaivers                            float64 `json:"incremental_premium_waivers"`
	IncrementalPremiumWaiversAdjusted                    float64 `json:"incremental_premium_waivers_adjusted"`
	InitialTemporaryPremiumWaivers                       float64 `json:"initial_temporary_premium_waivers"`
	InitialTemporaryPremiumWaiversAdjusted               float64 `json:"initial_temporary_premium_waivers_adjusted"`
	NaturalDeathsInForce                                 float64 `json:"natural_deaths_in_force"`
	NaturalDeathsInForceAdjusted                         float64 `json:"natural_deaths_in_force_adjusted"`
	NaturalDeathsPaidUp                                  float64 `json:"natural_deaths_paid_up"`
	NaturalDeathsPaidUpAdjusted                          float64 `json:"natural_deaths_paid_up_adjusted"`
	NaturalDeathsPremiumWaiver                           float64 `json:"natural_deaths_premium_waiver"`
	NaturalDeathsPremiumWaiverAdjusted                   float64 `json:"natural_deaths_premium_waiver_adjusted"`
	NaturalDeathsTemporaryWaivers                        float64 `json:"natural_deaths_temporary_waivers"`
	NaturalDeathsTemporaryWaiversAdjusted                float64 `json:"natural_deaths_temporary_waivers_adjusted"`
	NumberOfDeathsAccident                               float64 `json:"number_of_deaths_accident"`
	NumberOfDeathsAccidentAdjusted                       float64 `json:"number_of_deaths_accident_adjusted"`
	NumberOfAccidentDeathsPaidUp                         float64 `json:"number_of_accident_deaths_paid_up"`
	NumberOfAccidentDeathsPaidUpAdjusted                 float64 `json:"number_of_accident_deaths_paid_up_adjusted"`
	AccidentDeathsPremiumWaiver                          float64 `json:"accident_deaths_premium_waiver"`
	AccidentDeathsPremiumWaiverAdjusted                  float64 `json:"accident_deaths_premium_waiver_adjusted"`
	AccidentDeathsTemporaryPremiumWaiver                 float64 `json:"accident_deaths_temporary_premium_waiver"`
	AccidentDeathsTemporaryPremiumWaiverAdjusted         float64 `json:"accident_deaths_temporary_premium_waiver_adjusted"`
	AccidentDeathsInForce                                float64 `json:"accident_deaths_in_force"`
	AccidentDeathsInForceAdjusted                        float64 `json:"accident_deaths_in_force_adjusted"`
	NumberOfLapses                                       float64 `json:"number_of_lapses"`
	NumberOfLapsesAdjusted                               float64 `json:"number_of_lapses_adjusted"`
	TotalIncrementalAccidentalDeaths                     float64 `json:"total_incremental_accidental_deaths"`
	TotalIncrementalAccidentalDeathsAdjusted             float64 `json:"total_incremental_accidental_deaths_adjusted"`
	TotalIncrementalPaidupAccidentalDeaths               float64 `json:"total_incremental_paidup_accidental_deaths"`
	TotalIncrementalPaidupAccidentalDeathsAdjusted       float64 `json:"total_incremental_paidup_accidental_deaths_adjusted"`
	TotalSpouseIncrementalAccidentalDeaths               float64 `json:"total_spouse_incremental_accidental_deaths"`
	TotalSpouseIncrementalAccidentalDeathsAdjusted       float64 `json:"total_spouse_incremental_accidental_deaths_adjusted"`
	TotalChildIncrementalAccidentalDeaths                float64 `json:"total_child_incremental_accidental_deaths"`
	TotalChildIncrementalAccidentalDeathsAdjusted        float64 `json:"total_child_incremental_accidental_deaths_adjusted"`
	TotalSpouseIncrementalPaidupAccidentalDeaths         float64 `json:"total_spouse_incremental_paidup_accidental_deaths"`
	TotalSpouseIncrementalPaidupAccidentalDeathsAdjusted float64 `json:"total_spouse_incremental_paidup_accidental_deaths_adjusted"`
	TotalChildIncrementalPaidupAccidentalDeaths          float64 `json:"total_child_incremental_paidup_accidental_deaths"`
	TotalChildIncrementalPaidupAccidentalDeathsAdjusted  float64 `json:"total_child_incremental_paidup_accidental_deaths_adjusted"`
	TotalSpouseIncrementalPwAccidentalDeaths             float64 `json:"total_spouse_incremental_pw_accidental_deaths"`
	TotalSpouseIncrementalPwAccidentalDeathsAdjusted     float64 `json:"total_spouse_incremental_pw_accidental_deaths_adjusted"`
	TotalChildIncrementalPwAccidentalDeaths              float64 `json:"total_child_incremental_pw_accidental_deaths"`
	TotalChildIncrementalPwAccidentalDeathsAdjusted      float64 `json:"total_child_incremental_pw_accidental_deaths_adjusted"`
	TotalIncrementalLapses                               float64 `json:"total_incremental_lapses"`
	TotalIncrementalLapsesAdjusted                       float64 `json:"total_incremental_lapses_adjusted"`
	TotalIncrementalNaturalDeaths                        float64 `json:"total_incremental_natural_deaths"`
	TotalIncrementalNaturalDeathsAdjusted                float64 `json:"total_incremental_natural_deaths_adjusted"`
	TotalIncrementalPaidupNaturalDeaths                  float64 `json:"total_incremental_paidup_natural_deaths"`
	TotalIncrementalPaidupNaturalDeathsAdjusted          float64 `json:"total_incremental_paidup_natural_deaths_adjusted"`
	TotalSpouseIncrementalNaturalDeaths                  float64 `json:"total_spouse_incremental_natural_deaths"`
	TotalSpouseIncrementalNaturalDeathsAdjusted          float64 `json:"total_spouse_incremental_natural_deaths_adjusted"`
	TotalChildIncrementalNaturalDeaths                   float64 `json:"total_child_incremental_natural_deaths"`
	TotalChildIncrementalNaturalDeathsAdjusted           float64 `json:"total_child_incremental_natural_deaths_adjusted"`
	TotalSpouseIncrementalPaidupNaturalDeaths            float64 `json:"total_spouse_incremental_paidup_natural_deaths"`
	TotalSpouseIncrementalPaidupNaturalDeathsAdjusted    float64 `json:"total_spouse_incremental_paidup_natural_deaths_adjusted"`
	TotalChildIncrementalPaidupNaturalDeaths             float64 `json:"total_child_incremental_paidup_natural_deaths"`
	TotalChildIncrementalPaidupNaturalDeathsAdjusted     float64 `json:"total_child_incremental_paidup_natural_deaths_adjusted"`
	TotalSpouseIncrementalPwNaturalDeaths                float64 `json:"total_spouse_incremental_pw_natural_deaths"`
	TotalSpouseIncrementalPwNaturalDeathsAdjusted        float64 `json:"total_spouse_incremental_pw_natural_deaths_adjusted"`
	TotalChildIncrementalPwNaturalDeaths                 float64 `json:"total_child_incremental_pw_natural_deaths"`
	TotalChildIncrementalPwNaturalDeathsAdjusted         float64 `json:"total_child_incremental_pw_natural_deaths_adjusted"`
	NumberOfDisabilities                                 float64 `json:"number_of_disabilities"`
	NumberOfDisabilitiesAdjusted                         float64 `json:"number_of_disabilities_adjusted"`
	NumberOfRetrenchments                                float64 `json:"number_of_retrenchments"`
	NumberOfRetrenchmentsAdjusted                        float64 `json:"number_of_retrenchments_adjusted"`
	NumberOfMaturities                                   float64 `json:"number_of_maturities"`
	NumberOfMaturitiesAdjusted                           float64 `json:"number_of_maturities_adjusted"`
	TotalIncrementalDisabilities                         float64 `json:"total_incremental_disabilities"`
	TotalIncrementalDisabilitiesAdjusted                 float64 `json:"total_incremental_disabilities_adjusted"`
	TotalIncrementalRetrenchments                        float64 `json:"total_incremental_retrenchments"`
	TotalIncrementalRetrenchmentsAdjusted                float64 `json:"total_incremental_retrenchments_adjusted"`
	NumberOfEducatorsInWp                                float64 `json:"number_of_educators_in_wp"`
	NumberOfEducatorsInWpAdjusted                        float64 `json:"number_of_educators_in_wp_adjusted"`
	SumAssured                                           float64 `json:"sum_assured"`
	ChildSumAssured                                      float64 `json:"child_sum_assured"`
	CalculatedInstalment                                 float64 `json:"calculated_instalment"`
	OutstandingSumAssured                                float64 `json:"outstanding_sum_assured"`
	AdditionalSumAssured                                 float64 `json:"additional_sum_assured"`
	ChildAdditionalSumAssured                            float64 `json:"child_additional_sum_assured"`
	Premium                                              float64 `json:"premium"`
	ChildPremium                                         float64 `json:"child_premium"`
	PremiumAdjusted                                      float64 `json:"premium_adjusted"`
	PremiumIncome                                        float64 `json:"premium_income"`
	PremiumIncomeAdjusted                                float64 `json:"premium_income_adjusted"`
	PremiumNotReceived                                   float64 `json:"premium_not_received"`
	PremiumNotReceivedAdjusted                           float64 `json:"premium_not_received_adjusted"`
	Commission                                           float64 `json:"commission"`
	CommissionAdjusted                                   float64 `json:"commission_adjusted"`
	ClawBack                                             float64 `json:"claw_back"`
	ClawBackAdjusted                                     float64 `json:"claw_back_adjusted"`
	DeathOutgo                                           float64 `json:"death_outgo"`
	DeathOutgoAdjusted                                   float64 `json:"death_outgo_adjusted"`
	DisabilityOutgo                                      float64 `json:"disability_outgo"`
	DisabilityOutgoAdjusted                              float64 `json:"disability_outgo_adjusted"`
	RetrenchmentOutgo                                    float64 `json:"retrenchment_outgo"`
	RetrenchmentOutgoAdjusted                            float64 `json:"retrenchment_outgo_adjusted"`
	ChildDeathOutgo                                      float64 `json:"child_death_outgo"`
	ChildDeathOutgoAdjusted                              float64 `json:"child_death_outgo_adjusted"`
	SpouseDeathOutgo                                     float64 `json:"spouse_death_outgo"`
	SpouseDeathOutgoAdjusted                             float64 `json:"spouse_death_outgo_adjusted"`
	AccidentalDeathOutgo                                 float64 `json:"accidental_death_outgoing"`
	AccidentalDeathOutgoAdjusted                         float64 `json:"accidental_death_outgoing_adjusted"`
	ChildAccidentalDeathOutgo                            float64 `json:"child_accidental_death_outgo"`
	ChildAccidentalDeathOutgoAdjusted                    float64 `json:"child_accidental_death_outgo_adjusted"`
	SpouseAccidentalDeathOutgo                           float64 `json:"spouse_accidental_death_outgo"`
	SpouseAccidentalDeathOutgoAdjusted                   float64 `json:"spouse_accidental_death_outgo_adjusted"`
	Educator                                             float64 `json:"educator"`
	EducatorAdjusted                                     float64 `json:"educator_adjusted"`
	CashBackOnSurvival                                   float64 `json:"cash_back_on_survival"`
	CashBackOnSurvivalAdjusted                           float64 `json:"cash_back_on_survival_adjusted"`
	CashBackOnDeath                                      float64 `json:"cash_back_on_death"`
	CashBackOnDeathAdjusted                              float64 `json:"cash_back_on_death_adjusted"`
	RiderFuneral                                         float64 `json:"rider_funeral"`
	RiderFuneralAdjusted                                 float64 `json:"rider_funeral_adjusted"`
	Rider                                                float64 `json:"rider"`
	RiderAdjusted                                        float64 `json:"rider_adjusted"`
	Expenses                                             float64 `json:"expenses"`
	ExpensesAdjusted                                     float64 `json:"expenses_adjusted"`
	NetCashFlow                                          float64 `gorm:"index" json:"net_cash_flow"`
	NetCashFlowAdjusted                                  float64 `gorm:"index" json:"net_cash_flow_adjusted"`
	Reserves                                             float64 `gorm:"index" json:"reserves"`
	ReservesAdjusted                                     float64 `gorm:"index" json:"reserves_adjusted"`
	ChangeInReserves                                     float64 `gorm:"index" json:"change_in_reserves"`
	ChangeInReservesAdjusted                             float64 `gorm:"index" json:"change_in_reserves_adjusted"`
	InvestmentIncome                                     float64 `json:"investment_income"`
	InvestmentIncomeAdjusted                             float64 `json:"investment_income_adjusted"`
	Profit                                               float64 `json:"profit"`
	ProfitAdjusted                                       float64 `json:"profit_adjusted"`
	DiscountedPremiumIncome                              float64 `json:"discounted_premium_income"`
	DiscountedPremiumIncomeAdjusted                      float64 `json:"discounted_premium_income_adjusted"`
	DiscountedPremiumNotReceived                         float64 `json:"discounted_premium_not_received"`
	DiscountedPremiumNotReceivedAdjusted                 float64 `json:"discounted_premium_not_received_adjusted"`
	DiscountedCommission                                 float64 `json:"discounted_commission"`
	DiscountedCommissionAdjusted                         float64 `json:"discounted_commission_adjusted"`
	DiscountedClawBack                                   float64 `json:"discounted_claw_back"`
	DiscountedClawBackAdjusted                           float64 `json:"discounted_claw_back_adjusted"`
	DiscountedDeathOutgo                                 float64 `json:"discounted_death_outgoing"`
	DiscountedDeathOutgoAdjusted                         float64 `json:"discounted_death_outgoing_adjusted"`
	DiscountedDisabilityOutgo                            float64 `json:"discounted_disability_outgo"`
	DiscountedDisabilityOutgoAdjusted                    float64 `json:"discounted_disability_outgo_adjusted"`
	DiscountedRetrenchmentOutgo                          float64 `json:"discounted_retrenchment_outgo"`
	DiscountedRetrenchmentOutgoAdjusted                  float64 `json:"discounted_retrenchment_outgo_adjusted"`
	ChildDiscountedDeathOutgo                            float64 `json:"child_discounted_death_outgo"`
	ChildDiscountedDeathOutgoAdjusted                    float64 `json:"child_discounted_death_outgo_adjusted"`
	SpouseDiscountedDeathOutgo                           float64 `json:"spouse_discounted_death_outgo"`
	SpouseDiscountedDeathOutgoAdjusted                   float64 `json:"spouse_discounted_death_outgo_adjusted"`
	DiscountedAccidentalDeathOutgo                       float64 `json:"discounted_accidental_death_outgoing"`
	DiscountedAccidentalDeathOutgoAdjusted               float64 `json:"discounted_accidental_death_outgoing_adjusted"`
	ChildDiscountedAccidentalDeathOutgo                  float64 `json:"child_discounted_accidental_death_outgo"`
	ChildDiscountedAccidentalDeathOutgoAdjusted          float64 `json:"child_discounted_accidental_death_outgo_adjusted"`
	SpouseDiscountedAccidentalDeathOutgo                 float64 `json:"spouse_discounted_accidental_death_outgo"`
	SpouseDiscountedAccidentalDeathOutgoAdjusted         float64 `json:"spouse_discounted_accidental_death_outgo_adjusted"`
	DiscountedEducator                                   float64 `json:"discounted_educator"`
	DiscountedEducatorAdjusted                           float64 `json:"discounted_educator_adjusted"`
	DiscountedCashBackOnSurvival                         float64 `json:"discounted_cash_back_on_survival"`
	DiscountedCashBackOnSurvivalAdjusted                 float64 `json:"discounted_cash_back_on_survival_adjusted"`
	DiscountedCashBackOnDeath                            float64 `json:"discounted_cash_back_on_death"`
	DiscountedCashBackOnDeathAdjusted                    float64 `json:"discounted_cash_back_on_death_adjusted"`
	DiscountedRiderFuneral                               float64 `json:"discounted_rider_funeral"`
	DiscountedRiderFuneralAdjusted                       float64 `json:"discounted_rider_funeral_adjusted"`
	DiscountedRider                                      float64 `json:"discounted_rider"`
	DiscountedRiderAdjusted                              float64 `json:"discounted_rider_adjusted"`
	DiscountedExpenses                                   float64 `json:"discounted_expenses"`
	DiscountedExpensesAdjusted                           float64 `json:"discounted_expenses_adjusted"`
	DiscountedInvestmentIncome                           float64 `json:"discounted_investment_income"`
	DiscountedInvestmentIncomeAdjusted                   float64 `json:"discounted_investment_income_adjusted"`
	AnnuityFactor                                        float64 `json:"annuity_factor"`
	RenewalCommissionAnnuityFeeder                       float64 `json:"renewal_commission_annuity_feeder"`
	DiscountedRenewalCommissionAnnuityFeeder             float64 `json:"discounted_renewal_commission_annuity_feeder"`
	RenewalCommissionAnnuityFactor                       float64 `json:"renewal_commission_annuity_factor"`
	ChildAnnuityFactor                                   float64 `json:"child_annuity_factor"`
	DiscountedAnnuityFactor                              float64 `json:"discounted_annuity_factor"`
	DiscountedProfit                                     float64 `json:"discounted_profit"`
	DiscountedProfitAdjusted                             float64 `json:"discounted_profit_adjusted"`
	DiscountFactor                                       float64 `json:"discount_factor"`
	DiscountFactorAdjusted                               float64 `json:"discount_factor_adjusted"`
	ProfDiscountedPremium                                float64 `json:"prof_discounted_premium"`
	ProfDiscountedPremiumNotReceived                     float64 `json:"prof_discounted_premium_not_received"`
	ProfDiscountedDeath                                  float64 `json:"prof_discounted_death"`
	ProfDiscountedDisability                             float64 `json:"prof_discounted_disability"`
	ProfDiscountedRetrenchment                           float64 `json:"prof_discounted_retrenchment"`
	ProfDiscountedCommission                             float64 `json:"prof_discounted_commission"`
	ProfDiscountedClawback                               float64 `json:"prof_discounted_clawback"`
	ProfDiscountedAccidentalDeath                        float64 `json:"prof_discounted_accidental_death"`
	ProfDiscountedExpenses                               float64 `json:"prof_discounted_expenses"`
	ProfDiscountedChangeInReserve                        float64 `json:"prof_discounted_change_in_reserve"`
	ProfDiscountedInvestmentIncome                       float64 `json:"prof_discounted_investment_income"`
	ProfSpouseDiscountedDeath                            float64 `json:"prof_spouse_discounted_death"`
	ProfSpouseDiscountedAccidentalDeath                  float64 `json:"prof_spouse_discounted_accidental_death"`
	ProfChildDiscountedDeath                             float64 `json:"prof_child_discounted_death"`
	ProfChildDiscountedAccidentalDeath                   float64 `json:"prof_child_discounted_accidental_death"`
	ProfDiscountedEducator                               float64 `json:"prof_discounted_educator"`
	ProfDiscountedCashBackOnSurvival                     float64 `json:"prof_discounted_cash_back_on_survival"`
	ProfDiscountedCashBackOnDeath                        float64 `json:"prof_discounted_cash_back_on_death"`
	ProfDiscountedRiderFuneral                           float64 `json:"prof_discounted_rider_funeral"`
	ProfDiscountedRider                                  float64 `json:"prof_discounted_rider"`
	ProfRiskAdjustment                                   float64 `json:"prof_risk_adjustment"`
	ProfCSM                                              float64 `json:"prof_csm"`
	ProfLossComponent                                    float64 `json:"prof_loss_component"`
}

type AggregatedPricingPoint struct {
	ID                                             int     `json:"-" gorm:"primary_key"`
	SensitivityID                                  int     `json:"-"`
	JobProductID                                   int     `json:"-" gorm:"index"`
	ScenarioId                                     int     `json:"-" gorm:"index"`
	RunId                                          int     `json:"-"`
	RunType                                        int     `json:"-" gorm:"index"` // 0 for valuation, 1 for pricing
	ProductCode                                    string  `json:"product_code" gorm:"index"`
	PolicyNumber                                   string  `json:"policy_number" gorm:"index"`
	ProjectionMonth                                int     `json:"projection_month" gorm:"index"`
	ProjectionYear                                 int     `json:"projection_year"`
	ValuationTimeMonth                             int     `json:"valuation_time_month"`
	ValuationTimeYear                              float64 `json:"valuation_time_year"`
	SpouseSurvivalRate                             float64 `json:"spouse_survival_rate"`
	SpouseSurvivalRateAdjusted                     float64 `json:"spouse_survival_rate_adjusted"`
	ChildSurvivalRate                              float64 `json:"child_survival_rate"`
	ChildSurvivalRateAdjusted                      float64 `json:"child_survival_rate_adjusted"`
	InitialPolicy                                  float64 `gorm:"index" json:"initial_policy"`
	InitialPolicyAdjusted                          float64 `gorm:"index" json:"initial_policy_adjusted"`
	InitialPolicyInclRetrenchmentDecrement         float64 `json:"initial_policy_incl_retrenchment_decrement"`
	InitialPolicyInclRetrenchmentDecrementAdjusted float64 `json:"initial_policy_incl_retrenchment_decrement_adjusted"`
	SpouseNumberPolicies                           float64 `json:"spouse_number_policies"`
	SpouseNumberPoliciesAdjusted                   float64 `json:"spouse_number_policies_adjusted"`
	ChildNumberPolicies                            float64 `json:"child_number_policies"`
	ChildNumberPoliciesAdjusted                    float64 `json:"child_number_policies_adjusted"`
	NumberPaidUp                                   float64 `json:"number_paid_up"`
	NumberPaidUpAdjusted                           float64 `json:"number_paid_up_adjusted"`
	IncrementalPaidUp                              float64 `json:"incremental_paid_up"`
	IncrementalPaidUpAdjusted                      float64 `json:"incremental_paid_up_adjusted"`
	PostPaidupPolicyCount                          float64 `json:"post_paidup_policy_count"`
	PostPaidupPolicyCountAdjusted                  float64 `json:"post_paidup_policy_count_adjusted"`
	LastSurvivorPostPaidupPolicyCount              float64 `json:"last_survivor_post_paidup_policy_count"`
	LastSurvivorPostPaidupPolicyCountAdjusted      float64 `json:"last_survivor_post_paidup_policy_count_adjusted"`
	SpouseNumberOfPaidUps                          float64 `json:"spouse_number_of_paid_ups"`
	SpouseNumberOfPaidUpsAdjusted                  float64 `json:"spouse_number_of_paid_ups_adjusted"`
	ChildNumberOfPaidUps                           float64 `json:"child_number_of_paid_ups"`
	ChildNumberOfPaidUpsAdjusted                   float64 `json:"child_number_of_paid_ups_adjusted"`
	SpouseNumberOfPremiumWaivers                   float64 `json:"spouse_number_of_premium_waivers"`
	SpouseNumberOfPremiumWaiversAdjusted           float64 `json:"spouse_number_of_premium_waivers_adjusted"`
	ChildNumberOfPremiumWaivers                    float64 `json:"child_number_of_premium_waivers"`
	ChildNumberOfPremiumWaiversAdjusted            float64 `json:"child_number_of_premium_waivers_adjusted"`
	SpouseNumberOfPaidUpNaturalDeaths              float64 `json:"spouse_number_of_paid_up_natural_deaths"`
	SpouseNumberOfPaidUpNaturalDeathsAdjusted      float64 `json:"spouse_number_of_paid_up_natural_deaths_adjusted"`
	ChildNumberOfPaidUpNaturalDeaths               float64 `json:"child_number_of_paid_up_natural_deaths"`
	ChildNumberOfPaidUpNaturalDeathsAdjusted       float64 `json:"child_number_of_paid_up_natural_deaths_adjusted"`
	SpouseNumberOfPaidUpAccidentalDeaths           float64 `json:"spouse_number_of_paid_up_accidental_deaths"`
	SpouseNumberOfPaidUpAccidentalDeathsAdjusted   float64 `json:"spouse_number_of_paid_up_accidental_deaths_adjusted"`
	ChildNumberOfPaidUpAccidentalDeaths            float64 `json:"child_number_of_paid_up_accidental_deaths"`
	ChildNumberOfPaidUpAccidentalDeathsAdjusted    float64 `json:"child_number_of_paid_up_accidental_deaths_adjusted"`
	IncrementalPremiumWaivers                      float64 `json:"incremental_premium_waivers"`
	IncrementalPremiumWaiversAdjusted              float64 `json:"incremental_premium_waivers_adjusted"`
	InitialTemporaryPremiumWaivers                 float64 `json:"initial_temporary_premium_waivers"`
	InitialTemporaryPremiumWaiversAdjusted         float64 `json:"initial_temporary_premium_waivers_adjusted"`
	NaturalDeathsInForce                           float64 `json:"natural_deaths_in_force"`
	NaturalDeathsInForceAdjusted                   float64 `json:"natural_deaths_in_force_adjusted"`
	NaturalDeathsPaidUp                            float64 `json:"natural_deaths_paid_up"`
	NaturalDeathsPaidUpAdjusted                    float64 `json:"natural_deaths_paid_up_adjusted"`
	NaturalDeathsPremiumWaiver                     float64 `json:"natural_deaths_premium_waiver"`
	NaturalDeathsPremiumWaiverAdjusted             float64 `json:"natural_deaths_premium_waiver_adjusted"`
	NaturalDeathsTemporaryWaivers                  float64 `json:"natural_deaths_temporary_waivers"`
	NaturalDeathsTemporaryWaiversAdjusted          float64 `json:"natural_deaths_temporary_waivers_adjusted"`
	NumberOfDeathsAccident                         float64 `json:"number_of_deaths_accident"`
	NumberOfDeathsAccidentAdjusted                 float64 `json:"number_of_deaths_accident_adjusted"`
	NumberOfAccidentDeathsPaidUp                   float64 `json:"number_of_accident_deaths_paid_up"`
	NumberOfAccidentDeathsPaidUpAdjusted           float64 `json:"number_of_accident_deaths_paid_up_adjusted"`
	AccidentDeathsPremiumWaiver                    float64 `json:"accident_deaths_premium_waiver"`
	AccidentDeathsPremiumWaiverAdjusted            float64 `json:"accident_deaths_premium_waiver_adjusted"`
	AccidentDeathsTemporaryPremiumWaiver           float64 `json:"accident_deaths_temporary_premium_waiver"`
	AccidentDeathsTemporaryPremiumWaiverAdjusted   float64 `json:"accident_deaths_temporary_premium_waiver_adjusted"`
	AccidentDeathsInForce                          float64 `json:"accident_deaths_in_force"`
	AccidentDeathsInForceAdjusted                  float64 `json:"accident_deaths_in_force_adjusted"`
	NumberOfLapses                                 float64 `json:"number_of_lapses"`
	NumberOfLapsesAdjusted                         float64 `json:"number_of_lapses_adjusted"`
	TotalIncrementalAccidentalDeaths               float64 `json:"total_incremental_accidental_deaths"`
	TotalIncrementalAccidentalDeathsAdjusted       float64 `json:"total_incremental_accidental_deaths_adjusted"`
	TotalSpouseIncrementalAccidentalDeaths         float64 `json:"total_spouse_incremental_accidental_deaths"`
	TotalSpouseIncrementalAccidentalDeathsAdjusted float64 `json:"total_spouse_incremental_accidental_deaths_adjusted"`
	TotalChildIncrementalAccidentalDeaths          float64 `json:"total_child_incremental_accidental_deaths"`
	TotalChildIncrementalAccidentalDeathsAdjusted  float64 `json:"total_child_incremental_accidental_deaths_adjusted"`
	TotalIncrementalLapses                         float64 `json:"total_incremental_lapses"`
	TotalIncrementalLapsesAdjusted                 float64 `json:"total_incremental_lapses_adjusted"`
	TotalIncrementalNaturalDeaths                  float64 `json:"total_incremental_natural_deaths"`
	TotalIncrementalNaturalDeathsAdjusted          float64 `json:"total_incremental_natural_deaths_adjusted"`
	TotalSpouseIncrementalNaturalDeaths            float64 `json:"total_spouse_incremental_natural_deaths"`
	TotalSpouseIncrementalNaturalDeathsAdjusted    float64 `json:"total_spouse_incremental_natural_deaths_adjusted"`
	TotalChildIncrementalNaturalDeaths             float64 `json:"total_child_incremental_natural_deaths"`
	TotalChildIncrementalNaturalDeathsAdjusted     float64 `json:"total_child_incremental_natural_deaths_adjusted"`
	NumberOfDisabilities                           float64 `json:"number_of_disabilities"`
	NumberOfDisabilitiesAdjusted                   float64 `json:"number_of_disabilities_adjusted"`
	NumberOfRetrenchments                          float64 `json:"number_of_retrenchments"`
	NumberOfRetrenchmentsAdjusted                  float64 `json:"number_of_retrenchments_adjusted"`
	NumberOfMaturities                             float64 `json:"number_of_maturities"`
	NumberOfMaturitiesAdjusted                     float64 `json:"number_of_maturities_adjusted"`
	TotalIncrementalDisabilities                   float64 `json:"total_incremental_disabilities"`
	TotalIncrementalDisabilitiesAdjusted           float64 `json:"total_incremental_disabilities_adjusted"`
	TotalIncrementalRetrenchments                  float64 `json:"total_incremental_retrenchments"`
	TotalIncrementalRetrenchmentsAdjusted          float64 `json:"total_incremental_retrenchments_adjusted"`
	NumberOfEducatorsInWp                          float64 `json:"number_of_educators_in_wp"`
	NumberOfEducatorsInWpAdjusted                  float64 `json:"number_of_educators_in_wp_adjusted"`
	SumAssured                                     float64 `json:"sum_assured"`
	ChildSumAssured                                float64 `json:"child_sum_assured"`
	CalculatedInstalment                           float64 `json:"calculated_instalment"`
	OutstandingSumAssured                          float64 `json:"outstanding_sum_assured"`
	AdditionalSumAssured                           float64 `json:"additional_sum_assured"`
	ChildAdditionalSumAssured                      float64 `json:"child_additional_sum_assured"`
	Premium                                        float64 `json:"premium"`
	ChildPremium                                   float64 `json:"child_premium"`
	PremiumAdjusted                                float64 `json:"premium_adjusted"`
	PremiumIncome                                  float64 `json:"premium_income"`
	PremiumIncomeAdjusted                          float64 `json:"premium_income_adjusted"`
	PremiumNotReceived                             float64 `json:"premium_not_received"`
	PremiumNotReceivedAdjusted                     float64 `json:"premium_not_received_adjusted"`
	Commission                                     float64 `json:"commission"`
	CommissionAdjusted                             float64 `json:"commission_adjusted"`
	ClawBack                                       float64 `json:"claw_back"`
	ClawBackAdjusted                               float64 `json:"claw_back_adjusted"`
	DeathOutgo                                     float64 `json:"death_outgo"`
	DeathOutgoAdjusted                             float64 `json:"death_outgo_adjusted"`
	DisabilityOutgo                                float64 `json:"disability_outgo"`
	DisabilityOutgoAdjusted                        float64 `json:"disability_outgo_adjusted"`
	RetrenchmentOutgo                              float64 `json:"retrenchment_outgo"`
	RetrenchmentOutgoAdjusted                      float64 `json:"retrenchment_outgo_adjusted"`
	ChildDeathOutgo                                float64 `json:"child_death_outgo"`
	ChildDeathOutgoAdjusted                        float64 `json:"child_death_outgo_adjusted"`
	SpouseDeathOutgo                               float64 `json:"spouse_death_outgo"`
	SpouseDeathOutgoAdjusted                       float64 `json:"spouse_death_outgo_adjusted"`
	AccidentalDeathOutgo                           float64 `json:"accidental_death_outgoing"`
	AccidentalDeathOutgoAdjusted                   float64 `json:"accidental_death_outgoing_adjusted"`
	ChildAccidentalDeathOutgo                      float64 `json:"child_accidental_death_outgo"`
	ChildAccidentalDeathOutgoAdjusted              float64 `json:"child_accidental_death_outgo_adjusted"`
	SpouseAccidentalDeathOutgo                     float64 `json:"spouse_accidental_death_outgo"`
	SpouseAccidentalDeathOutgoAdjusted             float64 `json:"spouse_accidental_death_outgo_adjusted"`
	Educator                                       float64 `json:"educator"`
	EducatorAdjusted                               float64 `json:"educator_adjusted"`
	CashBackOnSurvival                             float64 `json:"cash_back_on_survival"`
	CashBackOnSurvivalAdjusted                     float64 `json:"cash_back_on_survival_adjusted"`
	CashBackOnDeath                                float64 `json:"cash_back_on_death"`
	CashBackOnDeathAdjusted                        float64 `json:"cash_back_on_death_adjusted"`
	RiderFuneral                                   float64 `json:"rider_funeral"`
	RiderFuneralAdjusted                           float64 `json:"rider_funeral_adjusted"`
	Rider                                          float64 `json:"rider"`
	RiderAdjusted                                  float64 `json:"rider_adjusted"`
	Expenses                                       float64 `json:"expenses"`
	ExpensesAdjusted                               float64 `json:"expenses_adjusted"`
	NetCashFlow                                    float64 `gorm:"index" json:"net_cash_flow"`
	NetCashFlowAdjusted                            float64 `gorm:"index" json:"net_cash_flow_adjusted"`
	Reserves                                       float64 `gorm:"index" json:"reserves"`
	ReservesAdjusted                               float64 `gorm:"index" json:"reserves_adjusted"`
	ChangeInReserves                               float64 `gorm:"index" json:"change_in_reserves"`
	ChangeInReservesAdjusted                       float64 `gorm:"index" json:"change_in_reserves_adjusted"`
	InvestmentIncome                               float64 `json:"investment_income"`
	InvestmentIncomeAdjusted                       float64 `json:"investment_income_adjusted"`
	Profit                                         float64 `json:"profit"`
	ProfitAdjusted                                 float64 `json:"profit_adjusted"`
	DiscountedPremiumIncome                        float64 `json:"discounted_premium_income"`
	DiscountedPremiumIncomeAdjusted                float64 `json:"discounted_premium_income_adjusted"`
	DiscountedPremiumNotReceived                   float64 `json:"discounted_premium_not_received"`
	DiscountedPremiumNotReceivedAdjusted           float64 `json:"discounted_premium_not_received_adjusted"`
	DiscountedCommission                           float64 `json:"discounted_commission"`
	DiscountedCommissionAdjusted                   float64 `json:"discounted_commission_adjusted"`
	DiscountedClawBack                             float64 `json:"discounted_claw_back"`
	DiscountedClawBackAdjusted                     float64 `json:"discounted_claw_back_adjusted"`
	DiscountedDeathOutgo                           float64 `json:"discounted_death_outgoing"`
	DiscountedDeathOutgoAdjusted                   float64 `json:"discounted_death_outgoing_adjusted"`
	DiscountedDisabilityOutgo                      float64 `json:"discounted_disability_outgo"`
	DiscountedDisabilityOutgoAdjusted              float64 `json:"discounted_disability_outgo_adjusted"`
	DiscountedRetrenchmentOutgo                    float64 `json:"discounted_retrenchment_outgo"`
	DiscountedRetrenchmentOutgoAdjusted            float64 `json:"discounted_retrenchment_outgo_adjusted"`
	ChildDiscountedDeathOutgo                      float64 `json:"child_discounted_death_outgo"`
	ChildDiscountedDeathOutgoAdjusted              float64 `json:"child_discounted_death_outgo_adjusted"`
	SpouseDiscountedDeathOutgo                     float64 `json:"spouse_discounted_death_outgo"`
	SpouseDiscountedDeathOutgoAdjusted             float64 `json:"spouse_discounted_death_outgo_adjusted"`
	DiscountedAccidentalDeathOutgo                 float64 `json:"discounted_accidental_death_outgoing"`
	DiscountedAccidentalDeathOutgoAdjusted         float64 `json:"discounted_accidental_death_outgoing_adjusted"`
	ChildDiscountedAccidentalDeathOutgo            float64 `json:"child_discounted_accidental_death_outgo"`
	ChildDiscountedAccidentalDeathOutgoAdjusted    float64 `json:"child_discounted_accidental_death_outgo_adjusted"`
	SpouseDiscountedAccidentalDeathOutgo           float64 `json:"spouse_discounted_accidental_death_outgo"`
	SpouseDiscountedAccidentalDeathOutgoAdjusted   float64 `json:"spouse_discounted_accidental_death_outgo_adjusted"`
	DiscountedEducator                             float64 `json:"discounted_educator"`
	DiscountedEducatorAdjusted                     float64 `json:"discounted_educator_adjusted"`
	DiscountedCashBackOnSurvival                   float64 `json:"discounted_cash_back_on_survival"`
	DiscountedCashBackOnSurvivalAdjusted           float64 `json:"discounted_cash_back_on_survival_adjusted"`
	DiscountedCashBackOnDeath                      float64 `json:"discounted_cash_back_on_death"`
	DiscountedCashBackOnDeathAdjusted              float64 `json:"discounted_cash_back_on_death_adjusted"`
	DiscountedRiderFuneral                         float64 `json:"discounted_rider_funeral"`
	DiscountedRiderFuneralAdjusted                 float64 `json:"discounted_rider_funeral_adjusted"`
	DiscountedRider                                float64 `json:"discounted_rider"`
	DiscountedRiderAdjusted                        float64 `json:"discounted_rider_adjusted"`
	DiscountedExpenses                             float64 `json:"discounted_expenses"`
	DiscountedExpensesAdjusted                     float64 `json:"discounted_expenses_adjusted"`
	DiscountedInvestmentIncome                     float64 `json:"discounted_investment_income"`
	DiscountedInvestmentIncomeAdjusted             float64 `json:"discounted_investment_income_adjusted"`
	DiscountedProfit                               float64 `json:"discounted_profit"`
	DiscountedProfitAdjusted                       float64 `json:"discounted_profit_adjusted"`
	ProfDiscountedPremium                          float64 `json:"prof_discounted_premium"`
	ProfDiscountedPremiumNotReceived               float64 `json:"prof_discounted_premium_not_received"`
	ProfDiscountedDeath                            float64 `json:"prof_discounted_death"`
	ProfDiscountedDisability                       float64 `json:"prof_discounted_disability"`
	ProfDiscountedRetrenchment                     float64 `json:"prof_discounted_retrenchment"`
	ProfDiscountedCommission                       float64 `json:"prof_discounted_commission"`
	ProfDiscountedClawback                         float64 `json:"prof_discounted_clawback"`
	ProfDiscountedAccidentalDeath                  float64 `json:"prof_discounted_accidental_death"`
	ProfDiscountedExpenses                         float64 `json:"prof_discounted_expenses"`
	ProfDiscountedChangeInReserve                  float64 `json:"prof_discounted_change_in_reserve"`
	ProfDiscountedInvestmentIncome                 float64 `json:"prof_discounted_investment_income"`
	ProfSpouseDiscountedDeath                      float64 `json:"prof_spouse_discounted_death"`
	ProfSpouseDiscountedAccidentalDeath            float64 `json:"prof_spouse_discounted_accidental_death"`
	ProfChildDiscountedDeath                       float64 `json:"prof_child_discounted_death"`
	ProfChildDiscountedAccidentalDeath             float64 `json:"prof_child_discounted_accidental_death"`
	ProfDiscountedEducator                         float64 `json:"prof_discounted_educator"`
	ProfDiscountedCashBackOnSurvival               float64 `json:"prof_discounted_cash_back_on_survival"`
	ProfDiscountedCashBackOnDeath                  float64 `json:"prof_discounted_cash_back_on_death"`
	ProfDiscountedRiderFuneral                     float64 `json:"prof_discounted_rider_funeral"`
	ProfDiscountedRider                            float64 `json:"prof_discounted_rider"`
	ProfRiskAdjustment                             float64 `json:"prof_risk_adjustment"`
	ProfCSM                                        float64 `json:"prof_csm"`
	ProfLossComponent                              float64 `json:"prof_loss_component"`
}

type Indicator struct {
	Name  string
	Value bool
}

type PricingRun struct {
	ID              int             `json:"id" gorm:"primary_key"`
	Name            string          `json:"name"`
	ProductId       int             `json:"product_id"`
	PricingConfig   []PricingConfig `json:"pricing_config" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RunSingle       bool            `json:"run_single"`
	RunGoalSeek     bool            `json:"run_goal_seek"`
	ProfitSignature bool            `json:"profit_signature"`
	Status          string          `json:"status"`
	FailureReason   string          `json:"failure_reason"`
	Progress        float64         `json:"progress"`
	RunTime         float64         `json:"run_time"`
	RunDate         time.Time       `json:"run_date"`
	User            string          `json:"user"`
	UserEmail       string          `json:"user_email"`
}

type PricingConfig struct {
	ID                   int    `gorm:"primary_key" json:"id"`
	PricingRunID         int    `json:"pricing_run_id"`
	Description          string `json:"description"`
	SpouseIndicator      bool   `json:"spouse_ind"`
	ChildIndicator       bool   `json:"child_ind"`
	EducatorIndicator    bool   `json:"educator_ind"`
	CashBackOnDeath      bool   `json:"cashback_on_death"`
	CashBackOnSurvival   bool   `json:"cashback_on_survival"`
	CashBack             bool   `json:"cash_back"`
	PaidupOnSurvival     bool   `json:"paidup_on_survival"`
	PremiumWaiverOnDeath bool   `json:"premium_waiver_on_death"`
	Cow                  bool   `json:"cow"`
	Grocery              bool   `json:"grocery"`
	Tombstone            bool   `json:"tombstone"`
	Repatriation         bool   `json:"repatriation"`
	Death                bool   `json:"death"`
	PermDisability       bool   `json:"perm_disability"`
	CriticalIllness      bool   `json:"critical_illness"`
	TempDisability       bool   `json:"temp_disability"`
	Retrenchment         bool   `json:"retrenchment"`
	Funeral              bool   `json:"funeral"`
	ParameterBasis       string `json:"parameter_basis"`
	ShockBasis           string `json:"shock_basis"`
	MpVersion            string `json:"mp_version"`
}

type PricingParameter struct {
	ProductCode                  string  `json:"product_code" csv:"product_code"`
	InitialCommissionPercentage1 float64 `json:"initial_commission_percentage1" csv:"initial_commission_percentage1"`
	InitialCommissionPercentage2 float64 `json:"initial_commission_percentage2" csv:"initial_commission_percentage2"`
	InitialCommissionRand        float64 `json:"initial_commission_rand" csv:"initial_commission_rand"`
	ClawbackPeriod               int     `json:"clawback_period" csv:"clawback_period"`
	RenewalCommissionPercentage  float64 `json:"renewal_commission_percentage" csv:"renewal_commission_percentage"`
	RenewalCommissionRand        float64 `json:"renewal_commission_rand" csv:"renewal_commission_rand"`
	HybridRenewalCommStartM      int     `json:"hybrid_renewal_comm_start_m" csv:"hybrid_renewal_comm_start_m"`
	HybridRenewalCommEndM        int     `json:"hybrid_renewal_comm_end_m" csv:"hybrid_renewal_comm_end_m"`
	NumberChildPerPolicy         float64 `json:"number_child_per_policy" csv:"number_child_per_policy"`
	AverageAgeAtEntryPerChild    int     `json:"average_age_at_entry_per_child" csv:"average_age_at_entry_per_child"`
	EducatorSumAssured           float64 `json:"educator_sum_assured" csv:"educator_sum_assured"`
	DistributionMale             float64 `json:"distribution_male" csv:"distribution_male"`
	DistributionFemale           float64 `json:"distribution_female" csv:"distribution_female"`
	TargetProfitMargin           float64 `json:"target_profit_margin" csv:"target_profit_margin"`
	GoalSeekStep                 float64 `json:"goal_seek_step" csv:"goal_seek_step"`
	GoalSeekMaxIterations        int     `json:"goal_seek_max_iterations" csv:"goal_seek_max_iterations"`
	YieldCurveCode               string  `json:"yield_curve_code" csv:"yield_curve_code"`
	RiskAdjustmentProp           float64 `json:"risk_adjustment_prop" csv:"risk_adjustment_prop"`
	Basis                        string  `json:"basis" csv:"basis"`
	//Year                         int     `json:"year" csv:"year"`
	Created int64 `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type ModelPointPricing struct {
	RunID                                int     `json:"run_id"`
	ProductCode                          string  `json:"product_code"`
	PolicyNumber                         string  `json:"policy_number"`
	PricingRunId                         int     `json:"pricing_run_id"`
	ScenarioID                           int     `json:"scenario_id" gorm:"index"`
	Age                                  int     `json:"age"`
	Gender                               string  `json:"gender"`
	SumAssured                           float64 `json:"sum_assured"`
	ScenarioDescription                  string  `json:"scenario_description"`
	CalculatedMinimumPremium             float64 `json:"calculated_minimum_premium"`
	CalculatedMinimumChildPremium        float64 `json:"calculated_minimum_child_premium"`
	CalculatedAnnualPremium              float64 `json:"calculated_annual_premium"`
	CalculatedPremiumRate                float64 `json:"calculated_premium_rate"`
	DiscountedPremiumIncome              float64 `json:"discounted_premium_income"`
	DiscountedPremiumNotReceived         float64 `json:"discounted_premium_not_received"`
	DiscountedCommission                 float64 `json:"discounted_commission"`
	DiscountedClawBack                   float64 `json:"discounted_claw_back"`
	DiscountedDeathOutgo                 float64 `json:"discounted_death_outgo"`
	DiscountedDisabilityOutgo            float64 `json:"discounted_disability_outgo"`
	DiscountedRetrenchmentOutgo          float64 `json:"discounted_retrenchment_outgo"`
	ChildDiscountedDeathOutgo            float64 `json:"child_discounted_death_outgo"`
	SpouseDiscountedDeathOutgo           float64 `json:"spouse_discounted_death_outgo"`
	DiscountedAccidentalDeathOutgo       float64 `json:"discounted_accidental_death_outgo"`
	ChildDiscountedAccidentalDeathOutgo  float64 `json:"child_discounted_accidental_death_outgo"`
	SpouseDiscountedAccidentalDeathOutgo float64 `json:"spouse_discounted_accidental_death_outgo"`
	DiscountedEducator                   float64 `json:"discounted_educator"`
	DiscountedCashBackOnSurvival         float64 `json:"discounted_cash_back_on_survival"`
	DiscountedCashBackOnDeath            float64 `json:"discounted_cash_back_on_death"`
	DiscountedRiderFuneral               float64 `json:"discounted_rider_funeral"`
	DiscountedRider                      float64 `json:"discounted_rider"`
	DiscountedExpenses                   float64 `json:"discounted_expenses"`
	DiscountedChangeInReserve            float64 `json:"discounted_change_in_reserve"`
	DiscountedInvestmentIncome           float64 `json:"discounted_investment_income"`
	DiscountedProfit                     float64 `json:"discounted_profit"`
	AnnuityFactor                        float64 `json:"annuity_factor"`
	ChildAnnuityFactor                   float64 `json:"child_annuity_factor"`
	Weighting                            float64 `json:"weighting"`
	ProfRiskAdjustment                   float64 `json:"prof_risk_adjustment"`
	ProfCSM                              float64 `json:"prof_csm"`
	ProfLossComponent                    float64 `json:"prof_loss_component"`
	Created                              int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type Profitability struct {
	RunID                    int     `json:"-"`
	ScenarioID               int     `json:"-" gor:"index"`
	ScenarioDescription      string  `json:"-"`
	WeightedDiscountedProfit float64 `json:"weighted_discounted_profit"`
	//WeightedDiscountedPremium            float64 `json:"weighted_discounted_premium"`
	WeightedDiscountedCommission         float64 `json:"weighted_discounted_commission"`
	WeightedDiscountedRisk               float64 `json:"weighted_discounted_risk"`
	WeightedDiscountedCashBackOnSurvival float64 `json:"weighted_discounted_cash_back_on_survival"`
	WeightedDiscountedCashBackOnDeath    float64 `json:"weighted_discounted_cash_back_on_death"`
	WeightedDiscountedRider              float64 `json:"weighted_discounted_rider"`
	WeightedDiscountedEducator           float64 `json:"weighted_discounted_educator"`
	WeightedDiscountedPremiumNotReceived float64 `json:"weighted_discounted_premium_not_received"`
	WeightedDiscountedExpenses           float64 `json:"weighted_discounted_expenses"`
	//WeightedDiscountedInvestmentIncome   float64 `json:"weighted_discounted_investment_income"`
	//WeightedDiscountedChangeInReserve    float64 `json:"weighted_discounted_change_in_reserve"`
	TotalValue float64 `json:"total_value"`
	Created    int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type PricingDistribution struct {
	Age          int     `json:"age"`
	Range1Male   float64 `json:"5000_male,omitempty"`
	Range1Female float64 `json:"5000_female,omitempty"`
	Range2Male   float64 `json:"10000_male,omitempty"`
	Range2Female float64 `json:"10000_female,omitempty"`
	Range3Male   float64 `json:"20000_male,omitempty"`
	Range3Female float64 `json:"20000_female,omitempty"`
	Range4Male   float64 `json:"50000_male,omitempty"`
	Range4Female float64 `json:"50000_female,omitempty"`
}

type PricingPolicyDemographic struct {
	ID               uint       `gorm:"primary_key" json:"-"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`
	DeletedAt        *time.Time `sql:"index" json:"-"`
	ProductCode      string     `json:"product_code" csv:"product_code"`
	ANB              int        `json:"anb" csv:"anb"`
	Gender           string     `json:"gender" csv:"gender"`
	SumAssured       float64    `json:"sum_assured" csv:"sum_assured"`
	Proportion       float64    `json:"proportion" csv:"proportion"`
	PricingMpVersion string     `json:"pricing_mp_version" csv:"pricing_mp_version"`
}

type PricingDistributionRange struct {
	SumAssuredValue float64
	RangeName       string
}

type CustomGormModel struct {
	ID        uint       `json:"-" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-";sql:"index"`
}

type PricingModelPointVersion struct {
	Version     string                     `json:"version"`
	ModelPoints []ProductPricingModelPoint `json:"model_points"`
	Count       int                        `json:"count"`
}
