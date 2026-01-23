package controllers

import (
	"database/sql"
	"net/http"
	"time"

	models "github.com/cvudumbarainformatika/backend/app/Models"
	services "github.com/cvudumbarainformatika/backend/app/Services"
	"github.com/cvudumbarainformatika/backend/config"
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// PDPIController handles PDPI API integration
type PDPIController struct {
	db          *sqlx.DB
	pdpiService *services.PDPIService
}

// NewPDPIController creates a new PDPI controller instance
func NewPDPIController(db *sqlx.DB, cfg *config.Config) *PDPIController {
	return &PDPIController{
		db:          db,
		pdpiService: services.NewPDPIService(&cfg.PDPI),
	}
}

// SyncMember syncs a PDPI member data with local database using user's email
// POST /api/v1/pdpi/sync-member
// Request body: { "email": "user@example.com" } or empty (will use authenticated user's email)
func (pc *PDPIController) SyncMember(c *gin.Context) {
	// Get authenticated user dari middleware
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "unauthorized", "User not authenticated", nil)
		return
	}

	// Convert userID to int64
	var userIDInt64 int64
	switch v := userIDVal.(type) {
	case int64:
		userIDInt64 = v
	case int:
		userIDInt64 = int64(v)
	case string:
		var err error
		userIDInt64, err = utils.StringToInt64(v)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "invalid_user_id", "Invalid user ID format", nil)
			return
		}
	default:
		utils.Error(c, http.StatusInternalServerError, "invalid_user_id", "Unsupported user ID type", nil)
		return
	}

	// Get user data
	user, err := models.FindByID(pc.db, userIDInt64)
	if err != nil || user == nil {
		utils.Error(c, http.StatusNotFound, "user_not_found", "User not found", nil)
		return
	}

	// Parse request (optional email)
	var reqBody struct {
		Email string `json:"email"`
	}
	_ = c.ShouldBindJSON(&reqBody)

	// Use user's email if not provided
	emailToSync := reqBody.Email
	if emailToSync == "" {
		emailToSync = user.Email
	}

	// Get member data from PDPI API
	pdpiMember, err := pc.pdpiService.GetMemberByEmail(emailToSync)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "member_not_found", "PDPI member not found: "+err.Error(), nil)
		return
	}

	// Convert PDPI API response to local model
	localMember := &models.PDPIMember{
		ID:             pdpiMember.ID,
		NPA:            pdpiMember.NPA,
		Nama:           pdpiMember.Nama,
		Gelar:          sql.NullString{String: pdpiMember.Gelar, Valid: pdpiMember.Gelar != ""},
		Gelar2:         sql.NullString{String: pdpiMember.Gelar2, Valid: pdpiMember.Gelar2 != ""},
		Email:          sql.NullString{String: pdpiMember.Email, Valid: pdpiMember.Email != ""},
		NoHP:           sql.NullString{String: pdpiMember.NoHP, Valid: pdpiMember.NoHP != ""},
		NIK:            sql.NullString{String: pdpiMember.NIK, Valid: pdpiMember.NIK != ""},
		JenisKelamin:   sql.NullString{String: pdpiMember.JenisKelamin, Valid: pdpiMember.JenisKelamin != ""},
		TempatLahir:    sql.NullString{String: pdpiMember.TempatLahir, Valid: pdpiMember.TempatLahir != ""},
		AlamatRumah:    sql.NullString{String: pdpiMember.AlamatRumah, Valid: pdpiMember.AlamatRumah != ""},
		Cabang:         sql.NullString{String: pdpiMember.Cabang, Valid: pdpiMember.Cabang != ""},
		Provinsi:       sql.NullString{String: pdpiMember.Provinsi, Valid: pdpiMember.Provinsi != ""},
		KotaKabupaten:  sql.NullString{String: pdpiMember.KotaKabupaten, Valid: pdpiMember.KotaKabupaten != ""},
		Status:         sql.NullString{String: pdpiMember.Status, Valid: pdpiMember.Status != ""},
		Alumni:         sql.NullString{String: pdpiMember.Alumni, Valid: pdpiMember.Alumni != ""},
		ThnLulus:       sql.NullInt64{Int64: int64(pdpiMember.ThnLulus), Valid: pdpiMember.ThnLulus > 0},
		TempatTugas:    sql.NullString{String: pdpiMember.TempatTugas, Valid: pdpiMember.TempatTugas != ""},
		TempatPraktek1: sql.NullString{String: pdpiMember.TempatPraktek1, Valid: pdpiMember.TempatPraktek1 != ""},
		TempatPraktek2: sql.NullString{String: pdpiMember.TempatPraktek2, Valid: pdpiMember.TempatPraktek2 != ""},
		Subspesialis:   sql.NullString{String: pdpiMember.Subspesialis, Valid: pdpiMember.Subspesialis != ""},
		NoSTR:          sql.NullString{String: pdpiMember.NoSTR, Valid: pdpiMember.NoSTR != ""},
		NoSIP:          sql.NullString{String: pdpiMember.NoSIP, Valid: pdpiMember.NoSIP != ""},
		SyncedAt:       sql.NullTime{Time: time.Now(), Valid: true},
	}

	// Parse dates if provided
	if pdpiMember.TglLahir != "" {
		if t, err := time.Parse("2006-01-02", pdpiMember.TglLahir); err == nil {
			localMember.TglLahir = sql.NullTime{Time: t, Valid: true}
		}
	}
	if pdpiMember.STRBerlakuSampai != "" {
		if t, err := time.Parse("2006-01-02", pdpiMember.STRBerlakuSampai); err == nil {
			localMember.STRBerlakuSampai = sql.NullTime{Time: t, Valid: true}
		}
	}
	if pdpiMember.SIPBerlakuSampai != "" {
		if t, err := time.Parse("2006-01-02", pdpiMember.SIPBerlakuSampai); err == nil {
			localMember.SIPBerlakuSampai = sql.NullTime{Time: t, Valid: true}
		}
	}

	// Link to current user if emails match
	if emailToSync == user.Email {
		localMember.UserID = sql.NullInt64{Int64: userIDInt64, Valid: true}
	}

	// Upsert to database
	err = models.UpsertPDPIMember(pc.db, localMember)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to save member data: "+err.Error(), nil)
		return
	}

	utils.Success(c, http.StatusOK, "Member synced successfully", gin.H{
		"member": localMember,
	})
}

// GetMembers retrieves PDPI members list from local database or API
// GET /api/v1/pdpi/members?source=local|api&page=&limit=&cabang=&provinsi=&status=&search=
func (pc *PDPIController) GetMembers(c *gin.Context) {
	source := c.DefaultQuery("source", "local") // local or api

	if source == "api" {
		// Fetch from PDPI API
		filter := services.MembersFilter{
			Page:     utils.QueryInt(c, "page", 1),
			Limit:    utils.QueryInt(c, "limit", 100),
			Cabang:   c.Query("cabang"),
			Provinsi: c.Query("provinsi"),
			Status:   c.Query("status"),
			Search:   c.Query("search"),
		}

		resp, err := pc.pdpiService.GetMembers(filter)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "pdpi_api_error", "Failed to fetch from PDPI API: "+err.Error(), nil)
			return
		}

		utils.Success(c, http.StatusOK, "Members fetched from PDPI API", gin.H{
			"members":    resp.Data,
			"pagination": resp.Pagination,
			"source":     "api",
		})
	} else {
		// Fetch from local database
		query := "SELECT * FROM pdpi_members WHERE 1=1"
		var args []interface{}

		// Apply filters
		if cabang := c.Query("cabang"); cabang != "" {
			query += " AND cabang = ?"
			args = append(args, cabang)
		}
		if provinsi := c.Query("provinsi"); provinsi != "" {
			query += " AND provinsi = ?"
			args = append(args, provinsi)
		}
		if status := c.Query("status"); status != "" {
			query += " AND status = ?"
			args = append(args, status)
		}
		if search := c.Query("search"); search != "" {
			query += " AND (nama LIKE ? OR email LIKE ? OR npa LIKE ?)"
			searchPattern := "%" + search + "%"
			args = append(args, searchPattern, searchPattern, searchPattern)
		}

		query += " ORDER BY nama ASC"

		// Pagination
		page := utils.QueryInt(c, "page", 1)
		limit := utils.QueryInt(c, "limit", 100)
		offset := (page - 1) * limit

		query += " LIMIT ? OFFSET ?"
		args = append(args, limit, offset)

		var members []models.PDPIMember
		err := pc.db.Select(&members, query, args...)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch members: "+err.Error(), nil)
			return
		}

		// Count total
		countQuery := "SELECT COUNT(*) FROM pdpi_members WHERE 1=1"
		var countArgs []interface{}
		// Reapply filters for count
		if cabang := c.Query("cabang"); cabang != "" {
			countQuery += " AND cabang = ?"
			countArgs = append(countArgs, cabang)
		}
		if provinsi := c.Query("provinsi"); provinsi != "" {
			countQuery += " AND provinsi = ?"
			countArgs = append(countArgs, provinsi)
		}
		if status := c.Query("status"); status != "" {
			countQuery += " AND status = ?"
			countArgs = append(countArgs, status)
		}
		if search := c.Query("search"); search != "" {
			countQuery += " AND (nama LIKE ? OR email LIKE ? OR npa LIKE ?)"
			searchPattern := "%" + search + "%"
			countArgs = append(countArgs, searchPattern, searchPattern, searchPattern)
		}

		var total int
		err = pc.db.Get(&total, countQuery, countArgs...)
		if err != nil {
			total = 0
		}

		totalPages := (total + limit - 1) / limit

		utils.Success(c, http.StatusOK, "Members fetched from local database", gin.H{
			"members": members,
			"pagination": gin.H{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": totalPages,
			},
			"source": "local",
		})
	}
}

// GetMemberByNPA retrieves a specific member by NPA
// GET /api/v1/pdpi/member/:npa
func (pc *PDPIController) GetMemberByNPA(c *gin.Context) {
	npa := c.Param("npa")
	if npa == "" {
		utils.Error(c, http.StatusBadRequest, "invalid_npa", "NPA is required", nil)
		return
	}

	// Try from local database first
	member, err := models.FindPDPIMemberByNPA(pc.db, npa)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Database error: "+err.Error(), nil)
		return
	}

	if member != nil {
		utils.Success(c, http.StatusOK, "Member found in local database", gin.H{
			"member": member,
			"source": "local",
		})
		return
	}

	// If not found locally, try PDPI API
	pdpiMember, err := pc.pdpiService.GetMemberByNPA(npa)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "member_not_found", "Member not found: "+err.Error(), nil)
		return
	}

	utils.Success(c, http.StatusOK, "Member found in PDPI API", gin.H{
		"member": pdpiMember,
		"source": "api",
	})
}

// GetMyMemberData retrieves PDPI member data for authenticated user
// GET /api/v1/pdpi/me
func (pc *PDPIController) GetMyMemberData(c *gin.Context) {
	// Get authenticated user
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "unauthorized", "User not authenticated", nil)
		return
	}

	// Convert userID to int64
	var userIDInt64 int64
	switch v := userIDVal.(type) {
	case int64:
		userIDInt64 = v
	case int:
		userIDInt64 = int64(v)
	case string:
		var err error
		userIDInt64, err = utils.StringToInt64(v)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "invalid_user_id", "Invalid user ID format", nil)
			return
		}
	default:
		utils.Error(c, http.StatusInternalServerError, "invalid_user_id", "Unsupported user ID type", nil)
		return
	}

	// Get user data
	user, err := models.FindByID(pc.db, userIDInt64)
	if err != nil || user == nil {
		utils.Error(c, http.StatusNotFound, "user_not_found", "User not found", nil)
		return
	}

	// Try to find member data in local database by email
	member, err := models.FindPDPIMemberByEmail(pc.db, user.Email)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Database error: "+err.Error(), nil)
		return
	}

	if member != nil {
		utils.Success(c, http.StatusOK, "Member data found", gin.H{
			"member": member,
			"source": "local",
		})
		return
	}

	// If not found locally, inform user to sync
	utils.Error(c, http.StatusNotFound, "member_not_synced", "Member data not synced yet. Please call /api/v1/pdpi/sync-member", nil)
}
