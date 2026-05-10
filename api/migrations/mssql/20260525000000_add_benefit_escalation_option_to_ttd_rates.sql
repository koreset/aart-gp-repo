-- Migration: add benefit_escalation_option to ttd_rates so TTD rate lookups
-- can vary by escalation choice, mirroring the existing PHI rate pattern.
-- Default '' (empty string) NOT NULL backfills all existing rows so current
-- TTD lookups continue to match when the scheme category leaves the new
-- TtdBenefitEscalation field empty.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('ttd_rates') AND name = 'benefit_escalation_option')
BEGIN
    ALTER TABLE ttd_rates
        ADD benefit_escalation_option NVARCHAR(32) NOT NULL CONSTRAINT df_ttd_rates_benefit_escalation_option DEFAULT '';
END;
