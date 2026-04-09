-- Migration for struct: Broker

-- Table: brokers

-- Ensure table exists
CREATE TABLE IF NOT EXISTS brokers (
    id SERIAL PRIMARY KEY
);

-- Add or modify column for field: Name
ALTER TABLE brokers ADD COLUMN IF NOT EXISTS name VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='brokers' AND column_name='name') THEN
        ALTER TABLE brokers ALTER COLUMN name TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ContactEmail
ALTER TABLE brokers ADD COLUMN IF NOT EXISTS contact_email VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='brokers' AND column_name='contact_email') THEN
        ALTER TABLE brokers ALTER COLUMN contact_email TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ContactNumber
ALTER TABLE brokers ADD COLUMN IF NOT EXISTS contact_number VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='brokers' AND column_name='contact_number') THEN
        ALTER TABLE brokers ALTER COLUMN contact_number TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: FSPNumber
ALTER TABLE brokers ADD COLUMN IF NOT EXISTS fsp_number VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='brokers' AND column_name='fsp_number') THEN
        ALTER TABLE brokers ALTER COLUMN fsp_number TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: FSPCategory
ALTER TABLE brokers ADD COLUMN IF NOT EXISTS fsp_category VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='brokers' AND column_name='fsp_category') THEN
        ALTER TABLE brokers ALTER COLUMN fsp_category TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: BinderAgreementRef
ALTER TABLE brokers ADD COLUMN IF NOT EXISTS binder_agreement_ref VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='brokers' AND column_name='binder_agreement_ref') THEN
        ALTER TABLE brokers ALTER COLUMN binder_agreement_ref TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: TiedAgentRef
ALTER TABLE brokers ADD COLUMN IF NOT EXISTS tied_agent_ref VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='brokers' AND column_name='tied_agent_ref') THEN
        ALTER TABLE brokers ALTER COLUMN tied_agent_ref TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: CreationDate
ALTER TABLE brokers ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='brokers' AND column_name='creation_date') THEN
        ALTER TABLE brokers ALTER COLUMN creation_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: CreatedBy
ALTER TABLE brokers ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='brokers' AND column_name='created_by') THEN
        ALTER TABLE brokers ALTER COLUMN created_by TYPE VARCHAR(255);
    END IF;
END $$;

