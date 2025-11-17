package service

import (
	"rms/internal/models"

	"gorm.io/gorm"
)

type CarService struct {
	DB *gorm.DB
}

func NewCarService(db *gorm.DB) *CarService {
	return &CarService{DB: db}
}

func (s *CarService) List() ([]models.Car, error) {
	var cars []models.Car
	if err := s.DB.Find(&cars).Error; err != nil {
		return nil, err
	}
	return cars, nil
}

func (s *CarService) Create(car *models.Car) error {
	return s.DB.Create(car).Error
}

func (s *CarService) Update(car *models.Car) error {
	return s.DB.Model(&models.Car{}).Where("id = ?", car.ID).Updates(car).Error
}

func (s *CarService) Delete(id uint) error {
	return s.DB.Delete(&models.Car{}, id).Error
}
