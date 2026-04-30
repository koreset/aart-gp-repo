-- Per-row credibility on each experience-rate override entry, plus
-- per-benefit credibility columns on the historical credibility audit so the
-- weight assumed for each benefit on a calc run is preserved.

ALTER TABLE group_pricing_experience_rate_overrides
    ADD COLUMN IF NOT EXISTS credibility DOUBLE PRECISION NOT NULL DEFAULT 0;

ALTER TABLE historical_credibility_data
    ADD COLUMN IF NOT EXISTS gla_credibility   DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE historical_credibility_data
    ADD COLUMN IF NOT EXISTS aagla_credibility DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE historical_credibility_data
    ADD COLUMN IF NOT EXISTS sgla_credibility  DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE historical_credibility_data
    ADD COLUMN IF NOT EXISTS ptd_credibility   DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE historical_credibility_data
    ADD COLUMN IF NOT EXISTS ttd_credibility   DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE historical_credibility_data
    ADD COLUMN IF NOT EXISTS phi_credibility   DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE historical_credibility_data
    ADD COLUMN IF NOT EXISTS ci_credibility    DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE historical_credibility_data
    ADD COLUMN IF NOT EXISTS fun_credibility   DOUBLE PRECISION NOT NULL DEFAULT 0;
