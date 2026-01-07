package requests

import (
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

// CreateAgendaRequest represents the request payload for creating a new agenda
type CreateAgendaRequest struct {
	Title           string  `json:"title" binding:"required,min=1,max=255"`
	Description     string  `json:"description" binding:"required,min=1"`
	Type            string  `json:"type" binding:"required,oneof=webinar workshop seminar kongres pelatihan"`
	Date            string  `json:"date" binding:"required"` // ISO8601 string
	EndDate         *string `json:"end_date" binding:"omitempty"`
	IsOnline        bool    `json:"is_online" binding:"omitempty"`
	Location        string  `json:"location" binding:"required"`
	Skp             float64 `json:"skp" binding:"omitempty"`
	Quota           int     `json:"quota" binding:"omitempty"`
	RegistrationURL string  `json:"registration_url" binding:"omitempty"`
	ImageURL        string  `json:"image_url" binding:"omitempty"`
	Fee             string  `json:"fee" binding:"omitempty"`
	Status          string  `json:"status" binding:"omitempty,oneof=draft published"`
	PublishedAt     *string `json:"published_at" binding:"omitempty"`
}

// Validate validates the CreateAgendaRequest
func (r *CreateAgendaRequest) Validate(c *gin.Context) error {
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

// UpdateAgendaRequest represents the request payload for updating an agenda
type UpdateAgendaRequest struct {
	Title           string  `json:"title" binding:"required,min=1,max=255"`
	Description     string  `json:"description" binding:"required,min=1"`
	Type            string  `json:"type" binding:"required,oneof=webinar workshop seminar kongres pelatihan"`
	Date            string  `json:"date" binding:"required"`
	EndDate         *string `json:"end_date" binding:"omitempty"`
	IsOnline        bool    `json:"is_online" binding:"omitempty"`
	Location        string  `json:"location" binding:"required"`
	Skp             float64 `json:"skp" binding:"omitempty"`
	Quota           int     `json:"quota" binding:"omitempty"`
	RegistrationURL string  `json:"registration_url" binding:"omitempty"`
	ImageURL        string  `json:"image_url" binding:"omitempty"`
	Fee             string  `json:"fee" binding:"omitempty"`
	Status          string  `json:"status" binding:"required,oneof=draft published"`
	PublishedAt     *string `json:"published_at" binding:"omitempty"`
}

// Validate validates the UpdateAgendaRequest
func (r *UpdateAgendaRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		utils.ValidationError(c, err.Error())
		return err
	}

	return nil
}

// PatchAgendaRequest represents the request payload for partial update (status, published_at, deleted_at)
type PatchAgendaRequest struct {
	Status      string  `json:"status" binding:"omitempty,oneof=draft published"`
	PublishedAt *string `json:"published_at" binding:"omitempty"`
	DeletedAt   *string `json:"deleted_at" binding:"omitempty"`
}

// Validate validates the PatchAgendaRequest
func (r *PatchAgendaRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		utils.ValidationError(c, err.Error())
		return err
	}

	return nil
}
