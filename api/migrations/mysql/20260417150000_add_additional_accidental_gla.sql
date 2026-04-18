-- Migration: add optional "Additional Accidental GLA" sub-benefit.
-- Re-uses the gla_rates / gla_aids_rates tables with a different benefit_type.
-- All GLA scheme parameters stay on the GLA fields — only the benefit_type
-- differs for the Additional Accidental layer.

--------------------------------------------------------------------------------
-- scheme_categories: toggle + benefit type for the Additional Accidental layer
--------------------------------------------------------------------------------

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='additional_accidental_gla_benefit'),
        'ALTER TABLE scheme_categories MODIFY COLUMN additional_accidental_gla_benefit TINYINT(1) DEFAULT 0;',
        'ALTER TABLE scheme_categories ADD COLUMN additional_accidental_gla_benefit TINYINT(1) DEFAULT 0;'
    )
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='additional_accidental_gla_benefit_type'),
        'ALTER TABLE scheme_categories MODIFY COLUMN additional_accidental_gla_benefit_type VARCHAR(255);',
        'ALTER TABLE scheme_categories ADD COLUMN additional_accidental_gla_benefit_type VARCHAR(255);'
    )
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

--------------------------------------------------------------------------------
-- member_rating_results: per-member Additional Accidental GLA outputs
--------------------------------------------------------------------------------

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='additional_accidental_gla_sum_assured'),
    'ALTER TABLE member_rating_results MODIFY COLUMN additional_accidental_gla_sum_assured DECIMAL(15,5);',
    'ALTER TABLE member_rating_results ADD COLUMN additional_accidental_gla_sum_assured DECIMAL(15,5);'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='additional_accidental_gla_capped_sum_assured'),
    'ALTER TABLE member_rating_results MODIFY COLUMN additional_accidental_gla_capped_sum_assured DECIMAL(15,5);',
    'ALTER TABLE member_rating_results ADD COLUMN additional_accidental_gla_capped_sum_assured DECIMAL(15,5);'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='additional_accidental_gla_qx'),
    'ALTER TABLE member_rating_results MODIFY COLUMN additional_accidental_gla_qx DECIMAL(15,5);',
    'ALTER TABLE member_rating_results ADD COLUMN additional_accidental_gla_qx DECIMAL(15,5);'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='additional_accidental_gla_aids_qx'),
    'ALTER TABLE member_rating_results MODIFY COLUMN additional_accidental_gla_aids_qx DECIMAL(15,5);',
    'ALTER TABLE member_rating_results ADD COLUMN additional_accidental_gla_aids_qx DECIMAL(15,5);'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='base_additional_accidental_gla_rate'),
    'ALTER TABLE member_rating_results MODIFY COLUMN base_additional_accidental_gla_rate DECIMAL(15,5);',
    'ALTER TABLE member_rating_results ADD COLUMN base_additional_accidental_gla_rate DECIMAL(15,5);'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='loaded_additional_accidental_gla_rate'),
    'ALTER TABLE member_rating_results MODIFY COLUMN loaded_additional_accidental_gla_rate DECIMAL(15,5);',
    'ALTER TABLE member_rating_results ADD COLUMN loaded_additional_accidental_gla_rate DECIMAL(15,5);'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='additional_accidental_gla_weighted_experience_crude_rate'),
    'ALTER TABLE member_rating_results MODIFY COLUMN additional_accidental_gla_weighted_experience_crude_rate DECIMAL(15,5);',
    'ALTER TABLE member_rating_results ADD COLUMN additional_accidental_gla_weighted_experience_crude_rate DECIMAL(15,5);'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='additional_accidental_gla_theoretical_rate'),
    'ALTER TABLE member_rating_results MODIFY COLUMN additional_accidental_gla_theoretical_rate DECIMAL(15,5);',
    'ALTER TABLE member_rating_results ADD COLUMN additional_accidental_gla_theoretical_rate DECIMAL(15,5);'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='additional_accidental_gla_experience_adjustment'),
    'ALTER TABLE member_rating_results MODIFY COLUMN additional_accidental_gla_experience_adjustment DECIMAL(15,5);',
    'ALTER TABLE member_rating_results ADD COLUMN additional_accidental_gla_experience_adjustment DECIMAL(15,5);'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_loaded_additional_accidental_gla_rate'),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_loaded_additional_accidental_gla_rate DECIMAL(15,5);',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_loaded_additional_accidental_gla_rate DECIMAL(15,5);'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='additional_accidental_gla_risk_premium'),
    'ALTER TABLE member_rating_results MODIFY COLUMN additional_accidental_gla_risk_premium DECIMAL(15,5);',
    'ALTER TABLE member_rating_results ADD COLUMN additional_accidental_gla_risk_premium DECIMAL(15,5);'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_additional_accidental_gla_risk_premium'),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_additional_accidental_gla_risk_premium DECIMAL(15,5);',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_additional_accidental_gla_risk_premium DECIMAL(15,5);'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='additional_accidental_gla_office_premium'),
    'ALTER TABLE member_rating_results MODIFY COLUMN additional_accidental_gla_office_premium DECIMAL(15,5);',
    'ALTER TABLE member_rating_results ADD COLUMN additional_accidental_gla_office_premium DECIMAL(15,5);'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_additional_accidental_gla_office_premium'),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_additional_accidental_gla_office_premium DECIMAL(15,5);',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_additional_accidental_gla_office_premium DECIMAL(15,5);'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
