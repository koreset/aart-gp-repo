-- Generated 2026-05-01T14:33:57+02:00 for dialect mysql

-- Migration for: GroupPricingSetting (table: group_pricing_settings)

ALTER TABLE group_pricing_settings ADD COLUMN risk_alr_ceiling_pct DOUBLE NOT NULL DEFAULT 100;
ALTER TABLE group_pricing_settings ADD COLUMN risk_alr_delta_pp DOUBLE NOT NULL DEFAULT 20;
ALTER TABLE group_pricing_settings ADD COLUMN risk_thresholds_updated_at DATETIME;
ALTER TABLE group_pricing_settings ADD COLUMN risk_thresholds_updated_by VARCHAR(255);

