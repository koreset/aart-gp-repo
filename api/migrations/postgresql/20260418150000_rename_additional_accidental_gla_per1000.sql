-- Migration: rename four Additional Accidental GLA rate-per-1000 columns on
-- member_rating_result_summaries to match GORM's NamingStrategy output
-- (`per1000_sa`, no underscore between `per` and `1000`).

ALTER TABLE member_rating_result_summaries
    RENAME COLUMN additional_accidental_gla_risk_rate_per_1000_sa
        TO additional_accidental_gla_risk_rate_per1000_sa;

ALTER TABLE member_rating_result_summaries
    RENAME COLUMN additional_accidental_gla_office_rate_per_1000_sa
        TO additional_accidental_gla_office_rate_per1000_sa;

ALTER TABLE member_rating_result_summaries
    RENAME COLUMN exp_additional_accidental_gla_risk_rate_per_1000_sa
        TO exp_additional_accidental_gla_risk_rate_per1000_sa;

ALTER TABLE member_rating_result_summaries
    RENAME COLUMN exp_additional_accidental_gla_office_rate_per_1000_sa
        TO exp_additional_accidental_gla_office_rate_per1000_sa;
