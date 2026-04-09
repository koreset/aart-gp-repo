package services

import (
	"api/models"
	"time"

	"gorm.io/gorm"
)

func CreateBordereauxTemplate(t *models.BordereauxTemplate) error {
	t.CreatedAt = time.Now()
	return DB.Create(t).Error
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

func UpdateBordereauxTemplate(id int, payload models.BordereauxTemplate) (models.BordereauxTemplate, error) {
	// Ensure we update the specific record
	payload.ID = id
	payload.UpdatedAt = time.Now()
	err := DB.Model(&models.BordereauxTemplate{}).Where("id = ?", id).Updates(payload).Error
	if err != nil {
		return payload, err
	}
	// Return the latest state
	return GetBordereauxTemplateByID(id)
}

func DeleteBordereauxTemplate(id int) error {
	return DB.Where("id = ?", id).Delete(&models.BordereauxTemplate{}).Error
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
