# Avatar Storage Security Guide

## Overview

Avatar files disimpan secara aman dengan Docker volume mounting yang persisten dan endpoint yang ter-proteksi.

## Architecture

```
┌─────────────────────┐
│   Client Request    │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────────────────────┐
│   Gin Router (Port 8080)            │
│  - PUT /api/v1/auth/profile         │
│  - GET /api/v1/avatars/:user_id     │
└──────────┬──────────────────────────┘
           │
           ▼
┌─────────────────────────────────────┐
│   Auth/Avatar Controller            │
│  - Upload Handler                   │
│  - File Serve Handler               │
└──────────┬──────────────────────────┘
           │
           ▼
┌─────────────────────────────────────┐
│   Docker Volume (avatars_data)      │
│   Location: /app/storage/avatars    │
│   - Persisten across rebuilds       │
│   - Shared between containers       │
└─────────────────────────────────────┘
```

## Development Setup

### 1. Docker Compose Configuration (docker-compose.yml)

```yaml
volumes:
  - avatars_data:/app/storage/avatars  # Development volume
```

**Keuntungan:**
- Data persisten saat container restart
- Data hilang saat `docker-compose down -v`
- Ideal untuk development

### 2. Running Development

```bash
# Start services
docker-compose up -d

# Upload avatar
curl -X PUT http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer TOKEN" \
  -F "name=Admin Pusat" \
  -F "avatar=@image.jpg"

# Access avatar
curl http://localhost:8080/api/v1/avatars/1
```

## Production Setup

### 1. Docker Compose Configuration (docker-compose.prod.yml)

```yaml
volumes:
  - avatars_prod_data:/app/storage/avatars  # Production volume
```

**Keuntungan:**
- Data persisten secara permanen
- Terpisah dari container lifecycle
- Bisa di-backup/restore terpisah
- Multi-container bisa share volume

### 2. Production Deployment

```bash
# Start services
docker-compose -f docker-compose.prod.yml up -d

# Verify volume
docker volume ls | grep avatars_prod_data

# Check storage
docker exec api-go-prod ls -lah /app/storage/avatars/
```

## Storage Structure

```
storage/
└── avatars/
    ├── avatar_1_1735427370.jpg     # User ID 1, uploaded at timestamp
    ├── avatar_1_1735427400.png     # User ID 1, new upload replaces old
    ├── avatar_2_1735427420.jpg     # User ID 2
    └── avatar_3_1735427450.gif     # User ID 3
```

**Naming Convention:** `avatar_{user_id}_{timestamp}.{ext}`

## API Endpoints

### Upload Avatar

```
PUT /api/v1/auth/profile
Authorization: Bearer {token}
Content-Type: multipart/form-data

Form Data:
- name: string (required)
- phone: string (optional)
- address: string (optional)
- bio: string (optional)
- avatar: file (optional, max 5MB)
```

**Response:**
```json
{
  "status": "success",
  "message": "Profile updated successfully",
  "data": {
    "id": 1,
    "name": "Admin Pusat",
    "avatar": "/api/v1/avatars/1",
    "avatar_path": "/storage/avatars/avatar_1_1735427370.jpg"
  }
}
```

### Get Avatar

```
GET /api/v1/avatars/:user_id
```

**Returns:** Binary image file with cache headers

**Response Headers:**
```
Cache-Control: public, max-age=86400
X-Content-Type-Options: nosniff
Content-Type: image/jpeg
```

### Get Avatar by Filename

```
GET /api/v1/avatars/file/:filename
```

**Security:** Filename must match pattern `avatar_*.{jpg,jpeg,png,gif,webp}`

## Security Features

### 1. File Validation
- ✅ Max size: 5MB
- ✅ Allowed types: JPEG, PNG, GIF, WebP
- ✅ MIME type validation

### 2. Path Protection
- ✅ Prevent directory traversal attacks
- ✅ Validate filename pattern
- ✅ Ensure path within avatars directory

### 3. Access Control
- ✅ Endpoint publicly accessible (can be restricted via middleware)
- ✅ Upload requires authentication (JWT)
- ✅ Only authenticated users can upload

### 4. File Naming
- ✅ Unique naming with user ID and timestamp
- ✅ Old avatar automatically deleted on new upload
- ✅ Prevents file overwrites

### 5. Docker Security
- ✅ Persistent volumes survive container recreation
- ✅ Volume independent from application code
- ✅ Can be mounted to multiple containers

## Backup & Restore

### Backup Avatar Data

```bash
# Export volume to tar
docker run --rm -v avatars_data:/data -v $(pwd):/backup \
  alpine tar czf /backup/avatars-backup.tar.gz -C /data .

# Or copy directly
docker run --rm -v avatars_data:/data -v $(pwd):/backup \
  alpine cp -r /data/* /backup/avatars/
```

### Restore Avatar Data

```bash
# Import from tar
docker run --rm -v avatars_data:/data -v $(pwd):/backup \
  alpine tar xzf /backup/avatars-backup.tar.gz -C /data

# Or copy back
docker run --rm -v avatars_data:/data -v $(pwd):/backup \
  alpine cp -r /backup/avatars/* /data/
```

## Troubleshooting

### 1. Avatar Tidak Persisten Setelah Rebuild

**Problem:** Avatar hilang setelah `docker-compose down`

**Solution:**
```bash
# Gunakan named volumes (default behavior)
docker volume ls

# Jangan gunakan -v flag saat down
docker-compose down  # ✅ Keep volumes
docker-compose down -v  # ❌ Remove volumes
```

### 2. Permission Denied Saat Upload

**Problem:** Container tidak bisa write ke volume

**Solution:**
```bash
# Pastikan directory ada dengan permission yang benar
docker exec api-go-prod mkdir -p /app/storage/avatars
docker exec api-go-prod chmod 755 /app/storage/avatars
```

### 3. Disk Space Issues

**Problem:** Volume terlalu besar

**Solution:**
```bash
# Check volume size
docker system df

# Find large files
docker exec api-go-prod find /app/storage/avatars -size +10M

# Manual cleanup (caution!)
docker run --rm -v avatars_data:/data \
  alpine find /data -name "avatar_*" -mtime +30 -delete
```

## Best Practices

### Development
1. ✅ Data otomatis hilang saat `docker-compose down -v` (clean slate)
2. ✅ Useful untuk testing upload/delete flows
3. ✅ Easy to reset: `docker volume rm avatars_data`

### Production
1. ✅ Always use named volumes (not bind mounts)
2. ✅ Regular backups: `tar czf avatars-$(date +%Y%m%d).tar.gz`
3. ✅ Monitor disk usage: `docker system df`
4. ✅ Implement cleanup policy untuk old avatars
5. ✅ Use separate storage for different environments
6. ✅ Consider external storage (S3) untuk scalability

## Migration from Static Serving

### Old Way (Not Recommended)
```
GET /storage/avatars/file.jpg  ❌
- No authorization
- Static file exposed
- Hard to track access
```

### New Way (Recommended)
```
GET /api/v1/avatars/1  ✅
- Controlled endpoint
- Can add authorization later
- Full audit trail possible
```

## Configuration Environment

### .env.example
```env
# Avatar Settings (future enhancement)
AVATAR_UPLOAD_ENABLED=true
AVATAR_MAX_SIZE=5242880  # 5MB in bytes
AVATAR_ALLOWED_TYPES=jpg,jpeg,png,gif,webp
AVATAR_STORAGE_PATH=/app/storage/avatars
AVATAR_PUBLIC_ACCESS=true  # Can be restricted per user
```

## Performance Optimization

### 1. Caching
```
Cache-Control: public, max-age=86400
- Cache 24 hours in browser
- Reduces server requests
- CDN-friendly headers
```

### 2. Image Optimization (Future)
```go
// Resize/compress on upload
- Original: 5MB JPEG
- Optimized: 200KB JPEG + 50KB thumbnail
- Save storage space
```

### 3. Lazy Loading (Frontend)
```javascript
<img loading="lazy" src="/api/v1/avatars/1" />
```

---

**Last Updated:** 2025-12-29
**Version:** 1.0
**Status:** Production Ready ✅
