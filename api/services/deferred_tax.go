package services

import (
	"api/models"
	"context"
	"time"

	"gorm.io/gorm"
)

func ComputeDeferredTax(req models.ComputeDeferredTaxRequest) ([]models.DeferredTaxEntry, models.DeferredTaxSummary, error) {
	ctx := context.Background()
	taxRate := req.TaxRate
	if taxRate <= 0 {
		taxRate = 0.27
	}

	var steps []models.AOSStepResult
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Where("csm_run_id = ?", req.RunID).Find(&steps).Error
	})
	if err != nil {
		return nil, models.DeferredTaxSummary{}, err
	}

	DB.Where("run_id = ?", req.RunID).Delete(&models.DeferredTaxEntry{})

	now := time.Now()
	entries := make([]models.DeferredTaxEntry, 0, len(steps))
	var summary models.DeferredTaxSummary
	summary.RunID = req.RunID
	summary.TaxRate = taxRate

	for _, s := range steps {
		csm := s.CSMBuildup - s.CSMRelease
		ra := s.RiskAdjustment
		loss := s.LossComponentBuildup

		dtlCSM := 0.0
		if csm > 0 {
			dtlCSM = csm * taxRate
		}
		dtlRA := 0.0
		if ra > 0 {
			dtlRA = ra * taxRate
		}
		dtaLoss := 0.0
		if loss > 0 {
			dtaLoss = loss * taxRate
		}
		net := (dtlCSM + dtlRA) - dtaLoss

		entries = append(entries, models.DeferredTaxEntry{
			RunID:          req.RunID,
			ProductCode:    s.ProductCode,
			IFRS17Group:    s.IFRS17Group,
			CSMAmount:      csm,
			RAAmount:       ra,
			LossComponent:  loss,
			TaxRate:        taxRate,
			DTLOnCSM:       dtlCSM,
			DTLOnRA:        dtlRA,
			DTAOnLoss:      dtaLoss,
			NetDeferredTax: net,
			ComputedAt:     now,
		})

		summary.TotalCSM += csm
		summary.TotalRA += ra
		summary.TotalLoss += loss
		summary.TotalDTL += dtlCSM + dtlRA
		summary.TotalDTA += dtaLoss
		summary.NetDeferredTax += net
		summary.GroupCount++
	}

	if len(entries) > 0 {
		if err2 := DB.Create(&entries).Error; err2 != nil {
			return nil, summary, err2
		}
	}
	return entries, summary, nil
}

func GetDeferredTaxEntries(runID int) ([]models.DeferredTaxEntry, models.DeferredTaxSummary, error) {
	ctx := context.Background()
	var entries []models.DeferredTaxEntry
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Where("run_id = ?", runID).Order("product_code, ifrs17_group").Find(&entries).Error
	})
	if err != nil {
		return nil, models.DeferredTaxSummary{}, err
	}
	var summary models.DeferredTaxSummary
	summary.RunID = runID
	for _, e := range entries {
		summary.TaxRate = e.TaxRate
		summary.TotalCSM += e.CSMAmount
		summary.TotalRA += e.RAAmount
		summary.TotalLoss += e.LossComponent
		summary.TotalDTL += e.DTLOnCSM + e.DTLOnRA
		summary.TotalDTA += e.DTAOnLoss
		summary.NetDeferredTax += e.NetDeferredTax
		summary.GroupCount++
	}
	return entries, summary, nil
}
