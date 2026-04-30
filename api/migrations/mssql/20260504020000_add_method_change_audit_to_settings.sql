-- Migration: track when each method choice on the singleton
-- group_pricing_settings row was last changed and by whom. Two independent
-- audit pairs so the FCL method and discount method timestamps move
-- independently. Nullable timestamps mean "never explicitly changed since
-- the row was seeded".

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'discount_method_updated_at')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD discount_method_updated_at DATETIME2 NULL;
END;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'discount_method_updated_by')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD discount_method_updated_by NVARCHAR(255) NOT NULL CONSTRAINT df_group_pricing_settings_discount_method_updated_by DEFAULT '';
END;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'fcl_method_updated_at')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD fcl_method_updated_at DATETIME2 NULL;
END;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_settings') AND name = 'fcl_method_updated_by')
BEGIN
    ALTER TABLE group_pricing_settings
        ADD fcl_method_updated_by NVARCHAR(255) NOT NULL CONSTRAINT df_group_pricing_settings_fcl_method_updated_by DEFAULT '';
END;
