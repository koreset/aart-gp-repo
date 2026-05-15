-- Migration: create salary_continuation_rates + reinsurance_salary_continuation_rates
-- rating tables for the Salary Continuation Benefit (SCB) — a new optional
-- benefit shown under PHI in the quote builder.
--
-- SCB mirrors the PHI rating factors but replaces `normal_retirement_age`
-- with `excess_period`. Both tables carry an indexed PERSISTED computed
-- `lookup_key` column so per-member rate fetches become single-column
-- index seeks (same pattern as phi_rates).
--
-- Also adds `scb_benefit` and `scb_excess_period` columns to scheme_categories.
--
-- The concatenation order and `|` separator in lookup_key MUST match the
-- string built by GetScbRate / GetReinsuranceScbRate in
-- api/services/group_pricing.go. MSSQL's CONCAT treats NULL as an empty
-- string, so no explicit COALESCE is needed.
-- Idempotent on re-runs.

-- 1) salary_continuation_rates table
IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'salary_continuation_rates')
BEGIN
    CREATE TABLE salary_continuation_rates (
        id                         BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT pk_salary_continuation_rates PRIMARY KEY,
        risk_rate_code             NVARCHAR(255) NULL,
        age_next_birthday          BIGINT NULL,
        gender                     NVARCHAR(255) NULL,
        occupation_class           BIGINT NULL,
        income_level               BIGINT NULL,
        deferred_period            BIGINT NULL,
        excess_period              BIGINT NULL,
        benefit_escalation_option  NVARCHAR(255) NULL,
        disability_definition      NVARCHAR(255) NULL,
        risk_type                  NVARCHAR(255) NULL,
        scb_rate                   FLOAT NULL,
        creation_date              DATETIMEOFFSET NULL,
        created_by                 NVARCHAR(255) NULL
    );
END;

-- 2) reinsurance_salary_continuation_rates table
IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'reinsurance_salary_continuation_rates')
BEGIN
    CREATE TABLE reinsurance_salary_continuation_rates (
        id                         BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT pk_reinsurance_salary_continuation_rates PRIMARY KEY,
        risk_rate_code             NVARCHAR(255) NULL,
        risk_type                  NVARCHAR(255) NULL,
        age_next_birthday          BIGINT NULL,
        gender                     NVARCHAR(255) NULL,
        occupation_class           NVARCHAR(255) NULL,
        income_level               NVARCHAR(255) NULL,
        deferred_period            BIGINT NULL,
        excess_period              BIGINT NULL,
        benefit_escalation_option  NVARCHAR(255) NULL,
        disability_definition      NVARCHAR(255) NULL,
        scb_rate                   FLOAT NULL,
        creation_date              DATETIMEOFFSET NULL,
        created_by                 NVARCHAR(255) NULL
    );
END;

-- 3) scheme_categories: scb_benefit + scb_excess_period
IF NOT EXISTS (SELECT 1 FROM sys.columns
               WHERE object_id = OBJECT_ID('scheme_categories')
                 AND name = 'scb_benefit')
BEGIN
    ALTER TABLE scheme_categories
        ADD scb_benefit BIT NOT NULL CONSTRAINT df_scheme_categories_scb_benefit DEFAULT 0;
END;

IF NOT EXISTS (SELECT 1 FROM sys.columns
               WHERE object_id = OBJECT_ID('scheme_categories')
                 AND name = 'scb_excess_period')
BEGIN
    ALTER TABLE scheme_categories
        ADD scb_excess_period BIGINT NOT NULL CONSTRAINT df_scheme_categories_scb_excess_period DEFAULT 0;
END;

-- 4) salary_continuation_rates.lookup_key
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

-- 5) reinsurance_salary_continuation_rates.lookup_key
IF NOT EXISTS (SELECT 1 FROM sys.columns
               WHERE object_id = OBJECT_ID('reinsurance_salary_continuation_rates')
                 AND name = 'lookup_key')
BEGIN
    ALTER TABLE reinsurance_salary_continuation_rates
        ADD lookup_key AS (
            CONCAT(
                risk_rate_code, '|',
                CAST(age_next_birthday AS VARCHAR(20)), '|',
                income_level, '|',
                gender, '|',
                occupation_class
            )
        ) PERSISTED;
END;

-- 6) Indexes
IF NOT EXISTS (SELECT 1 FROM sys.indexes
               WHERE name = 'idx_scb_rates_lookup_key'
                 AND object_id = OBJECT_ID('salary_continuation_rates'))
    CREATE INDEX idx_scb_rates_lookup_key
        ON salary_continuation_rates (lookup_key);

IF NOT EXISTS (SELECT 1 FROM sys.indexes
               WHERE name = 'idx_reins_scb_rates_lookup_key'
                 AND object_id = OBJECT_ID('reinsurance_salary_continuation_rates'))
    CREATE INDEX idx_reins_scb_rates_lookup_key
        ON reinsurance_salary_continuation_rates (lookup_key);
