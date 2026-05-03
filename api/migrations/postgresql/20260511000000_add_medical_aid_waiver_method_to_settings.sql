-- Migration: add medical aid waiver methodology selection to the singleton
-- group_pricing_settings row, with audit pair for change tracking. Default
-- 'formula' preserves the existing salary-based calculation behavior.

ALTER TABLE group_pricing_settings
    ADD COLUMN IF NOT EXISTS medical_aid_waiver_method VARCHAR(32) NOT NULL DEFAULT 'formula';

ALTER TABLE group_pricing_settings
    ADD COLUMN IF NOT EXISTS medical_aid_waiver_method_updated_at TIMESTAMP NULL;

ALTER TABLE group_pricing_settings
    ADD COLUMN IF NOT EXISTS medical_aid_waiver_method_updated_by VARCHAR(255) NOT NULL DEFAULT '';
