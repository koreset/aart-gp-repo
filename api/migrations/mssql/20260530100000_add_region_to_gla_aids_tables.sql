-- Migration: add region column to gla_aids_rates and reinsurance_gla_aids_rates
-- so AIDS Qx lookups can vary by scheme-category region. Default '' preserves
-- the column shape for existing rows, but row lookups now require an exact
-- region match — operators must reload the rate tables with region populated
-- per row before quotes will produce non-zero AIDS Qx under the new schema.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('gla_aids_rates') AND name = 'region')
BEGIN
    ALTER TABLE gla_aids_rates
        ADD region NVARCHAR(255) NOT NULL CONSTRAINT df_gla_aids_rates_region DEFAULT '';
END;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('reinsurance_gla_aids_rates') AND name = 'region')
BEGIN
    ALTER TABLE reinsurance_gla_aids_rates
        ADD region NVARCHAR(255) NOT NULL CONSTRAINT df_reinsurance_gla_aids_rates_region DEFAULT '';
END;
