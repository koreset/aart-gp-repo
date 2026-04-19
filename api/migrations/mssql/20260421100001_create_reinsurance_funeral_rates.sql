-- Migration for struct: ReinsuranceFuneralRate
-- Table: reinsurance_funeral_rates

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'reinsurance_funeral_rates')
BEGIN
    CREATE TABLE reinsurance_funeral_rates (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_funeral_rates' AND COLUMN_NAME = 'risk_rate_code')
    ALTER TABLE reinsurance_funeral_rates ADD risk_rate_code NVARCHAR(255);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_funeral_rates' AND COLUMN_NAME = 'age_next_birthday')
    ALTER TABLE reinsurance_funeral_rates ADD age_next_birthday INT;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_funeral_rates' AND COLUMN_NAME = 'gender')
    ALTER TABLE reinsurance_funeral_rates ADD gender NVARCHAR(255);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_funeral_rates' AND COLUMN_NAME = 'fun_qx')
    ALTER TABLE reinsurance_funeral_rates ADD fun_qx DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_funeral_rates' AND COLUMN_NAME = 'creation_date')
    ALTER TABLE reinsurance_funeral_rates ADD creation_date DATETIME;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_funeral_rates' AND COLUMN_NAME = 'created_by')
    ALTER TABLE reinsurance_funeral_rates ADD created_by NVARCHAR(255);
