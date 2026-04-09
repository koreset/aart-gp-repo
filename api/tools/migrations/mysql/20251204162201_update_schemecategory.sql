-- Migration for struct: SchemeCategory

-- Table: scheme_categories

-- Ensure table exists
CREATE TABLE IF NOT EXISTS scheme_categories (
    id INT AUTO_INCREMENT PRIMARY KEY
);

-- Drop FK constraint if it exists before modifying quote_id column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS WHERE TABLE_NAME='scheme_categories' AND CONSTRAINT_NAME='fk_group_pricing_quotes_scheme_categories' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories DROP FOREIGN KEY fk_group_pricing_quotes_scheme_categories;',
    'SELECT 1;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: QuoteId
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='quote_id' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN quote_id BIGINT;',
    'ALTER TABLE scheme_categories ADD COLUMN quote_id BIGINT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SchemeCategory
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='scheme_category' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN scheme_category VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN scheme_category VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Basis
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='basis' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN basis VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN basis VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: FreeCoverLimit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='free_cover_limit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN free_cover_limit DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN free_cover_limit DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdBenefit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ptd_benefit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ptd_benefit TINYINT(1);',
    'ALTER TABLE scheme_categories ADD COLUMN ptd_benefit TINYINT(1);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaBenefit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='gla_benefit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN gla_benefit TINYINT(1);',
    'ALTER TABLE scheme_categories ADD COLUMN gla_benefit TINYINT(1);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiBenefit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ci_benefit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ci_benefit TINYINT(1);',
    'ALTER TABLE scheme_categories ADD COLUMN ci_benefit TINYINT(1);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SglaBenefit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='sgla_benefit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN sgla_benefit TINYINT(1);',
    'ALTER TABLE scheme_categories ADD COLUMN sgla_benefit TINYINT(1);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiBenefit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='phi_benefit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN phi_benefit TINYINT(1);',
    'ALTER TABLE scheme_categories ADD COLUMN phi_benefit TINYINT(1);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdBenefit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ttd_benefit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ttd_benefit TINYINT(1);',
    'ALTER TABLE scheme_categories ADD COLUMN ttd_benefit TINYINT(1);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: FamilyFuneralBenefit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='family_funeral_benefit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN family_funeral_benefit TINYINT(1);',
    'ALTER TABLE scheme_categories ADD COLUMN family_funeral_benefit TINYINT(1);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaSalaryMultiple
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='gla_salary_multiple' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN gla_salary_multiple DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN gla_salary_multiple DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaTerminalIllnessBenefit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='gla_terminal_illness_benefit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN gla_terminal_illness_benefit VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN gla_terminal_illness_benefit VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaWaitingPeriod
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='gla_waiting_period' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN gla_waiting_period INT;',
    'ALTER TABLE scheme_categories ADD COLUMN gla_waiting_period INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaEducatorBenefit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='gla_educator_benefit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN gla_educator_benefit VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN gla_educator_benefit VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdRiskType
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ptd_risk_type' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ptd_risk_type VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN ptd_risk_type VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdBenefitType
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ptd_benefit_type' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ptd_benefit_type VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN ptd_benefit_type VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdSalaryMultiple
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ptd_salary_multiple' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ptd_salary_multiple DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN ptd_salary_multiple DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdDeferredPeriod
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ptd_deferred_period' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ptd_deferred_period INT;',
    'ALTER TABLE scheme_categories ADD COLUMN ptd_deferred_period INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdDisabilityDefinition
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ptd_disability_definition' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ptd_disability_definition VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN ptd_disability_definition VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdEducatorBenefit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ptd_educator_benefit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ptd_educator_benefit VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN ptd_educator_benefit VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiBenefitStructure
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ci_benefit_structure' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ci_benefit_structure VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN ci_benefit_structure VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiBenefitDefinition
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ci_benefit_definition' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ci_benefit_definition VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN ci_benefit_definition VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiMaxBenefit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ci_max_benefit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ci_max_benefit DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN ci_max_benefit DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiCriticalIllnessSalaryMultiple
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ci_critical_illness_salary_multiple' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ci_critical_illness_salary_multiple DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN ci_critical_illness_salary_multiple DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SglaSalaryMultiple
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='sgla_salary_multiple' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN sgla_salary_multiple DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN sgla_salary_multiple DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SglaMaxBenefit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='sgla_max_benefit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN sgla_max_benefit DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN sgla_max_benefit DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiRiskType
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='phi_risk_type' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN phi_risk_type VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN phi_risk_type VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiMaximumBenefit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='phi_maximum_benefit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN phi_maximum_benefit DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN phi_maximum_benefit DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiIncomeReplacementPercentage
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='phi_income_replacement_percentage' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN phi_income_replacement_percentage DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN phi_income_replacement_percentage DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiPremiumWaiver
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='phi_premium_waiver' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN phi_premium_waiver VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN phi_premium_waiver VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiMedicalAidPremiumWaiver
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='phi_medical_aid_premium_waiver' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN phi_medical_aid_premium_waiver VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN phi_medical_aid_premium_waiver VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiBenefitEscalation
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='phi_benefit_escalation' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN phi_benefit_escalation VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN phi_benefit_escalation VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiMaxPremiumWaiver
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='phi_max_premium_waiver' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN phi_max_premium_waiver DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN phi_max_premium_waiver DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiWaitingPeriod
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='phi_waiting_period' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN phi_waiting_period INT;',
    'ALTER TABLE scheme_categories ADD COLUMN phi_waiting_period INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiNumberMonthlyPayments
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='phi_number_monthly_payments' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN phi_number_monthly_payments INT;',
    'ALTER TABLE scheme_categories ADD COLUMN phi_number_monthly_payments INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiDeferredPeriod
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='phi_deferred_period' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN phi_deferred_period INT;',
    'ALTER TABLE scheme_categories ADD COLUMN phi_deferred_period INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiDisabilityDefinition
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='phi_disability_definition' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN phi_disability_definition VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN phi_disability_definition VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiNormalRetirementAge
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='phi_normal_retirement_age' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN phi_normal_retirement_age INT;',
    'ALTER TABLE scheme_categories ADD COLUMN phi_normal_retirement_age INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdRiskType
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ttd_risk_type' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ttd_risk_type VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN ttd_risk_type VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdMaximumBenefit
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ttd_maximum_benefit' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ttd_maximum_benefit DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN ttd_maximum_benefit DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdIncomeReplacementPercentage
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ttd_income_replacement_percentage' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ttd_income_replacement_percentage DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN ttd_income_replacement_percentage DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdPremiumWaiverPercentage
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ttd_premium_waiver_percentage' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ttd_premium_waiver_percentage DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN ttd_premium_waiver_percentage DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdWaitingPeriod
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ttd_waiting_period' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ttd_waiting_period INT;',
    'ALTER TABLE scheme_categories ADD COLUMN ttd_waiting_period INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdNumberMonthlyPayments
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ttd_number_monthly_payments' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ttd_number_monthly_payments DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN ttd_number_monthly_payments DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdDeferredPeriod
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ttd_deferred_period' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ttd_deferred_period INT;',
    'ALTER TABLE scheme_categories ADD COLUMN ttd_deferred_period INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdDisabilityDefinition
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ttd_disability_definition' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ttd_disability_definition VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN ttd_disability_definition VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: FamilyFuneralMainMemberFuneralSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='family_funeral_main_member_funeral_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN family_funeral_main_member_funeral_sum_assured DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN family_funeral_main_member_funeral_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: FamilyFuneralSpouseFuneralSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='family_funeral_spouse_funeral_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN family_funeral_spouse_funeral_sum_assured DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN family_funeral_spouse_funeral_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: FamilyFuneralChildrenFuneralSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='family_funeral_children_funeral_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN family_funeral_children_funeral_sum_assured DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN family_funeral_children_funeral_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: FamilyFuneralAdultDependantSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='family_funeral_adult_dependant_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN family_funeral_adult_dependant_sum_assured DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN family_funeral_adult_dependant_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: FamilyFuneralParentFuneralSumAssured
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='family_funeral_parent_funeral_sum_assured' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN family_funeral_parent_funeral_sum_assured DOUBLE;',
    'ALTER TABLE scheme_categories ADD COLUMN family_funeral_parent_funeral_sum_assured DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: FamilyFuneralMaxNumberChildren
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='family_funeral_max_number_children' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN family_funeral_max_number_children INT;',
    'ALTER TABLE scheme_categories ADD COLUMN family_funeral_max_number_children INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: FamilyFuneralMaxNumberAdultDependants
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='family_funeral_max_number_adult_dependants' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN family_funeral_max_number_adult_dependants INT;',
    'ALTER TABLE scheme_categories ADD COLUMN family_funeral_max_number_adult_dependants INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdAlias
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ptd_alias' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ptd_alias VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN ptd_alias VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiAlias
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ci_alias' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ci_alias VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN ci_alias VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SglaAlias
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='sgla_alias' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN sgla_alias VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN sgla_alias VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiAlias
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='phi_alias' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN phi_alias VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN phi_alias VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdAlias
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ttd_alias' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN ttd_alias VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN ttd_alias VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaAlias
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='gla_alias' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN gla_alias VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN gla_alias VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: FamilyFuneralAlias
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='family_funeral_alias' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_categories MODIFY COLUMN family_funeral_alias VARCHAR(255);',
    'ALTER TABLE scheme_categories ADD COLUMN family_funeral_alias VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

