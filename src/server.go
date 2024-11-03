package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"

	"github.com/cityracingteam/data-acq-backend/environment"
	"github.com/cityracingteam/data-acq-backend/routes"
)

func main() {
	// Init goth for user authentication
	openidConnect, err := openidConnect.New(
		environment.GetEnvOrDefault("OPENID_CONNECT_KEY"),
		environment.GetEnvOrDefault("OPENID_CONNECT_SECRET"),
		environment.GetEnvOrDefault("ENDPOINT")+"/auth/openid-connect/callback",
		environment.GetEnvOrDefault("OPENID_CONNECT_DISCOVERY_URL"))

	if err != nil {
		fmt.Println(err)
		return
	}
	goth.UseProviders(openidConnect)

	// Connect to database

	// Create a gin engine instance
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Register GraphQL routes
	r.POST("/graphql", routes.GraphqlHandler())
	r.GET("/graphql/playground", routes.PlaygroundHandler())

	// Register authentication routes
	r.GET("/auth/:provider/callback", routes.AuthCallbackHandler())
	r.GET("/auth/:provider", routes.AuthHandler())

	r.Run()
}
