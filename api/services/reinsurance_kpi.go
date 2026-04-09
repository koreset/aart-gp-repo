package services

import (
	"api/models"
	"math"
	"time"
)

// GetRIBordereauxKPIs computes all 7 §8.2 KPIs for an optional treaty and date window.
// All parameters are optional: zero treatyID means all treaties; empty period strings
// default to the current calendar year.
func GetRIBordereauxKPIs(treatyID int, periodFrom, periodTo string) (models.RIBordereauxKPIs, error) {
	now := time.Now()

	// Default period: current calendar year
	if periodFrom == "" {
		periodFrom = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	}
	if periodTo == "" {
		periodTo = now.Format("2006-01-02")
	}

	out := models.RIBordereauxKPIs{
		TreatyID:   treatyID,
		PeriodFrom: periodFrom,
		PeriodTo:   periodTo,
		ComputedAt: now.Format("2006-01-02T15:04:05Z"),
		// Targets from §8.2
		SubmissionTimelinessTarget: 90,
		ProcessingTimelinessTarget: 95,
		FirstTimeAcceptanceTarget:  85,
		AvgErrorResolutionTarget:   5,
		SettlementTimelinessTarget: 98,
		ClaimsCompletenessTarget:   100,
		OpenQueryBacklogTarget:     5,
	}

	// ── Load all runs in the window ───────────────────────────────────────────
	var allRuns []models.RIBordereauxRun
	q := DB.Where("period_start >= ? AND period_start <= ?", periodFrom, periodTo)
	if treatyID > 0 {
		q = q.Where("treaty_id = ?", treatyID)
	}
	if err := q.Find(&allRuns).Error; err != nil {
		return out, err
	}

	// ── KPI 1: Submission timeliness ──────────────────────────────────────────
	// % of runs submitted within 10 calendar days of PeriodEnd
	for _, r := range allRuns {
		if r.SubmittedAt == nil {
			continue
		}
		periodEnd, err := time.Parse("2006-01-02", r.PeriodEnd)
		if err != nil {
			continue
		}
		out.SubmissionTotal++
		daysDiff := r.SubmittedAt.Sub(periodEnd).Hours() / 24
		if daysDiff <= 10 {
			out.SubmissionOnTime++
		}
	}
	out.SubmissionTimelinessPct = pct(out.SubmissionOnTime, out.SubmissionTotal)

	// ── KPI 2: Processing timeliness ─────────────────────────────────────────
	// % of submitted runs acknowledged within 10 calendar days of submission
	for _, r := range allRuns {
		if r.SubmittedAt == nil || r.AcknowledgedAt == nil {
			continue
		}
		out.ProcessingTotal++
		daysDiff := r.AcknowledgedAt.Sub(*r.SubmittedAt).Hours() / 24
		if daysDiff <= 10 {
			out.ProcessingOnTime++
		}
	}
	out.ProcessingTimelinessPct = pct(out.ProcessingOnTime, out.ProcessingTotal)

	// ── KPI 3: First-time acceptance rate ────────────────────────────────────
	// % of submitted/acknowledged/settled runs that required no amendment (ParentRunID IS NULL)
	concludedStatuses := map[string]bool{"submitted": true, "acknowledged": true, "settled": true}
	for _, r := range allRuns {
		if !concludedStatuses[r.Status] {
			continue
		}
		out.FirstTimeTotal++
		if r.ParentRunID == nil {
			out.FirstTimeAccepted++
		}
	}
	out.FirstTimeAcceptancePct = pct(out.FirstTimeAccepted, out.FirstTimeTotal)

	// ── KPI 4: Average error resolution days ─────────────────────────────────
	// Avg (UpdatedAt - CreatedAt) for runs that went through validation and are now past it
	resolvedStatuses := map[string]bool{"validated": true, "submitted": true, "acknowledged": true, "settled": true}
	var totalDays float64
	for _, r := range allRuns {
		if !resolvedStatuses[r.Status] {
			continue
		}
		// Only count runs that actually went through validation (have results stored)
		var count int64
		DB.Model(&models.RIValidationResult{}).Where("run_id = ?", r.RunID).Count(&count)
		if count == 0 {
			continue
		}
		days := r.UpdatedAt.Sub(r.CreatedAt).Hours() / 24
		totalDays += days
		out.ErrorResolutionSamples++
	}
	if out.ErrorResolutionSamples > 0 {
		out.AvgErrorResolutionDays = math.Round(totalDays/float64(out.ErrorResolutionSamples)*10) / 10
	}

	// ── KPI 5: Settlement timeliness ─────────────────────────────────────────
	// % of technical accounts settled within 30 days of AgreedAt
	var accounts []models.TechnicalAccount
	taQ := DB.Where("period_start >= ? AND period_start <= ? AND status = 'settled'", periodFrom, periodTo)
	if treatyID > 0 {
		taQ = taQ.Where("treaty_id = ?", treatyID)
	}
	taQ.Find(&accounts)
	for _, a := range accounts {
		if a.AgreedAt == nil || a.SettledAt == nil {
			continue
		}
		out.SettlementTotal++
		daysDiff := a.SettledAt.Sub(*a.AgreedAt).Hours() / 24
		if daysDiff <= 30 {
			out.SettlementOnTime++
		}
	}
	out.SettlementTimelinessPct = pct(out.SettlementOnTime, out.SettlementTotal)

	// ── KPI 6: Claims register completeness ──────────────────────────────────
	// % of large-loss claims in claims runs that have a matching LargeClaimNotice
	var claimsRunIDs []string
	for _, r := range allRuns {
		if r.Type == "claims_run" {
			claimsRunIDs = append(claimsRunIDs, r.RunID)
		}
	}
	if len(claimsRunIDs) > 0 {
		var largeClaimRows []models.RIBordereauxClaimsRow
		DB.Where("run_id IN ? AND large_loss_flag = true", claimsRunIDs).Find(&largeClaimRows)
		out.ClaimsLargeLossTotal = len(largeClaimRows)

		for _, cr := range largeClaimRows {
			var noticeCount int64
			noticeQ := DB.Model(&models.LargeClaimNotice{}).Where("claim_number = ?", cr.ClaimNumber)
			if treatyID > 0 {
				noticeQ = noticeQ.Where("treaty_id = ?", treatyID)
			}
			noticeQ.Count(&noticeCount)
			if noticeCount > 0 {
				out.ClaimsWithNotice++
			}
		}
	}
	out.ClaimsCompletenessPct = pct(out.ClaimsWithNotice, out.ClaimsLargeLossTotal)

	// ── KPI 7: Open query backlog ─────────────────────────────────────────────
	// Runs stuck in validation_failed for more than 30 days
	cutoff := now.AddDate(0, 0, -30).Format("2006-01-02")
	backlogQ := DB.Model(&models.RIBordereauxRun{}).
		Where("status = 'validation_failed' AND created_at < ?", cutoff)
	if treatyID > 0 {
		backlogQ = backlogQ.Where("treaty_id = ?", treatyID)
	}
	var backlogCount int64
	backlogQ.Count(&backlogCount)
	out.OpenQueryBacklog = int(backlogCount)

	return out, nil
}

// pct returns 100 × numerator/denominator, or 100 if denominator is 0 (vacuously compliant)
func pct(numerator, denominator int) float64 {
	if denominator == 0 {
		return 100
	}
	return math.Round(float64(numerator)/float64(denominator)*1000) / 10
}
