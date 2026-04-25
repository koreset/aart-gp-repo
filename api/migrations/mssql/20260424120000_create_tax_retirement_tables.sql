-- Migration for struct: TaxRetirementTable
-- Table: tax_retirement_tables

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'tax_retirement_tables')
BEGIN
    CREATE TABLE tax_retirement_tables (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'tax_retirement_tables' AND COLUMN_NAME = 'risk_rate_code')
    ALTER TABLE tax_retirement_tables ADD risk_rate_code NVARCHAR(255);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'tax_retirement_tables' AND COLUMN_NAME = 'lower_bound')
    ALTER TABLE tax_retirement_tables ADD lower_bound DECIMAL(18,4);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'tax_retirement_tables' AND COLUMN_NAME = 'upper_bound')
    ALTER TABLE tax_retirement_tables ADD upper_bound DECIMAL(18,4);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'tax_retirement_tables' AND COLUMN_NAME = 'tax_rate')
    ALTER TABLE tax_retirement_tables ADD tax_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'tax_retirement_tables' AND COLUMN_NAME = 'cumulative_tax_relief')
    ALTER TABLE tax_retirement_tables ADD cumulative_tax_relief DECIMAL(18,4);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'tax_retirement_tables' AND COLUMN_NAME = 'creation_date')
    ALTER TABLE tax_retirement_tables ADD creation_date DATETIME;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'tax_retirement_tables' AND COLUMN_NAME = 'created_by')
    ALTER TABLE tax_retirement_tables ADD created_by NVARCHAR(255);
