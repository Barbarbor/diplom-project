package handlers

import (
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// CreateSurveyHandler контроллер для создания опроса
func CreateSurveyHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем user_id из контекста (установлено middleware-аутентификации)
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userID, ok := userIDInterface.(int)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user id"})
			return
		}

		survey, err := services.CreateSurvey(db, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Survey created successfully",
			"hash":    survey.Hash,
		})
	}
}

// GetSurveyHandler возвращает опрос по hash и email автора
func GetSurveyHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Param("hash")
		if hash == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Survey hash is required"})
			return
		}

		survey, email, err := services.GetSurveyByHash(db, hash)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found"})
			return
		}

		// Возвращаем только нужные поля
		c.JSON(http.StatusOK, gin.H{
			"survey": gin.H{
				"title":      survey.Title,
				"created_at": survey.CreatedAt,
				"updated_at": survey.UpdatedAt,
				"state":      survey.State,
				"creator":    email,
			},
		})
	}
}
