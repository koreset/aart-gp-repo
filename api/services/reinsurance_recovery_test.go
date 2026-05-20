package services

import (
	"errors"
	"fmt"
	"testing"

	"api/models"
)

func TestApplyTreatyCession(t *testing.T) {
	cases := []struct {
		name        string
		claimAmount float64
		treaty      models.ReinsuranceTreaty
		wantReq     bool
		wantAmount  float64
	}{
		{
			name:        "xl_risk: claim inside retention is not flagged",
			claimAmount: 50_000,
			treaty: models.ReinsuranceTreaty{
				TreatyType:  "xl_risk",
				XLRetention: 100_000,
				XLLimit:     900_000,
			},
			wantReq:    false,
			wantAmount: 0,
		},
		{
			name:        "xl_risk: claim above retention cedes excess",
			claimAmount: 250_000,
			treaty: models.ReinsuranceTreaty{
				TreatyType:  "xl_risk",
				XLRetention: 100_000,
				XLLimit:     900_000,
			},
			wantReq:    true,
			wantAmount: 150_000,
		},
		{
			name:        "xl_risk: ceded amount is capped at XLLimit",
			claimAmount: 2_000_000,
			treaty: models.ReinsuranceTreaty{
				TreatyType:  "xl_risk",
				XLRetention: 100_000,
				XLLimit:     500_000,
			},
			wantReq:    true,
			wantAmount: 500_000,
		},
		{
			name:        "surplus: claim inside retention is not flagged",
			claimAmount: 80_000,
			treaty: models.ReinsuranceTreaty{
				TreatyType:      "surplus",
				RetentionAmount: 100_000,
				SurplusLines:    5,
			},
			wantReq:    false,
			wantAmount: 0,
		},
		{
			name:        "surplus: surplus above retention cedes up to lines*retention",
			claimAmount: 400_000,
			treaty: models.ReinsuranceTreaty{
				TreatyType:      "surplus",
				RetentionAmount: 100_000,
				SurplusLines:    5,
			},
			wantReq:    true,
			wantAmount: 300_000,
		},
		{
			name:        "surplus: cession capped when claim exceeds retention*(1+lines)",
			claimAmount: 2_000_000,
			treaty: models.ReinsuranceTreaty{
				TreatyType:      "surplus",
				RetentionAmount: 100_000,
				SurplusLines:    5,
			},
			wantReq:    true,
			wantAmount: 500_000,
		},
		{
			name:        "proportional: flat retention percentage cedes the complement",
			claimAmount: 200_000,
			treaty: models.ReinsuranceTreaty{
				TreatyType:          "proportional",
				RetentionPercentage: 40,
			},
			wantReq:    true,
			wantAmount: 120_000,
		},
		{
			name:        "proportional: tiered cession applies Level1 proportion in band",
			claimAmount: 300_000,
			treaty: models.ReinsuranceTreaty{
				TreatyType:            "proportional",
				TreatyCode:            "TIER-A",
				Level1Lowerbound:      0,
				Level1Upperbound:      500_000,
				Level1CededProportion: 60,
			},
			wantReq:    true,
			wantAmount: 180_000,
		},
		{
			name:        "no cession parameters configured is treated as fully retained",
			claimAmount: 100_000,
			treaty: models.ReinsuranceTreaty{
				TreatyType: "proportional",
			},
			wantReq:    false,
			wantAmount: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotReq, gotAmount := applyTreatyCession(tc.claimAmount, tc.treaty)
			if gotReq != tc.wantReq {
				t.Errorf("Required: got %v, want %v", gotReq, tc.wantReq)
			}
			if !floatNear(gotAmount, tc.wantAmount, 0.01) {
				t.Errorf("Amount: got %.4f, want %.4f", gotAmount, tc.wantAmount)
			}
		})
	}
}

func floatNear(a, b, tol float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d <= tol
}

func TestRecoveryFromSelectorError(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		wantReason string
	}{
		{"none-linked is silent skip", ErrTreatyNoneLinked, "no treaty linked"},
		{"none-applicable is info skip", ErrTreatyNoneApplicable, "no applicable treaty for relevant date"},
		{"ambiguous is warn skip", ErrTreatyAmbiguous, "ambiguous treaty selection"},
		{"wrapped sentinels still map", fmt.Errorf("wrap: %w", ErrTreatyAmbiguous), "ambiguous treaty selection"},
		{"unknown error falls through", errors.New("connection refused"), "treaty selection failed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := recoveryFromSelectorError(tc.err, 1, 2)
			if got.SkipReason != tc.wantReason {
				t.Errorf("SkipReason = %q, want %q", got.SkipReason, tc.wantReason)
			}
			if got.Required || got.Amount != 0 {
				t.Errorf("expected zero result, got Required=%v Amount=%v", got.Required, got.Amount)
			}
		})
	}
}
