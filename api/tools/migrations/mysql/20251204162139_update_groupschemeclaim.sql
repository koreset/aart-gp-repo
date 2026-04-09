-- Migration for struct: GroupSchemeClaim

-- Table: group_scheme_claims

-- Ensure table exists
CREATE TABLE IF NOT EXISTS group_scheme_claims (
    id INT AUTO_INCREMENT PRIMARY KEY
);

-- Add or modify column for field: ClaimNumber
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='claim_number' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN claim_number VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN claim_number VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MemberIDNumber
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='member_id_number' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN member_id_number VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN member_id_number VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MemberName
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='member_name' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN member_name VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN member_name VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SchemeId
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='scheme_id' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN scheme_id INT;',
    'ALTER TABLE group_scheme_claims ADD COLUMN scheme_id INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SchemeName
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='scheme_name' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN scheme_name VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN scheme_name VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: BenefitType
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='benefit_type' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN benefit_type VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN benefit_type VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ClaimType
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='claim_type' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN claim_type VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN claim_type VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: DateOfEvent
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='date_of_event' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN date_of_event VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN date_of_event VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: DateNotified
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='date_notified' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN date_notified VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN date_notified VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ClaimAmount
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='claim_amount' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN claim_amount DOUBLE;',
    'ALTER TABLE group_scheme_claims ADD COLUMN claim_amount DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Priority
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='priority' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN priority VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN priority VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ClaimantName
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='claimant_name' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN claimant_name VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN claimant_name VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ClaimantIDNumber
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='claimant_id_number' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN claimant_id_number VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN claimant_id_number VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: RelationshipToMember
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='relationship_to_member' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN relationship_to_member VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN relationship_to_member VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ClaimantContactNumber
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='claimant_contact_number' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN claimant_contact_number VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN claimant_contact_number VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Description
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='description' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN description VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN description VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Status
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='status' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN status VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN status VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: DateRegistered
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='date_registered' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN date_registered VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN date_registered VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreationDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='creation_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN creation_date DATETIME;',
    'ALTER TABLE group_scheme_claims ADD COLUMN creation_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreatedBy
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='created_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN created_by VARCHAR(255);',
    'ALTER TABLE group_scheme_claims ADD COLUMN created_by VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Attachments
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='attachments' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN attachments TEXT;',
    'ALTER TABLE group_scheme_claims ADD COLUMN attachments TEXT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Assessments
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_scheme_claims' AND COLUMN_NAME='assessments' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_scheme_claims MODIFY COLUMN assessments TEXT;',
    'ALTER TABLE group_scheme_claims ADD COLUMN assessments TEXT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

