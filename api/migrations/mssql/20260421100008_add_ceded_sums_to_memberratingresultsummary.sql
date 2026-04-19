-- Migration: add ceded sum assured / ceded monthly benefit aggregates to
-- member_rating_result_summaries.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_result_summaries')
BEGIN
    CREATE TABLE member_rating_result_summaries (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_ceded_sum_assured')
    ALTER TABLE member_rating_result_summaries ADD total_gla_ceded_sum_assured FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_ceded_sum_assured')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_ceded_sum_assured FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ci_ceded_sum_assured')
    ALTER TABLE member_rating_result_summaries ADD total_ci_ceded_sum_assured FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_sgla_ceded_sum_assured')
    ALTER TABLE member_rating_result_summaries ADD total_sgla_ceded_sum_assured FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ttd_ceded_monthly_benefit')
    ALTER TABLE member_rating_result_summaries ADD total_ttd_ceded_monthly_benefit FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_phi_ceded_monthly_benefit')
    ALTER TABLE member_rating_result_summaries ADD total_phi_ceded_monthly_benefit FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_fun_ceded_sum_assured')
    ALTER TABLE member_rating_result_summaries ADD total_fun_ceded_sum_assured FLOAT;
