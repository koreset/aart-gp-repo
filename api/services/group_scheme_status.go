package services

import (
	"api/models"
	"api/utils"
	"fmt"
	"math"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func UpdateGroupSchemeStatuses() {
	var schemes []models.GroupScheme
	now := time.Now()

	// achieve the following
	// 1. Check for schemes whose cover end date is equal to or greater than today and status is InForce
	// 		a. Check if there is a renewal quote available, if it is, update the cover start and end date of the scheme and set the status to InForce
	// 		b. If there is no renewal quote available, set the status to Expired
	// 2. Check for schemes whose status is accepted and cover_start_date is <= today.
	// 	a. Set the status to InForce. Set cover start and end date as per the quote
	//
	// 3. For all InForce schemes, update earned premium (no audit)

	// 1) Check for schemes whose cover end date is greater than or equal to today

	if err := DB.Where("cover_end_date <= ? AND status = ?", now, models.StatusInForce).Find(&schemes).Error; err != nil {
		return
	}
	for _, s := range schemes {
		fmt.Printf("Group Scheme ID %d has expired\n", s.ID)
		var quote models.GroupPricingQuote
		var quoteStats models.GroupRiskQuoteStats
		DB.Where("scheme_id = ? and status = ? and scheme_quote_status = ?", s.ID, models.StatusAccepted, models.StatusNotInEffect).Order("commencement_date desc").First(&quote)

		if quote.ID > 0 {
			// expire the current quote and set the new quote as in force
			DB.Model(&models.GroupPricingQuote{}).Where("quote_name = ?", s.QuoteInForce).Update("scheme_quote_status", models.StatusExpired)
			fmt.Println(quote)
			//read quote stats and assign to current quote stats to scheme's stats
			DB.Where("quote_id  = ?", quote.ID).First(&quoteStats)

			if err := DB.Transaction(func(tx *gorm.DB) error {
				var before models.GroupScheme
				if err := tx.Where("id = ?", s.ID).First(&before).Error; err != nil {
					return err
				}
				updates := map[string]interface{}{
					"status":                   models.StatusInForce,
					"scheme_quote_status":      models.StatusInEffect,
					"cover_start_date":         quote.CommencementDate,
					"cover_end_date":           quote.CoverEndDate,
					"quote_in_force":           quote.QuoteName,
					"annual_premium":           quoteStats.AnnualPremium,
					"commission":               quoteStats.Commission,
					"member_count":             quoteStats.MemberCount,
					"expected_claims":          quoteStats.ExpectedClaims,
					"actual_claims":            0,
					"expected_expenses":        quoteStats.ExpectedExpenses,
					"expected_gla_claims":      quoteStats.ExpectedGlaClaims,
					"expected_ptd_claims":      quoteStats.ExpectedPtdClaims,
					"expected_ci_claims":       quoteStats.ExpectedCiClaims,
					"expected_sgla_claims":     quoteStats.ExpectedSglaClaims,
					"expected_ttd_claims":      quoteStats.ExpectedTtdClaims,
					"expected_phi_claims":      quoteStats.ExpectedPhiClaims,
					"expected_fun_claims":      quoteStats.ExpectedFunClaims,
					"gla_annual_premium":       quoteStats.GlaAnnualPremium,
					"ptd_annual_premium":       quoteStats.PtdAnnualPremium,
					"ci_annual_premium":        quoteStats.CiAnnualPremium,
					"sgla_annual_premium":      quoteStats.SglaAnnualPremium,
					"ttd_annual_premium":       quoteStats.TtdAnnualPremium,
					"phi_annual_premium":       quoteStats.PhiAnnualPremium,
					"funeral_annual_premium":   quoteStats.FuneralAnnualPremium,
					"active_scheme_categories": quote.SelectedSchemeCategories,
				}

				if err := tx.Model(&models.GroupScheme{}).Where("id = ?", s.ID).Updates(updates).Error; err != nil {
					return err
				}

				var after models.GroupScheme
				if err := tx.Where("id = ?", s.ID).First(&after).Error; err != nil {
					return err
				}

				statusAudit := models.GroupSchemeStatusAudit{
					SchemeID:      before.ID,
					OldStatus:     before.Status,
					NewStatus:     after.Status,
					StatusMessage: "Set in force by scheduler",
					ChangedBy:     "system",
					ChangedAt:     time.Now(),
				}
				if err := tx.Create(&statusAudit).Error; err != nil {
					return err
				}

				//// Generic audit
				if err := writeAudit(tx, AuditContext{
					Area:      "group-pricing",
					Entity:    "group_schemes",
					EntityID:  strconv.Itoa(s.ID),
					Action:    "UPDATE",
					ChangedBy: "system",
				}, before, after); err != nil {
					return err
				}

				if quote.QuoteType == "New Business" {
					var members []models.GPricingMemberDataInForce
					if err := tx.Where("scheme_id = ?", s.ID).Find(&members).Error; err != nil {
						return err
					}

					for _, m := range members {
						activity := models.MemberActivity{
							MemberID:       m.ID,
							MemberIDNumber: m.MemberIdNumber,
							Timestamp:      time.Now(),
							Type:           "enrollment",
							Title:          "Member Enrollment",
							Description:    "Member enrolled in scheme",
							PerformedBy:    "system",
						}
						if err := tx.Create(&activity).Error; err != nil {
							return err
						}
					}
				}

				return nil
			}); err != nil {
				fmt.Printf("[scheduler] renewal transition failed for scheme %d: %v\n", s.ID, err)
				continue
			}

			// update quote scheme quote status to In Effect
			DB.Model(&models.GroupPricingQuote{}).Where("id = ?", quote.ID).Update("scheme_quote_status", models.StatusInEffect)
			fmt.Printf("Group Scheme ID %d auto-renewed with Quote ID %d\n", s.ID, quote.ID)
		} else {
			// no renewal quote available, set status to Expired
			if err := DB.Transaction(func(tx *gorm.DB) error {
				var before models.GroupScheme
				if err := tx.Where("id = ?", s.ID).First(&before).Error; err != nil {
					return err
				}
				updates := map[string]interface{}{
					"status": models.StatusExpired,
				}
				if err := tx.Model(&models.GroupScheme{}).Where("id = ?", s.ID).Updates(updates).Error; err != nil {
					return err
				}
				var after models.GroupScheme

				if err := tx.Where("id = ?", s.ID).First(&after).Error; err != nil {
					return err
				}

				// set the associated quote to status Expired
				DB.Model(&models.GroupPricingQuote{}).Where("quote_name = ?", s.QuoteInForce).Update("scheme_quote_status", models.StatusExpired)
				return nil
			}); err != nil {
				fmt.Printf("[scheduler] expiry update failed for scheme %d: %v\n", s.ID, err)
			}
		}
	}

	// 2) Move Accepted -> InForce when cover_start_date <= now
	if err := DB.Where("status = ? AND cover_start_date <= ?", models.StatusAccepted, now).Find(&schemes).Error; err != nil {
		return
	}
	for _, s := range schemes {
		var quote models.GroupPricingQuote
		var quoteStats models.GroupRiskQuoteStats
		DB.Preload("SchemeCategories").Where("scheme_id = ? and status = ? and scheme_quote_status = ?", s.ID, models.StatusAccepted, models.StatusNotInEffect).Order("commencement_date desc").First(&quote)

		if quote.ID == 0 {
			fmt.Printf("[scheduler] no accepted quote found for scheme %d, skipping\n", s.ID)
			continue
		}

		//read quote stats and assign to current quote stats to scheme's stats
		DB.Where("quote_id  = ?", quote.ID).First(&quoteStats)

		if err := DB.Transaction(func(tx *gorm.DB) error {
			var before models.GroupScheme
			if err := tx.Where("id = ?", s.ID).First(&before).Error; err != nil {
				return err
			}

			if quote.CoverEndDate.IsZero() || quote.CommencementDate.IsZero() {
				quote.CoverEndDate = before.CoverEndDate
				quote.CommencementDate = before.CommencementDate
			}
			updates := map[string]interface{}{
				"status":                 models.StatusInForce,
				"scheme_quote_status":    models.StatusInEffect,
				"cover_start_date":       quote.CommencementDate,
				"cover_end_date":         quote.CoverEndDate,
				"quote_in_force":         quote.QuoteName,
				"annual_premium":         quoteStats.AnnualPremium,
				"commission":             quoteStats.Commission,
				"member_count":           quoteStats.MemberCount,
				"expected_claims":        quoteStats.ExpectedClaims,
				"actual_claims":          0,
				"expected_expenses":      quoteStats.ExpectedExpenses,
				"expected_gla_claims":    quoteStats.ExpectedGlaClaims,
				"expected_ptd_claims":    quoteStats.ExpectedPtdClaims,
				"expected_ci_claims":     quoteStats.ExpectedCiClaims,
				"expected_sgla_claims":   quoteStats.ExpectedSglaClaims,
				"expected_ttd_claims":    quoteStats.ExpectedTtdClaims,
				"expected_phi_claims":    quoteStats.ExpectedPhiClaims,
				"expected_fun_claims":    quoteStats.ExpectedFunClaims,
				"gla_annual_premium":     quoteStats.GlaAnnualPremium,
				"ptd_annual_premium":     quoteStats.PtdAnnualPremium,
				"ci_annual_premium":      quoteStats.CiAnnualPremium,
				"sgla_annual_premium":    quoteStats.SglaAnnualPremium,
				"ttd_annual_premium":     quoteStats.TtdAnnualPremium,
				"phi_annual_premium":     quoteStats.PhiAnnualPremium,
				"funeral_annual_premium": quoteStats.FuneralAnnualPremium,
				"expected_claims_ratio":  quoteStats.ExpectedClaimsRatio,
			}
			if err := tx.Model(&models.GroupScheme{}).Where("id = ?", s.ID).Updates(updates).Error; err != nil {
				return err
			}

			var after models.GroupScheme
			if err := tx.Where("id = ?", s.ID).First(&after).Error; err != nil {
				return err
			}

			// Specialized status audit
			statusAudit := models.GroupSchemeStatusAudit{
				SchemeID:      before.ID,
				OldStatus:     before.Status,
				NewStatus:     after.Status,
				StatusMessage: "Auto-activated (in force) by scheduler",
				ChangedBy:     "system",
				ChangedAt:     time.Now(),
			}
			if err := tx.Create(&statusAudit).Error; err != nil {
				return err
			}

			//// Generic audit
			//if err := writeAudit(tx, AuditContext{
			//	Area:      "group-pricing",
			//	Entity:    "group_schemes",
			//	EntityID:  strconv.Itoa(s.ID),
			//	Action:    "UPDATE",
			//	ChangedBy: "system",
			//}, before, after); err != nil {
			//	return err
			//}

			if quote.QuoteType == "New Business" {
				var members []models.GPricingMemberDataInForce
				if err := tx.Where("scheme_id = ?", s.ID).Find(&members).Error; err != nil {
					return err
				}

				for _, m := range members {
					activity := models.MemberActivity{
						MemberID:       m.ID,
						MemberIDNumber: m.MemberIdNumber,
						Timestamp:      time.Now(),
						Type:           "enrollment",
						Title:          "Member Enrollment",
						Description:    "Member enrolled in scheme",
						PerformedBy:    "system",
					}
					if err := tx.Create(&activity).Error; err != nil {
						return err
					}
				}
			}
			return nil
		}); err != nil {
			fmt.Printf("[scheduler] in-force transition failed for scheme %d: %v\n", s.ID, err)
			continue
		}

		// update quote scheme quote status to In Effect
		DB.Model(&models.GroupPricingQuote{}).Where("id = ?", quote.ID).Update("scheme_quote_status", models.StatusInEffect)
	}

	// 3) For all InForce schemes, update earned premium (no audit)
	if err := DB.Where("status = ?", models.StatusInForce).Find(&schemes).Error; err != nil {
		return
	}
	for _, scheme := range schemes {
		tempterm := scheme.CoverEndDate.Sub(scheme.CoverStartDate).Hours() / 24
		currentPeriodDuration := now.Sub(scheme.CoverStartDate).Hours() / 24
		scheme.RenewalDate = scheme.CoverEndDate.AddDate(0, 0, 1)
		if tempterm > 0 {
			tempEarnedProportion := math.Max(float64(currentPeriodDuration)/float64(tempterm), 0)
			scheme.DurationInForceDays = int(math.Max(float64(currentPeriodDuration), 0))
			scheme.EarnedPremium = utils.FloatPrecision(scheme.AnnualPremium*tempEarnedProportion, AccountingPrecision)
			DB.Save(&scheme)
		}
	}
}

func StartGroupSchemeStatusUpdater() {
	go func() {
		ticker := time.NewTicker(60 * time.Minute)
		defer ticker.Stop()
		for {
			fmt.Println("Updating group scheme statuses...")
			UpdateGroupSchemeStatuses()
			<-ticker.C
		}
	}()
}
