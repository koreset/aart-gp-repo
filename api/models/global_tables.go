package models

type ProductCommissionStructure struct {
	ProductCode                        string  `json:"product_code"  csv:"product_code" gorm:"primary_key;auto_increment:false;size:191"`
	CommissionType                     string  `json:"commission_type" csv:"commission_type" gorm:"primary_key;auto_increment:false;size:191"`
	InitialCommissionPercentage1       float64 `json:"initial_commission_percentage1" csv:"initial_commission_percentage1"`
	InitialCommissionPercentage2       float64 `json:"initial_commission_percentage2" csv:"initial_commission_percentage2"`
	InitialCommissionRand              float64 `json:"initial_commission_rand" csv:"initial_commission_rand"`
	InitialYear1CommissionPaymentMonth int     `json:"initial_year1_commission_payment_month" csv:"initial_year1_commission_payment_month"`
	InitialYear2CommissionPaymentMonth int     `json:"initial_year2_commission_payment_month" csv:"initial_year2_commission_payment_month"`
	ClawbackPeriod                     int     `json:"clawback_period" csv:"clawback_period"`
	RenewalCommissionPercentage        float64 `json:"renewal_commission_percentage" csv:"renewal_commission_percentage"`
	RenewalCommissionRand              float64 `json:"renewal_commission_rand" csv:"renewal_commission_rand"`
	HybridRenewalCommStartM            int     `json:"hybrid_renewal_comm_start_m" csv:"hybrid_renewal_comm_start_m"`
	HybridRenewalCommEndM              int     `json:"hybrid_renewal_comm_end_m" csv:"hybrid_renewal_comm_end_m"`
	Created                            int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
	CreatedBy                          string  `json:"created_by" csv:"created_by"`
}
