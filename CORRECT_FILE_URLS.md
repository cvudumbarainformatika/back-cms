# Correct File URLs - Important!

## ❌ WRONG URLs (Old Static File Serving)

```
❌ http://localhost:8080/storage/avatars/avatar_1_1767011202.png
❌ http://localhost:8080/storage/berita/berita_1_1767011202.jpg
❌ http://localhost:8080/storage/dokumen/dokumen_1_1767011202.pdf
```

**Why wrong:**
- Direct file access (no validation)
- Path traversal vulnerability
- No access control
- No audit trail
- Can be brute-forced

---

## ✅ CORRECT URLs (Secure API Endpoints)

```
✅ http://localhost:8080/api/v1/files/avatar/avatar_1_1767011202.png
✅ http://localhost:8080/api/v1/files/berita/berita_1_1767011202.jpg
✅ http://localhost:8080/api/v1/files/dokumen/dokumen_1_1767011202.pdf
✅ http://localhost:8080/api/v1/files/galeri/galeri_1_1767011202.jpg
✅ http://localhost:8080/api/v1/files/thumbnail/thumbnail_1_1767011202.jpg
✅ http://localhost:8080/api/v1/files/attachment/attachment_1_1767011202.zip
```

**Why correct:**
- ✅ Melalui FileController
- ✅ Filename validation
- ✅ Path traversal protection
- ✅ File type validation
- ✅ Cache headers
- ✅ Security headers
- ✅ Audit trail possible

---

## URL Format

```
/api/v1/files/{file_type}/{filename}

Pattern:
/api/v1/files/avatar/avatar_1_1767011202.png
              ↑       ↑
           Type    Filename
```

### File Types:
- `avatar` - Profile pictures
- `thumbnail` - Thumbnails
- `berita` - News images
- `dokumen` - Documents
- `galeri` - Gallery images
- `attachment` - Attachments

### Filename Format:
```
{type}_{identifier}_{timestamp}.{ext}

Examples:
avatar_1_1767011202.png           (user_id=1)
berita_5_1767011202.jpg           (article_id=5)
dokumen_2_1767011202.pdf          (doc_id=2)
galeri_3_1767011202.jpg           (gallery_id=3)
thumbnail_1_1767011202.jpg        (item_id=1)
attachment_10_1767011202.zip      (attach_id=10)
```

---

## API Endpoints Reference

### Serve File (Public)
```
GET /api/v1/files/:file_type/:filename

Example:
curl http://localhost:8080/api/v1/files/avatar/avatar_1_1767011202.png

Response: Binary image file with headers
Content-Type: image/png
Cache-Control: public, max-age=86400
X-Content-Type-Options: nosniff
```

### List Files (Admin)
```
GET /api/v1/files/:file_type/list

Example:
curl http://localhost:8080/api/v1/files/berita/list

Response: JSON array of files
[
  {
    "name": "berita_1_1767011202.jpg",
    "size": 1024000,
    "modified": "2025-12-29T18:30:00Z",
    "url": "/api/v1/files/berita/berita_1_1767011202.jpg"
  }
]
```

### Delete File (Admin)
```
DELETE /api/v1/files/:file_type/:filename

Example:
curl -X DELETE http://localhost:8080/api/v1/files/berita/berita_1_1767011202.jpg \
  -H "Authorization: Bearer ADMIN_TOKEN"

Response: JSON
{
  "status": "success",
  "message": "File deleted successfully",
  "data": {
    "filename": "berita_1_1767011202.jpg",
    "type": "berita"
  }
}
```

---

## Database Storage

When user uploads file, database stores:

```sql
-- Avatar example
UPDATE users SET avatar = '/api/v1/files/avatar/avatar_1_1767011202.png'

-- Berita example  
UPDATE berita SET thumbnail = '/api/v1/files/berita/berita_5_1767011202.jpg'

-- Document example
INSERT INTO documents (file_url) VALUES ('/api/v1/files/dokumen/dokumen_2_1767011202.pdf')
```

Note: Database stores the **public API URL**, not the file path!

---

## Frontend Usage

### Display Avatar Image

```html
<!-- JavaScript fetch -->
<img id="avatar" src="/api/v1/files/avatar/avatar_1_1767011202.png" alt="Avatar" />

<!-- Or from database -->
<script>
  fetch('/api/v1/auth/me')
    .then(r => r.json())
    .then(data => {
      document.getElementById('avatar').src = data.data.avatar;
      // Loads: /api/v1/files/avatar/avatar_1_1767011202.png
    });
</script>
```

### Display News Image

```html
<img src="/api/v1/files/berita/berita_5_1767011202.jpg" alt="News" />
```

### Download Document

```html
<a href="/api/v1/files/dokumen/dokumen_2_1767011202.pdf" download>
  Download PDF
</a>
```

---

## Response Headers

All file endpoints return secure headers:

```
HTTP/1.1 200 OK
Content-Type: image/png
Cache-Control: public, max-age=86400
X-Content-Type-Options: nosniff
Content-Length: 1024000
Date: Mon, 29 Dec 2025 18:30:00 GMT
```

### What They Mean:

- **Content-Type** - Correct MIME type (prevents execution)
- **Cache-Control** - Browser cache for 24 hours
- **X-Content-Type-Options: nosniff** - Prevent MIME type sniffing attack
- **Content-Length** - File size

---

## Security Validation

Each request goes through:

1. **Filename Validation**
   ```
   ✅ Must start with: avatar_, berita_, dokumen_, etc.
   ✅ Must have extension: .jpg, .png, .pdf, etc.
   ❌ No path separators: / or \
   ❌ No ../ sequences
   ```

2. **Path Traversal Check**
   ```go
   absPath, _ := filepath.Abs(filePath)
   absStoragePath, _ := filepath.Abs(storagePath)
   if !strings.HasPrefix(absPath, absStoragePath) {
       // Reject - path outside storage directory
   }
   ```

3. **File Type Validation**
   ```
   Check if file exists in correct storage directory
   Only serve if file belongs to requested type
   ```

---

## Comparison: Old vs New

| Feature | Old ❌ | New ✅ |
|---------|--------|--------|
| URL | `/storage/avatars/file.png` | `/api/v1/files/avatar/file.png` |
| Validation | None | Full |
| Security | Low | High |
| Path Protection | No | Yes |
| Access Control | No | Yes |
| Cache Headers | No | Yes |
| Audit Trail | No | Yes |
| Brute Force Safe | No | Yes |

---

## Migration Guide

If you have old URLs in database or frontend:

### Old URL:
```
/storage/avatars/avatar_1_1767011202.png
```

### Convert to New URL:
```
/api/v1/files/avatar/avatar_1_1767011202.png
```

### Script to Update Database

```sql
-- Update users table
UPDATE users 
SET avatar = CONCAT('/api/v1/files/avatar/', SUBSTRING_INDEX(avatar, '/', -1))
WHERE avatar LIKE '/storage/avatars/%';

-- Update berita table
UPDATE berita 
SET thumbnail = CONCAT('/api/v1/files/berita/', SUBSTRING_INDEX(thumbnail, '/', -1))
WHERE thumbnail LIKE '/storage/berita/%';
```

---

## Testing

### Test Avatar URL
```bash
curl -i http://localhost:8080/api/v1/files/avatar/avatar_1_1767011202.png

# Should return:
# HTTP/1.1 200 OK
# Content-Type: image/png
# ... image data ...
```

### Test Berita URL
```bash
curl -i http://localhost:8080/api/v1/files/berita/berita_5_1767011202.jpg

# Should return:
# HTTP/1.1 200 OK  
# Content-Type: image/jpeg
# ... image data ...
```

### Test Invalid URL (should fail)
```bash
curl -i http://localhost:8080/api/v1/files/avatar/../../etc/passwd

# Should return:
# HTTP/1.1 403 Forbidden
# {"status": "error", "message": "Invalid filename format"}
```

---

## Troubleshooting

### Problem: File URL returns 404
**Check:**
1. File exists in storage directory
2. Filename matches pattern: `{type}_{id}_{timestamp}.{ext}`
3. File type is correct (avatar, berita, etc.)

### Problem: Getting 403 Forbidden
**Check:**
1. Filename contains `../` or absolute path
2. Filename has invalid characters
3. File is outside storage directory

### Problem: Wrong MIME type
**Check:**
1. File extension matches content
2. File header/magic bytes match type
3. System correctly identifies file type

---

## Summary

✅ **Always use:** `/api/v1/files/{type}/{filename}`
❌ **Never use:** `/storage/{type}/{filename}`

The new approach provides:
- Security validation
- Path traversal protection
- Proper MIME types
- Cache headers
- Audit capability

---

**Important:** Remove any hardcoded `/storage/` URLs from your frontend and database!
