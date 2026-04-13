-- Migration for struct: CalculationJob

-- Table: calculation_jobs

-- Ensure table exists
CREATE TABLE IF NOT EXISTS calculation_jobs (
    id INT AUTO_INCREMENT PRIMARY KEY
);

-- Add or modify column for field: QuoteID
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='calculation_jobs' AND COLUMN_NAME='quote_id' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE calculation_jobs MODIFY COLUMN quote_id INT;',
    'ALTER TABLE calculation_jobs ADD COLUMN quote_id INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Basis
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='calculation_jobs' AND COLUMN_NAME='basis' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE calculation_jobs MODIFY COLUMN basis VARCHAR(255);',
    'ALTER TABLE calculation_jobs ADD COLUMN basis VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Credibility
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='calculation_jobs' AND COLUMN_NAME='credibility' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE calculation_jobs MODIFY COLUMN credibility DOUBLE;',
    'ALTER TABLE calculation_jobs ADD COLUMN credibility DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: UserEmail
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='calculation_jobs' AND COLUMN_NAME='user_email' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE calculation_jobs MODIFY COLUMN user_email VARCHAR(255);',
    'ALTER TABLE calculation_jobs ADD COLUMN user_email VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: UserName
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='calculation_jobs' AND COLUMN_NAME='user_name' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE calculation_jobs MODIFY COLUMN user_name VARCHAR(255);',
    'ALTER TABLE calculation_jobs ADD COLUMN user_name VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Status
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='calculation_jobs' AND COLUMN_NAME='status' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE calculation_jobs MODIFY COLUMN status VARCHAR(255);',
    'ALTER TABLE calculation_jobs ADD COLUMN status VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Error
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='calculation_jobs' AND COLUMN_NAME='error' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE calculation_jobs MODIFY COLUMN error VARCHAR(255);',
    'ALTER TABLE calculation_jobs ADD COLUMN error VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: QueuedAt
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='calculation_jobs' AND COLUMN_NAME='queued_at' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE calculation_jobs MODIFY COLUMN queued_at DATETIME;',
    'ALTER TABLE calculation_jobs ADD COLUMN queued_at DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: StartedAt
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='calculation_jobs' AND COLUMN_NAME='started_at' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE calculation_jobs MODIFY COLUMN started_at DATETIME;',
    'ALTER TABLE calculation_jobs ADD COLUMN started_at DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CompletedAt
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='calculation_jobs' AND COLUMN_NAME='completed_at' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE calculation_jobs MODIFY COLUMN completed_at DATETIME;',
    'ALTER TABLE calculation_jobs ADD COLUMN completed_at DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

