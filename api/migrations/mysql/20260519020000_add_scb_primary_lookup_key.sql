-- Follow-up to 20260519010000_add_salary_continuation_benefit_tables.sql:
-- on some MySQL installs the PREPARE/EXECUTE block that adds the
-- lookup_key column to salary_continuation_rates did not actually take
-- effect (the reinsurance variant in the same file succeeded). This
-- migration adds the column + index defensively on its own.
--
-- The concatenation order and `|` separator MUST match the lookupKey
-- string built by GetScbRate in api/services/group_pricing.go.
-- Idempotent on re-runs.

-- 1) salary_continuation_rates.lookup_key
SET @has_col := (SELECT COUNT(*) FROM information_schema.columns
                 WHERE table_schema = DATABASE()
                   AND table_name = 'salary_continuation_rates'
                   AND column_name = 'lookup_key');
SET @sql := IF(@has_col = 0,
  'ALTER TABLE salary_continuation_rates ADD COLUMN lookup_key VARCHAR(512) GENERATED ALWAYS AS (CONCAT(COALESCE(risk_rate_code, ''''), ''|'', COALESCE(risk_type, ''''), ''|'', CAST(age_next_birthday AS CHAR), ''|'', COALESCE(gender, ''''), ''|'', CAST(occupation_class AS CHAR), ''|'', CAST(income_level AS CHAR), ''|'', CAST(deferred_period AS CHAR), ''|'', CAST(excess_period AS CHAR), ''|'', COALESCE(benefit_escalation_option, ''''), ''|'', COALESCE(disability_definition, ''''))) STORED',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 2) Index on salary_continuation_rates.lookup_key
SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE()
               AND table_name = 'salary_continuation_rates'
               AND index_name = 'idx_scb_rates_lookup_key');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_scb_rates_lookup_key ON salary_continuation_rates (lookup_key)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
