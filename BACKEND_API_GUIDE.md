# Backend API Guide - Go Gin Starter Kit

## Struktur Aplikasi

```
app/
├── Exceptions/          # Error handling & middleware
├── Http/
│   ├── Controllers/     # Business logic (request -> response)
│   ├── Middleware/      # JWT, CORS, Rate Limiter, Logger
│   └── Requests/        # Request validation & binding
├── Models/              # Database models & queries
bootstrap/               # App initialization
config/                  # Configuration management
database/
├── migrations/          # Database schema
└── seeders/             # Database seeders
routes/                  # Route definitions
utils/                   # Helpers (JWT, pagination, response)
```

## 1. AUTHENTICATION SYSTEM

### 1.1 Login Flow

**Endpoint:** `POST /api/v1/auth/login`

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (Success - 200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "role": "admin",
      "id_hotel": 1,
      "id_wisata": 0
    }
  }
}
```

**Response (Error - 401):**
```json
{
  "success": false,
  "error": "invalid_credentials",
  "message": "Invalid username or password"
}
```

### 1.2 Registration Flow

**Endpoint:** `POST /api/v1/auth/register`

**Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (Success - 201):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

### 1.3 JWT Configuration

**File:** `.env`
```env
JWT_SECRET=your-secret-key-change-this-in-production
JWT_ACCESS_TOKEN_EXPIRATION=15          # Minutes
JWT_REFRESH_TOKEN_EXPIRATION=10080      # 7 days in minutes
```

**JWT Claims Structure:**
```go
type JWTClaims struct {
    UserID int64  `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}
```

---

## 2. USER CRUD API (WITH MIDDLEWARE & PAGINATION)

### 2.1 Get All Users (List)

**Endpoint:** `GET /api/v1/users/get-list?page=1&per_page=10&q=search&order_by=created_at&sort=desc`

**Headers Required:**
```
Authorization: Bearer {access_token}
```

**Query Parameters:**
- `page`: Page number (default: 1)
- `per_page`: Items per page (default: 10, max: 100)
- `q`: Search query (searches in name, hotel, wisata)
- `order_by`: Sort field (created_at, nama, nameUser, statusUser)
- `sort`: Sort direction (asc, desc)

**Response (200):**
```json
{
  "current_page": 1,
  "data": [
    {
      "id": 1,
      "nameUser": "john@example.com",
      "passUser": "$2a$10$...",
      "statusUser": "admin",
      "idHotel": 1,
      "idWisata": 0,
      "nameHotel": "Grand Hotel",
      "nameWisata": null
    }
  ],
  "first_page_url": "http://localhost:8080/api/v1/users/get-list?page=1&per_page=10",
  "from": 1,
  "last_page": 5,
  "next_page_url": "http://localhost:8080/api/v1/users/get-list?page=2&per_page=10",
  "path": "http://localhost:8080/api/v1/users/get-list",
  "per_page": 10,
  "prev_page_url": null,
  "to": 10,
  "total": 50
}
```

### 2.2 Get User By ID

**Endpoint:** `GET /api/v1/users/:id`

**Response (200):**
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "user": {
      "id": 1,
      "nameUser": "john@example.com",
      "statusUser": "admin",
      "idHotel": 1,
      "nameHotel": "Grand Hotel"
    }
  }
}
```

### 2.3 Create User

**Endpoint:** `POST /api/v1/users/create`

**Request:**
```json
{
  "email": "newuser@example.com",
  "password": "securepass123",
  "role": "admin",
  "idHotel": 1,
  "idWisata": 0
}
```

**Response (201):**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 2,
    "nameUser": "newuser@example.com",
    "email": "newuser@example.com",
    "statusUser": "admin",
    "idHotel": 1,
    "idWisata": 0,
    "nameHotel": "Grand Hotel",
    "nameWisata": null
  }
}
```

### 2.4 Update User

**Endpoint:** `PUT /api/v1/users/update/:id`

**Request:**
```json
{
  "email": "updated@example.com",
  "password": "newpassword123",
  "role": "user",
  "idHotel": 2,
  "idWisata": 1
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "nameUser": "updated@example.com",
    "statusUser": "user",
    "idHotel": 2,
    "idWisata": 1
  }
}
```

### 2.5 Delete User

**Endpoint:** `DELETE /api/v1/users/delete/:id`

**Response (200):**
```json
{
  "success": true,
  "message": "User deleted successfully",
  "data": null
}
```

---

## 3. MIDDLEWARE STACK

### 3.1 JWT Authentication Middleware

**File:** `app/Http/Middleware/jwt_auth.go`

**How it works:**
1. Extracts `Authorization: Bearer <token>` header
2. Validates token signature using JWT secret
3. Extracts user_id and email from claims
4. Sets values in context: `c.Set("user_id", claims.UserID)` dan `c.Set("user_email", claims.Email)`
5. If invalid, returns 401 Unauthorized

**Usage in routes:**
```go
protected := v1.Group("/")
protected.Use(middleware.JWTAuthMiddleware(cfg.JWT.Secret))
{
    // All routes here require valid JWT
    users := protected.Group("/users")
    users.GET("/get-list", userController.GetAllUsers)
}
```

### 3.2 CORS Middleware

**File:** `app/Http/Middleware/cors.go`

**Configuration (`.env`):**
```env
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
```

### 3.3 Rate Limiter Middleware

**File:** `app/Http/Middleware/rate_limiter.go`

**Configuration (`.env`):**
```env
RATE_LIMIT_ENABLED=true
RATE_LIMIT_MAX_REQUESTS=100
RATE_LIMIT_WINDOW_SECONDS=60
```

Limits 100 requests per 60 seconds per IP address.

### 3.4 Logger Middleware

**File:** `app/Http/Middleware/logger.go`

Logs setiap request dengan format:
```
[GIN] GET /api/v1/users 200 150ms
```

### 3.5 Error Handler Middleware

**File:** `app/Exceptions/handler.go`

Menangkap panic dan error, mengembalikan response JSON:
```json
{
  "success": false,
  "error": "internal_error",
  "message": "Internal server error"
}
```

---

## 4. REQUEST VALIDATION

### 4.1 Validation Pattern

**File:** `app/Http/Requests/login_request.go`

```go
type LoginRequest struct {
    Email    string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func (r *LoginRequest) Validate(c *gin.Context) error {
    if err := c.ShouldBindJSON(r); err != nil {
        utils.ValidationError(c, errStr)
        return err
    }
    return nil
}
```

**In Controller:**
```go
var req requests.LoginRequest
if err := req.Validate(c); err != nil {
    return  // Error response already sent by Validate()
}
// Continue with business logic
```

### 4.2 Binding Tags

- `required`: Field harus ada
- `email`: Harus format email
- `min=N`: Minimum length N
- `max=N`: Maximum length N
- `omitempty`: Optional field

---

## 5. PAGINATION HELPER

### 5.1 Get Pagination Parameters

**File:** `utils/pagination.go`

```go
func GetPaginationParams(c *gin.Context) (page int, limit int) {
    page = 1
    limit = 10
    
    if pageStr := c.Query("page"); pageStr != "" {
        if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
            page = p
        }
    }
    
    if limitStr := c.Query("limit"); limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
            limit = l
        }
    }
    
    return page, limit
}
```

### 5.2 Laravel-Style Pagination Response

```go
pagination := utils.CreateLaravelPagination(c, users, page, limit, total)
c.JSON(http.StatusOK, pagination)
```

Returns:
```json
{
  "current_page": 1,
  "data": [...],
  "first_page_url": "...",
  "from": 1,
  "last_page": 5,
  "next_page_url": "...",
  "path": "...",
  "per_page": 10,
  "prev_page_url": null,
  "to": 10,
  "total": 50
}
```

---

## 6. RESPONSE HELPER

### 6.1 Success Response

**File:** `utils/response.go`

```go
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
    c.JSON(statusCode, SuccessResponse{
        Success: true,
        Message: message,
        Data:    data,
    })
}
```

**Usage:**
```go
utils.Success(c, http.StatusOK, "User retrieved", gin.H{
    "user": user,
})
```

**Response:**
```json
{
  "success": true,
  "message": "User retrieved",
  "data": {
    "user": {...}
  }
}
```

### 6.2 Error Response

```go
func Error(c *gin.Context, statusCode int, err string, message string, details interface{}) {
    c.JSON(statusCode, ErrorResponse{
        Success: false,
        Error:   err,
        Message: message,
        Details: details,
    })
}
```

**Usage:**
```go
utils.Error(c, http.StatusNotFound, "user_not_found", "User not found", nil)
```

### 6.3 Validation Error Response

```go
func ValidationError(c *gin.Context, errors interface{}) {
    c.JSON(422, ErrorResponse{
        Success: false,
        Error:   "validation_failed",
        Message: "Validation failed",
        Details: errors,
    })
}
```

---

## 7. DATABASE MODELS

### 7.1 User Model

**File:** `app/Models/user.go`

```go
type User struct {
    ID         int     `db:"idUser" json:"id"`
    Name       string  `db:"nameUser" json:"nameUser"`
    Password   string  `db:"passUser" json:"passUser"`
    Status     string  `db:"statusUser" json:"statusUser"`
    IDHotel    int     `db:"idHotel" json:"idHotel"`
    IDWisata   int     `db:"idWisata" json:"idWisata"`
    HotelName  *string `db:"nameHotel" json:"nameHotel"`
    WisataName *string `db:"nameWisata" json:"nameWisata"`
}
```

**Methods:**
- `Create(db)`: Insert user
- `FindUserByID(db, id)`: Get user by ID with joins
- `FindUserByName(db, name)`: Get user by nameUser
- `GetUsers(db, params)`: List with filter & pagination
- `Update(db)`: Update user
- `HashPassword()`: Hash password dengan bcrypt
- `ComparePassword(plain)`: Compare bcrypt atau MD5 password
- `DeleteUserByID(db, id)`: Delete user

**Database Schema:**
```sql
CREATE TABLE users (
    idUser INT AUTO_INCREMENT PRIMARY KEY,
    nameUser VARCHAR(255) NOT NULL,
    passUser VARCHAR(255) NOT NULL,
    statusUser VARCHAR(50),
    idHotel INT DEFAULT 0,
    idWisata INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

---

## 8. APPLICATION BOOTSTRAP

### 8.1 Application Initialization

**File:** `bootstrap/app.go`

**Flow:**
1. Load configuration dari `.env`
2. Set Gin mode (release/debug)
3. Initialize Gin engine
4. Register global middleware
5. Connect ke database
6. Connect ke Redis
7. Setup routes

**Shutdown:**
- Close database connections
- Close Redis connection
- Graceful shutdown

### 8.2 Main Entry Point

**File:** `main.go`

```go
func main() {
    app, err := bootstrap.NewApplication()
    if err != nil {
        log.Fatalf("Failed to initialize application: %v", err)
    }
    
    // Run server
    go func() {
        if err := app.Run(); err != nil {
            log.Fatalf("Failed to run server: %v", err)
        }
    }()
    
    // Wait for interrupt signal (SIGINT, SIGTERM)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    // Graceful shutdown
    if err := app.Shutdown(); err != nil {
        log.Fatalf("Failed to shutdown application: %v", err)
    }
}
```

---

## 9. FLOW DIAGRAM - USER CRUD

### Create User Flow:

```
POST /api/v1/users/create
    ↓
[JWT Auth Middleware] - Validate token
    ↓
[Rate Limiter] - Check rate limit
    ↓
UserController.CreateUser()
    ↓
[Validate Request] - Email, Password, Role required
    ↓
[Check Existence] - User dengan email sudah ada?
    ↓
[Hash Password] - bcrypt hash
    ↓
[Database Insert] - Save ke database
    ↓
[Success Response] - Return created user (201)
```

### Get Users with Pagination Flow:

```
GET /api/v1/users/get-list?page=1&per_page=10
    ↓
[JWT Auth Middleware] - Validate token
    ↓
[Rate Limiter] - Check rate limit
    ↓
UserController.GetAllUsers()
    ↓
[Extract Params] - page, limit, q, order_by, sort
    ↓
[Build Query] - WHERE with search filter
    ↓
[Get Total Count] - SELECT COUNT(*)
    ↓
[Apply Pagination] - LIMIT ? OFFSET ?
    ↓
[Apply Sorting] - ORDER BY field ASC/DESC
    ↓
[Left Joins] - hotels dan wisata tables
    ↓
[Create Response] - Laravel-style pagination format
    ↓
[Success Response] - Return paginated data (200)
```

---

## 10. CONFIGURATION REFERENCE

### 10.1 .env Configuration

```env
# Application
APP_NAME=Go Gin Starter Kit
APP_ENV=local                           # local, production
APP_PORT=8080

# Database
DB_CONNECTION=mysql                     # mysql, postgres
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=starter_kit
DB_USERNAME=root
DB_PASSWORD=secret

# Connection Pool
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=300                # seconds

# JWT
JWT_SECRET=your-secret-key-change-this
JWT_ACCESS_TOKEN_EXPIRATION=15          # minutes
JWT_REFRESH_TOKEN_EXPIRATION=10080      # minutes (7 days)

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_MAX_REQUESTS=100
RATE_LIMIT_WINDOW_SECONDS=60

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

---

## 11. HTTP STATUS CODES REFERENCE

| Code | Usage |
|------|-------|
| 200 | Success (GET, PUT, DELETE) |
| 201 | Created (POST - create resource) |
| 204 | No Content |
| 400 | Bad Request (validation error) |
| 401 | Unauthorized (invalid/missing token) |
| 403 | Forbidden (no permission) |
| 404 | Not Found (resource not found) |
| 422 | Unprocessable Entity (validation error) |
| 429 | Too Many Requests (rate limit) |
| 500 | Internal Server Error |

---

## 12. TESTING THE API

### 12.1 Login & Get Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### 12.2 List Users (with token)

```bash
TOKEN="your-access-token-here"

curl -X GET "http://localhost:8080/api/v1/users/get-list?page=1&per_page=10" \
  -H "Authorization: Bearer $TOKEN"
```

### 12.3 Create User

```bash
curl -X POST http://localhost:8080/api/v1/users/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "email":"newuser@example.com",
    "password":"securepass123",
    "role":"admin",
    "idHotel":1,
    "idWisata":0
  }'
```

---

## Summary

Starter kit ini menyediakan:
- ✅ Complete authentication (Login, Register)
- ✅ JWT token generation & validation
- ✅ User CRUD API dengan paginasi
- ✅ Middleware stack (JWT, CORS, Rate Limiter, Logger, Error Handler)
- ✅ Request validation dengan binding tags
- ✅ Response helper (success, error, validation)
- ✅ Database abstraction dengan sqlx
- ✅ Configuration management dari .env
- ✅ Graceful shutdown
- ✅ Redis integration ready

Ready untuk dikembangkan lebih lanjut!
