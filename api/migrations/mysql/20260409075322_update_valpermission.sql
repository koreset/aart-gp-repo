-- Migration for struct: ValPermission

-- Table: val_permissions

-- Ensure table exists
CREATE TABLE IF NOT EXISTS val_permissions (
    id INT AUTO_INCREMENT PRIMARY KEY
);

-- Add or modify column for field: Name
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='val_permissions' AND COLUMN_NAME='name' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE val_permissions MODIFY COLUMN name VARCHAR(255);',
    'ALTER TABLE val_permissions ADD COLUMN name VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Slug
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='val_permissions' AND COLUMN_NAME='slug' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE val_permissions MODIFY COLUMN slug VARCHAR(255);',
    'ALTER TABLE val_permissions ADD COLUMN slug VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Description
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='val_permissions' AND COLUMN_NAME='description' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE val_permissions MODIFY COLUMN description VARCHAR(255);',
    'ALTER TABLE val_permissions ADD COLUMN description VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

