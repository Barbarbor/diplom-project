package middleware

import (
	"backend/internal/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// QuestionMiddleware ищет вопрос по questionId из URL и сохраняет его в контекст.
func QuestionMiddleware(questionRepo repositories.QuestionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		questionIDStr := c.Param("questionId")
		questionID, err := strconv.Atoi(questionIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
			return
		}

		question, err := questionRepo.GetQuestionByID(questionID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		}

		// Сохраняем найденный вопрос в контексте
		c.Set("question", question)
		c.Next()
	}
}
