package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"

	"github.com/cityracingteam/data-acq-backend/database"
	"github.com/cityracingteam/data-acq-backend/environment"
	"github.com/cityracingteam/data-acq-backend/graph/resolver"
	"github.com/cityracingteam/data-acq-backend/middleware"
	"github.com/cityracingteam/data-acq-backend/routes"
	"github.com/cityracingteam/data-acq-backend/util/jwt"
)

func main() {
	// Check that we are not in release mode without a session secret
	_, hasKey := os.LookupEnv("SESSION_SECRET")
	ginMode := os.Getenv("GIN_MODE")

	if !hasKey && ginMode == "release" {
		// quit the program if this is the case
		fmt.Println("[error]: GIN_MODE set to release but SESSION_SECRET not specified. This is an unsafe configuration.")
		return
	}
	// continue

	// Init goth for user authentication
	openidConnect, err := openidConnect.New(
		environment.GetEnvOrDefault("OPENID_CONNECT_KEY"),
		environment.GetEnvOrDefault("OPENID_CONNECT_SECRET"),
		environment.GetCallbackUri(),
		environment.GetEnvOrDefault("OPENID_CONNECT_DISCOVERY_URL"))

	if err != nil {
		fmt.Println(err)
		return
	}
	goth.UseProviders(openidConnect)

	// Override defaults for goth's cookiestore
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.MaxAge(86400)           // one day
	store.Options.Path = "/"      // default
	store.Options.HttpOnly = true // default, should always be enabled
	store.Options.Secure = (environment.GetEnvOrDefault("ENDPOINT_SCHEME") == "https")
	store.Options.Domain = environment.GetEnvOrDefault("DOMAIN")

	// Connect to the database and place the handle in the graphql resolver
	// So that is accessible when executing graphql requests in ctx.
	resolver := &resolver.Resolver{}
	db := database.Connect()
	resolver.UpdateDb(db)

	// Setup JWT (loads from database)
	jwt.Init(db)

	// Create a gin engine instance
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Register authentication routes
	r.GET("/auth/:provider/callback", routes.AuthCallbackHandler())
	r.GET("/auth/:provider", routes.AuthHandler())

	// Register routes that require the caller/user to be authenticated

	// Create a group for authorized routes
	authorized := r.Group("/")
	// Apply the requireAuth middleware to the authorized group
	authorized.Use(middleware.RequireAuth(db))
	{
		// Any routes registered here will require the user to be authenticated

		// Register GraphQL routes
		authorized.POST("/graphql", routes.GraphqlHandler())
		authorized.GET("/graphql/playground", routes.PlaygroundHandler())
	}

	r.Run()
}
