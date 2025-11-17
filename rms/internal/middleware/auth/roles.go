package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRoles checks that user role (set by JWT middleware) is one of allowed roles.
func RequireRoles(allowed ...string) gin.HandlerFunc {
	allowedMap := map[string]struct{}{}
	for _, r := range allowed {
		allowedMap[r] = struct{}{}
	}
	return func(c *gin.Context) {
		v, ok := c.Get(CtxUserRoleKey)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		role, _ := v.(string)
		if _, ok := allowedMap[role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
