package middleware

import (
	"net/http"
	"strings"

	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware validates JWT tokens from the Authorization header
func JWTAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "unauthorized",
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "unauthorized",
				"message": "Invalid authorization header format. Expected: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			// Check if the error is related to JSON parsing
			errStr := err.Error()
			if strings.Contains(errStr, "invalid character") && strings.Contains(errStr, "numeric literal") {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   "json_parsing_error",
					"message": "Invalid token format: JSON parsing error",
					"details": errStr,
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "unauthorized",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user data in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}
