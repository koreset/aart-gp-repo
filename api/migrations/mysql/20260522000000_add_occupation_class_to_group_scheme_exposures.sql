-- Denormalises occupation_class onto group_scheme_exposures so broker-scoped
-- exposure aggregations can group by occupation class without joining back
-- to group_pricing_quotes. Backfills existing rows from the parent quote and
-- adds a composite index on (financial_year, occupation_class). Idempotent.

-- 1) Add occupation_class column if missing.
SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'group_scheme_exposures'
               AND column_name = 'occupation_class');
SET @sql := IF(@col = 0,
  'ALTER TABLE group_scheme_exposures ADD COLUMN occupation_class INT NULL',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 2) Backfill from parent quote.
UPDATE group_scheme_exposures gse
JOIN group_pricing_quotes q ON q.id = gse.quote_id
SET gse.occupation_class = q.occupation_class
WHERE gse.occupation_class IS NULL;

-- 3) Create idx_gse_year_occclass on (financial_year, occupation_class).
SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE()
               AND table_name = 'group_scheme_exposures'
               AND index_name = 'idx_gse_year_occclass');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_gse_year_occclass ON group_scheme_exposures (financial_year, occupation_class)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
