package public

import (
	"net/http"

	appcfg "rms/internal/config"

	"github.com/gin-gonic/gin"
)

// WorkshopHandler returns static workshop info from config.
func WorkshopHandler(w appcfg.Workshop) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        w.Name,
			"description": w.Description,
			"image_path":  w.ImagePath,
			"owner":       w.Owner,
		})
	}
}
