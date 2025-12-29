# Test Avatar Upload - Check If Code Fix Works

## 🎯 Goal

Test upload avatar baru dari endpoint `/auth/profile` untuk verify code fix sudah bekerja.

---

## 📝 Test Cases

### Case 1: Update Profile WITHOUT Avatar (JSON)

```bash
curl -X PUT http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin Pusat",
    "phone": "081237660656",
    "address": "Jakarta",
    "bio": "Administrator"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "name": "Admin Pusat",
    "phone": "081237660656",
    "address": "Jakarta",
    "bio": "Administrator",
    "avatar": "/api/v1/files/avatar/avatar_1_1767013716.png"  ✅ (old data from DB)
  }
}
```

---

### Case 2: Upload New Avatar (Multipart with file)

```bash
curl -X PUT http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "name=Admin Pusat" \
  -F "phone=081237660656" \
  -F "address=Jakarta" \
  -F "bio=Administrator" \
  -F "avatar=@/path/to/new-image.jpg"
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "name": "Admin Pusat",
    "phone": "081237660656",
    "address": "Jakarta",
    "bio": "Administrator",
    "avatar": "/api/v1/files/avatar/avatar_1_1767013800.jpg"  ✅ (NEW URL format!)
  }
}
```

---

## 🔍 What to Check

After upload:

1. **Response avatar URL**
   - Check format: `/api/v1/files/avatar/{filename}`
   - NOT `/storage/avatars/{filename}`

2. **File exists**
   ```bash
   ls -la storage/avatars/ | grep avatar_1
   ```

3. **Can access via API**
   ```bash
   curl http://localhost:8080/api/v1/files/avatar/avatar_1_XXXXX.jpg
   # Should return image file
   ```

4. **Database updated**
   ```bash
   docker exec starter-mysql mysql -uadmin -psasa0102 sasacms -e \
   "SELECT id, avatar FROM users WHERE id = 1;"
   # Should show: /api/v1/files/avatar/avatar_1_XXXXX.jpg
   ```

---

## 📊 What Changed in Code

1. **UploadAvatar()** now returns:
   ```go
   ✅ /api/v1/files/avatar/{filename}
   ❌ /storage/avatars/{filename}
   ```

2. **UpdateProfile()** saves what UploadAvatar() returns
   - No URL transformation
   - Direct save to database

3. **JSON parsing** fixed
   - No unmarshal error when avatar is empty

---

## ✅ Expected Behavior

### First Time (Old Data in DB)
- Response shows old URL from database: `/storage/avatars/...`
- This is OK - it's data from previous upload

### After New Upload
- Response shows NEW URL format: `/api/v1/files/avatar/...`
- Database updated with new URL
- All future requests will show new format

---

## 🧪 Complete Test Workflow

```bash
# 1. Get current user (should show old avatar if exists)
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer TOKEN"
# Output: "avatar": "/storage/avatars/..."

# 2. Upload new avatar
curl -X PUT http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer TOKEN" \
  -F "name=Admin Pusat" \
  -F "avatar=@new-image.jpg"
# Output: "avatar": "/api/v1/files/avatar/..."  ✅

# 3. Get user again (should show new avatar URL)
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer TOKEN"
# Output: "avatar": "/api/v1/files/avatar/..."  ✅

# 4. Access avatar file
curl http://localhost:8080/api/v1/files/avatar/avatar_1_XXXXX.jpg
# Should return image file  ✅
```

---

## 🎯 How to Get TOKEN

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@pdpi.co.id",
    "password": "AdminPDPI2024!"
  }'

# Copy access_token dari response
# Use in Authorization: Bearer <token>
```

---

## ✅ Verification Checklist

After uploading new avatar:

- [ ] Response HTTP 200 OK
- [ ] Response avatar URL format: `/api/v1/files/avatar/...`
- [ ] File exists in: `storage/avatars/avatar_1_XXXXX.jpg`
- [ ] File accessible via: `http://localhost:8080/api/v1/files/avatar/avatar_1_XXXXX.jpg`
- [ ] Database shows new URL format

---

**Next Step:** Try upload new avatar dari frontend/curl dan share response! 🚀
