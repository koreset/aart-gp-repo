-- Migration for struct: PremiumLoading

-- Table: premium_loadings

-- Ensure table exists
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'premium_loadings')
BEGIN
    CREATE TABLE premium_loadings (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

-- Add or modify column for field: RiskRateCode
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'premium_loadings' AND COLUMN_NAME = 'risk_rate_code')
BEGIN
    ALTER TABLE premium_loadings ADD risk_rate_code NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE premium_loadings ALTER COLUMN risk_rate_code NVARCHAR(255);
END;

-- Add or modify column for field: Channel
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'premium_loadings' AND COLUMN_NAME = 'channel')
BEGIN
    ALTER TABLE premium_loadings ADD channel NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE premium_loadings ALTER COLUMN channel NVARCHAR(255);
END;

-- Add or modify column for field: SchemeSizeLevel
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'premium_loadings' AND COLUMN_NAME = 'scheme_size_level')
BEGIN
    ALTER TABLE premium_loadings ADD scheme_size_level INT;
END;
ELSE
BEGIN
    ALTER TABLE premium_loadings ALTER COLUMN scheme_size_level INT;
END;

-- Add or modify column for field: CommissionLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'premium_loadings' AND COLUMN_NAME = 'commission_loading')
BEGIN
    ALTER TABLE premium_loadings ADD commission_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE premium_loadings ALTER COLUMN commission_loading DECIMAL(15,5);
END;

-- Add or modify column for field: ExpenseLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'premium_loadings' AND COLUMN_NAME = 'expense_loading')
BEGIN
    ALTER TABLE premium_loadings ADD expense_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE premium_loadings ALTER COLUMN expense_loading DECIMAL(15,5);
END;

-- Add or modify column for field: AdminLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'premium_loadings' AND COLUMN_NAME = 'admin_loading')
BEGIN
    ALTER TABLE premium_loadings ADD admin_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE premium_loadings ALTER COLUMN admin_loading DECIMAL(15,5);
END;

-- Add or modify column for field: OtherLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'premium_loadings' AND COLUMN_NAME = 'other_loading')
BEGIN
    ALTER TABLE premium_loadings ADD other_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE premium_loadings ALTER COLUMN other_loading DECIMAL(15,5);
END;

-- Add or modify column for field: ProfitMargin
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'premium_loadings' AND COLUMN_NAME = 'profit_margin')
BEGIN
    ALTER TABLE premium_loadings ADD profit_margin DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE premium_loadings ALTER COLUMN profit_margin DECIMAL(15,5);
END;

-- Add or modify column for field: CreationDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'premium_loadings' AND COLUMN_NAME = 'creation_date')
BEGIN
    ALTER TABLE premium_loadings ADD creation_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE premium_loadings ALTER COLUMN creation_date DATETIME2;
END;

-- Add or modify column for field: CreatedBy
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'premium_loadings' AND COLUMN_NAME = 'created_by')
BEGIN
    ALTER TABLE premium_loadings ADD created_by NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE premium_loadings ALTER COLUMN created_by NVARCHAR(255);
END;

