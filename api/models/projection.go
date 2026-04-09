package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Projection struct {
	ID                                            int     `json:"-" gorm:"primary_key"`
	JobProductID                                  int     `json:"-" gorm:"index"`
	RunType                                       int     `json:"-" gorm:"index"` // 0 for valuation, 1 for pricing
	RunDate                                       string  `json:"run_date"`
	RunId                                         int     `json:"run_id" gorm:"index"`
	RunName                                       string  `json:"run_name"`
	ProductID                                     int     `json:"-"`
	ProductCode                                   string  `json:"product_code" gorm:"index"`
	PolicyNumber                                  string  `json:"policy_number" gorm:"index"`
	PolicyCount                                   float64 `json:"policy_count"`
	SpCode                                        int     `json:"sp_code"`
	IFRS17Group                                   string  `json:"ifrs_17_group"`
	ProjectionMonth                               int     `json:"projection_month" gorm:"index"`
	ProjectionYear                                int     `json:"projection_year"`
	ValuationTimeMonth                            int     `json:"valuation_time_month"`
	CalendarMonth                                 int     `json:"calendar_month"`
	ValuationTimeYear                             float64 `json:"valuation_time_year"`
	MainMemberAgeNextBirthday                     int     `json:"main_member_age_next_birthday"`
	AgeNextBirthday                               int     `json:"age_next_birthday"`
	AccidentProportion                            float64 `json:"accident_proportion"`
	InflationFactor                               float64 `json:"inflation_factor"`
	InflationFactorAdjusted                       float64 `json:"inflation_factor_adjusted"`
	AnnuityEscalationRate                         float64 `json:"annuity_escalation_rate"`
	PremiumEscalation                             float64 `json:"premium_escalation"`
	SumAssuredEscalation                          float64 `json:"sum_assured_escalation"`
	AnnuityEscalation                             float64 `json:"annuity_escalation"`
	LapseMargin                                   float64 `json:"lapse_margin"`
	PremiumWaiverOnFactor                         float64 `json:"premium_waiver_on_factor"`
	PaidUpOnFactor                                float64 `json:"paid_up_on_factor"`
	MainMemberMortalityRate                       float64 `json:"main_member_mortality_rate"`
	MainMemberMortalityRateAdjusted               float64 `json:"main_member_mortality_rate_adjusted"`
	BaseLapse                                     float64 `json:"base_lapse"`
	BaseLapseAdjusted                             float64 `json:"base_lapse_adjusted"`
	ContractingPartyAlivePortion                  float64 `json:"contracting_party_alive_portion"`
	ContractingPartyAlivePortionAdjusted          float64 `json:"contracting_party_alive_portion_adjusted"`
	ContractingPartyPolicyLapse                   float64 `json:"contracting_party_policy_lapse"`
	ContractingPartyPolicyLapseAdjusted           float64 `json:"contracting_party_policy_lapse_adjusted"`
	NonLifeMonthlyRiskRate                        float64 `json:"non_life_risk_rate"`
	NonLifeMonthlyRiskRateAdjusted                float64 `json:"non_life_risk_rate_adjusted"`
	BaseMortalityRate                             float64 `json:"base_mortality_rate"`
	BaseMortalityRateAdjusted                     float64 `json:"base_mortality_rate_adjusted"`
	BaseIndependentLapse                          float64 `json:"base_independent_lapse"`
	BaseIndependentLapseAdjusted                  float64 `json:"base_independent_lapse_adjusted"`
	BaseRetrenchmentRate                          float64 `json:"base_retrenchment_rate"`
	BaseRetrenchmentRateAdjusted                  float64 `json:"base_retrenchment_rate_adjusted"`
	BaseDisabilityIncidenceRate                   float64 `json:"base_disability_incidence_rate"`
	BaseDisabilityIncidenceRateAdjusted           float64 `json:"base_disability_incidence_rate_adjusted"`
	MainMemberMortalityRateByMonth                float64 `json:"main_member_mortality_rate_by_month"`
	MainMemberMortalityRateAdjustedByMonth        float64 `json:"main_member_mortality_rate_adjusted_by_month"`
	IndependentMortalityRateMonthly               float64 `json:"independent_mortality_rate_monthly"`
	IndependentMortalityRateAdjustedByMonth       float64 `json:"independent_mortality_rate_adjusted_by_month"`
	IndependentLapseMonthly                       float64 `json:"independent_lapse_monthly"`
	IndependentLapseMonthlyAdjusted               float64 `json:"independent_lapse_monthly_adjusted"`
	IndependentRetrenchmentMonthly                float64 `json:"independent_retrenchment_monthly"`
	IndependentRetrenchmentMonthlyAdjusted        float64 `json:"independent_retrenchment_monthly_adjusted"`
	IndependentDisabilityMonthly                  float64 `json:"independent_disability_monthly"`
	IndependentDisabilityMonthlyAdjusted          float64 `json:"independent_disability_monthly_adjusted"`
	MonthlyDependentMortality                     float64 `json:"monthly_dependent_mortality"`
	MonthlyDependentMortalityAdjusted             float64 `json:"monthly_dependent_mortality_adjusted"`
	MonthlyDependentLapse                         float64 `json:"monthly_dependent_lapse"`
	MonthlyDependentLapseAdjusted                 float64 `json:"monthly_dependent_lapse_adjusted"`
	MonthlyDependentRetrenchment                  float64 `json:"monthly_dependent_retrenchment"`
	MonthlyDependentRetrenchmentAdjusted          float64 `json:"monthly_dependent_retrenchment_adjusted"`
	MonthlyDependentDisability                    float64 `json:"monthly_dependent_disability"`
	MonthlyDependentDisabilityAdjusted            float64 `json:"monthly_dependent_disability_adjusted"`
	InitialPolicy                                 float64 `gorm:"index" json:"initial_policy"`
	InitialPolicyAdjusted                         float64 `gorm:"index" json:"initial_policy_adjusted"`
	InitialPaidUp                                 float64 `json:"initial_paid_up"`
	InitialPaidUpAdjusted                         float64 `json:"initial_paid_up_adjusted"`
	InitialPremiumWaivers                         float64 `json:"initial_premium_waivers"`
	InitialPremiumWaiversAdjusted                 float64 `json:"initial_premium_waivers_adjusted"`
	InitialTemporaryPremiumWaivers                float64 `json:"initial_temporary_premium_waivers"`
	InitialTemporaryPremiumWaiversAdjusted        float64 `json:"initial_temporary_premium_waivers_adjusted"`
	NaturalDeathsInForce                          float64 `json:"natural_deaths_in_force"`
	NaturalDeathsInForceAdjusted                  float64 `json:"natural_deaths_in_force_adjusted"`
	NaturalDeathsPaidUp                           float64 `json:"natural_deaths_paid_up"`
	NaturalDeathsPaidUpAdjusted                   float64 `json:"natural_deaths_paid_up_adjusted"`
	NaturalDeathsPremiumWaiver                    float64 `json:"natural_deaths_premium_waiver"`
	NaturalDeathsPremiumWaiverAdjusted            float64 `json:"natural_deaths_premium_waiver_adjusted"`
	NaturalDeathsTemporaryWaivers                 float64 `json:"natural_deaths_temporary_waivers"`
	NaturalDeathsTemporaryWaiversAdjusted         float64 `json:"natural_deaths_temporary_waivers_adjusted"`
	NumberOfAccidentDeaths                        float64 `json:"number_of_deaths_accident"`
	NumberOfAccidentDeathsAdjusted                float64 `json:"number_of_deaths_accident_adjusted"`
	AccidentDeathsPaidUp                          float64 `json:"accident_deaths_paid_up"`
	AccidentDeathsPaidUpAdjusted                  float64 `json:"accident_deaths_paid_up_adjusted"`
	AccidentDeathsPremiumWaiver                   float64 `json:"accident_deaths_premium_waiver"`
	AccidentDeathsPremiumWaiverAdjusted           float64 `json:"accident_deaths_premium_waiver_adjusted"`
	AccidentDeathsTemporaryPremiumWaiver          float64 `json:"accident_deaths_temporary_premium_waiver"`
	AccidentDeathsTemporaryPremiumWaiverAdjusted  float64 `json:"accident_deaths_temporary_premium_waiver_adjusted"`
	NumberOfLapses                                float64 `json:"number_of_lapses"`
	NumberOfLapsesAdjusted                        float64 `json:"number_of_lapses_adjusted"`
	NumberOfDisabilities                          float64 `json:"number_of_disabilities"`
	NumberOfDisabilitiesAdjusted                  float64 `json:"number_of_disabilities_adjusted"`
	NumberOfRetrenchments                         float64 `json:"number_of_retrenchments"`
	NumberOfRetrenchmentsAdjusted                 float64 `json:"number_of_retrenchments_adjusted"`
	NumberOfMaturities                            float64 `json:"number_of_maturities"`
	NumberOfMaturitiesAdjusted                    float64 `json:"number_of_maturities_adjusted"`
	TotalIncrementalNaturalDeaths                 float64 `json:"total_incremental_natural_deaths"`
	TotalIncrementalNaturalDeathsAdjusted         float64 `json:"total_incremental_natural_deaths_adjusted"`
	TotalIncrementalAccidentalDeaths              float64 `json:"total_incremental_accidental_deaths"`
	TotalIncrementalAccidentalDeathsAdjusted      float64 `json:"total_incremental_accidental_deaths_adjusted"`
	TotalIncrementalLapses                        float64 `json:"total_incremental_lapses"`
	TotalIncrementalLapsesAdjusted                float64 `json:"total_incremental_lapses_adjusted"`
	TotalIncrementalDisabilities                  float64 `json:"total_incremental_disabilities"`
	TotalIncrementalDisabilitiesAdjusted          float64 `json:"total_incremental_disabilities_adjusted"`
	TotalIncrementalRetrenchments                 float64 `json:"total_incremental_retrenchments"`
	TotalIncrementalRetrenchmentsAdjusted         float64 `json:"total_incremental_retrenchments_adjusted"`
	SumAssured                                    float64 `json:"sum_assured"`
	CalculatedInstalment                          float64 `json:"calculated_instalment"`
	OutstandingSumAssured                         float64 `json:"outstanding_sum_assured"`
	StandardAdditionalLumpSum                     float64 `json:"standard_additional_lump_sum"`
	RiderSumAssured                               float64 `json:"rider_sum_assured"`
	AnnuityIncome                                 float64 `json:"annuity_income"`
	Premium                                       float64 `json:"premium"`
	AllocatedPremium                              float64 `json:"allocated_premium"`
	PolicyFee                                     float64 `json:"policy_fee"`
	PremiumAdvisoryFee                            float64 `json:"premium_advisory_fee"`
	FundAdvisoryFee                               float64 `json:"fund_advisory_fee"`
	UnfundedUnitFundSom                           float64 `json:"unfunded_unit_fund_som"`
	FundInvestmentIncome                          float64 `json:"fund_investment_income"`
	ReversionaryBonus                             float64 `json:"reversionary_bonus"`
	TerminalBonus                                 float64 `json:"terminal_bonus"`
	FundAssetManagementCharge                     float64 `json:"fund_asset_management_charge"`
	FundRiskCharge                                float64 `json:"fund_risk_charge"`
	UnfundedUnitFundEom                           float64 `json:"unfunded_unit_fund_eom"`
	UnitGrowthRiskMargin                          float64 `json:"unit_growth_risk_margin"`
	BonusStabilisationAccount                     float64 `json:"bonus_stabilisation_account"`
	SumAtRisk                                     float64 `json:"sum_at_risk"`
	SurrenderPenalty                              float64 `json:"surrender_penalty"`
	SurrenderValue                                float64 `json:"surrender_value"`
	MaturityValue                                 float64 `json:"maturity_value"`
	PremiumIncome                                 float64 `json:"premium_income"`
	PremiumIncomeAdjusted                         float64 `json:"premium_income_adjusted"`
	EAllocatedPremiumIncome                       float64 `json:"e_allocated_premium_income"`
	EAllocatedPremiumIncomeAdjusted               float64 `json:"e_allocated_premium_income_adjusted"`
	EPolicyFee                                    float64 `json:"e_policy_fee"`
	EPolicyFeeAdjusted                            float64 `json:"e_policy_fee_adjusted"`
	EPremiumAdvisoryFee                           float64 `json:"e_premium_advisory_fee"`
	EPremiumAdvisoryFeeAdjusted                   float64 `json:"e_premium_advisory_fee_adjusted"`
	EFundAdvisoryFee                              float64 `json:"e_fund_advisory_fee"`
	EFundAdvisoryFeeAdjusted                      float64 `json:"e_fund_advisory_fee_adjusted"`
	EFundInvestmentIncome                         float64 `json:"e_fund_investment_income"`
	EFundInvestmentIncomeAdjusted                 float64 `json:"e_fund_investment_income_adjusted"`
	EFundAssetManagementCharge                    float64 `json:"e_fund_asset_management_charge"`
	EFundAssetManagementChargeAdjusted            float64 `json:"e_fund_asset_management_charge_adjusted"`
	EFundRiskCharge                               float64 `json:"e_fund_risk_charge"`
	EFundRiskChargeAdjusted                       float64 `json:"e_fund_risk_charge_adjusted"`
	EBsaShareholderCharge                         float64 `json:"e_bsa_shareholder_charge"`
	EBsaShareholderChargeAdjusted                 float64 `json:"e_bsa_shareholder_charge_adjusted"`
	PremiumNotReceived                            float64 `json:"premium_not_received"`
	PremiumNotReceivedAdjusted                    float64 `json:"premium_not_received_adjusted"`
	Commission                                    float64 `json:"commission"`
	CommissionAdjusted                            float64 `json:"commission_adjusted"`
	RenewalCommission                             float64 `json:"renewal_commission"`
	RenewalCommissionAdjusted                     float64 `json:"renewal_commission_adjusted"`
	ClawBack                                      float64 `json:"claw_back"`
	ClawBackAdjusted                              float64 `json:"claw_back_adjusted"`
	SurrenderOutgo                                float64 `json:"surrender_outgo"`
	SurrenderOutgoAdjusted                        float64 `json:"surrender_outgo_adjusted"`
	ESurrenderPenaltyCharge                       float64 `json:"e_surrender_penalty_charge"`
	ESurrenderPenaltyChargeAdjusted               float64 `json:"e_surrender_penalty_charge_adjusted"`
	EPremiumAdvisoryCost                          float64 `json:"e_premium_advisory_cost"`
	EPremiumAdvisoryCostAdjusted                  float64 `json:"e_premium_advisory_cost_adjusted"`
	EFundAdvisoryCost                             float64 `json:"e_fund_advisory_cost"`
	EFundAdvisoryCostAdjusted                     float64 `json:"e_fund_advisory_cost_adjusted"`
	EGuaranteeCost                                float64 `json:"e_guarantee_cost"`
	EGuaranteeCostAdjusted                        float64 `json:"e_guarantee_cost_adjusted"`
	DeathOutgo                                    float64 `json:"death_outgo"`
	DeathOutgoAdjusted                            float64 `json:"death_outgo_adjusted"`
	NonLifeClaimsOutgo                            float64 `json:"non_life_claims_outgo"`
	NonLifeClaimsOutgoAdjusted                    float64 `json:"non_life_claims_outgo_adjusted"`
	AccidentalDeathOutgo                          float64 `json:"accidental_death_outgo"`
	AccidentalDeathOutgoAdjusted                  float64 `json:"accidental_death_outgo_adjusted"`
	CashBackOnSurvival                            float64 `json:"cash_back_on_survival"`
	CashBackOnSurvivalAdjusted                    float64 `json:"cash_back_on_survival_adjusted"`
	CashBackOnDeath                               float64 `json:"cash_back_on_death"`
	CashBackOnDeathAdjusted                       float64 `json:"cash_back_on_death_adjusted"`
	DisabilityOutgo                               float64 `json:"disability_outgo"`
	DisabilityOutgoAdjusted                       float64 `json:"disability_outgo_adjusted"`
	RetrenchmentOutgo                             float64 `json:"retrenchment_outgo"`
	RetrenchmentOutgoAdjusted                     float64 `json:"retrenchment_outgo_adjusted"`
	AnnuityOutgo                                  float64 `json:"annuity_outgo"`
	AnnuityOutgoAdjusted                          float64 `json:"annuity_outgo_adjusted"`
	Rider                                         float64 `json:"rider"`
	RiderAdjusted                                 float64 `json:"rider_adjusted"`
	InitialExpenses                               float64 `json:"initial_expenses"`
	InitialExpensesAdjusted                       float64 `json:"initial_expenses_adjusted"`
	RenewalExpenses                               float64 `json:"renewal_expenses"`
	RenewalExpensesAdjusted                       float64 `json:"renewal_expenses_adjusted"`
	MaturityOutgo                                 float64 `json:"maturity_outgo"`
	MaturityOutgoAdjusted                         float64 `json:"maturity_outgo_adjusted"`
	NetCashFlow                                   float64 `gorm:"index" json:"net_cash_flow"`
	NetCashFlowAdjusted                           float64 `gorm:"index" json:"net_cash_flow_adjusted"`
	ValuationRate                                 float64 `json:"valuation_rate"`
	ValuationRateAdjusted                         float64 `json:"valuation_rate_adjusted"`
	Reserves                                      float64 `gorm:"index" json:"reserves"`
	ReservesAdjusted                              float64 `gorm:"index" json:"reserves_adjusted"`
	ChangeInReserves                              float64 `gorm:"index" json:"change_in_reserves"`
	ChangeInReservesAdjusted                      float64 `gorm:"index" json:"change_in_reserves_adjusted"`
	InvestmentIncome                              float64 `json:"investment_income"`
	InvestmentIncomeAdjusted                      float64 `json:"investment_income_adjusted"`
	ProfitAdjustment                              float64 `json:"profit_adjustment"`
	ProfitAdjustmentAdjusted                      float64 `json:"profit_adjustment_adjusted"`
	Profit                                        float64 `json:"profit"`
	ProfitAdjusted                                float64 `json:"profit_adjusted"`
	RiskDiscountRate                              float64 `json:"risk_discount_rate"`
	RiskDiscountRateAdjusted                      float64 `json:"risk_discount_rate_adjusted"`
	VIF                                           float64 `json:"vif"`
	VIFAdjusted                                   float64 `json:"vif_adjusted"`
	CorporateTax                                  float64 `json:"corporate_tax"`
	CorporateTaxAdjusted                          float64 `json:"corporate_tax_adjusted"`
	CoverageUnits                                 float64 `json:"coverage_units"`
	DiscountedPremiumIncome                       float64 `json:"discounted_premium_income"`
	DiscountedPremiumIncomeAdjusted               float64 `json:"discounted_premium_income_adjusted"`
	DiscountedPremiumNotReceived                  float64 `json:"discounted_premium_not_received"`
	DiscountedPremiumNotReceivedAdjusted          float64 `json:"discounted_premium_not_received_adjusted"`
	DiscountedUnallocatedPremiumIncome            float64 `json:"discounted_unallocated_premium_income"`
	DiscountedUnallocatedPremiumIncomeAdjusted    float64 `json:"discounted_unallocated_premium_income_adjusted"`
	DiscountedPolicyFee                           float64 `json:"discounted_policy_fee"`
	DiscountedPolicyFeeAdjusted                   float64 `json:"discounted_policy_fee_adjusted"`
	DiscountedPremiumAdvisoryFee                  float64 `json:"discounted_premium_advisory_fee"`
	DiscountedPremiumAdvisoryFeeAdjusted          float64 `json:"discounted_premium_advisory_fee_adjusted"`
	DiscountedFundAdvisoryFee                     float64 `json:"discounted_fund_advisory_fee"`
	DiscountedFundAdvisoryFeeAdjusted             float64 `json:"discounted_fund_advisory_fee_adjusted"`
	DiscountedFundAssetManagementCharges          float64 `json:"discounted_fund_asset_management_charges"`
	DiscountedFundAssetManagementChargesAdjusted  float64 `json:"discounted_fund_asset_management_charges_adjusted"`
	DiscountedFundRiskCharge                      float64 `json:"discounted_fund_risk_charge"`
	DiscountedFundRiskChargeAdjusted              float64 `json:"discounted_fund_risk_charge_adjusted"`
	DiscountedSurrenderPenaltyCharge              float64 `json:"discounted_surrender_penalty_charge"`
	DiscountedSurrenderPenaltyChargeAdjusted      float64 `json:"discounted_surrender_penalty_charge_adjusted"`
	DiscountedBsaShareholderCharge                float64 `json:"discounted_bsa_shareholder_charge"`
	DiscountedBsaShareholderChargeAdjusted        float64 `json:"discounted_bsa_shareholder_charge_adjusted"`
	DiscountedPremiumAdvisoryCost                 float64 `json:"discounted_premium_advisory_cost"`
	DiscountedPremiumAdvisoryCostAdjusted         float64 `json:"discounted_premium_advisory_cost_adjusted"`
	DiscountedFundAdvisoryCost                    float64 `json:"discounted_fund_advisory_cost"`
	DiscountedFundAdvisoryCostAdjusted            float64 `json:"discounted_fund_advisory_cost_adjusted"`
	DiscountedGuaranteeCost                       float64 `json:"discounted_guarantee_cost"`
	DiscountedGuaranteeCostAdjusted               float64 `json:"discounted_guarantee_cost_adjusted"`
	DiscountedSurrenderOutgo                      float64 `json:"discounted_surrender_outgo"`
	DiscountedSurrenderOutgoAdjusted              float64 `json:"discounted_surrender_outgo_adjusted"`
	DiscountedMaturityOutgo                       float64 `json:"discounted_maturity_outgo"`
	DiscountedMaturityOutgoAdjusted               float64 `json:"discounted_maturity_outgo_adjusted"`
	DiscountedCommission                          float64 `json:"discounted_commission"`
	DiscountedCommissionAdjusted                  float64 `json:"discounted_commission_adjusted"`
	DiscountedRenewalCommission                   float64 `json:"discounted_renewal_commission"`
	DiscountedRenewalCommissionAdjusted           float64 `json:"discounted_renewal_commission_adjusted"`
	DiscountedClawBack                            float64 `json:"discounted_claw_back"`
	DiscountedClawBackAdjusted                    float64 `json:"discounted_claw_back_adjusted"`
	DiscountedDeathOutgo                          float64 `json:"discounted_death_outgo"`
	DiscountedDeathOutgoAdjusted                  float64 `json:"discounted_death_outgo_adjusted"`
	DiscountedDisabilityOutgo                     float64 `json:"discounted_disability_outgo"`
	DiscountedDisabilityOutgoAdjusted             float64 `json:"discounted_disability_outgo_adjusted"`
	DiscountedRetrenchmentOutgo                   float64 `json:"discounted_retrenchment_outgo"`
	DiscountedRetrenchmentOutgoAdjusted           float64 `json:"discounted_retrenchment_outgo_adjusted"`
	DiscountedAnnuityOutgo                        float64 `json:"discounted_annuity_outgo"`
	DiscountedAnnuityOutgoAdjusted                float64 `json:"discounted_annuity_outgo_adjusted"`
	DiscountedAccidentalDeathOutgo                float64 `json:"discounted_accidental_death_outgo"`
	DiscountedAccidentalDeathOutgoAdjusted        float64 `json:"discounted_accidental_death_outgo_adjusted"`
	DiscountedCashBackOnSurvival                  float64 `json:"discounted_cash_back_on_survival"`
	DiscountedCashBackOnSurvivalAdjusted          float64 `json:"discounted_cash_back_on_survival_adjusted"`
	DiscountedCashBackOnDeath                     float64 `json:"discounted_cash_back_on_death"`
	DiscountedCashBackOnDeathAdjusted             float64 `json:"discounted_cash_back_on_death_adjusted"`
	DiscountedRider                               float64 `json:"discounted_rider"`
	DiscountedRiderAdjusted                       float64 `json:"discounted_rider_adjusted"`
	DiscountedInitialExpenses                     float64 `json:"discounted_initial_expenses"`
	DiscountedInitialExpensesAdjusted             float64 `json:"discounted_initial_expenses_adjusted"`
	DiscountedRenewalExpenses                     float64 `json:"discounted_renewal_expenses"`
	DiscountedRenewalExpensesAdjusted             float64 `json:"discounted_renewal_expenses_adjusted"`
	DiscountedInvestmentIncome                    float64 `json:"discounted_investment_income"`
	DiscountedInvestmentIncomeAdjusted            float64 `json:"discounted_investment_income_adjusted"`
	DiscountedProfitAdjustment                    float64 `json:"discounted_profit_adjustment"`
	DiscountedProfitAdjustmentAdjusted            float64 `json:"discounted_profit_adjustment_adjusted"`
	DiscountedProfit                              float64 `json:"discounted_profit"`
	DiscountedProfitAdjusted                      float64 `json:"discounted_profit_adjusted"`
	DiscountedCorporateTax                        float64 `json:"discounted_corporate_tax"`
	DiscountedCorporateTaxAdjusted                float64 `json:"discounted_corporate_tax_adjusted"`
	DiscountedAnnuityFactor                       float64 `json:"discounted_annuity_factor"`
	DiscountedCashOutflow                         float64 `json:"discounted_cash_outflow"`
	DiscountedCashOutflowExclAcquisition          float64 `json:"discounted_cash_outflow_excl_acquisition"`
	DiscountedAcquisitionCost                     float64 `json:"discounted_acquisition_cost"`
	DiscountedCashInflow                          float64 `json:"discounted_cash_inflow"`
	SumCoverageUnits                              float64 `json:"sum_coverage_units"`
	DiscountedCoverageUnits                       float64 `json:"discounted_coverage_units"`
	CededSumAssured                               float64 `json:"ceded_sum_assured"`
	ReinsurancePremium                            float64 `json:"reinsurance_premium"`
	ReinsurancePremiumAdjusted                    float64 `json:"reinsurance_premium_adjusted"`
	ReinsuranceCedingCommission                   float64 `json:"reinsurance_ceding_commission"`
	ReinsuranceCedingCommissionAdjusted           float64 `json:"reinsurance_ceding_commission_adjusted"`
	ReinsuranceClaims                             float64 `json:"reinsurance_claims"`
	ReinsuranceClaimsAdjusted                     float64 `json:"reinsurance_claims_adjusted"`
	DiscountedReinsurancePremium                  float64 `json:"discounted_reinsurance_premium"`
	DiscountedReinsurancePremiumAdjusted          float64 `json:"discounted_reinsurance_premium_adjusted"`
	DiscountedReinsuranceCedingCommission         float64 `json:"discounted_reinsurance_ceding_commission"`
	DiscountedReinsuranceCedingCommissionAdjusted float64 `json:"discounted_reinsurance_ceding_commission_adjusted"`
	DiscountedReinsuranceClaims                   float64 `json:"discounted_reinsurance_claims"`
	DiscountedReinsuranceClaimsAdjusted           float64 `json:"discounted_reinsurance_claims_adjusted"`
	NetReinsuranceCashflow                        float64 `json:"net_reinsurance_cashflow"`
	NetReinsuranceCashflowAdjusted                float64 `json:"net_reinsurance_cashflow_adjusted"`
	DiscountedNetReinsuranceCashflow              float64 `json:"discounted_net_reinsurance_cashflow"`
	DiscountedNetReinsuranceCashflowAdjusted      float64 `json:"discounted_net_reinsurance_cashflow_adjusted"`
	RunBasis                                      string  `json:"run_basis"`
}

type AggregatedReserves struct {
	ProjectionJobProduct JobProduct          `json:"projection_job"`
	AggregatedReserves   []AggregatedReserve `json:"aggregated_reserves"`
}
type AllAggregatedReserves struct {
	AllAggregatedReserves []AggregatedReserves `json:"all_aggregated_reserves"`
}

type AggregatedReserve struct {
	ProjectionMonth int     `json:"projection_month"`
	Reserves        float64 `json:"reserves"`
}
type AggregatedProjection struct {
	ID                                            int     `json:"-" gorm:"primary_key"`
	JobProductID                                  int     `json:"-" gorm:"index"`
	RunType                                       int     `json:"-" gorm:"index"` // 0 for valuation, 1 for pricing
	RunDate                                       string  `json:"run_date"`
	RunName                                       string  `json:"run_name"`
	RunId                                         int     `json:"run_id" gorm:"index"`
	ProductID                                     int     `json:"-"`
	ProductCode                                   string  `json:"product_code" gorm:"index"`
	PolicyNumber                                  string  `json:"-"`
	PolicyCount                                   float64 `json:"policy_count"`
	SpCode                                        int     `json:"sp_code" gorm:"index"`
	IFRS17Group                                   string  `json:"ifrs17_group"`
	ProjectionMonth                               int     `json:"projection_month" gorm:"index;column:projection_month"`
	ProjectionYear                                int     `json:"-"`
	CalendarMonth                                 int     `json:"calendar_month"`
	ValuationTimeMonth                            int     `json:"-"`
	MainMemberAgeNextBirthday                     int     `json:"-"`
	AgeNextBirthday                               int     `json:"-"`
	AccidentProportion                            float64 `json:"-"`
	InflationFactor                               float64 `json:"-"`
	InflationFactorAdjusted                       float64 `json:"-"`
	LapseMargin                                   float64 `json:"-"`
	PremiumWaiverOnFactor                         float64 `json:"-"`
	PaidUpOnFactor                                float64 `json:"-"`
	MainMemberMortalityRate                       float64 `json:"-"`
	MainMemberMortalityRateAdjusted               float64 `json:"-"`
	BaseLapse                                     float64 `json:"-"`
	BaseLapseAdjusted                             float64 `json:"-"`
	ContractingPartyAlivePortion                  float64 `json:"-"`
	ContractingPartyAlivePortionAdjusted          float64 `json:"-"`
	ContractingPartyPolicyLapse                   float64 `json:"-"`
	ContractingPartyPolicyLapseAdjusted           float64 `json:"-"`
	NonLifeMonthlyRiskRate                        float64 `json:"-"`
	NonLifeMonthlyRiskRateAdjusted                float64 `json:"-"`
	BaseMortalityRate                             float64 `json:"-"`
	BaseMortalityRateAdjusted                     float64 `json:"-"`
	BaseIndependentLapse                          float64 `json:"-"`
	BaseIndependentLapseAdjusted                  float64 `json:"-"`
	BaseRetrenchmentRate                          float64 `json:"-"`
	BaseRetrenchmentRateAdjusted                  float64 `json:"-"`
	BaseDisabilityIncidenceRate                   float64 `json:"-"`
	BaseDisabilityIncidenceRateAdjusted           float64 `json:"-"`
	MainMemberMortalityRateByMonth                float64 `json:"-"`
	MainMemberMortalityRateAdjustedByMonth        float64 `json:"-"`
	IndependentMortalityRateMonthly               float64 `json:"-"`
	IndependentMortalityRateAdjustedByMonth       float64 `json:"-"`
	IndependentLapseMonthly                       float64 `json:"-"`
	IndependentLapseMonthlyAdjusted               float64 `json:"-"`
	IndependentRetrenchmentMonthly                float64 `json:"-"`
	IndependentRetrenchmentMonthlyAdjusted        float64 `json:"-"`
	IndependentDisabilityMonthly                  float64 `json:"-"`
	IndependentDisabilityMonthlyAdjusted          float64 `json:"-"`
	MonthlyDependentMortality                     float64 `json:"-"`
	MonthlyDependentMortalityAdjusted             float64 `json:"-"`
	MonthlyDependentLapse                         float64 `json:"-"`
	MonthlyDependentLapseAdjusted                 float64 `json:"-"`
	MonthlyDependentRetrenchment                  float64 `json:"-"`
	MonthlyDependentRetrenchmentAdjusted          float64 `json:"-"`
	MonthlyDependentDisability                    float64 `json:"-"`
	MonthlyDependentDisabilityAdjusted            float64 `json:"-"`
	OutstandingSumAssured                         float64 `json:"-"`
	DiscountedBsaShareholderCharge                float64 `json:"discounted_bsa_shareholder_charge"`
	DiscountedBsaShareholderChargeAdjusted        float64 `json:"discounted_bsa_shareholder_charge_adjusted"`
	ValuationTimeYear                             float64 `json:"valuation_time_year"`
	AnnuityEscalationRate                         float64 `json:"annuity_escalation_rate"`
	PremiumEscalation                             float64 `json:"premium_escalation"`
	SumAssuredEscalation                          float64 `json:"sum_assured_escalation"`
	AnnuityEscalation                             float64 `json:"annuity_escalation"`
	InitialPolicy                                 float64 `json:"initial_policy"`
	InitialPolicyAdjusted                         float64 `json:"initial_policy_adjusted"`
	InitialPaidUp                                 float64 `json:"initial_paid_up"`
	InitialPaidUpAdjusted                         float64 `json:"initial_paid_up_adjusted"`
	InitialPremiumWaivers                         float64 `json:"initial_premium_waivers"`
	InitialPremiumWaiversAdjusted                 float64 `json:"initial_premium_waivers_adjusted"`
	InitialTemporaryPremiumWaivers                float64 `json:"initial_temporary_premium_waivers"`
	InitialTemporaryPremiumWaiversAdjusted        float64 `json:"initial_temporary_premium_waivers_adjusted"`
	NaturalDeathsInForce                          float64 `json:"natural_deaths_in_force"`
	NaturalDeathsInForceAdjusted                  float64 `json:"natural_deaths_in_force_adjusted"`
	NaturalDeathsPaidUp                           float64 `json:"natural_deaths_paid_up"`
	NaturalDeathsPaidUpAdjusted                   float64 `json:"natural_deaths_paid_up_adjusted"`
	NaturalDeathsPremiumWaiver                    float64 `json:"natural_deaths_premium_waiver"`
	NaturalDeathsPremiumWaiverAdjusted            float64 `json:"natural_deaths_premium_waiver_adjusted"`
	NaturalDeathsTemporaryWaivers                 float64 `json:"natural_deaths_temporary_waivers"`
	NaturalDeathsTemporaryWaiversAdjusted         float64 `json:"natural_deaths_temporary_waivers_adjusted"`
	NumberOfAccidentDeaths                        float64 `json:"number_of_accident_deaths"`
	NumberOfAccidentDeathsAdjusted                float64 `json:"number_of_accident_deaths_adjusted"`
	AccidentDeathsPaidUp                          float64 `json:"accident_deaths_paid_up"`
	AccidentDeathsPaidUpAdjusted                  float64 `json:"accident_deaths_paid_up_adjusted"`
	AccidentDeathsPremiumWaiver                   float64 `json:"accident_deaths_premium_waiver"`
	AccidentDeathsPremiumWaiverAdjusted           float64 `json:"accident_deaths_premium_waiver_adjusted"`
	AccidentDeathsTemporaryPremiumWaiver          float64 `json:"accident_deaths_temporary_premium_waiver"`
	AccidentDeathsTemporaryPremiumWaiverAdjusted  float64 `json:"accident_deaths_temporary_premium_waiver_adjusted"`
	NumberOfLapses                                float64 `json:"number_of_lapses"`
	NumberOfLapsesAdjusted                        float64 `json:"number_of_lapses_adjusted"`
	NumberOfDisabilities                          float64 `json:"number_of_disabilities"`
	NumberOfDisabilitiesAdjusted                  float64 `json:"number_of_disabilities_adjusted"`
	NumberOfRetrenchments                         float64 `json:"number_of_retrenchments"`
	NumberOfRetrenchmentsAdjusted                 float64 `json:"number_of_retrenchments_adjusted"`
	NumberOfMaturities                            float64 `json:"number_of_maturities"`
	NumberOfMaturitiesAdjusted                    float64 `json:"number_of_maturities_adjusted"`
	TotalIncrementalNaturalDeaths                 float64 `json:"total_incremental_natural_deaths"`
	TotalIncrementalNaturalDeathsAdjusted         float64 `json:"total_incremental_natural_deaths_adjusted"`
	TotalIncrementalAccidentalDeaths              float64 `json:"total_incremental_accidental_deaths"`
	TotalIncrementalAccidentalDeathsAdjusted      float64 `json:"total_incremental_accidental_deaths_adjusted"`
	TotalIncrementalLapses                        float64 `json:"total_incremental_lapses"`
	TotalIncrementalLapsesAdjusted                float64 `json:"total_incremental_lapses_adjusted"`
	TotalIncrementalDisabilities                  float64 `json:"total_incremental_disabilities"`
	TotalIncrementalDisabilitiesAdjusted          float64 `json:"total_incremental_disabilities_adjusted"`
	TotalIncrementalRetrenchments                 float64 `json:"total_incremental_retrenchments"`
	TotalIncrementalRetrenchmentsAdjusted         float64 `json:"total_incremental_retrenchments_adjusted"`
	SumAssured                                    float64 `json:"sum_assured"`
	CalculatedInstalment                          float64 `json:"calculated_instalment"`
	StandardAdditionalLumpSum                     float64 `json:"standard_additional_lump_sum"`
	RiderSumAssured                               float64 `json:"rider_sum_assured"`
	AnnuityIncome                                 float64 `json:"annuity_income"`
	Premium                                       float64 `json:"premium"`
	AllocatedPremium                              float64 `json:"allocated_premium"`
	PolicyFee                                     float64 `json:"policy_fee"`
	PremiumAdvisoryFee                            float64 `json:"premium_advisory_fee"`
	FundAdvisoryFee                               float64 `json:"fund_advisory_fee"`
	UnfundedUnitFundSom                           float64 `json:"unfunded_unit_fund_som"`
	FundInvestmentIncome                          float64 `json:"fund_investment_income"`
	ReversionaryBonus                             float64 `json:"reversionary_bonus"`
	TerminalBonus                                 float64 `json:"terminal_bonus"`
	FundAssetManagementCharge                     float64 `json:"fund_asset_management_charge"`
	FundRiskCharge                                float64 `json:"fund_risk_charge"`
	UnfundedUnitFundEom                           float64 `json:"unfunded_unit_fund_eom"`
	UnitGrowthRiskMargin                          float64 `json:"unit_growth_risk_margin"`
	BonusStabilisationAccount                     float64 `json:"bonus_stabilisation_account"`
	SumAtRisk                                     float64 `json:"sum_at_risk"`
	SurrenderPenalty                              float64 `json:"surrender_penalty"`
	SurrenderValue                                float64 `json:"surrender_value"`
	MaturityValue                                 float64 `json:"maturity_value"`
	PremiumIncome                                 float64 `json:"premium_income"`
	PremiumIncomeAdjusted                         float64 `json:"premium_income_adjusted"`
	EAllocatedPremiumIncome                       float64 `json:"e_allocated_premium_income"`
	EAllocatedPremiumIncomeAdjusted               float64 `json:"e_allocated_premium_income_adjusted"`
	EPolicyFee                                    float64 `json:"e_policy_fee"`
	EPolicyFeeAdjusted                            float64 `json:"e_policy_fee_adjusted"`
	EPremiumAdvisoryFee                           float64 `json:"e_premium_advisory_fee"`
	EPremiumAdvisoryFeeAdjusted                   float64 `json:"e_premium_advisory_fee_adjusted"`
	EFundAdvisoryFee                              float64 `json:"e_fund_advisory_fee"`
	EFundAdvisoryFeeAdjusted                      float64 `json:"e_fund_advisory_fee_adjusted"`
	EFundInvestmentIncome                         float64 `json:"e_fund_investment_income"`
	EFundInvestmentIncomeAdjusted                 float64 `json:"e_fund_investment_income_adjusted"`
	EFundAssetManagementCharge                    float64 `json:"e_fund_asset_management_charge"`
	EFundAssetManagementChargeAdjusted            float64 `json:"e_fund_asset_management_charge_adjusted"`
	EFundRiskCharge                               float64 `json:"e_fund_risk_charge"`
	EFundRiskChargeAdjusted                       float64 `json:"e_fund_risk_charge_adjusted"`
	EBsaShareholderCharge                         float64 `json:"e_bsa_shareholder_charge"`
	EBsaShareholderChargeAdjusted                 float64 `json:"e_bsa_shareholder_charge_adjusted"`
	PremiumNotReceived                            float64 `json:"premium_not_received"`
	PremiumNotReceivedAdjusted                    float64 `json:"premium_not_received_adjusted"`
	Commission                                    float64 `json:"commission"`
	CommissionAdjusted                            float64 `json:"commission_adjusted"`
	RenewalCommission                             float64 `json:"renewal_commission"`
	RenewalCommissionAdjusted                     float64 `json:"renewal_commission_adjusted"`
	ClawBack                                      float64 `json:"claw_back"`
	ClawBackAdjusted                              float64 `json:"claw_back_adjusted"`
	SurrenderOutgo                                float64 `json:"surrender_outgo"`
	SurrenderOutgoAdjusted                        float64 `json:"surrender_outgo_adjusted"`
	ESurrenderPenaltyCharge                       float64 `json:"e_surrender_penalty_charge"`
	ESurrenderPenaltyChargeAdjusted               float64 `json:"e_surrender_penalty_charge_adjusted"`
	EPremiumAdvisoryCost                          float64 `json:"e_premium_advisory_cost"`
	EPremiumAdvisoryCostAdjusted                  float64 `json:"e_premium_advisory_cost_adjusted"`
	EFundAdvisoryCost                             float64 `json:"e_fund_advisory_cost"`
	EFundAdvisoryCostAdjusted                     float64 `json:"e_fund_advisory_cost_adjusted"`
	EGuaranteeCost                                float64 `json:"e_guarantee_cost"`
	EGuaranteeCostAdjusted                        float64 `json:"e_guarantee_cost_adjusted"`
	NonLifeClaimsOutgo                            float64 `json:"non_life_claims_outgo"`
	NonLifeClaimsOutgoAdjusted                    float64 `json:"non_life_claims_outgo_adjusted"`
	DeathOutgo                                    float64 `json:"death_outgo"`
	DeathOutgoAdjusted                            float64 `json:"death_outgo_adjusted"`
	AccidentalDeathOutgo                          float64 `json:"accidental_death_outgo"`
	AccidentalDeathOutgoAdjusted                  float64 `json:"accidental_death_outgo_adjusted"`
	CashBackOnSurvival                            float64 `json:"cash_back_on_survival"`
	CashBackOnSurvivalAdjusted                    float64 `json:"cash_back_on_survival_adjusted"`
	CashBackOnDeath                               float64 `json:"cash_back_on_death"`
	CashBackOnDeathAdjusted                       float64 `json:"cash_back_on_death_adjusted"`
	DisabilityOutgo                               float64 `json:"disability_outgo"`
	DisabilityOutgoAdjusted                       float64 `json:"disability_outgo_adjusted"`
	DiscountedDisabilityOutgo                     float64 `json:"discounted_disability_outgo"`
	DiscountedDisabilityOutgoAdjusted             float64 `json:"discounted_disability_outgo_adjusted"`
	DiscountedRetrenchmentOutgo                   float64 `json:"discounted_retrenchment_outgo"`
	DiscountedRetrenchmentOutgoAdjusted           float64 `json:"discounted_retrenchment_outgo_adjusted"`
	RetrenchmentOutgo                             float64 `json:"retrenchment_outgo"`
	RetrenchmentOutgoAdjusted                     float64 `json:"retrenchment_outgo_adjusted"`
	AnnuityOutgo                                  float64 `json:"annuity_outgo"`
	AnnuityOutgoAdjusted                          float64 `json:"annuity_outgo_adjusted"`
	Rider                                         float64 `json:"rider"`
	RiderAdjusted                                 float64 `json:"rider_adjusted"`
	InitialExpenses                               float64 `json:"initial_expenses"`
	InitialExpensesAdjusted                       float64 `json:"initial_expenses_adjusted"`
	RenewalExpenses                               float64 `json:"renewal_expenses"`
	RenewalExpensesAdjusted                       float64 `json:"renewal_expenses_adjusted"`
	MaturityOutgo                                 float64 `json:"maturity_outgo"`
	MaturityOutgoAdjusted                         float64 `json:"maturity_outgo_adjusted"`
	NetCashFlow                                   float64 `json:"net_cash_flow"`
	NetCashFlowAdjusted                           float64 `json:"net_cash_flow_adjusted"`
	ValuationRate                                 float64 `json:"valuation_rate"`
	ValuationRateAdjusted                         float64 `json:"valuation_rate_adjusted"`
	Reserves                                      float64 `json:"reserves"`
	ReservesAdjusted                              float64 `json:"reserves_adjusted"`
	ChangeInReserves                              float64 `json:"change_in_reserves"`
	ChangeInReservesAdjusted                      float64 `json:"change_in_reserves_adjusted"`
	InvestmentIncome                              float64 `json:"investment_income"`
	InvestmentIncomeAdjusted                      float64 `json:"investment_income_adjusted"`
	ProfitAdjustment                              float64 `json:"profit_adjustment"`
	ProfitAdjustmentAdjusted                      float64 `json:"profit_adjustment_adjusted"`
	Profit                                        float64 `json:"profit"`
	ProfitAdjusted                                float64 `json:"profit_adjusted"`
	RiskDiscountRate                              float64 `json:"risk_discount_rate"`
	RiskDiscountRateAdjusted                      float64 `json:"risk_discount_rate_adjusted"`
	VIF                                           float64 `json:"vif"`
	VIFAdjusted                                   float64 `json:"vif_adjusted"`
	CorporateTax                                  float64 `json:"corporate_tax"`
	CorporateTaxAdjusted                          float64 `json:"corporate_tax_adjusted"`
	CoverageUnits                                 float64 `json:"coverage_units"`
	DiscountedPremiumIncome                       float64 `json:"discounted_premium_income"`
	DiscountedPremiumIncomeAdjusted               float64 `json:"discounted_premium_income_adjusted"`
	DiscountedPremiumNotReceived                  float64 `json:"discounted_premium_not_received"`
	DiscountedPremiumNotReceivedAdjusted          float64 `json:"discounted_premium_not_received_adjusted"`
	DiscountedUnallocatedPremiumIncome            float64 `json:"discounted_unallocated_premium_income"`
	DiscountedUnallocatedPremiumIncomeAdjusted    float64 `json:"discounted_unallocated_premium_income_adjusted"`
	DiscountedPolicyFee                           float64 `json:"discounted_policy_fee"`
	DiscountedPolicyFeeAdjusted                   float64 `json:"discounted_policy_fee_adjusted"`
	DiscountedPremiumAdvisoryFee                  float64 `json:"discounted_premium_advisory_fee"`
	DiscountedPremiumAdvisoryFeeAdjusted          float64 `json:"discounted_premium_advisory_fee_adjusted"`
	DiscountedFundAdvisoryFee                     float64 `json:"discounted_fund_advisory_fee"`
	DiscountedFundAdvisoryFeeAdjusted             float64 `json:"discounted_fund_advisory_fee_adjusted"`
	DiscountedFundAssetManagementCharges          float64 `json:"discounted_fund_asset_management_charges"`
	DiscountedFundAssetManagementChargesAdjusted  float64 `json:"discounted_fund_asset_management_charges_adjusted"`
	DiscountedFundRiskCharge                      float64 `json:"discounted_fund_risk_charge"`
	DiscountedFundRiskChargeAdjusted              float64 `json:"discounted_fund_risk_charge_adjusted"`
	DiscountedSurrenderPenaltyCharge              float64 `json:"discounted_surrender_penalty_charge"`
	DiscountedSurrenderPenaltyChargeAdjusted      float64 `json:"discounted_surrender_penalty_charge_adjusted"`
	DiscountedPremiumAdvisoryCost                 float64 `json:"discounted_premium_advisory_cost"`
	DiscountedPremiumAdvisoryCostAdjusted         float64 `json:"discounted_premium_advisory_cost_adjusted"`
	DiscountedFundAdvisoryCost                    float64 `json:"discounted_fund_advisory_cost"`
	DiscountedFundAdvisoryCostAdjusted            float64 `json:"discounted_fund_advisory_cost_adjusted"`
	DiscountedGuaranteeCost                       float64 `json:"discounted_guarantee_cost"`
	DiscountedGuaranteeCostAdjusted               float64 `json:"discounted_guarantee_cost_adjusted"`
	DiscountedSurrenderOutgo                      float64 `json:"discounted_surrender_outgo"`
	DiscountedSurrenderOutgoAdjusted              float64 `json:"discounted_surrender_outgo_adjusted"`
	DiscountedMaturityOutgo                       float64 `json:"discounted_maturity_outgo"`
	DiscountedMaturityOutgoAdjusted               float64 `json:"discounted_maturity_outgo_adjusted"`
	DiscountedCommission                          float64 `json:"discounted_commission"`
	DiscountedCommissionAdjusted                  float64 `json:"discounted_commission_adjusted"`
	DiscountedRenewalCommission                   float64 `json:"discounted_renewal_commission"`
	DiscountedRenewalCommissionAdjusted           float64 `json:"discounted_renewal_commission_adjusted"`
	DiscountedClawBack                            float64 `json:"discounted_claw_back"`
	DiscountedClawBackAdjusted                    float64 `json:"discounted_claw_back_adjusted"`
	DiscountedDeathOutgo                          float64 `json:"discounted_death_outgo"`
	DiscountedDeathOutgoAdjusted                  float64 `json:"discounted_death_outgo_adjusted"`
	DiscountedAccidentalDeathOutgo                float64 `json:"discounted_accidental_death_outgo"`
	DiscountedAccidentalDeathOutgoAdjusted        float64 `json:"discounted_accidental_death_outgo_adjusted"`
	DiscountedCashBackOnSurvival                  float64 `json:"discounted_cash_back_on_survival"`
	DiscountedCashBackOnSurvivalAdjusted          float64 `json:"discounted_cash_back_on_survival_adjusted"`
	DiscountedCashBackOnDeath                     float64 `json:"discounted_cash_back_on_death"`
	DiscountedCashBackOnDeathAdjusted             float64 `json:"discounted_cash_back_on_death_adjusted"`
	DiscountedRider                               float64 `json:"discounted_rider"`
	DiscountedRiderAdjusted                       float64 `json:"discounted_rider_adjusted"`
	DiscountedInitialExpenses                     float64 `json:"discounted_initial_expenses"`
	DiscountedInitialExpensesAdjusted             float64 `json:"discounted_initial_expenses_adjusted"`
	DiscountedRenewalExpenses                     float64 `json:"discounted_renewal_expenses"`
	DiscountedRenewalExpensesAdjusted             float64 `json:"discounted_renewal_expenses_adjusted"`
	DiscountedAnnuityOutgo                        float64 `json:"discounted_annuity_outgo"`
	DiscountedAnnuityOutgoAdjusted                float64 `json:"discounted_annuity_outgo_adjusted"`
	DiscountedInvestmentIncome                    float64 `json:"discounted_investment_income"`
	DiscountedInvestmentIncomeAdjusted            float64 `json:"discounted_investment_income_adjusted"`
	DiscountedProfitAdjustment                    float64 `json:"discounted_profit_adjustment"`
	DiscountedProfitAdjustmentAdjusted            float64 `json:"discounted_profit_adjustment_adjusted"`
	DiscountedProfit                              float64 `json:"discounted_profit"`
	DiscountedProfitAdjusted                      float64 `json:"discounted_profit_adjusted"`
	DiscountedCorporateTax                        float64 `json:"discounted_corporate_tax"`
	DiscountedCorporateTaxAdjusted                float64 `json:"discounted_corporate_tax_adjusted"`
	DiscountedAnnuityFactor                       float64 `json:"discounted_annuity_factor"`
	DiscountedCashOutflow                         float64 `json:"discounted_cash_outflow"`
	DiscountedCashOutflowExclAcquisition          float64 `json:"discounted_cash_outflow_excl_acquisition"`
	DiscountedAcquisitionCost                     float64 `json:"discounted_acquisition_cost"`
	DiscountedCashInflow                          float64 `json:"discounted_cash_inflow"`
	SumCoverageUnits                              float64 `json:"sum_coverage_units"`
	DiscountedCoverageUnits                       float64 `json:"discounted_coverage_units"`
	CededSumAssured                               float64 `json:"ceded_sum_assured"`
	ReinsurancePremium                            float64 `json:"reinsurance_premium"`
	ReinsurancePremiumAdjusted                    float64 `json:"reinsurance_premium_adjusted"`
	ReinsuranceCedingCommission                   float64 `json:"reinsurance_ceding_commission"`
	ReinsuranceCedingCommissionAdjusted           float64 `json:"reinsurance_ceding_commission_adjusted"`
	ReinsuranceClaims                             float64 `json:"reinsurance_claims"`
	ReinsuranceClaimsAdjusted                     float64 `json:"reinsurance_claims_adjusted"`
	NetReinsuranceCashflow                        float64 `json:"net_reinsurance_cashflow"`
	NetReinsuranceCashflowAdjusted                float64 `json:"net_reinsurance_cashflow_adjusted"`
	DiscountedReinsurancePremium                  float64 `json:"discounted_reinsurance_premium"`
	DiscountedReinsurancePremiumAdjusted          float64 `json:"discounted_reinsurance_premium_adjusted"`
	DiscountedReinsuranceCedingCommission         float64 `json:"discounted_reinsurance_ceding_commission"`
	DiscountedReinsuranceCedingCommissionAdjusted float64 `json:"discounted_reinsurance_ceding_commission_adjusted"`
	DiscountedReinsuranceClaims                   float64 `json:"discounted_reinsurance_claims"`
	DiscountedReinsuranceClaimsAdjusted           float64 `json:"discounted_reinsurance_claims_adjusted"`
	DiscountedNetReinsuranceCashflow              float64 `json:"discounted_net_reinsurance_cashflow"`
	DiscountedNetReinsuranceCashflowAdjusted      float64 `json:"discounted_net_reinsurance_cashflow_adjusted"`
	RunBasis                                      string  `json:"run_basis"`
}

type AggregatedProjectionData struct {
	RunDate                                       string  `json:"run_date"`
	RunName                                       string  `json:"run_name"`
	RunId                                         int     `json:"run_id"`
	ProductCode                                   string  `json:"product_code" gorm:"index"`
	PolicyCount                                   float64 `json:"policy_count"`
	SpCode                                        int     `json:"sp_code" gorm:"index"`
	IFRS17Group                                   string  `json:"ifrs17_group"`
	ProjectionMonth                               int     `json:"projection_month" gorm:"column:projection_month"`
	CalendarMonth                                 int     `json:"calendar_month"`
	DiscountedBsaShareholderCharge                float64 `json:"discounted_bsa_shareholder_charge"`
	DiscountedBsaShareholderChargeAdjusted        float64 `json:"discounted_bsa_shareholder_charge_adjusted"`
	ValuationTimeYear                             float64 `json:"valuation_time_year"`
	AnnuityEscalationRate                         float64 `json:"annuity_escalation_rate"`
	PremiumEscalation                             float64 `json:"premium_escalation"`
	SumAssuredEscalation                          float64 `json:"sum_assured_escalation"`
	AnnuityEscalation                             float64 `json:"annuity_escalation"`
	InitialPolicy                                 float64 `json:"initial_policy"`
	InitialPolicyAdjusted                         float64 `json:"initial_policy_adjusted"`
	InitialPaidUp                                 float64 `json:"initial_paid_up"`
	InitialPaidUpAdjusted                         float64 `json:"initial_paid_up_adjusted"`
	InitialPremiumWaivers                         float64 `json:"initial_premium_waivers"`
	InitialPremiumWaiversAdjusted                 float64 `json:"initial_premium_waivers_adjusted"`
	InitialTemporaryPremiumWaivers                float64 `json:"initial_temporary_premium_waivers"`
	InitialTemporaryPremiumWaiversAdjusted        float64 `json:"initial_temporary_premium_waivers_adjusted"`
	NaturalDeathsInForce                          float64 `json:"natural_deaths_in_force"`
	NaturalDeathsInForceAdjusted                  float64 `json:"natural_deaths_in_force_adjusted"`
	NaturalDeathsPaidUp                           float64 `json:"natural_deaths_paid_up"`
	NaturalDeathsPaidUpAdjusted                   float64 `json:"natural_deaths_paid_up_adjusted"`
	NaturalDeathsPremiumWaiver                    float64 `json:"natural_deaths_premium_waiver"`
	NaturalDeathsPremiumWaiverAdjusted            float64 `json:"natural_deaths_premium_waiver_adjusted"`
	NaturalDeathsTemporaryWaivers                 float64 `json:"natural_deaths_temporary_waivers"`
	NaturalDeathsTemporaryWaiversAdjusted         float64 `json:"natural_deaths_temporary_waivers_adjusted"`
	NumberOfAccidentDeaths                        float64 `json:"number_of_accident_deaths"`
	NumberOfAccidentDeathsAdjusted                float64 `json:"number_of_accident_deaths_adjusted"`
	AccidentDeathsPaidUp                          float64 `json:"accident_deaths_paid_up"`
	AccidentDeathsPaidUpAdjusted                  float64 `json:"accident_deaths_paid_up_adjusted"`
	AccidentDeathsPremiumWaiver                   float64 `json:"accident_deaths_premium_waiver"`
	AccidentDeathsPremiumWaiverAdjusted           float64 `json:"accident_deaths_premium_waiver_adjusted"`
	AccidentDeathsTemporaryPremiumWaiver          float64 `json:"accident_deaths_temporary_premium_waiver"`
	AccidentDeathsTemporaryPremiumWaiverAdjusted  float64 `json:"accident_deaths_temporary_premium_waiver_adjusted"`
	NumberOfLapses                                float64 `json:"number_of_lapses"`
	NumberOfLapsesAdjusted                        float64 `json:"number_of_lapses_adjusted"`
	NumberOfDisabilities                          float64 `json:"number_of_disabilities"`
	NumberOfDisabilitiesAdjusted                  float64 `json:"number_of_disabilities_adjusted"`
	NumberOfRetrenchments                         float64 `json:"number_of_retrenchments"`
	NumberOfRetrenchmentsAdjusted                 float64 `json:"number_of_retrenchments_adjusted"`
	NumberOfMaturities                            float64 `json:"number_of_maturities"`
	NumberOfMaturitiesAdjusted                    float64 `json:"number_of_maturities_adjusted"`
	TotalIncrementalNaturalDeaths                 float64 `json:"total_incremental_natural_deaths"`
	TotalIncrementalNaturalDeathsAdjusted         float64 `json:"total_incremental_natural_deaths_adjusted"`
	TotalIncrementalAccidentalDeaths              float64 `json:"total_incremental_accidental_deaths"`
	TotalIncrementalAccidentalDeathsAdjusted      float64 `json:"total_incremental_accidental_deaths_adjusted"`
	TotalIncrementalLapses                        float64 `json:"total_incremental_lapses"`
	TotalIncrementalLapsesAdjusted                float64 `json:"total_incremental_lapses_adjusted"`
	TotalIncrementalDisabilities                  float64 `json:"total_incremental_disabilities"`
	TotalIncrementalDisabilitiesAdjusted          float64 `json:"total_incremental_disabilities_adjusted"`
	TotalIncrementalRetrenchments                 float64 `json:"total_incremental_retrenchments"`
	TotalIncrementalRetrenchmentsAdjusted         float64 `json:"total_incremental_retrenchments_adjusted"`
	SumAssured                                    float64 `json:"sum_assured"`
	CalculatedInstalment                          float64 `json:"calculated_instalment"`
	StandardAdditionalLumpSum                     float64 `json:"standard_additional_lump_sum"`
	RiderSumAssured                               float64 `json:"rider_sum_assured"`
	AnnuityIncome                                 float64 `json:"annuity_income"`
	Premium                                       float64 `json:"premium"`
	AllocatedPremium                              float64 `json:"allocated_premium"`
	PolicyFee                                     float64 `json:"policy_fee"`
	PremiumAdvisoryFee                            float64 `json:"premium_advisory_fee"`
	FundAdvisoryFee                               float64 `json:"fund_advisory_fee"`
	UnfundedUnitFundSom                           float64 `json:"unfunded_unit_fund_som"`
	FundInvestmentIncome                          float64 `json:"fund_investment_income"`
	ReversionaryBonus                             float64 `json:"reversionary_bonus"`
	TerminalBonus                                 float64 `json:"terminal_bonus"`
	FundAssetManagementCharge                     float64 `json:"fund_asset_management_charge"`
	FundRiskCharge                                float64 `json:"fund_risk_charge"`
	UnfundedUnitFundEom                           float64 `json:"unfunded_unit_fund_eom"`
	UnitGrowthRiskMargin                          float64 `json:"unit_growth_risk_margin"`
	BonusStabilisationAccount                     float64 `json:"bonus_stabilisation_account"`
	SumAtRisk                                     float64 `json:"sum_at_risk"`
	SurrenderPenalty                              float64 `json:"surrender_penalty"`
	SurrenderValue                                float64 `json:"surrender_value"`
	MaturityValue                                 float64 `json:"maturity_value"`
	PremiumIncome                                 float64 `json:"premium_income"`
	PremiumIncomeAdjusted                         float64 `json:"premium_income_adjusted"`
	EAllocatedPremiumIncome                       float64 `json:"e_allocated_premium_income"`
	EAllocatedPremiumIncomeAdjusted               float64 `json:"e_allocated_premium_income_adjusted"`
	EPolicyFee                                    float64 `json:"e_policy_fee"`
	EPolicyFeeAdjusted                            float64 `json:"e_policy_fee_adjusted"`
	EPremiumAdvisoryFee                           float64 `json:"e_premium_advisory_fee"`
	EPremiumAdvisoryFeeAdjusted                   float64 `json:"e_premium_advisory_fee_adjusted"`
	EFundAdvisoryFee                              float64 `json:"e_fund_advisory_fee"`
	EFundAdvisoryFeeAdjusted                      float64 `json:"e_fund_advisory_fee_adjusted"`
	EFundInvestmentIncome                         float64 `json:"e_fund_investment_income"`
	EFundInvestmentIncomeAdjusted                 float64 `json:"e_fund_investment_income_adjusted"`
	EFundAssetManagementCharge                    float64 `json:"e_fund_asset_management_charge"`
	EFundAssetManagementChargeAdjusted            float64 `json:"e_fund_asset_management_charge_adjusted"`
	EFundRiskCharge                               float64 `json:"e_fund_risk_charge"`
	EFundRiskChargeAdjusted                       float64 `json:"e_fund_risk_charge_adjusted"`
	EBsaShareholderCharge                         float64 `json:"e_bsa_shareholder_charge"`
	EBsaShareholderChargeAdjusted                 float64 `json:"e_bsa_shareholder_charge_adjusted"`
	PremiumNotReceived                            float64 `json:"premium_not_received"`
	PremiumNotReceivedAdjusted                    float64 `json:"premium_not_received_adjusted"`
	Commission                                    float64 `json:"commission"`
	CommissionAdjusted                            float64 `json:"commission_adjusted"`
	RenewalCommission                             float64 `json:"renewal_commission"`
	RenewalCommissionAdjusted                     float64 `json:"renewal_commission_adjusted"`
	ClawBack                                      float64 `json:"claw_back"`
	ClawBackAdjusted                              float64 `json:"claw_back_adjusted"`
	SurrenderOutgo                                float64 `json:"surrender_outgo"`
	SurrenderOutgoAdjusted                        float64 `json:"surrender_outgo_adjusted"`
	ESurrenderPenaltyCharge                       float64 `json:"e_surrender_penalty_charge"`
	ESurrenderPenaltyChargeAdjusted               float64 `json:"e_surrender_penalty_charge_adjusted"`
	EPremiumAdvisoryCost                          float64 `json:"e_premium_advisory_cost"`
	EPremiumAdvisoryCostAdjusted                  float64 `json:"e_premium_advisory_cost_adjusted"`
	EFundAdvisoryCost                             float64 `json:"e_fund_advisory_cost"`
	EFundAdvisoryCostAdjusted                     float64 `json:"e_fund_advisory_cost_adjusted"`
	EGuaranteeCost                                float64 `json:"e_guarantee_cost"`
	EGuaranteeCostAdjusted                        float64 `json:"e_guarantee_cost_adjusted"`
	NonLifeClaimsOutgo                            float64 `json:"non_life_claims_outgo"`
	NonLifeClaimsOutgoAdjusted                    float64 `json:"non_life_claims_outgo_adjusted"`
	DeathOutgo                                    float64 `json:"death_outgo"`
	DeathOutgoAdjusted                            float64 `json:"death_outgo_adjusted"`
	AccidentalDeathOutgo                          float64 `json:"accidental_death_outgo"`
	AccidentalDeathOutgoAdjusted                  float64 `json:"accidental_death_outgo_adjusted"`
	CashBackOnSurvival                            float64 `json:"cash_back_on_survival"`
	CashBackOnSurvivalAdjusted                    float64 `json:"cash_back_on_survival_adjusted"`
	CashBackOnDeath                               float64 `json:"cash_back_on_death"`
	CashBackOnDeathAdjusted                       float64 `json:"cash_back_on_death_adjusted"`
	DisabilityOutgo                               float64 `json:"disability_outgo"`
	DisabilityOutgoAdjusted                       float64 `json:"disability_outgo_adjusted"`
	DiscountedDisabilityOutgo                     float64 `json:"discounted_disability_outgo"`
	DiscountedDisabilityOutgoAdjusted             float64 `json:"discounted_disability_outgo_adjusted"`
	DiscountedRetrenchmentOutgo                   float64 `json:"discounted_retrenchment_outgo"`
	DiscountedRetrenchmentOutgoAdjusted           float64 `json:"discounted_retrenchment_outgo_adjusted"`
	RetrenchmentOutgo                             float64 `json:"retrenchment_outgo"`
	RetrenchmentOutgoAdjusted                     float64 `json:"retrenchment_outgo_adjusted"`
	AnnuityOutgo                                  float64 `json:"annuity_outgo"`
	AnnuityOutgoAdjusted                          float64 `json:"annuity_outgo_adjusted"`
	Rider                                         float64 `json:"rider"`
	RiderAdjusted                                 float64 `json:"rider_adjusted"`
	InitialExpenses                               float64 `json:"initial_expenses"`
	InitialExpensesAdjusted                       float64 `json:"initial_expenses_adjusted"`
	RenewalExpenses                               float64 `json:"renewal_expenses"`
	RenewalExpensesAdjusted                       float64 `json:"renewal_expenses_adjusted"`
	MaturityOutgo                                 float64 `json:"maturity_outgo"`
	MaturityOutgoAdjusted                         float64 `json:"maturity_outgo_adjusted"`
	NetCashFlow                                   float64 `json:"net_cash_flow"`
	NetCashFlowAdjusted                           float64 `json:"net_cash_flow_adjusted"`
	ValuationRate                                 float64 `json:"valuation_rate"`
	ValuationRateAdjusted                         float64 `json:"valuation_rate_adjusted"`
	Reserves                                      float64 `json:"reserves"`
	ReservesAdjusted                              float64 `json:"reserves_adjusted"`
	ChangeInReserves                              float64 `json:"change_in_reserves"`
	ChangeInReservesAdjusted                      float64 `json:"change_in_reserves_adjusted"`
	InvestmentIncome                              float64 `json:"investment_income"`
	InvestmentIncomeAdjusted                      float64 `json:"investment_income_adjusted"`
	ProfitAdjustment                              float64 `json:"profit_adjustment"`
	ProfitAdjustmentAdjusted                      float64 `json:"profit_adjustment_adjusted"`
	Profit                                        float64 `json:"profit"`
	ProfitAdjusted                                float64 `json:"profit_adjusted"`
	RiskDiscountRate                              float64 `json:"risk_discount_rate"`
	RiskDiscountRateAdjusted                      float64 `json:"risk_discount_rate_adjusted"`
	VIF                                           float64 `json:"vif"`
	VIFAdjusted                                   float64 `json:"vif_adjusted"`
	CorporateTax                                  float64 `json:"corporate_tax"`
	CorporateTaxAdjusted                          float64 `json:"corporate_tax_adjusted"`
	CoverageUnits                                 float64 `json:"coverage_units"`
	DiscountedPremiumIncome                       float64 `json:"discounted_premium_income"`
	DiscountedPremiumIncomeAdjusted               float64 `json:"discounted_premium_income_adjusted"`
	DiscountedPremiumNotReceived                  float64 `json:"discounted_premium_not_received"`
	DiscountedPremiumNotReceivedAdjusted          float64 `json:"discounted_premium_not_received_adjusted"`
	DiscountedUnallocatedPremiumIncome            float64 `json:"discounted_unallocated_premium_income"`
	DiscountedUnallocatedPremiumIncomeAdjusted    float64 `json:"discounted_unallocated_premium_income_adjusted"`
	DiscountedPolicyFee                           float64 `json:"discounted_policy_fee"`
	DiscountedPolicyFeeAdjusted                   float64 `json:"discounted_policy_fee_adjusted"`
	DiscountedPremiumAdvisoryFee                  float64 `json:"discounted_premium_advisory_fee"`
	DiscountedPremiumAdvisoryFeeAdjusted          float64 `json:"discounted_premium_advisory_fee_adjusted"`
	DiscountedFundAdvisoryFee                     float64 `json:"discounted_fund_advisory_fee"`
	DiscountedFundAdvisoryFeeAdjusted             float64 `json:"discounted_fund_advisory_fee_adjusted"`
	DiscountedFundAssetManagementCharges          float64 `json:"discounted_fund_asset_management_charges"`
	DiscountedFundAssetManagementChargesAdjusted  float64 `json:"discounted_fund_asset_management_charges_adjusted"`
	DiscountedFundRiskCharge                      float64 `json:"discounted_fund_risk_charge"`
	DiscountedFundRiskChargeAdjusted              float64 `json:"discounted_fund_risk_charge_adjusted"`
	DiscountedSurrenderPenaltyCharge              float64 `json:"discounted_surrender_penalty_charge"`
	DiscountedSurrenderPenaltyChargeAdjusted      float64 `json:"discounted_surrender_penalty_charge_adjusted"`
	DiscountedPremiumAdvisoryCost                 float64 `json:"discounted_premium_advisory_cost"`
	DiscountedPremiumAdvisoryCostAdjusted         float64 `json:"discounted_premium_advisory_cost_adjusted"`
	DiscountedFundAdvisoryCost                    float64 `json:"discounted_fund_advisory_cost"`
	DiscountedFundAdvisoryCostAdjusted            float64 `json:"discounted_fund_advisory_cost_adjusted"`
	DiscountedGuaranteeCost                       float64 `json:"discounted_guarantee_cost"`
	DiscountedGuaranteeCostAdjusted               float64 `json:"discounted_guarantee_cost_adjusted"`
	DiscountedSurrenderOutgo                      float64 `json:"discounted_surrender_outgo"`
	DiscountedSurrenderOutgoAdjusted              float64 `json:"discounted_surrender_outgo_adjusted"`
	DiscountedMaturityOutgo                       float64 `json:"discounted_maturity_outgo"`
	DiscountedMaturityOutgoAdjusted               float64 `json:"discounted_maturity_outgo_adjusted"`
	DiscountedCommission                          float64 `json:"discounted_commission"`
	DiscountedCommissionAdjusted                  float64 `json:"discounted_commission_adjusted"`
	DiscountedRenewalCommission                   float64 `json:"discounted_renewal_commission"`
	DiscountedRenewalCommissionAdjusted           float64 `json:"discounted_renewal_commission_adjusted"`
	DiscountedClawBack                            float64 `json:"discounted_claw_back"`
	DiscountedClawBackAdjusted                    float64 `json:"discounted_claw_back_adjusted"`
	DiscountedDeathOutgo                          float64 `json:"discounted_death_outgo"`
	DiscountedDeathOutgoAdjusted                  float64 `json:"discounted_death_outgo_adjusted"`
	DiscountedAccidentalDeathOutgo                float64 `json:"discounted_accidental_death_outgo"`
	DiscountedAccidentalDeathOutgoAdjusted        float64 `json:"discounted_accidental_death_outgo_adjusted"`
	DiscountedCashBackOnSurvival                  float64 `json:"discounted_cash_back_on_survival"`
	DiscountedCashBackOnSurvivalAdjusted          float64 `json:"discounted_cash_back_on_survival_adjusted"`
	DiscountedCashBackOnDeath                     float64 `json:"discounted_cash_back_on_death"`
	DiscountedCashBackOnDeathAdjusted             float64 `json:"discounted_cash_back_on_death_adjusted"`
	DiscountedRider                               float64 `json:"discounted_rider"`
	DiscountedRiderAdjusted                       float64 `json:"discounted_rider_adjusted"`
	DiscountedInitialExpenses                     float64 `json:"discounted_initial_expenses"`
	DiscountedInitialExpensesAdjusted             float64 `json:"discounted_initial_expenses_adjusted"`
	DiscountedRenewalExpenses                     float64 `json:"discounted_renewal_expenses"`
	DiscountedRenewalExpensesAdjusted             float64 `json:"discounted_renewal_expenses_adjusted"`
	DiscountedAnnuityOutgo                        float64 `json:"discounted_annuity_outgo"`
	DiscountedAnnuityOutgoAdjusted                float64 `json:"discounted_annuity_outgo_adjusted"`
	DiscountedInvestmentIncome                    float64 `json:"discounted_investment_income"`
	DiscountedInvestmentIncomeAdjusted            float64 `json:"discounted_investment_income_adjusted"`
	DiscountedProfitAdjustment                    float64 `json:"discounted_profit_adjustment"`
	DiscountedProfitAdjustmentAdjusted            float64 `json:"discounted_profit_adjustment_adjusted"`
	DiscountedProfit                              float64 `json:"discounted_profit"`
	DiscountedProfitAdjusted                      float64 `json:"discounted_profit_adjusted"`
	DiscountedCorporateTax                        float64 `json:"discounted_corporate_tax"`
	DiscountedCorporateTaxAdjusted                float64 `json:"discounted_corporate_tax_adjusted"`
	DiscountedAnnuityFactor                       float64 `json:"discounted_annuity_factor"`
	DiscountedCashOutflow                         float64 `json:"discounted_cash_outflow"`
	DiscountedCashOutflowExclAcquisition          float64 `json:"discounted_cash_outflow_excl_acquisition"`
	DiscountedAcquisitionCost                     float64 `json:"discounted_acquisition_cost"`
	DiscountedCashInflow                          float64 `json:"discounted_cash_inflow"`
	SumCoverageUnits                              float64 `json:"sum_coverage_units"`
	DiscountedCoverageUnits                       float64 `json:"discounted_coverage_units"`
	CededSumAssured                               float64 `json:"ceded_sum_assured"`
	ReinsurancePremium                            float64 `json:"reinsurance_premium"`
	ReinsurancePremiumAdjusted                    float64 `json:"reinsurance_premium_adjusted"`
	ReinsuranceCedingCommission                   float64 `json:"reinsurance_ceding_commission"`
	ReinsuranceCedingCommissionAdjusted           float64 `json:"reinsurance_ceding_commission_adjusted"`
	ReinsuranceClaims                             float64 `json:"reinsurance_claims"`
	ReinsuranceClaimsAdjusted                     float64 `json:"reinsurance_claims_adjusted"`
	NetReinsuranceCashflow                        float64 `json:"net_reinsurance_cashflow"`
	NetReinsuranceCashflowAdjusted                float64 `json:"net_reinsurance_cashflow_adjusted"`
	DiscountedReinsurancePremium                  float64 `json:"discounted_reinsurance_premium"`
	DiscountedReinsurancePremiumAdjusted          float64 `json:"discounted_reinsurance_premium_adjusted"`
	DiscountedReinsuranceCedingCommission         float64 `json:"discounted_reinsurance_ceding_commission"`
	DiscountedReinsuranceCedingCommissionAdjusted float64 `json:"discounted_reinsurance_ceding_commission_adjusted"`
	DiscountedReinsuranceClaims                   float64 `json:"discounted_reinsurance_claims"`
	DiscountedReinsuranceClaimsAdjusted           float64 `json:"discounted_reinsurance_claims_adjusted"`
	DiscountedNetReinsuranceCashflow              float64 `json:"discounted_net_reinsurance_cashflow"`
	DiscountedNetReinsuranceCashflowAdjusted      float64 `json:"discounted_net_reinsurance_cashflow_adjusted"`
	RunBasis                                      string  `json:"run_basis"`
}

type ScopedAggregatedProjection struct {
	ID              int     `json:"id" gorm:"primary_key"`
	JobProductID    int     `json:"job_product_id" gorm:"index"`
	RunType         int     `json:"projection_type" gorm:"index"` // 0 for valuation, 1 for pricing
	RunDate         string  `json:"run_date"`
	ProductCode     string  `json:"product_code" gorm:"index"`
	PolicyCount     float64 `json:"policy_count"`
	SpCode          int     `json:"sp_code" gorm:"index"`
	IFRS17Group     string
	ProjectionMonth int `gorm:"column:projection_month"`
	PremiumIncome   float64
	//PremiumIncomeAdjusted                  float64
	PremiumNotReceivedLapse float64
	//PremiumNotReceivedLapseAdjusted        float64
	Commission float64
	//CommissionAdjusted                     float64
	RenewalCommission float64
	//RenewalCommissionAdjusted              float64
	ClawBack float64
	//ClawBackAdjusted                       float64
	NonLifeClaimsOutgo float64
	//NonLifeClaimsOutgoAdjusted             float64
	DeathOutgo float64
	//DeathOutgoAdjusted                     float64
	AccidentalDeathOutgo float64
	//AccidentalDeathOutgoAdjusted           float64
	CashBackOnSurvival float64
	//CashBackOnSurvivalAdjusted             float64
	CashBackOnDeath float64
	//CashBackOnDeathAdjusted                float64
	DisabilityOutgo float64
	//DisabilityOutgoAdjusted                float64
	RetrenchmentOutgo float64
	//RetrenchmentOutgoAdjusted              float64
	Rider float64
	//RiderAdjusted   float64
	InitialExpenses float64
	//InitialExpensesAdjusted                float64
	RenewalExpenses float64
	//RenewalExpensesAdjusted                float64
	NetCashFlow float64
	//NetCashFlowAdjusted                    float64
	Reserves float64
	//ReservesAdjusted                       float64
	ChangeInReserves float64
	//ChangeInReservesAdjusted               float64
	InvestmentIncome float64
	//InvestmentIncomeAdjusted               float64
	Profit float64
	//ProfitAdjusted                         float64
	CoverageUnits           float64
	DiscountedPremiumIncome float64 `json:"discounted_premium_income"`
	//DiscountedPremiumIncomeAdjusted        float64 `json:"discounted_premium_income_adjusted"`
	DiscountedPremiumNotReceived float64 `json:"discounted_premium_not_received"`
	//DiscountedPremiumNotReceivedAdjusted   float64 `json:"discounted_premium_not_received_adjusted"`
	DiscountedCommission float64 `json:"discounted_commission"`
	//DiscountedCommissionAdjusted           float64 `json:"discounted_commission_adjusted"`
	DiscountedRenewalCommission float64 `json:"discounted_renewal_commission"`
	//DiscountedRenewalCommissionAdjusted           float64 `json:"discounted_renewal_commission_adjusted"`
	DiscountedClawBack float64 `json:"discounted_claw_back"`
	//DiscountedClawBackAdjusted             float64 `json:"discounted_claw_back_adjusted"`
	DiscountedDeathOutgo float64 `json:"discounted_death_outgo"`
	//DiscountedDeathOutgoAdjusted           float64 `json:"discounted_death_outgo_adjusted"`
	DiscountedAccidentalDeathOutgo float64 `json:"discounted_accidental_death_outgo"`
	//DiscountedAccidentalDeathOutgoAdjusted float64 `json:"discounted_accidental_death_outgo_adjusted"`
	DiscountedMorbidityOutgo float64 `json:"discounted_morbidity_outgo"`
	//DiscountedMorbidityOutgoAdjusted       float64 `json:"discounted_morbidity_outgo_adjusted"`
	DiscountedNonLifeOutgo float64 `json:"discounted_non_life_outgo"`
	DiscountedOutgo        float64 `json:"discounted_outgo"`
	//DiscountedOutgoAdjusted                float64 `json:"discounted_outgo_adjusted"`
	DiscountedCashBackOnSurvival float64 `json:"discounted_cash_back_on_survival"`
	//DiscountedCashBackOnSurvivalAdjusted   float64 `json:"discounted_cash_back_on_survival_adjusted"`
	DiscountedCashBackOnDeath float64 `json:"discounted_cash_back_on_death"`
	//DiscountedCashBackOnDeathAdjusted      float64 `json:"discounted_cash_back_on_death_adjusted"`
	DiscountedSurrenderOutgo float64 `json:"discounted_surrender_outgo"`
	DiscountedAnnuityOutgo   float64 `json:"discounted_annuity_outgo"`
	DiscountedRider          float64 `json:"discounted_rider"`
	//DiscountedRiderAdjusted                float64 `json:"discounted_rider_adjusted"`
	DiscountedInitialExpenses float64 `json:"discounted_initial_expenses"`
	//DiscountedInitialExpensesAdjusted      float64 `json:"discounted_initial_expenses_adjusted"`
	DiscountedRenewalExpenses float64 `json:"discounted_renewal_expenses"`
	//DiscountedRenewalExpensesAdjusted      float64 `json:"discounted_renewal_expenses_adjusted"`
	DiscountedProfit float64 `json:"discounted_profit"`
	SumAtRisk        float64 `json:"sum_at_risk"`
	//DiscountedProfitAdjusted               float64 `json:"discounted_profit_adjusted"`
	DiscountedAnnuityFactor               float64 `json:"annuity_factor"`
	DiscountedCashOutflow                 float64 `json:"discounted_cash_outflow"`
	DiscountedCashOutflowExclAcquisition  float64 `json:"discounted_cash_outflow_excl_acquisition"`
	DiscountedAcquisitionCost             float64 `json:"discounted_acquisition_cost"`
	DiscountedCashInflow                  float64 `json:"discounted_cash_inflow"`
	SumCoverageUnits                      float64 `json:"sum_coverage_units"`
	DiscountedCoverageUnits               float64 `json:"discounted_coverage_units"`
	CededSumAssured                       float64 `json:"ceded_sum_assured"`
	ReinsurancePremium                    float64 `json:"reinsurance_premium"`
	ReinsuranceCedingCommission           float64 `json:"reinsurance_ceding_commission"`
	ReinsuranceClaims                     float64 `json:"reinsurance_claims"`
	NetReinsuranceCashflow                float64 `json:"net_reinsurance_cashflow"`
	DiscountedNetReinsuranceCashflow      float64 `json:"discounted_net_reinsurance_cashflow"`
	DiscountedReinsurancePremium          float64 `json:"discounted_reinsurance_premium"`
	DiscountedReinsuranceCedingCommission float64 `json:"discounted_reinsurance_ceding_commission"`
	DiscountedReinsuranceClaims           float64 `json:"discounted_reinsurance_claims"`
	RunBasis                              string  `json:"run_basis"`
	//Manual Import of SAP
	YieldCurveCode                          string  `json:"yield_curve_code" csv:"yield_curve_code"`
	DiscountedCashOutflowM0                 float64 `json:"discounted_cash_outflow_m0" csv:"discounted_cash_outflow_m0"`
	DiscountedCashOutflowM12                float64 `json:"discounted_cash_outflow_m12" csv:"discounted_cash_outflow_m12"`
	DiscountedCashOutflowExclAcquisitionM0  float64 `json:"discounted_cash_outflow_excl_acquisition_m0" csv:"discounted_cash_outflow_excl_acquisition_m0"`
	DiscountedCashOutflowExclAcquisitionM12 float64 `json:"discounted_cash_outflow_excl_acquisition_m12" csv:"discounted_cash_outflow_excl_acquisition_m12"`
	DiscountedAcquisitionCostM0             float64 `json:"discounted_acquisition_cost_m0" csv:"discounted_acquisition_cost_m0"`
	DiscountedAcquisitionCostM12            float64 `json:"discounted_acquisition_cost_m12" csv:"discounted_acquisition_cost_m12"`
	DiscountedCashInflowM0                  float64 `json:"discounted_cash_inflow_m0" csv:"discounted_cash_inflow_m0"`
	DiscountedCashInflowM12                 float64 `json:"discounted_cash_inflow_m12" csv:"discounted_cash_inflow_m12"`
	SumCoverageUnitsM0                      float64 `json:"sum_coverage_units_m0" csv:"sum_coverage_units_m0"`
	SumCoverageUnitsM12                     float64 `json:"sum_coverage_units_m12" csv:"sum_coverage_units_m12"`
	SumCoverageUnitsM24                     float64 `json:"sum_coverage_units_m24" csv:"sum_coverage_units_m24"`
	SumCoverageUnitsM36                     float64 `json:"sum_coverage_units_m36" csv:"sum_coverage_units_m36"`
	SumCoverageUnitsM48                     float64 `json:"sum_coverage_units_m48" csv:"sum_coverage_units_m48"`
	SumCoverageUnitsM60                     float64 `json:"sum_coverage_units_m60" csv:"sum_coverage_units_m60"`
	SumCoverageUnitsM72                     float64 `json:"sum_coverage_units_m72" csv:"sum_coverage_units_m72"`
	SumCoverageUnitsM84                     float64 `json:"sum_coverage_units_m84" csv:"sum_coverage_units_m84"`
	SumCoverageUnitsM96                     float64 `json:"sum_coverage_units_m96" csv:"sum_coverage_units_m96"`
	SumCoverageUnitsM108                    float64 `json:"sum_coverage_units_m108" csv:"sum_coverage_units_m108"`
	SumCoverageUnitsM120                    float64 `json:"sum_coverage_units_m120" csv:"sum_coverage_units_m120"`
	DiscountedCoverageUnitsM0               float64 `json:"discounted_coverage_units_m0" csv:"discounted_coverage_units_m0"`
	DiscountedCoverageUnitsM12              float64 `json:"discounted_coverage_units_m12" csv:"discounted_coverage_units_m12"`
	DiscountedCoverageUnitsM24              float64 `json:"discounted_coverage_units_m24" csv:"discounted_coverage_units_m24"`
	DiscountedCoverageUnitsM36              float64 `json:"discounted_coverage_units_m36" csv:"discounted_coverage_units_m36"`
	DiscountedCoverageUnitsM48              float64 `json:"discounted_coverage_units_m48" csv:"discounted_coverage_units_m48"`
	DiscountedCoverageUnitsM60              float64 `json:"discounted_coverage_units_m60" csv:"discounted_coverage_units_m60"`
	DiscountedCoverageUnitsM72              float64 `json:"discounted_coverage_units_m72" csv:"discounted_coverage_units_m72"`
	DiscountedCoverageUnitsM84              float64 `json:"discounted_coverage_units_m84" csv:"discounted_coverage_units_m84"`
	DiscountedCoverageUnitsM96              float64 `json:"discounted_coverage_units_m96" csv:"discounted_coverage_units_m96"`
	DiscountedCoverageUnitsM108             float64 `json:"discounted_coverage_units_m108" csv:"discounted_coverage_units_m108"`
	DiscountedCoverageUnitsM120             float64 `json:"discounted_coverage_units_m120" csv:"discounted_coverage_units_m120"`
	BelM0                                   float64 `json:"bel_m0" csv:"bel_m0"`
	BelM12                                  float64 `json:"bel_m12" csv:"bel_m12"`
	RiskAdjustmentM0                        float64 `json:"risk_adjustment_m0" csv:"risk_adjustment_m0"`
	RiskAdjustmentM12                       float64 `json:"risk_adjustment_m12" csv:"risk_adjustment_m12"`
	CashInflowM1                            float64 `json:"cash_inflow_m1" csv:"cash_inflow_m1"`
	CashInflowM2                            float64 `json:"cash_inflow_m2" csv:"cash_inflow_m2"`
	CashInflowM3                            float64 `json:"cash_inflow_m3" csv:"cash_inflow_m3"`
	CashInflowM4                            float64 `json:"cash_inflow_m4" csv:"cash_inflow_m4"`
	CashInflowM5                            float64 `json:"cash_inflow_m5" csv:"cash_inflow_m5"`
	CashInflowM6                            float64 `json:"cash_inflow_m6" csv:"cash_inflow_m6"`
	CashInflowM7                            float64 `json:"cash_inflow_m7" csv:"cash_inflow_m7"`
	CashInflowM8                            float64 `json:"cash_inflow_m8" csv:"cash_inflow_m8"`
	CashInflowM9                            float64 `json:"cash_inflow_m9" csv:"cash_inflow_m9"`
	CashInflowM10                           float64 `json:"cash_inflow_m10" csv:"cash_inflow_m10"`
	CashInflowM11                           float64 `json:"cash_inflow_m11" csv:"cash_inflow_m11"`
	CashInflowM12                           float64 `json:"cash_inflow_m12" csv:"cash_inflow_m12"`
	MortalityOutgoM1                        float64 `json:"mortality_outgo_m1" csv:"mortality_outgo_m1"`
	MortalityOutgoM2                        float64 `json:"mortality_outgo_m2" csv:"mortality_outgo_m2"`
	MortalityOutgoM3                        float64 `json:"mortality_outgo_m3" csv:"mortality_outgo_m3"`
	MortalityOutgoM4                        float64 `json:"mortality_outgo_m4" csv:"mortality_outgo_m4"`
	MortalityOutgoM5                        float64 `json:"mortality_outgo_m5" csv:"mortality_outgo_m5"`
	MortalityOutgoM6                        float64 `json:"mortality_outgo_m6" csv:"mortality_outgo_m6"`
	MortalityOutgoM7                        float64 `json:"mortality_outgo_m7" csv:"mortality_outgo_m7"`
	MortalityOutgoM8                        float64 `json:"mortality_outgo_m8" csv:"mortality_outgo_m8"`
	MortalityOutgoM9                        float64 `json:"mortality_outgo_m9" csv:"mortality_outgo_m9"`
	MortalityOutgoM10                       float64 `json:"mortality_outgo_m10" csv:"mortality_outgo_m10"`
	MortalityOutgoM11                       float64 `json:"mortality_outgo_m11" csv:"mortality_outgo_m11"`
	MortalityOutgoM12                       float64 `json:"mortality_outgo_m12" csv:"mortality_outgo_m12"`
	RetrenchmentOutgoM1                     float64 `json:"retrenchment_outgo_m1" csv:"retrenchment_outgo_m1"`
	RetrenchmentOutgoM2                     float64 `json:"retrenchment_outgo_m2" csv:"retrenchment_outgo_m2"`
	RetrenchmentOutgoM3                     float64 `json:"retrenchment_outgo_m3" csv:"retrenchment_outgo_m3"`
	RetrenchmentOutgoM4                     float64 `json:"retrenchment_outgo_m4" csv:"retrenchment_outgo_m4"`
	RetrenchmentOutgoM5                     float64 `json:"retrenchment_outgo_m5" csv:"retrenchment_outgo_m5"`
	RetrenchmentOutgoM6                     float64 `json:"retrenchment_outgo_m6" csv:"retrenchment_outgo_m6"`
	RetrenchmentOutgoM7                     float64 `json:"retrenchment_outgo_m7" csv:"retrenchment_outgo_m7"`
	RetrenchmentOutgoM8                     float64 `json:"retrenchment_outgo_m8" csv:"retrenchment_outgo_m8"`
	RetrenchmentOutgoM9                     float64 `json:"retrenchment_outgo_m9" csv:"retrenchment_outgo_m9"`
	RetrenchmentOutgoM10                    float64 `json:"retrenchment_outgo_m10" csv:"retrenchment_outgo_m10"`
	RetrenchmentOutgoM11                    float64 `json:"retrenchment_outgo_m11" csv:"retrenchment_outgo_m11"`
	RetrenchmentOutgoM12                    float64 `json:"retrenchment_outgo_m12" csv:"retrenchment_outgo_m12"`
	MorbidityOutgoM1                        float64 `json:"morbidity_outgo_m1" csv:"morbidity_outgo_m1"`
	MorbidityOutgoM2                        float64 `json:"morbidity_outgo_m2" csv:"morbidity_outgo_m2"`
	MorbidityOutgoM3                        float64 `json:"morbidity_outgo_m3" csv:"morbidity_outgo_m3"`
	MorbidityOutgoM4                        float64 `json:"morbidity_outgo_m4" csv:"morbidity_outgo_m4"`
	MorbidityOutgoM5                        float64 `json:"morbidity_outgo_m5" csv:"morbidity_outgo_m5"`
	MorbidityOutgoM6                        float64 `json:"morbidity_outgo_m6" csv:"morbidity_outgo_m6"`
	MorbidityOutgoM7                        float64 `json:"morbidity_outgo_m7" csv:"morbidity_outgo_m7"`
	MorbidityOutgoM8                        float64 `json:"morbidity_outgo_m8" csv:"morbidity_outgo_m8"`
	MorbidityOutgoM9                        float64 `json:"morbidity_outgo_m9" csv:"morbidity_outgo_m9"`
	MorbidityOutgoM10                       float64 `json:"morbidity_outgo_m10" csv:"morbidity_outgo_m10"`
	MorbidityOutgoM11                       float64 `json:"morbidity_outgo_m11" csv:"morbidity_outgo_m11"`
	MorbidityOutgoM12                       float64 `json:"morbidity_outgo_m12" csv:"morbidity_outgo_m12"`
	NonLifeClaimsOutgoM1                    float64 `json:"non_life_claims_outgo_m1" csv:"non_life_claims_outgo_m1"`
	NonLifeClaimsOutgoM2                    float64 `json:"non_life_claims_outgo_m2" csv:"non_life_claims_outgo_m2"`
	NonLifeClaimsOutgoM3                    float64 `json:"non_life_claims_outgo_m3" csv:"non_life_claims_outgo_m3"`
	NonLifeClaimsOutgoM4                    float64 `json:"non_life_claims_outgo_m4" csv:"non_life_claims_outgo_m4"`
	NonLifeClaimsOutgoM5                    float64 `json:"non_life_claims_outgo_m5" csv:"non_life_claims_outgo_m5"`
	NonLifeClaimsOutgoM6                    float64 `json:"non_life_claims_outgo_m6" csv:"non_life_claims_outgo_m6"`
	NonLifeClaimsOutgoM7                    float64 `json:"non_life_claims_outgo_m7" csv:"non_life_claims_outgo_m7"`
	NonLifeClaimsOutgoM8                    float64 `json:"non_life_claims_outgo_m8" csv:"non_life_claims_outgo_m8"`
	NonLifeClaimsOutgoM9                    float64 `json:"non_life_claims_outgo_m9" csv:"non_life_claims_outgo_m9"`
	NonLifeClaimsOutgoM10                   float64 `json:"non_life_claims_outgo_m10" csv:"non_life_claims_outgo_m10"`
	NonLifeClaimsOutgoM11                   float64 `json:"non_life_claims_outgo_m11" csv:"non_life_claims_outgo_m11"`
	NonLifeClaimsOutgoM12                   float64 `json:"non_life_claims_outgo_m12" csv:"non_life_claims_outgo_m12"`
	RenewalExpenseOutgoM1                   float64 `json:"renewal_expense_outgo_m1" csv:"renewal_expense_outgo_m1"`
	RenewalExpenseOutgoM2                   float64 `json:"renewal_expense_outgo_m2" csv:"renewal_expense_outgo_m2"`
	RenewalExpenseOutgoM3                   float64 `json:"renewal_expense_outgo_m3" csv:"renewal_expense_outgo_m3"`
	RenewalExpenseOutgoM4                   float64 `json:"renewal_expense_outgo_m4" csv:"renewal_expense_outgo_m4"`
	RenewalExpenseOutgoM5                   float64 `json:"renewal_expense_outgo_m5" csv:"renewal_expense_outgo_m5"`
	RenewalExpenseOutgoM6                   float64 `json:"renewal_expense_outgo_m6" csv:"renewal_expense_outgo_m6"`
	RenewalExpenseOutgoM7                   float64 `json:"renewal_expense_outgo_m7" csv:"renewal_expense_outgo_m7"`
	RenewalExpenseOutgoM8                   float64 `json:"renewal_expense_outgo_m8" csv:"renewal_expense_outgo_m8"`
	RenewalExpenseOutgoM9                   float64 `json:"renewal_expense_outgo_m9" csv:"renewal_expense_outgo_m9"`
	RenewalExpenseOutgoM10                  float64 `json:"renewal_expense_outgo_m10" csv:"renewal_expense_outgo_m10"`
	RenewalExpenseOutgoM11                  float64 `json:"renewal_expense_outgo_m11" csv:"renewal_expense_outgo_m11"`
	RenewalExpenseOutgoM12                  float64 `json:"renewal_expense_outgo_m12" csv:"renewal_expense_outgo_m12"`
	DiscountedEntityShareM0                 float64 `json:"discounted_entity_share_m0" csv:"discounted_entity_share_m0"`
	DiscountedEntityShareM12                float64 `json:"discounted_entity_share_m12" csv:"discounted_entity_share_m12"`
	DiscountedCashInflowsNotVaryM0          float64 `json:"discounted_cash_inflows_not_vary_m0" csv:"discounted_cash_inflows_not_vary_m0"`
	DiscountedCashInflowsNotVaryM12         float64 `json:"discounted_cash_inflows_not_vary_m12" csv:"discounted_cash_inflows_not_vary_m12"`
	TvogM0                                  float64 `json:"tvog_m0" csv:"tvog_m0"`
	TvogM12                                 float64 `json:"tvog_m12" csv:"tvog_m12"`
	UnitFundM0                              float64 `json:"unit_fund_m0" csv:"unit_fund_m0"`
	EntityShareM1                           float64 `json:"entity_share_m1" csv:"entity_share_m1"`
	EntityShareM2                           float64 `json:"entity_share_m2" csv:"entity_share_m2"`
	EntityShareM3                           float64 `json:"entity_share_m3" csv:"entity_share_m3"`
	EntityShareM4                           float64 `json:"entity_share_m4" csv:"entity_share_m4"`
	EntityShareM5                           float64 `json:"entity_share_m5" csv:"entity_share_m5"`
	EntityShareM6                           float64 `json:"entity_share_m6" csv:"entity_share_m6"`
	EntityShareM7                           float64 `json:"entity_share_m7" csv:"entity_share_m7"`
	EntityShareM8                           float64 `json:"entity_share_m8" csv:"entity_share_m8"`
	EntityShareM9                           float64 `json:"entity_share_m9" csv:"entity_share_m9"`
	EntityShareM10                          float64 `json:"entity_share_m10" csv:"entity_share_m10"`
	EntityShareM11                          float64 `json:"entity_share_m11" csv:"entity_share_m11"`
	EntityShareM12                          float64 `json:"entity_share_m12" csv:"entity_share_m12"`
	CashInflowsNotVaryM1                    float64 `json:"cash_inflows_not_vary_m1" csv:"cash_inflows_not_vary_m1"`
	CashInflowsNotVaryM2                    float64 `json:"cash_inflows_not_vary_m2" csv:"cash_inflows_not_vary_m2"`
	CashInflowsNotVaryM3                    float64 `json:"cash_inflows_not_vary_m3" csv:"cash_inflows_not_vary_m3"`
	CashInflowsNotVaryM4                    float64 `json:"cash_inflows_not_vary_m4" csv:"cash_inflows_not_vary_m4"`
	CashInflowsNotVaryM5                    float64 `json:"cash_inflows_not_vary_m5" csv:"cash_inflows_not_vary_m5"`
	CashInflowsNotVaryM6                    float64 `json:"cash_inflows_not_vary_m6" csv:"cash_inflows_not_vary_m6"`
	CashInflowsNotVaryM7                    float64 `json:"cash_inflows_not_vary_m7" csv:"cash_inflows_not_vary_m7"`
	CashInflowsNotVaryM8                    float64 `json:"cash_inflows_not_vary_m8" csv:"cash_inflows_not_vary_m8"`
	CashInflowsNotVaryM9                    float64 `json:"cash_inflows_not_vary_m9" csv:"cash_inflows_not_vary_m9"`
	CashInflowsNotVaryM10                   float64 `json:"cash_inflows_not_vary_m10" csv:"cash_inflows_not_vary_m10"`
	CashInflowsNotVaryM11                   float64 `json:"cash_inflows_not_vary_m11" csv:"cash_inflows_not_vary_m11"`
	CashInflowsNotVaryM12                   float64 `json:"cash_inflows_not_vary_m12" csv:"cash_inflows_not_vary_m12"`
	OtherCashOutflowsM1                     float64 `json:"other_cash_outflows_m1" csv:"other_cash_outflows_m1"`
	OtherCashOutflowsM2                     float64 `json:"other_cash_outflows_m2" csv:"other_cash_outflows_m2"`
	OtherCashOutflowsM3                     float64 `json:"other_cash_outflows_m3" csv:"other_cash_outflows_m3"`
	OtherCashOutflowsM4                     float64 `json:"other_cash_outflows_m4" csv:"other_cash_outflows_m4"`
	OtherCashOutflowsM5                     float64 `json:"other_cash_outflows_m5" csv:"other_cash_outflows_m5"`
	OtherCashOutflowsM6                     float64 `json:"other_cash_outflows_m6" csv:"other_cash_outflows_m6"`
	OtherCashOutflowsM7                     float64 `json:"other_cash_outflows_m7" csv:"other_cash_outflows_m7"`
	OtherCashOutflowsM8                     float64 `json:"other_cash_outflows_8" csv:"other_cash_outflows_8"`
	OtherCashOutflowsM9                     float64 `json:"other_cash_outflows_m9" csv:"other_cash_outflows_m9"`
	OtherCashOutflowsM10                    float64 `json:"other_cash_outflows_m10" csv:"other_cash_outflows_m10"`
	OtherCashOutflowsM11                    float64 `json:"other_cash_outflows_m11" csv:"other_cash_outflows_m11"`
	OtherCashOutflowsM12                    float64 `json:"other_cash_outflows_m12" csv:"other_cash_outflows_m12"`
	Treaty1BelM0                            float64 `json:"treaty1_bel_m0" csv:"treaty1_bel_m0"`
	Treaty1BelM12                           float64 `json:"treaty1_bel_m12" csv:"treaty1_bel_m12"`
	Treaty1RiskAdjustmentM0                 float64 `json:"treaty1_risk_adjustment_m0" csv:"treaty1_risk_adjustment_m0"`
	Treaty1RiskAdjustmentM12                float64 `json:"treaty1_risk_adjustment_m12" csv:"treaty1_risk_adjustment_m12"`
	Treaty1DiscountedCashOutflowM0          float64 `json:"treaty1_discounted_cash_outflow_m0" csv:"treaty1_discounted_cash_outflow_m0"`
	Treaty1DiscountedCashOutflowM12         float64 `json:"treaty1_discounted_cash_outflow_m12" csv:"treaty1_discounted_cash_outflow_m12"`
	Treaty1DiscountedCashInflowM0           float64 `json:"treaty1_discounted_cash_inflow_m0" csv:"treaty1_discounted_cash_inflow_m0"`
	Treaty1DiscountedCashInflowM12          float64 `json:"treaty1_discounted_cash_inflow_m12" csv:"treaty1_discounted_cash_inflow_m12"`
	Treaty2BelM0                            float64 `json:"treaty2_bel_m0" csv:"treaty2_bel_m0"`
	Treaty2BelM12                           float64 `json:"treaty2_bel_m12" csv:"treaty2_bel_m12"`
	Treaty2RiskAdjustmentM0                 float64 `json:"treaty2_risk_adjustment_m0" csv:"treaty2_risk_adjustment_m0"`
	Treaty2RiskAdjustmentM12                float64 `json:"treaty2_risk_adjustment_m12" csv:"treaty2_risk_adjustment_m12"`
	Treaty2DiscountedCashOutflowM0          float64 `json:"treaty2_discounted_cash_outflow_m0" csv:"treaty2_discounted_cash_outflow_m0"`
	Treaty2DiscountedCashOutflowM12         float64 `json:"treaty2_discounted_cash_outflow_m12" csv:"treaty2_discounted_cash_outflow_m12"`
	Treaty2DiscountedCashInflowM0           float64 `json:"treaty2_discounted_cash_inflow_m0" csv:"treaty2_discounted_cash_inflow_m0"`
	Treaty2DiscountedCashInflowM12          float64 `json:"treaty2_discounted_cash_inflow_m12" csv:"treaty2_discounted_cash_inflow_m12"`
	Treaty3BelM0                            float64 `json:"treaty3_bel_m0" csv:"treaty3_bel_m0"`
	Treaty3BelM12                           float64 `json:"treaty3_bel_m12" csv:"treaty3_bel_m12"`
	Treaty3RiskAdjustmentM0                 float64 `json:"treaty3_risk_adjustment_m0" csv:"treaty3_risk_adjustment_m0"`
	Treaty3RiskAdjustmentM12                float64 `json:"treaty3_risk_adjustment_m12" csv:"treaty3_risk_adjustment_m12"`
	Treaty3DiscountedCashOutflowM0          float64 `json:"treaty3_discounted_cash_outflow_m0" csv:"treaty3_discounted_cash_outflow_m0"`
	Treaty3DiscountedCashOutflowM12         float64 `json:"treaty3_discounted_cash_outflow_m12" csv:"treaty3_discounted_cash_outflow_m12"`
	Treaty3DiscountedCashInflowM0           float64 `json:"treaty3_discounted_cash_inflow_m0" csv:"treaty3_discounted_cash_inflow_m0"`
	Treaty3DiscountedCashInflowM12          float64 `json:"treaty3_discounted_cash_inflow_m12" csv:"treaty3_discounted_cash_inflow_m12"`
	Treaty1CashInflowM1                     float64 `json:"treaty1_cash_inflow_m1" csv:"treaty1_cash_inflow_m1"`
	Treaty1CashInflowM2                     float64 `json:"treaty1_cash_inflow_m2" csv:"treaty1_cash_inflow_m2"`
	Treaty1CashInflowM3                     float64 `json:"treaty1_cash_inflow_m3" csv:"treaty1_cash_inflow_m3"`
	Treaty1CashInflowM4                     float64 `json:"treaty1_cash_inflow_m4" csv:"treaty1_cash_inflow_m4"`
	Treaty1CashInflowM5                     float64 `json:"treaty1_cash_inflow_m5" csv:"treaty1_cash_inflow_m5"`
	Treaty1CashInflowM6                     float64 `json:"treaty1_cash_inflow_m6" csv:"treaty1_cash_inflow_m6"`
	Treaty1CashInflowM7                     float64 `json:"treaty1_cash_inflow_m7" csv:"treaty1_cash_inflow_m7"`
	Treaty1CashInflowM8                     float64 `json:"treaty1_cash_inflow_m8" csv:"treaty1_cash_inflow_m8"`
	Treaty1CashInflowM9                     float64 `json:"treaty1_cash_inflow_m9" csv:"treaty1_cash_inflow_m9"`
	Treaty1CashInflowM10                    float64 `json:"treaty1_cash_inflow_m10" csv:"treaty1_cash_inflow_m10"`
	Treaty1CashInflowM11                    float64 `json:"treaty1_cash_inflow_m11" csv:"treaty1_cash_inflow_m11"`
	Treaty1CashInflowM12                    float64 `json:"treaty1_cash_inflow_m12" csv:"treaty1_cash_inflow_m12"`
	Treaty1CashOutflowM1                    float64 `json:"treaty1_cash_outflow_m1" csv:"treaty1_cash_outflow_m1"`
	Treaty1CashOutflowM2                    float64 `json:"treaty1_cash_outflow_m2" csv:"treaty1_cash_outflow_m2"`
	Treaty1CashOutflowM3                    float64 `json:"treaty1_cash_outflow_m3" csv:"treaty1_cash_outflow_m3"`
	Treaty1CashOutflowM4                    float64 `json:"treaty1_cash_outflow_m4" csv:"treaty1_cash_outflow_m4"`
	Treaty1CashOutflowM5                    float64 `json:"treaty1_cash_outflow_m5" csv:"treaty1_cash_outflow_m5"`
	Treaty1CashOutflowM6                    float64 `json:"treaty1_cash_outflow_m6" csv:"treaty1_cash_outflow_m6"`
	Treaty1CashOutflowM7                    float64 `json:"treaty1_cash_outflow_m7" csv:"treaty1_cash_outflow_m7"`
	Treaty1CashOutflowM8                    float64 `json:"treaty1_cash_outflow_m8" csv:"treaty1_cash_outflow_m8"`
	Treaty1CashOutflowM9                    float64 `json:"treaty1_cash_outflow_m9" csv:"treaty1_cash_outflow_m9"`
	Treaty1CashOutflowM10                   float64 `json:"treaty1_cash_outflow_m10" csv:"treaty1_cash_outflow_m10"`
	Treaty1CashOutflowM11                   float64 `json:"treaty1_cash_outflow_m11" csv:"treaty1_cash_outflow_m11"`
	Treaty1CashOutflowM12                   float64 `json:"treaty1_cash_outflow_m12" csv:"treaty1_cash_outflow_m12"`
	Treaty2CashInflowM1                     float64 `json:"treaty2_cash_inflow_m1" csv:"treaty2_cash_inflow_m1"`
	Treaty2CashInflowM2                     float64 `json:"treaty2_cash_inflow_m2" csv:"treaty2_cash_inflow_m2"`
	Treaty2CashInflowM3                     float64 `json:"treaty2_cash_inflow_m3" csv:"treaty2_cash_inflow_m3"`
	Treaty2CashInflowM4                     float64 `json:"treaty2_cash_inflow_m4" csv:"treaty2_cash_inflow_m4"`
	Treaty2CashInflowM5                     float64 `json:"treaty2_cash_inflow_m5" csv:"treaty2_cash_inflow_m5"`
	Treaty2CashInflowM6                     float64 `json:"treaty2_cash_inflow_m6" csv:"treaty2_cash_inflow_m6"`
	Treaty2CashInflowM7                     float64 `json:"treaty2_cash_inflow_m7" csv:"treaty2_cash_inflow_m7"`
	Treaty2CashInflowM8                     float64 `json:"treaty2_cash_inflow_m8" csv:"treaty2_cash_inflow_m8"`
	Treaty2CashInflowM9                     float64 `json:"treaty2_cash_inflow_m9" csv:"treaty2_cash_inflow_m9"`
	Treaty2CashInflowM10                    float64 `json:"treaty2_cash_inflow_m10" csv:"treaty2_cash_inflow_m10"`
	Treaty2CashInflowM11                    float64 `json:"treaty2_cash_inflow_m11" csv:"treaty2_cash_inflow_m11"`
	Treaty2CashInflowM12                    float64 `json:"treaty2_cash_inflow_m12" csv:"treaty2_cash_inflow_m12"`
	Treaty2CashOutflowM1                    float64 `json:"treaty2_cash_outflow_m1" csv:"treaty2_cash_outflow_m1"`
	Treaty2CashOutflowM2                    float64 `json:"treaty2_cash_outflow_m2" csv:"treaty2_cash_outflow_m2"`
	Treaty2CashOutflowM3                    float64 `json:"treaty2_cash_outflow_m3" csv:"treaty2_cash_outflow_m3"`
	Treaty2CashOutflowM4                    float64 `json:"treaty2_cash_outflow_m4" csv:"treaty2_cash_outflow_m4"`
	Treaty2CashOutflowM5                    float64 `json:"treaty2_cash_outflow_m5" csv:"treaty2_cash_outflow_m5"`
	Treaty2CashOutflowM6                    float64 `json:"treaty2_cash_outflow_m6" csv:"treaty2_cash_outflow_m6"`
	Treaty2CashOutflowM7                    float64 `json:"treaty2_cash_outflow_m7" csv:"treaty2_cash_outflow_m7"`
	Treaty2CashOutflowM8                    float64 `json:"treaty2_cash_outflow_m8" csv:"treaty2_cash_outflow_m8"`
	Treaty2CashOutflowM9                    float64 `json:"treaty2_cash_outflow_m9" csv:"treaty2_cash_outflow_m9"`
	Treaty2CashOutflowM10                   float64 `json:"treaty2_cash_outflow_m10" csv:"treaty2_cash_outflow_m10"`
	Treaty2CashOutflowM11                   float64 `json:"treaty2_cash_outflow_m11" csv:"treaty2_cash_outflow_m11"`
	Treaty2CashOutflowM12                   float64 `json:"treaty2_cash_outflow_m12" csv:"treaty2_cash_outflow_m12"`
	Treaty3CashInflowM1                     float64 `json:"treaty3_cash_inflow_m1" csv:"treaty3_cash_inflow_m1"`
	Treaty3CashInflowM2                     float64 `json:"treaty3_cash_inflow_m2" csv:"treaty3_cash_inflow_m2"`
	Treaty3CashInflowM3                     float64 `json:"treaty3_cash_inflow_m3" csv:"treaty3_cash_inflow_m3"`
	Treaty3CashInflowM4                     float64 `json:"treaty3_cash_inflow_m4" csv:"treaty3_cash_inflow_m4"`
	Treaty3CashInflowM5                     float64 `json:"treaty3_cash_inflow_m5" csv:"treaty3_cash_inflow_m5"`
	Treaty3CashInflowM6                     float64 `json:"treaty3_cash_inflow_m6" csv:"treaty3_cash_inflow_m6"`
	Treaty3CashInflowM7                     float64 `json:"treaty3_cash_inflow_m7" csv:"treaty3_cash_inflow_m7"`
	Treaty3CashInflowM8                     float64 `json:"treaty3_cash_inflow_m8" csv:"treaty3_cash_inflow_m8"`
	Treaty3CashInflowM9                     float64 `json:"treaty3_cash_inflow_m9" csv:"treaty3_cash_inflow_m9"`
	Treaty3CashInflowM10                    float64 `json:"treaty3_cash_inflow_m10" csv:"treaty3_cash_inflow_m10"`
	Treaty3CashInflowM11                    float64 `json:"treaty3_cash_inflow_m11" csv:"treaty3_cash_inflow_m11"`
	Treaty3CashInflowM12                    float64 `json:"treaty3_cash_inflow_m12" csv:"treaty3_cash_inflow_m12"`
	Treaty3CashOutflowM1                    float64 `json:"treaty3_cash_outflow_m1" csv:"treaty3_cash_outflow_m1"`
	Treaty3CashOutflowM2                    float64 `json:"treaty3_cash_outflow_m2" csv:"treaty3_cash_outflow_m2"`
	Treaty3CashOutflowM3                    float64 `json:"treaty3_cash_outflow_m3" csv:"treaty3_cash_outflow_m3"`
	Treaty3CashOutflowM4                    float64 `json:"treaty3_cash_outflow_m4" csv:"treaty3_cash_outflow_m4"`
	Treaty3CashOutflowM5                    float64 `json:"treaty3_cash_outflow_m5" csv:"treaty3_cash_outflow_m5"`
	Treaty3CashOutflowM6                    float64 `json:"treaty3_cash_outflow_m6" csv:"treaty3_cash_outflow_m6"`
	Treaty3CashOutflowM7                    float64 `json:"treaty3_cash_outflow_m7" csv:"treaty3_cash_outflow_m7"`
	Treaty3CashOutflowM8                    float64 `json:"treaty3_cash_outflow_m8" csv:"treaty3_cash_outflow_m8"`
	Treaty3CashOutflowM9                    float64 `json:"treaty3_cash_outflow_m9" csv:"treaty3_cash_outflow_m9"`
	Treaty3CashOutflowM10                   float64 `json:"treaty3_cash_outflow_m10" csv:"treaty3_cash_outflow_m10"`
	Treaty3CashOutflowM11                   float64 `json:"treaty3_cash_outflow_m11" csv:"treaty3_cash_outflow_m11"`
	Treaty3CashOutflowM12                   float64 `json:"treaty3_cash_outflow_m12" csv:"treaty3_cash_outflow_m12"`

	RunId                                         int     `json:"-"`
	RunName                                       string  `json:"-"`
	ProductId                                     string  `json:"-"`
	PolicyNumber                                  string  `json:"-"`
	ProjectionYear                                int     `json:"-"`
	ValuationTimeMonth                            int     `json:"-"`
	CalendarMonth                                 int     `json:"-"`
	ValuationTimeYear                             float64 `json:"-"`
	MainMemberAgeNextBirthday                     int     `json:"-"`
	AgeNextBirthday                               int     `json:"-"`
	AccidentProportion                            float64 `json:"-"`
	InflationFactor                               float64 `json:"-"`
	InflationFactorAdjusted                       float64 `json:"-"`
	AnnuityEscalationRate                         float64 `json:"-"`
	PremiumEscalation                             float64 `json:"-"`
	SumAssuredEscalation                          float64 `json:"-"`
	AnnuityEscalation                             float64 `json:"-"`
	LapseMargin                                   float64 `json:"-"`
	PremiumWaiverOnFactor                         float64 `json:"-"`
	PaidUpOnFactor                                float64 `json:"-"`
	MainMemberMortalityRate                       float64 `json:"-"`
	MainMemberMortalityRateAdjusted               float64 `json:"-"`
	BaseLapse                                     float64 `json:"-"`
	BaseLapseAdjusted                             float64 `json:"-"`
	ContractingPartyAlivePortion                  float64 `json:"-"`
	ContractingPartyAlivePortionAdjusted          float64 `json:"-"`
	ContractingPartyPolicyLapse                   float64 `json:"-"`
	ContractingPartyPolicyLapseAdjusted           float64 `json:"-"`
	NonLifeMonthlyRiskRate                        float64 `json:"-"`
	NonLifeMonthlyRiskRateAdjusted                float64 `json:"-"`
	BaseMortalityRate                             float64 `json:"-"`
	BaseMortalityRateAdjusted                     float64 `json:"-"`
	BaseIndependentLapse                          float64 `json:"-"`
	BaseIndependentLapseAdjusted                  float64 `json:"-"`
	BaseRetrenchmentRate                          float64 `json:"-"`
	BaseRetrenchmentRateAdjusted                  float64 `json:"-"`
	BaseDisabilityIncidenceRate                   float64 `json:"-"`
	BaseDisabilityIncidenceRateAdjusted           float64 `json:"-"`
	MainMemberMortalityRateByMonth                float64 `json:"-"`
	MainMemberMortalityRateAdjustedByMonth        float64 `json:"-"`
	IndependentMortalityRateMonthly               float64 `json:"-"`
	IndependentMortalityRateAdjustedByMonth       float64 `json:"-"`
	IndependentLapseMonthly                       float64 `json:"-"`
	IndependentLapseMonthlyAdjusted               float64 `json:"-"`
	IndependentRetrenchmentMonthly                float64 `json:"-"`
	IndependentRetrenchmentMonthlyAdjusted        float64 `json:"-"`
	IndependentDisabilityMonthly                  float64 `json:"-"`
	IndependentDisabilityMonthlyAdjusted          float64 `json:"-"`
	MonthlyDependentMortality                     float64 `json:"-"`
	MonthlyDependentMortalityAdjusted             float64 `json:"-"`
	MonthlyDependentLapse                         float64 `json:"-"`
	MonthlyDependentLapseAdjusted                 float64 `json:"-"`
	MonthlyDependentRetrenchment                  float64 `json:"-"`
	MonthlyDependentRetrenchmentAdjusted          float64 `json:"-"`
	MonthlyDependentDisability                    float64 `json:"-"`
	MonthlyDependentDisabilityAdjusted            float64 `json:"-"`
	InitialPolicy                                 float64 `json:"-"`
	InitialPolicyAdjusted                         float64 `json:"-"`
	InitialPaidUp                                 float64 `json:"-"`
	InitialPaidUpAdjusted                         float64 `json:"-"`
	InitialPremiumWaivers                         float64 `json:"-"`
	InitialPremiumWaiversAdjusted                 float64 `json:"-"`
	InitialTemporaryPremiumWaivers                float64 `json:"-"`
	InitialTemporaryPremiumWaiversAdjusted        float64 `json:"-"`
	NaturalDeathsInForce                          float64 `json:"-"`
	NaturalDeathsInForceAdjusted                  float64 `json:"-"`
	NaturalDeathsPaidUp                           float64 `json:"-"`
	NaturalDeathsPaidUpAdjusted                   float64 `json:"-"`
	NaturalDeathsPremiumWaiver                    float64 `json:"-"`
	NaturalDeathsPremiumWaiverAdjusted            float64 `json:"-"`
	NaturalDeathsTemporaryWaivers                 float64 `json:"-"`
	NaturalDeathsTemporaryWaiversAdjusted         float64 `json:"-"`
	NumberOfAccidentDeaths                        float64 `json:"-"`
	NumberOfAccidentDeathsAdjusted                float64 `json:"-"`
	AccidentDeathsPaidUp                          float64 `json:"-"`
	AccidentDeathsPaidUpAdjusted                  float64 `json:"-"`
	AccidentDeathsPremiumWaiver                   float64 `json:"-"`
	AccidentDeathsPremiumWaiverAdjusted           float64 `json:"-"`
	AccidentDeathsTemporaryPremiumWaiver          float64 `json:"-"`
	AccidentDeathsTemporaryPremiumWaiverAdjusted  float64 `json:"-"`
	NumberOfLapses                                float64 `json:"-"`
	NumberOfLapsesAdjusted                        float64 `json:"-"`
	NumberOfDisabilities                          float64 `json:"-"`
	NumberOfDisabilitiesAdjusted                  float64 `json:"-"`
	NumberOfRetrenchments                         float64 `json:"-"`
	NumberOfRetrenchmentsAdjusted                 float64 `json:"-"`
	NumberOfMaturities                            float64 `json:"-"`
	NumberOfMaturitiesAdjusted                    float64 `json:"-"`
	TotalIncrementalNaturalDeaths                 float64 `json:"-"`
	TotalIncrementalNaturalDeathsAdjusted         float64 `json:"-"`
	TotalIncrementalAccidentalDeaths              float64 `json:"-"`
	TotalIncrementalAccidentalDeathsAdjusted      float64 `json:"-"`
	TotalIncrementalLapses                        float64 `json:"-"`
	TotalIncrementalLapsesAdjusted                float64 `json:"-"`
	TotalIncrementalDisabilities                  float64 `json:"-"`
	TotalIncrementalDisabilitiesAdjusted          float64 `json:"-"`
	TotalIncrementalRetrenchments                 float64 `json:"-"`
	TotalIncrementalRetrenchmentsAdjusted         float64 `json:"-"`
	SumAssured                                    float64 `json:"-"`
	CalculatedInstalment                          float64 `json:"-"`
	OutstandingSumAssured                         float64 `json:"-"`
	StandardAdditionalLumpSum                     float64 `json:"-"`
	RiderSumAssured                               float64 `json:"-"`
	AnnuityIncome                                 float64 `json:"-"`
	Premium                                       float64 `json:"-"`
	AllocatedPremium                              float64 `json:"-"`
	PolicyFee                                     float64 `json:"-"`
	PremiumAdvisoryFee                            float64 `json:"-"`
	FundAdvisoryFee                               float64 `json:"-"`
	UnfundedUnitFundSom                           float64 `json:"-"`
	FundInvestmentIncome                          float64 `json:"-"`
	ReversionaryBonus                             float64 `json:"-"`
	TerminalBonus                                 float64 `json:"-"`
	FundAssetManagementCharge                     float64 `json:"-"`
	FundRiskCharge                                float64 `json:"-"`
	UnfundedUnitFundEom                           float64 `json:"-"`
	UnitGrowthRiskMargin                          float64 `json:"-"`
	BonusStabilisationAccount                     float64 `json:"-"`
	SurrenderPenalty                              float64 `json:"-"`
	SurrenderValue                                float64 `json:"-"`
	MaturityValue                                 float64 `json:"-"`
	PremiumIncomeAdjusted                         float64 `json:"-"`
	EAllocatedPremiumIncome                       float64 `json:"-"`
	EAllocatedPremiumIncomeAdjusted               float64 `json:"-"`
	EPolicyFee                                    float64 `json:"-"`
	EPolicyFeeAdjusted                            float64 `json:"-"`
	EPremiumAdvisoryFee                           float64 `json:"-"`
	EPremiumAdvisoryFeeAdjusted                   float64 `json:"-"`
	EFundAdvisoryFee                              float64 `json:"-"`
	EFundAdvisoryFeeAdjusted                      float64 `json:"-"`
	EFundInvestmentIncome                         float64 `json:"-"`
	EFundInvestmentIncomeAdjusted                 float64 `json:"-"`
	EFundAssetManagementCharge                    float64 `json:"-"`
	EFundAssetManagementChargeAdjusted            float64 `json:"-"`
	EFundRiskCharge                               float64 `json:"-"`
	EFundRiskChargeAdjusted                       float64 `json:"-"`
	EBsaShareholderCharge                         float64 `json:"-"`
	EBsaShareholderChargeAdjusted                 float64 `json:"-"`
	PremiumNotReceived                            float64 `json:"-"`
	PremiumNotReceivedAdjusted                    float64 `json:"-"`
	CommissionAdjusted                            float64 `json:"-"`
	RenewalCommissionAdjusted                     float64 `json:"-"`
	ClawBackAdjusted                              float64 `json:"-"`
	SurrenderOutgo                                float64 `json:"-"`
	SurrenderOutgoAdjusted                        float64 `json:"-"`
	ESurrenderPenaltyCharge                       float64 `json:"-"`
	ESurrenderPenaltyChargeAdjusted               float64 `json:"-"`
	EPremiumAdvisoryCost                          float64 `json:"-"`
	EPremiumAdvisoryCostAdjusted                  float64 `json:"-"`
	EFundAdvisoryCost                             float64 `json:"-"`
	EFundAdvisoryCostAdjusted                     float64 `json:"-"`
	EGuaranteeCost                                float64 `json:"-"`
	EGuaranteeCostAdjusted                        float64 `json:"-"`
	DeathOutgoAdjusted                            float64 `json:"-"`
	NonLifeClaimsOutgoAdjusted                    float64 `json:"-"`
	AccidentalDeathOutgoAdjusted                  float64 `json:"-"`
	CashBackOnSurvivalAdjusted                    float64 `json:"-"`
	CashBackOnDeathAdjusted                       float64 `json:"-"`
	DisabilityOutgoAdjusted                       float64 `json:"-"`
	RetrenchmentOutgoAdjusted                     float64 `json:"-"`
	AnnuityOutgo                                  float64 `json:"-"`
	AnnuityOutgoAdjusted                          float64 `json:"-"`
	RiderAdjusted                                 float64 `json:"-"`
	InitialExpensesAdjusted                       float64 `json:"-"`
	RenewalExpensesAdjusted                       float64 `json:"-"`
	MaturityOutgo                                 float64 `json:"-"`
	MaturityOutgoAdjusted                         float64 `json:"-"`
	NetCashFlowAdjusted                           float64 `json:"-"`
	ValuationRate                                 float64 `json:"-"`
	ValuationRateAdjusted                         float64 `json:"-"`
	ReservesAdjusted                              float64 `json:"-"`
	ChangeInReservesAdjusted                      float64 `json:"-"`
	InvestmentIncomeAdjusted                      float64 `json:"-"`
	ProfitAdjustment                              float64 `json:"-"`
	ProfitAdjustmentAdjusted                      float64 `json:"-"`
	ProfitAdjusted                                float64 `json:"-"`
	RiskDiscountRate                              float64 `json:"-"`
	RiskDiscountRateAdjusted                      float64 `json:"-"`
	VIF                                           float64 `json:"-"`
	VIFAdjusted                                   float64 `json:"-"`
	CorporateTax                                  float64 `json:"-"`
	CorporateTaxAdjusted                          float64 `json:"-"`
	DiscountedPremiumIncomeAdjusted               float64 `json:"-"`
	DiscountedPremiumNotReceivedAdjusted          float64 `json:"-"`
	DiscountedUnallocatedPremiumIncome            float64 `json:"-"`
	DiscountedUnallocatedPremiumIncomeAdjusted    float64 `json:"-"`
	DiscountedPolicyFee                           float64 `json:"-"`
	DiscountedPolicyFeeAdjusted                   float64 `json:"-"`
	DiscountedPremiumAdvisoryFee                  float64 `json:"-"`
	DiscountedPremiumAdvisoryFeeAdjusted          float64 `json:"-"`
	DiscountedFundAdvisoryFee                     float64 `json:"-"`
	DiscountedFundAdvisoryFeeAdjusted             float64 `json:"-"`
	DiscountedFundAssetManagementCharges          float64 `json:"-"`
	DiscountedFundAssetManagementChargesAdjusted  float64 `json:"-"`
	DiscountedFundRiskCharge                      float64 `json:"-"`
	DiscountedFundRiskChargeAdjusted              float64 `json:"-"`
	DiscountedSurrenderPenaltyCharge              float64 `json:"-"`
	DiscountedSurrenderPenaltyChargeAdjusted      float64 `json:"-"`
	DiscountedBsaShareholderCharge                float64 `json:"-"`
	DiscountedBsaShareholderChargeAdjusted        float64 `json:"-"`
	DiscountedPremiumAdvisoryCost                 float64 `json:"-"`
	DiscountedPremiumAdvisoryCostAdjusted         float64 `json:"-"`
	DiscountedFundAdvisoryCost                    float64 `json:"-"`
	DiscountedFundAdvisoryCostAdjusted            float64 `json:"-"`
	DiscountedGuaranteeCost                       float64 `json:"-"`
	DiscountedGuaranteeCostAdjusted               float64 `json:"-"`
	DiscountedSurrenderOutgoAdjusted              float64 `json:"-"`
	DiscountedMaturityOutgo                       float64 `json:"-"`
	DiscountedMaturityOutgoAdjusted               float64 `json:"-"`
	DiscountedCommissionAdjusted                  float64 `json:"-"`
	DiscountedRenewalCommissionAdjusted           float64 `json:"-"`
	DiscountedClawBackAdjusted                    float64 `json:"-"`
	DiscountedDeathOutgoAdjusted                  float64 `json:"-"`
	DiscountedDisabilityOutgo                     float64 `json:"-"`
	DiscountedDisabilityOutgoAdjusted             float64 `json:"-"`
	DiscountedRetrenchmentOutgo                   float64 `json:"-"`
	DiscountedRetrenchmentOutgoAdjusted           float64 `json:"-"`
	DiscountedAnnuityOutgoAdjusted                float64 `json:"-"`
	DiscountedAccidentalDeathOutgoAdjusted        float64 `json:"-"`
	DiscountedCashBackOnSurvivalAdjusted          float64 `json:"-"`
	DiscountedCashBackOnDeathAdjusted             float64 `json:"-"`
	DiscountedRiderAdjusted                       float64 `json:"-"`
	DiscountedInitialExpensesAdjusted             float64 `json:"-"`
	DiscountedRenewalExpensesAdjusted             float64 `json:"-"`
	DiscountedInvestmentIncome                    float64 `json:"-"`
	DiscountedInvestmentIncomeAdjusted            float64 `json:"-"`
	DiscountedProfitAdjustment                    float64 `json:"-"`
	DiscountedProfitAdjustmentAdjusted            float64 `json:"-"`
	DiscountedProfitAdjusted                      float64 `json:"-"`
	DiscountedCorporateTax                        float64 `json:"-"`
	DiscountedCorporateTaxAdjusted                float64 `json:"-"`
	ReinsurancePremiumAdjusted                    float64 `json:"-"`
	ReinsuranceCedingCommissionAdjusted           float64 `json:"-"`
	ReinsuranceClaimsAdjusted                     float64 `json:"-"`
	DiscountedReinsurancePremiumAdjusted          float64 `json:"-"`
	DiscountedReinsuranceCedingCommissionAdjusted float64 `json:"-"`
	DiscountedReinsuranceClaimsAdjusted           float64 `json:"-"`
	NetReinsuranceCashflowAdjusted                float64 `json:"-"`
	DiscountedNetReinsuranceCashflowAdjusted      float64 `json:"-"`
}

type ManualSapList struct {
	//ID          int    `json:"id"`
	RunName   string `json:"run_name" csv:"run_name"`
	RunDate   string `json:"run_date" csv:"run_date"`
	CreatedAt string `json:"created_at" csv:"created_at" gorm:"index"`
	User      string `json:"user" csv:"user"`
	UserEmail string `json:"user_email" csv:"user_email"`
}

type ManualScopedAggregatedProjection struct {
	ID                      int `json:"id" gorm:"primary_key"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
	RunName                 string  `json:"run_name" csv:"run_name"`
	RunDate                 string  `json:"run_date" csv:"run_date"`
	ProductCode             string  `json:"product_code" csv:"product_code" gorm:"index"`
	SpCode                  int     `json:"sp_code" csv:"sp_code"`
	IFRS17Group             string  `json:"ifrs17_group" csv:"ifrs17_group"`
	ProjectionMonth         int     `json:"projection_month" gorm:"column:projection_month" csv:"projection_month"`
	PremiumIncome           float64 `json:"premium_income" csv:"premium_income"`
	PremiumNotReceivedLapse float64 `json:"premium_not_received_lapse" csv:"premium_not_received_lapse"`
	Commission              float64 `json:"commission" csv:"commission"`
	RenewalCommission       float64 `json:"renewal_commission" csv:"renewal_commission"`
	ClawBack                float64 `json:"claw_back" csv:"claw_back"`
	NonLifeClaimsOutgo      float64 `json:"non_life_claims_outgo" csv:"non_life_claims_outgo"`
	DeathOutgo              float64 `json:"death_outgo" csv:"death_outgo"`
	AccidentalDeathOutgo    float64 `json:"accidental_death_outgo" csv:"accidental_death_outgo"`
	CashBackOnSurvival      float64 `json:"cash_back_on_survival" csv:"cash_back_on_survival"`
	CashBackOnDeath         float64 `json:"cash_back_on_death" csv:"cash_back_on_death"`
	DisabilityOutgo         float64 `json:"disability_outgo" csv:"disability_outgo"`
	RetrenchmentOutgo       float64 `json:"retrenchment_outgo" csv:"retrenchment_outgo"`
	Rider                   float64 `json:"rider" csv:"rider"`
	//RiderAdjusted                         float64 `json:"rider_adjusted" csv:"rider_adjusted"`
	InitialExpenses                       float64 `json:"initial_expenses" csv:"initial_expenses"`
	RenewalExpenses                       float64 `json:"renewal_expenses" csv:"renewal_expenses"`
	NetCashFlow                           float64 `json:"net_cash_flow" csv:"net_cash_flow"`
	Reserves                              float64 `json:"reserves" csv:"reserves"`
	ChangeInReserves                      float64 `json:"change_in_reserves" csv:"change_in_reserves"`
	InvestmentIncome                      float64 `json:"investment_income" csv:"investment_income"`
	Profit                                float64 `json:"profit" csv:"profit"`
	CoverageUnits                         float64 `json:"coverage_units" csv:"coverage_units"`
	DiscountedPremiumIncome               float64 `json:"discounted_premium_income" csv:"discounted_premium_income"`
	DiscountedPremiumNotReceived          float64 `json:"discounted_premium_not_received" csv:"discounted_premium_not_received"`
	DiscountedCommission                  float64 `json:"discounted_commission" csv:"discounted_commission"`
	DiscountedClawBack                    float64 `json:"discounted_claw_back" csv:"discounted_claw_back"`
	DiscountedDeathOutgo                  float64 `json:"discounted_death_outgo" csv:"discounted_death_outgo"`
	DiscountedAccidentalDeathOutgo        float64 `json:"discounted_accidental_death_outgo" csv:"discounted_accidental_death_outgo"`
	DiscountedMorbidityOutgo              float64 `json:"discounted_morbidity_outgo" csv:"discounted_morbidity_outgo"`
	DiscountedOutgo                       float64 `json:"discounted_outgo" csv:"discounted_outgo"`
	DiscountedCashBackOnSurvival          float64 `json:"discounted_cash_back_on_survival" csv:"discounted_cash_back_on_survival"`
	DiscountedCashBackOnDeath             float64 `json:"discounted_cash_back_on_death" csv:"discounted_cash_back_on_death"`
	DiscountedRider                       float64 `json:"discounted_rider" csv:"discounted_rider"`
	DiscountedInitialExpenses             float64 `json:"discounted_initial_expenses" csv:"discounted_initial_expenses"`
	DiscountedRenewalExpenses             float64 `json:"discounted_renewal_expenses" csv:"discounted_renewal_expenses"`
	DiscountedProfit                      float64 `json:"discounted_profit" csv:"discounted_profit"`
	DiscountedAnnuityFactor               float64 `json:"annuity_factor" csv:"annuity_factor"`
	DiscountedCashOutflow                 float64 `json:"discounted_cash_outflow" csv:"discounted_cash_outflow"`
	DiscountedCashOutflowExclAcquisition  float64 `json:"discounted_cash_outflow_excl_acquisition" csv:"discounted_cash_outflow_excl_acquisition"`
	DiscountedAcquisitionCost             float64 `json:"discounted_acquisition_cost" csv:"discounted_acquisition_cost"`
	DiscountedCashInflow                  float64 `json:"discounted_cash_inflow" csv:"discounted_cash_inflow"`
	SumCoverageUnits                      float64 `json:"sum_coverage_units" csv:"sum_coverage_units"`
	DiscountedCoverageUnits               float64 `json:"discounted_coverage_units" csv:"discounted_coverage_units"`
	CededSumAssured                       float64 `json:"ceded_sum_assured" csv:"ceded_sum_assured"`
	ReinsurancePremium                    float64 `json:"reinsurance_premium" csv:"reinsurance_premium"`
	ReinsuranceCedingCommission           float64 `json:"reinsurance_ceding_commission" csv:"reinsurance_ceding_commission"`
	ReinsuranceClaims                     float64 `json:"reinsurance_claims" csv:"reinsurance_claims"`
	NetReinsuranceCashflow                float64 `json:"net_reinsurance_cashflow" csv:"net_reinsurance_cashflow"`
	DiscountedNetReinsuranceCashflow      float64 `json:"discounted_net_reinsurance_cashflow" csv:"discounted_net_reinsurance_cashflow"`
	DiscountedReinsurancePremium          float64 `json:"discounted_reinsurance_premium" csv:"discounted_reinsurance_premium"`
	DiscountedReinsuranceCedingCommission float64 `json:"discounted_reinsurance_ceding_commission" csv:"discounted_reinsurance_ceding_commission"`
	DiscountedReinsuranceClaims           float64 `json:"discounted_reinsurance_claims" csv:"discounted_reinsurance_claims"`
	RunBasis                              string  `json:"run_basis" csv:"run_basis"`
	// Manual Import SAP
	YieldCurveCode                          string  `json:"yield_curve_code" csv:"yield_curve_code"`
	DiscountedCashOutflowM0                 float64 `json:"discounted_cash_outflow_m0" csv:"discounted_cash_outflow_m0"`
	DiscountedCashOutflowM12                float64 `json:"discounted_cash_outflow_m12" csv:"discounted_cash_outflow_m12"`
	DiscountedCashOutflowExclAcquisitionM0  float64 `json:"discounted_cash_outflow_excl_acquisition_m0" csv:"discounted_cash_outflow_excl_acquisition_m0"`
	DiscountedCashOutflowExclAcquisitionM12 float64 `json:"discounted_cash_outflow_excl_acquisition_m12" csv:"discounted_cash_outflow_excl_acquisition_m12"`
	DiscountedAcquisitionCostM0             float64 `json:"discounted_acquisition_cost_m0" csv:"discounted_acquisition_cost_m0"`
	DiscountedAcquisitionCostM12            float64 `json:"discounted_acquisition_cost_m12" csv:"discounted_acquisition_cost_m12"`
	DiscountedCashInflowM0                  float64 `json:"discounted_cash_inflow_m0" csv:"discounted_cash_inflow_m0"`
	DiscountedCashInflowM12                 float64 `json:"discounted_cash_inflow_m12" csv:"discounted_cash_inflow_m12"`
	SumCoverageUnitsM0                      float64 `json:"sum_coverage_units_m0" csv:"sum_coverage_units_m0"`
	SumCoverageUnitsM12                     float64 `json:"sum_coverage_units_m12" csv:"sum_coverage_units_m12"`
	SumCoverageUnitsM24                     float64 `json:"sum_coverage_units_m24" csv:"sum_coverage_units_m24"`
	SumCoverageUnitsM36                     float64 `json:"sum_coverage_units_m36" csv:"sum_coverage_units_m36"`
	SumCoverageUnitsM48                     float64 `json:"sum_coverage_units_m48" csv:"sum_coverage_units_m48"`
	SumCoverageUnitsM60                     float64 `json:"sum_coverage_units_m60" csv:"sum_coverage_units_m60"`
	SumCoverageUnitsM72                     float64 `json:"sum_coverage_units_m72" csv:"sum_coverage_units_m72"`
	SumCoverageUnitsM84                     float64 `json:"sum_coverage_units_m84" csv:"sum_coverage_units_m84"`
	SumCoverageUnitsM96                     float64 `json:"sum_coverage_units_m96" csv:"sum_coverage_units_m96"`
	SumCoverageUnitsM108                    float64 `json:"sum_coverage_units_m108" csv:"sum_coverage_units_m108"`
	SumCoverageUnitsM120                    float64 `json:"sum_coverage_units_m120" csv:"sum_coverage_units_m120"`
	DiscountedCoverageUnitsM0               float64 `json:"discounted_coverage_units_m0" csv:"discounted_coverage_units_m0"`
	DiscountedCoverageUnitsM12              float64 `json:"discounted_coverage_units_m12" csv:"discounted_coverage_units_m12"`
	DiscountedCoverageUnitsM24              float64 `json:"discounted_coverage_units_m24" csv:"discounted_coverage_units_m24"`
	DiscountedCoverageUnitsM36              float64 `json:"discounted_coverage_units_m36" csv:"discounted_coverage_units_m36"`
	DiscountedCoverageUnitsM48              float64 `json:"discounted_coverage_units_m48" csv:"discounted_coverage_units_m48"`
	DiscountedCoverageUnitsM60              float64 `json:"discounted_coverage_units_m60" csv:"discounted_coverage_units_m60"`
	DiscountedCoverageUnitsM72              float64 `json:"discounted_coverage_units_m72" csv:"discounted_coverage_units_m72"`
	DiscountedCoverageUnitsM84              float64 `json:"discounted_coverage_units_m84" csv:"discounted_coverage_units_m84"`
	DiscountedCoverageUnitsM96              float64 `json:"discounted_coverage_units_m96" csv:"discounted_coverage_units_m96"`
	DiscountedCoverageUnitsM108             float64 `json:"discounted_coverage_units_m108" csv:"discounted_coverage_units_m108"`
	DiscountedCoverageUnitsM120             float64 `json:"discounted_coverage_units_m120" csv:"discounted_coverage_units_m120"`
	BelM0                                   float64 `json:"bel_m0" csv:"bel_m0"`
	BelM12                                  float64 `json:"bel_m12" csv:"bel_m12"`
	RiskAdjustmentM0                        float64 `json:"risk_adjustment_m0" csv:"risk_adjustment_m0"`
	RiskAdjustmentM12                       float64 `json:"risk_adjustment_m12" csv:"risk_adjustment_m12"`
	CashInflowM1                            float64 `json:"cash_inflow_m1" csv:"cash_inflow_m1"`
	CashInflowM2                            float64 `json:"cash_inflow_m2" csv:"cash_inflow_m2"`
	CashInflowM3                            float64 `json:"cash_inflow_m3" csv:"cash_inflow_m3"`
	CashInflowM4                            float64 `json:"cash_inflow_m4" csv:"cash_inflow_m4"`
	CashInflowM5                            float64 `json:"cash_inflow_m5" csv:"cash_inflow_m5"`
	CashInflowM6                            float64 `json:"cash_inflow_m6" csv:"cash_inflow_m6"`
	CashInflowM7                            float64 `json:"cash_inflow_m7" csv:"cash_inflow_m7"`
	CashInflowM8                            float64 `json:"cash_inflow_m8" csv:"cash_inflow_m8"`
	CashInflowM9                            float64 `json:"cash_inflow_m9" csv:"cash_inflow_m9"`
	CashInflowM10                           float64 `json:"cash_inflow_m10" csv:"cash_inflow_m10"`
	CashInflowM11                           float64 `json:"cash_inflow_m11" csv:"cash_inflow_m11"`
	CashInflowM12                           float64 `json:"cash_inflow_m12" csv:"cash_inflow_m12"`
	MortalityOutgoM1                        float64 `json:"mortality_outgo_m1" csv:"mortality_outgo_m1"`
	MortalityOutgoM2                        float64 `json:"mortality_outgo_m2" csv:"mortality_outgo_m2"`
	MortalityOutgoM3                        float64 `json:"mortality_outgo_m3" csv:"mortality_outgo_m3"`
	MortalityOutgoM4                        float64 `json:"mortality_outgo_m4" csv:"mortality_outgo_m4"`
	MortalityOutgoM5                        float64 `json:"mortality_outgo_m5" csv:"mortality_outgo_m5"`
	MortalityOutgoM6                        float64 `json:"mortality_outgo_m6" csv:"mortality_outgo_m6"`
	MortalityOutgoM7                        float64 `json:"mortality_outgo_m7" csv:"mortality_outgo_m7"`
	MortalityOutgoM8                        float64 `json:"mortality_outgo_m8" csv:"mortality_outgo_m8"`
	MortalityOutgoM9                        float64 `json:"mortality_outgo_m9" csv:"mortality_outgo_m9"`
	MortalityOutgoM10                       float64 `json:"mortality_outgo_m10" csv:"mortality_outgo_m10"`
	MortalityOutgoM11                       float64 `json:"mortality_outgo_m11" csv:"mortality_outgo_m11"`
	MortalityOutgoM12                       float64 `json:"mortality_outgo_m12" csv:"mortality_outgo_m12"`
	RetrenchmentOutgoM1                     float64 `json:"retrenchment_outgo_m1" csv:"retrenchment_outgo_m1"`
	RetrenchmentOutgoM2                     float64 `json:"retrenchment_outgo_m2" csv:"retrenchment_outgo_m2"`
	RetrenchmentOutgoM3                     float64 `json:"retrenchment_outgo_m3" csv:"retrenchment_outgo_m3"`
	RetrenchmentOutgoM4                     float64 `json:"retrenchment_outgo_m4" csv:"retrenchment_outgo_m4"`
	RetrenchmentOutgoM5                     float64 `json:"retrenchment_outgo_m5" csv:"retrenchment_outgo_m5"`
	RetrenchmentOutgoM6                     float64 `json:"retrenchment_outgo_m6" csv:"retrenchment_outgo_m6"`
	RetrenchmentOutgoM7                     float64 `json:"retrenchment_outgo_m7" csv:"retrenchment_outgo_m7"`
	RetrenchmentOutgoM8                     float64 `json:"retrenchment_outgo_m8" csv:"retrenchment_outgo_m8"`
	RetrenchmentOutgoM9                     float64 `json:"retrenchment_outgo_m9" csv:"retrenchment_outgo_m9"`
	RetrenchmentOutgoM10                    float64 `json:"retrenchment_outgo_m10" csv:"retrenchment_outgo_m10"`
	RetrenchmentOutgoM11                    float64 `json:"retrenchment_outgo_m11" csv:"retrenchment_outgo_m11"`
	RetrenchmentOutgoM12                    float64 `json:"retrenchment_outgo_m12" csv:"retrenchment_outgo_m12"`
	MorbidityOutgoM1                        float64 `json:"morbidity_outgo_m1" csv:"morbidity_outgo_m1"`
	MorbidityOutgoM2                        float64 `json:"morbidity_outgo_m2" csv:"morbidity_outgo_m2"`
	MorbidityOutgoM3                        float64 `json:"morbidity_outgo_m3" csv:"morbidity_outgo_m3"`
	MorbidityOutgoM4                        float64 `json:"morbidity_outgo_m4" csv:"morbidity_outgo_m4"`
	MorbidityOutgoM5                        float64 `json:"morbidity_outgo_m5" csv:"morbidity_outgo_m5"`
	MorbidityOutgoM6                        float64 `json:"morbidity_outgo_m6" csv:"morbidity_outgo_m6"`
	MorbidityOutgoM7                        float64 `json:"morbidity_outgo_m7" csv:"morbidity_outgo_m7"`
	MorbidityOutgoM8                        float64 `json:"morbidity_outgo_m8" csv:"morbidity_outgo_m8"`
	MorbidityOutgoM9                        float64 `json:"morbidity_outgo_m9" csv:"morbidity_outgo_m9"`
	MorbidityOutgoM10                       float64 `json:"morbidity_outgo_m10" csv:"morbidity_outgo_m10"`
	MorbidityOutgoM11                       float64 `json:"morbidity_outgo_m11" csv:"morbidity_outgo_m11"`
	MorbidityOutgoM12                       float64 `json:"morbidity_outgo_m12" csv:"morbidity_outgo_m12"`
	NonLifeClaimsOutgoM1                    float64 `json:"non_life_claims_outgo_m1" csv:"non_life_claims_outgo_m1"`
	NonLifeClaimsOutgoM2                    float64 `json:"non_life_claims_outgo_m2" csv:"non_life_claims_outgo_m2"`
	NonLifeClaimsOutgoM3                    float64 `json:"non_life_claims_outgo_m3" csv:"non_life_claims_outgo_m3"`
	NonLifeClaimsOutgoM4                    float64 `json:"non_life_claims_outgo_m4" csv:"non_life_claims_outgo_m4"`
	NonLifeClaimsOutgoM5                    float64 `json:"non_life_claims_outgo_m5" csv:"non_life_claims_outgo_m5"`
	NonLifeClaimsOutgoM6                    float64 `json:"non_life_claims_outgo_m6" csv:"non_life_claims_outgo_m6"`
	NonLifeClaimsOutgoM7                    float64 `json:"non_life_claims_outgo_m7" csv:"non_life_claims_outgo_m7"`
	NonLifeClaimsOutgoM8                    float64 `json:"non_life_claims_outgo_m8" csv:"non_life_claims_outgo_m8"`
	NonLifeClaimsOutgoM9                    float64 `json:"non_life_claims_outgo_m9" csv:"non_life_claims_outgo_m9"`
	NonLifeClaimsOutgoM10                   float64 `json:"non_life_claims_outgo_m10" csv:"non_life_claims_outgo_m10"`
	NonLifeClaimsOutgoM11                   float64 `json:"non_life_claims_outgo_m11" csv:"non_life_claims_outgo_m11"`
	NonLifeClaimsOutgoM12                   float64 `json:"non_life_claims_outgo_m12" csv:"non_life_claims_outgo_m12"`
	RenewalExpenseOutgoM1                   float64 `json:"renewal_expense_outgo_m1" csv:"renewal_expense_outgo_m1"`
	RenewalExpenseOutgoM2                   float64 `json:"renewal_expense_outgo_m2" csv:"renewal_expense_outgo_m2"`
	RenewalExpenseOutgoM3                   float64 `json:"renewal_expense_outgo_m3" csv:"renewal_expense_outgo_m3"`
	RenewalExpenseOutgoM4                   float64 `json:"renewal_expense_outgo_m4" csv:"renewal_expense_outgo_m4"`
	RenewalExpenseOutgoM5                   float64 `json:"renewal_expense_outgo_m5" csv:"renewal_expense_outgo_m5"`
	RenewalExpenseOutgoM6                   float64 `json:"renewal_expense_outgo_m6" csv:"renewal_expense_outgo_m6"`
	RenewalExpenseOutgoM7                   float64 `json:"renewal_expense_outgo_m7" csv:"renewal_expense_outgo_m7"`
	RenewalExpenseOutgoM8                   float64 `json:"renewal_expense_outgo_m8" csv:"renewal_expense_outgo_m8"`
	RenewalExpenseOutgoM9                   float64 `json:"renewal_expense_outgo_m9" csv:"renewal_expense_outgo_m9"`
	RenewalExpenseOutgoM10                  float64 `json:"renewal_expense_outgo_m10" csv:"renewal_expense_outgo_m10"`
	RenewalExpenseOutgoM11                  float64 `json:"renewal_expense_outgo_m11" csv:"renewal_expense_outgo_m11"`
	RenewalExpenseOutgoM12                  float64 `json:"renewal_expense_outgo_m12" csv:"renewal_expense_outgo_m12"`
	DiscountedEntityShareM0                 float64 `json:"discounted_entity_share_m0" csv:"discounted_entity_share_m0"`
	DiscountedEntityShareM12                float64 `json:"discounted_entity_share_m12" csv:"discounted_entity_share_m12"`
	DiscountedCashInflowsNotVaryM0          float64 `json:"discounted_cash_inflows_not_vary_m0" csv:"discounted_cash_inflows_not_vary_m0"`
	DiscountedCashInflowsNotVaryM12         float64 `json:"discounted_cash_inflows_not_vary_m12" csv:"discounted_cash_inflows_not_vary_m12"`
	TvogM0                                  float64 `json:"tvog_m0" csv:"tvog_m0"`
	TvogM12                                 float64 `json:"tvog_m12" csv:"tvog_m12"`
	UnitFundM0                              float64 `json:"unit_fund_m0" csv:"unit_fund_m0"`
	EntityShareM1                           float64 `json:"entity_share_m1" csv:"entity_share_m1"`
	EntityShareM2                           float64 `json:"entity_share_m2" csv:"entity_share_m2"`
	EntityShareM3                           float64 `json:"entity_share_m3" csv:"entity_share_m3"`
	EntityShareM4                           float64 `json:"entity_share_m4" csv:"entity_share_m4"`
	EntityShareM5                           float64 `json:"entity_share_m5" csv:"entity_share_m5"`
	EntityShareM6                           float64 `json:"entity_share_m6" csv:"entity_share_m6"`
	EntityShareM7                           float64 `json:"entity_share_m7" csv:"entity_share_m7"`
	EntityShareM8                           float64 `json:"entity_share_m8" csv:"entity_share_m8"`
	EntityShareM9                           float64 `json:"entity_share_m9" csv:"entity_share_m9"`
	EntityShareM10                          float64 `json:"entity_share_m10" csv:"entity_share_m10"`
	EntityShareM11                          float64 `json:"entity_share_m11" csv:"entity_share_m11"`
	EntityShareM12                          float64 `json:"entity_share_m12" csv:"entity_share_m12"`
	CashInflowsNotVaryM1                    float64 `json:"cash_inflows_not_vary_m1" csv:"cash_inflows_not_vary_m1"`
	CashInflowsNotVaryM2                    float64 `json:"cash_inflows_not_vary_m2" csv:"cash_inflows_not_vary_m2"`
	CashInflowsNotVaryM3                    float64 `json:"cash_inflows_not_vary_m3" csv:"cash_inflows_not_vary_m3"`
	CashInflowsNotVaryM4                    float64 `json:"cash_inflows_not_vary_m4" csv:"cash_inflows_not_vary_m4"`
	CashInflowsNotVaryM5                    float64 `json:"cash_inflows_not_vary_m5" csv:"cash_inflows_not_vary_m5"`
	CashInflowsNotVaryM6                    float64 `json:"cash_inflows_not_vary_m6" csv:"cash_inflows_not_vary_m6"`
	CashInflowsNotVaryM7                    float64 `json:"cash_inflows_not_vary_m7" csv:"cash_inflows_not_vary_m7"`
	CashInflowsNotVaryM8                    float64 `json:"cash_inflows_not_vary_m8" csv:"cash_inflows_not_vary_m8"`
	CashInflowsNotVaryM9                    float64 `json:"cash_inflows_not_vary_m9" csv:"cash_inflows_not_vary_m9"`
	CashInflowsNotVaryM10                   float64 `json:"cash_inflows_not_vary_m10" csv:"cash_inflows_not_vary_m10"`
	CashInflowsNotVaryM11                   float64 `json:"cash_inflows_not_vary_m11" csv:"cash_inflows_not_vary_m11"`
	CashInflowsNotVaryM12                   float64 `json:"cash_inflows_not_vary_m12" csv:"cash_inflows_not_vary_m12"`
	OtherCashOutflowsM1                     float64 `json:"other_cash_outflows_m1" csv:"other_cash_outflows_m1"`
	OtherCashOutflowsM2                     float64 `json:"other_cash_outflows_m2" csv:"other_cash_outflows_m2"`
	OtherCashOutflowsM3                     float64 `json:"other_cash_outflows_m3" csv:"other_cash_outflows_m3"`
	OtherCashOutflowsM4                     float64 `json:"other_cash_outflows_m4" csv:"other_cash_outflows_m4"`
	OtherCashOutflowsM5                     float64 `json:"other_cash_outflows_m5" csv:"other_cash_outflows_m5"`
	OtherCashOutflowsM6                     float64 `json:"other_cash_outflows_m6" csv:"other_cash_outflows_m6"`
	OtherCashOutflowsM7                     float64 `json:"other_cash_outflows_m7" csv:"other_cash_outflows_m7"`
	OtherCashOutflowsM8                     float64 `json:"other_cash_outflows_8" csv:"other_cash_outflows_8"`
	OtherCashOutflowsM9                     float64 `json:"other_cash_outflows_m9" csv:"other_cash_outflows_m9"`
	OtherCashOutflowsM10                    float64 `json:"other_cash_outflows_m10" csv:"other_cash_outflows_m10"`
	OtherCashOutflowsM11                    float64 `json:"other_cash_outflows_m11" csv:"other_cash_outflows_m11"`
	OtherCashOutflowsM12                    float64 `json:"other_cash_outflows_m12" csv:"other_cash_outflows_m12"`
	Treaty1BelM0                            float64 `json:"treaty1_bel_m0" csv:"treaty1_bel_m0"`
	Treaty1BelM12                           float64 `json:"treaty1_bel_m12" csv:"treaty1_bel_m12"`
	Treaty1RiskAdjustmentM0                 float64 `json:"treaty1_risk_adjustment_m0" csv:"treaty1_risk_adjustment_m0"`
	Treaty1RiskAdjustmentM12                float64 `json:"treaty1_risk_adjustment_m12" csv:"treaty1_risk_adjustment_m12"`
	Treaty1DiscountedCashOutflowM0          float64 `json:"treaty1_discounted_cash_outflow_m0" csv:"treaty1_discounted_cash_outflow_m0"`
	Treaty1DiscountedCashOutflowM12         float64 `json:"treaty1_discounted_cash_outflow_m12" csv:"treaty1_discounted_cash_outflow_m12"`
	Treaty1DiscountedCashInflowM0           float64 `json:"treaty1_discounted_cash_inflow_m0" csv:"treaty1_discounted_cash_inflow_m0"`
	Treaty1DiscountedCashInflowM12          float64 `json:"treaty1_discounted_cash_inflow_m12" csv:"treaty1_discounted_cash_inflow_m12"`
	Treaty2BelM0                            float64 `json:"treaty2_bel_m0" csv:"treaty2_bel_m0"`
	Treaty2BelM12                           float64 `json:"treaty2_bel_m12" csv:"treaty2_bel_m12"`
	Treaty2RiskAdjustmentM0                 float64 `json:"treaty2_risk_adjustment_m0" csv:"treaty2_risk_adjustment_m0"`
	Treaty2RiskAdjustmentM12                float64 `json:"treaty2_risk_adjustment_m12" csv:"treaty2_risk_adjustment_m12"`
	Treaty2DiscountedCashOutflowM0          float64 `json:"treaty2_discounted_cash_outflow_m0" csv:"treaty2_discounted_cash_outflow_m0"`
	Treaty2DiscountedCashOutflowM12         float64 `json:"treaty2_discounted_cash_outflow_m12" csv:"treaty2_discounted_cash_outflow_m12"`
	Treaty2DiscountedCashInflowM0           float64 `json:"treaty2_discounted_cash_inflow_m0" csv:"treaty2_discounted_cash_inflow_m0"`
	Treaty2DiscountedCashInflowM12          float64 `json:"treaty2_discounted_cash_inflow_m12" csv:"treaty2_discounted_cash_inflow_m12"`
	Treaty3BelM0                            float64 `json:"treaty3_bel_m0" csv:"treaty3_bel_m0"`
	Treaty3BelM12                           float64 `json:"treaty3_bel_m12" csv:"treaty3_bel_m12"`
	Treaty3RiskAdjustmentM0                 float64 `json:"treaty3_risk_adjustment_m0" csv:"treaty3_risk_adjustment_m0"`
	Treaty3RiskAdjustmentM12                float64 `json:"treaty3_risk_adjustment_m12" csv:"treaty3_risk_adjustment_m12"`
	Treaty3DiscountedCashOutflowM0          float64 `json:"treaty3_discounted_cash_outflow_m0" csv:"treaty3_discounted_cash_outflow_m0"`
	Treaty3DiscountedCashOutflowM12         float64 `json:"treaty3_discounted_cash_outflow_m12" csv:"treaty3_discounted_cash_outflow_m12"`
	Treaty3DiscountedCashInflowM0           float64 `json:"treaty3_discounted_cash_inflow_m0" csv:"treaty3_discounted_cash_inflow_m0"`
	Treaty3DiscountedCashInflowM12          float64 `json:"treaty3_discounted_cash_inflow_m12" csv:"treaty3_discounted_cash_inflow_m12"`
	Treaty1CashInflowM1                     float64 `json:"treaty1_cash_inflow_m1" csv:"treaty1_cash_inflow_m1"`
	Treaty1CashInflowM2                     float64 `json:"treaty1_cash_inflow_m2" csv:"treaty1_cash_inflow_m2"`
	Treaty1CashInflowM3                     float64 `json:"treaty1_cash_inflow_m3" csv:"treaty1_cash_inflow_m3"`
	Treaty1CashInflowM4                     float64 `json:"treaty1_cash_inflow_m4" csv:"treaty1_cash_inflow_m4"`
	Treaty1CashInflowM5                     float64 `json:"treaty1_cash_inflow_m5" csv:"treaty1_cash_inflow_m5"`
	Treaty1CashInflowM6                     float64 `json:"treaty1_cash_inflow_m6" csv:"treaty1_cash_inflow_m6"`
	Treaty1CashInflowM7                     float64 `json:"treaty1_cash_inflow_m7" csv:"treaty1_cash_inflow_m7"`
	Treaty1CashInflowM8                     float64 `json:"treaty1_cash_inflow_m8" csv:"treaty1_cash_inflow_m8"`
	Treaty1CashInflowM9                     float64 `json:"treaty1_cash_inflow_m9" csv:"treaty1_cash_inflow_m9"`
	Treaty1CashInflowM10                    float64 `json:"treaty1_cash_inflow_m10" csv:"treaty1_cash_inflow_m10"`
	Treaty1CashInflowM11                    float64 `json:"treaty1_cash_inflow_m11" csv:"treaty1_cash_inflow_m11"`
	Treaty1CashInflowM12                    float64 `json:"treaty1_cash_inflow_m12" csv:"treaty1_cash_inflow_m12"`
	Treaty1CashOutflowM1                    float64 `json:"treaty1_cash_outflow_m1" csv:"treaty1_cash_outflow_m1"`
	Treaty1CashOutflowM2                    float64 `json:"treaty1_cash_outflow_m2" csv:"treaty1_cash_outflow_m2"`
	Treaty1CashOutflowM3                    float64 `json:"treaty1_cash_outflow_m3" csv:"treaty1_cash_outflow_m3"`
	Treaty1CashOutflowM4                    float64 `json:"treaty1_cash_outflow_m4" csv:"treaty1_cash_outflow_m4"`
	Treaty1CashOutflowM5                    float64 `json:"treaty1_cash_outflow_m5" csv:"treaty1_cash_outflow_m5"`
	Treaty1CashOutflowM6                    float64 `json:"treaty1_cash_outflow_m6" csv:"treaty1_cash_outflow_m6"`
	Treaty1CashOutflowM7                    float64 `json:"treaty1_cash_outflow_m7" csv:"treaty1_cash_outflow_m7"`
	Treaty1CashOutflowM8                    float64 `json:"treaty1_cash_outflow_m8" csv:"treaty1_cash_outflow_m8"`
	Treaty1CashOutflowM9                    float64 `json:"treaty1_cash_outflow_m9" csv:"treaty1_cash_outflow_m9"`
	Treaty1CashOutflowM10                   float64 `json:"treaty1_cash_outflow_m10" csv:"treaty1_cash_outflow_m10"`
	Treaty1CashOutflowM11                   float64 `json:"treaty1_cash_outflow_m11" csv:"treaty1_cash_outflow_m11"`
	Treaty1CashOutflowM12                   float64 `json:"treaty1_cash_outflow_m12" csv:"treaty1_cash_outflow_m12"`
	Treaty2CashInflowM1                     float64 `json:"treaty2_cash_inflow_m1" csv:"treaty2_cash_inflow_m1"`
	Treaty2CashInflowM2                     float64 `json:"treaty2_cash_inflow_m2" csv:"treaty2_cash_inflow_m2"`
	Treaty2CashInflowM3                     float64 `json:"treaty2_cash_inflow_m3" csv:"treaty2_cash_inflow_m3"`
	Treaty2CashInflowM4                     float64 `json:"treaty2_cash_inflow_m4" csv:"treaty2_cash_inflow_m4"`
	Treaty2CashInflowM5                     float64 `json:"treaty2_cash_inflow_m5" csv:"treaty2_cash_inflow_m5"`
	Treaty2CashInflowM6                     float64 `json:"treaty2_cash_inflow_m6" csv:"treaty2_cash_inflow_m6"`
	Treaty2CashInflowM7                     float64 `json:"treaty2_cash_inflow_m7" csv:"treaty2_cash_inflow_m7"`
	Treaty2CashInflowM8                     float64 `json:"treaty2_cash_inflow_m8" csv:"treaty2_cash_inflow_m8"`
	Treaty2CashInflowM9                     float64 `json:"treaty2_cash_inflow_m9" csv:"treaty2_cash_inflow_m9"`
	Treaty2CashInflowM10                    float64 `json:"treaty2_cash_inflow_m10" csv:"treaty2_cash_inflow_m10"`
	Treaty2CashInflowM11                    float64 `json:"treaty2_cash_inflow_m11" csv:"treaty2_cash_inflow_m11"`
	Treaty2CashInflowM12                    float64 `json:"treaty2_cash_inflow_m12" csv:"treaty2_cash_inflow_m12"`
	Treaty2CashOutflowM1                    float64 `json:"treaty2_cash_outflow_m1" csv:"treaty2_cash_outflow_m1"`
	Treaty2CashOutflowM2                    float64 `json:"treaty2_cash_outflow_m2" csv:"treaty2_cash_outflow_m2"`
	Treaty2CashOutflowM3                    float64 `json:"treaty2_cash_outflow_m3" csv:"treaty2_cash_outflow_m3"`
	Treaty2CashOutflowM4                    float64 `json:"treaty2_cash_outflow_m4" csv:"treaty2_cash_outflow_m4"`
	Treaty2CashOutflowM5                    float64 `json:"treaty2_cash_outflow_m5" csv:"treaty2_cash_outflow_m5"`
	Treaty2CashOutflowM6                    float64 `json:"treaty2_cash_outflow_m6" csv:"treaty2_cash_outflow_m6"`
	Treaty2CashOutflowM7                    float64 `json:"treaty2_cash_outflow_m7" csv:"treaty2_cash_outflow_m7"`
	Treaty2CashOutflowM8                    float64 `json:"treaty2_cash_outflow_m8" csv:"treaty2_cash_outflow_m8"`
	Treaty2CashOutflowM9                    float64 `json:"treaty2_cash_outflow_m9" csv:"treaty2_cash_outflow_m9"`
	Treaty2CashOutflowM10                   float64 `json:"treaty2_cash_outflow_m10" csv:"treaty2_cash_outflow_m10"`
	Treaty2CashOutflowM11                   float64 `json:"treaty2_cash_outflow_m11" csv:"treaty2_cash_outflow_m11"`
	Treaty2CashOutflowM12                   float64 `json:"treaty2_cash_outflow_m12" csv:"treaty2_cash_outflow_m12"`
	Treaty3CashInflowM1                     float64 `json:"treaty3_cash_inflow_m1" csv:"treaty3_cash_inflow_m1"`
	Treaty3CashInflowM2                     float64 `json:"treaty3_cash_inflow_m2" csv:"treaty3_cash_inflow_m2"`
	Treaty3CashInflowM3                     float64 `json:"treaty3_cash_inflow_m3" csv:"treaty3_cash_inflow_m3"`
	Treaty3CashInflowM4                     float64 `json:"treaty3_cash_inflow_m4" csv:"treaty3_cash_inflow_m4"`
	Treaty3CashInflowM5                     float64 `json:"treaty3_cash_inflow_m5" csv:"treaty3_cash_inflow_m5"`
	Treaty3CashInflowM6                     float64 `json:"treaty3_cash_inflow_m6" csv:"treaty3_cash_inflow_m6"`
	Treaty3CashInflowM7                     float64 `json:"treaty3_cash_inflow_m7" csv:"treaty3_cash_inflow_m7"`
	Treaty3CashInflowM8                     float64 `json:"treaty3_cash_inflow_m8" csv:"treaty3_cash_inflow_m8"`
	Treaty3CashInflowM9                     float64 `json:"treaty3_cash_inflow_m9" csv:"treaty3_cash_inflow_m9"`
	Treaty3CashInflowM10                    float64 `json:"treaty3_cash_inflow_m10" csv:"treaty3_cash_inflow_m10"`
	Treaty3CashInflowM11                    float64 `json:"treaty3_cash_inflow_m11" csv:"treaty3_cash_inflow_m11"`
	Treaty3CashInflowM12                    float64 `json:"treaty3_cash_inflow_m12" csv:"treaty3_cash_inflow_m12"`
	Treaty3CashOutflowM1                    float64 `json:"treaty3_cash_outflow_m1" csv:"treaty3_cash_outflow_m1"`
	Treaty3CashOutflowM2                    float64 `json:"treaty3_cash_outflow_m2" csv:"treaty3_cash_outflow_m2"`
	Treaty3CashOutflowM3                    float64 `json:"treaty3_cash_outflow_m3" csv:"treaty3_cash_outflow_m3"`
	Treaty3CashOutflowM4                    float64 `json:"treaty3_cash_outflow_m4" csv:"treaty3_cash_outflow_m4"`
	Treaty3CashOutflowM5                    float64 `json:"treaty3_cash_outflow_m5" csv:"treaty3_cash_outflow_m5"`
	Treaty3CashOutflowM6                    float64 `json:"treaty3_cash_outflow_m6" csv:"treaty3_cash_outflow_m6"`
	Treaty3CashOutflowM7                    float64 `json:"treaty3_cash_outflow_m7" csv:"treaty3_cash_outflow_m7"`
	Treaty3CashOutflowM8                    float64 `json:"treaty3_cash_outflow_m8" csv:"treaty3_cash_outflow_m8"`
	Treaty3CashOutflowM9                    float64 `json:"treaty3_cash_outflow_m9" csv:"treaty3_cash_outflow_m9"`
	Treaty3CashOutflowM10                   float64 `json:"treaty3_cash_outflow_m10" csv:"treaty3_cash_outflow_m10"`
	Treaty3CashOutflowM11                   float64 `json:"treaty3_cash_outflow_m11" csv:"treaty3_cash_outflow_m11"`
	Treaty3CashOutflowM12                   float64 `json:"treaty3_cash_outflow_m12" csv:"treaty3_cash_outflow_m12"`
	User                                    string  `json:"user" csv:"user"`
	UserEmail                               string  `json:"user_email" csv:"user_email"`
}

type LICAggregatedProjections struct {
	ID                                  int    `json:"id" gorm:"primary_key"`
	JobProductID                        int    `json:"job_product_id" gorm:"index"`
	RunType                             int    `json:"projection_type" gorm:"index"` // 0 for valuation, 1 for pricing
	RunDate                             string `json:"run_date"`
	ProductCode                         string `json:"product_code" gorm:"index"`
	SpCode                              int
	IFRS17Group                         string
	ProjectionMonth                     int     `gorm:"column:projection_month"`
	InitialPolicy                       float64 `json:"initial_policy"`
	InitialPolicyAdjusted               float64 `json:"initial_policy_adjusted"`
	DisabilityOutgo                     float64
	DisabilityOutgoAdjusted             float64
	Rider                               float64
	RiderAdjusted                       float64
	Expenses                            float64
	ExpensesAdjusted                    float64
	NetCashFlow                         float64
	NetCashFlowAdjusted                 float64
	Reserves                            float64
	ReservesAdjusted                    float64
	CoverageUnits                       float64
	DiscountedMorbidityOutgo            float64 `json:"discounted_morbidity_outgo"`
	DiscountedMorbidityOutgoAdjusted    float64 `json:"discounted_morbidity_outgo_adjusted"`
	DiscountedOutgo                     float64 `json:"discounted_outgo"`
	DiscountedOutgoAdjusted             float64 `json:"discounted_outgo_adjusted"`
	DiscountedExpenses                  float64 `json:"discounted_expenses"`
	DiscountedExpensesAdjusted          float64 `json:"discounted_expenses_adjusted"`
	DiscountedAnnuityFactor             float64 `json:"annuity_factor"`
	DiscountedCashOutflow               float64 `json:"discounted_cash_outflow"`
	DiscountedCashInflow                float64 `json:"discounted_cash_inflow"`
	SumCoverageUnits                    float64 `json:"sum_coverage_units"`
	DiscountedCoverageUnits             float64 `json:"discounted_coverage_units"`
	CededSumAssured                     float64 `json:"ceded_sum_assured"`
	ReinsurancePremium                  float64 `json:"reinsurance_premium"`
	ReinsurancePremiumAdjusted          float64 `json:"reinsurance_premium_adjusted"`
	ReinsuranceCedingCommission         float64 `json:"reinsurance_ceding_commission"`
	ReinsuranceCedingCommissionAdjusted float64 `json:"reinsurance_ceding_commission_adjusted"`
	ReinsuranceClaims                   float64 `json:"reinsurance_claims"`
	ReinsuranceClaimsAdjusted           float64 `json:"reinsurance_claims_adjusted"`
}

type ProjectionJob struct {
	ID                int           `json:"id" gorm:"primary_key"`
	JobsTemplateID    int           `json:"jobs_template_id"`
	RunJobId          int           `json:"run_job_id"`
	Products          []JobProduct  `json:"products"`
	ProdIds           []int         `json:"prod_ids" gorm:"-"`
	ProdIdsJson       string        `json:"prod_ids_json" gorm:"column:prod_ids"`
	ProductID         int           `json:"product_id"`
	CreationDate      time.Time     `json:"creation_date"`
	RunTime           float64       `json:"run_time"`
	RunType           int           `json:"run_type"`
	RunDate           string        `json:"run_date"`
	RunName           string        `json:"run_name"`
	RunDescription    string        `json:"run_description"`
	Status            string        `json:"status"`
	StatusError       string        `json:"status_error"`
	TotalPoints       int           `json:"total_points"`
	PointsDone        int           `json:"points_done"`
	RunParameters     RunParameters `json:"run_parameters"`
	ShockSettingID    int           `json:"shock_setting_id"`
	AggregationPeriod int           `json:"aggregation_period"`
	YieldCurveBasis   string        `json:"yield_curve_basis"`
	YieldCurveMonth   int           `json:"yield_curve_month"`
	ModelPointVersion string        `json:"mp_version"`
	RunBasis          string        `json:"run_basis"`
	UserName          string        `json:"user_name"`
	UserEmail         string        `json:"user_email"`
}

func (p *ProjectionJob) AfterFind(tx *gorm.DB) error {
	if p.ProdIdsJson != "" {
		return json.Unmarshal([]byte(p.ProdIdsJson), &p.ProdIds)
	}
	return nil
}

type RunPayload struct {
	ID                int          `gorm:"primary_key"`
	RunJobID          int          `json:"run_job_id"`
	ProjectionJobID   int          `json:"projection_job_id"`
	JobsTemplateID    int          `json:"jobs_template_id"`
	ProdIds           []int        `json:"prod_ids"`
	RunType           int          `json:"run_type"`
	RunName           string       `json:"run_name"`
	RunSingle         bool         `json:"run_single"`
	RunDate           string       `json:"run_date"`
	UserName          string       `json:"user_name"`
	UserEmail         string       `json:"user_email"`
	ModelpointYear    int          `json:"modelpoint_year"`
	ModelPointVersion string       `json:"mp_version"`
	YieldcurveYear    int          `json:"yieldcurve_year"`
	YieldcurveMonth   int          `json:"yieldcurve_month"`
	ParameterYear     int          `json:"parameter_year"`
	MorbidityYear     int          `json:"morbidity_year"`
	MortalityYear     int          `json:"mortality_year"`
	LapseYear         int          `json:"lapse_year"`
	LapseMarginYear   int          `json:"lapse_margin_year"`
	RetrenchmentYear  int          `json:"retrenchment_year"`
	RunDescription    string       `json:"run_description"`
	IFRS17Indicator   bool         `json:"ifrs17_indicator"`
	LICIndicator      bool         `json:"lic_indicator"`
	SpCode            bool         `json:"spcode"`
	YieldCurveBasis   string       `json:"yield_curve_basis"`
	ShockSettings     ShockSetting `json:"shock_settings"`
	AggregationPeriod int          `json:"aggregation_period"`
	RunBasis          string       `json:"run_basis"`
	YearEndMonth      int          `json:"year_end_month"`
}

type JobProduct struct {
	ID              int    `json:"id" gorm:"primary_key"`
	ProductName     string `json:"product_name"`
	ProductID       int    `json:"product_id"`
	ProductCode     string `json:"product_code"`
	ProjectionJobID int    `json:"projection_job_id"`
	YieldCurveBasis string `json:"yield_curve_basis"`
	ModelPointCount int    `json:"model_point_count"`
	PointsDone      int    `json:"points_done"`
	JobStatus       string `json:"job_status"`
	JobStatusError  string `json:"job_status_error"`
}

type JobProductRunError struct {
	ID              int    `json:"id" gorm:"primary_key"`
	JobProductID    int    `json:"job_product_id"`
	ProjectionJobID int    `json:"projection_job_id"`
	ProductCode     string `json:"product_code"`
	Error           string `json:"error" gorm:"index:idx_failure_reason"`
	FailurePoint    string `json:"failure_point" gorm:"index:idx_failure_reason"`
}

type ProductRunParameter struct {
	ID                int    `gorm:"primary_key"`
	ProjectionJobID   int    `json:"projection_job_id"`
	RunType           int    `json:"run_type"`
	RunName           string `json:"run_name"`
	RunSingle         bool   `json:"run_single"`
	RunDate           string `json:"run_date"`
	User              string `json:"user"`
	ProductName       string `json:"product_name"`
	ProductCode       string `json:"product_code"`
	ModelpointYear    int    `json:"modelpoint_year"`
	ModelPointVersion string `json:"mp_version"`
	YieldcurveYear    int    `json:"yieldcurve_year"`
	YieldcurveMonth   int    `json:"yieldcurve_month"`
	ParameterYear     int    `json:"parameter_year"`
	TransitionYear    int    `json:"transition_year"`
	MorbidityYear     int    `json:"morbidity_year"`
	MortalityYear     int    `json:"mortality_year"`
	LapseYear         int    `json:"lapse_year"`
	LapseMarginYear   int    `json:"lapse_margin_year"`
	RetrenchmentYear  int    `json:"retrenchment_year"`
	RunDescription    string `json:"run_description"`
	IFRS17Indicator   bool   `json:"ifrs17_indicator"`
	IFRS17Group       bool   `json:"ifrs17group"`
	LICIndicator      bool   `json:"lic_indicator"`
	SpCode            int    `json:"spcode"`
	ShockSettingsID   int    `json:"shock_settings_id"`
	ShockSettingName  string `json:"shock_setting_name"`
	YieldCurveBasis   string `json:"yield_curve_basis"`
	AggregationPeriod int    `json:"aggregation_period"`
	RunBasis          string `json:"run_basis"`
	YearEndMonth      int    `json:"year_end_month"`
}

type CachedReserveResults struct {
	JobProductId int `json:"projection_job_id" gorm:"primary_key;auto_increment:false"`
	Variable     string
	ResultRange  int  `json:"result_range" gorm:"primary_key;auto_increment:false"`
	Result       JSON `sql:"type:json" json:"object"`
}

type ProductShock struct {
	//Year                          int     `json:"year" csv:"year"`
	//ProductCode                   string  `json:"product_code" csv:"product_code"`
	ProjectionMonth                int     `json:"projection_month" csv:"projection_month"`
	MultiplicativeMortality        float64 `json:"multiplicative_mortality" csv:"multiplicative_mortality"`
	AdditiveMortality              float64 `json:"additive_mortality" csv:"additive_mortality"`
	MultiplicativeLapse            float64 `json:"multiplicative_lapse" csv:"multiplicative_lapse"`
	AdditiveLapse                  float64 `json:"additive_lapse" csv:"additive_lapse"`
	MassLapse                      float64 `json:"mass_lapse" csv:"mass_lapse"`
	MultiplicativeDisability       float64 `json:"multiplicative_disability" csv:"multiplicative_disability"`
	AdditiveDisability             float64 `json:"additive_disability" csv:"additive_disability"`
	MultiplicativeRetrenchment     float64 `json:"multiplicative_retrenchment" csv:"multiplicative_retrenchment"`
	AdditiveRetrenchment           float64 `json:"additive_retrenchment" csv:"additive_retrenchment"`
	MultiplicativeCriticalIllness  float64 `json:"multiplicative_critical_illness" csv:"multiplicative_critical_illness"`
	AdditiveCriticalIllness        float64 `json:"additive_critical_illness" csv:"additive_critical_illness"`
	MultiplicativeYieldCurve       float64 `json:"multiplicative_yield_curve" csv:"multiplicative_yield_curve"`
	AdditiveYieldCurve             float64 `json:"additive_yield_curve" csv:"additive_yield_curve"`
	MultiplicativeExpense          float64 `json:"multiplicative_expense" csv:"multiplicative_expense"`
	AdditiveExpense                float64 `json:"additive_expense" csv:"additive_expense"`
	MultiplicativeInflation        float64 `json:"multiplicative_inflation" csv:"multiplicative_inflation"`
	AdditiveInflation              float64 `json:"additive_inflation" csv:"additive_inflation"`
	MortalityCatastropheFloor      float64 `json:"mortality_catastrophe_floor" csv:"mortality_catastrophe_floor"`
	MortalityCatastropheCeiling    float64 `json:"mortality_catastrophe_ceiling" csv:"mortality_catastrophe_ceiling"`
	CATScalar                      float64 `json:"cat_scalar" csv:"cat_scalar"`
	MortalityCatastropheMultiplier float64 `json:"mortality_catastrophe_multiplier" csv:"mortality_catastrophe_multiplier"`
	MorbidityCatastropheMultiplier float64 `json:"morbidity_catastrophe_multiplier" csv:"morbidity_catastrophe_multiplier"`
	ShockBasis                     string  `json:"shock_basis" csv:"shock_basis"`
	Created                        int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
	CreatedBy                      string  `json:"created_by" csv:"created_by"`
}

type ProductPricingShock struct {
	ProductCode                    string  `json:"product_code" csv:"product_code"`
	ProjectionMonth                int     `json:"projection_month" csv:"projection_month"`
	MultiplicativeMortality        float64 `json:"multiplicative_mortality" csv:"multiplicative_mortality"`
	AdditiveMortality              float64 `json:"additive_mortality" csv:"additive_mortality"`
	MultiplicativeLapse            float64 `json:"multiplicative_lapse" csv:"multiplicative_lapse"`
	AdditiveLapse                  float64 `json:"additive_lapse" csv:"additive_lapse"`
	MassLapse                      float64 `json:"mass_lapse" csv:"mass_lapse"`
	MultiplicativeDisability       float64 `json:"multiplicative_disability" csv:"multiplicative_disability"`
	AdditiveDisability             float64 `json:"additive_disability" csv:"additive_disability"`
	MultiplicativeRetrenchment     float64 `json:"multiplicative_retrenchment" csv:"multiplicative_retrenchment"`
	AdditiveRetrenchment           float64 `json:"additive_retrenchment" csv:"additive_retrenchment"`
	MultiplicativeCriticalIllness  float64 `json:"multiplicative_critical_illness" csv:"multiplicative_critical_illness"`
	AdditiveCriticalIllness        float64 `json:"additive_critical_illness" csv:"additive_critical_illness"`
	MultiplicativeYieldCurve       float64 `json:"multiplicative_yield_curve" csv:"multiplicative_yield_curve"`
	AdditiveYieldCurve             float64 `json:"additive_yield_curve" csv:"additive_yield_curve"`
	MultiplicativeExpense          float64 `json:"multiplicative_expense" csv:"multiplicative_expense"`
	AdditiveExpense                float64 `json:"additive_expense" csv:"additive_expense"`
	MultiplicativeInflation        float64 `json:"multiplicative_inflation" csv:"multiplicative_inflation"`
	AdditiveInflation              float64 `json:"additive_inflation" csv:"additive_inflation"`
	MortalityCatastropheFloor      float64 `json:"mortality_catastrophe_floor" csv:"mortality_catastrophe_floor"`
	MortalityCatastropheCeiling    float64 `json:"mortality_catastrophe_ceiling" csv:"mortality_catastrophe_ceiling"`
	MortalityCatastropheMultiplier float64 `json:"mortality_catastrophe_multiplier" csv:"mortality_catastrophe_multiplier"`
	MorbidityCatastropheMultiplier float64 `json:"morbidity_catastrophe_multiplier" csv:"morbidity_catastrophe_multiplier"`
	CatScalar                      float64 `json:"cat_scalar" csv:"cat_scalar"`
	ShockBasis                     string  `json:"shock_basis" csv:"shock_basis"`
}

type AvailableYieldResult struct {
	Year int
}

type AvailableLapseResult struct {
	Year string
}

type AvailableBasis struct {
	Basis string
	Name  string `json:"name"`
}

type JobsTemplate struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type JobTemplateContent struct {
	JobsTemplateID int    `json:"jobs_template_id"`
	Content        string `json:"content" gorm:"size:12000"`
}

type AggregatedVariableGroup struct {
	ID        int      `json:"id" gorm:"primary_key"`
	Name      string   `json:"name"`
	Variables []string `json:"variables" gorm:"serializer:json"`
}

type ExcelAggPayload struct {
	RunId           int    `json:"run_id"`
	ProductCode     string `json:"product_code"`
	SpCode          string `json:"sp_code"`
	JobProductId    int    `json:"job_product_id"`
	VariableGroupId int    `json:"variable_group_id"`
}
