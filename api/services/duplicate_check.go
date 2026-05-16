package services

import (
	"api/models"
	"errors"
	"fmt"
	"strings"
)

// ──────────────────────────────────────────────
// Duplicate detection (Phase 3)
// ──────────────────────────────────────────────
//
// Two duplicate classes are checked:
//
//   1. Same claim on a prior schedule
//        — hard block. Same claim cannot be paid twice. The new schedule
//          creation is refused outright with the offending schedule references.
//          A line in an archived or cancelled schedule does NOT block.
//
//   2. Same beneficiary on multiple lines within the same schedule
//        — soft block. Flagged on each affected line via
//          DuplicateBeneficiaryFlag; finance must explicitly clear the flag
//          before first authorisation. A single beneficiary getting two
//          legitimate payouts (e.g. funeral + GLA) is valid; the flag is
//          there so finance can confirm intent, not block automatically.

// ErrDuplicateClaim is returned by CheckClaimsNotOnPriorSchedule when at least
// one of the supplied claim IDs is already attached to a non-archived
// schedule.
var ErrDuplicateClaim = errors.New("one or more claims are already on an active payment schedule")

// CheckClaimsNotOnPriorSchedule rejects creation when any of the supplied
// claim IDs is already on a non-archived schedule. Called from
// CreatePaymentSchedule.
func CheckClaimsNotOnPriorSchedule(claimIDs []int) error {
	if len(claimIDs) == 0 {
		return nil
	}
	type hit struct {
		ClaimID        int
		ClaimNumber    string
		ScheduleNumber string
		Status         string
	}
	var hits []hit
	err := DB.Raw(`
		SELECT i.claim_id AS claim_id, i.claim_number AS claim_number,
		       s.schedule_number AS schedule_number, s.status AS status
		FROM claim_payment_schedule_items i
		JOIN claim_payment_schedules s ON s.id = i.schedule_id
		WHERE i.claim_id IN ?
		AND LOWER(s.status) NOT IN ('archived','cancelled')
		AND i.line_status IN ('pending','verified')
	`, claimIDs).Scan(&hits).Error
	if err != nil {
		return err
	}
	if len(hits) == 0 {
		return nil
	}
	parts := make([]string, 0, len(hits))
	for _, h := range hits {
		parts = append(parts, fmt.Sprintf("%s (on %s, %s)", h.ClaimNumber, h.ScheduleNumber, h.Status))
	}
	return fmt.Errorf("%w: %s", ErrDuplicateClaim, strings.Join(parts, "; "))
}

// FlagDuplicateBeneficiaries walks the supplied items and marks any whose
// beneficiary key (ID number, falling back to name) appears on another line
// in the same set. The flag is written back to the database for each
// affected row. Returns the count of items flagged so callers can include
// it in the genesis audit row.
func FlagDuplicateBeneficiaries(scheduleID int, items []models.ClaimPaymentScheduleItem) (int, error) {
	keyCounts := map[string][]int{}
	for _, it := range items {
		key := beneficiaryKey(it)
		if key == "" {
			continue
		}
		keyCounts[key] = append(keyCounts[key], it.ID)
	}
	var flagged int
	for _, ids := range keyCounts {
		if len(ids) < 2 {
			continue
		}
		if err := DB.Model(&models.ClaimPaymentScheduleItem{}).
			Where("id IN ?", ids).
			Update("duplicate_beneficiary_flag", true).Error; err != nil {
			return flagged, err
		}
		flagged += len(ids)
	}
	return flagged, nil
}

// FlagDuplicateBeneficiariesForSchedule re-runs the within-schedule check
// against the schedule's current item set. Used when entering finance review
// to catch beneficiary collisions that may not have been visible at create
// time (e.g. after a previous query/reject reshaped the schedule).
func FlagDuplicateBeneficiariesForSchedule(scheduleID int) (int, error) {
	var items []models.ClaimPaymentScheduleItem
	if err := DB.Where("schedule_id = ? AND line_status IN ?", scheduleID, []string{"pending", "verified"}).Find(&items).Error; err != nil {
		return 0, err
	}
	return FlagDuplicateBeneficiaries(scheduleID, items)
}

// ClearDuplicateBeneficiary marks a single line's within-schedule duplicate
// flag as reviewed and cleared. Does not remove the underlying flag (so the
// duplicate-warning UI still surfaces context); only `cleared` is toggled.
func ClearDuplicateBeneficiary(scheduleID, itemID int, user models.AppUser) error {
	return DB.Model(&models.ClaimPaymentScheduleItem{}).
		Where("id = ? AND schedule_id = ?", itemID, scheduleID).
		Updates(map[string]interface{}{
			"duplicate_beneficiary_cleared": true,
		}).Error
}

// outstandingDuplicateBeneficiaries returns the claim numbers whose duplicate
// flag is set but not yet cleared. Used by FinanceFirstAuthorise to refuse
// authorisation until finance has confirmed the duplicates are intentional.
func outstandingDuplicateBeneficiaries(scheduleID int) ([]string, error) {
	var rows []struct{ ClaimNumber string }
	err := DB.Model(&models.ClaimPaymentScheduleItem{}).
		Select("claim_number").
		Where("schedule_id = ? AND duplicate_beneficiary_flag = ? AND duplicate_beneficiary_cleared = ?", scheduleID, true, false).
		Where("line_status IN ?", []string{"pending", "verified"}).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.ClaimNumber)
	}
	return out, nil
}

// beneficiaryKey is the deduplication key. Identity number wins (uniquely
// identifies a person across name spellings); name is a fallback when ID
// is missing.
func beneficiaryKey(it models.ClaimPaymentScheduleItem) string {
	id := strings.TrimSpace(it.BeneficiaryIDNumber)
	if id != "" {
		return "id:" + strings.ToLower(id)
	}
	name := strings.TrimSpace(it.BeneficiaryName)
	if name != "" {
		return "name:" + strings.ToLower(name)
	}
	return ""
}
