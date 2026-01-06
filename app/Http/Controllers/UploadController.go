package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
)

// UploadController handles file uploads
type UploadController struct{}

// NewUploadController creates a new upload controller instance
func NewUploadController() *UploadController {
	return &UploadController{}
}

// UploadFile handles file upload
// POST /api/v1/upload
func (uc *UploadController) UploadFile(c *gin.Context) {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "no_file", "No file uploaded", nil)
		return
	}

	// Get file type from query or default to berita
	fileType := c.DefaultQuery("type", "berita")

	// Validate file type
	uploadType := utils.FileUploadType(fileType)
	config, exists := utils.FileUploadConfigs[uploadType]
	if !exists {
		utils.Error(c, http.StatusBadRequest, "invalid_type", "Invalid file type", nil)
		return
	}

	// Set storage path (lazy load from environment)
	config.StoragePath = utils.GetStoragePathForConfig(config)

	// Update config in map
	utils.FileUploadConfigs[uploadType] = config

	// Create upload service
	service, err := utils.NewFileUploadService(uploadType)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "service_error", err.Error(), nil)
		return
	}

	// Upload file with identifier (can be user ID or timestamp)
	identifier := fmt.Sprintf("%d", time.Now().Unix())
	uploadInfo, err := service.Upload(file, identifier)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "upload_error", err.Error(), nil)
		return
	}

	utils.Success(c, http.StatusOK, "File uploaded successfully", gin.H{
		"url":          uploadInfo.FileURL,
		"filename":     filepath.Base(uploadInfo.StoragePath),
		"originalName": uploadInfo.OriginalFilename,
		"size":         uploadInfo.FileSize,
		"type":         uploadInfo.MimeType,
	})
}
