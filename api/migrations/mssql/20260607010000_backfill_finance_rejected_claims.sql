-- One-shot backfill: surface historical finance rejections through the new
-- claim-side workflow. Before the finance-rejection feature shipped, both
-- query and reject paths returned the underlying claim to status='approved'
-- and recorded only an analytics row in claim_payment_schedule_queries.
-- This migration walks those rejected analytics rows and, for any claim
-- still sitting at 'approved' AND not on an active schedule, moves it to
-- 'finance_rejected' with the snapshot fields populated from the source
-- query row.
--
-- Idempotent: the WHERE filter on status = 'approved' means already-flagged
-- claims won't be touched a second time.

UPDATE gsc
SET status                            = 'finance_rejected',
    finance_rejected_at               = q.raised_at,
    finance_rejected_by               = q.raised_by,
    finance_rejection_reason_code     = q.reason_code,
    finance_rejection_notes           = q.notes,
    finance_rejection_schedule_number = s.schedule_number
FROM group_scheme_claims gsc
INNER JOIN (
    SELECT claim_id, MAX(raised_at) AS max_raised
    FROM claim_payment_schedule_queries
    WHERE LOWER(outcome) = 'rejected'
    GROUP BY claim_id
) latest ON latest.claim_id = gsc.id
INNER JOIN claim_payment_schedule_queries q
    ON q.claim_id = latest.claim_id
   AND q.raised_at = latest.max_raised
   AND LOWER(q.outcome) = 'rejected'
INNER JOIN claim_payment_schedules s ON s.id = q.schedule_id
WHERE LOWER(gsc.status) = 'approved'
  AND NOT EXISTS (
      SELECT 1
      FROM claim_payment_schedule_items i
      INNER JOIN claim_payment_schedules s2 ON s2.id = i.schedule_id
      WHERE i.claim_id = gsc.id
        AND i.line_status IN ('pending', 'verified')
        AND LOWER(s2.status) NOT IN ('archived', 'cancelled')
  );

INSERT INTO group_scheme_claim_status_audits (
    claim_id, old_status, new_status, status_message, changed_by, changed_at
)
SELECT gsc.id,
       'approved',
       'finance_rejected',
       'Backfill: historical finance rejection on ' + ISNULL(gsc.finance_rejection_schedule_number, '?') + ' — ' + ISNULL(gsc.finance_rejection_reason_code, ''),
       'system_backfill',
       SYSUTCDATETIME()
FROM group_scheme_claims gsc
WHERE LOWER(gsc.status) = 'finance_rejected'
  AND gsc.finance_rejected_by IS NOT NULL
  AND NOT EXISTS (
      SELECT 1 FROM group_scheme_claim_status_audits a
      WHERE a.claim_id = gsc.id
        AND a.changed_by = 'system_backfill'
        AND a.new_status = 'finance_rejected'
  );
