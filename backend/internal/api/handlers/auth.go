package handlers

import (
	auth "backend/internal/services/auth_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler структурирует обработку запросов для аутентификации.
type AuthHandler struct {
	authService *auth.AuthService // Заменили `AuthService` на `*AuthService`
}

// NewAuthHandler создаёт новый обработчик для аутентификации.
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register обрабатывает регистрацию.
func (h *AuthHandler) Register(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error: " + err.Error()})
		return
	}

	if len(body.Password) < 9 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 9 characters long"})
		return
	}

	userID, err := h.authService.RegisterUser(body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "userID": userID})
}

// Login обрабатывает логин.
func (h *AuthHandler) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error: " + err.Error()})
		return
	}

	token, err := h.authService.AuthenticateUser(body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("auth_token", token, 3600*24, "/", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func (h *AuthHandler) GetUser(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	user, err := h.authService.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}
