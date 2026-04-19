-- Migration for struct: TableConfigurationAuditLog
-- Table: table_configuration_audit_logs
--
-- Append-only audit history of every change to TableConfiguration.IsRequired.

CREATE TABLE IF NOT EXISTS table_configuration_audit_logs (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='table_configuration_audit_logs' AND COLUMN_NAME='table_type' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE table_configuration_audit_logs MODIFY COLUMN table_type VARCHAR(128);',
    'ALTER TABLE table_configuration_audit_logs ADD COLUMN table_type VARCHAR(128);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='table_configuration_audit_logs' AND COLUMN_NAME='event_type' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE table_configuration_audit_logs MODIFY COLUMN event_type VARCHAR(64);',
    'ALTER TABLE table_configuration_audit_logs ADD COLUMN event_type VARCHAR(64);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='table_configuration_audit_logs' AND COLUMN_NAME='old_value' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE table_configuration_audit_logs MODIFY COLUMN old_value TINYINT(1);',
    'ALTER TABLE table_configuration_audit_logs ADD COLUMN old_value TINYINT(1);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='table_configuration_audit_logs' AND COLUMN_NAME='new_value' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE table_configuration_audit_logs MODIFY COLUMN new_value TINYINT(1);',
    'ALTER TABLE table_configuration_audit_logs ADD COLUMN new_value TINYINT(1);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='table_configuration_audit_logs' AND COLUMN_NAME='changed_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE table_configuration_audit_logs MODIFY COLUMN changed_by VARCHAR(255);',
    'ALTER TABLE table_configuration_audit_logs ADD COLUMN changed_by VARCHAR(255);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='table_configuration_audit_logs' AND COLUMN_NAME='changed_at' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE table_configuration_audit_logs MODIFY COLUMN changed_at DATETIME;',
    'ALTER TABLE table_configuration_audit_logs ADD COLUMN changed_at DATETIME;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='table_configuration_audit_logs' AND COLUMN_NAME='details' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE table_configuration_audit_logs MODIFY COLUMN details TEXT;',
    'ALTER TABLE table_configuration_audit_logs ADD COLUMN details TEXT;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_NAME='table_configuration_audit_logs' AND INDEX_NAME='idx_table_configuration_audit_logs_table_type' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'CREATE INDEX idx_table_configuration_audit_logs_table_type ON table_configuration_audit_logs (table_type);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
