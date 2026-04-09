package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// ProductFamily example
type ProductFamily struct {
	ID       int       `gorm:"primary_key" example:"1" json:"id"`
	Name     string    `json:"name" example:"Health" json:"name"`
	Code     string    `json:"code" example:"HEAL" json:"code"`
	Products []Product `json:"products"`
}

type Product struct {
	ID                         int                         `json:"id" gorm:"primary_key"`
	ProductFamilyID            int                         `json:"product_family_id"`
	ProductName                string                      `json:"product_name" gorm:"unique"`
	ProductCode                string                      `json:"product_code" gorm:"unique;size:255"`
	AssumptionCode             string                      `json:"assumption_code"`
	ProductTransitionStates    []ProductTransitionState    `json:"product_transition_states"`
	ProductTransitions         []ProductTransition         `json:"product_transitions"`
	ProductRatingFactors       []ProductRatingFactor       `json:"product_rating_factors" gorm:"constraint:OnDelete:CASCADE"`
	ProductModelpointVariables []ProductModelpointVariable `json:"product_modelpoint_variables"`
	ProductFeatures            ProductFeatures             `json:"product_features" gorm:"foreignKey:ProductCode;references:product_code"`
	GlobalTables               []GlobalTable               `json:"global_tables"`
	ParametersArray            []string                    `json:"parameters_array" gorm:"-"`
	ProductParameters          ProductParameters           `json:"product_parameters" gorm:"foreignKey:ProductCode;references:product_code"`
	ProductTables              []ProductTable              `json:"product_tables"`
	ProductPricingTables       []ProductPricingTable       `json:"product_pricing_tables" gorm:"foreignKey:ProductID;references:id"`
	FeatureList                []string                    `json:"feature_list" gorm:"-"`
	Status                     string                      `json:"product_state"`
	CreatedAt                  time.Time                   `json:"created_at"`
	UpdatedAt                  time.Time                   `json:"updated_at"`
	CreatedBy                  string                      `json:"created_by"`
}

type GlobalTable struct {
	ID        int    `json:"id" gorm:"primary_key"`
	Name      string `json:"name"`
	ProductID int    `json:"product_id"`
}
type ProductTransitionState struct {
	ID        int    `json:"-" gorm:"primary_key"`
	State     string `json:"state" form:"state"`
	StateId   int    `json:"id"`
	ProductID int    `json:"product_id"`
}
type ProductTransition struct {
	ID              int    `json:"id" gorm:"primary_key"`
	ProductID       int    `json:"product_id"`
	Absorbing       bool   `json:"absorbing"`
	StartState      string `json:"start_state"`
	EndState        string `json:"end_state"`
	AssociatedTable string `json:"associated_table"`
}

type Fds struct {
	ID                    int    `json:"id" gorm:"primary_key"`
	ProductRatingFactorID int    `json:"product_rating_factor_id"`
	Factor                string `json:"factor"`
	Dimension             string `json:"dimension"`
}
type ProductRatingFactor struct {
	ID              int    `json:"id" gorm:"primary_key"`
	ProductID       int    `json:"product_id"`
	Fds             []Fds  `json:"fds" gorm:"constraint:OnDelete:SET NULL"`
	TransitionTable string `json:"table"`
}
type ProductModelpointVariable struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Type      string `json:"type"`
}

type ProductTable struct {
	ID        int    `gorm:"primary_key" json:"id"`
	ProductID int    `json:"product_id"`
	Class     string `json:"table_class"`
	Name      string `json:"table"`
	Populated bool   `json:"populated" gorm:"-"`
}

type ProductChildSumAssured struct {
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

func (ProductChildSumAssured) TableName() string {
	return "product_child_sum_assured"
}

type ProductChildAdditionalSumAssured struct {
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

func (ProductChildAdditionalSumAssured) TableName() string {
	return "product_child_additional_sum_assured"
}

type ProductAdditionalSumAssured struct {
	Year        int
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

func (ProductAdditionalSumAssured) TableName() string {
	return "product_additional_sum_assured"
}

type ProductRider struct {
	Year         int
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

type ProductClawback struct {
	ProductCode            string  `json:"product_code"`
	DurationInForceMonth   int     `json:"duration_in_force_month"`
	Year1InitialCommission float64 `json:"year_1_initial_commission"`
	Year2InitialCommission float64 `json:"year_2_initial_commission"`
}

func (ProductClawback) TableName() string {
	return "product_clawback"
}

type ProductLapseMargin struct {
	Year        int     `json:"year"`
	ProductCode string  `json:"product_code" gorm:"size:255"`
	Month       int     `json:"month"`
	Margin      float64 `json:"margin"`
	Basis       string  `json:"basis"`
}

func (prod *Product) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Where("product_id=?", prod.ID).Delete(&ProductModelpointVariable{})
	return
}

type ProductAccidentalBenefitMultiplier struct {
	ProductCode    string  `json:"product_code"`
	MainMember     float64 `json:"main_member"`
	Spouse         float64 `json:"spouse"`
	Child          float64 `json:"child"`
	Parent         float64 `json:"parent"`
	ExtendedFamily float64 `json:"extended_family"`
}

func (ProductAccidentalBenefitMultiplier) TableName() string {
	return "product_accident_benefit_multiplier"
}

type ProductReinsurance struct {
	ID                      int     `json:"id" csv:"id" gorm:"primary_key"`
	Year                    int     `json:"year"`
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
}

func (ProductReinsurance) TableName() string {
	return "product_reinsurance"
}

type Columnname struct {
	MortalityColumnName           string
	AccidentalMortalityColumnName string
	LapseColumnName               string
	RetrenchmentColumnName        string
	DisabilityColumnName          string
}

type PricingProductTableNames struct {
	MortalityTableName            string
	MortalityColumnName           string
	MortalityAccidentalTableName  string
	MortalityAccidentalColumnName string
	LapseTableName                string
	LapseColumnName               string
	DisabilityTableName           string
	DisabilityColumnName          string
	RetrenchmentTableName         string
	RetrenchmentColumnName        string
	LapseMonthCount               int
	LapseMarginMonthCount         int
	RetrenchmentTableRowCount     int
}
