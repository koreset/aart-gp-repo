-- Migration: add binder fee and outsource fee columns to group_pricing_quotes.

CREATE TABLE IF NOT EXISTS group_pricing_quotes (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='loadings_binder_fee' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE group_pricing_quotes MODIFY COLUMN loadings_binder_fee DOUBLE DEFAULT 0;', 'ALTER TABLE group_pricing_quotes ADD COLUMN loadings_binder_fee DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='loadings_outsource_fee' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE group_pricing_quotes MODIFY COLUMN loadings_outsource_fee DOUBLE DEFAULT 0;', 'ALTER TABLE group_pricing_quotes ADD COLUMN loadings_outsource_fee DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
