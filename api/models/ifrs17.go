package models

type BaseAosVariable struct {
	ID               int    `json:"id" gorm:"primary_key"`
	Name             string `json:"name" gorm:"unique_index:idx_aos_variable"`
	Time             string `json:"time"`
	Description      string `json:"description"`
	RunId            int    `json:"run_id"`
	RunName          string `json:"run_name"`
	InsuranceService string `json:"insurance_service"`
	Notes            string `json:"notes"`
	AssumptionBasis  string `json:"assumption_basis"`
}

type AosVariable struct {
	ID                 int    `json:"-" gorm:"primary_key"`
	AosVariableSetID   int    `json:"-"`
	AosVariableSetName string `json:"aos_variable_set_name"`
	CsmRunId           int    `json:"csm_run_id"`
	Name               string `json:"name"`
	Time               string `json:"time"`
	Description        string `json:"description"`
	RunId              int    `json:"run_id"`
	RunName            string `json:"run_name"`
	InsuranceService   string `json:"insurance_service"`
	Notes              string `json:"notes"`
	AssumptionBasis    string `json:"assumption_basis"`
	ExternalSapSource  bool   `json:"external_sap_source"`
}

type AosVariableSet struct {
	ID                 int           `json:"id" gorm:"primary_key"`
	ConfigurationName  string        `json:"configuration_name"`
	CoverageUnitOption string        `json:"coverage_unit_option"`
	ExternalSap        bool          `json:"external_sap"`
	AosVariables       []AosVariable `json:"aos_variables"`
}

type CsmAosVariable struct {
	ID               int    `json:"-" gorm:"primary_key"`
	CsmRunId         int    `json:"csm_run_id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	RunId            int    `json:"run_id"`
	InsuranceService string `json:"insurance_service"`
	Notes            string `json:"notes"`
	AssumptionBasis  string `json:"assumption_basis"`
}

type AOSStepResult struct {
	ID                            int     `json:"id,omitempty" gorm:"primary_key"`
	CsmRunID                      int     `json:"csm_run_id,omitempty"`
	RunDate                       string  `json:"run_date,omitempty"`
	ProductCode                   string  `json:"product_code,omitempty" gorm:"unique_index:idx_step"`
	IFRS17Group                   string  `json:"ifrs17_group,omitempty" gorm:"unique_index:idx_step"`
	Name                          string  `json:"name,omitempty" gorm:"unique_index:idx_step"`
	Description                   string  `json:"description,omitempty"`
	Time                          string  `json:"time,omitempty"`
	RunId                         int     `json:"run_id,omitempty"`
	BelOutflow                    float64 `json:"bel_outflow"`
	BelOutflowExclAcquisition     float64 `json:"bel_outflow_excl_acquisition"`
	BelAcquisitionCost            float64 `json:"bel_acquisition_cost"`
	BelInflow                     float64 `json:"bel_inflow"`
	BEL                           float64 `json:"bel"`
	BelInflowAt12                 float64 `json:"bel_inflow_at12"`
	BelOutflowAt12                float64 `json:"bel_outflow_at12"`
	BelOutflowExclAcquisitionAt12 float64 `json:"bel_outflow_excl_acquisition_at12"`
	BelAcquisitionCostAt12        float64 `json:"bel_acquisition_cost_at12"`
	BELAt12                       float64 `json:"bel_at12"`
	RiskAdjustment                float64 `json:"risk_adjustment"`
	RiskAdjustmentAt12            float64 `json:"risk_adjustment_at12"`
	BestEstimateLiabilityChange   float64 `json:"best_estimate_liability_change"`
	BelInflowChange               float64 `json:"bel_inflow_change"`
	BelOutflowChange              float64 `json:"bel_outflow_change"`
	RiskAdjustmentChange          float64 `json:"risk_adjustment_change"`
	LiabilityChange               float64 `json:"liability_change"`
	DACChange                     float64 `json:"dac_change"`
	CSMChange                     float64 `json:"csm_change"`
	LossComponentChange           float64 `json:"loss_component_change"`
	LossComponentUnwind           float64 `json:"loss_component_unwind"`
	PNLChange                     float64 `json:"pnl_change"`
	CSMBuildup                    float64 `json:"csm_buildup"`
	RiskAdjustmentBuildup         float64 `json:"risk_adjustment_buildup"`
	LossComponentBuildup          float64 `json:"loss_component_buildup"`
	DACBuildup                    float64 `json:"dac_buildup"`
	BelOutflowExclLC              float64 `json:"bel_outflow_excl_lc"`
	RiskAdjExclLC                 float64 `json:"risk_adj_excl_lc"`
	BelOutflowLC                  float64 `json:"bel_outflow_lc"`
	RiskAdjLC                     float64 `json:"risk_adj_lc"`
	SarActualClaimNetLc           float64 `json:"sar_actual_claim_net_lc"`
	SarActualExpenseNetLc         float64 `json:"sar_actual_expense_net_lc"`
	LcSar                         float64 `json:"lc_sar"`
	ActualNonAttributableExpenses float64 `json:"actual_non_attributable_expenses"`
	ExpectedRaNetOfLc             float64 `json:"expected_ra_net_of_lc"`
	SumCoverageUnits              float64 `json:"sum_coverage_units"`
	DiscountedCoverageUnits       float64 `json:"discounted_coverage_units"`
	SumCoverageUnitsAt12          float64 `json:"sum_coverage_units_at12"`
	DiscountedCoverageUnitsAt12   float64 `json:"discounted_coverage_units_at12"`
	CSMReleaseRatio               float64 `json:"csm_release_ratio"`
	CSMRelease                    float64 `json:"csm_release"`
	InterestAccretionFac          float64 `json:"interest_accretion_fac"`
	ExpectedCashInflow            float64 `json:"expected_cash_inflow"`
	ExpectedCashOutflow           float64 `json:"expected_cash_outflow"`
	ActualPremium                 float64 `json:"actual_premium"`
	PremiumDebtor                 float64 `json:"premium_debtor"`
	ExperiencePremiumVariance     float64 `json:"experience_premium_variance"`
	ReinsuranceBelOutflow         float64 `json:"reinsurance_bel_outflow"`
	ReinsuranceBelInflow          float64 `json:"reinsurance_bel_inflow"`
	ReinsuranceBel                float64 `json:"reinsurance_bel"`
	ReinsuranceRiskAdjustment     float64 `json:"reinsurance_risk_adjustment"`
	ReinsuranceBelOutflowAt12     float64 `json:"reinsurance_bel_outflow_at12"`
	ReinsuranceBelInflowAt12      float64 `json:"reinsurance_bel_inflow_at12"`
	ReinsuranceBelAt12            float64 `json:"reinsurance_bel_at12"`
	ReinsuranceRiskAdjustmentAt12 float64 `json:"reinsurance_risk_adjustment_at12"`
	ReinsBelOutflowChange         float64 `json:"reins_bel_outflow_change"`
	ReinsBelInflowChange          float64 `json:"reins_bel_inflow_change"`
	ReinsBelChange                float64 `json:"reins_bel_change"`
	ReinsRAChange                 float64 `json:"reins_ra_change"`
	ReinsCSMChange                float64 `json:"reins_csm_change"`
	ReinsCSMBuildup               float64 `json:"reins_csm_buildup"`
	ReinsPNLChange                float64 `json:"reins_pnl_change"`
	ReinsuranceCSM                float64 `json:"reinsurance_csm"`
	LossRecoveryComponentChange   float64 `json:"loss_recovery_component_change"`
	LossRecoveryComponentBuildup  float64 `json:"loss_recovery_component_buildup"`
	ReinsuranceRevenue            float64 `json:"reinsurance_revenue"`
	ReinsuranceServiceExpense     float64 `json:"reinsurance_service_expense"`
	ReinsuranceServiceResult      float64 `json:"reinsurance_service_result"`
	ReinsExpectedCashOutflow      float64 `json:"reins_expected_cash_outflow"`
	ReinsExpectedCashInflow       float64 `json:"reins_expected_cash_inflow"`
	IsOnerous                     bool    `json:"is_onerous"`
	OnerousReason                 string  `json:"onerous_reason"`
}

type PAAResult struct {
	ID                            int     `json:"id,omitempty" gorm:"primary_key"`
	CsmRunID                      int     `json:"csm_run_id,omitempty"`
	RunDate                       string  `json:"run_date,omitempty"`
	PortfolioName                 string  `json:"portfolio_name,omitempty"`
	ProductCode                   string  `json:"product_code" gorm:"unique_index:idx_step"`
	IFRS17Group                   string  `json:"ifrs17_group,omitempty" gorm:"unique_index:idx_step"`
	Treaty1IFRS17Group            string  `json:"treaty1_ifrs17_group,omitempty" gorm:"unique_index:idx_step"`
	Treaty2IFRS17Group            string  `json:"treaty2_ifrs17_group,omitempty" gorm:"unique_index:idx_step"`
	Treaty3IFRS17Group            string  `json:"treaty3_ifrs17_group,omitempty" gorm:"unique_index:idx_step"`
	Time                          string  `json:"time,omitempty"`
	RunId                         int     `json:"run_id,omitempty"`
	PremiumReceipt                float64 `json:"premium_receipt"`
	NBPremiumReceipt              float64 `json:"nb_premium_receipt"`
	TotalPremiumReceipt           float64 `json:"total_premium_receipt"`
	AcquisitionExpenses           float64 `json:"acquisition_expenses"`
	DeferredAcquisitionExpenses   float64 `json:"deferred_acquisition_expenses"`
	AmortisedAcquisitionCost      float64 `json:"amortised_acquisition_cost"`
	EarnedPremium                 float64 `json:"earned_premium"`
	PaaLiabilityRemainingCoverage float64 `json:"paa_liability_remaining_coverage"`
	GmmBel                        float64 `json:"gmm_bel"`
	GmmRiskAdjustment             float64 `json:"gmm_risk_adjustment"`
	GmmReserve                    float64 `json:"gmm_reserve"`
	InsuranceRevenue              float64 `json:"insurance_revenue"`
	IncurredClaims                float64 `json:"incurred_claims"`
	ClaimsPaid                    float64 `json:"claims_paid"`
	IncurredExpenses              float64 `json:"incurred_expenses"`
	PaaLossComponent              float64 `json:"paa_loss_component"`
	ReinsurancePremium            float64 `json:"reinsurance_premium"`
	//Treaty1PremiumPaid                 float64 `json:"treaty1_premium_paid"`
	//Treaty2PremiumPaid                 float64 `json:"treaty2_premium_paid"`
	//Treaty3PremiumPaid                 float64 `json:"treaty3_premium_paid"`
	PaaReinsuranceLrc float64 `json:"paa_reinsurance_lrc"`
	//PaaTreaty1Lrc                      float64 `json:"paa_treaty1_lrc"`
	//PaaTreaty2Lrc                      float64 `json:"paa_treaty2_lrc"`
	//PaaTreaty3Lrc                      float64 `json:"paa_treaty3_lrc"`
	//NonPerformanceTreaty1Lrc           float64 `json:"non_performance_treaty1_lrc"`
	//NonPerformanceTreaty2Lrc           float64 `json:"non_performance_treaty2_lrc"`
	NonPerformanceReinsuranceLrc float64 `json:"non_performance_reinsurance_lrc"`
	PaaReinsuranceDac            float64 `json:"paa_reinsurance_dac"`
	PaaLossRecoveryComponent     float64 `json:"paa_loss_recovery_component"`
	//PaaTreaty1LossRecoveryComponent   float64 `json:"paa_treaty1_loss_recovery_component"`
	//PaaTreaty2LossRecoveryComponent   float64 `json:"paa_treaty2_loss_recovery_component"`
	//PaaTreaty3LossRecoveryComponent   float64 `json:"paa_treaty3_loss_recovery_component"`
	PaaInitialLossRecoveryComponent        float64 `json:"paa_initial_loss_recovery_component"`
	PaaTreaty1InitialLossRecoveryComponent float64 `json:"paa_treaty1_initial_loss_recovery_component"`
	PaaTreaty2InitialLossRecoveryComponent float64 `json:"paa_treaty2_initial_loss_recovery_component"`
	PaaTreaty3InitialLossRecoveryComponent float64 `json:"paa_treaty3_initial_loss_recovery_component"`
	PaaLossRecoveryUnwind                  float64 `json:"paa_loss_recovery_unwind"`
	PaaTreaty1LossRecoveryUnwind           float64 `json:"paa_treaty1_loss_recovery_unwind"`
	PaaTreaty2LossRecoveryUnwind           float64 `json:"paa_treaty2_loss_recovery_unwind"`
	PaaTreaty3LossRecoveryUnwind           float64 `json:"paa_treaty3_loss_recovery_unwind"`

	AllocatedReinsurancePremium        float64 `json:"allocated_reinsurance_premium"`
	AllocatedTreaty1Premium            float64 `json:"allocated_treaty1_premium"`
	AllocatedTreaty2Premium            float64 `json:"allocated_treaty2_premium"`
	AllocatedTreaty3Premium            float64 `json:"allocated_treaty3_premium"`
	AllocatedReinsuranceFlatCommission float64 `json:"allocated_reinsurance_flat_commission"`
	//ReinsuranceClaims                  float64 `json:"reinsurance_claims"`
	ReinsuranceRecovery              float64 `json:"reinsurance_recovery"`
	Treaty1Recovery                  float64 `json:"treaty1_recovery"`
	Treaty2Recovery                  float64 `json:"treaty2_recovery"`
	Treaty3Recovery                  float64 `json:"treaty3_recovery"`
	ReinsuranceReinstatementPremium  float64 `json:"reinsurance_reinstatement_premium"`
	ReinsuranceUltimateLossRatio     float64 `json:"reinsurance_ultimate_loss_ratio"`
	ReinsuranceProvisionalCommission float64 `json:"reinsurance_provisional_commission"`
	ReinsuranceUltimateCommission    float64 `json:"reinsurance_ultimate_commission"`
	ReinsuranceProfitCommission      float64 `json:"reinsurance_profit_commission"`
	ReinsuranceTotalPaidToCedant     float64 `json:"reinsurance_total_paid_to_cedant"`
	ReinsuranceInvestmentComponent   float64 `json:"reinsurance_investment_component"`
	ReinsuranceRevenue               float64 `json:"reinsurance_revenue"`
	ReinsuranceServiceExpense        float64 `json:"reinsurance_service_expense"`
	ReinsuranceServiceResult         float64 `json:"reinsurance_service_result"`
	NonAttributableExpenses          float64 `json:"non_attributable_expenses"`
	IsOnerous                        bool    `json:"is_onerous"`
	OnerousReason                    string  `json:"onerous_reason"`
}

type PAAEligibilityTestResult struct {
	ID                        int     `json:"id,omitempty" gorm:"primary_key"`
	CsmRunID                  int     `json:"csm_run_id,omitempty"`
	ProjectionMonth           int     `json:"projection_month"`
	RunDate                   string  `json:"run_date,omitempty"`
	ProductCode               string  `json:"product_code" gorm:"unique_index:idx_step"`
	IFRS17Group               string  `json:"ifrs17_group,omitempty" gorm:"unique_index:idx_step"`
	Time                      string  `json:"time,omitempty"`
	RunId                     int     `json:"run_id,omitempty"`
	AmortisedAcquisitionCost  float64 `json:"amortised_acquisition_cost"`
	EarnedPremium             float64 `json:"earned_premium"`
	UnearnedPremiumReserve    float64 `json:"unearned_premium_reserve"`
	DeferredAcquisitionCost   float64 `json:"deferred_acquisition_cost"`
	PaaLossComponent          float64 `json:"paa_loss_component"`
	GmmBel                    float64 `json:"gmm_bel"`
	GmmRiskAdjustment         float64 `json:"gmm_risk_adjustment"`
	GmmReserve                float64 `json:"gmm_reserve"`
	GMMCsm                    float64 `json:"gmm_csm"`
	GMMCsmRelease             float64 `json:"gmm_csm_release"`
	GMMDacRelease             float64 `json:"gmm_dac_release"`
	PAALossComponentRelease   float64 `json:"paa_loss_component_release"`
	GMMExpectedOutflows       float64 `json:"gmm_expected_outflows"`
	GMMRiskAdjustmentRelease  float64 `json:"gmm_risk_adjustment_release"`
	PAARevenue                float64 `json:"paa_revenue"`
	GMMRevenue                float64 `json:"gmm_revenue"`
	RevenueVariance           float64 `json:"revenue_variance"`
	RevenueVarianceProportion float64 `json:"revenue_variance_proportion"`
	PAALrc                    float64 `json:"paa_lrc"`
	GMMLrc                    float64 `json:"gmm_lrc"`
	LRCVariance               float64 `json:"lrc_variance"`
	LRCVarianceProportion     float64 `json:"lrc_variance_proportion"`
}

type IncomeStatementEntry struct {
	LineItem    string  `json:"line_item"`
	Index       int     `json:"index"`
	PastYear    float64 `json:"past_year"`
	CurrentYear float64 `json:"current_year"`
	Notes       string  `json:"notes"`
	Reference   string  `json:"reference"`
	Style       string  `json:"style"`
	Type        string  `json:"type"`
}

type CSMModelPoint struct {
	Reserves                      float64 `json:"reserves"`
	PremiumIncome                 float64 `json:"premium_income"`
	AcquisitionExpenses           float64 `json:"acquisition_expenses"`
	MaintenanceExpenses           float64 `json:"maintenance_expenses"`
	DeathOutgo                    float64 `json:"death_outgo"`
	DisabilityOutgo               float64 `json:"disability_outgo"`
	RetrenchmentOutgo             float64 `json:"retrenchment_outgo"`
	CoverageUnits                 float64 `json:"coverage_units"`
	SumFutureCoverageUnits        float64 `json:"sum_future_coverage_units"`
	DiscountedFutureCoverageUnits float64 `json:"discounted_future_coverage_units"`
	DiscountedExpenses            float64 `json:"discounted_expenses"`
	DiscountedCashflowOutgo       float64 `json:"discounted_cashflow_outgo"`
	DiscountedDeathOutgo          float64 `json:"discounted_death_outgo"`
	DiscountedMorbidityOutgo      float64 `json:"discounted_morbidity_outgo"`
	DiscountedRetrenchmentOutgo   float64 `json:"discounted_retrenchment_outgo"`
	DiscountedSurrenderOutgo      float64 `json:"discounted_surrender_outgo"`
	DiscountedAnnuityIncome       float64 `json:"discounted_annuity_income"`
	CsmCarriedForward             float64 `json:"csm_carried_forward"`
}

type ProductList struct {
	ProductCode string `json:"product_code"`
}

type IFRS17List struct {
	Group string `json:"group"`
}

type GroupResults struct {
	IFRS17Group   string ` json:"ifrs17_group" gorm:"column:ifrs17_group"`
	ProductCode   string `json:"product_code" gorm:"product_code"`
	RunDate       string `json:"run_date" gorm:"run_date"`
	PortfolioName string `json:"portfolio_name" gorm:"portfolio_name"`
}

type ProductStepResult struct {
	ProductCode                   string  `json:"product_code"`
	IFRS17Group                   string  `json:"ifrs17_group,omitempty"`
	Name                          string  `json:"name"`
	BelInflow                     float64 `json:"bel_inflow"`
	BelOutflow                    float64 `json:"bel_outflow"`
	Bel                           float64 `json:"bel"`
	BelInflowAt12                 float64 `json:"bel_inflow_at12"`
	BelOutflowAt12                float64 `json:"bel_outflow_at12"`
	BelAt12                       float64 `json:"bel_at12"`
	RiskAdjustment                float64 `json:"risk_adjustment"`
	RiskAdjustmentAt12            float64 `json:"risk_adjustment_at12"`
	BestEstimateLiabilityChange   float64 `json:"best_estimate_liability_change"`
	BELInflowChange               float64 `json:"bel_inflow_change"`
	BELOutflowChange              float64 `json:"bel_outflow_change"`
	RiskAdjustmentChange          float64 `json:"risk_adjustment_change"`
	LiabilityChange               float64 `json:"liability_change"`
	DACChange                     float64 `json:"dac_change"`
	CsmChange                     float64 `json:"csm_change"`
	LossComponentChange           float64 `json:"loss_component_change"`
	LossComponentUnwind           float64 `json:"loss_component_unwind"`
	PnlChange                     float64 `json:"pnl_change"`
	CsmBuildup                    float64 `json:"csm_buildup"`
	RiskAdjustmentBuildup         float64 `json:"risk_adjustment_buildup"`
	LossComponentBuildup          float64 `json:"loss_component_buildup"`
	SarActualClaimNetLc           float64 `json:"sar_actual_claim_net_lc"`
	SarActualExpenseNetLc         float64 `json:"sar_actual_expense_net_lc"`
	ExpectedRaNetOfLc             float64 `json:"expected_ra_net_of_lc"`
	LcSar                         float64 `json:"lc_sar"`
	ActualNonAttributableExpenses float64 `json:"actual_non_attributable_expenses"`
	SumCoverageUnits              float64 `json:"sum_coverage_units"`
	DiscountedCoverageUnits       float64 `json:"discounted_coverage_units"`
	CSMReleaseRatio               float64 `json:"csm_release_ratio"`
	CSMRelease                    float64 `json:"csm_release"`
	InterestAccretionFac          float64 `json:"interest_accretion_fac"`
	ExpectedCashInflow            float64 `json:"expected_cash_inflow"`
	ExpectedCashOutflow           float64 `json:"expected_cash_outflow"`
	ActualPremium                 float64 `json:"actual_premium"`
	ExperiencePremiumVariance     float64 `json:"experience_premium_variance"`
	ReinsuranceBelOutflow         float64 `json:"reinsurance_bel_outflow"`
	ReinsuranceBelInflow          float64 `json:"reinsurance_bel_inflow"`
	ReinsuranceBel                float64 `json:"reinsurance_bel"`
	ReinsuranceRiskAdjustment     float64 `json:"reinsurance_risk_adjustment"`
	ReinsuranceBelOutflowAt12     float64 `json:"reinsurance_bel_outflow_at12"`
	ReinsuranceBelInflowAt12      float64 `json:"reinsurance_bel_inflow_at12"`
	ReinsuranceBelAt12            float64 `json:"reinsurance_bel_at12"`
	ReinsuranceRiskAdjustmentAt12 float64 `json:"reinsurance_risk_adjustment_at12"`
	ReinsBelOutflowChange         float64 `json:"reins_bel_outflow_change"`
	ReinsBelInflowChange          float64 `json:"reins_bel_inflow_change"`
	ReinsBelChange                float64 `json:"reins_bel_change"`
	ReinsRAChange                 float64 `json:"reins_ra_change"`
	ReinsCSMChange                float64 `json:"reins_csm_change"`
	ReinsCSMBuildup               float64 `json:"reins_csm_buildup"`
	ReinsPNLChange                float64 `json:"reins_pnl_change"`
	ReinsuranceCSM                float64 `json:"reinsurance_csm"`
	LossRecoveryComponentChange   float64 `json:"loss_recovery_component_change"`
	LossRecoveryComponentBuildup  float64 `json:"loss_recovery_component_buildup"`
	ReinsuranceRevenue            float64 `json:"reinsurance_revenue"`
	ReinsuranceServiceExpense     float64 `json:"reinsurance_service_expense"`
	ReinsuranceServiceResult      float64 `json:"reinsurance_service_result"`
	ReinsExpectedCashOutflow      float64 `json:"reins_expected_cash_outflow"`
	ReinsExpectedCashInflow       float64 `json:"reins_expected_cash_inflow"`
}

type LiabilityMovement struct {
	ID                    int     `json:"-" gorm:"primary_key"`
	CsmRunID              int     `json:"csm_run_id,omitempty"`
	RunDate               string  `json:"run_date"`
	ProductCode           string  `json:"product_code,omitempty"`
	Ifrs17Group           string  `json:"ifrs17_group,omitempty"`
	Code                  int     `json:"code"`
	Variable              string  `json:"variable"`
	Bel                   float64 `json:"bel"`
	RiskAdjustment        float64 `json:"risk_adjustment"`
	Csm                   float64 `json:"csm"`
	TotalLRC              float64 `json:"total_lrc"`
	IncurredClaimsBel     float64 `json:"incurred_claims_bel"`
	IncurredClaimsRiskAdj float64 `json:"incurred_claims_risk_adj"`
	TotalLIC              float64 `json:"total_lic"`
	FRACsm                float64 `json:"fra_csm"`
	MRACsm                float64 `json:"mra_csm"`
	FVACsm                float64 `json:"fva_csm"`
	LossComponent         float64 `json:"loss_component"`
	FRALossComponent      float64 `json:"fra_loss_component"`
	MRALossComponent      float64 `json:"mra_loss_component"`
	FVALossComponent      float64 `json:"fva_loss_component"`
	Reference             string  `json:"reference"`
}

type LiabilityMovementLine struct {
	Code                  int
	Variable              string  `json:"variable"`
	Bel                   float64 `json:"bel"`
	RiskAdjustment        float64 `json:"risk_adjustment"`
	Csm                   float64 `json:"csm"`
	TotalLRC              float64 `json:"total_lrc"`
	IncurredClaimsBel     float64 `json:"incurred_claims_bel"`
	IncurredClaimsRiskAdj float64 `json:"incurred_claims_risk_adj"`
	TotalLic              float64 `json:"total_lic"`
	FRACsm                float64 `json:"fra_csm"`
	MRACsm                float64 `json:"mra_csm"`
	FVACsm                float64 `json:"fva_csm"`
	LossComponent         float64 `json:"loss_component"`
	FRALossComponent      float64 `json:"fra_loss_component"`
	MRALossComponent      float64 `json:"mra_loss_component"`
	FVALossComponent      float64 `json:"fva_loss_component"`
	Reference             string  `json:"reference"`
}

type InsuranceRevenue struct {
	ID                            int     `json:"id" gorm:"primary_key"`
	RunDate                       string  `json:"run_date"`
	CsmRunID                      int     `json:"csm_run_id"`
	ProductCode                   string  `json:"product_code"`
	IFRS17Group                   string  `json:"ifrs17_group"`
	Code                          int     `json:"code"`
	Variable                      string  `json:"variable"`
	PostTransition                float64 `json:"post_transition"`
	FullRetrospectiveApproach     float64 `json:"full_retrospective_approach"`
	ModifiedRetrospectiveApproach float64 `json:"modified_retrospective_approach"`
	FairValueApproach             float64 `json:"fair_value_approach"`
	Reference                     string  `json:"reference"`
}

type InitialRecognition struct {
	ID                              int     `json:"id,omitempty" gorm:"primary_key"`
	RunDate                         string  `json:"run_date"`
	CsmRunID                        int     `json:"csm_run_id"`
	ProductCode                     string  `json:"product_code,omitempty"`
	IFRS17Group                     string  `json:"ifrs17_group,omitempty"`
	Code                            int     `json:"code"`
	Variable                        string  `json:"variable"`
	NoSignificantProbabilityOnerous float64 `json:"no_significant_probability_onerous"`
	Onerous                         float64 `json:"onerous"`
	Remaining                       float64 `json:"remaining"`
	BusinessAcquisitionTransfer     float64 `json:"business_acquisition_transfer"`
	Reference                       string  `json:"reference"`
}

type CsmProjection struct {
	ID                                         int     `json:"id,omitempty" gorm:"primary_key"`
	CsmRunID                                   int     `json:"csm_run_id,omitempty"`
	RunDate                                    string  `json:"run_date"`
	ProductCode                                string  `json:"product_code,omitempty"`
	IFRS17Group                                string  `json:"ifrs17_group,omitempty"`
	ProjectionMonth                            int     `json:"projection_month"`
	CsmTotal                                   float64 `json:"csm_total"`
	CsmReleaseTotal                            float64 `json:"csm_release_total"`
	CoverageUnits                              float64 `json:"coverage_units"`
	Csm                                        float64 `json:"csm"`
	CsmRelease                                 float64 `json:"csm_release"`
	CoverageUnitsFullRetrospectiveApproach     float64 `json:"coverage_units_full_retrospective_approach"`
	CsmFullRetrospectiveApproach               float64 `json:"csm_full_retrospective_approach"`
	CsmReleaseFullRetrospectiveApproach        float64 `json:"csm_release_full_retrospective_approach"`
	CoverageUnitsModifiedRetrospectiveApproach float64 `json:"coverage_units_modified_retrospective_approach"`
	CsmModifiedRetrospectiveApproach           float64 `json:"csm_modified_retrospective_approach"`
	CsmReleaseModifiedRetrospectiveApproach    float64 `json:"csm_release_modified_retrospective_approach"`
	CoverageUnitsFairValueApproach             float64 `json:"coverage_units_fair_value_approach"`
	CsmFairValueApproach                       float64 `json:"csm_fair_value_approach"`
	CsmReleaseFairValueApproach                float64 `json:"csm_release_fair_value_approach"`
	Reference                                  string  `json:"reference"`
}

type ChartAccountItem struct {
	Index                int    `json:"index" csv:"index"`
	AccountName          string `json:"account_name" csv:"account_name"`
	Description          string `json:"description" csv:"description"`
	AccountNumber        int    `json:"account_number" csv:"account_number"`
	AccountType          string `json:"account_type" csv:"account_type"`
	NormalAccountBalance string `json:"normal_account_balance" csv:"normal_account_balance"`
	MasterAccountType    string `json:"master_account_type" csv:"master_account_type"`
}

type RunComparisonMetric struct {
	Label       string  `json:"label"`
	RunAValue   float64 `json:"run_a_value"`
	RunBValue   float64 `json:"run_b_value"`
	Variance    float64 `json:"variance"`
	VariancePct float64 `json:"variance_pct"`
}
