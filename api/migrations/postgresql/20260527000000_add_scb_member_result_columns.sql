-- Add Salary Continuation Benefit (SCB) per-member and summary columns.
-- Mirrors the field structure of the conversion-on-withdrawal slices on
-- MemberRatingResult / MemberRatingResultSummary; SCB is tracked as a
-- reportable slice (NOT added to any group total).
--
-- Source-side rate tables and scheme_categories columns were added in
-- 20260519010000_add_salary_continuation_benefit_tables.sql.
-- Idempotent on re-runs.

ALTER TABLE member_rating_results
    ADD COLUMN IF NOT EXISTS scb_rate                       DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS base_scb_rate                  DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS loaded_scb_rate                DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS exp_adj_loaded_scb_rate        DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS scb_risk_premium               DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS exp_adj_scb_risk_premium       DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reins_scb_rate                 DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS base_reins_scb_rate            DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS loaded_reins_scb_rate          DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reins_scb_risk_premium         DOUBLE PRECISION NOT NULL DEFAULT 0;

ALTER TABLE member_rating_result_summaries
    ADD COLUMN IF NOT EXISTS total_scb_annual_risk_premium             DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS exp_adj_total_scb_annual_risk_premium     DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS proportion_scb_risk_premium_salary        DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS exp_adj_proportion_scb_risk_premium_salary DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS scb_risk_rate_per_1000_income             DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS exp_scb_risk_rate_per_1000_income         DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS final_scb_office_premium                  DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS total_reins_scb_annual_risk_premium       DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS final_reins_scb_office_premium            DOUBLE PRECISION NOT NULL DEFAULT 0;
