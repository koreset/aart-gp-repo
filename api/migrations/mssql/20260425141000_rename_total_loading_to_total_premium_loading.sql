-- Migration: align the member_rating_results column name with the GORM
-- struct field `TotalPremiumLoading` (default snake_case → total_premium_loading)
-- and the JSON/CSV tags. Earlier migrations created the column as
-- `total_loading`, which caused INSERTs to fail with
-- "Unknown column 'total_premium_loading' in 'field list'" on databases that
-- were provisioned via the manual migrations rather than GORM AutoMigrate.

-- Rename the legacy column if it exists and the canonical column does not.
IF EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'total_loading')
   AND NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'total_premium_loading')
BEGIN
    EXEC sp_rename 'member_rating_results.total_loading', 'total_premium_loading', 'COLUMN';
END;

-- Fallback: if neither column exists (e.g. a brand-new DB that skipped the
-- broken legacy migration), add the canonical column directly.
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'total_premium_loading')
BEGIN
    ALTER TABLE member_rating_results ADD total_premium_loading DECIMAL(15,5);
END;
