-- Migration for struct: GroupSchemeClaim

-- Table: group_scheme_claims

-- Ensure table exists
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'group_scheme_claims')
BEGIN
    CREATE TABLE group_scheme_claims (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

-- Add or modify column for field: ClaimNumber
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'claim_number')
BEGIN
    ALTER TABLE group_scheme_claims ADD claim_number NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN claim_number NVARCHAR(255);
END;

-- Add or modify column for field: MemberIDNumber
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'member_id_number')
BEGIN
    ALTER TABLE group_scheme_claims ADD member_id_number NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN member_id_number NVARCHAR(255);
END;

-- Add or modify column for field: MemberName
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'member_name')
BEGIN
    ALTER TABLE group_scheme_claims ADD member_name NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN member_name NVARCHAR(255);
END;

-- Add or modify column for field: SchemeId
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'scheme_id')
BEGIN
    ALTER TABLE group_scheme_claims ADD scheme_id INT;
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN scheme_id INT;
END;

-- Add or modify column for field: SchemeName
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'scheme_name')
BEGIN
    ALTER TABLE group_scheme_claims ADD scheme_name NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN scheme_name NVARCHAR(255);
END;

-- Add or modify column for field: BenefitType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'benefit_type')
BEGIN
    ALTER TABLE group_scheme_claims ADD benefit_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN benefit_type NVARCHAR(255);
END;

-- Add or modify column for field: ClaimType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'claim_type')
BEGIN
    ALTER TABLE group_scheme_claims ADD claim_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN claim_type NVARCHAR(255);
END;

-- Add or modify column for field: DateOfEvent
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'date_of_event')
BEGIN
    ALTER TABLE group_scheme_claims ADD date_of_event NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN date_of_event NVARCHAR(255);
END;

-- Add or modify column for field: DateNotified
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'date_notified')
BEGIN
    ALTER TABLE group_scheme_claims ADD date_notified NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN date_notified NVARCHAR(255);
END;

-- Add or modify column for field: ClaimAmount
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'claim_amount')
BEGIN
    ALTER TABLE group_scheme_claims ADD claim_amount DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN claim_amount DECIMAL(15,5);
END;

-- Add or modify column for field: Priority
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'priority')
BEGIN
    ALTER TABLE group_scheme_claims ADD priority NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN priority NVARCHAR(255);
END;

-- Add or modify column for field: ClaimantName
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'claimant_name')
BEGIN
    ALTER TABLE group_scheme_claims ADD claimant_name NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN claimant_name NVARCHAR(255);
END;

-- Add or modify column for field: ClaimantIDNumber
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'claimant_id_number')
BEGIN
    ALTER TABLE group_scheme_claims ADD claimant_id_number NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN claimant_id_number NVARCHAR(255);
END;

-- Add or modify column for field: RelationshipToMember
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'relationship_to_member')
BEGIN
    ALTER TABLE group_scheme_claims ADD relationship_to_member NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN relationship_to_member NVARCHAR(255);
END;

-- Add or modify column for field: ClaimantContactNumber
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'claimant_contact_number')
BEGIN
    ALTER TABLE group_scheme_claims ADD claimant_contact_number NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN claimant_contact_number NVARCHAR(255);
END;

-- Add or modify column for field: Description
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'description')
BEGIN
    ALTER TABLE group_scheme_claims ADD description NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN description NVARCHAR(255);
END;

-- Add or modify column for field: Status
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'status')
BEGIN
    ALTER TABLE group_scheme_claims ADD status NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN status NVARCHAR(255);
END;

-- Add or modify column for field: DateRegistered
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'date_registered')
BEGIN
    ALTER TABLE group_scheme_claims ADD date_registered NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN date_registered NVARCHAR(255);
END;

-- Add or modify column for field: CreationDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'creation_date')
BEGIN
    ALTER TABLE group_scheme_claims ADD creation_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN creation_date DATETIME2;
END;

-- Add or modify column for field: CreatedBy
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'created_by')
BEGIN
    ALTER TABLE group_scheme_claims ADD created_by NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN created_by NVARCHAR(255);
END;

-- Add or modify column for field: Attachments
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'attachments')
BEGIN
    ALTER TABLE group_scheme_claims ADD attachments NVARCHAR(MAX);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN attachments NVARCHAR(MAX);
END;

-- Add or modify column for field: Assessments
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_scheme_claims' AND COLUMN_NAME = 'assessments')
BEGIN
    ALTER TABLE group_scheme_claims ADD assessments NVARCHAR(MAX);
END;
ELSE
BEGIN
    ALTER TABLE group_scheme_claims ALTER COLUMN assessments NVARCHAR(MAX);
END;

