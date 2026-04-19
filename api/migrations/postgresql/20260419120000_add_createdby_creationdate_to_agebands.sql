-- Migration: add audit columns (created_by, creation_date) to
-- group_pricing_age_bands so uploaded age-band rows carry the same metadata
-- as every other rate table.

ALTER TABLE group_pricing_age_bands ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP NULL;
ALTER TABLE group_pricing_age_bands ADD COLUMN IF NOT EXISTS created_by VARCHAR(128);
