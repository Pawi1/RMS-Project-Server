package public

import (
	"net/http"

	"rms/internal/models"
	svc "rms/internal/service"

	"github.com/gin-gonic/gin"
)

// LoginHandler returns a gin handler that performs login using AuthService.
func LoginHandler(auth *svc.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.AccessRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		access, refresh, err := auth.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh})
	}
}
