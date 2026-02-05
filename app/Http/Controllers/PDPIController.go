package controllers

import (
	"fmt"
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
		Gelar:          utils.StringToPtr(pdpiMember.Gelar),
		Gelar2:         utils.StringToPtr(pdpiMember.Gelar2),
		Email:          utils.StringToPtr(pdpiMember.Email),
		NoHP:           utils.StringToPtr(pdpiMember.NoHP),
		NIK:            utils.StringToPtr(pdpiMember.NIK),
		JenisKelamin:   utils.StringToPtr(pdpiMember.JenisKelamin),
		TempatLahir:    utils.StringToPtr(pdpiMember.TempatLahir),
		AlamatRumah:    utils.StringToPtr(pdpiMember.AlamatRumah),
		Cabang:         utils.StringToPtr(pdpiMember.Cabang),
		Provinsi:       utils.StringToPtr(pdpiMember.Provinsi),
		KotaKabupaten:  utils.StringToPtr(pdpiMember.KotaKabupaten),
		Status:         utils.StringToPtr(pdpiMember.Status),
		Alumni:         utils.StringToPtr(pdpiMember.Alumni),
		ThnLulus:       utils.Int64ToPtr(int64(pdpiMember.ThnLulus)),
		TempatTugas:    utils.StringToPtr(pdpiMember.TempatTugas),
		TempatPraktek1: utils.StringToPtr(pdpiMember.TempatPraktek1),
		TempatPraktek2: utils.StringToPtr(pdpiMember.TempatPraktek2),
		Subspesialis:   utils.StringToPtr(pdpiMember.Subspesialis),
		NoSTR:          utils.StringToPtr(pdpiMember.NoSTR),
		NoSIP:          utils.StringToPtr(pdpiMember.NoSIP),
		SyncedAt:       utils.TimeToPtr(time.Now()),
	}

	// Parse dates if provided
	if pdpiMember.TglLahir != "" {
		if t, err := time.Parse("2006-01-02", pdpiMember.TglLahir); err == nil {
			localMember.TglLahir = utils.TimeToPtr(t)
		}
	}
	if pdpiMember.STRBerlakuSampai != "" {
		if t, err := time.Parse("2006-01-02", pdpiMember.STRBerlakuSampai); err == nil {
			localMember.STRBerlakuSampai = utils.TimeToPtr(t)
		}
	}
	if pdpiMember.SIPBerlakuSampai != "" {
		if t, err := time.Parse("2006-01-02", pdpiMember.SIPBerlakuSampai); err == nil {
			localMember.SIPBerlakuSampai = utils.TimeToPtr(t)
		}
	}

	// Link to current user if emails match
	if emailToSync == user.Email {
		localMember.UserID = utils.Int64ToPtr(userIDInt64)
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

// SyncAllMembers syncs ALL PDPI members from API to local database
// POST /api/v1/pdpi/sync-all-members
// This is an admin-only operation with pagination handling
func (pc *PDPIController) SyncAllMembers(c *gin.Context) {
	startTime := time.Now()

	// Statistics tracking
	var (
		totalFetched  int
		totalSynced   int
		totalFailed   int
		errorMessages []string
	)

	// Fetch all pages from PDPI API
	page := 1
	limit := 100 // Fetch 100 members per page
	hasMore := true

	for hasMore {
		// Fetch current page
		filter := services.MembersFilter{
			Page:  page,
			Limit: limit,
		}

		resp, err := pc.pdpiService.GetMembers(filter)
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to fetch page %d: %v", page, err)
			errorMessages = append(errorMessages, errorMsg)
			break
		}

		if !resp.Success {
			errorMsg := fmt.Sprintf("PDPI API returned unsuccessful response for page %d", page)
			errorMessages = append(errorMessages, errorMsg)
			break
		}

		// Process each member in current page
		for _, pdpiMember := range resp.Data {
			totalFetched++

			// Convert to local model
			localMember := &models.PDPIMember{
				ID:             pdpiMember.ID,
				NPA:            pdpiMember.NPA,
				Nama:           pdpiMember.Nama,
				Gelar:          utils.StringToPtr(pdpiMember.Gelar),
				Gelar2:         utils.StringToPtr(pdpiMember.Gelar2),
				Email:          utils.StringToPtr(pdpiMember.Email),
				NoHP:           utils.StringToPtr(pdpiMember.NoHP),
				NIK:            utils.StringToPtr(pdpiMember.NIK),
				JenisKelamin:   utils.StringToPtr(pdpiMember.JenisKelamin),
				TempatLahir:    utils.StringToPtr(pdpiMember.TempatLahir),
				AlamatRumah:    utils.StringToPtr(pdpiMember.AlamatRumah),
				Cabang:         utils.StringToPtr(pdpiMember.Cabang),
				Provinsi:       utils.StringToPtr(pdpiMember.Provinsi),
				KotaKabupaten:  utils.StringToPtr(pdpiMember.KotaKabupaten),
				Status:         utils.StringToPtr(pdpiMember.Status),
				Alumni:         utils.StringToPtr(pdpiMember.Alumni),
				ThnLulus:       utils.Int64ToPtr(int64(pdpiMember.ThnLulus)),
				TempatTugas:    utils.StringToPtr(pdpiMember.TempatTugas),
				TempatPraktek1: utils.StringToPtr(pdpiMember.TempatPraktek1),
				TempatPraktek2: utils.StringToPtr(pdpiMember.TempatPraktek2),
				Subspesialis:   utils.StringToPtr(pdpiMember.Subspesialis),
				NoSTR:          utils.StringToPtr(pdpiMember.NoSTR),
				NoSIP:          utils.StringToPtr(pdpiMember.NoSIP),
				SyncedAt:       utils.TimeToPtr(time.Now()),
			}

			// Parse dates
			if pdpiMember.TglLahir != "" {
				if t, err := time.Parse("2006-01-02", pdpiMember.TglLahir); err == nil {
					localMember.TglLahir = utils.TimeToPtr(t)
				}
			}
			if pdpiMember.STRBerlakuSampai != "" {
				if t, err := time.Parse("2006-01-02", pdpiMember.STRBerlakuSampai); err == nil {
					localMember.STRBerlakuSampai = utils.TimeToPtr(t)
				}
			}
			if pdpiMember.SIPBerlakuSampai != "" {
				if t, err := time.Parse("2006-01-02", pdpiMember.SIPBerlakuSampai); err == nil {
					localMember.SIPBerlakuSampai = utils.TimeToPtr(t)
				}
			}

			// Try to link to existing user by email
			if pdpiMember.Email != "" {
				user, _ := models.FindByEmail(pc.db, pdpiMember.Email)
				if user != nil {
					localMember.UserID = utils.Int64ToPtr(user.ID)
				}
			}

			// Upsert to database
			err = models.UpsertPDPIMember(pc.db, localMember)
			if err != nil {
				totalFailed++
				errorMsg := fmt.Sprintf("Failed to sync member %s (NPA: %s): %v", pdpiMember.Nama, pdpiMember.NPA, err)
				errorMessages = append(errorMessages, errorMsg)
				// Continue processing other members
				continue
			}

			totalSynced++
		}

		// Check if there are more pages
		// Use length of data instead of TotalPages since some APIs don't return reliable TotalPages
		if len(resp.Data) < limit {
			// If we got less data than requested, we've reached the end
			hasMore = false
		} else {
			// Continue to next page
			page++
		}
	}

	duration := time.Since(startTime)

	// Prepare response
	result := gin.H{
		"total_fetched": totalFetched,
		"total_synced":  totalSynced,
		"total_failed":  totalFailed,
		"duration_ms":   duration.Milliseconds(),
		"pages_fetched": page,
	}

	if len(errorMessages) > 0 {
		result["errors"] = errorMessages
	}

	if totalFailed > 0 {
		utils.Success(c, http.StatusOK, fmt.Sprintf("Sync completed with %d failures", totalFailed), result)
	} else {
		utils.Success(c, http.StatusOK, "All members synced successfully", result)
	}
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

		query += " ORDER BY CAST(npa AS UNSIGNED) ASC"

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

// SearchPublicMembers searches members for public directory
// GET /api/v1/members/search?nama=&cabang=&provinsi=&status=&page=&limit=
// This is a public endpoint for member directory feature
func (pc *PDPIController) SearchPublicMembers(c *gin.Context) {
	// Build query with filters
	query := "SELECT id, npa, nama, gelar, gelar2, email, cabang, provinsi, kota_kabupaten, status, tempat_praktek_1, tempat_praktek_2, alumni FROM pdpi_members WHERE 1=1"
	var args []interface{}

	// Apply filters
	if nama := c.Query("nama"); nama != "" {
		query += " AND nama LIKE ?"
		args = append(args, "%"+nama+"%")
	}
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

	// Order by name
	query += " ORDER BY CAST(npa AS UNSIGNED) ASC"

	// Pagination
	page := utils.QueryInt(c, "page", 1)
	limit := utils.QueryInt(c, "limit", 12)
	offset := (page - 1) * limit

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Execute query
	type PublicMember struct {
		ID             string  `db:"id" json:"id"`
		NPA            string  `db:"npa" json:"npa"`
		Nama           string  `db:"nama" json:"nama"`
		Gelar          *string `db:"gelar" json:"gelar"`
		Gelar2         *string `db:"gelar2" json:"gelar2"`
		JenisKelamin   *string `db:"jenis_kelamin" json:"jenis_kelamin"`
		Email          *string `db:"email" json:"email"`
		Cabang         *string `db:"cabang" json:"cabang"`
		Provinsi       *string `db:"provinsi" json:"provinsi"`
		KotaKabupaten  *string `db:"kota_kabupaten" json:"kota_kabupaten"`
		Status         *string `db:"status" json:"status"`
		TempatPraktek1 *string `db:"tempat_praktek_1" json:"tempat_praktek_1"`
		TempatPraktek2 *string `db:"tempat_praktek_2" json:"tempat_praktek_2"`
		Alumni         *string `db:"alumni" json:"alumni"`
	}

	// Query with COALESCE to handle NULL values
	query = "SELECT id, npa, nama, " +
		"COALESCE(gelar, '') as gelar, " +
		"COALESCE(gelar2, '') as gelar2, " +
		"COALESCE(jenis_kelamin, '') as jenis_kelamin, " +
		"COALESCE(email, '') as email, " +
		"COALESCE(cabang, '') as cabang, " +
		"COALESCE(provinsi, '') as provinsi, " +
		"COALESCE(kota_kabupaten, '') as kota_kabupaten, " +
		"COALESCE(status, '') as status, " +
		"COALESCE(tempat_praktek_1, '') as tempat_praktek_1, " +
		"COALESCE(tempat_praktek_2, '') as tempat_praktek_2, " +
		"COALESCE(alumni, '') as alumni " +
		"FROM pdpi_members WHERE nama IS NOT NULL AND nama != ''"
	args = []interface{}{}

	// Apply filters
	// Global search: search in nama, cabang, provinsi, tempat_praktek_1, tempat_praktek_2
	if search := c.Query("nama"); search != "" {
		query += " AND (nama LIKE ? OR cabang LIKE ? OR provinsi LIKE ? OR tempat_praktek_1 LIKE ? OR tempat_praktek_2 LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}
	// Specific filters (exact match)
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

	// Order by name
	query += " ORDER BY CAST(npa AS UNSIGNED) ASC"

	// Pagination
	page = utils.QueryInt(c, "page", 1)
	limit = utils.QueryInt(c, "limit", 12)
	offset = (page - 1) * limit

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	members := []PublicMember{} // Initialize as empty slice to return [] instead of null
	err := pc.db.Select(&members, query, args...)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch members: "+err.Error(), nil)
		return
	}

	// Count total for pagination
	countQuery := "SELECT COUNT(*) FROM pdpi_members WHERE nama IS NOT NULL AND nama != ''"
	var countArgs []interface{}

	// Reapply filters for count
	// Global search: search in nama, cabang, provinsi, tempat_praktek_1, tempat_praktek_2
	if search := c.Query("nama"); search != "" {
		countQuery += " AND (nama LIKE ? OR cabang LIKE ? OR provinsi LIKE ? OR tempat_praktek_1 LIKE ? OR tempat_praktek_2 LIKE ?)"
		searchPattern := "%" + search + "%"
		countArgs = append(countArgs, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}
	// Specific filters (exact match)
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

	var total int
	err = pc.db.Get(&total, countQuery, countArgs...)
	if err != nil {
		total = 0
	}

	totalPages := (total + limit - 1) / limit

	utils.Success(c, http.StatusOK, "Members retrieved successfully", gin.H{
		"members": members,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}
