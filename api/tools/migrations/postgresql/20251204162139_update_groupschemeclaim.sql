-- Migration for struct: GroupSchemeClaim

-- Table: group_scheme_claims

-- Ensure table exists
CREATE TABLE IF NOT EXISTS group_scheme_claims (
    id SERIAL PRIMARY KEY
);

-- Add or modify column for field: ClaimNumber
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS claim_number VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='claim_number') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN claim_number TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: MemberIDNumber
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS member_id_number VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='member_id_number') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN member_id_number TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: MemberName
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS member_name VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='member_name') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN member_name TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: SchemeId
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS scheme_id INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='scheme_id') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN scheme_id TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: SchemeName
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS scheme_name VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='scheme_name') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN scheme_name TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: BenefitType
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS benefit_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='benefit_type') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN benefit_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ClaimType
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS claim_type VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='claim_type') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN claim_type TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: DateOfEvent
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS date_of_event VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='date_of_event') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN date_of_event TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: DateNotified
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS date_notified VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='date_notified') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN date_notified TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ClaimAmount
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS claim_amount NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='claim_amount') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN claim_amount TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: Priority
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS priority VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='priority') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN priority TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ClaimantName
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS claimant_name VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='claimant_name') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN claimant_name TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ClaimantIDNumber
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS claimant_id_number VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='claimant_id_number') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN claimant_id_number TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: RelationshipToMember
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS relationship_to_member VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='relationship_to_member') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN relationship_to_member TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ClaimantContactNumber
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS claimant_contact_number VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='claimant_contact_number') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN claimant_contact_number TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Description
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS description VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='description') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN description TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Status
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS status VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='status') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN status TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: DateRegistered
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS date_registered VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='date_registered') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN date_registered TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: CreationDate
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='creation_date') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN creation_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: CreatedBy
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='created_by') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN created_by TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Attachments
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS attachments TEXT;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='attachments') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN attachments TYPE TEXT;
    END IF;
END $$;

-- Add or modify column for field: Assessments
ALTER TABLE group_scheme_claims ADD COLUMN IF NOT EXISTS assessments TEXT;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='group_scheme_claims' AND column_name='assessments') THEN
        ALTER TABLE group_scheme_claims ALTER COLUMN assessments TYPE TEXT;
    END IF;
END $$;

