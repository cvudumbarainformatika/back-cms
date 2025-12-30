package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	requests "github.com/cvudumbarainformatika/backend/app/Http/Requests"
	models "github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// UserController handles user management operations
type UserController struct {
	db *sqlx.DB
}

// NewUserController creates a new UserController instance
func NewUserController(db *sqlx.DB) *UserController {
	return &UserController{
		db: db,
	}
}

// GetList returns paginated list of users with optional filters
// GET /api/v1/users/get-lists?page=&per_page=&q=&role=&status=&cabang=&sort=&order=
func (uc *UserController) GetList(c *gin.Context) {
	// Get pagination parameters
	page, limit := utils.GetPaginationParams(c)
	
	// Get filter parameters from frontend
	search := c.Query("q")
	role := c.Query("role")
	status := c.Query("status")
	cabang := c.Query("cabang")
	
	// Get sort parameters (frontend sends sort=column, order=direction)
	orderBy := c.DefaultQuery("sort", "created_at")
	sortOrder := c.DefaultQuery("order", "desc")
	
	// Validate orderBy to prevent SQL injection (whitelist allowed columns)
	allowedColumns := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"role":       true,
		"status":     true,
		"cabang":     true,
		"phone":      true,
		"address":    true,
		"bio":        true,
		"avatar":     true,
		"created_at": true,
		"updated_at": true,
	}
	if !allowedColumns[orderBy] {
		orderBy = "created_at"
	}
	
	// Validate and normalize sort order
	sortOrder = strings.ToUpper(sortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	// Build query
	query := `SELECT id, name, email, role, status, cabang, phone, address, bio, avatar, created_at, updated_at FROM users WHERE 1=1`
	args := []interface{}{}

	// Add filters
	if role != "" {
		query += ` AND role = ?`
		args = append(args, role)
	}

	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}

	if cabang != "" {
		query += ` AND cabang = ?`
		args = append(args, cabang)
	}

	if search != "" {
		query += ` AND (name LIKE ? OR email LIKE ?)`
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// Count total records
	countQuery := `SELECT COUNT(*) as total FROM users WHERE 1=1`
	countArgs := []interface{}{}
	
	if role != "" {
		countQuery += ` AND role = ?`
		countArgs = append(countArgs, role)
	}
	if status != "" {
		countQuery += ` AND status = ?`
		countArgs = append(countArgs, status)
	}
	if cabang != "" {
		countQuery += ` AND cabang = ?`
		countArgs = append(countArgs, cabang)
	}
	if search != "" {
		countQuery += ` AND (name LIKE ? OR email LIKE ?)`
		searchPattern := "%" + search + "%"
		countArgs = append(countArgs, searchPattern, searchPattern)
	}

	var total int64
	err := uc.db.Get(&total, countQuery, countArgs...)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to count users", nil)
		return
	}

	// Add ordering and pagination
	query += ` ORDER BY ` + orderBy + ` ` + sortOrder + ` LIMIT ? OFFSET ?`
	offset := (page - 1) * limit
	args = append(args, limit, offset)

	// Fetch users
	var users []models.User
	err = uc.db.Select(&users, query, args...)
	if err != nil {
		// Log the actual error for debugging
		fmt.Printf("Query: %s\n", query)
		fmt.Printf("Args: %v\n", args)
		fmt.Printf("Error: %v\n", err)
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch users: "+err.Error(), nil)
		return
	}

	// Format response
	userResponses := make([]gin.H, len(users))
	for i, user := range users {
		userResponses[i] = gin.H{
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
		}
	}

	// Use standard pagination response format
	pagination := utils.OffsetPaginate(userResponses, page, limit, total)
	
	utils.Success(c, http.StatusOK, "Users fetched successfully", gin.H{
		"items":      pagination.Data,
		"pagination": pagination.Meta,
	})
}

// GetByID returns a single user by ID
// GET /api/v1/users/:id
func (uc *UserController) GetByID(c *gin.Context) {
	id := c.Param("id")

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid user ID", nil)
		return
	}

	user, err := models.FindByID(uc.db, userID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "user_not_found", "User not found", nil)
		return
	}

	utils.Success(c, http.StatusOK, "User retrieved successfully", gin.H{
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

// Create creates a new user
// POST /api/v1/users
func (uc *UserController) Create(c *gin.Context) {
	var req requests.CreateUserRequest

	if err := req.Validate(c); err != nil {
		return
	}

	// Check if email already exists
	existingUser, _ := models.FindByEmail(uc.db, req.Email)
	if existingUser != nil {
		utils.Error(c, http.StatusConflict, "email_exists", "Email already registered", nil)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "hash_error", "Failed to hash password", nil)
		return
	}

	// Create new user
	user := &models.User{
		Name:   req.Name,
		Email:  req.Email,
		Password: string(hashedPassword),
		Role:   req.Role,
		Status: req.Status,
	}

	// Set optional fields
	if req.Phone != "" {
		user.Phone.String = req.Phone
		user.Phone.Valid = true
	}
	if req.Address != "" {
		user.Address.String = req.Address
		user.Address.Valid = true
	}
	if req.Bio != "" {
		user.Bio.String = req.Bio
		user.Bio.Valid = true
	}
	if req.Cabang != "" {
		user.Cabang.String = req.Cabang
		user.Cabang.Valid = true
	}

	// Save to database
	err = user.Create(uc.db)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to create user", nil)
		return
	}

	utils.Success(c, http.StatusCreated, "User created successfully", gin.H{
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

// Update updates a user's information
// PUT /api/v1/users/:id
func (uc *UserController) Update(c *gin.Context) {
	id := c.Param("id")

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid user ID", nil)
		return
	}

	var req requests.UpdateUserRequest
	if err := req.Validate(c); err != nil {
		return
	}

	// Get existing user
	user, err := models.FindByID(uc.db, userID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "user_not_found", "User not found", nil)
		return
	}

	// Check if new email already exists (if email is being changed)
	if req.Email != user.Email {
		existingUser, _ := models.FindByEmail(uc.db, req.Email)
		if existingUser != nil {
			utils.Error(c, http.StatusConflict, "email_exists", "Email already registered", nil)
			return
		}
	}

	// Update fields
	user.Name = req.Name
	user.Email = req.Email
	user.Role = req.Role
	user.Status = req.Status

	// Update optional fields
	if req.Phone != "" {
		user.Phone.String = req.Phone
		user.Phone.Valid = true
	} else {
		user.Phone.Valid = false
	}

	if req.Address != "" {
		user.Address.String = req.Address
		user.Address.Valid = true
	} else {
		user.Address.Valid = false
	}

	if req.Bio != "" {
		user.Bio.String = req.Bio
		user.Bio.Valid = true
	} else {
		user.Bio.Valid = false
	}

	if req.Cabang != "" {
		user.Cabang.String = req.Cabang
		user.Cabang.Valid = true
	} else {
		user.Cabang.Valid = false
	}

	// Update password if provided
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "hash_error", "Failed to hash password", nil)
			return
		}
		user.Password = string(hashedPassword)
	}

	// Save to database
	err = user.Update(uc.db)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to update user", nil)
		return
	}

	utils.Success(c, http.StatusOK, "User updated successfully", gin.H{
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

// Patch updates only role and/or status of a user
// PATCH /api/v1/users/:id
func (uc *UserController) Patch(c *gin.Context) {
	id := c.Param("id")

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid user ID", nil)
		return
	}

	var req requests.PatchUserRequest
	if err := req.Validate(c); err != nil {
		return
	}

	// Get existing user
	user, err := models.FindByID(uc.db, userID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "user_not_found", "User not found", nil)
		return
	}

	// Update only role and status if provided
	if req.Role != "" {
		user.Role = req.Role
	}

	if req.Status != "" {
		user.Status = req.Status
	}

	// Save to database
	err = user.Update(uc.db)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to update user", nil)
		return
	}

	utils.Success(c, http.StatusOK, "User updated successfully", gin.H{
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

// Delete performs a soft delete on a user
// DELETE /api/v1/users/:id
func (uc *UserController) Delete(c *gin.Context) {
	id := c.Param("id")

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid user ID", nil)
		return
	}

	// Get existing user
	user, err := models.FindByID(uc.db, userID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "user_not_found", "User not found", nil)
		return
	}

	// Soft delete by setting status to deleted
	user.Status = "deleted"
	err = user.Update(uc.db)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to delete user", nil)
		return
	}

	utils.Success(c, http.StatusOK, "User deleted successfully", nil)
}

// Helper function to get string value from sql.NullString
func getStringValue(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
