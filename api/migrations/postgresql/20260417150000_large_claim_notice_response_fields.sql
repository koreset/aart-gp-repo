-- Migration: add reinsurer response fields to large_claim_notices
-- Purpose: P2-5 — capture the reinsurer's decision (accepted / rejected) on a
-- large-claim cession, plus who responded and when. See
-- services/reinsurance_cession.go AcceptLargeClaimNotice / RejectLargeClaimNotice /
-- QueryLargeClaimNotice.

ALTER TABLE large_claim_notices ADD COLUMN IF NOT EXISTS response_status VARCHAR(32) NOT NULL DEFAULT '';
ALTER TABLE large_claim_notices ADD COLUMN IF NOT EXISTS responded_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE large_claim_notices ADD COLUMN IF NOT EXISTS responded_by VARCHAR(255) NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_lcn_response_status ON large_claim_notices (response_status);
