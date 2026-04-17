package bav

import (
	"context"
	"errors"
	"time"
)

// Provider is the contract every BAV adapter must satisfy. Synchronous
// providers (e.g. VerifyNow) return StatusComplete or StatusFailed from
// Verify directly and return ErrNotSupported from Poll. Async providers
// may return StatusPending with a ProviderJobID that the caller later
// resolves through Poll.
type Provider interface {
	// Name returns the short identifier used in configuration (e.g.
	// "verifynow", "lexisnexis"). Must match the BAV_PROVIDER config value.
	Name() string

	// Verify performs a bank account verification. Adapters translate the
	// canonical VerifyRequest to their provider's wire format, dispatch the
	// call, and normalise the response into a VerifyResult.
	Verify(ctx context.Context, req VerifyRequest) (*VerifyResult, error)

	// Poll resolves the current state of a previously-dispatched async
	// verification identified by jobID. Synchronous providers must return
	// ErrNotSupported. Async providers return a VerifyResult whose Status
	// is either StatusPending (still in flight), StatusComplete, or
	// StatusFailed.
	Poll(ctx context.Context, jobID string) (*VerifyResult, error)
}

// Registry owns the single active Provider selected at startup from config
// and, optionally, a Logger that persists each verification call. Callers
// invoke Registry.Verify rather than talking to the provider directly so
// audit logging and idempotency-based deduplication apply uniformly.
type Registry struct {
	active Provider
	logger Logger
}

// NewRegistry constructs a Registry backed by the supplied active provider.
// Passing a nil provider is valid — Active will return nil and verification
// calls will fail with ErrProviderNotConfigured. Logging is disabled by
// default; attach a Logger with WithLogger for audit + dedup.
func NewRegistry(active Provider) *Registry {
	return &Registry{active: active}
}

// WithLogger installs a Logger on the Registry. Returns the receiver so it
// can be chained after NewRegistry.
func (r *Registry) WithLogger(l Logger) *Registry {
	if r == nil {
		return r
	}
	r.logger = l
	return r
}

// Active returns the currently configured provider, or nil if none is wired.
func (r *Registry) Active() Provider {
	if r == nil {
		return nil
	}
	return r.active
}

// Verify runs a single verification through the active provider, transparently
// applying two cross-cutting behaviours:
//  1. Deterministic idempotency — if the caller didn't set IdempotencyKey, one
//     is derived from (claim_id, attempt, provider, identifying banking
//     fields). Identical repeat calls reuse the same key.
//  2. Audit logging + success dedup — when a Logger is attached, a previously
//     cached success within DedupeWindow short-circuits the provider call
//     and every attempt (hit, success, failure) is persisted.
//
// Failures are never cached: they always re-hit the provider so users can
// retry after a transient error. Logger errors are swallowed so an audit
// failure never breaks the user-facing call.
func (r *Registry) Verify(ctx context.Context, req VerifyRequest) (*VerifyResult, error) {
	if r == nil || r.active == nil {
		return nil, ErrProviderNotConfigured
	}

	providerName := r.active.Name()
	if req.IdempotencyKey == "" {
		req.IdempotencyKey = req.DeriveIdempotencyKey(providerName)
	}

	if r.logger != nil {
		if cached, err := r.logger.LookupSuccessByKey(ctx, req.IdempotencyKey, DedupeWindow); err == nil && cached != nil {
			return cached, nil
		}
	}

	result, callErr := r.active.Verify(ctx, req)

	if r.logger != nil {
		entry := LogEntry{
			ClaimID:        req.ClaimID,
			Provider:       providerName,
			IdempotencyKey: req.IdempotencyKey,
			Request:        req,
			Result:         result,
			Err:            callErr,
			CreatedAt:      time.Now().UTC(),
		}
		if result != nil {
			entry.Status = result.Status
			entry.ProviderRequestID = result.ProviderRequestID
		} else {
			entry.Status = StatusFailed
		}
		_ = r.logger.Record(ctx, entry)
	}

	return result, callErr
}

// Poll resolves a pending async verification. Callers must have obtained
// jobID from a previous Verify call whose result was StatusPending.
//
// Only terminal polls (complete, failed, or transport error) are audit-
// logged. Intermediate StatusPending polls are intentionally skipped so
// a single 60s verification doesn't produce 20 near-identical rows; the
// initial pending response written by Verify already marks the job's
// start, and this call logs its end.
func (r *Registry) Poll(ctx context.Context, jobID string) (*VerifyResult, error) {
	if r == nil || r.active == nil {
		return nil, ErrProviderNotConfigured
	}
	providerName := r.active.Name()

	result, callErr := r.active.Poll(ctx, jobID)

	if r.logger != nil && isTerminalPoll(result, callErr) {
		entry := LogEntry{
			Provider:      providerName,
			ProviderJobID: jobID,
			Result:        result,
			Err:           callErr,
			CreatedAt:     time.Now().UTC(),
		}
		if result != nil {
			entry.Status = result.Status
			entry.ProviderRequestID = result.ProviderRequestID
		} else {
			entry.Status = StatusFailed
		}
		_ = r.logger.Record(ctx, entry)
	}

	return result, callErr
}

// isTerminalPoll reports whether a Poll outcome warrants an audit row.
// Transport errors and final provider states (complete/failed) are
// terminal; StatusPending results are not.
func isTerminalPoll(result *VerifyResult, err error) bool {
	if err != nil {
		return true
	}
	if result == nil {
		return true
	}
	return result.Status != StatusPending
}

// Sentinel errors returned by adapters. Callers should match with errors.Is
// rather than string comparisons so adapters can wrap these with additional
// context.
var (
	ErrUnauthorized          = errors.New("bav: unauthorized")
	ErrRateLimited           = errors.New("bav: rate limited")
	ErrProviderUnavailable   = errors.New("bav: provider unavailable")
	ErrInvalidInput          = errors.New("bav: invalid input")
	ErrProviderNotConfigured = errors.New("bav: provider not configured")
	ErrNotSupported          = errors.New("bav: operation not supported by provider")
)
