package requests

import (
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

// CreateBeritaRequest represents the request payload for creating a new berita
type CreateBeritaRequest struct {
	Title       string   `json:"title" binding:"required,min=1,max=255"`
	Excerpt     string   `json:"excerpt" binding:"required,min=1"`
	Content     string   `json:"content" binding:"required,min=1"`
	ImageURL    string   `json:"image_url" binding:"omitempty,max=255"`
	Category    string   `json:"category" binding:"required,oneof=umum ilmiah kegiatan pengumuman prestasi"`
	Author      string   `json:"author" binding:"required,min=1,max=255"`
	Status      string   `json:"status" binding:"omitempty,oneof=draft published"`
	Tags        []string `json:"tags" binding:"omitempty"`
	PublishedAt *string  `json:"published_at" binding:"omitempty"`
}

// Validate validates the CreateBeritaRequest
func (r *CreateBeritaRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		utils.ValidationError(c, err.Error())
		return err
	}

	// Set default status if not provided
	if r.Status == "" {
		r.Status = "draft"
	}

	return nil
}

// UpdateBeritaRequest represents the request payload for updating a berita
type UpdateBeritaRequest struct {
	Title       string   `json:"title" binding:"required,min=1,max=255"`
	Excerpt     string   `json:"excerpt" binding:"required,min=1"`
	Content     string   `json:"content" binding:"required,min=1"`
	ImageURL    string   `json:"image_url" binding:"omitempty,max=255"`
	Category    string   `json:"category" binding:"required,oneof=umum ilmiah kegiatan pengumuman prestasi"`
	Author      string   `json:"author" binding:"required,min=1,max=255"`
	Status      string   `json:"status" binding:"required,oneof=draft published"`
	Tags        []string `json:"tags" binding:"omitempty"`
	PublishedAt *string  `json:"published_at" binding:"omitempty"`
}

// Validate validates the UpdateBeritaRequest
func (r *UpdateBeritaRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		utils.ValidationError(c, err.Error())
		return err
	}

	return nil
}

// PatchBeritaRequest represents the request payload for partial update (status, published_at, deleted_at)
type PatchBeritaRequest struct {
	Status      string  `json:"status" binding:"omitempty,oneof=draft published"`
	PublishedAt *string `json:"published_at" binding:"omitempty"`
	DeletedAt   *string `json:"deleted_at" binding:"omitempty"`
}

// Validate validates the PatchBeritaRequest
func (r *PatchBeritaRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		utils.ValidationError(c, err.Error())
		return err
	}

	return nil
}
