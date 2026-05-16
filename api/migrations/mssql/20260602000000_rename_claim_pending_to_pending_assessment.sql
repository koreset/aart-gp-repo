-- Rename the legacy "pending" claim status to "pending_assessment" so the
-- claim lifecycle can distinguish drafts (capturer-controlled) from claims
-- that have been submitted for assessment. Only live status values on the
-- claims table are rewritten; status_audit rows preserve their historical
-- values so the audit timeline still reflects what was true at the time.
-- Idempotent on re-runs.

UPDATE group_scheme_claims
SET status = 'pending_assessment'
WHERE LOWER(status) = 'pending';
