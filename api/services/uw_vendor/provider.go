package uwvendor

import (
	"context"
	"errors"
)

// Provider is the contract every vendor adapter must satisfy. A provider
// is bound to a single Kind (a DocuSign adapter is an e-sign provider; an
// SMS gateway adapter is an SMS provider) and is registered once at boot.
//
// Synchronous adapters return StatusComplete from Submit directly; async
// ones return StatusAwaitingResponse and resolve via HandleWebhook (or via
// an internal Poll mechanism the adapter manages itself).
type Provider interface {
	// Name is the short identifier used in config and persisted on
	// VendorRequest.Provider. Must be unique per Kind.
	Name() string

	// Kind reports which service this provider implements.
	Kind() Kind

	// Submit dispatches the request. Adapters MUST set ExternalRequestID
	// on the result — webhooks reference it. The Registry takes care of
	// persisting the VendorRequest row before and after the call.
	Submit(ctx context.Context, req SubmitRequest) (*SubmitResult, error)
}

// WebhookCapable is an optional extension a Provider implements when the
// vendor delivers results via webhook. The Registry's Webhook dispatcher
// calls HandleWebhook with the verified, idempotent payload.
type WebhookCapable interface {
	Provider

	// VerifyWebhook validates the vendor's signature on an inbound payload.
	// Adapters check whatever HMAC header / shared secret the vendor uses.
	// Returning a non-nil error MUST prevent dispatch.
	VerifyWebhook(headers map[string][]string, body []byte) error

	// HandleWebhook decodes the payload into a WebhookOutcome. The
	// Registry uses the outcome to update the matching VendorRequest and
	// optionally persist an attachment on the originating case.
	HandleWebhook(ctx context.Context, body []byte) (*WebhookOutcome, error)
}

// Sentinel errors. Callers match with errors.Is so adapters can wrap
// these with vendor-specific context.
var (
	ErrProviderNotConfigured = errors.New("vendor: no provider configured for kind")
	ErrUnauthorized          = errors.New("vendor: unauthorized")
	ErrRateLimited           = errors.New("vendor: rate limited")
	ErrProviderUnavailable   = errors.New("vendor: provider unavailable")
	ErrInvalidInput          = errors.New("vendor: invalid input")
	ErrInvalidSignature      = errors.New("vendor: webhook signature invalid")
	ErrUnknownExternalID     = errors.New("vendor: unknown external request id")
	ErrAlreadyProcessed      = errors.New("vendor: webhook already processed")
)
