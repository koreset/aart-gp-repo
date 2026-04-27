-- Migration: add Final*BinderAmount + Final*OutsourcedAmount columns to
-- member_rating_result_summaries. These mirror the existing Exp*Binder /
-- Exp*Outsourced fields but are derived from the post-discount Final office
-- premium so the breakdown reconciles to FinalOfficePremium after a discount
-- is applied. Pre-discount they equal the Exp* values; on non-binder
-- distribution channels they remain 0.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_result_summaries')
BEGIN
    CREATE TABLE member_rating_result_summaries (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD final_gla_annual_binder_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD final_gla_annual_outsourced_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_add_acc_gla_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD final_add_acc_gla_annual_binder_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_add_acc_gla_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD final_add_acc_gla_annual_outsourced_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD final_ptd_annual_binder_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD final_ptd_annual_outsourced_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD final_ci_annual_binder_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD final_ci_annual_outsourced_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD final_sgla_annual_binder_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD final_sgla_annual_outsourced_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_tax_saver_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD final_tax_saver_annual_binder_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_tax_saver_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD final_tax_saver_annual_outsourced_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD final_ttd_annual_binder_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD final_ttd_annual_outsourced_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD final_phi_annual_binder_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD final_phi_annual_outsourced_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD final_fun_annual_binder_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD final_fun_annual_outsourced_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_educator_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD final_gla_educator_annual_binder_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_educator_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD final_gla_educator_annual_outsourced_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_educator_annual_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD final_ptd_educator_annual_binder_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_educator_annual_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD final_ptd_educator_annual_outsourced_amount FLOAT DEFAULT 0;
