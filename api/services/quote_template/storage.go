package quote_template

import (
	"fmt"

	"api/models"
	"api/services"
)

// SaveTemplate saves a new template version for an insurer, setting it as active
func SaveTemplate(insurerID int, filename string, blob []byte, uploadedBy string) (*models.InsurerQuoteTemplate, error) {
	if len(blob) == 0 {
		return nil, fmt.Errorf("template blob is empty")
	}

	// Start transaction
	tx := services.DB.Begin()

	// Get the highest version number for this insurer
	var maxVersion int
	err := tx.Model(&models.InsurerQuoteTemplate{}).
		Where("insurer_id = ?", insurerID).
		Select("COALESCE(MAX(version), 0)").
		Row().
		Scan(&maxVersion)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to query max version: %w", err)
	}

	newVersion := maxVersion + 1

	// Deactivate all previous templates for this insurer
	err = tx.Model(&models.InsurerQuoteTemplate{}).
		Where("insurer_id = ?", insurerID).
		Update("is_active", false).
		Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to deactivate previous templates: %w", err)
	}

	// Create new template
	template := &models.InsurerQuoteTemplate{
		InsurerID:  insurerID,
		Version:    newVersion,
		Filename:   filename,
		DocxBlob:   blob,
		SizeBytes:  len(blob),
		UploadedBy: uploadedBy,
		IsActive:   true,
	}

	err = tx.Create(template).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	// Commit transaction
	tx.Commit()

	return template, nil
}

// GetActiveTemplate returns the active template for an insurer, or (nil, nil) if none exists
func GetActiveTemplate(insurerID int) (*models.InsurerQuoteTemplate, error) {
	var template models.InsurerQuoteTemplate
	result := services.DB.
		Where("insurer_id = ? AND is_active = ?", insurerID, true).
		First(&template)

	if result.Error != nil {
		// Check if it's a "not found" error
		if result.RowsAffected == 0 {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query active template: %w", result.Error)
	}

	return &template, nil
}

// ListTemplates returns all templates for an insurer in descending version order (no blob)
func ListTemplates(insurerID int) ([]models.InsurerQuoteTemplate, error) {
	var templates []models.InsurerQuoteTemplate
	err := services.DB.
		Select("id", "insurer_id", "version", "filename", "size_bytes", "uploaded_by", "uploaded_at", "is_active").
		Where("insurer_id = ?", insurerID).
		Order("version DESC").
		Find(&templates).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	return templates, nil
}

// GetTemplate fetches a template by ID including the blob
func GetTemplate(id int) (*models.InsurerQuoteTemplate, error) {
	var template models.InsurerQuoteTemplate
	result := services.DB.First(&template, id)

	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query template: %w", result.Error)
	}

	return &template, nil
}

// ActivateVersion sets the given template as active and deactivates all others for that insurer
func ActivateVersion(insurerID, templateID int) error {
	// Get the template to verify it belongs to this insurer
	var template models.InsurerQuoteTemplate
	result := services.DB.First(&template, templateID)
	if result.Error != nil {
		return fmt.Errorf("template not found: %w", result.Error)
	}

	if template.InsurerID != insurerID {
		return fmt.Errorf("template does not belong to insurer")
	}

	// Start transaction
	tx := services.DB.Begin()

	// Deactivate all templates for this insurer
	err := tx.Model(&models.InsurerQuoteTemplate{}).
		Where("insurer_id = ?", insurerID).
		Update("is_active", false).
		Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to deactivate templates: %w", err)
	}

	// Activate the selected template
	err = tx.Model(&models.InsurerQuoteTemplate{}).
		Where("id = ?", templateID).
		Update("is_active", true).
		Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to activate template: %w", err)
	}

	tx.Commit()
	return nil
}
