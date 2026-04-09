package services

import (
	"api/models"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CreateIFRS17Amendment creates a new amendment record linking two CSM runs.
func CreateIFRS17Amendment(req models.CreateAmendmentRequest, user models.AppUser) (models.IFRS17Amendment, error) {
	var amendment models.IFRS17Amendment

	// Validate both run IDs exist.
	ctx := context.Background()
	var originalRun models.CsmRun
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Where("id = ?", req.OriginalRunID).First(&originalRun).Error
	})
	if err != nil {
		return amendment, errors.New("original run not found")
	}

	var amendmentRun models.CsmRun
	err = DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Where("id = ?", req.AmendmentRunID).First(&amendmentRun).Error
	})
	if err != nil {
		return amendment, errors.New("amendment run not found")
	}

	amendment = models.IFRS17Amendment{
		OriginalRunID:  req.OriginalRunID,
		AmendmentRunID: req.AmendmentRunID,
		AmendmentType:  req.AmendmentType,
		Reason:         req.Reason,
		Status:         "draft",
		CreatedBy:      user.UserName,
		CreatedAt:      time.Now(),
	}

	if err := DB.Create(&amendment).Error; err != nil {
		return amendment, err
	}
	return amendment, nil
}

// populateRunNames fills OriginalRunName and AmendmentRunName for a slice of amendments.
func populateRunNames(amendments []models.IFRS17Amendment) []models.IFRS17Amendment {
	ctx := context.Background()
	for i := range amendments {
		var orig models.CsmRun
		if err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
			return d.Where("id = ?", amendments[i].OriginalRunID).First(&orig).Error
		}); err == nil {
			amendments[i].OriginalRunName = orig.Name
		}

		var amend models.CsmRun
		if err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
			return d.Where("id = ?", amendments[i].AmendmentRunID).First(&amend).Error
		}); err == nil {
			amendments[i].AmendmentRunName = amend.Name
		}
	}
	return amendments
}

// GetAmendmentsForRun returns all amendments referencing a specific run.
func GetAmendmentsForRun(runID int) ([]models.IFRS17Amendment, error) {
	var amendments []models.IFRS17Amendment
	ctx := context.Background()
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Where("original_run_id = ? OR amendment_run_id = ?", runID, runID).
			Order("created_at DESC").Find(&amendments).Error
	})
	if err != nil {
		return amendments, err
	}
	amendments = populateRunNames(amendments)
	return amendments, nil
}

// GetAllAmendments returns all amendments optionally filtered by status.
func GetAllAmendments(status string) ([]models.IFRS17Amendment, error) {
	var amendments []models.IFRS17Amendment
	ctx := context.Background()
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		query := d.Order("created_at DESC")
		if status != "" {
			query = query.Where("status = ?", status)
		}
		return query.Find(&amendments).Error
	})
	if err != nil {
		return amendments, err
	}
	amendments = populateRunNames(amendments)
	return amendments, nil
}

// ApproveIFRS17Amendment approves an amendment record.
func ApproveIFRS17Amendment(id int, user models.AppUser) error {
	var amendment models.IFRS17Amendment
	if err := DB.First(&amendment, id).Error; err != nil {
		return err
	}
	now := time.Now()
	amendment.Status = "approved"
	amendment.ApprovedBy = user.UserName
	amendment.ApprovedAt = &now
	if err := DB.Save(&amendment).Error; err != nil {
		return err
	}
	go NotifyAmendmentApproved(amendment, user)
	return nil
}
