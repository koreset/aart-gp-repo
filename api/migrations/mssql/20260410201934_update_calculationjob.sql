-- Migration for struct: CalculationJob

-- Table: calculation_jobs

-- Ensure table exists
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'calculation_jobs')
BEGIN
    CREATE TABLE calculation_jobs (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

-- Add or modify column for field: QuoteID
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'calculation_jobs' AND COLUMN_NAME = 'quote_id')
BEGIN
    ALTER TABLE calculation_jobs ADD quote_id INT;
END;
ELSE
BEGIN
    ALTER TABLE calculation_jobs ALTER COLUMN quote_id INT;
END;

-- Add or modify column for field: Basis
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'calculation_jobs' AND COLUMN_NAME = 'basis')
BEGIN
    ALTER TABLE calculation_jobs ADD basis NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE calculation_jobs ALTER COLUMN basis NVARCHAR(255);
END;

-- Add or modify column for field: Credibility
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'calculation_jobs' AND COLUMN_NAME = 'credibility')
BEGIN
    ALTER TABLE calculation_jobs ADD credibility DECIMAL(15,5);
END;
ELSE
BEGIN
    ALTER TABLE calculation_jobs ALTER COLUMN credibility DECIMAL(15,5);
END;

-- Add or modify column for field: UserEmail
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'calculation_jobs' AND COLUMN_NAME = 'user_email')
BEGIN
    ALTER TABLE calculation_jobs ADD user_email NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE calculation_jobs ALTER COLUMN user_email NVARCHAR(255);
END;

-- Add or modify column for field: UserName
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'calculation_jobs' AND COLUMN_NAME = 'user_name')
BEGIN
    ALTER TABLE calculation_jobs ADD user_name NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE calculation_jobs ALTER COLUMN user_name NVARCHAR(255);
END;

-- Add or modify column for field: Status
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'calculation_jobs' AND COLUMN_NAME = 'status')
BEGIN
    ALTER TABLE calculation_jobs ADD status NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE calculation_jobs ALTER COLUMN status NVARCHAR(255);
END;

-- Add or modify column for field: Error
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'calculation_jobs' AND COLUMN_NAME = 'error')
BEGIN
    ALTER TABLE calculation_jobs ADD error NVARCHAR(255);
END;
ELSE
BEGIN
    ALTER TABLE calculation_jobs ALTER COLUMN error NVARCHAR(255);
END;

-- Add or modify column for field: QueuedAt
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'calculation_jobs' AND COLUMN_NAME = 'queued_at')
BEGIN
    ALTER TABLE calculation_jobs ADD queued_at DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE calculation_jobs ALTER COLUMN queued_at DATETIME2;
END;

-- Add or modify column for field: StartedAt
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'calculation_jobs' AND COLUMN_NAME = 'started_at')
BEGIN
    ALTER TABLE calculation_jobs ADD started_at DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE calculation_jobs ALTER COLUMN started_at DATETIME2;
END;

-- Add or modify column for field: CompletedAt
-- SQL Server: Add column if it doesn't exist
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'calculation_jobs' AND COLUMN_NAME = 'completed_at')
BEGIN
    ALTER TABLE calculation_jobs ADD completed_at DATETIME2;
END;
ELSE
BEGIN
    ALTER TABLE calculation_jobs ALTER COLUMN completed_at DATETIME2;
END;

