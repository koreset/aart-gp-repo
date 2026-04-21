-- Migration: make bordereauxes.date_of_birth nullable and clear zero dates.
-- The Bordereaux.DateOfBirth field is now *time.Time so members without a
-- recorded DOB serialise to NULL instead of MySQL's illegal '0000-00-00'.

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='bordereauxes' AND COLUMN_NAME='date_of_birth'),
        'ALTER TABLE bordereauxes MODIFY COLUMN date_of_birth DATETIME NULL DEFAULT NULL;',
        'ALTER TABLE bordereauxes ADD COLUMN date_of_birth DATETIME NULL DEFAULT NULL;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

UPDATE bordereauxes
SET date_of_birth = NULL
WHERE CAST(date_of_birth AS CHAR(19)) IN ('0000-00-00 00:00:00', '0001-01-01 00:00:00');
