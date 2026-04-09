-- Migration for struct: Broker

-- Table: brokers

-- Ensure table exists
CREATE TABLE IF NOT EXISTS brokers (
    id INT AUTO_INCREMENT PRIMARY KEY
);

-- Add or modify column for field: Name
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='brokers' AND COLUMN_NAME='name' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE brokers MODIFY COLUMN name VARCHAR(255);',
    'ALTER TABLE brokers ADD COLUMN name VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ContactEmail
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='brokers' AND COLUMN_NAME='contact_email' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE brokers MODIFY COLUMN contact_email VARCHAR(255);',
    'ALTER TABLE brokers ADD COLUMN contact_email VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ContactNumber
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='brokers' AND COLUMN_NAME='contact_number' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE brokers MODIFY COLUMN contact_number VARCHAR(255);',
    'ALTER TABLE brokers ADD COLUMN contact_number VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: FSPNumber
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='brokers' AND COLUMN_NAME='fsp_number' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE brokers MODIFY COLUMN fsp_number VARCHAR(255);',
    'ALTER TABLE brokers ADD COLUMN fsp_number VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: FSPCategory
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='brokers' AND COLUMN_NAME='fsp_category' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE brokers MODIFY COLUMN fsp_category VARCHAR(255);',
    'ALTER TABLE brokers ADD COLUMN fsp_category VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: BinderAgreementRef
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='brokers' AND COLUMN_NAME='binder_agreement_ref' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE brokers MODIFY COLUMN binder_agreement_ref VARCHAR(255);',
    'ALTER TABLE brokers ADD COLUMN binder_agreement_ref VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TiedAgentRef
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='brokers' AND COLUMN_NAME='tied_agent_ref' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE brokers MODIFY COLUMN tied_agent_ref VARCHAR(255);',
    'ALTER TABLE brokers ADD COLUMN tied_agent_ref VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreationDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='brokers' AND COLUMN_NAME='creation_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE brokers MODIFY COLUMN creation_date DATETIME;',
    'ALTER TABLE brokers ADD COLUMN creation_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreatedBy
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='brokers' AND COLUMN_NAME='created_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE brokers MODIFY COLUMN created_by VARCHAR(255);',
    'ALTER TABLE brokers ADD COLUMN created_by VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

