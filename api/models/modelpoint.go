package models

import (
	"errors"
	"fmt"
	"time"
)

type ProductModelPoint struct {
	Year                            int       `json:"year" gorm:"index:idx_year" csv:"year"`
	PolicyNumber                    string    `json:"policy_number" csv:"policy_number"`
	Spcode                          int       `json:"spcode" csv:"spcode"`
	IFRS17Group                     string    `json:"ifrs17_group" csv:"ifrs17_group"`
	ProductCode                     string    `json:"product_code" csv:"product_code"`
	LockedInYear                    int       `json:"locked_in_year" csv:"locked_in_year"`
	LockedInMonth                   int       `json:"locked_in_month" csv:"locked_in_month"`
	DurationInForceMonths           int       `json:"duration_in_force_months" csv:"duration_in_force_months"`
	AgeAtEntry                      int       `json:"age_at_entry" csv:"age_at_entry"`
	Gender                          string    `json:"gender" csv:"gender"`
	MainMemberAgeAtEntry            int       `json:"main_member_age_at_entry" csv:"main_member_age_at_entry"`
	MainMemberGender                string    `json:"main_member_gender" csv:"main_member_gender"`
	InitialPolicy                   int       `json:"initial_policy" csv:"initial_policy"`
	SumAssured                      float64   `json:"sum_assured" csv:"sum_assured"`
	UnitFund                        float64   `json:"unit_fund" csv:"unit_fund"`
	ReversionaryBonus               float64   `json:"reversionary_bonus" csv:"reversionary_bonus"`
	GuaranteedMaturityBenefit       float64   `json:"guaranteed_maturity_benefit" csv:"guaranteed_maturity_benefit"`
	BonusStabilisationAccount       float64   `json:"bonus_stabilisation_account" csv:"bonus_stabilisation_account"`
	OriginalLoan                    float64   `json:"original_loan" csv:"original_loan"`
	OutstandingLoan                 float64   `json:"outstanding_loan" csv:"outstanding_loan"`
	Interest                        float64   `json:"interest" csv:"interest"`
	Instalment                      float64   `json:"instalment" csv:"instalment"`
	AnnuityIncome                   float64   `json:"annuity_income" csv:"annuity_income"`
	LifeAnnuityAmount               float64   `json:"life_annuity_amount" csv:"life_annuity_amount"`
	LifeAnnuityPercentage           float64   `json:"life_annuity_percentage" csv:"life_annuity_percentage"`
	Term                            int       `json:"term" csv:"term"`
	OriginalTerm                    int       `json:"original_term" csv:"original_term"`
	OutstandingTermMonths           int       `json:"outstanding_term_months" csv:"outstanding_term_months"`
	AnnualPremium                   float64   `json:"annual_premium" csv:"annual_premium"`
	PremiumRate                     float64   `json:"premium_rate" csv:"premium_rate"`
	PremiumFrequency                int       `json:"premium_frequency" csv:"premium_frequency"`
	PremiumStatus                   float64   `json:"premium_status" csv:"premium_status"`
	CommissionType                  string    `json:"commission_type" csv:"commission_type"`
	WaitingPeriod                   int       `json:"waiting_period" csv:"waiting_period"`
	DeferredPeriod                  int       `json:"deferred_period" csv:"deferred_period"`
	MemberType                      string    `json:"member_type" csv:"member_type"`
	SumAssuredEscalation            float64   `json:"sum_assured_escalation" csv:"sum_assured_escalation"`
	PremiumEscalation               float64   `json:"premium_escalation" csv:"premium_escalation"`
	EscalationMonth                 int       `json:"escalation_month" csv:"escalation_month"`
	AnnuityEscalation               float64   `json:"annuity_escalation" csv:"annuity_escalation"`
	AnnuityEscalationMonth          int       `json:"annuity_escalation_month" csv:"annuity_escalation_month"`
	Plan                            string    `json:"plan" csv:"plan"`
	FundCode                        string    `json:"fund_code" csv:"fund_code"`
	BsaFundCode                     string    `json:"bsa_fund_code" csv:"bsa_fund_code"`
	MaturityBenefitCode             string    `json:"maturity_benefit_code" csv:"maturity_benefit_code"`
	SurrenderValueCode              string    `json:"surrender_value_code" csv:"surrender_value_code"`
	DistributionChannel             string    `json:"distribution_channel" csv:"distribution_channel"`
	TaxClass                        int       `json:"tax_class" csv:"tax_class"`
	PaidupOption                    bool      `json:"paidup_option" csv:"paidup_option"`
	PaidupIndicator                 bool      `json:"paidup_indicator" csv:"paidup_indicator"`
	TemporaryPremiumWaiverIndicator bool      `json:"temporary_premium_waiver_indicator" csv:"temporary_premium_waiver_indicator"`
	TemporaryPremiumWaiverMonthExit int       `json:"temporary_premium_waiver_month_exit" csv:"temporary_premium_waiver_month_exit"`
	ContinuityOrPremiumWaiverOption bool      `json:"continuity_or_premium_waiver_option" csv:"continuity_or_premium_waiver_option"`
	PremiumWaiverIndicator          bool      `json:"premium_waiver_indicator" csv:"premium_waiver_indicator"`
	EducatorOption                  bool      `json:"educator_option" csv:"educator_option"`
	EducatorWaitingPeriod           int       `json:"educator_waiting_period" csv:"educator_waiting_period"`
	Income                          int       `json:"income" csv:"income"`
	CashbackOption                  bool      `json:"cashback_option" csv:"cashback_option"`
	CashbackIndicator               bool      `json:"cashback_indicator" csv:"cashback_indicator"`
	Grocery                         bool      `json:"grocery" csv:"grocery"`
	Repatriation                    bool      `json:"repatriation" csv:"repatriation"`
	Tombstone                       bool      `json:"tombstone" csv:"tombstone"`
	CowBenefit                      bool      `json:"cow_benefit" csv:"cow_benefit"`
	AdditionalSumAssuredIndicator   bool      `json:"additional_sum_assured_indicator" csv:"additional_sum_assured_indicator"`
	PremiumHolidayUsed              int       `json:"premium_holiday_used" csv:"premium_holiday_used"`
	EducationLevel                  int       `json:"education_level" csv:"education_level"`
	SocioEconomicClass              int       `json:"socio_economic_class" csv:"socio_economic_class"`
	OccupationalClass               string    `json:"occupational_class" csv:"occupational_class"`
	SmokerStatus                    string    `json:"smoker_status" csv:"smoker_status"`
	SelectPeriod                    int       `json:"select_period" csv:"select_period"`
	DisabilityDefinition            int       `json:"disability_definition" csv:"disability_definition"`
	TreatyYear                      int       `json:"treaty_year" csv:"treaty_year"`
	Weighting                       float64   `json:"weighting" csv:"weighting"`
	PricingMpVersion                string    `json:"pricing_mp_version" csv:"pricing_mp_version"`
	MpVersion                       string    `json:"mp_version" csv:"mp_version" gorm:"index:idx_mp_version"`
	CreationDate                    time.Time `json:"creation_date" csv:"creation_date"`
}

type ProductModelPointVariableStats struct {
	ID             int     `json:"id" gorm:"primary_key"`
	ProductCode    string  `json:"product_code" csv:"product_code" gorm:"index:idx_product_code"`
	Variable       string  `json:"variable"`
	Min            float64 `json:"min"`
	Max            float64 `json:"max"`
	Sum            float64 `json:"sum"`
	Average        float64 `json:"average"`
	Male           float64 `json:"male"`
	Female         float64 `json:"female"`
	NumberOfZeroes int     `json:"number_of_zeroes"`
	NumberOfLives  int     `json:"number_of_lives"`
	Year           int     `json:"year" gorm:"index:idx_year"`
	Version        string  `json:"version" gorm:"index:idx_version"`
}

type PaaModelPointVariableStats struct {
	ID             int     `json:"id" gorm:"primary_key"`
	PortfolioName  string  `json:"portfolio_name" csv:"portfolio_name"`
	Year           int     `json:"year" csv:"year"`
	MpVersion      string  `json:"mp_version" csv:"mp_version"`
	Variable       string  `json:"variable"`
	Min            float64 `json:"min"`
	Max            float64 `json:"max"`
	Average        float64 `json:"average"`
	NumberOfZeroes int     `json:"number_of_zeroes"`
	EmptyValues    int     `json:"empty_values"`
	TotalCount     int     `json:"total_count"`
	DistinctValues int     `json:"distinct_values"`
}

type ConvertibleBoolean bool

func (bit ConvertibleBoolean) UnmarshalJSON(data []byte) error {
	asString := string(data)
	if asString == "1" || asString == "true" {
		bit = true
	} else if asString == "0" || asString == "false" {
		bit = false
	} else {
		return errors.New(fmt.Sprintf("Boolean unmarshal error: invalid input %s", asString))
	}
	return nil
}

//func (bit *ConvertibleBoolean) Scan(value interface{}) error {
//	val, ok := value.(uint8)
//	if !ok {
//		return errors.New("unsupported scan")
//	}
//
//	if val == 0 {
//		*bit = false
//		return nil
//	}
//
//	if val == 1 {
//		*bit = true
//		return nil
//	}
//	return errors.New("the value was not valid for a scan into ConvertibleBoolean")
//}
