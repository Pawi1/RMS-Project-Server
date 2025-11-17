package models

import "time"

type RepairOrder struct {
    ID          uint      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
    CarID       uint      `json:"car_id" gorm:"column:id_car"`
    Title       *string   `json:"title,omitempty" gorm:"column:title"`
    Description *string   `json:"description,omitempty" gorm:"column:description"`
    TotalCost   *float64  `json:"total_cost,omitempty" gorm:"column:total_cost"`
    CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (RepairOrder) TableName() string { return "RepairOrders" }