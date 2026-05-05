-- Add per-benefit scheme-size loading columns to scheme_size_levels.
-- These flow into each benefit's LoadedRate multiplier in
-- PopulateRatesPerMember (GLA/PTD/CI/PHI/TTD) and into the
-- TotalFuneralRiskPremium multiplier (Fun).

ALTER TABLE scheme_size_levels
    ADD COLUMN IF NOT EXISTS gla_loading     DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ptd_loading     DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ci_loading      DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS phi_loading     DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ttd_loading     DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS funeral_loading DOUBLE PRECISION NOT NULL DEFAULT 0;
