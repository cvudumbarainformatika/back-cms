package exceptions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a middleware that handles errors and panics
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Handle panic
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "internal_server_error",
					"message": "An unexpected error occurred",
				})
				c.Abort()
			}
		}()

		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Check if it's an AppError
			if appErr, ok := err.Err.(*AppError); ok {
				c.JSON(appErr.Code, gin.H{
					"success": false,
					"error":   http.StatusText(appErr.Code),
					"message": appErr.Message,
				})
				return
			}

			// Default error response
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "internal_server_error",
				"message": err.Error(),
			})
		}
	}
}
