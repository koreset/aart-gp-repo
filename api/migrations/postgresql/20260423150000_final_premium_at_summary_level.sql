-- Migration: collapse per-member Final* premiums and move Final premiums to
-- the summary only. On `member_rating_result_summaries`:
--   1) rename `final_adj_total_<x>_annual_office_premium` -> `final_<x>_office_premium`
--      (preserves data if old column exists)
--   2) ensure every Final*/ProportionFinal* column the model declares exists,
--      adding it with DOUBLE PRECISION DEFAULT 0 when missing.
-- On `member_rating_results`, drop the old per-member Final* office-premium
-- columns — Final now lives on the summary.

DO $$
BEGIN
    -- -----------------------------------------------------------------------
    -- Phase 1 — RENAMES on member_rating_result_summaries
    -- -----------------------------------------------------------------------
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_adj_total_gla_annual_office_premium')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_office_premium') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN final_adj_total_gla_annual_office_premium TO final_gla_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_adj_total_tax_saver_annual_office_premium')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_tax_saver_office_premium') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN final_adj_total_tax_saver_annual_office_premium TO final_tax_saver_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_adj_total_additional_accidental_gla_annual_office_premium')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_additional_accidental_gla_office_premium') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN final_adj_total_additional_accidental_gla_annual_office_premium TO final_additional_accidental_gla_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_adj_total_ptd_annual_office_premium')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_office_premium') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN final_adj_total_ptd_annual_office_premium TO final_ptd_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_adj_total_ci_annual_office_premium')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ci_office_premium') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN final_adj_total_ci_annual_office_premium TO final_ci_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_adj_total_sgla_annual_office_premium')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_sgla_office_premium') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN final_adj_total_sgla_annual_office_premium TO final_sgla_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_adj_total_ttd_annual_office_premium')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ttd_office_premium') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN final_adj_total_ttd_annual_office_premium TO final_ttd_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_adj_total_phi_annual_office_premium')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_phi_office_premium') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN final_adj_total_phi_annual_office_premium TO final_phi_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_adj_total_fun_annual_office_premium')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_fun_office_premium') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN final_adj_total_fun_annual_office_premium TO final_fun_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_adj_total_gla_educator_office_premium')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_educator_office_premium') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN final_adj_total_gla_educator_office_premium TO final_gla_educator_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_adj_total_ptd_educator_office_premium')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_educator_office_premium') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN final_adj_total_ptd_educator_office_premium TO final_ptd_educator_office_premium;
    END IF;

    -- -----------------------------------------------------------------------
    -- Phase 2 — ADD every Final*/ProportionFinal* column the model declares.
    -- Idempotent: each statement is guarded by NOT EXISTS. The office_premium
    -- columns are re-checked so that when Phase-1 rename was a no-op (legacy
    -- name not present) the column is still created.
    -- -----------------------------------------------------------------------

    -- Base benefit: office_premium (fallback ADD for each rename target)
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_office_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_office_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_tax_saver_office_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_tax_saver_office_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_additional_accidental_gla_office_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_additional_accidental_gla_office_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_office_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_office_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ci_office_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_office_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_sgla_office_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_office_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ttd_office_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_office_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_phi_office_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_office_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_fun_office_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_office_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_educator_office_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_educator_office_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_educator_office_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_educator_office_premium DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- Base benefit: risk_premium
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_risk_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_risk_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_tax_saver_risk_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_tax_saver_risk_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_additional_accidental_gla_risk_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_additional_accidental_gla_risk_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_risk_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_risk_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ci_risk_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_risk_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_sgla_risk_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_risk_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ttd_risk_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_risk_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_phi_risk_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_risk_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_fun_risk_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_risk_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_educator_risk_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_educator_risk_premium DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_educator_risk_premium') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_educator_risk_premium DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- Base benefit: proportion office_premium_salary (AAGLA uses short gorm name)
    -- proportion_final_gla_office_premium_salary is now computed inline in
    -- quote_template/schema.go and is no longer persisted.
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_tax_saver_office_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_tax_saver_office_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_add_acc_gla_office_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_add_acc_gla_office_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_ptd_office_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ptd_office_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_ci_office_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ci_office_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_sgla_office_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_sgla_office_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_ttd_office_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ttd_office_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_phi_office_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_phi_office_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_fun_office_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_fun_office_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_gla_educator_office_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_gla_educator_office_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_ptd_educator_office_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ptd_educator_office_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- Base benefit: proportion annual_risk_premium_salary (AAGLA uses short gorm name)
    -- proportion_final_gla_annual_risk_premium_salary is now computed inline in
    -- quote_template/schema.go and is no longer persisted.
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_tax_saver_annual_risk_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_tax_saver_annual_risk_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_add_acc_gla_ann_risk_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_add_acc_gla_ann_risk_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_ptd_annual_risk_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ptd_annual_risk_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_ci_annual_risk_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ci_annual_risk_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_sgla_annual_risk_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_sgla_annual_risk_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_ttd_annual_risk_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ttd_annual_risk_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_phi_annual_risk_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_phi_annual_risk_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_fun_annual_risk_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_fun_annual_risk_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_gla_educator_annual_risk_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_gla_educator_annual_risk_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='proportion_final_ptd_educator_annual_risk_premium_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ptd_educator_annual_risk_premium_salary DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- Base benefit: office_rate_per_1000_sa
    -- final_gla_office_rate_per1000_sa is now computed inline in
    -- quote_template/schema.go and is no longer persisted.
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_tax_saver_office_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_tax_saver_office_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_additional_accidental_gla_office_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_additional_accidental_gla_office_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_office_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_office_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ci_office_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_office_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_sgla_office_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_office_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ttd_office_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_office_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_phi_office_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_office_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_fun_office_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_office_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_educator_office_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_educator_office_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_educator_office_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_educator_office_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- Base benefit: risk_rate_per_1000_sa
    -- final_gla_risk_rate_per1000_sa is now computed inline in
    -- quote_template/schema.go and is no longer persisted.
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_tax_saver_risk_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_tax_saver_risk_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_additional_accidental_gla_risk_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_additional_accidental_gla_risk_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_risk_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_risk_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ci_risk_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_risk_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_sgla_risk_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_risk_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ttd_risk_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_risk_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_phi_risk_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_risk_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_fun_risk_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_risk_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_educator_risk_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_educator_risk_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_educator_risk_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_educator_risk_rate_per1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- -----------------------------------------------------------------------
    -- Phase 2b — Conversion / continuity slice columns (gorm-short names).
    -- 14 (benefit × slice) pairings × 6 columns each.
    -- -----------------------------------------------------------------------

    -- GLA conv_on_wdr
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_conv_on_wdr_risk_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_wdr_risk_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_conv_on_wdr_office_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_wdr_office_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_gla_conv_on_wdr_risk_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_conv_on_wdr_risk_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_gla_conv_on_wdr_office_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_conv_on_wdr_office_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_conv_on_wdr_risk_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_wdr_risk_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_conv_on_wdr_office_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_wdr_office_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- GLA conv_on_ret
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_conv_on_ret_risk_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_ret_risk_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_conv_on_ret_office_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_ret_office_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_gla_conv_on_ret_risk_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_conv_on_ret_risk_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_gla_conv_on_ret_office_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_conv_on_ret_office_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_conv_on_ret_risk_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_ret_risk_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_conv_on_ret_office_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_ret_office_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- GLA cont_dur_dis
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_cont_dur_dis_risk_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_cont_dur_dis_risk_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_cont_dur_dis_office_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_cont_dur_dis_office_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_gla_cont_dur_dis_risk_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_cont_dur_dis_risk_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_gla_cont_dur_dis_office_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_cont_dur_dis_office_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_cont_dur_dis_risk_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_cont_dur_dis_risk_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_cont_dur_dis_office_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_cont_dur_dis_office_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- GLA Educator conv_on_wdr
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_ed_conv_on_wdr_risk_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_wdr_risk_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_ed_conv_on_wdr_office_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_wdr_office_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_gla_ed_conv_on_wdr_risk_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_ed_conv_on_wdr_risk_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_gla_ed_conv_on_wdr_office_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_ed_conv_on_wdr_office_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_ed_conv_on_wdr_risk_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_wdr_risk_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_ed_conv_on_wdr_office_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_wdr_office_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- GLA Educator conv_on_ret
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_ed_conv_on_ret_risk_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_ret_risk_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_ed_conv_on_ret_office_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_ret_office_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_gla_ed_conv_on_ret_risk_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_ed_conv_on_ret_risk_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_gla_ed_conv_on_ret_office_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_ed_conv_on_ret_office_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_ed_conv_on_ret_risk_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_ret_risk_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_ed_conv_on_ret_office_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_ret_office_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- GLA Educator cont_dur_dis
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_ed_cont_dur_dis_risk_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_cont_dur_dis_risk_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_ed_cont_dur_dis_office_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_cont_dur_dis_office_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_gla_ed_cont_dur_dis_risk_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_ed_cont_dur_dis_risk_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_gla_ed_cont_dur_dis_office_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_ed_cont_dur_dis_office_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_ed_cont_dur_dis_risk_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_cont_dur_dis_risk_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_gla_ed_cont_dur_dis_office_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_cont_dur_dis_office_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- PTD conv_on_wdr
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_conv_on_wdr_risk_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_conv_on_wdr_risk_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_conv_on_wdr_office_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_conv_on_wdr_office_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_ptd_conv_on_wdr_risk_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ptd_conv_on_wdr_risk_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_ptd_conv_on_wdr_office_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ptd_conv_on_wdr_office_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_conv_on_wdr_risk_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_conv_on_wdr_risk_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_conv_on_wdr_office_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_conv_on_wdr_office_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- PTD Educator conv_on_wdr
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_ed_conv_on_wdr_risk_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_wdr_risk_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_ed_conv_on_wdr_office_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_wdr_office_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_ptd_ed_conv_on_wdr_risk_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ptd_ed_conv_on_wdr_risk_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_ptd_ed_conv_on_wdr_office_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ptd_ed_conv_on_wdr_office_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_ed_conv_on_wdr_risk_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_wdr_risk_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_ed_conv_on_wdr_office_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_wdr_office_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- PTD Educator conv_on_ret
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_ed_conv_on_ret_risk_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_ret_risk_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_ed_conv_on_ret_office_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_ret_office_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_ptd_ed_conv_on_ret_risk_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ptd_ed_conv_on_ret_risk_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_ptd_ed_conv_on_ret_office_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ptd_ed_conv_on_ret_office_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_ed_conv_on_ret_risk_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_ret_risk_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ptd_ed_conv_on_ret_office_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_ret_office_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- PHI conv_on_wdr
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_phi_conv_on_wdr_risk_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_conv_on_wdr_risk_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_phi_conv_on_wdr_office_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_conv_on_wdr_office_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_phi_conv_on_wdr_risk_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_phi_conv_on_wdr_risk_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_phi_conv_on_wdr_office_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_phi_conv_on_wdr_office_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_phi_conv_on_wdr_risk_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_conv_on_wdr_risk_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_phi_conv_on_wdr_office_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_conv_on_wdr_office_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- TTD conv_on_wdr
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ttd_conv_on_wdr_risk_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_conv_on_wdr_risk_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ttd_conv_on_wdr_office_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_conv_on_wdr_office_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_ttd_conv_on_wdr_risk_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ttd_conv_on_wdr_risk_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_ttd_conv_on_wdr_office_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ttd_conv_on_wdr_office_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ttd_conv_on_wdr_risk_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_conv_on_wdr_risk_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ttd_conv_on_wdr_office_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_conv_on_wdr_office_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- CI conv_on_wdr
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ci_conv_on_wdr_risk_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_conv_on_wdr_risk_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ci_conv_on_wdr_office_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_conv_on_wdr_office_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_ci_conv_on_wdr_risk_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ci_conv_on_wdr_risk_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_ci_conv_on_wdr_office_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ci_conv_on_wdr_office_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ci_conv_on_wdr_risk_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_conv_on_wdr_risk_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_ci_conv_on_wdr_office_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_conv_on_wdr_office_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- SGLA conv_on_wdr
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_sgla_conv_on_wdr_risk_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_conv_on_wdr_risk_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_sgla_conv_on_wdr_office_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_conv_on_wdr_office_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_sgla_conv_on_wdr_risk_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_sgla_conv_on_wdr_risk_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_sgla_conv_on_wdr_office_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_sgla_conv_on_wdr_office_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_sgla_conv_on_wdr_risk_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_conv_on_wdr_risk_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_sgla_conv_on_wdr_office_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_conv_on_wdr_office_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- FUN conv_on_wdr
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_fun_conv_on_wdr_risk_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_conv_on_wdr_risk_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_fun_conv_on_wdr_office_prem') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_conv_on_wdr_office_prem DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_fun_conv_on_wdr_risk_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_fun_conv_on_wdr_risk_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='prop_final_fun_conv_on_wdr_office_prem_salary') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_fun_conv_on_wdr_office_prem_salary DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_fun_conv_on_wdr_risk_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_conv_on_wdr_risk_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;
    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='final_fun_conv_on_wdr_office_rate_per_1000_sa') THEN
        ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_conv_on_wdr_office_rate_per_1000_sa DOUBLE PRECISION DEFAULT 0;
    END IF;

    -- -----------------------------------------------------------------------
    -- Phase 3 — DROPS on member_rating_results (Final now lives on the summary)
    -- -----------------------------------------------------------------------
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='final_gla_office_premium') THEN
        ALTER TABLE member_rating_results DROP COLUMN final_gla_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='final_tax_saver_office_premium') THEN
        ALTER TABLE member_rating_results DROP COLUMN final_tax_saver_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='final_additional_accidental_gla_office_premium') THEN
        ALTER TABLE member_rating_results DROP COLUMN final_additional_accidental_gla_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='final_ptd_office_premium') THEN
        ALTER TABLE member_rating_results DROP COLUMN final_ptd_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='final_ci_office_premium') THEN
        ALTER TABLE member_rating_results DROP COLUMN final_ci_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='final_spouse_gla_office_premium') THEN
        ALTER TABLE member_rating_results DROP COLUMN final_spouse_gla_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='final_ttd_office_premium') THEN
        ALTER TABLE member_rating_results DROP COLUMN final_ttd_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='final_phi_office_premium') THEN
        ALTER TABLE member_rating_results DROP COLUMN final_phi_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='final_total_funeral_office_cost') THEN
        ALTER TABLE member_rating_results DROP COLUMN final_total_funeral_office_cost;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='final_gla_educator_office_premium') THEN
        ALTER TABLE member_rating_results DROP COLUMN final_gla_educator_office_premium;
    END IF;
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='final_ptd_educator_office_premium') THEN
        ALTER TABLE member_rating_results DROP COLUMN final_ptd_educator_office_premium;
    END IF;
END $$;
