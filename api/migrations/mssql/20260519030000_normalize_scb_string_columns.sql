-- Normalize Salary Continuation Benefit string columns to NVARCHAR(255).
-- The original SCB migration's CREATE TABLE already declared NVARCHAR(255)
-- on MSSQL, so on most installs this is effectively a no-op; it is
-- included for parity with the MySQL/PostgreSQL files and to fix installs
-- whose tables were bootstrapped by GORM AutoMigrate (which maps `string`
-- to NVARCHAR(MAX)).
--
-- MSSQL allows ALTER COLUMN on a column referenced by a PERSISTED computed
-- column as long as the new type does not invalidate the expression.
-- NVARCHAR(255) is compatible with CONCAT(...).
--
-- One column per ALTER statement — MSSQL does not support multi-column
-- ALTER COLUMN in a single batch.

ALTER TABLE salary_continuation_rates ALTER COLUMN risk_rate_code            NVARCHAR(255) NULL;
ALTER TABLE salary_continuation_rates ALTER COLUMN gender                    NVARCHAR(255) NULL;
ALTER TABLE salary_continuation_rates ALTER COLUMN benefit_escalation_option NVARCHAR(255) NULL;
ALTER TABLE salary_continuation_rates ALTER COLUMN disability_definition     NVARCHAR(255) NULL;
ALTER TABLE salary_continuation_rates ALTER COLUMN risk_type                 NVARCHAR(255) NULL;
ALTER TABLE salary_continuation_rates ALTER COLUMN created_by                NVARCHAR(255) NULL;

ALTER TABLE reinsurance_salary_continuation_rates ALTER COLUMN risk_rate_code            NVARCHAR(255) NULL;
ALTER TABLE reinsurance_salary_continuation_rates ALTER COLUMN risk_type                 NVARCHAR(255) NULL;
ALTER TABLE reinsurance_salary_continuation_rates ALTER COLUMN gender                    NVARCHAR(255) NULL;
ALTER TABLE reinsurance_salary_continuation_rates ALTER COLUMN occupation_class          NVARCHAR(255) NULL;
ALTER TABLE reinsurance_salary_continuation_rates ALTER COLUMN income_level              NVARCHAR(255) NULL;
ALTER TABLE reinsurance_salary_continuation_rates ALTER COLUMN benefit_escalation_option NVARCHAR(255) NULL;
ALTER TABLE reinsurance_salary_continuation_rates ALTER COLUMN disability_definition     NVARCHAR(255) NULL;
ALTER TABLE reinsurance_salary_continuation_rates ALTER COLUMN created_by                NVARCHAR(255) NULL;
