package public

import (
	"net/http"

	"rms/internal/models"
	svc "rms/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type refreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshHandler returns a handler that verifies refresh token and returns a new access token.
func RefreshHandler(auth *svc.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rr refreshReq
		if err := c.ShouldBindJSON(&rr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		// parse token into RegisteredClaims so we can access Subject and ID
		rc := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(rr.RefreshToken, rc, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, jwt.ErrTokenMalformed
			}
			return auth.JWTSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		// rc now contains the claims
		// load user role from DB so the new access token contains role claim
		var u models.User
		if err := auth.DB.WithContext(c.Request.Context()).Select("id", "role").Where("id = ?", rc.Subject).First(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load user"})
			return
		}
		// issue a new access token including role
		access, err := auth.IssueAccessForSubject(c.Request.Context(), rc.Subject, rc.ID, string(u.Role))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"access_token": access})
	}
}
