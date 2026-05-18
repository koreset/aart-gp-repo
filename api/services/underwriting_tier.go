package services

import (
	"sort"

	"api/models"
)

// Underwriting tier values. Mirrored in the renderer (group_pricing UI) so any
// change here must be reflected in the frontend.
const (
	UnderwritingTierWithinFCL  = 0 // SA <= FCL — auto-accept, no evidence required
	UnderwritingTierShortForm  = 1 // SA modestly above FCL — short-form / tele-UW
	UnderwritingTierFullReview = 2 // SA significantly above FCL — human underwriter
)

// DefaultTierConfig returns the sane defaults applied when no
// UnderwritingTierConfig rows are loaded for the insurer:
//
//	ratio in [0, 1.0] → tier 0 (within FCL)
//	ratio in (1.0, 1.5] → tier 1 (short-form)
//	ratio  >  1.5      → tier 2 (full underwriting)
func DefaultTierConfig() []models.UnderwritingTierConfig {
	return []models.UnderwritingTierConfig{
		{Tier: UnderwritingTierWithinFCL, LowerMultiple: 0, UpperMultiple: 1.0, Active: true},
		{Tier: UnderwritingTierShortForm, LowerMultiple: 1.0, UpperMultiple: 1.5, Active: true},
		{Tier: UnderwritingTierFullReview, LowerMultiple: 1.5, UpperMultiple: 0, Active: true},
	}
}

// ClassifyTier returns (tier, ratio). It selects the highest tier across all
// supplied benefit sum-assured values and computes the corresponding
// ratio = SA / FCL using the benefit that produced the chosen tier.
//
// fcl <= 0 means no free-cover-limit is configured — any positive SA is
// classified as full underwriting (the safe default).
// All-zero SA across benefits returns tier 0 with ratio 0.
//
// cfg may be nil/empty; DefaultTierConfig() is used in that case.
func ClassifyTier(saByBenefit map[string]float64, fcl float64, cfg []models.UnderwritingTierConfig) (int, float64) {
	maxSA := 0.0
	for _, sa := range saByBenefit {
		if sa > maxSA {
			maxSA = sa
		}
	}
	if maxSA <= 0 {
		return UnderwritingTierWithinFCL, 0
	}
	if fcl <= 0 {
		return UnderwritingTierFullReview, 0
	}
	ratio := maxSA / fcl
	tier := tierFromRatio(ratio, cfg)
	return tier, ratio
}

// tierFromRatio walks the bands in ascending tier order and returns the first
// whose UpperMultiple covers the ratio (upper-inclusive). UpperMultiple == 0
// means "no upper bound" and acts as the catch-all. This makes boundary
// values fall into the LOWER tier — matching the existing strict-greater-than
// flag at group_pricing.go:3726 (SA == FCL is not flagged).
func tierFromRatio(ratio float64, cfg []models.UnderwritingTierConfig) int {
	if len(cfg) == 0 {
		cfg = DefaultTierConfig()
	}
	active := make([]models.UnderwritingTierConfig, 0, len(cfg))
	for _, row := range cfg {
		if row.Active {
			active = append(active, row)
		}
	}
	sort.SliceStable(active, func(i, j int) bool { return active[i].Tier < active[j].Tier })

	highest := UnderwritingTierWithinFCL
	for _, row := range active {
		if row.Tier > highest {
			highest = row.Tier
		}
		if row.UpperMultiple == 0 || ratio <= row.UpperMultiple {
			return row.Tier
		}
	}
	// No catch-all configured and ratio fell past every band — fail safe to
	// the highest configured tier (member routes to a human).
	return highest
}
