package models

type ProductFeatures struct {
	ID                              int    `json:"-" gorm:"primary_key"`
	ProductCode                     string `json:"-" gorm:"unique;size:255"`
	WholeOfLife                     bool   `json:"whole_of_life"`
	NonLife                         bool   `json:"non_life"`
	FuneralCover                    bool   `json:"funeral_cover"`
	CreditLife                      bool   `json:"credit_life"`
	SaOutstandingLoan               bool   `json:"sa_outstanding_loan"`
	SaFixedBaseLumpSum              bool   `json:"sa_fixed_base_lump_sum"`
	PremiumHoliday                  bool   `json:"premium_holiday"`
	ProductLevelEscalations         bool   `json:"product_level_escalations"`
	AgeRatedEscalations             bool   `json:"age_rated_escalations"`
	StandardAdditionalLumpSum       bool   `json:"standard_additional_lump_sum"`
	RiderBenefit                    bool   `json:"rider_benefit"`
	AccidentalDeathBenefit          bool   `json:"accidental_death_benefit"`
	LapseDependentOnCpDeath         bool   `json:"lapse_dependent_on_cp_death"`
	CashBackSurvival                bool   `json:"cash_back_survival"`
	CashBackOnDeath                 bool   `json:"cash_back_on_death"`
	RetrenchmentBenefit             bool   `json:"retrenchment_benefit"`
	TemporaryDisabilityBenefit      bool   `json:"temporary_disability_benefit"`
	CreditLifeFlatPremium           bool   `json:"credit_life_flat_premium"`
	CreditLifeDecreasingPremium     bool   `json:"credit_life_decreasing_premium"`
	OsProjPvMethod                  bool   `json:"os_proj_pv_method"`
	ProportionalReinsurance         bool   `json:"proportional_reinsurance"`
	JointLife                       bool   `json:"joint_life"`
	AnnuityIncome                   bool   `json:"annuity_income"`
	AnnuityCalendarMonthEscalations bool   `json:"annuity_calendar_month_escalations"`
	FullyPaidup                     bool   `json:"fully_paidup"`
	ReducedPaidup                   bool   `json:"reduced_paidup"`
	PremiumWaiverOnDeath            bool   `json:"premium_waiver_on_death"`
	PremiumWaiverOnDisability       bool   `json:"premium_waiver_on_disability"`
	SurrenderBenefit                bool   `json:"surrender_benefit"`
	MaturityBenefit                 bool   `json:"maturity_benefit"`
	UnitFund                        bool   `json:"unit_fund"`
	WithProfit                      bool   `json:"with_profit"`
	FundRiskCharge                  bool   `json:"fund_risk_charge"`
	AdvisoryFeeCharge               bool   `json:"advisory_fee_charge"`
	PermanentDisabilityBenefit      bool   `json:"permanent_disability_benefit"`
	DeathBenefit                    bool   `json:"death_benefit"`
	CriticalIllnessBenefit          bool   `json:"critical_illness_benefit"`
	CreditLifeFuneralOption         bool   `json:"credit_life_funeral_option"`
	SpecialDecrementMargin          bool   `json:"special_decrement_margin"`
	RenewableProfitAdjustment       bool   `json:"renewable_profit_adjustment"`
	SurrenderValueQuadraticFormula  bool   `json:"surrender_value_quadratic_formula"`
	MaturityBenefitPattern          bool   `json:"maturity_benefit_pattern"`
}

type ProductNonLifeRating struct {
	ID             int     `json:"-" gorm:"primary_key"`
	Year           int     `json:"year"`
	ProductCode    string  `json:"product_code"`
	DurationIfM    int     `json:"duration_if_m"`
	AnnualRiskRate float64 `json:"annual_risk_rate"`
}

type ProductBenefitMultiplier struct {
	ID                      int     `json:"-" gorm:"primary_key"`
	ProductCode             string  `json:"product_code" csv:"product_code"`
	Plan                    string  `json:"plan" csv:"plan"`
	AnnuityIncomeMultiplier float64 `json:"annuity_income_multiplier" csv:"annuity_income_multiplier"`
	SalaryMultiplier        float64 `json:"salary_multiplier" csv:"salary_multiplier"`
}

func (ProductBenefitMultiplier) TableName() string {
	return "product_benefit_multiplier"
}
