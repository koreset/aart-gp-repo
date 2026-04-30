-- Per-row credibility on each experience-rate override entry, plus
-- per-benefit credibility columns on the historical credibility audit so the
-- weight assumed for each benefit on a calc run is preserved.

IF NOT EXISTS (
    SELECT 1 FROM sys.columns
    WHERE object_id = OBJECT_ID('group_pricing_experience_rate_overrides')
      AND name = 'credibility'
)
BEGIN
    ALTER TABLE group_pricing_experience_rate_overrides
        ADD credibility FLOAT NOT NULL CONSTRAINT df_gpero_credibility DEFAULT 0;
END;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('historical_credibility_data') AND name = 'gla_credibility')
BEGIN
    ALTER TABLE historical_credibility_data
        ADD gla_credibility FLOAT NOT NULL CONSTRAINT df_hcd_gla_cred DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('historical_credibility_data') AND name = 'aagla_credibility')
BEGIN
    ALTER TABLE historical_credibility_data
        ADD aagla_credibility FLOAT NOT NULL CONSTRAINT df_hcd_aagla_cred DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('historical_credibility_data') AND name = 'sgla_credibility')
BEGIN
    ALTER TABLE historical_credibility_data
        ADD sgla_credibility FLOAT NOT NULL CONSTRAINT df_hcd_sgla_cred DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('historical_credibility_data') AND name = 'ptd_credibility')
BEGIN
    ALTER TABLE historical_credibility_data
        ADD ptd_credibility FLOAT NOT NULL CONSTRAINT df_hcd_ptd_cred DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('historical_credibility_data') AND name = 'ttd_credibility')
BEGIN
    ALTER TABLE historical_credibility_data
        ADD ttd_credibility FLOAT NOT NULL CONSTRAINT df_hcd_ttd_cred DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('historical_credibility_data') AND name = 'phi_credibility')
BEGIN
    ALTER TABLE historical_credibility_data
        ADD phi_credibility FLOAT NOT NULL CONSTRAINT df_hcd_phi_cred DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('historical_credibility_data') AND name = 'ci_credibility')
BEGIN
    ALTER TABLE historical_credibility_data
        ADD ci_credibility FLOAT NOT NULL CONSTRAINT df_hcd_ci_cred DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('historical_credibility_data') AND name = 'fun_credibility')
BEGIN
    ALTER TABLE historical_credibility_data
        ADD fun_credibility FLOAT NOT NULL CONSTRAINT df_hcd_fun_cred DEFAULT 0;
END;
