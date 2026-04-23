-- Migration for struct: SchemeCategory

-- Table: scheme_categories

-- Ensure table exists
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'scheme_categories')
BEGIN
    CREATE TABLE scheme_categories (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

-- Add or modify column for field: QuoteId
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'quote_id')
BEGIN
    ALTER TABLE scheme_categories ADD quote_id INT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN quote_id INT;
END;

-- Add or modify column for field: SchemeCategory
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'scheme_category')
BEGIN
    ALTER TABLE scheme_categories ADD scheme_category NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN scheme_category NVARCHAR(255);
END;

-- Add or modify column for field: Basis
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'basis')
BEGIN
    ALTER TABLE scheme_categories ADD basis NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN basis NVARCHAR(255);
END;

-- Add or modify column for field: FreeCoverLimit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'free_cover_limit')
BEGIN
    ALTER TABLE scheme_categories ADD free_cover_limit DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN free_cover_limit DECIMAL(15,5);
END;

-- Add or modify column for field: PtdBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ptd_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD ptd_benefit BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ptd_benefit BIT;
END;

-- Add or modify column for field: GlaBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD gla_benefit BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_benefit BIT;
END;

-- Add or modify column for field: CiBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ci_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD ci_benefit BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ci_benefit BIT;
END;

-- Add or modify column for field: SglaBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'sgla_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD sgla_benefit BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN sgla_benefit BIT;
END;

-- Add or modify column for field: PhiBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD phi_benefit BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_benefit BIT;
END;

-- Add or modify column for field: TtdBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ttd_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD ttd_benefit BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ttd_benefit BIT;
END;

-- Add or modify column for field: FamilyFuneralBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'family_funeral_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD family_funeral_benefit BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN family_funeral_benefit BIT;
END;

-- Add or modify column for field: GlaSalaryMultiple
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_salary_multiple')
BEGIN
    ALTER TABLE scheme_categories ADD gla_salary_multiple DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_salary_multiple DECIMAL(15,5);
END;

-- Add or modify column for field: GlaTerminalIllnessBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_terminal_illness_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD gla_terminal_illness_benefit NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_terminal_illness_benefit NVARCHAR(255);
END;

-- Add or modify column for field: GlaWaitingPeriod
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_waiting_period')
BEGIN
    ALTER TABLE scheme_categories ADD gla_waiting_period INT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_waiting_period INT;
END;

-- Add or modify column for field: GlaEducatorBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_educator_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD gla_educator_benefit NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_educator_benefit NVARCHAR(255);
END;

-- Add or modify column for field: GlaEducatorBenefitType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_educator_benefit_type')
BEGIN
    ALTER TABLE scheme_categories ADD gla_educator_benefit_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_educator_benefit_type NVARCHAR(255);
END;

-- Add or modify column for field: GlaBenefitType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_benefit_type')
BEGIN
    ALTER TABLE scheme_categories ADD gla_benefit_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_benefit_type NVARCHAR(255);
END;

-- Add or modify column for field: GlaConversionOnWithdrawal
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_conversion_on_withdrawal')
BEGIN
    ALTER TABLE scheme_categories ADD gla_conversion_on_withdrawal BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_conversion_on_withdrawal BIT;
END;

-- Add or modify column for field: GlaConversionOnRetirement
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_conversion_on_retirement')
BEGIN
    ALTER TABLE scheme_categories ADD gla_conversion_on_retirement BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_conversion_on_retirement BIT;
END;

-- Add or modify column for field: AdditionalAccidentalGlaBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'additional_accidental_gla_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD additional_accidental_gla_benefit BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_accidental_gla_benefit BIT;
END;

-- Add or modify column for field: AdditionalAccidentalGlaBenefitType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'additional_accidental_gla_benefit_type')
BEGIN
    ALTER TABLE scheme_categories ADD additional_accidental_gla_benefit_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_accidental_gla_benefit_type NVARCHAR(255);
END;

-- Add or modify column for field: TaxSaverBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'tax_saver_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD tax_saver_benefit BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN tax_saver_benefit BIT;
END;

-- Add or modify column for field: AdditionalGlaCoverBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'additional_gla_cover_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD additional_gla_cover_benefit BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_benefit BIT;
END;

-- Add or modify column for field: AdditionalGlaCoverAgeBandSource
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'additional_gla_cover_age_band_source')
BEGIN
    ALTER TABLE scheme_categories ADD additional_gla_cover_age_band_source NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_age_band_source NVARCHAR(255);
END;

-- Add or modify column for field: AdditionalGlaCoverAgeBandType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'additional_gla_cover_age_band_type')
BEGIN
    ALTER TABLE scheme_categories ADD additional_gla_cover_age_band_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_age_band_type NVARCHAR(255);
END;

-- Add or modify column for field: AdditionalGlaCoverCustomAgeBands
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'additional_gla_cover_custom_age_bands')
BEGIN
    ALTER TABLE scheme_categories ADD additional_gla_cover_custom_age_bands text;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_custom_age_bands text;
END;

-- Add or modify column for field: AdditionalGlaCoverBandRates
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'additional_gla_cover_band_rates')
BEGIN
    ALTER TABLE scheme_categories ADD additional_gla_cover_band_rates text;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_band_rates text;
END;

-- Add or modify column for field: AdditionalGlaCoverMalePropUsed
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'additional_gla_cover_male_prop_used')
BEGIN
    ALTER TABLE scheme_categories ADD additional_gla_cover_male_prop_used NVARCHAR(MAX);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_male_prop_used NVARCHAR(MAX);
END;

-- Add or modify column for field: PtdRiskType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ptd_risk_type')
BEGIN
    ALTER TABLE scheme_categories ADD ptd_risk_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ptd_risk_type NVARCHAR(255);
END;

-- Add or modify column for field: PtdBenefitType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ptd_benefit_type')
BEGIN
    ALTER TABLE scheme_categories ADD ptd_benefit_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ptd_benefit_type NVARCHAR(255);
END;

-- Add or modify column for field: PtdSalaryMultiple
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ptd_salary_multiple')
BEGIN
    ALTER TABLE scheme_categories ADD ptd_salary_multiple DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ptd_salary_multiple DECIMAL(15,5);
END;

-- Add or modify column for field: PtdDeferredPeriod
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ptd_deferred_period')
BEGIN
    ALTER TABLE scheme_categories ADD ptd_deferred_period INT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ptd_deferred_period INT;
END;

-- Add or modify column for field: PtdDisabilityDefinition
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ptd_disability_definition')
BEGIN
    ALTER TABLE scheme_categories ADD ptd_disability_definition NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ptd_disability_definition NVARCHAR(255);
END;

-- Add or modify column for field: PtdEducatorBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ptd_educator_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD ptd_educator_benefit NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ptd_educator_benefit NVARCHAR(255);
END;

-- Add or modify column for field: PtdEducatorBenefitType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ptd_educator_benefit_type')
BEGIN
    ALTER TABLE scheme_categories ADD ptd_educator_benefit_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ptd_educator_benefit_type NVARCHAR(255);
END;

-- Add or modify column for field: PtdConversionOnWithdrawal
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ptd_conversion_on_withdrawal')
BEGIN
    ALTER TABLE scheme_categories ADD ptd_conversion_on_withdrawal BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ptd_conversion_on_withdrawal BIT;
END;

-- Add or modify column for field: CiBenefitStructure
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ci_benefit_structure')
BEGIN
    ALTER TABLE scheme_categories ADD ci_benefit_structure NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ci_benefit_structure NVARCHAR(255);
END;

-- Add or modify column for field: CiBenefitDefinition
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ci_benefit_definition')
BEGIN
    ALTER TABLE scheme_categories ADD ci_benefit_definition NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ci_benefit_definition NVARCHAR(255);
END;

-- Add or modify column for field: CiMaxBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ci_max_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD ci_max_benefit DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ci_max_benefit DECIMAL(15,5);
END;

-- Add or modify column for field: CiCriticalIllnessSalaryMultiple
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ci_critical_illness_salary_multiple')
BEGIN
    ALTER TABLE scheme_categories ADD ci_critical_illness_salary_multiple DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ci_critical_illness_salary_multiple DECIMAL(15,5);
END;

-- Add or modify column for field: CiConversionOnWithdrawal
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ci_conversion_on_withdrawal')
BEGIN
    ALTER TABLE scheme_categories ADD ci_conversion_on_withdrawal BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ci_conversion_on_withdrawal BIT;
END;

-- Add or modify column for field: SglaSalaryMultiple
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'sgla_salary_multiple')
BEGIN
    ALTER TABLE scheme_categories ADD sgla_salary_multiple DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN sgla_salary_multiple DECIMAL(15,5);
END;

-- Add or modify column for field: SglaMaxBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'sgla_max_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD sgla_max_benefit DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN sgla_max_benefit DECIMAL(15,5);
END;

-- Add or modify column for field: PhiRiskType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_risk_type')
BEGIN
    ALTER TABLE scheme_categories ADD phi_risk_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_risk_type NVARCHAR(255);
END;

-- Add or modify column for field: PhiMaximumBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_maximum_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD phi_maximum_benefit DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_maximum_benefit DECIMAL(15,5);
END;

-- Add or modify column for field: PhiIncomeReplacementPercentage
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_income_replacement_percentage')
BEGIN
    ALTER TABLE scheme_categories ADD phi_income_replacement_percentage DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_income_replacement_percentage DECIMAL(15,5);
END;

-- Add or modify column for field: PhiUseTieredIncomeReplacementRatio
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_use_tiered_income_replacement_ratio')
BEGIN
    ALTER TABLE scheme_categories ADD phi_use_tiered_income_replacement_ratio BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_use_tiered_income_replacement_ratio BIT;
END;

-- Add or modify column for field: PhiTieredIncomeReplacementType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_tiered_income_replacement_type')
BEGIN
    ALTER TABLE scheme_categories ADD phi_tiered_income_replacement_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_tiered_income_replacement_type NVARCHAR(255);
END;

-- Add or modify column for field: PhiPremiumWaiver
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_premium_waiver')
BEGIN
    ALTER TABLE scheme_categories ADD phi_premium_waiver NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_premium_waiver NVARCHAR(255);
END;

-- Add or modify column for field: PhiMedicalAidPremiumWaiver
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_medical_aid_premium_waiver')
BEGIN
    ALTER TABLE scheme_categories ADD phi_medical_aid_premium_waiver NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_medical_aid_premium_waiver NVARCHAR(255);
END;

-- Add or modify column for field: PhiBenefitEscalation
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_benefit_escalation')
BEGIN
    ALTER TABLE scheme_categories ADD phi_benefit_escalation NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_benefit_escalation NVARCHAR(255);
END;

-- Add or modify column for field: PhiMaxPremiumWaiver
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_max_premium_waiver')
BEGIN
    ALTER TABLE scheme_categories ADD phi_max_premium_waiver DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_max_premium_waiver DECIMAL(15,5);
END;

-- Add or modify column for field: PhiWaitingPeriod
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_waiting_period')
BEGIN
    ALTER TABLE scheme_categories ADD phi_waiting_period INT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_waiting_period INT;
END;

-- Add or modify column for field: PhiNumberMonthlyPayments
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_number_monthly_payments')
BEGIN
    ALTER TABLE scheme_categories ADD phi_number_monthly_payments INT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_number_monthly_payments INT;
END;

-- Add or modify column for field: PhiDeferredPeriod
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_deferred_period')
BEGIN
    ALTER TABLE scheme_categories ADD phi_deferred_period INT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_deferred_period INT;
END;

-- Add or modify column for field: PhiDisabilityDefinition
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_disability_definition')
BEGIN
    ALTER TABLE scheme_categories ADD phi_disability_definition NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_disability_definition NVARCHAR(255);
END;

-- Add or modify column for field: PhiNormalRetirementAge
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_normal_retirement_age')
BEGIN
    ALTER TABLE scheme_categories ADD phi_normal_retirement_age INT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_normal_retirement_age INT;
END;

-- Add or modify column for field: TtdRiskType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ttd_risk_type')
BEGIN
    ALTER TABLE scheme_categories ADD ttd_risk_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ttd_risk_type NVARCHAR(255);
END;

-- Add or modify column for field: TtdMaximumBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ttd_maximum_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD ttd_maximum_benefit DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ttd_maximum_benefit DECIMAL(15,5);
END;

-- Add or modify column for field: TtdIncomeReplacementPercentage
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ttd_income_replacement_percentage')
BEGIN
    ALTER TABLE scheme_categories ADD ttd_income_replacement_percentage DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ttd_income_replacement_percentage DECIMAL(15,5);
END;

-- Add or modify column for field: TtdUseTieredIncomeReplacementRatio
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ttd_use_tiered_income_replacement_ratio')
BEGIN
    ALTER TABLE scheme_categories ADD ttd_use_tiered_income_replacement_ratio BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ttd_use_tiered_income_replacement_ratio BIT;
END;

-- Add or modify column for field: TtdTieredIncomeReplacementType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ttd_tiered_income_replacement_type')
BEGIN
    ALTER TABLE scheme_categories ADD ttd_tiered_income_replacement_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ttd_tiered_income_replacement_type NVARCHAR(255);
END;

-- Add or modify column for field: TtdPremiumWaiverPercentage
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ttd_premium_waiver_percentage')
BEGIN
    ALTER TABLE scheme_categories ADD ttd_premium_waiver_percentage DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ttd_premium_waiver_percentage DECIMAL(15,5);
END;

-- Add or modify column for field: TtdWaitingPeriod
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ttd_waiting_period')
BEGIN
    ALTER TABLE scheme_categories ADD ttd_waiting_period INT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ttd_waiting_period INT;
END;

-- Add or modify column for field: TtdNumberMonthlyPayments
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ttd_number_monthly_payments')
BEGIN
    ALTER TABLE scheme_categories ADD ttd_number_monthly_payments DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ttd_number_monthly_payments DECIMAL(15,5);
END;

-- Add or modify column for field: TtdDeferredPeriod
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ttd_deferred_period')
BEGIN
    ALTER TABLE scheme_categories ADD ttd_deferred_period INT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ttd_deferred_period INT;
END;

-- Add or modify column for field: TtdDisabilityDefinition
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ttd_disability_definition')
BEGIN
    ALTER TABLE scheme_categories ADD ttd_disability_definition NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ttd_disability_definition NVARCHAR(255);
END;

-- Add or modify column for field: FamilyFuneralMainMemberFuneralSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'family_funeral_main_member_funeral_sum_assured')
BEGIN
    ALTER TABLE scheme_categories ADD family_funeral_main_member_funeral_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN family_funeral_main_member_funeral_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: FamilyFuneralSpouseFuneralSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'family_funeral_spouse_funeral_sum_assured')
BEGIN
    ALTER TABLE scheme_categories ADD family_funeral_spouse_funeral_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN family_funeral_spouse_funeral_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: FamilyFuneralChildrenFuneralSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'family_funeral_children_funeral_sum_assured')
BEGIN
    ALTER TABLE scheme_categories ADD family_funeral_children_funeral_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN family_funeral_children_funeral_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: FamilyFuneralAdultDependantSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'family_funeral_adult_dependant_sum_assured')
BEGIN
    ALTER TABLE scheme_categories ADD family_funeral_adult_dependant_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN family_funeral_adult_dependant_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: FamilyFuneralParentFuneralSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'family_funeral_parent_funeral_sum_assured')
BEGIN
    ALTER TABLE scheme_categories ADD family_funeral_parent_funeral_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN family_funeral_parent_funeral_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: FamilyFuneralMaxNumberChildren
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'family_funeral_max_number_children')
BEGIN
    ALTER TABLE scheme_categories ADD family_funeral_max_number_children INT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN family_funeral_max_number_children INT;
END;

-- Add or modify column for field: FamilyFuneralMaxNumberAdultDependants
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'family_funeral_max_number_adult_dependants')
BEGIN
    ALTER TABLE scheme_categories ADD family_funeral_max_number_adult_dependants INT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN family_funeral_max_number_adult_dependants INT;
END;

-- Add or modify column for field: ExtendedFamilyBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'extended_family_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD extended_family_benefit BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN extended_family_benefit BIT;
END;

-- Add or modify column for field: ExtendedFamilyAgeBandSource
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'extended_family_age_band_source')
BEGIN
    ALTER TABLE scheme_categories ADD extended_family_age_band_source NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN extended_family_age_band_source NVARCHAR(255);
END;

-- Add or modify column for field: ExtendedFamilyAgeBandType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'extended_family_age_band_type')
BEGIN
    ALTER TABLE scheme_categories ADD extended_family_age_band_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN extended_family_age_band_type NVARCHAR(255);
END;

-- Add or modify column for field: ExtendedFamilyCustomAgeBands
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'extended_family_custom_age_bands')
BEGIN
    ALTER TABLE scheme_categories ADD extended_family_custom_age_bands text;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN extended_family_custom_age_bands text;
END;

-- Add or modify column for field: ExtendedFamilyPricingMethod
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'extended_family_pricing_method')
BEGIN
    ALTER TABLE scheme_categories ADD extended_family_pricing_method NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN extended_family_pricing_method NVARCHAR(255);
END;

-- Add or modify column for field: ExtendedFamilySumsAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'extended_family_sums_assured')
BEGIN
    ALTER TABLE scheme_categories ADD extended_family_sums_assured text;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN extended_family_sums_assured text;
END;

-- Add or modify column for field: ExtendedFamilyBandRates
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'extended_family_band_rates')
BEGIN
    ALTER TABLE scheme_categories ADD extended_family_band_rates text;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN extended_family_band_rates text;
END;

-- Add or modify column for field: PtdAlias
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ptd_alias')
BEGIN
    ALTER TABLE scheme_categories ADD ptd_alias NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ptd_alias NVARCHAR(255);
END;

-- Add or modify column for field: CiAlias
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ci_alias')
BEGIN
    ALTER TABLE scheme_categories ADD ci_alias NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ci_alias NVARCHAR(255);
END;

-- Add or modify column for field: SglaAlias
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'sgla_alias')
BEGIN
    ALTER TABLE scheme_categories ADD sgla_alias NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN sgla_alias NVARCHAR(255);
END;

-- Add or modify column for field: PhiAlias
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_alias')
BEGIN
    ALTER TABLE scheme_categories ADD phi_alias NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_alias NVARCHAR(255);
END;

-- Add or modify column for field: TtdAlias
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ttd_alias')
BEGIN
    ALTER TABLE scheme_categories ADD ttd_alias NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ttd_alias NVARCHAR(255);
END;

-- Add or modify column for field: GlaAlias
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_alias')
BEGIN
    ALTER TABLE scheme_categories ADD gla_alias NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_alias NVARCHAR(255);
END;

-- Add or modify column for field: FamilyFuneralAlias
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'family_funeral_alias')
BEGIN
    ALTER TABLE scheme_categories ADD family_funeral_alias NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN family_funeral_alias NVARCHAR(255);
END;

-- Add or modify column for field: Region
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'region')
BEGIN
    ALTER TABLE scheme_categories ADD region NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN region NVARCHAR(255);
END;

-- Add or modify column for field: GlaEducatorConversionOnWithdrawal
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_ed_conv_on_wdr')
BEGIN
    ALTER TABLE scheme_categories ADD gla_ed_conv_on_wdr BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_ed_conv_on_wdr BIT;
END;

-- Add or modify column for field: GlaEducatorConversionOnRetirement
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_ed_conv_on_ret')
BEGIN
    ALTER TABLE scheme_categories ADD gla_ed_conv_on_ret BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_ed_conv_on_ret BIT;
END;

-- Add or modify column for field: GlaEducatorContinuityDuringDisability
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_ed_cont_dur_dis')
BEGIN
    ALTER TABLE scheme_categories ADD gla_ed_cont_dur_dis BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_ed_cont_dur_dis BIT;
END;

-- Add or modify column for field: PtdEducatorConversionOnWithdrawal
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ptd_ed_conv_on_wdr')
BEGIN
    ALTER TABLE scheme_categories ADD ptd_ed_conv_on_wdr BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ptd_ed_conv_on_wdr BIT;
END;

-- Add or modify column for field: PtdEducatorConversionOnRetirement
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ptd_ed_conv_on_ret')
BEGIN
    ALTER TABLE scheme_categories ADD ptd_ed_conv_on_ret BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ptd_ed_conv_on_ret BIT;
END;

-- Add or modify column for field: PhiConversionOnWithdrawal
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_conversion_on_withdrawal')
BEGIN
    ALTER TABLE scheme_categories ADD phi_conversion_on_withdrawal BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN phi_conversion_on_withdrawal BIT;
END;

-- Add or modify column for field: SglaConversionOnWithdrawal
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'sgla_conversion_on_withdrawal')
BEGIN
    ALTER TABLE scheme_categories ADD sgla_conversion_on_withdrawal BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN sgla_conversion_on_withdrawal BIT;
END;

-- Add or modify column for field: FunConversionOnWithdrawal
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'fun_conversion_on_withdrawal')
BEGIN
    ALTER TABLE scheme_categories ADD fun_conversion_on_withdrawal BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN fun_conversion_on_withdrawal BIT;
END;

-- Add or modify column for field: TtdConversionOnWithdrawal
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ttd_conversion_on_withdrawal')
BEGIN
    ALTER TABLE scheme_categories ADD ttd_conversion_on_withdrawal BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN ttd_conversion_on_withdrawal BIT;
END;

-- Add or modify column for field: GlaContinuityDuringDisability
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_continuity_during_disability')
BEGIN
    ALTER TABLE scheme_categories ADD gla_continuity_during_disability BIT;
END;
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN gla_continuity_during_disability BIT;
END;

