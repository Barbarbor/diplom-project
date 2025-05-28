package api

import (
	"backend/internal/api/handlers"
	"backend/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	profileHandler *handlers.ProfileHandler,
	surveyHandler *handlers.SurveyHandler,
	questionHandler *handlers.QuestionHandler,
	optionHandler *handlers.OptionHandler,
	interviewHandler *handlers.InterviewHandler,
	surveyAccessMiddleware gin.HandlerFunc,
	questionMiddleware gin.HandlerFunc,
	optionMiddleware gin.HandlerFunc,
	interviewMiddleware gin.HandlerFunc,
) {
	api := router.Group("/api", middleware.I18nMiddleware())
	{
		// Authorization routes
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.GET("/user", middleware.AuthMiddleware(), authHandler.GetUser)
		}

		// Profile routes
		profileRoutes := api.Group("/profile", middleware.AuthMiddleware())
		{
			profileRoutes.GET("", profileHandler.GetProfile)
			profileRoutes.PUT("", profileHandler.UpdateProfile)
		}

		// Survey routes
		surveyRoutes := api.Group("/surveys", middleware.AuthMiddleware())
		{
			surveyRoutes.POST("", surveyHandler.CreateSurvey)
			surveyRoutes.GET("", surveyHandler.GetSurveys)
			surveyProtected := surveyRoutes.Group("/:hash", surveyAccessMiddleware)
			{
				surveyProtected.GET("", surveyHandler.GetSurvey)
				surveyProtected.PATCH("", surveyHandler.UpdateSurvey)
				surveyProtected.POST("/publish", surveyHandler.PublishSurvey)
				surveyProtected.PUT("/restore", surveyHandler.RestoreSurvey)
				questionRoutes := surveyProtected.Group("/question")
				{
					questionRoutes.POST("", questionHandler.CreateQuestion)
					questionProtected := questionRoutes.Group("/:questionId", questionMiddleware)
					{
						questionProtected.POST("/option", optionHandler.CreateOption)
						questionProtected.PATCH("", questionHandler.UpdateQuestion)
						questionProtected.PATCH("/type", questionHandler.UpdateQuestionType)
						questionProtected.PATCH("/order", questionHandler.UpdateQuestionOrder)
						questionProtected.PATCH("/extra_params", questionHandler.UpdateExtraParams)
						questionProtected.PUT("/restore", questionHandler.RestoreQuestion)
						questionProtected.DELETE("", questionHandler.DeleteQuestion)
						optionProtected := questionProtected.Group("/option/:optionId", optionMiddleware)
						{
							optionProtected.PATCH("/order", optionHandler.UpdateOptionOrder)
							optionProtected.DELETE("", optionHandler.DeleteOption)
							optionProtected.PATCH("", optionHandler.UpdateOption)
						}
					}
				}
			}
		}

		interviewRoutes := api.Group("/interview/:hash")
		{
			interviewRoutes.POST("/start", interviewHandler.StartInterview)
			interviewGroup := interviewRoutes.Group("", interviewMiddleware)
			{
				interviewGroup.GET("/survey", interviewHandler.GetSurveyWithAnswers)
				interviewRoutes.PATCH("/:questionId/answer", interviewHandler.UpdateQuestionAnswer)
				interviewRoutes.POST("/finish", interviewHandler.FinishInterview)
			}
		}
		// statsRoutes := api.Group("/stats/:hash")
		// {
		// 	statsRoutes.GET("", statsHandler.GetSurveyStats)
		// }
	}
}
