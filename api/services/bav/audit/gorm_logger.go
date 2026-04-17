// Package audit provides a GORM-backed implementation of bav.Logger that
// persists each bank account verification call to the
// bav_verification_logs table.
package audit

import (
	"api/models"
	"api/services/bav"
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

// GormLogger persists bav verification attempts to the database and serves
// the 24-hour dedup lookup via the idempotency key index.
type GormLogger struct {
	db *gorm.DB
}

// NewGormLogger wires a GORM-backed bav.Logger. db is the primary
// application *gorm.DB handle.
func NewGormLogger(db *gorm.DB) *GormLogger {
	return &GormLogger{db: db}
}

// LookupSuccessByKey returns the most recent successful verification for
// the given idempotency key within the supplied window, or (nil, nil) on
// miss. The returned VerifyResult is reconstituted from the stored
// response payload.
func (g *GormLogger) LookupSuccessByKey(ctx context.Context, key string, window time.Duration) (*bav.VerifyResult, error) {
	if g == nil || g.db == nil || key == "" {
		return nil, nil
	}

	var row models.BAVVerificationLog
	cutoff := time.Now().UTC().Add(-window)
	err := g.db.WithContext(ctx).
		Where("idempotency_key = ? AND status = ? AND created_at >= ?", key, string(bav.StatusComplete), cutoff).
		Order("created_at DESC").
		First(&row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if len(row.ResponsePayload) == 0 {
		return nil, nil
	}
	return unmarshalStoredResult(row.ResponsePayload)
}

// Record writes a single verification attempt. Any storage error is logged
// locally but not propagated — an audit failure must not break the user's
// verify call.
func (g *GormLogger) Record(ctx context.Context, entry bav.LogEntry) error {
	if g == nil || g.db == nil {
		return nil
	}

	reqBytes, _ := json.Marshal(entry.Request)
	var respBytes []byte
	if entry.Result != nil {
		respBytes, _ = marshalResultForStorage(entry.Result)
	}

	errMsg := ""
	if entry.Err != nil {
		errMsg = entry.Err.Error()
		if len(errMsg) > 1024 {
			errMsg = errMsg[:1024]
		}
	}

	row := models.BAVVerificationLog{
		ClaimID:           entry.ClaimID,
		Provider:          entry.Provider,
		ProviderRequestID: entry.ProviderRequestID,
		IdempotencyKey:    entry.IdempotencyKey,
		Status:            string(entry.Status),
		RequestPayload:    models.JSON(reqBytes),
		ResponsePayload:   models.JSON(respBytes),
		ErrorMessage:      errMsg,
		CreatedAt:         entry.CreatedAt,
	}

	if err := g.db.WithContext(ctx).Create(&row).Error; err != nil {
		log.Printf("bav/audit: failed to persist verification log: %v", err)
		return err
	}
	return nil
}

// persistShape is the storage-side projection of bav.VerifyResult. Unlike
// the v2 wire shape, it INCLUDES RawPayload and ProviderStatusText so dedup
// cache hits can hand back a fully-reconstituted result. Changing this
// struct is a storage migration — backfilling or soft-failing old rows.
type persistShape struct {
	Status             bav.Status   `json:"status"`
	Verified           bool         `json:"verified"`
	Summary            string       `json:"summary"`
	AccountFound       bav.TriState `json:"accountFound"`
	AccountOpen        bav.TriState `json:"accountOpen"`
	IdentityMatch      bav.TriState `json:"identityMatch"`
	AccountTypeMatch   bav.TriState `json:"accountTypeMatch"`
	AcceptsCredits     bav.TriState `json:"acceptsCredits"`
	AcceptsDebits      bav.TriState `json:"acceptsDebits"`
	Provider           string       `json:"provider"`
	ProviderRequestID  string       `json:"providerRequestId"`
	ProviderJobID      string       `json:"providerJobId,omitempty"`
	ProviderStatusText string       `json:"providerStatusText,omitempty"`
	RawPayload         []byte       `json:"rawPayload,omitempty"`
}

// marshalResultForStorage serialises VerifyResult using persistShape so the
// raw provider payload survives the round-trip.
func marshalResultForStorage(r *bav.VerifyResult) ([]byte, error) {
	return json.Marshal(persistShape{
		Status:             r.Status,
		Verified:           r.Verified,
		Summary:            r.Summary,
		AccountFound:       r.AccountFound,
		AccountOpen:        r.AccountOpen,
		IdentityMatch:      r.IdentityMatch,
		AccountTypeMatch:   r.AccountTypeMatch,
		AcceptsCredits:     r.AcceptsCredits,
		AcceptsDebits:      r.AcceptsDebits,
		Provider:           r.Provider,
		ProviderRequestID:  r.ProviderRequestID,
		ProviderJobID:      r.ProviderJobID,
		ProviderStatusText: r.ProviderStatusText,
		RawPayload:         r.RawPayload,
	})
}

// unmarshalStoredResult is the inverse of marshalResultForStorage.
func unmarshalStoredResult(data []byte) (*bav.VerifyResult, error) {
	var p persistShape
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &bav.VerifyResult{
		Status:             p.Status,
		Verified:           p.Verified,
		Summary:            p.Summary,
		AccountFound:       p.AccountFound,
		AccountOpen:        p.AccountOpen,
		IdentityMatch:      p.IdentityMatch,
		AccountTypeMatch:   p.AccountTypeMatch,
		AcceptsCredits:     p.AcceptsCredits,
		AcceptsDebits:      p.AcceptsDebits,
		Provider:           p.Provider,
		ProviderRequestID:  p.ProviderRequestID,
		ProviderJobID:      p.ProviderJobID,
		ProviderStatusText: p.ProviderStatusText,
		RawPayload:         p.RawPayload,
	}, nil
}

var _ bav.Logger = (*GormLogger)(nil)
