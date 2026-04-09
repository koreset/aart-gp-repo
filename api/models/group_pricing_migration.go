package models

// MigrationValidationError represents a single validation issue found during migration.
type MigrationValidationError struct {
	Template string `json:"template"`
	Row      int    `json:"row"`
	Column   string `json:"column"`
	Message  string `json:"message"`
}

// MigrationResult is returned by both the validate and execute migration endpoints.
type MigrationResult struct {
	Valid            bool                       `json:"valid"`
	Errors           []MigrationValidationError `json:"errors"`
	SchemeCount      int                        `json:"scheme_count"`
	CategoryCount    int                        `json:"category_count"`
	MemberCount      int                        `json:"member_count"`
	BeneficiaryCount int                        `json:"beneficiary_count"`
	ClaimsCount      int                        `json:"claims_count"`
	CreatedSchemeIDs []int                      `json:"created_scheme_ids,omitempty"`
	CreatedQuoteIDs  []int                      `json:"created_quote_ids,omitempty"`
}

// SchemeSetupRow is the CSV shadow struct for Template 1 (scheme_setup.csv).
type SchemeSetupRow struct {
	SchemeName             string  `csv:"scheme_name"`
	DistributionChannel    string  `csv:"distribution_channel"`
	BrokerName             string  `csv:"broker_name"`
	BrokerEmail            string  `csv:"broker_email"`
	ContactPerson          string  `csv:"contact_person"`
	ContactEmail           string  `csv:"contact_email"`
	CommencementDate       CsvTime `csv:"commencement_date"`
	CoverStartDate         CsvTime `csv:"cover_start_date"`
	CoverEndDate           CsvTime `csv:"cover_end_date"`
	RenewalDate            CsvTime `csv:"renewal_date"`
	Industry               string  `csv:"industry"`
	OccupationClass        int     `csv:"occupation_class"`
	Currency               string  `csv:"currency"`
	NormalRetirementAge    int     `csv:"normal_retirement_age"`
	FreeCoverLimit         float64 `csv:"free_cover_limit"`
	RiskRateCode           string  `csv:"risk_rate_code"`
	ExperienceRating       string  `csv:"experience_rating"`
	UseGlobalSalaryMultiple bool   `csv:"use_global_salary_multiple"`
}

// SchemeCategoryRow is the CSV shadow struct for Template 2 (scheme_categories.csv).
type SchemeCategoryRow struct {
	SchemeName     string  `csv:"scheme_name"`
	CategoryName   string  `csv:"category_name"`
	Basis          string  `csv:"basis"`
	FreeCoverLimit float64 `csv:"free_cover_limit"`
	Region         string  `csv:"region"`
	// GLA
	GlaBenefit                bool    `csv:"gla_benefit"`
	GlaSalaryMultiple         float64 `csv:"gla_salary_multiple"`
	GlaBenefitType            string  `csv:"gla_benefit_type"`
	GlaTerminalIllnessBenefit string  `csv:"gla_terminal_illness_benefit"`
	GlaWaitingPeriod          int     `csv:"gla_waiting_period"`
	GlaEducatorBenefit        string  `csv:"gla_educator_benefit"`
	// SGLA
	SglaBenefit        bool    `csv:"sgla_benefit"`
	SglaSalaryMultiple float64 `csv:"sgla_salary_multiple"`
	SglaMaxBenefit     float64 `csv:"sgla_max_benefit"`
	// PTD
	PtdBenefit              bool    `csv:"ptd_benefit"`
	PtdSalaryMultiple       float64 `csv:"ptd_salary_multiple"`
	PtdBenefitType          string  `csv:"ptd_benefit_type"`
	PtdRiskType             string  `csv:"ptd_risk_type"`
	PtdDeferredPeriod       int     `csv:"ptd_deferred_period"`
	PtdDisabilityDefinition string  `csv:"ptd_disability_definition"`
	// CI
	CiBenefit            bool    `csv:"ci_benefit"`
	CiSalaryMultiple     float64 `csv:"ci_salary_multiple"`
	CiBenefitStructure   string  `csv:"ci_benefit_structure"`
	CiBenefitDefinition  string  `csv:"ci_benefit_definition"`
	CiMaxBenefit         float64 `csv:"ci_max_benefit"`
	// TTD
	TtdBenefit                      bool    `csv:"ttd_benefit"`
	TtdRiskType                     string  `csv:"ttd_risk_type"`
	TtdMaximumBenefit               float64 `csv:"ttd_maximum_benefit"`
	TtdIncomeReplacementPercentage  float64 `csv:"ttd_income_replacement_percentage"`
	TtdWaitingPeriod                int     `csv:"ttd_waiting_period"`
	TtdDeferredPeriod               int     `csv:"ttd_deferred_period"`
	// PHI
	PhiBenefit                      bool    `csv:"phi_benefit"`
	PhiRiskType                     string  `csv:"phi_risk_type"`
	PhiMaximumBenefit               float64 `csv:"phi_maximum_benefit"`
	PhiIncomeReplacementPercentage  float64 `csv:"phi_income_replacement_percentage"`
	PhiWaitingPeriod                int     `csv:"phi_waiting_period"`
	PhiDeferredPeriod               int     `csv:"phi_deferred_period"`
	PhiNormalRetirementAge          int     `csv:"phi_normal_retirement_age"`
	PhiPremiumWaiver                string  `csv:"phi_premium_waiver"`
	// Family Funeral
	FamilyFuneralBenefit           bool    `csv:"family_funeral_benefit"`
	FuneralMainMemberSumAssured    float64 `csv:"funeral_main_member_sum_assured"`
	FuneralSpouseSumAssured        float64 `csv:"funeral_spouse_sum_assured"`
	FuneralChildrenSumAssured      float64 `csv:"funeral_children_sum_assured"`
	FuneralMaxChildren             int     `csv:"funeral_max_children"`
}

// BeneficiaryRow is the CSV shadow struct for Template 4 (beneficiaries.csv).
type BeneficiaryRow struct {
	SchemeName           string  `csv:"scheme_name"`
	MemberIdNumber       string  `csv:"member_id_number"`
	BeneficiaryFullName  string  `csv:"beneficiary_full_name"`
	Relationship         string  `csv:"relationship"`
	IdType               string  `csv:"id_type"`
	IdNumber             string  `csv:"id_number"`
	Gender               string  `csv:"gender"`
	DateOfBirth          CsvTime `csv:"date_of_birth"`
	AllocationPercentage float64 `csv:"allocation_percentage"`
	BenefitTypes         string  `csv:"benefit_types"`
	ContactNumber        string  `csv:"contact_number"`
	Email                string  `csv:"email"`
	Address              string  `csv:"address"`
	BankName             string  `csv:"bank_name"`
	BranchCode           string  `csv:"branch_code"`
	AccountNumber        string  `csv:"account_number"`
	AccountType          string  `csv:"account_type"`
	GuardianName         string  `csv:"guardian_name"`
	GuardianRelationship string  `csv:"guardian_relationship"`
	GuardianIdNumber     string  `csv:"guardian_id_number"`
	GuardianContact      string  `csv:"guardian_contact"`
}

// MemberDataRow is the CSV shadow struct for Template 3 (member_data.csv).
// Mirrors the memberDataCSV struct used in SaveQuoteTables().
type MemberDataRow struct {
	Year                         int     `csv:"year"`
	SchemeName                   string  `csv:"scheme_name"`
	MemberName                   string  `csv:"member_name"`
	MemberIdNumber               string  `csv:"member_id_number"`
	MemberIdType                 string  `csv:"member_id_type"`
	SchemeCategory               string  `csv:"scheme_category"`
	Gender                       string  `csv:"gender"`
	Email                        string  `csv:"email"`
	EmployeeNumber               string  `csv:"employee_number"`
	DateOfBirth                  CsvTime `csv:"date_of_birth"`
	AnnualSalary                 float64 `csv:"annual_salary"`
	ContributionWaiverProportion float64 `csv:"contribution_waiver_proportion"`
	EntryDate                    CsvTime `csv:"entry_date"`
	ExitDate                     CsvTime `csv:"exit_date"`
	EffectiveExitDate            CsvTime `csv:"effective_exit_date"`
	BenefitsGlaMultiple          float64 `csv:"benefits_gla_multiple"`
	BenefitsSglaMultiple         float64 `csv:"benefits_sgla_multiple"`
	BenefitsPtdMultiple          float64 `csv:"benefits_ptd_multiple"`
	BenefitsCiMultiple           float64 `csv:"benefits_ci_multiple"`
	BenefitsTtdMultiple          float64 `csv:"benefits_ttd_multiple"`
	BenefitsPhiMultiple          float64 `csv:"benefits_phi_multiple"`
}
