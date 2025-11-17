package operator

import (
	"net/http"

	"rms/internal/models"
	svc "rms/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateTaskHandler returns a handler that creates a RepairTask.
// Caller roles are enforced at the router level.
func CreateTaskHandler(db *gorm.DB) gin.HandlerFunc {
	taskSvc := svc.NewTaskService(db)
	return func(c *gin.Context) {
		var in models.RepairTask
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		if in.OrderID == 0 || in.MechanicID == 0 || in.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "order_id, mechanic_id and title are required"})
			return
		}
		if err := taskSvc.Create(&in); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		c.JSON(http.StatusCreated, in)
	}
}
