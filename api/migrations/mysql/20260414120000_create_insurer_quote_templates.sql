-- Migration for struct: InsurerQuoteTemplate

-- Table: insurer_quote_templates

-- Ensure table exists
CREATE TABLE IF NOT EXISTS insurer_quote_templates (
    id INT AUTO_INCREMENT PRIMARY KEY
);

-- Add or modify column for field: InsurerID
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='insurer_id' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE insurer_quote_templates MODIFY COLUMN insurer_id INT NOT NULL;',
    'ALTER TABLE insurer_quote_templates ADD COLUMN insurer_id INT NOT NULL;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add index on insurer_id
SET @index_name = 'idx_insurer_id';
SET @index_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_NAME='insurer_quote_templates' AND INDEX_NAME=@index_name AND TABLE_SCHEMA = DATABASE());
SET @sql = IF(@index_exists > 0, 'SELECT 1', CONCAT('ALTER TABLE insurer_quote_templates ADD INDEX ', @index_name, ' (insurer_id)'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Version
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='version' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE insurer_quote_templates MODIFY COLUMN version INT;',
    'ALTER TABLE insurer_quote_templates ADD COLUMN version INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Filename
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='filename' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE insurer_quote_templates MODIFY COLUMN filename VARCHAR(255);',
    'ALTER TABLE insurer_quote_templates ADD COLUMN filename VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: DocxBlob
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='docx_blob' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE insurer_quote_templates MODIFY COLUMN docx_blob LONGBLOB;',
    'ALTER TABLE insurer_quote_templates ADD COLUMN docx_blob LONGBLOB;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SizeBytes
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='size_bytes' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE insurer_quote_templates MODIFY COLUMN size_bytes INT;',
    'ALTER TABLE insurer_quote_templates ADD COLUMN size_bytes INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: UploadedBy
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='uploaded_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE insurer_quote_templates MODIFY COLUMN uploaded_by VARCHAR(255);',
    'ALTER TABLE insurer_quote_templates ADD COLUMN uploaded_by VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: UploadedAt
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='uploaded_at' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE insurer_quote_templates MODIFY COLUMN uploaded_at DATETIME;',
    'ALTER TABLE insurer_quote_templates ADD COLUMN uploaded_at DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: IsActive
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='is_active' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE insurer_quote_templates MODIFY COLUMN is_active BOOLEAN;',
    'ALTER TABLE insurer_quote_templates ADD COLUMN is_active BOOLEAN;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add index on is_active
SET @index_name = 'idx_is_active';
SET @index_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_NAME='insurer_quote_templates' AND INDEX_NAME=@index_name AND TABLE_SCHEMA = DATABASE());
SET @sql = IF(@index_exists > 0, 'SELECT 1', CONCAT('ALTER TABLE insurer_quote_templates ADD INDEX ', @index_name, ' (is_active)'));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
