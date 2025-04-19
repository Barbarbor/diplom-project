package handlers

import (
	"backend/internal/domain"
	question "backend/internal/services/question_service"
	"backend/pkg/i18n"
	"net/http"
	"strconv"

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T("question.handler.invalidData")})
		return
	}
	questionType := c.Query("type")
	if questionType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T("question.handler.invalidType")})
		return
	}

	question, err := h.service.CreateQuestion(survey.ID, domain.QuestionType(questionType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.T("question.handler.invalidType")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"question": question})
}

// UpdateQuestionType обрабатывает запрос на обновление типа вопроса.
// Маршрут: PUT /api/surveys/:hash/question/:questionId/type?newType=...
func (h *QuestionHandler) UpdateQuestionType(c *gin.Context) {
	newTypeStr := c.Query("newType")
	if newTypeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New question type is required"})
		return
	}
	newType := domain.QuestionType(newTypeStr)

	// Извлекаем текущий вопрос из контекста (установлен в middleware)
	qData, _ := c.Get("question")

	question, ok := qData.(*domain.SurveyQuestionTemp)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid question data"})
		return
	}

	if newType == question.Type {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can`t change question type to the same type"})
		return
	}
	// Вызываем сервис для обновления типа, передавая текущее состояние
	if err := h.service.UpdateQuestionType(question.ID, newType, string(question.QuestionState)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question type updated successfully"})
}

// UpdateQuestionLabelHandler обновляет только label вопроса.
// Маршрут: PUT /api/surveys/:hash/question/:questionId/label
func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	questionIDStr := c.Param("questionId")
	questionID, _ := strconv.Atoi(questionIDStr)

	var body struct {
		Label string `json:"label" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	err := h.service.UpdateQuestion(questionID, body.Label)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question label updated successfully"})
}

// UpdateQuestionOrderHandler обновляет порядок вопроса.
// Маршрут: PUT /api/surveys/:hash/question/:questionId/order
func (h *QuestionHandler) UpdateQuestionOrder(c *gin.Context) {
	// Извлекаем объект вопроса из контекста, который должен быть установлен QuestionMiddleware
	questionData, _ := c.Get("question")

	q, ok := questionData.(*domain.SurveyQuestionTemp)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid question data"})
		return
	}

	// Допустим, surveyID можно извлечь также из контекста (например, как часть объекта опроса)
	surveyData, _ := c.Get("survey")

	survey, ok := surveyData.(*domain.Survey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid survey data"})
		return
	}

	var body struct {
		NewOrder int `json:"new_order" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Проверяем, что новый порядок допустим (не равен текущему и не выходит за границы)
	if body.NewOrder <= 0 || body.NewOrder == q.QuestionOrder {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid new order value"})
		return
	}

	// Передаем параметры: questionID, новый порядок, текущий порядок, surveyID.
	if err := h.service.UpdateQuestionOrder(q.ID, body.NewOrder, q.QuestionOrder, survey.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question order updated successfully"})
}

// DELETE /api/surveys/:hash/question/:questionId
func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	qData, _ := c.Get("question")
	q := qData.(*domain.SurveyQuestionTemp)

	if err := h.service.DeleteQuestion(q.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Question deleted"})
}

// PUT /api/surveys/:hash/question/:questionId/restore
func (h *QuestionHandler) RestoreQuestion(c *gin.Context) {
	qData, _ := c.Get("question")
	q := qData.(*domain.SurveyQuestionTemp)

	if err := h.service.RestoreQuestion(q.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Question restored"})
}
