package routes

import (
	"notion_saver/src/controllers"
	"notion_saver/src/utils"

	"github.com/gin-gonic/gin"
)

// Define all the API routes
func Router(server *utils.Server, router *gin.Engine) {
	routes := router.Group("/api/v1")
	// Saves
	{
		routes.POST("/saves", controllers.NewSave(server))
	}
	// Test
	{
		routes.GET("/test", controllers.TestHammer(server))
	}
}
