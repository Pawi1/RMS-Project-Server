package models

import "time"

type Mechanic struct {
    ID         uint      `json:"id" gorm:"column:id;primaryKey"` // PK = FK -> Users.id
    HourlyRate *float64  `json:"hourly_rate,omitempty" gorm:"column:hourly_rate"`
    CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
    UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Mechanic) TableName() string { return "Mechanic" }