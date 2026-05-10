-- Migration: add benefit_escalation_option to ttd_rates so TTD rate lookups
-- can vary by escalation choice, mirroring the existing PHI rate pattern.
-- Default '' (empty string) NOT NULL backfills all existing rows so current
-- TTD lookups continue to match when the scheme category leaves the new
-- TtdBenefitEscalation field empty.

ALTER TABLE ttd_rates
    ADD COLUMN benefit_escalation_option VARCHAR(32) NOT NULL DEFAULT '';
