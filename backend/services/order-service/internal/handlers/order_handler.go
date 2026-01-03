package handlers

import (
	"net/http"
	"order-service/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService services.IOrderService
}

func NewOrderHandler(orderService services.IOrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order for the authenticated user
// @Tags orders
// @Accept  json
// @Produce  json
// @Param   order body services.CreateOrderRequest true "Create Order"
// @Success 201 {object} models.Order
// @Failure 400 {object} handlers.GenericErrorResponse
// @Security ApiKeyAuth
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := h.orderService.ValidateUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req services.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Get an order by ID
// @Tags orders
// @Produce  json
// @Param   id path int true "Order ID"
// @Success 200 {object} models.Order
// @Failure 404 {object} handlers.GenericErrorResponse
// @Security ApiKeyAuth
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := h.orderService.ValidateUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrder(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if order.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized access to order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrder godoc
// @Summary Update an order
// @Description Update an order's status
// @Tags orders
// @Accept  json
// @Produce  json
// @Param   id path int true "Order ID"
// @Param   order body services.UpdateOrderRequest true "Update Order"
// @Success 200 {object} models.Order
// @Failure 400 {object} handlers.GenericErrorResponse
// @Security ApiKeyAuth
// @Router /orders/{id} [put]
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := h.orderService.ValidateUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	// Check ownership before updating
	existingOrder, err := h.orderService.GetOrder(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if existingOrder.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized access to order"})
		return
	}

	var req services.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.UpdateOrder(uint(orderID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// DeleteOrder godoc
// @Summary Delete an order
// @Description Delete an order by ID
// @Tags orders
// @Produce  json
// @Param   id path int true "Order ID"
// @Success 200 {object} handlers.GenericSuccessResponse
// @Failure 400 {object} handlers.GenericErrorResponse
// @Security ApiKeyAuth
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := h.orderService.ValidateUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	// Check ownership before deleting
	existingOrder, err := h.orderService.GetOrder(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if existingOrder.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized access to order"})
		return
	}

	if err := h.orderService.DeleteOrder(uint(orderID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted successfully"})
}

// GetOrders godoc
// @Summary Get all orders for the authenticated user
// @Description Get all orders for the authenticated user
// @Tags orders
// @Produce  json
// @Param   offset query int false "Offset"
// @Param   limit query int false "Limit"
// @Success 200 {object} handlers.GetOrdersResponse
// @Failure 500 {object} handlers.GenericErrorResponse
// @Security ApiKeyAuth
// @Router /orders [get]
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := h.orderService.ValidateUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if limit > 100 {
		limit = 100
	}

	orders, err := h.orderService.GetOrdersByUserID(userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}
