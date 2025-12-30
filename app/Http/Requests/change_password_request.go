package requests

import (
	"errors"
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

// ChangePasswordRequest represents the request payload for changing password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required,min=6"`
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
}

// Validate validates the ChangePasswordRequest
func (r *ChangePasswordRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		utils.ValidationError(c, "Format request tidak valid")
		return err
	}

	// Validate that new password is different from current
	if r.CurrentPassword == r.NewPassword {
		utils.ValidationError(c, "Password baru harus berbeda dengan password saat ini")
		return errors.New("password baru harus berbeda dengan password saat ini")
	}

	return nil
}
