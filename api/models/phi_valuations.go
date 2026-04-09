package models

import "time"

type PhiYieldCurve struct {
	ID           int `json:"id" gorm:"primary_key"`
	Year         int
	Month        int `json:"month" csv:"month" gorm:"primary_key;auto_increment:false;column:month"`
	Version      string
	ProjTime     int       `json:"proj_time" csv:"proj_time" gorm:"primary_key;auto_increment:false;column:proj_time"`
	NominalRate  float64   `json:"nominal_rate" csv:"nominal_rate"`
	Inflation    float64   `json:"inflation" csv:"inflation"`
	CreatedBy    string    `json:"created_by" csv:"created_by"`
	CreationDate time.Time `json:"creation_date" gorm:"autoCreateTime"`
}

type PhiParameter struct {
	ID                                      int `json:"id" gorm:"primary_key"`
	Year                                    int
	Version                                 string
	ProductCode                             string    `json:"product_code" csv:"product_code"`
	Spcode                                  int       `json:"spcode" csv:"spcode"`
	AnnualRenewalExpenseAmount              float64   `json:"annual_renewal_expense_amount" csv:"annual_renewal_expense_amount"`
	ClaimsExpenseProportion                 float64   `json:"claims_expense_proportion" csv:"claims_expense_proportion"`
	AnnualClaimsExpenseAmount               float64   `json:"annual_claims_expense_amount" csv:"annual_claims_expense_amount"`
	AttributableExpenseProportion           float64   `json:"attributable_expense_proportion" csv:"attributable_expense_proportion"`
	NotifiedAcceptanceRatio                 float64   `json:"notified_acceptance_ratio" csv:"notified_acceptance_ratio"`
	PendingAcceptanceRatio                  float64   `json:"pending_acceptance_ratio" csv:"pending_acceptance_ratio"`
	AnnuityCalendarMonthEscalationIndicator bool      `json:"annuity_calendar_month_escalation_indicator" csv:"annuity_calendar_month_escalation_indicator"`
	CreatedBy                               string    `json:"created_by" csv:"created_by"`
	CreationDate                            time.Time `json:"creation_date" gorm:"autoCreateTime"`
}

type PhiRecoveryRate struct {
	ID            int       `json:"id" gorm:"primary_key"`
	Year          int       `json:"year" csv:"year"`
	Version       string    `json:"version" csv:"version"`
	Anb           int       `json:"anb" csv:"anb"`
	Gender        string    `json:"gender" csv:"gender"`
	DurationIfM   int       `json:"duration_if_m" csv:"duration_if_m"`
	WaitingPeriod int       `json:"waiting_period" csv:"waiting_period"`
	RecoveryRate  float64   `json:"recovery_rate" csv:"recovery_rate"`
	CreatedBy     string    `json:"created_by" csv:"created_by"`
	CreationDate  time.Time `json:"creation_date" gorm:"autoCreateTime"`
}

type PhiMortality struct {
	ID            int       `json:"id" gorm:"primary_key"`
	Year          int       `json:"year" csv:"year"`
	Version       string    `json:"version" csv:"version"`
	Anb           int       `json:"anb" csv:"anb"`
	Gender        string    `json:"gender" csv:"gender"`
	MortalityRate float64   `json:"mortality_rate" csv:"mortality_rate"`
	CreatedBy     string    `json:"created_by" csv:"created_by"`
	CreationDate  time.Time `json:"creation_date"`
}
type PhiShock struct {
	ID                             int       `json:"id" gorm:"primary_key"`
	Year                           int       `json:"year" csv:"year"`
	Version                        string    `json:"version" csv:"version"`
	ProjectionMonth                int       `json:"projection_month" csv:"projection_month"`
	MultiplicativeMortality        float64   `json:"multiplicative_mortality" csv:"multiplicative_mortality"`
	AdditiveMortality              float64   `json:"additive_mortality" csv:"additive_mortality"`
	MultiplicativeRecovery         float64   `json:"multiplicative_recovery" csv:"multiplicative_recovery"`
	AdditiveRecovery               float64   `json:"additive_recovery" csv:"additive_recovery"`
	MultiplicativeYieldCurve       float64   `json:"multiplicative_yield_curve" csv:"multiplicative_yield_curve"`
	AdditiveYieldCurve             float64   `json:"additive_yield_curve" csv:"additive_yield_curve"`
	MultiplicativeExpense          float64   `json:"multiplicative_expense" csv:"multiplicative_expense"`
	AdditiveExpense                float64   `json:"additive_expense" csv:"additive_expense"`
	MultiplicativeInflation        float64   `json:"multiplicative_inflation" csv:"multiplicative_inflation"`
	AdditiveInflation              float64   `json:"additive_inflation" csv:"additive_inflation"`
	MortalityCatastropheFloor      float64   `json:"mortality_catastrophe_floor" csv:"mortality_catastrophe_floor"`
	MortalityCatastropheCeiling    float64   `json:"mortality_catastrophe_ceiling" csv:"mortality_catastrophe_ceiling"`
	CATScalar                      float64   `json:"cat_scalar" csv:"cat_scalar"`
	MortalityCatastropheMultiplier float64   `json:"mortality_catastrophe_multiplier" csv:"mortality_catastrophe_multiplier"`
	MorbidityCatastropheMultiplier float64   `json:"morbidity_catastrophe_multiplier" csv:"morbidity_catastrophe_multiplier"`
	ShockBasis                     string    `json:"shock_basis" csv:"shock_basis"`
	CreatedBy                      string    `json:"created_by" csv:"created_by"`
	CreationDate                   time.Time `json:"creation_date" gorm:"autoCreateTime"`
}

type PhiModelPoint struct {
	ID                     int       `json:"id" gorm:"primary_key"`
	Year                   int       `json:"year" csv:"year"`
	Version                string    `json:"version" csv:"version"`
	PolicyNumber           string    `json:"policy_number" csv:"policy_number"`
	Spcode                 int       `json:"spcode" csv:"spcode"`
	IFRS17Group            string    `json:"ifrs17_group" csv:"ifrs17_group"`
	ProductCode            string    `json:"product_code" csv:"product_code"`
	LockedInYear           int       `json:"locked_in_year" csv:"locked_in_year"`
	LockedInMonth          int       `json:"locked_in_month" csv:"locked_in_month"`
	DurationInForceMonths  int       `json:"duration_in_force_months" csv:"duration_in_force_months"`
	AgeAtEntry             int       `json:"age_at_entry" csv:"age_at_entry"`
	Gender                 string    `json:"gender" csv:"gender"`
	InitialPolicy          int       `json:"initial_policy" csv:"initial_policy"`
	AnnuityIncome          float64   `json:"annuity_income" csv:"annuity_income"`
	AnnuityIncomeFrequency int       `json:"annuity_income_frequency" csv:"annuity_income_frequency"`
	TermInMonths           int       `json:"term_in_months" csv:"term_in_months"`
	WaitingPeriod          int       `json:"waiting_period" csv:"waiting_period"`
	AnnuityEscalation      float64   `json:"annuity_escalation" csv:"annuity_escalation"`
	AnnuityEscalationMonth int       `json:"annuity_escalation_month" csv:"annuity_escalation_month"`
	Plan                   string    `json:"plan" csv:"plan"`
	ClaimStatus            string    `json:"claim_status" csv:"claim_status"`
	SocioEconomicClass     int       `json:"socio_economic_class" csv:"socio_economic_class"`
	OccupationalClass      string    `json:"occupational_class" csv:"occupational_class"`
	DisabilityDefinition   int       `json:"disability_definition" csv:"disability_definition"`
	TreatyCode             string    `json:"treaty_code" csv:"treaty_code"`
	CreatedBy              string    `json:"created_by" csv:"created_by"`
	CreationDate           time.Time `json:"creation_date" `
}

type PhiReinsurance struct {
	ID                    int       `json:"id"  gorm:"primary_key" `
	Year                  int       `json:"year" csv:"year"`
	ProductCode           string    `json:"product_code" csv:"product_code"`
	Version               string    `json:"version" csv:"version"`
	Level1CededProportion float64   `json:"level1_ceded_proportion" csv:"level1_ceded_proportion"`
	Level1Lowerbound      float64   `json:"level1_lowerbound" csv:"level1_lowerbound"`
	Level1Upperbound      float64   `json:"level1_upperbound" csv:"level1_upperbound"`
	Level2CededProportion float64   `json:"level2_ceded_proportion" csv:"level2_ceded_proportion"`
	Level2Lowerbound      float64   `json:"level2_lowerbound" csv:"level2_lowerbound"`
	Level2Upperbound      float64   `json:"level2_upperbound" csv:"level2_upperbound"`
	Level3CededProportion float64   `json:"level3_ceded_proportion" csv:"level3_ceded_proportion"`
	Level3Lowerbound      float64   `json:"level3_lowerbound" csv:"level3_lowerbound"`
	Level3Upperbound      float64   `json:"level3_upperbound" csv:"level3_upperbound"`
	CedingCommission      float64   `json:"ceding_commission" csv:"ceding_commission"`
	LeadReProportion      float64   `json:"lead_re_proportion" csv:"lead_re_proportion"`
	Re2Proportion         float64   `json:"re2_proportion" csv:"re2_proportion"`
	Re3Proportion         float64   `json:"re3_proportion" csv:"re3_proportion"`
	CreatedBy             string    `json:"created_by" csv:"created_by" `
	CreationDate          time.Time `json:"creation_date" gorm:"autoCreateTime" `
}

type PhiShockSetting struct {
	ID                   int       `json:"id" gorm:"primary_key"`
	Name                 string    `json:"name" gorm:"unique"`
	Description          string    `json:"description"`
	Mortality            bool      `json:"mortality"`
	RealYieldCurve       bool      `json:"real_yield_curve"`
	NominalYieldCurve    bool      `json:"nominal_yield_curve"`
	Expense              bool      `json:"expense"`
	Inflation            bool      `json:"inflation"`
	MortalityCatastrophe bool      `json:"mortality_catastrophe"`
	MorbidityCatastrophe bool      `json:"morbidity_catastrophe"`
	Recovery             bool      `json:"recovery"`
	ShockBasis           string    `json:"shock_basis"`
	Year                 int       `json:"year"`
	CreatedBy            string    `json:"created_by" csv:"created_by"`
	CreationDate         time.Time `json:"creation_date" gorm:"autoCreateTime"`
}

type PhiProjections struct {
	ID                                       int     `json:"-" gorm:"primary_key"`
	JobProductID                             int     `json:"-" gorm:"index"`
	RunDate                                  string  `json:"run_date"`
	RunId                                    int     `json:"run_id" gorm:"index"`
	RunName                                  string  `json:"run_name"`
	ProductCode                              string  `json:"product_code" gorm:"index"`
	PolicyNumber                             string  `json:"policy_number" gorm:"index"`
	SpCode                                   int     `json:"sp_code"`
	IFRS17Group                              string  `json:"ifrs_17_group"`
	ProjectionMonth                          int     `json:"projection_month" gorm:"index"`
	ProjectionYear                           int     `json:"projection_year"`
	ValuationTimeMonth                       int     `json:"valuation_time_month"`
	CalendarMonth                            int     `json:"calendar_month"`
	ValuationTimeYear                        float64 `json:"valuation_time_year"`
	AgeNextBirthday                          int     `json:"age_next_birthday"`
	InflationFactor                          float64 `json:"inflation_factor"`
	InflationFactorAdjusted                  float64 `json:"inflation_factor_adjusted"`
	AnnuityEscalationRate                    float64 `json:"annuity_escalation_rate"`
	AnnuityEscalation                        float64 `json:"annuity_escalation"`
	BaseMortalityRate                        float64 `json:"base_mortality_rate"`
	BaseMortalityRateAdjusted                float64 `json:"base_mortality_rate_adjusted"`
	BaseRecoveryRate                         float64 `json:"base_recovery_rate"`
	BaseRecoveryRateAdjusted                 float64 `json:"base_recovery_rate_adjusted"`
	IndependentMortalityRateMonthly          float64 `json:"independent_mortality_rate_monthly"`
	IndependentMortalityRateAdjustedByMonth  float64 `json:"independent_mortality_rate_adjusted_by_month"`
	IndependentRecoveryRateMonthly           float64 `json:"independent_recovery_rate_monthly"`
	IndependentRecoveryRateMonthlyAdjusted   float64 `json:"independent_recovery_rate_monthly_adjusted"`
	MonthlyDependentMortality                float64 `json:"monthly_dependent_mortality"`
	MonthlyDependentMortalityAdjusted        float64 `json:"monthly_dependent_mortality_adjusted"`
	MonthlyDependentRecovery                 float64 `json:"monthly_dependent_recovery"`
	MonthlyDependentRecoveryAdjusted         float64 `json:"monthly_dependent_recovery_adjusted"`
	InitialPolicy                            float64 `gorm:"index" json:"initial_policy"`
	InitialPolicyAdjusted                    float64 `gorm:"index" json:"initial_policy_adjusted"`
	NumberOfDeathsInForce                    float64 `json:"number_of_deaths_in_force"`
	NumberOfDeathsInForceAdjusted            float64 `json:"number_of_deaths_in_force_adjusted"`
	NumberOfRecoveries                       float64 `json:"number_of_recoveries"`
	NumberOfRecoveriesAdjusted               float64 `json:"number_of_recoveries_adjusted"`
	NumberOfMaturities                       float64 `json:"number_of_maturities"`
	NumberOfMaturitiesAdjusted               float64 `json:"number_of_maturities_adjusted"`
	IncrementalDeaths                        float64 `json:"incremental_deaths"`
	IncrementalDeathsAdjusted                float64 `json:"incremental_deaths_adjusted"`
	IncrementalRecoveries                    float64 `json:"incremental_recoveries"`
	IncrementalRecoveriesAdjusted            float64 `json:"incremental_recoveries_adjusted"`
	SumAssured                               float64 `json:"sum_assured"`
	RiderSumAssured                          float64 `json:"rider_sum_assured"`
	AnnuityIncome                            float64 `json:"annuity_income"`
	Premium                                  float64 `json:"premium"`
	MaturityValue                            float64 `json:"maturity_value"`
	PremiumIncome                            float64 `json:"premium_income"`
	PremiumIncomeAdjusted                    float64 `json:"premium_income_adjusted"`
	SurrenderOutgo                           float64 `json:"surrender_outgo"`
	SurrenderOutgoAdjusted                   float64 `json:"surrender_outgo_adjusted"`
	DeathOutgo                               float64 `json:"death_outgo"`
	DeathOutgoAdjusted                       float64 `json:"death_outgo_adjusted"`
	RecoveryOutgo                            float64 `json:"retrenchment_outgo"`
	RecoveryOutgoAdjusted                    float64 `json:"retrenchment_outgo_adjusted"`
	AnnuityOutgo                             float64 `json:"annuity_outgo"`
	AnnuityOutgoAdjusted                     float64 `json:"annuity_outgo_adjusted"`
	Rider                                    float64 `json:"rider"`
	RiderAdjusted                            float64 `json:"rider_adjusted"`
	InitialExpenses                          float64 `json:"initial_expenses"`
	InitialExpensesAdjusted                  float64 `json:"initial_expenses_adjusted"`
	RenewalExpenses                          float64 `json:"renewal_expenses"`
	RenewalExpensesAdjusted                  float64 `json:"renewal_expenses_adjusted"`
	MaturityOutgo                            float64 `json:"maturity_outgo"`
	MaturityOutgoAdjusted                    float64 `json:"maturity_outgo_adjusted"`
	NetCashFlow                              float64 `gorm:"index" json:"net_cash_flow"`
	NetCashFlowAdjusted                      float64 `gorm:"index" json:"net_cash_flow_adjusted"`
	ValuationRate                            float64 `json:"valuation_rate"`
	ValuationRateAdjusted                    float64 `json:"valuation_rate_adjusted"`
	Reserves                                 float64 `gorm:"index" json:"reserves"`
	ReservesAdjusted                         float64 `gorm:"index" json:"reserves_adjusted"`
	ChangeInReserves                         float64 `gorm:"index" json:"change_in_reserves"`
	ChangeInReservesAdjusted                 float64 `gorm:"index" json:"change_in_reserves_adjusted"`
	InvestmentIncome                         float64 `json:"investment_income"`
	InvestmentIncomeAdjusted                 float64 `json:"investment_income_adjusted"`
	Profit                                   float64 `json:"profit"`
	ProfitAdjusted                           float64 `json:"profit_adjusted"`
	RiskDiscountRate                         float64 `json:"risk_discount_rate"`
	RiskDiscountRateAdjusted                 float64 `json:"risk_discount_rate_adjusted"`
	VIF                                      float64 `json:"vif"`
	VIFAdjusted                              float64 `json:"vif_adjusted"`
	CorporateTax                             float64 `json:"corporate_tax"`
	CorporateTaxAdjusted                     float64 `json:"corporate_tax_adjusted"`
	CoverageUnits                            float64 `json:"coverage_units"`
	DiscountedPremiumIncome                  float64 `json:"discounted_premium_income"`
	DiscountedPremiumIncomeAdjusted          float64 `json:"discounted_premium_income_adjusted"`
	DiscountedSurrenderOutgo                 float64 `json:"discounted_surrender_outgo"`
	DiscountedSurrenderOutgoAdjusted         float64 `json:"discounted_surrender_outgo_adjusted"`
	DiscountedMaturityOutgo                  float64 `json:"discounted_maturity_outgo"`
	DiscountedMaturityOutgoAdjusted          float64 `json:"discounted_maturity_outgo_adjusted"`
	DiscountedDeathOutgo                     float64 `json:"discounted_death_outgo"`
	DiscountedDeathOutgoAdjusted             float64 `json:"discounted_death_outgo_adjusted"`
	DiscountedRecoveryOutgo                  float64 `json:"discounted_recovery_outgo"`
	DiscountedRecoveryOutgoAdjusted          float64 `json:"discounted_recovery_outgo_adjusted"`
	DiscountedAnnuityOutgo                   float64 `json:"discounted_annuity_outgo"`
	DiscountedAnnuityOutgoAdjusted           float64 `json:"discounted_annuity_outgo_adjusted"`
	DiscountedRider                          float64 `json:"discounted_rider"`
	DiscountedRiderAdjusted                  float64 `json:"discounted_rider_adjusted"`
	DiscountedInitialExpenses                float64 `json:"discounted_initial_expenses"`
	DiscountedInitialExpensesAdjusted        float64 `json:"discounted_initial_expenses_adjusted"`
	DiscountedRenewalExpenses                float64 `json:"discounted_renewal_expenses"`
	DiscountedRenewalExpensesAdjusted        float64 `json:"discounted_renewal_expenses_adjusted"`
	DiscountedInvestmentIncome               float64 `json:"discounted_investment_income"`
	DiscountedInvestmentIncomeAdjusted       float64 `json:"discounted_investment_income_adjusted"`
	DiscountedProfit                         float64 `json:"discounted_profit"`
	DiscountedProfitAdjusted                 float64 `json:"discounted_profit_adjusted"`
	DiscountedCorporateTax                   float64 `json:"discounted_corporate_tax"`
	DiscountedCorporateTaxAdjusted           float64 `json:"discounted_corporate_tax_adjusted"`
	DiscountedAnnuityFactor                  float64 `json:"discounted_annuity_factor"`
	DiscountedCashOutflow                    float64 `json:"discounted_cash_outflow"`
	DiscountedCashOutflowExclAcquisition     float64 `json:"discounted_cash_outflow_excl_acquisition"`
	DiscountedAcquisitionCost                float64 `json:"discounted_acquisition_cost"`
	DiscountedCashInflow                     float64 `json:"discounted_cash_inflow"`
	SumCoverageUnits                         float64 `json:"sum_coverage_units"`
	DiscountedCoverageUnits                  float64 `json:"discounted_coverage_units"`
	CededAnnuityIncome                       float64 `json:"ceded_annuity_income"`
	LeadReAnnuityIncome                      float64 `json:"lead_re_annuity_income"`
	Re2AnnuityIncome                         float64 `json:"re_2_annuity_income"`
	Re3AnnuityIncome                         float64 `json:"re_3_annuity_income"`
	LeadReAnnuityOutgo                       float64 `json:"lead_re_annuity_outgo"`
	LeadReAnnuityOutgoAdjusted               float64 `json:"lead_re_annuity_outgo_adjusted"`
	Re2AnnuityOutgo                          float64 `json:"re_2_annuity_outgo"`
	Re2AnnuityOutgoAdjusted                  float64 `json:"re_2_annuity_outgo_adjusted"`
	Re3AnnuityOutgo                          float64 `json:"re_3_annuity_outgo"`
	Re3AnnuityOutgoAdjusted                  float64 `json:"re_3_annuity_outgo_adjusted"`
	DiscountedLeadReAnnuityOutgo             float64 `json:"discounted_lead_re_annuity_outgo"`
	DiscountedLeadReAnnuityOutgoAdjusted     float64 `json:"discounted_lead_re_annuity_outgo_adjusted"`
	DiscountedRe2AnnuityOutgo                float64 `json:"discounted_re_2_annuity_outgo"`
	DiscountedRe2AnnuityOutgoAdjusted        float64 `json:"discounted_re_2_annuity_outgo_adjusted"`
	DiscountedRe3AnnuityOutgo                float64 `json:"discounted_re_3_annuity_outgo"`
	DiscountedRe3AnnuityOutgoAdjusted        float64 `json:"discounted_re_3_annuity_outgo_adjusted"`
	NetReinsuranceCashflow                   float64 `json:"net_reinsurance_cashflow"`
	NetReinsuranceCashflowAdjusted           float64 `json:"net_reinsurance_cashflow_adjusted"`
	DiscountedNetLeadReCashflow              float64 `json:"discounted_net_lead_re_cashflow"`
	DiscountedNetLeadReCashflowAdjusted      float64 `json:"discounted_net_lead_re_cashflow_adjusted"`
	DiscountedNetRe2Cashflow                 float64 `json:"discounted_net_re_2_cashflow"`
	DiscountedNetRe2CashflowAdjusted         float64 `json:"discounted_net_re_2_cashflow_adjusted"`
	DiscountedNetRe3Cashflow                 float64 `json:"discounted_net_re_3_cashflow"`
	DiscountedNetRe3CashflowAdjusted         float64 `json:"discounted_net_re_3_cashflow_adjusted"`
	DiscountedNetReinsuranceCashflow         float64 `json:"discounted_net_reinsurance_cashflow"`
	DiscountedNetReinsuranceCashflowAdjusted float64 `json:"discounted_net_reinsurance_cashflow_adjusted"`
	RunBasis                                 string  `json:"run_basis"`
}

type PhiAggregatedProjections struct {
	ID                                       int     `json:"-" gorm:"primary_key"`
	JobProductID                             int     `json:"-" gorm:"index"`
	RunDate                                  string  `json:"run_date"`
	RunId                                    int     `json:"run_id" gorm:"index"`
	RunName                                  string  `json:"run_name"`
	ProductCode                              string  `json:"product_code" gorm:"index"`
	PolicyNumber                             string  `json:"policy_number" gorm:"index"`
	SpCode                                   int     `json:"sp_code"`
	IFRS17Group                              string  `json:"ifrs_17_group"`
	ProjectionMonth                          int     `json:"projection_month" gorm:"index"`
	ProjectionYear                           int     `json:"projection_year"`
	ValuationTimeMonth                       int     `json:"valuation_time_month"`
	CalendarMonth                            int     `json:"calendar_month"`
	ValuationTimeYear                        float64 `json:"valuation_time_year"`
	AgeNextBirthday                          int     `json:"age_next_birthday"`
	InflationFactor                          float64 `json:"inflation_factor"`
	InflationFactorAdjusted                  float64 `json:"inflation_factor_adjusted"`
	AnnuityEscalationRate                    float64 `json:"annuity_escalation_rate"`
	AnnuityEscalation                        float64 `json:"annuity_escalation"`
	BaseMortalityRate                        float64 `json:"base_mortality_rate"`
	BaseMortalityRateAdjusted                float64 `json:"base_mortality_rate_adjusted"`
	BaseRecoveryRate                         float64 `json:"base_recovery_rate"`
	BaseRecoveryRateAdjusted                 float64 `json:"base_recovery_rate_adjusted"`
	IndependentMortalityRateMonthly          float64 `json:"independent_mortality_rate_monthly"`
	IndependentMortalityRateAdjustedByMonth  float64 `json:"independent_mortality_rate_adjusted_by_month"`
	IndependentRecoveryRateMonthly           float64 `json:"independent_recovery_rate_monthly"`
	IndependentRecoveryRateMonthlyAdjusted   float64 `json:"independent_recovery_rate_monthly_adjusted"`
	MonthlyDependentMortality                float64 `json:"monthly_dependent_mortality"`
	MonthlyDependentMortalityAdjusted        float64 `json:"monthly_dependent_mortality_adjusted"`
	MonthlyDependentRecovery                 float64 `json:"monthly_dependent_recovery"`
	MonthlyDependentRecoveryAdjusted         float64 `json:"monthly_dependent_recovery_adjusted"`
	InitialPolicy                            float64 `gorm:"index" json:"initial_policy"`
	InitialPolicyAdjusted                    float64 `gorm:"index" json:"initial_policy_adjusted"`
	NumberOfDeathsInForce                    float64 `json:"number_of_deaths_in_force"`
	NumberOfDeathsInForceAdjusted            float64 `json:"number_of_deaths_in_force_adjusted"`
	NumberOfRecoveries                       float64 `json:"number_of_recoveries"`
	NumberOfRecoveriesAdjusted               float64 `json:"number_of_recoveries_adjusted"`
	NumberOfMaturities                       float64 `json:"number_of_maturities"`
	NumberOfMaturitiesAdjusted               float64 `json:"number_of_maturities_adjusted"`
	IncrementalDeaths                        float64 `json:"incremental_deaths"`
	IncrementalDeathsAdjusted                float64 `json:"incremental_deaths_adjusted"`
	IncrementalRecoveries                    float64 `json:"incremental_recoveries"`
	IncrementalRecoveriesAdjusted            float64 `json:"incremental_recoveries_adjusted"`
	SumAssured                               float64 `json:"sum_assured"`
	RiderSumAssured                          float64 `json:"rider_sum_assured"`
	AnnuityIncome                            float64 `json:"annuity_income"`
	Premium                                  float64 `json:"premium"`
	MaturityValue                            float64 `json:"maturity_value"`
	PremiumIncome                            float64 `json:"premium_income"`
	PremiumIncomeAdjusted                    float64 `json:"premium_income_adjusted"`
	SurrenderOutgo                           float64 `json:"surrender_outgo"`
	SurrenderOutgoAdjusted                   float64 `json:"surrender_outgo_adjusted"`
	DeathOutgo                               float64 `json:"death_outgo"`
	DeathOutgoAdjusted                       float64 `json:"death_outgo_adjusted"`
	RecoveryOutgo                            float64 `json:"retrenchment_outgo"`
	RecoveryOutgoAdjusted                    float64 `json:"retrenchment_outgo_adjusted"`
	AnnuityOutgo                             float64 `json:"annuity_outgo"`
	AnnuityOutgoAdjusted                     float64 `json:"annuity_outgo_adjusted"`
	Rider                                    float64 `json:"rider"`
	RiderAdjusted                            float64 `json:"rider_adjusted"`
	InitialExpenses                          float64 `json:"initial_expenses"`
	InitialExpensesAdjusted                  float64 `json:"initial_expenses_adjusted"`
	RenewalExpenses                          float64 `json:"renewal_expenses"`
	RenewalExpensesAdjusted                  float64 `json:"renewal_expenses_adjusted"`
	MaturityOutgo                            float64 `json:"maturity_outgo"`
	MaturityOutgoAdjusted                    float64 `json:"maturity_outgo_adjusted"`
	NetCashFlow                              float64 `gorm:"index" json:"net_cash_flow"`
	NetCashFlowAdjusted                      float64 `gorm:"index" json:"net_cash_flow_adjusted"`
	ValuationRate                            float64 `json:"valuation_rate"`
	ValuationRateAdjusted                    float64 `json:"valuation_rate_adjusted"`
	Reserves                                 float64 `gorm:"index" json:"reserves"`
	ReservesAdjusted                         float64 `gorm:"index" json:"reserves_adjusted"`
	ChangeInReserves                         float64 `gorm:"index" json:"change_in_reserves"`
	ChangeInReservesAdjusted                 float64 `gorm:"index" json:"change_in_reserves_adjusted"`
	InvestmentIncome                         float64 `json:"investment_income"`
	InvestmentIncomeAdjusted                 float64 `json:"investment_income_adjusted"`
	Profit                                   float64 `json:"profit"`
	ProfitAdjusted                           float64 `json:"profit_adjusted"`
	RiskDiscountRate                         float64 `json:"risk_discount_rate"`
	RiskDiscountRateAdjusted                 float64 `json:"risk_discount_rate_adjusted"`
	VIF                                      float64 `json:"vif"`
	VIFAdjusted                              float64 `json:"vif_adjusted"`
	CorporateTax                             float64 `json:"corporate_tax"`
	CorporateTaxAdjusted                     float64 `json:"corporate_tax_adjusted"`
	CoverageUnits                            float64 `json:"coverage_units"`
	DiscountedPremiumIncome                  float64 `json:"discounted_premium_income"`
	DiscountedPremiumIncomeAdjusted          float64 `json:"discounted_premium_income_adjusted"`
	DiscountedSurrenderOutgo                 float64 `json:"discounted_surrender_outgo"`
	DiscountedSurrenderOutgoAdjusted         float64 `json:"discounted_surrender_outgo_adjusted"`
	DiscountedMaturityOutgo                  float64 `json:"discounted_maturity_outgo"`
	DiscountedMaturityOutgoAdjusted          float64 `json:"discounted_maturity_outgo_adjusted"`
	DiscountedDeathOutgo                     float64 `json:"discounted_death_outgo"`
	DiscountedDeathOutgoAdjusted             float64 `json:"discounted_death_outgo_adjusted"`
	DiscountedRecoveryOutgo                  float64 `json:"discounted_recovery_outgo"`
	DiscountedRecoveryOutgoAdjusted          float64 `json:"discounted_recovery_outgo_adjusted"`
	DiscountedAnnuityOutgo                   float64 `json:"discounted_annuity_outgo"`
	DiscountedAnnuityOutgoAdjusted           float64 `json:"discounted_annuity_outgo_adjusted"`
	DiscountedRider                          float64 `json:"discounted_rider"`
	DiscountedRiderAdjusted                  float64 `json:"discounted_rider_adjusted"`
	DiscountedInitialExpenses                float64 `json:"discounted_initial_expenses"`
	DiscountedInitialExpensesAdjusted        float64 `json:"discounted_initial_expenses_adjusted"`
	DiscountedRenewalExpenses                float64 `json:"discounted_renewal_expenses"`
	DiscountedRenewalExpensesAdjusted        float64 `json:"discounted_renewal_expenses_adjusted"`
	DiscountedInvestmentIncome               float64 `json:"discounted_investment_income"`
	DiscountedInvestmentIncomeAdjusted       float64 `json:"discounted_investment_income_adjusted"`
	DiscountedProfit                         float64 `json:"discounted_profit"`
	DiscountedProfitAdjusted                 float64 `json:"discounted_profit_adjusted"`
	DiscountedCorporateTax                   float64 `json:"discounted_corporate_tax"`
	DiscountedCorporateTaxAdjusted           float64 `json:"discounted_corporate_tax_adjusted"`
	DiscountedAnnuityFactor                  float64 `json:"discounted_annuity_factor"`
	DiscountedCashOutflow                    float64 `json:"discounted_cash_outflow"`
	DiscountedCashOutflowExclAcquisition     float64 `json:"discounted_cash_outflow_excl_acquisition"`
	DiscountedAcquisitionCost                float64 `json:"discounted_acquisition_cost"`
	DiscountedCashInflow                     float64 `json:"discounted_cash_inflow"`
	SumCoverageUnits                         float64 `json:"sum_coverage_units"`
	DiscountedCoverageUnits                  float64 `json:"discounted_coverage_units"`
	CededAnnuityIncome                       float64 `json:"ceded_annuity_income"`
	LeadReAnnuityIncome                      float64 `json:"lead_re_annuity_income"`
	Re2AnnuityIncome                         float64 `json:"re_2_annuity_income"`
	Re3AnnuityIncome                         float64 `json:"re_3_annuity_income"`
	LeadReAnnuityOutgo                       float64 `json:"lead_re_annuity_outgo"`
	LeadReAnnuityOutgoAdjusted               float64 `json:"lead_re_annuity_outgo_adjusted"`
	Re2AnnuityOutgo                          float64 `json:"re_2_annuity_outgo"`
	Re2AnnuityOutgoAdjusted                  float64 `json:"re_2_annuity_outgo_adjusted"`
	Re3AnnuityOutgo                          float64 `json:"re_3_annuity_outgo"`
	Re3AnnuityOutgoAdjusted                  float64 `json:"re_3_annuity_outgo_adjusted"`
	DiscountedLeadReAnnuityOutgo             float64 `json:"discounted_lead_re_annuity_outgo"`
	DiscountedLeadReAnnuityOutgoAdjusted     float64 `json:"discounted_lead_re_annuity_outgo_adjusted"`
	DiscountedRe2AnnuityOutgo                float64 `json:"discounted_re_2_annuity_outgo"`
	DiscountedRe2AnnuityOutgoAdjusted        float64 `json:"discounted_re_2_annuity_outgo_adjusted"`
	DiscountedRe3AnnuityOutgo                float64 `json:"discounted_re_3_annuity_outgo"`
	DiscountedRe3AnnuityOutgoAdjusted        float64 `json:"discounted_re_3_annuity_outgo_adjusted"`
	NetReinsuranceCashflow                   float64 `json:"net_reinsurance_cashflow"`
	NetReinsuranceCashflowAdjusted           float64 `json:"net_reinsurance_cashflow_adjusted"`
	DiscountedNetLeadReCashflow              float64 `json:"discounted_net_lead_re_cashflow"`
	DiscountedNetLeadReCashflowAdjusted      float64 `json:"discounted_net_lead_re_cashflow_adjusted"`
	DiscountedNetRe2Cashflow                 float64 `json:"discounted_net_re_2_cashflow"`
	DiscountedNetRe2CashflowAdjusted         float64 `json:"discounted_net_re_2_cashflow_adjusted"`
	DiscountedNetRe3Cashflow                 float64 `json:"discounted_net_re_3_cashflow"`
	DiscountedNetRe3CashflowAdjusted         float64 `json:"discounted_net_re_3_cashflow_adjusted"`
	DiscountedNetReinsuranceCashflow         float64 `json:"discounted_net_reinsurance_cashflow"`
	DiscountedNetReinsuranceCashflowAdjusted float64 `json:"discounted_net_reinsurance_cashflow_adjusted"`
	RunBasis                                 string  `json:"run_basis"`
}

type PhiRunConfig struct {
	ID                uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name              string    `json:"name" gorm:"not null;uniqueIndex;size:255"`
	RunName           string    `json:"run_name" gorm:"size:255"`
	ModelPointYear    int       `json:"model_point_year"`
	ModelPointVersion string    `json:"model_point_version" gorm:"size:255"`
	ParameterYear     int       `json:"parameter_year"`
	ParameterVersion  string    `json:"parameter_version" gorm:"size:255"`
	MortalityYear     int       `json:"mortality_year"`
	MortalityVersion  string    `json:"mortality_version" gorm:"size:255"`
	RecoveryYear      int       `json:"recovery_year"`
	RecoveryVersion   string    `json:"recovery_version" gorm:"size:255"`
	YieldCurveYear    int       `json:"yield_curve_year"`
	YieldCurveVersion string    `json:"yield_curve_version" gorm:"size:255"`
	ShockSettingsID   int       `json:"shock_settings_id"`
	ShockSettingsName string    `json:"shock_settings_name" gorm:"size:255"`
	AggregationPeriod int       `json:"aggregation_period"`
	YearEndMonth      int       `json:"year_end_month"`
	RunSingle         bool      `json:"run_single"`
	CreatedBy         string    `json:"created_by" gorm:"size:255"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type PhiScopedAggregatedProjections struct {
	ID                                       int     `json:"-" gorm:"primary_key"`
	JobProductID                             int     `json:"-" gorm:"index"`
	RunDate                                  string  `json:"run_date"`
	RunId                                    int     `json:"run_id" gorm:"index"`
	RunName                                  string  `json:"run_name"`
	ProductCode                              string  `json:"product_code" gorm:"index"`
	PolicyNumber                             string  `json:"policy_number" gorm:"index"`
	SpCode                                   int     `json:"sp_code"`
	IFRS17Group                              string  `json:"ifrs_17_group"`
	ProjectionMonth                          int     `json:"projection_month" gorm:"index"`
	ProjectionYear                           int     `json:"projection_year"`
	ValuationTimeMonth                       int     `json:"valuation_time_month"`
	CalendarMonth                            int     `json:"calendar_month"`
	ValuationTimeYear                        float64 `json:"valuation_time_year"`
	AgeNextBirthday                          int     `json:"age_next_birthday"`
	InflationFactor                          float64 `json:"inflation_factor"`
	InflationFactorAdjusted                  float64 `json:"inflation_factor_adjusted"`
	AnnuityEscalationRate                    float64 `json:"annuity_escalation_rate"`
	AnnuityEscalation                        float64 `json:"annuity_escalation"`
	BaseMortalityRate                        float64 `json:"base_mortality_rate"`
	BaseMortalityRateAdjusted                float64 `json:"base_mortality_rate_adjusted"`
	BaseRecoveryRate                         float64 `json:"base_recovery_rate"`
	BaseRecoveryRateAdjusted                 float64 `json:"base_recovery_rate_adjusted"`
	IndependentMortalityRateMonthly          float64 `json:"independent_mortality_rate_monthly"`
	IndependentMortalityRateAdjustedByMonth  float64 `json:"independent_mortality_rate_adjusted_by_month"`
	IndependentRecoveryRateMonthly           float64 `json:"independent_recovery_rate_monthly"`
	IndependentRecoveryRateMonthlyAdjusted   float64 `json:"independent_recovery_rate_monthly_adjusted"`
	MonthlyDependentMortality                float64 `json:"monthly_dependent_mortality"`
	MonthlyDependentMortalityAdjusted        float64 `json:"monthly_dependent_mortality_adjusted"`
	MonthlyDependentRecovery                 float64 `json:"monthly_dependent_recovery"`
	MonthlyDependentRecoveryAdjusted         float64 `json:"monthly_dependent_recovery_adjusted"`
	InitialPolicy                            float64 `gorm:"index" json:"initial_policy"`
	InitialPolicyAdjusted                    float64 `gorm:"index" json:"initial_policy_adjusted"`
	NumberOfDeathsInForce                    float64 `json:"number_of_deaths_in_force"`
	NumberOfDeathsInForceAdjusted            float64 `json:"number_of_deaths_in_force_adjusted"`
	NumberOfRecoveries                       float64 `json:"number_of_recoveries"`
	NumberOfRecoveriesAdjusted               float64 `json:"number_of_recoveries_adjusted"`
	NumberOfMaturities                       float64 `json:"number_of_maturities"`
	NumberOfMaturitiesAdjusted               float64 `json:"number_of_maturities_adjusted"`
	IncrementalDeaths                        float64 `json:"incremental_deaths"`
	IncrementalDeathsAdjusted                float64 `json:"incremental_deaths_adjusted"`
	IncrementalRecoveries                    float64 `json:"incremental_recoveries"`
	IncrementalRecoveriesAdjusted            float64 `json:"incremental_recoveries_adjusted"`
	SumAssured                               float64 `json:"sum_assured"`
	RiderSumAssured                          float64 `json:"rider_sum_assured"`
	AnnuityIncome                            float64 `json:"annuity_income"`
	Premium                                  float64 `json:"premium"`
	MaturityValue                            float64 `json:"maturity_value"`
	PremiumIncome                            float64 `json:"premium_income"`
	PremiumIncomeAdjusted                    float64 `json:"premium_income_adjusted"`
	SurrenderOutgo                           float64 `json:"surrender_outgo"`
	SurrenderOutgoAdjusted                   float64 `json:"surrender_outgo_adjusted"`
	DeathOutgo                               float64 `json:"death_outgo"`
	DeathOutgoAdjusted                       float64 `json:"death_outgo_adjusted"`
	RecoveryOutgo                            float64 `json:"retrenchment_outgo"`
	RecoveryOutgoAdjusted                    float64 `json:"retrenchment_outgo_adjusted"`
	AnnuityOutgo                             float64 `json:"annuity_outgo"`
	AnnuityOutgoAdjusted                     float64 `json:"annuity_outgo_adjusted"`
	Rider                                    float64 `json:"rider"`
	RiderAdjusted                            float64 `json:"rider_adjusted"`
	InitialExpenses                          float64 `json:"initial_expenses"`
	InitialExpensesAdjusted                  float64 `json:"initial_expenses_adjusted"`
	RenewalExpenses                          float64 `json:"renewal_expenses"`
	RenewalExpensesAdjusted                  float64 `json:"renewal_expenses_adjusted"`
	MaturityOutgo                            float64 `json:"maturity_outgo"`
	MaturityOutgoAdjusted                    float64 `json:"maturity_outgo_adjusted"`
	NetCashFlow                              float64 `gorm:"index" json:"net_cash_flow"`
	NetCashFlowAdjusted                      float64 `gorm:"index" json:"net_cash_flow_adjusted"`
	ValuationRate                            float64 `json:"valuation_rate"`
	ValuationRateAdjusted                    float64 `json:"valuation_rate_adjusted"`
	Reserves                                 float64 `gorm:"index" json:"reserves"`
	ReservesAdjusted                         float64 `gorm:"index" json:"reserves_adjusted"`
	ChangeInReserves                         float64 `gorm:"index" json:"change_in_reserves"`
	ChangeInReservesAdjusted                 float64 `gorm:"index" json:"change_in_reserves_adjusted"`
	InvestmentIncome                         float64 `json:"investment_income"`
	InvestmentIncomeAdjusted                 float64 `json:"investment_income_adjusted"`
	Profit                                   float64 `json:"profit"`
	ProfitAdjusted                           float64 `json:"profit_adjusted"`
	RiskDiscountRate                         float64 `json:"risk_discount_rate"`
	RiskDiscountRateAdjusted                 float64 `json:"risk_discount_rate_adjusted"`
	VIF                                      float64 `json:"vif"`
	VIFAdjusted                              float64 `json:"vif_adjusted"`
	CorporateTax                             float64 `json:"corporate_tax"`
	CorporateTaxAdjusted                     float64 `json:"corporate_tax_adjusted"`
	CoverageUnits                            float64 `json:"coverage_units"`
	DiscountedPremiumIncome                  float64 `json:"discounted_premium_income"`
	DiscountedPremiumIncomeAdjusted          float64 `json:"discounted_premium_income_adjusted"`
	DiscountedSurrenderOutgo                 float64 `json:"discounted_surrender_outgo"`
	DiscountedSurrenderOutgoAdjusted         float64 `json:"discounted_surrender_outgo_adjusted"`
	DiscountedMaturityOutgo                  float64 `json:"discounted_maturity_outgo"`
	DiscountedMaturityOutgoAdjusted          float64 `json:"discounted_maturity_outgo_adjusted"`
	DiscountedDeathOutgo                     float64 `json:"discounted_death_outgo"`
	DiscountedDeathOutgoAdjusted             float64 `json:"discounted_death_outgo_adjusted"`
	DiscountedRecoveryOutgo                  float64 `json:"discounted_recovery_outgo"`
	DiscountedRecoveryOutgoAdjusted          float64 `json:"discounted_recovery_outgo_adjusted"`
	DiscountedAnnuityOutgo                   float64 `json:"discounted_annuity_outgo"`
	DiscountedAnnuityOutgoAdjusted           float64 `json:"discounted_annuity_outgo_adjusted"`
	DiscountedRider                          float64 `json:"discounted_rider"`
	DiscountedRiderAdjusted                  float64 `json:"discounted_rider_adjusted"`
	DiscountedInitialExpenses                float64 `json:"discounted_initial_expenses"`
	DiscountedInitialExpensesAdjusted        float64 `json:"discounted_initial_expenses_adjusted"`
	DiscountedRenewalExpenses                float64 `json:"discounted_renewal_expenses"`
	DiscountedRenewalExpensesAdjusted        float64 `json:"discounted_renewal_expenses_adjusted"`
	DiscountedInvestmentIncome               float64 `json:"discounted_investment_income"`
	DiscountedInvestmentIncomeAdjusted       float64 `json:"discounted_investment_income_adjusted"`
	DiscountedProfit                         float64 `json:"discounted_profit"`
	DiscountedProfitAdjusted                 float64 `json:"discounted_profit_adjusted"`
	DiscountedCorporateTax                   float64 `json:"discounted_corporate_tax"`
	DiscountedCorporateTaxAdjusted           float64 `json:"discounted_corporate_tax_adjusted"`
	DiscountedAnnuityFactor                  float64 `json:"discounted_annuity_factor"`
	DiscountedCashOutflow                    float64 `json:"discounted_cash_outflow"`
	DiscountedCashOutflowExclAcquisition     float64 `json:"discounted_cash_outflow_excl_acquisition"`
	DiscountedAcquisitionCost                float64 `json:"discounted_acquisition_cost"`
	DiscountedCashInflow                     float64 `json:"discounted_cash_inflow"`
	SumCoverageUnits                         float64 `json:"sum_coverage_units"`
	DiscountedCoverageUnits                  float64 `json:"discounted_coverage_units"`
	CededAnnuityIncome                       float64 `json:"ceded_annuity_income"`
	LeadReAnnuityIncome                      float64 `json:"lead_re_annuity_income"`
	Re2AnnuityIncome                         float64 `json:"re_2_annuity_income"`
	Re3AnnuityIncome                         float64 `json:"re_3_annuity_income"`
	LeadReAnnuityOutgo                       float64 `json:"lead_re_annuity_outgo"`
	LeadReAnnuityOutgoAdjusted               float64 `json:"lead_re_annuity_outgo_adjusted"`
	Re2AnnuityOutgo                          float64 `json:"re_2_annuity_outgo"`
	Re2AnnuityOutgoAdjusted                  float64 `json:"re_2_annuity_outgo_adjusted"`
	Re3AnnuityOutgo                          float64 `json:"re_3_annuity_outgo"`
	Re3AnnuityOutgoAdjusted                  float64 `json:"re_3_annuity_outgo_adjusted"`
	DiscountedLeadReAnnuityOutgo             float64 `json:"discounted_lead_re_annuity_outgo"`
	DiscountedLeadReAnnuityOutgoAdjusted     float64 `json:"discounted_lead_re_annuity_outgo_adjusted"`
	DiscountedRe2AnnuityOutgo                float64 `json:"discounted_re_2_annuity_outgo"`
	DiscountedRe2AnnuityOutgoAdjusted        float64 `json:"discounted_re_2_annuity_outgo_adjusted"`
	DiscountedRe3AnnuityOutgo                float64 `json:"discounted_re_3_annuity_outgo"`
	DiscountedRe3AnnuityOutgoAdjusted        float64 `json:"discounted_re_3_annuity_outgo_adjusted"`
	NetReinsuranceCashflow                   float64 `json:"net_reinsurance_cashflow"`
	NetReinsuranceCashflowAdjusted           float64 `json:"net_reinsurance_cashflow_adjusted"`
	DiscountedNetLeadReCashflow              float64 `json:"discounted_net_lead_re_cashflow"`
	DiscountedNetLeadReCashflowAdjusted      float64 `json:"discounted_net_lead_re_cashflow_adjusted"`
	DiscountedNetRe2Cashflow                 float64 `json:"discounted_net_re_2_cashflow"`
	DiscountedNetRe2CashflowAdjusted         float64 `json:"discounted_net_re_2_cashflow_adjusted"`
	DiscountedNetRe3Cashflow                 float64 `json:"discounted_net_re_3_cashflow"`
	DiscountedNetRe3CashflowAdjusted         float64 `json:"discounted_net_re_3_cashflow_adjusted"`
	DiscountedNetReinsuranceCashflow         float64 `json:"discounted_net_reinsurance_cashflow"`
	DiscountedNetReinsuranceCashflowAdjusted float64 `json:"discounted_net_reinsurance_cashflow_adjusted"`
	RunBasis                                 string  `json:"run_basis"`
}
