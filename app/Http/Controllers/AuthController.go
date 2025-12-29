package controllers

import (
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
		Role:     "user",
		Status:   "active",
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

	utils.Success(c, http.StatusCreated, "User registered successfully", gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    ac.config.JWT.AccessTokenExpiration * 60, // Convert minutes to seconds
	})
}

// Login handles user login
// POST /api/v1/auth/login
func (ac *AuthController) Login(c *gin.Context) {
	var req requests.LoginRequest

	// Validate request
	if err := req.Validate(c); err != nil {
		return
	}

	// Find user by email
	user, err := models.FindByEmail(ac.db, req.Email)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to retrieve user", nil)
		return
	}

	if user == nil {
		utils.Error(c, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password", nil)
		return
	}

	// Check if user is active
	if user.Status != "active" {
		utils.Error(c, http.StatusUnauthorized, "user_inactive", "User account is not active", nil)
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password", nil)
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

	utils.Success(c, http.StatusOK, "Login successful", gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    ac.config.JWT.AccessTokenExpiration * 60, // Convert minutes to seconds
	})
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

	utils.Success(c, http.StatusOK, "User information retrieved successfully", gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
		"status": user.Status,
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
