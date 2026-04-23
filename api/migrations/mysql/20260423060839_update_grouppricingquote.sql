-- Migration for struct: GroupPricingQuote

-- Table: group_pricing_quotes

-- Ensure table exists
CREATE TABLE IF NOT EXISTS group_pricing_quotes (
    id INT AUTO_INCREMENT PRIMARY KEY
);

-- Add or modify column for field: QuoteName
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='quote_name' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN quote_name VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN quote_name VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Basis
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='basis' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN basis VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN basis VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreationDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='creation_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN creation_date datetime;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN creation_date datetime;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: QuoteType
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='quote_type' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN quote_type VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN quote_type VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SchemeName
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='scheme_name' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN scheme_name VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN scheme_name VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SchemeID
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='scheme_id' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN scheme_id INT;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN scheme_id INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SchemeContact
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='scheme_contact' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN scheme_contact VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN scheme_contact VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SchemeEmail
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='scheme_email' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN scheme_email VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN scheme_email VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ID
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='broker_id' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN broker_id INT;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN broker_id INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Name
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='broker_name' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN broker_name VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN broker_name VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: DistributionChannel
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='distribution_channel' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN distribution_channel VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN distribution_channel VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ObligationType
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='obligation_type' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN obligation_type VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN obligation_type VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CommencementDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='commencement_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN commencement_date DATETIME;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN commencement_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CoverEndDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='cover_end_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN cover_end_date DATETIME;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN cover_end_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Industry
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='industry' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN industry VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN industry VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: OccupationClass
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='occupation_class' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN occupation_class INT;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN occupation_class INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: FreeCoverLimit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='free_cover_limit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN free_cover_limit DOUBLE;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN free_cover_limit DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Currency
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='currency' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN currency VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN currency VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExchangeRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='exchange_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN exchange_rate INT;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN exchange_rate INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: NormalRetirementAge
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='normal_retirement_age' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN normal_retirement_age INT;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN normal_retirement_age INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExperienceRating
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='experience_rating' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN experience_rating VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN experience_rating VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreatedBy
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='created_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN created_by VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN created_by VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Reviewer
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='reviewer' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN reviewer VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN reviewer VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ApprovedBy
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='approved_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN approved_by VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN approved_by VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SentBy
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='sent_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN sent_by VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN sent_by VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ModifiedBy
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='modified_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN modified_by VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN modified_by VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ModificationDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='modification_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN modification_date DATETIME;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN modification_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Status
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='status' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN status VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN status VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MemberDataCount
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='member_data_count' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN member_data_count INT;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN member_data_count INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ClaimsExperienceCount
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='claims_experience_count' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN claims_experience_count INT;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN claims_experience_count INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MemberRatingResultCount
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='member_rating_result_count' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN member_rating_result_count INT;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN member_rating_result_count INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MemberPremiumScheduleCount
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='member_premium_schedule_count' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN member_premium_schedule_count INT;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN member_premium_schedule_count INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: BordereauxCount
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='bordereaux_count' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN bordereaux_count INT;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN bordereaux_count INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: UseGlobalSalaryMultiple
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='use_global_salary_multiple' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN use_global_salary_multiple TINYINT(1);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN use_global_salary_multiple TINYINT(1);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SelectedSchemeCategories
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='selected_scheme_categories' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN selected_scheme_categories json;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN selected_scheme_categories json;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SchemeCategories
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='scheme_categories' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN scheme_categories TEXT;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN scheme_categories TEXT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CommissionLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='loadings_commission_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN loadings_commission_loading DOUBLE;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN loadings_commission_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ProfitLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='loadings_profit_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN loadings_profit_loading DOUBLE;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN loadings_profit_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpenseLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='loadings_expense_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN loadings_expense_loading DOUBLE;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN loadings_expense_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: AdminLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='loadings_admin_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN loadings_admin_loading DOUBLE;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN loadings_admin_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ContingencyLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='loadings_contingency_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN loadings_contingency_loading DOUBLE;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN loadings_contingency_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: OtherLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='loadings_other_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN loadings_other_loading DOUBLE;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN loadings_other_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Discount
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='loadings_discount' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN loadings_discount DOUBLE;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN loadings_discount DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: BinderFee
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='loadings_binder_fee' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN loadings_binder_fee DOUBLE;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN loadings_binder_fee DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: OutsourceFee
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='loadings_outsource_fee' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN loadings_outsource_fee DOUBLE;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN loadings_outsource_fee DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MemberAverageAge
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='member_average_age' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN member_average_age INT;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN member_average_age INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MemberAverageIncome
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='member_average_income' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN member_average_income DOUBLE;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN member_average_income DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MemberMaleFemaleDistribution
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='member_male_female_distribution' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN member_male_female_distribution DOUBLE;',
    'ALTER TABLE group_pricing_quotes ADD COLUMN member_male_female_distribution DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MemberIndicativeData
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='member_indicative_data' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN member_indicative_data TINYINT(1);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN member_indicative_data TINYINT(1);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: RiskRateCode
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='risk_rate_code' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN risk_rate_code VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN risk_rate_code VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SchemeQuoteStatus
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_pricing_quotes' AND COLUMN_NAME='scheme_quote_status' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_pricing_quotes MODIFY COLUMN scheme_quote_status VARCHAR(255);',
    'ALTER TABLE group_pricing_quotes ADD COLUMN scheme_quote_status VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

