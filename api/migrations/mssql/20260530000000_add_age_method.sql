-- Migration: add age calculation methodology selection to the singleton
-- group_pricing_settings row, with audit pair for change tracking. Default
-- 'age_next_birthday' preserves historical behaviour (member age is rounded up
-- once the commencement date crosses the birthday in the same year). Switching
-- to 'age_last_birthday' uses the floored age:
--   ROUNDDOWN((12*(YEAR(CommenDate)-YEAR(DoB)) + (MONTH(CommenDate)-MONTH(DoB)))/12, 0)

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'age_method')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD age_method NVARCHAR(32) NOT NULL CONSTRAINT df_group_pricing_settings_age_method DEFAULT 'age_next_birthday';
END;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'age_method_updated_at')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD age_method_updated_at DATETIME2 NULL;
END;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'age_method_updated_by')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD age_method_updated_by NVARCHAR(255) NOT NULL CONSTRAINT df_group_pricing_settings_age_method_updated_by DEFAULT '';
END;
