package services

import (
	"fmt"
	"math"
	"strings"
	"time"

	"api/models"

	"gorm.io/gorm"
)

// ──────────────────────────────────────────────
// Pre-authorisation verification suite (Phase 5)
// ──────────────────────────────────────────────
//
// Three pure checks that the schedule UI surfaces to finance reviewers:
//
//   1. Cross-claim duplicate detection — has the same person (by ID) or the
//      same bank account been paid on a different claim?
//   2. Amount drift — does the schedule line's Gross differ from the figure
//      the assessor explicitly approved?
//   3. Banking verification freshness — when was banking last verified and
//      is that result stale?
//
// Each function is pure (DB read-only) so the schedule-creation path can
// compute the data into RiskFlags / snapshot columns, and a separate HTTP
// endpoint can call the same function for live lookups in the UI.

const (
	// amountDriftTolerance is the minimum absolute delta (in currency units)
	// before a non-zero drift is reported. Guards against floating-point
	// noise on identical figures.
	amountDriftTolerance = 1.0

	// bavStalenessWindow is how long a "verified" BAV result stays fresh.
	// Beyond this window the UI surfaces an amber chip prompting a re-verify.
	bavStalenessWindow = 30 * 24 * time.Hour

	// crossClaimDupeRefLimit truncates the historical references stored in
	// RiskFlags so the JSON column stays small. The drawer / popover can
	// fetch the full list via a separate endpoint if needed later.
	crossClaimDupeRefLimit = 5
)

// inFlightOrPaidStatuses are the GroupSchemeClaim.Status values that count as
// "previously paid OR currently being paid" for the cross-claim duplicate
// check. Catching the in-flight ones prevents simultaneous double-payments
// before they reach the bank.
var inFlightOrPaidStatuses = []string{"paid", "submitted_to_bank", "submitted_for_payment"}

// DupeRef is a slim historical reference returned by the cross-claim
// duplicate scan. The UI renders these inside a tooltip — keep the fields
// short and presentable.
type DupeRef struct {
	ClaimID     int       `json:"claim_id"`
	ClaimNumber string    `json:"claim_number"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	Date        string    `json:"date"`
}

// BAVStatus is the digested banking verification state for a single claim.
// Stale collapses two cases that prompt the same UI response: a verified
// result older than bavStalenessWindow, or a recent banking change since
// verification.
type BAVStatus struct {
	HasResult          bool       `json:"has_result"`
	Status             string     `json:"status"` // complete | failed | pending | ""
	Verified           bool       `json:"verified"`
	VerifiedAt         *time.Time `json:"verified_at"`
	Stale              bool       `json:"stale"`
	StaleReason        string     `json:"stale_reason,omitempty"`
	ProviderRequestID  string     `json:"provider_request_id,omitempty"`
	LastAttempt        int        `json:"last_attempt"`
}

// CheckCrossClaimDuplicates scans for other in-flight or already-paid claims
// sharing this claim's ID number or bank account number. Returns separate
// slices so the UI can show two distinct chips with their own tooltips.
//
// Empty fields are skipped (no point in matching everyone with a blank ID).
// The current claim itself is always excluded from the match. Returned slices
// are truncated to crossClaimDupeRefLimit so the RiskFlags JSON stays compact.
func CheckCrossClaimDuplicates(db *gorm.DB, claim models.GroupSchemeClaim) (idHits, accountHits []DupeRef, err error) {
	if db == nil {
		db = DB
	}

	idNumber := strings.TrimSpace(claim.ClaimantIDNumber)
	if idNumber != "" {
		idHits, err = lookupDupesByField(db, "claimant_id_number", idNumber, claim.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("lookup id duplicates: %w", err)
		}
	}

	account := strings.TrimSpace(claim.BankAccountNumber)
	if account != "" {
		accountHits, err = lookupDupesByField(db, "bank_account_number", account, claim.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("lookup account duplicates: %w", err)
		}
	}

	return idHits, accountHits, nil
}

func lookupDupesByField(db *gorm.DB, column, value string, excludeID int) ([]DupeRef, error) {
	type row struct {
		ID             int
		ClaimNumber    string
		ClaimAmount    float64
		Status         string
		DateRegistered string
	}
	var rows []row
	err := db.Table("group_scheme_claims").
		Select("id, claim_number, claim_amount, status, date_registered").
		Where(column+" = ?", value).
		Where("id <> ?", excludeID).
		Where("status IN ?", inFlightOrPaidStatuses).
		Order("creation_date DESC").
		Limit(crossClaimDupeRefLimit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]DupeRef, 0, len(rows))
	for _, r := range rows {
		out = append(out, DupeRef{
			ClaimID:     r.ID,
			ClaimNumber: r.ClaimNumber,
			Amount:      r.ClaimAmount,
			Status:      r.Status,
			Date:        r.DateRegistered,
		})
	}
	return out, nil
}

// FormatDupeRefs converts DupeRef rows into the short human-readable strings
// the UI shows inside tooltips, e.g. "CLM-12345 · R 5,000 · paid (2026-03-12)".
// Lives here so the schedule-creation path and any future re-render share one
// format.
func FormatDupeRefs(refs []DupeRef) []string {
	out := make([]string, 0, len(refs))
	for _, r := range refs {
		date := strings.TrimSpace(r.Date)
		if date != "" {
			out = append(out, fmt.Sprintf("%s · R %.2f · %s (%s)", r.ClaimNumber, r.Amount, r.Status, date))
		} else {
			out = append(out, fmt.Sprintf("%s · R %.2f · %s", r.ClaimNumber, r.Amount, r.Status))
		}
	}
	return out
}

// LatestApprovedAmountForClaim returns the most-recent ApprovedAmount stamped
// on this claim's assessments, or 0 when no approval has been recorded.
// Schedule creation snapshots this onto the line item so the drift check has
// a stable source of truth even if the assessment is later edited.
func LatestApprovedAmountForClaim(db *gorm.DB, claimID int) float64 {
	if db == nil {
		db = DB
	}
	var assessment models.GroupSchemeClaimAssessment
	err := db.Where("claim_id = ? AND approved_amount IS NOT NULL", claimID).
		Order("approved_at DESC, updated_at DESC").
		Select("id, approved_amount").
		First(&assessment).Error
	if err != nil {
		return 0
	}
	return assessment.ApprovedAmount
}

// ComputeAmountDrift returns the delta between the schedule line's Gross and
// the snapshotted ApprovedAmount, plus a boolean that's true only when the
// absolute delta exceeds amountDriftTolerance. Returns (0, false) when no
// approved amount was snapshotted — legacy lines stay silent rather than
// firing false positives.
func ComputeAmountDrift(item models.ClaimPaymentScheduleItem) (delta float64, drifted bool) {
	if item.ApprovedAmountSnapshot == 0 {
		return 0, false
	}
	delta = item.GrossAmount - item.ApprovedAmountSnapshot
	if math.Abs(delta) <= amountDriftTolerance {
		return 0, false
	}
	return delta, true
}

// LatestBankVerification returns the most-recent BAV outcome for a claim,
// digested into the UI-ready BAVStatus shape. Stale=true when the verified
// result is older than bavStalenessWindow OR when the line carries the
// banking_change_30d risk flag (passed in by the caller — the function
// itself only knows about the BAV log).
func LatestBankVerification(db *gorm.DB, claimID int, bankingChanged bool) (*BAVStatus, error) {
	if db == nil {
		db = DB
	}
	var row models.BAVVerificationLog
	err := db.Where("claim_id = ?", claimID).
		Order("created_at DESC").
		First(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &BAVStatus{HasResult: false}, nil
		}
		return nil, err
	}

	status := &BAVStatus{
		HasResult:         true,
		Status:            row.Status,
		Verified:          strings.EqualFold(row.Status, "complete"),
		VerifiedAt:        &row.CreatedAt,
		ProviderRequestID: row.ProviderRequestID,
	}

	if status.Verified {
		if time.Since(row.CreatedAt) > bavStalenessWindow {
			status.Stale = true
			status.StaleReason = "verified result is older than 30 days — re-verify recommended"
		}
		if bankingChanged {
			status.Stale = true
			if status.StaleReason != "" {
				status.StaleReason += "; "
			}
			status.StaleReason += "banking details changed since this verification — re-verify recommended"
		}
	}

	// Best-effort Attempt extraction so the controller can compute next
	// attempt without re-parsing JSON. Stored under the IdempotencyKey's
	// derivation; the log row itself doesn't carry Attempt as a column, so
	// we approximate by counting prior attempts for this claim.
	var attemptCount int64
	_ = db.Model(&models.BAVVerificationLog{}).Where("claim_id = ?", claimID).Count(&attemptCount).Error
	status.LastAttempt = int(attemptCount)

	return status, nil
}

// LastReverifyWithinWindow returns true when the most recent BAV log row for
// a claim was created within the supplied window. The HTTP re-verify handler
// uses this as a server-side cost guard so accidental double-clicks don't
// hit the provider twice in quick succession.
func LastReverifyWithinWindow(db *gorm.DB, claimID int, window time.Duration) (bool, *time.Time, error) {
	if db == nil {
		db = DB
	}
	var row models.BAVVerificationLog
	err := db.Where("claim_id = ?", claimID).
		Order("created_at DESC").
		First(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil, nil
		}
		return false, nil, err
	}
	if time.Since(row.CreatedAt) < window {
		return true, &row.CreatedAt, nil
	}
	return false, &row.CreatedAt, nil
}

// outstandingAmountDrifts returns the claim numbers on a schedule whose
// snapshotted approved amount differs from the schedule gross by more than
// amountDriftTolerance AND that have not yet been acknowledged. Used by
// FinanceFirstAuthorise to refuse authorisation when drift is pending review.
func outstandingAmountDrifts(scheduleID int) ([]string, error) {
	var items []models.ClaimPaymentScheduleItem
	err := DB.Where("schedule_id = ? AND amount_drift_resolved = ?", scheduleID, false).
		Where("line_status IN ?", []string{"pending", "verified"}).
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(items))
	for _, it := range items {
		if _, drifted := ComputeAmountDrift(it); drifted {
			out = append(out, it.ClaimNumber)
		}
	}
	return out, nil
}

// SnapshotApprovedAmountIntoFlags decorates the RiskFlags struct with cross-
// claim duplicate signals derived from the given hits. Lives next to the
// other risk-flag composition logic so the schedule-creation code can stay
// readable — one call instead of four field assignments per claim.
func DecorateRiskFlagsWithCrossClaim(flags RiskFlags, idHits, accountHits []DupeRef) RiskFlags {
	if len(idHits) > 0 {
		flags.IDPaidBefore = true
		flags.IDPaidBeforeRefs = FormatDupeRefs(idHits)
	}
	if len(accountHits) > 0 {
		flags.AccountUsedBefore = true
		flags.AccountUsedBeforeRefs = FormatDupeRefs(accountHits)
	}
	return flags
}
