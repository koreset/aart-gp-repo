-- Add educator base loading columns to general_loadings. Folded into the
-- educator multiplier in computeEducatorLoadedRates alongside the educator
-- conversion / continuity slice loadings (per risk_rate_code, age, gender).

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('general_loadings') AND name = 'gla_educator_loading_rate')
BEGIN
    ALTER TABLE general_loadings
        ADD gla_educator_loading_rate FLOAT NOT NULL CONSTRAINT df_gl_gla_educator_loading_rate DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('general_loadings') AND name = 'ptd_educator_loading_rate')
BEGIN
    ALTER TABLE general_loadings
        ADD ptd_educator_loading_rate FLOAT NOT NULL CONSTRAINT df_gl_ptd_educator_loading_rate DEFAULT 0;
END;
