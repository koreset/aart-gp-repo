package models

//type DiscountCurve struct {
//	Year                  int     `json:"year" gorm:"primary_key;auto_increment:false"`
//	ProjectionTime         int     `json:"proj_time" gorm:"primary_key;auto_increment:false;column:proj_time"`
//	NominalRate           float64 `json:"nominal_rate"`
//	RealRate              float64 `json:"real_rate"`
//	Inflation             float64 `json:"inflation"`
//	DividendYield         float64 `json:"dividend_yield"`
//	RentalYield           float64 `json:"rental_yield"`
//	RealisedCapitalGain   float64 `json:"realised_capital_gain"`
//	UnrealisedCapitalGain float64 `json:"unrealised_capital_gain"`
//}

type YieldCurve struct {
	Year           int     `json:"year" csv:"year" gorm:"primary_key;auto_increment:false"`
	YieldCurveCode string  `json:"yield_curve_code" csv:"yield_curve_code" gorm:"primary_key;auto_increment:false;column:yield_curve_code"`
	ProjectionTime int     `json:"proj_time" csv:"proj_time" gorm:"primary_key;auto_increment:false;column:proj_time"`
	Month          int     `json:"month" csv:"month" gorm:"primary_key;auto_increment:false;column:month"`
	NominalRate    float64 `json:"nominal_rate" csv:"nominal_rate"`
	Inflation      float64 `json:"inflation" csv:"inflation"`
	//Basis          string  `json:"basis" csv:"basis" gorm:"primary_key;auto_increment:false;column:basis"`
	Created int64 `json:"created" csv:"created" gorm:"autoCreateTime"`
}

func (YieldCurve) TableName() string {
	return "yield_curve"
}
