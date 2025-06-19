package handlers

import (
	"backend/internal/domain"
	interview "backend/internal/services/interview_service"
	survey "backend/internal/services/survey_service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InterviewHandler struct {
	service       *interview.InterviewService
	surveyService *survey.SurveyService
}

func NewInterviewHandler(service *interview.InterviewService, surveyService *survey.SurveyService) *InterviewHandler {
	return &InterviewHandler{service: service, surveyService: surveyService}
}

// StartInterview handles the request to start an interview
func (h *InterviewHandler) StartInterview(c *gin.Context) {
	hash := c.Param("hash")
	interviewID := c.Query("interviewId")

	if interviewID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "interviewId is required"})
		return
	}
	isDemoStr := c.Query("isDemo")
	isDemo, err := strconv.ParseBool(isDemoStr)
	if err != nil {
		isDemo = false // Default to false if parsing fails
	}

	err = h.service.StartInterview(hash, interviewID, isDemo)

	if err != nil {
		switch err {
		case domain.ErrSurveyNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found"})
		case domain.ErrInterviewAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Interview already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start interview"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Interview started successfully"})
}

func (h *InterviewHandler) GetSurveyWithAnswers(c *gin.Context) {
	// Получаем interview из контекста
	interview, exists := c.Get("interview")
	if !exists {
		c.JSON(500, gin.H{"error": "Interview not found in context"})
		return
	}

	// Приводим interview к типу *domain.SurveyInterview
	surveyInterview, ok := interview.(*domain.SurveyInterview)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid interview type in context"})
		return
	}

	// Проверка статуса интервью
	if surveyInterview.Status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Вы уже завершили опрос"})
		return
	}

	// Извлекаем параметры
	surveyID := surveyInterview.SurveyID
	interviewID := surveyInterview.ID
	isDemo := surveyInterview.IsDemo

	// Вызываем сервис
	questions, err := h.surveyService.GetSurveyQuestionsWithAnswers(surveyID, interviewID, isDemo)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get survey questions with answers"})
		return
	}

	// Возвращаем результат клиенту
	c.JSON(200, gin.H{"questions": questions})
}

func (h *InterviewHandler) UpdateQuestionAnswer(c *gin.Context) {
	// Извлекаем интервью из контекста
	interview, exists := c.Get("interview")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Interview not found in context"})
		return
	}
	surveyInterview, ok := interview.(*domain.SurveyInterview)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid interview type"})
		return
	}

	// Проверка статуса интервью
	if surveyInterview.Status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Вы уже завершили опрос"})
		return
	}

	// Извлекаем questionId из URL
	questionIDStr := c.Param("questionId")
	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	// Извлекаем тело запроса с ответом
	var requestBody struct {
		Answer string `json:"answer"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Вызываем сервис для обновления или создания ответа
	err = h.surveyService.UpdateQuestionAnswer(surveyInterview.ID, questionID, requestBody.Answer, surveyInterview.IsDemo)
	if err != nil {
		switch err {
		case domain.ErrQuestionNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update answer"})
		}
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Answer updated successfully"})
}

func (h *InterviewHandler) FinishInterview(c *gin.Context) {
	// Извлекаем интервью из контекста
	interview, exists := c.Get("interview")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "interview not found in context"})
		return
	}
	surveyInterview, ok := interview.(*domain.SurveyInterview)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid interview type"})
		return
	}

	// Проверка статуса интервью
	if surveyInterview.Status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Вы уже завершили опрос"})
		return
	}

	// Завершаем интервью через сервис
	if err := h.surveyService.FinishInterview(surveyInterview.SurveyID, surveyInterview.ID, surveyInterview.IsDemo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "interview completed successfully"})
}
