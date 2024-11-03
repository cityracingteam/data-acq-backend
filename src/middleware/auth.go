package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to fetch the user without a full re-authenticate
		if user, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
			// user is authorized, continue

			// Set goth.User variable inside context
			// so we can access it from inside a request handler
			c.Set("user", user)

			c.Next() // pass through to next middleware (will execute request)
		} else {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "route requires authorization",
			})
		}
	}
}
