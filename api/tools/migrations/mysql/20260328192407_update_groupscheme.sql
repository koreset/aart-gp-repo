-- Migration for struct: GroupScheme

-- Table: group_schemes

-- Ensure table exists
CREATE TABLE IF NOT EXISTS group_schemes (
    id INT AUTO_INCREMENT PRIMARY KEY
);

-- Add or modify column for field: Name
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='name' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN name VARCHAR(255);',
    'ALTER TABLE group_schemes ADD COLUMN name VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: BrokerId
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='broker_id' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN broker_id INT;',
    'ALTER TABLE group_schemes ADD COLUMN broker_id INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: DistributionChannel
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='distribution_channel' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN distribution_channel VARCHAR(255);',
    'ALTER TABLE group_schemes ADD COLUMN distribution_channel VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ContactPerson
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='contact_person' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN contact_person VARCHAR(255);',
    'ALTER TABLE group_schemes ADD COLUMN contact_person VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ContactEmail
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='contact_email' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN contact_email VARCHAR(255);',
    'ALTER TABLE group_schemes ADD COLUMN contact_email VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: DurationInForceDays
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='duration_in_force_days' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN duration_in_force_days INT;',
    'ALTER TABLE group_schemes ADD COLUMN duration_in_force_days INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: RenewalDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='renewal_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN renewal_date DATETIME;',
    'ALTER TABLE group_schemes ADD COLUMN renewal_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MemberCount
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='member_count' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN member_count DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN member_count DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: AnnualPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='annual_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN annual_premium DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN annual_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: GlaAnnualPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='gla_annual_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN gla_annual_premium DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN gla_annual_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PtdAnnualPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='ptd_annual_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN ptd_annual_premium DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN ptd_annual_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CiAnnualPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='ci_annual_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN ci_annual_premium DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN ci_annual_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SglaAnnualPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='sgla_annual_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN sgla_annual_premium DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN sgla_annual_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TtdAnnualPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='ttd_annual_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN ttd_annual_premium DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN ttd_annual_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: PhiAnnualPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='phi_annual_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN phi_annual_premium DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN phi_annual_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: FuneralAnnualPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='funeral_annual_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN funeral_annual_premium DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN funeral_annual_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Commission
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='commission' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN commission DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN commission DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: EarnedPremium
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='earned_premium' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN earned_premium DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN earned_premium DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpectedExpenses
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='expected_expenses' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN expected_expenses DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN expected_expenses DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpectedGlaClaims
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='expected_gla_claims' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN expected_gla_claims DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN expected_gla_claims DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpectedPtdClaims
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='expected_ptd_claims' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN expected_ptd_claims DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN expected_ptd_claims DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpectedCiClaims
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='expected_ci_claims' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN expected_ci_claims DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN expected_ci_claims DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpectedSglaClaims
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='expected_sgla_claims' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN expected_sgla_claims DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN expected_sgla_claims DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpectedTtdClaims
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='expected_ttd_claims' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN expected_ttd_claims DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN expected_ttd_claims DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpectedPhiClaims
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='expected_phi_claims' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN expected_phi_claims DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN expected_phi_claims DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpectedFunClaims
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='expected_fun_claims' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN expected_fun_claims DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN expected_fun_claims DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpectedClaims
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='expected_claims' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN expected_claims DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN expected_claims DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ActualClaims
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='actual_claims' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN actual_claims DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN actual_claims DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpectedClaimsRatio
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='expected_claims_ratio' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN expected_claims_ratio DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN expected_claims_ratio DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ActualClaimsRatio
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='actual_claims_ratio' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN actual_claims_ratio DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN actual_claims_ratio DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpectedLossRatio
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='expected_loss_ratio' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN expected_loss_ratio DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN expected_loss_ratio DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ActualLossRatio
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='actual_loss_ratio' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN actual_loss_ratio DOUBLE;',
    'ALTER TABLE group_schemes ADD COLUMN actual_loss_ratio DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: InForce
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='in_force' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN in_force TINYINT(1);',
    'ALTER TABLE group_schemes ADD COLUMN in_force TINYINT(1);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Status
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='status' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN status VARCHAR(255);',
    'ALTER TABLE group_schemes ADD COLUMN status VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: NewBusiness
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='new_business' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN new_business TINYINT(1);',
    'ALTER TABLE group_schemes ADD COLUMN new_business TINYINT(1);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SchemeStatusMessage
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='scheme_status_message' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN scheme_status_message VARCHAR(255);',
    'ALTER TABLE group_schemes ADD COLUMN scheme_status_message VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreationDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='creation_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN creation_date DATETIME;',
    'ALTER TABLE group_schemes ADD COLUMN creation_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreatedBy
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='created_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN created_by VARCHAR(255);',
    'ALTER TABLE group_schemes ADD COLUMN created_by VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: QuoteId
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='quote_id' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN quote_id INT;',
    'ALTER TABLE group_schemes ADD COLUMN quote_id INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Quote
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='quote' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN quote TEXT;',
    'ALTER TABLE group_schemes ADD COLUMN quote TEXT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: QuoteInForce
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='quote_in_force' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN quote_in_force VARCHAR(255);',
    'ALTER TABLE group_schemes ADD COLUMN quote_in_force VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ActiveSchemeCategories
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='active_scheme_categories' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN active_scheme_categories json;',
    'ALTER TABLE group_schemes ADD COLUMN active_scheme_categories json;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CoverStartDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='cover_start_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN cover_start_date DATETIME;',
    'ALTER TABLE group_schemes ADD COLUMN cover_start_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CoverEndDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='cover_end_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN cover_end_date DATETIME;',
    'ALTER TABLE group_schemes ADD COLUMN cover_end_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CommencementDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='commencement_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN commencement_date DATETIME;',
    'ALTER TABLE group_schemes ADD COLUMN commencement_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SchemeQuoteStatus
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='scheme_quote_status' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE group_schemes MODIFY COLUMN scheme_quote_status VARCHAR(255);',
    'ALTER TABLE group_schemes ADD COLUMN scheme_quote_status VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

