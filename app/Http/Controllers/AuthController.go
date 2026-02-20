package controllers

import (
	"database/sql"
	"net/http"
	"strings"

	requests "github.com/cvudumbarainformatika/backend/app/Http/Requests"
	models "github.com/cvudumbarainformatika/backend/app/Models"
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

	// Helper function for post-login actions (token gen only)
	handleSuccess := func(u *models.User) {
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
	}

	if user != nil {
		// Condition 1: User Exists in Local DB
		// Verify Password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err == nil {
			// Password Match! Login Success (Local Auth)
			handleSuccess(user)
			return
		}
	}

	// If user not found or password invalid, return error
	// PDPI Hybrid Auth is DISABLED due to migration to Supabase
	utils.Error(c, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password", nil)
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
