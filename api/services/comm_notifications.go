package services

import (
	"api/log"
	"api/models"
	"context"
	"time"

	"gorm.io/gorm"
)

// CreateNotification persists a notification and pushes it via WebSocket if the recipient is online.
func CreateNotification(req models.CreateNotificationRequest) (models.Notification, error) {
	n := models.Notification{
		RecipientEmail: req.RecipientEmail,
		SenderEmail:    req.SenderEmail,
		SenderName:     req.SenderName,
		Type:           req.Type,
		Title:          req.Title,
		Body:           req.Body,
		ObjectType:     req.ObjectType,
		ObjectID:       req.ObjectID,
	}
	if err := DB.Create(&n).Error; err != nil {
		return n, err
	}

	// Push via WebSocket
	if hub := GetHub(); hub != nil {
		hub.SendToUser(req.RecipientEmail, WSEnvelope{
			Type:    WSNotification,
			Payload: n,
		})
	}

	return n, nil
}

// GetNotifications returns a paginated list of notifications for a user.
func GetNotifications(email string, filter models.NotificationFilter) (models.NotificationListResponse, error) {
	var resp models.NotificationListResponse
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 100 {
		filter.PageSize = 20
	}
	resp.Page = filter.Page
	resp.PageSize = filter.PageSize

	err := DBReadWithResilience(context.Background(), func(db *gorm.DB) error {
		query := db.Where("recipient_email = ?", email)
		if filter.IsRead != nil {
			query = query.Where("is_read = ?", *filter.IsRead)
		}
		if filter.Type != "" {
			query = query.Where("type = ?", filter.Type)
		}

		var total int64
		if err := query.Model(&models.Notification{}).Count(&total).Error; err != nil {
			return err
		}
		resp.Total = total

		offset := (filter.Page - 1) * filter.PageSize
		return query.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&resp.Notifications).Error
	})

	return resp, err
}

// GetUnreadNotificationCount returns the count of unread notifications for a user.
func GetUnreadNotificationCount(email string) (int64, error) {
	var count int64
	err := DBReadWithResilience(context.Background(), func(db *gorm.DB) error {
		return db.Model(&models.Notification{}).
			Where("recipient_email = ? AND is_read = ?", email, false).
			Count(&count).Error
	})
	return count, err
}

// MarkNotificationAsRead marks a single notification as read.
func MarkNotificationAsRead(id int, email string) error {
	now := time.Now()
	return DB.Model(&models.Notification{}).
		Where("id = ? AND recipient_email = ?", id, email).
		Updates(map[string]interface{}{"is_read": true, "read_at": now}).Error
}

// MarkAllNotificationsAsRead marks all notifications as read for a user.
func MarkAllNotificationsAsRead(email string) error {
	now := time.Now()
	return DB.Model(&models.Notification{}).
		Where("recipient_email = ? AND is_read = ?", email, false).
		Updates(map[string]interface{}{"is_read": true, "read_at": now}).Error
}

// DeleteNotification removes a notification by ID (only if it belongs to the user).
func DeleteNotification(id int, email string) error {
	result := DB.Where("id = ? AND recipient_email = ?", id, email).Delete(&models.Notification{})
	if result.RowsAffected == 0 {
		log.WithField("notification_id", id).Warn("Notification not found or not owned by user")
	}
	return result.Error
}
