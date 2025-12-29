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
		// Try generic storage path first
		if genPath := os.Getenv("STORAGE_BASE_PATH"); genPath != "" {
			storagePath = filepath.Join(genPath, "avatars")
		} else {
			storagePath = "./storage/avatars"  // fallback untuk development
		}
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

	// Read ALL file content into memory first
	fileContent, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	if len(fileContent) == 0 {
		return "", fmt.Errorf("uploaded file is empty")
	}

	log.Printf("[Avatar] File content read: %d bytes", len(fileContent))

	// Validate MIME type from Content-Type header
	mimeType := file.Header.Get("Content-Type")
	if !isAllowedMimeType(mimeType) {
		return "", fmt.Errorf("file type not allowed. Allowed types: jpeg, png, gif, webp")
	}

	// Get storage path from environment with better fallback
	avatarDir := getAvatarStoragePath()
	log.Printf("[Avatar] Using storage path: %s", avatarDir)
	
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
		log.Printf("[Avatar] ERROR creating file %s: %v", filePath, err)
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	log.Printf("[Avatar] File created successfully: %s", filePath)

	// Write file content to disk
	bytesCopied, err := dst.Write(fileContent)
	if err != nil {
		log.Printf("[Avatar] ERROR writing file content: %v (bytes written: %d)", err, bytesCopied)
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	log.Printf("[Avatar] File content written successfully: %d bytes", bytesCopied)

	// Flush to disk
	if err := dst.Sync(); err != nil {
		log.Printf("[Avatar] WARNING: Failed to sync file: %v", err)
	}

	// Verify file exists and has content
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("[Avatar] ERROR: File not found after creation: %v", err)
		return "", fmt.Errorf("file verification failed: %v", err)
	}

	log.Printf("[Avatar] File verified: %s (size: %d bytes)", filePath, fileInfo.Size())

	if fileInfo.Size() == 0 {
		log.Printf("[Avatar] ERROR: File is empty! Size: 0 bytes")
		return "", fmt.Errorf("uploaded file is empty")
	}

	// Return API endpoint URL (not storage path)
	// Format: /api/v1/files/avatar/{filename}
	apiURL := fmt.Sprintf("/api/v1/files/avatar/%s", filename)
	
	log.Printf("[Avatar] SUCCESS: Avatar uploaded: %s (stored at: %s, size: %d bytes)", apiURL, filePath, fileInfo.Size())
	return apiURL, nil
}

// DeleteAvatar deletes an avatar file
func DeleteAvatar(avatarPath string) error {
	if avatarPath == "" {
		return nil
	}

	// Build full path from storage directory
	storagePath := GetStoragePath()
	
	// Extract filename from various path formats
	var fullPath string
	
	// Format 1: API URL /api/v1/files/avatar/avatar_1_xxx.jpg
	if strings.HasPrefix(avatarPath, "/api/v1/files/avatar/") {
		filename := filepath.Base(avatarPath)
		fullPath = filepath.Join(storagePath, filename)
		log.Printf("[Avatar] Converting API URL to storage path: %s → %s", avatarPath, fullPath)
	} else if strings.HasPrefix(avatarPath, "/storage/avatars/") {
		// Format 2: Old storage path /storage/avatars/avatar_1_xxx.jpg
		filename := filepath.Base(avatarPath)
		fullPath = filepath.Join(storagePath, filename)
		log.Printf("[Avatar] Converting storage path to full path: %s → %s", avatarPath, fullPath)
	} else if !strings.ContainsAny(avatarPath, "/\\") {
		// Format 3: Just filename avatar_1_xxx.jpg
		fullPath = filepath.Join(storagePath, avatarPath)
		log.Printf("[Avatar] Using filename with storage path: %s", fullPath)
	} else {
		// Format 4: Absolute or relative path
		fullPath = avatarPath
		log.Printf("[Avatar] Using as-is path: %s", fullPath)
	}

	log.Printf("[Avatar] Attempting to delete: %s", fullPath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Printf("[Avatar] File does not exist (OK): %s", fullPath)
		return nil // File doesn't exist, no error
	}

	// Delete file
	if err := os.Remove(fullPath); err != nil {
		log.Printf("[Avatar] Error deleting file: %v", err)
		return fmt.Errorf("failed to delete avatar file: %v", err)
	}

	log.Printf("[Avatar] Deleted successfully: %s", fullPath)
	return nil
}

// getAvatarStoragePath returns the avatar storage path with detailed logging
func getAvatarStoragePath() string {
	// Check specific avatar path first
	if avatarPath := os.Getenv("AVATAR_STORAGE_PATH"); avatarPath != "" {
		log.Printf("[Avatar] Found AVATAR_STORAGE_PATH env var: %s", avatarPath)
		return avatarPath
	}

	// Check generic storage path
	if basePath := os.Getenv("STORAGE_BASE_PATH"); basePath != "" {
		path := filepath.Join(basePath, "avatars")
		log.Printf("[Avatar] Using STORAGE_BASE_PATH with /avatars: %s", path)
		return path
	}

	// Default fallback
	log.Printf("[Avatar] Using default fallback path: ./storage/avatars")
	return "./storage/avatars"
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
