package requests

import (
	"net/http"

	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

type CreateDocumentRequest struct {
	Name       string `form:"name" binding:"required"`
	Type       string `form:"type" binding:"required"`
	ValidUntil string `form:"valid_until"` // Formatted as YYYY-MM-DD
}

func (r *CreateDocumentRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBind(r); err != nil {
		utils.ValidationError(c, err.Error())
		return err
	}

	// Validate file upload separately in the controller
	_, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "file_required", "File dokumen wajib diunggah", nil)
		return err
	}

	return nil
}
