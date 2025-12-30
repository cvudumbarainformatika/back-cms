package requests

import (
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

// CreateUserRequest represents the request payload for creating a new user
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=member admin_cabang admin_wilayah admin_pusat"`
	Status   string `json:"status" binding:"required,oneof=active pending inactive"`
	Phone    string `json:"phone" binding:"max=20"`
	Address  string `json:"address" binding:"max=500"`
	Bio      string `json:"bio" binding:"max=1000"`
	Cabang   string `json:"cabang" binding:"max=255"`
}

// Validate validates the CreateUserRequest
func (r *CreateUserRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		utils.ValidationError(c, err.Error())
		return err
	}

	return nil
}
