package requests

import (
	"strings"

	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

// CreateExampleRequest represents the create example request payload
type CreateExampleRequest struct {
	Name   string `json:"name" binding:"required,min=3"`
	Email  string `json:"email" binding:"required,email"`
	Status string `json:"status" binding:"required"`
}

// Validate validates and binds the create example request
func (r *CreateExampleRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
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

// UpdateExampleRequest represents the update example request payload
type UpdateExampleRequest struct {
	Name   string `json:"name" binding:"omitempty,min=3"`
	Email  string `json:"email" binding:"omitempty,email"`
	Status string `json:"status" binding:"omitempty"`
}

// Validate validates and binds the update example request
func (r *UpdateExampleRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
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
