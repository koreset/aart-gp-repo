-- Migration: add request_payload to generated_bordereauxes so draft rows can
-- be regenerated in place using the original request parameters.

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'generated_bordereauxes' AND COLUMN_NAME = 'request_payload'),
        'SELECT 1;',
        'ALTER TABLE generated_bordereauxes ADD COLUMN request_payload JSON NULL;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
