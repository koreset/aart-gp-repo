package audit

import (
	"api/services/bav"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestLookupSuccessByKey_NilLoggerReturnsMiss(t *testing.T) {
	var g *GormLogger
	got, err := g.LookupSuccessByKey(context.Background(), "key", time.Hour)
	if err != nil || got != nil {
		t.Fatalf("want (nil, nil) on nil logger; got (%v, %v)", got, err)
	}
}

func TestLookupSuccessByKey_NilDBReturnsMiss(t *testing.T) {
	g := NewGormLogger(nil)
	got, err := g.LookupSuccessByKey(context.Background(), "key", time.Hour)
	if err != nil || got != nil {
		t.Fatalf("want (nil, nil) on nil db; got (%v, %v)", got, err)
	}
}

func TestLookupSuccessByKey_EmptyKeyReturnsMiss(t *testing.T) {
	g := NewGormLogger(nil)
	got, err := g.LookupSuccessByKey(context.Background(), "", time.Hour)
	if err != nil || got != nil {
		t.Fatalf("want (nil, nil) on empty key; got (%v, %v)", got, err)
	}
}

func TestRecord_NilDBIsNoop(t *testing.T) {
	g := NewGormLogger(nil)
	err := g.Record(context.Background(), bav.LogEntry{
		Provider:  "test",
		Status:    bav.StatusComplete,
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("Record should be a no-op with nil db; got err=%v", err)
	}
}

func TestMarshalResultForStorage_RoundTrip(t *testing.T) {
	raw := []byte(`{"verifynow":"raw"}`)
	in := &bav.VerifyResult{
		Status:             bav.StatusComplete,
		Verified:           true,
		Summary:            "ok",
		AccountFound:       bav.TriYes,
		AccountOpen:        bav.TriYes,
		IdentityMatch:      bav.TriUnknown,
		AccountTypeMatch:   bav.TriYes,
		AcceptsCredits:     bav.TriYes,
		AcceptsDebits:      bav.TriNo,
		Provider:           "verifynow",
		ProviderRequestID:  "req-1",
		ProviderStatusText: "Verified",
		RawPayload:         raw,
	}

	bytes, err := marshalResultForStorage(in)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	out, err := unmarshalStoredResult(bytes)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if out.Status != in.Status || out.Verified != in.Verified {
		t.Errorf("status/verified mismatch: got %+v", out)
	}
	if out.IdentityMatch != bav.TriUnknown {
		t.Errorf("IdentityMatch = %q, want unknown", out.IdentityMatch)
	}
	if string(out.RawPayload) != string(raw) {
		t.Errorf("RawPayload round-trip failed: got %q, want %q", out.RawPayload, raw)
	}
	if out.ProviderStatusText != in.ProviderStatusText {
		t.Errorf("ProviderStatusText round-trip: got %q, want %q", out.ProviderStatusText, in.ProviderStatusText)
	}
}

// TestMarshalResultForStorage_ExcludesZeroProviderJobID confirms ProviderJobID
// is omitted when empty so sync-provider logs don't carry a noisy empty
// field. Async providers (Phase 6) will populate it and rely on its presence.
func TestMarshalResultForStorage_ExcludesZeroProviderJobID(t *testing.T) {
	bytes, err := marshalResultForStorage(&bav.VerifyResult{Status: bav.StatusComplete})
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if containsField(bytes, "providerJobId") {
		t.Errorf("expected providerJobId to be omitted when empty; payload=%s", bytes)
	}
}

func containsField(raw []byte, name string) bool {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return false
	}
	_, ok := m[name]
	return ok
}

// Compile-time assertion that errors package is still used — guards against
// silent regressions when we extend this file.
var _ = errors.New
