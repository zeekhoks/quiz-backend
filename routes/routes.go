package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zeekhoks/quiz-backend/controllers"
	"github.com/zeekhoks/quiz-backend/middleware"
)

func GetRouter() *gin.Engine {

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")

	router.Use(cors.New(config))

	apiGroup := router.Group("/api")

	apiGroup.POST("/user", controllers.CreateNewUser())
	apiGroup.POST("/login", middleware.BasicAuth(), controllers.LoginHandler())

	apiGroup.POST("/questions", middleware.UserExtractor(), middleware.AdminCheck(), controllers.UploadQuestionHandler())
	apiGroup.GET("/questions", middleware.UserExtractor(), middleware.AdminCheck(), controllers.GetDisplayQuestionsByTopicHandler())

	apiGroup.GET("/topics", middleware.UserExtractor(), controllers.GetAllTopics())
	apiGroup.POST("/quiz", middleware.UserExtractor(), controllers.GenerateQuizHandler())
	apiGroup.POST("/quiz/:id/response", middleware.UserExtractor(), controllers.SubmitAnswerHandler())
	apiGroup.GET("/quiz/:id/result", middleware.UserExtractor(), controllers.QuizResultHandler())

	return router
}
