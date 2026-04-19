-- Migration: drop legacy grade proportion columns from funeral_parameters.
-- Those fields belong to the Educator Rates table, not Funeral Parameters,
-- and were confusing users viewing the funeral_parameters CSV template.

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('funeral_parameters') AND name = 'grade0_proportion')
BEGIN
    ALTER TABLE funeral_parameters DROP COLUMN grade0_proportion;
END

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('funeral_parameters') AND name = 'grade1_7_proportion')
BEGIN
    ALTER TABLE funeral_parameters DROP COLUMN grade1_7_proportion;
END

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('funeral_parameters') AND name = 'grade_8_12_proportion')
BEGIN
    ALTER TABLE funeral_parameters DROP COLUMN grade_8_12_proportion;
END

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('funeral_parameters') AND name = 'tertiary_proportion')
BEGIN
    ALTER TABLE funeral_parameters DROP COLUMN tertiary_proportion;
END
