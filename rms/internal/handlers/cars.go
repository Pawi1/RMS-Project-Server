package handlers

import (
	"net/http"
	"strconv"

	"rms/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewCarsHandlers returns handlers bound to DB
func NewCarsHandlers(db *gorm.DB) *carsHandler {
	return &carsHandler{db: db}
}

type carsHandler struct {
	db *gorm.DB
}

func (h *carsHandler) List(c *gin.Context) {
	var cars []models.Car
	if err := h.db.Find(&cars).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, cars)
}

func (h *carsHandler) Create(c *gin.Context) {
	var in models.Car
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	if err := h.db.Create(&in).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusCreated, in)
}

func (h *carsHandler) Patch(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var in models.Car
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	in.ID = uint(id)
	if err := h.db.Model(&models.Car{}).Where("id = ?", id).Updates(&in).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, in)
}

func (h *carsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if err := h.db.Delete(&models.Car{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.Status(http.StatusNoContent)
}
