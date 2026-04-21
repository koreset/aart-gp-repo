-- Migration: add optional TaxSaver rider fields. The tax-saver loading
-- itself lives on general_loadings (per age/gender/risk_rate_code), so
-- scheme_categories only carries the opt-in flag.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'scheme_categories')
BEGIN
    CREATE TABLE scheme_categories (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'tax_saver_benefit')
    ALTER TABLE scheme_categories ADD tax_saver_benefit BIT DEFAULT 0;

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'general_loadings')
BEGIN
    CREATE TABLE general_loadings (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'tax_saver_loading_rate')
    ALTER TABLE general_loadings ADD tax_saver_loading_rate FLOAT;

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_results')
BEGIN
    CREATE TABLE member_rating_results (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'tax_saver_loading')
    ALTER TABLE member_rating_results ADD tax_saver_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'tax_saver_risk_premium')
    ALTER TABLE member_rating_results ADD tax_saver_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_tax_saver_risk_premium')
    ALTER TABLE member_rating_results ADD exp_adj_tax_saver_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'tax_saver_office_premium')
    ALTER TABLE member_rating_results ADD tax_saver_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_tax_saver_office_premium')
    ALTER TABLE member_rating_results ADD exp_adj_tax_saver_office_premium FLOAT;

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_result_summaries')
BEGIN
    CREATE TABLE member_rating_result_summaries (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'tax_saver_benefit')
    ALTER TABLE member_rating_result_summaries ADD tax_saver_benefit BIT DEFAULT 0;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_tax_saver_annual_risk_premium')
    ALTER TABLE member_rating_result_summaries ADD total_tax_saver_annual_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_tax_saver_annual_office_premium')
    ALTER TABLE member_rating_result_summaries ADD total_tax_saver_annual_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_tax_saver_annual_risk_premium')
    ALTER TABLE member_rating_result_summaries ADD exp_total_tax_saver_annual_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_tax_saver_annual_office_premium')
    ALTER TABLE member_rating_result_summaries ADD exp_total_tax_saver_annual_office_premium FLOAT;
