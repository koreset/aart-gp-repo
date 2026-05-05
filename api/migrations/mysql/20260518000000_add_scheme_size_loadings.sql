-- Add per-benefit scheme-size loading columns to scheme_size_levels.
-- These flow into each benefit's LoadedRate multiplier in
-- PopulateRatesPerMember (GLA/PTD/CI/PHI/TTD) and into the
-- TotalFuneralRiskPremium multiplier (Fun).

ALTER TABLE scheme_size_levels
    ADD COLUMN gla_loading     DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN ptd_loading     DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN ci_loading      DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN phi_loading     DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN ttd_loading     DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN funeral_loading DOUBLE NOT NULL DEFAULT 0;
