package requests

import (
	"io"
	"strings"

	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required"` // Changed from "required,email" to just "required" to allow username
	Password string `json:"password" binding:"required"`
}

// Validate validates and binds the login request
func (r *LoginRequest) Validate(c *gin.Context) error {
	// Read the raw body to help with debugging
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.Error(c, 400, "request_error", "Failed to read request body", gin.H{
			"details": err.Error(),
		})
		return err
	}

	// Put the body back for further processing
	c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	if err := c.ShouldBindJSON(r); err != nil {
		// Check if the error is related to JSON parsing, specifically numeric literals
		errStr := err.Error()
		if strings.Contains(errStr, "invalid character") && strings.Contains(errStr, "numeric literal") {
			utils.Error(c, 400, "json_parsing_error", "Invalid JSON in request body. Check for malformed numeric values.", gin.H{
				"raw_body": string(bodyBytes),
				"details":  errStr,
			})
			return err
		}
		// Check for other JSON parsing errors and treat them as validation errors
		if strings.Contains(errStr, "invalid character") {
			utils.Error(c, 400, "json_parsing_error", "Invalid JSON in request body. Check for malformed JSON.", gin.H{
				"raw_body": string(bodyBytes),
				"details":  errStr,
			})
			return err
		}
		utils.ValidationError(c, errStr)
		return err
	}
	return nil
}
