-- Migration: create cpi_indices table.
--
-- Stores published Consumer Price Index values per (year, month). Used by
-- the Regular Income Claims tooling to apply benefit escalations and by
-- actuarial extracts that need historical CPI alongside claim history.
-- Uploads upsert on (year_index, month_index), so re-importing a file
-- updates existing rows in place.

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'cpi_indices')
BEGIN
    CREATE TABLE cpi_indices (
        id           BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT pk_cpi_indices PRIMARY KEY,
        year_index   INT NOT NULL,
        month_index  INT NOT NULL,
        cpi_index    FLOAT NOT NULL,
        created      DATETIMEOFFSET NULL CONSTRAINT df_cpi_indices_created DEFAULT (SYSDATETIMEOFFSET())
    );
END;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_cpi_indices_year_month')
    CREATE UNIQUE INDEX idx_cpi_indices_year_month
        ON cpi_indices(year_index, month_index);
