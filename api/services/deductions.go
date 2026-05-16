package services

import (
	"api/models"

	"gorm.io/gorm"
)

// Deductions captures the per-line deductions applied to a claim's gross
// amount before payment: outstanding premium arrears, policy loans, and
// tax withheld (IT3(a) implications for SA life claims). Each comes from
// a different source system and is computed best-effort here so the
// schedule line carries a complete breakdown for finance to verify.
type Deductions struct {
	PremiumArrears float64 `json:"premium_arrears"`
	PolicyLoan     float64 `json:"policy_loan"`
	TaxWithheld    float64 `json:"tax_withheld"`
}

// Total returns the sum of all deductions.
func (d Deductions) Total() float64 {
	return d.PremiumArrears + d.PolicyLoan + d.TaxWithheld
}

// ComputeDeductions returns the deductions applicable to a single claim.
//
// Phase 1 intentionally returns zero-value deductions: the premium ledger
// and policy-loan ledger lookups depend on data feeds that aren't wired in
// yet. The structure is in place so future integrations (or manual entry
// while the schedule is still in `draft`) can populate the values without
// any schema change.
//
// db is accepted so a transaction handle can be passed in once real lookups
// land; the no-op implementation ignores it.
func ComputeDeductions(claim models.GroupSchemeClaim, db *gorm.DB) Deductions {
	_ = db
	return Deductions{}
}
