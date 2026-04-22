-- Migration: rename educator GLA/PTD rate-per-1000 columns on
-- member_rating_result_summaries so they match GORM's NamingStrategy output
-- (`per1000_sa`, no underscore between `per` and `1000`). The
-- 20260421170000 migration used `per_1000_sa` which GORM cannot find on
-- INSERT/UPDATE — causing "Unknown column" 1054 errors during
-- CalculateGroupPricingQuote. Follows the precedent of
-- 20260418150000_rename_additional_accidental_gla_per1000.sql.

-- If the (incorrect) underscored column exists AND the correct one does not,
-- rename it. Otherwise no-op. Safe to re-run.

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='gla_educator_risk_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='gla_educator_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN gla_educator_risk_rate_per_1000_sa gla_educator_risk_rate_per1000_sa DOUBLE;',
    'SELECT 1;'
));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='gla_educator_office_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='gla_educator_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN gla_educator_office_rate_per_1000_sa gla_educator_office_rate_per1000_sa DOUBLE;',
    'SELECT 1;'
));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_gla_educator_risk_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_gla_educator_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN exp_gla_educator_risk_rate_per_1000_sa exp_gla_educator_risk_rate_per1000_sa DOUBLE;',
    'SELECT 1;'
));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_gla_educator_office_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_gla_educator_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN exp_gla_educator_office_rate_per_1000_sa exp_gla_educator_office_rate_per1000_sa DOUBLE;',
    'SELECT 1;'
));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ptd_educator_risk_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ptd_educator_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN ptd_educator_risk_rate_per_1000_sa ptd_educator_risk_rate_per1000_sa DOUBLE;',
    'SELECT 1;'
));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ptd_educator_office_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ptd_educator_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN ptd_educator_office_rate_per_1000_sa ptd_educator_office_rate_per1000_sa DOUBLE;',
    'SELECT 1;'
));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ptd_educator_risk_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ptd_educator_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN exp_ptd_educator_risk_rate_per_1000_sa exp_ptd_educator_risk_rate_per1000_sa DOUBLE;',
    'SELECT 1;'
));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ptd_educator_office_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ptd_educator_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN exp_ptd_educator_office_rate_per_1000_sa exp_ptd_educator_office_rate_per1000_sa DOUBLE;',
    'SELECT 1;'
));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
