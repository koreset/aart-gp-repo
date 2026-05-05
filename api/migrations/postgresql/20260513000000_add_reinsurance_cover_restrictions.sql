-- Migration: create reinsurance_cover_restrictions, a per-(risk_rate_code,
-- scheme size band, benefit type) ceiling on the covered sum assured /
-- covered income for each member.
-- Lookup at calc time: rows where risk_rate_code matches the quote's code AND
-- min_scheme_size <= MemberCount AND (max_scheme_size = 0 OR MemberCount <= max_scheme_size).
-- max_scheme_size = 0 represents the open-ended top size band.
-- maximum_cover = 0 represents no restriction for that combination.

CREATE TABLE IF NOT EXISTS reinsurance_cover_restrictions (
    id              BIGSERIAL PRIMARY KEY,
    risk_rate_code  VARCHAR(255) NOT NULL DEFAULT '',
    benefit_type    VARCHAR(32)  NOT NULL DEFAULT '',
    min_scheme_size BIGINT       NOT NULL DEFAULT 0,
    max_scheme_size BIGINT       NOT NULL DEFAULT 0,
    maximum_cover   DOUBLE PRECISION NOT NULL DEFAULT 0,
    creation_date   TIMESTAMP WITH TIME ZONE,
    created_by      VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_rcr_lookup
    ON reinsurance_cover_restrictions (risk_rate_code, benefit_type, min_scheme_size, max_scheme_size);
