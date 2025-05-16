package middleware

import (
	"backend/internal/domain"
	"backend/internal/repositories"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// QuestionMiddleware ищет вопрос по questionId из URL и сохраняет его в контекст.
func QuestionMiddleware(questionRepo repositories.QuestionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		questionIDStr := c.Param("questionId")
		surveyData, _ := c.Get("survey")
		survey, ok := surveyData.(*domain.Survey)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid survey data"})
			return
		}

		surveyID := survey.ID

		questionID, err := strconv.Atoi(questionIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
			return
		}
		fmt.Print("questionId", questionID, " surveyId", surveyID)
		question, err := questionRepo.GetQuestionByID(questionID, surveyID)
		fmt.Println(" question", question)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		}

		// Сохраняем найденный вопрос в контексте
		c.Set("question", question)
		c.Next()
	}
}
