package handlers

import (
	"auth-service/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.IAuthService
}

func NewAuthHandler(authService services.IAuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   register body services.RegisterRequest true "Register"
// @Success 201 {object} services.AuthResponse
// @Failure 400 {object} handlers.GenericErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Login a user
// @Description Login a user with email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   login body services.LoginRequest true "Login"
// @Success 200 {object} services.AuthResponse
// @Failure 401 {object} handlers.GenericErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ValidateToken godoc
// @Summary Validate a token
// @Description Validate a token
// @Tags auth
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} models.User
// @Failure 401 {object} handlers.GenericErrorResponse
// @Router /auth/validate [get]
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
		return
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	user, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if user != nil {
		c.Writer.Header().Set("X-USER-ID", strconv.Itoa(int(user.ID)))
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags auth
// @Produce  json
// @Param   id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} handlers.GenericErrorResponse
// @Router /auth/users/{id} [get]
func (h *AuthHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.authService.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
