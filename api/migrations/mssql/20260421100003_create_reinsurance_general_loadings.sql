-- Migration for struct: ReinsuranceGeneralLoading
-- Table: reinsurance_general_loadings

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'reinsurance_general_loadings')
BEGIN
    CREATE TABLE reinsurance_general_loadings (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_general_loadings' AND COLUMN_NAME = 'risk_rate_code')
    ALTER TABLE reinsurance_general_loadings ADD risk_rate_code NVARCHAR(255);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_general_loadings' AND COLUMN_NAME = 'age')
    ALTER TABLE reinsurance_general_loadings ADD age INT;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_general_loadings' AND COLUMN_NAME = 'gender')
    ALTER TABLE reinsurance_general_loadings ADD gender NVARCHAR(255);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_general_loadings' AND COLUMN_NAME = 'gla_contigency_loading_rate')
    ALTER TABLE reinsurance_general_loadings ADD gla_contigency_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_general_loadings' AND COLUMN_NAME = 'ptd_contigency_loading_rate')
    ALTER TABLE reinsurance_general_loadings ADD ptd_contigency_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_general_loadings' AND COLUMN_NAME = 'ci_contigency_loading_rate')
    ALTER TABLE reinsurance_general_loadings ADD ci_contigency_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_general_loadings' AND COLUMN_NAME = 'ttd_contigency_loading_rate')
    ALTER TABLE reinsurance_general_loadings ADD ttd_contigency_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_general_loadings' AND COLUMN_NAME = 'phi_contigency_loading_rate')
    ALTER TABLE reinsurance_general_loadings ADD phi_contigency_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_general_loadings' AND COLUMN_NAME = 'fun_contigency_loading_rate')
    ALTER TABLE reinsurance_general_loadings ADD fun_contigency_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_general_loadings' AND COLUMN_NAME = 'continuation_loading_rate')
    ALTER TABLE reinsurance_general_loadings ADD continuation_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_general_loadings' AND COLUMN_NAME = 'terminal_illness_loading_rate')
    ALTER TABLE reinsurance_general_loadings ADD terminal_illness_loading_rate DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_general_loadings' AND COLUMN_NAME = 'ptd_accelerated_benefit_discount')
    ALTER TABLE reinsurance_general_loadings ADD ptd_accelerated_benefit_discount DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_general_loadings' AND COLUMN_NAME = 'ci_accelerated_benefit_discount')
    ALTER TABLE reinsurance_general_loadings ADD ci_accelerated_benefit_discount DECIMAL(15,5);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_general_loadings' AND COLUMN_NAME = 'creation_date')
    ALTER TABLE reinsurance_general_loadings ADD creation_date DATETIME;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'reinsurance_general_loadings' AND COLUMN_NAME = 'created_by')
    ALTER TABLE reinsurance_general_loadings ADD created_by NVARCHAR(255);
