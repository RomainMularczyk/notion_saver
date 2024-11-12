package routes

import (
	"github.com/gin-gonic/gin"
	"notion_saver/src/controllers"
)

// Define all the routes in the application
func Router() {
	router := gin.Default()
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
