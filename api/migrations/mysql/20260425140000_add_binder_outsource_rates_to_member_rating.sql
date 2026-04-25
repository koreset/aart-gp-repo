-- Migration: persist binder-fee and outsource-fee rates on each member rating
-- result row so the per-member loading breakdown is auditable end-to-end.

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='binder_fee_rate' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1;', 'ALTER TABLE member_rating_results ADD COLUMN binder_fee_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='outsource_fee_rate' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1;', 'ALTER TABLE member_rating_results ADD COLUMN outsource_fee_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
