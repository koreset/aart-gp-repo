package services

import (
	"api/models"
	"context"
	"gorm.io/gorm"
)

// GetSARBMappings returns all SARB code mappings.
func GetSARBMappings() ([]models.SARBCodeMapping, error) {
	var mappings []models.SARBCodeMapping
	ctx := context.Background()
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Order("regulatory_return, return_code").Find(&mappings).Error
	})
	return mappings, err
}

// UpsertSARBMapping creates or updates a SARB code mapping.
func UpsertSARBMapping(req models.UpsertSARBMappingRequest, user models.AppUser) (models.SARBCodeMapping, error) {
	mapping := models.SARBCodeMapping{
		ID:                     req.ID,
		LineItem:               req.LineItem,
		Description:            req.Description,
		IFRS17BalanceSheetItem: req.IFRS17BalanceSheetItem,
		RegulatoryReturn:       req.RegulatoryReturn,
		ReturnCode:             req.ReturnCode,
		DebitCreditIndicator:   req.DebitCreditIndicator,
		Notes:                  req.Notes,
		CreatedBy:              user.UserName,
	}
	err := DB.Save(&mapping).Error
	return mapping, err
}

// DeleteSARBMapping deletes a mapping by ID.
func DeleteSARBMapping(id int) error {
	return DB.Delete(&models.SARBCodeMapping{}, id).Error
}

// GenerateSARBReport produces a regulatory return report for a given run.
// It reads BalanceSheet/TrialBalance amounts and maps them to return codes.
// Since we don't have a separate balance sheet model, we aggregate from AOSStepResult
// and PAA results based on IFRS17BalanceSheetItem keyword matching.
func GenerateSARBReport(runID int) ([]models.SARBReportRow, error) {
	ctx := context.Background()

	// Fetch mappings
	var mappings []models.SARBCodeMapping
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Order("regulatory_return, return_code").Find(&mappings).Error
	})
	if err != nil {
		return nil, err
	}

	// Aggregate key IFRS17 amounts from AOSStepResult for this run
	type Aggregates struct {
		TotalBEL          float64
		TotalRA           float64
		TotalCSM          float64
		TotalLRC          float64
		TotalReinsBEL     float64
		TotalReinsRA      float64
		TotalReinsCSM     float64
		TotalInsRevenue   float64
		TotalReinsRevenue float64
	}
	var agg Aggregates
	var steps []models.AOSStepResult
	err = DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Where("csm_run_id = ?", runID).Find(&steps).Error
	})
	if err != nil {
		return nil, err
	}
	for _, s := range steps {
		agg.TotalBEL += s.BEL
		agg.TotalRA += s.RiskAdjustment
		agg.TotalCSM += s.CSMBuildup - s.CSMRelease
		agg.TotalReinsBEL += s.ReinsuranceBel
		agg.TotalReinsRA += s.ReinsuranceRiskAdjustment
		agg.TotalReinsCSM += s.ReinsuranceCSM
		agg.TotalReinsRevenue += s.ReinsuranceRevenue
	}
	agg.TotalLRC = agg.TotalBEL + agg.TotalRA + agg.TotalCSM

	// Map IFRS17BalanceSheetItem keywords to computed amounts
	amountFor := func(item string) float64 {
		switch item {
		case "BEL":
			return agg.TotalBEL
		case "RA":
			return agg.TotalRA
		case "CSM":
			return agg.TotalCSM
		case "LRC":
			return agg.TotalLRC
		case "Reinsurance BEL":
			return agg.TotalReinsBEL
		case "Reinsurance RA":
			return agg.TotalReinsRA
		case "Reinsurance CSM":
			return agg.TotalReinsCSM
		case "Reinsurance Revenue":
			return agg.TotalReinsRevenue
		default:
			return 0
		}
	}

	// Build report rows
	rows := make([]models.SARBReportRow, 0, len(mappings))
	for _, m := range mappings {
		rows = append(rows, models.SARBReportRow{
			ReturnCode:           m.ReturnCode,
			RegulatoryReturn:     m.RegulatoryReturn,
			LineItem:             m.LineItem,
			Description:          m.Description,
			DebitCreditIndicator: m.DebitCreditIndicator,
			Amount:               amountFor(m.IFRS17BalanceSheetItem),
			Notes:                m.Notes,
		})
	}
	return rows, nil
}

// SeedDefaultSARBMappings inserts a standard set of RC&AP and J200 mappings if none exist.
func SeedDefaultSARBMappings(user models.AppUser) error {
	var count int64
	DB.Model(&models.SARBCodeMapping{}).Count(&count)
	if count > 0 {
		return nil
	}
	defaults := []models.SARBCodeMapping{
		{LineItem: "Insurance Contract Liabilities — BEL", IFRS17BalanceSheetItem: "BEL", RegulatoryReturn: "RC&AP", ReturnCode: "RC01", DebitCreditIndicator: "C", Description: "Best Estimate Liability on insurance contracts", CreatedBy: user.UserName},
		{LineItem: "Insurance Contract Liabilities — RA", IFRS17BalanceSheetItem: "RA", RegulatoryReturn: "RC&AP", ReturnCode: "RC02", DebitCreditIndicator: "C", Description: "Risk Adjustment for non-financial risk", CreatedBy: user.UserName},
		{LineItem: "Insurance Contract Liabilities — CSM", IFRS17BalanceSheetItem: "CSM", RegulatoryReturn: "RC&AP", ReturnCode: "RC03", DebitCreditIndicator: "C", Description: "Contractual Service Margin", CreatedBy: user.UserName},
		{LineItem: "Insurance Contract Liabilities — Total LRC", IFRS17BalanceSheetItem: "LRC", RegulatoryReturn: "RC&AP", ReturnCode: "RC04", DebitCreditIndicator: "C", Description: "Total Liability for Remaining Coverage", CreatedBy: user.UserName},
		{LineItem: "Reinsurance Contract Assets — BEL", IFRS17BalanceSheetItem: "Reinsurance BEL", RegulatoryReturn: "RC&AP", ReturnCode: "RC10", DebitCreditIndicator: "D", Description: "Reinsurance best estimate asset", CreatedBy: user.UserName},
		{LineItem: "Reinsurance Contract Assets — RA", IFRS17BalanceSheetItem: "Reinsurance RA", RegulatoryReturn: "RC&AP", ReturnCode: "RC11", DebitCreditIndicator: "D", Description: "Reinsurance risk adjustment asset", CreatedBy: user.UserName},
		{LineItem: "Reinsurance Contract Assets — CSM", IFRS17BalanceSheetItem: "Reinsurance CSM", RegulatoryReturn: "RC&AP", ReturnCode: "RC12", DebitCreditIndicator: "D", Description: "Reinsurance CSM asset", CreatedBy: user.UserName},
		{LineItem: "Insurance Revenue — Reinsurance", IFRS17BalanceSheetItem: "Reinsurance Revenue", RegulatoryReturn: "J200", ReturnCode: "J200-I01", DebitCreditIndicator: "C", Description: "Reinsurance income", CreatedBy: user.UserName},
		{LineItem: "Balance Sheet — LRC (J200)", IFRS17BalanceSheetItem: "LRC", RegulatoryReturn: "J200", ReturnCode: "J200-L01", DebitCreditIndicator: "C", Description: "Insurance contract liabilities (J200 form)", CreatedBy: user.UserName},
		{LineItem: "Balance Sheet — Reins BEL (J200)", IFRS17BalanceSheetItem: "Reinsurance BEL", RegulatoryReturn: "J200", ReturnCode: "J200-A01", DebitCreditIndicator: "D", Description: "Reinsurance contract assets (J200 form)", CreatedBy: user.UserName},
	}
	return DB.Create(&defaults).Error
}
