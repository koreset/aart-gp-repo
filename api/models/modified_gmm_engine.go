package models

import (
	"strings"
	"time"
)

type CsvTime time.Time

// tryParseDate tries multiple common date layouts and returns the first match.
func tryParseDate(value string) (time.Time, error) {
    layouts := []string{
        time.RFC3339,        // 2006-01-02T15:04:05Z07:00
        "2006-01-02",       // 1986-05-25
        "2006/01/02",       // 1986/05/25
        "02/01/2006",       // 25/05/1986 (dd/mm/yyyy)
        "2/1/2006",         // 5/5/1986 (d/m/yyyy)
        "01/02/2006",       // 05/25/1986 (mm/dd/yyyy)
        "1/2/2006",         // 5/5/1986 (m/d/yyyy)
        "02/01/06",         // 25/05/86 (dd/mm/yy)
        "2/1/06",           // 5/5/86 (d/m/yy)
        "01/02/06",         // 05/25/86 (mm/dd/yy)
        "1/2/06",           // 5/5/86 (m/d/yy)
        "02-01-2006",       // 25-05-1986
        "01-02-2006",       // 05-25-1986
        "02-01-06",         // 25-05-86
        "01-02-06",         // 05-25-86
    }
    var lastErr error
    for _, layout := range layouts {
        if t, err := time.Parse(layout, value); err == nil {
            return t, nil
        } else {
            lastErr = err
        }
    }
    return time.Time{}, lastErr
}

func (c *CsvTime) UnmarshalCSV(b []byte) error {
    value := strings.Trim(string(b), `"`) //get rid of "
    value = strings.TrimSpace(value)
    if value == "" || value == "null" || value == "0" { // allow blanks/zeros
        return nil
    }

    t, err := tryParseDate(value)
    if err != nil {
        return err
    }
    *c = CsvTime(t) //set result using the pointer
    return nil
}

func (c *CsvTime) UnmarshalJSON(b []byte) error {
    value := strings.Trim(string(b), `"`) //get rid of "
    value = strings.TrimSpace(value)
    if value == "" || value == "null" || value == "0" {
        return nil
    }

    t, err := tryParseDate(value)
    if err != nil {
        return err
    }
    *c = CsvTime(t) //set result using the pointer
    return nil
}

func (c CsvTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format("2006-01-02") + `"`), nil
}

func (c CsvTime) MarshalCSV() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format("2006-01-02") + `"`), nil
}

type ModifiedGMMModelPoint struct {
	ID                        int     `json:"-" gorm:"primary_key"`
	PaaPortfolioID            int     `json:"-" gorm:"index"`
	PaaPortfolioName          string  `json:"paa_portfolio_name" gorm:"index"`
	ProductCode               string  `json:"product_code" csv:"product_code"`
	SubProductCode            string  `json:"sub_product_code" csv:"sub_product_code"`
	TermMonths                float64 `json:"term_months" csv:"term_months"`
	LockedInYear              int     `json:"locked_in_year" csv:"locked_in_year"`
	LockedInMonth             int     `json:"locked_in_month" csv:"locked_in_month"`
	PolicyNumber              string  `json:"policy_number" csv:"policy_number"`
	IFRS17Group               string  `json:"ifrs17_group" csv:"ifrs17_group"`
	IFRS17GroupTreaty1        string  `json:"ifrs17_group_treaty1" csv:"ifrs17_group_treaty1"`
	IFRS17GroupTreaty2        string  `json:"ifrs17_group_treaty2" csv:"ifrs17_group_treaty2"`
	IFRS17GroupTreaty3        string  `json:"ifrs17_group_treaty3" csv:"ifrs17_group_treaty3"`
	AnnualPremium             float64 `json:"annual_premium" csv:"annual_premium"`
	WrittenPremium            float64 `json:"written_premium" csv:"written_premium"`
	Frequency                 int     `json:"frequency" csv:"frequency"`
	DurationInForceMonths     int     `json:"duration_in_force_months" csv:"duration_in_force_months"`
	DistributionChannel       string  `json:"distribution_channel" csv:"distribution_channel"`
	CommencementDate          string  `json:"commencement_date" csv:"commencement_date"`
	CoverEndDate              string  `json:"cover_end_date" csv:"cover_end_date"`
	Status                    string  `json:"status" csv:"status"`
	OriginalLoan              float64 `json:"original_loan" cs	v:"original_loan"`
	OutstandingLoan           float64 `json:"outstanding_loan" csv:"outstanding_loan"`
	AnnualInterestRate        float64 `json:"annual_interest_rate" csv:"annual_interest_rate"`
	MonthlyInstalment         float64 `json:"monthly_instalment" csv:"monthly_instalment"`
	OutstandingLoanTermMonths int     `json:"outstanding_loan_term_months" csv:"outstanding_loan_term_months"`
	Year                      int     `json:"year" gorm:"index" csv:"year"`
	Month                     int     `json:"month" gorm:"index" csv:"month"`
	MpVersion                 string  `json:"mp_version" gorm:"index" csv:"mp_version"`
	Created                   int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
	CreatedBy                 string  `json:"created_by" csv:"created_by"`
}

type ModifiedGMMParameter struct {
	ID                                         int     `json:"-" gorm:"primary_key"`
	PortfolioName                              string  `json:"portfolio_name" csv:"portfolio_name"`
	ProductName                                string  `json:"product_name" csv:"product_name"`
	ProductCode                                string  `json:"product_code" csv:"product_code"`
	SubProductCode                             string  `json:"sub_product_code" csv:"sub_product_code"`
	YieldCurveCode                             string  `json:"yield_curve_code" csv:"yield_curve_code"`
	Year                                       int     `json:"year" csv:"year"`
	PremiumEarningPatternCode                  string  `json:"premium_earning_pattern_code" csv:"premium_earning_pattern_code"`
	ClaimsProportion                           float64 `json:"claims_proportion" csv:"claims_proportion"`
	ClaimsExpenseProportion                    float64 `json:"claims_expense_proportion" csv:"claims_expense_proportion"`
	MaintenanceExpenseProportion               float64 `json:"maintenance_expense_proportion" csv:"maintenance_expense_proportion"`
	MaintenanceAnnualExpenseAmount             float64 `json:"maintenance_annual_expense_amount" csv:"maintenance_annual_expense_amount"`
	InitialCommissionAmount                    float64 `json:"initial_commission_amount" csv:"initial_commission_amount"`
	InitialYr1CommissionProportion             float64 `json:"initial_yr_1_commission_proportion" csv:"initial_yr_1_commission_proportion"`
	InitialYr2CommissionProportion             float64 `json:"initial_yr_2_commission_proportion" csv:"initial_yr_2_commission_proportion"`
	RenewalAnnualCommissionAmount              float64 `json:"renewal_annual_commission_amount" csv:"renewal_annual_commission_amount"`
	RenewalCommissionProportion                float64 `json:"renewal_commission_proportion" csv:"renewal_commission_proportion"`
	InitialExpenseProportion                   float64 `json:"initial_expense_proportion" csv:"initial_expense_proportion"`
	InitialExpenseAmount                       float64 `json:"initial_expense_amount" csv:"initial_expense_amount"`
	IBNRProportion                             float64 `json:"ibnr_proportion" csv:"ibnr_proportion"`
	DacBuildupIndicator                        bool    `json:"dac_buildup_indicator" csv:"dac_buildup_indicator"`
	ReinsuranceTreaty1ClaimsProportion         float64 `json:"reinsurance_treaty1_claims_proportion" csv:"reinsurance_treaty1_claims_proportion"`
	ReinsuranceTreaty1PremiumProportion        float64 `json:"reinsurance_treaty1_premium_proportion" csv:"reinsurance_treaty1_premium_proportion"`
	ReinsuranceTreaty2ClaimsProportion         float64 `json:"reinsurance_treaty2_claims_proportion" csv:"reinsurance_treaty2_claims_proportion"`
	ReinsuranceTreaty2PremiumProportion        float64 `json:"reinsurance_treaty2_premium_proportion" csv:"reinsurance_treaty2_premium_proportion"`
	ReinsuranceTreaty3ClaimsProportion         float64 `json:"reinsurance_treaty3_claims_proportion" csv:"reinsurance_treaty3_claims_proportion"`
	ReinsuranceTreaty3PremiumProportion        float64 `json:"reinsurance_treaty3_premium_proportion" csv:"reinsurance_treaty3_premium_proportion"`
	RiskAdjustmentProportion                   float64 `json:"risk_adjustment_proportion" csv:"risk_adjustment_proportion"`
	ReinsuranceTreaty1RiskAdjustmentProportion float64 `json:"reinsurance_treaty1_risk_adjustment_proportion" csv:"reinsurance_treaty1_risk_adjustment_proportion"`
	ReinsuranceTreaty2RiskAdjustmentProportion float64 `json:"reinsurance_treaty2_risk_adjustment_proportion" csv:"reinsurance_treaty2_risk_adjustment_proportion"`
	ReinsuranceTreaty3RiskAdjustmentProportion float64 `json:"reinsurance_treaty3_risk_adjustment_proportion" csv:"reinsurance_treaty3_risk_adjustment_proportion"`
	Created                                    int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
	CreatedBy                                  string  `json:"created_by" csv:"created_by"`
}

type ModifiedGMMScopedAggregation struct {
	ID                                       int     `json:"-" csv:"id" gorm:"primary_key"`
	RunDate                                  string  `json:"run_date" csv:"run_date"`
	RunID                                    int     `json:"-" csv:"run_id" gorm:"index"`
	JobRunId                                 int     `json:"-" csv:"job_run_id"`
	RunName                                  string  `json:"run_name" csv:"run_name"`
	PortfolioName                            string  `json:"portfolio_name" csv:"portfolio_name"`
	PortfolioId                              int     `json:"-" csv:"portfolio_id"`
	ProjectionMonth                          int     `json:"projection_month" csv:"projection_month"`
	ProductCode                              string  `json:"product_code" csv:"product_code"`
	IFRS17Group                              string  `json:"ifrs17_group" csv:"ifrs17_group"`
	Treaty1IFRS17Group                       string  `json:"treaty1_ifrs17_group" csv:"treaty1_ifrs17_group"`
	Treaty2IFRS17Group                       string  `json:"treaty2_ifrs17_group" csv:"treaty2_ifrs17_group"`
	Treaty3IFRS17Group                       string  `json:"treaty3_ifrs17_group" csv:"treaty3_ifrs17_group"`
	LockedinYear                             string  `json:"lockedin_year" csv:"lockedin_year"`
	LockedinMonth                            string  `json:"lockedin_month" csv:"lockedin_month"`
	YieldCurveCode                           string  `json:"yield_curve_code" csv:"yield_curve_code"`
	PremiumReceipt                           float64 `json:"premium_receipt" csv:"premium_receipt"`
	NbPremiumReceipt                         float64 `json:"nb_premium_receipt" csv:"nb_premium_receipt"`
	InitialCommissionOutgo                   float64 `json:"initial_commission_outgo"`
	RenewalCommissionOutgo                   float64 `json:"renewal_commission_outgo"`
	InitialExpenseOutgo                      float64 `json:"initial_expense_outgo"`
	CurrentPeriodEarnedPremium               float64 `json:"current_period_earned_premium"`
	CurrentPeriodAmortisedAcquisition        float64 `json:"current_period_amortised_acquisition"`
	CurrentPeriodInsuranceRevenue            float64 `json:"current_period_insurance_revenue"`
	EarnedPremium                            float64 `json:"earned_premium"`
	UnearnedPremium                          float64 `json:"unearned_premium"`
	SumFutureEarnedPremium                   float64 `json:"sum_future_earned_premium" csv:"sum_future_earned_premium"`
	SumFutureAcquisitionCost                 float64 `json:"sum_future_acquisition_cost" csv:"sum_future_acquisition_cost"`
	DiscountedEarnedPremiumCurrent           float64 `json:"discounted_earned_premium_current" csv:"discounted_earned_premium_current"`
	DiscountedEarnedPremiumLockedin          float64 `json:"discounted_earned_premium_lockedin" csv:"discounted_earned_premium_lockedin"`
	ClaimsOutgo                              float64 `json:"claims_outgo" csv:"claims_outgo"`
	ClaimsExpenseOutgo                       float64 `json:"claims_expense_outgo" csv:"claims_expense_outgo"`
	MaintenanceExpenseOutgo                  float64 `json:"maintenance_expense_outgo" csv:"maintenance_expense_outgo"`
	CashOutflow                              float64 `json:"cash_outflow" csv:"cash_outflow"`
	NetCashFlow                              float64 `json:"net_cash_flow" csv:"net_cash_flow"`
	SumFutureNetCashFlows                    float64 `json:"sum_future_net_cash_flows"`
	DiscountedNetCashFlowsCurrent            float64 `json:"discounted_net_cash_flows_current" csv:"discounted_net_cash_flows_current"`
	DiscountedNetCashFlowsLockedin           float64 `json:"discounted_net_cash_flows_lockedin" csv:"discounted_net_cash_flows_lockedin"`
	DiscountedDACLockedin                    float64 `json:"discounted_dac_lockedin"`
	DiscountedClaimsOutgoLockedin            float64 `json:"discounted_claims_outgo_lockedin"`
	RiskAdjustment                           float64 `json:"risk_adjustment " csv:"risk_adjustment"`
	ModifiedGMMCsm                           float64 `json:"modified_gmm_csm " csv:"modified_gmm_csm"`
	CoverageUnits                            float64 `json:"coverage_units " csv:"coverage_units"`
	SumFutureCashOutflows                    float64 `json:"sum_future_cash_outflows" csv:"sum_future_cash_outflows"`
	DiscountedCashOutflows                   float64 `json:"discounted_cash_outflows" csv:"discounted_cash_outflows"`
	SumFutureCashInflows                     float64 `json:"sum_future_cash_inflows" csv:"sum_future_cash_inflows"`
	DiscountedCashInflows                    float64 `json:"discounted_cash_inflows" csv:"discounted_cash_inflows"`
	OutstandingLoan                          float64 `json:"outstanding_loan" csv:"outstanding_loan"`
	OriginalLoan                             float64 `json:"original_loan" csv:"original_loan"`
	Treaty1CurrentPeriodEarnedPremium        float64 `json:"treaty1_current_period_earned_premium"`
	Treaty1CurrentPeriodAmortisedAcquisition float64 `json:"treaty1_current_period_amortised_acquisition"`
	Treaty2CurrentPeriodEarnedPremium        float64 `json:"treaty2_current_period_earned_premium"`
	Treaty2CurrentPeriodAmortisedAcquisition float64 `json:"treaty2_current_period_amortised_acquisition"`
	Treaty3CurrentPeriodEarnedPremium        float64 `json:"treaty3_current_period_earned_premium"`
	Treaty3CurrentPeriodAmortisedAcquisition float64 `json:"treaty3_current_period_amortised_acquisition"`
	WrittenPremium                           float64 `json:"written_premium"`
	Treaty1WrittenPremium                    float64 `json:"treaty1_written_premium"`
	Treaty2WrittenPremium                    float64 `json:"treaty2_written_premium"`
	Treaty3WrittenPremium                    float64 `json:"treaty3_written_premium"`
	Treaty1DiscountedCashOutflowsLockedin    float64 `json:"treaty_1_discounted_cash_outflows_lockedin" csv:"treaty_1_discounted_cash_outflows_lockedin"`
	Treaty2DiscountedCashOutflowsLockedin    float64 `json:"treaty_2_discounted_cash_outflows_lockedin" csv:"treaty_2_discounted_cash_outflows_lockedin"`
	Treaty3DiscountedCashOutflowsLockedin    float64 `json:"treaty3_discounted_cash_outflows_lockedin" csv:"treaty_3_discounted_cash_outflows_lockedin"`
	IrNbPremiumReceipt                       float64 `json:"ir_nb_premium_receipt"`
	IrInitialCommissionOutgo                 float64 `json:"ir_initial_commission_outgo"`
	IrInitialExpenseOutgo                    float64 `json:"ir_initial_expense_outgo"`
	IrEarnedPremium                          float64 `json:"ir_earned_premium"`
	IrClaimsOutgo                            float64 `json:"ir_claims_outgo"`
	IrMaintenanceExpenseOutgo                float64 `json:"ir_maintenance_expense_outgo"`
	IrNetCashFlow                            float64 `json:"ir_net_cash_flow"`
	IrDiscountedEarnedPremium                float64 `json:"ir_discounted_earned_premium"`
	IrDiscountedClaimsOutgoLockedin          float64 `json:"ir_discounted_claims_outgo_lockedin"`
	IrDiscountedCashOutflowsLockedin         float64 `json:"ir_discounted_cash_outflows_lockedin"`
	IrDiscountedNetCashFlowsLockedin         float64 `json:"ir_discounted_net_cash_flows_lockedin"`
	IrRiskAdjustment                         float64 `json:"ir_risk_adjustment"`
	IrTreaty1DiscountedClaimsOutgo           float64 `json:"ir_treaty1_discounted_claims_outgo"`
	IrTreaty2DiscountedClaimsOutgo           float64 `json:"ir_treaty2_discounted_claims_outgo"`
	IrTreaty3DiscountedClaimsOutgo           float64 `json:"ir_treaty3_discounted_claims_outgo"`
	Created                                  int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type ModifiedGMMProjection struct {
	ID                                       int     `json:"id" csv:"id" gorm:"primary_key"`
	RunDate                                  string  `json:"run_date" csv:"run_date"`
	RunID                                    int     `json:"run_id" csv:"run_id" gorm:"index"`
	RunName                                  string  `json:"run_name" csv:"run_name"`
	PortfolioName                            string  `json:"portfolio_name" csv:"portfolio_name"`
	PortfolioId                              int     `json:"portfolio_id" csv:"portfolio_id"`
	Year                                     int     `json:"year" csv:"year"`
	ProductCode                              string  `json:"product_code" csv:"product_code"`
	PolicyNumber                             string  `json:"policy_number" csv:"policy_number"`
	IFRS17Group                              string  `json:"ifrs17_group" csv:"ifrs17_group"`
	Treaty1IFRS17Group                       string  `json:"treaty1_ifrs17_group" csv:"treaty1_ifrs17_group"`
	Treaty2IFRS17Group                       string  `json:"treaty2_ifrs17_group" csv:"treaty2_ifrs17_group"`
	Treaty3IFRS17Group                       string  `json:"treaty3_ifrs17_group" csv:"treaty3_ifrs17_group"`
	LockedinYear                             string  `json:"lockedin_year" csv:"lockedin_year"`
	LockedinMonth                            string  `json:"lockedin_month" csv:"lockedin_month"`
	ValuationDate                            string  `json:"valuation_date"`
	ProjectionMonth                          int     `json:"projection_month"`
	ValuationMonth                           int     `json:"valuation_month"`
	DaysCovered                              int     `json:"days_covered"`
	LapseRate                                float64 `json:"lapse_rate"`
	InforcePolicyCountSM                     float64 `json:"inforce_policy_count_sm"`
	InforcePolicyCount                       float64 `json:"inforce_policy_count"`
	PremiumReceipt                           float64 `json:"premium_receipt"`
	NbPremiumReceipt                         float64 `json:"nb_premium_receipt"`
	TotalPremiumReceipt                      float64 `json:"total_premium_receipt"`
	DiscountedTotalPremiumReceiptCurrent     float64 `json:"discounted_total_premium_receipt_current"`
	DiscountedTotalPremiumReceiptLockedin    float64 `json:"discounted_total_premium_receipt_lockedin"`
	DiscountedDACLockedin                    float64 `json:"discounted_dac_lockedin"`
	CurrentRate                              float64 `json:"current_rate"`
	LockedInRate                             float64 `json:"locked_in_rate"`
	CurrentPeriodEarnedPremium               float64 `json:"current_period_earned_premium"`
	CurrentPeriodAmortisedAcquisition        float64 `json:"current_period_amortised_acquisition"`
	CurrentPeriodInsuranceRevenue            float64 `json:"current_period_insurance_revenue"`
	EarnedPremium                            float64 `json:"earned_premium"`
	SumFutureEarnedPremium                   float64 `json:"sum_future_earned_premium"`
	SumFutureAcquisitionCost                 float64 `json:"sum_future_acquisition_cost"`
	DiscountedEarnedPremiumCurrent           float64 `json:"discounted_earned_premium_current"`
	DiscountedEarnedPremiumLockedin          float64 `json:"discounted_earned_premium_lockedin"`
	ClaimsOutgo                              float64 `json:"claims_outgo"`
	ClaimsExpenseOutgo                       float64 `json:"claims_expense_outgo"`
	InitialCommissionOutgo                   float64 `json:"initial_commission_outgo"`
	RenewalCommissionOutgo                   float64 `json:"renewal_commission_outgo"`
	InitialExpenseOutgo                      float64 `json:"initial_expense_outgo"`
	MaintenanceExpenseOutgo                  float64 `json:"maintenance_expense_outgo"`
	OutstandingLoan                          float64 `json:"outstanding_loan"`
	OriginalLoan                             float64 `json:"original_loan"`
	CashOutflow                              float64 `json:"cash_outflow"`
	NetCashFlow                              float64 `json:"net_cash_flow"`
	SumFutureCashOutflows                    float64 `json:"sum_future_cash_outflows"`
	DiscountedCashOutflows                   float64 `json:"discounted_cash_outflows"`
	DiscountedCashOutflowsLockedin           float64 `json:"discounted_cash_outflows_lockedin"`
	DiscountedClaimsOutgoLockedin            float64 `json:"discounted_claims_outgo_lockedin"`
	SumFutureNetCashFlows                    float64 `json:"sum_future_net_cash_flows"`
	DiscountedNetCashFlowsCurrent            float64 `json:"discounted_net_cash_flows_current"`
	DiscountedNetCashFlowsLockedin           float64 `json:"discounted_net_cash_flows_lockedin"`
	RiskAdjustment                           float64 `json:"risk_adjustment"`
	ModifiedGMMCsm                           float64 `json:"modified_gmm_csm " csv:"modified_gmm_csm"`
	CoverageUnits                            float64 `json:"coverage_units " csv:"coverage_units"`
	Treaty1PremiumReceipt                    float64 `json:"treaty1_premium_receipt"`
	Treaty1NbPremiumReceipt                  float64 `json:"treaty1_nb_premium_receipt"`
	Treaty1TotalPremiumReceipt               float64 `json:"treaty1_total_premium_receipt"`
	Treaty1CurrentDiscount                   float64 `json:"treaty1_current_discount"`
	Treaty1LockedInDiscount                  float64 `json:"treaty1_locked_in_discount"`
	Treaty1EarnedPremium                     float64 `json:"treaty1_earned_premium"`
	Treaty1SumFutureEarnedPremium            float64 `json:"treaty1_sum_future_earned_premium"`
	Treaty1DiscountedEarnedPremiumCurrent    float64 `json:"treaty1_discounted_earned_premium_current"`
	Treaty1DiscountedEarnedPremiumLockedin   float64 `json:"treaty1_discounted_earned_premium_lockedin"`
	Treaty1ClaimsOutgo                       float64 `json:"treaty1_claims_outgo"`
	Treaty1CashOutflow                       float64 `json:"treaty1_cash_outflow"`
	Treaty1NetCashFlow                       float64 `json:"treaty1_net_cash_flow"`
	Treaty1SumFutureCashOutflows             float64 `json:"treaty1_sum_future_cash_outflows"`
	Treaty1DiscountedCashOutflows            float64 `json:"treaty1_discounted_cash_outflows"`
	Treaty1DiscountedCashOutflowsLockedin    float64 `json:"treaty1_discounted_cash_outflows_lockedin"`
	Treaty1SumFutureNetCashFlows             float64 `json:"treaty1_sum_future_net_cash_flows"`
	Treaty1DiscountedNetCashFlowsCurrent     float64 `json:"treaty1_discounted_net_cash_flows_current"`
	Treaty1DiscountedNetCashFlowsLockedin    float64 `json:"treaty1_discounted_future_net_cash_flows_lockedin"`
	Treaty1RiskAdjustment                    float64 `json:"treaty1_risk_adjustment"`
	Treaty2PremiumReceipt                    float64 `json:"treaty2_premium_receipt"`
	Treaty2NbPremiumReceipt                  float64 `json:"treaty2_nb_premium_receipt"`
	Treaty2TotalPremiumReceipt               float64 `json:"treaty2_total_premium_receipt"`
	Treaty2CurrentDiscount                   float64 `json:"treaty2_current_discount"`
	Treaty2LockedInDiscount                  float64 `json:"treaty2_locked_in_discount"`
	Treaty2EarnedPremium                     float64 `json:"treaty2_earned_premium"`
	Treaty2SumFutureEarnedPremium            float64 `json:"treaty2_sum_future_earned_premium"`
	Treaty2DiscountedEarnedPremiumCurrent    float64 `json:"treaty2_discounted_earned_premium_current"`
	Treaty2DiscountedEarnedPremiumLockedin   float64 `json:"treaty2_discounted_earned_premium_lockedin"`
	Treaty2ClaimsOutgo                       float64 `json:"treaty2_claims_outgo"`
	Treaty2CashOutflow                       float64 `json:"treaty2_cash_outflow"`
	Treaty2NetCashFlow                       float64 `json:"treaty2_net_cash_flow"`
	Treaty2SumFutureCashOutflows             float64 `json:"treaty2_sum_future_cash_outflows"`
	Treaty2DiscountedCashOutflows            float64 `json:"treaty2_discounted_cash_outflows"`
	Treaty2DiscountedCashOutflowsLockedin    float64 `json:"treaty2_discounted_cash_outflows_lockedin"`
	Treaty2SumFutureNetCashFlows             float64 `json:"treaty2_sum_future_net_cash_flows"`
	Treaty2DiscountedNetCashFlowsCurrent     float64 `json:"treaty2_discounted_net_cash_flows_current"`
	Treaty2DiscountedNetCashFlowsLockedin    float64 `json:"treaty2_discounted_net_cash_flows_lockedin"`
	Treaty2RiskAdjustment                    float64 `json:"treaty2_risk_adjustment"`
	Treaty3PremiumReceipt                    float64 `json:"treaty3_premium_receipt"`
	Treaty3NbPremiumReceipt                  float64 `json:"treaty3_nb_premium_receipt"`
	Treaty3TotalPremiumReceipt               float64 `json:"treaty3_total_premium_receipt"`
	Treaty3CurrentDiscount                   float64 `json:"treaty3_current_discount"`
	Treaty3LockedInDiscount                  float64 `json:"treaty3_locked_in_discount"`
	Treaty3EarnedPremium                     float64 `json:"treaty3_earned_premium"`
	Treaty3SumFutureEarnedPremium            float64 `json:"treaty3_sum_future_earned_premium"`
	Treaty3DiscountedEarnedPremiumCurrent    float64 `json:"treaty3_discounted_earned_premium_current"`
	Treaty3DiscountedEarnedPremiumLockedin   float64 `json:"treaty3_discounted_earned_premium_lockedin"`
	Treaty3ClaimsOutgo                       float64 `json:"treaty3_claims_outgo"`
	Treaty3CashOutflow                       float64 `json:"treaty3_cash_outflow"`
	Treaty3NetCashFlow                       float64 `json:"treaty3_net_cash_flow"`
	Treaty3SumFutureCashOutflows             float64 `json:"treaty3_sum_future_cash_outflows"`
	Treaty3DiscountedCashOutflows            float64 `json:"treaty3_discounted_cash_outflows"`
	Treaty3DiscountedCashOutflowsLockedin    float64 `json:"treaty3_discounted_cash_outflows_lockedin"`
	Treaty3SumFutureNetCashFlows             float64 `json:"treaty3_sum_future_net_cash_flows"`
	Treaty3DiscountedNetCashFlowsCurrent     float64 `json:"treaty3_discounted_net_cash_flows_current"`
	Treaty3DiscountedNetCashFlowsLockedin    float64 `json:"treaty3_discounted_net_cash_flows_lockedin"`
	Treaty3RiskAdjustment                    float64 `json:"treaty3_risk_adjustment"`
	Treaty1CurrentPeriodEarnedPremium        float64 `json:"treaty1_current_period_earned_premium"`
	Treaty1CurrentPeriodAmortisedAcquisition float64 `json:"treaty1_current_period_amortised_acquisition"`
	Treaty2CurrentPeriodEarnedPremium        float64 `json:"treaty2_current_period_earned_premium"`
	Treaty2CurrentPeriodAmortisedAcquisition float64 `json:"treaty2_current_period_amortised_acquisition"`
	Treaty3CurrentPeriodEarnedPremium        float64 `json:"treaty3_current_period_earned_premium"`
	Treaty3CurrentPeriodAmortisedAcquisition float64 `json:"treaty3_current_period_amortised_acquisition"`
	WrittenPremium                           float64 `json:"written_premium"`
	Treaty1WrittenPremium                    float64 `json:"treaty1_written_premium"`
	Treaty2WrittenPremium                    float64 `json:"treaty2_written_premium"`
	Treaty3WrittenPremium                    float64 `json:"treaty3_written_premium"`
	CededPremiumReceipt                      float64 `json:"ceded_premium_receipt"`
	CededNbPremiumReceipt                    float64 `json:"ceded_nb_premium_receipt"`
	TotalCededPremiumReceipt                 float64 `json:"total_ceded_premium_receipt"`
	CededCurrentDiscount                     float64 `json:"ceded_current_discount"`
	CededLockedInDiscount                    float64 `json:"ceded_locked_in_discount"`
	CededEarnedPremium                       float64 `json:"ceded_earned_premium"`
	CededSumFutureEarnedPremium              float64 `json:"ceded_sum_future_earned_premium"`
	CededDiscountedEarnedPremiumCurrent      float64 `json:"ceded_discounted_earned_premium_current"`
	CededDiscountedEarnedPremiumLockedin     float64 `json:"ceded_discounted_earned_premium_lockedin"`
	CededClaimsOutgo                         float64 `json:"ceded_claims_outgo"`
	CededCashOutflow                         float64 `json:"ceded_cash_outflow"`
	CededNetCashFlow                         float64 `json:"ceded_net_cash_flow"`
	CededSumFutureCashOutflows               float64 `json:"ceded_sum_future_cash_outflows"`
	CededDiscountedCashOutflows              float64 `json:"ceded_discounted_cash_outflows"`
	CededDiscountedCashOutflowsLockedin      float64 `json:"ceded_discounted_cash_outflows_lockedin"`
	CededSumFutureNetCashFlows               float64 `json:"ceded_sum_future_net_cash_flows"`
	CededDiscountedNetCashFlowsCurrent       float64 `json:"ceded_discounted_net_cash_flows_current"`
	CededDiscountedNetCashFlowsLockedin      float64 `json:"ceded_discounted_net_cash_flows_lockedin"`
	CededRiskAdjustment                      float64 `json:"ceded_risk_adjustment"`
	IrInforcePolicyCountSM                   float64 `json:"ir_inforce_policy_count_sm"`
	IrLapseRate                              float64 `json:"ir_lapse_rate"`
	IrValuationMonth                         int     `json:"ir_valuation_month"`
	IrNbPremiumReceipt                       float64 `json:"ir_nb_premium_receipt"`
	IrInitialCommissionOutgo                 float64 `json:"ir_initial_commission_outgo"`
	IrRenewalCommissionOutgo                 float64 `json:"ir_renewal_commission_outgo"`
	IrInitialExpenseOutgo                    float64 `json:"ir_initial_expense_outgo"`
	IrEarnedPremium                          float64 `json:"ir_earned_premium"`
	IrClaimsOutgo                            float64 `json:"ir_claims_outgo"`
	IrClaimsExpenseOutgo                     float64 `json:"ir_claims_expense_outgo"`
	IrMaintenanceExpenseOutgo                float64 `json:"ir_maintenance_expense_outgo"`
	IrNetCashFlow                            float64 `json:"ir_net_cash_flow"`
	IrDiscountedEarnedPremium                float64 `json:"ir_discounted_earned_premium"`
	IrDiscountedClaimsOutgoLockedin          float64 `json:"ir_discounted_claims_outgo_lockedin"`
	IrDiscountedCashOutflowsLockedin         float64 `json:"ir_discounted_cash_outflows_lockedin"`
	IrDiscountedNetCashFlowsLockedin         float64 `json:"ir_discounted_net_cash_flows_lockedin"`
	IrRiskAdjustment                         float64 `json:"ir_risk_adjustment"`
	IrTreaty1DiscountedClaimsOutgo           float64 `json:"ir_treaty1_discounted_claims_outgo"`
	IrTreaty2DiscountedClaimsOutgo           float64 `json:"ir_treaty2_discounted_claims_outgo"`
	IrTreaty3DiscountedClaimsOutgo           float64 `json:"ir_treaty3_discounted_claims_outgo"`
	Created                                  int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type AggregatedModifiedGMMProjection struct {
	ID                                       int     `json:"-" csv:"id" gorm:"primary_key"`
	RunDate                                  string  `json:"run_date" csv:"run_date"`
	RunID                                    int     `json:"-" csv:"run_id" gorm:"index"`
	JobRunId                                 int     `json:"-" csv:"job_run_id"`
	RunName                                  string  `json:"run_name" csv:"run_name"`
	PortfolioName                            string  `json:"portfolio_name" csv:"portfolio_name"`
	PortfolioId                              int     `json:"-" csv:"portfolio_id"`
	Year                                     int     `json:"year" csv:"year"`
	ProductCode                              string  `json:"product_code" csv:"product_code"`
	PolicyNumber                             string  `json:"policy_number" csv:"policy_number"`
	IFRS17Group                              string  `json:"ifrs17_group" csv:"ifrs17_group"`
	Treaty1IFRS17Group                       string  `json:"treaty1_ifrs17_group" csv:"treaty1_ifrs17_group"`
	Treaty2IFRS17Group                       string  `json:"treaty2_ifrs17_group" csv:"treaty2_ifrs17_group"`
	Treaty3IFRS17Group                       string  `json:"treaty3_ifrs17_group" csv:"treaty3_ifrs17_group"`
	LockedinYear                             string  `json:"lockedin_year" csv:"lockedin_year"`
	LockedinMonth                            string  `json:"lockedin_month" csv:"lockedin_month"`
	YieldCurveCode                           string  `json:"yield_curve_code" csv:"yield_curve_code"`
	ValuationDate                            string  `json:"valuation_date"`
	ProjectionMonth                          int     `json:"projection_month"`
	ValuationMonth                           int     `json:"valuation_month"`
	DaysCovered                              int     `json:"days_covered"`
	LapseRate                                float64 `json:"lapse_rate"`
	InforcePolicyCountSM                     float64 `json:"inforce_policy_count_sm"`
	InforcePolicyCount                       float64 `json:"inforce_policy_count"`
	PremiumReceipt                           float64 `json:"premium_receipt"`
	NbPremiumReceipt                         float64 `json:"nb_premium_receipt"`
	TotalPremiumReceipt                      float64 `json:"total_premium_receipt"`
	DiscountedTotalPremiumReceiptCurrent     float64 `json:"discounted_total_premium_receipt_current"`
	DiscountedTotalPremiumReceiptLockedin    float64 `json:"discounted_total_premium_receipt_lockedin"`
	DiscountedDACLockedin                    float64 `json:"discounted_dac_lockedin"`
	CurrentRate                              float64 `json:"current_rate"`
	LockedInRate                             float64 `json:"locked_in_rate"`
	CurrentPeriodEarnedPremium               float64 `json:"current_period_earned_premium"`
	CurrentPeriodAmortisedAcquisition        float64 `json:"current_period_amortised_acquisition"`
	CurrentPeriodInsuranceRevenue            float64 `json:"current_period_insurance_revenue"`
	EarnedPremium                            float64 `json:"earned_premium"`
	SumFutureEarnedPremium                   float64 `json:"sum_future_earned_premium"`
	SumFutureAcquisitionCost                 float64 `json:"sum_future_acquisition_cost"`
	DiscountedEarnedPremiumCurrent           float64 `json:"discounted_earned_premium_current"`
	DiscountedEarnedPremiumLockedin          float64 `json:"discounted_earned_premium_lockedin"`
	ClaimsOutgo                              float64 `json:"claims_outgo"`
	ClaimsExpenseOutgo                       float64 `json:"claims_expense_outgo"`
	InitialCommissionOutgo                   float64 `json:"initial_commission_outgo"`
	RenewalCommissionOutgo                   float64 `json:"renewal_commission_outgo"`
	InitialExpenseOutgo                      float64 `json:"initial_expense_outgo"`
	MaintenanceExpenseOutgo                  float64 `json:"maintenance_expense_outgo"`
	OutstandingLoan                          float64 `json:"outstanding_loan"`
	OriginalLoan                             float64 `json:"original_loan"`
	CashOutflow                              float64 `json:"cash_outflow"`
	NetCashFlow                              float64 `json:"net_cash_flow"`
	SumFutureCashOutflows                    float64 `json:"sum_future_cash_outflows"`
	DiscountedCashOutflows                   float64 `json:"discounted_cash_outflows"`
	DiscountedCashOutflowsLockedin           float64 `json:"discounted_cash_outflows_lockedin"`
	SumFutureNetCashFlows                    float64 `json:"sum_future_net_cash_flows"`
	DiscountedNetCashFlowsCurrent            float64 `json:"discounted_net_cash_flows_current"`
	DiscountedNetCashFlowsLockedin           float64 `json:"discounted_net_cash_flows_lockedin"`
	RiskAdjustment                           float64 `json:"risk_adjustment"`
	DiscountedClaimsOutgoLockedin            float64 `json:"discounted_claims_outgo_lockedin"`
	CoverageUnits                            float64 `json:"coverage_units " csv:"coverage_units"`
	DiscountedCoverageUnits                  float64 `json:"discounted_coverage_units " csv:"discounted_coverage_units"`
	CSMAllocationRatio                       float64 `json:"csm_allocation_ratio " csv:"csm_allocation_ratio"`
	ModifiedGMMCsm                           float64 `json:"modified_gmm_csm " csv:"modified_gmm_csm"`
	ModifiedGMMCsmRelease                    float64 `json:"modified_gmm_csm_release " csv:"modified_gmm_csm_release"`
	Treaty1PremiumReceipt                    float64 `json:"treaty1_premium_receipt"`
	Treaty1NbPremiumReceipt                  float64 `json:"treaty1_nb_premium_receipt"`
	Treaty1TotalPremiumReceipt               float64 `json:"treaty1_total_premium_receipt"`
	Treaty1CurrentDiscount                   float64 `json:"treaty1_current_discount"`
	Treaty1LockedInDiscount                  float64 `json:"treaty1_locked_in_discount"`
	Treaty1EarnedPremium                     float64 `json:"treaty1_earned_premium"`
	Treaty1SumFutureEarnedPremium            float64 `json:"treaty1_sum_future_earned_premium"`
	Treaty1DiscountedEarnedPremiumCurrent    float64 `json:"treaty1_discounted_earned_premium_current"`
	Treaty1DiscountedEarnedPremiumLockedin   float64 `json:"treaty1_discounted_earned_premium_lockedin"`
	Treaty1ClaimsOutgo                       float64 `json:"treaty1_claims_outgo"`
	Treaty1CashOutflow                       float64 `json:"treaty1_cash_outflow"`
	Treaty1NetCashFlow                       float64 `json:"treaty1_net_cash_flow"`
	Treaty1SumFutureCashOutflows             float64 `json:"treaty1_sum_future_cash_outflows"`
	Treaty1DiscountedCashOutflows            float64 `json:"treaty1_discounted_cash_outflows"`
	Treaty1DiscountedCashOutflowsLockedin    float64 `json:"treaty1_discounted_cash_outflows_lockedin"`
	Treaty1SumFutureNetCashFlows             float64 `json:"treaty1_sum_future_net_cash_flows"`
	Treaty1DiscountedNetCashFlowsCurrent     float64 `json:"treaty1_discounted_net_cash_flows_current"`
	Treaty1DiscountedNetCashFlowsLockedin    float64 `json:"treaty1_discounted_future_net_cash_flows_lockedin"`
	Treaty1RiskAdjustment                    float64 `json:"treaty1_risk_adjustment"`
	Treaty2PremiumReceipt                    float64 `json:"treaty2_premium_receipt"`
	Treaty2NbPremiumReceipt                  float64 `json:"treaty2_nb_premium_receipt"`
	Treaty2TotalPremiumReceipt               float64 `json:"treaty2_total_premium_receipt"`
	Treaty2CurrentDiscount                   float64 `json:"treaty2_current_discount"`
	Treaty2LockedInDiscount                  float64 `json:"treaty2_locked_in_discount"`
	Treaty2EarnedPremium                     float64 `json:"treaty2_earned_premium"`
	Treaty2SumFutureEarnedPremium            float64 `json:"treaty2_sum_future_earned_premium"`
	Treaty2DiscountedEarnedPremiumCurrent    float64 `json:"treaty2_discounted_earned_premium_current"`
	Treaty2DiscountedEarnedPremiumLockedin   float64 `json:"treaty2_discounted_earned_premium_lockedin"`
	Treaty2ClaimsOutgo                       float64 `json:"treaty2_claims_outgo"`
	Treaty2CashOutflow                       float64 `json:"treaty2_cash_outflow"`
	Treaty2NetCashFlow                       float64 `json:"treaty2_net_cash_flow"`
	Treaty2SumFutureCashOutflows             float64 `json:"treaty2_sum_future_cash_outflows"`
	Treaty2DiscountedCashOutflows            float64 `json:"treaty2_discounted_cash_outflows"`
	Treaty2DiscountedCashOutflowsLockedin    float64 `json:"treaty2_discounted_cash_outflows_lockedin"`
	Treaty2SumFutureNetCashFlows             float64 `json:"treaty2_sum_future_net_cash_flows"`
	Treaty2DiscountedNetCashFlowsCurrent     float64 `json:"treaty2_discounted_net_cash_flows_current"`
	Treaty2DiscountedNetCashFlowsLockedin    float64 `json:"treaty2_discounted_net_cash_flows_lockedin"`
	Treaty2RiskAdjustment                    float64 `json:"treaty2_risk_adjustment"`
	Treaty3PremiumReceipt                    float64 `json:"treaty3_premium_receipt"`
	Treaty3NbPremiumReceipt                  float64 `json:"treaty3_nb_premium_receipt"`
	Treaty3TotalPremiumReceipt               float64 `json:"treaty3_total_premium_receipt"`
	Treaty3CurrentDiscount                   float64 `json:"treaty3_current_discount"`
	Treaty3LockedInDiscount                  float64 `json:"treaty3_locked_in_discount"`
	Treaty3EarnedPremium                     float64 `json:"treaty3_earned_premium"`
	Treaty3SumFutureEarnedPremium            float64 `json:"treaty3_sum_future_earned_premium"`
	Treaty3DiscountedEarnedPremiumCurrent    float64 `json:"treaty3_discounted_earned_premium_current"`
	Treaty3DiscountedEarnedPremiumLockedin   float64 `json:"treaty3_discounted_earned_premium_lockedin"`
	Treaty3ClaimsOutgo                       float64 `json:"treaty3_claims_outgo"`
	Treaty3CashOutflow                       float64 `json:"treaty3_cash_outflow"`
	Treaty3NetCashFlow                       float64 `json:"treaty3_net_cash_flow"`
	Treaty3SumFutureCashOutflows             float64 `json:"treaty3_sum_future_cash_outflows"`
	Treaty3DiscountedCashOutflows            float64 `json:"treaty3_discounted_cash_outflows"`
	Treaty3DiscountedCashOutflowsLockedin    float64 `json:"treaty3_discounted_cash_outflows_lockedin"`
	Treaty3SumFutureNetCashFlows             float64 `json:"treaty3_sum_future_net_cash_flows"`
	Treaty3DiscountedNetCashFlowsCurrent     float64 `json:"treaty3_discounted_net_cash_flows_current"`
	Treaty3DiscountedNetCashFlowsLockedin    float64 `json:"treaty3_discounted_net_cash_flows_lockedin"`
	Treaty3RiskAdjustment                    float64 `json:"treaty3_risk_adjustment"`
	Treaty1CurrentPeriodEarnedPremium        float64 `json:"treaty1_current_period_earned_premium"`
	Treaty1CurrentPeriodAmortisedAcquisition float64 `json:"treaty1_current_period_amortised_acquisition"`
	Treaty2CurrentPeriodEarnedPremium        float64 `json:"treaty2_current_period_earned_premium"`
	Treaty2CurrentPeriodAmortisedAcquisition float64 `json:"treaty2_current_period_amortised_acquisition"`
	Treaty3CurrentPeriodEarnedPremium        float64 `json:"treaty3_current_period_earned_premium"`
	Treaty3CurrentPeriodAmortisedAcquisition float64 `json:"treaty3_current_period_amortised_acquisition"`
	WrittenPremium                           float64 `json:"written_premium"`
	Treaty1WrittenPremium                    float64 `json:"treaty1_written_premium"`
	Treaty2WrittenPremium                    float64 `json:"treaty2_written_premium"`
	Treaty3WrittenPremium                    float64 `json:"treaty3_written_premium"`
	CededPremiumReceipt                      float64 `json:"ceded_premium_receipt"`
	CededNbPremiumReceipt                    float64 `json:"ceded_nb_premium_receipt"`
	TotalCededPremiumReceipt                 float64 `json:"total_ceded_premium_receipt"`
	CededCurrentDiscount                     float64 `json:"ceded_current_discount"`
	CededLockedInDiscount                    float64 `json:"ceded_locked_in_discount"`
	CededEarnedPremium                       float64 `json:"ceded_earned_premium"`
	CededSumFutureEarnedPremium              float64 `json:"ceded_sum_future_earned_premium"`
	CededDiscountedEarnedPremiumCurrent      float64 `json:"ceded_discounted_earned_premium_current"`
	CededDiscountedEarnedPremiumLockedin     float64 `json:"ceded_discounted_earned_premium_lockedin"`
	CededClaimsOutgo                         float64 `json:"ceded_claims_outgo"`
	CededCashOutflow                         float64 `json:"ceded_cash_outflow"`
	CededNetCashFlow                         float64 `json:"ceded_net_cash_flow"`
	CededSumFutureCashOutflows               float64 `json:"ceded_sum_future_cash_outflows"`
	CededDiscountedCashOutflows              float64 `json:"ceded_discounted_cash_outflows"`
	CededDiscountedCashOutflowsLockedin      float64 `json:"ceded_discounted_cash_outflows_lockedin"`
	CededSumFutureNetCashFlows               float64 `json:"ceded_sum_future_net_cash_flows"`
	CededDiscountedNetCashFlowsCurrent       float64 `json:"ceded_discounted_net_cash_flows_current"`
	CededDiscountedNetCashFlowsLockedin      float64 `json:"ceded_discounted_net_cash_flows_lockedin"`
	CededRiskAdjustment                      float64 `json:"ceded_risk_adjustment"`
	IrInforcePolicyCountSM                   float64 `json:"ir_inforce_policy_count_sm"`
	IrLapseRate                              float64 `json:"ir_lapse_rate"`
	IrValuationMonth                         int     `json:"ir_valuation_month"`
	IrNbPremiumReceipt                       float64 `json:"ir_nb_premium_receipt"`
	IrInitialCommissionOutgo                 float64 `json:"ir_initial_commission_outgo"`
	IrRenewalCommissionOutgo                 float64 `json:"ir_renewal_commission_outgo"`
	IrInitialExpenseOutgo                    float64 `json:"ir_initial_expense_outgo"`
	IrEarnedPremium                          float64 `json:"ir_earned_premium"`
	IrClaimsOutgo                            float64 `json:"ir_claims_outgo"`
	IrClaimsExpenseOutgo                     float64 `json:"ir_claims_expense_outgo"`
	IrMaintenanceExpenseOutgo                float64 `json:"ir_maintenance_expense_outgo"`
	IrNetCashFlow                            float64 `json:"ir_net_cash_flow"`
	IrDiscountedEarnedPremium                float64 `json:"ir_discounted_earned_premium"`
	IrDiscountedClaimsOutgoLockedin          float64 `json:"ir_discounted_claims_outgo_lockedin"`
	IrDiscountedCashOutflowsLockedin         float64 `json:"ir_discounted_cash_outflows_lockedin"`
	IrDiscountedNetCashFlowsLockedin         float64 `json:"ir_discounted_net_cash_flows_lockedin"`
	IrRiskAdjustment                         float64 `json:"ir_risk_adjustment"`
	IrTreaty1EarnedPremium                   float64 `json:"ir_treaty1_earned_premium"`
	IrTreaty2EarnedPremium                   float64 `json:"ir_treaty2_earned_premium"`
	IrTreaty3EarnedPremium                   float64 `json:"ir_treaty3_earned_premium"`
	IrTreaty1ClaimsOutgo                     float64 `json:"ir_treaty1_claims_outgo"`
	IrTreaty2ClaimsOutgo                     float64 `json:"ir_treaty2_claims_outgo"`
	IrTreaty3ClaimsOutgo                     float64 `json:"ir_treaty3_claims_outgo"`
	IrTreaty1DiscountedClaimsOutgo           float64 `json:"ir_treaty1_discounted_claims_outgo"`
	IrTreaty2DiscountedClaimsOutgo           float64 `json:"ir_treaty2_discounted_claims_outgo"`
	IrTreaty3DiscountedClaimsOutgo           float64 `json:"ir_treaty3_discounted_claims_outgo"`
	Created                                  int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type ModifiedDiscountedValues struct {
	TotalPremiumReceipt                    float64
	DiscountedTotalPremiumReceiptCurrent   float64
	DiscountedTotalPremiumReceiptLockedin  float64
	EarnedPremium                          float64
	SumFutureEarnedPremium                 float64
	DiscountedEarnedPremiumCurrent         float64
	DiscountedEarnedPremiumLockedin        float64
	CashOutflow                            float64
	ClaimsOutgo                            float64
	InitialCommissionOutgo                 float64
	InitialExpenseOutgo                    float64
	ClaimsExpenseOutgo                     float64
	SumFutureCashOutflows                  float64
	DiscountedCashOutflows                 float64
	DiscountedCashOutflowsLockedin         float64
	NetCashFlow                            float64
	SumFutureNetCashFlows                  float64
	SumFutureAcquisitionCost               float64
	DiscountedNetCashFlowsCurrent          float64
	DiscountedNetCashFlowsLockedin         float64
	DiscountedDACLockedin                  float64
	RiskAdjustment                         float64
	DiscountedClaimsOutgoLockedin          float64
	CoverageUnits                          float64
	DiscountedCoverageUnits                float64
	Treaty1PremiumReceipt                  float64
	Treaty1TotalPremiumReceipt             float64
	Treaty1EarnedPremium                   float64
	Treaty1SumFutureEarnedPremium          float64
	Treaty1DiscountedEarnedPremiumCurrent  float64
	Treaty1DiscountedEarnedPremiumLockedin float64
	Treaty1CashOutflow                     float64
	Treaty1SumFutureCashOutflows           float64
	Treaty1DiscountedCashOutflows          float64
	Treaty1DiscountedCashOutflowsLockedin  float64
	Treaty1NetCashFlow                     float64
	Treaty1SumFutureNetCashFlows           float64
	Treaty1DiscountedNetCashFlowsCurrent   float64
	Treaty1DiscountedNetCashFlowsLockedin  float64
	Treaty1RiskAdjustment                  float64
	Treaty2PremiumReceipt                  float64
	Treaty2TotalPremiumReceipt             float64
	Treaty2EarnedPremium                   float64
	Treaty2CashOutflow                     float64
	Treaty2NetCashFlow                     float64
	Treaty2SumFutureEarnedPremium          float64
	Treaty2DiscountedEarnedPremiumCurrent  float64
	Treaty2DiscountedEarnedPremiumLockedin float64
	Treaty2SumFutureCashOutflows           float64
	Treaty2DiscountedCashOutflows          float64
	Treaty2DiscountedCashOutflowsLockedin  float64
	Treaty2SumFutureNetCashFlows           float64
	Treaty2DiscountedNetCashFlowsCurrent   float64
	Treaty2DiscountedNetCashFlowsLockedin  float64
	Treaty2RiskAdjustment                  float64
	Treaty3PremiumReceipt                  float64
	Treaty3TotalPremiumReceipt             float64
	Treaty3EarnedPremium                   float64
	Treaty3CashOutflow                     float64
	Treaty3NetCashFlow                     float64
	Treaty3SumFutureEarnedPremium          float64
	Treaty3DiscountedEarnedPremiumCurrent  float64
	Treaty3DiscountedEarnedPremiumLockedin float64
	Treaty3SumFutureCashOutflows           float64
	Treaty3DiscountedCashOutflows          float64
	Treaty3DiscountedCashOutflowsLockedin  float64
	Treaty3SumFutureNetCashFlows           float64
	Treaty3DiscountedNetCashFlowsCurrent   float64
	Treaty3DiscountedNetCashFlowsLockedin  float64
	Treaty3RiskAdjustment                  float64
	TotalCededPremiumReceipt               float64
	CededEarnedPremium                     float64
	CededCashOutflow                       float64
	CededNetCashFlow                       float64
	CededSumFutureEarnedPremium            float64
	CededDiscountedEarnedPremiumCurrent    float64
	CededDiscountedEarnedPremiumLockedin   float64
	CededSumFutureCashOutflows             float64
	CededDiscountedCashOutflows            float64
	CededDiscountedCashOutflowsLockedin    float64
	CededSumFutureNetCashFlows             float64
	CededDiscountedNetCashFlowsCurrent     float64
	CededDiscountedNetCashFlowsLockedin    float64
	CededRiskAdjustment                    float64
	IrNbPremiumReceipt                     float64
	IrInitialCommissionOutgo               float64
	IrRenewalCommissionOutgo               float64
	IrInitialExpenseOutgo                  float64
	IrEarnedPremium                        float64
	IrClaimsOutgo                          float64
	IrClaimsExpenseOutgo                   float64
	IrMaintenanceExpenseOutgo              float64
	IrDiscountedEarnedPremium              float64
	IrDiscountedNetCashFlowsLockedin       float64
	IrRiskAdjustment                       float64
	IrDiscountedClaimsOutgoLockedin        float64
	IrDiscountedCashOutflowsLockedin       float64
	IrTreaty1ClaimsOutgo                   float64
	IrTreaty2ClaimsOutgo                   float64
	IrTreaty3ClaimsOutgo                   float64
	IrTreaty1DiscountedClaimsOutgo         float64
	IrTreaty2DiscountedClaimsOutgo         float64
	IrTreaty3DiscountedClaimsOutgo         float64
}

type MgmmRun struct {
	ID               int             `json:"id" gorm:"primary_key"`
	Name             string          `json:"name" csv:"name"`
	RunDate          string          `json:"run_date" csv:"run_date"`
	RunType          string          `json:"run_type" csv:"run_type"`
	UserName         string          `json:"user_name" csv:"user_name"`
	UserEmail        string          `json:"user_email" csv:"user_email"`
	GMMRunSettings   []GMMRunSetting `json:"gmm_run_settings" csv:"gmm_run_settings"`
	Ifrs17Ready      bool            `json:"ifrs17_ready" csv:"ifrs17_ready"`
	TotalRecords     int             `json:"total_records" csv:"total_records"`
	ProcessedRecords int             `json:"processed_records" csv:"processed_records"`
	ProcessingStatus string          `json:"processing_status" csv:"processing_status"`
	FailureReason    string          `json:"failure_reason" csv:"failure_reason"`
	RunTime          float64         `json:"run_time" csv:"run_time"`
	CreationDate     time.Time       `json:"creation_date" csv:"creation_date"`
}

type MgmmRunPayload struct {
	ID                  int                     `json:"id" gorm:"primary_key"`
	Name                string                  `json:"name" csv:"name"`
	RunDate             string                  `json:"run_date" csv:"run_date"`
	RunType             string                  `json:"run_type" csv:"run_type"`
	UserName            string                  `json:"user_name" csv:"user_name"`
	UserEmail           string                  `json:"user_email" csv:"user_email"`
	ModelPoint          int                     `json:"model_point" csv:"model_point"`
	ModelPointVersion   string                  `json:"model_point_version" csv:"model_point_version"`
	ExpectedClaims      int                     `json:"expected_claims" csv:"expected_claims"`
	ClaimsExpenses      int                     `json:"claims_expenses" csv:"claims_expenses"`
	MaintenanceExpenses int                     `json:"maintenance_expense" csv:"maintenance_expense"`
	AcquisitionExpenses int                     `json:"acquisition_expense" csv:"acquisition_expense"`
	RiskAdjustment      int                     `json:"risk_adjustment" csv:"risk_adjustment"`
	Reinsurance         int                     `json:"reinsurance" csv:"reinsurance"`
	PremiumEarning      int                     `json:"premium_earning" csv:"premium_earning"`
	ShockSettingID      int                     `json:"shock_setting_id" csv:"shock_setting_id"`
	ShockSettings       ModifiedGMMShockSetting `json:"shock_setting" csv:"shock_setting" gorm:"-"`
	IFRS17Aggregation   bool                    `json:"ifrs17_aggregation" csv:"ifrs17_aggregation"`
	YearEndMonth        int                     `json:"year_end_month" csv:"year_end_month"`
	ProjectionPeriod    int                     `json:"projection_period" csv:"projection_period"`
	TotalRecords        int                     `json:"total_records" csv:"total_records"`
	ProcessedRecords    int                     `json:"processed_records" csv:"processed_records"`
	ProcessingStatus    string                  `json:"processing_status" csv:"processing_status"`
	RunTime             float64                 `json:"run_time" csv:"run_time"`
	RunSingle           bool                    `json:"run_single" csv:"run_single"`
	IndividualResults   bool                    `json:"individual_results" csv:"individual_results"`
	Portfolios          []GMMRunSetting         `json:"portfolios" csv:"portfolios"`
	CreationDate        time.Time               `json:"creation_date"`
}

type GMMRunSetting struct {
	ID                int                     `json:"id" gorm:"primary_key"`
	MgmmRunID         int                     `json:"mgmm_run_id"`
	Name              string                  `json:"name" csv:"name"`
	RunDate           string                  `json:"run_date" csv:"run_date"`
	RunType           string                  `json:"run_type" csv:"run_type"`
	UserName          string                  `json:"user_name" csv:"user_name"`
	UserEmail         string                  `json:"user_email" csv:"user_email"`
	PortfolioName     string                  `json:"portfolio_name" csv:"portfolio_name"`
	PortfolioId       int                     `json:"portfolio_id"`
	Description       string                  `json:"description" csv:"description"`
	ModelPoint        int                     `json:"model_point" csv:"model_point"`
	ModelPointVersion string                  `json:"model_point_version" csv:"model_point_version"`
	PremiumEarning    int                     `json:"premium_earning" csv:"premium_earning"`
	ParameterYear     int                     `json:"parameter_year" csv:"parameter_year"`
	ShockSettingID    int                     `json:"shock_setting_id" csv:"shock_setting_id"`
	ShockSettings     ModifiedGMMShockSetting `json:"shock_setting" csv:"shock_setting" gorm:"-"`
	IFRS17Aggregation bool                    `json:"ifrs17_aggregation" csv:"ifrs17_aggregation"`
	YearEndMonth      int                     `json:"year_end_month" csv:"year_end_month"`
	ProjectionPeriod  int                     `json:"projection_period" csv:"projection_period"`
	TotalRecords      int                     `json:"total_records" csv:"total_records"`
	ProcessedRecords  int                     `json:"processed_records" csv:"processed_records"`
	ProcessingStatus  string                  `json:"processing_status" csv:"processing_status"`
	FailureReason     string                  `json:"failure_reason" csv:"failure_reason"`
	RunTime           float64                 `json:"run_time" csv:"run_time"`
	RunSingle         bool                    `json:"run_single" csv:"run_single"`
	IndividualResults bool                    `json:"individual_results" csv:"individual_results"`
	CreationDate      time.Time               `json:"creation_date"`
}

type GMMRunSettingDisplay struct {
	ID                  int                     `json:"id" gorm:"primary_key"`
	Name                string                  `json:"name" csv:"name"`
	RunDate             string                  `json:"run_date" csv:"run_date"`
	RunType             string                  `json:"run_type" csv:"run_type"`
	PortfolioName       string                  `json:"portfolio_name" csv:"portfolio_name"`
	PortfolioId         int                     `json:"portfolio_id"`
	Description         string                  `json:"description" csv:"description"`
	ModelPoint          int                     `json:"model_point" csv:"model_point"`
	DiscountCurve       int                     `json:"discount_curve" csv:"discount_curve"`
	ExpectedClaims      int                     `json:"expected_claims" csv:"expected_claims"`
	ClaimsExpenses      int                     `json:"claims_expenses" csv:"claims_expenses"`
	MaintenanceExpenses int                     `json:"maintenance_expense" csv:"maintenance_expense"`
	AcquisitionExpenses int                     `json:"acquisition_expense" csv:"acquisition_expense"`
	RiskAdjustment      int                     `json:"risk_adjustment" csv:"risk_adjustment"`
	Reinsurance         int                     `json:"reinsurance" csv:"reinsurance"`
	PremiumEarning      int                     `json:"premium_earning" csv:"premium_earning"`
	ShockSettingID      int                     `json:"shock_setting_id" csv:"shock_setting_id"`
	ShockSettings       ModifiedGMMShockSetting `json:"shock_setting" csv:"shock_setting" gorm:"-"`
	IFRS17Aggregation   bool                    `json:"ifrs17_aggregation" csv:"ifrs17_aggregation"`
	YearEndMonth        int                     `json:"year_end_month" csv:"year_end_month"`
	TotalRecords        int                     `json:"total_records" csv:"total_records"`
	ProcessedRecords    int                     `json:"processed_records" csv:"processed_records"`
	ProcessingStatus    string                  `json:"processing_status" csv:"processing_status"`
	RunTime             float64                 `json:"run_time" csv:"run_time"`
	RunSingle           bool                    `json:"run_single" csv:"run_single"`
	IndividualResults   bool                    `json:"individual_results" csv:"individual_results"`
	CreationDate        time.Time               `json:"creation_date"`
}

type ModifiedGMMShockSetting struct {
	ID            int    `json:"id" gorm:"primary_key"`
	Name          string `json:"name" gorm:"unique"`
	Description   string `json:"description"`
	DiscountCurve bool   `json:"discount_curve"`
	Claims        bool   `json:"claims" csv:"claims"`
	Expenses      bool   `json:"expenses" csv:"expenses"`
	ShockBasis    string `json:"shock_basis" csv:"shock_basis"`
	Year          int    `json:"year" csv:"year"`
}

type ModifiedGMMShock struct {
	ID                          int     `json:"id" gorm:"primary_key"`
	ShockBasis                  string  `json:"shock_basis" csv:"shock_basis"`
	ProjectionMonth             int     `json:"projection_month" csv:"projection_month"`
	MultiplicativeClaims        float64 `json:"multiplicative_claims" csv:"multiplicative_claims"`
	AdditiveClaims              float64 `json:"additive_claims" csv:"additive_claims"`
	MultiplicativeDiscountCurve float64 `json:"multiplicative_discount_curve" csv:"multiplicative_discount_curve"`
	AdditiveDiscountCurve       float64 `json:"additive_discount_curve" csv:"additive_discount_curve"`
	MultiplicativeExpenses      float64 `json:"multiplicative_expenses" csv:"multiplicative_expenses"`
	AdditiveExpenses            float64 `json:"additive_expenses" csv:"additive_expenses"`
	Created                     int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
	CreatedBy                   string  `json:"created_by" csv:"created_by"`
}

type TableMetaData struct {
	TableType   string                   `json:"table_type"`
	TableName   string                   `json:"table_name"`
	Category    string                   `json:"category"`
	Data        []map[string]interface{} `json:"data"`
	Populated   bool                     `json:"populated"`
	TableKey    string                   `json:"table_key,omitempty"`    // canonical statName, e.g. "regionLoading"
	DeleteKey   string                   `json:"delete_key,omitempty"`   // slug expected by rate-tables/:table_type endpoints, e.g. "reinsurancegeneralloadings"
	IsRequired  bool                     `json:"is_required"`            // false = skip DB read, downstream reads zero
	UpdatedBy   string                   `json:"updated_by,omitempty"`   // last user to toggle IsRequired
	UpdatedAt   *time.Time               `json:"updated_at,omitempty"`
}

type PaaPortfolio struct {
	ID                    int                     `json:"id" csv:"id" gorm:"primary_key"`
	Name                  string                  `json:"name" csv:"name" gorm:"unique"`
	DiscountOption        string                  `json:"discount_option" csv:"discount_option"`
	PremiumEarningPattern string                  `json:"premium_earning_pattern" csv:"premium_earning_pattern"`
	InsuranceType         string                  `json:"insurance_type" csv:"insurance_type"`
	ModelPoints           []ModifiedGMMModelPoint `json:"model_points" csv:"model_points"`
	YearVersions          []PAAYearVersion        `json:"year_versions" csv:"year_versions" gorm:"-"`
	ModelPointYears       []int                   `json:"model_point_years" csv:"model_point_years" gorm:"-"`
}

type PAAYearVersion struct {
	ID            int    `json:"id" csv:"id" gorm:"primary_key"`
	PortfolioName string `json:"portfolio_name" csv:"portfolio_name"`
	PortfolioId   int    `json:"portfolio_id"`
	Year          int    `json:"year" csv:"year"`
	MpVersion     string `json:"mp_version" csv:"mp_version"`
	Count         int    `json:"count" csv:"count"`
}

type PremiumEarningPattern struct {
	ID                        int     `json:"id" csv:"id" gorm:"primary_key"`
	Year                      int     `json:"year" csv:"year"`
	DurationInForce           int     `json:"duration_in_force" csv:"duration_in_force"`
	PremiumEarningPatternCode string  `json:"premium_earning_pattern_code" csv:"premium_earning_pattern_code"`
	RiskUnit                  float64 `json:"risk_unit" csv:"risk_unit"`
	Created                   int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
	CreatedBy                 string  `json:"created_by" csv:"created_by"`
}

type PAALapse struct {
	ID                   int     `json:"id" csv:"id" gorm:"primary_key"`
	Year                 int     `json:"year" csv:"year"`
	ProductCode          string  `json:"product_code" csv:"product_code"`
	DurationInForceMonth int     `json:"duration_in_force_month" csv:"duration_in_force_month"`
	DistributionChannel  string  `json:"distribution_channel" csv:"distribution_channel"`
	LapseRate            float64 `json:"lapse_rate" csv:"lapse_rate"`
	Created              int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
	CreatedBy            string  `json:"created_by" csv:"created_by"`
}

type PAAFinance struct {
	ID                      int     `json:"id" gorm:"primary_key"`
	Year                    int     `json:"year"`
	Version                 string  `json:"version"`
	PortfolioName           string  `json:"portfolio_name" csv:"portfolio_name"`
	ProductCode             string  `json:"product_code" csv:"product_code"`
	IFRS17Group             string  `json:"ifrs17_group" csv:"ifrs17_group"`
	IFStatus                string  `json:"if_status" csv:"if_status"`
	ExitStatus              string  `json:"exit_status" csv:"exit_status"`
	ActualPremiumReceipt    float64 `json:"actual_premium_receipt" csv:"actual_premium_receipt"`
	PremiumDebtors          float64 `json:"premium_debtors" csv:"premium_debtors"`
	PremiumRefund           float64 `json:"premium_refund" csv:"premium_refund"`
	AcquisitionCostPaid     float64 `json:"acquisition_cost_paid" csv:"acquisition_cost_paid"`
	CommissionClawback      float64 `json:"commission_clawback" csv:"commission_clawback"`
	AttributableExpenses    float64 `json:"attributable_expenses" csv:"attributable_expenses"`
	NonAttributableExpenses float64 `json:"non_attributable_expenses" csv:"non_attributable_expenses"`
	ClaimsExpenses          float64 `json:"claims_expenses" csv:"claims_expenses"`
	ReinsurancePremium      float64 `json:"reinsurance_premium" csv:"reinsurance_premium"`
	ReinstatementPremium    float64 `json:"reinstatement_premium" csv:"reinstatement_premium"`
	ReinsuranceRecovery     float64 `json:"reinsurance_recovery" csv:"reinsurance_recovery"`
	Treaty1Recovery         float64 `json:"treaty_1_recovery" csv:"treaty_1_recovery"`
	Treaty2Recovery         float64 `json:"treaty_2_recovery" csv:"treaty_2_recovery"`
	Treaty3Recovery         float64 `json:"treaty_3_recovery" csv:"treaty_3_recovery"`
	ProvisionalCommission   float64 `json:"provisional_commission" csv:"provisional_commission"`
	IncurredClaims          float64 `json:"incurred_claims" csv:"incurred_claims"`
	ClaimsPaid              float64 `json:"claims_paid" csv:"claims_paid"`
	DacBuildupIndicator     bool    `json:"dac_buildup_indicator" csv:"dac_buildup_indicator"`
	Created                 int64   `json:"created" gorm:"autoCreateTime"`
	CreatedBy               string  `json:"created_by" `
}

type ReinsuranceParameter struct {
	ID                        int     `json:"id" csv:"id" gorm:"primary_key"`
	Year                      int     `json:"year" csv:"year"`
	PortfolioName             string  `json:"portfolio_name" csv:"portfolio_name"`
	ProductCode               string  `json:"product_code" csv:"product_code"`
	ReinsuranceInwardOutward  string  `json:"reinsurance_inward_outward" csv:"reinsurance_inward_outward"`
	UlrLowerboundRate         float64 `json:"ulr_lowerbound_rate" csv:"ulr_lowerbound_rate"`
	UlrUpperboundRate         float64 `json:"ulr_upperbound_rate" csv:"ulr_upperbound_rate"`
	SlidingScaleMinRate       float64 `json:"sliding_scale_min_rate" csv:"sliding_scale_min_rate"`
	SlidingScaleMaxRate       float64 `json:"sliding_scale_max_rate" csv:"sliding_scale_max_rate"`
	ProvisionalCommissionRate float64 `json:"provisional_commission_rate" csv:"provisional_commission_rate"`
	ProfitCommissionRate      float64 `json:"profit_commission_rate" csv:"profit_commission_rate"`
	ProfitCommissionModel     string  `json:"profit_commission_model" csv:"profit_commission_model"`
	Created                   int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
	CreatedBy                 string  `json:"created_by" csv:"created_by"`
}

type PaaYieldCurve struct {
	ID   int `json:"id" csv:"id" gorm:"primary_key"`
	Year int `json:"year" csv:"year" gorm:"primary_key;auto_increment:false"`
	//ExpConfigurationName  string  `json:"portfolio_name" csv:"portfolio_name"`
	YieldCurveCode string `json:"yield_curve_code" csv:"yield_curve_code"`
	//ExpConfigurationId    int     `json:"portfolio_id,omitempty" csv:"portfolio_id"`
	ProjectionTime int     `json:"proj_time" csv:"proj_time" gorm:"primary_key;auto_increment:false;column:proj_time"`
	Month          int     `json:"month" csv:"month" gorm:"primary_key;auto_increment:false;column:month"`
	NominalRate    float64 `json:"nominal_rate" csv:"nominal_rate"`
	Inflation      float64 `json:"inflation" csv:"inflation"`
	Created        int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
	CreatedBy      string  `json:"created_by" csv:"created_by"`
}

func (PaaYieldCurve) TableName() string {
	return "paa_yield_curve"
}
