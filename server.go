package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cityracingteam/data-acq-backend/routes"
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

	r.POST("/graphql", routes.GraphqlHandler())
	r.GET("/graphql/playground", routes.PlaygroundHandler())
	r.Run()
}
