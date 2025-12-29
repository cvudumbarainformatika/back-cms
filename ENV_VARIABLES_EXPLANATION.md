# Environment Variables Explanation

## ❌ Mistake You Caught

```yaml
❌ WRONG (Incomplete):
environment:
  AVATAR_STORAGE_PATH: /app/storage/avatars
  # ... volumes untuk semua type tapi hanya 1 env var!

volumes:
  - avatars_data:/app/storage/avatars
  - thumbnails_data:/app/storage/thumbnails
  - berita_data:/app/storage/berita
  - dokumen_data:/app/storage/dokumen
  - galeri_data:/app/storage/galeri
  - attachments_data:/app/storage/attachments
```

**Problem:** 
- Volumes mounted untuk semua types ✅
- Tapi Go app hanya tahu path untuk avatar ❌
- Untuk thumbnail, berita, dll → Go app akan error atau create di default path ❌

---

## ✅ Correct Approach

```yaml
✅ CORRECT (Complete):
environment:
  # Storage paths for ALL file types
  AVATAR_STORAGE_PATH: /app/storage/avatars
  STORAGE_THUMBNAILS_PATH: /app/storage/thumbnails
  STORAGE_BERITA_PATH: /app/storage/berita
  STORAGE_DOKUMEN_PATH: /app/storage/dokumen
  STORAGE_GALERI_PATH: /app/storage/galeri
  STORAGE_ATTACHMENT_PATH: /app/storage/attachments

volumes:
  - avatars_data:/app/storage/avatars
  - thumbnails_data:/app/storage/thumbnails
  - berita_data:/app/storage/berita
  - dokumen_data:/app/storage/dokumen
  - galeri_data:/app/storage/galeri
  - attachments_data:/app/storage/attachments
```

**Why:**
- 1 env var per volume ✅
- Go app tahu path untuk setiap type ✅
- Scalable: tambah type → tambah 1 env var + 1 volume ✅

---

## How It Works

### Docker-compose side:
```yaml
volumes:
  - avatars_data:/app/storage/avatars
  └─ Mounting named volume 'avatars_data' 
     to container path '/app/storage/avatars'
```

### Go App side:
```go
// file_upload_service.go
func getStoragePathForType(fileType string) string {
    envVar := fmt.Sprintf("STORAGE_%s_PATH", strings.ToUpper(fileType))
    // For "berita": looks for STORAGE_BERITA_PATH
    storagePath := os.Getenv(envVar)
    if storagePath == "" {
        // Fallback if env not set
        baseStorage := GetStoragePath()
        storagePath = filepath.Join(baseStorage, fileType)
    }
    return storagePath
}
```

### Flow:
```
1. Docker compose sets env var:
   STORAGE_BERITA_PATH=/app/storage/berita

2. Go app reads env var:
   os.Getenv("STORAGE_BERITA_PATH")
   → Returns: /app/storage/berita

3. Go app uses path:
   filePath := filepath.Join(storagePath, filename)
   → /app/storage/berita/berita_1_xxx.jpg

4. Docker volume mount:
   Container /app/storage/berita
   ← mounted from → VPS /data/berita (managed volume)

5. File saved:
   /app/storage/berita/berita_1_xxx.jpg
   = /data/berita/berita_1_xxx.jpg (on VPS)
```

---

## Environment Variable Naming Convention

```
STORAGE_{TYPE}_PATH

Avatar:     AVATAR_STORAGE_PATH  (special case, legacy name)
Thumbnail:  STORAGE_THUMBNAILS_PATH
Berita:     STORAGE_BERITA_PATH
Dokumen:    STORAGE_DOKUMEN_PATH
Galeri:     STORAGE_GALERI_PATH
Attachment: STORAGE_ATTACHMENT_PATH
```

**Why different names?**
- `AVATAR_STORAGE_PATH` - untuk backward compatibility
- `STORAGE_*_PATH` - pattern untuk new types
- Konsisten dan mudah diidentifikasi

---

## Complete List (Dev vs Prod)

### Development (.env)
```env
# Base path (fallback)
STORAGE_BASE_PATH=./storage

# Individual types (optional in dev, used if set)
AVATAR_STORAGE_PATH=./storage/avatars
STORAGE_THUMBNAILS_PATH=./storage/thumbnails
STORAGE_BERITA_PATH=./storage/berita
STORAGE_DOKUMEN_PATH=./storage/dokumen
STORAGE_GALERI_PATH=./storage/galeri
STORAGE_ATTACHMENT_PATH=./storage/attachments
```

### Production (docker-compose.prod.yml)
```yaml
environment:
  AVATAR_STORAGE_PATH: /app/storage/avatars
  STORAGE_THUMBNAILS_PATH: /app/storage/thumbnails
  STORAGE_BERITA_PATH: /app/storage/berita
  STORAGE_DOKUMEN_PATH: /app/storage/dokumen
  STORAGE_GALERI_PATH: /app/storage/galeri
  STORAGE_ATTACHMENT_PATH: /app/storage/attachments
```

---

## Why Both Must Match

### Scenario 1: Volume tapi no env var
```yaml
environment:
  AVATAR_STORAGE_PATH: /app/storage/avatars  ✅
  # Missing: STORAGE_BERITA_PATH ❌

volumes:
  - berita_data:/app/storage/berita  ✅
```

**Result:**
- Go app uploads berita
- Looks for STORAGE_BERITA_PATH → not found
- Fallback: ./storage/berita
- File saved to: ./storage/berita/ (relative to working dir)
- Not saved to: /app/storage/berita (where volume is mounted)
- File lost when container restarts! ❌

### Scenario 2: Env var tapi no volume
```yaml
environment:
  STORAGE_BERITA_PATH: /app/storage/berita  ✅

volumes:
  # Missing: berita_data:/app/storage/berita ❌
```

**Result:**
- Go app uploads berita
- Reads STORAGE_BERITA_PATH = /app/storage/berita
- Creates /app/storage/berita/file.jpg
- File saved to container filesystem (not mounted)
- File lost when container stops! ❌

### Scenario 3: Both match correctly ✅
```yaml
environment:
  STORAGE_BERITA_PATH: /app/storage/berita  ✅

volumes:
  - berita_data:/app/storage/berita  ✅
```

**Result:**
- Go app uploads berita to /app/storage/berita/file.jpg
- Docker mounts volume: berita_data ← /app/storage/berita
- File persisted on VPS (in managed volume) ✅
- File survives container restart ✅

---

## Verification Checklist

For each file type, verify:

```
Type: avatar
  ✅ Environment variable in docker-compose.yml
  ✅ Environment variable in docker-compose.prod.yml
  ✅ Volume mount in docker-compose.yml
  ✅ Volume mount in docker-compose.prod.yml
  ✅ Volume definition in volumes section
  ✅ Env var name matches pattern: {TYPE}_STORAGE_PATH or STORAGE_{TYPE}_PATH
  ✅ Path in env var matches path in volume mount

Type: thumbnail
  ... (same checks)

Type: berita
  ... (same checks)

Type: dokumen
  ... (same checks)

Type: galeri
  ... (same checks)

Type: attachment
  ... (same checks)
```

---

## Example: Adding New File Type "video"

### 1. Add to docker-compose.yml

```yaml
environment:
  # ... existing vars ...
  STORAGE_VIDEO_PATH: /app/storage/videos  # ← ADD THIS

volumes:
  # ... existing volumes ...
  - videos_data:/app/storage/videos  # ← ADD THIS
```

### 2. Add to volumes section

```yaml
volumes:
  # ... existing volumes ...
  videos_data:  # ← ADD THIS
```

### 3. Add to docker-compose.prod.yml

```yaml
environment:
  # ... existing vars ...
  STORAGE_VIDEO_PATH: /app/storage/videos  # ← ADD THIS

volumes:
  # ... existing volumes ...
  - videos_prod_data:/app/storage/videos  # ← ADD THIS
```

### 4. Add to volumes section

```yaml
volumes:
  # ... existing volumes ...
  videos_prod_data:  # ← ADD THIS
```

### 5. Add to file_upload_service.go

```go
FileUploadConfigs[FileTypeVideo] = FileUploadConfig{
    StoragePath: getStoragePathForType("videos"),
    // ... other config
}
```

### That's it! ✅

Now your app will:
- Read STORAGE_VIDEO_PATH from env
- Use /app/storage/videos for uploads
- Mount to videos_data volume
- Persist files across rebuilds

---

## Key Takeaway

✅ **Rule:** 1 Environment Variable per Volume

```
1 volume mount ← needs → 1 env variable
```

If you forget the env var:
- App doesn't know where to save files
- Falls back to default/relative paths
- Files get lost on container restart

**Always keep them in sync!**

---

## Files to Check

```
✅ .env.example - All env vars defined
✅ docker-compose.yml - All env vars set + all volumes mounted
✅ docker-compose.prod.yml - All env vars set + all volumes mounted
✅ file_upload_service.go - getStoragePathForType() uses env vars
```

---

**Thank you for catching this! Attention to detail is critical for production systems!** 🎯
