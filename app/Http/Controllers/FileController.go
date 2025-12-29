package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

// FileController handles file serving for all upload types
type FileController struct{}

// NewFileController creates a new file controller instance
func NewFileController() *FileController {
	return &FileController{}
}

// ServeFile serves a file by type and filename
// GET /api/v1/files/:file_type/:filename
func (fc *FileController) ServeFile(c *gin.Context) {
	fileType := c.Param("file_type")
	filename := c.Param("filename")

	if fileType == "" || filename == "" {
		utils.Error(c, http.StatusBadRequest, "invalid_request", "File type and filename are required", nil)
		return
	}

	// Validate file type
	fileUploadType := utils.FileUploadType(fileType)
	config, exists := utils.FileUploadConfigs[fileUploadType]
	if !exists {
		utils.Error(c, http.StatusBadRequest, "invalid_file_type", "Invalid file type: "+fileType, nil)
		return
	}

	// Validate filename format (prevent directory traversal)
	if !isValidFilename(filename) {
		utils.Error(c, http.StatusForbidden, "invalid_filename", "Invalid filename format", nil)
		return
	}

	// Get storage path (lazy loaded from environment)
	storagePath := utils.GetStoragePathForConfig(config)

	// Build safe file path
	filePath := filepath.Join(storagePath, filename)

	// Double-check: ensure path is within storage directory
	absPath, _ := filepath.Abs(filePath)
	absStoragePath, _ := filepath.Abs(storagePath)
	if !strings.HasPrefix(absPath, absStoragePath) {
		utils.Error(c, http.StatusForbidden, "access_denied", "Access denied", nil)
		return
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.Error(c, http.StatusNotFound, "not_found", "File not found", nil)
		return
	}

	// Set cache headers based on file type
	cacheControl := "public, max-age=86400"  // 24 hours default
	if fileType == "attachment" || fileType == "dokumen" {
		cacheControl = "public, max-age=604800"  // 7 days for documents
	}

	c.Header("Cache-Control", cacheControl)
	c.Header("X-Content-Type-Options", "nosniff")

	// Serve file
	c.File(filePath)
}

// ListFiles lists all files of a specific type (admin only)
// GET /api/v1/files/:file_type/list
func (fc *FileController) ListFiles(c *gin.Context) {
	fileType := c.Param("file_type")

	// TODO: Add admin authorization check
	// if !isAdmin(c) {
	//     utils.Error(c, http.StatusForbidden, "forbidden", "Admin access required", nil)
	//     return
	// }

	fileUploadType := utils.FileUploadType(fileType)
	config, exists := utils.FileUploadConfigs[fileUploadType]
	if !exists {
		utils.Error(c, http.StatusBadRequest, "invalid_file_type", "Invalid file type", nil)
		return
	}

	// Get storage path (lazy loaded from environment)
	storagePath := utils.GetStoragePathForConfig(config)

	// Read directory
	files, err := os.ReadDir(storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			utils.Success(c, http.StatusOK, "No files found", []string{})
			return
		}
		utils.Error(c, http.StatusInternalServerError, "read_error", "Failed to read files", nil)
		return
	}

	var fileList []gin.H
	for _, file := range files {
		if !file.IsDir() {
			info, _ := file.Info()
			fileList = append(fileList, gin.H{
				"name":       file.Name(),
				"size":       info.Size(),
				"modified":   info.ModTime(),
				"url":        "/api/v1/files/" + fileType + "/" + file.Name(),
			})
		}
	}

	utils.Success(c, http.StatusOK, "Files retrieved", fileList)
}

// DeleteFile deletes a file (admin only)
// DELETE /api/v1/files/:file_type/:filename
func (fc *FileController) DeleteFile(c *gin.Context) {
	fileType := c.Param("file_type")
	filename := c.Param("filename")

	// TODO: Add admin authorization check
	// if !isAdmin(c) {
	//     utils.Error(c, http.StatusForbidden, "forbidden", "Admin access required", nil)
	//     return
	// }

	if fileType == "" || filename == "" {
		utils.Error(c, http.StatusBadRequest, "invalid_request", "File type and filename are required", nil)
		return
	}

	// Create service
	service, err := utils.NewFileUploadService(utils.FileUploadType(fileType))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_file_type", err.Error(), nil)
		return
	}

	// Delete file
	if err := service.Delete(filename); err != nil {
		utils.Error(c, http.StatusInternalServerError, "delete_error", err.Error(), nil)
		return
	}

	utils.Success(c, http.StatusOK, "File deleted successfully", gin.H{
		"filename": filename,
		"type":     fileType,
	})
}

// isValidFilename validates filename format
func isValidFilename(filename string) bool {
	// Must start with file type prefix
	hasPrefixPattern := false
	for _, fileType := range []string{"avatar_", "thumbnail_", "berita_", "dokumen_", "galeri_", "attachment_"} {
		if strings.HasPrefix(filename, fileType) {
			hasPrefixPattern = true
			break
		}
	}

	if !hasPrefixPattern {
		return false
	}

	// Must have extension
	if !strings.Contains(filename, ".") {
		return false
	}

	// No directory separators
	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return false
	}

	return true
}
