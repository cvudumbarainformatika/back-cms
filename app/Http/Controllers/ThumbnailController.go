package controllers

import (
	"net/http"
	"strconv"
	"strings"

	models "github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"encoding/json"
	"time"
	"context"
)

type ThumbnailController struct {
	DB    *sqlx.DB
	Redis *redis.Client
}

func NewThumbnailController(db *sqlx.DB, rdb *redis.Client) *ThumbnailController {
	return &ThumbnailController{DB: db, Redis: rdb}
}

func (tc *ThumbnailController) Index(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)
	category := c.Query("category")
	if category != "" {
		category = strings.ToUpper(category)
	}
	search := c.Query("search")
	offset := (page - 1) * limit

	items, total, err := models.GetAllThumbnails(tc.DB, category, search, offset, limit)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", err.Error(), nil)
		return
	}

	pagination := utils.OffsetPaginate(items, page, limit, total)
	utils.Success(c, http.StatusOK, "Thumbnails fetched", gin.H{
		"items":      pagination.Data,
		"pagination": pagination.Meta,
	})
}

func (tc *ThumbnailController) Categories(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "thumbnails:categories"

	// Try cache
	if tc.Redis != nil {
		cached, err := tc.Redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var cats []string
			if err := json.Unmarshal([]byte(cached), &cats); err == nil {
				utils.Success(c, http.StatusOK, "Categories fetched (cached)", cats)
				return
			}
		}
	}

	cats, err := models.GetThumbnailCategories(tc.DB)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", err.Error(), nil)
		return
	}

	// Save to cache
	if tc.Redis != nil {
		data, _ := json.Marshal(cats)
		tc.Redis.Set(ctx, cacheKey, data, 24*time.Hour)
	}

	utils.Success(c, http.StatusOK, "Categories fetched", cats)
}

func (tc *ThumbnailController) Store(c *gin.Context) {
	var item models.Thumbnail
	if err := c.ShouldBindJSON(&item); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_input", err.Error(), nil)
		return
	}

	item.Category = strings.ToUpper(item.Category)

	if err := item.Create(tc.DB); err != nil {
		utils.Error(c, http.StatusInternalServerError, "insert_error", err.Error(), nil)
		return
	}

	tc.invalidateCache()
	utils.Success(c, http.StatusCreated, "Thumbnail created", item)
}

func (tc *ThumbnailController) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var item models.Thumbnail
	if err := c.ShouldBindJSON(&item); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_input", err.Error(), nil)
		return
	}
	item.ID = id
	item.Category = strings.ToUpper(item.Category)

	if err := item.Update(tc.DB); err != nil {
		utils.Error(c, http.StatusInternalServerError, "update_error", err.Error(), nil)
		return
	}

	tc.invalidateCache()
	utils.Success(c, http.StatusOK, "Thumbnail updated", item)
}

func (tc *ThumbnailController) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := models.DeleteThumbnail(tc.DB, id); err != nil {
		utils.Error(c, http.StatusInternalServerError, "delete_error", err.Error(), nil)
		return
	}
	tc.invalidateCache()
	utils.Success(c, http.StatusOK, "Thumbnail deleted", nil)
}

func (tc *ThumbnailController) DeleteCategory(c *gin.Context) {
	category := c.Param("category")
	if err := models.DeleteThumbnailsByCategory(tc.DB, category); err != nil {
		utils.Error(c, http.StatusInternalServerError, "delete_error", err.Error(), nil)
		return
	}
	tc.invalidateCache()
	utils.Success(c, http.StatusOK, "Category and associated thumbnails deleted", nil)
}

func (tc *ThumbnailController) PublicGrouped(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "thumbnails:public_grouped"

	if tc.Redis != nil {
		cached, err := tc.Redis.Get(ctx, cacheKey).Result()
		if err == nil {
			c.Data(http.StatusOK, "application/json", []byte(cached))
			return
		}
	}

	// Fetch all thumbnails (ordered by created_at)
	items, _, err := models.GetAllThumbnails(tc.DB, "", "", 0, 1000)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", err.Error(), nil)
		return
	}

	// Group by category
	grouped := make(map[string][]models.Thumbnail)
	for _, item := range items {
		cat := strings.ToUpper(item.Category)
		grouped[cat] = append(grouped[cat], item)
	}

	resp := gin.H{
		"success": true,
		"data":    grouped,
	}

	if tc.Redis != nil {
		data, _ := json.Marshal(resp)
		tc.Redis.Set(ctx, cacheKey, data, 24*time.Hour)
	}

	c.JSON(http.StatusOK, resp)
}

func (tc *ThumbnailController) invalidateCache() {
	if tc.Redis != nil {
		ctx := context.Background()
		tc.Redis.Del(ctx, "thumbnails:categories")
		tc.Redis.Del(ctx, "thumbnails:public_grouped")
	}
}
