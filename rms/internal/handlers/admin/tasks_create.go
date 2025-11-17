package admin

import (
	operatorh "rms/internal/handlers/operator"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// CreateUserHandler for admin simply reuses operator logic (admins can create users too).
func CreateTaskHandler(db *gorm.DB) gin.HandlerFunc {
	return operatorh.CreateTaskHandler(db)
}
