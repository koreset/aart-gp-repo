-- Migration for struct: CalculationJob

-- Table: calculation_jobs

-- Ensure table exists
CREATE TABLE IF NOT EXISTS calculation_jobs (
    id SERIAL PRIMARY KEY
);

-- Add or modify column for field: QuoteID
ALTER TABLE calculation_jobs ADD COLUMN IF NOT EXISTS quote_id INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='calculation_jobs' AND column_name='quote_id') THEN
        ALTER TABLE calculation_jobs ALTER COLUMN quote_id TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: Basis
ALTER TABLE calculation_jobs ADD COLUMN IF NOT EXISTS basis VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='calculation_jobs' AND column_name='basis') THEN
        ALTER TABLE calculation_jobs ALTER COLUMN basis TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Credibility
ALTER TABLE calculation_jobs ADD COLUMN IF NOT EXISTS credibility NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='calculation_jobs' AND column_name='credibility') THEN
        ALTER TABLE calculation_jobs ALTER COLUMN credibility TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: UserEmail
ALTER TABLE calculation_jobs ADD COLUMN IF NOT EXISTS user_email VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='calculation_jobs' AND column_name='user_email') THEN
        ALTER TABLE calculation_jobs ALTER COLUMN user_email TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: UserName
ALTER TABLE calculation_jobs ADD COLUMN IF NOT EXISTS user_name VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='calculation_jobs' AND column_name='user_name') THEN
        ALTER TABLE calculation_jobs ALTER COLUMN user_name TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Status
ALTER TABLE calculation_jobs ADD COLUMN IF NOT EXISTS status VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='calculation_jobs' AND column_name='status') THEN
        ALTER TABLE calculation_jobs ALTER COLUMN status TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Error
ALTER TABLE calculation_jobs ADD COLUMN IF NOT EXISTS error VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='calculation_jobs' AND column_name='error') THEN
        ALTER TABLE calculation_jobs ALTER COLUMN error TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: QueuedAt
ALTER TABLE calculation_jobs ADD COLUMN IF NOT EXISTS queued_at TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='calculation_jobs' AND column_name='queued_at') THEN
        ALTER TABLE calculation_jobs ALTER COLUMN queued_at TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: StartedAt
ALTER TABLE calculation_jobs ADD COLUMN IF NOT EXISTS started_at TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='calculation_jobs' AND column_name='started_at') THEN
        ALTER TABLE calculation_jobs ALTER COLUMN started_at TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: CompletedAt
ALTER TABLE calculation_jobs ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='calculation_jobs' AND column_name='completed_at') THEN
        ALTER TABLE calculation_jobs ALTER COLUMN completed_at TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

