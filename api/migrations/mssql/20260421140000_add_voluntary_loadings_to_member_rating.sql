-- Migration: add per-benefit voluntary loading columns to member_rating_results.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_results')
BEGIN
    CREATE TABLE member_rating_results (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_voluntary_loading')
    ALTER TABLE member_rating_results ADD gla_voluntary_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_voluntary_loading')
    ALTER TABLE member_rating_results ADD ptd_voluntary_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_voluntary_loading')
    ALTER TABLE member_rating_results ADD ci_voluntary_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_voluntary_loading')
    ALTER TABLE member_rating_results ADD ttd_voluntary_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_voluntary_loading')
    ALTER TABLE member_rating_results ADD phi_voluntary_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'fun_voluntary_loading')
    ALTER TABLE member_rating_results ADD fun_voluntary_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_gla_voluntary_loading')
    ALTER TABLE member_rating_results ADD reins_gla_voluntary_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_ptd_voluntary_loading')
    ALTER TABLE member_rating_results ADD reins_ptd_voluntary_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_ci_voluntary_loading')
    ALTER TABLE member_rating_results ADD reins_ci_voluntary_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_ttd_voluntary_loading')
    ALTER TABLE member_rating_results ADD reins_ttd_voluntary_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_phi_voluntary_loading')
    ALTER TABLE member_rating_results ADD reins_phi_voluntary_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_fun_voluntary_loading')
    ALTER TABLE member_rating_results ADD reins_fun_voluntary_loading FLOAT;
