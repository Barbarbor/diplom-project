package handlers

import (
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(router *gin.Engine, db *sqlx.DB) {
	api := router.Group("/api")
	{
		// Authorization routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", RegisterHandler(db))
			auth.POST("/login", LoginHandler(db))
			auth.GET("/user", middleware.AuthMiddleware(), GetUserHandler(db))
		}

		// Profile routes
		profile := api.Group("/profile", middleware.AuthMiddleware())
		{
			profile.GET("", GetProfileHandler(db))
			profile.PUT("", UpdateProfileHandler(db))
		}
		surveys := api.Group("/surveys", middleware.AuthMiddleware())
		{
			surveys.POST("", CreateSurveyHandler(db))
			surveys.GET("/:hash", GetSurveyHandler(db))
		}
	}
}
