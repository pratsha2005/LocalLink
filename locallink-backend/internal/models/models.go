// internal/models/models.go
package models

import "time"


type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Don't send password hash to client
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Product struct {
	ID          int       `json:"id"`
	ProducerID  int       `json:"producerId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	Latitude    float64   `json:"latitude"`  // For JSON
	Longitude   float64   `json:"longitude"` // For JSON
	CreatedAt   time.Time `json:"createdAt"`
}

// For creating a user
type RegisterUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // 'producer' or 'buyer'
}

// For logging in
type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}