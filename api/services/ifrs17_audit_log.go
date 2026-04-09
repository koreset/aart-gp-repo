package services

import (
	"api/models"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// LogIFRS17Event writes an audit log entry. It is fire-and-forget: errors are
// printed to stdout and never propagated to the caller.
func LogIFRS17Event(eventType, objectType, objectName string, objectID int, user models.AppUser, details string) {
	entry := models.IFRS17AuditLog{
		EventType:  eventType,
		ObjectType: objectType,
		ObjectID:   objectID,
		ObjectName: objectName,
		ChangedBy:  user.UserName,
		ChangedAt:  time.Now(),
		Details:    details,
	}
	if err := DB.Create(&entry).Error; err != nil {
		fmt.Println("LogIFRS17Event: failed to write audit log:", err)
	}
}

// GetIFRS17AuditLog returns audit log entries with optional filtering.
// eventType, from, and to are all optional (empty string = no filter).
// Results are ordered newest-first.
func GetIFRS17AuditLog(eventType, from, to string) ([]models.IFRS17AuditLog, error) {
	var logs []models.IFRS17AuditLog
	ctx := context.Background()
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		q := d.Model(&models.IFRS17AuditLog{})
		if eventType != "" {
			q = q.Where("event_type = ?", eventType)
		}
		if from != "" {
			q = q.Where("changed_at >= ?", from)
		}
		if to != "" {
			q = q.Where("changed_at <= ?", to)
		}
		return q.Order("changed_at DESC").Find(&logs).Error
	})
	if err != nil {
		return nil, err
	}
	return logs, nil
}
