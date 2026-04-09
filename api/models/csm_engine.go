package models

import (
	"time"
)

type RAConfiguration struct {
	ProductCode        string  `json:"product_code" gorm:"primary_key"`
	Spcode             int     `json:"spcode"`
	LevelOfAggregation string  `json:"level_of_aggregation"`
	RiskEntityShare    float64 `json:"risk_entity_share"`
}

type RiskDriver struct {
	ProductCode   string `json:"product_code" gorm:"primary_key"`
	MortalityRisk string `json:"mortality_risk" csv:"mortality_risk" `
	MorbidityRisk string `json:"morbidity_risk" csv:"morbidity_risk"`
	LongevityRisk string `json:"longevity_risk" csv:"longevity_risk"`
	ExpenseRisk   string `json:"expense_risk" csv:"expense_risk"`
	LapseRisk     string `json:"lapse_risk" csv:"lapse_risk"`
	Catastrophe   string `json:"catastrophe" csv:"catastrophe"`
	Operational   string `json:"operational" csv:"operational"`
}

type RiskAdjustmentFactor struct {
	ID            int     `json:"-" gorm:"primary_key"`
	ProductCode   string  `json:"product_code" csv:"product_code"`
	MortalityRisk float64 `json:"mortality_risk" csv:"mortality_risk"`
	MorbidityRisk float64 `json:"morbidity_risk" csv:"morbidity_risk"`
	LongevityRisk float64 `json:"longevity_risk" csv:"longevity_risk"`
	ExpenseRisk   float64 `json:"expense_risk" csv:"expense_risk"`
	LapseRisk     float64 `json:"lapse_risk" csv:"lapse_risk"`
	Catastrophe   float64 `json:"catastrophe" csv:"catastrophe"`
	Operational   float64 `json:"operational" csv:"operational"`
	Year          int     `json:"year" csv:"year"`
	Version       string  `json:"version" csv:"version"`
	Created       int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
	CreatedBy     string  `json:"created_by" csv:"created_by"`
}

type RiskAdjustmentDriver struct {
	ID       int    `json:"-" gorm:"primary_key"`
	RiskType string `json:"risk_type" csv:"risk_type"`
}

type FinanceVariables struct {
	ID                               int     `json:"-" gorm:"primary_key"`
	ProductCode                      string  `json:"product_code" csv:"product_code"`
	IFRS17Group                      string  `json:"ifrs17_group" csv:"ifrs17_group"`
	IFStatus                         string  `json:"if_status" csv:"if_status"`
	DurationInForceMonths            int     `json:"duration_in_force_months" csv:"duration_in_force_months"`
	LockedInYear                     int     `json:"locked_in_year" csv:"locked_in_year"`
	LockedInMonth                    int     `json:"locked_in_month" csv:"locked_in_month"`
	ActualPremiumIncome              float64 `json:"actual_premium_income" csv:"actual_premium_income"`
	PremiumDebtors                   float64 `json:"premium_debtors"`
	ActualAttributableExpenses       float64 `json:"actual_attributable_expenses" csv:"actual_attributable_expenses"`
	ActualNonAttributableExpenses    float64 `json:"actual_non_attributable_expenses" csv:"actual_non_attributable_expenses"`
	ActualMortalityClaimsIncurred    float64 `json:"actual_mortality_claims_incurred" csv:"actual_mortality_claims_incurred"`
	ActualRetrenchmentClaimsIncurred float64 `json:"actual_retrenchment_claims_incurred" csv:"actual_retrenchment_claims_incurred"`
	ActualMorbidityClaimsIncurred    float64 `json:"actual_morbidity_claims_incurred" csv:"actual_morbidity_claims_incurred"`
	ActualNonLifeClaimsIncurred      float64 `json:"actual_non_life_claims_incurred" csv:"actual_non_life_claims_incurred"`
	ActualAcquisitionExpenses        float64 `json:"actual_acquisition_expenses" csv:"actual_acquisition_expenses"`
	RiskMitigationProportion         float64 `json:"risk_mitigation_proportion" csv:"risk_mitigation_proportion"`
	YieldCurveBasis                  string  `json:"yield_curve_basis" csv:"yield_curve_basis"`
	YieldCurveCode                   string  `json:"yield_curve_code" csv:"yield_curve_code"`
	Year                             int     `json:"year" csv:"year"`
	Version                          string  `json:"version" csv:"version"`
	Created                          int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
	CreatedBy                        string  `json:"created_by" csv:"created_by"`
}

type AOSConfiguration struct {
	RunName        string `json:"run_name" gorm:"primary_key"`
	Description    string `json:"description"`
	ModelPointFile int    `json:"model_point_file"`
	Mortality      int    `json:"mortality"`
	Lapse          int    `json:"lapse"`
	Morbidity      int    `json:"morbidity"`
	Retrenchment   int    `json:"retrenchment"`
	Parameter      int    `json:"parameter"`
	YieldCurve     int    `json:"yield_curve"`
	ChangeType     string `json:"change_type"`
	Created        int64  `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type JournalTransactions struct {
	ID                               int     `json:"id" gorm:"primary_key"`
	CsmRunID                         int     `json:"csm_run_id"`
	RunDate                          string  `json:"run_date"`
	MeasurementType                  string  `json:"measurement_type"`
	ProductCode                      string  `json:"product_code"`
	IFRS17Group                      string  `json:"ifrs_17_group"`
	CsmRelease                       float64 `json:"csm_release"`
	DacRelease                       float64 `json:"dac_release"`
	RAChange                         float64 `json:"ra_change"`
	ExpectedBenefits                 float64 `json:"expected_benefits"`
	Expenses                         float64 `json:"expenses"`
	PremiumVariance                  float64 `json:"premium_variance"`
	ClaimsIncurred                   float64 `json:"claims_incurred"`
	ExpensesIncurred                 float64 `json:"expenses_incurred"`
	AmortizationAcquisitionCF        float64 `json:"amortization_acquisition_cf"`
	LossComponentFutureServiceChange float64 `json:"loss_component"`
	LossComponentOnInitialRecog      float64 `json:"loss_component_on_initial_recog"`
	LossComponentUnwind              float64 `json:"loss_component_unwind"`
	InsuranceFinanceExpense          float64 `json:"insurance_finance_expense"`
	PaaEarnedPremium                 float64 `json:"paa_earned_premium"`
	PaaLossComponent                 float64 `json:"paa_loss_component"`
	PaaLossRecoveryComponent         float64 `json:"paa_loss_recovery_component"`
	PaaLossComponentAdjustment       float64 `json:"paa_loss_component_adjustment"`
	PaaLossRecoveryAdjustment        float64 `json:"paa_loss_recovery_adjustment"`
	PaaLossRecoveryUnwind            float64 `json:"paa_loss_recovery_unwind"`
	PaaReinsurancePremium            float64 `json:"paa_reinsurance_premium"`
	ReinsuranceFlatCommission        float64 `json:"reinsurance_flat_commission"`
	ReinsuranceRecovery              float64 `json:"reinsurance_recovery"`
	ReinsuranceReinstatementPremium  float64 `json:"reinsurance_reinstatement_premium"`
	ReinsuranceProvisionalCommission float64 `json:"reinsurance_provisional_commission"`
	ReinsuranceUltimateCommission    float64 `json:"reinsurance_ultimate_commission"`
	ReinsuranceProfitCommission      float64 `json:"reinsurance_profit_commission"`
	ReinsuranceInvestmentComponent   float64 `json:"reinsurance_investment_component"`
	NonAttributableExpensesIncurred  float64 `json:"non_attributable_expenses_incurred"`
}

type ConfigRunId struct {
	RunId   int    `json:"run_id"`
	RunName string `json:"run_name"`
}

type LicRun struct {
	ID                   int    `json:"id,omitempty"`
	RunName              string `json:"run_name"`
	RunDate              string `json:"run_date"`
	OpeningBalanceDate   string `json:"opening_balance_date"`
	LicConfigurationName string `json:"lic_configuration_name"`
	LicConfigurationId   int    `json:"lic_configuration_id"`
}
type CsmRun struct {
	ID                           int        `json:"id,omitempty" gorm:"primary_key"`
	Name                         string     `json:"name,omitempty"`
	RunDate                      string     `json:"run_date,omitempty"`
	RunID                        int        `json:"run_id,omitempty"`
	UserName                     string     `json:"user_name,omitempty"`
	UserEmail                    string     `json:"user_email,omitempty"`
	PaaRunId                     int        `json:"paa_run_id,omitempty"`
	PaaEligibilityTest           bool       `json:"paa_eligibility_test"`
	PaaRunName                   string     `json:"paa_run_name,omitempty" gorm:"-"`
	MeasurementType              string     `json:"measurement_type,omitempty"`
	ConfigurationName            string     `json:"configuration_name,omitempty"`
	TransitionType               string     `json:"transition_type,omitempty"`
	FinanceYear                  int        `json:"finance_year"`
	RiskAdjustmentYear           int        `json:"risk_adjustment_year"`
	RiskAdjustmentVersion        string     `json:"risk_adjustment_version"`
	OpeningBalanceDate           string     `json:"opening_bal_date,omitempty"`
	OpeningRiskAdjustmentVersion string     `json:"opening_risk_adjustment_version"`
	LicConfig                    string     `json:"lic_config,omitempty"`
	CreationDate                 time.Time  `json:"creation_date,omitempty"`
	ProcessingStatus             string     `json:"processing_status,omitempty"`
	FailureReason                string     `json:"failure_reason,omitempty"`
	ProcessedGroups              int        `json:"processed_groups,omitempty"`
	TotalGroups                  int        `json:"total_groups,omitempty"`
	RunTime                      float64    `json:"run_time,omitempty"`
	ManualSap                    bool       `json:"manual_sap"`
	FinanceVersion               string     `json:"finance_version"`
	// Approval workflow fields
	ApprovalStatus               string     `json:"approval_status,omitempty" gorm:"default:draft"`
	ReviewedBy                   string     `json:"reviewed_by,omitempty"`
	ReviewedAt                   *time.Time `json:"reviewed_at,omitempty"`
	ReviewNotes                  string     `json:"review_notes,omitempty"`
	ApprovedBy                   string     `json:"approved_by,omitempty"`
	ApprovedAt                   *time.Time `json:"approved_at,omitempty"`
	ApproveNotes                 string     `json:"approve_notes,omitempty"`
	LockedBy                     string     `json:"locked_by,omitempty"`
	LockedAt                     *time.Time `json:"locked_at,omitempty"`
	ReturnReason                 string     `json:"return_reason,omitempty"`
}

// CsmRunApprovalRequest is the request body for approval workflow actions
type CsmRunApprovalRequest struct {
	Notes  string `json:"notes"`
	Reason string `json:"reason"`
}

type JTReportEntry struct {
	ProductCode          string  `json:"product_code"`
	IFRS17Group          string  `json:"ifrs17_group"`
	PostingDate          string  `json:"posting_date"`
	MeasurementModel     string  `json:"measurement_model"`
	AosStep              int     `json:"aos_step"`
	AccountNumber        int     `json:"account_number"`
	AccountTitle         string  `json:"account_title"`
	MasterAccountType    string  `json:"master_account_type"`
	Debit                float64 `json:"debit"`
	Credit               float64 `json:"credit"`
	NormalAccountBalance string  `json:"normal_account_balance"`
	AccountType          string  `json:"account_type"`
	ReportBundle         int     `json:"report_bundle"`
}

type JournalEntryResult struct {
	RunDate string
}

type SubLedgerReportEntry struct {
	ProductCode                         string  `json:"product_code"`
	IFRS17Group                         string  `json:"ifrs17_group"`
	PostingDate                         string  `json:"posting_date"`
	AccountNumber                       int     `json:"account_number"`
	LedgerName                          string  `json:"ledger_name"`
	MasterAccountType                   string  `json:"master_account_type"`
	ContraAccountTransactionDescription string  `json:"contra_account_transaction_description"`
	PostReference                       string  `json:"post_reference"`
	Debit                               float64 `json:"debit"`
	Credit                              float64 `json:"credit"`
	Balance                             float64 `json:"balance"`
}

type TrialBalanceReportEntry struct {
	ProductCode   string  `json:"product_code"`
	IFRS17Group   string  `json:"ifrs17_group"`
	PostingDate   string  `json:"posting_date"`
	AccountNumber int     `json:"account_number"`
	AccountName   string  `json:"account_name"`
	Debit         float64 `json:"debit"`
	Credit        float64 `json:"credit"`
	Amount        float64 `json:"amount"`
}

type BalanceSheetSummaryRecord struct {
	AccountName  string  `json:"account_name"`
	CurrentYear  float64 `json:"current_year"`
	PreviousYear float64 `json:"previous_year"`
	Notes        string  `json:"notes"`
}

type BalanceSheetRecord struct {
	ID                                  int     `json:"id" gorm:"primary_key"`
	CsmRunID                            int     `json:"csm_run_id"`
	ProductCode                         string  `json:"product_code" csv:"product_code"`
	MeasurementType                     string  `json:"measurement_type" csv:"measurement_type"`
	IFRS17Group                         string  `json:"ifrs17_group" csv:"ifrs17_group"`
	Date                                string  `json:"date" csv:"date" `
	BELOutflow                          float64 `json:"bel_outflow" csv:"bel_outflow" `
	BELInflow                           float64 `json:"bel_inflow" csv:"bel_inflow" `
	RiskAdjustment                      float64 `json:"risk_adjustment" csv:"risk_adjustment" `
	Treaty1BELOutflow                   float64 `json:"treaty1_bel_outflow" csv:"treaty1_bel_outflow" `
	Treaty1BELInflow                    float64 `json:"treaty1_bel_inflow" csv:"treaty1_bel_inflow" `
	Treaty1RiskAdjustment               float64 `json:"treaty1_risk_adjustment" csv:"treaty1_risk_adjustment" `
	Treaty2BELOutflow                   float64 `json:"treaty2_bel_outflow" csv:"treaty2_bel_outflow" `
	Treaty2BELInflow                    float64 `json:"treaty2_bel_inflow" csv:"treaty2_bel_inflow" `
	Treaty2RiskAdjustment               float64 `json:"treaty2_risk_adjustment" csv:"treaty2_risk_adjustment" `
	Treaty3BELOutflow                   float64 `json:"treaty3_bel_outflow" csv:"treaty3_bel_outflow" `
	Treaty3BELInflow                    float64 `json:"treaty3_bel_inflow" csv:"treaty3_bel_inflow" `
	Treaty3RiskAdjustment               float64 `json:"treaty3_risk_adjustment" csv:"treaty3_risk_adjustment" `
	PostTransitionCsm                   float64 `json:"post_transition_csm" csv:"post_transition_csm" `
	PostTransitionLc                    float64 `json:"post_transition_lc" csv:"post_transition_lc"`
	PostTransitionDAC                   float64 `json:"post_transition_dac" csv:"post_transition_dac"`
	FullyRetrospectiveCsm               float64 `json:"fully_retrospective_csm" csv:"fully_retrospective_csm"`
	FullyRetrospectiveLc                float64 `json:"fully_retrospective_lc" csv:"fully_retrospective_lc"`
	FullyRetrospectiveDAC               float64 `json:"fully_retrospective_dac" csv:"fully_retrospective_dac"`
	ModifiedRetrospectiveCsm            float64 `json:"modified_retrospective_csm" csv:"modified_retrospective_csm"`
	ModifiedRetrospectiveLc             float64 `json:"modified_retrospective_lc" csv:"modified_retrospective_lc"`
	ModifiedRetrospectiveDAC            float64 `json:"modified_retrospective_dac" csv:"modified_retrospective_dac"`
	PostTransitionTreaty1Csm            float64 `json:"post_transition_treaty1_csm" csv:"post_transition_treaty1_csm"`
	PostTransitionTreaty2Csm            float64 `json:"post_transition_treaty2_csm" csv:"post_transition_treaty2_csm"`
	PostTransitionTreaty3Csm            float64 `json:"post_transition_treaty3_csm" csv:"post_transition_treaty3_csm"`
	FullyRetrospectiveTreaty1Csm        float64 `json:"fully_retrospective_treaty1_csm" csv:"fully_retrospective_treaty1_csm"`
	FullyRetrospectiveTreaty2Csm        float64 `json:"fully_retrospective_treaty2_csm" csv:"fully_retrospective_treaty2_csm"`
	FullyRetrospectiveTreaty3Csm        float64 `json:"fully_retrospective_treaty3_csm" csv:"fully_retrospective_treaty3_csm"`
	ModifiedRetrospectiveTreaty1Csm     float64 `json:"modified_retrospective_treaty1_csm" csv:"modified_retrospective_treaty1_csm"`
	ModifiedRetrospectiveTreaty2Csm     float64 `json:"modified_retrospective_treaty2_csm" csv:"modified_retrospective_treaty2_csm"`
	ModifiedRetrospectiveTreaty3Csm     float64 `json:"modified_retrospective_treaty3_csm" csv:"modified_retrospective_treaty3_csm"`
	FairValueCsm                        float64 `json:"fair_value_csm" csv:"fair_value_csm"`
	FairValueLc                         float64 `json:"fair_value_lc" csv:"fair_value_lc"`
	FairValueDAC                        float64 `json:"fair_value_dac" csv:"fair_value_dac"`
	FairValueTreaty1Csm                 float64 `json:"fair_value_treaty1_csm" csv:"fair_value_treaty1_csm"`
	FairValueTreaty2Csm                 float64 `json:"fair_value_treaty2_csm" csv:"fair_value_treaty2_csm"`
	FairValueTreaty3Csm                 float64 `json:"fair_value_treaty3_csm" csv:"fair_value_treaty3_csm"`
	PostTransitionLossRecoveryComponent float64 `json:"post_transition_loss_recovery_component" csv:"post_transition_loss_recovery_component"`
	FRALossRecoveryComponent            float64 `json:"fra_loss_recovery_component" csv:"fra_loss_recovery_component"`
	MRALossRecoveryComponent            float64 `json:"mra_loss_recovery_component" csv:"mra_loss_recovery_component"`
	FVALossRecoveryComponent            float64 `json:"fva_loss_recovery_component" csv:"fva_loss_recovery_component"`
	PaaLiabilityRemainingCoverage       float64 `json:"paa_liability_remaining_coverage" csv:"paa_liability_remaining_coverage"`
	PaaDAC                              float64 `json:"paa_dac" csv:"paa_dac"`
	PaaReinsuranceLRC                   float64 `json:"paa_reinsurance_lrc" csv:"paa_reinsurance_lrc"`
	PaaReinsuranceDAC                   float64 `json:"paa_reinsurance_dac" csv:"paa_reinsurance_dac"`
	PaaLossComponent                    float64 `json:"paa_loss_component" csv:"paa_loss_component"`
	PaaLossRecoveryComponent            float64 `json:"paa_loss_recovery_component" csv:"paa_loss_recovery_component"`
	IBNR                                float64 `json:"ibnr" csv:"ibnr"`
	IBNRRiskAdjustment                  float64 `json:"ibnr_risk_adjustment" csv:"ibnr_risk_adjustment"`
	OutstandingClaimsReserve            float64 `json:"outstanding_claims_reserve" csv:"outstanding_claims_reserve"`
	CashBack                            float64 `json:"cash_back" csv:"cash_back"`
	Treaty1IBNR                         float64 `json:"treaty1_ibnr" csv:"treaty1_ibnr"`
	Treaty1IBNRRiskAdjustment           float64 `json:"treaty1_ibnr_risk_adjustment" csv:"treaty1_ibnr_risk_adjustment"`
	Treaty2IBNR                         float64 `json:"treaty2_ibnr" csv:"treaty1_ibnr"`
	Treaty2IBNRRiskAdjustment           float64 `json:"treaty2_ibnr_risk_adjustment" csv:"treaty1_ibnr_risk_adjustment"`
	Treaty3IBNR                         float64 `json:"treaty3_ibnr" csv:"treaty1_ibnr"`
	Treaty3IBNRRiskAdjustment           float64 `json:"treaty3_ibnr_risk_adjustment" csv:"treaty1_ibnr_risk_adjustment"`
	Year                                int     `json:"year" csv:"year"`
	Version                             string  `json:"version" csv:"version"`
	Created                             int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type PAABuildUp struct {
	ID                                    int     `json:"-" gorm:"primary_key"`
	Name                                  string  `json:"name"`
	RunDate                               string  `json:"run_date,omitempty"`
	RunId                                 int     `json:"run_id,omitempty"`
	PortfolioName                         string  `json:"portfolio_name,omitempty"`
	ProductCode                           string  `json:"product_code,omitempty"`
	IFRS17Group                           string  `json:"ifrs17_group,omitempty"`
	VariableChange                        float64 `json:"variable_change"`
	PaaLrcBuildup                         float64 `json:"paa_lrc_buildup"`
	DacBuildup                            float64 `json:"dac_buildup"`
	ModifiedGMMBel                        float64 `json:"modified_gmm_bel"`
	ModifiedGMMRiskAdjustment             float64 `json:"modified_gmm_risk_adjustment"`
	ModifiedGMMReserve                    float64 `json:"modified_gmm_reserve"`
	InitialRecognitionLossComponent       float64 `json:"initial_recognition_loss_component"`
	LossComponentUnwind                   float64 `json:"loss_component_unwind"`
	LossComponentBuildup                  float64 `json:"loss_component_buildup"`
	LossComponentAdjustment               float64 `json:"loss_component_adjustment"`
	InsuranceRevenue                      float64 `json:"insurance_revenue"`
	InsuranceServiceExpense               float64 `json:"insurance_service_expense"`
	PaaReinsuranceLrcBuildup              float64 `json:"paa_reinsurance_lrc_buildup"`
	PaaReinsuranceDacBuildup              float64 `json:"paa_reinsurance_dac_buildup"`
	InitialRecognitionLossRecovery        float64 `json:"initial_recognition_loss_recovery"`
	LossRecoveryUnwind                    float64 `json:"loss_recovery_unwind"`
	LossRecoveryBuildup                   float64 `json:"loss_recovery_buildup"`
	LossRecoveryAdjustment                float64 `json:"loss_recovery_adjustment"`
	CoverageUnits                         float64 `json:"coverage_units"`
	AcquisitionCostAmortizationProportion float64 `json:"acquisition_cost_amortization_proportion"`
	EarnedPremiumProportion               float64 `json:"earned_premium_proportion"`
}

type LicJournalTransactions struct {
	ID                       int     `json:"id" gorm:"primary_key"`
	LicRunID                 int     `json:"csm_run_id"`
	RunDate                  string  `json:"run_date"`
	ProductCode              string  `json:"product_code"`
	IFRS17Group              string  `json:"ifrs_17_group"`
	LicFutureCashFlowsChange float64 `json:"lic_future_cash_flows_change"`
	LicExperienceVariance    float64 `json:"lic_experience_variance"`
	IBNRIncurredClaims       float64 `json:"ibnr_incurred_claims"`
}
