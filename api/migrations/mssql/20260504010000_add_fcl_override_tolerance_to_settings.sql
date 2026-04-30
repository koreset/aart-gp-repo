-- Migration: add Free Cover Limit override tolerance to the singleton
-- group_pricing_settings table. This is the fractional headroom above the
-- maximum allowed free cover limit (configured per risk rate on the
-- restrictions table) that a quote-level override is allowed to claim before
-- being clamped. Default 0.2 means "allow up to 20% above the ceiling".

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'fcl_override_tolerance')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD fcl_override_tolerance FLOAT NOT NULL CONSTRAINT df_group_pricing_settings_fcl_override_tolerance DEFAULT 0.2;
END;
