package requests

import (
	"strings"

	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Validate validates and binds the register request
func (r *RegisterRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		// Check if the error is related to JSON parsing, specifically numeric literals
		errStr := err.Error()
		if strings.Contains(errStr, "invalid character") && strings.Contains(errStr, "numeric literal") {
			utils.Error(c, 400, "json_parsing_error", "Invalid JSON in request body. Check for malformed numeric values.", gin.H{
				"details": errStr,
			})
			return err
		}
		utils.ValidationError(c, errStr)
		return err
	}
	return nil
}
