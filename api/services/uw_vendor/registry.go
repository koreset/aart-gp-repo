package uwvendor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"

	"api/models"
)

// Registry owns the active Provider per Kind. Boot wiring registers each
// one; controllers call Registry.Submit / Registry.IngestWebhook so the
// audit, idempotency and attachment-dispatch behaviours apply uniformly.
type Registry struct {
	mu        sync.RWMutex
	db        *gorm.DB
	providers map[Kind]Provider
}

// NewRegistry returns an empty Registry. db is used to persist
// VendorRequest / VendorWebhook rows.
func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{db: db, providers: make(map[Kind]Provider)}
}

// Register installs a provider for its Kind. Re-registering replaces any
// existing provider for the same Kind — useful for tests; in production
// each kind is registered exactly once at boot.
func (r *Registry) Register(p Provider) {
	if r == nil || p == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[p.Kind()] = p
}

// Active returns the provider configured for a kind, or nil if none.
func (r *Registry) Active(k Kind) Provider {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.providers[k]
}

// Kinds returns the kinds for which a provider is registered. Helpful for
// surfacing capability to the renderer ("don't show the pathology button
// if no pathology provider is wired").
func (r *Registry) Kinds() []Kind {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Kind, 0, len(r.providers))
	for k := range r.providers {
		out = append(out, k)
	}
	return out
}

// Submit dispatches a request through the registered provider, persists
// the VendorRequest audit row, and returns the persisted row to the
// caller.
//
// Behavioural contract:
//   - A VendorRequest row is created in queued status BEFORE the vendor
//     call so a crashed call leaves a forensic trail.
//   - On success the row is updated to in_flight / awaiting_response /
//     complete per the provider's response.
//   - On failure the row moves to failed and the error is recorded on
//     ErrorMessage so the renderer can surface it.
//   - Caller email is recorded on RequestedBy.
func (r *Registry) Submit(ctx context.Context, req SubmitRequest, actorEmail string) (*models.VendorRequest, error) {
	if r == nil {
		return nil, ErrProviderNotConfigured
	}
	if !IsValidKind(req.Kind) {
		return nil, ErrInvalidInput
	}
	p := r.Active(req.Kind)
	if p == nil {
		return nil, ErrProviderNotConfigured
	}

	metaJSON, _ := json.Marshal(req.Metadata)
	row := models.VendorRequest{
		Kind:               string(req.Kind),
		Provider:           p.Name(),
		CaseID:             req.CaseID,
		QuoteID:            req.QuoteID,
		Subject:            req.Subject,
		Body:               req.Body,
		MetadataJSON:       string(metaJSON),
		RequestPayloadHash: DeriveRequestPayloadHash(req),
		Status:             string(StatusQueued),
		RequestedBy:        actorEmail,
	}
	if err := r.db.Create(&row).Error; err != nil {
		return nil, fmt.Errorf("persist queued request: %w", err)
	}

	result, callErr := p.Submit(ctx, req)
	if callErr != nil {
		row.Status = string(StatusFailed)
		row.ErrorMessage = callErr.Error()
		_ = r.db.Save(&row).Error
		return &row, callErr
	}
	if result != nil {
		row.ExternalRequestID = result.ExternalRequestID
		row.Status = string(result.Status)
		row.CostCents = result.CostCents
		if len(result.RawPayload) > 0 {
			row.ResponseJSON = string(result.RawPayload)
		}
		if result.Status == StatusComplete {
			now := time.Now()
			row.CompletedAt = &now
		}
	}
	if err := r.db.Save(&row).Error; err != nil {
		return &row, fmt.Errorf("persist post-submit state: %w", err)
	}
	return &row, nil
}

// IngestWebhook persists the raw inbound payload, verifies the vendor's
// signature, then dispatches to the provider's HandleWebhook. Duplicate
// deliveries (same idempotency key) are recognised and short-circuited.
//
// Per-step semantics:
//  1. Persist VendorWebhook row in not-processed state so we can replay.
//  2. If the same idempotency key already exists and was processed, no-op
//     and return ErrAlreadyProcessed (controller maps to 200 OK with a
//     `replayed=true` flag).
//  3. VerifyWebhook → on failure mark webhook ProcessError and return.
//  4. HandleWebhook → update VendorRequest.Status and persist any
//     delivered attachment on the originating case.
func (r *Registry) IngestWebhook(ctx context.Context, kind Kind, providerName string, headers map[string][]string, body []byte) (*models.VendorWebhook, error) {
	if r == nil {
		return nil, ErrProviderNotConfigured
	}
	p := r.Active(kind)
	if p == nil || p.Name() != providerName {
		return nil, ErrProviderNotConfigured
	}
	wp, ok := p.(WebhookCapable)
	if !ok {
		return nil, fmt.Errorf("vendor: provider %s/%s does not accept webhooks", kind, providerName)
	}

	// Extract external request id from the body so it can be persisted on
	// the audit row even when verification fails (so we can debug).
	externalID := sniffExternalID(body)
	idempotencyKey := DeriveWebhookIdempotencyKey(providerName, externalID, body)

	row := models.VendorWebhook{
		Provider:          providerName,
		Kind:              string(kind),
		ExternalRequestID: externalID,
		IdempotencyKey:    idempotencyKey,
		BodySHA256:        sha256Hex(body),
		RawBody:           string(body),
		SignatureHeader:   firstHeader(headers, "X-Signature"),
	}

	var existing models.VendorWebhook
	if err := r.db.Where("idempotency_key = ?", idempotencyKey).First(&existing).Error; err == nil {
		if existing.Processed {
			return &existing, ErrAlreadyProcessed
		}
		// First persistence failed previously; reuse the existing row.
		row = existing
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := r.db.Create(&row).Error; err != nil {
			return nil, fmt.Errorf("persist webhook: %w", err)
		}
	} else {
		return nil, fmt.Errorf("lookup webhook: %w", err)
	}

	if err := wp.VerifyWebhook(headers, body); err != nil {
		row.ProcessError = err.Error()
		_ = r.db.Save(&row).Error
		return &row, fmt.Errorf("%w: %v", ErrInvalidSignature, err)
	}

	outcome, err := wp.HandleWebhook(ctx, body)
	if err != nil {
		row.ProcessError = err.Error()
		_ = r.db.Save(&row).Error
		return &row, err
	}

	if outcome != nil {
		if err := r.applyOutcome(*outcome, providerName, actorFromHeaders(headers)); err != nil {
			row.ProcessError = err.Error()
			_ = r.db.Save(&row).Error
			return &row, err
		}
	}

	now := time.Now()
	row.Processed = true
	row.ProcessedAt = &now
	if err := r.db.Save(&row).Error; err != nil {
		return &row, fmt.Errorf("mark processed: %w", err)
	}
	return &row, nil
}

// applyOutcome links the WebhookOutcome back to its originating
// VendorRequest, updates the status, and optionally persists an
// attachment on the case (medical reports, e-signed forms, etc.).
func (r *Registry) applyOutcome(outcome WebhookOutcome, providerName, actor string) error {
	if outcome.ExternalRequestID == "" {
		return ErrUnknownExternalID
	}
	var req models.VendorRequest
	if err := r.db.Where("external_request_id = ? AND provider = ?", outcome.ExternalRequestID, providerName).First(&req).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUnknownExternalID
		}
		return err
	}
	req.Status = string(outcome.Status)
	if outcome.Status == StatusComplete {
		now := time.Now()
		req.CompletedAt = &now
	}
	if outcome.Message != "" {
		req.ErrorMessage = ""
	}
	if err := r.db.Save(&req).Error; err != nil {
		return err
	}
	if outcome.Attachment != nil && req.CaseID > 0 {
		if err := persistAttachment(r.db, req.CaseID, outcome.Attachment, actor); err != nil {
			return fmt.Errorf("persist attachment: %w", err)
		}
	}
	return nil
}
