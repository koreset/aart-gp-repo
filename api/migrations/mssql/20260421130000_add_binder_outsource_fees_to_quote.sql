-- Migration: add binder fee and outsource fee columns to group_pricing_quotes.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'group_pricing_quotes')
BEGIN
    CREATE TABLE group_pricing_quotes (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'loadings_binder_fee')
    ALTER TABLE group_pricing_quotes ADD loadings_binder_fee FLOAT DEFAULT 0;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'group_pricing_quotes' AND COLUMN_NAME = 'loadings_outsource_fee')
    ALTER TABLE group_pricing_quotes ADD loadings_outsource_fee FLOAT DEFAULT 0;
