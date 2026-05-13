-- Add an indexed, stored generated column `lookup_key` on phi_rates and
-- reinsurance_phi_rates so the per-member rate fetches in GetPhiRate /
-- GetReinsurancePhiRate become single-column index seeks instead of
-- multi-column full scans on a table with no covering index.
--
-- The concatenation order and `|` separator MUST match the lookupKey
-- string built in api/services/group_pricing.go.
-- Idempotent on re-runs.

-- 1) phi_rates.lookup_key
ALTER TABLE phi_rates
    ADD COLUMN IF NOT EXISTS lookup_key TEXT
    GENERATED ALWAYS AS (
        COALESCE(risk_rate_code, '') || '|' ||
        COALESCE(risk_type, '') || '|' ||
        age_next_birthday::TEXT || '|' ||
        COALESCE(gender, '') || '|' ||
        occupation_class::TEXT || '|' ||
        income_level::TEXT || '|' ||
        waiting_period::TEXT || '|' ||
        deferred_period::TEXT || '|' ||
        normal_retirement_age::TEXT || '|' ||
        COALESCE(benefit_escalation_option, '') || '|' ||
        COALESCE(disability_definition, '')
    ) STORED;

-- 2) reinsurance_phi_rates.lookup_key
ALTER TABLE reinsurance_phi_rates
    ADD COLUMN IF NOT EXISTS lookup_key TEXT
    GENERATED ALWAYS AS (
        COALESCE(risk_rate_code, '') || '|' ||
        age_next_birthday::TEXT || '|' ||
        income_level::TEXT || '|' ||
        COALESCE(gender, '') || '|' ||
        occupation_class::TEXT
    ) STORED;

-- 3) Indexes
CREATE INDEX IF NOT EXISTS idx_phi_rates_lookup_key
    ON phi_rates (lookup_key);

CREATE INDEX IF NOT EXISTS idx_reins_phi_rates_lookup_key
    ON reinsurance_phi_rates (lookup_key);
