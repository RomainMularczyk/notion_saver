package controllers

import (
	"notion_saver/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestHammer(server *utils.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		notion := utils.NewNotion(server)
		for i := 0; i < 40; i++ {
			id := "39b6bbca-9c5e-80dd-9aaf-cc6bcd825847"
			notion.GetPageBlocks(uuid.New(), &id)
		}
	}
}
