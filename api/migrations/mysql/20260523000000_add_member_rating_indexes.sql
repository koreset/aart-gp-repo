-- Adds composite indexes on member_rating_result_summaries and
-- member_rating_results so the GP dashboard's status/year rollups and the
-- per-quote demographic aggregations hit index range scans instead of full
-- table / full-quote scans.
--
-- String columns participating in new indexes are widened from TEXT-family
-- to sized VARCHAR first, since MySQL cannot index TEXT/BLOB columns
-- without a prefix length. Idempotent on re-runs: every step checks
-- information_schema before acting.

-- ─────────────────────────────────────────────────────────────────────────
-- 1) Widen string columns on member_rating_result_summaries to VARCHAR.
-- ─────────────────────────────────────────────────────────────────────────

-- if_status → VARCHAR(64)
SET @needs_alter := (SELECT COUNT(*) FROM information_schema.columns
                     WHERE table_schema = DATABASE()
                       AND table_name = 'member_rating_result_summaries'
                       AND column_name = 'if_status'
                       AND LOWER(data_type) IN ('text', 'longtext', 'mediumtext', 'tinytext'));
SET @sql := IF(@needs_alter > 0,
  'ALTER TABLE member_rating_result_summaries MODIFY COLUMN if_status VARCHAR(64) NOT NULL DEFAULT ''''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- quote_type → VARCHAR(64)
SET @needs_alter := (SELECT COUNT(*) FROM information_schema.columns
                     WHERE table_schema = DATABASE()
                       AND table_name = 'member_rating_result_summaries'
                       AND column_name = 'quote_type'
                       AND LOWER(data_type) IN ('text', 'longtext', 'mediumtext', 'tinytext'));
SET @sql := IF(@needs_alter > 0,
  'ALTER TABLE member_rating_result_summaries MODIFY COLUMN quote_type VARCHAR(64) NOT NULL DEFAULT ''''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ─────────────────────────────────────────────────────────────────────────
-- 2) Widen string columns on member_rating_results to VARCHAR.
-- ─────────────────────────────────────────────────────────────────────────

-- category → VARCHAR(128)
SET @needs_alter := (SELECT COUNT(*) FROM information_schema.columns
                     WHERE table_schema = DATABASE()
                       AND table_name = 'member_rating_results'
                       AND column_name = 'category'
                       AND LOWER(data_type) IN ('text', 'longtext', 'mediumtext', 'tinytext'));
SET @sql := IF(@needs_alter > 0,
  'ALTER TABLE member_rating_results MODIFY COLUMN category VARCHAR(128) NOT NULL DEFAULT ''''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- age_band → VARCHAR(64)
SET @needs_alter := (SELECT COUNT(*) FROM information_schema.columns
                     WHERE table_schema = DATABASE()
                       AND table_name = 'member_rating_results'
                       AND column_name = 'age_band'
                       AND LOWER(data_type) IN ('text', 'longtext', 'mediumtext', 'tinytext'));
SET @sql := IF(@needs_alter > 0,
  'ALTER TABLE member_rating_results MODIFY COLUMN age_band VARCHAR(64) NOT NULL DEFAULT ''''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- gender → VARCHAR(32)
SET @needs_alter := (SELECT COUNT(*) FROM information_schema.columns
                     WHERE table_schema = DATABASE()
                       AND table_name = 'member_rating_results'
                       AND column_name = 'gender'
                       AND LOWER(data_type) IN ('text', 'longtext', 'mediumtext', 'tinytext'));
SET @sql := IF(@needs_alter > 0,
  'ALTER TABLE member_rating_results MODIFY COLUMN gender VARCHAR(32) NOT NULL DEFAULT ''''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ─────────────────────────────────────────────────────────────────────────
-- 3) Create indexes on member_rating_result_summaries.
-- ─────────────────────────────────────────────────────────────────────────

-- idx_mrrs_status_type_creation on (if_status, quote_type, creation_date)
SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND index_name = 'idx_mrrs_status_type_creation');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_mrrs_status_type_creation ON member_rating_result_summaries (if_status, quote_type, creation_date)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- idx_mrrs_creation_date on (creation_date)
SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND index_name = 'idx_mrrs_creation_date');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_mrrs_creation_date ON member_rating_result_summaries (creation_date)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- idx_mrrs_quote_creation on (quote_id, creation_date)
SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND index_name = 'idx_mrrs_quote_creation');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_mrrs_quote_creation ON member_rating_result_summaries (quote_id, creation_date)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ─────────────────────────────────────────────────────────────────────────
-- 4) Create indexes on member_rating_results.
-- ─────────────────────────────────────────────────────────────────────────

-- idx_mrr_quote_category on (quote_id, category)
SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND index_name = 'idx_mrr_quote_category');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_mrr_quote_category ON member_rating_results (quote_id, category)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- idx_mrr_quote_age_gender on (quote_id, age_band, gender)
SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND index_name = 'idx_mrr_quote_age_gender');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_mrr_quote_age_gender ON member_rating_results (quote_id, age_band, gender)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
