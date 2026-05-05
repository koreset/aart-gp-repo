-- Migration: add risk profile variation tolerance (in percent) to the
-- singleton group_pricing_settings row. Surfaces in the Acceptance Form text
-- on the generated quote PDF / DOCX. Default 7 preserves the historical
-- hardcoded "7%" wording.

ALTER TABLE group_pricing_settings ADD COLUMN risk_profile_variation_tolerance_pct DOUBLE NOT NULL DEFAULT 7;
