-- Migration: rename educator GLA/PTD rate-per-1000 columns on
-- member_rating_result_summaries so they match GORM's NamingStrategy output
-- (`per1000_sa`, no underscore between `per` and `1000`). The
-- 20260421170000 migration used `per_1000_sa` which GORM cannot find on
-- INSERT/UPDATE — causing column-not-found errors during
-- CalculateGroupPricingQuote. Follows the precedent of
-- 20260418150000_rename_additional_accidental_gla_per1000.sql.

DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='gla_educator_risk_rate_per_1000_sa')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='gla_educator_risk_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN gla_educator_risk_rate_per_1000_sa TO gla_educator_risk_rate_per1000_sa;
    END IF;

    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='gla_educator_office_rate_per_1000_sa')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='gla_educator_office_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN gla_educator_office_rate_per_1000_sa TO gla_educator_office_rate_per1000_sa;
    END IF;

    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='exp_gla_educator_risk_rate_per_1000_sa')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='exp_gla_educator_risk_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN exp_gla_educator_risk_rate_per_1000_sa TO exp_gla_educator_risk_rate_per1000_sa;
    END IF;

    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='exp_gla_educator_office_rate_per_1000_sa')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='exp_gla_educator_office_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN exp_gla_educator_office_rate_per_1000_sa TO exp_gla_educator_office_rate_per1000_sa;
    END IF;

    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='ptd_educator_risk_rate_per_1000_sa')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='ptd_educator_risk_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN ptd_educator_risk_rate_per_1000_sa TO ptd_educator_risk_rate_per1000_sa;
    END IF;

    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='ptd_educator_office_rate_per_1000_sa')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='ptd_educator_office_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN ptd_educator_office_rate_per_1000_sa TO ptd_educator_office_rate_per1000_sa;
    END IF;

    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='exp_ptd_educator_risk_rate_per_1000_sa')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='exp_ptd_educator_risk_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN exp_ptd_educator_risk_rate_per_1000_sa TO exp_ptd_educator_risk_rate_per1000_sa;
    END IF;

    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='exp_ptd_educator_office_rate_per_1000_sa')
        AND NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_result_summaries' AND column_name='exp_ptd_educator_office_rate_per1000_sa') THEN
        ALTER TABLE member_rating_result_summaries RENAME COLUMN exp_ptd_educator_office_rate_per_1000_sa TO exp_ptd_educator_office_rate_per1000_sa;
    END IF;
END $$;
