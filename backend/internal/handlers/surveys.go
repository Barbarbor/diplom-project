package handlers

import (
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetSurveysHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		surveys, err := services.GetSurveys(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, surveys)
	}
}
