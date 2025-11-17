package operator

import (
	"net/http"

	"rms/internal/models"
	svc "rms/internal/service"

	"github.com/gin-gonic/gin"
)

type createClientReq struct {
	UserName     *string `json:"user_name"`
	Email        string  `json:"email" binding:"required"`
	Password     *string `json:"password"`
	PasswordHash *string `json:"password_hash"`
}

// CreateClientHandler returns handler that creates a client user.
func CreateClientHandler(authSvc *svc.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createClientReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		// prepare hash
		hash := ""
		if req.PasswordHash != nil && *req.PasswordHash != "" {
			hash = *req.PasswordHash
		} else if req.Password != nil {
			h, err := svc.HashPassword(*req.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
				return
			}
			hash = h
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password or password_hash required"})
			return
		}

		// caller role info is optional here; router already restricts caller
		u := models.User{
			UserName:     req.UserName,
			Email:        req.Email,
			PasswordHash: hash,
			Role:         models.Role("client"),
		}
		if err := authSvc.DB.Create(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		c.JSON(http.StatusCreated, u)
	}
}
