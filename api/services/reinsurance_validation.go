package services

import (
	"api/models"
	"fmt"
	"math"
	"strings"
	"time"
)

// ValidateRIBordereaux runs the 3-level validation pipeline against a run and persists results.
// Level 1 — Structural (completeness, duplicate run detection, treaty metadata)
// Level 2 — Data integrity (arithmetic, date logic, exchange rate range)
// Level 3 — Business rules (treaty limits, sanctions, large-loss recovery)
func ValidateRIBordereaux(runID string) (models.ValidationSummary, error) {
	run, err := GetRIBordereauxRunByID(runID)
	if err != nil {
		return models.ValidationSummary{}, fmt.Errorf("run not found: %w", err)
	}
	if run.Status == "submitted" || run.Status == "acknowledged" || run.Status == "settled" {
		return models.ValidationSummary{}, fmt.Errorf("cannot validate a run with status '%s'", run.Status)
	}

	// Mark as validating
	DB.Model(&run).Update("status", "validating")

	// Clear any prior results for a re-run
	DB.Where("run_id = ?", runID).Delete(&models.RIValidationResult{})

	treaty, err := GetTreatyByID(run.TreatyID)
	if err != nil {
		return models.ValidationSummary{}, fmt.Errorf("treaty not found: %w", err)
	}

	var results []models.RIValidationResult

	// ── Level 1: Structural ───────────────────────────────────────────────────
	results = append(results, validateL1Run(run, treaty)...)

	if run.Type == "member_census" {
		memberRows, _ := GetRIBordereauxMemberRows(runID)
		results = append(results, validateL1Members(runID, memberRows)...)
		results = append(results, validateL2Members(runID, memberRows)...)
		results = append(results, validateL3Members(runID, memberRows, treaty)...)
	} else {
		claimsRows, _ := GetRIBordereauxClaimsRows(runID)
		results = append(results, validateL1Claims(runID, claimsRows)...)
		results = append(results, validateL2Claims(runID, claimsRows)...)
		results = append(results, validateL3Claims(runID, claimsRows, treaty)...)
	}

	// Persist all findings
	for i := range results {
		DB.Create(&results[i])
	}

	// Tally
	summary := models.ValidationSummary{RunID: runID, Results: results}
	for _, r := range results {
		summary.Total++
		switch r.Severity {
		case "critical":
			summary.Critical++
		case "major":
			summary.Major++
		case "minor":
			summary.Minor++
		}
	}

	// Final status: any critical → validation_failed, else validated
	if summary.Critical > 0 {
		summary.Status = "validation_failed"
	} else {
		summary.Status = "validated"
	}
	DB.Model(&run).Update("status", summary.Status)

	return summary, nil
}

// GetRIValidationResults returns persisted validation results for a run
func GetRIValidationResults(runID string) (models.ValidationSummary, error) {
	var results []models.RIValidationResult
	if err := DB.Where("run_id = ?", runID).Order("level, severity, row_index").Find(&results).Error; err != nil {
		return models.ValidationSummary{}, err
	}
	run, err := GetRIBordereauxRunByID(runID)
	if err != nil {
		return models.ValidationSummary{}, err
	}
	summary := models.ValidationSummary{RunID: runID, Status: run.Status, Results: results}
	for _, r := range results {
		summary.Total++
		switch r.Severity {
		case "critical":
			summary.Critical++
		case "major":
			summary.Major++
		case "minor":
			summary.Minor++
		}
	}
	return summary, nil
}

// ── Level 1: Structural ───────────────────────────────────────────────────────

func validateL1Run(run models.RIBordereauxRun, treaty models.ReinsuranceTreaty) []models.RIValidationResult {
	var out []models.RIValidationResult

	// Treaty must be active
	if treaty.Status != "active" {
		out = append(out, finding(run.RunID, "run", 0, 1, "critical", "treaty_status",
			"L1-001", fmt.Sprintf("Treaty status is '%s'; only active treaties may be submitted", treaty.Status)))
	}

	// Duplicate run check: same treaty + period + type, status not draft, different runID, not an amendment of this run
	var dupCount int64
	DB.Model(&models.RIBordereauxRun{}).
		Where("treaty_id = ? AND period_start = ? AND period_end = ? AND type = ? AND status NOT IN ? AND run_id != ?",
			run.TreatyID, run.PeriodStart, run.PeriodEnd, run.Type,
			[]string{"draft", "validating"}, run.RunID).
		Count(&dupCount)
	if dupCount > 0 {
		out = append(out, finding(run.RunID, "run", 0, 1, "major", "run_id",
			"L1-002", fmt.Sprintf("A non-draft run for this treaty and period (%s) already exists", run.PeriodLabel)))
	}

	// Period completeness
	if run.PeriodStart == "" || run.PeriodEnd == "" {
		out = append(out, finding(run.RunID, "run", 0, 1, "critical", "period",
			"L1-003", "Run period start or end date is missing"))
	}

	// Currency: must be a 3-letter uppercase ISO code
	currency := treaty.Currency
	if len(currency) != 3 || currency != strings.ToUpper(currency) {
		out = append(out, finding(run.RunID, "run", 0, 1, "minor", "currency",
			"L1-004", fmt.Sprintf("Treaty currency '%s' does not appear to be a valid ISO 4217 code", currency)))
	}

	return out
}

func validateL1Members(runID string, rows []models.RIBordereauxMemberRow) []models.RIValidationResult {
	var out []models.RIValidationResult
	for i, r := range rows {
		idx := i + 1
		if strings.TrimSpace(r.MemberIDNumber) == "" {
			out = append(out, finding(runID, "member", idx, 1, "critical", "member_id_number",
				"L1-101", fmt.Sprintf("Row %d: member ID number is missing", idx)))
		}
		if strings.TrimSpace(r.MemberName) == "" {
			out = append(out, finding(runID, "member", idx, 1, "major", "member_name",
				"L1-102", fmt.Sprintf("Row %d: member name is missing", idx)))
		}
		if strings.TrimSpace(r.DateOfBirth) == "" {
			out = append(out, finding(runID, "member", idx, 1, "major", "date_of_birth",
				"L1-103", fmt.Sprintf("Row %d: date of birth is missing", idx)))
		}
		if strings.TrimSpace(r.BenefitCode) == "" {
			out = append(out, finding(runID, "member", idx, 1, "major", "benefit_code",
				"L1-104", fmt.Sprintf("Row %d: benefit code is missing", idx)))
		}
		if r.SumAssured <= 0 {
			out = append(out, finding(runID, "member", idx, 1, "critical", "sum_assured",
				"L1-105", fmt.Sprintf("Row %d: sum assured must be greater than zero", idx)))
		}
		if r.GrossPremium <= 0 {
			out = append(out, finding(runID, "member", idx, 1, "critical", "gross_premium",
				"L1-106", fmt.Sprintf("Row %d: gross premium must be greater than zero", idx)))
		}
	}
	return out
}

func validateL1Claims(runID string, rows []models.RIBordereauxClaimsRow) []models.RIValidationResult {
	var out []models.RIValidationResult
	for i, r := range rows {
		idx := i + 1
		if strings.TrimSpace(r.ClaimNumber) == "" {
			out = append(out, finding(runID, "claims", idx, 1, "critical", "claim_number",
				"L1-201", fmt.Sprintf("Row %d: claim number is missing", idx)))
		}
		if strings.TrimSpace(r.MemberIDNumber) == "" {
			out = append(out, finding(runID, "claims", idx, 1, "major", "member_id_number",
				"L1-202", fmt.Sprintf("Row %d: member ID number is missing", idx)))
		}
		if strings.TrimSpace(r.DateOfEvent) == "" {
			out = append(out, finding(runID, "claims", idx, 1, "critical", "date_of_event",
				"L1-203", fmt.Sprintf("Row %d: date of event is missing", idx)))
		}
		if strings.TrimSpace(r.BenefitCode) == "" {
			out = append(out, finding(runID, "claims", idx, 1, "major", "benefit_code",
				"L1-204", fmt.Sprintf("Row %d: benefit code is missing", idx)))
		}
		if r.GrossClaimAmount <= 0 {
			out = append(out, finding(runID, "claims", idx, 1, "critical", "gross_claim_amount",
				"L1-205", fmt.Sprintf("Row %d: gross claim amount must be greater than zero", idx)))
		}
	}
	return out
}

// ── Level 2: Data Integrity ───────────────────────────────────────────────────

func validateL2Members(runID string, rows []models.RIBordereauxMemberRow) []models.RIValidationResult {
	var out []models.RIValidationResult
	today := time.Now()

	for i, r := range rows {
		idx := i + 1

		// Premium arithmetic: CededPremium + RetainedPremium ≈ GrossPremium (tolerance 0.01)
		if r.GrossPremium > 0 {
			diff := math.Abs(r.GrossPremium - r.CededPremium - r.RetainedPremium)
			if diff > 0.01 {
				out = append(out, finding(runID, "member", idx, 2, "major", "gross_premium",
					"L2-101", fmt.Sprintf("Row %d: premium arithmetic error — gross %.2f ≠ ceded %.2f + retained %.2f (diff %.4f)",
						idx, r.GrossPremium, r.CededPremium, r.RetainedPremium, diff)))
			}
		}

		// Exchange rate must be positive and within plausible range
		if r.ExchangeRate <= 0 || r.ExchangeRate > 10000 {
			out = append(out, finding(runID, "member", idx, 2, "minor", "exchange_rate",
				"L2-102", fmt.Sprintf("Row %d: exchange rate %.4f is outside plausible range (0–10000)", idx, r.ExchangeRate)))
		}

		// Entry date must not be in the future
		if r.EntryDate != "" {
			if ed, err := time.Parse("2006-01-02", r.EntryDate); err == nil {
				if ed.After(today) {
					out = append(out, finding(runID, "member", idx, 2, "minor", "entry_date",
						"L2-103", fmt.Sprintf("Row %d: entry date %s is in the future", idx, r.EntryDate)))
				}
			}
		}

		// Date of birth must be plausible (born after 1900, not in future)
		if r.DateOfBirth != "" {
			if dob, err := time.Parse("2006-01-02", r.DateOfBirth); err == nil {
				if dob.Year() < 1900 || dob.After(today) {
					out = append(out, finding(runID, "member", idx, 2, "major", "date_of_birth",
						"L2-104", fmt.Sprintf("Row %d: date of birth %s is implausible", idx, r.DateOfBirth)))
				}
			} else {
				out = append(out, finding(runID, "member", idx, 2, "minor", "date_of_birth",
					"L2-105", fmt.Sprintf("Row %d: date of birth '%s' cannot be parsed", idx, r.DateOfBirth)))
			}
		}

		// Ceded + Retained SA should equal Sum Assured for quota-share/surplus (not XL)
		if r.CessionBasis == "quota_share" || r.CessionBasis == "surplus" {
			diff := math.Abs(r.SumAssured - r.RetentionAmount - r.CededAmount)
			if diff > 1.0 {
				out = append(out, finding(runID, "member", idx, 2, "major", "sum_assured",
					"L2-106", fmt.Sprintf("Row %d: sum assured %.2f ≠ retention %.2f + ceded %.2f (diff %.2f)",
						idx, r.SumAssured, r.RetentionAmount, r.CededAmount, diff)))
			}
		}
	}
	return out
}

func validateL2Claims(runID string, rows []models.RIBordereauxClaimsRow) []models.RIValidationResult {
	var out []models.RIValidationResult
	today := time.Now()

	for i, r := range rows {
		idx := i + 1

		// GrossPaidLosses + GrossOutstandingReserve ≈ GrossClaimAmount (tolerance 0.01)
		if r.GrossClaimAmount > 0 {
			diff := math.Abs(r.GrossClaimAmount - r.GrossPaidLosses - r.GrossOutstandingReserve)
			if diff > 0.01 {
				out = append(out, finding(runID, "claims", idx, 2, "major", "gross_claim_amount",
					"L2-201", fmt.Sprintf("Row %d: claim split error — gross %.2f ≠ paid %.2f + outstanding %.2f (diff %.4f)",
						idx, r.GrossClaimAmount, r.GrossPaidLosses, r.GrossOutstandingReserve, diff)))
			}
		}

		// CededClaimAmount must not exceed GrossClaimAmount
		if r.CededClaimAmount > r.GrossClaimAmount {
			out = append(out, finding(runID, "claims", idx, 2, "critical", "ceded_claim_amount",
				"L2-202", fmt.Sprintf("Row %d: ceded amount %.2f exceeds gross claim amount %.2f",
					idx, r.CededClaimAmount, r.GrossClaimAmount)))
		}

		// DateNotified must not precede DateOfEvent
		if r.DateOfEvent != "" && r.DateNotified != "" {
			doe, errE := time.Parse("2006-01-02", r.DateOfEvent)
			don, errN := time.Parse("2006-01-02", r.DateNotified)
			if errE == nil && errN == nil && don.Before(doe) {
				out = append(out, finding(runID, "claims", idx, 2, "major", "date_notified",
					"L2-203", fmt.Sprintf("Row %d: date notified %s is before date of event %s",
						idx, r.DateNotified, r.DateOfEvent)))
			}
		}

		// DateOfEvent must not be in the future
		if r.DateOfEvent != "" {
			if doe, err := time.Parse("2006-01-02", r.DateOfEvent); err == nil {
				if doe.After(today) {
					out = append(out, finding(runID, "claims", idx, 2, "minor", "date_of_event",
						"L2-204", fmt.Sprintf("Row %d: date of event %s is in the future", idx, r.DateOfEvent)))
				}
			}
		}

		// Recoveries must not exceed gross paid losses
		if r.Recoveries > r.GrossPaidLosses && r.GrossPaidLosses > 0 {
			out = append(out, finding(runID, "claims", idx, 2, "major", "recoveries",
				"L2-205", fmt.Sprintf("Row %d: recoveries %.2f exceed gross paid losses %.2f",
					idx, r.Recoveries, r.GrossPaidLosses)))
		}
	}
	return out
}

// ── Level 3: Business Rules ───────────────────────────────────────────────────

func validateL3Members(runID string, rows []models.RIBordereauxMemberRow, treaty models.ReinsuranceTreaty) []models.RIValidationResult {
	var out []models.RIValidationResult

	for i, r := range rows {
		idx := i + 1

		// Sanctions-flagged members are a critical stop
		if r.SanctionsScreeningStatus == "flagged" {
			out = append(out, finding(runID, "member", idx, 3, "critical", "sanctions_screening_status",
				"L3-101", fmt.Sprintf("Row %d: member %s (%s) is sanctions-flagged and must be excluded before submission",
					idx, r.MemberName, r.MemberIDNumber)))
		}

		// Sanctions-pending members warrant a major warning
		if r.SanctionsScreeningStatus == "pending" {
			out = append(out, finding(runID, "member", idx, 3, "major", "sanctions_screening_status",
				"L3-102", fmt.Sprintf("Row %d: member %s (%s) has a pending sanctions screening — resolve before submission",
					idx, r.MemberName, r.MemberIDNumber)))
		}

		// For XL treaties, all members should be below retention (no cession)
		// but a ceded amount > 0 on XL member is an error
		if (treaty.TreatyType == "xl_risk" || treaty.TreatyType == "xl_event") && r.CededAmount > 0 {
			out = append(out, finding(runID, "member", idx, 3, "minor", "ceded_amount",
				"L3-103", fmt.Sprintf("Row %d: XL treaty members should have zero ceded SA; use claims bordereaux for XL recoveries", idx)))
		}

		// For surplus treaties: SumAssured must exceed retention to have any cession
		if treaty.TreatyType == "surplus" && r.CededAmount > 0 && r.SumAssured <= treaty.RetentionAmount {
			out = append(out, finding(runID, "member", idx, 3, "major", "ceded_amount",
				"L3-104", fmt.Sprintf("Row %d: sum assured %.2f ≤ treaty retention %.2f but ceded amount is %.2f",
					idx, r.SumAssured, treaty.RetentionAmount, r.CededAmount)))
		}

		// CessionBasis must align with treaty type
		expectedBasis := cessionBasisForTreaty(treaty.TreatyType)
		if expectedBasis != "" && r.CessionBasis != expectedBasis {
			out = append(out, finding(runID, "member", idx, 3, "minor", "cession_basis",
				"L3-105", fmt.Sprintf("Row %d: cession basis '%s' does not match treaty type '%s' (expected '%s')",
					idx, r.CessionBasis, treaty.TreatyType, expectedBasis)))
		}
	}
	return out
}

func validateL3Claims(runID string, rows []models.RIBordereauxClaimsRow, treaty models.ReinsuranceTreaty) []models.RIValidationResult {
	var out []models.RIValidationResult

	// Treaty expiry vs run period is a run-level check but most useful in L3
	for i, r := range rows {
		idx := i + 1

		// XL: claims not exceeding retention should have zero cession (IsBelowRetention=true)
		// If IsBelowRetention=false but CededClaimAmount=0, that's an inconsistency
		if !r.IsBelowRetention && r.CededClaimAmount == 0 && r.GrossClaimAmount > 0 {
			out = append(out, finding(runID, "claims", idx, 3, "major", "ceded_claim_amount",
				"L3-201", fmt.Sprintf("Row %d: claim %.2f is above retention but ceded amount is zero",
					idx, r.GrossClaimAmount)))
		}

		// Large loss: if paid and recovery expected but RecoveryReceived = 0, warn
		if r.LargeLossFlag && r.GrossPaidLosses > 0 && r.RecoveryReceived == 0 {
			out = append(out, finding(runID, "claims", idx, 3, "major", "recovery_received",
				"L3-202", fmt.Sprintf("Row %d: large loss claim %s is paid but no reinsurer recovery recorded",
					idx, r.ClaimNumber)))
		}

		// Cat event code: if LargeLossFlag=true and CatastropheEventCode is blank, warn
		if r.LargeLossFlag && strings.TrimSpace(r.CatastropheEventCode) == "" && treaty.TreatyType == "catastrophe_xl" {
			out = append(out, finding(runID, "claims", idx, 3, "minor", "catastrophe_event_code",
				"L3-203", fmt.Sprintf("Row %d: catastrophe XL treaty but no cat event code provided", idx)))
		}

		// IBNR claims should not have paid losses
		if r.IBNRFlag && r.GrossPaidLosses > 0 {
			out = append(out, finding(runID, "claims", idx, 3, "major", "ibnr_flag",
				"L3-204", fmt.Sprintf("Row %d: claim %s is flagged as IBNR but has gross paid losses %.2f",
					idx, r.ClaimNumber, r.GrossPaidLosses)))
		}
	}
	return out
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func finding(runID, rowType string, rowIndex, level int, severity, fieldName, errorCode, message string) models.RIValidationResult {
	return models.RIValidationResult{
		RunID:     runID,
		RowType:   rowType,
		RowIndex:  rowIndex,
		Level:     level,
		Severity:  severity,
		FieldName: fieldName,
		ErrorCode: errorCode,
		Message:   message,
	}
}

func cessionBasisForTreaty(treatyType string) string {
	switch treatyType {
	case "quota_share":
		return "quota_share"
	case "surplus":
		return "surplus"
	case "xl_risk", "xl_event":
		return "xl_excess"
	}
	return ""
}
