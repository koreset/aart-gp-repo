-- Migration: add binder fee and outsource fee columns to group_pricing_quotes.
-- These percentages are applied when distribution_channel = 'binder' so the
-- binderholder can iterate on a competitive final rate.

CREATE TABLE IF NOT EXISTS group_pricing_quotes (
    id SERIAL PRIMARY KEY
);

ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS loadings_binder_fee NUMERIC(20,6) DEFAULT 0;
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS loadings_outsource_fee NUMERIC(20,6) DEFAULT 0;
