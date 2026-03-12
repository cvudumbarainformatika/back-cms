package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

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

// NormalizePhoneNumber cleans and formats Indonesian phone numbers for WhatsApp
func NormalizePhoneNumber(phone string) string {
	// 1. Remove non-numeric characters
	reg := regexp.MustCompile("[^0-9]")
	cleaned := reg.ReplaceAllString(phone, "")

	if cleaned == "" {
		return ""
	}

	// 2. Handle leading '0' (e.g., 0812...) -> replace with '62'
	if strings.HasPrefix(cleaned, "0") {
		return "62" + cleaned[1:]
	}

	// 3. Handle leading '8' (e.g., 8122...) -> prepend '62'
	if strings.HasPrefix(cleaned, "8") {
		return "62" + cleaned
	}

	// 4. If it starts with '62', keep it (e.g., 62812...)
	if strings.HasPrefix(cleaned, "62") {
		return cleaned
	}

	return cleaned
}
