package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func AuthCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Authenticated sucessfully, placeholder response
		c.JSON(200, gin.H{
			"user": user,
		})
	}
}

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add the provider specified in /auth/:provider to context so gothic can read it
		// Based of <https://github.com/markbates/goth/issues/264#issuecomment-997328048>
		provider := c.Param("provider")
		c.Request = c.Request.WithContext(context.WithValue(context.Background(), "provider", provider))

		// Process the (newly updated) request
		if user, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
			// Authenticated sucessfully, placeholder response
			c.JSON(200, gin.H{
				"user": user,
			})
		} else {
			// Not authenticated
			gothic.BeginAuthHandler(c.Writer, c.Request)
		}
	}
}
