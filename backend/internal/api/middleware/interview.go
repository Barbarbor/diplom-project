package middleware

import (
	"backend/internal/domain"
	"backend/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InterviewMiddleware(repo repositories.InterviewRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		interviewID := c.Query("interviewId")

		if interviewID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "interviewId is required"})
			return
		}

		interview, err := repo.GetInterviewByID(interviewID)
		if err != nil {
			if err == domain.ErrInterviewNotFound {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Interview not found"})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch interview"})
			}
			return
		}

		// Add the interview to the context for downstream handlers
		c.Set("interview", interview)
		c.Next()
	}
}
