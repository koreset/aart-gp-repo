-- Migration: add Free Cover Limit method selector + max-cover scaling factor.
-- Adds fcl_method to the singleton group_pricing_settings table (alongside the
-- existing discount_method) and fcl_maximum_cover_scaling_factor to
-- group_pricing_parameters. The new outlier FCL method uses the latter to
-- bound the limit at max(SA) * factor.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'fcl_method')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD fcl_method NVARCHAR(32) NOT NULL CONSTRAINT df_group_pricing_settings_fcl_method DEFAULT 'percentile';
END;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_parameters') AND name = 'fcl_maximum_cover_scaling_factor')
BEGIN
    ALTER TABLE group_pricing_parameters
        ADD fcl_maximum_cover_scaling_factor FLOAT NOT NULL CONSTRAINT df_group_pricing_parameters_fcl_max_cover DEFAULT 0;
END;
