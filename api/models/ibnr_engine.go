package models

import (
	"time"
)

type IbnrYieldCurve struct {
	Year           int     `json:"year" csv:"year" gorm:"primary_key;auto_increment:false"`
	YieldCurveCode string  `json:"yield_curve_code" csv:"yield_curve_code" gorm:"primary_key;auto_increment:false;column:yield_curve_code"`
	ProjectionTime int     `json:"proj_time" csv:"proj_time" gorm:"primary_key;auto_increment:false;column:proj_time"`
	Month          int     `json:"month" csv:"month" gorm:"primary_key;auto_increment:false;column:month"`
	NominalRate    float64 `json:"nominal_rate" csv:"nominal_rate"`
	Inflation      float64 `json:"inflation" csv:"inflation"`
	//Basis          string  `json:"basis" csv:"basis"`
	Created int64 `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type IBNRShockSetting struct {
	ID            int    `json:"id" gorm:"primary_key"`
	Name          string `json:"name" gorm:"unique"`
	Description   string `json:"description"`
	DiscountCurve bool   `json:"discount_curve"`
	Claims        bool   `json:"claims" csv:"claims"`
	Expenses      bool   `json:"expenses" csv:"expenses"`
	ShockBasis    string `json:"shock_basis" csv:"shock_basis"`
	Year          int    `json:"year" csv:"year"`
}

type IBNRIncurredClaims struct {
	ID                  int     `json:"-" gorm:"primary_key"`
	ProductCode         string  `json:"-" csv:"product_code"`
	AccidentYear        int     `json:"accident_year" csv:"accident_year"`
	TotalClaimsIncurred float64 `json:"total_claims_incurred" csv:"total_claims_incurred"`
	ClaimCount          int     `json:"claim_count" csv:"claim_count"`
}

type IBNRPaidVsOutstandingClaims struct {
	ID                     int     `json:"-" gorm:"primary_key"`
	ProductCode            string  `json:"product_code" csv:"product_code"`
	AccidentYear           int     `json:"accident_year" csv:"accident_year"`
	TotalClaimsPaid        float64 `json:"total_claims_paid" csv:"total_claims_paid"`
	TotalClaimsOutstanding float64 `json:"total_claims_outstanding" csv:"total_claims_outstanding"`
}

type IBNRProportionOutstandingClaims struct {
	ID                          int     `json:"-" gorm:"primary_key"`
	ProductCode                 string  `json:"product_code" csv:"product_code"`
	AccidentYear                int     `json:"accident_year" csv:"accident_year"`
	PercentageOutstandingClaims float64 `json:"percentage_outstanding_claims" csv:"percentage_outstanding_claims"`
}

type IBNRAverageClaimAmount struct {
	ID                 int     `json:"-" gorm:"primary_key"`
	ProductCode        string  `json:"product_code" csv:"product_code"`
	AccidentYear       int     `json:"accident_year" csv:"accident_year"`
	AverageClaimAmount float64 `json:"average_claim_amount" csv:"average_claim_amount"`
}

type IBNRClaimHistorySummary struct {
	ID                             int     `json:"-" gorm:"primary_key"`
	ProductCode                    string  `json:"product_code" csv:"product_code"`
	AccidentYear                   int     `json:"accident_year" csv:"accident_year"`
	CurrentTotalIncurredClaims     float64 `json:"current_total_incurred_claims" csv:"current_total_incurred_claims"`
	PreviousTotalIncurredClaims    float64 `json:"previous_total_incurred_claims" csv:"previous_total_incurred_claims"`
	CurrentTotalPaidClaims         float64 `json:"current_total_paid_claims" csv:"current_total_paid_claims"`
	PreviousTotalPaidClaims        float64 `json:"previous_total_paid_claims" csv:"previous_total_paid_claims"`
	CurrentTotalOutstandingClaims  float64 `json:"current_total_outstanding_claims" csv:"current_total_outstanding_claims"`
	PreviousTotalOutstandingClaims float64 `json:"previous_total_outstanding_claims" csv:"previous_total_outstanding_claims"`
}

// IBNRTableWithChart wraps any summary table with an embedded chart configuration
// so the frontend can auto-render the appropriate chart without hardcoded logic.
type IBNRTableWithChart struct {
	TableName   string                 `json:"table_name"`
	ChartConfig map[string]interface{} `json:"chart_config"`
	Data        interface{}            `json:"data"`
}

// IBNRLossRatioSummary — Loss Ratio by Accident Year (CAS / BF prerequisite)
type IBNRLossRatioSummary struct {
	AccidentYear      int     `json:"accident_year"`
	EarnedPremium     float64 `json:"earned_premium"`
	TotalIncurred     float64 `json:"total_incurred"`
	ActualLossRatio   float64 `json:"actual_loss_ratio"`
	ExpectedLossRatio float64 `json:"expected_loss_ratio"`
	UltimateLossRatio float64 `json:"ultimate_loss_ratio"`
}

// IBNRUltimateLossSummary — Ultimate Loss Reconciliation by Accident Year (IFoA GIRO)
type IBNRUltimateLossSummary struct {
	AccidentYear       int     `json:"accident_year"`
	PaidToDate         float64 `json:"paid_to_date"`
	OutstandingAtVal   float64 `json:"outstanding_at_valuation"`
	IbnrBestEstimate   float64 `json:"ibnr_best_estimate"`
	UltimateLoss       float64 `json:"ultimate_loss"`
	PctPaidToUltimate  float64 `json:"pct_paid_to_ultimate"`
	ProportionRunOff   float64 `json:"proportion_run_off"`
	EffectiveMethod    string  `json:"effective_method"`
}

// IBNRFrequencySeveritySummary — Claim Frequency & Severity Trend (Munich / Clark LDF)
type IBNRFrequencySeveritySummary struct {
	AccidentYear    int     `json:"accident_year"`
	EarnedPremium   float64 `json:"earned_premium"`
	ClaimCount      int     `json:"claim_count"`
	TotalIncurred   float64 `json:"total_incurred"`
	ClaimFrequency  float64 `json:"claim_frequency"`
	ClaimSeverity   float64 `json:"claim_severity"`
}

// IBNRDevelopmentMaturitySummary — Development Maturity by Accident Year (ASISA SAP 5)
type IBNRDevelopmentMaturitySummary struct {
	AccidentYear            int     `json:"accident_year"`
	ProportionRunOff        float64 `json:"proportion_run_off"`
	ProportionNotRunOff     float64 `json:"proportion_not_run_off"`
	EffectiveMethod         string  `json:"effective_method"`
	ChainLadderEligible     bool    `json:"chain_ladder_eligible"`
}

// IBNRReserveAdequacySummary — Actual vs Expected Loss Diagnostics (ASISA SAP 5)
type IBNRReserveAdequacySummary struct {
	AccidentYear          int     `json:"accident_year"`
	ActualClaims          float64 `json:"actual_claims"`
	PredictedTotalLoss    float64 `json:"predicted_total_loss"`
	ActualVsExpectedRatio float64 `json:"actual_vs_expected_ratio"`
	PredictedLossRatio    float64 `json:"predicted_loss_ratio"`
	CredibilityWeightCL   float64 `json:"credibility_weight_cl"`
	CredibilityWeightBF   float64 `json:"credibility_weight_bf"`
}

// IBNRMethodSensitivitySummary — IBNR Method Comparison as % deviation from BEL (actuarial report)
type IBNRMethodSensitivitySummary struct {
	ProductCode          string  `json:"product_code"`
	IbnrBel              float64 `json:"ibnr_bel"`
	ChainLadderIbnr      float64 `json:"chain_ladder_ibnr"`
	ClVsBelPct           float64 `json:"cl_vs_bel_pct"`
	BfIbnr               float64 `json:"bf_ibnr"`
	BfVsBelPct           float64 `json:"bf_vs_bel_pct"`
	CapeCodIbnr          float64 `json:"cape_cod_ibnr"`
	CcVsBelPct           float64 `json:"cc_vs_bel_pct"`
	CombinedClBfIbnr     float64 `json:"combined_cl_bf_ibnr"`
	CombinedClCcIbnr     float64 `json:"combined_cl_cape_cod_ibnr"`
	BootstrapIbnr        float64 `json:"bootstrap_ibnr"`
	BootstrapVsBelPct    float64 `json:"bootstrap_vs_bel_pct"`
	MackModelIbnr        float64 `json:"mack_model_ibnr"`
	MackVsBelPct         float64 `json:"mack_vs_bel_pct"`
}

// IBNRStochasticSummary — Stochastic vs Deterministic Reserve Summary (IFRS 17 para 37)
type IBNRStochasticSummary struct {
	ProductCode               string  `json:"product_code"`
	IbnrBel                   float64 `json:"ibnr_bel"`
	BootstrapMean             float64 `json:"bootstrap_mean"`
	BootstrapStdDev           float64 `json:"bootstrap_std_dev"`
	BootstrapCV               float64 `json:"bootstrap_cv"`
	BootstrapNthPercentile    float64 `json:"bootstrap_nth_percentile"`
	MackMean                  float64 `json:"mack_mean"`
	MackStdDev                float64 `json:"mack_std_dev"`
	MackCV                    float64 `json:"mack_cv"`
	MackNthPercentile         float64 `json:"mack_nth_percentile"`
	StochasticAdequacyMargin  float64 `json:"stochastic_adequacy_margin"`
}

type IBNRShock struct {
	ID                          int     `json:"id" gorm:"primary_key"`
	Year                        int     `json:"year" csv:"year"`
	ShockBasis                  string  `json:"shock_basis" csv:"shock_basis"`
	MultiplicativeClaims        float64 `json:"multiplicative_claims" csv:"multiplicative_claims"`
	AdditiveClaims              float64 `json:"additive_claims" csv:"additive_claims"`
	MultiplicativeDiscountCurve float64 `json:"multiplicative_discount_curve" csv:"multiplicative_discount_curve"`
	AdditiveDiscountCurve       float64 `json:"additive_discount_curve" csv:"additive_discount_curve"`
	MultiplicativeExpenses      float64 `json:"multiplicative_expenses" csv:"multiplicative_expenses"`
	AdditiveExpenses            float64 `json:"additive_expenses" csv:"additive_expenses"`
	Created                     int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type LICParameter struct {
	ID                                 int     `json:"id" gorm:"primary_key"`
	Year                               int     `json:"year" csv:"year"`
	PortfolioName                      string  `json:"portfolio_name" csv:"portfolio_name"`
	ProductCode                        string  `json:"product_code" csv:"product_code"`
	ProductName                        string  `json:"product_name" csv:"product_name"`
	UnderwritingYear                   int     `json:"underwriting_year" csv:"underwriting_year"`
	UnderwritingMonth                  int     `json:"underwriting_month" csv:"underwriting_month"`
	YieldCurveCode                     string  `json:"yield_curve_code" csv:"yield_curve_code"`
	Basis                              string  `json:"basis" csv:"basis"`
	ExpectedLossRatio                  float64 `json:"expected_loss_ratio" csv:"expected_loss_ratio"`
	Margin                             float64 `json:"margin" csv:"margin"`
	RiskAdjustmentConfidenceLevel      float64 `json:"risk_adjustment_confidence_level" csv:"risk_adjustment_confidence_level"`
	RiskAdjustmentFactor               float64 `json:"risk_adjustment_factor" csv:"risk_adjustment_factor"`
	AllocatedClaimsExpenseProportion   float64 `json:"allocated_claims_expense_proportion" csv:"allocated_claims_expense_proportion"`
	AllocatedClaimsExpenseAmount       float64 `json:"allocated_claims_expense_amount" csv:"allocated_claims_expense_amount"`
	UnallocatedClaimsExpenseProportion float64 `json:"unallocated_claims_expense_proportion" csv:"unallocated_claims_expense_proportion"`
	UnallocatedClaimsExpenseAmount     float64 `json:"unallocated_claims_expense_amount" csv:"unallocated_claims_expense_amount"`
	Treaty1RecoverableProportion       float64 `json:"treaty_1_recoverable_proportion" csv:"treaty_1_recoverable_proportion"`
	Treaty2RecoverableProportion       float64 `json:"treaty_2_recoverable_proportion" csv:"treaty_2_recoverable_proportion"`
	Treaty3RecoverableProportion       float64 `json:"treaty_3_recoverable_proportion" csv:"treaty_3_recoverable_proportion"`
	Created                            int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type LicClaimsYearVersion struct {
	ID            int    `json:"id" csv:"id" gorm:"primary_key"`
	PortfolioName string `json:"portfolio_name" csv:"portfolio_name"`
	PortfolioId   int    `json:"portfolio_id"`
	Year          int    `json:"year" csv:"year"`
	VersionName   string `json:"version_name" csv:"version_name"`
	Count         int    `json:"count" csv:"count"`
}

type LicEarnedPremiumYearVersion struct {
	ID            int    `json:"id" csv:"id" gorm:"primary_key"`
	PortfolioName string `json:"portfolio_name" csv:"portfolio_name"`
	PortfolioId   int    `json:"portfolio_id"`
	Year          int    `json:"year" csv:"year"`
	VersionName   string `json:"version_name" csv:"version_name"`
	Count         int    `json:"count" csv:"count"`
}

type LICClaimsInput struct {
	ID               int     `json:"id" gorm:"primary_key" csv:"-"`
	Year             int     `json:"year" csv:"year"`
	VersionName      string  `json:"version_name" csv:"version_name"`
	ClaimNumber      string  `json:"claim_number" csv:"claim_number"`
	LicPortfolioID   int     `json:"lic_portfolio_id" csv:"lic_portfolio_id"`
	PortfolioName    string  `json:"portfolio_name" csv:"portfolio_name"`
	ProductName      string  `json:"product_name" csv:"product_name"`
	ProductCode      string  `json:"product_code" csv:"product_code"`
	PolicyNumber     string  `json:"policy_number" csv:"policy_number"`
	DamageYear       int     `json:"damage_year" csv:"damage_year"`
	DamageMonth      int     `json:"damage_month" csv:"damage_month"`
	ReportedYear     int     `json:"reported_year" csv:"reported_year"`
	ReportedMonth    int     `json:"reported_month" csv:"reported_month"`
	LockedInYear     int     `json:"locked_in_year" csv:"locked_in_year"`
	LockedInMonth    int     `json:"locked_in_month" csv:"locked_in_month"`
	SettlementYear   int     `json:"settlement_year" csv:"settlement_year"`
	SettlementMonth  int     `json:"settlement_month" csv:"settlement_month"`
	ClaimAmount      float64 `json:"claim_amount" csv:"claim_amount"`
	AssessmentCost   float64 `json:"assessment_cost" csv:"assessment_cost"`
	PaidClaims       float64 `json:"paid_claims" csv:"paid_claims"`
	IFRS17Group      string  `json:"ifrs17_group" csv:"ifrs17_group"`
	CauseDamage      string  `json:"cause_damage" csv:"cause_damage"`
	ClaimStatus      string  `json:"claim_status" csv:"claim_status"`
	CoverType        string  `json:"cover_type" csv:"cover_type"`
	Deductible       string  `json:"deductible" csv:"deductible"`
	Location         string  `json:"location" csv:"location"`
	UnderwritingYear int     `json:"underwriting_year" csv:"underwriting_year"`
	Created          int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type LICEarnedPremium struct {
	ID              int     `json:"id" gorm:"primary_key"`
	PortfolioName   string  `json:"portfolio_name" csv:"portfolio_name"`
	IBNRPortfolioID int     `json:"ibnr_portfolio_id" csv:"ibnr_portfolio_id"`
	ProductCode     string  `json:"product_code" csv:"product_code"`
	YearIndex       int     `json:"year_index" csv:"year_index"`
	Month           int     `json:"month" csv:"month"`
	EarnedPremium   float64 `json:"earned_premium" csv:"earned_premium"`
	Year            int     `json:"year" csv:"year"`
	VersionName     string  `json:"version_name" csv:"version_name"`
	Created         int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

func (LICEarnedPremium) TableName() string {
	return "lic_earned_premiums"
}

type LicCPI struct {
	ID         int     `json:"id" gorm:"primary_key"`
	YearIndex  int     `json:"year_index" csv:"year_index"`
	MonthIndex int     `json:"month_index" csv:"month_index"`
	CpiIndex   float64 `json:"cpi_index" csv:"cpi_index"`
	Created    int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type LicBaseVariable struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
	//Description     string `json:"description"`
	RunName         string `json:"run_name"`
	RunId           int    `json:"run_id"`
	Notes           string `json:"notes"`
	AssumptionBasis string `json:"assumption_basis"`
}

type LicVariable struct {
	ID                int    `json:"id" gorm:"primary_key"`
	LicVariableSetID  int    `json:"lic_variable_set_id"`
	LicBaseVariableID string `json:"lic_base_variable_id"`
	Name              string `json:"name"`
	//Description       string `json:"description"`
	RunName         string `json:"run_name"`
	RunId           int    `json:"run_id"`
	Notes           string `json:"notes"`
	AssumptionBasis string `json:"assumption_basis"`
}

type LicBuildupVariable struct {
	ID                int    `json:"id" gorm:"primary_key"`
	LicVariableSetID  int    `json:"lic_variable_set_id"`
	LicBaseVariableID string `json:"lic_base_variable_id"`
	Name              string `json:"name"`
	//Description       string `json:"description"`
	RunName         string `json:"run_name"`
	RunId           int    `json:"run_id"`
	Notes           string `json:"notes"`
	AssumptionBasis string `json:"assumption_basis"`
}

type LicBuildupResult struct {
	ID                               int     `json:"id" gorm:"primary_key"`
	ConfigurationId                  int     `json:"configuration_id"`
	PortfolioId                      int     `json:"portfolio_id"`
	PortfolioName                    string  `json:"portfolio_name"`
	ProductCode                      string  `json:"product_code"`
	ConfigurationName                string  `json:"configuration_name"`
	RunDate                          string  `json:"run_date"`
	RunId                            int     `json:"run_id"`
	Name                             string  `json:"name"`
	Description                      string  `json:"description"`
	IBNR                             float64 `json:"ibnr"`
	ReportedClaims                   float64 `json:"reported_claims"`
	CashBack                         float64 `json:"cash_back"`
	UnadjustedLossAdjustmentExpenses float64 `json:"unadjusted_loss_adjustment_expenses"`
	RiskAdjustment                   float64 `json:"risk_adjustment"`
	LicBuildup                       float64 `json:"lic_buildup"`
	VariableChange                   float64 `json:"variable_change"`
	Pnl                              float64 `json:"Pnl"`
	Treaty1ReportedClaims            float64 `json:"treaty1_reported_claims"`
	Treaty1Ibnr                      float64 `json:"treaty1_ibnr"`
	Treaty1IbnrRiskAdjustment        float64 `json:"treaty1_ibnr_risk_adjustment"`
	Treaty2ReportedClaims            float64 `json:"treaty2_reported_claims"`
	Treaty2Ibnr                      float64 `json:"treaty2_ibnr"`
	Treaty2IbnrRiskAdjustment        float64 `json:"treaty2_ibnr_risk_adjustment"`
	Treaty3ReportedClaims            float64 `json:"treaty3_reported_claims"`
	Treaty3Ibnr                      float64 `json:"treaty3_ibnr"`
	Treaty3IbnrRiskAdjustment        float64 `json:"treaty3_ibnr_risk_adjustment"`
	Notes                            string  `json:"notes"`
	AssumptionBasis                  string  `json:"assumption_basis"`
	IBNRAt12                         float64 `json:"ibnr_at12"`
	IBNRRiskAdjustmentAt12           float64 `json:"ibnr_risk_adjustment_at12"`
	IFRS17Group                      string  `json:"ifrs17_group"`
}

type LicVariableSet struct {
	ID                int           `json:"id" gorm:"primary_key"`
	ConfigurationName string        `json:"configuration_name"`
	RunType           string        `json:"run_type"`
	LicVariables      []LicVariable `json:"lic_variables"`
}

type LicPortfolio struct {
	ID                       int                           `json:"id" gorm:"primary_key"`
	Name                     string                        `json:"name" gorm:"unique"`
	DiscountOption           string                        `json:"discount_option"`
	PolicyCount              []LicPolicyCount              `json:"policy_count" gorm:"-"`
	EarnedPremiumCount       []LicPolicyCount              `json:"earned_premium_count" gorm:"-"`
	ClaimsYearVersion        []LicClaimsYearVersion        `json:"claims_year_version" gorm:"-"`
	EarnedPremiumYearVersion []LicEarnedPremiumYearVersion `json:"earned_premium_year_version" gorm:"-"`
	ClaimsInput              []LICClaimsInput              `json:"claims_input"`
}
type LicPolicyCount struct {
	Year  int `json:"year"`
	Count int `json:"count"`
}

type PAAModelPointCount struct {
	Year      int    `json:"year"`
	MpVersion string `json:"mp_version"`
	Count     int    `json:"count"`
}

type LicModelPoint struct { // calculated values from claims input
	ID int `json:"id" gorm:"primary_key"`
	//Year                 int     `json:"year" csv:"year"`
	RunDate              string  `json:"run_date" csv:"run_date"`
	RunID                int     `json:"run_id" csv:"run_id"`
	ClaimNumber          string  `json:"claim_number"`
	LicPortfolioID       int     `json:"lic_portfolio_id"`
	PortfolioName        string  `json:"portfolio_name"`
	ProductName          string  `json:"product_name"`
	ProductCode          string  `json:"product_code"`
	IFRS17Group          string  `json:"ifrs17_group"`
	PolicyNumber         string  `json:"policy_number"`
	CPIDenominator       float64 `json:"cpi_denominator"`
	CPINumerator         float64 `json:"cpi_numerator"`
	UnInflatedClaim      float64 `json:"un_inflated_claim"`
	PaidClaims           float64 `json:"paid_claims"`
	AssessmentCost       float64 `json:"assessment_cost"`
	TotalUnInflatedClaim float64 `json:"total_un_inflated_claim"`
	TotalInflatedClaim   float64 `json:"total_inflated_claim"`
	AccidentYear         int     `json:"accident_year"`
	AccidentMonth        int     `json:"incident_month"`
	ReportingYear        int     `json:"reporting_year"`
	ReportingMonth       int     `json:"reporting_month"`
	ReportingDelay       int     `json:"reporting_delay"`
	SettlementYear       int     `json:"settlement_year"`
	SettlementMonth      int     `json:"settlement_month"`
	EarnedPremium        float64 `json:"earned_premium"`
	UWYear               int     `json:"uw_year"`
	VersionName          string  `json:"version_name"`
}

type LicTriangulation struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"-"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicTriangulationClaimCount struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicTriangulationAverageClaim struct {
	ID             int     `json:"id" gorm:"primary_key"`
	RunDate        string  `json:"run_date,omitempty"`
	RunID          int     `json:"run_id"`
	PortfolioName  string  `json:"portfolio_name,omitempty"`
	LicPortfolioId int     `json:"lic_portfolio_id,omitempty"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicCumulativeTriangulation struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicCumulativeTriangulationClaimCount struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicCumulativeTriangulationAverageClaim struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type IBNRRunSetting struct {
	ID                   int       `json:"id" gorm:"primary_key"`
	RunName              string    `json:"run_name" csv:"run_name"`
	RunDate              string    `json:"run_date" csv:"run_date"`
	RunType              string    `json:"run_type" csv:"run_type"`
	PortfolioName        string    `json:"portfolio_name" csv:"portfolio_name"`
	PortfolioId          int       `json:"portfolio_id"`
	Description          string    `json:"description" csv:"description"`
	ClaimsDataYear       int       `json:"claims_data_year" csv:"model_point"`
	ClaimsInputVersion   string    `json:"claims_input_version" csv:"claims_input_version"`
	ParameterYear        int       `json:"parameter_year" csv:"parameter_year"`
	YieldCurveYear       int       `json:"yield_curve_year" csv:"yield_curve_year"`
	YieldCurveMonth      int       `json:"yield_curve_month" csv:"yield_curve_month"`
	InflationIndicator   bool      `json:"inflation_indicator" csv:"inflation_indicator"`
	BootStrapIndicator   bool      `json:"bootstrap_indicator" csv:"bootstrap_indicator"`
	Simulations          int       `json:"simulations" csv:"simulations"`
	MackModelIndicator   bool      `json:"mack_model_indicator" csv:"mack_model_indicator"`
	DistributionModel    string    `json:"distribution_model" csv:"distribution_model"`
	MackModelSimulations int       `json:"mack_model_simulations" csv:"mack_model_simulations"`
	CalculationInterval  string    `json:"calculation_interval" csv:"calculation_interval"`
	InflationDate        string    `json:"inflation_date" csv:"inflation_date"`
	DataInputStartDate   string    `json:"data_input_start_date" csv:"data_input_start_date"`
	DataInputEndDate     string    `json:"data_input_end_date" csv:"data_input_end_date"`
	IBNRMethod           string    `json:"ibnr_method" csv:"ibnr_method"`
	Basis                string    `json:"basis" csv:"basis"`
	TotalRecords         int       `json:"total_records" csv:"total_records"`
	ProcessedRecords     int       `json:"processed_records" csv:"processed_records"`
	ProcessingStatus     string    `json:"processing_status" csv:"processing_status"`
	FailureReason        string    `json:"failure_reason" csv:"failure_reason"`
	RunTime              float64   `json:"run_time" csv:"run_time"`
	RerunIndicator       bool      `json:"rerun_indicator" csv:"rerun_indicator"`
	CreationDate         time.Time `json:"creation_date"`
	UserName             string    `json:"user_name" csv:"user_name"`
	UserEmail            string    `json:"user_email" csv:"user_email"`
}

// LicIbnrMethodAssignment stores per-product IBNR method overrides for a portfolio.
// When no record exists for a product, the run-level IBNRMethod is used (preserving
// the existing single-method behaviour).
type LicIbnrMethodAssignment struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	PortfolioID      int       `json:"portfolio_id"`
	PortfolioName    string    `json:"portfolio_name"`
	ProductCode      string    `json:"product_code"`
	// AccidentYearFrom / AccidentYearTo define an optional AY range for within-triangle segmentation.
	// Both 0 = product-level assignment (applies to all accident years in this product).
	// Non-zero = applies only when accidentYear ∈ [AccidentYearFrom, AccidentYearTo].
	AccidentYearFrom int       `json:"accident_year_from"`
	AccidentYearTo   int       `json:"accident_year_to"`
	Method           string    `json:"ibnr_method"`
	AssignmentType  string    `json:"assignment_type"` // "manual" | "rule"
	// Rule thresholds captured at assignment time for auditability
	MinDataYears    int       `json:"min_data_years"`
	MaxCVThreshold  float64   `json:"max_cv_threshold"`
	// Computed data-quality signals saved alongside rule assignments
	ActualDataYears int       `json:"actual_data_years"`
	ActualCV        float64   `json:"actual_cv"`
	RuleReason      string    `json:"rule_reason"` // plain-English justification
	// Audit trail
	CreatedBy       string    `json:"created_by"`
	UserName        string    `json:"user_name"`
	// IsActive: false means the assignment is ignored by the engine (soft disable)
	IsActive        bool      `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type LicDevelopmentVariables struct {
	ColumnSum                        float64 `json:"column_sum"`
	LastValue                        float64 `json:"last_value"`
	ColumnSumExclLastValue           float64 `json:"column_sum_excl_last_value"`
	DevelopmentFactors               float64 `json:"development_factors"`
	CumulativeDevelopmentFactors     float64 `json:"cumulative_development_factors"`
	CumulativeDevelopmentProportion  float64 `json:"cumulative_development_proportion"`
	IncrementalDevelopmentProportion float64 `json:"incremental_development_proportion"`
	Inverse                          float64 `json:"inverse"`
	Runoff                           float64 `json:"runoff"`
	ProportionRunoff                 float64 `json:"proportion_runoff"`
	ProportionNotRunoff              float64 `json:"proportion_not_runoff"`
}

type LicDevelopmentFactor struct {
	ID                  int     `json:"-" gorm:"primary_key" csv:"-"`
	RunDate             string  `json:"-"  csv:"run_date"`
	RunID               int     `json:"-" csv:"run_id"`
	PortfolioName       string  `json:"-" csv:"portfolio_name"`
	ProductCode         string  `json:"product_code,omitempty" csv:"product_code"`
	DevelopmentVariable string  `json:"development_variable" csv:"development_variable"`
	Rd0                 float64 `json:"rd0" csv:"rd0"`
	Rd1                 float64 `json:"rd1" csv:"rd1"`
	Rd2                 float64 `json:"rd2" csv:"rd2"`
	Rd3                 float64 `json:"rd3" csv:"rd3"`
	Rd4                 float64 `json:"rd4" csv:"rd4"`
	Rd5                 float64 `json:"rd5" csv:"rd5"`
	Rd6                 float64 `json:"rd6" csv:"rd6"`
	Rd7                 float64 `json:"rd7" csv:"rd7"`
	Rd8                 float64 `json:"rd8" csv:"rd8"`
	Rd9                 float64 `json:"rd9" csv:"rd9"`
	Rd10                float64 `json:"rd10" csv:"rd10"`
	Rd11                float64 `json:"rd11" csv:"rd11"`
	Rd12                float64 `json:"rd12" csv:"rd12"`
	Rd13                float64 `json:"rd13" csv:"rd13"`
	Rd14                float64 `json:"rd14" csv:"rd14"`
	Rd15                float64 `json:"rd15" csv:"rd15"`
	Rd16                float64 `json:"rd16" csv:"rd16"`
	Rd17                float64 `json:"rd17" csv:"rd17"`
	Rd18                float64 `json:"rd18" csv:"rd18"`
	Rd19                float64 `json:"rd19" csv:"rd19"`
	Rd20                float64 `json:"rd20" csv:"rd20"`
	Rd21                float64 `json:"rd21" csv:"rd21"`
	Rd22                float64 `json:"rd22" csv:"rd22"`
	Rd23                float64 `json:"rd23" csv:"rd23"`
	Rd24                float64 `json:"rd24" csv:"rd24"`
	Rd25                float64 `json:"rd25" csv:"rd25"`
	Rd26                float64 `json:"rd26" csv:"rd26"`
	Rd27                float64 `json:"rd27" csv:"rd27"`
	Rd28                float64 `json:"rd28" csv:"rd28"`
	Rd29                float64 `json:"rd29" csv:"rd29"`
	Rd30                float64 `json:"rd30" csv:"rd30"`
	Rd31                float64 `json:"rd31" csv:"rd31"`
	Rd32                float64 `json:"rd32" csv:"rd32"`
	Rd33                float64 `json:"rd33" csv:"rd33"`
	Rd34                float64 `json:"rd34" csv:"rd34"`
	Rd35                float64 `json:"rd35" csv:"rd35"`
	Rd36                float64 `json:"rd36" csv:"rd36"`
	Rd37                float64 `json:"rd37" csv:"rd37"`
	Rd38                float64 `json:"rd38" csv:"rd38"`
	Rd39                float64 `json:"rd39" csv:"rd39"`
	Rd40                float64 `json:"rd40" csv:"rd40"`
	Rd41                float64 `json:"rd41" csv:"rd41"`
	Rd42                float64 `json:"rd42" csv:"rd42"`
	Rd43                float64 `json:"rd43" csv:"rd43"`
	Rd44                float64 `json:"rd44" csv:"rd44"`
	Rd45                float64 `json:"rd45" csv:"rd45"`
	Rd46                float64 `json:"rd46" csv:"rd46"`
	Rd47                float64 `json:"rd47" csv:"rd47"`
	Rd48                float64 `json:"rd48" csv:"rd48"`
	Rd49                float64 `json:"rd49" csv:"rd49"`
	Rd50                float64 `json:"rd50" csv:"rd50"`
	Rd51                float64 `json:"rd51" csv:"rd51"`
	Rd52                float64 `json:"rd52" csv:"rd52"`
	Rd53                float64 `json:"rd53" csv:"rd53"`
	Rd54                float64 `json:"rd54" csv:"rd54"`
	Rd55                float64 `json:"rd55" csv:"rd55"`
	Rd56                float64 `json:"rd56" csv:"rd56"`
	Rd57                float64 `json:"rd57" csv:"rd57"`
	Rd58                float64 `json:"rd58" csv:"rd58"`
	Rd59                float64 `json:"rd59" csv:"rd59"`
	Rd60                float64 `json:"rd60" csv:"rd60"`
	Rd61                float64 `json:"rd61" csv:"rd61"`
	Rd62                float64 `json:"rd62" csv:"rd62"`
	Rd63                float64 `json:"rd63" csv:"rd63"`
	Rd64                float64 `json:"rd64" csv:"rd64"`
	Rd65                float64 `json:"rd65" csv:"rd65"`
	Rd66                float64 `json:"rd66" csv:"rd66"`
	Rd67                float64 `json:"rd67" csv:"rd67"`
	Rd68                float64 `json:"rd68" csv:"rd68"`
	Rd69                float64 `json:"rd69" csv:"rd69"`
	Rd70                float64 `json:"rd70" csv:"rd70"`
	Rd71                float64 `json:"rd71" csv:"rd71"`
	Rd72                float64 `json:"rd72" csv:"rd72"`
	Rd73                float64 `json:"rd73" csv:"rd73"`
	Rd74                float64 `json:"rd74" csv:"rd74"`
	Rd75                float64 `json:"rd75" csv:"rd75"`
	Rd76                float64 `json:"rd76" csv:"rd76"`
	Rd77                float64 `json:"rd77" csv:"rd77"`
	Rd78                float64 `json:"rd78" csv:"rd78"`
	Rd79                float64 `json:"rd79" csv:"rd79"`
	Rd80                float64 `json:"rd80" csv:"rd80"`
	Rd81                float64 `json:"rd81" csv:"rd81"`
	Rd82                float64 `json:"rd82" csv:"rd82"`
	Rd83                float64 `json:"rd83" csv:"rd83"`
	Rd84                float64 `json:"rd84" csv:"rd84"`
	Rd85                float64 `json:"rd85" csv:"rd85"`
	Rd86                float64 `json:"rd86" csv:"rd86"`
	Rd87                float64 `json:"rd87" csv:"rd87"`
	Rd88                float64 `json:"rd88" csv:"rd88"`
	Rd89                float64 `json:"rd89" csv:"rd89"`
	Rd90                float64 `json:"rd90" csv:"rd90"`
	Rd91                float64 `json:"rd91" csv:"rd91"`
	Rd92                float64 `json:"rd92" csv:"rd92"`
	Rd93                float64 `json:"rd93" csv:"rd93"`
	Rd94                float64 `json:"rd94" csv:"rd94"`
	Rd95                float64 `json:"rd95" csv:"rd95"`
	Rd96                float64 `json:"rd96" csv:"rd96"`
	Rd97                float64 `json:"rd97" csv:"rd97"`
	Rd98                float64 `json:"rd98" csv:"rd98"`
	Rd99                float64 `json:"rd99" csv:"rd99"`
	Rd100               float64 `json:"rd100" csv:"rd100"`
	Rd101               float64 `json:"rd101" csv:"rd101"`
	Rd102               float64 `json:"rd102" csv:"rd102"`
	Rd103               float64 `json:"rd103" csv:"rd103"`
	Rd104               float64 `json:"rd104" csv:"rd104"`
	Rd105               float64 `json:"rd105" csv:"rd105"`
	Rd106               float64 `json:"rd106" csv:"rd106"`
	Rd107               float64 `json:"rd107" csv:"rd107"`
	Rd108               float64 `json:"rd108" csv:"rd108"`
	Rd109               float64 `json:"rd109" csv:"rd109"`
	Rd110               float64 `json:"rd110" csv:"rd110"`
	Rd111               float64 `json:"rd111" csv:"rd111"`
	Rd112               float64 `json:"rd112" csv:"rd112"`
	Rd113               float64 `json:"rd113" csv:"rd113"`
	Rd114               float64 `json:"rd114" csv:"rd114"`
	Rd115               float64 `json:"rd115" csv:"rd115"`
	Rd116               float64 `json:"rd116" csv:"rd116"`
	Rd117               float64 `json:"rd117" csv:"rd117"`
	Rd118               float64 `json:"rd118" csv:"rd118"`
	Rd119               float64 `json:"rd119" csv:"rd119"`
	Rd120               float64 `json:"rd120" csv:"rd120"`
}

type LicDevelopmentFactorClaimCount struct {
	ID                  int     `json:"-" gorm:"primary_key"`
	RunDate             string  `json:"-"`
	RunID               int     `json:"-"`
	PortfolioName       string  `json:"-"`
	ProductCode         string  `json:"product_code,omitempty"`
	DevelopmentVariable string  `json:"development_variable"`
	Rd0                 float64 `json:"rd0"`
	Rd1                 float64 `json:"rd1"`
	Rd2                 float64 `json:"rd2"`
	Rd3                 float64 `json:"rd3"`
	Rd4                 float64 `json:"rd4"`
	Rd5                 float64 `json:"rd5"`
	Rd6                 float64 `json:"rd6"`
	Rd7                 float64 `json:"rd7"`
	Rd8                 float64 `json:"rd8"`
	Rd9                 float64 `json:"rd9"`
	Rd10                float64 `json:"rd10"`
	Rd11                float64 `json:"rd11"`
	Rd12                float64 `json:"rd12"`
	Rd13                float64 `json:"rd13"`
	Rd14                float64 `json:"rd14"`
	Rd15                float64 `json:"rd15"`
	Rd16                float64 `json:"rd16"`
	Rd17                float64 `json:"rd17"`
	Rd18                float64 `json:"rd18"`
	Rd19                float64 `json:"rd19"`
	Rd20                float64 `json:"rd20"`
	Rd21                float64 `json:"rd21"`
	Rd22                float64 `json:"rd22"`
	Rd23                float64 `json:"rd23"`
	Rd24                float64 `json:"rd24"`
	Rd25                float64 `json:"rd25"`
	Rd26                float64 `json:"rd26"`
	Rd27                float64 `json:"rd27"`
	Rd28                float64 `json:"rd28"`
	Rd29                float64 `json:"rd29"`
	Rd30                float64 `json:"rd30"`
	Rd31                float64 `json:"rd31"`
	Rd32                float64 `json:"rd32"`
	Rd33                float64 `json:"rd33"`
	Rd34                float64 `json:"rd34"`
	Rd35                float64 `json:"rd35"`
	Rd36                float64 `json:"rd36"`
	Rd37                float64 `json:"rd37"`
	Rd38                float64 `json:"rd38"`
	Rd39                float64 `json:"rd39"`
	Rd40                float64 `json:"rd40"`
	Rd41                float64 `json:"rd41"`
	Rd42                float64 `json:"rd42"`
	Rd43                float64 `json:"rd43"`
	Rd44                float64 `json:"rd44"`
	Rd45                float64 `json:"rd45"`
	Rd46                float64 `json:"rd46"`
	Rd47                float64 `json:"rd47"`
	Rd48                float64 `json:"rd48"`
	Rd49                float64 `json:"rd49"`
	Rd50                float64 `json:"rd50"`
	Rd51                float64 `json:"rd51"`
	Rd52                float64 `json:"rd52"`
	Rd53                float64 `json:"rd53"`
	Rd54                float64 `json:"rd54"`
	Rd55                float64 `json:"rd55"`
	Rd56                float64 `json:"rd56"`
	Rd57                float64 `json:"rd57"`
	Rd58                float64 `json:"rd58"`
	Rd59                float64 `json:"rd59"`
	Rd60                float64 `json:"rd60"`
	Rd61                float64 `json:"rd61"`
	Rd62                float64 `json:"rd62"`
	Rd63                float64 `json:"rd63"`
	Rd64                float64 `json:"rd64"`
	Rd65                float64 `json:"rd65"`
	Rd66                float64 `json:"rd66"`
	Rd67                float64 `json:"rd67"`
	Rd68                float64 `json:"rd68"`
	Rd69                float64 `json:"rd69"`
	Rd70                float64 `json:"rd70"`
	Rd71                float64 `json:"rd71"`
	Rd72                float64 `json:"rd72"`
	Rd73                float64 `json:"rd73"`
	Rd74                float64 `json:"rd74"`
	Rd75                float64 `json:"rd75"`
	Rd76                float64 `json:"rd76"`
	Rd77                float64 `json:"rd77"`
	Rd78                float64 `json:"rd78"`
	Rd79                float64 `json:"rd79"`
	Rd80                float64 `json:"rd80"`
	Rd81                float64 `json:"rd81"`
	Rd82                float64 `json:"rd82"`
	Rd83                float64 `json:"rd83"`
	Rd84                float64 `json:"rd84"`
	Rd85                float64 `json:"rd85"`
	Rd86                float64 `json:"rd86"`
	Rd87                float64 `json:"rd87"`
	Rd88                float64 `json:"rd88"`
	Rd89                float64 `json:"rd89"`
	Rd90                float64 `json:"rd90"`
	Rd91                float64 `json:"rd91"`
	Rd92                float64 `json:"rd92"`
	Rd93                float64 `json:"rd93"`
	Rd94                float64 `json:"rd94"`
	Rd95                float64 `json:"rd95"`
	Rd96                float64 `json:"rd96"`
	Rd97                float64 `json:"rd97"`
	Rd98                float64 `json:"rd98"`
	Rd99                float64 `json:"rd99"`
	Rd100               float64 `json:"rd100"`
	Rd101               float64 `json:"rd101"`
	Rd102               float64 `json:"rd102"`
	Rd103               float64 `json:"rd103"`
	Rd104               float64 `json:"rd104"`
	Rd105               float64 `json:"rd105"`
	Rd106               float64 `json:"rd106"`
	Rd107               float64 `json:"rd107"`
	Rd108               float64 `json:"rd108"`
	Rd109               float64 `json:"rd109"`
	Rd110               float64 `json:"rd110"`
	Rd111               float64 `json:"rd111"`
	Rd112               float64 `json:"rd112"`
	Rd113               float64 `json:"rd113"`
	Rd114               float64 `json:"rd114"`
	Rd115               float64 `json:"rd115"`
	Rd116               float64 `json:"rd116"`
	Rd117               float64 `json:"rd117"`
	Rd118               float64 `json:"rd118"`
	Rd119               float64 `json:"rd119"`
	Rd120               float64 `json:"rd120"`
}

type LicDevelopmentFactorAverageClaim struct {
	ID                  int     `json:"-" gorm:"primary_key"`
	RunDate             string  `json:"-"`
	RunID               int     `json:"-"`
	PortfolioName       string  `json:"-"`
	ProductCode         string  `json:"product_code,omitempty"`
	DevelopmentVariable string  `json:"development_variable"`
	Rd0                 float64 `json:"rd0"`
	Rd1                 float64 `json:"rd1"`
	Rd2                 float64 `json:"rd2"`
	Rd3                 float64 `json:"rd3"`
	Rd4                 float64 `json:"rd4"`
	Rd5                 float64 `json:"rd5"`
	Rd6                 float64 `json:"rd6"`
	Rd7                 float64 `json:"rd7"`
	Rd8                 float64 `json:"rd8"`
	Rd9                 float64 `json:"rd9"`
	Rd10                float64 `json:"rd10"`
	Rd11                float64 `json:"rd11"`
	Rd12                float64 `json:"rd12"`
	Rd13                float64 `json:"rd13"`
	Rd14                float64 `json:"rd14"`
	Rd15                float64 `json:"rd15"`
	Rd16                float64 `json:"rd16"`
	Rd17                float64 `json:"rd17"`
	Rd18                float64 `json:"rd18"`
	Rd19                float64 `json:"rd19"`
	Rd20                float64 `json:"rd20"`
	Rd21                float64 `json:"rd21"`
	Rd22                float64 `json:"rd22"`
	Rd23                float64 `json:"rd23"`
	Rd24                float64 `json:"rd24"`
	Rd25                float64 `json:"rd25"`
	Rd26                float64 `json:"rd26"`
	Rd27                float64 `json:"rd27"`
	Rd28                float64 `json:"rd28"`
	Rd29                float64 `json:"rd29"`
	Rd30                float64 `json:"rd30"`
	Rd31                float64 `json:"rd31"`
	Rd32                float64 `json:"rd32"`
	Rd33                float64 `json:"rd33"`
	Rd34                float64 `json:"rd34"`
	Rd35                float64 `json:"rd35"`
	Rd36                float64 `json:"rd36"`
	Rd37                float64 `json:"rd37"`
	Rd38                float64 `json:"rd38"`
	Rd39                float64 `json:"rd39"`
	Rd40                float64 `json:"rd40"`
	Rd41                float64 `json:"rd41"`
	Rd42                float64 `json:"rd42"`
	Rd43                float64 `json:"rd43"`
	Rd44                float64 `json:"rd44"`
	Rd45                float64 `json:"rd45"`
	Rd46                float64 `json:"rd46"`
	Rd47                float64 `json:"rd47"`
	Rd48                float64 `json:"rd48"`
	Rd49                float64 `json:"rd49"`
	Rd50                float64 `json:"rd50"`
	Rd51                float64 `json:"rd51"`
	Rd52                float64 `json:"rd52"`
	Rd53                float64 `json:"rd53"`
	Rd54                float64 `json:"rd54"`
	Rd55                float64 `json:"rd55"`
	Rd56                float64 `json:"rd56"`
	Rd57                float64 `json:"rd57"`
	Rd58                float64 `json:"rd58"`
	Rd59                float64 `json:"rd59"`
	Rd60                float64 `json:"rd60"`
	Rd61                float64 `json:"rd61"`
	Rd62                float64 `json:"rd62"`
	Rd63                float64 `json:"rd63"`
	Rd64                float64 `json:"rd64"`
	Rd65                float64 `json:"rd65"`
	Rd66                float64 `json:"rd66"`
	Rd67                float64 `json:"rd67"`
	Rd68                float64 `json:"rd68"`
	Rd69                float64 `json:"rd69"`
	Rd70                float64 `json:"rd70"`
	Rd71                float64 `json:"rd71"`
	Rd72                float64 `json:"rd72"`
	Rd73                float64 `json:"rd73"`
	Rd74                float64 `json:"rd74"`
	Rd75                float64 `json:"rd75"`
	Rd76                float64 `json:"rd76"`
	Rd77                float64 `json:"rd77"`
	Rd78                float64 `json:"rd78"`
	Rd79                float64 `json:"rd79"`
	Rd80                float64 `json:"rd80"`
	Rd81                float64 `json:"rd81"`
	Rd82                float64 `json:"rd82"`
	Rd83                float64 `json:"rd83"`
	Rd84                float64 `json:"rd84"`
	Rd85                float64 `json:"rd85"`
	Rd86                float64 `json:"rd86"`
	Rd87                float64 `json:"rd87"`
	Rd88                float64 `json:"rd88"`
	Rd89                float64 `json:"rd89"`
	Rd90                float64 `json:"rd90"`
	Rd91                float64 `json:"rd91"`
	Rd92                float64 `json:"rd92"`
	Rd93                float64 `json:"rd93"`
	Rd94                float64 `json:"rd94"`
	Rd95                float64 `json:"rd95"`
	Rd96                float64 `json:"rd96"`
	Rd97                float64 `json:"rd97"`
	Rd98                float64 `json:"rd98"`
	Rd99                float64 `json:"rd99"`
	Rd100               float64 `json:"rd100"`
	Rd101               float64 `json:"rd101"`
	Rd102               float64 `json:"rd102"`
	Rd103               float64 `json:"rd103"`
	Rd104               float64 `json:"rd104"`
	Rd105               float64 `json:"rd105"`
	Rd106               float64 `json:"rd106"`
	Rd107               float64 `json:"rd107"`
	Rd108               float64 `json:"rd108"`
	Rd109               float64 `json:"rd109"`
	Rd110               float64 `json:"rd110"`
	Rd111               float64 `json:"rd111"`
	Rd112               float64 `json:"rd112"`
	Rd113               float64 `json:"rd113"`
	Rd114               float64 `json:"rd114"`
	Rd115               float64 `json:"rd115"`
	Rd116               float64 `json:"rd116"`
	Rd117               float64 `json:"rd117"`
	Rd118               float64 `json:"rd118"`
	Rd119               float64 `json:"rd119"`
	Rd120               float64 `json:"rd120"`
}

type LicCumulativeProjection struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicCumulativeProjectionClaimCount struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicCumulativeProjectionAverageClaim struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicCumulativeProjectionAveragetoTotalClaim struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicIncrementalExpectedSimulation struct {
	ID             int     `json:"id" gorm:"primary_key"`
	RunDate        string  `json:"run_date,omitempty"`
	RunID          int     `json:"run_id"`
	PortfolioName  string  `json:"portfolio_name,omitempty"`
	LicPortfolioId int     `json:"lic_portfolio_id,omitempty"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicExpectedSimulation struct {
	ID             int     `json:"id" gorm:"primary_key"`
	RunDate        string  `json:"run_date,omitempty"`
	RunID          int     `json:"run_id"`
	PortfolioName  string  `json:"portfolio_name,omitempty"`
	LicPortfolioId int     `json:"lic_portfolio_id,omitempty"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicStandardisedResiduals struct {
	ID             int     `json:"id" gorm:"primary_key"`
	RunDate        string  `json:"run_date,omitempty"`
	RunID          int     `json:"run_id"`
	PortfolioName  string  `json:"portfolio_name,omitempty"`
	LicPortfolioId int     `json:"lic_portfolio_id,omitempty"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicIncrementalProjectionAveragetoTotalClaim struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicIncrementalInflatedAveragetoTotalClaim struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicDiscountedIncrementalInflatedAveragetoTotalClaim struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicIbnrReserve struct {
	ID                                     int     `json:"-" gorm:"primary_key"`
	RunDate                                string  `json:"-"`
	RunID                                  int     `json:"-"`
	PortfolioName                          string  `json:"-"`
	LicPortfolioId                         int     `json:"-"`
	ProductCode                            string  `json:"product_code,omitempty"`
	AccidentYear                           int     `json:"accident_year"`
	AccidentMonth                          int     `json:"accident_month"`
	ChainLadderIbnr                        float64 `json:"chain_ladder_ibnr"`
	ChainLadderAverageCostPerClaimIbnr     float64 `json:"chain_ladder_average_cost_per_claim_ibnr"`
	EarnedPremium                          float64 `json:"earned_premium"`
	ExpectedTotalLoss                      float64 `json:"expected_total_loss"`
	ProportionNotRunoff                    float64 `json:"proportion_not_runoff"`
	BfIbnr                                 float64 `json:"bf_ibnr"`
	ProportionRunoff                       float64 `json:"proportion_runoff"`
	PremiumRunoff                          float64 `json:"premium_runoff"`
	ActualClaims                           float64 `json:"actual_claims"`
	PredictedLossRatio                     float64 `json:"predicted_loss_ratio"`
	PredictedTotalLoss                     float64 `json:"predicted_total_loss"`
	CapeCodIbnr                            float64 `json:"cape_cod_ibnr"`
	MackModelIbnr                          float64 `json:"mack_model_ibnr"`
	RiskAdjustment                         float64 `json:"risk_adjustment"`
	ChainLadderIbnrAt12                    float64 `json:"chain_ladder_ibnr_at12"`
	BfIbnrAt12                             float64 `json:"bf_ibnr_at12"`
	CapeCodIbnrAt12                        float64 `json:"cape_cod_ibnr_at12"`
	ChainLadderAverageCostPerClaimIbnrAt12 float64 `json:"chain_ladder_average_cost_per_claim_ibnr_at12"`
	CombinedClBfIbnr                       float64 `json:"combined_cl_bf_ibnr"`
	CombinedClBfIbnrAt12                   float64 `json:"combined_cl_bf_ibnr_at12"`
	CombinedClCapeCodIbnr                  float64 `json:"combined_cl_cape_cod_ibnr"`
	CombinedClCapeCodIbnrAt12              float64 `json:"combined_cl_cape_cod_ibnr_at12"`
	// EffectiveMethod records which IBNR method was actually applied to this accident year.
	// Populated at run time; enables per-AY transparency in the Credibility Breakdown view.
	EffectiveMethod                        string  `json:"effective_method"`
}

type LicDevelopmentFactorGraphData struct {
	DevelopmentVariable string  `json:"development_variable"`
	Rd0                 float64 `json:"rd0"`
	Rd1                 float64 `json:"rd1"`
	Rd2                 float64 `json:"rd2"`
	Rd3                 float64 `json:"rd3"`
	Rd4                 float64 `json:"rd4"`
	Rd5                 float64 `json:"rd5"`
	Rd6                 float64 `json:"rd6"`
	Rd7                 float64 `json:"rd7"`
	Rd8                 float64 `json:"rd8"`
	Rd9                 float64 `json:"rd9"`
	Rd10                float64 `json:"rd10"`
	Rd11                float64 `json:"rd11"`
	Rd12                float64 `json:"rd12"`
	//Rd13                float64 `json:"rd13"`
	//Rd14                float64 `json:"rd14"`
	//Rd15                float64 `json:"rd15"`
	//Rd16                float64 `json:"rd16"`
	//Rd17                float64 `json:"rd17"`
	//Rd18                float64 `json:"rd18"`
	//Rd19                float64 `json:"rd19"`
	//Rd20                float64 `json:"rd20"`
	//Rd21                float64 `json:"rd21"`
	//Rd22                float64 `json:"rd22"`
	//Rd23                float64 `json:"rd23"`
	//Rd24                float64 `json:"rd24"`
	//Rd25                float64 `json:"rd25"`
	//Rd26                float64 `json:"rd26"`
	//Rd27                float64 `json:"rd27"`
	//Rd28                float64 `json:"rd28"`
	//Rd29                float64 `json:"rd29"`
	//Rd30                float64 `json:"rd30"`
	//Rd31                float64 `json:"rd31"`
	//Rd32                float64 `json:"rd32"`
	//Rd33                float64 `json:"rd33"`
	//Rd34                float64 `json:"rd34"`
	//Rd35                float64 `json:"rd35"`
	//Rd36                float64 `json:"rd36"`
	//Rd37                float64 `json:"rd37"`
	//Rd38                float64 `json:"rd38"`
	//Rd39                float64 `json:"rd39"`
	//Rd40                float64 `json:"rd40"`
	//Rd41                float64 `json:"rd41"`
	//Rd42                float64 `json:"rd42"`
	//Rd43                float64 `json:"rd43"`
	//Rd44                float64 `json:"rd44"`
	//Rd45                float64 `json:"rd45"`
	//Rd46                float64 `json:"rd46"`
	//Rd47                float64 `json:"rd47"`
	//Rd48                float64 `json:"rd48"`
	//Rd49                float64 `json:"rd49"`
	//Rd50                float64 `json:"rd50"`
	//Rd51                float64 `json:"rd51"`
	//Rd52                float64 `json:"rd52"`
	//Rd53                float64 `json:"rd53"`
	//Rd54                float64 `json:"rd54"`
	//Rd55                float64 `json:"rd55"`
	//Rd56                float64 `json:"rd56"`
	//Rd57                float64 `json:"rd57"`
	//Rd58                float64 `json:"rd58"`
	//Rd59                float64 `json:"rd59"`
	//Rd60                float64 `json:"rd60"`
	//Rd61                float64 `json:"rd61"`
	//Rd62                float64 `json:"rd62"`
	//Rd63                float64 `json:"rd63"`
	//Rd64                float64 `json:"rd64"`
	//Rd65                float64 `json:"rd65"`
	//Rd66                float64 `json:"rd66"`
	//Rd67                float64 `json:"rd67"`
	//Rd68                float64 `json:"rd68"`
	//Rd69                float64 `json:"rd69"`
	//Rd70                float64 `json:"rd70"`
	//Rd71                float64 `json:"rd71"`
	//Rd72                float64 `json:"rd72"`
	//Rd73                float64 `json:"rd73"`
	//Rd74                float64 `json:"rd74"`
	//Rd75                float64 `json:"rd75"`
	//Rd76                float64 `json:"rd76"`
	//Rd77                float64 `json:"rd77"`
	//Rd78                float64 `json:"rd78"`
	//Rd79                float64 `json:"rd79"`
	//Rd80                float64 `json:"rd80"`
	//Rd81                float64 `json:"rd81"`
	//Rd82                float64 `json:"rd82"`
	//Rd83                float64 `json:"rd83"`
	//Rd84                float64 `json:"rd84"`
	//Rd85                float64 `json:"rd85"`
	//Rd86                float64 `json:"rd86"`
	//Rd87                float64 `json:"rd87"`
	//Rd88                float64 `json:"rd88"`
	//Rd89                float64 `json:"rd89"`
	//Rd90                float64 `json:"rd90"`
	//Rd91                float64 `json:"rd91"`
	//Rd92                float64 `json:"rd92"`
	//Rd93                float64 `json:"rd93"`
	//Rd94                float64 `json:"rd94"`
	//Rd95                float64 `json:"rd95"`
	//Rd96                float64 `json:"rd96"`
	//Rd97                float64 `json:"rd97"`
	//Rd98                float64 `json:"rd98"`
	//Rd99                float64 `json:"rd99"`
	//Rd100               float64 `json:"rd100"`
	//Rd101               float64 `json:"rd101"`
	//Rd102               float64 `json:"rd102"`
	//Rd103               float64 `json:"rd103"`
	//Rd104               float64 `json:"rd104"`
	//Rd105               float64 `json:"rd105"`
	//Rd106               float64 `json:"rd106"`
	//Rd107               float64 `json:"rd107"`
	//Rd108               float64 `json:"rd108"`
	//Rd109               float64 `json:"rd109"`
	//Rd110               float64 `json:"rd110"`
	//Rd111               float64 `json:"rd111"`
	//Rd112               float64 `json:"rd112"`
	//Rd113               float64 `json:"rd113"`
	//Rd114               float64 `json:"rd114"`
	//Rd115               float64 `json:"rd115"`
	//Rd116               float64 `json:"rd116"`
	//Rd117               float64 `json:"rd117"`
	//Rd118               float64 `json:"rd118"`
	//Rd119               float64 `json:"rd119"`
	//Rd120               float64 `json:"rd120"`
}

type IbnrReserveReport struct {
	ID                                 int     `json:"-" gorm:"primary_key"`
	RunDate                            string  `json:"-"`
	RunID                              int     `json:"-"`
	PortfolioName                      string  `json:"-"`
	LicPortfolioId                     int     `json:"-"`
	ProductCode                        string  `json:"product_code,omitempty"`
	IbnrBel                            float64 `json:"ibnr_bel"`
	IBNRRiskAdjustment                 float64 `json:"ibnr_risk_adjustment"`
	IbnrReserve                        float64 `json:"ibnr_reserve"`
	ChainLadderIbnr                    float64 `json:"chain_ladder_ibnr"`
	BootstrapIbnr                      float64 `json:"bootstrap_ibnr"`
	BootstrapNthPercentileIbnr         float64 `json:"bootstrap_nth_percentile_ibnr"`
	MackModelIbnr                      float64 `json:"mack_model_ibnr"`
	MackModelNthPercentileIbnr         float64 `json:"mack_model_nth_percentile_ibnr"`
	ChainLadderAverageCostPerClaimIbnr float64 `json:"chain_ladder_average_cost_per_claim_ibnr"`
	BfIbnr                             float64 `json:"bf_ibnr"`
	CapeCodIbnr                        float64 `json:"cape_cod_ibnr"`
	IbnrBelAt12                        float64 `json:"ibnr_bel_at12"`
	IbnrRiskAdjustmentAt12             float64 `json:"ibnr_risk_adjustment_at12"`
	CombinedClBfIbnr                   float64 `json:"combined_cl_bf_ibnr"`
	CombinedClCapeCodIbnr              float64 `json:"combined_cl_cape_cod_ibnr"`
}

type LicRunSetting struct {
	ID                   int       `json:"id" gorm:"primary_key"`
	RunName              string    `json:"run_name" csv:"run_name"`
	RunDate              string    `json:"run_date" csv:"run_date"`
	OpeningBalanceDate   string    `json:"opening_balance_date" csv:"opening_balance_date"`
	UserName             string    `json:"user_name" csv:"user_name"`
	UserEmail            string    `json:"user_email" csv:"user_email"`
	LicConfigurationName string    `json:"lic_configuration_name" csv:"lic_configuration_name"`
	LicConfigurationId   int       `json:"lic_configuration_id" csv:"lic_configuration_id"`
	LicParameterYear     int       `json:"lic_parameter_year" csv:"lic_parameter_year"`
	LicParameterVersion  string    `json:"lic_parameter_version" csv:"lic_parameter_version"`
	Description          string    `json:"description" csv:"description"`
	TotalRecords         int       `json:"total_records" csv:"total_records"`
	ProcessedRecords     int       `json:"processed_records" csv:"processed_records"`
	ProcessingStatus     string    `json:"processing_status" csv:"processing_status"`
	RunFailed            bool      `json:"run_failed" csv:"run_failed"`
	RunFailureReason     string    `json:"run_failure_reason" csv:"run_failure_reason"`
	RunTime              float64   `json:"run_time" csv:"run_time"`
	CreationDate         time.Time `json:"creation_date"`
}

type LicGeneratedRandomResiduals struct {
	ID             int     `json:"id" gorm:"primary_key"`
	RunDate        string  `json:"run_date,omitempty"`
	RunId          int     `json:"run_id"`
	PortfolioName  string  `json:"portfolio_name,omitempty"`
	LicPortfolioId int     `json:"lic_portfolio_id,omitempty"`
	ProductCode    string  `json:"product_code,omitempty"`
	ResidualNumber int     `json:"residual_number"`
	Residual       float64 `json:"residual"`
	RandomResidual float64 `json:"random_residual"`
}

type LicRandomResiduals struct {
	ID             int     `json:"id" gorm:"primary_key"`
	RunDate        string  `json:"run_date,omitempty"`
	RunID          int     `json:"run_id"`
	PortfolioName  string  `json:"portfolio_name,omitempty"`
	LicPortfolioId int     `json:"lic_portfolio_id,omitempty"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicBootStrappedIncremental struct {
	ID             int     `json:"id" gorm:"primary_key"`
	RunDate        string  `json:"run_date,omitempty"`
	RunID          int     `json:"run_id"`
	PortfolioName  string  `json:"portfolio_name,omitempty"`
	LicPortfolioId int     `json:"lic_portfolio_id,omitempty"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicBootStrappedCumulative struct {
	ID             int     `json:"id" gorm:"primary_key"`
	RunDate        string  `json:"run_date,omitempty"`
	RunID          int     `json:"run_id"`
	PortfolioName  string  `json:"portfolio_name,omitempty"`
	LicPortfolioId int     `json:"lic_portfolio_id,omitempty"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicBootstrappedDevelopmentFactor struct {
	ID                  int     `json:"id" gorm:"primary_key"`
	RunDate             string  `json:"run_date,omitempty"`
	RunID               int     `json:"run_id"`
	PortfolioName       string  `json:"portfolio_name,omitempty"`
	ProductCode         string  `json:"product_code,omitempty"`
	DevelopmentVariable string  `json:"development_variable"`
	Rd0                 float64 `json:"rd0"`
	Rd1                 float64 `json:"rd1"`
	Rd2                 float64 `json:"rd2"`
	Rd3                 float64 `json:"rd3"`
	Rd4                 float64 `json:"rd4"`
	Rd5                 float64 `json:"rd5"`
	Rd6                 float64 `json:"rd6"`
	Rd7                 float64 `json:"rd7"`
	Rd8                 float64 `json:"rd8"`
	Rd9                 float64 `json:"rd9"`
	Rd10                float64 `json:"rd10"`
	Rd11                float64 `json:"rd11"`
	Rd12                float64 `json:"rd12"`
	Rd13                float64 `json:"rd13"`
	Rd14                float64 `json:"rd14"`
	Rd15                float64 `json:"rd15"`
	Rd16                float64 `json:"rd16"`
	Rd17                float64 `json:"rd17"`
	Rd18                float64 `json:"rd18"`
	Rd19                float64 `json:"rd19"`
	Rd20                float64 `json:"rd20"`
	Rd21                float64 `json:"rd21"`
	Rd22                float64 `json:"rd22"`
	Rd23                float64 `json:"rd23"`
	Rd24                float64 `json:"rd24"`
	Rd25                float64 `json:"rd25"`
	Rd26                float64 `json:"rd26"`
	Rd27                float64 `json:"rd27"`
	Rd28                float64 `json:"rd28"`
	Rd29                float64 `json:"rd29"`
	Rd30                float64 `json:"rd30"`
	Rd31                float64 `json:"rd31"`
	Rd32                float64 `json:"rd32"`
	Rd33                float64 `json:"rd33"`
	Rd34                float64 `json:"rd34"`
	Rd35                float64 `json:"rd35"`
	Rd36                float64 `json:"rd36"`
	Rd37                float64 `json:"rd37"`
	Rd38                float64 `json:"rd38"`
	Rd39                float64 `json:"rd39"`
	Rd40                float64 `json:"rd40"`
	Rd41                float64 `json:"rd41"`
	Rd42                float64 `json:"rd42"`
	Rd43                float64 `json:"rd43"`
	Rd44                float64 `json:"rd44"`
	Rd45                float64 `json:"rd45"`
	Rd46                float64 `json:"rd46"`
	Rd47                float64 `json:"rd47"`
	Rd48                float64 `json:"rd48"`
	Rd49                float64 `json:"rd49"`
	Rd50                float64 `json:"rd50"`
	Rd51                float64 `json:"rd51"`
	Rd52                float64 `json:"rd52"`
	Rd53                float64 `json:"rd53"`
	Rd54                float64 `json:"rd54"`
	Rd55                float64 `json:"rd55"`
	Rd56                float64 `json:"rd56"`
	Rd57                float64 `json:"rd57"`
	Rd58                float64 `json:"rd58"`
	Rd59                float64 `json:"rd59"`
	Rd60                float64 `json:"rd60"`
	Rd61                float64 `json:"rd61"`
	Rd62                float64 `json:"rd62"`
	Rd63                float64 `json:"rd63"`
	Rd64                float64 `json:"rd64"`
	Rd65                float64 `json:"rd65"`
	Rd66                float64 `json:"rd66"`
	Rd67                float64 `json:"rd67"`
	Rd68                float64 `json:"rd68"`
	Rd69                float64 `json:"rd69"`
	Rd70                float64 `json:"rd70"`
	Rd71                float64 `json:"rd71"`
	Rd72                float64 `json:"rd72"`
	Rd73                float64 `json:"rd73"`
	Rd74                float64 `json:"rd74"`
	Rd75                float64 `json:"rd75"`
	Rd76                float64 `json:"rd76"`
	Rd77                float64 `json:"rd77"`
	Rd78                float64 `json:"rd78"`
	Rd79                float64 `json:"rd79"`
	Rd80                float64 `json:"rd80"`
	Rd81                float64 `json:"rd81"`
	Rd82                float64 `json:"rd82"`
	Rd83                float64 `json:"rd83"`
	Rd84                float64 `json:"rd84"`
	Rd85                float64 `json:"rd85"`
	Rd86                float64 `json:"rd86"`
	Rd87                float64 `json:"rd87"`
	Rd88                float64 `json:"rd88"`
	Rd89                float64 `json:"rd89"`
	Rd90                float64 `json:"rd90"`
	Rd91                float64 `json:"rd91"`
	Rd92                float64 `json:"rd92"`
	Rd93                float64 `json:"rd93"`
	Rd94                float64 `json:"rd94"`
	Rd95                float64 `json:"rd95"`
	Rd96                float64 `json:"rd96"`
	Rd97                float64 `json:"rd97"`
	Rd98                float64 `json:"rd98"`
	Rd99                float64 `json:"rd99"`
	Rd100               float64 `json:"rd100"`
	Rd101               float64 `json:"rd101"`
	Rd102               float64 `json:"rd102"`
	Rd103               float64 `json:"rd103"`
	Rd104               float64 `json:"rd104"`
	Rd105               float64 `json:"rd105"`
	Rd106               float64 `json:"rd106"`
	Rd107               float64 `json:"rd107"`
	Rd108               float64 `json:"rd108"`
	Rd109               float64 `json:"rd109"`
	Rd110               float64 `json:"rd110"`
	Rd111               float64 `json:"rd111"`
	Rd112               float64 `json:"rd112"`
	Rd113               float64 `json:"rd113"`
	Rd114               float64 `json:"rd114"`
	Rd115               float64 `json:"rd115"`
	Rd116               float64 `json:"rd116"`
	Rd117               float64 `json:"rd117"`
	Rd118               float64 `json:"rd118"`
	Rd119               float64 `json:"rd119"`
	Rd120               float64 `json:"rd120"`
}

type LicBootStrappedCumulativeProjection struct {
	ID             int     `json:"id" gorm:"primary_key"`
	RunDate        string  `json:"run_date,omitempty"`
	RunID          int     `json:"run_id"`
	PortfolioName  string  `json:"portfolio_name,omitempty"`
	LicPortfolioId int     `json:"lic_portfolio_id,omitempty"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicBootStrappedIncrementalProjection struct {
	ID             int     `json:"id" gorm:"primary_key"`
	RunDate        string  `json:"run_date,omitempty"`
	RunID          int     `json:"run_id"`
	PortfolioName  string  `json:"portfolio_name,omitempty"`
	LicPortfolioId int     `json:"lic_portfolio_id,omitempty"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicBootStrappedIncrementalInflatedProjection struct {
	ID             int     `json:"id" gorm:"primary_key"`
	RunDate        string  `json:"run_date,omitempty"`
	RunID          int     `json:"run_id"`
	PortfolioName  string  `json:"portfolio_name,omitempty"`
	LicPortfolioId int     `json:"lic_portfolio_id,omitempty"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicBootStrappedIncrementalInflatedDiscountedProjection struct {
	ID             int     `json:"id" gorm:"primary_key"`
	RunDate        string  `json:"run_date,omitempty"`
	RunID          int     `json:"run_id"`
	PortfolioName  string  `json:"portfolio_name,omitempty"`
	LicPortfolioId int     `json:"lic_portfolio_id,omitempty"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicBootstrappedResults struct {
	ID               int     `json:"id" gorm:"primary_key"`
	RunDate          string  `json:"run_date,omitempty"`
	RunID            int     `json:"run_id"`
	PortfolioName    string  `json:"portfolio_name,omitempty"`
	LicPortfolioId   int     `json:"lic_portfolio_id,omitempty"`
	ProductCode      string  `json:"product_code,omitempty"`
	SimulationNumber int     `json:"simulation_number"`
	Reserve          float64 `json:"reserve"`
}

type LicBootstrappedResultSummary struct {
	ID                int     `json:"-" gorm:"primary_key"`
	RunDate           string  `json:"-"`
	RunID             int     `json:"-"`
	PortfolioName     string  `json:"-"`
	LicPortfolioId    int     `json:"-"`
	ProductCode       string  `json:"product_code,omitempty"`
	Mean              float64 `json:"mean"`
	Median            float64 `json:"median"`
	Percentile        float64 `json:"percentile"`
	StandardDeviation float64 `json:"standard_deviation"`
	Minimum           float64 `json:"minimum"`
	Maximum           float64 `json:"maximum"`
}

type IbnrFrequency struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	Group          int     `json:"group"`
	Reserve        float64 `json:"reserve"`
	Frequency      int     `json:"frequency"`
}

type LicIncrementalProjection struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicIncrementalInflated struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicDiscountedIncrementalInflated struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicIndividualDevelopmentFactors struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicBiasAdjustmentFactor struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicMackModelCalculatedParameters struct {
	ID                int     `json:"-" gorm:"primary_key"`
	RunDate           string  `json:"-"`
	RunID             int     `json:"-"`
	PortfolioName     string  `json:"-"`
	ProductCode       string  `json:"product_code,omitempty"`
	MackModelVariable string  `json:"mack_model_variable"`
	Rd0               float64 `json:"rd0"`
	Rd1               float64 `json:"rd1"`
	Rd2               float64 `json:"rd2"`
	Rd3               float64 `json:"rd3"`
	Rd4               float64 `json:"rd4"`
	Rd5               float64 `json:"rd5"`
	Rd6               float64 `json:"rd6"`
	Rd7               float64 `json:"rd7"`
	Rd8               float64 `json:"rd8"`
	Rd9               float64 `json:"rd9"`
	Rd10              float64 `json:"rd10"`
	Rd11              float64 `json:"rd11"`
	Rd12              float64 `json:"rd12"`
	Rd13              float64 `json:"rd13"`
	Rd14              float64 `json:"rd14"`
	Rd15              float64 `json:"rd15"`
	Rd16              float64 `json:"rd16"`
	Rd17              float64 `json:"rd17"`
	Rd18              float64 `json:"rd18"`
	Rd19              float64 `json:"rd19"`
	Rd20              float64 `json:"rd20"`
	Rd21              float64 `json:"rd21"`
	Rd22              float64 `json:"rd22"`
	Rd23              float64 `json:"rd23"`
	Rd24              float64 `json:"rd24"`
	Rd25              float64 `json:"rd25"`
	Rd26              float64 `json:"rd26"`
	Rd27              float64 `json:"rd27"`
	Rd28              float64 `json:"rd28"`
	Rd29              float64 `json:"rd29"`
	Rd30              float64 `json:"rd30"`
	Rd31              float64 `json:"rd31"`
	Rd32              float64 `json:"rd32"`
	Rd33              float64 `json:"rd33"`
	Rd34              float64 `json:"rd34"`
	Rd35              float64 `json:"rd35"`
	Rd36              float64 `json:"rd36"`
	Rd37              float64 `json:"rd37"`
	Rd38              float64 `json:"rd38"`
	Rd39              float64 `json:"rd39"`
	Rd40              float64 `json:"rd40"`
	Rd41              float64 `json:"rd41"`
	Rd42              float64 `json:"rd42"`
	Rd43              float64 `json:"rd43"`
	Rd44              float64 `json:"rd44"`
	Rd45              float64 `json:"rd45"`
	Rd46              float64 `json:"rd46"`
	Rd47              float64 `json:"rd47"`
	Rd48              float64 `json:"rd48"`
	Rd49              float64 `json:"rd49"`
	Rd50              float64 `json:"rd50"`
	Rd51              float64 `json:"rd51"`
	Rd52              float64 `json:"rd52"`
	Rd53              float64 `json:"rd53"`
	Rd54              float64 `json:"rd54"`
	Rd55              float64 `json:"rd55"`
	Rd56              float64 `json:"rd56"`
	Rd57              float64 `json:"rd57"`
	Rd58              float64 `json:"rd58"`
	Rd59              float64 `json:"rd59"`
	Rd60              float64 `json:"rd60"`
	Rd61              float64 `json:"rd61"`
	Rd62              float64 `json:"rd62"`
	Rd63              float64 `json:"rd63"`
	Rd64              float64 `json:"rd64"`
	Rd65              float64 `json:"rd65"`
	Rd66              float64 `json:"rd66"`
	Rd67              float64 `json:"rd67"`
	Rd68              float64 `json:"rd68"`
	Rd69              float64 `json:"rd69"`
	Rd70              float64 `json:"rd70"`
	Rd71              float64 `json:"rd71"`
	Rd72              float64 `json:"rd72"`
	Rd73              float64 `json:"rd73"`
	Rd74              float64 `json:"rd74"`
	Rd75              float64 `json:"rd75"`
	Rd76              float64 `json:"rd76"`
	Rd77              float64 `json:"rd77"`
	Rd78              float64 `json:"rd78"`
	Rd79              float64 `json:"rd79"`
	Rd80              float64 `json:"rd80"`
	Rd81              float64 `json:"rd81"`
	Rd82              float64 `json:"rd82"`
	Rd83              float64 `json:"rd83"`
	Rd84              float64 `json:"rd84"`
	Rd85              float64 `json:"rd85"`
	Rd86              float64 `json:"rd86"`
	Rd87              float64 `json:"rd87"`
	Rd88              float64 `json:"rd88"`
	Rd89              float64 `json:"rd89"`
	Rd90              float64 `json:"rd90"`
	Rd91              float64 `json:"rd91"`
	Rd92              float64 `json:"rd92"`
	Rd93              float64 `json:"rd93"`
	Rd94              float64 `json:"rd94"`
	Rd95              float64 `json:"rd95"`
	Rd96              float64 `json:"rd96"`
	Rd97              float64 `json:"rd97"`
	Rd98              float64 `json:"rd98"`
	Rd99              float64 `json:"rd99"`
	Rd100             float64 `json:"rd100"`
	Rd101             float64 `json:"rd101"`
	Rd102             float64 `json:"rd102"`
	Rd103             float64 `json:"rd103"`
	Rd104             float64 `json:"rd104"`
	Rd105             float64 `json:"rd105"`
	Rd106             float64 `json:"rd106"`
	Rd107             float64 `json:"rd107"`
	Rd108             float64 `json:"rd108"`
	Rd109             float64 `json:"rd109"`
	Rd110             float64 `json:"rd110"`
	Rd111             float64 `json:"rd111"`
	Rd112             float64 `json:"rd112"`
	Rd113             float64 `json:"rd113"`
	Rd114             float64 `json:"rd114"`
	Rd115             float64 `json:"rd115"`
	Rd116             float64 `json:"rd116"`
	Rd117             float64 `json:"rd117"`
	Rd118             float64 `json:"rd118"`
	Rd119             float64 `json:"rd119"`
	Rd120             float64 `json:"rd120"`
}

type LicBiasAdjustedResiduals struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicMeanBiasAdjustedResiduals struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicPseudoRatios struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicMackModelSimulatedDevelopmentFactor struct {
	ID                int     `json:"-" gorm:"primary_key"`
	RunDate           string  `json:"-"`
	RunID             int     `json:"-"`
	PortfolioName     string  `json:"-"`
	ProductCode       string  `json:"product_code,omitempty"`
	MackModelVariable string  `json:"mack_model_variable"`
	Rd0               float64 `json:"rd0"`
	Rd1               float64 `json:"rd1"`
	Rd2               float64 `json:"rd2"`
	Rd3               float64 `json:"rd3"`
	Rd4               float64 `json:"rd4"`
	Rd5               float64 `json:"rd5"`
	Rd6               float64 `json:"rd6"`
	Rd7               float64 `json:"rd7"`
	Rd8               float64 `json:"rd8"`
	Rd9               float64 `json:"rd9"`
	Rd10              float64 `json:"rd10"`
	Rd11              float64 `json:"rd11"`
	Rd12              float64 `json:"rd12"`
	Rd13              float64 `json:"rd13"`
	Rd14              float64 `json:"rd14"`
	Rd15              float64 `json:"rd15"`
	Rd16              float64 `json:"rd16"`
	Rd17              float64 `json:"rd17"`
	Rd18              float64 `json:"rd18"`
	Rd19              float64 `json:"rd19"`
	Rd20              float64 `json:"rd20"`
	Rd21              float64 `json:"rd21"`
	Rd22              float64 `json:"rd22"`
	Rd23              float64 `json:"rd23"`
	Rd24              float64 `json:"rd24"`
	Rd25              float64 `json:"rd25"`
	Rd26              float64 `json:"rd26"`
	Rd27              float64 `json:"rd27"`
	Rd28              float64 `json:"rd28"`
	Rd29              float64 `json:"rd29"`
	Rd30              float64 `json:"rd30"`
	Rd31              float64 `json:"rd31"`
	Rd32              float64 `json:"rd32"`
	Rd33              float64 `json:"rd33"`
	Rd34              float64 `json:"rd34"`
	Rd35              float64 `json:"rd35"`
	Rd36              float64 `json:"rd36"`
	Rd37              float64 `json:"rd37"`
	Rd38              float64 `json:"rd38"`
	Rd39              float64 `json:"rd39"`
	Rd40              float64 `json:"rd40"`
	Rd41              float64 `json:"rd41"`
	Rd42              float64 `json:"rd42"`
	Rd43              float64 `json:"rd43"`
	Rd44              float64 `json:"rd44"`
	Rd45              float64 `json:"rd45"`
	Rd46              float64 `json:"rd46"`
	Rd47              float64 `json:"rd47"`
	Rd48              float64 `json:"rd48"`
	Rd49              float64 `json:"rd49"`
	Rd50              float64 `json:"rd50"`
	Rd51              float64 `json:"rd51"`
	Rd52              float64 `json:"rd52"`
	Rd53              float64 `json:"rd53"`
	Rd54              float64 `json:"rd54"`
	Rd55              float64 `json:"rd55"`
	Rd56              float64 `json:"rd56"`
	Rd57              float64 `json:"rd57"`
	Rd58              float64 `json:"rd58"`
	Rd59              float64 `json:"rd59"`
	Rd60              float64 `json:"rd60"`
	Rd61              float64 `json:"rd61"`
	Rd62              float64 `json:"rd62"`
	Rd63              float64 `json:"rd63"`
	Rd64              float64 `json:"rd64"`
	Rd65              float64 `json:"rd65"`
	Rd66              float64 `json:"rd66"`
	Rd67              float64 `json:"rd67"`
	Rd68              float64 `json:"rd68"`
	Rd69              float64 `json:"rd69"`
	Rd70              float64 `json:"rd70"`
	Rd71              float64 `json:"rd71"`
	Rd72              float64 `json:"rd72"`
	Rd73              float64 `json:"rd73"`
	Rd74              float64 `json:"rd74"`
	Rd75              float64 `json:"rd75"`
	Rd76              float64 `json:"rd76"`
	Rd77              float64 `json:"rd77"`
	Rd78              float64 `json:"rd78"`
	Rd79              float64 `json:"rd79"`
	Rd80              float64 `json:"rd80"`
	Rd81              float64 `json:"rd81"`
	Rd82              float64 `json:"rd82"`
	Rd83              float64 `json:"rd83"`
	Rd84              float64 `json:"rd84"`
	Rd85              float64 `json:"rd85"`
	Rd86              float64 `json:"rd86"`
	Rd87              float64 `json:"rd87"`
	Rd88              float64 `json:"rd88"`
	Rd89              float64 `json:"rd89"`
	Rd90              float64 `json:"rd90"`
	Rd91              float64 `json:"rd91"`
	Rd92              float64 `json:"rd92"`
	Rd93              float64 `json:"rd93"`
	Rd94              float64 `json:"rd94"`
	Rd95              float64 `json:"rd95"`
	Rd96              float64 `json:"rd96"`
	Rd97              float64 `json:"rd97"`
	Rd98              float64 `json:"rd98"`
	Rd99              float64 `json:"rd99"`
	Rd100             float64 `json:"rd100"`
	Rd101             float64 `json:"rd101"`
	Rd102             float64 `json:"rd102"`
	Rd103             float64 `json:"rd103"`
	Rd104             float64 `json:"rd104"`
	Rd105             float64 `json:"rd105"`
	Rd106             float64 `json:"rd106"`
	Rd107             float64 `json:"rd107"`
	Rd108             float64 `json:"rd108"`
	Rd109             float64 `json:"rd109"`
	Rd110             float64 `json:"rd110"`
	Rd111             float64 `json:"rd111"`
	Rd112             float64 `json:"rd112"`
	Rd113             float64 `json:"rd113"`
	Rd114             float64 `json:"rd114"`
	Rd115             float64 `json:"rd115"`
	Rd116             float64 `json:"rd116"`
	Rd117             float64 `json:"rd117"`
	Rd118             float64 `json:"rd118"`
	Rd119             float64 `json:"rd119"`
	Rd120             float64 `json:"rd120"`
}

type LicMackCumulativeProjection struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicMackReserves struct {
	ID             int     `json:"id" gorm:"primary_key"`
	RunDate        string  `json:"run_date,omitempty"`
	RunID          int     `json:"run_id"`
	PortfolioName  string  `json:"portfolio_name,omitempty"`
	LicPortfolioId int     `json:"lic_portfolio_id,omitempty"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicMackSimulationResults struct {
	ID               int     `json:"id" gorm:"primary_key"`
	RunDate          string  `json:"run_date,omitempty"`
	RunID            int     `json:"run_id"`
	PortfolioName    string  `json:"portfolio_name,omitempty"`
	LicPortfolioId   int     `json:"lic_portfolio_id,omitempty"`
	ProductCode      string  `json:"product_code,omitempty"`
	SimulationNumber int     `json:"simulation_number"`
	Reserve          float64 `json:"reserve"`
}

type LicMackSimulationSummaryStats struct {
	ID                int     `json:"-" gorm:"primary_key"`
	RunDate           string  `json:"-"`
	RunID             int     `json:"-"`
	PortfolioName     string  `json:"-"`
	LicPortfolioId    int     `json:"-"`
	ProductCode       string  `json:"product_code,omitempty"`
	Mean              float64 `json:"mean"`
	Median            float64 `json:"median"`
	Percentile        float64 `json:"percentile"`
	StandardDeviation float64 `json:"standard_deviation"`
	Minimum           float64 `json:"minimum"`
	Maximum           float64 `json:"maximum"`
}

type MackIbnrFrequency struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	Group          int     `json:"group"`
	Reserve        float64 `json:"reserve"`
	Frequency      int     `json:"frequency"`
}

type LicLogNormalSigmas struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicLogNormalMeans struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}

type LicLogNormalStandardDeviations struct {
	ID             int     `json:"-" gorm:"primary_key"`
	RunDate        string  `json:"-"`
	RunID          int     `json:"-"`
	PortfolioName  string  `json:"-"`
	LicPortfolioId int     `json:"-"`
	ProductCode    string  `json:"product_code,omitempty"`
	AccidentYear   int     `json:"accident_year"`
	AccidentMonth  int     `json:"accident_month"`
	EarnedPremium  float64 `json:"earned_premium"`
	Rd0            float64 `json:"rd0"`
	Rd1            float64 `json:"rd1"`
	Rd2            float64 `json:"rd2"`
	Rd3            float64 `json:"rd3"`
	Rd4            float64 `json:"rd4"`
	Rd5            float64 `json:"rd5"`
	Rd6            float64 `json:"rd6"`
	Rd7            float64 `json:"rd7"`
	Rd8            float64 `json:"rd8"`
	Rd9            float64 `json:"rd9"`
	Rd10           float64 `json:"rd10"`
	Rd11           float64 `json:"rd11"`
	Rd12           float64 `json:"rd12"`
	Rd13           float64 `json:"rd13"`
	Rd14           float64 `json:"rd14"`
	Rd15           float64 `json:"rd15"`
	Rd16           float64 `json:"rd16"`
	Rd17           float64 `json:"rd17"`
	Rd18           float64 `json:"rd18"`
	Rd19           float64 `json:"rd19"`
	Rd20           float64 `json:"rd20"`
	Rd21           float64 `json:"rd21"`
	Rd22           float64 `json:"rd22"`
	Rd23           float64 `json:"rd23"`
	Rd24           float64 `json:"rd24"`
	Rd25           float64 `json:"rd25"`
	Rd26           float64 `json:"rd26"`
	Rd27           float64 `json:"rd27"`
	Rd28           float64 `json:"rd28"`
	Rd29           float64 `json:"rd29"`
	Rd30           float64 `json:"rd30"`
	Rd31           float64 `json:"rd31"`
	Rd32           float64 `json:"rd32"`
	Rd33           float64 `json:"rd33"`
	Rd34           float64 `json:"rd34"`
	Rd35           float64 `json:"rd35"`
	Rd36           float64 `json:"rd36"`
	Rd37           float64 `json:"rd37"`
	Rd38           float64 `json:"rd38"`
	Rd39           float64 `json:"rd39"`
	Rd40           float64 `json:"rd40"`
	Rd41           float64 `json:"rd41"`
	Rd42           float64 `json:"rd42"`
	Rd43           float64 `json:"rd43"`
	Rd44           float64 `json:"rd44"`
	Rd45           float64 `json:"rd45"`
	Rd46           float64 `json:"rd46"`
	Rd47           float64 `json:"rd47"`
	Rd48           float64 `json:"rd48"`
	Rd49           float64 `json:"rd49"`
	Rd50           float64 `json:"rd50"`
	Rd51           float64 `json:"rd51"`
	Rd52           float64 `json:"rd52"`
	Rd53           float64 `json:"rd53"`
	Rd54           float64 `json:"rd54"`
	Rd55           float64 `json:"rd55"`
	Rd56           float64 `json:"rd56"`
	Rd57           float64 `json:"rd57"`
	Rd58           float64 `json:"rd58"`
	Rd59           float64 `json:"rd59"`
	Rd60           float64 `json:"rd60"`
	Rd61           float64 `json:"rd61"`
	Rd62           float64 `json:"rd62"`
	Rd63           float64 `json:"rd63"`
	Rd64           float64 `json:"rd64"`
	Rd65           float64 `json:"rd65"`
	Rd66           float64 `json:"rd66"`
	Rd67           float64 `json:"rd67"`
	Rd68           float64 `json:"rd68"`
	Rd69           float64 `json:"rd69"`
	Rd70           float64 `json:"rd70"`
	Rd71           float64 `json:"rd71"`
	Rd72           float64 `json:"rd72"`
	Rd73           float64 `json:"rd73"`
	Rd74           float64 `json:"rd74"`
	Rd75           float64 `json:"rd75"`
	Rd76           float64 `json:"rd76"`
	Rd77           float64 `json:"rd77"`
	Rd78           float64 `json:"rd78"`
	Rd79           float64 `json:"rd79"`
	Rd80           float64 `json:"rd80"`
	Rd81           float64 `json:"rd81"`
	Rd82           float64 `json:"rd82"`
	Rd83           float64 `json:"rd83"`
	Rd84           float64 `json:"rd84"`
	Rd85           float64 `json:"rd85"`
	Rd86           float64 `json:"rd86"`
	Rd87           float64 `json:"rd87"`
	Rd88           float64 `json:"rd88"`
	Rd89           float64 `json:"rd89"`
	Rd90           float64 `json:"rd90"`
	Rd91           float64 `json:"rd91"`
	Rd92           float64 `json:"rd92"`
	Rd93           float64 `json:"rd93"`
	Rd94           float64 `json:"rd94"`
	Rd95           float64 `json:"rd95"`
	Rd96           float64 `json:"rd96"`
	Rd97           float64 `json:"rd97"`
	Rd98           float64 `json:"rd98"`
	Rd99           float64 `json:"rd99"`
	Rd100          float64 `json:"rd100"`
	Rd101          float64 `json:"rd101"`
	Rd102          float64 `json:"rd102"`
	Rd103          float64 `json:"rd103"`
	Rd104          float64 `json:"rd104"`
	Rd105          float64 `json:"rd105"`
	Rd106          float64 `json:"rd106"`
	Rd107          float64 `json:"rd107"`
	Rd108          float64 `json:"rd108"`
	Rd109          float64 `json:"rd109"`
	Rd110          float64 `json:"rd110"`
	Rd111          float64 `json:"rd111"`
	Rd112          float64 `json:"rd112"`
	Rd113          float64 `json:"rd113"`
	Rd114          float64 `json:"rd114"`
	Rd115          float64 `json:"rd115"`
	Rd116          float64 `json:"rd116"`
	Rd117          float64 `json:"rd117"`
	Rd118          float64 `json:"rd118"`
	Rd119          float64 `json:"rd119"`
	Rd120          float64 `json:"rd120"`
}
