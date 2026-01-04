package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"order-service/internal/models"
	"order-service/internal/services"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(userID uint, req *services.CreateOrderRequest) (*models.Order, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) GetOrder(orderID uint) (*models.Order, error) {
	args := m.Called(orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) UpdateOrder(orderID uint, req *services.UpdateOrderRequest) (*models.Order, error) {
	args := m.Called(orderID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) DeleteOrder(orderID uint) error {
	args := m.Called(orderID)
	return args.Error(0)
}

func (m *MockOrderService) GetOrdersByUserID(userID uint, offset, limit int) ([]*models.Order, error) {
	args := m.Called(userID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Order), args.Error(1)
}

func (m *MockOrderService) ValidateUserID(userIDStr string) (uint, error) {
	args := m.Called(userIDStr)
	return uint(args.Int(0)), args.Error(1)
}

func TestOrderHandler_CreateOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockOrderService := new(MockOrderService)
		orderHandler := NewOrderHandler(mockOrderService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "1")

		reqBody := &services.CreateOrderRequest{
			OrderItems: []services.CreateOrderItemRequest{
				{
					ProductID:   "test-product-id",
					ProductName: "Test Product",
					Quantity:    1,
					UnitPrice:   10.0,
				},
			},
		}
		reqBytes, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		mockOrder := &models.Order{
			ID:          1,
			UserID:      1,
			OrderNumber: "ORD-123",
			Status:      "pending",
			TotalAmount: 10.0,
			Currency:    "USD",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			OrderItems: []models.OrderItem{
				{
					ID:          1,
					OrderID:     1,
					ProductID:   "test-product-id",
					ProductName: "Test Product",
					Quantity:    1,
					UnitPrice:   10.0,
					TotalPrice:  10.0,
				},
			},
		}
		mockOrderService.On("ValidateUserID", "1").Return(1, nil)
		mockOrderService.On("CreateOrder", uint(1), reqBody).Return(mockOrder, nil)

		orderHandler.CreateOrder(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp models.Order
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, mockOrder.ID, resp.ID)
		assert.Equal(t, mockOrder.UserID, resp.UserID)
		assert.Equal(t, mockOrder.OrderNumber, resp.OrderNumber)
		assert.Equal(t, mockOrder.Status, resp.Status)
		assert.Equal(t, mockOrder.TotalAmount, resp.TotalAmount)
		assert.Equal(t, mockOrder.Currency, resp.Currency)
		assert.Equal(t, len(mockOrder.OrderItems), len(resp.OrderItems))
		mockOrderService.AssertExpectations(t)
	})

	t.Run("invalid user id", func(t *testing.T) {
		mockOrderService := new(MockOrderService)
		orderHandler := NewOrderHandler(mockOrderService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "invalid")

		mockOrderService.On("ValidateUserID", "invalid").Return(0, errors.New("invalid user ID"))

		orderHandler.CreateOrder(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestOrderHandler_GetOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockOrderService := new(MockOrderService)
		orderHandler := NewOrderHandler(mockOrderService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		mockOrder := &models.Order{
			ID:          1,
			UserID:      1,
			OrderNumber: "ORD-123",
			Status:      "pending",
			TotalAmount: 10.0,
			Currency:    "USD",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		mockOrderService.On("GetOrder", uint(1)).Return(mockOrder, nil)

		orderHandler.GetOrder(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.Order
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, mockOrder.ID, resp.ID)
		assert.Equal(t, mockOrder.UserID, resp.UserID)
		assert.Equal(t, mockOrder.OrderNumber, resp.OrderNumber)
		assert.Equal(t, mockOrder.Status, resp.Status)
		assert.Equal(t, mockOrder.TotalAmount, resp.TotalAmount)
		assert.Equal(t, mockOrder.Currency, resp.Currency)
		mockOrderService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockOrderService := new(MockOrderService)
		orderHandler := NewOrderHandler(mockOrderService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		mockOrderService.On("GetOrder", uint(1)).Return(nil, errors.New("order not found"))

		orderHandler.GetOrder(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestOrderHandler_UpdateOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockOrderService := new(MockOrderService)
		orderHandler := NewOrderHandler(mockOrderService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		reqBody := &services.UpdateOrderRequest{
			Status: "shipped",
		}
		reqBytes, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPut, "/orders/1", bytes.NewBuffer(reqBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		mockOrder := &models.Order{
			ID:          1,
			UserID:      1,
			OrderNumber: "ORD-123",
			Status:      "shipped",
			TotalAmount: 10.0,
			Currency:    "USD",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		mockOrderService.On("UpdateOrder", uint(1), reqBody).Return(mockOrder, nil)

		orderHandler.UpdateOrder(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.Order
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, mockOrder.Status, resp.Status)
		mockOrderService.AssertExpectations(t)
	})
}

func TestOrderHandler_DeleteOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockOrderService := new(MockOrderService)
		orderHandler := NewOrderHandler(mockOrderService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		mockOrderService.On("DeleteOrder", uint(1)).Return(nil)

		orderHandler.DeleteOrder(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestOrderHandler_GetOrders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockOrderService := new(MockOrderService)
		orderHandler := NewOrderHandler(mockOrderService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "1")
		c.Request, _ = http.NewRequest(http.MethodGet, "/orders?offset=0&limit=10", nil)

		mockOrders := []*models.Order{
			{
				ID:          1,
				UserID:      1,
				OrderNumber: "ORD-123",
				Status:      "pending",
				TotalAmount: 10.0,
				Currency:    "USD",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}
		mockOrderService.On("ValidateUserID", "1").Return(1, nil)
		mockOrderService.On("GetOrdersByUserID", uint(1), 0, 10).Return(mockOrders, nil)

		orderHandler.GetOrders(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp gin.H
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp["orders"])
		mockOrderService.AssertExpectations(t)
	})
}