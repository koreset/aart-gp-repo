-- Migration: add reinsurance premium aggregates & proportions to member_rating_result_summaries.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_result_summaries')
BEGIN
    CREATE TABLE member_rating_result_summaries (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_reinsurance_premium')
    ALTER TABLE member_rating_result_summaries ADD total_gla_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_reinsurance_premium')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ci_reinsurance_premium')
    ALTER TABLE member_rating_result_summaries ADD total_ci_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_sgla_reinsurance_premium')
    ALTER TABLE member_rating_result_summaries ADD total_sgla_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_phi_reinsurance_premium')
    ALTER TABLE member_rating_result_summaries ADD total_phi_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ttd_reinsurance_premium')
    ALTER TABLE member_rating_result_summaries ADD total_ttd_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_fun_reinsurance_premium')
    ALTER TABLE member_rating_result_summaries ADD total_fun_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'gla_reinsurance_premium_proportion')
    ALTER TABLE member_rating_result_summaries ADD gla_reinsurance_premium_proportion FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'ptd_reinsurance_premium_proportion')
    ALTER TABLE member_rating_result_summaries ADD ptd_reinsurance_premium_proportion FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'ci_reinsurance_premium_proportion')
    ALTER TABLE member_rating_result_summaries ADD ci_reinsurance_premium_proportion FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'sgla_reinsurance_premium_proportion')
    ALTER TABLE member_rating_result_summaries ADD sgla_reinsurance_premium_proportion FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'phi_reinsurance_premium_proportion')
    ALTER TABLE member_rating_result_summaries ADD phi_reinsurance_premium_proportion FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'ttd_reinsurance_premium_proportion')
    ALTER TABLE member_rating_result_summaries ADD ttd_reinsurance_premium_proportion FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'fun_reinsurance_premium_proportion')
    ALTER TABLE member_rating_result_summaries ADD fun_reinsurance_premium_proportion FLOAT;
