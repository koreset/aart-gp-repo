-- Migration for struct: ValUserRole

-- Table: val_user_roles

-- Ensure table exists
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'val_user_roles')
BEGIN
    CREATE TABLE val_user_roles (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

-- Add or modify column for field: RoleName
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'val_user_roles' AND COLUMN_NAME = 'role_name')
BEGIN
    ALTER TABLE val_user_roles ADD role_name NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE val_user_roles ALTER COLUMN role_name NVARCHAR(255);
END;

-- Add or modify column for field: Description
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'val_user_roles' AND COLUMN_NAME = 'description')
BEGIN
    ALTER TABLE val_user_roles ADD description NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE val_user_roles ALTER COLUMN description NVARCHAR(255);
END;

-- Add or modify column for field: Permissions
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'val_user_roles' AND COLUMN_NAME = 'permissions')
BEGIN
    ALTER TABLE val_user_roles ADD permissions NVARCHAR(MAX);
END;
ELSE
BEGIN
    ALTER TABLE val_user_roles ALTER COLUMN permissions NVARCHAR(MAX);
END;

