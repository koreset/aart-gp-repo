package services

import (
	appLog "api/log"
	"api/models"
	"fmt"
	"time"
)

// Backfill markers. v1 only scaled binder/outsource and re-ran
// recomputeFinalPremiumsAndCommission. v2 additionally re-runs
// applySchemeWideCommission so persisted Exp*CommissionAmount values are
// rewritten against the new SchemeTotalLoading() denominator. v2 detects
// whether v1 already ran on this install: if so, it SKIPS the binder/outsource
// scaling (already applied) and only re-runs the commission allocation.
// Without that detection, running v2 after v1 would double-scale the binder
// and outsource amounts.
const (
	officePremiumDenominatorBackfillMarkerV1 = "20260429120001_backfill_office_premium_denominator"
	officePremiumDenominatorBackfillMarkerV2 = "20260429120001_backfill_office_premium_denominator_v2"
)

// BackfillOfficePremiumDenominator brings every existing quote's persisted
// premium fields onto the new SchemeTotalLoading() denominator
// (expense + profit + admin + other + binder + outsource).
//
// Two things were wrong on legacy data:
//
//  1. Per-member binder/outsource amounts (and thus the summary aggregates
//     ExpTotal*BinderAmount / ExpTotal*OutsourcedAmount / Total* siblings)
//     were computed via the pre-fix schemeLoadingFromQuote() which only
//     summed expense + profit. Multiplying these by
//     scale = (1 - (expense + profit)) / (1 - SchemeTotalLoading())
//     restores office × rate = amount under the new denominator.
//
//  2. Exp*CommissionAmount was allocated by applySchemeWideCommission using
//     the pre-fix ComputeOfficePremium(risk, summary) which divided by
//     (1 - (expense + profit)) only. Re-running applySchemeWideCommission
//     under the new denominator rewrites these slices so Exp*Commission
//     and Final*Commission agree when Discount == 0 (the invariant the
//     QuoteBenefitSummary buildRow comment relies on).
//
// Idempotent: gates execution on a marker row in the migrations table.
// Detects whether v1 already scaled the binder/outsource amounts and skips
// re-scaling in that case so re-running v2 against a v1-treated database
// doesn't double-scale.
func BackfillOfficePremiumDenominator() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	var v2 Migration
	if err := DB.Where("version = ?", officePremiumDenominatorBackfillMarkerV2).First(&v2).Error; err == nil {
		appLog.Info("Office-premium denominator backfill v2 already applied; skipping")
		return nil
	}

	appLog.Info("Running office-premium denominator backfill v2")

	var v1 Migration
	v1Already := DB.Where("version = ?", officePremiumDenominatorBackfillMarkerV1).First(&v1).Error == nil

	if !v1Already {
		appLog.Info("v1 marker absent — scaling summary binder/outsource amounts")
		if err := scaleSummaryBinderOutsourceAmounts(); err != nil {
			return fmt.Errorf("scale summary binder/outsource amounts: %w", err)
		}
	} else {
		appLog.Info("v1 marker present — binder/outsource scaling already applied, skipping to commission re-allocation")
	}

	if err := refreshAllQuotesPremiumsAndCommission(); err != nil {
		return fmt.Errorf("refresh premiums and commission: %w", err)
	}

	marker := Migration{
		Version:   officePremiumDenominatorBackfillMarkerV2,
		Name:      "backfill office premium denominator v2 (admin/other/binder/outsource + commission re-allocation)",
		AppliedAt: time.Now(),
	}
	if err := DB.Create(&marker).Error; err != nil {
		return fmt.Errorf("record backfill marker: %w", err)
	}

	appLog.Info("Office-premium denominator backfill v2 completed")
	return nil
}

// scaleSummaryBinderOutsourceAmounts multiplies every persisted
// *_binder_amount / *_outsourced_amount column on member_rating_result_summaries
// by (1 - (expense_loading + profit_loading)) / (1 - SchemeTotalLoading()).
// On rows where the two denominators are equal (non-binder channels prior to
// this change, plus any rows the calc pipeline already wrote with the new
// formula), scale = 1 and the row is unchanged.
func scaleSummaryBinderOutsourceAmounts() error {
	var summaries []models.MemberRatingResultSummary
	if err := DB.Find(&summaries).Error; err != nil {
		return err
	}

	for i := range summaries {
		s := &summaries[i]
		oldDenom := s.ExpenseLoading + s.ProfitLoading
		newDenom := s.SchemeTotalLoading()
		if oldDenom == newDenom {
			continue
		}
		oldFactor := 1.0 - oldDenom
		newFactor := 1.0 - newDenom
		if oldFactor <= 0 || newFactor <= 0 {
			continue
		}
		scale := oldFactor / newFactor
		applyScaleToBinderOutsource(s, scale)
		if err := DB.Save(s).Error; err != nil {
			return fmt.Errorf("save scaled summary %d: %w", s.ID, err)
		}
	}
	return nil
}

// applyScaleToBinderOutsource multiplies every binder/outsource amount field
// (Total*, ExpTotal*, ExpAdjTotal*, Final*) on the given summary by `scale`.
// Final* are scaled here too so the row stays self-consistent before
// recomputeFinalPremiumsAndCommission overwrites them in pass 2.
func applyScaleToBinderOutsource(s *models.MemberRatingResultSummary, scale float64) {
	s.TotalGlaAnnualBinderAmount *= scale
	s.TotalGlaAnnualOutsourcedAmount *= scale
	s.ExpTotalGlaAnnualBinderAmount *= scale
	s.ExpTotalGlaAnnualOutsourcedAmount *= scale
	s.TotalAdditionalAccidentalGlaAnnualBinderAmount *= scale
	s.TotalAdditionalAccidentalGlaAnnualOutsourcedAmt *= scale
	s.ExpTotalAdditionalAccidentalGlaAnnualBinderAmount *= scale
	s.ExpTotalAdditionalAccidentalGlaAnnualOutsourcedAmt *= scale
	s.TotalPtdAnnualBinderAmount *= scale
	s.TotalPtdAnnualOutsourcedAmount *= scale
	s.ExpTotalPtdAnnualBinderAmount *= scale
	s.ExpTotalPtdAnnualOutsourcedAmount *= scale
	s.TotalCiAnnualBinderAmount *= scale
	s.TotalCiAnnualOutsourcedAmount *= scale
	s.ExpTotalCiAnnualBinderAmount *= scale
	s.ExpTotalCiAnnualOutsourcedAmount *= scale
	s.TotalSglaAnnualBinderAmount *= scale
	s.TotalSglaAnnualOutsourcedAmount *= scale
	s.ExpTotalSglaAnnualBinderAmount *= scale
	s.ExpTotalSglaAnnualOutsourcedAmount *= scale
	s.TotalTtdAnnualBinderAmount *= scale
	s.TotalTtdAnnualOutsourcedAmount *= scale
	s.ExpTotalTtdAnnualBinderAmount *= scale
	s.ExpTotalTtdAnnualOutsourcedAmount *= scale
	s.TotalPhiAnnualBinderAmount *= scale
	s.TotalPhiAnnualOutsourcedAmount *= scale
	s.ExpTotalPhiAnnualBinderAmount *= scale
	s.ExpTotalPhiAnnualOutsourcedAmount *= scale
	s.TotalFunAnnualBinderAmount *= scale
	s.TotalFunAnnualOutsourcedAmount *= scale
	s.ExpTotalFunAnnualBinderAmount *= scale
	s.ExpTotalFunAnnualOutsourcedAmount *= scale
	s.TotalGlaEducatorBinderAmount *= scale
	s.TotalGlaEducatorOutsourcedAmount *= scale
	s.ExpAdjTotalGlaEducatorBinderAmount *= scale
	s.ExpAdjTotalGlaEducatorOutsourcedAmount *= scale
	s.TotalPtdEducatorBinderAmount *= scale
	s.TotalPtdEducatorOutsourcedAmount *= scale
	s.ExpAdjTotalPtdEducatorBinderAmount *= scale
	s.ExpAdjTotalPtdEducatorOutsourcedAmount *= scale
	s.TotalAnnualBinderAmount *= scale
	s.TotalAnnualOutsourcedAmount *= scale
}

// refreshAllQuotesPremiumsAndCommission re-runs applySchemeWideCommission on
// every quote so the persisted Exp*CommissionAmount fields get rewritten
// against the new ComputeOfficePremium values. applySchemeWideCommission
// internally calls recomputeFinalPremiumsAndCommission, which refreshes
// Final*OfficePremium / Final*Commission / Final*Binder / Final*Outsource.
//
// Without re-running the Exp* commission allocation, with Discount == 0 the
// Pre side (Exp*) and Final side disagree because Exp* commission slices were
// sized against the old denominator while Final* commission slices are sized
// against the new one — the buildRow invariants break.
func refreshAllQuotesPremiumsAndCommission() error {
	var quoteIDs []int
	if err := DB.Model(&models.MemberRatingResultSummary{}).
		Distinct("quote_id").Pluck("quote_id", &quoteIDs).Error; err != nil {
		return err
	}

	logger := appLog.WithField("backfill", "office_premium_denominator")
	for _, qID := range quoteIDs {
		var quote models.GroupPricingQuote
		if err := DB.First(&quote, qID).Error; err != nil {
			appLog.WithField("quote_id", qID).
				WithField("error", err.Error()).
				Warn("Skipping quote during refresh — quote not found")
			continue
		}
		// applySchemeWideCommission rewrites Exp*Commission *and* internally
		// calls recomputeFinalPremiumsAndCommission, so one call refreshes
		// both Pre-side and Final-side commission allocations consistently.
		if err := applySchemeWideCommission(qID, quote, logger); err != nil {
			appLog.WithField("quote_id", qID).
				WithField("error", err.Error()).
				Warn("Skipping quote during refresh — commission re-allocation failed")
			continue
		}
	}
	return nil
}
