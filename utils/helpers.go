package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// QueryInt extracts an integer query parameter with a default value
func QueryInt(c *gin.Context, key string, defaultValue int) int {
	if valueStr := c.Query(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// StringToInt64 converts  string to int64
func StringToInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to int64: %w", err)
	}
	return i, nil
}

// Int64ToString converts int64 to string
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
