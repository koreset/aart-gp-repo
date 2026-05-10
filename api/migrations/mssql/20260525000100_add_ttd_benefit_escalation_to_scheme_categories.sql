-- Migration: add ttd_benefit_escalation to scheme_categories so the user's
-- choice of TTD escalation option flows from the configuration UI into
-- GetTtdRate(), mirroring how phi_benefit_escalation already works.
-- Default '' (empty string) NOT NULL preserves current behaviour for every
-- existing scheme category row.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'ttd_benefit_escalation')
BEGIN
    ALTER TABLE scheme_categories
        ADD ttd_benefit_escalation NVARCHAR(32) NOT NULL CONSTRAINT df_scheme_categories_ttd_benefit_escalation DEFAULT '';
END;
