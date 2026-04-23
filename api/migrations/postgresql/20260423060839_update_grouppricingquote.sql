-- Migration for struct: GroupPricingQuote

-- Table: group_pricing_quotes

-- Ensure table exists
CREATE TABLE IF NOT EXISTS group_pricing_quotes (
    id SERIAL PRIMARY KEY
);

-- Add or modify column for field: QuoteName
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS quote_name VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='quote_name') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN quote_name TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Basis
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS basis VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='basis') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN basis TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: CreationDate
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS creation_date datetime;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='creation_date') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN creation_date TYPE datetime;
    END IF;
END $$;

-- Add or modify column for field: QuoteType
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS quote_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='quote_type') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN quote_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: SchemeName
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS scheme_name VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='scheme_name') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN scheme_name TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: SchemeID
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS scheme_id INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='scheme_id') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN scheme_id TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: SchemeContact
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS scheme_contact VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='scheme_contact') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN scheme_contact TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: SchemeEmail
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS scheme_email VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='scheme_email') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN scheme_email TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ID
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS broker_id INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='broker_id') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN broker_id TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: Name
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS broker_name VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='broker_name') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN broker_name TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: DistributionChannel
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS distribution_channel VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='distribution_channel') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN distribution_channel TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ObligationType
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS obligation_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='obligation_type') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN obligation_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: CommencementDate
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS commencement_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='commencement_date') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN commencement_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: CoverEndDate
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS cover_end_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='cover_end_date') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN cover_end_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: Industry
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS industry VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='industry') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN industry TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: OccupationClass
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS occupation_class INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='occupation_class') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN occupation_class TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: FreeCoverLimit
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS free_cover_limit NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='free_cover_limit') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN free_cover_limit TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: Currency
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS currency VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='currency') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN currency TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ExchangeRate
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS exchange_rate INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='exchange_rate') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN exchange_rate TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: NormalRetirementAge
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS normal_retirement_age INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='normal_retirement_age') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN normal_retirement_age TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: ExperienceRating
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS experience_rating VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='experience_rating') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN experience_rating TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: CreatedBy
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='created_by') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN created_by TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Reviewer
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS reviewer VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='reviewer') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN reviewer TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ApprovedBy
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS approved_by VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='approved_by') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN approved_by TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: SentBy
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS sent_by VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='sent_by') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN sent_by TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ModifiedBy
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS modified_by VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='modified_by') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN modified_by TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ModificationDate
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS modification_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='modification_date') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN modification_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: Status
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS status VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='status') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN status TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: MemberDataCount
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS member_data_count INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='member_data_count') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN member_data_count TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: ClaimsExperienceCount
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS claims_experience_count INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='claims_experience_count') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN claims_experience_count TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: MemberRatingResultCount
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS member_rating_result_count INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='member_rating_result_count') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN member_rating_result_count TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: MemberPremiumScheduleCount
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS member_premium_schedule_count INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='member_premium_schedule_count') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN member_premium_schedule_count TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: BordereauxCount
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS bordereaux_count INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='bordereaux_count') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN bordereaux_count TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: UseGlobalSalaryMultiple
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS use_global_salary_multiple BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='use_global_salary_multiple') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN use_global_salary_multiple TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: SelectedSchemeCategories
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS selected_scheme_categories json;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='selected_scheme_categories') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN selected_scheme_categories TYPE json;
    END IF;
END $$;

-- Add or modify column for field: SchemeCategories
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS scheme_categories TEXT;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='scheme_categories') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN scheme_categories TYPE TEXT;
    END IF;
END $$;

-- Add or modify column for field: CommissionLoading
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS loadings_commission_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='loadings_commission_loading') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_commission_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ProfitLoading
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS loadings_profit_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='loadings_profit_loading') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_profit_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpenseLoading
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS loadings_expense_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='loadings_expense_loading') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_expense_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: AdminLoading
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS loadings_admin_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='loadings_admin_loading') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_admin_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ContingencyLoading
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS loadings_contingency_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='loadings_contingency_loading') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_contingency_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: OtherLoading
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS loadings_other_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='loadings_other_loading') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_other_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: Discount
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS loadings_discount NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='loadings_discount') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_discount TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: BinderFee
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS loadings_binder_fee NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='loadings_binder_fee') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_binder_fee TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: OutsourceFee
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS loadings_outsource_fee NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='loadings_outsource_fee') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_outsource_fee TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: MemberAverageAge
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS member_average_age INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='member_average_age') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN member_average_age TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: MemberAverageIncome
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS member_average_income NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='member_average_income') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN member_average_income TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: MemberMaleFemaleDistribution
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS member_male_female_distribution NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='member_male_female_distribution') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN member_male_female_distribution TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: MemberIndicativeData
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS member_indicative_data BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='member_indicative_data') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN member_indicative_data TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: RiskRateCode
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS risk_rate_code VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='risk_rate_code') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN risk_rate_code TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: SchemeQuoteStatus
ALTER TABLE group_pricing_quotes ADD COLUMN IF NOT EXISTS scheme_quote_status VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_pricing_quotes' AND column_name='scheme_quote_status') THEN
        ALTER TABLE group_pricing_quotes ALTER COLUMN scheme_quote_status TYPE VARCHAR(255);
    END IF;
END $$;

