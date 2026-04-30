-- Persist the manually-entered credibility (0-1) the actuary supplies
-- alongside experience-rate overrides. Saved to HistoricalCredibilityData as
-- ManuallyAddedCredibility on each calc run for future reference. Only
-- meaningful when experience_rating = 'Override'; left at 0 otherwise.

ALTER TABLE group_pricing_quotes
    ADD COLUMN experience_override_credibility DOUBLE NOT NULL DEFAULT 0;
