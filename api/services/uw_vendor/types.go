// Package vendor defines a single, kind-agnostic abstraction for outbound
// third-party service calls (pathology pulls, GP-record requests,
// e-signature envelopes, SMS sends). The framework mirrors services/bav
// (canonical types + Provider interface + Registry + audit logger) but
// uses one set of types parameterised by Kind so we don't end up with four
// copies of the same plumbing.
//
// Vendor-specific adapters live under services/vendor/providers/ and
// translate the canonical SubmitRequest/SubmitResult to and from each
// vendor's wire format.
package uwvendor

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// Kind identifies the service the request asks the vendor for. The string
// values are persisted on VendorRequest.Kind and exposed to the renderer.
type Kind string

const (
	KindPathology Kind = "pathology"
	KindGPRecords Kind = "gp_records"
	KindEsign     Kind = "esign"
	KindSMS       Kind = "sms"
)

// AllKinds is the canonical list used for validation and UI dropdowns.
var AllKinds = []Kind{KindPathology, KindGPRecords, KindEsign, KindSMS}

// IsValidKind reports whether k is a recognised Kind.
func IsValidKind(k Kind) bool {
	for _, v := range AllKinds {
		if v == k {
			return true
		}
	}
	return false
}

// Status describes the lifecycle of a vendor request.
type Status string

const (
	StatusQueued           Status = "queued"
	StatusInFlight         Status = "in_flight"
	StatusAwaitingResponse Status = "awaiting_response"
	StatusComplete         Status = "complete"
	StatusFailed           Status = "failed"
	StatusCancelled        Status = "cancelled"
)

// SubmitRequest is the canonical input every Provider.Submit consumes.
// Adapters translate it to the vendor's wire format.
//
//   - Subject is the addressee: an email address for e-sign, a phone for
//     SMS, a member ID for pathology / GP pulls. Mock and real providers
//     are free to interpret it as appropriate for their Kind.
//   - Body is the human-readable instruction or template text (e.g. the
//     SMS body, the pathology test panel name).
//   - Metadata carries per-Kind extras the canonical fields can't express;
//     it's serialised to VendorRequest.MetadataJSON for audit.
type SubmitRequest struct {
	Kind     Kind
	CaseID   int
	QuoteID  int
	Subject  string
	Body     string
	Metadata map[string]any
}

// SubmitResult is the canonical Provider.Submit response. ExternalRequestID
// is the vendor's identifier — webhooks carry it back so callers can find
// the originating VendorRequest row.
type SubmitResult struct {
	ExternalRequestID string
	Status            Status
	ProviderName      string
	CostCents         int
	RawPayload        json.RawMessage
}

// WebhookOutcome is what HandleWebhook returns to the Registry. The
// Registry uses it to update the matching VendorRequest, optionally
// attach a delivered file to the case, and emit a renderer-facing event.
type WebhookOutcome struct {
	ExternalRequestID string
	Status            Status
	Attachment        *AttachmentPayload
	Message           string
}

// AttachmentPayload represents a file delivered by the vendor (a pathology
// PDF, an e-signed consent form, a GP-records bundle). The Registry
// persists it as a UnderwritingCaseAttachment on the originating case.
type AttachmentPayload struct {
	Kind        string // medical_report | consent | disclosure | actively_at_work
	FileName    string
	ContentType string
	Body        []byte
}

// DeriveRequestPayloadHash produces a stable SHA-256 of the canonical
// request inputs so identical retries can be deduped against a cached
// success. Metadata is JSON-serialised in deterministic order so the hash
// is stable across re-invocations.
func DeriveRequestPayloadHash(req SubmitRequest) string {
	metaJSON := stableJSON(req.Metadata)
	h := sha256.New()
	_, _ = fmt.Fprintf(
		h,
		"%s|%d|%d|%s|%s|%s",
		req.Kind,
		req.CaseID,
		req.QuoteID,
		req.Subject,
		req.Body,
		metaJSON,
	)
	return hex.EncodeToString(h.Sum(nil))
}

// DeriveWebhookIdempotencyKey produces a SHA-256 of (provider + external
// request id + body) so re-deliveries of the same webhook are recognised
// even when the vendor's idempotency header is missing.
func DeriveWebhookIdempotencyKey(provider, externalID string, body []byte) string {
	h := sha256.New()
	_, _ = fmt.Fprintf(h, "%s|%s|", provider, externalID)
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

// stableJSON marshals a map with keys in sorted order so the same logical
// content always produces the same byte sequence — important for hashing.
func stableJSON(m map[string]any) string {
	if len(m) == 0 {
		return ""
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Inline insertion sort — input is tiny (per-request metadata).
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j-1] > keys[j]; j-- {
			keys[j-1], keys[j] = keys[j], keys[j-1]
		}
	}
	var b strings.Builder
	b.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			b.WriteByte(',')
		}
		kj, _ := json.Marshal(k)
		vj, _ := json.Marshal(m[k])
		b.Write(kj)
		b.WriteByte(':')
		b.Write(vj)
	}
	b.WriteByte('}')
	return b.String()
}
