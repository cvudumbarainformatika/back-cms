# Scalable File Upload System - Complete Guide

## Overview

Sistem file upload yang scalable untuk menangani berbagai jenis file:
- 👤 Avatar (profile pictures)
- 📸 Thumbnail (berita, galeri)
- 📰 Berita (news images)
- 📄 Dokumen (PDF, Word)
- 🖼️ Galeri (gallery images)
- 📎 Attachment (various documents)

---

## Architecture

### File Types & Configuration

```
┌──────────────────────────────────────────────────────────────┐
│ FileUploadService (Generic)                                  │
│                                                               │
│ Supported Types:                                             │
│ - avatar      (5MB,   JPG/PNG/GIF/WebP)                     │
│ - thumbnail   (10MB,  JPG/PNG/WebP, create thumb)           │
│ - berita      (20MB,  JPG/PNG/WebP, create thumb)           │
│ - dokumen     (50MB,  PDF/DOC/DOCX)                         │
│ - galeri      (25MB,  JPG/PNG/WebP, create thumb)           │
│ - attachment  (100MB, PDF/JPG/PNG/ZIP)                      │
└──────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────┐
│ Storage Structure (VPS)                                      │
│                                                               │
│ /data/                                                       │
│ ├── avatars/           → avatar_1_1735427370.jpg            │
│ ├── thumbnails/        → thumbnail_1_1735427370.jpg         │
│ ├── berita/            → berita_1_1735427370.jpg            │
│ ├── dokumen/           → dokumen_1_1735427370.pdf           │
│ ├── galeri/            → galeri_1_1735427370.jpg            │
│ └── attachments/       → attachment_1_1735427370.zip        │
└──────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────┐
│ Docker Volume Mounting                                       │
│                                                               │
│ Container /app/storage/avatars    ← /data/avatars           │
│ Container /app/storage/thumbnails ← /data/thumbnails        │
│ etc...                                                       │
└──────────────────────────────────────────────────────────────┘
```

---

## File Upload Service

### FileUploadType Constants

```go
FileTypeAvatar       = "avatar"
FileTypeThumbnail    = "thumbnail"
FileTypeBerita       = "berita"
FileTypeDokumen      = "dokumen"
FileTypeGaleri       = "galeri"
FileTypeAttachment   = "attachment"
```

### Configuration Structure

```go
type FileUploadConfig struct {
    StoragePath   string   // Directory path
    MaxSize       int64    // Max file size
    AllowedTypes  []string // MIME types
    AllowedExts   []string // Extensions
    CreateThumb   bool     // Create thumbnail?
    ThumbWidth    int      // Thumbnail width
    ThumbHeight   int      // Thumbnail height
}
```

### Usage Example

```go
// Create service for specific file type
service, err := utils.NewFileUploadService(utils.FileTypeBerita)
if err != nil {
    return err
}

// Upload file
uploadInfo, err := service.Upload(file, userID)
if err != nil {
    return err
}

// Access file info
fmt.Println(uploadInfo.FileURL)        // /api/v1/files/berita/berita_1_xxx.jpg
fmt.Println(uploadInfo.StoragePath)    // /app/storage/berita/berita_1_xxx.jpg
fmt.Println(uploadInfo.ThumbnailURL)   // /api/v1/files/berita/berita_1_xxx_thumb.jpg
```

---

## API Endpoints

### Upload File (Endpoint Examples)

**Upload Avatar:**
```
PUT /api/v1/auth/profile
Authorization: Bearer {token}
Content-Type: multipart/form-data

Form Data:
  name: "User Name"
  avatar: <file>
```

**Upload Berita Image:**
```
POST /api/v1/berita
Authorization: Bearer {token}
Content-Type: multipart/form-data

Form Data:
  title: "Berita Title"
  content: "Content"
  thumbnail: <file>
```

### Serve File

**Get Any File:**
```
GET /api/v1/files/{file_type}/{filename}

Examples:
GET /api/v1/files/avatar/avatar_1_1735427370.jpg
GET /api/v1/files/berita/berita_1_1735427370.jpg
GET /api/v1/files/dokumen/dokumen_1_1735427370.pdf
GET /api/v1/files/attachment/attachment_1_1735427370.zip
```

**List Files (Admin):**
```
GET /api/v1/files/{file_type}/list

Example:
GET /api/v1/files/berita/list
Response: [
  {
    "name": "berita_1_1735427370.jpg",
    "size": 1024000,
    "modified": "2025-12-29T18:30:00Z",
    "url": "/api/v1/files/berita/berita_1_1735427370.jpg"
  }
]
```

**Delete File (Admin):**
```
DELETE /api/v1/files/{file_type}/{filename}

Example:
DELETE /api/v1/files/berita/berita_1_1735427370.jpg
```

---

## Environment Variables

### Development (.env)

```env
# Base storage path
STORAGE_BASE_PATH=./storage

# Avatar
AVATAR_STORAGE_PATH=./storage/avatars
AVATAR_UPLOAD_ENABLED=true
AVATAR_MAX_SIZE=5242880

# Thumbnail
STORAGE_THUMBNAILS_PATH=./storage/thumbnails
THUMBNAIL_UPLOAD_ENABLED=true
THUMBNAIL_MAX_SIZE=10485760

# Berita
STORAGE_BERITA_PATH=./storage/berita
BERITA_UPLOAD_ENABLED=true
BERITA_MAX_SIZE=20971520

# Dokumen
STORAGE_DOKUMEN_PATH=./storage/dokumen
DOKUMEN_UPLOAD_ENABLED=true
DOKUMEN_MAX_SIZE=52428800

# Galeri
STORAGE_GALERI_PATH=./storage/galeri
GALERI_UPLOAD_ENABLED=true
GALERI_MAX_SIZE=26214400

# Attachment
STORAGE_ATTACHMENT_PATH=./storage/attachments
ATTACHMENT_UPLOAD_ENABLED=true
ATTACHMENT_MAX_SIZE=104857600
```

### Production (docker-compose.prod.yml)

```yaml
environment:
  STORAGE_BASE_PATH: /app/storage
  AVATAR_STORAGE_PATH: /app/storage/avatars
  STORAGE_THUMBNAILS_PATH: /app/storage/thumbnails
  STORAGE_BERITA_PATH: /app/storage/berita
  STORAGE_DOKUMEN_PATH: /app/storage/dokumen
  STORAGE_GALERI_PATH: /app/storage/galeri
  STORAGE_ATTACHMENT_PATH: /app/storage/attachments
```

---

## File Size Limits

| Type | Max Size | Use Case |
|------|----------|----------|
| Avatar | 5 MB | Profile pictures |
| Thumbnail | 10 MB | Berita/galeri thumbnails |
| Berita | 20 MB | News article images |
| Dokumen | 50 MB | PDF/Word documents |
| Galeri | 25 MB | Gallery images |
| Attachment | 100 MB | General attachments |

---

## Allowed File Types

### Image Types
- **avatar**: JPG, JPEG, PNG, GIF, WebP
- **thumbnail**: JPG, JPEG, PNG, WebP
- **berita**: JPG, JPEG, PNG, WebP
- **galeri**: JPG, JPEG, PNG, WebP

### Document Types
- **dokumen**: PDF, DOC, DOCX
- **attachment**: PDF, JPG, JPEG, PNG, ZIP

---

## Docker Volume Configuration

### Development (docker-compose.yml)

```yaml
volumes:
  - avatars_data:/app/storage/avatars
  - thumbnails_data:/app/storage/thumbnails
  - berita_data:/app/storage/berita
  - dokumen_data:/app/storage/dokumen
  - galeri_data:/app/storage/galeri
  - attachments_data:/app/storage/attachments
```

### Production (docker-compose.prod.yml)

```yaml
volumes:
  - avatars_prod_data:/app/storage/avatars
  - thumbnails_prod_data:/app/storage/thumbnails
  - berita_prod_data:/app/storage/berita
  - dokumen_prod_data:/app/storage/dokumen
  - galeri_prod_data:/app/storage/galeri
  - attachments_prod_data:/app/storage/attachments
```

---

## File Naming Convention

Each uploaded file follows this pattern:
```
{file_type}_{identifier}_{timestamp}.{extension}

Examples:
avatar_1_1735427370.jpg              (user_id = 1)
thumbnail_1_1735427370.jpg           (article_id = 1)
berita_5_1735427370.jpg              (article_id = 5)
dokumen_2_1735427370.pdf             (document_id = 2)
galeri_3_1735427370.jpg              (gallery_id = 3)
attachment_10_1735427370.zip         (attachment_id = 10)
```

---

## Security Features

### ✅ Path Protection
- Directory traversal prevention
- Filename pattern validation
- Path bounds checking
- No `../` or absolute paths allowed

### ✅ File Validation
- MIME type checking
- File extension validation
- File size limits per type
- Magic number verification (future)

### ✅ Access Control
- JWT authentication for upload
- Old files auto-deleted on replacement
- Admin-only delete/list operations
- Cache control headers

### ✅ Docker Security
- Named volumes (not bind mounts)
- Separate volume per file type
- Volume independent from app code
- Backup/restore capability

---

## Usage Examples

### 1. Upload Avatar

```bash
curl -X PUT http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer TOKEN" \
  -F "name=User Name" \
  -F "avatar=@avatar.jpg"

# Response
{
  "data": {
    "avatar": "/api/v1/files/avatar/avatar_1_1735427370.jpg",
    "avatar_path": "/app/storage/avatars/avatar_1_1735427370.jpg"
  }
}
```

### 2. Upload Berita Thumbnail

```bash
curl -X POST http://localhost:8080/api/v1/berita \
  -H "Authorization: Bearer TOKEN" \
  -F "title=News Title" \
  -F "content=Content here" \
  -F "thumbnail=@image.jpg"

# Response
{
  "data": {
    "thumbnail_url": "/api/v1/files/berita/berita_1_1735427370.jpg",
    "thumbnail_path": "/app/storage/berita/berita_1_1735427370.jpg"
  }
}
```

### 3. Upload Document

```bash
curl -X POST http://localhost:8080/api/v1/dokumen/upload \
  -H "Authorization: Bearer TOKEN" \
  -F "title=Document" \
  -F "file=@document.pdf"

# Response
{
  "data": {
    "file_url": "/api/v1/files/dokumen/dokumen_1_1735427370.pdf",
    "file_path": "/app/storage/dokumen/dokumen_1_1735427370.pdf"
  }
}
```

### 4. Access File

```bash
# Get avatar
curl http://localhost:8080/api/v1/files/avatar/avatar_1_1735427370.jpg

# Get berita image
curl http://localhost:8080/api/v1/files/berita/berita_1_1735427370.jpg

# Get document
curl http://localhost:8080/api/v1/files/dokumen/dokumen_1_1735427370.pdf
```

### 5. List Files (Admin)

```bash
curl http://localhost:8080/api/v1/files/berita/list \
  -H "Authorization: Bearer ADMIN_TOKEN"

# Response
{
  "data": [
    {
      "name": "berita_1_1735427370.jpg",
      "size": 1024000,
      "modified": "2025-12-29T18:30:00Z",
      "url": "/api/v1/files/berita/berita_1_1735427370.jpg"
    }
  ]
}
```

---

## Backup & Restore

### Backup All Storage

```bash
# Backup specific file type
docker run --rm -v berita_data:/data -v $(pwd):/backup \
  alpine tar czf /backup/berita-backup.tar.gz -C /data .

# Backup all
for type in avatars thumbnails berita dokumen galeri attachments; do
  docker run --rm -v ${type}_data:/data -v $(pwd):/backup \
    alpine tar czf /backup/${type}-backup.tar.gz -C /data .
done
```

### Restore Files

```bash
# Restore specific type
docker run --rm -v berita_data:/data -v $(pwd):/backup \
  alpine tar xzf /backup/berita-backup.tar.gz -C /data

# Restore all
for type in avatars thumbnails berita dokumen galeri attachments; do
  docker run --rm -v ${type}_data:/data -v $(pwd):/backup \
    alpine tar xzf /backup/${type}-backup.tar.gz -C /data
done
```

---

## Cache Headers

| File Type | Cache Duration |
|-----------|-----------------|
| avatar | 24 hours |
| thumbnail | 24 hours |
| berita | 24 hours |
| galeri | 24 hours |
| dokumen | 7 days |
| attachment | 7 days |

---

## Future Enhancements

- [ ] Image resizing on upload
- [ ] Thumbnail generation (for images)
- [ ] Image compression (WebP conversion)
- [ ] Virus scanning
- [ ] CDN integration (S3, Cloudinary)
- [ ] Automatic cleanup of old files
- [ ] File versioning
- [ ] Audit logging

---

## Troubleshooting

### Problem: File not found after upload
**Solution:** Check if volume is mounted correctly in Docker

### Problem: Permission denied writing files
**Solution:** Ensure storage directory has correct permissions (755)

### Problem: Disk space full
**Solution:** Check volume size and implement cleanup policy

### Problem: Files lost after rebuild
**Solution:** Ensure named volumes are used (not inline volumes)

---

## Related Documentation

- `STORAGE_ARCHITECTURE.md` - Container vs VPS explanation
- `COMMON_MISTAKES_FIXED.md` - Common pitfalls
- `AVATAR_STORAGE_GUIDE.md` - Avatar-specific guide

---

**Status:** ✅ Production Ready

This system is ready for production deployment with multiple file types and secure storage!
