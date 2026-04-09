package models

import "time"

type ModelPointCount struct {
	Year    int    `json:"year"`
	Version string `json:"version"`
	Count   int    `json:"count"`
}

type ProductModelPointCount struct {
	ProductId   int    `json:"product_id"`
	ProductCode string `json:"product_code"`
	Year        int    `json:"year"`
	Version     string `json:"version"`
	Count       int    `json:"count"`
}

type RunJob struct {
	ID             int          `gorm:"primary_key" json:"id"`
	JobsTemplateID int          `json:"jobs_template_id"`
	Jobs           []RunPayload `json:"jobs"`
	UserEmail      string       `json:"user_email"`
	UserName       string       `json:"user_name"`
	User           AppUser      `json:"user"`
}

type ShockSetting struct {
	ID                   int    `json:"id" gorm:"primary_key"`
	Name                 string `json:"name" gorm:"unique"`
	Description          string `json:"description"`
	Mortality            bool   `json:"mortality"`
	Disability           bool   `json:"disability"`
	Lapse                bool   `json:"lapse"`
	CriticalIllness      bool   `json:"critical_illness"`
	RealYieldCurve       bool   `json:"real_yield_curve"`
	NominalYieldCurve    bool   `json:"nominal_yield_curve"`
	Retrenchment         bool   `json:"retrenchment"`
	Expense              bool   `json:"expense"`
	Inflation            bool   `json:"inflation"`
	MortalityCatastrophe bool   `json:"mortality_catastrophe"`
	MorbidityCatastrophe bool   `json:"morbidity_catastrophe"`
	ShockBasis           string `json:"shock_basis"`
	Year                 int    `json:"year"`
}

type RunParameters struct {
	ID                int          `gorm:"primary_key"`
	ProjectionJobID   int          `json:"projection_job_id"`
	RunType           int          `json:"run_type"`
	RunName           string       `json:"run_name"`
	RunSingle         bool         `json:"run_single"`
	RunDate           string       `json:"run_date"`
	UserEmail         string       `json:"user_email"`
	UserName          string       `json:"user_name"`
	ModelpointYear    int          `json:"modelpoint_year"`
	ModelPointVersion string       `json:"mp_version"`
	YieldcurveYear    int          `json:"yieldcurve_year"`
	YieldcurveMonth   int          `json:"yieldcurve_month"`
	YieldCurveBasis   string       `json:"yield_curve_basis"`
	ParameterYear     int          `json:"parameter_year"`
	TransitionYear    int          `json:"transition_year"`
	MorbidityYear     int          `json:"morbidity_year"`
	MortalityYear     int          `json:"mortality_year"`
	LapseYear         int          `json:"lapse_year"`
	LapseMarginYear   int          `json:"lapse_margin_year"`
	RetrenchmentYear  int          `json:"retrenchment_year"`
	RunDescription    string       `json:"run_description"`
	IFRS17Indicator   bool         `json:"ifrs17_indicator"`
	IFRS17Group       bool         `json:"ifrs17group"`
	LICIndicator      bool         `json:"lic_indicator"`
	SpCode            int          `json:"spcode"`
	ShockSettingsID   int          `json:"shock_settings_id"`
	ShockSettings     ShockSetting `json:"shock_settings" gorm:"-"`
	AggregationPeriod int          `json:"aggregation_period"`
	RunBasis          string       `json:"run_basis"`
	YearEndMonth      int          `json:"year_end_month"`
}

type PhiRunParameters struct {
	ID                int             `json:"id" gorm:"primary_key"`
	ProjectionJobID   int             `json:"projection_job_id"`
	RunName           string          `json:"run_name"`
	RunDate           string          `json:"run_date"`
	RunTime           float64         `json:"run_time"`
	UserEmail         string          `json:"user_email"`
	UserName          string          `json:"user_name"`
	ModelPointYear    int             `json:"modelpoint_year"`
	ModelPointVersion string          `json:"modelpoint_version"`
	YieldCurveYear    int             `json:"yield_curve_year"`
	YieldCurveVersion string          `json:"yield_curve_version"`
	ParameterYear     int             `json:"parameter_year"`
	ParameterVersion  string          `json:"parameter_version"`
	RecoveryYear      int             `json:"recovery_year"`
	RecoveryVersion   string          `json:"recovery_version"`
	MortalityYear     int             `json:"mortality_year"`
	MortalityVersion  string          `json:"mortality_version"`
	RunDescription    string          `json:"run_description"`
	ShockSettingsID   int             `json:"shock_settings_id"`
	ShockSettings     PhiShockSetting `json:"shock_settings" gorm:"-"`
	AggregationPeriod int             `json:"aggregation_period"`
	YearEndMonth      int             `json:"year_end_month"`
	RunSingle         bool            `json:"run_single"`
	ModelPointCount   int             `json:"model_point_count"`
	PointsDone        int             `json:"points_done"`
	JobStatus         string          `json:"job_status"`
	Status            string          `json:"status"`
	JobStatusError    string          `json:"job_status_error"`
	CreatedBy         string          `json:"created_by" csv:"created_by"`
	CreationDate      time.Time       `json:"creation_date"`
}

type RunPhiJob struct {
	ID        int                `gorm:"primary_key" json:"id"`
	Jobs      []PhiRunParameters `json:"jobs"`
	UserEmail string             `json:"user_email"`
	UserName  string             `json:"user_name"`
	User      AppUser            `json:"user"`
}

type UserToken struct {
	Subject     string `gorm:"primary_key"`
	TokenString string `gorm:"type:text"`
}

type RunList struct {
	RunDate string `json:"run_date"`
	RunName string `json:"run_name"`
}

type PortfolioRuns struct {
	Name string `json:"name"`
}
