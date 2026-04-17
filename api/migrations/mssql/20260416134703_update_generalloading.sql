-- Migration for struct: GeneralLoading

-- Table: general_loadings

-- Ensure table exists
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'general_loadings')
BEGIN
    CREATE TABLE general_loadings (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

-- Add or modify column for field: RiskRateCode
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'risk_rate_code')
BEGIN
    ALTER TABLE general_loadings ADD risk_rate_code NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE general_loadings ALTER COLUMN risk_rate_code NVARCHAR(255);
END;

-- Add or modify column for field: Age
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'age')
BEGIN
    ALTER TABLE general_loadings ADD age INT;
END;
ELSE
BEGIN
    ALTER TABLE general_loadings ALTER COLUMN age INT;
END;

-- Add or modify column for field: Gender
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'gender')
BEGIN
    ALTER TABLE general_loadings ADD gender NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE general_loadings ALTER COLUMN gender NVARCHAR(255);
END;

-- Add or modify column for field: GlaContigencyLoadingRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'gla_contigency_loading_rate')
BEGIN
    ALTER TABLE general_loadings ADD gla_contigency_loading_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE general_loadings ALTER COLUMN gla_contigency_loading_rate DECIMAL(15,5);
END;

-- Add or modify column for field: PtdContigencyLoadingRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'ptd_contigency_loading_rate')
BEGIN
    ALTER TABLE general_loadings ADD ptd_contigency_loading_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE general_loadings ALTER COLUMN ptd_contigency_loading_rate DECIMAL(15,5);
END;

-- Add or modify column for field: CiContigencyLoadingRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'ci_contigency_loading_rate')
BEGIN
    ALTER TABLE general_loadings ADD ci_contigency_loading_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE general_loadings ALTER COLUMN ci_contigency_loading_rate DECIMAL(15,5);
END;

-- Add or modify column for field: TtdContigencyLoadingRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'ttd_contigency_loading_rate')
BEGIN
    ALTER TABLE general_loadings ADD ttd_contigency_loading_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE general_loadings ALTER COLUMN ttd_contigency_loading_rate DECIMAL(15,5);
END;

-- Add or modify column for field: PhiContigencyLoadingRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'phi_contigency_loading_rate')
BEGIN
    ALTER TABLE general_loadings ADD phi_contigency_loading_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE general_loadings ALTER COLUMN phi_contigency_loading_rate DECIMAL(15,5);
END;

-- Add or modify column for field: FunContigencyLoadingRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'fun_contigency_loading_rate')
BEGIN
    ALTER TABLE general_loadings ADD fun_contigency_loading_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE general_loadings ALTER COLUMN fun_contigency_loading_rate DECIMAL(15,5);
END;

-- Add or modify column for field: ContinuationLoadingRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'continuation_loading_rate')
BEGIN
    ALTER TABLE general_loadings ADD continuation_loading_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE general_loadings ALTER COLUMN continuation_loading_rate DECIMAL(15,5);
END;

-- Add or modify column for field: TerminalIllnessLoadingRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'terminal_illness_loading_rate')
BEGIN
    ALTER TABLE general_loadings ADD terminal_illness_loading_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE general_loadings ALTER COLUMN terminal_illness_loading_rate DECIMAL(15,5);
END;

-- Add or modify column for field: PtdAcceleratedBenefitDiscount
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'ptd_accelerated_benefit_discount')
BEGIN
    ALTER TABLE general_loadings ADD ptd_accelerated_benefit_discount DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE general_loadings ALTER COLUMN ptd_accelerated_benefit_discount DECIMAL(15,5);
END;

-- Add or modify column for field: CiAcceleratedBenefitDiscount
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'ci_accelerated_benefit_discount')
BEGIN
    ALTER TABLE general_loadings ADD ci_accelerated_benefit_discount DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE general_loadings ALTER COLUMN ci_accelerated_benefit_discount DECIMAL(15,5);
END;

-- Add or modify column for field: CreationDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'creation_date')
BEGIN
    ALTER TABLE general_loadings ADD creation_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE general_loadings ALTER COLUMN creation_date DATETIME2;
END;

-- Add or modify column for field: CreatedBy
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'created_by')
BEGIN
    ALTER TABLE general_loadings ADD created_by NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE general_loadings ALTER COLUMN created_by NVARCHAR(255);
END;

