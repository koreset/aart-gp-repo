-- Audit columns for the Additional GLA Cover smoothed-rate save endpoint.
-- The smoothed values themselves are persisted inside the existing JSON
-- snapshot column (additional_gla_cover_band_rates); these two columns
-- record who last touched the smoothing for the category and when.

ALTER TABLE scheme_categories ADD COLUMN additional_gla_smoothed_updated_at DATETIME NULL;
ALTER TABLE scheme_categories ADD COLUMN additional_gla_smoothed_updated_by VARCHAR(255) NOT NULL DEFAULT '';
