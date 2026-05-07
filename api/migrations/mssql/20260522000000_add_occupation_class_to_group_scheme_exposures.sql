-- Denormalises occupation_class onto group_scheme_exposures so broker-scoped
-- exposure aggregations can group by occupation class without joining back
-- to group_pricing_quotes. Backfills existing rows from the parent quote and
-- adds a composite index on (financial_year, occupation_class). Idempotent.
--
-- The runner splits on ';' and executes each statement in its own batch, so
-- referencing the newly-added column from a subsequent UPDATE is safe.

-- 1) Add occupation_class column if missing.
IF NOT EXISTS (SELECT 1 FROM sys.columns
               WHERE object_id = OBJECT_ID('group_scheme_exposures')
                 AND name = 'occupation_class')
  ALTER TABLE group_scheme_exposures ADD occupation_class INT NULL;

-- 2) Backfill from parent quote.
UPDATE gse
SET gse.occupation_class = q.occupation_class
FROM group_scheme_exposures gse
JOIN group_pricing_quotes q ON q.id = gse.quote_id
WHERE gse.occupation_class IS NULL;

-- 3) Create idx_gse_year_occclass on (financial_year, occupation_class).
IF NOT EXISTS (SELECT 1 FROM sys.indexes
               WHERE name = 'idx_gse_year_occclass'
                 AND object_id = OBJECT_ID('group_scheme_exposures'))
  CREATE INDEX idx_gse_year_occclass
    ON group_scheme_exposures (financial_year, occupation_class);
