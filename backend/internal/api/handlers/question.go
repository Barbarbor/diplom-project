package handlers

import (
	"backend/internal/domain"
	question "backend/internal/services/question_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// QuestionHandler обрабатывает запросы, связанные с вопросами
type QuestionHandler struct {
	service *question.QuestionService
}

// NewQuestionHandler создаёт новый обработчик вопросов
func NewQuestionHandler(service *question.QuestionService) *QuestionHandler {
	return &QuestionHandler{service: service}
}

// CreateQuestion создаёт новый вопрос в опросе
func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	surveyData, _ := c.Get("survey")

	survey, ok := surveyData.(*domain.Survey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid survey data"})
		return
	}
	questionType := c.Query("type")
	if questionType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Question type is required"})
		return
	}

	question, err := h.service.CreateQuestion(survey.ID, domain.QuestionType(questionType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"question": question})
}
