package handlers

import (
	"backend/internal/models"
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

func CreateSurveyHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var survey models.Survey
		if err := c.ShouldBindJSON(&survey); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Создание нового опроса
		if err := services.CreateSurvey(db, survey); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Возвращаем успешно созданный опрос
		c.JSON(http.StatusCreated, survey)
	}
}
