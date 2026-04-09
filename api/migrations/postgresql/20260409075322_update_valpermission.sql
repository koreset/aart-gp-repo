-- Migration for struct: ValPermission

-- Table: val_permissions

-- Ensure table exists
CREATE TABLE IF NOT EXISTS val_permissions (
    id SERIAL PRIMARY KEY
);

-- Add or modify column for field: Name
ALTER TABLE val_permissions ADD COLUMN IF NOT EXISTS name VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='val_permissions' AND column_name='name') THEN
        ALTER TABLE val_permissions ALTER COLUMN name TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Slug
ALTER TABLE val_permissions ADD COLUMN IF NOT EXISTS slug VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='val_permissions' AND column_name='slug') THEN
        ALTER TABLE val_permissions ALTER COLUMN slug TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Description
ALTER TABLE val_permissions ADD COLUMN IF NOT EXISTS description VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='val_permissions' AND column_name='description') THEN
        ALTER TABLE val_permissions ALTER COLUMN description TYPE VARCHAR(255);
    END IF;
END $$;

