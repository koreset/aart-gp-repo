-- Migration: add ttd_benefit_escalation to scheme_categories so the user's
-- choice of TTD escalation option flows from the configuration UI into
-- GetTtdRate(), mirroring how phi_benefit_escalation already works.
-- Default '' (empty string) NOT NULL preserves current behaviour for every
-- existing scheme category row.

ALTER TABLE scheme_categories
    ADD COLUMN ttd_benefit_escalation VARCHAR(32) NOT NULL DEFAULT '';
