-- Audit columns for the Additional GLA Cover smoothed-rate save endpoint.
-- The smoothed values themselves are persisted inside the existing JSON
-- snapshot column (additional_gla_cover_band_rates); these two columns
-- record who last touched the smoothing for the category and when.

IF NOT EXISTS (
    SELECT 1 FROM sys.columns
    WHERE object_id = OBJECT_ID('scheme_categories')
      AND name = 'additional_gla_smoothed_updated_at'
)
BEGIN
    ALTER TABLE scheme_categories
        ADD additional_gla_smoothed_updated_at DATETIME2 NULL;
END;

IF NOT EXISTS (
    SELECT 1 FROM sys.columns
    WHERE object_id = OBJECT_ID('scheme_categories')
      AND name = 'additional_gla_smoothed_updated_by'
)
BEGIN
    ALTER TABLE scheme_categories
        ADD additional_gla_smoothed_updated_by NVARCHAR(255) NOT NULL
            CONSTRAINT df_sc_agla_smoothed_updated_by DEFAULT '';
END;
