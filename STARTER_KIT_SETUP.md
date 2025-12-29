# Go Gin Starter Kit - Setup Guide

Starter kit yang bersih dan siap untuk development dengan struktur yang lengkap untuk authentication, CRUD API, middleware, pagination, dan response helpers.

## ✅ What's Included

### Infrastructure & Setup
- ✅ Gin web framework dengan middleware stack
- ✅ JWT authentication (access & refresh tokens)
- ✅ CORS middleware
- ✅ Rate limiter middleware
- ✅ Logger middleware
- ✅ Error handler middleware
- ✅ Database abstraction dengan sqlx (MySQL & PostgreSQL)
- ✅ Redis integration ready
- ✅ Graceful shutdown
- ✅ Configuration management (.env)

### Utilities
- ✅ JWT token generation & validation
- ✅ Pagination helpers (offset-based & cursor-based)
- ✅ Response helpers (success, error, validation)
- ✅ Request validation dengan binding tags

### Documentation
- ✅ `BACKEND_API_GUIDE.md` - Complete API reference dengan examples
- ✅ `STARTER_KIT_SETUP.md` - Setup & development guide
- ✅ Template files untuk controller, model, dan request

---

## 🚀 Quick Start

### 1. Setup Environment

Copy `.env.example` ke `.env` dan sesuaikan:

```bash
cp .env.example .env
```

Edit `.env`:
```env
APP_NAME=My API
APP_ENV=local
APP_PORT=8080

DB_CONNECTION=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=my_database
DB_USERNAME=root
DB_PASSWORD=secret

JWT_SECRET=your-super-secret-key-change-this
JWT_ACCESS_TOKEN_EXPIRATION=15
JWT_REFRESH_TOKEN_EXPIRATION=10080

# Optional
REDIS_HOST=localhost
REDIS_PORT=6379
```

### 2. Setup Database

Buat database baru:
```sql
CREATE DATABASE my_database;
```

### 3. Run Application

```bash
# Install dependencies
go mod download

# Run application
go run main.go

# Output: Starting Go Gin Starter Kit on :8080 (env: local)
```

### 4. Test Health Endpoint

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "ok",
  "database": "connected"
}
```

---

## 📁 Project Structure

```
app/
├── Exceptions/              # Error handling
│   ├── errors.go           # Error definitions
│   └── handler.go          # Global error handler
│
├── Http/
│   ├── Controllers/        # API logic (ADD YOUR CONTROLLERS HERE)
│   │   └── CONTROLLER_TEMPLATE.go
│   │
│   ├── Middleware/         # Request middleware
│   │   ├── jwt_auth.go     # JWT validation
│   │   ├── cors.go         # CORS handling
│   │   ├── rate_limiter.go # Rate limiting
│   │   ├── logger.go       # Request logging
│   │   └── authorization.go
│   │
│   └── Requests/           # Request validation (ADD YOUR REQUESTS HERE)
│       ├── login_request.go
│       ├── register_request.go
│       └── EXAMPLE_REQUEST_TEMPLATE.go
│
├── Models/                 # Database models (ADD YOUR MODELS HERE)
│   └── MODEL_TEMPLATE.go
│
bootstrap/                  # Application initialization
│   └── app.go
│
config/                     # Configuration
│   └── config.go
│
database/
├── database.go            # Database connection
├── redis.go               # Redis connection
├── migrations/            # SQL migrations (ADD YOUR MIGRATIONS HERE)
└── seeders/              # Database seeders (ADD YOUR SEEDERS HERE)

routes/
└── api.go                # Route definitions

utils/
├── jwt.go                # JWT utilities
├── pagination.go         # Pagination helpers
└── response.go          # Response formatting

main.go                   # Entry point
```

---

## 🔧 Development Workflow

### 1. Buat Migration SQL

**File:** `database/migrations/001_create_examples_table.sql`

```sql
CREATE TABLE IF NOT EXISTS examples (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

Run migration manually di database atau via SQL client.

### 2. Buat Model

**File:** `app/Models/example.go`

Copy template dari `MODEL_TEMPLATE.go` dan sesuaikan:

```go
package models

import (
    "database/sql"
    "github.com/cvudumbarainformatika/backend/utils"
    "github.com/jmoiron/sqlx"
)

type Example struct {
    ID    int    `db:"id" json:"id"`
    Name  string `db:"name" json:"name"`
    Email string `db:"email" json:"email"`
}

// Implement Create, FindByID, GetAll, Update, Delete methods
```

### 3. Buat Request Validation

**File:** `app/Http/Requests/create_example_request.go`

Copy template dari `EXAMPLE_REQUEST_TEMPLATE.go`:

```go
type CreateExampleRequest struct {
    Name  string `json:"name" binding:"required,min=3"`
    Email string `json:"email" binding:"required,email"`
}

func (r *CreateExampleRequest) Validate(c *gin.Context) error {
    if err := c.ShouldBindJSON(r); err != nil {
        utils.ValidationError(c, err.Error())
        return err
    }
    return nil
}
```

### 4. Buat Controller

**File:** `app/Http/Controllers/example_controller.go`

Copy template dari `CONTROLLER_TEMPLATE.go`:

```go
type ExampleController struct {
    db *sqlx.DB
}

func (ec *ExampleController) GetAll(c *gin.Context) {
    // Implementation
}
```

### 5. Daftarkan Routes

**File:** `routes/api.go`

```go
protected := v1.Group("/")
protected.Use(middleware.JWTAuthMiddleware(cfg.JWT.Secret))
{
    exampleController := controllers.NewExampleController(db)
    examples := protected.Group("/examples")
    {
        examples.GET("/get-list", exampleController.GetAll)
        examples.GET("/:id", exampleController.GetByID)
        examples.POST("/create", exampleController.Create)
        examples.PUT("/update/:id", exampleController.Update)
        examples.DELETE("/delete/:id", exampleController.Delete)
    }
}
```

---

## 🔐 Authentication Flow

### 1. Register (Public)

```bash
POST /api/v1/auth/register

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

Implement di controller:
```go
// authController.Register - implement this
```

### 2. Login (Public)

```bash
POST /api/v1/auth/login

{
  "email": "john@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

Implement di controller:
```go
// authController.Login - implement this
```

### 3. Gunakan Token di Protected Routes

```bash
GET /api/v1/examples/get-list

Headers:
Authorization: Bearer {access_token}
```

JWT middleware akan:
1. Extract token dari header
2. Validate signature
3. Check expiration
4. Set user data di context
5. Allow request atau return 401

---

## 📝 Binding Tags Reference

### Common Tags

```go
type User struct {
    Name     string `json:"name" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Age      int    `json:"age" binding:"omitempty,min=0,max=150"`
    Phone    string `json:"phone" binding:"omitempty"`
}
```

| Tag | Usage |
|-----|-------|
| `required` | Field harus ada |
| `omitempty` | Field optional |
| `email` | Harus format email |
| `min=N` | Minimum length/value N |
| `max=N` | Maximum length/value N |
| `numeric` | Harus angka |
| `alpha` | Hanya huruf |
| `alphanum` | Huruf & angka |
| `url` | Harus URL valid |
| `uuid` | Harus UUID valid |

---

## 🛠️ Utils Usage

### Response Helpers

```go
// Success response
utils.Success(c, http.StatusOK, "User retrieved", gin.H{
    "user": user,
})

// Error response
utils.Error(c, http.StatusNotFound, "user_not_found", "User not found", nil)

// Validation error
utils.ValidationError(c, gin.H{
    "email": "email is required",
})
```

### Pagination Helpers

```go
// Get pagination params
page, limit := utils.GetPaginationParams(c)  // Returns page, limit with defaults

// Get filter params (includes pagination + search + sort)
params := utils.GetFilterParams(c)  // FilterParams{Page, PerPage, Q, OrderBy, Sort}

// Create pagination response
pagination := utils.CreateLaravelPagination(c, data, page, limit, total)
c.JSON(http.StatusOK, pagination)
```

### JWT Helpers

```go
// Generate token
token, err := utils.GenerateAccessToken(userID, email, secret, 15)

// Validate token
claims, err := utils.ValidateToken(tokenString, secret)
if err != nil {
    // Token invalid or expired
}
```

---

## 🧪 Testing

### Test Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Test Protected Endpoint

```bash
TOKEN="your-access-token-here"

curl -X GET http://localhost:8080/api/v1/examples/get-list \
  -H "Authorization: Bearer $TOKEN"
```

### Test Pagination

```bash
curl -X GET "http://localhost:8080/api/v1/examples/get-list?page=1&per_page=10&q=search&sort=desc" \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📚 Next Steps

1. **Implement Auth**: Buat login/register controller
2. **Create First API**: Buat CRUD endpoint untuk resource pertama
3. **Add Database Migrations**: Create tables yang diperlukan
4. **Test Endpoints**: Pastikan semua working dengan Postman/curl
5. **Setup CI/CD**: Deploy ke production

---

## 🐛 Troubleshooting

### Database Connection Error
- Pastikan database sudah berjalan
- Cek `.env` database config
- Verify database user & password

### JWT Token Invalid
- Verify `JWT_SECRET` di `.env`
- Check token format: `Authorization: Bearer {token}`
- Verify token belum expired

### CORS Error
- Check `CORS_ALLOWED_ORIGINS` di `.env`
- Ensure frontend domain di allowlist

### Rate Limit Error (429)
- Wait atau adjust `RATE_LIMIT_MAX_REQUESTS` di `.env`

---

## 📖 References

- See `BACKEND_API_GUIDE.md` for complete API documentation
- See `CONTROLLER_TEMPLATE.go` for controller example
- See `MODEL_TEMPLATE.go` for model example
- See `EXAMPLE_REQUEST_TEMPLATE.go` for validation example

---

## 📄 License

This starter kit is ready for development. Feel free to customize as needed!
