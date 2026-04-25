-- Migration: add the non-Exp Total*CommissionAmount mirrors plus the per-category
-- ExpTotalCommission roll-up to member_rating_result_summaries. Parallel to the
-- existing Exp*CommissionAmount columns added in 20260421150000 so the scheme-wide
-- commission can be distributed on both the Exp and Total premium bases.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_result_summaries')
BEGIN
    CREATE TABLE member_rating_result_summaries (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_annual_premium_excl_funeral')
    ALTER TABLE member_rating_result_summaries ADD total_annual_premium_excl_funeral FLOAT DEFAULT 0;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_annual_premium_excl_funeral')
    ALTER TABLE member_rating_result_summaries ADD exp_total_annual_premium_excl_funeral FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_commission')
    ALTER TABLE member_rating_result_summaries ADD exp_total_commission FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD total_gla_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_add_acc_gla_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD total_add_acc_gla_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ci_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD total_ci_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_sgla_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD total_sgla_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ttd_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD total_ttd_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_phi_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD total_phi_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_fun_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD total_fun_annual_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_educator_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD total_gla_educator_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_educator_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_educator_commission_amount FLOAT;
