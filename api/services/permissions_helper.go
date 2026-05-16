package services

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// PermissionResolution is the exportable shape of a per-request
// permission lookup result. Cached in gin context so the middleware and
// downstream controller code share one DB hit. The cache key is a plain
// string so both the routes and controllers packages can read it
// without importing each other.
type PermissionResolution struct {
	HasRole bool
	Slugs   map[string]bool
}

// PermissionResolutionCacheKey is the gin-context key the middleware
// uses when caching a resolved permission set. Exported so service-layer
// helpers (and controllers) can read the same cache and avoid a second
// DB hit on the same request.
const PermissionResolutionCacheKey = "__userPermissionsSvc"

// ResolveUserPermissions returns the current request's permission set.
// Cached in gin context for the remainder of the request. Lookup order
// matches the routes middleware: prefer X-License-Id (frontend supplies
// this), fall back to JWT email. Returns a zero-value (no role, empty
// slugs) on lookup error so callers can default to "deny extras".
func ResolveUserPermissions(c *gin.Context) PermissionResolution {
	if cached, ok := c.Get(PermissionResolutionCacheKey); ok {
		if pr, ok := cached.(PermissionResolution); ok {
			return pr
		}
	}

	pr := PermissionResolution{Slugs: map[string]bool{}}

	var hasRole bool
	var slugs []string
	var err error

	if licenseId := strings.TrimSpace(c.GetHeader("X-License-Id")); licenseId != "" {
		hasRole, slugs, err = GetPermissionsForLicense(licenseId)
	} else {
		email, _ := c.Get("userEmail")
		emailStr, _ := email.(string)
		if emailStr == "" {
			c.Set(PermissionResolutionCacheKey, pr)
			return pr
		}
		hasRole, slugs, err = GetPermissionsForEmail(emailStr)
	}
	if err != nil {
		c.Set(PermissionResolutionCacheKey, pr)
		return pr
	}

	pr.HasRole = hasRole
	for _, s := range slugs {
		pr.Slugs[s] = true
	}
	c.Set(PermissionResolutionCacheKey, pr)
	return pr
}

// UserHasPermission reports whether the active request's user has the
// given permission slug. Mirrors RequirePermission semantics:
//   - user has no role assigned (fresh-install bootstrap) → true
//   - user has system:admin → true
//   - user has the specific slug → true
//   - otherwise false
func UserHasPermission(c *gin.Context, slug string) bool {
	pr := ResolveUserPermissions(c)
	if !pr.HasRole {
		return true
	}
	if pr.Slugs["system:admin"] {
		return true
	}
	return pr.Slugs[slug]
}
