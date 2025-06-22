package model

import "time"

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents an order in the system
type Order struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	UserID      uint        `json:"user_id" gorm:"not null"`
	ProductName string      `json:"product_name" gorm:"not null"`
	Price       float64     `json:"price" gorm:"not null"`
	Quantity    int         `json:"quantity" gorm:"not null"`
	Status      OrderStatus `json:"status" gorm:"not null;default:'pending'"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// CreateOrderRequest represents the request to create a new order
type CreateOrderRequest struct {
	UserID      uint    `json:"user_id" binding:"required" example:"1"`
	ProductName string  `json:"product_name" binding:"required" example:"iPhone 15"`
	Price       float64 `json:"price" binding:"required,gt=0" example:"999.99"`
	Quantity    int     `json:"quantity" binding:"required,gt=0" example:"1"`
}

// UpdateOrderStatusRequest represents the request to update order status
type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required" example:"confirmed"`
}

// OrderResponse represents the response for order operations
type OrderResponse struct {
	ID          uint        `json:"id" example:"1"`
	UserID      uint        `json:"user_id" example:"1"`
	ProductName string      `json:"product_name" example:"iPhone 15"`
	Price       float64     `json:"price" example:"999.99"`
	Quantity    int         `json:"quantity" example:"1"`
	Status      OrderStatus `json:"status" example:"pending"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// OrderWithUser represents an order with user information
type OrderWithUser struct {
	OrderResponse
	User *UserInfo `json:"user,omitempty"`
}

// UserInfo represents basic user information for order context
type UserInfo struct {
	ID        uint   `json:"id" example:"1"`
	Username  string `json:"username" example:"johndoe"`
	Email     string `json:"email" example:"john@example.com"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
} 