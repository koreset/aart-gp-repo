-- Migration for struct: MemberRatingResult

-- Table: member_rating_results

-- Ensure table exists
CREATE TABLE IF NOT EXISTS member_rating_results (
    id SERIAL PRIMARY KEY
);

-- Add or modify column for field: FinancialYear
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS financial_year INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='financial_year') THEN
        ALTER TABLE member_rating_results ALTER COLUMN financial_year TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: SchemeId
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS scheme_id INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='scheme_id') THEN
        ALTER TABLE member_rating_results ALTER COLUMN scheme_id TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: QuoteId
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS quote_id INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='quote_id') THEN
        ALTER TABLE member_rating_results ALTER COLUMN quote_id TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: Category
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS category VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='category') THEN
        ALTER TABLE member_rating_results ALTER COLUMN category TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: MemberName
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS member_name VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='member_name') THEN
        ALTER TABLE member_rating_results ALTER COLUMN member_name TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: MemberCount
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS member_count INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='member_count') THEN
        ALTER TABLE member_rating_results ALTER COLUMN member_count TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: Gender
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gender VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gender') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gender TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: DateOfBirth
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS date_of_birth TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='date_of_birth') THEN
        ALTER TABLE member_rating_results ALTER COLUMN date_of_birth TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: IsOriginalMember
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS is_original_member BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='is_original_member') THEN
        ALTER TABLE member_rating_results ALTER COLUMN is_original_member TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: EntryDate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS entry_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='entry_date') THEN
        ALTER TABLE member_rating_results ALTER COLUMN entry_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: ExitDate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exit_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exit_date') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exit_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: ExpCredibility
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_credibility NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_credibility') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_credibility TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ManuallyAddedCredibility
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS manually_added_credibility NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='manually_added_credibility') THEN
        ALTER TABLE member_rating_results ALTER COLUMN manually_added_credibility TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: AnnualSalary
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS annual_salary NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='annual_salary') THEN
        ALTER TABLE member_rating_results ALTER COLUMN annual_salary TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: IncomeLevel
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS income_level INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='income_level') THEN
        ALTER TABLE member_rating_results ALTER COLUMN income_level TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: GlaSalaryMultiple
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_salary_multiple NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_salary_multiple') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_salary_multiple TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: SglaSalaryMultiple
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS sgla_salary_multiple NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='sgla_salary_multiple') THEN
        ALTER TABLE member_rating_results ALTER COLUMN sgla_salary_multiple TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdSalaryMultiple
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_salary_multiple NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ptd_salary_multiple') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ptd_salary_multiple TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiSalaryMultiple
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_salary_multiple NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ci_salary_multiple') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ci_salary_multiple TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: Occupation
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS occupation VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='occupation') THEN
        ALTER TABLE member_rating_results ALTER COLUMN occupation TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: OccupationClass
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS occupation_class INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='occupation_class') THEN
        ALTER TABLE member_rating_results ALTER COLUMN occupation_class TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: Industry
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS industry VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='industry') THEN
        ALTER TABLE member_rating_results ALTER COLUMN industry TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: AgeNextBirthday
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS age_next_birthday INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='age_next_birthday') THEN
        ALTER TABLE member_rating_results ALTER COLUMN age_next_birthday TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: AgeBand
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS age_band VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='age_band') THEN
        ALTER TABLE member_rating_results ALTER COLUMN age_band TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: SpouseGender
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_gender VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='spouse_gender') THEN
        ALTER TABLE member_rating_results ALTER COLUMN spouse_gender TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: SpouseAgeNextBirthday
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_age_next_birthday INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='spouse_age_next_birthday') THEN
        ALTER TABLE member_rating_results ALTER COLUMN spouse_age_next_birthday TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: AverageDependantAgeNextBirthday
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS average_dependant_age_next_birthday NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='average_dependant_age_next_birthday') THEN
        ALTER TABLE member_rating_results ALTER COLUMN average_dependant_age_next_birthday TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: AverageChildAgeNextBirthday
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS average_child_age_next_birthday NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='average_child_age_next_birthday') THEN
        ALTER TABLE member_rating_results ALTER COLUMN average_child_age_next_birthday TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: AverageNumberDependants
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS average_number_dependants NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='average_number_dependants') THEN
        ALTER TABLE member_rating_results ALTER COLUMN average_number_dependants TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: AverageNumberChildren
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS average_number_children NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='average_number_children') THEN
        ALTER TABLE member_rating_results ALTER COLUMN average_number_children TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CalculatedFreeCoverLimit
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS calculated_free_cover_limit NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='calculated_free_cover_limit') THEN
        ALTER TABLE member_rating_results ALTER COLUMN calculated_free_cover_limit TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: AppliedFreeCoverLimit
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS applied_free_cover_limit NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='applied_free_cover_limit') THEN
        ALTER TABLE member_rating_results ALTER COLUMN applied_free_cover_limit TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaSumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaCappedSumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_capped_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_capped_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_capped_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpenseLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS expense_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='expense_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN expense_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: AdminLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS admin_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='admin_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN admin_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CommissionLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS commission_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='commission_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN commission_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ProfitLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS profit_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='profit_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN profit_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: OtherLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS other_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='other_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN other_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: Discount
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS discount NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='discount') THEN
        ALTER TABLE member_rating_results ALTER COLUMN discount TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TotalPremiumLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS total_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='total_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN total_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaRegionLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_region_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_region_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_region_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaAidsRegionLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_aids_region_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_aids_region_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_aids_region_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdRegionLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_region_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ptd_region_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ptd_region_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiRegionLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_region_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ci_region_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ci_region_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdRegionLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_region_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ttd_region_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ttd_region_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiRegionLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_region_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='phi_region_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN phi_region_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: FunRegionLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS fun_region_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='fun_region_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN fun_region_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: FunAidsRegionLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS fun_aids_region_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='fun_aids_region_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN fun_aids_region_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaIndustryLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_industry_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_industry_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_industry_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdIndustryLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_industry_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ptd_industry_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ptd_industry_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiIndustryLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_industry_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ci_industry_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ci_industry_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdIndustryLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_industry_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ttd_industry_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ttd_industry_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiIndustryLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_industry_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='phi_industry_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN phi_industry_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaContingencyLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_contingency_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_contingency_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_contingency_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdContingencyLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_contingency_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ptd_contingency_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ptd_contingency_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiContingencyLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_contingency_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ci_contingency_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ci_contingency_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdContingencyLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_contingency_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ttd_contingency_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ttd_contingency_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiContingencyLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_contingency_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='phi_contingency_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN phi_contingency_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: FunContingencyLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS fun_contingency_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='fun_contingency_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN fun_contingency_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ContinuationLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS continuation_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='continuation_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN continuation_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaQx
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_qx NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_qx') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_qx TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaAidsQx
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_aids_qx NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_aids_qx') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_aids_qx TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: BaseGlaRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS base_gla_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='base_gla_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN base_gla_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaTerminalIllnessLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_terminal_illness_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_terminal_illness_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_terminal_illness_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: LoadedGlaRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS loaded_gla_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='loaded_gla_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN loaded_gla_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaWeightedExperienceCrudeRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_weighted_experience_crude_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_weighted_experience_crude_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_weighted_experience_crude_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaTheoreticalRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_theoretical_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_theoretical_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_theoretical_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdExperienceCrudeRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_experience_crude_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ptd_experience_crude_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ptd_experience_crude_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdTheoreticalRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_theoretical_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ptd_theoretical_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ptd_theoretical_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiExperienceCrudeRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_experience_crude_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ci_experience_crude_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ci_experience_crude_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiTheoreticalRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_theoretical_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ci_theoretical_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ci_theoretical_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjLoadedGlaRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_loaded_gla_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_loaded_gla_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_loaded_gla_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaExperienceAdjustment
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_experience_adjustment NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_experience_adjustment') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_experience_adjustment TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_risk_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_risk_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_risk_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjGlaRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_risk_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_gla_risk_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_gla_risk_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='gla_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN gla_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjGlaOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_gla_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_gla_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdSumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ptd_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ptd_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdCappedSumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_capped_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ptd_capped_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ptd_capped_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: BasePtdRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS base_ptd_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='base_ptd_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN base_ptd_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ptd_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ptd_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: LoadedPtdRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS loaded_ptd_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='loaded_ptd_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN loaded_ptd_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdExperienceAdjustment
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_experience_adjustment NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ptd_experience_adjustment') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ptd_experience_adjustment TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjLoadedPtdRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_loaded_ptd_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_loaded_ptd_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_loaded_ptd_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_risk_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ptd_risk_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ptd_risk_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjPtdRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ptd_risk_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_ptd_risk_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_ptd_risk_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ptd_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ptd_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjPtdOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ptd_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_ptd_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_ptd_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiSumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ci_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ci_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiCappedSumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_capped_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ci_capped_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ci_capped_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: BaseCiRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS base_ci_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='base_ci_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN base_ci_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ci_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ci_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: LoadedCiRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS loaded_ci_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='loaded_ci_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN loaded_ci_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiExperienceAdjustment
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_experience_adjustment NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ci_experience_adjustment') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ci_experience_adjustment TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjLoadedCiRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_loaded_ci_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_loaded_ci_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_loaded_ci_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_risk_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ci_risk_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ci_risk_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjCiRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ci_risk_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_ci_risk_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_ci_risk_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ci_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ci_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjCiOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ci_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_ci_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_ci_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: SpouseGlaSumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_gla_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='spouse_gla_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN spouse_gla_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: SpouseGlaCappedSumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_gla_capped_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='spouse_gla_capped_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN spouse_gla_capped_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: SpouseGlaQx
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_gla_qx NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='spouse_gla_qx') THEN
        ALTER TABLE member_rating_results ALTER COLUMN spouse_gla_qx TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: SpouseGlaAidsQx
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_gla_aids_qx NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='spouse_gla_aids_qx') THEN
        ALTER TABLE member_rating_results ALTER COLUMN spouse_gla_aids_qx TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: BaseSpouseGlaRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS base_spouse_gla_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='base_spouse_gla_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN base_spouse_gla_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: SpouseGlaLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_gla_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='spouse_gla_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN spouse_gla_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: LoadedSpouseGlaRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS loaded_spouse_gla_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='loaded_spouse_gla_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN loaded_spouse_gla_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjLoadedSpouseGlaRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_loaded_spouse_gla_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_loaded_spouse_gla_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_loaded_spouse_gla_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: SpouseGlaRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_gla_risk_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='spouse_gla_risk_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN spouse_gla_risk_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: SpouseGlaOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_gla_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='spouse_gla_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN spouse_gla_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjSpouseGlaOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_spouse_gla_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_spouse_gla_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_spouse_gla_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdIncome
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_income NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ttd_income') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ttd_income TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdCappedIncome
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_capped_income NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ttd_capped_income') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ttd_capped_income TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdNumberOfMonthlyPayments
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_number_of_monthly_payments NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ttd_number_of_monthly_payments') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ttd_number_of_monthly_payments TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: IncomeReplacementRatio
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS income_replacement_ratio NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='income_replacement_ratio') THEN
        ALTER TABLE member_rating_results ALTER COLUMN income_replacement_ratio TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: BaseTtdRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS base_ttd_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='base_ttd_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN base_ttd_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ttd_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ttd_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: LoadedTtdRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS loaded_ttd_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='loaded_ttd_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN loaded_ttd_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdExperienceAdjustment
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_experience_adjustment NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ttd_experience_adjustment') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ttd_experience_adjustment TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjLoadedTtdRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_loaded_ttd_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_loaded_ttd_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_loaded_ttd_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_risk_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ttd_risk_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ttd_risk_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjTtdRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ttd_risk_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_ttd_risk_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_ttd_risk_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='ttd_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN ttd_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjTtdOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ttd_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_ttd_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_ttd_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjSpouseGlaRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_spouse_gla_risk_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_spouse_gla_risk_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_spouse_gla_risk_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiIncome
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_income NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='phi_income') THEN
        ALTER TABLE member_rating_results ALTER COLUMN phi_income TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiCappedIncome
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_capped_income NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='phi_capped_income') THEN
        ALTER TABLE member_rating_results ALTER COLUMN phi_capped_income TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiContributionWaiver
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_contribution_waiver NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='phi_contribution_waiver') THEN
        ALTER TABLE member_rating_results ALTER COLUMN phi_contribution_waiver TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiMedicalAidWaiver
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_medical_aid_waiver NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='phi_medical_aid_waiver') THEN
        ALTER TABLE member_rating_results ALTER COLUMN phi_medical_aid_waiver TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiMonthlyBenefit
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_monthly_benefit NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='phi_monthly_benefit') THEN
        ALTER TABLE member_rating_results ALTER COLUMN phi_monthly_benefit TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiAnnuityFactor
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_annuity_factor NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='phi_annuity_factor') THEN
        ALTER TABLE member_rating_results ALTER COLUMN phi_annuity_factor TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: BasePhiRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS base_phi_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='base_phi_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN base_phi_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiSalaryLevel
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_salary_level NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='phi_salary_level') THEN
        ALTER TABLE member_rating_results ALTER COLUMN phi_salary_level TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiLoading
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='phi_loading') THEN
        ALTER TABLE member_rating_results ALTER COLUMN phi_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: LoadedPhiRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS loaded_phi_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='loaded_phi_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN loaded_phi_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiExperienceAdjustment
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_experience_adjustment NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='phi_experience_adjustment') THEN
        ALTER TABLE member_rating_results ALTER COLUMN phi_experience_adjustment TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjLoadedPhiRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_loaded_phi_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_loaded_phi_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_loaded_phi_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_risk_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='phi_risk_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN phi_risk_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjPhiRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_phi_risk_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_phi_risk_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_phi_risk_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='phi_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN phi_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjPhiOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_phi_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_phi_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_phi_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: MemberFuneralSumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS member_funeral_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='member_funeral_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN member_funeral_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: MainMemberFuneralBaseRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS main_member_funeral_base_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='main_member_funeral_base_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN main_member_funeral_base_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: MainMemberFuneralRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS main_member_funeral_cost NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='main_member_funeral_cost') THEN
        ALTER TABLE member_rating_results ALTER COLUMN main_member_funeral_cost TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: MainMemberFuneralOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS main_member_funeral_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='main_member_funeral_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN main_member_funeral_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: SpouseFuneralSumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_funeral_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='spouse_funeral_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN spouse_funeral_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: SpouseFuneralBaseRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_funeral_base_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='spouse_funeral_base_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN spouse_funeral_base_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: SpouseFuneralRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_funeral_cost NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='spouse_funeral_cost') THEN
        ALTER TABLE member_rating_results ALTER COLUMN spouse_funeral_cost TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: SpouseFuneralOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_funeral_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='spouse_funeral_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN spouse_funeral_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ChildFuneralBaseRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS child_funeral_base_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='child_funeral_base_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN child_funeral_base_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ChildFuneralSumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS child_funeral_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='child_funeral_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN child_funeral_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ChildFuneralRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS children_funeral_cost NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='children_funeral_cost') THEN
        ALTER TABLE member_rating_results ALTER COLUMN children_funeral_cost TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ChildrenFuneralOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS children_funeral_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='children_funeral_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN children_funeral_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ParentFuneralBaseRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS dependant_funeral_base_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='dependant_funeral_base_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN dependant_funeral_base_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ParentFuneralSumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS dependant_funeral_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='dependant_funeral_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN dependant_funeral_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ParentFuneralRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS dependants_funeral_cost NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='dependants_funeral_cost') THEN
        ALTER TABLE member_rating_results ALTER COLUMN dependants_funeral_cost TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ParentFuneralOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS dependants_funeral_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='dependants_funeral_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN dependants_funeral_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ParentFuneralSumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS parent_funeral_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='parent_funeral_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN parent_funeral_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TotalFuneralRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS total_funeral_risk_cost NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='total_funeral_risk_cost') THEN
        ALTER TABLE member_rating_results ALTER COLUMN total_funeral_risk_cost TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjTotalFuneralRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_total_funeral_risk_cost NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_total_funeral_risk_cost') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_total_funeral_risk_cost TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TotalFuneralOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS total_funeral_office_cost NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='total_funeral_office_cost') THEN
        ALTER TABLE member_rating_results ALTER COLUMN total_funeral_office_cost TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjTotalFuneralOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_total_funeral_office_cost NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_total_funeral_office_cost') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_total_funeral_office_cost TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: Grade0SumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS grade0_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='grade0_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN grade0_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: Grade17SumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS grade17_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='grade17_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN grade17_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: Grade812SumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS grade812_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='grade812_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN grade812_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TertiarySumAssured
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS tertiary_sum_assured NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='tertiary_sum_assured') THEN
        ALTER TABLE member_rating_results ALTER COLUMN tertiary_sum_assured TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: Grade0RiskRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS grade0_risk_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='grade0_risk_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN grade0_risk_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: Grade17RiskRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS grade17_risk_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='grade17_risk_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN grade17_risk_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: Grade812RiskRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS grade812_risk_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='grade812_risk_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN grade812_risk_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TertiaryRiskRate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS tertiary_risk_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='tertiary_risk_rate') THEN
        ALTER TABLE member_rating_results ALTER COLUMN tertiary_risk_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: EducatorRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS educator_risk_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='educator_risk_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN educator_risk_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: EducatorOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS educator_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='educator_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN educator_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjEducatorRiskPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_educator_risk_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_educator_risk_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_educator_risk_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpAdjEducatorOfficePremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_educator_office_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exp_adj_educator_office_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exp_adj_educator_office_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExceedsNormalRetirementAgeIndicator
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exceeds_normal_retirement_age_indicator INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exceeds_normal_retirement_age_indicator') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exceeds_normal_retirement_age_indicator TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: ExceedsFreeCoverLimitIndicator
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exceeds_free_cover_limit_indicator INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='exceeds_free_cover_limit_indicator') THEN
        ALTER TABLE member_rating_results ALTER COLUMN exceeds_free_cover_limit_indicator TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: FuneralExperienceAdjustedAnnualPremium
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS funeral_experience_adjusted_annual_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='funeral_experience_adjusted_annual_premium') THEN
        ALTER TABLE member_rating_results ALTER COLUMN funeral_experience_adjusted_annual_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CreationDate
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='creation_date') THEN
        ALTER TABLE member_rating_results ALTER COLUMN creation_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: CreatedBy
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='member_rating_results' AND column_name='created_by') THEN
        ALTER TABLE member_rating_results ALTER COLUMN created_by TYPE VARCHAR(255);
    END IF;
END $$;

