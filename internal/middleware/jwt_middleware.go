package middleware

import (
	"net/http"
	"os"
	"pos-be/internal/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "Missing or invalid token")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "secret123"
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("roles", claims["roles"])
		c.Set("permissions", claims["permissions"])
		c.Set("tenant_id", claims["tenant_id"])

		c.Next()
	}
}

// hasRole checks if rolesClaim (from JWT or context) contains target role slug.
func hasRole(rolesClaim interface{}, target string) bool {
	switch v := rolesClaim.(type) {
	case []string:
		for _, s := range v {
			if s == target {
				return true
			}
		}
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok && s == target {
				return true
			}
		}
	case jwt.MapClaims:
		// not expected here, but handle defensively
		if raw, ok := v["roles"]; ok {
			return hasRole(raw, target)
		}
	default:
		// try string (single role)
		if s, ok := v.(string); ok && s == target {
			return true
		}
	}
	return false
}

// hasPermission checks if permissionsClaim (from JWT or context) contains target permission slug.
func hasPermission(permissionsClaim interface{}, target string) bool {
	switch v := permissionsClaim.(type) {
	case []string:
		for _, s := range v {
			if s == target {
				return true
			}
		}
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok && s == target {
				return true
			}
		}
	case jwt.MapClaims:
		if raw, ok := v["permissions"]; ok {
			return hasPermission(raw, target)
		}
	default:
		if s, ok := v.(string); ok && s == target {
			return true
		}
	}
	return false
}

// RequirePermission ensures the JWT contains the required permission slug.
func RequirePermission(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// admin bypass
		rolesClaim, _ := ctx.Get("roles")
		if hasRole(rolesClaim, "admin") {
			ctx.Next()
			return
		}

		permissionsClaim, exists := ctx.Get("permissions")
		if !exists {
			response.Error(ctx, http.StatusForbidden, "Forbidden: missing permissions")
			ctx.Abort()
			return
		}

		if !hasPermission(permissionsClaim, permission) {
			response.Error(ctx, http.StatusForbidden, "Forbidden: insufficient permission")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// TenantRequired ensures the token includes a tenant_id (non-nil), used for tenant-scoped routes.
func TenantRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t, exists := ctx.Get("tenant_id")
		if !exists || t == nil {
			response.Error(ctx, http.StatusForbidden, "Forbidden: tenant required in token")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
