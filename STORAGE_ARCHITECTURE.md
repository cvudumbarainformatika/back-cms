# Storage Architecture - Explaining Container vs VPS

## ❌ SALAH PAHAM

```
User mikir: "Storage ada di container"

❌ KELIRU!
```

---

## ✅ REALITAS

### Diagram Real Case:

```
VPS/Server (Linux)
├── /data/avatars/          ← Storage FISIK di VPS
│   ├── avatar_1_xxx.jpg
│   └── avatar_2_xxx.jpg
│
└── Docker
    └── Container (App)
        ├── /app/storage/avatars  ← Mount point (soft link)
        │   → Points to /data/avatars on VPS
        └── Process: Go App
            └── Write to /app/storage/avatars/
                └── Actually writes to /data/avatars/ on VPS
```

### Analogi:

```
SHORTCUT di Windows:
- Folder actual: C:\Users\Admin\Pictures
- Shortcut at: Desktop\Pictures → C:\Users\Admin\Pictures
- Buka shortcut = buka actual folder

DOCKER VOLUME:
- Folder actual: /data/avatars (VPS)
- Mount at: /app/storage/avatars (Container)
- Write to mount point = write to actual folder
```

---

## 🔄 Real Case: User Upload Foto

### Step by Step:

```
1️⃣  Nuxt Frontend (Browser)
    └─ User select image.jpg
    └─ POST /api/v1/auth/profile
    └─ Send file binary data
       
2️⃣  Go Backend Container
    ├─ Receive multipart request
    ├─ Validate file
    ├─ Generate filename: avatar_1_1735427370.jpg
    ├─ WRITE to /app/storage/avatars/avatar_1_1735427370.jpg
    │  ↓ (This is actually writing to VPS!)
    └─ Database: UPDATE users SET avatar = '/storage/avatars/avatar_1_1735427370.jpg'
    
3️⃣  Docker Volume Mount
    ├─ Container path: /app/storage/avatars
    ├─ VPS path: /data/avatars
    ├─ When container write to /app/storage/avatars/file.jpg
    └─ Actually saved at VPS: /data/avatars/file.jpg
    
4️⃣  Restart Container
    ├─ Container stop & removed
    ├─ New container start
    ├─ /app/storage/avatars mounted again to /data/avatars
    ├─ Old files still there! ✅
    ├─ Can read files immediately
    └─ GET /api/v1/avatars/1 → serves file
```

---

## ❌ KESALAHAN UMUM

### 1. Hardcode Path

```go
❌ SALAH (Current Implementation - but works for dev):
filePath := filepath.Join("storage/avatars", filename)

✅ BENAR:
storagePath := os.Getenv("STORAGE_PATH")
if storagePath == "" {
    storagePath = "./storage/avatars"  // fallback
}
filePath := filepath.Join(storagePath, filename)
```

### 2. Assumsi Storage di Container

```go
❌ SALAH:
"Saat container di-rebuild, file hilang"
→ Hanya benar jika tidak ada volume mount

✅ BENAR dengan volume:
- Container rebuild: file TETAP ada
- Karena file di VPS, bukan di container
- Volume hanya re-mount saja
```

### 3. Tidak Define Env Variable

```go
❌ SALAH - Hardcode:
avatarDir := "storage/avatars"

✅ BENAR - Environment:
avatarDir := os.Getenv("AVATAR_STORAGE_PATH")
if avatarDir == "" {
    avatarDir = "./storage/avatars"
}
```

---

## 🔧 Perbaikan yang Perlu Dilakukan

### 1. Update utils/file.go

```go
// BEFORE (Hardcode):
avatarDir := "storage/avatars"

// AFTER (Configurable):
avatarDir := os.Getenv("AVATAR_STORAGE_PATH")
if avatarDir == "" {
    avatarDir = "./storage/avatars"  // fallback untuk dev
}
```

### 2. Update AvatarController.go

```go
// BEFORE:
avatarDir := "storage/avatars"

// AFTER:
avatarDir := os.Getenv("AVATAR_STORAGE_PATH")
if avatarDir == "" {
    avatarDir = "./storage/avatars"
}
```

### 3. Add to .env.example

```env
# Storage Configuration
AVATAR_STORAGE_PATH=./storage/avatars
AVATAR_UPLOAD_MAX_SIZE=5242880
AVATAR_ALLOWED_TYPES=jpg,jpeg,png,gif,webp
```

### 4. Add to docker-compose.yml

```yaml
environment:
  AVATAR_STORAGE_PATH: /app/storage/avatars
  AVATAR_UPLOAD_ENABLED: true
```

### 5. Add to docker-compose.prod.yml

```yaml
environment:
  AVATAR_STORAGE_PATH: /app/storage/avatars
  AVATAR_UPLOAD_ENABLED: true
```

---

## 📊 Container vs VPS Storage

### Skenario 1: TANPA Volume

```
❌ Problem:
VPS /data/avatars/avatar_1.jpg
Container /app/storage/ → EMPTY

Why? Container memiliki filesystem sendiri
```

### Skenario 2: DENGAN Volume (Current)

```
✅ Correct:
VPS /data/avatars/avatar_1.jpg
Container /app/storage/ → LINK ke /data/avatars/

Docker Compose:
volumes:
  - /data/avatars:/app/storage/avatars

Result: File PERSISTEN ✅
```

### Skenario 3: Named Volume

```
✅ Also Correct:
Docker volume: avatars_data

docker-compose.yml:
volumes:
  - avatars_data:/app/storage/avatars

Result: File PERSISTEN ✅
```

---

## 🎯 Key Points untuk Diingat

### ✅ Yang Benar:

1. **Storage di VPS** - /data/avatars atau /var/lib/docker/volumes/...
2. **Container hanya akses via mount** - /app/storage/avatars
3. **File persisten** - Karena ada di VPS, bukan di container
4. **Rebuild container** - File tetap ada, volume hanya re-mount

### ❌ Yang Salah:

1. **Hardcode path** - Tidak fleksibel
2. **Anggap storage di container** - Akan hilang saat container stop
3. **Tidak backup volume** - File bisa hilang
4. **Tidak config env** - Sulit untuk dev vs prod

---

## 🚀 Correct Implementation Path

```
USER (Nuxt Browser)
    ↓
    POST /api/v1/auth/profile
    {file: image.jpg}
    ↓
Go Backend
    ├─ Read env: AVATAR_STORAGE_PATH = /app/storage/avatars
    ├─ Create filename: avatar_1_1735427370.jpg
    ├─ Create fullpath: /app/storage/avatars/avatar_1_1735427370.jpg
    ├─ Write file
    │  (Docker internally routes to /data/avatars/avatar_1_1735427370.jpg)
    ├─ Save to DB: avatar = /storage/avatars/avatar_1_1735427370.jpg
    └─ Return response
       
    ↓
Docker Volume Mount
    ├─ Container: /app/storage/avatars
    ├─ VPS/Host: /data/avatars (or managed volume)
    └─ File actually saved at: /data/avatars/avatar_1_1735427370.jpg
    
    ↓
VPS Filesystem
    └─ PHYSICAL FILE ✅
       /data/avatars/avatar_1_1735427370.jpg
       - Persisten selamanya
       - Survive container restart
       - Can be backed up
       - Can be shared with other containers
```

---

## 📝 Configuration Strategy

### Development (.env):
```env
AVATAR_STORAGE_PATH=./storage/avatars
# Local relative path - works with docker volume
```

### Production (.env.production):
```env
AVATAR_STORAGE_PATH=/app/storage/avatars
# Absolute path in container, but mounted from VPS /data/avatars
```

### Both environments:
```
✅ Same code
✅ Different paths via env
✅ Volume mount handles the rest
```

---

## ✅ Takeaways

1. **Storage fisik selalu di VPS** ✅
2. **Container cuma akses via mount point** ✅
3. **File persisten karena di VPS, bukan container** ✅
4. **Gunakan env variable, jangan hardcode** ✅
5. **Volume mount adalah "soft link" ke VPS storage** ✅

**Kalau lupa:** Container adalah tempat app jalan, BUKAN tempat data disimpan!
