package services

import (
	"order-service/internal/middleware"
	"order-service/internal/models"
	"order-service/internal/repositories"
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo *repositories.OrderRepository
}

func NewOrderService(orderRepo *repositories.OrderRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
	}
}

type CreateOrderRequest struct {
	OrderItems []CreateOrderItemRequest `json:"order_items" binding:"required,min=1"`
}

type CreateOrderItemRequest struct {
	ProductID   string  `json:"product_id" binding:"required"`
	ProductName string  `json:"product_name" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required,min=1"`
	UnitPrice   float64 `json:"unit_price" binding:"required,min=0"`
}

type UpdateOrderRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed shipped delivered cancelled"`
}

func (s *OrderService) CreateOrder(userID uint, req *CreateOrderRequest) (*models.Order, error) {
	start := time.Now()
	defer func() {
		middleware.RecordDatabaseQuery("order_create", time.Since(start))
	}()

	order := &models.Order{
		UserID:      userID,
		OrderNumber: s.orderRepo.GenerateOrderNumber(),
		Status:      "pending",
		Currency:    "USD",
	}

	itemMap := make(map[string]*models.OrderItem)
	
	for _, item := range req.OrderItems {
		totalPrice := float64(item.Quantity) * item.UnitPrice
		
		if existingItem, exists := itemMap[item.ProductID]; exists {
			existingItem.Quantity += item.Quantity
			existingItem.TotalPrice += totalPrice
		} else {
			orderItem := &models.OrderItem{
				ProductID:   item.ProductID,
				ProductName: item.ProductName,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
				TotalPrice:  totalPrice,
			}
			itemMap[item.ProductID] = orderItem
		}
	}

	var orderItems []models.OrderItem
	var totalAmount float64

	for _, item := range itemMap {
		orderItems = append(orderItems, *item)
		totalAmount += item.TotalPrice
	}

	order.TotalAmount = totalAmount
	order.OrderItems = orderItems

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	middleware.RecordOrderCreated()

	return s.orderRepo.GetByID(order.ID)
}

func (s *OrderService) GetOrder(orderID uint) (*models.Order, error) {
	start := time.Now()
	defer func() {
		middleware.RecordDatabaseQuery("order_get", time.Since(start))
	}()

	return s.orderRepo.GetByID(orderID)
}

func (s *OrderService) GetOrdersByUserID(userID uint, offset, limit int) ([]*models.Order, error) {
	start := time.Now()
	defer func() {
		middleware.RecordDatabaseQuery("orders_list", time.Since(start))
	}()

	return s.orderRepo.GetByUserID(userID, offset, limit)
}

func (s *OrderService) UpdateOrder(orderID uint, req *UpdateOrderRequest) (*models.Order, error) {
	start := time.Now()
	defer func() {
		middleware.RecordDatabaseQuery("order_update", time.Since(start))
	}()

	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	order.Status = req.Status

	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return s.orderRepo.GetByID(orderID)
}

func (s *OrderService) DeleteOrder(orderID uint) error {
	start := time.Now()
	defer func() {
		middleware.RecordDatabaseQuery("order_delete", time.Since(start))
	}()

	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		return err
	}

	if order.Status == "shipped" || order.Status == "delivered" {
		return errors.New("cannot delete order that is shipped or delivered")
	}

	return s.orderRepo.Delete(orderID)
}

func (s *OrderService) ValidateUserID(userIDStr string) (uint, error) {
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, errors.New("invalid user ID")
	}
	return uint(userID), nil
}
