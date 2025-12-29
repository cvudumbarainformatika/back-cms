# Storage & Git Ignore Guide

## ❓ Pertanyaan: Apakah storage/ disimpan ke git?

**Jawaban:** ❌ **TIDAK** - storage/ sudah di-ignore oleh `.gitignore`

---

## 📋 .gitignore Configuration

### Storage Directories

```gitignore
# Storage directories (user uploaded files)
storage/
!storage/.gitkeep
```

**Apa artinya:**
- `storage/` - Ignore semua file di storage folder
- `!storage/.gitkeep` - TAPI, jangan ignore file `.gitkeep`

---

## 📁 Storage Structure

```
storage/
├── .gitkeep                    # Committed to git (folder marker)
├── avatars/
│   ├── .gitkeep               # Committed to git (folder marker)
│   ├── avatar_1_1767011202.png  # ❌ NOT committed (user file)
│   └── avatar_2_1767011202.jpg  # ❌ NOT committed (user file)
├── thumbnails/
│   ├── .gitkeep               # Committed to git
│   └── thumbnail_1_xxx.jpg    # ❌ NOT committed
├── berita/
│   ├── .gitkeep               # Committed to git
│   └── berita_5_xxx.jpg       # ❌ NOT committed
├── dokumen/
│   ├── .gitkeep               # Committed to git
│   └── dokumen_1_xxx.pdf      # ❌ NOT committed
├── galeri/
│   ├── .gitkeep               # Committed to git
│   └── galeri_3_xxx.jpg       # ❌ NOT committed
└── attachments/
    ├── .gitkeep               # Committed to git
    └── attachment_1_xxx.zip   # ❌ NOT committed
```

---

## ✅ What Gets Committed

### Committed to Git:
- ✅ `storage/` (directory)
- ✅ `storage/.gitkeep` (folder marker)
- ✅ `storage/avatars/.gitkeep` (folder marker)
- ✅ `storage/thumbnails/.gitkeep` (folder marker)
- ✅ `storage/berita/.gitkeep` (folder marker)
- ✅ `storage/dokumen/.gitkeep` (folder marker)
- ✅ `storage/galeri/.gitkeep` (folder marker)
- ✅ `storage/attachments/.gitkeep` (folder marker)

### NOT Committed to Git:
- ❌ `storage/avatars/avatar_1_xxx.png` (user uploaded)
- ❌ `storage/berita/berita_5_xxx.jpg` (user uploaded)
- ❌ `storage/dokumen/dokumen_1_xxx.pdf` (user uploaded)
- ❌ Any actual uploaded files

---

## 🎯 Why .gitkeep?

Git tidak track empty directories. `.gitkeep` adalah trick untuk:

1. **Preserve folder structure** - Folder tetap ada meski kosong
2. **Enable cloning** - Clone repo, storage folders sudah ada
3. **Ready for uploads** - User uploads langsung bisa berjalan

### Tanpa .gitkeep:
```bash
git clone repo
# storage/ folder tidak ada!
# Upload akan error: directory doesn't exist
```

### Dengan .gitkeep:
```bash
git clone repo
# storage/ folder sudah ada
# All subdirectories ready
# Upload berjalan lancar ✅
```

---

## 🧪 Verification

### Check git status
```bash
# Should NOT show avatar_1_xxx.png files
git status

# Should ONLY show .gitkeep files
# Output example:
# On branch main
# nothing to commit, working tree clean
```

### Check what's ignored
```bash
# See all ignored files in storage/
git check-ignore -v storage/avatars/*

# Output:
# storage/avatars/avatar_1_1767011202.png ignore via storage/
# storage/avatars/avatar_2_xxx.jpg ignore via storage/
```

### List only committed storage files
```bash
git ls-files storage/

# Should ONLY show:
# storage/.gitkeep
# storage/avatars/.gitkeep
# storage/berita/.gitkeep
# storage/dokumen/.gitkeep
# storage/galeri/.gitkeep
# storage/thumbnails/.gitkeep
# storage/attachments/.gitkeep
```

---

## 🚀 Setup for New Clone

### After git clone:
```bash
# 1. Repository cloned
git clone your-repo.git
cd your-repo

# 2. Storage folders already exist (via .gitkeep)
ls -la storage/
# Output:
# total 56
# drwxr-xr-x   9 user  staff   288 Dec 29 19:58 .
# drwxr-xr-x  63 user  staff  2016 Dec 29 20:10 ..
# drwxr-xr-x   3 user  staff    96 Dec 29 20:10 avatars
# drwxr-xr-x   3 user  staff    96 Dec 29 20:10 berita
# ... etc

# 3. Ready to start server
docker-compose up -d
go run main.go

# 4. Users can upload immediately ✅
```

---

## 💾 Backup Strategy

### For Production:

**Git Repository (Committed):**
- Code files ✅
- Configuration ✅
- .gitkeep markers ✅

**Separate Backup (NOT in git):**
- User uploaded files ❌ (too large)
- Database backups ❌ (too large)
- Storage volumes ❌ (too large)

**Backup Commands:**
```bash
# Backup user uploads
docker run --rm -v avatars_data:/data -v $(pwd):/backup \
  alpine tar czf /backup/avatars-$(date +%Y%m%d).tar.gz -C /data .

# Backup all storage
for type in avatars thumbnails berita dokumen galeri attachments; do
  docker run --rm -v ${type}_data:/data -v $(pwd):/backup \
    alpine tar czf /backup/${type}-$(date +%Y%m%d).tar.gz -C /data .
done

# Backup database
docker exec mysql mysqldump -uadmin -psasa0102 sasacms > backup-$(date +%Y%m%d).sql
```

---

## ⚠️ Common Mistakes

### ❌ Mistake 1: Accidentally commit storage files
```bash
git add -A
git commit -m "Add all files"
# ❌ Storage files might get added!

# Fix: Use git add with specific paths
git add app/ routes/ config/
git add docker-compose.yml .env.example
# ✅ Never accidentally add storage/
```

### ❌ Mistake 2: Remove .gitkeep files
```bash
rm storage/avatars/.gitkeep
git add -A
git commit -m "cleanup"
# ❌ When cloned, folders won't exist!

# Fix: Always keep .gitkeep files
git checkout storage/avatars/.gitkeep
```

### ❌ Mistake 3: Wrong .gitignore pattern
```gitignore
❌ storage  (without slash - might match file named "storage")
❌ /storage/avatars/*  (too specific, gitkeep won't work)

✅ storage/  (directory)
✅ !storage/.gitkeep  (exception)
```

---

## 🔄 Workflow

### Developer 1: Clone and upload
```bash
git clone repo
cd repo
docker-compose up -d
go run main.go

# Upload avatar
curl -X PUT .../auth/profile -F "avatar=@myavatar.jpg"
# File saved to: storage/avatars/avatar_1_xxx.jpg
```

### Developer 2: Clone and can't see the file
```bash
git clone repo
cd repo
# storage/avatars/avatar_1_xxx.jpg doesn't exist ✅ (by design)
# But storage/avatars/ folder exists ✅ (via .gitkeep)
# Each environment has its own uploads
```

### Production: Backup and restore
```bash
# Backup uploads
tar czf uploads-backup.tar.gz storage/

# Deploy new version
git pull
docker-compose pull
docker-compose up -d

# Restore uploads
tar xzf uploads-backup.tar.gz
```

---

## 📊 Summary Table

| Item | Git Committed | Location |
|------|---------------|----------|
| storage/ folder | ✅ Yes | repository |
| .gitkeep files | ✅ Yes | repository |
| User uploads | ❌ No | Docker volume / VPS |
| Code files | ✅ Yes | repository |
| Configuration | ✅ Yes (as .example) | repository |

---

## ✅ Best Practices

1. **Always keep .gitkeep** - Preserves folder structure
2. **Never commit user uploads** - Use separate backup
3. **Use gitignore pattern** - `storage/` + `!storage/.gitkeep`
4. **Backup separately** - User files ≠ code
5. **Document backup strategy** - In case of disaster recovery
6. **Test on fresh clone** - Make sure .gitkeep works

---

## 🚨 Emergency: Storage Lost

If storage folders somehow deleted:

```bash
# Recreate structure
mkdir -p storage/{avatars,thumbnails,berita,dokumen,galeri,attachments}

# Add .gitkeep markers
for dir in storage/avatars storage/thumbnails storage/berita \
           storage/dokumen storage/galeri storage/attachments; do
  touch "$dir/.gitkeep"
done

# Restore from backup
tar xzf uploads-backup.tar.gz

# Verify
ls -la storage/avatars/
```

---

**Status:** ✅ Storage properly ignored, .gitkeep preserved!
