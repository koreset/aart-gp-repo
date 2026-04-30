-- Migration: add Free Cover Limit method selector + max-cover scaling factor.
-- Adds fcl_method to the singleton group_pricing_settings table (alongside the
-- existing discount_method) and fcl_maximum_cover_scaling_factor to
-- group_pricing_parameters. The new outlier FCL method uses the latter to
-- bound the limit at max(SA) * factor.

ALTER TABLE group_pricing_settings
    ADD COLUMN IF NOT EXISTS fcl_method VARCHAR(32) NOT NULL DEFAULT 'percentile';

ALTER TABLE group_pricing_parameters
    ADD COLUMN IF NOT EXISTS fcl_maximum_cover_scaling_factor DOUBLE PRECISION NOT NULL DEFAULT 0;
