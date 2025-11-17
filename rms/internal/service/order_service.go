package service

import (
	"rms/internal/models"

	"gorm.io/gorm"
)

type OrderService struct {
	DB *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService { return &OrderService{DB: db} }

func (s *OrderService) Create(o *models.RepairOrder) error {
	return s.DB.Create(o).Error
}
