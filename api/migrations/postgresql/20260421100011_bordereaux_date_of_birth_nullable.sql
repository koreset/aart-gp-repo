-- Migration: make bordereauxes.date_of_birth nullable and clear zero dates.
-- The Bordereaux.DateOfBirth field is now *time.Time so members without a
-- recorded DOB serialise to NULL.

ALTER TABLE bordereauxes ALTER COLUMN date_of_birth DROP NOT NULL;

UPDATE bordereauxes
SET date_of_birth = NULL
WHERE date_of_birth = '0001-01-01 00:00:00';
