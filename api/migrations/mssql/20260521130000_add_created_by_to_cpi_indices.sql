-- Migration: add created_by column to cpi_indices so the UI can show who
-- registered or uploaded each CPI value.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE Name = N'created_by' AND Object_ID = Object_ID(N'cpi_indices'))
    ALTER TABLE cpi_indices ADD created_by NVARCHAR(128) NOT NULL CONSTRAINT df_cpi_indices_created_by DEFAULT '';
