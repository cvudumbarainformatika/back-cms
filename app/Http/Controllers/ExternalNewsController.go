package controllers

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	models "github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/cvudumbarainformatika/backend/app/Services"
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type ExternalNewsController struct {
	DB             *sqlx.DB
	RSSService     *services.RSSService
	ScraperService *services.ScraperService
}

func NewExternalNewsController(db *sqlx.DB, rdb *redis.Client) *ExternalNewsController {
	return &ExternalNewsController{
		DB:             db,
		RSSService:     services.NewRSSService(db, rdb),
		ScraperService: services.NewScraperService(),
	}
}

func (c *ExternalNewsController) Index(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	dbCtx, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	isImported := ctx.Query("is_imported")

	items, total, err := models.GetAllExternalNews(dbCtx, c.DB, limit, offset, isImported)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if items == nil {
		items = []models.ExternalNews{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": items,
			"pagination": gin.H{
				"total": total,
				"page":  page,
				"limit": limit,
			},
		},
	})
}



func (c *ExternalNewsController) Import(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	// Get target status from query, default to draft
	targetStatus := ctx.DefaultQuery("status", "draft")

	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Get external news data
	var ext models.ExternalNews
	err := c.DB.GetContext(dbCtx, &ext, "SELECT * FROM external_news WHERE id = ?", id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Berita luar tidak ditemukan"})
		return
	}

	// 1.5 Scrape detail IF full_content is empty
	bestContent := ext.FullContent.String
	if bestContent == "" {
		fmt.Printf("[Import] Scraping content on-demand for: %s\n", ext.URL)
		scraped, err := c.ScraperService.ScrapeDetail(ext.URL)
		if err == nil {
			bestContent = scraped
			// Update cache in DB
			models.UpdateExternalNewsContent(dbCtx, c.DB, id, scraped)
		} else {
			fmt.Printf("[Import Warning] Scraping failed: %v\n", err)
			bestContent = ext.Description.String
		}
	}

	// 2. Map to local Berita model
	slug := utils.GenerateSlug(ext.Title)
	
	// Create pointer strings for model compatibility
	pExcerpt := ext.Description.String
	pImage := "" 
	pCat := "RSS"
	pAuthor := ext.Source
	pStatus := targetStatus

	// Try to extract image from description if pImage is empty
	if pImage == "" {
		re := regexp.MustCompile(`<img[^>]+src="([^">]+)"`)
		match := re.FindStringSubmatch(ext.Description.String)
		if len(match) > 1 {
			pImage = match[1]
		}
	}

	// Construct content with backlink to source
	content := fmt.Sprintf("<div>%s</div><br/><p><i>Baca selengkapnya di: <a href=\"%s\" target=\"_blank\" rel=\"noopener noreferrer\">%s</a></i></p>", 
		bestContent, ext.URL, ext.Source)

	// Set PublishedAt if status is published
	var publishedAt *time.Time
	if targetStatus == "published" {
		now := time.Now()
		publishedAt = &now
	}

	berita := models.Berita{
		Title:       ext.Title,   // string
		Slug:        slug,        // string
		Excerpt:     &pExcerpt,   // *string
		Content:     content,     // string
		ImageURL:    &pImage,     // *string
		Category:    &pCat,       // *string
		Author:      &pAuthor,    // *string
		Status:      &pStatus,    // *string
		PublishedAt: publishedAt, // *time.Time
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 3. Save to database
	res, err := c.DB.NamedExecContext(dbCtx, `
		INSERT INTO berita (title, slug, excerpt, content, image_url, category, author, status, published_at, created_at, updated_at)
		VALUES (:title, :slug, :excerpt, :content, :image_url, :category, :author, :status, :published_at, :created_at, :updated_at)
	`, &berita)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan berita: " + err.Error()})
		return
	}

	newID, _ := res.LastInsertId()

	// 4. Mark as imported
	c.DB.ExecContext(dbCtx, "UPDATE external_news SET is_imported = 1 WHERE id = ?", id)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berita berhasil diimpor sebagai " + targetStatus,
		"data": gin.H{
			"id": newID,
		},
	})
}
