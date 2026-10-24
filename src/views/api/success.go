package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// API response on successful read access
func ReadSuccess(c *gin.Context, data any) *gin.Context {
	c.JSON(
		http.StatusOK,
		gin.H{
			"data": data,
		},
	)
	return c
}

// API response on successful creation access
func CreateSuccess(c *gin.Context, data any) *gin.Context {
	c.JSON(
		http.StatusCreated,
		gin.H{
			"data": data,
		},
	)
	return c
}

// API response on successful delete access
func DeleteSuccess(c *gin.Context, dataType string) *gin.Context {
	c.JSON(
		http.StatusOK,
		gin.H{
			"message": fmt.Sprintf("%s was delete successfully.", dataType),
		},
	)
	return c
}
