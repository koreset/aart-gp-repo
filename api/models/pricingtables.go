package models

import "gorm.io/gorm"

type ProductPricingModelPoint struct {
	PolicyNumber                    string  `json:"policy_number" csv:"policy_number"`
	ProductCode                     string  `json:"product_code" csv:"product_code"`
	Spcode                          int     `json:"spcode" csv:"spcode"`
	IFRS17Group                     string  `json:"ifrs17_group" csv:"ifrs17_group"`
	InitialPolicy                   int     `json:"initial_policy" csv:"initial_policy"`
	LockedInYear                    int     `json:"locked_in_year" csv:"locked_in_year"`
	LockedInMonth                   int     `json:"locked_in_month" csv:"locked_in_month"`
	DurationInForceMonths           int     `json:"duration_in_force_months" csv:"duration_in_force_months"`
	AgeAtEntry                      int     `json:"age_at_entry" csv:"age_at_entry"`
	Gender                          string  `json:"gender" csv:"gender"`
	MainMemberAgeAtEntry            int     `json:"main_member_age_at_entry" csv:"main_member_age_at_entry"`
	MainMemberGender                string  `json:"main_member_gender" csv:"main_member_gender"`
	SumAssured                      float64 `json:"sum_assured" csv:"sum_assured"`
	UnitFund                        float64 `json:"unit_fund" csv:"unit_fund"`
	ReversionaryBonus               float64 `json:"reversionary_bonus" csv:"reversionary_bonus"`
	GuaranteedMaturityBenefit       float64 `json:"guaranteed_maturity_benefit" csv:"guaranteed_maturity_benefit"`
	BonusStabilisationAccount       float64 `json:"bonus_stabilisation_account" csv:"bonus_stabilisation_account"`
	OriginalLoan                    float64 `json:"original_loan" csv:"original_loan"`
	OutstandingLoan                 float64 `json:"outstanding_loan" csv:"outstanding_loan"`
	Interest                        float64 `json:"interest" csv:"interest"`
	Instalment                      float64 `json:"instalment" csv:"instalment"`
	AnnuityIncome                   float64 `json:"annuity_income" csv:"annuity_income"`
	LifeAnnuityAmount               float64 `json:"life_annuity_amount" csv:"life_annuity_amount"`
	LifeAnnuityPercentage           float64 `json:"life_annuity_percentage" csv:"life_annuity_percentage"`
	Term                            int     `json:"term" csv:"term"`
	OriginalTerm                    int     `json:"original_term" csv:"original_term"`
	OutstandingTermMonths           int     `json:"outstanding_term_months" csv:"outstanding_term_months"`
	AnnualPremium                   float64 `json:"annual_premium" csv:"annual_premium"`
	PremiumRate                     float64 `json:"premium_rate" csv:"premium_rate"`
	PremiumFrequency                int     `json:"premium_frequency" csv:"premium_frequency"`
	PremiumStatus                   float64 `json:"premium_status" csv:"premium_status"`
	CommissionType                  string  `json:"commission_type" csv:"commission_type"`
	WaitingPeriod                   int     `json:"waiting_period" csv:"waiting_period"`
	DeferredPeriod                  int     `json:"deferred_period" csv:"deferred_period"`
	MemberType                      string  `json:"member_type" csv:"member_type"`
	SumAssuredEscalation            float64 `json:"sum_assured_escalation" csv:"sum_assured_escalation"`
	PremiumEscalation               float64 `json:"premium_escalation" csv:"premium_escalation"`
	EscalationMonth                 int     `json:"escalation_month" csv:"escalation_month"`
	AnnuityEscalation               float64 `json:"annuity_escalation" csv:"annuity_escalation"`
	AnnuityEscalationMonth          float64 `json:"annuity_escalation_month" csv:"annuity_escalation_month"`
	Plan                            string  `json:"plan" csv:"plan"`
	FundCode                        string  `json:"fund_code" csv:"fund_code"`
	BsaFundCode                     string  `json:"bsa_fund_code" csv:"bsa_fund_code"`
	MaturityBenefitCode             string  `json:"maturity_benefit_code" csv:"maturity_benefit_code"`
	SurrenderValueCode              string  `json:"surrender_value_code" csv:"surrender_value_code"`
	DistributionChannel             string  `json:"distribution_channel" csv:"distribution_channel"`
	TaxClass                        int     `json:"tax_class" csv:"tax_class"`
	PaidupOption                    bool    `json:"paidup_option" csv:"paidup_option"`
	PaidupIndicator                 bool    `json:"paidup_indicator" csv:"paidup_indicator"`
	TemporaryPremiumWaiverIndicator bool    `json:"temporary_premium_waiver_indicator" csv:"temporary_premium_waiver_indicator"`
	TemporaryPremiumWaiverMonthExit int     `json:"temporary_premium_waiver_month_exit" csv:"temporary_premium_waiver_month_exit"`
	ContinuityOrPremiumWaiverOption bool    `json:"continuity_or_premium_waiver_option" csv:"continuity_or_premium_waiver_option"`
	PremiumWaiverIndicator          bool    `json:"premium_waiver_indicator" csv:"premium_waiver_indicator"`
	EducatorOption                  bool    `json:"educator_option" csv:"educator_option"`
	EducatorWaitingPeriod           int     `json:"educator_waiting_period" csv:"educator_waiting_period"`
	Income                          int     `json:"income" csv:"income"`
	CashbackOption                  bool    `json:"cashback_option" csv:"cashback_option"`
	CashbackIndicator               bool    `json:"cashback_indicator" csv:"cashback_indicator"`
	Grocery                         bool    `json:"grocery" csv:"grocery"`
	Repatriation                    bool    `json:"repatriation" csv:"repatriation"`
	Tombstone                       bool    `json:"tombstone" csv:"tombstone"`
	CowBenefit                      bool    `json:"cow_benefit" csv:"cow_benefit"`
	AdditionalSumAssuredIndicator   bool    `json:"additional_sum_assured_indicator" csv:"additional_sum_assured_indicator"`
	PremiumHolidayUsed              int     `json:"premium_holiday_used" csv:"premium_holiday_used"`
	EducationLevel                  int     `json:"education_level" csv:"education_level"`
	SocioEconomicClass              int     `json:"socio_economic_class" csv:"socio_economic_class"`
	OccupationalClass               string  `json:"occupational_class" csv:"occupational_class"`
	SmokerStatus                    string  `json:"smoker_status" csv:"smoker_status"`
	SelectPeriod                    int     `json:"select_period" csv:"select_period"`
	DisabilityDefinition            int     `json:"disability_definition" csv:"disability_definition"`
	Weighting                       float64 `json:"weighting" csv:"weighting"`
	TreatyYear                      float64 `json:"treaty_year" csv:"treaty_year"`
	PricingMpVersion                string  `json:"pricing_mp_version" gorm:"index"  csv:"pricing_mp_version"`
	Created                         int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type ProductPricingMargins struct {
	ProductCode        string  `json:"product_code" gorm:"primary_key;auto_increment:false" `
	MortalityMargin    float64 `json:"mortality_margin" csv:"mortality_margin"`
	MorbidityMargin    float64 `json:"morbidity_margin" csv:"morbidity_margin"`
	ExpenseMargin      float64 `json:"expense_margin" csv:"expense_margin"`
	RetrenchmentMargin float64 `json:"retrenchment_margin" csv:"retrenchment_margin"`
	InflationMargin    float64 `json:"inflation_margin" csv:"inflation_margin"`
	InvestmentMargin   float64 `json:"investment_margin" csv:"investment_margin"`
	Year               float64 `json:"year" csv:"year"`
	Created            int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type ProductPricingNewBusinessProfile struct {
	gorm.Model
	ProductCode      string  `json:"product_code" csv:"product_code"`
	PricingMpVersion string  `json:"pricing_mp_version" csv:"pricing_mp_version"`
	NbYear1          float64 `json:"nb_year_1" csv:"nb_year_1"`
	NbYear2          float64 `json:"nb_year_2" csv:"nb_year_2"`
	NbYear3          float64 `json:"nb_year_3" csv:"nb_year_3"`
	NbYear4          float64 `json:"nb_year_4" csv:"nb_year_4"`
	NbYear5          float64 `json:"nb_year_5" csv:"nb_year_5"`
	NbYear6          float64 `json:"nb_year_6" csv:"nb_year_6"`
	NbYear7          float64 `json:"nb_year_7" csv:"nb_year_7"`
	NbYear8          float64 `json:"nb_year_8" csv:"nb_year_8"`
	NbYear9          float64 `json:"nb_year_9" csv:"nb_year_9"`
	NbYear10         float64 `json:"nb_year_10" csv:"nb_year_10"`
	Created          int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type ProductPricingProfitMargin struct {
	CustomGormModel
	ProductCode  string  `json:"product_code" csv:"product_code"`
	ANB          float64 `json:"anb" csv:"anb"`
	SumAssured   float64 `json:"sum_assured" csv:"sum_assured"`
	ProfitMargin float64 `json:"profit_margin" csv:"profit_margin"`
	Created      int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

type ProductPricingProductLevelEscalation struct {
	ID                       int     `gorm:"primary_key" json:"id"`
	ProductCode              string  `json:"product_code" csv:"product_code" gorm:"size:255"`
	Basis                    string  `json:"basis" gorm:"primary_key;auto_increment:false" csv:"basis"` //Basis                    string  `json:"basis" gorm:"uniqueIndex" csv:"basis"`
	ANB                      float64 `json:"anb" csv:"anb"`
	SumAssuredEscalationRate float64 `json:"sum_assured_escalation_rate" csv:"sum_assured_escalation_rate"`
	PremiumEscalationRate    float64 `json:"premium_escalation_rate" csv:"premium_escalation_rate"`
	Created                  int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}
type ProductPricingParameters struct {
	ProductCode                         string  `json:"product_code" gorm:"primary_key;auto_increment:false;size:255" csv:"product_code"`
	Basis                               string  `json:"basis" gorm:"primary_key;auto_increment:false" csv:"basis"`
	InitialExpensePercentage            float64 `json:"initial_expense_percentage" csv:"initial_expense_percentage"`
	InitialExpenseRand                  float64 `json:"initial_expense_rand" csv:"initial_expense_rand"`
	RenewalExpensePercentage            float64 `json:"renewal_expense_percentage" csv:"renewal_expense_percentage"`
	RenewalExpenseRand                  float64 `json:"renewal_expense_rand" csv:"renewal_expense_rand"`
	PremiumWaiverWaitingPeriod          int     `json:"premium_waiver_waiting_period" csv:"premium_waiver_waiting_period"`
	PaidUpOnSurvivalWaitingPeriod       int     `json:"paid_up_on_survival_waiting_period" csv:"paid_up_on_survival_waiting_period"`
	EducatorSumAssuredPaymentAge        int     `json:"educator_sum_assured_payment_age" csv:"educator_sum_assured_payment_age"`
	ChildExitAge                        int     `json:"child_exit_age" csv:"child_exit_age"`
	Timing                              float64 `json:"timing" csv:"timing"`
	PremiumWaiverSumAssuredFactor       float64 `json:"premium_waiver_sum_assured_factor" csv:"premium_waiver_sum_assured_factor"`
	PremiumWaiverAdjustedSumassuredTerm float64 `json:"premium_waiver_adjusted_sumassured_term" csv:"premium_waiver_adjusted_sumassured_term"`
	PremiumWaiverExitAge                float64 `json:"premium_waiver_exit_age" csv:"premium_waiver_exit_age"`
	CashbackOnSurvivalPeriod            float64 `json:"cashback_on_survival_period" csv:"cashback_on_survival_period"`
	CashbackOnSurvivalTerm              float64 `json:"cashback_on_survival_term" csv:"cashback_on_survival_term"`
	PaidupEffectiveAge                  int     `json:"paidup_effective_age" csv:"paidup_effective_age"`
	PaidupEffectiveDuration             float64 `json:"paidup_effective_duration" csv:"paidup_effective_duration"`
	CashbackOnSurvivalRatio             float64 `json:"cashback_on_survival_ratio" csv:"cashback_on_survival_ratio"`
	CashbackOnDeathRatio                float64 `json:"cashback_on_death_ratio" csv:"cashback_on_death_ratio"`
	CashbackOnDeathTerm                 float64 `json:"cashback_on_death_term" csv:"cashback_on_death_term"`
	ClaimsExpensePercentage             float64 `json:"claims_expense_percentage" csv:"claims_expense_percentage"`
	ClaimsExpenseRand                   float64 `json:"claims_expense_rand" csv:"claims_expense_rand"`
	StandardAdditionalSumAssured        float64 `json:"standard_additional_sum_assured" csv:"standard_additional_sum_assured"`
	RiderBenefits                       float64 `json:"rider_benefits" csv:"rider_benefits"`
	GroceryBenefitMultiplier            int     `json:"grocery_benefit_multiplier" csv:"grocery_benefit_multiplier"`
	RepatriationAssumption              float64 `json:"repatriation_assumption" csv:"repatriation_assumption"`
	OriginalSumAssured                  int     `json:"original_sum_assured" csv:"original_sum_assured"`
	RetrenchmentQualificationFactor     float64 `json:"retrenchment_qualification_factor" csv:"retrenchment_qualification_factor"`
	MainMemberIndicator                 bool    `json:"main_member_indicator" csv:"main_member_indicator"`
	OtherLivesIndicator                 bool    `json:"other_lives_indicator" csv:"other_lives_indicator"`
	NumberRetrenchmentInstalments       float64 `json:"number_retrenchment_instalments" csv:"number_retrenchment_instalments"`
	RiderFuneralSaProportion            float64 `json:"rider_funeral_sa_proportion" csv:"rider_funeral_sa_proportion"`
	RiderFuneralMinimumBenefit          float64 `json:"rider_funeral_minimum_benefit" csv:"rider_funeral_minimum_benefit"`
	RiderFuneralMaximumBenefit          float64 `json:"rider_funeral_maximum_benefit" csv:"rider_funeral_maximum_benefit"`
	PremiumHolidayWaitingPeriod         int     `json:"premium_holiday_waiting_period" csv:"premium_holiday_waiting_period"`
	MaximumPremiumHolidays              int     `json:"maximum_premium_holidays" csv:"maximum_premium_holidays"`
	PremiumHolidayCycle                 int     `json:"premium_holiday_cycle" csv:"premium_holiday_cycle"`
	CorporateTaxRate                    float64 `json:"corporate_tax_rate" csv:"corporate_tax_rate"`
	ShareholdersRequiredMargin          float64 `json:"shareholders_required_margin" csv:"shareholders_required_margin"`
	Created                             int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
	AccidentalBenefitMultiplier         float64 `json:"-" csv:"-"`
	CalculatedTerm                      int     `json:"-" csv:"-"`
	CalculatedInstalment                float64 `json:"-" csv:"-"`
}

type ProductPricingTable struct {
	ID        int    `gorm:"primary_key" json:"id"`
	ProductID int    `json:"product_id"`
	Class     string `json:"table_class"`
	Name      string `json:"table"`
	Populated bool   `json:"populated" gorm:"-"`
}

type ProductPricingChildSumAssured struct {
	ProductCode string `gorm:"size:255"`
	Age         int
	A           float64
	B           float64
	C           float64
	D           float64
	E           float64
	F           float64
	G           float64
	H           float64
	I           float64
	J           float64
	K           float64
	L           float64
	M           float64
	N           float64
	O           float64
}

func (ProductPricingChildSumAssured) TableName() string {
	return "product_pricing_child_sum_assured"
}

type ProductPricingChildAdditionalSumAssured struct {
	ProductCode string `gorm:"size:255"`
	Age         int
	A           float64
	B           float64
	C           float64
	D           float64
	E           float64
	F           float64
	G           float64
	H           float64
	I           float64
	J           float64
	K           float64
	L           float64
	M           float64
	N           float64
	O           float64
}

func (ProductPricingChildAdditionalSumAssured) TableName() string {
	return "product_pricing_child_additional_sum_assured"
}

type ProductPricingAdditionalSumAssured struct {
	ProductCode string `gorm:"size:255"`
	A           float64
	B           float64
	C           float64
	D           float64
	E           float64
	F           float64
	G           float64
	H           float64
	I           float64
	J           float64
	K           float64
	L           float64
	M           float64
	N           float64
	O           float64
}

func (ProductPricingAdditionalSumAssured) TableName() string {
	return "product_pricing_funeral_service"
}

type ProductPricingRider struct {
	ProductCode  string  `json:"product_code"`
	RiderBenefit string  `json:"rider_benefit"`
	A            float64 `json:"a"`
	B            float64 `json:"b"`
	C            float64 `json:"c"`
	D            float64 `json:"d"`
	E            float64 `json:"e"`
	F            float64 `json:"f"`
	G            float64 `json:"g"`
	H            float64 `json:"h"`
	I            float64 `json:"i"`
	J            float64 `json:"j"`
	K            float64 `json:"k"`
	L            float64 `json:"l"`
	M            float64 `json:"m"`
	N            float64 `json:"n"`
	O            float64 `json:"o"`
}

type ProductPricingClawback struct {
	ProductCode            string  `json:"product_code"`
	DurationInForceMonth   int     `json:"duration_in_force_month"`
	Year1InitialCommission float64 `json:"year_1_initial_commission"`
	Year2InitialCommission float64 `json:"year_2_initial_commission"`
}

func (ProductPricingClawback) TableName() string {
	return "product_pricing_clawback"
}

type ProductPricingLapseMargin struct {
	ProductCode string  `json:"product_code" gorm:"size:255"`
	Month       int     `json:"month"`
	Margin      float64 `json:"margin"`
	Basis       string  `json:"basis"`
}

type PricingYieldCurve struct {
	ID             int     `gorm:"primary_key" json:"id"`
	YieldCurveCode string  `json:"yield_curve_code" csv:"yield_curve_code"`
	ProjectionTime int     `json:"proj_time" csv:"proj_time" gorm:"primary_key;auto_increment:false;column:proj_time"`
	Month          int     `json:"month" csv:"month" gorm:"primary_key;auto_increment:false;column:month"`
	NominalRate    float64 `json:"nominal_rate" csv:"nominal_rate"`
	Inflation      float64 `json:"inflation" csv:"inflation"`
	//Basis          string  `json:"basis" csv:"basis"`
}

func (PricingYieldCurve) TableName() string {
	return "pricing_yield_curve"
}

//func (ProductPricingNewBusinessProfile) TableName() string {
//	return "product_pricing_new_business_profiles"
//}

type PricingRetrenchmentRate struct {
	ID                   uint    `gorm:"primary_key" json:"id"`
	DurationInForceYears uint    `json:"duration_in_force_years" gorm:"column:duration_if_y"`
	Value                float64 `json:"value" gorm:"column:retr_rate"`
}

type ProductPricingAccidentalBenefitMultiplier struct {
	ProductCode    string  `json:"product_code"`
	MainMember     float64 `json:"main_member"`
	Spouse         float64 `json:"spouse"`
	Child          float64 `json:"child"`
	Parent         float64 `json:"parent"`
	ExtendedFamily float64 `json:"extended_family"`
}

func (ProductPricingAccidentalBenefitMultiplier) TableName() string {
	return "product_pricing_accident_benefit_multiplier"
}

type ProductPricingReinsurance struct {
	Year                    int     `json:"year"`
	ID                      int     `json:"id" csv:"id" gorm:"primary_key"`
	ProductCode             string  `json:"product_code"`
	TreatyYear              int     `json:"treaty_year"`
	FlatAnnualReinsPremRate float64 `json:"flat_annual_reins_prem_rate"`
	Level1CededProportion   float64 `json:"level1_ceded_proportion"`
	Level1Lowerbound        float64 `json:"level1_lowerbound"`
	Level1Upperbound        float64 `json:"level1_upperbound"`
	Level2CededProportion   float64 `json:"level2_ceded_proportion"`
	Level2Lowerbound        float64 `json:"level2_lowerbound"`
	Level2Upperbound        float64 `json:"level2_upperbound"`
	Level3CededProportion   float64 `json:"level3_ceded_proportion"`
	Level3Lowerbound        float64 `json:"level3_lowerbound"`
	Level3Upperbound        float64 `json:"level3_upperbound"`
	CedingCommission        float64 `json:"ceding_commission"`
	Created                 int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}

func (ProductPricingReinsurance) TableName() string {
	return "product_pricing_reinsurance"
}
