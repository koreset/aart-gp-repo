package services

import (
	"api/models"
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"
)

// RiskFlags captures the per-line-item risk signals that finance reviews
// before authorising payment. Stored on ClaimPaymentScheduleItem.RiskFlags as
// JSON so new flags can be added without a schema change.
type RiskFlags struct {
	BankingChange30d    bool   `json:"banking_change_30d"`
	Contestable         bool   `json:"contestable"`
	RecentReinstatement bool   `json:"recent_reinstatement"`
	FraudRiskLevel      string `json:"fraud_risk_level"`
	// Cross-claim duplicate signals (Phase 5). Each Refs slice carries a
	// human-readable list of prior claim references the UI can render in a
	// tooltip — e.g. "CLM-12345 · R 5,000 · paid". Truncated server-side to
	// the most recent few entries so the JSON column stays small.
	IDPaidBefore          bool     `json:"id_paid_before"`
	IDPaidBeforeRefs      []string `json:"id_paid_before_refs,omitempty"`
	AccountUsedBefore     bool     `json:"account_used_before"`
	AccountUsedBeforeRefs []string `json:"account_used_before_refs,omitempty"`
}

// ComputeRiskFlags derives the risk flags for a claim at schedule-generation
// time. Best-effort: a missing signal returns the safe default (false / "")
// rather than erroring. The audit/policy data we don't have yet is left as
// false so finance still sees the row but isn't blocked.
func ComputeRiskFlags(claim models.GroupSchemeClaim, db *gorm.DB) RiskFlags {
	if db == nil {
		db = DB
	}
	flags := RiskFlags{}

	flags.BankingChange30d = bankingChangedRecently(claim, db)
	flags.Contestable = isContestable(claim)
	// recent_reinstatement requires a policy-admin integration that isn't
	// in scope for Phase 1; default to false until that lands.
	flags.RecentReinstatement = false
	flags.FraudRiskLevel = latestFraudRiskLevel(claim.ID, db)

	// Phase 5: cross-claim duplicate signals. Best-effort — if the lookup
	// errors (e.g. transient DB blip) we leave the flags clear rather than
	// blocking the schedule generation.
	if idHits, accountHits, err := CheckCrossClaimDuplicates(db, claim); err == nil {
		flags = DecorateRiskFlagsWithCrossClaim(flags, idHits, accountHits)
	}

	return flags
}

// MarshalRiskFlags serialises flags to the JSON byte form persisted on the
// schedule item.
func MarshalRiskFlags(f RiskFlags) models.JSON {
	b, err := json.Marshal(f)
	if err != nil {
		return models.JSON("{}")
	}
	return models.JSON(b)
}

// UnmarshalRiskFlags is the inverse of MarshalRiskFlags. Returns a zero-value
// RiskFlags (all defaults) when the input is empty or malformed so callers
// don't need to plumb errors through every read site.
func UnmarshalRiskFlags(raw models.JSON) RiskFlags {
	var f RiskFlags
	if len(raw) == 0 {
		return f
	}
	if err := json.Unmarshal([]byte(raw), &f); err != nil {
		return RiskFlags{}
	}
	return f
}

// bankingChangedRecently looks for a status audit row mentioning a banking
// change within the last 30 days. Until a dedicated banking-detail audit
// table exists, this is a best-effort string match — false negatives are
// preferable to false positives here.
func bankingChangedRecently(claim models.GroupSchemeClaim, db *gorm.DB) bool {
	cutoff := time.Now().AddDate(0, 0, -30)
	var count int64
	err := db.Model(&models.GroupSchemeClaimStatusAudit{}).
		Where("claim_id = ? AND changed_at >= ?", claim.ID, cutoff).
		Where("LOWER(status_message) LIKE ?", "%bank%").
		Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

// isContestable returns true when the claim event occurred within the
// two-year contestability window of the policy registration. We use
// DateRegistered on the claim as a proxy for the policy inception until
// a richer policy-data link exists.
func isContestable(claim models.GroupSchemeClaim) bool {
	const layout = "2006-01-02"
	eventStr := strings.TrimSpace(claim.DateOfEvent)
	regStr := strings.TrimSpace(claim.DateRegistered)
	if eventStr == "" || regStr == "" {
		return false
	}
	event, err1 := time.Parse(layout, eventStr)
	reg, err2 := time.Parse(layout, regStr)
	if err1 != nil || err2 != nil {
		return false
	}
	// Event within 2 years of registration → contestable.
	return event.Sub(reg) < 2*365*24*time.Hour
}

// latestFraudRiskLevel returns the most recent FraudRiskLevel value from
// this claim's assessments, or "" if none recorded.
func latestFraudRiskLevel(claimID int, db *gorm.DB) string {
	var assessment models.GroupSchemeClaimAssessment
	err := db.Where("claim_id = ?", claimID).
		Order("creation_date DESC").
		Select("id, fraud_risk_level").
		First(&assessment).Error
	if err != nil {
		return ""
	}
	return assessment.FraudRiskLevel
}
