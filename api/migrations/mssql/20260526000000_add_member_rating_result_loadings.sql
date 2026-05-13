-- Add per-member scheme-size and educator base loading columns to
-- member_rating_results. Mirrors the source-side columns added in
-- 20260518000000_add_scheme_size_loadings.sql and
-- 20260518010000_add_educator_loadings.sql.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'gla_scheme_size_loading')
BEGIN
    ALTER TABLE member_rating_results
        ADD gla_scheme_size_loading FLOAT NOT NULL CONSTRAINT df_mrr_gla_scheme_size_loading DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'ptd_scheme_size_loading')
BEGIN
    ALTER TABLE member_rating_results
        ADD ptd_scheme_size_loading FLOAT NOT NULL CONSTRAINT df_mrr_ptd_scheme_size_loading DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'ci_scheme_size_loading')
BEGIN
    ALTER TABLE member_rating_results
        ADD ci_scheme_size_loading FLOAT NOT NULL CONSTRAINT df_mrr_ci_scheme_size_loading DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'ttd_scheme_size_loading')
BEGIN
    ALTER TABLE member_rating_results
        ADD ttd_scheme_size_loading FLOAT NOT NULL CONSTRAINT df_mrr_ttd_scheme_size_loading DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'phi_scheme_size_loading')
BEGIN
    ALTER TABLE member_rating_results
        ADD phi_scheme_size_loading FLOAT NOT NULL CONSTRAINT df_mrr_phi_scheme_size_loading DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'fun_scheme_size_loading')
BEGIN
    ALTER TABLE member_rating_results
        ADD fun_scheme_size_loading FLOAT NOT NULL CONSTRAINT df_mrr_fun_scheme_size_loading DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'gla_educator_loading')
BEGIN
    ALTER TABLE member_rating_results
        ADD gla_educator_loading FLOAT NOT NULL CONSTRAINT df_mrr_gla_educator_loading DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'ptd_educator_loading')
BEGIN
    ALTER TABLE member_rating_results
        ADD ptd_educator_loading FLOAT NOT NULL CONSTRAINT df_mrr_ptd_educator_loading DEFAULT 0;
END;
