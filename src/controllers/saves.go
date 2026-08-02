package controllers

import (
	"notion_saver/src/services"
	"notion_saver/src/utils"
	"notion_saver/src/views/api"

	"github.com/gin-gonic/gin"
)

func NewSave(server *utils.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		saveService := services.NewSaveService(server)
		newSave, err := saveService.New(server)

		if err != nil {
			api.ServerDataError(c, err)
		}

		api.ReadSuccess(c, newSave)
	}
}
