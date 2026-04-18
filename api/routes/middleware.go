package routes

import (
	"api/log"
	"api/models"
	"api/services"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// RequestLoggerMiddleware logs all incoming requests with their method, path, status code, and response time
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Create a request ID and add it to the context
		requestID := log.NewRequestID()
		c.Set("requestID", requestID)

		// Set request ID in response headers
		c.Writer.Header().Set("X-Request-ID", requestID)

		// Process request
		c.Next()

		// Calculate response time
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// Get status code
		statusCode := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get user info if available
		userEmail, _ := c.Get("userEmail")
		userName, _ := c.Get("userName")

		// Create context with request ID and user info
		ctx := context.Background()
		ctx = context.WithValue(ctx, log.RequestIDKey, requestID)
		if userEmail != nil {
			ctx = context.WithValue(ctx, log.UserEmailKey, userEmail.(string))
		}
		if userName != nil {
			ctx = context.WithValue(ctx, log.UserNameKey, userName.(string))
		}

		// Log request details with appropriate log level based on status code
		fields := log.WithContext(ctx).WithFields(map[string]interface{}{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status_code": statusCode,
			"latency_ms":  latency.Milliseconds(),
			"client_ip":   clientIP,
		})

		if c.Request.URL.Path != "/health" && c.Request.URL.Path != "/valuations/jobs" {
			if statusCode >= 500 {
				fields.Error("Server error occurred while processing request")
			} else if statusCode >= 400 {
				fields.Warn("Client error occurred while processing request")
			} else {
				fields.Info("Request processed successfully")
			}

		}
	}
}

//func AuthMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		var userToken models.UserToken
//		tokenString := strings.Trim(c.GetHeader("Authorization"), "Bearer")
//		tokenString = strings.TrimSpace(tokenString)
//
//		err := services.DB.Where("token_string = ? ", tokenString).First(&userToken).Error
//		if err == nil && userToken.TokenString != "" {
//			storedToken, err := jwt1.ParseString(tokenString)
//			fmt.Println(storedToken.Expiration().Local().Unix())
//			fmt.Println(time.Now().Local().Unix())
//			if storedToken != nil && storedToken.Expiration().Local().Unix() < time.Now().Local().Unix() {
//				log.Error().Err(err).Msg("JWT validation error")
//				c.Abort()
//				c.Writer.WriteHeader(http.StatusUnauthorized)
//				c.Writer.Write([]byte("Access token has expired"))
//			} else {
//				c.Set("groups", storedToken.Claims(context.Background()))
//				c.Set("user", storedToken.Subject())
//				if err != nil {
//					fmt.Println(err)
//				}
//			}
//			return
//		}
//
//		toValidate := map[string]string{}
//		toValidate["aud"] = "api://aart-app"
//		toValidate["cid"] = "0oa2e2l7tp5yK97tb357"
//
//		jwtVerifierSetup := jwtverifier.JwtVerifier{
//			Issuer:           "https://dev-376454.okta.com/oauth2/aus3j00ke8tb8O6dT357",
//			ClaimsToValidate: toValidate,
//		}
//
//		verifier := jwtVerifierSetup.New()
//
//		token, err := verifier.VerifyAccessToken(tokenString)
//		if err != nil {
//			log.Error().Err(err).Msg("JWT validation error")
//			c.Abort()
//			c.Writer.WriteHeader(http.StatusUnauthorized)
//			c.Writer.Write([]byte("Unauthorized"))
//			return
//		} else {
//			fmt.Println(token.Claims)
//		}
//
//		c.Set("groups", token.Claims["scp"])
//		c.Set("claims", token.Claims["sub"])
//		cache.Set(token.Claims["sub"], tokenString, 1)
//		err = services.StoreUserToken(tokenString, token.Claims["sub"].(string))
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//}

func GetActiveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Authorization") != "" {
			tokenString := c.GetHeader("Authorization")

			userJwt := strings.Split(tokenString, "Bearer ")[1]
			claims := jwt.MapClaims{}
			jwt.ParseWithClaims(userJwt, claims, func(token *jwt.Token) (interface{}, error) {
				return nil, nil // or return nil, nil if absolutely necessary
			})

			// Get request ID from context if available
			requestID, exists := c.Get("requestID")
			var ctx context.Context
			if exists {
				ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
			} else {
				ctx = context.Background()
			}

			// for now we ignore the error
			//if err != nil {
			//	log.WithContext(ctx).WithField("error", err.Error()).Error("Authentication failed: error parsing JWT token")
			//	c.Abort()
			//	c.Writer.WriteHeader(http.StatusUnauthorized)
			//	c.Writer.Write([]byte("error parsing JWT token"))
			//	return
			//}

			if claims == nil {
				log.WithContext(ctx).Error("Authentication failed: valid user token could not be found")
				c.Abort()
				c.Writer.WriteHeader(http.StatusUnauthorized)
				c.Writer.Write([]byte("a valid user token could not be found"))
				return
			}

			activeUser, ok := claims["user"].(map[string]interface{})
			if !ok {
				log.WithContext(ctx).Error("Authentication failed: user claim not found in token")
				c.Abort()
				c.Writer.WriteHeader(http.StatusUnauthorized)
				c.Writer.Write([]byte("user claim not found in token"))
				return
			}

			email, emailOk := activeUser["Email"].(string)
			fullName, nameOk := activeUser["FullName"].(string)

			if !emailOk || !nameOk {
				log.WithContext(ctx).Error("Authentication failed: email or full name not found in user claim")
				c.Abort()
				c.Writer.WriteHeader(http.StatusUnauthorized)
				c.Writer.Write([]byte("email or full name not found in user claim"))
				return
			}

			c.Set("userEmail", email)
			c.Set("userName", fullName)
			user := models.AppUser{
				UserEmail: email,
				UserName:  fullName,
			}
			c.Set("user", user)

			// Update context with user info and log successful authentication
			//ctx = log.ContextWithUserInfo(ctx, email, fullName)
			//log.WithContext(ctx).Info("User authenticated successfully")

		} else {
			// Get request ID from context if available
			requestID, exists := c.Get("requestID")
			var ctx context.Context
			if exists {
				ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
			} else {
				ctx = context.Background()
			}

			log.WithContext(ctx).Error("Authentication failed: no authorization header provided")
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("this action is unauthorized. No valid or licensed user in the header"))
			return
		}
	}
}

// userPermissionsCacheKey holds the result of looking up the caller's
// permissions on the Gin context, so a request passing through multiple
// RequirePermission middlewares hits the DB once.
const userPermissionsCacheKey = "__userPermissions"

type userPermissions struct {
	hasRole bool
	slugs   map[string]bool
}

// resolveUserPermissions looks up (and caches on c) the active user's
// permission set. Never panics — on any lookup error it logs and returns a
// zero value so the caller can decide how to fail.
//
// Lookup priority: X-License-Id header (matches the frontend's
// loadUserPermissions, which resolves by license_id). Falls back to the JWT
// email if the header is missing — supports older clients and avoids
// breaking non-Electron callers.
func resolveUserPermissions(c *gin.Context) userPermissions {
	if cached, ok := c.Get(userPermissionsCacheKey); ok {
		if up, ok := cached.(userPermissions); ok {
			return up
		}
	}

	up := userPermissions{slugs: map[string]bool{}}

	var hasRole bool
	var slugs []string
	var err error

	if licenseId := strings.TrimSpace(c.GetHeader("X-License-Id")); licenseId != "" {
		hasRole, slugs, err = services.GetPermissionsForLicense(licenseId)
	} else {
		email, _ := c.Get("userEmail")
		emailStr, _ := email.(string)
		if emailStr == "" {
			c.Set(userPermissionsCacheKey, up)
			return up
		}
		hasRole, slugs, err = services.GetPermissionsForEmail(emailStr)
	}

	if err != nil {
		log.WithField("error", err.Error()).Error("RequirePermission: failed to resolve user permissions")
		c.Set(userPermissionsCacheKey, up)
		return up
	}
	up.hasRole = hasRole
	for _, s := range slugs {
		up.slugs[s] = true
	}
	c.Set(userPermissionsCacheKey, up)
	return up
}

// RequirePermission returns a middleware that enforces the given permission
// slug. Matches the frontend's usePermissionCheck semantics:
//   - user has no role assigned → allow (fresh-install bootstrap)
//   - user has system:admin → allow
//   - user has the specific slug → allow
//   - otherwise 403.
func RequirePermission(slug string) gin.HandlerFunc {
	return func(c *gin.Context) {
		up := resolveUserPermissions(c)

		if !up.hasRole {
			c.Next()
			return
		}
		if up.slugs["system:admin"] || up.slugs[slug] {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "You do not have the required permission: " + slug,
		})
	}
}

// RequireEntitlement returns a middleware that checks the X-Entitlements header
// for a specific entitlement string. Users with "all-features" always pass.
func RequireEntitlement(entitlement string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("X-Entitlements")
		for _, e := range strings.Split(header, ",") {
			e = strings.TrimSpace(e)
			if e == entitlement || e == "all-features" {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "You do not have the required entitlement: " + entitlement,
		})
	}
}

func ActivityLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user info from context
		userEmail, _ := c.Get("userEmail")
		userName, _ := c.Get("userName")

		// Create context with user info
		ctx := context.Background()
		if userEmail != nil {
			ctx = context.WithValue(ctx, log.UserEmailKey, userEmail.(string))
		}
		if userName != nil {
			ctx = context.WithValue(ctx, log.UserNameKey, userName.(string))
		}

		log.WithContext(ctx).Info("User activity logged")
	}
}
