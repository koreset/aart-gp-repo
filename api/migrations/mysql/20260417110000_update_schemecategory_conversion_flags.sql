-- Migration: add per-benefit conversion flags to scheme_categories.
--   gla_conversion_on_withdrawal  (GLA)
--   gla_conversion_on_retirement  (GLA, retirement-only)
--   ptd_conversion_on_withdrawal  (PTD)
--   ci_conversion_on_withdrawal   (CI)

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='gla_conversion_on_withdrawal'),
        'ALTER TABLE scheme_categories MODIFY COLUMN gla_conversion_on_withdrawal TINYINT(1) DEFAULT 0;',
        'ALTER TABLE scheme_categories ADD COLUMN gla_conversion_on_withdrawal TINYINT(1) DEFAULT 0;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='gla_conversion_on_retirement'),
        'ALTER TABLE scheme_categories MODIFY COLUMN gla_conversion_on_retirement TINYINT(1) DEFAULT 0;',
        'ALTER TABLE scheme_categories ADD COLUMN gla_conversion_on_retirement TINYINT(1) DEFAULT 0;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='ptd_conversion_on_withdrawal'),
        'ALTER TABLE scheme_categories MODIFY COLUMN ptd_conversion_on_withdrawal TINYINT(1) DEFAULT 0;',
        'ALTER TABLE scheme_categories ADD COLUMN ptd_conversion_on_withdrawal TINYINT(1) DEFAULT 0;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='ci_conversion_on_withdrawal'),
        'ALTER TABLE scheme_categories MODIFY COLUMN ci_conversion_on_withdrawal TINYINT(1) DEFAULT 0;',
        'ALTER TABLE scheme_categories ADD COLUMN ci_conversion_on_withdrawal TINYINT(1) DEFAULT 0;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
