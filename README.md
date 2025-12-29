# Go Gin Backend Starter Kit

A production-ready, clean Go Gin REST API starter kit with complete authentication, middleware stack, pagination, and database abstraction.

## ✨ Features

### Infrastructure
- ✅ **JWT Authentication** - Access & refresh tokens with bcrypt password hashing
- ✅ **Middleware Stack** - JWT auth, CORS, rate limiting, logging, error handling
- ✅ **Database Abstraction** - MySQL & PostgreSQL support with sqlx
- ✅ **Redis Integration** - Ready for caching & sessions
- ✅ **Configuration Management** - Environment-based config from .env
- ✅ **Graceful Shutdown** - Proper cleanup on application termination

### Request/Response System
- ✅ **Request Validation** - Binding tags & structured validation
- ✅ **Response Formatting** - Consistent success/error responses
- ✅ **Pagination Helpers** - Offset, cursor, and Laravel-style pagination
- ✅ **Error Handling** - Global error handler with proper HTTP status codes

### Development Ready
- ✅ **Template Files** - Controller, model, request, and migration templates
- ✅ **Comprehensive Documentation** - Complete API guides & patterns
- ✅ **Clean Scripts** - Utilities for cleanup and git setup
- ✅ **Best Practices** - Implemented throughout the codebase

## 🚀 Quick Start

### Prerequisites
- Go 1.23+ 
- MySQL 5.7+ or PostgreSQL 10+
- Redis (optional)

### Setup

1. **Clone or initialize the starter kit**
```bash
git clone <your-starter-kit-url>
cd your-project
```

2. **Setup environment**
```bash
cp .env.example .env
# Edit .env with your database credentials & JWT secret
```

3. **Install dependencies**
```bash
go mod download
```

4. **Run application**
```bash
go run main.go
```

5. **Test health endpoint**
```bash
curl http://localhost:8080/health
```

## 📁 Project Structure

```
app/
├── Exceptions/              # Error handling
├── Http/
│   ├── Controllers/        # API logic (ADD YOUR CONTROLLERS HERE)
│   │   └── CONTROLLER_TEMPLATE.go
│   ├── Middleware/         # Request middleware (JWT, CORS, etc.)
│   └── Requests/           # Request validation (ADD YOUR REQUESTS HERE)
│       └── EXAMPLE_REQUEST_TEMPLATE.go
├── Models/                 # Database models (ADD YOUR MODELS HERE)
│   └── MODEL_TEMPLATE.go

bootstrap/                  # Application initialization
config/                     # Configuration management
database/
├── database.go            # Database connection
├── redis.go               # Redis connection
├── migrations/            # SQL migrations (ADD YOUR MIGRATIONS HERE)
│   └── MIGRATION_TEMPLATE.sql
└── seeders/              # Database seeders (ADD YOUR SEEDERS HERE)

routes/                     # Route definitions
utils/                      # Helpers (JWT, pagination, response)
main.go                     # Entry point
```

## 📚 Documentation

- **[BACKEND_API_GUIDE.md](./BACKEND_API_GUIDE.md)** - Complete API reference with authentication flow, CRUD patterns, middleware, pagination, and response helpers
- **[STARTER_KIT_SETUP.md](./STARTER_KIT_SETUP.md)** - Step-by-step setup guide and development workflow
- **[STARTER_KIT_SUMMARY.md](./STARTER_KIT_SUMMARY.md)** - Quick reference for key patterns and checklist

## 🔧 Development Workflow

### Creating a New API

1. **Create migration** (database schema)
```bash
# Edit database/migrations/001_create_your_table.sql
```

2. **Create model**
```bash
# Copy app/Models/MODEL_TEMPLATE.go to app/Models/your_model.go
# Implement queries (Create, FindByID, GetAll, Update, Delete)
```

3. **Create request validation**
```bash
# Copy app/Http/Requests/EXAMPLE_REQUEST_TEMPLATE.go
# Implement validation rules
```

4. **Create controller**
```bash
# Copy app/Http/Controllers/CONTROLLER_TEMPLATE.go
# Implement CRUD operations
```

5. **Register routes**
```bash
# Edit routes/api.go
# Add your controller routes
```

6. **Test**
```bash
go run main.go
# Test with curl or Postman
```

## 🔐 Authentication

### Login
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

Response:
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

### Using Token
```bash
GET /api/v1/users
Authorization: Bearer {access_token}
```

## 📋 Environment Configuration

```env
# Application
APP_NAME=My API
APP_ENV=local
APP_PORT=8080

# Database
DB_CONNECTION=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=my_database
DB_USERNAME=root
DB_PASSWORD=secret

# JWT
JWT_SECRET=your-super-secret-key
JWT_ACCESS_TOKEN_EXPIRATION=15
JWT_REFRESH_TOKEN_EXPIRATION=10080

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_MAX_REQUESTS=100
RATE_LIMIT_WINDOW_SECONDS=60

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS

# Redis (Optional)
REDIS_HOST=localhost
REDIS_PORT=6379
```

## 🧹 Cleanup & Reset

To reset to clean starter kit state (removing all implementations):

```bash
chmod +x clean_starter.sh
./clean_starter.sh
```

## 📦 Git Setup

To initialize as a new git repository:

```bash
chmod +x init_starter_git.sh
./init_starter_git.sh
```

## 🧪 Testing

### Test with curl

```bash
# Health check
curl http://localhost:8080/health

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"pass"}'

# Get list (requires token)
TOKEN="your-access-token-here"
curl -X GET "http://localhost:8080/api/v1/your-resource?page=1&per_page=10" \
  -H "Authorization: Bearer $TOKEN"
```

## 🎯 Template Files

Ready-to-use templates for rapid development:

- `app/Http/Controllers/CONTROLLER_TEMPLATE.go` - CRUD controller template
- `app/Models/MODEL_TEMPLATE.go` - Database model template
- `app/Http/Requests/EXAMPLE_REQUEST_TEMPLATE.go` - Request validation template
- `database/migrations/MIGRATION_TEMPLATE.sql` - Database migration template

Copy and customize these templates for your own resources.

## 🏗️ Architecture

### Request Flow
```
Request
  ↓
[CORS Middleware]
  ↓
[Rate Limiter]
  ↓
[Logger Middleware]
  ↓
[JWT Auth Middleware] ← If protected route
  ↓
[Controller]
  ↓
[Validation]
  ↓
[Database Query]
  ↓
[Response Formatter]
  ↓
Response
```

### Response Format

**Success (200)**
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {...}
}
```

**Error (4xx/5xx)**
```json
{
  "success": false,
  "error": "error_code",
  "message": "Human readable message",
  "details": null
}
```

**Validation Error (422)**
```json
{
  "success": false,
  "error": "validation_failed",
  "message": "Validation failed",
  "details": {
    "field_name": "validation error message"
  }
}
```

## 📝 Best Practices Implemented

1. **Separation of Concerns** - Controllers, models, requests, middleware clearly separated
2. **Error Handling** - Structured error responses with proper HTTP status codes
3. **Security** - JWT authentication, password hashing, CORS, rate limiting
4. **Scalability** - Connection pooling, pagination, Redis ready
5. **Maintainability** - Clear structure, reusable helpers, comprehensive documentation
6. **Configuration** - Environment-based config from .env
7. **Logging** - Request logging middleware for debugging
8. **Validation** - Binding tags for request validation

## 🚀 Production Deployment

Before deploying to production:

1. Update `.env` with production values
2. Set `APP_ENV=production` (enables Gin release mode)
3. Use strong `JWT_SECRET` (generate with: `openssl rand -base64 32`)
4. Setup database backups
5. Configure proper CORS origins
6. Setup logging & monitoring
7. Use HTTPS
8. Rate limiting appropriate to your use case

## 🐛 Troubleshooting

### Database Connection Error
- Verify MySQL/PostgreSQL is running
- Check database credentials in `.env`
- Verify database user has proper permissions

### JWT Token Invalid
- Verify `JWT_SECRET` matches between token generation and validation
- Check token format: `Authorization: Bearer {token}`
- Verify token hasn't expired

### CORS Error
- Check `CORS_ALLOWED_ORIGINS` in `.env`
- Ensure frontend domain is in the allowlist

### Rate Limit Error (429)
- Wait or adjust `RATE_LIMIT_MAX_REQUESTS` in `.env`

## 📖 Additional Resources

- [Go Language Docs](https://golang.org/doc/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [sqlx Documentation](https://github.com/jmoiron/sqlx)
- [JWT Go](https://github.com/golang-jwt/jwt)

## 📄 License

This starter kit is provided as-is for development use.

## 🤝 Contributing

Feel free to customize and extend this starter kit for your needs. Follow the patterns established in the template files.

## 📞 Support

Refer to the documentation files:
- `BACKEND_API_GUIDE.md` - Complete API architecture
- `STARTER_KIT_SETUP.md` - Setup & development workflow
- Template files - Code examples and patterns

---

**Ready to build your next Go backend? Start coding!** 🚀
