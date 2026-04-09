-- Migration for struct: ValUserRole

-- Table: val_user_roles

-- Ensure table exists
CREATE TABLE IF NOT EXISTS val_user_roles (
    id INT AUTO_INCREMENT PRIMARY KEY
);

-- Add or modify column for field: RoleName
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='val_user_roles' AND COLUMN_NAME='role_name' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE val_user_roles MODIFY COLUMN role_name VARCHAR(255);',
    'ALTER TABLE val_user_roles ADD COLUMN role_name VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Description
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='val_user_roles' AND COLUMN_NAME='description' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE val_user_roles MODIFY COLUMN description VARCHAR(255);',
    'ALTER TABLE val_user_roles ADD COLUMN description VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Permissions
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='val_user_roles' AND COLUMN_NAME='permissions' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE val_user_roles MODIFY COLUMN permissions TEXT;',
    'ALTER TABLE val_user_roles ADD COLUMN permissions TEXT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

