package middleware

import (
	"backend/internal/repositories"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SurveyAccessMiddleware проверяет, есть ли у пользователя доступ к опросу
func SurveyAccessMiddleware(surveyRepo repositories.SurveyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Param("hash")
		fmt.Print("hash", hash)
		if hash == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Survey hash is required"})
			return
		}

		// Получаем опрос по хэшу
		survey, email, err := surveyRepo.GetSurveyByHash(hash)
		fmt.Print("survey", survey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Survey not found"})
			return
		}

		// Проверяем, есть ли у пользователя доступ к опросу
		userID, _ := c.Get("user_id")

		hasAccess, err := surveyRepo.CheckUserAccess(userID.(int), survey.ID)

		fmt.Print("access?", hasAccess)
		if err != nil || !hasAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		// Сохраняем информацию об опросе в контексте и информацию о создателе опроса (email)
		c.Set("survey", survey)
		c.Set("surveyAuthor", email)
		c.Next()
	}
}
