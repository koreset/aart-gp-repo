-- Migration: move educator benefit selection from group_pricing_parameters
-- to the per-scheme-category level, and enforce uniqueness on the
-- educator benefit structure.

-- ---------------------------------------------------------------------------
-- 1. scheme_categories: new per-benefit educator type columns
-- ---------------------------------------------------------------------------
SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='gla_educator_benefit_type'),
        'ALTER TABLE scheme_categories MODIFY COLUMN gla_educator_benefit_type VARCHAR(255)',
        'ALTER TABLE scheme_categories ADD COLUMN gla_educator_benefit_type VARCHAR(255)'
    )
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='ptd_educator_benefit_type'),
        'ALTER TABLE scheme_categories MODIFY COLUMN ptd_educator_benefit_type VARCHAR(255)',
        'ALTER TABLE scheme_categories ADD COLUMN ptd_educator_benefit_type VARCHAR(255)'
    )
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ---------------------------------------------------------------------------
-- 2. educator_benefit_structures: composite unique index on (rrc, code)
-- ---------------------------------------------------------------------------
-- Convert risk_rate_code and educator_benefit_code to VARCHAR(255) first.
-- MySQL cannot index TEXT/LONGTEXT columns without a prefix length, and the
-- columns in older databases were auto-created as LONGTEXT by GORM.
ALTER TABLE educator_benefit_structures MODIFY COLUMN risk_rate_code VARCHAR(255);
ALTER TABLE educator_benefit_structures MODIFY COLUMN educator_benefit_code VARCHAR(255);

-- Guard against pre-existing duplicates; fail loudly if any are present.
SET @dup_count := (
    SELECT COUNT(*) FROM (
        SELECT 1
        FROM educator_benefit_structures
        GROUP BY risk_rate_code, educator_benefit_code
        HAVING COUNT(*) > 1
    ) AS d
);
SET @msg := CONCAT('educator_benefit_structures has ', @dup_count, ' duplicate (risk_rate_code, educator_benefit_code) group(s); dedupe before re-running.');
SET @sql := IF(@dup_count > 0, CONCAT("SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '", @msg, "'"), 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='educator_benefit_structures' AND INDEX_NAME='idx_educator_rrc_code'),
        'SELECT 1',
        'CREATE UNIQUE INDEX idx_educator_rrc_code ON educator_benefit_structures (risk_rate_code, educator_benefit_code)'
    )
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ---------------------------------------------------------------------------
-- 3. group_pricing_parameters: drop educator_benefit_code
-- ---------------------------------------------------------------------------
SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='group_pricing_parameters' AND COLUMN_NAME='educator_benefit_code'),
        'ALTER TABLE group_pricing_parameters DROP COLUMN educator_benefit_code',
        'SELECT 1'
    )
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
