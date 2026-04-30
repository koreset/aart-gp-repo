-- Migration: add Free Cover Limit override tolerance to the singleton
-- group_pricing_settings table. This is the fractional headroom above the
-- maximum allowed free cover limit (configured per risk rate on the
-- restrictions table) that a quote-level override is allowed to claim before
-- being clamped. Default 0.2 means "allow up to 20% above the ceiling".

ALTER TABLE group_pricing_settings
    ADD COLUMN IF NOT EXISTS fcl_override_tolerance DOUBLE PRECISION NOT NULL DEFAULT 0.2;
