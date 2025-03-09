package handlers

import (
	"backend/internal/domain"
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
	userIDInterface, _ := c.Get("user_id")

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
	// Получаем данные опроса, установленные middleware
	surveyData, _ := c.Get("survey")

	survey, ok := surveyData.(*domain.Survey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid survey data"})
		return
	}

	// Получаем email создателя из контекста
	creator, exists := c.Get("surveyAuthor")
	if !exists {
		creator = "unknown"
	}

	c.JSON(http.StatusOK, gin.H{
		"survey": gin.H{
			"title":      survey.Title,
			"created_at": survey.CreatedAt,
			"updated_at": survey.UpdatedAt,
			"state":      survey.State,
			"creator":    creator, // email автора
		},
	})
}

func (h *SurveyHandler) GetSurveys(c *gin.Context) {
	// Получаем user_id из контекста (middleware-аутентификации)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	summaries, err := h.surveyService.GetSurveysByAuthor(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch surveys"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"surveys": summaries})
}
