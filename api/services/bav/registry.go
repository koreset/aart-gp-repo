package bav

import "context"

// defaultRegistry holds the process-wide active provider registry. Set once
// at startup via SetDefault; read by Active() and the package-level Verify
// convenience. Not protected by a mutex — writes must happen before the HTTP
// server accepts requests.
var defaultRegistry *Registry

// SetDefault installs the supplied registry as the process-wide default.
// Intended to be called exactly once from application bootstrap.
func SetDefault(r *Registry) {
	defaultRegistry = r
}

// Default returns the process-wide registry, or nil if none has been set.
func Default() *Registry {
	return defaultRegistry
}

// Active returns the active provider from the process-wide default registry,
// or nil if no registry has been installed or no provider is configured.
func Active() Provider {
	return defaultRegistry.Active()
}

// Verify is a convenience wrapper around the process-wide Registry's Verify.
// Returns ErrProviderNotConfigured when no registry is installed.
func Verify(ctx context.Context, req VerifyRequest) (*VerifyResult, error) {
	if defaultRegistry == nil {
		return nil, ErrProviderNotConfigured
	}
	return defaultRegistry.Verify(ctx, req)
}

// Poll resolves a previously-issued async verification job via the
// process-wide Registry.
func Poll(ctx context.Context, jobID string) (*VerifyResult, error) {
	if defaultRegistry == nil {
		return nil, ErrProviderNotConfigured
	}
	return defaultRegistry.Poll(ctx, jobID)
}
