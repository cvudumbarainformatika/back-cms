package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole checks if the authenticated user has one of the required roles
// Note: This is a placeholder implementation. In a real application, you would
// fetch the user's roles from the database and check against the required roles.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by JWTAuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		// TODO: Fetch user roles from database using userID
		// For now, this is a placeholder that always allows access
		_ = userID

		// TODO: Check if user has one of the required roles
		// hasRole := false
		// for _, role := range roles {
		//     if userHasRole(userID, role) {
		//         hasRole = true
		//         break
		//     }
		// }
		//
		// if !hasRole {
		//     c.JSON(http.StatusForbidden, gin.H{
		//         "success": false,
		//         "error":   "forbidden",
		//         "message": "Insufficient permissions",
		//     })
		//     c.Abort()
		//     return
		// }

		c.Next()
	}
}

// RequirePermission checks if the authenticated user has one of the required permissions
// Note: This is a placeholder implementation. In a real application, you would
// fetch the user's permissions from the database and check against the required permissions.
func RequirePermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by JWTAuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		// TODO: Fetch user permissions from database using userID
		// For now, this is a placeholder that always allows access
		_ = userID

		// TODO: Check if user has one of the required permissions
		// hasPermission := false
		// for _, permission := range permissions {
		//     if userHasPermission(userID, permission) {
		//         hasPermission = true
		//         break
		//     }
		// }
		//
		// if !hasPermission {
		//     c.JSON(http.StatusForbidden, gin.H{
		//         "success": false,
		//         "error":   "forbidden",
		//         "message": "Insufficient permissions",
		//     })
		//     c.Abort()
		//     return
		// }

		c.Next()
	}
}
