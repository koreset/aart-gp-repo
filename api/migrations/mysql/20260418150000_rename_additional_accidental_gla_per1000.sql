-- Migration: rename four Additional Accidental GLA rate-per-1000 columns on
-- member_rating_result_summaries so they match GORM's NamingStrategy output
-- (`per1000_sa`, no underscore between `per` and `1000`). The earlier
-- 20260418140000 migration used `per_1000_sa` which GORM cannot find on INSERT.

-- If the (incorrect) underscored column exists AND the correct one does not,
-- rename it. Otherwise no-op. Safe to re-run.

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='additional_accidental_gla_risk_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='additional_accidental_gla_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN additional_accidental_gla_risk_rate_per_1000_sa additional_accidental_gla_risk_rate_per1000_sa DECIMAL(15,5);',
    'SELECT 1;'
));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='additional_accidental_gla_office_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='additional_accidental_gla_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN additional_accidental_gla_office_rate_per_1000_sa additional_accidental_gla_office_rate_per1000_sa DECIMAL(15,5);',
    'SELECT 1;'
));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_additional_accidental_gla_risk_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_additional_accidental_gla_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN exp_additional_accidental_gla_risk_rate_per_1000_sa exp_additional_accidental_gla_risk_rate_per1000_sa DECIMAL(15,5);',
    'SELECT 1;'
));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_additional_accidental_gla_office_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_additional_accidental_gla_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN exp_additional_accidental_gla_office_rate_per_1000_sa exp_additional_accidental_gla_office_rate_per1000_sa DECIMAL(15,5);',
    'SELECT 1;'
));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
