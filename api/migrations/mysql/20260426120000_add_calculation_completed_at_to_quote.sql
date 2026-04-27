-- Migration: add `calculation_completed_at` to group_pricing_quotes so the
-- UI can display when the underlying calculation last produced results.
-- Nullable on purpose: pre-existing quotes have no such timestamp until
-- they are recalculated.

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS
           WHERE TABLE_NAME='group_pricing_quotes'
             AND COLUMN_NAME='calculation_completed_at'
             AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN calculation_completed_at DATETIME NULL;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
