package handlers

import (
	"backend/internal/domain"
	option "backend/internal/services/option_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OptionHandler обрабатывает запросы, связанные с вариантами ответов.
type OptionHandler struct {
	optionService *option.OptionService
}

// NewOptionHandler создает новый обработчик вариантов ответов.
func NewOptionHandler(optionService *option.OptionService) *OptionHandler {
	return &OptionHandler{optionService: optionService}
}

// CreateOption обрабатывает POST-запрос на создание варианта ответа.
// Маршрут: POST /api/surveys/:hash/question/:questionId/option
func (h *OptionHandler) CreateOption(c *gin.Context) {
	// Извлекаем вопрос из контекста, установленный QuestionMiddleware
	questionData, _ := c.Get("question")

	question, ok := questionData.(*domain.SurveyQuestionTemp)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid question data"})
		return
	}

	// Используем тип вопроса из найденного объекта
	option, err := h.optionService.CreateOption(question.ID, question.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"option": option})
}
