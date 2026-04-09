-- Migration for struct: TaxTable

-- Table: tax_tables

-- Ensure table exists
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'tax_tables')
BEGIN
    CREATE TABLE tax_tables (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

-- Add or modify column for field: RiskRateCode
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'tax_tables' AND COLUMN_NAME = 'risk_rate_code')
BEGIN
    ALTER TABLE tax_tables ADD risk_rate_code NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE tax_tables ALTER COLUMN risk_rate_code NVARCHAR(255);
END;

-- Add or modify column for field: Level
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'tax_tables' AND COLUMN_NAME = 'level')
BEGIN
    ALTER TABLE tax_tables ADD level INT;
END;
ELSE
BEGIN
    ALTER TABLE tax_tables ALTER COLUMN level INT;
END;

-- Add or modify column for field: Min
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'tax_tables' AND COLUMN_NAME = 'min')
BEGIN
    ALTER TABLE tax_tables ADD min DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE tax_tables ALTER COLUMN min DECIMAL(15,5);
END;

-- Add or modify column for field: Max
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'tax_tables' AND COLUMN_NAME = 'max')
BEGIN
    ALTER TABLE tax_tables ADD max DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE tax_tables ALTER COLUMN max DECIMAL(15,5);
END;

-- Add or modify column for field: TaxRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'tax_tables' AND COLUMN_NAME = 'tax_rate')
BEGIN
    ALTER TABLE tax_tables ADD tax_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE tax_tables ALTER COLUMN tax_rate DECIMAL(15,5);
END;

-- Add or modify column for field: CreationDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'tax_tables' AND COLUMN_NAME = 'creation_date')
BEGIN
    ALTER TABLE tax_tables ADD creation_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE tax_tables ALTER COLUMN creation_date DATETIME2;
END;

-- Add or modify column for field: CreatedBy
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'tax_tables' AND COLUMN_NAME = 'created_by')
BEGIN
    ALTER TABLE tax_tables ADD created_by NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE tax_tables ALTER COLUMN created_by NVARCHAR(255);
END;

