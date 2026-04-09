package models

import "time"

type DeferredTaxEntry struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	RunID          int       `json:"run_id"`
	ProductCode    string    `json:"product_code"`
	IFRS17Group    string    `json:"ifrs17_group"`
	CSMAmount      float64   `json:"csm_amount"`
	RAAmount       float64   `json:"ra_amount"`
	LossComponent  float64   `json:"loss_component"`
	TaxRate        float64   `json:"tax_rate"`
	DTLOnCSM       float64   `json:"dtl_on_csm"`
	DTLOnRA        float64   `json:"dtl_on_ra"`
	DTAOnLoss      float64   `json:"dta_on_loss"`
	NetDeferredTax float64   `json:"net_deferred_tax"`
	ComputedAt     time.Time `json:"computed_at"`
}

type DeferredTaxSummary struct {
	RunID          int     `json:"run_id"`
	TaxRate        float64 `json:"tax_rate"`
	TotalCSM       float64 `json:"total_csm"`
	TotalRA        float64 `json:"total_ra"`
	TotalLoss      float64 `json:"total_loss"`
	TotalDTL       float64 `json:"total_dtl"`
	TotalDTA       float64 `json:"total_dta"`
	NetDeferredTax float64 `json:"net_deferred_tax"`
	GroupCount     int     `json:"group_count"`
}

type ComputeDeferredTaxRequest struct {
	RunID   int     `json:"run_id"`
	TaxRate float64 `json:"tax_rate"`
}
