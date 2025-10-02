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
		c.Set("role_id", claims["role_id"])

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRoleID, exists := ctx.Get("role_id")
		if !exists {
			response.Error(ctx, http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		switch roleID := userRoleID.(type) {
		case float64: // dari JWT MapClaims
			if int(roleID) != 1 {
				response.Error(ctx, http.StatusForbidden, "Forbidden: Admin only")
				ctx.Abort()
				return
			}
		case uint: // kalau dari DB atau langsung set
			if roleID != 1 {
				response.Error(ctx, http.StatusForbidden, "Forbidden: Admin only")
				ctx.Abort()
				return
			}
		default:
			response.Error(ctx, http.StatusForbidden, "Forbidden: Invalid role type")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func AdminOrStaffOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRoleID, exists := ctx.Get("role_id")
		if !exists {
			response.Error(ctx, http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		// Cek apakah role_id bukan 1 (admin) dan bukan 2 (staff/kasir)
		switch roleID := userRoleID.(type) {
		case float64: // jwt.MapClaims biasanya decode ke float64
			if roleID != 1 && roleID != 2 {
				response.Error(ctx, http.StatusForbidden, "Forbidden: Admin or Staff only")
				ctx.Abort()
				return
			}
		case uint:
			if roleID != 1 && roleID != 2 {
				response.Error(ctx, http.StatusForbidden, "Forbidden: Admin or Staff only")
				ctx.Abort()
				return
			}
		default:
			response.Error(ctx, http.StatusForbidden, "Forbidden: Invalid role")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
