package user

import (
	"net/http"

	"rms/internal/middleware/auth"
	"rms/internal/models"
	svc "rms/internal/service"

	"github.com/gin-gonic/gin"
)

type patchUserReq struct {
	UserName *string `json:"user_name"`
	Email    *string `json:"email"`
}

func PatchMeHandler(authSvc *svc.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req patchUserReq
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
		updates := map[string]interface{}{}
		if req.UserName != nil {
			updates["user_name"] = req.UserName
		}
		if req.Email != nil {
			updates["email"] = req.Email
		}
		if len(updates) > 0 {
			if err := authSvc.DB.Model(&u).Updates(updates).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update user"})
				return
			}
		}
		c.JSON(http.StatusOK, u)
	}
}
