-- Denormalises occupation_class onto group_scheme_exposures so broker-scoped
-- exposure aggregations can group by occupation class without joining back
-- to group_pricing_quotes. Backfills existing rows from the parent quote and
-- adds a composite index on (financial_year, occupation_class). Idempotent.

ALTER TABLE group_scheme_exposures
  ADD COLUMN IF NOT EXISTS occupation_class INTEGER;

UPDATE group_scheme_exposures gse
SET occupation_class = q.occupation_class
FROM group_pricing_quotes q
WHERE q.id = gse.quote_id
  AND gse.occupation_class IS NULL;

CREATE INDEX IF NOT EXISTS idx_gse_year_occclass
  ON group_scheme_exposures (financial_year, occupation_class);
