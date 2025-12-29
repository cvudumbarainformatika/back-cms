# Fix Database Avatar URLs - Quick Guide

## 🎯 Current Issue

Database masih menyimpan avatar dengan URL lama:

```
❌ /storage/avatars/avatar_1_1767013716.png
```

Response endpoint masih return nilai dari database (URL lama).

---

## ✅ Solution: Update Database

### Run This SQL Command

```sql
UPDATE users 
SET avatar = CONCAT('/api/v1/files/avatar/', SUBSTRING_INDEX(avatar, '/', -1))
WHERE avatar LIKE '/storage/avatars/%';
```

**Penjelasan:**
- `SUBSTRING_INDEX(avatar, '/', -1)` → Extract filename (bagian terakhir setelah `/`)
  - `/storage/avatars/avatar_1_1767013716.png` → `avatar_1_1767013716.png`
- `CONCAT('/api/v1/files/avatar/', ...)` → Build new URL
  - `avatar_1_1767011716.png` → `/api/v1/files/avatar/avatar_1_1767011716.png`

---

## 🚀 Run Command

### Option 1: Direct Command

```bash
docker exec starter-mysql mysql -uadmin -psasa0102 sasacms -e \
"UPDATE users SET avatar = CONCAT('/api/v1/files/avatar/', SUBSTRING_INDEX(avatar, '/', -1)) WHERE avatar LIKE '/storage/avatars/%';"
```

### Option 2: Via MySQL Client

```bash
# Start MySQL client
docker exec -it starter-mysql mysql -uadmin -psasa0102 sasacms

# Run query
UPDATE users 
SET avatar = CONCAT('/api/v1/files/avatar/', SUBSTRING_INDEX(avatar, '/', -1))
WHERE avatar LIKE '/storage/avatars/%';

# Exit
exit
```

### Option 3: Via SQL File

Create file `fix_avatars.sql`:
```sql
UPDATE users 
SET avatar = CONCAT('/api/v1/files/avatar/', SUBSTRING_INDEX(avatar, '/', -1))
WHERE avatar LIKE '/storage/avatars/%';
```

Then run:
```bash
docker exec starter-mysql mysql -uadmin -psasa0102 sasacms < fix_avatars.sql
```

---

## 🧪 Verify Before Update

```bash
# See what will be updated
docker exec starter-mysql mysql -uadmin -psasa0102 sasacms -e \
"SELECT id, email, avatar FROM users WHERE avatar LIKE '/storage/avatars/%';"

# Output example:
# id | email             | avatar
# 1  | admin@pdpi.co.id  | /storage/avatars/avatar_1_1767013716.png
```

---

## ✅ Verify After Update

```bash
# Check updated values
docker exec starter-mysql mysql -uadmin -psasa0102 sasacms -e \
"SELECT id, email, avatar FROM users WHERE id = 1;"

# Output should be:
# id | email             | avatar
# 1  | admin@pdpi.co.id  | /api/v1/files/avatar/avatar_1_1767013716.png
```

---

## 🧪 Test API After Fix

```bash
# Get profile
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN"

# Response should show:
{
  "avatar": "/api/v1/files/avatar/avatar_1_1767013716.png"  ✅
}

# Test file access
curl http://localhost:8080/api/v1/files/avatar/avatar_1_1767013716.png
# Should return image file ✅
```

---

## 🔄 Complete Fix Checklist

- [ ] Backup database (optional)
  ```bash
  docker exec starter-mysql mysqldump -uadmin -psasa0102 sasacms > backup.sql
  ```

- [ ] Verify what will be updated
  ```bash
  docker exec starter-mysql mysql -uadmin -psasa0102 sasacms -e \
  "SELECT COUNT(*) as count FROM users WHERE avatar LIKE '/storage/avatars/%';"
  ```

- [ ] Run update command (choose one option above)

- [ ] Verify update
  ```bash
  docker exec starter-mysql mysql -uadmin -psasa0102 sasacms -e \
  "SELECT id, avatar FROM users WHERE avatar NOT NULL;"
  ```

- [ ] Test API endpoint

- [ ] Test file serving

---

## 📊 Before & After

**Before:**
```
Database: /storage/avatars/avatar_1_1767013716.png
Response: /storage/avatars/avatar_1_1767013716.png
Access:   ❌ Direct file serving (unsafe)
```

**After:**
```
Database: /api/v1/files/avatar/avatar_1_1767013716.png
Response: /api/v1/files/avatar/avatar_1_1767013716.png
Access:   ✅ Via FileController (secure)
```

---

## 🚨 Rollback (If Needed)

```bash
# Restore from backup
docker exec starter-mysql mysql -uadmin -psasa0102 sasacms < backup.sql

# Or reverse the update
UPDATE users 
SET avatar = REPLACE(avatar, '/api/v1/files/avatar/', '/storage/avatars/')
WHERE avatar LIKE '/api/v1/files/avatar/%';
```

---

**After running the SQL command, all avatars will use the correct API URL format!** ✅
