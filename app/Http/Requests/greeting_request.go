package requests

import (
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

// CreateGreetingRequest represents the request payload for creating a new greeting
type CreateGreetingRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=255"`
	Content  string `json:"content" binding:"required,min=1"`
	ImageURL string `json:"image_url" binding:"omitempty,max=255"`
	IsActive bool   `json:"is_active" binding:"omitempty"`
}

// Validate validates the CreateGreetingRequest
func (r *CreateGreetingRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		utils.ValidationError(c, err.Error())
		return err
	}
	return nil
}

// UpdateGreetingRequest represents the request payload for updating a greeting
type UpdateGreetingRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=255"`
	Content  string `json:"content" binding:"required,min=1"`
	ImageURL string `json:"image_url" binding:"omitempty,max=255"`
	IsActive bool   `json:"is_active" binding:"omitempty"`
}

// Validate validates the UpdateGreetingRequest
func (r *UpdateGreetingRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		utils.ValidationError(c, err.Error())
		return err
	}
	return nil
}
