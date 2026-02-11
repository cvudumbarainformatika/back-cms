package controllers

import (
	"fmt"
	"net/http"

	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type DashboardController struct {
	db *sqlx.DB
}

func NewDashboardController(db *sqlx.DB) *DashboardController {
	return &DashboardController{db: db}
}

// GetStats returns dashboard statistics
// GET /api/v1/dashboard/stats
func (dc *DashboardController) GetStats(c *gin.Context) {
	// Get user info from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "unauthorized", "User not authenticated", nil)
		return
	}

	// Fetch user role from database to ensure accuracy
	// Context role might be missing or incomplete
	var role string
	err := dc.db.Get(&role, "SELECT role FROM users WHERE id = ?", userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch user role", nil)
		return
	}

	var articleCount int
	var agendaCount int
	var memberCount int

	// 1. Article Count
	// Logic: Admin (any admin role) sees total, Member sees their own
	// Check if role starts with "admin" or is "superadmin"
	isAdmin := role == "superadmin" || (len(role) >= 5 && role[:5] == "admin")

	if isAdmin {
		query := "SELECT COUNT(*) FROM berita WHERE deleted_at IS NULL"
		err = dc.db.Get(&articleCount, query)
	} else {
		// Convert userID to string for author_id comparison
		userIDStr := fmt.Sprintf("%v", userID)
		query := "SELECT COUNT(*) FROM berita WHERE author_id = ? AND deleted_at IS NULL"
		err = dc.db.Get(&articleCount, query, userIDStr)
	}

	if err != nil {
		// Log error but continue with 0
		fmt.Printf("Error counting articles: %v\n", err)
		articleCount = 0
	}

	// 2. Agenda Count (Total Active)
	// Logic: Count all non-deleted agenda
	agendaQuery := "SELECT COUNT(*) FROM agenda WHERE deleted_at IS NULL"
	err = dc.db.Get(&agendaCount, agendaQuery)
	if err != nil {
		fmt.Printf("Error counting agenda: %v\n", err)
		agendaCount = 0
	}

	// 3. Member Count (Total Directory)
	// Logic: Count all profiles in pdpi_members
	memberQuery := "SELECT COUNT(*) FROM pdpi_members"
	err = dc.db.Get(&memberCount, memberQuery)
	if err != nil {
		fmt.Printf("Error counting members: %v\n", err)
		memberCount = 0
	}

	utils.Success(c, http.StatusOK, "Dashboard stats retrieved", gin.H{
		"article_count": articleCount,
		"agenda_count":  agendaCount,
		"member_count":  memberCount,
	})
}
