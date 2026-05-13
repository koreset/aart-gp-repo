-- Add an indexed, persisted computed column `lookup_key` on phi_rates and
-- reinsurance_phi_rates so the per-member rate fetches in GetPhiRate /
-- GetReinsurancePhiRate become single-column index seeks instead of
-- multi-column full scans on a table with no covering index.
--
-- The concatenation order and `|` separator MUST match the lookupKey
-- string built in api/services/group_pricing.go. MSSQL's CONCAT treats
-- NULL as an empty string, so no explicit COALESCE is needed.
-- Idempotent on re-runs.

-- 1) phi_rates.lookup_key
IF NOT EXISTS (SELECT 1 FROM sys.columns
               WHERE object_id = OBJECT_ID('phi_rates')
                 AND name = 'lookup_key')
BEGIN
    ALTER TABLE phi_rates
        ADD lookup_key AS (
            CONCAT(
                risk_rate_code, '|',
                risk_type, '|',
                CAST(age_next_birthday AS VARCHAR(10)), '|',
                gender, '|',
                CAST(occupation_class AS VARCHAR(10)), '|',
                CAST(income_level AS VARCHAR(10)), '|',
                CAST(waiting_period AS VARCHAR(10)), '|',
                CAST(deferred_period AS VARCHAR(10)), '|',
                CAST(normal_retirement_age AS VARCHAR(10)), '|',
                benefit_escalation_option, '|',
                disability_definition
            )
        ) PERSISTED;
END;

-- 2) reinsurance_phi_rates.lookup_key
IF NOT EXISTS (SELECT 1 FROM sys.columns
               WHERE object_id = OBJECT_ID('reinsurance_phi_rates')
                 AND name = 'lookup_key')
BEGIN
    ALTER TABLE reinsurance_phi_rates
        ADD lookup_key AS (
            CONCAT(
                risk_rate_code, '|',
                CAST(age_next_birthday AS VARCHAR(10)), '|',
                CAST(income_level AS VARCHAR(10)), '|',
                gender, '|',
                CAST(occupation_class AS VARCHAR(10))
            )
        ) PERSISTED;
END;

-- 3) Indexes
IF NOT EXISTS (SELECT 1 FROM sys.indexes
               WHERE name = 'idx_phi_rates_lookup_key'
                 AND object_id = OBJECT_ID('phi_rates'))
    CREATE INDEX idx_phi_rates_lookup_key
        ON phi_rates (lookup_key);

IF NOT EXISTS (SELECT 1 FROM sys.indexes
               WHERE name = 'idx_reins_phi_rates_lookup_key'
                 AND object_id = OBJECT_ID('reinsurance_phi_rates'))
    CREATE INDEX idx_reins_phi_rates_lookup_key
        ON reinsurance_phi_rates (lookup_key);
