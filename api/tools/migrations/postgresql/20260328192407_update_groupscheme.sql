-- Migration for struct: GroupScheme

-- Table: group_schemes

-- Ensure table exists
CREATE TABLE IF NOT EXISTS group_schemes (
    id SERIAL PRIMARY KEY
);

-- Add or modify column for field: Name
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS name VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='name') THEN
        ALTER TABLE group_schemes ALTER COLUMN name TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: BrokerId
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS broker_id INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='broker_id') THEN
        ALTER TABLE group_schemes ALTER COLUMN broker_id TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: DistributionChannel
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS distribution_channel VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='distribution_channel') THEN
        ALTER TABLE group_schemes ALTER COLUMN distribution_channel TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ContactPerson
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS contact_person VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='contact_person') THEN
        ALTER TABLE group_schemes ALTER COLUMN contact_person TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ContactEmail
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS contact_email VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='contact_email') THEN
        ALTER TABLE group_schemes ALTER COLUMN contact_email TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: DurationInForceDays
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS duration_in_force_days INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='duration_in_force_days') THEN
        ALTER TABLE group_schemes ALTER COLUMN duration_in_force_days TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: RenewalDate
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS renewal_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='renewal_date') THEN
        ALTER TABLE group_schemes ALTER COLUMN renewal_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: MemberCount
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS member_count NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='member_count') THEN
        ALTER TABLE group_schemes ALTER COLUMN member_count TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: AnnualPremium
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS annual_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='annual_premium') THEN
        ALTER TABLE group_schemes ALTER COLUMN annual_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: GlaAnnualPremium
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS gla_annual_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='gla_annual_premium') THEN
        ALTER TABLE group_schemes ALTER COLUMN gla_annual_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdAnnualPremium
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS ptd_annual_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='ptd_annual_premium') THEN
        ALTER TABLE group_schemes ALTER COLUMN ptd_annual_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiAnnualPremium
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS ci_annual_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='ci_annual_premium') THEN
        ALTER TABLE group_schemes ALTER COLUMN ci_annual_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: SglaAnnualPremium
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS sgla_annual_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='sgla_annual_premium') THEN
        ALTER TABLE group_schemes ALTER COLUMN sgla_annual_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdAnnualPremium
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS ttd_annual_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='ttd_annual_premium') THEN
        ALTER TABLE group_schemes ALTER COLUMN ttd_annual_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiAnnualPremium
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS phi_annual_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='phi_annual_premium') THEN
        ALTER TABLE group_schemes ALTER COLUMN phi_annual_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: FuneralAnnualPremium
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS funeral_annual_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='funeral_annual_premium') THEN
        ALTER TABLE group_schemes ALTER COLUMN funeral_annual_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: Commission
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS commission NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='commission') THEN
        ALTER TABLE group_schemes ALTER COLUMN commission TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: EarnedPremium
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS earned_premium NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='earned_premium') THEN
        ALTER TABLE group_schemes ALTER COLUMN earned_premium TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpectedExpenses
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS expected_expenses NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='expected_expenses') THEN
        ALTER TABLE group_schemes ALTER COLUMN expected_expenses TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpectedGlaClaims
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS expected_gla_claims NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='expected_gla_claims') THEN
        ALTER TABLE group_schemes ALTER COLUMN expected_gla_claims TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpectedPtdClaims
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS expected_ptd_claims NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='expected_ptd_claims') THEN
        ALTER TABLE group_schemes ALTER COLUMN expected_ptd_claims TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpectedCiClaims
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS expected_ci_claims NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='expected_ci_claims') THEN
        ALTER TABLE group_schemes ALTER COLUMN expected_ci_claims TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpectedSglaClaims
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS expected_sgla_claims NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='expected_sgla_claims') THEN
        ALTER TABLE group_schemes ALTER COLUMN expected_sgla_claims TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpectedTtdClaims
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS expected_ttd_claims NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='expected_ttd_claims') THEN
        ALTER TABLE group_schemes ALTER COLUMN expected_ttd_claims TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpectedPhiClaims
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS expected_phi_claims NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='expected_phi_claims') THEN
        ALTER TABLE group_schemes ALTER COLUMN expected_phi_claims TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpectedFunClaims
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS expected_fun_claims NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='expected_fun_claims') THEN
        ALTER TABLE group_schemes ALTER COLUMN expected_fun_claims TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpectedClaims
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS expected_claims NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='expected_claims') THEN
        ALTER TABLE group_schemes ALTER COLUMN expected_claims TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ActualClaims
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS actual_claims NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='actual_claims') THEN
        ALTER TABLE group_schemes ALTER COLUMN actual_claims TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpectedClaimsRatio
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS expected_claims_ratio NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='expected_claims_ratio') THEN
        ALTER TABLE group_schemes ALTER COLUMN expected_claims_ratio TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ActualClaimsRatio
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS actual_claims_ratio NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='actual_claims_ratio') THEN
        ALTER TABLE group_schemes ALTER COLUMN actual_claims_ratio TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpectedLossRatio
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS expected_loss_ratio NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='expected_loss_ratio') THEN
        ALTER TABLE group_schemes ALTER COLUMN expected_loss_ratio TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ActualLossRatio
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS actual_loss_ratio NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='actual_loss_ratio') THEN
        ALTER TABLE group_schemes ALTER COLUMN actual_loss_ratio TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: InForce
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS in_force BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='in_force') THEN
        ALTER TABLE group_schemes ALTER COLUMN in_force TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: Status
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS status VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='status') THEN
        ALTER TABLE group_schemes ALTER COLUMN status TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: NewBusiness
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS new_business BOOLEAN;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='new_business') THEN
        ALTER TABLE group_schemes ALTER COLUMN new_business TYPE BOOLEAN;
    END IF;
END $$;

-- Add or modify column for field: SchemeStatusMessage
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS scheme_status_message VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='scheme_status_message') THEN
        ALTER TABLE group_schemes ALTER COLUMN scheme_status_message TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: CreationDate
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='creation_date') THEN
        ALTER TABLE group_schemes ALTER COLUMN creation_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: CreatedBy
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='created_by') THEN
        ALTER TABLE group_schemes ALTER COLUMN created_by TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: QuoteId
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS quote_id INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='quote_id') THEN
        ALTER TABLE group_schemes ALTER COLUMN quote_id TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: Quote
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS quote TEXT;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='quote') THEN
        ALTER TABLE group_schemes ALTER COLUMN quote TYPE TEXT;
    END IF;
END $$;

-- Add or modify column for field: QuoteInForce
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS quote_in_force VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='quote_in_force') THEN
        ALTER TABLE group_schemes ALTER COLUMN quote_in_force TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ActiveSchemeCategories
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS active_scheme_categories json;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='active_scheme_categories') THEN
        ALTER TABLE group_schemes ALTER COLUMN active_scheme_categories TYPE json;
    END IF;
END $$;

-- Add or modify column for field: CoverStartDate
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS cover_start_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='cover_start_date') THEN
        ALTER TABLE group_schemes ALTER COLUMN cover_start_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: CoverEndDate
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS cover_end_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='cover_end_date') THEN
        ALTER TABLE group_schemes ALTER COLUMN cover_end_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: CommencementDate
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS commencement_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='commencement_date') THEN
        ALTER TABLE group_schemes ALTER COLUMN commencement_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: SchemeQuoteStatus
ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS scheme_quote_status VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_schemes' AND column_name='scheme_quote_status') THEN
        ALTER TABLE group_schemes ALTER COLUMN scheme_quote_status TYPE VARCHAR(255);
    END IF;
END $$;

