# Common Mistakes & Fixes - Avatar Storage

## ❌ Kesalahan Umum yang Sering Dilakukan

### 1. Hardcode Path di Go Code

**❌ SALAH (Before):**
```go
avatarDir := "storage/avatars"  // Hardcoded!
filePath := filepath.Join(avatarDir, filename)
```

**✅ BENAR (After):**
```go
avatarDir := utils.GetStoragePath()  // From environment
filePath := filepath.Join(avatarDir, filename)
```

### 2. Tidak Distinguish antara Container Path dan VPS Path

**❌ SALAH - Paham:**
"Berarti storage ada di container ya?"

**✅ BENAR - Realitas:**
- Storage fisik: **di VPS** (`/data/avatars`)
- Container akses via: **mount point** (`/app/storage/avatars`)
- Container bukan tempat menyimpan data!

### 3. Hardcode di Multiple Places

**❌ SALAH:**
```
utils/file.go:        avatarDir := "storage/avatars"
AvatarController.go:  avatarDir := "storage/avatars"
Lagi di tempat lain:  avatarDir := "storage/avatars"
```

**✅ BENAR:**
```go
// Satu tempat saja di utils/file.go
func GetStoragePath() string {
    storagePath := os.Getenv("AVATAR_STORAGE_PATH")
    if storagePath == "" {
        storagePath = "./storage/avatars"
    }
    return storagePath
}

// Digunakan di mana-mana
avatarDir := utils.GetStoragePath()
```

### 4. Tidak Set Environment Variable

**❌ SALAH:**
```bash
# Jalankan tanpa env var
./backend

# App gunakan hardcode "./storage/avatars"
# Padahal di Docker harus "/app/storage/avatars"
```

**✅ BENAR:**
```bash
# docker-compose.yml
environment:
  AVATAR_STORAGE_PATH: /app/storage/avatars

# Backend baca dari env
storagePath := os.Getenv("AVATAR_STORAGE_PATH")
```

---

## 🔧 What We Fixed

### 1. Created GetStoragePath() Helper

```go
func GetStoragePath() string {
    storagePath := os.Getenv("AVATAR_STORAGE_PATH")
    if storagePath == "" {
        storagePath = "./storage/avatars"  // fallback
    }
    return storagePath
}
```

**Used in:**
- `utils/file.go` - UploadAvatar()
- `utils/file.go` - DeleteAvatar()
- `app/Http/Controllers/AvatarController.go` - GetAvatar()
- `app/Http/Controllers/AvatarController.go` - GetAvatarByName()

### 2. Added Environment Variables

**`.env.example`:**
```env
AVATAR_STORAGE_PATH=./storage/avatars
AVATAR_UPLOAD_ENABLED=true
AVATAR_MAX_SIZE=5242880
AVATAR_ALLOWED_TYPES=jpg,jpeg,png,gif,webp
```

**`docker-compose.yml`:**
```yaml
environment:
  AVATAR_STORAGE_PATH: /app/storage/avatars
```

**`docker-compose.prod.yml`:**
```yaml
environment:
  AVATAR_STORAGE_PATH: /app/storage/avatars
```

### 3. Improved Path Handling

**Before (Wrong):**
```go
avatarDir := "storage/avatars"
```

**After (Correct):**
```go
avatarDir := utils.GetStoragePath()
// Returns based on environment:
// - Dev: ./storage/avatars
// - Prod: /app/storage/avatars
```

---

## 📊 Real Case Walkthrough

### User Upload Flow (CORRECT NOW):

```
1️⃣  Nuxt Frontend
    └─ POST /api/v1/auth/profile with file
    
2️⃣  Go Backend - UpdateProfile()
    ├─ Get AVATAR_STORAGE_PATH from env
    ├─ Get path: AVATAR_STORAGE_PATH = /app/storage/avatars
    ├─ Create filename: avatar_1_1735427370.jpg
    ├─ Call UploadAvatar(file, userID)
    
3️⃣  UploadAvatar() in utils/file.go
    ├─ avatarDir := utils.GetStoragePath()
    ├─ // Returns /app/storage/avatars
    ├─ MkdirAll(avatarDir, 0755)
    ├─ Create file: /app/storage/avatars/avatar_1_1735427370.jpg
    └─ Return path: /storage/avatars/avatar_1_1735427370.jpg
    
4️⃣  Docker Volume Mount
    ├─ Container: /app/storage/avatars
    ├─ VPS: /data/avatars (via volume mount)
    ├─ File written to container path
    ├─ Docker mounts translate it to VPS path
    └─ Physical file: /data/avatars/avatar_1_1735427370.jpg ✅
    
5️⃣  Database
    └─ UPDATE users SET avatar = /storage/avatars/avatar_1_1735427370.jpg
    
6️⃣  Later - Container Restart
    ├─ Container deleted
    ├─ New container created
    ├─ Docker volume re-mounted
    ├─ /app/storage/avatars → /data/avatars
    ├─ Old files still there! ✅
    ├─ GET /api/v1/avatars/1
    └─ File served successfully ✅
```

---

## ✅ Verification Checklist

- [ ] `GetStoragePath()` exists in `utils/file.go`
- [ ] Used in `UploadAvatar()` and `DeleteAvatar()`
- [ ] Used in `AvatarController.GetAvatar()` and `GetAvatarByName()`
- [ ] `AVATAR_STORAGE_PATH` added to `.env.example`
- [ ] `docker-compose.yml` sets `AVATAR_STORAGE_PATH: /app/storage/avatars`
- [ ] `docker-compose.prod.yml` sets `AVATAR_STORAGE_PATH: /app/storage/avatars`
- [ ] Volume mounted in both compose files
- [ ] No hardcoded `"storage/avatars"` strings remain

---

## 🚀 Testing

### Development
```bash
# Start services
docker-compose up -d

# Build and run
go build -o backend main.go
AVATAR_STORAGE_PATH=./storage/avatars ./backend

# Upload
curl -X PUT http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer TOKEN" \
  -F "avatar=@image.jpg"

# Access
curl http://localhost:8080/api/v1/avatars/1
```

### Production
```bash
# Start services
docker-compose -f docker-compose.prod.yml up -d

# Verify env var
docker exec api-go-prod echo $AVATAR_STORAGE_PATH
# Output: /app/storage/avatars

# Check storage
docker exec api-go-prod ls -lah /app/storage/avatars/
```

---

## 📝 Key Takeaways

1. **Never hardcode paths** - Use environment variables
2. **Centralize path logic** - `GetStoragePath()` is the single source of truth
3. **Dev ≠ Prod** - Different paths but same code
4. **Volume is the bridge** - Connects container to VPS storage
5. **Fallback is safe** - If env not set, use sensible default

---

## 🔗 Related Files

- `AVATAR_STORAGE_GUIDE.md` - Full architecture explanation
- `STORAGE_ARCHITECTURE.md` - Container vs VPS detailed explanation
- `.env.example` - Configuration reference
- `docker-compose.yml` - Development setup
- `docker-compose.prod.yml` - Production setup

---

**Status:** ✅ All common mistakes fixed and documented
