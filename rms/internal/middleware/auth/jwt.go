package auth

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    svc "rms/internal/service"
)

// keys used in Gin context
const (
    CtxUserIDKey  = "userID"
    CtxUserRoleKey = "userRole"
)

// RequireAuth returns a middleware that verifies Bearer access token and sets user id/role in context.
func RequireAuth(auth *svc.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        authz := c.GetHeader("Authorization")
        if authz == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization"})
            return
        }
        parts := strings.Fields(authz)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
            return
        }
        tokenStr := parts[1]
        // parse into MapClaims so we can read custom fields like "role"
        mc := jwt.MapClaims{}
        token, err := jwt.ParseWithClaims(tokenStr, mc, func(t *jwt.Token) (interface{}, error) {
            if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
                return nil, jwt.ErrTokenMalformed
            }
            return auth.JWTSecret, nil
        })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }
        // extract subject and role if present
        if sub, ok := mc["sub"].(string); ok {
            c.Set(CtxUserIDKey, sub)
        }
        if role, ok := mc["role"].(string); ok {
            c.Set(CtxUserRoleKey, role)
        }
        c.Next()
    }
}