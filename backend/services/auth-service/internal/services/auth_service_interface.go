package services

import "auth-service/internal/models"

type IAuthService interface {
	Register(req *RegisterRequest) (*AuthResponse, error)
	Login(req *LoginRequest) (*AuthResponse, error)
	ValidateToken(tokenString string) (*models.User, error)
	GetUserByID(userID uint) (*models.User, error)
}
