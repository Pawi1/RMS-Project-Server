package models

import "time"

type Car struct {
    ID                uint      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
    OwnerID           *uint     `json:"owner_id,omitempty" gorm:"column:id_owner"`
    VIN               *string   `json:"vin,omitempty" gorm:"column:vin;uniqueIndex"`
    PlateNumber       *string   `json:"plate_number,omitempty" gorm:"column:plate_number"`
    Make              *string   `json:"make,omitempty" gorm:"column:make"`
    Model             *string   `json:"model,omitempty" gorm:"column:model"`
    Year              *int      `json:"year,omitempty" gorm:"column:year"`
    LastMileage       *int      `json:"last_mileage,omitempty" gorm:"column:last_mileage"`
    FuelType          *string   `json:"fuel_type,omitempty" gorm:"column:fuel_type"`
    EngineCapacity    *float64  `json:"engine_capacity,omitempty" gorm:"column:engine_capacity"`
    EngineType        *string   `json:"engine_type,omitempty" gorm:"column:engine_type"`
    DefaultHourlyRate *float64  `json:"default_hourly_rate,omitempty" gorm:"column:default_hourly_rate"`
    Notes             *string   `json:"notes,omitempty" gorm:"column:notes"`
    CreatedAt         time.Time `json:"created_at" gorm:"column:created_at"`
    UpdatedAt         time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Car) TableName() string { return "Cars" }