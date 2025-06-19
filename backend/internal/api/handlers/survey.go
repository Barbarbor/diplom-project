package handlers

import (
	"backend/internal/domain"
	survey "backend/internal/services/survey_service"
	"backend/pkg/i18n"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// SurveyHandler структурирует обработку запросов для опросов.
type SurveyHandler struct {
	surveyService *survey.SurveyService
	db            *sqlx.DB
}

// NewSurveyHandler создаёт новый обработчик опросов.
func NewSurveyHandler(surveyService *survey.SurveyService, db *sqlx.DB) *SurveyHandler {
	return &SurveyHandler{surveyService: surveyService, db: db}
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
	surveyData, _ := c.Get("survey")

	survey, ok := surveyData.(*domain.Survey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T("survey.handler.invalidData")})
		return
	}

	// Если middleware уже установила email автора (например, "surveyAuthor")
	creator, exists := c.Get("surveyAuthor")
	if !exists {
		creator = "unknown"
	}

	// Получаем список вопросов для опроса
	questions, err := h.surveyService.GetQuestionsForSurvey(survey.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T("question.handler.notFound")})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"survey": gin.H{
			"title":      survey.Title,
			"created_at": survey.CreatedAt,
			"updated_at": survey.UpdatedAt,
			"hash":       survey.Hash,
			"state":      survey.State,
			"creator":    creator,
			"questions":  questions,
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

// PATCH /api/surveys/:hash
func (h *SurveyHandler) UpdateSurvey(c *gin.Context) {
	var body struct {
		Title string `json:"title"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	surveyData, _ := c.Get("survey")
	survey := surveyData.(*domain.Survey) // из middleware
	if err := h.surveyService.UpdateSurvey(survey.ID, body.Title); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// POST /api/surveys/:hash/publish
func (h *SurveyHandler) PublishSurvey(c *gin.Context) {
	surveyData, _ := c.Get("survey")
	survey := surveyData.(*domain.Survey)
	if err := h.surveyService.PublishSurvey(survey.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// PUT /api/surveys/:hash/restore
func (h *SurveyHandler) RestoreSurvey(c *gin.Context) {
	// middleware уже положил *domain.Survey
	raw, _ := c.Get("survey")
	surveyObj := raw.(*domain.Survey)

	if err := h.surveyService.RestoreSurveyByID(surveyObj.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetSurveyStats - хэндлер для маршрута /stats/:hash
func (h *SurveyHandler) GetSurveyStats(c *gin.Context) {
	surveyData, _ := c.Get("survey")
	survey := surveyData.(*domain.Survey)
	stats, err := h.surveyService.GetSurveyStats(survey.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetSurveyAccessHandler handles GET /access
func (h *SurveyHandler) GetSurveyAccessHandler(c *gin.Context) {
	surveyData, _ := c.Get("survey")
	survey := surveyData.(*domain.Survey)
	emails, err := h.surveyService.GetAccessList(survey.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, emails)
}

// AddSurveyAccessHandler handles POST /access?email=<userEmailToAdd>
func (h *SurveyHandler) AddSurveyAccessHandler(c *gin.Context) {
	surveyData, _ := c.Get("survey")
	survey := surveyData.(*domain.Survey)
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	err := h.surveyService.AddEditAccess(survey.ID, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// RemoveSurveyAccessHandler handles DELETE /access?email=<userEmailToDelete>
func (h *SurveyHandler) RemoveSurveyAccessHandler(c *gin.Context) {
	surveyData, _ := c.Get("survey")
	survey := surveyData.(*domain.Survey)
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, _ := userIDInterface.(int)
	creator, exists := c.Get("surveyAuthor")
	creatorParsed := creator.(string)
	if !exists {
		creator = "unknown"
	}

	err := h.surveyService.RemoveEditAccess(survey.ID, email, userID, creatorParsed)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
