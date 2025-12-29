package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/cvudumbarainformatika/backend/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Visitor represents a rate limiter for a specific IP
type Visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter manages rate limiting for multiple visitors
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	config   config.RateLimitConfig
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(cfg config.RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		config:   cfg,
	}

	// Start cleanup goroutine
	go rl.cleanupVisitors()

	return rl
}

// getVisitor returns the rate limiter for a specific IP
func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	visitor, exists := rl.visitors[ip]
	if !exists {
		// Create new limiter for this IP
		// rate.Every calculates the interval between events
		limiter := rate.NewLimiter(
			rate.Every(time.Duration(rl.config.WindowSeconds)*time.Second/time.Duration(rl.config.MaxRequests)),
			rl.config.MaxRequests,
		)
		rl.visitors[ip] = &Visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	// Update last seen time
	visitor.lastSeen = time.Now()
	return visitor.limiter
}

// cleanupVisitors removes old visitors to prevent memory leaks
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		rl.mu.Lock()
		for ip, visitor := range rl.visitors {
			if time.Since(visitor.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Middleware returns the rate limiting middleware
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.config.Enabled {
			c.Next()
			return
		}

		// Get client IP
		ip := c.ClientIP()

		// Get rate limiter for this IP
		limiter := rl.getVisitor(ip)

		// Check if request is allowed
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "rate_limit_exceeded",
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
