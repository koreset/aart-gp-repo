-- Migration for struct: SchemeSizeLevel

-- Table: scheme_size_levels

-- Ensure table exists
CREATE TABLE IF NOT EXISTS scheme_size_levels (
    id SERIAL PRIMARY KEY
);

-- Add or modify column for field: RiskRateCode
ALTER TABLE scheme_size_levels ADD COLUMN IF NOT EXISTS risk_rate_code VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_size_levels' AND column_name='risk_rate_code') THEN
        ALTER TABLE scheme_size_levels ALTER COLUMN risk_rate_code TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: MinCount
ALTER TABLE scheme_size_levels ADD COLUMN IF NOT EXISTS min_count INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_size_levels' AND column_name='min_count') THEN
        ALTER TABLE scheme_size_levels ALTER COLUMN min_count TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: MaxCount
ALTER TABLE scheme_size_levels ADD COLUMN IF NOT EXISTS max_count INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_size_levels' AND column_name='max_count') THEN
        ALTER TABLE scheme_size_levels ALTER COLUMN max_count TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: SizeLevel
ALTER TABLE scheme_size_levels ADD COLUMN IF NOT EXISTS size_level INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_size_levels' AND column_name='size_level') THEN
        ALTER TABLE scheme_size_levels ALTER COLUMN size_level TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: CreationDate
ALTER TABLE scheme_size_levels ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_size_levels' AND column_name='creation_date') THEN
        ALTER TABLE scheme_size_levels ALTER COLUMN creation_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: CreatedBy
ALTER TABLE scheme_size_levels ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='scheme_size_levels' AND column_name='created_by') THEN
        ALTER TABLE scheme_size_levels ALTER COLUMN created_by TYPE VARCHAR(255);
    END IF;
END $$;

