package handlers

import (
	survey "backend/internal/services/survey_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SurveyHandler структурирует обработку запросов для опросов.
type SurveyHandler struct {
	surveyService *survey.SurveyService
}

// NewSurveyHandler создаёт новый обработчик опросов.
func NewSurveyHandler(surveyService *survey.SurveyService) *SurveyHandler {
	return &SurveyHandler{
		surveyService: surveyService,
	}
}

// CreateSurvey создает новый опрос и возвращает его hash.
func (h *SurveyHandler) CreateSurvey(c *gin.Context) {
	// Извлекаем user_id из контекста (установлено middleware-аутентификации)
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

	survey, err := h.surveyService.CreateSurvey(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Survey created successfully",
		"hash":    survey.Hash,
	})
}

// GetSurvey возвращает опрос по hash и добавляет поле creator (user_email) из контекста.
func (h *SurveyHandler) GetSurvey(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Survey hash is required"})
		return
	}

	survey, err := h.surveyService.GetSurveyByHash(hash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"survey": gin.H{
			"title":      survey.Survey.Title,
			"created_at": survey.Survey.CreatedAt,
			"updated_at": survey.Survey.UpdatedAt,
			"state":      survey.Survey.State,
			"creator":    survey.CreatorEmail,
		},
	})
}
