-- Migration: add per-benefit conversion flags to scheme_categories.
--   gla_conversion_on_withdrawal  (GLA)
--   gla_conversion_on_retirement  (GLA, retirement-only)
--   ptd_conversion_on_withdrawal  (PTD)
--   ci_conversion_on_withdrawal   (CI)

ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_conversion_on_withdrawal BOOLEAN DEFAULT FALSE;
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_conversion_on_withdrawal') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_conversion_on_withdrawal TYPE BOOLEAN USING gla_conversion_on_withdrawal::BOOLEAN;
    END IF;
END $$;

ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_conversion_on_retirement BOOLEAN DEFAULT FALSE;
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_conversion_on_retirement') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_conversion_on_retirement TYPE BOOLEAN USING gla_conversion_on_retirement::BOOLEAN;
    END IF;
END $$;

ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_conversion_on_withdrawal BOOLEAN DEFAULT FALSE;
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ptd_conversion_on_withdrawal') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ptd_conversion_on_withdrawal TYPE BOOLEAN USING ptd_conversion_on_withdrawal::BOOLEAN;
    END IF;
END $$;

ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ci_conversion_on_withdrawal BOOLEAN DEFAULT FALSE;
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ci_conversion_on_withdrawal') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ci_conversion_on_withdrawal TYPE BOOLEAN USING ci_conversion_on_withdrawal::BOOLEAN;
    END IF;
END $$;
