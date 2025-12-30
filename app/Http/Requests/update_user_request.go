package requests

import (
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

// UpdateUserRequest represents the request payload for updating a user (PUT)
type UpdateUserRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"max=255"` // Optional, can be empty to skip password update
	Role     string `json:"role" binding:"required,oneof=member admin_cabang admin_wilayah admin_pusat"`
	Status   string `json:"status" binding:"required,oneof=active pending inactive deleted"`
	Phone    string `json:"phone" binding:"max=20"`
	Address  string `json:"address" binding:"max=500"`
	Bio      string `json:"bio" binding:"max=1000"`
	Cabang   string `json:"cabang" binding:"max=255"`
}

// Validate validates the UpdateUserRequest
func (r *UpdateUserRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		utils.ValidationError(c, err.Error())
		return err
	}

	return nil
}
