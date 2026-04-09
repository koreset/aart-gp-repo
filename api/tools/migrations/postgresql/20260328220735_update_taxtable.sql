-- Migration for struct: TaxTable

-- Table: tax_tables

-- Ensure table exists
CREATE TABLE IF NOT EXISTS tax_tables (
    id SERIAL PRIMARY KEY
);

-- Add or modify column for field: RiskRateCode
ALTER TABLE tax_tables ADD COLUMN IF NOT EXISTS risk_rate_code VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='tax_tables' AND column_name='risk_rate_code') THEN
        ALTER TABLE tax_tables ALTER COLUMN risk_rate_code TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Level
ALTER TABLE tax_tables ADD COLUMN IF NOT EXISTS level INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='tax_tables' AND column_name='level') THEN
        ALTER TABLE tax_tables ALTER COLUMN level TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: Min
ALTER TABLE tax_tables ADD COLUMN IF NOT EXISTS min NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='tax_tables' AND column_name='min') THEN
        ALTER TABLE tax_tables ALTER COLUMN min TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: Max
ALTER TABLE tax_tables ADD COLUMN IF NOT EXISTS max NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='tax_tables' AND column_name='max') THEN
        ALTER TABLE tax_tables ALTER COLUMN max TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TaxRate
ALTER TABLE tax_tables ADD COLUMN IF NOT EXISTS tax_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='tax_tables' AND column_name='tax_rate') THEN
        ALTER TABLE tax_tables ALTER COLUMN tax_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CreationDate
ALTER TABLE tax_tables ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='tax_tables' AND column_name='creation_date') THEN
        ALTER TABLE tax_tables ALTER COLUMN creation_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: CreatedBy
ALTER TABLE tax_tables ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='tax_tables' AND column_name='created_by') THEN
        ALTER TABLE tax_tables ALTER COLUMN created_by TYPE VARCHAR(255);
    END IF;
END $$;

