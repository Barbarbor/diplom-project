package handlers

import (
	"backend/internal/domain"
	survey "backend/internal/services/survey_service"
	"backend/pkg/redisclient"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

func (h *SurveyHandler) GetSurvey(c *gin.Context) {
	// Извлекаем опрос из контекста, установленный middleware
	surveyData, exists := c.Get("survey")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found in context"})
		return
	}
	survey, ok := surveyData.(*domain.Survey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid survey data"})
		return
	}

	// Если middleware уже установила email автора (например, "surveyAuthor")
	creator, exists := c.Get("surveyAuthor")
	if !exists {
		creator = "unknown"
	}

	// Формируем ключ для Redis. Здесь используем "survey:<id>"
	cacheKey := fmt.Sprintf("survey:%d", survey.ID)

	cachedData, err := redisclient.Client.Get(redisclient.Ctx, cacheKey).Bytes()
	if err == nil && len(cachedData) > 0 {
		c.Data(http.StatusOK, "application/json", cachedData)
		return
	}

	// Если в кеше нет, получаем список вопросов из сервиса
	questions, err := h.surveyService.GetQuestionsForSurvey(survey.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}

	// Формируем итоговый ответ
	responseBody := gin.H{
		"survey": gin.H{
			"title":      survey.Title,
			"created_at": survey.CreatedAt,
			"updated_at": survey.UpdatedAt,
			"hash":       survey.Hash,
			"state":      survey.State,
			"creator":    creator,
			"questions":  questions,
		},
	}

	// Сериализуем ответ в JSON
	responseJSON, err := json.Marshal(responseBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal response"})
		return
	}

	// Сохраняем данные в Redis с заданным TTL (например, 5 минут)
	redisclient.Client.Set(redisclient.Ctx, cacheKey, responseJSON, 24*31*365*time.Hour)

	// Отправляем сформированный ответ
	c.Data(http.StatusOK, "application/json", responseJSON)
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
