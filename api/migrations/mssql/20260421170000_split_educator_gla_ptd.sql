-- Migration: split educator benefit tracking into GLA and PTD components.
-- The combined Educator* columns stay as the sum; the new columns let the
-- business attribute the educator premium between GLA-educator and
-- PTD-educator and expose buildup fields (premium, %salary, rate per 1000).

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_results')
BEGIN
    CREATE TABLE member_rating_results (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_educator_risk_premium')
    ALTER TABLE member_rating_results ADD gla_educator_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_educator_office_premium')
    ALTER TABLE member_rating_results ADD gla_educator_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_educator_risk_premium')
    ALTER TABLE member_rating_results ADD exp_adj_gla_educator_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_educator_office_premium')
    ALTER TABLE member_rating_results ADD exp_adj_gla_educator_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_educator_risk_premium')
    ALTER TABLE member_rating_results ADD ptd_educator_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_educator_office_premium')
    ALTER TABLE member_rating_results ADD ptd_educator_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ptd_educator_risk_premium')
    ALTER TABLE member_rating_results ADD exp_adj_ptd_educator_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ptd_educator_office_premium')
    ALTER TABLE member_rating_results ADD exp_adj_ptd_educator_office_premium FLOAT;

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_result_summaries')
BEGIN
    CREATE TABLE member_rating_result_summaries (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_educator_sum_assured')
    ALTER TABLE member_rating_result_summaries ADD total_educator_sum_assured FLOAT;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_educator_risk_premium')
    ALTER TABLE member_rating_result_summaries ADD total_gla_educator_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_educator_office_premium')
    ALTER TABLE member_rating_result_summaries ADD total_gla_educator_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_educator_risk_premium')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_educator_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_educator_office_premium')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_educator_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'proportion_gla_educator_risk_premium_salary')
    ALTER TABLE member_rating_result_summaries ADD proportion_gla_educator_risk_premium_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'proportion_gla_educator_office_premium_salary')
    ALTER TABLE member_rating_result_summaries ADD proportion_gla_educator_office_premium_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_proportion_gla_educator_risk_premium_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_proportion_gla_educator_risk_premium_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_proportion_gla_educator_office_premium_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_proportion_gla_educator_office_premium_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'gla_educator_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD gla_educator_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'gla_educator_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD gla_educator_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_gla_educator_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_gla_educator_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_gla_educator_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_gla_educator_office_rate_per_1000_sa FLOAT;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_educator_risk_premium')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_educator_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_educator_office_premium')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_educator_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_ptd_educator_risk_premium')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_ptd_educator_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_ptd_educator_office_premium')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_ptd_educator_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'proportion_ptd_educator_risk_premium_salary')
    ALTER TABLE member_rating_result_summaries ADD proportion_ptd_educator_risk_premium_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'proportion_ptd_educator_office_premium_salary')
    ALTER TABLE member_rating_result_summaries ADD proportion_ptd_educator_office_premium_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_proportion_ptd_educator_risk_premium_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_proportion_ptd_educator_risk_premium_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_proportion_ptd_educator_office_premium_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_proportion_ptd_educator_office_premium_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'ptd_educator_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD ptd_educator_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'ptd_educator_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD ptd_educator_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_ptd_educator_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_ptd_educator_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_ptd_educator_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_ptd_educator_office_rate_per_1000_sa FLOAT;

-- Split binder / outsource per-member columns
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_educator_binder_amount')
    ALTER TABLE member_rating_results ADD gla_educator_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_educator_outsourced_amount')
    ALTER TABLE member_rating_results ADD gla_educator_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_educator_binder_amount')
    ALTER TABLE member_rating_results ADD exp_adj_gla_educator_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_educator_outsourced_amount')
    ALTER TABLE member_rating_results ADD exp_adj_gla_educator_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_educator_binder_amount')
    ALTER TABLE member_rating_results ADD ptd_educator_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_educator_outsourced_amount')
    ALTER TABLE member_rating_results ADD ptd_educator_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ptd_educator_binder_amount')
    ALTER TABLE member_rating_results ADD exp_adj_ptd_educator_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ptd_educator_outsourced_amount')
    ALTER TABLE member_rating_results ADD exp_adj_ptd_educator_outsourced_amount FLOAT;

-- Split binder / outsource summary totals
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_educator_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD total_gla_educator_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_educator_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD total_gla_educator_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_educator_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_educator_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_educator_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_educator_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_educator_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_educator_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_educator_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_educator_outsourced_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_ptd_educator_binder_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_ptd_educator_binder_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_ptd_educator_outsourced_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_ptd_educator_outsourced_amount FLOAT;

-- Split commission summary totals
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_educator_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_educator_commission_amount FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_ptd_educator_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_ptd_educator_commission_amount FLOAT;
