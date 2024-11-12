package utils

import (
    "github.com/gin-gonic/gin"
    "github.com/a-h/templ"
)

func Render(
    c *gin.Context, 
    status int,
    template templ.Component,
) error {
    c.Status(status)
    return template.Render(c.Request.Context(), c.Writer)
}
