-- Migration: persist binder-fee and outsource-fee rates on each member rating
-- result summary row alongside the existing scheme-level loadings.

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='binder_fee_rate' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN binder_fee_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='outsource_fee_rate' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN outsource_fee_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
