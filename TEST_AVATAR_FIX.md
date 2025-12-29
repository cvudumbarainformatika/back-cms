# Test Avatar File Fix

## 🎯 What Was Fixed

**Problem:** File baru tidak tersimpan ke disk
- Database: `/api/v1/files/avatar/avatar_1_1767015536.jpeg` ✅
- Storage: `avatar_1_1767011202.png` ❌ (old file)

**Root Cause:** `DeleteAvatar()` tidak bisa parse API URL format `/api/v1/files/avatar/...`

**Solution:** Improve `DeleteAvatar()` untuk handle multiple path formats:
- `/api/v1/files/avatar/avatar_1_xxx.jpg` (NEW format)
- `/storage/avatars/avatar_1_xxx.jpg` (OLD format)
- `avatar_1_xxx.jpg` (Filename only)

---

## 🚀 Test Steps

### Step 1: Restart Backend

```bash
pkill -9 -f "./backend" || true
sleep 1
./backend &
sleep 2
```

### Step 2: Get Token

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@pdpi.co.id",
    "password": "AdminPDPI2024!"
  }' | jq -r '.data.access_token')

echo "Token: $TOKEN"
```

### Step 3: Upload New Avatar

```bash
curl -X PUT http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer $TOKEN" \
  -F "name=Admin Pusat" \
  -F "phone=081237660656" \
  -F "avatar=@/path/to/new-image.jpg"
```

Expected response:
```json
{
  "success": true,
  "data": {
    "avatar": "/api/v1/files/avatar/avatar_1_NEW_TIMESTAMP.jpg"
  }
}
```

### Step 4: Verify File Saved

```bash
ls -lah storage/avatars/

# Should show BOTH:
# avatar_1_1767011202.png (OLD - should be deleted)
# avatar_1_NEW_TIMESTAMP.jpg (NEW - just uploaded)
```

**Expected:** Only NEW file should exist
- ✅ avatar_1_NEW_TIMESTAMP.jpg (from new upload)
- ❌ avatar_1_1767011202.png (should be deleted)

### Step 5: Access File via API

```bash
curl http://localhost:8080/api/v1/files/avatar/avatar_1_NEW_TIMESTAMP.jpg \
  -o /tmp/downloaded.jpg

file /tmp/downloaded.jpg
# Should show: image/jpeg or image/png
```

### Step 6: Verify Database

```bash
docker exec starter-mysql mysql -uadmin -psasa0102 sasacms -e \
"SELECT id, avatar FROM users WHERE id = 1;"

# Should show:
# id | avatar
# 1  | /api/v1/files/avatar/avatar_1_NEW_TIMESTAMP.jpg
```

---

## 🔍 Debug Logs

Check backend logs for DeleteAvatar operations:

```bash
# If running in foreground, look for:
# [Avatar] Converting API URL to storage path: /api/v1/files/avatar/avatar_1_OLD.png → ./storage/avatars/avatar_1_OLD.png
# [Avatar] Deleted successfully: ./storage/avatars/avatar_1_OLD.png
```

These logs confirm:
1. Old avatar path recognized ✅
2. Converted to correct storage path ✅
3. Old file successfully deleted ✅

---

## ✅ Expected Final State

After successful upload:

```
Database:
  avatar = /api/v1/files/avatar/avatar_1_NEW_TIMESTAMP.jpg ✅

Storage:
  storage/avatars/
  ├── .gitkeep
  └── avatar_1_NEW_TIMESTAMP.jpg ✅ (only new file)

File Access:
  curl http://localhost:8080/api/v1/files/avatar/avatar_1_NEW_TIMESTAMP.jpg
  → Returns image file ✅
```

---

## 🚨 If Still Not Working

### Check 1: Old file not deleted
```bash
ls -la storage/avatars/
# If old file still exists, DeleteAvatar failed silently
# Check backend logs for error messages
```

### Check 2: File not created
```bash
ls -la storage/avatars/avatar_1_*
# If no new file, UploadAvatar failed
# Check backend response for error details
```

### Check 3: Wrong storage path
```bash
grep "Using storage path" /tmp/backend.log
# Should show: Using storage path: ./storage/avatars
```

---

## 📊 Path Format Handling

DeleteAvatar now handles all these formats:

```
Input: /api/v1/files/avatar/avatar_1_1767015536.jpeg
→ Extract: avatar_1_1767015536.jpeg
→ Join with storage: ./storage/avatars/avatar_1_1767015536.jpeg ✅

Input: /storage/avatars/avatar_1_1767011202.png
→ Extract: avatar_1_1767011202.png
→ Join with storage: ./storage/avatars/avatar_1_1767011202.png ✅

Input: avatar_1_1767015536.jpeg
→ Join directly: ./storage/avatars/avatar_1_1767015536.jpeg ✅
```

---

**After successful test, old database data still needs migration (separate task).** ✅
