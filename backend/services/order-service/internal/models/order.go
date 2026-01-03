package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	OrderNumber string         `json:"order_number" gorm:"uniqueIndex;size:255;not null"`
	Status      string         `json:"status" gorm:"default:'pending'"`
	TotalAmount float64        `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	Currency    string         `json:"currency" gorm:"default:'USD'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	OrderItems  []OrderItem    `json:"order_items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	OrderID     uint           `json:"order_id" gorm:"not null;index"`
	ProductID   string         `json:"product_id" gorm:"size:255;not null"`
	ProductName string         `json:"product_name" gorm:"not null"`
	Quantity    int            `json:"quantity" gorm:"not null"`
	UnitPrice   float64        `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	TotalPrice  float64        `json:"total_price" gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
