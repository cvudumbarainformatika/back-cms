package controllers

import (
	"net/http"
	"strconv"
	"time"

	requests "github.com/cvudumbarainformatika/backend/app/Http/Requests"
	models "github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// AgendaController handles agenda (event) operations
type AgendaController struct {
	db *sqlx.DB
}

// NewAgendaController creates a new AgendaController instance
func NewAgendaController(db *sqlx.DB) *AgendaController {
	return &AgendaController{
		db: db,
	}
}

// GetList returns paginated list of agenda with optional filters
// GET /api/v1/agenda?page=&limit=&type=&status=&upcoming=&sort=&order=
func (ac *AgendaController) GetList(c *gin.Context) {
	// Get pagination parameters
	page, limit := utils.GetPaginationParams(c)

	// Get filter parameters
	agendaType := c.Query("type")
	status := c.Query("status")
	upcomingStr := c.Query("upcoming")
	upcoming := upcomingStr == "true"

	// Get sort parameters
	// Default sort depends on 'upcoming'. Handled in model, but we pass filters.

	filters := map[string]interface{}{
		"type":     agendaType,
		"status":   status,
		"upcoming": upcoming,
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch agenda
	agendaList, total, err := models.GetAllAgenda(ac.db, filters, offset, limit)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch agenda: "+err.Error(), nil)
		return
	}

	// Format response
	agendaResponses := make([]gin.H, len(agendaList))
	for i, agenda := range agendaList {
		agendaResponses[i] = formatAgendaResponse(agenda)
	}

	// Use standard pagination response format
	pagination := utils.OffsetPaginate(agendaResponses, page, limit, total)

	utils.Success(c, http.StatusOK, "Agenda fetched successfully", gin.H{
		"items":      pagination.Data,
		"pagination": pagination.Meta,
	})
}

// GetBySlug returns a single agenda by slug
// GET /api/v1/agenda/:slug
func (ac *AgendaController) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	agenda, err := models.FindAgendaBySlug(ac.db, slug)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch agenda", nil)
		return
	}

	if agenda == nil {
		utils.Error(c, http.StatusNotFound, "agenda_not_found", "Agenda not found", nil)
		return
	}

	utils.Success(c, http.StatusOK, "Agenda retrieved successfully", formatAgendaResponse(*agenda))
}

// Create creates a new agenda
// POST /api/v1/agenda
func (ac *AgendaController) Create(c *gin.Context) {
	var req requests.CreateAgendaRequest

	if err := req.Validate(c); err != nil {
		return
	}

	// Generate slug from title
	slug := utils.GenerateSlug(req.Title)

	// Check if slug already exists
	existing, _ := models.FindAgendaBySlug(ac.db, slug)
	if existing != nil {
		slug = slug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	// Parse dates
	eventDate, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_date", "Invalid date format (RFC3339 required)", nil)
		return
	}

	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		parsedEnd, err := time.Parse(time.RFC3339, *req.EndDate)
		if err == nil {
			endDate = &parsedEnd
		}
	}

	// Create agenda model
	agenda := &models.Agenda{
		Slug:            slug,
		Title:           req.Title,
		Description:     req.Description,
		Type:            req.Type,
		Date:            eventDate,
		EndDate:         endDate,
		IsOnline:        req.IsOnline,
		Location:        req.Location,
		SKP:             req.Skp, // Note: Model field is SKP (capitalized in previous thought, check model file)
		Quota:           req.Quota,
		RegistrationURL: req.RegistrationURL,
		ImageURL:        req.ImageURL,
		Fee:             req.Fee,
		Status:          req.Status,
	}

	// Set published_at if status is published
	if req.Status == "published" {
		if req.PublishedAt != nil && *req.PublishedAt != "" {
			publishedAt, err := time.Parse(time.RFC3339, *req.PublishedAt)
			if err == nil {
				agenda.PublishedAt = &publishedAt
			}
		}
		if agenda.PublishedAt == nil {
			now := time.Now()
			agenda.PublishedAt = &now
		}
	}

	// Save to database
	err = agenda.Create(ac.db)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to create agenda: "+err.Error(), nil)
		return
	}

	utils.Success(c, http.StatusCreated, "Agenda created successfully", formatAgendaResponse(*agenda))
}

// Update updates an agenda
// PUT /api/v1/agenda/:id
func (ac *AgendaController) Update(c *gin.Context) {
	id := c.Param("id")

	agendaID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid agenda ID", nil)
		return
	}

	var req requests.UpdateAgendaRequest
	if err := req.Validate(c); err != nil {
		return
	}

	// Get existing agenda
	agenda, err := models.FindAgendaByID(ac.db, agendaID)
	if err != nil || agenda == nil {
		utils.Error(c, http.StatusNotFound, "agenda_not_found", "Agenda not found", nil)
		return
	}

	// Parse dates
	eventDate, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_date", "Invalid date format", nil)
		return
	}

	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		parsedEnd, err := time.Parse(time.RFC3339, *req.EndDate)
		if err == nil {
			endDate = &parsedEnd
		}
	} else {
		endDate = nil
	}

	// Update fields
	agenda.Title = req.Title
	agenda.Description = req.Description
	agenda.Type = req.Type
	agenda.Date = eventDate
	agenda.EndDate = endDate
	agenda.IsOnline = req.IsOnline
	agenda.Location = req.Location
	agenda.SKP = req.Skp
	agenda.Quota = req.Quota
	agenda.RegistrationURL = req.RegistrationURL
	agenda.ImageURL = req.ImageURL
	agenda.Fee = req.Fee
	agenda.Status = req.Status

	// Update published_at logic
	if req.Status == "published" {
		if req.PublishedAt != nil && *req.PublishedAt != "" {
			publishedAt, err := time.Parse(time.RFC3339, *req.PublishedAt)
			if err == nil {
				agenda.PublishedAt = &publishedAt
			}
		} else if agenda.PublishedAt == nil {
			now := time.Now()
			agenda.PublishedAt = &now
		}
	}

	// Save to database
	err = agenda.Update(ac.db)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to update agenda: "+err.Error(), nil)
		return
	}

	utils.Success(c, http.StatusOK, "Agenda updated successfully", formatAgendaResponse(*agenda))
}

// Patch performs partial update on an agenda (status, published_at, deleted_at)
// PATCH /api/v1/agenda/:id
func (ac *AgendaController) Patch(c *gin.Context) {
	id := c.Param("id")

	agendaID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid agenda ID", nil)
		return
	}

	var req requests.PatchAgendaRequest
	if err := req.Validate(c); err != nil {
		return
	}

	// Get existing agenda
	agenda, err := models.FindAgendaByID(ac.db, agendaID)
	if err != nil || agenda == nil {
		utils.Error(c, http.StatusNotFound, "agenda_not_found", "Agenda not found", nil)
		return
	}

	// Update status if provided
	if req.Status != "" {
		agenda.Status = req.Status
		if req.Status == "published" && agenda.PublishedAt == nil {
			now := time.Now()
			agenda.PublishedAt = &now
		}
	}

	// Update published_at if provided
	if req.PublishedAt != nil {
		if *req.PublishedAt == "" {
			agenda.PublishedAt = nil
		} else {
			publishedAt, err := time.Parse(time.RFC3339, *req.PublishedAt)
			if err == nil {
				agenda.PublishedAt = &publishedAt
			}
		}
	}

	// Handle soft delete/restore
	if req.DeletedAt != nil {
		if *req.DeletedAt == "" {
			// Restore: set deleted_at to NULL
			agenda.DeletedAt = nil
		} else {
			// Soft delete
			deletedAt, err := time.Parse(time.RFC3339, *req.DeletedAt)
			if err == nil {
				agenda.DeletedAt = &deletedAt
			} else {
				now := time.Now()
				agenda.DeletedAt = &now
			}
		}
	}

	// Save to database
	err = agenda.Update(ac.db)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to update agenda: "+err.Error(), nil)
		return
	}

	utils.Success(c, http.StatusOK, "Agenda updated successfully", formatAgendaResponse(*agenda))
}

// Delete performs soft delete on an agenda
// DELETE /api/v1/agenda/:id
func (ac *AgendaController) Delete(c *gin.Context) {
	id := c.Param("id")

	agendaID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid agenda ID", nil)
		return
	}

	// Get existing agenda
	agenda, err := models.FindAgendaByID(ac.db, agendaID)
	if err != nil || agenda == nil {
		utils.Error(c, http.StatusNotFound, "agenda_not_found", "Agenda not found", nil)
		return
	}

	// Soft delete
	err = agenda.Delete(ac.db)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to delete agenda: "+err.Error(), nil)
		return
	}

	utils.Success(c, http.StatusOK, "Agenda deleted successfully", nil)
}

// GetTypes returns all unique agenda types
// GET /api/v1/agenda/types
func (ac *AgendaController) GetTypes(c *gin.Context) {
	types, err := models.GetAgendaTypes(ac.db)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch types", nil)
		return
	}

	utils.Success(c, http.StatusOK, "Agenda types fetched successfully", gin.H{
		"types": types,
	})
}

// Helper function to format agenda response
func formatAgendaResponse(agenda models.Agenda) gin.H {
	response := gin.H{
		"id":               agenda.ID,
		"slug":             agenda.Slug,
		"title":            agenda.Title,
		"description":      agenda.Description,
		"type":             agenda.Type,
		"date":             agenda.Date,
		"end_date":         agenda.EndDate,
		"is_online":        agenda.IsOnline,
		"location":         agenda.Location,
		"skp":              agenda.SKP,
		"quota":            agenda.Quota,
		"registration_url": agenda.RegistrationURL,
		"image_url":        agenda.ImageURL,
		"fee":              agenda.Fee,
		"status":           agenda.Status,
		"created_at":       agenda.CreatedAt,
		"updated_at":       agenda.UpdatedAt,
	}

	if agenda.PublishedAt != nil {
		response["published_at"] = agenda.PublishedAt
	}

	if agenda.DeletedAt != nil {
		response["deleted_at"] = agenda.DeletedAt
	}

	return response
}
