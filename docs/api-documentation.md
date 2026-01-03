# API Documentation

## Overview

This document provides comprehensive API documentation for the Auth Service and Order Service microservices. Both services are RESTful APIs built with Go and Gin framework, using MySQL for data persistence.

## Base URL

- **API Gateway**: `http://localhost:8080`
- **Auth Service**: `http://localhost:8081` (direct)
- **Order Service**: `http://localhost:8082` (direct)

## Authentication

All order endpoints require JWT authentication. Include the following headers:

```
Authorization: Bearer <jwt_token>
```
**Headers for direct access order service**
```
X-User-ID: <user_id>
```



## Auth Service Endpoints

### 1. Register User

**Endpoint**: `POST /api/v1/auth/register`

**Description**: Creates a new user account

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response** (201 Created):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true,
    "created_at": "2025-12-31T09:22:27.662Z",
    "updated_at": "2025-12-31T09:22:27.662Z"
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid input data or user already exists
- `500 Internal Server Error`: Server error

### 2. Login User

**Endpoint**: `POST /api/v1/auth/login`

**Description**: Authenticates a user and returns JWT token

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response** (200 OK):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true,
    "created_at": "2025-12-31T09:22:27.662Z",
    "updated_at": "2025-12-31T09:22:27.662Z"
  }
}
```

**Error Responses**:
- `401 Unauthorized`: Invalid credentials
- `400 Bad Request`: Invalid input data
- `500 Internal Server Error`: Server error

### 3. Validate Token

**Endpoint**: `GET /api/v1/auth/validate`

**Description**: Validates a JWT token and returns user information

**Headers**:
```
Authorization: Bearer <jwt_token>
```

**Response** (200 OK):
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true,
    "created_at": "2025-12-31T09:22:27.662Z",
    "updated_at": "2025-12-31T09:22:27.662Z"
  }
}
```
**Error Responses**:
- `401 Unauthorized`: Invalid or expired token
- `400 Bad Request`: Missing authorization header



## Order Service Endpoints

### 1. Create Order

**Endpoint**: `POST /api/v1/orders`

**Description**: Creates a new order with items

**Headers**:
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body**:
```json
{
  "order_items": [
    {
      "product_id": "1",
      "product_name": "rice",
      "quantity": 2,
      "unit_price": 10.50
    },
    {
      "product_id": "2",
      "product_name": "bread",
      "quantity": 1,
      "unit_price": 3.25
    }
  ]
}
```

**Response** (201 Created):
```json
{
  "id": 1,
  "user_id": 4,
  "order_number": "ORD-20251231092622",
  "status": "pending",
  "total_amount": 24.25,
  "currency": "USD",
  "created_at": "2025-12-31T09:26:22.948Z",
  "updated_at": "2025-12-31T09:26:22.948Z",
  "order_items": [
    {
      "id": 1,
      "order_id": 1,
      "product_id": "1",
      "product_name": "rice",
      "quantity": 2,
      "unit_price": 10.50,
      "total_price": 21.00,
      "created_at": "2025-12-31T09:26:22.956Z",
      "updated_at": "2025-12-31T09:26:22.956Z"
    },
    {
      "id": 2,
      "order_id": 1,
      "product_id": "2",
      "product_name": "bread",
      "quantity": 1,
      "unit_price": 3.25,
      "total_price": 3.25,
      "created_at": "2025-12-31T09:26:22.956Z",
      "updated_at": "2025-12-31T09:26:22.956Z"
    }
  ]
}
```

**Error Responses**:
- `400 Bad Request`: Invalid input data, authentication issues
- `401 Unauthorized`: Invalid or missing authentication
- `500 Internal Server Error`: Server error

### 2. Get All Orders

**Endpoint**: `GET /api/v1/orders`

**Description**: Retrieves all orders for the authenticated user

**Headers**:
```
Authorization: Bearer <jwt_token>
```

**Query Parameters**:
- `offset` (integer, optional): Number of orders to skip (default: 0)
- `limit` (integer, optional): Maximum number of orders to return (default: 10, max: 100)

**Response** (200 OK):
```json
{
  "orders": [
    {
      "id": 1,
      "user_id": 4,
      "order_number": "ORD-20251231092622",
      "status": "pending",
      "total_amount": 24.25,
      "currency": "USD",
      "created_at": "2025-12-31T09:26:22.948Z",
      "updated_at": "2025-12-31T09:26:22.948Z",
      "order_items": [
        {
          "id": 1,
          "order_id": 1,
          "product_id": "1",
          "product_name": "rice",
          "quantity": 2,
          "unit_price": 10.50,
          "total_price": 21.00,
          "created_at": "2025-12-31T09:26:22.956Z",
          "updated_at": "2025-12-31T09:26:22.956Z"
        }
      ]
    }
  ]
}
```

**Error Responses**:
- `401 Unauthorized`: Invalid or missing authentication
- `500 Internal Server Error`: Server error

### 3. Get Order by ID

**Endpoint**: `GET /api/v1/orders/{id}`

**Description**: Retrieves a specific order by ID

**Headers**:
```
Authorization: Bearer <jwt_token>
```

**Path Parameters**:
- `id` (integer): Order ID

**Response** (200 OK):
```json
{
  "id": 1,
  "user_id": 4,
  "order_number": "ORD-20251231092622",
  "status": "pending",
  "total_amount": 24.25,
  "currency": "USD",
  "created_at": "2025-12-31T09:26:22.948Z",
  "updated_at": "2025-12-31T09:26:22.948Z",
  "order_items": [
    {
      "id": 1,
      "order_id": 1,
      "product_id": "1",
      "product_name": "rice",
      "quantity": 2,
      "unit_price": 10.50,
      "total_price": 21.00,
      "created_at": "2025-12-31T09:26:22.956Z",
      "updated_at": "2025-12-31T09:26:22.956Z"
    }
  ]
}
```

**Error Responses**:
- `401 Unauthorized`: Invalid or missing authentication
- `404 Not Found`: Order not found
- `400 Bad Request`: Invalid order ID
- `500 Internal Server Error`: Server error

### 4. Update Order

**Endpoint**: `PUT /api/v1/orders/{id}`

**Description**: Updates an order's status

**Headers**:
```
Authorization: Bearer <jwt_token>

```

**Path Parameters**:
- `id` (integer): Order ID

**Request Body**:
```json
{
  "status": "confirmed"
}
```

**Valid Status Values**:
- `pending`
- `confirmed`
- `shipped`
- `delivered`
- `cancelled`

**Response** (200 OK):
```json
{
  "id": 1,
  "user_id": 4,
  "order_number": "ORD-20251231092622",
  "status": "confirmed",
  "total_amount": 24.25,
  "currency": "USD",
  "created_at": "2025-12-31T09:26:22.948Z",
  "updated_at": "2025-12-31T09:30:15.123Z",
  "order_items": [
    {
      "id": 1,
      "order_id": 1,
      "product_id": "1",
      "product_name": "rice",
      "quantity": 2,
      "unit_price": 10.50,
      "total_price": 21.00,
      "created_at": "2025-12-31T09:26:22.956Z",
      "updated_at": "2025-12-31T09:26:22.956Z"
    }
  ]
}
```

**Error Responses**:
- `401 Unauthorized`: Invalid or missing authentication
- `404 Not Found`: Order not found
- `400 Bad Request`: Invalid order ID or status value
- `500 Internal Server Error`: Server error

### 5. Delete Order

**Endpoint**: `DELETE /api/v1/orders/{id}`

**Description**: Deletes an order (soft delete)

**Headers**:
```
Authorization: Bearer <jwt_token>
```

**Path Parameters**:
- `id` (integer): Order ID

**Response** (200 OK):
```json
{
  "message": "order deleted successfully"
}
```

**Error Responses**:
- `401 Unauthorized`: Invalid or missing authentication
- `404 Not Found`: Order not found
- `400 Bad Request`: Cannot delete order that is shipped or delivered
- `500 Internal Server Error`: Server error

## Health Check

### Gateway Health Check

**Endpoint**: `GET /health`

**Description**: Checks if the API gateway is healthy

**Response** (200 OK):
```json
{
  "status": "healthy",
  "service": "api-gateway"
}
```

### Service Health Checks

**Auth Service**: `GET http://localhost:8081/health`
**Order Service**: `GET http://localhost:8082/health`

**Response** (200 OK):
```json
{
  "status": "healthy",
  "service": "auth-service"  // or "order-service"
}
```

## Error Response Format

All error responses follow this format:

```json
{
  "error": "Error message describing what went wrong"
}
```

