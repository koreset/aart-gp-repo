package models

import (
	"time"
)

type ExpRunGroup struct {
	ID               int                     `json:"id" gorm:"primary_key"`
	Name             string                  `json:"name" gorm:"unique"`
	ExpRuns          []ExpAnalysisRunSetting `json:"exp_runs" gorm:"-"`
	CreationDate     time.Time               `json:"creation_date,omitempty"`
	ProcessingStatus string                  `json:"processing_status" csv:"processing_status"`
	ProcessedRecords float64                 `json:"processed_records" csv:"processed_records"`
	TotalRecords     int                     `json:"total_records" csv:"total_records"`
	FailureReason    string                  `json:"failure_reason" csv:"failure_reason"`
	RunDuration      float64                 `json:"run_duration" csv:"run_duration"`
	UserName         string                  `json:"user_name" csv:"user_name"`
	UserEmail        string                  `json:"email" csv:"email"`
	Status           string                  `json:"status" csv:"status"`
}

type ExpAnalysisRunSetting struct {
	ID                   int       `json:"id" gorm:"primary_key"`
	ExpRunGroupID        int       `json:"run_group_id"`
	ExpRunGroupName      string    `json:"run_group_name" `
	RunName              string    `json:"run_name"`
	RunDate              string    `json:"run_date"`
	RunType              string    `json:"run_type"`
	Description          string    `json:"description"`
	PeriodStartDate      string    `json:"period_start_date,omitempty"`
	PeriodEndDate        string    `json:"period_end_date,omitempty"`
	ExpConfigurationName string    `json:"exp_configuration_name"`
	ExpConfigurationId   int       `json:"exp_configuration_id"`
	ExposureDataYear     int       `json:"exposure_data_year"`
	ExposureDataVersion  string    `json:"exposure_data_version"`
	ActualDataYear       int       `json:"actual_data_year" csv:"actual_data_year"`
	ActualDataVersion    string    `json:"actual_data_version" csv:"actual_data_version"`
	LapseDataYear        int       `json:"lapse_data_year" csv:"lapse_data_year"`
	LapseDataVersion     string    `json:"lapse_data_version" csv:"lapse_data_version"`
	MortalityDataYear    int       `json:"mortality_data_year" csv:"mortality_data_year"`
	MortalityDataVersion string    `json:"mortality_data_version" csv:"mortality_data_version"`
	AgeBandVersion       string    `json:"age_band_version" csv:"age_band_version"`
	FailureReason        string    `json:"failure_reason" csv:"failure_reason"`
	TotalRecords         int       `json:"total_records" csv:"total_records"`
	ProcessedRecords     float64   `json:"processed_records" csv:"processed_records"`
	ProcessingStatus     string    `json:"processing_status" csv:"processing_status"`
	RunDuration          float64   `json:"run_duration" csv:"run_duration"`
	CreationDate         time.Time `json:"creation_date,omitempty"`
	UserName             string    `json:"user_name" csv:"user_name"`
	UserEmail            string    `json:"user_email" csv:"user_email"`
	Status               string    `json:"status" csv:"status"`
}

type ExpExposureData struct {
	ID                         int     `json:"id" gorm:"primary_key"`
	ExpConfigurationID         int     `json:"-" csv:"exp_configuration_id"`
	PolicyNumber               string  `json:"policy_number" csv:"policy_number"`
	MemberType                 string  `json:"member_type" csv:"member_type"`
	Gender                     string  `json:"gender" csv:"gender"`
	SumAssured                 float64 `json:"sum_assured" csv:"sum_assured"`
	DateOfBirth                string  `json:"date_of_birth" csv:"date_of_birth"`
	CommencementDate           string  `json:"commencement_date" csv:"commencement_date"`
	ClaimsDate                 string  `json:"claims_date" csv:"claims_date"`
	LapseDate                  string  `json:"lapse_date" csv:"lapse_date"`
	ExitDate                   string  `json:"exit_date" csv:"exit_date"`
	ExitCode                   string  `json:"exit_code" csv:"exit_code"`
	PremiumCeaseDate           string  `json:"premium_cease_date" csv:"premium_cease_date"`
	DisabilityBenefitCeaseDate string  `json:"disability_benefit_cease_date" csv:"disability_benefit_cease_date"`
	LifeCoverCeaseDate         string  `json:"life_cover_cease_date" csv:"life_cover_cease_date"`
	ProductName                string  `json:"product_name" csv:"product_name"`
	CreatedBy                  string  `json:"created_by" csv:"created_by"`
	Version                    string  `json:"version" csv:"version"`
	Year                       int     `json:"year" csv:"year"`
}

type ExpActualData struct {
	ID                   int     `json:"id" gorm:"primary_key"`
	ExpConfigurationID   int     `json:"-" csv:"exp_configuration_id"`
	ProductName          string  `json:"product_name" csv:"product_name"`
	PremiumIncome        float64 `json:"premium_income" csv:"premium_income"`
	ClaimsCount          int     `json:"claims_count" csv:"claims_count"`
	ExitSumAssuredClaims float64 `json:"exit_sum_assured_claims" csv:"exit_sum_assured_claims"`
	ExitPremiumClaims    float64 `json:"exit_premium_claims" csv:"exit_premium_claims"`
	ClaimsPaid           float64 `json:"claims_paid" csv:"claims_paid"`
	OutstandingClaims    float64 `json:"outstanding_claims" csv:"outstanding_claims"`
	LapseCount           int     `json:"lapse_count" csv:"lapse_count"`
	ExitSumAssuredLapses float64 `json:"exit_sum_assured_lapses" csv:"exit_sum_assured_lapses"`
	CreatedBy            string  `json:"created_by" csv:"created_by"`
	Version              string  `json:"version" csv:"version"`
	Year                 int     `json:"year" csv:"year"`
}

type ExpConfiguration struct {
	ID                           int                              `json:"id" gorm:"primary_key"`
	Name                         string                           `json:"name" gorm:"unique"`
	ExpDataYearVersions          []ExpExpDataYearVersion          `json:"exp_data_year_versions" gorm:"-"`
	ExpDataYears                 []int                            `json:"exp_data_years" gorm:"-"`
	ActualDataYearVersions       []ExpActualDataYearVersion       `json:"actual_data_year_versions" gorm:"-"`
	ActualDataYears              []int                            `json:"actual_data_years" gorm:"-"`
	CurrentLapseYearVersions     []ExpCurrentLapseYearVersion     `json:"current_lapse_year_versions" gorm:"-"`
	CurrentMortalityYearVersions []ExpCurrentMortalityYearVersion `json:"current_mortality_year_versions" gorm:"-"`
}

type ExpExpDataYearVersion struct {
	ID                int    `json:"id" csv:"id" gorm:"primary_key"`
	ConfigurationName string `json:"configuration_name" csv:"configuration_name"`
	ConfigurationId   int    `json:"configuration_id"`
	Year              int    `json:"year" csv:"year"`
	Version           string `json:"version" csv:"version"`
	Count             int    `json:"count" csv:"count"`
}

type ExpActualDataYearVersion struct {
	ID                int    `json:"id" csv:"id" gorm:"primary_key"`
	ConfigurationName string `json:"configuration_name" csv:"configuration_name"`
	ConfigurationId   int    `json:"configuration_id"`
	Year              int    `json:"year" csv:"year"`
	Version           string `json:"version" csv:"version"`
	Count             int    `json:"count" csv:"count"`
}

type ExpCurrentMortalityYearVersion struct {
	ID                int    `json:"id" csv:"id" gorm:"primary_key"`
	ConfigurationName string `json:"configuration_name" csv:"configuration_name"`
	ConfigurationId   int    `json:"configuration_id"`
	Year              int    `json:"year" csv:"year"`
	Version           string `json:"version" csv:"version"`
	Count             int    `json:"count" csv:"count"`
}

type ExpCurrentLapseYearVersion struct {
	ID                int    `json:"id" csv:"id" gorm:"primary_key"`
	ConfigurationName string `json:"configuration_name" csv:"configuration_name"`
	ConfigurationId   int    `json:"configuration_id"`
	Year              int    `json:"year" csv:"year"`
	Version           string `json:"version" csv:"version"`
	Count             int    `json:"count" csv:"count"`
}

type ExpExposureActualYears struct {
	ExposureYears  []int `json:"exposure_years"`
	ActualYears    []int `json:"actual_years"`
	LapseYears     []int `json:"lapse_years"`
	MortalityYears []int `json:"mortality_years"`
}

type TableDataPayload struct {
	Year            int    `json:"year"`
	Version         string `json:"version"`
	ConfigurationId int    `json:"configuration_id"`
	TableName       string `json:"table_name"`
}

type ExposureModelPoint struct { // calculated values from ExpExposureData
	ID                               int     `json:"id" gorm:"primary_key"`
	ExpRunID                         int     `json:"exp_run_id" gorm:"index"`
	ExpRunGroupID                    int     `json:"exp_run_group_id" gorm:"index"`
	ProductName                      string  `json:"product_name" csv:"product_name"`
	ExpConfigurationID               int     `json:"exp_configuration_id" csv:"exp_configuration_id"`
	PolicyNumber                     string  `json:"policy_number" csv:"policy_number"`
	MemberType                       string  `json:"member_type" csv:"member_type"`
	Gender                           string  `json:"gender" csv:"gender"`
	SumAssured                       float64 `json:"sum_assured" csv:"sum_assured"`
	DateOfBirth                      string  `json:"date_of_birth" csv:"date_of_birth"`
	CommencementDate                 string  `json:"commencement_date" csv:"commencement_date"`
	DisabilityBenefitCeaseDate       string  `json:"disability_benefit_cease_date" csv:"disability_benefit_cease_date"`
	ClaimsDate                       string  `json:"claims_date" csv:"claims_date"`
	LapseDate                        string  `json:"lapse_date" csv:"lapse_date"`
	ExitCode                         string  `json:"exit_code" csv:"exit_code"`
	ExitDate                         string  `json:"exit_date" csv:"exit_date"`
	PremiumCeaseDate                 string  `json:"premium_cease_date" csv:"premium_cease_date"`
	LifeCoverCeaseDate               string  `json:"life_cover_cease_date" csv:"life_cover_cease_date"`
	AgeAtStart                       int     `json:"age_at_start" csv:"age_at_start"`
	AgeAtExit                        int     `json:"age_at_exit" csv:"age_at_exit"`
	DateAtStart                      string  `json:"date_at_start" csv:"date_at_start"`
	DateReachingAgeNext              string  `json:"date_reaching_age_next" csv:"date_reaching_age_next"`
	DateAtEndAgeX                    string  `json:"date_at_end_age_x" csv:"date_at_end_age_x"`
	MaxDuration                      float64 `json:"max_duration" csv:"max_duration"`
	DateAtEnd                        string  `json:"date_at_end" csv:"date_at_end"`
	AgeXExposure                     float64 `json:"age_x_exposure" csv:"age_x_exposure"`
	AgeNextExposure                  float64 `json:"age_next_exposure" csv:"age_next_exposure"`
	AmountAgeXExposure               float64 `json:"amount_age_x_exposure" csv:"amount_age_x_exposure"`
	AmountAgeNextExposure            float64 `json:"amount_age_next_exposure" csv:"amount_age_next_exposure"`
	ExpectedClaimCountMale           float64 `json:"expected_claim_count_male" csv:"expected_claim_count_male"`
	ExpectedClaimCountFemale         float64 `json:"expected_claim_count_female" csv:"expected_claim_count_female"`
	ExpectedClaimAmountMale          float64 `json:"expected_claim_amount_male" csv:"expected_claim_amount_male"`
	ExpectedClaimAmountFemale        float64 `json:"expected_claim_amount_female" csv:"expected_claim_amount_female"`
	ExpectedAgeNextClaimCountMale    float64 `json:"expected_age_next_claim_count_male" csv:"expected_age_next_claim_count_male"`
	ExpectedAgeNextClaimCountFemale  float64 `json:"expected_age_next_claim_count_female" csv:"expected_age_next_claim_count_female"`
	ExpectedAgeNextClaimAmountMale   float64 `json:"expected_age_next_claim_amount_male" csv:"expected_age_next_claim_amount_male"`
	ExpectedAgeNextClaimAmountFemale float64 `json:"expected_age_next_claim_amount_female" csv:"expected_age_next_claim_amount_female"`
	RelevantDeathClaimIndicator      int     `json:"relevant_death_claim_indicator" csv:"relevant_death_claim_indicator"`
	RelevantDisabilityClaimIndicator int     `json:"relevant_disability_claim_indicator" csv:"relevant_disability_claim_indicator"`
	RelevantLapseIndicator           int     `json:"relevant_lapse_indicator" csv:"relevant_lapse_indicator"`
	AgeXBand                         int     `json:"age_x_band" csv:"age_x_band"`
	AgeNextBand                      int     `json:"age_next_band" csv:"age_next_band"`
	DurationIfAtStart                float64 `json:"duration_if_at_start" csv:"duration_if_at_start"`
	DurationIfATEnd                  float64 `json:"duration_if_at_end" csv:"duration_if_at_end"`
	TotalDurationExposure            float64 `json:"total_duration_exposure" csv:"total_duration_exposure"`
	DurationInYear1                  float64 `json:"duration_in_year1" csv:"duration_in_year1"`
	DurationInYear2                  float64 `json:"duration_in_year2" csv:"duration_in_year2"`
	DurationInYear3                  float64 `json:"duration_in_year3" csv:"duration_in_year3"`
	DurationInYear4                  float64 `json:"duration_in_year4" csv:"duration_in_year4"`
	DurationInYear5Plus              float64 `json:"duration_in_year5_plus" csv:"duration_in_year5_plus"`
	DurationIfAtEndY                 float64 `json:"duration_if_at_end_y" csv:"duration_if_at_end_y"`
	CreatedBy                        string  `json:"created_by" csv:"created_by"`
	Version                          string  `json:"version" csv:"version"`
	Year                             int     `json:"year" csv:"year"`
}

type ExpCrudeResult struct { // calculated values from ExpExposureData
	ID                        int     `json:"id" gorm:"primary_key"`
	ExpRunID                  int     `json:"exp_run_id" gorm:"exp_run_id"`
	ExpRunGroupID             int     `json:"exp_run_group_id" gorm:"exp_run_group_id"`
	Age                       int     `json:"age" csv:"age"`
	AnnualQxMale              float64 `json:"annual_qx_male" csv:"annual_qx_male"`
	AnnualQxFemale            float64 `json:"annual_qx_female" csv:"annual_qx_female"`
	MonthlyQxMale             float64 `json:"monthly_qx_male" csv:"monthly_qx_male"`
	MonthlyQxFemale           float64 `json:"monthly_qx_female" csv:"monthly_qx_female"`
	ActualClaimCountMale      float64 `json:"actual_claim_count_male" csv:"actual_claim_count_male"`
	ActualClaimCountFemale    float64 `json:"actual_claim_count_female" csv:"actual_claim_count_female"`
	ActualClaimAmountMale     float64 `json:"actual_claim_amount_male" csv:"actual_claim_amount_male"`
	ActualClaimAmountFemale   float64 `json:"actual_claim_amount_female" csv:"actual_claim_amount_female"`
	ExposureCountMale         float64 `json:"exposure_count_male" csv:"exposure_count_male"`
	ExposureCountFemale       float64 `json:"exposure_count_female" csv:"exposure_count_female"`
	ExposureAmountMale        float64 `json:"exposure_amount_male" csv:"exposure_amount_male"`
	ExposureAmountFemale      float64 `json:"exposure_amount_female" csv:"exposure_amount_female"`
	ExpectedClaimCountMale    float64 `json:"expected_claim_count_male" csv:"expected_claim_count_male"`
	ExpectedClaimCountFemale  float64 `json:"expected_claim_count_female" csv:"expected_claim_count_female"`
	ExpectedClaimAmountMale   float64 `json:"expected_claim_amount_male" csv:"expected_claim_amount_male"`
	ExpectedClaimAmountFemale float64 `json:"expected_claim_amount_female" csv:"expected_claim_amount_female"`
	CrudeRatesLivesMale       float64 `json:"crude_rates_lives_male" csv:"crude_rates_lives_male"`
	CrudeRatesLivesFemale     float64 `json:"crude_rates_lives_female" csv:"crude_rates_lives_female"`
	CrudeRatesAmountMale      float64 `json:"crude_rates_amount_male" csv:"crude_rates_amount_male"`
	CrudeRatesAmountFemale    float64 `json:"crude_rates_amount_female" csv:"crude_rates_amount_female"`
	DataPointCount            float64 `json:"data_point_count" csv:"data_point_count"`
	CreatedBy                 string  `json:"created_by" csv:"created_by"`
	Version                   string  `json:"version" csv:"version"`
	Year                      int     `json:"year" csv:"year"`
}

type TotalMortalityExpAnalysisResult struct { // calculated values from ExpExposureData
	ID                                  int     `json:"id" gorm:"primary_key"`
	ExpRunID                            int     `json:"exp_run_id" gorm:"exp_run_id"`
	ExpRunGroupID                       int     `json:"exp_run_group_id" gorm:"exp_run_group_id"`
	Period                              string  `json:"period" gorm:"period"`
	TotalExposureCountMale              float64 `json:"total_exposure_count_male" csv:"total_exposure_count_male"`
	TotalExposureCountFemale            float64 `json:"total_exposure_count_female" csv:"total_exposure_count_female"`
	TotalExposureAmountMale             float64 `json:"total_exposure_amount_male" csv:"total_exposure_amount_male"`
	TotalExposureAmountFemale           float64 `json:"total_exposure_amount_female" csv:"total_exposure_amount_female"`
	TotalExpectedClaimCountMale         float64 `json:"total_expected_claim_count_male" csv:"total_expected_claim_count_male"`
	TotalExpectedClaimCountFemale       float64 `json:"total_expected_claim_count_female" csv:"total_expected_claim_count_female"`
	TotalExpectedClaimAmountMale        float64 `json:"total_expected_claim_amount_male" csv:"total_expected_claim_amount_male"`
	TotalExpectedClaimAmountFemale      float64 `json:"total_expected_claim_amount_female" csv:"total_expected_claim_amount_female"`
	TotalActualClaimCountMale           float64 `json:"total_actual_claim_count_male" csv:"total_actual_claim_count_male"`
	TotalActualClaimCountFemale         float64 `json:"total_actual_claim_count_female" csv:"total_actual_claim_count_female"`
	TotalActualClaimAmountMale          float64 `json:"total_actual_claim_amount_male" csv:"total_actual_claim_amount_male"`
	TotalActualClaimAmountFemale        float64 `json:"total_actual_claim_amount_female" csv:"total_actual_claim_amount_female"`
	TotalExpectedCrudeRatesLivesMale    float64 `json:"total_expected_crude_rates_lives_male" csv:"total_expected_crude_rates_lives_male"`
	TotalExpectedCrudeRatesLivesFemale  float64 `json:"total_expected_crude_rates_lives_female" csv:"total_expected_crude_rates_lives_female"`
	TotalExpectedCrudeRatesAmountMale   float64 `json:"total_expected_crude_rates_amount_male" csv:"total_expected_crude_rates_amount_male"`
	TotalExpectedCrudeRatesAmountFemale float64 `json:"total_expected_crude_rates_amount_female" csv:"total_expected_crude_rates_amount_female"`
	TotalActualCrudeRatesLivesMale      float64 `json:"total_actual_crude_rates_lives_male" csv:"total_actual_crude_rates_lives_male"`
	TotalActualCrudeRatesLivesFemale    float64 `json:"actual_crude_rates_lives_female" csv:"actual_crude_rates_lives_female"`
	TotalActualCrudeRatesAmountMale     float64 `json:"actual_crude_rates_amount_male" csv:"actual_crude_rates_amount_male"`
	TotalActualCrudeRatesAmountFemale   float64 `json:"actual_crude_rates_amount_female" csv:"actual_crude_rates_amount_female"`
	CreatedBy                           string  `json:"created_by" csv:"created_by"`
	Version                             string  `json:"version" csv:"version"`
	Year                                int     `json:"year" csv:"year"`
}

type TotalLapseExpAnalysisResult struct { // calculated values from ExpExposureData
	ID                           int     `json:"id" gorm:"primary_key"`
	ExpRunID                     int     `json:"exp_run_id" gorm:"exp_run_id"`
	ExpRunGroupID                int     `json:"exp_run_group_id" gorm:"exp_run_group_id"`
	Period                       string  `json:"period" gorm:"period"`
	TotalYear1Exposure           float64 `json:"total_year_1_exposure" gorm:"total_year_1_exposure"`
	TotalActualYear1Lapses       float64 `json:"total_actual_year_1_lapses" gorm:"total_actual_year_1_lapses"`
	ActualYear1Ux                float64 `json:"actual_year_1_ux" gorm:"actual_year_1_ux"`
	ActualYear1MonthlyRate       float64 `json:"actual_year_1_monthly_rate" gorm:"actual_year_1_monthly_rate"`
	ActualYear1AnnualRate        float64 `json:"actual_year_1_annual_rate" gorm:"actual_year_1_annual_rate"`
	TotalExpectedYear1Lapses     float64 `json:"total_expected_year_1_lapses" gorm:"total_expected_year_1_lapses"`
	ExpectedYear1Ux              float64 `json:"expected_year_1_ux" gorm:"expected_year_1_ux"`
	ExpectedYear1MonthlyRate     float64 `json:"expected_year_1_monthly_rate" gorm:"expected_year_1_monthly_rate"`
	ExpectedYear1AnnualRate      float64 `json:"expected_year_1_annual_rate" gorm:"expected_year_1_annual_rate"`
	TotalYear2Exposure           float64 `json:"total_year_2_exposure" gorm:"total_year_2_exposure"`
	TotalActualYear2Lapses       float64 `json:"total_actual_year_2_lapses" gorm:"total_actual_year_2_lapses"`
	ActualYear2Ux                float64 `json:"actual_year_2_ux" gorm:"actual_year_2_ux"`
	ActualYear2MonthlyRate       float64 `json:"actual_year_2_monthly_rate" gorm:"actual_year_2_monthly_rate"`
	ActualYear2AnnualRate        float64 `json:"actual_year_2_annual_rate" gorm:"actual_year_2_annual_rate"`
	TotalExpectedYear2Lapses     float64 `json:"total_expected_year_2_lapses" gorm:"total_expected_year_2_lapses"`
	ExpectedYear2Ux              float64 `json:"expected_year_2_ux" gorm:"expected_year_2_ux"`
	ExpectedYear2MonthlyRate     float64 `json:"expected_year_2_monthly_rate" gorm:"expected_year_2_monthly_rate"`
	ExpectedYear2AnnualRate      float64 `json:"expected_year_2_annual_rate" gorm:"expected_year_2_annual_rate"`
	TotalYear3Exposure           float64 `json:"total_year_3_exposure" gorm:"total_year_3_exposure"`
	TotalActualYear3Lapses       float64 `json:"total_actual_year_3_lapses" gorm:"total_actual_year_3_lapses"`
	ActualYear3Ux                float64 `json:"actual_year_3_ux" gorm:"actual_year_3_ux"`
	ActualYear3MonthlyRate       float64 `json:"actual_year_3_monthly_rate" gorm:"actual_year_3_monthly_rate"`
	ActualYear3AnnualRate        float64 `json:"actual_year_3_annual_rate" gorm:"actual_year_3_annual_rate"`
	TotalExpectedYear3Lapses     float64 `json:"total_expected_year_3_lapses" gorm:"total_expected_year_3_lapses"`
	ExpectedYear3Ux              float64 `json:"expected_year_3_ux" gorm:"expected_year_3_ux"`
	ExpectedYear3MonthlyRate     float64 `json:"expected_year_3_monthly_rate" gorm:"expected_year_3_monthly_rate"`
	ExpectedYear3AnnualRate      float64 `json:"expected_year_3_annual_rate" gorm:"expected_year_3_annual_rate"`
	TotalYear4Exposure           float64 `json:"year_4_exposure" gorm:"year_4_exposure"`
	TotalActualYear4Lapses       float64 `json:"total_actual_year_4_lapses" gorm:"total_actual_year_4_lapses"`
	ActualYear4Ux                float64 `json:"actual_year_4_ux" gorm:"actual_year_4_ux"`
	ActualYear4MonthlyRate       float64 `json:"actual_year_4_monthly_rate" gorm:"actual_year_4_monthly_rate"`
	ActualYear4AnnualRate        float64 `json:"actual_year_4_annual_rate" gorm:"actual_year_4_annual_rate"`
	TotalExpectedYear4Lapses     float64 `json:"total_expected_year_4_lapses" gorm:"total_expected_year_4_lapses"`
	ExpectedYear4Ux              float64 `json:"expected_year_4_ux" gorm:"expected_year_4_ux"`
	ExpectedYear4MonthlyRate     float64 `json:"expected_year_4_monthly_rate" gorm:"expected_year_4_monthly_rate"`
	ExpectedYear4AnnualRate      float64 `json:"expected_year_4_annual_rate" gorm:"expected_year_4_annual_rate"`
	TotalYear5PlusExposure       float64 `json:"total_year_5_plus_exposure" gorm:"total_year_5_plus_exposure"`
	TotalActualYear5PlusLapses   float64 `json:"total_actual_year_5_plus_lapses" gorm:"total_actual_year_5_plus_lapses"`
	ActualYear5PlusUx            float64 `json:"actual_year_5_plus_ux" gorm:"actual_year_5_plus_ux"`
	ActualYear5PlusMonthlyRate   float64 `json:"actual_year_5_plus_monthly_rate" gorm:"actual_year_5_plus_monthly_rate"`
	ActualYear5PlusAnnualRate    float64 `json:"actual_year_5_plus_annual_rate" gorm:"actual_year_5_plus_annual_rate"`
	TotalExpectedYear5PlusLapses float64 `json:"total_expected_year_5_plus_lapses" gorm:"total_expected_year_5_plus_lapses"`
	ExpectedYear5PlusUx          float64 `json:"expected_year_5_plus_ux" gorm:"expected_year_5_plus_ux"`
	ExpectedYear5PlusMonthlyRate float64 `json:"expected_year_5_plus_monthly_rate" gorm:"expected_year_5_plus_monthly_rate"`
	ExpectedYear5PlusAnnualRate  float64 `json:"expected_year_5_plus_annual_rate" gorm:"expected_year_5_plus_annual_rate"`
	CreatedBy                    string  `json:"created_by" csv:"created_by"`
	Version                      string  `json:"version" csv:"version"`
	Year                         int     `json:"year" csv:"year"`
}

type ExpAgeBand struct {
	ID      int    `json:"-" gorm:"primary_key"`
	Name    string `json:"name" csv:"name"`
	MinAge  int    `json:"min_age" csv:"min_age"`
	MaxAge  int    `json:"max_age" csv:"max_age"`
	Year    int    `json:"year" csv:"year"`
	Version string `json:"version" csv:"version"`
}

type ExpActualsVsExpected struct {
	ID                      int     `json:"-" gorm:"primary_key"`
	ExpRunId                int     `json:"-" gorm:"exp_run_id"`
	AgeNext                 string  `json:"age_next" csv:"age_next"`
	NumDeathsTotal          float64 `json:"num_deaths_total" csv:"num_deaths_total"`
	ExpectedNumDeathsTotal  float64 `json:"expected_num_deaths_total" csv:"expected_num_deaths_total"`
	LivesExposureTotal      float64 `json:"lives_exposure_total" csv:"lives_exposure_total"`
	CrudeRateTotal          float64 `json:"crude_rate_total" csv:"crude_rate_total"`
	ExpectedRateTotal       float64 `json:"expected_rate_total" csv:"expected_rate_total"`
	AvETotal                float64 `json:"ave_total" csv:"ave_total"`
	NumDeathsMale           float64 `json:"num_deaths_male" csv:"num_deaths_male"`
	ExpectedNumDeathsMale   float64 `json:"expected_num_deaths_male" csv:"expected_num_deaths_male"`
	LivesExposureMale       float64 `json:"lives_exposure_male" csv:"lives_exposure_male"`
	CrudeRateMale           float64 `json:"crude_rate_male" csv:"crude_rate_male"`
	ExpectedRateMale        float64 `json:"expected_rate_male" csv:"expected_rate_male"`
	AvEMale                 float64 `json:"ave_male" csv:"ave_male"`
	NumDeathsFemale         float64 `json:"num_deaths_female" csv:"num_deaths_female"`
	ExpectedNumDeathsFemale float64 `json:"expected_num_deaths_female" csv:"expected_num_deaths_female"`
	LivesExposureFemale     float64 `json:"lives_exposure_female" csv:"lives_exposure_female"`
	CrudeRateFemale         float64 `json:"crude_rate_female" csv:"crude_rate_female"`
	ExpectedRateFemale      float64 `json:"expected_rate_female" csv:"expected_rate_female"`
	AvEFemale               float64 `json:"ave_female" csv:"ave_female"`
}

type ExpExposureDataSummary struct {
	ID               int     `json:"-" gorm:"primary_key"`
	ExpRunId         int     `json:"-" gorm:"exp_run_id"`
	AgeNext          string  `json:"age_next" csv:"age_next"`
	Male             float64 `json:"num_deaths_total" csv:"num_deaths_total"`
	Female           float64 `json:"expected_num_deaths_total" csv:"expected_num_deaths_total"`
	TotalCount       float64 `json:"total_count" csv:"total_count"`
	MaleProportion   float64 `json:"male_proportion" csv:"male_proportion"`
	FemaleProportion float64 `json:"female_proportion" csv:"female_proportion"`
	TotalProportion  float64 `json:"total_proportion" csv:"total_proportion"`
	AvEFemale        float64 `json:"ave_female" csv:"ave_female"`
}

type ExpCurrentLapse struct {
	ID                  int     `json:"id" gorm:"primary_key"`
	DurationInForceYear int     `json:"duration_in_force_year" csv:"duration_in_force_year"`
	LapseRate           float64 `json:"lapse_rate" csv:"lapse_rate"`
	Version             string  `json:"version" csv:"version"`
	Year                int     `json:"year" csv:"year"`
	CreatedBy           string  `json:"created_by" csv:"created_by"`
}

type ExpCurrentMortality struct {
	ID        int     `json:"id" gorm:"primary_key"`
	Anb       int     `json:"anb" csv:"anb"`
	Gender    string  `json:"gender" csv:"gender"`
	Qx        float64 `json:"qx" csv:"qx"`
	Version   string  `json:"version" csv:"version"`
	Year      int     `json:"year" csv:"year"`
	CreatedBy string  `json:"created_by" csv:"created_by"`
}

type ExpLapseCrudeResult struct { // calculated values from ExpExposureData
	ID                  int     `json:"-" gorm:"primary_key"`
	ExpRunID            int     `json:"exp_run_id" gorm:"exp_run_id"`
	ExpRunGroupID       int     `json:"exp_run_group_id" gorm:"exp_run_group_id"`
	DurationYear        int     `json:"duration_year" csv:"duration_year"`
	CentralExposure     float64 `json:"central_exposure" csv:"central_exposure"`
	ActualLapses        float64 `json:"actual_lapses" csv:"actual_lapses"`
	ExpectedLapses      float64 `json:"expected_lapses" csv:"expected_lapses"`
	ActualUx            float64 `json:"actual_ux" csv:"actual_ux"`
	ExpectedUx          float64 `json:"expected_ux" csv:"expected_ux"`
	ActualAnnualRate    float64 `json:"actual_annual_rate" csv:"actual_annual_rate"`
	ExpectedAnnualRate  float64 `json:"expected_annual_rate" csv:"expected_annual_rate"`
	ActualMonthlyRate   float64 `json:"actual_monthly_rate" csv:"actual_monthly_rate"`
	ExpectedMonthlyRate float64 `json:"expected_monthly_rate" csv:"expected_monthly_rate"`
	TestStatisticZ      float64 `json:"test_statistic_z" csv:"test_statistic_z"`
	DataPointCount      float64 `json:"data_point_count" csv:"data_point_count"`
	CreatedBy           string  `json:"created_by" csv:"created_by"`
	Version             string  `json:"version" csv:"version"`
	Year                int     `json:"year" csv:"year"`
}

type ExpLapseCrudeResultSummary struct { // calculated values from ExpExposureData
	DurationYear        string  `json:"duration_year" csv:"duration_year"`
	CentralExposure     float64 `json:"central_exposure" csv:"central_exposure"`
	ActualLapses        float64 `json:"actual_lapses" csv:"actual_lapses"`
	ExpectedLapses      float64 `json:"expected_lapses" csv:"expected_lapses"`
	ActualUx            float64 `json:"actual_ux" csv:"actual_ux"`
	ExpectedUx          float64 `json:"expected_ux" csv:"expected_ux"`
	ActualAnnualRate    float64 `json:"actual_annual_rate" csv:"actual_annual_rate"`
	ExpectedAnnualRate  float64 `json:"expected_annual_rate" csv:"expected_annual_rate"`
	ActualMonthlyRate   float64 `json:"actual_monthly_rate" csv:"actual_monthly_rate"`
	ExpectedMonthlyRate float64 `json:"expected_monthly_rate" csv:"expected_monthly_rate"`
	TestStatisticZ      float64 `json:"test_statistic_z" csv:"test_statistic_z"`
}

type AGGridColumnDef struct {
	HeaderName     string      `json:"headerName"`
	Field          string      `json:"field"`
	RowGroup       bool        `json:"rowGroup,omitempty"` // Will not be used in this version
	Hide           bool        `json:"hide,omitempty"`     // Will not be used in this version
	Pinned         string      `json:"pinned,omitempty"`
	MinWidth       int         `json:"minWidth,omitempty"`
	ValueFormatter string      `json:"valueFormatter,omitempty"`
	CellStyle      interface{} `json:"cellStyle,omitempty"`
}

type APIResponse struct {
	ColumnDefs                   []AGGridColumnDef        `json:"columnDefs"`
	MaleRowData                  []map[string]interface{} `json:"maleRowData"`
	FemaleRowData                []map[string]interface{} `json:"femaleRowData"`
	CombinedRowData              []map[string]interface{} `json:"combinedRowData"`
	LiveExposureColumnDefs       []AGGridColumnDef        `json:"liveExposureColumnDefs,omitempty"`
	LiveExposureRowData          []map[string]interface{} `json:"liveExposureRowData,omitempty"`
	GenderShareColumnDefs        []AGGridColumnDef        `json:"genderShareColumnDefs,omitempty"`
	GenderShareRowData           []map[string]interface{} `json:"genderShareRowData,omitempty"`
	AgeBandExposureColumnDefs    []AGGridColumnDef        `json:"ageBandExposureColumnDefs,omitempty"`
	AgeBandExposureRowData       []map[string]interface{} `json:"ageBandExposureRowData,omitempty"`
	AgeBandShareColumnDefs       []AGGridColumnDef        `json:"ageBandShareColumnDefs,omitempty"`
	AgeBandShareRowData          []map[string]interface{} `json:"ageBandShareRowData,omitempty"`
}

type PeriodOutputData struct {
	PeriodLabel string
	Male        TableRowData
	Female      TableRowData
	Combined    TableRowData
}

type TableRowData struct {
	LivesExposure              float64
	AELives                    float64
	ExpectedClaims             float64
	ActualClaims               float64
	CrudeRateLives             float64
	ExpectedCrudeMortalityRate float64
	ActualCrudeMortalityRate   float64
	AELivesBasis               float64
}

// TableRowOutput represents a single row in the output table (for AG-Grid)
type TableRowOutput struct {
	Duration                  string  `json:"duration"`
	ActualLapses              float64 `json:"actual_lapses"`
	ActualExposure            float64 `json:"actual_exposure"`
	ActualUX                  float64 `json:"actual_ux"`
	ActualAnnualRate          string  `json:"actual_annual_rate"`         // Formatted as "X.XX%"
	ActualActualMonthlyRate   string  `json:"actual_actual_monthly_rate"` // Formatted as "X.XX%"
	ExpectedLapses            float64 `json:"expected_lapses"`
	ExpectedExposure          float64 `json:"expected_exposure"`
	ExpectedUX                float64 `json:"expected_ux"`
	ExpectedAnnualRate        string  `json:"expected_annual_rate"`         // Formatted as "X.XX%"
	ExpectedActualMonthlyRate string  `json:"expected_actual_monthly_rate"` // Formatted as "X.XX%"
}

// TransformedTableSet holds the data for both Actual and Expected tables for one input record
type TransformedTableSet struct {
	ActualData   []TableRowOutput `json:"actual_data"`
	ExpectedData []TableRowOutput `json:"expected_data"`
}
