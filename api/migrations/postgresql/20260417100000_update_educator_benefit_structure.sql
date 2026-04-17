-- Migration: move educator benefit selection from group_pricing_parameters
-- to the per-scheme-category level, and enforce uniqueness on the
-- educator benefit structure.
--
-- Steps:
--   1. Add gla_educator_benefit_type and ptd_educator_benefit_type columns
--      to scheme_categories.
--   2. Add a composite unique index (risk_rate_code, educator_benefit_code)
--      to educator_benefit_structures.
--   3. Drop educator_benefit_code from group_pricing_parameters.

-- ---------------------------------------------------------------------------
-- 1. scheme_categories: new per-benefit educator type columns
-- ---------------------------------------------------------------------------
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_educator_benefit_type VARCHAR(255);
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_educator_benefit_type') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_educator_benefit_type TYPE VARCHAR(255);
    END IF;
END $$;

ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_educator_benefit_type VARCHAR(255);
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ptd_educator_benefit_type') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ptd_educator_benefit_type TYPE VARCHAR(255);
    END IF;
END $$;

-- ---------------------------------------------------------------------------
-- 2. educator_benefit_structures: composite unique index on (rrc, code)
-- ---------------------------------------------------------------------------
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM educator_benefit_structures
        GROUP BY risk_rate_code, educator_benefit_code
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION 'educator_benefit_structures has duplicate (risk_rate_code, educator_benefit_code) rows; dedupe before re-running this migration.';
    END IF;
END $$;

CREATE UNIQUE INDEX IF NOT EXISTS idx_educator_rrc_code
    ON educator_benefit_structures (risk_rate_code, educator_benefit_code);

-- ---------------------------------------------------------------------------
-- 3. group_pricing_parameters: drop educator_benefit_code
-- ---------------------------------------------------------------------------
ALTER TABLE group_pricing_parameters DROP COLUMN IF EXISTS educator_benefit_code;
