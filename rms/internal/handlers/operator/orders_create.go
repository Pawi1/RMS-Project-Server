package operator

import (
	"net/http"

	"rms/internal/models"
	svc "rms/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateOrderHandler returns a handler that creates a RepairOrder.
// Allowed caller roles are enforced at the router level.
func CreateOrderHandler(db *gorm.DB) gin.HandlerFunc {
	orderSvc := svc.NewOrderService(db)
	return func(c *gin.Context) {
		var in models.RepairOrder
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		if in.CarID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "car_id is required"})
			return
		}
		if err := orderSvc.Create(&in); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		c.JSON(http.StatusCreated, in)
	}
}
