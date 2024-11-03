package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/cityracingteam/data-acq-backend/models"
	"github.com/cityracingteam/data-acq-backend/util/jwt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	// ints without a value default to 0; start at 1

	RESULT_USER_FETCHED_SUCESSFULLY      int = 1
	RESULT_NO_AUTH_HEADER                int = 2
	RESULT_INVALID_AUTH_HEADER           int = 3
	RESULT_NO_TOKEN_IN_AUTH_HEADER       int = 4
	RESULT_INVALID_JWT                   int = 5
	RESULT_VALID_JWT_USER_DOES_NOT_EXIST int = 6
)

// Middleware to decode the jwt and place the associated user data into the context
func RequireAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var result int
		var user *models.User

		// Check if there is a jwt present in the Authorization header
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" {
			// Header is present.
			splitHeader := strings.Split(authHeader, "Bearer ")
			if len(splitHeader) == 2 {
				if splitHeader[1] != "" {
					tokenStr := splitHeader[1]

					token, _, err := jwt.VerifyJwt(tokenStr)
					if err == nil {
						// JWT was verified successfully.
						// Let's fetch the user for this token.
						userRes, err2 := jwt.GetUserForJwt(db, token)
						if err2 == nil {
							// User was fetched successfully
							user = userRes
							result = RESULT_USER_FETCHED_SUCESSFULLY

						} else {
							// JWT was valid but user did not exist.
							result = RESULT_VALID_JWT_USER_DOES_NOT_EXIST
							log.Fatal("Valid JWT presented but user could not be fetched from database!")
						}
					} else {
						// JWT was not valid
						result = RESULT_INVALID_JWT
					}
				} else {
					// Authorization header contains "Bearer " but no token, ignore.
					result = RESULT_NO_TOKEN_IN_AUTH_HEADER
				}
			} else {
				// Auth header is invalid (doesn't contain "Bearer ")
				result = RESULT_INVALID_AUTH_HEADER
			}
		} else {
			// Authorization header is not present, ignore.
			result = RESULT_NO_AUTH_HEADER

		}

		// If user is present, set context
		if user != nil {
			// Get a copy of the parent context with the user key set
			ctx := context.WithValue(c.Request.Context(), "user", user)
			// Create a new request with the new context
			newRequest := c.Request.WithContext(ctx)
			// Overwrite the original request with the new one
			c.Request = newRequest
		}

		if result != 0 {
			// Add result to context
			// Get a copy of the parent context with the result key set
			ctx := context.WithValue(c.Request.Context(), "result", result)
			// Create a new request with the new context
			newRequest := c.Request.WithContext(ctx)
			// Overwrite the original request with the new one
			c.Request = newRequest
		}

		// Call the next handler (including the target itself)
		c.Next()

		// Request has been executed
	}
}
