-- Migration: add per-benefit binder and outsource fee amount columns to member_rating_results.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_results')
BEGIN
    CREATE TABLE member_rating_results (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_binder_amount')
    ALTER TABLE member_rating_results ADD gla_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_outsourced_amount')
    ALTER TABLE member_rating_results ADD gla_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_binder_amount')
    ALTER TABLE member_rating_results ADD exp_adj_gla_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_outsourced_amount')
    ALTER TABLE member_rating_results ADD exp_adj_gla_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'additional_accidental_gla_binder_amount')
    ALTER TABLE member_rating_results ADD additional_accidental_gla_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'additional_accidental_gla_outsourced_amount')
    ALTER TABLE member_rating_results ADD additional_accidental_gla_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_add_acc_gla_binder_amount')
    ALTER TABLE member_rating_results ADD exp_adj_add_acc_gla_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_add_acc_gla_outsourced_amount')
    ALTER TABLE member_rating_results ADD exp_adj_add_acc_gla_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_binder_amount')
    ALTER TABLE member_rating_results ADD ptd_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_outsourced_amount')
    ALTER TABLE member_rating_results ADD ptd_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ptd_binder_amount')
    ALTER TABLE member_rating_results ADD exp_adj_ptd_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ptd_outsourced_amount')
    ALTER TABLE member_rating_results ADD exp_adj_ptd_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_binder_amount')
    ALTER TABLE member_rating_results ADD ci_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_outsourced_amount')
    ALTER TABLE member_rating_results ADD ci_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ci_binder_amount')
    ALTER TABLE member_rating_results ADD exp_adj_ci_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ci_outsourced_amount')
    ALTER TABLE member_rating_results ADD exp_adj_ci_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_gla_binder_amount')
    ALTER TABLE member_rating_results ADD spouse_gla_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_gla_outsourced_amount')
    ALTER TABLE member_rating_results ADD spouse_gla_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_spouse_gla_binder_amount')
    ALTER TABLE member_rating_results ADD exp_adj_spouse_gla_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_spouse_gla_outsourced_amount')
    ALTER TABLE member_rating_results ADD exp_adj_spouse_gla_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_binder_amount')
    ALTER TABLE member_rating_results ADD ttd_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_outsourced_amount')
    ALTER TABLE member_rating_results ADD ttd_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ttd_binder_amount')
    ALTER TABLE member_rating_results ADD exp_adj_ttd_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ttd_outsourced_amount')
    ALTER TABLE member_rating_results ADD exp_adj_ttd_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_binder_amount')
    ALTER TABLE member_rating_results ADD phi_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_outsourced_amount')
    ALTER TABLE member_rating_results ADD phi_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_phi_binder_amount')
    ALTER TABLE member_rating_results ADD exp_adj_phi_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_phi_outsourced_amount')
    ALTER TABLE member_rating_results ADD exp_adj_phi_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'main_member_funeral_binder_amount')
    ALTER TABLE member_rating_results ADD main_member_funeral_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'main_member_funeral_outsourced_amount')
    ALTER TABLE member_rating_results ADD main_member_funeral_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_funeral_binder_amount')
    ALTER TABLE member_rating_results ADD spouse_funeral_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_funeral_outsourced_amount')
    ALTER TABLE member_rating_results ADD spouse_funeral_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'children_funeral_binder_amount')
    ALTER TABLE member_rating_results ADD children_funeral_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'children_funeral_outsourced_amount')
    ALTER TABLE member_rating_results ADD children_funeral_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'dependants_funeral_binder_amount')
    ALTER TABLE member_rating_results ADD dependants_funeral_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'dependants_funeral_outsourced_amount')
    ALTER TABLE member_rating_results ADD dependants_funeral_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'total_funeral_binder_amount')
    ALTER TABLE member_rating_results ADD total_funeral_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'total_funeral_outsourced_amount')
    ALTER TABLE member_rating_results ADD total_funeral_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_total_funeral_binder_amount')
    ALTER TABLE member_rating_results ADD exp_adj_total_funeral_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_total_funeral_outsourced_amount')
    ALTER TABLE member_rating_results ADD exp_adj_total_funeral_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'educator_binder_amount')
    ALTER TABLE member_rating_results ADD educator_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'educator_outsourced_amount')
    ALTER TABLE member_rating_results ADD educator_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_educator_binder_amount')
    ALTER TABLE member_rating_results ADD exp_adj_educator_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_educator_outsourced_amount')
    ALTER TABLE member_rating_results ADD exp_adj_educator_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'total_binder_amount')
    ALTER TABLE member_rating_results ADD total_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'total_outsourced_amount')
    ALTER TABLE member_rating_results ADD total_outsourced_amount FLOAT;
