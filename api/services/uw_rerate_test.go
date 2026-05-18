package services

import (
	"math"
	"testing"
	"time"

	"api/models"
)

func TestEffectiveSA(t *testing.T) {
	tests := []struct {
		name        string
		originalCap float64
		uwCap       float64
		declined    bool
		want        float64
	}{
		{"no uw cap → original", 500_000, 0, false, 500_000},
		{"uw cap lower → uw cap", 500_000, 300_000, false, 300_000},
		{"uw cap higher → original (min)", 500_000, 800_000, false, 500_000},
		{"declined → 0 even with cap", 500_000, 800_000, true, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := effectiveSA(tc.originalCap, tc.uwCap, tc.declined)
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestComputeUWAdjustedRiskPremium_NoAdjustments(t *testing.T) {
	row := models.MemberRatingResult{
		ExpAdjLoadedGlaRate:       0.005,
		GlaCappedSumAssured:       500_000,
		ExpAdjLoadedPtdRate:       0.003,
		PtdCappedSumAssured:       400_000,
		ExpAdjLoadedCiRate:        0.004,
		CiCappedSumAssured:        300_000,
		ExpAdjLoadedSpouseGlaRate: 0.002,
		SpouseGlaCappedSumAssured: 200_000,
	}
	// No UW fields set → premium is original rates × original caps.
	got := computeUWAdjustedRiskPremium(row)
	want := 0.005*500_000 + 0.003*400_000 + 0.004*300_000 + 0.002*200_000
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestComputeUWAdjustedRiskPremium_LoadingApplied(t *testing.T) {
	row := models.MemberRatingResult{
		ExpAdjLoadedGlaRate: 0.005,
		GlaCappedSumAssured: 500_000,
		UWGlaLoading:        25, // +25%
	}
	got := computeUWAdjustedRiskPremium(row)
	want := 0.005 * 1.25 * 500_000
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestComputeUWAdjustedRiskPremium_CoverCapApplied(t *testing.T) {
	row := models.MemberRatingResult{
		ExpAdjLoadedGlaRate: 0.005,
		GlaCappedSumAssured: 1_000_000,
		UWGlaCoverCap:       400_000, // caps SA below original
	}
	got := computeUWAdjustedRiskPremium(row)
	want := 0.005 * 400_000
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestComputeUWAdjustedRiskPremium_DeclinedBenefit(t *testing.T) {
	row := models.MemberRatingResult{
		ExpAdjLoadedGlaRate: 0.005,
		GlaCappedSumAssured: 500_000,
		UWGlaDeclined:       true, // SA becomes 0
		ExpAdjLoadedPtdRate: 0.003,
		PtdCappedSumAssured: 400_000,
	}
	got := computeUWAdjustedRiskPremium(row)
	// GLA contribution = 0; PTD contribution unchanged.
	want := 0.003 * 400_000
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestApplyDecisionsToRow_NilDecisions_ResetsFields(t *testing.T) {
	row := models.MemberRatingResult{
		UWGlaLoading:  25,
		UWGlaCoverCap: 100_000,
		UWGlaDeclined: true,
	}
	applyDecisionsToRow(&row, nil)
	if row.UWGlaLoading != 0 || row.UWGlaCoverCap != 0 || row.UWGlaDeclined {
		t.Errorf("nil decisions should reset fields, got %+v", row)
	}
}

func TestApplyDecisionsToRow_AcceptOutcome(t *testing.T) {
	row := models.MemberRatingResult{}
	dec := map[benefitKey]models.UnderwritingDecision{
		benefitGla: {
			BenefitType:    "gla",
			Outcome:        models.UWOutcomeAccept,
			LoadingPercent: 50,
			CoverCap:       250_000,
			CreationDate:   time.Now(),
		},
	}
	applyDecisionsToRow(&row, dec)
	if row.UWGlaLoading != 50 {
		t.Errorf("loading not set, got %v", row.UWGlaLoading)
	}
	if row.UWGlaCoverCap != 250_000 {
		t.Errorf("cap not set, got %v", row.UWGlaCoverCap)
	}
	if row.UWGlaDeclined {
		t.Errorf("accept should not mark declined")
	}
}

func TestApplyDecisionsToRow_DeclineZerosLoadingAndCap(t *testing.T) {
	row := models.MemberRatingResult{}
	dec := map[benefitKey]models.UnderwritingDecision{
		benefitPtd: {
			BenefitType:    "ptd",
			Outcome:        models.UWOutcomeDecline,
			LoadingPercent: 50, // should be cleared
			CoverCap:       250_000, // should be cleared
		},
	}
	applyDecisionsToRow(&row, dec)
	if !row.UWPtdDeclined {
		t.Errorf("decline should mark declined")
	}
	if row.UWPtdLoading != 0 || row.UWPtdCoverCap != 0 {
		t.Errorf("decline should zero loading and cap, got %+v", row)
	}
}

func TestApplyDecisionsToRow_PostponeBehavesLikeDecline(t *testing.T) {
	// Postpone is treated as decline for premium math (no cover until
	// evidence comes back). Documented as the conservative default.
	row := models.MemberRatingResult{}
	dec := map[benefitKey]models.UnderwritingDecision{
		benefitCi: {
			BenefitType: "ci",
			Outcome:     models.UWOutcomePostpone,
		},
	}
	applyDecisionsToRow(&row, dec)
	if !row.UWCiDeclined {
		t.Errorf("postpone should set declined flag for SA-zeroing math")
	}
}

func TestSumExistingTotalPremium_PreferUWAdjustedWhenSet(t *testing.T) {
	summaries := []models.MemberRatingResultSummary{
		{TotalAnnualPremium: 100_000, UWAdjustedTotalAnnualPremium: 120_000},
		{TotalAnnualPremium: 50_000, UWAdjustedTotalAnnualPremium: 0},
	}
	got := sumExistingTotalPremium(summaries)
	// First row uses uw_adjusted (120k); second falls back to original (50k).
	want := 170_000.0
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("got %v want %v", got, want)
	}
}
