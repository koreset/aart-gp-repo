-- Migration: add admin_loading + other_loading columns to
-- member_rating_result_summaries. These let SchemeTotalLoading() include the
-- full premium-loading sum (expense + profit + admin + other + binder +
-- outsource) at the summary level, matching the rating-phase
-- TotalPremiumLoading on MemberRatingResult.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='admin_loading' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN admin_loading DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='other_loading' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN other_loading DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
