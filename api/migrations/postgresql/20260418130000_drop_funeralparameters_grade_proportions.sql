-- Migration: drop legacy grade proportion columns from funeral_parameters.
-- Those fields belong to the Educator Rates table, not Funeral Parameters,
-- and were confusing users viewing the funeral_parameters CSV template.

ALTER TABLE funeral_parameters DROP COLUMN IF EXISTS grade0_proportion;
ALTER TABLE funeral_parameters DROP COLUMN IF EXISTS grade1_7_proportion;
ALTER TABLE funeral_parameters DROP COLUMN IF EXISTS grade_8_12_proportion;
ALTER TABLE funeral_parameters DROP COLUMN IF EXISTS tertiary_proportion;
