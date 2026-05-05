-- Migration: add risk profile variation tolerance (in percent) to the
-- singleton group_pricing_settings row. Surfaces in the Acceptance Form text
-- on the generated quote PDF / DOCX. Default 7 preserves the historical
-- hardcoded "7%" wording.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'risk_profile_variation_tolerance_pct')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD risk_profile_variation_tolerance_pct FLOAT NOT NULL CONSTRAINT df_group_pricing_settings_risk_profile_variation_tolerance_pct DEFAULT 7;
END;
