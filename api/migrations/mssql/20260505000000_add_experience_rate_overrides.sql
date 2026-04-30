-- Migration: per-(quote, scheme_category, benefit) experience-rate overrides.
-- Used when the quote's experience_rating field is set to 'Override': the
-- actuary supplies a mode (theoretical | experience_rated) and an
-- override_rate per (category, benefit) so the calculation engine bypasses
-- claims-based experience adjustment and uses the override rate directly.

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'group_pricing_experience_rate_overrides')
BEGIN
    CREATE TABLE group_pricing_experience_rate_overrides (
        id              INT IDENTITY(1,1) NOT NULL CONSTRAINT pk_gpero PRIMARY KEY,
        quote_id        INT NOT NULL,
        scheme_category NVARCHAR(128) NOT NULL,
        benefit         NVARCHAR(8) NOT NULL,
        mode            NVARCHAR(24) NOT NULL CONSTRAINT df_gpero_mode DEFAULT 'theoretical',
        override_rate   FLOAT NOT NULL CONSTRAINT df_gpero_override_rate DEFAULT 0,
        created_at      DATETIME2 NOT NULL CONSTRAINT df_gpero_created_at DEFAULT SYSUTCDATETIME(),
        created_by      NVARCHAR(128) NOT NULL CONSTRAINT df_gpero_created_by DEFAULT '',
        updated_at      DATETIME2 NOT NULL CONSTRAINT df_gpero_updated_at DEFAULT SYSUTCDATETIME(),
        updated_by      NVARCHAR(128) NOT NULL CONSTRAINT df_gpero_updated_by DEFAULT '',
        CONSTRAINT ux_group_pricing_exp_overrides UNIQUE (quote_id, scheme_category, benefit)
    );

    CREATE INDEX idx_gpero_quote ON group_pricing_experience_rate_overrides(quote_id);
END;
