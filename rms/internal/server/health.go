package server

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type HealthHandler struct {
    DB *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
    return &HealthHandler{DB: db}
}

func (h *HealthHandler) Get(c *gin.Context) {
    sqlDB, err := h.DB.DB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "db_error"})
        return
    }
    if err := sqlDB.Ping(); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db_down"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}