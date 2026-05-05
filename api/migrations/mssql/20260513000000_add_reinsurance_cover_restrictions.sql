-- Migration: create reinsurance_cover_restrictions, a per-(risk_rate_code,
-- scheme size band, benefit type) ceiling on the covered sum assured /
-- covered income for each member.
-- Lookup at calc time: rows where risk_rate_code matches the quote's code AND
-- min_scheme_size <= MemberCount AND (max_scheme_size = 0 OR MemberCount <= max_scheme_size).
-- max_scheme_size = 0 represents the open-ended top size band.
-- maximum_cover = 0 represents no restriction for that combination.

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'reinsurance_cover_restrictions')
BEGIN
    CREATE TABLE reinsurance_cover_restrictions (
        id              BIGINT       IDENTITY(1,1) NOT NULL PRIMARY KEY,
        risk_rate_code  NVARCHAR(255) NOT NULL CONSTRAINT df_rcr_risk_rate_code  DEFAULT '',
        benefit_type    NVARCHAR(32)  NOT NULL CONSTRAINT df_rcr_benefit_type    DEFAULT '',
        min_scheme_size BIGINT        NOT NULL CONSTRAINT df_rcr_min_scheme_size DEFAULT 0,
        max_scheme_size BIGINT        NOT NULL CONSTRAINT df_rcr_max_scheme_size DEFAULT 0,
        maximum_cover   FLOAT         NOT NULL CONSTRAINT df_rcr_maximum_cover   DEFAULT 0,
        creation_date   DATETIME2     NULL,
        created_by      NVARCHAR(255) NOT NULL CONSTRAINT df_rcr_created_by      DEFAULT ''
    );

    CREATE INDEX idx_rcr_lookup
        ON reinsurance_cover_restrictions (risk_rate_code, benefit_type, min_scheme_size, max_scheme_size);
END;
