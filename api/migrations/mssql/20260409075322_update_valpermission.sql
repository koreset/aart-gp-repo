-- Migration for struct: ValPermission

-- Table: val_permissions

-- Ensure table exists
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'val_permissions')
BEGIN
    CREATE TABLE val_permissions (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

-- Add or modify column for field: Name
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'val_permissions' AND COLUMN_NAME = 'name')
BEGIN
    ALTER TABLE val_permissions ADD name NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE val_permissions ALTER COLUMN name NVARCHAR(255);
END;

-- Add or modify column for field: Slug
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'val_permissions' AND COLUMN_NAME = 'slug')
BEGIN
    ALTER TABLE val_permissions ADD slug NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE val_permissions ALTER COLUMN slug NVARCHAR(255);
END;

-- Add or modify column for field: Description
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'val_permissions' AND COLUMN_NAME = 'description')
BEGIN
    ALTER TABLE val_permissions ADD description NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE val_permissions ALTER COLUMN description NVARCHAR(255);
END;

