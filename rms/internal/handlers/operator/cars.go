package operator

import (
	"net/http"
	"strconv"

	"rms/internal/models"
	svc "rms/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewCarsHandlers returns handlers bound to DB
func NewCarsHandlers(db *gorm.DB) *carsHandler {
	return &carsHandler{carSvc: svc.NewCarService(db)}
}

type carsHandler struct {
	carSvc *svc.CarService
}

func (h *carsHandler) ListCars(c *gin.Context) {
	cars, err := h.carSvc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, cars)
}

func (h *carsHandler) CreateCar(c *gin.Context) {
	var in models.Car
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	if err := h.carSvc.Create(&in); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusCreated, in)
}

func (h *carsHandler) PatchCar(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var in models.Car
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	in.ID = uint(id)
	if err := h.carSvc.Update(&in); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, in)
}

func (h *carsHandler) DeleteCar(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if err := h.carSvc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.Status(http.StatusNoContent)
}
