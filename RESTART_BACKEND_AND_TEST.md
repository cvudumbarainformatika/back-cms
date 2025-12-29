# Restart Backend and Test Avatar Upload

## 🚀 Restart Backend

```bash
# Kill old process
pkill -9 -f "./backend" || true

# Wait a moment
sleep 1

# Start new backend
./backend &

# Verify running
sleep 2
curl http://localhost:8080/health
```

Or if using docker:

```bash
docker-compose restart app
sleep 3
curl http://localhost:8080/health
```

---

## 🧪 Test Avatar Upload Again

### Step 1: Login to get token

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@pdpi.co.id",
    "password": "AdminPDPI2024!"
  }' | jq '.data.access_token'

# Save the token
TOKEN="your-token-here"
```

### Step 2: Upload new avatar with multipart

```bash
curl -X PUT http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer $TOKEN" \
  -F "name=Admin Pusat" \
  -F "phone=081237660656" \
  -F "avatar=@/path/to/image.jpg"
```

Expected response:
```json
{
  "success": true,
  "data": {
    "avatar": "/api/v1/files/avatar/avatar_1_1767015087.png",
    ...
  }
}
```

### Step 3: Check if file saved

```bash
# Check storage directory
ls -la storage/avatars/

# Should show:
# avatar_1_1767015087.png  (NEW file with new timestamp)
```

### Step 4: Access file via API

```bash
curl http://localhost:8080/api/v1/files/avatar/avatar_1_1767015087.png \
  -o downloaded-image.jpg

# Should save the image file
ls -la downloaded-image.jpg
```

### Step 5: Check database

```bash
docker exec starter-mysql mysql -uadmin -psasa0102 sasacms -e \
"SELECT id, avatar FROM users WHERE id = 1;"

# Should show:
# id | avatar
# 1  | /api/v1/files/avatar/avatar_1_1767015087.png
```

---

## 🔍 Debug: Check Logs

If something goes wrong, check backend logs:

```bash
# If running in foreground
# Look for: [Avatar] Using storage path: ...
# This will show which path was used

# If running in background
tail -100 /tmp/backend.log | grep Avatar

# Or start backend with output
./backend 2>&1 | grep -i avatar
```

Look for log lines:
```
[Avatar] Found AVATAR_STORAGE_PATH env var: ...
[Avatar] Using STORAGE_BASE_PATH with /avatars: ...
[Avatar] Using default fallback path: ./storage/avatars
[Avatar] Using storage path: ...
[Avatar] Avatar uploaded successfully: ...
```

---

## ✅ Expected Behavior After Fix

1. **Upload new avatar**
   - File saved to: `storage/avatars/avatar_1_XXXXX.jpg`
   - Response includes: `/api/v1/files/avatar/avatar_1_XXXXX.jpg`

2. **Access file**
   - `http://localhost:8080/api/v1/files/avatar/avatar_1_XXXXX.jpg`
   - Should return image file (not 404)

3. **Database**
   - Stores: `/api/v1/files/avatar/avatar_1_XXXXX.jpg`
   - Next response returns same URL

---

## 🚨 Troubleshooting

### Error: File not found (404)
**Cause:** Storage path not set correctly

**Check:**
```bash
# Backend logs should show:
# [Avatar] Found AVATAR_STORAGE_PATH env var: ...
# or
# [Avatar] Using STORAGE_BASE_PATH with /avatars: ...

# If shows default path, env vars not set!
```

**Fix:**
```bash
# Check .env has AVATAR_STORAGE_PATH
grep AVATAR_STORAGE_PATH .env

# If not there, add it:
echo "AVATAR_STORAGE_PATH=./storage/avatars" >> .env

# Restart backend
```

### Error: Cannot create file
**Cause:** Directory doesn't exist or permissions issue

**Check:**
```bash
ls -la storage/
mkdir -p storage/avatars
chmod 755 storage/avatars
```

### Error: JSON unmarshal error
**Cause:** Sending JSON with avatar field

**Fix:**
```bash
# Use multipart for avatar upload:
curl -X PUT ... -F "avatar=@image.jpg"

# Or remove avatar field from JSON:
curl -X PUT ... -d '{"name":"...", "phone":"..."}'
```

---

## ✅ Checklist After Restart

- [ ] Backend running: `curl http://localhost:8080/health` → 200 OK
- [ ] Upload new avatar works
- [ ] File saved to storage/avatars/
- [ ] API endpoint returns file
- [ ] Database shows new URL format
- [ ] All tests pass

---

**After restart and successful upload, old database data will still show old URLs. To fix that, run SQL migration command (separate task).** ✅
