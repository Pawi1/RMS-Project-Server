package user

import (
	"net/http"

	"rms/internal/middleware/auth"
	"rms/internal/models"
	svc "rms/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type changePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func ChangePasswordHandler(authSvc *svc.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req changePasswordReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		uidv, ok := c.Get(auth.CtxUserIDKey)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		uid := uidv.(string)
		var u models.User
		if err := authSvc.DB.Where("id = ?", uid).First(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.OldPassword)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		nh, err := svc.HashPassword(req.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
			return
		}
		if err := authSvc.DB.Model(&u).Update("password_hash", nh).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update password"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
