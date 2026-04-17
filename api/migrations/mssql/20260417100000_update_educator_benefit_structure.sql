-- Migration: move educator benefit selection from group_pricing_parameters
-- to the per-scheme-category level, and enforce uniqueness on the
-- educator benefit structure.

-- ---------------------------------------------------------------------------
-- 1. scheme_categories: new per-benefit educator type columns
-- ---------------------------------------------------------------------------
IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'gla_educator_benefit_type')
BEGIN
    ALTER TABLE scheme_categories ADD gla_educator_benefit_type NVARCHAR(255);
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_educator_benefit_type NVARCHAR(255);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'ptd_educator_benefit_type')
BEGIN
    ALTER TABLE scheme_categories ADD ptd_educator_benefit_type NVARCHAR(255);
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ptd_educator_benefit_type NVARCHAR(255);
END

-- ---------------------------------------------------------------------------
-- 2. educator_benefit_structures: composite unique index on (rrc, code)
-- ---------------------------------------------------------------------------
-- Abort loudly if duplicates exist so they can be cleaned up first.
IF EXISTS (
    SELECT 1
    FROM educator_benefit_structures
    GROUP BY risk_rate_code, educator_benefit_code
    HAVING COUNT(*) > 1
)
BEGIN
    RAISERROR('educator_benefit_structures has duplicate (risk_rate_code, educator_benefit_code) rows; dedupe before re-running this migration.', 16, 1);
END

IF NOT EXISTS(SELECT 1 FROM sys.indexes WHERE object_id = OBJECT_ID('educator_benefit_structures') AND name = 'idx_educator_rrc_code')
BEGIN
    CREATE UNIQUE INDEX idx_educator_rrc_code
        ON educator_benefit_structures (risk_rate_code, educator_benefit_code);
END

-- ---------------------------------------------------------------------------
-- 3. group_pricing_parameters: drop educator_benefit_code
-- ---------------------------------------------------------------------------
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_parameters') AND name = 'educator_benefit_code')
BEGIN
    -- Drop any default constraint tied to the column first (SQL Server quirk).
    DECLARE @constraint_name SYSNAME;
    SELECT @constraint_name = dc.name
    FROM sys.default_constraints dc
    INNER JOIN sys.columns c ON c.default_object_id = dc.object_id
    WHERE c.object_id = OBJECT_ID('group_pricing_parameters')
      AND c.name = 'educator_benefit_code';
    IF @constraint_name IS NOT NULL
    BEGIN
        EXEC('ALTER TABLE group_pricing_parameters DROP CONSTRAINT ' + @constraint_name);
    END

    ALTER TABLE group_pricing_parameters DROP COLUMN educator_benefit_code;
END
