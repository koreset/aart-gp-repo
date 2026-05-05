-- Add educator base loading columns to general_loadings. Folded into the
-- educator multiplier in computeEducatorLoadedRates alongside the educator
-- conversion / continuity slice loadings (per risk_rate_code, age, gender).

ALTER TABLE general_loadings
    ADD COLUMN IF NOT EXISTS gla_educator_loading_rate DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ptd_educator_loading_rate DOUBLE PRECISION NOT NULL DEFAULT 0;
