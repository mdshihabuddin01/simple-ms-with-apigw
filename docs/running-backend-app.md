# Running and Testing the Services

This document provides instructions on how to run the `auth-service` and `order-service` applications and their respective unit tests.

## Auth Service

### Running the Application

To run the `auth-service`, navigate to its directory and use the `go run` command:

```bash
cd backend/services/auth-service
go run cmd/main.go
```

### Running Unit Tests

To run the unit tests for the `auth-service`, navigate to its directory and use the `go test` command:

```bash
cd backend/services/auth-service
go test ./...
```

## Order Service

### Running the Application

To run the `order-service`, navigate to its directory and use the `go run` command:

```bash
cd backend/services/order-service
go run cmd/main.go
```

### Running Unit Tests

To run the unit tests for the `order-service`, navigate to its directory and use the `go test` command:

```bash
cd backend/services/order-service
go test ./...
```
### Running in Docker Compose
```bash
cd backend
docker compose up -d
```