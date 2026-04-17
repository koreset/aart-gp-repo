package providers

import (
	"api/services/bav"
	"fmt"
	"time"
)

// Config is the provider-agnostic configuration consumed by NewRegistry.
// Callers populate it from their environment (env vars, config file, etc.)
// and NewRegistry picks the appropriate adapter based on Provider.
type Config struct {
	Provider          string
	APIKey            string
	BaseURL           string
	Mode              string
	OAuthClientID     string
	OAuthClientSecret string
	OAuthTokenURL     string
	Timeout           time.Duration
	// MockAsync enables the pending → complete flow when Provider == "mock".
	// Wired from MOCK_BAV_ASYNC; ignored for non-mock providers.
	MockAsync bool
}

// NewRegistry builds a bav.Registry for the provider named in cfg.Provider.
// An empty Provider defaults to "verifynow". A missing APIKey is tolerated
// — the returned registry's Active() will still be non-nil and the adapter
// will surface ErrProviderNotConfigured at call time.
func NewRegistry(cfg Config) (*bav.Registry, error) {
	name := cfg.Provider
	if name == "" {
		name = "verifynow"
	}

	switch name {
	case "verifynow":
		return bav.NewRegistry(NewVerifyNow(VerifyNowConfig{
			APIKey:  cfg.APIKey,
			Mode:    cfg.Mode,
			BaseURL: cfg.BaseURL,
			Timeout: cfg.Timeout,
		})), nil
	case "mock":
		return bav.NewRegistry(NewMock(MockConfig{Async: cfg.MockAsync})), nil
	default:
		return nil, fmt.Errorf("bav: unknown provider %q", name)
	}
}
