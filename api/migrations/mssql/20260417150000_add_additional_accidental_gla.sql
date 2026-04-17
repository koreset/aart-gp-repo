-- Migration: add optional "Additional Accidental GLA" sub-benefit.
-- Re-uses the gla_rates / gla_aids_rates tables with a different benefit_type.
-- All GLA scheme parameters (salary multiple, terminal illness, waiting period,
-- loadings, educator, conversion) stay on the GLA fields — only the benefit_type
-- differs for the Additional Accidental layer.

--------------------------------------------------------------------------------
-- scheme_categories: toggle + benefit type for the Additional Accidental layer
--------------------------------------------------------------------------------

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'additional_accidental_gla_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD additional_accidental_gla_benefit BIT DEFAULT 0;
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_accidental_gla_benefit BIT;
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'additional_accidental_gla_benefit_type')
BEGIN
    ALTER TABLE scheme_categories ADD additional_accidental_gla_benefit_type NVARCHAR(255);
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_accidental_gla_benefit_type NVARCHAR(255);
END

--------------------------------------------------------------------------------
-- member_rating_results: per-member Additional Accidental GLA outputs
--------------------------------------------------------------------------------

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'additional_accidental_gla_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD additional_accidental_gla_sum_assured DECIMAL(15,5);
END
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN additional_accidental_gla_sum_assured DECIMAL(15,5);
END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'additional_accidental_gla_capped_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD additional_accidental_gla_capped_sum_assured DECIMAL(15,5);
END
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN additional_accidental_gla_capped_sum_assured DECIMAL(15,5);
END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'additional_accidental_gla_qx')
BEGIN
    ALTER TABLE member_rating_results ADD additional_accidental_gla_qx DECIMAL(15,5);
END
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN additional_accidental_gla_qx DECIMAL(15,5);
END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'additional_accidental_gla_aids_qx')
BEGIN
    ALTER TABLE member_rating_results ADD additional_accidental_gla_aids_qx DECIMAL(15,5);
END
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN additional_accidental_gla_aids_qx DECIMAL(15,5);
END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'base_additional_accidental_gla_rate')
BEGIN
    ALTER TABLE member_rating_results ADD base_additional_accidental_gla_rate DECIMAL(15,5);
END
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN base_additional_accidental_gla_rate DECIMAL(15,5);
END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'loaded_additional_accidental_gla_rate')
BEGIN
    ALTER TABLE member_rating_results ADD loaded_additional_accidental_gla_rate DECIMAL(15,5);
END
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN loaded_additional_accidental_gla_rate DECIMAL(15,5);
END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'additional_accidental_gla_weighted_experience_crude_rate')
BEGIN
    ALTER TABLE member_rating_results ADD additional_accidental_gla_weighted_experience_crude_rate DECIMAL(15,5);
END
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN additional_accidental_gla_weighted_experience_crude_rate DECIMAL(15,5);
END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'additional_accidental_gla_theoretical_rate')
BEGIN
    ALTER TABLE member_rating_results ADD additional_accidental_gla_theoretical_rate DECIMAL(15,5);
END
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN additional_accidental_gla_theoretical_rate DECIMAL(15,5);
END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'additional_accidental_gla_experience_adjustment')
BEGIN
    ALTER TABLE member_rating_results ADD additional_accidental_gla_experience_adjustment DECIMAL(15,5);
END
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN additional_accidental_gla_experience_adjustment DECIMAL(15,5);
END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_loaded_additional_accidental_gla_rate')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_loaded_additional_accidental_gla_rate DECIMAL(15,5);
END
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_loaded_additional_accidental_gla_rate DECIMAL(15,5);
END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'additional_accidental_gla_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD additional_accidental_gla_risk_premium DECIMAL(15,5);
END
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN additional_accidental_gla_risk_premium DECIMAL(15,5);
END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_additional_accidental_gla_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_additional_accidental_gla_risk_premium DECIMAL(15,5);
END
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_additional_accidental_gla_risk_premium DECIMAL(15,5);
END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'additional_accidental_gla_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD additional_accidental_gla_office_premium DECIMAL(15,5);
END
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN additional_accidental_gla_office_premium DECIMAL(15,5);
END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_additional_accidental_gla_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_additional_accidental_gla_office_premium DECIMAL(15,5);
END
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_additional_accidental_gla_office_premium DECIMAL(15,5);
END
