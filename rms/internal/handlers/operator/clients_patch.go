package operator

import (
	"net/http"
	"strconv"

	"rms/internal/models"
	svc "rms/internal/service"

	"github.com/gin-gonic/gin"
)

type patchClientReq struct {
	UserName *string `json:"user_name"`
	Email    *string `json:"email"`
}

// PatchClientHandler allows operator/editor to update client users (restricted by router)
func PatchClientHandler(authSvc *svc.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var req patchClientReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		var u models.User
		if err := authSvc.DB.First(&u, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		// ensure target is client
		if string(u.Role) != "client" {
			c.JSON(http.StatusForbidden, gin.H{"error": "can only modify client users"})
			return
		}
		if req.UserName != nil {
			u.UserName = req.UserName
		}
		if req.Email != nil {
			u.Email = *req.Email
		}
		if err := authSvc.DB.Save(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		c.JSON(http.StatusOK, u)
	}
}
