package models

import "time"

type RepairTask struct {
    ID          uint      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
    OrderID     uint      `json:"order_id" gorm:"column:order_id"`
    MechanicID  uint      `json:"mechanic_id" gorm:"column:mechanic_id"`
    Title       string    `json:"title" gorm:"column:title;not null"`
    Description *string   `json:"description,omitempty" gorm:"column:description"`
    Hours       *float64  `json:"hours,omitempty" gorm:"column:hours"`
    ImagePath   *string   `json:"image_path,omitempty" gorm:"column:image_path"`
    CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (RepairTask) TableName() string { return "RepairTasks" }