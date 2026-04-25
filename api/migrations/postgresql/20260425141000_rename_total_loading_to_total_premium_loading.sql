-- Migration: align the member_rating_results column name with the GORM
-- struct field `TotalPremiumLoading` (default snake_case → total_premium_loading)
-- and the JSON/CSV tags. Earlier migrations created the column as
-- `total_loading`, which caused INSERTs to fail with
-- "Unknown column 'total_premium_loading' in 'field list'" on databases that
-- were provisioned via the manual migrations rather than GORM AutoMigrate.

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'member_rating_results' AND column_name = 'total_premium_loading'
    ) THEN
        IF EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'member_rating_results' AND column_name = 'total_loading'
        ) THEN
            ALTER TABLE member_rating_results RENAME COLUMN total_loading TO total_premium_loading;
        ELSE
            ALTER TABLE member_rating_results ADD COLUMN total_premium_loading NUMERIC(15,5);
        END IF;
    END IF;
END $$;
