package routes

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/cityracingteam/data-acq-backend/graph"
	"github.com/cityracingteam/data-acq-backend/graph/resolver"
	"github.com/gin-gonic/gin"
)

// Defining the GraphQL handler
func GraphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
