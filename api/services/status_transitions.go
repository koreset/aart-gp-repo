package services

import (
	"fmt"
	"strings"
)

// Status constants grouped by entity. String-typed on purpose — keeping the
// field type as `string` on GORM models avoids a serialization change across
// the entire module. Use the consts when writing new code so typos surface at
// compile time; the transition maps below are the source of truth for which
// moves are legal.

// GeneratedBordereaux lifecycle.
const (
	StatusGeneratedDraft     = "draft"
	StatusGeneratedGenerated = "generated"
	StatusGeneratedReviewed  = "reviewed"
	StatusGeneratedApproved  = "approved"
	StatusGeneratedSubmitted = "submitted"
	StatusGeneratedConfirmed = "confirmed"
	StatusGeneratedCancelled = "cancelled"
	StatusGeneratedFailed    = "failed"
)

// BordereauxConfirmation lifecycle.
const (
	StatusConfirmationPending     = "pending"
	StatusConfirmationProcessed   = "processed"
	StatusConfirmationMatched     = "matched"
	StatusConfirmationDiscrepancy = "discrepancy"
	StatusConfirmationError       = "error"
)

// BordereauxReconciliationResult lifecycle.
const (
	StatusReconMatched     = "matched"
	StatusReconDiscrepancy = "discrepancy"
	StatusReconMissing     = "missing"
	StatusReconExtra       = "extra"
	StatusReconEscalated   = "escalated"
	StatusReconResolved    = "resolved"
)

// BordereauxDeadline lifecycle.
const (
	StatusDeadlinePending  = "pending"
	StatusDeadlineReceived = "received"
	StatusDeadlineOverdue  = "overdue"
	StatusDeadlineWaived   = "waived"
)

// EmployerSubmission lifecycle (documented inline in the model).
const (
	StatusSubmissionPendingReceipt = "pending_receipt"
	StatusSubmissionReceived       = "received"
	StatusSubmissionUnderReview    = "under_review"
	StatusSubmissionQueriesRaised  = "queries_raised"
	StatusSubmissionAccepted       = "accepted"
	StatusSubmissionRejected       = "rejected"
)

// LargeClaimNotice send/ack lifecycle.
const (
	StatusNoticePending      = "pending"
	StatusNoticeSent         = "sent"
	StatusNoticeAcknowledged = "acknowledged"
	StatusNoticeLate         = "late"
	StatusNoticeQueried      = "queried"
)

// LargeClaimNotice ResponseStatus — reinsurer decision field (separate from
// Status which tracks send/ack lifecycle).
const (
	ResponseNoticeAccepted = "accepted"
	ResponseNoticeRejected = "rejected"
)

// RIBordereauxRun lifecycle.
const (
	StatusRIRunDraft            = "draft"
	StatusRIRunGenerated        = "generated"
	StatusRIRunValidating       = "validating"
	StatusRIRunValidated        = "validated"
	StatusRIRunValidationFailed = "validation_failed"
	StatusRIRunSubmitted        = "submitted"
	StatusRIRunAcknowledged     = "acknowledged"
	StatusRIRunQueried          = "queried"
	StatusRIRunSettled          = "settled"
)

// ClaimNotificationLog lifecycle.
const (
	StatusClaimNotifPending      = "pending"
	StatusClaimNotifSent         = "sent"
	StatusClaimNotifAcknowledged = "acknowledged"
	StatusClaimNotifOverdue      = "overdue"
)

// Transition maps. Outer key = current status; inner set = legal next states.
// An empty "" current status means "new record being initialised" and is
// treated as always-allowed to any of the inner values.

var generatedBordereauxTransitions = map[string]map[string]bool{
	"": {StatusGeneratedDraft: true, StatusGeneratedGenerated: true, StatusGeneratedFailed: true},
	StatusGeneratedDraft: {
		StatusGeneratedGenerated: true,
		StatusGeneratedCancelled: true,
	},
	StatusGeneratedGenerated: {
		StatusGeneratedReviewed:  true,
		StatusGeneratedDraft:     true,
		StatusGeneratedCancelled: true,
		StatusGeneratedFailed:    true,
	},
	StatusGeneratedReviewed: {
		StatusGeneratedApproved:  true,
		StatusGeneratedDraft:     true,
		StatusGeneratedCancelled: true,
	},
	StatusGeneratedApproved: {
		StatusGeneratedSubmitted: true,
		StatusGeneratedDraft:     true,
	},
	StatusGeneratedSubmitted: {
		StatusGeneratedConfirmed: true,
		StatusGeneratedFailed:    true,
	},
	StatusGeneratedConfirmed: {},
	StatusGeneratedCancelled: {},
	StatusGeneratedFailed:    {StatusGeneratedDraft: true},
}

var reconciliationResultTransitions = map[string]map[string]bool{
	"": {
		StatusReconMatched:     true,
		StatusReconDiscrepancy: true,
		StatusReconMissing:     true,
		StatusReconExtra:       true,
	},
	StatusReconMatched:     {StatusReconDiscrepancy: true, StatusReconEscalated: true, StatusReconResolved: true},
	StatusReconDiscrepancy: {StatusReconEscalated: true, StatusReconResolved: true, StatusReconMatched: true},
	StatusReconMissing:     {StatusReconEscalated: true, StatusReconResolved: true},
	StatusReconExtra:       {StatusReconEscalated: true, StatusReconResolved: true},
	StatusReconEscalated:   {StatusReconResolved: true},
	StatusReconResolved:    {},
}

var deadlineTransitions = map[string]map[string]bool{
	"": {StatusDeadlinePending: true},
	StatusDeadlinePending: {
		StatusDeadlineReceived: true,
		StatusDeadlineOverdue:  true,
		StatusDeadlineWaived:   true,
	},
	StatusDeadlineOverdue: {
		StatusDeadlineReceived: true,
		StatusDeadlineWaived:   true,
		StatusDeadlinePending:  true,
	},
	StatusDeadlineReceived: {StatusDeadlinePending: true},
	StatusDeadlineWaived:   {StatusDeadlinePending: true},
}

var employerSubmissionTransitions = map[string]map[string]bool{
	"": {StatusSubmissionPendingReceipt: true, StatusSubmissionReceived: true},
	StatusSubmissionPendingReceipt: {StatusSubmissionReceived: true, StatusSubmissionRejected: true},
	StatusSubmissionReceived:       {StatusSubmissionUnderReview: true, StatusSubmissionRejected: true},
	StatusSubmissionUnderReview:    {StatusSubmissionQueriesRaised: true, StatusSubmissionAccepted: true, StatusSubmissionRejected: true},
	StatusSubmissionQueriesRaised:  {StatusSubmissionAccepted: true, StatusSubmissionRejected: true, StatusSubmissionUnderReview: true},
	StatusSubmissionAccepted:       {},
	StatusSubmissionRejected:       {},
}

var largeClaimNoticeTransitions = map[string]map[string]bool{
	"": {StatusNoticePending: true, StatusNoticeSent: true},
	StatusNoticePending: {
		StatusNoticeSent:    true,
		StatusNoticeLate:    true,
		StatusNoticeQueried: true,
	},
	StatusNoticeSent: {
		StatusNoticeAcknowledged: true,
		StatusNoticeLate:         true,
		StatusNoticeQueried:      true,
	},
	StatusNoticeAcknowledged: {StatusNoticeQueried: true},
	StatusNoticeLate:         {StatusNoticeAcknowledged: true, StatusNoticeQueried: true, StatusNoticeSent: true},
	StatusNoticeQueried:      {StatusNoticeAcknowledged: true, StatusNoticeSent: true},
}

var riBordereauxRunTransitions = map[string]map[string]bool{
	"": {StatusRIRunDraft: true, StatusRIRunGenerated: true},
	StatusRIRunDraft:      {StatusRIRunGenerated: true, StatusRIRunValidating: true},
	StatusRIRunGenerated:  {StatusRIRunValidating: true, StatusRIRunSubmitted: true, StatusRIRunDraft: true},
	StatusRIRunValidating: {StatusRIRunValidated: true, StatusRIRunValidationFailed: true},
	StatusRIRunValidated:  {StatusRIRunSubmitted: true, StatusRIRunValidating: true},
	StatusRIRunValidationFailed: {
		StatusRIRunValidating: true,
		StatusRIRunDraft:      true,
	},
	StatusRIRunSubmitted:    {StatusRIRunAcknowledged: true, StatusRIRunQueried: true},
	StatusRIRunAcknowledged: {StatusRIRunSettled: true, StatusRIRunQueried: true},
	StatusRIRunQueried:      {StatusRIRunAcknowledged: true, StatusRIRunSubmitted: true},
	StatusRIRunSettled:      {},
}

var claimNotificationLogTransitions = map[string]map[string]bool{
	"": {StatusClaimNotifPending: true, StatusClaimNotifSent: true},
	StatusClaimNotifPending: {
		StatusClaimNotifSent:    true,
		StatusClaimNotifOverdue: true,
	},
	StatusClaimNotifSent: {
		StatusClaimNotifAcknowledged: true,
		StatusClaimNotifOverdue:      true,
	},
	StatusClaimNotifOverdue:      {StatusClaimNotifSent: true, StatusClaimNotifAcknowledged: true},
	StatusClaimNotifAcknowledged: {},
}

// validateTransition is the shared primitive. Accepts the current status, the
// desired next status, and the transition map for an entity; returns a
// descriptive error if the move is illegal. Same-state "transitions"
// (from == to) are allowed to support idempotent setters.
func validateTransition(entity string, table map[string]map[string]bool, from, to string) error {
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)
	if to == "" {
		return fmt.Errorf("%s status cannot be empty", entity)
	}
	if from == to {
		return nil
	}
	allowed, ok := table[from]
	if !ok {
		return fmt.Errorf("%s has unknown current status %q", entity, from)
	}
	if !allowed[to] {
		return fmt.Errorf("%s cannot transition from %q to %q", entity, from, to)
	}
	return nil
}

// Per-entity wrappers. Controllers and services use these; they provide
// entity-specific error prefixes so callers know what failed.

func ValidateGeneratedBordereauxTransition(from, to string) error {
	return validateTransition("generated_bordereaux", generatedBordereauxTransitions, from, to)
}

func ValidateReconciliationResultTransition(from, to string) error {
	return validateTransition("reconciliation_result", reconciliationResultTransitions, from, to)
}

func ValidateDeadlineTransition(from, to string) error {
	return validateTransition("bordereaux_deadline", deadlineTransitions, from, to)
}

func ValidateEmployerSubmissionTransition(from, to string) error {
	return validateTransition("employer_submission", employerSubmissionTransitions, from, to)
}

func ValidateLargeClaimNoticeTransition(from, to string) error {
	return validateTransition("large_claim_notice", largeClaimNoticeTransitions, from, to)
}

// IsKnownLargeClaimNoticeStatus exposes the set of known statuses so the
// UpdateLargeClaimNotice handler can reject arbitrary user input without
// needing to know the current row state.
func IsKnownLargeClaimNoticeStatus(s string) bool {
	if _, ok := largeClaimNoticeTransitions[s]; ok {
		return true
	}
	return false
}

func ValidateRIBordereauxRunTransition(from, to string) error {
	return validateTransition("ri_bordereaux_run", riBordereauxRunTransitions, from, to)
}

func ValidateClaimNotificationLogTransition(from, to string) error {
	return validateTransition("claim_notification_log", claimNotificationLogTransitions, from, to)
}
