package services

import (
	"errors"
	"fmt"
	"time"

	"api/models"
	"api/services/email"

	"gorm.io/gorm"
)

// ListEmailTemplates returns all templates for a license, ordered by code.
func ListEmailTemplates(licenseId string) ([]models.EmailTemplate, error) {
	var list []models.EmailTemplate
	err := DB.Where("license_id = ?", licenseId).Order("code ASC").Find(&list).Error
	return list, err
}

// GetEmailTemplate returns a single template by (license, code).
func GetEmailTemplate(licenseId, code string) (models.EmailTemplate, error) {
	var t models.EmailTemplate
	err := DB.Where("license_id = ? AND code = ?", licenseId, code).First(&t).Error
	return t, err
}

// SaveEmailTemplateInput is the sanitised form accepted by Save.
type SaveEmailTemplateInput struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	SubjectTemplate string `json:"subject_template"`
	BodyTemplate    string `json:"body_template"`
	AttachmentsSpec string `json:"attachments_spec"`
	Status          string `json:"status"`
}

// CreateEmailTemplate inserts a new template. Returns an error if a template
// with the same code already exists for the license.
func CreateEmailTemplate(licenseId string, in SaveEmailTemplateInput, user models.AppUser) (models.EmailTemplate, error) {
	if licenseId == "" {
		return models.EmailTemplate{}, errors.New("license_id is required")
	}
	if in.Code == "" {
		return models.EmailTemplate{}, errors.New("code is required")
	}
	status := normaliseTemplateStatus(in.Status)
	_, err := GetEmailTemplate(licenseId, in.Code)
	if err == nil {
		return models.EmailTemplate{}, fmt.Errorf("template with code %q already exists", in.Code)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.EmailTemplate{}, err
	}
	row := models.EmailTemplate{
		LicenseId:       licenseId,
		Code:            in.Code,
		Name:            in.Name,
		Description:     in.Description,
		SubjectTemplate: in.SubjectTemplate,
		BodyTemplate:    in.BodyTemplate,
		AttachmentsSpec: in.AttachmentsSpec,
		Status:          status,
		UpdatedBy:       user.UserName,
		UpdatedAt:       time.Now(),
	}
	if err := DB.Create(&row).Error; err != nil {
		return row, err
	}
	return row, nil
}

// UpdateEmailTemplate updates a template in place, keyed by (license, code).
func UpdateEmailTemplate(licenseId, code string, in SaveEmailTemplateInput, user models.AppUser) (models.EmailTemplate, error) {
	existing, err := GetEmailTemplate(licenseId, code)
	if err != nil {
		return existing, err
	}
	existing.Name = in.Name
	existing.Description = in.Description
	existing.SubjectTemplate = in.SubjectTemplate
	existing.BodyTemplate = in.BodyTemplate
	existing.AttachmentsSpec = in.AttachmentsSpec
	if in.Status != "" {
		existing.Status = normaliseTemplateStatus(in.Status)
	}
	existing.UpdatedBy = user.UserName
	existing.UpdatedAt = time.Now()
	if err := DB.Save(&existing).Error; err != nil {
		return existing, err
	}
	return existing, nil
}

// DeleteEmailTemplate removes a template. Outbox rows keep their template_code
// string as-is; they're historical records.
func DeleteEmailTemplate(licenseId, code string) error {
	res := DB.Where("license_id = ? AND code = ?", licenseId, code).Delete(&models.EmailTemplate{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// PreviewEmailTemplate renders a template in place (no persistence) using the
// supplied vars. Used by the admin UI's preview pane.
func PreviewEmailTemplate(tpl models.EmailTemplate, vars map[string]interface{}) (email.Rendered, error) {
	return email.Preview(tpl, vars)
}

func normaliseTemplateStatus(s string) string {
	switch s {
	case models.EmailTemplateActive:
		return models.EmailTemplateActive
	default:
		return models.EmailTemplateDraft
	}
}
