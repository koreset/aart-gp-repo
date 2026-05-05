-- Migration: persist per-benefit total annual salary on
-- member_rating_result_summaries. These exclude members whose age has passed
-- the matching restriction.<benefit>_max_cover_age and serve as the
-- denominator for Proportion*RiskPremiumSalary fields, so an aged-out member
-- contributes 0 to both numerator (no premium) and denominator (no salary).
-- TotalAnnualSalary remains the unfiltered gross sum.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_annual_salary_gla')
BEGIN
    ALTER TABLE member_rating_result_summaries
        ADD total_annual_salary_gla  FLOAT NOT NULL CONSTRAINT df_mrrs_total_annual_salary_gla  DEFAULT 0,
            total_annual_salary_ptd  FLOAT NOT NULL CONSTRAINT df_mrrs_total_annual_salary_ptd  DEFAULT 0,
            total_annual_salary_ci   FLOAT NOT NULL CONSTRAINT df_mrrs_total_annual_salary_ci   DEFAULT 0,
            total_annual_salary_sgla FLOAT NOT NULL CONSTRAINT df_mrrs_total_annual_salary_sgla DEFAULT 0,
            total_annual_salary_ttd  FLOAT NOT NULL CONSTRAINT df_mrrs_total_annual_salary_ttd  DEFAULT 0,
            total_annual_salary_phi  FLOAT NOT NULL CONSTRAINT df_mrrs_total_annual_salary_phi  DEFAULT 0,
            total_annual_salary_fun  FLOAT NOT NULL CONSTRAINT df_mrrs_total_annual_salary_fun  DEFAULT 0;
END;
