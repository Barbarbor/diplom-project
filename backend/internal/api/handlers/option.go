package handlers

import (
	"backend/internal/domain"
	option "backend/internal/services/option_service"
	"backend/pkg/i18n"
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

// PATCH /api/.../option/:optionId/order
func (h *OptionHandler) UpdateOptionOrder(c *gin.Context) {
	optData, _ := c.Get("option")
	opt := optData.(*domain.OptionTemp)

	var body struct {
		NewOrder int `json:"new_order" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T("option.handler.invalidData")})
		return
	}
	// Проверяем, что новый порядок допустим (не равен текущему и не выходит за границы)
	if body.NewOrder <= 0 || body.NewOrder == opt.OptionOrder {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid new order value"})
		return
	}
	if err := h.optionService.UpdateOptionOrder(opt, body.NewOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// DELETE /api/.../option/:optionId
func (h *OptionHandler) DeleteOption(c *gin.Context) {
	questionData, _ := c.Get("question")
	question := questionData.(*domain.SurveyQuestionTemp)
	optData, _ := c.Get("option")
	opt := optData.(*domain.OptionTemp)

	if err := h.optionService.DeleteOption(opt, question.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// PATCH /api/.../option/:optionId
func (h *OptionHandler) UpdateOption(c *gin.Context) {
	questionData, _ := c.Get("question")
	question := questionData.(*domain.SurveyQuestionTemp)
	optData, _ := c.Get("option")
	opt := optData.(*domain.OptionTemp)

	var body struct {
		Label string `json:"label"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.T("option.handler.invalidData")})
		return
	}
	if err := h.optionService.UpdateOptionLabel(opt, body.Label, question.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
