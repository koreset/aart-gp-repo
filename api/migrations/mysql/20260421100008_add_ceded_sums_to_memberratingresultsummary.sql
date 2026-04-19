-- Migration: add ceded sum assured / ceded monthly benefit aggregates to
-- member_rating_result_summaries so the Reinsurance Premium Summary view can
-- surface totals per benefit, per category.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_ceded_sum_assured' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_gla_ceded_sum_assured DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_gla_ceded_sum_assured DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ptd_ceded_sum_assured' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ptd_ceded_sum_assured DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ptd_ceded_sum_assured DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ci_ceded_sum_assured' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ci_ceded_sum_assured DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ci_ceded_sum_assured DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_sgla_ceded_sum_assured' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_sgla_ceded_sum_assured DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_sgla_ceded_sum_assured DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ttd_ceded_monthly_benefit' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ttd_ceded_monthly_benefit DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ttd_ceded_monthly_benefit DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_phi_ceded_monthly_benefit' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_phi_ceded_monthly_benefit DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_phi_ceded_monthly_benefit DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_fun_ceded_sum_assured' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_fun_ceded_sum_assured DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_fun_ceded_sum_assured DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

