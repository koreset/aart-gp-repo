-- Follow-up to 20260519010000_add_salary_continuation_benefit_tables.sql:
-- defensively re-adds salary_continuation_rates.lookup_key for installs
-- where the original ALTER didn't take effect. Idempotent via
-- ADD COLUMN IF NOT EXISTS / CREATE INDEX IF NOT EXISTS.
--
-- The concatenation order and `|` separator MUST match the lookupKey
-- string built by GetScbRate in api/services/group_pricing.go.

ALTER TABLE salary_continuation_rates
    ADD COLUMN IF NOT EXISTS lookup_key TEXT
    GENERATED ALWAYS AS (
        COALESCE(risk_rate_code, '') || '|' ||
        COALESCE(risk_type, '') || '|' ||
        age_next_birthday::TEXT || '|' ||
        COALESCE(gender, '') || '|' ||
        occupation_class::TEXT || '|' ||
        income_level::TEXT || '|' ||
        deferred_period::TEXT || '|' ||
        excess_period::TEXT || '|' ||
        COALESCE(benefit_escalation_option, '') || '|' ||
        COALESCE(disability_definition, '')
    ) STORED;

CREATE INDEX IF NOT EXISTS idx_scb_rates_lookup_key
    ON salary_continuation_rates (lookup_key);
