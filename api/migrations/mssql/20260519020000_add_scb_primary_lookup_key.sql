-- Follow-up to 20260519010000_add_salary_continuation_benefit_tables.sql:
-- defensively re-adds salary_continuation_rates.lookup_key for installs
-- where the original ALTER didn't take effect. Idempotent via
-- IF NOT EXISTS guards on sys.columns and sys.indexes.
--
-- The concatenation order and `|` separator MUST match the lookupKey
-- string built by GetScbRate in api/services/group_pricing.go. MSSQL's
-- CONCAT treats NULL as an empty string, so no explicit COALESCE needed.

IF NOT EXISTS (SELECT 1 FROM sys.columns
               WHERE object_id = OBJECT_ID('salary_continuation_rates')
                 AND name = 'lookup_key')
BEGIN
    ALTER TABLE salary_continuation_rates
        ADD lookup_key AS (
            CONCAT(
                risk_rate_code, '|',
                risk_type, '|',
                CAST(age_next_birthday AS VARCHAR(20)), '|',
                gender, '|',
                CAST(occupation_class AS VARCHAR(20)), '|',
                CAST(income_level AS VARCHAR(20)), '|',
                CAST(deferred_period AS VARCHAR(20)), '|',
                CAST(excess_period AS VARCHAR(20)), '|',
                benefit_escalation_option, '|',
                disability_definition
            )
        ) PERSISTED;
END;

IF NOT EXISTS (SELECT 1 FROM sys.indexes
               WHERE name = 'idx_scb_rates_lookup_key'
                 AND object_id = OBJECT_ID('salary_continuation_rates'))
    CREATE INDEX idx_scb_rates_lookup_key
        ON salary_continuation_rates (lookup_key);
