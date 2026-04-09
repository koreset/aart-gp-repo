package models

type LICClaimsAnalysisOfChange struct {
	ID                       int     `json:"id" gorm:"primary_key;auto_increment:false" csv:"id"`
	ProductCode              string  `json:"product_code" gorm:"primary_key;auto_increment:false" csv:"product_code"`
	UnderwritingYear         int     `json:"underwriting_year" gorm:"primary_key;auto_increment:false" csv:"underwriting_year"`
	ExpectedClaimDevelopment float64 `json:"expected_claim_development" csv:"expected_claim_development"`
	ActualClaimDevelopment   float64 `json:"actual_claim_development" csv:"actual_claim_development"`
}

type Lic2Parameter struct {
	ID                               int     `json:"id" gorm:"primary_key" csv:"id"`
	PortfolioName                    string  `json:"portfolio_name" csv:"portfolio_name"`
	ProductName                      string  `json:"product_name" csv:"product_name"`
	ProductCode                      string  `json:"product_code" csv:"product_code"`
	Basis                            string  `json:"basis" csv:"basis"`
	Treaty1ClaimsProportion          float64 `json:"treaty1_claims_proportion" csv:"treaty1_claims_proportion"`
	Treaty2ClaimsProportion          float64 `json:"treaty2_claims_proportion" csv:"treaty2_claims_proportion"`
	Treaty3ClaimsProportion          float64 `json:"treaty3_claims_proportion" csv:"treaty3_claims_proportion"`
	PrevTreaty1NonPerformanceRate    float64 `json:"prev_treaty1_non_performance_rate" csv:"prev_treaty1_non_performance_rate"`
	CurrTreaty1NonPerformanceRate    float64 `json:"curr_treaty1_non_performance_rate" csv:"curr_treaty1_non_performance_rate"`
	PrevTreaty2NonPerformanceRate    float64 `json:"prev_treaty2_non_performance_rate" csv:"prev_treaty2_non_performance_rate"`
	CurrTreaty2NonPerformanceRate    float64 `json:"curr_treaty2_non_performance_rate" csv:"curr_treaty2_non_performance_rate"`
	PrevTreaty3NonPerformanceRate    float64 `json:"prev_treaty3_non_performance_rate" csv:"prev_treaty3_non_performance_rate"`
	CurrTreaty3NonPerformanceRate    float64 `json:"curr_treaty3_non_performance_rate" csv:"curr_treaty3_non_performance_rate"`
	AdditionalIncurredClaims         float64 `json:"additional_incurred_claims" csv:"additional_incurred_claims"`
	AdditionalOldClaimsPaid          float64 `json:"additional_old_claims_paid" csv:"additional_old_claims_paid"`
	AdditionalNewClaimsPaid          float64 `json:"additional_new_claims_paid" csv:"additional_new_claims_paid"`
	ReportedClaimsEstimateAdjustment float64 `json:"reported_claims_estimate_adjustment" csv:"reported_claims_estimate_adjustment"`
	EstimatedCashbackReserveRelease  float64 `json:"estimated_cashback_reserve_release" csv:"estimated_cashback_reserve_release"`
	CashbackEstimateAdjustment       float64 `json:"cashback_estimate_adjustment" csv:"cashback_estimate_adjustment"`
	Year                             int     `json:"year" csv:"year"`
	Version                          string  `json:"version" csv:"version"`
	Created                          int64   `json:"created" csv:"created" gorm:"autoCreateTime"`
}
