package handlers

import "order-service/internal/models"

// GenericErrorResponse represents a generic error response.
// @name GenericErrorResponse
type GenericErrorResponse struct {
	Error string `json:"error"`
}

// GenericSuccessResponse represents a generic success response.
// @name GenericSuccessResponse
type GenericSuccessResponse struct {
	Message string `json:"message"`
}

// GetOrdersResponse represents the response for getting all orders.
// @name GetOrdersResponse
type GetOrdersResponse struct {
	Orders []*models.Order `json:"orders"`
}
