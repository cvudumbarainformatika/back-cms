package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

// AvatarAuthMiddleware validates avatar access permissions
// Allows:
// 1. Public access to own avatar
// 2. Admin access to any avatar
// 3. Configurable public/private access via environment
func AvatarAuthMiddleware(isPublic bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If avatar is configured as public, allow all access
		if isPublic {
			c.Next()
			return
		}

		// Extract user_id from URL parameter (for profile endpoint)
		userIDStr := c.Param("user_id")
		if userIDStr == "" {
			utils.Error(c, http.StatusBadRequest, "invalid_request", "User ID is required", nil)
			c.Abort()
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "invalid_user_id", "Invalid user ID format", nil)
			c.Abort()
			return
		}

		// Get authenticated user ID from JWT context
		authUserID, exists := c.Get("user_id")
		if !exists {
			// Allow unauthenticated access if avatar is semi-public
			// This can be configured based on your requirements
			utils.Error(c, http.StatusUnauthorized, "unauthorized", "Authentication required to access avatars", nil)
			c.Abort()
			return
		}

		authUserIDInt64, ok := authUserID.(int64)
		if !ok {
			utils.Error(c, http.StatusInternalServerError, "internal_error", "Invalid user ID format in token", nil)
			c.Abort()
			return
		}

		// Get user role from JWT context
		userRole, _ := c.Get("user_role")

		// Allow if:
		// 1. User accessing own avatar
		// 2. User is admin
		if authUserIDInt64 != userID && !isAdmin(userRole) {
			utils.Error(c, http.StatusForbidden, "forbidden", "You don't have permission to access this avatar", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// isAdmin checks if user role is admin
func isAdmin(role interface{}) bool {
	if role == nil {
		return false
	}

	roleStr, ok := role.(string)
	if !ok {
		return false
	}

	// Check for admin variations
	return strings.Contains(strings.ToLower(roleStr), "admin")
}
