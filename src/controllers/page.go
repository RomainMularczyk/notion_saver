package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"notion_saver/src/models"
	"notion_saver/src/repositories"
	"notion_saver/src/services"
)

func GetPage(c *gin.Context) {
	id := c.Param("id")
	queryPage := repositories.QueryPage()
	page := queryPage.RetrieveOne(id)
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": http.StatusOK,
			"data": page,
		},
	)
}

func GetAllPages(c *gin.Context) {
	queryPage := repositories.QueryPage()
	allPages := queryPage.RetrieveAll()
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": http.StatusOK,
			"data": allPages,
		},
	)
}

func AddPages(c *gin.Context) {
	var pages []models.Page
	if err := c.ShouldBindJSON(&pages); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "An error occurred with the data format.",
			},
		)
	} else {
		queryPage := repositories.QueryPage()
		queryPage.CreateMany(pages)
		c.JSON(
			http.StatusCreated,
			gin.H{
				"code":    http.StatusCreated,
				"message": "Pages were created successfully.",
				"data":    pages,
			},
		)
	}
}

func AddPage(c *gin.Context) {
	var page models.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "An error occurred with the data format.",
			},
		)
	} else {
		services.CreatePage(page)
		c.JSON(
			http.StatusCreated,
			gin.H{
				"code":    http.StatusCreated,
				"message": "Page was created successfully.",
				"data":    page,
			},
		)
	}
}

func DeletePage(c *gin.Context) {
	id := c.Param("id")
	queryPage := repositories.QueryPage()
	err := queryPage.DeleteOne(id)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": fmt.Sprintf("An error occurred when deleting the page with Id: %s.", id),
			},
		)
	} else {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":    http.StatusOK,
				"message": "Page was deleted successfully.",
			},
		)
	}
}
