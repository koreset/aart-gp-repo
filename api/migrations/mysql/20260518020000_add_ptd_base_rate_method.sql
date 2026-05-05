-- Migration: add PTD base rate methodology selection to the singleton
-- group_pricing_settings row, with audit pair for change tracking. Default
-- 'ptd_only' preserves historical behaviour (BasePtdRate uses ptd_rate alone,
-- excluding GLA AIDS rate). Switching to 'ptd_plus_gla_aids' mirrors the GLA
-- pattern: BasePtdRate += gla_aids_rate × (1 + GlaAidsRegionLoading).

ALTER TABLE group_pricing_settings
    ADD COLUMN ptd_base_rate_method VARCHAR(32) NOT NULL DEFAULT 'ptd_only',
    ADD COLUMN ptd_base_rate_method_updated_at DATETIME NULL,
    ADD COLUMN ptd_base_rate_method_updated_by VARCHAR(255) NOT NULL DEFAULT '';
