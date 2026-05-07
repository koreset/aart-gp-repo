-- Backfill commission_structures.holder_name to lower(trim(holder_name)) so
-- lookups no longer depend on case-insensitive collation. The Go service
-- layer now normalises holder_name on every read and write, mirroring the
-- existing channel handling.
--
-- Idempotent: re-running rewrites every row to the same already-normalised
-- value. Safe under the existing uniqueIndex (channel, holder_name,
-- lower_bound) — on a _ci collation, mixed-case duplicates could not have
-- been inserted in the first place.

UPDATE commission_structures
SET holder_name = LOWER(TRIM(holder_name));
