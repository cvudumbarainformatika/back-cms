package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

// AvatarController handles avatar file operations
type AvatarController struct{}

// NewAvatarController creates a new avatar controller instance
func NewAvatarController() *AvatarController {
	return &AvatarController{}
}

// GetAvatar retrieves and serves an avatar file
// GET /api/v1/avatars/:user_id
// Optional: can be public or protected based on config
func (ac *AvatarController) GetAvatar(c *gin.Context) {
	userID := c.Param("user_id")

	if userID == "" {
		utils.Error(c, http.StatusBadRequest, "invalid_request", "User ID is required", nil)
		return
	}

	// Build safe file path - prevent directory traversal attacks
	// Avatar files follow naming convention: avatar_{user_id}_*.ext
	avatarDir := utils.GetStoragePath()

	// List files matching the user's avatar pattern
	files, err := os.ReadDir(avatarDir)
	if err != nil {
		if os.IsNotExist(err) {
			utils.Error(c, http.StatusNotFound, "not_found", "Avatar not found", nil)
			return
		}
		utils.Error(c, http.StatusInternalServerError, "read_error", "Failed to read avatar", nil)
		return
	}

	// Find the latest avatar file for this user
	var latestFile *os.DirEntry
	prefix := "avatar_" + userID + "_"

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), prefix) {
			if latestFile == nil {
				f := file
				latestFile = &f
			}
		}
	}

	if latestFile == nil {
		utils.Error(c, http.StatusNotFound, "not_found", "Avatar not found for this user", nil)
		return
	}

	// Serve the file
	filePath := filepath.Join(avatarDir, (*latestFile).Name())

	// Add security headers
	c.Header("Cache-Control", "public, max-age=86400") // Cache for 24 hours
	c.Header("X-Content-Type-Options", "nosniff")      // Prevent MIME type sniffing

	// Serve file
	c.File(filePath)
}

// GetAvatarByName retrieves avatar by specific filename
// GET /api/v1/avatars/file/:filename
// More restrictive - only serves files matching avatar pattern
func (ac *AvatarController) GetAvatarByName(c *gin.Context) {
	filename := c.Param("filename")

	if filename == "" {
		utils.Error(c, http.StatusBadRequest, "invalid_request", "Filename is required", nil)
		return
	}

	// Validate filename - must match avatar naming pattern
	// Pattern: avatar_{user_id}_{timestamp}.{ext}
	if !isValidAvatarFilename(filename) {
		utils.Error(c, http.StatusForbidden, "invalid_file", "Invalid avatar filename", nil)
		return
	}

	// Build safe path using configured storage path
	avatarDir := utils.GetStoragePath()
	filePath := filepath.Join(avatarDir, filename)

	// Ensure path is within avatars directory (prevent directory traversal)
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_path", "Invalid file path", nil)
		return
	}

	absAvatarDir, _ := filepath.Abs("storage/avatars")
	if !strings.HasPrefix(absPath, absAvatarDir) {
		utils.Error(c, http.StatusForbidden, "invalid_path", "Access denied", nil)
		return
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.Error(c, http.StatusNotFound, "not_found", "Avatar file not found", nil)
		return
	}

	// Add security headers
	c.Header("Cache-Control", "public, max-age=86400")
	c.Header("X-Content-Type-Options", "nosniff")

	// Serve file
	c.File(filePath)
}

// isValidAvatarFilename validates avatar filename format
func isValidAvatarFilename(filename string) bool {
	// Must start with "avatar_"
	if !strings.HasPrefix(filename, "avatar_") {
		return false
	}

	// Must have extension
	if !strings.Contains(filename, ".") {
		return false
	}

	// Only allow specific extensions
	ext := strings.ToLower(filepath.Ext(filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	return allowedExts[ext]
}
