// Package sanctions defines a slim adapter contract for sanctions / PEP /
// adverse-media screening. Phase 3 ships with a "manual" provider — finance
// records the outcome by hand — and the interface here is shaped so a real
// upstream (LexisNexis, Refinitiv, Dow Jones, etc.) can be plugged in later
// without touching the calling code.
//
// Calling convention is intentionally narrower than services/bav: a single
// synchronous Screen call, no async polling. Real providers that are
// inherently async should wrap their HTTP roundtrip in their adapter and
// return when they have a terminal answer or a clear pending state.
package sanctions

import (
	"context"
	"errors"
)

// Subject describes who is being screened. All real providers want at least
// a name; IdentityNumber improves match quality enormously when supplied.
type Subject struct {
	FullName       string
	IdentityNumber string
	DateOfBirth    string // ISO YYYY-MM-DD if known
	Country        string // ISO 3166-1 alpha-2
}

// Status enumerates the canonical screening outcomes. Mirrors the values
// stored on models.SanctionsScreening.Status.
type Status string

const (
	StatusPending     Status = "pending"
	StatusClear       Status = "clear"
	StatusHit         Status = "hit"
	StatusManualClear Status = "manual_clear"
	StatusSkipped     Status = "skipped"
)

// Result is the normalised provider response.
type Result struct {
	Status      Status
	ProviderRef string // upstream lookup id, for audit
	HitSummary  string // short human-readable description of matches
	// MatchCount is the count of distinct hits — useful for sanity-checking
	// "clear" results from providers whose default response is verbose.
	MatchCount int
}

// Provider is the contract a real sanctions adapter must satisfy.
type Provider interface {
	// Name returns the short identifier persisted on
	// models.SanctionsScreening.Provider (e.g. "manual", "lexisnexis").
	Name() string

	// Screen performs a single screening call. Synchronous providers return
	// StatusClear or StatusHit directly; the manual provider always returns
	// StatusPending so finance can record the outcome through Record.
	Screen(ctx context.Context, subj Subject) (*Result, error)
}

// Sentinel errors. Matched with errors.Is so adapters can wrap.
var (
	ErrProviderNotConfigured = errors.New("sanctions: provider not configured")
	ErrUnsupported           = errors.New("sanctions: operation not supported")
)

// active holds the currently registered provider. Phase 3 wires the manual
// provider at startup via Use(). When a real provider lands, it just calls
// Use() with its own implementation.
var active Provider

// Use installs the active provider. Safe to call once at startup; not safe
// for concurrent re-installation.
func Use(p Provider) {
	active = p
}

// Active returns the registered provider or nil if none is configured.
func Active() Provider {
	return active
}

// Screen is the package-level entry point. Callers (controllers, lifecycle
// services) use this so swapping the upstream is invisible to them.
func Screen(ctx context.Context, subj Subject) (*Result, error) {
	if active == nil {
		return nil, ErrProviderNotConfigured
	}
	return active.Screen(ctx, subj)
}
