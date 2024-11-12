package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"notion_saver/src/models"
	"notion_saver/src/repositories"
	"notion_saver/src/services"
)

func GetSave(c *gin.Context) {

}

func GetLatestSave(c *gin.Context) {
	querySave := repositories.QuerySave()
	latestSave, _ := querySave.RetrieveLatest()
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": http.StatusOK,
			"data": latestSave,
		},
	)
}

func GetAllSaves(c *gin.Context) {
	querySave := repositories.QuerySave()
	allSaves := querySave.RetrieveAll()
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": http.StatusOK,
			"data": allSaves,
		},
	)
}

func AddSave(c *gin.Context) {
	var save models.Save
	if err := c.ShouldBindJSON(&save); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "An error occurred with the data format.",
			},
		)
	} else {
		services.CreateSave(save)
		c.JSON(
			http.StatusCreated,
			gin.H{
				"code":    http.StatusCreated,
				"message": "Save was created successfully.",
				"data":    save,
			},
		)
	}
}

func DeleteSave(c *gin.Context) {

}
