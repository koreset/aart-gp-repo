-- Migration: create cpi_indices table.
--
-- Stores published Consumer Price Index values per (year, month). Used by
-- the Regular Income Claims tooling to apply benefit escalations and by
-- actuarial extracts that need historical CPI alongside claim history.
-- Uploads upsert on (year_index, month_index), so re-importing a file
-- updates existing rows in place.

CREATE TABLE IF NOT EXISTS cpi_indices (
    id           BIGSERIAL PRIMARY KEY,
    year_index   INT NOT NULL,
    month_index  INT NOT NULL,
    cpi_index    DOUBLE PRECISION NOT NULL,
    created      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_cpi_indices_year_month
    ON cpi_indices(year_index, month_index);
