package controllers

import (
	"net/http"
	"strconv"
	"strings"

	models "github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/cvudumbarainformatika/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type MemberController struct {
	DB *sqlx.DB
}

// NewMemberController creates a new member controller
func NewMemberController(db *sqlx.DB) *MemberController {
	return &MemberController{DB: db}
}

// GetMembers returns paginated list of PDPI members with search and filters
func (c *MemberController) GetMembers(ctx *gin.Context) {
	// Pagination
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Search
	search := strings.TrimSpace(ctx.Query("search"))

	// Filters
	cabang := strings.TrimSpace(ctx.Query("cabang"))
	provinsi := strings.TrimSpace(ctx.Query("provinsi"))
	status := strings.TrimSpace(ctx.Query("status"))

	// Build query
	query := "SELECT * FROM pdpi_members WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM pdpi_members WHERE 1=1"
	args := []interface{}{}

	// Search condition (nama, npa, email)
	if search != "" {
		searchCondition := " AND (nama LIKE ? OR npa LIKE ? OR email LIKE ?)"
		query += searchCondition
		countQuery += searchCondition
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}

	// Filter by cabang
	if cabang != "" {
		query += " AND cabang = ?"
		countQuery += " AND cabang = ?"
		args = append(args, cabang)
	}

	// Filter by provinsi
	if provinsi != "" {
		query += " AND provinsi = ?"
		countQuery += " AND provinsi = ?"
		args = append(args, provinsi)
	}

	// Filter by status
	if status != "" {
		query += " AND status = ?"
		countQuery += " AND status = ?"
		args = append(args, status)
	}

	// Get total count
	var total int
	err := c.DB.Get(&total, countQuery, args...)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "database_error", err.Error(), nil)
		return
	}

	// Add ordering and pagination
	query += " ORDER BY npa ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Get members
	members := make([]models.PDPIMember, 0)
	err = c.DB.Select(&members, query, args...)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "database_error", err.Error(), nil)
		return
	}

	// Calculate pagination metadata
	totalPages := (total + limit - 1) / limit

	utils.Success(ctx, http.StatusOK, "Members fetched successfully", gin.H{
		"items": members,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetMemberByID returns a single member by ID
func (c *MemberController) GetMemberByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var member models.PDPIMember
	query := "SELECT * FROM pdpi_members WHERE id = ? LIMIT 1"
	err := c.DB.Get(&member, query, id)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils.Error(ctx, http.StatusNotFound, "member_not_found", "Member not found", nil)
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "database_error", err.Error(), nil)
		return
	}

	utils.Success(ctx, http.StatusOK, "Member fetched successfully", member)
}

// UpdateMember updates a member (limited fields)
func (c *MemberController) UpdateMember(ctx *gin.Context) {
	id := ctx.Param("id")

	// Check if member exists
	var existing models.PDPIMember
	err := c.DB.Get(&existing, "SELECT * FROM pdpi_members WHERE id = ? LIMIT 1", id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils.Error(ctx, http.StatusNotFound, "member_not_found", "Member not found", nil)
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "database_error", err.Error(), nil)
		return
	}

	// Parse request body
	var input struct {
		UserID *int64 `json:"user_id"` // Allow linking to user account
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "invalid_input", err.Error(), nil)
		return
	}

	// Update only allowed fields
	query := "UPDATE pdpi_members SET user_id = ?, updated_at = NOW() WHERE id = ?"
	_, err = c.DB.Exec(query, input.UserID, id)

	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "database_error", err.Error(), nil)
		return
	}

	utils.Success(ctx, http.StatusOK, "Member updated successfully", nil)
}

// GetFilterOptions returns available filter options (cabang, provinsi, status)
func (c *MemberController) GetFilterOptions(ctx *gin.Context) {
	var cabangList []string
	var provinsiList []string
	var statusList []string

	// Get distinct cabang
	c.DB.Select(&cabangList, "SELECT DISTINCT cabang FROM pdpi_members WHERE cabang IS NOT NULL AND cabang != '' ORDER BY cabang")

	// Get distinct provinsi
	c.DB.Select(&provinsiList, "SELECT DISTINCT provinsi FROM pdpi_members WHERE provinsi IS NOT NULL AND provinsi != '' ORDER BY provinsi")

	// Get distinct status
	c.DB.Select(&statusList, "SELECT DISTINCT status FROM pdpi_members WHERE status IS NOT NULL AND status != '' ORDER BY status")

	utils.Success(ctx, http.StatusOK, "Filter options fetched successfully", gin.H{
		"cabang":   cabangList,
		"provinsi": provinsiList,
		"status":   statusList,
	})
}
