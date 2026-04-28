-- Migration: singleton group_pricing_settings table for global group-pricing
-- toggles. Today it carries the discount calculation method (loading_adjustment
-- vs prorata) read by recomputeFinalPremiumsAndCommission and ApplyDiscountToQuote.

CREATE TABLE IF NOT EXISTS group_pricing_settings (
    id              BIGSERIAL PRIMARY KEY,
    discount_method VARCHAR(32) NOT NULL DEFAULT 'loading_adjustment',
    updated_at      TIMESTAMP,
    updated_by      VARCHAR(255)
);
