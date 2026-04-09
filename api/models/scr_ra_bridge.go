package models

import "time"

// SCRRABridgeEntry stores manually entered SCR-based risk margin per product/period.
type SCRRABridgeEntry struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductCode   string    `json:"product_code"`
	Period        string    `json:"period"`          // e.g. "2024-12"
	SCRRiskMargin float64   `json:"scr_risk_margin"` // SCR-based risk margin (PA method)
	Notes         string    `json:"notes"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// SCRRABridgeRow is the reconciliation output — RA vs SCR risk margin.
type SCRRABridgeRow struct {
	ProductCode   string  `json:"product_code"`
	Period        string  `json:"period"`
	IFRS17RA      float64 `json:"ifrs17_ra"`
	SCRRiskMargin float64 `json:"scr_risk_margin"`
	Variance      float64 `json:"variance"`    // IFRS17RA - SCRRiskMargin
	VariancePct   float64 `json:"variance_pct"` // Variance / SCRRiskMargin * 100
	Notes         string  `json:"notes"`
}

type CreateSCREntryRequest struct {
	ProductCode   string  `json:"product_code"`
	Period        string  `json:"period"`
	SCRRiskMargin float64 `json:"scr_risk_margin"`
	Notes         string  `json:"notes"`
}
