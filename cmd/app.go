package cmd

import (
	"log"
	"os"

	"github.com/JerryJeager/ohara-be/manualwire"
	"github.com/JerryJeager/ohara-be/middleware"
	"github.com/gin-gonic/gin"
)

func ExecuteApiRoutes() {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome",
		})
	})

	userController := manualwire.GetUserController()
	documentController := manualwire.GetDocumentController()
	websiteController := manualwire.GetWebsiteController()

	api := router.Group("/api/v1")
	users := api.Group("/users")
	documents := api.Group("/documents")
	websites := api.Group("/websites")

	users.POST("/signup", userController.CreateUser)
	users.POST("/verify-email", userController.VerifyUserEmail)
	users.POST("/login", userController.Login)

	documents.GET("/embed", documentController.EmbedDocument)
	documents.GET("/query/:website_id", documentController.QueryDocument)
	documents.GET("/chunk", documentController.ChunkDocument)

	websites.POST("", websiteController.CreateWebsite)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panic("failed to run server")
	}
}
