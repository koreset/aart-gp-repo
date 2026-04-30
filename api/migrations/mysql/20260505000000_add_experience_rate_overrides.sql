-- Migration: per-(quote, scheme_category, benefit) experience-rate overrides.
-- Used when the quote's experience_rating field is set to 'Override': the
-- actuary supplies a mode (theoretical | experience_rated) and an
-- override_rate per (category, benefit) so the calculation engine bypasses
-- claims-based experience adjustment and uses the override rate directly.

CREATE TABLE IF NOT EXISTS group_pricing_experience_rate_overrides (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    quote_id        INT NOT NULL,
    scheme_category VARCHAR(128) NOT NULL,
    benefit         VARCHAR(8) NOT NULL,
    mode            VARCHAR(24) NOT NULL DEFAULT 'theoretical',
    override_rate   DOUBLE NOT NULL DEFAULT 0,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      VARCHAR(128) NOT NULL DEFAULT '',
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by      VARCHAR(128) NOT NULL DEFAULT '',
    UNIQUE KEY ux_group_pricing_exp_overrides (quote_id, scheme_category, benefit),
    KEY idx_gpero_quote (quote_id)
);
