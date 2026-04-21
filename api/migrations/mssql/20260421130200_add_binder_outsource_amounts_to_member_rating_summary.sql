-- Migration: add per-benefit binder and outsource fee aggregate columns to member_rating_result_summaries.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_result_summaries')
BEGIN
    CREATE TABLE member_rating_result_summaries (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD total_gla_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD total_gla_annual_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_gla_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_gla_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_gla_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_gla_annual_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_add_acc_gla_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD total_add_acc_gla_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_add_acc_gla_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD total_add_acc_gla_annual_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_add_acc_gla_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_add_acc_gla_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_add_acc_gla_annual_outsourced_amt')
    ALTER TABLE member_rating_result_summaries ADD exp_total_add_acc_gla_annual_outsourced_amt FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_annual_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_ptd_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_ptd_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_ptd_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_ptd_annual_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ci_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD total_ci_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ci_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD total_ci_annual_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_ci_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_ci_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_ci_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_ci_annual_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_sgla_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD total_sgla_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_sgla_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD total_sgla_annual_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_sgla_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_sgla_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_sgla_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_sgla_annual_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ttd_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD total_ttd_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ttd_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD total_ttd_annual_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_ttd_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_ttd_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_ttd_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_ttd_annual_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_phi_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD total_phi_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_phi_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD total_phi_annual_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_phi_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_phi_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_phi_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_phi_annual_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_fun_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD total_fun_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_fun_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD total_fun_annual_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_fun_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_fun_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_fun_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_fun_annual_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_educator_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD total_educator_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_educator_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD total_educator_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_educator_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_educator_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_educator_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_educator_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD total_annual_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD total_annual_outsourced_amount FLOAT;
