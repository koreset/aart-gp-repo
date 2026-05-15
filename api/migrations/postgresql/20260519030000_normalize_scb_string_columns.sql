-- Normalize Salary Continuation Benefit string columns to VARCHAR(255).
-- If the table was bootstrapped via GORM AutoMigrate the columns came up
-- as TEXT; this realigns them with the migration's original intent.
--
-- ALTER COLUMN ... TYPE is allowed in PostgreSQL even when the column is
-- referenced by a STORED generated column, provided the new type is
-- compatible with the expression (TEXT/VARCHAR are both string types and
-- accepted by || / COALESCE).
--
-- Re-running this migration is harmless — converting VARCHAR(255) to
-- VARCHAR(255) is a no-op.

ALTER TABLE salary_continuation_rates
    ALTER COLUMN risk_rate_code            TYPE VARCHAR(255),
    ALTER COLUMN gender                    TYPE VARCHAR(255),
    ALTER COLUMN benefit_escalation_option TYPE VARCHAR(255),
    ALTER COLUMN disability_definition     TYPE VARCHAR(255),
    ALTER COLUMN risk_type                 TYPE VARCHAR(255),
    ALTER COLUMN created_by                TYPE VARCHAR(255);

ALTER TABLE reinsurance_salary_continuation_rates
    ALTER COLUMN risk_rate_code            TYPE VARCHAR(255),
    ALTER COLUMN risk_type                 TYPE VARCHAR(255),
    ALTER COLUMN gender                    TYPE VARCHAR(255),
    ALTER COLUMN occupation_class          TYPE VARCHAR(255),
    ALTER COLUMN income_level              TYPE VARCHAR(255),
    ALTER COLUMN benefit_escalation_option TYPE VARCHAR(255),
    ALTER COLUMN disability_definition     TYPE VARCHAR(255),
    ALTER COLUMN created_by                TYPE VARCHAR(255);
