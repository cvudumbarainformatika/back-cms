package controllers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	models "github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type ContentController struct {
	DB *sqlx.DB
}

func NewContentController(db *sqlx.DB) *ContentController {
	return &ContentController{DB: db}
}

// Request/Response DTOs
type ContentImage struct {
	Src string `json:"src"`
}

type ContentBadge struct {
	Label string `json:"label"`
}

type ContentInput struct {
	Slug        string                `json:"slug" binding:"required"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Body        string                `json:"body"`
	HTML        string                `json:"html"`
	Date        string                `json:"date"` // ISO string
	Image       ContentImage          `json:"image"`
	Badge       ContentBadge          `json:"badge"`
	Authors     models.ContentAuthors `json:"authors"`
}

func (c *ContentController) InitTable() {
	schema := `
	CREATE TABLE IF NOT EXISTS content_pages (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		slug VARCHAR(255) NOT NULL UNIQUE,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		body LONGTEXT,
		html LONGTEXT,
		date DATETIME,
		image_src VARCHAR(255),
		badge_label VARCHAR(100),
		authors JSON,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_slug (slug)
	);`

	_, err := c.DB.Exec(schema)
	if err != nil {
		// Just log error, cannot return to Gin context here
		// In a real app we might want to panic or log fatal
	}
}

func (c *ContentController) GetContentBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	// Trim leading slash or trailing if any
	slug = strings.Trim(slug, "/")

	var content models.ContentPage
	query := `SELECT * FROM content_pages WHERE slug = ? LIMIT 1`
	err := c.DB.Get(&content, query, slug)

	if err == sql.ErrNoRows {
		// Return 404
		utils.Error(ctx, http.StatusNotFound, "content_not_found", "Content not found", nil)
		return
	} else if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "database_error", err.Error(), nil)
		return
	}

	// Map to response structure expected by frontend
	// Using anonymous struct for response
	response := struct {
		models.ContentPage
		Image ContentImage `json:"image"`
		Badge ContentBadge `json:"badge"`
	}{
		ContentPage: content,
		Image:       ContentImage{Src: content.ImageSrc},
		Badge:       ContentBadge{Label: content.BadgeLabel},
	}

	utils.Success(ctx, http.StatusOK, "Content fetched successfully", response)
}

func (c *ContentController) SaveContent(ctx *gin.Context) {
	var input ContentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "invalid_input", err.Error(), nil)
		return
	}

	// Parse date
	parsedDate, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		parsedDate, err = time.Parse("2006-01-02", input.Date)
		if err != nil {
			parsedDate = time.Now()
		}
	}

	// Clean slug
	input.Slug = strings.Trim(input.Slug, "/")

	// Check if exists
	var existsID int64
	checkQuery := `SELECT id FROM content_pages WHERE slug = ?`
	err = c.DB.Get(&existsID, checkQuery, input.Slug)

	if err == sql.ErrNoRows {
		// Insert
		insertQuery := `
			INSERT INTO content_pages (slug, title, description, body, html, date, image_src, badge_label, authors, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		`
		res, err := c.DB.Exec(insertQuery,
			input.Slug, input.Title, input.Description, input.Body, input.HTML,
			parsedDate, input.Image.Src, input.Badge.Label, input.Authors,
		)
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, "insert_error", err.Error(), nil)
			return
		}

		id, _ := res.LastInsertId()
		utils.Success(ctx, http.StatusCreated, "Content created successfully", gin.H{"id": id})

	} else if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "database_error", err.Error(), nil)
		return
	} else {
		// Update
		updateQuery := `
			UPDATE content_pages 
			SET title=?, description=?, body=?, html=?, date=?, image_src=?, badge_label=?, authors=?, updated_at=NOW()
			WHERE id=?
		`
		_, err := c.DB.Exec(updateQuery,
			input.Title, input.Description, input.Body, input.HTML,
			parsedDate, input.Image.Src, input.Badge.Label, input.Authors,
			existsID,
		)
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, "update_error", err.Error(), nil)
			return
		}

		utils.Success(ctx, http.StatusOK, "Content updated successfully", gin.H{"id": existsID})
	}
}
