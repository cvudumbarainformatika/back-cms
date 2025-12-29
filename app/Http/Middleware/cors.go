package middleware

import (
	"github.com/cvudumbarainformatika/backend/config"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(cfg config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers
		if len(cfg.AllowedOrigins) > 0 {
			origin := c.Request.Header.Get("Origin")
			for _, allowedOrigin := range cfg.AllowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
					break
				}
			}
		}

		if len(cfg.AllowedMethods) > 0 {
			methods := ""
			for i, method := range cfg.AllowedMethods {
				if i > 0 {
					methods += ", "
				}
				methods += method
			}
			c.Writer.Header().Set("Access-Control-Allow-Methods", methods)
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
