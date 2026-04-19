-- Migration for struct: TableConfiguration
-- Table: table_configurations
--
-- Per-table-type "is required" flag for Group Pricing rating tables.
-- Rows are seeded at startup by services.EnsureTableConfigurations().

CREATE TABLE IF NOT EXISTS table_configurations (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='table_configurations' AND COLUMN_NAME='table_type' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE table_configurations MODIFY COLUMN table_type VARCHAR(128);',
    'ALTER TABLE table_configurations ADD COLUMN table_type VARCHAR(128);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='table_configurations' AND COLUMN_NAME='display_name' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE table_configurations MODIFY COLUMN display_name VARCHAR(255);',
    'ALTER TABLE table_configurations ADD COLUMN display_name VARCHAR(255);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='table_configurations' AND COLUMN_NAME='category' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE table_configurations MODIFY COLUMN category VARCHAR(64);',
    'ALTER TABLE table_configurations ADD COLUMN category VARCHAR(64);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='table_configurations' AND COLUMN_NAME='is_required' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE table_configurations MODIFY COLUMN is_required TINYINT(1) DEFAULT 1;',
    'ALTER TABLE table_configurations ADD COLUMN is_required TINYINT(1) DEFAULT 1;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='table_configurations' AND COLUMN_NAME='updated_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE table_configurations MODIFY COLUMN updated_by VARCHAR(255);',
    'ALTER TABLE table_configurations ADD COLUMN updated_by VARCHAR(255);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='table_configurations' AND COLUMN_NAME='updated_at' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE table_configurations MODIFY COLUMN updated_at DATETIME;',
    'ALTER TABLE table_configurations ADD COLUMN updated_at DATETIME;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='table_configurations' AND COLUMN_NAME='created_at' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE table_configurations MODIFY COLUMN created_at DATETIME;',
    'ALTER TABLE table_configurations ADD COLUMN created_at DATETIME;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_NAME='table_configurations' AND INDEX_NAME='idx_table_configurations_table_type' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'CREATE UNIQUE INDEX idx_table_configurations_table_type ON table_configurations (table_type);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
