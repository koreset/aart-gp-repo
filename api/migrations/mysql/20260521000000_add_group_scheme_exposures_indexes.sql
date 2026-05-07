-- Adds two composite indexes on group_scheme_exposures so the GP dashboard's
-- region and industry-by-age queries hit index range scans instead of full
-- table scans. quote_status is widened to a sized VARCHAR first because
-- MySQL (unlike Postgres) cannot index TEXT/BLOB columns without a prefix
-- length. Idempotent on re-runs.

-- 1) Convert quote_status from TEXT/LONGTEXT to VARCHAR(64) if needed.
SET @needs_alter := (SELECT COUNT(*) FROM information_schema.columns
                     WHERE table_schema = DATABASE()
                       AND table_name = 'group_scheme_exposures'
                       AND column_name = 'quote_status'
                       AND LOWER(data_type) IN ('text', 'longtext', 'mediumtext', 'tinytext'));
SET @sql := IF(@needs_alter > 0,
  'ALTER TABLE group_scheme_exposures MODIFY COLUMN quote_status VARCHAR(64) NOT NULL DEFAULT ''''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 2) Create idx_gse_year_quote on (financial_year, quote_id) if missing.
SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE()
               AND table_name = 'group_scheme_exposures'
               AND index_name = 'idx_gse_year_quote');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_gse_year_quote ON group_scheme_exposures (financial_year, quote_id)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 3) Create idx_gse_year_status on (financial_year, quote_status) if missing.
SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE()
               AND table_name = 'group_scheme_exposures'
               AND index_name = 'idx_gse_year_status');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_gse_year_status ON group_scheme_exposures (financial_year, quote_status)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
