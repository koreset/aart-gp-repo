-- Migration: persist per-benefit total annual salary on
-- member_rating_result_summaries. These exclude members whose age has passed
-- the matching restriction.<benefit>_max_cover_age and serve as the
-- denominator for Proportion*RiskPremiumSalary fields, so an aged-out member
-- contributes 0 to both numerator (no premium) and denominator (no salary).
-- TotalAnnualSalary remains the unfiltered gross sum.

ALTER TABLE member_rating_result_summaries
    ADD COLUMN IF NOT EXISTS total_annual_salary_gla  DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS total_annual_salary_ptd  DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS total_annual_salary_ci   DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS total_annual_salary_sgla DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS total_annual_salary_ttd  DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS total_annual_salary_phi  DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS total_annual_salary_fun  DOUBLE PRECISION NOT NULL DEFAULT 0;
