package services

import (
    "api/log"
    "api/models"
    "errors"
)

// CreateSchemeType creates a new scheme type
func CreateSchemeType(st *models.SchemeType) (*models.SchemeType, error) {
    if st.Name == "" {
        return nil, errors.New("name is required")
    }
    if err := DB.Create(st).Error; err != nil {
        log.WithField("error", err.Error()).Error("failed to create scheme type")
        return nil, err
    }
    return st, nil
}

// GetSchemeTypes returns all scheme types
func GetSchemeTypes() ([]models.SchemeType, error) {
    list := []models.SchemeType{}
    if err := DB.Order("name asc").Find(&list).Error; err != nil {
        return nil, err
    }
    return list, nil
}

// GetSchemeType returns a scheme type by id
func GetSchemeType(id int) (*models.SchemeType, error) {
    var st models.SchemeType
    if err := DB.First(&st, id).Error; err != nil {
        return nil, err
    }
    return &st, nil
}

// UpdateSchemeType updates name/description
func UpdateSchemeType(id int, payload *models.SchemeType) (*models.SchemeType, error) {
    var st models.SchemeType
    if err := DB.First(&st, id).Error; err != nil {
        return nil, err
    }
    if payload.Name != "" {
        st.Name = payload.Name
    }
    st.Description = payload.Description
    if err := DB.Save(&st).Error; err != nil {
        return nil, err
    }
    return &st, nil
}

// DeleteSchemeType deletes a scheme type. If in use, returns error.
func DeleteSchemeType(id int) error {
    // Check reference from quotes if a SchemeTypeID field exists
    type QuoteRef struct{ Count int }
    var cnt int64
    if DB.Migrator().HasColumn(&models.GroupPricingQuote{}, "scheme_type_id") {
        DB.Model(&models.GroupPricingQuote{}).Where("scheme_type_id = ?", id).Count(&cnt)
        if cnt > 0 {
            return errors.New("cannot delete scheme type that is referenced by one or more quotes")
        }
    }
    return DB.Delete(&models.SchemeType{}, id).Error
}
