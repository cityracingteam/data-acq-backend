package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database

	// Create a gin engine instance
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
