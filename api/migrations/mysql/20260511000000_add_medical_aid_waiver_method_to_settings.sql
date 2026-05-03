-- Migration: add medical aid waiver methodology selection to the singleton
-- group_pricing_settings row, with audit pair for change tracking. Default
-- 'formula' preserves the existing salary-based calculation behavior.

ALTER TABLE group_pricing_settings
    ADD COLUMN medical_aid_waiver_method VARCHAR(32) NOT NULL DEFAULT 'formula',
    ADD COLUMN medical_aid_waiver_method_updated_at DATETIME NULL,
    ADD COLUMN medical_aid_waiver_method_updated_by VARCHAR(255) NOT NULL DEFAULT '';
