-- Backfill commission_structures.holder_name to lower(ltrim(rtrim(holder_name)))
-- so lookups no longer depend on case-insensitive collation. The Go service
-- layer now normalises holder_name on every read and write, mirroring the
-- existing channel handling.
--
-- Idempotent: re-running rewrites every row to the same already-normalised
-- value. LTRIM(RTRIM(...)) is used in place of TRIM() for compatibility with
-- SQL Server versions older than 2017.

UPDATE commission_structures
SET holder_name = LOWER(LTRIM(RTRIM(holder_name)));
