-- Migration: add per-benefit commission amount columns and scheme-wide commission totals to member_rating_result_summaries.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_result_summaries')
BEGIN
    CREATE TABLE member_rating_result_summaries (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_gla_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_gla_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_add_acc_gla_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_add_acc_gla_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_ptd_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_ptd_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_ci_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_ci_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_sgla_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_sgla_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_ttd_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_ttd_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_phi_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_phi_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_fun_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_fun_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_educator_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_educator_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'scheme_total_commission')
    ALTER TABLE member_rating_result_summaries ADD scheme_total_commission FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'scheme_total_commission_rate')
    ALTER TABLE member_rating_result_summaries ADD scheme_total_commission_rate FLOAT;
