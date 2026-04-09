-- Migration for struct: PremiumLoading

-- Table: premium_loadings

-- Ensure table exists
CREATE TABLE IF NOT EXISTS premium_loadings (
    id SERIAL PRIMARY KEY
);

-- Add or modify column for field: RiskRateCode
ALTER TABLE premium_loadings ADD COLUMN IF NOT EXISTS risk_rate_code VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='premium_loadings' AND column_name='risk_rate_code') THEN
        ALTER TABLE premium_loadings ALTER COLUMN risk_rate_code TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Channel
ALTER TABLE premium_loadings ADD COLUMN IF NOT EXISTS channel VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='premium_loadings' AND column_name='channel') THEN
        ALTER TABLE premium_loadings ALTER COLUMN channel TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: SchemeSizeLevel
ALTER TABLE premium_loadings ADD COLUMN IF NOT EXISTS scheme_size_level INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='premium_loadings' AND column_name='scheme_size_level') THEN
        ALTER TABLE premium_loadings ALTER COLUMN scheme_size_level TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: CommissionLoading
ALTER TABLE premium_loadings ADD COLUMN IF NOT EXISTS commission_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='premium_loadings' AND column_name='commission_loading') THEN
        ALTER TABLE premium_loadings ALTER COLUMN commission_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ExpenseLoading
ALTER TABLE premium_loadings ADD COLUMN IF NOT EXISTS expense_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='premium_loadings' AND column_name='expense_loading') THEN
        ALTER TABLE premium_loadings ALTER COLUMN expense_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: AdminLoading
ALTER TABLE premium_loadings ADD COLUMN IF NOT EXISTS admin_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='premium_loadings' AND column_name='admin_loading') THEN
        ALTER TABLE premium_loadings ALTER COLUMN admin_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: OtherLoading
ALTER TABLE premium_loadings ADD COLUMN IF NOT EXISTS other_loading NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='premium_loadings' AND column_name='other_loading') THEN
        ALTER TABLE premium_loadings ALTER COLUMN other_loading TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ProfitMargin
ALTER TABLE premium_loadings ADD COLUMN IF NOT EXISTS profit_margin NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='premium_loadings' AND column_name='profit_margin') THEN
        ALTER TABLE premium_loadings ALTER COLUMN profit_margin TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CreationDate
ALTER TABLE premium_loadings ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='premium_loadings' AND column_name='creation_date') THEN
        ALTER TABLE premium_loadings ALTER COLUMN creation_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: CreatedBy
ALTER TABLE premium_loadings ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='premium_loadings' AND column_name='created_by') THEN
        ALTER TABLE premium_loadings ALTER COLUMN created_by TYPE VARCHAR(255);
    END IF;
END $$;

