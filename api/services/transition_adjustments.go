package services

import (
	"api/models"
	"context"
	"time"

	"gorm.io/gorm"
)

// GetTransitionAdjustments returns all transition adjustments optionally filtered by status.
func GetTransitionAdjustments(status string) ([]models.TransitionAdjustment, error) {
	var records []models.TransitionAdjustment
	ctx := context.Background()
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		query := d.Order("created_at DESC")
		if status != "" {
			query = query.Where("status = ?", status)
		}
		return query.Find(&records).Error
	})
	return records, err
}

// ApproveTransitionAdjustment sets a transition adjustment's status to approved.
func ApproveTransitionAdjustment(id int, notes string, user models.AppUser) error {
	var record models.TransitionAdjustment
	if err := DB.First(&record, id).Error; err != nil {
		return err
	}
	now := time.Now()
	record.Status = "approved"
	record.ReviewedBy = user.UserName
	record.ReviewedAt = &now
	record.Notes = notes
	return DB.Save(&record).Error
}

// RejectTransitionAdjustment sets a transition adjustment's status to rejected.
func RejectTransitionAdjustment(id int, notes string, user models.AppUser) error {
	var record models.TransitionAdjustment
	if err := DB.First(&record, id).Error; err != nil {
		return err
	}
	now := time.Now()
	record.Status = "rejected"
	record.ReviewedBy = user.UserName
	record.ReviewedAt = &now
	record.Notes = notes
	return DB.Save(&record).Error
}

// BulkApproveTransitionAdjustments approves multiple transition adjustments.
func BulkApproveTransitionAdjustments(ids []int, user models.AppUser) error {
	for _, id := range ids {
		if err := ApproveTransitionAdjustment(id, "", user); err != nil {
			return err
		}
	}
	return nil
}
