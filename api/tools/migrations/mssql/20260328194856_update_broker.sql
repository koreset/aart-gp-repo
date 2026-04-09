-- Migration for struct: Broker

-- Table: brokers

-- Ensure table exists
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'brokers')
BEGIN
    CREATE TABLE brokers (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

-- Add or modify column for field: Name
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'brokers' AND COLUMN_NAME = 'name')
BEGIN
    ALTER TABLE brokers ADD name NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE brokers ALTER COLUMN name NVARCHAR(255);
END;

-- Add or modify column for field: ContactEmail
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'brokers' AND COLUMN_NAME = 'contact_email')
BEGIN
    ALTER TABLE brokers ADD contact_email NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE brokers ALTER COLUMN contact_email NVARCHAR(255);
END;

-- Add or modify column for field: ContactNumber
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'brokers' AND COLUMN_NAME = 'contact_number')
BEGIN
    ALTER TABLE brokers ADD contact_number NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE brokers ALTER COLUMN contact_number NVARCHAR(255);
END;

-- Add or modify column for field: FSPNumber
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'brokers' AND COLUMN_NAME = 'fsp_number')
BEGIN
    ALTER TABLE brokers ADD fsp_number NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE brokers ALTER COLUMN fsp_number NVARCHAR(255);
END;

-- Add or modify column for field: FSPCategory
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'brokers' AND COLUMN_NAME = 'fsp_category')
BEGIN
    ALTER TABLE brokers ADD fsp_category NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE brokers ALTER COLUMN fsp_category NVARCHAR(255);
END;

-- Add or modify column for field: BinderAgreementRef
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'brokers' AND COLUMN_NAME = 'binder_agreement_ref')
BEGIN
    ALTER TABLE brokers ADD binder_agreement_ref NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE brokers ALTER COLUMN binder_agreement_ref NVARCHAR(255);
END;

-- Add or modify column for field: TiedAgentRef
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'brokers' AND COLUMN_NAME = 'tied_agent_ref')
BEGIN
    ALTER TABLE brokers ADD tied_agent_ref NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE brokers ALTER COLUMN tied_agent_ref NVARCHAR(255);
END;

-- Add or modify column for field: CreationDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'brokers' AND COLUMN_NAME = 'creation_date')
BEGIN
    ALTER TABLE brokers ADD creation_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE brokers ALTER COLUMN creation_date DATETIME2;
END;

-- Add or modify column for field: CreatedBy
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'brokers' AND COLUMN_NAME = 'created_by')
BEGIN
    ALTER TABLE brokers ADD created_by NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE brokers ALTER COLUMN created_by NVARCHAR(255);
END;

