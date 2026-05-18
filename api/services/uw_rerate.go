package services

import (
	"fmt"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"

	"api/models"
)

// benefitKey is the canonical lowercase identifier used by UnderwritingDecision
// rows: gla, ptd, ci, sgla. Other benefit codes (educator, scb, etc.) are
// ignored in Phase 4 — they're documentary on the case file but not folded
// into the UW-adjusted totals.
type benefitKey string

const (
	benefitGla  benefitKey = "gla"
	benefitPtd  benefitKey = "ptd"
	benefitCi   benefitKey = "ci"
	benefitSgla benefitKey = "sgla"
)

// memberDecisionMap is `(MemberName + Category) -> benefit -> Decision`.
// The latest decision per (member, benefit) wins.
type memberDecisionMap map[string]map[benefitKey]models.UnderwritingDecision

// ApplyDecisionsAndReRate rebuilds the UW-adjusted premium for a quote by
// walking its MemberRatingResult rows, folding in every UnderwritingDecision,
// summing into the category summaries, and writing a QuoteReRateEvent.
// Pushes a WSQuoteReRated envelope to the broker (quote creator) and the
// triggering user when finished.
//
// `reason` is recorded on the event and shown in the renderer banner.
// `caseID` is optional context — non-zero when the trigger was a specific
// decision write.
//
// Returns the resulting QuoteReRateEvent (with previous/new premium delta).
func ApplyDecisionsAndReRate(quoteID int, user models.AppUser, reason string, caseID int) (*models.QuoteReRateEvent, error) {
	if quoteID <= 0 {
		return nil, fmt.Errorf("invalid quote id")
	}

	var quote models.GroupPricingQuote
	if err := DB.Where("id = ?", quoteID).First(&quote).Error; err != nil {
		return nil, fmt.Errorf("load quote: %w", err)
	}

	decisions, err := loadLatestDecisionsForQuote(quoteID)
	if err != nil {
		return nil, err
	}

	// Load summaries (one per category) and rating rows. The rating set is
	// the source of truth for per-member rates and original capped SAs.
	var summaries []models.MemberRatingResultSummary
	if err := DB.Where("quote_id = ?", quoteID).Find(&summaries).Error; err != nil {
		return nil, fmt.Errorf("load summaries: %w", err)
	}
	var rows []models.MemberRatingResult
	if err := DB.Where("quote_id = ?", quoteID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("load member rows: %w", err)
	}

	previousPremium := sumExistingTotalPremium(summaries)

	summaryByCategory := map[string]*models.MemberRatingResultSummary{}
	for i := range summaries {
		summaryByCategory[summaries[i].Category] = &summaries[i]
		summaries[i].UWAdjustedTotalAnnualPremium = 0
		summaries[i].UWAdjustedTotalGlaCappedSumAssured = 0
		summaries[i].UWAdjustedTotalPtdCappedSumAssured = 0
		summaries[i].UWAdjustedTotalCiCappedSumAssured = 0
		summaries[i].UWAdjustedTotalSglaCappedSumAssured = 0
	}

	now := time.Now()

	if err := DB.Transaction(func(tx *gorm.DB) error {
		for i := range rows {
			row := &rows[i]
			summary := summaryByCategory[row.Category]
			if summary == nil {
				continue
			}
			memberDecs := decisions[memberDecisionKey(row.MemberName, row.Category)]
			applyDecisionsToRow(row, memberDecs)
			adjustedRiskPremium := computeUWAdjustedRiskPremium(*row)
			row.UWAdjustedAnnualOfficePremium = models.ComputeOfficePremium(adjustedRiskPremium, summary)

			summary.UWAdjustedTotalAnnualPremium += row.UWAdjustedAnnualOfficePremium
			summary.UWAdjustedTotalGlaCappedSumAssured += effectiveSA(row.GlaCappedSumAssured, row.UWGlaCoverCap, row.UWGlaDeclined)
			summary.UWAdjustedTotalPtdCappedSumAssured += effectiveSA(row.PtdCappedSumAssured, row.UWPtdCoverCap, row.UWPtdDeclined)
			summary.UWAdjustedTotalCiCappedSumAssured += effectiveSA(row.CiCappedSumAssured, row.UWCiCoverCap, row.UWCiDeclined)
			summary.UWAdjustedTotalSglaCappedSumAssured += effectiveSA(row.SpouseGlaCappedSumAssured, row.UWSglaCoverCap, row.UWSglaDeclined)

			if err := tx.Save(row).Error; err != nil {
				return fmt.Errorf("save member row %d: %w", row.ID, err)
			}
		}
		for i := range summaries {
			summaries[i].UWReRatedAt = &now
			summaries[i].UWReRatedBy = user.UserEmail
			if err := tx.Save(&summaries[i]).Error; err != nil {
				return fmt.Errorf("save summary: %w", err)
			}
		}
		quote.RatingVersion++
		if err := tx.Model(&models.GroupPricingQuote{}).
			Where("id = ?", quote.ID).
			Update("rating_version", quote.RatingVersion).Error; err != nil {
			return fmt.Errorf("bump rating version: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	newPremium := 0.0
	for i := range summaries {
		newPremium += summaries[i].UWAdjustedTotalAnnualPremium
	}

	event := models.QuoteReRateEvent{
		QuoteID:         quoteID,
		TriggeredBy:     user.UserEmail,
		PreviousPremium: previousPremium,
		NewPremium:      newPremium,
		PremiumDelta:    newPremium - previousPremium,
		Reason:          reason,
		CaseID:          caseID,
		RatingVersion:   quote.RatingVersion,
	}
	if err := DB.Create(&event).Error; err != nil {
		return nil, fmt.Errorf("write event: %w", err)
	}

	pushQuoteReRatedEvent(quote, event, user)
	return &event, nil
}

// loadLatestDecisionsForQuote builds memberDecisionMap, keeping only the
// most-recent UnderwritingDecision per (case, benefit). Decisions on cases
// that belong to a different quote are skipped by the join.
func loadLatestDecisionsForQuote(quoteID int) (memberDecisionMap, error) {
	var cases []models.UnderwritingCase
	if err := DB.Preload("Decisions").
		Where("quote_id = ?", quoteID).
		Find(&cases).Error; err != nil {
		return nil, fmt.Errorf("load cases: %w", err)
	}
	out := memberDecisionMap{}
	for _, c := range cases {
		key := memberDecisionKey(c.MemberName, c.Category)
		if _, ok := out[key]; !ok {
			out[key] = map[benefitKey]models.UnderwritingDecision{}
		}
		// Decisions are append-only; keep the newest per benefit.
		latest := map[benefitKey]models.UnderwritingDecision{}
		for _, d := range c.Decisions {
			bk := benefitKey(strings.ToLower(d.BenefitType))
			cur, ok := latest[bk]
			if !ok || d.CreationDate.After(cur.CreationDate) || (d.CreationDate.Equal(cur.CreationDate) && d.ID > cur.ID) {
				latest[bk] = d
			}
		}
		for bk, d := range latest {
			out[key][bk] = d
		}
	}
	return out, nil
}

func memberDecisionKey(name, category string) string {
	return name + "|" + category
}

// applyDecisionsToRow folds the per-benefit decision adjustments into the
// rating row's UW* fields. Outcomes:
//   - accept   → loading + cap applied
//   - postpone → treated as decline for premium purposes (no cover until
//                evidence comes back; safer than charging an in-doubt premium)
//   - decline  → declined flag set; SA contribution is zero
func applyDecisionsToRow(row *models.MemberRatingResult, dec map[benefitKey]models.UnderwritingDecision) {
	if dec == nil {
		// Reset UW fields so a previously-decided case that's been deleted
		// no longer affects the row.
		row.UWGlaLoading, row.UWGlaCoverCap, row.UWGlaDeclined = 0, 0, false
		row.UWPtdLoading, row.UWPtdCoverCap, row.UWPtdDeclined = 0, 0, false
		row.UWCiLoading, row.UWCiCoverCap, row.UWCiDeclined = 0, 0, false
		row.UWSglaLoading, row.UWSglaCoverCap, row.UWSglaDeclined = 0, 0, false
		return
	}
	apply := func(d models.UnderwritingDecision, loading, cap *float64, declined *bool) {
		declinedOutcome := d.Outcome == models.UWOutcomeDecline || d.Outcome == models.UWOutcomePostpone
		*declined = declinedOutcome
		if declinedOutcome {
			*loading = 0
			*cap = 0
			return
		}
		*loading = d.LoadingPercent
		*cap = d.CoverCap
	}
	if d, ok := dec[benefitGla]; ok {
		apply(d, &row.UWGlaLoading, &row.UWGlaCoverCap, &row.UWGlaDeclined)
	} else {
		row.UWGlaLoading, row.UWGlaCoverCap, row.UWGlaDeclined = 0, 0, false
	}
	if d, ok := dec[benefitPtd]; ok {
		apply(d, &row.UWPtdLoading, &row.UWPtdCoverCap, &row.UWPtdDeclined)
	} else {
		row.UWPtdLoading, row.UWPtdCoverCap, row.UWPtdDeclined = 0, 0, false
	}
	if d, ok := dec[benefitCi]; ok {
		apply(d, &row.UWCiLoading, &row.UWCiCoverCap, &row.UWCiDeclined)
	} else {
		row.UWCiLoading, row.UWCiCoverCap, row.UWCiDeclined = 0, 0, false
	}
	if d, ok := dec[benefitSgla]; ok {
		apply(d, &row.UWSglaLoading, &row.UWSglaCoverCap, &row.UWSglaDeclined)
	} else {
		row.UWSglaLoading, row.UWSglaCoverCap, row.UWSglaDeclined = 0, 0, false
	}
}

// computeUWAdjustedRiskPremium re-derives the member's pre-loading risk
// premium across the 4 underwritten benefits with UW adjustments applied.
// Returns the sum of per-benefit `ExpAdjLoadedRate × (1+loading) × effectiveSA`.
func computeUWAdjustedRiskPremium(row models.MemberRatingResult) float64 {
	gla := row.ExpAdjLoadedGlaRate * (1 + row.UWGlaLoading/100) * effectiveSA(row.GlaCappedSumAssured, row.UWGlaCoverCap, row.UWGlaDeclined)
	ptd := row.ExpAdjLoadedPtdRate * (1 + row.UWPtdLoading/100) * effectiveSA(row.PtdCappedSumAssured, row.UWPtdCoverCap, row.UWPtdDeclined)
	ci := row.ExpAdjLoadedCiRate * (1 + row.UWCiLoading/100) * effectiveSA(row.CiCappedSumAssured, row.UWCiCoverCap, row.UWCiDeclined)
	sgla := row.ExpAdjLoadedSpouseGlaRate * (1 + row.UWSglaLoading/100) * effectiveSA(row.SpouseGlaCappedSumAssured, row.UWSglaCoverCap, row.UWSglaDeclined)
	return gla + ptd + ci + sgla
}

// effectiveSA returns the SA used in UW-adjusted premium math: 0 if the
// benefit is declined; otherwise min(originalCap, uwCap) where uwCap is
// ignored when 0.
func effectiveSA(originalCap, uwCap float64, declined bool) float64 {
	if declined {
		return 0
	}
	if uwCap > 0 {
		return math.Min(originalCap, uwCap)
	}
	return originalCap
}

func sumExistingTotalPremium(summaries []models.MemberRatingResultSummary) float64 {
	total := 0.0
	for _, s := range summaries {
		if s.UWAdjustedTotalAnnualPremium > 0 {
			total += s.UWAdjustedTotalAnnualPremium
			continue
		}
		total += s.TotalAnnualPremium
	}
	return total
}

func pushQuoteReRatedEvent(quote models.GroupPricingQuote, event models.QuoteReRateEvent, user models.AppUser) {
	hub := GetHub()
	if hub == nil {
		return
	}
	payload := QuoteReRatedPayload{
		QuoteID:         quote.ID,
		PreviousPremium: event.PreviousPremium,
		NewPremium:      event.NewPremium,
		Delta:           event.PremiumDelta,
		Reason:          event.Reason,
		CaseID:          event.CaseID,
		RatingVersion:   event.RatingVersion,
		TriggeredBy:     event.TriggeredBy,
		TriggeredAt:     event.TriggeredAt.Format(time.RFC3339),
	}
	env := WSEnvelope{Type: WSQuoteReRated, Payload: payload}
	// Push to the broker on the quote (CreatedBy is a name not an email, so
	// fall back to ModifiedBy if it looks like an email) and the triggering
	// user. De-dup if they coincide.
	recipients := map[string]struct{}{user.UserEmail: {}}
	if strings.Contains(quote.ModifiedBy, "@") {
		recipients[quote.ModifiedBy] = struct{}{}
	}
	if strings.Contains(quote.CreatedBy, "@") {
		recipients[quote.CreatedBy] = struct{}{}
	}
	for email := range recipients {
		if email == "" {
			continue
		}
		hub.SendToUser(email, env)
	}
}
