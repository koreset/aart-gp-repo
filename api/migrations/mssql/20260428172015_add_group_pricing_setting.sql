-- Migration: singleton group_pricing_settings table for global group-pricing
-- toggles. Today it carries the discount calculation method (loading_adjustment
-- vs prorata) read by recomputeFinalPremiumsAndCommission and ApplyDiscountToQuote.

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'group_pricing_settings')
BEGIN
    CREATE TABLE group_pricing_settings (
        id              BIGINT IDENTITY(1,1) PRIMARY KEY,
        discount_method NVARCHAR(32) NOT NULL CONSTRAINT df_group_pricing_settings_method DEFAULT 'loading_adjustment',
        updated_at      DATETIME2,
        updated_by      NVARCHAR(255)
    );
END;
