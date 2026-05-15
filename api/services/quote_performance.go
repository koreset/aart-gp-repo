package services

import (
	"api/models"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// maxExtractRows caps the row count of the extract xlsx export so a
// pathological filter selecting the entire quote book can't blow up
// excelize's in-memory representation. Callers receive a 413 above this.
const maxExtractRows = 50_000

// secondsDiffExpr returns a dialect-appropriate SQL fragment that
// evaluates to the elapsed seconds between two DATETIME columns.
// Caller is responsible for ensuring both columns are non-null when used
// inside an AVG/SUM — wrap with CASE WHEN ... IS NOT NULL THEN ... END.
func secondsDiffExpr(startCol, endCol string) string {
	switch DbBackend {
	case "mssql":
		return fmt.Sprintf("DATEDIFF(SECOND, %s, %s)", startCol, endCol)
	default:
		return fmt.Sprintf("TIMESTAMPDIFF(SECOND, %s, %s)", startCol, endCol)
	}
}

// regionAggExpr returns a dialect-appropriate SQL fragment that
// aggregates distinct region values into a comma-separated string.
func regionAggExpr(col string) string {
	switch DbBackend {
	case "mssql":
		return fmt.Sprintf("STRING_AGG(DISTINCT %s, ',')", col)
	default:
		return fmt.Sprintf("GROUP_CONCAT(DISTINCT %s ORDER BY %s SEPARATOR ',')", col, col)
	}
}

// applyPerformanceFilters mutates the supplied tx to apply the shared
// dashboard filter set. All filters are AND'd across dimensions and OR'd
// within a single multi-select dimension. Region is filtered via EXISTS
// against scheme_categories rather than a JOIN to avoid row-multiplication
// before aggregation.
func applyPerformanceFilters(tx *gorm.DB, q models.QuotePerformanceQuery, submittedDateFilter bool) *gorm.DB {
	if q.From != nil {
		if submittedDateFilter {
			tx = tx.Where("group_pricing_quotes.submitted_at >= ?", q.From)
		} else {
			tx = tx.Where("group_pricing_quotes.creation_date >= ?", q.From)
		}
	}
	if q.To != nil {
		if submittedDateFilter {
			tx = tx.Where("group_pricing_quotes.submitted_at <= ?", q.To)
		} else {
			tx = tx.Where("group_pricing_quotes.creation_date <= ?", q.To)
		}
	}
	if len(q.Users) > 0 {
		tx = tx.Where("group_pricing_quotes.created_by IN ?", q.Users)
	}
	if len(q.QuoteType) > 0 {
		tx = tx.Where("group_pricing_quotes.quote_type IN ?", q.QuoteType)
	}
	if len(q.DistributionChannel) > 0 {
		tx = tx.Where("group_pricing_quotes.distribution_channel IN ?", q.DistributionChannel)
	}
	if len(q.Region) > 0 {
		tx = tx.Where("EXISTS (SELECT 1 FROM scheme_categories sc WHERE sc.quote_id = group_pricing_quotes.id AND sc.region IN ?)", q.Region)
	}
	return tx
}

// GetQuotePerformanceKpis returns the per-user leaderboard rows powering
// the dashboard's main table. Counts and status mix come from
// group_pricing_quotes; annual premium comes from group_risk_quote_stats
// (LEFT JOIN because a quote that hasn't been calculated yet has no
// stats row). Cycle-time averages skip rows where either endpoint is
// NULL via CASE-WHEN so they're not depressed by in-flight quotes.
func GetQuotePerformanceKpis(q models.QuotePerformanceQuery) ([]models.QuotePerformanceKpis, error) {
	type row struct {
		UserName                string
		TotalQuotes             int64
		DraftCount              int64
		SubmittedCount          int64
		ApprovedCount           int64
		RejectedCount           int64
		AcceptedCount           int64
		InForceCount            int64
		AvgTimeToSubmitSecs     *float64
		AvgTimeToApproveSecs    *float64
		AvgTimeToAcceptSecs     *float64
		AvgTotalCycleSecs       *float64
		TotalAnnualPremium      *float64
		PipelineAnnualPremium   *float64
		AvgQuoteValue           *float64
	}

	timeToSubmit := secondsDiffExpr("group_pricing_quotes.creation_date", "group_pricing_quotes.submitted_at")
	timeToApprove := secondsDiffExpr("group_pricing_quotes.submitted_at", "group_pricing_quotes.approved_at")
	timeToAccept := secondsDiffExpr("group_pricing_quotes.approved_at", "group_pricing_quotes.accepted_at")
	totalCycle := secondsDiffExpr("group_pricing_quotes.submitted_at", "group_pricing_quotes.accepted_at")

	// "Ever reached" semantics for conversion math: a quote that's now
	// in_force was definitely submitted at some point — so it should count
	// toward submitted_count, approved_count and accepted_count. Without
	// this, a user whose quotes have all already moved past 'submitted'
	// shows submitted=0 and a nonsense 0% approval rate.
	// draft_count keeps snapshot semantics (it's the in-flight queue).
	selectClause := strings.Join([]string{
		"COALESCE(NULLIF(group_pricing_quotes.created_by, ''), '(unassigned)') AS user_name",
		"COUNT(DISTINCT group_pricing_quotes.id) AS total_quotes",
		"SUM(CASE WHEN group_pricing_quotes.status = 'draft' THEN 1 ELSE 0 END) AS draft_count",
		"SUM(CASE WHEN group_pricing_quotes.submitted_at IS NOT NULL OR group_pricing_quotes.status IN ('submitted','approved','rejected','accepted','in_force') THEN 1 ELSE 0 END) AS submitted_count",
		"SUM(CASE WHEN group_pricing_quotes.approved_at IS NOT NULL OR group_pricing_quotes.status IN ('approved','accepted','in_force') THEN 1 ELSE 0 END) AS approved_count",
		"SUM(CASE WHEN group_pricing_quotes.rejected_at IS NOT NULL OR group_pricing_quotes.status = 'rejected' THEN 1 ELSE 0 END) AS rejected_count",
		"SUM(CASE WHEN group_pricing_quotes.accepted_at IS NOT NULL OR group_pricing_quotes.status IN ('accepted','in_force') THEN 1 ELSE 0 END) AS accepted_count",
		"SUM(CASE WHEN group_pricing_quotes.in_force_at IS NOT NULL OR group_pricing_quotes.status = 'in_force' THEN 1 ELSE 0 END) AS in_force_count",
		fmt.Sprintf("AVG(CASE WHEN group_pricing_quotes.submitted_at IS NOT NULL THEN %s END) AS avg_time_to_submit_secs", timeToSubmit),
		fmt.Sprintf("AVG(CASE WHEN group_pricing_quotes.approved_at IS NOT NULL AND group_pricing_quotes.submitted_at IS NOT NULL THEN %s END) AS avg_time_to_approve_secs", timeToApprove),
		fmt.Sprintf("AVG(CASE WHEN group_pricing_quotes.accepted_at IS NOT NULL AND group_pricing_quotes.approved_at IS NOT NULL THEN %s END) AS avg_time_to_accept_secs", timeToAccept),
		fmt.Sprintf("AVG(CASE WHEN group_pricing_quotes.accepted_at IS NOT NULL AND group_pricing_quotes.submitted_at IS NOT NULL THEN %s END) AS avg_total_cycle_secs", totalCycle),
		"SUM(CASE WHEN group_pricing_quotes.status IN ('accepted', 'in_force') THEN COALESCE(grqs.annual_premium, 0) ELSE 0 END) AS total_annual_premium",
		"SUM(CASE WHEN group_pricing_quotes.status IN ('submitted', 'approved') THEN COALESCE(grqs.annual_premium, 0) ELSE 0 END) AS pipeline_annual_premium",
		"AVG(NULLIF(grqs.annual_premium, 0)) AS avg_quote_value",
	}, ", ")

	var rows []row
	tx := DB.Table("group_pricing_quotes").
		Select(selectClause).
		Joins("LEFT JOIN group_risk_quote_stats grqs ON grqs.quote_id = group_pricing_quotes.id").
		Group("COALESCE(NULLIF(group_pricing_quotes.created_by, ''), '(unassigned)')")

	tx = applyPerformanceFilters(tx, q, false)

	if err := tx.Scan(&rows).Error; err != nil {
		return nil, err
	}

	// Pull SLA breach counts per user in scope so we can stitch them
	// into the leaderboard without a second round-trip from the client.
	breachByUser, transitionByUser, err := slaBreachCountsByUser(q)
	if err != nil {
		return nil, err
	}

	result := make([]models.QuotePerformanceKpis, 0, len(rows))
	for _, r := range rows {
		kpi := models.QuotePerformanceKpis{
			UserName:       r.UserName,
			TotalQuotes:    r.TotalQuotes,
			DraftCount:     r.DraftCount,
			SubmittedCount: r.SubmittedCount,
			ApprovedCount:  r.ApprovedCount,
			RejectedCount:  r.RejectedCount,
			AcceptedCount:  r.AcceptedCount,
			InForceCount:   r.InForceCount,
		}
		if r.SubmittedCount > 0 {
			kpi.ApprovalRate = float64(r.ApprovedCount) / float64(r.SubmittedCount)
			kpi.ConversionRate = float64(r.AcceptedCount) / float64(r.SubmittedCount)
			kpi.RejectionRate = float64(r.RejectedCount) / float64(r.SubmittedCount)
		}
		if r.ApprovedCount > 0 {
			kpi.AcceptanceRate = float64(r.AcceptedCount) / float64(r.ApprovedCount)
		}
		kpi.AvgTimeToSubmitHrs = secsToHours(r.AvgTimeToSubmitSecs)
		kpi.AvgTimeToApproveHrs = secsToHours(r.AvgTimeToApproveSecs)
		kpi.AvgTimeToAcceptHrs = secsToHours(r.AvgTimeToAcceptSecs)
		kpi.AvgTotalCycleHrs = secsToHours(r.AvgTotalCycleSecs)
		kpi.TotalAnnualPremium = derefFloat(r.TotalAnnualPremium)
		kpi.PipelineAnnualPremium = derefFloat(r.PipelineAnnualPremium)
		kpi.AvgQuoteValue = derefFloat(r.AvgQuoteValue)
		kpi.SlaBreachCount = breachByUser[r.UserName]
		kpi.SlaTransitionCount = transitionByUser[r.UserName]
		if kpi.SlaTransitionCount > 0 {
			kpi.SlaCompliancePct = 1 - float64(kpi.SlaBreachCount)/float64(kpi.SlaTransitionCount)
		}
		result = append(result, kpi)
	}

	sort.Slice(result, func(i, j int) bool { return result[i].TotalQuotes > result[j].TotalQuotes })
	return result, nil
}

// GetQuoteFunnel returns one row per workflow stage describing how many
// quotes have ever reached that stage and the average dwell time in the
// prior stage before reaching it. Dwell time is sourced from the audit
// table's precomputed duration_from_prev_secs, excluding synthetic
// backfill rows so estimates don't pollute the chart.
func GetQuoteFunnel(q models.QuotePerformanceQuery) ([]models.FunnelStage, error) {
	type row struct {
		Stage         string
		Count         int64
		AvgDwellSecs  *float64
	}

	stages := []models.Status{
		models.StatusDraft,
		models.StatusSubmitted,
		models.StatusApproved,
		models.StatusAccepted,
		models.StatusInForce,
	}

	out := make([]models.FunnelStage, 0, len(stages))
	for _, stage := range stages {
		var r row
		r.Stage = string(stage)

		// Count quotes that have ever reached this stage. For draft we
		// take "every quote exists in draft at some point" → total count.
		countTx := DB.Table("group_pricing_quotes")
		countTx = applyPerformanceFilters(countTx, q, false)
		if stage == models.StatusDraft {
			countTx = countTx.Select("COUNT(*)")
		} else {
			countTx = countTx.Where(stageReachedClause(stage))
			countTx = countTx.Select("COUNT(*)")
		}
		if err := countTx.Row().Scan(&r.Count); err != nil {
			return nil, err
		}

		// Average dwell: time spent in the previous stage. From the
		// audit table where new_status = stage and synthetic = false.
		if stage != models.StatusDraft {
			dwellTx := DB.Table("group_pricing_quote_status_audits a").
				Joins("JOIN group_pricing_quotes ON group_pricing_quotes.id = a.quote_id").
				Where("a.new_status = ?", stage).
				Where("a.synthetic = ?", false).
				Where("a.duration_from_prev_secs > 0")
			dwellTx = applyPerformanceFilters(dwellTx, q, false)
			dwellTx = dwellTx.Select("AVG(a.duration_from_prev_secs)")
			_ = dwellTx.Row().Scan(&r.AvgDwellSecs)
		}

		out = append(out, models.FunnelStage{
			Stage:         r.Stage,
			Count:         r.Count,
			AvgDwellHours: secsToHours(r.AvgDwellSecs),
		})
	}
	return out, nil
}

// stageReachedClause returns the WHERE-fragment that determines whether
// a quote has reached a given stage. A quote has reached an approved
// stage if its approved_at is non-null OR its current status implies the
// stage was passed through (e.g. an accepted quote was definitely
// approved at some point, even on a backfilled timestamp).
func stageReachedClause(stage models.Status) string {
	switch stage {
	case models.StatusSubmitted:
		return "group_pricing_quotes.submitted_at IS NOT NULL OR group_pricing_quotes.status IN ('submitted','approved','rejected','accepted','in_force')"
	case models.StatusApproved:
		return "group_pricing_quotes.approved_at IS NOT NULL OR group_pricing_quotes.status IN ('approved','accepted','in_force')"
	case models.StatusAccepted:
		return "group_pricing_quotes.accepted_at IS NOT NULL OR group_pricing_quotes.status IN ('accepted','in_force')"
	case models.StatusInForce:
		return "group_pricing_quotes.in_force_at IS NOT NULL OR group_pricing_quotes.status = 'in_force'"
	default:
		return "1=1"
	}
}

// GetQuoteTrend buckets submitted / approved / accepted / rejected
// counts by day, week, or month. Sourced from the per-status milestone
// timestamps on the quote so the timeline reflects actual workflow
// activity, not creation date.
func GetQuoteTrend(q models.QuotePerformanceQuery, bucket string) ([]models.TrendBucket, error) {
	bucketExpr, err := bucketExprForBackend(bucket)
	if err != nil {
		return nil, err
	}

	type row struct {
		Bucket    string
		Submitted int64
		Approved  int64
		Accepted  int64
		Rejected  int64
	}

	// Build a UNION-style accounting where each milestone contributes
	// to exactly one column. Easier to express as separate sub-queries.
	type stageRow struct {
		Bucket string
		Count  int64
	}

	gather := func(timestampCol string) (map[string]int64, error) {
		tx := DB.Table("group_pricing_quotes").
			Select(fmt.Sprintf("%s AS bucket, COUNT(*) AS count", strings.ReplaceAll(bucketExpr, "$col$", timestampCol))).
			Where(timestampCol + " IS NOT NULL").
			Group(strings.ReplaceAll(bucketExpr, "$col$", timestampCol))
		tx = applyPerformanceFilters(tx, q, false)
		var rows []stageRow
		if err := tx.Scan(&rows).Error; err != nil {
			return nil, err
		}
		m := make(map[string]int64, len(rows))
		for _, r := range rows {
			m[r.Bucket] += r.Count
		}
		return m, nil
	}

	submitted, err := gather("submitted_at")
	if err != nil {
		return nil, err
	}
	approved, err := gather("approved_at")
	if err != nil {
		return nil, err
	}
	accepted, err := gather("accepted_at")
	if err != nil {
		return nil, err
	}
	rejected, err := gather("rejected_at")
	if err != nil {
		return nil, err
	}

	buckets := map[string]bool{}
	for k := range submitted {
		buckets[k] = true
	}
	for k := range approved {
		buckets[k] = true
	}
	for k := range accepted {
		buckets[k] = true
	}
	for k := range rejected {
		buckets[k] = true
	}

	out := make([]models.TrendBucket, 0, len(buckets))
	for k := range buckets {
		out = append(out, models.TrendBucket{
			Bucket:    k,
			Submitted: submitted[k],
			Approved:  approved[k],
			Accepted:  accepted[k],
			Rejected:  rejected[k],
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Bucket < out[j].Bucket })
	return out, nil
}

func bucketExprForBackend(bucket string) (string, error) {
	switch bucket {
	case "day", "":
		if DbBackend == "mssql" {
			return "CONVERT(VARCHAR(10), $col$, 23)", nil
		}
		return "DATE_FORMAT($col$, '%Y-%m-%d')", nil
	case "week":
		if DbBackend == "mssql" {
			return "CONCAT(DATEPART(YEAR, $col$), '-W', RIGHT('00' + CAST(DATEPART(WEEK, $col$) AS VARCHAR(2)), 2))", nil
		}
		return "DATE_FORMAT($col$, '%x-W%v')", nil
	case "month":
		if DbBackend == "mssql" {
			return "CONVERT(VARCHAR(7), $col$, 23)", nil
		}
		return "DATE_FORMAT($col$, '%Y-%m')", nil
	default:
		return "", fmt.Errorf("invalid bucket: %q (expected day, week, or month)", bucket)
	}
}

// GetQuoteSlaBreaches returns per-transition and per-user SLA breach
// counts. A breach is an audit row where the actual duration exceeded
// the configured target hours for that (from_status, to_status,
// quote_type) tuple. Synthetic backfill rows are excluded so estimated
// timestamps don't inflate breach counts.
func GetQuoteSlaBreaches(q models.QuotePerformanceQuery) (models.SlaBreachSummary, error) {
	summary := models.SlaBreachSummary{}

	targets, err := ListQuoteSlaTargets()
	if err != nil {
		return summary, err
	}
	// Build a lookup keyed by (from, to, quote_type). The "" quote_type
	// is the default fallback when no type-specific row matches.
	type key struct{ from, to, qtype string }
	targetMap := map[key]models.QuoteSlaTarget{}
	for _, t := range targets {
		if !t.Active {
			continue
		}
		targetMap[key{string(t.FromStatus), string(t.ToStatus), t.QuoteType}] = t
	}
	resolveTarget := func(from, to, qtype string) (models.QuoteSlaTarget, bool) {
		if t, ok := targetMap[key{from, to, qtype}]; ok {
			return t, true
		}
		if t, ok := targetMap[key{from, to, ""}]; ok {
			return t, true
		}
		return models.QuoteSlaTarget{}, false
	}

	// Pull all real (non-synthetic) audit rows in scope plus the
	// quote_type so we can resolve the right SLA target per row.
	type auditRow struct {
		FromStatus string
		ToStatus   string
		QuoteType  string
		ChangedBy  string
		Duration   int64
	}
	tx := DB.Table("group_pricing_quote_status_audits a").
		Select("COALESCE(a.old_status, '') AS from_status, a.new_status AS to_status, group_pricing_quotes.quote_type AS quote_type, a.changed_by AS changed_by, a.duration_from_prev_secs AS duration").
		Joins("JOIN group_pricing_quotes ON group_pricing_quotes.id = a.quote_id").
		Where("a.synthetic = ?", false).
		Where("a.duration_from_prev_secs > 0")
	tx = applyPerformanceFilters(tx, q, false)

	var rows []auditRow
	if err := tx.Scan(&rows).Error; err != nil {
		return summary, err
	}

	type bkey struct{ from, to string }
	breachByTrans := map[bkey]int64{}
	totalByTrans := map[bkey]int64{}
	targetByTrans := map[bkey]float64{}
	breachByUserMap := map[string]int64{}
	totalByUserMap := map[string]int64{}

	for _, r := range rows {
		target, ok := resolveTarget(r.FromStatus, r.ToStatus, r.QuoteType)
		if !ok {
			continue
		}
		bk := bkey{r.FromStatus, r.ToStatus}
		totalByTrans[bk]++
		targetByTrans[bk] = target.TargetHours
		totalByUserMap[r.ChangedBy]++
		if float64(r.Duration) > target.TargetHours*3600 {
			breachByTrans[bk]++
			breachByUserMap[r.ChangedBy]++
		}
	}

	for bk, total := range totalByTrans {
		summary.BreachesByTransition = append(summary.BreachesByTransition, models.SlaBreachByTransition{
			FromStatus:      models.Status(bk.from),
			ToStatus:        models.Status(bk.to),
			BreachCount:     breachByTrans[bk],
			TransitionCount: total,
			TargetHours:     targetByTrans[bk],
		})
	}
	sort.Slice(summary.BreachesByTransition, func(i, j int) bool {
		a, b := summary.BreachesByTransition[i], summary.BreachesByTransition[j]
		if a.FromStatus != b.FromStatus {
			return a.FromStatus < b.FromStatus
		}
		return a.ToStatus < b.ToStatus
	})

	for user, total := range totalByUserMap {
		summary.BreachesByUser = append(summary.BreachesByUser, models.SlaBreachByUser{
			UserName:        user,
			BreachCount:     breachByUserMap[user],
			TransitionCount: total,
		})
	}
	sort.Slice(summary.BreachesByUser, func(i, j int) bool {
		return summary.BreachesByUser[i].BreachCount > summary.BreachesByUser[j].BreachCount
	})

	return summary, nil
}

// slaBreachCountsByUser is the helper called from GetQuotePerformanceKpis
// so the per-user breach + transition counts can be merged into the
// leaderboard without a second client round-trip.
func slaBreachCountsByUser(q models.QuotePerformanceQuery) (breaches, transitions map[string]int64, err error) {
	sum, err := GetQuoteSlaBreaches(q)
	if err != nil {
		return nil, nil, err
	}
	breaches = make(map[string]int64, len(sum.BreachesByUser))
	transitions = make(map[string]int64, len(sum.BreachesByUser))
	for _, b := range sum.BreachesByUser {
		breaches[b.UserName] = b.BreachCount
		transitions[b.UserName] = b.TransitionCount
	}
	return breaches, transitions, nil
}

// ExtractQuotes returns a paginated extract grid row set plus the total
// match count for the supplied filter. Used by both the in-app grid and
// the xlsx export (which pages through the result set).
func ExtractQuotes(f models.QuoteExtractFilter) ([]models.QuoteExtractRow, int64, error) {
	tx := buildExtractQuery(f)

	var total int64
	if err := tx.Session(&gorm.Session{}).
		Select("COUNT(DISTINCT group_pricing_quotes.id)").
		Row().Scan(&total); err != nil {
		return nil, 0, err
	}

	if f.PageSize <= 0 {
		f.PageSize = 50
	}
	if f.Page <= 0 {
		f.Page = 1
	}

	regionsExpr := regionAggExpr("sc.region")
	cycleExpr := secondsDiffExpr("group_pricing_quotes.submitted_at", "group_pricing_quotes.accepted_at")

	selectClause := strings.Join([]string{
		"group_pricing_quotes.id",
		"group_pricing_quotes.quote_name",
		"group_pricing_quotes.quote_type",
		"group_pricing_quotes.scheme_name",
		"group_pricing_quotes.industry",
		fmt.Sprintf("%s AS regions", regionsExpr),
		"group_pricing_quotes.distribution_channel",
		"group_pricing_quotes.status",
		"group_pricing_quotes.created_by",
		"group_pricing_quotes.reviewer",
		"group_pricing_quotes.approved_by",
		"group_pricing_quotes.creation_date",
		"group_pricing_quotes.submitted_at",
		"group_pricing_quotes.approved_at",
		"group_pricing_quotes.accepted_at",
		"group_pricing_quotes.in_force_at",
		"group_pricing_quotes.rejected_at",
		"COALESCE(grqs.annual_premium, 0) AS annual_premium",
		"COALESCE(grqs.member_count, 0) AS member_count",
		fmt.Sprintf("CASE WHEN group_pricing_quotes.submitted_at IS NOT NULL AND group_pricing_quotes.accepted_at IS NOT NULL THEN %s ELSE NULL END AS cycle_secs", cycleExpr),
	}, ", ")

	orderClause := sanitizeExtractOrderBy(f.OrderBy)

	type scanRow struct {
		ID                  int
		QuoteName           string
		QuoteType           string
		SchemeName          string
		Industry            string
		Regions             *string
		DistributionChannel string
		Status              models.Status
		CreatedBy           string
		Reviewer            string
		ApprovedBy          string
		CreationDate        time.Time
		SubmittedAt         *time.Time
		ApprovedAt          *time.Time
		AcceptedAt          *time.Time
		InForceAt           *time.Time
		RejectedAt          *time.Time
		AnnualPremium       float64
		MemberCount         int
		CycleSecs           *float64
	}

	var raw []scanRow
	err := tx.Select(selectClause).
		Joins("LEFT JOIN scheme_categories sc ON sc.quote_id = group_pricing_quotes.id").
		Joins("LEFT JOIN group_risk_quote_stats grqs ON grqs.quote_id = group_pricing_quotes.id").
		Group("group_pricing_quotes.id, group_pricing_quotes.quote_name, group_pricing_quotes.quote_type, group_pricing_quotes.scheme_name, group_pricing_quotes.industry, group_pricing_quotes.distribution_channel, group_pricing_quotes.status, group_pricing_quotes.created_by, group_pricing_quotes.reviewer, group_pricing_quotes.approved_by, group_pricing_quotes.creation_date, group_pricing_quotes.submitted_at, group_pricing_quotes.approved_at, group_pricing_quotes.accepted_at, group_pricing_quotes.in_force_at, group_pricing_quotes.rejected_at, grqs.annual_premium, grqs.member_count").
		Order(orderClause).
		Limit(f.PageSize).
		Offset((f.Page - 1) * f.PageSize).
		Scan(&raw).Error

	if err != nil {
		return nil, 0, err
	}

	out := make([]models.QuoteExtractRow, 0, len(raw))
	for _, r := range raw {
		row := models.QuoteExtractRow{
			ID:                  r.ID,
			QuoteName:           r.QuoteName,
			QuoteType:           r.QuoteType,
			SchemeName:          r.SchemeName,
			Industry:            r.Industry,
			DistributionChannel: r.DistributionChannel,
			Status:              r.Status,
			CreatedBy:           r.CreatedBy,
			Reviewer:            r.Reviewer,
			ApprovedBy:          r.ApprovedBy,
			CreationDate:        r.CreationDate,
			SubmittedAt:         r.SubmittedAt,
			ApprovedAt:          r.ApprovedAt,
			AcceptedAt:          r.AcceptedAt,
			InForceAt:           r.InForceAt,
			RejectedAt:          r.RejectedAt,
			AnnualPremium:       r.AnnualPremium,
			MemberCount:         r.MemberCount,
		}
		if r.Regions != nil {
			row.Regions = *r.Regions
		}
		if r.CycleSecs != nil {
			hrs := *r.CycleSecs / 3600.0
			row.CycleHours = &hrs
		}
		out = append(out, row)
	}

	return out, total, nil
}

// buildExtractQuery applies every filter from QuoteExtractFilter to a
// fresh base query on group_pricing_quotes. The returned tx still needs
// SELECT / JOIN / GROUP BY / ORDER / LIMIT applied by the caller.
func buildExtractQuery(f models.QuoteExtractFilter) *gorm.DB {
	tx := DB.Table("group_pricing_quotes")
	if len(f.CreatedBy) > 0 {
		tx = tx.Where("group_pricing_quotes.created_by IN ?", f.CreatedBy)
	}
	if len(f.Reviewer) > 0 {
		tx = tx.Where("group_pricing_quotes.reviewer IN ?", f.Reviewer)
	}
	if len(f.Status) > 0 {
		tx = tx.Where("group_pricing_quotes.status IN ?", f.Status)
	}
	if len(f.QuoteType) > 0 {
		tx = tx.Where("group_pricing_quotes.quote_type IN ?", f.QuoteType)
	}
	if len(f.Industry) > 0 {
		tx = tx.Where("group_pricing_quotes.industry IN ?", f.Industry)
	}
	if len(f.DistributionChannel) > 0 {
		tx = tx.Where("group_pricing_quotes.distribution_channel IN ?", f.DistributionChannel)
	}
	if len(f.Region) > 0 {
		tx = tx.Where("EXISTS (SELECT 1 FROM scheme_categories sc2 WHERE sc2.quote_id = group_pricing_quotes.id AND sc2.region IN ?)", f.Region)
	}
	if f.From != nil {
		tx = tx.Where("group_pricing_quotes.creation_date >= ?", f.From)
	}
	if f.To != nil {
		tx = tx.Where("group_pricing_quotes.creation_date <= ?", f.To)
	}
	if f.MinAnnualPremium != nil {
		tx = tx.Where("EXISTS (SELECT 1 FROM group_risk_quote_stats grqs2 WHERE grqs2.quote_id = group_pricing_quotes.id AND grqs2.annual_premium >= ?)", *f.MinAnnualPremium)
	}
	if f.MaxAnnualPremium != nil {
		tx = tx.Where("EXISTS (SELECT 1 FROM group_risk_quote_stats grqs2 WHERE grqs2.quote_id = group_pricing_quotes.id AND grqs2.annual_premium <= ?)", *f.MaxAnnualPremium)
	}
	return tx
}

// sanitizeExtractOrderBy whitelists order-by clauses so the extract
// endpoint can't be tricked into SQL injection via the OrderBy field.
// Falls back to creation_date DESC when input is unrecognised.
func sanitizeExtractOrderBy(input string) string {
	allowed := map[string]string{
		"id":             "group_pricing_quotes.id",
		"quote_name":     "group_pricing_quotes.quote_name",
		"creation_date":  "group_pricing_quotes.creation_date",
		"submitted_at":   "group_pricing_quotes.submitted_at",
		"approved_at":    "group_pricing_quotes.approved_at",
		"accepted_at":    "group_pricing_quotes.accepted_at",
		"status":         "group_pricing_quotes.status",
		"created_by":     "group_pricing_quotes.created_by",
		"annual_premium": "annual_premium",
		"member_count":   "member_count",
	}
	field := strings.TrimSpace(input)
	dir := "DESC"
	if idx := strings.Index(field, " "); idx > 0 {
		dirPart := strings.ToUpper(strings.TrimSpace(field[idx+1:]))
		if dirPart == "ASC" || dirPart == "DESC" {
			dir = dirPart
		}
		field = field[:idx]
	}
	col, ok := allowed[field]
	if !ok {
		return "group_pricing_quotes.creation_date DESC"
	}
	return col + " " + dir
}

// ExportQuoteExtractXlsx renders the extract grid as a .xlsx workbook
// for download. Sheet1 is the data, frozen top row, date columns
// formatted, header row bold.
func ExportQuoteExtractXlsx(rows []models.QuoteExtractRow) ([]byte, error) {
	if len(rows) > maxExtractRows {
		return nil, fmt.Errorf("extract too large: %d rows exceeds the %d-row cap, narrow your filter", len(rows), maxExtractRows)
	}

	f := excelize.NewFile()
	sheet := "Quotes"
	idx, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(idx)
	_ = f.DeleteSheet("Sheet1")

	headers := []string{
		"ID", "Quote Number", "Quote Type", "Scheme Name", "Industry", "Regions",
		"Distribution Channel", "Status", "Created By", "Reviewer", "Approved By",
		"Created", "Submitted", "Approved", "Accepted", "In Force", "Rejected",
		"Annual Premium", "Member Count", "Cycle (hrs)",
	}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheet, cell, h)
	}

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E0E7FF"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	endCol, _ := excelize.CoordinatesToCellName(len(headers), 1)
	_ = f.SetCellStyle(sheet, "A1", endCol, headerStyle)

	dateStyle, _ := f.NewStyle(&excelize.Style{NumFmt: 22}) // m/d/yy h:mm

	for ri, r := range rows {
		row := ri + 2
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", row), r.ID)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", row), r.QuoteName)
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", row), r.QuoteType)
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", row), r.SchemeName)
		_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", row), r.Industry)
		_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", row), r.Regions)
		_ = f.SetCellValue(sheet, fmt.Sprintf("G%d", row), r.DistributionChannel)
		_ = f.SetCellValue(sheet, fmt.Sprintf("H%d", row), string(r.Status))
		_ = f.SetCellValue(sheet, fmt.Sprintf("I%d", row), r.CreatedBy)
		_ = f.SetCellValue(sheet, fmt.Sprintf("J%d", row), r.Reviewer)
		_ = f.SetCellValue(sheet, fmt.Sprintf("K%d", row), r.ApprovedBy)
		setDateCell(f, sheet, fmt.Sprintf("L%d", row), &r.CreationDate, dateStyle)
		setDateCell(f, sheet, fmt.Sprintf("M%d", row), r.SubmittedAt, dateStyle)
		setDateCell(f, sheet, fmt.Sprintf("N%d", row), r.ApprovedAt, dateStyle)
		setDateCell(f, sheet, fmt.Sprintf("O%d", row), r.AcceptedAt, dateStyle)
		setDateCell(f, sheet, fmt.Sprintf("P%d", row), r.InForceAt, dateStyle)
		setDateCell(f, sheet, fmt.Sprintf("Q%d", row), r.RejectedAt, dateStyle)
		_ = f.SetCellValue(sheet, fmt.Sprintf("R%d", row), r.AnnualPremium)
		_ = f.SetCellValue(sheet, fmt.Sprintf("S%d", row), r.MemberCount)
		if r.CycleHours != nil {
			_ = f.SetCellValue(sheet, fmt.Sprintf("T%d", row), *r.CycleHours)
		}
	}

	// Freeze the header row.
	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func setDateCell(f *excelize.File, sheet, cell string, t *time.Time, styleID int) {
	if t == nil || t.IsZero() {
		return
	}
	_ = f.SetCellValue(sheet, cell, *t)
	_ = f.SetCellStyle(sheet, cell, cell, styleID)
}

// ListQuoteSlaTargets returns every configured SLA target, ordered by
// from_status then to_status so the admin settings UI is stable.
func ListQuoteSlaTargets() ([]models.QuoteSlaTarget, error) {
	var targets []models.QuoteSlaTarget
	if err := DB.Order("from_status, to_status, quote_type").Find(&targets).Error; err != nil {
		return nil, err
	}
	return targets, nil
}

// UpsertQuoteSlaTarget creates or updates a SLA target row. The
// (from_status, to_status, quote_type) tuple is the natural key — if a
// row with that tuple already exists and the caller didn't supply an ID,
// the existing row is updated in place so admins editing through the
// settings UI don't accidentally duplicate.
func UpsertQuoteSlaTarget(t models.QuoteSlaTarget, user models.AppUser) (models.QuoteSlaTarget, error) {
	t.UpdatedBy = user.UserName
	t.UpdatedAt = time.Now()
	if t.WarningPctOfSla <= 0 || t.WarningPctOfSla > 1 {
		t.WarningPctOfSla = 0.8
	}
	if t.TargetHours <= 0 {
		return t, fmt.Errorf("target_hours must be greater than 0")
	}
	if t.FromStatus == "" || t.ToStatus == "" {
		return t, fmt.Errorf("from_status and to_status are required")
	}

	if t.ID == 0 {
		var existing models.QuoteSlaTarget
		if err := DB.Where("from_status = ? AND to_status = ? AND quote_type = ?", t.FromStatus, t.ToStatus, t.QuoteType).First(&existing).Error; err == nil {
			t.ID = existing.ID
		}
	}

	if err := DB.Save(&t).Error; err != nil {
		return t, err
	}
	return t, nil
}

// DeactivateQuoteSlaTarget soft-disables a SLA target by flipping its
// active flag. Used by the settings UI's "delete" action so historical
// breach calculations against deactivated targets remain reproducible.
func DeactivateQuoteSlaTarget(id int, user models.AppUser) error {
	return DB.Model(&models.QuoteSlaTarget{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"active":     false,
			"updated_by": user.UserName,
			"updated_at": time.Now(),
		}).Error
}

// secsToHours converts a nullable seconds value (raw DB AVG) into hours
// (float), returning 0 when nil — leaderboard rows for users with no
// completed transitions show 0 cycle time, not "null".
func secsToHours(secs *float64) float64 {
	if secs == nil {
		return 0
	}
	return *secs / 3600.0
}

func derefFloat(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}
