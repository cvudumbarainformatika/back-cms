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

// FileUploadType defines the type of file being uploaded
type FileUploadType string

const (
	FileTypeAvatar        FileUploadType = "avatar"
	FileTypeThumbnail     FileUploadType = "thumbnail"
	FileTypeBerita        FileUploadType = "berita"
	FileTypeDokumen       FileUploadType = "dokumen"
	FileTypeGaleri        FileUploadType = "galeri"
	FileTypeAttachment    FileUploadType = "attachment"
)

// FileUploadConfig contains configuration for each file type
type FileUploadConfig struct {
	FileType     FileUploadType // Type of file
	StoragePath  string         // Directory path for storage (lazy-loaded)
	MaxSize      int64          // Max file size in bytes
	AllowedTypes []string       // Allowed MIME types
	AllowedExts  []string       // Allowed file extensions
	CreateThumb  bool           // Whether to create thumbnail
	ThumbWidth   int            // Thumbnail width
	ThumbHeight  int            // Thumbnail height
}

// FileUploadConfigs contains all file type configurations (without StoragePath - lazy loaded)
var FileUploadConfigs = map[FileUploadType]FileUploadConfig{
	FileTypeAvatar: {
		FileType:     FileTypeAvatar,
		MaxSize:      5 * 1024 * 1024, // 5MB
		AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
		AllowedExts:  []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		CreateThumb:  false,
	},
	FileTypeThumbnail: {
		FileType:     FileTypeThumbnail,
		MaxSize:      10 * 1024 * 1024, // 10MB
		AllowedTypes: []string{"image/jpeg", "image/png", "image/webp"},
		AllowedExts:  []string{".jpg", ".jpeg", ".png", ".webp"},
		CreateThumb:  true,
		ThumbWidth:   300,
		ThumbHeight:  300,
	},
	FileTypeBerita: {
		FileType:     FileTypeBerita,
		MaxSize:      20 * 1024 * 1024, // 20MB
		AllowedTypes: []string{"image/jpeg", "image/png", "image/webp"},
		AllowedExts:  []string{".jpg", ".jpeg", ".png", ".webp"},
		CreateThumb:  true,
		ThumbWidth:   800,
		ThumbHeight:  600,
	},
	FileTypeDokumen: {
		FileType:     FileTypeDokumen,
		MaxSize:      50 * 1024 * 1024, // 50MB
		AllowedTypes: []string{"application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
		AllowedExts:  []string{".pdf", ".doc", ".docx"},
		CreateThumb:  false,
	},
	FileTypeGaleri: {
		FileType:     FileTypeGaleri,
		MaxSize:      25 * 1024 * 1024, // 25MB
		AllowedTypes: []string{"image/jpeg", "image/png", "image/webp"},
		AllowedExts:  []string{".jpg", ".jpeg", ".png", ".webp"},
		CreateThumb:  true,
		ThumbWidth:   500,
		ThumbHeight:  500,
	},
	FileTypeAttachment: {
		FileType:     FileTypeAttachment,
		MaxSize:      100 * 1024 * 1024, // 100MB
		AllowedTypes: []string{"application/pdf", "image/jpeg", "image/png", "application/zip"},
		AllowedExts:  []string{".pdf", ".jpg", ".jpeg", ".png", ".zip"},
		CreateThumb:  false,
	},
}

// getStoragePathForType returns the storage path for a file type from environment (lazy loaded)
func getStoragePathForType(fileType string) string {
	// Special case untuk avatar (legacy name: AVATAR_STORAGE_PATH, bukan STORAGE_AVATAR_PATH)
	var envVar string
	if fileType == "avatar" {
		envVar = "AVATAR_STORAGE_PATH"
	} else {
		envVar = fmt.Sprintf("STORAGE_%s_PATH", strings.ToUpper(fileType))
	}
	
	storagePath := os.Getenv(envVar)
	if storagePath == "" {
		// Fallback to base storage path
		baseStorage := GetStoragePath()
		storagePath = filepath.Join(baseStorage, fileType)
	}
	log.Printf("[Storage] %s path resolved to: %s (from env: %s)", fileType, storagePath, envVar)
	return storagePath
}

// GetStoragePathForConfig returns the storage path for a config (with lazy loading)
func GetStoragePathForConfig(config FileUploadConfig) string {
	if config.StoragePath != "" {
		return config.StoragePath
	}
	// Lazy load from environment
	fileTypeStr := strings.ToLower(string(config.FileType))
	return getStoragePathForType(fileTypeStr)
}

// UploadedFileInfo contains info about uploaded file
type UploadedFileInfo struct {
	OriginalFilename string
	StoragePath      string // Full path where file is stored
	FileURL          string // Public URL to access file
	FileSize         int64
	MimeType         string
	UploadedAt       time.Time
	ThumbnailURL     string // Thumbnail URL if created
}

// FileUploadService handles file uploads
type FileUploadService struct {
	config FileUploadConfig
	fileType FileUploadType
}

// NewFileUploadService creates a new file upload service
func NewFileUploadService(fileType FileUploadType) (*FileUploadService, error) {
	config, exists := FileUploadConfigs[fileType]
	if !exists {
		return nil, fmt.Errorf("unknown file type: %s", fileType)
	}

	return &FileUploadService{
		config: config,
		fileType: fileType,
	}, nil
}

// Upload handles file upload with validation and storage
func (s *FileUploadService) Upload(file *multipart.FileHeader, identifier string) (*UploadedFileInfo, error) {
	// Validate file size
	if file.Size > s.config.MaxSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.config.MaxSize)
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Validate MIME type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	mimeType := file.Header.Get("Content-Type")
	if !s.isAllowedMimeType(mimeType) {
		return nil, fmt.Errorf("file type not allowed. Allowed types: %v", s.config.AllowedTypes)
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !s.isAllowedExtension(ext) {
		return nil, fmt.Errorf("file extension not allowed. Allowed extensions: %v", s.config.AllowedExts)
	}

	// Reset file pointer
	src.Seek(0, 0)

	// Create storage directory
	if err := os.MkdirAll(s.config.StoragePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %v", err)
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s_%s_%d%s", string(s.fileType), identifier, timestamp, ext)
	filePath := filepath.Join(s.config.StoragePath, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// Build public URL
	fileURL := fmt.Sprintf("/api/v1/files/%s/%s", s.fileType, filename)

	log.Printf("[%s] File uploaded: %s (size: %d bytes)", s.fileType, filename, file.Size)

	uploadInfo := &UploadedFileInfo{
		OriginalFilename: file.Filename,
		StoragePath:      filePath,
		FileURL:          fileURL,
		FileSize:         file.Size,
		MimeType:         mimeType,
		UploadedAt:       time.Now(),
	}

	// TODO: Create thumbnail if needed
	// if s.config.CreateThumb {
	//     thumbURL, err := s.createThumbnail(filePath, filename)
	//     if err == nil {
	//         uploadInfo.ThumbnailURL = thumbURL
	//     }
	// }

	return uploadInfo, nil
}

// Delete removes uploaded file
func (s *FileUploadService) Delete(filename string) error {
	if filename == "" {
		return nil
	}

	filePath := filepath.Join(s.config.StoragePath, filename)

	// Safety check: ensure path is within storage directory
	absPath, _ := filepath.Abs(filePath)
	absStoragePath, _ := filepath.Abs(s.config.StoragePath)
	if !strings.HasPrefix(absPath, absStoragePath) {
		return fmt.Errorf("invalid file path")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("[%s] File not found for deletion: %s", s.fileType, filename)
		return nil
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	log.Printf("[%s] File deleted: %s", s.fileType, filename)

	// TODO: Delete thumbnail if exists
	// if s.config.CreateThumb {
	//     thumbFilename := strings.TrimSuffix(filename, ext) + "_thumb" + ext
	//     os.Remove(filepath.Join(s.config.StoragePath, thumbFilename))
	// }

	return nil
}

// GetFilePath returns the full storage path for a filename
func (s *FileUploadService) GetFilePath(filename string) (string, error) {
	filePath := filepath.Join(s.config.StoragePath, filename)

	// Safety check
	absPath, _ := filepath.Abs(filePath)
	absStoragePath, _ := filepath.Abs(s.config.StoragePath)
	if !strings.HasPrefix(absPath, absStoragePath) {
		return "", fmt.Errorf("invalid file path")
	}

	return filePath, nil
}

// Helper methods
func (s *FileUploadService) isAllowedMimeType(mimeType string) bool {
	for _, allowed := range s.config.AllowedTypes {
		if mimeType == allowed {
			return true
		}
	}
	return false
}

func (s *FileUploadService) isAllowedExtension(ext string) bool {
	for _, allowed := range s.config.AllowedExts {
		if ext == allowed {
			return true
		}
	}
	return false
}
