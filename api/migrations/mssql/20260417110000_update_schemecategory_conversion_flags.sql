-- Migration: add per-benefit conversion flags to scheme_categories.
--   gla_conversion_on_withdrawal  (GLA)
--   gla_conversion_on_retirement  (GLA, retirement-only)
--   ptd_conversion_on_withdrawal  (PTD)
--   ci_conversion_on_withdrawal   (CI)

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'gla_conversion_on_withdrawal')
BEGIN
    ALTER TABLE scheme_categories ADD gla_conversion_on_withdrawal BIT DEFAULT 0;
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_conversion_on_withdrawal BIT;
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'gla_conversion_on_retirement')
BEGIN
    ALTER TABLE scheme_categories ADD gla_conversion_on_retirement BIT DEFAULT 0;
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_conversion_on_retirement BIT;
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'ptd_conversion_on_withdrawal')
BEGIN
    ALTER TABLE scheme_categories ADD ptd_conversion_on_withdrawal BIT DEFAULT 0;
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ptd_conversion_on_withdrawal BIT;
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'ci_conversion_on_withdrawal')
BEGIN
    ALTER TABLE scheme_categories ADD ci_conversion_on_withdrawal BIT DEFAULT 0;
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ci_conversion_on_withdrawal BIT;
END
