-- Migration: add optional "Additional Accidental GLA" sub-benefit.
-- Re-uses the gla_rates / gla_aids_rates tables with a different benefit_type.
-- All GLA scheme parameters stay on the GLA fields — only the benefit_type
-- differs for the Additional Accidental layer.

--------------------------------------------------------------------------------
-- scheme_categories
--------------------------------------------------------------------------------

ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_accidental_gla_benefit BOOLEAN DEFAULT FALSE;
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='additional_accidental_gla_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN additional_accidental_gla_benefit TYPE BOOLEAN USING additional_accidental_gla_benefit::BOOLEAN;
    END IF;
END $$;

ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_accidental_gla_benefit_type VARCHAR(255);

--------------------------------------------------------------------------------
-- member_rating_results
--------------------------------------------------------------------------------

ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS additional_accidental_gla_sum_assured NUMERIC(15,5);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS additional_accidental_gla_capped_sum_assured NUMERIC(15,5);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS additional_accidental_gla_qx NUMERIC(15,5);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS additional_accidental_gla_aids_qx NUMERIC(15,5);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS base_additional_accidental_gla_rate NUMERIC(15,5);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS loaded_additional_accidental_gla_rate NUMERIC(15,5);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS additional_accidental_gla_weighted_experience_crude_rate NUMERIC(15,5);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS additional_accidental_gla_theoretical_rate NUMERIC(15,5);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS additional_accidental_gla_experience_adjustment NUMERIC(15,5);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_loaded_additional_accidental_gla_rate NUMERIC(15,5);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS additional_accidental_gla_risk_premium NUMERIC(15,5);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_additional_accidental_gla_risk_premium NUMERIC(15,5);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS additional_accidental_gla_office_premium NUMERIC(15,5);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_additional_accidental_gla_office_premium NUMERIC(15,5);
