package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	requests "github.com/cvudumbarainformatika/backend/app/Http/Requests"
	models "github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// BeritaController handles berita (news) operations
type BeritaController struct {
	db *sqlx.DB
}

// NewBeritaController creates a new BeritaController instance
func NewBeritaController(db *sqlx.DB) *BeritaController {
	return &BeritaController{
		db: db,
	}
}

// GetList returns paginated list of berita with optional filters
// GET /api/v1/berita?page=&limit=&category=&author=&status=&search=&sort=&order=
func (bc *BeritaController) GetList(c *gin.Context) {
	// Get pagination parameters
	page, limit := utils.GetPaginationParams(c)

	// Get filter parameters
	category := c.Query("category")
	author := c.Query("author")
	status := c.Query("status")
	search := c.Query("search")

	// Get sort parameters
	orderBy := c.DefaultQuery("sort", "created_at")
	sortOrder := c.DefaultQuery("order", "desc")

	// Validate orderBy to prevent SQL injection
	allowedColumns := map[string]bool{
		"id":           true,
		"title":        true,
		"category":     true,
		"author":       true,
		"status":       true,
		"views":        true,
		"published_at": true,
		"created_at":   true,
		"updated_at":   true,
	}
	if !allowedColumns[orderBy] {
		orderBy = "created_at"
	}

	// Validate and normalize sort order
	sortOrder = strings.ToUpper(sortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	// Build query
	query := `SELECT id, slug, title, excerpt, content, image_url, category, author, status, views, published_at, created_at, updated_at, deleted_at 
	          FROM berita WHERE deleted_at IS NULL`
	args := []interface{}{}

	// Add filters
	if category != "" {
		query += ` AND category = ?`
		args = append(args, category)
	}

	if author != "" {
		query += ` AND author = ?`
		args = append(args, author)
	}

	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}

	if search != "" {
		query += ` AND (title LIKE ? OR excerpt LIKE ? OR content LIKE ?)`
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}

	// Count total records
	countQuery := `SELECT COUNT(*) FROM berita WHERE deleted_at IS NULL`
	countArgs := []interface{}{}

	if category != "" {
		countQuery += ` AND category = ?`
		countArgs = append(countArgs, category)
	}
	if author != "" {
		countQuery += ` AND author = ?`
		countArgs = append(countArgs, author)
	}
	if status != "" {
		countQuery += ` AND status = ?`
		countArgs = append(countArgs, status)
	}
	if search != "" {
		countQuery += ` AND (title LIKE ? OR excerpt LIKE ? OR content LIKE ?)`
		searchPattern := "%" + search + "%"
		countArgs = append(countArgs, searchPattern, searchPattern, searchPattern)
	}

	var total int64
	err := bc.db.Get(&total, countQuery, countArgs...)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to count berita", nil)
		return
	}

	// Add ordering and pagination
	query += ` ORDER BY ` + orderBy + ` ` + sortOrder + ` LIMIT ? OFFSET ?`
	offset := (page - 1) * limit
	args = append(args, limit, offset)

	// Fetch berita
	var beritaList []models.Berita
	err = bc.db.Select(&beritaList, query, args...)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch berita: "+err.Error(), nil)
		return
	}

	// Format response
	beritaResponses := make([]gin.H, len(beritaList))
	for i, berita := range beritaList {
		beritaResponses[i] = formatBeritaResponse(berita, false) // false = tidak include content
	}

	// Use standard pagination response format
	pagination := utils.OffsetPaginate(beritaResponses, page, limit, total)

	utils.Success(c, http.StatusOK, "Berita fetched successfully", gin.H{
		"items":      pagination.Data,
		"pagination": pagination.Meta,
	})
}

// GetBySlug returns a single berita by slug
// GET /api/v1/berita/:slug
func (bc *BeritaController) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	berita, err := models.FindBeritaBySlug(bc.db, slug)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch berita", nil)
		return
	}

	if berita == nil {
		utils.Error(c, http.StatusNotFound, "berita_not_found", "Berita not found", nil)
		return
	}

	// Increment views
	_, _ = bc.db.Exec(`UPDATE berita SET views = views + 1 WHERE id = ?`, berita.ID)

	utils.Success(c, http.StatusOK, "Berita retrieved successfully", formatBeritaResponse(*berita, true))
}

// Create creates a new berita
// POST /api/v1/berita
func (bc *BeritaController) Create(c *gin.Context) {
	var req requests.CreateBeritaRequest

	if err := req.Validate(c); err != nil {
		return
	}

	// Generate slug from title
	slug := utils.GenerateSlug(req.Title)

	// Check if slug already exists
	existing, _ := models.FindBeritaBySlug(bc.db, slug)
	if existing != nil {
		// Add timestamp to make it unique
		slug = slug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	// Create berita model
	berita := &models.Berita{
		Slug:     slug,
		Title:    req.Title,
		Excerpt:  req.Excerpt,
		Content:  req.Content,
		ImageURL: req.ImageURL,
		Category: req.Category,
		Author:   req.Author,
		Status:   req.Status,
		Views:    0,
		Tags:     req.Tags,
	}

	// Set published_at if status is published
	if req.Status == "published" {
		if req.PublishedAt != nil && *req.PublishedAt != "" {
			publishedAt, err := time.Parse(time.RFC3339, *req.PublishedAt)
			if err == nil {
				berita.PublishedAt = &publishedAt
			}
		}
		if berita.PublishedAt == nil {
			now := time.Now()
			berita.PublishedAt = &now
		}
	}

	// Save to database
	err := berita.Create(bc.db)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to create berita: "+err.Error(), nil)
		return
	}

	utils.Success(c, http.StatusCreated, "Berita created successfully", formatBeritaResponse(*berita, true))
}

// Update updates a berita
// PUT /api/v1/berita/:id
func (bc *BeritaController) Update(c *gin.Context) {
	id := c.Param("id")

	beritaID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid berita ID", nil)
		return
	}

	var req requests.UpdateBeritaRequest
	if err := req.Validate(c); err != nil {
		return
	}

	// Get existing berita
	berita, err := models.FindBeritaByID(bc.db, beritaID)
	if err != nil || berita == nil {
		utils.Error(c, http.StatusNotFound, "berita_not_found", "Berita not found", nil)
		return
	}

	// Update fields
	berita.Title = req.Title
	berita.Excerpt = req.Excerpt
	berita.Content = req.Content
	berita.ImageURL = req.ImageURL
	berita.Category = req.Category
	berita.Author = req.Author
	berita.Status = req.Status
	berita.Tags = req.Tags

	// Update published_at if status changed to published
	if req.Status == "published" {
		if req.PublishedAt != nil && *req.PublishedAt != "" {
			publishedAt, err := time.Parse(time.RFC3339, *req.PublishedAt)
			if err == nil {
				berita.PublishedAt = &publishedAt
			}
		} else if berita.PublishedAt == nil {
			now := time.Now()
			berita.PublishedAt = &now
		}
	}

	// Save to database
	err = berita.Update(bc.db)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to update berita: "+err.Error(), nil)
		return
	}

	utils.Success(c, http.StatusOK, "Berita updated successfully", formatBeritaResponse(*berita, true))
}

// Patch performs partial update on a berita (status, published_at, deleted_at)
// PATCH /api/v1/berita/:id
func (bc *BeritaController) Patch(c *gin.Context) {
	id := c.Param("id")

	beritaID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid berita ID", nil)
		return
	}

	var req requests.PatchBeritaRequest
	if err := req.Validate(c); err != nil {
		return
	}

	// Get existing berita
	berita, err := models.FindBeritaByID(bc.db, beritaID)
	if err != nil || berita == nil {
		utils.Error(c, http.StatusNotFound, "berita_not_found", "Berita not found", nil)
		return
	}

	// Update status if provided
	if req.Status != "" {
		berita.Status = req.Status
		if req.Status == "published" && berita.PublishedAt == nil {
			now := time.Now()
			berita.PublishedAt = &now
		}
	}

	// Update published_at if provided
	if req.PublishedAt != nil {
		if *req.PublishedAt == "" {
			berita.PublishedAt = nil
		} else {
			publishedAt, err := time.Parse(time.RFC3339, *req.PublishedAt)
			if err == nil {
				berita.PublishedAt = &publishedAt
			}
		}
	}

	// Handle soft delete/restore
	if req.DeletedAt != nil {
		if *req.DeletedAt == "" {
			// Restore: set deleted_at to NULL
			berita.DeletedAt = nil
		} else {
			// Soft delete
			deletedAt, err := time.Parse(time.RFC3339, *req.DeletedAt)
			if err == nil {
				berita.DeletedAt = &deletedAt
			} else {
				now := time.Now()
				berita.DeletedAt = &now
			}
		}
	}

	// Save to database
	err = berita.Update(bc.db)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to update berita: "+err.Error(), nil)
		return
	}

	utils.Success(c, http.StatusOK, "Berita updated successfully", formatBeritaResponse(*berita, true))
}

// Delete performs soft delete on a berita
// DELETE /api/v1/berita/:id
func (bc *BeritaController) Delete(c *gin.Context) {
	id := c.Param("id")

	beritaID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid berita ID", nil)
		return
	}

	// Get existing berita
	berita, err := models.FindBeritaByID(bc.db, beritaID)
	if err != nil || berita == nil {
		utils.Error(c, http.StatusNotFound, "berita_not_found", "Berita not found", nil)
		return
	}

	// Soft delete
	err = berita.Delete(bc.db)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to delete berita: "+err.Error(), nil)
		return
	}

	utils.Success(c, http.StatusOK, "Berita deleted successfully", nil)
}

// GetCategories returns all unique categories
// GET /api/v1/berita/categories
func (bc *BeritaController) GetCategories(c *gin.Context) {
	categories, err := models.GetBeritaCategories(bc.db)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch categories", nil)
		return
	}

	utils.Success(c, http.StatusOK, "Categories fetched successfully", gin.H{
		"categories": categories,
	})
}

// Helper function to format berita response
func formatBeritaResponse(berita models.Berita, includeContent bool) gin.H {
	response := gin.H{
		"id":         berita.ID,
		"slug":       berita.Slug,
		"title":      berita.Title,
		"excerpt":    berita.Excerpt,
		"image_url":  berita.ImageURL,
		"category":   berita.Category,
		"author":     berita.Author,
		"status":     berita.Status,
		"views":      berita.Views,
		"created_at": berita.CreatedAt,
		"updated_at": berita.UpdatedAt,
	}

	// Include content only for detail view
	if includeContent {
		response["content"] = berita.Content
		response["tags"] = berita.Tags
	}

	// Add published_at if not nil
	if berita.PublishedAt != nil {
		response["published_at"] = berita.PublishedAt
	}

	// Add deleted_at if not nil (for admin view)
	if berita.DeletedAt != nil {
		response["deleted_at"] = berita.DeletedAt
	}

	return response
}
