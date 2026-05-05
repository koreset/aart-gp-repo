-- Migration: add PTD base rate methodology selection to the singleton
-- group_pricing_settings row, with audit pair for change tracking. Default
-- 'ptd_only' preserves historical behaviour (BasePtdRate uses ptd_rate alone,
-- excluding GLA AIDS rate). Switching to 'ptd_plus_gla_aids' mirrors the GLA
-- pattern: BasePtdRate += gla_aids_rate × (1 + GlaAidsRegionLoading).

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'ptd_base_rate_method')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD ptd_base_rate_method NVARCHAR(32) NOT NULL CONSTRAINT df_group_pricing_settings_ptd_base_rate_method DEFAULT 'ptd_only';
END;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'ptd_base_rate_method_updated_at')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD ptd_base_rate_method_updated_at DATETIME2 NULL;
END;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'ptd_base_rate_method_updated_by')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD ptd_base_rate_method_updated_by NVARCHAR(255) NOT NULL CONSTRAINT df_group_pricing_settings_ptd_base_rate_method_updated_by DEFAULT '';
END;
