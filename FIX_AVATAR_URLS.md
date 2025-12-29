# Fix Avatar URLs in Database

## 🎯 Problem

Avatar URLs dalam database tersimpan dengan format lama:

```
❌ /storage/avatars/avatar_1_1767013716.png
```

Seharusnya format baru (API endpoint):

```
✅ /api/v1/files/avatar/avatar_1_1767013716.png
```

---

## ✅ Solution Applied

### Code Changes:

1. **utils/file.go** - `UploadAvatar()` sekarang return API URL:
```go
// Before: /storage/avatars/avatar_1_xxx.png
// After:  /api/v1/files/avatar/avatar_1_xxx.png
apiURL := fmt.Sprintf("/api/v1/files/avatar/%s", filename)
return apiURL, nil
```

2. **AuthController.go** - `UpdateProfile()` dan `Me()` return correct URL

### Result:
- ✅ New uploads akan automatically use correct URL
- ❌ Old uploads in database still have old URL format

---

## 🔧 Fix Existing Data

### SQL Query to Update All Avatar URLs

```sql
UPDATE users 
SET avatar = CONCAT('/api/v1/files/avatar/', SUBSTRING_INDEX(avatar, '/', -1))
WHERE avatar LIKE '/storage/avatars/%';
```

**Penjelasan:**
- `LIKE '/storage/avatars/%'` - Find all old format URLs
- `SUBSTRING_INDEX(avatar, '/', -1)` - Extract filename (last part after /)
- `CONCAT('/api/v1/files/avatar/', ...)` - Build new URL

---

## 🧪 Verify Before Update

```sql
-- See what will be updated
SELECT id, email, avatar FROM users WHERE avatar LIKE '/storage/avatars/%';

-- Example output:
-- id | email                  | avatar
-- 1  | admin@pdpi.co.id       | /storage/avatars/avatar_1_1767013716.png
```

---

## 🚀 Steps to Apply Fix

### Step 1: Backup database
```bash
docker exec starter-mysql mysqldump -uadmin -psasa0102 sasacms > backup-before-fix.sql
```

### Step 2: Run update query
```bash
docker exec starter-mysql mysql -uadmin -psasa0102 sasacms -e "UPDATE users SET avatar = CONCAT('/api/v1/files/avatar/', SUBSTRING_INDEX(avatar, '/', -1)) WHERE avatar LIKE '/storage/avatars/%';"
```

### Step 3: Verify update
```bash
docker exec starter-mysql mysql -uadmin -psasa0102 sasacms -e "SELECT id, email, avatar FROM users WHERE avatar NOT NULL;"
```

### Step 4: Test API
```bash
# Get user profile
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer TOKEN"

# Should return:
# "avatar": "/api/v1/files/avatar/avatar_1_1767013716.png"
```

---

## 📝 Complete SQL Script

```sql
-- Backup old URLs (optional)
CREATE TABLE avatar_urls_backup AS 
SELECT id, email, avatar FROM users WHERE avatar IS NOT NULL;

-- Update to new format
UPDATE users 
SET avatar = CONCAT('/api/v1/files/avatar/', SUBSTRING_INDEX(avatar, '/', -1))
WHERE avatar LIKE '/storage/avatars/%';

-- Verify
SELECT id, email, avatar FROM users WHERE avatar IS NOT NULL;
```

---

## ✅ Verification Checklist

After running the update:

- [ ] All avatars starting with `/storage/avatars/` are updated
- [ ] New format: `/api/v1/files/avatar/avatar_X_XXX.ext`
- [ ] Files still exist in storage directory
- [ ] API endpoint `/api/v1/files/avatar/avatar_X_XXX.png` works
- [ ] Frontend can load images with new URLs

---

## 🔄 Going Forward

### New Uploads
- ✅ Automatically use correct URL format
- ✅ Stored as: `/api/v1/files/avatar/{filename}`

### Frontend
- ✅ Can use avatar URL directly from database
- ✅ No need to transform URLs

### Example Response
```json
{
  "id": 1,
  "name": "Admin Pusat",
  "email": "admin@pdpi.co.id",
  "avatar": "/api/v1/files/avatar/avatar_1_1767013716.png",
  "phone": "081237660656",
  "address": "Jakarta"
}
```

---

## 🚨 Rollback (If Something Goes Wrong)

```bash
# Restore from backup
docker exec starter-mysql mysql -uadmin -psasa0102 sasacms < backup-before-fix.sql

# Or restore individual records
mysql -uadmin -psasa0102 sasacms -e "RESTORE FROM avatar_urls_backup;"
```

---

## 📊 Summary

**Current Status:**
- ❌ Old data: `/storage/avatars/...`
- ✅ New uploads: `/api/v1/files/avatar/...`

**After Fix:**
- ✅ All data: `/api/v1/files/avatar/...`
- ✅ Consistent across entire system

---

**Next Steps:**
1. Run backup
2. Run SQL update query
3. Verify in database
4. Test API endpoint
5. Confirm frontend works
