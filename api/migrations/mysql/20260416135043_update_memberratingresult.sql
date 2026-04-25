-- Migration for struct: MemberRatingResult

-- Table: member_rating_results

-- Ensure table exists
CREATE TABLE IF NOT EXISTS member_rating_results (
    id INT AUTO_INCREMENT PRIMARY KEY
);

-- Add or modify column for field: FinancialYear
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='financial_year' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN financial_year INT;',
    'ALTER TABLE member_rating_results ADD COLUMN financial_year INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SchemeId
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='scheme_id' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN scheme_id INT;',
    'ALTER TABLE member_rating_results ADD COLUMN scheme_id INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: QuoteId
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='quote_id' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN quote_id INT;',
    'ALTER TABLE member_rating_results ADD COLUMN quote_id INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Category
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='category' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN category VARCHAR(255);',
    'ALTER TABLE member_rating_results ADD COLUMN category VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MemberName
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='member_name' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN member_name VARCHAR(255);',
    'ALTER TABLE member_rating_results ADD COLUMN member_name VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MemberCount
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='member_count' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN member_count INT;',
    'ALTER TABLE member_rating_results ADD COLUMN member_count INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Gender
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gender' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN gender VARCHAR(255);',
    'ALTER TABLE member_rating_results ADD COLUMN gender VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: DateOfBirth
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='date_of_birth' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN date_of_birth DATETIME;',
    'ALTER TABLE member_rating_results ADD COLUMN date_of_birth DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: IsOriginalMember
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='is_original_member' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN is_original_member TINYINT(1);',
    'ALTER TABLE member_rating_results ADD COLUMN is_original_member TINYINT(1);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: EntryDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='entry_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN entry_date DATETIME;',
    'ALTER TABLE member_rating_results ADD COLUMN entry_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExitDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exit_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exit_date DATETIME;',
    'ALTER TABLE member_rating_results ADD COLUMN exit_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpCredibility
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_credibility' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_credibility DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_credibility DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ManuallyAddedCredibility
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='manually_added_credibility' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN manually_added_credibility DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN manually_added_credibility DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: AnnualSalary
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='annual_salary' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN annual_salary DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN annual_salary DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: IncomeLevel
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='income_level' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN income_level INT;',
    'ALTER TABLE member_rating_results ADD COLUMN income_level INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaSalaryMultiple
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_salary_multiple' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN gla_salary_multiple DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN gla_salary_multiple DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SglaSalaryMultiple
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='sgla_salary_multiple' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN sgla_salary_multiple DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN sgla_salary_multiple DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdSalaryMultiple
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_salary_multiple' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ptd_salary_multiple DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ptd_salary_multiple DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiSalaryMultiple
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ci_salary_multiple' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ci_salary_multiple DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ci_salary_multiple DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Occupation
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='occupation' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN occupation VARCHAR(255);',
    'ALTER TABLE member_rating_results ADD COLUMN occupation VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: OccupationClass
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='occupation_class' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN occupation_class INT;',
    'ALTER TABLE member_rating_results ADD COLUMN occupation_class INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Industry
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='industry' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN industry VARCHAR(255);',
    'ALTER TABLE member_rating_results ADD COLUMN industry VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: AgeNextBirthday
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='age_next_birthday' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN age_next_birthday INT;',
    'ALTER TABLE member_rating_results ADD COLUMN age_next_birthday INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: AgeBand
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='age_band' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN age_band VARCHAR(255);',
    'ALTER TABLE member_rating_results ADD COLUMN age_band VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SpouseGender
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_gender' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN spouse_gender VARCHAR(255);',
    'ALTER TABLE member_rating_results ADD COLUMN spouse_gender VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SpouseAgeNextBirthday
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_age_next_birthday' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN spouse_age_next_birthday INT;',
    'ALTER TABLE member_rating_results ADD COLUMN spouse_age_next_birthday INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: AverageDependantAgeNextBirthday
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='average_dependant_age_next_birthday' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN average_dependant_age_next_birthday DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN average_dependant_age_next_birthday DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: AverageChildAgeNextBirthday
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='average_child_age_next_birthday' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN average_child_age_next_birthday DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN average_child_age_next_birthday DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: AverageNumberDependants
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='average_number_dependants' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN average_number_dependants DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN average_number_dependants DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: AverageNumberChildren
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='average_number_children' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN average_number_children DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN average_number_children DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CalculatedFreeCoverLimit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='calculated_free_cover_limit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN calculated_free_cover_limit DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN calculated_free_cover_limit DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: AppliedFreeCoverLimit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='applied_free_cover_limit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN applied_free_cover_limit DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN applied_free_cover_limit DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN gla_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN gla_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaCappedSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_capped_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN gla_capped_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN gla_capped_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpenseLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='expense_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN expense_loading DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN expense_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: AdminLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='admin_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN admin_loading DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN admin_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CommissionLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='commission_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN commission_loading DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN commission_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ProfitLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='profit_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN profit_loading DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN profit_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: OtherLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='other_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN other_loading DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN other_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Discount
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='discount' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN discount DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN discount DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TotalPremiumLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='total_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN total_loading DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN total_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaQx
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_qx' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN gla_qx DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN gla_qx DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaAidsQx
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_aids_qx' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN gla_aids_qx DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN gla_aids_qx DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: BaseGlaRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='base_gla_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN base_gla_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN base_gla_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN gla_loading DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN gla_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaTerminalIllnessLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_terminal_illness_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN gla_terminal_illness_loading DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN gla_terminal_illness_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: LoadedGlaRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='loaded_gla_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN loaded_gla_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN loaded_gla_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaWeightedExperienceCrudeRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_weighted_experience_crude_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN gla_weighted_experience_crude_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN gla_weighted_experience_crude_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaTheoreticalRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_theoretical_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN gla_theoretical_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN gla_theoretical_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdExperienceCrudeRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_experience_crude_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ptd_experience_crude_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ptd_experience_crude_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdTheoreticalRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_theoretical_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ptd_theoretical_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ptd_theoretical_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiExperienceCrudeRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ci_experience_crude_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ci_experience_crude_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ci_experience_crude_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiTheoreticalRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ci_theoretical_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ci_theoretical_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ci_theoretical_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjLoadedGlaRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_loaded_gla_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_loaded_gla_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_loaded_gla_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaExperienceAdjustment
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_experience_adjustment' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN gla_experience_adjustment DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN gla_experience_adjustment DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_risk_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN gla_risk_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN gla_risk_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjGlaRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_risk_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_gla_risk_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_gla_risk_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN gla_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN gla_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjGlaOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_gla_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_gla_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ptd_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ptd_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdCappedSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_capped_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ptd_capped_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ptd_capped_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: BasePtdRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='base_ptd_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN base_ptd_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN base_ptd_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ptd_loading DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ptd_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: LoadedPtdRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='loaded_ptd_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN loaded_ptd_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN loaded_ptd_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdExperienceAdjustment
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_experience_adjustment' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ptd_experience_adjustment DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ptd_experience_adjustment DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjLoadedPtdRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_loaded_ptd_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_loaded_ptd_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_loaded_ptd_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_risk_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ptd_risk_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ptd_risk_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjPtdRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ptd_risk_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ptd_risk_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ptd_risk_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ptd_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ptd_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjPtdOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ptd_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ptd_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ptd_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ci_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ci_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ci_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiCappedSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ci_capped_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ci_capped_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ci_capped_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: BaseCiRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='base_ci_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN base_ci_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN base_ci_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ci_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ci_loading DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ci_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: LoadedCiRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='loaded_ci_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN loaded_ci_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN loaded_ci_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiExperienceAdjustment
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ci_experience_adjustment' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ci_experience_adjustment DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ci_experience_adjustment DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjLoadedCiRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_loaded_ci_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_loaded_ci_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_loaded_ci_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ci_risk_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ci_risk_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ci_risk_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjCiRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ci_risk_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ci_risk_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ci_risk_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ci_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ci_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ci_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjCiOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ci_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ci_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ci_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SpouseGlaSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_gla_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN spouse_gla_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN spouse_gla_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SpouseGlaCappedSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_gla_capped_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN spouse_gla_capped_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN spouse_gla_capped_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SpouseGlaQx
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_gla_qx' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN spouse_gla_qx DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN spouse_gla_qx DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SpouseGlaAidsQx
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_gla_aids_qx' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN spouse_gla_aids_qx DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN spouse_gla_aids_qx DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: BaseSpouseGlaRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='base_spouse_gla_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN base_spouse_gla_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN base_spouse_gla_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SpouseGlaLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_gla_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN spouse_gla_loading DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN spouse_gla_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: LoadedSpouseGlaRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='loaded_spouse_gla_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN loaded_spouse_gla_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN loaded_spouse_gla_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjLoadedSpouseGlaRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_loaded_spouse_gla_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_loaded_spouse_gla_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_loaded_spouse_gla_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SpouseGlaRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_gla_risk_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN spouse_gla_risk_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN spouse_gla_risk_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SpouseGlaOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_gla_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN spouse_gla_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN spouse_gla_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjSpouseGlaOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_spouse_gla_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_spouse_gla_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_spouse_gla_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdIncome
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_income' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ttd_income DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ttd_income DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdCappedIncome
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_capped_income' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ttd_capped_income DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ttd_capped_income DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdNumberOfMonthlyPayments
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_number_of_monthly_payments' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ttd_number_of_monthly_payments DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ttd_number_of_monthly_payments DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: IncomeReplacementRatio
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='income_replacement_ratio' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN income_replacement_ratio DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN income_replacement_ratio DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: BaseTtdRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='base_ttd_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN base_ttd_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN base_ttd_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ttd_loading DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ttd_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: LoadedTtdRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='loaded_ttd_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN loaded_ttd_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN loaded_ttd_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdExperienceAdjustment
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_experience_adjustment' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ttd_experience_adjustment DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ttd_experience_adjustment DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjLoadedTtdRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_loaded_ttd_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_loaded_ttd_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_loaded_ttd_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_risk_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ttd_risk_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ttd_risk_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjTtdRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ttd_risk_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ttd_risk_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ttd_risk_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN ttd_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN ttd_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjTtdOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ttd_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ttd_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ttd_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjSpouseGlaRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_spouse_gla_risk_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_spouse_gla_risk_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_spouse_gla_risk_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiIncome
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_income' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN phi_income DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN phi_income DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiCappedIncome
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_capped_income' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN phi_capped_income DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN phi_capped_income DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiContributionWaiver
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_contribution_waiver' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN phi_contribution_waiver DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN phi_contribution_waiver DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiMedicalAidWaiver
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_medical_aid_waiver' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN phi_medical_aid_waiver DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN phi_medical_aid_waiver DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiMonthlyBenefit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_monthly_benefit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN phi_monthly_benefit DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN phi_monthly_benefit DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiAnnuityFactor
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_annuity_factor' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN phi_annuity_factor DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN phi_annuity_factor DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: BasePhiRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='base_phi_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN base_phi_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN base_phi_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiSalaryLevel
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_salary_level' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN phi_salary_level DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN phi_salary_level DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN phi_loading DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN phi_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: LoadedPhiRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='loaded_phi_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN loaded_phi_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN loaded_phi_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiExperienceAdjustment
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_experience_adjustment' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN phi_experience_adjustment DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN phi_experience_adjustment DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjLoadedPhiRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_loaded_phi_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_loaded_phi_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_loaded_phi_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_risk_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN phi_risk_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN phi_risk_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjPhiRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_phi_risk_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_phi_risk_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_phi_risk_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN phi_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN phi_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjPhiOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_phi_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_phi_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_phi_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MemberFuneralSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='member_funeral_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN member_funeral_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN member_funeral_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MainMemberFuneralBaseRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='main_member_funeral_base_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN main_member_funeral_base_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN main_member_funeral_base_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MainMemberFuneralRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='main_member_funeral_cost' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN main_member_funeral_cost DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN main_member_funeral_cost DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MainMemberFuneralOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='main_member_funeral_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN main_member_funeral_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN main_member_funeral_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SpouseFuneralSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_funeral_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN spouse_funeral_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN spouse_funeral_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SpouseFuneralBaseRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_funeral_base_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN spouse_funeral_base_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN spouse_funeral_base_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SpouseFuneralRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_funeral_cost' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN spouse_funeral_cost DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN spouse_funeral_cost DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SpouseFuneralOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_funeral_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN spouse_funeral_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN spouse_funeral_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ChildFuneralBaseRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='child_funeral_base_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN child_funeral_base_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN child_funeral_base_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ChildFuneralSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='child_funeral_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN child_funeral_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN child_funeral_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ChildFuneralRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='children_funeral_cost' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN children_funeral_cost DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN children_funeral_cost DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ChildrenFuneralOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='children_funeral_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN children_funeral_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN children_funeral_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ParentFuneralBaseRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='dependant_funeral_base_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN dependant_funeral_base_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN dependant_funeral_base_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ParentFuneralSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='dependant_funeral_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN dependant_funeral_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN dependant_funeral_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ParentFuneralRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='dependants_funeral_cost' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN dependants_funeral_cost DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN dependants_funeral_cost DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ParentFuneralOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='dependants_funeral_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN dependants_funeral_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN dependants_funeral_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ParentFuneralSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='parent_funeral_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN parent_funeral_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN parent_funeral_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TotalFuneralRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='total_funeral_risk_cost' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN total_funeral_risk_cost DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN total_funeral_risk_cost DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjTotalFuneralRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_total_funeral_risk_cost' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_total_funeral_risk_cost DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_total_funeral_risk_cost DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TotalFuneralOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='total_funeral_office_cost' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN total_funeral_office_cost DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN total_funeral_office_cost DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjTotalFuneralOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_total_funeral_office_cost' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_total_funeral_office_cost DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_total_funeral_office_cost DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Grade0SumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='grade0_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN grade0_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN grade0_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Grade17SumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='grade17_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN grade17_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN grade17_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Grade812SumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='grade812_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN grade812_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN grade812_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TertiarySumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='tertiary_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN tertiary_sum_assured DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN tertiary_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Grade0RiskRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='grade0_risk_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN grade0_risk_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN grade0_risk_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Grade17RiskRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='grade17_risk_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN grade17_risk_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN grade17_risk_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Grade812RiskRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='grade812_risk_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN grade812_risk_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN grade812_risk_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TertiaryRiskRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='tertiary_risk_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN tertiary_risk_rate DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN tertiary_risk_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: EducatorRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='educator_risk_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN educator_risk_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN educator_risk_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: EducatorOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='educator_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN educator_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN educator_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjEducatorRiskPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_educator_risk_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_educator_risk_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_educator_risk_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpAdjEducatorOfficePremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_educator_office_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_educator_office_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN exp_adj_educator_office_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExceedsNormalRetirementAgeIndicator
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exceeds_normal_retirement_age_indicator' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exceeds_normal_retirement_age_indicator INT;',
    'ALTER TABLE member_rating_results ADD COLUMN exceeds_normal_retirement_age_indicator INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExceedsFreeCoverLimitIndicator
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exceeds_free_cover_limit_indicator' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN exceeds_free_cover_limit_indicator INT;',
    'ALTER TABLE member_rating_results ADD COLUMN exceeds_free_cover_limit_indicator INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: FuneralExperienceAdjustedAnnualPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='funeral_experience_adjusted_annual_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN funeral_experience_adjusted_annual_premium DOUBLE;',
    'ALTER TABLE member_rating_results ADD COLUMN funeral_experience_adjusted_annual_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreationDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='creation_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN creation_date DATETIME;',
    'ALTER TABLE member_rating_results ADD COLUMN creation_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreatedBy
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='created_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE member_rating_results MODIFY COLUMN created_by VARCHAR(255);',
    'ALTER TABLE member_rating_results ADD COLUMN created_by VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

