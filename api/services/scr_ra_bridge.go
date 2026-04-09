package services

import (
	"api/models"
	"context"

	"gorm.io/gorm"
)

// CreateSCREntry saves a manually entered SCR risk margin record.
func CreateSCREntry(req models.CreateSCREntryRequest, user models.AppUser) (models.SCRRABridgeEntry, error) {
	entry := models.SCRRABridgeEntry{
		ProductCode:   req.ProductCode,
		Period:        req.Period,
		SCRRiskMargin: req.SCRRiskMargin,
		Notes:         req.Notes,
		CreatedBy:     user.UserName,
	}
	err := DB.Save(&entry).Error
	return entry, err
}

// GetSCREntries returns all SCR entries, optionally filtered by product.
func GetSCREntries(productCode string) ([]models.SCRRABridgeEntry, error) {
	var entries []models.SCRRABridgeEntry
	ctx := context.Background()
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		q := d.Order("period DESC")
		if productCode != "" {
			q = q.Where("product_code = ?", productCode)
		}
		return q.Find(&entries).Error
	})
	return entries, err
}

// DeleteSCREntry deletes a SCR entry by ID.
func DeleteSCREntry(id int) error {
	return DB.Delete(&models.SCRRABridgeEntry{}, id).Error
}

// GetSCRRABridge builds the reconciliation report for a given CSM run.
// It aggregates IFRS 17 RA from AOSStepResult by product, then joins with SCR entries.
func GetSCRRABridge(runID int) ([]models.SCRRABridgeRow, error) {
	ctx := context.Background()

	// Fetch the run to get run_date / period
	var run models.CsmRun
	err := DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.First(&run, runID).Error
	})
	if err != nil {
		return nil, err
	}

	// Aggregate RA by product for this run
	type raAgg struct {
		ProductCode string
		TotalRA     float64
	}
	var agg []raAgg
	err = DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Model(&models.AOSStepResult{}).
			Select("product_code, SUM(risk_adjustment) as total_ra").
			Where("csm_run_id = ?", runID).
			Group("product_code").
			Scan(&agg).Error
	})
	if err != nil {
		return nil, err
	}

	// Fetch SCR entries for this period
	var scrEntries []models.SCRRABridgeEntry
	err = DBReadWithResilience(ctx, func(d *gorm.DB) error {
		return d.Where("period = ?", run.RunDate).Find(&scrEntries).Error
	})
	if err != nil {
		return nil, err
	}

	// Build SCR map
	scrMap := map[string]models.SCRRABridgeEntry{}
	for _, e := range scrEntries {
		scrMap[e.ProductCode] = e
	}

	// Build reconciliation rows
	rows := make([]models.SCRRABridgeRow, 0, len(agg))
	for _, a := range agg {
		scr := scrMap[a.ProductCode]
		variance := a.TotalRA - scr.SCRRiskMargin
		var variancePct float64
		if scr.SCRRiskMargin != 0 {
			variancePct = variance / scr.SCRRiskMargin * 100
		}
		rows = append(rows, models.SCRRABridgeRow{
			ProductCode:   a.ProductCode,
			Period:        run.RunDate,
			IFRS17RA:      a.TotalRA,
			SCRRiskMargin: scr.SCRRiskMargin,
			Variance:      variance,
			VariancePct:   variancePct,
			Notes:         scr.Notes,
		})
	}
	return rows, nil
}
