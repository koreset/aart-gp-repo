-- Add per-benefit scheme-size loading columns to scheme_size_levels.
-- These flow into each benefit's LoadedRate multiplier in
-- PopulateRatesPerMember (GLA/PTD/CI/PHI/TTD) and into the
-- TotalFuneralRiskPremium multiplier (Fun).

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_size_levels') AND name = 'gla_loading')
BEGIN
    ALTER TABLE scheme_size_levels
        ADD gla_loading FLOAT NOT NULL CONSTRAINT df_ssl_gla_loading DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_size_levels') AND name = 'ptd_loading')
BEGIN
    ALTER TABLE scheme_size_levels
        ADD ptd_loading FLOAT NOT NULL CONSTRAINT df_ssl_ptd_loading DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_size_levels') AND name = 'ci_loading')
BEGIN
    ALTER TABLE scheme_size_levels
        ADD ci_loading FLOAT NOT NULL CONSTRAINT df_ssl_ci_loading DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_size_levels') AND name = 'phi_loading')
BEGIN
    ALTER TABLE scheme_size_levels
        ADD phi_loading FLOAT NOT NULL CONSTRAINT df_ssl_phi_loading DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_size_levels') AND name = 'ttd_loading')
BEGIN
    ALTER TABLE scheme_size_levels
        ADD ttd_loading FLOAT NOT NULL CONSTRAINT df_ssl_ttd_loading DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_size_levels') AND name = 'funeral_loading')
BEGIN
    ALTER TABLE scheme_size_levels
        ADD funeral_loading FLOAT NOT NULL CONSTRAINT df_ssl_funeral_loading DEFAULT 0;
END;
