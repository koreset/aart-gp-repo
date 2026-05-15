-- Migration: add region column to gla_aids_rates and reinsurance_gla_aids_rates
-- so AIDS Qx lookups can vary by scheme-category region. Default '' preserves
-- the column shape for existing rows, but row lookups now require an exact
-- region match — operators must reload the rate tables with region populated
-- per row before quotes will produce non-zero AIDS Qx under the new schema.

ALTER TABLE gla_aids_rates
    ADD COLUMN IF NOT EXISTS region VARCHAR(255) NOT NULL DEFAULT '';

ALTER TABLE reinsurance_gla_aids_rates
    ADD COLUMN IF NOT EXISTS region VARCHAR(255) NOT NULL DEFAULT '';
