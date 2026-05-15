-- Migration: create salary_continuation_rates + reinsurance_salary_continuation_rates
-- rating tables for the Salary Continuation Benefit (SCB) — a new optional
-- benefit shown under PHI in the quote builder.
--
-- SCB mirrors the PHI rating factors but replaces `normal_retirement_age`
-- with `excess_period`. Both tables carry an indexed STORED generated
-- `lookup_key` column so per-member rate fetches become single-column
-- index seeks (same pattern as phi_rates).
--
-- Also adds `scb_benefit` and `scb_excess_period` columns to scheme_categories
-- so a quote can opt-in to SCB and record the chosen excess period.
--
-- The concatenation order and `|` separator in lookup_key MUST match the
-- string built by GetScbRate / GetReinsuranceScbRate in
-- api/services/group_pricing.go.
-- Idempotent on re-runs.

-- ─────────────────────────────────────────────────────────────────────────
-- 1) salary_continuation_rates table
-- ─────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS salary_continuation_rates (
    id                         BIGINT AUTO_INCREMENT PRIMARY KEY,
    risk_rate_code             VARCHAR(255),
    age_next_birthday          BIGINT,
    gender                     VARCHAR(255),
    occupation_class           BIGINT,
    income_level               BIGINT,
    deferred_period            BIGINT,
    excess_period              BIGINT,
    benefit_escalation_option  VARCHAR(255),
    disability_definition      VARCHAR(255),
    risk_type                  VARCHAR(255),
    scb_rate                   DOUBLE,
    creation_date              DATETIME(3),
    created_by                 VARCHAR(255)
);

-- ─────────────────────────────────────────────────────────────────────────
-- 2) reinsurance_salary_continuation_rates table
-- ─────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS reinsurance_salary_continuation_rates (
    id                         BIGINT AUTO_INCREMENT PRIMARY KEY,
    risk_rate_code             VARCHAR(255),
    risk_type                  VARCHAR(255),
    age_next_birthday          BIGINT,
    gender                     VARCHAR(255),
    occupation_class           VARCHAR(255),
    income_level               VARCHAR(255),
    deferred_period            BIGINT,
    excess_period              BIGINT,
    benefit_escalation_option  VARCHAR(255),
    disability_definition      VARCHAR(255),
    scb_rate                   DOUBLE,
    creation_date              DATETIME(3),
    created_by                 VARCHAR(255)
);

-- ─────────────────────────────────────────────────────────────────────────
-- 3) scheme_categories: scb_benefit + scb_excess_period
-- ─────────────────────────────────────────────────────────────────────────
SET @has_col := (SELECT COUNT(*) FROM information_schema.columns
                 WHERE table_schema = DATABASE()
                   AND table_name = 'scheme_categories'
                   AND column_name = 'scb_benefit');
SET @sql := IF(@has_col = 0,
  'ALTER TABLE scheme_categories ADD COLUMN scb_benefit TINYINT(1) NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @has_col := (SELECT COUNT(*) FROM information_schema.columns
                 WHERE table_schema = DATABASE()
                   AND table_name = 'scheme_categories'
                   AND column_name = 'scb_excess_period');
SET @sql := IF(@has_col = 0,
  'ALTER TABLE scheme_categories ADD COLUMN scb_excess_period BIGINT NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ─────────────────────────────────────────────────────────────────────────
-- 4) salary_continuation_rates.lookup_key
-- ─────────────────────────────────────────────────────────────────────────
SET @has_col := (SELECT COUNT(*) FROM information_schema.columns
                 WHERE table_schema = DATABASE()
                   AND table_name = 'salary_continuation_rates'
                   AND column_name = 'lookup_key');
SET @sql := IF(@has_col = 0,
  'ALTER TABLE salary_continuation_rates ADD COLUMN lookup_key VARCHAR(512) GENERATED ALWAYS AS (CONCAT(COALESCE(risk_rate_code, ''''), ''|'', COALESCE(risk_type, ''''), ''|'', CAST(age_next_birthday AS CHAR), ''|'', COALESCE(gender, ''''), ''|'', CAST(occupation_class AS CHAR), ''|'', CAST(income_level AS CHAR), ''|'', CAST(deferred_period AS CHAR), ''|'', CAST(excess_period AS CHAR), ''|'', COALESCE(benefit_escalation_option, ''''), ''|'', COALESCE(disability_definition, ''''))) STORED',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ─────────────────────────────────────────────────────────────────────────
-- 5) reinsurance_salary_continuation_rates.lookup_key
-- ─────────────────────────────────────────────────────────────────────────
SET @has_col := (SELECT COUNT(*) FROM information_schema.columns
                 WHERE table_schema = DATABASE()
                   AND table_name = 'reinsurance_salary_continuation_rates'
                   AND column_name = 'lookup_key');
SET @sql := IF(@has_col = 0,
  'ALTER TABLE reinsurance_salary_continuation_rates ADD COLUMN lookup_key VARCHAR(512) GENERATED ALWAYS AS (CONCAT(COALESCE(risk_rate_code, ''''), ''|'', CAST(age_next_birthday AS CHAR), ''|'', COALESCE(income_level, ''''), ''|'', COALESCE(gender, ''''), ''|'', COALESCE(occupation_class, ''''))) STORED',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ─────────────────────────────────────────────────────────────────────────
-- 6) Index on salary_continuation_rates.lookup_key
-- ─────────────────────────────────────────────────────────────────────────
SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE()
               AND table_name = 'salary_continuation_rates'
               AND index_name = 'idx_scb_rates_lookup_key');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_scb_rates_lookup_key ON salary_continuation_rates (lookup_key)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ─────────────────────────────────────────────────────────────────────────
-- 7) Index on reinsurance_salary_continuation_rates.lookup_key
-- ─────────────────────────────────────────────────────────────────────────
SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE()
               AND table_name = 'reinsurance_salary_continuation_rates'
               AND index_name = 'idx_reins_scb_rates_lookup_key');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_reins_scb_rates_lookup_key ON reinsurance_salary_continuation_rates (lookup_key)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
