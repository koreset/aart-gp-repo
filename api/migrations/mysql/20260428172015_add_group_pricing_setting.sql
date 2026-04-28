-- Generated 2026-04-28T17:20:15+02:00 for dialect mysql

-- Create table: group_pricing_settings
CREATE TABLE group_pricing_settings (
    id BIGINT AUTO_INCREMENT,
    discount_method VARCHAR(32) NOT NULL DEFAULT 'loading_adjustment',
    updated_at DATETIME,
    updated_by VARCHAR(255),
    PRIMARY KEY (id)
);



