#!/bin/bash

################################################################################
# Clean Starter Kit Script
# 
# Usage: ./clean_starter.sh
#
# This script removes all custom implementations and resets the backend to
# a clean starter kit state, keeping only the core infrastructure.
#
# It will remove:
#   - All controllers (app/Http/Controllers/*.go except templates)
#   - All models (app/Models/*.go except templates)
#   - All request files (app/Http/Requests/*.go except login/register)
#   - All database migrations
#   - All database seeders
#   - All custom implementation in routes/api.go
#
# What it keeps:
#   - Middleware stack (JWT, CORS, Rate Limiter, Logger, Error Handler)
#   - Template files (for reference)
#   - Documentation files
#   - Configuration & utilities
#   - Core infrastructure
#
################################################################################

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if script is run from project root
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ Error: go.mod not found. Please run this script from project root directory.${NC}"
    exit 1
fi

echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║           🧹 GO GIN STARTER KIT - CLEANUP SCRIPT              ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Confirmation
echo -e "${YELLOW}⚠️  This will delete ALL custom implementations!${NC}"
echo ""
echo "Files that will be DELETED:"
echo "  • All controllers (except templates)"
echo "  • All models (except templates)"
echo "  • All request files (except login/register templates)"
echo "  • All database migrations"
echo "  • All database seeders"
echo "  • Custom routes setup"
echo ""
echo "Files that will be KEPT:"
echo "  • Middleware stack"
echo "  • Template files"
echo "  • Documentation"
echo "  • Core infrastructure"
echo ""
read -p "Are you sure? Type 'yes' to continue: " -r
echo
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo -e "${YELLOW}Cancelled.${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}Starting cleanup...${NC}"
echo ""

# Counter for deleted files
DELETED_COUNT=0

# Function to safely delete file
delete_file() {
    if [ -f "$1" ]; then
        rm "$1"
        echo -e "${GREEN}✓${NC} Deleted: $1"
        ((DELETED_COUNT++))
    fi
}

# Function to safely delete files matching pattern
delete_pattern() {
    local pattern=$1
    local description=$2
    local count=0
    
    while IFS= read -r -d '' file; do
        delete_file "$file"
        ((count++))
    done < <(find . -name "$pattern" -type f -print0 2>/dev/null)
    
    if [ $count -gt 0 ]; then
        echo -e "${BLUE}  ($description: $count files)${NC}"
    fi
}

# Delete Controllers (except templates)
echo -e "${BLUE}Deleting controllers...${NC}"
find app/Http/Controllers -name "*.go" -type f ! -name "CONTROLLER_TEMPLATE.go" -delete
echo -e "${GREEN}✓${NC} Deleted all controllers (except template)"
((DELETED_COUNT++))

# Delete Models (except templates)
echo -e "${BLUE}Deleting models...${NC}"
find app/Models -name "*.go" -type f ! -name "MODEL_TEMPLATE.go" -delete
echo -e "${GREEN}✓${NC} Deleted all models (except template)"
((DELETED_COUNT++))

# Delete Request files (except templates and login/register)
echo -e "${BLUE}Deleting request files...${NC}"
find app/Http/Requests -name "*.go" -type f \
    ! -name "EXAMPLE_REQUEST_TEMPLATE.go" \
    ! -name "login_request.go" \
    ! -name "register_request.go" \
    -delete
# Remove subdirectories if they exist
[ -d "app/Http/Requests/User" ] && rm -rf "app/Http/Requests/User"
[ -d "app/Http/Requests/DataWisata" ] && rm -rf "app/Http/Requests/DataWisata"
[ -d "app/Http/Requests/DataWisataUpdateLog" ] && rm -rf "app/Http/Requests/DataWisataUpdateLog"
[ -d "app/Http/Requests/Wisata" ] && rm -rf "app/Http/Requests/Wisata"
echo -e "${GREEN}✓${NC} Deleted all request files (except templates and login/register)"
((DELETED_COUNT++))

# Delete Migrations
echo -e "${BLUE}Deleting migrations...${NC}"
find database/migrations -name "*.sql" -type f ! -name "MIGRATION_TEMPLATE.sql" -delete
echo -e "${GREEN}✓${NC} Deleted all migrations (except template)"
((DELETED_COUNT++))

# Delete Seeders
echo -e "${BLUE}Deleting seeders...${NC}"
find database/seeders -name "*.go" -type f -delete
echo -e "${GREEN}✓${NC} Deleted all seeders"
((DELETED_COUNT++))

# Reset routes/api.go to clean state
echo -e "${BLUE}Resetting routes/api.go...${NC}"
cat > routes/api.go << 'EOF'
package routes

import (
	middleware "github.com/cvudumbarainformatika/backend/app/Http/Middleware"
	"github.com/cvudumbarainformatika/backend/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, db *sqlx.DB, redis *redis.Client, cfg *config.Config) {
	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// ==============================
		// Public Routes (No Auth Required)
		// ==============================
		// TODO: Add auth routes here (login, register)
		// Example:
		// auth := v1.Group("/auth")
		// {
		//     auth.POST("/register", authController.Register)
		//     auth.POST("/login", authController.Login)
		// }

		// ==============================
		// Protected Routes (JWT Required)
		// ==============================
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuthMiddleware(cfg.JWT.Secret))
		{
			// TODO: Add your protected routes here
			// Example: User CRUD routes
			// users := protected.Group("/users")
			// {
			//     users.GET("/get-list", userController.GetAllUsers)
			//     users.GET("/:id", userController.GetUserByID)
			//     users.POST("/create", userController.CreateUser)
			//     users.PUT("/update/:id", userController.UpdateUser)
			//     users.DELETE("/delete/:id", userController.DeleteUser)
			// }
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "ok",
			"database": "connected",
		})
	})
}
EOF
echo -e "${GREEN}✓${NC} Reset routes/api.go to clean state"
((DELETED_COUNT++))

# Clean up .env if exists (optional - user may want to keep it)
if [ -f ".env" ]; then
    echo ""
    read -p "Do you want to keep .env file? (y/n) " -r
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        rm ".env"
        echo -e "${GREEN}✓${NC} Deleted .env"
        ((DELETED_COUNT++))
    else
        echo -e "${YELLOW}ⓘ${NC} Kept .env file"
    fi
fi

# Summary
echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}✓ Cleanup Complete!${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo "Total files deleted/reset: $DELETED_COUNT"
echo ""
echo -e "${GREEN}Your starter kit is ready for a new project!${NC}"
echo ""
echo "Next steps:"
echo "  1. Copy .env.example to .env: ${BLUE}cp .env.example .env${NC}"
echo "  2. Update .env with your database credentials"
echo "  3. Read STARTER_KIT_SETUP.md for development guide"
echo "  4. Copy template files and start building!"
echo ""
echo "Templates available in:"
echo "  • app/Http/Controllers/CONTROLLER_TEMPLATE.go"
echo "  • app/Models/MODEL_TEMPLATE.go"
echo "  • app/Http/Requests/EXAMPLE_REQUEST_TEMPLATE.go"
echo "  • database/migrations/MIGRATION_TEMPLATE.sql"
echo ""
