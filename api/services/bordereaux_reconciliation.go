package services

import (
	appLog "api/log"
	"api/models"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// escalationSLA returns the SLA duration for a given priority. Unknown values
// fall back to the medium default so every escalation has a due date.
func escalationSLA(priority string) time.Duration {
	switch strings.ToLower(strings.TrimSpace(priority)) {
	case "high":
		return 24 * time.Hour
	case "low":
		return 7 * 24 * time.Hour
	default: // medium or unspecified
		return 72 * time.Hour
	}
}

// ResolveDiscrepancyRequest is the payload for resolving a reconciliation discrepancy.
type ResolveDiscrepancyRequest struct {
	Resolution    string `json:"resolution" binding:"required"` // accept_expected | accept_actual | manual_override
	OverrideValue string `json:"override_value"`
	Notes         string `json:"notes"`
}

// EscalateDiscrepancyRequest is the payload for escalating a discrepancy.
type EscalateDiscrepancyRequest struct {
	EscalateTo string `json:"escalate_to" binding:"required"`
	Reason     string `json:"reason" binding:"required"`
	Priority   string `json:"priority"` // low | medium | high
}

// AddDiscrepancyCommentRequest is the payload for adding a comment to a discrepancy.
type AddDiscrepancyCommentRequest struct {
	Comment string `json:"comment" binding:"required"`
}

// ResolveDiscrepancy marks a single reconciliation result as resolved.
func ResolveDiscrepancy(id int, req ResolveDiscrepancyRequest, user models.AppUser) (models.BordereauxReconciliationResult, error) {
	var result models.BordereauxReconciliationResult
	if err := DB.First(&result, id).Error; err != nil {
		return result, fmt.Errorf("reconciliation result %d not found", id)
	}

	before := result

	result.IsResolved = true
	result.Status = "resolved"
	if req.Notes != "" {
		if result.Comments != "" {
			result.Comments += "\n"
		}
		result.Comments += fmt.Sprintf("[%s] %s resolved (%s): %s",
			time.Now().Format("2006-01-02 15:04"), user.UserName, req.Resolution, req.Notes)
	} else {
		if result.Comments != "" {
			result.Comments += "\n"
		}
		result.Comments += fmt.Sprintf("[%s] %s resolved (%s)",
			time.Now().Format("2006-01-02 15:04"), user.UserName, req.Resolution)
	}

	if err := DB.Save(&result).Error; err != nil {
		return result, fmt.Errorf("failed to resolve discrepancy: %w", err)
	}

	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "bordereaux_reconciliation_results",
		EntityID:  strconv.Itoa(id),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, result)

	// Update parent confirmation stats
	updateConfirmationStats(result.BordereauxConfirmationID)

	return result, nil
}

// EscalateDiscrepancy flags a discrepancy for escalation, assigns it to a
// target user/role, computes an SLA due date from the priority, and fires a
// notification to the assignee. Notification delivery failures are logged but
// do not block the escalation write — the primary action must succeed even
// if routing is misconfigured.
func EscalateDiscrepancy(id int, req EscalateDiscrepancyRequest, user models.AppUser) (models.BordereauxReconciliationResult, error) {
	var result models.BordereauxReconciliationResult
	if err := DB.First(&result, id).Error; err != nil {
		return result, fmt.Errorf("reconciliation result %d not found", id)
	}

	before := result

	now := time.Now()
	due := now.Add(escalationSLA(req.Priority))
	assigneeEmail := ResolveEscalationTarget(req.EscalateTo)

	result.Status = "escalated"
	result.AssignedTo = req.EscalateTo
	result.Priority = strings.ToLower(strings.TrimSpace(req.Priority))
	result.EscalatedBy = user.UserName
	result.EscalatedAt = &now
	result.DueDate = &due

	note := fmt.Sprintf("[%s] Escalated to %s by %s (priority: %s, due: %s): %s",
		now.Format("2006-01-02 15:04"), req.EscalateTo, user.UserName,
		result.Priority, due.Format("2006-01-02 15:04"), req.Reason)
	if result.Comments != "" {
		result.Comments += "\n"
	}
	result.Comments += note

	if err := DB.Save(&result).Error; err != nil {
		return result, fmt.Errorf("failed to escalate discrepancy: %w", err)
	}

	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "bordereaux_reconciliation_results",
		EntityID:  strconv.Itoa(id),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, result)

	go func() {
		if err := NotifyDiscrepancyEscalated(result, assigneeEmail, user, req.Reason); err != nil {
			appLog.WithFields(map[string]interface{}{
				"result_id":      result.ID,
				"assignee_email": assigneeEmail,
				"error":          err.Error(),
			}).Error("Escalation notification delivery failed")
		}
	}()

	return result, nil
}

// ListEscalationsRequest is the payload for filtering the escalation queue.
// All fields are optional; omitted filters are ignored.
type ListEscalationsRequest struct {
	AssignedTo  string
	Priority    string
	OverdueOnly bool
}

// ListEscalations returns all reconciliation results currently in the
// escalated state, optionally filtered by assignee, priority, or overdue-ness.
// Overdue is computed relative to time.Now() against the stored DueDate.
func ListEscalations(filter ListEscalationsRequest) ([]models.BordereauxReconciliationResult, error) {
	q := DB.Where("status = ?", "escalated")
	if filter.AssignedTo != "" {
		q = q.Where("assigned_to = ?", filter.AssignedTo)
	}
	if filter.Priority != "" {
		q = q.Where("priority = ?", strings.ToLower(strings.TrimSpace(filter.Priority)))
	}
	if filter.OverdueOnly {
		q = q.Where("due_date IS NOT NULL AND due_date < ?", time.Now())
	}
	var results []models.BordereauxReconciliationResult
	if err := q.Order("due_date ASC").Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// AddDiscrepancyComment appends a timestamped comment to a reconciliation result.
func AddDiscrepancyComment(id int, req AddDiscrepancyCommentRequest, user models.AppUser) (models.BordereauxReconciliationResult, error) {
	var result models.BordereauxReconciliationResult
	if err := DB.First(&result, id).Error; err != nil {
		return result, fmt.Errorf("reconciliation result %d not found", id)
	}

	before := result

	note := fmt.Sprintf("[%s] %s: %s", time.Now().Format("2006-01-02 15:04"), user.UserName, req.Comment)
	if result.Comments != "" {
		result.Comments += "\n"
	}
	result.Comments += note

	if err := DB.Save(&result).Error; err != nil {
		return result, fmt.Errorf("failed to add comment: %w", err)
	}

	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "bordereaux_reconciliation_results",
		EntityID:  strconv.Itoa(id),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, result)

	return result, nil
}

// ConfirmReconciliation marks all results for a confirmation as resolved (bulk accept).
func ConfirmReconciliation(confirmationID int, user models.AppUser) error {
	var confirmation models.BordereauxConfirmation
	if err := DB.First(&confirmation, confirmationID).Error; err != nil {
		return fmt.Errorf("confirmation %d not found", confirmationID)
	}

	note := fmt.Sprintf("[%s] Bulk confirmed by %s", time.Now().Format("2006-01-02 15:04"), user.UserName)

	result := DB.Model(&models.BordereauxReconciliationResult{}).
		Where("bordereaux_confirmation_id = ? AND is_resolved = ?", confirmationID, false).
		Updates(map[string]interface{}{
			"is_resolved": true,
			"status":      "resolved",
			"comments":    gorm.Expr("CONCAT(COALESCE(comments,''), ?)", "\n"+note),
		})
	if result.Error != nil {
		return fmt.Errorf("failed to confirm reconciliation: %w", result.Error)
	}

	// Update confirmation status
	confirmation.Status = "matched"
	now := time.Now()
	confirmation.LastReconciled = &now
	DB.Save(&confirmation)

	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "bordereaux_confirmations",
		EntityID:  strconv.Itoa(confirmationID),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, struct{}{}, map[string]interface{}{
		"action":          "bulk_confirm",
		"confirmation_id": confirmationID,
		"rows_resolved":   result.RowsAffected,
	})

	return nil
}

// ReprocessReconciliation re-runs the reconciliation engine for a given confirmation.
func ReprocessReconciliation(confirmationID int, user models.AppUser) (*ReconciliationSummary, error) {
	var confirmation models.BordereauxConfirmation
	if err := DB.First(&confirmation, confirmationID).Error; err != nil {
		return nil, fmt.Errorf("confirmation %d not found", confirmationID)
	}

	// Clear existing results for a clean re-run
	if err := DB.Where("bordereaux_confirmation_id = ?", confirmationID).
		Delete(&models.BordereauxReconciliationResult{}).Error; err != nil {
		return nil, fmt.Errorf("failed to clear existing results: %w", err)
	}

	// Reset confirmation status
	confirmation.Status = "pending"
	confirmation.MatchedCount = 0
	confirmation.DiscrepancyCount = 0
	DB.Save(&confirmation)

	// Re-run reconciliation atomically so cleared results + new results + confirmation
	// status all commit together.
	var summary *ReconciliationSummary
	txErr := DB.Transaction(func(tx *gorm.DB) error {
		s, err := ReconcileConfirmation(context.Background(), tx, &confirmation)
		if err != nil {
			return err
		}
		summary = s
		return nil
	})
	if txErr != nil {
		return nil, fmt.Errorf("reprocessing failed: %w", txErr)
	}

	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "bordereaux_confirmations",
		EntityID:  strconv.Itoa(confirmationID),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, struct{}{}, map[string]interface{}{
		"action":  "reprocess",
		"summary": summary,
	})

	return summary, nil
}

// AddReconciliationNote persists a free-text note against a confirmation in
// the dedicated bordereaux_confirmation_notes table. Legacy _note synthetic
// rows in bordereaux_reconciliation_results are no longer written; the filter
// on them in updateConfirmationStats is kept as a safety net for historic data.
func AddReconciliationNote(confirmationID int, note string, user models.AppUser) (models.BordereauxConfirmation, error) {
	var confirmation models.BordereauxConfirmation
	if err := DB.First(&confirmation, confirmationID).Error; err != nil {
		return confirmation, fmt.Errorf("confirmation %d not found", confirmationID)
	}
	entry := models.BordereauxConfirmationNote{
		BordereauxConfirmationID: confirmationID,
		Note:                     note,
		CreatedBy:                user.UserName,
	}
	if err := DB.Create(&entry).Error; err != nil {
		return confirmation, fmt.Errorf("failed to add note: %w", err)
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "bordereaux_confirmation_notes",
		EntityID:  strconv.Itoa(entry.ID),
		Action:    "CREATE",
		ChangedBy: user.UserName,
	}, struct{}{}, entry)
	return confirmation, nil
}

// GetReconciliationNotes returns the notes for a confirmation, newest first.
func GetReconciliationNotes(confirmationID int) ([]models.BordereauxConfirmationNote, error) {
	var notes []models.BordereauxConfirmationNote
	if err := DB.Where("bordereaux_confirmation_id = ?", confirmationID).
		Order("created_at DESC").
		Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

// updateConfirmationStats recalculates matched/discrepancy counts for a confirmation.
func updateConfirmationStats(confirmationID int) {
	var matched, discrepancy int64
	DB.Model(&models.BordereauxReconciliationResult{}).
		Where("bordereaux_confirmation_id = ? AND status NOT IN ('note') AND is_resolved = ?", confirmationID, true).
		Count(&matched)
	DB.Model(&models.BordereauxReconciliationResult{}).
		Where("bordereaux_confirmation_id = ? AND status NOT IN ('note','resolved','matched') AND is_resolved = ?", confirmationID, false).
		Count(&discrepancy)

	updates := map[string]interface{}{
		"matched_count":     int(matched),
		"discrepancy_count": int(discrepancy),
	}
	if discrepancy == 0 {
		updates["status"] = "matched"
		now := time.Now()
		DB.Model(&models.BordereauxConfirmation{}).Where("id = ?", confirmationID).Updates(map[string]interface{}{
			"matched_count":     int(matched),
			"discrepancy_count": 0,
			"status":            "matched",
			"last_reconciled":   &now,
		})
	} else {
		DB.Model(&models.BordereauxConfirmation{}).Where("id = ?", confirmationID).Updates(updates)
	}
}
