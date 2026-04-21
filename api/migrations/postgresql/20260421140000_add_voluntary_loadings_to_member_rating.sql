-- Migration: add per-benefit voluntary loading columns to
-- member_rating_results. The voluntary-loading Go struct fields were
-- introduced alongside the voluntary-loading rate tables but the matching
-- ALTER for member_rating_results was never shipped, so GORM INSERTs fail
-- on any DB provisioned before this fix.

CREATE TABLE IF NOT EXISTS member_rating_results (
    id SERIAL PRIMARY KEY
);

ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_voluntary_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_voluntary_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_voluntary_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_voluntary_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_voluntary_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS fun_voluntary_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_gla_voluntary_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_ptd_voluntary_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_ci_voluntary_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_ttd_voluntary_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_phi_voluntary_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_fun_voluntary_loading NUMERIC(20,6);
