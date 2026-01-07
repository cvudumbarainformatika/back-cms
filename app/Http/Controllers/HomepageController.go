package controllers

import (
	"net/http"

	requests "github.com/cvudumbarainformatika/backend/app/Http/Requests"
	models "github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type HomepageController struct {
	db *sqlx.DB
}

func NewHomepageController(db *sqlx.DB) *HomepageController {
	return &HomepageController{db: db}
}

func (hc *HomepageController) Get(c *gin.Context) {
	homepage, err := models.GetHomepage(hc.db)
	if err != nil {
		// If not found, return default structure or empty
		defaultContent := gin.H{
			"hero": gin.H{
				"title":       "Leading Respiratory Science",
				"description": "Wadah profesional kesehatan paru dan respirasi...",
				"images":      []string{},
			},
			"stats": []interface{}{},
			"seo": gin.H{
				"title":       "PDPI",
				"description": "Perhimpunan Dokter Paru Indonesia",
			},
		}
		utils.Success(c, http.StatusOK, "Default homepage data (DB empty)", defaultContent)
		return
	}

	utils.Success(c, http.StatusOK, "Homepage data fetched", homepage.Content)
}

func (hc *HomepageController) Update(c *gin.Context) {
	var req requests.UpdateHomepageRequest
	if err := req.Validate(c); err != nil {
		return
	}

	// Prepare content map
	content := map[string]interface{}{
		"hero":  req.Hero,
		"stats": req.Stats,
		"seo":   req.SEO,
	}

	homepage, err := models.UpdateHomepage(hc.db, content)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to update homepage", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Homepage updated successfully", homepage.Content)
}
