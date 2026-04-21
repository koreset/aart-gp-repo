package services

import (
	"api/models"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// validateBordereauxFieldMappings rejects templates that reuse the same
// source_field more than once — each source field may appear at most once per
// template.
func validateBordereauxFieldMappings(mappings []models.BordereauxFieldMapping) error {
	seen := make(map[string]struct{}, len(mappings))
	for _, m := range mappings {
		if m.SourceField == "" {
			continue
		}
		if _, ok := seen[m.SourceField]; ok {
			return fmt.Errorf("duplicate source field: %s", m.SourceField)
		}
		seen[m.SourceField] = struct{}{}
	}
	return nil
}

func CreateBordereauxTemplate(t *models.BordereauxTemplate, user models.AppUser) error {
	if err := validateBordereauxFieldMappings(t.FieldMappings); err != nil {
		return err
	}
	t.CreatedAt = time.Now()
	if err := DB.Create(t).Error; err != nil {
		return err
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "bordereaux_templates",
		EntityID:  strconv.Itoa(t.ID),
		Action:    "CREATE",
		ChangedBy: user.UserName,
	}, struct{}{}, *t)
	return nil
}

func GetBordereauxTemplates() ([]models.BordereauxTemplate, error) {
	var list []models.BordereauxTemplate
	err := DB.Find(&list).Error
	return list, err
}

func GetBordereauxTemplateByID(id int) (models.BordereauxTemplate, error) {
	var t models.BordereauxTemplate
	err := DB.Where("id = ?", id).First(&t).Error
	return t, err
}

func UpdateBordereauxTemplate(id int, payload models.BordereauxTemplate, user models.AppUser) (models.BordereauxTemplate, error) {
	if err := validateBordereauxFieldMappings(payload.FieldMappings); err != nil {
		return payload, err
	}
	before, _ := GetBordereauxTemplateByID(id)
	// Ensure we update the specific record
	payload.ID = id
	payload.UpdatedAt = time.Now()
	err := DB.Model(&models.BordereauxTemplate{}).Where("id = ?", id).Updates(payload).Error
	if err != nil {
		return payload, err
	}
	// Return the latest state
	updated, err := GetBordereauxTemplateByID(id)
	if err == nil {
		_ = writeAudit(DB, AuditContext{
			Area:      "group-pricing",
			Entity:    "bordereaux_templates",
			EntityID:  strconv.Itoa(id),
			Action:    "UPDATE",
			ChangedBy: user.UserName,
		}, before, updated)
	}
	return updated, err
}

func DeleteBordereauxTemplate(id int, user models.AppUser) error {
	before, _ := GetBordereauxTemplateByID(id)
	if err := DB.Where("id = ?", id).Delete(&models.BordereauxTemplate{}).Error; err != nil {
		return err
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "bordereaux_templates",
		EntityID:  strconv.Itoa(id),
		Action:    "DELETE",
		ChangedBy: user.UserName,
	}, before, struct{}{})
	return nil
}

func IncrementTemplateUsage(id int) error {
	var usageCount int64
	DB.Model(&models.BordereauxTemplate{}).Where("id = ?", id).Select("usage_count").Row().Scan(&usageCount)

	result := DB.Model(&models.BordereauxTemplate{}).Where("id = ?", id).Updates(map[string]interface{}{"usage_count": usageCount + 1})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
