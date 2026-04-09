-- Migration for struct: OrgUser

-- Table: org_users

-- Ensure table exists
CREATE TABLE IF NOT EXISTS org_users (
    id INT AUTO_INCREMENT PRIMARY KEY
);

-- Add or modify column for field: Name
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='org_users' AND COLUMN_NAME='name' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE org_users MODIFY COLUMN name VARCHAR(255);',
    'ALTER TABLE org_users ADD COLUMN name VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Email
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='org_users' AND COLUMN_NAME='email' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE org_users MODIFY COLUMN email VARCHAR(255);',
    'ALTER TABLE org_users ADD COLUMN email VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: LicenseId
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='org_users' AND COLUMN_NAME='license_id' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE org_users MODIFY COLUMN license_id VARCHAR(255);',
    'ALTER TABLE org_users ADD COLUMN license_id VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Role
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='org_users' AND COLUMN_NAME='role' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE org_users MODIFY COLUMN role VARCHAR(255);',
    'ALTER TABLE org_users ADD COLUMN role VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GPRole
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='org_users' AND COLUMN_NAME='gp_role' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE org_users MODIFY COLUMN gp_role VARCHAR(255);',
    'ALTER TABLE org_users ADD COLUMN gp_role VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GPRoleId
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='org_users' AND COLUMN_NAME='gp_role_id' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE org_users MODIFY COLUMN gp_role_id INT;',
    'ALTER TABLE org_users ADD COLUMN gp_role_id INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ValRole
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='org_users' AND COLUMN_NAME='val_role' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE org_users MODIFY COLUMN val_role VARCHAR(255);',
    'ALTER TABLE org_users ADD COLUMN val_role VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ValRoleId
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='org_users' AND COLUMN_NAME='val_role_id' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE org_users MODIFY COLUMN val_role_id INT;',
    'ALTER TABLE org_users ADD COLUMN val_role_id INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Organisation
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='org_users' AND COLUMN_NAME='organisation' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE org_users MODIFY COLUMN organisation VARCHAR(255);',
    'ALTER TABLE org_users ADD COLUMN organisation VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

