package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	requests "github.com/cvudumbarainformatika/backend/app/Http/Requests"
	models "github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/cvudumbarainformatika/backend/config"
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type DocumentController struct {
	db     *sqlx.DB
	config *config.Config
}

func NewDocumentController(db *sqlx.DB, cfg *config.Config) *DocumentController {
	return &DocumentController{
		db:     db,
		config: cfg,
	}
}

// GetList retrieves documents. Admins see all, users see their own.
// GET /api/v1/documents
func (dc *DocumentController) GetList(c *gin.Context) {
	// Authentication checks
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "unauthorized", "User not authenticated", nil)
		return
	}
	userID := userIDInterface.(int64)

	roleInterface, roleExists := c.Get("user_role")
	if !roleExists {
		utils.Error(c, http.StatusUnauthorized, "unauthorized", "User role not found", nil)
		return
	}
	role := roleInterface.(string)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	status := c.Query("status")
	docType := c.Query("type")

	filters := map[string]interface{}{}
	if status != "" {
		filters["status"] = status
	}
	if docType != "" {
		filters["type"] = docType
	}

	var documents []models.Document
	var total int64
	var err error

	if role == "admin" || role == "super_admin" {
		// Admin logic: If a specific user_id is requested, fetch theirs. Otherwise fetch all.
		targetUserIDStr := c.Query("user_id")
		if targetUserIDStr != "" {
			targetUserID, errConv := strconv.ParseInt(targetUserIDStr, 10, 64)
			if errConv == nil {
				documents, total, err = models.GetUserDocuments(dc.db, targetUserID, filters, offset, limit)
			} else {
				utils.Error(c, http.StatusBadRequest, "invalid_request", "Invalid user ID", nil)
				return
			}
		} else {
			// For admins to see EVERYTHING, we need a slight query modification.
			// Currently models.GetUserDocuments enforces user_id = ?.
			// Let's implement a quick custom query for admin viewing all active documents.
			query := `SELECT id, user_id, name, type, valid_until, status, file_url, created_at, updated_at FROM documents WHERE deleted_at IS NULL`
			countQuery := `SELECT COUNT(*) FROM documents WHERE deleted_at IS NULL`

			var args []interface{}
			if docType != "" {
				query += ` AND type = ?`
				countQuery += ` AND type = ?`
				args = append(args, docType)
			}
			if status != "" {
				query += ` AND status = ?`
				countQuery += ` AND status = ?`
				args = append(args, status)
			}

			err = dc.db.Get(&total, countQuery, args...)
			if err == nil {
				query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
				args = append(args, limit, offset)
				err = dc.db.Select(&documents, query, args...)
			}
		}
	} else {
		// Regular user logic
		documents, total, err = models.GetUserDocuments(dc.db, userID, filters, offset, limit)
	}

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to retrieve documents", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Documents retrieved successfully", gin.H{
		"items": documents,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// Upload handles creating a new document
// POST /api/v1/documents
func (dc *DocumentController) Upload(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "unauthorized", "User not authenticated", nil)
		return
	}
	userID := userIDInterface.(int64)

	var req requests.CreateDocumentRequest
	if err := req.Validate(c); err != nil {
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "file_error", "Failed to get file", err.Error())
		return
	}

	// Use FileUploadService
	uploadType := utils.FileTypeDokumen
	uploadConfig, exists := utils.FileUploadConfigs[uploadType]
	if !exists {
		utils.Error(c, http.StatusInternalServerError, "config_error", "Document upload config not found", nil)
		return
	}
	uploadConfig.StoragePath = utils.GetStoragePathForConfig(uploadConfig)
	utils.FileUploadConfigs[uploadType] = uploadConfig

	service, err := utils.NewFileUploadService(uploadType)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "service_error", "Failed to initialize upload service", nil)
		return
	}

	identifier := fmt.Sprintf("%d_%d", userID, time.Now().Unix())
	uploadInfo, err := service.Upload(file, identifier)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "upload_error", err.Error(), nil)
		return
	}

	fileUrl := uploadInfo.FileURL

	var validUntil *time.Time
	if req.ValidUntil != "" && req.ValidUntil != "null" {
		parsedDate, errParse := time.Parse("2006-01-02", req.ValidUntil)
		if errParse == nil {
			validUntil = &parsedDate
		}
	}

	document := &models.Document{
		UserID:     userID,
		Name:       req.Name,
		Type:       req.Type,
		ValidUntil: validUntil,
		Status:     "pending", // default status
		FileURL:    fileUrl,
	}

	if err := document.Create(dc.db); err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Gagal menyimpan data dokumen", err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, "Dokumen berhasil diunggah", document)
}

// Delete removes a document
// DELETE /api/v1/documents/:id
func (dc *DocumentController) Delete(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "unauthorized", "User not authenticated", nil)
		return
	}
	userID := userIDInterface.(int64)

	roleInterface, _ := c.Get("user_role")
	role := roleInterface.(string)

	docID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid document ID", nil)
		return
	}

	document, err := models.FindDocumentByID(dc.db, docID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to retrieve document", err.Error())
		return
	}
	if document == nil {
		utils.Error(c, http.StatusNotFound, "not_found", "Document not found", nil)
		return
	}

	// Authorization check
	if role != "admin" && role != "super_admin" && document.UserID != userID {
		utils.Error(c, http.StatusForbidden, "forbidden", "You do not have permission to delete this document", nil)
		return
	}

	if err := document.Delete(dc.db); err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to delete document", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Dokumen berhasil dihapus", nil)
}
