-- Migration for struct: GeneralLoading (voluntary loadings)
-- Table: general_loadings

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'gla_voluntary_loading_rate')
BEGIN
    ALTER TABLE general_loadings ADD gla_voluntary_loading_rate DECIMAL(15,5);
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'ptd_voluntary_loading_rate')
BEGIN
    ALTER TABLE general_loadings ADD ptd_voluntary_loading_rate DECIMAL(15,5);
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'ci_voluntary_loading_rate')
BEGIN
    ALTER TABLE general_loadings ADD ci_voluntary_loading_rate DECIMAL(15,5);
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'ttd_voluntary_loading_rate')
BEGIN
    ALTER TABLE general_loadings ADD ttd_voluntary_loading_rate DECIMAL(15,5);
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'phi_voluntary_loading_rate')
BEGIN
    ALTER TABLE general_loadings ADD phi_voluntary_loading_rate DECIMAL(15,5);
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'fun_voluntary_loading_rate')
BEGIN
    ALTER TABLE general_loadings ADD fun_voluntary_loading_rate DECIMAL(15,5);
END;
