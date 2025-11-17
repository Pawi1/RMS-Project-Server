package models

import "time"

type Person struct {
    UserID    uint      `json:"user_id" gorm:"column:user_id;primaryKey"`
    FirstName *string   `json:"first_name,omitempty" gorm:"column:first_name"`
    LastName  *string   `json:"last_name,omitempty" gorm:"column:last_name"`
    Phone     *string   `json:"phone,omitempty" gorm:"column:phone"`
    CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
    UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Person) TableName() string { return "Persons" }