-- Persist the manually-entered credibility (0-1) the actuary supplies
-- alongside experience-rate overrides. Saved to HistoricalCredibilityData as
-- ManuallyAddedCredibility on each calc run for future reference. Only
-- meaningful when experience_rating = 'Override'; left at 0 otherwise.

IF NOT EXISTS (
    SELECT 1 FROM sys.columns
    WHERE object_id = OBJECT_ID('group_pricing_quotes')
      AND name = 'experience_override_credibility'
)
BEGIN
    ALTER TABLE group_pricing_quotes
        ADD experience_override_credibility FLOAT NOT NULL CONSTRAINT df_gpq_experience_override_credibility DEFAULT 0;
END;
