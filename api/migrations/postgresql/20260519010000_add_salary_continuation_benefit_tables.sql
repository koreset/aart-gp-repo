-- Migration: create salary_continuation_rates + reinsurance_salary_continuation_rates
-- rating tables for the Salary Continuation Benefit (SCB) — a new optional
-- benefit shown under PHI in the quote builder.
--
-- SCB mirrors the PHI rating factors but replaces `normal_retirement_age`
-- with `excess_period`. Both tables carry an indexed STORED generated
-- `lookup_key` column so per-member rate fetches become single-column
-- index seeks (same pattern as phi_rates).
--
-- Also adds `scb_benefit` and `scb_excess_period` columns to scheme_categories.
--
-- The concatenation order and `|` separator in lookup_key MUST match the
-- string built by GetScbRate / GetReinsuranceScbRate in
-- api/services/group_pricing.go.
-- Idempotent on re-runs.

-- 1) salary_continuation_rates table
CREATE TABLE IF NOT EXISTS salary_continuation_rates (
    id                         SERIAL PRIMARY KEY,
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
    scb_rate                   DOUBLE PRECISION,
    creation_date              TIMESTAMP WITH TIME ZONE,
    created_by                 VARCHAR(255)
);

-- 2) reinsurance_salary_continuation_rates table
CREATE TABLE IF NOT EXISTS reinsurance_salary_continuation_rates (
    id                         SERIAL PRIMARY KEY,
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
    scb_rate                   DOUBLE PRECISION,
    creation_date              TIMESTAMP WITH TIME ZONE,
    created_by                 VARCHAR(255)
);

-- 3) scheme_categories: scb_benefit + scb_excess_period
ALTER TABLE scheme_categories
    ADD COLUMN IF NOT EXISTS scb_benefit BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE scheme_categories
    ADD COLUMN IF NOT EXISTS scb_excess_period BIGINT NOT NULL DEFAULT 0;

-- 4) salary_continuation_rates.lookup_key
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

-- 5) reinsurance_salary_continuation_rates.lookup_key
ALTER TABLE reinsurance_salary_continuation_rates
    ADD COLUMN IF NOT EXISTS lookup_key TEXT
    GENERATED ALWAYS AS (
        COALESCE(risk_rate_code, '') || '|' ||
        age_next_birthday::TEXT || '|' ||
        COALESCE(income_level, '') || '|' ||
        COALESCE(gender, '') || '|' ||
        COALESCE(occupation_class, '')
    ) STORED;

-- 6) Indexes
CREATE INDEX IF NOT EXISTS idx_scb_rates_lookup_key
    ON salary_continuation_rates (lookup_key);

CREATE INDEX IF NOT EXISTS idx_reins_scb_rates_lookup_key
    ON reinsurance_salary_continuation_rates (lookup_key);
