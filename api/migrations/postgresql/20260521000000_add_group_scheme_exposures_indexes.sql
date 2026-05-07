-- Adds two composite indexes on group_scheme_exposures so the GP dashboard's
-- region and industry-by-age queries hit index range scans instead of full
-- table scans. Idempotent on re-runs.

CREATE INDEX IF NOT EXISTS idx_gse_year_quote
  ON group_scheme_exposures (financial_year, quote_id);

CREATE INDEX IF NOT EXISTS idx_gse_year_status
  ON group_scheme_exposures (financial_year, quote_status);
