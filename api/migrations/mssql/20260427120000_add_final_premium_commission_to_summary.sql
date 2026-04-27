-- Migration: add Discount + Final* office-premium and commission columns to
-- member_rating_result_summaries. Final*OfficePremium is persisted as the
-- post-discount pre-comm office premium (Exp*RiskPremium /
-- (1 - (SchemeTotalLoading + Discount))) plus its per-benefit commission
-- slice, so the Final* values include commission and reconcile to
-- final_total_annual_premium{,_excl_funeral}. Exp* values stay pre-commission
-- and frozen at calc time, so pre-discount Final - Exp == commission slice.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_result_summaries')
BEGIN
    CREATE TABLE member_rating_result_summaries (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='discount')
    ALTER TABLE member_rating_result_summaries ADD discount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_annual_office_premium')
    ALTER TABLE member_rating_result_summaries ADD final_gla_annual_office_premium FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_add_acc_gla_annual_office_premium')
    ALTER TABLE member_rating_result_summaries ADD final_add_acc_gla_annual_office_premium FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_annual_office_premium')
    ALTER TABLE member_rating_result_summaries ADD final_ptd_annual_office_premium FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_annual_office_premium')
    ALTER TABLE member_rating_result_summaries ADD final_ci_annual_office_premium FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_annual_office_premium')
    ALTER TABLE member_rating_result_summaries ADD final_sgla_annual_office_premium FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_tax_saver_annual_office_premium')
    ALTER TABLE member_rating_result_summaries ADD final_tax_saver_annual_office_premium FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_annual_office_premium')
    ALTER TABLE member_rating_result_summaries ADD final_ttd_annual_office_premium FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_annual_office_premium')
    ALTER TABLE member_rating_result_summaries ADD final_phi_annual_office_premium FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_annual_office_premium')
    ALTER TABLE member_rating_result_summaries ADD final_fun_annual_office_premium FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_educator_annual_office_premium')
    ALTER TABLE member_rating_result_summaries ADD final_gla_educator_annual_office_premium FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_educator_annual_office_premium')
    ALTER TABLE member_rating_result_summaries ADD final_ptd_educator_annual_office_premium FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD final_gla_annual_commission_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_add_acc_gla_annual_comm_amount')
    ALTER TABLE member_rating_result_summaries ADD final_add_acc_gla_annual_comm_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD final_ptd_annual_commission_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD final_ci_annual_commission_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD final_sgla_annual_commission_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_tax_saver_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD final_tax_saver_annual_commission_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD final_ttd_annual_commission_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD final_phi_annual_commission_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD final_fun_annual_commission_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_educator_annual_comm_amount')
    ALTER TABLE member_rating_result_summaries ADD final_gla_educator_annual_comm_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_educator_annual_comm_amount')
    ALTER TABLE member_rating_result_summaries ADD final_ptd_educator_annual_comm_amount FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_total_annual_premium_excl_funeral')
    ALTER TABLE member_rating_result_summaries ADD final_total_annual_premium_excl_funeral FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_total_annual_premium')
    ALTER TABLE member_rating_result_summaries ADD final_total_annual_premium FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_scheme_total_commission')
    ALTER TABLE member_rating_result_summaries ADD final_scheme_total_commission FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_scheme_total_commission_rate')
    ALTER TABLE member_rating_result_summaries ADD final_scheme_total_commission_rate FLOAT DEFAULT 0;
