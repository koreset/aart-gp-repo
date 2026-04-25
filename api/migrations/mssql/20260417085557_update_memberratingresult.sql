-- Migration for struct: MemberRatingResult

-- Table: member_rating_results

-- Ensure table exists
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_results')
BEGIN
    CREATE TABLE member_rating_results (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

-- Add or modify column for field: FinancialYear
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'financial_year')
BEGIN
    ALTER TABLE member_rating_results ADD financial_year INT;
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN financial_year INT;
END;

-- Add or modify column for field: SchemeId
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'scheme_id')
BEGIN
    ALTER TABLE member_rating_results ADD scheme_id INT;
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN scheme_id INT;
END;

-- Add or modify column for field: QuoteId
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'quote_id')
BEGIN
    ALTER TABLE member_rating_results ADD quote_id INT;
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN quote_id INT;
END;

-- Add or modify column for field: Category
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'category')
BEGIN
    ALTER TABLE member_rating_results ADD category NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN category NVARCHAR(255);
END;

-- Add or modify column for field: MemberName
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'member_name')
BEGIN
    ALTER TABLE member_rating_results ADD member_name NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN member_name NVARCHAR(255);
END;

-- Add or modify column for field: MemberCount
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'member_count')
BEGIN
    ALTER TABLE member_rating_results ADD member_count INT;
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN member_count INT;
END;

-- Add or modify column for field: Gender
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gender')
BEGIN
    ALTER TABLE member_rating_results ADD gender NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gender NVARCHAR(255);
END;

-- Add or modify column for field: DateOfBirth
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'date_of_birth')
BEGIN
    ALTER TABLE member_rating_results ADD date_of_birth DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN date_of_birth DATETIME2;
END;

-- Add or modify column for field: IsOriginalMember
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'is_original_member')
BEGIN
    ALTER TABLE member_rating_results ADD is_original_member BIT;
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN is_original_member BIT;
END;

-- Add or modify column for field: EntryDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'entry_date')
BEGIN
    ALTER TABLE member_rating_results ADD entry_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN entry_date DATETIME2;
END;

-- Add or modify column for field: ExitDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exit_date')
BEGIN
    ALTER TABLE member_rating_results ADD exit_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exit_date DATETIME2;
END;

-- Add or modify column for field: ExpCredibility
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_credibility')
BEGIN
    ALTER TABLE member_rating_results ADD exp_credibility DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_credibility DECIMAL(15,5);
END;

-- Add or modify column for field: ManuallyAddedCredibility
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'manually_added_credibility')
BEGIN
    ALTER TABLE member_rating_results ADD manually_added_credibility DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN manually_added_credibility DECIMAL(15,5);
END;

-- Add or modify column for field: AnnualSalary
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'annual_salary')
BEGIN
    ALTER TABLE member_rating_results ADD annual_salary DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN annual_salary DECIMAL(15,5);
END;

-- Add or modify column for field: IncomeLevel
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'income_level')
BEGIN
    ALTER TABLE member_rating_results ADD income_level INT;
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN income_level INT;
END;

-- Add or modify column for field: GlaSalaryMultiple
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_salary_multiple')
BEGIN
    ALTER TABLE member_rating_results ADD gla_salary_multiple DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_salary_multiple DECIMAL(15,5);
END;

-- Add or modify column for field: SglaSalaryMultiple
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'sgla_salary_multiple')
BEGIN
    ALTER TABLE member_rating_results ADD sgla_salary_multiple DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN sgla_salary_multiple DECIMAL(15,5);
END;

-- Add or modify column for field: PtdSalaryMultiple
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_salary_multiple')
BEGIN
    ALTER TABLE member_rating_results ADD ptd_salary_multiple DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ptd_salary_multiple DECIMAL(15,5);
END;

-- Add or modify column for field: CiSalaryMultiple
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_salary_multiple')
BEGIN
    ALTER TABLE member_rating_results ADD ci_salary_multiple DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ci_salary_multiple DECIMAL(15,5);
END;

-- Add or modify column for field: Occupation
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'occupation')
BEGIN
    ALTER TABLE member_rating_results ADD occupation NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN occupation NVARCHAR(255);
END;

-- Add or modify column for field: OccupationClass
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'occupation_class')
BEGIN
    ALTER TABLE member_rating_results ADD occupation_class INT;
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN occupation_class INT;
END;

-- Add or modify column for field: Industry
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'industry')
BEGIN
    ALTER TABLE member_rating_results ADD industry NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN industry NVARCHAR(255);
END;

-- Add or modify column for field: AgeNextBirthday
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'age_next_birthday')
BEGIN
    ALTER TABLE member_rating_results ADD age_next_birthday INT;
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN age_next_birthday INT;
END;

-- Add or modify column for field: AgeBand
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'age_band')
BEGIN
    ALTER TABLE member_rating_results ADD age_band NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN age_band NVARCHAR(255);
END;

-- Add or modify column for field: SpouseGender
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_gender')
BEGIN
    ALTER TABLE member_rating_results ADD spouse_gender NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN spouse_gender NVARCHAR(255);
END;

-- Add or modify column for field: SpouseAgeNextBirthday
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_age_next_birthday')
BEGIN
    ALTER TABLE member_rating_results ADD spouse_age_next_birthday INT;
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN spouse_age_next_birthday INT;
END;

-- Add or modify column for field: AverageDependantAgeNextBirthday
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'average_dependant_age_next_birthday')
BEGIN
    ALTER TABLE member_rating_results ADD average_dependant_age_next_birthday DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN average_dependant_age_next_birthday DECIMAL(15,5);
END;

-- Add or modify column for field: AverageChildAgeNextBirthday
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'average_child_age_next_birthday')
BEGIN
    ALTER TABLE member_rating_results ADD average_child_age_next_birthday DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN average_child_age_next_birthday DECIMAL(15,5);
END;

-- Add or modify column for field: AverageNumberDependants
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'average_number_dependants')
BEGIN
    ALTER TABLE member_rating_results ADD average_number_dependants DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN average_number_dependants DECIMAL(15,5);
END;

-- Add or modify column for field: AverageNumberChildren
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'average_number_children')
BEGIN
    ALTER TABLE member_rating_results ADD average_number_children DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN average_number_children DECIMAL(15,5);
END;

-- Add or modify column for field: CalculatedFreeCoverLimit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'calculated_free_cover_limit')
BEGIN
    ALTER TABLE member_rating_results ADD calculated_free_cover_limit DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN calculated_free_cover_limit DECIMAL(15,5);
END;

-- Add or modify column for field: AppliedFreeCoverLimit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'applied_free_cover_limit')
BEGIN
    ALTER TABLE member_rating_results ADD applied_free_cover_limit DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN applied_free_cover_limit DECIMAL(15,5);
END;

-- Add or modify column for field: GlaSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD gla_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: GlaCappedSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_capped_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD gla_capped_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_capped_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: ExpenseLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'expense_loading')
BEGIN
    ALTER TABLE member_rating_results ADD expense_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN expense_loading DECIMAL(15,5);
END;

-- Add or modify column for field: AdminLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'admin_loading')
BEGIN
    ALTER TABLE member_rating_results ADD admin_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN admin_loading DECIMAL(15,5);
END;

-- Add or modify column for field: CommissionLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'commission_loading')
BEGIN
    ALTER TABLE member_rating_results ADD commission_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN commission_loading DECIMAL(15,5);
END;

-- Add or modify column for field: ProfitLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'profit_loading')
BEGIN
    ALTER TABLE member_rating_results ADD profit_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN profit_loading DECIMAL(15,5);
END;

-- Add or modify column for field: OtherLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'other_loading')
BEGIN
    ALTER TABLE member_rating_results ADD other_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN other_loading DECIMAL(15,5);
END;

-- Add or modify column for field: Discount
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'discount')
BEGIN
    ALTER TABLE member_rating_results ADD discount DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN discount DECIMAL(15,5);
END;

-- Add or modify column for field: TotalLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'total_loading')
BEGIN
    ALTER TABLE member_rating_results ADD total_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN total_loading DECIMAL(15,5);
END;

-- Add or modify column for field: GlaRegionLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_region_loading')
BEGIN
    ALTER TABLE member_rating_results ADD gla_region_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_region_loading DECIMAL(15,5);
END;

-- Add or modify column for field: GlaAidsRegionLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_aids_region_loading')
BEGIN
    ALTER TABLE member_rating_results ADD gla_aids_region_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_aids_region_loading DECIMAL(15,5);
END;

-- Add or modify column for field: PtdRegionLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_region_loading')
BEGIN
    ALTER TABLE member_rating_results ADD ptd_region_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ptd_region_loading DECIMAL(15,5);
END;

-- Add or modify column for field: CiRegionLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_region_loading')
BEGIN
    ALTER TABLE member_rating_results ADD ci_region_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ci_region_loading DECIMAL(15,5);
END;

-- Add or modify column for field: TtdRegionLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_region_loading')
BEGIN
    ALTER TABLE member_rating_results ADD ttd_region_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ttd_region_loading DECIMAL(15,5);
END;

-- Add or modify column for field: PhiRegionLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_region_loading')
BEGIN
    ALTER TABLE member_rating_results ADD phi_region_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN phi_region_loading DECIMAL(15,5);
END;

-- Add or modify column for field: FunRegionLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'fun_region_loading')
BEGIN
    ALTER TABLE member_rating_results ADD fun_region_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN fun_region_loading DECIMAL(15,5);
END;

-- Add or modify column for field: FunAidsRegionLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'fun_aids_region_loading')
BEGIN
    ALTER TABLE member_rating_results ADD fun_aids_region_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN fun_aids_region_loading DECIMAL(15,5);
END;

-- Add or modify column for field: GlaIndustryLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_industry_loading')
BEGIN
    ALTER TABLE member_rating_results ADD gla_industry_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_industry_loading DECIMAL(15,5);
END;

-- Add or modify column for field: PtdIndustryLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_industry_loading')
BEGIN
    ALTER TABLE member_rating_results ADD ptd_industry_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ptd_industry_loading DECIMAL(15,5);
END;

-- Add or modify column for field: CiIndustryLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_industry_loading')
BEGIN
    ALTER TABLE member_rating_results ADD ci_industry_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ci_industry_loading DECIMAL(15,5);
END;

-- Add or modify column for field: TtdIndustryLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_industry_loading')
BEGIN
    ALTER TABLE member_rating_results ADD ttd_industry_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ttd_industry_loading DECIMAL(15,5);
END;

-- Add or modify column for field: PhiIndustryLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_industry_loading')
BEGIN
    ALTER TABLE member_rating_results ADD phi_industry_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN phi_industry_loading DECIMAL(15,5);
END;

-- Add or modify column for field: GlaContingencyLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_contingency_loading')
BEGIN
    ALTER TABLE member_rating_results ADD gla_contingency_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_contingency_loading DECIMAL(15,5);
END;

-- Add or modify column for field: PtdContingencyLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_contingency_loading')
BEGIN
    ALTER TABLE member_rating_results ADD ptd_contingency_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ptd_contingency_loading DECIMAL(15,5);
END;

-- Add or modify column for field: CiContingencyLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_contingency_loading')
BEGIN
    ALTER TABLE member_rating_results ADD ci_contingency_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ci_contingency_loading DECIMAL(15,5);
END;

-- Add or modify column for field: TtdContingencyLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_contingency_loading')
BEGIN
    ALTER TABLE member_rating_results ADD ttd_contingency_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ttd_contingency_loading DECIMAL(15,5);
END;

-- Add or modify column for field: PhiContingencyLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_contingency_loading')
BEGIN
    ALTER TABLE member_rating_results ADD phi_contingency_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN phi_contingency_loading DECIMAL(15,5);
END;

-- Add or modify column for field: FunContingencyLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'fun_contingency_loading')
BEGIN
    ALTER TABLE member_rating_results ADD fun_contingency_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN fun_contingency_loading DECIMAL(15,5);
END;

-- Add or modify column for field: ContinuationLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'continuation_loading')
BEGIN
    ALTER TABLE member_rating_results ADD continuation_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN continuation_loading DECIMAL(15,5);
END;

-- Add or modify column for field: GlaQx
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_qx')
BEGIN
    ALTER TABLE member_rating_results ADD gla_qx DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_qx DECIMAL(15,5);
END;

-- Add or modify column for field: GlaAidsQx
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_aids_qx')
BEGIN
    ALTER TABLE member_rating_results ADD gla_aids_qx DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_aids_qx DECIMAL(15,5);
END;

-- Add or modify column for field: BaseGlaRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'base_gla_rate')
BEGIN
    ALTER TABLE member_rating_results ADD base_gla_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN base_gla_rate DECIMAL(15,5);
END;

-- Add or modify column for field: GlaLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_loading')
BEGIN
    ALTER TABLE member_rating_results ADD gla_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_loading DECIMAL(15,5);
END;

-- Add or modify column for field: GlaTerminalIllnessLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_terminal_illness_loading')
BEGIN
    ALTER TABLE member_rating_results ADD gla_terminal_illness_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_terminal_illness_loading DECIMAL(15,5);
END;

-- Add or modify column for field: LoadedGlaRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'loaded_gla_rate')
BEGIN
    ALTER TABLE member_rating_results ADD loaded_gla_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN loaded_gla_rate DECIMAL(15,5);
END;

-- Add or modify column for field: GlaWeightedExperienceCrudeRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_weighted_experience_crude_rate')
BEGIN
    ALTER TABLE member_rating_results ADD gla_weighted_experience_crude_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_weighted_experience_crude_rate DECIMAL(15,5);
END;

-- Add or modify column for field: GlaTheoreticalRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_theoretical_rate')
BEGIN
    ALTER TABLE member_rating_results ADD gla_theoretical_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_theoretical_rate DECIMAL(15,5);
END;

-- Add or modify column for field: PtdExperienceCrudeRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_experience_crude_rate')
BEGIN
    ALTER TABLE member_rating_results ADD ptd_experience_crude_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ptd_experience_crude_rate DECIMAL(15,5);
END;

-- Add or modify column for field: PtdTheoreticalRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_theoretical_rate')
BEGIN
    ALTER TABLE member_rating_results ADD ptd_theoretical_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ptd_theoretical_rate DECIMAL(15,5);
END;

-- Add or modify column for field: CiExperienceCrudeRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_experience_crude_rate')
BEGIN
    ALTER TABLE member_rating_results ADD ci_experience_crude_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ci_experience_crude_rate DECIMAL(15,5);
END;

-- Add or modify column for field: CiTheoreticalRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_theoretical_rate')
BEGIN
    ALTER TABLE member_rating_results ADD ci_theoretical_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ci_theoretical_rate DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjLoadedGlaRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_loaded_gla_rate')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_loaded_gla_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_loaded_gla_rate DECIMAL(15,5);
END;

-- Add or modify column for field: GlaExperienceAdjustment
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_experience_adjustment')
BEGIN
    ALTER TABLE member_rating_results ADD gla_experience_adjustment DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_experience_adjustment DECIMAL(15,5);
END;

-- Add or modify column for field: GlaRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD gla_risk_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_risk_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjGlaRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_gla_risk_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_gla_risk_premium DECIMAL(15,5);
END;

-- Add or modify column for field: GlaOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD gla_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN gla_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjGlaOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_gla_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_gla_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: PtdSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD ptd_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ptd_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: PtdCappedSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_capped_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD ptd_capped_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ptd_capped_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: BasePtdRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'base_ptd_rate')
BEGIN
    ALTER TABLE member_rating_results ADD base_ptd_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN base_ptd_rate DECIMAL(15,5);
END;

-- Add or modify column for field: PtdLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_loading')
BEGIN
    ALTER TABLE member_rating_results ADD ptd_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ptd_loading DECIMAL(15,5);
END;

-- Add or modify column for field: LoadedPtdRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'loaded_ptd_rate')
BEGIN
    ALTER TABLE member_rating_results ADD loaded_ptd_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN loaded_ptd_rate DECIMAL(15,5);
END;

-- Add or modify column for field: PtdExperienceAdjustment
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_experience_adjustment')
BEGIN
    ALTER TABLE member_rating_results ADD ptd_experience_adjustment DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ptd_experience_adjustment DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjLoadedPtdRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_loaded_ptd_rate')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_loaded_ptd_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_loaded_ptd_rate DECIMAL(15,5);
END;

-- Add or modify column for field: PtdRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD ptd_risk_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ptd_risk_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjPtdRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ptd_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_ptd_risk_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_ptd_risk_premium DECIMAL(15,5);
END;

-- Add or modify column for field: PtdOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD ptd_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ptd_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjPtdOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ptd_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_ptd_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_ptd_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: CiSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD ci_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ci_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: CiCappedSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_capped_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD ci_capped_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ci_capped_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: BaseCiRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'base_ci_rate')
BEGIN
    ALTER TABLE member_rating_results ADD base_ci_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN base_ci_rate DECIMAL(15,5);
END;

-- Add or modify column for field: CiLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_loading')
BEGIN
    ALTER TABLE member_rating_results ADD ci_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ci_loading DECIMAL(15,5);
END;

-- Add or modify column for field: LoadedCiRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'loaded_ci_rate')
BEGIN
    ALTER TABLE member_rating_results ADD loaded_ci_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN loaded_ci_rate DECIMAL(15,5);
END;

-- Add or modify column for field: CiExperienceAdjustment
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_experience_adjustment')
BEGIN
    ALTER TABLE member_rating_results ADD ci_experience_adjustment DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ci_experience_adjustment DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjLoadedCiRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_loaded_ci_rate')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_loaded_ci_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_loaded_ci_rate DECIMAL(15,5);
END;

-- Add or modify column for field: CiRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD ci_risk_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ci_risk_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjCiRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ci_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_ci_risk_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_ci_risk_premium DECIMAL(15,5);
END;

-- Add or modify column for field: CiOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD ci_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ci_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjCiOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ci_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_ci_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_ci_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: SpouseGlaSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_gla_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD spouse_gla_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN spouse_gla_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: SpouseGlaCappedSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_gla_capped_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD spouse_gla_capped_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN spouse_gla_capped_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: SpouseGlaQx
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_gla_qx')
BEGIN
    ALTER TABLE member_rating_results ADD spouse_gla_qx DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN spouse_gla_qx DECIMAL(15,5);
END;

-- Add or modify column for field: SpouseGlaAidsQx
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_gla_aids_qx')
BEGIN
    ALTER TABLE member_rating_results ADD spouse_gla_aids_qx DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN spouse_gla_aids_qx DECIMAL(15,5);
END;

-- Add or modify column for field: BaseSpouseGlaRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'base_spouse_gla_rate')
BEGIN
    ALTER TABLE member_rating_results ADD base_spouse_gla_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN base_spouse_gla_rate DECIMAL(15,5);
END;

-- Add or modify column for field: SpouseGlaLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_gla_loading')
BEGIN
    ALTER TABLE member_rating_results ADD spouse_gla_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN spouse_gla_loading DECIMAL(15,5);
END;

-- Add or modify column for field: LoadedSpouseGlaRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'loaded_spouse_gla_rate')
BEGIN
    ALTER TABLE member_rating_results ADD loaded_spouse_gla_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN loaded_spouse_gla_rate DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjLoadedSpouseGlaRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_loaded_spouse_gla_rate')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_loaded_spouse_gla_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_loaded_spouse_gla_rate DECIMAL(15,5);
END;

-- Add or modify column for field: SpouseGlaRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_gla_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD spouse_gla_risk_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN spouse_gla_risk_premium DECIMAL(15,5);
END;

-- Add or modify column for field: SpouseGlaOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_gla_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD spouse_gla_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN spouse_gla_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjSpouseGlaOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_spouse_gla_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_spouse_gla_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_spouse_gla_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: TtdIncome
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_income')
BEGIN
    ALTER TABLE member_rating_results ADD ttd_income DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ttd_income DECIMAL(15,5);
END;

-- Add or modify column for field: TtdCappedIncome
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_capped_income')
BEGIN
    ALTER TABLE member_rating_results ADD ttd_capped_income DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ttd_capped_income DECIMAL(15,5);
END;

-- Add or modify column for field: TtdNumberOfMonthlyPayments
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_number_of_monthly_payments')
BEGIN
    ALTER TABLE member_rating_results ADD ttd_number_of_monthly_payments DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ttd_number_of_monthly_payments DECIMAL(15,5);
END;

-- Add or modify column for field: IncomeReplacementRatio
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'income_replacement_ratio')
BEGIN
    ALTER TABLE member_rating_results ADD income_replacement_ratio DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN income_replacement_ratio DECIMAL(15,5);
END;

-- Add or modify column for field: BaseTtdRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'base_ttd_rate')
BEGIN
    ALTER TABLE member_rating_results ADD base_ttd_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN base_ttd_rate DECIMAL(15,5);
END;

-- Add or modify column for field: TtdLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_loading')
BEGIN
    ALTER TABLE member_rating_results ADD ttd_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ttd_loading DECIMAL(15,5);
END;

-- Add or modify column for field: LoadedTtdRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'loaded_ttd_rate')
BEGIN
    ALTER TABLE member_rating_results ADD loaded_ttd_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN loaded_ttd_rate DECIMAL(15,5);
END;

-- Add or modify column for field: TtdExperienceAdjustment
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_experience_adjustment')
BEGIN
    ALTER TABLE member_rating_results ADD ttd_experience_adjustment DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ttd_experience_adjustment DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjLoadedTtdRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_loaded_ttd_rate')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_loaded_ttd_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_loaded_ttd_rate DECIMAL(15,5);
END;

-- Add or modify column for field: TtdRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD ttd_risk_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ttd_risk_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjTtdRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ttd_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_ttd_risk_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_ttd_risk_premium DECIMAL(15,5);
END;

-- Add or modify column for field: TtdOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD ttd_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN ttd_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjTtdOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ttd_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_ttd_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_ttd_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjSpouseGlaRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_spouse_gla_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_spouse_gla_risk_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_spouse_gla_risk_premium DECIMAL(15,5);
END;

-- Add or modify column for field: PhiIncome
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_income')
BEGIN
    ALTER TABLE member_rating_results ADD phi_income DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN phi_income DECIMAL(15,5);
END;

-- Add or modify column for field: PhiCappedIncome
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_capped_income')
BEGIN
    ALTER TABLE member_rating_results ADD phi_capped_income DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN phi_capped_income DECIMAL(15,5);
END;

-- Add or modify column for field: PhiContributionWaiver
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_contribution_waiver')
BEGIN
    ALTER TABLE member_rating_results ADD phi_contribution_waiver DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN phi_contribution_waiver DECIMAL(15,5);
END;

-- Add or modify column for field: PhiMedicalAidWaiver
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_medical_aid_waiver')
BEGIN
    ALTER TABLE member_rating_results ADD phi_medical_aid_waiver DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN phi_medical_aid_waiver DECIMAL(15,5);
END;

-- Add or modify column for field: PhiMonthlyBenefit
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_monthly_benefit')
BEGIN
    ALTER TABLE member_rating_results ADD phi_monthly_benefit DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN phi_monthly_benefit DECIMAL(15,5);
END;

-- Add or modify column for field: PhiAnnuityFactor
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_annuity_factor')
BEGIN
    ALTER TABLE member_rating_results ADD phi_annuity_factor DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN phi_annuity_factor DECIMAL(15,5);
END;

-- Add or modify column for field: BasePhiRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'base_phi_rate')
BEGIN
    ALTER TABLE member_rating_results ADD base_phi_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN base_phi_rate DECIMAL(15,5);
END;

-- Add or modify column for field: PhiSalaryLevel
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_salary_level')
BEGIN
    ALTER TABLE member_rating_results ADD phi_salary_level DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN phi_salary_level DECIMAL(15,5);
END;

-- Add or modify column for field: PhiLoading
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_loading')
BEGIN
    ALTER TABLE member_rating_results ADD phi_loading DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN phi_loading DECIMAL(15,5);
END;

-- Add or modify column for field: LoadedPhiRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'loaded_phi_rate')
BEGIN
    ALTER TABLE member_rating_results ADD loaded_phi_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN loaded_phi_rate DECIMAL(15,5);
END;

-- Add or modify column for field: PhiExperienceAdjustment
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_experience_adjustment')
BEGIN
    ALTER TABLE member_rating_results ADD phi_experience_adjustment DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN phi_experience_adjustment DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjLoadedPhiRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_loaded_phi_rate')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_loaded_phi_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_loaded_phi_rate DECIMAL(15,5);
END;

-- Add or modify column for field: PhiRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD phi_risk_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN phi_risk_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjPhiRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_phi_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_phi_risk_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_phi_risk_premium DECIMAL(15,5);
END;

-- Add or modify column for field: PhiOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD phi_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN phi_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjPhiOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_phi_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_phi_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_phi_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: MemberFuneralSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'member_funeral_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD member_funeral_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN member_funeral_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: MainMemberFuneralBaseRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'main_member_funeral_base_rate')
BEGIN
    ALTER TABLE member_rating_results ADD main_member_funeral_base_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN main_member_funeral_base_rate DECIMAL(15,5);
END;

-- Add or modify column for field: MainMemberFuneralRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'main_member_funeral_cost')
BEGIN
    ALTER TABLE member_rating_results ADD main_member_funeral_cost DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN main_member_funeral_cost DECIMAL(15,5);
END;

-- Add or modify column for field: MainMemberFuneralOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'main_member_funeral_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD main_member_funeral_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN main_member_funeral_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: SpouseFuneralSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_funeral_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD spouse_funeral_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN spouse_funeral_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: SpouseFuneralBaseRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_funeral_base_rate')
BEGIN
    ALTER TABLE member_rating_results ADD spouse_funeral_base_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN spouse_funeral_base_rate DECIMAL(15,5);
END;

-- Add or modify column for field: SpouseFuneralRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_funeral_cost')
BEGIN
    ALTER TABLE member_rating_results ADD spouse_funeral_cost DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN spouse_funeral_cost DECIMAL(15,5);
END;

-- Add or modify column for field: SpouseFuneralOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'spouse_funeral_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD spouse_funeral_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN spouse_funeral_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ChildFuneralBaseRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'child_funeral_base_rate')
BEGIN
    ALTER TABLE member_rating_results ADD child_funeral_base_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN child_funeral_base_rate DECIMAL(15,5);
END;

-- Add or modify column for field: ChildFuneralSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'child_funeral_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD child_funeral_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN child_funeral_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: ChildFuneralRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'children_funeral_cost')
BEGIN
    ALTER TABLE member_rating_results ADD children_funeral_cost DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN children_funeral_cost DECIMAL(15,5);
END;

-- Add or modify column for field: ChildrenFuneralOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'children_funeral_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD children_funeral_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN children_funeral_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ParentFuneralBaseRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'dependant_funeral_base_rate')
BEGIN
    ALTER TABLE member_rating_results ADD dependant_funeral_base_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN dependant_funeral_base_rate DECIMAL(15,5);
END;

-- Add or modify column for field: ParentFuneralSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'dependant_funeral_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD dependant_funeral_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN dependant_funeral_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: ParentFuneralRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'dependants_funeral_cost')
BEGIN
    ALTER TABLE member_rating_results ADD dependants_funeral_cost DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN dependants_funeral_cost DECIMAL(15,5);
END;

-- Add or modify column for field: ParentFuneralOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'dependants_funeral_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD dependants_funeral_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN dependants_funeral_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ParentFuneralSumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'parent_funeral_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD parent_funeral_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN parent_funeral_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: TotalFuneralRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'total_funeral_risk_cost')
BEGIN
    ALTER TABLE member_rating_results ADD total_funeral_risk_cost DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN total_funeral_risk_cost DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjTotalFuneralRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_total_funeral_risk_cost')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_total_funeral_risk_cost DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_total_funeral_risk_cost DECIMAL(15,5);
END;

-- Add or modify column for field: TotalFuneralOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'total_funeral_office_cost')
BEGIN
    ALTER TABLE member_rating_results ADD total_funeral_office_cost DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN total_funeral_office_cost DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjTotalFuneralOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_total_funeral_office_cost')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_total_funeral_office_cost DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_total_funeral_office_cost DECIMAL(15,5);
END;

-- Add or modify column for field: Grade0SumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'grade0_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD grade0_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN grade0_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: Grade17SumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'grade17_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD grade17_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN grade17_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: Grade812SumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'grade812_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD grade812_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN grade812_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: TertiarySumAssured
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'tertiary_sum_assured')
BEGIN
    ALTER TABLE member_rating_results ADD tertiary_sum_assured DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN tertiary_sum_assured DECIMAL(15,5);
END;

-- Add or modify column for field: Grade0RiskRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'grade0_risk_rate')
BEGIN
    ALTER TABLE member_rating_results ADD grade0_risk_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN grade0_risk_rate DECIMAL(15,5);
END;

-- Add or modify column for field: Grade17RiskRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'grade17_risk_rate')
BEGIN
    ALTER TABLE member_rating_results ADD grade17_risk_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN grade17_risk_rate DECIMAL(15,5);
END;

-- Add or modify column for field: Grade812RiskRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'grade812_risk_rate')
BEGIN
    ALTER TABLE member_rating_results ADD grade812_risk_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN grade812_risk_rate DECIMAL(15,5);
END;

-- Add or modify column for field: TertiaryRiskRate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'tertiary_risk_rate')
BEGIN
    ALTER TABLE member_rating_results ADD tertiary_risk_rate DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN tertiary_risk_rate DECIMAL(15,5);
END;

-- Add or modify column for field: EducatorRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'educator_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD educator_risk_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN educator_risk_premium DECIMAL(15,5);
END;

-- Add or modify column for field: EducatorOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'educator_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD educator_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN educator_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjEducatorRiskPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_educator_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_educator_risk_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_educator_risk_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExpAdjEducatorOfficePremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_educator_office_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_educator_office_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exp_adj_educator_office_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExceedsNormalRetirementAgeIndicator
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exceeds_normal_retirement_age_indicator')
BEGIN
    ALTER TABLE member_rating_results ADD exceeds_normal_retirement_age_indicator INT;
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exceeds_normal_retirement_age_indicator INT;
END;

-- Add or modify column for field: ExceedsFreeCoverLimitIndicator
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exceeds_free_cover_limit_indicator')
BEGIN
    ALTER TABLE member_rating_results ADD exceeds_free_cover_limit_indicator INT;
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN exceeds_free_cover_limit_indicator INT;
END;

-- Add or modify column for field: FuneralExperienceAdjustedAnnualPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'funeral_experience_adjusted_annual_premium')
BEGIN
    ALTER TABLE member_rating_results ADD funeral_experience_adjusted_annual_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN funeral_experience_adjusted_annual_premium DECIMAL(15,5);
END;

-- Add or modify column for field: CreationDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'creation_date')
BEGIN
    ALTER TABLE member_rating_results ADD creation_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN creation_date DATETIME2;
END;

-- Add or modify column for field: CreatedBy
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'created_by')
BEGIN
    ALTER TABLE member_rating_results ADD created_by NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE member_rating_results ALTER COLUMN created_by NVARCHAR(255);
END;

