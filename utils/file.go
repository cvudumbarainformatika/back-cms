package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
	"log"
)

// AllowedAvatarMimes allowed MIME types for avatar upload
var AllowedAvatarMimes = []string{"image/jpeg", "image/png", "image/gif", "image/webp"}

// MaxAvatarSize maximum file size for avatar in bytes (5MB)
const MaxAvatarSize = 5 * 1024 * 1024

// GetStoragePath returns the configured storage path for avatars
func GetStoragePath() string {
	storagePath := os.Getenv("AVATAR_STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage/avatars"  // fallback untuk development
	}
	return storagePath
}

// UploadAvatar uploads an avatar file and returns the relative URL path
func UploadAvatar(file *multipart.FileHeader, userID int64) (string, error) {
	// Validate file size
	if file.Size > MaxAvatarSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of 5MB")
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Validate MIME type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	mimeType := file.Header.Get("Content-Type")
	if !isAllowedMimeType(mimeType) {
		return "", fmt.Errorf("file type not allowed. Allowed types: jpeg, png, gif, webp")
	}

	// Reset file pointer after reading for validation
	src.Seek(0, 0)

	// Get storage path from environment or use default
	avatarDir := GetStoragePath()
	log.Printf("Using avatar storage path: %s", avatarDir)
	
	// Create avatars directory if it doesn't exist
	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create avatars directory: %v", err)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("avatar_%d_%d%s", userID, timestamp, ext)
	filePath := filepath.Join(avatarDir, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// Return relative URL path (use configured storage path prefix)
	// This should match the API endpoint path
	relativePath := strings.ReplaceAll(filePath, "\\", "/")
	// If it's an absolute path, convert to relative for API response
	if strings.HasPrefix(relativePath, "/app") {
		relativePath = strings.TrimPrefix(relativePath, "/app")
	}
	if !strings.HasPrefix(relativePath, "/") {
		relativePath = "/" + relativePath
	}
	
	log.Printf("Avatar uploaded successfully: %s", relativePath)
	return relativePath, nil
}

// DeleteAvatar deletes an avatar file
func DeleteAvatar(avatarPath string) error {
	if avatarPath == "" {
		return nil
	}

	// Build full path from storage directory
	storagePath := GetStoragePath()
	
	// If path is relative, join with storage path
	fullPath := avatarPath
	if !strings.HasPrefix(avatarPath, "/") && !strings.HasPrefix(avatarPath, storagePath) {
		fullPath = filepath.Join(storagePath, avatarPath)
	} else if strings.HasPrefix(avatarPath, "/storage") {
		// Convert API path to actual storage path
		// /storage/avatars/avatar_1.jpg → {STORAGE_PATH}/avatar_1.jpg
		filename := filepath.Base(avatarPath)
		fullPath = filepath.Join(storagePath, filename)
	}

	log.Printf("Attempting to delete avatar: %s (from storage: %s)", fullPath, storagePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Printf("Avatar file does not exist: %s", fullPath)
		return nil // File doesn't exist, no error
	}

	// Delete file
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete avatar file: %v", err)
	}

	log.Printf("Avatar deleted successfully: %s", fullPath)
	return nil
}

// isAllowedMimeType checks if MIME type is allowed
func isAllowedMimeType(mimeType string) bool {
	for _, allowed := range AllowedAvatarMimes {
		if mimeType == allowed {
			return true
		}
	}
	return false
}
