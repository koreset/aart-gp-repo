-- Migration for struct: ValUserRole

-- Table: val_user_roles

-- Ensure table exists
CREATE TABLE IF NOT EXISTS val_user_roles (
    id SERIAL PRIMARY KEY
);

-- Add or modify column for field: RoleName
ALTER TABLE val_user_roles ADD COLUMN IF NOT EXISTS role_name VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='val_user_roles' AND column_name='role_name') THEN
        ALTER TABLE val_user_roles ALTER COLUMN role_name TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Description
ALTER TABLE val_user_roles ADD COLUMN IF NOT EXISTS description VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='val_user_roles' AND column_name='description') THEN
        ALTER TABLE val_user_roles ALTER COLUMN description TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Permissions
ALTER TABLE val_user_roles ADD COLUMN IF NOT EXISTS permissions TEXT;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='val_user_roles' AND column_name='permissions') THEN
        ALTER TABLE val_user_roles ALTER COLUMN permissions TYPE TEXT;
    END IF;
END $$;

