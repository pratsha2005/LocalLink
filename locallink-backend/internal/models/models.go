package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
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
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Order struct {
	ID         int         `json:"id"`
	BuyerID    int         `json:"buyerId"`
	ProducerID int         `json:"producerId"`
	TotalPrice float64     `json:"totalPrice"`
	Status     string      `json:"status"`
	CreatedAt  time.Time   `json:"createdAt"`
	Items      []OrderItem `json:"items"`
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"orderId"`
	ProductID int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Review struct {
	ID        int       `json:"id"`
	ProductID int       `json:"productId"`
	UserID    int       `json:"userId"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}

// Input Structs
type RegisterUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateOrderInput struct {
	ProducerID int `json:"producerId"`
	Items      []struct {
		ProductID int `json:"productId"`
		Quantity  int `json:"quantity"`
	} `json:"items"`
}

type CreateReviewInput struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type UpdateUserInput struct {
	Name *string `json:"name"`
}

type UpdateProductInput struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Quantity    *int     `json:"quantity"`
}

type UpdateOrderStatusInput struct {
	Status string `json:"status"`
}