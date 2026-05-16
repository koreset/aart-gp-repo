package services

import (
	"api/models"
	"api/services/sanctions"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ──────────────────────────────────────────────
// Sanctions screening orchestration (Phase 3)
// ──────────────────────────────────────────────
//
// Wraps the slim services/sanctions provider interface around the DB-backed
// SanctionsScreening rows so callers (lifecycle, UI) can stay provider-agnostic.

// ScreenLineItem runs the active sanctions provider against a single line and
// upserts a SanctionsScreening row. Returns the persisted row. Called from
// the UI ("Screen now" button) and ahead of first finance authorisation.
//
// "screening row already cleared" is treated as idempotent — re-running over
// a manual_clear row keeps the manual outcome.
func ScreenLineItem(scheduleID, itemID int, user models.AppUser) (models.SanctionsScreening, error) {
	provider := sanctions.Active()
	if provider == nil {
		return models.SanctionsScreening{}, sanctions.ErrProviderNotConfigured
	}

	var item models.ClaimPaymentScheduleItem
	if err := DB.Where("id = ? AND schedule_id = ?", itemID, scheduleID).First(&item).Error; err != nil {
		return models.SanctionsScreening{}, err
	}

	subj := sanctions.Subject{
		FullName:       firstNonEmpty(item.BeneficiaryName, item.MemberName),
		IdentityNumber: firstNonEmpty(item.BeneficiaryIDNumber, item.MemberIDNumber),
		Country:        "ZA",
	}
	result, err := provider.Screen(context.Background(), subj)
	if err != nil {
		return models.SanctionsScreening{}, fmt.Errorf("provider %s failed: %w", provider.Name(), err)
	}
	if result == nil {
		return models.SanctionsScreening{}, errors.New("provider returned no result")
	}

	now := time.Now()

	// Upsert by (schedule_item_id, provider).
	var existing models.SanctionsScreening
	lookup := DB.Where("schedule_item_id = ? AND provider = ?", itemID, provider.Name()).First(&existing).Error
	if lookup == nil {
		// Preserve manual_clear if finance has already overridden a previous hit.
		if existing.Status == string(sanctions.StatusManualClear) {
			return existing, nil
		}
		updates := map[string]interface{}{
			"schedule_id":  scheduleID,
			"claim_id":     item.ClaimID,
			"status":       string(result.Status),
			"provider_ref": result.ProviderRef,
			"hit_summary":  result.HitSummary,
			"screened_by":  user.UserName,
			"screened_at":  &now,
		}
		if err := DB.Model(&existing).Updates(updates).Error; err != nil {
			return existing, err
		}
		return GetSanctionsScreeningByID(existing.ID)
	}
	if !errors.Is(lookup, gorm.ErrRecordNotFound) {
		return models.SanctionsScreening{}, lookup
	}

	row := models.SanctionsScreening{
		ScheduleID:     scheduleID,
		ScheduleItemID: itemID,
		ClaimID:        item.ClaimID,
		Provider:       provider.Name(),
		Status:         string(result.Status),
		ProviderRef:    result.ProviderRef,
		HitSummary:     result.HitSummary,
		ScreenedBy:     user.UserName,
		ScreenedAt:     &now,
	}
	if err := DB.Create(&row).Error; err != nil {
		return row, err
	}
	return row, nil
}

// GetSanctionsScreeningByID fetches a single screening row.
func GetSanctionsScreeningByID(id int) (models.SanctionsScreening, error) {
	var row models.SanctionsScreening
	err := DB.First(&row, id).Error
	return row, err
}

// ListSanctionsScreeningsForSchedule returns every screening row attached to
// a schedule, newest first. Drives the schedule's sanctions overview.
func ListSanctionsScreeningsForSchedule(scheduleID int) ([]models.SanctionsScreening, error) {
	var rows []models.SanctionsScreening
	err := DB.Where("schedule_id = ?", scheduleID).Order("updated_at DESC").Find(&rows).Error
	return rows, err
}

// RecordSanctionsOutcomeRequest is the inbound payload for a manual outcome.
type RecordSanctionsOutcomeRequest struct {
	Status string `json:"status"` // clear | hit | manual_clear
	Notes  string `json:"notes"`
}

// RecordSanctionsOutcome lets finance manually mark a line as clear, hit, or
// manual_clear (the latter being the "I've reviewed the hit and confirmed
// it's a false positive" path). Always upserts; if no provider row exists
// yet, one is created with provider="manual".
func RecordSanctionsOutcome(scheduleID, itemID int, req RecordSanctionsOutcomeRequest, user models.AppUser) (models.SanctionsScreening, error) {
	allowed := map[string]bool{
		string(sanctions.StatusClear):       true,
		string(sanctions.StatusHit):         true,
		string(sanctions.StatusManualClear): true,
	}
	if !allowed[req.Status] {
		return models.SanctionsScreening{}, fmt.Errorf("invalid sanctions status %q", req.Status)
	}

	var item models.ClaimPaymentScheduleItem
	if err := DB.Where("id = ? AND schedule_id = ?", itemID, scheduleID).First(&item).Error; err != nil {
		return models.SanctionsScreening{}, err
	}

	providerName := "manual"
	if p := sanctions.Active(); p != nil {
		providerName = p.Name()
	}

	now := time.Now()
	var existing models.SanctionsScreening
	lookup := DB.Where("schedule_item_id = ? AND provider = ?", itemID, providerName).First(&existing).Error
	if errors.Is(lookup, gorm.ErrRecordNotFound) {
		row := models.SanctionsScreening{
			ScheduleID:     scheduleID,
			ScheduleItemID: itemID,
			ClaimID:        item.ClaimID,
			Provider:       providerName,
			Status:         req.Status,
			Notes:          req.Notes,
			ScreenedBy:     user.UserName,
			ScreenedAt:     &now,
		}
		if req.Status == string(sanctions.StatusManualClear) || req.Status == string(sanctions.StatusClear) {
			row.ClearedBy = user.UserName
			row.ClearedAt = &now
		}
		if err := DB.Create(&row).Error; err != nil {
			return row, err
		}
		return row, nil
	}
	if lookup != nil {
		return models.SanctionsScreening{}, lookup
	}

	updates := map[string]interface{}{
		"status": req.Status,
		"notes":  req.Notes,
	}
	if req.Status == string(sanctions.StatusManualClear) || req.Status == string(sanctions.StatusClear) {
		updates["cleared_by"] = user.UserName
		updates["cleared_at"] = &now
	}
	if err := DB.Model(&existing).Updates(updates).Error; err != nil {
		return existing, err
	}
	return GetSanctionsScreeningByID(existing.ID)
}

// blockingSanctionsItems returns the line-item IDs whose latest screening
// status is "hit" or whose only screening row is still "pending". Used by
// FinanceFirstAuthorise to refuse authorisation until all blockers are
// resolved.
func blockingSanctionsItems(scheduleID int) ([]string, error) {
	type row struct {
		ID          int
		ClaimNumber string
		Status      string
	}

	// Find items whose newest screening row is a blocker.
	var rows []row
	query := `
		SELECT i.id AS id, i.claim_number AS claim_number, s.status AS status
		FROM claim_payment_schedule_items i
		LEFT JOIN sanctions_screenings s ON s.schedule_item_id = i.id
		WHERE i.schedule_id = ?
		AND i.line_status IN ('pending','verified')
	`
	if err := DB.Raw(query, scheduleID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	// One item may have multiple provider rows; pick the worst across all.
	worst := map[int]string{}
	claim := map[int]string{}
	for _, r := range rows {
		claim[r.ID] = r.ClaimNumber
		curr := worst[r.ID]
		worst[r.ID] = pickWorseStatus(curr, r.Status)
	}

	var problems []string
	for itemID, status := range worst {
		if isSanctionsBlocker(status) {
			problems = append(problems, fmt.Sprintf("%s (%s)", claim[itemID], statusLabel(status)))
		}
	}
	return problems, nil
}

// pickWorseStatus returns whichever of a and b should block authorisation.
// Order, worst-first: "hit" > "pending" > "" > "skipped" > "clear" > "manual_clear".
func pickWorseStatus(a, b string) string {
	rank := func(s string) int {
		switch strings.ToLower(s) {
		case "hit":
			return 5
		case "pending":
			return 4
		case "":
			return 3
		case "skipped":
			return 2
		case "clear":
			return 1
		case "manual_clear":
			return 0
		}
		return 3
	}
	if rank(a) >= rank(b) {
		return a
	}
	return b
}

func isSanctionsBlocker(status string) bool {
	switch strings.ToLower(status) {
	case "hit", "pending", "skipped", "":
		return true
	}
	return false
}

func statusLabel(s string) string {
	if s == "" {
		return "not screened"
	}
	return s
}

func firstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}
