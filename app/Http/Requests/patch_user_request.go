package requests

import (
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

// PatchUserRequest represents the request payload for patching a user (PATCH)
// Used to update only role and/or status
type PatchUserRequest struct {
	Role   string `json:"role" binding:"omitempty,oneof=member admin_cabang admin_wilayah admin_pusat"`
	Status string `json:"status" binding:"omitempty,oneof=active pending inactive deleted"`
}

// Validate validates the PatchUserRequest
func (r *PatchUserRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		utils.ValidationError(c, err.Error())
		return err
	}

	// At least one field should be provided
	if r.Role == "" && r.Status == "" {
		utils.ValidationError(c, "At least one field (role or status) must be provided")
		return nil
	}

	return nil
}
