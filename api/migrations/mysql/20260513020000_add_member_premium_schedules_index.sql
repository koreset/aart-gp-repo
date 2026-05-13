-- Add a composite index on member_premium_schedules so that
-- refreshMemberPremiumSchedules' per-quote ORDER BY (category, member_name)
-- and the per-quote DELETE both hit an index seek + ordered scan instead
-- of a full table scan + filesort. The existing single-column quote_id
-- index doesn't cover the ORDER BY, so the 108k-row read sorts in memory.
--
-- TEXT-family string columns are widened to sized VARCHAR first, since
-- MySQL can't index TEXT without a prefix length. Idempotent on re-runs.

-- ─────────────────────────────────────────────────────────────────────────
-- 1) Widen string columns participating in the new index.
-- ─────────────────────────────────────────────────────────────────────────

-- category → VARCHAR(128)
SET @needs_alter := (SELECT COUNT(*) FROM information_schema.columns
                     WHERE table_schema = DATABASE()
                       AND table_name = 'member_premium_schedules'
                       AND column_name = 'category'
                       AND LOWER(data_type) IN ('text', 'longtext', 'mediumtext', 'tinytext'));
SET @sql := IF(@needs_alter > 0,
  'ALTER TABLE member_premium_schedules MODIFY COLUMN category VARCHAR(128) NOT NULL DEFAULT ''''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- member_name → VARCHAR(255)
SET @needs_alter := (SELECT COUNT(*) FROM information_schema.columns
                     WHERE table_schema = DATABASE()
                       AND table_name = 'member_premium_schedules'
                       AND column_name = 'member_name'
                       AND LOWER(data_type) IN ('text', 'longtext', 'mediumtext', 'tinytext'));
SET @sql := IF(@needs_alter > 0,
  'ALTER TABLE member_premium_schedules MODIFY COLUMN member_name VARCHAR(255) NOT NULL DEFAULT ''''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ─────────────────────────────────────────────────────────────────────────
-- 2) Create composite index.
-- ─────────────────────────────────────────────────────────────────────────

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE()
               AND table_name = 'member_premium_schedules'
               AND index_name = 'idx_mps_quote_category_member');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_mps_quote_category_member ON member_premium_schedules (quote_id, category, member_name)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
