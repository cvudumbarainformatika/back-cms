package requests

import (
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

type UpdateHomepageRequest struct {
	Hero  map[string]interface{}   `json:"hero" binding:"required"`
	Stats []map[string]interface{} `json:"stats" binding:"required"`
	SEO   map[string]interface{}   `json:"seo" binding:"required"`
}

func (r *UpdateHomepageRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		utils.Error(c, 400, "validation_error", err.Error(), nil)
		return err
	}
	return nil
}
