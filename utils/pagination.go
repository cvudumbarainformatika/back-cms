package utils

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OffsetPaginationMeta contains metadata for offset-based pagination
type OffsetPaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	LastPage    int   `json:"last_page"`
}

// OffsetPaginationResponse represents a paginated response with offset
type OffsetPaginationResponse struct {
	Data interface{}          `json:"data"`
	Meta OffsetPaginationMeta `json:"meta"`
}

// CursorPaginationResponse represents a paginated response with cursor
type CursorPaginationResponse struct {
	Data       interface{} `json:"data"`
	NextCursor *int64      `json:"next_cursor"`
}

// OffsetPaginate creates an offset-based pagination response
func OffsetPaginate(data interface{}, page, limit int, total int64) OffsetPaginationResponse {
	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	return OffsetPaginationResponse{
		Data: data,
		Meta: OffsetPaginationMeta{
			CurrentPage: page,
			PerPage:     limit,
			Total:       total,
			LastPage:    lastPage,
		},
	}
}

// CursorPaginate creates a cursor-based pagination response
func CursorPaginate(data interface{}, nextCursor *int64) CursorPaginationResponse {
	return CursorPaginationResponse{
		Data:       data,
		NextCursor: nextCursor,
	}
}

// LaravelPaginationResponse represents a paginated response matching Laravel's format
type LaravelPaginationResponse struct {
	CurrentPage  int         `json:"current_page"`
	Data         interface{} `json:"data"`
	FirstPageURL string      `json:"first_page_url"`
	From         int         `json:"from"`
	LastPage     int         `json:"last_page"`
	NextPageURL  *string     `json:"next_page_url"`
	Path         string      `json:"path"`
	PerPage      int         `json:"per_page"`
	PrevPageURL  *string     `json:"prev_page_url"`
	To           int         `json:"to"`
	Total        int64       `json:"total"`
}

// CreateLaravelPagination creates a Laravel-style pagination response
func CreateLaravelPagination(c *gin.Context, data interface{}, page, limit int, total int64) LaravelPaginationResponse {
	baseURL := "http://" + c.Request.Host + c.Request.URL.Path
	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	from := (page-1)*limit + 1
	to := page * limit
	if to > int(total) {
		to = int(total)
	}
	if total == 0 {
		from = 0
		to = 0
	}

	var nextPageURL, prevPageURL *string

	if page < lastPage {
		url := fmt.Sprintf("%s?page=%d&per_page=%d", baseURL, page+1, limit)
		nextPageURL = &url
	}

	if page > 1 {
		url := fmt.Sprintf("%s?page=%d&per_page=%d", baseURL, page-1, limit)
		prevPageURL = &url
	}

	firstPageURL := fmt.Sprintf("%s?page=1&per_page=%d", baseURL, limit)

	return LaravelPaginationResponse{
		CurrentPage:  page,
		Data:         data,
		FirstPageURL: firstPageURL,
		From:         from,
		LastPage:     lastPage,
		NextPageURL:  nextPageURL,
		Path:         baseURL,
		PerPage:      limit,
		PrevPageURL:  prevPageURL,
		To:           to,
		Total:        total,
	}
}

// GetPaginationParams extracts pagination parameters from query string
func GetPaginationParams(c *gin.Context) (page int, limit int) {
	page = 1
	limit = 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Check for both 'limit' and 'per_page'
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	} else if limitStr := c.Query("per_page"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	return page, limit
}

// GetCursorParams extracts cursor pagination parameters from query string
func GetCursorParams(c *gin.Context) (cursor int64, limit int) {
	cursor = 0
	limit = 10

	if cursorStr := c.Query("cursor"); cursorStr != "" {
		if cur, err := strconv.ParseInt(cursorStr, 10, 64); err == nil && cur > 0 {
			cursor = cur
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	return cursor, limit
}

// FilterParams holds common filter parameters
type FilterParams struct {
	Page    int
	PerPage int
	Q       string
	OrderBy string
	Sort    string
}

// GetFilterParams extracts filter parameters from query string
func GetFilterParams(c *gin.Context) FilterParams {
	page, limit := GetPaginationParams(c)

	return FilterParams{
		Page:    page,
		PerPage: limit,
		Q:       c.Query("q"),
		OrderBy: c.DefaultQuery("order_by", "created_at"),
		Sort:    c.DefaultQuery("sort", "desc"),
	}
}
