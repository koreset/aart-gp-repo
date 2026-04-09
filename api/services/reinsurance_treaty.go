package services

import (
	"api/models"
	"fmt"
	"time"
)

// CreateTreaty creates a new reinsurance treaty
func CreateTreaty(req models.CreateTreatyRequest, user models.AppUser) (models.ReinsuranceTreaty, error) {
	// Enforce composite uniqueness: a treaty number is only unique per line of business
	var count int64
	lob := req.LineOfBusiness
	if lob == "" {
		lob = "group_life"
	}
	DB.Model(&models.ReinsuranceTreaty{}).
		Where("treaty_number = ? AND line_of_business = ?", req.TreatyNumber, lob).
		Count(&count)
	if count > 0 {
		return models.ReinsuranceTreaty{}, fmt.Errorf(
			"treaty %q already exists for line of business %q; use a different treaty number or a different LoB",
			req.TreatyNumber, lob,
		)
	}

	currency := req.Currency
	if currency == "" {
		currency = "ZAR"
	}
	frequency := req.PremiumPaymentFrequency
	if frequency == "" {
		frequency = "monthly"
	}
	notificationDays := req.ClaimsNotificationDays
	if notificationDays == 0 {
		notificationDays = 30
	}
	t := models.ReinsuranceTreaty{
		TreatyNumber:              req.TreatyNumber,
		TreatyName:                req.TreatyName,
		ReinsurerName:             req.ReinsurerName,
		ReinsurerCode:             req.ReinsurerCode,
		BrokerName:                req.BrokerName,
		TreatyType:                req.TreatyType,
		LineOfBusiness:            req.LineOfBusiness,
		EffectiveDate:             req.EffectiveDate,
		ExpiryDate:                req.ExpiryDate,
		RenewalDate:               req.RenewalDate,
		Status:                    "draft",
		Currency:                  currency,
		RetentionType:             req.RetentionType,
		RetentionAmount:           req.RetentionAmount,
		RetentionPercentage:       req.RetentionPercentage,
		SurplusLines:              req.SurplusLines,
		XLRetention:               req.XLRetention,
		XLLimit:                   req.XLLimit,
		XLLayerFrom:               req.XLLayerFrom,
		XLLayerTo:                 req.XLLayerTo,
		AggregateAnnualLimit:      req.AggregateAnnualLimit,
		ProfitCommissionRate:      req.ProfitCommissionRate,
		ReinsuranceCommissionRate: req.ReinsuranceCommissionRate,
		PremiumPaymentFrequency:   frequency,
		ClaimsNotificationDays:    notificationDays,
		LargeClaimsThreshold:      req.LargeClaimsThreshold,
		DeltaReporting:            req.DeltaReporting,
		Notes:                     req.Notes,
		CreatedBy:                 user.UserName,
		// Three-Tier Reinsurance Structure
		TreatyCode:                  req.TreatyCode,
		RiskPremiumBasisIndicator:   req.RiskPremiumBasisIndicator,
		FlatAnnualReinsPremRate:     req.FlatAnnualReinsPremRate,
		Level1CededProportion:       req.Level1CededProportion,
		Level1Lowerbound:            req.Level1Lowerbound,
		Level1Upperbound:            req.Level1Upperbound,
		Level2CededProportion:       req.Level2CededProportion,
		Level2Lowerbound:            req.Level2Lowerbound,
		Level2Upperbound:            req.Level2Upperbound,
		Level3CededProportion:       req.Level3CededProportion,
		Level3Lowerbound:            req.Level3Lowerbound,
		Level3Upperbound:            req.Level3Upperbound,
		IncomeLevel1CededProportion: req.IncomeLevel1CededProportion,
		IncomeLevel1Lowerbound:      req.IncomeLevel1Lowerbound,
		IncomeLevel1Upperbound:      req.IncomeLevel1Upperbound,
		IncomeLevel2CededProportion: req.IncomeLevel2CededProportion,
		IncomeLevel2Lowerbound:      req.IncomeLevel2Lowerbound,
		IncomeLevel2Upperbound:      req.IncomeLevel2Upperbound,
		IncomeLevel3CededProportion: req.IncomeLevel3CededProportion,
		IncomeLevel3Lowerbound:      req.IncomeLevel3Lowerbound,
		IncomeLevel3Upperbound:      req.IncomeLevel3Upperbound,
		LeadReinsurerShare:          req.LeadReinsurerShare,
		LeadReinsurerCode:           req.LeadReinsurerCode,
		NonLeadReinsurer1Share:      req.NonLeadReinsurer1Share,
		NonLeadReinsurer1Code:       req.NonLeadReinsurer1Code,
		NonLeadReinsurer2Share:      req.NonLeadReinsurer2Share,
		NonLeadReinsurer2Code:       req.NonLeadReinsurer2Code,
		NonLeadReinsurer3Share:      req.NonLeadReinsurer3Share,
		NonLeadReinsurer3Code:       req.NonLeadReinsurer3Code,
		CedingCommission:            req.CedingCommission,
	}
	if err := DB.Create(&t).Error; err != nil {
		return t, fmt.Errorf("failed to create treaty: %w", err)
	}
	return t, nil
}

// GetTreaties returns treaties filtered by optional status, type and reinsurer
func GetTreaties(status, treatyType, reinsurer string) ([]models.ReinsuranceTreaty, error) {
	var treaties []models.ReinsuranceTreaty
	q := DB.Order("created_at DESC")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if treatyType != "" {
		q = q.Where("treaty_type = ?", treatyType)
	}
	if reinsurer != "" {
		q = q.Where("reinsurer_name LIKE ?", "%"+reinsurer+"%")
	}
	if err := q.Find(&treaties).Error; err != nil {
		return nil, err
	}
	return treaties, nil
}

// GetTreatyByID returns a single treaty by ID
func GetTreatyByID(id int) (models.ReinsuranceTreaty, error) {
	var t models.ReinsuranceTreaty
	if err := DB.First(&t, id).Error; err != nil {
		return t, fmt.Errorf("treaty %d not found: %w", id, err)
	}
	return t, nil
}

// UpdateTreaty updates treaty fields
func UpdateTreaty(id int, req models.UpdateTreatyRequest, user models.AppUser) (models.ReinsuranceTreaty, error) {
	t, err := GetTreatyByID(id)
	if err != nil {
		return t, err
	}
	if req.TreatyName != "" {
		t.TreatyName = req.TreatyName
	}
	if req.ReinsurerName != "" {
		t.ReinsurerName = req.ReinsurerName
	}
	if req.ReinsurerCode != "" {
		t.ReinsurerCode = req.ReinsurerCode
	}
	if req.BrokerName != "" {
		t.BrokerName = req.BrokerName
	}
	if req.TreatyType != "" {
		t.TreatyType = req.TreatyType
	}
	if req.LineOfBusiness != "" {
		t.LineOfBusiness = req.LineOfBusiness
	}
	if req.EffectiveDate != "" {
		t.EffectiveDate = req.EffectiveDate
	}
	if req.ExpiryDate != "" {
		t.ExpiryDate = req.ExpiryDate
	}
	if req.RenewalDate != "" {
		t.RenewalDate = req.RenewalDate
	}
	if req.Status != "" {
		t.Status = req.Status
	}
	if req.Currency != "" {
		t.Currency = req.Currency
	}
	if req.RetentionType != "" {
		t.RetentionType = req.RetentionType
	}
	if req.RetentionAmount > 0 {
		t.RetentionAmount = req.RetentionAmount
	}
	if req.RetentionPercentage > 0 {
		t.RetentionPercentage = req.RetentionPercentage
	}
	if req.SurplusLines > 0 {
		t.SurplusLines = req.SurplusLines
	}
	if req.XLRetention > 0 {
		t.XLRetention = req.XLRetention
	}
	if req.XLLimit > 0 {
		t.XLLimit = req.XLLimit
	}
	if req.XLLayerFrom > 0 {
		t.XLLayerFrom = req.XLLayerFrom
	}
	if req.XLLayerTo > 0 {
		t.XLLayerTo = req.XLLayerTo
	}
	if req.AggregateAnnualLimit > 0 {
		t.AggregateAnnualLimit = req.AggregateAnnualLimit
	}
	if req.ProfitCommissionRate > 0 {
		t.ProfitCommissionRate = req.ProfitCommissionRate
	}
	if req.ReinsuranceCommissionRate > 0 {
		t.ReinsuranceCommissionRate = req.ReinsuranceCommissionRate
	}
	if req.PremiumPaymentFrequency != "" {
		t.PremiumPaymentFrequency = req.PremiumPaymentFrequency
	}
	if req.ClaimsNotificationDays > 0 {
		t.ClaimsNotificationDays = req.ClaimsNotificationDays
	}
	if req.LargeClaimsThreshold > 0 {
		t.LargeClaimsThreshold = req.LargeClaimsThreshold
	}
	t.DeltaReporting = req.DeltaReporting
	if req.Notes != "" {
		t.Notes = req.Notes
	}
	// Three-Tier Reinsurance Structure Updates
	if req.TreatyCode != "" {
		t.TreatyCode = req.TreatyCode
	}
	if req.RiskPremiumBasisIndicator != "" {
		t.RiskPremiumBasisIndicator = req.RiskPremiumBasisIndicator
	}
	if req.FlatAnnualReinsPremRate >= 0 {
		t.FlatAnnualReinsPremRate = req.FlatAnnualReinsPremRate
	}
	if req.Level1CededProportion >= 0 {
		t.Level1CededProportion = req.Level1CededProportion
	}
	if req.Level1Lowerbound >= 0 {
		t.Level1Lowerbound = req.Level1Lowerbound
	}
	if req.Level1Upperbound >= 0 {
		t.Level1Upperbound = req.Level1Upperbound
	}
	if req.Level2CededProportion >= 0 {
		t.Level2CededProportion = req.Level2CededProportion
	}
	if req.Level2Lowerbound >= 0 {
		t.Level2Lowerbound = req.Level2Lowerbound
	}
	if req.Level2Upperbound >= 0 {
		t.Level2Upperbound = req.Level2Upperbound
	}
	if req.Level3CededProportion >= 0 {
		t.Level3CededProportion = req.Level3CededProportion
	}
	if req.Level3Lowerbound >= 0 {
		t.Level3Lowerbound = req.Level3Lowerbound
	}
	if req.Level3Upperbound >= 0 {
		t.Level3Upperbound = req.Level3Upperbound
	}
	if req.IncomeLevel1CededProportion >= 0 {
		t.IncomeLevel1CededProportion = req.IncomeLevel1CededProportion
	}
	if req.IncomeLevel1Lowerbound >= 0 {
		t.IncomeLevel1Lowerbound = req.IncomeLevel1Lowerbound
	}
	if req.IncomeLevel1Upperbound >= 0 {
		t.IncomeLevel1Upperbound = req.IncomeLevel1Upperbound
	}
	if req.IncomeLevel2CededProportion >= 0 {
		t.IncomeLevel2CededProportion = req.IncomeLevel2CededProportion
	}
	if req.IncomeLevel2Lowerbound >= 0 {
		t.IncomeLevel2Lowerbound = req.IncomeLevel2Lowerbound
	}
	if req.IncomeLevel2Upperbound >= 0 {
		t.IncomeLevel2Upperbound = req.IncomeLevel2Upperbound
	}
	if req.IncomeLevel3CededProportion >= 0 {
		t.IncomeLevel3CededProportion = req.IncomeLevel3CededProportion
	}
	if req.IncomeLevel3Lowerbound >= 0 {
		t.IncomeLevel3Lowerbound = req.IncomeLevel3Lowerbound
	}
	if req.IncomeLevel3Upperbound >= 0 {
		t.IncomeLevel3Upperbound = req.IncomeLevel3Upperbound
	}
	if req.LeadReinsurerShare >= 0 {
		t.LeadReinsurerShare = req.LeadReinsurerShare
	}
	if req.LeadReinsurerCode != "" {
		t.LeadReinsurerCode = req.LeadReinsurerCode
	}
	if req.NonLeadReinsurer1Share >= 0 {
		t.NonLeadReinsurer1Share = req.NonLeadReinsurer1Share
	}
	if req.NonLeadReinsurer1Code != "" {
		t.NonLeadReinsurer1Code = req.NonLeadReinsurer1Code
	}
	if req.NonLeadReinsurer2Share >= 0 {
		t.NonLeadReinsurer2Share = req.NonLeadReinsurer2Share
	}
	if req.NonLeadReinsurer2Code != "" {
		t.NonLeadReinsurer2Code = req.NonLeadReinsurer2Code
	}
	if req.NonLeadReinsurer3Share >= 0 {
		t.NonLeadReinsurer3Share = req.NonLeadReinsurer3Share
	}
	if req.NonLeadReinsurer3Code != "" {
		t.NonLeadReinsurer3Code = req.NonLeadReinsurer3Code
	}
	if req.CedingCommission >= 0 {
		t.CedingCommission = req.CedingCommission
	}
	t.UpdatedBy = user.UserName
	t.UpdatedAt = time.Now()
	if err := DB.Save(&t).Error; err != nil {
		return t, fmt.Errorf("failed to update treaty: %w", err)
	}
	return t, nil
}

// DeleteTreaty deletes a treaty by ID
func DeleteTreaty(id int) error {
	if err := DB.Delete(&models.ReinsuranceTreaty{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete treaty: %w", err)
	}
	return nil
}

// LinkSchemeToTreaty creates a scheme–treaty link
func LinkSchemeToTreaty(treatyID int, req models.LinkSchemeRequest, user models.AppUser) (models.TreatySchemeLink, error) {
	link := models.TreatySchemeLink{
		TreatyID:        treatyID,
		SchemeID:        req.SchemeID,
		SchemeName:      req.SchemeName,
		CessionOverride: req.CessionOverride,
		EffectiveDate:   req.EffectiveDate,
		ExpiryDate:      req.ExpiryDate,
		CreatedBy:       user.UserName,
	}
	if err := DB.Create(&link).Error; err != nil {
		return link, fmt.Errorf("failed to link scheme to treaty: %w", err)
	}
	return link, nil
}

// BulkLinkSchemesToTreaty links multiple schemes to a treaty in one operation, skipping already-linked ones
func BulkLinkSchemesToTreaty(treatyID int, req models.BulkLinkSchemesRequest, user models.AppUser) (int, error) {
	created := 0
	for _, schemeID := range req.SchemeIDs {
		var existing int64
		DB.Model(&models.TreatySchemeLink{}).
			Where("treaty_id = ? AND scheme_id = ?", treatyID, schemeID).
			Count(&existing)
		if existing > 0 {
			continue
		}
		var scheme models.GroupScheme
		DB.First(&scheme, schemeID)
		link := models.TreatySchemeLink{
			TreatyID:        treatyID,
			SchemeID:        schemeID,
			SchemeName:      scheme.Name,
			CessionOverride: req.CessionOverride,
			EffectiveDate:   req.EffectiveDate,
			CreatedBy:       user.UserName,
		}
		if err := DB.Create(&link).Error; err == nil {
			created++
		}
	}
	return created, nil
}

// GetTreatySchemeLinks returns all scheme links for a treaty
func GetTreatySchemeLinks(treatyID int) ([]models.TreatySchemeLink, error) {
	var links []models.TreatySchemeLink
	if err := DB.Where("treaty_id = ?", treatyID).Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

// RemoveSchemeTreatyLink removes a scheme–treaty link
func RemoveSchemeTreatyLink(linkID int) error {
	if err := DB.Delete(&models.TreatySchemeLink{}, linkID).Error; err != nil {
		return fmt.Errorf("failed to remove scheme link: %w", err)
	}
	return nil
}

// BulkRemoveSchemeLinks removes multiple scheme–treaty links by their IDs
func BulkRemoveSchemeLinks(linkIDs []int) (int64, error) {
	if len(linkIDs) == 0 {
		return 0, nil
	}
	result := DB.Where("id IN ?", linkIDs).Delete(&models.TreatySchemeLink{})
	if result.Error != nil {
		return 0, fmt.Errorf("failed to bulk remove scheme links: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// GetActiveTreatiesForScheme returns active treaties linked to a scheme
func GetActiveTreatiesForScheme(schemeID int) ([]models.ReinsuranceTreaty, error) {
	var treaties []models.ReinsuranceTreaty
	err := DB.Joins("JOIN treaty_scheme_links ON treaty_scheme_links.treaty_id = reinsurance_treaties.id").
		Where("treaty_scheme_links.scheme_id = ? AND reinsurance_treaties.status = 'active'", schemeID).
		Find(&treaties).Error
	if err != nil {
		return nil, err
	}
	return treaties, nil
}

// GetTreatyStats returns a summary count of treaties by status including expiry warnings
func GetTreatyStats() (models.TreatyStats, error) {
	var treaties []models.ReinsuranceTreaty
	if err := DB.Find(&treaties).Error; err != nil {
		return models.TreatyStats{}, err
	}
	stats := models.TreatyStats{Total: len(treaties)}
	today := time.Now()
	in60 := today.AddDate(0, 0, 60).Format("2006-01-02")
	todayStr := today.Format("2006-01-02")
	for _, t := range treaties {
		switch t.Status {
		case "active":
			stats.Active++
			if t.ExpiryDate >= todayStr && t.ExpiryDate <= in60 {
				stats.ExpiringIn60Days++
			}
		case "draft":
			stats.Draft++
		case "expired":
			stats.Expired++
		}
	}
	return stats, nil
}

// GetRunOffTreaties returns all treaties with is_run_off = true, optionally filtered by status
func GetRunOffTreaties(status string) ([]models.ReinsuranceTreaty, error) {
	var treaties []models.ReinsuranceTreaty
	q := DB.Where("is_run_off = ?", true).Order("run_off_start_date DESC")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Find(&treaties).Error; err != nil {
		return nil, err
	}
	return treaties, nil
}
