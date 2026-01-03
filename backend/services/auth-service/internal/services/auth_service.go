package services

import (
	"auth-service/internal/auth"
	"auth-service/internal/middleware"
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"errors"
	"time"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repositories.UserRepository
	jwtSvc   *auth.JWTService
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSvc *auth.JWTService) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtSvc:   jwtSvc,
	}
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *models.User `json:"user"`
}

func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	start := time.Now()
	defer func() {
		middleware.RecordDatabaseQuery("user_register", time.Since(start))
	}()

	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := s.jwtSvc.GenerateToken(user.ID, user.Email, user.IsActive)
	if err != nil {
		return nil, err
	}

	middleware.RecordAuthTokenIssued()

	user.Password = ""
	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	start := time.Now()
	defer func() {
		middleware.RecordDatabaseQuery("user_login", time.Since(start))
	}()

	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	if err := auth.CheckPassword(user.Password, req.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwtSvc.GenerateToken(user.ID, user.Email, user.IsActive)
	if err != nil {
		return nil, err
	}

	middleware.RecordAuthTokenIssued()

	user.Password = ""
	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	claims, err := s.jwtSvc.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if !claims.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	user.Password = ""
	return user, nil
}

func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	start := time.Now()
	defer func() {
		middleware.RecordDatabaseQuery("user_get", time.Since(start))
	}()

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}
