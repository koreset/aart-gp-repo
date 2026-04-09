package models

type ProductParameters struct {
	Year           int    `json:"year" gorm:"primary_key;auto_increment:false" csv:"year"`
	ProductCode    string `json:"product_code" gorm:"primary_key;auto_increment:false;size:191"  csv:"product_code"`
	YieldCurveCode string `json:"yield_curve_code" csv:"yield_curve_code"`
	Basis          string `json:"basis" gorm:"primary_key;auto_increment:false;size:191" csv:"basis"`
	MarginBasis    string `json:"margin_basis" csv:"margin_basis"`
	//InitialCommissionPercentage1         float64 `json:"initial_commission_percentage1" csv:"initial_commission_percentage1"`
	//InitialCommissionPercentage2         float64 `json:"initial_commission_percentage2" csv:"initial_commission_percentage2"`
	//InitialCommissionRand                float64 `json:"initial_commission_rand" csv:"initial_commission_rand"`
	//ClawbackPeriod                       int     `json:"clawback_period" csv:"clawback_period"`
	//RenewalCommissionPercentage          float64 `json:"renewal_commission_percentage" csv:"renewal_commission_percentage"`
	//RenewalCommissionRand                float64 `json:"renewal_commission_rand" csv:"renewal_commission_rand"`
	InitialExpensePercentage             float64 `json:"initial_expense_percentage" csv:"initial_expense_percentage"`
	InitialExpenseRand                   float64 `json:"initial_expense_rand" csv:"initial_expense_rand"`
	RenewalExpensePercentage             float64 `json:"renewal_expense_percentage" csv:"renewal_expense_percentage"`
	RenewalExpenseRand                   float64 `json:"renewal_expense_rand" csv:"renewal_expense_rand"`
	InitialAttributableExpenseProportion float64 `json:"initial_attributable_expense_proportion" csv:"initial_attributable_expense_proportion"`
	RenewalAttributableExpenseProportion float64 `json:"renewal_attributable_expense_proportion" csv:"renewal_attributable_expense_proportion"`
	ClaimsAttributableExpenseProportion  float64 `json:"claims_attributable_expense_proportion" csv:"claims_attributable_expense_proportion"`
	PremiumWaiverWaitingPeriod           int     `json:"premium_waiver_waiting_period" csv:"premium_waiver_waiting_period"`
	EducatorSumAssuredPaymentAge         int     `json:"educator_sum_assured_payment_age" csv:"educator_sum_assured_payment_age"`
	PaidUpOnSurvivalWaitingPeriod        int     `json:"paid_up_on_survival_waiting_period" csv:"paid_up_on_survival_waiting_period"`
	ChildExitAge                         int     `json:"child_exit_age" csv:"child_exit_age"`
	Timing                               float64 `json:"timing" csv:"timing"`
	PremiumWaiverSumAssuredFactor        float64 `json:"premium_waiver_sum_assured_factor" csv:"premium_waiver_sum_assured_factor"`
	PremiumWaiverAdjustedSumassuredTerm  float64 `json:"premium_waiver_adjusted_sumassured_term" csv:"premium_waiver_adjusted_sumassured_term"`
	PremiumWaiverExitAge                 float64 `json:"premium_waiver_exit_age" csv:"premium_waiver_exit_age"`
	CashbackOnSurvivalPeriod             float64 `json:"cashback_on_survival_period" csv:"cashback_on_survival_period"`
	CashbackOnSurvivalTerm               float64 `json:"cashback_on_survival_term" csv:"cashback_on_survival_term"`
	CashbackOnDeathTerm                  float64 `json:"cashback_on_death_term" csv:"cashback_on_death_term"`
	PaidupEffectiveAge                   int     `json:"paidup_effective_age" csv:"paidup_effective_age"`
	PaidupEffectiveDuration              float64 `json:"paidup_effective_duration" csv:"paidup_effective_duration"`
	CashbackOnSurvivalRatio              float64 `json:"cashback_on_survival_ratio" csv:"cashback_on_survival_ratio"`
	CashbackOnDeathRatio                 float64 `json:"cashback_on_death_ratio" csv:"cashback_on_death_ratio"`
	ClaimsExpensePercentage              float64 `json:"claims_expense_percentage" csv:"claims_expense_percentage"`
	ClaimsExpenseRand                    float64 `json:"claims_expense_rand" csv:"claims_expense_rand"`
	StandardAdditionalSumAssured         float64 `json:"standard_additional_sum_assured" csv:"standard_additional_sum_assured"`
	RetrenchmentQualificationFactor      float64 `json:"retrenchment_qualification_factor" csv:"retrenchment_qualification_factor"`
	RepatriationAssumption               float64 `json:"repatriation_assumption" csv:"repatriation_assumption"`
	MainMemberIndicator                  bool    `json:"main_member_indicator" csv:"main_member_indicator" gorm:"type:bool"`
	OtherLivesIndicator                  bool    `json:"other_lives_indicator" csv:"other_lives_indicator"`
	NumberRetrenchmentInstalments        float64 `json:"number_retrenchment_instalments" csv:"number_retrenchment_instalments"`
	RiderFuneralSaProportion             float64 `json:"rider_funeral_sa_proportion" csv:"rider_funeral_sa_proportion"`
	RiderFuneralMinimumBenefit           float64 `json:"rider_funeral_minimum_benefit" csv:"rider_funeral_minimum_benefit"`
	RiderFuneralMaximumBenefit           float64 `json:"rider_funeral_maximum_benefit" csv:"rider_funeral_maximum_benefit"`
	PremiumHolidayWaitingPeriod          int     `json:"premium_holiday_waiting_period" csv:"premium_holiday_waiting_period"`
	MaximumPremiumHolidays               int     `json:"maximum_premium_holidays" csv:"maximum_premium_holidays"`
	PremiumHolidayCycle                  int     `json:"premium_holiday_cycle" csv:"premium_holiday_cycle"`
	IbnrProportion                       float64 `json:"ibnr_proportion" csv:"ibnr_proportion"`
	//HybridRenewalCommStartM              int     `json:"hybrid_renewal_comm_start_m" csv:"hybrid_renewal_comm_start_m"`
	//HybridRenewalCommEndM                int     `json:"hybrid_renewal_comm_end_m" csv:"hybrid_renewal_comm_end_m"`
	LossRatio                   float64 `json:"loss_ratio" csv:"loss_ratio"`
	CorporateTaxRate            float64 `json:"corporate_tax_rate" csv:"corporate_tax_rate"`
	ShareholdersRequiredMargin  float64 `json:"shareholders_required_margin" csv:"shareholders_required_margin"`
	SpecialMarginCode           string  `json:"special_margin_code" csv:"special_margin_code"`
	ProfitAdjustmentCode        string  `json:"profit_adjustment_code" csv:"profit_adjustment_code"`
	Created                     int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
	CreatedBy                   string  `json:"created_by" csv:"created_by"`
	CalculatedTerm              int     `json:"-" csv:"-"`
	AccidentalBenefitMultiplier float64 `json:"-" csv:"-"`
	BenefitMultiplier           float64 `json:"-" csv:"-"`
	FlatAnnualReinsPremRate     float64 `json:"-" csv:"-"`
	Level1CededProp             float64 `json:"-" csv:"-"`
	Level1Lowerbound            float64 `json:"-" csv:"-"`
	Level1Upperbound            float64 `json:"-" csv:"-"`
	Level2CededProp             float64 `json:"-" csv:"-"`
	Level2Lowerbound            float64 `json:"-" csv:"-"`
	Level2Upperbound            float64 `json:"-" csv:"-"`
	Level3CededProp             float64 `json:"-" csv:"-"`
	Level3Lowerbound            float64 `json:"-" csv:"-"`
	Level3Upperbound            float64 `json:"-" csv:"-"`
	CedingCommission            float64 `json:"-" csv:"-"`
	RiderSumAssured             float64 `json:"-" csv:"-"`
}
