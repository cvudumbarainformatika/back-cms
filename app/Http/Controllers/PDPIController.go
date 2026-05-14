package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	models "github.com/cvudumbarainformatika/backend/app/Models"
	services "github.com/cvudumbarainformatika/backend/app/Services"
	"github.com/cvudumbarainformatika/backend/config"
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
)

// PDPIController handles PDPI API integration
type PDPIController struct {
	db          *sqlx.DB
	pdpiService *services.PDPIService
	waba360     *config.WABA360Config

	localMembersByNPA   map[string]*models.PDPIMember
	localMembersByEmail map[string]*models.PDPIMember
}

// NewPDPIController creates a new PDPI controller instance
func NewPDPIController(db *sqlx.DB, cfg *config.Config) *PDPIController {
	return &PDPIController{
		db:          db,
		pdpiService: services.NewPDPIService(&cfg.PDPI),
		waba360:     &cfg.WABA360,
	}
}
func KeepExistingString(newVal, oldVal *string) *string {
	if newVal != nil && *newVal != "" {
		return newVal
	}
	return oldVal
}

func (pc *PDPIController) prepareLocalMemberCache() error {

	var localMembers []models.PDPIMember

	err := pc.db.Select(
		&localMembers,
		"SELECT * FROM pdpi_members",
	)
	if err != nil {
		return err
	}

	pc.localMembersByNPA = make(map[string]*models.PDPIMember)
	pc.localMembersByEmail = make(map[string]*models.PDPIMember)

	for i := range localMembers {

		member := &localMembers[i]

		if member.NPA != "" {
			pc.localMembersByNPA[member.NPA] = member
		}

		if member.Email != nil {
			pc.localMembersByEmail[*member.Email] = member
		}
	}

	return nil
}
// Helper to map Supabase response to Local Model
func (pc *PDPIController) mapSupabaseToLocal(pdpiMember services.SupabaseMember) *models.PDPIMember {
	// localMember := &models.PDPIMember{
	// 	ID:       pdpiMember.ID,
	// 	Nama:     pdpiMember.Nama,
	// 	SyncedAt: utils.TimeToPtr(time.Now()),
	// }
	var existing *models.PDPIMember

	npa := utils.PtrToString(pdpiMember.NPA)

	if npa != "" {
		existing = pc.localMembersByNPA[npa]
	}

	if existing == nil && pdpiMember.Email != nil {
		existing = pc.localMembersByEmail[*pdpiMember.Email]
	}

	var localMember *models.PDPIMember

	if existing != nil {
		localMember = existing
	} else {
		localMember = &models.PDPIMember{}
	}

	localMember.ID = pdpiMember.ID
	localMember.Nama = pdpiMember.Nama
	localMember.SyncedAt = utils.TimeToPtr(time.Now())

	
	// Basic Fields
	localMember.NPA = utils.PtrToString(pdpiMember.NPA) // Helper to handle nil pointer
	if pdpiMember.NPANumeric != nil {
		localMember.NPANumeric = pdpiMember.NPANumeric
	}
	if pdpiMember.Foto != nil {
		localMember.Foto = pdpiMember.Foto
	}
	if pdpiMember.Gelar != nil {
		localMember.Gelar = pdpiMember.Gelar
	}
	if pdpiMember.Gelar2 != nil {
		localMember.Gelar2 = pdpiMember.Gelar2
	}
	if pdpiMember.Email != nil {
		// localMember.Email = pdpiMember.Email
		localMember.Email = KeepExistingString(
												pdpiMember.Email,
												existing.Email,
											)
	}
	if pdpiMember.NoHP != nil {
		// localMember.NoHP = utils.StringToPtr(utils.NormalizePhoneNumber(*pdpiMember.NoHP))
		localMember.NoHP = KeepExistingString(
												pdpiMember.NoHP,
												existing.NoHP,
											)
	}
	
	if pdpiMember.NIK != nil {
		localMember.NIK = pdpiMember.NIK
	}
	if pdpiMember.JenisKelamin != nil {
		localMember.JenisKelamin = pdpiMember.JenisKelamin
	}
	if pdpiMember.TempatLahir != nil {
		localMember.TempatLahir = pdpiMember.TempatLahir
	}
	if pdpiMember.AlamatRumah != nil {
		localMember.AlamatRumah = pdpiMember.AlamatRumah
	}
	if pdpiMember.Cabang != nil {
		localMember.Cabang = pdpiMember.Cabang
	}
	if pdpiMember.Provinsi != nil {
		localMember.Provinsi = pdpiMember.Provinsi
	}
	if pdpiMember.KotaKabupaten != nil {
		localMember.KotaKabupaten = pdpiMember.KotaKabupaten
	}
	if pdpiMember.KotaKabupatenKantor != nil {
		localMember.KotaKabupatenKantor = pdpiMember.KotaKabupatenKantor
	}
	if pdpiMember.ProvinsiKantor != nil {
		localMember.ProvinsiKantor = pdpiMember.ProvinsiKantor
	}
	if pdpiMember.Status != nil {
		localMember.Status = pdpiMember.Status
	}
	if pdpiMember.Alumni != nil {
		localMember.Alumni = pdpiMember.Alumni
	}
	if pdpiMember.ThnLulus != nil {
		localMember.ThnLulus = pdpiMember.ThnLulus
	}
	if pdpiMember.TempatTugas != nil {
		localMember.TempatTugas = pdpiMember.TempatTugas
	}
	if pdpiMember.Subspesialis != nil {
		localMember.Subspesialis = pdpiMember.Subspesialis
	}
	if pdpiMember.GelarFISR != nil {
		localMember.GelarFISR = pdpiMember.GelarFISR
	}

	// Praktek 1
	if pdpiMember.TempatPraktek1 != nil {
		localMember.TempatPraktek1 = pdpiMember.TempatPraktek1
	}
	if pdpiMember.TempatPraktek1Tipe != nil {
		localMember.TempatPraktek1Tipe = pdpiMember.TempatPraktek1Tipe
	}
	if pdpiMember.TempatPraktek1Tipe2 != nil {
		localMember.TempatPraktek1Tipe2 = pdpiMember.TempatPraktek1Tipe2
	}
	if pdpiMember.TempatPraktek1Alkes != nil {
		localMember.TempatPraktek1Alkes = pdpiMember.TempatPraktek1Alkes
	}
	if pdpiMember.TempatPraktek1Alkes2 != nil {
		localMember.TempatPraktek1Alkes2 = pdpiMember.TempatPraktek1Alkes2
	}

	// Praktek 2
	if pdpiMember.TempatPraktek2 != nil {
		localMember.TempatPraktek2 = pdpiMember.TempatPraktek2
	}
	if pdpiMember.TempatPraktek2Tipe != nil {
		localMember.TempatPraktek2Tipe = pdpiMember.TempatPraktek2Tipe
	}
	if pdpiMember.TempatPraktek2Tipe2 != nil {
		localMember.TempatPraktek2Tipe2 = pdpiMember.TempatPraktek2Tipe2
	}
	if pdpiMember.TempatPraktek2Alkes != nil {
		localMember.TempatPraktek2Alkes = pdpiMember.TempatPraktek2Alkes
	}
	if pdpiMember.TempatPraktek2Alkes2 != nil {
		localMember.TempatPraktek2Alkes2 = pdpiMember.TempatPraktek2Alkes2
	}
	if pdpiMember.KotaKabupatenPraktek2 != nil {
		localMember.KotaKabupatenPraktek2 = pdpiMember.KotaKabupatenPraktek2
	}
	if pdpiMember.ProvinsiPraktek2 != nil {
		localMember.ProvinsiPraktek2 = pdpiMember.ProvinsiPraktek2
	}

	// Praktek 3
	if pdpiMember.TempatPraktek3 != nil {
		localMember.TempatPraktek3 = pdpiMember.TempatPraktek3
	}
	if pdpiMember.TempatPraktek3Tipe != nil {
		localMember.TempatPraktek3Tipe = pdpiMember.TempatPraktek3Tipe
	}
	if pdpiMember.TempatPraktek3Tipe2 != nil {
		localMember.TempatPraktek3Tipe2 = pdpiMember.TempatPraktek3Tipe2
	}
	if pdpiMember.TempatPraktek3Alkes != nil {
		localMember.TempatPraktek3Alkes = pdpiMember.TempatPraktek3Alkes
	}
	if pdpiMember.TempatPraktek3Alkes2 != nil {
		localMember.TempatPraktek3Alkes2 = pdpiMember.TempatPraktek3Alkes2
	}
	if pdpiMember.KotaKabupatenPraktek3 != nil {
		localMember.KotaKabupatenPraktek3 = pdpiMember.KotaKabupatenPraktek3
	}
	if pdpiMember.ProvinsiPraktek3 != nil {
		localMember.ProvinsiPraktek3 = pdpiMember.ProvinsiPraktek3
	}

	// STR SIP
	if pdpiMember.NoSTR != nil {
		localMember.NoSTR = pdpiMember.NoSTR
	}
	if pdpiMember.NoSIP != nil {
		localMember.NoSIP = pdpiMember.NoSIP
	}

	// Parse Dates
	if pdpiMember.TglLahir != nil && *pdpiMember.TglLahir != "" {
		if t, err := time.Parse("2006-01-02", *pdpiMember.TglLahir); err == nil {
			localMember.TglLahir = utils.TimeToPtr(t)
		}
	}
	if pdpiMember.STRBerlakuSampai != nil && *pdpiMember.STRBerlakuSampai != "" {
		if t, err := time.Parse("2006-01-02", *pdpiMember.STRBerlakuSampai); err == nil {
			localMember.STRBerlakuSampai = utils.TimeToPtr(t)
		}
	}
	if pdpiMember.SIPBerlakuSampai != nil && *pdpiMember.SIPBerlakuSampai != "" {
		if t, err := time.Parse("2006-01-02", *pdpiMember.SIPBerlakuSampai); err == nil {
			localMember.SIPBerlakuSampai = utils.TimeToPtr(t)
		}
	}

	return localMember
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

	// fetch existing from db
	err := pc.prepareLocalMemberCache()
	if err != nil {
		return 
	} 
	
	// Get user data
	user, err := models.FindByID(pc.db, userIDInt64)
	if err != nil || user == nil {
		utils.Error(c, http.StatusNotFound, "user_not_found", "User not found", nil)
		return
	}

	
	
	// Parse request (optional email or npa)
	var reqBody struct {
		Email string `json:"email"`
		NPA   string `json:"npa"`
	}
	_ = c.ShouldBindJSON(&reqBody)

	var pdpiMember *services.SupabaseMember
	var fetchErr error

	if reqBody.NPA != "" {
		pdpiMember, fetchErr = pc.pdpiService.FetchMemberByNPA(reqBody.NPA)
	} else {
		// Use user's email if not provided
		emailToSync := reqBody.Email
		if emailToSync == "" {
			emailToSync = user.Email
		}
		pdpiMember, fetchErr = pc.pdpiService.FetchMemberByEmail(emailToSync)
	}

	if fetchErr != nil {
		utils.Error(c, http.StatusNotFound, "member_not_found", "PDPI member not found in Supabase: "+fetchErr.Error(), nil)
		return
	}

	// Map and Save
	localMember := pc.mapSupabaseToLocal(*pdpiMember)

	// Link to current user if emails match
	if localMember.Email != nil && *localMember.Email == user.Email {
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

// SyncAllMembers syncs ALL PDPI members from Supabase API to local database
// POST /api/v1/pdpi/sync-all-members
// This is an admin-only operation
func (pc *PDPIController) SyncAllMembers(c *gin.Context) {
	startTime := time.Now()

	// Statistics tracking
	var (
		totalFetched  int
		totalSynced   int
		totalFailed   int
		errorMessages []string
	)

	// fetch existing from db
	err := pc.prepareLocalMemberCache()
	if err != nil {
		return 
	}
	// Fetch all members from Supabase via Service
	supabaseMembers, err := pc.pdpiService.FetchMembersFromSupabase()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "supabase_error", "Failed to fetch members from Supabase: "+err.Error(), nil)
		return
	}

	totalFetched = len(supabaseMembers)

	// Process each member
	for _, pdpiMember := range supabaseMembers {
		// Map to local model
		localMember := pc.mapSupabaseToLocal(pdpiMember)

		// Try to link to existing user by email
		if localMember.Email != nil && *localMember.Email != "" {
			user, _ := models.FindByEmail(pc.db, *localMember.Email)
			if user != nil {
				localMember.UserID = utils.Int64ToPtr(user.ID)
			}
		}

		// Upsert to database
		err = models.UpsertPDPIMember(pc.db, localMember)
		if err != nil {
			totalFailed++
			continue
		}

		totalSynced++
	}

	duration := time.Since(startTime)

	// Prepare response
	result := gin.H{
		"total_fetched": totalFetched,
		"total_synced":  totalSynced,
		"total_failed":  totalFailed,
		"duration_ms":   duration.Milliseconds(),
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
		utils.Error(c, http.StatusBadRequest, "deprecated", "Direct API access is deprecated on this version. Please sync members first then use source=local", nil)
		return
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

	// If not found locally, inform user to sync
	utils.Error(c, http.StatusNotFound, "member_not_found", "Member not found in local database. Please sync member data first.", nil)
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
	query := "SELECT id, npa, nama, gelar, gelar2, email, cabang, provinsi, kota_kabupaten, status, tempat_praktek_1, tempat_praktek_2, alumni, foto FROM pdpi_members WHERE 1=1"
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
		Foto           *string `db:"foto" json:"foto"`
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
		"COALESCE(alumni, '') as alumni, " +
		"COALESCE(foto, '') as foto " +
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

// ImportExcel imports member data from uploaded Excel file
// POST /api/v1/pdpi/import-excel
func (pc *PDPIController) ImportExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_file", "Failed to get uploaded file", nil)
		return
	}

	// Open the uploaded file
	f, err := file.Open()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "error", "Failed to open file", nil)
		return
	}
	defer f.Close()

	excel, err := excelize.OpenReader(f)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "error", "Failed to parse excel", nil)
		return
	}
	defer excel.Close()

	sheets := excel.GetSheetList()
	if len(sheets) == 0 {
		utils.Error(c, http.StatusBadRequest, "error", "No sheets found in excel", nil)
		return
	}

	rows, err := excel.GetRows(sheets[0])
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "error", "Failed to read rows", nil)
		return
	}

	if len(rows) < 2 {
		utils.Error(c, http.StatusBadRequest, "error", "Excel file is empty", nil)
		return
	}

	// Stats
	var updated, created, failed int
	
	// Column mapping based on inspection:
	// 0: NPA, 1: display_name, 2: user_email, 3: pdpi_cabang, 4: no WA
	
	for i, row := range rows {
		if i == 0 { // Skip header
			continue
		}
		
		if len(row) < 1 || row[0] == "" { // Skip empty NPA
			continue
		}

		npa := row[0]
		
		// Map data
		member := &models.PDPIMember{
			NPA: npa,
		}
		
		if len(row) > 1 { member.Nama = row[1] }
		if len(row) > 2 && row[2] != "" { member.Email = utils.StringToPtr(row[2]) }
		if len(row) > 3 && row[3] != "" { member.Cabang = utils.StringToPtr(row[3]) }
		if len(row) > 4 && row[4] != "" { 
			noHP := row[4]
			member.NoHP = utils.StringToPtr(utils.NormalizePhoneNumber(noHP)) 
		}
		
		// Set status to Aktif by default for imports
		activeStatus := "Aktif"
		member.Status = &activeStatus
		member.SyncedAt = utils.TimeToPtr(time.Now())

		// Check if exists to determine if it's update or create
		existing, _ := models.FindPDPIMemberByNPA(pc.db, npa)
		if existing != nil {
			member.ID = existing.ID
			updated++
		} else {
			created++
		}

		err := models.UpsertByNPA(pc.db, member)
		if err != nil {
			failed++
			log.Printf("Failed to upsert member %s: %v", npa, err)
		}
	}

	utils.Success(c, http.StatusOK, "Import completed", gin.H{
		"updated": updated,
		"created": created,
		"failed":  failed,
		"rows_processed": len(rows) - 1,
	})
}

type BirthdayMember struct {
	ID     string  `db:"id"`
	Nama   string  `db:"nama"`
	Email  *string `db:"email"`
	NoHP   *string `db:"no_hp"`
	Gelar  *string `db:"gelar"`
	Gelar2 *string `db:"gelar2"`
}

// TestSendBirthdayGreeting sends a birthday greeting test to a specific member by ID
// This is a debug/test endpoint - should only be accessible to admins
func (pc *PDPIController) TestSendBirthdayGreeting(c *gin.Context) {
	memberID := c.Param("id")
	if memberID == "" {
		utils.Error(c, http.StatusBadRequest, "validation_error", "Member ID is required", nil)
		return
	}

	member := BirthdayMember{}
	// Fetch member from database
	query := "SELECT id, nama, email, no_hp, gelar, gelar2 FROM pdpi_members WHERE id = ? LIMIT 1"
	// member := &struct {
	// 	ID     string
	// 	Nama   string
	// 	Email  *string
	// 	NoHP   *string
	// 	Gelar  *string
	// 	Gelar2 *string
	// }{}

	err := pc.db.Get(&member, query, memberID)
	fmt.Printf("RESULT: %+v\n", member)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "not_found", "Member not found", nil)
		return
	}

	// Initialize WhatsApp service
	waService := services.NewWhatsAppService(*pc.waba360)

	// Prepare data
	var phoneNumber string
	if member.NoHP != nil && *member.NoHP != "" {
		phoneNumber = *member.NoHP
	} else {
		utils.Error(c, http.StatusBadRequest, "validation_error", "Member has no phone number", nil)
		return
	}

	// Normalize phone number
	phoneNumber = utils.NormalizePhoneNumber(phoneNumber)

	// Prepare greeting data
	gelarStr := ""
	if member.Gelar != nil {
		gelarStr = *member.Gelar
		// Remove trailing dot if exists
		// gelarStr = strings.TrimSuffix(gelarStr, ".")
		gelarStr = strings.TrimRight(gelarStr, ". ")
	}

	gelar2Str := ""
	if member.Gelar2 != nil {
		gelar2Str = *member.Gelar2
	}

	// Send birthday greeting via WhatsApp using "ultah" template
	waErr := waService.SendUltah(phoneNumber, gelarStr, member.Nama, gelar2Str)
	if waErr != nil {
		utils.Error(c, http.StatusInternalServerError, "whatsapp_error", "Failed to send WhatsApp message: "+waErr.Error(), nil)
		return
	}

	if gelarStr == "" {
	gelarStr = "-"
}

if gelar2Str == "" {
	gelar2Str = "-"
}

	utils.Success(c, http.StatusOK, "Birthday greeting sent successfully", gin.H{
		"member_id": memberID,
		"nama":      member.Nama,
		"phone":     phoneNumber,
		"gelar":     gelarStr,
		"gelar2":    gelar2Str,
	})
}
