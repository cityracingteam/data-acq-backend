package routes

import (
	"context"

	"github.com/cityracingteam/data-acq-backend/models"
	"github.com/cityracingteam/data-acq-backend/util/jwt"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func AuthCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Authenticated sucessfully!

		// Get a models.User object
		user := models.GetUserFromGoth(&gothUser)
		// Use said object to issue a short-lived access token
		token, err := jwt.NewAccessJwt(*user)

		// All good, return a success and the access token
		if err == nil {
			c.JSON(200, gin.H{
				"success":      true,
				"access_token": token,
			})
		} else {
			// Error issuing token
			c.JSON(500, gin.H{
				"error": "error issuing jwt",
			})
		}
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
