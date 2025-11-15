package main

import (
	// "github.com/golang-jwt/jwt/v4"
	"time"
)

// Users
type User struct {
	ID           uint      `json:"id" db:"id"`
	UserName     string    `json:"user_name,omitempty" db:"user_name"`
	Email        string    `json:"email,omitempty" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // not serialized to JSON
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// Persons
type Person struct {
	UserID    uint      `json:"user_id" db:"user_id"`
	FirstName string    `json:"first_name,omitempty" db:"first_name"`
	LastName  string    `json:"last_name,omitempty" db:"last_name"`
	Phone     string    `json:"phone,omitempty" db:"phone"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// Cars
type Car struct {
	ID                uint      `json:"id" db:"id"`
	OwnerID           *uint     `json:"owner_id,omitempty" db:"id_owner"` // nullable (ON DELETE SET NULL)
	VIN               string    `json:"vin,omitempty" db:"vin"`
	PlateNumber       string    `json:"plate_number,omitempty" db:"plate_number"`
	Make              string    `json:"make,omitempty" db:"make"`
	Model             string    `json:"model,omitempty" db:"model"`
	Year              *int      `json:"year,omitempty" db:"year"`
	LastMileage       *int      `json:"last_mileage,omitempty" db:"last_mileage"`
	FuelType          string    `json:"fuel_type,omitempty" db:"fuel_type"`
	EngineCapacity    *float64  `json:"engine_capacity,omitempty" db:"engine_capacity"`
	EngineType        string    `json:"engine_type,omitempty" db:"engine_type"`
	DefaultHourlyRate *float64  `json:"default_hourly_rate,omitempty" db:"default_hourly_rate"`
	Notes             string    `json:"notes,omitempty" db:"notes"`
	CreatedAt         time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// Mechanic
type Mechanic struct {
	ID         uint      `json:"id" db:"id"` // FK -> Users.id, also PK
	HourlyRate *float64  `json:"hourly_rate,omitempty" db:"hourly_rate"`
	CreatedAt  time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// RepairOrders
type RepairOrder struct {
	ID          uint      `json:"id" db:"id"`
	CarID       uint      `json:"car_id" db:"id_car"`
	Title       string    `json:"title,omitempty" db:"title"`
	Description string    `json:"description,omitempty" db:"description"`
	TotalCost   *float64  `json:"total_cost,omitempty" db:"total_cost"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// RepairTasks
type RepairTask struct {
	ID          uint      `json:"id" db:"id"`
	OrderID     uint      `json:"order_id" db:"order_id"`
	MechanicID  uint      `json:"mechanic_id" db:"mechanic_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description,omitempty" db:"description"`
	Hours       *float64  `json:"hours,omitempty" db:"hours"`
	ImagePath   string    `json:"image_path,omitempty" db:"image_path"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// AccessRequest represents a login request
type AccessRequest struct {
	Email    string      `json:"email"`    // User email
	Password string      `json:"password"` // User password
	Argument interface{} `json:"argument"` // Additional data for the request
}

// // Claims represents JWT claims for authentication
// type Claims struct {
// 	Email string `json:"email"` // User email
// 	Role  string `json:"role"`  // User role
// 	jwt.StandardClaims
// }

// // Input represents a password change request
// type Input struct {
// 	OldPassword string `json:"old_password"` // Current password
// 	NewPassword string `json:"new_password"` // New password
// }
