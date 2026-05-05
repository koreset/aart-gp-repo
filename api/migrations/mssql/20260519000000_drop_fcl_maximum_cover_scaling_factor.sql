-- Migration: drop the per-parameter fcl_maximum_cover_scaling_factor column.
-- The Statistical Outlier FCL method now derives its max-cover cap from the
-- system-wide GroupPricingSetting.FCLOverrideTolerance, so the per-parameter
-- field is redundant. See api/services/group_pricing.go.

IF EXISTS (
    SELECT 1 FROM sys.default_constraints
    WHERE name = 'df_group_pricing_parameters_fcl_max_cover'
)
BEGIN
    ALTER TABLE group_pricing_parameters
        DROP CONSTRAINT df_group_pricing_parameters_fcl_max_cover;
END;

IF EXISTS (
    SELECT 1 FROM sys.columns
    WHERE object_id = OBJECT_ID('group_pricing_parameters')
      AND name = 'fcl_maximum_cover_scaling_factor'
)
BEGIN
    ALTER TABLE group_pricing_parameters
        DROP COLUMN fcl_maximum_cover_scaling_factor;
END;
