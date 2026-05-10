-- Backfills four columns on member_rating_result_summaries that exist on
-- the MemberRatingResultSummary Go struct but were never delivered via a
-- migration file. Fresh-install DBs picked them up via the bootstrap
-- AutoMigrate path; existing DBs past bootstrap did not, which broke the
-- slim /result-summary endpoint with "Unknown column" errors after it
-- started fetching them by name. Idempotent on re-runs: every step checks
-- information_schema before acting, so DBs that already have the columns
-- (because they were created via AutoMigrate or a prior hand-applied ALTER)
-- treat this migration as a no-op.

-- additional_gla_cover_age_band_type VARCHAR(64) NOT NULL DEFAULT ''
SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND column_name = 'additional_gla_cover_age_band_type');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN additional_gla_cover_age_band_type VARCHAR(64) NOT NULL DEFAULT ''''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- additional_gla_cover_male_prop_used DOUBLE NULL (Go *float64 → nullable)
SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND column_name = 'additional_gla_cover_male_prop_used');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN additional_gla_cover_male_prop_used DOUBLE NULL',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- extended_family_age_band_type VARCHAR(64) NOT NULL DEFAULT ''
SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND column_name = 'extended_family_age_band_type');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN extended_family_age_band_type VARCHAR(64) NOT NULL DEFAULT ''''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- extended_family_pricing_method VARCHAR(64) NOT NULL DEFAULT ''
SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND column_name = 'extended_family_pricing_method');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN extended_family_pricing_method VARCHAR(64) NOT NULL DEFAULT ''''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
