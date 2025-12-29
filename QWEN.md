# Project Overview: Juwita Tourism Information System - Backend

## Project Summary
This is a backend RESTful API server for the Juwita tourism information system built with Go using the Gin framework. It's designed as a standalone API server with JWT authentication, MySQL/PostgreSQL support, and Docker integration.

The backend follows the API starter kit template found at `github.com/yourusername/go-gin-restful-api-starter-kit`, which implements a Laravel-inspired project structure with clean separation of concerns.

## Technologies & Dependencies
- **Language**: Go 1.23+
- **Web Framework**: Gin (v1.10.1)
- **Database**: sqlx (v1.4.0) with MySQL (go-sql-driver/mysql v1.9.3) and PostgreSQL (lib/pq v1.10.9) support
- **Authentication**: JWT (golang-jwt/jwt/v5 v5.3.0) with bcrypt password hashing
- **Configuration**: godotenv (v1.5.1) for environment variable management
- **Containerization**: Docker & Docker Compose with multi-stage builds

## Project Architecture
```
.
├── main.go                          # Application entry point
├── .env                             # Environment configuration
├── .env.example                     # Environment template
│
├── app/
│   ├── Models/                      # Data models (users.go, post.go)
│   ├── Http/
│   │   ├── Controllers/             # Request handlers (auth_controller.go, post_controller.go)
│   │   ├── Requests/                # Request validation structures
│   │   └── Middleware/              # HTTP middleware (CORS, rate limiter, JWT auth)
│   └── Exceptions/                  # Error handling
│
├── bootstrap/                       # Application initialization
├── config/                          # Configuration loader
├── database/                        # Database connection and migrations
│   ├── migrations/                  # SQL migration files
│   └── seeders/                     # Seed data
├── routes/                          # API route definitions
├── storage/logs/                    # Application logs
└── utils/                           # Utility functions (JWT, responses, pagination)
```

## Key Features
- JWT-based authentication with access/refresh tokens
- Rate limiting per IP address
- CORS configuration
- Request validation
- Pagination (both offset-based and cursor-based)
- MySQL and PostgreSQL support
- Database connection pooling with configurable settings
- Docker support (development and production)
- Hot reload during development (using Air)
- Comprehensive error handling

## Environment Configuration
Configuration is managed through environment variables loaded from a `.env` file:

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Application name | Go Gin Starter Kit |
| `APP_ENV` | Environment (local, production) | local |
| `APP_PORT` | Server port | 8080 |
| `DB_CONNECTION` | Database driver | mysql |
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 3306 |
| `DB_DATABASE` | Database name | starter_kit |
| `DB_USERNAME` | Database username | starter |
| `DB_PASSWORD` | Database password | secret |
| `JWT_SECRET` | JWT signing secret | dev-secret-key-change-in-production |
| `JWT_ACCESS_TOKEN_EXPIRATION` | Access token expiration (minutes) | 15 |
| `JWT_REFRESH_TOKEN_EXPIRATION` | Refresh token expiration (minutes) | 10080 |
| `RATE_LIMIT_ENABLED` | Enable rate limiting | true |
| `CORS_ALLOWED_ORIGINS` | Allowed CORS origins | http://localhost:3000 |

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login

### Posts (Protected)
- `GET /api/v1/posts` - Get posts with pagination
- `POST /api/v1/posts` - Create a new post

### Health Check
- `GET /health` - Health check endpoint

## Database Schema
The application uses two main tables:

**Users Table**: Stores user information (id, name, email, password, timestamps)
**Posts Table**: Stores posts associated with users (id, user_id, title, content, timestamps)

Migration files are located in `database/migrations/`.

## Development Commands
- `make dev` - Start development environment with Docker
- `make dev-logs` - View development logs
- `make dev-down` - Stop development environment
- `make prod` - Start production environment
- `make prod-logs` - View production logs
- `make prod-down` - Stop production environment
- `make build` - Build application locally
- `make test` - Run tests
- `make db-migrate` - Run migrations in development
- `make db-migrate-prod` - Run migrations in production

## Security Considerations
- Passwords are hashed using bcrypt with cost factor 10
- JWT tokens signed with HS256
- Access tokens short-lived (15 minutes default)
- Refresh tokens long-lived (7 days default)
- Rate limiting prevents brute force attacks
- Parameterized SQL queries prevent injection
- Configurable CORS per environment

## Docker Configuration
- Development environment uses `docker-compose.yml` with hot reload
- Production environment uses `docker-compose.prod.yml`
- Development builds from `Dockerfile.dev`
- Production builds from `Dockerfile` (multi-stage for minimal image)

## Testing
Run tests using: `go test ./...` or `make test`

## Deployment
For production deployments, ensure to:
1. Create a `.env` file with production values
2. Update JWT_SECRET with a strong random key
3. Configure CORS_ALLOWED_ORIGINS appropriately
4. Run `make prod` to start services
5. Manually run database migrations with `make db-migrate-prod`

## Project Context
This backend is part of the Juwita Tourism Information System (Juwita - Sistem Informasi Kunjungan Wisatawan Mancanegara dan Dalam Negeri), which also includes a native PHP component in a sibling directory. The Go backend provides a modern API layer while the PHP component serves as a legacy management interface. Both components interact with the same database.