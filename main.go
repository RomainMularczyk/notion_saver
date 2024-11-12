package main

import (
	"github.com/gin-gonic/gin"
	"notion_saver/src/controllers"
	"notion_saver/src/database"
)

func main() {
	// Create Gin engine
	router := gin.Default()

	// Migrate databse schemas
	database.MigrateSchemas()

	// Define routing rules
	routes := router.Group("/")
	{
		routes.GET("/", controllers.Home)
		routes.GET("/notion", controllers.Notion)
	}
	// Saves
	{
		routes.POST("/saves", controllers.AddSave)
		routes.GET("/saves", controllers.GetAllSaves)
		routes.GET("/saves/latest", controllers.GetLatestSave)
	}
	// Pages
	{
		routes.POST("/pages/:id", controllers.AddPage)
		routes.POST("/pages", controllers.AddPages)
		routes.GET("/pages/:id", controllers.GetPage)
		routes.GET("/pages", controllers.GetAllPages)
		routes.DELETE("/pages/:id", controllers.DeletePage)
	}
	router.Run(":8080")
}
