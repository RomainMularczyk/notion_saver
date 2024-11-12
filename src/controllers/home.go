package controllers

import (
    "github.com/gin-gonic/gin"
    "notion_saver/src/templates"
    "notion_saver/src/utils"
)

func Home(c *gin.Context) {
    utils.Render(c, 200, templates.Home())
}
