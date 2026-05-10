-- Adds composite indexes on member_rating_result_summaries and
-- member_rating_results so the GP dashboard's status/year rollups and the
-- per-quote demographic aggregations hit index range scans instead of full
-- table / full-quote scans. Idempotent on re-runs.

CREATE INDEX IF NOT EXISTS idx_mrrs_status_type_creation
  ON member_rating_result_summaries (if_status, quote_type, creation_date);

CREATE INDEX IF NOT EXISTS idx_mrrs_creation_date
  ON member_rating_result_summaries (creation_date);

CREATE INDEX IF NOT EXISTS idx_mrrs_quote_creation
  ON member_rating_result_summaries (quote_id, creation_date);

CREATE INDEX IF NOT EXISTS idx_mrr_quote_category
  ON member_rating_results (quote_id, category);

CREATE INDEX IF NOT EXISTS idx_mrr_quote_age_gender
  ON member_rating_results (quote_id, age_band, gender);
