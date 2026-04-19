-- Migration: add reinsurance rate/loading/premium fields to member_rating_results
-- Adds the 62 columns introduced alongside the reinsurance rate tables so that
-- GORM INSERTs from MemberRatingResult succeed against existing DBs.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_results')
BEGIN
    CREATE TABLE member_rating_results (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_gla_industry_loading')
    ALTER TABLE member_rating_results ADD reins_gla_industry_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_ptd_industry_loading')
    ALTER TABLE member_rating_results ADD reins_ptd_industry_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_ci_industry_loading')
    ALTER TABLE member_rating_results ADD reins_ci_industry_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_ttd_industry_loading')
    ALTER TABLE member_rating_results ADD reins_ttd_industry_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_phi_industry_loading')
    ALTER TABLE member_rating_results ADD reins_phi_industry_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_gla_region_loading')
    ALTER TABLE member_rating_results ADD reins_gla_region_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_gla_aids_region_loading')
    ALTER TABLE member_rating_results ADD reins_gla_aids_region_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_ptd_region_loading')
    ALTER TABLE member_rating_results ADD reins_ptd_region_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_ci_region_loading')
    ALTER TABLE member_rating_results ADD reins_ci_region_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_ttd_region_loading')
    ALTER TABLE member_rating_results ADD reins_ttd_region_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_phi_region_loading')
    ALTER TABLE member_rating_results ADD reins_phi_region_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_fun_region_loading')
    ALTER TABLE member_rating_results ADD reins_fun_region_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_fun_aids_region_loading')
    ALTER TABLE member_rating_results ADD reins_fun_aids_region_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_gla_contingency_loading')
    ALTER TABLE member_rating_results ADD reins_gla_contingency_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_ptd_contingency_loading')
    ALTER TABLE member_rating_results ADD reins_ptd_contingency_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_ci_contingency_loading')
    ALTER TABLE member_rating_results ADD reins_ci_contingency_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_ttd_contingency_loading')
    ALTER TABLE member_rating_results ADD reins_ttd_contingency_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_phi_contingency_loading')
    ALTER TABLE member_rating_results ADD reins_phi_contingency_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_fun_contingency_loading')
    ALTER TABLE member_rating_results ADD reins_fun_contingency_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_continuation_loading')
    ALTER TABLE member_rating_results ADD reins_continuation_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_gla_terminal_illness_loading')
    ALTER TABLE member_rating_results ADD reins_gla_terminal_illness_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_gla_qx')
    ALTER TABLE member_rating_results ADD reins_gla_qx FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_gla_aids_qx')
    ALTER TABLE member_rating_results ADD reins_gla_aids_qx FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'base_reins_gla_rate')
    ALTER TABLE member_rating_results ADD base_reins_gla_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'loaded_reins_gla_rate')
    ALTER TABLE member_rating_results ADD loaded_reins_gla_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_ptd_rate')
    ALTER TABLE member_rating_results ADD reins_ptd_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'base_reins_ptd_rate')
    ALTER TABLE member_rating_results ADD base_reins_ptd_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'loaded_reins_ptd_rate')
    ALTER TABLE member_rating_results ADD loaded_reins_ptd_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_ci_rate')
    ALTER TABLE member_rating_results ADD reins_ci_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'base_reins_ci_rate')
    ALTER TABLE member_rating_results ADD base_reins_ci_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'loaded_reins_ci_rate')
    ALTER TABLE member_rating_results ADD loaded_reins_ci_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'base_reins_ttd_rate')
    ALTER TABLE member_rating_results ADD base_reins_ttd_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'loaded_reins_ttd_rate')
    ALTER TABLE member_rating_results ADD loaded_reins_ttd_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_phi_rate')
    ALTER TABLE member_rating_results ADD reins_phi_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'base_reins_phi_rate')
    ALTER TABLE member_rating_results ADD base_reins_phi_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'loaded_reins_phi_rate')
    ALTER TABLE member_rating_results ADD loaded_reins_phi_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_spouse_gla_qx')
    ALTER TABLE member_rating_results ADD reins_spouse_gla_qx FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_spouse_gla_aids_qx')
    ALTER TABLE member_rating_results ADD reins_spouse_gla_aids_qx FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'reins_spouse_gla_loading')
    ALTER TABLE member_rating_results ADD reins_spouse_gla_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'base_reins_spouse_gla_rate')
    ALTER TABLE member_rating_results ADD base_reins_spouse_gla_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'loaded_reins_spouse_gla_rate')
    ALTER TABLE member_rating_results ADD loaded_reins_spouse_gla_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'main_member_reinsurance_base_rate')
    ALTER TABLE member_rating_results ADD main_member_reinsurance_base_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'main_member_reinsurance_rate')
    ALTER TABLE member_rating_results ADD main_member_reinsurance_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_reinsurance_base_rate')
    ALTER TABLE member_rating_results ADD spouse_reinsurance_base_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_reinsurance_rate')
    ALTER TABLE member_rating_results ADD spouse_reinsurance_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'child_reinsurance_base_rate')
    ALTER TABLE member_rating_results ADD child_reinsurance_base_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'child_reinsurance_rate')
    ALTER TABLE member_rating_results ADD child_reinsurance_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'parent_reinsurance_base_rate')
    ALTER TABLE member_rating_results ADD parent_reinsurance_base_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'parent_reinsurance_rate')
    ALTER TABLE member_rating_results ADD parent_reinsurance_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'dependant_reinsurance_base_rate')
    ALTER TABLE member_rating_results ADD dependant_reinsurance_base_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'dependant_reinsurance_rate')
    ALTER TABLE member_rating_results ADD dependant_reinsurance_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_reinsurance_premium')
    ALTER TABLE member_rating_results ADD gla_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_reinsurance_premium')
    ALTER TABLE member_rating_results ADD ptd_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_reinsurance_premium')
    ALTER TABLE member_rating_results ADD ci_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_gla_reinsurance_premium')
    ALTER TABLE member_rating_results ADD spouse_gla_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_reinsurance_premium')
    ALTER TABLE member_rating_results ADD ttd_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_reinsurance_premium')
    ALTER TABLE member_rating_results ADD phi_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'main_member_reinsurance_premium')
    ALTER TABLE member_rating_results ADD main_member_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_reinsurance_premium')
    ALTER TABLE member_rating_results ADD spouse_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'child_reinsurance_premium')
    ALTER TABLE member_rating_results ADD child_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'parent_reinsurance_premium')
    ALTER TABLE member_rating_results ADD parent_reinsurance_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'dependant_reinsurance_premium')
    ALTER TABLE member_rating_results ADD dependant_reinsurance_premium FLOAT;
