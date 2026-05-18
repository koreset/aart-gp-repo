package services

import (
	"math"
	"testing"

	"api/models"
)

func TestClassifyTier_Defaults(t *testing.T) {
	tests := []struct {
		name      string
		sa        map[string]float64
		fcl       float64
		wantTier  int
		wantRatio float64
	}{
		{
			name:      "all zero SA → within FCL, ratio 0",
			sa:        map[string]float64{"gla": 0, "ptd": 0, "ci": 0, "sgla": 0},
			fcl:       500_000,
			wantTier:  UnderwritingTierWithinFCL,
			wantRatio: 0,
		},
		{
			name:      "SA below FCL → within FCL",
			sa:        map[string]float64{"gla": 400_000},
			fcl:       500_000,
			wantTier:  UnderwritingTierWithinFCL,
			wantRatio: 0.8,
		},
		{
			name:      "SA exactly at FCL (boundary) → within FCL",
			sa:        map[string]float64{"gla": 500_000},
			fcl:       500_000,
			wantTier:  UnderwritingTierWithinFCL,
			wantRatio: 1.0,
		},
		{
			name:      "SA modestly above FCL → short-form",
			sa:        map[string]float64{"gla": 600_000},
			fcl:       500_000,
			wantTier:  UnderwritingTierShortForm,
			wantRatio: 1.2,
		},
		{
			name:      "SA at short-form upper boundary (1.5x) → short-form",
			sa:        map[string]float64{"gla": 750_000},
			fcl:       500_000,
			wantTier:  UnderwritingTierShortForm,
			wantRatio: 1.5,
		},
		{
			name:      "SA significantly above FCL → full UW",
			sa:        map[string]float64{"gla": 900_000},
			fcl:       500_000,
			wantTier:  UnderwritingTierFullReview,
			wantRatio: 1.8,
		},
		{
			name:      "max across benefits drives tier (PTD highest)",
			sa:        map[string]float64{"gla": 400_000, "ptd": 800_000, "ci": 200_000},
			fcl:       500_000,
			wantTier:  UnderwritingTierFullReview,
			wantRatio: 1.6,
		},
		{
			name:      "FCL == 0 with positive SA → fail-safe full UW",
			sa:        map[string]float64{"gla": 100_000},
			fcl:       0,
			wantTier:  UnderwritingTierFullReview,
			wantRatio: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotTier, gotRatio := ClassifyTier(tc.sa, tc.fcl, nil)
			if gotTier != tc.wantTier {
				t.Errorf("tier = %d, want %d", gotTier, tc.wantTier)
			}
			if math.Abs(gotRatio-tc.wantRatio) > 1e-9 {
				t.Errorf("ratio = %v, want %v", gotRatio, tc.wantRatio)
			}
		})
	}
}

func TestClassifyTier_CustomConfig(t *testing.T) {
	cfg := []models.UnderwritingTierConfig{
		{Tier: UnderwritingTierWithinFCL, LowerMultiple: 0, UpperMultiple: 1.0, Active: true},
		{Tier: UnderwritingTierShortForm, LowerMultiple: 1.0, UpperMultiple: 2.0, Active: true},
		{Tier: UnderwritingTierFullReview, LowerMultiple: 2.0, UpperMultiple: 0, Active: true},
	}
	// Ratio of 1.8 would be full-UW under defaults (>1.5) but short-form here.
	tier, ratio := ClassifyTier(map[string]float64{"gla": 900_000}, 500_000, cfg)
	if tier != UnderwritingTierShortForm {
		t.Errorf("with widened short-form band, ratio 1.8 should be short-form, got tier %d", tier)
	}
	if math.Abs(ratio-1.8) > 1e-9 {
		t.Errorf("ratio = %v, want 1.8", ratio)
	}
}

func TestClassifyTier_InactiveRowsIgnored(t *testing.T) {
	cfg := []models.UnderwritingTierConfig{
		{Tier: UnderwritingTierWithinFCL, LowerMultiple: 0, UpperMultiple: 1.0, Active: true},
		{Tier: UnderwritingTierShortForm, LowerMultiple: 1.0, UpperMultiple: 1.5, Active: false},
		{Tier: UnderwritingTierFullReview, LowerMultiple: 1.0, UpperMultiple: 0, Active: true},
	}
	// Short-form band is inactive, so ratio 1.2 should fall to full-review.
	tier, _ := ClassifyTier(map[string]float64{"gla": 600_000}, 500_000, cfg)
	if tier != UnderwritingTierFullReview {
		t.Errorf("inactive short-form band, ratio 1.2 should fall to full-UW, got tier %d", tier)
	}
}
