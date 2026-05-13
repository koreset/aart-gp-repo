-- Add an indexed, stored generated column `lookup_key` on phi_rates and
-- reinsurance_phi_rates so the per-member rate fetches in GetPhiRate /
-- GetReinsurancePhiRate become single-column index seeks instead of
-- multi-column full scans on a table with no covering index.
--
-- The concatenation order and `|` separator MUST match the lookupKey
-- string built in api/services/group_pricing.go.
-- Idempotent on re-runs.

-- ─────────────────────────────────────────────────────────────────────────
-- 1) phi_rates.lookup_key
-- ─────────────────────────────────────────────────────────────────────────

SET @has_col := (SELECT COUNT(*) FROM information_schema.columns
                 WHERE table_schema = DATABASE()
                   AND table_name = 'phi_rates'
                   AND column_name = 'lookup_key');
SET @sql := IF(@has_col = 0,
  'ALTER TABLE phi_rates ADD COLUMN lookup_key VARCHAR(512) GENERATED ALWAYS AS (CONCAT(COALESCE(risk_rate_code, ''''), ''|'', COALESCE(risk_type, ''''), ''|'', CAST(age_next_birthday AS CHAR), ''|'', COALESCE(gender, ''''), ''|'', CAST(occupation_class AS CHAR), ''|'', CAST(income_level AS CHAR), ''|'', CAST(waiting_period AS CHAR), ''|'', CAST(deferred_period AS CHAR), ''|'', CAST(normal_retirement_age AS CHAR), ''|'', COALESCE(benefit_escalation_option, ''''), ''|'', COALESCE(disability_definition, ''''))) STORED',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ─────────────────────────────────────────────────────────────────────────
-- 2) reinsurance_phi_rates.lookup_key
-- ─────────────────────────────────────────────────────────────────────────

SET @has_col := (SELECT COUNT(*) FROM information_schema.columns
                 WHERE table_schema = DATABASE()
                   AND table_name = 'reinsurance_phi_rates'
                   AND column_name = 'lookup_key');
SET @sql := IF(@has_col = 0,
  'ALTER TABLE reinsurance_phi_rates ADD COLUMN lookup_key VARCHAR(512) GENERATED ALWAYS AS (CONCAT(COALESCE(risk_rate_code, ''''), ''|'', CAST(age_next_birthday AS CHAR), ''|'', CAST(income_level AS CHAR), ''|'', COALESCE(gender, ''''), ''|'', CAST(occupation_class AS CHAR))) STORED',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ─────────────────────────────────────────────────────────────────────────
-- 3) Index on phi_rates.lookup_key
-- ─────────────────────────────────────────────────────────────────────────

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE()
               AND table_name = 'phi_rates'
               AND index_name = 'idx_phi_rates_lookup_key');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_phi_rates_lookup_key ON phi_rates (lookup_key)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ─────────────────────────────────────────────────────────────────────────
-- 4) Index on reinsurance_phi_rates.lookup_key
-- ─────────────────────────────────────────────────────────────────────────

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE()
               AND table_name = 'reinsurance_phi_rates'
               AND index_name = 'idx_reins_phi_rates_lookup_key');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_reins_phi_rates_lookup_key ON reinsurance_phi_rates (lookup_key)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
