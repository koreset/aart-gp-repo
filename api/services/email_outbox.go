package services

import (
	"errors"
	"fmt"
	"time"

	"api/models"

	"gorm.io/gorm"
)

// ListEmailOutboxInput filters and paginates the outbox list endpoint.
type ListEmailOutboxInput struct {
	LicenseId string
	Status    string // "", "pending", "sending", "sent", "failed"
	Page      int
	PageSize  int
}

// ListEmailOutboxResult is the paginated result.
type ListEmailOutboxResult struct {
	Items      []models.EmailOutbox `json:"items"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// ListEmailOutbox returns a paginated slice of outbox rows, newest first.
func ListEmailOutbox(in ListEmailOutboxInput) (ListEmailOutboxResult, error) {
	if in.LicenseId == "" {
		return ListEmailOutboxResult{}, errors.New("license_id is required")
	}
	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize < 1 || in.PageSize > 200 {
		in.PageSize = 25
	}

	q := DB.Model(&models.EmailOutbox{}).Where("license_id = ?", in.LicenseId)
	if in.Status != "" {
		q = q.Where("status = ?", in.Status)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return ListEmailOutboxResult{}, err
	}

	var items []models.EmailOutbox
	if err := q.Order("id DESC").
		Offset((in.Page - 1) * in.PageSize).
		Limit(in.PageSize).
		Find(&items).Error; err != nil {
		return ListEmailOutboxResult{}, err
	}

	totalPages := int((total + int64(in.PageSize) - 1) / int64(in.PageSize))
	return ListEmailOutboxResult{
		Items:      items,
		Total:      total,
		Page:       in.Page,
		PageSize:   in.PageSize,
		TotalPages: totalPages,
	}, nil
}

// RetryEmailOutbox resets a failed (or idle pending) row so the worker picks
// it up on its next poll. Returns an error if the row isn't in a state that
// can be retried.
func RetryEmailOutbox(licenseId string, id int) (models.EmailOutbox, error) {
	var row models.EmailOutbox
	if err := DB.Where("id = ? AND license_id = ?", id, licenseId).First(&row).Error; err != nil {
		return row, err
	}
	if row.Status == models.EmailOutboxSent {
		return row, errors.New("already sent")
	}
	if row.Status == models.EmailOutboxSending {
		return row, errors.New("currently sending — wait for the attempt to finish")
	}
	updates := map[string]interface{}{
		"status":          models.EmailOutboxPending,
		"next_attempt_at": time.Now(),
		"last_error":      "",
	}
	if err := DB.Model(&models.EmailOutbox{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return row, fmt.Errorf("update row: %w", err)
	}
	return GetEmailOutboxByID(licenseId, id)
}

// GetEmailOutboxByID returns one outbox row, scoped to the caller's license.
func GetEmailOutboxByID(licenseId string, id int) (models.EmailOutbox, error) {
	var row models.EmailOutbox
	err := DB.Where("id = ? AND license_id = ?", id, licenseId).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return row, err
	}
	return row, err
}
