package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"notion_saver/src/queries"
)

func Notion(c *gin.Context) {
	result := queries.GetAllPages()
	fmt.Println(result)
}
