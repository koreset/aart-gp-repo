-- Migration: add medical aid waiver methodology selection to the singleton
-- group_pricing_settings row, with audit pair for change tracking. Default
-- 'formula' preserves the existing salary-based calculation behavior.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'medical_aid_waiver_method')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD medical_aid_waiver_method NVARCHAR(32) NOT NULL CONSTRAINT df_group_pricing_settings_medical_aid_waiver_method DEFAULT 'formula';
END;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'medical_aid_waiver_method_updated_at')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD medical_aid_waiver_method_updated_at DATETIME2 NULL;
END;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'medical_aid_waiver_method_updated_by')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD medical_aid_waiver_method_updated_by NVARCHAR(255) NOT NULL CONSTRAINT df_group_pricing_settings_medical_aid_waiver_method_updated_by DEFAULT '';
END;
