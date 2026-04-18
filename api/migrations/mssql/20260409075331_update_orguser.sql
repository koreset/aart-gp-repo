-- Migration for struct: OrgUser

-- Table: org_users

-- Ensure table exists
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'org_users')
BEGIN
    CREATE TABLE org_users (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

-- Add or modify column for field: Name
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'org_users' AND COLUMN_NAME = 'name')
BEGIN
    ALTER TABLE org_users ADD name NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE org_users ALTER COLUMN name NVARCHAR(255);
END;

-- Add or modify column for field: Email
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'org_users' AND COLUMN_NAME = 'email')
BEGIN
    ALTER TABLE org_users ADD email NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE org_users ALTER COLUMN email NVARCHAR(255);
END;

-- Add or modify column for field: LicenseId
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'org_users' AND COLUMN_NAME = 'license_id')
BEGIN
    ALTER TABLE org_users ADD license_id NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE org_users ALTER COLUMN license_id NVARCHAR(255);
END;

-- Add or modify column for field: Role
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'org_users' AND COLUMN_NAME = 'role')
BEGIN
    ALTER TABLE org_users ADD role NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE org_users ALTER COLUMN role NVARCHAR(255);
END;

-- Add or modify column for field: GPRole
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'org_users' AND COLUMN_NAME = 'gp_role')
BEGIN
    ALTER TABLE org_users ADD gp_role NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE org_users ALTER COLUMN gp_role NVARCHAR(255);
END;

-- Add or modify column for field: GPRoleId
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'org_users' AND COLUMN_NAME = 'gp_role_id')
BEGIN
    ALTER TABLE org_users ADD gp_role_id INT;
END;
ELSE
BEGIN
    ALTER TABLE org_users ALTER COLUMN gp_role_id INT;
END;

-- Add or modify column for field: Organisation
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'org_users' AND COLUMN_NAME = 'organisation')
BEGIN
    ALTER TABLE org_users ADD organisation NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE org_users ALTER COLUMN organisation NVARCHAR(255);
END;

