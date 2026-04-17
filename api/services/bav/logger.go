package bav

import (
	"context"
	"time"
)

// DedupeWindow is how far back Registry.Verify looks for a cached successful
// verification with the same idempotency key. Matches the 24-hour requirement
// in the Phase 5 plan.
const DedupeWindow = 24 * time.Hour

// LogEntry is the payload Registry passes to Logger.Record for each
// verification attempt — successes and failures alike. ProviderJobID is
// populated for Poll entries; IdempotencyKey for Verify entries.
type LogEntry struct {
	ClaimID           *int
	Provider          string
	ProviderRequestID string
	ProviderJobID     string
	IdempotencyKey    string
	Status            Status
	Request           VerifyRequest
	Result            *VerifyResult
	Err               error
	CreatedAt         time.Time
}

// Logger persists verification attempts and serves dedup lookups. A nil
// Logger on Registry disables persistence and dedup entirely — useful for
// unit tests that don't stand up a database.
type Logger interface {
	// LookupSuccessByKey returns a previously-cached successful VerifyResult
	// for the given idempotency key if one was recorded within window. Only
	// successful calls are cached — failures never short-circuit retries.
	// Returns (nil, nil) on miss; (nil, err) on storage error.
	LookupSuccessByKey(ctx context.Context, key string, window time.Duration) (*VerifyResult, error)

	// Record persists a single verification attempt. Called with both
	// success and failure entries. Errors from Record should be logged by
	// the implementation but should not propagate, since a logging failure
	// must not break the user-facing verify call.
	Record(ctx context.Context, entry LogEntry) error
}
