-- Migration for struct: GroupScheme

-- Table: group_schemes

-- Ensure table exists
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'group_schemes')
BEGIN
    CREATE TABLE group_schemes (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

-- Add or modify column for field: Name
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'name')
BEGIN
    ALTER TABLE group_schemes ADD name NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN name NVARCHAR(255);
END;

-- Add or modify column for field: BrokerId
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'broker_id')
BEGIN
    ALTER TABLE group_schemes ADD broker_id INT;
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN broker_id INT;
END;

-- Add or modify column for field: DistributionChannel
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'distribution_channel')
BEGIN
    ALTER TABLE group_schemes ADD distribution_channel NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN distribution_channel NVARCHAR(255);
END;

-- Add or modify column for field: ContactPerson
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'contact_person')
BEGIN
    ALTER TABLE group_schemes ADD contact_person NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN contact_person NVARCHAR(255);
END;

-- Add or modify column for field: ContactEmail
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'contact_email')
BEGIN
    ALTER TABLE group_schemes ADD contact_email NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN contact_email NVARCHAR(255);
END;

-- Add or modify column for field: DurationInForceDays
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'duration_in_force_days')
BEGIN
    ALTER TABLE group_schemes ADD duration_in_force_days INT;
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN duration_in_force_days INT;
END;

-- Add or modify column for field: RenewalDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'renewal_date')
BEGIN
    ALTER TABLE group_schemes ADD renewal_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN renewal_date DATETIME2;
END;

-- Add or modify column for field: MemberCount
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'member_count')
BEGIN
    ALTER TABLE group_schemes ADD member_count DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN member_count DECIMAL(15,5);
END;

-- Add or modify column for field: AnnualPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'annual_premium')
BEGIN
    ALTER TABLE group_schemes ADD annual_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN annual_premium DECIMAL(15,5);
END;

-- Add or modify column for field: GlaAnnualPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'gla_annual_premium')
BEGIN
    ALTER TABLE group_schemes ADD gla_annual_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN gla_annual_premium DECIMAL(15,5);
END;

-- Add or modify column for field: PtdAnnualPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'ptd_annual_premium')
BEGIN
    ALTER TABLE group_schemes ADD ptd_annual_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN ptd_annual_premium DECIMAL(15,5);
END;

-- Add or modify column for field: CiAnnualPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'ci_annual_premium')
BEGIN
    ALTER TABLE group_schemes ADD ci_annual_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN ci_annual_premium DECIMAL(15,5);
END;

-- Add or modify column for field: SglaAnnualPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'sgla_annual_premium')
BEGIN
    ALTER TABLE group_schemes ADD sgla_annual_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN sgla_annual_premium DECIMAL(15,5);
END;

-- Add or modify column for field: TtdAnnualPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'ttd_annual_premium')
BEGIN
    ALTER TABLE group_schemes ADD ttd_annual_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN ttd_annual_premium DECIMAL(15,5);
END;

-- Add or modify column for field: PhiAnnualPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'phi_annual_premium')
BEGIN
    ALTER TABLE group_schemes ADD phi_annual_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN phi_annual_premium DECIMAL(15,5);
END;

-- Add or modify column for field: FuneralAnnualPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'funeral_annual_premium')
BEGIN
    ALTER TABLE group_schemes ADD funeral_annual_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN funeral_annual_premium DECIMAL(15,5);
END;

-- Add or modify column for field: Commission
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'commission')
BEGIN
    ALTER TABLE group_schemes ADD commission DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN commission DECIMAL(15,5);
END;

-- Add or modify column for field: EarnedPremium
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'earned_premium')
BEGIN
    ALTER TABLE group_schemes ADD earned_premium DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN earned_premium DECIMAL(15,5);
END;

-- Add or modify column for field: ExpectedExpenses
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'expected_expenses')
BEGIN
    ALTER TABLE group_schemes ADD expected_expenses DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN expected_expenses DECIMAL(15,5);
END;

-- Add or modify column for field: ExpectedGlaClaims
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'expected_gla_claims')
BEGIN
    ALTER TABLE group_schemes ADD expected_gla_claims DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN expected_gla_claims DECIMAL(15,5);
END;

-- Add or modify column for field: ExpectedPtdClaims
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'expected_ptd_claims')
BEGIN
    ALTER TABLE group_schemes ADD expected_ptd_claims DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN expected_ptd_claims DECIMAL(15,5);
END;

-- Add or modify column for field: ExpectedCiClaims
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'expected_ci_claims')
BEGIN
    ALTER TABLE group_schemes ADD expected_ci_claims DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN expected_ci_claims DECIMAL(15,5);
END;

-- Add or modify column for field: ExpectedSglaClaims
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'expected_sgla_claims')
BEGIN
    ALTER TABLE group_schemes ADD expected_sgla_claims DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN expected_sgla_claims DECIMAL(15,5);
END;

-- Add or modify column for field: ExpectedTtdClaims
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'expected_ttd_claims')
BEGIN
    ALTER TABLE group_schemes ADD expected_ttd_claims DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN expected_ttd_claims DECIMAL(15,5);
END;

-- Add or modify column for field: ExpectedPhiClaims
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'expected_phi_claims')
BEGIN
    ALTER TABLE group_schemes ADD expected_phi_claims DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN expected_phi_claims DECIMAL(15,5);
END;

-- Add or modify column for field: ExpectedFunClaims
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'expected_fun_claims')
BEGIN
    ALTER TABLE group_schemes ADD expected_fun_claims DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN expected_fun_claims DECIMAL(15,5);
END;

-- Add or modify column for field: ExpectedClaims
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'expected_claims')
BEGIN
    ALTER TABLE group_schemes ADD expected_claims DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN expected_claims DECIMAL(15,5);
END;

-- Add or modify column for field: ActualClaims
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'actual_claims')
BEGIN
    ALTER TABLE group_schemes ADD actual_claims DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN actual_claims DECIMAL(15,5);
END;

-- Add or modify column for field: ExpectedClaimsRatio
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'expected_claims_ratio')
BEGIN
    ALTER TABLE group_schemes ADD expected_claims_ratio DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN expected_claims_ratio DECIMAL(15,5);
END;

-- Add or modify column for field: ActualClaimsRatio
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'actual_claims_ratio')
BEGIN
    ALTER TABLE group_schemes ADD actual_claims_ratio DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN actual_claims_ratio DECIMAL(15,5);
END;

-- Add or modify column for field: ExpectedLossRatio
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'expected_loss_ratio')
BEGIN
    ALTER TABLE group_schemes ADD expected_loss_ratio DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN expected_loss_ratio DECIMAL(15,5);
END;

-- Add or modify column for field: ActualLossRatio
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'actual_loss_ratio')
BEGIN
    ALTER TABLE group_schemes ADD actual_loss_ratio DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN actual_loss_ratio DECIMAL(15,5);
END;

-- Add or modify column for field: InForce
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'in_force')
BEGIN
    ALTER TABLE group_schemes ADD in_force BIT;
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN in_force BIT;
END;

-- Add or modify column for field: Status
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'status')
BEGIN
    ALTER TABLE group_schemes ADD status NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN status NVARCHAR(255);
END;

-- Add or modify column for field: NewBusiness
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'new_business')
BEGIN
    ALTER TABLE group_schemes ADD new_business BIT;
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN new_business BIT;
END;

-- Add or modify column for field: SchemeStatusMessage
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'scheme_status_message')
BEGIN
    ALTER TABLE group_schemes ADD scheme_status_message NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN scheme_status_message NVARCHAR(255);
END;

-- Add or modify column for field: CreationDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'creation_date')
BEGIN
    ALTER TABLE group_schemes ADD creation_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN creation_date DATETIME2;
END;

-- Add or modify column for field: CreatedBy
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'created_by')
BEGIN
    ALTER TABLE group_schemes ADD created_by NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN created_by NVARCHAR(255);
END;

-- Add or modify column for field: QuoteId
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'quote_id')
BEGIN
    ALTER TABLE group_schemes ADD quote_id INT;
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN quote_id INT;
END;

-- Add or modify column for field: Quote
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'quote')
BEGIN
    ALTER TABLE group_schemes ADD quote NVARCHAR(MAX);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN quote NVARCHAR(MAX);
END;

-- Add or modify column for field: QuoteInForce
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'quote_in_force')
BEGIN
    ALTER TABLE group_schemes ADD quote_in_force NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN quote_in_force NVARCHAR(255);
END;

-- Add or modify column for field: ActiveSchemeCategories
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'active_scheme_categories')
BEGIN
    ALTER TABLE group_schemes ADD active_scheme_categories json;
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN active_scheme_categories json;
END;

-- Add or modify column for field: CoverStartDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'cover_start_date')
BEGIN
    ALTER TABLE group_schemes ADD cover_start_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN cover_start_date DATETIME2;
END;

-- Add or modify column for field: CoverEndDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'cover_end_date')
BEGIN
    ALTER TABLE group_schemes ADD cover_end_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN cover_end_date DATETIME2;
END;

-- Add or modify column for field: CommencementDate
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'commencement_date')
BEGIN
    ALTER TABLE group_schemes ADD commencement_date DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN commencement_date DATETIME2;
END;

-- Add or modify column for field: SchemeQuoteStatus
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_schemes' AND COLUMN_NAME = 'scheme_quote_status')
BEGIN
    ALTER TABLE group_schemes ADD scheme_quote_status NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE group_schemes ALTER COLUMN scheme_quote_status NVARCHAR(255);
END;

