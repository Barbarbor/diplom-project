package middleware

import (
	"backend/internal/domain"
	"backend/internal/repositories"
	"backend/pkg/i18n"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OptionMiddleware ищет опцию по optionId из URL и сохраняет её в контекст.
// Предполагается, что до этого уже отработал QuestionMiddleware и в контексте лежит "question" типа *domain.SurveyQuestionTemp.
func OptionMiddleware(optionRepo repositories.OptionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Сначала достаём вопрос из контекста
		questionData, exists := c.Get("question")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": i18n.T("option.handler.invalidQuestionContext"),
			})
			return
		}
		question, ok := questionData.(*domain.SurveyQuestionTemp)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": i18n.T("option.handler.invalidQuestionContext"),
			})
			return
		}

		// Достаём optionId из URL
		optionIDStr := c.Param("optionId")
		optionID, err := strconv.Atoi(optionIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": i18n.T("option.handler.invalidOptionID"),
			})
			return
		}

		// Ищем опцию в репозитории, ограничиваясь текущим question.ID
		opt, err := optionRepo.GetOptionById(question.ID, optionID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": i18n.T("option.handler.notFound"),
			})
			return
		}

		// Кладём найденную опцию в контекст
		c.Set("option", opt)
		c.Next()
	}
}
