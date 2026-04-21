-- Migration: add optional TaxSaver rider fields. The tax-saver loading
-- itself lives on general_loadings (per age/gender/risk_rate_code), so
-- scheme_categories only carries the opt-in flag.

CREATE TABLE IF NOT EXISTS scheme_categories (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='tax_saver_benefit' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE scheme_categories MODIFY COLUMN tax_saver_benefit TINYINT(1) DEFAULT 0;', 'ALTER TABLE scheme_categories ADD COLUMN tax_saver_benefit TINYINT(1) DEFAULT 0;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS general_loadings (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='general_loadings' AND COLUMN_NAME='tax_saver_loading_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE general_loadings MODIFY COLUMN tax_saver_loading_rate DOUBLE;', 'ALTER TABLE general_loadings ADD COLUMN tax_saver_loading_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS member_rating_results (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='tax_saver_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN tax_saver_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN tax_saver_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='tax_saver_risk_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN tax_saver_risk_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN tax_saver_risk_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_tax_saver_risk_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_tax_saver_risk_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_tax_saver_risk_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='tax_saver_office_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN tax_saver_office_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN tax_saver_office_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_tax_saver_office_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_tax_saver_office_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_tax_saver_office_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='tax_saver_benefit' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN tax_saver_benefit TINYINT(1) DEFAULT 0;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN tax_saver_benefit TINYINT(1) DEFAULT 0;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_tax_saver_annual_risk_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_tax_saver_annual_risk_premium DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_tax_saver_annual_risk_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_tax_saver_annual_office_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_tax_saver_annual_office_premium DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_tax_saver_annual_office_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_tax_saver_annual_risk_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_tax_saver_annual_risk_premium DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_tax_saver_annual_risk_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_tax_saver_annual_office_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_tax_saver_annual_office_premium DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_tax_saver_annual_office_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
