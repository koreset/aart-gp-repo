package models

import "time"

// ReinsuranceTreaty defines the terms of a reinsurance contract
type ReinsuranceTreaty struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement"`
	// TreatyNumber is unique per Line of Business — same number may exist for a different LoB
	TreatyNumber  string `json:"treaty_number" gorm:"uniqueIndex:treaty_lob_unique;size:100"`
	TreatyName    string `json:"treaty_name" gorm:"default:''"`
	ReinsurerName string `json:"reinsurer_name" gorm:"default:''"`
	ReinsurerCode string `json:"reinsurer_code" gorm:"default:''"`
	BrokerName    string `json:"broker_name" gorm:"default:''"`
	// TreatyType: proportional | xl_risk | xl_event | stop_loss | catastrophe_xl
	TreatyType string `json:"treaty_type" gorm:"default:'proportional'"`
	// LineOfBusiness: group_life | group_disability | funeral | group_health | credit_life
	// Part of composite unique index with TreatyNumber
	LineOfBusiness string `json:"line_of_business" gorm:"uniqueIndex:treaty_lob_unique;default:'group_life'"`
	EffectiveDate  string `json:"effective_date" gorm:"default:''"`
	ExpiryDate     string `json:"expiry_date" gorm:"default:''"`
	RenewalDate    string `json:"renewal_date" gorm:"default:''"`
	// Status: draft | active | expired | cancelled | under_negotiation
	Status   string `json:"status" gorm:"default:'draft'"`
	Currency string `json:"currency" gorm:"default:'ZAR'"`
	// RetentionType: fixed_amount | percentage
	RetentionType        string  `json:"retention_type" gorm:"default:'percentage'"`
	RetentionAmount      float64 `json:"retention_amount" gorm:"default:0"`
	RetentionPercentage  float64 `json:"retention_percentage" gorm:"default:0"`
	SurplusLines         int     `json:"surplus_lines" gorm:"default:0"`
	XLRetention          float64 `json:"xl_retention" gorm:"default:0"`
	XLLimit              float64 `json:"xl_limit" gorm:"default:0"`
	XLLayerFrom          float64 `json:"xl_layer_from" gorm:"default:0"`
	XLLayerTo            float64 `json:"xl_layer_to" gorm:"default:0"`
	AggregateAnnualLimit float64 `json:"aggregate_annual_limit" gorm:"default:0"`
	// Three-Tier Reinsurance Structure (Sum Assured Tiers)
	TreatyCode                string  `json:"treaty_code" gorm:"default:''"`
	RiskPremiumBasisIndicator string  `json:"risk_premium_basis_indicator" gorm:"default:'default'"`
	FlatAnnualReinsPremRate   float64 `json:"flat_annual_reins_prem_rate" gorm:"default:0"`
	Level1CededProportion     float64 `json:"level1_ceded_proportion" gorm:"default:0"`
	Level1Lowerbound          float64 `json:"level1_lowerbound" gorm:"default:0"`
	Level1Upperbound          float64 `json:"level1_upperbound" gorm:"default:0"`
	Level2CededProportion     float64 `json:"level2_ceded_proportion" gorm:"default:0"`
	Level2Lowerbound          float64 `json:"level2_lowerbound" gorm:"default:0"`
	Level2Upperbound          float64 `json:"level2_upperbound" gorm:"default:0"`
	Level3CededProportion     float64 `json:"level3_ceded_proportion" gorm:"default:0"`
	Level3Lowerbound          float64 `json:"level3_lowerbound" gorm:"default:0"`
	Level3Upperbound          float64 `json:"level3_upperbound" gorm:"default:0"`
	// Income-Based Tiers
	IncomeLevel1CededProportion float64 `json:"income_level1_ceded_proportion" gorm:"default:0"`
	IncomeLevel1Lowerbound      float64 `json:"income_level1_lowerbound" gorm:"default:0"`
	IncomeLevel1Upperbound      float64 `json:"income_level1_upperbound" gorm:"default:0"`
	IncomeLevel2CededProportion float64 `json:"income_level2_ceded_proportion" gorm:"default:0"`
	IncomeLevel2Lowerbound      float64 `json:"income_level2_lowerbound" gorm:"default:0"`
	IncomeLevel2Upperbound      float64 `json:"income_level2_upperbound" gorm:"default:0"`
	IncomeLevel3CededProportion float64 `json:"income_level3_ceded_proportion" gorm:"default:0"`
	IncomeLevel3Lowerbound      float64 `json:"income_level3_lowerbound" gorm:"default:0"`
	IncomeLevel3Upperbound      float64 `json:"income_level3_upperbound" gorm:"default:0"`
	// Multi-Reinsurer Structure
	LeadReinsurerShare        float64 `json:"lead_reinsurer_share" gorm:"default:0"`
	LeadReinsurerCode         string  `json:"lead_reinsurer_code" gorm:"default:''"`
	NonLeadReinsurer1Share    float64 `json:"non_lead_reinsurer1_share" gorm:"default:0"`
	NonLeadReinsurer1Code     string  `json:"non_lead_reinsurer1_code" gorm:"default:''"`
	NonLeadReinsurer2Share    float64 `json:"non_lead_reinsurer2_share" gorm:"default:0"`
	NonLeadReinsurer2Code     string  `json:"non_lead_reinsurer2_code" gorm:"default:''"`
	NonLeadReinsurer3Share    float64 `json:"non_lead_reinsurer3_share" gorm:"default:0"`
	NonLeadReinsurer3Code     string  `json:"non_lead_reinsurer3_code" gorm:"default:''"`
	CedingCommission          float64 `json:"ceding_commission" gorm:"default:0"`
	ProfitCommissionRate      float64 `json:"profit_commission_rate" gorm:"default:0"`
	ReinsuranceCommissionRate float64 `json:"reinsurance_commission_rate" gorm:"default:0"`
	// PremiumPaymentFrequency: monthly | quarterly | annual
	PremiumPaymentFrequency string  `json:"premium_payment_frequency" gorm:"default:'monthly'"`
	ClaimsNotificationDays  int     `json:"claims_notification_days" gorm:"default:30"`
	LargeClaimsThreshold    float64 `json:"large_claims_threshold" gorm:"default:0"`
	// DeltaReporting: true = send only changed members; false = full census
	DeltaReporting bool `json:"delta_reporting" gorm:"default:false"`
	// IsRunOff: treaty is in run-off (no new business, run-off processing cycle applies)
	IsRunOff        bool      `json:"is_run_off" gorm:"default:false"`
	RunOffStartDate string    `json:"run_off_start_date" gorm:"default:''"`
	Notes           string    `json:"notes" gorm:"type:text"`
	CreatedBy       string    `json:"created_by" gorm:"default:''"`
	UpdatedBy       string    `json:"updated_by" gorm:"default:''"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TreatySchemeLink links a scheme to a reinsurance treaty with optional cession override
type TreatySchemeLink struct {
	ID              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	TreatyID        int       `json:"treaty_id" gorm:"index"`
	SchemeID        int       `json:"scheme_id" gorm:"index"`
	SchemeName      string    `json:"scheme_name" gorm:"default:''"`
	CessionOverride float64   `json:"cession_override" gorm:"default:0"`
	EffectiveDate   string    `json:"effective_date" gorm:"default:''"`
	ExpiryDate      string    `json:"expiry_date" gorm:"default:''"`
	CreatedBy       string    `json:"created_by" gorm:"default:''"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateTreatyRequest struct {
	TreatyNumber              string  `json:"treaty_number" binding:"required"`
	TreatyName                string  `json:"treaty_name" binding:"required"`
	ReinsurerName             string  `json:"reinsurer_name" binding:"required"`
	ReinsurerCode             string  `json:"reinsurer_code"`
	BrokerName                string  `json:"broker_name"`
	TreatyType                string  `json:"treaty_type" binding:"required"`
	LineOfBusiness            string  `json:"line_of_business"`
	EffectiveDate             string  `json:"effective_date" binding:"required"`
	ExpiryDate                string  `json:"expiry_date" binding:"required"`
	RenewalDate               string  `json:"renewal_date"`
	Currency                  string  `json:"currency"`
	RetentionType             string  `json:"retention_type"`
	RetentionAmount           float64 `json:"retention_amount"`
	RetentionPercentage       float64 `json:"retention_percentage"`
	SurplusLines              int     `json:"surplus_lines"`
	XLRetention               float64 `json:"xl_retention"`
	XLLimit                   float64 `json:"xl_limit"`
	XLLayerFrom               float64 `json:"xl_layer_from"`
	XLLayerTo                 float64 `json:"xl_layer_to"`
	AggregateAnnualLimit      float64 `json:"aggregate_annual_limit"`
	ProfitCommissionRate      float64 `json:"profit_commission_rate"`
	ReinsuranceCommissionRate float64 `json:"reinsurance_commission_rate"`
	PremiumPaymentFrequency   string  `json:"premium_payment_frequency"`
	ClaimsNotificationDays    int     `json:"claims_notification_days"`
	LargeClaimsThreshold      float64 `json:"large_claims_threshold"`
	DeltaReporting            bool    `json:"delta_reporting"`
	IsRunOff                  bool    `json:"is_run_off"`
	RunOffStartDate           string  `json:"run_off_start_date"`
	Notes                     string  `json:"notes"`
	// Three-Tier Reinsurance Structure
	TreatyCode                  string  `json:"treaty_code"`
	RiskPremiumBasisIndicator   string  `json:"risk_premium_basis_indicator"`
	FlatAnnualReinsPremRate     float64 `json:"flat_annual_reins_prem_rate"`
	Level1CededProportion       float64 `json:"level1_ceded_proportion"`
	Level1Lowerbound            float64 `json:"level1_lowerbound"`
	Level1Upperbound            float64 `json:"level1_upperbound"`
	Level2CededProportion       float64 `json:"level2_ceded_proportion"`
	Level2Lowerbound            float64 `json:"level2_lowerbound"`
	Level2Upperbound            float64 `json:"level2_upperbound"`
	Level3CededProportion       float64 `json:"level3_ceded_proportion"`
	Level3Lowerbound            float64 `json:"level3_lowerbound"`
	Level3Upperbound            float64 `json:"level3_upperbound"`
	IncomeLevel1CededProportion float64 `json:"income_level1_ceded_proportion"`
	IncomeLevel1Lowerbound      float64 `json:"income_level1_lowerbound"`
	IncomeLevel1Upperbound      float64 `json:"income_level1_upperbound"`
	IncomeLevel2CededProportion float64 `json:"income_level2_ceded_proportion"`
	IncomeLevel2Lowerbound      float64 `json:"income_level2_lowerbound"`
	IncomeLevel2Upperbound      float64 `json:"income_level2_upperbound"`
	IncomeLevel3CededProportion float64 `json:"income_level3_ceded_proportion"`
	IncomeLevel3Lowerbound      float64 `json:"income_level3_lowerbound"`
	IncomeLevel3Upperbound      float64 `json:"income_level3_upperbound"`
	LeadReinsurerShare          float64 `json:"lead_reinsurer_share"`
	LeadReinsurerCode           string  `json:"lead_reinsurer_code"`
	NonLeadReinsurer1Share      float64 `json:"non_lead_reinsurer1_share"`
	NonLeadReinsurer1Code       string  `json:"non_lead_reinsurer1_code"`
	NonLeadReinsurer2Share      float64 `json:"non_lead_reinsurer2_share"`
	NonLeadReinsurer2Code       string  `json:"non_lead_reinsurer2_code"`
	NonLeadReinsurer3Share      float64 `json:"non_lead_reinsurer3_share"`
	NonLeadReinsurer3Code       string  `json:"non_lead_reinsurer3_code"`
	CedingCommission            float64 `json:"ceding_commission"`
}

type UpdateTreatyRequest struct {
	TreatyName                string  `json:"treaty_name"`
	ReinsurerName             string  `json:"reinsurer_name"`
	ReinsurerCode             string  `json:"reinsurer_code"`
	BrokerName                string  `json:"broker_name"`
	TreatyType                string  `json:"treaty_type"`
	LineOfBusiness            string  `json:"line_of_business"`
	EffectiveDate             string  `json:"effective_date"`
	ExpiryDate                string  `json:"expiry_date"`
	RenewalDate               string  `json:"renewal_date"`
	Status                    string  `json:"status"`
	Currency                  string  `json:"currency"`
	RetentionType             string  `json:"retention_type"`
	RetentionAmount           float64 `json:"retention_amount"`
	RetentionPercentage       float64 `json:"retention_percentage"`
	SurplusLines              int     `json:"surplus_lines"`
	XLRetention               float64 `json:"xl_retention"`
	XLLimit                   float64 `json:"xl_limit"`
	XLLayerFrom               float64 `json:"xl_layer_from"`
	XLLayerTo                 float64 `json:"xl_layer_to"`
	AggregateAnnualLimit      float64 `json:"aggregate_annual_limit"`
	ProfitCommissionRate      float64 `json:"profit_commission_rate"`
	ReinsuranceCommissionRate float64 `json:"reinsurance_commission_rate"`
	PremiumPaymentFrequency   string  `json:"premium_payment_frequency"`
	ClaimsNotificationDays    int     `json:"claims_notification_days"`
	LargeClaimsThreshold      float64 `json:"large_claims_threshold"`
	DeltaReporting            bool    `json:"delta_reporting"`
	IsRunOff                  bool    `json:"is_run_off"`
	RunOffStartDate           string  `json:"run_off_start_date"`
	Notes                     string  `json:"notes"`
	// Three-Tier Reinsurance Structure
	TreatyCode                  string  `json:"treaty_code"`
	RiskPremiumBasisIndicator   string  `json:"risk_premium_basis_indicator"`
	FlatAnnualReinsPremRate     float64 `json:"flat_annual_reins_prem_rate"`
	Level1CededProportion       float64 `json:"level1_ceded_proportion"`
	Level1Lowerbound            float64 `json:"level1_lowerbound"`
	Level1Upperbound            float64 `json:"level1_upperbound"`
	Level2CededProportion       float64 `json:"level2_ceded_proportion"`
	Level2Lowerbound            float64 `json:"level2_lowerbound"`
	Level2Upperbound            float64 `json:"level2_upperbound"`
	Level3CededProportion       float64 `json:"level3_ceded_proportion"`
	Level3Lowerbound            float64 `json:"level3_lowerbound"`
	Level3Upperbound            float64 `json:"level3_upperbound"`
	IncomeLevel1CededProportion float64 `json:"income_level1_ceded_proportion"`
	IncomeLevel1Lowerbound      float64 `json:"income_level1_lowerbound"`
	IncomeLevel1Upperbound      float64 `json:"income_level1_upperbound"`
	IncomeLevel2CededProportion float64 `json:"income_level2_ceded_proportion"`
	IncomeLevel2Lowerbound      float64 `json:"income_level2_lowerbound"`
	IncomeLevel2Upperbound      float64 `json:"income_level2_upperbound"`
	IncomeLevel3CededProportion float64 `json:"income_level3_ceded_proportion"`
	IncomeLevel3Lowerbound      float64 `json:"income_level3_lowerbound"`
	IncomeLevel3Upperbound      float64 `json:"income_level3_upperbound"`
	LeadReinsurerShare          float64 `json:"lead_reinsurer_share"`
	LeadReinsurerCode           string  `json:"lead_reinsurer_code"`
	NonLeadReinsurer1Share      float64 `json:"non_lead_reinsurer1_share"`
	NonLeadReinsurer1Code       string  `json:"non_lead_reinsurer1_code"`
	NonLeadReinsurer2Share      float64 `json:"non_lead_reinsurer2_share"`
	NonLeadReinsurer2Code       string  `json:"non_lead_reinsurer2_code"`
	NonLeadReinsurer3Share      float64 `json:"non_lead_reinsurer3_share"`
	NonLeadReinsurer3Code       string  `json:"non_lead_reinsurer3_code"`
	CedingCommission            float64 `json:"ceding_commission"`
}

type LinkSchemeRequest struct {
	SchemeID        int     `json:"scheme_id" binding:"required"`
	SchemeName      string  `json:"scheme_name"`
	CessionOverride float64 `json:"cession_override"`
	EffectiveDate   string  `json:"effective_date"`
	ExpiryDate      string  `json:"expiry_date"`
}

// BulkLinkSchemesRequest links multiple schemes to a treaty at once
type BulkLinkSchemesRequest struct {
	SchemeIDs       []int   `json:"scheme_ids" binding:"required,min=1"`
	CessionOverride float64 `json:"cession_override"`
	EffectiveDate   string  `json:"effective_date"`
}

type TreatyStats struct {
	Total            int `json:"total"`
	Active           int `json:"active"`
	Draft            int `json:"draft"`
	Expired          int `json:"expired"`
	ExpiringIn60Days int `json:"expiring_in_60_days"`
}
