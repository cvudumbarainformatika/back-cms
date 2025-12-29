package requests

import (
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

// UpdateProfileRequest represents the request payload for updating user profile
type UpdateProfileRequest struct {
	Name    string                `json:"name" form:"name" binding:"required,min=1,max=255"`
	Phone   string                `json:"phone" form:"phone" binding:"max=20"`
	Address string                `json:"address" form:"address" binding:"max=500"`
	Bio     string                `json:"bio" form:"bio" binding:"max=1000"`
	Avatar  *multipart.FileHeader `form:"avatar"`
}

// Validate validates the UpdateProfileRequest
// Supports both JSON and multipart/form-data
func (r *UpdateProfileRequest) Validate(c *gin.Context) error {
	contentType := c.ContentType()

	// Handle JSON request
	if contentType == "application/json" {
		if err := c.ShouldBindJSON(r); err != nil {
			utils.ValidationError(c, err.Error())
			return err
		}
	} else {
		// Handle form data (multipart or urlencoded)
		if err := c.ShouldBind(r); err != nil {
			utils.ValidationError(c, err.Error())
			return err
		}
	}
	return nil
}
