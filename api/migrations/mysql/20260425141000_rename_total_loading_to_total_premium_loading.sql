-- Migration: align the member_rating_results column name with the GORM
-- struct field `TotalPremiumLoading` (default snake_case → total_premium_loading)
-- and the JSON/CSV tags. Earlier migrations created the column as
-- `total_loading`, which caused INSERTs to fail with
-- "Unknown column 'total_premium_loading' in 'field list'" on databases that
-- were provisioned via the manual migrations rather than GORM AutoMigrate.

-- Three-way decision:
--   1. canonical `total_premium_loading` already exists → no-op
--   2. legacy `total_loading` exists (and canonical does not) → rename it
--   3. neither exists (fresh DB that skipped the broken migration) → add it
SET @s = (SELECT
    IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='total_premium_loading' AND TABLE_SCHEMA = DATABASE()),
        'SELECT 1;',
        IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='total_loading' AND TABLE_SCHEMA = DATABASE()),
            'ALTER TABLE member_rating_results CHANGE COLUMN total_loading total_premium_loading DOUBLE;',
            'ALTER TABLE member_rating_results ADD COLUMN total_premium_loading DOUBLE;')));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
