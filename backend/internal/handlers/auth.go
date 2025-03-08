package handlers

import (
	"backend/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// RegisterHandler обрабатывает регистрацию с валидацией
func RegisterHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		// Если входные данные некорректны, Gin вернёт ошибку валидации
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error: " + err.Error()})
			return
		}

		// Дополнительная проверка пароля, если нужно (например, минимальная длина)
		if len(body.Password) < 9 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 9 characters long"})
			return
		}

		// Регистрируем пользователя
		userID, err := auth.RegisterUser(db, body.Email, body.Password)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "userID": userID})
	}
}

// LoginHandler обрабатывает логин
func LoginHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error: " + err.Error()})
			return
		}

		// Пытаемся аутентифицировать пользователя
		token, err := auth.AuthenticateUser(db, body.Email, body.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.SetCookie("auth_token", token, 3600*24, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
	}
}
func GetUserHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		var user struct {
			ID    int    `db:"id"`
			Email string `db:"email"`
		}

		err := db.Get(&user, "SELECT id, email FROM users WHERE id = $1", userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch user data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}
