// Package sms defines a slim adapter contract for outbound SMS. Phase 4
// ships with a "log" provider — every send is appended to the app log and
// returned successful — so claimant-notification logic can run end-to-end
// without an upstream SMS account. The interface is shaped so a real
// provider (Clickatell, Twilio, BulkSMS, etc.) drops in via Use() without
// changes to the calling code.
package sms

import (
	"context"
	"errors"
)

// Message describes a single outbound SMS.
type Message struct {
	To   string // E.164 if possible; the provider normalises
	Body string
	// Reference is an arbitrary correlation id (e.g. claim number) that the
	// adapter may pass through to the upstream so the message is searchable
	// in their dashboard.
	Reference string
}

// Result is the normalised provider response. ProviderRef is the upstream
// message id when one is returned (search key in the provider dashboard).
type Result struct {
	ProviderRef string
	Status      string // "queued" | "delivered" | "failed"
}

// Provider is the contract every SMS adapter must satisfy.
type Provider interface {
	// Name returns the short identifier persisted on the SMS audit row.
	Name() string

	// Send dispatches a single message synchronously. Adapters that talk to
	// async upstreams should wrap the upstream's pending state and return
	// Status "queued" with the upstream id; delivery is observed via a
	// separate webhook outside this interface.
	Send(ctx context.Context, msg Message) (*Result, error)
}

// Sentinel errors — matched with errors.Is so adapters can wrap.
var (
	ErrProviderNotConfigured = errors.New("sms: provider not configured")
	ErrInvalidRecipient      = errors.New("sms: invalid recipient")
)

// active holds the currently registered provider.
var active Provider

// Use installs the active provider. Safe to call once at startup; not safe
// for concurrent re-installation.
func Use(p Provider) { active = p }

// Active returns the currently registered provider, or nil.
func Active() Provider { return active }

// Send is the package-level entry point.
func Send(ctx context.Context, msg Message) (*Result, error) {
	if active == nil {
		return nil, ErrProviderNotConfigured
	}
	return active.Send(ctx, msg)
}
