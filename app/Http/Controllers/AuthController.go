package controllers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	requests "github.com/cvudumbarainformatika/backend/app/Http/Requests"
	models "github.com/cvudumbarainformatika/backend/app/Models"
	services "github.com/cvudumbarainformatika/backend/app/Services"
	"github.com/cvudumbarainformatika/backend/config"
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// AuthController handles authentication-related requests
type AuthController struct {
	db     *sqlx.DB
	config *config.Config
}

// NewAuthController creates a new AuthController instance
func NewAuthController(db *sqlx.DB, cfg *config.Config) *AuthController {
	return &AuthController{
		db:     db,
		config: cfg,
	}
}

// Register handles user registration
// POST /api/v1/auth/register
func (ac *AuthController) Register(c *gin.Context) {
	var req requests.RegisterRequest

	// Validate request
	if err := req.Validate(c); err != nil {
		return
	}

	// Check if user already exists
	existingUser, err := models.FindByEmail(ac.db, req.Email)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to check user existence", nil)
		return
	}

	if existingUser != nil {
		utils.Error(c, http.StatusConflict, "user_exists", "User with this email already exists", nil)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "hashing_error", "Failed to hash password", nil)
		return
	}

	// Create new user
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "member",
		Status:   "pending",
	}

	if err := user.Create(ac.db); err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			utils.Error(c, http.StatusConflict, "user_exists", "User with this email already exists", nil)
			return
		}
		utils.Error(c, http.StatusInternalServerError, "registration_error", "Failed to register user", nil)
		return
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, ac.config.JWT.Secret, ac.config.JWT.AccessTokenExpiration)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "token_error", "Failed to generate access token", nil)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, ac.config.JWT.Secret, ac.config.JWT.RefreshTokenExpiration)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "token_error", "Failed to generate refresh token", nil)
		return
	}

	// Helper function to get string value
	getStringValue := func(ns sql.NullString) string {
		if ns.Valid {
			return ns.String
		}
		return ""
	}

	utils.Success(c, http.StatusCreated, "User registered successfully", gin.H{
		"user": gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"phone":      getStringValue(user.Phone),
			"address":    getStringValue(user.Address),
			"bio":        getStringValue(user.Bio),
			"avatar":     getStringValue(user.Avatar),
			"cabang":     getStringValue(user.Cabang),
			"role":       user.Role,
			"status":     user.Status,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    ac.config.JWT.AccessTokenExpiration * 60, // Convert minutes to seconds
	})
}

// Login handles user login with Hybrid Local-First Authentication Strategy
// POST /api/v1/auth/login
// Strategy:
//  1. Check if user exists in Local DB
//  2. If user exists:
//     a. Verify password (bcrypt)
//     b. If valid -> Login Success (Local Auth)
//     c. If invalid -> Try PDPI Auth (Hybrid Check for Password Change)
//     i. If PDPI success -> Update local password -> Login Success
//     ii. If PDPI fail -> Return Invalid Password
//  3. If user does NOT exist:
//     a. Try PDPI Auth
//     b. If PDPI success -> Auto-Create User -> Login Success
//     c. If PDPI fail -> Return Invalid Credentials
func (ac *AuthController) Login(c *gin.Context) {
	var req requests.LoginRequest

	// Validate request
	if err := req.Validate(c); err != nil {
		return
	}

	// Step 1: Check Local Database
	user, err := models.FindByEmail(ac.db, req.Email)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to retrieve user", nil)
		return
	}

	var pdpiResp *services.PDPILoginResponse
	var pdpiErr error
	pdpiService := services.NewPDPIService(&ac.config.PDPI)

	// Helper function for post-login actions (sync & token gen)
	handleSuccess := func(u *models.User, fromPDPI bool) {
		// Verify status
		if u.Status == "pending" {
			utils.Error(c, http.StatusUnauthorized, "user_pending", "Menunggu verifikasi Admin", nil)
			return
		}
		if u.Status != "active" {
			utils.Error(c, http.StatusUnauthorized, "user_inactive", "Harap mendaftar terlebih dahulu!", nil)
			return
		}

		// Generate Tokens
		accessToken, err := utils.GenerateAccessToken(u.ID, u.Email, ac.config.JWT.Secret, ac.config.JWT.AccessTokenExpiration)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "token_error", "Failed to generate access token", nil)
			return
		}
		refreshToken, err := utils.GenerateRefreshToken(u.ID, u.Email, ac.config.JWT.Secret, ac.config.JWT.RefreshTokenExpiration)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "token_error", "Failed to generate refresh token", nil)
			return
		}

		// Helper string value using local scope since the original closure is gone
		getStringValue := func(ns sql.NullString) string {
			if ns.Valid {
				return ns.String
			}
			return ""
		}

		utils.Success(c, http.StatusOK, "Login successful", gin.H{
			"user": gin.H{
				"id":         u.ID,
				"name":       u.Name,
				"email":      u.Email,
				"phone":      getStringValue(u.Phone),
				"address":    getStringValue(u.Address),
				"bio":        getStringValue(u.Bio),
				"avatar":     getStringValue(u.Avatar),
				"cabang":     getStringValue(u.Cabang),
				"role":       u.Role,
				"status":     u.Status,
				"created_at": u.CreatedAt,
				"updated_at": u.UpdatedAt,
			},
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expires_in":    ac.config.JWT.AccessTokenExpiration * 60,
		})

		// Background: Sync PDPI member data if logged in via PDPI or periodically
		// For now, if we have PDPI response, we sync the member data to local DB
		if fromPDPI && pdpiResp != nil && pdpiResp.Success && pdpiResp.Member.NPA != "" {
			// Sync logic (simplified from previous implementation)
			localMember := &models.PDPIMember{
				ID:             pdpiResp.Member.ID,
				NPA:            pdpiResp.Member.NPA,
				Nama:           pdpiResp.Member.Nama,
				Gelar:          sql.NullString{String: pdpiResp.Member.Gelar, Valid: pdpiResp.Member.Gelar != ""},
				Gelar2:         sql.NullString{String: pdpiResp.Member.Gelar2, Valid: pdpiResp.Member.Gelar2 != ""},
				Email:          sql.NullString{String: pdpiResp.Member.Email, Valid: pdpiResp.Member.Email != ""},
				NoHP:           sql.NullString{String: pdpiResp.Member.NoHP, Valid: pdpiResp.Member.NoHP != ""},
				NIK:            sql.NullString{String: pdpiResp.Member.NIK, Valid: pdpiResp.Member.NIK != ""},
				JenisKelamin:   sql.NullString{String: pdpiResp.Member.JenisKelamin, Valid: pdpiResp.Member.JenisKelamin != ""},
				TempatLahir:    sql.NullString{String: pdpiResp.Member.TempatLahir, Valid: pdpiResp.Member.TempatLahir != ""},
				AlamatRumah:    sql.NullString{String: pdpiResp.Member.AlamatRumah, Valid: pdpiResp.Member.AlamatRumah != ""},
				Cabang:         sql.NullString{String: pdpiResp.Member.Cabang, Valid: pdpiResp.Member.Cabang != ""},
				Provinsi:       sql.NullString{String: pdpiResp.Member.Provinsi, Valid: pdpiResp.Member.Provinsi != ""},
				KotaKabupaten:  sql.NullString{String: pdpiResp.Member.KotaKabupaten, Valid: pdpiResp.Member.KotaKabupaten != ""},
				Status:         sql.NullString{String: pdpiResp.Member.Status, Valid: pdpiResp.Member.Status != ""},
				Alumni:         sql.NullString{String: pdpiResp.Member.Alumni, Valid: pdpiResp.Member.Alumni != ""},
				ThnLulus:       sql.NullInt64{Int64: int64(pdpiResp.Member.ThnLulus), Valid: pdpiResp.Member.ThnLulus > 0},
				TempatTugas:    sql.NullString{String: pdpiResp.Member.TempatTugas, Valid: pdpiResp.Member.TempatTugas != ""},
				TempatPraktek1: sql.NullString{String: pdpiResp.Member.TempatPraktek1, Valid: pdpiResp.Member.TempatPraktek1 != ""},
				TempatPraktek2: sql.NullString{String: pdpiResp.Member.TempatPraktek2, Valid: pdpiResp.Member.TempatPraktek2 != ""},
				Subspesialis:   sql.NullString{String: pdpiResp.Member.Subspesialis, Valid: pdpiResp.Member.Subspesialis != ""},
				NoSTR:          sql.NullString{String: pdpiResp.Member.NoSTR, Valid: pdpiResp.Member.NoSTR != ""},
				NoSIP:          sql.NullString{String: pdpiResp.Member.NoSIP, Valid: pdpiResp.Member.NoSIP != ""},
				UserID:         sql.NullInt64{Int64: u.ID, Valid: true},
				SyncedAt:       sql.NullTime{Time: time.Now(), Valid: true},
			}

			// Parse dates
			if t, err := time.Parse("2006-01-02", pdpiResp.Member.TglLahir); err == nil {
				localMember.TglLahir = sql.NullTime{Time: t, Valid: true}
			}
			if t, err := time.Parse("2006-01-02", pdpiResp.Member.STRBerlakuSampai); err == nil {
				localMember.STRBerlakuSampai = sql.NullTime{Time: t, Valid: true}
			}
			if t, err := time.Parse("2006-01-02", pdpiResp.Member.SIPBerlakuSampai); err == nil {
				localMember.SIPBerlakuSampai = sql.NullTime{Time: t, Valid: true}
			}

			// Upsert to DB
			_ = models.UpsertPDPIMember(ac.db, localMember)
		}
	}

	if user != nil {
		// Condition 1: User Exists in Local DB

		// 1a. Verify Password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err == nil {
			// Password Match! Login Success (Local Auth)
			handleSuccess(user, false) // Not explicit from PDPI, using local data
			return
		}

		// Password Invalid: Try PDPI Auth (Hybrid Check)
		// Maybe user changed password in PDPI portal
		pdpiResp, pdpiErr = pdpiService.Login(req.Email, req.Password)
		if pdpiErr == nil && pdpiResp != nil && pdpiResp.Success {
			// PDPI Success! User changed password.
			// Update local password
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			user.Password = string(hashedPassword)
			_ = user.Update(ac.db) // Update password in DB

			handleSuccess(user, true) // From PDPI
			return
		}

		// Both Local and PDPI failed
		utils.Error(c, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password", nil)
		return

	} else {
		// Condition 2: User NOT Found in Local DB
		// Try PDPI Auth
		pdpiResp, pdpiErr = pdpiService.Login(req.Email, req.Password)

		if pdpiErr == nil && pdpiResp != nil && pdpiResp.Success {
			// PDPI Success! Auto-Create User
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				utils.Error(c, http.StatusInternalServerError, "hashing_error", "Failed to hash password", nil)
				return
			}

			user = &models.User{
				Name:     pdpiResp.Member.Nama,
				Email:    req.Email,
				Password: string(hashedPassword),
				Role:     "member", // Default role
				Status:   "active", // Auto-active
			}

			if pdpiResp.Member.NoHP != "" {
				user.Phone = sql.NullString{String: pdpiResp.Member.NoHP, Valid: true}
			}
			if pdpiResp.Member.Cabang != "" {
				user.Cabang = sql.NullString{String: pdpiResp.Member.Cabang, Valid: true}
			}

			if err := user.Create(ac.db); err != nil {
				utils.Error(c, http.StatusInternalServerError, "registration_error", "Failed to create user from PDPI data", nil)
				return
			}

			handleSuccess(user, true) // From PDPI
			return
		}

		// PDPI also failed
		utils.Error(c, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password", nil)
		return
	}
}

// Me retrieves the current authenticated user's information
// GET /api/v1/auth/me
func (ac *AuthController) Me(c *gin.Context) {
	// Get user ID from JWT middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "unauthorized", "User not authenticated", nil)
		return
	}

	// Type assert to int64
	userIDInt64, ok := userID.(int64)
	if !ok {
		utils.Error(c, http.StatusInternalServerError, "internal_error", "Invalid user ID format", nil)
		return
	}

	// Fetch user from database
	user, err := models.FindByID(ac.db, userIDInt64)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to retrieve user", nil)
		return
	}

	if user == nil {
		utils.Error(c, http.StatusNotFound, "user_not_found", "User not found", nil)
		return
	}

	// Helper function to get string value
	getStringValue := func(ns sql.NullString) string {
		if ns.Valid {
			return ns.String
		}
		return ""
	}

	utils.Success(c, http.StatusOK, "User information retrieved successfully", gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"phone":      getStringValue(user.Phone),
		"address":    getStringValue(user.Address),
		"bio":        getStringValue(user.Bio),
		"avatar":     getStringValue(user.Avatar),
		"cabang":     getStringValue(user.Cabang),
		"role":       user.Role,
		"status":     user.Status,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// Refresh generates a new access token using a valid refresh token
// POST /api/v1/auth/refresh
func (ac *AuthController) Refresh(c *gin.Context) {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	// Validate refresh token
	claims, err := utils.ValidateToken(req.RefreshToken, ac.config.JWT.Secret)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "invalid_token", "Invalid or expired refresh token", nil)
		return
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(claims.UserID, claims.Email, ac.config.JWT.Secret, ac.config.JWT.AccessTokenExpiration)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "token_error", "Failed to generate access token", nil)
		return
	}

	utils.Success(c, http.StatusOK, "Access token refreshed successfully", gin.H{
		"access_token": accessToken,
		"expires_in":   ac.config.JWT.AccessTokenExpiration * 60, // Convert minutes to seconds
	})
}

// Logout handles user logout
// POST /api/v1/auth/logout
func (ac *AuthController) Logout(c *gin.Context) {
	// In a simple JWT implementation, logout is typically handled client-side
	// by removing the token from storage. However, this endpoint can be used
	// to invalidate tokens server-side by storing them in a blacklist (e.g., Redis)

	utils.Success(c, http.StatusOK, "Logout successful", gin.H{})
}

// ChangePassword changes user password
// POST /api/v1/auth/profile/change-password
func (ac *AuthController) ChangePassword(c *gin.Context) {
	var req requests.ChangePasswordRequest

	// Validate request
	if err := req.Validate(c); err != nil {
		return
	}

	// Get user ID from JWT middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "unauthorized", "User not authenticated", nil)
		return
	}

	// Type assert to int64
	userIDInt64, ok := userID.(int64)
	if !ok {
		utils.Error(c, http.StatusInternalServerError, "internal_error", "Invalid user ID format", nil)
		return
	}

	// Fetch user from database
	user, err := models.FindByID(ac.db, userIDInt64)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to retrieve user", nil)
		return
	}

	if user == nil {
		utils.Error(c, http.StatusNotFound, "user_not_found", "User not found", nil)
		return
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		utils.Error(c, http.StatusUnauthorized, "invalid_password", "Current password is incorrect", nil)
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "hash_error", "Failed to hash password", nil)
		return
	}

	// Update user password
	user.Password = string(hashedPassword)
	if err := user.Update(ac.db); err != nil {
		utils.Error(c, http.StatusInternalServerError, "update_error", "Failed to update password", nil)
		return
	}

	utils.Success(c, http.StatusOK, "Password changed successfully", gin.H{
		"id":      user.ID,
		"email":   user.Email,
		"message": "Your password has been changed successfully",
	})
}

// UpdateProfile updates the current authenticated user's profile
// PUT /api/v1/auth/profile
// Supports both JSON and multipart/form-data requests
func (ac *AuthController) UpdateProfile(c *gin.Context) {
	var req requests.UpdateProfileRequest

	// Validate request
	if err := req.Validate(c); err != nil {
		return
	}

	// Get user ID from JWT middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "unauthorized", "User not authenticated", nil)
		return
	}

	// Type assert to int64
	userIDInt64, ok := userID.(int64)
	if !ok {
		utils.Error(c, http.StatusInternalServerError, "internal_error", "Invalid user ID format", nil)
		return
	}

	// Fetch user from database
	user, err := models.FindByID(ac.db, userIDInt64)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to retrieve user", nil)
		return
	}

	if user == nil {
		utils.Error(c, http.StatusNotFound, "user_not_found", "User not found", nil)
		return
	}

	// Handle avatar: file upload takes precedence over JSON string
	if req.Avatar != nil {
		// Delete old avatar if exists
		if user.Avatar.Valid && user.Avatar.String != "" {
			utils.DeleteAvatar(user.Avatar.String)
		}

		// Upload new avatar
		avatarPath, err := utils.UploadAvatar(req.Avatar, userIDInt64)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "upload_error", "Failed to upload avatar: "+err.Error(), nil)
			return
		}
		user.Avatar.String = avatarPath
		user.Avatar.Valid = true
	}
	// If no file upload, avatar field from JSON stays as is (can be empty string or existing URL)

	// Update user profile fields
	user.Name = req.Name
	user.Phone.String = req.Phone
	user.Phone.Valid = req.Phone != ""
	user.Address.String = req.Address
	user.Address.Valid = req.Address != ""
	user.Bio.String = req.Bio
	user.Bio.Valid = req.Bio != ""
	user.Cabang.String = req.Cabang
	user.Cabang.Valid = req.Cabang != ""

	// Save updated user
	if err := user.UpdateProfile(ac.db); err != nil {
		utils.Error(c, http.StatusInternalServerError, "update_error", "Failed to update profile", nil)
		return
	}

	// Helper function to get string value
	getStringValue := func(ns sql.NullString) string {
		if ns.Valid {
			return ns.String
		}
		return ""
	}

	utils.Success(c, http.StatusOK, "Profile updated successfully", gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"phone":      getStringValue(user.Phone),
		"address":    getStringValue(user.Address),
		"bio":        getStringValue(user.Bio),
		"avatar":     getStringValue(user.Avatar), // Already contains API endpoint URL
		"cabang":     getStringValue(user.Cabang),
		"role":       user.Role,
		"status":     user.Status,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}
