package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// Success sends a successful JSON response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error JSON response
func Error(c *gin.Context, statusCode int, err string, message string, details interface{}) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   err,
		Message: message,
		Details: details,
	})
}

// ValidationError sends a validation error response (422)
func ValidationError(c *gin.Context, errors interface{}) {
	c.JSON(422, ErrorResponse{
		Success: false,
		Error:   "validation_failed",
		Message: "Validation failed",
		Details: errors,
	})
}

// SetRedisCache stores data in Redis with JSON serialization
func SetRedisCache(redisClient *redis.Client, ctx context.Context, key string, data interface{}, expiration int) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return redisClient.Set(ctx, key, jsonData, 0).Err()
}

// GetCurrentTimeString returns the current time in "YYYY-MM-DD HH:MM:SS" format for MySQL datetime
func GetCurrentTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GetCurrentTimeForMySQL returns the current time in MySQL datetime format
func GetCurrentTimeForMySQL() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GetCurrentUnixTimestamp returns the current Unix timestamp as float64
func GetCurrentUnixTimestamp() float64 {
	return float64(time.Now().Unix())
}
