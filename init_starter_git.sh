#!/bin/bash

################################################################################
# Initialize Starter Kit as Git Repository
#
# Usage: ./init_starter_git.sh
#
# This script initializes the current directory as a new Git repository
# configured as a Go Backend Starter Kit with all necessary setup.
#
# It will:
#   - Initialize git repository
#   - Configure git settings (user.name, user.email)
#   - Create initial commits
#   - Setup git hooks (optional)
#
################################################################################

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if already a git repo
if [ -d ".git" ]; then
    echo -e "${YELLOW}⚠️  This directory is already a Git repository.${NC}"
    read -p "Do you want to reinitialize? (y/n) " -r
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}Cancelled.${NC}"
        exit 1
    fi
    rm -rf .git
fi

echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║        📦 GO GIN STARTER KIT - GIT INITIALIZATION              ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Git user configuration
echo -e "${BLUE}Git Configuration:${NC}"
echo ""

read -p "Git user name (leave empty to skip): " GIT_USER_NAME
read -p "Git user email (leave empty to skip): " GIT_USER_EMAIL

echo ""

# Initialize repository
echo -e "${BLUE}Initializing Git repository...${NC}"
git init
echo -e "${GREEN}✓${NC} Git repository initialized"

# Configure user if provided
if [ -n "$GIT_USER_NAME" ]; then
    git config user.name "$GIT_USER_NAME"
    echo -e "${GREEN}✓${NC} Git user name set to: $GIT_USER_NAME"
fi

if [ -n "$GIT_USER_EMAIL" ]; then
    git config user.email "$GIT_USER_EMAIL"
    echo -e "${GREEN}✓${NC} Git user email set to: $GIT_USER_EMAIL"
fi

echo ""
echo -e "${BLUE}Setting up Git configuration...${NC}"

# Configure helpful settings
git config core.safecrlf true
git config core.ignorecase false
git config pull.rebase false

echo -e "${GREEN}✓${NC} Git configuration set"

echo ""
echo -e "${BLUE}Adding files to Git...${NC}"

# Add all files
git add .
echo -e "${GREEN}✓${NC} All files added"

# Create initial commit
echo ""
git commit -m "Initial commit: Go Gin Backend Starter Kit

- Clean infrastructure ready for new projects
- Authentication system with JWT
- Middleware stack (JWT, CORS, Rate Limiter, Logger, Error Handler)
- Database abstraction with sqlx (MySQL & PostgreSQL)
- Request/Response helpers
- Pagination utilities
- Template files for rapid development
- Comprehensive documentation"

echo -e "${GREEN}✓${NC} Initial commit created"

echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}✓ Git Repository Initialized!${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Show git info
echo "Current Git status:"
git log --oneline -1
echo ""
echo "Repository information:"
echo "  Branch: $(git rev-parse --abbrev-ref HEAD)"
echo "  Total commits: $(git rev-list --all --count)"
echo ""

# Remote setup instructions
echo -e "${YELLOW}Next steps:${NC}"
echo ""
echo "1. Add remote repository:"
echo "   ${BLUE}git remote add origin <your-repository-url>${NC}"
echo ""
echo "2. Push to remote:"
echo "   ${BLUE}git branch -M main${NC}"
echo "   ${BLUE}git push -u origin main${NC}"
echo ""
echo "3. For new projects, use:"
echo "   ${BLUE}./clean_starter.sh${NC} to reset to clean state"
echo ""
echo "📚 Documentation files:"
echo "   • BACKEND_API_GUIDE.md - Complete API reference"
echo "   • STARTER_KIT_SETUP.md - Setup & development guide"
echo "   • STARTER_KIT_SUMMARY.md - Quick reference"
echo ""
