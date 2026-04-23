-- Migration for struct: SchemeCategory

-- Table: scheme_categories

-- Ensure table exists
CREATE TABLE IF NOT EXISTS scheme_categories (
    id SERIAL PRIMARY KEY
);

-- Add or modify column for field: QuoteId
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS quote_id INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='quote_id') THEN
        ALTER TABLE scheme_categories ALTER COLUMN quote_id TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: SchemeCategory
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS scheme_category VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='scheme_category') THEN
        ALTER TABLE scheme_categories ALTER COLUMN scheme_category TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Basis
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS basis VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='basis') THEN
        ALTER TABLE scheme_categories ALTER COLUMN basis TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: FreeCoverLimit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS free_cover_limit NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='free_cover_limit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN free_cover_limit TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_benefit BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ptd_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ptd_benefit TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: GlaBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_benefit BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_benefit TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: CiBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ci_benefit BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ci_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ci_benefit TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: SglaBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS sgla_benefit BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='sgla_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN sgla_benefit TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: PhiBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_benefit BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_benefit TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: TtdBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ttd_benefit BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ttd_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ttd_benefit TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: FamilyFuneralBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS family_funeral_benefit BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='family_funeral_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN family_funeral_benefit TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: GlaSalaryMultiple
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_salary_multiple NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_salary_multiple') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_salary_multiple TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaTerminalIllnessBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_terminal_illness_benefit VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_terminal_illness_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_terminal_illness_benefit TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: GlaWaitingPeriod
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_waiting_period INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_waiting_period') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_waiting_period TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: GlaEducatorBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_educator_benefit VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_educator_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_educator_benefit TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: GlaEducatorBenefitType
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_educator_benefit_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_educator_benefit_type') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_educator_benefit_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: GlaBenefitType
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_benefit_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_benefit_type') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_benefit_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: GlaConversionOnWithdrawal
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_conversion_on_withdrawal BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_conversion_on_withdrawal') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_conversion_on_withdrawal TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: GlaConversionOnRetirement
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_conversion_on_retirement BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_conversion_on_retirement') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_conversion_on_retirement TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: AdditionalAccidentalGlaBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_accidental_gla_benefit BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='additional_accidental_gla_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN additional_accidental_gla_benefit TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: AdditionalAccidentalGlaBenefitType
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_accidental_gla_benefit_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='additional_accidental_gla_benefit_type') THEN
        ALTER TABLE scheme_categories ALTER COLUMN additional_accidental_gla_benefit_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: TaxSaverBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS tax_saver_benefit BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='tax_saver_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN tax_saver_benefit TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: AdditionalGlaCoverBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_gla_cover_benefit BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='additional_gla_cover_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_benefit TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: AdditionalGlaCoverAgeBandSource
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_gla_cover_age_band_source VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='additional_gla_cover_age_band_source') THEN
        ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_age_band_source TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: AdditionalGlaCoverAgeBandType
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_gla_cover_age_band_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='additional_gla_cover_age_band_type') THEN
        ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_age_band_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: AdditionalGlaCoverCustomAgeBands
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_gla_cover_custom_age_bands text;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='additional_gla_cover_custom_age_bands') THEN
        ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_custom_age_bands TYPE text;
    END IF;
END $$;

-- Add or modify column for field: AdditionalGlaCoverBandRates
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_gla_cover_band_rates text;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='additional_gla_cover_band_rates') THEN
        ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_band_rates TYPE text;
    END IF;
END $$;

-- Add or modify column for field: AdditionalGlaCoverMalePropUsed
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_gla_cover_male_prop_used TEXT;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='additional_gla_cover_male_prop_used') THEN
        ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_male_prop_used TYPE TEXT;
    END IF;
END $$;

-- Add or modify column for field: PtdRiskType
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_risk_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ptd_risk_type') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ptd_risk_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: PtdBenefitType
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_benefit_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ptd_benefit_type') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ptd_benefit_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: PtdSalaryMultiple
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_salary_multiple NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ptd_salary_multiple') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ptd_salary_multiple TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdDeferredPeriod
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_deferred_period INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ptd_deferred_period') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ptd_deferred_period TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: PtdDisabilityDefinition
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_disability_definition VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ptd_disability_definition') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ptd_disability_definition TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: PtdEducatorBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_educator_benefit VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ptd_educator_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ptd_educator_benefit TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: PtdEducatorBenefitType
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_educator_benefit_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ptd_educator_benefit_type') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ptd_educator_benefit_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: PtdConversionOnWithdrawal
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_conversion_on_withdrawal BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ptd_conversion_on_withdrawal') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ptd_conversion_on_withdrawal TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: CiBenefitStructure
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ci_benefit_structure VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ci_benefit_structure') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ci_benefit_structure TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: CiBenefitDefinition
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ci_benefit_definition VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ci_benefit_definition') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ci_benefit_definition TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: CiMaxBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ci_max_benefit NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ci_max_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ci_max_benefit TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiCriticalIllnessSalaryMultiple
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ci_critical_illness_salary_multiple NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ci_critical_illness_salary_multiple') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ci_critical_illness_salary_multiple TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiConversionOnWithdrawal
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ci_conversion_on_withdrawal BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ci_conversion_on_withdrawal') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ci_conversion_on_withdrawal TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: SglaSalaryMultiple
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS sgla_salary_multiple NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='sgla_salary_multiple') THEN
        ALTER TABLE scheme_categories ALTER COLUMN sgla_salary_multiple TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: SglaMaxBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS sgla_max_benefit NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='sgla_max_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN sgla_max_benefit TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiRiskType
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_risk_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_risk_type') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_risk_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: PhiMaximumBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_maximum_benefit NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_maximum_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_maximum_benefit TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiIncomeReplacementPercentage
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_income_replacement_percentage NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_income_replacement_percentage') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_income_replacement_percentage TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiUseTieredIncomeReplacementRatio
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_use_tiered_income_replacement_ratio BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_use_tiered_income_replacement_ratio') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_use_tiered_income_replacement_ratio TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: PhiTieredIncomeReplacementType
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_tiered_income_replacement_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_tiered_income_replacement_type') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_tiered_income_replacement_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: PhiPremiumWaiver
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_premium_waiver VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_premium_waiver') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_premium_waiver TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: PhiMedicalAidPremiumWaiver
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_medical_aid_premium_waiver VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_medical_aid_premium_waiver') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_medical_aid_premium_waiver TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: PhiBenefitEscalation
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_benefit_escalation VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_benefit_escalation') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_benefit_escalation TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: PhiMaxPremiumWaiver
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_max_premium_waiver NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_max_premium_waiver') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_max_premium_waiver TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiWaitingPeriod
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_waiting_period INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_waiting_period') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_waiting_period TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: PhiNumberMonthlyPayments
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_number_monthly_payments INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_number_monthly_payments') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_number_monthly_payments TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: PhiDeferredPeriod
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_deferred_period INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_deferred_period') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_deferred_period TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: PhiDisabilityDefinition
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_disability_definition VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_disability_definition') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_disability_definition TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: PhiNormalRetirementAge
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_normal_retirement_age INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_normal_retirement_age') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_normal_retirement_age TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: TtdRiskType
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ttd_risk_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ttd_risk_type') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ttd_risk_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: TtdMaximumBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ttd_maximum_benefit NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ttd_maximum_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ttd_maximum_benefit TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdIncomeReplacementPercentage
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ttd_income_replacement_percentage NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ttd_income_replacement_percentage') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ttd_income_replacement_percentage TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdUseTieredIncomeReplacementRatio
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ttd_use_tiered_income_replacement_ratio BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ttd_use_tiered_income_replacement_ratio') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ttd_use_tiered_income_replacement_ratio TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: TtdTieredIncomeReplacementType
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ttd_tiered_income_replacement_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ttd_tiered_income_replacement_type') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ttd_tiered_income_replacement_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: TtdPremiumWaiverPercentage
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ttd_premium_waiver_percentage NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ttd_premium_waiver_percentage') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ttd_premium_waiver_percentage TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdWaitingPeriod
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ttd_waiting_period INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ttd_waiting_period') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ttd_waiting_period TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: TtdNumberMonthlyPayments
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ttd_number_monthly_payments NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ttd_number_monthly_payments') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ttd_number_monthly_payments TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdDeferredPeriod
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ttd_deferred_period INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ttd_deferred_period') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ttd_deferred_period TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: TtdDisabilityDefinition
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ttd_disability_definition VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ttd_disability_definition') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ttd_disability_definition TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: FamilyFuneralMainMemberFuneralSumAssured
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS family_funeral_main_member_funeral_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='family_funeral_main_member_funeral_sum_assured') THEN
        ALTER TABLE scheme_categories ALTER COLUMN family_funeral_main_member_funeral_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: FamilyFuneralSpouseFuneralSumAssured
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS family_funeral_spouse_funeral_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='family_funeral_spouse_funeral_sum_assured') THEN
        ALTER TABLE scheme_categories ALTER COLUMN family_funeral_spouse_funeral_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: FamilyFuneralChildrenFuneralSumAssured
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS family_funeral_children_funeral_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='family_funeral_children_funeral_sum_assured') THEN
        ALTER TABLE scheme_categories ALTER COLUMN family_funeral_children_funeral_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: FamilyFuneralAdultDependantSumAssured
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS family_funeral_adult_dependant_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='family_funeral_adult_dependant_sum_assured') THEN
        ALTER TABLE scheme_categories ALTER COLUMN family_funeral_adult_dependant_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: FamilyFuneralParentFuneralSumAssured
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS family_funeral_parent_funeral_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='family_funeral_parent_funeral_sum_assured') THEN
        ALTER TABLE scheme_categories ALTER COLUMN family_funeral_parent_funeral_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: FamilyFuneralMaxNumberChildren
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS family_funeral_max_number_children INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='family_funeral_max_number_children') THEN
        ALTER TABLE scheme_categories ALTER COLUMN family_funeral_max_number_children TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: FamilyFuneralMaxNumberAdultDependants
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS family_funeral_max_number_adult_dependants INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='family_funeral_max_number_adult_dependants') THEN
        ALTER TABLE scheme_categories ALTER COLUMN family_funeral_max_number_adult_dependants TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: ExtendedFamilyBenefit
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS extended_family_benefit BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='extended_family_benefit') THEN
        ALTER TABLE scheme_categories ALTER COLUMN extended_family_benefit TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: ExtendedFamilyAgeBandSource
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS extended_family_age_band_source VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='extended_family_age_band_source') THEN
        ALTER TABLE scheme_categories ALTER COLUMN extended_family_age_band_source TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ExtendedFamilyAgeBandType
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS extended_family_age_band_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='extended_family_age_band_type') THEN
        ALTER TABLE scheme_categories ALTER COLUMN extended_family_age_band_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ExtendedFamilyCustomAgeBands
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS extended_family_custom_age_bands text;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='extended_family_custom_age_bands') THEN
        ALTER TABLE scheme_categories ALTER COLUMN extended_family_custom_age_bands TYPE text;
    END IF;
END $$;

-- Add or modify column for field: ExtendedFamilyPricingMethod
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS extended_family_pricing_method VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='extended_family_pricing_method') THEN
        ALTER TABLE scheme_categories ALTER COLUMN extended_family_pricing_method TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ExtendedFamilySumsAssured
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS extended_family_sums_assured text;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='extended_family_sums_assured') THEN
        ALTER TABLE scheme_categories ALTER COLUMN extended_family_sums_assured TYPE text;
    END IF;
END $$;

-- Add or modify column for field: ExtendedFamilyBandRates
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS extended_family_band_rates text;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='extended_family_band_rates') THEN
        ALTER TABLE scheme_categories ALTER COLUMN extended_family_band_rates TYPE text;
    END IF;
END $$;

-- Add or modify column for field: PtdAlias
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_alias VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ptd_alias') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ptd_alias TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: CiAlias
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ci_alias VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ci_alias') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ci_alias TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: SglaAlias
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS sgla_alias VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='sgla_alias') THEN
        ALTER TABLE scheme_categories ALTER COLUMN sgla_alias TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: PhiAlias
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_alias VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_alias') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_alias TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: TtdAlias
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ttd_alias VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ttd_alias') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ttd_alias TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: GlaAlias
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_alias VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_alias') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_alias TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: FamilyFuneralAlias
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS family_funeral_alias VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='family_funeral_alias') THEN
        ALTER TABLE scheme_categories ALTER COLUMN family_funeral_alias TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Region
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS region VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='region') THEN
        ALTER TABLE scheme_categories ALTER COLUMN region TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: GlaEducatorConversionOnWithdrawal
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_ed_conv_on_wdr BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_ed_conv_on_wdr') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_ed_conv_on_wdr TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: GlaEducatorConversionOnRetirement
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_ed_conv_on_ret BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_ed_conv_on_ret') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_ed_conv_on_ret TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: GlaEducatorContinuityDuringDisability
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_ed_cont_dur_dis BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_ed_cont_dur_dis') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_ed_cont_dur_dis TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: PtdEducatorConversionOnWithdrawal
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_wdr BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ptd_ed_conv_on_wdr') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ptd_ed_conv_on_wdr TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: PtdEducatorConversionOnRetirement
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_ret BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ptd_ed_conv_on_ret') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ptd_ed_conv_on_ret TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: PhiConversionOnWithdrawal
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_conversion_on_withdrawal BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='phi_conversion_on_withdrawal') THEN
        ALTER TABLE scheme_categories ALTER COLUMN phi_conversion_on_withdrawal TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: SglaConversionOnWithdrawal
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS sgla_conversion_on_withdrawal BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='sgla_conversion_on_withdrawal') THEN
        ALTER TABLE scheme_categories ALTER COLUMN sgla_conversion_on_withdrawal TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: FunConversionOnWithdrawal
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS fun_conversion_on_withdrawal BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='fun_conversion_on_withdrawal') THEN
        ALTER TABLE scheme_categories ALTER COLUMN fun_conversion_on_withdrawal TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: TtdConversionOnWithdrawal
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ttd_conversion_on_withdrawal BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='ttd_conversion_on_withdrawal') THEN
        ALTER TABLE scheme_categories ALTER COLUMN ttd_conversion_on_withdrawal TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: GlaContinuityDuringDisability
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_continuity_during_disability BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_categories' AND column_name='gla_continuity_during_disability') THEN
        ALTER TABLE scheme_categories ALTER COLUMN gla_continuity_during_disability TYPE BOOLEAN;
    END IF;
END $$;

