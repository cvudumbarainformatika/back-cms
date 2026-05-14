package controllers

import (
	"net/http"

	models "github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type TypeArtikelController struct {
	DB *sqlx.DB
}

func NewTypeArtikelController(db *sqlx.DB) *TypeArtikelController {
	return &TypeArtikelController{
		DB: db,
	}
}

// GET /api/v1/typeartikel
func (tc *TypeArtikelController) GetAll(c *gin.Context) {
	items, err := models.GetAllTypeArtikel(tc.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    items,
	})
}

// POST /api/v1/typeartikel
func (tc *TypeArtikelController) Create(c *gin.Context) {
	var payload struct {
		TypeArtikel string `json:"typeartikel"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Payload tidak valid",
		})
		return
	}

	item := &models.TypeArtikel{
		TypeArtikel: payload.TypeArtikel,
	}

	if err := item.Create(tc.DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal tambah type artikel",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil tambah type artikel",
		"data":    item,
	})
}
