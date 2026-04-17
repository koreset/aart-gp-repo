package services

import (
	"api/models"
	"fmt"
	"math"
	"strconv"
)

// FieldDiff is a single before/after field change used across the diff payload.
type FieldDiff struct {
	Field  string      `json:"field"`
	Before interface{} `json:"before"`
	After  interface{} `json:"after"`
}

// MemberRowDiff represents a member row that changed between two runs.
type MemberRowDiff struct {
	Key     string                     `json:"key"`
	Before  models.RIBordereauxMemberRow `json:"before"`
	After   models.RIBordereauxMemberRow `json:"after"`
	Changes []FieldDiff                  `json:"changes"`
}

type ClaimRowDiff struct {
	Key     string                       `json:"key"`
	Before  models.RIBordereauxClaimsRow `json:"before"`
	After   models.RIBordereauxClaimsRow `json:"after"`
	Changes []FieldDiff                  `json:"changes"`
}

// RIBordereauxRunDiff is the payload returned by DiffRIBordereauxRuns.
// Members/Claims are grouped by semantic change so UIs can render them under
// add/remove/change tabs without re-sorting client-side.
type RIBordereauxRunDiff struct {
	Base        models.RIBordereauxRun `json:"base"`
	Amendment   models.RIBordereauxRun `json:"amendment"`
	SummaryDiff []FieldDiff            `json:"summary_diff"`
	Members     struct {
		Added   []models.RIBordereauxMemberRow `json:"added"`
		Removed []models.RIBordereauxMemberRow `json:"removed"`
		Changed []MemberRowDiff                `json:"changed"`
	} `json:"members"`
	Claims struct {
		Added   []models.RIBordereauxClaimsRow `json:"added"`
		Removed []models.RIBordereauxClaimsRow `json:"removed"`
		Changed []ClaimRowDiff                 `json:"changed"`
	} `json:"claims"`
}

// resolveParentRunID looks up the string run_id for a parent referenced by
// ParentRunID (the numeric DB id).
func resolveParentRunID(parentID uint) (string, error) {
	var parent models.RIBordereauxRun
	if err := DB.Where("id = ?", parentID).First(&parent).Error; err != nil {
		return "", fmt.Errorf("parent run (id=%d) not found: %w", parentID, err)
	}
	return parent.RunID, nil
}

// DiffRIBordereauxRuns compares two RI bordereaux runs and returns a structured
// diff covering header aggregates plus per-row additions, removals, and field
// changes. If againstRunID is empty, the base run's ParentRunID is used so the
// natural "what changed in this amendment?" query doesn't need a second lookup
// from the caller.
//
// Conceptually the diff reads: `base` → `amendment`. Pass the amendment run as
// fromRunID and (optionally) the parent as againstRunID to render "what's new
// in this version?" in the UI.
func DiffRIBordereauxRuns(fromRunID, againstRunID string) (RIBordereauxRunDiff, error) {
	var out RIBordereauxRunDiff

	amendment, err := GetRIBordereauxRunByID(fromRunID)
	if err != nil {
		return out, err
	}
	out.Amendment = amendment

	if againstRunID == "" {
		if amendment.ParentRunID == nil {
			return out, fmt.Errorf("run %s has no parent to diff against; pass ?against=<run_id>", fromRunID)
		}
		parentRunID, err := resolveParentRunID(*amendment.ParentRunID)
		if err != nil {
			return out, err
		}
		againstRunID = parentRunID
	}

	base, err := GetRIBordereauxRunByID(againstRunID)
	if err != nil {
		return out, err
	}
	out.Base = base

	out.SummaryDiff = diffRunHeader(base, amendment)

	// Member rows
	baseMembers, err := GetRIBordereauxMemberRows(base.RunID)
	if err != nil {
		return out, fmt.Errorf("load base member rows: %w", err)
	}
	amendMembers, err := GetRIBordereauxMemberRows(amendment.RunID)
	if err != nil {
		return out, fmt.Errorf("load amendment member rows: %w", err)
	}
	out.Members.Added, out.Members.Removed, out.Members.Changed = diffMemberRows(baseMembers, amendMembers)

	// Claim rows
	baseClaims, err := GetRIBordereauxClaimsRows(base.RunID)
	if err != nil {
		return out, fmt.Errorf("load base claim rows: %w", err)
	}
	amendClaims, err := GetRIBordereauxClaimsRows(amendment.RunID)
	if err != nil {
		return out, fmt.Errorf("load amendment claim rows: %w", err)
	}
	out.Claims.Added, out.Claims.Removed, out.Claims.Changed = diffClaimRows(baseClaims, amendClaims)

	return out, nil
}

// diffRunHeader walks the business-relevant summary fields and emits a FieldDiff
// for each change. Infrastructure fields (ID, timestamps, file paths) are not
// compared because they always differ between versions.
func diffRunHeader(base, amend models.RIBordereauxRun) []FieldDiff {
	diffs := []FieldDiff{}
	if base.Status != amend.Status {
		diffs = append(diffs, FieldDiff{"status", base.Status, amend.Status})
	}
	if base.TotalLives != amend.TotalLives {
		diffs = append(diffs, FieldDiff{"total_lives", base.TotalLives, amend.TotalLives})
	}
	if base.TotalCededLives != amend.TotalCededLives {
		diffs = append(diffs, FieldDiff{"total_ceded_lives", base.TotalCededLives, amend.TotalCededLives})
	}
	if !floatsEqual(base.GrossPremium, amend.GrossPremium) {
		diffs = append(diffs, FieldDiff{"gross_premium", base.GrossPremium, amend.GrossPremium})
	}
	if !floatsEqual(base.CededPremium, amend.CededPremium) {
		diffs = append(diffs, FieldDiff{"ceded_premium", base.CededPremium, amend.CededPremium})
	}
	if !floatsEqual(base.RetainedPremium, amend.RetainedPremium) {
		diffs = append(diffs, FieldDiff{"retained_premium", base.RetainedPremium, amend.RetainedPremium})
	}
	if !floatsEqual(base.GrossClaimsIncurred, amend.GrossClaimsIncurred) {
		diffs = append(diffs, FieldDiff{"gross_claims_incurred", base.GrossClaimsIncurred, amend.GrossClaimsIncurred})
	}
	if !floatsEqual(base.CededClaimsIncurred, amend.CededClaimsIncurred) {
		diffs = append(diffs, FieldDiff{"ceded_claims_incurred", base.CededClaimsIncurred, amend.CededClaimsIncurred})
	}
	if base.AmendmentNotes != amend.AmendmentNotes {
		diffs = append(diffs, FieldDiff{"amendment_notes", base.AmendmentNotes, amend.AmendmentNotes})
	}
	return diffs
}

func riMemberKey(r models.RIBordereauxMemberRow) string {
	return strconv.Itoa(r.SchemeID) + "|" + r.MemberIDNumber + "|" + r.BenefitCode
}

func diffMemberRows(base, amend []models.RIBordereauxMemberRow) (added, removed []models.RIBordereauxMemberRow, changed []MemberRowDiff) {
	baseMap := make(map[string]models.RIBordereauxMemberRow, len(base))
	for _, r := range base {
		baseMap[riMemberKey(r)] = r
	}
	amendKeys := make(map[string]bool, len(amend))
	for _, a := range amend {
		k := riMemberKey(a)
		amendKeys[k] = true
		b, ok := baseMap[k]
		if !ok {
			added = append(added, a)
			continue
		}
		if diffs := diffMemberRow(b, a); len(diffs) > 0 {
			changed = append(changed, MemberRowDiff{Key: k, Before: b, After: a, Changes: diffs})
		}
	}
	for k, b := range baseMap {
		if !amendKeys[k] {
			removed = append(removed, b)
		}
	}
	return
}

func diffMemberRow(b, a models.RIBordereauxMemberRow) []FieldDiff {
	var diffs []FieldDiff
	if b.MemberName != a.MemberName {
		diffs = append(diffs, FieldDiff{"member_name", b.MemberName, a.MemberName})
	}
	if !floatsEqual(b.SumAssured, a.SumAssured) {
		diffs = append(diffs, FieldDiff{"sum_assured", b.SumAssured, a.SumAssured})
	}
	if !floatsEqual(b.AnnualSalary, a.AnnualSalary) {
		diffs = append(diffs, FieldDiff{"annual_salary", b.AnnualSalary, a.AnnualSalary})
	}
	if !floatsEqual(b.GrossPremium, a.GrossPremium) {
		diffs = append(diffs, FieldDiff{"gross_premium", b.GrossPremium, a.GrossPremium})
	}
	if !floatsEqual(b.CededPremium, a.CededPremium) {
		diffs = append(diffs, FieldDiff{"ceded_premium", b.CededPremium, a.CededPremium})
	}
	if !floatsEqual(b.CededAmount, a.CededAmount) {
		diffs = append(diffs, FieldDiff{"ceded_amount", b.CededAmount, a.CededAmount})
	}
	if b.MemberStatus != a.MemberStatus {
		diffs = append(diffs, FieldDiff{"member_status", b.MemberStatus, a.MemberStatus})
	}
	if b.ChangeType != a.ChangeType {
		diffs = append(diffs, FieldDiff{"change_type", b.ChangeType, a.ChangeType})
	}
	return diffs
}

func riClaimKey(r models.RIBordereauxClaimsRow) string {
	if r.ClaimNumber != "" {
		return r.ClaimNumber
	}
	return strconv.Itoa(r.SchemeID) + "|" + r.MemberIDNumber + "|" + r.DateOfEvent
}

func diffClaimRows(base, amend []models.RIBordereauxClaimsRow) (added, removed []models.RIBordereauxClaimsRow, changed []ClaimRowDiff) {
	baseMap := make(map[string]models.RIBordereauxClaimsRow, len(base))
	for _, r := range base {
		baseMap[riClaimKey(r)] = r
	}
	amendKeys := make(map[string]bool, len(amend))
	for _, a := range amend {
		k := riClaimKey(a)
		amendKeys[k] = true
		b, ok := baseMap[k]
		if !ok {
			added = append(added, a)
			continue
		}
		if diffs := diffClaimRow(b, a); len(diffs) > 0 {
			changed = append(changed, ClaimRowDiff{Key: k, Before: b, After: a, Changes: diffs})
		}
	}
	for k, b := range baseMap {
		if !amendKeys[k] {
			removed = append(removed, b)
		}
	}
	return
}

func diffClaimRow(b, a models.RIBordereauxClaimsRow) []FieldDiff {
	var diffs []FieldDiff
	if !floatsEqual(b.GrossClaimAmount, a.GrossClaimAmount) {
		diffs = append(diffs, FieldDiff{"gross_claim_amount", b.GrossClaimAmount, a.GrossClaimAmount})
	}
	if !floatsEqual(b.CededClaimAmount, a.CededClaimAmount) {
		diffs = append(diffs, FieldDiff{"ceded_claim_amount", b.CededClaimAmount, a.CededClaimAmount})
	}
	if !floatsEqual(b.RecoveryReceived, a.RecoveryReceived) {
		diffs = append(diffs, FieldDiff{"recovery_received", b.RecoveryReceived, a.RecoveryReceived})
	}
	if b.ClaimStatus != a.ClaimStatus {
		diffs = append(diffs, FieldDiff{"claim_status", b.ClaimStatus, a.ClaimStatus})
	}
	if b.IsBelowRetention != a.IsBelowRetention {
		diffs = append(diffs, FieldDiff{"is_below_retention", b.IsBelowRetention, a.IsBelowRetention})
	}
	if !floatsEqual(b.GrossPaidLosses, a.GrossPaidLosses) {
		diffs = append(diffs, FieldDiff{"gross_paid_losses", b.GrossPaidLosses, a.GrossPaidLosses})
	}
	if !floatsEqual(b.GrossOutstandingReserve, a.GrossOutstandingReserve) {
		diffs = append(diffs, FieldDiff{"gross_outstanding_reserve", b.GrossOutstandingReserve, a.GrossOutstandingReserve})
	}
	return diffs
}

// floatsEqual ignores sub-cent noise. Reuses a conservative threshold that
// matches the reconciliation default tolerance.
func floatsEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.001
}
