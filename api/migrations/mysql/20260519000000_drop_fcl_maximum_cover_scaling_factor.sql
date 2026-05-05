-- Migration: drop the per-parameter fcl_maximum_cover_scaling_factor column.
-- The Statistical Outlier FCL method now derives its max-cover cap from the
-- system-wide GroupPricingSetting.FCLOverrideTolerance, so the per-parameter
-- field is redundant. See api/services/group_pricing.go.

ALTER TABLE group_pricing_parameters
    DROP COLUMN fcl_maximum_cover_scaling_factor;
