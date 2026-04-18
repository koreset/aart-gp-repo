package services

import (
	appLog "api/log"
	"api/models"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// CreateClaimNotification creates a new notification log entry.
// Validates that the claim_number exists on the specified scheme.
func CreateClaimNotification(req models.CreateClaimNotificationRequest, user models.AppUser) (models.ClaimNotificationLog, error) {
	// Validate claim exists
	if req.SchemeID > 0 && req.ClaimNumber != "" {
		var count int64
		DB.Model(&models.GroupSchemeClaim{}).
			Where("scheme_id = ? AND claim_number = ?", req.SchemeID, req.ClaimNumber).
			Count(&count)
		if count == 0 {
			return models.ClaimNotificationLog{}, fmt.Errorf("claim %s not found on scheme %d", req.ClaimNumber, req.SchemeID)
		}
	} else if req.ClaimNumber != "" {
		var count int64
		DB.Model(&models.GroupSchemeClaim{}).
			Where("claim_number = ?", req.ClaimNumber).
			Count(&count)
		if count == 0 {
			return models.ClaimNotificationLog{}, fmt.Errorf("claim %s not found", req.ClaimNumber)
		}
	}

	n := models.ClaimNotificationLog{
		ClaimID:          req.ClaimID,
		ClaimNumber:      req.ClaimNumber,
		SchemeID:         req.SchemeID,
		SchemeName:       req.SchemeName,
		ReinsurerName:    req.ReinsurerName,
		ReinsurerCode:    req.ReinsurerCode,
		NotificationType: req.NotificationType,
		Status:           "pending",
		DueDate:          req.DueDate,
		Notes:            req.Notes,
		CreatedBy:        user.UserName,
	}
	if err := DB.Create(&n).Error; err != nil {
		return n, err
	}
	return n, nil
}

// GetClaimNotifications returns a paginated, filtered list of notification logs.
// Pending and sent records whose due date has passed are automatically updated to overdue.
func GetClaimNotifications(claimID, schemeID, page, pageSize int, reinsurerCode, notificationType, status, claimNumber string) ([]models.ClaimNotificationLog, int64, error) {
	var records []models.ClaimNotificationLog
	var total int64
	ctx := context.Background()

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 500 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		q := d.Model(&models.ClaimNotificationLog{})
		if claimID > 0 {
			q = q.Where("claim_id = ?", claimID)
		}
		if schemeID > 0 {
			q = q.Where("scheme_id = ?", schemeID)
		}
		if reinsurerCode != "" {
			q = q.Where("reinsurer_code = ?", reinsurerCode)
		}
		if notificationType != "" {
			q = q.Where("notification_type = ?", notificationType)
		}
		if status != "" {
			q = q.Where("status = ?", status)
		}
		if claimNumber != "" {
			q = q.Where("claim_number LIKE ?", "%"+claimNumber+"%")
		}
		if err := q.Count(&total).Error; err != nil {
			return err
		}
		return q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error
	})
	if err != nil {
		return nil, 0, err
	}

	// Auto-mark overdue: pending or sent records whose due date has passed
	now := time.Now()
	for i, r := range records {
		if (r.Status == "pending" || r.Status == "sent") && r.DueDate != "" {
			if dueDate, err := time.Parse("2006-01-02", r.DueDate); err == nil {
				if now.After(dueDate) {
					records[i].Status = "overdue"
					DB.Model(&models.ClaimNotificationLog{}).Where("id = ?", r.ID).Update("status", "overdue")
				}
			}
		}
	}
	return records, total, nil
}

// MarkNotificationSent marks a notification as sent and records the timestamp.
func MarkNotificationSent(id int, notes string, user models.AppUser) (models.ClaimNotificationLog, error) {
	var n models.ClaimNotificationLog
	if err := DB.First(&n, id).Error; err != nil {
		return n, fmt.Errorf("notification not found")
	}
	if err := ValidateClaimNotificationLogTransition(n.Status, StatusClaimNotifSent); err != nil {
		return n, err
	}
	before := n
	now := time.Now()
	n.Status = StatusClaimNotifSent
	n.SentAt = &now
	if notes != "" {
		if n.Notes != "" {
			n.Notes += "\n" + notes
		} else {
			n.Notes = notes
		}
	}
	if err := DB.Save(&n).Error; err != nil {
		return n, err
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "claim_notification_logs",
		EntityID:  strconv.Itoa(id),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, n)
	return n, nil
}

// MarkNotificationAcknowledged marks a notification as acknowledged by the reinsurer.
func MarkNotificationAcknowledged(id int, notes string, user models.AppUser) (models.ClaimNotificationLog, error) {
	var n models.ClaimNotificationLog
	if err := DB.First(&n, id).Error; err != nil {
		return n, fmt.Errorf("notification not found")
	}
	if err := ValidateClaimNotificationLogTransition(n.Status, StatusClaimNotifAcknowledged); err != nil {
		return n, err
	}
	before := n
	now := time.Now()
	n.Status = StatusClaimNotifAcknowledged
	n.AcknowledgedAt = &now
	if notes != "" {
		if n.Notes != "" {
			n.Notes += "\n" + notes
		} else {
			n.Notes = notes
		}
	}
	if err := DB.Save(&n).Error; err != nil {
		return n, err
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "claim_notification_logs",
		EntityID:  strconv.Itoa(id),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, n)
	return n, nil
}

// DeleteClaimNotification deletes a notification that is still in pending status.
func DeleteClaimNotification(id int) error {
	var n models.ClaimNotificationLog
	if err := DB.First(&n, id).Error; err != nil {
		return fmt.Errorf("notification not found")
	}
	if n.Status != "pending" {
		return fmt.Errorf("only pending notifications can be deleted (current status: %s)", n.Status)
	}
	return DB.Delete(&n).Error
}

// GenerateMonthEndNotifications auto-generates notifications for all open claims on a scheme.
// Looks up the scheme's reinsurance treaty to populate reinsurer info and SLA-based due dates.
//
// Rules per claim:
//   - Registered this period → "initial"
//   - Status paid/declined/rejected/closed, no prior "final" → "final"
//   - All other open statuses → "status_update" (deduplicated per period)
func GenerateMonthEndNotifications(schemeID, month, year int, user models.AppUser) ([]models.ClaimNotificationLog, error) {
	if month == 0 {
		month = int(time.Now().Month())
	}
	if year == 0 {
		year = time.Now().Year()
	}

	var claims []models.GroupSchemeClaim
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("scheme_id = ?", schemeID).Find(&claims).Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch claims: %w", err)
	}

	// Look up active treaty linked to this scheme via TreatySchemeLink junction
	var treaty models.ReinsuranceTreaty
	var hasTreaty bool
	treaties, treatyErr := GetActiveTreatiesForScheme(schemeID)
	if treatyErr == nil && len(treaties) > 0 {
		treaty = treaties[0]
		hasTreaty = true
	}

	periodStart := fmt.Sprintf("%04d-%02d-01", year, month)
	nextMonth := month + 1
	nextYear := year
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}
	periodEnd := fmt.Sprintf("%04d-%02d-01", nextYear, nextMonth)

	// Compute due date from treaty SLA or fall back to end of month
	notificationDays := 30
	if hasTreaty && treaty.ClaimsNotificationDays > 0 {
		notificationDays = treaty.ClaimsNotificationDays
	}
	dueDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).
		AddDate(0, 0, notificationDays-1).Format("2006-01-02")

	var created []models.ClaimNotificationLog
	for _, claim := range claims {
		var notifType string

		switch {
		case claim.DateRegistered >= periodStart && claim.DateRegistered < periodEnd:
			notifType = "initial"

		case claim.Status == "paid" || claim.Status == "declined" || claim.Status == "rejected" || claim.Status == "closed":
			var existingFinal int64
			DB.Model(&models.ClaimNotificationLog{}).
				Where("claim_number = ? AND notification_type = 'final'", claim.ClaimNumber).
				Count(&existingFinal)
			if existingFinal > 0 {
				continue
			}
			notifType = "final"

		case claim.Status == "pending" || claim.Status == "registered" || claim.Status == "open" ||
			claim.Status == "approved" || claim.Status == "in_progress" || claim.Status == "under_assessment":
			// Deduplicate: skip if a status_update already exists for this period
			var existing int64
			DB.Model(&models.ClaimNotificationLog{}).
				Where("claim_number = ? AND notification_type = 'status_update' AND due_date >= ? AND due_date < ?",
					claim.ClaimNumber, periodStart, periodEnd).
				Count(&existing)
			if existing > 0 {
				continue
			}
			notifType = "status_update"

		default:
			continue
		}

		n := models.ClaimNotificationLog{
			ClaimID:          claim.ID,
			ClaimNumber:      claim.ClaimNumber,
			SchemeID:         schemeID,
			SchemeName:       claim.SchemeName,
			NotificationType: notifType,
			Status:           "pending",
			DueDate:          dueDate,
			CreatedBy:        user.UserName,
		}
		// Populate reinsurer info from treaty if available
		if hasTreaty {
			n.ReinsurerName = treaty.ReinsurerName
			n.ReinsurerCode = treaty.ReinsurerCode
		}
		if err := DB.Create(&n).Error; err != nil {
			appLog.WithField("claim_number", claim.ClaimNumber).Warn("Failed to create month-end notification: " + err.Error())
			continue
		}
		created = append(created, n)
	}

	return created, nil
}

// GetNotificationStats returns counts by status for the given scheme (0 = all schemes).
func GetNotificationStats(schemeID int) (models.NotificationStats, error) {
	var stats models.NotificationStats
	ctx := context.Background()

	type Row struct {
		Status string
		Count  int
	}
	var rows []Row
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		q := d.Model(&models.ClaimNotificationLog{})
		if schemeID > 0 {
			q = q.Where("scheme_id = ?", schemeID)
		}
		return q.Select("status, count(*) as count").Group("status").Scan(&rows).Error
	})
	if err != nil {
		return stats, err
	}

	for _, r := range rows {
		stats.Total += r.Count
		switch r.Status {
		case "pending":
			stats.Pending = r.Count
		case "sent":
			stats.Sent = r.Count
		case "acknowledged":
			stats.Acknowledged = r.Count
		case "overdue":
			stats.Overdue = r.Count
		}
	}
	return stats, nil
}

// ClaimLookupItem is a lightweight struct for claim number dropdowns.
type ClaimLookupItem struct {
	ID          int    `json:"id"`
	ClaimNumber string `json:"claim_number"`
	MemberName  string `json:"member_name"`
	Status      string `json:"status"`
}

// GetClaimsByScheme returns a lightweight list of claims for a given scheme (for dropdown use).
func GetClaimsByScheme(schemeID int) ([]ClaimLookupItem, error) {
	var items []ClaimLookupItem
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Model(&models.GroupSchemeClaim{}).
			Select("id, claim_number, member_name, status").
			Where("scheme_id = ?", schemeID).
			Order("claim_number ASC").
			Scan(&items).Error
	})
	return items, err
}

// SweepOverdueNotifications bulk-updates all pending/sent notifications whose due date has passed.
// Designed to be called from a background ticker so overdue state is set even if nobody opens the UI.
func SweepOverdueNotifications() {
	result := DB.Model(&models.ClaimNotificationLog{}).
		Where("status IN ('pending','sent') AND due_date != '' AND due_date < ?", time.Now().Format("2006-01-02")).
		Update("status", "overdue")
	if result.Error != nil {
		appLog.Warn("Overdue notification sweep failed: " + result.Error.Error())
	} else if result.RowsAffected > 0 {
		appLog.WithField("count", result.RowsAffected).Info("Marked notifications as overdue")
	}
}

// StartNotificationOverdueSweeper runs a background goroutine that marks overdue notifications every 15 minutes.
func StartNotificationOverdueSweeper() {
	BackfillReconciliationItems()
	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()
		for {
			SweepOverdueNotifications()
			<-ticker.C
		}
	}()
}

// ExportNotificationsCSV returns CSV bytes for all notifications matching the given filters.
func ExportNotificationsCSV(schemeID int, status, notificationType string) ([]byte, error) {
	var records []models.ClaimNotificationLog
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		q := d.Order("created_at DESC")
		if schemeID > 0 {
			q = q.Where("scheme_id = ?", schemeID)
		}
		if status != "" {
			q = q.Where("status = ?", status)
		}
		if notificationType != "" {
			q = q.Where("notification_type = ?", notificationType)
		}
		return q.Find(&records).Error
	})
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"ID", "Claim Number", "Scheme", "Reinsurer", "Reinsurer Code", "Type", "Status", "Due Date", "Sent At", "Acknowledged At", "Notes", "Created By", "Created At"})
	for _, r := range records {
		sentAt := ""
		if r.SentAt != nil {
			sentAt = r.SentAt.Format("2006-01-02")
		}
		ackAt := ""
		if r.AcknowledgedAt != nil {
			ackAt = r.AcknowledgedAt.Format("2006-01-02")
		}
		_ = w.Write([]string{
			fmt.Sprintf("%d", r.ID),
			r.ClaimNumber,
			r.SchemeName,
			r.ReinsurerName,
			r.ReinsurerCode,
			r.NotificationType,
			r.Status,
			r.DueDate,
			sentAt,
			ackAt,
			r.Notes,
			r.CreatedBy,
			r.CreatedAt.Format("2006-01-02 15:04"),
		})
	}
	w.Flush()
	return buf.Bytes(), w.Error()
}
