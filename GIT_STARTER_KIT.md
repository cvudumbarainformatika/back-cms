# Go Gin Backend Starter Kit - Complete for Git

Ini adalah dokumentasi lengkap untuk menggunakan folder ini sebagai Go Backend starter kit yang akan disimpan di git.

## 📋 File Structure - Ready for Git

```
.
├── .env.example                    # Environment template
├── .gitattributes                  # Git line endings
├── .gitignore                      # Git exclusions
├── README.md                       # Project overview
├── CHANGELOG.md                    # Version history
├── DEPLOYMENT.md                   # Deployment guide
├── CONTRIBUTING.md                 # Contributing guidelines
├── Makefile                        # Development tasks
├── clean_starter.sh                # Reset to clean state
├── init_starter_git.sh             # Initialize git repo
│
├── app/
│   ├── Exceptions/                 # Error handling
│   ├── Http/
│   │   ├── Controllers/            # ADD YOUR CONTROLLERS
│   │   │   └── CONTROLLER_TEMPLATE.go
│   │   ├── Middleware/             # (Production-ready)
│   │   └── Requests/               # ADD YOUR REQUESTS
│   │       └── EXAMPLE_REQUEST_TEMPLATE.go
│   └── Models/                     # ADD YOUR MODELS
│       └── MODEL_TEMPLATE.go
│
├── bootstrap/                      # App initialization
├── config/                         # Configuration
├── database/
│   ├── database.go                 # (Production-ready)
│   ├── redis.go                    # (Production-ready)
│   ├── migrations/                 # ADD YOUR MIGRATIONS
│   │   └── MIGRATION_TEMPLATE.sql
│   └── seeders/                    # ADD YOUR SEEDERS
│
├── routes/                         # Route definitions
├── utils/                          # (Production-ready)
│
├── BACKEND_API_GUIDE.md            # API documentation
├── STARTER_KIT_SETUP.md            # Setup guide
├── STARTER_KIT_SUMMARY.md          # Quick reference
└── GIT_STARTER_KIT.md              # This file
```

## 🚀 Quick Start for New Projects

### Step 1: Clone This Starter Kit
```bash
git clone <your-starter-kit-repo> my-new-project
cd my-new-project
```

### Step 2: Setup Environment
```bash
cp .env.example .env
# Edit .env with your database credentials
```

### Step 3: Install & Run
```bash
make install
make run
```

### Step 4: Start Building
```bash
# Copy templates for your first API
cp app/Http/Controllers/CONTROLLER_TEMPLATE.go app/Http/Controllers/product_controller.go
cp app/Models/MODEL_TEMPLATE.go app/Models/product.go
# ... implement your API
```

## 📚 Documentation Files

| File | Purpose |
|------|---------|
| `README.md` | Project overview & quick start |
| `BACKEND_API_GUIDE.md` | Complete API architecture & patterns (15 KB) |
| `STARTER_KIT_SETUP.md` | Step-by-step setup & development guide (10 KB) |
| `STARTER_KIT_SUMMARY.md` | Quick reference & patterns |
| `CHANGELOG.md` | Version history & features |
| `DEPLOYMENT.md` | Production deployment guide |
| `CONTRIBUTING.md` | Contributing guidelines |
| `GIT_STARTER_KIT.md` | This file - Git setup guide |

## 🔧 Utility Scripts

### `./init_starter_git.sh`
Initialize git repository for this starter kit:
```bash
./init_starter_git.sh
```

What it does:
- Initialize git repository
- Ask for your name & email
- Add all files
- Create initial commit
- Show instructions for remote setup

### `./clean_starter.sh`
Reset to clean starter kit state (after development):
```bash
./clean_starter.sh
```

What it removes:
- All custom controllers
- All custom models
- All custom request files
- All database migrations
- All database seeders
- Custom routes

What it keeps:
- Middleware stack
- Core infrastructure
- Template files
- Documentation

## 📦 Git Setup Options

### Option 1: Using Script (Recommended)
```bash
./init_starter_git.sh
```

Then:
```bash
git remote add origin <your-repo-url>
git branch -M main
git push -u origin main
```

### Option 2: Manual Setup
```bash
git init
git config user.name "Your Name"
git config user.email "your@email.com"
git add .
git commit -m "Initial commit: Go Gin Backend Starter Kit"
git remote add origin <your-repo-url>
git branch -M main
git push -u origin main
```

## 🎯 Template Files - Copy & Customize

### 1. Controller Template
```bash
cp app/Http/Controllers/CONTROLLER_TEMPLATE.go \
   app/Http/Controllers/your_controller.go
```

### 2. Model Template
```bash
cp app/Models/MODEL_TEMPLATE.go \
   app/Models/your_model.go
```

### 3. Request Validation Template
```bash
cp app/Http/Requests/EXAMPLE_REQUEST_TEMPLATE.go \
   app/Http/Requests/create_your_request.go
```

### 4. Migration Template
```bash
cp database/migrations/MIGRATION_TEMPLATE.sql \
   database/migrations/001_create_your_table.sql
```

## 🏗️ Development Workflow

### Create a New API
1. **Migration** → `database/migrations/001_your_table.sql`
2. **Model** → `app/Models/your_model.go`
3. **Request** → `app/Http/Requests/create_your_request.go`
4. **Controller** → `app/Http/Controllers/your_controller.go`
5. **Routes** → Update `routes/api.go`
6. **Test** → `curl http://localhost:8080/api/v1/your-endpoint`

### Example Commands
```bash
# Development
make install                # Install dependencies
make run                    # Run app
make build                  # Build executable
make test                   # Run tests
make fmt                    # Format code

# Git setup
make setup-git              # Initialize git
make reset-starter          # Reset to clean state

# Docker
make docker-build           # Build image
make docker-run             # Start containers
make docker-stop            # Stop containers
```

## 📖 Key Documentation

### For Understanding Architecture
→ **BACKEND_API_GUIDE.md**
- Complete API reference
- Authentication flow
- CRUD patterns
- Middleware explanation
- Pagination helpers
- Response formatting

### For Getting Started
→ **STARTER_KIT_SETUP.md**
- Step-by-step setup
- Development workflow
- Creating first API
- Best practices
- Troubleshooting

### For Quick Reference
→ **STARTER_KIT_SUMMARY.md**
- Architecture overview
- Key patterns
- Checklist for new projects
- Quick patterns

### For Production
→ **DEPLOYMENT.md**
- Pre-deployment checklist
- Docker deployment
- Cloud platforms (Heroku, AWS, DigitalOcean)
- SSL/TLS setup
- Monitoring & logging

## ✨ What's Included

### Production-Ready Infrastructure
- ✅ JWT Authentication (access & refresh tokens)
- ✅ Middleware Stack (JWT, CORS, Rate Limiter, Logger, Error Handler)
- ✅ Database Abstraction (MySQL & PostgreSQL with sqlx)
- ✅ Redis Integration (ready to use)
- ✅ Configuration Management (from .env)
- ✅ Graceful Shutdown (proper cleanup)

### Request/Response System
- ✅ Request Validation (binding tags)
- ✅ Response Formatting (success, error, validation)
- ✅ Pagination Helpers (offset, cursor, Laravel-style)
- ✅ Error Handling (global error handler)

### Development Tools
- ✅ Template Files (controller, model, request, migration)
- ✅ Comprehensive Documentation (7 guides)
- ✅ Utility Scripts (git setup, cleanup)
- ✅ Makefile (common tasks)
- ✅ Git Configuration (.gitignore, .gitattributes)

## 🔐 Environment Configuration

Copy `.env.example` to `.env` and update:

```env
APP_NAME=My API
APP_ENV=local
APP_PORT=8080

DB_CONNECTION=mysql
DB_HOST=localhost
DB_DATABASE=my_database
DB_USERNAME=root
DB_PASSWORD=secret

JWT_SECRET=generate-with-openssl-rand-base64-32
```

For production, update these values appropriately.

## 📝 Git Best Practices

### Commit Messages
```
feat(auth): add login endpoint
fix(user): correct validation error
docs(readme): update setup instructions
style(format): format code
refactor(api): improve error handling
test(user): add user model tests
chore(deps): update dependencies
```

### Workflow
```bash
# Create feature branch
git checkout -b feat/your-feature

# Make changes & commit
git add .
git commit -m "feat: your feature description"

# Push to remote
git push origin feat/your-feature

# Create Pull Request on GitHub/GitLab
```

## 🧪 Testing Your Setup

### Test Health Endpoint
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

### Create Test API
```bash
# Copy templates
cp app/Http/Controllers/CONTROLLER_TEMPLATE.go app/Http/Controllers/test_controller.go
cp app/Models/MODEL_TEMPLATE.go app/Models/test_model.go

# Edit files (change struct names, etc.)

# Register in routes/api.go

# Test
go run main.go
```

## 🚀 Deployment

See `DEPLOYMENT.md` for:
- Pre-deployment checklist
- Docker deployment
- Cloud platforms
- SSL/TLS setup
- Monitoring
- Troubleshooting

## 📞 Support

Refer to documentation:
1. **Setup Issues** → See STARTER_KIT_SETUP.md
2. **Architecture Questions** → See BACKEND_API_GUIDE.md
3. **Code Patterns** → See STARTER_KIT_SUMMARY.md
4. **Deployment** → See DEPLOYMENT.md

## ✅ Next Steps

1. **Initialize Git** (if not done):
   ```bash
   ./init_starter_git.sh
   ```

2. **Push to GitHub/GitLab**:
   ```bash
   git remote add origin <your-repo-url>
   git push -u origin main
   ```

3. **Your starter kit is ready!**
   - Clone for new projects: `git clone <repo> <project-name>`
   - Start building: Copy templates & customize
   - Deploy: Follow DEPLOYMENT.md

4. **Share with your team** or use for all Go backend projects!

---

**Version:** 1.0.0  
**Status:** Production Ready  
**Last Updated:** 2024-12-29

---

**Happy coding! 🚀**
