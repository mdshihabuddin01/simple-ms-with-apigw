package services

import "order-service/internal/models"

type IOrderService interface {
	CreateOrder(userID uint, req *CreateOrderRequest) (*models.Order, error)
	GetOrder(orderID uint) (*models.Order, error)
	UpdateOrder(orderID uint, req *UpdateOrderRequest) (*models.Order, error)
	DeleteOrder(orderID uint) error
	GetOrdersByUserID(userID uint, offset, limit int) ([]*models.Order, error)
	ValidateUserID(userIDStr string) (uint, error)
}
