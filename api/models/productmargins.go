package models

type ProductMargins struct {
	ID                 int     `gorm:"primary_key"`
	Year               int     `json:"year"  csv:"year"`
	ProductCode        string  `json:"product_code"  csv:"product_code"`
	MortalityMargin    float64 `json:"mortality_margin" csv:"mortality_margin"`
	MorbidityMargin    float64 `json:"morbidity_margin" csv:"morbidity_margin"`
	ExpenseMargin      float64 `json:"expense_margin" csv:"expense_margin"`
	RetrenchmentMargin float64 `json:"retrenchment_margin" csv:"retrenchment_margin"`
	InflationMargin    float64 `json:"inflation_margin" csv:"inflation_margin"`
	InvestmentMargin   float64 `json:"investment_margin" csv:"investment_margin"`
	Basis              string  `json:"basis" csv:"basis"`
	Created            int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}
