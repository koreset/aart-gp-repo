-- Migration: make bordereauxes.date_of_birth nullable and clear zero dates.
-- The Bordereaux.DateOfBirth field is now *time.Time so members without a
-- recorded DOB serialise to NULL.

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('bordereauxes') AND name = 'date_of_birth')
BEGIN
    ALTER TABLE bordereauxes ALTER COLUMN date_of_birth DATETIME2 NULL;
END

UPDATE bordereauxes
SET date_of_birth = NULL
WHERE date_of_birth = '0001-01-01';
