package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// API response on server error linked to the database connection
func ServerDatabaseError(c *gin.Context) *gin.Context {
	c.JSON(
		http.StatusUnprocessableEntity,
		gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "An error occurred with the database connection.",
		},
	)
	return c
}

// API response on server error linked to data handling issues
func ServerDataError(c *gin.Context, err error) *gin.Context {
	c.JSON(
		http.StatusInternalServerError,
		gin.H{
			"code": http.StatusInternalServerError,
			"message": fmt.Sprintf(
				"An error occurred when creating data: %v",
				err,
			),
		},
	)
	return c
}

// API response on server error linked to WebSocket connection
func ServerWebsocketError(c *gin.Context) *gin.Context {
	c.JSON(
		http.StatusInternalServerError,
		gin.H{
			"code":    http.StatusInternalServerError,
			"message": "An error occurred when trying to initialize websocket connection.",
		},
	)
	return c
}

// API response on client error linked to the data format provided
func ClientDataFormatError(c *gin.Context) *gin.Context {
	c.JSON(
		http.StatusUnprocessableEntity,
		gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "An error occurred with the data format provided.",
		},
	)
	return c
}
