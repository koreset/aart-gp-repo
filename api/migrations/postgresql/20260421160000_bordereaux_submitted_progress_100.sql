-- Migration: submission is now the final step in the bordereaux cycle, so
-- already-submitted rows should reflect 100% progress (previously 66%).
-- Idempotent: no-op after the first run since no submitted row can drop below 100.

UPDATE generated_bordereauxes
SET progress = 100
WHERE status = 'submitted' AND progress < 100;
