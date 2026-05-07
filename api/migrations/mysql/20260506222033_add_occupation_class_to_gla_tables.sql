-- Generated 2026-05-06T22:20:33+02:00 for dialect mysql

-- Migration for: GlaRate (table: gla_rates)

ALTER TABLE gla_rates ADD COLUMN occupation_class BIGINT;

-- Migration for: GlaAidsRate (table: gla_aids_rates)

ALTER TABLE gla_aids_rates ADD COLUMN occupation_class BIGINT;

-- Migration for: ReinsuranceGlaRate (table: reinsurance_gla_rates)

ALTER TABLE reinsurance_gla_rates ADD COLUMN occupation_class BIGINT;

-- Migration for: ReinsuranceGlaAidsRate (table: reinsurance_gla_aids_rates)

ALTER TABLE reinsurance_gla_aids_rates ADD COLUMN occupation_class BIGINT;

