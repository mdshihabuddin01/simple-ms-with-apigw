package handlers

import (
	"auth-service/internal/models"
	"auth-service/internal/services"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock of IAuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req *services.RegisterRequest) (*services.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(req *services.LoginRequest) (*services.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func (m *MockAuthService) ValidateToken(token string) (*models.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		authHandler := NewAuthHandler(mockAuthService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := &services.RegisterRequest{
			Email:     "test@example.com",
			Password:  "password",
			FirstName: "test",
			LastName:  "user",
		}
		reqBytes, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		mockResponse := &services.AuthResponse{
			Token: "some-jwt-token",
			User: &models.User{
				ID:        1,
				Email:     "test@example.com",
				FirstName: "test",
				LastName:  "user",
			},
		}
		mockAuthService.On("Register", reqBody).Return(mockResponse, nil)

		authHandler.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp services.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, mockResponse, &resp)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("bad request", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		authHandler := NewAuthHandler(mockAuthService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		authHandler.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		authHandler := NewAuthHandler(mockAuthService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := &services.LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		}
		reqBytes, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		mockResponse := &services.AuthResponse{
			Token: "some-jwt-token",
			User: &models.User{
				ID:    1,
				Email: "test@example.com",
			},
		}
		mockAuthService.On("Login", reqBody).Return(mockResponse, nil)

		authHandler.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp services.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, mockResponse, &resp)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("bad request", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		authHandler := NewAuthHandler(mockAuthService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		authHandler.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_ValidateToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		authHandler := NewAuthHandler(mockAuthService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/validate", nil)
		c.Request.Header.Set("Authorization", "Bearer some-jwt-token")

		mockUser := &models.User{
			ID:        1,
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockAuthService.On("ValidateToken", "some-jwt-token").Return(mockUser, nil)

		authHandler.ValidateToken(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp gin.H
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp["user"])
		mockAuthService.AssertExpectations(t)
	})

	t.Run("missing token", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		authHandler := NewAuthHandler(mockAuthService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/validate", nil)

		authHandler.ValidateToken(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthHandler_GetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		authHandler := NewAuthHandler(mockAuthService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		mockUser := &models.User{
			ID:        1,
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockAuthService.On("GetUserByID", uint(1)).Return(mockUser, nil)

		authHandler.GetUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.User
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, mockUser.ID, resp.ID)
		assert.Equal(t, mockUser.Email, resp.Email)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockAuthService := new(MockAuthService)
		authHandler := NewAuthHandler(mockAuthService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		mockAuthService.On("GetUserByID", uint(1)).Return(nil, errors.New("user not found"))

		authHandler.GetUser(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}