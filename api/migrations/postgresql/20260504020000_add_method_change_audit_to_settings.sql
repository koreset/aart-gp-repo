-- Migration: track when each method choice on the singleton
-- group_pricing_settings row was last changed and by whom. Two independent
-- audit pairs so the FCL method and discount method timestamps move
-- independently. Nullable timestamps mean "never explicitly changed since
-- the row was seeded".

ALTER TABLE group_pricing_settings
    ADD COLUMN IF NOT EXISTS discount_method_updated_at TIMESTAMP NULL;

ALTER TABLE group_pricing_settings
    ADD COLUMN IF NOT EXISTS discount_method_updated_by VARCHAR(255) NOT NULL DEFAULT '';

ALTER TABLE group_pricing_settings
    ADD COLUMN IF NOT EXISTS fcl_method_updated_at TIMESTAMP NULL;

ALTER TABLE group_pricing_settings
    ADD COLUMN IF NOT EXISTS fcl_method_updated_by VARCHAR(255) NOT NULL DEFAULT '';
