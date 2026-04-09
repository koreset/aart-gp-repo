package models

import "time"

// TechnicalAccount is the formal settlement statement exchanged between cedant and reinsurer
type TechnicalAccount struct {
	ID                    int     `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountNumber         string  `json:"account_number" gorm:"uniqueIndex;size:100"`
	TreatyID              int     `json:"treaty_id" gorm:"index"`
	TreatyNumber          string  `json:"treaty_number" gorm:"default:''"`
	ReinsurerName         string  `json:"reinsurer_name" gorm:"default:''"`
	PeriodStart           string  `json:"period_start" gorm:"default:''"`
	PeriodEnd             string  `json:"period_end" gorm:"default:''"`
	CededPremiumEarned    float64 `json:"ceded_premium_earned" gorm:"default:0"`
	ReinsuranceCommission float64 `json:"reinsurance_commission" gorm:"default:0"`
	NetCededPremium       float64 `json:"net_ceded_premium" gorm:"default:0"`
	CededClaimsPaid       float64 `json:"ceded_claims_paid" gorm:"default:0"`
	CededIBNR             float64 `json:"ceded_ibnr" gorm:"default:0"`
	TotalCededClaims      float64 `json:"total_ceded_claims" gorm:"default:0"`
	ProfitCommission      float64 `json:"profit_commission" gorm:"default:0"`
	// NetBalance: positive = cedant owes reinsurer; negative = reinsurer owes cedant
	NetBalance float64 `json:"net_balance" gorm:"default:0"`
	// Status: draft | issued | agreed | settled | disputed
	Status       string     `json:"status" gorm:"default:'draft'"`
	IssuedAt     *time.Time `json:"issued_at"`
	AgreedAt     *time.Time `json:"agreed_at"`
	SettledAt    *time.Time `json:"settled_at"`
	DisputeNotes string     `json:"dispute_notes" gorm:"type:text"`
	// DisputeStatus: none | stage1 | stage2 | stage3 | resolved
	DisputeStatus      string     `json:"dispute_status" gorm:"default:'none'"`
	DisputeEscalatedAt *time.Time `json:"dispute_escalated_at"`
	Notes              string     `json:"notes" gorm:"type:text"`
	CreatedBy          string     `json:"created_by" gorm:"default:''"`
	UpdatedBy          string     `json:"updated_by" gorm:"default:''"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// SettlementPayment records a payment against a technical account
type SettlementPayment struct {
	ID                 int     `json:"id" gorm:"primaryKey;autoIncrement"`
	TechnicalAccountID int     `json:"technical_account_id" gorm:"index"`
	PaymentDate        string  `json:"payment_date" gorm:"default:''"`
	Amount             float64 `json:"amount" gorm:"default:0"`
	// Direction: cedant_to_ri | ri_to_cedant
	Direction     string    `json:"direction" gorm:"default:''"`
	PaymentMethod string    `json:"payment_method" gorm:"default:''"`
	Reference     string    `json:"reference" gorm:"default:''"`
	Notes         string    `json:"notes" gorm:"type:text"`
	RecordedBy    string    `json:"recorded_by" gorm:"default:''"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GenerateTechnicalAccountRequest struct {
	TreatyID    int    `json:"treaty_id" binding:"required"`
	PeriodStart string `json:"period_start" binding:"required"`
	PeriodEnd   string `json:"period_end" binding:"required"`
	Notes       string `json:"notes"`
}

type UpdateTechnicalAccountRequest struct {
	Status       string  `json:"status"`
	DisputeNotes string  `json:"dispute_notes"`
	Notes        string  `json:"notes"`
	CededIBNR    float64 `json:"ceded_ibnr"`
}

type RecordSettlementPaymentRequest struct {
	TechnicalAccountID int     `json:"technical_account_id" binding:"required"`
	PaymentDate        string  `json:"payment_date" binding:"required"`
	Amount             float64 `json:"amount" binding:"required"`
	Direction          string  `json:"direction" binding:"required"`
	PaymentMethod      string  `json:"payment_method"`
	Reference          string  `json:"reference"`
	Notes              string  `json:"notes"`
}

type EscalateDisputeRequest struct {
	Stage string `json:"stage" binding:"required"` // stage1 | stage2 | stage3
	Notes string `json:"notes"`
}

type ResolveDisputeRequest struct {
	Notes string `json:"notes"`
}

type SettlementStats struct {
	Total           int     `json:"total"`
	Draft           int     `json:"draft"`
	Issued          int     `json:"issued"`
	Agreed          int     `json:"agreed"`
	Settled         int     `json:"settled"`
	Disputed        int     `json:"disputed"`
	NetOwedToRI     float64 `json:"net_owed_to_ri"`
	NetOwedByCedant float64 `json:"net_owed_by_cedant"`
}
