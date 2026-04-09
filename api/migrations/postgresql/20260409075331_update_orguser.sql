-- Migration for struct: OrgUser

-- Table: org_users

-- Ensure table exists
CREATE TABLE IF NOT EXISTS org_users (
    id SERIAL PRIMARY KEY
);

-- Add or modify column for field: Name
ALTER TABLE org_users ADD COLUMN IF NOT EXISTS name VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='org_users' AND column_name='name') THEN
        ALTER TABLE org_users ALTER COLUMN name TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Email
ALTER TABLE org_users ADD COLUMN IF NOT EXISTS email VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='org_users' AND column_name='email') THEN
        ALTER TABLE org_users ALTER COLUMN email TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: LicenseId
ALTER TABLE org_users ADD COLUMN IF NOT EXISTS license_id VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='org_users' AND column_name='license_id') THEN
        ALTER TABLE org_users ALTER COLUMN license_id TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Role
ALTER TABLE org_users ADD COLUMN IF NOT EXISTS role VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='org_users' AND column_name='role') THEN
        ALTER TABLE org_users ALTER COLUMN role TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: GPRole
ALTER TABLE org_users ADD COLUMN IF NOT EXISTS gp_role VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='org_users' AND column_name='gp_role') THEN
        ALTER TABLE org_users ALTER COLUMN gp_role TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: GPRoleId
ALTER TABLE org_users ADD COLUMN IF NOT EXISTS gp_role_id INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='org_users' AND column_name='gp_role_id') THEN
        ALTER TABLE org_users ALTER COLUMN gp_role_id TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: ValRole
ALTER TABLE org_users ADD COLUMN IF NOT EXISTS val_role VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='org_users' AND column_name='val_role') THEN
        ALTER TABLE org_users ALTER COLUMN val_role TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: ValRoleId
ALTER TABLE org_users ADD COLUMN IF NOT EXISTS val_role_id INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='org_users' AND column_name='val_role_id') THEN
        ALTER TABLE org_users ALTER COLUMN val_role_id TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: Organisation
ALTER TABLE org_users ADD COLUMN IF NOT EXISTS organisation VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='org_users' AND column_name='organisation') THEN
        ALTER TABLE org_users ALTER COLUMN organisation TYPE VARCHAR(255);
    END IF;
END $$;

