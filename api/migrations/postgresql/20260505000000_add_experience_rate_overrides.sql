-- Migration: per-(quote, scheme_category, benefit) experience-rate overrides.
-- Used when the quote's experience_rating field is set to 'Override': the
-- actuary supplies a mode (theoretical | experience_rated) and an
-- override_rate per (category, benefit) so the calculation engine bypasses
-- claims-based experience adjustment and uses the override rate directly.

CREATE TABLE IF NOT EXISTS group_pricing_experience_rate_overrides (
    id              SERIAL PRIMARY KEY,
    quote_id        INTEGER NOT NULL,
    scheme_category VARCHAR(128) NOT NULL,
    benefit         VARCHAR(8) NOT NULL,
    mode            VARCHAR(24) NOT NULL DEFAULT 'theoretical',
    override_rate   DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      VARCHAR(128) NOT NULL DEFAULT '',
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by      VARCHAR(128) NOT NULL DEFAULT '',
    CONSTRAINT ux_group_pricing_exp_overrides UNIQUE (quote_id, scheme_category, benefit)
);

CREATE INDEX IF NOT EXISTS idx_gpero_quote ON group_pricing_experience_rate_overrides(quote_id);
