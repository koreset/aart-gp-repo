-- Generated 2026-05-06T22:41:33+02:00 for dialect mysql

-- Migration for: IndustryLoading (table: industry_loadings)

ALTER TABLE industry_loadings ADD COLUMN income_level BIGINT;

-- Migration for: RegionLoading (table: region_loadings)

ALTER TABLE region_loadings ADD COLUMN income_level BIGINT;
ALTER TABLE region_loadings ADD COLUMN occupation_class BIGINT;

-- Migration for: ReinsuranceIndustryLoading (table: reinsurance_industry_loadings)

ALTER TABLE reinsurance_industry_loadings ADD COLUMN income_level BIGINT;

-- Migration for: ReinsuranceRegionLoading (table: reinsurance_region_loadings)

ALTER TABLE reinsurance_region_loadings ADD COLUMN income_level BIGINT;
ALTER TABLE reinsurance_region_loadings ADD COLUMN occupation_class BIGINT;

