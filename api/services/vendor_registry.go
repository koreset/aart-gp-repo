package services

import (
	"sync"

	"api/services/uw_vendor"
	"api/services/uw_vendor/providers"
)

// vendorRegistry is the single Registry shared across the API. Initialised
// lazily on first call so test bootstrap that skips InitVendorRegistry
// still gets a working (empty) instance.
var (
	vendorRegistryOnce sync.Once
	vendorRegistry     *uwvendor.Registry
)

// VendorRegistry returns the shared Registry, constructing it on first
// call. Tests may call InitVendorRegistry to replace the providers with
// custom fixtures before exercising the route handlers.
func VendorRegistry() *uwvendor.Registry {
	vendorRegistryOnce.Do(func() {
		vendorRegistry = uwvendor.NewRegistry(DB)
		registerDefaultVendorProviders(vendorRegistry)
	})
	return vendorRegistry
}

// InitVendorRegistry resets the registry and re-registers a fresh set of
// mock providers. Useful for tests that need a known state per Kind.
func InitVendorRegistry() *uwvendor.Registry {
	vendorRegistry = uwvendor.NewRegistry(DB)
	registerDefaultVendorProviders(vendorRegistry)
	vendorRegistryOnce.Do(func() {}) // ensure the once.Do is consumed
	return vendorRegistry
}

// registerDefaultVendorProviders wires up the bundled mock provider for
// every Kind. Production deployments replace these with real adapters
// (DocuSign / Twilio / LANCET / etc.) before the API listens — the call
// site is intentionally a single function so swapping is trivial.
func registerDefaultVendorProviders(r *uwvendor.Registry) {
	for _, k := range uwvendor.AllKinds {
		r.Register(providers.NewMock(providers.MockConfig{Kind: k, Async: true}))
	}
}
