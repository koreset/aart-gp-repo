-- Migration for struct: ReinsuranceIndustryLoading
-- Table: reinsurance_industry_loadings

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'reinsurance_industry_loadings')
BEGIN
    CREATE TABLE reinsurance_industry_loadings (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_industry_loadings' AND COLUMN_NAME = 'risk_rate_code')
    ALTER TABLE reinsurance_industry_loadings ADD risk_rate_code NVARCHAR(255);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_industry_loadings' AND COLUMN_NAME = 'occupation_class')
    ALTER TABLE reinsurance_industry_loadings ADD occupation_class INT;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_industry_loadings' AND COLUMN_NAME = 'gender')
    ALTER TABLE reinsurance_industry_loadings ADD gender NVARCHAR(255);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_industry_loadings' AND COLUMN_NAME = 'gla_industry_loading_rate')
    ALTER TABLE reinsurance_industry_loadings ADD gla_industry_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_industry_loadings' AND COLUMN_NAME = 'ptd_industry_loading_rate')
    ALTER TABLE reinsurance_industry_loadings ADD ptd_industry_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_industry_loadings' AND COLUMN_NAME = 'ci_industry_loading_rate')
    ALTER TABLE reinsurance_industry_loadings ADD ci_industry_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_industry_loadings' AND COLUMN_NAME = 'ttd_industry_loading_rate')
    ALTER TABLE reinsurance_industry_loadings ADD ttd_industry_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_industry_loadings' AND COLUMN_NAME = 'phi_industry_loading_rate')
    ALTER TABLE reinsurance_industry_loadings ADD phi_industry_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_industry_loadings' AND COLUMN_NAME = 'creation_date')
    ALTER TABLE reinsurance_industry_loadings ADD creation_date DATETIME;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_industry_loadings' AND COLUMN_NAME = 'created_by')
    ALTER TABLE reinsurance_industry_loadings ADD created_by NVARCHAR(255);
