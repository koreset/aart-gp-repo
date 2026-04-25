-- Migration: align member_rating_results funeral columns with the Go model.
-- The Go model went through a "Cost → RiskPremium / OfficePremium" rename
-- and a "dependant(s) → parent" rename, plus added FinalTotalFuneralOfficePremium.
-- The previous migrations only updated comments, not column names.

-- Helper procedure used to rename a column only if the source exists and the
-- target does not. MySQL 8.0+ supports `RENAME COLUMN`.
DROP PROCEDURE IF EXISTS __mrr_rename_col;
DELIMITER $$
CREATE PROCEDURE __mrr_rename_col(IN old_name VARCHAR(64), IN new_name VARCHAR(64))
BEGIN
    DECLARE old_exists INT DEFAULT 0;
    DECLARE new_exists INT DEFAULT 0;
    SELECT COUNT(*) INTO old_exists
    FROM information_schema.columns
    WHERE table_schema = DATABASE() AND table_name = 'member_rating_results' AND column_name = old_name;
    SELECT COUNT(*) INTO new_exists
    FROM information_schema.columns
    WHERE table_schema = DATABASE() AND table_name = 'member_rating_results' AND column_name = new_name;
    IF old_exists = 1 AND new_exists = 0 THEN
        SET @sql := CONCAT('ALTER TABLE member_rating_results RENAME COLUMN `', old_name, '` TO `', new_name, '`');
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
    END IF;
END$$
DELIMITER ;

CALL __mrr_rename_col('main_member_funeral_cost', 'main_member_funeral_risk_premium');
CALL __mrr_rename_col('spouse_funeral_cost', 'spouse_funeral_risk_premium');
CALL __mrr_rename_col('children_funeral_cost', 'child_funeral_risk_premium');
CALL __mrr_rename_col('dependants_funeral_cost', 'parent_funeral_risk_premium');
CALL __mrr_rename_col('total_funeral_risk_cost', 'total_funeral_risk_premium');
CALL __mrr_rename_col('exp_adj_total_funeral_risk_cost', 'exp_adj_total_funeral_risk_premium');
CALL __mrr_rename_col('total_funeral_office_cost', 'total_funeral_office_premium');
CALL __mrr_rename_col('exp_adj_total_funeral_office_cost', 'exp_adj_total_funeral_office_premium');
CALL __mrr_rename_col('dependant_funeral_base_rate', 'parent_funeral_base_rate');

DROP PROCEDURE IF EXISTS __mrr_rename_col;

-- Add the new column the model introduced but no migration created.
SET @sql := IF(
    (SELECT COUNT(*) FROM information_schema.columns
        WHERE table_schema = DATABASE() AND table_name = 'member_rating_results' AND column_name = 'final_total_funeral_office_premium') = 0,
    'ALTER TABLE member_rating_results ADD COLUMN final_total_funeral_office_premium DECIMAL(15,5)',
    'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Drop obsolete dependant_funeral_sum_assured (parent_funeral_sum_assured exists already).
SET @sql := IF(
    (SELECT COUNT(*) FROM information_schema.columns
        WHERE table_schema = DATABASE() AND table_name = 'member_rating_results' AND column_name = 'dependant_funeral_sum_assured') = 1,
    'ALTER TABLE member_rating_results DROP COLUMN dependant_funeral_sum_assured',
    'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
