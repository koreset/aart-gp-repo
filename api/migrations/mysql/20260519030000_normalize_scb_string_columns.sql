-- Normalize Salary Continuation Benefit string columns from LONGTEXT
-- (assigned by GORM AutoMigrate's default `string` mapping when the table
-- was bootstrapped via EnsureGPTableStats) to the VARCHAR(255) shape the
-- original migration intended.
--
-- MODIFY COLUMN is safe here even though `lookup_key` is a STORED generated
-- column that references these source columns: MySQL re-validates the
-- expression in place, and CONCAT(...) accepts both LONGTEXT and VARCHAR
-- inputs so no drop/recreate dance is required.
--
-- Re-running this migration is harmless — MODIFY of an already-VARCHAR(255)
-- column is a no-op.

ALTER TABLE salary_continuation_rates
    MODIFY COLUMN risk_rate_code            VARCHAR(255),
    MODIFY COLUMN gender                    VARCHAR(255),
    MODIFY COLUMN benefit_escalation_option VARCHAR(255),
    MODIFY COLUMN disability_definition     VARCHAR(255),
    MODIFY COLUMN risk_type                 VARCHAR(255),
    MODIFY COLUMN created_by                VARCHAR(255);

ALTER TABLE reinsurance_salary_continuation_rates
    MODIFY COLUMN risk_rate_code            VARCHAR(255),
    MODIFY COLUMN risk_type                 VARCHAR(255),
    MODIFY COLUMN gender                    VARCHAR(255),
    MODIFY COLUMN occupation_class          VARCHAR(255),
    MODIFY COLUMN income_level              VARCHAR(255),
    MODIFY COLUMN benefit_escalation_option VARCHAR(255),
    MODIFY COLUMN disability_definition     VARCHAR(255),
    MODIFY COLUMN created_by                VARCHAR(255);
