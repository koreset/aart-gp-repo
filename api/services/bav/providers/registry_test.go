package providers

import (
	"api/services/bav"
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewRegistry_DefaultsToVerifyNow(t *testing.T) {
	reg, err := NewRegistry(Config{APIKey: "k"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reg.Active() == nil {
		t.Fatal("Active() should be non-nil")
	}
	if got := reg.Active().Name(); got != "verifynow" {
		t.Errorf("provider name = %q, want verifynow", got)
	}
}

func TestNewRegistry_ExplicitVerifyNow(t *testing.T) {
	reg, err := NewRegistry(Config{
		Provider: "verifynow",
		APIKey:   "k",
		BaseURL:  "https://example.test",
		Mode:     "test",
		Timeout:  10 * time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reg.Active().Name() != "verifynow" {
		t.Fatal("expected verifynow provider")
	}
}

func TestNewRegistry_UnknownProvider(t *testing.T) {
	_, err := NewRegistry(Config{Provider: "does-not-exist"})
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestNewRegistry_MissingAPIKeySurfacesAtCallTime(t *testing.T) {
	reg, err := NewRegistry(Config{Provider: "verifynow"})
	if err != nil {
		t.Fatalf("unexpected error at build time: %v", err)
	}
	// Build succeeds — startup doesn't fail when API key is absent. The
	// adapter surfaces ErrProviderNotConfigured when Verify is actually
	// called. Phase 8 will add a startup warning; Phase 3 only requires
	// that nothing panics.
	bav.SetDefault(reg)
	t.Cleanup(func() { bav.SetDefault(nil) })
	_, err = bav.Verify(context.Background(), bav.VerifyRequest{})
	if !errors.Is(err, bav.ErrProviderNotConfigured) {
		t.Fatalf("err = %v, want ErrProviderNotConfigured", err)
	}
}
