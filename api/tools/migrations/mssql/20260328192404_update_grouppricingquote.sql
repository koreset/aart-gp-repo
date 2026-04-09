-- Migration for struct: GroupPricingQuote

-- Table: group_pricing_quotes

-- Ensure table exists
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'group_pricing_quotes')
BEGIN
    CREATE TABLE group_pricing_quotes (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

-- Add or modify column for field: QuoteName
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'quote_name')
BEGIN
    ALTER TABLE group_pricing_quotes ADD quote_name NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN quote_name NVARCHAR(255);
END;

-- Add or modify column for field: Basis
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'basis')
BEGIN
    ALTER TABLE group_pricing_quotes ADD basis NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN basis NVARCHAR(255);
END;

-- Add or modify column for field: CreationDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'creation_date')
BEGIN
    ALTER TABLE group_pricing_quotes ADD creation_date datetime;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN creation_date datetime;
END;

-- Add or modify column for field: QuoteType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'quote_type')
BEGIN
    ALTER TABLE group_pricing_quotes ADD quote_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN quote_type NVARCHAR(255);
END;

-- Add or modify column for field: SchemeName
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'scheme_name')
BEGIN
    ALTER TABLE group_pricing_quotes ADD scheme_name NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN scheme_name NVARCHAR(255);
END;

-- Add or modify column for field: SchemeID
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'scheme_id')
BEGIN
    ALTER TABLE group_pricing_quotes ADD scheme_id INT;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN scheme_id INT;
END;

-- Add or modify column for field: SchemeContact
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'scheme_contact')
BEGIN
    ALTER TABLE group_pricing_quotes ADD scheme_contact NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN scheme_contact NVARCHAR(255);
END;

-- Add or modify column for field: SchemeEmail
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'scheme_email')
BEGIN
    ALTER TABLE group_pricing_quotes ADD scheme_email NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN scheme_email NVARCHAR(255);
END;

-- Add or modify column for field: ID
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'broker_id')
BEGIN
    ALTER TABLE group_pricing_quotes ADD broker_id INT;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN broker_id INT;
END;

-- Add or modify column for field: Name
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'broker_name')
BEGIN
    ALTER TABLE group_pricing_quotes ADD broker_name NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN broker_name NVARCHAR(255);
END;

-- Add or modify column for field: DistributionChannel
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'distribution_channel')
BEGIN
    ALTER TABLE group_pricing_quotes ADD distribution_channel NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN distribution_channel NVARCHAR(255);
END;

-- Add or modify column for field: ObligationType
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'obligation_type')
BEGIN
    ALTER TABLE group_pricing_quotes ADD obligation_type NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN obligation_type NVARCHAR(255);
END;

-- Add or modify column for field: CommencementDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'commencement_date')
BEGIN
    ALTER TABLE group_pricing_quotes ADD commencement_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN commencement_date DATETIME2;
END;

-- Add or modify column for field: CoverEndDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'cover_end_date')
BEGIN
    ALTER TABLE group_pricing_quotes ADD cover_end_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN cover_end_date DATETIME2;
END;

-- Add or modify column for field: Industry
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'industry')
BEGIN
    ALTER TABLE group_pricing_quotes ADD industry NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN industry NVARCHAR(255);
END;

-- Add or modify column for field: OccupationClass
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'occupation_class')
BEGIN
    ALTER TABLE group_pricing_quotes ADD occupation_class INT;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN occupation_class INT;
END;

-- Add or modify column for field: FreeCoverLimit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'free_cover_limit')
BEGIN
    ALTER TABLE group_pricing_quotes ADD free_cover_limit DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN free_cover_limit DECIMAL(15,5);
END;

-- Add or modify column for field: Currency
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'currency')
BEGIN
    ALTER TABLE group_pricing_quotes ADD currency NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN currency NVARCHAR(255);
END;

-- Add or modify column for field: ExchangeRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'exchange_rate')
BEGIN
    ALTER TABLE group_pricing_quotes ADD exchange_rate INT;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN exchange_rate INT;
END;

-- Add or modify column for field: NormalRetirementAge
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'normal_retirement_age')
BEGIN
    ALTER TABLE group_pricing_quotes ADD normal_retirement_age INT;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN normal_retirement_age INT;
END;

-- Add or modify column for field: ExperienceRating
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'experience_rating')
BEGIN
    ALTER TABLE group_pricing_quotes ADD experience_rating NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN experience_rating NVARCHAR(255);
END;

-- Add or modify column for field: CreatedBy
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'created_by')
BEGIN
    ALTER TABLE group_pricing_quotes ADD created_by NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN created_by NVARCHAR(255);
END;

-- Add or modify column for field: Reviewer
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'reviewer')
BEGIN
    ALTER TABLE group_pricing_quotes ADD reviewer NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN reviewer NVARCHAR(255);
END;

-- Add or modify column for field: ApprovedBy
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'approved_by')
BEGIN
    ALTER TABLE group_pricing_quotes ADD approved_by NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN approved_by NVARCHAR(255);
END;

-- Add or modify column for field: SentBy
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'sent_by')
BEGIN
    ALTER TABLE group_pricing_quotes ADD sent_by NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN sent_by NVARCHAR(255);
END;

-- Add or modify column for field: ModifiedBy
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'modified_by')
BEGIN
    ALTER TABLE group_pricing_quotes ADD modified_by NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN modified_by NVARCHAR(255);
END;

-- Add or modify column for field: ModificationDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'modification_date')
BEGIN
    ALTER TABLE group_pricing_quotes ADD modification_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN modification_date DATETIME2;
END;

-- Add or modify column for field: Status
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'status')
BEGIN
    ALTER TABLE group_pricing_quotes ADD status NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN status NVARCHAR(255);
END;

-- Add or modify column for field: MemberDataCount
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'member_data_count')
BEGIN
    ALTER TABLE group_pricing_quotes ADD member_data_count INT;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN member_data_count INT;
END;

-- Add or modify column for field: ClaimsExperienceCount
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'claims_experience_count')
BEGIN
    ALTER TABLE group_pricing_quotes ADD claims_experience_count INT;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN claims_experience_count INT;
END;

-- Add or modify column for field: MemberRatingResultCount
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'member_rating_result_count')
BEGIN
    ALTER TABLE group_pricing_quotes ADD member_rating_result_count INT;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN member_rating_result_count INT;
END;

-- Add or modify column for field: MemberPremiumScheduleCount
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'member_premium_schedule_count')
BEGIN
    ALTER TABLE group_pricing_quotes ADD member_premium_schedule_count INT;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN member_premium_schedule_count INT;
END;

-- Add or modify column for field: BordereauxCount
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'bordereaux_count')
BEGIN
    ALTER TABLE group_pricing_quotes ADD bordereaux_count INT;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN bordereaux_count INT;
END;

-- Add or modify column for field: UseGlobalSalaryMultiple
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'use_global_salary_multiple')
BEGIN
    ALTER TABLE group_pricing_quotes ADD use_global_salary_multiple BIT;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN use_global_salary_multiple BIT;
END;

-- Add or modify column for field: ContinuationOption
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'continuation_option')
BEGIN
    ALTER TABLE group_pricing_quotes ADD continuation_option BIT;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN continuation_option BIT;
END;

-- Add or modify column for field: SelectedSchemeCategories
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'selected_scheme_categories')
BEGIN
    ALTER TABLE group_pricing_quotes ADD selected_scheme_categories json;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN selected_scheme_categories json;
END;

-- Add or modify column for field: SchemeCategories
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'scheme_categories')
BEGIN
    ALTER TABLE group_pricing_quotes ADD scheme_categories NVARCHAR(MAX);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN scheme_categories NVARCHAR(MAX);
END;

-- Add or modify column for field: CommissionLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'loadings_commission_loading')
BEGIN
    ALTER TABLE group_pricing_quotes ADD loadings_commission_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_commission_loading DECIMAL(15,5);
END;

-- Add or modify column for field: ProfitLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'loadings_profit_loading')
BEGIN
    ALTER TABLE group_pricing_quotes ADD loadings_profit_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_profit_loading DECIMAL(15,5);
END;

-- Add or modify column for field: ExpenseLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'loadings_expense_loading')
BEGIN
    ALTER TABLE group_pricing_quotes ADD loadings_expense_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_expense_loading DECIMAL(15,5);
END;

-- Add or modify column for field: AdminLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'loadings_admin_loading')
BEGIN
    ALTER TABLE group_pricing_quotes ADD loadings_admin_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_admin_loading DECIMAL(15,5);
END;

-- Add or modify column for field: ContingencyLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'loadings_contingency_loading')
BEGIN
    ALTER TABLE group_pricing_quotes ADD loadings_contingency_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_contingency_loading DECIMAL(15,5);
END;

-- Add or modify column for field: OtherLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'loadings_other_loading')
BEGIN
    ALTER TABLE group_pricing_quotes ADD loadings_other_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_other_loading DECIMAL(15,5);
END;

-- Add or modify column for field: Discount
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'loadings_discount')
BEGIN
    ALTER TABLE group_pricing_quotes ADD loadings_discount DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN loadings_discount DECIMAL(15,5);
END;

-- Add or modify column for field: MemberAverageAge
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'member_average_age')
BEGIN
    ALTER TABLE group_pricing_quotes ADD member_average_age INT;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN member_average_age INT;
END;

-- Add or modify column for field: MemberAverageIncome
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'member_average_income')
BEGIN
    ALTER TABLE group_pricing_quotes ADD member_average_income DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN member_average_income DECIMAL(15,5);
END;

-- Add or modify column for field: MemberMaleFemaleDistribution
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'member_male_female_distribution')
BEGIN
    ALTER TABLE group_pricing_quotes ADD member_male_female_distribution DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN member_male_female_distribution DECIMAL(15,5);
END;

-- Add or modify column for field: MemberIndicativeData
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'member_indicative_data')
BEGIN
    ALTER TABLE group_pricing_quotes ADD member_indicative_data BIT;
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN member_indicative_data BIT;
END;

-- Add or modify column for field: RiskRateCode
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'risk_rate_code')
BEGIN
    ALTER TABLE group_pricing_quotes ADD risk_rate_code NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN risk_rate_code NVARCHAR(255);
END;

-- Add or modify column for field: SchemeQuoteStatus
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'scheme_quote_status')
BEGIN
    ALTER TABLE group_pricing_quotes ADD scheme_quote_status NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_pricing_quotes ALTER COLUMN scheme_quote_status NVARCHAR(255);
END;

