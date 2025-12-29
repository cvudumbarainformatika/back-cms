# Configuration Verification Checklist

## Quick Verification

Run these commands untuk verify setup:

### 1. Check docker-compose.yml

```bash
# Verify all env vars set for dev
grep "STORAGE_" docker-compose.yml | grep environment -A 20

# Expected output:
#   AVATAR_STORAGE_PATH: /app/storage/avatars
#   STORAGE_THUMBNAILS_PATH: /app/storage/thumbnails
#   STORAGE_BERITA_PATH: /app/storage/berita
#   STORAGE_DOKUMEN_PATH: /app/storage/dokumen
#   STORAGE_GALERI_PATH: /app/storage/galeri
#   STORAGE_ATTACHMENT_PATH: /app/storage/attachments
```

```bash
# Verify all volumes mounted for dev
grep -A 10 "volumes:" docker-compose.yml | grep "storage"

# Expected output:
#   - avatars_data:/app/storage/avatars
#   - thumbnails_data:/app/storage/thumbnails
#   - berita_data:/app/storage/berita
#   - dokumen_data:/app/storage/dokumen
#   - galeri_data:/app/storage/galeri
#   - attachments_data:/app/storage/attachments
```

### 2. Check docker-compose.prod.yml

```bash
# Verify all env vars set for prod
grep "STORAGE_" docker-compose.prod.yml | grep environment -A 20

# Should have same 6 paths
```

```bash
# Verify all volumes mounted for prod
grep -A 15 "volumes:" docker-compose.prod.yml | grep "_prod_data"

# Expected output:
#   - avatars_prod_data:/app/storage/avatars
#   - thumbnails_prod_data:/app/storage/thumbnails
#   - berita_prod_data:/app/storage/berita
#   - dokumen_prod_data:/app/storage/dokumen
#   - galeri_prod_data:/app/storage/galeri
#   - attachments_prod_data:/app/storage/attachments
```

### 3. Check .env.example

```bash
# Count storage paths
grep "STORAGE_" .env.example | wc -l

# Should be at least 6 (or 7 with STORAGE_BASE_PATH)
```

---

## Complete Checklist

### File Types Configuration

**Avatar:**
- [ ] AVATAR_STORAGE_PATH in .env.example
- [ ] AVATAR_STORAGE_PATH in docker-compose.yml environment
- [ ] AVATAR_STORAGE_PATH in docker-compose.prod.yml environment
- [ ] avatars_data volume in docker-compose.yml
- [ ] avatars_prod_data volume in docker-compose.prod.yml

**Thumbnail:**
- [ ] STORAGE_THUMBNAILS_PATH in .env.example
- [ ] STORAGE_THUMBNAILS_PATH in docker-compose.yml environment
- [ ] STORAGE_THUMBNAILS_PATH in docker-compose.prod.yml environment
- [ ] thumbnails_data volume in docker-compose.yml
- [ ] thumbnails_prod_data volume in docker-compose.prod.yml

**Berita:**
- [ ] STORAGE_BERITA_PATH in .env.example
- [ ] STORAGE_BERITA_PATH in docker-compose.yml environment
- [ ] STORAGE_BERITA_PATH in docker-compose.prod.yml environment
- [ ] berita_data volume in docker-compose.yml
- [ ] berita_prod_data volume in docker-compose.prod.yml

**Dokumen:**
- [ ] STORAGE_DOKUMEN_PATH in .env.example
- [ ] STORAGE_DOKUMEN_PATH in docker-compose.yml environment
- [ ] STORAGE_DOKUMEN_PATH in docker-compose.prod.yml environment
- [ ] dokumen_data volume in docker-compose.yml
- [ ] dokumen_prod_data volume in docker-compose.prod.yml

**Galeri:**
- [ ] STORAGE_GALERI_PATH in .env.example
- [ ] STORAGE_GALERI_PATH in docker-compose.yml environment
- [ ] STORAGE_GALERI_PATH in docker-compose.prod.yml environment
- [ ] galeri_data volume in docker-compose.yml
- [ ] galeri_prod_data volume in docker-compose.prod.yml

**Attachment:**
- [ ] STORAGE_ATTACHMENT_PATH in .env.example
- [ ] STORAGE_ATTACHMENT_PATH in docker-compose.yml environment
- [ ] STORAGE_ATTACHMENT_PATH in docker-compose.prod.yml environment
- [ ] attachments_data volume in docker-compose.yml
- [ ] attachments_prod_data volume in docker-compose.prod.yml

---

## Code Implementation

### File Upload Service
- [ ] FileUploadType constants defined
- [ ] FileUploadConfigs map has all 6 types
- [ ] getStoragePathForType() function exists
- [ ] Each type has correct MaxSize
- [ ] Each type has correct AllowedTypes
- [ ] Each type has correct AllowedExts

### File Controller
- [ ] ServeFile() method implemented
- [ ] ListFiles() method implemented
- [ ] DeleteFile() method implemented
- [ ] Path traversal protection in place

### Routes
- [ ] FileController imported
- [ ] GET /files/:file_type/:filename route added
- [ ] GET /files/:file_type/list route added
- [ ] DELETE /files/:file_type/:filename route added (for future)

---

## Environment Variables Sync

### docker-compose.yml Environment Block

```yaml
environment:
  ✅ AVATAR_STORAGE_PATH: /app/storage/avatars
  ✅ STORAGE_THUMBNAILS_PATH: /app/storage/thumbnails
  ✅ STORAGE_BERITA_PATH: /app/storage/berita
  ✅ STORAGE_DOKUMEN_PATH: /app/storage/dokumen
  ✅ STORAGE_GALERI_PATH: /app/storage/galeri
  ✅ STORAGE_ATTACHMENT_PATH: /app/storage/attachments
```

### docker-compose.yml Volumes Block

```yaml
volumes:
  ✅ - avatars_data:/app/storage/avatars
  ✅ - thumbnails_data:/app/storage/thumbnails
  ✅ - berita_data:/app/storage/berita
  ✅ - dokumen_data:/app/storage/dokumen
  ✅ - galeri_data:/app/storage/galeri
  ✅ - attachments_data:/app/storage/attachments
```

### docker-compose.yml Named Volumes

```yaml
volumes:
  ✅ avatars_data:
  ✅ thumbnails_data:
  ✅ berita_data:
  ✅ dokumen_data:
  ✅ galeri_data:
  ✅ attachments_data:
```

### docker-compose.prod.yml Environment Block

```yaml
environment:
  ✅ AVATAR_STORAGE_PATH: /app/storage/avatars
  ✅ STORAGE_THUMBNAILS_PATH: /app/storage/thumbnails
  ✅ STORAGE_BERITA_PATH: /app/storage/berita
  ✅ STORAGE_DOKUMEN_PATH: /app/storage/dokumen
  ✅ STORAGE_GALERI_PATH: /app/storage/galeri
  ✅ STORAGE_ATTACHMENT_PATH: /app/storage/attachments
```

### docker-compose.prod.yml Volumes Block

```yaml
volumes:
  ✅ - avatars_prod_data:/app/storage/avatars
  ✅ - thumbnails_prod_data:/app/storage/thumbnails
  ✅ - berita_prod_data:/app/storage/berita
  ✅ - dokumen_prod_data:/app/storage/dokumen
  ✅ - galeri_prod_data:/app/storage/galeri
  ✅ - attachments_prod_data:/app/storage/attachments
```

### docker-compose.prod.yml Named Volumes

```yaml
volumes:
  ✅ avatars_prod_data:
  ✅ thumbnails_prod_data:
  ✅ berita_prod_data:
  ✅ dokumen_prod_data:
  ✅ galeri_prod_data:
  ✅ attachments_prod_data:
```

---

## Testing Checklist

### Unit Testing
- [ ] FileUploadService.Upload() for each type
- [ ] FileUploadService.Delete() for each type
- [ ] FileController.ServeFile() for each type
- [ ] File validation working for each type
- [ ] Size limits enforced per type

### Integration Testing
- [ ] Upload avatar via /api/v1/auth/profile
- [ ] Upload berita via /api/v1/berita
- [ ] Upload dokumen via /api/v1/dokumen
- [ ] Serve files via /api/v1/files/:type/:filename
- [ ] List files via /api/v1/files/:type/list

### Docker Testing
- [ ] docker-compose up -d works
- [ ] All 6 volumes mounted correctly
- [ ] Env vars passed to container
- [ ] Can read env vars inside container
- [ ] Files persist after container restart

### Production Testing
- [ ] docker-compose -f docker-compose.prod.yml up -d works
- [ ] All 6 prod volumes created
- [ ] Env vars in prod compose correct
- [ ] Files persist across container rebuilds
- [ ] Old files still accessible after deployment

---

## Common Issues & Fixes

### Issue: Env var not found
```bash
# Check if set
docker exec container_name env | grep STORAGE_

# If not found, check docker-compose.yml syntax
# (indentation, colons, etc.)
```

### Issue: Files lost after restart
```bash
# Check volume mount
docker inspect container_name | grep Mounts

# Should show:
# "Source": "docker volume name"
# "Destination": "/app/storage/xxx"
```

### Issue: Can't write to storage
```bash
# Check permissions
docker exec container_name ls -la /app/storage/

# Should show 755 permissions
# If not: docker exec container_name chmod 755 /app/storage/
```

---

## Quick Commands

```bash
# Verify all volumes
docker volume ls | grep avatars
docker volume ls | grep thumbnails
docker volume ls | grep berita
docker volume ls | grep dokumen
docker volume ls | grep galeri
docker volume ls | grep attachment

# Check volume data
docker run --rm -v avatars_data:/data alpine ls /data

# Backup all volumes
for type in avatars thumbnails berita dokumen galeri attachments; do
  docker run --rm -v ${type}_data:/data -v $(pwd):/backup \
    alpine tar czf /backup/${type}-$(date +%Y%m%d).tar.gz -C /data .
done

# Check docker-compose syntax
docker-compose config

# Check production compose syntax
docker-compose -f docker-compose.prod.yml config
```

---

## Summary

**Total items to verify: 30+**

Run this to count completed items:
```bash
grep "✅" VERIFICATION_CHECKLIST.md | wc -l
```

Should have:
- 5 checks per file type × 6 types = 30 ✅
- 3 code implementation items = 3 ✅
- **Total: 33+ items**

---

**When all items are checked ✅, you're production-ready!**
