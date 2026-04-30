-- Per-row credibility on each experience-rate override entry, plus
-- per-benefit credibility columns on the historical credibility audit so the
-- weight assumed for each benefit on a calc run is preserved.

ALTER TABLE group_pricing_experience_rate_overrides
    ADD COLUMN credibility DOUBLE NOT NULL DEFAULT 0;

ALTER TABLE historical_credibility_data
    ADD COLUMN gla_credibility   DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN aagla_credibility DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN sgla_credibility  DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN ptd_credibility   DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN ttd_credibility   DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN phi_credibility   DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN ci_credibility    DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN fun_credibility   DOUBLE NOT NULL DEFAULT 0;
