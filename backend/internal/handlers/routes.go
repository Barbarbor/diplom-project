package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(router *gin.Engine, db *sqlx.DB) {
	api := router.Group("/api")
	{
		api.GET("/surveys", GetSurveysHandler(db))
		api.POST("/surveys", CreateSurveyHandler(db))
	}
}
