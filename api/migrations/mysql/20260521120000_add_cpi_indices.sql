-- Migration: create cpi_indices table.
--
-- Stores published Consumer Price Index values per (year, month). Used by
-- the Regular Income Claims tooling to apply benefit escalations and by
-- actuarial extracts that need historical CPI alongside claim history.
-- Uploads upsert on (year_index, month_index), so re-importing a file
-- updates existing rows in place.

CREATE TABLE IF NOT EXISTS cpi_indices (
    id           BIGINT AUTO_INCREMENT PRIMARY KEY,
    year_index   INT NOT NULL,
    month_index  INT NOT NULL,
    cpi_index    DOUBLE NOT NULL,
    created      DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3)
);

CREATE UNIQUE INDEX idx_cpi_indices_year_month
    ON cpi_indices(year_index, month_index);
