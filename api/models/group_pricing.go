package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"time"
)

type Status string

const (
	StatusDraft         Status = "draft"
	StatusSubmitted     Status = "submitted"
	StatusApproved      Status = "approved"
	StatusAccepted      Status = "accepted"
	StatusRejected      Status = "rejected"
	StatusInForce       Status = "in_force"
	StatusOutOfForce    Status = "out_of_force"
	StatusCancelled     Status = "cancelled"
	StatusInProgress    Status = "in_progress"
	StatusExpired       Status = "expired"
	StatusLapsed        Status = "lapsed"
	StatusPendingReview Status = "pending_review"
	StatusNotTakenUp    Status = "not_taken_up"
	StatusQuoted        Status = "quoted"
	StatusInEffect      Status = "in_effect"
	StatusNotInEffect   Status = "not_in_effect"
	StatusActive        Status = "active"
)

type DistributionChannel string

const (
	ChannelBroker    DistributionChannel = "broker"
	ChannelDirect    DistributionChannel = "direct"
	ChannelBinder    DistributionChannel = "binder"
	ChannelTiedAgent DistributionChannel = "tied_agent"
)

// CalculationJob represents a queued quote calculation request.
type CalculationJob struct {
	ID          int        `json:"id" gorm:"primaryKey"`
	QuoteID     int        `json:"quote_id" gorm:"index"`
	Basis       string     `json:"basis"`
	Credibility float64    `json:"credibility"`
	UserEmail   string     `json:"user_email"`
	UserName    string     `json:"user_name"`
	Status      string     `json:"status" gorm:"index;default:queued"` // queued, processing, completed, failed
	Error       string     `json:"error"`
	QueuedAt    time.Time  `json:"queued_at" gorm:"autoCreateTime"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

// GroupPricingQuote struct
type GroupPricingQuote struct {
	ID                           int                       `json:"id" gorm:"primary_key"`
	QuoteName                    string                    `json:"quote_name"`
	Basis                        string                    `json:"basis"`
	CreationDate                 time.Time                 `json:"creation_date" gorm:"type:datetime"`
	QuoteType                    string                    `json:"quote_type"`
	SchemeName                   string                    `json:"scheme_name"`
	SchemeID                     int                       `json:"scheme_id"`
	SchemeContact                string                    `json:"scheme_contact"`
	SchemeEmail                  string                    `json:"scheme_email"`
	QuoteBroker                  QuoteBroker               `json:"quote_broker" gorm:"embedded;embeddedPrefix:broker_"`
	DistributionChannel          DistributionChannel       `json:"distribution_channel" gorm:"size:20;default:'broker'"`
	ObligationType               string                    `json:"obligation_type"`
	CommencementDate             time.Time                 `json:"commencement_date"`
	CoverEndDate                 time.Time                 `json:"cover_end_date"`
	Industry                     string                    `json:"industry"`
	OccupationClass              int                       `json:"occupation_class"`
	FreeCoverLimit               float64                   `json:"free_cover_limit"`
	Currency                     string                    `json:"currency"`
	ExchangeRate                 int                       `json:"exchangeRate"`
	NormalRetirementAge          int                       `json:"normal_retirement_age"`
	ExperienceRating             string                    `json:"experience_rating"`
	CreatedBy                    string                    `json:"created_by"`
	Reviewer                     string                    `json:"reviewer"`
	ApprovedBy                   string                    `json:"approved_by"`
	SentBy                       string                    `json:"sent_by"`
	ModifiedBy                   string                    `json:"modified_by"`
	ModificationDate             time.Time                 `json:"modification_date"`
	Status                       Status                    `json:"status"`
	MemberDataCount              int                       `json:"member_data_count"`
	ClaimsExperienceCount        int                       `json:"claims_experience_count"`
	MemberRatingResultCount      int                       `json:"member_rating_result_count"`
	MemberPremiumScheduleCount   int                       `json:"member_premium_schedule_count"`
	BordereauxCount              int                       `json:"bordereaux_count"`
	UseGlobalSalaryMultiple      bool                      `json:"use_global_salary_multiple"`
	ContinuationOption           bool                      `json:"continuation_option"`
	SelectedSchemeCategories     StringArray               `json:"selected_scheme_categories" gorm:"type:json"`
	SchemeCategories             []SchemeCategory          `json:"scheme_categories" gorm:"foreignKey:QuoteId"`
	Loadings                     Loadings                  `json:"loadings" gorm:"embedded;embeddedPrefix:loadings_"`
	MemberAverageAge             int                       `json:"member_average_age" gorm:"column:member_average_age"`
	MemberAverageIncome          float64                   `json:"member_average_income" gorm:"column:member_average_income"`
	MemberMaleFemaleDistribution float64                   `json:"member_male_female_distribution" gorm:"column:member_male_female_distribution"`
	MemberIndicativeData         bool                      `json:"member_indicative_data"`
	MemberIndicativeDataSet      []MemberIndicativeDataSet `json:"member_indicative_data_set" gorm:"-"`
	RiskRateCode                 string                    `json:"risk_rate_code" gorm:"column:risk_rate_code"`
	SchemeQuoteStatus            Status                    `json:"scheme_quote_status" gorm:"column:scheme_quote_status"`
	EditMode                     bool                      `json:"edit_mode" gorm:"-"`
}

type GroupRiskQuoteStats struct {
	// Use GORM v2 tags to ensure auto-incrementing primary key
	ID                   int       `json:"id" gorm:"primaryKey;autoIncrement"`
	QuoteID              int       `json:"quote_id"`
	MemberCount          int       `json:"member_count"`
	AnnualPremium        float64   `json:"annual_premium"`
	Commission           float64   `json:"commission"`
	ExpectedExpenses     float64   `json:"expected_expenses"`
	ExpectedClaims       float64   `json:"expected_claims"`
	ExpectedGlaClaims    float64   `json:"expected_gla_claims"`
	ExpectedPtdClaims    float64   `json:"expected_ptd_claims"`
	ExpectedCiClaims     float64   `json:"expected_ci_claims"`
	ExpectedSglaClaims   float64   `json:"expected_sgla_claims"`
	ExpectedTtdClaims    float64   `json:"expected_ttd_claims"`
	ExpectedPhiClaims    float64   `json:"expected_phi_claims"`
	ExpectedFunClaims    float64   `json:"expected_fun_claims"`
	GlaAnnualPremium     float64   `json:"gla_annual_premium"`
	PtdAnnualPremium     float64   `json:"ptd_annual_premium"`
	CiAnnualPremium      float64   `json:"ci_annual_premium"`
	SglaAnnualPremium    float64   `json:"sgla_annual_premium"`
	TtdAnnualPremium     float64   `json:"ttd_annual_premium"`
	PhiAnnualPremium     float64   `json:"phi_annual_premium"`
	FuneralAnnualPremium float64   `json:"funeral_annual_premium"`
	ExpectedClaimsRatio  float64   `json:"expected_claims_ratio"`
	CoverStartDate       time.Time `json:"cover_start_date"`
	CoverEndDate         time.Time `json:"cover_end_date"`
	CreationDate         time.Time `json:"creation_date"`
	Creator              string    `json:"creator"`
}

type StringArray []string

// Implement the driver.Valuer interface
func (s StringArray) Value() (driver.Value, error) {
	// Marshal to JSON
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	// Return as string, NOT []byte
	return string(b), nil
}

// Implement the sql.Scanner interface
func (s *StringArray) Scan(value interface{}) error {
	var bytes []byte

	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return fmt.Errorf("cannot convert %T to StringArray", value)
	}

	return json.Unmarshal(bytes, s)
}

type SchemeCategory struct {
	ID                                       int     `json:"id" gorm:"primary_key"`
	QuoteId                                  int     `json:"quote_id"`
	SchemeCategory                           string  `json:"scheme_category"`
	Basis                                    string  `json:"basis"`
	FreeCoverLimit                           float64 `json:"free_cover_limit"`
	PtdBenefit                               bool    `json:"ptd_benefit"`
	GlaBenefit                               bool    `json:"gla_benefit"`
	CiBenefit                                bool    `json:"ci_benefit"`
	SglaBenefit                              bool    `json:"sgla_benefit"`
	PhiBenefit                               bool    `json:"phi_benefit"`
	TtdBenefit                               bool    `json:"ttd_benefit"`
	FamilyFuneralBenefit                     bool    `json:"family_funeral_benefit"`
	GlaSalaryMultiple                        float64 `json:"gla_salary_multiple"`
	GlaTerminalIllnessBenefit                string  `json:"gla_terminal_illness_benefit"`
	GlaWaitingPeriod                         int     `json:"gla_waiting_period"`
	GlaEducatorBenefit                       string  `json:"gla_educator_benefit"`
	GlaBenefitType                           string  `json:"gla_benefit_type"`
	PtdRiskType                              string  `json:"ptd_risk_type"`
	PtdBenefitType                           string  `json:"ptd_benefit_type"`
	PtdSalaryMultiple                        float64 `json:"ptd_salary_multiple"`
	PtdDeferredPeriod                        int     `json:"ptd_deferred_period"`
	PtdDisabilityDefinition                  string  `json:"ptd_disability_definition"`
	PtdEducatorBenefit                       string  `json:"ptd_educator_benefit"`
	CiBenefitStructure                       string  `json:"ci_benefit_structure"`
	CiBenefitDefinition                      string  `json:"ci_benefit_definition"`
	CiMaxBenefit                             float64 `json:"ci_max_benefit"`
	CiCriticalIllnessSalaryMultiple          float64 `json:"ci_critical_illness_salary_multiple"`
	SglaSalaryMultiple                       float64 `json:"sgla_salary_multiple"`
	SglaMaxBenefit                           float64 `json:"sgla_max_benefit"`
	PhiRiskType                              string  `json:"phi_risk_type"`
	PhiMaximumBenefit                        float64 `json:"phi_maximum_benefit"`
	PhiIncomeReplacementPercentage           float64 `json:"phi_income_replacement_percentage"`
	PhiUseTieredIncomeReplacementRatio       bool    `json:"phi_use_tiered_income_replacement_ratio"`
	PhiTieredIncomeReplacementType           string  `json:"phi_tiered_income_replacement_type"`
	PhiPremiumWaiver                         string  `json:"phi_premium_waiver"`
	PhiMedicalAidPremiumWaiver               string  `json:"phi_medical_aid_premium_waiver"`
	PhiBenefitEscalation                     string  `json:"phi_benefit_escalation"`
	PhiMaxPremiumWaiver                      float64 `json:"phi_max_premium_waiver"`
	PhiWaitingPeriod                         int     `json:"phi_waiting_period"`
	PhiNumberMonthlyPayments                 int     `json:"phi_number_monthly_payments"`
	PhiDeferredPeriod                        int     `json:"phi_deferred_period"`
	PhiDisabilityDefinition                  string  `json:"phi_disability_definition"`
	PhiNormalRetirementAge                   int     `json:"phi_normal_retirement_age"`
	TtdRiskType                              string  `json:"ttd_risk_type"`
	TtdMaximumBenefit                        float64 `json:"ttd_maximum_benefit"`
	TtdIncomeReplacementPercentage           float64 `json:"ttd_income_replacement_percentage"`
	TtdUseTieredIncomeReplacementRatio       bool    `json:"ttd_use_tiered_income_replacement_ratio"`
	TtdTieredIncomeReplacementType           string  `json:"ttd_tiered_income_replacement_type"`
	TtdPremiumWaiverPercentage               float64 `json:"ttd_premium_waiver_percentage"`
	TtdWaitingPeriod                         int     `json:"ttd_waiting_period"`
	TtdNumberMonthlyPayments                 float64 `json:"ttd_number_monthly_payments"`
	TtdDeferredPeriod                        int     `json:"ttd_deferred_period"`
	TtdDisabilityDefinition                  string  `json:"ttd_disability_definition"`
	FamilyFuneralMainMemberFuneralSumAssured float64 `json:"family_funeral_main_member_funeral_sum_assured"`
	FamilyFuneralSpouseFuneralSumAssured     float64 `json:"family_funeral_spouse_funeral_sum_assured"`
	FamilyFuneralChildrenFuneralSumAssured   float64 `json:"family_funeral_children_funeral_sum_assured"`
	FamilyFuneralAdultDependantSumAssured    float64 `json:"family_funeral_adult_dependant_sum_assured"`
	FamilyFuneralParentFuneralSumAssured     float64 `json:"family_funeral_parent_funeral_sum_assured"`
	FamilyFuneralMaxNumberChildren           int     `json:"family_funeral_max_number_children"`
	FamilyFuneralMaxNumberAdultDependants    int     `json:"family_funeral_max_number_adult_dependants"`
	PtdAlias                                 string  `json:"ptd_alias"`
	CiAlias                                  string  `json:"ci_alias"`
	SglaAlias                                string  `json:"sgla_alias"`
	PhiAlias                                 string  `json:"phi_alias"`
	TtdAlias                                 string  `json:"ttd_alias"`
	GlaAlias                                 string  `json:"gla_alias"`
	FamilyFuneralAlias                       string  `json:"family_funeral_alias"`
	Region                                   string  `json:"region"`
}

type QuoteBroker struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Gla struct {
	SalaryMultiple         float64 `json:"salary_multiple"`
	TerminalIllnessBenefit string  `json:"terminal_illness_benefit"`
	WaitingPeriod          int     `json:"waiting_period"`
	CoverTerminationAge    int     `json:"cover_termination_age"`
	EducatorBenefit        string  `json:"educator_benefit"`
}
type Ptd struct {
	RiskType             string  `json:"risk_type"`
	SalaryMultiple       float64 `json:"salary_multiple"`
	BenefitType          string  `json:"benefit_type"`
	CoverTerminationAge  int     `json:"cover_termination_age"`
	DeferredPeriod       int     `json:"deferred_period"`
	WaitingPeriod        int     `json:"waiting_period"`
	DisabilityDefinition string  `json:"disability_definition"`
	EducatorBenefit      string  `json:"educator_benefit"`
}
type Ci struct {
	CriticalIllnessSalaryMultiple float64 `json:"critical_illness_salary_multiple"`
	MaxBenefit                    float64 `json:"max_benefit"`
	BenefitStructure              string  `json:"benefit_structure"`
	BenefitDefinition             string  `json:"benefit_definition"`
	CoverTerminationAge           int     `json:"cover_termination_age"`
}
type Sgla struct {
	SglaSalaryMultiple  float64 `json:"sgla_salary_multiple"`
	MaxBenefit          float64 `json:"max_benefit"`
	CoverTerminationAge int     `json:"cover_termination_age"`
}
type Phi struct {
	Benefit                     string  `json:"benefit"`
	RiskType                    string  `json:"risk_type"`
	MaximumBenefit              float64 `json:"maximum_benefit"`
	IncomeReplacementPercentage float64 `json:"income_replacement_percentage"`
	PremiumWaiver               string  `json:"premium_waiver"`
	PremiumWaiverPercentage     float64 `json:"premium_waiver_percentage"`
	BenefitEscalationOption     string  `json:"benefit_escalation_option"`
	MaxPremiumWaiver            float64 `json:"max_premium_waiver"`
	//NumberMonthlyPayments       int     `json:"number_monthly_payments"`
	CoverTerminationAge     int    `json:"cover_termination_age"`
	WaitingPeriod           int    `json:"waiting_period"`
	DeferredPeriod          int    `json:"deferred_period"`
	DisabilityDefinition    string `json:"disability_definition"`
	MedicalAidPremiumWaiver string `json:"medical_aid_premium_waiver"`
	NormalRetirementAge     int    `json:"normal_retirement_age"`
}

type Ttd struct {
	Benefit                     string  `json:"benefit"`
	RiskType                    string  `json:"risk_type"`
	MaximumBenefit              float64 `json:"maximum_benefit"`
	IncomeReplacementPercentage float64 `json:"income_replacement_percentage"`
	//PremiumWaiver               string  `json:"premium_waiver"`
	//PremiumWaiverPercentage     float64 `json:"premium_waiver_percentage"`
	//EscalationPercentage        string  `json:"escalation_percentage"`
	//MaxPremiumWaiver            float64 `json:"max_premium_waiver"`
	NumberMonthlyPayments int    `json:"number_monthly_payments"`
	CoverTerminationAge   int    `json:"cover_termination_age"`
	WaitingPeriod         int    `json:"waiting_period"`
	DeferredPeriod        int    `json:"deferred_period"`
	DisabilityDefinition  string `json:"disability_definition"`
}
type FamilyFuneral struct {
	MainMemberFuneralSumAssured float64 `json:"main_member_funeral_sum_assured"`
	SpouseFuneralSumAssured     float64 `json:"spouse_funeral_sum_assured"`
	ChildrenFuneralSumAssured   float64 `json:"children_funeral_sum_assured"`
	AdultDependantSumAssured    float64 `json:"adult_dependant_sum_assured"`
	ParentFuneralSumAssured     float64 `json:"parent_funeral_sum_assured"`
	MaxNumberChildren           int     `json:"max_number_children"`
	MaxNumberAdultDependants    int     `json:"max_number_adult_dependants"`
}
type Loadings struct {
	CommissionLoading  float64 `json:"commission_loading"`
	ProfitLoading      float64 `json:"profit_loading"`
	ExpenseLoading     float64 `json:"expense_loading"`
	AdminLoading       float64 `json:"admin_loading"`
	ContingencyLoading float64 `json:"contingency_loading"`
	OtherLoading       float64 `json:"other_loading"`
	Discount           float64 `json:"discount"`
}

type Broker struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	ContactEmail       string    `json:"contact_email" gorm:"size:255;uniqueIndex"`
	ContactNumber      string    `json:"contact_number"`
	FSPNumber          string    `json:"fsp_number" gorm:"size:50"`
	FSPCategory        string    `json:"fsp_category" gorm:"size:50"`
	BinderAgreementRef string    `json:"binder_agreement_ref" gorm:"size:100"`
	TiedAgentRef       string    `json:"tied_agent_ref" gorm:"size:100"`
	CreationDate       time.Time `json:"creation_date" gorm:"autoCreateTime"`
	CreatedBy          string    `json:"created_by"`
}

type Reinsurer struct {
	ID                 int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name               string     `json:"name" gorm:"size:255;not null"`
	Code               string     `json:"code" gorm:"size:100;uniqueIndex;not null"`
	ContactEmail       string     `json:"contact_email" gorm:"size:255"`
	ContactPerson      string     `json:"contact_person" gorm:"size:255"`
	CreationDate       time.Time  `json:"creation_date" gorm:"autoCreateTime"`
	CreatedBy          string     `json:"created_by" gorm:"size:100"`
	IsActive           bool       `json:"is_active" gorm:"default:true;not null"`
	DeactivatedAt      *time.Time `json:"deactivated_at,omitempty"`
	DeactivatedBy      string     `json:"deactivated_by,omitempty" gorm:"size:100"`
	DeactivationReason string     `json:"deactivation_reason,omitempty" gorm:"type:text"`
}

type GPricingMemberData struct {
	ID                int            `json:"id" gorm:"primary_key"`
	Year              int            `json:"year" csv:"year"`
	SchemeName        string         `json:"scheme_name" csv:"scheme_name"`
	MemberName        string         `json:"member_name" csv:"member_name"`
	MemberIdNumber    string         `json:"member_id_number" csv:"member_id_number"`
	MemberIdType      string         `json:"member_id_type" csv:"member_id_type"`
	SchemeCategory    string         `json:"scheme_category" csv:"scheme_category"`
	Gender            string         `json:"gender" csv:"gender"`
	DateOfBirth       time.Time      `json:"date_of_birth" csv:"date_of_birth"`
	AnnualSalary      float64        `json:"annual_salary" csv:"annual_salary"`
	AddressLine1      string         `json:"address_line_1"`
	AddressLine2      string         `json:"address_line_2"`
	City              string         `json:"city"`
	Province          string         `json:"province"`
	PostalCode        string         `json:"postal_code"`
	PhoneNumber       string         `json:"phone_number"`
	Email             string         `json:"email"`
	EmployeeNumber    string         `json:"employee_number"`
	Occupation        string         `json:"occupation"`
	OccupationalClass string         `json:"occupational_class"`
	Benefits          MemberBenefits `json:"benefits" gorm:"embedded;embeddedPrefix:benefits_"`
	//GlaSalaryMultiple            float64        `json:"gla_salary_multiple" csv:"gla_salary_multiple"`
	//SglaSalaryMultiple           float64        `json:"sgla_salary_multiple" csv:"sgla_salary_multiple"`
	//PtdSalaryMultiple            float64        `json:"ptd_salary_multiple" csv:"ptd_salary_multiple"`
	//CiSalaryMultiple             float64        `json:"ci_salary_multiple" csv:"ci_salary_multiple"`
	ContributionWaiverProportion float64    `json:"contribution_waiver_proportion" csv:"contribution_waiver_proportion"`
	CreationDate                 time.Time  `json:"creation_date" gorm:"autoCreateTime"`
	EntryDate                    time.Time  `json:"entry_date" csv:"entry_date" gorm:"autoCreateTime"`
	ExitDate                     *time.Time `json:"exit_date" csv:"exit_date"`
	EffectiveExitDate            *time.Time `json:"effective_exit_date" csv:"effective_exit_date"`
	CreatedBy                    string     `json:"created_by"`
	QuoteId                      int        `json:"quote_id" gorm:"index"`
	SchemeId                     int        `json:"scheme_id"`
	Status                       string     `json:"status"`
	IsOriginalMember             bool       `json:"is_original_member"`
}

type GPricingMemberDataInForce struct {
	ID                    int            `json:"id" gorm:"primary_key"`
	Year                  int            `json:"year" csv:"year"`
	SchemeName            string         `json:"scheme_name" csv:"scheme_name"`
	MemberName            string         `json:"member_name" csv:"member_name"`
	MemberIdNumber        string         `json:"member_id_number" csv:"member_id_number"`
	MemberIdType          string         `json:"member_id_type" csv:"member_id_type"`
	SchemeCategory        string         `json:"scheme_category" csv:"scheme_category"`
	SchemeCategoryDetails SchemeCategory `json:"scheme_category_details,omitempty" gorm:"-"`
	Gender                string         `json:"gender" csv:"gender"`
	DateOfBirth           time.Time      `json:"date_of_birth" csv:"date_of_birth"`
	AnnualSalary          float64        `json:"annual_salary" csv:"annual_salary"`
	AddressLine1          string         `json:"address_line_1"`
	AddressLine2          string         `json:"address_line_2"`
	City                  string         `json:"city"`
	Province              string         `json:"province"`
	PostalCode            string         `json:"postal_code"`
	PhoneNumber           string         `json:"phone_number"`
	Email                 string         `json:"email"`
	EmployeeNumber        string         `json:"employee_number"`
	Occupation            string         `json:"occupation"`
	OccupationalClass     string         `json:"occupational_class"`
	Benefits              MemberBenefits `json:"benefits" csv:"benefits" gorm:"embedded;embeddedPrefix:benefits_"`
	//GlaSalaryMultiple            float64        `json:"gla_salary_multiple" csv:"gla_salary_multiple"`
	//SglaSalaryMultiple           float64        `json:"sgla_salary_multiple" csv:"sgla_salary_multiple"`
	//PtdSalaryMultiple            float64        `json:"ptd_salary_multiple" csv:"ptd_salary_multiple"`
	//CiSalaryMultiple             float64        `json:"ci_salary_multiple" csv:"ci_salary_multiple"`
	ContributionWaiverProportion float64    `json:"contribution_waiver_proportion,string" csv:"contribution_waiver_proportion"`
	CreationDate                 time.Time  `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	EntryDate                    time.Time  `json:"entry_date" csv:"entry_date" gorm:"autoCreateTime"`
	ExitDate                     *time.Time `json:"exit_date" csv:"exit_date"`
	EffectiveExitDate            *time.Time `json:"effective_exit_date" csv:"effective_exit_date"`
	CreatedBy                    string     `json:"created_by" csv:"created_by"`
	QuoteId                      int        `json:"quote_id" csv:"quote_id" gorm:"index"`
	SchemeId                     int        `json:"scheme_id" csv:"scheme_id"`
	Status                       string     `json:"status" csv:"status"`
	IsOriginalMember             bool       `json:"is_original_member" csv:"is_original_member"`
}

// UnmarshalJSON implements custom parsing to accept both full RFC3339 datetimes
// and simple date strings (YYYY-MM-DD) for time fields. This makes the API input
// tolerant of date-only values coming from clients.
func (g *GPricingMemberDataInForce) UnmarshalJSON(data []byte) error {
	type Alias GPricingMemberDataInForce
	// Use shadow fields for date/time inputs as strings so we can parse flexibly
	var aux struct {
		Alias
		DateOfBirth       *string `json:"date_of_birth"`
		CreationDate      *string `json:"creation_date"`
		EntryDate         *string `json:"entry_date"`
		ExitDate          *string `json:"exit_date"`
		EffectiveExitDate *string `json:"effective_exit_date"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Copy non-time fields first
	*g = GPricingMemberDataInForce(aux.Alias)

	// Helper to parse multiple layouts
	parse := func(s *string) (time.Time, error) {
		if s == nil {
			return time.Time{}, nil
		}
		str := *s
		if str == "" {
			return time.Time{}, nil
		}
		layouts := []string{
			time.RFC3339Nano,
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02",
		}
		var lastErr error
		for _, layout := range layouts {
			if t, err := time.Parse(layout, str); err == nil {
				return t, nil
			} else {
				lastErr = err
			}
		}
		return time.Time{}, lastErr
	}

	if t, err := parse(aux.DateOfBirth); err == nil {
		if !t.IsZero() {
			g.DateOfBirth = t
		}
	} else {
		return fmt.Errorf("invalid date_of_birth: %w", err)
	}

	if t, err := parse(aux.CreationDate); err == nil {
		if !t.IsZero() {
			g.CreationDate = t
		}
	} else {
		return fmt.Errorf("invalid creation_date: %w", err)
	}

	if t, err := parse(aux.EntryDate); err == nil {
		if !t.IsZero() {
			g.EntryDate = t
		}
	} else {
		return fmt.Errorf("invalid entry_date: %w", err)
	}

	if t, err := parse(aux.ExitDate); err == nil {
		if !t.IsZero() {
			g.ExitDate = &t
		}
	} else {
		return fmt.Errorf("invalid exit_date: %w", err)
	}

	if t, err := parse(aux.EffectiveExitDate); err == nil {
		if !t.IsZero() {
			g.EffectiveExitDate = &t
		}
	} else {
		return fmt.Errorf("invalid effective_exit_date: %w", err)
	}

	return nil
}

type MemberBenefits struct {
	GlaEnabled   bool    `json:"gla_enabled"`
	GlaMultiple  float64 `json:"gla_multiple"`
	SglaEnabled  bool    `json:"sgla_enabled"`
	SglaMultiple float64 `json:"sgla_multiple"`
	PtdEnabled   bool    `json:"ptd_enabled"`
	PtdMultiple  float64 `json:"ptd_multiple"`
	CiEnabled    bool    `json:"ci_enabled"`
	CiMultiple   float64 `json:"ci_multiple"`
	TtdEnabled   bool    `json:"ttd_enabled"`
	TtdMultiple  float64 `json:"ttd_multiple"`
	PhiEnabled   bool    `json:"phi_enabled"`
	PhiMultiple  float64 `json:"phi_multiple"`
	GffEnabled   bool    `json:"gff_enabled"`
}

type MemberAddress struct {
	ID                          int    `json:"id" gorm:"primary_key"`
	GPricingMemberDataInForceID int    `json:"g_pricing_member_data_in_force_id" gorm:"index"`
	AddressLine1                string `json:"address_line_1"`
	AddressLine2                string `json:"address_line_2"`
	City                        string `json:"city"`
	Province                    string `json:"province"`
	PostalCode                  string `json:"postal_code"`
}

type GroupPricingClaimsExperience struct {
	ID                 int       `json:"id" gorm:"primary_key"`
	Year               int       `json:"year" csv:"year"`
	SchemeName         string    `json:"scheme_name" csv:"-"`
	StartDate          string    `json:"start_date" csv:"start_date"`
	EndDate            string    `json:"end_date" csv:"end_date"`
	TotalGlaSumAssured float64   `json:"total_gla_sum_assured" csv:"total_gla_sum_assured"`
	TotalPtdSumAssured float64   `json:"total_ptd_sum_assured" csv:"total_ptd_sum_assured"`
	TotalCiSumAssured  float64   `json:"total_ci_sum_assured" csv:"total_ci_sum_assured"`
	NumberOfMembers    int       `json:"number_of_members" csv:"number_of_members"`
	NumberOfGlaClaims  int       `json:"number_of_gla_claims" csv:"number_of_gla_claims"`
	GlaClaimsAmount    float64   `json:"gla_claims_amount" csv:"gla_claims_amount"`
	NumberOfPtdClaims  int       `json:"number_of_ptd_claims" csv:"number_of_ptd_claims"`
	PtdClaimsAmount    float64   `json:"ptd_claims_amount" csv:"ptd_claims_amount"`
	NumberOfCiClaims   int       `json:"number_of_ci_claims" csv:"number_of_ci_claims"`
	CiClaimsAmount     float64   `json:"ci_claims_amount" csv:"ci_claims_amount"`
	NumberOfPhiClaims  int       `json:"number_of_phi_claims" csv:"number_of_phi_claims"`
	PhiClaimsAmount    float64   `json:"phi_claims_amount" csv:"phi_claims_amount"`
	Weighting          float64   `json:"weighting" csv:"weighting"`
	CreationDate       time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy          string    `json:"created_by" csv:"created_by"`
	QuoteId            int       `json:"quote_id" csv:"-"`
}

type MemberRatingResult struct {
	FinancialYear                   int       `json:"financial_year" csv:"financial_year"`
	SchemeId                        int       `json:"-" csv:"scheme_id"`
	QuoteId                         int       `json:"-" csv:"quote_id" gorm:"index"`
	Category                        string    `json:"category" csv:"category"`
	MemberName                      string    `json:"member_name" csv:"member_name"`
	MemberCount                     int       `json:"member_count" csv:"member_count"`
	Gender                          string    `json:"gender" csv:"gender"`
	DateOfBirth                     time.Time `json:"date_of_birth" csv:"date_of_birth"`
	IsOriginalMember                bool      `json:"is_original_member" csv:"is_original_member"`
	EntryDate                       time.Time `json:"entry_date" csv:"entry_date" gorm:"autoCreateTime"`
	ExitDate                        time.Time `json:"exit_date" csv:"exit_date" gorm:"autoCreateTime"`
	ExpCredibility                  float64   `json:"exp_credibility" csv:"exp_credibility"`
	ManuallyAddedCredibility        float64   `json:"manually_added_credibility" csv:"manually_added_credibility"`
	AnnualSalary                    float64   `json:"annual_salary" csv:"annual_salary"`
	IncomeLevel                     int       `json:"income_level" csv:"income_level"`
	GlaSalaryMultiple               float64   `json:"gla_salary_multiple" csv:"gla_salary_multiple"`
	SglaSalaryMultiple              float64   `json:"sgla_salary_multiple" csv:"sgla_salary_multiple"`
	PtdSalaryMultiple               float64   `json:"ptd_salary_multiple" csv:"ptd_salary_multiple"`
	CiSalaryMultiple                float64   `json:"ci_salary_multiple" csv:"ci_salary_multiple"`
	Occupation                      string    `json:"occupation" csv:"occupation"`
	OccupationClass                 int       `json:"occupation_class" csv:"occupation_class"`
	Industry                        string    `json:"industry" csv:"industry"`
	AgeNextBirthday                 int       `json:"age_next_birthday" csv:"age_next_birthday"`
	AgeBand                         string    `json:"age_band" csv:"age_band"`
	SpouseGender                    string    `json:"spouse_gender" csv:"spouse_gender"`
	SpouseAgeNextBirthday           int       `json:"spouse_age_next_birthday" csv:"spouse_age_next_birthday"`
	AverageDependantAgeNextBirthday float64   `json:"average_dependant_age_next_birthday" csv:"average_dependant_age_next_birthday"`
	AverageChildAgeNextBirthday     float64   `json:"average_child_age_next_birthday" csv:"average_child_age_next_birthday"`
	AverageNumberDependants         float64   `json:"average_number_dependants" csv:"average_number_dependants"`
	AverageNumberChildren           float64   `json:"average_number_children" csv:"average_number_children"`
	CalculatedFreeCoverLimit        float64   `json:"calculated_free_cover_limit" csv:"calculated_free_cover_limit"`
	AppliedFreeCoverLimit           float64   `json:"applied_free_cover_limit" csv:"applied_free_cover_limit"`
	GlaSumAssured                   float64   `json:"gla_sum_assured" csv:"gla_sum_assured"`
	GlaCappedSumAssured             float64   `json:"gla_capped_sum_assured" csv:"gla_capped_sum_assured"`
	ExpenseLoading                  float64   `json:"expense_loading" csv:"expense_loading"`
	AdminLoading                    float64   `json:"admin_loading" csv:"admin_loading"`
	CommissionLoading               float64   `json:"commission_loading" csv:"commission_loading"`
	ProfitLoading                   float64   `json:"profit_loading" csv:"profit_loading"`
	OtherLoading                    float64   `json:"other_loading" csv:"other_loading"`
	Discount                        float64   `json:"discount" csv:"discount"`
	TotalLoading                    float64   `json:"total_loading" csv:"total_loading"`

	BaseGlaRate                    float64 `json:"base_gla_rate" csv:"base_gla_rate"`
	GlaLoading                     float64 `json:"gla_loading" csv:"gla_loading"`
	GlaTerminalIllnessLoading      float64 `json:"gla_terminal_illness_loading" csv:"gla_terminal_illness_loading"`
	LoadedGlaRate                  float64 `json:"loaded_gla_rate" csv:"loaded_gla_rate"`
	GlaWeightedExperienceCrudeRate float64 `json:"gla_weighted_experience_crude_rate" csv:"gla_weighted_experience_crude_rate"`
	GlaTheoreticalRate             float64 `json:"gla_theoretical_rate" csv:"gla_theoretical_rate"`
	PtdExperienceCrudeRate         float64 `json:"ptd_experience_crude_rate" csv:"ptd_experience_crude_rate"`
	PtdTheoreticalRate             float64 `json:"ptd_theoretical_rate" csv:"ptd_theoretical_rate"`
	CiExperienceCrudeRate          float64 `json:"ci_experience_crude_rate" csv:"ci_experience_crude_rate"`
	CiTheoreticalRate              float64 `json:"ci_theoretical_rate" csv:"ci_theoretical_rate"`
	ExpAdjLoadedGlaRate            float64 `json:"exp_adj_loaded_gla_rate" csv:"exp_adj_loaded_gla_rate"`
	GlaExperienceAdjustment        float64 `json:"gla_experience_adjustment" csv:"gla_experience_adjustment"`
	GlaRiskPremium                 float64 `json:"gla_risk_premium" csv:"gla_risk_premium"`
	ExpAdjGlaRiskPremium           float64 `json:"exp_adj_gla_risk_premium" csv:"exp_adj_gla_risk_premium"`
	GlaOfficePremium               float64 `json:"gla_office_premium" csv:"gla_office_premium"`
	ExpAdjGlaOfficePremium         float64 `json:"exp_adj_gla_office_premium" csv:"exp_adj_gla_office_premium"`

	PtdSumAssured           float64 `json:"ptd_sum_assured" csv:"ptd_sum_assured"`
	PtdCappedSumAssured     float64 `json:"ptd_capped_sum_assured" csv:"ptd_capped_sum_assured"`
	BasePtdRate             float64 `json:"base_ptd_rate" csv:"base_ptd_rate"`
	PtdLoading              float64 `json:"ptd_loading" csv:"ptd_loading"`
	LoadedPtdRate           float64 `json:"loaded_ptd_rate" csv:"loaded_ptd_rate"`
	PtdExperienceAdjustment float64 `json:"ptd_experience_adjustment" csv:"ptd_experience_adjustment"`
	ExpAdjLoadedPtdRate     float64 `json:"exp_adj_loaded_ptd_rate" csv:"exp_adj_loaded_ptd_rate"`
	PtdRiskPremium          float64 `json:"ptd_risk_premium" csv:"ptd_risk_premium"`
	ExpAdjPtdRiskPremium    float64 `json:"exp_adj_ptd_risk_premium" csv:"exp_adj_ptd_risk_premium"`
	PtdOfficePremium        float64 `json:"ptd_office_premium" csv:"ptd_office_premium"`
	ExpAdjPtdOfficePremium  float64 `json:"exp_adj_ptd_office_premium" csv:"exp_adj_ptd_office_premium"`

	CiSumAssured           float64 `json:"ci_sum_assured" csv:"ci_sum_assured"`
	CiCappedSumAssured     float64 `json:"ci_capped_sum_assured" csv:"ci_capped_sum_assured"`
	BaseCiRate             float64 `json:"base_ci_rate" csv:"base_ci_rate"`
	CiLoading              float64 `json:"ci_loading" csv:"ci_loading"`
	LoadedCiRate           float64 `json:"loaded_ci_rate" csv:"loaded_ci_rate"`
	CiExperienceAdjustment float64 `json:"ci_experience_adjustment" csv:"ci_experience_adjustment"`
	ExpAdjLoadedCiRate     float64 `json:"exp_adj_loaded_ci_rate" csv:"exp_adj_loaded_ci_rate"`
	CiRiskPremium          float64 `json:"ci_risk_premium" csv:"ci_risk_premium"`
	ExpAdjCiRiskPremium    float64 `json:"exp_adj_ci_risk_premium" csv:"exp_adj_ci_risk_premium"`
	CiOfficePremium        float64 `json:"ci_office_premium" csv:"ci_office_premium"`
	ExpAdjCiOfficePremium  float64 `json:"exp_adj_ci_office_premium" csv:"exp_adj_ci_office_premium"`

	SpouseGlaSumAssured          float64 `json:"spouse_gla_sum_assured" csv:"spouse_gla_sum_assured"`
	SpouseGlaCappedSumAssured    float64 `json:"spouse_gla_capped_sum_assured" csv:"spouse_gla_capped_sum_assured"`
	BaseSpouseGlaRate            float64 `json:"base_spouse_gla_rate" csv:"base_spouse_gla_rate"`
	SpouseGlaLoading             float64 `json:"spouse_gla_loading" csv:"spouse_gla_loading"`
	LoadedSpouseGlaRate          float64 `json:"loaded_spouse_gla_rate" csv:"loaded_spouse_gla_rate"`
	ExpAdjLoadedSpouseGlaRate    float64 `json:"exp_adj_loaded_spouse_gla_rate" csv:"exp_adj_loaded_spouse_gla_rate"`
	SpouseGlaRiskPremium         float64 `json:"spouse_gla_risk_premium" csv:"spouse_gla_risk_premium"`
	SpouseGlaOfficePremium       float64 `json:"spouse_gla_office_premium" csv:"spouse_gla_office_premium"`
	ExpAdjSpouseGlaOfficePremium float64 `json:"exp_adj_spouse_gla_office_premium" csv:"exp_adj_spouse_gla_office_premium"`

	TtdIncome                  float64 `json:"ttd_income" csv:"ttd_income"`
	TtdCappedIncome            float64 `json:"ttd_capped_income" csv:"ttd_capped_income"`
	TtdNumberOfMonthlyPayments float64 `json:"ttd_number_of_monthly_payments" csv:"ttd_number_of_monthly_payments"`
	IncomeReplacementRatio     float64 `json:"income_replacement_ratio" csv:"income_replacement_ratio"`
	BaseTtdRate                float64 `json:"base_ttd_rate" csv:"base_ttd_rate"`
	TtdLoading                 float64 `json:"ttd_loading" csv:"ttd_loading"`
	LoadedTtdRate              float64 `json:"loaded_ttd_rate" csv:"loaded_ttd_rate"`
	TtdExperienceAdjustment    float64 `json:"ttd_experience_adjustment" csv:"ttd_experience_adjustment"`
	ExpAdjLoadedTtdRate        float64 `json:"exp_adj_loaded_ttd_rate" csv:"exp_adj_loaded_ttd_rate"`
	TtdRiskPremium             float64 `json:"ttd_risk_premium" csv:"ttd_risk_premium"`
	ExpAdjTtdRiskPremium       float64 `json:"exp_adj_ttd_risk_premium" csv:"exp_adj_ttd_risk_premium"`
	TtdOfficePremium           float64 `json:"ttd_office_premium" csv:"ttd_office_premium"`
	ExpAdjTtdOfficePremium     float64 `json:"exp_adj_ttd_office_premium" csv:"exp_adj_ttd_office_premium"`
	ExpAdjSpouseGlaRiskPremium float64 `json:"exp_adj_spouse_gla_risk_premium" csv:"exp_adj_spouse_gla_risk_premium"`

	PhiIncome               float64 `json:"phi_income" csv:"phi_income"`
	PhiCappedIncome         float64 `json:"phi_capped_income" csv:"phi_capped_income"`
	PhiContributionWaiver   float64 `json:"phi_contribution_waiver" csv:"phi_contribution_waiver"`
	PhiMedicalAidWaiver     float64 `json:"phi_medical_aid_waiver" csv:"phi_medical_aid_waiver"`
	PhiMonthlyBenefit       float64 `json:"phi_monthly_benefit" csv:"phi_monthly_benefit"`
	PhiAnnuityFactor        float64 `json:"phi_annuity_factor" csv:"phi_annuity_factor"`
	BasePhiRate             float64 `json:"base_phi_rate" csv:"base_phi_rate"`
	PhiSalaryLevel          float64 `json:"phi_salary_level" csv:"phi_salary_level"`
	PhiLoading              float64 `json:"phi_loading" csv:"phi_loading"`
	LoadedPhiRate           float64 `json:"loaded_phi_rate" csv:"loaded_phi_rate"`
	PhiExperienceAdjustment float64 `json:"phi_experience_adjustment" csv:"phi_experience_adjustment"`
	ExpAdjLoadedPhiRate     float64 `json:"exp_adj_loaded_phi_rate" csv:"exp_adj_loaded_phi_rate"`
	PhiRiskPremium          float64 `json:"phi_risk_premium" csv:"phi_risk_premium"`
	ExpAdjPhiRiskPremium    float64 `json:"exp_adj_phi_risk_premium" csv:"exp_adj_phi_risk_premium"`
	PhiOfficePremium        float64 `json:"phi_office_premium" csv:"phi_office_premium"`
	ExpAdjPhiOfficePremium  float64 `json:"exp_adj_phi_office_premium" csv:"exp_adj_phi_office_premium"`

	MemberFuneralSumAssured        float64 `json:"member_funeral_sum_assured" csv:"member_funeral_sum_assured"`
	MainMemberFuneralCost          float64 `json:"main_member_funeral_cost" csv:"main_member_funeral_cost"`
	MainMemberFuneralOfficePremium float64 `json:"main_member_funeral_office_premium" csv:"main_member_funeral_office_premium"`

	//MarriageProportion float64 `json:"marriage_proportion" csv:"marriage_proportion"`

	SpouseFuneralSumAssured    float64 `json:"spouse_funeral_sum_assured" csv:"spouse_funeral_sum_assured"`
	SpouseFuneralCost          float64 `json:"spouse_funeral_cost" csv:"spouse_funeral_cost"`
	SpouseFuneralOfficePremium float64 `json:"spouse_funeral_office_premium" csv:"spouse_funeral_office_premium"`

	ChildFuneralBaseRate         float64 `json:"child_funeral_base_rate" csv:"child_funeral_base_rate"`
	ChildFuneralSumAssured       float64 `json:"child_funeral_sum_assured" csv:"child_funeral_sum_assured"`
	ChildrenFuneralCost          float64 `json:"children_funeral_cost" csv:"children_funeral_cost"`
	ChildrenFuneralOfficePremium float64 `json:"children_funeral_office_premium" csv:"children_funeral_office_premium"`

	DependantFuneralBaseRate       float64 `json:"dependant_funeral_base_rate" csv:"dependant_funeral_base_rate"`
	DependantFuneralSumAssured     float64 `json:"dependant_funeral_sum_assured" csv:"dependant_funeral_sum_assured"`
	DependantsFuneralCost          float64 `json:"dependants_funeral_cost" csv:"dependants_funeral_cost"`
	DependantsFuneralOfficePremium float64 `json:"dependants_funeral_office_premium" csv:"dependants_funeral_office_premium"`

	ParentFuneralSumAssured      float64 `json:"parent_funeral_sum_assured" csv:"parent_funeral_sum_assured"`
	TotalFuneralRiskCost         float64 `json:"total_funeral_risk_cost" csv:"total_funeral_risk_cost"`
	ExpAdjTotalFuneralRiskCost   float64 `json:"exp_adj_total_funeral_risk_cost" csv:"exp_adj_total_funeral_risk_cost"`
	TotalFuneralOfficeCost       float64 `json:"total_funeral_office_cost" csv:"total_funeral_office_cost"`
	ExpAdjTotalFuneralOfficeCost float64 `json:"exp_adj_total_funeral_office_cost" csv:"exp_adj_total_funeral_office_cost"`

	Grade0SumAssured            float64 `json:"grade_0_sum_assured" csv:"grade_0_sum_assured"`
	Grade17SumAssured           float64 `json:"grade_1_7_sum_assured" csv:"grade_1_7_sum_assured"`
	Grade812SumAssured          float64 `json:"grade_8_12_sum_assured" csv:"grade_8_12_sum_assured"`
	TertiarySumAssured          float64 `json:"tertiary_sum_assured" csv:"tertiary_sum_assured"`
	Grade0RiskRate              float64 `json:"grade_0_risk_rate" csv:"grade_0_risk_rate"`
	Grade17RiskRate             float64 `json:"grade_1_7_risk_rate" csv:"grade_1_7_risk_rate"`
	Grade812RiskRate            float64 `json:"grade_8_12_risk_rate" csv:"grade_8_12_risk_rate"`
	TertiaryRiskRate            float64 `json:"tertiary_risk_rate" csv:"tertiary_risk_rate"`
	EducatorRiskPremium         float64 `json:"educator_risk_premium" csv:"educator_risk_premium"`
	EducatorOfficePremium       float64 `json:"educator_office_premium" csv:"educator_office_premium"`
	ExpAdjEducatorRiskPremium   float64 `json:"exp_adj_educator_risk_premium" csv:"exp_adj_educator_risk_premium"`
	ExpAdjEducatorOfficePremium float64 `json:"exp_adj_educator_office_premium" csv:"exp_adj_educator_office_premium"`

	ExceedsNormalRetirementAgeIndicator    int       `json:"exceeds_normal_retirement_age_indicator" csv:"exceeds_normal_retirement_age_indicator"`
	ExceedsFreeCoverLimitIndicator         int       `json:"exceeds_free_cover_limit_indicator" csv:"exceeds_free_cover_limit_indicator"`
	FuneralExperienceAdjustedAnnualPremium float64   `json:"funeral_experience_adjusted_annual_premium" csv:"funeral_experience_adjusted_annual_premium"`
	CreationDate                           time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy                              string    `json:"created_by" csv:"created_by"`
}

type MemberRatingResultSummary struct {
	ID                                      int     `json:"id" gorm:"primary_key"`
	QuoteId                                 int     `json:"quote_id" csv:"quote_id" gorm:"index"`
	SchemeId                                int     `json:"scheme_id" csv:"scheme_id"`
	Category                                string  `json:"category" csv:"category"`
	FinancialYear                           int     `json:"financial_year" csv:"financial_year"`
	MemberCount                             float64 `json:"member_count" csv:"member_count"`
	TotalAnnualSalary                       float64 `json:"total_annual_salary" csv:"total_annual_salary"`
	IfStatus                                Status  `json:"if_status" csv:"if_status"`
	QuoteType                               string  `json:"quote_type" csv:"quote_type"`
	FreeCoverLimit                          float64 `json:"free_cover_limit" csv:"free_cover_limit"`
	ExpenseLoading                          float64 `json:"expense_loading" csv:"expense_loading"`
	CommissionLoading                       float64 `json:"commission_loading" csv:"commission_loading"`
	ProfitLoading                           float64 `json:"profit_loading" csv:"profit_loading"`
	AcceleratedBenefitDiscount              float64 `json:"accelerated_benefit_discount" csv:"accelerated_benefit_discount"`
	MinGlaSumAssured                        float64 `json:"min_gla_sum_assured" csv:"min_gla_sum_assured"`
	MaxGlaSumAssured                        float64 `json:"max_gla_sum_assured" csv:"max_gla_sum_assured"`
	MaxGlaCappedSumAssured                  float64 `json:"max_gla_capped_sum_assured" csv:"max_gla_capped_sum_assured"`
	TotalGlaSumAssured                      float64 `json:"total_gla_sum_assured" csv:"total_gla_sum_assured"`
	TotalGlaCappedSumAssured                float64 `json:"total_gla_capped_sum_assured" csv:"total_gla_capped_sum_assured"`
	AverageGlaCappedSumAssured              float64 `json:"average_gla_capped_sum_assured" csv:"average_gla_capped_sum_assured"`
	TotalGlaRiskRate                        float64 `json:"total_gla_risk_rate" csv:"total_gla_risk_rate"`
	TotalGlaAnnualRiskPremium               float64 `json:"total_gla_annual_risk_premium" csv:"total_gla_annual_risk_premium"`
	GlaRiskRatePer1000SA                    float64 `json:"gla_risk_rate_per_1000_sa" csv:"gla_risk_rate_per_1000_sa"`
	ProportionGlaAnnualRiskPremiumSalary    float64 `json:"proportion_gla_annual_risk_premium_salary" csv:"proportion_gla_annual_risk_premium_salary"`
	TotalGlaAnnualOfficePremium             float64 `json:"total_gla_annual_office_premium" csv:"total_gla_annual_office_premium"`
	GlaOfficeRatePer1000SA                  float64 `json:"gla_office_rate_per_1000_sa" csv:"gla_office_rate_per_1000_sa"`
	ProportionGlaOfficePremiumSalary        float64 `json:"proportion_gla_office_premium_salary" csv:"proportion_gla_office_premium_salary"`
	ExpTotalGlaRiskRate                     float64 `json:"exp_total_gla_risk_rate" csv:"exp_total_gla_risk_rate"`
	ExpTotalGlaAnnualRiskPremium            float64 `json:"exp_total_gla_annual_risk_premium" csv:"exp_total_gla_annual_risk_premium"`
	ExpGlaRiskRatePer1000SA                 float64 `json:"exp_gla_risk_rate_per_1000_sa" csv:"exp_gla_risk_rate_per_1000_sa"`
	ExpProportionGlaAnnualRiskPremiumSalary float64 `json:"exp_proportion_gla_annual_risk_premium_salary" csv:"exp_proportion_gla_annual_risk_premium_salary"`
	ExpTotalGlaAnnualOfficePremium          float64 `json:"exp_total_gla_annual_office_premium" csv:"exp_total_gla_annual_office_premium"`
	ExpGlaOfficeRatePer1000SA               float64 `json:"exp_gla_office_rate_per_1000_sa" csv:"exp_gla_office_rate_per_1000_sa"`
	ExpProportionGlaOfficePremiumSalary     float64 `json:"exp_proportion_gla_office_premium_salary" csv:"exp_proportion_gla_office_premium_salary"`

	MinPtdSumAssured                        float64 `json:"min_ptd_sum_assured" csv:"min_ptd_sum_assured"`
	MaxPtdSumAssured                        float64 `json:"max_ptd_sum_assured" csv:"max_ptd_sum_assured"`
	MaxPtdCappedSumAssured                  float64 `json:"max_ptd_capped_sum_assured" csv:"max_ptd_capped_sum_assured"`
	TotalPtdSumAssured                      float64 `json:"total_ptd_sum_assured" csv:"total_ptd_sum_assured"`
	TotalPtdCappedSumAssured                float64 `json:"total_ptd_capped_sum_assured" csv:"total_ptd_capped_sum_assured"`
	AveragePtdCappedSumAssured              float64 `json:"average_ptd_capped_sum_assured" csv:"average_ptd_capped_sum_assured"`
	TotalPtdRiskRate                        float64 `json:"total_ptd_risk_rate" csv:"total_ptd_risk_rate"`
	TotalPtdAnnualRiskPremium               float64 `json:"total_ptd_annual_risk_premium" csv:"total_ptd_annual_risk_premium"`
	PtdRiskRatePer1000SA                    float64 `json:"ptd_risk_rate_per_1000_sa" csv:"ptd_risk_rate_per_1000_sa"`
	ProportionPtdAnnualRiskPremiumSalary    float64 `json:"proportion_ptd_annual_risk_premium_salary" csv:"proportion_ptd_annual_risk_premium_salary"`
	TotalPtdAnnualOfficePremium             float64 `json:"total_ptd_annual_office_premium" csv:"total_ptd_annual_office_premium"`
	PtdOfficeRatePer1000SA                  float64 `json:"ptd_office_rate_per_1000_sa" csv:"ptd_office_rate_per_1000_sa"`
	ProportionPtdOfficePremiumSalary        float64 `json:"proportion_ptd_office_premium_salary" csv:"proportion_ptd_office_premium_salary"`
	ExpTotalPtdRiskRate                     float64 `json:"exp_total_ptd_risk_rate" csv:"exp_total_ptd_risk_rate"`
	ExpTotalPtdAnnualRiskPremium            float64 `json:"exp_total_ptd_annual_risk_premium" csv:"exp_total_ptd_annual_risk_premium"`
	ExpPtdRiskRatePer1000SA                 float64 `json:"exp_ptd_risk_rate_per_1000_sa" csv:"exp_ptd_risk_rate_per_1000_sa"`
	ExpProportionPtdAnnualRiskPremiumSalary float64 `json:"exp_proportion_ptd_annual_risk_premium_salary" csv:"exp_proportion_ptd_annual_risk_premium_salary"`
	ExpTotalPtdAnnualOfficePremium          float64 `json:"exp_total_ptd_annual_office_premium" csv:"exp_total_ptd_annual_office_premium"`
	ExpPtdOfficeRatePer1000SA               float64 `json:"exp_ptd_office_rate_per_1000_sa" csv:"exp_ptd_office_rate_per_1000_sa"`
	ExpProportionPtdOfficePremiumSalary     float64 `json:"exp_proportion_ptd_office_premium_salary" csv:"exp_proportion_ptd_office_premium_salary"`

	MinCiSumAssured                        float64 `json:"min_ci_sum_assured" csv:"min_ci_sum_assured"`
	MaxCiSumAssured                        float64 `json:"max_ci_sum_assured" csv:"max_ci_sum_assured"`
	MaxCiCappedSumAssured                  float64 `json:"max_ci_capped_sum_assured" csv:"max_ci_capped_sum_assured"`
	TotalCiSumAssured                      float64 `json:"total_ci_sum_assured" csv:"total_ci_sum_assured"`
	TotalCiCappedSumAssured                float64 `json:"total_ci_capped_sum_assured" csv:"total_ci_capped_sum_assured"`
	AverageCiCappedSumAssured              float64 `json:"average_ci_capped_sum_assured" csv:"average_ci_capped_sum_assured"`
	TotalCiRiskRate                        float64 `json:"total_ci_risk_rate" csv:"total_ci_risk_rate"`
	TotalCiAnnualRiskPremium               float64 `json:"total_ci_annual_risk_premium" csv:"total_ci_annual_risk_premium"`
	CiRiskRatePer1000SA                    float64 `json:"ci_risk_rate_per_1000_sa" csv:"ci_risk_rate_per_1000_sa"`
	ProportionCiAnnualRiskPremiumSalary    float64 `json:"proportion_ci_annual_risk_premium_salary" csv:"proportion_ci_annual_risk_premium_salary"`
	TotalCiAnnualOfficePremium             float64 `json:"total_ci_annual_office_premium" csv:"total_ci_annual_office_premium"`
	CiOfficeRatePer1000SA                  float64 `json:"ci_office_rate_per_1000_sa" csv:"ci_office_rate_per_1000_sa"`
	ProportionCiOfficePremiumSalary        float64 `json:"proportion_ci_office_premium_salary" csv:"proportion_ci_office_premium_salary"`
	ExpTotalCiRiskRate                     float64 `json:"exp_total_ci_risk_rate" csv:"exp_total_ci_risk_rate"`
	ExpTotalCiAnnualRiskPremium            float64 `json:"exp_total_ci_annual_risk_premium" csv:"exp_total_ci_annual_risk_premium"`
	ExpCiRiskRatePer1000SA                 float64 `json:"exp_ci_risk_rate_per_1000_sa" csv:"exp_ci_risk_rate_per_1000_sa"`
	ExpProportionCiAnnualRiskPremiumSalary float64 `json:"exp_proportion_ci_annual_risk_premium_salary" csv:"exp_proportion_ci_annual_risk_premium_salary"`
	ExpTotalCiAnnualOfficePremium          float64 `json:"exp_total_ci_annual_office_premium" csv:"exp_total_ci_annual_office_premium"`
	ExpCiOfficeRatePer1000SA               float64 `json:"exp_ci_office_rate_per_1000_sa" csv:"exp_ci_office_rate_per_1000_sa"`
	ExpProportionCiOfficePremiumSalary     float64 `json:"exp_proportion_ci_office_premium_salary" csv:"exp_proportion_ci_office_premium_salary"`

	MinSglaSumAssured                        float64 `json:"min_sgla_sum_assured" csv:"min_sgla_sum_assured"`
	MaxSglaSumAssured                        float64 `json:"max_sgla_sum_assured" csv:"max_sgla_sum_assured"`
	MaxSglaCappedSumAssured                  float64 `json:"max_sgla_capped_sum_assured" csv:"max_sgla_capped_sum_assured"`
	TotalSglaSumAssured                      float64 `json:"total_sgla_sum_assured" csv:"total_sgla_sum_assured"`
	TotalSglaCappedSumAssured                float64 `json:"total_sgla_capped_sum_assured" csv:"total_sgla_capped_sum_assured"`
	AverageSglaCappedSumAssured              float64 `json:"average_sgla_capped_sum_assured" csv:"average_sgla_capped_sum_assured"`
	TotalSglaRiskRate                        float64 `json:"total_sgla_risk_rate" csv:"total_sgla_risk_rate"`
	TotalSglaAnnualRiskPremium               float64 `json:"total_sgla_annual_risk_premium" csv:"total_sgla_annual_risk_premium"`
	SglaRiskRatePer1000SA                    float64 `json:"sgla_risk_rate_per_1000_sa" csv:"sgla_risk_rate_per_1000_sa"`
	ProportionSglaAnnualRiskPremiumSalary    float64 `json:"proportion_sgla_annual_risk_premium_salary" csv:"proportion_sgla_annual_risk_premium_salary"`
	TotalSglaAnnualOfficePremium             float64 `json:"total_sgla_annual_office_premium" csv:"total_sgla_annual_office_premium"`
	SglaOfficeRatePer1000SA                  float64 `json:"sgla_office_rate_per_1000_sa" csv:"sgla_office_rate_per_1000_sa"`
	ProportionSglaOfficePremiumSalary        float64 `json:"proportion_sgla_office_premium_salary" csv:"proportion_sgla_office_premium_salary"`
	ExpTotalSglaRiskRate                     float64 `json:"exp_total_sgla_risk_rate" csv:"exp_total_sgla_risk_rate"`
	ExpTotalSglaAnnualRiskPremium            float64 `json:"exp_total_sgla_annual_risk_premium" csv:"exp_total_sgla_annual_risk_premium"`
	ExpSglaRiskRatePer1000SA                 float64 `json:"exp_sgla_risk_rate_per_1000_sa" csv:"exp_sgla_risk_rate_per_1000_sa"`
	ExpProportionSglaAnnualRiskPremiumSalary float64 `json:"exp_proportion_sgla_annual_risk_premium_salary" csv:"exp_proportion_sgla_annual_risk_premium_salary"`
	ExpTotalSglaAnnualOfficePremium          float64 `json:"exp_total_sgla_annual_office_premium" csv:"exp_total_sgla_annual_office_premium"`
	ExpSglaOfficeRatePer1000SA               float64 `json:"exp_sgla_office_rate_per_1000_sa" csv:"exp_sgla_office_rate_per_1000_sa"`
	ExpProportionSglaOfficePremiumSalary     float64 `json:"exp_proportion_sgla_office_premium_salary" csv:"exp_proportion_sgla_office_premium_salary"`

	MinTtdIncome                            float64 `json:"min_ttd_income" csv:"min_ttd_income"`
	MaxTtdIncome                            float64 `json:"max_ttd_income" csv:"max_ttd_income"`
	MaxTtdCappedIncome                      float64 `json:"max_ttd_capped_income" csv:"max_ttd_capped_income"`
	TotalTtdIncome                          float64 `json:"total_ttd_income" csv:"total_ttd_income"`
	TotalTtdCappedIncome                    float64 `json:"total_ttd_capped_income" csv:"total_ttd_capped_income"`
	AverageTtdCappedIncome                  float64 `json:"average_ttd_capped_income" csv:"average_ttd_capped_income"`
	TotalTtdRiskRate                        float64 `json:"total_ttd_risk_rate" csv:"total_ttd_risk_rate"`
	TotalTtdAnnualRiskPremium               float64 `json:"total_ttd_annual_risk_premium" csv:"total_ttd_annual_risk_premium"`
	TtdRiskRatePer1000SA                    float64 `json:"ttd_risk_rate_per_1000_sa" csv:"ttd_risk_rate_per_1000_sa"`
	ProportionTtdAnnualRiskPremiumSalary    float64 `json:"proportion_ttd_annual_risk_premium_salary" csv:"proportion_ttd_annual_risk_premium_salary"`
	TotalTtdAnnualOfficePremium             float64 `json:"total_ttd_annual_office_premium" csv:"total_ttd_annual_office_premium"`
	TtdOfficeRatePer1000SA                  float64 `json:"ttd_office_rate_per_1000_sa" csv:"ttd_office_rate_per_1000_sa"`
	ProportionTtdOfficePremiumSalary        float64 `json:"proportion_ttd_office_premium_salary" csv:"proportion_ttd_office_premium_salary"`
	ExpTotalTtdRiskRate                     float64 `json:"exp_total_ttd_risk_rate" csv:"exp_total_ttd_risk_rate"`
	ExpTotalTtdAnnualRiskPremium            float64 `json:"exp_total_ttd_annual_risk_premium" csv:"exp_total_ttd_annual_risk_premium"`
	ExpTtdRiskRatePer1000SA                 float64 `json:"exp_ttd_risk_rate_per_1000_sa" csv:"exp_ttd_risk_rate_per_1000_sa"`
	ExpProportionTtdAnnualRiskPremiumSalary float64 `json:"exp_proportion_ttd_annual_risk_premium_salary" csv:"exp_proportion_ttd_annual_risk_premium_salary"`
	ExpTotalTtdAnnualOfficePremium          float64 `json:"exp_total_ttd_annual_office_premium" csv:"exp_total_ttd_annual_office_premium"`
	ExpTtdOfficeRatePer1000SA               float64 `json:"exp_ttd_office_rate_per_1000_sa" csv:"exp_ttd_office_rate_per_1000_sa"`
	ExpProportionTtdOfficePremiumSalary     float64 `json:"exp_proportion_ttd_office_premium_salary" csv:"exp_proportion_ttd_office_premium_salary"`

	MinPhiIncome                            float64 `json:"min_phi_income" csv:"min_phi_income"`
	MaxPhiIncome                            float64 `json:"max_phi_income" csv:"max_phi_income"`
	MaxPhiCappedIncome                      float64 `json:"max_phi_capped_income" csv:"max_phi_capped_income"`
	TotalPhiIncome                          float64 `json:"total_phi_income" csv:"total_phi_income"`
	TotalPhiCappedIncome                    float64 `json:"total_phi_capped_income" csv:"total_phi_capped_income"`
	AveragePhiCappedIncome                  float64 `json:"average_phi_capped_income" csv:"average_phi_capped_income"`
	TotalPhiContributionWaiver              float64 `json:"total_phi_contribution_waiver" csv:"total_phi_contribution_waiver"`
	TotalPhiMedicalAidWaiver                float64 `json:"total_phi_medical_aid_waiver" csv:"total_phi_medical_aid_waiver"`
	TotalPhiMonthlyBenefit                  float64 `json:"total_phi_monthly_benefit" csv:"total_phi_monthly_benefit"`
	TotalPhiRiskRate                        float64 `json:"total_phi_risk_rate" csv:"total_phi_risk_rate"`
	TotalPhiAnnualRiskPremium               float64 `json:"total_phi_annual_risk_premium" csv:"total_phi_annual_risk_premium"`
	PhiRiskRatePer1000SA                    float64 `json:"phi_risk_rate_per_1000_sa" csv:"phi_risk_rate_per_1000_sa"`
	ProportionPhiAnnualRiskPremiumSalary    float64 `json:"proportion_phi_annual_risk_premium_salary" csv:"proportion_phi_annual_risk_premium_salary"`
	TotalPhiAnnualOfficePremium             float64 `json:"total_phi_annual_office_premium" csv:"total_phi_annual_office_premium"`
	PhiOfficeRatePer1000SA                  float64 `json:"phi_office_rate_per_1000_sa" csv:"phi_office_rate_per_1000_sa"`
	ProportionPhiOfficePremiumSalary        float64 `json:"proportion_phi_office_premium_salary" csv:"proportion_phi_office_premium_salary"`
	ExpTotalPhiRiskRate                     float64 `json:"exp_total_phi_risk_rate" csv:"exp_total_phi_risk_rate"`
	ExpTotalPhiAnnualRiskPremium            float64 `json:"exp_total_phi_annual_risk_premium" csv:"exp_total_phi_annual_risk_premium"`
	ExpPhiRiskRatePer1000SA                 float64 `json:"exp_phi_risk_rate_per_1000_sa" csv:"exp_phi_risk_rate_per_1000_sa"`
	ExpProportionPhiAnnualRiskPremiumSalary float64 `json:"exp_proportion_phi_annual_risk_premium_salary" csv:"exp_proportion_phi_annual_risk_premium_salary"`
	ExpTotalPhiAnnualOfficePremium          float64 `json:"exp_total_phi_annual_office_premium" csv:"exp_total_phi_annual_office_premium"`
	ExpPhiOfficeRatePer1000SA               float64 `json:"exp_phi_office_rate_per_1000_sa" csv:"exp_phi_office_rate_per_1000_sa"`
	ExpProportionPhiOfficePremiumSalary     float64 `json:"exp_proportion_phi_office_premium_salary" csv:"exp_proportion_phi_office_premium_salary"`

	//TotalFunRiskRate                        float64 `json:"total_fun_risk_rate" csv:"total_fun_risk_rate"`
	TotalFunAnnualRiskPremium float64 `json:"total_fun_annual_risk_premium" csv:"total_fun_annual_risk_premium"`
	//FunRiskRatePer1000SA                    float64 `json:"fun_risk_rate_per_1000_sa" csv:"fun_risk_rate_per_1000_sa"`
	ProportionFunAnnualRiskPremiumSalary float64 `json:"proportion_fun_annual_risk_premium_salary" csv:"proportion_fun_annual_risk_premium_salary"`
	TotalFunAnnualOfficePremium          float64 `json:"total_fun_annual_office_premium" csv:"total_fun_annual_office_premium"`
	//FunOfficeRatePer1000SA                  float64 `json:"fun_office_rate_per_1000_sa" csv:"fun_office_rate_per_1000_sa"`
	ProportionFunOfficePremiumSalary float64 `json:"proportion_fun_office_premium_salary" csv:"proportion_fun_office_premium_salary"`
	//ExpTotalFunRiskRate                     float64 `json:"exp_total_fun_risk_rate" csv:"exp_total_fun_risk_rate"`
	ExpTotalFunAnnualRiskPremium float64 `json:"exp_total_fun_annual_risk_premium" csv:"exp_total_fun_annual_risk_premium"`
	//ExpFunRiskRatePer1000SA                 float64 `json:"exp_fun_risk_rate_per_1000_sa" csv:"exp_fun_risk_rate_per_1000_sa"`
	ExpProportionFunAnnualRiskPremiumSalary float64 `json:"exp_proportion_fun_annual_risk_premium_salary" csv:"exp_proportion_fun_annual_risk_premium_salary"`
	ExpTotalFunAnnualOfficePremium          float64 `json:"exp_total_fun_annual_office_premium" csv:"exp_total_fun_annual_office_premium"`
	//ExpFunOfficeRatePer1000SA               float64 `json:"exp_fun_office_rate_per_1000_sa" csv:"exp_fun_office_rate_per_1000_sa"`
	ExpProportionFunOfficePremiumSalary float64 `json:"exp_proportion_fun_office_premium_salary" csv:"exp_proportion_fun_office_premium_salary"`
	TotalFunAnnualPremiumPerMember      float64 `json:"total_fun_annual_premium_per_member" csv:"total_fun_annual_premium_per_member"`
	TotalFunMonthlyPremiumPerMember     float64 `json:"total_fun_monthly_premium_per_member" csv:"total_fun_monthly_premium_per_member"`
	ExpTotalFunAnnualPremiumPerMember   float64 `json:"exp_total_fun_annual_premium_per_member" csv:"exp_total_fun_annual_premium_per_member"`
	ExpTotalFunMonthlyPremiumPerMember  float64 `json:"exp_total_fun_monthly_premium_per_member" csv:"exp_total_fun_monthly_premium_per_member"`

	TotalRiskWeightedEducatorSumAssured         float64 `json:"total_risk_weighted_educator_sum_assured" csv:"total_risk_weighted_educator_sum_assured"`
	TotalEducatorRiskPremium                    float64 `json:"total_educator_risk_premium" csv:"total_educator_risk_premium"`
	TotalEducatorOfficePremium                  float64 `json:"total_educator_office_premium" csv:"total_educator_office_premium"`
	ExpAdjTotalEducatorRiskPremium              float64 `json:"exp_adj_total_educator_risk_premium" csv:"exp_adj_total_educator_risk_premium"`
	ExpAdjTotalEducatorOfficePremium            float64 `json:"exp_adj_total_educator_office_premium" csv:"exp_adj_total_educator_office_premium"`
	ProportionEducatorRiskPremiumSalary         float64 `json:"proportion_educator_risk_premium_salary" csv:"proportion_educator_risk_premium_salary"`
	ProportionEducatorOfficePremiumSalary       float64 `json:"proportion_educator_office_premium_salary" csv:"proportion_educator_office_premium_salary"`
	ExpAdjProportionEducatorRiskPremiumSalary   float64 `json:"exp_adj_proportion_educator_risk_premium_salary" csv:"exp_adj_proportion_educator_risk_premium_salary"`
	ExpAdjProportionEducatorOfficePremiumSalary float64 `json:"exp_adj_proportion_educator_office_premium_salary" csv:"exp_adj_proportion_educator_office_premium_salary"`

	ExceedsNormalRetirementAgeIndicator int `json:"exceeds_normal_retirement_age_indicator" csv:"exceeds_normal_retirement_age_indicator"`
	ExceedsFreeCoverLimitIndicator      int `json:"exceeds_free_cover_limit_indicator" csv:"exceeds_free_cover_limit_indicator"`
	//TotalGlaExperienceAdjustedAnnualPremium     float64 `json:"total_gla_experience_adjusted_annual_premium" csv:"total_gla_experience_adjusted_annual_premium"`
	//TotalPtdExperienceAdjustedAnnualPremium     float64 `json:"total_ptd_experience_adjusted_annual_premium" csv:"total_ptd_experience_adjusted_annual_premium"`
	//TotalTtdExperienceAdjustedAnnualPremium     float64 `json:"total_ttd_experience_adjusted_annual_premium" csv:"total_ttd_experience_adjusted_annual_premium"`
	//TotalPhiExperienceAdjustedAnnualPremium     float64 `json:"total_phi_experience_adjusted_annual_premium" csv:"total_phi_experience_adjusted_annual_premium"`
	//TotalCiExperienceAdjustedAnnualPremium      float64 `json:"total_ci_experience_adjusted_annual_premium" csv:"total_ci_experience_adjusted_annual_premium"`
	//TotalSpouseExperienceAdjustedAnnualPremium  float64 `json:"total_spouse_experience_adjusted_annual_premium" csv:"total_spouse_experience_adjusted_annual_premium"`
	//TotalFuneralExperienceAdjustedAnnualPremium float64 `json:"total_funeral_experience_adjusted_annual_premium" csv:"total_funeral_experience_adjusted_annual_premium"`
	TotalAnnualPremiumExcludingFuneral         float64 `json:"total_annual_premium_excluding_funeral" csv:"total_annual_premium_excluding_funeral"`
	TotalSumAssured                            float64 `json:"total_sum_assured" csv:"total_sum_assured"`
	TotalAnnualPremium                         float64 `json:"total_annual_premium" csv:"total_annual_premium" gorm:"default:0"`
	ExpTotalAnnualPremiumExclFuneral           float64 `json:"exp_total_annual_premium_excl_funeral" csv:"exp_total_annual_premium_excl_funeral"`
	ProportionExpTotalPremiumExclFuneralSalary float64 `json:"proportion_exp_total_premium_excl_funeral_salary" csv:"proportion_exp_total_premium_excl_funeral_salary"`
	TotalCommission                            float64 `json:"total_commission" csv:"total_commission"`
	TotalExpenses                              float64 `json:"total_expenses" csv:"total_expenses"`
	TotalExpectedClaims                        float64 `json:"total_expected_claims" csv:"total_expected_claims"`

	AnnualGlaExperienceWeightedRate    float64   `json:"annual_gla_experience_weighted_rate" csv:"annual_gla_experience_weighted_rate"`
	AnnualPtdExperienceWeightedRate    float64   `json:"annual_ptd_experience_weighted_rate" csv:"annual_ptd_experience_weighted_rate"`
	AnnualCiExperienceWeightedRate     float64   `json:"annual_ci_experience_weighted_rate" csv:"annual_ci_experience_weighted_rate"`
	CredibilityRate                    float64   `json:"credibility_rate" csv:"credibility_rate"`
	ManuallyAddedCredibility           float64   `json:"manually_added_credibility" csv:"manually_added_credibility"`
	PremiumRatesGuaranteedPeriodMonths int       `json:"premium_rates_guaranteed_period_months" csv:"premium_rates_guaranteed_period_months"`
	QuoteValidityPeriodMonths          int       `json:"quote_validity_period_months" csv:"quote_validity_period_months"`
	CreationDate                       time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy                          string    `json:"created_by" csv:"created_by"`
}

type MovementMemberRatingResult struct {
	FinancialYear                       int     `json:"financial_year" csv:"financial_year"`
	SchemeName                          string  `json:"-" csv:"scheme_name"`
	Category                            string  `json:"category" csv:"category"`
	QuoteId                             int     `json:"-" csv:"quote_id" gorm:"index"`
	MemberName                          string  `json:"member_name" csv:"member_name"`
	Gender                              string  `json:"gender" csv:"gender"`
	DateOfBirth                         string  `json:"date_of_birth" csv:"date_of_birth"`
	MovementType                        string  `json:"movement_type" csv:"movement_type"`
	AnnualSalary                        float64 `json:"annual_salary" csv:"annual_salary"`
	BenefitSalaryMultiple               float64 `json:"benefit_salary_multiple" csv:"benefit_salary_multiple"`
	Occupation                          string  `json:"occupation" csv:"occupation"`
	OccupationClass                     int     `json:"occupation_class" csv:"occupation_class"`
	Industry                            string  `json:"industry" csv:"industry"`
	AgeNextBirthday                     int     `json:"age_next_birthday" csv:"age_next_birthday"`
	AgeBand                             string  `json:"age_band" csv:"age_band"`
	SpouseGender                        string  `json:"spouse_gender" csv:"spouse_gender"`
	SpouseAgeNextBirthday               int     `json:"spouse_age_next_birthday" csv:"spouse_age_next_birthday"`
	AverageDependantAgeNextBirthday     float64 `json:"average_dependant_age_next_birthday" csv:"average_dependant_age_next_birthday"`
	AverageChildAgeNextBirthday         float64 `json:"average_child_age_next_birthday" csv:"average_child_age_next_birthday"`
	AverageNumberDependants             float64 `json:"average_number_dependants" csv:"average_number_dependants"`
	AverageNumberChildren               float64 `json:"average_number_children" csv:"average_number_children"`
	GlaSumAssured                       float64 `json:"gla_sum_assured" csv:"gla_sum_assured"`
	GlaCappedSumAssured                 float64 `json:"gla_capped_sum_assured" csv:"gla_capped_sum_assured"`
	PtdSumAssured                       float64 `json:"ptd_sum_assured" csv:"ptd_sum_assured"`
	PtdCappedSumAssured                 float64 `json:"ptd_capped_sum_assured" csv:"ptd_capped_sum_assured"`
	CiSumAssured                        float64 `json:"ci_sum_assured" csv:"ci_sum_assured"`
	CiCappedSumAssured                  float64 `json:"ci_capped_sum_assured" csv:"ci_capped_sum_assured"`
	SpouseGlaSumAssured                 float64 `json:"spouse_gla_sum_assured" csv:"spouse_gla_sum_assured"`
	SpouseGlaCappedSumAssured           float64 `json:"spouse_gla_capped_sum_assured" csv:"spouse_gla_capped_sum_assured"`
	TtdIncome                           float64 `json:"ttd_income" csv:"ttd_income"`
	TtdCappedIncome                     float64 `json:"ttd_capped_income" csv:"ttd_capped_income"`
	TtdNumberOfMonthlyPayments          float64 `json:"ttd_number_of_monthly_payments" csv:"ttd_number_of_monthly_payments"`
	IncomeReplacementRatio              float64 `json:"income_replacement_ratio" csv:"income_replacement_ratio"`
	PhiIncome                           float64 `json:"phi_income" csv:"phi_income"`
	PhiCappedIncome                     float64 `json:"phi_capped_income" csv:"phi_capped_income"`
	PhiContributionWaiver               float64 `json:"phi_contribution_waiver" csv:"phi_contribution_waiver"`
	PhiMedicalAidWaiver                 float64 `json:"phi_medical_aid_waiver" csv:"phi_medical_aid_waiver"`
	PhiMonthlyBenefit                   float64 `json:"phi_monthly_benefit" csv:"phi_monthly_benefit"`
	PhiAnnuityFactor                    float64 `json:"phi_annuity_factor" csv:"phi_annuity_factor"`
	BaseGlaRate                         float64 `json:"base_gla_rate" csv:"base_gla_rate"`
	GlaLoading                          float64 `json:"gla_loading" csv:"gla_loading"`
	GlaTerminalIllnessLoading           float64 `json:"gla_terminal_illness_loading" csv:"gla_terminal_illness_loading"`
	LoadedGlaRate                       float64 `json:"loaded_gla_rate" csv:"loaded_gla_rate"`
	ExpCredibility                      float64 `json:"exp_credibility" csv:"exp_credibility"`
	GlaExpRate                          float64 `json:"gla_exp_rate" csv:"gla_exp_rate"`
	GlaExpAdjustedRate                  float64 `json:"gla_exp_adjusted_rate" csv:"gla_exp_adjusted_rate"`
	ExpAdjLoadedGlaRate                 float64 `json:"exp_adj_loaded_gla_rate" csv:"exp_adj_loaded_gla_rate"`
	BasePtdRate                         float64 `json:"base_ptd_rate" csv:"base_ptd_rate"`
	PtdLoading                          float64 `json:"ptd_loading" csv:"ptd_loading"`
	LoadedPtdRate                       float64 `json:"loaded_ptd_rate" csv:"loaded_ptd_rate"`
	ExpAdjLoadedPtdRate                 float64 `json:"exp_adj_loaded_ptd_rate" csv:"exp_adj_loaded_ptd_rate"`
	BaseTtdRate                         float64 `json:"base_ttd_rate" csv:"base_ttd_rate"`
	TtdLoading                          float64 `json:"ttd_loading" csv:"ttd_loading"`
	LoadedTtdRate                       float64 `json:"loaded_ttd_rate" csv:"loaded_ttd_rate"`
	ExpAdjLoadedTtdRate                 float64 `json:"exp_adj_loaded_ttd_rate" csv:"exp_adj_loaded_ttd_rate"`
	BasePhiRate                         float64 `json:"base_phi_rate" csv:"base_phi_rate"`
	PhiSalaryLevel                      float64 `json:"phi_salary_level" csv:"phi_salary_level"`
	PhiLoading                          float64 `json:"phi_loading" csv:"phi_loading"`
	LoadedPhiRate                       float64 `json:"loaded_phi_rate" csv:"loaded_phi_rate"`
	ExpAdjLoadedPhiRate                 float64 `json:"exp_adj_loaded_phi_rate" csv:"exp_adj_loaded_phi_rate"`
	BaseCiRate                          float64 `json:"base_ci_rate" csv:"base_ci_rate"`
	CiLoading                           float64 `json:"ci_loading" csv:"ci_loading"`
	LoadedCiRate                        float64 `json:"loaded_ci_rate" csv:"loaded_ci_rate"`
	ExpAdjLoadedCiRate                  float64 `json:"exp_adj_loaded_ci_rate" csv:"exp_adj_loaded_ci_rate"`
	BaseSpouseGlaRate                   float64 `json:"base_spouse_gla_rate" csv:"base_spouse_gla_rate"`
	SpouseGlaLoading                    float64 `json:"spouse_gla_loading" csv:"spouse_gla_loading"`
	LoadedSpouseGlaRate                 float64 `json:"loaded_spouse_gla_rate" csv:"loaded_spouse_gla_rate"`
	ExpAdjLoadedSpouseGlaRate           float64 `json:"exp_adj_loaded_spouse_gla_rate" csv:"exp_adj_loaded_spouse_gla_rate"`
	ExpenseLoading                      float64 `json:"expense_loading" csv:"expense_loading"`
	CommissionLoading                   float64 `json:"commission_loading" csv:"commission_loading"`
	ProfitLoading                       float64 `json:"profit_loading" csv:"profit_loading"`
	Discount                            float64 `json:"discount" csv:"discount"`
	GlaRiskPremium                      float64 `json:"gla_risk_premium" csv:"gla_risk_premium"`
	GlaOfficePremium                    float64 `json:"gla_office_premium" csv:"gla_office_premium"`
	ExpAdjGlaRiskPremium                float64 `json:"exp_adj_gla_risk_premium" csv:"exp_adj_gla_risk_premium"`
	ExpAdjGlaOfficePremium              float64 `json:"exp_adj_gla_office_premium" csv:"exp_adj_gla_office_premium"`
	PtdRiskPremium                      float64 `json:"ptd_risk_premium" csv:"ptd_risk_premium"`
	PtdOfficePremium                    float64 `json:"ptd_office_premium" csv:"ptd_office_premium"`
	ExpAdjPtdRiskPremium                float64 `json:"exp_adj_ptd_risk_premium" csv:"exp_adj_ptd_risk_premium"`
	ExpAdjPtdOfficePremium              float64 `json:"exp_adj_ptd_office_premium" csv:"exp_adj_ptd_office_premium"`
	TtdRiskPremium                      float64 `json:"ttd_risk_premium" csv:"ttd_risk_premium"`
	TtdOfficePremium                    float64 `json:"ttd_office_premium" csv:"ttd_office_premium"`
	ExpAdjTtdRiskPremium                float64 `json:"exp_adj_ttd_risk_premium" csv:"exp_adj_ttd_risk_premium"`
	ExpAdjTtdOfficePremium              float64 `json:"exp_adj_ttd_office_premium" csv:"exp_adj_ttd_office_premium"`
	PhiRiskPremium                      float64 `json:"phi_risk_premium" csv:"phi_risk_premium"`
	PhiOfficePremium                    float64 `json:"phi_office_premium" csv:"phi_office_premium"`
	ExpAdjPhiRiskPremium                float64 `json:"exp_adj_phi_risk_premium" csv:"exp_adj_phi_risk_premium"`
	ExpAdjPhiOfficePremium              float64 `json:"exp_adj_phi_office_premium" csv:"exp_adj_phi_office_premium"`
	CiRiskPremium                       float64 `json:"ci_risk_premium" csv:"ci_risk_premium"`
	CiOfficePremium                     float64 `json:"ci_office_premium" csv:"ci_office_premium"`
	ExpAdjCiRiskPremium                 float64 `json:"exp_adj_ci_risk_premium" csv:"exp_adj_ci_risk_premium"`
	ExpAdjCiOfficePremium               float64 `json:"exp_adj_ci_office_premium" csv:"exp_adj_ci_office_premium"`
	SpouseGlaRiskPremium                float64 `json:"spouse_gla_risk_premium" csv:"spouse_gla_risk_premium"`
	SpouseGlaOfficePremium              float64 `json:"spouse_gla_office_premium" csv:"spouse_gla_office_premium"`
	ExpAdjSpouseGlaRiskPremium          float64 `json:"exp_adj_spouse_gla_risk_premium" csv:"exp_adj_spouse_gla_risk_premium"`
	ExpAdjSpouseGlaOfficePremium        float64 `json:"exp_adj_spouse_gla_office_premium" csv:"exp_adj_spouse_gla_office_premium"`
	MarriageProportion                  float64 `json:"marriage_proportion" csv:"marriage_proportion"`
	ChildFuneralBaseRate                float64 `json:"child_funeral_base_rate" csv:"child_funeral_base_rate"`
	ChildFuneralSumAssured              float64 `json:"child_funeral_sum_assured" csv:"child_funeral_sum_assured"`
	DependantFuneralBaseRate            float64 `json:"dependant_funeral_base_rate" csv:"dependant_funeral_base_rate"`
	DependantFuneralSumAssured          float64 `json:"dependant_funeral_sum_assured" csv:"dependant_funeral_sum_assured"`
	MainMemberFuneralCost               float64 `json:"main_member_funeral_cost" csv:"main_member_funeral_cost"`
	SpouseFuneralCost                   float64 `json:"spouse_funeral_cost" csv:"spouse_funeral_cost"`
	ChildrenFuneralCost                 float64 `json:"children_funeral_cost" csv:"children_funeral_cost"`
	DependantsFuneralCost               float64 `json:"dependants_funeral_cost" csv:"dependants_funeral_cost"`
	TotalFuneralRiskCost                float64 `json:"total_funeral_risk_cost" csv:"total_funeral_risk_cost"`
	TotalFuneralOfficeCost              float64 `json:"total_funeral_office_cost" csv:"total_funeral_office_cost"`
	ExpAdjTotalFuneralRiskCost          float64 `json:"exp_adj_total_funeral_risk_cost" csv:"exp_adj_total_funeral_risk_cost"`
	ExpAdjTotalFuneralOfficeCost        float64 `json:"exp_adj_total_funeral_office_cost" csv:"exp_adj_total_funeral_office_cost"`
	ExceedsNormalRetirementAgeIndicator int     `json:"exceeds_normal_retirement_age_indicator" csv:"exceeds_normal_retirement_age_indicator"`
	ExceedsFreeCoverLimitIndicator      int     `json:"exceeds_free_cover_limit_indicator" csv:"exceeds_free_cover_limit_indicator"`
	GlaExperienceAdjustment             float64 `json:"gla_experience_adjustment" csv:"gla_experience_adjustment"`
	PtdExperienceAdjustment             float64 `json:"ptd_experience_adjustment" csv:"ptd_experience_adjustment"`
	CiExperienceAdjustment              float64 `json:"ci_experience_adjustment" csv:"ci_experience_adjustment"`
	PhiExperienceAdjustment             float64 `json:"phi_experience_adjustment" csv:"phi_experience_adjustment"`
	TtdExperienceAdjustment             float64 `json:"ttd_experience_adjustment" csv:"ttd_experience_adjustment"`
	//GlaExperienceAdjustedAnnualRate        float64 `json:"gla_experience_adjusted_annual_rate" csv:"gla_experience_adjusted_annual_rate"`
	//GlaExperienceAdjustedAnnualPremium     float64 `json:"gla_experience_adjusted_annual_premium" csv:"gla_experience_adjusted_annual_premium"`
	//PtdExperienceAdjustedAnnualPremium     float64 `json:"ptd_experience_adjusted_annual_premium" csv:"age_next_birthday"`
	//TtdExperienceAdjustedAnnualPremium     float64 `json:"ttd_experience_adjusted_annual_premium" csv:"ttd_experience_adjusted_annual_premium"`
	//PhiExperienceAdjustedAnnualPremium     float64 `json:"phi_experience_adjusted_annual_premium" csv:"phi_experience_adjusted_annual_premium"`
	//CiExperienceAdjustedAnnualPremium      float64 `json:"ci_experience_adjusted_annual_premium" csv:"ci_experience_adjusted_annual_premium"`
	//SpouseExperienceAdjustedAnnualPremium  float64 `json:"spouse_experience_adjusted_annual_premium" csv:"spouse_experience_adjusted_annual_premium"`
	FuneralExperienceAdjustedAnnualPremium float64   `json:"funeral_experience_adjusted_annual_premium" csv:"funeral_experience_adjusted_annual_premium"`
	CreationDate                           time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy                              string    `json:"created_by" csv:"created_by"`
}

type GlaRate struct {
	ID              int       `json:"id" gorm:"primary_key"`
	RiskRateCode    string    `json:"risk_rate_code" csv:"risk_rate_code"`
	AgeNextBirthday int       `json:"age_next_birthday" csv:"age_next_birthday"`
	IncomeLevel     string    `json:"income_level" csv:"income_level"`
	Gender          string    `json:"gender" csv:"gender"`
	WaitingPeriod   int       `json:"waiting_period" csv:"waiting_period"`
	BenefitType     string    `json:"benefit_type" csv:"benefit_type"`
	Qx              float64   `json:"qx" csv:"qx"`
	CreationDate    time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy       string    `json:"created_by" csv:"created_by"`
}

type PtdRate struct {
	ID                   int       `json:"-" gorm:"primary_key"`
	RiskRateCode         string    `json:"risk_rate_code" csv:"risk_rate_code"`
	AgeNextBirthday      int       `json:"age_next_birthday" csv:"age_next_birthday"`
	Gender               string    `json:"gender" csv:"gender"`
	OccupationClass      int       `json:"occupation_class" csv:"occupation_class"`
	IncomeLevel          int       `json:"income_level" csv:"income_level"`
	WaitingPeriod        int       `json:"waiting_period" csv:"waiting_period"`
	DeferredPeriod       int       `json:"deferred_period" csv:"deferred_period"`
	DisabilityDefinition string    `json:"disability_definition" csv:"disability_definition"`
	PtdRate              float64   `json:"ptd_rate" csv:"ptd_rate"`
	RiskType             string    `json:"risk_type" csv:"risk_type"`
	CreationDate         time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy            string    `json:"created_by" csv:"created_by"`
}

type CiRate struct {
	ID                int       `json:"id" gorm:"primary_key"`
	RiskRateCode      string    `json:"risk_rate_code" csv:"risk_rate_code"`
	AgeNextBirthday   int       `json:"age_next_birthday" csv:"age_next_birthday"`
	Gender            string    `json:"gender" csv:"gender"`
	OccupationClass   int       `json:"occupation_class" csv:"occupation_class"`
	IncomeLevel       int       `json:"income_level" csv:"income_level"`
	WaitingPeriod     int       `json:"waiting_period" csv:"waiting_period"`
	BenefitDefinition string    `json:"benefit_definition" csv:"benefit_definition"`
	CiRate            float64   `json:"ci_rate" csv:"ci_rate"`
	CreationDate      time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy         string    `json:"created_by" csv:"created_by"`
}

type AccidentalTtdRate struct {
	ID              int       `json:"id" gorm:"primary_key"`
	RiskRateCode    string    `json:"risk_rate_code" csv:"risk_rate_code"`
	AgeNextBirthday int       `json:"age_next_birthday" csv:"age_next_birthday"`
	Gender          string    `json:"gender" csv:"gender"`
	AccTtdRate      float64   `json:"acc_ttd_rate" csv:"acc_ttd_rate"`
	CreationDate    time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy       string    `json:"created_by" csv:"created_by"`
}

type TtdRate struct {
	ID                   int       `json:"id" gorm:"primary_key"`
	RiskRateCode         string    `json:"risk_rate_code" csv:"risk_rate_code"`
	AgeNextBirthday      int       `json:"age_next_birthday" csv:"age_next_birthday"`
	Gender               string    `json:"gender" csv:"gender"`
	OccupationClass      int       `json:"occupation_class" csv:"occupation_class"`
	IncomeLevel          int       `json:"income_level" csv:"income_level"`
	WaitingPeriod        int       `json:"waiting_period" csv:"waiting_period"`
	DeferredPeriod       int       `json:"deferred_period" csv:"deferred_period"`
	DisabilityDefinition string    `json:"disability_definition" csv:"disability_definition"`
	RiskType             string    `json:"risk_type" csv:"risk_type"`
	TtdRate              float64   `json:"ttd_rate" csv:"ttd_rate"`
	CreationDate         time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy            string    `json:"created_by" csv:"created_by"`
}

type PhiRate struct {
	ID                      int       `json:"-" gorm:"primary_key"`
	RiskRateCode            string    `json:"risk_rate_code" csv:"risk_rate_code"`
	AgeNextBirthday         int       `json:"age_next_birthday" csv:"age_next_birthday"`
	Gender                  string    `json:"gender" csv:"gender"`
	OccupationClass         int       `json:"occupation_class" csv:"occupation_class"`
	IncomeLevel             int       `json:"income_level" csv:"income_level"`
	WaitingPeriod           int       `json:"waiting_period" csv:"waiting_period"`
	DeferredPeriod          int       `json:"deferred_period" csv:"deferred_period"`
	NormalRetirementAge     int       `json:"normal_retirement_age" csv:"normal_retirement_age"`
	BenefitEscalationOption string    `json:"benefit_escalation_option" csv:"benefit_escalation_option"`
	DisabilityDefinition    string    `json:"disability_definition" csv:"disability_definition"`
	RiskType                string    `json:"risk_type" csv:"risk_type"`
	PhiRate                 float64   `json:"phi_rate" csv:"phi_rate"`
	CreationDate            time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy               string    `json:"created_by" csv:"created_by"`
}

type ChildMortality struct {
	ID              int       `json:"id" gorm:"primary_key"`
	RiskRateCode    string    `json:"risk_rate_code" csv:"risk_rate_code"`
	AgeNextBirthday int       `json:"age_next_birthday" csv:"age_next_birthday"`
	ChildRate       float64   `json:"child_rate" csv:"child_rate"`
	CreationDate    time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy       string    `json:"created_by" csv:"created_by"`
}

type IndustryLoading struct {
	ID                     int       `json:"id" gorm:"primary_key"`
	RiskRateCode           string    `json:"risk_rate_code" csv:"risk_rate_code"`
	OccupationClass        int       `json:"occupation_class" csv:"occupation_class"`
	Gender                 string    `json:"gender" csv:"gender"`
	GlaIndustryLoadingRate float64   `json:"gla_industry_loading_rate" csv:"gla_industry_loading_rate"`
	PtdIndustryLoadingRate float64   `json:"ptd_industry_loading_rate" csv:"ptd_industry_loading_rate"`
	CiIndustryLoadingRate  float64   `json:"ci_industry_loading_rate" csv:"ci_industry_loading_rate"`
	TtdIndustryLoadingRate float64   `json:"ttd_industry_loading_rate" csv:"ttd_industry_loading_rate"`
	PhiIndustryLoadingRate float64   `json:"phi_industry_loading_rate" csv:"phi_industry_loading_rate"`
	CreationDate           time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy              string    `json:"created_by" csv:"created_by"`
}

type FuneralParameters struct {
	ID int `json:"id" gorm:"primary_key"`
	//Year                int       `json:"year" csv:"year"`
	RiskRateCode        string    `json:"risk_rate_code" csv:"risk_rate_code"`
	AgeNextBirthday     int       `json:"age_next_birthday" csv:"age_next_birthday"`
	ProportionMarried   float64   `json:"proportion_married" csv:"proportion_married"`
	AverageChildAge     float64   `json:"average_child_age" csv:"average_child_age"`
	AverageDependantAge float64   `json:"average_dependant_age" csv:"average_dependant_age"`
	NumberChildren      float64   `json:"number_children" csv:"number_children"`
	NumberDependants    float64   `json:"number_dependants" csv:"number_dependants"`
	CreationDate        time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy           string    `json:"created_by" csv:"created_by"`
}

type GroupPricingReinsuranceStructure struct {
	ID                                   int       `json:"id" gorm:"primary_key"`
	RiskRateCode                         string    `json:"risk_rate_code" csv:"risk_rate_code"`
	Basis                                string    `json:"basis" csv:"basis"`
	TreatyCode                           string    `json:"treaty_code" csv:"treaty_code"`
	RiskPremiumBasisIndicator            bool      `json:"risk_premium_basis_indicator" csv:"risk_premium_basis_indicator"`
	FuneralReinsuranceInclusionIndicator bool      `json:"funeral_reinsurance_inclusion_indicator" csv:"funeral_reinsurance_inclusion_indicator"`
	FlatAnnualReinsPremRate              float64   `json:"flat_annual_reins_prem_rate" csv:"flat_annual_reins_prem_rate"`
	Level1CededProportion                float64   `json:"level1_ceded_proportion" csv:"level1_ceded_proportion"`
	Level1Lowerbound                     float64   `json:"level1_lowerbound" csv:"level1_lowerbound"`
	Level1Upperbound                     float64   `json:"level1_upperbound" csv:"level1_upperbound"`
	Level2CededProportion                float64   `json:"level2_ceded_proportion" csv:"level2_ceded_proportion"`
	Level2Lowerbound                     float64   `json:"level2_lowerbound" csv:"level2_lowerbound"`
	Level2Upperbound                     float64   `json:"level2_upperbound" csv:"level2_upperbound"`
	Level3CededProportion                float64   `json:"level3_ceded_proportion" csv:"level3_ceded_proportion"`
	Level3Lowerbound                     float64   `json:"level3_lowerbound" csv:"level3_lowerbound"`
	Level3Upperbound                     float64   `json:"level3_upperbound" csv:"level3_upperbound"`
	IncomeLevel1CededProportion          float64   `json:"income_level1_ceded_proportion" csv:"income_level1_ceded_proportion"`
	IncomeLevel1Lowerbound               float64   `json:"income_level1_lowerbound" csv:"income_level1_lowerbound"`
	IncomeLevel1Upperbound               float64   `json:"income_level1_upperbound" csv:"income_level1_upperbound"`
	IncomeLevel2CededProportion          float64   `json:"income_level2_ceded_proportion" csv:"income_level2_ceded_proportion"`
	IncomeLevel2Lowerbound               float64   `json:"income_level2_lowerbound" csv:"income_level2_lowerbound"`
	IncomeLevel2Upperbound               float64   `json:"income_level2_upperbound" csv:"income_level2_upperbound"`
	IncomeLevel3CededProportion          float64   `json:"income_level3_ceded_proportion" csv:"income_level3_ceded_proportion"`
	IncomeLevel3Lowerbound               float64   `json:"income_level3_lowerbound" csv:"income_level3_lowerbound"`
	IncomeLevel3Upperbound               float64   `json:"income_level3_upperbound" csv:"income_level3_upperbound"`
	LeadReinsurerShare                   float64   `json:"lead_reinsurer_share" csv:"lead_reinsurer_share"`
	NonLeadReinsurer1Share               float64   `json:"non_lead_reinsurer1_share" csv:"non_lead_reinsurer1_share"`
	NonLeadReinsurer2Share               float64   `json:"non_lead_reinsurer2_share" csv:"non_lead_reinsurer2_share"`
	NonLeadReinsurer3Share               float64   `json:"non_lead_reinsurer3_share" csv:"non_lead_reinsurer3_share"`
	CedingCommission                     float64   `json:"ceding_commission" csv:"ceding_commission"`
	CreationDate                         time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy                            string    `json:"created_by" csv:"created_by"`
}

type GroupPricingParameters struct {
	ID int `json:"id" gorm:"primary_key"`
	//Year                               int       `json:"year" csv:"year"`
	Basis                              string    `json:"basis" csv:"basis"`
	RiskRateCode                       string    `json:"risk_rate_code" csv:"risk_rate_code"`
	TreatyCode                         string    `json:"treaty_code" csv:"treaty_code"`
	EducatorBenefitCode                string    `json:"educator_benefit_code" csv:"educator_benefit_code"`
	SpouseAgeGap                       int       `json:"spouse_age_gap" csv:"spouse_age_gap"`
	ReinsurerProfitLoading             float64   `json:"reinsurer_profit_loading" csv:"reinsurer_profit_loading"`
	RiskMarginRate                     float64   `json:"risk_margin_rate" csv:"risk_margin_rate"`
	IsLumpsumReinsGLADependent         bool      `json:"is_lumpsum_reins_gla_dependent" csv:"is_lumpsum_reins_gla_dependent"`
	PremiumRatesGuaranteedPeriodMonths int       `json:"premium_rates_guaranteed_period_months" csv:"premium_rates_guaranteed_period_months"`
	QuoteValidityPeriodMonths          int       `json:"quote_validity_period_months" csv:"quote_validity_period_months"`
	AnnualExpenseAmount                float64   `json:"annual_expense_amount" csv:"annual_expense_amount"`
	FullCredibilityThreshold           float64   `json:"full_credibility_threshold" csv:"full_credibility_threshold"`
	FreeCoverLimitScalingFactor        float64   `json:"free_cover_limit_scaling_factor" csv:"free_cover_limit_scaling_factor"`
	FreeCoverLimitPercentile           float64   `json:"free_cover_limit_percentile" csv:"free_cover_limit_percentile"`
	FreeCoverLimitNearestMultiple      float64   `json:"free_cover_limit_nearest_multiple" csv:"free_cover_limit_nearest_multiple"`
	GlobalGlaExperienceRate            float64   `json:"global_gla_experience_rate" csv:"global_gla_experience_rate"`
	GlobalPtdExperienceRate            float64   `json:"global_ptd_experience_rate" csv:"global_ptd_experience_rate"`
	GlobalCiExperienceRate             float64   `json:"global_ci_experience_rate" csv:"global_ci_experience_rate"`
	MedicalAidWaiverProportion         float64   `json:"medical_aid_waiver_proportion" csv:"medical_aid_waiver_proportion"`
	MedicalAidWaiverAmount             float64   `json:"medical_aid_waiver_amount" csv:"medical_aid_waiver_amount"`
	TtdNumberMonthlyPayments           float64   `json:"ttd_number_monthly_payments" csv:"ttd_number_monthly_payments"`
	CreationDate                       time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy                          string    `json:"created_by" csv:"created_by"`
}

type GroupScheme struct {
	ID                     int                 `json:"id" gorm:"primary_key"`
	Name                   string              `json:"name" csv:"name"`
	BrokerId               int                 `json:"broker_id" csv:"broker_id"`
	Broker                 Broker              `json:"broker" csv:"broker" gorm:"-"`
	DistributionChannel    DistributionChannel `json:"distribution_channel" csv:"distribution_channel" gorm:"size:20;default:'broker'"`
	ContactPerson          string              `json:"contact_person" csv:"contact_person"`
	ContactEmail           string              `json:"contact_email" csv:"contact_email"`
	DurationInForceDays    int                 `json:"duration_in_force_days" csv:"duration_in_force_days"`
	RenewalDate            time.Time           `json:"renewal_date" csv:"renewal_date"`
	MemberCount            float64             `json:"member_count" csv:"member_count"`
	AnnualPremium          float64             `json:"annual_premium" csv:"annual_premium"`
	GlaAnnualPremium       float64             `json:"gla_annual_premium" csv:"gla_annual_premium"`
	PtdAnnualPremium       float64             `json:"ptd_annual_premium" csv:"ptd_annual_premium"`
	CiAnnualPremium        float64             `json:"ci_annual_premium" csv:"ci_annual_premium"`
	SglaAnnualPremium      float64             `json:"sgla_annual_premium" csv:"sgla_annual_premium"`
	TtdAnnualPremium       float64             `json:"ttd_annual_premium" csv:"ttd_annual_premium"`
	PhiAnnualPremium       float64             `json:"phi_annual_premium" csv:"phi_annual_premium"`
	FuneralAnnualPremium   float64             `json:"funeral_annual_premium" csv:"funeral_annual_premium"`
	Commission             float64             `json:"commission" csv:"commission"`
	EarnedPremium          float64             `json:"earned_premium" csv:"earned_premium"`
	ExpectedExpenses       float64             `json:"expected_expenses" csv:"expected_expenses"`
	ExpectedGlaClaims      float64             `json:"expected_gla_claims" csv:"expected_gla_claims"`
	ExpectedPtdClaims      float64             `json:"expected_ptd_claims" csv:"expected_ptd_claims"`
	ExpectedCiClaims       float64             `json:"expected_ci_claims" csv:"expected_ci_claims"`
	ExpectedSglaClaims     float64             `json:"expected_sgla_claims" csv:"expected_sgla_claims"`
	ExpectedTtdClaims      float64             `json:"expected_ttd_claims" csv:"expected_ttd_claims"`
	ExpectedPhiClaims      float64             `json:"expected_phi_claims" csv:"expected_phi_claims"`
	ExpectedFunClaims      float64             `json:"expected_fun_claims" csv:"expected_fun_claims"`
	ExpectedClaims         float64             `json:"expected_claims" csv:"expected_claims"`
	ActualClaims           float64             `json:"actual_claims" csv:"actual_claims"`
	ExpectedClaimsRatio    float64             `json:"expected_claims_ratio" csv:"expected_claims_ratio"`
	ActualClaimsRatio      float64             `json:"actual_claims_ratio" csv:"actual_claims_ratio"`
	ExpectedLossRatio      float64             `json:"expected_loss_ratio" csv:"expected_loss_ratio"`
	ActualLossRatio        float64             `json:"actual_loss_ratio" csv:"actual_loss_ratio"`
	InForce                bool                `json:"in_force" csv:"in_force"`
	Status                 Status              `json:"status" csv:"status"`
	NewBusiness            bool                `json:"new_business" csv:"new_business"`
	SchemeStatusMessage    string              `json:"scheme_status_message" csv:"scheme_status_message"`
	CreationDate           time.Time           `json:"creation_date" gorm:"autoCreateTime"`
	CreatedBy              string              `json:"created_by"`
	QuoteId                int                 `json:"quote_id"`
	Quote                  GroupPricingQuote   `json:"quote" gorm:"foreignKey:QuoteId"`
	QuoteInForce           string              `json:"quote_in_force"`
	ActiveSchemeCategories StringArray         `json:"active_scheme_categories" gorm:"type:json"`
	CoverStartDate         time.Time           `json:"cover_start_date" csv:"cover_start_date"`
	CoverEndDate           time.Time           `json:"cover_end_date" csv:"cover_end_date"`
	CommencementDate       time.Time           `json:"commencement_date" csv:"commencement_date"`
	SchemeQuoteStatus      string              `json:"scheme_quote_status"`
	HasTreatyLink          bool                `json:"has_treaty_link" gorm:"-"`
}

type GroupPricingInsurerDetail struct {
	ID                    int       `json:"id" gorm:"primary_key"`
	Name                  string    `json:"name" csv:"name"`
	ContactPerson         string    `json:"contact_person" csv:"contact_person"`
	AddressLine1          string    `json:"address_line_1" csv:"address_line_1"`
	AddressLine2          string    `json:"address_line_2" csv:"address_line_2"`
	AddressLine3          string    `json:"address_line_3" csv:"address_line_3"`
	PostCode              string    `json:"post_code" csv:"post_code"`
	Province              string    `json:"province" csv:"province"`
	City                  string    `json:"city" csv:"city"`
	Country               string    `json:"country" csv:"country"`
	Telephone             string    `json:"telephone" csv:"telephone"`
	Email                 string    `json:"email" csv:"email"`
	CreationDate          time.Time `json:"creation_date" gorm:"autoCreateTime"`
	Logo                  []byte    `json:"logo"`
	LogoMimeType          string    `json:"logo_mime_type"`
	YearEndMonth          int       `json:"year_end_month" csv:"year_end_month"`
	IntroductoryText      string    `json:"introductory_text" csv:"introductory_text"`
	GeneralProvisionsText string    `json:"general_provisions_text" csv:"general_provisions_text"`
	OnRiskLetterText      string    `json:"on_risk_letter_text" csv:"on_risk_letter_text"`
}

// InsurerQuoteTemplate represents a custom DOCX template uploaded by an insurer
type InsurerQuoteTemplate struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	InsurerID  int       `json:"insurer_id" gorm:"index;not null"`
	Version    int       `json:"version"`              // 1, 2, 3, ... per insurer
	Filename   string    `json:"filename"`             // original upload filename
	DocxBlob   []byte    `json:"-"`                    // never serialised in JSON
	SizeBytes  int       `json:"size_bytes"`
	UploadedBy string    `json:"uploaded_by"`
	UploadedAt time.Time `json:"uploaded_at" gorm:"autoCreateTime"`
	IsActive   bool      `json:"is_active" gorm:"index"` // exactly one true per insurer
}

// InsurerOnRiskLetterTemplate represents a custom DOCX template for On Risk letters
type InsurerOnRiskLetterTemplate struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	InsurerID  int       `json:"insurer_id" gorm:"index;not null"`
	Version    int       `json:"version"`
	Filename   string    `json:"filename"`
	DocxBlob   []byte    `json:"-"`
	SizeBytes  int       `json:"size_bytes"`
	UploadedBy string    `json:"uploaded_by"`
	UploadedAt time.Time `json:"uploaded_at" gorm:"autoCreateTime"`
	IsActive   bool      `json:"is_active" gorm:"index"`
}

// OnRiskLetter records the issuance of an On Risk letter when a quote is accepted.
type OnRiskLetter struct {
	ID               int       `json:"id" gorm:"primaryKey;autoIncrement"`
	QuoteID          int       `json:"quote_id" gorm:"index;not null"`
	SchemeID         int       `json:"scheme_id" gorm:"index"`
	LetterDate       time.Time `json:"letter_date"`
	CommencementDate time.Time `json:"commencement_date"`
	CoverEndDate     time.Time `json:"cover_end_date"`
	GeneratedBy      string    `json:"generated_by"`
	LetterReference  string    `json:"letter_reference" gorm:"type:varchar(255);uniqueIndex"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type MemberPremiumSchedule struct {
	SchemeId                       int       `json:"-" csv:"scheme_id"`
	QuoteId                        int       `json:"-" csv:"quote_id" gorm:"index"`
	Category                       string    `json:"category" csv:"category"`
	MemberName                     string    `json:"member_name" csv:"member_name"`
	Gender                         string    `json:"gender" csv:"gender"`
	EntryDate                      time.Time `json:"entry_date" csv:"entry_date" gorm:"autoCreateTime"`
	ExitDate                       time.Time `json:"exit_date" csv:"exit_date" gorm:"autoCreateTime"`
	IsOriginalMember               bool      `json:"is_original_member" csv:"is_original_member"`
	GlaCoveredSumAssured           float64   `json:"gla_covered_sum_assured" csv:"gla_covered_sum_assured"`
	GlaAnnualPremium               float64   `json:"gla_annual_premium" csv:"gla_annual_premium"`
	PtdCoveredSumAssured           float64   `json:"ptd_covered_sum_assured" csv:"ptd_covered_sum_assured"`
	PtdAnnualPremium               float64   `json:"ptd_annual_premium" csv:"ptd_annual_premium"`
	CiCoveredSumAssured            float64   `json:"ci_covered_sum_assured" csv:"ci_covered_sum_assured"`
	CiAnnualPremium                float64   `json:"ci_annual_premium" csv:"ci_annual_premium"`
	TtdCoveredIncome               float64   `json:"ttd_covered_income" csv:"ttd_covered_income"`
	TtdAnnualPremium               float64   `json:"ttd_annual_premium" csv:"ttd_annual_premium"`
	SpouseGlaCoveredSumAssured     float64   `json:"spouse_gla_covered_sum_assured" csv:"spouse_gla_covered_sum_assured"`
	SpouseGlaAnnualPremium         float64   `json:"spouse_gla_annual_premium" csv:"spouse_gla_annual_premium"`
	PhiCoveredIncome               float64   `json:"phi_covered_income" csv:"phi_covered_income"`
	PhiAnnualPremium               float64   `json:"phi_annual_premium" csv:"phi_annual_premium"`
	MainMemberFuneralSumAssured    float64   `json:"main_member_funeral_sum_assured" csv:"main_member_funeral_sum_assured"`
	MainMemberFuneralAnnualPremium float64   `json:"main_member_funeral_annual_premium" csv:"main_member_funeral_annual_premium"`
	SpouseFuneralSumAssured        float64   `json:"spouse_funeral_sum_assured" csv:"spouse_funeral_sum_assured"`
	SpouseFuneralAnnualPremium     float64   `json:"spouse_funeral_annual_premium" csv:"spouse_funeral_annual_premium"`
	ChildFuneralSumAssured         float64   `json:"child_funeral_sum_assured" csv:"child_funeral_sum_assured"`
	ChildrenFuneralAnnualPremium   float64   `json:"children_funeral_annual_premium" csv:"children_funeral_annual_premium"`
	DependantsFuneralSumAssured    float64   `json:"dependants_funeral_sum_assured" csv:"dependants_funeral_sum_assured"`
	DependantsFuneralAnnualPremium float64   `json:"dependants_funeral_annual_premium" csv:"dependants_funeral_annual_premium"`
	TotalAnnualPremiumPayable      float64   `json:"total_annual_premium_payable" csv:"total_annual_premium_payable"`
	CreationDate                   time.Time `json:"creation_date" gorm:"autoCreateTime"`
	CreatedBy                      string    `json:"created_by"`
	Logo                           []byte    `json:"logo"`
}

type IncomeLevel struct {
	ID int `json:"id" gorm:"primary_key"`
	//Year         int       `json:"year" csv:"year"`
	RiskRateCode string    `json:"risk_rate_code" csv:"risk_rate_code"`
	MinIncome    float64   `json:"min_income" csv:"min_income"`
	MaxIncome    float64   `json:"max_income" csv:"max_income"`
	Level        int       `json:"level" csv:"level"`
	CreationDate time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy    string    `json:"created_by" csv:"created_by"`
}

type SchemeCategoryMaster struct {
	ID          int       `json:"id" gorm:"primary_key"`
	InsurerId   int       `json:"insurer_id" csv:"insurer_id"`
	Name        string    `json:"name" csv:"name"`
	Description string    `json:"description" csv:"description"`
	Active      bool      `json:"active" csv:"active"` // true: active, false: inactive
	CreatedBy   string    `json:"created_by" csv:"created_by"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type OccupationClass struct {
	ID           int       `json:"id" gorm:"primary_key"`
	RiskRateCode string    `json:"risk_rate_code" csv:"risk_rate_code"`
	Industry     string    `json:"industry" csv:"industry"`
	Category     string    `json:"category" csv:"category"`
	Class        int       `json:"class" csv:"class"`
	CreationDate time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy    string    `json:"created_by" csv:"created_by"`
}

type EducatorBenefitStructure struct {
	ID int `json:"id" gorm:"primary_key"`
	//Year                                int       `json:"year" csv:"year"`
	RiskRateCode                        string    `json:"risk_rate_code" csv:"risk_rate_code"`
	EducatorBenefitCode                 string    `json:"educator_benefit_code" csv:"educator_benefit_code"`
	Grade0MaxTuitionPerYear             float64   `json:"grade0_max_tuition_per_year" csv:"grade0_max_tuition_per_year"`
	Grade0MaxCoverageYears              float64   `json:"grade0_max_coverage_years" csv:"grade0_max_coverage_years"`
	Grade17MaxTuitionPerYear            float64   `json:"grade17_max_tuition_per_year" csv:"grade17_max_tuition_per_year"`
	Grade17MaxCoverageYears             float64   `json:"grade17_max_coverage_years" csv:"grade17_max_coverage_years"`
	Grade812MaxTuitionPerYear           float64   `json:"grade812_max_tuition_per_year" csv:"grade812_max_tuition_per_year"`
	Grade812MaxCoverageYears            float64   `json:"grade812_max_coverage_years" csv:"grade812_max_coverage_years"`
	TertiaryMaxTuitionPerYear           float64   `json:"tertiary_max_tuition_per_year" csv:"tertiary_max_tuition_per_year"`
	TertiaryMaxCoverageYears            float64   `json:"tertiary_max_coverage_years" csv:"tertiary_max_coverage_years"`
	MaxBookAllowanceProportion          float64   `json:"max_book_allowance_proportion" csv:"max_book_allowance_proportion"`
	MaxBookAllowanceAmount              float64   `json:"max_book_allowance_amount" csv:"max_book_allowance_amount"`
	MaxAccommodationAllowanceProportion float64   `json:"max_accommodation_allowance_proportion" csv:"max_accommodation_allowance_proportion"`
	MaxAccommodationAllowanceAmount     float64   `json:"max_accommodation_allowance_amount" csv:"max_accommodation_allowance_amount"`
	CreationDate                        time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy                           string    `json:"created_by" csv:"created_by"`
}

type EducatorRate struct {
	ID int `json:"id" gorm:"primary_key"`
	//Year                  int       `json:"year" csv:"year"`
	RiskRateCode          string    `json:"risk_rate_code" csv:"risk_rate_code"`
	AgeNextBirthday       float64   `json:"age_next_birthday" csv:"age_next_birthday"`
	AverageChildAge       float64   `json:"average_child_age" csv:"average_child_age"`
	AverageNumberChildren float64   `json:"average_number_children" csv:"average_number_children"`
	Grade0RiskRate        float64   `json:"grade0_risk_rate" csv:"grade0_risk_rate"`
	Grade17RiskRate       float64   `json:"grade17_risk_rate" csv:"grade17_risk_rate"`
	Grade812RiskRate      float64   `json:"grade812_risk_rate" csv:"grade812_risk_rate"`
	TertiaryRiskRate      float64   `json:"tertiary_risk_rate" csv:"tertiary_risk_rate"`
	CreationDate          time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy             string    `json:"created_by" csv:"created_by"`
}

type Escalations struct {
	ID                int     `json:"-" gorm:"primary_key"`
	Name              string  `json:"name" csv:"name"`
	EscalationCode    int     `json:"escalation_code" csv:"escalation_code"`
	MinEscalationRate float64 `json:"min_escalation_rate" csv:"min_escalation_rate"`
	MaxEscalationRate float64 `json:"max_escalation_rate" csv:"max_escalation_rate"`
	IsCpiDependent    bool    `json:"is_cpi_dependent" csv:"is_cpi_dependent"`
	//Year              int       `json:"year" csv:"year"`
	//RiskRateCode string    `json:"risk_rate_code" csv:"risk_rate_code"`
	CreationDate time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy    string    `json:"created_by" csv:"created_by"`
}

type FuneralAidsRate struct {
	ID              int       `json:"id" gorm:"primary_key"`
	RiskRateCode    string    `json:"risk_rate_code" csv:"risk_rate_code"`
	AgeNextBirthday int       `json:"age_next_birthday" csv:"age_next_birthday"`
	Gender          string    `json:"gender" csv:"gender"`
	FunAidsQx       float64   `json:"fun_aids_qx" csv:"fun_aids_qx"`
	CreationDate    time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy       string    `json:"created_by" csv:"created_by"`
}

type FuneralRate struct {
	ID              int       `json:"id" gorm:"primary_key"`
	RiskRateCode    string    `json:"risk_rate_code" csv:"risk_rate_code"`
	AgeNextBirthday int       `json:"age_next_birthday" csv:"age_next_birthday"`
	Gender          string    `json:"gender" csv:"gender"`
	FunQx           float64   `json:"fun_qx" csv:"fun_qx"`
	CreationDate    time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy       string    `json:"created_by" csv:"created_by"`
}

type GeneralLoading struct {
	ID                            int       `json:"id" gorm:"primary_key"`
	RiskRateCode                  string    `json:"risk_rate_code" csv:"risk_rate_code"`
	Age                           int       `json:"age" csv:"age"`
	Gender                        string    `json:"gender" csv:"gender"`
	ContigencyLoadingRate         float64   `json:"contigency_loading_rate" csv:"contigency_loading_rate"`
	ContinuationLoadingRate       float64   `json:"continuation_loading_rate" csv:"continuation_loading_rate"`
	TerminalIllnessLoadingRate    float64   `json:"terminal_illness_loading_rate" csv:"terminal_illness_loading_rate"`
	PtdAcceleratedBenefitDiscount float64   `json:"ptd_accelerated_benefit_discount" csv:"ptd_accelerated_benefit_discount"`
	CiAcceleratedBenefitDiscount  float64   `json:"ci_accelerated_benefit_discount" csv:"ci_accelerated_benefit_discount"`
	CreationDate                  time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy                     string    `json:"created_by" csv:"created_by"`
}

type GlaAidsRate struct {
	ID              int       `json:"id" gorm:"primary_key"`
	RiskRateCode    string    `json:"risk_rate_code" csv:"risk_rate_code"`
	AgeNextBirthday int       `json:"age_next_birthday" csv:"age_next_birthday"`
	Gender          string    `json:"gender" csv:"gender"`
	GlaAidsQx       float64   `json:"gla_aids_qx" csv:"gla_aids_qx"`
	CreationDate    time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy       string    `json:"created_by" csv:"created_by"`
}

type RegionLoading struct {
	ID                       int       `json:"id" gorm:"primary_key"`
	RiskRateCode             string    `json:"risk_rate_code" csv:"risk_rate_code"`
	Region                   string    `json:"region" csv:"region"`
	Gender                   string    `json:"gender" csv:"gender"`
	GlaRegionLoadingRate     float64   `json:"gla_region_loading_rate" csv:"gla_region_loading_rate"`
	GlaAidsRegionLoadingRate float64   `json:"gla_aids_region_loading_rate" csv:"gla_aids_region_loading_rate"`
	PtdRegionLoadingRate     float64   `json:"ptd_region_loading_rate" csv:"ptd_region_loading_rate"`
	CiRegionLoadingRate      float64   `json:"ci_region_loading_rate" csv:"ci_region_loading_rate"`
	TtdRegionLoadingRate     float64   `json:"ttd_region_loading_rate" csv:"ttd_region_loading_rate"`
	PhiRegionLoadingRate     float64   `json:"phi_region_loading_rate" csv:"phi_region_loading_rate"`
	FunRegionLoadingRate     float64   `json:"fun_region_loading_rate" csv:"fun_region_loading_rate"`
	FunAidsRegionLoadingRate float64   `json:"fun_aids_region_loading_rate" csv:"fun_aids_region_loading_rate"`
	CreationDate             time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy                string    `json:"created_by" csv:"created_by"`
}

type TieredIncomeReplacement struct {
	ID                     int       `json:"id" gorm:"primary_key"`
	RiskRateCode           string    `json:"risk_rate_code" csv:"risk_rate_code"`
	AnnualLowerBound       float64   `json:"annual_lower_bound" csv:"annual_lower_bound"`
	AnnualUpperBound       float64   `json:"annual_upper_bound" csv:"annual_upper_bound"`
	IncomeReplacementRatio float64   `json:"income_replacement_ratio" csv:"income_replacement_ratio"`
	CreationDate           time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy              string    `json:"created_by" csv:"created_by"`
}

type CustomTieredIncomeReplacement struct {
	ID                     int       `json:"id" gorm:"primaryKey;autoIncrement"`
	SchemeName             string    `json:"scheme_name" csv:"scheme_name" gorm:"index:idx_custom_tir_scheme_rrc"`
	RiskRateCode           string    `json:"risk_rate_code" csv:"risk_rate_code" gorm:"index:idx_custom_tir_scheme_rrc"`
	AnnualLowerBound       float64   `json:"annual_lower_bound" csv:"annual_lower_bound"`
	AnnualUpperBound       float64   `json:"annual_upper_bound" csv:"annual_upper_bound"`
	IncomeReplacementRatio float64   `json:"income_replacement_ratio" csv:"income_replacement_ratio"`
	RequestedBy            string    `json:"requested_by"`
	CreationDate           time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy              string    `json:"created_by" csv:"created_by"`
}

type DiscountAuthority struct {
	ID           int       `json:"id" gorm:"primary_key"`
	RiskRateCode string    `json:"risk_rate_code" csv:"risk_rate_code"`
	Role         string    `json:"role" csv:"role"`
	MaxDiscount  float64   `json:"max_discount" csv:"max_discount"`
	CreationDate time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy    string    `json:"created_by" csv:"created_by"`
}

type Restriction struct {
	ID                                  int       `json:"id" gorm:"primary_key"`
	RiskRateCode                        string    `json:"risk_rate_code" csv:"risk_rate_code"`
	SevereIllnessMaximumBenefit         float64   `json:"severe_illness_maximum_benefit" csv:"severe_illness_maximum_benefit"`
	SpouseGlaMaximumBenefit             float64   `json:"spouse_gla_maximum_benefit" csv:"spouse_gla_maximum_benefit"`
	PhiMaximumMonthlyBenefit            float64   `json:"phi_maximum_monthly_benefit" csv:"phi_maximum_monthly_benefit"`
	PhiMaximumMonthlyContributionWaiver float64   `json:"phi_maximum_monthly_contribution_waiver" csv:"phi_maximum_monthly_contribution_waiver"`
	TtdMaximumMonthlyBenefit            float64   `json:"ttd_maximum_monthly_benefit" csv:"ttd_maximum_monthly_benefit"`
	MaxMedicalAidWaiver                 float64   `json:"max_medical_aid_waiver" csv:"max_medical_aid_waiver"`
	MinEntryAge                         int       `json:"min_entry_age" csv:"min_entry_age"`
	MaxEntryAge                         int       `json:"max_entry_age" csv:"max_entry_age"`
	GlaMaxCoverAge                      int       `json:"gla_max_cover_age" csv:"gla_max_cover_age"`
	PtdMaxCoverAge                      int       `json:"ptd_max_cover_age" csv:"ptd_max_cover_age"`
	CiMaxCoverAge                       int       `json:"ci_max_cover_age" csv:"ci_max_cover_age"`
	PhiMaxCoverAge                      int       `json:"phi_max_cover_age" csv:"phi_max_cover_age"`
	TtdMaxCoverAge                      int       `json:"ttd_max_cover_age" csv:"ttd_max_cover_age"`
	FunMaxCoverAge                      int       `json:"fun_max_cover_age" csv:"fun_max_cover_age"`
	CreationDate                        time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy                           string    `json:"created_by" csv:"created_by"`
}

type ReinsuranceGlaRate struct {
	ID              int       `json:"id" gorm:"primary_key"`
	RiskRateCode    string    `json:"risk_rate_code" csv:"risk_rate_code"`
	AgeNextBirthday int       `json:"age_next_birthday" csv:"age_next_birthday"`
	IncomeLevel     string    `json:"income_level" csv:"income_level"`
	Gender          string    `json:"gender" csv:"gender"`
	WaitingPeriod   int       `json:"waiting_period" csv:"waiting_period"`
	ReQx            float64   `json:"re_qx" csv:"re_qx"`
	CreationDate    time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy       string    `json:"created_by" csv:"created_by"`
}

type ReinsuranceCiRate struct {
	ID                int       `json:"id" gorm:"primary_key"`
	RiskRateCode      string    `json:"risk_rate_code" csv:"risk_rate_code"`
	AgeNextBirthday   int       `json:"age_next_birthday" csv:"age_next_birthday"`
	Gender            string    `json:"gender" csv:"gender"`
	OccupationClass   string    `json:"occupation_class" csv:"occupation_class"`
	IncomeLevel       string    `json:"income_level" csv:"income_level"`
	WaitingPeriod     int       `json:"waiting_period" csv:"waiting_period"`
	BenefitDefinition string    `json:"benefit_definition" csv:"benefit_definition"`
	CiRate            float64   `json:"ci_rate" csv:"ci_rate"`
	CreationDate      time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy         string    `json:"created_by" csv:"created_by"`
}

type ReinsurancePtdRate struct {
	ID                   int       `json:"id" gorm:"primary_key"`
	RiskRateCode         string    `json:"risk_rate_code" csv:"risk_rate_code"`
	RiskType             string    `json:"risk_type" csv:"risk_type"`
	AgeNextBirthday      int       `json:"age_next_birthday" csv:"age_next_birthday"`
	Gender               string    `json:"gender" csv:"gender"`
	OccupationClass      string    `json:"occupation_class" csv:"occupation_class"`
	IncomeLevel          string    `json:"income_level" csv:"income_level"`
	WaitingPeriod        int       `json:"waiting_period" csv:"waiting_period"`
	DeferredPeriod       int       `json:"deferred_period" csv:"deferred_period"`
	DisabilityDefinition string    `json:"disability_definition" csv:"disability_definition"`
	PtdRate              float64   `json:"ptd_rate" csv:"ptd_rate"`
	CreationDate         time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy            string    `json:"created_by" csv:"created_by"`
}

type ReinsurancePhiRate struct {
	ID                      int       `json:"id" gorm:"primary_key"`
	RiskRateCode            string    `json:"risk_rate_code" csv:"risk_rate_code"`
	RiskType                string    `json:"risk_type" csv:"risk_type"`
	AgeNextBirthday         int       `json:"age_next_birthday" csv:"age_next_birthday"`
	Gender                  string    `json:"gender" csv:"gender"`
	OccupationClass         string    `json:"occupation_class" csv:"occupation_class"`
	IncomeLevel             string    `json:"income_level" csv:"income_level"`
	WaitingPeriod           int       `json:"waiting_period" csv:"waiting_period"`
	DeferredPeriod          int       `json:"deferred_period" csv:"deferred_period"`
	NormalRetirementAge     int       `json:"normal_retirement_age" csv:"normal_retirement_age"`
	BenefitEscalationOption string    `json:"benefit_escalation_option" csv:"benefit_escalation_option"`
	DisabilityDefinition    string    `json:"disability_definition" csv:"disability_definition"`
	PhiRate                 float64   `json:"phi_rate" csv:"phi_rate"`
	CreationDate            time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy               string    `json:"created_by" csv:"created_by"`
}

type PremiumLoading struct {
	ID                    int       `json:"id" gorm:"primary_key"`
	RiskRateCode          string    `json:"risk_rate_code" csv:"risk_rate_code"`
	Channel               string    `json:"channel" csv:"channel"`
	SchemeSizeLevel       int       `json:"scheme_size_level" csv:"scheme_size_level"`
	CommissionLoading     float64   `json:"commission_loading" csv:"commission_loading"`
	ExpenseLoading        float64   `json:"expense_loading" csv:"expense_loading"`
	AdminLoading          float64   `json:"admin_loading" csv:"admin_loading"`
	OtherLoading          float64   `json:"other_loading" csv:"other_loading"`
	ProfitLoading         float64   `json:"profit_loading" csv:"profit_loading"`
	MinimumPremiumLoading float64   `json:"minimum_premium_loading" csv:"minimum_premium_loading" gorm:"not null;default:0"`
	CreationDate          time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy             string    `json:"created_by" csv:"created_by"`
}

type SchemeSizeLevel struct {
	ID           int       `json:"id" gorm:"primary_key"`
	RiskRateCode string    `json:"risk_rate_code" csv:"risk_rate_code"`
	MinCount     int       `json:"min_count" csv:"min_count"`
	MaxCount     int       `json:"max_count" csv:"max_count"`
	SizeLevel    int       `json:"size_level" csv:"size_level"`
	CreationDate time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy    string    `json:"created_by" csv:"created_by"`
}

type TaxTable struct {
	ID           int       `json:"id" gorm:"primary_key"`
	RiskRateCode string    `json:"risk_rate_code" csv:"risk_rate_code"`
	Level        int       `json:"level" csv:"level"`
	Min          float64   `json:"min" csv:"min"`
	Max          float64   `json:"max" csv:"max"`
	TaxRate      float64   `json:"tax_rate" csv:"tax_rate"`
	CreationDate time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy    string    `json:"created_by" csv:"created_by"`
}

type Bordereaux struct {
	ID                            int        `json:"-" gorm:"primaryKey;autoIncrement"`
	SchemeId                      int        `json:"-" csv:"scheme_id"`
	QuoteId                       int        `json:"-" csv:"quote_id" gorm:"index"`
	Month                         time.Month `json:"-" csv:"month" gorm:"index"`
	Year                          int        `json:"-" csv:"year" gorm:"index"`
	Period                        string     `json:"-" csv:"period" gorm:"index"`
	Category                      string     `json:"category" csv:"category"`
	EmployeeNumber                string     `json:"employee_number" csv:"employee_number"`
	MemberIdNumber                string     `json:"member_id_number" csv:"member_id_number"`
	MemberName                    string     `json:"member_name" csv:"member_name"`
	Gender                        string     `json:"gender" csv:"gender"`
	DateOfBirth                   time.Time  `json:"date_of_birth" csv:"date_of_birth"`
	AgeNextBirthday               float64    `json:"age_next_birthday" csv:"age_next_birthday"`
	EntryDate                     time.Time  `json:"entry_date" csv:"entry_date" gorm:"autoCreateTime"`
	ExitDate                      time.Time  `json:"exit_date" csv:"exit_date" gorm:"autoCreateTime"`
	IsOriginalMember              bool       `json:"is_original_member" csv:"is_original_member"`
	RenewalDate                   string     `json:"renewal_date" csv:"renewal_date"`
	Currency                      string     `json:"currency" csv:"currency"`
	AnnualSalary                  float64    `json:"annual_salary" csv:"annual_salary"`
	Industry                      string     `json:"industry" csv:"industry"`
	IndustryClass                 string     `json:"industry_class" csv:"industry_class"`
	GlaMultiple                   float64    `json:"gla_multiple" csv:"gla_multiple"`
	GlaCoveredSumAssured          float64    `json:"gla_covered_sum_assured" csv:"gla_covered_sum_assured"`
	GlaRetainedSumAssured         float64    `json:"gla_retained_sum_assured" csv:"gla_retained_sum_assured"`
	GlaCededSumAssured            float64    `json:"gla_ceded_sum_assured" csv:"gla_ceded_sum_assured"`
	LoadedGlaRiskRate             float64    `json:"loaded_gla_risk_rate" csv:"loaded_gla_risk_rate"`
	ExpAdjLoadedGlaRiskRate       float64    `json:"exp_adj_loaded_gla_risk_rate" csv:"exp_adj_loaded_gla_risk_rate"`
	GlaRetainedRiskPremium        float64    `json:"gla_retained_risk_premium" csv:"gla_retained_risk_premium"`
	GlaCededRiskPremium           float64    `json:"gla_ceded_risk_premium" csv:"gla_ceded_risk_premium"`
	GlaAnnualPremium              float64    `json:"gla_annual_premium" csv:"gla_annual_premium"`
	PtdMultiple                   float64    `json:"ptd_multiple" csv:"ptd_multiple"`
	PtdCoveredSumAssured          float64    `json:"ptd_covered_sum_assured" csv:"ptd_covered_sum_assured"`
	PtdRetainedSumAssured         float64    `json:"ptd_retained_sum_assured" csv:"ptd_retained_sum_assured"`
	PtdCededSumAssured            float64    `json:"ptd_ceded_sum_assured" csv:"ptd_ceded_sum_assured"`
	LoadedPtdRiskRate             float64    `json:"loaded_ptd_risk_rate" csv:"loaded_ptd_risk_rate"`
	ExpAdjLoadedPtdRiskRate       float64    `json:"exp_adj_loaded_ptd_risk_rate" csv:"exp_adj_loaded_ptd_risk_rate"`
	PtdRetainedRiskPremium        float64    `json:"ptd_retained_risk_premium" csv:"ptd_retained_risk_premium"`
	PtdCededRiskPremium           float64    `json:"ptd_ceded_risk_premium" csv:"ptd_ceded_risk_premium"`
	PtdAnnualPremium              float64    `json:"ptd_annual_premium" csv:"ptd_annual_premium"`
	CiMultiple                    float64    `json:"ci_multiple" csv:"ci_multiple"`
	CiCoveredSumAssured           float64    `json:"ci_covered_sum_assured" csv:"ci_covered_sum_assured"`
	CiRetainedSumAssured          float64    `json:"ci_retained_sum_assured" csv:"ci_retained_sum_assured"`
	CiCededSumAssured             float64    `json:"ci_ceded_sum_assured" csv:"ci_ceded_sum_assured"`
	LoadedCiRiskRate              float64    `json:"loaded_ci_risk_rate" csv:"loaded_ci_risk_rate"`
	ExpAdjLoadedCiRiskRate        float64    `json:"exp_adj_loaded_ci_risk_rate" csv:"exp_adj_loaded_ci_risk_rate"`
	CiRetainedRiskPremium         float64    `json:"ci_retained_risk_premium" csv:"ci_retained_risk_premium"`
	CiCededRiskPremium            float64    `json:"ci_ceded_risk_premium" csv:"ci_ceded_risk_premium"`
	CiAnnualPremium               float64    `json:"ci_annual_premium" csv:"ci_annual_premium"`
	SglaMultiple                  float64    `json:"sgla_multiple" csv:"sgla_multiple"`
	SglaCoveredSumAssured         float64    `json:"sgla_covered_sum_assured" csv:"sgla_covered_sum_assured"`
	SglaRetainedSumAssured        float64    `json:"sgla_retained_sum_assured" csv:"sgla_retained_sum_assured"`
	SglaCededSumAssured           float64    `json:"sgla_ceded_sum_assured" csv:"sgla_ceded_sum_assured"`
	LoadedSglaRiskRate            float64    `json:"loaded_sgla_risk_rate" csv:"loaded_sgla_risk_rate"`
	ExpAdjLoadedSglaRiskRate      float64    `json:"exp_adj_loaded_sgla_risk_rate" csv:"exp_adj_loaded_sgla_risk_rate"`
	SglaRetainedRiskPremium       float64    `json:"sgla_retained_risk_premium" csv:"sgla_retained_risk_premium"`
	SglaCededRiskPremium          float64    `json:"sgla_ceded_risk_premium" csv:"sgla_ceded_risk_premium"`
	SglaAnnualPremium             float64    `json:"sgla_annual_premium" csv:"sgla_annual_premium"`
	TtdReplacementMultiple        float64    `json:"ttd_replacement_multiple" csv:"ttd_replacement_multiple"`
	TtdMonthlyBenefit             float64    `json:"ttd_monthly_benefit" csv:"ttd_monthly_benefit"`
	TtdRetainedMonthlyBenefit     float64    `json:"ttd_retained_monthly_benefit" csv:"ttd_retained_monthly_benefit"`
	TtdCededMonthlyBenefit        float64    `json:"ttd_ceded_monthly_benefit" csv:"ttd_ceded_monthly_benefit"`
	LoadedTtdRiskRate             float64    `json:"loaded_ttd_risk_rate" csv:"loaded_ttd_risk_rate"`
	ExpAdjLoadedTtdRiskRate       float64    `json:"exp_adj_loaded_ttd_risk_rate" csv:"exp_adj_loaded_ttd_risk_rate"`
	TtdRetainedRiskPremium        float64    `json:"ttd_retained_risk_premium" csv:"ttd_retained_risk_premium"`
	TtdCededRiskPremium           float64    `json:"ttd_ceded_risk_premium" csv:"ttd_ceded_risk_premium"`
	PhiReplacementMultiple        float64    `json:"phi_replacement_multiple" csv:"phi_replacement_multiple"`
	PhiMonthlyBenefit             float64    `json:"phi_monthly_benefit" csv:"phi_monthly_benefit"`
	PhiRetainedMonthlyBenefit     float64    `json:"phi_retained_monthly_benefit" csv:"phi_retained_monthly_benefit"`
	PhiCededMonthlyBenefit        float64    `json:"phi_ceded_monthly_benefit" csv:"phi_ceded_monthly_benefit"`
	LoadedPhiRiskRate             float64    `json:"loaded_phi_risk_rate" csv:"loaded_phi_risk_rate"`
	ExpAdjLoadedPhiRiskRate       float64    `json:"exp_adj_loaded_phi_risk_rate" csv:"exp_adj_loaded_phi_risk_rate"`
	PhiRetainedRiskPremium        float64    `json:"phi_retained_risk_premium" csv:"phi_retained_risk_premium"`
	PhiCededRiskPremium           float64    `json:"phi_ceded_risk_premium" csv:"phi_ceded_risk_premium"`
	MainMemberFuneralSumAssured   float64    `json:"main_member_funeral_sum_assured" csv:"main_member_funeral_sum_assured"`
	MainMemberRetainedSumAssured  float64    `json:"main_member_retained_sum_assured" csv:"main_member_retained_sum_assured"`
	MainMemberCededSumAssured     float64    `json:"main_member_ceded_sum_assured" csv:"main_member_ceded_sum_assured"`
	MainMemberRiskRate            float64    `json:"main_member_risk_rate" csv:"main_member_risk_rate"`
	MainMemberRetainedRiskPremium float64    `json:"main_member_retained_risk_premium" csv:"main_member_retained_risk_premium"`
	MainMemberCededRiskPremium    float64    `json:"main_member_ceded_risk_premium" csv:"main_member_ceded_risk_premium"`
	SpouseFuneralSumAssured       float64    `json:"spouse_funeral_sum_assured" csv:"spouse_funeral_sum_assured"`
	SpouseRetainedSumAssured      float64    `json:"spouse_retained_sum_assured" csv:"spouse_retained_sum_assured"`
	SpouseCededSumAssured         float64    `json:"spouse_ceded_sum_assured" csv:"spouse_ceded_sum_assured"`
	SpouseRiskRate                float64    `json:"spouse_risk_rate" csv:"spouse_risk_rate"`
	SpouseRetainedRiskPremium     float64    `json:"spouse_retained_risk_premium" csv:"spouse_retained_risk_premium"`
	SpouseCededRiskPremium        float64    `json:"spouse_ceded_risk_premium" csv:"spouse_ceded_risk_premium"`
	ChildFuneralSumAssured        float64    `json:"child_funeral_sum_assured" csv:"child_funeral_sum_assured"`
	ChildRetainedSumAssured       float64    `json:"child_retained_sum_assured" csv:"child_retained_sum_assured"`
	ChildCededSumAssured          float64    `json:"child_ceded_sum_assured" csv:"child_ceded_sum_assured"`
	ChildRiskRate                 float64    `json:"child_risk_rate" csv:"child_risk_rate"`
	ChildRetainedRiskPremium      float64    `json:"child_retained_risk_premium" csv:"child_retained_risk_premium"`
	ChildCededRiskPremium         float64    `json:"child_ceded_risk_premium" csv:"child_ceded_risk_premium"`
	ParentFuneralSumAssured       float64    `json:"parent_funeral_sum_assured" csv:"parent_funeral_sum_assured"`
	ParentRetainedSumAssured      float64    `json:"parent_retained_sum_assured" csv:"parent_retained_sum_assured"`
	ParentCededSumAssured         float64    `json:"parent_ceded_sum_assured" csv:"parent_ceded_sum_assured"`
	ParentRiskRate                float64    `json:"parent_risk_rate" csv:"parent_risk_rate"`
	ParentRetainedRiskPremium     float64    `json:"parent_retained_risk_premium" csv:"parent_retained_risk_premium"`
	ParentCededRiskPremium        float64    `json:"parent_ceded_risk_premium" csv:"parent_ceded_risk_premium"`
	DependantFuneralSumAssured    float64    `json:"dependant_funeral_sum_assured" csv:"dependant_funeral_sum_assured"`
	DependantRetainedSumAssured   float64    `json:"dependant_retained_sum_assured" csv:"dependant_retained_sum_assured"`
	DependantCededSumAssured      float64    `json:"dependant_ceded_sum_assured" csv:"dependant_ceded_sum_assured"`
	DependantRiskRate             float64    `json:"dependant_risk_rate" csv:"dependant_risk_rate"`
	DependantRetainedRiskPremium  float64    `json:"dependant_retained_risk_premium" csv:"dependant_retained_risk_premium"`
	DependantCededRiskPremium     float64    `json:"dependant_ceded_risk_premium" csv:"dependant_ceded_risk_premium"`
	CreationDate                  time.Time  `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy                     string     `json:"created_by" csv:"created_by"`
}

type CardData struct {
	Title    string      `json:"title"`
	Value    interface{} `json:"value"`
	DataType string      `json:"data_type"`
	Flex     int         `json:"flex"`
}

type RevenueBenefit struct {
	Benefit string  `json:"type"`
	Revenue float64 `json:"revenue"`
	Claims  float64 `json:"claims"`
}
type GroupPricingIncomeComponent struct {
	Benefit  string  `json:"type"`
	Expected float64 `json:"expected"`
	Actual   float64 `json:"actual"`
}

type GroupSchemeExposure struct {
	SchemeName       string    `json:"scheme_name"`
	QuoteId          int       `json:"-" csv:"quote_id"`
	FinancialYear    int       `json:"financial_year"`
	Industry         string    `json:"industry"`
	Benefit          string    `json:"benefit"`
	AgeBand          string    `json:"age_band"`
	MinAge           int       `json:"min_age"`
	MaxAge           int       `json:"max_age"`
	MaleSumAssured   float64   `json:"male_sum_assured"`
	FemaleSumAssured float64   `json:"female_sum_assured"`
	TotalSumAssured  float64   `json:"total_sum_assured"`
	QuoteStatus      string    `json:"quote_status" csv:"quote_status"`
	CreationDate     time.Time `json:"creation_date" csv:"creation_date" gorm:"autoCreateTime"`
	CreatedBy        string    `json:"created_by" csv:"created_by"`
}

type GroupPricingAgeBands struct {
	ID     int    `json:"id" gorm:"primary_key"`
	Name   string `json:"name"`
	MinAge int    `json:"min_age"`
	MaxAge int    `json:"max_age"`
}

type GroupBusinessBenefits struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

type GroupSchemeClaim struct {
	ID                       int                             `json:"id" gorm:"primary_key"`
	ClaimNumber              string                          `json:"claim_number"`
	MemberIDNumber           string                          `json:"member_id_number"`
	MemberName               string                          `json:"member_name"`
	SchemeId                 int                             `json:"scheme_id"`
	SchemeName               string                          `json:"scheme_name"`
	BenefitAlias             string                          `json:"benefit_alias"`
	BenefitName              string                          `json:"benefit_name"`
	BenefitCode              string                          `json:"benefit_code"`
	MemberType               string                          `json:"member_type"`
	DateOfEvent              string                          `json:"date_of_event"`
	DateNotified             string                          `json:"date_notified"`
	CauseType                string                          `json:"cause_type"`
	ClaimAmount              float64                         `json:"claim_amount"`
	Priority                 string                          `json:"priority"`
	ClaimantName             string                          `json:"claimant_name"`
	ClaimantIDNumber         string                          `json:"claimant_id_number"`
	RelationshipToMember     string                          `json:"relationship_to_member"`
	ClaimantContactNumber    string                          `json:"claimant_contact_number"`
	BankName                 string                          `json:"bank_name"`
	BankBranchCode           string                          `json:"bank_branch_code"`
	BankAccountNumber        string                          `json:"bank_account_number"`
	BankAccountType          string                          `json:"bank_account_type"`
	AccountHolderName        string                          `json:"account_holder_name"`
	SupportingDocuments      []SupportingDocument            `json:"supporting_documents" gorm:"-"`
	Description              string                          `json:"description"`
	Status                   string                          `json:"status"`
	DateRegistered           string                          `json:"date_registered"`
	MissingRequiredDocuments StringArray                     `json:"missing_required_documents" gorm:"type:json"`
	CreationDate             time.Time                       `json:"creation_date" gorm:"autoCreateTime"`
	CreatedBy                string                          `json:"created_by"`
	Attachments              []GroupSchemeClaimAttachment    `json:"attachments" gorm:"foreignKey:ClaimID;references:ID"`
	Assessments              []GroupSchemeClaimAssessment    `json:"assessments" gorm:"foreignKey:ClaimID;references:ID"`
	Communications           []GroupSchemeClaimCommunication `json:"communications" gorm:"foreignKey:ClaimID;references:ID"`
	// Declines keeps records of claim decline decisions
	Declines []GroupSchemeClaimDecline `json:"declines" gorm:"foreignKey:ClaimID;references:ID"`
	// StatusAudits keeps the history of status changes for this claim
	StatusAudits []GroupSchemeClaimStatusAudit `json:"status_audits" gorm:"foreignKey:ClaimID;references:ID"`
}

type SupportingDocument struct {
	DocumentType string `json:"document_type"`
	DocumentName string `json:"document_name"`
	FileIndex    int    `json:"file_index"`
}

// GroupSchemeClaimAttachment stores metadata for uploaded claim supporting documents
type GroupSchemeClaimAttachment struct {
	ID           int       `json:"id" gorm:"primary_key"`
	ClaimID      int       `json:"claim_id" gorm:"index"`
	DocumentType string    `json:"document_type"`
	DocumentName string    `json:"document_name"`
	FileName     string    `json:"filename"`
	ContentType  string    `json:"content_type"`
	SizeBytes    int64     `json:"size_bytes"`
	StoragePath  string    `json:"storage_path"`
	UploadedAt   time.Time `json:"uploaded_at" gorm:"autoCreateTime"`
	UploadedBy   string    `json:"uploaded_by"`
	// ViewerURL is a computed field providing a direct URL that a document viewer
	// (e.g., iframe) can use to access the attachment. Not persisted to DB.
	ViewerURL string `json:"viewer_url" gorm:"-"`
}

// JSONMapBool is a JSON-serializable map[string]bool that implements
// the sql driver Valuer and Scanner interfaces so GORM can persist it.
type JSONMapBool map[string]bool

// Value implements driver.Valuer. It marshals the map to JSON for storage.
func (m JSONMapBool) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	b, err := json.Marshal(map[string]bool(m))
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan implements sql.Scanner. It unmarshals JSON from the DB into the map.
func (m *JSONMapBool) Scan(value any) error {
	if value == nil {
		*m = nil
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type for JSONMapBool Scan: %T", value)
	}
	if len(data) == 0 {
		*m = nil
		return nil
	}
	var tmp map[string]bool
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*m = JSONMapBool(tmp)
	return nil
}

// GroupSchemeClaimAssessment stores assessment details for a claim
type GroupSchemeClaimAssessment struct {
	ID                    int         `json:"id" gorm:"primary_key"`
	ClaimID               int         `json:"claim_id" gorm:"index;not null"`
	AssessorName          string      `json:"assessor_name"`
	AssessmentDate        string      `json:"assessment_date"`
	AssessmentOutcome     string      `json:"assessment_outcome"`
	RecommendedAmount     float64     `json:"recommended_amount"`
	MedicalOfficer        string      `json:"medical_officer"`
	MedicalAssessmentDate *string     `json:"medical_assessment_date"`
	DisabilityPercentage  *string     `json:"disability_percentage"`
	MedicalCondition      string      `json:"medical_condition"`
	MedicalNotes          string      `json:"medical_notes"`
	DocumentsVerified     JSONMapBool `json:"documents_verified" gorm:"type:json"`
	FraudRiskLevel        string      `json:"fraud_risk_level"`
	RequiresInvestigation bool        `json:"requires_investigation"`
	RiskNotes             string      `json:"risk_notes"`
	AssessmentNotes       string      `json:"assessment_notes"`
	NextActions           string      `json:"next_actions"`
	Checklist             JSONMapBool `json:"checklist" gorm:"type:json"`
	AssessmentTimestamp   *time.Time  `json:"assessment_timestamp"`
	CreatedBy             string      `json:"created_by"`
	CreationDate          time.Time   `json:"creation_date" gorm:"autoCreateTime"`
	UpdatedAt             time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

// GroupSchemeClaimCommunication stores communication logs between assessors and claimant
type GroupSchemeClaimCommunication struct {
	ID        int        `json:"id" gorm:"primary_key"`
	ClaimID   int        `json:"claim_id" gorm:"index;not null"`
	Method    string     `json:"method"` // 'Email', 'Phone', 'SMS', 'Letter', 'In Person'
	Message   string     `json:"message"`
	Timestamp *time.Time `json:"timestamp"` // Optional time of the communication event (could differ from created_at)
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy string     `json:"created_by"`
}

// GroupSchemeClaimDecline stores decline records for a claim
// Table name will be inferred by GORM as group_scheme_claim_declines
type GroupSchemeClaimDecline struct {
	ID                         int        `json:"id" gorm:"primary_key"`
	ClaimID                    int        `json:"claim_id" gorm:"index;not null"`
	PrimaryReason              string     `json:"primary_reason"`
	DetailedReason             string     `json:"detailed_reason"`
	AssessmentReference        string     `json:"assessment_reference"`
	RequiresMemberNotification bool       `json:"requires_member_notification"`
	InternalNotes              string     `json:"internal_notes"`
	DeclinedBy                 string     `json:"declined_by"`
	DeclinedAt                 *time.Time `json:"declined_at"`
	CreatedAt                  time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

// GroupSchemeClaimStatusAudit keeps a history of status changes for a GroupSchemeClaim
// Table name will be inferred by GORM as group_scheme_claim_status_audits
type GroupSchemeClaimStatusAudit struct {
	ID            int       `json:"id" gorm:"primary_key"`
	ClaimID       int       `json:"claim_id" gorm:"index"`
	OldStatus     string    `json:"old_status"`
	NewStatus     string    `json:"new_status"`
	StatusMessage string    `json:"status_message"`
	ChangedBy     string    `json:"changed_by"`
	ChangedAt     time.Time `json:"changed_at" gorm:"autoCreateTime"`
}

// ClaimPaymentSchedule represents a batch of approved claims collected for payment disbursement.
// After creation the included claims are moved to "submitted_for_payment".
// Once a proof of payment is uploaded and confirmed, claims move to "paid".
// Table name: claim_payment_schedules
type ClaimPaymentSchedule struct {
	ID             int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ScheduleNumber string `json:"schedule_number" gorm:"type:varchar(191);uniqueIndex"`
	Description    string `json:"description"`
	// Status lifecycle: draft → submitted → confirmed
	Status           string                     `json:"status"`
	TotalAmount      float64                    `json:"total_amount"`
	ClaimsCount      int                        `json:"claims_count"`
	ExportedAt       *time.Time                 `json:"exported_at"`
	ExportedBy       string                     `json:"exported_by"`
	ACBFileGenerated bool                       `json:"acb_file_generated"`
	ACBGeneratedAt   *time.Time                 `json:"acb_generated_at"`
	ACBGeneratedBy   string                     `json:"acb_generated_by"`
	BankProfileID    *int                       `json:"bank_profile_id"`
	CreatedBy        string                     `json:"created_by"`
	CreatedAt        time.Time                  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time                  `json:"updated_at" gorm:"autoUpdateTime"`
	Items            []ClaimPaymentScheduleItem `json:"items" gorm:"foreignKey:ScheduleID;references:ID"`
	ProofOfPayments  []ClaimPaymentProof        `json:"proof_of_payments" gorm:"foreignKey:ScheduleID;references:ID"`
}

// ClaimPaymentScheduleItem links a single claim to a payment schedule.
// Table name: claim_payment_schedule_items
type ClaimPaymentScheduleItem struct {
	ID                int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ScheduleID        int       `json:"schedule_id" gorm:"index;not null"`
	ClaimID           int       `json:"claim_id" gorm:"index;not null"`
	ClaimNumber       string    `json:"claim_number"`
	MemberName        string    `json:"member_name"`
	MemberIDNumber    string    `json:"member_id_number"`
	BenefitName       string    `json:"benefit_name"`
	SchemeName        string    `json:"scheme_name"`
	SchemeID          int       `json:"scheme_id"`
	ClaimAmount       float64   `json:"claim_amount"`
	BankName          string    `json:"bank_name"`
	BankBranchCode    string    `json:"bank_branch_code"`
	BankAccountNumber string    `json:"bank_account_number"`
	BankAccountType   string    `json:"bank_account_type"`
	AccountHolderName string    `json:"account_holder_name"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// ACBBankProfile stores company-level config for ACB file generation. One per source bank account.
type ACBBankProfile struct {
	ID                int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ProfileName       string    `json:"profile_name"`
	BankName          string    `json:"bank_name"`
	UserCode          string    `json:"user_code"`
	UserBranchCode    string    `json:"user_branch_code"`
	UserAccountNumber string    `json:"user_account_number"`
	UserAccountType   string    `json:"user_account_type"`
	BankTypeCode      string    `json:"bank_type_code"`
	ServiceType       string    `json:"service_type"`
	GenerationNumber  int       `json:"generation_number"`
	IsActive          bool      `json:"is_active" gorm:"default:true"`
	CreatedBy         string    `json:"created_by"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// ACBFileRecord is an audit trail for every generated ACB file.
type ACBFileRecord struct {
	ID               int                  `json:"id" gorm:"primaryKey;autoIncrement"`
	ScheduleID       int                  `json:"schedule_id" gorm:"index"`
	BankProfileID    int                  `json:"bank_profile_id"`
	FileName         string               `json:"file_name"`
	FilePath         string               `json:"file_path"`
	ActionDate       string               `json:"action_date"`
	TransactionCount int                  `json:"transaction_count"`
	TotalAmount      float64              `json:"total_amount"`
	HashTotal        int64                `json:"hash_total" gorm:"type:bigint"`
	GenerationNumber int                  `json:"generation_number"`
	Status           string               `json:"status"`
	IsRetry          bool                 `json:"is_retry"`
	GeneratedBy      string               `json:"generated_by"`
	GeneratedAt      time.Time            `json:"generated_at" gorm:"autoCreateTime"`
	ReconciledAt     *time.Time           `json:"reconciled_at"`
	ReconciledBy     string               `json:"reconciled_by"`
	Schedule         ClaimPaymentSchedule `json:"schedule" gorm:"foreignKey:ScheduleID;references:ID"`
	BankProfile      ACBBankProfile       `json:"bank_profile" gorm:"foreignKey:BankProfileID;references:ID"`
}

// ACBReconciliationResult stores per-transaction outcome from bank response parsing.
type ACBReconciliationResult struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ACBFileID      int       `json:"acb_file_id" gorm:"index"`
	ScheduleItemID int       `json:"schedule_item_id" gorm:"index"`
	ClaimID        int       `json:"claim_id" gorm:"index"`
	ClaimNumber    string    `json:"claim_number"`
	AccountNumber  string    `json:"account_number"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"`
	FailureReason  string    `json:"failure_reason"`
	BankReference  string    `json:"bank_reference"`
	ResponseCode   string    `json:"response_code"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// CreateBankProfileRequest is the inbound payload for creating a bank profile.
type CreateBankProfileRequest struct {
	ProfileName       string `json:"profile_name"`
	BankName          string `json:"bank_name"`
	UserCode          string `json:"user_code"`
	UserBranchCode    string `json:"user_branch_code"`
	UserAccountNumber string `json:"user_account_number"`
	UserAccountType   string `json:"user_account_type"`
	BankTypeCode      string `json:"bank_type_code"`
	ServiceType       string `json:"service_type"`
}

// UpdateBankProfileRequest is the inbound payload for updating a bank profile.
type UpdateBankProfileRequest struct {
	ProfileName       *string `json:"profile_name"`
	BankName          *string `json:"bank_name"`
	UserCode          *string `json:"user_code"`
	UserBranchCode    *string `json:"user_branch_code"`
	UserAccountNumber *string `json:"user_account_number"`
	UserAccountType   *string `json:"user_account_type"`
	BankTypeCode      *string `json:"bank_type_code"`
	ServiceType       *string `json:"service_type"`
	IsActive          *bool   `json:"is_active"`
}

// GenerateACBRequest is the inbound payload for generating an ACB file.
type GenerateACBRequest struct {
	BankProfileID int    `json:"bank_profile_id"`
	ActionDate    string `json:"action_date"`
}

// RetryFailedRequest is the inbound payload for retrying failed payments.
type RetryFailedRequest struct {
	ItemIDs []int `json:"item_ids"`
}

// ACBReconciliationSummary is the response for ACB reconciliation results.
type ACBReconciliationSummary struct {
	TotalTransactions int     `json:"total_transactions"`
	Paid              int     `json:"paid"`
	Failed            int     `json:"failed"`
	Unmatched         int     `json:"unmatched"`
	TotalPaid         float64 `json:"total_paid"`
	TotalFailed       float64 `json:"total_failed"`
}

// ClaimPaymentProof stores uploaded proof-of-payment documents for a payment schedule.
// Once uploaded the schedule is marked "confirmed" and all its claims move to "paid".
// Table name: claim_payment_proofs
type ClaimPaymentProof struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ScheduleID  int       `json:"schedule_id" gorm:"index;not null"`
	FileName    string    `json:"file_name"`
	ContentType string    `json:"content_type"`
	SizeBytes   int64     `json:"size_bytes"`
	StoragePath string    `json:"storage_path"`
	Notes       string    `json:"notes"`
	UploadedBy  string    `json:"uploaded_by"`
	UploadedAt  time.Time `json:"uploaded_at" gorm:"autoCreateTime"`
}

type PremiumBordereauxData struct {
	ID                                 uint      `gorm:"primaryKey" json:"id"`
	BordereauxID                       string    `gorm:"index" json:"bordereaux_id"`
	SchemeName                         string    `json:"scheme_name"`
	MemberName                         string    `json:"member_name"`
	MemberIdNumber                     string    `json:"member_id_number"`
	EmployeeNumber                     string    `json:"employee_number"`
	Category                           string    `json:"category"`
	Period                             string    `json:"period"`
	Month                              string    `json:"month"`
	Year                               int       `json:"year"`
	Age                                int       `json:"age"`
	AnnualSalary                       float64   `json:"annual_salary"`
	GlaAnnualPremium                   float64   `json:"gla_annual_premium"`
	PtdAnnualPremium                   float64   `json:"ptd_annual_premium"`
	CiAnnualPremium                    float64   `json:"ci_annual_premium"`
	PhiAnnualPremium                   float64   `json:"phi_annual_premium"`
	TotalAnnualPremium                 float64   `json:"total_annual_premium"`
	TotalAnnualPremiumExcludingFuneral float64   `json:"total_annual_premium_excluding_funeral"`
	TotalAnnualFuneralPremium          float64   `json:"total_annual_funeral_premium"`
	CreatedAt                          time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type MemberBordereauxData struct {
	ID                     uint      `gorm:"primaryKey" json:"id"`
	BordereauxID           string    `gorm:"index" json:"bordereaux_id"`
	SchemeName             string    `json:"scheme_name"`
	MemberName             string    `json:"member_name"`
	EmployeeNumber         string    `json:"employee_number"`
	MemberIdNumber         string    `json:"member_id_number"`
	Category               string    `json:"category"`
	Period                 string    `json:"period"`
	Month                  string    `json:"month"`
	Year                   int       `json:"year"`
	Gender                 string    `json:"gender"`
	DateOfBirth            time.Time `json:"date_of_birth"`
	AnnualSalary           float64   `json:"annual_salary"`
	GlaMultiple            float64   `json:"gla_multiple"`
	GlaCoveredSumAssured   float64   `json:"gla_covered_sum_assured"`
	GlaRetainedSumAssured  float64   `json:"gla_retained_sum_assured"`
	GlaCededSumAssured     float64   `json:"gla_ceded_sum_assured"`
	PtdMultiple            float64   `json:"ptd_multiple"`
	PtdCoveredSumAssured   float64   `json:"ptd_covered_sum_assured"`
	PtdRetainedSumAssured  float64   `json:"ptd_retained_sum_assured"`
	PtdCededSumAssured     float64   `json:"ptd_ceded_sum_assured"`
	CiMultiple             float64   `json:"ci_multiple"`
	CiCoveredSumAssured    float64   `json:"ci_covered_sum_assured"`
	CiRetainedSumAssured   float64   `json:"ci_retained_sum_assured"`
	CiCededSumAssured      float64   `json:"ci_ceded_sum_assured"`
	SglaMultiple           float64   `json:"sgla_multiple"`
	SglaCoveredSumAssured  float64   `json:"sgla_covered_sum_assured"`
	SglaRetainedSumAssured float64   `json:"sgla_retained_sum_assured"`
	SglaCededSumAssured    float64   `json:"sgla_ceded_sum_assured"`
	PhiMultiple            float64   `json:"phi_multiple"`
	PhiCoveredIncome       float64   `json:"phi_covered_income"`
	PhiRetainedIncome      float64   `json:"phi_retained_income"`
	PhiCededIncome         float64   `json:"phi_ceded_income"`
	TtdMultiple            float64   `json:"ttd_multiple"`
	TtdCoveredIncome       float64   `json:"ttd_covered_income"`
	TtdRetainedIncome      float64   `json:"ttd_retained_income"`
	TtdCededIncome         float64   `json:"ttd_ceded_income"`
	MmFuneralSumAssured    float64   `json:"mm_funeral_sum_assured"`
	SpFuneralSumAssured    float64   `json:"sp_funeral_sum_assured"`
	ChFuneralSumAssured    float64   `json:"ch_funeral_sum_assured"`
	ParFuneralSumAssured   float64   `json:"par_funeral_sum_assured"`
	DepFuneralSumAssured   float64   `json:"dep_funeral_sum_assured"`
	CreatedAt              time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type ClaimBordereauxData struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	BordereauxID   string     `gorm:"index" json:"bordereaux_id"`
	SchemeName     string     `json:"scheme_name"`
	MemberName     string     `json:"member_name"`
	MemberIdNumber string     `json:"member_id_number"`
	ClaimNumber    string     `json:"claim_number"`
	Category       string     `json:"category"`
	Period         string     `json:"period"`
	Month          string     `json:"month"`
	Year           int        `json:"year"`
	EventDate      *time.Time `json:"event_date"`
	ClaimAmount    float64    `json:"claim_amount"`
	ClaimType      string     `json:"claim_type"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

type GroupBenefitMapper struct {
	ID               int    `json:"id" gorm:"primary_key"`
	BenefitName      string `json:"benefit_name" gorm:"type:varchar(255);uniqueIndex"`
	BenefitCode      string `json:"benefit_code"`
	BenefitAlias     string `json:"benefit_alias"`
	BenefitAliasCode string `json:"benefit_alias_code"`
	IsMapped         bool   `json:"is_mapped"`
}

type HistoricalCredibilityData struct {
	ID                       int       `json:"id" gorm:"primary_key"`
	QuoteID                  int       `json:"quote_id"`
	Basis                    string    `json:"basis"`
	CreationDate             time.Time `json:"creation_date"`
	QuoteType                string    `json:"quote_type"`
	SchemeName               string    `json:"scheme_name"`
	SchemeID                 int       `json:"scheme_id"`
	Year                     int       `json:"year"`
	DurationInForce          float64   `json:"duration_in_force"`
	MemberCount              int       `json:"member_count"`
	ClaimCount               int       `json:"claim_count"`
	ExperiencePeriod         float64   `json:"experience_period"`
	CalculatedCredibility    float64   `json:"calculated_credibility"`
	ManuallyAddedCredibility float64   `json:"manually_added_credibility"`
	WeightedLifeYears        float64   `json:"weighted_life_years"`
	FullCredibilityThreshold float64   `json:"full_credibility_threshold"`
	AnnualGlaExperienceRate  float64   `json:"annual_gla_experience_rate"`
	AnnualPtdExperienceRate  float64   `json:"annual_ptd_experience_rate"`
	AnnualCiExperienceRate   float64   `json:"annual_ci_experience_rate"`
}

type MemberIndicativeDataSet struct {
	ID                           int     `json:"-" gorm:"primary_key"`
	QuoteID                      int     `json:"quote_id"`
	SchemeCategory               string  `json:"scheme_category"`
	MemberAverageAge             int     `json:"member_average_age"`
	MemberAverageIncome          float64 `json:"member_average_income"`
	MemberDataCount              int     `json:"member_data_count"`
	MemberMaleFemaleDistribution float64 `json:"member_male_female_distribution"`
}

type SchemeStatusUpdate struct {
	Status              Status `json:"status" binding:"required"`
	SchemeStatusMessage string `json:"scheme_status_message"`
}

// GroupSchemeStatusAudit keeps a history of status changes for a GroupScheme
// Table name will be inferred by GORM as group_scheme_status_audits
// If you use different naming conventions, ensure migrations match this structure
type GroupSchemeStatusAudit struct {
	ID            int       `json:"id" gorm:"primary_key"`
	SchemeID      int       `json:"scheme_id" gorm:"index"`
	OldStatus     Status    `json:"old_status"`
	NewStatus     Status    `json:"new_status"`
	StatusMessage string    `json:"status_message"`
	ChangedBy     string    `json:"changed_by"`
	ChangedAt     time.Time `json:"changed_at" gorm:"autoCreateTime"`
}

// MemberActivity stores structured history for a scheme member.
// Matches the format requested for scheme member history.
type MemberActivity struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	MemberID       int       `json:"member_id" gorm:"index"`
	MemberIDNumber string    `json:"member_id_number" gorm:"index"`
	Timestamp      time.Time `json:"timestamp" gorm:"column:timestamp;autoCreateTime"`
	Type           string    `json:"type" gorm:"column:type"` // e.g., enrollment, salary_change, claim, etc.
	Title          string    `json:"title" gorm:"column:title"`
	Description    string    `json:"description" gorm:"column:description"`
	Details        JSON      `json:"details" gorm:"type:json"`
	PerformedBy    string    `json:"performedBy" gorm:"column:performed_by"`
}

func (MemberActivity) TableName() string {
	return "member_activities"
}

// / Added from CoPilot
type ClaimRegistration struct {
	// Member Information
	MemberIDNumber string `json:"member_id_number" binding:"required"`
	MemberName     string `json:"member_name"`
	SchemeID       int    `json:"scheme_id" binding:"required"`
	SchemeName     string `json:"scheme_name"`

	// Claim Information
	BenefitType    BenefitType `json:"benefit_type" binding:"required"`
	BenefitName    string      `json:"benefit_name"`
	BenefitCode    string      `json:"benefit_code"`
	BenefitAlias   string      `json:"benefit_alias"`
	MemberType     string      `json:"member_type" binding:"required"`
	DateOfEvent    string      `json:"date_of_event" binding:"required"`
	DateNotified   string      `json:"date_notified" binding:"required"`
	ClaimAmount    float64     `json:"claim_amount" binding:"required,gt=0"`
	Priority       string      `json:"priority" binding:"required"`
	Status         string      `json:"status"`
	DateRegistered string      `json:"date_registered"`
	ClaimNumber    string      `json:"claim_number"`

	// Claimant Information (for dependants)
	ClaimantName          string `json:"claimant_name"`
	ClaimantIDNumber      string `json:"claimant_id_number"`
	RelationshipToMember  string `json:"relationship_to_member"`
	ClaimantContactNumber string `json:"claimant_contact_number"`

	// Documentation
	Description              string   `json:"description" binding:"required"`
	MissingRequiredDocuments []string `json:"missing_required_documents"`
}

// BenefitType represents the benefit type selection
type BenefitType struct {
	Value BenefitValue `json:"value"`
	Title string       `json:"title"`
}

type BenefitValue struct {
	BenefitCode  string `json:"benefit_code"`
	BenefitName  string `json:"benefit_name"`
	BenefitAlias string `json:"benefit_alias"`
	IsMapped     bool   `json:"is_mapped"`
}

// DocumentUpload represents uploaded documents with their types
type DocumentUpload struct {
	File         *multipart.FileHeader `json:"-"`
	DocumentType string                `json:"document_type"`
	DocumentName string                `json:"document_name"`
	FileName     string                `json:"file_name"`
	FileSize     int64                 `json:"file_size"`
}

// BenefitDocumentType represents a document type required for a specific benefit code
type BenefitDocumentType struct {
	ID          int    `json:"id" gorm:"primary_key"`
	BenefitCode string `json:"benefit_code" gorm:"index"` // e.g., GLA, SGLA, PTD, etc.
	Code        string `json:"code"`                      // e.g., claim_form, death_certificate
	Name        string `json:"name"`                      // e.g., Claim Form (official insurer form)
	Required    bool   `json:"required"`
}

// AuditLog represents a generic audit log entry.
// Used for tracking changes across different areas of the application.
type AuditLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Area       string    `json:"area" gorm:"index"`
	Entity     string    `json:"entity" gorm:"index"`
	EntityID   string    `json:"entity_id" gorm:"index"`
	Action     string    `json:"action"` // CREATE, UPDATE, DELETE
	Route      string    `json:"route"`
	RequestID  string    `json:"request_id"`
	PrevValues string    `json:"prev_values" gorm:"type:longtext"`
	NewValues  string    `json:"new_values" gorm:"type:longtext"`
	Diff       string    `json:"diff" gorm:"type:longtext"`
	ChangedBy  string    `json:"changed_by"`
	ChangedAt  time.Time `json:"changed_at" gorm:"index"`
}

// SumsAssuredResult holds the covered, retained, and ceded sums assured for all benefits
type SumsAssuredResult struct {
	// GLA - Group Life Assurance
	GlaCoveredSumAssured  float64
	GlaRetainedSumAssured float64
	GlaCededSumAssured    float64

	// PTD - Permanent Total Disability
	PtdCoveredSumAssured  float64
	PtdRetainedSumAssured float64
	PtdCededSumAssured    float64

	// CI - Critical Illness
	CiCoveredSumAssured  float64
	CiRetainedSumAssured float64
	CiCededSumAssured    float64

	// SGLA - Spouse Group Life Assurance
	SglaCoveredSumAssured  float64
	SglaRetainedSumAssured float64
	SglaCededSumAssured    float64

	// PHI - Permanent Health Insurance (monthly benefit)
	PhiMonthlyBenefit         float64
	PhiRetainedMonthlyBenefit float64
	PhiCededMonthlyBenefit    float64

	// TTD - Temporary Total Disability (monthly benefit)
	TtdMonthlyBenefit         float64
	TtdRetainedMonthlyBenefit float64
	TtdCededMonthlyBenefit    float64

	//Funeral Sums Assured
	MmFuneralSumAssured  float64
	SpFuneralSumAssured  float64
	ChFuneralSumAssured  float64
	ParFuneralSumAssured float64
	DepFuneralSumAssured float64

	MmRetainedFuneralSumAssured  float64
	SpRetainedFuneralSumAssured  float64
	ChRetainedFuneralSumAssured  float64
	ParRetainedFuneralSumAssured float64
	DepRetainedFuneralSumAssured float64

	MmCededFuneralSumAssured  float64
	SpCededFuneralSumAssured  float64
	ChCededFuneralSumAssured  float64
	ParCededFuneralSumAssured float64
	DepCededFuneralSumAssured float64
}

type PremiumComputation struct {
	// GLA - Group Life Assurance
	GlaRiskPremium   float64
	GlaOfficePremium float64

	// PTD - Permanent Total Disability
	PtdRiskPremium   float64
	PtdOfficePremium float64

	// CI - Critical Illness
	CiRiskPremium   float64
	CiOfficePremium float64

	// SGLA - Spouse Group Life Assurance
	SglaRiskPremium   float64
	SglaOfficePremium float64

	// PHI - Permanent Health Insurance (monthly benefit)
	PhiRiskPremium   float64
	PhiOfficePremium float64

	// TTD - Temporary Total Disability (monthly benefit)
	TtdRiskPremium   float64
	TtdOfficePremium float64

	//Funeral
	FuneralRiskPremium   float64
	FuneralOfficePremium float64

	TotalRiskPremiumExclFun float64
	TotalRiskFuneralPremium float64
	TotalRiskPremium        float64

	TotalOfficePremiumExclFun float64
	TotalOfficePremium        float64
	TotalOfficeFuneralPremium float64
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

type BenefitTreatyMap struct {
	ID       int
	Name     string
	TreatyID int
}

type MonthlyQuoteTrend struct {
	Month              int     `json:"month"`
	MonthName          string  `json:"month_name"`
	NewBusinessCount   int64   `json:"new_business_count"`
	RenewalCount       int64   `json:"renewal_count"`
	TotalCount         int64   `json:"total_count"`
	NewBusinessPremium float64 `json:"new_business_premium"`
	RenewalPremium     float64 `json:"renewal_premium"`
}

type BrokerMetric struct {
	BrokerID       int     `json:"broker_id"`
	BrokerName     string  `json:"broker_name"`
	TotalQuotes    int64   `json:"total_quotes"`
	AcceptedQuotes int64   `json:"accepted_quotes"`
	ConversionRate float64 `json:"conversion_rate"`
	TotalPremium   float64 `json:"total_premium"`
}

type QuoteFunnelStage struct {
	Stage   string  `json:"stage"`
	Count   int64   `json:"count"`
	Premium float64 `json:"premium"`
}

type DashboardPricingMetrics struct {
	AvgGlaRatePer1000  float64 `json:"avg_gla_rate_per_1000"`
	AvgPtdRatePer1000  float64 `json:"avg_ptd_rate_per_1000"`
	AvgCiRatePer1000   float64 `json:"avg_ci_rate_per_1000"`
	AvgSglaRatePer1000 float64 `json:"avg_sgla_rate_per_1000"`
	AvgDiscount        float64 `json:"avg_discount"`
	AvgCommissionPct   float64 `json:"avg_commission_pct"`
	ExpectedLossRatio  float64 `json:"expected_loss_ratio"`
	ActualLossRatio    float64 `json:"actual_loss_ratio"`
}

// WinProbabilityModel stores persisted logistic regression weights and normalisation parameters.
type WinProbabilityModel struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Weights        string    `json:"weights" gorm:"type:text"`          // JSON array float64
	FeatureNames   string    `json:"feature_names" gorm:"type:text"`    // JSON array string
	FeatureMeans   string    `json:"feature_means" gorm:"type:text"`    // JSON array float64
	FeatureStdDevs string    `json:"feature_std_devs" gorm:"type:text"` // JSON array float64
	TrainingSize   int       `json:"training_size"`
	Accuracy       float64   `json:"accuracy"`
	AUC            float64   `json:"auc"`
	TrainedAt      time.Time `json:"trained_at"`
	CreatedBy      string    `json:"created_by"`
}

// QuoteWinProbability stores the per-quote win probability score.
type QuoteWinProbability struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	QuoteID     int       `json:"quote_id" gorm:"uniqueIndex"`
	Score       float64   `json:"score"`                         // 0.0–1.0
	ScorePct    float64   `json:"score_pct"`                     // 0–100
	Band        string    `json:"band"`                          // low/medium/high/very_high
	TopFeatures string    `json:"top_features" gorm:"type:text"` // JSON [{name,contribution}]
	ModelID     int       `json:"model_id"`
	ScoredAt    time.Time `json:"scored_at"`
}

// GPTableStat tracks the row count for each Group Pricing rate table.
// It is updated on every upload and delete, allowing GetGPTableMetaData to
// resolve "populated?" with a single query instead of 32 sequential COUNTs.
type GPTableStat struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TableName string    `gorm:"uniqueIndex;not null;size:100" json:"table_name"`
	RowCount  int64     `gorm:"not null;default:0" json:"row_count"`
	UpdatedAt time.Time `json:"updated_at"`
}
