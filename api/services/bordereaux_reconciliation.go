package services

import (
	"api/models"
	"context"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

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

// EscalateDiscrepancy flags a discrepancy for escalation.
func EscalateDiscrepancy(id int, req EscalateDiscrepancyRequest, user models.AppUser) (models.BordereauxReconciliationResult, error) {
	var result models.BordereauxReconciliationResult
	if err := DB.First(&result, id).Error; err != nil {
		return result, fmt.Errorf("reconciliation result %d not found", id)
	}

	before := result

	result.Status = "escalated"
	note := fmt.Sprintf("[%s] Escalated to %s by %s (priority: %s): %s",
		time.Now().Format("2006-01-02 15:04"), req.EscalateTo, user.UserName, req.Priority, req.Reason)
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

	return result, nil
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

	// Re-run reconciliation
	summary, err := ReconcileConfirmation(context.Background(), &confirmation)
	if err != nil {
		return nil, fmt.Errorf("reprocessing failed: %w", err)
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

// AddReconciliationNote appends a note to a confirmation's status field.
func AddReconciliationNote(confirmationID int, note string, user models.AppUser) (models.BordereauxConfirmation, error) {
	var confirmation models.BordereauxConfirmation
	if err := DB.First(&confirmation, confirmationID).Error; err != nil {
		return confirmation, fmt.Errorf("confirmation %d not found", confirmationID)
	}

	// Store notes in the status field as a pragmatic approach since there's no dedicated notes column
	// The note is appended with timestamp and user
	entry := fmt.Sprintf("[%s] %s: %s", time.Now().Format("2006-01-02 15:04"), user.UserName, note)

	// We'll use the reconciliation result comments on a synthetic "note" record
	noteResult := models.BordereauxReconciliationResult{
		BordereauxConfirmationID: confirmationID,
		GeneratedBordereauxID:    confirmation.GeneratedBordereauxID,
		Field:                    "_note",
		Status:                   "note",
		Comments:                 entry,
		IsResolved:               true,
	}
	if err := DB.Create(&noteResult).Error; err != nil {
		return confirmation, fmt.Errorf("failed to add note: %w", err)
	}

	return confirmation, nil
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
