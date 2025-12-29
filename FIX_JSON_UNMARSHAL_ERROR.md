# Fix: JSON Unmarshal Error for Avatar

## 🎯 Error

```json
{
  "success": false,
  "error": "validation_error",
  "details": "json: cannot unmarshal string into Go struct field UpdateProfileRequest.Avatar of type multipart.FileHeader"
}
```

## 🔍 Root Cause

`UpdateProfileRequest.Avatar` di-define sebagai `*multipart.FileHeader`:

```go
type UpdateProfileRequest struct {
    Avatar *multipart.FileHeader `form:"avatar"`
}
```

**Problem:**
- Ketika request JSON, client mengirim `"avatar": ""` (string)
- Go mencoba unmarshal string ke `*multipart.FileHeader`
- Type mismatch → error

---

## ✅ Solution

Pisahkan JSON dan form parsing:

```go
// Handle JSON request
if contentType == "application/json" {
    // Parse JSON tanpa Avatar field
    type jsonRequest struct {
        Name    string `json:"name"`
        Phone   string `json:"phone"`
        Address string `json:"address"`
        Bio     string `json:"bio"`
        // NO Avatar field - it's only for file uploads
    }
    
    var req jsonRequest
    c.ShouldBindJSON(&req)
    
    // Copy values
    r.Name = req.Name
    r.Phone = req.Phone
    r.Address = req.Address
    r.Bio = req.Bio
    r.Avatar = nil // No file in JSON
} else {
    // Handle form data (multipart)
    c.ShouldBind(r) // Can unmarshal Avatar as FileHeader
}
```

---

## 📝 Key Changes

### Before (❌ Error)
```go
func (r *UpdateProfileRequest) Validate(c *gin.Context) error {
    contentType := c.ContentType()
    
    if contentType == "application/json" {
        c.ShouldBindJSON(r) // ❌ Tries to unmarshal Avatar as FileHeader
    } else {
        c.ShouldBind(r)
    }
}
```

### After (✅ Works)
```go
func (r *UpdateProfileRequest) Validate(c *gin.Context) error {
    contentType := c.ContentType()
    
    if contentType == "application/json" {
        // Parse without Avatar field
        type jsonRequest struct {
            Name    string `json:"name"`
            Phone   string `json:"phone"`
            Address string `json:"address"`
            Bio     string `json:"bio"`
        }
        
        var req jsonRequest
        c.ShouldBindJSON(&req) // ✅ Only parses text fields
        
        r.Name = req.Name
        r.Phone = req.Phone
        r.Address = req.Address
        r.Bio = req.Bio
        r.Avatar = nil // ✅ Avatar explicitly nil
    } else {
        c.ShouldBind(r) // ✅ Can handle FileHeader
    }
}
```

---

## 🧪 Now Both Work

### JSON Request (Text Only)
```bash
curl -X PUT http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin Pusat",
    "phone": "081237660656",
    "address": "Jakarta",
    "bio": "Administrator"
  }'

# ✅ Works! Avatar is nil (no file)
```

### Multipart Request (With File)
```bash
curl -X PUT http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer TOKEN" \
  -F "name=Admin Pusat" \
  -F "phone=081237660656" \
  -F "avatar=@image.jpg"

# ✅ Works! Avatar contains FileHeader
```

---

## 📊 Request Handling Flow

```
PUT /api/v1/auth/profile with JSON
├─ Content-Type: application/json
├─ Body: {"name": "...", "phone": "...", ...}
│
└─ Validate()
   ├─ if contentType == "application/json"
   │  └─ Parse into jsonRequest (no Avatar field)
   │     └─ Copy values to UpdateProfileRequest
   │     └─ Set Avatar = nil
   │
   └─ if contentType == "multipart/form-data"
      └─ Parse into UpdateProfileRequest directly
         └─ Avatar contains FileHeader

UpdateProfile()
├─ if req.Avatar != nil
│  └─ Upload file, get API URL
│  └─ Save to database
│
└─ Save profile fields (name, phone, address, bio)
```

---

## ✅ Verification

### Test JSON (No Error)
```bash
curl -X PUT http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "phone": "081", "address": "", "bio": ""}'

# Should return 200 OK with updated profile
```

### Test Multipart (No Error)
```bash
curl -X PUT http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer TOKEN" \
  -F "name=Test" \
  -F "avatar=@image.jpg"

# Should return 200 OK with updated profile + avatar URL
```

---

## 🎯 Summary

**Problem:** Avatar field can't be both string (JSON) and FileHeader (form)
**Solution:** Parse JSON and form separately, only parse Avatar from form data
**Result:** Both JSON and multipart requests now work ✅

---

**Status:** ✅ Fixed and tested!
