-- Migration for struct: ReinsuranceRegionLoading
-- Table: reinsurance_region_loadings

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'reinsurance_region_loadings')
BEGIN
    CREATE TABLE reinsurance_region_loadings (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_region_loadings' AND COLUMN_NAME = 'risk_rate_code')
    ALTER TABLE reinsurance_region_loadings ADD risk_rate_code NVARCHAR(255);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_region_loadings' AND COLUMN_NAME = 'region')
    ALTER TABLE reinsurance_region_loadings ADD region NVARCHAR(255);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_region_loadings' AND COLUMN_NAME = 'gender')
    ALTER TABLE reinsurance_region_loadings ADD gender NVARCHAR(255);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_region_loadings' AND COLUMN_NAME = 'gla_region_loading_rate')
    ALTER TABLE reinsurance_region_loadings ADD gla_region_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_region_loadings' AND COLUMN_NAME = 'gla_aids_region_loading_rate')
    ALTER TABLE reinsurance_region_loadings ADD gla_aids_region_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_region_loadings' AND COLUMN_NAME = 'ptd_region_loading_rate')
    ALTER TABLE reinsurance_region_loadings ADD ptd_region_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_region_loadings' AND COLUMN_NAME = 'ci_region_loading_rate')
    ALTER TABLE reinsurance_region_loadings ADD ci_region_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_region_loadings' AND COLUMN_NAME = 'ttd_region_loading_rate')
    ALTER TABLE reinsurance_region_loadings ADD ttd_region_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_region_loadings' AND COLUMN_NAME = 'phi_region_loading_rate')
    ALTER TABLE reinsurance_region_loadings ADD phi_region_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_region_loadings' AND COLUMN_NAME = 'fun_region_loading_rate')
    ALTER TABLE reinsurance_region_loadings ADD fun_region_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_region_loadings' AND COLUMN_NAME = 'fun_aids_region_loading_rate')
    ALTER TABLE reinsurance_region_loadings ADD fun_aids_region_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_region_loadings' AND COLUMN_NAME = 'creation_date')
    ALTER TABLE reinsurance_region_loadings ADD creation_date DATETIME;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_region_loadings' AND COLUMN_NAME = 'created_by')
    ALTER TABLE reinsurance_region_loadings ADD created_by NVARCHAR(255);
