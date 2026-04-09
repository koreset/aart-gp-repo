-- Migration for struct: SchemeSizeLevel

-- Table: scheme_size_levels

-- Ensure table exists
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'scheme_size_levels')
BEGIN
    CREATE TABLE scheme_size_levels (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

-- Add or modify column for field: RiskRateCode
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_size_levels' AND COLUMN_NAME = 'risk_rate_code')
BEGIN
    ALTER TABLE scheme_size_levels ADD risk_rate_code NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_size_levels ALTER COLUMN risk_rate_code NVARCHAR(255);
END;

-- Add or modify column for field: MinCount
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_size_levels' AND COLUMN_NAME = 'min_count')
BEGIN
    ALTER TABLE scheme_size_levels ADD min_count INT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_size_levels ALTER COLUMN min_count INT;
END;

-- Add or modify column for field: MaxCount
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_size_levels' AND COLUMN_NAME = 'max_count')
BEGIN
    ALTER TABLE scheme_size_levels ADD max_count INT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_size_levels ALTER COLUMN max_count INT;
END;

-- Add or modify column for field: SizeLevel
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_size_levels' AND COLUMN_NAME = 'size_level')
BEGIN
    ALTER TABLE scheme_size_levels ADD size_level INT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_size_levels ALTER COLUMN size_level INT;
END;

-- Add or modify column for field: CreationDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_size_levels' AND COLUMN_NAME = 'creation_date')
BEGIN
    ALTER TABLE scheme_size_levels ADD creation_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE scheme_size_levels ALTER COLUMN creation_date DATETIME2;
END;

-- Add or modify column for field: CreatedBy
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_size_levels' AND COLUMN_NAME = 'created_by')
BEGIN
    ALTER TABLE scheme_size_levels ADD created_by NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_size_levels ALTER COLUMN created_by NVARCHAR(255);
END;

