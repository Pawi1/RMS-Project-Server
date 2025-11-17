package models

import "time"

type User struct {
    ID           uint      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
    UserName     *string   `json:"user_name,omitempty" gorm:"column:user_name;uniqueIndex"`
    Email        string    `json:"email" gorm:"column:email;uniqueIndex"`
    PasswordHash string    `json:"-" gorm:"column:password_hash"`
    Role         Role      `json:"role" gorm:"column:role;type:text;not null;default:client"`
    CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
    UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (User) TableName() string { return "Users" }