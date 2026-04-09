package models

type ConsolidateResult struct {
	Month        int     `json:"month" gorm:"primary_key"`
	PremBudget   float64 `json:"prem_budget"`
	PremActual   float64 `json:"prem_actual"`
	CommBudget   float64 `json:"comm_budget"`
	CommActual   float64 `json:"comm_actual"`
	BinderBudget float64 `json:"binder_budget"`
	BinderActual float64 `json:"binder_actual"`
	ExpBudget    float64 `json:"exp_budget"`
	ExpActual    float64 `json:"exp_actual"`
	ClaimBudget  float64 `json:"claim_budget"`
	ClaimActual  float64 `json:"claim_actual"`
}

type CumulativeConsolidatedResult struct {
	Month        int     `json:"month" gorm:"primary_key"`
	PremBudget   float64 `json:"prem_budget"`
	PremActual   float64 `json:"prem_actual"`
	CommBudget   float64 `json:"comm_budget"`
	CommActual   float64 `json:"comm_actual"`
	BinderBudget float64 `json:"binder_budget"`
	BinderActual float64 `json:"binder_actual"`
	ExpBudget    float64 `json:"exp_budget"`
	ExpActual    float64 `json:"exp_actual"`
	ClaimBudget  float64 `json:"claim_budget"`
	ClaimActual  float64 `json:"claim_actual"`
}

type AnnualConsolidatedResult struct {
	Year         int     `json:"year" gorm:"primary_key"`
	PremBudget   float64 `json:"prem_budget"`
	PremActual   float64 `json:"prem_actual"`
	CommBudget   float64 `json:"comm_budget"`
	CommActual   float64 `json:"comm_actual"`
	BinderBudget float64 `json:"binder_budget"`
	BinderActual float64 `json:"binder_actual"`
	ExpBudget    float64 `json:"exp_budget"`
	ExpActual    float64 `json:"exp_actual"`
	ClaimBudget  float64 `json:"claim_budget"`
	ClaimActual  float64 `json:"claim_actual"`
}

type ParameterDisclosureItem struct {
	Variable     string      `json:"variable"`
	CurrentYear  interface{} `json:"current_year"`
	PreviousYear interface{} `json:"previous_year"`
	Variance     float64     `json:"variance"`
	Change       float64     `json:"change"`
	IsNumeric    bool        `json:"is_numeric"`
}

type MortalityDisclosureItem struct {
	ANB         int     `json:"anb"`
	Gender      string  `json:"gender"`
	CurrentYear float64 `json:"current_year"`
	PastYear    float64 `json:"past_year"`
}

type YieldDisclosureItem struct {
	ProjectionTime         int     `json:"projection_time"`
	CurrentYearNominalRate float64 `json:"current_year_nominal_rate"`
	PastYearNominalRate    float64 `json:"past_year_nominal_rate"`
	CurrentYearInflation   float64 `json:"current_year_inflation"`
	PastYearInflation      float64 `json:"past_year_inflation"`
}

type LapseDisclosureItem struct {
	DurationInForceMonths int     `json:"duration_in_force_months"`
	CurrentYear           float64 `json:"current_year"`
	PastYear              float64 `json:"past_year"`
	Variance              float64 `json:"variance"`
	Change                float64 `json:"change"`
}

type ReserveSummary struct {
	ValuationDate                  string  `json:"valuation_date" csv:"valuation_date"`
	ProductCode                    string  `json:"product_code" csv:"product_code"`
	SpCode                         int     `json:"sp_code" csv:"sp_code"`
	ProductName                    string  `json:"product_name,omitempty" csv:"product_name"`
	RunName                        string  `json:"run_name" csv:"run_name"`
	PolicyCount                    float64 `json:"policy_count" csv:"policy_count"`
	InitialPolicy                  float64 `json:"initial_policy" csv:"initial_policy"`
	SumAssured                     float64 `json:"sum_assured" csv:"sum_assured"`
	UnfundedUnitFund               float64 `json:"unfunded_unit_fund" csv:"unfunded_unit_fund"`
	BonusStabilisationAccount      float64 `json:"bonus_stabilisation_account" csv:"bonus_stabilisation_account"`
	AnnuityIncome                  float64 `json:"annuity_income" csv:"annuity_income	"`
	AnnualPremium                  float64 `json:"annual_premium" csv:"annual_premium"`
	BEL                            float64 `json:"bel" csv:"bel"`
	BELAdjusted                    float64 `json:"bel_adjusted" csv:"bel_adjusted"`
	ReinsuranceBEL                 float64 `json:"reinsurance_bel" csv:"reinsurance_bel"`
	ReinsuranceBELAdjusted         float64 `json:"reinsurance_bel_adjusted" csv:"reinsurance_bel_adjusted"`
	Ibnr                           float64 `json:"ibnr" csv:"ibnr"`
	Upr                            float64 `json:"upr" csv:"upr"`
	Vif                            float64 `json:"vif" csv:"vif"`
	VifAdjusted                    float64 `json:"vif_adjusted" csv:"vif_adjusted"`
	DiscountedCorporateTax         float64 `json:"discounted_corporate_tax" csv:"discounted_corporate_tax"`
	DiscountedCorporateTaxAdjusted float64 `json:"discounted_corporate_tax_adjusted" csv:"discounted_corporate_tax_adjusted"`
	DiscountedPremiumIncome        float64 `json:"discounted_premium_income" csv:"discounted_premium_income"`
}

type PAAReserveSummary struct {
	ValuationDate                     string  `json:"valuation_date" csv:"valuation_date"`
	RunName                           string  `json:"run_name" csv:"run_name"`
	PortfolioName                     string  `json:"portfolio_name" csv:"portfolio_name"`
	ProductCode                       string  `json:"product_code" csv:"product_code"`
	IFRS17Group                       string  `json:"ifrs17_group" csv:"ifrs17_group"`
	PolicyCount                       float64 `json:"policy_count" csv:"policy_count"`
	SumFutureEarnedPremium            float64 `json:"sum_future_earned_premium" csv:"sum_future_earned_premium"`
	CurrentPeriodEarnedPremium        float64 `json:"current_period_earned_premium" csv:"current_period_earned_premium"`
	CurrentPeriodAmortisedAcquisition float64 `json:"current_period_amortised_acquisition" csv:"current_period_amortised_acquisition"`
	CurrentPeriodInsuranceRevenue     float64 `json:"current_period_insurance_revenue" csv:"current_period_insurance_revenue"`
	GMMBel                            float64 `json:"gmm_bel" csv:"gmm_bel"`
	RiskAdjustment                    float64 `json:"risk_adjustment" csv:"risk_adjustment"`
}
