package service

import (
	"rms/internal/models"

	"gorm.io/gorm"
)

type TaskService struct {
	DB *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService { return &TaskService{DB: db} }

func (s *TaskService) Create(t *models.RepairTask) error {
	return s.DB.Create(t).Error
}
