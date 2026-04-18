package services

import (
	appLog "api/log"
	"api/models"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// CessionResult holds the computed cession for one member's sum assured
type CessionResult struct {
	CessionBasis    string
	RetentionAmount float64
	CededAmount     float64
	CessionPct      float64
}

// CalculateMemberCession applies treaty terms to compute cession for a single member
func CalculateMemberCession(sumAssured float64, treaty models.ReinsuranceTreaty) CessionResult {
	// Check if three-tier structure is configured (at least Level 1 has bounds and proportion)
	if treaty.Level1Upperbound > 0 && treaty.Level1CededProportion > 0 {
		return calculateTieredCession(sumAssured, treaty)
	}

	// Legacy cession logic
	switch treaty.TreatyType {
	case "proportional", "quota_share":
		// Three-tier structure is required for proportional treaties
		// If not configured, return all retained (no cession)
		return CessionResult{
			CessionBasis:    "proportional",
			RetentionAmount: sumAssured,
			CededAmount:     0,
			CessionPct:      0,
		}
	case "surplus":
		retention := treaty.RetentionAmount
		if sumAssured <= retention {
			return CessionResult{
				CessionBasis:    "surplus",
				RetentionAmount: sumAssured,
				CededAmount:     0,
				CessionPct:      0,
			}
		}
		surplus := sumAssured - retention
		maxCession := retention * float64(treaty.SurplusLines)
		if surplus > maxCession {
			surplus = maxCession
		}
		pct := (surplus / sumAssured) * 100
		return CessionResult{
			CessionBasis:    "surplus",
			RetentionAmount: retention,
			CededAmount:     surplus,
			CessionPct:      pct,
		}
	case "xl_risk", "xl_event":
		// XL cession is computed per claim, not per sum assured
		return CessionResult{
			CessionBasis:    "xl_excess",
			RetentionAmount: sumAssured,
			CededAmount:     0,
			CessionPct:      0,
		}
	default:
		// Three-tier structure is required for proportional/quota_share treaties
		// If not configured, return all retained (no cession)
		return CessionResult{
			CessionBasis:    "quota_share",
			RetentionAmount: sumAssured,
			CededAmount:     0,
			CessionPct:      0,
		}
	}
}

// calculateTieredCession applies three-tier cession structure based on sum assured bands
func calculateTieredCession(sumAssured float64, treaty models.ReinsuranceTreaty) CessionResult {
	var cessionPct float64
	var tierName string

	// Determine which tier applies - check from Level 3 down to Level 1
	if treaty.Level3Upperbound > 0 && sumAssured >= treaty.Level3Lowerbound && sumAssured <= treaty.Level3Upperbound {
		cessionPct = treaty.Level3CededProportion
		tierName = "level3"
	} else if treaty.Level2Upperbound > 0 && sumAssured >= treaty.Level2Lowerbound && sumAssured <= treaty.Level2Upperbound {
		cessionPct = treaty.Level2CededProportion
		tierName = "level2"
	} else if sumAssured >= treaty.Level1Lowerbound && sumAssured <= treaty.Level1Upperbound {
		cessionPct = treaty.Level1CededProportion
		tierName = "level1"
	} else {
		// Sum assured falls outside all configured tiers - no cession
		return CessionResult{
			CessionBasis:    "tiered_outside_range",
			RetentionAmount: sumAssured,
			CededAmount:     0,
			CessionPct:      0,
		}
	}

	ceded := sumAssured * (cessionPct / 100)
	return CessionResult{
		CessionBasis:    fmt.Sprintf("tiered_%s", tierName),
		RetentionAmount: sumAssured - ceded,
		CededAmount:     ceded,
		CessionPct:      cessionPct,
	}
}

// CalculateMemberCessionWithIncome applies income-based tiered cession
func CalculateMemberCessionWithIncome(annualIncome float64, treaty models.ReinsuranceTreaty) CessionResult {
	// Check if income-based tiers are configured
	if treaty.IncomeLevel1Upperbound == 0 {
		return CessionResult{
			CessionBasis:    "no_income_tiers",
			RetentionAmount: annualIncome,
			CededAmount:     0,
			CessionPct:      0,
		}
	}

	var cessionPct float64
	var tierName string

	// Determine which income tier applies
	if treaty.IncomeLevel3Upperbound > 0 && annualIncome >= treaty.IncomeLevel3Lowerbound && annualIncome <= treaty.IncomeLevel3Upperbound {
		cessionPct = treaty.IncomeLevel3CededProportion
		tierName = "income_level3"
	} else if treaty.IncomeLevel2Upperbound > 0 && annualIncome >= treaty.IncomeLevel2Lowerbound && annualIncome <= treaty.IncomeLevel2Upperbound {
		cessionPct = treaty.IncomeLevel2CededProportion
		tierName = "income_level2"
	} else if annualIncome >= treaty.IncomeLevel1Lowerbound && annualIncome <= treaty.IncomeLevel1Upperbound {
		cessionPct = treaty.IncomeLevel1CededProportion
		tierName = "income_level1"
	} else {
		// Income falls outside all configured tiers
		return CessionResult{
			CessionBasis:    "income_outside_range",
			RetentionAmount: annualIncome,
			CededAmount:     0,
			CessionPct:      0,
		}
	}

	ceded := annualIncome * (cessionPct / 100)
	return CessionResult{
		CessionBasis:    tierName,
		RetentionAmount: annualIncome - ceded,
		CededAmount:     ceded,
		CessionPct:      cessionPct,
	}
}

// CalculateClaimCession computes the ceded portion of a claim for XL or proportional treaties.
// For proportional/quota_share treaties the three-tier sum-assured structure is applied to the
// claim amount so that the reinsurer's share mirrors the cession proportion that was in effect
// when the risk was accepted.
func CalculateClaimCession(claimAmount float64, treaty models.ReinsuranceTreaty) (ceded, retention float64, belowRetention bool) {
	switch treaty.TreatyType {
	case "xl_risk", "xl_event":
		if claimAmount <= treaty.XLRetention {
			return 0, treaty.XLRetention, true
		}
		cededAmt := claimAmount - treaty.XLRetention
		if treaty.XLLimit > 0 && cededAmt > treaty.XLLimit {
			cededAmt = treaty.XLLimit
		}
		return cededAmt, treaty.XLRetention, false

	default: // proportional, quota_share, surplus
		// Apply three-tier cession structure if configured
		if treaty.Level1Upperbound > 0 && treaty.Level1CededProportion > 0 {
			res := calculateTieredCession(claimAmount, treaty)
			return res.CededAmount, res.RetentionAmount, res.CededAmount == 0
		}
		// Surplus treaties use retention amount
		if treaty.TreatyType == "surplus" && treaty.RetentionAmount > 0 {
			if claimAmount <= treaty.RetentionAmount {
				return 0, claimAmount, true
			}
			surplus := claimAmount - treaty.RetentionAmount
			maxCession := treaty.RetentionAmount * float64(treaty.SurplusLines)
			if surplus > maxCession {
				surplus = maxCession
			}
			return surplus, claimAmount - surplus, false
		}
		// Flat retention percentage fallback
		if treaty.RetentionPercentage > 0 {
			cededAmt := claimAmount * (1 - treaty.RetentionPercentage/100)
			return cededAmt, claimAmount - cededAmt, false
		}
		return 0, claimAmount, true
	}
}

// generateBPR creates a unique Bordereaux Processing Reference e.g. BPR-202503-001
func generateBPR() (string, error) {
	now := time.Now()
	datePart := fmt.Sprintf("%d%02d", now.Year(), now.Month())
	for i := 0; i < 10; i++ {
		suffix := fmt.Sprintf("%03d", rand.Intn(999)+1)
		candidate := "BPR-" + datePart + "-" + suffix
		var count int64
		DB.Model(&models.RIBordereauxRun{}).Where("bpr = ?", candidate).Count(&count)
		if count == 0 {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("could not generate unique BPR")
}

// generateRIRunID creates a unique run ID e.g. RIBRD-202503-001
func generateRIRunID(runType string) (string, error) {
	prefix := "RIBRD"
	if runType == "claims_run" {
		prefix = "RICL"
	}
	now := time.Now()
	datePart := fmt.Sprintf("%d%02d", now.Year(), now.Month())
	for i := 0; i < 5; i++ {
		suffix := fmt.Sprintf("%03d", rand.Intn(999)+1)
		candidate := prefix + "-" + datePart + "-" + suffix
		var count int64
		DB.Model(&models.RIBordereauxRun{}).Where("run_id = ?", candidate).Count(&count)
		if count == 0 {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("could not generate unique run ID")
}

// riResolvePeriod returns (periodStart, periodEnd, periodLabel) strings for a given month/year
func riResolvePeriod(month, year int) (string, string, string) {
	if month == 0 {
		month = int(time.Now().Month())
	}
	if year == 0 {
		year = time.Now().Year()
	}
	periodStart := fmt.Sprintf("%d-%02d-01", year, month)
	lastDay := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
	periodEnd := fmt.Sprintf("%d-%02d-%02d", year, month, lastDay)
	periodLabel := fmt.Sprintf("%s %d", time.Month(month).String(), year)
	return periodStart, periodEnd, periodLabel
}

// memberGlaSumAssured estimates a member's GLA sum assured from salary × multiple
func memberGlaSumAssured(m models.GPricingMemberDataInForce) float64 {
	if m.Benefits.GlaEnabled && m.Benefits.GlaMultiple > 0 {
		return m.AnnualSalary * m.Benefits.GlaMultiple
	}
	// Fallback: use annual salary as proxy
	return m.AnnualSalary
}

// resolveSchemeIDs returns the effective scheme IDs for a bordereaux run.
// If the caller supplied an explicit list, that list is returned as-is.
// If the list is empty, all schemes linked to the treaty are used instead.
func resolveSchemeIDs(treatyID int, requested []int) ([]int, error) {
	if len(requested) > 0 {
		return requested, nil
	}
	links, err := GetTreatySchemeLinks(treatyID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch treaty scheme links: %w", err)
	}
	if len(links) == 0 {
		return nil, fmt.Errorf("no schemes are linked to treaty %d; link at least one scheme before generating", treatyID)
	}
	ids := make([]int, len(links))
	for i, l := range links {
		ids[i] = l.SchemeID
	}
	return ids, nil
}

// GenerateRIMemberBordereaux creates an RIBordereauxRun of type member_census
func GenerateRIMemberBordereaux(req models.GenerateRIBordereauxRequest, user models.AppUser) (models.RIBordereauxRun, error) {
	treaty, err := GetTreatyByID(req.TreatyID)
	if err != nil {
		return models.RIBordereauxRun{}, fmt.Errorf("treaty not found: %w", err)
	}

	schemeIDs, err := resolveSchemeIDs(req.TreatyID, req.SchemeIDs)
	if err != nil {
		return models.RIBordereauxRun{}, err
	}

	runID, err := generateRIRunID("member_census")
	if err != nil {
		return models.RIBordereauxRun{}, err
	}
	bpr, err := generateBPR()
	if err != nil {
		return models.RIBordereauxRun{}, err
	}

	periodStart, periodEnd, periodLabel := riResolvePeriod(req.Month, req.Year)
	periodStartTime, _ := time.Parse("2006-01-02", periodStart)

	schemeIDStrs := make([]string, len(schemeIDs))
	for i, id := range schemeIDs {
		schemeIDStrs[i] = strconv.Itoa(id)
	}

	run := models.RIBordereauxRun{
		RunID:         runID,
		TreatyID:      treaty.ID,
		TreatyNumber:  treaty.TreatyNumber,
		ReinsurerName: treaty.ReinsurerName,
		PeriodStart:   periodStart,
		PeriodEnd:     periodEnd,
		PeriodLabel:   periodLabel,
		Type:          "member_census",
		SchemeIDs:     "[" + strings.Join(schemeIDStrs, ",") + "]",
		Status:        "generated",
		GeneratedBy:   user.UserName,
		RunVersion:    1,
		BPR:           bpr,
	}

	var totalLives, totalCededLives int
	var totalGross, totalCeded, totalRetained float64

	for _, schemeID := range schemeIDs {
		var members []models.GPricingMemberDataInForce
		if err := DB.Where("scheme_id = ?", schemeID).Find(&members).Error; err != nil {
			appLog.WithField("scheme_id", schemeID).Warn("failed to fetch members for RI bordereaux")
			continue
		}

		var scheme models.GroupScheme
		DB.First(&scheme, schemeID)
		schemeName := scheme.Name

		for _, m := range members {
			sumAssured := memberGlaSumAssured(m)
			cr := CalculateMemberCession(sumAssured, treaty)

			changeType := "nil"
			if !m.EntryDate.IsZero() && m.EntryDate.After(periodStartTime) {
				changeType = "new"
			}
			memberStatus := "in_force"
			if changeType == "new" {
				memberStatus = "new_entrant"
			}

			// Simplified monthly premium estimate: sum_assured * 0.1% per month
			premium := sumAssured * 0.001
			cededPremium := premium * (cr.CessionPct / 100)
			retainedPremium := premium - cededPremium

			row := models.RIBordereauxMemberRow{
				RunID:                    runID,
				SchemeID:                 schemeID,
				SchemeName:               schemeName,
				MemberIDNumber:           m.MemberIdNumber,
				MemberName:               m.MemberName,
				DateOfBirth:              m.DateOfBirth.Format("2006-01-02"),
				Age:                      m.Year,
				Gender:                   m.Gender,
				EntryDate:                m.EntryDate.Format("2006-01-02"),
				BenefitCode:              m.SchemeCategory,
				BenefitName:              m.SchemeCategory,
				SumAssured:               sumAssured,
				AnnualSalary:             m.AnnualSalary,
				GrossPremium:             premium,
				CededPremium:             cededPremium,
				RetainedPremium:          retainedPremium,
				RetentionAmount:          cr.RetentionAmount,
				CededAmount:              cr.CededAmount,
				CessionBasis:             cr.CessionBasis,
				ChangeType:               changeType,
				MemberStatus:             memberStatus,
				SanctionsScreeningStatus: "cleared",
				PolicyType:               "group",
				ExchangeRate:             1.0,
				TreatySection:            treaty.TreatyNumber,
			}
			DB.Create(&row)

			totalLives++
			if cr.CededAmount > 0 {
				totalCededLives++
			}
			totalGross += premium
			totalCeded += cededPremium
			totalRetained += retainedPremium
		}
	}

	run.TotalLives = totalLives
	run.TotalCededLives = totalCededLives
	run.GrossPremium = totalGross
	run.CededPremium = totalCeded
	run.RetainedPremium = totalRetained

	if err := DB.Create(&run).Error; err != nil {
		return run, fmt.Errorf("failed to create RI run: %w", err)
	}
	return run, nil
}

// GenerateRIClaimsBordereaux creates an RIBordereauxRun of type claims_run
func GenerateRIClaimsBordereaux(req models.GenerateRIBordereauxRequest, user models.AppUser) (models.RIBordereauxRun, error) {
	treaty, err := GetTreatyByID(req.TreatyID)
	if err != nil {
		return models.RIBordereauxRun{}, fmt.Errorf("treaty not found: %w", err)
	}

	schemeIDs, err := resolveSchemeIDs(req.TreatyID, req.SchemeIDs)
	if err != nil {
		return models.RIBordereauxRun{}, err
	}

	runID, err := generateRIRunID("claims_run")
	if err != nil {
		return models.RIBordereauxRun{}, err
	}
	bpr, err := generateBPR()
	if err != nil {
		return models.RIBordereauxRun{}, err
	}

	periodStart, periodEnd, periodLabel := riResolvePeriod(req.Month, req.Year)

	schemeIDStrs := make([]string, len(schemeIDs))
	for i, id := range schemeIDs {
		schemeIDStrs[i] = strconv.Itoa(id)
	}

	run := models.RIBordereauxRun{
		RunID:         runID,
		TreatyID:      treaty.ID,
		TreatyNumber:  treaty.TreatyNumber,
		ReinsurerName: treaty.ReinsurerName,
		PeriodStart:   periodStart,
		PeriodEnd:     periodEnd,
		PeriodLabel:   periodLabel,
		Type:          "claims_run",
		SchemeIDs:     "[" + strings.Join(schemeIDStrs, ",") + "]",
		Status:        "generated",
		GeneratedBy:   user.UserName,
		RunVersion:    1,
		BPR:           bpr,
	}

	var totalGrossClaims, totalCededClaims float64

	for _, schemeID := range schemeIDs {
		var claims []models.GroupSchemeClaim
		if err := DB.Where("scheme_id = ? AND date_registered >= ? AND date_registered <= ?",
			schemeID, periodStart, periodEnd).Find(&claims).Error; err != nil {
			appLog.WithField("scheme_id", schemeID).Warn("failed to fetch claims for RI bordereaux")
			continue
		}

		var scheme models.GroupScheme
		DB.First(&scheme, schemeID)

		for _, cl := range claims {
			ceded, retention, below := CalculateClaimCession(cl.ClaimAmount, treaty)

			// Derive paid vs outstanding split from claim status
			var grossPaid, grossOutstanding float64
			if cl.Status == "paid" || cl.Status == "approved" {
				grossPaid = cl.ClaimAmount
			} else {
				grossOutstanding = cl.ClaimAmount
			}

			largeLoss := treaty.LargeClaimsThreshold > 0 && cl.ClaimAmount >= treaty.LargeClaimsThreshold

			row := models.RIBordereauxClaimsRow{
				RunID:                   runID,
				SchemeID:                schemeID,
				SchemeName:              scheme.Name,
				ClaimNumber:             cl.ClaimNumber,
				MemberIDNumber:          cl.MemberIDNumber,
				MemberName:              cl.MemberName,
				DateOfEvent:             cl.DateOfEvent,
				DateNotified:            cl.DateNotified,
				BenefitCode:             cl.BenefitCode,
				GrossClaimAmount:        cl.ClaimAmount,
				ExcessRetention:         retention,
				CededClaimAmount:        ceded,
				ClaimStatus:             cl.Status,
				IsBelowRetention:        below,
				CauseOfLoss:             cl.BenefitCode,
				GrossPaidLosses:         grossPaid,
				GrossOutstandingReserve: grossOutstanding,
				IBNRFlag:                false,
				LargeLossFlag:           largeLoss,
			}
			DB.Create(&row)
			totalGrossClaims += cl.ClaimAmount
			totalCededClaims += ceded
		}
	}

	run.GrossClaimsIncurred = totalGrossClaims
	run.CededClaimsIncurred = totalCededClaims

	if err := DB.Create(&run).Error; err != nil {
		return run, fmt.Errorf("failed to create RI claims run: %w", err)
	}
	return run, nil
}

// GetRIBordereauxRuns returns runs filtered by optional treaty, type and status
func GetRIBordereauxRuns(treatyID int, runType, status string) ([]models.RIBordereauxRun, error) {
	var runs []models.RIBordereauxRun
	q := DB.Order("created_at DESC")
	if treatyID > 0 {
		q = q.Where("treaty_id = ?", treatyID)
	}
	if runType != "" {
		q = q.Where("type = ?", runType)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Find(&runs).Error; err != nil {
		return nil, err
	}
	return runs, nil
}

// GetRIBordereauxRunByID returns a single run by RunID string
func GetRIBordereauxRunByID(runID string) (models.RIBordereauxRun, error) {
	var run models.RIBordereauxRun
	if err := DB.Where("run_id = ?", runID).First(&run).Error; err != nil {
		return run, fmt.Errorf("run %s not found: %w", runID, err)
	}
	return run, nil
}

// GetRIBordereauxMemberRows returns persisted member rows for a run
func GetRIBordereauxMemberRows(runID string) ([]models.RIBordereauxMemberRow, error) {
	var rows []models.RIBordereauxMemberRow
	if err := DB.Where("run_id = ?", runID).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetRIBordereauxClaimsRows returns persisted claims rows for a run
func GetRIBordereauxClaimsRows(runID string) ([]models.RIBordereauxClaimsRow, error) {
	var rows []models.RIBordereauxClaimsRow
	if err := DB.Where("run_id = ?", runID).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// SubmitRIBordereaux marks a generated run as submitted
func SubmitRIBordereaux(req models.SubmitRIBordereauxRequest, user models.AppUser) (models.RIBordereauxRun, error) {
	var run models.RIBordereauxRun
	if err := DB.Where("run_id = ?", req.RunID).First(&run).Error; err != nil {
		return run, fmt.Errorf("run not found: %w", err)
	}
	if err := ValidateRIBordereauxRunTransition(run.Status, StatusRIRunSubmitted); err != nil {
		return run, err
	}
	now := time.Now()
	run.Status = StatusRIRunSubmitted
	run.SubmittedAt = &now
	run.SubmittedBy = user.UserName
	run.UpdatedAt = now
	if err := DB.Save(&run).Error; err != nil {
		return run, fmt.Errorf("failed to submit RI bordereaux: %w", err)
	}
	go NotifyRIBordereauxSubmitted(run, user)
	return run, nil
}

// AcknowledgeRIBordereaux marks a submitted run as acknowledged by the reinsurer
func AcknowledgeRIBordereaux(runID string, user models.AppUser) (models.RIBordereauxRun, error) {
	var run models.RIBordereauxRun
	if err := DB.Where("run_id = ?", runID).First(&run).Error; err != nil {
		return run, fmt.Errorf("run not found: %w", err)
	}
	if err := ValidateRIBordereauxRunTransition(run.Status, StatusRIRunAcknowledged); err != nil {
		return run, err
	}
	now := time.Now()
	run.Status = StatusRIRunAcknowledged
	run.AcknowledgedAt = &now
	run.AcknowledgedBy = user.UserName
	run.UpdatedAt = now
	if err := DB.Save(&run).Error; err != nil {
		return run, fmt.Errorf("failed to acknowledge RI bordereaux: %w", err)
	}
	go NotifyRIBordereauxAcknowledged(run, user)
	return run, nil
}

// AcknowledgeRIBordereauxReceipt logs formal receipt of reinsurer confirmation without changing status
func AcknowledgeRIBordereauxReceipt(runID string, req models.AcknowledgeReceiptRequest, user models.AppUser) (models.RIBordereauxRun, error) {
	var run models.RIBordereauxRun
	if err := DB.Where("run_id = ?", runID).First(&run).Error; err != nil {
		return run, fmt.Errorf("run not found: %w", err)
	}
	// Receipt acknowledgement is a no-op if the run is already acknowledged;
	// otherwise require a legal transition into acknowledged.
	if run.Status != StatusRIRunAcknowledged {
		if err := ValidateRIBordereauxRunTransition(run.Status, StatusRIRunAcknowledged); err != nil {
			return run, err
		}
	}
	receivedDate := req.ReceivedDate
	if receivedDate == "" {
		receivedDate = time.Now().Format("2006-01-02")
	}
	now := time.Now()
	run.ReceivedDate = receivedDate
	run.AcknowledgedBy = user.UserName
	run.AcknowledgedAt = &now
	run.Status = StatusRIRunAcknowledged
	run.UpdatedAt = now
	if err := DB.Save(&run).Error; err != nil {
		return run, fmt.Errorf("failed to log receipt acknowledgement: %w", err)
	}
	return run, nil
}

// AmendRIBordereaux creates a new versioned run cloned from an existing submitted/acknowledged run
func AmendRIBordereaux(runID string, req models.AmendRIBordereauxRequest, user models.AppUser) (models.RIBordereauxRun, error) {
	parent, err := GetRIBordereauxRunByID(runID)
	if err != nil {
		return models.RIBordereauxRun{}, fmt.Errorf("parent run not found: %w", err)
	}
	if parent.Status != "submitted" && parent.Status != "acknowledged" {
		return models.RIBordereauxRun{}, fmt.Errorf("only submitted or acknowledged runs can be amended (current: %s)", parent.Status)
	}

	newRunID, err := generateRIRunID(parent.Type)
	if err != nil {
		return models.RIBordereauxRun{}, err
	}
	newBPR, err := generateBPR()
	if err != nil {
		return models.RIBordereauxRun{}, err
	}

	parentID := uint(parent.ID)
	newRun := models.RIBordereauxRun{
		RunID:               newRunID,
		TreatyID:            parent.TreatyID,
		TreatyNumber:        parent.TreatyNumber,
		ReinsurerName:       parent.ReinsurerName,
		PeriodStart:         parent.PeriodStart,
		PeriodEnd:           parent.PeriodEnd,
		PeriodLabel:         parent.PeriodLabel,
		Type:                parent.Type,
		SchemeIDs:           parent.SchemeIDs,
		Status:              "generated",
		GeneratedBy:         user.UserName,
		RunVersion:          parent.RunVersion + 1,
		AmendmentNotes:      req.AmendmentNotes,
		ParentRunID:         &parentID,
		BPR:                 newBPR,
		TotalLives:          parent.TotalLives,
		TotalCededLives:     parent.TotalCededLives,
		GrossPremium:        parent.GrossPremium,
		CededPremium:        parent.CededPremium,
		RetainedPremium:     parent.RetainedPremium,
		GrossClaimsIncurred: parent.GrossClaimsIncurred,
		CededClaimsIncurred: parent.CededClaimsIncurred,
	}
	if err := DB.Create(&newRun).Error; err != nil {
		return newRun, fmt.Errorf("failed to create amendment run: %w", err)
	}

	// Copy rows from parent into the new run
	if parent.Type == "member_census" {
		var rows []models.RIBordereauxMemberRow
		DB.Where("run_id = ?", runID).Find(&rows)
		for _, r := range rows {
			r.ID = 0
			r.RunID = newRunID
			DB.Create(&r)
		}
	} else {
		var rows []models.RIBordereauxClaimsRow
		DB.Where("run_id = ?", runID).Find(&rows)
		for _, r := range rows {
			r.ID = 0
			r.RunID = newRunID
			DB.Create(&r)
		}
	}

	return newRun, nil
}

// GetRIBordereauxStats returns a summary of RI bordereaux runs
func GetRIBordereauxStats(treatyID int) (models.RIBordereauxStats, error) {
	var runs []models.RIBordereauxRun
	q := DB.Model(&models.RIBordereauxRun{})
	if treatyID > 0 {
		q = q.Where("treaty_id = ?", treatyID)
	}
	if err := q.Find(&runs).Error; err != nil {
		return models.RIBordereauxStats{}, err
	}
	stats := models.RIBordereauxStats{TotalRuns: len(runs)}
	for _, r := range runs {
		switch r.Type {
		case "member_census":
			stats.MemberCensus++
		case "claims_run":
			stats.ClaimsRuns++
		}
		if r.Status == "submitted" || r.Status == "acknowledged" {
			stats.Submitted++
		}
		stats.TotalCededPremium += r.CededPremium
	}
	return stats, nil
}

// MonitorLargeClaims scans all claims for treaty-linked schemes and creates LargeClaimNotice records
func MonitorLargeClaims(req models.MonitorLargeClaimsRequest) (int, error) {
	treaty, err := GetTreatyByID(req.TreatyID)
	if err != nil {
		return 0, err
	}
	if treaty.LargeClaimsThreshold <= 0 {
		return 0, fmt.Errorf("treaty has no large claims threshold configured")
	}

	links, err := GetTreatySchemeLinks(treaty.ID)
	if err != nil {
		return 0, err
	}

	created := 0
	today := time.Now().Format("2006-01-02")

	for _, link := range links {
		var claims []models.GroupSchemeClaim
		if err := DB.Where("scheme_id = ? AND claim_amount >= ?", link.SchemeID, treaty.LargeClaimsThreshold).
			Find(&claims).Error; err != nil {
			continue
		}
		for _, cl := range claims {
			var existing int64
			DB.Model(&models.LargeClaimNotice{}).
				Where("treaty_id = ? AND claim_id = ?", treaty.ID, cl.ID).
				Count(&existing)
			if existing > 0 {
				continue
			}

			ceded, retention, _ := CalculateClaimCession(cl.ClaimAmount, treaty)
			dueDateParsed, parseErr := time.Parse("2006-01-02", cl.DateRegistered)
			if parseErr != nil {
				dueDateParsed = time.Now()
			}
			dueDate := dueDateParsed.AddDate(0, 0, treaty.ClaimsNotificationDays).Format("2006-01-02")

			lateFlag := today > dueDate
			status := "pending"
			if lateFlag {
				status = "late"
			}

			notice := models.LargeClaimNotice{
				TreatyID:             treaty.ID,
				TreatyNumber:         treaty.TreatyNumber,
				ClaimID:              cl.ID,
				ClaimNumber:          cl.ClaimNumber,
				SchemeID:             link.SchemeID,
				SchemeName:           cl.SchemeName,
				ReinsurerName:        treaty.ReinsurerName,
				EventDate:            cl.DateOfEvent,
				NotifiedDate:         cl.DateNotified,
				DueDate:              dueDate,
				BenefitCode:          cl.BenefitCode,
				GrossClaimAmount:     cl.ClaimAmount,
				ExcessAmount:         retention,
				EstimatedCededAmount: ceded,
				Status:               status,
				LateFlag:             lateFlag,
			}
			if err := DB.Create(&notice).Error; err == nil {
				created++
			}
		}
	}
	return created, nil
}

// GetLargeClaimNotices returns notices filtered by treaty, scheme and status; auto-marks overdue
func GetLargeClaimNotices(treatyID, schemeID int, status string) ([]models.LargeClaimNotice, error) {
	today := time.Now().Format("2006-01-02")
	// Auto-mark overdue pending notices
	DB.Model(&models.LargeClaimNotice{}).
		Where("status = 'pending' AND due_date < ?", today).
		Updates(map[string]interface{}{"status": "late", "late_flag": true})

	var notices []models.LargeClaimNotice
	q := DB.Order("created_at DESC")
	if treatyID > 0 {
		q = q.Where("treaty_id = ?", treatyID)
	}
	if schemeID > 0 {
		q = q.Where("scheme_id = ?", schemeID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Find(&notices).Error; err != nil {
		return nil, err
	}
	return notices, nil
}

// UpdateLargeClaimNotice updates status, notes, and timestamps on a notice
func UpdateLargeClaimNotice(id int, req models.UpdateLargeClaimNoticeRequest, user models.AppUser) (models.LargeClaimNotice, error) {
	var notice models.LargeClaimNotice
	if err := DB.First(&notice, id).Error; err != nil {
		return notice, fmt.Errorf("notice %d not found: %w", id, err)
	}
	if req.Status != "" {
		if !IsKnownLargeClaimNoticeStatus(req.Status) {
			return notice, fmt.Errorf("unknown large-claim notice status %q", req.Status)
		}
		if err := ValidateLargeClaimNoticeTransition(notice.Status, req.Status); err != nil {
			return notice, err
		}
		notice.Status = req.Status
		if req.Status == StatusNoticeSent && notice.SentAt == nil {
			now := time.Now()
			notice.SentAt = &now
		}
		if req.Status == StatusNoticeAcknowledged && notice.AcknowledgedAt == nil {
			now := time.Now()
			notice.AcknowledgedAt = &now
		}
	}
	if req.SentAt != nil {
		notice.SentAt = req.SentAt
	}
	if req.AcknowledgedAt != nil {
		notice.AcknowledgedAt = req.AcknowledgedAt
	}
	if req.QueryDetails != "" {
		notice.QueryDetails = req.QueryDetails
	}
	if req.ResponseNotes != "" {
		notice.ResponseNotes = req.ResponseNotes
	}
	notice.UpdatedAt = time.Now()
	if err := DB.Save(&notice).Error; err != nil {
		return notice, fmt.Errorf("failed to update notice: %w", err)
	}
	return notice, nil
}

// LargeClaimResponseRequest is shared by the accept / reject / query endpoints.
// Empty fields are ignored so callers only pass what's relevant.
type LargeClaimResponseRequest struct {
	Notes          string `json:"notes"`
	AcceptedAmount float64 `json:"accepted_amount"`
	QueryDetails   string `json:"query_details"`
	Reason         string `json:"reason"`
}

// AcceptLargeClaimNotice records a reinsurer's acceptance of a large-claim
// cession. Sets ResponseStatus=accepted and flips Status to acknowledged so
// the ceding-side dashboards reflect that the notice has been responded to.
// Idempotent in the sense that re-accepting just refreshes RespondedAt / notes.
func AcceptLargeClaimNotice(id int, req LargeClaimResponseRequest, user models.AppUser) (models.LargeClaimNotice, error) {
	var notice models.LargeClaimNotice
	if err := DB.First(&notice, id).Error; err != nil {
		return notice, fmt.Errorf("notice %d not found: %w", id, err)
	}
	if notice.Status != StatusNoticeAcknowledged {
		if err := ValidateLargeClaimNoticeTransition(notice.Status, StatusNoticeAcknowledged); err != nil {
			return notice, err
		}
	}
	before := notice
	now := time.Now()

	notice.ResponseStatus = ResponseNoticeAccepted
	notice.RespondedAt = &now
	notice.RespondedBy = user.UserName
	notice.Status = StatusNoticeAcknowledged
	if notice.AcknowledgedAt == nil {
		notice.AcknowledgedAt = &now
	}
	if req.Notes != "" {
		notice.ResponseNotes = appendResponseNote(notice.ResponseNotes, user.UserName, "accepted", req.Notes)
	}
	if req.AcceptedAmount > 0 {
		notice.EstimatedCededAmount = req.AcceptedAmount
	}
	notice.UpdatedAt = now
	if err := DB.Save(&notice).Error; err != nil {
		return notice, fmt.Errorf("failed to accept notice: %w", err)
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "large_claim_notices",
		EntityID:  fmt.Sprintf("%d", notice.ID),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, notice)
	return notice, nil
}

// RejectLargeClaimNotice records a reinsurer's rejection of a cession. The
// reason is mandatory (validated at the controller layer); it is appended to
// ResponseNotes for audit trail alongside who rejected and when.
func RejectLargeClaimNotice(id int, req LargeClaimResponseRequest, user models.AppUser) (models.LargeClaimNotice, error) {
	var notice models.LargeClaimNotice
	if err := DB.First(&notice, id).Error; err != nil {
		return notice, fmt.Errorf("notice %d not found: %w", id, err)
	}
	if notice.Status != StatusNoticeAcknowledged {
		if err := ValidateLargeClaimNoticeTransition(notice.Status, StatusNoticeAcknowledged); err != nil {
			return notice, err
		}
	}
	before := notice
	now := time.Now()

	notice.ResponseStatus = ResponseNoticeRejected
	notice.RespondedAt = &now
	notice.RespondedBy = user.UserName
	notice.Status = StatusNoticeAcknowledged
	if notice.AcknowledgedAt == nil {
		notice.AcknowledgedAt = &now
	}
	if req.Reason != "" {
		notice.ResponseNotes = appendResponseNote(notice.ResponseNotes, user.UserName, "rejected", req.Reason)
	}
	notice.UpdatedAt = now
	if err := DB.Save(&notice).Error; err != nil {
		return notice, fmt.Errorf("failed to reject notice: %w", err)
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "large_claim_notices",
		EntityID:  fmt.Sprintf("%d", notice.ID),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, notice)
	return notice, nil
}

// QueryLargeClaimNotice records a reinsurer's query on the notice. Does not
// set ResponseStatus — the reinsurer has not decided yet; Status moves to
// "queried" so the ceding team knows to reply.
func QueryLargeClaimNotice(id int, req LargeClaimResponseRequest, user models.AppUser) (models.LargeClaimNotice, error) {
	var notice models.LargeClaimNotice
	if err := DB.First(&notice, id).Error; err != nil {
		return notice, fmt.Errorf("notice %d not found: %w", id, err)
	}
	if notice.Status != StatusNoticeQueried {
		if err := ValidateLargeClaimNoticeTransition(notice.Status, StatusNoticeQueried); err != nil {
			return notice, err
		}
	}
	before := notice
	now := time.Now()

	notice.Status = StatusNoticeQueried
	notice.RespondedAt = &now
	notice.RespondedBy = user.UserName
	if req.QueryDetails != "" {
		notice.QueryDetails = req.QueryDetails
	}
	notice.UpdatedAt = now
	if err := DB.Save(&notice).Error; err != nil {
		return notice, fmt.Errorf("failed to record query: %w", err)
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "large_claim_notices",
		EntityID:  fmt.Sprintf("%d", notice.ID),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, notice)
	return notice, nil
}

// appendResponseNote builds a timestamped entry for ResponseNotes so the full
// history is preserved across multiple interactions rather than overwritten.
func appendResponseNote(existing, user, action, body string) string {
	entry := fmt.Sprintf("[%s] %s %s: %s",
		time.Now().Format("2006-01-02 15:04"), user, action, body)
	if existing == "" {
		return entry
	}
	return existing + "\n" + entry
}

// GetLargeClaimStats returns count summary of large claim notices
func GetLargeClaimStats(treatyID int) (models.LargeClaimStats, error) {
	var notices []models.LargeClaimNotice
	q := DB.Model(&models.LargeClaimNotice{})
	if treatyID > 0 {
		q = q.Where("treaty_id = ?", treatyID)
	}
	if err := q.Find(&notices).Error; err != nil {
		return models.LargeClaimStats{}, err
	}
	stats := models.LargeClaimStats{Total: len(notices)}
	for _, n := range notices {
		switch n.Status {
		case "pending":
			stats.Pending++
		case "sent":
			stats.Sent++
		case "acknowledged":
			stats.Acknowledged++
		case "late":
			stats.Late++
		case "queried":
			stats.Queried++
		}
	}
	return stats, nil
}

// GetCatastropheClaimsRows returns claims rows that have a non-empty catastrophe_event_code,
// optionally filtered by code, treaty, and period.
func GetCatastropheClaimsRows(catEventCode string, treatyID int, periodFrom, periodTo string) ([]models.RIBordereauxClaimsRow, error) {
	var rows []models.RIBordereauxClaimsRow
	q := DB.Where("catastrophe_event_code != ''")
	if catEventCode != "" {
		q = q.Where("catastrophe_event_code = ?", catEventCode)
	}
	if treatyID > 0 {
		var runIDs []string
		DB.Model(&models.RIBordereauxRun{}).Where("treaty_id = ?", treatyID).Pluck("id", &runIDs)
		if len(runIDs) > 0 {
			q = q.Where("run_id IN ?", runIDs)
		} else {
			return rows, nil // treaty has no runs, return empty
		}
	}
	if periodFrom != "" {
		var runIDs []string
		DB.Model(&models.RIBordereauxRun{}).Where("period_end >= ?", periodFrom).Pluck("id", &runIDs)
		if len(runIDs) > 0 {
			q = q.Where("run_id IN ?", runIDs)
		}
	}
	if periodTo != "" {
		var runIDs []string
		DB.Model(&models.RIBordereauxRun{}).Where("period_start <= ?", periodTo).Pluck("id", &runIDs)
		if len(runIDs) > 0 {
			q = q.Where("run_id IN ?", runIDs)
		}
	}
	if err := q.Order("catastrophe_event_code, date_of_event").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
