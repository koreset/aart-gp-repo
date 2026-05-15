package services

import (
	"api/models"
	"time"

	"gorm.io/gorm"
)

// RecordQuoteStatusChange appends one row to the
// group_pricing_quote_status_audits table for a quote whose status just
// transitioned from oldStatus to newStatus. It computes the elapsed
// seconds since the prior audit row for the same quote (or 0 if this is
// the first real transition) and stores it on the row so SLA-breach
// queries don't need a self-join over the audit table.
//
// Call this inside the same transaction as the quote update so a rolled-
// back save also rolls back the audit row. The synthetic flag is left
// false; only the backfill migration writes synthetic = true rows.
func RecordQuoteStatusChange(tx *gorm.DB, quoteID int, oldStatus, newStatus models.Status, message, changedBy string) error {
	var prev models.GroupPricingQuoteStatusAudit
	var durationSecs int64
	if err := tx.Where("quote_id = ?", quoteID).Order("changed_at DESC").First(&prev).Error; err == nil {
		durationSecs = int64(time.Since(prev.ChangedAt).Seconds())
		if durationSecs < 0 {
			durationSecs = 0
		}
	}
	return tx.Create(&models.GroupPricingQuoteStatusAudit{
		QuoteID:              quoteID,
		OldStatus:            oldStatus,
		NewStatus:            newStatus,
		StatusMessage:        message,
		ChangedBy:            changedBy,
		ChangedAt:            time.Now(),
		DurationFromPrevSecs: durationSecs,
		Synthetic:            false,
	}).Error
}

// applyQuoteStatusTimestamp sets the per-status milestone timestamp
// matching newStatus on the supplied quote pointer. No-op for statuses
// that don't have a dedicated milestone column (e.g. draft, in_progress).
// Returns the time it stamped so callers can pass the same now value
// into the audit helper for a consistent audit trail.
func applyQuoteStatusTimestamp(quote *models.GroupPricingQuote, newStatus models.Status) time.Time {
	now := time.Now()
	switch newStatus {
	case models.StatusSubmitted, models.StatusPendingReview:
		quote.SubmittedAt = &now
	case models.StatusApproved:
		quote.ApprovedAt = &now
	case models.StatusRejected:
		quote.RejectedAt = &now
	case models.StatusAccepted:
		quote.AcceptedAt = &now
	case models.StatusInForce:
		quote.InForceAt = &now
	}
	return now
}

// GetGroupPricingQuoteStatusAudit returns the full status-change history
// for a single quote, most-recent first. Used by the dashboard's status-
// history drill-down.
func GetGroupPricingQuoteStatusAudit(quoteId string) ([]models.GroupPricingQuoteStatusAudit, error) {
	var audits []models.GroupPricingQuoteStatusAudit
	if err := DB.Where("quote_id = ?", quoteId).Order("changed_at DESC").Find(&audits).Error; err != nil {
		return nil, err
	}
	return audits, nil
}
