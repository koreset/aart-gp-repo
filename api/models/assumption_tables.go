package models

type Termination struct {
}

type ProductSpecialDecrementMargin struct { // may need to rename to special decrement adjustments instead of margins
	ID                  int     `gorm:"primary_key" json:"id"`
	ProductCode         string  `json:"product_code" gorm:"index" csv:"product_code"`
	MemberType          string  `json:"member_type" gorm:"index" csv:"member_type"`
	SpecialMarginCode   string  `json:"special_margin_code" gorm:"index" csv:"special_margin_code"`
	Anb                 int     `json:"anb" gorm:"index" csv:"anb"`
	MortalityMargin     float64 `json:"mortality_margin" csv:"mortality_margin"`
	MorbidityMargin     float64 `json:"morbidity_margin" csv:"morbidity_margin"`
	RetrenchmentMargin  float64 `json:"retrenchment_margin" csv:"retrenchment_margin"`
	MortalityTableProp  float64 `json:"mortality_table_prop" csv:"mortality_table_prop"`
	DisabilityTableProp float64 `json:"disability_table_prop" csv:"disability_table_prop"`
	LapseTableProp      float64 `json:"lapse_table_prop" csv:"lapse_table_prop"`
	Basis               string  `json:"basis" gorm:"index" csv:"basis"`
}

type ProductPricingSpecialDecrementMargin struct {
	ID                  int     `gorm:"primary_key" json:"id"`
	ProductCode         string  `json:"product_code" csv:"product_code"`
	MemberType          string  `json:"member_type" csv:"member_type"`
	SpecialMarginCode   string  `json:"special_margin_code" csv:"special_margin_code"`
	Anb                 int     `json:"anb" csv:"anb"`
	MortalityMargin     float64 `json:"mortality_margin" csv:"mortality_margin"`
	MorbidityMargin     float64 `json:"morbidity_margin" csv:"morbidity_margin"`
	RetrenchmentMargin  float64 `json:"retrenchment_margin" csv:"retrenchment_margin"`
	MortalityTableProp  float64 `json:"mortality_table_prop" csv:"mortality_table_prop"`
	DisabilityTableProp float64 `json:"disability_table_prop" csv:"disability_table_prop"`
	LapseTableProp      float64 `json:"lapse_table_prop" csv:"lapse_table_prop"`
	Basis               string  `json:"basis" csv:"basis"`
}

type ProductRenewableProfitAdjustment struct {
	ID                        int     `gorm:"primary_key" json:"id"`
	ProductCode               string  `json:"product_code" csv:"product_code"`
	DurationIfY               float64 `json:"duration_if_y" csv:"duration_if_y"`
	ProfitAdjustmentCode      string  `json:"profit_adjustment_code" csv:"profit_adjustment_code"`
	RenewableProfitAdjustment float64 `json:"renewable_profit_adjustment" csv:"renewable_profit_adjustment"`
}

type ProductPricingRenewableProfitAdjustment struct {
	ID                        int     `gorm:"primary_key" json:"id"`
	ProductCode               string  `json:"product_code" csv:"product_code"`
	DurationIfY               float64 `json:"duration_if_y" csv:"duration_if_y"`
	ProfitAdjustmentCode      string  `json:"profit_adjustment_code" csv:"profit_adjustment_code"`
	RenewableProfitAdjustment float64 `json:"renewable_profit_adjustment" csv:"renewable_profit_adjustment"`
}

type ProductInvestmentReturn struct {
	ID                       int     `gorm:"primary_key" json:"id"`
	ProductCode              string  `json:"product_code" csv:"product_code"`
	FundYear                 int     `json:"fund_year" csv:"fund_year"`
	FundCode                 string  `json:"fund_code" csv:"fund_code"`
	FundReversionaryBonus    float64 `json:"fund_reversionary_bonus" csv:"fund_reversionary_bonus"`
	FundTerminalBonus        float64 `json:"fund_terminal_bonus" csv:"fund_terminal_bonus"`
	FundGrowthMargin         float64 `json:"fund_growth_margin" csv:"fund_growth_margin"`
	CashYieldGap             float64 `json:"cash_yield_gap" csv:"cash_yield_gap"`
	CorporateBondRiskPremium float64 `json:"corporate_bond_risk_premium" csv:"corporate_bond_risk_premium"`
	EquityRiskPremium        float64 `json:"equity_risk_premium" csv:"equity_risk_premium"`
	PropertyRiskPremium      float64 `json:"property_risk_premium" csv:"property_risk_premium"`
	Year                     int     `json:"year" csv:"year"`
}

type ProductUnitFundCharge struct {
	ID                               int     `gorm:"primary_key" json:"id"`
	ProductCode                      string  `json:"product_code" csv:"product_code"`
	FundCode                         string  `json:"fund_code" csv:"fund_code"`
	DurationIfY                      int     `json:"duration_if_y" csv:"duration_if_y"`
	PremiumAdvisoryFeeRate           float64 `json:"premium_advisory_fee_rate" csv:"premium_advisory_fee_rate"`
	FundAdvisoryFeeRate              float64 `json:"fund_advisory_fee_rate" csv:"fund_advisory_fee_rate"`
	AnnualAdvisoryFeeAmount          float64 `json:"annual_advisory_fee_amount" csv:"annual_advisory_fee_amount"`
	PremiumPolicyFeeRate             float64 `json:"premium_policy_fee_rate" csv:"premium_policy_fee_rate"`
	AnnualPolicyFeeAmount            float64 `json:"annual_policy_fee_amount" csv:"annual_policy_fee_amount"`
	FundManagementChargeRate         float64 `json:"fund_management_charge_rate" csv:"fund_management_charge_rate"`
	AnnualFundManagementChargeAmount float64 `json:"annual_fund_management_charge_amount" csv:"annual_fund_management_charge_amount"`
	PremiumAllocationRate            float64 `json:"premium_allocation_rate" csv:"premium_allocation_rate"`
	FundSurrenderPenaltyRate         float64 `json:"fund_surrender_penalty_rate" csv:"fund_surrender_penalty_rate"`
	FundMinimumSurrenderValueRate    float64 `json:"fund_minimum_surrender_value_rate" csv:"fund_minimum_surrender_value_rate"`
	BidOfferSpread                   float64 `json:"bid_offer_spread" csv:"bid_offer_spread"`
	MarketValueAdjustment            float64 `json:"market_value_adjustment" csv:"market_value_adjustment"`
	PartSurrenderFeeRate             float64 `json:"part_surrender_fee_rate" csv:"part_surrender_fee_rate"`
	BsaShareholderFeeRate            float64 `json:"bsa_shareholder_fee_rate" csv:"bsa_shareholder_fee_rate"`
}

type ProductFundAssetDistribution struct {
	ID             int     `gorm:"primary_key" json:"id"`
	ProductCode    string  `json:"product_code" csv:"product_code"`
	FundCode       string  `json:"fund_code" csv:"fund_code"`
	GovernmentBond float64 `json:"government_bond" csv:"government_bond"`
	Cash           float64 `json:"cash" csv:"cash"`
	CorporateBond  float64 `json:"corporate_bond" csv:"corporate_bond"`
	Equity         float64 `json:"equity" csv:"equity"`
	Property       float64 `json:"property" csv:"property"`
}

type ProductMaturityPattern struct {
	ID                  int     `gorm:"primary_key" json:"id"`
	ProductCode         string  `json:"product_code" csv:"product_code"`
	MaturityPatternCode string  `json:"maturity_pattern_code" csv:"maturity_pattern_code"`
	RemainingTermYear   int     `json:"remaining_term_year" csv:"remaining_term_year"`
	PartMaturityRate    float64 `json:"part_maturity_rate" csv:"part_maturity_rate"`
}

type ProductSurrenderValueCoefficient struct {
	ID                                 int     `gorm:"primary_key" json:"id"`
	ProductCode                        string  `json:"product_code" csv:"product_code"`
	SurrenderValueCode                 string  `json:"surrender_value_code" csv:"surrender_value_code"`
	A                                  float64 `json:"a" csv:"a"`
	B                                  float64 `json:"b" csv:"b"`
	C                                  float64 `json:"c" csv:"c"`
	SurrenderPayoutWaitingPeriodMonths int     `json:"surrender_payout_waiting_period_months" csv:"surrender_payout_waiting_period_months"`
}
