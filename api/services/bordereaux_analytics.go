package services

import (
	"api/models"
	"fmt"
	"strings"
	"time"
)

// BordereauxAnalyticsFilters is the input to GetBordereauxAnalytics.
// Period accepts last_7_days | last_30_days | last_quarter | last_year | ytd.
// If Period == "custom" or From/To are set, those override the period window.
type BordereauxAnalyticsFilters struct {
	Period   string
	From     time.Time
	To       time.Time
	SchemeID int
}

type BordereauxTopMetrics struct {
	TotalPremiumVolume  float64 `json:"total_premium_volume"`
	TotalClaimVolume    float64 `json:"total_claim_volume"`
	LossRatio           float64 `json:"loss_ratio"`
	ClaimsCount         int     `json:"claims_count"`
	AvgClaimSize        float64 `json:"avg_claim_size"`
	BordereauxGenerated int     `json:"bordereaux_generated"`
	AvgMatchScore       float64 `json:"avg_match_score"`
}

type BordereauxBenefitRow struct {
	Benefit        string  `json:"benefit"`
	ClaimCount     int     `json:"claim_count"`
	ClaimVolume    float64 `json:"claim_volume"`
	AvgClaimAmount float64 `json:"avg_claim_amount"`
	ApprovedCount  int     `json:"approved_count"`
	ApprovalRate   float64 `json:"approval_rate"`
}

type BordereauxInsurerRow struct {
	InsurerName      string  `json:"insurer_name"`
	BordereauxCount  int     `json:"bordereaux_count"`
	RecordCount      int     `json:"record_count"`
	MarketSharePct   float64 `json:"market_share_pct"`
	AvgMatchScore    float64 `json:"avg_match_score"`
	DiscrepancyCount int     `json:"discrepancy_count"`
}

type BordereauxMonthlyRow struct {
	Month            string  `json:"month"` // YYYY-MM
	BordereauxCount  int     `json:"bordereaux_count"`
	PremiumVolume    float64 `json:"premium_volume"`
	ClaimVolume      float64 `json:"claim_volume"`
	AvgMatchScore    float64 `json:"avg_match_score"`
	DiscrepancyCount int     `json:"discrepancy_count"`
}

type BordereauxAnalyticsResponse struct {
	PeriodFrom         time.Time              `json:"period_from"`
	PeriodTo           time.Time              `json:"period_to"`
	TopMetrics         BordereauxTopMetrics   `json:"top_metrics"`
	BenefitPerformance []BordereauxBenefitRow `json:"benefit_performance"`
	InsurerMetrics     []BordereauxInsurerRow `json:"insurer_metrics"`
	MonthlyKPIs        []BordereauxMonthlyRow `json:"monthly_kpis"`
}

// resolveAnalyticsWindow turns a human period string into a from/to pair. Explicit
// filter.From and filter.To override the period if non-zero.
func resolveAnalyticsWindow(f BordereauxAnalyticsFilters) (time.Time, time.Time) {
	now := time.Now()
	to := now
	var from time.Time
	switch strings.ToLower(strings.TrimSpace(f.Period)) {
	case "last_7_days":
		from = now.AddDate(0, 0, -7)
	case "last_quarter":
		from = now.AddDate(0, 0, -90)
	case "last_year":
		from = now.AddDate(-1, 0, 0)
	case "ytd":
		from = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	case "last_30_days", "":
		from = now.AddDate(0, 0, -30)
	default:
		from = now.AddDate(0, 0, -30)
	}
	if !f.From.IsZero() {
		from = f.From
	}
	if !f.To.IsZero() {
		to = f.To
	}
	return from, to
}

// GetBordereauxAnalytics produces the aggregated metrics consumed by
// BordereauxAnalyticsDashboard.vue. Every value comes from live DB state; if
// a metric has no honest data source it is simply left at zero rather than
// fabricated.
func GetBordereauxAnalytics(f BordereauxAnalyticsFilters) (BordereauxAnalyticsResponse, error) {
	from, to := resolveAnalyticsWindow(f)
	resp := BordereauxAnalyticsResponse{PeriodFrom: from, PeriodTo: to}

	top, err := analyticsTopMetrics(from, to, f.SchemeID)
	if err != nil {
		return resp, err
	}
	resp.TopMetrics = top

	if resp.BenefitPerformance, err = analyticsBenefitPerformance(from, to, f.SchemeID); err != nil {
		return resp, err
	}
	if resp.InsurerMetrics, err = analyticsInsurerMetrics(from, to, f.SchemeID); err != nil {
		return resp, err
	}
	if resp.MonthlyKPIs, err = analyticsMonthlyKPIs(from, to, f.SchemeID); err != nil {
		return resp, err
	}
	return resp, nil
}

func analyticsTopMetrics(from, to time.Time, schemeID int) (BordereauxTopMetrics, error) {
	var m BordereauxTopMetrics

	// Premium volume: sum TotalAnnualPremium across premium bordereaux rows whose
	// parent GeneratedBordereaux was created in window. Joining via GeneratedID
	// (string) — PremiumBordereauxData.BordereauxID references it.
	premiumQ := DB.Table("premium_bordereaux_data AS p").
		Joins("JOIN generated_bordereauxes AS g ON g.generated_id = p.bordereaux_id").
		Where("g.created_at BETWEEN ? AND ?", from, to)
	if schemeID > 0 {
		// PremiumBordereauxData has no direct scheme_id; filter via scheme_name join
		// is deferred (scheme filtering is primarily for claims/generated rows).
		_ = schemeID
	}
	if err := premiumQ.Select("COALESCE(SUM(p.total_annual_premium), 0)").Row().Scan(&m.TotalPremiumVolume); err != nil {
		return m, fmt.Errorf("premium volume: %w", err)
	}

	// Claim stats (count, volume, avg, approval ratio drivers). Uses CreationDate
	// because DateOfEvent is a string.
	claimQ := DB.Model(&models.GroupSchemeClaim{}).
		Where("creation_date BETWEEN ? AND ?", from, to)
	if schemeID > 0 {
		claimQ = claimQ.Where("scheme_id = ?", schemeID)
	}
	var totalClaimVolume float64
	var claimCount int64
	if err := claimQ.Select("COALESCE(SUM(claim_amount), 0)").Row().Scan(&totalClaimVolume); err != nil {
		return m, fmt.Errorf("claim volume: %w", err)
	}
	// Re-issue because the Row().Scan above consumed the statement.
	claimQ = DB.Model(&models.GroupSchemeClaim{}).
		Where("creation_date BETWEEN ? AND ?", from, to)
	if schemeID > 0 {
		claimQ = claimQ.Where("scheme_id = ?", schemeID)
	}
	if err := claimQ.Count(&claimCount).Error; err != nil {
		return m, fmt.Errorf("claim count: %w", err)
	}
	m.TotalClaimVolume = totalClaimVolume
	m.ClaimsCount = int(claimCount)
	if claimCount > 0 {
		m.AvgClaimSize = totalClaimVolume / float64(claimCount)
	}
	if m.TotalPremiumVolume > 0 {
		m.LossRatio = (totalClaimVolume / m.TotalPremiumVolume) * 100
	}

	// Bordereaux generated + avg match score in window.
	genQ := DB.Model(&models.GeneratedBordereaux{}).
		Where("created_at BETWEEN ? AND ?", from, to)
	var bordCount int64
	if err := genQ.Count(&bordCount).Error; err != nil {
		return m, fmt.Errorf("bordereaux count: %w", err)
	}
	m.BordereauxGenerated = int(bordCount)

	var avgScore float64
	if err := DB.Model(&models.BordereauxConfirmation{}).
		Where("imported_at BETWEEN ? AND ?", from, to).
		Select("COALESCE(AVG(match_score), 0)").Row().Scan(&avgScore); err != nil {
		return m, fmt.Errorf("avg match score: %w", err)
	}
	m.AvgMatchScore = avgScore
	return m, nil
}

func analyticsBenefitPerformance(from, to time.Time, schemeID int) ([]BordereauxBenefitRow, error) {
	type row struct {
		Benefit       string
		ClaimCount    int
		ClaimVolume   float64
		ApprovedCount int
	}
	q := DB.Model(&models.GroupSchemeClaim{}).
		Select("benefit_alias AS benefit, COUNT(*) AS claim_count, COALESCE(SUM(claim_amount),0) AS claim_volume, SUM(CASE WHEN LOWER(status) IN ('approved','paid') THEN 1 ELSE 0 END) AS approved_count").
		Where("creation_date BETWEEN ? AND ? AND benefit_alias != ''", from, to).
		Group("benefit_alias").
		Order("claim_volume DESC")
	if schemeID > 0 {
		q = q.Where("scheme_id = ?", schemeID)
	}
	var rows []row
	if err := q.Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("benefit performance: %w", err)
	}
	out := make([]BordereauxBenefitRow, 0, len(rows))
	for _, r := range rows {
		var avg, rate float64
		if r.ClaimCount > 0 {
			avg = r.ClaimVolume / float64(r.ClaimCount)
			rate = (float64(r.ApprovedCount) / float64(r.ClaimCount)) * 100
		}
		out = append(out, BordereauxBenefitRow{
			Benefit:        r.Benefit,
			ClaimCount:     r.ClaimCount,
			ClaimVolume:    r.ClaimVolume,
			AvgClaimAmount: avg,
			ApprovedCount:  r.ApprovedCount,
			ApprovalRate:   rate,
		})
	}
	return out, nil
}

func analyticsInsurerMetrics(from, to time.Time, schemeID int) ([]BordereauxInsurerRow, error) {
	// Aggregate GeneratedBordereaux + its BordereauxConfirmations by insurer name.
	type row struct {
		InsurerName      string
		BordereauxCount  int
		RecordCount      int
		MatchScore       float64
		DiscrepancyCount int
	}
	q := DB.Table("generated_bordereauxes AS g").
		Select(`g.insurer_name AS insurer_name,
			COUNT(DISTINCT g.id) AS bordereaux_count,
			COALESCE(SUM(g.records),0) AS record_count,
			COALESCE(AVG(c.match_score),0) AS match_score,
			COALESCE(SUM(c.discrepancy_count),0) AS discrepancy_count`).
		Joins("LEFT JOIN bordereaux_confirmations AS c ON c.generated_bordereaux_id = g.generated_id").
		Where("g.created_at BETWEEN ? AND ? AND g.insurer_name != ''", from, to).
		Group("g.insurer_name").
		Order("record_count DESC")
	var rows []row
	if err := q.Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("insurer metrics: %w", err)
	}
	var totalRecords int
	for _, r := range rows {
		totalRecords += r.RecordCount
	}
	out := make([]BordereauxInsurerRow, 0, len(rows))
	for _, r := range rows {
		share := 0.0
		if totalRecords > 0 {
			share = (float64(r.RecordCount) / float64(totalRecords)) * 100
		}
		out = append(out, BordereauxInsurerRow{
			InsurerName:      r.InsurerName,
			BordereauxCount:  r.BordereauxCount,
			RecordCount:      r.RecordCount,
			MarketSharePct:   share,
			AvgMatchScore:    r.MatchScore,
			DiscrepancyCount: r.DiscrepancyCount,
		})
	}
	return out, nil
}

func analyticsMonthlyKPIs(from, to time.Time, schemeID int) ([]BordereauxMonthlyRow, error) {
	// Walk each calendar month in the window and aggregate independently. A single
	// GROUP BY TO_CHAR works on Postgres/MySQL but not on SQL Server without
	// different syntax; iterating month-by-month keeps the query cross-DB.
	months := monthRange(from, to)
	out := make([]BordereauxMonthlyRow, 0, len(months))
	for _, m := range months {
		start := m
		end := m.AddDate(0, 1, 0)

		var bordCount int64
		bordQ := DB.Model(&models.GeneratedBordereaux{}).
			Where("created_at >= ? AND created_at < ?", start, end)
		if err := bordQ.Count(&bordCount).Error; err != nil {
			return nil, fmt.Errorf("monthly bordereaux: %w", err)
		}

		var premium float64
		if err := DB.Table("premium_bordereaux_data AS p").
			Joins("JOIN generated_bordereauxes AS g ON g.generated_id = p.bordereaux_id").
			Where("g.created_at >= ? AND g.created_at < ?", start, end).
			Select("COALESCE(SUM(p.total_annual_premium),0)").Row().Scan(&premium); err != nil {
			return nil, fmt.Errorf("monthly premium: %w", err)
		}

		var claimVol float64
		claimQ := DB.Model(&models.GroupSchemeClaim{}).
			Where("creation_date >= ? AND creation_date < ?", start, end)
		if schemeID > 0 {
			claimQ = claimQ.Where("scheme_id = ?", schemeID)
		}
		if err := claimQ.Select("COALESCE(SUM(claim_amount),0)").Row().Scan(&claimVol); err != nil {
			return nil, fmt.Errorf("monthly claim: %w", err)
		}

		var avgMatch float64
		var discrepancy int64
		confQ := DB.Model(&models.BordereauxConfirmation{}).
			Where("imported_at >= ? AND imported_at < ?", start, end)
		if err := confQ.Select("COALESCE(AVG(match_score),0)").Row().Scan(&avgMatch); err != nil {
			return nil, fmt.Errorf("monthly match score: %w", err)
		}
		confQ = DB.Model(&models.BordereauxConfirmation{}).
			Where("imported_at >= ? AND imported_at < ?", start, end)
		if err := confQ.Select("COALESCE(SUM(discrepancy_count),0)").Row().Scan(&discrepancy); err != nil {
			return nil, fmt.Errorf("monthly discrepancy count: %w", err)
		}

		out = append(out, BordereauxMonthlyRow{
			Month:            start.Format("2006-01"),
			BordereauxCount:  int(bordCount),
			PremiumVolume:    premium,
			ClaimVolume:      claimVol,
			AvgMatchScore:    avgMatch,
			DiscrepancyCount: int(discrepancy),
		})
	}
	return out, nil
}

// monthRange returns the first-of-month timestamps covering [from, to], capped
// at 24 months to avoid absurd queries over open-ended windows.
func monthRange(from, to time.Time) []time.Time {
	start := time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, from.Location())
	end := time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, to.Location())
	var months []time.Time
	for m := start; !m.After(end) && len(months) < 24; m = m.AddDate(0, 1, 0) {
		months = append(months, m)
	}
	return months
}
